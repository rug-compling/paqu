package main

//. Imports

import (
	"database/sql"
	"fmt"
	"html"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

//. Main

func stats(q *Context) {

	var chClose <-chan bool
	if f, ok := q.w.(http.CloseNotifier); ok {
		chClose = f.CloseNotify()
	} else {
		chClose = make(<-chan bool)
	}

	now := time.Now()

	download := false
	if first(q.r, "d") != "" {
		download = true
	}

	option := make(map[string]string)
	for _, t := range []string{"word", "postag", "rel", "hpostag", "hword"} {
		option[t] = first(q.r, t)
	}
	if option["word"] == "" && option["postag"] == "" && option["rel"] == "" && option["hpostag"] == "" && option["hword"] == "" {
		http.Error(q.w, "Query ontbreekt", http.StatusPreconditionFailed)
		return
	}

	prefix := first(q.r, "db")
	if prefix == "" {
		http.Error(q.w, "Geen corpus opgegeven", http.StatusPreconditionFailed)
		return
	}
	if !q.prefixes[prefix] {
		http.Error(q.w, "Ongeldig corpus", http.StatusPreconditionFailed)
		return
	}

	// BEGIN: Opstellen van de query
	// Deze code moet gelijk zijn aan die in het programma 'lassystats'
	parts := make([]string, 0, 6)
	for _, p := range []string{"", "h"} {
		if option[p+"word"] != "" {
			wrd := option[p+"word"]
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
					logerr(ConnectionClosed)
					return
				default:
				}
				rows, err := q.db.Query(fmt.Sprintf("SELECT `lemma` FROM `%s_c_%s_word` WHERE `word` = %q",
					Cfg.Prefix, prefix, wrd))
				if err != nil {
					http.Error(q.w, err.Error(), http.StatusInternalServerError)
					logerr(err)
					return
				}
				lset := make(map[string]bool)
				for rows.Next() {
					err := rows.Scan(&s)
					if err != nil {
						http.Error(q.w, err.Error(), http.StatusInternalServerError)
						logerr(err)
						return
					}
					for _, i := range strings.Split(s, "\t") {
						lset[fmt.Sprintf("%q", i)] = true
					}
				}
				err = rows.Err()
				if err != nil {
					http.Error(q.w, err.Error(), http.StatusInternalServerError)
					logerr(err)
					return
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
	if option["postag"] != "" {
		parts = append(parts, fmt.Sprintf("`postag` = %q", option["postag"]))
	}
	if option["rel"] != "" {
		parts = append(parts, fmt.Sprintf("`rel` = %q", option["rel"]))
	}
	if option["hpostag"] != "" {
		parts = append(parts, fmt.Sprintf("`hpostag` = %q", option["hpostag"]))
	}
	query := strings.Join(parts, " AND ")
	// EINDE: Opstellen van de query

	// BEGIN UITVOER

	if download {
		q.w.Header().Set("Content-type", "text/plain; charset=utf-8")
	} else {
		q.w.Header().Set("Content-type", "text/html; charset=utf-8")
	}
	q.w.Header().Set("Cache-Control", "no-cache")
	q.w.Header().Add("Pragma", "no-cache")

	// DEBUG: HTML-uitvoer van de query
	if !download {
		fmt.Fprint(q.w, "<div style=\"font-family:monospace\">\n", html.EscapeString(query), "\n</div><p>\n")
	}

	// Aantal zinnen die matchen met de query
	select {
	case <-chClose:
		logerr(ConnectionClosed)
		return
	default:
	}
	rows, err := timeoutQuery(q, chClose, "SELECT 1 FROM `"+Cfg.Prefix+"_c_"+prefix+"_deprel` WHERE "+
		query+" GROUP BY `arch`,`file`")
	if err != nil {
		interneFoutRegel(q, err, !download)
		logerr(err)
		return
	}
	counter := 0
	for rows.Next() {
		counter++
	}
	if err != nil {
		interneFoutRegel(q, err, !download)
		logerr(err)
		return
	}

	if download {
		fmt.Fprintf(q.w, "# %d zinnen\t\n", counter)
	} else {
		fmt.Fprintln(q.w, "Aantal gevonden zinnen:", iformat(counter))
	}

	// Tellingen van onderdelen
	for i, ww := range []string{"word", "lemma", "postag", "rel", "hword", "hlemma", "hpostag"} {
		var j, count int
		var s, p, limit string
		if download {
			fmt.Fprintln(q.w, "# "+ww+"\t")
		} else {
			if i == 0 {
				fmt.Fprintln(q.w, "<p>"+YELLOW+"<b>word</b></span>: ")
			} else if i == 4 {
				fmt.Fprintln(q.w, "<p>"+GREEN+"<b>hword</b></span>: ")
			} else {
				fmt.Fprintln(q.w, "<p><b>"+ww+"</b>: ")
			}
			limit = " LIMIT " + fmt.Sprint(WRDMAX)
		}
		select {
		case <-chClose:
			logerr(ConnectionClosed)
			return
		default:
		}
		rows, err := timeoutQuery(q, chClose, "SELECT count(*), `"+ww+"` FROM `"+Cfg.Prefix+"_c_"+prefix+
			"_deprel` WHERE "+query+" GROUP BY `"+ww+"` COLLATE 'utf8_bin' ORDER BY 1 DESC, 2"+limit)
		if err != nil {
			interneFoutRegel(q, err, !download)
			logerr(err)
			return
		}
		for rows.Next() {
			err := rows.Scan(&j, &s)
			if err != nil {
				interneFoutRegel(q, err, !download)
				logerr(err)
				return
			}
			if s == "" {
				s = "\"\""
			}
			if download {
				fmt.Fprintf(q.w, "%d\t%s\n", j, s)
			} else {
				fmt.Fprint(q.w, p, j, "&times;&nbsp;", html.EscapeString(s))
				p = ", "
				count++
			}
		}
		err = rows.Err()
		if err != nil {
			interneFoutRegel(q, err, !download)
			logerr(err)
			return
		}
		if count == WRDMAX {
			fmt.Fprint(q.w, ", ...")
		}
		if !download {
			fmt.Fprint(q.w, "\n<BR>\n")
		}
	}

	if !download {
		fmt.Fprintf(q.w,
			"<hr>tijd: %s\n<p>\n<a href=\"/stats?%s&amp;d=1\" target=\"_blank\">download</a>\n",
			time.Now().Sub(now),
			q.r.URL.RawQuery)
	}
}

func interneFoutRegel(q *Context, err error, is_html bool) {
	s := err.Error()
	if is_html {
		s = html.EscapeString(s)
	}
	fmt.Fprintln(q.w, "Interne fout:", s)
}

// cancel query niet alleen bij timeout, maar ook als request wordt verbroken
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
		t := fmt.Sprintf(" MAX_STATEMENT_TIME = %n000 ", Cfg.Querytimeout)
		query = strings.Replace(query, " ", t, 1)
		timeout = false // laat timeout door MySQL-server doen
	}

	ch := make(chan bool, 1)
	defer func() { ch <- true }()
	go cancelQuery(id, timeout, ch, chClose)

	return q.db.Query(query)
}

// cancel query niet alleen bij timeout, maar ook als request wordt verbroken
func cancelQuery(id string, timeout bool, ch chan bool, chClose <-chan bool) {

	if timeout && Cfg.Querytimeout > 0 {
		select {
		case <-ch:
			// query is klaar
			return
		case <-chClose:
			// verbinding is gesloten: cancel query
		case <-time.After(time.Duration(Cfg.Querytimeout) * time.Second):
			// timeout: cancel query
		}
	} else {
		select {
		case <-ch:
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
