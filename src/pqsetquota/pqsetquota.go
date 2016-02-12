package main

//. Imports

import (
	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pebbe/util"

	"database/sql"
	"fmt"
	"os"
	"path"
	"strconv"
)

//. Types

type Config struct {
	Login  string
	Prefix string
}

var (
	DefaultPaquDir string
)

func main() {

	if len(os.Args) < 3 {
		fmt.Printf("\nSyntax: %s quotum e-mail...\n\n", os.Args[0])
		return
	}

	var Cfg Config
	paqudir := os.Getenv("PAQU")
	if paqudir == "" {
		if DefaultPaquDir != "" {
			paqudir = DefaultPaquDir
		} else {
			paqudir = path.Join(os.Getenv("HOME"), ".paqu")
		}
	}
	_, err := toml.DecodeFile(path.Join(paqudir, "setup.toml"), &Cfg)
	util.CheckErr(err)

	if Cfg.Login[0] == '$' {
		Cfg.Login = os.Getenv(Cfg.Login[1:])
	}
	db, err := sql.Open("mysql", Cfg.Login+"?charset=utf8&parseTime=true&loc=Europe%2FAmsterdam")
	util.CheckErr(err)

	quotum, err := strconv.Atoi(os.Args[1])
	util.CheckErr(err)
	if quotum < 0 {
		fmt.Printf("Quotum kan niet negatief zijn\n")
		return
	}

	for _, user := range os.Args[2:] {
		result, err := db.Exec(fmt.Sprintf("UPDATE `%s_users` SET `quotum` = %d WHERE `mail` = %q", Cfg.Prefix, quotum, user))
		util.CheckErr(err)
		if n, _ := result.RowsAffected(); n > 0 {
			fmt.Printf("Quotum aangepast voor: %s\n", user)
		} else {
			fmt.Printf("Gebruiker niet gevonden: %s\n", user)
		}
	}
}
