package alud

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"sort"
)

// meest linkse node
func nLeft(nodes []interface{}) *nodeType {
	switch len(nodes) {
	case 0:
		return noNode
	case 1:
		return nodes[0].(*nodeType)
	}
	sort.Slice(nodes, func(i, j int) bool {
		// solve cases where begin is identical (hij is en blijft omstreden)??
		ii := nodes[i].(*nodeType)
		jj := nodes[j].(*nodeType)
		if ii.Begin != jj.Begin {
			return ii.Begin < jj.Begin // ints
		}
		//if ii.End != jj.End {
		return ii.End < jj.End // ints
		//}
		//return ii.ID > jj.ID // ints, omgekeerd
	})
	return nodes[0].(*nodeType)
}

// meest linkse node als []interface{}, met lengte 0 of 1
func ifLeft(nodes []interface{}) []interface{} {
	n := nLeft(nodes)
	if n == noNode {
		return []interface{}{}
	}
	return []interface{}{n}
}

// eerste interface{} als []interface{}, met lengte 0 of 1
func if1(nodes []interface{}) []interface{} {
	if len(nodes) > 1 {
		return nodes[:1]
	}
	return nodes
}

// laatste interface{} als []interface{}, met lengte 0 of 1
func ifZ(nodes []interface{}) []interface{} {
	if len(nodes) > 0 {
		return nodes[len(nodes)-1:]
	}
	return []interface{}{}
}

// eerste node
func n1(nodes []interface{}) *nodeType {
	//return nLeft(nodes)

	if len(nodes) > 0 {
		return nodes[0].(*nodeType)
	}
	return noNode
}

// laatste node
func nZ(nodes []interface{}) *nodeType {
	if len(nodes) > 0 {
		return nodes[len(nodes)-1].(*nodeType)
	}
	return noNode
}

// eerste int
func i1(ii []interface{}) int {
	if len(ii) > 0 {
		return ii[0].(int)
	}
	return error_no_value
}

/*
// laatste int
func iZ(ii []interface{}) int {
	if l := len(ii); l > 0 {
		return ii[l-1].(int)
	}
	return error_NO_VALUE
}
*/

func depthCheck(q *context) {
	q.depth++
	if q.depth == 1000 {
		panic("Recursion depth limit")
	}
}

func dump(alpino *alpino_ds) {
	b, err := xml.MarshalIndent(alpino, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println("<?xml version=\"1.0\"?>\n" + string(b))
}

// nodes[0] -> node
// nodes[1] -> head
// nodes[2] -> gap
// nodes[3] -> subj
func trace(r interface{}, s string, q *context, nodes ...*nodeType) traceT {
	tr := traceType{s: s}
	if len(nodes) > 0 && nodes[0] != nil {
		tr.node = nodes[0]
	}
	if len(nodes) > 1 && nodes[1] != nil {
		tr.head = nodes[1]
	}
	if len(nodes) > 2 && nodes[2] != nil {
		tr.gap = nodes[2]
	}
	if len(nodes) > 3 && nodes[3] != nil {
		tr.subj = nodes[3]
	}
	switch rr := r.(type) {
	case traceT:
		rr.trace = append(rr.trace, tr)
		return rr
	default:
		return traceT{msg: fmt.Sprint(rr), debugs: q.debugs, trace: []traceType{tr}}
	}
}

func untrace(r interface{}) string {
	var buf bytes.Buffer

	switch rr := r.(type) {
	case traceT:

		fmt.Fprint(&buf, rr.msg)
		if node := rr.trace[len(rr.trace)-1].node; node != nil && node.Word != "" {
			fmt.Fprintf(&buf, " for %s:%s", number(node.End), node.Word)
		}
		for _, d := range rr.debugs {
			fmt.Fprintf(&buf, "\n  # debug: %s", d)
		}
		for _, t := range rr.trace {
			buf.WriteString("\n    in " + t.s)
			for ii, n := range []*nodeType{t.node, t.head, t.gap, t.subj} {
				if n == nil {
					continue
				}
				fmt.Fprintf(
					&buf,
					"\n        %s -- id:%s  begin:%s  end:%s",
					[]string{"node", "head", "gap ", "subj"}[ii],
					number(n.ID),
					number(n.Begin),
					number(n.End))
				if a := n.Word; a != "" {
					fmt.Fprintf(&buf, "  word:%s", a)
				}
				if a := n.Pt; a != "" {
					fmt.Fprintf(&buf, "  pt:%s", a)
				}
				if a := n.Cat; a != "" {
					fmt.Fprintf(&buf, "  cat:%s", a)
				}
				if a := n.Rel; a != "" {
					fmt.Fprintf(&buf, "  rel:%s", a)
				}
				if a := n.Conjtype; a != "" {
					fmt.Fprintf(&buf, "  conjtype:%s", a)
				}
				if a := n.Genus; a != "" {
					fmt.Fprintf(&buf, "  genus:%s", a)
				}
				if a := n.Getal; a != "" {
					fmt.Fprintf(&buf, "  getal:%s", a)
				}
				if a := n.Graad; a != "" {
					fmt.Fprintf(&buf, "  graad:%s", a)
				}
				if a := n.Index; a > 0 {
					fmt.Fprintf(&buf, "  index:%d", a)
				}
				if a := n.Lemma; a != "" {
					fmt.Fprintf(&buf, "  lemma:%s", a)
				}
				if a := n.Lwtype; a != "" {
					fmt.Fprintf(&buf, "  lwtype:%s", a)
				}
				if a := n.Naamval; a != "" {
					fmt.Fprintf(&buf, "  naamval:%s", a)
				}
				if a := n.Ntype; a != "" {
					fmt.Fprintf(&buf, "  ntype:%s", a)
				}
				if a := n.Numtype; a != "" {
					fmt.Fprintf(&buf, "  numtype:%s", a)
				}
				if a := n.Pdtype; a != "" {
					fmt.Fprintf(&buf, "  pdtype:%s", a)
				}
				if a := n.Persoon; a != "" {
					fmt.Fprintf(&buf, "  persoon:%s", a)
				}
				if a := n.Pvagr; a != "" {
					fmt.Fprintf(&buf, "  pvagr:%s", a)
				}
				if a := n.Pvtijd; a != "" {
					fmt.Fprintf(&buf, "  pvtijd:%s", a)
				}
				if a := n.Sc; a != "" {
					fmt.Fprintf(&buf, "  sc:%s", a)
				}
				if a := n.Spectype; a != "" {
					fmt.Fprintf(&buf, "  spectype:%s", a)
				}
				if a := n.Vwtype; a != "" {
					fmt.Fprintf(&buf, "  vwtype:%s", a)
				}
				if a := n.Wvorm; a != "" {
					fmt.Fprintf(&buf, "  wvorm:%s", a)
				}
			}
		}
	default:
		fmt.Fprint(&buf, rr)
	}
	return buf.String()
}
