package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	attribs := strings.Fields(string(data))

	fmt.Print(`//
// THIS IS A GENERATED FILE. DO NOT EDIT.
//

package main

import (
	"fmt"
)

type FullNode struct {
`)
	for _, attrib := range attribs {
		up := upper(attrib)
		typ := "string"
		if attrib == "begin" || attrib == "end" {
			typ = "int   "
		}
		fmt.Printf("\t%-12s %s `xml:\"%s,attr\"`\n", up, typ, attrib)
	}
	fmt.Print(`}

var NodeTags = []string{
`)
	for _, attrib := range attribs {
		fmt.Printf("\t\"%s\",\n", attrib)
	}
	fmt.Print(`}

func getAttr(attr string, n *FullNode) string {
	switch attr {
`)
	for _, attrib := range attribs {
		up := upper(attrib)
		fmt.Printf("\tcase \"%s\":\n\t\treturn ", attrib)
		if attrib == "begin" || attrib == "end" {
			fmt.Printf("fmt.Sprint(n.%s)\n", up)
		} else {
			fmt.Printf("n.%s\n", up)
		}

	}
	fmt.Print(`	}
	return ""
}
`)
}

func upper(s string) string {
	return strings.Replace(strings.Title(strings.Replace(strings.Replace(s, "_", " ", -1), "-", " ", -1)), " ", "", -1)
}
