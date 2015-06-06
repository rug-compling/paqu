package main

//. Imports

import (
	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pebbe/util"

	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
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
	paqudir := os.Getenv("PAQU")
	if paqudir == "" {
		paqudir = path.Join(os.Getenv("HOME"), ".paqu")
	}
	_, err := toml.DecodeFile(path.Join(paqudir, "setup.toml"), &Cfg)
	util.CheckErr(err)

	util.CheckErr(os.Chdir(path.Join(paqudir, "data")))

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
			"SELECT `id`, `owner`, `status`, `nword`, `active` FROM `%s_info` WHERE `owner` != \"none\" ORDER BY `active` DESC",
			Cfg.Prefix))
	util.CheckErr(err)

	corpora := make(map[string]int)
	words := make(map[string]int)

	fmt.Print("\nCORPORA\n\nlaatst gebruikt\tid\t\t\ttokens\t\tstatus\t\teigenaar\n")
	for rows.Next() {
		var id, owner, status string
		var nword int
		var active time.Time
		util.CheckErr(rows.Scan(&id, &owner, &status, &nword, &active))
		fmt.Printf("%v\t%-23s\t%-15s\t%-10s\t%v\n",
			date(active),
			id,
			fmt.Sprint(nword),
			status,
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

	dirnames := make([]string, 0)
	for id, val := range disk {
		if val {
			dirnames = append(dirnames, id)
		}
	}
	if len(dirnames) == 0 {
		return
	}
	sort.Strings(dirnames)
	fmt.Print("\nVERWEESDE DIRECTORY'S\n\n")
	for _, d := range dirnames {
		fmt.Println(d)
	}
	fmt.Println()

}

func date(t time.Time) string {
	return fmt.Sprintf("%02d-%02d-%d", t.Day(), t.Month(), t.Year())
}

func dbopen() (*sql.DB, error) {
	return sql.Open("mysql", Cfg.Login+"?charset=utf8&parseTime=true&loc=Europe%2FAmsterdam")
}
