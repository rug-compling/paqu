// +build !nodbxml

package main

import (
	"github.com/pebbe/dbxml"

	"fmt"
	"os"
)

func main() {
	db, err := dbxml.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = db.Prepare(os.Args[2])
	if err != nil {
		fmt.Println(err)
		db.Close()
		return
	}

	db.Close()

	fmt.Println("OK")
}
