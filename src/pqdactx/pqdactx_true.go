// +build !nodbxml

package main

import (
	pqnode "github.com/rug-compling/paqu/internal/node"

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
	Node     *pqnode.Node  `xml:"node,omitempty"`
	Sentence *SentType     `xml:"sentence,omitempty"`
	Comments *CommentsType `xml:"comments,omitempty"`
	Root     []*UdNodeType `xml:"root,omitempty"`
	Conllu   *ConlluType   `xml:"conllu,omitempty"`
}

type UdNodeType struct {
	XMLName xml.Name

	RecursionLimit string `xml:"recursion_limit,attr,omitempty"`

	Ud    string `xml:"ud,attr,omitempty"`
	Id    string `xml:"id,attr,omitempty"`
	Eid   string `xml:"eid,attr,omitempty"`
	Form  string `xml:"form,attr,omitempty"`
	Lemma string `xml:"lemma,attr,omitempty"`
	Upos  string `xml:"upos,attr,omitempty"`
	pqnode.FeatsType
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

	UdNodes []byte `xml:",innerxml"`
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
	Build string `xml:"build,attr,omitempty"`
	Date  string `xml:"date,attr,omitempty"`
	Cats  string `xml:"cats,attr,omitempty"`
	Skips string `xml:"skips,attr,omitempty"`
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
	size, _ := db1.Size()
	db2, err := dbxml.OpenReadWrite(os.Args[2])
	x(err)
	docs, err := db1.All()
	x(err)
	teller := 0
	for docs.Next() {
		teller++
		name := docs.Name()
		fmt.Printf("\r\033[K[%d/%d] %s ", teller, size, name)
		content := docs.Content()
		alpino := Alpino_ds{}
		err = xml.Unmarshal([]byte(content), &alpino)
		x(err)
		if expand(&alpino) {
			content = format(alpino)
		}
		x(db2.PutXml(name, content, false))
	}
	fmt.Print("\r\033[K")
	db2.Close()
	db1.Close()
	st, err := os.Stat(os.Args[1])
	x(err)
	x(os.Chmod(os.Args[2], st.Mode()))
}

func expand(alpino *Alpino_ds) bool {
	refs := make(map[string]*pqnode.Node)
	getIndexed(alpino.Node, refs)
	if len(refs) == 0 {
		return false
	}
	alpino.Version = "X-" + alpino.Version
	expandNode(alpino.Node, refs)
	return true
}

func getIndexed(node *pqnode.Node, nodes map[string]*pqnode.Node) {
	if node.Index != "" && (node.NodeList != nil || node.Word != "") {
		nodes[node.Index] = node
	}
	if node.NodeList != nil {
		for _, n := range node.NodeList {
			getIndexed(n, nodes)
		}
	}
}

func expandNode(n *pqnode.Node, nodes map[string]*pqnode.Node) {
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

	pqnode.CopyNodeOnEmpty(n, o)
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
