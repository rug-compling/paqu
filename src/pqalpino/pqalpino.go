package main

import (
	"github.com/pebbe/util"
	"github.com/rug-compling/alpinods"
	"github.com/rug-compling/alud/v2"

	"bytes"
	"encoding/json"
	"encoding/xml"
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
	opt_X = flag.Bool("X", false, "true: derive extra attributes, like is_np, is_vorfeld...")

	x            = util.CheckErr
	reParser     = regexp.MustCompile(`<parser.*?>`)
	alpino_build string
	filename     string
	lastdir      string

	idxnodes map[int]*alpinods.Node
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
  -X        : Maak extra attributen: is_np, is_vorfeld, is_nachfeld (default: nee)

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
line too long: %d tokens (%d allowed)
Q#%s|skipped|??|????
**** parsed %s (line number %d)
`,
						label, lineno,
						n, *opt_n,
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
				"LD_LIBRARY_PATH=" + os.Getenv("LD_LIBRARY_PATH"),
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
		xml := setBuild(string(b))
		if *opt_X {
			xml = doExtra(xml)
		}
		if *opt_U || *opt_u != "" {
			s, err := alud.UdAlpino([]byte(xml), filename, "")
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
		`{"request":"parse", "data_type":%q, "timeout":%d, "parser":%q, "max_tokens":%d, "ud":false}`,
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
				if *opt_X {
					line.Alpino_ds = doExtra(line.Alpino_ds)
				}
				if *opt_U || *opt_u != "" {
					s, err := alud.UdAlpino([]byte(line.Alpino_ds), filename, "")
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

func doExtra(data string) string {
	var alpino alpinods.AlpinoDS
	x(xml.Unmarshal([]byte(data), &alpino))

	alpino.Version = alpinods.DtdVersion

	idxnodes = make(map[int]*alpinods.Node)
	prepare(alpino.Node)
	doNp(alpino.Node)
	doVorfeld(alpino.Node)
	doNachfeld(alpino.Node)
	return alpino.String()
}

func prepare(node *alpinods.Node) {
	node.IsNp = ""
	node.IsNachfeld = ""
	node.IsVorfeld = ""
	if node.Index != 0 && (node.Cat != "" || node.Pt != "") {
		idxnodes[node.Index] = node
	}
	if node.Node != nil {
		for _, n := range node.Node {
			prepare(n)
		}
	}
}

func idx(node *alpinods.Node) *alpinods.Node {
	if n, ok := idxnodes[node.Index]; ok {
		return n
	}
	return node
}

// is_np

var (
	rels map[int][]string
)

func doNp(node *alpinods.Node) {
	rels = make(map[int][]string)
	prepareNp1(node)
	prepareNp2(node)
}

func prepareNp1(node *alpinods.Node) {
	if rels[node.ID] == nil {
		rels[node.ID] = make([]string, 0)
	}
	rels[node.ID] = append(rels[node.ID], node.Rel)
	if n := idx(node); n != node {
		if rels[n.ID] == nil {
			rels[n.ID] = make([]string, 0)
		}
		rels[n.ID] = append(rels[n.ID], node.Rel)
	}
	if node.Node != nil {
		for _, n := range node.Node {
			prepareNp1(n)
		}
	}
}

func prepareNp2(node *alpinods.Node) {
	if isNP(node) {
		node.IsNp = "true"
	}
	if node.Node != nil {
		for _, n := range node.Node {
			doNp(n)
		}
	}
}

func isNP(node *alpinods.Node) bool {

	if node.Cat == "" && node.Pt == "" {
		return false
	}

	if node.Cat == "np" {
		return true
	}

	if node.Lcat == "np" && otherString(rels[node.ID], "hd", "mwp") {
		return true
	}

	if node.Pt == "n" && otherString(rels[node.ID], "hd") {
		return true
	}

	if node.Pt == "vnw" && node.Pdtype == "pron" && otherString(rels[node.ID], "hd") {
		return true
	}

	if node.Cat == "mwu" && hasString(rels[node.ID], "su", "obj1", "obj2", "app") {
		return true
	}

	if node.Node != nil {
		for _, n := range node.Node {
			if n.Rel != "cnj" {
				continue
			}
			n = idx(n)
			if isNP(n) {
				return true
			}
		}
	}

	return false
}

// is er een string in ss die gelijk is aan een string in s ?
func hasString(ss []string, s ...string) bool {
	for _, s1 := range ss {
		for _, s2 := range s {
			if s1 == s2 {
				return true
			}
		}
	}
	return false
}

// is er een string in ss die ongelijk is aan alle strings in s ?
func otherString(ss []string, s ...string) bool {
LOOP:
	for _, s1 := range ss {
		for _, s2 := range s {
			if s1 == s2 {
				continue LOOP
			}
		}
		return true
	}
	return false
}

// is_vorfeld

var (
	vorfeld     map[int]map[int]bool
	vorfeldSkip map[int]map[int]bool
)

func doVorfeld(node *alpinods.Node) {
	vorfeld = make(map[int]map[int]bool)
	vorfeldSkip = make(map[int]map[int]bool)
	prepareVorfeld1(node)
	prepareVorfeld2(node)
	prepareVorfeld3(node)
}

func prepareVorfeld1(node *alpinods.Node) {
	vorfeld[node.ID] = make(map[int]bool)
	vorfeldSkip[node.ID] = make(map[int]bool)
	if node.Node != nil {
		for _, n := range node.Node {
			prepareVorfeld1(n)
		}
	}
}

func prepareVorfeld2(node *alpinods.Node) {
	if node.Cat == "smain" {
		smainVorfeld(node)
	} else if node.Cat == "whq" {
		whqVorfeld(node)
	}
	if node.Node != nil {
		for _, n := range node.Node {
			prepareVorfeld2(n)
		}
	}
}

func prepareVorfeld3(node *alpinods.Node) {
	if node.Cat != "" || node.Pt != "" {
		for id := range vorfeld[node.ID] {
			if !vorfeldSkip[node.ID][id] {
				node.IsVorfeld = "true"
				break
			}
		}
	}
	if node.Node != nil {
		for _, n := range node.Node {
			prepareVorfeld3(n)
		}
	}
}

func smainVorfeld(node *alpinods.Node) {
	if node.Node != nil {
		for _, n := range node.Node {
			if n.Rel == "hd" {
				// NIET alleen primary links
				n = idx(n)
				if n.Word != "" {
					for _, topic := range findTopic(node, n.Begin, false) {
						if checkTopic(topic, node, n.Begin) {
							vorfeld[topic.ID][node.ID] = true
						} else {
							vorfeldSkip[topic.ID][node.ID] = true
						}
					}
				}
			}
		}
	}
}

func whqVorfeld(node *alpinods.Node) {
	if node.Node != nil {
		for _, n := range node.Node {
			// NIET alleen primary links
			rel := n.Rel
			n = idx(n)
			if rel == "body" && n.Cat == "sv1" {
				for _, n2 := range n.Node {
					if n2.Rel == "hd" {
						// NIET alleen primary links
						n2 = idx(n2)
						if n2.Word != "" {
							for _, topic := range findTopic(node, n2.Begin, true) {
								if checkTopic(topic, node, n2.Begin) {
									vorfeld[topic.ID][node.ID] = true
								} else {
									vorfeldSkip[topic.ID][node.ID] = true
								}
							}
						}
					}
				}
			}
		}
	}
}

func findTopic(node *alpinods.Node, begin int, needWhd bool) []*alpinods.Node {
	topics := make([]*alpinods.Node, 0)

	// hier: inclusief topnode
	//if isTopic(node, begin) {
	//      topics = append(topics, node)
	//}

	if node.Node != nil {
		for _, n := range node.Node {

			if needWhd && n.Rel != "whd" {
				continue
			}

			// hier: exclusief topnode
			if isTopic(n, begin) {
				topics = append(topics, n)
			}

			// ALLEEN primary links
			for _, topic := range findTopic(n, begin, false) {
				topics = append(topics, topic)
			}
		}
	}
	return topics
}

func isTopic(node *alpinods.Node, begin int) bool {
	if node.Begin < begin && node.End <= begin {
		return true
	}
	if node.Lemma != "" || node.Cat == "mwu" {
		if node.Begin < begin {
			return true
		}
		return false
	}

	if node.Node != nil {
		for _, n := range node.Node {
			if n.Rel == "hd" || n.Rel == "cmp" || n.Rel == "crd" {
				// NIET alleen primary links
				n = idx(n)
				if (n.Lemma != "" || n.Cat == "mwu") && n.Begin < begin {
					return true
				}
			}
		}
	}
	return false
}

func checkTopic(topic, node *alpinods.Node, begin int) bool {
	// alle nodes tussen node (exclusief) en topic (exclusief)
	nodes := make(map[*alpinods.Node]bool)
	nodePath(node, topic, nodes)

	for n := range nodes {
		if isTopic(n, begin) {
			return false
		}
	}

	return true
}

func nodePath(top, bottom *alpinods.Node, nodes map[*alpinods.Node]bool) bool {
	retval := false
	if top.Node != nil {
		for _, n := range top.Node {
			// NIET alleen primaire links (lijkt niet logisch, maar werkt toch beter)
			n = idx(n)
			if n == bottom {
				retval = true
			} else if nodePath(n, bottom, nodes) {
				nodes[n] = true
				retval = true
			}
		}
	}
	return retval
}

// is_nachfeld

var (
	vpCat = map[string]bool{
		"inf":   true,
		"ti":    true,
		"ssub":  true,
		"oti":   true,
		"ppart": true,
	}
	rHead = map[string]bool{
		"hd":   true,
		"cmp":  true,
		"mwp":  true,
		"crd":  true,
		"rhd":  true,
		"whd":  true,
		"nucl": true,
		"dp":   true,
	}
)

func doNachfeld(node *alpinods.Node) {
	findNachfeld(node)
	if node.Node != nil {
		for _, n := range node.Node {
			doNachfeld(n)
		}
	}
}

func findNachfeld(node *alpinods.Node) {
	if !vpCat[node.Cat] {
		return
	}

	if node.Node == nil || len(node.Node) == 0 {
		return
	}

	// zoek begin van head
	headBegin := -1
	for _, n := range node.Node {
		n = idx(n)
		if n.Rel == "hd" {
			headBegin = n.Begin
			break
		}
	}
	if headBegin < 0 {
		return
	}
	// zoek nachfeld
	for _, n := range node.Node {
		n = idx(n)
		if n.Rel != "hd" {
			setNachfeld(n, headBegin, node)
		}
	}
}

func setNachfeld(node *alpinods.Node, begin int, vp *alpinods.Node) {
	if node.Rel != "hd" && node.Rel != "svp" && node.Cat != "inf" && node.Cat != "ppart" {

		var skip func(*alpinods.Node, int) bool
		skip = func(current *alpinods.Node, state int) bool {
			current = idx(current)
			if current == node {
				return state == 2
			}
			switch state {
			case 0:
				state = 1
			case 1:
				if vpCat[current.Cat] {
					state = 2
				}
			}
			if current.Node != nil {
				for _, n := range current.Node {
					if skip(n, state) {
						return true
					}
				}
			}
			return false
		}

		n2 := make([]*alpinods.Node, 0)
		if node.Node != nil {
			for _, nn2 := range node.Node {
				nn2 = idx(nn2)
				if rHead[nn2.Rel] {
					n2 = append(n2, nn2)
				}
			}
		}
		if len(n2) == 0 {
			if begin < node.Begin {
				if !skip(vp, 0) {
					node.IsNachfeld = "true"
				}
				return
			}
		} else {
			ok := true
			for _, nn2 := range n2 {
				if begin >= nn2.Begin {
					ok = false
					break
				}
			}
			if ok {
				if !skip(vp, 0) {
					node.IsNachfeld = "true"
				}
				return
			}
		}
	}
	if vpCat[node.Cat] {
		return
	}
	if node.Node == nil {
		return
	}
	for _, n := range node.Node {
		n = idx(n)
		setNachfeld(n, begin, vp)
	}
}

// END is_np, is_vorfeld, is_nachfeld
//
////////////////
