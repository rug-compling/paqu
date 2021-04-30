package main

import (
	"github.com/rug-compling/paqu/internal/dir"

	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pebbe/util"

	"bytes"
	"compress/gzip"
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

var (
	x = util.CheckErr
)

func main() {

	var Cfg Config
	_, err := TomlDecodeFile(filepath.Join(dir.Config, "setup.toml"), &Cfg)
	x(err)

	if Cfg.Login[0] == '$' {
		Cfg.Login = os.Getenv(Cfg.Login[1:])
	}
	db, err := sql.Open("mysql", Cfg.Login+"?charset=utf8&parseTime=true&loc=Europe%2FAmsterdam&sql_mode=''")
	x(err)
	defer db.Close()

	changed := false

	////////////////////////////////////////////////////////////////

	// tabel <prefix>_version aanmaken

	version := 0

	rows, err := db.Query("SELECT `version` FROM `" + Cfg.Prefix + "_version` WHERE `id` = 1 LIMIT 0, 1")
	if err == nil {
		if rows.Next() {
			x(rows.Scan(&version))
			rows.Close()
		} else {
			x(fmt.Errorf("Missing row"))
		}
	} else {
		_, err = db.Exec("CREATE TABLE " + Cfg.Prefix + "_version (id int NOT NULL, version int NOT NULL DEFAULT 0, UNIQUE INDEX (id))")
		x(err)
		_, err = db.Exec("INSERT `" + Cfg.Prefix + "_version` (`id`,`version`) VALUES (1,0);")
		x(err)
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
			x(err)
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
			x(err)
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
		x(err)
		for rows.Next() {
			var tp, def string
			x(rows.Scan(&tp, &def))
			if def == "QUEUING" && strings.Contains(tp, "QUEUING") {
				ok = true
			}
		}
		if !ok {
			_, err := db.Exec("ALTER TABLE `" + Cfg.Prefix +
				"_info` CHANGE `status` `status` ENUM('QUEUED', 'WORKING', 'FINISHED', 'FAILED', 'QUEUING') NOT NULL DEFAULT 'QUEUING'")
			x(err)
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
			x(err)
			changed = true
		}

	}

	////////////////////////////////////////////////////////////////

	// tabellen *_deprel
	// veld `idd` toevoegen

	if version < 3 {

		tables := make([]string, 0)
		rows, err = db.Query("SELECT `id` FROM `" + Cfg.Prefix + "_info`")
		x(err)
		for rows.Next() {
			var t string
			x(rows.Scan(&t))
			tables = append(tables, t)
		}
		x(rows.Err())

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
			if util.WarnErr(err, tb) != nil {
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
		x(err)
		for rows.Next() {
			var id, o string
			x(rows.Scan(&id, &o))
			if strings.Contains(o, "@") {
				tables = append(tables, id)
			}
		}
		x(rows.Err())

		for _, table := range tables {
			tb := Cfg.Prefix + "_c_" + table + "_file"
			rows, err = db.Query("SELECT `id`,`file` FROM `" + tb + "`")
			fmt.Println("Upgrade prefix in tabel", tb, "...")
			if util.WarnErr(err, tb) == nil {
				// Misschien bestaat de tabel helemaal niet, omdat er een fout was met het corpus
				p := "/data/" + table + "/xml"
				ln := len(p)
				for rows.Next() {
					var id, filename string
					x(rows.Scan(&id, &filename))
					i := strings.Index(filename, p)
					if i >= 0 {
						name := "$$" + filename[i+ln:]
						_, err = db.Exec(fmt.Sprintf("UPDATE `%s` SET `file` = %q WHERE `id` = %q", tb, name, id))
						x(err)
						changed = true
					}
				}
				x(rows.Err())
			}
		}
	}

	////////////////////////////////////////////////////////////////

	// tabel <prefix>_info
	// velden `info` en `infop` toevoegen

	if version < 4 {

		rows, err = db.Query("SELECT `info`, `infop` FROM `" + Cfg.Prefix + "_info` LIMIT 0, 1")
		if err == nil {
			rows.Close()
		} else {
			_, err := db.Exec("ALTER TABLE `" + Cfg.Prefix + "_info` ADD `info` TEXT NOT NULL DEFAULT ''")
			x(err)
			_, err = db.Exec("ALTER TABLE `" + Cfg.Prefix + "_info` ADD `infop` TEXT NOT NULL DEFAULT ''")
			x(err)
			changed = true
		}

	}

	////////////////////////////////////////////////////////////////

	// tabel <prefix>_info
	// veld `hasud` toevoegen

	if version < 5 {
		rows, err = db.Query("SELECT `hasud` FROM `" + Cfg.Prefix + "_info` LIMIT 0, 1")
		if err == nil {
			rows.Close()
		} else {
			_, err := db.Exec("ALTER TABLE `" + Cfg.Prefix + "_info` ADD `hasud` BOOLEAN NOT NULL DEFAULT 0 AFTER `hasmeta`")
			x(err)
			changed = true

			IDs := make([]string, 0)
			rows, err = db.Query("SELECT `id` FROM `" + Cfg.Prefix + "_info`")
			x(err)
			for rows.Next() {
				var id string
				x(rows.Scan(&id))
				IDs = append(IDs, id)
			}
			for _, id := range IDs {
				archID := -1
				fileID := -1
				rows, err := db.Query("SELECT `arch`,`file` FROM `" + Cfg.Prefix + "_c_" + id + "_sent` LIMIT 0, 1")
				if util.WarnErr(err, id) != nil {
					// Misschien bestaat de tabel helemaal niet, omdat er een fout was met het corpus
					continue
				}
				for rows.Next() {
					x(rows.Scan(&archID, &fileID))
				}
				if fileID < 0 {
					x(fmt.Errorf("No sentence found for %s", id))
				}
				var archname, filename string
				if archID >= 0 {
					rows, err := db.Query(fmt.Sprintf("SELECT `arch` FROM `%s_c_%s_arch` WHERE `id` = %d", Cfg.Prefix, id, archID))
					x(err)
					for rows.Next() {
						x(rows.Scan(&archname))
					}
					if archname == "" {
						x(fmt.Errorf("No arch %d found for %s", archID, id))
					}
				}
				rows, err = db.Query(fmt.Sprintf("SELECT `file` FROM `%s_c_%s_file` WHERE `id` = %d", Cfg.Prefix, id, fileID))
				x(err)
				for rows.Next() {
					x(rows.Scan(&filename))
				}
				if filename == "" {
					x(fmt.Errorf("No file found for %s", id))
				}
				if strings.HasPrefix(filename, "$$") {
					filename = dir.Data + "/data/" + id + "/xml" + filename[2:]
				}

				var b []byte
				if archname != "" {
					b, err = get_dact(archname, filename)
					if util.WarnErr(err, id, archname, filename) != nil {
						continue
					}
				} else {
					b, err = ioutil.ReadFile(filename)
					if err != nil {
						fp, err := os.Open(filename + ".gz")
						x(err)
						rd, err := gzip.NewReader(fp)
						x(err)
						b, err = ioutil.ReadAll(rd)
						x(err)
						x(rd.Close())
						x(fp.Close())
					}
				}
				if bytes.Contains(b, []byte("<conllu")) {
					_, err := db.Exec(fmt.Sprintf("UPDATE `%s_info` SET `hasud` = 1 WHERE `id` = %q", Cfg.Prefix, id))
					x(err)
				}
			}
		}
	}

	////////////////////////////////////////////////////////////////

	// upgrade naar versie 5

	if version < 5 {

		fmt.Printf("Upgrade from version %d to 5\n", version)

		result, err := db.Exec(fmt.Sprintf("UPDATE `%s_version` SET `version` = 5 WHERE `id` = 1", Cfg.Prefix))
		x(err)
		n, err := result.RowsAffected()
		x(err)
		if n < 1 {
			x(fmt.Errorf("Version update failed"))
		}

		version = 5
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
