package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	fmt.Print(`//
// THIS IS A GENERATED FILE. DO NOT EDIT.
//

package main

import (
	"github.com/pebbe/util"

	"encoding/base64"
)
`)

	for _, filename := range os.Args[1:] {
		fmt.Printf("\nvar %s = `", varname(filename))
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatal(err)
		}
		encoded := base64.StdEncoding.EncodeToString(data)
		n := len(encoded)
		for i := 0; i < n; i += 80 {
			j := i + 80
			if j > n {
				j = n
			}
			fmt.Print("\n", encoded[i:j])
		}
		fmt.Println("`")
	}

	fmt.Println("\nfunc init() {\n\tvar b []byte\n\tvar err error")
	for _, filename := range os.Args[1:] {
		name := varname(filename)
		fmt.Printf("\tb, err = base64.StdEncoding.DecodeString(%s)\n", name)
		fmt.Println("\tutil.CheckErr(err)")
		fmt.Printf("\t%s = string(b)\n", name)
	}
	fmt.Println("}")
}

func varname(s string) string {
	return "file__" + strings.Replace(strings.Replace(s, ".", "__", -1), "-", "__", -1)
}
