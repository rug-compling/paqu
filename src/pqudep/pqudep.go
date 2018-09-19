package main

/*
#cgo LDFLAGS: -lxqilla
#include <xqilla/xqilla-xqc.h>
#include <stdlib.h>

XQC_Implementation *impl;
XQC_Expression *expr;
XQC_DynamicContext *context;
XQC_Sequence *seq, *doc;
XQC_Error err;
const char *value;
int done;

int init(char const *xquery)
{
  // XQilla specific way to create an XQC_Implementation struct
  impl = createXQillaXQCImplementation(XQC_VERSION_NUMBER);
  if(impl == 0) return 1;

  // Parse an XQuery expression
  err = impl->prepare(impl, xquery, 0, &expr);
  if(err != 0) return err;

  return 0;
}

int parse(char const *xml) {
  // Parse a document
  err = impl->parse_document(impl, xml, &doc);
  if(err != 0) return err;

  // Create a dynamic context
  err = expr->create_context(expr, &context);
  if(err != 0) return err;

  // Set the document as the context item
  doc->next(doc);
  context->set_context_item(context, doc);

  // Execute the query
  err = expr->execute(expr, context, &seq);
  if(err != 0) return err;

  return 0;
}

int next() {
  done = 1;
  if((err = seq->next(seq)) == XQC_NO_ERROR) {
    seq->string_value(seq, &value);
    done = 0;
  }

  if(err == XQC_END_OF_SEQUENCE)
    err = XQC_NO_ERROR;

  if (done) {
    seq->free(seq);
    context->free(context);
    doc->free(doc);
  }

  return err;
}

*/
import "C"

import (
	"github.com/pebbe/dbxml"
	"github.com/pebbe/util"

	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

type Alpino_ds struct {
	XMLName  xml.Name      `xml:"alpino_ds"`
	Version  string        `xml:"version,attr,omitempty"`
	Metadata *MetadataType `xml:"metadata,omitempty"`
	Parser   *ParserType   `xml:"parser,omitempty"`
	Node     *NodeType     `xml:"node,omitempty"`
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

type NodeType struct {
	FullNode
	Ud       *UdType     `xml:"ud,omitempty"`
	NodeList []*NodeType `xml:"node"`
	skip     bool
}

type UdType struct {
	Id    string `xml:"id,attr,omitempty"`
	Form  string `xml:"form,attr,omitempty"`
	Lemma string `xml:"lemma,attr,omitempty"`
	Upos  string `xml:"upos,attr,omitempty"`
	Xpos  string `xml:"xpos,attr,omitempty"`
	FeatsType
	Head   string     `xml:"head,attr,omitempty"`
	Deprel string     `xml:"deprel,attr,omitempty"`
	Deps   []DepsType `xml:"deps,omitempty"`
	Misc   string     `xml:"misc,attr,omitempty"`
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

type DepsType struct {
	Id     string `xml:"id,attr,omitempty"`
	Head   string `xml:"head,attr,omitempty"`
	Deprel string `xml:"deprel,attr,omitempty"`
}

type ConlluType struct {
	Conllu string `xml:",cdata"`
	Status string `xml:"status,attr,omitempty"`
	Error  string `xml:"error,attr,omitempty"`
}

var (
	x     = util.CheckErr
	opt_l = flag.String("l", "", "filelist")
)

func usage() {
	fmt.Printf(`
Usage: %s [-l filelist] [file.(xml|dact)]...

  -l: file with list of names of xml and/or dact files
`)
}

func main() {

	flag.Usage = usage
	flag.Parse()

	if flag.NArg() == 0 && *opt_l == "" {
		usage()
		return
	}

	cs := C.CString(udep)
	if e := C.init(cs); e != 0 {
		x(fmt.Errorf("C.init: %d", e))
	}
	C.free(unsafe.Pointer(cs))

	if *opt_l != "" {
		fp, err := os.Open(*opt_l)
		x(err)
		scanner := bufio.NewScanner(fp)
		for scanner.Scan() {
			filename := strings.TrimSpace(scanner.Text())
			if filename != "" {
				doFile(filename)
			}
		}
		x(scanner.Err())
	}

	for _, filename := range flag.Args() {
		doFile(filename)
	}
	fmt.Println()
}

func doFile(filename string) {
	if strings.HasSuffix(filename, ".dact") {
		db1, err := dbxml.OpenRead(filename)
		x(err)
		db2, err := dbxml.OpenReadWrite(filename + ".tmp")
		x(err)
		docs, err := db1.All()
		for docs.Next() {
			f := docs.Name()
			xml := docs.Content()
			fmt.Printf("%s - %-8s\r", filename, f)
			result, err := doXml(xml, f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s - %s: %v\n", filename, f, err)
			}
			x(db2.PutXml(f, result, false))
		}
		db2.Close()
		db1.Close()
	} else {
		fmt.Printf("%s%8s\r", filename, "")
		b, err := ioutil.ReadFile(filename)
		x(err)
		result, err := doXml(string(b), filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", filename, err)
		}
		fp, err := os.Create(filename + ".tmp")
		x(err)
		fmt.Fprintln(fp, result)
		fp.Close()
	}
}

func doXml(document, filename string) (result string, err error) {

	var alpino Alpino_ds
	lines := make([]string, 0)

	defer func() {
		alpino.Conllu.Conllu = "\n" + strings.Join(lines, "\n") + "\n"
		if err == nil {
			alpino.Conllu.Status = "OK"
		} else {
			// alpino.Conllu.Conllu = "\n" + strings.Join(lines, "\n") + "\n"
			alpino.Conllu.Status = "error"
			alpino.Conllu.Error = err.Error()
			clean(alpino.Node)
		}
		minify(alpino)
		result = format(alpino)
	}()

	err = xml.Unmarshal([]byte(document), &alpino)
	if err != nil {
		return
	}
	reset(alpino.Node)
	alpino.Conllu = &ConlluType{}

	cs := C.CString(document)
	e := C.parse(cs)
	C.free(unsafe.Pointer(cs))
	if e != 0 {
		err = fmt.Errorf("C.parse %s: %d", filename, e)
		return
	}
	for {
		if e := C.next(); e != 0 {
			err = fmt.Errorf("C.next %s: %d", filename, e)
			return
		}
		if C.done != 0 {
			break
		}
		lines = append(lines, C.GoString(C.value))
	}

	valid := make(map[string]bool)
	valid["0"] = true // root
	copies := make(map[string]string)
	for i, line := range lines {
		a := strings.Split(line, "\t")
		if len(a) != 10 {
			err = fmt.Errorf("Line %d: Wrong number of fields", i+1)
			return
		}
		valid[a[0]] = true
		items := getItems(a[9])
		if n, ok := items["CopiedFrom"]; ok {
			copies[n] = a[0]
		}
	}

	for i, line := range lines {
		a := strings.Split(line, "\t")
		if strings.Contains(a[0], "-") {
			continue
		}

		items := getItems(a[9])
		es, isCopy := items["CopiedFrom"]
		if !isCopy {
			es = a[0]
		}
		e, _ := strconv.Atoi(es)
		node := getNode(alpino.Node, e)
		if node == nil {
			err = fmt.Errorf("Line %d: Node '%s' not found", i+1, es)
			return
		}

		if a[8] != "_" {
			node.Ud.Deps = make([]DepsType, 0)
			for _, deps := range strings.Split(a[8], "|") {
				b := strings.SplitN(deps, ":", 2)
				if len(b) != 2 {
					err = fmt.Errorf("Line %d: Not a valid dependency: %s", i+1, deps)
					return
				}
				if !valid[b[0]] {
					err = fmt.Errorf("Line %d: Not a valid head: %s", i+1, b[0])
					return
				}
				node.Ud.Deps = append(node.Ud.Deps, DepsType{
					Id:     a[0],
					Head:   b[0],
					Deprel: b[1],
				})
			}
		}

		if isCopy {
			continue
		}

		if !valid[a[6]] {
			err = fmt.Errorf("Line %d: Not a valid head: %s", i+1, a[6])
			return
		}

		node.Ud.Id = a[0]
		node.Ud.Form = a[1]
		node.Ud.Lemma = a[2]
		node.Ud.Upos = a[3]
		node.Ud.Xpos = a[4]
		node.Ud.Head = a[6]
		node.Ud.Deprel = a[7]

		feats := getItems(a[5])
		node.Ud.Abbr = feats["Abbr"]
		node.Ud.Case = feats["Case"]
		node.Ud.Definite = feats["Definite"]
		node.Ud.Degree = feats["Degree"]
		node.Ud.Foreign = feats["Foreign"]
		node.Ud.Gender = feats["Gender"]
		node.Ud.Number = feats["Number"]
		node.Ud.Person = feats["Person"]
		node.Ud.PronType = feats["PronType"]
		node.Ud.Reflex = feats["Reflex"]
		node.Ud.Tense = feats["Tense"]
		node.Ud.VerbForm = feats["VerbForm"]

	}

	return
}

func format(alpino Alpino_ds) string {
	b, err := xml.MarshalIndent(&alpino, "", "  ")
	x(err)
	s := "<?xml version=\"1.0\"?>\n" + string(b)

	// shorten
	for _, v := range []string{"meta", "parser", "node", "deps"} {
		s = strings.Replace(s, "></"+v+">", "/>", -1)
	}

	// namespace
	s = strings.Replace(s, "<alpino_ds", "<alpino_ds xmlns:ud=\"http://www.let.rug.nl/alfa/unidep/\"", 1)
	for _, v := range []string{"ud", "deps", "conllu"} {
		s = strings.Replace(s, "<"+v, "<ud:"+v, -1)
		s = strings.Replace(s, "</"+v, "</ud:"+v, -1)
	}

	return s
}

func getNode(node *NodeType, end int) *NodeType {
	if node == nil {
		return nil
	}
	if node.End == end && node.Word != "" {
		return node
	}
	if node.NodeList != nil {
		for _, n := range node.NodeList {
			if n2 := getNode(n, end); n2 != nil {
				return n2
			}
		}
	}
	return nil
}

func getItems(s string) map[string]string {
	m := make(map[string]string)
	if s == "_" {
		return m
	}
	for _, item := range strings.Split(s, "|") {
		a := strings.SplitN(item, "=", 2)
		if len(a) == 2 {
			m[strings.TrimSpace(a[0])] = strings.TrimSpace(a[1])
		}
	}
	return m
}

/*
  bedoeling:
  - zorg ervoor dat ",omitempty" werkt
*/
func minify(alpino Alpino_ds) {
	if alpino.Metadata != nil && (alpino.Metadata.Meta == nil || len(alpino.Metadata.Meta) == 0) {
		alpino.Metadata = nil
	}
	if alpino.Parser != nil && alpino.Parser.Cats == "" && alpino.Parser.Skips == "" {
		alpino.Parser = nil
	}
	if alpino.Sentence != nil && alpino.Sentence.Sent == "" && alpino.Sentence.SentId == "" {
		alpino.Sentence = nil
	}
	if alpino.Comments != nil && (alpino.Comments.Comment == nil || len(alpino.Comments.Comment) == 0) {
		alpino.Comments = nil
	}
	minifyNode(alpino.Node)
}

func minifyNode(node *NodeType) {
	if node == nil {
		return
	}
	if node.Ud != nil {
		if node.Ud.Id == "" {
			node.Ud = nil
		} else {
			if len(node.Ud.Deps) == 0 {
				node.Ud.Deps = nil
			}
		}
	}
	if node.NodeList != nil {
		for _, n := range node.NodeList {
			minifyNode(n)
		}
	}
}

/*
 bedoeling:
 - na fout, verwijder alle UD-informatie
*/
func clean(node *NodeType) {
	node.Ud = nil
	for _, n := range node.NodeList {
		clean(n)
	}
}

/*
  bedoeling:
  - verwijder eventuele oude gegevens
  - init als nog niet aanwezig
*/
func reset(node *NodeType) {
	node.Ud = &UdType{Deps: make([]DepsType, 0)}
	for _, n := range node.NodeList {
		reset(n)
	}
}
