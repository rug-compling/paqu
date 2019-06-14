package main

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
	"path/filepath"
	"regexp"
	"strings"
)

type Alpino_ds struct {
	XMLName  xml.Name      `xml:"alpino_ds"`
	Metadata *MetadataType `xml:"metadata"`
	Parser   *ParserType   `xml:"parser"`
	Sentence *SentType     `xml:"sentence"`
	Comments []string      `xml:"comments>comment"`
	Conllu   *ConlluType   `xml:"conllu"`
}

type MetadataType struct {
	Meta []MetaType `xml:"meta"`
}

type SentType struct {
	Sent   string `xml:",chardata"`
	SentId string `xml:"sentid,attr"`
}

type MetaType struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type ParserType struct {
	Build string `xml:"build,attr"`
	Date  string `xml:"date,attr"`
	Cats  string `xml:"cats,attr"`
	Skips string `xml:"skips,attr"`
}

type ConlluType struct {
	Conllu string `xml:",cdata"`
	Status string `xml:"status,attr"`
	Error  string `xml:"error,attr"`
	Auto   string `xml:"auto,attr"`
}

var (
	x = util.CheckErr

	opt_l = flag.String("l", "", "filelist")

	reUnQ1 = regexp.MustCompile(`( ,,) (.*?) ('' )`)
	reUnQ2 = regexp.MustCompile(`( ") (.*?) (" )`)
	reUnQ3 = regexp.MustCompile("( [`']) (.*?) (' )")
	reUnQ4 = regexp.MustCompile(`( ‘) (.*?) (’ )`)
	reUnQ5 = regexp.MustCompile(`( [“„”]) (.*?) (” )`)
	reUnP1 = regexp.MustCompile(`([\[({]) `)
	reUnP2 = regexp.MustCompile(` ([\])}:;.,!?])`)
)

func usage() {
	p := filepath.Base(os.Args[0])
	fmt.Printf(`
Usage, examples:

  %s file.(xml|dact)...
  %s -l filelist
  find . -name '*.xml' | %s

  -l filelist : file with list of names of xml and/or dact files

`, p, p, p)
}

func main() {

	flag.Usage = usage
	flag.Parse()

	if flag.NArg() == 0 && *opt_l == "" && util.IsTerminal(os.Stdin) {
		usage()
		return
	}

	filenames := make([]string, 0)

	if flag.NArg() == 0 && *opt_l == "" {
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

	for _, filename := range filenames {
		doFile(filename)
	}
}

func doFile(filename string) {
	if strings.HasSuffix(filename, ".dact") {
		db, err := dbxml.OpenRead(filename)
		x(err)
		docs, err := db.All()
		x(err)
		for docs.Next() {
			f := docs.Name()
			xml := docs.Content()
			doXml(xml, filename, f)
		}
		db.Close()
	} else {

		var b []byte

		fp, err := os.Open(filename)
		x(err)
		rd, err := gzip.NewReader(fp)
		if err == nil {
			b, err = ioutil.ReadAll(rd)
			rd.Close()
			fp.Close()
		} else {
			fp.Close()
			b, err = ioutil.ReadFile(filename)
		}
		x(err)

		doXml(string(b), "", filename)
	}
}

func doXml(document, archname, filename string) {
	var alpino Alpino_ds
	x(xml.Unmarshal([]byte(document), &alpino))
	if alpino.Conllu == nil {
		return
	}
	if archname == "" {
		fmt.Println("# source =", filename)
	} else {
		fmt.Println("# source =", archname, "::", filename)
	}
	if id := alpino.Sentence.SentId; id != "" {
		fmt.Println("# sent_id =", id)
	}
	if t := alpino.Sentence.Sent; t != "" {
		t = reUnQ1.ReplaceAllString(" "+t+" ", `$1$2$3`)
		t = reUnQ2.ReplaceAllString(t, `$1$2$3`)
		t = reUnQ3.ReplaceAllString(t, `$1$2$3`)
		t = reUnQ4.ReplaceAllString(t, `$1$2$3`)
		t = reUnQ5.ReplaceAllString(t, `$1$2$3`)
		t = reUnP1.ReplaceAllString(t, `$1`)
		t = reUnP2.ReplaceAllString(t, `$1`)
		fmt.Println("# text =", strings.TrimSpace(t))
	}
	if s := alpino.Conllu.Status; s != "" {
		fmt.Println("# conllu_status =", s)
	}
	if e := alpino.Conllu.Error; e != "" {
		fmt.Println("# conllu_error =", e)
	}
	if a := alpino.Conllu.Auto; a != "" {
		fmt.Println("# conllu_auto =", a)
	}
	if p := alpino.Parser; p != nil {
		if pp := p.Build; pp != "" {
			fmt.Println("# parser_build =", pp)
		}
		if pp := p.Date; pp != "" {
			fmt.Println("# parser_date =", pp)
		}
		if pp := p.Cats; pp != "" {
			fmt.Println("# parser_cats =", pp)
		}
		if pp := p.Skips; pp != "" {
			fmt.Println("# parser_skips =", pp)
		}
	}
	if m := alpino.Metadata; m != nil {
		for _, mm := range m.Meta {
			fmt.Printf("# meta_%s = %s\n", mm.Name, mm.Value)
		}
	}
	for i, com := range alpino.Comments {
		fmt.Printf("# comment_%d = %s\n", i+1, com)
	}

	fmt.Println(strings.TrimSpace(alpino.Conllu.Conllu))
	fmt.Println()
}
