package main

import (
	"bytes"
	"fmt"
	"html"
	"net/http"
	"strings"
)

func share(q *Context) {

	if !q.auth {
		http.Error(q.w, "Je bent niet ingelogd", http.StatusUnauthorized)
		return
	}

	id := first(q.r, "id")
	if !q.myprefixes[id] {
		http.Error(q.w, "Dat is niet je corpus", http.StatusUnauthorized)
		return
	}

	if q.protected[id] {
		http.Error(q.w, "Dat is afgeleid van een corpus dat niet van jou is", http.StatusUnauthorized)
		return
	}

	/*
		parameters:
		- corpusomschrijving, htmlescaped
		- aantal zinnen, als string met punten
		- corpus-id
		- radio none
		- radio some
		- textarea, htmlescaped
		- radio, all

		radio, leeg of: " checked="checked
	*/

	mails := make([]string, 0)
	rows, err := q.db.Query(fmt.Sprintf(
		"SELECT `user` FROM `%s_corpora` WHERE `prefix` = %q AND `user` != \"all\" AND `user` != %q ORDER BY `user`", Cfg.Prefix, id, q.user))
	if err != nil {
		http.Error(q.w, err.Error(), http.StatusInternalServerError)
		logerr(err)
		return
	}
	for rows.Next() {
		var mail string
		err := rows.Scan(&mail)
		if err != nil {
			http.Error(q.w, err.Error(), http.StatusInternalServerError)
			logerr(err)
			return
		}
		mails = append(mails, mail)
	}
	users := strings.Join(mails, "\n")

	checked := `" checked="checked`
	checked1 := ""
	checked2 := ""
	checked3 := ""
	switch q.shared[id] {
	case "PRIVATE":
		checked1 = checked
	case "SHARED":
		checked2 = checked
	case "PUBLIC":
		checked3 = checked
	}
	writeHead(q, "Corpus delen", 0)
	fmt.Fprintf(q.w, `<h1>Corpus delen</h1>
    Corpus:
    <blockquote>
      <b>%s</b><br>
      %s zinnen
    </blockquote>
    Delen met:
    <blockquote>
      <form action="share2" method="post" enctype="multipart/form-data">
        <input type="hidden" name="corpus" value="%s">
        <input type="radio" name="share" value="none%s">niemand
        <p>
          <input type="radio" name="share" value="some%s">deze mensen (e-mailadressen):
          <br>
          <textarea id="mail" style="margin-left:4em;" name="mail" rows="8" cols="40">%s</textarea>
        <p>
          <input type="radio" name="share" value="all%s">iedereen die is ingelogd
        <p>
          <input type="submit">
      </form>
    </blockquote>
  </body>
<html>
`,
		html.EscapeString(q.desc[id]),
		iformat(q.lines[id]),
		id,
		checked1,
		checked2,
		html.EscapeString(users),
		checked3)
}

func share2(q *Context) {

	if !q.auth {
		http.Error(q.w, "Je bent niet ingelogd", http.StatusUnauthorized)
		return
	}

	id := firstf(q.form, "corpus")
	if !q.myprefixes[id] {
		http.Error(q.w, "Dat is niet je corpus", http.StatusUnauthorized)
		return
	}

	var sh string
	switch firstf(q.form, "share") {
	case "none":
		sh = "PRIVATE"
	case "some":
		sh = "SHARED"
	case "all":
		sh = "PUBLIC"
	default:
		writeHtml(q, "Fout", "Waarde voor share is ongeldig of ontbreekt")
		return
	}

	_, err := q.db.Exec(fmt.Sprintf("DELETE FROM `%s_corpora` WHERE `prefix` = %q", Cfg.Prefix, id))
	if err != nil {
		http.Error(q.w, err.Error(), http.StatusInternalServerError)
		logerr(err)
		return
	}

	u := q.user
	if sh == "PUBLIC" {
		u = "all"
	}
	_, err = q.db.Exec(fmt.Sprintf("INSERT INTO `%s_corpora` (`user`, `prefix`) VALUES (%q, %q)", Cfg.Prefix, u, id))
	if err != nil {
		http.Error(q.w, err.Error(), http.StatusInternalServerError)
		logerr(err)
		return
	}

	e := 0
	if sh == "SHARED" {
		e = 1
	}
	fouten := make([]string, 0)
	seen := make(map[string]bool)
	count := 0
	for _, user := range strings.Fields(strings.Replace(firstf(q.form, "mail"), ",", " ", -1)) {
		user = strings.ToLower(user)
		if user == "all" || user == q.user {
			continue
		}
		if seen[user] {
			continue
		}
		count++
		seen[user] = true
		rows, err := q.db.Query(fmt.Sprintf("SELECT `mail` FROM `%s_users` WHERE `mail` = %q", Cfg.Prefix, user))
		if err != nil {
			http.Error(q.w, err.Error(), http.StatusInternalServerError)
			logerr(err)
			return
		}
		if rows.Next() {
			rows.Close()
		} else {
			fouten = append(fouten, user)
		}
		_, err = q.db.Exec(fmt.Sprintf("INSERT INTO `%s_corpora` (`user`, `prefix`, `enabled`) VALUES (%q, %q, %d)", Cfg.Prefix, user, id, e))
		if err != nil {
			http.Error(q.w, err.Error(), http.StatusInternalServerError)
			logerr(err)
			return
		}
	}

	if sh == "SHARED" && count == 0 {
		sh = "PRIVATE"
	}

	_, err = q.db.Exec(fmt.Sprintf("UPDATE `%s_info` SET `shared` = %q WHERE `id` = %q", Cfg.Prefix, sh, id))
	if err != nil {
		http.Error(q.w, err.Error(), http.StatusInternalServerError)
		logerr(err)
		return
	}

	var wie string
	switch sh {
	case "PRIVATE":
		wie = "niemand"
	case "SHARED":
		s := "personen"
		if count == 1 {
			s = "persoon"
		}
		wie = fmt.Sprintf("%d %s", count, s)
	case "PUBLIC":
		wie = "iedereen"
	}

	var buf bytes.Buffer

	fmt.Fprintf(&buf, "Corpus <b>%s</b> gedeeld met %s", html.EscapeString(q.desc[id]), html.EscapeString(wie))

	if len(fouten) > 0 {
		s := "Onbekende personen"
		if len(fouten) == 1 {
			s = "Onbekend persoon"
		}
		fmt.Fprintf(&buf, "<div class=\"warning\">\n%s:\n<ul>\n", s)
		for _, f := range fouten {
			fmt.Fprintf(&buf, "<li>%s\n", html.EscapeString(f))
		}
		fmt.Fprint(&buf, "</ul>\n</div>\n")
	}

	writeHtml(
		q,
		"Corpus gedeeld",
		buf.String())
}
