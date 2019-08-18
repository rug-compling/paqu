package alud

import (
	"encoding/xml"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

// updates to the output
const VersionMajor = 4

// updates to the package API (unlikely)
const VersionMinor = 2

const (
	UDVersionMajor = 2
	UDVersionMinor = 2 // not yet 3
)

// options can be or'ed as last argument to Ud()
const (
	OPT_DEBUG                  = 1 << iota // include debug messages in comments
	OPT_NO_COMMENTS                        // don't include comments
	OPT_NO_DETOKENIZE                      // don't try to restore detokenized sentence
	OPT_NO_ENHANCED                        // skip enhanced dependencies
	OPT_NO_FIX_PUNCT                       // don't fix punctuation
	OPT_NO_FIX_MISPLACED_HEADS             // don't fix misplaced heads in coordination
	OPT_PANIC                              // panic on error (for development)
)

const (
	underscore = -1000 * (iota + 1)
	empty_head
	error_no_head
	error_no_value
)

var (
	noNode = &nodeType{
		Begin:               -1000,
		End:                 -1000,
		udCopiedFrom:        -1000,
		Id:                  -1,
		Node:                []*nodeType{},
		axParent:            []interface{}{},
		axAncestors:         []interface{}{},
		axChildren:          []interface{}{},
		axDescendants:       []interface{}{},
		axDescendantsOrSelf: []interface{}{},
		udHeadPosition:      error_no_head,
		udEHeadPosition:     error_no_head,
	}
	versionID = fmt.Sprintf("ALUD%d.%d", VersionMajor, VersionMinor)
)

func init() {
	noNode.parent = noNode
}

// Version identifier in format ALUDmajor.minor
func VersionID() string {
	return versionID
}


// Derive Universal Dependencies from parsed sentence in alpino_ds format.
func Ud(alpino_doc []byte, filename string, options int) (conllu string, err error) {
	conllu, _, err = ud(alpino_doc, filename, options)
	return
}

func ud(alpino_doc []byte, filename string, options int) (conllu string, q *context, err error) {
	if options&OPT_PANIC == 0 {
		defer func() {
			if r := recover(); r != nil {
				conllu = ""
				err = fmt.Errorf("%v", r)
			}
		}()
	}
	conllu, q, err = udTry(alpino_doc, filename, options)
	return
}

func udTry(alpino_doc []byte, filename string, options int) (conllu string, q *context, err error) {

	var alpino alpino_ds
	err = xml.Unmarshal(alpino_doc, &alpino)
	if err != nil {
		return "", nil, err
	}

	if alpino.Sentence.SentId == "" {
		id := filepath.Base(filename)
		if strings.HasSuffix(id, ".xml") {
			id = id[:len(id)-4]
		}
		alpino.Sentence.SentId = id
	}

	var walk func(*nodeType)
	walk = func(node *nodeType) {
		node.Begin *= 1000
		node.End *= 1000
		if node.Node == nil {
			node.Node = make([]*nodeType, 0)
		} else {
			for _, n := range node.Node {
				walk(n)
			}
		}
	}
	walk(alpino.Node)

	q = &context{
		alpino:   &alpino,
		filename: filename,
		sentence: alpino.Sentence.Sent,
		sentid:   alpino.Sentence.SentId,
		varroot:  []interface{}{alpino.Node},
		swapped:  [][2]*nodeType{},
	}

	inspect(q)

	if options&OPT_NO_FIX_MISPLACED_HEADS == 0 {
		fixMisplacedHeadsInCoordination(q)
	}
	addPosTags(q)
	addFeatures(q)
	addDependencyRelations(q)
	if options&OPT_NO_ENHANCED == 0 {
		enhancedDependencies(q)
	}
	if options&OPT_NO_FIX_PUNCT == 0 {
		fixpunct(q)
	}
	if options&OPT_NO_DETOKENIZE == 0 {
		untokenize(q)
	}
	return conll(q, options), q, nil
}

func inspect(q *context) {
	allnodes := make([]*nodeType, 0)
	varallnodes := make([]interface{}, 0)
	ptnodes := make([]*nodeType, 0)
	varindexnodes := make([]interface{}, 0)

	var walk func(*nodeType)
	walk = func(node *nodeType) {

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

	for _, node := range allnodes {
		node.axAncestors = make([]interface{}, 0)
		if node != q.alpino.Node {
			node.axAncestors = append(node.axAncestors, node.parent)
			node.axAncestors = append(node.axAncestors, node.parent.axAncestors...)
			if node.axAncestors[len(node.axAncestors)-1] != q.alpino.Node {
				// zou niet mogelijk moeten zijn
				panic("Missing ancestors in " + q.filename)
			}
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
