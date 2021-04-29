package main

import (
	"github.com/pebbe/util"
	"github.com/rug-compling/alpinods"

	"crypto/md5"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	noNode = &nodeType{
		NodeAttributes: alpinods.NodeAttributes{
			Begin: -1,
			End:   -1,
			ID:    -1,
		},
		Node:                []*nodeType{},
		axParent:            []interface{}{},
		axAncestors:         []interface{}{},
		axAncestorsOrSelf:   []interface{}{},
		axChildren:          []interface{}{},
		axDescendants:       []interface{}{},
		axDescendantsOrSelf: []interface{}{},
	}

	posCount    = make(map[string]int)
	postagCount = make(map[string]int)
	ptCount     = make(map[string]int)
	types       = make(map[string]bool)

	sentenceCount = 0
	lengths       = make(map[string]map[int]int)
	sentences     = make(map[string]int)
	items         = make(map[string]int)

	has_his      bool
	has_parser   bool
	has_pos_verb bool
	has_qm       bool
	has_sc       bool
	has_yn       bool

	first = true

	tokens   = 0
	tokenlen = 0

	macroRE  = regexp.MustCompile(`([a-zA-Z][_a-zA-Z0-9]*)\s*=\s*"""((?s:.*?))"""`)
	macroKY  = regexp.MustCompile(`%[a-zA-Z][_a-zA-Z0-9]*%`)
	macroCOM = regexp.MustCompile(`(?m:^\s*#.*)`)

	rules = make(map[string]string)

	x = util.CheckErr
)

func do_spod(data []byte) {
	sentenceCount++

	var alpino alpino_ds
	x(xml.Unmarshal(data, &alpino))

	// Extra node bovenaan vanwege gedoe met //node
	alpino.Node = &nodeType{
		NodeAttributes: alpinods.NodeAttributes{
			Begin: alpino.Node.Begin,
			End:   alpino.Node.End,
			ID:    -2, // ??? TODO
		},
		Node: []*nodeType{alpino.Node},
	}

	var walk func(*nodeType)
	walk = func(node *nodeType) {
		if p := node.Pos; p != "" {
			posCount[p] = posCount[p] + 1
		}
		if p := node.Postag; p != "" {
			postagCount[p] = postagCount[p] + 1
		}
		if p := node.Pt; p != "" {
			ptCount[p] = ptCount[p] + 1
		} else if node.Word != "" {
			ptCount["na"] = ptCount["na"] + 1
		}
		if node.Node == nil {
			node.Node = make([]*nodeType, 0)
		} else {
			for _, n := range node.Node {
				walk(n)
			}
		}
	}
	walk(alpino.Node)

	qq := &context{
		alpino:   &alpino,
		sentence: alpino.Sentence.Sent,
		sentid:   alpino.Sentence.SentId,
		varroot:  []interface{}{alpino.Node},
	}

	inspect(qq)

	for _, node := range qq.ptnodes {
		if node.Pt != "let" {
			types[node.Word] = true
			tokens++
			tokenlen += utf8.RuneCountInString(strings.Replace(node.Word, "ij", "y", -1))
		}
	}

	for _, spod := range spods {

		if spod.special == "hidden1" || spod.special == "attr" {
			continue
		}

		if first {
			lengths[spod.lbl] = make(map[int]int)
		}

		if spod.special == "parser" {
			if alpino.Parser == nil {
				continue
			}
			p := alpino.Parser
			cats, err := strconv.Atoi(p.Cats)
			if err != nil {
				cats = -1
			}
			skips, err := strconv.Atoi(p.Skips)
			if err != nil {
				skips = -1
			}
			switch spod.lbl {
			case "ok":
				if cats == 1 && skips == 0 {
					sentences[spod.lbl] = sentences[spod.lbl] + 1
					items[spod.lbl] = items[spod.lbl] + 1
				}
			case "cats0":
				if cats == 0 {
					sentences[spod.lbl] = sentences[spod.lbl] + 1
					items[spod.lbl] = items[spod.lbl] + 1
				}
			case "cats1":
				if cats == 1 {
					sentences[spod.lbl] = sentences[spod.lbl] + 1
					items[spod.lbl] = items[spod.lbl] + 1
				}
			case "cats2":
				if cats == 2 {
					sentences[spod.lbl] = sentences[spod.lbl] + 1
					items[spod.lbl] = items[spod.lbl] + 1
				}
			case "cats3":
				if cats == 3 {
					sentences[spod.lbl] = sentences[spod.lbl] + 1
					items[spod.lbl] = items[spod.lbl] + 1
				}
			case "cats4":
				if cats > 3 {
					sentences[spod.lbl] = sentences[spod.lbl] + 1
					items[spod.lbl] = items[spod.lbl] + 1
				}
			case "skips0":
				if skips == 0 {
					sentences[spod.lbl] = sentences[spod.lbl] + 1
					items[spod.lbl] = items[spod.lbl] + 1
				}
			case "skips1":
				if skips == 1 {
					sentences[spod.lbl] = sentences[spod.lbl] + 1
					items[spod.lbl] = items[spod.lbl] + 1
				}
			case "skips2":
				if skips == 2 {
					sentences[spod.lbl] = sentences[spod.lbl] + 1
					items[spod.lbl] = items[spod.lbl] + 1
				}
			case "skips3":
				if skips == 3 {
					sentences[spod.lbl] = sentences[spod.lbl] + 1
					items[spod.lbl] = items[spod.lbl] + 1
				}
			case "skips4":
				if skips > 3 {
					sentences[spod.lbl] = sentences[spod.lbl] + 1
					items[spod.lbl] = items[spod.lbl] + 1
				}
			}
			continue
		}

		results, err := spod2xpath[spod.lbl].do(qq)
		x(err)
		if results == nil || len(results) == 0 {
			continue
		}
		sentences[spod.lbl] = sentences[spod.lbl] + 1
		items[spod.lbl] = items[spod.lbl] + len(results)
		for _, r := range results {
			node := r.(*nodeType)
			lengths[spod.lbl][node.size] = lengths[spod.lbl][node.size] + 1
		}

	}
	first = false

}

func inspect(q *context) {
	allnodes := make([]*nodeType, 0)
	varallnodes := make([]interface{}, 0)
	ptnodes := make([]*nodeType, 0)
	varindexnodes := make([]interface{}, 0)

	indextable := make(map[int]*nodeType)

	if q.alpino.Parser != nil && q.alpino.Parser.Cats != "" && q.alpino.Parser.Skips != "" {
		has_parser = true
	}

	var walk func(*nodeType)
	walk = func(node *nodeType) {

		if node.Pos == "verb" {
			has_pos_verb = true
		}
		if node.His != "" {
			has_his = true
		}
		if node.Word == "?" {
			has_qm = true
		}
		if node.Stype == "ynquestion" {
			has_yn = true
		}
		if node.Sc != "" {
			has_sc = true
		}

		if node.Index > 0 && (node.Word != "" || len(node.Node) > 0) {
			indextable[node.Index] = node
		}

		// bug in Alpino: missing pt
		if node.Word != "" && node.Pt == "" {
			node.Pt = strings.ToLower(strings.Split(node.Postag, "(")[0])
			if node.Pt == "" {
				node.Pt = "na"
			}
		}

		allnodes = append(allnodes, node)
		varallnodes = append(varallnodes, node)
		if node.Pt != "" {
			ptnodes = append(ptnodes, node)
		}
		if node.Index > 0 {
			varindexnodes = append(varindexnodes, node)
		}
		for _, n := range node.Node {
			n.parent = node
			n.axParent = []interface{}{node}
			walk(n)
		}
		node.axChildren = make([]interface{}, 0)
		node.axDescendants = make([]interface{}, 0)
		node.axDescendantsOrSelf = make([]interface{}, 1)
		node.axDescendantsOrSelf[0] = node
		for _, n := range node.Node {
			node.axChildren = append(node.axChildren, n)
			node.axDescendants = append(node.axDescendants, n)
			node.axDescendants = append(node.axDescendants, n.axDescendants...)
			node.axDescendantsOrSelf = append(node.axDescendantsOrSelf, n.axDescendantsOrSelf...) // niet n
		}
	}
	walk(q.alpino.Node)
	q.alpino.Node.parent = noNode
	q.alpino.Node.axParent = []interface{}{}

	var found map[int]bool
	walk = func(node *nodeType) {
		if node.Index > 0 {
			node = indextable[node.Index]
		}
		if node.Word != "" {
			found[node.ID] = true
		} else {
			for _, n := range node.Node {
				walk(n)
			}
		}
	}
	for _, node := range allnodes {
		found = make(map[int]bool)
		walk(node)
		node.size = len(found)
	}

	for _, node := range allnodes {
		node.axAncestors = make([]interface{}, 0)
		node.axAncestorsOrSelf = make([]interface{}, 0)
		node.axAncestorsOrSelf = append(node.axAncestorsOrSelf, node)
		if node != q.alpino.Node {
			node.axAncestors = append(node.axAncestors, node.parent)
			node.axAncestors = append(node.axAncestors, node.parent.axAncestors...)
			if node.axAncestors[len(node.axAncestors)-1] != q.alpino.Node {
				// zou niet mogelijk moeten zijn
				panic("Missing ancestors in " + q.filename)
			}
			node.axAncestorsOrSelf = append(node.axAncestorsOrSelf, node.parent.axAncestorsOrSelf...)
		}
	}

	sort.Slice(ptnodes, func(i, j int) bool {
		return ptnodes[i].End < ptnodes[j].End
	})
	varptnodes := make([]interface{}, len(ptnodes))
	for i, node := range ptnodes {
		varptnodes[i] = node
	}

	q.allnodes = allnodes
	q.varallnodes = varallnodes
	q.varindexnodes = varindexnodes
	q.ptnodes = ptnodes
	q.varptnodes = varptnodes

}

func spod_save() {

	dir := filepath.Join(paqudatadir, "data", prefix, "spod")

	x(os.MkdirAll(dir, 0o777))

	fp, err := os.Create(filepath.Join(dir, "stats"))
	x(err)
	fmt.Fprintf(fp, "%8d\t\t\tzinnen\tzinnen\n", sentenceCount)
	fmt.Fprintf(fp, "%8d\t\t\twoorden\twoorden\n", tokens)
	fmt.Fprintf(fp, "%8.4f\t\t\ttt\ttypes per token\n", float64(len(types))/float64(tokens))
	fmt.Fprintf(fp, "%8.4f\t\t\twz\twoorden per zin\n", float64(tokens)/float64(sentenceCount))
	fmt.Fprintf(fp, "%8.4f\t\t\tlw\tletters per woord\n", float64(tokenlen)/float64(tokens))
	fp.Close()

	loadMacros()

	for _, spod := range spods {
		filename := filepath.Join(dir, spod_fingerprint(spod.xpath, spod.method))
		fp, err := os.Create(filename)
		x(err)
		if spod.special == "hidden1" {
			var found bool
			switch spod.lbl {
			case "has_his":
				found = has_his
			case "has_parser":
				found = has_parser
			case "has_pos_verb":
				found = has_pos_verb
			case "has_qm":
				found = has_qm
			case "has_sc":
				found = has_sc
			case "has_yn":
				found = has_yn
			}
			if found {
				fmt.Fprintf(fp, "%s\t1\t1\t1:1\n", spod.lbl)
			} else {
				fmt.Fprintf(fp, "%s\t0\t0\t\n", spod.lbl)
			}
			fp.Close()
			continue
		}
		if spod.special == "attr" {
			var count map[string]int
			switch spod.lbl {
			case "pos":
				count = posCount
			case "postag":
				count = postagCount
			case "pt":
				count = ptCount
			}

			keys := make([]string, 0)
			sum := 0
			for key, value := range count {
				keys = append(keys, key)
				sum += value
			}
			sort.Strings(keys)
			for _, key := range keys {
				f := float64(count[key]) / float64(sum)
				fmt.Fprintf(fp, "%d\t%.3f\t%.2f%%\t%s\n", count[key], f, f*100.0, key)
			}
			fp.Close()
			continue
		}
		if spod.special == "parser" {
			if !has_parser {
				fp.Close()
				os.Remove(filename)
				continue
			}
		} else if spod.special == "qm -yn" {
			if !(has_qm || !has_yn) {
				fp.Close()
				os.Remove(filename)
				continue
			}
		} else if spod.special == "pos_verb" {
			if !has_pos_verb {
				fp.Close()
				os.Remove(filename)
				continue
			}
		} else if spod.special == "sc" {
			if !has_sc {
				fp.Close()
				os.Remove(filename)
				continue
			}
		} else if spod.special == "parser" {
			if !has_parser {
				fp.Close()
				os.Remove(filename)
				continue
			}
		} else if spod.special == "his" {
			if !has_his {
				fp.Close()
				os.Remove(filename)
				continue
			}
		}
		fmt.Fprintf(fp, "%s\t%d\t%d\t",
			spod.lbl,
			sentences[spod.lbl],
			items[spod.lbl])
		if spod.special == "parser" {
			if sentences[spod.lbl] > 0 {
				fmt.Fprintf(fp, "0:%d", sentences[spod.lbl])
			}
		} else {
			lens := make([]int, 0)
			for key := range lengths[spod.lbl] {
				lens = append(lens, key)
			}
			sort.Ints(lens)
			p := ""
			for _, key := range lens {
				fmt.Fprintf(fp, "%s%d:%d", p, key, lengths[spod.lbl][key])
				p = ","
			}
		}
		fmt.Fprintln(fp)
		fp.Close()
	}

}

func spod_fingerprint(xpath string, method string) string {
	query := macroKY.ReplaceAllStringFunc(xpath, func(s string) string {
		return rules[s[1:len(s)-1]]
	})
	query = strings.Join(strings.Fields(query), " ")
	return fmt.Sprintf("%x", md5.Sum([]byte(query+method)))
}

func loadMacros() {

	for _, set := range macroRE.FindAllStringSubmatch(macroCOM.ReplaceAllLiteralString(file__macros__txt, ""), -1) {
		s := strings.Replace(set[2], "\r\n", "\n", -1)
		s = strings.Replace(s, "\n\r", "\n", -1)
		s = strings.Replace(s, "\r", "\n", -1)
		rules[set[1]] = s
	}

	for key := range rules {
		for {
			rule := macroKY.ReplaceAllStringFunc(rules[key], func(s string) string {
				return rules[s[1:len(s)-1]]
			})
			if rule == rules[key] {
				break
			}
			if len(rule) > 100000 {
				rules[key] = "RECURSIONLIMIT"
				break
			}
			rules[key] = rule
		}
	}

}
