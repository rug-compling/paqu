package main

import (
	"fmt"
	"html"
	"net/http"
	"strings"
)

// TAB: metadata
func metadata(q *Context) {

	if first(q.r, "d") != "" {
		metadownload(q)
		return
	}

	prefix := getprefix(q)
	if !q.prefixes[prefix] {
		http.Error(q.w, "Invalid corpus: "+prefix, http.StatusPreconditionFailed)
		return
	}

	// HTML-uitvoer van begin van de pagina
	writeHead(q, "", 3)

	fmt.Fprint(q.w, `
<form action="metadata" method="get" accept-charset="utf-8">
corpus: <select name="db">
`)
	html_opts(q, q.opt_dbmeta, getprefix(q), "corpus")
	fmt.Fprintln(q.w, "</select>")
	if q.auth {
		fmt.Fprintln(q.w, "<a href=\"corpuslijst?meta=1\">meer/minder</a>")
	}
	fmt.Fprint(q.w, `
       <p>
       <input type="submit" value="Kies corpus">
	   </form>
       <p>
	   `)

	if !q.hasmeta[prefix] {
		html_footer(q)
		return
	}
	fmt.Fprintln(q.w, "<hr><p>")

	metas := getMeta(q, prefix)
	for _, meta := range metas {
		fmt.Fprintf(q.w, "<b>%s</b><p>\n", html.EscapeString(meta.name))
		o := `idx`
		align := "right"
		limit := ""
		if meta.mtype == "TEXT" {
			o = "`n` DESC, `idx`"
			align = "left"
			limit = fmt.Sprintf(" LIMIT 0, %d", METAMAX)
		}
		rows, err := q.db.Query(fmt.Sprintf(
			"SELECT `text`, `n` FROM `%s_c_%s_mval` WHERE `id`=%d ORDER BY %s%s",
			Cfg.Prefix, prefix,
			meta.id,
			o,
			limit))
		if doErr(q, err) {
			return
		}
		fmt.Fprintln(q.w, "<table>")
		count := 0
		for rows.Next() {
			count++
			var s string
			var n int
			if doErr(q, rows.Scan(&s, &n)) {
				return
			}
			cnil := ""
			if s == "" {
				s = "(leeg)"
				cnil = " nil"
			} else {
				s = html.EscapeString(s)
			}
			fmt.Fprintf(q.w, "<tr><td class=\"right\">%s<td class=\"%s%s\">%s\n", iformat(n), align, cnil, s)
		}
		if meta.mtype == "TEXT" && count == METAMAX {
			fmt.Fprintln(q.w, "<tr><td><td>...")
		}
		fmt.Fprintln(q.w, "</table><p>")
		if doErr(q, rows.Err()) {
			return
		}
	}
	fmt.Fprintf(q.w,
		"<a href=\"metadata?%s&amp;d=1\">download</a><p>\n",
		strings.Replace(q.r.URL.RawQuery, "&", "&amp;", -1))

	fmt.Fprintln(q.w, `<p><b>In ontwikkeling</b>`)

	html_footer(q)
}

func metadownload(q *Context) {

	prefix := getprefix(q)
	if !q.hasmeta[prefix] {
		http.Error(q.w, "Invalid corpus: "+prefix, http.StatusPreconditionFailed)
		return
	}

	q.w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	q.w.Header().Set("Content-Disposition", "attachment; filename=metadata.txt")
	cache(q)

	metas := getMeta(q, prefix)
	for _, meta := range metas {
		fmt.Fprintf(q.w, "# %s\n", meta.name)
		o := `idx`
		if meta.mtype == "TEXT" {
			o = "`n` DESC, `idx`"
		}
		rows, err := q.db.Query(fmt.Sprintf(
			"SELECT `text`, `n` FROM `%s_c_%s_mval` WHERE `id`=%d ORDER BY %s",
			Cfg.Prefix, prefix,
			meta.id,
			o))
		if doErr(q, err) {
			return
		}
		for rows.Next() {
			var s string
			var n int
			if doErr(q, rows.Scan(&s, &n)) {
				return
			}
			fmt.Fprintf(q.w, "%d\t%s\n", n, s)
		}
		if doErr(q, rows.Err()) {
			return
		}
	}
}
