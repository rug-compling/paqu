package main

import (
	"github.com/pebbe/dbxml"
	"github.com/pebbe/util"
	"github.com/rug-compling/alud/v2"

	"bufio"
	"compress/gzip"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

type conlluType struct {
	Conllu string `xml:",cdata"`
	Status string `xml:"status,attr"`
	Auto   string `xml:"auto,attr"`
}

var (
	x     = util.CheckErr
	opt_i = flag.Bool("i", false, "ignore existing cdata from <conllu>")
	opt_k = flag.Bool("k", false, "keep existing cdata from <conllu>")
	opt_l = flag.String("l", "", "filelist")
	opt_o = flag.Bool("o", false, "overwrite")
	opt_p = flag.String("p", "", "prefix")
	opt_v = flag.Bool("v", false, "version")

	versionID = alud.VersionID()

	chQuit = make(chan bool)
)

func usage() {
	p := filepath.Base(os.Args[0])
	fmt.Printf(`
Usage, examples:

  %s file.(xml|xml.gz|dact)...
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
  -v : print version and exit

`, p, p, p)
}

func main() {

	flag.Usage = usage
	flag.Parse()

	if *opt_v {
		fmt.Println(versionID)
		return
	}

	if flag.NArg() == 0 && *opt_l == "" && util.IsTerminal(os.Stdin) {
		usage()
		return
	}

	go func() {
		chSignal := make(chan os.Signal, 1)
		signal.Notify(chSignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		sig := <-chSignal
		close(chQuit)
		fmt.Printf("\r\033[KSignal: %v\n", sig)
	}()

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
			if strings.TrimSpace(string(b)) == versionID {
				os.Rename(filename+".tmp.partial", filename+".tmp") // negeer fout
			}
		}
		os.Remove(filename + ".tmp.partial") // negeer fout

		db2, err := dbxml.OpenReadWrite(filename + ".tmp")
		x(err)
		docs, err := db1.All()
		x(err)
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
			fmt.Fprintln(fp, versionID)
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

	reuse := false
	auto := ""
	for {
		if *opt_i {
			break
		}
		i := strings.Index(document, "<conllu")
		if i < 0 {
			break
		}
		j := strings.Index(document, "</conllu>")
		if j < 0 {
			break
		}
		s := document[i : j+9]
		var conllu conlluType
		err := xml.Unmarshal([]byte(s), &conllu)
		if err != nil || strings.TrimSpace(conllu.Conllu) == "" || conllu.Status != "OK" {
			break
		}
		auto = conllu.Auto
		if *opt_k {
			reuse = true
			break
		}
		if strings.HasPrefix(conllu.Auto, "ALUD") || strings.HasPrefix(conllu.Auto, "PQU") {
			break
		}
		reuse = true
		break
	}

	var err error
	if reuse {
		result, err = alud.Alpino([]byte(document), "", auto)
	} else {
		result, err = alud.UdAlpino([]byte(document), filename)
	}
	if err != nil {
		if result == "" {
			x(err)
		}
		if archname == "" {
			fmt.Fprintln(os.Stderr, ">>>", filename)
		} else {
			fmt.Fprintln(os.Stderr, ">>>", archname, "/", filename)
		}
		fmt.Fprintln(os.Stderr, "^^^", err)
		fmt.Fprintln(os.Stderr)
	}
	return result
}
