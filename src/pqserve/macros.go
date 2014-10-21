package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

type MacroResult struct {
	Err    string   `json:"err"`
	Macros string   `json:"macros"`
	Keys   []string `json:"keys"`
}

type Macros struct {
	rules map[string]string
	keys  []string
	time  time.Time
}

var (
	macroRE = regexp.MustCompile(`([a-zA-Z][_a-zA-Z0-9]*)\s*=\s*"""((?s:.*?))"""`)
	macroKY = regexp.MustCompile(`%[a-zA-Z][_a-zA-Z0-9]*%`)

	macroLock sync.Mutex
	macroMap  = make(map[string]Macros)
)

func savemacros(q *Context) {
	if !has_dbxml {
		http.NotFound(q.w, q.r)
		return
	}

	macros := firstf(q.form, "macrotext")

	result := MacroResult{Keys: make([]string, 0)}

MACROLOOP:
	for {
		if !q.auth {
			result.Err = "Je bent niet ingelogd"
			break
		}

		if len(macros) > 65535 {
			result.Err = "Te groot (max is 65535 tekens)"
			break
		}

		s := macroRE.ReplaceAllLiteralString(macros, "")
		if t := strings.TrimSpace(s); t != "" {
			result.Err = fmt.Sprintf("Overgebleven tekst %q", t)
			break
		}

		rules := make(map[string]string)
		for _, set := range macroRE.FindAllStringSubmatch(macros, -1) {
			rules[set[1]] = set[2]
		}
		for key := range rules {
			for i := 0; i < 101; i++ {
				if i == 100 || len(rules[key]) > 65535 {
					result.Err = fmt.Sprintf("Te diepe recursie voor %q", key)
					break MACROLOOP
				}
				rule := macroKY.ReplaceAllStringFunc(rules[key], func(s string) string {
					return rules[s[1:len(s)-1]]
				})
				if rule == rules[key] {
					break
				}
				rules[key] = rule
			}
		}

		_, err := q.db.Exec(fmt.Sprintf("DELETE FROM `%s_macros` WHERE `user` = %q", Cfg.Prefix, q.user))
		if err != nil {
			result.Err = "Databasefout: " + err.Error()
			logerr(err)
			break
		}

		_, err = q.db.Exec(fmt.Sprintf("INSERT INTO `%s_macros` (`user`, `macros`) VALUES (%q, %q)", Cfg.Prefix, q.user, macros))
		if err != nil {
			result.Err = "Databasefout: " + err.Error()
			logerr(err)
			break
		}

		result.Macros = macros
		for _, set := range macroRE.FindAllStringSubmatch(macros, -1) {
			result.Keys = append(result.Keys, set[1])
		}
		sort.Strings(result.Keys)

		macroLock.Lock()
		delete(macroMap, q.user)
		macroLock.Unlock()
		break

	}

	b, _ := json.Marshal(result)

	fmt.Fprint(q.w, `<!DOCTYPE html>
<html>
<head>
<title></title>
<script type="text/javascript"><!--
window.parent._fn.update(`)
	fmt.Fprintln(q.w, string(b))
	fmt.Fprint(q.w, `);
//--></script>
</head>
<body></body>
</html>
`)

}

func clearMacros() {
	for {
		time.Sleep(time.Hour)

		macroLock.Lock()
		then := time.Now().Add(-1 * time.Hour)
		for user := range macroMap {
			if then.After(macroMap[user].time) {
				delete(macroMap, user)
			}
		}
		macroLock.Unlock()
	}
}

func loadMacros(q *Context) {
	macroLock.Lock()
	defer macroLock.Unlock()
	if m, ok := macroMap[q.user]; ok {
		m.time = time.Now()
		return
	}

	text := ""
	rows, err := q.db.Query(fmt.Sprintf("SELECT `macros` FROM `%s_macros` WHERE `user` = %q", Cfg.Prefix, q.user))
	if err == nil && rows.Next() {
		rows.Scan(&text)
		rows.Close()
	}

	macros := Macros{
		rules: make(map[string]string),
		time:  time.Now(),
	}

	for _, set := range macroRE.FindAllStringSubmatch(text, -1) {
		macros.rules[set[1]] = set[2]
	}

	for key := range macros.rules {
		for {
			rule := macroKY.ReplaceAllStringFunc(macros.rules[key], func(s string) string {
				return macros.rules[s[1:len(s)-1]]
			})
			if rule == macros.rules[key] {
				break
			}
			macros.rules[key] = rule
		}
	}

	macros.keys = make([]string, 0, len(macros.rules))
	for key := range macros.rules {
		macros.keys = append(macros.keys, key)
	}
	sort.Strings(macros.keys)

	macroMap[q.user] = macros

}

func getMacrosKeys(q *Context) []string {
	loadMacros(q)
	macroLock.Lock()
	defer macroLock.Unlock()
	if m, ok := macroMap[q.user]; ok {
		return m.keys
	}
	return []string{}
}

func getMacrosRules(q *Context) map[string]string {
	loadMacros(q)
	macroLock.Lock()
	defer macroLock.Unlock()
	if m, ok := macroMap[q.user]; ok {
		return m.rules
	}
	return map[string]string{}
}
