package spod

import (
	"github.com/rug-compling/paqu/internal/file"

	"crypto/md5"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	macroRE  = regexp.MustCompile(`([a-zA-Z][_a-zA-Z0-9]*)\s*=\s*"""((?s:.*?))"""`)
	macroKY  = regexp.MustCompile(`%[a-zA-Z][_a-zA-Z0-9]*%`)
	macroCOM = regexp.MustCompile(`(?m:^\s*#.*)`)
	rules    = make(map[string]string)
	index    = make(map[string]int)
	hash     = make(map[string]string)
	has      = make(map[string]int)
)

func init() {
	loadMacros()
	for idx, spod := range Spods {
		index[spod.Lbl] = idx
		if strings.HasPrefix(spod.Lbl, "has_") {
			has[spod.Lbl[4:]] = idx
		}
		query := macroKY.ReplaceAllStringFunc(spod.Xpath, func(s string) string {
			return rules[s[1:len(s)-1]]
		})
		query = strings.Join(strings.Fields(query), " ")
		hash[spod.Lbl] = fmt.Sprintf("%x", md5.Sum([]byte(query+spod.Method)))
	}
}

func loadMacros() {

	for _, set := range macroRE.FindAllStringSubmatch(macroCOM.ReplaceAllLiteralString(file.File__macros__txt, ""), -1) {
		s := strings.Replace(set[2], "\r\n", "\n", -1)
		s = strings.Replace(s, "\n\r", "\n", -1)
		s = strings.Replace(s, "\r", "\n", -1)
		rules[set[1]] = s
	}

	for key := range rules {
		for {
			rule := macroKY.ReplaceAllStringFunc(rules[key], func(s string) string {
				return rules[s[1:len(s)-1]]
			})
			if rule == rules[key] {
				break
			}
			if len(rule) > 100000 {
				rules[key] = "RECURSIONLIMIT"
				break
			}
			rules[key] = rule
		}
	}

}

func Hash(label string) string {
	h, ok := hash[label]
	if !ok {
		panic("Undefined label " + label)
	}
	return h
}

func Has(label string, dir string) (bool, error) {
	if label == "pos" || label == "postag" {
		return false, nil
	}
	idx, ok := index[label]
	if !ok {
		panic("Undefined label " + label)
	}
	spod := Spods[idx]

	if spod.Special == "" {
		return true, nil
	}

	for _, special := range strings.Fields(spod.Special) {
		neg := false
		if special[0] == '-' {
			neg = true
			special = special[1:]
		}
		pri, ok := has[special]
		if !ok {
			return true, nil
		}
		b, err := ioutil.ReadFile(filepath.Join(dir, Hash(Spods[pri].Lbl)))
		if err != nil {
			return false, err
		}
		fields := strings.Fields(string(b))
		if len(fields) < 2 {
			return false, fmt.Errorf("No field")
		}
		if fields[1] == "0" == neg {
			return true, nil
		}
	}
	return false, nil
}
