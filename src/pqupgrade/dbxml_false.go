// +build nodbxml

package main

import (
	"errors"
)

var (
	errGeenDact = errors.New("Het dact-formaat wordt niet ondersteund")
)

func get_dact(archive, filename string) ([]byte, error) {
	return []byte{}, errGeenDact
}
