package main

//. Imports

import (
	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pebbe/util"

	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

//. Types

type Config struct {
	Login  string
	Prefix string
	Conllu bool

	Sh   string
	Path string
}

var (
	Cfg Config

	x = util.CheckErr
)

//. Main

func main() {
	configfile := filepath.Join(paquconfigdir, "setup.toml")

	var re *regexp.Regexp
	var err error
	if len(os.Args) > 1 {
		re, err = regexp.Compile(os.Args[1])
		x(err)
	} else {
		fmt.Println("Met regexp: alleen corpora die matchen")
	}

	_, err = TomlDecodeFile(configfile, &Cfg)
	x(err)

	if !Cfg.Conllu {
		x(fmt.Errorf("Option 'conllu' in '%s' is false", configfile))
	}

	x(os.Chdir(filepath.Join(paqudatadir, "data")))

	corporall := make([]string, 0)
	db, err := dbopen()
	x(err)
	rows, err := db.Query(
		fmt.Sprintf(
			"SELECT `id` FROM `%s_info` WHERE `owner` != \"none\" AND `status` = \"FINISHED\" ORDER BY `id`",
			Cfg.Prefix))
	x(err)
	for rows.Next() {
		var id string
		x(rows.Scan(&id))
		corporall = append(corporall, id)
	}
	x(rows.Err())
	db.Close()

	b, err := exec.Command("pqudep", "-v").Output()
	x(err)
	version := strings.TrimSpace(string(b))

	corpora := make([]string, 0)
	for _, corpus := range corporall {
		if re != nil && !re.MatchString(corpus) {
			continue
		}
		up_to_date := false
		for {
			b, err := ioutil.ReadFile(corpus + "/conllu.version")
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

	for teller, corpus := range corpora {
		fmt.Printf("Updating: [%d/%d] %s\n", teller+1, len(corpora), corpus)

		cmd := exec.Command("sh", "-c", fmt.Sprintf("find %s -name '*.xml.gz' | sort | pqudep -o", corpus))
		cmd.Stdout = os.Stdout
		out, err := cmd.StderrPipe()
		x(err)
		x(cmd.Start())
		scanner := bufio.NewScanner(out)
		fp, err := os.Create(corpus + "/conllu.err.tmp")
		x(err)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasSuffix(line, ".xml.gz") {
				line = line[:len(line)-3]
			}
			fmt.Fprintln(fp, line)
		}
		fp.Close()
		x(scanner.Err())
		x(cmd.Wait())
		x(os.Rename(corpus+"/conllu.err.tmp", corpus+"/conllu.err"))

		cmd = exec.Command("sh", "-c", fmt.Sprintf("pqudep -o %s/data.dact", corpus))
		cmd.Stdout = os.Stdout
		x(cmd.Run())

		if _, err = os.Stat(corpus + "/data.dactx"); err == nil {
			cmd = exec.Command("sh", "-c", fmt.Sprintf("pqdactx %s/data.dact %s/data.dactx.tmp", corpus, corpus))
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			x(cmd.Run())
			x(os.Rename(corpus+"/data.dactx.tmp", corpus+"/data.dactx"))
		}

		fp, err = os.Create(corpus + "/conllu.version")
		x(err)
		fmt.Fprintln(fp, version)
		fp.Close()
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
