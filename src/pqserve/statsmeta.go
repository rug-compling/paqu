package main

import (
	"bytes"
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"
)

func statsmeta(q *Context) {

	var buf bytes.Buffer

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

	metas := getMeta(q, prefix)

	query, err := makeQuery(q, prefix, chClose)
	if err != nil {
		http.Error(q.w, err.Error(), http.StatusInternalServerError)
		logerr(err)
		return
	}

	// BEGIN UITVOER

	if download {
		q.w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		q.w.Header().Set("Content-Disposition", "attachment; filename=telling.txt")
		cache(q)
	} else {
		q.w.Header().Set("Content-Type", "text/html; charset=utf-8")
		cache(q)
		fmt.Fprint(q.w, `<!DOCTYPE html>
<html>
<head>
<title></title>
<script type="text/javascript"><!--
function f(s) {
    window.parent._fn.updatemeta(s);
}
//--></script>
</head>
<body">
<script type="text/javascript">
window.parent._fn.startedmeta();
</script>
`)
	}

	// DEBUG: HTML-uitvoer van de query
	if !download {
		fmt.Fprint(&buf, "<div style=\"font-family:monospace\">\n", html.EscapeString(query), "\n</div><p>\n")
		updateText(q, buf.String())
		buf.Reset()
	}

	// Tellingen van onderdelen
	for _, meta := range metas {
		ww := meta[0]
		for run := 0; run < 2; run++ {
			var j, count int
			var s, p, limit string
			if download {
				if run == 0 {
					fmt.Fprintln(q.w, "# "+ww+" per item\t")
				} else {
					fmt.Fprintln(q.w, "# "+ww+" per zin\t")
				}
			} else {
				if run == 0 {
					fmt.Fprintln(&buf, "<p><b>"+html.EscapeString(ww)+"</b> per item: ")
				} else {
					fmt.Fprintln(&buf, "<p><b>"+html.EscapeString(ww)+"</b> per zin: ")
				}
				limit = " LIMIT " + fmt.Sprint(WRDMAX)
			}
			select {
			case <-chClose:
				logerr(errConnectionClosed)
				return
			default:
			}
			val := "tval"
			if meta[1] == "INT" {
				val = "ival"
			} else if meta[1] == "FLOAT" {
				val = "fval"
			} else if meta[1] == "DATE" || meta[1] == "DATETIME" {
				val = "dval"
			}
			var qu string
			if run == 0 {
				qu = fmt.Sprintf(
					"SELECT count(*), `%s` FROM `%s_c_%s_deprel_meta` WHERE `name` = %q AND %s GROUP BY 2 ORDER BY 1 DESC, 2",
					val,
					Cfg.Prefix, prefix,
					ww,
					query)
			} else {
				qu = fmt.Sprintf(
					"SELECT DISTINCT `arch`,`file`,`%s` FROM `%s_c_%s_deprel_meta` WHERE `name` = %q AND %s",
					val,
					Cfg.Prefix, prefix,
					ww,
					query)
				qu = fmt.Sprintf(
					"SELECT COUNT(`a`.`%s`), `a`.`%s` FROM ( %s ) `a` GROUP BY 2 ORDER BY 1 DESC, 2",
					val, val, qu)
			}
			rows, err := timeoutQuery(q, chClose, qu+limit)
			if err != nil {
				updateError(q, err, !download)
				completedmeta(q, download)
				logerr(err)
				return
			}
			for rows.Next() {
				var err error
				if meta[1] == "DATE" {
					var t time.Time
					err = rows.Scan(&j, &t)
					s = printDate(t, false)
				} else if meta[1] == "DATETIME" {
					var t time.Time
					err = rows.Scan(&j, &t)
					s = printDate(t, true)
				} else {
					err = rows.Scan(&j, &s)
				}
				if err != nil {
					updateError(q, err, !download)
					completedmeta(q, download)
					logerr(err)
					return
				}
				if s == "" {
					s = "\"\""
				}
				s = unHigh(s)
				if download {
					fmt.Fprintf(q.w, "%d\t%s\n", j, s)
				} else {
					fmt.Fprint(&buf, p, j, "&times;&nbsp;", html.EscapeString(s))
					p = ", "
					count++
				}
			}
			err = rows.Err()
			if err != nil {
				updateError(q, err, !download)
				completedmeta(q, download)
				logerr(err)
				return
			}
			if !download {
				if count == WRDMAX {
					fmt.Fprint(&buf, ", ...")
				}
				fmt.Fprint(&buf, "\n<BR>\n")
				updateText(q, buf.String())
				buf.Reset()
			}
		}
	}

	if !download {
		fmt.Fprintf(&buf,
			"<hr>tijd: %s\n<p>\n<a href=\"statsmeta?%s&amp;d=1\">download</a>\n",
			tijd(time.Now().Sub(now)),
			strings.Replace(q.r.URL.RawQuery, "&", "&amp;", -1))
		updateText(q, buf.String())
		completedmeta(q, download)
	}
}

func completedmeta(q *Context, download bool) {
	if download {
		return
	}
	fmt.Fprintf(q.w, `<script type="text/javascript">
window.parent._fn.completedmeta();
</script>
`)
}
