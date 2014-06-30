// +build nodbxml

package main

import (
	"errors"
)

const (
	has_dbxml = false
)

func get_dact(archive, filename string) ([]byte, error) {
	return []byte{}, errors.New("Het dact-formaat wordt niet ondersteund")
}

func makeDact(dact, xml string, chKill chan bool) error {
	return nil
}
