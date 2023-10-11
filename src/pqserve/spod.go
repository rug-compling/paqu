package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
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

	"github.com/rug-compling/paqu/internal/dir"
	pqnode "github.com/rug-compling/paqu/internal/node"
	pqspod "github.com/rug-compling/paqu/internal/spod"

	"github.com/pebbe/dbxml"
)

type spodQueueItem struct {
	ch    chan bool
	owner string
}

type spod_writer struct {
	header map[string][]string
	buffer bytes.Buffer
}

var (
	spodSemaphore chan bool
	spodNext      = make(chan bool)
	spodQueue     = make([]*spodQueueItem, 0)
	spodQueueMu   sync.Mutex

	spodMu      sync.Mutex
	spodWorking = make(map[string]bool)
	spodRE      = regexp.MustCompile(`([0-9]+)[^0-9]+([0-9]+)`)
	spodREerr   = regexp.MustCompile(`\[err:`)
)

func spod_init() {
	spodSemaphore = make(chan bool, Cfg.Maxspodjob)
	for i, spod := range pqspod.Spods {
		pqspod.Spods[i].Xpath = strings.TrimSpace(spod.Xpath)
	}
	go spod_starter()
}

func spod_main(q *Context) {
	writeHead(q, "Syntactic profiler of Dutch", 5)

	fmt.Fprintln(q.w, "Syntactic profiler of Dutch<p>")

	if Cfg.Maxspodlines > 0 {
		fmt.Fprintf(q.w, "Corpora met meer dan %s zinnen zijn niet beschikbaar voor dit onderdeel<p>", iformat(Cfg.Maxspodlines))
	}

	fmt.Fprint(q.w, `
<script type="text/javascript" src="jquery.js"></script>
<script type="text/javascript"><!--
  var xpaths = [
`)

	p := ""
	first := -1
	for i, spod := range pqspod.Spods {
		fmt.Fprintf(q.w, "%s['%s', '%s', '%s']", p, url.QueryEscape(spod.Xpath), spod.Method, spod.Lbl)
		p = ",\n"
		if first < 0 && !strings.HasPrefix(spod.Special, "hidden") {
			first = i
		}
	}
	fmt.Fprint(q.w, `
];
  var indexen = {
`)
	p = ""
	for i, spod := range pqspod.Spods {
		fmt.Fprintf(q.w, "%s'%s': %d", p, spod.Lbl, i)
		p = ",\n"
	}
	fmt.Fprint(q.w, `
};
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
var queryvisible = false;
var queryid;
function info(id) {
    uninfo();
    var e = $(id);
    e.show();
    e.css("zIndex", 9999);
    queryvisible = true;
    queryid = id;
}
function uninfo() {
    if (queryvisible) {
        var e = $(queryid);
        e.hide();
        e.css("zIndex", 1);
        queryvisible = false;
    }
}
$(document).mouseup(
  function(e) {
    uninfo();
  });
function getChoices() {
  var res = [];
  for (var i = `, first, `; i < `, len(pqspod.Spods), `; i++) {
    var e = document.getElementById('i' + i);
    if (e.checked) {
      res.push(xpaths[i][2]);
    }
  }
  return res;
}
function setChoices(c) {
  var ee = document.getElementsByClassName("spodblockinner");
  for (var i = 0; i < ee.length; i++) {
    if (! ee[i].classList.contains('hide')) {
      ee[i].classList.add('hide');
    }
  }
  var e = document.getElementById("spodform").elements;
  for (var i = `, first, `; i < `, len(pqspod.Spods), `; i++) {
    e[i+1].checked = false;
  }
  for (var i = 0; i < c.length; i++) {
    var idx = indexen[c[i]];
    if (idx) {
      var el = e[idx + 1];
      el.checked = true;
      while (! el.classList.contains('spodblockinner')) {
        el = el.parentNode;
      }
      if (el.classList.contains('hide')) {
        el.classList.remove('hide');
      }
    }
  }
}
function optsave(n) {
  var c = getChoices();
  localStorage.setItem(
      "paqu-spod-"+n,
      JSON.stringify(c)
  );
  document.getElementById('op'+n).disabled = (c.length == 0);
}
function optload(n) {
  var storageContent = localStorage.getItem("paqu-spod-" + n);
  if (storageContent !== undefined) {
    setChoices(JSON.parse(storageContent) || []);
  }
}
//--></script>
<form id="spodform" action="spodform" method="get" accept-charset="utf-8" target="_blank">
<a href="corpusinfo">[?]</a> corpus: <select name="db">
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
<a href="javascript:alles(%d, %d, true)">alles</a> &mdash;
<a href="javascript:niets(%d, %d, true)">niets</a> &mdash;
<a href="javascript:omkeren(%d, %d)">omkeren</a>
`, first, len(pqspod.Spods), first, len(pqspod.Spods), first, len(pqspod.Spods))

	inTable := false
	blocknum := 0
	for i, spod := range pqspod.Spods {
		if strings.HasPrefix(spod.Special, "hidden") {
			fmt.Fprintf(q.w, `
<input type="hidden" name="i%d" value="t">
`,
				i)
			continue
		}
		spodtext := strings.Replace(spod.Text, "|", "", -1)
		if spod.Header != "" {
			if inTable {
				fmt.Fprintln(q.w, "</table></div></div>")
			} else {
				inTable = true
			}

			var j int
			for j = i + 1; j < len(pqspod.Spods); j++ {
				if pqspod.Spods[j].Header != "" {
					break
				}
			}
			blocknum++
			a := strings.SplitN(spod.Header, "//", 2)
			extra := ""
			if len(a) == 2 {
				// we kunnen hier geen <button> gebruiken omdat die de nummering in de war schopt
				extra = fmt.Sprintf(`
<a href="javascript:info('#h%d')" class="spodi">?</a>
<div id="h%d" class="submenu a9999">
  <div class="queryhelp">
    <h4>%s</h4>
%s
  </div>
</div>
`, blocknum, blocknum, html.EscapeString(a[0]), strings.TrimSpace(a[1]))
			}
			fmt.Fprintf(q.w, `
<div class="spodblock">
<a href="javascript:hider('spodblock%d')">%s</a>%s
<div class="spodblockinner hide" id="spodblock%d">
`, blocknum, spodEscape(a[0]), extra, blocknum)
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
			spodEscape(spodtext))
	}

	fmt.Fprintf(q.w, `</table></div></div>
<div class="spod"></div>
<div id="spodfoot"><div id="spodinfoot">
<div>
hergebruiken<br>
<button type="button" disabled onclick="optload(1)" id="op1">1</button><button type="button" disabled onclick="optload(2)" id="op2">2</button><button type="button" disabled onclick="optload(3)" id="op3">3</button><br>
<button type="button" disabled onclick="optload(4)" id="op4">4</button><button type="button" disabled onclick="optload(5)" id="op5">5</button><button type="button" disabled onclick="optload(6)" id="op6">6</button>
</div>
<div>
bewaren<br>
<button type="button" onclick="optsave(1)">1</button><button type="button" onclick="optsave(2)">2</button><button type="button" onclick="optsave(3)">3</button><br>
<button type="button" onclick="optsave(4)">4</button><button type="button" onclick="optsave(5)">5</button><button type="button" onclick="optsave(6)">6</button>
</div>
uitvoer: <select name="out">
<option value="html">HTML</option>
<option value="text">Teksttabel</option>
<option value="tbl">Gedetailleerde tabel</option>
<option value="tbllen">Gedetailleerde tabel, inclusief lengtes</option>
<option value="tbls">Gedetailleerde tabel, inclusief zinnen</option>
<option value="tbllens">Gedetailleerde tabel, inclusief lengtes en zinnen</option>
</select>
<p>
<input type="submit" value="verzenden">
</div></div>
</form>
<script type="text/javascript"><!--
for (var n = 1; n < 7; n++) {
  var storageContent = localStorage.getItem("paqu-spod-" + n);
  if (storageContent !== undefined) {
    var d = JSON.parse(storageContent);
    if (d && d.length > 0) {
      document.getElementById('op'+n).disabled = false;
    }
  }
}
//--></script>
`)

	fmt.Fprintln(q.w, "</body></html>")
}

func spod_form(q *Context) {
	seen := ""

	outFormat := first(q.r, "out")

	doHtml := outFormat == "html"
	doTable := strings.HasPrefix(outFormat, "tbl")

	if doHtml {
		writeHead(q, "SPOD -- Resultaat", 0)
		defer func() {
			fmt.Fprintf(q.w, "</body>\n</html>\n")
		}()
	} else if doTable {
		contentType(q, "text/tab-separated-values; charset=utf-8")
	} else {
		contentType(q, "text/plain; charset=utf-8")
		nocache(q)
	}

	db := first(q.r, "db")
	if !q.spodprefixes[db] {
		fmt.Fprintln(q.w, "Ongeldig corpus:", db)
		return
	}

	if doTable {
		q.w.Header().Set("Content-Disposition", "attachment; filename=spod-table.tsv")
		spod_table(q, db, strings.HasPrefix(outFormat, "tbllen"), strings.HasSuffix(outFormat, "s"))
		return
	}

	available := map[string]bool{
		"":     true,
		"attr": true,
	}

	rows, err := sqlDB.Query("SELECT `nline`,`owner` from `" + Cfg.Prefix + "_info` WHERE `id` = \"" + db + "\"")
	if sysErr(err) {
		fmt.Fprintln(q.w, err)
		return
	}
	nlines := 0
	var owner string
	for rows.Next() {
		err = rows.Scan(&nlines, &owner)
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
		fmt.Fprintf(q.w, "Corpus: %s\n<p>\n", spodEscape(q.desc[db]))
		allDone = spod_stats(q, db, owner, true)
		if !allDone {
			fmt.Fprintln(q.w, "???<br>")
		}
	} else {
		fmt.Fprintln(q.w, "# corpus:", q.desc[db])
		fmt.Fprintln(q.w, "# waarde\t         \tlabel\tomschrijving")
		fmt.Fprintln(q.w, "## Stats")
		allDone = spod_stats(q, db, owner, false)
		if !allDone {
			fmt.Fprintln(q.w, "???")
		}
	}

	inTable := false
	inAttr := false
	inData := false

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
			spodEscape(q.desc[db])+`
    </div>
  </div>
</div>
<script src="https://d3js.org/d3.v4.min.js"></script>
<script>
function vb(i) {
    window.open("xpath?db=`+db+`&xpath=" + xpaths[i][0] + "&mt=" + xpaths[i][1]);
}
function vb2(i, a, v) {
	if (a == 'topcat') {
        window.open("xpath?db=`+db+`&xpath=%2falpino_ds%2fnode%2fnode%5b%40cat%3d%22" + v + "%22%5d&mt=" + xpaths[i][1]);
	} else {
        window.open("xpath?db=`+db+`&xpath=%2f%2fnode%5b%40" + a + "%3d%22" + v + "%22%5d&mt=" + xpaths[i][1]);
	}
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

    var ticks = [];
    var step = Math.round(dat.length / 20);
    if (step < 1) {
        step = 1;
    }
    for (var i = step; i <= dat.length; i += step) {
        ticks.push(""+i);
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
        .call(d3.axisBottom(x).tickValues(ticks))
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
`)
	}
	worddata := make([]string, 0)
	wordtitles := make([]string, 0)
	header := ""
	spod_in_use := make(map[string]bool)

	for idx, spod := range pqspod.Spods {
		spod_in_use[pqspod.Hash(spod.Lbl)] = true

		if strings.HasPrefix(spod.Special, "hidden") {
			lines, _, _, done, err := spod_get(q, db, idx, owner)
			if err == nil && done {
				available[strings.SplitN(spod.Lbl, "_", 2)[1]] = lines > 0
			}
			continue
		}

		if spod.Header != "" {
			a := strings.SplitN(spod.Header, "//", 2)
			header = strings.TrimSpace(a[0])
		}
		spodtext := spod.Text
		if strings.Contains(spodtext, "||") || !strings.Contains(spodtext, "|") {
			spodtext = strings.Replace(spodtext, "||", "|", 1)
			seen = ""
		}
		if first(q.r, fmt.Sprintf("i%d", idx)) == "t" {
			if !inTable {
				inTable = true
				if doHtml {
					fmt.Fprintln(q.w, `<div class="max100"><table class="spod">`)
				}
			}
			if spod.Special == "attr" {
				if !inAttr {
					inAttr = true
					if doHtml {
						fmt.Fprintln(q.w, `<tr><th colspan="2" class="r b">attributen<th colspan="2" class="r"><th colspan="2"></tr>
`)
					} else {
						fmt.Fprintln(q.w, "# attributen att/totaal\t\tlabel\tomschrijving\t")
					}
				}
				if doHtml {
					fmt.Fprintf(q.w, "<tr><th colspan=\"2\" class=\"r\"><th colspan=\"2\" class=\"r\"><th colspan=\"2\" class=\"left\">%s</tr>\n", spodEscape(spod.Text))
				}
			} else {
				if inAttr {
					inAttr = false
				}
				if !inData {
					inData = true
					if doHtml {
						fmt.Fprintln(q.w, `<tr><th colspan="2" class="r b">zinnen<th class="r b">items<th class="r b">woorden<th><th></tr>
`)
					} else {
						fmt.Fprintln(q.w, "# zinnen zinnen/totaal\titems\tlabel\tomschrijving\twoordtelling")
					}
				}
			}
			if doHtml {
				a := strings.SplitN(spodtext, "|", 2)
				if len(a) == 2 {
					if a[0] == seen {
						spodtext = "— " + a[1]
					}
					seen = a[0]
				}
			} else {
				spodtext = strings.Replace(spodtext, "|", "", 1)
			}
			if header != "" && spod.Special != "attr" {
				if doHtml {
					fmt.Fprintf(q.w, "<tr><th colspan=\"2\" class=\"r\"><th class=\"r\"><th class=\"r\"><th colspan=\"2\" class=\"left\">%s</tr>\n", spodEscape(header))
				} else {
					fmt.Fprintln(q.w, "##", header)
				}
				header = ""
			}
			avail, ok := available[spod.Special]
			if specials := strings.Fields(spod.Special); !ok && len(specials) > 1 {
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
					available[spod.Special] = true
					avail = true
					ok = true
				} else if allok {
					available[spod.Special] = false
					avail = false
					ok = true
				}
			}
			if ok && !avail {
				if doHtml {
					fmt.Fprintf(q.w, "<tr><td colspan=\"4\" class=\"center r\"><em>niet voor dit corpus</em><td><td>%s</tr>\n", spodEscape(spodtext))
				} else {
					fmt.Fprintf(q.w, "### niet voor dit corpus\t%s\n", spodtext)
				}
				continue
			}
			var lines, items int
			var wcount string
			var done bool
			var err error
			if avail {
				lines, items, wcount, done, err = spod_get(q, db, idx, owner)
			}
			if err != nil {
				if doHtml {
					fmt.Fprintf(q.w, "<tr><td colspan=\"3\"><b>%s</b>", spodEscape(err.Error()))
				} else {
					fmt.Fprint(q.w, "#", err)
				}
				allDone = false
			} else if done && available[spod.Special] {
				if doHtml {
					if spod.Special == "attr" {
						for _, line := range strings.Split(wcount, "\n") {
							aa := strings.Fields(line)
							if len(aa) == 4 {
								fmt.Fprintf(q.w, `<tr><td class="right">%s<td class="right r">%s<td colspan="2" class="r">
<td><a href="javascript:vb2(%d, '%s','%s')">vb</a><td>%s</tr>
`,
									aa[0], aa[2], idx, spod.Lbl, url.QueryEscape(aa[3]), html.EscapeString(aa[3]))
							}
						}
					} else if spod.Special == "parser" {
						fmt.Fprintf(q.w, "<tr><td class=\"right\">%d<td class=\"right r\">%.2f%%",
							lines, float64(lines)/float64(nlines)*100.0)
					} else if spod.Special == "his" {
						fmt.Fprintf(q.w, "<tr><td><td class=\"r\"><td class=\"right r\">%d",
							items)
					} else {
						fmt.Fprintf(q.w, "<tr><td class=\"right\">%d<td class=\"right r\">%.2f%%<td class=\"right r\">%d",
							lines, float64(lines)/float64(nlines)*100.0, items)
					}
				} else {
					if spod.Special == "attr" {
						for _, line := range strings.Split(wcount, "\n") {
							aa := strings.Fields(line)
							if len(aa) == 4 {
								fmt.Fprintf(q.w, "%s\t%-15s\t\t%s.%s\tattribuut %s: %s\t\n",
									aa[0], aa[1], spod.Lbl, aa[3], spod.Lbl, aa[3])
							}
						}
					} else {
						v := fmt.Sprintf("%.3g", float64(lines)/float64(nlines))
						fmt.Fprintf(q.w, "%d\t%-15s\t%d", lines, v, items)
					}
				}
			} else {
				wcount = "???"
				if doHtml {
					if spod.Special == "attr" {
						fmt.Fprintf(q.w, "<tr><td class=\"right\">???<td class=\"right r\">???<td class=\"right r\" colspan=\"2\"><td><td>???</tr>")
					} else if spod.Special == "parser" {
						fmt.Fprintf(q.w, "<tr><td class=\"right\">???<td class=\"right r\">???")
					} else if spod.Special == "his" {
						fmt.Fprintf(q.w, "<tr><td><td class=\"r\"><td class=\"right r\">???")
					} else {
						fmt.Fprintf(q.w, "<tr><td class=\"right\">???<td class=\"right r\">???<td class=\"right r\">???")
					}
				} else {
					if spod.Special == "attr" {
						fmt.Fprintf(q.w, "???\t???       \t\t%s.???\tattribuut %s: ???\n", spod.Lbl, spod.Lbl)
						// TODO
					} else {
						fmt.Fprint(q.w, "???\t???       \t???")
					}
				}
				allDone = false
			}
			if doHtml {
				if spod.Special == "parser" {
					fmt.Fprintf(
						q.w,
						"<td class=\"r\"><td class=\"r\"><td><a href=\"javascript:vb(%d)\">vb</a><td>%s\n",
						idx, spodEscape(spodtext))
				} else if spod.Special == "attr" {
					// niks
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
							"<td class=\"right r\">%s<td><a href=\"javascript:vb(%d)\">vb</a><td>%s\n",
							t, idx, spodEscape(spodtext))
					} else {
						fmt.Fprintf(
							q.w,
							"<td class=\"right r\"><a href=\"javascript:wg(%d)\">%.1f</a><td>\n",
							len(worddata), float64(sum)/float64(n))

						fmt.Fprintf(
							q.w,
							"<a href=\"javascript:vb(%d)\">vb</a><td>%s\n",
							idx, spodEscape(spodtext))

						worddata = append(worddata, "[["+
							strings.Replace(
								strings.Replace(wcount, ",", "],[", -1),
								":", ",", -1)+"]]")
						wordtitles = append(wordtitles, jsstringsEscape(spod.Text))
					}
				}
			} else {
				if spod.Special == "attr" {
					// niks
				} else {
					fmt.Fprintf(q.w, "\t%s\t%s\t%s\n", spod.Lbl, spodtext, wcount)
				}
			}
		}
	}
	if doHtml {
		if inTable {
			fmt.Fprintln(q.w, "</table></div>")
		}
		fmt.Fprintln(q.w, "<script>var data = [")
		fmt.Fprintln(q.w, strings.Join(worddata, ",\n"))
		fmt.Fprintln(q.w, "];\nvar titles = [")
		fmt.Fprintln(q.w, strings.Join(wordtitles, ",\n"))
		fmt.Fprintln(q.w, "];\nvar xpaths = [")
		p := ""
		for _, spod := range pqspod.Spods {
			fmt.Fprintf(q.w, "%s['%s', '%s']", p, url.QueryEscape(spod.Xpath), spod.Method)
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
	dirpath := filepath.Join(dir.Data, "data", db, "spod")
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

func spod_stats(q *Context, db string, owner string, doHtml bool) bool {
	spodMu.Lock()
	defer spodMu.Unlock()

	key := db
	if spodWorking[key] {
		return false
	}

	dirpath := filepath.Join(dir.Data, "data", db, "spod")
	os.MkdirAll(dirpath, 0700)
	filename := filepath.Join(dirpath, "stats")

	data, err := ioutil.ReadFile(filename)
	if err == nil {
		if doHtml {
			fmt.Fprintln(q.w, "<table class=\"compact\">")
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				a := strings.Split(line, "\t")
				if len(a) == 5 {
					fmt.Fprintf(q.w, "<tr><td class=\"right\">%s<td>%s\n", a[0], spodEscape(a[4]))
				}
			}
			fmt.Fprintln(q.w, "</table><p>")
		} else {
			q.w.Write(data)
		}
		return true
	}

	spodWorking[key] = true
	go func() {
		go spod_stats_work(q, db, owner, filename)
		spodMu.Lock()
		delete(spodWorking, key)
		spodMu.Unlock()
	}()

	return false
}

func spod_get(q *Context, db string, item int, owner string) (lines int, items int, wcount string, done bool, err error) {
	spodMu.Lock()
	defer spodMu.Unlock()

	key := fmt.Sprintf("%s\t%d", db, item)
	if spodWorking[key] {
		return 0, 0, "", false, nil
	}

	dirpath := filepath.Join(dir.Data, "data", db, "spod")
	os.MkdirAll(dirpath, 0700)

	fingerprint := pqspod.Hash(pqspod.Spods[item].Lbl)
	filename := filepath.Join(dirpath, fingerprint)
	data, err := ioutil.ReadFile(filename)
	if err == nil && len(data) > 0 {
		if pqspod.Spods[item].Special == "attr" {
			return 0, 0, string(data), true, nil
		} else {
			a := strings.Fields(string(data))
			if len(a) == 3 {
				a = append(a, "1:0")
			}
			if len(a) == 4 {
				if a[0] != pqspod.Spods[item].Lbl {
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
	}

	spodWorking[key] = true
	go func() {
		spod_work(q, key, filename, db, owner, item)
		spodMu.Lock()
		delete(spodWorking, key)
		spodMu.Unlock()
	}()

	return 0, 0, "", false, nil
}

func spod_work(q *Context, key string, filename string, db string, owner string, item int) {
	spod_schedule(owner)
	defer func() {
		<-spodSemaphore
	}()

	chLog <- fmt.Sprintf("SPOD: work: %s %s", db, pqspod.Spods[item].Lbl)
	defer func() {
		chLog <- fmt.Sprintf("FINISHED SPOD: work: %s %s", db, pqspod.Spods[item].Lbl)
	}()

	var u string
	onlyone := pqspod.Spods[item].Special == "hidden1"
	attr := pqspod.Spods[item].Special == "attr"
	if onlyone {
		u = fmt.Sprintf("http://localhost/?db=%s&xpath=%s&mt=%s&xn=1",
			db,
			url.QueryEscape(pqspod.Spods[item].Xpath),
			pqspod.Spods[item].Method)
	} else if attr {
		lbl := pqspod.Spods[item].Lbl
		if lbl == "topcat" {
			lbl = "cat"
		}
		u = fmt.Sprintf("http://localhost/?db=%s&xpath=%s&mt=%s&attr1=%s&d=1",
			db,
			url.QueryEscape(pqspod.Spods[item].Xpath),
			pqspod.Spods[item].Method,
			lbl)
	} else {
		u = fmt.Sprintf("http://localhost/?db=%s&xpath=%s&mt=%s&attr1=word_is_&d=1",
			db,
			url.QueryEscape(pqspod.Spods[item].Xpath),
			pqspod.Spods[item].Method)
	}

	r, err := http.NewRequest("GET", u, nil)
	if sysErr(err) {
		return
	}
	if sysErr(r.ParseForm()) {
		return
	}
	w := spod_writer{header: make(map[string][]string)}
	myQ := Context{
		r:            r,
		w:            &w,
		user:         q.user,
		auth:         q.auth,
		sec:          q.sec,
		quotum:       q.quotum,
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
			fmt.Fprintf(fp, "%s\t0\t0\t\n", pqspod.Spods[item].Lbl)
		} else {
			fmt.Fprintf(fp, "%s\t1\t1\t1:1\n", pqspod.Spods[item].Lbl)
		}
		fp.Close()
		return
	}

	if attr {
		xpathstats(&myQ)
		scanner := bufio.NewScanner(&w.buffer)
		scanner.Scan()
		scanner.Scan()
		at := make([]string, 0)
		sum := 0
		for scanner.Scan() {
			line := strings.Fields(scanner.Text())
			at = append(at, line[2]+"\t"+line[0])
			a, _ := strconv.Atoi(line[0])
			sum += a
		}
		sort.Strings(at)
		for _, line := range at {
			aa := strings.Fields(line)
			n, _ := strconv.Atoi(aa[1])
			fmt.Fprintf(fp, "%d\t%.3f\t%.2f%%\t%s\n", n, float64(n)/float64(sum), float64(n)/float64(sum)*100.0, aa[0])
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
			fmt.Fprintf(fp, "%s\t%s\t%s\t", pqspod.Spods[item].Lbl, match[1], match[2])
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

func spod_stats_work(q *Context, dbname string, owner string, outname string) {
	spod_schedule(owner)
	defer func() {
		<-spodSemaphore
	}()

	chLog <- fmt.Sprintf("SPOD: stats: %s", dbname)
	defer func() {
		chLog <- fmt.Sprintf("FINISHED SPOD: stats: %s", dbname)
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

	archnames := make([]string, 0)
	rows, err := sqlDB.Query(fmt.Sprintf("SELECT `arch` FROM `%s_c_%s_arch`", Cfg.Prefix, dbname))
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
		rows, err := sqlDB.Query(fmt.Sprintf("SELECT `file` FROM `%s_c_%s_file`", Cfg.Prefix, dbname))
		if xx(err) {
			return
		}
		for rows.Next() {
			var s string
			if xx(rows.Scan(&s)) {
				rows.Close()
				return
			}
			if strings.HasPrefix(s, "$$/") {
				s = filepath.Join(dir.Data, "data", dbname, "xml", s[3:])
			}
			filenames = append(filenames, s)
		}
		if xx(rows.Err()) {
			return
		}
	}

	var processNode func(node *pqnode.Node)
	processNode = func(node *pqnode.Node) {
		if node.Word != "" && node.Pt != "let" {
			wordCount++
			runeCount += utf8.RuneCountInString(strings.Replace(node.Word, "ij", "y", -1))
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
		"%8d\t\t\tzinnen\tzinnen\n"+
			"%8d\t\t\twoorden\twoorden\n"+
			"%8.4f\t\t\ttt\ttypes per token\n"+
			"%8.4f\t\t\twz\twoorden per zin\n"+
			"%8.4f\t\t\tlw\tletters per woord\n",
		lineCount,
		wordCount,
		float64(len(tokens))/float64(wordCount),
		float64(wordCount)/float64(lineCount),
		float64(runeCount)/float64(wordCount))
	fp.Close()
}

func spod_list(q *Context) {
	contentType(q, "text/plain; charset=utf-8")
	nocache(q)

	for _, spod := range pqspod.Spods {
		if strings.HasPrefix(spod.Special, "hidden") {
			continue
		}
		header := spod.Header
		if i := strings.Index(header, "//"); i > 0 {
			header = header[:i]
		}
		spodtext := strings.Replace(spod.Text, "|", "", -1)
		if header != "" {
			fmt.Fprint(q.w, "\n\n", header, "\n", strings.Repeat("=", len(header)), "\n\n")
		}
		fmt.Fprint(q.w,
			"\n",
			spod.Lbl, ": ", spodtext, "\n",
			strings.Repeat("-", len(spod.Lbl)+2+len(spodtext)), "\n\n",
			spod.Xpath, "\n\n")
	}
}

func jsstringsEscape(s string) string {
	return "\"" +
		strings.Replace(
			strings.Replace(
				strings.Replace(s, "\\", "\\\\", -1),
				"\"", "\\\"", -1),
			"|", "", -1) +
		"\""
}

func spodEscape(s string) string {
	s = html.EscapeString(s)
	if strings.Contains(s, "|") {
		a := strings.SplitN(s, "|", 2)
		s = "<u>" + a[0] + "</u>" + a[1]
	}
	return s
}

func spod_schedule(owner string) {
	ch := make(chan bool)
	item := spodQueueItem{
		ch:    ch,
		owner: owner,
	}

	spodQueueMu.Lock()
	if n := len(spodQueue); n == 0 || spodQueue[n-1].owner == owner {
		spodQueue = append(spodQueue, &item)
	} else {
		idx := 0
		for i := n - 1; i >= 0; i-- {
			if spodQueue[i].owner == owner {
				idx = i + 1
				break
			}
		}
		seen := make(map[string]bool)
		for i := idx; i < n; i++ {
			if o := spodQueue[i].owner; !seen[o] {
				seen[o] = true
				idx = i + 1
			}
		}
		if idx >= n {
			spodQueue = append(spodQueue, &item)
		} else {
			spodQueue = append(spodQueue, nil)
			copy(spodQueue[idx+1:], spodQueue[idx:n])
			spodQueue[idx] = &item
		}
	}
	spodQueueMu.Unlock()

	go func() {
		spodNext <- true
	}()

	<-ch
}

func spod_starter() {
	for {
		<-spodNext
		spodSemaphore <- true
		spodQueueMu.Lock()
		close(spodQueue[0].ch)
		spodQueue = spodQueue[1:]
		spodQueueMu.Unlock()
	}
}
