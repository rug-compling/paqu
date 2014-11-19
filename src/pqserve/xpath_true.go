// +build !nodbxml

package main

import (
	"github.com/pebbe/dbxml"

	"bytes"
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	reXpath    = regexp.MustCompile(`'[^']*'|"[^"]*"|@[a-z][-_a-z]*|\$[a-z][-_a-zA-Z0-9]*|[a-zA-Z][-_a-zA-Z]*:*(\s*\()?`)
	keyTags    = make(map[string]bool)
	xpathNames = map[string]bool{
		"ancestor-or-self::":   true,
		"ancestor::":           true,
		"attribute::":          true,
		"child::":              true,
		"descendant-or-self::": true,
		"descendant::":         true,
		"following-sibling::":  true,
		"following::":          true,
		"namespace::":          true,
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
		"mod":                  true,
		"ne":                   true,
		"node":                 true,
		"of":                   true,
		"only":                 true,
		"or":                   true,
		"precedes":             true,
		"return":               true,
		"satisfies":            true,
		"some":                 true,
		"then":                 true,
		"to":                   true,
		"treat":                true,
		"type":                 true,
		"union":                true,
	}
)

func init() {
	for _, tag := range NodeTags {
		keyTags[tag] = true
	}
}

func xpathcheck(q *Context) {
	contentType(q, "text/plain")

	query := first(q.r, "xpath")
	if query == "" {
		cache(q)
		fmt.Fprintln(q.w, "0")
		return
	}

	ch := true
	if strings.Contains(query, "%") {
		rules := getMacrosRules(q)
		query = macroKY.ReplaceAllStringFunc(query, func(s string) string {
			if ch {
				nocache(q)
				ch = false
			}
			return rules[s[1:len(s)-1]]
		})
	}
	if ch {
		cache(q)
	}

	lvl := 0
PARTLOOP:
	for _, part := range strings.Split(query, "+|+") {

		if strings.TrimSpace(part) == "" {
			fmt.Fprintln(q.w, "2")
			return
		}

		// syntactisch fout -> 2
		if part == "." || part == "/" || dbxml.Check(part) != nil {
			fmt.Fprintln(q.w, "2")
			return
		}

		// geen resultaat -> 1
		for i, s := range reXpath.FindAllString(part, -1) {
			if i == 0 && s == "alpino_ds" {
				continue
			}
			if s[0] == '\'' || s[0] == '"' || s[0] == '$' {
				continue
			}
			if s[0] == '@' {
				if keyTags[s[1:]] {
					continue
				}
				lvl = 1
				continue PARTLOOP
			}
			if strings.HasSuffix(s, "(") {
				continue
			}
			if !xpathNames[s] {
				lvl = 1
				continue PARTLOOP
			}
		}

	}

	fmt.Fprintln(q.w, lvl)
}

// TAB: xpath
func xpath(q *Context) {

	prefix := getprefix(q)
	if !q.prefixes[prefix] {
		http.Error(q.w, "Invalid corpus: "+prefix, http.StatusPreconditionFailed)
		return
	}

	// HTML-uitvoer van begin van de pagina
	writeHead(q, "", 2)
	html_xpath_header(q)

	// HTML-uitvoer van het formulier
	// Returnwaarde is true als er een query was gedefinieerd
	has_query := html_xpath_form(q)

	// Als er geen query is gedefinieerd, HTML-uitvoer van korte helptekst, pagina-einde, en exit
	if !has_query {
		html_xpath_uitleg(q)
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
	if offset < 0 {
		offset = 0
	}

	fmt.Fprintln(q.w, "<hr>")

	now := time.Now()

	var owner string
	var nlines uint64
	rows, err := q.db.Query(fmt.Sprintf("SELECT `owner`,`nline` FROM `%s_info` WHERE `id` = %q", Cfg.Prefix, prefix))
	if doErr(q, err) {
		return
	}
	for rows.Next() {
		if doErr(q, rows.Scan(&owner, &nlines)) {
			rows.Close()
			return
		}
	}
	if doErr(q, rows.Err()) {
		return
	}

	dactfiles := make([]string, 0)
	global := false
	if strings.Contains(owner, "@") {
		dactfiles = append(dactfiles, path.Join(paqudir, "data", prefix, "data.dact"))
	} else {
		global = true
		rows, err := q.db.Query(fmt.Sprintf("SELECT `arch` FROM `%s_c_%s_arch` ORDER BY `id`", Cfg.Prefix, prefix))
		if doErr(q, err) {
			return
		}
		for rows.Next() {
			var s string
			if doErr(q, rows.Scan(&s)) {
				rows.Close()
				return
			}
			if strings.HasSuffix(s, ".dact") {
				dactfiles = append(dactfiles, s)
			}
		}
		if doErr(q, rows.Err()) {
			return
		}
	}

	if len(dactfiles) == 0 {
		fmt.Fprintln(q.w, "Er zijn geen dact-bestanden voor dit corpus")
		return
	}

	fmt.Fprintf(q.w, "<ol start=\"%d\" id=\"ol\" class=\"xpath\">\n</ol>\n", offset+1)

	fmt.Fprintln(q.w, "<div id=\"loading\"><img src=\"busy.gif\"> <span></span></div>")
	if ff, ok := q.w.(http.Flusher); ok {
		ff.Flush()
	}

	found := false
	curno := 0
	filename := ""
	curdac := ""
	xmlall := ""
	xmlparts := make([]string, 0)
	query := first(q.r, "xpath")
	fullquery := query
	if strings.Contains(query, "%") {
		rules := getMacrosRules(q)
		fullquery = macroKY.ReplaceAllStringFunc(query, func(s string) string {
			return rules[s[1:len(s)-1]]
		})
	}

	queryparts := strings.Split(fullquery, "+|+")

	var seen uint64
	for _, dactfile := range dactfiles {
		select {
		case <-chClose:
			logerr(errConnectionClosed)
			return
		default:
		}

		if seen > 0 {
			fmt.Fprintf(q.w, `<script type="text/javascript"><!--
$('#loading span').html('%.1f%%');
//--></script>
`, float64(seen)*100/float64(nlines))
			if ff, ok := q.w.(http.Flusher); ok {
				ff.Flush()
			}
		}

		db, err := dbxml.Open(dactfile)
		if doErr(q, err) {
			return
		}

		qu, err := db.Prepare(queryparts[0])
		if err != nil {
			fmt.Fprintln(q.w, html.EscapeString(err.Error()))
			db.Close()
			clearLoading(q.w)
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

		docs, err := qu.Run()
		if err != nil {
			fmt.Fprintln(q.w, html.EscapeString(err.Error()))
			db.Close()
			done <- true
			clearLoading(q.w)
			return
		}
		filename = ""
	NEXTDOC:
		for docs.Next() {
			name := docs.Name()
			newdoc := false
			if name != filename {
				if found && curno > offset && curno <= offset+ZINMAX*2 {
					found = false
					xpath_result(q, curno, curdac, filename, xmlall, xmlparts, prefix, global)
					xmlparts = xmlparts[0:0]
				}
				if len(queryparts) == 1 {
					curno++
				}
				curdac = dactfile
				filename = name
				newdoc = true
			}
			if len(queryparts) == 1 {
				found = true
				if curno > offset+ZINMAX*2 {
					docs.Close()
				} else {
					if curno > offset && curno <= offset+ZINMAX*2 {
						xmlall = docs.Content()
						xmlparts = append(xmlparts, docs.Match())
					}
				}
			} else if newdoc {
				newdoc = false
				doctxt := fmt.Sprintf("[dbxml:metadata('dbxml:name')=%q]", name)
				for i := 1; i < len(queryparts)-1; i++ {
					docs2, err := db.Query(doctxt + queryparts[i])
					if err != nil {
						fmt.Fprintln(q.w, html.EscapeString(err.Error()))
						docs.Close()
						db.Close()
						done <- true
						clearLoading(q.w)
						return
					}
					if !docs2.Next() {
						continue NEXTDOC
					}
					docs2.Close()
				}
				docs2, err := db.Query(doctxt + queryparts[len(queryparts)-1])
				if err != nil {
					fmt.Fprintln(q.w, html.EscapeString(err.Error()))
					docs.Close()
					db.Close()
					done <- true
					clearLoading(q.w)
					return
				}
				found = false
				for docs2.Next() {
					if !found {
						found = true
						curno++
						if curno > offset+ZINMAX*2 {
							docs.Close()
						}
					}
					if curno > offset && curno <= offset+ZINMAX*2 {
						xmlall = docs2.Content()
						xmlparts = append(xmlparts, docs2.Match())
					} else {
						docs2.Close()
					}
				}
			}
		}
		if err := docs.Error(); err != nil {
			logerr(err)
		}
		if n, err := db.Size(); err == nil {
			seen += n
		}
		db.Close()
		done <- true
		select {
		case <-interrupted:
			return
		default:
		}

		if found && curno > offset && curno <= offset+ZINMAX*2 {
			found = false
			xpath_result(q, curno, curdac, filename, xmlall, xmlparts, prefix, global)
			xmlparts = xmlparts[0:0]
		}
		if curno > offset+ZINMAX*2 {
			break
		}
	}

	clearLoading(q.w)

	if curno == 0 {
		fmt.Fprintf(q.w, "geen match gevonden")
	}

	// Links naar volgende en vorige pagina's met resultaten
	qs := "xpath=" + urlencode(query)
	if offset > 0 || curno > offset+ZINMAX*2 {
		if offset > 0 {
			fmt.Fprintf(q.w, "<a href=\"/xpath?%s&amp;offset=%d\">vorige</a>", qs, offset-ZINMAX*2)
		} else {
			fmt.Fprint(q.w, "vorige")
		}
		fmt.Fprint(q.w, " | ")
		if curno > offset+ZINMAX*2 {
			fmt.Fprintf(q.w, "<a href=\"/xpath?%s&amp;offset=%d\">volgende</a>", qs, offset+ZINMAX*2)
		} else {
			fmt.Fprint(q.w, "volgende")
		}
	}

	fmt.Fprintln(q.w, "<hr><small>tijd:", tijd(time.Now().Sub(now)), "</small><hr>")

	if curno == 0 {
		html_footer(q)
		return
	}

	// Links naar statistieken
	fmt.Fprintf(q.w, `<p>
		<div id="xstats">
		<form action="xpathstats" target="xframe" name="xstatsform" onsubmit="javascript:return xstatftest()">
		<input type="hidden" name="xpath" value="%s">
		<input type="hidden" name="db" value="%s">
		Selecteer &eacute;&eacute;n tot drie attributen:<br>
`, html.EscapeString(query), html.EscapeString(prefix))

	for i := 1; i <= 3; i++ {
		fmt.Fprintf(q.w, "<select name=\"attr%d\">\n<option value=\"\">--</option>\n", i)
		for _, s := range NodeTags {
			fmt.Fprintf(q.w, "<option>%s</option>\n", s)
		}
		fmt.Fprintln(q.w, "</select>")
	}

	fmt.Fprint(q.w, `
		<p>
		<input type="submit" value="tellingen">
		</form>
		<p>
        <iframe src="leeg.html" name="xframe" class="hide"></iframe>
        <div id="result"></div>
		</div>
`)

	html_footer(q)

}

//. HTML

func html_xpath_header(q *Context) {
	fmt.Fprint(q.w, `
<script type="text/javascript" src="jquery.js"></script>
<script type="text/javascript"><!--

  function formclear(f) {
    f.xpath.value = "";
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

  var result;
  var at1, at2, at3;
  var xquery;
  window._fn = {
    update: function(data) {
      result.html(data);
    },
    update2: function(data) {
      if (data.err == "") {
          $('#macromsg').addClass('hide');
          $('#macrotext').text(data.macros);
          $('#macrotext').val(data.macros);
          macros = data.keys;
      } else {
          $('#macromsg').removeClass('hide').text("Fout: " + data.err);
      }
      disableSave();
    }
  }
  function xstatftest() {
    var n = 0;
    if (at1.selectedIndex > 0) { n++; }
    if (at2.selectedIndex > 0) { n++; }
    if (at3.selectedIndex > 0) { n++; }
    if (n < 1) {
      alert("Geen attribuut geselecteerd");
      return false;
    }
    setCookie("xpattr1", at1.selectedIndex, 14);
    setCookie("xpattr2", at2.selectedIndex, 14);
    setCookie("xpattr3", at3.selectedIndex, 14);
    return true;
  }
  function setForm() {
    try {
      at1.selectedIndex = getCookie("xpattr1");
      at2.selectedIndex = getCookie("xpattr2");
      at3.selectedIndex = getCookie("xpattr3");
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
       $('#btExpand').removeClass('hide');
    } else {
       $('#btExpand').addClass('hide');
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

  function macroExpand() {
    $.ajax("macroexpand?" + xquery.serialize())
      .done(function(data) {
         $("#macroInner").text(data);
      }).fail(function(e) {
         $("#macroInner").text(e.responseText);
      })
      $('#macroOuter').removeClass('hide');
  }

  function openMacro() {
      $('#openmacro').addClass('hide');
      $('#macros').removeClass('hide');
  }

  function sluitMacro() {
      $('#openmacro').removeClass('hide');
      $('#macros').addClass('hide');
      $('#macromsg').addClass('hide');
  }

  function enableDeleteSave(e) {
      var k = e.which;
      if (k == 8 || k == 46) {
        enableSave();
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

  function init() {
    xquery = $('#xquery');
    xquery.on('keyup', qcheck);
    xquery.on('change', qcheck);
    xquery.on('click', qcheck);
    qcheckdo();
    $('#openmacro').on('click', openMacro);
    $('#sluitmacro').on('click', sluitMacro);
    $('#macroreset').on('click', disableSave);
    $('#macrotext').on('keypress', enableSave);
    $('#macrotext').on('keyup', enableDeleteSave);
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
    $('#btExpand').on('click', macroExpand);
    $('#btClose').on('click', function() { $('#macroOuter').addClass('hide'); });
  }

  $(document).ready(function() {
    result = $('#result');
    try {
      var f = document.forms["xstatsform"];
      at1 = f["attr1"];
      at2 = f["attr2"];
      at3 = f["attr3"];
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
<a href="http://rug-compling.github.io/dact/cookbook/" target="_blanc">Dact Cookbook</a>
`)
}

func html_xpath_form(q *Context) (has_query bool) {
	has_query = true
	if first(q.r, "xpath") == "" {
		has_query = false
	}

	if q.auth {
		macros := ""
		rows, err := q.db.Query(fmt.Sprintf("SELECT `macros` FROM `%s_macros` WHERE `user` = %q", Cfg.Prefix, q.user))
		if err == nil {
			if rows.Next() {
				rows.Scan(&macros)
				rows.Close()
			}
		}
		fmt.Fprintf(q.w, `
<button id="openmacro">Macro's</button>
<div id="macros" class="hide">
<button id="sluitmacro">Macro's sluiten</button>
<div class="hide warning" id="macromsg">
</div>
<p>
<form action="savemacros" method="post" target="hiddenframe" enctype="multipart/form-data">
<input type="submit" value="Uploaden" id="macrofilesave" disabled="disabled">
<input type="file" id="macrofilename" value="Bestand kiezen" name="macrotext" id="macrofilename">
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
corpus: <select name="db">
`)
	html_opts(q, q.opt_db, getprefix(q), "corpus")
	fmt.Fprintln(q.w, "</select>")
	if q.auth {
		fmt.Fprintln(q.w, "<a href=\"corpuslijst\">meer/minder</a>")
	}
	fmt.Fprintf(q.w, `<p>
		XPATH query (<a href="http://rug-compling.github.io/dact/cookbook/" target="_blanc">voorbeelden</a>):<br>
		<textarea name="xpath" rows="6" cols="80" maxlength="1200" id="xquery">%s</textarea>
		`, html.EscapeString(first(q.r, "xpath")))
	fmt.Fprint(q.w, `<p>
           <input type="submit" value="Zoeken">
           <input type="button" value="Wissen" onClick="javascript:formclear(form)">
           <input type="reset" value="Reset" onClick="javascript:qcheck()">
           <input type="button" id="btExpand" class="hide" value="Toon macro-expansie">
       </form>
       <div id="macroOuter" class="hide">
       <div id="macroInner"></div>
       <button id="btClose">Sluiten</button>
       </div>
<script type="text/javascript" src="jquery.textcomplete.js"></script>
<script type="text/javascript"><!--
var begin = ['//node', '/alpino_ds/node'];
var other = ['/node',
        "@aform",
        "@begin",
        "@buiging",
        "@case",
        "@cat",
        "@comparative",
        "@conjtype",
        "@def",
        "@dial",
        "@end",
        "@frame",
        "@gen",
        "@genus",
        "@getal",
        "@getal-n",
        "@graad",
        "@id",
        "@index",
        "@infl",
        "@lcat",
        "@lemma",
        "@lwtype",
        "@mwu_root",
        "@mwu_sense",
        "@naamval",
        "@neclass",
        "@npagr",
        "@ntype",
        "@num",
        "@numtype",
        "@pb",
        "@pdtype",
        "@per",
        "@persoon",
        "@pos",
        "@positie",
        "@postag",
        "@pt",
        "@pvagr",
        "@pvtijd",
        "@refl",
        "@rel",
        "@root",
        "@sc",
        "@sense",
        "@special",
        "@spectype",
        "@status",
        "@tense",
        "@vform",
        "@vwtype",
        "@vztype",
        "@wh",
        "@wk",
        "@word",
        "@wvorm"];

var axis = [
		"node",
		"and",
		"div",
		"mod",
		"or",
		"fn:",
		"ancestor-or-self::",
		"ancestor::",
		"attribute::",
		"child::",
		"descendant-or-self::",
		"descendant::",
		"following-sibling::",
		"following::",
		"namespace::",
		"parent::",
		"preceding-sibling::",
		"preceding::",
		"self::"];

var macros = [`)

	keys := getMacrosKeys(q)
	p := ","
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
    match: /^(\/\/?[_a-z]*)$/,
    search: function (term, callback) {
        callback($.map(begin, function (e) {
            return e.indexOf(term) === 0 ? e : null;
        }));
    },
    replace: function (value) {
        return value + ' ';
    },
    index: 1
},
{
    match: /([\/@][-_a-z]*)$/,
    search: function (term, callback) {
        callback($.map(other, function (e) {
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
    match: /::([a-z]*)$/,
    search: function (term, callback) {
        callback($.map(['node'], function (e) {
            return e.indexOf(term) === 0 ? e : null;
        }));
    },
    replace: function (value) {
        return '::' + value + ' ';
    },
    index: 1,
    context: outText
},
{
    match: /^((.|\n)*[ \n'")\]])$/,
    search: function (term, callback) {
        var chars = [];
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
        var rt = [];
        var s = "";
        while (chars.length > 0) {
            s += chars.shift();
            rt.push(s);
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
    maxCount: 100,
    debounce: 100,
});

init();

//--></script>
`)

	return
}

func xpath_result(q *Context, curno int, dactfile, filename, xmlall string, xmlparts []string, prefix string, global bool) {
	alpino := Alpino_ds{}
	err := xml.Unmarshal([]byte(xmlall), &alpino)
	if err != nil {
		fmt.Fprintf(q.w, "FOUT bij parsen van XML: %s\n", html.EscapeString(err.Error()))
		return
	}
	woorden := strings.Fields(alpino.Sentence)

	lvl := make([]int, len(woorden)+1)
	ids := make([]string, len(xmlparts))

	for i, part := range xmlparts {
		alp := Alpino_ds{}
		err := xml.Unmarshal([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<alpino_ds version="1.3">
`+part+`
</alpino_ds>`), &alp)
		if err != nil {
			fmt.Fprintf(q.w, "FOUT bij parsen van XML: %s\n", html.EscapeString(err.Error()))
			return
		}
		if alp.Node0 != nil {
			ids[i] = alp.Node0.Id
			lvl1 := make([]int, len(woorden)+1)
			alpscan(alp.Node0, alpino.Node0, lvl1)
			for j, n := range lvl1 {
				lvl[j] += n
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
		fmt.Fprintf(&buf, html.EscapeString(woord))
		for l > lvl[i+1] {
			l--
			fmt.Fprint(&buf, "</span>")
		}
		fmt.Fprint(&buf, " ")
	}

	fmt.Fprintf(&buf, "\n<a href=\"/tree?db=%s&amp;names=true&amp;mwu=false&amp;arch=%s&amp;file=%s&amp;global=%v&amp;marknodes=%s\" class=\"ico\">&#10020;</a>\n",
		prefix,
		html.EscapeString(dactfile),
		html.EscapeString(filename),
		global,
		strings.Join(ids, ","))

	fmt.Fprintf(q.w, `<script type="text/javascript"><!--
$('ol').append(%q);
//--></script>
`, buf.String())

	if ff, ok := q.w.(http.Flusher); ok {
		ff.Flush()
	}
}

// zet de teller voor alle woorden onder node 1 hoger
func alpscan(node, node0 *Node, lvl1 []int) {
	if node == nil {
		return
	}
	if strings.TrimSpace(node.Word) != "" {
		lvl1[node.Begin] = 1
	}
	if idx, err := strconv.Atoi(node.Index); err == nil && strings.TrimSpace(node.Word) == "" && len(node.NodeList) == 0 {
		alpscan(alpindex(idx, node0), node0, lvl1)
	}
	for _, n := range node.NodeList {
		alpscan(n, node0, lvl1)
	}
}

// vind de node met een index
func alpindex(idx int, node *Node) *Node {
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
