// +build spodmake

package main

import (
	"github.com/pebbe/util"

	"bytes"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type Line struct {
	lvl int
	s   string
}

var (
	macroRE  = regexp.MustCompile(`([a-zA-Z][_a-zA-Z0-9]*)\s*=\s*"""((?s:.*?))"""`)
	macroKY  = regexp.MustCompile(`%[a-zA-Z][_a-zA-Z0-9]*%`)
	macroCOM = regexp.MustCompile(`(?m:^\s*#.*)`)

	reComment = regexp.MustCompile("(?s:\\(:.*?:\\))")
	reNumber  = regexp.MustCompile("number\\(([^()]*)\\)")
	reRemove  = regexp.MustCompile("[^a-zA-Z0-9]")
	reSet     = regexp.MustCompile("(?s:=\\s*\\(\\s*(\".*?\")\\s*\\))")

	hash = make(map[string]string)

	x = util.CheckErr
)

func main() {

	fmt.Print(`// GENERATED FILE. DO NOT EDIT.

package main

var spod2xpath = map[string]*xPath{
`)

	macros := make(map[string]string)

	b, err := ioutil.ReadFile(os.Args[1])
	x(err)

	for _, set := range macroRE.FindAllStringSubmatch(macroCOM.ReplaceAllLiteralString(string(b), ""), -1) {
		s := strings.Replace(set[2], "\r\n", "\n", -1)
		s = strings.Replace(s, "\n\r", "\n", -1)
		s = strings.Replace(s, "\r", "\n", -1)
		macros[set[1]] = s
	}

	for key := range macros {
		for {
			rule := macroKY.ReplaceAllStringFunc(macros[key], func(s string) string {
				return macros[s[1:len(s)-1]]
			})
			if rule == macros[key] {
				break
			}
			if len(rule) > 100000 {
				x(fmt.Errorf("RECURSIONLIMIT"))
			}
			macros[key] = rule
		}
	}

	for _, spod := range spods {
		if spod.special == "hidden1" || spod.special == "attr" || spod.special == "parser" {
			continue
		}

		if spod.method != SPOD_STD {
			x(fmt.Errorf("Method not supported in %s", spod.lbl))
		}
		query := macroKY.ReplaceAllStringFunc(spod.xpath, func(s string) string {
			return macros[s[1:len(s)-1]]
		})

		s2 := reComment.ReplaceAllLiteralString(query, "")

		// 'string' -> "string"
		s2 = strings.Replace(s2, "'", `"`, -1)

		// [1] -> [first()]
		s2 = strings.Replace(s2, "[1]", "[first()]", -1)

		// [3] -> [third()]
		s2 = strings.Replace(s2, "[3]", "[third()]", -1)

		// node[3][bla] -> node[3]/self::node[bla]
		s2 = strings.Replace(s2, "][", "]/self::node[", -1)

		// =("aap","noot","mies")  ->  ="hashcode"
		s2 = reSet.ReplaceAllStringFunc(s2, func(s string) string {
			mm := reSet.FindStringSubmatch(s)
			h := fmt.Sprintf("%x", md5.Sum([]byte(mm[1])))
			hash[h] = mm[1]
			return `="` + h + `"`
		})

		// number(..) => ..
		s2 = reNumber.ReplaceAllString(s2, " $1 ")

		cmd := exec.Command("testXPath", "--tree", s2)
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		x(cmd.Run())
		so := stdout.String()
		se := strings.Replace(stderr.String(), "xmlXPathCompiledEval: evaluation failed\n", "", 1)
		se = strings.Replace(se, "XPath error : Undefined variable\n", "", 1)

		if se != "" {
			x(fmt.Errorf("%s\n\n%s\n", se, s2))
		}

		fmt.Printf("%q: %s,\n", spod.lbl, parse(so))

	}

	fmt.Println("}")
}

func parse(s string) string {

	out := make([]string, 0)

	lines := make([]Line, 0)
	for _, l := range strings.Split(s, "\n") {
		if len(l) == 0 || l[0] != ' ' {
			continue
		}
		lvl := (len(l) - len(strings.TrimLeft(l, " "))) / 2
		lines = append(lines, Line{
			s:   strings.TrimSpace(l),
			lvl: lvl,
		})
	}

	if len(lines) == 0 {
		x(fmt.Errorf("COMPILE ERROR"))
	}

	out = append(out, "&xPath{")
	lvl := 0
	prev := make([]int, 1)
	for _, line := range lines {
		for lvl > line.lvl {
			out = append(out, "},")
			lvl--
			prev = prev[:len(prev)-1]
		}
		if lvl == line.lvl {
			out = append(out, "},")
		} else {
			lvl = line.lvl
		}
		words := strings.Fields(line.s)
		cmd := strings.ToLower(words[0])
		str := "d" + strings.Title(cmd)
		cmd0 := cmd
		if len(prev) <= lvl {
			prev = append(prev, 1)
			cmd = "arg1"
		} else {
			prev[lvl] += 1
			cmd = fmt.Sprint("arg", prev[lvl])
		}
		out = append(out, fmt.Sprintf("%s: &%s{", cmd, str))
		if cmd0 == "variable" {
			out = append(out, fmt.Sprintf("VAR: %s,", words[1]))
		} else if cmd0 == "elem" {
			if words[4] == "string" {
				data, ok := hash[words[6]]
				if !ok {
					data = `"` + words[6] + `"`
				}
				out = append(out, fmt.Sprintf("DATA: []interface{}{%s},", data))
			} else if words[4] == "number" {
				out = append(out, fmt.Sprintf("DATA: []interface{}{%s},", words[6]))
			} else {
				panic("Unknown elem type: " + strings.Join(words, " "))
			}
		} else {
			if len(words) > 1 {
				args := strings.Join(words[1:], " ")
				args = strings.Replace(args, "'", "", -1)
				args = strings.Replace(args, " :", "", -1)
				args = strings.Replace(args, "<", "lt", -1)
				args = strings.Replace(args, ">", "gt", -1)
				args = strings.Replace(args, "(", " ", -1)
				args = strings.Replace(args, ")", "", -1)
				args = strings.Replace(args, "=", "is", -1)
				args = strings.Replace(args, "--", "minmin", -1)
				args = strings.Replace(args, "deep-equal", "deep equal", -1)
				args = strings.Replace(args, "-with", " with", -1)
				args = strings.Replace(args, "-or-self", " or self", -1)
				args = strings.Replace(args, "+", "plus", -1)
				args = strings.Replace(args, " name node", "", -1)
				args = strings.Replace(args, "local:internal_head_position", "local internal head position", -1)
				args = strings.Replace(reRemove.ReplaceAllStringFunc(args, func(s string) string {
					return fmt.Sprintf("_%02x", s[0])
				}), "_20", "__", -1)
				if args == "_2d" {
					args = "minus"
				}
				out = append(out, fmt.Sprintf("ARG: %s__%s,", cmd0, args))
			}
		}
	}
	for lvl > 0 {
		out = append(out, "},")
		lvl--
	}
	out = append(out, "}")

	return strings.Join(out, "\n")
}
