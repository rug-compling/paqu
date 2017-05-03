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
	XMLName  xml.Name `xml:"alpino_ds"`
	Version  string   `xml:"version,attr,omitempty"`
	Metadata []MetaT  `xml:"metadata>meta,omitempty"`
	Node0    *Node    `xml:"node,omitempty"`
	Sentence SentT    `xml:"sentence,omitempty"`
	Comments []string `xml:"comments>comment,omitempty"`
}

type SentT struct {
	Sent   string `xml:",chardata"`
	SentId string `xml:"sentid,attr,omitempty"`
}

type MetaT struct {
	Type  string `xml:"type,attr,omitempty"`
	Name  string `xml:"name,attr,omitempty"`
	Value string `xml:"value,attr,omitempty"`
}

type Node struct {
	FullNode
	NodeList []*Node `xml:"node"`
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
	db1, err := dbxml.Open(os.Args[1])
	x(err)
	db2, err := dbxml.Open(os.Args[2])
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
			b, err := xml.MarshalIndent(&alpino, "", "  ")
			x(err)
			content = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
				strings.Replace(
					strings.Replace(string(b), "  <metadata></metadata>\n", "", 1),
					"  <comments></comments>\n", "", 1) + "\n"
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
	getIndexed(alpino.Node0, refs)
	if len(refs) == 0 {
		return false
	}
	alpino.Version = "X-" + alpino.Version
	expandNode(alpino.Node0, refs)
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

	copyNodeOnEmpty(n, o)
}
