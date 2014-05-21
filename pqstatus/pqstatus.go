package main

//. Imports

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/pebbe/util"

	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

//. Types

type Config struct {
	Data string

	Login  string
	Prefix string

	Sh   string
	Path string
}

var (
	Cfg Config
)

//. Main

func main() {
	data, err := ioutil.ReadFile(os.Args[1])
	util.CheckErr(err)
	util.CheckErr(json.Unmarshal(data, &Cfg))

	util.CheckErr(os.Chdir(Cfg.Data))

	diskuse := make(map[string]int)
	totaldiskuse := 0
	cmd := exec.Command(
		Cfg.Sh,
		"-c",
		fmt.Sprintf(
			`PATH=%s; export PATH; `+
				`du --block-size=1 --max-depth=1 .`,
			Cfg.Path))
	b, err := cmd.Output()
	util.CheckErr(err)
	for _, line := range strings.Split(string(b), "\n") {
		a := strings.Fields(line)
		if len(a) != 2 {
			continue
		}
		value, err := strconv.Atoi(a[0])
		util.CheckErr(err)
		if a[1] == "." {
			totaldiskuse = value
		} else {
			diskuse[a[1][2:]] = value
		}
	}

	db, err := dbopen()
	util.CheckErr(err)

	dbuse := make(map[string]int)
	totaldbuse := 0
	rows, err := db.Query("SELECT TABLE_NAME, DATA_LENGTH, INDEX_LENGTH FROM `information_schema`.`TABLES` WHERE TABLE_NAME LIKE '" +
		Cfg.Prefix + "\\_c\\_%';")
	util.CheckErr(err)
	for rows.Next() {
		var id string
		var d, i int
		util.CheckErr(rows.Scan(&id, &d, &i))
		id = strings.Split(id, "_")[2]
		dbuse[id] += (d + i)
		totaldbuse += (d + i)
	}

	useract := make(map[string]time.Time)
	rows, err = db.Query(
		fmt.Sprintf(
			"SELECT `mail`, `active` FROM `%s_users`",
			Cfg.Prefix))
	util.CheckErr(err)
	for rows.Next() {
		var id string
		var act time.Time
		util.CheckErr(rows.Scan(&id, &act))
		useract[id] = act
	}

	rows, err = db.Query(
		fmt.Sprintf(
			"SELECT `id`, `description`, `owner`, `status`, `nline`, `shared`, `msg`, `active` FROM `%s_info` ORDER BY `owner`, `description`",
			Cfg.Prefix))
	util.CheckErr(err)

	queued := 0
	working := 0
	finished := 0
	failed := 0
	nowners := 0
	prevowner := ""
	disksum := 0
	dbsum := 0
	for rows.Next() {
		var id, description, owner, status, shared, msg string
		var zinnen int
		var active time.Time
		util.CheckErr(rows.Scan(&id, &description, &owner, &status, &zinnen, &shared, &msg, &active))
		if owner != prevowner {
			if prevowner != "" {
				fmt.Printf("  Totaal disk: %s\n  Totaal db  : %s\n", size(disksum), size(dbsum))
			}
			disksum = 0
			dbsum = 0
			prevowner = owner
			nowners++

			fmt.Printf("\nGebruiker: %s\n", owner)
			if t, ok := useract[owner]; ok {
				fmt.Printf("   Actief: %v\n", t)
			}
			fmt.Println()
		}
		dbsum += dbuse[id]
		disksum += diskuse[id]
		fmt.Printf("  Corpus: %s\n", id)
		fmt.Printf("   Titel: %s\n", description)
		fmt.Printf("  Status: %s  %s\n", status, msg)
		fmt.Printf("  Zinnen: %s\n", iformat(zinnen))
		fmt.Printf("  Shared: %s\n", shared)
		fmt.Printf("    Disk: %s\n", size(diskuse[id]))
		fmt.Printf("      Db: %s\n", size(dbuse[id]))
		fmt.Printf("  Actief: %v\n\n", active)
		switch status {
		case "QUEUED":
			queued++
		case "WORKING":
			working++
		case "FINISHED":
			finished++
		case "FAILED":
			failed++
		}
		delete(diskuse, id)
	}
	util.CheckErr(rows.Err())
	if prevowner != "" {
		fmt.Printf("  Totaal disk: %s\n  Totaal db  : %s\n", size(disksum), size(dbsum))
	}

	fmt.Printf("\nAantal gebruikers:   %d\n", nowners)
	fmt.Printf("Totaal schijfruimte: %s\n", size(totaldiskuse))
	fmt.Printf("Totaal database:     %s\n", size(totaldbuse))

	fmt.Printf(`
Opdrachten:

%8d  QUEUED
%8d  WORKING
%8d  FINISHED
%8d  FAILED

`, queued, working, finished, failed)

	if len(diskuse) > 0 {
		fmt.Println("Ongebruikte directory's:\n")
		for d := range diskuse {
			fmt.Println("  ", d)
		}
		fmt.Println()
	}
}

func size(v int) string {
	if v < 1024 {
		return fmt.Sprint(v, " b")
	} else if v < 1024*1024 {
		return fmt.Sprintf("%.1f Kb", float64(v)/1024.0)
	} else if v < 1024*1024*1024 {
		return fmt.Sprintf("%.1f Mb", float64(v)/1024.0/1024.0)
	} else {
		return fmt.Sprintf("%.1f Gb", float64(v)/1024.0/1024.0/1024.0)
	}
}

func iformat(i int) string {
	s1 := fmt.Sprint(i)
	s2 := ""
	for n := len(s1); n > 3; n = len(s1) {
		s2 = "." + s1[n-3:n] + s2
		s1 = s1[0 : n-3]
	}
	return s1 + s2
}

func dbopen() (*sql.DB, error) {
	return sql.Open("mysql", Cfg.Login+"?charset=utf8mb4,utf8&parseTime=true&loc=Europe%2FAmsterdam")
}
