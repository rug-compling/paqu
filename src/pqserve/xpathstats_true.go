// +build !nodbxml

package main

import (
	"github.com/pebbe/dbxml"

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

type Alpino_ds_full_node struct {
	XMLName xml.Name  `xml:"alpino_ds"`
	Node0   *FullNode `xml:"node"`
}

type FullNode struct {
	Aform       string `xml:"aform,attr"`
	Begin       string `xml:"begin,attr"`
	Buiging     string `xml:"buiging,attr"`
	Case        string `xml:"case,attr"`
	Cat         string `xml:"cat,attr"`
	Comparative string `xml:"comparative,attr"`
	Conjtype    string `xml:"conjtype,attr"`
	Def         string `xml:"def,attr"`
	Dial        string `xml:"dial,attr"`
	End         string `xml:"end,attr"`
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
		return n.Begin
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
		return n.End
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

func xpathstats(q *Context) {

	download := false
	if first(q.r, "d") != "" {
		download = true
	}

	if download {
		q.w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	} else {
		q.w.Header().Set("Content-Type", "text/html; charset=utf-8")
	}
	q.w.Header().Set("Cache-Control", "no-cache")
	q.w.Header().Add("Pragma", "no-cache")

	prefix := getprefix(q)
	if !q.prefixes[prefix] {
		fmt.Fprintf(q.w, "Invalid corpus: "+html.EscapeString(prefix))
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
		fmt.Fprintln(q.w, "Geen attribuut gekozen")
		return
	}
	fmt.Println(attr)

	query := first(q.r, "xpath")

	if query == "" {
		fmt.Fprintln(q.w, "Query ontbreekt")
		return
	}

	var owner string
	rows, err := q.db.Query(fmt.Sprintf("SELECT `owner` FROM `%s_info` WHERE `id` = %q", Cfg.Prefix, prefix))
	if err != nil {
		interneFoutRegel(q, err, !download)
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
		interneFoutRegel(q, err, !download)
		logerr(err)
		return
	}

	dactfiles := make([]string, 0)
	if strings.Contains(owner, "@") {
		dactfiles = append(dactfiles, path.Join(paqudir, "data", prefix, "data.dact"))
	} else {
		rows, err := q.db.Query(fmt.Sprintf("SELECT `arch` FROM `%s_c_%s_arch` ORDER BY `id`", Cfg.Prefix, prefix))
		if err != nil {
			interneFoutRegel(q, err, !download)
			logerr(err)
			return
		}
		for rows.Next() {
			var s string
			err := rows.Scan(&s)
			if err != nil {
				interneFoutRegel(q, err, !download)
				logerr(err)
				rows.Close()
				return
			}
			if strings.HasSuffix(s, ".dact") {
				dactfiles = append(dactfiles, s)
			}
		}
		if err := rows.Err(); err != nil {
			interneFoutRegel(q, err, !download)
			logerr(err)
			return
		}
	}

	if len(dactfiles) == 0 {
		fmt.Fprintln(q.w, "Er zijn geen dact-bestanden voor dit corpus")
		return
	}

	now := time.Now()

	sums := make(map[string]int)

	for _, dactfile := range dactfiles {
		select {
		case <-chClose:
			logerr(errConnectionClosed)
			return
		default:
		}
		db, err := dbxml.Open(dactfile)
		if err != nil {
			interneFoutRegel(q, err, !download)
			logerr(err)
			return
		}
		docs, err := db.Query(query)
		if err != nil {
			interneFoutRegel(q, err, !download)
			return
		}
		for docs.Next() {
			select {
			case <-chClose:
				docs.Close()
				db.Close()
				logerr(errConnectionClosed)
				return
			default:
			}
			alp := Alpino_ds_full_node{}
			err := xml.Unmarshal([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<alpino_ds version="1.3">
`+docs.Match()+`
</alpino_ds>`), &alp)
			if err != nil {
				interneFoutRegel(q, err, !download)
				logerr(err)
				return
			}
			sums[getAttr(attr[0], alp.Node0)+"\t"+getAttr(attr[1], alp.Node0)+"\t"+getAttr(attr[2], alp.Node0)]++
		}
		db.Close()
	}

	attrList := make([]*AttrItem, 0, len(sums))
	for key, value := range sums {
		attrList = append(attrList, &AttrItem{value, key})
	}
	sort.Sort(AttrItems(attrList))

	if download {
		fmt.Fprintf(q.w, "# %d combinaties\t\n", len(attrList))
	} else {
		fmt.Fprintln(q.w, "Aantal gevonden combinaties:", iformat(len(attrList)))
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
		fmt.Fprint(q.w, "<table>\n<tr><th>")
	}
	for i := 0; i < nAttr; i++ {
		if download {
			fmt.Fprintf(q.w, "\t%s", attr[i])
		} else {
			fmt.Fprintf(q.w, "<th>%s", html.EscapeString(attr[i]))
		}
	}
	fmt.Fprintln(q.w)

	for n, a := range attrList {
		if !download && n == WRDMAX {
			fmt.Fprintln(q.w, "<tr><td class=\"right\">...")
			for i := 0; i < nAttr; i++ {
				fmt.Fprintln(q.w, "<td class=\"nil\">...")
			}
			break
		}
		if download {
			fmt.Fprint(q.w, a.n)
		} else {
			fmt.Fprintf(q.w, "<tr><td class=\"right\">%d\n", a.n)
		}
		v := strings.Split(a.a, "\t")
		for i := 0; i < nAttr; i++ {
			if v[i] == "" {
				if download {
					fmt.Fprintf(q.w, "\tNIL")
				} else {
					fmt.Fprintln(q.w, "<td><span class=\"nil\">&mdash;</span>")
				}
			} else {
				if download {
					fmt.Fprintf(q.w, "\t%s", v[i])
				} else {
					fmt.Fprintf(q.w, "<td>%s\n", html.EscapeString(v[i]))
				}
			}
		}
		if download {
			fmt.Fprintln(q.w)
		}
	}

	if !download {
		fmt.Fprintf(q.w,
			"</table>\n<hr>tijd: %s\n<p>\n<a href=\"xpathstats?%s&amp;d=1\" target=\"_blank\">download</a>\n",
			time.Now().Sub(now),
			strings.Replace(q.r.URL.RawQuery, "&", "&amp;", -1))
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
