package main

//. Imports

import (
	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pebbe/compactcorpus"
	"github.com/pebbe/util"

	"bufio"
	"bytes"
	"compress/gzip"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

//. Types

type Alpino_ds struct {
	XMLName  xml.Name `xml:"alpino_ds"`
	Meta     []MetaT  `xml:"metadata>meta"`
	Node0    *Node    `xml:"node"`
	Sentence string   `xml:"sentence"`
}

type MetaT struct {
	Type  string `xml:"type,attr"`
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
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

type MetaIdx struct {
	Id   int
	Type string
}

//. Constanten

const (
	DPRL = iota
	SENT
	WORD
	FILE
	ARCH
	META
	MIDX
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
	topmidx  = -1
	lastarch string

	buffer       [7]bytes.Buffer
	buf_has_data [7]bool

	db_append       bool
	db_overwrite    bool
	db_exists       bool
	db_makeindex    bool
	db_updatestatus bool
	db_strippath    bool
	db_decode       bool

	rePath      *regexp.Regexp
	reFilecodes = regexp.MustCompile("_[0-9A-F][0-9A-F]|__")

	memstats runtime.MemStats

	Cfg Config

	utfRE = regexp.MustCompile("&#[0-9]+;|[^\001-\uFFFF]")

	meta = make(map[string]MetaIdx)
)

//. Main

func main() {

	now := time.Now()

	buffer[DPRL].Grow(50000)
	buffer[SENT].Grow(50000)
	buffer[WORD].Grow(50000)
	buffer[FILE].Grow(50000)
	buffer[ARCH].Grow(50000)
	buffer[META].Grow(50000)
	buffer[MIDX].Grow(50000)

	db_makeindex = true
	db_updatestatus = true
	for len(os.Args) > 1 {
		if os.Args[1] == "-a" {
			db_append = true
			db_overwrite = false
		} else if os.Args[1] == "-w" {
			db_append = false
			db_overwrite = true
		} else if os.Args[1] == "-i" {
			db_makeindex = false
		} else if os.Args[1] == "-s" {
			db_updatestatus = false
		} else if os.Args[1] == "-p" && len(os.Args) > 2 {
			db_strippath = true
			os.Args = append(os.Args[:1], os.Args[2:]...)
			rePath = regexp.MustCompile(os.Args[1])
		} else if os.Args[1] == "-d" {
			db_decode = true
		} else {
			break
		}
		os.Args = append(os.Args[:1], os.Args[2:]...)
	}

	if len(os.Args) != 5 || util.IsTerminal(os.Stdin) {
		fmt.Printf(`
Syntax: %s [-a] [-w] [-i] [-s] [-p regexp] [-d] id description owner public < bestandnamen

Opties:

 -a : toevoegen aan bestaande database
 -w : bestaande database overschrijven
 -i : geen tabel van woord naar lemmas aanmaken
 -s : status niet bijwerken als klaar
 -p : prefix die van bestandnaam wordt gestript voor label
 -d : bestandnaam decoderen voor label

  id:
  description:
  owner:       'none' of een e-mailadres
  public:      0 (private) of 1 (public)


`, os.Args[0])
		return
	}

	prefix = strings.TrimSpace(os.Args[1])
	desc = strings.TrimSpace(os.Args[2])
	owner = strings.TrimSpace(os.Args[3])
	public = strings.TrimSpace(os.Args[4])

	paqudir := os.Getenv("PAQU")
	if paqudir == "" {
		paqudir = path.Join(os.Getenv("HOME"), ".paqu")
	}
	_, err := toml.DecodeFile(path.Join(paqudir, "setup.toml"), &Cfg)
	util.CheckErr(err)

	if desc == "" {
		util.CheckErr(fmt.Errorf("De omschrijving mag niet leeg zijn"))
	}

	if owner != "none" && strings.Index(owner, "@") < 0 {
		util.CheckErr(fmt.Errorf("De eigenaar moet 'none' zijn of een e-mailadres"))
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
		fmt.Println("Verbinding met database wordt gesloten...")
		util.CheckErr(db.Close())
	}()

	//
	// kijk of de database al bestaat
	//

	rows, err := db.Query("SELECT `begin` FROM `" + Cfg.Prefix + "_c_" + prefix + "_deprel` LIMIT 0, 1;")
	if err == nil && rows.Next() {
		rows.Close()
		if !(db_append || db_overwrite) {
			util.CheckErr(fmt.Errorf("De database bestaat al, en er is geen optie -a of -w"))
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
			rows, err = db.Query("SELECT MAX(id) FROM " + Cfg.Prefix + "_c_" + prefix + "_midx")
			util.CheckErr(err)
			if rows.Next() {
				if rows.Scan(&topmidx) != nil {
					topmidx = -1
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
				DROP INDEX arch,
				DROP INDEX lbl;`)
			fmt.Println("Verwijderen index uit " + Cfg.Prefix + "_c_" + prefix + "_file ...")
			db.Exec(`ALTER TABLE ` + Cfg.Prefix + "_c_" + prefix + `_file
				DROP INDEX id`)
			fmt.Println("Verwijderen index uit " + Cfg.Prefix + "_c_" + prefix + "_arch ...")
			db.Exec(`ALTER TABLE ` + Cfg.Prefix + "_c_" + prefix + `_arch
				DROP INDEX id;`)
			fmt.Println("Verwijderen indexen uit " + Cfg.Prefix + "_c_" + prefix + "_meta ...")
			db.Exec(`ALTER TABLE ` + Cfg.Prefix + "_c_" + prefix + `_meta
				DROP INDEX id,
				DROP INDEX file,
				DROP INDEX arch,
				DROP INDEX tval,
				DROP INDEX ival,
				DROP INDEX fval,
				DROP INDEX dval,
				DROP INDEX idx;`)
			fmt.Println("Verwijderen indexen uit " + Cfg.Prefix + "_c_" + prefix + "_midx ...")
			db.Exec(`ALTER TABLE ` + Cfg.Prefix + "_c_" + prefix + `_midx
				DROP INDEX id,
				DROP INDEX name;`)
		}
	}

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

	// oude tabellen weggooien
	if !db_exists || db_overwrite {
		_, err := db.Exec(fmt.Sprintf(
			"DROP TABLE IF EXISTS `%s_c_%s_deprel`, `%s_c_%s_sent`, `%s_c_%s_file`, `%s_c_%s_arch`, `%s_c_%s_meta`, `%s_c_%s_midx`;",
			Cfg.Prefix,
			prefix,
			Cfg.Prefix,
			prefix,
			Cfg.Prefix,
			prefix,
			Cfg.Prefix,
			prefix,
			Cfg.Prefix,
			prefix,
			Cfg.Prefix,
			prefix))
		util.CheckErr(err)
		// nieuwe tabellen aanmaken
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
			arch int          NOT NULL,
			file int          NOT NULL,
			sent text         NOT NULL,
			lbl  varchar(260) NOT NULL)
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
		_, err = db.Exec(`CREATE TABLE ` + Cfg.Prefix + "_c_" + prefix + `_midx (
			id   int          NOT NULL,
			type enum('TEXT','INT','FLOAT','DATE','DATETIME') NOT NULL DEFAULT 'TEXT',
			name varchar(128) NOT NULL)
			DEFAULT CHARACTER SET utf8
			DEFAULT COLLATE utf8_unicode_ci;`)
		util.CheckErr(err)
		_, err = db.Exec(`CREATE TABLE ` + Cfg.Prefix + "_c_" + prefix + `_meta (
			id   int          NOT NULL,
			arch int          NOT NULL,
			file int          NOT NULL,
			tval varchar(128) NOT NULL DEFAULT "",
			ival int          NOT NULL DEFAULT 0,
			fval float        NOT NULL DEFAULT 0.0,
			dval datetime     NOT NULL DEFAULT "1000-01-01 00:00:00",
			idx  int          NOT NULL DEFAULT -1)
			DEFAULT CHARACTER SET utf8
			DEFAULT COLLATE utf8_unicode_ci;`)
		util.CheckErr(err)
	}

	//
	// Bestandnamen van stdin inlezen en verwerken.
	//

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
		lowername := strings.ToLower(filename)
		if strings.HasSuffix(lowername, ".xml") {
			data, err := ioutil.ReadFile(filename)
			util.CheckErr(err)
			do_data("", filename, data)
		} else if strings.HasSuffix(lowername, ".xml.gz") {
			fp, err := os.Open(filename)
			util.CheckErr(err)
			r, err := gzip.NewReader(fp)
			util.CheckErr(err)
			data, err := ioutil.ReadAll(r)
			r.Close()
			fp.Close()
			util.CheckErr(err)
			do_data("", filename[:len(filename)-3], data)
		} else if has_dbxml && strings.HasSuffix(lowername, ".dact") {
			do_dact(filename)
		} else if strings.HasSuffix(lowername, ".data.dz") {
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
			util.CheckErr(fmt.Errorf("Ongeldige extensie voor bestand '%s'", filename))
		}
	}
	util.CheckErr(scanner.Err())

	// stuur laatste data uit buffers naar de database
	buf_flush(DPRL)
	buf_flush(SENT)
	buf_flush(FILE)
	buf_flush(ARCH)
	buf_flush(META)
	buf_flush(MIDX)

	_, err = db.Exec("COMMIT;")
	util.CheckErr(err)

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
		ADD INDEX (arch),
		ADD INDEX (lbl);`)
	util.CheckErr(err)

	fmt.Println("Aanmaken index op " + Cfg.Prefix + "_c_" + prefix + "_file ...")
	_, err = db.Exec(`ALTER TABLE ` + Cfg.Prefix + "_c_" + prefix + `_file
		ADD UNIQUE INDEX (id)`)
	util.CheckErr(err)

	fmt.Println("Aanmaken index op " + Cfg.Prefix + "_c_" + prefix + "_arch ...")
	_, err = db.Exec(`ALTER TABLE ` + Cfg.Prefix + "_c_" + prefix + `_arch
		ADD UNIQUE INDEX (id);`)
	util.CheckErr(err)

	fmt.Println("Aanmaken indexen op " + Cfg.Prefix + "_c_" + prefix + "_midx ...")
	_, err = db.Exec(`ALTER TABLE ` + Cfg.Prefix + "_c_" + prefix + `_midx
		ADD INDEX (name),
		ADD UNIQUE INDEX (id);`)
	util.CheckErr(err)

	fmt.Println("Aanmaken indexen op " + Cfg.Prefix + "_c_" + prefix + "_meta ...")
	_, err = db.Exec(`ALTER TABLE ` + Cfg.Prefix + "_c_" + prefix + `_meta
		ADD INDEX (id),
		ADD INDEX (file),
		ADD INDEX (arch),
		ADD INDEX (tval),
		ADD INDEX (ival),
		ADD INDEX (fval),
		ADD INDEX (dval),
		ADD INDEX (idx);`)
	util.CheckErr(err)

	// tijd voor aanmaken tabellen <prefix>_deprel en <prefix>_sent
	fmt.Println("Tijd:", time.Now().Sub(now))

	showmemstats()

	//
	// tabel <prefix>_word aanmaken
	//

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
	_, err = db.Exec("COMMIT;")
	util.CheckErr(err)

	fmt.Println("Aanmaken index op " + Cfg.Prefix + "_c_" + prefix + "_word ...")
	_, err = db.Exec(`ALTER TABLE ` + Cfg.Prefix + "_c_" + prefix + `_word
		ADD UNIQUE INDEX (word);`)
	util.CheckErr(err)

	fmt.Println("Tijd:", time.Now().Sub(now))

	//
	// ranges
	//

	fmt.Println("Ranges bepalen voor " + Cfg.Prefix + "_c_" + prefix + "_meta ...")

	_, err = db.Exec(fmt.Sprintf(
		"DROP TABLE IF EXISTS `%s_c_%s_mval`, %s_c_%s_minf; ",
		Cfg.Prefix,
		prefix,
		Cfg.Prefix,
		prefix))
	util.CheckErr(err)

	_, err = db.Exec(`CREATE TABLE ` + Cfg.Prefix + "_c_" + prefix + `_mval (
			id    int          NOT NULL DEFAULT 0,
			idx   int          NOT NULL DEFAULT 0,
			text  varchar(260) NOT NULL DEFAULT 0,
			n     int          NOT NULL DEFAULT 0)
			DEFAULT CHARACTER SET utf8
			DEFAULT COLLATE utf8_unicode_ci;`)
	util.CheckErr(err)

	_, err = db.Exec(`CREATE TABLE ` + Cfg.Prefix + "_c_" + prefix + `_minf (
			id      int      NOT NULL DEFAULT 0,
			indexed boolean  NOT NULL DEFAULT 1,
			size    int      NOT NULL DEFAULT 0,
			dmin    datetime NOT NULL DEFAULT "1000-01-01 00:00:00",
			dmax    datetime NOT NULL DEFAULT "1000-01-01 00:00:00",
			dtype   int      NOT NULL DEFAULT 0,
			fmin    float    NOT NULL DEFAULT 0.0,
			fstep   float    NOT NULL DEFAULT 0.0,
			imin    int      NOT NULL DEFAULT 0,
			istep   int      NOT NULL DEFAULT 0);`)
	util.CheckErr(err)

	metas := make([]string, 0)
	metat := make(map[string]string)
	metai := make(map[string]int)
	rows, err = db.Query(fmt.Sprintf("SELECT `id`,`name`,`type` FROM `%s_c_%s_midx` ORDER BY 2", Cfg.Prefix, prefix))
	util.CheckErr(err)
	for rows.Next() {
		var i int
		var n, t string
		util.CheckErr(rows.Scan(&i, &n, &t))
		metas = append(metas, n)
		metat[n] = t
		metai[n] = i
	}
	util.CheckErr(rows.Err())
	for _, meta := range metas {
		idx := make(map[int]string)
		switch metat[meta] {
		case "TEXT":
			rows, err := db.Query(fmt.Sprintf(
				"SELECT DISTINCT `tval` FROM `%s_c_%s_meta` WHERE `id` = %d ORDER BY 1",
				Cfg.Prefix, prefix,
				metai[meta]))
			util.CheckErr(err)
			ix := 0
			for rows.Next() {
				var s string
				util.CheckErr(rows.Scan(&s))
				idx[ix] = s
				_, err = db.Exec(fmt.Sprintf(
					"UPDATE `%s_c_%s_meta` SET `idx` = %d WHERE `id` = %d AND `tval` = %q",
					Cfg.Prefix, prefix, ix, metai[meta], s))
				util.CheckErr(err)
				ix++
			}
			util.CheckErr(rows.Err())
		case "INT":
			rows, err := db.Query(fmt.Sprintf(
				"SELECT MIN(`ival`), MAX(`ival`), COUNT(DISTINCT `ival`) FROM `%s_c_%s_meta` WHERE `id` = %d",
				Cfg.Prefix, prefix,
				metai[meta]))
			util.CheckErr(err)
			var v1, v2, vx int
			for rows.Next() {
				rows.Scan(&v1, &v2, &vx)
			}
			ir := newIrange(v1, v2, vx)
			indexed := 0
			if ir.indexed {
				indexed = 1
			}
			_, err = db.Exec(fmt.Sprintf(
				"INSERT `%s_c_%s_minf` (`id`,`imin`,`istep`,`indexed`,`size`) VALUES (%d,%d,%d,%d,%d)",
				Cfg.Prefix, prefix,
				metai[meta], ir.min, ir.step, indexed, len(ir.s)))
			util.CheckErr(err)
			rows, err = db.Query(fmt.Sprintf(
				"SELECT DISTINCT `ival` FROM `%s_c_%s_meta` WHERE `id` = %d",
				Cfg.Prefix, prefix,
				metai[meta]))
			util.CheckErr(err)
			var v int
			for rows.Next() {
				util.CheckErr(rows.Scan(&v))
				s, ix := ir.value(v)
				idx[ix] = s
				_, err = db.Exec(fmt.Sprintf(
					"UPDATE `%s_c_%s_meta` SET `idx` = %d WHERE `id` = %d AND `ival` = %d",
					Cfg.Prefix, prefix,
					ix,
					metai[meta],
					v))
				util.CheckErr(err)
			}
			util.CheckErr(rows.Err())
		case "FLOAT":
			rows, err := db.Query(fmt.Sprintf(
				"SELECT MIN(`fval`), MAX(`fval`) FROM `%s_c_%s_meta` WHERE `id` = %d",
				Cfg.Prefix, prefix,
				metai[meta]))
			util.CheckErr(err)
			var v1, v2 float64
			for rows.Next() {
				rows.Scan(&v1, &v2)
			}
			fr := newFrange(v1, v2)
			_, err = db.Exec(fmt.Sprintf(
				"INSERT `%s_c_%s_minf` (`id`,`fmin`,`fstep`,`size`) VALUES (%d,%g,%g,%d)",
				Cfg.Prefix, prefix,
				metai[meta], fr.min, fr.step, len(fr.s)))
			util.CheckErr(err)
			_, err = db.Exec(fmt.Sprintf(
				"UPDATE `%s_c_%s_meta` SET `idx` = FLOOR((`fval` - %g) / %g) WHERE `id` = %d",
				Cfg.Prefix, prefix,
				fr.min,
				fr.step,
				metai[meta]))
			util.CheckErr(err)
			rows, err = db.Query(fmt.Sprintf(
				"SELECT DISTINCT `idx` FROM `%s_c_%s_meta` WHERE `id` = %d",
				Cfg.Prefix, prefix,
				metai[meta]))
			util.CheckErr(err)
			for rows.Next() {
				var i int
				util.CheckErr(rows.Scan(&i))
				idx[i] = fr.s[i]
			}
			util.CheckErr(err)
		case "DATE", "DATETIME":
			dis := "0"
			if metat[meta] == "DATE" {
				dis = "COUNT(DISTINCT `dval`)"
			}
			rows, err := db.Query(fmt.Sprintf(
				"SELECT MIN(`dval`), MAX(`dval`), %s FROM `%s_c_%s_meta` WHERE `id` = %d",
				dis,
				Cfg.Prefix, prefix,
				metai[meta]))
			util.CheckErr(err)
			var v1, v2 time.Time
			var i int
			for rows.Next() {
				rows.Scan(&v1, &v2, &i)
			}
			dr := newDrange(v1, v2, i, metat[meta] == "DATETIME")
			indexed := 0
			if dr.indexed {
				indexed = 1
			}
			_, err = db.Exec(fmt.Sprintf(
				"INSERT `%s_c_%s_minf` (`id`,`dmin`,`dmax`,`dtype`,`indexed`,`size`) VALUES (%d,\"%04d-%02d-%02d %02d:%02d:%02d\",\"%04d-%02d-%02d %02d:%02d:%02d\",%d,%d,%d)",
				Cfg.Prefix, prefix,
				metai[meta],
				dr.min.Year(), dr.min.Month(), dr.min.Day(), dr.min.Hour(), dr.min.Minute(), dr.min.Second(),
				dr.max.Year(), dr.max.Month(), dr.max.Day(), dr.max.Hour(), dr.max.Minute(), dr.max.Second(),
				dr.r, indexed, len(dr.s)))
			util.CheckErr(err)
			rows, err = db.Query(fmt.Sprintf(
				"SELECT `dval` FROM `%s_c_%s_meta` WHERE `id` = %d",
				Cfg.Prefix, prefix,
				metai[meta]))
			util.CheckErr(err)
			var v time.Time
			for rows.Next() {
				util.CheckErr(rows.Scan(&v))
				s, ix := dr.value(v)
				idx[ix] = s
				_, err = db.Exec(fmt.Sprintf(
					"UPDATE `%s_c_%s_meta` SET `idx` = %d WHERE `id` = %d AND `dval` = \"%04d-%02d-%02d %02d-%02d-%02d\"",
					Cfg.Prefix, prefix,
					ix,
					metai[meta],
					v.Year(), v.Month(), v.Day(),
					v.Hour(), v.Minute(), v.Second()))
				util.CheckErr(err)
			}
			util.CheckErr(rows.Err())
		}
		_, err = db.Exec("COMMIT;")
		util.CheckErr(err)

		// zinnen waarvoor geen metadata is, die toevoegen
		_, err = db.Exec(fmt.Sprintf(
			"INSERT `%s_c_%s_meta` (`id`,`arch`,`file`,`idx`)"+
				"SELECT DISTINCT %d, `arch`, `file`, 2147483647 FROM `%s_c_%s_sent` `s` WHERE NOT EXISTS ( "+
				"SELECT `arch`, `file` FROM `%s_c_%s_meta` `m` WHERE `s`.`arch`=`m`.`arch` AND `s`.`file`=`m`.`file` AND `id`=%d )",
			Cfg.Prefix, prefix,
			metai[meta], Cfg.Prefix, prefix,
			Cfg.Prefix, prefix, metai[meta]))
		util.CheckErr(err)
		_, err = db.Exec("COMMIT;")
		util.CheckErr(err)
		// kijk of er echt metadata is toegevoegd
		rows, err = db.Query(fmt.Sprintf(
			"SELECT DISTINCT 1 FROM `%s_c_%s_meta` WHERE `id`=%d AND `idx`=2147483647",
			Cfg.Prefix, prefix,
			metai[meta]))
		util.CheckErr(err)
		for rows.Next() {
			idx[2147483647] = ""
		}
		util.CheckErr(rows.Err())

		for ix := range idx {
			_, err = db.Exec(fmt.Sprintf(
				"INSERT `%s_c_%s_mval` (`id`,`idx`,`text`) VALUES (%d,%d,%q)",
				Cfg.Prefix, prefix,
				metai[meta],
				ix,
				idx[ix]))
			util.CheckErr(err)
		}
		_, err = db.Exec("COMMIT;")
		util.CheckErr(err)
	}

	fmt.Println("Aanmaken indexen op " + Cfg.Prefix + "_c_" + prefix + "_mval ...")
	_, err = db.Exec(`ALTER TABLE ` + Cfg.Prefix + "_c_" + prefix + `_mval
		ADD INDEX (id),
		ADD INDEX (idx);`)
	util.CheckErr(err)

	fmt.Println("Aanmaken index op " + Cfg.Prefix + "_c_" + prefix + "_minf ...")
	_, err = db.Exec(`ALTER TABLE ` + Cfg.Prefix + "_c_" + prefix + `_minf
		ADD INDEX (id);`)
	util.CheckErr(err)

	fmt.Println("Telling van ranges ...")
	for _, meta := range metas {
		sums := make(map[int]int)
		rows, err := db.Query(fmt.Sprintf(
			"SELECT COUNT(`idx`),`idx` FROM `%s_c_%s_meta` WHERE `id` = %d GROUP BY `idx`",
			Cfg.Prefix, prefix,
			metai[meta]))
		util.CheckErr(err)
		for rows.Next() {
			var c, i int
			util.CheckErr(rows.Scan(&c, &i))
			sums[i] = c
		}
		util.CheckErr(rows.Err())

		for s := range sums {
			_, err = db.Exec(fmt.Sprintf(
				"UPDATE `%s_c_%s_mval` SET `n` = %d WHERE `id` = %d AND `idx` = %d",
				Cfg.Prefix, prefix,
				sums[s],
				metai[meta], s))
			util.CheckErr(err)
		}
		_, err = db.Exec("COMMIT;")
		util.CheckErr(err)
	}

	//
	// zet info over corpus in de database
	//

	lines := 0
	rows, err = db.Query("SELECT COUNT(*) FROM " + Cfg.Prefix + "_c_" + prefix + "_sent")
	util.CheckErr(err)
	if rows.Next() {
		util.CheckErr(rows.Scan(&lines))
		rows.Close()
	}
	hasmeta := 0
	if len(metas) > 0 {
		hasmeta = 1
	} else {
		db.Exec(fmt.Sprintf(
			"DROP TABLE IF EXISTS `%s_c_%s_meta`, `%s_c_%s_midx`, `%s_c_%s_minf`, `%s_c_%s_mval`;",
			Cfg.Prefix,
			prefix,
			Cfg.Prefix,
			prefix,
			Cfg.Prefix,
			prefix,
			Cfg.Prefix,
			prefix))
	}
	if db_updatestatus {
		_, err = db.Exec(fmt.Sprintf("UPDATE `%s_info` SET `status` = \"FINISHED\", `nline` = %d, `active` = NOW(), `hasmeta` = %d WHERE `id` = %q",
			Cfg.Prefix, lines, hasmeta, prefix))
	} else {
		_, err = db.Exec(fmt.Sprintf("UPDATE `%s_info` SET `nline` = %d, `active` = NOW(), `hasmeta` = %d WHERE `id` = %q",
			Cfg.Prefix, lines, hasmeta, prefix))
	}
	util.CheckErr(err)

	user := owner
	if public == "1" {
		user = "all"
	}
	_, err = db.Exec(fmt.Sprintf("INSERT `%s_corpora` (`user`, `prefix`) VALUES (%q, %q);", Cfg.Prefix, user, prefix))
	util.CheckErr(err)

	_, err = db.Exec("COMMIT;")
	util.CheckErr(err)

	//fmt.Println("Bijwerken menu's voor postag, rel en hpostag ...")
	//tags()

	// totale tijd
	fmt.Println("Tijd:", time.Now().Sub(now))
	showmemstats()

	sizes()
}

//. Verwerking

func utfFunc(b []byte) []byte {
	if b[0] == '&' {
		i, err := strconv.Atoi(string(b[2 : len(b)-1]))
		if err != nil || (err == nil && i < 65536) {
			return b
		}
		bb := []byte("&amp;")
		bb = append(bb, b[1:]...)
		return bb
	}

	u, _ := utf8.DecodeRune(b)
	return []byte(fmt.Sprintf("&amp;#%d;", u))
}

// Verwerk een enkel xml-bestand
func do_data(archname, filename string, data []byte) {

	// MySQL tot versie 5.5.3 kan niet met tekens boven U+FFFF overweg
	data = utfRE.ReplaceAllFunc(data, utfFunc)

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

	for _, m := range alpino.Meta {
		if m.Type != "text" && m.Type != "int" && m.Type != "float" && m.Type != "date" && m.Type != "datetime" {
			util.CheckErr(fmt.Errorf("Ongeldig type in %s||%s: %s", archname, filename, m.Type))
		}
		if _, ok := meta[m.Name]; !ok {
			topmidx++
			meta[m.Name] = MetaIdx{Id: topmidx, Type: m.Type}
			midx_buf_put(topmidx, m.Name, m.Type)
		}
		mi := meta[m.Name]
		if m.Type != mi.Type {
			util.CheckErr(fmt.Errorf("Ongeldig type in %s||%s: %s, eerder gedefinieerd als %s", archname, filename, m.Type, mi.Type))
		}
		var txt string
		var intval int
		var floatval float64
		dateval := "1000-01-01 00:00:00"
		if m.Type == "int" {
			intval, err = strconv.Atoi(m.Value)
			// 2147483647 is gereserveerd voor intern gebruik
			if intval < -2147483648 || intval > 2147483646 {
				util.CheckErr(fmt.Errorf("Integer niet in bereik -2147483648 - 2147483646: %d", intval))
			}
			util.CheckErr(err)
		} else if m.Type == "float" {
			floatval, err = strconv.ParseFloat(m.Value, 32) // 32 is dezelfde precisie als gebruikt door MySQL
			if floatval < -math.MaxFloat32 || floatval > math.MaxFloat32 {
				util.CheckErr(fmt.Errorf("Float niet in bereik %g - %g: %g", -math.MaxFloat32, math.MaxFloat32, floatval))
			}
			if floatval > -math.SmallestNonzeroFloat32 && floatval < math.SmallestNonzeroFloat32 && floatval != 0 {
				util.CheckErr(fmt.Errorf("Float te klein: %g", floatval))
			}
			util.CheckErr(err)
		} else if m.Type == "date" {
			t, err := time.Parse("2006-01-02", m.Value)
			util.CheckErr(err)
			year := t.Year()
			if year < 1000 || year > 9999 {
				util.CheckErr(fmt.Errorf("Jaartal niet in bereik 1000 - 9999: %d", year))
			}
			dateval = fmt.Sprintf("%04d-%02d-%02d 00:00:00", year, t.Month(), t.Day())
		} else if m.Type == "datetime" {
			t, err := time.Parse("2006-01-02 15:04", m.Value)
			if err != nil {
				t, err = time.Parse("2006-01-02 15:04:05", m.Value)
			}
			util.CheckErr(err)
			year := t.Year()
			if year < 1000 || year > 9999 {
				util.CheckErr(fmt.Errorf("Jaartal niet in bereik 1000-9999: %d", year))
			}
			dateval = fmt.Sprintf("%04d-%02d-%02d %02d:%02d:00", year, t.Month(), t.Day(), t.Hour(), t.Minute())
		} else {
			txt = m.Value
		}
		meta_buf_put(mi.Id, arch, file, txt, intval, floatval, dateval)
	}

	// zin opslaan in <prefix>_sent
	sent_buf_put(arch, file, alpino.Sentence, archname, filename)

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
	// negeer woorden met relatie == "--" en pt == "let"
	if node.Word != "" && !(node.Rel == "--" && node.Postag == "let") && !node.used {
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
		// negeer woorden met relatie == "--" en pt == "let"
		if node.Rel == "--" && node.Postag == "let" {
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

func repl_filecode(s string) string {
	if s == "__" {
		return "_"
	}
	i, _ := strconv.ParseInt(s[1:], 16, 0)
	b := []byte{byte(i)}
	return string(b)
}

// Zet een zin in de buffer.
// Als de buffer vol raakt, stuur alles naar de database.
func sent_buf_put(arch, file int, sentence, archname, filename string) {

	lbl := noext(filename, ".xml")
	if archname != "" {
		lbl = noext(archname, ".dact", "data.dz") + "::" + lbl
	}
	if db_strippath {
		lbl = rePath.ReplaceAllLiteralString(lbl, "")
	}
	if db_decode {
		lbl = reFilecodes.ReplaceAllStringFunc(lbl, repl_filecode)
	}

	komma := ","
	if !buf_has_data[SENT] {
		komma = ""
		buf_has_data[SENT] = true
		fmt.Fprintf(&buffer[SENT], "INSERT `%s_c_%s_sent` (`arch`,`file`,`sent`,`lbl`) VALUES", Cfg.Prefix, prefix)
	}
	fmt.Fprintf(&buffer[SENT], "%s\n(%d,%d,%q,%q)", komma, arch, file, sentence, lbl)
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

// Zet een meta-item in de buffer.
// Als de buffer vol raakt, stuur alles naar de database.
func meta_buf_put(id int, arch int, file int, txt string, intval int, floatval float64, dateval string) {
	komma := ","
	if !buf_has_data[META] {
		komma = ""
		buf_has_data[META] = true
		fmt.Fprintf(&buffer[META], "INSERT `%s_c_%s_meta` (`id`,`arch`,`file`,`tval`,`ival`,`fval`,`dval`) VALUES", Cfg.Prefix, prefix)
	}
	fmt.Fprintf(&buffer[META], "%s\n(%d,%d,%d,%q,%d,%g,%q)", komma, id, arch, file, txt, intval, float32(floatval), dateval)
	if buffer[META].Len() > 49500 {
		buf_flush(META)
	}
}

// Zet een meta-indexitem in de buffer.
// Als de buffer vol raakt, stuur alles naar de database.
func midx_buf_put(id int, name string, mtype string) {
	komma := ","
	if !buf_has_data[MIDX] {
		komma = ""
		buf_has_data[MIDX] = true
		fmt.Fprintf(&buffer[MIDX], "INSERT `%s_c_%s_midx` (`id`,`type`,`name`) VALUES", Cfg.Prefix, prefix)
	}
	fmt.Fprintf(&buffer[MIDX], "%s\n(%d,%q,%q)", komma, id, strings.ToUpper(mtype), name)
	if buffer[MIDX].Len() > 49500 {
		buf_flush(MIDX)
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
		Cfg.Prefix + "\\_c\\_" + prefix + "\\_%' ORDER BY 1;")
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
	db, err := sql.Open("mysql", Cfg.Login+"?charset=utf8&parseTime=true&loc=Europe%2FAmsterdam")
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

func noext(name string, ext ...string) string {
	s := strings.ToLower(name)
	for _, e := range ext {
		if strings.HasSuffix(s, e) {
			return name[:len(s)-len(e)]
		}
	}
	return name
}
