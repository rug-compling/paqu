// +build !nodbxml

package main

import (
	"github.com/pebbe/dbxml"
	"github.com/pebbe/util"

	"fmt"
)

const (
	has_dbxml = true
)

func do_dact(filename string) {
	reader, err := dbxml.Open(filename)
	util.CheckErr(err)
	fmt.Println(">>>", filename)
	docs, err := reader.All()
	util.CheckErr(err)
	for docs.Next() {
		do_data(filename, docs.Name(), []byte(docs.Content()))
	}
	showmemstats()
	reader.Close()
}
