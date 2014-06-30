// +build !nodbxml

package main

import (
	"github.com/pebbe/dbxml"

	"errors"
	"io/ioutil"
	"os"
	"path"
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

func makeDact(dact, xml string, chKill chan bool) error {
	files, err := ioutil.ReadDir(xml)
	if err != nil {
		return err
	}

	os.Remove(dact)
	db, err := dbxml.Open(dact)
	if err != nil {
		return err
	}
	defer db.Close()

	for _, file := range files {

		select {
		case <-chGlobalExit:
			return errors.New("Global Exit")
		case <-chKill:
			return errors.New("Killed")
		default:
		}

		name := file.Name()
		data, err := ioutil.ReadFile(path.Join(xml, name))
		if err != nil {
			return err
		}
		err = db.PutXml(name, string(data), false)
		if err != nil {
			return err
		}
	}
	return nil
}
