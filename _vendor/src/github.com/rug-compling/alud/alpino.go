package alud

import (
	"encoding/xml"
	"regexp"
	"strconv"
	"strings"
)

var (
	reShorted  = regexp.MustCompile(`></(meta|parser|node|dep|acl|advcl|advmod|amod|appos|aux|case|cc|ccomp|clf|compound|conj|cop|csubj|det|discourse|dislocated|expl|fixed|flat|goeswith|iobj|list|mark|nmod|nsubj|nummod|obj|obl|orphan|parataxis|punct|ref|reparandum|root|vocative|xcomp)>`)
	reNoConllu = regexp.MustCompile(`><!\[CDATA\[\s*\]\]></conllu>`)
)

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
		}
	}
	for _, node := range q.allnodes {
		node.Ud = &udType{Dep: make([]depType, 0)}
		node.Begin /= 1000
		node.End /= 1000
	}
	q.alpino.UdNodes = []*udNodeType{}
	q.alpino.Conllu = &conlluType{Auto: versionID}
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

	return s
}

func alpinoDo(conllu string, q *context) {

	alpino := q.alpino

	lines := make([]string, 0)

	for _, line := range strings.Split(conllu, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}

	udNodeList := make([]*udNodeType, 0)
	eudNodeList := make([]*udNodeType, 0)

	for _, line := range lines {
		a := strings.Split(line, "\t")
		items := getItems(a[9])
		es, isCopy := items["CopiedFrom"]
		if !isCopy {
			es = a[0]
		}
		e, _ := strconv.Atoi(es)
		node := getNode(alpino.Node, e)

		if a[8] != "_" {
			for _, deps := range strings.Split(a[8], "|") {
				dep := strings.SplitN(deps, ":", 2)
				var aux string
				dd := strings.SplitN(dep[1], ":", 2)
				if len(dd) > 1 {
					aux = dd[1]
				}
				node.Ud.Dep = append(node.Ud.Dep, depType{
					Id:         a[0],
					Head:       dep[0],
					Deprel:     dep[1],
					DeprelMain: dd[0],
					DeprelAux:  aux,
					Elided:     strings.Contains(a[0], "."),
				})
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

			featsType: node.Ud.featsType,

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
		}
		udNodeList = append(udNodeList, &ud)

		for _, dep := range node.Ud.Dep {

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

				featsType: node.Ud.featsType,

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
			}
			eudNodeList = append(eudNodeList, &ud)
		}

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
	q.alpino.Conllu.Status = "OK"
	q.alpino.Conllu.Conllu = "\n" + strings.TrimSpace(conllu) + "\n"
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
	if node == nil {
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
