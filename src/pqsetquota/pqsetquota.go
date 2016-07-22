package main

//. Imports

import (
	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pebbe/util"

	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
			paqudir = filepath.Join(os.Getenv("HOME"), ".paqu")
		}
	}
	_, err := TomlDecodeFile(filepath.Join(paqudir, "setup.toml"), &Cfg)
	util.CheckErr(err)

	if Cfg.Login[0] == '$' {
		Cfg.Login = os.Getenv(Cfg.Login[1:])
	}
	db, err := sql.Open("mysql", Cfg.Login+"?charset=utf8&parseTime=true&loc=Europe%2FAmsterdam&sql_mode=''")
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
