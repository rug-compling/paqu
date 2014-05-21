package main

import (
	"github.com/dchest/authcookie"

	"fmt"
	"html"
	"math/rand"
	"net/http"
	"path"
	"strings"
	"time"
)

func login(q *Context) {

	mail := first(q.r, "mail")
	pw := first(q.r, "pw")

	rows, err := q.db.Query(fmt.Sprintf("SELECT 1 FROM `%s_users` WHERE `mail` = %q AND `pw` = %q", Cfg.Prefix, mail, pw))
	if err != nil {
		http.Error(q.w, err.Error(), http.StatusInternalServerError)
		logerr(err)
		return
	}

	if rows.Next() {
		rows.Close()
		_, err := q.db.Exec(fmt.Sprintf("UPDATE `%s_users` SET `pw` = '' WHERE `mail` = %q", Cfg.Prefix, mail))
		if err != nil {
			http.Error(q.w, err.Error(), http.StatusInternalServerError)
			logerr(err)
			return
		}
		q.auth = true
		q.user = mail
		setcookie(q)
		writeHtml(q, "OK", "Je bent ingelogd", ".")
	} else {
		writeHtml(q, "Fout", "Log-in mislukt", ".")
	}
}

func login1(q *Context) {

	mail := strings.ToLower(first(q.r, "mail"))
	if mail == "" {
		writeHtml(q, "Fout", "E-mailadres ontbreeekt", ".")
		return
	}
	if !reMail.MatchString(mail) {
		writeHtml(q, "Fout", "Dat ziet er niet uit als een geldig e-mailadress", ".")
		return
	}

	_, err := q.db.Exec(fmt.Sprintf("DELETE FROM `%s_users` WHERE `mail` = %q", Cfg.Prefix, mail))
	if err != nil {
		http.Error(q.w, err.Error(), http.StatusInternalServerError)
		logerr(err)
		return
	}

	a := make([]byte, 16)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 16; i++ {
		a[i] = byte(97 + rnd.Intn(24))
	}
	auth := string(a)
	_, err = q.db.Exec(fmt.Sprintf("INSERT INTO `%s_users` (`mail`, `pw`) VALUES (%q, %q)", Cfg.Prefix, mail, auth))
	if err != nil {
		http.Error(q.w, err.Error(), http.StatusInternalServerError)
		logerr(err)
		return
	}

	err = sendmail(
		mail,
		"Log in",
		fmt.Sprintf(
			"Visit this URL to log in: %s?mail=%s&pw=%s",
			path.Join(Cfg.Url, "login"), urlencode(mail), urlencode(auth)))
	if err != nil {
		http.Error(q.w, err.Error(), http.StatusInternalServerError)
		logerr(err)
		return
	}

	writeHtml(
		q,
		"Mail verzonden",
		fmt.Sprintf("Een bericht is verstuurd naar %s. Ga naar de link in dat bericht om in te loggen.", html.EscapeString(mail)),
		".")
}

func logout(q *Context) {
	http.SetCookie(q.w, &http.Cookie{Name: "wordrel-auth", Path: "/", MaxAge: -1})
	http.SetCookie(q.w, &http.Cookie{Name: "wordrel-mail", Path: "/", MaxAge: -1})
	writeHtml(q, "Uitgelogd", "Je bent uitgelogd", ".")
}

func setcookie(q *Context) {
	if q.auth {
		exp := time.Now().AddDate(0, 0, 14)
		au := authcookie.New(q.user, exp, []byte(Cfg.Secret+q.user+getRemote(q)))
		http.SetCookie(q.w, &http.Cookie{Name: "wordrel-auth", Value: au, Path: Cfg.Cookiepath, Expires: exp})
		http.SetCookie(q.w, &http.Cookie{Name: "wordrel-mail", Value: q.user, Path: Cfg.Cookiepath, Expires: exp})
	}
}

func getRemote(q *Context) string {
	/* TODO
	if Cfg.Forwarded {
		return ...
	}
	*/
	if Cfg.Remote {
		return strings.Split(q.r.RemoteAddr, ":")[0]
	}
	return ""
}
