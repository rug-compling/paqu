package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

func tei(infile, outfile string) error {

	fpin, err := os.Open(infile)
	if err != nil {
		return err
	}
	defer fpin.Close()

	fpout, err := os.Create(outfile)
	if err != nil {
		return err
	}
	defer fpout.Close()

	d := xml.NewDecoder(fpin)
	var inS, inW, inPC bool
	var label string
	var teller uint64
	words := make([]string, 0, 500)
	for {
		tt, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if t, ok := tt.(xml.StartElement); ok {
			switch t.Name.Local {
			case "s":
				teller++
				label = fmt.Sprintf("s.%d", teller)
				for _, e := range t.Attr {
					if e.Name.Local == "id" {
						label = e.Value
						break
					}
				}
				inS = true
				inW = false
				inPC = false
			case "w":
				inW = true
				inPC = false
			case "pc":
				inPC = true
				inW = false
			}
		} else if t, ok := tt.(xml.EndElement); ok {
			switch t.Name.Local {
			case "s":
				if inS {
					if len(words) > 0 {
						fmt.Fprintf(fpout, "%s|%s\n", label, strings.Join(words, " "))
						words = words[0:0]
					}
				}
				inS = false
				inW = false
				inPC = false
			case "w":
				inW = false
			case "pc":
				inPC = false
			}
		} else if t, ok := tt.(xml.CharData); ok {
			if inS && (inW || inPC) {
				words = append(words, alpinoEscape(string(t)))
				inW = false
				inPC = false
			}
		}
	}
	return nil
}
