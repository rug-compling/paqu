package main

import (
	"github.com/pebbe/util"
	"github.com/rug-compling/alud/v2"

	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"
)

type Response struct {
	Code            int
	Status          string
	Message         string
	Id              string
	Interval        int
	Number_of_lines int
	Timeout         int
	Max_tokens      int
	Finished        bool
	Batch           []Line
}

type Line struct {
	Line_status string
	Line_number int
	Label       string
	Sentence    string
	Alpino_ds   string
	Log         string
}

type AlpinoInfo struct {
	ParserBuild string `json:"parser_build"`
}

var (
	opt_d = flag.String("d", "xml", "directory voor uitvoer")
	opt_e = flag.String("e", "half", "escape level: none / half / full")
	opt_l = flag.Bool("l", false, "true: één zin per regel; false: doorlopende tekst")
	opt_L = flag.String("L", "doc", "prefix voor labels")
	opt_n = flag.Int("n", 0, "maximum aantal tokens per regel")
	opt_p = flag.String("p", "", "alternatieve parser")
	opt_q = flag.Bool("q", false, "true: quiet")
	opt_s = flag.String("s", "", "URL van Alpino-server")
	opt_t = flag.Int("t", 900, "time-out in seconden per regel")
	opt_T = flag.Bool("T", false, "true: zinnen zijn getokeniseerd")
	opt_u = flag.String("u", "", "output file for UD errors (impliceert -U)")
	opt_U = flag.Bool("U", false, "true: derive Universal Dependencies")

	x             = util.CheckErr
	reParser      = regexp.MustCompile(`<parser.*?>`)
	reParserBuild = regexp.MustCompile(`<parser[^>]*?build=.*?>`)
	alpino_build  string
	filename      string
	lastdir       string
)

func usage() {
	fmt.Fprintf(os.Stderr, `
Syntax: %s [opties] datafile

Optie:

  -s string : Alpino-server, zie: https://github.com/rug-compling/alpino-api
              Als deze ontbreekt wordt een lokale versie van Alpino gebruikt

Zonder gebruik van Alpino-server:

  De tekst moet bestaan uit één zin per regel, getokeniseerd, met of zonder labels.

Met gebruik van Alpino-server:

  De tekst kan verschillende vormen hebben.

Overige opties:

  -d string : Directory waar uitvoer wordt geplaatst (default: xml)
  -e string : Escape level: none / half / full (default: half)
  -n int    : Maximum aantal tokens per regel (default: 0 = geen limiet)
  -p string : Alternatieve parser, zoals qa (default: geen)
  -t int    : Time-out per regel (default: 900)
  -U        : Derive Universal Dependencies (UD) (default: nee)
  -u string : Output file for UD errors (impliceert -U) (default: geen)

Opties alleen van toepassing bij gebruik van Alpino-server:

  -l        : Eén zin per regel (default: doorlopende tekst)
  -L string : Prefix voor labels (default: doc)
  -q        : Stil
  -T        : Zinnen zijn getokeniseerd (default: niet getokeniseerd)

`, os.Args[0])
}

func main() {

	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		usage()
		return
	}
	filename = flag.Arg(0)

	// PARSEN

	if *opt_s == "" {
		doLocal()
	} else {
		info, err := doServerInfo()
		if err == nil {
			doServer(info)
			return
		}
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, "Fallback to local Alpino")
		if *opt_l && *opt_T {
			doLocal()
		} else {
			fmt.Fprintln(os.Stderr, "Lokale versie van Alpino vereist getokeniseerde tekst, één zin per regel")
		}
	}
}

func doLocal() {

	b, err := ioutil.ReadFile(os.Getenv("ALPINO_HOME") + "/version")
	x(err)
	alpino_build = strings.TrimSpace(string(b))

	var fpud *os.File
	if *opt_u != "" {
		fpud, err = os.Create(*opt_u)
		x(err)
	}

	var fpin, fpout *os.File
	var errval error
	tmpfile := filename + ".part"
	var maxtok, parser string
	if *opt_n > 0 {
		maxtok = fmt.Sprint("max_sentence_length=", *opt_n)
	}
	if *opt_p != "" {
		parser = "application_type=" + *opt_p
	}
	defer func() {
		if fpin != nil {
			fpin.Close()
		}
		if fpout != nil {
			fpout.Close()
		}
		os.Remove(tmpfile)
		if errval != io.EOF {
			x(errval)
		}
	}()
	fpin, errval = os.Open(filename)
	if errval != nil {
		return
	}
	rd := util.NewReaderSize(fpin, 5000)
	lineno := 0
	for {
		line, err := rd.ReadLineString()
		if err != nil && err != io.EOF {
			errval = err
			return
		}
		if err == nil && lineno == 0 {
			fpout, errval = os.Create(tmpfile)
			if errval != nil {
				return
			}
		}
		if err == nil {
			if strings.HasPrefix(line, "%") {
				continue
			}
			lineno++
			var label string
			a := strings.SplitN(line, "|", 2)
			if len(a) == 2 {
				a[0] = strings.TrimSpace(a[0])
				a[1] = strings.TrimSpace(a[1])
				if a[0] == "" {
					a[0] = fmt.Sprint(lineno)
				}
				label = a[0]
				line = a[0] + "|" + escape(a[1])
			} else {
				label = fmt.Sprint(lineno)
				line = label + "|" + escape(line)
			}
			if *opt_n > 0 {
				if n := len(strings.Fields(line)); n > *opt_n {
					fmt.Fprintf(os.Stderr, `**** parsing %s (line number %d)
line too long: %d tokens
Q#%s|skipped|??|????
**** parsed %s (line number %d)
`,
						label, lineno,
						n,
						line,
						label, lineno)
					continue
				}
			}
			fmt.Fprintln(fpout, line)
			dirname := filepath.Dir(filepath.Join(*opt_d, label))
			if dirname != lastdir {
				lastdir = dirname
				os.MkdirAll(dirname, 0777)
			}
		}
		if (err == io.EOF && lineno%10000 != 0) || lineno%10000 == 0 {
			fpout.Close()
			fpout = nil
			cmd := exec.Command(
				"/bin/bash",
				"-c",
				fmt.Sprintf(
					"$ALPINO_HOME/bin/Alpino -veryfast -flag treebank %s debug=1 end_hook=xml user_max=%d %s %s -parse < %s",
					*opt_d, *opt_t*1000, maxtok, parser, tmpfile))
			cmd.Env = []string{
				"ALPINO_HOME=" + os.Getenv("ALPINO_HOME"),
				"PATH=" + os.Getenv("ALPINO_HOME") + "/bin:" + os.Getenv("PATH"),
				"LANG=en_US.utf8",
				"LANGUAGE=en_US.utf8",
				"LC_ALL=en_US.utf8",
			}
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			errval = cmd.Run()
			if errval != nil {
				return
			}
		}
		if err == io.EOF {
			break
		}
	}
	// UD, version en date invoegen
	filenames := make([]string, 0)
	x(filepath.Walk(*opt_d, func(path string, info os.FileInfo, err error) error {
		x(err)
		if !info.IsDir() && strings.HasSuffix(path, ".xml") {
			filenames = append(filenames, path)
		}
		return nil
	}))
	for _, filename := range filenames {
		b, err := ioutil.ReadFile(filename)
		x(err)
		if reParserBuild.Match(b) {
			continue
		}
		xml := setBuild(string(b))
		if *opt_U || *opt_u != "" {
			s, err := alud.UdAlpino([]byte(xml), filename)
			if s != "" {
				xml = s
			}
			if err != nil && *opt_u != "" {
				fmt.Fprintln(fpud, ">>>", filename)
				fmt.Fprintln(fpud, "^^^", err)
				fmt.Fprintln(fpud)
			}
		}
		fp, err := os.Create(filename)
		x(err)
		fp.WriteString(xml)
		fp.Close()
	}
	if *opt_u != "" {
		x(fpud.Close())
	}
}

func doServerInfo() (*AlpinoInfo, error) {

	buf1 := bytes.NewBufferString(`{"request":"info"}`)
	resp, err := http.Post(*opt_s, "application/json", buf1)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	var info AlpinoInfo
	err = json.Unmarshal(data, &info)
	return &info, err
}

func doServer(info *AlpinoInfo) {
	alpino_build = info.ParserBuild

	var fpud *os.File
	var err error
	if *opt_u != "" {
		fpud, err = os.Create(*opt_u)
		x(err)
	}

	var buf bytes.Buffer
	var dataType string
	if *opt_l {
		if *opt_T {
			dataType = "lines tokens " + *opt_e
		} else {
			dataType = "lines"
		}
	} else {
		if *opt_L == "" {
			dataType = "text doc"
		} else {
			dataType = "text " + *opt_L
		}
	}
	fmt.Fprintf(
		&buf,
		`{"request":"parse", "data_type":%q, "timeout":%d, "parser":%q, "max_tokens":%d. "ud":false}`,
		dataType,
		*opt_t,
		*opt_p,
		*opt_n)
	fp, err := os.Open(filename)
	x(err)
	_, err = io.Copy(&buf, fp)
	fp.Close()
	x(err)
	resp, err := http.Post(*opt_s, "application/json", &buf)
	util.CheckErr(err)
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	util.CheckErr(err)
	var response Response
	err = json.Unmarshal(data, &response)
	util.CheckErr(err)
	if response.Code > 299 {
		x(fmt.Errorf("%d %s -- %s", response.Code, response.Status, response.Message))
	}
	maxinterval := (response.Interval * 100) / 80
	totallines := response.Number_of_lines
	id := response.Id
	if !*opt_q {
		if response.Timeout > 0 {
			fmt.Printf("timeout: %ds\n", response.Timeout)
		}
		if response.Max_tokens > 0 {
			fmt.Printf("max tokens: %d\n", response.Max_tokens)
		}
		fmt.Println(totallines)
	}

	go func() {
		chSignal := make(chan os.Signal, 1)
		signal.Notify(chSignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		sig := <-chSignal
		fmt.Fprintf(os.Stderr, "Signal: %v\n", sig)

		var buf bytes.Buffer
		fmt.Fprintf(&buf, `{"request":"cancel", "id":%q}`, id)
		resp, err := http.Post(*opt_s, "application/json", &buf)
		util.CheckErr(err)
		_, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		util.CheckErr(err)

		os.Exit(0)
	}()

	seen := 0
	interval := 2
	incr := true
	moment := time.Now()
	for {
		if interval > maxinterval {
			interval = maxinterval
		}
		if interval > 120 {
			interval = 120
		}

		if sleep := time.Duration(interval)*time.Second - time.Now().Sub(moment); sleep > 0 {
			time.Sleep(sleep)
		}

		var buf bytes.Buffer
		fmt.Fprintf(&buf, `{"request":"output", "id":%q}`, id)
		resp, err := http.Post(*opt_s, "application/json", &buf)
		util.CheckErr(err)

		moment = time.Now()

		data, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		util.CheckErr(err)

		var response Response
		err = json.Unmarshal(data, &response)
		util.CheckErr(err)
		if response.Code > 299 {
			x(fmt.Errorf("%d %s -- %s", response.Code, response.Status, response.Message))
		}
		seen += len(response.Batch)
		if !*opt_q {
			if totallines > 0 {
				fmt.Println(totallines - seen)
			} else {
				fmt.Println(seen)
			}
		}
		for _, line := range response.Batch {
			if line.Line_status == "ok" {
				if line.Label == "" {
					line.Label = fmt.Sprint(line.Line_number)
				}
				filename := filepath.Join(*opt_d, line.Label+".xml")
				if *opt_U || *opt_u != "" {
					s, err := alud.UdAlpino([]byte(line.Alpino_ds), filename)
					if s != "" {
						line.Alpino_ds = s
					}
					if err != nil && *opt_u != "" {
						fmt.Fprintln(fpud, ">>>", filename)
						fmt.Fprintln(fpud, "^^^", err)
						fmt.Fprintln(fpud)
					}
				}
				dirname := filepath.Dir(filename)
				if dirname != lastdir {
					lastdir = dirname
					os.MkdirAll(dirname, 0777)
				}
				fp, err := os.Create(filename)
				x(err)
				fmt.Fprintln(fp, setBuild(line.Alpino_ds))
				fp.Close()
			} else {
				fmt.Fprintf(os.Stderr, `**** parsing %s (line number %d)
%s
Q#%s|%s|%s|??|????
**** parsed %s (line number %d)
`,
					line.Label, line.Line_number,
					line.Log,
					line.Label, line.Sentence, line.Line_status,
					line.Label, line.Line_number)
			}
		}

		if response.Finished {
			break
		}
		if incr && totallines > 0 && len(response.Batch) > totallines-seen {
			incr = false
			interval *= totallines - seen
			interval /= len(response.Batch)
			if interval < 2 {
				interval = 2
			}
		}
		if incr {
			interval = (3 * interval) / 2
		}
	}
	if *opt_u != "" {
		x(fpud.Close())
	}
}

func escape(s string) string {
	if *opt_e == "none" {
		return s
	}
	words := strings.Fields(s)
	for i, word := range words {
		switch word {
		case `[`:
			words[i] = `\[`
		case `]`:
			words[i] = `\]`
		case `\[`:
			if *opt_e == "full" {
				words[i] = `\\[`
			}
		case `\]`:
			if *opt_e == "full" {
				words[i] = `\\]`
			}
		}
	}
	return strings.Join(words, " ")
}

func setBuild(xml string) string {
	return reParser.ReplaceAllStringFunc(xml, func(s string) string {
		if !strings.Contains(s, "date=") {
			s = s[:7] + fmt.Sprintf(" date=%q", time.Now().Format(time.RFC3339)) + s[7:]
		}
		if !strings.Contains(s, "build=") {
			s = s[:7] + fmt.Sprintf(" build=%q", alpino_build) + s[7:]
		}
		return s
	})
}
