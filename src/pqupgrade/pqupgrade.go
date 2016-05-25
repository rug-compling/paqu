package main

import (
	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pebbe/util"

	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Login  string
	Prefix string
}

var (
	DefaultPaquDir string
)

func main() {

	paqudir := os.Getenv("PAQU")
	if paqudir == "" {
		if DefaultPaquDir != "" {
			paqudir = DefaultPaquDir
		} else {
			paqudir = filepath.Join(os.Getenv("HOME"), ".paqu")
		}
	}
	var Cfg Config
	_, err := toml.DecodeFile(filepath.Join(paqudir, "setup.toml"), &Cfg)
	util.CheckErr(err)

	if Cfg.Login[0] == '$' {
		Cfg.Login = os.Getenv(Cfg.Login[1:])
	}
	db, err := sql.Open("mysql", Cfg.Login+"?charset=utf8&parseTime=true&loc=Europe%2FAmsterdam&sql_mode=''")
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
	// veld `hasmeta` toevoegen

	rows, err = db.Query("SELECT `hasmeta` FROM `" + Cfg.Prefix + "_info` LIMIT 0, 1")
	if err == nil {
		rows.Close()
	} else {
		_, err := db.Exec("ALTER TABLE `" + Cfg.Prefix + "_info` ADD `hasmeta` BOOLEAN NOT NULL DEFAULT '0'")
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

	// tabel <prefix>_info
	// veld `attr` verwijderen

	rows, err = db.Query("SELECT `attr` FROM `" + Cfg.Prefix + "_info` LIMIT 0, 1")
	if err == nil {
		rows.Close()
		_, err := db.Exec("ALTER TABLE `" + Cfg.Prefix + "_info` DROP `attr`")
		util.CheckErr(err)
		changed = true
	}

	////////////////////////////////////////////////////////////////

	// tabellen *_deprel
	// veld `idd` toevoegen

	tables := make([]string, 0)
	rows, err = db.Query("SELECT `id` FROM `" + Cfg.Prefix + "_info`")
	util.CheckErr(err)
	for rows.Next() {
		var t string
		util.CheckErr(rows.Scan(&t))
		tables = append(tables, t)
	}
	util.CheckErr(rows.Err())

	for _, table := range tables {
		tb := Cfg.Prefix + "_c_" + table + "_deprel"
		rows, err = db.Query("SELECT `idd` FROM `" + tb + "`")
		if err == nil {
			rows.Close()
			continue
		}
		fmt.Print("Toevoegen van kolom `idd` aan tabel ", tb, "...")
		_, err := db.Exec("ALTER TABLE `" + tb + "` ADD `idd` INT NOT NULL AUTO_INCREMENT PRIMARY KEY FIRST")
		fmt.Println()
		if util.WarnErr(err) != nil {
			// Misschien bestaat de tabel helemaal niet, omdat er een fout was met het corpus
			continue
		}
		changed = true
	}

	////////////////////////////////////////////////////////////////

	if changed {
		fmt.Println("Database is aangepast")
	} else {
		fmt.Println("Niets veranderd")
	}
}
