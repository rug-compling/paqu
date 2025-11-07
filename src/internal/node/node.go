package node

type Node struct {
	FullNode
	Ud       *UdType `xml:"ud,omitempty"`
	NodeList []*Node `xml:"node"`
	SkipThis bool    `xml:"-"`
}

type UdType struct {
	Id    string `xml:"id,attr,omitempty"`
	Form  string `xml:"form,attr,omitempty"`
	Lemma string `xml:"lemma,attr,omitempty"`
	Upos  string `xml:"upos,attr,omitempty"`
	FeatsType
	Head       string    `xml:"head,attr,omitempty"`
	Deprel     string    `xml:"deprel,attr,omitempty"`
	DeprelMain string    `xml:"deprel_main,attr,omitempty"`
	DeprelAux  string    `xml:"deprel_aux,attr,omitempty"`
	Dep        []DepType `xml:"dep,omitempty"`
}

type FeatsType struct {
	Abbr     string `xml:"Abbr,attr,omitempty"`
	Case     string `xml:"Case,attr,omitempty"`
	Definite string `xml:"Definite,attr,omitempty"`
	Degree   string `xml:"Degree,attr,omitempty"`
	ExtPos   string `xml:"ExtPos,attr,omitempty"`
	Foreign  string `xml:"Foreign,attr,omitempty"`
	Gender   string `xml:"Gender,attr,omitempty"`
	Mood     string `xml:"Mood,attr,omitempty"`
	Number   string `xml:"Number,attr,omitempty"`
	Person   string `xml:"Person,attr,omitempty"`
	PronType string `xml:"PronType,attr,omitempty"`
	Reflex   string `xml:"Reflex,attr,omitempty"`
	Tense    string `xml:"Tense,attr,omitempty"`
	VerbForm string `xml:"VerbForm,attr,omitempty"`
}

type DepType struct {
	Id         string `xml:"id,attr,omitempty"`
	Head       string `xml:"head,attr,omitempty"`
	Deprel     string `xml:"deprel,attr,omitempty"`
	DeprelMain string `xml:"deprel_main,attr,omitempty"`
	DeprelAux  string `xml:"deprel_aux,attr,omitempty"`
	Elided     bool   `xml:"elided,attr,omitempty"`
}
