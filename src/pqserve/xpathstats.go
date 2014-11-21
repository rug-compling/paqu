package main

import (
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
)

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
	Wvorm       string `xml:"wvorm,attr"`
	other       map[string]string
}

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
	if n.other != nil {
		return n.other[attr]
	}
	return ""
}

type NodeTT Node

func (x *Node) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		if n := attr.Name.Local; !keyTags[n] {
			if x.other == nil {
				x.other = make(map[string]string)
			}
			x.other[n] = attr.Value
		}
	}
	return d.DecodeElement((*NodeTT)(x), &start)
}

func updateText(q *Context, s string) {
	fmt.Fprintf(q.w, `<script type="text/javascript">
f(%q);
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

func init() {
	for _, tag := range NodeTags {
		keyTags[tag] = true
	}
}
