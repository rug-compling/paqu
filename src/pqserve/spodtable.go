package main

import (
	"github.com/rug-compling/paqu/internal/dir"
	pqspod "github.com/rug-compling/paqu/internal/spod"

	"github.com/pebbe/dbxml"
	"github.com/rug-compling/alpinods"

	"encoding/xml"
	"fmt"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	// niet "na" (staat voor ontbrekend, oude bug in Alpino)
	ptlist = []string{"adj", "bw", "let", "lid", "n", "spec", "tsw", "tw", "vg", "vnw", "vz", "ww"}

	noNode = &pqspod.NodeType{
		NodeAttributes: alpinods.NodeAttributes{
			Begin: -1,
			End:   -1,
			ID:    -1,
		},
		Node:                []*pqspod.NodeType{},
		AxParent:            []interface{}{},
		AxAncestors:         []interface{}{},
		AxAncestorsOrSelf:   []interface{}{},
		AxChildren:          []interface{}{},
		AxDescendants:       []interface{}{},
		AxDescendantsOrSelf: []interface{}{},
	}
)

func spod_info(q *Context) {
	q.w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	cache(q)

	fmt.Fprint(q.w, `Beschrijving van kolommen in gedetailleerde SPOD-tabel

    sentence.id     sentence ID

    tokens          aantal tokens in de zin, zonder leestekens
    tokens.len      gemiddeld aantal letters per token, zonder leestekens

De volgende kolommen bevatten het aantal woorden met het betreffende attribuut @pt per zin (ook leestekens)

    pt.adj          bijvoeglijk naamwoord
    pt.bw           bijwoord
    pt.let          leesteken
    pt.lid          lidwoord
    pt.n            zelfstandig naamwoord
    pt.spec         speciale woordsoorten
    pt.tsw          tussenwerpsel
    pt.tw           telwoord
    pt.vg           voegwoord
    pt.vnw          voornaamwoord
    pt.vz           voorzetsel
    pt.ww           werkwoord

De volgende kolommen bevatten het aantal nodes per zin (label zonder toevoeging .len)
en het gemiddeld aantal woorden onder elke gevonden node (label met toevoeging .len)
Het laatste alleen voor een gedetailleerde tabel inclusief lengtes.
`)

	for _, spod := range pqspod.Spods {
		if spod.Special == "parser" {
			break
		}
		if spod.Special == "hidden1" || spod.Special == "attr" {
			continue
		}
		if spod.Header != "" {
			h := spod.Header
			if i := strings.Index(h, "//"); i > 0 {
				h = h[:i]
			}
			fmt.Fprintf(q.w, "\n    %s\n\n", h)
		}
		fmt.Fprintf(q.w, "    %-16s%s\n", spod.Lbl, strings.Replace(spod.Text, "|", "", -1))
	}

	fmt.Fprint(q.w, `
Parser succes:

    ok              de waarde 1 als de zin correct geparst is (parse bestaat uit een enkel deel, en geen woord is overgeslagen)
    cats            parse bestaat uit dit aantal losse delen
    skips           aantal overgeslagen woorden

Onbekende woorden per zin:

`)
	for _, spod := range pqspod.Spods {
		if spod.Special != "his" {
			continue
		}
		fmt.Fprintf(q.w, "    %-16s%s\n", spod.Lbl, strings.Replace(spod.Text, "|", "", -1))
	}

	fmt.Fprint(q.w, `
Tot slot, afhankelijk van het corpus volgt alle metadata van de zin

    meta.*

`)
}

func spod_table(q *Context, prefix string, length bool) {

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

	metalist := make([]string, 0)
	rows, err = sqlDB.Query(fmt.Sprintf("SELECT `name` FROM `%s_c_%s_midx` order by `name`", Cfg.Prefix, prefix))
	if err == nil {
		for rows.Next() {
			var name string
			err = rows.Scan(&name)
			if logerr(err) {
				rows.Close()
				return
			}
			metalist = append(metalist, name)
		}
		err = rows.Err()
		if logerr(err) {
			return
		}
	}

	directory := filepath.Join(dir.Data, "data", prefix, "spod")
	opts := make(map[string]bool)
	optlist := make([]string, 0)
	skiplist := make([]string, 0)
	unknown := make([]string, 0)
	unavail := make([]string, 0)
	cats := false
	skips := false
	for idx, spod := range pqspod.Spods {
		if spod.Special == "hidden1" {
			continue
		}
		has, err := pqspod.Has(spod.Lbl, directory)
		haserr := err != nil
		if err != nil {
			has = true
		}
		if first(q.r, fmt.Sprintf("i%d", idx)) == "t" {
			if spod.Lbl == "postag" || spod.Lbl == "pos" {
				unavail = append(unavail, spod.Lbl)
				continue
			}
			if spod.Lbl == "pt" {
				opts["pt"] = true
				optlist = append(optlist, "pt."+strings.Join(ptlist, "\tpt."))
				continue
			}
			if spod.Special == "parser" && strings.HasPrefix(spod.Lbl, "cats") {
				if !cats {
					cats = true
					if has {
						opts["cats1"] = true
						optlist = append(optlist, "cats")
						if haserr {
							unknown = append(unknown, "cats")
						}
					} else {
						skiplist = append(skiplist, "cats")
					}
				}
				continue
			}
			if spod.Special == "parser" && strings.HasPrefix(spod.Lbl, "skips") {
				if !skips {
					skips = true
					if has {
						opts["skips1"] = true
						optlist = append(optlist, "skips")
						if haserr {
							unknown = append(unknown, "skips")
						}
					} else {
						skiplist = append(skiplist, "skips")
					}
				}
				continue
			}
			if has {
				opts[spod.Lbl] = true
				optlist = append(optlist, spod.Lbl)
				if length && spod.Special != "attr" && spod.Special != "parser" {
					optlist = append(optlist, spod.Lbl+".len")
				}
				if haserr {
					unknown = append(unknown, spod.Lbl)
				}
			} else {
				skiplist = append(skiplist, spod.Lbl)
			}
		}
	}

	fmt.Fprintln(q.w, "# corpus:", q.desc[prefix])
	if len(unavail) > 0 {
		fmt.Fprintln(q.w, "# niet beschikbaar:", strings.Join(unavail, " "))
	}
	if len(skiplist) > 0 {
		fmt.Fprintln(q.w, "# niet voor dit corpus:", strings.Join(skiplist, " "))
	}
	if len(unknown) > 0 {
		fmt.Fprintln(q.w, "# beschikbaarheid nog onbekend:", strings.Join(unknown, " "))
	}
	fmt.Fprintf(q.w, "# beschrijving van kolommen: %sspodinfo\n", Cfg.Url)

	fmt.Fprint(q.w, "sentence.id\ttokens\ttokens.len")
	if len(optlist) > 0 {
		fmt.Fprint(q.w, "\t"+strings.Join(optlist, "\t"))
	}
	if len(metalist) > 0 {
		fmt.Fprint(q.w, "\tmeta."+strings.Join(metalist, "\tmeta."))
	}
	fmt.Fprintln(q.w)

	dactfiles := make([]string, 0)
	//global := false
	if strings.Contains(owner, "@") {
		dactfiles = append(dactfiles, filepath.Join(dir.Data, "data", prefix, "data.dact"))
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
			if !spod_table_file(q, docs.Name(), docs.Value(), opts, metalist, length) {
				db.Close()
				return
			}
			if ff, ok := q.w.(http.Flusher); ok {
				ff.Flush()
			}
		}
		db.Close()
	}
}

func spod_table_file(q *Context, filename string, contents string, opts map[string]bool, metalist []string, length bool) bool {

	defer fmt.Fprintln(q.w)

	meta := make(map[string][]string)
	for _, m := range metalist {
		meta[m] = make([]string, 0)
	}

	var alpino pqspod.Alpino_ds
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

	if alpino.Metadata != nil && alpino.Metadata.Meta != nil {
		for _, m := range alpino.Metadata.Meta {
			meta[m.Name] = append(meta[m.Name], m.Value)
		}
	}

	// Extra node bovenaan vanwege gedoe met //node
	alpino.Node = &pqspod.NodeType{
		NodeAttributes: alpinods.NodeAttributes{
			Begin: alpino.Node.Begin,
			End:   alpino.Node.End,
			ID:    -2, // ??? TODO
		},
		Node: []*pqspod.NodeType{alpino.Node},
	}

	ptCount := make(map[string]int)

	var walk func(*pqspod.NodeType)
	walk = func(node *pqspod.NodeType) {
		ptCount[node.Pt] = ptCount[node.Pt] + 1
		if node.Node == nil {
			node.Node = make([]*pqspod.NodeType, 0)
		} else {
			for _, n := range node.Node {
				walk(n)
			}
		}
	}
	walk(alpino.Node)

	qq := &pqspod.Context{
		Alpino:   &alpino,
		Filename: filename,
		Sentence: alpino.Sentence.Sent,
		Sentid:   alpino.Sentence.SentId,
		Varroot:  []interface{}{alpino.Node},
	}

	inspect(qq)

	tokens := 0
	tokenlen := 0
	for _, node := range qq.Ptnodes {
		if node.Pt != "let" {
			tokens++
			tokenlen += utf8.RuneCountInString(strings.Replace(node.Word, "ij", "y", -1))
		}
	}
	if tokens == 0 {
		_, err = fmt.Fprintf(q.w, "%s\t0\tNA", alpino.Sentence.SentId)
	} else {
		_, err = fmt.Fprintf(q.w, "%s\t%d\t%s", alpino.Sentence.SentId, tokens, spodfloat(float64(tokenlen)/float64(tokens)))
	}
	if logerr(err) {
		return false
	}

SPODS:
	for _, spod := range pqspod.Spods {
		if !opts[spod.Lbl] {
			continue
		}
		if spod.Lbl == "pt" {
			for _, attr := range ptlist {
				fmt.Fprintf(q.w, "\t%d", ptCount[attr])
			}
			continue
		}
		if spod.Special == "parser" {
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
			switch spod.Lbl {
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

		results, err := pqspod.Spod2xpath[spod.Lbl].Do(qq)
		if err != nil {
			fmt.Fprint(q.w, "\tERROR\tNA")
			continue SPODS
		}
		seen := make(map[int]*pqspod.NodeType)
		totalSize := 0
		if results != nil {
			for _, result := range results {
				node, ok := result.(*pqspod.NodeType)
				if !ok {
					fmt.Fprintf(q.w, "\tNA\t%T", result)
					continue SPODS
				}
				if _, ok := seen[node.ID]; !ok {
					seen[node.ID] = node
					totalSize += node.NodeSize
				}
			}
		}
		if length {
			if len(seen) == 0 {
				fmt.Fprint(q.w, "\t0\tNA")
				continue
			}
			fmt.Fprintf(q.w, "\t%d\t%s", len(seen), spodfloat(float64(totalSize)/float64(len(seen))))
		} else {
			fmt.Fprintf(q.w, "\t%d", len(seen))
		}
	}

	for _, m := range metalist {
		fmt.Fprintf(q.w, "\t%s", strings.Join(meta[m], "|"))
	}

	return true
}

func inspect(q *pqspod.Context) {
	allnodes := make([]*pqspod.NodeType, 0)
	varallnodes := make([]interface{}, 0)
	ptnodes := make([]*pqspod.NodeType, 0)
	varindexnodes := make([]interface{}, 0)

	indextable := make(map[int]*pqspod.NodeType)

	var walk func(*pqspod.NodeType)
	walk = func(node *pqspod.NodeType) {

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
			n.Parent = node
			n.AxParent = []interface{}{node}
			walk(n)
		}
		node.AxChildren = make([]interface{}, 0)
		node.AxDescendants = make([]interface{}, 0)
		node.AxDescendantsOrSelf = make([]interface{}, 1)
		node.AxDescendantsOrSelf[0] = node
		for _, n := range node.Node {
			node.AxChildren = append(node.AxChildren, n)
			node.AxDescendants = append(node.AxDescendants, n)
			node.AxDescendants = append(node.AxDescendants, n.AxDescendants...)
			node.AxDescendantsOrSelf = append(node.AxDescendantsOrSelf, n.AxDescendantsOrSelf...) // niet n
		}
	}
	walk(q.Alpino.Node)
	q.Alpino.Node.Parent = noNode
	q.Alpino.Node.AxParent = []interface{}{}

	var found map[int]bool
	walk = func(node *pqspod.NodeType) {
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
		node.NodeSize = len(found)
	}

	for _, node := range allnodes {
		node.AxAncestors = make([]interface{}, 0)
		node.AxAncestorsOrSelf = make([]interface{}, 0)
		node.AxAncestorsOrSelf = append(node.AxAncestorsOrSelf, node)
		if node != q.Alpino.Node {
			node.AxAncestors = append(node.AxAncestors, node.Parent)
			node.AxAncestors = append(node.AxAncestors, node.Parent.AxAncestors...)
			if node.AxAncestors[len(node.AxAncestors)-1] != q.Alpino.Node {
				// zou niet mogelijk moeten zijn
				panic("Missing ancestors in " + q.Filename)
			}
			node.AxAncestorsOrSelf = append(node.AxAncestorsOrSelf, node.Parent.AxAncestorsOrSelf...)
		}
	}

	sort.Slice(ptnodes, func(i, j int) bool {
		return ptnodes[i].End < ptnodes[j].End
	})
	varptnodes := make([]interface{}, len(ptnodes))
	for i, node := range ptnodes {
		varptnodes[i] = node
	}

	q.Allnodes = allnodes
	q.Varallnodes = varallnodes
	q.Varindexnodes = varindexnodes
	q.Ptnodes = ptnodes
	q.Varptnodes = varptnodes

}

func spodfloat(f float64) string {
	s := fmt.Sprintf("%.3f", f)
	n := len(s)
	for i := 1; i < 3; i++ {
		if s[n-i] == '0' {
			s = s[:n-i]
		} else {
			break
		}
	}
	return s
}
