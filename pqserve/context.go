package main

import (
	"github.com/dchest/authcookie"

	"database/sql"
	"fmt"
	"mime/multipart"
	"net/http"
	"path"
	"strings"
)

// The context acts as global store for a single request

type Context struct {
	w          http.ResponseWriter
	r          *http.Request
	user       string
	auth       bool
	quotum     int
	db         *sql.DB
	opt_db     []string
	prefixes   map[string]bool
	myprefixes map[string]bool
	desc       map[string]string
	lines      map[string]int
	shared     map[string]string
	form       *multipart.Form
}

func handleStatic(url string, handler func(*Context)) {
	url = path.Join("/", url)
	http.HandleFunc(
		url,
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != url {
				http.NotFound(w, r)
				return
			}
			q := &Context{
				w: w,
				r: r,
			}
			handler(q)
		})
}

func handleFunc(url string, handler func(*Context)) {
	url = path.Join("/", url)
	http.HandleFunc(
		url,
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != url {
				http.NotFound(w, r)
				return
			}

			q := &Context{
				w:          w,
				r:          r,
				opt_db:     make([]string, 0),
				prefixes:   make(map[string]bool),
				myprefixes: make(map[string]bool),
				desc:       make(map[string]string),
				lines:      make(map[string]int),
				shared:     make(map[string]string),
			}

			// Maak verbinding met database
			var err error
			q.db, err = dbopen()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				logerr(err)
				return
			}
			defer q.db.Close()

			// Is de gebruiker ingelogd?
			if auth, err := r.Cookie("paqu-auth"); err == nil {
				q.user = authcookie.Login(auth.Value, []byte(getRemote(q)+Cfg.Secret))
			}
			if q.user != "" {
				rows, err := q.db.Query(fmt.Sprintf("SELECT SQL_CACHE `quotum` FROM `%s_users` WHERE `mail` = %q", Cfg.Prefix, q.user))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					logerr(err)
					return
				}
				if !rows.Next() {
					q.user = ""
				} else {
					err := rows.Scan(&q.quotum)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						logerr(err)
						return
					}
					rows.Close()
					q.auth = true
					_, err = q.db.Exec(fmt.Sprintf("UPDATE `%s_users` SET `active` = NOW() WHERE `mail` = %q", Cfg.Prefix, q.user))
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						logerr(err)
						return
					}
				}
			}

			// Laad lijsten van corpora
			s := "\"Z\""
			where := ""
			o := "2"
			if q.auth {
				s = fmt.Sprintf("IF(`i`.`owner` = \"none\", \"A\", IF(`i`.`owner` = %q, \"B\", \"C\")) ", q.user)
				where = fmt.Sprintf(" OR `c`.`user` = %q", q.user)
				o = "6, 2"
			}
			rows, err := q.db.Query(fmt.Sprintf(
				"SELECT SQL_CACHE `i`.`id`, `i`.`description`, `i`.`nline`, `i`.`owner`, `i`.`shared`,  "+s+
					"FROM `%s_info` `i`, `%s_corpora` `c` "+
					"WHERE `c`.`enabled` = 1 AND "+
					"`i`.`status` = \"FINISHED\" AND `i`.`id` = `c`.`prefix` AND ( `c`.`user` = \"all\"%s ) "+
					"ORDER BY "+o,
				Cfg.Prefix,
				Cfg.Prefix,
				where))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				logerr(err)
				return
			}
			var id, desc, owner, shared, group string
			var zinnen int
			for rows.Next() {
				err := rows.Scan(&id, &desc, &zinnen, &owner, &shared, &group)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					logerr(err)
					return
				}
				if group == "C" {
					q.opt_db = append(q.opt_db, fmt.Sprintf("C%s %s \u2014 %s \u2014 %s zinnen", id, desc, coded(owner), iformat(zinnen)))
				} else {
					q.opt_db = append(q.opt_db, fmt.Sprintf("%s%s %s \u2014 %s zinnen", group, id, desc, iformat(zinnen)))
				}
				q.prefixes[id] = true
				q.desc[id] = desc
				q.lines[id] = zinnen
				q.shared[id] = shared
				if q.auth && owner == q.user {
					q.myprefixes[id] = true
				}
			}

			// Verwerk input
			switch r.Method {
			case "GET":
				err = r.ParseForm()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					logerr(err)
					return
				}
			case "POST":
				reader, err := r.MultipartReader()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					logerr(err)
					return
				}
				q.form, err = reader.ReadForm(Cfg.Maxmem / 10)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					logerr(err)
					return
				}
			default:
				http.Error(w, "Method "+r.Method+" is not allowd", http.StatusMethodNotAllowed)
				return
			}

			// Update login-cookies
			setcookie(q)

			handler(q)
		})
}

func coded(s string) string {
	p1 := strings.Index(s, "@")
	p2 := strings.LastIndex(s, ".")
	if p1 < 0 || p2 < 0 {
		return s
	}
	return s[0:p1+1] + ".." + s[p2:len(s)]
}
