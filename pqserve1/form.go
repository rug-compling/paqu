package main

import (
	"errors"
	"fmt"
	"html"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//. Main

func form(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	// HTML-uitvoer van begin van de pagina
	html_header(w)

	err := r.ParseForm()
	if err != nil {
		writeErr(w, err)
		return
	}

	// HTML-uitvoer van het formulier
	// Returnwaarde is true als er een query was gedefinieerd
	has_query := html_form(w, r)

	// Als er geen query is gedefinieerd, HTML-uitvoer van korte helptekst, pagina-einde, en exit
	if !has_query {
		html_uitleg(w)
		html_footer(w)
		return
	}

	prefix := first(r, "db")
	if prefix == "" && len(opt_db) > 0 {
		prefix = opt_dbc[0].Id
	}
	if !prefixes[prefix] {
		writeErr(w, errors.New("Invalid corpus: "+prefix))
		return
	}

	offset := 0
	o, err := strconv.Atoi(first(r, "offset"))
	if err == nil {
		offset = o
	}

	db, err := connect()
	if err != nil {
		writeErr(w, err)
		return
	}
	defer db.Close()

	fmt.Fprintln(w, "<hr>")

	// BEGIN: Opstellen van de query
	// Deze code moet gelijk zijn aan die in het programma 'lassystats'
	parts := make([]string, 0, 6)
	for _, p := range []string{"", "h"} {
		if first(r, p+"word") != "" {
			wrd := first(r, p+"word")
			if wrd[0] == '+' {
				parts = append(parts, fmt.Sprintf("`"+p+"lemma` = %q", wrd[1:]))
			} else if wrd[0] == '@' {
				parts = append(parts, fmt.Sprintf("`"+p+"root` = %q", wrd[1:]))
			} else if wrd[0] == '=' {
				parts = append(parts, fmt.Sprintf("`"+p+"word` = %q COLLATE \"utf8_bin\"", wrd[1:]))
			} else if wrd[0] == '?' {
				parts = append(parts, fmt.Sprintf("`"+p+"word` = %q", wrd[1:]))
			} else if strings.Index(wrd, "%") >= 0 {
				parts = append(parts, fmt.Sprintf("`"+p+"word` LIKE %q", wrd))
			} else {
				var s string
				rows, err := db.Query(fmt.Sprintf("SELECT `lemma` FROM `%s_c_%s_word` WHERE `word` = %q", PRE, prefix, wrd))
				if err != nil {
					writeErr(w, err)
					return
				}
				lset := make(map[string]bool)
				for rows.Next() {
					err := rows.Scan(&s)
					if err != nil {
						writeErr(w, err)
						return
					}
					for _, i := range strings.Split(s, "\t") {
						lset[fmt.Sprintf("%q", i)] = true
					}
				}
				err = rows.Err()
				if err != nil {
					writeErr(w, err)
					return
				}
				if len(lset) > 0 {
					ll := make([]string, 0, len(lset))
					for key := range lset {
						ll = append(ll, key)
					}
					if len(ll) == 1 {
						parts = append(parts, "`"+p+"lemma` = "+ll[0])
					} else {
						parts = append(parts, "`"+p+"lemma` IN ("+strings.Join(ll, ", ")+")")
					}
				} else {
					parts = append(parts, fmt.Sprintf("`"+p+"word` = %q", wrd))
				}
			}
		}
	}
	if s := first(r, "postag"); s != "" {
		parts = append(parts, fmt.Sprintf("`postag` = %q", s))
	}
	if s := first(r, "rel"); s != "" {
		parts = append(parts, fmt.Sprintf("`rel` = %q", s))
	}
	if s := first(r, "hpostag"); s != "" {
		parts = append(parts, fmt.Sprintf("`hpostag` = %q", s))
	}
	query := strings.Join(parts, " AND ")
	// EINDE: Opstellen van de query

	// DEBUG: HTML-uitvoer van de query
	fmt.Fprint(w, "<div style=\"font-family:monospace\">\n", html.EscapeString(query), "\n</div><p>\n")

	// Maximaal ZINMAX matchende xml-bestanden opvragen
	rows, err := db.Query(
		"SELECT `arch`,`file` FROM `" + PRE + "_c_" + prefix + "_deprel` WHERE " + query + " GROUP BY `arch`,`file` LIMIT " + fmt.Sprint(offset) + ", " + fmt.Sprint(ZINMAX))
	if err != nil {
		writeErr(w, err)
		return
	}
	zinnen := make([]*Sentence, 0, ZINMAX)
	var a, f int
	for rows.Next() {
		err := rows.Scan(&a, &f)
		if err != nil {
			writeErr(w, err)
			return
		}
		zinnen = append(zinnen, &Sentence{arch: a, file: f, items: make([]Row, 0)})
	}
	err = rows.Err()
	if err != nil {
		writeErr(w, err)
		return
	}

	// Gegevens bij gevonden xml-bestanden opvragen
	for _, zin := range zinnen {
		// Zin bij xml-bestand opvragen
		var s string
		rows, err := db.Query(fmt.Sprintf("SELECT `sent` FROM `%s_c_%s_sent` WHERE `arch` = %d AND `file`= %d", PRE, prefix, zin.arch, zin.file))
		if err != nil {
			writeErr(w, err)
			return
		}
		if rows.Next() {
			err := rows.Scan(&s)
			if err != nil {
				writeErr(w, err)
				return
			}
			zin.words = strings.Fields(s)
			rows.Close()
		} else {
			fmt.Printf("Missing sentence for file id %d\n", zin.file)
			os.Exit(1)
		}

		// Matchende dependency relations bij xml-bestand opvragen
		rows, err = db.Query(fmt.Sprintf(
			"SELECT `word`,`lemma`,`postag`,`rel`,`hpostag`,`hlemma`,`hword`,`begin`,`end`,`hbegin`,`hend`,`mark` FROM `%s_c_%s_deprel` WHERE `arch` = %d AND `file`= %d AND %s ORDER BY `begin`,`hbegin`,`rel`", PRE, prefix, zin.arch, zin.file, query))
		if err != nil {
			writeErr(w, err)
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
			if err != nil {
				writeErr(w, err)
				return
			}
			zin.items = append(zin.items, r)
		}
		err = rows.Err()
		if err != nil {
			writeErr(w, err)
			return
		}
	}

	// Verwerking en HTML-uitvoer van zinnen en dependency relations
	fmt.Fprintln(w, "<ol>")
	for i, zin := range zinnen {

		// Nodes die gemarkeerd moeten worden in de boom, voor link naar 'lassytree'
		mark := make(map[string]bool)

		// Begin zin
		fmt.Fprintf(w, "<li value=\"%d\" class=\"li1\">", i+offset+1)

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
				fmt.Fprint(w, YELGRN)
			} else if idx[j]&1 != 0 {
				fmt.Fprint(w, YELLOW)
			} else if idx[j]&2 != 0 {
				fmt.Fprint(w, GREEN)
			}
			fmt.Fprint(w, word)
			if idx[j]&4 != 0 {
				fmt.Fprint(w, "</span>")
			}
			fmt.Fprint(w, " ")
		}

		// HTML-uitvoer van link naar 'lassytree'
		ms := make([]string, 0, len(mark))
		for m := range mark {
			ms = append(ms, m)
		}
		fmt.Fprintf(w, "\n(<a href=\"/tree?db=%s&amp;arch=%d&amp;file=%d&amp;yl=%s&amp;gr=%s&amp;ms=%s\">boom</a>)\n<ul>\n",
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
			q := make_query_string(item.word, item.postag, item.rel, item.hpostag, item.hword, prefix)
			fmt.Fprintf(w, "<li class=\"li2\"><a href=\"/?%s\">%s</a>\n", q, html.EscapeString(s))
		}

		// Einde zin + dependency relations
		fmt.Fprint(w, "</ul>\n\n")
	}
	fmt.Fprint(w, "</ol>\n<p>\n")

	// Links naar volgende en vorige pagina's met resultaten
	qs := make_query_string(
		first(r, "word"),
		first(r, "postag"),
		first(r, "rel"),
		first(r, "hpostag"),
		first(r, "hword"),
		first(r, "db"))
	if offset > 0 || len(zinnen) == ZINMAX {
		if offset > 0 {
			fmt.Fprintf(w, "<a href=\"/?%s&amp;offset=%d\">vorige</a>", qs, offset-ZINMAX)
		} else {
			fmt.Fprint(w, "vorige")
		}
		fmt.Fprint(w, " | ")
		if len(zinnen) == ZINMAX {
			fmt.Fprintf(w, "<a href=\"/?%s&amp;offset=%d\">volgende</a>", qs, offset+ZINMAX)
		} else {
			fmt.Fprint(w, "volgende")
		}
	}

	fmt.Fprintln(w, "<hr><small>tijd:", time.Now().Sub(now), "</small><hr>")

	// Link naar statistieken
	fmt.Fprintf(w, `<p>
		<div id="stats">
		<button onclick="javascript:$.fn.stats('/stats?word=%s&amp;postag=%s&amp;rel=%s&amp;hpostag=%s&amp;hword=%s&amp;db=%s')">statistiek</button>
		</div>
`,
		urlencode(first(r, "word")),
		urlencode(first(r, "postag")),
		urlencode(first(r, "rel")),
		urlencode(first(r, "hpostag")),
		urlencode(first(r, "hword")),
		urlencode(first(r, "db")))

	html_footer(w)
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

func html_header(w http.ResponseWriter) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head>
<meta name="robots" content="noindex,nofollow">
<title>PaQu1</title>
<script type="text/javascript" src="/jquery.js"></script>
<script type="text/javascript"><!--
  $.fn.stats = function(url) {
    $("#stats").html("Loading...");
    $.get(url, function(data) {
      $("#stats").html(data);
    }).fail(function(e) {
      $("#stats").html(e.responseText);
    });
  }
  formclear = function(f) {
    f.word.value = "";
    f.postag.value = "";
    f.rel.value = "";
    f.hpostag.value = "";
    f.hword.value = "";
  }
  //--></script>
<link rel="stylesheet" type="text/css" href="/wordrel.css">
</head>
<body>
`)
}

func html_form(w http.ResponseWriter, r *http.Request) (has_query bool) {

	has_query = true
	if first(r, "word") == "" &&
		first(r, "postag") == "" &&
		first(r, "rel") == "" &&
		first(r, "hpostag") == "" &&
		first(r, "hword") == "" {
		has_query = false
		fmt.Fprintln(w, "<img src=\"paqu.png\" class=\"logo\"><br>")
	}

	fmt.Fprint(w, `
<form action="/" method="get" accept-charset="utf-8">
<table>
<tr>
  <td colspan="3" style="padding-bottom:1em">corpus: <select name="db">
`)
	html_opts(w, opt_db, first(r, "db"), "")
	fmt.Fprintf(w, `
</select>
<tr>
  <td style="background-color: yellow">woord
  <td>
  <td style="background-color: lightgreen">hoofdwoord
<tr>
  <td><input type="text" name="word" size="20" value="%s">
`, html.EscapeString(first(r, "word")))
	fmt.Fprint(w, `
  <td>
    <select name="rel">
`)
	html_opts(w, opt_rel, first(r, "rel"), "relatie")
	fmt.Fprintf(w, `
    </select>
  <td><input type="text" name="hword" size="20" value="%s">
`, html.EscapeString(first(r, "hword")))
	fmt.Fprint(w, `
<tr>
  <td>
    <select name="postag" style="width: 100%">
`)
	html_opts(w, opt_postag, first(r, "postag"), "postag")
	fmt.Fprintf(w, `
    </select>
  <td>
  <td>
    <select name="hpostag" style="width:100%%" >
`)
	html_opts(w, opt_hpostag, first(r, "hpostag"), "postag")
	fmt.Fprint(w, `
    </select>
<tr>
  <td style="padding-top:1em">
    <input type="button" value="help" onClick="javascript:window.open('/info.html')">
  <td colspan="2" class="right" style="padding-top:1em">
    <input type="submit" value="Zoeken">
    <input type="button" value="Wissen" onClick="javascript:formclear(form)">
    <input type="reset" value="Reset">
</table>
</form>
`)

	return
}

func html_uitleg(w http.ResponseWriter) {
	fmt.Fprint(w, `
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
<a href="/info.html">Meer info...</a>
`)
}

func html_footer(w http.ResponseWriter) {
	fmt.Fprint(w, `
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

func html_opts(w http.ResponseWriter, opts []string, value, title string) {
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
				fmt.Fprintf(w, "    <optgroup label=\"&mdash; %s &mdash;\">\n", t)
			}
			opt = opt[1:]
		}
		sel := ""
		if opt == value {
			sel = " selected=\"selected\""
		}
		if opt == "" {
			fmt.Fprintf(w, "    <option value=\"\"%s>&mdash; %s &mdash;</option>\n", sel, title)
		} else if txt == "" {
			fmt.Fprintf(w, "    <option%s>%s</option>\n", sel, opt)
		} else {
			fmt.Fprintf(w, "    <option value=\"%s\"%s>%s</option>\n", opt, sel, html.EscapeString(txt))
		}
	}
}
