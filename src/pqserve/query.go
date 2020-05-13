package main

//. Imports

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// sys errors moeten gelogd worden
// user errors hoeven alleen aan gebruiker gemeld te worden (de gebruiker deed iets fout)

// Opstellen van query voor hoofdformulier, en voor stats
func makeQuery(q *Context, prefix, table string, chClose <-chan bool) (string, string, error /*user*/, error /*sys*/) {
	return makeQueryDo(q, prefix, table, chClose, false)
}

func makeQueryF(q *Context, prefix, table string, chClose <-chan bool) (string, string, error /*user*/, error /*sys*/) {
	return makeQueryDo(q, prefix, table, chClose, true)
}

func makeQueryDo(q *Context, prefix, table string, chClose <-chan bool, form bool) (string, string, error /*user*/, error /*sys*/) {

	if table != "" {
		table = "`" + table + "`."
	}

	var frst func(s string) string
	if form {
		frst = func(s string) string {
			return firstf(q.form, s)
		}
	} else {
		frst = func(s string) string {
			return first(q.r, s)
		}
	}

	parts := make([]string, 0, 6)
	for _, p := range []string{"", "h"} {
		if frst(p+"word") != "" {
			wrd := setHigh(frst(p + "word"))
			if wrd[0] == '+' {
				parts = append(parts, fmt.Sprintf(table+"`"+p+"lemma` = %q", wrd[1:]))
			} else if wrd[0] == '@' {
				parts = append(parts, fmt.Sprintf(table+"`"+p+"root` = %q", wrd[1:]))
			} else if wrd[0] == '=' {
				parts = append(parts, fmt.Sprintf(table+"`"+p+"word` = %q COLLATE \"utf8_bin\"", wrd[1:]))
			} else if wrd[0] == '?' {
				parts = append(parts, fmt.Sprintf(table+"`"+p+"word` = %q", wrd[1:]))
			} else if strings.Index(wrd, "%") >= 0 {
				parts = append(parts, fmt.Sprintf(table+"`"+p+"word` LIKE %q", wrd))
			} else {
				var s string
				select {
				case <-chClose:
					return "", "", nil, errConnectionClosed
				default:
				}
				rows, err := sqlDB.Query(fmt.Sprintf("SELECT `lemma` FROM `%s_c_%s_word` WHERE `word` = %q",
					Cfg.Prefix, prefix, wrd))
				if err != nil {
					return "", "", nil, err
				}
				lset := make(map[string]bool)
				for rows.Next() {
					err := rows.Scan(&s)
					if err != nil {
						return "", "", nil, err
					}
					for _, i := range strings.Split(s, "\t") {
						lset[fmt.Sprintf("%q", i)] = true
					}
				}
				err = rows.Err()
				if err != nil {
					return "", "", nil, err
				}
				if len(lset) > 0 {
					ll := make([]string, 0, len(lset))
					for key := range lset {
						ll = append(ll, key)
					}
					if len(ll) == 1 {
						parts = append(parts, table+"`"+p+"lemma` = "+ll[0])
					} else {
						parts = append(parts, table+"`"+p+"lemma` IN ("+strings.Join(ll, ", ")+")")
					}
				} else {
					parts = append(parts, fmt.Sprintf(table+"`"+p+"word` = %q", wrd))
				}
			}
		}
	}
	if s := frst("postag"); s != "" {
		if s == "(leeg)" {
			s = ""
		}
		parts = append(parts, fmt.Sprintf(table+"`postag` = %q", s))
	}
	if s := frst("rel"); s != "" {
		parts = append(parts, fmt.Sprintf(table+"`rel` = %q", s))
	}
	if s := frst("hpostag"); s != "" {
		if s == "(leeg)" {
			s = ""
		}
		parts = append(parts, fmt.Sprintf(table+"`hpostag` = %q", s))
	}

	query := strings.Join(parts, " AND ")

	joins := make([]string, 0)
	if m := frst("meta"); m != "" {
		meta, n, usererr, syserr := sqlmeta(q, prefix, m)
		if usererr != nil || syserr != nil {
			return "", "", usererr, syserr
		}
		for i := 1; i <= n; i++ {
			if table == "" {
				joins = append(joins, fmt.Sprintf("JOIN `%s_c_%s_meta` `meta%d` USING(`arch`,`file`)", Cfg.Prefix, prefix, i))
			} else {
				joins = append(joins, fmt.Sprintf("JOIN `%s_c_%s_meta` `meta%d` ON (`meta%d`.`arch` = %s`arch` AND `meta%d`.`file`=%s`file`)",
					Cfg.Prefix, prefix, i, i, table, i, table))
			}
		}
		if query != "" {
			query = query + " AND ( " + meta + " )"
		} else {
			query = meta
		}
	}

	return query, strings.Join(joins, " "), nil, nil
}

// Query uitvoeren. Cancel query niet alleen bij timeout, maar ook als request wordt verbroken.
func timeoutQuery(q *Context, chClose <-chan bool, query string) (*sql.Rows, error) {

	chLog <- "QUERY: " + query
	a := make([]byte, 16)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 16; i++ {
		a[i] = byte(97 + rnd.Intn(24))
	}
	id := "/* " + string(a) + " */"
	query = strings.Replace(query, " ", " "+id+" ", 1)

	timeout := true
	if Cfg.Querytimeout > 0 && hasMaxExecutionTime {
		t := fmt.Sprintf(" /*+ MAX_EXECUTION_TIME(%d000) */ ", Cfg.Querytimeout)
		query = strings.Replace(query, " ", t, 1)
		timeout = false // laat timeout door MySQL-server doen
	}
	if Cfg.Querytimeout > 0 && hasMaxStatementTime {
		query = fmt.Sprintf("SET STATEMENT max_statement_time=%d FOR %s", Cfg.Querytimeout, query)
		timeout = false // laat timeout door MariaDB-server doen
	}

	chFinished := make(chan bool)
	defer close(chFinished)
	go cancelQuery(id, timeout, chFinished, chClose)

	return sqlDB.Query(query)
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

	// deze query geeft ook zichzelf als resultaat
	rows, err := sqlDB.Query("SELECT `ID` FROM `information_schema`.`PROCESSLIST` WHERE `INFO` LIKE \"%" + id + "%\"")
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
		sqlDB.Exec("KILL QUERY " + id)
	}
}
