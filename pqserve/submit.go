package main

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

var (
	tr = map[string]string{
		"FINISHED": "gereed",
		"QUEUED":   "wachtrij",
		"WORKING":  "bezig",
		"FAILED":   "fout",
		"PRIVATE":  "priv&eacute;",
		"PUBLIC":   "openbaar",
		"SHARED":   "gedeeld",
	}
)

func submit(q *Context) {

	if !q.auth {
		http.Error(q.w, "Je bent niet ingelogd", http.StatusUnauthorized)
		return
	}

	writeHead(q, "Mijn corpora")
	fmt.Fprint(q.w, `
<div class="home">
<a href=".">Start</a>
</div>
<h1>Mijn corpora</h1>
`)

	rows, err := q.db.Query(
		fmt.Sprintf(
			"SELECT `id`, `description`, `status`, `nline`, `nword`, `msg`, `shared` FROM `%s_info` WHERE `owner` = \"%s\" ORDER BY `description`",
			Cfg.Prefix,
			q.user))
	if doErr(q, err) {
		return
	}

	n := 0
	gebruikt := 0
	var id, desc, status, msg, shared string
	var zinnen, woorden int
	for rows.Next() {
		err := rows.Scan(&id, &desc, &status, &zinnen, &woorden, &msg, &shared)
		if err != nil {
			if n > 0 {
				fmt.Fprintln(q.w, "</table>")
			}
			doErr(q, err)
			return
		}
		gebruikt += woorden
		n++
		if n == 1 {
			fmt.Fprintln(q.w, "<table border=\"1\" cellpadding=\"4\">")
		}
		class := ""
		if status == "FINISHED" {
			class = " class=\"ok\""
		} else if status == "FAILED" {
			class = " class=\"error\""
		}
		fmt.Fprintf(q.w, "<tr valign=\"top\"><td%s>%s\n", class, vertaal(status))
		if status == "FINISHED" {
			fmt.Fprintf(q.w, "<td><a href=\"share?id=%s\">%s</a>\n", id, vertaal(shared))
		} else {
			fmt.Fprint(q.w, "<td>\n")
		}
		fmt.Fprintf(q.w, "<td>%s\n", html.EscapeString(desc))
		if status == "FINISHED" {
			fmt.Fprintf(q.w, "<td class=\"right\">%s zinnen\n", iformat(zinnen))
		} else {
			fmt.Fprintln(q.w, "<td>")
		}
		if status == "FINISHED" {
			fmt.Fprintf(q.w, "<td><a href=\"rename?id=%s\">hernoemen</a>\n", id)
		} else {
			fmt.Fprint(q.w, "<td>\n")
		}
		fmt.Fprintf(q.w, "<td><a href=\"delete?id=%s\">verwijderen</a>\n", id)
		fmt.Fprintf(q.w, "<td>%s\n", html.EscapeString(msg))
	}

	if n == 0 {
		fmt.Fprintln(q.w, "Je hebt nog geen corpora")
	} else {
		fmt.Fprintln(q.w, "</table>")
	}

	fmt.Fprintln(q.w, "<h2>Nieuw corpus maken</h2>")
	if q.quotum > 0 {
		fmt.Fprintf(q.w, "Je hebt nog ruimte voor %d woorden (tokens)\n<p>\n", q.quotum - gebruikt)
	}
	fmt.Fprint(q.w, `
    <form action="submit" method="post" enctype="multipart/form-data">
        De tekst die je uploadt moet platte tekst zijn, zonder opmaak (geen Word of zo), gecodeerd in utf-8.
        <p>
    Titel:<br>
	<input type="text" name="title">
    <p>
	Upload document:<br>
	<input type="file" name="data">
        <p>
        Structuur van document:<br>
	<select name="how">
	  <option value="run">Doorlopende tekst</option>
	  <option value="line">Een zin per regel</option>
	</select>
      <p>
	<input type="submit">
    </form>
</body>
</html>
`)

}

func submit2(q *Context) {

	if !q.auth {
		http.Error(q.w, "Je bent niet ingelogd", http.StatusUnauthorized)
		return
	}

	how := firstf(q.form, "how")
	title := firstf(q.form, "title")

	if title == "" {
		http.Error(q.w, "Titel ontbreekt", http.StatusPreconditionFailed)
		return
	}

	dirname := reNoAz.ReplaceAllString(strings.ToLower(title), "")
	if len(dirname) > 20 {
		dirname = dirname[:20]
	} else if dirname == "" {
		dirname = "a"
	}

	dirnameLock.Lock()
	defer dirnameLock.Unlock()
	for i := 0; true; i++ {
		d := dirname + abc(i)
		rows, err := q.db.Query(fmt.Sprintf("SELECT 1 FROM `%s_info` WHERE `id` = %q", Cfg.Prefix, d))
		if err != nil {
			http.Error(q.w, err.Error(), http.StatusInternalServerError)
			logerr(err)
			return
		}
		if rows.Next() {
			rows.Close()
			continue
		}
		dirname = d
		break
	}
	fulldirname := path.Join(paqudir, "data", dirname)
	err := os.Mkdir(fulldirname, 0700)
	if err != nil {
		http.Error(q.w, err.Error(), http.StatusInternalServerError)
		logerr(err)
		return
	}

	if len(q.form.File["data"]) > 0 {
		fpout, err := os.Create(path.Join(fulldirname, "data"))
		if err != nil {
			http.Error(q.w, err.Error(), http.StatusInternalServerError)
			logerr(err)
			return
		}
		defer fpout.Close()
		fpin, err := q.form.File["data"][0].Open()
		if err != nil {
			http.Error(q.w, err.Error(), http.StatusInternalServerError)
			logerr(err)
			return
		}
		defer fpin.Close()
		_, err = io.Copy(fpout, fpin)
		if err != nil {
			http.Error(q.w, err.Error(), http.StatusInternalServerError)
			logerr(err)
			return
		}
	}

	_, err = q.db.Exec(fmt.Sprintf("INSERT %s_info (id, description, owner, status, params) VALUES (%q, %q, %q, \"QUEUED\", %q);",
		Cfg.Prefix,
		dirname, title, q.user, how))
	if err != nil {
		http.Error(q.w, err.Error(), http.StatusInternalServerError)
		logerr(err)
		return
	}

	logf("QUEUED: " + dirname)
	processes[dirname] = &Process{
		id:     dirname,
		chKill: make(chan bool, 10),
		queued: true,
	}
	go func() {
		chWork <- processes[dirname]
	}()

	writeHtml(
		q,
		"Document word verwerkt",
		`
Je document wordt verwerkt. Als het klaar is zie je op de hoofdpagina een nieuw corpus bij de databases.
<p>
Let op: Dit kan even duren. Minuten, uren, of dagen, afhankelijk van de grootte van je document.
<p>
<b>Je krijgt een e-mail als het corpus klaar is.</b>
`,
		"corpora")
}

func abc(i int) string {
	b := make([]byte, 0)
	for {
		j := i % 26
		i = i / 26
		b = append(b, byte(j+97))
		if i == 0 {
			break
		}
	}
	return string(b)
}

func vertaal(s string) string {
	if s1, ok := tr[s]; ok {
		return s1
	}
	return s
}
