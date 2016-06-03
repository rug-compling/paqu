package main

import (
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

type FoliaState struct {
	inUtt  bool
	inS    bool
	inW    bool
	inT    bool
	inSkip bool
}

func folia(prefix string, fpin io.Reader, fpout io.Writer) error {

	statestack := make([]FoliaState, 1, 10)

	d := xml.NewDecoder(fpin)
	var label string
	text := make([]byte, 0)
	var teller, uttteller uint64

	for {
		tt, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if t, ok := tt.(xml.StartElement); ok {

			state := statestack[len(statestack)-1]

			hasClass := false
			id := ""
			src := ""
			for _, e := range t.Attr {
				switch e.Name.Local {
				case "id":
					id = e.Value
				case "src":
					src = e.Value
				case "class":
					hasClass = true
				case "auth":
					state.inSkip = true
				}
			}

			switch t.Name.Local {
			case "metadata":
				fmt.Fprintf(fpout, "##META %s\n", hex.EncodeToString([]byte("text metadata.src = "+src)))
			case "s":
				teller++
				if !state.inSkip && !state.inS {
					if id == "" {
						label = fmt.Sprintf("%ss.%d", prefix, teller)
					} else {
						label = id
					}
					text = text[0:0]
					state.inS = true
					state.inW = false
					state.inT = false
				}
			case "utt":
				uttteller++
				if !state.inSkip && !state.inS { // inS, niet inUtt
					if id == "" {
						label = fmt.Sprintf("%sutt.%d", prefix, uttteller)
					} else {
						label = id
					}
					text = text[0:0]
					state.inUtt = true
					state.inW = false
					state.inT = false
				}
			case "w":
				state.inW = true
				state.inT = false
			case "t":
				if !hasClass {
					state.inT = true
				}
			case "morpheme", "str":
				state.inSkip = true
			}

			if _, ok := tt.(xml.EndElement); !ok {
				statestack = append(statestack, state)
			}
		} else if t, ok := tt.(xml.EndElement); ok {
			state := statestack[len(statestack)-1]
			statestack = statestack[0 : len(statestack)-1]
			if !state.inSkip {
				switch t.Name.Local {
				case "s", "utt":
					if !statestack[len(statestack)-1].inS {
						if t := strings.TrimSpace(string(text)); t != "" {
							fmt.Fprintf(fpout, "##PAQULBL %s\n%s\n", hex.EncodeToString([]byte(label)), t)
							text = text[0:0]
						}
					}
				}
			}
		} else if t, ok := tt.(xml.CharData); ok {
			state := statestack[len(statestack)-1]
			if !state.inSkip && (state.inS || state.inUtt) && state.inW && state.inT {
				ww := make([]string, 0, 1)
				for _, w := range strings.Fields(string(t)) {
					ww = append(ww, alpinoEscape(w))
				}
				if len(ww) > 0 {
					text = append(text, []byte(strings.Join(ww, " ")+" ")...)
				}
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
