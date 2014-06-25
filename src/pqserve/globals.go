package main

//. Imports

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
	"net"
	"regexp"
	"sync"
	"time"
)

//. Types

type Config struct {
	Port int
	Url  string

	Default string

	Mailfrom string
	Smtpserv string
	Smtpuser string
	Smtppass string

	Login  string
	Prefix string

	Maxjob int
	Maxwrd int
	Maxmem int64 // voor het parsen post data, in bytes

	Secret string

	Sh      string
	Path    string
	Alpino  string
	Timeout int

	Https     bool
	Httpdual  bool
	Remote    bool
	Forwarded bool

	Querytimeout int // in secondes

	Directories []string

	View   []ViewType
	Access []AccessType
}

type ViewType struct {
	Allow bool
	Addr  []string
	all   bool
	ip    []net.IP
	ipnet []*net.IPNet
}

type AccessType struct {
	Allow bool
	Mail  []string
	all   bool
	re    []*regexp.Regexp
}

// een dependency relation, geretourneerd door SQL
type Row struct {
	word    string
	lemma   string
	postag  string
	rel     string
	hpostag string
	hlemma  string
	hword   string
	begin   int
	end     int
	hbegin  int
	hend    int
	mark    string
}

type Sentence struct {
	arch  int
	file  int
	words []string // de zin opgesplitst in woorden
	items []Row    // alle matchende dependency relations voor deze zin
}

type Alpino_ds struct {
	XMLName  xml.Name `xml:"alpino_ds"`
	Node0    *Node    `xml:"node"`
	Sentence string   `xml:"sentence"`
}

type Node struct {
	Id       string  `xml:"id,attr"`
	Index    string  `xml:"index,attr"`
	Cat      string  `xml:"cat,attr"`
	Pt       string  `xml:"pt,attr"`
	Word     string  `xml:"word,attr"`
	Lemma    string  `xml:"lemma,attr"`
	Postag   string  `xml:"postag,attr"`
	Rel      string  `xml:"rel,attr"`
	Begin    int     `xml:"begin,attr"`
	End      int     `xml:"end,attr"`
	NodeList []*Node `xml:"node"`
	skip     bool
}

type TreeContext struct {
	yellow map[int]bool
	green  map[int]bool
	marks  map[string]bool
	refs   map[string]bool
	graph  bytes.Buffer // definitie dot-bestand
	start  int
	words  []string
}

type Process struct {
	id     string
	nr     int
	chKill chan bool
	killed bool
	queued bool
	lock   sync.Mutex
}

type ProcessMap map[string]*Process

//. Constanten

const (
	MAXTITLELEN = 64
	ZINMAX      = 10
	WRDMAX      = 250
	YELLOW      = "<span style=\"background-color: yellow\">"
	GREEN       = "<span style=\"background-color: lightgreen\">"
	YELGRN      = "<span style=\"background-color: cyan\">"
)

//. Variabelen

var (
	paqudir    string
	cookiepath string

	tlsConfig = &tls.Config{}

	Cfg     Config
	verbose bool

	chLog = make(chan string)

	chWork    = make(chan *Process, 10000) // gebufferd om er voor te zorgen dat FIFO
	processes = make(map[string]*Process)

	dirnameLock sync.Mutex
	quotumLock  sync.Mutex
	processLock sync.RWMutex

	reQuote = regexp.MustCompile("%[0-9A-F-a-f][0-9A-F-a-f]") // voor urldecode()
	reMail  = regexp.MustCompile("^[-.a-z0-9!#$%&'*+/=?^_`{|}~]+@[-.a-z0-9]+$")
	reNoAz  = regexp.MustCompile("[^a-z]+")

	opt_postag  = []string{"", "adj", "bw", "let", "lid", "mwu", "n", "spec", "tsw", "tw", "vg", "vnw", "vz", "ww"}
	opt_hpostag = []string{"", "adj", "bw", "let", "lid", "mwu", "n", "spec", "tsw", "tw", "vg", "vnw", "vz", "ww"}
	opt_rel     = []string{"", "Aapp", "Adet", "Ahdf", "Ald", "Ame", "Amod", "Aobcomp", "Aobj1", "Aobj2", "Apc", "Apobj1",
		"Apredc", "Apredm", "Ase", "Asu", "Asup", "Asvp", "Avc", "Bbody/cmp", "Bbody/rhd", "Bbody/whd", "Bcnj/cnj",
		"Bcnj/crd", "Bcrd/cnj", "Bcrd/crd", "Bmod/cmp", "Bnucl/dlink", "Bobj1/su", "Bobj2/su", "Bsu/obj1", "Bsu/obj2",
		"Ccmp/-", "Cdlink/-", "Cdp/-", "Chd/-", "Cnucl/-", "Csat/-", "Ctag/-", "Dapp/cnj", "Dapp/crd", "Dapp/mod",
		"Dcnj/app", "Dcnj/det", "Dcnj/mod", "Dcrd/-", "Dcrd/app", "Dcrd/det", "Dcrd/mod", "Ddet/-", "Ddet/cnj",
		"Ddet/crd", "Ddet/mod", "Dmod/-", "Dmod/app", "Dmod/cnj", "Dmod/crd", "Dmod/det", "Dmod/mod", "Dmod/rhd",
		"Dmod/whd", "Dmwp/-", "Dobj1/-", "Dsat/dlink", "Dsu/-", "Dtag/dlink", "Dwhd/-"}

	ConnectionClosed = errors.New("Verbinding gesloten")

	versionstring string
	version       [3]int

	hasMaxStatementTime bool

	taskWaitNr int
	taskWorkNr int

	started = time.Now()

	wg           sync.WaitGroup
	wgLogger     sync.WaitGroup
	chGlobalExit = make(chan bool)
	chLoggerExit = make(chan bool)
)

func (p ProcessMap) String() string {
	processLock.RLock()
	defer processLock.RUnlock()
	var buf bytes.Buffer
	var comma string
	fmt.Fprint(&buf, "[")
	for key, val := range p {
		st := "working"
		if val.killed {
			st = "killed"
		} else if val.queued {
			st = "queued"
		}
		fmt.Fprintf(&buf, "%s{\"id\":%q,\"status\":%q}", comma, key, st)
		comma = ","
	}
	fmt.Fprint(&buf, "]")
	return buf.String()
}