package main

import (
	"github.com/microcosm-cc/bluemonday"
	"gopkg.in/russross/blackfriday.v2"

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
  Optioneel, toelichting:<br>
  <textarea rows="6" cols="80" id="infotext" name="infotext" maxlength="4000" placeholder="tekst in markdown-formaat">%s</textarea>
<p>
<input type="submit">
</form>
</body>
</html>
`, id, html.EscapeString(q.desc[id]), MAXTITLELEN+MAXTITLELEN/4, MAXTITLELEN, html.EscapeString(q.infos[id]))
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

	info := strings.Replace(first(q.r, "infotext"), "\r\n", "\n", -1)
	unsafe := blackfriday.Run([]byte(info))
	infop := strings.TrimSpace(string(bluemonday.UGCPolicy().SanitizeBytes(unsafe)))
	_, err := q.db.Exec(fmt.Sprintf("UPDATE `%s_info` SET `info` = %q WHERE `id` = %q", Cfg.Prefix, info, id))
	if err != nil {
		http.Error(q.w, err.Error(), http.StatusInternalServerError)
		logerr(err)
		return
	}
	_, err = q.db.Exec(fmt.Sprintf("UPDATE `%s_info` SET `infop` = %q WHERE `id` = %q", Cfg.Prefix, infop, id))
	if err != nil {
		http.Error(q.w, err.Error(), http.StatusInternalServerError)
		logerr(err)
		return
	}

	http.Redirect(q.w, q.r, urlJoin(Cfg.Url, "corpora"), http.StatusTemporaryRedirect)
}
