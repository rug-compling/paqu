package main

import (
	"fmt"
	"html"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// TAB: begin
func home(q *Context) {

	prefix := getprefix(q)
	if !q.prefixes[prefix] {
		http.Error(q.w, "Invalid corpus: "+prefix, http.StatusPreconditionFailed)
		return
	}

	// HTML-uitvoer van begin van de pagina
	writeHead(q, "", 1)
	html_header(q)

	// HTML-uitvoer van het formulier
	// Returnwaarde is true als er een query was gedefinieerd
	has_query := html_form(q)

	// Als er geen query is gedefinieerd, HTML-uitvoer van korte helptekst, pagina-einde, en exit
	if !has_query {
		html_uitleg(q)
		html_footer(q)
		return
	}

	var chClose <-chan bool
	if f, ok := q.w.(http.CloseNotifier); ok {
		chClose = f.CloseNotify()
	} else {
		chClose = make(<-chan bool)
	}

	_, err := q.db.Exec(fmt.Sprintf("UPDATE `%s_info` SET `active` = NOW() WHERE `id` = %q", Cfg.Prefix, prefix))
	if doErr(q, err) {
		return
	}

	offset := 0
	o, err := strconv.Atoi(first(q.r, "offset"))
	if err == nil {
		offset = o
	}

	fmt.Fprintln(q.w, "<hr>")

	query, err := makeQuery(q, prefix, chClose)
	if doErr(q, err) {
		return
	}

	// DEBUG: HTML-uitvoer van de query
	fmt.Fprint(q.w, "<div style=\"font-family:monospace\">\n", html.EscapeString(query), "\n</div><p>\n")

	fmt.Fprint(q.w, "<div id=\"busy\"><img src=\"busy.gif\" alt=\"aan het werk...\"></div>\n")

	if ff, ok := q.w.(http.Flusher); ok {
		ff.Flush()
	}

	now := time.Now()

	// Maximaal ZINMAX matchende xml-bestanden opvragen
	select {
	case <-chClose:
		logerr(ConnectionClosed)
		return
	default:
	}
	rows, err := timeoutQuery(q, chClose,
		"SELECT `arch`,`file` FROM `"+Cfg.Prefix+"_c_"+prefix+"_deprel` WHERE "+query+" GROUP BY `arch`,`file` LIMIT "+fmt.Sprint(offset)+", "+fmt.Sprint(ZINMAX))
	if doErr(q, err) {
		busyClear(q)
		return
	}
	zinnen := make([]*Sentence, 0, ZINMAX)
	var a, f int
	for rows.Next() {
		err := rows.Scan(&a, &f)
		if doErr(q, err) {
			busyClear(q)
			return
		}
		zinnen = append(zinnen, &Sentence{arch: a, file: f, items: make([]Row, 0)})
	}
	err = rows.Err()
	if doErr(q, err) {
		busyClear(q)
		return
	}

	// Gegevens bij gevonden xml-bestanden opvragen
	for _, zin := range zinnen {
		// Zin bij xml-bestand opvragen
		var s string
		select {
		case <-chClose:
			logerr(ConnectionClosed)
			return
		default:
		}
		rows, err := q.db.Query(fmt.Sprintf("SELECT `sent` FROM `%s_c_%s_sent` WHERE `arch` = %d AND `file`= %d",
			Cfg.Prefix, prefix, zin.arch, zin.file))
		if doErr(q, err) {
			busyClear(q)
			return
		}
		if rows.Next() {
			err := rows.Scan(&s)
			if doErr(q, err) {
				busyClear(q)
				return
			}
			zin.words = strings.Fields(s)
			rows.Close()
		} else {
			fmt.Printf("Missing sentence for file id %d\n", zin.file)
			os.Exit(1)
		}

		// Matchende dependency relations bij xml-bestand opvragen
		select {
		case <-chClose:
			logerr(ConnectionClosed)
			return
		default:
		}
		rows, err = q.db.Query(fmt.Sprintf(
			"SELECT `word`,`lemma`,`postag`,`rel`,`hpostag`,`hlemma`,`hword`,`begin`,`end`,`hbegin`,`hend`,`mark` FROM `%s_c_%s_deprel` WHERE `arch` = %d AND `file`= %d AND %s ORDER BY `begin`,`hbegin`,`rel`",
			Cfg.Prefix, prefix, zin.arch, zin.file, query))
		if doErr(q, err) {
			busyClear(q)
			return
		}
		for rows.Next() {
			var r Row
			err := rows.Scan(
				&r.word,
				&r.lemma,
				&r.postag,
				&r.rel,
				&r.hpostag,
				&r.hlemma,
				&r.hword,
				&r.begin,
				&r.end,
				&r.hbegin,
				&r.hend,
				&r.mark)
			if doErr(q, err) {
				busyClear(q)
				return
			}
			zin.items = append(zin.items, r)
		}
		err = rows.Err()
		if doErr(q, err) {
			busyClear(q)
			return
		}
	}

	busyClear(q)
	// Verwerking en HTML-uitvoer van zinnen en dependency relations
	fmt.Fprintln(q.w, "<ol>")
	for i, zin := range zinnen {

		// Nodes die gemarkeerd moeten worden in de boom, voor link naar 'lassytree'
		mark := make(map[string]bool)

		// Begin zin
		fmt.Fprintf(q.w, "<li value=\"%d\" class=\"li1\">", i+offset+1)

		// Posities verzamelen waar een markering moet beginnen / eindigen
		idx := make([]int, len(zin.words))
		for _, item := range zin.items {
			idx[item.begin] |= 1 // markering begin gewoon woord
			idx[item.end-1] |= 4 // markering einde woord
			if item.hword != "" {
				idx[item.hbegin] |= 2 // markering begin hoofdwoord
				idx[item.hend-1] |= 4 // markering einde woord
			}
		}

		// Posities van woorden en hoofdwoorden in de boom, voor link naar 'lassytree'
		yellows := make([]string, 0, len(zin.words))
		greens := make([]string, 0, len(zin.words))

		// Alle woorden in de zin verwerken
		for j, word := range zin.words {

			// Voor link naar 'lassytree'
			if idx[j]&1 != 0 {
				yellows = append(yellows, fmt.Sprint(j))
				get_path(zin, j, mark) // path van woord naar hoofdwoord toevoegen aan map 'mark'
			}
			if idx[j]&2 != 0 {
				greens = append(greens, fmt.Sprint(j))
			}

			// HTML-uitvoer van woord in zin, incl begin/eind kleur
			if idx[j]&3 == 3 {
				fmt.Fprint(q.w, YELGRN)
			} else if idx[j]&1 != 0 {
				fmt.Fprint(q.w, YELLOW)
			} else if idx[j]&2 != 0 {
				fmt.Fprint(q.w, GREEN)
			}
			fmt.Fprint(q.w, word)
			if idx[j]&4 != 0 {
				fmt.Fprint(q.w, "</span>")
			}
			fmt.Fprint(q.w, " ")
		}

		// HTML-uitvoer van link naar 'lassytree'
		ms := make([]string, 0, len(mark))
		for m := range mark {
			ms = append(ms, m)
		}
		fmt.Fprintf(q.w, "\n<a href=\"/tree?db=%s&amp;arch=%d&amp;file=%d&amp;yl=%s&amp;gr=%s&amp;ms=%s\">&nbsp;&#9741;&nbsp;</a>\n<ul>\n",
			prefix,
			zin.arch,
			zin.file,
			strings.Join(yellows, ","),
			strings.Join(greens, ","),
			strings.Join(ms, ","))

		// HTML-uitvoer van alle matchende dependency relations voor huidige zin
		seen := make(map[string]bool)
		for _, item := range zin.items {
			s := fmt.Sprintf("%s:%s — %s — %s:%s", item.word, item.postag, item.rel, item.hword, item.hpostag)
			if seen[s] {
				continue
			}
			seen[s] = true
			qq := make_query_string(item.word, item.postag, item.rel, item.hpostag, item.hword, prefix)
			fmt.Fprintf(q.w, "<li class=\"li2\"><a href=\"/?%s\">%s</a>\n", qq, html.EscapeString(s))
		}

		// Einde zin + dependency relations
		fmt.Fprint(q.w, "</ul>\n\n")
	}
	fmt.Fprint(q.w, "</ol>\n<p>\n")

	// Links naar volgende en vorige pagina's met resultaten
	qs := make_query_string(
		first(q.r, "word"),
		first(q.r, "postag"),
		first(q.r, "rel"),
		first(q.r, "hpostag"),
		first(q.r, "hword"),
		first(q.r, "db"))
	if offset > 0 || len(zinnen) == ZINMAX {
		if offset > 0 {
			fmt.Fprintf(q.w, "<a href=\"/?%s&amp;offset=%d\">vorige</a>", qs, offset-ZINMAX)
		} else {
			fmt.Fprint(q.w, "vorige")
		}
		fmt.Fprint(q.w, " | ")
		if len(zinnen) == ZINMAX {
			fmt.Fprintf(q.w, "<a href=\"/?%s&amp;offset=%d\">volgende</a>", qs, offset+ZINMAX)
		} else {
			fmt.Fprint(q.w, "volgende")
		}
	}

	fmt.Fprintln(q.w, "<hr><small>tijd:", time.Now().Sub(now), "</small><hr>")

	// Links naar statistieken
	fmt.Fprintf(q.w, `<p>
		<div id="stats">
		<button onclick="javascript:$.fn.stats('stats?word=%s&amp;postag=%s&amp;rel=%s&amp;hpostag=%s&amp;hword=%s&amp;db=%s')">statistiek &mdash; algemeen</button>
		</div>
`,
		urlencode(first(q.r, "word")),
		urlencode(first(q.r, "postag")),
		urlencode(first(q.r, "rel")),
		urlencode(first(q.r, "hpostag")),
		urlencode(first(q.r, "hword")),
		urlencode(prefix))

	fmt.Fprintf(q.w, `<p>
		<div id="statsrel">
		<form action="javascript:$.fn.statsrel()" name="statsrelform" onsubmit="javascript:return statftest()">
		<input type="hidden" name="word" value="%s">
		<input type="hidden" name="postag" value="%s">
		<input type="hidden" name="rel" value="%s">
		<input type="hidden" name="hpostag" value="%s">
		<input type="hidden" name="hword" value="%s">
		<input type="hidden" name="db" value="%s">
		Selecteer twee of meer elementen om ze te koppelen:
		<p>
		<table>
		<tr>
		  <td style="background-color: yellow"><input type="checkbox" name="cword" value="1">woord
		  <td>
		  <td style="background-color: lightgreen"><input type="checkbox" name="chword" value="1">hoofdwoord
		<tr>
		  <td><input type="checkbox" name="clemma" value="1">lemma
		  <td><input type="checkbox" name="crel" value="1">relatie
		  <td><input type="checkbox" name="chlemma" value="1">lemma
		<tr>
		  <td><input type="checkbox" name="cpostag" value="1">postag
		  <td>
		  <td><input type="checkbox" name="chpostag" value="1">postag
		</table>
		<p>
		<input type="submit" value="statistiek &mdash; gerelateerd">
		</form>
		<p>
		<div id="statresults">
		</div>
		</div>
`,
		html.EscapeString(first(q.r, "word")),
		html.EscapeString(first(q.r, "postag")),
		html.EscapeString(first(q.r, "rel")),
		html.EscapeString(first(q.r, "hpostag")),
		html.EscapeString(first(q.r, "hword")),
		html.EscapeString(prefix))

	html_footer(q)

}

//. Hulpfuncties

func get_path(zin *Sentence, idx int, mark map[string]bool) {
	for _, item := range zin.items {
		if item.begin == idx {
			if item.mark != "" {
				for _, m := range strings.Split(item.mark, ",") {
					mark[m] = true
				}
			}
		}
	}
}

func make_query_string(word, postag, rel, hpostag, hword, db string) string {
	return fmt.Sprintf(
		"word=%s&amp;postag=%s&amp;rel=%s&amp;hpostag=%s&amp;hword=%s&amp;db=%s",
		urlencode(word),
		urlencode(postag),
		urlencode(rel),
		urlencode(hpostag),
		urlencode(hword),
		urlencode(db))
}

//. HTML

func html_header(q *Context) {
	fmt.Fprint(q.w, `
<script type="text/javascript" src="jquery.js"></script>
<script type="text/javascript"><!--
  $.fn.stats = function(url) {
    $("#stats").html('<img src="busy.gif" alt="aan het werk...">');
    $.get(url, function(data) {
      $("#stats").html(data);
    }).fail(function(e) {
      $("#stats").html(e.responseText);
    });
  }
  $.fn.statsrel = function() {
    $("#statresults").html('<img src="busy.gif">');
    $.get("statsrel?" + $(document.statsrelform).serialize(), function(data) {
      $("#statresults").html(data);
    }).fail(function(e) {
      $("#statresults").html(e.responseText);
    });
  }
  function formclear(f) {
    f.word.value = "";
    f.postag.value = "";
    f.rel.value = "";
    f.hpostag.value = "";
    f.hword.value = "";
  }
  function statftest() {
    var f = document.statsrelform;
    var n = 0;
    if (f.cword.checked   ) { n++; }
    if (f.clemma.checked  ) { n++; }
    if (f.cpostag.checked ) { n++; }
    if (f.crel.checked    ) { n++; }
    if (f.chword.checked  ) { n++; }
    if (f.chlemma.checked ) { n++; }
    if (f.chpostag.checked) { n++; }
    if (n < 2) {
      alert("Selecteer minimaal twee elementen");
      return false;
    }
    return true;
  }
  //--></script>
`)
}

func html_uitleg(q *Context) {
	fmt.Fprint(q.w, `
<p>
<hr>
<p>
Met deze toepassing kun je zoeken naar woord-paren in delen van de
Lassy treebanks. De Lassy treebanks bestaan uit Nederlandstalige
zinnen die voorzien zijn van hun syntactische ontleding. In de Lassy
Klein treebank (1 miljoen woorden) is voor elke zin de
syntactische ontleding handmatig gecheckt. In de Lassy Groot treebank
(700 miljoen woorden) is de syntactische ontleding automatisch
toegevoegd door de automatische parser Alpino.
<p>
<a href="info.html">Meer info...</a>
`)
}

func html_footer(q *Context) {
	fmt.Fprint(q.w, `
<p>
<hr>
<div class="foot">
mede mogelijk gemaakt door:
<p>
<a href="http://www.clarin.nl/" target="_blank"><img src="/clarinnl.png" class="noborder" alt="Clarin NL"></a>
</div>
</body>
</html>
`)
}

func html_form(q *Context) (has_query bool) {
	has_query = true
	if first(q.r, "word") == "" &&
		first(q.r, "postag") == "" &&
		first(q.r, "rel") == "" &&
		first(q.r, "hpostag") == "" &&
		first(q.r, "hword") == "" {
		has_query = false
	}

	fmt.Fprint(q.w, `
<form action="." method="get" accept-charset="utf-8">
<table>
<tr>
  <td colspan="3" style="padding-bottom:1em">corpus: <select name="db">
`)
	html_opts(q, q.opt_db, getprefix(q), "corpus")
	fmt.Fprintln(q.w, "</select>")
	if q.auth {
		fmt.Fprintln(q.w, "<a href=\"corpuslijst\">meer/minder</a>")
	}
	fmt.Fprintf(q.w, `
	   <tr>
		 <td style="background-color: yellow">woord
		 <td>
		 <td style="background-color: lightgreen">hoofdwoord
	   <tr>
		 <td><input type="text" name="word" size="20" value="%s">
	   `, html.EscapeString(first(q.r, "word")))
	fmt.Fprint(q.w, `
		 <td>
		   <select name="rel">
	   `)
	html_opts(q, opt_rel, first(q.r, "rel"), "relatie")
	fmt.Fprintf(q.w, `
		   </select>
		 <td><input type="text" name="hword" size="20" value="%s">
	   `, html.EscapeString(first(q.r, "hword")))
	fmt.Fprint(q.w, `
	   <tr>
		 <td>
		   <select name="postag" style="width: 100%">
	   `)
	html_opts(q, opt_postag, first(q.r, "postag"), "postag")
	fmt.Fprintf(q.w, `
		   </select>
		 <td>
		 <td>
		   <select name="hpostag" style="width:100%%" >
	   `)
	html_opts(q, opt_hpostag, first(q.r, "hpostag"), "postag")
	fmt.Fprint(q.w, `
		   </select>
	   <tr>
		 <td style="padding-top:1em">
		   <input type="button" value="help" onClick="javascript:window.open('info.html')">
		 <td colspan="2" class="right" style="padding-top:1em">
		   <input type="submit" value="Zoeken">
		   <input type="button" value="Wissen" onClick="javascript:formclear(form)">
		   <input type="reset" value="Reset">
	   </table>
	   </form>
	   `)

	return
}

func html_opts(q *Context, opts []string, value, title string) {
	var c byte
	for _, optstr := range opts {
		var opt, txt string
		p := strings.Fields(optstr)
		if len(p) == 1 {
			opt = optstr
		} else if len(p) > 1 {
			opt = p[0]
			txt = strings.Join(p[1:], " ")
		}
		if title == "relatie" && opt != "" {
			if opt[0] != c {
				c = opt[0]
				t := ""
				switch c {
				case 'A':
					t = "relaties met gewoon hoofd"
				case 'B':
					t = "andere relaties"
				case 'C':
					t = "enkelzijdige relaties"
				case 'D':
					t = "zeldzame relaties"
				}
				fmt.Fprintf(q.w, "    <optgroup label=\"&mdash; %s &mdash;\">\n", t)
			}
			opt = opt[1:]
		} else if title == "corpus" && opt != "" {
			if opt[0] != c {
				c = opt[0]
				t := ""
				switch c {
				case 'A':
					t = "algemene corpora"
				case 'B':
					t = "mijn corpora"
				case 'C':
					t = "corpora gedeeld door anderen"
				}
				if c != 'Z' {
					fmt.Fprintf(q.w, "    <optgroup label=\"&mdash; %s &mdash;\">\n", t)
				}
			}
			opt = opt[1:]
		}
		sel := ""
		if opt == value {
			sel = " selected=\"selected\""
		}
		if opt == "" {
			fmt.Fprintf(q.w, "    <option value=\"\"%s>&mdash; %s &mdash;</option>\n", sel, title)
		} else if txt == "" {
			fmt.Fprintf(q.w, "    <option%s>%s</option>\n", sel, opt)
		} else {
			fmt.Fprintf(q.w, "    <option value=\"%s\"%s>%s</option>\n", opt, sel, html.EscapeString(txt))
		}
	}
}

func getprefix(q *Context) string {
	var db string
	if p := first(q.r, "db"); p != "" {
		if p == "_default" {
			db = Cfg.Default
		} else {
			db = p
		}
	} else if prev, err := q.r.Cookie("paqu-prev"); err == nil {
		if p := prev.Value; q.prefixes[p] {
			db = p
		}
	}
	if db == "" {
		db = Cfg.Default
	}
	exp := time.Now().AddDate(0, 0, 14)
	http.SetCookie(q.w, &http.Cookie{Name: "paqu-prev", Value: db, Path: cookiepath, Expires: exp})
	return db
}

func busyClear(q *Context) {
	fmt.Fprint(q.w, `<script type="text/javascript"><!--
document.getElementById('busy').className = 'hide';
//--></script>
`)
}
