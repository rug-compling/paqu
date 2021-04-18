package main

import (
	"github.com/rug-compling/alpinods"

	"encoding/xml"
	"fmt"
	"os"
	"runtime/debug"
	"sort"
	"strings"
)

type indexType int

type doer interface {
	do(subdoc []interface{}, q *context) []interface{}
}

type context struct {
	alpino        *alpino_ds
	filename      string
	sentence      string
	sentid        string
	depth         int
	allnodes      []*nodeType
	ptnodes       []*nodeType
	varallnodes   []interface{}
	varindexnodes []interface{}
	varptnodes    []interface{}
	varroot       []interface{}
}

type alpino_ds struct {
	XMLName  xml.Name      `xml:"alpino_ds"`
	Version  string        `xml:"version,attr,omitempty"`
	Metadata *metadataType `xml:"metadata,omitempty"`
	Parser   *parserType   `xml:"parser,omitempty"`
	Node     *nodeType     `xml:"node,omitempty"`
	Sentence *sentType     `xml:"sentence,omitempty"`
}

type metadataType struct {
	Meta []metaType `xml:"meta,omitempty"`
}

type metaType struct {
	Type  string `xml:"type,attr,omitempty"`
	Name  string `xml:"name,attr,omitempty"`
	Value string `xml:"value,attr,omitempty"`
}

type parserType struct {
	Build string `xml:"build,attr,omitempty"`
	Date  string `xml:"date,attr,omitempty"`
	Cats  string `xml:"cats,attr,omitempty"`
	Skips string `xml:"skips,attr,omitempty"`
}

type sentType struct {
	Sent   string `xml:",chardata"`
	SentId string `xml:"sentid,attr,omitempty"`
}

type nodeType struct {
	alpinods.NodeAttributes
	Node   []*nodeType `xml:"node"`
	parent *nodeType
	size   int

	axParent            []interface{}
	axAncestors         []interface{}
	axAncestorsOrSelf   []interface{}
	axChildren          []interface{}
	axDescendants       []interface{}
	axDescendantsOrSelf []interface{}
}

const (
	cmp__lt = iota
	cmp__gt
	collect__ancestors__node
	collect__ancestors__or__self__node
	collect__attributes__begin
	collect__attributes__end
	collect__attributes__cat
	collect__attributes__frame
	collect__attributes__graad
	collect__attributes__his
	collect__attributes__index
	collect__attributes__lcat
	collect__attributes__lemma
	collect__attributes__naamval
	collect__attributes__pdtype
	collect__attributes__pt
	collect__attributes__pvagr
	collect__attributes__pvtijd
	collect__attributes__rel
	collect__attributes__sc
	collect__attributes__spectype
	collect__attributes__stype
	collect__attributes__tense
	collect__attributes__vwtype
	collect__attributes__word
	collect__attributes__wvorm
	collect__child__node
	collect__descendant__node
	collect__descendant__or__self__node
	collect__descendant__or__self__type__node
	collect__parent__type__node
	collect__parent__node
	collect__self__all__node
	collect__self__node
	equal__is
	function__contains__2__args
	function__count__1__args
	function__ends__with__2__args
	function__first__0__args // LET OP: extra gebruik in (*dCollect).do()
	function__third__0__args // LET OP: extra gebruik in (*dCollect).do()
	function__last__0__args  // LET OP: extra gebruik in (*dCollect).do()
	// function__local__internal__head__position__1__args
	function__not__1__args
	function__starts__with__2__args
	plus__plus
	plus__minus
)

var (
	nTRUE  = []interface{}{true}
	nFALSE = []interface{}{}
)

type dAnd struct {
	arg1 doer
	arg2 doer
}

func (d *dAnd) do(subdoc []interface{}, q *context) []interface{} {
	for _, a := range []doer{d.arg1, d.arg2} {
		if r := a.do(subdoc, q); len(r) == 0 {
			return nFALSE
		}
	}
	return nTRUE
}

type dArg struct {
	arg1 doer
	arg2 doer
}

func (d *dArg) do(subdoc []interface{}, q *context) []interface{} {
	result := []interface{}{}
	for i, a := range []doer{d.arg1, d.arg2} {
		if i == 0 || a != nil {
			// TODO: waarom flatten?
			result = append(result, flatten(a.do(subdoc, q)))
		}
	}
	return result
}

type dCmp struct {
	ARG  int
	arg1 doer
	arg2 doer
}

func (d *dCmp) do(subdoc []interface{}, q *context) []interface{} {
	arg1 := d.arg1.do(subdoc, q)
	arg2 := d.arg2.do(subdoc, q)
	switch d.ARG {
	case cmp__lt: // <
		for _, a1 := range arg1 {
			for _, a2 := range arg2 {
				switch a1t := a1.(type) {
				case int:
					if a1t < a2.(int) {
						return nTRUE
					}
				case string:
					if a1t < a2.(string) {
						return nTRUE
					}
				default:
					panic(fmt.Sprintf("Cmp<: Missing case for type %T in %s", a1, q.filename))
				}
			}
		}
		return nFALSE
	case cmp__gt: // >
		for _, a1 := range arg1 {
			for _, a2 := range arg2 {
				switch a1t := a1.(type) {
				case int:
					if a1t > a2.(int) {
						return nTRUE
					}
				case string:
					if a1t > a2.(string) {
						return nTRUE
					}
				default:
					panic(fmt.Sprintf("Cmp>: Missing case for type %T in %s", a1, q.filename))
				}
			}
		}
		return nFALSE
	default:
		panic("Cmp: Missing case in " + q.filename)
	}
}

type dCollect struct {
	ARG  int
	arg1 doer
	arg2 doer
}

func (d *dCollect) do(subdoc []interface{}, q *context) []interface{} {

	lists := [][]interface{}{}

	result1 := []interface{}{}
	for _, r := range d.arg1.do(subdoc, q) {
		switch d.ARG {
		case collect__ancestors__node:
			lists = append(lists, r.(*nodeType).axAncestors)
		case collect__ancestors__or__self__node:
			lists = append(lists, r.(*nodeType).axAncestorsOrSelf)
		case collect__attributes__begin:
			if i := r.(*nodeType).Begin; i >= 0 {
				result1 = append(result1, i)
			}
		case collect__attributes__cat:
			if i := r.(*nodeType).Cat; i != "" {
				result1 = append(result1, i)
			}
		case collect__attributes__frame:
			if i := r.(*nodeType).Frame; i != "" {
				result1 = append(result1, i)
			}
		case collect__attributes__graad:
			if i := r.(*nodeType).Graad; i != "" {
				result1 = append(result1, i)
			}
		case collect__attributes__end:
			if i := r.(*nodeType).End; i >= 0 {
				result1 = append(result1, i)
			}
		case collect__attributes__his:
			if i := r.(*nodeType).His; i != "" {
				result1 = append(result1, i)
			}
		case collect__attributes__index:
			if i := r.(*nodeType).Index; i > 0 {
				result1 = append(result1, i)
			}
		case collect__attributes__lcat:
			if i := r.(*nodeType).Lcat; i != "" {
				result1 = append(result1, i)
			}
		case collect__attributes__lemma:
			if i := r.(*nodeType).Lemma; i != "" {
				result1 = append(result1, i)
			}
		case collect__attributes__naamval:
			if i := r.(*nodeType).Naamval; i != "" {
				result1 = append(result1, i)
			}
		case collect__attributes__pdtype:
			if i := r.(*nodeType).Pdtype; i != "" {
				result1 = append(result1, i)
			}
		case collect__attributes__pt:
			if i := r.(*nodeType).Pt; i != "" {
				result1 = append(result1, i)
			}
		case collect__attributes__pvagr:
			if i := r.(*nodeType).Pvagr; i != "" {
				result1 = append(result1, i)
			}
		case collect__attributes__pvtijd:
			if i := r.(*nodeType).Pvtijd; i != "" {
				result1 = append(result1, i)
			}
		case collect__attributes__rel:
			if i := r.(*nodeType).Rel; i != "" {
				result1 = append(result1, i)
			}
		case collect__attributes__sc:
			if i := r.(*nodeType).Sc; i != "" {
				result1 = append(result1, i)
			}
		case collect__attributes__spectype:
			if i := r.(*nodeType).Spectype; i != "" {
				result1 = append(result1, i)
			}
		case collect__attributes__stype:
			if i := r.(*nodeType).Stype; i != "" {
				result1 = append(result1, i)
			}
		case collect__attributes__tense:
			if i := r.(*nodeType).Stype; i != "" {
				result1 = append(result1, i)
			}
		case collect__attributes__vwtype:
			if i := r.(*nodeType).Vwtype; i != "" {
				result1 = append(result1, i)
			}
		case collect__attributes__word:
			if i := r.(*nodeType).Word; i != "" {
				result1 = append(result1, i)
			}
		case collect__attributes__wvorm:
			if i := r.(*nodeType).Wvorm; i != "" {
				result1 = append(result1, i)
			}
		case collect__child__node:
			lists = append(lists, r.(*nodeType).axChildren)
		case collect__descendant__node:
			lists = append(lists, r.(*nodeType).axDescendants)
		case collect__descendant__or__self__type__node, collect__descendant__or__self__node:
			lists = append(lists, r.(*nodeType).axDescendantsOrSelf)
		case collect__parent__type__node, collect__parent__node:
			lists = append(lists, r.(*nodeType).axParent)
		case collect__self__all__node, collect__self__node:
			result1 = append(result1, r) // TODO: correct?
		default:
			panic("Collect: Missing case in " + q.filename)
		}
	}

	if d.arg2 == nil {
		for _, list := range lists {
			result1 = append(result1, list...)
		}
		return result1
	}

	if len(result1) > 0 {
		l := lists
		lists = [][]interface{}{result1}
		lists = append(lists, l...)
	}

	result2 := []interface{}{}

	/*
			if pred, ok := d.arg2.(*dPredicate); ok {
				result2 := []interface{}{}
				for _, list := range lists {
					ll := pred.do(list, q)
					if len(ll) == 0 {
						continue
					}
					if _, ok := ll[0].(*nodeType); ok {
						result2 = append(result2, ll...)
					} else {
						result2 = append(result2, list...)
					}
				}
				return flatten(result2)
			}

		panic("NOT HERE")
		// TODO: panic????

	*/

	if p, ok := d.arg2.(*dPredicate); ok {
		if f, ok := p.arg2.(*dFunction); ok {
			if f.ARG == function__first__0__args || f.ARG == function__third__0__args || f.ARG == function__last__0__args {
				for _, list := range lists {
					r1 := []interface{}{}
					for _, e := range list {
						if len(p.arg1.do([]interface{}{e}, q)) > 0 {
							r1 = append(r1, e)
						}
					}
					if len(r1) == 0 {
						continue
					}
					switch f.ARG {
					case function__first__0__args:
						result2 = append(result2, r1[0])
					case function__third__0__args:
						if len(r1) > 2 {
							result2 = append(result2, r1[2])
						}
					case function__last__0__args:
						result2 = append(result2, r1[len(r1)-1])
					default:
						panic("Collect: Missing case for nested index in " + q.filename)
					}
				}
				return result2
			}
		}
	}

	for _, list := range lists {
		for _, e := range list {
			for _, r2 := range d.arg2.do([]interface{}{e}, q) {
				if idx, ok := r2.(indexType); ok {
					switch idx {
					case 1:
						result2 = append(result2, list[0])
					case 3:
						if len(list) > 2 {
							result2 = append(result2, list[2])
						}
					case -1:
						result2 = append(result2, list[len(list)-1])
					default:
						panic("Collect: Missing case for plain index in " + q.filename)
					}
					continue
				}
				result2 = append(result2, e)
			}
		}
	}

	return result2
}

type dElem struct {
	DATA []interface{}
	arg1 doer
}

func (d *dElem) do(subdoc []interface{}, q *context) []interface{} {

	if d.arg1 == nil {
		return d.DATA
	}

	/*
		TODO: arg1 negeren: klopt dat? zie bijvoorbeeld: foo[@a + 10 = @b]

		Waarvoor dient arg1 dan?

		  SORT
		    COLLECT  'child' 'name' 'node' foo
		      NODE
		      PREDICATE
		        EQUAL =
		          PLUS +
		            COLLECT  'attributes' 'name' 'node' a
		              NODE
		            ELEM Object is a number : 10
		              COLLECT  'attributes' 'name' 'node' a
		                NODE
		          COLLECT  'attributes' 'name' 'node' b
		            NODE
	*/
	return d.DATA
}

type dEqual struct {
	ARG  int
	arg1 doer
	arg2 doer
}

func (d *dEqual) do(subdoc []interface{}, q *context) []interface{} {
	result := make([]interface{}, 0)
	switch d.ARG {
	case equal__is:
		a1 := d.arg1.do(subdoc, q)
		a2 := d.arg2.do(subdoc, q)
		for _, aa1 := range a1 {
			for _, aa2 := range a2 {
				if aa1 == aa2 {
					result = append(result, true)
				}
			}
		}
		return result
	default:
		panic("Equal: Missing case in " + q.filename)
	}
}

type dFilter struct {
	arg1 doer
	arg2 doer
}

func (d *dFilter) do(subdoc []interface{}, q *context) []interface{} {

	result := []interface{}{}
	r1 := d.arg1.do(subdoc, q)
	for _, r := range r1 {
		r2 := d.arg2.do([]interface{}{r}, q)
		if len(r2) > 0 {
			if idx, ok := r2[0].(indexType); ok {
				switch idx {
				case 1:
					return []interface{}{r1[0]}
				case 3:
					if len(r1) > 2 {
						return []interface{}{r1[2]}
					} else {
						return []interface{}{}
					}
				case -1:
					return []interface{}{r1[len(r1)-1]}
				default:
					panic(fmt.Sprintf("Filter: Missing case for index %d in %s", int(idx), q.filename))
				}
			} else {
				result = append(result, r)
			}
		}
	}
	return result
}

type dFunction struct {
	ARG  int
	arg1 doer
}

func (d *dFunction) do(subdoc []interface{}, q *context) []interface{} {

	var r []interface{}
	if d.arg1 != nil {
		r = d.arg1.do(subdoc, q)
	} else {
		r = subdoc
	}

	switch d.ARG {
	case function__contains__2__args:
		for _, s1 := range r[1].([]interface{}) {
			for _, s0 := range r[0].([]interface{}) {
				if strings.Contains(s0.(string), s1.(string)) {
					return nTRUE
				}
			}
		}
		return nFALSE
	case function__count__1__args:
		return []interface{}{len(r[0].([]interface{}))}
	case function__ends__with__2__args:
		for _, s1 := range r[1].([]interface{}) {
			for _, s0 := range r[0].([]interface{}) {
				if strings.HasSuffix(s0.(string), s1.(string)) {
					return nTRUE
				}
			}
		}
		return nFALSE
	case function__first__0__args:
		return []interface{}{indexType(1)}
	case function__third__0__args:
		return []interface{}{indexType(3)}
	case function__last__0__args:
		return []interface{}{indexType(-1)}
		/*
			case function__local__internal__head__position__1__args:
				head, err := internalHeadPositionWithError(r[0].([]interface{}), q)
				if err == nil {
					return []interface{}{head}
				}
				return []interface{}{}
		*/
	case function__not__1__args:
		if len(r[0].([]interface{})) == 0 {
			return nTRUE
		}
		return nFALSE
	case function__starts__with__2__args:
		for _, s1 := range r[1].([]interface{}) {
			for _, s0 := range r[0].([]interface{}) {
				if strings.HasPrefix(s0.(string), s1.(string)) {
					return nTRUE
				}
			}
		}
		return nFALSE
	default:
		panic("Function: Missing case in " + q.filename)
	}
}

type dNode struct {
}

func (d *dNode) do(subdoc []interface{}, q *context) []interface{} {
	return subdoc
}

type dOr struct {
	arg1 doer
	arg2 doer
}

func (d *dOr) do(subdoc []interface{}, q *context) []interface{} {
	for _, a := range []doer{d.arg1, d.arg2} {
		if r := a.do(subdoc, q); len(r) > 0 {
			return nTRUE
		}
	}
	return nFALSE
}

type dPlus struct {
	ARG  int
	arg1 doer
	arg2 doer
}

func (d *dPlus) do(subdoc []interface{}, q *context) []interface{} {
	switch d.ARG {
	case plus__plus:
		result := []interface{}{}
		a1 := d.arg1.do(subdoc, q)
		a2 := d.arg2.do(subdoc, q)
		for _, aa1 := range a1 {
			for _, aa2 := range a2 {
				result = append(result, aa1.(int)+aa2.(int))
			}
		}
		return result
	case plus__minus:
		result := []interface{}{}
		a1 := d.arg1.do(subdoc, q)
		a2 := d.arg2.do(subdoc, q)
		for _, aa1 := range a1 {
			for _, aa2 := range a2 {
				result = append(result, aa1.(int)-aa2.(int))
			}
		}
		return result
	default:
		panic("Plus: Missing case in " + q.filename)
	}
}

type dPredicate struct {
	arg1 doer
	arg2 doer
}

func (d *dPredicate) do(subdoc []interface{}, q *context) []interface{} {

	result := d.arg1.do(subdoc, q)

	if d.arg2 == nil || len(result) == 0 {
		return result
	}

	p := make([]interface{}, 0)
	switch result[0].(indexType) {
	case 1:
		p = append(p, subdoc[0])
	case 3:
		if len(subdoc) > 2 {
			p = append(p, subdoc[2])
		}
	case -1:
		p = append(p, subdoc[len(subdoc)-1])
	}
	return d.arg2.do(p, q)

	idx := d.arg2.do(result, q)[0].(indexType) // TODO: altijd een index?
	switch idx {
	case 1:
		return []interface{}{result[0]}
	case 3:
		if len(result) > 2 {
			return []interface{}{result[2]}
		}
		return nFALSE
	case -1:
		return []interface{}{result[len(result)-1]}
	default:
		panic(fmt.Sprintf("Predicate arg2: Missing case for index %d in %s", int(idx), q.filename))
	}

	/*
		result2 := make([]interface{}, 0)
		for _, r := range result {
			result2 = append(result2, d.arg2.do([]interface{}{r}, q)...)
		}
		return result2
	*/

	result2 := d.arg2.do(result, q)
	return result2

	/*

		if r, ok := result[0].(indexType); ok {
			idx := int(r)
			if idx < 0 {
				idx = len(subdoc) + idx
			} else {
				idx--
			}
			if idx < 0 || idx >= len(subdoc) {
				return []interface{}{}
			}
			result = []interface{}{subdoc[idx]}
		}

	*/

	/*
		if _, ok := d.arg1.(*dPredicate); ok {
			var idx int
			switch t := result[0].(type) {
			case indexType:
				idx = int(t)
			case int:
				idx = t
			default:
				panic(fmt.Sprintf("invalid type %T", result[0]))
			}


			if idx > len(subdoc) {
				return []interface{}{}
			}
			return d.arg2.do([]interface{}{subdoc[idx-1]}, q)
		}
	*/

}

type dRoot struct {
}

func (d *dRoot) do(subdoc []interface{}, q *context) []interface{} {
	return q.varroot
}

type dSort struct {
	arg1 doer
}

func (d *dSort) do(subdoc []interface{}, q *context) []interface{} {
	result := d.arg1.do(subdoc, q)
	if len(result) < 2 {
		return result
	}

	if _, ok := result[0].([]interface{}); ok {
		res := make([]interface{}, 0)
		for _, r := range result {
			res = append(res, r.([]interface{})...)
		}
		result = res
		if len(result) < 2 {
			return result
		}
	}

	switch result[0].(type) {
	case *nodeType:
		sort.Slice(result, func(i, j int) bool {
			return result[i].(*nodeType).ID < result[j].(*nodeType).ID
		})
		for i := 1; i < len(result); i++ {
			if result[i].(*nodeType) == result[i-1].(*nodeType) {
				result = append(result[:i], result[i+1:]...)
				i--
			}
		}
	case string:
		sort.Slice(result, func(i, j int) bool {
			return result[i].(string) < result[j].(string)
		})
		for i := 1; i < len(result); i++ {
			if result[i].(string) == result[i-1].(string) {
				result = append(result[:i], result[i+1:]...)
				i--
			}
		}
	case int:
		sort.Slice(result, func(i, j int) bool {
			return result[i].(int) < result[j].(int)
		})
		for i := 1; i < len(result); i++ {
			if result[i].(int) == result[i-1].(int) {
				result = append(result[:i], result[i+1:]...)
				i--
			}
		}
	default:
		panic(fmt.Sprintf("Sort: Missing case for type %T in %s", result[0], q.filename))
	}
	return result
}

type dVariable struct {
	VAR interface{}
}

func (d *dVariable) do(subdoc []interface{}, q *context) []interface{} {
	switch t := d.VAR.(type) {
	case []interface{}:
		return t
	case *nodeType:
		return []interface{}{t}
	case []*nodeType:
		ii := make([]interface{}, len(t))
		for i, v := range t {
			ii[i] = v
		}
		return ii
	case int:
		return []interface{}{t}
	case string:
		return []interface{}{t}
	default:
		panic(fmt.Sprintf("Variable: Missing case for type %T in %s", t, q.filename))
	}
}

type xPath struct {
	arg1 doer
}

func (d *xPath) do(q *context) (result []interface{}, err error) {

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
			fmt.Fprintf(os.Stderr, "\n\n>>> %s -- %v\n\n", q.sentid, err)
			debug.PrintStack()
			fmt.Fprintln(os.Stderr, "\n<<<\n\n")
		}
	}()

	result = d.arg1.do([]interface{}{}, q)

	return
}

////////////////////////////////////////////////////////////////

func list(i interface{}) []interface{} {
	switch ii := i.(type) {
	case []interface{}:
		return ii
	case []*nodeType:
		doc := []interface{}{}
		for _, n := range ii {
			doc = append(doc, n)
		}
		return doc
	default:
		return []interface{}{ii}
	}
}

func flatten(aa []interface{}) []interface{} {
	result := make([]interface{}, 0)
	for _, a := range aa {
		switch t := a.(type) {
		case []interface{}:
			result = append(result, flatten(t)...)
		default:
			result = append(result, a)
		}
	}
	return result
}
