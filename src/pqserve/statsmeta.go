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

		// ranges definieren

		var ir *irange
		var fr *frange
		var dr *drange
		ww := meta[0]
		if meta[1] == "INT" {
			rows, err := q.db.Query(fmt.Sprintf(
				"SELECT SQL_CACHE MIN(`ival`), MAX(`ival`), COUNT(DISTINCT `ival`) FROM `%s_c_%s_meta` JOIN `%s_c_%s_midx` USING (`id`) WHERE `name` = %q",
				Cfg.Prefix, prefix,
				Cfg.Prefix, prefix,
				ww))
			if err != nil {
				updateError(q, err, !download)
				completedmeta(q, download)
				logerr(err)
				return
			}
			var v1, v2, vx int
			for rows.Next() {
				rows.Scan(&v1, &v2, &vx)
			}
			ir = newIrange(v1, v2, vx)
		} else if meta[1] == "FLOAT" {
			rows, err := q.db.Query(fmt.Sprintf(
				"SELECT SQL_CACHE MIN(`fval`), MAX(`fval`) FROM `%s_c_%s_meta` JOIN `%s_c_%s_midx` USING (`id`) WHERE `name` = %q",
				Cfg.Prefix, prefix,
				Cfg.Prefix, prefix,
				ww))
			if err != nil {
				updateError(q, err, !download)
				completedmeta(q, download)
				logerr(err)
				return
			}
			var v1, v2 float64
			for rows.Next() {
				rows.Scan(&v1, &v2)
			}
			fr = newFrange(v1, v2)
		} else if meta[1] == "DATE" || meta[1] == "DATETIME" {
			dis := "0"
			if meta[1] == "DATE" {
				dis = "COUNT(DISTINCT `dval`)"
			}
			rows, err := timeoutQuery(q, chClose, fmt.Sprintf(
				"SELECT SQL_CACHE MIN(`dval`), MAX(`dval`), %s FROM `%s_c_%s_meta` JOIN `%s_c_%s_midx` USING (`id`) WHERE `name` = %q",
				dis,
				Cfg.Prefix, prefix,
				Cfg.Prefix, prefix,
				ww))
			if err != nil {
				updateError(q, err, !download)
				completedmeta(q, download)
				logerr(err)
				return
			}
			var v1, v2 time.Time
			var i int
			for rows.Next() {
				rows.Scan(&v1, &v2, &i)
			}
			dr = newDrange(v1, v2, i, meta[1] == "DATETIME")
		}

		// telling van metadata over alle zinnen

		fsum := 0
		fcount := make(map[string]int)
		val := "`tval`"
		if meta[1] == "INT" {
			val = ir.sql()
		} else if meta[1] == "FLOAT" {
			val = fr.sql()
		} else if meta[1] == "DATE" || meta[1] == "DATETIME" {
			val = dr.sql()
		}
		// hier moet (de view van) de join met *_deprel gebruikt worden, omdat mogelijk niet alle zinnen trippels hebben
		// anders zou je een join van alleen *_meta en *_midx kunnen gebruiken
		qu := fmt.Sprintf(
			"SELECT DISTINCT `arch`,`file`,%s AS `val` FROM `%s_c_%s_deprel_meta` WHERE `name` = %q",
			val,
			Cfg.Prefix, prefix,
			ww)
		qu = fmt.Sprintf(
			"SELECT SQL_CACHE COUNT(`a`.`val`), `a`.`val` FROM ( %s ) `a` GROUP BY 2",
			qu)
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
			var s string
			var j int
			var err error
			if meta[1] == "DATE" || meta[1] == "DATETIME" {
				var v time.Time
				err = rows.Scan(&j, &v)
				s, _ = dr.value(v)
			} else if meta[1] == "INT" {
				var v int
				err = rows.Scan(&j, &v)
				s, _ = ir.value(v)
			} else if meta[1] == "FLOAT" {
				var v float64
				err = rows.Scan(&j, &v)
				s, _ = fr.value(v)
			} else {
				err = rows.Scan(&j, &s)
			}
			if err != nil {
				updateError(q, err, !download)
				completedmeta(q, download)
				logerr(err)
				return
			}
			s = unHigh(s)
			fsum += j
			fcount[s] = j
		}

		// telling van metadata in matchende zinnen

		var limit, order string
		if meta[1] == "TEXT" {
			order = "1 DESC, 2"
			if !download {
				limit = "LIMIT " + fmt.Sprint(WRDMAX)
			}
		} else {
			order = "2"
		}
		for run := 0; run < 2; run++ {
			lines := make([]Statline, 0)
			var j, count int
			var s string
			var sum = 0.0
			if download {
				if run == 0 {
					fmt.Fprintln(q.w, "# "+ww+" per item\t")
				} else {
					fmt.Fprintln(q.w, "# "+ww+" per zin\t")
				}
			} else {
				if run == 0 {
					fmt.Fprintln(&buf, "<p><b>"+html.EscapeString(ww)+"</b><table><tr><td>per item:<table class=\"right\">")
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
			val := "`tval`"
			if meta[1] == "INT" {
				val = ir.sql()
			} else if meta[1] == "FLOAT" {
				val = fr.sql()
			} else if meta[1] == "DATE" || meta[1] == "DATETIME" {
				val = dr.sql()
			}
			var qu string
			if run == 0 {
				qu = fmt.Sprintf(
					"SELECT COUNT(*), %s FROM `%s_c_%s_deprel_meta` WHERE `name` = %q AND %s GROUP BY 2 ORDER BY %s %s",
					val,
					Cfg.Prefix, prefix,
					ww,
					query,
					order,
					limit)
			} else {
				qu = fmt.Sprintf(
					"SELECT DISTINCT `arch`,`file`,%s AS `val` FROM `%s_c_%s_deprel_meta` WHERE `name` = %q AND %s",
					val,
					Cfg.Prefix, prefix,
					ww,
					query)
				qu = fmt.Sprintf(
					"SELECT COUNT(`a`.`val`), `a`.`val` FROM ( %s ) `a` GROUP BY 2 ORDER BY %s %s",
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
				var err error
				if meta[1] == "DATE" || meta[1] == "DATETIME" {
					var v time.Time
					err = rows.Scan(&j, &v)
					s, _ = dr.value(v)
				} else if meta[1] == "INT" {
					var v int
					err = rows.Scan(&j, &v)
					s, _ = ir.value(v)
				} else if meta[1] == "FLOAT" {
					var v float64
					err = rows.Scan(&j, &v)
					s, _ = fr.value(v)
				} else {
					err = rows.Scan(&j, &s)
				}
				if err != nil {
					updateError(q, err, !download)
					completedmeta(q, download)
					logerr(err)
					return
				}
				s = unHigh(s)
				lines = append(lines, Statline{s, j})
				sum += float64(j) / (float64(fcount[s]) / float64(fsum))
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
					p = float64(line.i) / (float64(fcount[line.s]) / float64(fsum)) / sum * 100.0
				}
				if download {
					if run == 1 {
						fmt.Fprintf(q.w, "%d\t%.2f%%\t%s\n", line.i, float32(p), line.s)
					} else {
						fmt.Fprintf(q.w, "%d\t%s\n", line.i, line.s)
					}
				} else {
					td := "<td>"
					if meta[1] == "TEXT" {
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
				if count == WRDMAX {
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
