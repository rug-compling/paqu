package main

import (
	"bytes"
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"
)

func stats(q *Context) {

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
    window.parent._fn.update(s);
}
//--></script>
</head>
<body">
<script type="text/javascript">
window.parent._fn.started();
</script>
`)
	}

	// DEBUG: HTML-uitvoer van de query
	if !download {
		fmt.Fprint(&buf, "<div style=\"font-family:monospace\">\n", html.EscapeString(query), "\n</div><p>\n")
		updateText(q, buf.String())
		buf.Reset()
	}

	// Aantal zinnen die matchen met de query
	select {
	case <-chClose:
		logerr(errConnectionClosed)
		return
	default:
	}
	rows, err := timeoutQuery(q, chClose, "SELECT DISTINCT `arch`,`file` FROM `"+Cfg.Prefix+"_c_"+prefix+"_deprel` WHERE "+
		query)
	if err != nil {
		updateError(q, err, !download)
		completed(q, download)
		logerr(err)
		return
	}
	counter := 0
	for rows.Next() {
		counter++
	}
	if err != nil {
		updateError(q, err, !download)
		completed(q, download)
		logerr(err)
		return
	}

	if download {
		fmt.Fprintf(q.w, "# %d zinnen\t\n", counter)
	} else {
		fmt.Fprintln(&buf, "Aantal gevonden zinnen:", iformat(counter))
		updateText(q, buf.String())
		buf.Reset()
	}

	// Tellingen van onderdelen
	for i, ww := range []string{"word", "lemma", "postag", "rel", "hword", "hlemma", "hpostag"} {
		var j, count int
		var s, p, limit string
		if download {
			fmt.Fprintln(q.w, "# "+ww+"\t")
		} else {
			if i == 0 {
				fmt.Fprintln(&buf, "<p>"+YELLOW+"<b>word</b></span>: ")
			} else if i == 4 {
				fmt.Fprintln(&buf, "<p>"+GREEN+"<b>hword</b></span>: ")
			} else {
				fmt.Fprintln(&buf, "<p><b>"+ww+"</b>: ")
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
			updateError(q, err, !download)
			completed(q, download)
			logerr(err)
			return
		}
		for rows.Next() {
			err := rows.Scan(&j, &s)
			if err != nil {
				updateError(q, err, !download)
				completed(q, download)
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
			completed(q, download)
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

	if !download {
		fmt.Fprintf(&buf,
			"<hr>tijd: %s\n<p>\n<a href=\"stats?%s&amp;d=1\">download</a>\n",
			tijd(time.Now().Sub(now)),
			strings.Replace(q.r.URL.RawQuery, "&", "&amp;", -1))
		updateText(q, buf.String())
		completed(q, download)
	}
}

func completed(q *Context, download bool) {
	if download {
		return
	}
	fmt.Fprintf(q.w, `<script type="text/javascript">
window.parent._fn.completed();
</script>
`)
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

	hasmeta := hasMeta(q, prefix)
	var metas []MetaType
	if hasmeta {
		metas = getMeta(q, prefix)
	}
	metai := make(map[string]int)
	metat := make(map[string]string)
	for _, meta := range metas {
		metai[meta.name] = meta.id
		metat[meta.name] = meta.mtype
	}

	qselect := "COUNT(*)"
	qfrom := fmt.Sprintf("`%s_c_%s_deprel` `a`", Cfg.Prefix, prefix)
	qwhere := ""
	ncols := 0
	fields := make([]interface{}, 0)
	fields = append(fields, new(int))

	cols := make([]string, 1, 8)
	aligns := make([]string, 1, 8)
	aligns[0] = "right"
	for _, t := range []string{"word", "lemma", "postag", "rel", "hword", "hlemma", "hpostag"} {
		if first(q.r, "c"+t) == "1" {
			qselect += ",`a`.`" + t + "`"
			ncols++
			cols = append(cols, t)
			aligns = append(aligns, "left")
			fields = append(fields, new(string))
		}
	}

	for _, meta := range q.r.Form["cmeta"] {
		cols = append(cols, "meta:"+meta)
		if metat[meta] == "TEXT" {
			aligns = append(aligns, "left")
		} else {
			aligns = append(aligns, "right")
		}
		ncols++
		fields = append(fields, new(string))
		table := fmt.Sprintf("m%d", metai[meta])
		qfrom += fmt.Sprintf(" JOIN ( `%s_c_%s_meta` `%s`", Cfg.Prefix, prefix, table)
		qfrom += fmt.Sprintf(" JOIN `%s_c_%s_mval` `%sv` USING(`id`,`idx`)", Cfg.Prefix, prefix, table)
		qfrom += " ) USING(`arch`,`file`)"
		qwhere += fmt.Sprintf(" AND `%s`.`id` = %d", table, metai[meta])
		qselect += ", `" + table + "v`.`text`"
	}

	query, err := makeQuery(q, prefix, "a", chClose)
	if err != nil {
		http.Error(q.w, err.Error(), http.StatusInternalServerError)
		logerr(err)
		return
	}

	qgroupby := ""
	qorder := "1 DESC"
	for i := 0; i < ncols; i++ {
		n := fmt.Sprint(i + 2)
		if i > 0 {
			qgroupby += ","
		}
		qgroupby += n
		qorder += "," + n
	}
	fullquery := fmt.Sprintf("SELECT %s FROM %s WHERE %s %s GROUP BY %s ORDER BY %s",
		qselect, qfrom, query, qwhere, qgroupby, qorder)

	qword := urlencode(first(q.r, "word"))
	qpostag := urlencode(first(q.r, "postag"))
	qrel := urlencode(first(q.r, "rel"))
	qhword := urlencode(first(q.r, "hword"))
	qhpostag := urlencode(first(q.r, "hpostag"))
	qdb := urlencode(first(q.r, "db"))

	// BEGIN UITVOER

	if download {
		q.w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		q.w.Header().Set("Content-Disposition", "attachment; filename=telling.txt")
	} else {
		q.w.Header().Set("Content-Type", "text/html; charset=utf-8")
	}
	cache(q)

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

	if !download {
		fullquery += fmt.Sprintf(" LIMIT %d", WRDMAX+1)
	}

	rows, err := timeoutQuery(q, chClose, fullquery)
	if err != nil {
		interneFoutRegel(q, err, !download)
		logerr(err)
		return
	}

	if !download {
		fmt.Fprintln(q.w, "<table class=\"breed\"><tr class=\"odd\">")
		for _, c := range cols {
			if strings.HasPrefix(c, "meta:") {
				c = c[5:]
			}
			fmt.Fprintln(q.w, "<th>"+c)
		}
	} else {
		fmt.Fprint(q.w, "aantal")
		for _, c := range cols[1:] {
			fmt.Fprint(q.w, "\t"+c)
		}
		fmt.Fprintln(q.w)
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
			if n%2 == 0 {
				fmt.Fprintln(q.w, "<tr class=\"odd\">")
			} else {
				fmt.Fprintln(q.w, "<tr>")
			}
		}
		for i, e := range fields {
			var value string
			switch v := e.(type) {
			case *string:
				value = unHigh(*v)
			case *int:
				value = fmt.Sprint(*v)
			}
			if !download && n > WRDMAX {
				if i == 0 {
					value = " " // met spatie!
				} else {
					value = "..."
				}
			}
			if !download {
				var a1, a2 string
				if i == 0 && n <= WRDMAX {
					for j := len(cols) - 1; j > 0; j-- { // van achter naar voor zodat word prioriteit krijgt over lemma
						if sp, ok := fields[j].(*string); ok {
							s := *sp
							switch cols[j] {
							case "word":
								qword = urlencode("=" + unHigh(s))
							case "lemma":
								qword = urlencode("+" + unHigh(s))
							case "postag":
								qpostag = urlencode(s)
							case "rel":
								qrel = urlencode(s)
							case "hword":
								qhword = urlencode("=" + unHigh(s))
							case "hlemma":
								qhword = urlencode("+" + unHigh(s))
							case "hpostag":
								qhpostag = urlencode(s)
								if qhpostag == "" {
									qhpostag = "--LEEG--"
								}
							}
						}
					}
					a1 = fmt.Sprintf("<a href=\".?db=%s&amp;word=%s&amp;postag=%s&amp;rel=%s&amp;hword=%s&amp;hpostag=%s\">",
						qdb, qword, qpostag, qrel, qhword, qhpostag)
					a2 = "</a>"
				}
				var s string
				var c string
				if value == "" {
					s = "(leeg)"
					c = " nil"
				} else {
					s = html.EscapeString(value)
				}
				fmt.Fprintf(q.w, " <td class=\"%s%s\">%s%s%s\n", aligns[i], c, a1, s, a2)
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
			"<hr>tijd: %s\n<p>\n<a href=\"statsrel?%s&amp;d=1\">download</a>\n",
			tijd(time.Now().Sub(now)),
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
