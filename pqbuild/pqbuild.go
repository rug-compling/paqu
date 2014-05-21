package main

//. Imports

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/pebbe/compactcorpus"
	"github.com/pebbe/dbxml"
	"github.com/pebbe/util"

	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"
)

//. Types

type Alpino_ds struct {
	XMLName  xml.Name `xml:"alpino_ds"`
	Node0    *Node    `xml:"node"`
	Sentence string   `xml:"sentence"`
}

type Node struct {
	Id       string  `xml:"id,attr"`
	Word     string  `xml:"word,attr"`
	Lemma    string  `xml:"lemma,attr"`
	Root     string  `xml:"root,attr"`
	Pos      string  `xml:"pos,attr"`
	Postag   string  `xml:"pt,attr"` // pt i.p.v. postag
	Rel      string  `xml:"rel,attr"`
	Cat      string  `xml:"cat,attr"`
	Begin    int     `xml:"begin,attr"`
	End      int     `xml:"end,attr"`
	Index    int     `xml:"index,attr"`
	NodeList []*Node `xml:"node"`
	used     bool
}

// Een node met het pad naar de node
type NodePath struct {
	node *Node
	path []string
}

// Een dependency relation
type Deprel struct {
	word, lemma, root, postag, rel      string
	hword, hlemma, hroot, hpostag, hrel string
	begin, hbegin, end, hend            int
	mark                                string
}

type Deprels []*Deprel

type Config struct {
	Login  string
	Prefix string
}

//. Constanten

const (
	DPRL = iota
	SENT
	WORD
	FILE
	ARCH
)

//. Variabelen

var (
	db *sql.DB

	lineno = 0

	refnodes []*Node   // reset per zin
	deprels  []*Deprel // reset per zin

	targets = []string{"hd", "cmp", "crd", "dlink", "rhd", "whd"}

	configfile string
	prefix     string
	desc       string
	owner      string
	public     string

	topfile  = -1
	toparch  = -1
	lastarch string

	buffer       [5]bytes.Buffer
	buf_has_data [5]bool

	db_append    bool
	db_overwrite bool
	db_exists    bool
	db_makeindex bool

	memstats runtime.MemStats

	Cfg Config
)

//. Main

func main() {

	now := time.Now()

	buffer[DPRL].Grow(50000)
	buffer[SENT].Grow(50000)
	buffer[WORD].Grow(50000)
	buffer[FILE].Grow(50000)
	buffer[ARCH].Grow(50000)

	db_makeindex = true
	for len(os.Args) > 1 {
		if os.Args[1] == "-a" {
			db_append = true
			db_overwrite = false
		} else if os.Args[1] == "-w" {
			db_append = false
			db_overwrite = true
		} else if os.Args[1] == "-i" {
			db_makeindex = false
		} else {
			break
		}
		os.Args = append(os.Args[:1], os.Args[2:]...)
	}

	if len(os.Args) != 6 || util.IsTerminal(os.Stdin) {
		fmt.Printf(`
Syntax: %s [-a] [-w] [-i] configfile id description owner public < bestandnamen

Opties:

 -a : toevoegen aan bestaande database
 -w : bestaande database overschrijven
 -i : geen tabel van woord naar lemmas aanmaken

  configfile:
  id:
  description:
  owner:       e-mail address
  public:      0 (private) or 1 (public)


`, os.Args[0])
		return
	}

	configfile = strings.TrimSpace(os.Args[1])
	prefix = strings.TrimSpace(os.Args[2])
	desc = strings.TrimSpace(os.Args[3])
	owner = strings.TrimSpace(os.Args[4])
	public = strings.TrimSpace(os.Args[5])

	data, err := ioutil.ReadFile(configfile)
	util.CheckErr(err)
	util.CheckErr(json.Unmarshal(data, &Cfg))

	if desc == "" {
		util.CheckErr(fmt.Errorf("De omschrijving mag niet leeg zijn"))
	}

	if owner == "" {
		util.CheckErr(fmt.Errorf("De eigenaar mag niet leeg zijn"))
	}

	if prefix == "" {
		util.CheckErr(fmt.Errorf("De id mag niet leeg zijn"))
	}

	for _, c := range prefix {
		if c < 'a' || c > 'z' {
			util.CheckErr(fmt.Errorf("Ongeldige tekens in '%s'. Alleen kleine letters a tot z mogen.", prefix))
		}
	}

	db = connect()
	defer func() {
		fmt.Println("Closing database connection...")
		util.CheckErr(db.Close())
	}()

	// kijk of de database al bestaat
	rows, err := db.Query("SELECT `begin` FROM `" + Cfg.Prefix + "_c_" + prefix + "_deprel` LIMIT 0, 1;")
	if err == nil && rows.Next() {
		rows.Close()
		if !(db_append || db_overwrite) {
			util.CheckErr(fmt.Errorf("Database exists and no -a or -w specified"))
		}
		db_exists = true

		if db_append {
			rows, err := db.Query("SELECT MAX(id) FROM " + Cfg.Prefix + "_c_" + prefix + "_arch")
			util.CheckErr(err)
			if rows.Next() {
				if rows.Scan(&toparch) != nil {
					toparch = -1
				}
				rows.Close()
			}
			rows, err = db.Query("SELECT MAX(id) FROM " + Cfg.Prefix + "_c_" + prefix + "_file")
			util.CheckErr(err)
			if rows.Next() {
				if rows.Scan(&topfile) != nil {
					topfile = -1
				}
				rows.Close()
			}
			fmt.Println("Verwijderen indexen uit " + Cfg.Prefix + "_c_" + prefix + "_deprel ...")
			db.Exec(`ALTER TABLE ` + Cfg.Prefix + "_c_" + prefix + `_deprel
				DROP INDEX word,
				DROP INDEX lemma,
				DROP INDEX root,
				DROP INDEX postag,
				DROP INDEX rel,
				DROP INDEX hword,
				DROP INDEX hlemma,
				DROP INDEX hroot,
				DROP INDEX hpostag,
				DROP INDEX file,
				DROP INDEX arch;`)
			fmt.Println("Verwijderen indexen uit " + Cfg.Prefix + "_c_" + prefix + "_sent ...")
			db.Exec(`ALTER TABLE ` + Cfg.Prefix + "_c_" + prefix + `_sent
				DROP INDEX file,
				DROP INDEX arch;`)
			fmt.Println("Verwijderen index uit " + Cfg.Prefix + "_c_" + prefix + "_file ...")
			db.Exec(`ALTER TABLE ` + Cfg.Prefix + "_c_" + prefix + `_file
				DROP INDEX id`)
			fmt.Println("Verwijderen index uit " + Cfg.Prefix + "_c_" + prefix + "_arch ...")
			db.Exec(`ALTER TABLE ` + Cfg.Prefix + "_c_" + prefix + `_arch
				DROP INDEX id;`)
		}
	}

	// set up database

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS ` + Cfg.Prefix + `_info (
		id          varchar(128) NOT NULL,
		description varchar(128) NOT NULL COLLATE utf8_unicode_ci,
		owner       varchar(128) NOT NULL DEFAULT 'none',
		status      enum('QUEUED','WORKING','FINISHED','FAILED') NOT NULL DEFAULT 'QUEUED',
		msg         varchar(256) NOT NULL,
		nline       int          NOT NULL DEFAULT 0,
		params      varchar(128) NOT NULL,
		shared      enum('PRIVATE','PUBLIC','SHARED') NOT NULL DEFAULT 'PRIVATE',
		created     timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
		active      datetime     NOT NULL,
		UNIQUE INDEX (id),
		INDEX (description),
		INDEX (owner),
		INDEX (status),
		INDEX (created),
		INDEX (active))
		DEFAULT CHARACTER SET utf8;`)
	util.CheckErr(err)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS ` + Cfg.Prefix + `_corpora (
		user    varchar(64) NOT NULL,
		prefix  varchar(64) NOT NULL,
		enabled tinyint     NOT NULL DEFAULT 1,
		INDEX (user),
		INDEX (prefix),
		INDEX (enabled))
		DEFAULT CHARACTER SET utf8;`)
	util.CheckErr(err)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS ` + Cfg.Prefix + `_users (
		mail   varchar(64) NOT NULL,
		pw     char(16)    NOT NULL,
		active datetime    NOT NULL,
		UNIQUE INDEX (mail),
		INDEX (active))
		DEFAULT CHARACTER SET utf8;`)
	util.CheckErr(err)

	share := "PRIVATE"
	if public == "1" {
		share = "PUBLIC"
	}

	db.Exec(fmt.Sprintf("INSERT `%s_info` (`id`) VALUES (%q);", Cfg.Prefix, prefix)) // negeer fout
	_, err = db.Exec(fmt.Sprintf("UPDATE `%s_info` SET `description` = %q, `owner` = %q, `status` = \"WORKING\", `shared` = %q WHERE `id` = %q",
		Cfg.Prefix, desc, owner, share, prefix))
	util.CheckErr(err)

	_, err = db.Exec("DELETE FROM " + Cfg.Prefix + "_corpora WHERE `prefix` = \"" + prefix + "\";")
	util.CheckErr(err)

	// tabellen <prefix>_deprel en <prefix>_sent aanmaken
	if !db_exists || db_overwrite {
		_, err := db.Exec(fmt.Sprintf(
			"DROP TABLE IF EXISTS `%s_c_%s_deprel`, `%s_c_%s_sent`, `%s_c_%s_file`, `%s_c_%s_arch`; ",
			Cfg.Prefix,
			prefix,
			Cfg.Prefix,
			prefix,
			Cfg.Prefix,
			prefix,
			Cfg.Prefix,
			prefix))
		util.CheckErr(err)
		_, err = db.Exec(`CREATE TABLE ` + Cfg.Prefix + "_c_" + prefix + `_deprel (
			word    varchar(128) NOT NULL,
			lemma   varchar(128) NOT NULL,
			root    varchar(128) NOT NULL,
			postag  varchar(64)  NOT NULL,
			rel     varchar(64)  NOT NULL,
			hword   varchar(128) NOT NULL,
			hlemma  varchar(128) NOT NULL,
			hroot   varchar(128) NOT NULL,
			hpostag varchar(64)  NOT NULL,
			arch    int          NOT NULL,
			file    int          NOT NULL,
			begin   int          NOT NULL,
			end     int          NOT NULL,
			hbegin  int          NOT NULL,
			hend    int          NOT NULL,
			mark    varchar(128) NOT NULL)
			DEFAULT CHARACTER SET utf8
			DEFAULT COLLATE utf8_unicode_ci;`)

		util.CheckErr(err)
		_, err = db.Exec(`CREATE TABLE ` + Cfg.Prefix + "_c_" + prefix + `_sent (
			arch int  NOT NULL,
			file int  NOT NULL,
			sent text NOT NULL)
			DEFAULT CHARACTER SET utf8
			DEFAULT COLLATE utf8_unicode_ci;`)
		util.CheckErr(err)
		_, err = db.Exec(`CREATE TABLE ` + Cfg.Prefix + "_c_" + prefix + `_file (
			id   int          NOT NULL,
			file varchar(260) NOT NULL)
			DEFAULT CHARACTER SET utf8;`)
		util.CheckErr(err)
		_, err = db.Exec(`CREATE TABLE ` + Cfg.Prefix + "_c_" + prefix + `_arch (
			id   int          NOT NULL,
			arch varchar(260) NOT NULL)
			DEFAULT CHARACTER SET utf8;`)
		util.CheckErr(err)
	}

	/*
		Bestandnamen van stdin inlezen en verwerken.
	*/
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		filename := strings.TrimSpace(scanner.Text())
		if filename == "" {
			continue
		}
		var err error
		filename, err = filepath.Abs(filename)
		util.CheckErr(err)
		filename = filepath.Clean(filename)
		if strings.HasSuffix(filename, ".xml") {
			data, err := ioutil.ReadFile(filename)
			util.CheckErr(err)
			do_data("", filename, data)
		} else if strings.HasSuffix(filename, ".dact") {
			reader, err := dbxml.Open(filename)
			util.CheckErr(err)
			fmt.Println(">>>", filename)
			docs, err := reader.All()
			util.CheckErr(err)
			for docs.Next() {
				do_data(filename, docs.Name(), []byte(docs.Content()))
			}
			showmemstats()
			reader.Close()
		} else if strings.HasSuffix(filename, ".data.dz") {
			reader, err := compactcorpus.Open(filename)
			util.CheckErr(err)
			fmt.Println(">>>", filename)
			docs, err := reader.NewRange()
			util.CheckErr(err)
			for docs.HasNext() {
				name, xml := docs.Next()
				do_data(filename, name, xml)
			}
			showmemstats()
		} else {
			util.CheckErr(fmt.Errorf("Invalid suffix for file '%s'", filename))
		}
	}
	util.CheckErr(scanner.Err())

	// stuur laatste data uit buffers naar de database
	buf_flush(DPRL)
	buf_flush(SENT)
	buf_flush(FILE)
	buf_flush(ARCH)
	db.Exec("COMMIT;")

	fmt.Println("Tijd:", time.Now().Sub(now))

	if !db_makeindex {
		sizes()
		return
	}

	fmt.Println("Aanmaken indexen op " + Cfg.Prefix + "_c_" + prefix + "_deprel ...")
	_, err = db.Exec(`ALTER TABLE ` + Cfg.Prefix + "_c_" + prefix + `_deprel
		ADD INDEX (word),
		ADD INDEX (lemma),
		ADD INDEX (root),
		ADD INDEX (postag),
		ADD INDEX (rel),
		ADD INDEX (hword),
		ADD INDEX (hlemma),
		ADD INDEX (hroot),
		ADD INDEX (hpostag),
		ADD INDEX (file),
		ADD INDEX (arch);`)
	util.CheckErr(err)

	fmt.Println("Aanmaken indexen op " + Cfg.Prefix + "_c_" + prefix + "_sent ...")
	_, err = db.Exec(`ALTER TABLE ` + Cfg.Prefix + "_c_" + prefix + `_sent
		ADD INDEX (file),
		ADD INDEX (arch);`)
	util.CheckErr(err)

	fmt.Println("Aanmaken index op " + Cfg.Prefix + "_c_" + prefix + "_file ...")
	_, err = db.Exec(`ALTER TABLE ` + Cfg.Prefix + "_c_" + prefix + `_file
		ADD UNIQUE INDEX (id)`)
	util.CheckErr(err)

	fmt.Println("Aanmaken index op " + Cfg.Prefix + "_c_" + prefix + "_arch ...")
	_, err = db.Exec(`ALTER TABLE ` + Cfg.Prefix + "_c_" + prefix + `_arch
		ADD UNIQUE INDEX (id);`)
	util.CheckErr(err)

	// tijd voor aanmaken tabellen <prefix>_deprel en <prefix>_sent
	fmt.Println("Tijd:", time.Now().Sub(now))

	showmemstats()

	// tabel <prefix>_word aanmaken
	_, err = db.Exec(fmt.Sprintf(
		"DROP TABLE IF EXISTS `%s_c_%s_word`;",
		Cfg.Prefix,
		prefix))
	util.CheckErr(err)
	_, err = db.Exec(`CREATE TABLE ` + Cfg.Prefix + "_c_" + prefix + `_word (
		word  varchar(128) NOT NULL,
		lemma varchar(1024) NOT NULL)
		DEFAULT CHARACTER SET utf8
		DEFAULT COLLATE utf8_unicode_ci;`)
	util.CheckErr(err)

	/*
		Inlezen van woorden die kandidaat zijn voor het zoeken via een lemma.
		De lijst moet in één keer ingelezen worden, omdat er anders door een time-out
		slechts een klein deel van de woorden wordt verwerkt.
		Nu maar hopen dat de complete woordenlijst in het geheugen past.
	*/
	fmt.Println("Tellingen van woorden opvragen ...")
	rows, err = db.Query("SELECT count(*), `word` FROM `" + Cfg.Prefix + "_c_" + prefix +
		"_deprel` WHERE `postag` IN (\"adj\", \"n\", \"ww\") GROUP BY `word` HAVING count(*) >= 10 ORDER BY `word`")
	util.CheckErr(err)
	woorden := make([]string, 0)
	var woord string
	for rows.Next() {
		var i int
		util.CheckErr(rows.Scan(&i, &woord))
		woorden = append(woorden, woord)
	}
	util.CheckErr(rows.Err())

	// zoek de lemma's bij elk woord
	fmt.Println("Zoeken naar lemma's bij woorden ...")
	for idx, woord := range woorden {
		var s, p string

		if n := len(woorden) - idx; n%100 == 0 {
			fmt.Printf(" %d   \r", n)
		}

		lemmas := make([]string, 0)

		/*
			word -> lemma
			Deze stap is simpel: kijk voor elk woord met welk lemma het voorkomt.
			Dit werkt prima voor LassyDevelop
		*/
		rows, err := db.Query(fmt.Sprintf(
			"SELECT `lemma` FROM `"+Cfg.Prefix+"_c_"+prefix+"_deprel` WHERE `word` = %q GROUP BY `lemma`;",
			woord))
		util.CheckErr(err)
		for rows.Next() {
			util.CheckErr(rows.Scan(&s))
			lemmas = append(lemmas, s)
		}
		util.CheckErr(rows.Err())

		/*
			word -> root+postag -> lemma
			In LassyLarge zijn vaak geen goede lemma's opgenomen. Het woord 'mannen' geeft lemma 'mannen'.
			De oplossing hier is te zoeken via root. Het woord 'mannen' geeft root 'man'. De root 'man' geeft
			lemma's 'man' en 'mannen'.
			Gevonden roots worden alleen gebruikt als ze ook dezelfde postag hebben. Dit voorkomt dat je voor het
			woord 'fietst' (ww) het lemma 'fiets' (n) krijgt.
		*/
		roots := make([][2]string, 0)
		rows, err = db.Query(fmt.Sprintf(
			"SELECT `root`,`postag` FROM `"+Cfg.Prefix+"_c_"+prefix+"_deprel` WHERE `word` = %q GROUP BY `root`,`postag`;",
			woord))
		util.CheckErr(err)
		for rows.Next() {
			util.CheckErr(rows.Scan(&s, &p))
			roots = append(roots, [2]string{s, p})
		}
		util.CheckErr(rows.Err())
		for _, root := range roots {
			rows, err := db.Query(fmt.Sprintf(
				"SELECT `lemma` FROM `"+Cfg.Prefix+"_c_"+prefix+"_deprel` WHERE `root` = %q AND `postag` = %q GROUP BY `lemma`;",
				root[0], root[1]))
			util.CheckErr(err)
			for rows.Next() {
				util.CheckErr(rows.Scan(&s))
				if !has(lemmas, s) {
					lemmas = append(lemmas, s)
				}
			}
			util.CheckErr(rows.Err())
		}

		/* stuur woord met lemma's naar de databasebuffer */
		sort.Strings(lemmas)
		word_buf_put(woord, lemmas)
	}

	// stuur laatste data uit buffer naar de database
	buf_flush(WORD)
	db.Exec("COMMIT;")

	fmt.Println("Aanmaken index op " + Cfg.Prefix + "_c_" + prefix + "_word ...")
	_, err = db.Exec(`ALTER TABLE ` + Cfg.Prefix + "_c_" + prefix + `_word
		ADD UNIQUE INDEX (word);`)
	util.CheckErr(err)

	// zet info over corpus in de database

	lines := 0
	rows, err = db.Query("SELECT COUNT(*) FROM " + Cfg.Prefix + "_c_" + prefix + "_sent")
	util.CheckErr(err)
	if rows.Next() {
		util.CheckErr(rows.Scan(&lines))
		rows.Close()
	}
	_, err = db.Exec(fmt.Sprintf("UPDATE `%s_info` SET `status` = \"FINISHED\", `nline` = %d, `active` = NOW() WHERE `id` = %q",
		Cfg.Prefix, lines, prefix))
	util.CheckErr(err)

	user := owner
	if public == "1" {
		user = "all"
	}
	_, err = db.Exec(fmt.Sprintf("INSERT `%s_corpora` (`user`, `prefix`) VALUES (%q, %q);", Cfg.Prefix, user, prefix))
	util.CheckErr(err)

	//fmt.Println("Bijwerken menu's voor postag, rel en hpostag ...")
	//tags()

	// totale tijd
	fmt.Println("Tijd:", time.Now().Sub(now))
	showmemstats()

	sizes()
}

//. Verwerking

// Verwerk een enkel xml-bestand
func do_data(archname, filename string, data []byte) {

	arch := -1
	if archname != "" {
		if archname != lastarch {
			toparch++
			lastarch = archname
			arch_buf_put(archname, toparch)
		}
		arch = toparch
	}

	topfile++
	file := topfile
	file_buf_put(filename, file)

	lineno++
	if lineno%1000 == 0 {
		fmt.Println(lineno, filename)
	}

	alpino := Alpino_ds{}
	err := xml.Unmarshal(data, &alpino)
	util.CheckErr(err)

	// zin opslaan in <prefix>_sent
	sent_buf_put(arch, file, alpino.Sentence)

	// multi-word units "ineenvouwen"
	mwu(alpino.Node0)

	/*
		Zoek alle referenties. Dit zijn nodes met een attribuut "index".
		Sla deze op in een tabel: index -> *Node
	*/
	refnodes = make([]*Node, len(strings.Fields(alpino.Sentence)))
	prepare(alpino.Node0)

	/*
		Zoek alle dependency relations, en sla die op in de tabel 'deprels'
	*/
	deprels = make([]*Deprel, 0, len(strings.Fields(alpino.Sentence)))
	traverse(alpino.Node0)

	/*
		Sla alle resterende woorden, waarvoor geen dependency relation is gevonden op in de tabel 'deprels'.
	*/
	traverse2(alpino.Node0)

	/*
		Sorteer de lijst met dependency relations.
		Stuur elk naar de databasebuffer, waarbij dubbelen worden overgeslagen.
	*/
	sort.Sort(Deprels(deprels))
	for i, d := range deprels {
		d1 := d
		if i != 0 {
			d1 = deprels[i-1]
		}
		rel := d.rel + "/" + d.hrel
		if d.hrel == "" {
			rel = d.rel + "/-"
		} else if d.hrel == "hd" {
			rel = d.rel
		}
		if i == 0 || d.begin != d1.begin || d.hbegin != d1.hbegin || d.rel != d1.rel || d.hrel != d1.hrel || d.mark != d1.mark {
			deprel_buf_put(arch, file, rel, d)
		}
	}
}

// Multi-word units "ineenvouwen".
// Grotendeels dezelfde code als in lassytree.go. Zie dat bestand voor beschrijving van verschillen.
func mwu(node *Node) {
	if node.Cat == "mwu" {
		/*
			Voorwaardes:
			- De dochters moeten op elkaar aansluiten, en heel het bereik van de parentnode beslaan.
			- Er mogen geen indexen in gebruikt zijn (met of zonder inhoud)
		*/
		p1 := node.Begin
		p2 := node.End
		ok := true
		for _, n := range node.NodeList {
			if n.Index > 0 {
				ok = false
				break
			}
			if n.Begin != p1 {
				ok = false
				break
			}
			p1 = n.End
		}
		if ok && p1 == p2 {
			node.Cat = ""
			node.Postag = "mwu"
			words := make([]string, 0, node.End-node.Begin)
			lemmas := make([]string, 0, node.End-node.Begin)
			roots := make([]string, 0, node.End-node.Begin)
			for _, n := range node.NodeList {
				words = append(words, n.Word)
				lemmas = append(lemmas, n.Lemma)
				roots = append(roots, n.Root)
			}
			node.Word = strings.Join(words, " ")
			node.Lemma = strings.Join(lemmas, " ")
			node.Root = strings.Join(roots, " ")
			node.NodeList = node.NodeList[0:0]
		}
	}
	for _, n := range node.NodeList {
		mwu(n)
	}

}

// Zoek alle referenties. Dit zijn nodes met een attribuut "index".
// Sla deze op in een tabel 'refnames': index -> *Node
func prepare(node *Node) {
	if node.Index > 0 && (len(node.Word) != 0 || len(node.NodeList) != 0) {
		for len(refnodes) <= node.Index {
			refnodes = append(refnodes, nil)
		}
		refnodes[node.Index] = node
	}
	for _, n := range node.NodeList {
		prepare(n)
	}
}

// Zoek alle dependency relations, en sla die op in de tabel 'deprels'
func traverse(node *Node) {
	if len(node.NodeList) == 0 {
		return
	}

	// Zoek hoofd-dochter. Dit is de eerste van 'targets': hd cmp crd dlink rhd whd
	idx := -1
TARGET:
	for _, target := range targets {
		for i, n := range node.NodeList {
			if n.Rel == target {
				idx = i
				break TARGET
			}
		}
	}
	if idx >= 0 {
		heads := find_head(node.NodeList[idx])
		for i, n := range node.NodeList {
			if i == idx {
				continue
			}
			for _, np2 := range find_head(n) {
				n2 := np2.node
				for _, headpath := range heads {
					head := headpath.node
					lassy_deprel(n2.Word, n2.Lemma, n2.Root, n2.Postag, n.Rel, // n.Rel, dus niet n2.Rel !
						head.Word, head.Lemma, head.Root, head.Postag, node.NodeList[idx].Rel,
						n2.Begin, n2.End, head.Begin, head.End,
						np2.path, headpath.path)
					n2.used = true
				}
			}
		}
	}

	// Zoek su-dochter met obj1-dochter of obj2-dochter
	idx = -1
	for i, n := range node.NodeList {
		if n.Rel == "su" {
			idx = i
			break
		}
	}
	if idx >= 0 {
		subjs := find_head(node.NodeList[idx])
		for _, obj := range node.NodeList {
			if obj.Rel != "obj1" && obj.Rel != "obj2" {
				continue
			}
			for _, op := range find_head(obj) {
				o := op.node
				for _, sup := range subjs {
					su := sup.node
					lassy_deprel(o.Word, o.Lemma, o.Root, o.Postag, obj.Rel,
						su.Word, su.Lemma, su.Root, su.Postag, "su",
						o.Begin, o.End, su.Begin, su.End,
						op.path, sup.path)
					o.used = true
					lassy_deprel(su.Word, su.Lemma, su.Root, su.Postag, "su",
						o.Word, o.Lemma, o.Root, o.Postag, obj.Rel,
						su.Begin, su.End, o.Begin, o.End,
						sup.path, op.path)
					su.used = true
				}
			}
		}
	}

	// cat conj: alles kan head zijn
	if node.Cat == "conj" {
		heads := make([][]*NodePath, len(node.NodeList))
		for i, n1 := range node.NodeList {
			heads[i] = find_head(n1)
		}
		for i := 1; i < len(heads); i++ {
			for j := 0; j < i; j++ {
				for _, np1 := range heads[i] {
					n1 := np1.node
					for _, np2 := range heads[j] {
						n2 := np2.node
						lassy_deprel(n1.Word, n1.Lemma, n1.Root, n1.Postag, node.NodeList[i].Rel,
							n2.Word, n2.Lemma, n2.Root, n2.Postag, node.NodeList[j].Rel,
							n1.Begin, n1.End, n2.Begin, n2.End,
							np1.path, np2.path)
						n1.used = true
						lassy_deprel(n2.Word, n2.Lemma, n2.Root, n2.Postag, node.NodeList[j].Rel,
							n1.Word, n1.Lemma, n1.Root, n1.Postag, node.NodeList[i].Rel,
							n2.Begin, n2.End, n1.Begin, n1.End,
							np2.path, np1.path)
						n2.used = true
					}
				}
			}
		}
	}

	for _, n := range node.NodeList {
		traverse(n)
	}

}

// Sla alle resterende woorden, waarvoor geen dependency relation is gevonden op in de tabel 'deprels'.
func traverse2(node *Node) {
	// negeer woorden met relatie == "--"
	if node.Word != "" && node.Rel != "--" && !node.used {
		lassy_deprel(node.Word, node.Lemma, node.Root, node.Postag, node.Rel,
			"", "", "", "", "", node.Begin, node.End, 0, 0, []string{}, []string{})
	}
	for _, n := range node.NodeList {
		traverse2(n)
	}
}

// Geef een lijst van alle dochters van node die als head kunnen optreden.
// Bij elke dochter, geef ook het pad dat naar die node leidde.
func find_head(node *Node) []*NodePath {
	path := []string{node.Id}

	/*
		Als we bij een index zijn, spring naar de node met de definitie voor deze index.
		(Dat kan de node zelf zijn.)
		De node waarnaar gesprongen wordt wordt niet opgenomen in het pad. Dat is iets
		wat in het programma lassytree opgelost moet worden. Wel opnemen in het pad zorgt
		voor problemen die in het programma lassytree veel moeilijker zijn op te lossen.
	*/
	if node.Index > 0 {
		node = refnodes[node.Index]
	}

	/*
		Als het woord niet leeg is, dan hebben we een terminal bereikt.
	*/
	if node.Word != "" {
		// negeer woorden met relatie == "--"
		if node.Rel == "--" {
			return []*NodePath{}
		}
		return []*NodePath{&NodePath{node: node, path: path}}
	}

	/*
		Als de node categorie "conj" heeft, dan kan elke dochter een head zijn.
		Geef een lijst van de heads van alle dochters.
	*/
	if node.Cat == "conj" {
		nodes := make([]*NodePath, 0, len(node.NodeList))
		for _, n := range node.NodeList {
			for _, n2 := range find_head(n) {
				p2 := make([]string, len(n2.path))
				copy(p2, n2.path)
				for _, p := range path {
					p2 = append(p2, p)
				}
				nodes = append(nodes, &NodePath{node: n2.node, path: p2})
			}
		}
		return nodes
	}

	/*
		Zoek hoofd-dochter. Dit is de eerste van 'targets': hd cmp crd dlink rhd whd
	*/
	for _, target := range targets {
		for _, n := range node.NodeList {
			if n.Rel == target {
				nodes := make([]*NodePath, 0)
				for _, n2 := range find_head(n) {
					p2 := make([]string, len(n2.path))
					copy(p2, n2.path)
					for _, p := range path {
						p2 = append(p2, p)
					}
					nodes = append(nodes, &NodePath{node: n2.node, path: p2})
				}
				return nodes
			}
		}
	}

	// Geen hoofd gevonden: retourneer lege lijst
	return []*NodePath{}
}

// Voeg een dependency relation toe aan de lijst 'deprels'
func lassy_deprel(word, lemma, root, postag, rel, hword, hlemma, hroot, hpostag, hrel string, begin, end, hbegin, hend int, path, hpath []string) {
	/*
		Alle elementen uit paden van woord en hoofdwoord samenvoegen.
		De topnode komt niet in het pad. Voor elke geregistreerde node in het pad
		loopt er een link van die node naar de parent node.
		De lijst met nodes wordt gesorteerd om duplicaten van een dependency relation
		te kunnen vinden.
	*/
	p := ""
	if len(path) > 0 {
		marks := make(map[string]bool)
		for _, i := range path {
			marks[i] = true
		}
		for _, i := range hpath {
			marks[i] = true
		}
		ms := make([]string, 0, len(marks))
		for m := range marks {
			ms = append(ms, m)
		}
		sort.Sort(sort.StringSlice(ms))
		p = strings.Join(ms, ",")
	}

	deprels = append(deprels, &Deprel{
		word:    word,
		lemma:   lemma,
		root:    root,
		postag:  postag,
		rel:     rel,
		hword:   hword,
		hlemma:  hlemma,
		hroot:   hroot,
		hpostag: hpostag,
		hrel:    hrel,
		begin:   begin,
		end:     end,
		hbegin:  hbegin,
		hend:    hend,
		mark:    p,
	})
}

//. Buffering

// Zet een zin in de buffer.
// Als de buffer vol raakt, stuur alles naar de database.
func sent_buf_put(arch, file int, sentence string) {
	komma := ","
	if !buf_has_data[SENT] {
		komma = ""
		buf_has_data[SENT] = true
		fmt.Fprintf(&buffer[SENT], "INSERT `%s_c_%s_sent` (`arch`,`file`,`sent`) VALUES", Cfg.Prefix, prefix)
	}
	fmt.Fprintf(&buffer[SENT], "%s\n(%d,%d,%q)", komma, arch, file, sentence)
	if buffer[SENT].Len() > 49500 {
		buf_flush(SENT)
	}
}

// Zet een dependency relation in de buffer.
// Als de buffer vol raakt, stuur alles naar de database.
func deprel_buf_put(arch, file int, rel string, d *Deprel) {
	komma := ","
	if !buf_has_data[DPRL] {
		komma = ""
		buf_has_data[DPRL] = true
		fmt.Fprintf(&buffer[DPRL], `INSERT %s_c_%s_deprel
		( word,   lemma,   root,   postag, rel,   hword,   hlemma,   hroot,   hpostag, arch, file,   begin,   end,   hbegin,   hend,   mark) VALUES`, Cfg.Prefix, prefix)
	}
	fmt.Fprintf(&buffer[DPRL], "%s\n(%q,%q,%q,%q,%q,%q,%q,%q,%q,%d,%d,%d,%d,%d,%d,%q)", komma,
		d.word, d.lemma, d.root, d.postag, rel, d.hword, d.hlemma, d.hroot, d.hpostag, arch, file, d.begin, d.end, d.hbegin, d.hend, d.mark)
	if buffer[DPRL].Len() > 49500 {
		buf_flush(DPRL)
	}
}

// Zet een woord met lemma's in de buffer.
// Als de buffer vol raakt, stuur alles naar de database.
// Als er maar 1 lemma is, en dat is gelijk aan het woord, negeer het dan.
// Gebruik een simpele collocatie, die de meest-voorkomende varianten pakt.
func word_buf_put(woord string, lemmas []string) {
	if len(lemmas) == 1 && simple(woord) == simple(lemmas[0]) {
		return
	}
	komma := ","
	if !buf_has_data[WORD] {
		komma = ""
		buf_has_data[WORD] = true
		fmt.Fprintf(&buffer[WORD], "INSERT `%s_c_%s_word` (`word`,`lemma`) VALUES", Cfg.Prefix, prefix)
	}
	fmt.Fprintf(&buffer[WORD], "%s\n(%q,%q)", komma, woord, strings.Join(lemmas, "\t"))
	if buffer[WORD].Len() > 49500 {
		buf_flush(WORD)
	}
}

// Zet een archiefnaam in de buffer.
// Als de buffer vol raakt, stuur alles naar de database.
func arch_buf_put(name string, n int) {
	komma := ","
	if !buf_has_data[ARCH] {
		komma = ""
		buf_has_data[ARCH] = true
		fmt.Fprintf(&buffer[ARCH], "INSERT `%s_c_%s_arch` (`id`,`arch`) VALUES", Cfg.Prefix, prefix)
	}
	fmt.Fprintf(&buffer[ARCH], "%s\n(%d,%q)", komma, n, name)
	if buffer[ARCH].Len() > 49500 {
		buf_flush(ARCH)
	}
}

// Zet een xml-filenaam in de buffer.
// Als de buffer vol raakt, stuur alles naar de database.
func file_buf_put(name string, n int) {
	komma := ","
	if !buf_has_data[FILE] {
		komma = ""
		buf_has_data[FILE] = true
		fmt.Fprintf(&buffer[FILE], "INSERT `%s_c_%s_file` (`id`,`file`) VALUES", Cfg.Prefix, prefix)
	}
	fmt.Fprintf(&buffer[FILE], "%s\n(%d,%q)", komma, n, name)
	if buffer[FILE].Len() > 49500 {
		buf_flush(FILE)
	}
}

// Stuur een volle buffer naar de database.
// Reset buffer.
func buf_flush(buf int) {
	if !buf_has_data[buf] {
		return
	}
	buffer[buf].WriteString(";")
	_, err := db.Exec(buffer[buf].String())
	util.CheckErr(err)
	buffer[buf].Reset()
	buf_has_data[buf] = false
}

//. Functies voor sorteren van deprels

func (d Deprels) Len() int { return len(d) }

func (d Deprels) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d Deprels) Less(i, j int) bool {
	if d[i].begin != d[j].begin {
		return d[i].begin < d[j].begin
	}
	if d[i].hbegin != d[j].hbegin {
		return d[i].hbegin < d[j].hbegin
	}
	if d[i].rel != d[j].rel {
		return d[i].rel < d[j].rel
	}
	if d[i].hrel != d[j].hrel {
		return d[i].hrel < d[j].hrel
	}
	return d[i].mark < d[j].mark
}

//. Hulpfuncties

// Afdrukken van grootte van de tabellen.
func sizes() {
	var d, i, total int
	var t string
	rows, err := db.Query("SELECT TABLE_NAME, DATA_LENGTH, INDEX_LENGTH FROM `information_schema`.`TABLES` WHERE TABLE_NAME LIKE '" +
		Cfg.Prefix + "_c_" + prefix + "_%' ORDER BY 1;")
	util.CheckErr(err)
	fmt.Println()
	for rows.Next() {
		util.CheckErr(rows.Scan(&t, &d, &i))
		fmt.Printf("%s: %s + %s = %s\n", t, value(d), value(i), value(d+i))
		total += d
		total += i
	}
	util.CheckErr(rows.Err())
	fmt.Printf("Totaal: %v\n\n", value(total))
}

// Formatteren van een aantal bytes
func value(i int) string {
	if i > 1024*1024*1024*1024*10 {
		return fmt.Sprintf("%.1f Tb", float64(i)/(1024*1024*1024*1024))
	}
	if i > 1024*1024*1024*10 {
		return fmt.Sprintf("%.1f Gb", float64(i)/(1024*1024*1024))
	}
	if i > 1024*1024*10 {
		return fmt.Sprintf("%.1f Mb", float64(i)/(1024*1024))
	}
	if i > 1024*10 {
		return fmt.Sprintf("%.1f Kb", float64(i)/1024)
	}
	return fmt.Sprintf("%d b", i)
}

func connect() *sql.DB {
	db, err := sql.Open("mysql", Cfg.Login+"?charset=utf8mb4,utf8&parseTime=true&loc=Europe%2FAmsterdam")
	util.CheckErr(err)
	return db
}

// Zit string in []string ?
func has(ss []string, s string) bool {
	for _, s1 := range ss {
		if s1 == s {
			return true
		}
	}
	return false
}

//. Simplify spelling

var reChr = regexp.MustCompile("[^ -~]")

var translate = map[string]string{
	"À": "A",
	"Á": "A",
	"Â": "A",
	"Ã": "A",
	"Ä": "A",
	"Å": "A",
	"Ç": "C",
	"È": "E",
	"É": "E",
	"Ê": "E",
	"Ë": "E",
	"Ì": "I",
	"Í": "I",
	"Î": "I",
	"Ï": "I",
	"Ñ": "N",
	"Ò": "O",
	"Ó": "O",
	"Ô": "O",
	"Õ": "O",
	"Ö": "O",
	"Ø": "O",
	"Ù": "U",
	"Ú": "U",
	"Û": "U",
	"Ü": "U",
	"Ý": "Y",
	"à": "a",
	"á": "a",
	"â": "a",
	"ã": "a",
	"ä": "a",
	"å": "a",
	"ç": "c",
	"è": "e",
	"é": "e",
	"ê": "e",
	"ë": "e",
	"ì": "i",
	"í": "i",
	"î": "i",
	"ï": "i",
	"ñ": "n",
	"ò": "o",
	"ó": "o",
	"ô": "o",
	"õ": "o",
	"ö": "o",
	"ø": "o",
	"ù": "u",
	"ú": "u",
	"û": "u",
	"ü": "u",
	"ý": "y",
	"ÿ": "y",
}

func simpleMap(s string) string {
	if s2, ok := translate[s]; ok {
		return s2
	}
	return s
}

func simple(s string) string {
	return strings.ToLower(reChr.ReplaceAllStringFunc(s, simpleMap))
}

////////////////////////////////////////////////////////////////

func tags() {

	var lines int64

	skip := map[string]bool{
		"":     true,
		"--/-": true,
	}

	targets := []string{"postag", "hpostag", "rel"}

	keys := make(map[string]map[string]int64)
	for _, t := range targets {
		keys[t] = make(map[string]int64)
	}

	prefixes := make([]string, 0)
	rows, err := db.Query("SELECT `id` FROM `" + Cfg.Prefix + "_info`;")
	util.CheckErr(err)
	for rows.Next() {
		var val string
		util.CheckErr(rows.Scan(&val))
		prefixes = append(prefixes, val)
	}
	util.CheckErr(rows.Err())

	_, err = db.Exec("DROP TABLE IF EXISTS `" + Cfg.Prefix + "_postag`, `" + Cfg.Prefix + "_hpostag`, `" + Cfg.Prefix + "_rel`;")
	util.CheckErr(err)

	for _, t := range targets {
		_, err = db.Exec(`CREATE TABLE ` + Cfg.Prefix + `_` + t + ` (
			tag varchar(128))
			DEFAULT CHARACTER SET utf8;`)
	}

	for _, p := range prefixes {
		rows, err := db.Query("SELECT count(*) FROM `" + Cfg.Prefix + "_c_" + p + "_sent`;")
		util.CheckErr(err)
		for rows.Next() {
			var val int
			util.CheckErr(rows.Scan(&val))
			lines += int64(val)
		}
		util.CheckErr(rows.Err())

		for _, t := range targets {
			rows, err := db.Query("SELECT count(*),`" + t + "` FROM `" + Cfg.Prefix + "_c_" + p + "_deprel` group by `" + t + "`;")
			util.CheckErr(err)
			var key string
			var val int
			for rows.Next() {
				util.CheckErr(rows.Scan(&val, &key))
				keys[t][key] += int64(val)
			}
			util.CheckErr(rows.Err())
		}
	}
	for _, t := range targets {
		var buf bytes.Buffer
		items := make([]string, 0, len(keys[t]))
		for key := range keys[t] {
			if t != "rel" {
				items = append(items, key)
			} else {
				if keys[t][key]*1000 <= lines {
					items = append(items, "D"+key)
				} else if strings.Index(key, "/-") > 0 {
					items = append(items, "C"+key)
				} else if strings.Index(key, "/") > 0 {
					items = append(items, "B"+key)
				} else {
					items = append(items, "A"+key)
				}
			}
		}
		sort.Strings(items)
		fmt.Fprint(&buf, "INSERT "+Cfg.Prefix+"_"+t+" VALUES (\"\")")
		for _, item := range items {
			if !skip[item] {
				fmt.Fprintf(&buf, ",\n(%q)", item)
			}
		}
		fmt.Fprintln(&buf, ";")
		_, err := db.Exec(buf.String())
		util.CheckErr(err)
	}
}

func showmemstats() {
	fmt.Print("Geheugen: ")
	runtime.ReadMemStats(&memstats)
	v := memstats.Sys
	if v < 1024 {
		fmt.Print(v, " b\n")
	} else if v < 1024*1024 {
		fmt.Printf("%.1f Kb\n", float64(v)/1024.0)
	} else if v < 1024*1024*1024 {
		fmt.Printf("%.1f Mb\n", float64(v)/1024.0/1024.0)
	} else {
		fmt.Printf("%.1f Gb\n", float64(v)/1024.0/1024.0/1024.0)
	}
}

func iformat(i int) string {
	s1 := fmt.Sprint(i)
	s2 := ""
	for n := len(s1); n > 3; n = len(s1) {
		s2 = "." + s1[n-3:n] + s2
		s1 = s1[0 : n-3]
	}
	return s1 + s2
}
