package main

import (
	"fmt"
	"html"
	"net/http"
	"strconv"
	//"strings"
)

// TAB: browse (zinnen)
func browse(q *Context) {

	prefix := getprefix(q)
	if !q.prefixes[prefix] {
		http.Error(q.w, "Invalid corpus: "+prefix, http.StatusPreconditionFailed)
		return
	}

	// HTML-uitvoer van begin van de pagina
	writeHead(q, "Zinnen", 2)
	fmt.Fprint(q.w, `
<script type="text/javascript"><!--
  function formclear(f) {
    f.lbl.value = "";
  }
  //--></script>
`)

	// HTML-uitvoer van het formulier
	html_browse_form(q)

	// Maximaal 2*ZINMAX matchende xml-bestanden opvragen

	offset := 0
	o, err := strconv.Atoi(first(q.r, "offset"))
	if err == nil {
		offset = o
	}

	lbl := first(q.r, "lbl")
	query := ""
	if lbl != "" {
		query = fmt.Sprintf("WHERE `lbl` LIKE %q", lbl)
	}

	fmt.Fprintf(q.w, "<ol start=\"%d\">\n", offset+1)
	rows, err := q.db.Query(
		fmt.Sprintf(
			"SELECT `arch`,`file`,`sent`,`lbl` FROM `%s_c_%s_sent` %s LIMIT %d,%d",
			Cfg.Prefix,
			prefix,
			query,
			offset,
			2*ZINMAX))
	if doErr(q, err) {
		return
	}
	nzin := 0
	for rows.Next() {
		nzin++
		var arch, file int
		var sent, lbl string
		err := rows.Scan(&arch, &file, &sent, &lbl)
		if doErr(q, err) {
			return
		}
		fmt.Fprintf(q.w, "<li>%s\n", html.EscapeString(sent))
		fmt.Fprintf(q.w, "<a href=\"/tree?db=%s&amp;arch=%d&amp;file=%d\">&nbsp;&#9741;&nbsp;</a>\n",
			prefix, arch, file)

	}
	fmt.Fprint(q.w, "</ol>\n<p>\n")

	// Links naar volgende en vorige pagina's met resultaten
	qs := fmt.Sprintf("db=%s&amp;lbl=%s", urlencode(prefix), urlencode(lbl))
	if offset > 0 || nzin == 2*ZINMAX {
		if offset > 0 {
			fmt.Fprintf(q.w, "<a href=\"/browse?%s&amp;offset=%d\">vorige</a>", qs, offset-2*ZINMAX)
		} else {
			fmt.Fprint(q.w, "vorige")
		}
		fmt.Fprint(q.w, " | ")
		if nzin == 2*ZINMAX {
			fmt.Fprintf(q.w, "<a href=\"/browse?%s&amp;offset=%d\">volgende</a>", qs, offset+2*ZINMAX)
		} else {
			fmt.Fprint(q.w, "volgende")
		}
	}

	fmt.Fprint(q.w, "</body>\n</html>\n")

}

func html_browse_form(q *Context) {

	fmt.Fprint(q.w, `
<form action="browse" method="get" accept-charset="utf-8">
corpus: <select name="db">
`)
	html_opts(q, q.opt_db, getprefix(q), "corpus")
	fmt.Fprintf(q.w, `</select>
<p>
label: <input type="text" name="lbl" size="20" value="%s">
`, html.EscapeString(first(q.r, "lbl")))

	fmt.Fprintf(q.w, `<p>
		<input type="submit" value="Selecteren">
		<input type="button" value="Wissen" onClick="javascript:formclear(form)">
		<input type="reset" value="Reset">
	   </form>
`)
}
