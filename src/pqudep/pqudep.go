package main

/*
#cgo LDFLAGS: -lxqilla
#include <xqilla/xqilla-xqc.h>
#include <stdlib.h>
#include <string.h>

XQC_Implementation *impl;
XQC_Expression *expr;
XQC_DynamicContext *context;
XQC_StaticContext *static_context;
XQC_Sequence *seq, *doc;
XQC_Error err;
const char *value;
char *error_value;
int error_len;
int done;

void my_handler(XQC_ErrorHandler *handler,
                XQC_Error error,
                const char *error_uri,
                const char *error_localname,
                const char *description,
                XQC_Sequence *error_object)
{
    int n;
    if (!description) {
      description = "unknown error";
    }
    n = strlen(description) + 1;
    if (n > error_len) {
      error_value = (char *) realloc (error_value, n * sizeof (char));
      error_len = n;
    }
    strcpy (error_value, description);
}

XQC_ErrorHandler my_handler_s = {
    error: my_handler
};

int init(char const *xquery)
{
  error_value = (char *) malloc (sizeof (char));
  error_value[0] = '\0';
  error_len = 1;

  // XQilla specific way to create an XQC_Implementation struct
  impl = createXQillaXQCImplementation(XQC_VERSION_NUMBER);
  if(impl == 0) return 1;

  err = impl->create_context(impl, &static_context);
  if(err != 0) return err;

  static_context->set_error_handler(static_context, &my_handler_s);

  // Parse an XQuery expression
  err = impl->prepare(impl, xquery, static_context, &expr);
  if(err != 0) return err;

  return 0;
}

int parse(char const *xml) {
  XQC_Sequence
     *value;
  char const
     *lib[] = { "lib" },
     *any[] = { "_" },
     *yes[] = { "yes" };

  // Parse a document
  err = impl->parse_document(impl, xml, &doc);
  if(err != 0) return err;

  // Create a dynamic context
  err = expr->create_context(expr, &context);
  if(err != 0) return err;

  context->set_error_handler(context, &my_handler_s);

  impl->create_string_sequence(impl, any, 1, &value);
  if(err != 0) return err;
  err = context->set_variable(context, "", "DIR", value);
  if(err != 0) return err;

  impl->create_string_sequence(impl, lib, 1, &value);
  if(err != 0) return err;
  err = context->set_variable(context, "", "MODE", value);
  if(err != 0) return err;

  impl->create_string_sequence(impl, yes, 1, &value);
  if(err != 0) return err;
  err = context->set_variable(context, "", "ENHANCED", value);
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
	"compress/gzip"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

const (
	VERSIONs        = "PQU%d.%d"
	VERSIONxq       = int(2) // ophogen als xquery-script veranderd is, en dan de volgende resetten naar 0
	VERSIONxml      = int(0) // ophogen als xml-formaat is veranderd
	ALPINO_DS_MAJOR = int(1)
	ALPINO_DS_MINOR = int(9)
)

type Alpino_ds struct {
	XMLName  xml.Name      `xml:"alpino_ds"`
	Version  string        `xml:"version,attr,omitempty"`
	Metadata *MetadataType `xml:"metadata,omitempty"`
	Parser   *ParserType   `xml:"parser,omitempty"`
	Node     *NodeType     `xml:"node,omitempty"`
	Sentence *SentType     `xml:"sentence,omitempty"`
	Comments *CommentsType `xml:"comments,omitempty"`
	Root     []*UdNodeType `xml:"root,omitempty"`
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
	Build string `xml:"build,attr,omitempty"`
	Date  string `xml:"date,attr,omitempty"`
	Cats  string `xml:"cats,attr,omitempty"`
	Skips string `xml:"skips,attr,omitempty"`
}

type NodeType struct {
	FullNode
	Ud       *UdType     `xml:"ud,omitempty"`
	NodeList []*NodeType `xml:"node"`
}

type UdNodeType struct {
	XMLName xml.Name

	RecursionLimit string `xml:"recursion_limit,attr,omitempty"`
	recursion      []string

	Enhanced bool `xml:"enhanced,attr,omitempty"`

	Id    string `xml:"id,attr,omitempty"`
	Form  string `xml:"form,attr,omitempty"`
	Lemma string `xml:"lemma,attr,omitempty"`
	Upos  string `xml:"upos,attr,omitempty"`
	FeatsType
	Head       string `xml:"head,attr,omitempty"`
	Deprel     string `xml:"deprel,attr,omitempty"`
	DeprelMain string `xml:"deprel_main,attr,omitempty"`
	DeprelAux  string `xml:"deprel_aux,attr,omitempty"`

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

	UdNodes []*UdNodeType `xml:",omitempty"`
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
	Id         string `xml:"id,attr,omitempty"`
	Head       string `xml:"head,attr,omitempty"`
	Deprel     string `xml:"deprel,attr,omitempty"`
	DeprelMain string `xml:"deprel_main,attr,omitempty"`
	DeprelAux  string `xml:"deprel_aux,attr,omitempty"`
	Elided     bool   `xml:"elided,attr,omitempty"`
}

type ConlluType struct {
	Conllu string `xml:",cdata"`
	Status string `xml:"status,attr,omitempty"`
	Error  string `xml:"error,attr,omitempty"`
	Auto   string `xml:"auto,attr,omitempty"`
}

var (
	x     = util.CheckErr
	opt_i = flag.Bool("i", false, "ignore existing cdata from <conllu>")
	opt_k = flag.Bool("k", false, "keep existing cdata from <conllu>")
	opt_l = flag.String("l", "", "filelist")
	opt_o = flag.Bool("o", false, "overwrite")
	opt_p = flag.String("p", "", "prefix")
	opt_s = flag.String("s", "", "external script")
	opt_v = flag.Bool("v", false, "version")

	reShorted = regexp.MustCompile(`></(meta|parser|node|dep|acl|advcl|advmod|amod|appos|aux|case|cc|ccomp|clf|compound|conj|cop|csubj|det|discourse|dislocated|expl|fixed|flat|goeswith|iobj|list|mark|nmod|nsubj|nummod|obj|obl|orphan|parataxis|punct|ref|reparandum|root|vocative|xcomp)>`)

	reJunk = regexp.MustCompile(`(?s:<ud:ud.*?</ud:ud>)|(?s:<ud:conllu.*?</ud:conllu>)`)
	chQuit = make(chan bool)
)

func usage() {
	p := filepath.Base(os.Args[0])
	fmt.Printf(`
Usage, examples:

  %s file.(xml|dact)...
  %s -l filelist
  find . -name '*.xml' | %s

  -l filelist : file with list of names of xml and/or dact files

What happens if a file already contains CoNLL-U data?
With option -i the data is ignored.
With option -k the data is kept.
In other cases, the data is ignored only if it was generated with
an older version of the pqudep program.

More options:

  -o : overwrite original file (default: save with .tmp)
  -p prefix : remove prefix from filename in stderr
  -s script : use external script instead of built-in (implies -i, for development)
  -v : print version and exit

`, p, p, p)
}

func main() {

	flag.Usage = usage
	flag.Parse()

	if *opt_v {
		fmt.Printf(VERSIONs+"\n", VERSIONxq, VERSIONxml)
		return
	}

	if flag.NArg() == 0 && *opt_l == "" && util.IsTerminal(os.Stdin) {
		usage()
		return
	}

	if *opt_s != "" {
		b, err := ioutil.ReadFile(*opt_s)
		x(err)
		udep = string(b)
		*opt_i = true
		*opt_k = false
	}

	go func() {
		chSignal := make(chan os.Signal, 1)
		signal.Notify(chSignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		sig := <-chSignal
		close(chQuit)
		fmt.Printf("\r\033[KSignal: %v\n", sig)
	}()

	cs := C.CString(udep)
	if e := C.init(cs); e != 0 {
		x(fmt.Errorf("C.init: [%d] %s", e, C.GoString(C.error_value)))
	}
	C.free(unsafe.Pointer(cs))

	filenames := make([]string, 0)

	if !util.IsTerminal(os.Stdin) {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			filename := strings.TrimSpace(scanner.Text())
			if filename != "" {
				filenames = append(filenames, filename)
			}
		}
		x(scanner.Err())
	}

	if *opt_l != "" {
		fp, err := os.Open(*opt_l)
		x(err)
		scanner := bufio.NewScanner(fp)
		for scanner.Scan() {
			filename := strings.TrimSpace(scanner.Text())
			if filename != "" {
				filenames = append(filenames, filename)
			}
		}
		x(scanner.Err())
	}

	for _, filename := range flag.Args() {
		filenames = append(filenames, filename)
	}

LOOP:
	for i, filename := range filenames {
		doFile(filename, i+1, len(filenames))
		select {
		case <-chQuit:
			break LOOP
		default:
		}
	}
	fmt.Print("\r\033[K")
}

func doFile(filename string, index, length int) {
	if strings.HasSuffix(filename, ".dact") {
		db1, err := dbxml.OpenRead(filename)
		x(err)
		size, err := db1.Size()
		x(err)

		os.Remove(filename + ".tmp") // negeer fout

		// is er een gedeeltelijke versie die we kunnen hergebruiken?
		b, err := ioutil.ReadFile(filename + ".tmp.version")
		if err == nil {
			os.Remove(filename + ".tmp.version")
			if fmt.Sprintf(VERSIONs, VERSIONxq, VERSIONxml) == strings.TrimSpace(string(b)) {
				os.Rename(filename+".tmp.partial", filename+".tmp") // negeer fout
			}
		}
		os.Remove(filename + ".tmp.partial") // negeer fout

		db2, err := dbxml.OpenReadWrite(filename + ".tmp")
		x(err)
		docs, err := db1.All()
		teller := 0
		quit := false
		for docs.Next() {
			teller++
			f := docs.Name()
			fmt.Printf("\r\033[K[%d/%d] %s -> [%d/%d] %s ", index, length, filename, teller, size, f)
			if _, err := db2.Get(f); err == nil {
				continue
			}
			xml := docs.Content()
			result := doXml(xml, filename, f)
			select {
			case <-chQuit:
				quit = true
				docs.Close()
			default:
				x(db2.PutXml(f, result, false))
			}
		}
		db2.Close()
		db1.Close()
		if quit {
			os.Rename(filename+".tmp", filename+".tmp.partial")
			fp, err := os.Create(filename + ".tmp.version")
			x(err)
			fmt.Fprintf(fp, VERSIONs+"\n", VERSIONxq, VERSIONxml)
			fp.Close()
			return
		}
	} else {
		fmt.Printf("\r\033[K[%d/%d] %s ", index, length, filename)

		var b []byte
		gz := false

		fp, err := os.Open(filename)
		x(err)
		rd, err := gzip.NewReader(fp)
		if err == nil {
			gz = true
			b, err = ioutil.ReadAll(rd)
			rd.Close()
			fp.Close()
		} else {
			fp.Close()
			b, err = ioutil.ReadFile(filename)
		}
		x(err)

		result := doXml(string(b), "", filename)
		select {
		case <-chQuit:
			return
		default:
		}

		fp, err = os.Create(filename + ".tmp")
		x(err)
		if gz {
			w := gzip.NewWriter(fp)
			fmt.Fprintln(w, result)
			w.Close()
		} else {
			fmt.Fprintln(fp, result)
		}
		fp.Close()
	}
	if *opt_o {
		os.Rename(filename+".tmp", filename)
	}
}

func doXml(document, archname, filename string) (result string) {

	var alpino Alpino_ds
	var lineno int
	var err error
	lines := make([]string, 0)

	defer func() {
		if alpino.Conllu == nil {
			alpino.Conllu = &ConlluType{}
		}
		alpino.Conllu.Conllu = "\n" + strings.Join(lines, "\n") + "\n"
		if err == nil {
			alpino.Conllu.Status = "OK"
		} else {
			if *opt_p != "" {
				if archname == "" {
					filename = strings.Replace(filename, *opt_p, "", 1)
				} else {
					archname = strings.Replace(archname, *opt_p, "", 1)
				}
			}
			if archname == "" {
				fmt.Fprintln(os.Stderr, ">>>", filename)
			} else {
				fmt.Fprintln(os.Stderr, ">>>", archname, "/", filename)
			}
			if alpino.Sentence == nil {
				alpino.Sentence = &SentType{}
			}
			if id := alpino.Sentence.SentId; id != "" {
				fmt.Fprintln(os.Stderr, "# sent_id =", id)
			}
			if t := alpino.Sentence.Sent; t != "" {
				fmt.Fprintln(os.Stderr, "# text =", t)
			}
			if lineno == 0 {
				fmt.Fprintln(os.Stderr, "^^^", err)
			}
			for i, line := range lines {
				fmt.Fprintln(os.Stderr, line)
				if i+1 == lineno {
					fmt.Fprintln(os.Stderr, "^^^", err)
				}
			}
			fmt.Fprintln(os.Stderr)
			alpino.Conllu.Status = "error"
			alpino.Conllu.Error = fmt.Sprintf("Line %d: %v", lineno, err)
			if alpino.Node != nil {
				clean(alpino.Node)
			}
		}
		minify(alpino)
		result = format(alpino)
	}()

	document = reJunk.ReplaceAllString(document, "") // erg oude versie
	err = xml.Unmarshal([]byte(document), &alpino)
	if err != nil {
		lineno = 0
		return
	}

	reset(alpino.Node)
	if oldVersion(alpino.Version) {
		alpino.Version = fmt.Sprintf("%d.%d", ALPINO_DS_MAJOR, ALPINO_DS_MINOR)
	}
	if alpino.Conllu == nil {
		alpino.Conllu = &ConlluType{}
	} else {
		if *opt_i {
			alpino.Conllu = &ConlluType{}
		} else if !*opt_k {
			if strings.HasPrefix(alpino.Conllu.Auto, "PQU") && (isOld(alpino.Conllu.Auto) || alpino.Conllu.Status != "OK") {
				alpino.Conllu = &ConlluType{}
			}
		}
	}
	alpino.Conllu.Status = ""
	alpino.Conllu.Error = ""

	for _, line := range strings.Split(alpino.Conllu.Conllu, "\n") {
		line = strings.TrimSpace(line)
		if line != "" && line[0] != '#' {
			lines = append(lines, line)
		}
	}

	if len(lines) == 0 || strings.HasPrefix(alpino.Conllu.Auto, "PQU") {
		if *opt_s != "" {
			alpino.Conllu.Auto = "script:" + fmt.Sprintf(*opt_s)
		} else {
			alpino.Conllu.Auto = fmt.Sprintf(VERSIONs, VERSIONxq, VERSIONxml)
		}
	}

	if len(lines) == 0 {
		cs := C.CString(document)
		e := C.parse(cs)
		C.free(unsafe.Pointer(cs))
		if e != 0 {
			err = fmt.Errorf("C.parse: [%d] %s", e, C.GoString(C.error_value))
			lineno = 0
			return
		}
		for {
			if e := C.next(); e != 0 {
				err = fmt.Errorf("C.next: [%d] %s", e, C.GoString(C.error_value))
				lineno = len(lines)
				return
			}
			if C.done != 0 {
				break
			}
			lines = append(lines, C.GoString(C.value))
		}
	}

	valid := make(map[string]bool)
	valid["0"] = true // root
	copies := make(map[string]string)
	for i, line := range lines {
		a := strings.Split(line, "\t")
		if len(a) != 10 {
			err = fmt.Errorf("Wrong number of fields")
			lineno = i + 1
			return
		}
		if strings.Contains(a[3], "ERROR") {
			err = fmt.Errorf("Invalid UPOS value %v", a[3])
			lineno = i + 1
			return
		}
		if strings.Contains(a[5], "ERROR") {
			err = fmt.Errorf("Invalid FEAT value %v", a[5])
			lineno = i + 1
			return
		}
		if strings.Contains(a[6], "ERROR") {
			err = fmt.Errorf("Invalid HEAD value %v", a[6])
			lineno = i + 1
			return
		}
		if strings.Contains(a[7], "ERROR") {
			err = fmt.Errorf("Invalid DEPREL value %v", a[7])
			lineno = i + 1
			return
		}
		if a[7] == "root" && a[6] != "0" || a[7] != "root" && a[6] == "0" {
			err = fmt.Errorf("Invalid HEAD/DEPREL combination %v/%v", a[6], a[7])
			lineno = i + 1
			return
		}
		if strings.Contains(a[8], "ERROR") {
			err = fmt.Errorf("Invalid DEPS value %v", a[8])
			lineno = i + 1
			return
		}
		valid[a[0]] = true
		items := getItems(a[9])
		if n, ok := items["CopiedFrom"]; ok {
			copies[n] = a[0]
		}
	}

	udNodeList := make([]*UdNodeType, 0)
	eudNodeList := make([]*UdNodeType, 0)

	for i, line := range lines {
		lineno = i + 1
		a := strings.Split(line, "\t")
		for i := range a {
			a[i] = strings.TrimSpace(a[i])
		}
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
			err = fmt.Errorf("Node '%s' not found", es)
			return
		}

		if a[8] != "_" {
			for _, deps := range strings.Split(a[8], "|") {
				dep := strings.SplitN(deps, ":", 2)
				if len(dep) != 2 {
					err = fmt.Errorf("Not a valid dependency: %s", dep)
					return
				}
				if !valid[dep[0]] {
					err = fmt.Errorf("Not a valid head: %s", dep[0])
					return
				}
				var aux string
				dd := strings.SplitN(dep[1], ":", 2)
				if len(dd) > 1 {
					aux = dd[1]
				}
				if dd[0] == "root" && dep[0] != "0" || dd[0] != "root" && dep[0] == "0" {
					err = fmt.Errorf("Invalid combination in DEPS %v", deps)
					return
				}
				node.Ud.Dep = append(node.Ud.Dep, DepType{
					Id:         a[0],
					Head:       dep[0],
					Deprel:     dep[1],
					DeprelMain: dd[0],
					DeprelAux:  aux,
					Elided:     strings.Contains(a[0], "."),
				})
			}
		}

		if isCopy {
			continue
		}

		if !valid[a[6]] {
			err = fmt.Errorf("Not a valid head: %s", a[6])
			return
		}

		node.Ud.Id = a[0]
		node.Ud.Form = a[1]
		node.Ud.Lemma = a[2]
		node.Ud.Upos = noe(a[3])
		node.Ud.Head = noe(a[6])
		node.Ud.Deprel = noe(a[7])

		dd := strings.SplitN(node.Ud.Deprel, ":", 2)
		node.Ud.DeprelMain = dd[0]
		if len(dd) > 1 {
			node.Ud.DeprelAux = dd[1]
		}

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

		ud := UdNodeType{
			recursion: make([]string, 0),

			XMLName:   xml.Name{Local: node.Ud.DeprelMain},
			Id:        node.Ud.Id,
			Form:      node.Ud.Form,
			Lemma:     node.Ud.Lemma,
			Upos:      node.Ud.Upos,
			Head:      node.Ud.Head,
			Deprel:    node.Ud.Deprel,
			DeprelAux: node.Ud.DeprelAux,

			FeatsType: node.Ud.FeatsType,

			Buiging:  node.Buiging,
			Conjtype: node.Conjtype,
			Dial:     node.Dial,
			Genus:    node.Genus,
			Getal:    node.Getal,
			GetalN:   node.GetalN,
			Graad:    node.Graad,
			Lwtype:   node.Lwtype,
			Naamval:  node.Naamval,
			Npagr:    node.Npagr,
			Ntype:    node.Ntype,
			Numtype:  node.Numtype,
			Pdtype:   node.Pdtype,
			Persoon:  node.Persoon,
			Positie:  node.Positie,
			Pt:       node.Pt,
			Pvagr:    node.Pvagr,
			Pvtijd:   node.Pvtijd,
			Spectype: node.Spectype,
			Status:   node.Status,
			Vwtype:   node.Vwtype,
			Vztype:   node.Vztype,
			Wvorm:    node.Wvorm,
		}
		udNodeList = append(udNodeList, &ud)

		for _, dep := range node.Ud.Dep {

			ud := UdNodeType{
				recursion: make([]string, 0),

				XMLName:   xml.Name{Local: dep.DeprelMain},
				Id:        dep.Id,
				Form:      node.Ud.Form,
				Lemma:     node.Ud.Lemma,
				Upos:      node.Ud.Upos,
				Head:      dep.Head,
				Deprel:    dep.Deprel,
				DeprelAux: dep.DeprelAux,

				Enhanced: true,

				FeatsType: node.Ud.FeatsType,

				Buiging:  node.Buiging,
				Conjtype: node.Conjtype,
				Dial:     node.Dial,
				Genus:    node.Genus,
				Getal:    node.Getal,
				GetalN:   node.GetalN,
				Graad:    node.Graad,
				Lwtype:   node.Lwtype,
				Naamval:  node.Naamval,
				Npagr:    node.Npagr,
				Ntype:    node.Ntype,
				Numtype:  node.Numtype,
				Pdtype:   node.Pdtype,
				Persoon:  node.Persoon,
				Positie:  node.Positie,
				Pt:       node.Pt,
				Pvagr:    node.Pvagr,
				Pvtijd:   node.Pvtijd,
				Spectype: node.Spectype,
				Status:   node.Status,
				Vwtype:   node.Vwtype,
				Vztype:   node.Vztype,
				Wvorm:    node.Wvorm,
			}
			eudNodeList = append(eudNodeList, &ud)
		}

	}

	alpino.Root = make([]*UdNodeType, 0)

	for _, n := range udNodeList {
		if n.Head == "0" {
			alpino.Root = append(alpino.Root, n)
		}
	}

	for _, n := range eudNodeList {
		if n.Head == "0" {
			alpino.Root = append(alpino.Root, n)
		}
	}

	for i, root := range alpino.Root {
		var items []*UdNodeType
		if i == 0 {
			items = udNodeList
		} else {
			items = eudNodeList
		}
		root.UdNodes = make([]*UdNodeType, 0)
		expand(root, items)
	}

	return
}

func expand(udnode *UdNodeType, items []*UdNodeType) {
	for _, item := range items {
		if item.Head == udnode.Id {
			it := new(UdNodeType)
			*it = *item
			it.UdNodes = make([]*UdNodeType, 0)
			it.recursion = append([]string{udnode.Id}, udnode.recursion...)
			udnode.UdNodes = append(udnode.UdNodes, it)
		}
	}
	for _, un := range udnode.UdNodes {
		if recursionLimit(un.recursion) {
			un.RecursionLimit = "TOO DEEP"
			un.Head = ""
		} else {
			expand(un, items)
		}
	}
	udnode.Head = ""
}

func recursionLimit(s []string) bool {
	if len(s) < 2 {
		return false
	}
	found := 0
	for i := 1; i < len(s)-1; i++ {
		if s[i] == s[0] && s[i+1] == s[1] {
			found++
		}
	}
	return found > 1
}

func noe(s string) string {
	if s == "_" {
		return ""
	}
	return s
}

func format(alpino Alpino_ds) string {
	b, err := xml.MarshalIndent(&alpino, "", "  ")
	x(err)
	s := "<?xml version=\"1.0\"?>\n" + string(b)

	// shorten
	s = reShorted.ReplaceAllString(s, "/>")

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
	if alpino.Parser != nil && alpino.Parser.Build == "" && alpino.Parser.Date == "" && alpino.Parser.Cats == "" && alpino.Parser.Skips == "" {
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
			if len(node.Ud.Dep) == 0 {
				node.Ud.Dep = nil
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
	node.Ud = &UdType{Dep: make([]DepType, 0)}
	for _, n := range node.NodeList {
		reset(n)
	}
}

func oldVersion(ver string) bool {
	vv := strings.Split(ver, ".")
	if len(vv) < 2 {
		return true
	}
	v1, err := strconv.Atoi(vv[0])
	if err != nil {
		return true
	}
	v2, err := strconv.Atoi(vv[1])
	if err != nil {
		return true
	}
	if v1 < ALPINO_DS_MAJOR || (v1 == ALPINO_DS_MAJOR && v2 < ALPINO_DS_MINOR) {
		return true
	}
	return false
}

func isOld(ver string) bool {
	vv := strings.Split(ver[3:], ".")
	i, err := strconv.Atoi(vv[0])
	return err != nil || i < VERSIONxq
}
