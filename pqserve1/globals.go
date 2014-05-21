package main

//. Imports

import (
	"bytes"
	"encoding/xml"
	"regexp"
)

//. Types

// een dependency relation, geretourneerd door SQL
type Row struct {
	word    string
	lemma   string
	postag  string
	rel     string
	hpostag string
	hlemma  string
	hword   string
	begin   int
	end     int
	hbegin  int
	hend    int
	mark    string
}

type Sentence struct {
	arch  int
	file  int
	words []string // de zin opgesplitst in woorden
	items []Row    // alle matchende dependency relations voor deze zin
}

type Alpino_ds struct {
	XMLName  xml.Name `xml:"alpino_ds"`
	Node0    *Node    `xml:"node"`
	Sentence string   `xml:"sentence"`
}

type Node struct {
	Id       string  `xml:"id,attr"`
	Index    string  `xml:"index,attr"`
	Cat      string  `xml:"cat,attr"`
	Pt       string  `xml:"pt,attr"`
	Word     string  `xml:"word,attr"`
	Lemma    string  `xml:"lemma,attr"`
	Postag   string  `xml:"postag,attr"`
	Rel      string  `xml:"rel,attr"`
	Begin    int     `xml:"begin,attr"`
	End      int     `xml:"end,attr"`
	NodeList []*Node `xml:"node"`
	skip     bool
}

type TreeContext struct {
	yellow map[int]bool
	green  map[int]bool
	marks  map[string]bool
	refs   map[string]bool
	graph  bytes.Buffer // definitie dot-bestand
	start  int
	words  []string
}

type Corpus struct {
	Id    string
	Title string
}

type Corpora []*Corpus

//. Constanten

const (
	ZINMAX = 10
	WRDMAX = 250
	YELLOW = "<span style=\"background-color: yellow\">"
	GREEN  = "<span style=\"background-color: lightgreen\">"
	YELGRN = "<span style=\"background-color: cyan\">"

	PRE = "wordrel"
)

//. Variabelen

var (
	login string

	reQuote = regexp.MustCompile("%[0-9A-F-a-f][0-9A-F-a-f]") // voor urldecode()

	opt_postag  = make([]string, 0)
	opt_hpostag = make([]string, 0)
	opt_rel     = make([]string, 0)

	opt_db   = []string{}
	opt_dbc  = []*Corpus{}
	prefixes = map[string]bool{}
)
