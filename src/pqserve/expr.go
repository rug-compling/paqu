package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	reDate      = regexp.MustCompile("^[_1-9][_0-9][_0-9][_0-9]-[_0-9][_0-9]-[_0-9][_0-9]$")
	reDateTime1 = regexp.MustCompile("^[_1-9][_0-9][_0-9][_0-9]-[_0-9][_0-9]-[_0-9][_0-9]  ?[_0-9][_0-9]?:[_0-9][_0-9]$")
	reDateTime2 = regexp.MustCompile("^[_1-9][_0-9][_0-9][_0-9]-[_0-9][_0-9]-[_0-9][_0-9]  ?[_0-9][_0-9]?:[_0-9][_0-9]:[_0-9][_0-9]$")
	reQQ        = regexp.MustCompile(`^[^\\ \t\r\n\f"()|&<>=~%]+$`)
)

func sqlmeta(q *Context, prefix string, text string) (query string, njoins int, usererr error, syserr error) {

	////////////////////////////////////////////////////////////////

	// tokenize -> tokens

	var tokens []string

	inWord := false
	inStr := false
	inCmp := false
	inQuote := false
	var tokbuf bytes.Buffer
	for _, c := range strings.TrimSpace(strings.Replace(text, "\r\n", "\n", -1)) {
		if inQuote {
			tokbuf.WriteRune(c)
			inQuote = false
			continue
		}
		if c == '\\' && inStr {
			inQuote = true
			continue
		}
		switch c {
		case ' ', '\t', '\n', '\r', '\f':
			if inStr {
				tokbuf.WriteRune(c)
			} else if inWord || inCmp {
				tokbuf.WriteString("\a")
				inWord = false
				inCmp = false
			}
		case '"':
			// aanhalingstekens bewaren, met name voor onderscheid tussen ( en "("
			if inStr {
				tokbuf.WriteString("\"\a")
			} else if inWord || inCmp {
				tokbuf.WriteString("\a")
				inWord = false
				inCmp = false
			}
			if !inStr {
				tokbuf.WriteString("\"")
			}
			inStr = !inStr
		case '(', ')', '|', '&':
			if inWord || inCmp {
				tokbuf.WriteString("\a")
				inWord = false
				inCmp = false
			}
			tokbuf.WriteRune(c)
			if !inStr {
				tokbuf.WriteString("\a")
			}
		case '<', '>':
			if inWord || inCmp {
				tokbuf.WriteString("\a")
				inWord = false
				inCmp = false
			}
			tokbuf.WriteRune(c)
			if !inStr {
				inCmp = true
			}
		case '=', '~', '%':
			if inWord {
				tokbuf.WriteString("\a")
				inWord = false
			}
			tokbuf.WriteRune(c)
			inCmp = false
			if !inStr {
				tokbuf.WriteString("\a")
			}
		default:
			tokbuf.WriteRune(c)
			if !inStr {
				inWord = true
			}
		}
	}
	if inStr {
		return "", 0, fmt.Errorf("Ongesloten string"), nil
	}

	toks := tokbuf.String()
	for strings.HasSuffix(toks, "\a") {
		toks = toks[:len(toks)-1]
	}
	tokens = strings.Split(toks, "\a")

	////////////////////////////////////////////////////////////////

	// info over metadata inlezen -> ids, types

	ids := make(map[string]int)
	types := make(map[string]string)

	hasmeta := false
	rows, err := sqlDB.Query(fmt.Sprintf("SELECT `hasmeta` FROM `%s_info` WHERE `id` = %q", Cfg.Prefix, prefix))
	if err != nil {
		return "", 0, nil, err
	}
	for rows.Next() {
		err := rows.Scan(&hasmeta)
		if err != nil {
			rows.Close()
			return "", 0, nil, err
		}
	}
	err = rows.Err()
	if err != nil {
		return "", 0, nil, err
	}
	if !hasmeta {
		return "", 0, fmt.Errorf("Corpus %q heeft geen metadata", prefix), nil
	}

	rows, err = sqlDB.Query(fmt.Sprintf("SELECT `id`,`type`,`name` FROM `%s_c_%s_midx`", Cfg.Prefix, prefix))
	if err != nil {
		return "", 0, nil, err
	}
	for rows.Next() {
		var i int
		var t, n string
		err := rows.Scan(&i, &t, &n)
		if err != nil {
			rows.Close()
			return "", 0, nil, err
		}
		n = unHigh(n)
		ids[n] = i
		types[n] = t
	}
	err = rows.Err()
	if err != nil {
		return "", 0, nil, err
	}

	////////////////////////////////////////////////////////////////

	// parse

	lvl := 0
	maxlvl := 0
	for _, token := range tokens {
		if token == "(" {
			lvl++
			if lvl > maxlvl {
				maxlvl = lvl
			}
		} else if token == ")" {
			lvl--
			if lvl < 0 {
				return "", 0, fmt.Errorf("Ongeldig sluithaakje"), nil
			}
		}
	}
	if lvl > 0 {
		return "", 0, fmt.Errorf("Ontbrekend(e) sluithaakje(s)"), nil
	}
	start := make([]int, maxlvl+1)
	current := make([]int, maxlvl+1)
	max := make([]int, maxlvl+1)

	var sqlbuf bytes.Buffer
	var id int
	var typ, cmp, collate string
	var pair [2]string

	state := 1
	for _, token := range tokens {
		switch state {
		case 1:
			if token == "(" {
				sqlbuf.WriteString("(")
				start[lvl+1] = current[lvl]
				current[lvl+1] = current[lvl]
				max[lvl+1] = current[lvl]
				lvl++
			} else {
				if token[0] == '"' {
					token = token[1 : len(token)-1]
				}
				id = ids[token]
				typ = types[token]
				if typ == "" {
					return "", 0, fmt.Errorf("Onbekende metadata: %q", token), nil
				}
				state = 2
			}
		case 2:
			collate = ""
			if token == "in" {
				cmp = token
			} else if (token == "<=" || token == ">=" || token == "<" || token == ">" || token == "=") && typ != "BOOL" {
				cmp = token
				collate = ` COLLATE "utf8_bin"`
			} else if token == "=" && typ == "BOOL" {
				cmp = "="
			} else if token == "%" && (typ == "TEXT" || typ == "DATE" || typ == "DATETIME") {
				cmp = "LIKE"
			} else if token == "~" && typ == "TEXT" {
				cmp = "REGEXP"
			} else {
				return "", 0, fmt.Errorf("Ongeldige operator: %q", token), nil
			}
			if cmp == "in" {
				state = 21
			} else {
				state = 3
			}
		case 3:
			current[lvl]++
			if current[lvl] > max[lvl] {
				max[lvl] = current[lvl]
			}
			n := current[lvl]
			fmt.Fprintf(&sqlbuf, "`meta%d`.`id` = %d AND ", n, id)
			if token == "nil" && cmp == "=" {
				fmt.Fprintf(&sqlbuf, "`meta%d`.`idx` = 2147483647", n)
			} else {
				fmt.Fprintf(&sqlbuf, "`meta%d`.`idx` != 2147483647 AND ", n)
				if token[0] == '"' {
					token = token[1 : len(token)-1]
				}
				if cmp == "REGEXP" {
					_, err := regexp.Compile(token)
					if err != nil {
						return "", 0, err, nil
					}
				}
				switch typ {
				case "TEXT":
					fmt.Fprintf(&sqlbuf, "`meta%d`.`tval` %s %q%s", n, cmp, token, collate)
				case "BOOL":
					i := 0
					if isTrue[token] {
						i = 1
					}
					fmt.Fprintf(&sqlbuf, "`meta%d`.`ival` %s %d", n, cmp, i)
				case "INT":
					i, err := strconv.Atoi(token)
					if err != nil {
						return "", 0, err, nil
					}
					fmt.Fprintf(&sqlbuf, "`meta%d`.`ival` %s %d", n, cmp, i)
				case "FLOAT":
					f, err := strconv.ParseFloat(token, 64)
					if err != nil {
						return "", 0, err, nil
					}
					fmt.Fprintf(&sqlbuf, "`meta%d`.`fval` %s %g", n, cmp, float32(f))
				case "DATE":
					if !reDate.MatchString(token) {
						return "", 0, fmt.Errorf("Ongeldige datum: %s", token), nil
					}
					fmt.Fprintf(&sqlbuf, "`meta%d`.`dval` %s \"%s 00:00:00\"", n, cmp, token)
				case "DATETIME":
					if reDateTime1.MatchString(token) {
						token = token + ":00"
					} else if !reDateTime2.MatchString(token) {
						return "", 0, fmt.Errorf("Ongeldige datum en tijd: %s", token), nil
					}
					fmt.Fprintf(&sqlbuf, "`meta%d`.`dval` %s \"%s\"", n, cmp, token)
				}
			}
			state = 4
		case 4:
			if token == ")" {
				sqlbuf.WriteString(")")
				lvl--
				current[lvl] = max[lvl+1]
				if current[lvl] > max[lvl] {
					max[lvl] = current[lvl]
				}
			} else {
				t := strings.ToLower(token)
				if t == "&" || t == "en" || t == "and" {
					sqlbuf.WriteString("\nAND ")
				} else if t == "|" || t == "of" || t == "or" {
					sqlbuf.WriteString("\nOR ")
					current[lvl] = start[lvl]
				} else {
					return "", 0, fmt.Errorf("Ongeldige operator %q (verwacht: AND, OR)", token), nil
				}
				state = 1
			}
		case 21:
			pair[0] = token
			state = 22
		case 22:
			pair[1] = token
			current[lvl]++
			if current[lvl] > max[lvl] {
				max[lvl] = current[lvl]
			}
			n := current[lvl]
			fmt.Fprintf(&sqlbuf, "`meta%d`.`id` = %d AND `meta%d`.`idx` != 2147483647 AND ", n, id, n)
			var f string
			for i := 0; i < 2; i++ {
				if pair[i][0] == '"' {
					pair[i] = pair[i][1 : len(pair[i])-1]
				}
				switch typ {
				case "TEXT":
					pair[i] = fmt.Sprintf("%q", pair[i])
					f = "tval"
				case "BOOL":
					f = "ival"
				case "INT":
					_, err := strconv.Atoi(pair[i])
					if err != nil {
						return "", 0, err, nil
					}
					f = "ival"
				case "FLOAT":
					_, err := strconv.ParseFloat(pair[i], 64)
					if err != nil {
						return "", 0, err, nil
					}
					f = "fval"
				case "DATE":
					if !reDate.MatchString(pair[i]) {
						return "", 0, fmt.Errorf("Ongeldige datum: %s", token), nil
					}
					pair[i] = `"` + pair[i] + `"`
					f = "dval"
				case "DATETIME":
					if reDateTime1.MatchString(pair[i]) {
						pair[i] = pair[i] + ":00"
					} else if !reDateTime2.MatchString(pair[i]) {
						return "", 0, fmt.Errorf("Ongeldige datum en tijd: %s", pair[i]), nil
					}
					pair[i] = `"` + pair[i] + `"`
					f = "dval"
				}
			}
			fmt.Fprintf(&sqlbuf, "`meta%d`.`%s` BETWEEN %s AND %s", n, f, pair[0], pair[1])
			state = 4
		}
	}

	if state != 4 {
		return "", 0, fmt.Errorf("Onvolledige parse"), nil
	}

	return sqlbuf.String(), max[0], nil, nil
}

func qquote(s string) string {
	if reQQ.MatchString(s) {
		return s
	}
	return fmt.Sprintf("%q", s)
}
