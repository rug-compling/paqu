package main

import (
	// "github.com/kr/pretty"
	"github.com/pebbe/dbxml"
	"github.com/rug-compling/alpinods"

	"encoding/xml"
	"fmt"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

var (
	ptlist = []string{"adj", "bw", "let", "lid", "n", "na", "spec", "tsw", "tw", "vg", "vnw", "vz", "ww"}

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
)

func spod_table(q *Context, prefix string) {

	var owner string
	rows, err := sqlDB.Query(fmt.Sprintf("SELECT `owner` FROM `%s_info` WHERE `id` = %q", Cfg.Prefix, prefix))
	if logerr(err) {
		return
	}
	for rows.Next() {
		err = rows.Scan(&owner)
		if logerr(err) {
			rows.Close()
			return
		}
	}
	err = rows.Err()
	if logerr(err) {
		return
	}

	opts := make(map[string]bool)
	optlist := make([]string, 0)
	cats := false
	skips := false
	for idx, spod := range spods {
		if spod.special == "hidden1" {
			continue
		}
		if spod.lbl == "postag" || spod.lbl == "pos" {
			continue
		}
		if first(q.r, fmt.Sprintf("i%d", idx)) == "t" {
			if spod.lbl == "pt" {
				opts["pt"] = true
				optlist = append(optlist, "pt."+strings.Join(ptlist, "\tpt."))
				continue
			}
			if spod.special == "parser" && strings.HasPrefix(spod.lbl, "cats") {
				if !cats {
					cats = true
					opts["cats1"] = true
					optlist = append(optlist, "cats")
				}
				continue
			}
			if spod.special == "parser" && strings.HasPrefix(spod.lbl, "skips") {
				if !skips {
					skips = true
					opts["skips1"] = true
					optlist = append(optlist, "skips")
				}
				continue
			}
			opts[spod.lbl] = true
			optlist = append(optlist, spod.lbl)
			if spod.special != "attr" && spod.special != "parser" {
				optlist = append(optlist, spod.lbl+".len")
			}
		}
	}

	fmt.Fprintf(q.w, "sentence.id\ttokens\ttokens.len\t%s\n", strings.Join(optlist, "\t"))

	dactfiles := make([]string, 0)
	//global := false
	if strings.Contains(owner, "@") {
		dactfiles = append(dactfiles, filepath.Join(paqudatadir, "data", prefix, "data.dact"))
	} else {
		//global = true
		rows, err = sqlDB.Query(fmt.Sprintf("SELECT `arch` FROM `%s_c_%s_arch` ORDER BY `id`", Cfg.Prefix, prefix))
		if logerr(err) {
			return
		}
		for rows.Next() {
			var s string
			err = rows.Scan(&s)
			if logerr(err) {
				rows.Close()
				return
			}
			if strings.HasSuffix(s, ".dact") {
				dactfiles = append(dactfiles, s)
			}
		}
		err = rows.Err()
		if logerr(err) {
			return
		}
	}

	//n := 0

	for _, dactfile := range dactfiles {
		db, err := dbxml.OpenRead(dactfile)
		if logerr(err) {
			return
		}
		docs, err := db.All()
		if logerr(err) {
			db.Close()
			return
		}
		for docs.Next() {
			if !spod_table_file(q, docs.Name(), docs.Value(), opts) {
				db.Close()
				return
			}
			if ff, ok := q.w.(http.Flusher); ok {
				ff.Flush()
			}
			//	n++
			//	if n == 10 {
			//		db.Close()
			//		return
			//	}
		}

		db.Close()
	}

}

func spod_table_file(q *Context, filename string, contents string, opts map[string]bool) bool {

	defer fmt.Fprintln(q.w)

	var alpino alpino_ds
	err := xml.Unmarshal([]byte(contents), &alpino)
	if logerr(err) {
		return false
	}

	if alpino.Sentence.SentId == "" {
		id := filepath.Base(filename)
		if strings.HasSuffix(id, ".xml") {
			id = id[:len(id)-4]
		}
		alpino.Sentence.SentId = id
	}

	ptCount := make(map[string]int)

	var walk func(*nodeType)
	walk = func(node *nodeType) {
		ptCount[node.Pt] = ptCount[node.Pt] + 1
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
		filename: filename,
		sentence: alpino.Sentence.Sent,
		sentid:   alpino.Sentence.SentId,
		varroot:  []interface{}{alpino.Node},
	}

	inspect(qq)

	tokens := 0
	tokenlen := 0
	for _, node := range qq.ptnodes {
		if node.Pt != "let" {
			tokens++
			tokenlen += len(node.Word)
		}
	}
	_, err = fmt.Fprintf(q.w, "%s\t%d\t%.1f", alpino.Sentence.SentId, tokens, float64(tokenlen)/float64(tokens))
	if logerr(err) {
		return false
	}

SPODS:
	for _, spod := range spods {
		if !opts[spod.lbl] {
			continue
		}
		if spod.lbl == "pt" {
			for _, attr := range ptlist {
				fmt.Fprintf(q.w, "\t%d", ptCount[attr])
			}
			continue
		}
		if spod.special == "parser" {
			if alpino.Parser == nil {
				fmt.Fprint(q.w, "\tNA")
				continue
			}
			cats, e2 := strconv.Atoi(alpino.Parser.Cats)
			skips, e1 := strconv.Atoi(alpino.Parser.Skips)
			if e1 != nil || e2 != nil {
				fmt.Fprint(q.w, "\tNA")
				continue
			}
			switch spod.lbl {
			case "ok":
				if cats == 1 && skips == 0 {
					fmt.Fprint(q.w, "\t1")
				} else {
					fmt.Fprint(q.w, "\t0")
				}
			case "cats1":
				fmt.Fprintf(q.w, "\t%d", cats)
			case "skips1":
				fmt.Fprintf(q.w, "\t%d", skips)
			}
			// TODO
			continue
		}

		results, err := spod2xpath[spod.lbl].do(qq)
		if err != nil {
			fmt.Fprint(q.w, "\tERROR\tNA")
			continue SPODS
		}
		seen := make(map[int]*nodeType)
		totalSize := 0
		if results != nil {
			for _, result := range results {
				node, ok := result.(*nodeType)
				if !ok {
					fmt.Fprintf(q.w, "\tNA\t%T", result)
					continue SPODS
				}
				if _, ok := seen[node.ID]; !ok {
					seen[node.ID] = node
					totalSize += node.size
				}
			}
		}
		if len(seen) == 0 {
			fmt.Fprint(q.w, "\t0\tNA")
			continue
		}
		fmt.Fprintf(q.w, "\t%d\t%.1f", len(seen), float64(totalSize)/float64(len(seen)))
	}

	return true
}

func inspect(q *context) {
	allnodes := make([]*nodeType, 0)
	varallnodes := make([]interface{}, 0)
	ptnodes := make([]*nodeType, 0)
	varindexnodes := make([]interface{}, 0)

	indextable := make(map[int]*nodeType)

	var walk func(*nodeType)
	walk = func(node *nodeType) {

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
