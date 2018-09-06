package main

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
	"strings"
)

type Config struct {
	Login  string
	Prefix string
}

func main() {

	var Cfg Config
	_, err := TomlDecodeFile(filepath.Join(paquconfigdir, "setup.toml"), &Cfg)
	util.CheckErr(err)

	if Cfg.Login[0] == '$' {
		Cfg.Login = os.Getenv(Cfg.Login[1:])
	}
	db, err := sql.Open("mysql", Cfg.Login+"?charset=utf8&parseTime=true&loc=Europe%2FAmsterdam&sql_mode=''")
	util.CheckErr(err)
	defer db.Close()

	changed := false

	////////////////////////////////////////////////////////////////

	// tabel <prefix>_version aanmaken

	version := 0

	rows, err := db.Query("SELECT `version` FROM `" + Cfg.Prefix + "_version` WHERE `id` = 1 LIMIT 0, 1")
	if err == nil {
		if rows.Next() {
			util.CheckErr(rows.Scan(&version))
			rows.Close()
		} else {
			util.CheckErr(fmt.Errorf("Missing row"))
		}
	} else {
		_, err = db.Exec("CREATE TABLE " + Cfg.Prefix + "_version (id int NOT NULL, version int NOT NULL DEFAULT 0, UNIQUE INDEX (id))")
		util.CheckErr(err)
		_, err = db.Exec("INSERT `" + Cfg.Prefix + "_version` (`id`,`version`) VALUES (1,0);")
		util.CheckErr(err)
		changed = true
	}

	////////////////////////////////////////////////////////////////

	// tabel <prefix>_info
	// veld `protected` toevoegen

	if version < 3 {

		rows, err = db.Query("SELECT `protected` FROM `" + Cfg.Prefix + "_info` LIMIT 0, 1")
		if err == nil {
			rows.Close()
		} else {
			_, err := db.Exec("ALTER TABLE `" + Cfg.Prefix + "_info` ADD `protected` BOOLEAN NOT NULL DEFAULT '0'")
			util.CheckErr(err)
			changed = true
		}

	}

	////////////////////////////////////////////////////////////////

	// tabel <prefix>_info
	// veld `hasmeta` toevoegen

	if version < 3 {

		rows, err = db.Query("SELECT `hasmeta` FROM `" + Cfg.Prefix + "_info` LIMIT 0, 1")
		if err == nil {
			rows.Close()
		} else {
			_, err := db.Exec("ALTER TABLE `" + Cfg.Prefix + "_info` ADD `hasmeta` BOOLEAN NOT NULL DEFAULT '0'")
			util.CheckErr(err)
			changed = true
		}

	}

	////////////////////////////////////////////////////////////////

	// tabel <prefix>_info
	// optie "QUEUING" toevoegen aan veld `status`, en default maken

	if version < 3 {

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

	}

	////////////////////////////////////////////////////////////////

	// tabel <prefix>_info
	// veld `attr` verwijderen

	if version < 3 {

		rows, err = db.Query("SELECT `attr` FROM `" + Cfg.Prefix + "_info` LIMIT 0, 1")
		if err == nil {
			rows.Close()
			_, err := db.Exec("ALTER TABLE `" + Cfg.Prefix + "_info` DROP `attr`")
			util.CheckErr(err)
			changed = true
		}

	}

	////////////////////////////////////////////////////////////////

	// tabellen *_deprel
	// veld `idd` toevoegen

	if version < 3 {

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

	}

	////////////////////////////////////////////////////////////////

	// prefix van filename vervangen door $$ voor user-corpora

	if version < 3 {

		tables := make([]string, 0)
		rows, err = db.Query("SELECT `id`,`owner` FROM `" + Cfg.Prefix + "_info`")
		util.CheckErr(err)
		for rows.Next() {
			var id, o string
			util.CheckErr(rows.Scan(&id, &o))
			if strings.Contains(o, "@") {
				tables = append(tables, id)
			}
		}
		util.CheckErr(rows.Err())

		for _, table := range tables {
			tb := Cfg.Prefix + "_c_" + table + "_file"
			rows, err = db.Query("SELECT `id`,`file` FROM `" + tb + "`")
			fmt.Println("Upgrade prefix in tabel", tb, "...")
			if util.WarnErr(err) == nil {
				// Misschien bestaat de tabel helemaal niet, omdat er een fout was met het corpus
				p := "/data/" + table + "/xml"
				ln := len(p)
				for rows.Next() {
					var id, filename string
					util.CheckErr(rows.Scan(&id, &filename))
					i := strings.Index(filename, p)
					if i >= 0 {
						name := "$$" + filename[i+ln:]
						_, err = db.Exec(fmt.Sprintf("UPDATE `%s` SET `file` = %q WHERE `id` = %q", tb, name, id))
						util.CheckErr(err)
						changed = true
					}
				}
				util.CheckErr(rows.Err())
			}
		}
	}

	////////////////////////////////////////////////////////////////

	// upgrade naar versie 3

	if version < 3 {

		fmt.Printf("Upgrade from version %d to 3\n", version)

		result, err := db.Exec(fmt.Sprintf("UPDATE `%s_version` SET `version` = 3 WHERE `id` = 1", Cfg.Prefix))
		util.CheckErr(err)
		n, err := result.RowsAffected()
		util.CheckErr(err)
		if n < 1 {
			util.CheckErr(fmt.Errorf("Version update failed"))
		}

		version = 3
		changed = true
	}

	////////////////////////////////////////////////////////////////
	//
	// Is pqinit aangepast aan het laatste versienummer?
	//
	////////////////////////////////////////////////////////////////

	if changed {
		fmt.Println("Database is aangepast")
	} else {
		fmt.Println("Niets veranderd")
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
