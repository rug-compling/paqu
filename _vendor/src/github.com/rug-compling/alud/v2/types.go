package alud

import (
	"github.com/rug-compling/alpinods"

	"encoding/xml"
)

type context struct {
	alpino        *alpino_ds
	filename      string
	sentence      string
	sentid        string
	debugs        []string
	depth         int
	allnodes      []*nodeType
	ptnodes       []*nodeType
	varallnodes   []interface{}
	varindexnodes []interface{}
	varptnodes    []interface{}
	varroot       []interface{}
	swapped       [][2]*nodeType
}

type traceT struct {
	msg    string
	debugs []string
	trace  []traceType
}

type traceType struct {
	s    string
	node *nodeType
	head *nodeType
	gap  *nodeType
	subj *nodeType
}

type alpino_ds struct {
	XMLName  xml.Name      `xml:"alpino_ds"`
	Version  string        `xml:"version,attr,omitempty"`
	Metadata *metadataType `xml:"metadata,omitempty"`
	Parser   *parserType   `xml:"parser,omitempty"`
	Node     *nodeType     `xml:"node,omitempty"`
	Sentence *sentType     `xml:"sentence,omitempty"`
	Comments *commentsType `xml:"comments,omitempty"`
	UdNodes  []*udNodeType `xml:"root,omitempty"`
	Conllu   *conlluType   `xml:"conllu,omitempty"`
}

type metadataType struct {
	Meta []metaType `xml:"meta,omitempty"`
}

type commentsType struct {
	Comment []string `xml:"comment,omitempty"`
}

type sentType struct {
	Sent   string `xml:",chardata"`
	SentId string `xml:"sentid,attr,omitempty"`
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

type depType struct {
	Id         string `xml:"id,attr,omitempty"`
	Head       string `xml:"head,attr,omitempty"`
	Deprel     string `xml:"deprel,attr,omitempty"`
	DeprelMain string `xml:"deprel_main,attr,omitempty"`
	DeprelAux  string `xml:"deprel_aux,attr,omitempty"`
	Elided     bool   `xml:"elided,attr,omitempty"`
}

type conlluType struct {
	Conllu string `xml:",cdata"`
	Status string `xml:"status,attr,omitempty"`
	Error  string `xml:"error,attr,omitempty"`
	Auto   string `xml:"auto,attr,omitempty"`
}

type nodeType struct {
	alpinods.NodeAttributes
	Data   []*Data     `xml:"data,omitempty"`
	Node   []*nodeType `xml:"node"`
	parent *nodeType

	Ud *udType `xml:"ud,omitempty"`

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
	udNoSpaceAfter   bool
	udNumber         string
	udOldState       *nodeType
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

type Data struct {
	Name string `xml:"name,attr,omitempty"`
	Data string `xml:",chardata"`
}

type udNodeType struct {
	XMLName xml.Name

	RecursionLimit string `xml:"recursion_limit,attr,omitempty"`
	recursion      []string

	Ud string `xml:"ud,attr,omitempty"`

	Id    string `xml:"id,attr,omitempty"`
	Form  string `xml:"form,attr,omitempty"`
	Lemma string `xml:"lemma,attr,omitempty"`
	Upos  string `xml:"upos,attr,omitempty"`
	alpinods.Feats
	Head      string `xml:"head,attr,omitempty"`
	Deprel    string `xml:"deprel,attr,omitempty"`
	DeprelAux string `xml:"deprel_aux,attr,omitempty"`

	alpinods.DeprelAttributes

	UdNodes []*udNodeType `xml:",omitempty"`
}

type udType struct {
	Id    string `xml:"id,attr,omitempty"`
	Form  string `xml:"form,attr,omitempty"`
	Lemma string `xml:"lemma,attr,omitempty"`
	Upos  string `xml:"upos,attr,omitempty"`
	alpinods.Feats
	Head       string    `xml:"head,attr,omitempty"`
	Deprel     string    `xml:"deprel,attr,omitempty"`
	DeprelMain string    `xml:"deprel_main,attr,omitempty"`
	DeprelAux  string    `xml:"deprel_aux,attr,omitempty"`
	Dep        []depType `xml:"dep,omitempty"`
}

type nodeDepType struct {
	node *nodeType
	dep  *depType
}
