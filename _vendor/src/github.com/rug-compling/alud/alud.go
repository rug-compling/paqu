package alud

import (
	"encoding/xml"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

// updates to the output
const VersionMajor = int(1)

// updates to the package API (unlikely)
const VersionMinor = int(0)

const (
	error_EXTERNAL_HEAD_MUST_HAVE_ONE_ARG = -1000 * (iota + 1)
	error_MORE_THAN_ONE_INTERNAL_HEAD_POSITION_FOUND
	error_NO_EXTERNAL_HEAD
	error_NO_HEAD_FOUND
	error_NO_INTERNAL_HEAD
	error_NO_INTERNAL_HEAD_IN_GAPPED_CONSTITUENT
	error_NO_INTERNAL_HEAD_POSITION_FOUND
	error_NO_VALUE
	error_RECURSION_LIMIT
	underscore
	empty_head
)

type context struct {
	alpino        *alpino_ds
	filename      string
	sentence      string
	sentid        string
	debugs        []string
	warnings      []string
	depth         int
	allnodes      []*nodeType
	ptnodes       []*nodeType
	varallnodes   []interface{}
	varindexnodes []interface{}
	varptnodes    []interface{}
	varroot       []interface{}
}

type alpino_ds struct {
	XMLName  xml.Name  `xml:"alpino_ds"`
	Node     *nodeType `xml:"node,omitempty"`
	Sentence *sentType `xml:"sentence,omitempty"`
}

type sentType struct {
	Sent   string `xml:",chardata"`
	SentId string `xml:"sentid,attr,omitempty"`
}

type nodeType struct {
	Begin    int         `xml:"begin,attr"`
	Cat      string      `xml:"cat,attr,omitempty"`
	Conjtype string      `xml:"conjtype,attr,omitempty"`
	End      int         `xml:"end,attr"`
	Genus    string      `xml:"genus,attr,omitempty"`
	Getal    string      `xml:"getal,attr,omitempty"`
	Graad    string      `xml:"graad,attr,omitempty"`
	Id       int         `xml:"id,attr,omitempty"`
	Index    int         `xml:"index,attr,omitempty"`
	Lemma    string      `xml:"lemma,attr,omitempty"`
	Lwtype   string      `xml:"lwtype,attr,omitempty"`
	Naamval  string      `xml:"naamval,attr,omitempty"`
	Ntype    string      `xml:"ntype,attr,omitempty"`
	Numtype  string      `xml:"numtype,attr,omitempty"`
	Pdtype   string      `xml:"pdtype,attr,omitempty"`
	Persoon  string      `xml:"persoon,attr,omitempty"`
	Postag   string      `xml:"postag,attr,omitempty"`
	Pt       string      `xml:"pt,attr,omitempty"`
	Pvagr    string      `xml:"pvagr,attr,omitempty"`
	Pvtijd   string      `xml:"pvtijd,attr,omitempty"`
	Rel      string      `xml:"rel,attr,omitempty"`
	Sc       string      `xml:"sc,attr,omitempty"`
	Spectype string      `xml:"spectype,attr,omitempty"`
	Vwtype   string      `xml:"vwtype,attr,omitempty"`
	Word     string      `xml:"word,attr,omitempty"`
	Wvorm    string      `xml:"wvorm,attr,omitempty"`
	Node     []*nodeType `xml:"node"`
	parent   *nodeType

	// als je hier iets aan toevoegt, dan ook toevoegen in emptyheads-in.go in functie reconstructEmptyHead
	udAbbr           string
	udCase           string
	udCopiedFrom     int
	udDefinite       string
	udDegree         string
	udEHeadPosition  int
	udERelation      string
	udEnhanced       string
	udFirstWordBegin int
	udForeign        string
	udGender         string
	udHeadPosition   int
	udNumber         string
	udPerson         string
	udPos            string
	udPoss           string
	udPronType       string
	udReflex         string
	udRelation       string
	udTense          string
	udVerbForm       string

	axParent            []interface{}
	axAncestors         []interface{}
	axChildren          []interface{}
	axDescendants       []interface{}
	axDescendantsOrSelf []interface{}
}

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
		udHeadPosition:      error_NO_EXTERNAL_HEAD,
		udEHeadPosition:     error_NO_EXTERNAL_HEAD,
	}
)

func init() {
	noNode.parent = noNode
}

// Derive Universal Dependencies form parsed sentence in alpino_ds format.
func Ud(alpino_doc []byte, filename string) (conllu string, err error) {

	defer func() {
		if r := recover(); r != nil {
			conllu = ""
			err = fmt.Errorf("%v", r)
		}
	}()

	conllu, err = UdTry(alpino_doc, filename)

	return
}

// Like Ud(), but may panic. Used for development.
func UdTry(alpino_doc []byte, filename string) (conllu string, err error) {

	var alpino alpino_ds
	err = xml.Unmarshal(alpino_doc, &alpino)
	if err != nil {
		return "", err
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

	q := &context{
		alpino:   &alpino,
		filename: filename,
		sentence: alpino.Sentence.Sent,
		sentid:   alpino.Sentence.SentId,
		varroot:  []interface{}{alpino.Node},
		warnings: []string{},
	}

	inspect(q)

	fixMisplacedHeadsInCoordination(q)
	addPosTags(q)
	addFeatures(q)
	addDependencyRelations(q)
	enhancedDependencies(q)
	fixpunct(q)
	untokenize(q)
	return conll(q), nil
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
