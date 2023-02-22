//go:build ignore
// +build ignore

package main

import (
	"github.com/rug-compling/alpinods"

	"fmt"
	"sort"
	"strings"
)

func main() {
	attribs := []string{
		"other_id",
		"word_is_",
	}
	for _, a := range alpinods.Attribs() {
		attribs = append(attribs, a.Name)
	}
	sort.Strings(attribs)

	fmt.Print(`//
// THIS IS A GENERATED FILE. DO NOT EDIT.
//

package node

import (
	"fmt"
)

type FullNode struct {
`)
	for _, attrib := range attribs {
		up := upper(attrib)
		typ := "string"
		omit := ",omitempty"
		if attrib == "begin" || attrib == "end" {
			typ = "int   "
			omit = ""
		}
		fmt.Printf("\t%-12s %s `xml:\"%s,attr%s\"`\n", up, typ, attrib, omit)
	}
	fmt.Print(`}

var NodeTags = []string{
`)
	for _, attrib := range attribs {
		if !strings.HasSuffix(attrib, "_") {
			fmt.Printf("\t\"%s\",\n", attrib)
		}
	}
	fmt.Print(`}

func GetAttr(attr string, n *FullNode) string {
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

func CopyNodeOnEmpty(dst, src *Node) {
`)
	for _, attrib := range attribs {
		up := upper(attrib)
		null := `""`
		if attrib == "begin" || attrib == "end" {
			null = `0`
		}
		fmt.Printf("\tif dst.%s == %s {\n\t\tdst.%s = src.%s\n\t}\n", up, null, up, up)
	}
	fmt.Println(`
	if dst.NodeList == nil {
		dst.NodeList = src.NodeList
	}
	if dst.Ud == nil || dst.Ud.Id == "" {
		dst.Ud = src.Ud
	}
}
`)
}

func upper(s string) string {
	return strings.Replace(strings.Title(strings.Replace(strings.Replace(s, "_", " ", -1), "-", " ", -1)), " ", "", -1)
}
