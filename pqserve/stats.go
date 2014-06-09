package main

import (
	"fmt"
	"html"
	"net/http"
	"time"
)

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

	query, err := makeQuery(q, prefix, chClose)
	if err != nil {
		http.Error(q.w, err.Error(), http.StatusInternalServerError)
		logerr(err)
		return
	}

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
