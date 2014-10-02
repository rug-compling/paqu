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
	"sort"
	"strings"
	"time"
)

var NodeTags = []string{
	"aform",
	"begin",
	"buiging",
	"case",
	"cat",
	"comparative",
	"conjtype",
	"def",
	"dial",
	"end",
	"frame",
	"gen",
	"genus",
	"getal",
	"getal-n",
	"graad",
	"id",
	"index",
	"infl",
	"lcat",
	"lemma",
	"lwtype",
	"mwu_root",
	"mwu_sense",
	"naamval",
	"neclass",
	"npagr",
	"ntype",
	"num",
	"numtype",
	"pb",
	"pdtype",
	"per",
	"persoon",
	"pos",
	"positie",
	"postag",
	"pt",
	"pvagr",
	"pvtijd",
	"refl",
	"rel",
	"root",
	"sc",
	"sense",
	"special",
	"spectype",
	"status",
	"tense",
	"vform",
	"vwtype",
	"vztype",
	"wh",
	"wk",
	"word",
	"wvorm",
}

type AttrItem struct {
	n int
	a string
}

type AttrItems []*AttrItem

type ValueItem struct {
	begin int
	end   int
	value string
}

type ValueItems []*ValueItem

type Alpino_ds_full_node struct {
	XMLName xml.Name  `xml:"alpino_ds"`
	Node0   *FullNode `xml:"node"`
}

type FullNode struct {
	Aform       string `xml:"aform,attr"`
	Begin       int    `xml:"begin,attr"`
	Buiging     string `xml:"buiging,attr"`
	Case        string `xml:"case,attr"`
	Cat         string `xml:"cat,attr"`
	Comparative string `xml:"comparative,attr"`
	Conjtype    string `xml:"conjtype,attr"`
	Def         string `xml:"def,attr"`
	Dial        string `xml:"dial,attr"`
	End         int    `xml:"end,attr"`
	Frame       string `xml:"frame,attr"`
	Gen         string `xml:"gen,attr"`
	Genus       string `xml:"genus,attr"`
	Getal       string `xml:"getal,attr"`
	GetalN      string `xml:"getal-n,attr"`
	Graad       string `xml:"graad,attr"`
	Id          string `xml:"id,attr"`
	Index       string `xml:"index,attr"`
	Infl        string `xml:"infl,attr"`
	Lcat        string `xml:"lcat,attr"`
	Lemma       string `xml:"lemma,attr"`
	Lwtype      string `xml:"lwtype,attr"`
	MwuRoot     string `xml:"mwu_root,attr"`
	MwuSense    string `xml:"mwu_sense,attr"`
	Naamval     string `xml:"naamval,attr"`
	Neclass     string `xml:"neclass,attr"`
	Npagr       string `xml:"npagr,attr"`
	Ntype       string `xml:"ntype,attr"`
	Num         string `xml:"num,attr"`
	Numtype     string `xml:"numtype,attr"`
	Pb          string `xml:"pb,attr"`
	Pdtype      string `xml:"pdtype,attr"`
	Per         string `xml:"per,attr"`
	Persoon     string `xml:"persoon,attr"`
	Pos         string `xml:"pos,attr"`
	Positie     string `xml:"positie,attr"`
	Postag      string `xml:"postag,attr"`
	Pt          string `xml:"pt,attr"`
	Pvagr       string `xml:"pvagr,attr"`
	Pvtijd      string `xml:"pvtijd,attr"`
	Refl        string `xml:"refl,attr"`
	Rel         string `xml:"rel,attr"`
	Root        string `xml:"root,attr"`
	Sc          string `xml:"sc,attr"`
	Sense       string `xml:"sense,attr"`
	Special     string `xml:"special,attr"`
	Spectype    string `xml:"spectype,attr"`
	Status      string `xml:"status,attr"`
	Tense       string `xml:"tense,attr"`
	Vform       string `xml:"vform,attr"`
	Vwtype      string `xml:"vwtype,attr"`
	Vztype      string `xml:"vztype,attr"`
	Wh          string `xml:"wh,attr"`
	Wk          string `xml:"wk,attr"`
	Word        string `xml:"word,attr"`
	Wvorm       string `xml:"wvorm"`
}

func getAttr(attr string, n *FullNode) string {
	switch attr {
	case "aform":
		return n.Aform
	case "begin":
		return fmt.Sprint(n.Begin)
	case "buiging":
		return n.Buiging
	case "case":
		return n.Case
	case "cat":
		return n.Cat
	case "comparative":
		return n.Comparative
	case "conjtype":
		return n.Conjtype
	case "def":
		return n.Def
	case "dial":
		return n.Dial
	case "end":
		return fmt.Sprint(n.End)
	case "frame":
		return n.Frame
	case "gen":
		return n.Gen
	case "genus":
		return n.Genus
	case "getal":
		return n.Getal
	case "getal-n":
		return n.GetalN
	case "graad":
		return n.Graad
	case "id":
		return n.Id
	case "index":
		return n.Index
	case "infl":
		return n.Infl
	case "lcat":
		return n.Lcat
	case "lemma":
		return n.Lemma
	case "lwtype":
		return n.Lwtype
	case "mwu_root":
		return n.MwuRoot
	case "mwu_sense":
		return n.MwuSense
	case "naamval":
		return n.Naamval
	case "neclass":
		return n.Neclass
	case "npagr":
		return n.Npagr
	case "ntype":
		return n.Ntype
	case "num":
		return n.Num
	case "numtype":
		return n.Numtype
	case "pb":
		return n.Pb
	case "pdtype":
		return n.Pdtype
	case "per":
		return n.Per
	case "persoon":
		return n.Persoon
	case "pos":
		return n.Pos
	case "positie":
		return n.Positie
	case "postag":
		return n.Postag
	case "pt":
		return n.Pt
	case "pvagr":
		return n.Pvagr
	case "pvtijd":
		return n.Pvtijd
	case "refl":
		return n.Refl
	case "rel":
		return n.Rel
	case "root":
		return n.Root
	case "sc":
		return n.Sc
	case "sense":
		return n.Sense
	case "special":
		return n.Special
	case "spectype":
		return n.Spectype
	case "status":
		return n.Status
	case "tense":
		return n.Tense
	case "vform":
		return n.Vform
	case "vwtype":
		return n.Vwtype
	case "vztype":
		return n.Vztype
	case "wh":
		return n.Wh
	case "wk":
		return n.Wk
	case "word":
		return n.Word
	case "wvorm":
		return n.Wvorm
	}
	return ""
}

func getDeepAttr(attr string, n *Node, values *[]*ValueItem) {
	if n.Index != "" && len(n.NodeList) == 0 && n.Word == "" {
		*values = append(*values, &ValueItem{-1, -1, n.Index})
		return
	}
	if s := getAttr(attr, &n.FullNode); s != "" {
		*values = append(*values, &ValueItem{n.Begin, n.End, s})
		return
	}
	for _, n2 := range n.NodeList {
		getDeepAttr(attr, n2, values)
	}
}

func getIndexValue(idx, attr string, n *Node, values *[]*ValueItem) bool {
	if n.Index == idx && (len(n.NodeList) > 0 || n.Word != "") {
		getDeepAttr(attr, n, values)
		return true
	}
	for _, n2 := range n.NodeList {
		if getIndexValue(idx, attr, n2, values) {
			return true
		}
	}
	return false
}

func getFullAttr(attr string, n, top *Node) string {

	if s := getAttr(attr, &n.FullNode); s != "" {
		return s
	}
	values := make([]*ValueItem, 0)
	getDeepAttr(attr, n, &values)

	if len(values) == 0 {
		return ""
	}

	for _, v := range values {
		if v.begin < 0 {
			getIndexValue(v.value, attr, top, &values)
		}
	}

	sort.Sort(ValueItems(values))

	s := make([]string, 0, len(values))
	p := -1
	q := 0
	for _, v := range values {
		if v.begin > p {
			p = v.begin
			if v.begin > q && len(s) > 0 {
				s = append(s, "...")
			}
			q = v.end
			s = append(s, v.value)
		}
	}
	return "  " + strings.Join(s, " ")
}

func updateText(q *Context, s string) {
	fmt.Fprintf(q.w, `<script type="text/javascript">
window.parent._fn.update(%q);
</script>
`, s)
	if ff, ok := q.w.(http.Flusher); ok {
		ff.Flush()
	}
}

func updateError(q *Context, err error, is_html bool) {
	s := err.Error()
	if is_html {
		updateText(q, "Interne fout: "+html.EscapeString(s))
	} else {
		fmt.Fprintln(q.w, "Interne fout:", s)
	}
}

func xpathstats(q *Context) {

	download := false
	if first(q.r, "d") != "" {
		download = true
	}

	if download {
		contentType(q, "text/plain; charset=utf-8")
		cache(q)
	} else {
		contentType(q, "text/html; charset=utf-8")
		cache(q)
		fmt.Fprint(q.w, `<!DOCTYPE html>
<html>
<head>
<title></title>
</head>
<body">
`)
		updateText(q, `<img src="busy.gif" alt="aan het werk...">`)
	}

	prefix := getprefix(q)
	if !q.prefixes[prefix] {
		if download {
			fmt.Fprintf(q.w, "Invalid corpus: "+prefix)
		} else {
			updateText(q, "Invalid corpus: "+html.EscapeString(prefix))
			fmt.Fprintln(q.w, "</body>\n</html>")
		}
		return
	}

	var chClose <-chan bool
	if f, ok := q.w.(http.CloseNotifier); ok {
		chClose = f.CloseNotify()
	} else {
		chClose = make(<-chan bool)
	}

	attr := make([]string, 3)
	attr[0], attr[1], attr[2] = first(q.r, "attr1"), first(q.r, "attr2"), first(q.r, "attr3")
	if attr[1] == "" {
		attr[1], attr[2] = attr[2], ""
	}
	if attr[0] == "" {
		attr[0], attr[1], attr[2] = attr[1], attr[2], ""
	}
	if attr[0] == "" {
		if download {
			fmt.Fprintln(q.w, "Geen attribuut gekozen")
		} else {
			updateText(q, "Geen attribuut gekozen")
			fmt.Fprintln(q.w, "</body>\n</html>")
		}
		return
	}

	query := first(q.r, "xpath")

	if query == "" {
		if download {
			fmt.Fprintln(q.w, "Query ontbreekt")
		} else {
			updateText(q, "Query ontbreekt")
			fmt.Fprintln(q.w, "</body>\n</html>")
		}
		return
	}

	var owner string
	rows, err := q.db.Query(fmt.Sprintf("SELECT `owner` FROM `%s_info` WHERE `id` = %q", Cfg.Prefix, prefix))
	if err != nil {
		updateError(q, err, !download)
		logerr(err)
		return
	}
	for rows.Next() {
		if doErr(q, rows.Scan(&owner)) {
			rows.Close()
			return
		}
	}
	if err := rows.Err(); err != nil {
		updateError(q, err, !download)
		logerr(err)
		return
	}

	dactfiles := make([]string, 0)
	if strings.Contains(owner, "@") {
		dactfiles = append(dactfiles, path.Join(paqudir, "data", prefix, "data.dact"))
	} else {
		rows, err := q.db.Query(fmt.Sprintf("SELECT `arch` FROM `%s_c_%s_arch` ORDER BY `id`", Cfg.Prefix, prefix))
		if err != nil {
			updateError(q, err, !download)
			logerr(err)
			return
		}
		for rows.Next() {
			var s string
			err := rows.Scan(&s)
			if err != nil {
				updateError(q, err, !download)
				logerr(err)
				rows.Close()
				return
			}
			if strings.HasSuffix(s, ".dact") {
				dactfiles = append(dactfiles, s)
			}
		}
		if err := rows.Err(); err != nil {
			updateError(q, err, !download)
			logerr(err)
			return
		}
	}

	if len(dactfiles) == 0 {
		fmt.Fprintln(q.w, "Er zijn geen dact-bestanden voor dit corpus")
		return
	}

	now := time.Now()
	now2 := time.Now()

	sums := make(map[string]int)
	count := 0
	tooMany := false
	for _, dactfile := range dactfiles {
		if !download && time.Now().Sub(now2) > 2*time.Second {
			xpathout(q, sums, attr, count, tooMany, now, download, false)
			now2 = time.Now()
		}
		if tooMany {
			break
		}
		select {
		case <-chClose:
			logerr(errConnectionClosed)
			return
		default:
		}
		db, err := dbxml.Open(dactfile)
		if err != nil {
			updateError(q, err, !download)
			logerr(err)
			return
		}

		qu, err := db.Prepare(query)
		if err != nil {
			updateError(q, err, !download)
			db.Close()
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
			updateError(q, err, !download)
			db.Close()
			return
		}
		for docs.Next() {
			if !download && time.Now().Sub(now2) > 2*time.Second {
				xpathout(q, sums, attr, count, tooMany, now, download, false)
				now2 = time.Now()
			}
			alpino := Alpino_ds{}
			alp := Alpino_ds{}
			err := xml.Unmarshal([]byte(docs.Content()), &alpino)
			if err != nil {
				updateError(q, err, !download)
				logerr(err)
				docs.Close()
				db.Close()
				return
			}
			err = xml.Unmarshal([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<alpino_ds version="1.3">
`+docs.Match()+`
</alpino_ds>`), &alp)
			if err != nil {
				updateError(q, err, !download)
				logerr(err)
				docs.Close()
				db.Close()
				return
			}
			sums[getFullAttr(attr[0], alp.Node0, alpino.Node0)+"\t"+
				getFullAttr(attr[1], alp.Node0, alpino.Node0)+"\t"+
				getFullAttr(attr[2], alp.Node0, alpino.Node0)]++
			count++
			if len(sums) >= 100000 {
				tooMany = true
				docs.Close()
			}
		}
		if err := docs.Error(); err != nil {
			logerr(err)
		}
		db.Close()
		done <- true
		select {
		case <-interrupted:
			return
		default:
		}
	}

	xpathout(q, sums, attr, count, tooMany, now, download, true)

	if !download {
		fmt.Fprintln(q.w, "</body>\n</html>")
	}

}

func xpathout(q *Context, sums map[string]int, attr []string, count int, tooMany bool, now time.Time, download bool, final bool) {

	attrList := make([]*AttrItem, 0, len(sums))
	for key, value := range sums {
		attrList = append(attrList, &AttrItem{value, key})
	}
	sort.Sort(AttrItems(attrList))

	var buf bytes.Buffer

	if download {
		if tooMany {
			fmt.Fprintln(q.w, "# ONDERBROKEN VANWEGE TE VEEL COMBINATIES")
		}
		fmt.Fprintf(q.w, "# %d matches in %d combinaties\n", count, len(attrList))
	} else {
		if tooMany {
			fmt.Fprintln(&buf, "<div class=\"warning\">Onderbroken vanwege te veel combinaties</div>")
		}
		f := ""
		if !final {
			f = `<img src="busy.gif" alt="aan het werk...">`
		}
		fmt.Fprintf(&buf, `<table>
<tr><td>Matches:<td class="right">%s<td rowspan="2" width="50">%s
<tr><td>Combinaties:<td class="right">%s
<tr><td>Tijd:<td colspan="2">%s
</table>
<p>
`, iformat(count), f, iformat(len(attrList)), time.Now().Sub(now),
		)
	}

	nAttr := 0
	for i := 0; i < 3; i++ {
		if attr[i] != "" {
			nAttr = i + 1
		}
	}

	if download {
		fmt.Fprint(q.w, "aantal")
	} else {
		fmt.Fprint(&buf, "<table class=\"breed\">\n<tr class=\"odd\"><th><th>")
	}
	for i := 0; i < nAttr; i++ {
		if download {
			fmt.Fprintf(q.w, "\t%s", attr[i])
		} else {
			fmt.Fprintf(&buf, "<th>%s", html.EscapeString(attr[i]))
		}
	}
	fmt.Fprintln(q.w)

	for n, a := range attrList {
		o := ""
		if n%2 == 1 {
			o = " class=\"odd\""
		}
		if !download && n == WRDMAX {
			fmt.Fprintf(&buf, "<tr%s><td class=\"right\">...<td class=\"right\">...\n", o)
			for i := 0; i < nAttr; i++ {
				fmt.Fprintln(&buf, "<td class=\"nil\">...")
			}
			break
		}
		if download {
			fmt.Fprintf(q.w, "%d\t%.1f%%", a.n, float64(a.n)/float64(count)*100)
		} else {
			fmt.Fprintf(&buf, "<tr%s><td class=\"right\">%d<td class=\"right\">%.1f%%\n", o, a.n, float64(a.n)/float64(count)*100)
		}
		v := strings.Split(a.a, "\t")
		for i := 0; i < nAttr; i++ {
			if strings.TrimSpace(v[i]) == "" {
				if download {
					fmt.Fprintf(q.w, "\tNIL")
				} else {
					fmt.Fprintln(&buf, "<td class=\"nil\">(leeg)")
				}
			} else {
				if download {
					fmt.Fprintf(q.w, "\t%s", strings.TrimSpace(v[i]))
				} else {
					if strings.HasPrefix(v[i], "  ") {
						fmt.Fprintf(&buf, "<td class=\"multi\">%s\n", html.EscapeString(strings.TrimSpace(v[i])))
					} else {
						fmt.Fprintf(&buf, "<td>%s\n", html.EscapeString(v[i]))
					}
				}
			}
		}
		if download {
			fmt.Fprintln(q.w)
		}
	}

	if !download {
		fmt.Fprintf(&buf,
			"</table>\n<hr><a href=\"xpathstats?%s&amp;d=1\" target=\"_blank\">download</a>\n",
			strings.Replace(q.r.URL.RawQuery, "&", "&amp;", -1))
		updateText(q, buf.String())
	}
}

func (x AttrItems) Less(i, j int) bool {
	if x[i].n != x[j].n {
		return x[i].n > x[j].n
	}
	return x[i].a < x[j].a
}

func (x AttrItems) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

func (x AttrItems) Len() int {
	return len(x)
}

func (x ValueItems) Less(i, j int) bool {
	if x[i].begin < x[j].begin {
		return true
	}
	return x[i].end < x[j].end
}

func (x ValueItems) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

func (x ValueItems) Len() int {
	return len(x)
}
