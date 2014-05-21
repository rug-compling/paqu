package main

import (
	"fmt"
	"html"
	"net/http"
	"strings"
)

func rename(q *Context) {
	if !q.auth {
		http.Error(q.w, "Je bent niet ingelogd", http.StatusUnauthorized)
		return
	}

	id := first(q.r, "id")
	if !q.myprefixes[id] {
		http.Error(q.w, "Dat is niet je corpus", http.StatusUnauthorized)
		return
	}

	q.w.Header().Set("Content-type", "text/html; charset=utf-8")
	q.w.Header().Set("Cache-Control", "no-cache")
	q.w.Header().Add("Pragma", "no-cache")
	fmt.Fprintf(q.w, file__rename__html, id, html.EscapeString(q.desc[id]))
}

func rename2(q *Context) {
	if !q.auth {
		http.Error(q.w, "Je bent niet ingelogd", http.StatusUnauthorized)
		return
	}

	id := first(q.r, "id")
	if !q.myprefixes[id] {
		http.Error(q.w, "Dat is niet je corpus", http.StatusUnauthorized)
		return
	}

	d2 := strings.TrimSpace(first(q.r, "desc"))
	if d2 == "" {
		writeHtml(q, "Corpus niet hernoemd", "Corpus is niet hernoemd", "corpora")
	} else {
		_, err := q.db.Exec(fmt.Sprintf("UPDATE `%s_info` SET `description` = %q WHERE `id` = %q", Cfg.Prefix, d2, id))
		if err != nil {
			http.Error(q.w, err.Error(), http.StatusInternalServerError)
			logerr(err)
			return
		}
		writeHtml(
			q,
			"Corpus hernoemd",
			fmt.Sprintf(`
Van: <em>%s</em><br>
Naar: <em>%s</em>
`, html.EscapeString(q.desc[id]), html.EscapeString(d2)),
			"corpora")

	}
}
