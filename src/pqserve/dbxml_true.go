// +build !nodbxml

package main

import (
	"github.com/pebbe/dbxml"
)

const (
	has_dbxml = true
)

func get_dact(archive, filename string) ([]byte, error) {
	reader, err := dbxml.Open(archive)
	if err != nil {
		return []byte{}, err
	}
	defer reader.Close()
	d, err := reader.Get(filename)
	if err != nil {
		return []byte{}, err
	}
	return []byte(d), nil
}









