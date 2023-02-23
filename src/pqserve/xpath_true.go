// +build !nodbxml

package main

import (
	"github.com/rug-compling/paqu/internal/dir"
	pqnode "github.com/rug-compling/paqu/internal/node"

	"github.com/pebbe/dbxml"

	"bytes"
	"crypto/md5"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	pqbugtest string
	pqxok     string
	reXpath   = regexp.MustCompile(
		`'[^']*'|"[^"]*"|@[_:a-zA-ZÀ-ÖØ-öø-ÿ][-._:a-zA-ZÀ-ÖØ-öø-ÿ0-9]*|\$[a-z][-_a-zA-Z0-9]*|[a-zA-Z][-_a-zA-Z]*:*(\s*\()?`)
	keyTags = make(map[string]bool)
	udTags  = map[string]bool{
		"Abbr":            true,
		"Case":            true,
		"Definite":        true,
		"Degree":          true,
		"Foreign":         true,
		"Gender":          true,
		"Number":          true,
		"Person":          true,
		"PronType":        true,
		"Reflex":          true,
		"Tense":           true,
		"VerbForm":        true,
		"auto":            true,
		"deprel":          true,
		"deprel_aux":      true,
		"deprel_main":     true,
		"elided":          true,
		"error":           true,
		"form":            true,
		"head":            true,
		"id":              true,
		"lemma":           true,
		"recursion_limit": true,
		"status":          true,
		"ud":              true,
		"upos":            true,
	}
	xpathNames = map[string]bool{
		"alpino_ds":            true,
		"ancestor-or-self::":   true,
		"ancestor::":           true,
		"attribute::":          true,
		"child::":              true,
		"descendant-or-self::": true,
		"descendant::":         true,
		"following-sibling::":  true,
		"following::":          true,
		"parent::":             true,
		"preceding-sibling::":  true,
		"preceding::":          true,
		"self::":               true,
		"fn:":                  true,
		"and":                  true,
		"as":                   true,
		"assert":               true,
		"at":                   true,
		"attribute":            true,
		"cast":                 true,
		"div":                  true,
		"element":              true,
		"else":                 true,
		"eq":                   true,
		"every":                true,
		"except":               true,
		"follows":              true,
		"for":                  true,
		"ge":                   true,
		"gt":                   true,
		"in":                   true,
		"instance":             true,
		"intersect":            true,
		"item":                 true,
		"le":                   true,
		"lt":                   true,
		"meta":                 true,
		"metadata":             true,
		"mod":                  true,
		"ne":                   true,
		"node":                 true,
		"of":                   true,
		"only":                 true,
		"or":                   true,
		"parser":               true,
		"precedes":             true,
		"return":               true,
		"satisfies":            true,
		"sentence":             true,
		"some":                 true,
		"then":                 true,
		"to":                   true,
		"treat":                true,
		"type":                 true,
		"union":                true,
		"conllu":               true,
		"dep":                  true,
		"ud":                   true,

		"acl":        true,
		"advcl":      true,
		"advmod":     true,
		"amod":       true,
		"appos":      true,
		"aux":        true,
		"case":       true,
		"cc":         true,
		"ccomp":      true,
		"clf":        true,
		"compound":   true,
		"conj":       true,
		"cop":        true,
		"csubj":      true,
		"det":        true,
		"discourse":  true,
		"dislocated": true,
		"expl":       true,
		"fixed":      true,
		"flat":       true,
		"goeswith":   true,
		"iobj":       true,
		"list":       true,
		"mark":       true,
		"nmod":       true,
		"nsubj":      true,
		"nummod":     true,
		"obj":        true,
		"obl":        true,
		"orphan":     true,
		"parataxis":  true,
		"punct":      true,
		"ref":        true,
		"reparandum": true,
		"root":       true,
		"vocative":   true,
		"xcomp":      true,
	}

	filters = map[[3]bool]string{
		[3]bool{true, false, false}: "/(self::node|self::ud|self::dep)",
		[3]bool{false, true, false}: "/self::*[@ud='basic']",
		[3]bool{false, false, true}: "/self::*[@ud='enhanced']",
		[3]bool{true, true, false}:  "/(self::node|self::ud|self::dep|self::*[@ud='basic'])",
		[3]bool{true, false, true}:  "/(self::node|self::ud|self::dep|self::*[@ud='enhanced'])",
		[3]bool{false, true, true}:  "/self::*[@ud]",
	}
)

func xpathcheck(q *Context) {

	if pqxok == "" {
		for _, d := range strings.Split(Cfg.Path, string(filepath.ListSeparator)) {
			xok := filepath.Join(d, "pqxok")
			fi, err := os.Stat(xok)
			if err != nil {
				continue
			}
			if (fi.Mode() | 0111) != 0 {
				pqxok = xok
				break
			}
		}
		if pqxok == "" {
			chLog <- "ERROR: Geen path naar pqxok"
			pqxok = "pqxok"
		}
	}

	contentType(q, "text/plain")

	query := first(q.r, "xpath")
	if query == "" {
		cache(q)
		fmt.Fprintln(q.w, "0")
		return
	}

	nocache(q)

	if strings.Contains(query, "%") {
		rules := getMacrosRules(q)
		query = macroKY.ReplaceAllStringFunc(query, func(s string) string {
			return rules[s[1:len(s)-1]]
		})
	}

	parts := make([]string, 0)
	for _, part := range strings.Split(query, "+|+") {
		part = strings.TrimSpace(part)
		if part == "" || part == "." || part == "/" {
			fmt.Fprintln(q.w, "2")
			return
		} else {
			parts = append(parts, part)
		}
	}
	t, e := exec.Command(pqxok, parts...).Output()
	if e != nil || strings.TrimSpace(string(t)) != "OK" {
		fmt.Fprintln(q.w, "2")
		return
	}

	lvl := 0
PARTLOOP:
	for _, part := range strings.Split(query, "+|+") {

		// geen resultaat -> 1
		for _, s := range reXpath.FindAllString(part, -1) {
			if s == "alpino_ds" {
				continue
			}
			if s[0] == '\'' || s[0] == '"' || s[0] == '$' {
				continue
			}
			if s[0] == '@' {
				if keyTags[s[1:]] {
					continue
				}
				if udTags[s[1:]] {
					continue
				}
				if s == "@type" || s == "@name" || s == "@value" || s == "@sentid" || s == "@cats" || s == "@skips" || s == "@version" {
					continue
				}
				lvl = 1
				break PARTLOOP
			}

			if strings.HasSuffix(s, "(") {
				continue
			}
			if !xpathNames[s] {
				lvl = 1
				break PARTLOOP
			}
		}

	}

	fmt.Fprintln(q.w, lvl)
}

func xpathdl(q *Context) {

	if sysErr(q.r.ParseForm()) {
		return
	}

	query := firstf(q.form, "xpath")
	prefix := firstf(q.form, "db")
	methode := firstf(q.form, "mt")
	step, err := strconv.Atoi(firstf(q.form, "step"))
	if err != nil || step < 1 {
		step = 1
	}
	gettextcount := firstf(q.form, "what") == "textcount"

	var chClose <-chan bool
	if f, ok := q.w.(http.CloseNotifier); ok {
		chClose = f.CloseNotify()
	} else {
		chClose = make(<-chan bool)
	}

	q.w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	q.w.Header().Set("Content-Disposition", "attachment; filename=uitvoer.txt")

	//curno, hash, methode, loading, errval := xpath_do_search(q, query, prefix, methode, 0, 99999999, chClose, false, step, textcount)
	_, _, _, _, _ = xpath_do_search(q, query, prefix, methode, 0, 99999999, chClose, false, step, gettextcount)
}

// TAB: xpath
func xpath(q *Context) {

	prefix := getprefix(q)
	if !q.prefixes[prefix] {
		http.Error(q.w, "Invalid corpus: "+prefix, http.StatusPreconditionFailed)
		return
	}

	xpathmax := getxpathmax(q)

	methode := first(q.r, "mt")
	if methode != "dx" {
		methode = "std"
	}

	var errval error
	var db *dbxml.Db
	var docs *dbxml.Docs
	var loading bool
	defer func() {
		if docs != nil {
			docs.Close()
		}
		if db != nil {
			db.Close()
		}
		if loading {
			clearLoading(q.w)
		}
		if errval != nil {
			fmt.Fprintf(q.w, "<div class=\"error\">FOUT: %s</div>\n", html.EscapeString(errval.Error()))
		}
		html_footer(q)
	}()

	// HTML-uitvoer van begin van de pagina
	writeHead(q, "XPath", 2)
	html_xpath_header(q)

	// HTML-uitvoer van het formulier
	// Returnwaarde is true als er een query was gedefinieerd
	has_query, fff := html_xpath_form(q, xpathmax)

	// Als er geen query is gedefinieerd, HTML-uitvoer van korte helptekst, pagina-einde, en exit
	if !has_query {
		html_xpath_uitleg(q)
		return
	}

	var chClose <-chan bool
	if f, ok := q.w.(http.CloseNotifier); ok {
		chClose = f.CloseNotify()
	} else {
		chClose = make(<-chan bool)
	}

	_, errval = sqlDB.Exec(fmt.Sprintf("UPDATE `%s_info` SET `active` = NOW() WHERE `id` = %q", Cfg.Prefix, prefix))
	if logerr(errval) {
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

	// hier begint het

	query := first(q.r, "xpath")
	oriquery := query
	query += filters[fff]
	now := time.Now()
	curno, hash, methode, loading, errval := xpath_do_search(q, query, prefix, methode, offset, xpathmax, chClose, true, 1, false)
	if errval != nil {
		return
	}

	// hier stopt het

	if curno == 0 {
		// commentaar voor SPOD
		fmt.Fprintln(q.w, "<!--NOMATCH-->\ngeen match gevonden")
	}

	// Links naar volgende en vorige pagina's met resultaten
	qs := "xpath=" + urlencode(oriquery) + "&amp;mt=" + methode
	if fff[0] {
		qs += "&A=1"
	}
	if fff[1] {
		qs += "&U=1"
	}
	if fff[2] {
		qs += "&E=1"
	}
	if offset > 0 || curno > offset+xpathmax {
		if offset > 0 {
			fmt.Fprintf(q.w, "<a href=\"xpath?%s&amp;offset=%d\">vorige</a>", qs, offset-xpathmax)
		} else {
			fmt.Fprint(q.w, "vorige")
		}
		fmt.Fprint(q.w, " | ")
		if curno > offset+xpathmax {
			fmt.Fprintf(q.w, "<a href=\"xpath?%s&amp;offset=%d\">volgende</a>", qs, offset+xpathmax)
		} else {
			fmt.Fprint(q.w, "volgende")
		}
	}

	if curno > 0 {
		fmt.Fprintf(q.w, `<p>
<form action="dlx" method="POST" accept-charset="UTF-8" enctype="multipart/form-data">
<input type="hidden" name="xpath" value="%s">
<input type="hidden" name="db" value="%s">
<input type="hidden" name="mt" value="%s">
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
</select>
<input type="submit" value="downloaden">
</form>
`,
			html.EscapeString(query),
			html.EscapeString(prefix),
			methode)
	}

	if q.auth && curno > 0 {
		fmt.Fprintf(q.w, `<p>
<form action="xsavez" method="POST" accept-charset="UTF-8" enctype="multipart/form-data">
<input type="hidden" name="xpath" value="%s">
<input type="hidden" name="db" value="%s">
<input type="hidden" name="mt" value="%s">
<input type="submit" value="nieuw corpus maken op basis van deze zoekopdracht">
</form>
`,
			html.EscapeString(query),
			html.EscapeString(prefix),
			methode)
	}

	fmt.Fprintln(q.w, "<hr><small>tijd:", tijd(time.Now().Sub(now)), "</small>")

	if curno == 0 {
		return
	}

	fmt.Fprintln(q.w, "<hr>")

	var metas []MetaType
	if q.hasmeta[prefix] {
		metas = getMeta(q, prefix)
	}

	// Links naar statistieken
	fmt.Fprintf(q.w, `
        <p>
		<div id="xstats">
		<form action="javascript:$.fn.xpathstats()" name="xstatsform">
        <input type="hidden" name="hash" value="%s">
		<input type="hidden" name="xpath" value="%s">
		<input type="hidden" name="db" value="%s">
		<input type="hidden" name="mt" value="%s">
		Selecteer nul tot vijf attributen:
        <p>
`, hash, html.EscapeString(query), html.EscapeString(prefix), methode)

	for i := 1; i <= 5; i++ {

		fmt.Fprintf(q.w, "<select name=\"attr%d\">\n<option value=\"\">--</option>\n", i)
		if q.hasmeta[prefix] {
			fmt.Fprintln(q.w, "<optgroup label=\"&mdash; metadata &mdash;\">")
			for _, m := range metas {
				fmt.Fprintf(q.w, "<option value=\":%s\">%s</option>\n", html.EscapeString(m.name), html.EscapeString(m.name))
			}
			fmt.Fprintln(q.w, "</optgroup>")
		}
		fmt.Fprintln(q.w, "<optgroup label=\"&mdash; node: attributen &mdash;\">")
		for _, s := range pqnode.NodeTags {
			fmt.Fprintf(q.w, "<option>%s</option>\n", s)
		}
		fmt.Fprintln(q.w, "<optgroup label=\"&mdash; ud: attributen &mdash;\">")
		for _, s := range UdTags {
			fmt.Fprintf(q.w, "<option value=\"ud:%s\">%s</option>\n", s, s)
		}
		fmt.Fprintln(q.w, "</optgroup>")
		fmt.Fprintln(q.w, "<optgroup label=\"&mdash; dep: attributen &mdash;\">")
		for _, s := range DepTags {
			fmt.Fprintf(q.w, "<option value=\"dep:%s\">%s</option>\n", s, s)
		}
		fmt.Fprintln(q.w, "</optgroup>")
		fmt.Fprintln(q.w, "</select>")
	}

	fmt.Fprint(q.w, `
		<p>
		<input type="submit" value="doe telling">
		</form>
		<p>
        <iframe src="leeg.html" id="xframe" class="hide"></iframe>
        <div id="result" class="hide"></div>
`)
	if q.hasmeta[prefix] {
		metahelp(q)
		fmt.Fprintln(q.w, `<p>
            <div id="statsmeta" class="hide">
            <div id="innermetatop"></div>
            <div id="metacount" class="hide">
            <table>
            <tr><td>items:<td class="right" id="metacount1">
            <tr><td>zinnen:<td class="right" id="metacount2">
            </table>
            </div>
            <div id="innermeta"></div>
            <img src="busy.gif" id="busymeta" class="hide" alt="aan het werk..." style="margin-top:1em">
            </div>`)
	}
	fmt.Fprintln(q.w, "</div>")

}

func getxpathmax(q *Context) int {
	xn := first(q.r, "xn")
	if xn == "" {
		if x, err := q.r.Cookie("paqu-xn"); err == nil {
			xn = x.Value
		}
	}

	xi, err := strconv.Atoi(xn)
	if err != nil {
		xi = ZINMAX * 2
	} else if xi < 10 {
		xi = 10
	} else if xi > 500 {
		xi = 500
	}
	xn = fmt.Sprint(xi)

	exp := time.Now().AddDate(0, 0, 14)
	http.SetCookie(q.w, &http.Cookie{Name: "paqu-xn", Value: xn, Path: cookiepath, Expires: exp})
	return xi
}

//. HTML

func html_xpath_header(q *Context) {
	fmt.Fprint(q.w, `
<script type="text/javascript" src="jquery.js"></script>
<script type="text/javascript"><!--

  hexEncode = function(s){
    var hex, i;
    var result = "";
    for (i = 0; i < s.length; i++) {
        hex = s.charCodeAt(i).toString(16);
        result += ("000" + hex).slice(-4);
    }
    return result.toLowerCase();
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

  var metarun = true;
  $.fn.xpathstats = function() {

    var n = 0;
    var val = "";
    if (at1.selectedIndex > 0) { n++; val = at1.value; }
    if (at2.selectedIndex > 0) { n++; val = at2.value; }
    if (at3.selectedIndex > 0) { n++; val = at3.value;  }
    if (at4.selectedIndex > 0) { n++; val = at4.value;  }
    if (at5.selectedIndex > 0) { n++; val = at5.value;  }
    /*
    if (n < 1) {
      alert("Geen attribuut geselecteerd");
      return;
    }
    */
    setCookie("paqu-xpattr1", hexEncode(at1.value), 14);
    setCookie("paqu-xpattr2", hexEncode(at2.value), 14);
    setCookie("paqu-xpattr3", hexEncode(at3.value), 14);
    setCookie("paqu-xpattr4", hexEncode(at4.value), 14);
    setCookie("paqu-xpattr5", hexEncode(at5.value), 14);

    if (n == 1 && val.substring(0, 1) == ":") {
        val = val.substring(1)
        $('#result').addClass('hide');
        if (!metarun) {
            $('.metasub').addClass('hide');
        }
        $('#statsmeta').removeClass('hide');
        if (metarun) {
            $('#xframe').attr('src', "xstatsmeta?item=" + val + "&" + $(document.xstatsform).serialize());
        } else {
            $('#meta' + hexEncode(val)).removeClass('hide');
        }
    } else {
        $('#statsmeta').addClass('hide');
        $('#result').html('');
        $('#result').removeClass('hide');
        $('#xframe').attr('src', "xpathstats?" + $(document.xstatsform).serialize());
    }

  }

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

  function fillMeta(idx) {
    var ida = "#meta" + idx + "a";
    var idb = "#meta" + idx + "b";
    var fl = metavars[idx].fl;
    var lbl = metavars[idx].lbl;
    var max = metavars[idx].max;
    var da = metavars[idx].a;
    var db = metavars[idx].b;
    var ac = metavars[idx].ac;
    var bc = metavars[idx].bc;
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
         metavars[idx].a = sortMeta(da, 0);
         metavars[idx].ac = 0;
         fillMeta(idx);
      });
    $(ida + ' td.b').on('click', function() {
         metavars[idx].a = sortMeta(da, 1);
         metavars[idx].ac = 1;
         fillMeta(idx);
      });
    $(idb + ' td.a').on('click', function() {
         metavars[idx].b = sortMeta(db, 0);
         metavars[idx].bc = 0;
         fillMeta(idx);
      });
    $(idb + ' td.b').on('click', function() {
         metavars[idx].b = sortMeta(db, 1);
         metavars[idx].bc = 1;
         fillMeta(idx);
      });
    $(idb + ' td.c').on('click', function() {
         metavars[idx].b = sortMeta(db, 2);
         metavars[idx].bc = 2;
         fillMeta(idx);
      });
  }

  var metavisible = false;
  function metahelp() {
    var e = $("#helpmeta");
    e.show();
    e.css("zIndex", 9999);
    metavisible = true;
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
    });

  function formclear(f) {
    f.xpath.value = "";
    f.A.checked = false;
    f.U.checked = false;
    f.E.checked = false;
    xquery.css('background-color', '#ffffff');
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

  var resultmetatop;
  var resultmeta;
  var busymeta;
  var metacount;
  var metacount1;
  var metacount2;
  var result;
  var at1, at2, at3, at4, at5;
  var xquery;
  var macrotext;
  var metadn = 0;
  var metavars = [];
  var aligns;
  var labels;
  var isidx;
  var data;
  var curcol;
  var hasattribs;

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

  function resultset(c) {
    data.lines.sort(function(a, b) {
        if (c == 0) {
            var r = b[0][1] - a[0][1];
            if (r != 0) { return r; }
        } else if (isidx[c]) {
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
        for (i in isidx) {
            if (i == c) { continue; }
            if (i == 0) {
                var r = b[0][1] - a[0][1];
                if (r != 0) { return r; }
            } else if (isidx[i]) {
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
    curcol = c;
    fillresult();
  }

  function fillresult() {
     if (data.toomany) {
         result.prepend('<div class="warning">Onderbroken vanwege te veel combinaties</div>');
         data.toomany = false;
     }
     $('#resultmatches').html(data.matches);
     $('#resultlinecount').html(data.linecount);
     if (hasattribs) {
         $('#resultcombis').html(data.combis);
     }
     $('#resulttijd').html(data.tijd);
     if (data.final) {
         $('#resultbusy').html('');
     } else {
         if (typeof data.perc !== 'undefined') {
             $('#resultbusy').html('<img src="busy.gif" alt="aan het werk..."> ' + data.perc);
         }
     }
     if (!hasattribs) {
         return;
     }
     var t = $('#resultlines');
     var em
     if (curcol == 0) { em = ' em'; } else { em = ""; }
     var s = '<tr class="odd"><th id="c0" class="link' + em + '">items<th>';
     var i;
     for (i = 1; i < aligns.length; i++) {
         if (curcol == i) { em = ' em'; } else { em = ""; }
         s += '<th id="c' + i + '" class="link ' + aligns[i] + em + '">' + labels[i];
     }
     t.html(s);
     $('#c0').on('click', function() { resultset(0); });
     $('#c1').on('click', function() { resultset(1); });
     $('#c2').on('click', function() { resultset(2); });
     $('#c3').on('click', function() { resultset(3); });
     $('#c4').on('click', function() { resultset(4); });
     $('#c5').on('click', function() { resultset(5); });
     var odd;
     for (i in data.lines) {
         if (i % 2 == 1) { odd = ' class="odd"'; } else { odd = ""; }
         if (i == 250) {
             s = '<tr' + odd + '><td><td>';
             for (j = 1; j < aligns.length; j++) {
                 s += '<td class="' + aligns[j] + '">...';
             }
             t.append(s);
             break;
         }
         s = '<tr' + odd + '><td class="right' + '">' + ival(data.lines[i][0][1]) + '<td class="right">' + data.lines[i][0][0];
         for (j = 1; j < aligns.length; j++) {
             var w = data.lines[i][j][0];
             var cl = "";
             if (w.substr(0, 2) == "  ") {
                 cl = " multi";
             } else if (w == "") {
                 cl = " nil";
                 w = "(leeg)";
             }
             s += '<td class="' + aligns[j] + cl + '">' + escapeHtml(w);
         }
         t.append(s);
     }
  }

  window._fn = {
    setmetaval: function(value) {
      metadn = value;
    },
    setmetavars: function(idx, lbl, fl, max, ac, bc) {
      metavars[idx] = {};
      metavars[idx].lbl = lbl;
      metavars[idx].fl = fl;
      metavars[idx].max = max;
      metavars[idx].ac = ac;
      metavars[idx].bc = bc;
    },
    setmetalines: function(idx, a, b) {
      metavars[idx].a = a;
      metavars[idx].b = b;
    },
    makemetatable: function(idx) {
      fillMeta(idx);
    },
    updatemeta: function(data) {
      resultmeta.append(data);
    },
    updatemetatop: function(data) {
      resultmetatop.append(data);
    },
    countmeta: function(i, j) {
      metacount1.html(i);
      metacount2.html(j);
    },
    startedmeta: function() {
      resultmeta.html('');
      busymeta.removeClass('hide');
      metacount.removeClass('hide');
    },
    completedmeta: function() {
      busymeta.addClass('hide');
      metarun = false;
    },
    init: function(o) {
      curcol = 0;
      aligns = o.aligns;
      labels = o.labels;
      isidx = o.isidx;
      hasattribs = (labels.length > 1);
      if (!hasattribs) {
        result.html(
'<table>\n' +
'<tr><td>items:<td class="right" id="resultmatches">0<td rowspan="3" id="resultbusy"><img src="busy.gif" alt="aan het werk...">\n' +
'<tr><td>zinnen:<td class="right" id="resultlinecount">0\n' +
'<tr><td>tijd:<td class="right" id="resulttijd">0s\n' +
'</table>\n');
      } else {
        result.html(
'<table>\n' +
'<tr><td>items:<td class="right" id="resultmatches">0<td rowspan="3" id="resultbusy"><img src="busy.gif" alt="aan het werk...">\n' +
'<tr><td>zinnen:<td class="right" id="resultlinecount">0\n' +
'<tr><td>combinaties:<td class="right" id="resultcombis">0\n' +
'<tr><td>tijd:<td class="right" id="resulttijd">0s\n' +
'</table>\n' +
'<p>\n' +
'<table class="breed" id="resultlines">\n' +
'</table>\n' +
'<hr><a href="xpathstats?' + o.download + '">download</a>\n');
      }
    },
    update: function(o) {
        data = o;
        if (curcol == 0) {
            fillresult();
        } else {
            resultset(curcol);
        }
    },
    error: function(o) {
        $('#resultbusy').html('');
        result.prepend('<div class="error">Error: ' + o + '</div>');
    },
    update2: function(data) {
      if (data.err == "") {
          $('#macromsg').slideUp(200);
          $('#macrotext').text(data.macros);
          $('#macrotext').val(data.macros);
          macros = data.keys;
      } else {
          $('#macromsg').text("Fout: " + data.err).slideDown(400);
      }
      disableSave();
    }
  }
  function xstatftest() {
    var n = 0;
    if (at1.selectedIndex > 0) { n++; }
    if (at2.selectedIndex > 0) { n++; }
    if (at3.selectedIndex > 0) { n++; }
    if (at4.selectedIndex > 0) { n++; }
    if (at5.selectedIndex > 0) { n++; }
    /*
    if (n < 1) {
      alert("Geen attribuut geselecteerd");
      return false;
    }
    */
    setCookie("paqu-xpattr1", hexEncode(at1.value), 14);
    setCookie("paqu-xpattr2", hexEncode(at2.value), 14);
    setCookie("paqu-xpattr3", hexEncode(at3.value), 14);
    setCookie("paqu-xpattr4", hexEncode(at4.value), 14);
    setCookie("paqu-xpattr5", hexEncode(at5.value), 14);
    return true;
  }
  function setForm() {
    at1.selectedIndex = 0;
    at2.selectedIndex = 0;
    at3.selectedIndex = 0;
    at4.selectedIndex = 0;
    at5.selectedIndex = 0;
    try {
      var a = getCookie("paqu-xpattr1");
      if (a != "" ) {
          at1.value = hexDecode(a);
      }
    } catch (e) { }
    try {
      var a = getCookie("paqu-xpattr2");
      if (a != "" ) {
          at2.value = hexDecode(a);
      }
    } catch (e) { }
    try {
      var a = getCookie("paqu-xpattr3");
      if (a != "" ) {
          at3.value = hexDecode(a);
      }
    } catch (e) { }
    try {
      var a = getCookie("paqu-xpattr4");
      if (a != "" ) {
          at4.value = hexDecode(a);
      }
    } catch (e) { }
    try {
      var a = getCookie("paqu-xpattr5");
      if (a != "" ) {
          at5.value = hexDecode(a);
      }
    } catch (e) { }
  }

  var lastcall;
  var timer;
  var reMacro = /%[a-zA-Z][-_a-zA-Z0-9]*%/;
  function qcheck() {
    try {
      window.clearTimeout(timer);
    } catch (e) { }
    timer = window.setTimeout(function(){qcheckdo()}, 200);
  }
  function qcheckdo() {
    if (lastcall) {
      try {
        lastcall.abort();
      }
      catch(err) {}
    }
    if (reMacro.test(xquery.val())) {
       $('#btExpand').show(400);
    } else {
       $('#btExpand').hide(200);
    }
    lastcall = $.ajax("xpathcheck?" + xquery.serialize())
      .done(function(data) {
        r = parseInt(data);
        if (r == 0) {
          xquery.css('background-color', '#ffffff');
        } else if (r == 1) {
          xquery.css('background-color', '#ffff80');
        } else if (r == 2) {
          xquery.css('background-color', '#ffa0a0');
        }
      }).fail(function(e) {
          xquery.setStyle('background-color', '#d0d0d0');
      })
      .always(function() {
        lastcall = null;
      });
  }

  var expandLvl = 1;
  var macroOri = "";
  function macroExpand(reset) {
    if (reset) {
        expandLvl = 1;
        macroOri = xquery.serialize();
    } else {
        expandLvl++;
    }
    $('#macroOuter').slideUp(200);
    $.ajax("macroexpand?lvl=" + expandLvl + "&" + macroOri)
      .done(function(data) {
         $("#macroInner").text(data);
         $('#macroOuter').slideDown(400);
         if (reMacro.test(data)) {
           $('#btExpandNxt').show(400);
         } else {
           $('#btExpandNxt').hide(200);
         }
      }).fail(function(e) {
         $("#macroInner").text(e.responseText);
         $('#macroOuter').slideDown(400);
      })
  }

  var macroIsOpen = false;
  function openMacro() {
    if (macroIsOpen) {
      $("#openmacro").text("Macro's")
      $('#macros').slideUp(200);
      macroIsOpen = false;
    } else {
      msetsize();
      $("#openmacro").text("Macro's sluiten")
      $('#macros').slideDown(400);
      macroIsOpen = true;
    }
  }

  function enableSave() {
    $('#macrosave').removeAttr('disabled');
    $('#macrosave').addClass('bold');
  }

  function disableSave() {
    $('#macrosave').attr("disabled", "disabled");
    $('#macrosave').removeClass('bold');
  }

  function xresize() {
    localStorage.setItem(
    "paqu-xpath-x",
    JSON.stringify({
      w: xquery.outerWidth(),
      h: xquery.outerHeight()
    }));
  }

  function mresize() {
    if (macroIsOpen) {
      localStorage.setItem(
        "paqu-xpath-m",
        JSON.stringify({
          w: macrotext.outerWidth(),
          h: macrotext.outerHeight()
        }));
    }
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

  function xsetsize() {
    var storageContent = localStorage.getItem("paqu-xpath-x");
    if (storageContent !== undefined) {
      var d = JSON.parse(storageContent) || {};
      if (d['w']) {
        xquery.outerWidth(limit(d['w'], 400, 1000));
      }
      if (d['h']) {
        xquery.outerHeight(limit(d['h'], 80, 600));
      }
    }
  }

  function msetsize() {
    var storageContent = localStorage.getItem("paqu-xpath-m");
    if (storageContent !== undefined) {
      var d = JSON.parse(storageContent) || {};
      if (d['w']) {
        macrotext.outerWidth(limit(d['w'], 400, 1000));
      }
      if (d['h']) {
        macrotext.outerHeight(limit(d['h'], 80, 600));
      }
    }
  }

  function init() {
    xquery = $('#xquery');
    macrotext = $('#macrotext');
    xsetsize();
    xquery.on('keyup', qcheck);
    xquery.on('change', qcheck);
    xquery.on('click', qcheck);
    xquery.on('mouseup', xresize);
    macrotext.on('mouseup', mresize);
    qcheckdo();
    $('#openmacro').on('click', openMacro);
    $('#macroreset').on('click', disableSave);
    $('#macrotext').on('input propertychange', enableSave);
    $('#macrofilename').on('change', function() {
        if ($('#macrofilename').val() == "") {
          $('#macrofilesave').attr('disabled', 'disabled');
          $('#macrofilesave').removeClass('bold');
        } else {
          $('#macrofilesave').removeAttr('disabled');
          $('#macrofilesave').addClass('bold');
        }
    });
    $('#macrofilesave').on('click', function() {
          $('#macrofilesave').removeClass('bold');
    });
    $('#btExpand').on('click', function() { macroExpand(true) });
    $('#btExpandNxt').on('click', function() { macroExpand(false) });
    $('#btClose').on('click', function() { $('#macroOuter').slideUp(200); });
    $('#btXCopy').on('click', function() {
        var txt = $('#macroInner').text();
        if (txt.length > 1200) {
            txt = txt.replace(/   */g, ' ');
        }
        if (txt.length > 1200) {
            alert("De tekst is te lang voor het invoerveld");
        } else {
            $('#xquery').val(txt);
            $('#macroOuter').slideUp(200);
            if (!reMacro.test(txt)) {
              $('#btExpand').hide(200);
            }
        }
    });
  }

  $(document).ready(function() {
    result = $('#result');
    resultmeta = $('#innermeta');
    resultmetatop = $('#innermetatop');
    busymeta = $('#busymeta');
    metacount = $('#metacount');
    metacount1 = $('#metacount1');
    metacount2 = $('#metacount2');
    try {
      var f = document.forms["xstatsform"];
      at1 = f["attr1"];
      at2 = f["attr2"];
      at3 = f["attr3"];
      at4 = f["attr4"];
      at5 = f["attr5"];
      setForm();
    } catch (e) {}
  });

  //--></script>
`)
}

func html_xpath_uitleg(q *Context) {
	fmt.Fprint(q.w, `
<p>
<hr>
<p>
PaQu ondersteunt XPath2 met dezelfde extensies als Dact: macro's en pipelines.
<p>
Voorbeelden, zie:
<a href="http://rug-compling.github.io/dact/cookbook/" target="_blank">Dact Cookbook</a>
`)
}

func html_xpath_form(q *Context, xpathmax int) (has_query bool, filter [3]bool) {
	has_query = true
	if first(q.r, "xpath") == "" {
		has_query = false
	}
	methode := first(q.r, "mt")

	filter = [3]bool{
		first(q.r, "A") != "",
		first(q.r, "U") != "",
		first(q.r, "E") != "",
	}

	if q.auth {
		macros := ""
		rows, err := sqlDB.Query(fmt.Sprintf("SELECT `macros` FROM `%s_macros` WHERE `user` = %q", Cfg.Prefix, q.user))
		if err == nil {
			if rows.Next() {
				rows.Scan(&macros)
				rows.Close()
			}
		}
		fmt.Fprintf(q.w, `
<button id="openmacro">Macro's</button>
<div id="macros" style="display:none">
<div class="warning" id="macromsg" style="display:none"></div>
<p>
<form action="savemacros" method="post" target="hiddenframe" enctype="multipart/form-data">
<input type="submit" value="Uploaden" id="macrofilesave" disabled="disabled">
<input type="file" name="macrotext" id="macrofilename">
</form>
<p>
<form action="savemacros" method="post" target="hiddenframe" enctype="multipart/form-data">
<textarea rows="6" cols="80" id="macrotext" name="macrotext">%s</textarea><br>
<input type="submit" value="Opslaan" id="macrosave" disabled="disabled">
<input type="reset" value="Reset" id="macroreset">
<input type="button" value="Download" onclick="window.location.assign('downloadmacros')">
</form>
<p>
<hr>
</div>
<iframe src="leeg.html" name="hiddenframe" class="hide"></iframe>
<p>
`, html.EscapeString(macros))
	}

	fmt.Fprint(q.w, `
<form action="xpath" method="get" accept-charset="utf-8">
<a href="corpusinfo">[?]</a> corpus: <select name="db">
`)
	html_opts(q, q.opt_db, getprefix(q), "corpus")
	fmt.Fprintln(q.w, "</select>")
	if q.auth {
		fmt.Fprintln(q.w, "<a href=\"corpuslijst\">meer/minder</a>")
	}
	fmt.Fprintf(q.w, `<p>
		XPATH query (<a href="http://rug-compling.github.io/dact/cookbook/" target="_blank">voorbeelden</a>)
(<a href="macros.txt" target="_blank">ingebouwde macro's</a>):
<br>
		<textarea name="xpath" rows="6" cols="80" maxlength="1200" id="xquery">%s</textarea>
		<p>
		`, html.EscapeString(first(q.r, "xpath")))
	fmt.Fprintf(q.w, `filter:
<input type="checkbox" id="cbA" name="A"%v> <label for="cbA">Alpino</label> &nbsp;
<input type="checkbox" id="cbU" name="U"%v> <label for="cbU">Basic UD</label> &nbsp;
<input type="checkbox" id="cbE" name="E"%v> <label for="cbE">Enhanced UD</label>
(<a href="info.html#subset" target="_blank">uitleg</a>)
<p>
`,
		ifelse(filter[0], " checked", ""),
		ifelse(filter[1], " checked", ""),
		ifelse(filter[2], " checked", ""))
	if Cfg.Dactx {
		selected := ""
		if methode == "dx" {
			selected = " selected"
		}
		fmt.Fprintf(q.w, `methode: <select name="mt">
			<option value="std">standaard</option>
			<option value="dx"%s>ge&euml;xpandeerde indexnodes</option>
			</select> (<a href="info.html#expanded" target="_blank">uitleg</a>)
			<p>`, selected)
	}
	fmt.Fprintln(q.w, `aantal: <select name="xn">`)
	for _, i := range []int{10, 20, 50, 100, 200, 500} {
		if i == xpathmax {
			fmt.Fprintf(q.w, "<option selected>%d</option>\n", i)
		} else {
			fmt.Fprintf(q.w, "<option>%d</option>\n", i)
		}
	}
	fmt.Fprint(q.w, `
		</select>
		<p>
           <input type="submit" value="Zoeken">
           <input type="button" value="Wissen" onClick="javascript:formclear(form)">
           <input type="reset" value="Reset" onClick="javascript:qcheck()">
           <input type="button" id="btExpand" value="Toon macro-expansie" style="display:none">
       </form>
       <div id="macroOuter" style="display:none">
       <div id="macroInner"></div>
       <button id="btClose">Sluiten</button>
       <button id="btExpandNxt" style="display:none">Toon verdere macro-expansie</button>
       <button id="btXCopy">Kopieer naar invoer</button>
       </div>
`)
	fmt.Fprint(q.w, `<script type="text/javascript" src="jquery.textcomplete.js"></script>
<script type="text/javascript"><!--
var attribs = ["@cats","@skips","@name","@value","@sentid"`)
	for _, a := range pqnode.NodeTags {
		fmt.Fprintf(q.w, ",\n%q", "@"+a)
	}
	for a := range udTags {
		fmt.Fprintf(q.w, ",\n%q", "@"+a)
	}
	fmt.Fprint(q.w, `].sort();

var axis = [`)
	p := ""
	for key := range xpathNames {
		fmt.Fprintf(q.w, "%s\n%q", p, key)
		p = ","
	}
	fmt.Fprint(q.w,
		`].sort();

var macros = [`)

	keys := getMacrosKeys(q)
	p = ","
	for i, key := range keys {
		if i == len(keys)-1 {
			p = ""
		}
		fmt.Fprintf(q.w, "%q%s\n", key, p)
	}

	fmt.Fprint(q.w, `];

function outText(text) {
    var state = 0;
    var i, j;
    for (i = 0, j = text.length; i < j; i++) {
        var c = text.charAt(i);
        if (state == 0) {
           if (c == "'") {
               state = 1;
           } else if (c == '"') {
               state = 2;
           }
        } else if (state == 1) {
           if (c == "'") {
               state = 0;
           }
        } else {
           if (c == '"') {
               state = 0;
           }
        }
    }
    return state == 0;
}

function inMacro(text) {
    var state = false;
    var i, j;
    for (i = 0, j = text.length; i < j; i++) {
        if (text.charAt(i) == "%") {
            state = !state;
        }
    }
    return state;
}

$('#xquery').textcomplete([
{
    match: /(@[-_a-zA-Z0-9]*)$/,
    search: function (term, callback) {
        callback($.map(attribs, function (e) {
            return e.indexOf(term) === 0 ? e : null;
        }));
    },
    replace: function (value) {
        return value + ' ';
    },
    index: 1,
    context: outText
},
{
    match: /%(|[a-zA-Z][_a-zA-Z0-9]*)$/,
    search: function (term, callback) {
        callback($.map(macros, function (e) {
            return e.indexOf(term) === 0 ? e : null;
        }));
    },
    replace: function (value) {
        return '%' + value + '% ';
    },
    index: 1,
    context: inMacro
},
{
    match: /\b([a-z][-a-z]*)$/,
    search: function (term, callback) {
        callback($.map(axis, function (e) {
            return e.indexOf(term) === 0 ? e : null;
        }));
    },
    replace: function (value) {
        if (value.indexOf(":") > 0) {
            return value;
        } else {
            return value + ' ';
        }
    },
    index: 1,
    context: outText
},
{
    match: /^((.|\n)*(['")\]]|@[-_a-zA-Z0-9]+ |[%0-9] ))$/,
    search: function (term, callback) {
        var chars = [];
        var chars2 = [];
        var i, j;
        var state = 0;

        for (i = 0, j = term.length; i < j; i++) {
            var c = term.charAt(i);
            if (state == 0) {
               if (c == "'") {
                   state = 1;
               } else if (c == '"') {
                   state = 2;
               } else if (c == '[') {
                   chars.unshift(']');
               } else if (c == ']') {
                   if (chars.length == 0 || chars[0] != ']') {
                       callback([]);
                       return;
                   }
                   chars.shift();
               } else if (c == '(') {
                   chars.unshift(')');
               } else if (c == ')') {
                   if (chars.length == 0 || chars[0] != ')') {
                       callback([]);
                       return;
                   }
                   chars.shift();
               }
            } else if (state == 1) {
               if (c == "'") {
                   state = 0;
               }
            } else {
               if (c == '"') {
                   state = 0;
               }
            }
        }

        var alltext = document.getElementById('xquery').value;
        state = 0;
        for (i = alltext.length-1, j = term.length; i >= j; i--) {
            var c = alltext.charAt(i);
            if (state == 0) {
               if (c == "'") {
                   state = 1;
               } else if (c == '"') {
                   state = 2;
               } else if (c == ']') {
                   chars2.unshift(']');
               } else if (c == '[') {
                   if (chars2.length == 0 || chars2[0] != ']') {
                       chars2 = [];
                   } else {
                       chars2.shift();
                   }
               } else if (c == ')') {
                   chars2.unshift(')');
               } else if (c == '(') {
                   if (chars2.length == 0 || chars2[0] != ')') {
                       chars2 = [];
                   } else {
                       chars2.shift();
                   }
               }
            } else if (state == 1) {
               if (c == "'") {
                   state = 0;
               }
            } else {
               if (c == '"') {
                   state = 0;
               }
            }
        }

        var rt = [];
        if (chars.length > 0 && chars2.length > 0) {
             if (chars2[0] != chars[0]) {
                rt.push(chars[0]);
             }
        } else {
            var s = "";
            while (chars.length > 0) {
                s += chars.shift() + " ";
                rt.push(s);
            }
        }
        callback(rt);
    },
    index: 1,
    replace: function (value) {
        return "$1" + value;
    },
    context: outText
}],
{
    maxCount: 200,
    debounce: 300
});

init();
//--></script>
`)
	return
}

func xpath_result(q *Context, curno int, dactfile, filename, xmlall string, xmlparts []string, prefix string, global bool) {

	fout := func(err error, s string) {
		s2 := fmt.Sprintf("<li>FOUT bij parsen van XML: %s\n<pre>%s</pre></li>\n", html.EscapeString(err.Error()), html.EscapeString(s))
		fmt.Fprintf(q.w, `<script type="text/javascript"><!--
$('ol').append(%q);
//--></script>
`, s2)
		if ff, ok := q.w.(http.Flusher); ok {
			ff.Flush()
		}
	}

	seen := make(map[string]bool)
	alpino := Alpino_ds{}
	err := xml.Unmarshal([]byte(xmlall), &alpino)
	if err != nil {
		fout(err, xmlall)
		return
	}
	woorden := strings.Fields(alpino.Sentence)

	lvl := make([]int, len(woorden)+1)
	ids := make([]string, len(xmlparts))
	ud1 := make([]string, 0)
	ud2 := make([]string, 0)

	for i, part := range xmlparts {

		var isUd, isDep, isId, isEid bool
		var ID, Head, Deprel string

		if strings.HasPrefix(part, "<node") {
			// isNode
		} else if strings.HasPrefix(part, "<ud") {
			isUd = true
		} else if strings.HasPrefix(part, "<dep") {
			isDep = true
		} else {
			var alpino_test Alpino_test
			err := xml.Unmarshal([]byte(part), &alpino_test)
			if err != nil {
				fout(err, part)
				return
			}
			if alpino_test.Id != "" {
				ID = alpino_test.Id
				Head = alpino_test.Head
				Deprel = alpino_test.Deprel
				if alpino_test.Ud == "enhanced" {
					isEid = true
				} else {
					isId = true
				}
			}
		}

		alp := Alpino_ds{}

		if isId {
			alp.Node0 = findUdId(alpino.Node0, ID)
			isUd = true
		} else if isEid {
			alp.Node0 = findDepId(alpino.Node0, ID, Head, Deprel)
			isDep = true
		}

		if isUd {
			if !isId {
				err := xml.Unmarshal([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<alpino_ds version="1.3">
<node>
`+part+`
</node>
</alpino_ds>`), &alp)
				if err != nil {
					fout(err, part)
					return
				}
			}
			ud1 = append(ud1, alp.Node0.Ud.Id+":"+alp.Node0.Ud.Head+":"+alp.Node0.Ud.Deprel)
			if i, err := strconv.Atoi(alp.Node0.Ud.Id); err == nil && i > 0 && i <= len(woorden) {
				lvl[i-1]++
			}
		} else if isDep {
			if !isEid {
				err := xml.Unmarshal([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<alpino_ds version="1.3">
<node>
<ud>
`+part+`
</ud>
</node>
</alpino_ds>`), &alp)
				if err != nil {
					fout(err, part)
					return
				}
			}
			for _, dep := range alp.Node0.Ud.Dep {
				ud2 = append(ud2, dep.Id+":"+dep.Head+":"+dep.Deprel)
				if i, err := strconv.Atoi(dep.Id); err == nil && i > 0 && i <= len(woorden) {
					lvl[i-1]++
				}
			}
		} else {
			err := xml.Unmarshal([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<alpino_ds version="1.3">
`+part+`
</alpino_ds>`), &alp)
			if err != nil {
				fout(err, part)
				return
			}
			if alp.Node0 != nil {
				ids[i] = alp.Node0.Id
				if i, err := strconv.Atoi(alp.Node0.Index); err == nil && alp.Node0.Word == "" && len(alp.Node0.NodeList) == 0 {
					alp.Node0 = alpindex(i, alpino.Node0)
				}
				sid := alp.Node0.Id
				if alp.Node0.OtherId != "" {
					sid = alp.Node0.OtherId
				}
				if !seen[sid] {
					seen[sid] = true
					lvl1 := make([]bool, len(woorden)+1)
					alpscan(alp.Node0, alpino.Node0, lvl1)
					for j, n := range lvl1 {
						if n {
							lvl[j]++
						}
					}
				}
			}
		}
	}

	var buf bytes.Buffer

	fmt.Fprint(&buf, "<li>")
	l := 0
	for i, woord := range woorden {
		for l < lvl[i] {
			l++
			fmt.Fprintf(&buf, "<span class=\"c%d\">", l)
		}
		fmt.Fprint(&buf, html.EscapeString(woord))
		for l > lvl[i+1] {
			l--
			fmt.Fprint(&buf, "</span>")
		}
		fmt.Fprint(&buf, " ")
	}

	if strings.HasPrefix(dactfile, dir.Data+"/") {
		dactfile = strings.Replace(dactfile, dir.Data, "$$", 1)
	}

	fmt.Fprintf(&buf, "\n<a href=\"tree?db=%s&amp;names=true&amp;mwu=false&amp;arch=%s&amp;file=%s&amp;global=%v&amp;marknodes=%s&amp;ud1=%s&amp;ud2=%s\" class=\"ico\" target=\"_blank\">&#10020;</a>\n",
		prefix,
		html.EscapeString(dactfile),
		html.EscapeString(filename),
		global,
		strings.Join(ids, ","),
		strings.Join(ud1, "|"),
		strings.Join(ud2, "|"))

	fmt.Fprintf(q.w, `<script type="text/javascript"><!--
$('ol').append(%q);
//--></script>
`, buf.String())

	if ff, ok := q.w.(http.Flusher); ok {
		ff.Flush()
	}
}

// zet de waarde voor alle woorden onder node op true
func alpscan(node, node0 *pqnode.Node, lvl1 []bool) {
	if node == nil {
		return
	}
	if strings.TrimSpace(node.Word) != "" {
		lvl1[node.Begin] = true
	}
	if idx, err := strconv.Atoi(node.Index); err == nil && strings.TrimSpace(node.Word) == "" && len(node.NodeList) == 0 {
		alpscan(alpindex(idx, node0), node0, lvl1)
	}
	for _, n := range node.NodeList {
		alpscan(n, node0, lvl1)
	}
}

// vind de node met een index
func alpindex(idx int, node *pqnode.Node) *pqnode.Node {
	if i, err := strconv.Atoi(node.Index); err == nil && i == idx && (strings.TrimSpace(node.Word) != "" || len(node.NodeList) > 1) {
		return node
	}
	for _, n := range node.NodeList {
		if n2 := alpindex(idx, n); n2 != nil {
			return n2
		}
	}
	return nil
}

func clearLoading(w http.ResponseWriter) {
	fmt.Fprint(w, `
<script type="text/javascript"><!--
$('#loading').addClass('hide');
//--></script>
`)
}

var reBugtest = regexp.MustCompile(`\[err:[A-Z]+[0-9]+\]`)

func bugtest(filename, xpath string) error {

	if pqbugtest == "" {
		for _, d := range strings.Split(Cfg.Path, string(filepath.ListSeparator)) {
			bt := filepath.Join(d, "pqbugtest")
			fi, err := os.Stat(bt)
			if err != nil {
				continue
			}
			if (fi.Mode() | 0111) != 0 {
				pqbugtest = bt
				break
			}
		}
		if pqbugtest == "" {
			chLog <- "ERROR: Geen path naar pqbugtest"
			pqbugtest = "pqbugtest"
		}
	}

	b, err := exec.Command(pqbugtest, filename, xpath).CombinedOutput()
	if err != nil {
		return err
	}
	s := strings.TrimSpace(strings.Replace(string(b), "\n", " ", -1))
	if s == "OK" {
		return nil
	}
	e := errors.New(s)
	if !reBugtest.MatchString(s) {
		logerr(errors.New("BUGTEST: " + strings.Replace(xpath, "\n", " ", -1)))
		logerr(e)
	}
	return e
}

func findUdId(node *pqnode.Node, ID string) *pqnode.Node {
	if node.Ud != nil && node.Ud.Id == ID {
		return node
	}
	for _, n := range node.NodeList {
		if ud := findUdId(n, ID); ud != nil {
			return ud
		}
	}
	return nil
}

func findDepId(node *pqnode.Node, ID string, head string, deprel string) *pqnode.Node {
	if node.Ud != nil && node.Ud.Dep != nil {
		for _, d := range node.Ud.Dep {
			if d.Id == ID && d.Head == head && d.Deprel == deprel {
				if len(node.Ud.Dep) == 1 {
					return node
				}
				n := *node // kopie
				n.NodeList = nil
				ud := *node.Ud // kopie
				ud.Dep = []pqnode.DepType{d}
				n.Ud = &ud
				return &n
			}
		}
	}
	for _, n := range node.NodeList {
		if dep := findDepId(n, ID, head, deprel); dep != nil {
			return dep
		}
	}
	return nil
}

func xpath_do_search(q *Context, query string, prefix string, methode string, offset int, xpathmax int, chClose <-chan bool, doHtml bool, step int, gettextcount bool) (curno int, hash string, methode2 string, loading bool, errval error) {

	methode2 = methode

	var owner string
	var nlines uint64
	var rows *sql.Rows
	rows, errval = sqlDB.Query(fmt.Sprintf("SELECT `owner`,`nline` FROM `%s_info` WHERE `id` = %q", Cfg.Prefix, prefix))
	if logerr(errval) {
		return
	}
	for rows.Next() {
		errval = rows.Scan(&owner, &nlines)
		if logerr(errval) {
			rows.Close()
			return
		}
	}
	errval = rows.Err()
	if logerr(errval) {
		return
	}

	dactfiles := make([]string, 0)
	global := false
	if strings.Contains(owner, "@") {
		dactfiles = append(dactfiles, filepath.Join(dir.Data, "data", prefix, "data.dact"))
	} else {
		global = true
		rows, errval = sqlDB.Query(fmt.Sprintf("SELECT `arch` FROM `%s_c_%s_arch` ORDER BY `id`", Cfg.Prefix, prefix))
		if logerr(errval) {
			return
		}
		for rows.Next() {
			var s string
			errval = rows.Scan(&s)
			if logerr(errval) {
				rows.Close()
				return
			}
			if strings.HasSuffix(s, ".dact") {
				dactfiles = append(dactfiles, s)
			}
		}
		errval = rows.Err()
		if logerr(errval) {
			return
		}
	}

	if len(dactfiles) == 0 {
		fmt.Fprintln(q.w, "Er zijn geen dact-bestanden voor dit corpus")
		return
	}

	if doHtml {
		fmt.Fprintf(q.w, "<ol start=\"%d\" id=\"ol\" class=\"xpath\">\n</ol>\n", offset+1)

		fmt.Fprintln(q.w, "<div id=\"loading\"><img src=\"busy.gif\" alt=\"[bezig]\"> <span></span></div>")
		if ff, ok := q.w.(http.Flusher); ok {
			ff.Flush()
		}
		loading = true

		if !q.hasud[prefix] {
			fmt.Fprintln(q.w, `<script type="text/javascript"><!--
$('#ol').before('<div class="warning">Zoeken naar <b>Universal Dependencies</b> is <b>niet beschikbaar</b> voor dit corpus.</div>');
//--></script>`)
			if ff, ok := q.w.(http.Flusher); ok {
				ff.Flush()
			}
		}

	}

	found := false
	curno = 0
	filename := ""
	curdac := ""
	xmlall := ""
	xmlparts := make([]string, 0)
	fullquery := query
	if strings.Contains(query, "%") {
		rules := getMacrosRules(q)
		fullquery = macroKY.ReplaceAllStringFunc(query, func(s string) string {
			return rules[s[1:len(s)-1]]
		})
		hash = fmt.Sprintf("%x", md5.Sum([]byte(fullquery)))
	}

	queryparts := strings.Split(fullquery, "+|+")

	textcount := 0

	var seen uint64
	for i, dactfile := range dactfiles {
		select {
		case <-chClose:
			logerr(errConnectionClosed)
			return
		default:
		}

		dactname := ""
		if len(dactfiles) > 1 {
			dactname = path.Base(dactfile) + "//"
		}

		if Cfg.Dactx && methode == "dx" {
			if i == 0 {
				if _, err := os.Stat(dactfile + "x"); err != nil {
					methode = "std"
					methode2 = "std"
					if doHtml {
						fmt.Fprintln(q.w, `<script type="text/javascript"><!--
$('#ol').before('<div class="warning">Geen ge&euml;xpandeerde indexnodes beschikbaar voor dit corpus.<br>De standaardmethode wordt gebruikt.</div>');
//--></script>`)
					}
				}
			}
			if methode == "dx" {
				dactfile += "x"
			}
		}

		if curno > offset+xpathmax {
			break
		}

		if doHtml && seen > 0 {
			fmt.Fprintf(q.w, `<script type="text/javascript"><!--
$('#loading span').html('%.1f%%');
//--></script>
`, float64(seen)*100/float64(nlines))
			if ff, ok := q.w.(http.Flusher); ok {
				ff.Flush()
			}
		}

		errval = bugtest(dactfile, queryparts[0])
		if errval != nil {
			return
		}

		var db *dbxml.Db
		db, errval = dbxml.OpenRead(dactfile)
		if logerr(errval) {
			return
		}

		var qu *dbxml.Query
		qu, errval = db.Prepare(queryparts[0])
		if logerr(errval) {
			return
		}
		done := make(chan bool, 1)
		interrupted := make(chan bool, 1)
		go func() {
			select {
			case <-chClose:
				interrupted <- true
				logerr(errConnectionClosed)
				qu.Cancel()
			case <-done:
			}
		}()

		var docs *dbxml.Docs
		docs, errval = qu.Run()
		if logerr(errval) {
			done <- true
			return
		}
		filename = ""
		prevresult := ""
		var variants map[string]bool
	NEXTDOC:
		for docs.Next() {
			name := docs.Name()
			newdoc := false
			if name != filename {
				if found && curno > offset && curno <= offset+xpathmax {
					found = false
					if doHtml {
						xpath_result(q, curno, curdac, filename, xmlall, xmlparts, prefix, global)
					}
					xmlparts = xmlparts[0:0]
				}
				if len(queryparts) == 1 {
					curno++
					if curno > offset+xpathmax {
						docs.Close()
						continue
					}
				}
				curdac = dactfile
				filename = name
				newdoc = true
				if gettextcount {
					if prevresult != "" {
						fmt.Fprintf(q.w, prevresult, len(variants))
						if ff, ok := q.w.(http.Flusher); ok {
							ff.Flush()
						}
						prevresult = ""
					}
					variants = make(map[string]bool)
				}
			}
			if len(queryparts) == 1 {
				if doHtml {
					found = true
					if curno > offset+xpathmax {
						docs.Close()
					} else {
						if curno > offset && curno <= offset+xpathmax {
							xmlall = docs.Content()
							xmlparts = append(xmlparts, docs.Match())
						}
					}
				} else {
					if gettextcount {
						variants[docs.Match()] = true
					}
					if newdoc {
						textcount++
						if (textcount % step) == 0 {
							var alpino Alpino_ds
							err := xml.Unmarshal([]byte(docs.Content()), &alpino)
							name := docs.Name()
							if strings.HasSuffix(name, ".xml") {
								name = name[:len(name)-4]
							}
							if err == nil {
								if gettextcount {
									prevresult = fmt.Sprintf("%s%s\t%%d\t%s\n", dactname, name, alpino.Sentence)
								} else {
									fmt.Fprintf(q.w, "%s%s|%s\n", dactname, name, alpino.Sentence)
								}
							} else {
								if gettextcount {
									prevresult = fmt.Sprintf("%s%s\t%%d\t%v\n", dactname, name, err)
								} else {
									fmt.Fprintf(q.w, "%s%s|%v\n", dactname, name, err)
								}
							}
							if !gettextcount {
								if ff, ok := q.w.(http.Flusher); ok {
									ff.Flush()
								}
							}
						}
					}
				}
			} else if newdoc {
				newdoc = false
				doctxt := fmt.Sprintf("[dbxml:metadata('dbxml:name')=%q]", name)
				var docs2 *dbxml.Docs
				for i := 1; i < len(queryparts)-1; i++ {
					docs2, errval = db.Query(doctxt + queryparts[i])
					if logerr(errval) {
						done <- true
						return
					}
					if !docs2.Next() {
						continue NEXTDOC
					}
					docs2.Close()
				}
				docs2, errval = db.Query(doctxt + queryparts[len(queryparts)-1])
				if logerr(errval) {
					done <- true
					return
				}
				found = false
				for docs2.Next() {
					if doHtml {
						if !found {
							found = true
							curno++
							if curno > offset+xpathmax {
								docs.Close()
							}
						}
						if curno > offset && curno <= offset+xpathmax {
							xmlall = docs2.Content()
							xmlparts = append(xmlparts, docs2.Match())
						} else {
							docs2.Close()
						}
					} else {
						if gettextcount {
							variants[docs2.Match()] = true
						}
						textcount++
						if (textcount % step) == 0 {
							var alpino Alpino_ds
							err := xml.Unmarshal([]byte(docs2.Content()), &alpino)
							name := docs2.Name()
							if strings.HasSuffix(name, ".xml") {
								name = name[:len(name)-4]
							}
							if err == nil {
								if gettextcount {
									prevresult = fmt.Sprintf("%s%s\t%%d\t%s\n", dactname, name, alpino.Sentence)
								} else {
									fmt.Fprintf(q.w, "%s%s|%s\n", dactname, name, alpino.Sentence)
								}
							} else {
								if gettextcount {
									prevresult = fmt.Sprintf("%s%s\t%%d\t%v\n", dactname, name, err)
								} else {
									fmt.Fprintf(q.w, "%s%s|%v\n", dactname, name, err)
								}
							}
							if !gettextcount {
								if ff, ok := q.w.(http.Flusher); ok {
									ff.Flush()
								}
							}
						}
						docs2.Close()
					}
				}
			}
		} // for docs.Next()

		if gettextcount {
			if prevresult != "" {
				fmt.Fprintf(q.w, prevresult, len(variants))
				if ff, ok := q.w.(http.Flusher); ok {
					ff.Flush()
				}
			}
		}

		errval = docs.Error()
		docs = nil
		if logerr(errval) {
			done <- true
			return
		}

		if len(dactfiles) > 1 {
			if n, err := db.Size(); err == nil {
				seen += n
			}
		}
		db.Close()
		db = nil
		done <- true
		select {
		case <-interrupted:
			return
		default:
		}

		if found && curno > offset && curno <= offset+xpathmax {
			found = false
			if doHtml {
				xpath_result(q, curno, curdac, filename, xmlall, xmlparts, prefix, global)
			}
			xmlparts = xmlparts[0:0]
		}
		if curno > offset+xpathmax {
			break
		}
	} // for _, dactfile := range dactfiles

	if doHtml {
		clearLoading(q.w)
		loading = false
	}

	return
}
