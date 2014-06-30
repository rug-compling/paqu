// +build nodbxml

package main

import (
	"errors"
)

const (
	has_dbxml = false
)

var (
	errGeenDact = errors.New("Het dact-formaat wordt niet ondersteund")
)

func get_dact(archive, filename string) ([]byte, error) {
	return []byte{}, errGeenDact
}

func makeDact(dact, xml string, chKill chan bool) error {
	return errGeenDact
}

func unpackDact(data, xmldir, dact, stderr string, chKill chan bool) (tokens, nline int, err error) {
	return 0, 0, errGeenDact
}
