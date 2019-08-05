package alud

import (
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
		//return ii.Id > jj.Id // ints, omgekeerd
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
	return error_NO_VALUE
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

func depthCheck(q *context, s string) {
	q.depth++
	if q.depth == 1000 {
		panic("Recursion depth limit for " + s)
	}
}

/*
func dump(alpino *Alpino_ds) {
	b, err := xml.MarshalIndent(alpino, "", "  ")
	x(err)
	s := strings.Replace(string(b), "000", "", -1)
	fmt.Println("<?xml version=\"1.0\"?>\n" + s)
}
*/
