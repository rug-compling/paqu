/*

Package alpinods implements the type AlpinoDS that can be used to
marshal and unmarshall XML files in the alpino_ds format, using the
package encoding/xml.

The field UserData in several places can be used for storing
processing information.

About Alpino, see: https://www.let.rug.nl/~vannoord/alp/Alpino/

*/
package alpinods

import (
	"encoding/xml"
	"strings"
)

// DtdVersion defines the highest supported alpino_ds.dtd version.
const DtdVersion = "1.12"

// The AlpinoDS type encodes a complete document in the alpino_ds XML
// format.
type AlpinoDS struct {
	XMLName  xml.Name  `xml:"alpino_ds"`
	Version  string    `xml:"version,attr,omitempty"`
	Metadata *Metadata `xml:"metadata,omitempty"`
	Parser   *Parser   `xml:"parser,omitempty"`
	Node     *Node     `xml:"node,omitempty"`
	Sentence *Sentence `xml:"sentence,omitempty"`
	Comments *Comments `xml:"comments,omitempty"`
	Root     []*Deprel `xml:"root,omitempty"`
	Conllu   *Conllu   `xml:"conllu,omitempty"`

	UserData interface{} `xml:"-"`
}

// The Metadata type encodes `/alpino_ds/metadata`.
type Metadata struct {
	Meta []Meta `xml:"meta,omitempty"`
}

// The Meta type encodes `/alpino_ds/metadata/meta`.
type Meta struct {
	Type  string `xml:"type,attr"`
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`

	UserData interface{} `xml:"-"`
}

// The Parser type encodes `/alpino_ds/parser`.
type Parser struct {
	Build string `xml:"build,attr,omitempty"`
	Date  string `xml:"date,attr,omitempty"`
	Cats  string `xml:"cats,attr,omitempty"`
	Skips string `xml:"skips,attr,omitempty"`
}

// The Comments type encodes `/alpino_ds/comments`.
type Comments struct {
	Comment []string `xml:"comment,omitempty"`
}

// The Sentence type encodes `/alpino_ds/sentence`.
type Sentence struct {
	Sentence string `xml:",chardata"`
	SentID   string `xml:"sentid,attr,omitempty"`
}

// The Conllu type encodes `/alpino_ds/conllu`.
type Conllu struct {
	Conllu string `xml:",cdata"`
	Status string `xml:"status,attr,omitempty"`
	Error  string `xml:"error,attr,omitempty"`
	Auto   string `xml:"auto,attr,omitempty"`
}

// The NodeAttributes type encodes the attributes for
// `/alpino_ds//node`.
type NodeAttributes struct {
	Aform        string `xml:"aform,attr,omitempty"`
	Begin        int    `xml:"begin,attr"`
	Buiging      string `xml:"buiging,attr,omitempty"`
	Case         string `xml:"case,attr,omitempty"`
	Cat          string `xml:"cat,attr,omitempty"`
	Comparative  string `xml:"comparative,attr,omitempty"`
	Conjtype     string `xml:"conjtype,attr,omitempty"`
	Def          string `xml:"def,attr,omitempty"`
	Dial         string `xml:"dial,attr,omitempty"`
	Dscmanual    string `xml:"dscmanual,attr,omitempty"`
	Dscsense     string `xml:"dscsense,attr,omitempty"`
	End          int    `xml:"end,attr"`
	Frame        string `xml:"frame,attr,omitempty"`
	Gen          string `xml:"gen,attr,omitempty"`
	Genus        string `xml:"genus,attr,omitempty"`
	Getal        string `xml:"getal,attr,omitempty"`
	GetalN       string `xml:"getal-n,attr,omitempty"`
	Graad        string `xml:"graad,attr,omitempty"`
	His          string `xml:"his,attr,omitempty"`
	His1         string `xml:"his_1,attr,omitempty"`
	His2         string `xml:"his_2,attr,omitempty"`
	His11        string `xml:"his_1_1,attr,omitempty"`
	His12        string `xml:"his_1_2,attr,omitempty"`
	His21        string `xml:"his_2_1,attr,omitempty"`
	His22        string `xml:"his_2_2,attr,omitempty"`
	His111       string `xml:"his_1_1_1,attr,omitempty"`
	His112       string `xml:"his_1_1_2,attr,omitempty"`
	His121       string `xml:"his_1_2_1,attr,omitempty"`
	His122       string `xml:"his_1_2_2,attr,omitempty"`
	His211       string `xml:"his_2_1_1,attr,omitempty"`
	His212       string `xml:"his_2_1_2,attr,omitempty"`
	His221       string `xml:"his_2_2_1,attr,omitempty"`
	His222       string `xml:"his_2_2_2,attr,omitempty"`
	ID           int    `xml:"id,attr"`
	Iets         string `xml:"iets,attr,omitempty"`
	Index        int    `xml:"index,attr,omitempty"`
	Infl         string `xml:"infl,attr,omitempty"`
	Lcat         string `xml:"lcat,attr,omitempty"`
	Lemma        string `xml:"lemma,attr,omitempty"`
	Lwtype       string `xml:"lwtype,attr,omitempty"`
	MwuRoot      string `xml:"mwu_root,attr,omitempty"`
	MwuSense     string `xml:"mwu_sense,attr,omitempty"`
	Naamval      string `xml:"naamval,attr,omitempty"`
	Neclass      string `xml:"neclass,attr,omitempty"`
	Npagr        string `xml:"npagr,attr,omitempty"`
	Ntype        string `xml:"ntype,attr,omitempty"`
	Num          string `xml:"num,attr,omitempty"`
	Numtype      string `xml:"numtype,attr,omitempty"`
	Pb           string `xml:"pb,attr,omitempty"`
	Pdtype       string `xml:"pdtype,attr,omitempty"`
	Per          string `xml:"per,attr,omitempty"`
	Personalized string `xml:"personalized,attr,omitempty"`
	Persoon      string `xml:"persoon,attr,omitempty"`
	Pos          string `xml:"pos,attr,omitempty"`
	Positie      string `xml:"positie,attr,omitempty"`
	Postag       string `xml:"postag,attr,omitempty"`
	Pron         string `xml:"pron,attr,omitempty"`
	Pt           string `xml:"pt,attr,omitempty"`
	Pvagr        string `xml:"pvagr,attr,omitempty"`
	Pvtijd       string `xml:"pvtijd,attr,omitempty"`
	Refl         string `xml:"refl,attr,omitempty"`
	Rel          string `xml:"rel,attr,omitempty"`
	Rnum         string `xml:"rnum,attr,omitempty"`
	Root         string `xml:"root,attr,omitempty"`
	Sc           string `xml:"sc,attr,omitempty"`
	Sense        string `xml:"sense,attr,omitempty"`
	SonarNe      string `xml:"sonar_ne,attr,omitempty"`
	SonarNeClass string `xml:"sonar_ne_class,attr,omitempty"`
	SonarNeBegin string `xml:"sonar_ne_begin,attr,omitempty"`
	SonarNeEnd   string `xml:"sonar_ne_end,attr,omitempty"`
	Special      string `xml:"special,attr,omitempty"`
	Spectype     string `xml:"spectype,attr,omitempty"`
	Status       string `xml:"status,attr,omitempty"`
	Stype        string `xml:"stype,attr,omitempty"`
	Tense        string `xml:"tense,attr,omitempty"`
	Vform        string `xml:"vform,attr,omitempty"`
	Vwtype       string `xml:"vwtype,attr,omitempty"`
	Vztype       string `xml:"vztype,attr,omitempty"`
	Wh           string `xml:"wh,attr,omitempty"`
	Wk           string `xml:"wk,attr,omitempty"`
	Word         string `xml:"word,attr,omitempty"`
	Wvorm        string `xml:"wvorm,attr,omitempty"`
}

// The Node type encodes `/alpino_ds//node`.
type Node struct {
	NodeAttributes
	Data []*Data `xml:"data,omitempty"`
	Ud   *Ud     `xml:"ud,omitempty"`
	Node []*Node `xml:"node"`

	UserData interface{} `xml:"-"`
}

// The Data type encodes `alpino_ds//node/data`
type Data struct {
	Name string `xml:"name,attr,omitempty"`
	Data string `xml:",chardata"`

	UserData interface{} `xml:"-"`
}

// The Ud type encodes `/alpino_ds//node/ud`.
type Ud struct {
	ID    string `xml:"id,attr,omitempty"`
	Form  string `xml:"form,attr,omitempty"`
	Lemma string `xml:"lemma,attr,omitempty"`
	Upos  string `xml:"upos,attr,omitempty"`
	Feats
	Head       string `xml:"head,attr,omitempty"`
	Deprel     string `xml:"deprel,attr,omitempty"`
	DeprelMain string `xml:"deprel_main,attr,omitempty"`
	DeprelAux  string `xml:"deprel_aux,attr,omitempty"`
	Dep        []Dep  `xml:"dep,omitempty"`

	UserData interface{} `xml:"-"`
}

// The Dep type encodes `/alpino_ds//node/ud/dep`.
type Dep struct {
	ID         string `xml:"id,attr,omitempty"`
	Head       string `xml:"head,attr,omitempty"`
	Deprel     string `xml:"deprel,attr,omitempty"`
	DeprelMain string `xml:"deprel_main,attr,omitempty"`
	DeprelAux  string `xml:"deprel_aux,attr,omitempty"`
	Elided     bool   `xml:"elided,attr,omitempty"`

	UserData interface{} `xml:"-"`
}

// The Feats type encodes the standard UD features that are used in
// Alpino.
type Feats struct {
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

// The DeprelAttributes type encodes attributes that are copied from
// `/alpino_ds//node` into `/alpino_ds/root` and its descendants.
type DeprelAttributes struct {
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
}

// The Deprel type encodes `/alpino_ds/root` and its descendants.
type Deprel struct {
	XMLName xml.Name

	RecursionLimit string `xml:"recursion_limit,attr,omitempty"`

	Ud    string `xml:"ud,attr,omitempty"`
	ID    string `xml:"id,attr,omitempty"`
	Form  string `xml:"form,attr,omitempty"`
	Lemma string `xml:"lemma,attr,omitempty"`
	Upos  string `xml:"upos,attr,omitempty"`
	Feats
	Head      string `xml:"head,attr,omitempty"`
	Deprel    string `xml:"deprel,attr,omitempty"`
	DeprelAux string `xml:"deprel_aux,attr,omitempty"`

	DeprelAttributes

	Deps []*Deprel `xml:",any"`

	UserData interface{} `xml:"-"`
}

// The String method returns the AlpinoDS type as a complete, cleaned-up, and formatted XML document.
func (a AlpinoDS) String() string {
	b, err := xml.MarshalIndent(a, "", "  ")
	if err != nil {
		panic(err) // This should never happen!
	}
	s := string(b)
	for _, a := range []string{"parser", "meta", "node", "dep", "data"} {
		s = strings.Replace(s, "></"+a+">", "/>", -1)
	}
	return "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" + s + "\n"
}
