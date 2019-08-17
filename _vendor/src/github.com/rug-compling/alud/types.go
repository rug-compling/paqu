package alud

import (
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

type trace struct {
	s    string
	node *nodeType
	head *nodeType
	gap  *nodeType
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
type featsType struct {
	Abbr     string `xml:"Abbr,attr,omitempty"`
	Case     string `xml:"Case,attr,omitempty"`
	Definite string `xml:"Definite,attr,omitempty"`
	Degree   string `xml:"Degree,attr,omitempty"`
	Foreign  string `xml:"Foreign,attr,omitempty"`
	Gender   string `xml:"Gender,attr,omitempty"`
	Number   string `xml:"Number,attr,omitempty"`
	Person   string `xml:"Person,attr,omitempty"`
	PronType string `xml:"PronType,attr,omitempty"`
	Reflex   string `xml:"Reflex,attr,omitempty"`
	Tense    string `xml:"Tense,attr,omitempty"`
	VerbForm string `xml:"VerbForm,attr,omitempty"`
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
	Begin  int    `xml:"begin,attr"`
	End    int    `xml:"end,attr"`
	Id     int    `xml:"id,attr,omitempty"`
	Index  int    `xml:"index,attr,omitempty"`
	Lemma  string `xml:"lemma,attr,omitempty"`
	Postag string `xml:"postag,attr,omitempty"`
	Pt     string `xml:"pt,attr,omitempty"`
	Rel    string `xml:"rel,attr,omitempty"`
	Word   string `xml:"word,attr,omitempty"`
	fullNode
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

type udNodeType struct {
	XMLName xml.Name

	RecursionLimit string `xml:"recursion_limit,attr,omitempty"`
	recursion      []string

	Ud string `xml:"ud,attr,omitempty"`

	Id    string `xml:"id,attr,omitempty"`
	Form  string `xml:"form,attr,omitempty"`
	Lemma string `xml:"lemma,attr,omitempty"`
	Upos  string `xml:"upos,attr,omitempty"`
	featsType
	Head      string `xml:"head,attr,omitempty"`
	Deprel    string `xml:"deprel,attr,omitempty"`
	DeprelAux string `xml:"deprel_aux,attr,omitempty"`

	Buiging  string `xml:"buiging,attr,omitempty"`
	Conjtype string `xml:"conjtype,attr,omitempty"`
	Dial     string `xml:"dial,attr,omitempty"`
	Genus    string `xml:"genus,attr,omitempty"`
	Getal    string `xml:"getal,attr,omitempty"`
	GetalN   string `xml:"getal-n,attr,omitempty"`
	Graad    string `xml:"graad,attr,omitempty"`
	Lwtype   string `xml:"lwtype,attr,omitempty"`
	Naamval  string `xml:"naamval,attr,omitempty"`
	Npagr    string `xml:"npagr,attr,omitempty"`
	Ntype    string `xml:"ntype,attr,omitempty"`
	Numtype  string `xml:"numtype,attr,omitempty"`
	Pdtype   string `xml:"pdtype,attr,omitempty"`
	Persoon  string `xml:"persoon,attr,omitempty"`
	Positie  string `xml:"positie,attr,omitempty"`
	Pt       string `xml:"pt,attr,omitempty"`
	Pvagr    string `xml:"pvagr,attr,omitempty"`
	Pvtijd   string `xml:"pvtijd,attr,omitempty"`
	Spectype string `xml:"spectype,attr,omitempty"`
	Status   string `xml:"status,attr,omitempty"`
	Vwtype   string `xml:"vwtype,attr,omitempty"`
	Vztype   string `xml:"vztype,attr,omitempty"`
	Wvorm    string `xml:"wvorm,attr,omitempty"`

	UdNodes []*udNodeType `xml:",omitempty"`
}

type udType struct {
	Id    string `xml:"id,attr,omitempty"`
	Form  string `xml:"form,attr,omitempty"`
	Lemma string `xml:"lemma,attr,omitempty"`
	Upos  string `xml:"upos,attr,omitempty"`
	featsType
	Head       string    `xml:"head,attr,omitempty"`
	Deprel     string    `xml:"deprel,attr,omitempty"`
	DeprelMain string    `xml:"deprel_main,attr,omitempty"`
	DeprelAux  string    `xml:"deprel_aux,attr,omitempty"`
	Dep        []depType `xml:"dep,omitempty"`
}
