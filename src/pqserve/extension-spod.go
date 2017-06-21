// +build extension

package main

import (
	"github.com/pebbe/dbxml"

	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"
)

const (
	SPOD_MAX = 8 // maximaal aantal processen
	SPOD_STD = "std"
	SPOD_DX  = "dx"
)

const (
	SPOD_X_NONE = iota
	SPOD_X_PARSER
	SPOD_X_PARSER_CATS
	SPOD_X_PARSER_SKIPS
	SPOD_X_HIS
)

type Spod struct {
	header  string
	xpath   string
	method  string
	lbl     string
	text    string
	special int
}

type spod_writer struct {
	header map[string][]string
	buffer bytes.Buffer
}

var (
	spods = []Spod{
		{
			"Hoofdzinnen",
			`//node[@cat="smain"]`,
			SPOD_STD,
			"smain",
			"mededelende hoofdzinnen",
			SPOD_X_NONE,
		},
		{
			"",
			`//node[@cat="whq"]`,
			SPOD_STD,
			"whq",
			"vraagzinnen (wh)",
			SPOD_X_NONE,
		},
		{
			"",
			`%PQ_janee_vragen%`,
			SPOD_STD,
			"janee",
			"ja/nee vragen",
			SPOD_X_NONE,
		},
		{
			"",
			`%PQ_imperatieven%`,
			SPOD_STD,
			"imp",
			"imperatieven",
			SPOD_X_NONE,
		},
		{
			"Bijzinnen",
			`%PQ_ingebedde_vraagzinnen%`,
			SPOD_STD,
			"whsub",
			"ingebedde vraagzinnen",
			SPOD_X_NONE,
		},
		{
			"",
			`//node[@cat='cp' and node[@rel='body' and @cat='ssub']]`,
			SPOD_STD,
			"ssub",
			"finiete bijzinnen",
			SPOD_X_NONE,
		},
		{
			"",
			`//node[@cat='cp' and node[@rel='cmp' and @lemma='dat'] and node[@rel='body' and @cat='ssub']]`,
			SPOD_STD,
			"ssubdat",
			"finiete bijzinnen met \"dat\"",
			SPOD_X_NONE,
		},
		{
			"",
			`//node[@cat='cp' and node[@rel='cmp' and @lemma='of'] and node[@rel='body' and @cat='ssub']]`,
			SPOD_STD,
			"ssubof",
			"finiete bijzinnen met \"of\"",
			SPOD_X_NONE,
		},
		{
			"",
			`//node[@cat='cp' and node[@rel='cmp' and not(@lemma='of' or @lemma='dat')] and node[@rel='body' and @cat='ssub']]`,
			SPOD_STD,
			"ssubcmp",
			"finiete bijzinnen met andere voegwoorden",
			SPOD_X_NONE,
		},
		{
			"",
			`//node[@cat='oti']`,
			SPOD_STD,
			"oti",
			"infiniete bijzinnen met \"om\"",
			SPOD_X_NONE,
		},
		{
			"",
			`//node[@cat='ti' and not(../@cat='oti' or ../@cat='cp')]`,
			SPOD_STD,
			"tite",
			"infiniete bijzinnen met alleen \"te\"",
			SPOD_X_NONE,
		},
		{
			"",
			`//node[@cat='cp' and node[@rel='body' and @cat='ti']]`,
			SPOD_STD,
			"ti",
			"infiniete bijzinnen met ander voorzetsel",
			SPOD_X_NONE,
		},
		{
			"",
			`%PQ_relatieve_bijzinnen%`,
			SPOD_STD,
			"relssub",
			"relatieve bijzinnen",
			SPOD_X_NONE,
		},
		{
			"",
			`%PQ_free_relatives%`,
			SPOD_STD,
			"whrel",
			"free relatives",
			SPOD_X_NONE,
		},
		{
			"Woordgroepen",
			`//node[%PQ_np%]`,
			SPOD_STD,
			"np",
			"np",
			SPOD_X_NONE,
		},
		{
			"",
			`//node[@cat='pp']`,
			SPOD_STD,
			"pp",
			"pp",
			SPOD_X_NONE,
		},
		{
			"",
			`//node[@cat='ap' or @pt='adj' and not(@rel='hd')]`,
			SPOD_STD,
			"ap",
			"ap",
			SPOD_X_NONE,
		},
		{
			"",
			`//node[@cat='advp' or @pt='bw' and not(@rel='hd')]`,
			SPOD_STD,
			"advp",
			"advp",
			SPOD_X_NONE,
		},
		{
			"Werkwoorden",
			`//node[@rel="svp" and @cat]`,
			SPOD_STD,
			"vwuit",
			"vaste werkwoordelijke uitdrukkingen",
			SPOD_X_NONE,
		},
		{
			"",
			`%PQ_groen%`,
			SPOD_STD,
			"groen",
			"groene werkwoordsvolgorde",
			SPOD_X_NONE,
		},
		{
			"",
			`%PQ_rood%`,
			SPOD_STD,
			"rood",
			"rode werkwoordsvolgorde",
			SPOD_X_NONE,
		},
		{
			"",
			`//node[%PQ_dep_node_in_verbcluster%]`,
			SPOD_STD,
			"wwclus",
			"werkwoordsclusters",
			SPOD_X_NONE,
		},
		{
			"",
			`%PQ_cross_serial_verbcluster%`,
			SPOD_STD,
			"accinf",
			"accusativus cum infinitivo -- werkt niet met Corpus Gesproken Nederlands",
			SPOD_X_NONE,
		},
		{
			"",
			`//node[%PQ_passive%]`,
			SPOD_STD,
			"passive",
			"passief -- werkt niet met Corpus Gesproken Nederlands",
			SPOD_X_NONE,
		},
		{
			"",
			`//node[%PQ_impersonal_passive%]`,
			SPOD_STD,
			"nppas",
			"niet-persoonlijke passief -- werkt niet met Corpus Gesproken Nederlands",
			SPOD_X_NONE,
		},
		{
			"Inbedding",
			`//node[%PQ_finiete_inbedding0%]`,
			SPOD_STD,
			"inb0",
			"geen inbedding",
			SPOD_X_NONE,
		},
		{
			"",
			`//node[%PQ_finiete_inbedding1%]`,
			SPOD_STD,
			"inb1",
			"minstens 1 finiete zinsinbedding",
			SPOD_X_NONE,
		},
		{
			"",
			`//node[%PQ_finiete_inbedding2%]`,
			SPOD_STD,
			"inb2",
			"minstens 2 finiete zinsinbeddingen",
			SPOD_X_NONE,
		},
		{
			"",
			`//node[%PQ_finiete_inbedding3%]`,
			SPOD_STD,
			"inb3",
			"minstens 3 finiete zinsinbeddingen",
			SPOD_X_NONE,
		},
		{
			"",
			`//node[%PQ_finiete_inbedding4%]`,
			SPOD_STD,
			"inb4",
			"minstens 4 finiete zinsinbeddingen",
			SPOD_X_NONE,
		},
		{
			"",
			`//node[%PQ_finiete_inbedding5%]`,
			SPOD_STD,
			"inb5",
			"minstens 5 finiete zinsinbeddingen",
			SPOD_X_NONE,
		},
		{
			"",
			`//node[%PQ_finiete_inbedding6%]`,
			SPOD_STD,
			"inb6",
			"minstens 6 finiete zinsinbeddingen",
			SPOD_X_NONE,
		},
		{
			"",
			`//node[%PQ_finiete_inbedding7%]`,
			SPOD_STD,
			"inb7",
			"minstens 7 finiete zinsinbeddingen",
			SPOD_X_NONE,
		},
		{
			"",
			`//node[%PQ_finiete_inbedding8%]`,
			SPOD_STD,
			"inb8",
			"minstens 8 finiete zinsinbeddingen",
			SPOD_X_NONE,
		},
		{
			"Topicalizatie en Extractie",
			`%PQ_vorfeld_np_subject%`,
			SPOD_STD,
			"nptsub",
			"np-topic is subject",
			SPOD_X_NONE,
		},
		{
			"",
			`%PQ_vorfeld_np_no_subject%`,
			SPOD_STD,
			"nptnsub",
			"np-topic is geen subject",
			SPOD_X_NONE,
		},
		{
			"",
			`%PQ_vorfeld_non_local%`,
			SPOD_STD,
			"tnonloc",
			"topic is niet lokaal",
			SPOD_X_NONE,
		},
		{
			"",
			`%PQ_local_extraction%`,
			SPOD_STD,
			"locext",
			"lokale extractie",
			SPOD_X_NONE,
		},
		{
			"",
			`%PQ_non_local_extraction%`,
			SPOD_STD,
			"nlocext",
			"niet-lokale extractie",
			SPOD_X_NONE,
		},
		{
			"Parser succes -- werkt niet met Alpino Treebank, Corpus Gesproken Nederlands, Lassy Klein",
			`//parser[@cats="1" and @skips="0"]`,
			SPOD_STD,
			"ok",
			"volledige parse",
			SPOD_X_PARSER,
		},
		{
			"Parser succes: geparste delen",
			`//parser[@cats="0"]`,
			SPOD_STD,
			"cats0",
			"geen enkel deel is geparst",
			SPOD_X_PARSER_CATS,
		},
		{
			"",
			`//parser[@cats="1"]`,
			SPOD_STD,
			"cats1",
			"parse bestaat uit één deel",
			SPOD_X_PARSER_CATS,
		},
		{
			"",
			`//parser[@cats="2"]`,
			SPOD_STD,
			"cats2",
			"parse bestaat uit twee losse delen",
			SPOD_X_PARSER_CATS,
		},
		{
			"",
			`//parser[@cats="3"]`,
			SPOD_STD,
			"cats3",
			"parse bestaat uit drie losse delen",
			SPOD_X_PARSER_CATS,
		},
		{
			"",
			`//parser[number(@cats) > 3]`,
			SPOD_STD,
			"cats4",
			"parse bestaat uit vier of meer losse delen",
			SPOD_X_PARSER_CATS,
		},
		{
			"parser success: overgeslagen woorden",
			`//parser[@skips="0"]`,
			SPOD_STD,
			"skips0",
			"geen enkel woord is overgeslagen",
			SPOD_X_PARSER_SKIPS,
		},
		{
			"",
			`//parser[@skips="1"]`,
			SPOD_STD,
			"skips1",
			"een van de woorden is overgeslagen",
			SPOD_X_PARSER_SKIPS,
		},
		{
			"",
			`//parser[@skips="2"]`,
			SPOD_STD,
			"skips2",
			"twee van de woorden zijn overgeslagen",
			SPOD_X_PARSER_SKIPS,
		},
		{
			"",
			`//parser[@skips="3"]`,
			SPOD_STD,
			"skips3",
			"drie van de woorden zijn overgeslagen",
			SPOD_X_PARSER_SKIPS,
		},
		{
			"",
			`//parser[number(@skips) > 3]`,
			SPOD_STD,
			"skips4",
			"vier of meer van de woorden zijn overgeslagen",
			SPOD_X_PARSER_SKIPS,
		},
		{
			"Onbekende woorden -- werkt niet met Alpino Treebank, Corpus Gesproken Nederlands, Lassy Klein",
			`//node[@his]`,
			SPOD_STD,
			"his",
			"nodes met attribuut @his",
			SPOD_X_HIS,
		},
		{
			"",
			`//node[@his="normal"]`,
			SPOD_STD,
			"normal",
			"woorden uit het woordenboek of de namenlijst",
			SPOD_X_HIS,
		},
		{
			"",
			`//node[@his and not(@his="normal")]`,
			SPOD_STD,
			"onbeken",
			"woorden niet direct uit het woordenboek",
			SPOD_X_HIS,
		},
		{
			"",
			`//node[@his="compound"]`,
			SPOD_STD,
			"compoun",
			"woorden herkend als samenstelling",
			SPOD_X_HIS,
		},
		{
			"",
			`//node[@his="name"]`,
			SPOD_STD,
			"name",
			"woorden herkend als naam (maar niet uit namenlijst)",
			SPOD_X_HIS,
		},
		{
			"",
			`//node[@his and not(@his="normal" or @his="compound" or @his="name")]`,
			SPOD_STD,
			"noun",
			"onbekende woorden die niet als samenstelling of naam werden herkend",
			SPOD_X_HIS,
		},
	}
	spodMu        sync.Mutex
	spodSemaphore = make(chan bool, SPOD_MAX)
	spodWorking   = make(map[string]bool)
	spodRE        = regexp.MustCompile(`([0-9]+)[^0-9]+([0-9]+)`)
	spodREerr     = regexp.MustCompile(`\[err:`)
)

func extension_menu(q *Context, tab int) {
	s := ""
	if tab == 6 {
		s = " class=\"selected\""
	}
	fmt.Fprintln(q.w, "<a href=\"ext?ext=spod\""+s+">SPOD</a>")
}

func extension(q *Context) {

	for i, db := range q.opt_db {
		dbs := strings.Fields(db)
		if dbs[0] == "Alassywiki" || dbs[0] == "Zlassywiki" { // wel/niet ingelogd
			if i == len(q.opt_db)-1 {
				q.opt_db = q.opt_db[:i]
			} else {
				q.opt_db = append(q.opt_db[:i], q.opt_db[i+1:]...)
			}
			break
		}
	}
	delete(q.prefixes, "lassywiki")

	switch first(q.r, "ext") {
	case "spod":
		extension_spod_main(q)
	case "spodform":
		extension_spod_form(q)
	case "spodlist":
		extension_spod_list(q)
	default:
		contentType(q, "text/plain; charset=utf-8")
		nocache(q)
		fmt.Fprint(q.w, `
Ongedefinieerde extensie.
`)
	}
}

func extension_spod_main(q *Context) {

	writeHead(q, "Syntactic profile of Dutch", 6)

	fmt.Fprintln(q.w, "Syntactic profile of Dutch<p>")

	fmt.Fprintln(q.w, "<div class=\"warning\">In ontwikkeling</div>")

	fmt.Fprint(q.w, `
<script type="text/javascript"><!--
  var xpaths = [
`)

	p := ""
	for _, spod := range spods {
		fmt.Fprintf(q.w, "%s['%s', '%s']", p, url.QueryEscape(spod.xpath), spod.method)
		p = ",\n"
	}
	fmt.Fprint(q.w, `
];
function vb(i) {
    var e = document.getElementById("spodform").elements
    window.open("xpath?db=" + e[1].value + "&xpath=" + xpaths[i][0] + "&mt=" + xpaths[i][1]);
}
function alles(n, m) {
    var e = document.getElementById("spodform").elements
    for (var i = n; i < m; i++) {
        e[i+2].checked = true;
    }
}
function niets(n, m) {
    var e = document.getElementById("spodform").elements
    for (var i = n; i < m; i++) {
        e[i+2].checked = false;
    }
}
function omkeren(n, m) {
    var e = document.getElementById("spodform").elements
    for (var i = n; i < m; i++) {
        e[i+2].checked = !e[i+2].checked;
    }
}
//--></script>
<form id="spodform" action="ext" method="get" accept-charset="utf-8" target="_blank">
<input type="hidden" name="ext" value="spodform">
corpus: <select name="db">
`)

	html_opts(q, q.opt_db, getprefix(q), "corpus")
	fmt.Fprintln(q.w, "</select>")
	if q.auth {
		fmt.Fprintln(q.w, "<a href=\"corpuslijst\">meer/minder</a>")
	}

	fmt.Fprintf(q.w, `
<p>
<a href="ext?ext=spodlist" target="_blank">lijst van queries</a>
<p>
<a href="javascript:alles(0, %d)">alles</a> &mdash;
<a href="javascript:niets(0, %d)">niets</a> &mdash;
<a href="javascript:omkeren(0, %d)">omkeren</a>
`, len(spods), len(spods), len(spods))

	for i, spod := range spods {
		if spod.header != "" {
			if i > 0 {
				fmt.Fprintln(q.w, "</table>")
			}
			var j int
			for j = i + 1; j < len(spods); j++ {
				if spods[j].header != "" {
					break
				}
			}
			fmt.Fprintln(q.w, "<p><b>"+html.EscapeString(spod.header)+"</b><br>")
			fmt.Fprintf(q.w, `
<a href="javascript:alles(%d, %d)">alles</a> &mdash;
<a href="javascript:niets(%d, %d)">niets</a> &mdash;
<a href="javascript:omkeren(%d, %d)">omkeren</a>
<p>
<table class="breed">`,
				i, j, i, j, i, j)
		}
		fmt.Fprintf(q.w, `<tr>
  <td><input type="checkbox" name="i%d" value="t">
  <td><a href="javascript:vb(%d)" target="_blank">vb</a>
  <td>%s
`,
			i,
			i,
			html.EscapeString(spod.text))
	}

	fmt.Fprintf(q.w, `</table>
<div class="spod"></div>
<div id="spodfoot"><div id="spodinfoot">
uitvoer: <select name="out">
<option value="html">HTML</option>
<option value="text">Teksttabel</option>
</select>
<p>
<input type="submit" value="verzenden">
</div></div>
</form>
`)

	fmt.Fprintln(q.w, "</body></html>")

}

func extension_spod_form(q *Context) {

	doHtml := first(q.r, "out") == "html"

	if doHtml {
		writeHead(q, "SPOD", 0)
		defer func() {
			fmt.Fprintf(q.w, "</body>\n</html>\n")
		}()
	} else {
		contentType(q, "text/plain; charset=utf-8")
		nocache(q)
	}

	db := first(q.r, "db")
	if !q.prefixes[db] {
		fmt.Fprintln(q.w, "Ongeldig corpus:", db)
		return
	}

	dbase, err := dbopen()
	if sysErr(err) {
		return
	}
	defer dbase.Close()
	rows, err := dbase.Query("SELECT `nline` from `" + Cfg.Prefix + "_info` WHERE `id` = \"" + db + "\"")
	if sysErr(err) {
		fmt.Fprintln(q.w, err)
		return
	}
	nlines := 0
	for rows.Next() {
		err = rows.Scan(&nlines)
		rows.Close()
		if sysErr(err) {
			fmt.Fprintln(q.w, err)
			return
		}
	}
	if nlines == 0 {
		err = fmt.Errorf("nlines not found")
		sysErr(err)
		fmt.Fprintln(q.w, err)
		return
	}

	mt := first(q.r, "mt")
	if mt != "std" && mt != "dx" {
		mt = "std"
	}

	var allDone bool
	if doHtml {
		fmt.Fprintf(q.w, "Corpus: %s\n<p>\n", html.EscapeString(q.desc[db]))
		allDone = spod_stats(q, db, true)
		if !allDone {
			fmt.Fprintln(q.w, "NA<br>")
		}
	} else {
		fmt.Fprintln(q.w, "# corpus:", q.desc[db])
		fmt.Fprintln(q.w, "# waarde\t         \tlabel\tomschrijving")
		fmt.Fprintln(q.w, "## Stats")
		allDone = spod_stats(q, db, false)
		if !allDone {
			fmt.Fprintln(q.w, "NA")
		}
	}
	if doHtml {
		fmt.Fprintln(q.w, `
<style>
.max100 {
    max-width: 100%;
    overflow: auto;
}
/* The Modal (background) */
.modal {
    display: none; /* Hidden by default */
    position: fixed; /* Stay in place */
    z-index: 1; /* Sit on top */
    left: 0;
    top: 0;
    max-width: 100%;
    width: 100%; /* Full width */
    height: 100%; /* Full height */
    background-color: rgb(0,0,0); /* Fallback color */
    background-color: rgba(0,0,0,0.4); /* Black w/ opacity */
}
.modal-header {
    padding: 1em 2em;
    background-color: #315b7d;
    color: #b5cde1;
}
.modal-header h2 {
    font-size: x-large;
    font-weight: normal;
    text-align: center;
}
/* Modal Body */
.modal-body {
    width: 100%;
    overflow: auto;
    padding: 1em 2em;
    text-align: center;
}
/* Modal Footer */
.modal-footer {
    padding: 1em 2em;
    background-color: #315b7d;
    color: #b5cde1;
    text-align: center;
}
/* Modal Content */
.modal-content {
    position: relative;
    background-color: #fefefe;
    margin: auto;
    padding: 0;
    border: 1px solid #b5cde1;
    width: auto;
    max-width: 94%;
    box-shadow: 0 4px 8px 0 rgba(0,0,0,0.2),0 6px 20px 0 rgba(0,0,0,0.19);
    -webkit-animation-name: animatetop;
    -webkit-animation-duration: 0.4s;
    -moz-animation-name: animatetop;
    -moz-animation-duration: 0.4s;
    -o-animation-name: animatetop;
    -o-animation-duration: 0.4s;
    animation-name: animatetop;
    animation-duration: 0.4s
}

/* Add Animation */
@-webkit-keyframes animatetop {
    from {top: -300px; opacity: 0}
    to {top: 0; opacity: 1}
}
@-moz-keyframes animatetop {
    from {top: -300px; opacity: 0}
    to {top: 0; opacity: 1}
}
@-o-keyframes animatetop {
    from {top: -300px; opacity: 0}
    to {top: 0; opacity: 1}
}

@keyframes animatetop {
    from {top: -300px; opacity: 0}
    to {top: 0; opacity: 1}
}

/* The Close Button */
.close {
    color: #b5cde1;
    float: right;
    font-size: 28px;
    font-weight: bold;
}

.close:hover,
.close:focus {
    color: black;
    text-decoration: none;
    cursor: pointer;
}
.modal-body rect {
  fill: #4682b4;
}

.modal-body text {
  font-family: sans-serif;
  font-size: small;
}
</style>
<div id="myModal" class="modal">
  <!-- Modal content -->
  <div class="modal-content">
    <div class="modal-header">
      <span class="close">&times;</span>
      <h2 id="innerTitle"></h2>
    </div>
    <div class="modal-body"><svg width="960" height="500"></svg></div>
    <div class="modal-footer">Corpus: `+
			html.EscapeString(q.desc[db])+`
    </div>
  </div>
</div>
<script src="https://d3js.org/d3.v4.min.js"></script>
<script>
var modal = document.getElementById('myModal');
var span = document.getElementsByClassName("close")[0];
function wg(idx) {

    var d0 = data[idx];
    var prev = 0;
    var s = "woorden,frequentie";
    for (var i = 0; i < d0.length; i++) {
        while (prev < d0[i][0] - 1) {
            prev++;
            s += "\n" + prev + ",0";
        }
        s += "\n" + d0[i][0] + "," + d0[i][1];
        prev = d0[i][0];
    }
    var dat = d3.csvParse(s);
    for (var i = 0; i < dat.length; i++) {
        if (dat[i].frequentie) {
            dat[i].frequentie = + dat[i].frequentie;
        }
    }

    var step = Math.round(dat.length / 20);
    if (step > 1) {
        for (var i = 0; i < dat.length; i++) {
            if (i % step != step - 1) {
                if (dat[i].woorden) {
                    dat[i].woorden = " ".repeat(i);
                }
            }
        }
    }

    d3.selectAll("svg > *").remove()

	var svg = d3.select("svg"),
	    margin = {top: 20, right: 20, bottom: 40, left: 80},
	    width = +svg.attr("width") - margin.left - margin.right,
	    height = +svg.attr("height") - margin.top - margin.bottom;

	var x = d3.scaleBand().rangeRound([0, width]).padding(0.1),
	    y = d3.scaleLinear().rangeRound([height, 0]);

	var g = svg.append("g")
	    .attr("transform", "translate(" + margin.left + "," + margin.top + ")");

    var ymax = d3.max(dat, function(d) { return d.frequentie; });

    x.domain(dat.map(function(d) { return d.woorden; }));
    y.domain([0, ymax]);

    g.append("g")
        .attr("class", "axis axis--x")
        .attr("transform", "translate(0," + height + ")")
        .call(d3.axisBottom(x))
        .append("text")
        .attr("x", width / 2)
        .attr("y", 0)
        .attr("dy", "37px")
        .attr("text-anchor", "center")
        .attr("fill", "black")
        .text("aantal woorden");

    g.append("g")
        .attr("class", "axis axis--y")
        .call(ymax < 10 ? d3.axisLeft(y).ticks(ymax) : d3.axisLeft(y).tickFormat(d3.format("d")))
        .append("text")
        .attr("transform", "rotate(-90)")
        .attr("y", 6)
        .attr("dy", "0.71em")
        .attr("text-anchor", "end")
        .attr("fill", "black")
        .text("frequentie");

    g.selectAll(".bar")
      .data(dat)
      .enter().append("rect")
      .attr("class", "bar")
      .attr("x", function(d) { return x(d.woorden); })
      .attr("y", function(d) { return y(d.frequentie); })
      .attr("width", x.bandwidth())
      .attr("height", function(d) { return height - y(d.frequentie); });

    d3.selectAll("#innerTitle").text(titles[idx]);
    modal.style.display = "block";
}
span.onclick = function() {
    modal.style.display = "none";
}
window.onclick = function(event) {
    if (event.target == modal) {
        modal.style.display = "none";
    }
}
</script>
<div class="max100"><table><tr><th colspan="2">zinnen<th>items<th>woorden<th></tr>
`)
	} else {
		fmt.Fprintln(q.w, "# zinnen zinnen/totaal\titems\tlabel\tomschrijving\twoordtelling")
	}
	worddata := make([]string, 0)
	wordtitles := make([]string, 0)
	header := ""
	spod_in_use := make(map[string]bool)
	hiscount := -1

	use_his := false
	for i, spod := range spods {
		if spod.special == SPOD_X_HIS && first(q.r, fmt.Sprintf("i%d", i)) == "t" {
			use_his = true
		}
	}

	for i, spod := range spods {
		spod_in_use[spod_fingerprint(i)] = true
		if spod.header != "" {
			header = spod.header
		}
		if (use_his && spod.lbl == "his") || first(q.r, fmt.Sprintf("i%d", i)) == "t" {
			if header != "" {
				if doHtml {
					fmt.Fprintf(q.w, "<tr><th colspan=\"4\"><th class=\"left\">%s</tr>\n", html.EscapeString(header))
				} else {
					fmt.Fprintln(q.w, "##", header)
				}
				header = ""
			}
			lines, items, wcount, done, err := spod_get(q, db, i)
			if err != nil {
				if doHtml {
					fmt.Fprintf(q.w, "<tr><td colspan=\"3\"><b>%s</b></tr>\n", html.EscapeString(err.Error()))

				} else {
					fmt.Fprint(q.w, "#", err)
				}
			} else if done {
				if spod.lbl == "his" {
					hiscount = items
				}
				if doHtml {
					if spod.special == SPOD_X_PARSER || spod.special == SPOD_X_PARSER_CATS || spod.special == SPOD_X_PARSER_SKIPS {
						fmt.Fprintf(q.w, "<tr><td class=\"right\">%d<td class=\"right\">%.2f%%",
							lines, float64(lines)/float64(nlines)*100.0)
						if spod.special == SPOD_X_PARSER {
							fmt.Fprint(q.w, "<td><td><td>")
						} else {
							fmt.Fprintf(q.w, "<td><td><td><span style=\"display:inline-block;width:100px;height:1em;background-color:white\"><span style=\"display:inline-block;width:%dpx;height:1em;background-color:blue\"></span></span>",
								int(float64(lines)/float64(nlines)*100.0+0.5))
						}
					} else if spod.special == SPOD_X_HIS {
						fmt.Fprintf(q.w, "<tr><td><td><td class=\"right\">%d",
							items)
					} else {
						fmt.Fprintf(q.w, "<tr><td class=\"right\">%d<td class=\"right\">%.2f%%<td class=\"right\">%d",
							lines, float64(lines)/float64(nlines)*100.0, items)
					}
				} else {
					v := fmt.Sprintf("%.3g", float64(lines)/float64(nlines))
					fmt.Fprintf(q.w, "%d\t%-15s\t%d", lines, v, items)
				}
			} else {
				wcount = "NA"
				if doHtml {
					if spod.special == SPOD_X_PARSER || spod.special == SPOD_X_PARSER_CATS || spod.special == SPOD_X_PARSER_SKIPS {
						fmt.Fprintf(q.w, "<tr><td class=\"right\">NA<td class=\"right\">NA<td><td><td>")
					} else if spod.special == SPOD_X_HIS {
						fmt.Fprintf(q.w, "<tr><td><td><td class=\"right\">NA")
					} else {
						fmt.Fprintf(q.w, "<tr><td class=\"right\">NA<td class=\"right\">NA<td class=\"right\">NA")
					}
				} else {
					fmt.Fprint(q.w, "NA\tNA       \tNA")
				}
				allDone = false
			}
			if doHtml {
				if spod.special == SPOD_X_PARSER || spod.special == SPOD_X_PARSER_CATS || spod.special == SPOD_X_PARSER_SKIPS {
					fmt.Fprintf(q.w, " %s\n", html.EscapeString(spod.text))
				} else {
					counts := strings.Split(wcount, ",")
					sum := 0
					n := 0
					for _, count := range counts {
						a := strings.Split(count, ":")
						if len(a) == 2 {
							i, _ := strconv.Atoi(a[0])
							j, _ := strconv.Atoi(a[1])
							sum += i * j
							n += j
						}
					}
					if sum == 0 {
						t := "&mdash;"
						if wcount == "NA" {
							t = "NA"
						}
						fmt.Fprintf(q.w, "<td class=\"right\">%s<td>%s\n", t, html.EscapeString(spod.text))
					} else {
						fmt.Fprintf(q.w, "<td class=\"right\"><a href=\"javascript:wg(%d)\">%.1f</a><td>\n",
							len(worddata), float64(sum)/float64(n))

						if spod.special == SPOD_X_HIS && hiscount > 0 {
							fmt.Fprintf(q.w, "<span style=\"display:inline-block;width:100px;height:1em;background-color:white\"><span style=\"display:inline-block;width:%dpx;height:1em;background-color:blue\"></span></span> %.1f%% ",
								int(float64(items)/float64(hiscount)*100.0+0.5),
								float64(items)/float64(hiscount)*100.0)
						}

						fmt.Fprintf(q.w, "%s\n", html.EscapeString(spod.text))

						worddata = append(worddata, "[["+
							strings.Replace(
								strings.Replace(wcount, ",", "],[", -1),
								":", ",", -1)+"]]")
						wordtitles = append(wordtitles, jsstringsEscape(spod.text))
					}
				}
			} else {
				fmt.Fprintf(q.w, "\t%s\t%s\t%s\n", spod.lbl, spod.text, wcount)
			}
		}
	}
	if doHtml {
		fmt.Fprintln(q.w, "</table></div><script>var data = [")
		fmt.Fprintln(q.w, strings.Join(worddata, ",\n"))
		fmt.Fprintln(q.w, "];\nvar titles = [")
		fmt.Fprintln(q.w, strings.Join(wordtitles, ",\n"))
		fmt.Fprintln(q.w, "];</script>")
	}
	if !allDone {
		if doHtml {
			fmt.Fprintln(q.w, "<div class=\"info\"><button type=\"button\" onClick=\"location.reload(true)\">Herladen</button> &mdash; Herlaad de pagina om de ontbrekende resultaten te krijgen</div>")
		} else {
			fmt.Fprintln(q.w, "#\n#  -->  HERLAAD DE PAGINA OM DE ONTBREKENDE RESULTATEN TE KRIJGEN  <--\n#")
		}
	}

	// oude varianten verwijderen
	dirpath := filepath.Join(paqudir, "data", db, "spod")
	files, err := ioutil.ReadDir(dirpath)
	if sysErr(err) {
		return
	}
	spod_in_use["stats"] = true
	for _, file := range files {
		filename := file.Name()
		if !spod_in_use[filename] {
			err := os.Remove(filepath.Join(dirpath, filename))
			sysErr(err)
		}
	}
}

func spod_stats(q *Context, db string, doHtml bool) bool {
	spodMu.Lock()
	defer spodMu.Unlock()

	key := db
	if spodWorking[key] {
		return false
	}

	dirpath := filepath.Join(paqudir, "data", db, "spod")
	os.MkdirAll(dirpath, 0700)
	filename := filepath.Join(dirpath, "stats")

	data, err := ioutil.ReadFile(filename)
	if err == nil {
		if doHtml {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				a := strings.Split(line, "\t")
				if len(a) == 5 {
					fmt.Fprintf(q.w, "%s %s<br>\n", a[0], html.EscapeString(a[4]))
				}
			}
		} else {
			q.w.Write(data)
		}
		return true
	}

	spodWorking[key] = true
	go func() {
		go spod_stats_work(q, db, filename)
		spodMu.Lock()
		delete(spodWorking, key)
		spodMu.Unlock()
	}()

	return false
}

func spod_get(q *Context, db string, item int) (lines int, items int, wcount string, done bool, err error) {
	spodMu.Lock()
	defer spodMu.Unlock()

	key := fmt.Sprintf("%s\t%d", db, item)
	if spodWorking[key] {
		return 0, 0, "", false, nil
	}

	dirpath := filepath.Join(paqudir, "data", db, "spod")
	os.MkdirAll(dirpath, 0700)

	fingerprint := spod_fingerprint(item)
	filename := filepath.Join(dirpath, fingerprint)
	data, err := ioutil.ReadFile(filename)
	if err == nil {
		a := strings.Fields(string(data))
		if len(a) == 3 {
			a = append(a, "1:0")
		}
		if len(a) == 4 {
			if a[0] != spods[item].lbl {
				return 0, 0, "", false, fmt.Errorf("ERROR: invalid label %q", a[0])
			}
			lines, err := strconv.Atoi(a[1])
			if err != nil {
				return 0, 0, "", false, err
			}
			items, err := strconv.Atoi(a[2])
			if err != nil {
				return 0, 0, "", false, err
			}
			return lines, items, a[3], true, nil
		} else {
			return 0, 0, "", false, fmt.Errorf("ERROR: invalid data %q", string(data))
		}
	}

	spodWorking[key] = true
	go func() {
		spod_work(q, key, filename, db, item)
		spodMu.Lock()
		delete(spodWorking, key)
		spodMu.Unlock()
	}()

	return 0, 0, "", false, nil
}

func spod_work(q *Context, key string, filename string, db string, item int) {
	spodSemaphore <- true
	defer func() {
		<-spodSemaphore
	}()

	r, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"http://localhost/?db=%s&xpath=%s&mt=%s&attr1=word_is_&d=1",
			db,
			url.QueryEscape(spods[item].xpath),
			spods[item].method),
		nil)
	if sysErr(err) {
		return
	}
	if sysErr(r.ParseForm()) {
		return
	}
	w := spod_writer{header: make(map[string][]string)}
	dbase, err := dbopen()
	if sysErr(err) {
		return
	}
	defer dbase.Close()
	myQ := Context{
		r:          r,
		w:          &w,
		user:       q.user,
		auth:       q.auth,
		sec:        q.sec,
		quotum:     q.quotum,
		db:         dbase,
		opt_db:     q.opt_db,
		opt_dbmeta: q.opt_dbmeta,
		ignore:     q.ignore,
		prefixes:   q.prefixes,
		myprefixes: q.myprefixes,
		protected:  q.prefixes,
		hasmeta:    q.hasmeta,
		desc:       q.desc,
		lines:      q.lines,
		shared:     q.shared,
		params:     q.params,
		form:       nil,
	}
	xpathstats(&myQ)

	fp, err := os.Create(filename)
	if sysErr(err) {
		return
	}

	scanner := bufio.NewScanner(&w.buffer)

	scanner.Scan()
	s := scanner.Text()

	if spodREerr.MatchString(s) {
		fmt.Fprintf(fp, "ERROR: %s\n", s)
		chLog <- "ERROR spodREerr.MatchString: " + s
	} else {
		match := spodRE.FindStringSubmatch(s)
		if len(match) == 3 {
			fmt.Fprintf(fp, "%s\t%s\t%s\t", spods[item].lbl, match[1], match[2])
			// skip regel
			scanner.Scan()
			counts := make(map[int]int)
			var n int
			var err error
			for scanner.Scan() {
				s := scanner.Text()
				a := strings.Fields(s)
				n, err = strconv.Atoi(a[0])
				if err != nil {
					fmt.Fprintf(fp, "ERROR: no number found in %q\n", s)
					chLog <- "ERROR Atoi: " + err.Error()
					break
				}
				m := 0
				for _, i := range a[2:] {
					if i == "+" {
						m++
					}
				}
				counts[m] = counts[m] + n
			}
			if err == nil {
				keys := make([]int, 0)
				for key := range counts {
					keys = append(keys, key)
				}
				sort.Ints(keys)
				p := ""
				for _, key := range keys {
					fmt.Fprintf(fp, "%s%d:%d", p, key, counts[key])
					p = ","
				}
			}
			fmt.Fprintln(fp)
		} else {
			fmt.Fprintf(fp, "ERROR: no match found in %q\n", s)
			chLog <- "ERROR spodREerr.FindStringSubmatch: " + s
		}
	}
	fp.Close()
}

func (s *spod_writer) Header() http.Header {
	return http.Header(s.header)
}

func (s *spod_writer) Write(b []byte) (int, error) {
	return s.buffer.Write(b)
}

func (s *spod_writer) WriteHeader(i int) {
}

func spod_stats_work(q *Context, dbname string, outname string) {
	spodSemaphore <- true
	defer func() {
		<-spodSemaphore
	}()

	xx := func(err error) bool {
		if err == nil {
			return false
		}
		fp, _ := os.Create(outname)
		fmt.Fprintln(fp, "ERROR:", err)
		fp.Close()
		return true
	}

	lineCount := 0
	wordCount := 0
	runeCount := 0
	tokens := make(map[string]bool)

	db, err := dbopen()
	if xx(err) {
		return
	}
	defer func() {
		if db != nil {
			db.Close()
		}
	}()

	archnames := make([]string, 0)
	rows, err := db.Query(fmt.Sprintf("SELECT `arch` FROM `%s_c_%s_arch`", Cfg.Prefix, dbname))
	if xx(err) {
		return
	}
	for rows.Next() {
		var s string
		if xx(rows.Scan(&s)) {
			rows.Close()
			return
		}
		archnames = append(archnames, s)
	}
	if xx(rows.Err()) {
		return
	}

	filenames := make([]string, 0)
	if len(archnames) == 0 {
		rows, err := db.Query(fmt.Sprintf("SELECT `file` FROM `%s_c_%s_file`", Cfg.Prefix, dbname))
		if xx(err) {
			return
		}
		for rows.Next() {
			var s string
			if xx(rows.Scan(&s)) {
				rows.Close()
				return
			}
			filenames = append(filenames, s)
		}
		if xx(rows.Err()) {
			return
		}
	}

	db.Close()
	db = nil

	var processNode func(node *Node)
	processNode = func(node *Node) {
		if node.Word != "" && node.Pt != "let" {
			wordCount++
			runeCount += utf8.RuneCountInString(node.Word)
			tokens[node.Word] = true
		}
		for _, n := range node.NodeList {
			processNode(n)
		}
	}

	processFile := func(data []byte) error {
		lineCount++
		var alpino Alpino_ds
		err := xml.Unmarshal(data, &alpino)
		if err != nil {
			return err
		}
		processNode(alpino.Node0)
		return nil
	}

	if len(archnames) > 0 {
		for _, archname := range archnames {
			db, err := dbxml.Open(archname)
			if xx(err) {
				return
			}
			docs, err := db.All()
			if xx(err) {
				db.Close()
				return
			}
			for docs.Next() {
				if xx(processFile([]byte(docs.Content()))) {
					db.Close()
					return
				}
			}
			db.Close()
		}
	} else {
		for _, filename := range filenames {
			data, err := ioutil.ReadFile(filename)
			if err != nil {
				fp, err := os.Open(filename + ".gz")
				if xx(err) {
					return
				}
				r, err := gzip.NewReader(fp)
				if xx(err) {
					fp.Close()
					return
				}
				data, err = ioutil.ReadAll(r)
				fp.Close()
				if xx(err) {
					return
				}
			}
			if xx(processFile(data)) {
				return
			}
		}
	}

	fp, _ := os.Create(outname)
	fmt.Fprintf(fp,
		"%8.4f\t\t\ttt\ttypes per tokens\n"+
			"%8.4f\t\t\twz\twoorden per zin\n"+
			"%8.4f\t\t\tlw\tletters per woord\n",
		float64(len(tokens))/float64(wordCount),
		float64(wordCount)/float64(lineCount),
		float64(runeCount)/float64(wordCount))
	fp.Close()
}

func extension_spod_list(q *Context) {

	contentType(q, "text/plain; charset=utf-8")
	nocache(q)

	for _, spod := range spods {
		if spod.header != "" {
			fmt.Fprint(q.w, "\n\n", spod.header, "\n", strings.Repeat("=", len(spod.header)), "\n\n")
		}
		fmt.Fprint(q.w,
			"\n",
			spod.lbl, ": ", spod.text, "\n",
			strings.Repeat("-", len(spod.lbl)+2+len(spod.text)), "\n\n",
			spod.xpath, "\n\n")
	}
}

func spod_fingerprint(item int) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(file__macros__txt+spods[item].xpath+spods[item].method)))
}

func jsstringsEscape(s string) string {
	return "\"" +
		strings.Replace(
			strings.Replace(s, "\\", "\\\\", -1),
			"\"", "\\\"", -1) +
		"\""

}
