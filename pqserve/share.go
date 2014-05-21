package main

import (
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
	q.w.Header().Set("Content-type", "text/html; charset=utf-8")
	q.w.Header().Set("Cache-Control", "no-cache")
	q.w.Header().Add("Pragma", "no-cache")
	fmt.Fprintf(q.w, file__share__html,
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
		writeHtml(q, "Fout", "Waarde voor share is ongeldig of ontbreekt", "corpora")
		return
	}

	_, err := q.db.Exec(fmt.Sprintf("UPDATE `%s_info` SET `shared` = %q WHERE `id` = %q", Cfg.Prefix, sh, id))
	if err != nil {
		http.Error(q.w, err.Error(), http.StatusInternalServerError)
		logerr(err)
		return
	}

	_, err = q.db.Exec(fmt.Sprintf("DELETE FROM `%s_corpora` WHERE `prefix` = %q", Cfg.Prefix, id))
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
	userlist := make([]string, 0)
	fouten := make([]string, 0)
	seen := make(map[string]bool)
	for _, user := range strings.Fields(strings.Replace(firstf(q.form, "mail"), ",", " ", -1)) {
		user = strings.ToLower(user)
		if user == "all" || user == q.user {
			continue
		}
		if seen[user] {
			continue
		}
		seen[user] = true
		rows, err := q.db.Query(fmt.Sprintf("SELECT `mail` FROM `%s_users` WHERE `mail` = %q", Cfg.Prefix, user))
		if err != nil {
			http.Error(q.w, err.Error(), http.StatusInternalServerError)
			logerr(err)
			return
		}
		if rows.Next() {
			rows.Close()
			_, err = q.db.Exec(fmt.Sprintf("INSERT INTO `%s_corpora` (`user`, `prefix`, `enabled`) VALUES (%q, %q, %d)", Cfg.Prefix, user, id, e))
			if err != nil {
				http.Error(q.w, err.Error(), http.StatusInternalServerError)
				logerr(err)
				return
			}
			userlist = append(userlist, user)
		} else {
			fouten = append(fouten, user)
		}
	}

	fout := ""
	if len(fouten) > 0 {
		fout = "\n<p>\nOnbekende gebruiker(s): " + html.EscapeString(strings.Join(fouten, " "))
	}

	var wie string
	switch sh {
	case "PRIVATE":
		wie = "niemand"
	case "SHARED":
		wie = "deze mensen: " + strings.Join(userlist, " ")
	case "PUBLIC":
		wie = "iedereen"
	}

	writeHtml(
		q,
		"Corpus gedeeld",
		fmt.Sprintf("Corpus <b>%s</b> gedeeld met %s%s", html.EscapeString(q.desc[id]), html.EscapeString(wie), fout),
		"corpora")
}
