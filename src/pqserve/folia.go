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
	var label, wid string
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
			for _, e := range t.Attr {
				switch e.Name.Local {
				case "id":
					id = e.Value
				case "class":
					hasClass = true
				case "auth":
					state.inSkip = true
				}
			}

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
				wid = id
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
				case "w":
					text = append(text, ' ')
				case "s", "utt":
					if !statestack[len(statestack)-1].inS && strings.TrimSpace(string(text)) != "" {
						words := make([]string, 0)
						for _, w := range strings.Fields(string(text)) {
							words = append(words, alpinoEscape(w))
						}
						fmt.Fprintf(fpout, "##PAQULBL %s\n%s\n", hex.EncodeToString([]byte(label)), strings.Join(words, " "))
						text = text[0:0]
					}
				}
			}
		} else if t, ok := tt.(xml.CharData); ok {
			state := statestack[len(statestack)-1]
			if !state.inSkip && (state.inS || state.inUtt) && state.inW && state.inT {
				s := alpinoEscape(string(t))
				if Cfg.Alpino15 && wid != "" {
					s = fmt.Sprintf("[ @id %s ] %s", alpinoEscape(wid), s)
				}
				text = append(text, []byte(s)...)
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
