package main

//. Imports

import (
	"github.com/rug-compling/paqu/internal/dir"

	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pebbe/util"

	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
	_, err := TomlDecodeFile(filepath.Join(dir.Config, "setup.toml"), &Cfg)
	util.CheckErr(err)

	if Cfg.Login[0] == '$' {
		Cfg.Login = os.Getenv(Cfg.Login[1:])
	}
	db, err := sql.Open("mysql", Cfg.Login+"?charset=utf8&parseTime=true&loc=Europe%2FAmsterdam&sql_mode=''")
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

	_, err = db.Exec(fmt.Sprintf(
		"DROP TABLE IF EXISTS `%s_c_%s_deprel`, `%s_c_%s_sent`, `%s_c_%s_file`, `%s_c_%s_arch`, `%s_c_%s_word`, "+
			"`%s_c_%s_meta`, `%s_c_%s_midx`, `%s_c_%s_minf`, `%s_c_%s_mval`",
		Cfg.Prefix, corpus,
		Cfg.Prefix, corpus,
		Cfg.Prefix, corpus,
		Cfg.Prefix, corpus,
		Cfg.Prefix, corpus,
		Cfg.Prefix, corpus,
		Cfg.Prefix, corpus,
		Cfg.Prefix, corpus,
		Cfg.Prefix, corpus))
	util.CheckErr(err)
	fmt.Printf("DROP TABLE IF EXISTS `%s_c_%s_*` (9): ok\n", Cfg.Prefix, corpus)

	util.CheckErr(os.RemoveAll(filepath.Join(dir.Data, "data", corpus)))
	fmt.Printf("rm -r %s: ok\n", filepath.Join(dir.Data, "data", corpus))
}

func rijen(r sql.Result) int64 {
	n, _ := r.RowsAffected()
	return n
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
