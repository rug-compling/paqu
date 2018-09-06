package main

//. Imports

import (
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
	"strings"
)

//. Types

type Config struct {
	Login  string
	Prefix string
}

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("\nSyntax: %s e-mail\n\n", os.Args[0])
		return
	}

	var Cfg Config
	_, err := TomlDecodeFile(filepath.Join(paquconfigdir, "setup.toml"), &Cfg)
	util.CheckErr(err)

	if Cfg.Login[0] == '$' {
		Cfg.Login = os.Getenv(Cfg.Login[1:])
	}
	db, err := sql.Open("mysql", Cfg.Login+"?charset=utf8&parseTime=true&loc=Europe%2FAmsterdam&sql_mode=''")
	util.CheckErr(err)

	user := os.Args[1]

	if strings.Index(user, "@") < 0 {
		fmt.Printf("Gebruiker %s kan niet verwijderd worden\n", user)
		return
	}

	corpora := make([]string, 0)
	rows, err := db.Query(fmt.Sprintf("SELECT `id`, `status` FROM `%s_info` WHERE `owner` = %q", Cfg.Prefix, user))
	util.CheckErr(err)
	for rows.Next() {
		var id, status string
		util.CheckErr(rows.Scan(&id, &status))
		if status != "FINISHED" && status != "FAILED" {
			fmt.Println("Gebruiker heeft lopende jobs of jobs in wachtrij. Gebruiker kan niet verwijderd worden.")
			return
		}
		corpora = append(corpora, id)
	}

	for _, corpus := range corpora {
		fmt.Printf("Verwijdering corpus: %s\n", corpus)

		_, err := db.Exec(fmt.Sprintf("DELETE FROM `%s_corpora` WHERE `prefix` = %q", Cfg.Prefix, corpus))
		util.CheckErr(err)

		_, err = db.Exec(fmt.Sprintf("DELETE FROM `%s_ignore` WHERE `prefix` = %q", Cfg.Prefix, corpus))
		util.CheckErr(err)

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

		util.CheckErr(os.RemoveAll(filepath.Join(paqudatadir, "data", corpus)))

		// deze pas als de rest goed ging
		_, err = db.Exec(fmt.Sprintf("DELETE FROM `%s_info` WHERE `id` = %q", Cfg.Prefix, corpus))
		util.CheckErr(err)
	}

	_, err = db.Exec(fmt.Sprintf("DELETE FROM `%s_corpora` WHERE `user` = %q", Cfg.Prefix, user))
	util.CheckErr(err)

	_, err = db.Exec(fmt.Sprintf("DELETE FROM `%s_ignore` WHERE `user` = %q", Cfg.Prefix, user))
	util.CheckErr(err)

	_, err = db.Exec(fmt.Sprintf("DELETE FROM `%s_macros` WHERE `user` = %q", Cfg.Prefix, user))
	util.CheckErr(err)

	result, err := db.Exec(fmt.Sprintf("DELETE FROM `%s_users` WHERE `mail` = %q", Cfg.Prefix, user))
	util.CheckErr(err)
	if n, _ := result.RowsAffected(); n > 0 {
		fmt.Printf("Verwijdering gebruiker: %s\n", user)
	} else {
		fmt.Printf("Gebruiker niet gevonden: %s\n", user)
	}

	util.CheckErr(os.RemoveAll(filepath.Join(paqudatadir, "folia", hex.EncodeToString([]byte(user)))))
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
