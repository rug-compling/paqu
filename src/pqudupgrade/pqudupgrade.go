package main

//. Imports

import (
	"github.com/rug-compling/paqu/internal/dir"

	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pebbe/dbxml"
	"github.com/pebbe/util"

	"bufio"
	"bytes"
	"compress/gzip"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

//. Types

type Config struct {
	Login  string
	Prefix string
	Conllu bool
	Dactx  bool

	Sh   string
	Path string
}

type Status struct {
	Xml      string
	Dact     string
	filename string
}

type Corpus struct {
	id     string
	params string
}

var (
	Cfg         Config
	reFilecodes = regexp.MustCompile("_[0-9A-F][0-9A-F]|__")
	x           = util.CheckErr
)

//. Main

func main() {

	if len(os.Args) != 2 {
		fmt.Printf(`
Usage: %s regexp

`, os.Args[0])
		return
	}

	re, err := regexp.Compile(os.Args[1])
	x(err)

	configfile := filepath.Join(dir.Config, "setup.toml")
	_, err = TomlDecodeFile(configfile, &Cfg)
	x(err)

	if !Cfg.Conllu {
		x(fmt.Errorf("Option 'conllu' in '%s' is false", configfile))
	}

	x(os.Chdir(filepath.Join(dir.Data, "data")))

	corporall := make([]Corpus, 0)
	db, err := dbopen()
	x(err)
	rows, err := db.Query(
		fmt.Sprintf(
			"SELECT `id`,`params` FROM `%s_info` WHERE `owner` != \"none\" AND `owner` != \"auto\" AND `owner` != \"manual\" AND `status` = \"FINISHED\" ORDER BY `id`",
			Cfg.Prefix))
	x(err)
	for rows.Next() {
		var id, params string
		x(rows.Scan(&id, &params))
		corporall = append(corporall, Corpus{id: id, params: params})
	}
	x(rows.Err())
	db.Close()

	b, err := exec.Command("pqudep", "-v").Output()
	x(err)
	version := strings.TrimSpace(string(b))

	corpora := make([]Corpus, 0)
	for _, corpus := range corporall {
		if !re.MatchString(corpus.id) {
			continue
		}
		up_to_date := false
		for {
			b, err := ioutil.ReadFile(corpus.id + "/conllu.version")
			if err != nil {
				break
			}
			s := strings.TrimSpace(string(b))
			if s != version {
				break
			}
			up_to_date = true
			break
		}
		if !up_to_date {
			corpora = append(corpora, corpus)
		}
	}

	var lockfile string

	go func() {
		chSignal := make(chan os.Signal, 1)
		signal.Notify(chSignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		<-chSignal
		os.Remove(lockfile)
		os.Exit(1)
	}()

	for teller, corpus := range corpora {

		lockfile = corpus.id + "/lock"
		if os.Symlink(fmt.Sprintf("%s.%d", lockfile, os.Getpid()), lockfile) != nil {
			fmt.Printf("Locked:     [%d/%d] %s\n", teller+1, len(corpora), corpus.id)
			continue
		}

		up_to_date := false
		for {
			b, err := ioutil.ReadFile(corpus.id + "/conllu.version")
			if err != nil {
				break
			}
			s := strings.TrimSpace(string(b))
			if s != version {
				break
			}
			up_to_date = true
			break
		}
		if up_to_date {
			fmt.Printf("Up to date: [%d/%d] %s\n", teller+1, len(corpora), corpus.id)
			os.Remove(lockfile)
			continue
		}

		fmt.Printf("Updating:   [%d/%d] %s\n", teller+1, len(corpora), corpus.id)

		status := readStatus(corpus.id)

		if status.Xml != version {
			cmd := exec.Command("sh", "-c", fmt.Sprintf("find %s -name '*.xml.gz' | sort | pqudep -o", corpus.id))
			cmd.Stdout = os.Stdout
			out, err := cmd.StderrPipe()
			x(err)
			x(cmd.Start())
			scanner := bufio.NewScanner(out)
			fp, err := os.Create(corpus.id + "/conllu.err.tmp")
			x(err)
			for scanner.Scan() {
				line := strings.TrimRight(scanner.Text(), "\n")
				if strings.HasSuffix(line, ".xml.gz") {
					line = line[:len(line)-3]
				}
				fmt.Fprintln(fp, line)
			}
			fp.Close()
			x(scanner.Err())
			x(cmd.Wait())
			x(os.Rename(corpus.id+"/conllu.err.tmp", corpus.id+"/conllu.err"))

			status.Xml = version
			saveStatus(status)
		}

		if status.Dact != version {

			p := ""
			if strings.Contains(corpus.params, "-lbl") || strings.HasPrefix(corpus.params, "folia") || strings.HasPrefix(corpus.params, "tei") {
				p = "-"
			} else if strings.HasPrefix(corpus.params, "xmlzip") || strings.HasPrefix(corpus.params, "dact") {
				p = "/"
			}
			makeDact(corpus.id+"/data.dact.tmp", corpus.id+"/xml", p)
			os.Rename(corpus.id+"/data.dact.tmp", corpus.id+"/data.dact")

			status.Dact = version
			saveStatus(status)
		}

		if Cfg.Dactx {
			os.Remove(corpus.id + "/data.dactx.tmp")
			cmd := exec.Command("sh", "-c", fmt.Sprintf("pqdactx %s/data.dact %s/data.dactx.tmp", corpus.id, corpus.id))
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			x(cmd.Run())
			x(os.Rename(corpus.id+"/data.dactx.tmp", corpus.id+"/data.dactx"))
		}

		removeStatus(status)

		fp, err := os.Create(corpus.id + "/conllu.version")
		x(err)
		fmt.Fprintln(fp, version)
		fp.Close()

		os.Remove(lockfile)
	}
}

func dbopen() (*sql.DB, error) {
	if Cfg.Login[0] == '$' {
		Cfg.Login = os.Getenv(Cfg.Login[1:])
	}
	return sql.Open("mysql", Cfg.Login+"?charset=utf8&parseTime=true&loc=Europe%2FAmsterdam&sql_mode=''")
}

func TomlDecodeFile(fpath string, v interface{}) (toml.MetaData, error) {
	bs, err := ioutil.ReadFile(fpath)
	if err != nil {
		return toml.MetaData{}, err
	}
	// skip BOM (berucht op Windows)
	if bytes.HasPrefix(bs, []byte{239, 187, 191}) {
		bs = bs[3:]
	}
	return toml.Decode(string(bs), v)
}

func readStatus(corpus string) Status {

	status := Status{
		filename: corpus + "/pqudupgrade.json",
	}

	b, err := ioutil.ReadFile(status.filename)
	if err == nil {
		x(json.Unmarshal(b, &status))
	}

	return status
}

func saveStatus(status Status) {
	b, err := json.Marshal(status)
	x(err)
	fp, err := os.Create(status.filename)
	fp.Write(b)
	fp.Close()
}

func removeStatus(status Status) {
	os.Remove(status.filename)
}

func makeDact(dact, xml string, stripchar string) {
	files := filenames2(xml)

	os.Remove(dact)
	db, err := dbxml.OpenReadWrite(dact)
	x(err)
	defer db.Close()

	for index, gzname := range files {

		name := gzname[:len(gzname)-3]
		fp, err := os.Open(filepath.Join(xml, gzname))
		x(err)
		rd, err := gzip.NewReader(fp)
		x(err)
		data, err := ioutil.ReadAll(rd)
		x(err)
		rd.Close()
		fp.Close()

		name = decode_filename(name)
		if stripchar != "" {
			name = name[1+strings.Index(name, stripchar):]
		}
		fmt.Printf("\r\033[Kdata.dact -> [%d/%d] %s ", index, len(files), name)
		x(db.PutXml(name, strings.TrimSpace(string(data)), false))

	}
	fmt.Print("\r\033[K")
}

// alle bestandsnamen van all subdirectories van de gegeven directory)
func filenames2(dirname string) []string {
	fnames := make([]string, 0)
	dirs, err := ioutil.ReadDir(dirname)
	x(err)

	for _, dir := range dirs {
		dname := dir.Name()
		files, err := ioutil.ReadDir(filepath.Join(dirname, dname))
		x(err)
		for _, file := range files {
			if name := file.Name(); strings.HasSuffix(name, ".xml.gz") {
				fnames = append(fnames, filepath.Join(dname, name))
			}
		}
	}
	return fnames
}

func repl_filecode(s string) string {
	if s == "__" {
		return "_"
	}
	i, _ := strconv.ParseInt(s[1:], 16, 0)
	b := []byte{byte(i)}
	return string(b)
}

func decode_filename(s string) string {
	if s == "_" {
		return ""
	}

	return reFilecodes.ReplaceAllStringFunc(s, repl_filecode)
}
