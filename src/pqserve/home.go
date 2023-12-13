package main

import (
	"fmt"
	"html"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// TAB: begin
func home(q *Context) {
	prefix := getprefix(q)
	if !q.prefixes[prefix] {
		http.Error(q.w, "Invalid corpus: "+prefix, http.StatusPreconditionFailed)
		return
	}

	zinmax := getzinmax(q)

	// HTML-uitvoer van begin van de pagina
	writeHead(q, "Zoeken", 1)
	html_header(q)

	// HTML-uitvoer van het formulier
	// Returnwaarde is true als er een query was gedefinieerd
	has_query := html_form(q, prefix, zinmax)

	// Als er geen query is gedefinieerd, HTML-uitvoer van korte helptekst, pagina-einde, en exit
	if !has_query {
		html_uitleg(q)
		q.w.Write([]byte(Cfg.Message))
		html_footer(q)
		return
	}

	var chClose <-chan bool
	if f, ok := q.w.(http.CloseNotifier); ok {
		chClose = f.CloseNotify()
	} else {
		chClose = make(<-chan bool)
	}

	_, err := sqlDB.Exec(fmt.Sprintf("UPDATE `%s_info` SET `active` = NOW() WHERE `id` = %q", Cfg.Prefix, prefix))
	if doErr(q, err) {
		return
	}

	offset := 0
	o, err := strconv.Atoi(first(q.r, "offset"))
	if err == nil {
		offset = o
	}
	if offset < 0 {
		offset = 0
	}

	fmt.Fprintln(q.w, "<hr>")

	query, joins, usererr, syserr := makeQuery(q, prefix, "", chClose)
	if doErr(q, syserr) {
		return
	}
	if userErr(q, usererr) {
		return
	}

	// DEBUG: HTML-uitvoer van de query
	fmt.Fprint(q.w, "<div style=\"font-family:monospace\">\n", html.EscapeString(query), "\n</div><p>\n")

	fmt.Fprint(q.w, "<div id=\"busy1\"><img src=\"busy.gif\" alt=\"aan het werk...\"></div>\n")

	if ff, ok := q.w.(http.Flusher); ok {
		ff.Flush()
	}

	now := time.Now()

	// Maximaal ZINMAX matchende xml-bestanden opvragen
	select {
	case <-chClose:
		logerr(errConnectionClosed)
		return
	default:
	}

	// Om resultaten te krijgen die gegarandeerd correct zijn zou je "ORDER BY 1,2" moeten toevoegen, maar
	// dat maakt het veel trager (vooral als er heel veel hits zijn), en zo lijkt het ook goed te werken.
	rows, err := timeoutQuery(q, chClose,
		"SELECT DISTINCT `arch`,`file` FROM `"+Cfg.Prefix+"_c_"+prefix+"_deprel` "+joins+" WHERE "+query+" LIMIT "+fmt.Sprint(offset)+", "+fmt.Sprint(zinmax))
	if doErr(q, err) {
		busyClear(q)
		return
	}
	zinnen := make([]*Sentence, 0, zinmax)
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
			logerr(errConnectionClosed)
			return
		default:
		}
		rows, err := sqlDB.Query(fmt.Sprintf("SELECT `sent` FROM `%s_c_%s_sent` WHERE `arch` = %d AND `file`= %d",
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
			doErr(q, fmt.Errorf("Zin ontbreekt voor bestand met id %v\n", zin.file))
			busyClear(q)
			return
		}

		// Matchende dependency relations bij xml-bestand opvragen
		select {
		case <-chClose:
			logerr(errConnectionClosed)
			return
		default:
		}
		rows, err = sqlDB.Query(fmt.Sprintf(
			"SELECT `word`,`lemma`,`postag`,`rel`,`hpostag`,`hlemma`,`hword`,`begin`,`end`,`hbegin`,`hend`,`mark` FROM `%s_c_%s_deprel` "+joins+" WHERE `arch` = %d AND `file`= %d AND ( %s ) ORDER BY `begin`,`hbegin`,`rel`",
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
			fmt.Fprint(q.w, html.EscapeString(unHigh(word)))
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
		fmt.Fprintf(q.w, "\n<a href=\"tree?db=%s&amp;arch=%d&amp;file=%d&amp;yl=%s&amp;gr=%s&amp;ms=%s\" class=\"ico\" target=\"_blank\">&#10020;</a>\n<ul>\n",
			prefix,
			zin.arch,
			zin.file,
			strings.Join(yellows, ","),
			strings.Join(greens, ","),
			strings.Join(ms, ","))

		// HTML-uitvoer van alle matchende dependency relations voor huidige zin
		seen := make(map[string]bool)
		lines := make([][2]string, 0)
		links := make(map[string][]string)
		for _, item := range zin.items {
			s := fmt.Sprintf("%s:%s — %s — %s:%s",
				html.EscapeString(unHigh(item.word)), item.postag,
				item.rel,
				html.EscapeString(unHigh(item.hword)), item.hpostag)
			if !seen[s] {
				seen[s] = true
				qq := make_query_string(item.word, item.postag, item.rel, item.hpostag, item.hword, first(q.r, "meta"), prefix)
				line := fmt.Sprintf("<li class=\"li2\"><a href=\"?%s\">%s</a>\n", qq, s)
				lines = append(lines, [2]string{s, line})
				links[s] = make([]string, 0)
			}
			// Link naar bomen voor individuele matchende dependency relation
			green := ""
			if item.hword != "" {
				green = fmt.Sprint(item.hbegin)
			}
			links[s] = append(links[s],
				fmt.Sprintf("<a href=\"tree?db=%s&amp;arch=%d&amp;file=%d&amp;yl=%d&amp;gr=%s&amp;ms=%s\" class=\"icol\" target=\"_blank\">&#10020;</a>\n",
					prefix,
					zin.arch,
					zin.file,
					item.begin,
					green,
					item.mark))

		}
		for _, line := range lines {
			fmt.Fprint(q.w, line[1])
			for _, link := range links[line[0]] {
				fmt.Fprint(q.w, link)
			}
		}

		// Einde zin + dependency relations
		fmt.Fprint(q.w, "</ul>\n\n")
	}
	fmt.Fprint(q.w, "</ol>\n<p>\n")

	if offset == 0 && len(zinnen) == 0 {
		fmt.Fprintln(q.w, "geen match gevonden")
	}

	defer html_footer(q)

	// Links naar volgende en vorige pagina's met resultaten
	qs := make_query_string(
		first(q.r, "word"),
		first(q.r, "postag"),
		first(q.r, "rel"),
		first(q.r, "hpostag"),
		first(q.r, "hword"),
		first(q.r, "meta"),
		first(q.r, "db"))
	if offset > 0 || len(zinnen) == zinmax {
		if offset > 0 {
			fmt.Fprintf(q.w, "<a href=\"?%s&amp;offset=%d\">vorige</a>", qs, offset-zinmax)
		} else {
			fmt.Fprint(q.w, "vorige")
		}
		fmt.Fprint(q.w, " | ")
		if len(zinnen) == zinmax {
			fmt.Fprintf(q.w, "<a href=\"?%s&amp;offset=%d\">volgende</a>", qs, offset+zinmax)
		} else {
			fmt.Fprint(q.w, "volgende")
		}
	}

	if offset > 0 || len(zinnen) > 0 {
		fmt.Fprintf(q.w, `<p>
<form action="dl" method="POST" accept-charset="UTF-8" enctype="multipart/form-data">
<input type="hidden" name="word" value="%s">
<input type="hidden" name="postag" value="%s">
<input type="hidden" name="rel" value="%s">
<input type="hidden" name="hpostag" value="%s">
<input type="hidden" name="hword" value="%s">
<input type="hidden" name="meta" value="%s">
<input type="hidden" name="db" value="%s">
<select name="step">
<option value="1">alles</option>
<option value="2">1 per 2</option>
<option value="5">1 per 5</option>
<option value="10">1 per 10</option>
<option value="20">1 per 20</option>
<option value="50">1 per 50</option>
<option value="100">1 per 100</option>
<option value="200">1 per 200</option>
<option value="500">1 per 500</option>
<option value="1000">1 per `+iformat(1000)+`</option>
<option value="2000">1 per `+iformat(2000)+`</option>
<option value="5000">1 per `+iformat(5000)+`</option>
<option value="10000">1 per `+iformat(10000)+`</option>
</select>
<select name="what">
<option value="text">zinnen</option>
<option value="textcount">zinnen met telling</option>
<option value="data">alle data</option>
</select>
<input type="submit" value="downloaden">
</form>
`,
			html.EscapeString(first(q.r, "word")),
			html.EscapeString(first(q.r, "postag")),
			html.EscapeString(first(q.r, "rel")),
			html.EscapeString(first(q.r, "hpostag")),
			html.EscapeString(first(q.r, "hword")),
			html.EscapeString(first(q.r, "meta")),
			html.EscapeString(prefix))
	}

	if q.auth && (offset > 0 || len(zinnen) > 0) {
		fmt.Fprintf(q.w, `<p>
<form action="savez" method="POST" accept-charset="UTF-8" enctype="multipart/form-data">
<input type="hidden" name="word" value="%s">
<input type="hidden" name="postag" value="%s">
<input type="hidden" name="rel" value="%s">
<input type="hidden" name="hpostag" value="%s">
<input type="hidden" name="hword" value="%s">
<input type="hidden" name="meta" value="%s">
<input type="hidden" name="db" value="%s">
<input type="submit" value="nieuw corpus maken op basis van deze zoekopdracht">
</form>
`,
			html.EscapeString(first(q.r, "word")),
			html.EscapeString(first(q.r, "postag")),
			html.EscapeString(first(q.r, "rel")),
			html.EscapeString(first(q.r, "hpostag")),
			html.EscapeString(first(q.r, "hword")),
			html.EscapeString(first(q.r, "meta")),
			html.EscapeString(prefix))
	}

	fmt.Fprintln(q.w, "<hr><small>tijd:", tijd(time.Now().Sub(now)), "</small>")

	if offset == 0 && len(zinnen) == 0 {
		return
	}

	// Links naar statistieken

	if q.hasmeta[prefix] {
		metahelp(q)
	}

	fmt.Fprintf(q.w, `<p>
        <div id="statsrel">
        <form action="javascript:$.fn.statsrel()" name="statsrelform" id="statsrelform">
        <input type="hidden" name="word" value="%s">
        <input type="hidden" name="postag" value="%s">
        <input type="hidden" name="rel" value="%s">
        <input type="hidden" name="hpostag" value="%s">
        <input type="hidden" name="hword" value="%s">
        <input type="hidden" name="meta" value="%s">
        <input type="hidden" name="db" value="%s">
        Selecteer &eacute;&eacute;n tot vijf elementen:
        <p>
        <table>
        <tr style="vertical-align:top"><td>
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
        <td>
`,
		html.EscapeString(first(q.r, "word")),
		html.EscapeString(first(q.r, "postag")),
		html.EscapeString(first(q.r, "rel")),
		html.EscapeString(first(q.r, "hpostag")),
		html.EscapeString(first(q.r, "hword")),
		html.EscapeString(first(q.r, "meta")),
		html.EscapeString(prefix))

	if q.hasmeta[prefix] {
		for _, meta := range getMeta(q, prefix) {
			fmt.Fprintf(q.w, "<input type=\"checkbox\" name=\"cmeta\" value=\"%s\">%s<br>\n", html.EscapeString(meta.name), html.EscapeString(meta.name))
		}
	}

	fmt.Fprint(q.w, `
        </table>
        <p>
        <input type="submit" id="statsrelsubmit" value="doe telling">
        </form>
        <p>
        <div id="statstel" class="hide">
            <table>
            <tr><td>items:<td class="right" id="statstel1"><img src="busy.gif" alt="aan het werk...">
            <tr><td>zinnen:<td class="right" id="statstel2"><img src="busy.gif" alt="aan het werk...">
            </table>
            </div>
        <p>
        <div id="statresults">
        </div>
        </div>
<script type="text/javascript"><!--
  hexEncode = function(s) {
    var hex, i;
    var result = "";
    for (i = 0; i < s.length; i++) {
      hex = s.charCodeAt(i).toString(16);
      result += ("000"+hex).slice(-4);
    }
    return result;
  }
  hexDecode = function(s) {
    var j;
    var hexes = s.match(/.{1,4}/g) || [];
    var back = "";
    for(j = 0; j < hexes.length; j++) {
      back += String.fromCharCode(parseInt(hexes[j], 16));
    }
    return back;
  }
  function setCookie(cname, cvalue, exdays) {
    var d = new Date();
    d.setTime(d.getTime() + (exdays*24*60*60*1000));
    var expires = "expires="+d.toUTCString();
    document.cookie = cname + "=" + cvalue + "; " + expires;
  }
  function getCookie(cname) {
    var name = cname + "=";
    var ca = document.cookie.split(';');
    for(var i=0; i<ca.length; i++) {
        var c = ca[i];
        while (c.charAt(0)==' ') c = c.substring(1);
        if (c.indexOf(name) != -1) return c.substring(name.length,c.length);
    }
    return "";
  }
  function statsrelformcheck() {
    var f = document.forms["statsrelform"];
    var n = 0;
    if (f.cword.checked   ) { n++; }
    if (f.clemma.checked  ) { n++; }
    if (f.cpostag.checked ) { n++; }
    if (f.crel.checked    ) { n++; }
    if (f.chword.checked  ) { n++; }
    if (f.chlemma.checked ) { n++; }
    if (f.chpostag.checked) { n++; }
    if (typeof f.cmeta !== 'undefined' && f.cmeta !== null) {
      if (f.cmeta.checked) {
        n++;
      } else {
        for (i = 0; i < f.cmeta.length; i++) {
           if (f.cmeta[i].checked) { n++; }
        }
      }
    }
    if (n < 1 || n > 5) {
      $('#statsrelsubmit').prop('disabled', true);
    } else {
      $('#statsrelsubmit').prop('disabled', false);
      var el = $('#statsrelform input:checked');
      var items = []
      for (i = 0; i < el.length; i++) {
          items.push(hexEncode(el[i].name));
      }
      setCookie("paqu-basic", items.join('|'), 14);
    }
  }
  $('#statsrelform input').on('change', function (e) {
    statsrelformcheck();
  });
  var ck = getCookie('paqu-basic').split('|');
  var f = document.forms["statsrelform"];
  for (i = 0; i < ck.length; i++) {
    try {
        f[hexDecode(ck[i])].checked = true;
    }
    catch(err) {}
  }
  statsrelformcheck();
//--></script>
`)
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

func make_query_string(word, postag, rel, hpostag, hword, meta, db string) string {
	return fmt.Sprintf(
		"word=%s&amp;postag=%s&amp;rel=%s&amp;hpostag=%s&amp;hword=%s&amp;meta=%s&amp;db=%s",
		urlencode(unHigh(word)),
		urlencode(postag),
		urlencode(rel),
		urlencode(hpostag),
		urlencode(unHigh(hword)),
		urlencode(meta),
		urlencode(db))
}

func make_query_string_js(word, postag, rel, hpostag, hword, meta, db string) string {
	return fmt.Sprintf(
		"word=%s&postag=%s&rel=%s&hpostag=%s&hword=%s&meta=%s&db=%s",
		urlencode(unHigh(word)),
		urlencode(postag),
		urlencode(rel),
		urlencode(hpostag),
		urlencode(unHigh(hword)),
		urlencode(meta),
		urlencode(db))
}

//. HTML

func html_header(q *Context) {
	qs := make_query_string_js(
		first(q.r, "word"),
		first(q.r, "postag"),
		first(q.r, "rel"),
		first(q.r, "hpostag"),
		first(q.r, "hword"),
		first(q.r, "meta"),
		first(q.r, "db"))
	fmt.Fprint(q.w, `
<script type="text/javascript" src="jquery.js"></script>
<script type="text/javascript"><!--

  function ival(i) {
      var s1 = "".concat(i);
      var s2 = "";
      for (var n = s1.length; n > 3; n = s1.length) {
         s2 = "&#8239;".concat(s1.substr(n-3, n), s2);
         s1 = s1.substr(0, n-3);
      }
      return s1.concat(s2);
  }

  function sortMeta(tbl, colno) {
      var n = tbl[0].length - 2;
      if (colno == n) {
        tbl.sort(function(a, b) {return (a[n] - b[n]);});
      } else {
        tbl.sort(function(a, b) {
          var c = b[colno] - a[colno];
          if (c == 0) {
             return (a[n] - b[n]);
          }
          return c;
        });
      }
      return tbl;
  }

  function fillMeta() {
    var ida = "#metaa";
    var idb = "#metab";
    var fl = metavars.fl;
    var lbl = metavars.lbl;
    var max = metavars.max;
    var da = metavars.a;
    var db = metavars.b;
    var ac = metavars.ac;
    var bc = metavars.bc;
    var ac0 = "";
    var ac1 = "";
    var bc0 = "";
    var bc1 = "";
    var bc2 = "";
    if (ac == 0) { ac0 = " em"; }
    if (ac == 1) { ac1 = " em"; }
    if (bc == 0) { bc0 = " em"; }
    if (bc == 1) { bc1 = " em"; }
    if (bc == 2) { bc2 = " em"; }
    var a = $(ida);
    var b = $(idb);
    a.html('<tr><td class="link a'+ac0+'">aantal<td class="link b '+fl+ac1+'">'+lbl+'\n');
    b.html('<tr><td class="link a'+bc0+'">aantal<td class="link b'+bc1+'">per&nbsp;'+ival(metadn)+'<td class="link c '+fl+bc2+'">'+lbl+'\n');
    for (i in da) {
       if (i > max) {
         a.append('<tr><td><td class="' + fl + '">...\n');
         b.append('<tr><td><td>...<td class="' + fl + '">...\n');
         break;
       }
       var cl = '';
       var v = da[i][2];
       if (da[i][1] == 2147483647) {
         cl = ' nil';
         v = '(leeg)';
       }
       a.append('<tr><td>' + ival(da[i][0]) + '<td class="' + fl + cl + '">' + v + '\n');
       cl = '';
       v = db[i][3];
       if (db[i][2] == 2147483647) {
         cl = ' nil';
         v = '(leeg)';
       }
       b.append('<tr><td>' + ival(db[i][0]) + '<td>' + ival(db[i][1]) + '<td class="' + fl + cl + '">' + v + '\n');
    }
    $(ida + ' td.a').on('click', function() {
         metavars.a = sortMeta(da, 0);
         metavars.ac = 0;
         fillMeta();
      });
    $(ida + ' td.b').on('click', function() {
         metavars.a = sortMeta(da, 1);
         metavars.ac = 1;
         fillMeta();
      });
    $(idb + ' td.a').on('click', function() {
         metavars.b = sortMeta(db, 0);
         metavars.bc = 0;
         fillMeta();
      });
    $(idb + ' td.b').on('click', function() {
         metavars.b = sortMeta(db, 1);
         metavars.bc = 1;
         fillMeta();
      });
    $(idb + ' td.c').on('click', function() {
         metavars.b = sortMeta(db, 2);
         metavars.bc = 2;
         fillMeta();
      });
  }

  function formclear(f) {
    f.word.value = "";
    f.postag.value = "";
    f.rel.value = "";
    f.hpostag.value = "";
    f.hword.value = "";
    f.meta.value = "";
  }

  var metadn = 0;
  var metavars = {};
  var metavisible = false;
  var queryvisible = false;

  function metahelp() {
    var e = $("#helpmeta");
    e.show();
    e.css("zIndex", 9999);
    metavisible = true;
    return false;
  }
  function queryhelp() {
    var e = $("#helpquery");
    e.show();
    e.css("zIndex", 9999);
    queryvisible = true;
    return false;
  }

  $(document).mouseup(
    function(e) {
      if (metavisible) {
        var e = $("#helpmeta");
        e.hide();
        e.css("zIndex", 1);
        metavisible = false;
      }
      if (queryvisible) {
        var e = $("#helpquery");
        e.hide();
        e.css("zIndex", 1);
        queryvisible = false;
      }
    });

  var entityMap = {
    "&": "&amp;",
    "<": "&lt;",
    ">": "&gt;",
    '"': '&quot;',
    "'": '&#39;',
    "/": '&#x2F;'
  };

  function escapeHtml(string) {
    return String(string).replace(/[&<>"'\/]/g, function (s) {
      return entityMap[s];
    });
  }

  var lastcall = null;
  var statsreldata;
  var statsrelcol = 0;
  var firsttime = true;
  $.fn.statsrel = function() {
    if (lastcall) {
      try {
        lastcall.abort();
      }
      catch(err) {}
    }
    if (firsttime) {
        firsttime = false;
        $('#statstel').removeClass('hide');
        $.ajax("statstel?t=1&`+qs+`")
          .done(function(data) {
            $("#statstel1").html(data);
          }).fail(function(e) {
            $("#statstel1").html(escapeHtml(e.responseText));
          });
        $.ajax("statstel?t=2&`+qs+`")
          .done(function(data) {
            $("#statstel2").html(data);
          }).fail(function(e) {
            $("#statstel2").html(escapeHtml(e.responseText));
          });
    }
    $("#statresults").html('<img src="busy.gif">');
    if ($('#statsrelform input:checked').length == 1) {
      if ($('#statsrelform input[name=cmeta]:checked').length == 1) {
        var val = $('#statsrelform input[name=cmeta]:checked')[0].value;
        lastcall = $.ajax("statsmeta?item=" + val + "&" + $(document.statsrelform).serialize())
        .done(function(data) {
          var e = $("#statresults");
          if (data.err != "") {
            $("#statresults").html('<div class="error">Fout: ' + escapeHtml(data.err) + '</div>');
            return;
          }
          e.html('<!-- <div style="font-family:monospace">' + escapeHtml(data.query) + '</div> -->\n' +
            '<p>\n' +
            '<b>' + escapeHtml(val) + '</b>' +
            ' &mdash; <a href="javascript:void(0)" onclick="javascript:metahelp()">toelichting</a>\n' +
            '<p>\n' +
            '<table>\n' +
            '  <tr>\n' +
            '   <td>per item:\n' +
            '     <table class="right" id="metaa">\n' +
            '     </table>\n' +
            '   <td class="next">per zin:\n' +
            '     <table class="right" id="metab">\n' +
            '     </table>\n' +
            '</table>\n' +
            '<hr>tijd: ' + data.tijd + '<p><a href="statsmeta?' + data.download + '">download</a>\n');
          metadn = data.n;
          metavars.lbl = data.value;
          metavars.fl = data.fl;
          metavars.max = data.max;
          metavars.ac = data.ac;
          metavars.bc = data.bc;
          metavars.a = data.lines[0];
          metavars.b = data.lines[1];
          fillMeta();
        }).fail(function(e) {
          $("#statresults").html('<div class="error">Fout: ' + escapeHtml(e.responseText) + '</div>');
        })
        .always(function() {
          lastcall = null;
        });
      } else {
        var val = $('#statsrelform input:checked')[0].name;
        lastcall = $.ajax("stats?item=" + val + "&" + $(document.statsrelform).serialize())
        .done(function(data) {
          $("#statresults").html(data);
        }).fail(function(e) {
          $("#statresults").html('<div class="error">Fout: ' + escapeHtml(e.responseText) + '</div>');
        })
        .always(function() {
          lastcall = null;
        });
      }
    } else {
      lastcall = $.ajax("statsrel?" + $(document.statsrelform).serialize())
      .done(function(data) {
        statsreldata = data;
        var e = $("#statresults");
        e.html('<!-- <div style="font-family:monospace">' + data.query +
          '</div> --> <p><table class="breed"></table>');
        if (data.toomany) {
            e.append('<div class="warning">Te veel treffers. Bij het sorteren kunnen treffers met lagere aantallen ontbreken.</div>');
        }
        e.append('<hr>tijd: ' + data.tijd +
          '<p><a href="statsrel?' + data.download + '">download</a>');
        statsrelcol = 0;
        fillstatsrel();
      }).fail(function(e) {
        $("#statresults").html('<div class="error">Fout: ' + escapeHtml(e.responseText) + '</div>');
      })
      .always(function() {
        lastcall = null;
      });
    }
  }

  function statrelset(c) {
    statsreldata.lines.sort(function(a, b) {
        if (c == 0) {
            var r = b[0][1] - a[0][1];
            if (r != 0) { return r; }
        } else if (statsreldata.isint[c]) {
            var r = a[c][1] - b[c][1];
            if (r != 0) { return r; }
        } else {
            if (a[c][0] == "" && b[c][0] != "") {
                return 1;
            }
            if (a[c][0] != "" && b[c][0] == "") {
                return -1;
            }
            if (a[c][0] < b[c][0]) {
                return -1;
            }
            if (a[c][0] > b[c][0]) {
                return 1;
            }
        }
        for (i in statsreldata.isint) {
            if (i == c) { continue; }
            if (i == 0) {
                var r = b[0][1] - a[0][1];
                if (r != 0) { return r; }
            } else if (statsreldata.isint[i]) {
                var r = a[i][1] - b[i][1];
                if (r != 0) { return r; }
            } else {
                if (a[i][0] == "" && b[i][0] != "") {
                    return 1;
                }
                if (a[i][0] != "" && b[i][0] == "") {
                    return -1;
                }
                if (a[i][0] < b[i][0]) {
                    return -1;
                }
                if (a[i][0] > b[i][0]) {
                    return 1;
                }
            }
        }
        return 0;
    });
    statsrelcol = c;
    fillstatsrel();
  }

  function statsrellink(i) {
    $('#sr' + i).on('click', function() { statrelset(i); });
  }

  function fillstatsrel() {
    var t = $("#statresults table");
    var s = '<tr class="odd">';
    for (i in statsreldata.aligns) {
        var em = "";
        if (i == statsrelcol) {
            em = " em";
        }
        s += '<th class="' + statsreldata.aligns[i] + em + ' link" id="sr' + i + '">' + statsreldata.labels[i];
    }
    t.html(s);
    for (i in statsreldata.aligns) {
        statsrellink(i)
    }
    for (j in statsreldata.lines) {
        if (j % 2 == 1) {
            s = '<tr class="odd">';
        } else {
            s = '<tr>';
        }
        if (j > 250) {
            s += '<td>';
            for (i = 1; i < statsreldata.aligns.length; i++) {
                s += '<td class="' + statsreldata.aligns[i] + '">...';
            }
            t.append(s);
            break;
        }
        for (i in statsreldata.aligns) {
            if (i == 0) {
                s += '<td class="right">';
                if (statsreldata.lines[j][i][0] == "") {
                    s += statsreldata.lines[j][i][1];
                } else {
                    s += '<a href=".?' + statsreldata.lines[j][i][0] + '">' + statsreldata.lines[j][i][1] + '</a>';
                }
            } else {
                s += '<td class="' + statsreldata.aligns[i];
                if (statsreldata.lines[j][i][0] == "") {
                    s += ' nil">(leeg)';
                } else {
                    s += '">' + escapeHtml(statsreldata.lines[j][i][0]);
                }
            }
        }
        t.append(s);
    }
  }

  function mresize() {
    localStorage.setItem(
      "paqu-home",
      JSON.stringify({
        w: metaarea.outerWidth(),
        h: metaarea.outerHeight()
      }));
  }

  function limit(v, min, max) {
    v = Number(v);
    if (v < min) {
      return min;
    }
    if (v > max) {
      return max;
    }
    return v;
  }

  function msetsize() {
    var storageContent = localStorage.getItem("paqu-home");
    if (storageContent !== undefined) {
        var d = JSON.parse(storageContent) || {};
        if (d['w']) {
          metaarea.outerWidth(limit(d['w'], 400, 1000));
        }
        if (d['h']) {
          metaarea.outerHeight(limit(d['h'], 80, 600));
        }
    }
  }

  //--></script>
`)
}

func html_uitleg(q *Context) {
	fmt.Fprint(q.w, `
<p>
<hr>
<p>
Met <abbr title="Parse &amp; Query">PaQu</abbr> kun je zoeken in
syntactisch geannoteerde Nederlandstalige corpora.
PaQu ondersteunt twee manieren van zoeken. Met de eerste,
eenvoudige, manier kun je naar woordparen zoeken, met daarbij
eventueel hun syntactische relatie. De tweede, ingewikkeldere, manier
gebruikt de zoektaal XPath.
<p>
Wanneer je inlogt kun je ook je eigen teksten laten annoteren.
<p>
<a href="info.html">Meer info...</a>
`)
}

func html_footer(q *Context) {
	fmt.Fprint(q.w, `
</body>
</html>
`)
}

func html_form(q *Context, prefix string, maxzin int) (has_query bool) {
	has_query = true
	if first(q.r, "word") == "" &&
		first(q.r, "postag") == "" &&
		first(q.r, "rel") == "" &&
		first(q.r, "hpostag") == "" &&
		first(q.r, "hword") == "" &&
		first(q.r, "meta") == "" {
		has_query = false
	}

	fmt.Fprint(q.w, `
<form action="." method="get" accept-charset="utf-8">
<a href="corpusinfo">[?]</a> corpus: <select name="db">
`)
	html_opts(q, q.opt_db, prefix, "corpus")
	fmt.Fprintln(q.w, "</select>")
	if q.auth {
		fmt.Fprintln(q.w, "<a href=\"corpuslijst\">meer/minder</a>")
	}
	fmt.Fprintf(q.w, `<p><table class="home">
       <tr>
         <td style="background-color: yellow">woord
         <td>
         <td style="background-color: lightgreen">hoofdwoord
       <tr>
         <td><input type="text" name="word" size="12" value="%s">
       `, html.EscapeString(first(q.r, "word")))
	fmt.Fprint(q.w, `
         <td>
           <select name="rel">
       `)
	html_opts(q, opt_rel, first(q.r, "rel"), "relatie")
	fmt.Fprintf(q.w, `
           </select>
         <td><input type="text" name="hword" size="12" value="%s">
       `, html.EscapeString(first(q.r, "hword")))
	fmt.Fprint(q.w, `
       <tr>
         <td>
           <select name="postag" style="width: 100%">
       `)
	html_opts(q, opt_postag, first(q.r, "postag"), "postag")
	fmt.Fprint(q.w, `
           </select>
         <td>
         <td>
           <select name="hpostag" style="width:100%" >
       `)
	html_opts(q, opt_hpostag, first(q.r, "hpostag"), "postag")
	fmt.Fprint(q.w, `
           </select>
       <tr>
         <td colspan="3"><span class="ie">Metadata:<br></span>
           <textarea rows="3" cols="40" name="meta" placeholder="metadata" id="metaarea">`+html.EscapeString(first(q.r, "meta"))+`</textarea><br>
           <small><a href="javascript:void(0)" onclick="javascript:queryhelp()">voorbeelden</a></small>
       <tr>
         <td colspan="3">aantal: <select name="sn">
`)
	for _, i := range []int{10, 20, 50, 100, 200, 500} {
		if i == maxzin {
			fmt.Fprintf(q.w, "<option selected>%d</option>\n", i)
		} else {
			fmt.Fprintf(q.w, "<option>%d</option>\n", i)
		}
	}
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
<script type="text/javascript"><!--
  var metaarea = $('#metaarea');
  metaarea.on('mouseup', mresize);
  msetsize();
//--></script>
<div class="submenu a9999" id="helpquery">
<div class="queryhelp">
<b>Voorbeelden van zoeken met metadata</b>
<dl>
<dt>Zinnen moeten aan beide voorwaardes voldoen
<dd>country = BE &amp; sex = female
<dt>Zinnen moeten aan een van de voorwaardes voldoen
<dd>country = BE | country = NL
<dt>Gebruik haakjes om prioriteit aan te geven
<dd>( country = BE | country = NL ) &amp; sex = female
<dt>Zet woorden met speciale tekens tussen dubbele aanhalingstekens
<dd>&quot;Type(A)&quot; = &quot;Eerste category&quot;
<dt>Vergelijkingen
<dd>temperature &gt; 20
<dd>temperature &gt;= 20
<dd>temperature &lt; 20
<dd>temperature &lt;= 20
<dt>Dit doet niet wat je verwacht als een zin meerdere waardes kan hebben voor een meta-attribuut
<dd>hits >= 10 & hits <= 20
<dt>Zo zoek je naar een waarde tussen 10 en 20, inclusief
<dd>hits in 10 20
<dt>Zoeken zonder onderscheid van hoofdletters en kleine letters, en met SQL-jokers % en _
<dd>name % jan
<dd>name % j_n
<dd>name % &quot;jan%&quot;
<dt>Zoeken met een reguliere expressie
<dd>name ~ [Jj]anss?ens?
<dt>Zoeken naar een ontbrekende waarde
<dd>name = nil
<dt>Zoeken naar de naam <em>nil</em>
<dd>name = &quot;nil&quot;
</dl>
</div>
</div>

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
		disabled := ""
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
			if opt[0] == '-' {
				disabled = " disabled=\"disabled\""
				opt = opt[1:]
			}
			if opt[0] != c {
				c = opt[0]
				t := ""
				switch c {
				case 'A':
					t = "algemene corpora &mdash; handmatig geannoteerd"
				case 'B':
					t = "algemene corpora &mdash; automatisch geannoteerd"
				case 'C':
					t = "algemene corpora"
				case 'D':
					t = "mijn corpora"
				case 'E':
					t = "corpora gedeeld door anderen"
				}
				if c != 'Z' {
					fmt.Fprintf(q.w, "    <optgroup label=\"&mdash; %s &mdash;\">\n", t)
				}
			}
			opt = opt[1:]
		}
		sel := ""
		if opt == value && disabled == "" {
			sel = " selected=\"selected\""
		}
		if opt == "" {
			fmt.Fprintf(q.w, "    <option value=\"\"%s%s>&mdash; %s &mdash;</option>\n", sel, disabled, title)
		} else if txt == "" {
			fmt.Fprintf(q.w, "    <option%s%s>%s</option>\n", sel, disabled, opt)
		} else {
			fmt.Fprintf(q.w, "    <option value=\"%s\"%s%s>%s</option>\n", opt, sel, disabled, html.EscapeString(txt))
		}
	}
}

func getprefix(q *Context) string {
	db := "_invalid"
	for _, p := range strings.Fields(Cfg.Default) {
		if q.prefixes[p] {
			db = p
			break
		}
	}
	if p := first(q.r, "db"); p != "" {
		if p != "_default" {
			db = p
		}
	} else if prev, err := q.r.Cookie("paqu-prev"); err == nil {
		if p := prev.Value; q.prefixes[p] {
			db = p
		}
	}
	exp := time.Now().AddDate(0, 0, 14)
	http.SetCookie(q.w, &http.Cookie{Name: "paqu-prev", Value: db, Path: cookiepath, Expires: exp})
	return db
}

func getzinmax(q *Context) int {
	sn := first(q.r, "sn")
	if sn == "" {
		if s, err := q.r.Cookie("paqu-sn"); err == nil {
			sn = s.Value
		}
	}

	si, err := strconv.Atoi(sn)
	if err != nil {
		si = ZINMAX
	} else if si < 10 {
		si = 10
	} else if si > 500 {
		si = 500
	}
	sn = fmt.Sprint(si)

	exp := time.Now().AddDate(0, 0, 14)
	http.SetCookie(q.w, &http.Cookie{Name: "paqu-sn", Value: sn, Path: cookiepath, Expires: exp})
	return si
}

func busyClear(q *Context) {
	fmt.Fprint(q.w, `<script type="text/javascript"><!--
$('#busy1').addClass('hide');
//--></script>
`)
}

var (
	reUnHigh  = regexp.MustCompile("&#[0-9]+;")
	reSetHigh = regexp.MustCompile("[^\001-\uFFFF]")
)

func unHighFunc(s string) string {
	c, _ := strconv.Atoi(s[2 : len(s)-1])
	return fmt.Sprintf("%c", c)
}

func setHighFunc(s string) string {
	u, _ := utf8.DecodeRuneInString(s)
	return fmt.Sprintf("&#%d;", u)
}

func unHigh(s string) string {
	return reUnHigh.ReplaceAllStringFunc(s, unHighFunc)
}

func setHigh(s string) string {
	return reSetHigh.ReplaceAllStringFunc(s, setHighFunc)
}

func homedl(q *Context) {
	prefix := firstf(q.form, "db")
	if prefix == "" {
		http.Error(q.w, "Geen corpus opgegeven", http.StatusPreconditionFailed)
		return
	}
	if !q.prefixes[prefix] {
		http.Error(q.w, "Ongeldig corpus: "+prefix, http.StatusPreconditionFailed)
		return
	}

	var chClose <-chan bool
	if f, ok := q.w.(http.CloseNotifier); ok {
		chClose = f.CloseNotify()
	} else {
		chClose = make(<-chan bool)
	}

	query, joins, usererr, syserr := makeQueryF(q, prefix, "", chClose)
	if doErr(q, syserr) {
		return
	}
	if userErr(q, usererr) {
		return
	}

	q.w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	q.w.Header().Set("Content-Disposition", "attachment; filename=uitvoer.txt")

	var alldata, textcount bool
	var t1, t2 string
	var cc int
	if what := firstf(q.form, "what"); what == "data" {
		alldata = true
	} else if what == "textcount" {
		textcount = true
		t1 = ", count(*)"
		t2 = " GROUP BY `arch`,`file`"
		cc = 1
	}

	columns := []string{
		"word", "begin", "end", "postag",
		"rel",
		"hword", "hbegin", "hend", "hpostag",
		"lbl", "sent",
	}

	rows, err := timeoutQuery(q, chClose,
		"SELECT `"+strings.Join(columns, "`,`")+"`"+t1+" FROM `"+Cfg.Prefix+"_c_"+prefix+"_deprel` "+
			joins+" JOIN `"+Cfg.Prefix+"_c_"+prefix+"_sent` USING (`arch`,`file`) WHERE "+query+t2)
	if doErr(q, err) {
		return
	}

	if alldata {
		fmt.Fprintln(q.w, strings.Join(columns, "\t"))
	}

	step, _ := strconv.Atoi(firstf(q.form, "step"))
	if step < 1 {
		step = 1
	}
	lineno := 0

	items := make([]string, len(columns)+cc)
	fields := make([]interface{}, len(columns)+cc)
	for i := range items {
		fields[i] = &items[i]
	}
	prev := ""
	for rows.Next() {
		err := rows.Scan(fields...)
		if doErr(q, err) {
			return
		}
		if items[9] == prev {
			continue
		}
		lineno++
		if lineno%step == 0 {
			if alldata {
				fmt.Fprintln(q.w, strings.Join(items, "\t"))
			} else {
				if textcount {
					fmt.Fprintln(q.w, items[9]+"\t"+items[11]+"\t"+items[10])
				} else {
					fmt.Fprintln(q.w, items[9]+"|"+items[10])
				}
			}
		}
		if !alldata {
			prev = items[9]
		}
		select {
		case <-chClose:
			logerr(errConnectionClosed)
			rows.Close()
			return
		default:
		}
	}
}
