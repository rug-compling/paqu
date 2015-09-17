package main

import (
	"fmt"
	"net/http"
)

// TAB: begin
func metadata(q *Context) {

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
	   `)

	fmt.Fprintln(q.w, `<p><b>In ontwikkeling</b>`)

	html_footer(q)
}
