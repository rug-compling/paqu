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
	SPOD_STD = "std"
	SPOD_DX  = "dx"
)

type Spod struct {
	header  string
	xpath   string
	method  string
	lbl     string
	text    string
	special string
}

type spod_writer struct {
	header map[string][]string
	buffer bytes.Buffer
}

var (
	spods = []Spod{
		// hidden objecten moeten aan het begin, omdat anderen ervan afhankelijk zijn
		{
			"",
			`//node [@pos="verb"]`, // spatie i.v.m. mogelijk duplicaat
			SPOD_STD,
			"has_pos_verb",
			"has @pos=verb",
			"hidden1",
		},
		{
			"",
			"//parser [@cats and @skips]", // spatie i.v.m. mogelijk duplicaat
			SPOD_STD,
			"has_parser",
			"has <parser>",
			"hidden1",
		},
		{
			"",
			"//node [@his]", // spatie i.v.m. mogelijk duplicaat
			SPOD_STD,
			"has_his",
			"has @his",
			"hidden1",
		},
		{
			"",
			`//node [@word="?"]`, // spatie i.v.m. mogelijk duplicaat
			SPOD_STD,
			"has_qm",
			"has @word=?",
			"hidden1",
		},
		{
			"",
			`//node [@stype="ynquestion"]`, // spatie i.v.m. mogelijk duplicaat
			SPOD_STD,
			"has_yn",
			"has @stype=ynquestion",
			"hidden1",
		},
		{
			"Hoofdzinnen",
			`//node[@cat="smain"]`,
			SPOD_STD,
			"smain",
			"mededelende hoofdzinnen",
			"",
		},
		{
			"",
			`//node[@cat="whq"]`,
			SPOD_STD,
			"whq",
			"vraagzinnen (wh)",
			"",
		},
		{
			"",
			`%PQ_janee_vragen%`,
			SPOD_STD,
			"janee",
			"ja/nee vragen",
			"qm -yn",
		},
		{
			"",
			`%PQ_imperatieven%`,
			SPOD_STD,
			"imp",
			"imperatieven",
			"",
		},
		{
			"Bijzinnen",
			`%PQ_ingebedde_vraagzinnen%`,
			SPOD_STD,
			"whsub",
			"ingebedde vraagzinnen",
			"",
		},
		{
			"",
			`//node[@cat='cp' and node[@rel='body' and @cat='ssub']]`,
			SPOD_STD,
			"ssub",
			"finiete bijzinnen",
			"",
		},
		{
			"",
			`//node[@cat='cp' and node[@rel='cmp' and @lemma='dat'] and node[@rel='body' and @cat='ssub']]`,
			SPOD_STD,
			"ssubdat",
			"finiete bijzinnen met \"dat\"",
			"",
		},
		{
			"",
			`//node[@cat='cp' and node[@rel='cmp' and @lemma='of'] and node[@rel='body' and @cat='ssub']]`,
			SPOD_STD,
			"ssubof",
			"finiete bijzinnen met \"of\"",
			"",
		},
		{
			"",
			`//node[@cat='cp' and node[@rel='cmp' and not(@lemma='of' or @lemma='dat')] and node[@rel='body' and @cat='ssub']]`,
			SPOD_STD,
			"ssubcmp",
			"finiete bijzinnen met andere voegwoorden",
			"",
		},
		{
			"",
			`//node[@cat='oti']`,
			SPOD_STD,
			"oti",
			"infiniete bijzinnen met \"om\"",
			"",
		},
		{
			"",
			`//node[@cat='ti' and not(../@cat='oti' or ../@cat='cp')]`,
			SPOD_STD,
			"tite",
			"infiniete bijzinnen met alleen \"te\"",
			"",
		},
		{
			"",
			`//node[@cat='cp' and node[@rel='body' and @cat='ti']]`,
			SPOD_STD,
			"ti",
			"infiniete bijzinnen met ander voorzetsel",
			"",
		},
		{
			"",
			`%PQ_relatieve_bijzinnen%`,
			SPOD_STD,
			"relssub",
			"relatieve bijzinnen",
			"",
		},
		{
			"",
			`%PQ_free_relatives%`,
			SPOD_STD,
			"whrel",
			"free relatives",
			"",
		},
		{
			"Woordgroepen",
			`//node[%PQ_np%]`,
			SPOD_STD,
			"np",
			"np",
			"",
		},
		{
			"",
			`//node[@cat='pp']`,
			SPOD_STD,
			"pp",
			"pp",
			"",
		},
		{
			"",
			`//node[@cat='ap' or @pt='adj' and not(@rel='hd')]`,
			SPOD_STD,
			"ap",
			"ap",
			"",
		},
		{
			"",
			`//node[@cat='advp' or @pt='bw' and not(@rel='hd')]`,
			SPOD_STD,
			"advp",
			"advp",
			"",
		},
		{
			"Werkwoorden",
			`//node[@rel="svp" and @cat]`,
			SPOD_STD,
			"vwuit",
			"vaste werkwoordelijke uitdrukkingen",
			"",
		},
		{
			"",
			`%PQ_groen%`,
			SPOD_STD,
			"groen",
			"groene werkwoordsvolgorde",
			"",
		},
		{
			"",
			`%PQ_rood%`,
			SPOD_STD,
			"rood",
			"rode werkwoordsvolgorde",
			"",
		},
		{
			"",
			`//node[%PQ_dep_node_in_verbcluster%]`,
			SPOD_STD,
			"wwclus",
			"werkwoordsclusters",
			"",
		},
		{
			"",
			`%PQ_cross_serial_verbcluster%`,
			SPOD_STD,
			"accinf",
			"accusativus cum infinitivo",
			"pos",
		},
		{
			"",
			`//node[%PQ_passive%]`,
			SPOD_STD,
			"passive",
			"passief",
			"pos",
		},
		{
			"",
			`//node[%PQ_impersonal_passive%]`,
			SPOD_STD,
			"nppas",
			"niet-persoonlijke passief",
			"pos",
		},
		{
			"Inbedding",
			`//node[%PQ_finiete_inbedding0%]`,
			SPOD_STD,
			"inb0",
			"geen inbedding",
			"",
		},
		{
			"",
			`//node[%PQ_finiete_inbedding1%]`,
			SPOD_STD,
			"inb1",
			"minstens 1 finiete zinsinbedding",
			"",
		},
		{
			"",
			`//node[%PQ_finiete_inbedding2%]`,
			SPOD_STD,
			"inb2",
			"minstens 2 finiete zinsinbeddingen",
			"",
		},
		{
			"",
			`//node[%PQ_finiete_inbedding3%]`,
			SPOD_STD,
			"inb3",
			"minstens 3 finiete zinsinbeddingen",
			"",
		},
		{
			"",
			`//node[%PQ_finiete_inbedding4%]`,
			SPOD_STD,
			"inb4",
			"minstens 4 finiete zinsinbeddingen",
			"",
		},
		{
			"",
			`//node[%PQ_finiete_inbedding5%]`,
			SPOD_STD,
			"inb5",
			"minstens 5 finiete zinsinbeddingen",
			"",
		},
		{
			"",
			`//node[%PQ_finiete_inbedding6%]`,
			SPOD_STD,
			"inb6",
			"minstens 6 finiete zinsinbeddingen",
			"",
		},
		{
			"",
			`//node[%PQ_finiete_inbedding7%]`,
			SPOD_STD,
			"inb7",
			"minstens 7 finiete zinsinbeddingen",
			"",
		},
		{
			"",
			`//node[%PQ_finiete_inbedding8%]`,
			SPOD_STD,
			"inb8",
			"minstens 8 finiete zinsinbeddingen",
			"",
		},
		{
			"Topicalizatie en Extractie",
			`%PQ_vorfeld_np_subject%`,
			SPOD_STD,
			"nptsub",
			"np-topic is subject",
			"",
		},
		{
			"",
			`%PQ_vorfeld_np_no_subject%`,
			SPOD_STD,
			"nptnsub",
			"np-topic is geen subject",
			"",
		},
		{
			"",
			`%PQ_vorfeld_non_local%`,
			SPOD_STD,
			"tnonloc",
			"topic is niet lokaal",
			"",
		},
		{
			"",
			`%PQ_local_extraction%`,
			SPOD_STD,
			"locext",
			"lokale extractie",
			"",
		},
		{
			"",
			`%PQ_non_local_extraction%`,
			SPOD_STD,
			"nlocext",
			"niet-lokale extractie",
			"",
		},
		{
			"Parser succes",
			`//parser[@cats="1" and @skips="0"]`,
			SPOD_STD,
			"ok",
			"volledige parse",
			"parser",
		},
		{
			"Parser succes: geparste delen",
			`//parser[@cats="0"]`,
			SPOD_STD,
			"cats0",
			"geen enkel deel is geparst",
			"parser",
		},
		{
			"",
			`//parser[@cats="1"]`,
			SPOD_STD,
			"cats1",
			"parse bestaat uit één deel",
			"parser",
		},
		{
			"",
			`//parser[@cats="2"]`,
			SPOD_STD,
			"cats2",
			"parse bestaat uit twee losse delen",
			"parser",
		},
		{
			"",
			`//parser[@cats="3"]`,
			SPOD_STD,
			"cats3",
			"parse bestaat uit drie losse delen",
			"parser",
		},
		{
			"",
			`//parser[number(@cats) > 3]`,
			SPOD_STD,
			"cats4",
			"parse bestaat uit vier of meer losse delen",
			"parser",
		},
		{
			"parser success: overgeslagen woorden",
			`//parser[@skips="0"]`,
			SPOD_STD,
			"skips0",
			"geen enkel woord is overgeslagen",
			"parser",
		},
		{
			"",
			`//parser[@skips="1"]`,
			SPOD_STD,
			"skips1",
			"een van de woorden is overgeslagen",
			"parser",
		},
		{
			"",
			`//parser[@skips="2"]`,
			SPOD_STD,
			"skips2",
			"twee van de woorden zijn overgeslagen",
			"parser",
		},
		{
			"",
			`//parser[@skips="3"]`,
			SPOD_STD,
			"skips3",
			"drie van de woorden zijn overgeslagen",
			"parser",
		},
		{
			"",
			`//parser[number(@skips) > 3]`,
			SPOD_STD,
			"skips4",
			"vier of meer van de woorden zijn overgeslagen",
			"parser",
		},
		{
			"Onbekende woorden",
			`//node[@his]`,
			SPOD_STD,
			"his",
			"\"alle\" woorden (nodes met attribuut @his)",
			"his",
		},
		{
			"",
			`//node[@his="normal"]`,
			SPOD_STD,
			"normal",
			"woorden uit het woordenboek of de namenlijst",
			"his",
		},
		{
			"",
			`//node[@his and not(@his="normal")]`,
			SPOD_STD,
			"onbeken",
			"woorden niet direct uit het woordenboek",
			"his",
		},
		{
			"",
			`//node[@his="compound"]`,
			SPOD_STD,
			"compoun",
			"woorden herkend als samenstelling",
			"his",
		},
		{
			"",
			`//node[@his="name"]`,
			SPOD_STD,
			"name",
			"woorden herkend als naam (maar niet uit namenlijst)",
			"his",
		},
		{
			"",
			`//node[@his and not(@his="normal" or @his="compound" or @his="name")]`,
			SPOD_STD,
			"noun",
			"onbekende woorden die niet als samenstelling of naam werden herkend",
			"his",
		},
	}
	spodMu        sync.Mutex
	spodSemaphore chan bool
	spodWorking   = make(map[string]bool)
	spodRE        = regexp.MustCompile(`([0-9]+)[^0-9]+([0-9]+)`)
	spodREerr     = regexp.MustCompile(`\[err:`)
)

func spod_init() {
	spodSemaphore = make(chan bool, Cfg.Maxspodjob)
}

func spod_main(q *Context) {

	writeHead(q, "Syntactic profile of Dutch", 5)

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
    window.open("xpath?db=" + e[0].value + "&xpath=" + xpaths[i][0] + "&mt=" + xpaths[i][1]);
}
function alles(n, m, o) {
    var e = document.getElementById("spodform").elements
    for (var i = n; i < m; i++) {
        e[i+1].checked = true;
    }
    if (o) {
        var ee = document.getElementsByClassName("spodblockinner");
        for (var i = 0; i < ee.length; i++) {
            if (ee[i].classList.contains('hide')) {
                ee[i].classList.remove('hide');
            }
        }
    }
}
function niets(n, m, o) {
    var e = document.getElementById("spodform").elements
    for (var i = n; i < m; i++) {
        e[i+1].checked = false;
    }
    if (o) {
        var ee = document.getElementsByClassName("spodblockinner");
        for (var i = 0; i < ee.length; i++) {
            if (! ee[i].classList.contains('hide')) {
                ee[i].classList.add('hide');
            }
        }
    }
}
function omkeren(n, m) {
    var e = document.getElementById("spodform").elements
    for (var i = n; i < m; i++) {
        e[i+1].checked = !e[i+1].checked;
    }
}
function hider(id) {
    var e = document.getElementById(id);
    e.classList.toggle("hide");
}
//--></script>
<form id="spodform" action="spodform" method="get" accept-charset="utf-8" target="_blank">
corpus: <select name="db">
`)

	html_opts(q, q.opt_dbspod, getprefix(q), "corpus")
	fmt.Fprintln(q.w, "</select>")
	if q.auth {
		fmt.Fprintln(q.w, "<a href=\"corpuslijst\">meer/minder</a>")
	}

	fmt.Fprintf(q.w, `
<p>
<a href="spodlist" target="_blank">lijst van queries</a>
<p>
<a href="javascript:alles(0, %d, true)">alles</a> &mdash;
<a href="javascript:niets(0, %d, true)">niets</a> &mdash;
<a href="javascript:omkeren(0, %d)">omkeren</a>
`, len(spods), len(spods), len(spods))

	inTable := false
	blocknum := 0
	for i, spod := range spods {
		if strings.HasPrefix(spod.special, "hidden") {
			fmt.Fprintf(q.w, `
<input type="hidden" name="i%d" value="t">
`,
				i)
			continue
		}
		if spod.header != "" {
			if inTable {
				fmt.Fprintln(q.w, "</table></div></div>")
			} else {
				inTable = true
			}

			var j int
			for j = i + 1; j < len(spods); j++ {
				if spods[j].header != "" {
					break
				}
			}
			blocknum++
			fmt.Fprintf(q.w, `
<div class="spodblock">
<a href="javascript:hider('spodblock%d')">%s</a>
<div class="spodblockinner hide" id="spodblock%d">
`, blocknum, html.EscapeString(spod.header), blocknum)
			fmt.Fprintf(q.w, `
<a href="javascript:alles(%d, %d)">alles</a> &mdash;
<a href="javascript:niets(%d, %d)">niets</a> &mdash;
<a href="javascript:omkeren(%d, %d)">omkeren</a>
<p>
<table class="breed">`,
				i, j, i, j, i, j)
		}
		fmt.Fprintf(q.w, `<tr>
  <td><input type="checkbox" name="i%d" id="i%d" value="t">
  <td><a href="javascript:vb(%d)">vb</a>
  <td><label for="i%d">%s</label>
`,
			i,
			i,
			i,
			i,
			html.EscapeString(spod.text))
	}

	fmt.Fprintf(q.w, `</table></div></div>
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

func spod_form(q *Context) {

	doHtml := first(q.r, "out") == "html"

	if doHtml {
		writeHead(q, "SPOD -- Resultaat", 0)
		defer func() {
			fmt.Fprintf(q.w, "</body>\n</html>\n")
		}()
	} else {
		contentType(q, "text/plain; charset=utf-8")
		nocache(q)
	}

	db := first(q.r, "db")
	if !q.spodprefixes[db] {
		fmt.Fprintln(q.w, "Ongeldig corpus:", db)
		return
	}

	available := map[string]bool{
		"": true,
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
			fmt.Fprintln(q.w, "???<br>")
		}
	} else {
		fmt.Fprintln(q.w, "# corpus:", q.desc[db])
		fmt.Fprintln(q.w, "# waarde\t         \tlabel\tomschrijving")
		fmt.Fprintln(q.w, "## Stats")
		allDone = spod_stats(q, db, false)
		if !allDone {
			fmt.Fprintln(q.w, "???")
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
    <div class="modal-footer">`+
			html.EscapeString(q.desc[db])+`
    </div>
  </div>
</div>
<script src="https://d3js.org/d3.v4.min.js"></script>
<script>
function vb(i) {
    window.open("xpath?db=`+db+`&xpath=" + xpaths[i][0] + "&mt=" + xpaths[i][1]);
}

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
        var space = "";
        for (var i = 0; i < dat.length; i++) {
            if (i % step != step - 1) {
                if (dat[i].woorden) {
                    space += " ";
                    dat[i].woorden = space;
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

	for idx, spod := range spods {
		spod_in_use[spod_fingerprint(idx)] = true

		if strings.HasPrefix(spod.special, "hidden") {
			lines, _, _, done, err := spod_get(q, db, idx)
			if err == nil && done {
				available[strings.Split(spod.lbl, "_")[1]] = lines > 0
			}
			continue
		}

		if spod.header != "" {
			header = spod.header
		}
		if first(q.r, fmt.Sprintf("i%d", idx)) == "t" {
			if header != "" {
				if doHtml {
					fmt.Fprintf(q.w, "<tr><th colspan=\"4\"><th class=\"left\">%s</tr>\n", html.EscapeString(header))
				} else {
					fmt.Fprintln(q.w, "##", header)
				}
				header = ""
			}
			avail, ok := available[spod.special]
			if specials := strings.Fields(spod.special); !ok && len(specials) > 1 {
				allok := true
				someavail := false
				for _, spec := range specials {
					neg := false
					if spec[0] == '-' {
						spec = spec[1:]
						neg = true
					}
					av, ok := available[spec]
					if !ok {
						allok = false
						continue
					}
					if neg {
						av = !av
					}
					if av {
						someavail = true
						break
					}
				}
				if someavail {
					available[spod.special] = true
					avail = true
					ok = true
				} else if allok {
					available[spod.special] = false
					avail = false
					ok = true
				}
			}
			if ok && !avail {
				if doHtml {
					fmt.Fprintf(q.w, "<tr><td colspan=\"4\" class=\"center\"><em>niet voor dit corpus</em><td>%s</tr>\n", spod.text)
				} else {
					fmt.Fprintf(q.w, "### niet voor dit corpus\t%s\n", spod.text)
				}
				continue
			}
			var lines, items int
			var wcount string
			var done bool
			var err error
			if avail {
				lines, items, wcount, done, err = spod_get(q, db, idx)
			}
			if err != nil {
				if doHtml {
					fmt.Fprintf(q.w, "<tr><td colspan=\"3\"><b>%s</b>", html.EscapeString(err.Error()))

				} else {
					fmt.Fprint(q.w, "#", err)
				}
				allDone = false
			} else if done && available[spod.special] {
				if doHtml {
					if spod.special == "parser" {
						fmt.Fprintf(q.w, "<tr><td class=\"right\">%d<td class=\"right\">%.2f%%",
							lines, float64(lines)/float64(nlines)*100.0)
					} else if spod.special == "his" {
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
				wcount = "???"
				if doHtml {
					if spod.special == "parser" {
						fmt.Fprintf(q.w, "<tr><td class=\"right\">???<td class=\"right\">???")
					} else if spod.special == "his" {
						fmt.Fprintf(q.w, "<tr><td><td><td class=\"right\">???")
					} else {
						fmt.Fprintf(q.w, "<tr><td class=\"right\">???<td class=\"right\">???<td class=\"right\">???")
					}
				} else {
					fmt.Fprint(q.w, "???\t???       \t???")
				}
				allDone = false
			}
			if doHtml {
				if spod.special == "parser" {
					fmt.Fprintf(
						q.w,
						"<td><td><td><a href=\"javascript:vb(%d)\">vb</a> &nbsp; %s\n",
						idx, html.EscapeString(spod.text))
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
						if wcount == "???" {
							t = "???"
						}
						fmt.Fprintf(
							q.w,
							"<td class=\"right\">%s<td><a href=\"javascript:vb(%d)\">vb</a> &nbsp; %s\n",
							t, idx, html.EscapeString(spod.text))
					} else {
						fmt.Fprintf(
							q.w,
							"<td class=\"right\"><a href=\"javascript:wg(%d)\">%.1f</a><td>\n",
							len(worddata), float64(sum)/float64(n))

						fmt.Fprintf(
							q.w,
							"<a href=\"javascript:vb(%d)\">vb</a> &nbsp; %s\n",
							idx, html.EscapeString(spod.text))

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
		fmt.Fprintln(q.w, "];\nvar xpaths = [")
		p := ""
		for _, spod := range spods {
			fmt.Fprintf(q.w, "%s['%s', '%s']", p, url.QueryEscape(spod.xpath), spod.method)
			p = ",\n"
		}
		fmt.Fprintln(q.w, "];</script>")
	}
	if !allDone {
		if doHtml {
			fmt.Fprintln(q.w, "<div class=\"footspace\"></div><div id=\"spodinfo\" class=\"info\"><button type=\"button\" onClick=\"location.reload(true)\">Herladen</button> &mdash; Herlaad de pagina om de ontbrekende resultaten te krijgen</div>")
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
				os.Remove(filename)
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
			os.Remove(filename)
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

	var u string
	onlyone := spods[item].special == "hidden1"
	if onlyone {
		u = fmt.Sprintf("http://localhost/?db=%s&xpath=%s&mt=%s&xn=1",
			db,
			url.QueryEscape(spods[item].xpath),
			spods[item].method)
	} else {
		u = fmt.Sprintf("http://localhost/?db=%s&xpath=%s&mt=%s&attr1=word_is_&d=1",
			db,
			url.QueryEscape(spods[item].xpath),
			spods[item].method)
	}

	r, err := http.NewRequest("GET", u, nil)
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
		r:            r,
		w:            &w,
		user:         q.user,
		auth:         q.auth,
		sec:          q.sec,
		quotum:       q.quotum,
		db:           dbase,
		opt_db:       q.opt_db,
		opt_dbmeta:   q.opt_dbmeta,
		ignore:       q.ignore,
		prefixes:     q.prefixes,
		myprefixes:   q.myprefixes,
		spodprefixes: q.spodprefixes,
		protected:    q.protected,
		hasmeta:      q.hasmeta,
		desc:         q.desc,
		lines:        q.lines,
		shared:       q.shared,
		params:       q.params,
		form:         nil,
	}

	fp, err := os.Create(filename)
	if sysErr(err) {
		return
	}

	if onlyone {
		xpath(&myQ)
		if strings.Contains(w.buffer.String(), "<!--NOMATCH-->") {
			fmt.Fprintf(fp, "%s\t0\t0\t\n", spods[item].lbl)
		} else {
			fmt.Fprintf(fp, "%s\t1\t1\t1:1\n", spods[item].lbl)
		}
		fp.Close()
		return
	}

	xpathstats(&myQ)

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
			db, err := dbxml.OpenRead(archname)
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

func spod_list(q *Context) {

	contentType(q, "text/plain; charset=utf-8")
	nocache(q)

	for _, spod := range spods {
		if strings.HasPrefix(spod.special, "hidden") {
			continue
		}
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
	rules := getMacrosRules(&Context{})
	query := macroKY.ReplaceAllStringFunc(spods[item].xpath, func(s string) string {
		return rules[s[1:len(s)-1]]
	})
	query = strings.Join(strings.Fields(query), " ")
	return fmt.Sprintf("%x", md5.Sum([]byte(query+spods[item].method)))
}

func jsstringsEscape(s string) string {
	return "\"" +
		strings.Replace(
			strings.Replace(s, "\\", "\\\\", -1),
			"\"", "\\\"", -1) +
		"\""

}
