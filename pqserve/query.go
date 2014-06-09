package main

//. Imports

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Opstellen van query voor hoofdformulier, en voor stats
func makeQuery(q *Context, prefix string, chClose <-chan bool) (string, error) {

	parts := make([]string, 0, 6)
	for _, p := range []string{"", "h"} {
		if first(q.r, p+"word") != "" {
			wrd := first(q.r, p+"word")
			if wrd[0] == '+' {
				parts = append(parts, fmt.Sprintf("`"+p+"lemma` = %q", wrd[1:]))
			} else if wrd[0] == '@' {
				parts = append(parts, fmt.Sprintf("`"+p+"root` = %q", wrd[1:]))
			} else if wrd[0] == '=' {
				parts = append(parts, fmt.Sprintf("`"+p+"word` = %q COLLATE \"utf8_bin\"", wrd[1:]))
			} else if wrd[0] == '?' {
				parts = append(parts, fmt.Sprintf("`"+p+"word` = %q", wrd[1:]))
			} else if strings.Index(wrd, "%") >= 0 {
				parts = append(parts, fmt.Sprintf("`"+p+"word` LIKE %q", wrd))
			} else {
				var s string
				select {
				case <-chClose:
					return "", ConnectionClosed
				default:
				}
				rows, err := q.db.Query(fmt.Sprintf("SELECT `lemma` FROM `%s_c_%s_word` WHERE `word` = %q",
					Cfg.Prefix, prefix, wrd))
				if err != nil {
					return "", err
				}
				lset := make(map[string]bool)
				for rows.Next() {
					err := rows.Scan(&s)
					if err != nil {
						return "", err
					}
					for _, i := range strings.Split(s, "\t") {
						lset[fmt.Sprintf("%q", i)] = true
					}
				}
				err = rows.Err()
				if err != nil {
					return "", err
				}
				if len(lset) > 0 {
					ll := make([]string, 0, len(lset))
					for key := range lset {
						ll = append(ll, key)
					}
					if len(ll) == 1 {
						parts = append(parts, "`"+p+"lemma` = "+ll[0])
					} else {
						parts = append(parts, "`"+p+"lemma` IN ("+strings.Join(ll, ", ")+")")
					}
				} else {
					parts = append(parts, fmt.Sprintf("`"+p+"word` = %q", wrd))
				}
			}
		}
	}
	if s := first(q.r, "postag"); s != "" {
		parts = append(parts, fmt.Sprintf("`postag` = %q", s))
	}
	if s := first(q.r, "rel"); s != "" {
		parts = append(parts, fmt.Sprintf("`rel` = %q", s))
	}
	if s := first(q.r, "hpostag"); s != "" {
		parts = append(parts, fmt.Sprintf("`hpostag` = %q", s))
	}

	return strings.Join(parts, " AND "), nil
}

// Query uitvoeren. Cancel query niet alleen bij timeout, maar ook als request wordt verbroken.
func timeoutQuery(q *Context, chClose <-chan bool, query string) (*sql.Rows, error) {

	a := make([]byte, 16)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 16; i++ {
		a[i] = byte(97 + rnd.Intn(24))
	}
	id := "/* " + string(a) + " */"
	query = strings.Replace(query, " ", " "+id+" ", 1)

	timeout := true
	if Cfg.Querytimeout > 0 && hasMaxStatementTime {
		t := fmt.Sprintf(" MAX_STATEMENT_TIME = %d000 ", Cfg.Querytimeout)
		query = strings.Replace(query, " ", t, 1)
		timeout = false // laat timeout door MySQL-server doen
	}

	chFinished := make(chan bool)
	defer close(chFinished)
	go cancelQuery(id, timeout, chFinished, chClose)

	return q.db.Query(query)
}

// hulpfunctie voor timeoutQuery
func cancelQuery(id string, timeout bool, chFinished chan bool, chClose <-chan bool) {

	if timeout && Cfg.Querytimeout > 0 {
		select {
		case <-chFinished:
			// query is klaar
			return
		case <-chClose:
			// verbinding is gesloten: cancel query
		case <-time.After(time.Duration(Cfg.Querytimeout) * time.Second):
			// timeout: cancel query
		}
	} else {
		select {
		case <-chFinished:
			// query is klaar
			return
		case <-chClose:
			// verbinding is gesloten: cancel query
		}
	}

	db, err := dbopen()
	if err != nil {
		logerr(err)
		return
	}
	defer db.Close()

	// deze query geeft ook zichzelf als resultaat
	rows, err := db.Query("SELECT `ID` FROM `information_schema`.`PROCESSLIST` WHERE `INFO` LIKE \"%" + id + "%\"")
	if err != nil {
		logerr(err)
		return
	}
	ids := make([]string, 0, 2)
	for rows.Next() {
		var s string
		err := rows.Scan(&s)
		if err != nil {
			logerr(err)
			return
		}
		ids = append(ids, s)
	}
	for _, id := range ids {
		db.Exec("KILL QUERY " + id)
	}
}
