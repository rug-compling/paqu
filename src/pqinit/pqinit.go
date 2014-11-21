package main

import (
	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pebbe/util"

	"database/sql"
	"fmt"
	"os"
	"path"
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

 -w : bestaande database overschrijven

`, os.Args[0])
		return
	}

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
		status      enum('QUEUED','WORKING','FINISHED','FAILED') NOT NULL DEFAULT 'QUEUED',
		msg         varchar(256) NOT NULL,
		nline       int          NOT NULL DEFAULT 0,
		nword       int          NOT NULL DEFAULT 0,
		params      varchar(128) NOT NULL,
		shared      enum('PRIVATE','PUBLIC','SHARED') NOT NULL DEFAULT 'PRIVATE',
		created     timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
		active      datetime     NOT NULL,
		attr        text         NOT NULL,
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
		active datetime     NOT NULL,
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
