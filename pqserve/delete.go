package main

import (
	"fmt"
	"html"
	"net/http"
	"os"
	"path"
	"time"
)

func myCorpus(q *Context, id string) (bool, string) {
	rows, err := q.db.Query(fmt.Sprintf(
		"SELECT `description` FROM `%s_info` WHERE `id` = %q AND `owner` = %q",
		Cfg.Prefix, id, q.user))
	if err != nil {
		logerr(err)
		return false, ""
	}
	var desc string
	if rows.Next() {
		err := rows.Scan(&desc)
		if err == nil {
			rows.Close()
			return true, desc
		} else {
			logerr(err)
		}
	}
	return false, ""
}

func remove(q *Context) {

	if !q.auth {
		http.Error(q.w, "Je bent niet ingelogd", http.StatusUnauthorized)
		return
	}

	id := first(q.r, "id")

	my, desc := myCorpus(q, id)
	if !my {
		http.Error(q.w, "Dat is niet je corpus", http.StatusUnauthorized)
		return
	}

	writeHead(q, "Verwijder corpus")
	fmt.Fprintf(q.w, `
<h1>Verwijder corpus</h1>
Corpus: <b>%s</b>
<p>
Weet je zeker dat je dit corpus wilt verwijderen?
<p>
<form action="delete2">
<input type="hidden" name="id" value="%v">
<button type="submit">Verwijder corpus</button>
</form>
<div class="next">
<a href="corpora">Terug</a>
</div>
</body>
</html>
`,
		html.EscapeString(desc),
		id)

}

func remove2(q *Context) {

	if !q.auth {
		http.Error(q.w, "Je bent niet ingelogd", http.StatusUnauthorized)
		return
	}

	id := first(q.r, "id")

	my, desc := myCorpus(q, id)
	if !my {
		http.Error(q.w, "Dat is niet je corpus", http.StatusUnauthorized)
		return
	}

	kill(id)

	go func() {
		p := path.Join(Cfg.Data, id)
		p2 := path.Join(Cfg.Data, "_invalid_"+id)
		err := os.Rename(p, p2)
		if err != nil {
			logf("os.Rename(%q, %q) error: %v", p, p2, err)
		}
		for i := 0; i < 12; i++ {
			err = os.RemoveAll(p2)
			if err == nil {
				break
			}
			time.Sleep(10 * time.Second)
		}
		if err != nil {
			logf("os.RemoveAll(%q) error: %v", p2, err)
		} else {
			logf("os.RemoveAll(%q) OK", p2)
		}
	}()

	_, err := q.db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%s_c_%s_deprel`, `%s_c_%s_sent`, `%s_c_%s_file`, `%s_c_%s_arch` , `%s_c_%s_word`",
		Cfg.Prefix, id,
		Cfg.Prefix, id,
		Cfg.Prefix, id,
		Cfg.Prefix, id,
		Cfg.Prefix, id))
	logerr(err)
	_, err = q.db.Exec(fmt.Sprintf("DELETE FROM `%s_corpora` WHERE `prefix` = %q", Cfg.Prefix, id))
	logerr(err)
	_, err = q.db.Exec(fmt.Sprintf("DELETE FROM `%s_info` WHERE `id` = %q", Cfg.Prefix, id))
	logerr(err)

	logf("DELETED: %v", id)

	writeHtml(
		q,
		"Corpus verwijderd",
		fmt.Sprintf("Het corpus <em>%s</em> is verwijderd", html.EscapeString(desc)),
		"corpora")
}
