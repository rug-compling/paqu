// +build !nodbxml

package main

import (
	"github.com/pebbe/dbxml"

	"fmt"
	"os"
)

func main() {

	ok := true
	for i, arg := range os.Args[1:] {
		err := dbxml.Check(arg)
		if err != nil {
			fmt.Println(i+1, err)
			ok = false
		}
	}

	if ok {
		fmt.Println("OK")
	}
}
