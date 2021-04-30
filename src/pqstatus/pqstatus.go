package main

//. Imports

import (
	"github.com/rug-compling/paqu/internal/dir"

	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pebbe/util"

	"bytes"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

//. Types

type Config struct {
	Login  string
	Prefix string

	Sh   string
	Path string
}

var (
	Cfg Config
)

//. Main

func main() {
	_, err := TomlDecodeFile(filepath.Join(dir.Config, "setup.toml"), &Cfg)
	util.CheckErr(err)

	util.CheckErr(os.Chdir(filepath.Join(dir.Data, "data")))

	disk := make(map[string]bool)
	files, err := ioutil.ReadDir(".")
	util.CheckErr(err)
	for _, i := range files {
		disk[i.Name()] = true
	}

	db, err := dbopen()
	util.CheckErr(err)
	defer db.Close()

	rows, err := db.Query(
		fmt.Sprintf(
			"SELECT `id`, `owner`, `status`, `nword`, `active` FROM `%s_info` WHERE `owner` != \"none\" AND `owner` != \"auto\" AND `owner` != \"manual\" ORDER BY `active` DESC",
			Cfg.Prefix))
	util.CheckErr(err)

	corpora := make(map[string]int)
	words := make(map[string]int)

	fmt.Print("\nCORPORA\n\nlaatst gebruikt\tid\t\t\ttokens\t\tstatus\t\tconllu\teigenaar\n")
	for rows.Next() {
		var id, owner, status string
		var nword int
		var active time.Time
		util.CheckErr(rows.Scan(&id, &owner, &status, &nword, &active))
		conllu := ""
		b, err := ioutil.ReadFile(id + "/conllu.version")
		if err == nil {
			conllu = strings.TrimSpace(string(b))
		}
		fmt.Printf("%v\t%-23s\t%-15s\t%-10s\t%s\t%v\n",
			date(active),
			id,
			fmt.Sprint(nword),
			status,
			conllu,
			owner)
		disk[id] = false
		corpora[owner]++
		words[owner] += nword
	}

	rows, err = db.Query(
		fmt.Sprintf(
			"SELECT `mail`, `active`, `quotum` FROM `%s_users` ORDER BY `active` DESC",
			Cfg.Prefix))
	util.CheckErr(err)
	fmt.Print("\n\nGEBRUIKERS\n\nlaatst actief\tquotum\t\tcorpora\ttokens\t\tmail\n")
	for rows.Next() {
		var mail string
		var quotum int
		var active time.Time
		util.CheckErr(rows.Scan(&mail, &active, &quotum))
		fmt.Printf("%v\t%-15s\t%v\t%-15s\t%v\n",
			date(active),
			fmt.Sprint(quotum),
			corpora[mail],
			fmt.Sprint(words[mail]),
			mail)
	}

	fmt.Println()

	rows, err = db.Query(
		fmt.Sprintf(
			"SELECT `id` FROM `%s_info` WHERE `owner` = \"none\" OR `owner` = \"auto\" OR `owner` = \"manual\"",
			Cfg.Prefix))
	util.CheckErr(err)
	globals := make(map[string]bool)
	for rows.Next() {
		var id string
		util.CheckErr(rows.Scan(&id))
		globals[id] = true
	}
	dirnames := make([]string, 0)
	for id, val := range disk {
		if val && !globals[id] {
			dirnames = append(dirnames, id)
		}
	}
	if len(dirnames) != 0 {
		sort.Strings(dirnames)
		fmt.Print("\nVERWEESDE DIRECTORY'S\n\n")
		for _, d := range dirnames {
			fmt.Println(d)
		}
		fmt.Println()
	}

	if os.Chdir(filepath.Join(dir.Data, "folia")) != nil {
		return
	}
	dirnames = dirnames[0:0]
	files, err = ioutil.ReadDir(".")
	util.CheckErr(err)
	for _, i := range files {
		if i.IsDir() {
			dirnames = append(dirnames, i.Name())
		}
	}
	if len(dirnames) == 0 {
		return
	}
	fmt.Print("\nFOLIADATA IN KB\n\n")
	sort.Strings(dirnames)
	for _, name := range dirnames {
		b, err := hex.DecodeString(name)
		util.CheckErr(err)
		fmt.Printf("%.1f\t%s\n", float64(folia(name))/1024, string(b))
	}
	fmt.Println()
}

func date(t time.Time) string {
	return fmt.Sprintf("%02d-%02d-%d", t.Day(), t.Month(), t.Year())
}

func dbopen() (*sql.DB, error) {
	if Cfg.Login[0] == '$' {
		Cfg.Login = os.Getenv(Cfg.Login[1:])
	}
	return sql.Open("mysql", Cfg.Login+"?charset=utf8&parseTime=true&loc=Europe%2FAmsterdam&sql_mode=''")
}

func folia(dir string) int64 {
	var sum int64
	files, err := ioutil.ReadDir(dir)
	util.CheckErr(err)
	for _, file := range files {
		if file.IsDir() {
			sum += folia(filepath.Join(dir, file.Name()))
		} else {
			sum += file.Size()
		}
	}
	return sum
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
