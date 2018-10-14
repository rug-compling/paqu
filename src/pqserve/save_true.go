// +build !nodbxml

package main

import (
	"github.com/pebbe/dbxml"

	"fmt"
	"strings"
)

func saveOpenDact(q *Context, prefix string, arch int) (interface{}, string) {
	rows, err := q.db.Query(fmt.Sprintf("SELECT `arch` FROM `%s_c_%s_arch` WHERE `id` = %d", Cfg.Prefix, prefix, arch))
	if doErr(q, err) {
		return nil, ""
	}
	var filename string
	for rows.Next() {
		err := rows.Scan(&filename)
		if doErr(q, err) {
			rows.Close()
			return nil, ""
		}
	}
	if doErr(q, rows.Err()) {
		return nil, ""
	}
	if filename == "" {
		return nil, ""
	}

	if strings.HasPrefix(filename, "$$") {
		filename = paqudatadir + filename[2:]
	}
	db, err := dbxml.OpenRead(filename)
	if doErr(q, err) {
		return nil, ""
	}

	return db, filename
}

func saveGetDact(q *Context, dact interface{}, filename string) []byte {

	d, ok := dact.(*dbxml.Db)
	if !ok {
		return []byte{}
	}

	s, err := d.Get(filename)
	if doErr(q, err) {
		return []byte{}
	}

	return []byte(s)
}

func saveCloseDact(dact interface{}) {
	if d, ok := dact.(*dbxml.Db); ok {
		d.Close()
	}
}
