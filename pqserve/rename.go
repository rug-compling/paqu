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

	writeHead(q, "Corpus hernoemen", 0)
	fmt.Fprintf(q.w, `
<h1>Corpus hernoemen</h1>
<form action="rename2" method="get" accept-charset="utf-8">
  <input type="hidden" name="id" value="%s">
  Nieuwe naam: <input type="text" name="desc" value="%s" size="%d" maxlength="%d">
<p>
<input type="submit">
</form>
</body>
</html>
`, id, html.EscapeString(q.desc[id]), MAXTITLELEN+MAXTITLELEN/4, MAXTITLELEN)
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

	d2 := maxtitlelen(strings.TrimSpace(first(q.r, "desc")))
	if d2 != "" {
		_, err := q.db.Exec(fmt.Sprintf("UPDATE `%s_info` SET `description` = %q WHERE `id` = %q", Cfg.Prefix, d2, id))
		if err != nil {
			http.Error(q.w, err.Error(), http.StatusInternalServerError)
			logerr(err)
			return
		}
	}
	http.Redirect(q.w, q.r, urlJoin(Cfg.Url, "corpora"), http.StatusTemporaryRedirect)
}
