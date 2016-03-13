package main

import (
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

func folia(prefix string, fpin io.Reader, fpout io.Writer) error {

	d := xml.NewDecoder(fpin)
	var inS, inW, inT, inCorrection, inOriginal bool
	var label, wid string
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
			case "metadata":
				var src, typ string
				for _, e := range t.Attr {
					switch e.Name.Local {
					case "src":
						src = e.Value
					case "type":
						typ = e.Value
					}
				}
				fmt.Fprintf(fpout, "##META %s\n", hex.EncodeToString([]byte("text metadata.src = "+src)))
				fmt.Fprintf(fpout, "##META %s\n", hex.EncodeToString([]byte("text metadata.type = "+typ)))
			case "whitespace":
				if len(words) > 0 {
					fmt.Fprintf(fpout, "##PAQULBL %s\n%s\n", hex.EncodeToString([]byte(label)), strings.Join(words, " "))
					words = words[0:0]
					label += ".b"
				}
			case "s":
				teller++
				label = fmt.Sprintf("%ss.%d", prefix, teller)
				for _, e := range t.Attr {
					if e.Name.Local == "id" {
						label = e.Value
						break
					}
				}
				inS = true
				inW = false
				inT = false
				inCorrection = false
				inOriginal = false
			case "w":
				wid = ""
				for _, e := range t.Attr {
					if e.Name.Local == "id" {
						wid = e.Value
						break
					}
				}
				inW = true
				inT = false
			case "t":
				inT = true
			case "correction":
				inCorrection = true
				inOriginal = false
			case "original":
				if inCorrection {
					inOriginal = true
				}
			}
		} else if t, ok := tt.(xml.EndElement); ok {
			switch t.Name.Local {
			case "s":
				if inS {
					if len(words) > 0 {
						fmt.Fprintf(fpout, "##PAQULBL %s\n%s\n", hex.EncodeToString([]byte(label)), strings.Join(words, " "))
						words = words[0:0]
					}
				}
				inS = false
				inW = false
				inT = false
			case "w":
				inW = false
				inT = false
			case "t":
				inT = false
			case "correction":
				inCorrection = false
				inOriginal = false
			case "original":
				inOriginal = false
			}
		} else if t, ok := tt.(xml.CharData); ok {
			if inS && inW && inT && !inOriginal {
				s := alpinoEscape(string(t))
				if Cfg.Alpino15 && wid != "" {
					s = fmt.Sprintf("[ @id %s ] %s", alpinoEscape(wid), s)
				}
				words = append(words, s)
				wid = ""
				inW = false
				inT = false
			}
		}
	}
	return nil
}

func alpinoEscape(s string) string {
	switch s {
	case `[`:
		return `\[`
	case `]`:
		return `\]`
	case `\[`:
		return `\\[`
	case `\]`:
		return `\\]`
	}
	return s
}
