// +build !nodbxml

package main

import (
	"github.com/pebbe/dbxml"
	"github.com/pebbe/util"

	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

type Alpino_ds struct {
	XMLName  xml.Name      `xml:"alpino_ds"`
	Version  string        `xml:"version,attr,omitempty"`
	Metadata *MetadataType `xml:"metadata,omitempty"`
	Parser   *ParserType   `xml:"parser,omitempty"`
	Node     *Node         `xml:"node,omitempty"`
	Sentence *SentType     `xml:"sentence,omitempty"`
	Comments *CommentsType `xml:"comments,omitempty"`
	Conllu   *ConlluType   `xml:"conllu,omitempty"`
}

type MetadataType struct {
	Meta []MetaType `xml:"meta,omitempty"`
}

type CommentsType struct {
	Comment []string `xml:"comment,omitempty"`
}

type SentType struct {
	Sent   string `xml:",chardata"`
	SentId string `xml:"sentid,attr,omitempty"`
}

type MetaType struct {
	Type  string `xml:"type,attr,omitempty"`
	Name  string `xml:"name,attr,omitempty"`
	Value string `xml:"value,attr,omitempty"`
}

type ParserType struct {
	Cats  string `xml:"cats,attr,omitempty"`
	Skips string `xml:"skips,attr,omitempty"`
}

type Node struct {
	FullNode
	Ud       *UdType `xml:"ud,omitempty"`
	NodeList []*Node `xml:"node"`
}

type UdType struct {
	Id    string `xml:"id,attr,omitempty"`
	Form  string `xml:"form,attr,omitempty"`
	Lemma string `xml:"lemma,attr,omitempty"`
	Upos  string `xml:"upos,attr,omitempty"`
	Xpos  string `xml:"xpos,attr,omitempty"`
	FeatsType
	Head   string    `xml:"head,attr,omitempty"`
	Deprel string    `xml:"deprel,attr,omitempty"`
	Dep    []DepType `xml:"dep,omitempty"`
	Misc   string    `xml:"misc,attr,omitempty"`
}

type FeatsType struct {
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

type DepType struct {
	Id     string `xml:"id,attr,omitempty"`
	Head   string `xml:"head,attr,omitempty"`
	Deprel string `xml:"deprel,attr,omitempty"`
}

type ConlluType struct {
	Conllu string `xml:",cdata"`
	Status string `xml:"status,attr,omitempty"`
	Error  string `xml:"error,attr,omitempty"`
	Auto   string `xml:"auto,attr,omitempty"`
}

var (
	x = util.CheckErr
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf(`
Syntax: %s infile.dact outfile.dactx

`, os.Args[0])
		return
	}
	db1, err := dbxml.OpenRead(os.Args[1])
	x(err)
	db2, err := dbxml.OpenReadWrite(os.Args[2])
	x(err)
	docs, err := db1.All()
	x(err)
	for docs.Next() {
		name := docs.Name()
		fmt.Fprintln(os.Stderr, name)
		content := docs.Content()
		alpino := Alpino_ds{}
		err = xml.Unmarshal([]byte(content), &alpino)
		x(err)
		if expand(&alpino) {
			content = format(alpino)
		}
		x(db2.PutXml(name, content, false))
	}
	db2.Close()
	db1.Close()
	st, err := os.Stat(os.Args[1])
	x(err)
	x(os.Chmod(os.Args[2], st.Mode()))
}

func expand(alpino *Alpino_ds) bool {
	refs := make(map[string]*Node)
	getIndexed(alpino.Node, refs)
	if len(refs) == 0 {
		return false
	}
	alpino.Version = "X-" + alpino.Version
	expandNode(alpino.Node, refs)
	return true
}

func getIndexed(node *Node, nodes map[string]*Node) {
	if node.Index != "" && (node.NodeList != nil || node.Word != "") {
		nodes[node.Index] = node
	}
	if node.NodeList != nil {
		for _, n := range node.NodeList {
			getIndexed(n, nodes)
		}
	}
}

func expandNode(n *Node, nodes map[string]*Node) {
	if n.NodeList != nil {
		for _, node := range n.NodeList {
			expandNode(node, nodes)
		}
	}

	if n.Index == "" || n.NodeList != nil || n.Word != "" {
		return
	}

	o, ok := nodes[n.Index]
	if !ok {
		fmt.Fprintln(os.Stderr, "Missing node")
		return
	}

	n.OtherId = o.Id
	n.Ud = o.Ud

	copyNodeOnEmpty(n, o)
}

func format(alpino Alpino_ds) string {
	b, err := xml.MarshalIndent(&alpino, "", "  ")
	x(err)
	s := "<?xml version=\"1.0\"?>\n" + string(b)

	// shorten
	for _, v := range []string{"meta", "parser", "node", "dep"} {
		s = strings.Replace(s, "></"+v+">", "/>", -1)
	}

	return s
}
