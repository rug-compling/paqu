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
	"strings"
)

//. Types

type Config struct {
	Login  string
	Prefix string
}

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("\nSyntax: %s corpus-ID\n\n", os.Args[0])
		return
	}

	var Cfg Config
	paqudir := os.Getenv("PAQU")
	if paqudir == "" {
		paqudir = path.Join(os.Getenv("HOME"), ".paqu")
	}
	_, err := toml.DecodeFile(path.Join(paqudir, "setup.toml"), &Cfg)
	util.CheckErr(err)

	db, err := sql.Open("mysql", Cfg.Login+"?charset=utf8&parseTime=true&loc=Europe%2FAmsterdam")
	util.CheckErr(err)

	corpus := os.Args[1]

	rows, err := db.Query(fmt.Sprintf("SELECT `owner` FROM `%s_info` WHERE `id` = %q", Cfg.Prefix, corpus))
	util.CheckErr(err)

	owner := ""
	if rows.Next() {
		util.CheckErr(rows.Scan(&owner))
		rows.Close()
	} else {
		fmt.Printf("Corpus niet gevonden in tabel `%s_info`\n", Cfg.Prefix)
	}

	n, err := db.Exec(fmt.Sprintf("DELETE FROM `%s_corpora` WHERE `prefix` = %q", Cfg.Prefix, corpus))
	util.CheckErr(err)
	fmt.Printf("DELETE FROM `%s_corpora`: %#v rij(en)\n", Cfg.Prefix, rijen(n))

	n, err = db.Exec(fmt.Sprintf("DELETE FROM `%s_ignore` WHERE `prefix` = %q", Cfg.Prefix, corpus))
	util.CheckErr(err)
	fmt.Printf("DELETE FROM `%s_ignore`: %#v rij(en)\n", Cfg.Prefix, rijen(n))

	n, err = db.Exec(fmt.Sprintf("DELETE FROM `%s_info` WHERE `id` = %q", Cfg.Prefix, corpus))
	util.CheckErr(err)
	fmt.Printf("DELETE FROM `%s_info`: %#v rij(en)\n", Cfg.Prefix, rijen(n))

	n, err = db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%s_c_%s_deprel`, `%s_c_%s_sent`, `%s_c_%s_file`, `%s_c_%s_arch`, `%s_c_%s_word`",
		Cfg.Prefix, corpus,
		Cfg.Prefix, corpus,
		Cfg.Prefix, corpus,
		Cfg.Prefix, corpus,
		Cfg.Prefix, corpus))
	util.CheckErr(err)
	fmt.Printf("DROP TABLE IF EXISTS `%s_c_%s_*` (5): ok\n", Cfg.Prefix, corpus)

	if strings.Contains(owner, "@") {
		util.CheckErr(os.RemoveAll(path.Join(paqudir, "data", corpus)))
		fmt.Printf("rm -r %s: ok\n", path.Join(paqudir, "data", corpus))
	}

}

func rijen(r sql.Result) int64 {
	n, _ := r.RowsAffected()
	return n
}
