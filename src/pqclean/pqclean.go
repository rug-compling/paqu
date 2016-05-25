package main

//. Imports

import (
	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pebbe/util"

	"database/sql"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
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
	DefaultPaquDir string
	Cfg            Config
)

//. Main

func main() {

	if len(os.Args) != 2 || os.Args[1] != "-c" {
		fmt.Printf("\nSyntax: %s -c\n\nDit verwijdert alle gebruikers zonder corpora die twee maanden niet actief zijn geweest\n\n", os.Args[0])
		return
	}

	paqudir := os.Getenv("PAQU")
	if paqudir == "" {
		if DefaultPaquDir != "" {
			paqudir = DefaultPaquDir
		} else {
			paqudir = filepath.Join(os.Getenv("HOME"), ".paqu")
		}
	}
	_, err := toml.DecodeFile(filepath.Join(paqudir, "setup.toml"), &Cfg)
	util.CheckErr(err)

	db, err := dbopen()
	util.CheckErr(err)
	defer db.Close()

	rows, err := db.Query(
		fmt.Sprintf(
			"SELECT `owner` FROM `%s_info`",
			Cfg.Prefix))
	util.CheckErr(err)

	owners := make(map[string]bool)

	for rows.Next() {
		var owner string
		util.CheckErr(rows.Scan(&owner))
		owners[owner] = true
	}

	verwijder := make([]string, 0)

	rows, err = db.Query(
		fmt.Sprintf(
			"SELECT `mail` FROM `%s_users` WHERE `active` < %q",
			Cfg.Prefix,
			datetime(time.Now().AddDate(0, -2, 0))))
	util.CheckErr(err)
	for rows.Next() {
		var mail string
		util.CheckErr(rows.Scan(&mail))
		if !owners[mail] {
			verwijder = append(verwijder, mail)
		}
	}

	for _, user := range verwijder {
		fmt.Println("Verwijderen:", user)
		_, err = db.Exec(fmt.Sprintf("DELETE FROM `%s_corpora` WHERE `user` = %q", Cfg.Prefix, user))
		util.CheckErr(err)
		_, err = db.Exec(fmt.Sprintf("DELETE FROM `%s_ignore` WHERE `user` = %q", Cfg.Prefix, user))
		util.CheckErr(err)
		_, err = db.Exec(fmt.Sprintf("DELETE FROM `%s_macros` WHERE `user` = %q", Cfg.Prefix, user))
		util.CheckErr(err)
		_, err = db.Exec(fmt.Sprintf("DELETE FROM `%s_users` WHERE `mail` = %q", Cfg.Prefix, user))
		util.CheckErr(err)
		util.CheckErr(os.RemoveAll(filepath.Join(paqudir, "folia", hex.EncodeToString([]byte(user)))))
	}
}

func datetime(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

func dbopen() (*sql.DB, error) {
	if Cfg.Login[0] == '$' {
		Cfg.Login = os.Getenv(Cfg.Login[1:])
	}
	return sql.Open("mysql", Cfg.Login+"?charset=utf8&parseTime=true&loc=Europe%2FAmsterdam&sql_mode=''")
}
