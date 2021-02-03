package alud

import (
	"github.com/rug-compling/alpinods"

	"encoding/xml"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	reShorted  = regexp.MustCompile(`></(meta|parser|node|data|dep|acl|advcl|advmod|amod|appos|aux|case|cc|ccomp|clf|compound|conj|cop|csubj|det|discourse|dislocated|expl|fixed|flat|goeswith|iobj|list|mark|nmod|nsubj|nummod|obj|obl|orphan|parataxis|punct|ref|reparandum|root|vocative|xcomp)>`)
	reNoConllu = regexp.MustCompile(`><!\[CDATA\[\s*\]\]></conllu>`)

	reEnts = regexp.MustCompile("&#(34|38|39|60|62);")
	ents   = map[string]string{
		"&#34;": "&quot;",
		"&#38;": "&amp;",
		"&#39;": "&apos;",
		"&#60;": "&lt;",
		"&#62;": "&gt;",
	}
)

// Insert given Universal Dependencies into alpino_ds format.
//
// Use UD info from alpino_doc if conllu is "".
//
// The conllu format is not checked for correctness. Garbage in, garbage out.
//
// The value from auto is copied to the output.
func Alpino(alpino_doc []byte, conllu, auto string) (alpino string, err error) {

	defer func() {
		if r := recover(); r != nil {
			alpino = ""
			err = fmt.Errorf("%v", r)
		}
	}()

	var alp alpino_ds
	if err = xml.Unmarshal(alpino_doc, &alp); err != nil {
		return "", err
	}
	if strings.TrimSpace(conllu) == "" && alp.Conllu != nil {
		conllu = alp.Conllu.Conllu
	}
	var reset func(*nodeType)
	reset = func(node *nodeType) {
		node.Ud = &udType{Dep: make([]depType, 0)}
		if node.Node != nil {
			for _, n := range node.Node {
				reset(n)
			}
		}
	}
	reset(alp.Node)
	alp.UdNodes = []*udNodeType{}
	alp.Conllu = &conlluType{Conllu: conllu, Auto: auto}
	alpinoDo(conllu, &alp)
	alpino, err = alpinoFormat(&alp), nil
	return // geen argumenten i.v.m. recover
}

// Derive Universal Dependencies and insert into alpino_ds format.
//
// If sentid is "" it is derived from the filename.
//
// When err is not nil and alpino is not "" it contains the err in the alpino_ds format.
func UdAlpino(alpino_doc []byte, filename, sentid string) (alpino string, err error) {
	conllu, q, err := ud(alpino_doc, filename, sentid, OPT_NO_COMMENTS|OPT_NO_DETOKENIZE)

	if err == nil {
		alpinoRestore(q)
		alpinoDo(conllu, q.alpino)
		return alpinoFormat(q.alpino), nil
	}

	var alp alpino_ds
	if xml.Unmarshal(alpino_doc, &alp) != nil {
		return "", err
	}

	e := err.Error()
	i := strings.Index(e, "\n")
	if i > 0 {
		e = e[:i]
	}

	var r func(*nodeType)
	r = func(node *nodeType) {
		node.Ud = nil
		for _, n := range node.Node {
			r(n)
		}
	}
	if alp.Node != nil {
		r(alp.Node)
	}

	if alp.Sentence.SentId == "" {
		id := filepath.Base(filename)
		if strings.HasSuffix(id, ".xml") {
			id = id[:len(id)-4]
		}
		alp.Sentence.SentId = id
	}
	alp.UdNodes = []*udNodeType{}
	alp.Conllu = &conlluType{
		Status: "error",
		Error:  e,
		Auto:   VersionID(),
		Conllu: " ", // spatie is nodig, wordt later verwijderd
	}
	return alpinoFormat(&alp), err
}

func alpinoRestore(q *context) {
	for i := len(q.swapped) - 1; i >= 0; i-- {
		swap(q.swapped[i][1], q.swapped[i][0])
	}
	for _, n := range q.varindexnodes {
		node := n.(*nodeType)
		if node.udOldState != nil {
			node.Begin = node.udOldState.Begin
			node.End = node.udOldState.End
			node.Word = node.udOldState.Word
			node.Lemma = node.udOldState.Lemma
			node.Postag = node.udOldState.Postag
			node.Pt = node.udOldState.Pt
			node.Cat = node.udOldState.Cat
			node.Node = node.udOldState.Node
		}
	}
	for _, node := range q.allnodes {
		node.Ud = &udType{Dep: make([]depType, 0)}
		node.Begin /= 1000
		node.End /= 1000
		node.ID /= 1000
	}
	q.alpino.UdNodes = []*udNodeType{}
	q.alpino.Conllu = &conlluType{Auto: VersionID()}
}

func alpinoFormat(alpino *alpino_ds) string {
	var v1, v2 int
	if a := strings.Split(alpino.Version, "."); len(a) > 1 {
		var err error
		if v1, err = strconv.Atoi(a[0]); err != nil {
			v1 = 0
		}
		if v2, err = strconv.Atoi(a[1]); err != nil {
			v2 = 0
		}
	}
	if v1 < 1 || (v1 == 1 && v2 < 10) {
		alpino.Version = "1.10"
	}

	b, _ := xml.MarshalIndent(alpino, "", "  ")
	s := "<?xml version=\"1.0\"?>\n" + string(b)

	// shorten
	s = reShorted.ReplaceAllString(s, "/>")
	s = reNoConllu.ReplaceAllString(s, "/>")

	// standard XML entities
	s = reEnts.ReplaceAllStringFunc(s, func(s1 string) string {
		return ents[s1]
	})

	return s
}

func alpinoDo(conllu string, alpino *alpino_ds) {

	lines := make([]string, 0)

	for _, line := range strings.Split(conllu, "\n") {
		line = strings.TrimSpace(line)
		if line != "" && line[0] != '#' {
			lines = append(lines, line)
		}
	}

	udNodeList := make([]*udNodeType, 0)
	eudNodeList := make([]*udNodeType, 0)
	nodeDepList := make([]nodeDepType, 0)

	for _, line := range lines {
		a := strings.Split(line, "\t")
		if strings.Contains(a[0], "-") {
			continue
		}
		items := getItems(a[9])
		es, isCopy := items["CopiedFrom"]
		if !isCopy {
			if strings.Contains(a[0], ".") {
				panic("Missing CopiedFrom for ID=" + a[0])
			}
			es = a[0]
		}
		e, err := strconv.Atoi(es)
		if err != nil {
			panic(fmt.Sprintf("%v for ID=%s", err, a[0]))
		}
		node := getNode(alpino.Node, e)

		if a[8] != "_" {
			for _, deps := range strings.Split(a[8], "|") {
				dep := strings.SplitN(deps, ":", 2)
				var aux string
				dd := strings.SplitN(dep[1], ":", 2)
				if len(dd) > 1 {
					aux = dd[1]
				}
				d := depType{
					Id:         a[0],
					Head:       dep[0],
					Deprel:     dep[1],
					DeprelMain: dd[0],
					DeprelAux:  aux,
					Elided:     strings.Contains(a[0], "."),
				}
				node.Ud.Dep = append(node.Ud.Dep, d)
				nodeDepList = append(nodeDepList, nodeDepType{node: node, dep: &d})
			}
		}

		if isCopy {
			continue
		}

		node.Ud.Id = a[0]
		node.Ud.Form = a[1]
		node.Ud.Lemma = a[2]
		node.Ud.Upos = noe(a[3])
		node.Ud.Head = noe(a[6])
		node.Ud.Deprel = noe(a[7])

		dd := strings.SplitN(node.Ud.Deprel, ":", 2)
		node.Ud.DeprelMain = dd[0]
		if len(dd) > 1 {
			node.Ud.DeprelAux = dd[1]
		}

		feats := getItems(a[5])
		node.Ud.Abbr = feats["Abbr"]
		node.Ud.Case = feats["Case"]
		node.Ud.Definite = feats["Definite"]
		node.Ud.Degree = feats["Degree"]
		node.Ud.Foreign = feats["Foreign"]
		node.Ud.Gender = feats["Gender"]
		node.Ud.Number = feats["Number"]
		node.Ud.Person = feats["Person"]
		node.Ud.PronType = feats["PronType"]
		node.Ud.Reflex = feats["Reflex"]
		node.Ud.Tense = feats["Tense"]
		node.Ud.VerbForm = feats["VerbForm"]

		ud := udNodeType{
			recursion: make([]string, 0),

			XMLName:   xml.Name{Local: node.Ud.DeprelMain},
			Id:        node.Ud.Id,
			Form:      node.Ud.Form,
			Lemma:     node.Ud.Lemma,
			Upos:      node.Ud.Upos,
			Head:      node.Ud.Head,
			Deprel:    node.Ud.Deprel,
			DeprelAux: node.Ud.DeprelAux,
			Ud:        "basic",

			Feats: node.Ud.Feats,

			DeprelAttributes: alpinods.DeprelAttributes{
				Buiging:  node.Buiging,
				Conjtype: node.Conjtype,
				Dial:     node.Dial,
				Genus:    node.Genus,
				Getal:    node.Getal,
				GetalN:   node.GetalN,
				Graad:    node.Graad,
				Lwtype:   node.Lwtype,
				Naamval:  node.Naamval,
				Npagr:    node.Npagr,
				Ntype:    node.Ntype,
				Numtype:  node.Numtype,
				Pdtype:   node.Pdtype,
				Persoon:  node.Persoon,
				Positie:  node.Positie,
				Pt:       node.Pt,
				Pvagr:    node.Pvagr,
				Pvtijd:   node.Pvtijd,
				Spectype: node.Spectype,
				Status:   node.Status,
				Vwtype:   node.Vwtype,
				Vztype:   node.Vztype,
				Wvorm:    node.Wvorm,
			},
		}
		udNodeList = append(udNodeList, &ud)
	}

	for _, nd := range nodeDepList {

		node := nd.node
		dep := nd.dep

		ud := udNodeType{
			recursion: make([]string, 0),

			XMLName:   xml.Name{Local: dep.DeprelMain},
			Id:        dep.Id,
			Form:      node.Ud.Form,
			Lemma:     node.Ud.Lemma,
			Upos:      node.Ud.Upos,
			Head:      dep.Head,
			Deprel:    dep.Deprel,
			DeprelAux: dep.DeprelAux,
			Ud:        "enhanced",

			Feats: node.Ud.Feats,

			DeprelAttributes: alpinods.DeprelAttributes{
				Buiging:  node.Buiging,
				Conjtype: node.Conjtype,
				Dial:     node.Dial,
				Genus:    node.Genus,
				Getal:    node.Getal,
				GetalN:   node.GetalN,
				Graad:    node.Graad,
				Lwtype:   node.Lwtype,
				Naamval:  node.Naamval,
				Npagr:    node.Npagr,
				Ntype:    node.Ntype,
				Numtype:  node.Numtype,
				Pdtype:   node.Pdtype,
				Persoon:  node.Persoon,
				Positie:  node.Positie,
				Pt:       node.Pt,
				Pvagr:    node.Pvagr,
				Pvtijd:   node.Pvtijd,
				Spectype: node.Spectype,
				Status:   node.Status,
				Vwtype:   node.Vwtype,
				Vztype:   node.Vztype,
				Wvorm:    node.Wvorm,
			},
		}
		eudNodeList = append(eudNodeList, &ud)
	}

	alpino.UdNodes = make([]*udNodeType, 0)

	for _, n := range udNodeList {
		if n.Head == "0" {
			alpino.UdNodes = append(alpino.UdNodes, n)
		}
	}

	for _, n := range eudNodeList {
		if n.Head == "0" {
			alpino.UdNodes = append(alpino.UdNodes, n)
		}
	}

	for i, root := range alpino.UdNodes {
		var items []*udNodeType
		if i == 0 {
			items = udNodeList
		} else {
			items = eudNodeList
		}
		root.UdNodes = make([]*udNodeType, 0)
		expand(root, items)
	}
	minify(alpino)
	alpino.Conllu.Status = "OK"
	alpino.Conllu.Conllu = "\n" + strings.TrimSpace(conllu) + "\n"
}

/*
  bedoeling:
  - zorg ervoor dat ",omitempty" werkt
*/
func minify(alpino *alpino_ds) {
	if alpino.Metadata != nil && (alpino.Metadata.Meta == nil || len(alpino.Metadata.Meta) == 0) {
		alpino.Metadata = nil
	}
	if alpino.Parser != nil && alpino.Parser.Build == "" && alpino.Parser.Date == "" && alpino.Parser.Cats == "" && alpino.Parser.Skips == "" {
		alpino.Parser = nil
	}
	if alpino.Sentence != nil && alpino.Sentence.Sent == "" && alpino.Sentence.SentId == "" {
		alpino.Sentence = nil
	}
	if alpino.Comments != nil && (alpino.Comments.Comment == nil || len(alpino.Comments.Comment) == 0) {
		alpino.Comments = nil
	}
	minifyNode(alpino.Node)
}

func minifyNode(node *nodeType) {
	if node == nil {
		return
	}
	if node.Ud != nil {
		if node.Ud.Id == "" {
			node.Ud = nil
		} else {
			if len(node.Ud.Dep) == 0 {
				node.Ud.Dep = nil
			}
		}
	}
	if node.Node != nil {
		for _, n := range node.Node {
			minifyNode(n)
		}
	}
}

func noe(s string) string {
	if s == "_" {
		return ""
	}
	return s
}

func getNode(node *nodeType, end int) *nodeType {
	if node == nil || node.End < end || node.Begin >= end {
		return nil
	}
	if node.End == end && node.Word != "" {
		return node
	}
	if node.Node != nil {
		for _, n := range node.Node {
			if n2 := getNode(n, end); n2 != nil {
				return n2
			}
		}
	}
	return nil
}

func getItems(s string) map[string]string {
	m := make(map[string]string)
	if s == "_" {
		return m
	}
	for _, item := range strings.Split(s, "|") {
		a := strings.SplitN(item, "=", 2)
		if len(a) == 2 {
			m[strings.TrimSpace(a[0])] = strings.TrimSpace(a[1])
		}
	}
	return m
}

func expand(udnode *udNodeType, items []*udNodeType) {
	for _, item := range items {
		if item.Head == udnode.Id {
			it := new(udNodeType)
			*it = *item
			it.UdNodes = make([]*udNodeType, 0)
			it.recursion = append([]string{udnode.Id}, udnode.recursion...)
			udnode.UdNodes = append(udnode.UdNodes, it)
		}
	}
	for _, un := range udnode.UdNodes {
		if recursionLimit(un.recursion) {
			un.RecursionLimit = "TOO DEEP"
		} else {
			expand(un, items)
		}
	}
}

func recursionLimit(s []string) bool {
	if len(s) < 2 {
		return false
	}
	found := 0
	for i := 1; i < len(s)-1; i++ {
		if s[i] == s[0] && s[i+1] == s[1] {
			found++
		}
	}
	return found > 1
}
