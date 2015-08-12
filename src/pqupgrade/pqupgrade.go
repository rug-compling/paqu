package main

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

type Config struct {
	Login  string
	Prefix string
}

func main() {

	paqudir := os.Getenv("PAQU")
	if paqudir == "" {
		paqudir = path.Join(os.Getenv("HOME"), ".paqu")
	}
	var Cfg Config
	_, err := toml.DecodeFile(path.Join(paqudir, "setup.toml"), &Cfg)
	util.CheckErr(err)

	db, err := sql.Open("mysql", Cfg.Login+"?charset=utf8&parseTime=true&loc=Europe%2FAmsterdam")
	util.CheckErr(err)
	defer db.Close()

	changed := false

	////////////////////////////////////////////////////////////////

	// tabel <prefix>_info
	// veld `protected` toevoegen

	rows, err := db.Query("SELECT `protected` FROM `" + Cfg.Prefix + "_info` LIMIT 0, 1")
	if err == nil {
		rows.Close()
	} else {
		_, err := db.Exec("ALTER TABLE `" + Cfg.Prefix + "_info` ADD `protected` BOOLEAN NOT NULL DEFAULT '0'")
		util.CheckErr(err)
		changed = true
	}

	////////////////////////////////////////////////////////////////

	// tabel <prefix>_info
	// optie "QUEUING" toevoegen aan veld `status`, en default maken

	ok := false
	rows, err = db.Query("SELECT `COLUMN_TYPE`,`COLUMN_DEFAULT` FROM `information_schema`.`COLUMNS` WHERE `TABLE_NAME` = \"" +
		Cfg.Prefix + "_info\" AND `COLUMN_NAME` = \"status\"")
	util.CheckErr(err)
	for rows.Next() {
		var tp, def string
		util.CheckErr(rows.Scan(&tp, &def))
		if def == "QUEUING" && strings.Contains(tp, "QUEUING") {
			ok = true
		}
	}
	if !ok {
		_, err := db.Exec("ALTER TABLE `" + Cfg.Prefix +
			"_info` CHANGE `status` `status` ENUM('QUEUED', 'WORKING', 'FINISHED', 'FAILED', 'QUEUING') NOT NULL DEFAULT 'QUEUING'")
		util.CheckErr(err)
		changed = true

	}

	////////////////////////////////////////////////////////////////

	if changed {
		fmt.Println("Database is aangepast")
	} else {
		fmt.Println("Niets veranderd")
	}
}
