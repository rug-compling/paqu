package main

import (
	"fmt"
	"html"
	"net/http"
	"strings"
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
		q.w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	} else {
		q.w.Header().Set("Content-Type", "text/html; charset=utf-8")
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
		logerr(errConnectionClosed)
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
			logerr(errConnectionClosed)
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
			"<hr>tijd: %s\n<p>\n<a href=\"stats?%s&amp;d=1\" target=\"_blank\">download</a>\n",
			time.Now().Sub(now),
			strings.Replace(q.r.URL.RawQuery, "&", "&amp;", -1))
	}
}

func statsrel(q *Context) {

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

	items := make([]string, 0, 7)
	cols := make([]string, 1, 8)
	for _, t := range []string{"word", "lemma", "postag", "rel", "hword", "hlemma", "hpostag"} {
		if first(q.r, "c"+t) == "1" {
			items = append(items, "`"+t+"`")
			cols = append(cols, t)
		}
	}
	st := strings.Join(items, ",")

	query, err := makeQuery(q, prefix, chClose)
	if err != nil {
		http.Error(q.w, err.Error(), http.StatusInternalServerError)
		logerr(err)
		return
	}

	qword := urlencode(first(q.r, "word"))
	qpostag := urlencode(first(q.r, "postag"))
	qrel := urlencode(first(q.r, "rel"))
	qhword := urlencode(first(q.r, "hword"))
	qhpostag := urlencode(first(q.r, "hpostag"))
	qdb := urlencode(first(q.r, "db"))

	// BEGIN UITVOER

	if download {
		q.w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	} else {
		q.w.Header().Set("Content-Type", "text/html; charset=utf-8")
	}
	q.w.Header().Set("Cache-Control", "no-cache")
	q.w.Header().Add("Pragma", "no-cache")

	// DEBUG: HTML-uitvoer van de query
	if !download {
		fmt.Fprint(q.w, "<div style=\"font-family:monospace\">\n", html.EscapeString(query), "\n</div><p>\n")
	}

	select {
	case <-chClose:
		logerr(errConnectionClosed)
		return
	default:
	}

	limit := ""
	if !download {
		limit = fmt.Sprintf(" LIMIT %d", WRDMAX + 1)
	}

	rows, err := timeoutQuery(q, chClose, "SELECT count(*),"+st+" FROM `"+Cfg.Prefix+"_c_"+prefix+"_deprel` WHERE "+
		query+" GROUP BY "+st+" ORDER BY 1 DESC,"+st+limit)
	if err != nil {
		interneFoutRegel(q, err, !download)
		logerr(err)
		return
	}

	if !download {
		fmt.Fprintln(q.w, "<table><tr>")
		for _, c := range cols {
			fmt.Fprintln(q.w, "<th>"+c)
		}
	} else {
		fmt.Fprint(q.w, "aantal")
		for _, c := range cols[1:] {
			fmt.Fprint(q.w, "\t"+c)
		}
		fmt.Fprintln(q.w)
	}

	fields := make([]interface{}, 1+len(items))
	fields[0] = new(string)
	for i := 0; i < len(items); i++ {
		fields[i+1] = new(string)
	}
	n := 0
	for rows.Next() {
		n++
		err := rows.Scan(fields...)
		if err != nil {
			interneFoutRegel(q, err, !download)
			logerr(err)
			return
		}
		if !download {
			fmt.Fprintln(q.w, "<tr>")
		}
		for i, e := range fields {
			value := *(e.(*string))
			if !download && n > WRDMAX {
				value = "..."
			}
			if !download {
				var a1, a2, class string
				if i == 0 {
					for j := len(cols) - 1; j > 0; j-- { // van achter naar voor zodat word prioriteit krijgt over lemma
						s := *fields[j].(*string)
						switch cols[j] {
						case "word":
							qword = urlencode("=" + s)
						case "lemma":
							qword = urlencode("+" + s)
						case "postag":
							qpostag = urlencode(s)
						case "rel":
							qrel = urlencode(s)
						case "hword":
							qhword = urlencode("=" + s)
						case "hlemma":
							qhword = urlencode("+" + s)
						case "hpostag":
							qhpostag = urlencode(s)
						}
					}
					a1 = fmt.Sprintf("<a href=\".?db=%s&amp;word=%s&amp;postag=%s&amp;rel=%s&amp;hword=%s&amp;hpostag=%s\">",
						qdb, qword, qpostag, qrel, qhword, qhpostag)
					a2 = "</a>"
					class = " class=\"right\""
				}
				fmt.Fprintf(q.w, " <td%s>%s%s%s\n", class, a1, html.EscapeString(value), a2)
			} else {
				var t string
				if i != 0 {
					t = "\t"
				}
				fmt.Fprintf(q.w, t+value)
			}
		}
		if download {
			fmt.Fprintln(q.w)
		}
	}

	if download {
		return
	}

	fmt.Fprintln(q.w, "</table>")

	if !download {
		fmt.Fprintf(q.w,
			"<hr>tijd: %s\n<p>\n<a href=\"statsrel?%s&amp;d=1\" target=\"_blank\">download</a>\n",
			time.Now().Sub(now),
			strings.Replace(q.r.URL.RawQuery, "&", "&amp;", -1))
	}
}

func interneFoutRegel(q *Context, err error, is_html bool) {
	s := err.Error()
	if is_html {
		s = html.EscapeString(s)
	}
	fmt.Fprintln(q.w, "Interne fout:", s)
}
