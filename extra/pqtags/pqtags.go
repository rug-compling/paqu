/*

Dit programma is gebruikt om de code te genereren voor
de menu's voor postags en rels, momenteel deel van
../src/pqserve/globals.go

Gebruik: pqtags corpus...

Voorbeeld: pqtags lassysmall lassywiki

Hierin is lassysmall het corpus Lassy Klein, en
lassywiki het Wikipedia-gedeelte van Lassy Groot.
De laatste is volledig automatisch geparst en bevat
daarom, door fouten van de parser, combinaties van
relaties die je niet in de eerste ziet. Dit zijn meest
"zeldzame" relaties, en die wil je in het menu hebben om
fouten in de parse te kunnen opsporen.

*/

package main

import (
	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pebbe/util"

	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Config struct {
	Login  string
	Prefix string
}

var (
	paqudir string
	Cfg     Config
	x       = util.CheckErr
)

func main() {

	paqudir = os.Getenv("PAQU")
	if paqudir == "" {
		paqudir = filepath.Join(os.Getenv("HOME"), ".paqu")
	}
	tom := filepath.Join(paqudir, "setup.toml")
	_, err := toml.DecodeFile(tom, &Cfg)
	x(err)

	var lines int64

	db, err := dbopen()
	x(err)

	/*
		"--/-" werd oorspronkelijk weggelaten, maar vanwege
		redenen later toch toegevoegd aan het menu
	*/
	skip := map[string]bool{
		"": true,
		// "--/-": true,
	}

	targets := []string{"postag", "hpostag", "rel"}

	keys := make(map[string]map[string]int64)
	for _, t := range targets {
		keys[t] = make(map[string]int64)
	}
	for _, p := range os.Args[1:] {
		rows, err := db.Query("SELECT count(*) FROM `" + Cfg.Prefix + "_c_" + p + "_sent`;")
		x(err)
		if rows.Next() {
			var val int
			x(rows.Scan(&val))
			lines += int64(val)
		}
		x(rows.Err())

		for _, t := range targets {
			rows, err := db.Query("SELECT count(*),`" + t + "` FROM `" + Cfg.Prefix + "_c_" + p + "_deprel` group by `" + t + "`;")
			x(err)
			var key string
			var val int
			for rows.Next() {
				x(rows.Scan(&val, &key))
				keys[t][key] += int64(val)
			}
			x(rows.Err())
		}
	}
	for _, t := range targets {
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
		fmt.Print("\topt_", t, " = []string{\"\"")
		if t != "rel" {
			fmt.Print(", \"(leeg)\"")
		}
		for _, item := range items {
			if !skip[item] {
				fmt.Print(", \"", item, "\"")
			}
		}
		fmt.Println("}")
	}
}

func dbopen() (*sql.DB, error) {
	login := Cfg.Login
	if login[0] == '$' {
		login = os.Getenv(login[1:])
	}
	return sql.Open("mysql", login+"?charset=utf8&parseTime=true&loc=Europe%2FAmsterdam&sql_mode=''")
}
