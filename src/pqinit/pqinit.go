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
)

type Config struct {
	Login  string
	Prefix string
}

func main() {

	db_overwrite := false
	for len(os.Args) > 1 {
		if os.Args[1] == "-w" {
			db_overwrite = true
		} else {
			break
		}
		os.Args = append(os.Args[:1], os.Args[2:]...)
	}

	if len(os.Args) > 1 {
		fmt.Printf(`
Syntax: %s [-w]

 -w : bestaande database overschrijven : OUDE GEGEVENS WORDEN GEWIST

`, os.Args[0])
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
	defer db.Close()

	if db_overwrite {
		_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%s_corpora`, `%s_info`, `%s_ignore`, `%s_users`, `%s_macros`; ",
			Cfg.Prefix, Cfg.Prefix, Cfg.Prefix, Cfg.Prefix, Cfg.Prefix))
		util.CheckErr(err)
	}

	_, err = db.Exec(`CREATE TABLE ` + Cfg.Prefix + `_corpora (
		user    varchar(128) NOT NULL,
		prefix  varchar(128) NOT NULL,
		enabled tinyint      NOT NULL DEFAULT 1,
		INDEX (user),
		INDEX (prefix),
		INDEX (enabled))
		DEFAULT CHARACTER SET utf8;`)
	util.CheckErr(err)

	_, err = db.Exec(`CREATE TABLE ` + Cfg.Prefix + `_info (
		id          varchar(128) NOT NULL,
		description varchar(128) NOT NULL COLLATE utf8_unicode_ci,
		owner       varchar(128) NOT NULL DEFAULT 'none',
		status      enum('QUEUED','WORKING','FINISHED','FAILED','QUEUING') NOT NULL DEFAULT 'QUEUING',
		msg         varchar(256) NOT NULL,
		nline       int          NOT NULL DEFAULT 0,
		nword       int          NOT NULL DEFAULT 0,
		params      varchar(128) NOT NULL,
		shared      enum('PRIVATE','PUBLIC','SHARED') NOT NULL DEFAULT 'PRIVATE',
		created     timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
		active      datetime     NOT NULL DEFAULT "1000-01-01 00:00:00",
        protected   boolean      NOT NULL DEFAULT 0,
        hasmeta     boolean      NOT NULL DEFAULT 0,
		version     int          NOT NULL DEFAULT 0,
		UNIQUE INDEX (id),
		INDEX (description),
		INDEX (owner),
		INDEX (status),
		INDEX (created),
		INDEX (active))
		DEFAULT CHARACTER SET utf8;`)
	util.CheckErr(err)

	_, err = db.Exec(`CREATE TABLE ` + Cfg.Prefix + `_ignore (
		user    varchar(128) NOT NULL,
		prefix  varchar(128) NOT NULL,
		INDEX (user),
		INDEX (prefix));`)
	util.CheckErr(err)

	_, err = db.Exec(`CREATE TABLE ` + Cfg.Prefix + `_users (
		mail   varchar(128) NOT NULL,
		sec    char(16)     NOT NULL,
		pw     char(16)     NOT NULL,
		active datetime     NOT NULL DEFAULT "1000-01-01 00:00:00",
		quotum int          NOT NULL DEFAULT 0,
		UNIQUE INDEX (mail),
		INDEX (active))
		DEFAULT CHARACTER SET utf8;`)
	util.CheckErr(err)

	_, err = db.Exec(`CREATE TABLE ` + Cfg.Prefix + `_macros (
		user   varchar(128) NOT NULL,
		macros text         NOT NULL,
		UNIQUE INDEX (user))
		DEFAULT CHARACTER SET utf8;`)
	util.CheckErr(err)
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
