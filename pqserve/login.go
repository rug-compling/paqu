package main

import (
	"github.com/dchest/authcookie"

	"fmt"
	"html"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func login(q *Context) {

	mail := first(q.r, "mail")
	pw := first(q.r, "pw")
	sec := ""

	if pw == "" {
		pw = "none" // anders kan iemand na een eerdere inlog zonder password inloggen
	}

	rows, err := q.db.Query(fmt.Sprintf("SELECT `sec` FROM `%s_users` WHERE `mail` = %q AND `pw` = %q", Cfg.Prefix, mail, pw))
	if err != nil {
		http.Error(q.w, err.Error(), http.StatusInternalServerError)
		logerr(err)
		return
	}

	if rows.Next() {
		err := rows.Scan(&sec)
		if err != nil {
			http.Error(q.w, err.Error(), http.StatusInternalServerError)
			logerr(err)
			return
		}
		rows.Close()
		_, err = q.db.Exec(fmt.Sprintf("UPDATE `%s_users` SET `pw` = '' WHERE `mail` = %q", Cfg.Prefix, mail))
		if err != nil {
			http.Error(q.w, err.Error(), http.StatusInternalServerError)
			logerr(err)
			return
		}
		q.auth = true
		q.user = mail
		q.sec = sec
		setcookie(q)
		writeHtml(q, "OK", "Je bent ingelogd")
	} else {
		writeHtml(q, "Fout", "Log-in mislukt")
	}
}

func login1(q *Context) {

	mail := strings.ToLower(first(q.r, "mail"))

	if !accessLogin(mail) {
		logf("LOGIN DENIED: %s %s %s", q.r.RemoteAddr, q.r.Method, q.r.URL)
		http.Error(q.w, "Access denied", http.StatusForbidden)
		return
	}

	if mail == "" {
		writeHtml(q, "Fout", "E-mailadres ontbreekt")
		return
	}
	if !reMail.MatchString(mail) {
		writeHtml(q, "Fout", "Dat ziet er niet uit als een geldig e-mailadress")
		return
	}

	auth := rand16()
	sec := rand16()

	rows, err := q.db.Query(fmt.Sprintf("SELECT * from `%s_users` WHERE `mail` = %q", Cfg.Prefix, mail))
	if err != nil {
		http.Error(q.w, err.Error(), http.StatusInternalServerError)
		logerr(err)
		return
	}
	if rows.Next() {
		rows.Close()
		_, err = q.db.Exec(fmt.Sprintf("UPDATE `%s_users` SET `pw` = %q, `sec` = %q WHERE `mail` = %q", Cfg.Prefix, auth, sec, mail))
	} else {
		_, err = q.db.Exec(fmt.Sprintf("INSERT INTO `%s_users` (`mail`, `pw`, `sec`, `quotum`) VALUES (%q, %q, %q, %d)", Cfg.Prefix, mail, auth, sec, Cfg.Maxwrd))
	}
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
			urlJoin(Cfg.Url, "login"), urlencode(mail), urlencode(auth)))
	if err != nil {
		http.Error(q.w, err.Error(), http.StatusInternalServerError)
		logerr(err)
		return
	}

	writeHtml(
		q,
		"Mail verzonden",
		fmt.Sprintf("Een bericht is verstuurd naar %s. Ga naar de link in dat bericht om in te loggen.", html.EscapeString(mail)))
}

func logout(q *Context) {
	if q.auth {
		q.db.Exec(fmt.Sprintf("UPDATE `%s_users` SET `sec` = \"x\" WHERE `mail` = %q", Cfg.Prefix, q.user))
		q.auth = false
	}
	http.SetCookie(q.w, &http.Cookie{Name: "paqu-auth", Path: "/", MaxAge: -1})
	writeHtml(q, "Uitgelogd", "Je bent uitgelogd")
}

func setcookie(q *Context) {
	if q.auth {
		exp := time.Now().AddDate(0, 0, 14)
		au := authcookie.New(q.sec+"|"+q.user, exp, []byte(getRemote(q)+Cfg.Secret))
		http.SetCookie(q.w, &http.Cookie{Name: "paqu-auth", Value: au, Path: cookiepath, Expires: exp})
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

func rand16() string {
	a := make([]byte, 16)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 16; i++ {
		a[i] = byte(97 + rnd.Intn(26))
	}
	return string(a)
}
