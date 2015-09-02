package main

import (
	"bytes"
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"
)

type Statline struct {
	s string
	i int
	r float64
}

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

	query, err := makeQuery(q, prefix, "", chClose)
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

		// telling van metadata in matchende zinnen

		for run := 0; run < 2; run++ {
			lines := make([]Statline, 0)
			var count int
			var sum = 0.0
			if download {
				if run == 0 {
					fmt.Fprintln(q.w, "# "+meta.name+" per item\t")
				} else {
					fmt.Fprintln(q.w, "# "+meta.name+" per zin\t")
				}
			} else {
				if run == 0 {
					fmt.Fprintln(&buf, "<p><b>"+html.EscapeString(meta.name)+"</b><table><tr><td>per item:<table class=\"right\">")
				} else {
					fmt.Fprintln(&buf, "<td>per zin:<table class=\"right\">")
				}
			}
			select {
			case <-chClose:
				logerr(errConnectionClosed)
				return
			default:
			}
			var qu string
			var limit, order string
			if meta.mtype == "TEXT" {
				order = "1 DESC, `idx`"
				if !download {
					limit = "LIMIT " + fmt.Sprint(METAMAX)
				}
			} else {
				order = "`idx`"
			}
			if run == 0 {
				qu = fmt.Sprintf(
					"SELECT COUNT(`text`), `text`, 0 "+
						"FROM `%s_c_%s_deprel` "+
						"JOIN `%s_c_%s_meta` USING(`arch`,`file`) "+
						"JOIN `%s_c_%s_mval` USING (`id`,`idx`) "+
						"WHERE `id` = %d AND %s "+
						"GROUP BY `text` ORDER BY %s %s",
					Cfg.Prefix, prefix,
					Cfg.Prefix, prefix,
					Cfg.Prefix, prefix,
					meta.id,
					query,
					order,
					limit)
			} else {
				qu = fmt.Sprintf(
					"SELECT DISTINCT `arch`,`file`,`idx`,`text`,`rtrip` "+
						"FROM `%s_c_%s_deprel` "+
						"JOIN `%s_c_%s_meta` USING(`arch`,`file`) "+
						"JOIN `%s_c_%s_mval` USING (`id`,`idx`) "+
						"WHERE `id` = %d AND %s",
					Cfg.Prefix, prefix,
					Cfg.Prefix, prefix,
					Cfg.Prefix, prefix,
					meta.id,
					query)
				qu = fmt.Sprintf(
					"SELECT COUNT(`a`.`text`), `a`.`text`,`a`.`rtrip` FROM ( %s ) `a` GROUP BY `a`.`idx` ORDER BY %s %s",
					qu,
					order,
					limit)
			}
			rows, err := timeoutQuery(q, chClose, qu)
			if err != nil {
				updateError(q, err, !download)
				completedmeta(q, download)
				logerr(err)
				return
			}
			for rows.Next() {
				select {
				case <-chClose:
					rows.Close()
					logerr(errConnectionClosed)
					return
				default:
				}
				var cnt int
				var text string
				var r float64
				err := rows.Scan(&cnt, &text, &r)
				if err != nil {
					updateError(q, err, !download)
					completedmeta(q, download)
					logerr(err)
					return
				}
				text = unHigh(text)
				lines = append(lines, Statline{text, cnt, float64(cnt) / r})
				sum += float64(cnt) / r
			}
			err = rows.Err()
			if err != nil {
				updateError(q, err, !download)
				completedmeta(q, download)
				logerr(err)
				return
			}
			for _, line := range lines {
				var p float64
				if run == 1 {
					p = float64(line.r) / sum * 100.0
				}
				if download {
					if run == 1 {
						fmt.Fprintf(q.w, "%d\t%.2f%%\t%s\n", line.i, float32(p), line.s)
					} else {
						fmt.Fprintf(q.w, "%d\t%s\n", line.i, line.s)
					}
				} else {
					td := "<td>"
					if meta.mtype == "TEXT" {
						td = "<td class=\"left\">"
						count++
					}
					fmt.Fprintln(&buf, "<tr><td>", line.i)
					if run == 1 {
						fmt.Fprintf(&buf, "<td>%.2f%%", float32(p))
					}
					fmt.Fprintln(&buf, td, html.EscapeString(line.s))
				}
			}
			if !download {
				if count == METAMAX {
					fmt.Fprint(&buf, "<tr><td><td class=\"left\">...")
				}
				fmt.Fprintln(&buf, "</table>")
				if run == 1 {
					fmt.Fprintln(&buf, "</table>")
					updateText(q, buf.String())
					buf.Reset()
				}
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
