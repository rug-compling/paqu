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

	reUnQ1 = regexp.MustCompile(`( (?:,,|'')) (.*?) ('' )`)
	reUnQ2 = regexp.MustCompile(`( ") (.*?) (" )`)
	reUnQ3 = regexp.MustCompile("( [`']) (.*?) (' )")
	reUnQ4 = regexp.MustCompile(`( [‘’]) (.*?) (’ )`)
	reUnQ5 = regexp.MustCompile(`( [“„”]) (.*?) (” )`)
	reUnP1 = regexp.MustCompile(`([\[({]) `)
	reUnP2 = regexp.MustCompile(` ([\])}:;.,!?])`)
)

func usage() {
	p := filepath.Base(os.Args[0])
	fmt.Printf(`
Usage, examples:

  %s file.(xml|xml.gz|dact)...
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
	err := xml.Unmarshal([]byte(document), &alpino)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR in (%s) %s: %v\n", archname, filename, err)
		return
	}
	if alpino.Conllu == nil {
		return
	}

	lines := strings.Split(strings.TrimSpace(alpino.Conllu.Conllu), "\n")

	// untokenize
	alines := make([][]string, len(lines))
	for i, line := range lines {
		alines[i] = strings.Fields(line)
	}
	words1 := make([]string, 0, len(lines))
	for _, line := range alines {
		if len(line) != 10 {
			fmt.Fprintf(os.Stderr, "ERROR in (%s) %s: %d columns\n", archname, filename, len(line))
			return
		}
		if !strings.Contains(line[0], ".") {
			words1 = append(words1, line[1])
		}
	}
	sent1 := strings.Join(words1, " ")
	sent2 := reUnQ1.ReplaceAllString(" "+sent1+" ", `$1$2$3`)
	sent2 = reUnQ2.ReplaceAllString(sent2, `$1$2$3`)
	sent2 = reUnQ3.ReplaceAllString(sent2, `$1$2$3`)
	sent2 = reUnQ4.ReplaceAllString(sent2, `$1$2$3`)
	sent2 = reUnQ5.ReplaceAllString(sent2, `$1$2$3`)
	sent2 = reUnP1.ReplaceAllString(sent2, `$1`)
	sent2 = reUnP2.ReplaceAllString(sent2, `$1`)
	sent2 = strings.TrimSpace(sent2)
	noSpace := make(map[string]bool)
	words2 := strings.Fields(sent2)
	j := 0
	for i, word1 := range words1 {
		if word1 != words2[j] {
			noSpace[fmt.Sprint(i+1)] = true
			words2[j] = words2[j][len(word1):]
		} else {
			j++
		}
	}
	for i, line := range alines {
		if noSpace[line[0]] {
			if alines[i][9] == "_" {
				alines[i][9] = "SpaceAfter=No"
			} else {
				alines[i][9] += "|SpaceAfter=No"
			}
			lines[i] = strings.Join(alines[i], "\t")
		}
	}

	if archname != "" {
		fmt.Println("# archive =", archname)
	}
	fmt.Println("# source =", filename)
	if id := alpino.Sentence.SentId; id != "" {
		fmt.Println("# sent_id =", id)
	}
	fmt.Println("# text =", sent2)
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
		coms := strings.Split(com, "\n")
		if len(coms) == 1 {
			fmt.Printf("# comment_%d = %s\n", i+1, com)
		} else {
			for j, cc := range coms {
				fmt.Printf("# comment_%d.%02d = %s\n", i+1, j+1, cc)
			}
		}
	}

	fmt.Println(strings.Join(lines, "\n"))
	fmt.Println()
}
