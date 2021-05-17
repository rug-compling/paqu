package main

import (
	"github.com/rug-compling/paqu/internal/ranges"

	"github.com/dchest/authcookie"

	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

// The context acts as global store for a single request

type Context struct {
	w            http.ResponseWriter
	r            *http.Request
	user         string
	auth         bool
	sec          string
	quotum       int
	opt_db       []string
	opt_dbmeta   []string
	opt_dbspod   []string
	ignore       map[string]bool
	prefixes     map[string]bool
	spodprefixes map[string]bool
	myprefixes   map[string]bool
	protected    map[string]bool
	hasmeta      map[string]bool
	hasud        map[string]bool
	desc         map[string]string
	lines        map[string]int
	words        map[string]int
	shared       map[string]string
	params       map[string]string
	infos        map[string]string
	infops       map[string]string
	owners       map[string]string
	dates        map[string]time.Time
	form         *multipart.Form
}

// Wrap handler in minimale context, net genoeg voor afhandelen statische pagina's
func handleStatic(url string, handler func(*Context)) {
	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}
	http.HandleFunc(
		url,
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != url && !strings.HasSuffix(url, "/") {
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

// Wrap handler in complete context
func handleFunc(url string, handler func(*Context), options *HandlerOptions) {
	if options == nil {
		options = &HandlerOptions{}
	}
	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}
	http.HandleFunc(
		url,
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != url && (url == "/" || !strings.HasSuffix(url, "/")) {
				http.NotFound(w, r)
				return
			}

			q := &Context{
				w:            w,
				r:            r,
				opt_db:       make([]string, 0),
				opt_dbmeta:   make([]string, 0),
				opt_dbspod:   make([]string, 0),
				prefixes:     make(map[string]bool),
				myprefixes:   make(map[string]bool),
				spodprefixes: make(map[string]bool),
				protected:    make(map[string]bool),
				hasmeta:      make(map[string]bool),
				hasud:        make(map[string]bool),
				desc:         make(map[string]string),
				lines:        make(map[string]int),
				words:        make(map[string]int),
				shared:       make(map[string]string),
				params:       make(map[string]string),
				infos:        make(map[string]string),
				infops:       make(map[string]string),
				owners:       make(map[string]string),
				dates:        make(map[string]time.Time),
			}

			// Is de gebruiker ingelogd?
			if auth, err := r.Cookie("paqu-auth"); err == nil {
				s := strings.SplitN(authcookie.Login(auth.Value, []byte(getRemote(q)+Cfg.Secret)), "|", 2)
				if len(s) == 2 {
					q.user = s[1]
					q.sec = s[0]
				}
			}
			if q.user != "" {
				rows, err := sqlDB.Query(fmt.Sprintf(
					"SELECT `quotum` FROM `%s_users` WHERE `mail` = %q AND `sec` = %q", Cfg.Prefix, q.user, q.sec))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					logerr(err)
					return
				}
				if !rows.Next() {
					q.user = ""
				} else {
					err := rows.Scan(&q.quotum)
					rows.Close()
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						logerr(err)
						return
					}
					q.auth = true
					_, err = sqlDB.Exec(fmt.Sprintf("UPDATE `%s_users` SET `active` = NOW() WHERE `mail` = %q", Cfg.Prefix, q.user))
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						logerr(err)
						return
					}
				}
			}

			// Laad lijsten van corpora

			q.ignore = make(map[string]bool)
			if q.auth {
				rows, err := sqlDB.Query(fmt.Sprintf("SELECT `prefix` FROM `%s_ignore` WHERE `user` = %q", Cfg.Prefix, q.user))
				if err != nil {
					http.Error(q.w, err.Error(), http.StatusInternalServerError)
					logerr(err)
					return
				}
				for rows.Next() {
					var s string
					err := rows.Scan(&s)
					if err != nil {
						http.Error(q.w, err.Error(), http.StatusInternalServerError)
						logerr(err)
						return
					}
					q.ignore[s] = true
				}
			}

			s := "IF(`i`.`owner` = \"none\", \"C\", IF(`i`.`owner` = \"auto\", \"B\",  IF(`i`.`owner` = \"manual\", \"A\", \"Z\")))"
			where := ""
			if q.auth {
				s = fmt.Sprintf("IF(`i`.`owner` = \"none\", \"C\", IF(`i`.`owner` = \"auto\", \"B\", IF(`i`.`owner` = \"manual\", \"A\", IF(`i`.`owner` = %q, \"D\", \"E\"))))", q.user)
				where = fmt.Sprintf(" OR `c`.`user` = %q", q.user)
			}
			rows, err := sqlDB.Query(fmt.Sprintf(
				"SELECT `i`.`id`, `i`.`description`, `i`.`nline`, `i`.`nword`, `i`.`owner`, `i`.`shared`, `i`.`params`,  "+s+", `i`.`protected`, `i`.`hasmeta`,`i`.`hasud`, `i`.`info`, `i`.`infop`, `i`.`created` "+
					"FROM `%s_info` `i`, `%s_corpora` `c` "+
					"WHERE `c`.`enabled` = 1 AND "+
					"`i`.`status` = \"FINISHED\" AND `i`.`id` = `c`.`prefix` AND ( `c`.`user` = \"all\"%s ) "+
					"ORDER BY 8, 2",
				Cfg.Prefix,
				Cfg.Prefix,
				where))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				logerr(err)
				return
			}
			var id, desc, owner, shared, params, info, infop, group string
			var zinnen, woorden, protected, hasmeta, hasud int
			var date time.Time
			for rows.Next() {
				err := rows.Scan(&id, &desc, &zinnen, &woorden, &owner, &shared, &params, &group, &protected, &hasmeta, &hasud, &info, &infop, &date)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					logerr(err)
					return
				}
				if group == "E" {
					if !q.ignore[id] {
						q.opt_db = append(q.opt_db, fmt.Sprintf("E%s %s", id, desc))
						q.prefixes[id] = true
						if hasmeta > 0 {
							q.opt_dbmeta = append(q.opt_dbmeta, fmt.Sprintf("E%s %s", id, desc))
						} else {
							q.opt_dbmeta = append(q.opt_dbmeta, fmt.Sprintf("-E%s %s", id, desc))
						}
						if Cfg.Maxspodlines < 1 || zinnen <= Cfg.Maxspodlines {
							q.opt_dbspod = append(q.opt_dbspod, fmt.Sprintf("E%s %s", id, desc))
							q.spodprefixes[id] = true
						} else {
							q.opt_dbspod = append(q.opt_dbspod, fmt.Sprintf("-E%s %s", id, desc))
						}
						q.owners[id] = owner
					}
				} else if q.auth || owner == "none" || owner == "auto" || owner == "manual" {
					q.opt_db = append(q.opt_db, fmt.Sprintf("%s%s %s", group, id, desc))
					q.prefixes[id] = true
					if hasmeta > 0 {
						q.opt_dbmeta = append(q.opt_dbmeta, fmt.Sprintf("%s%s %s", group, id, desc))
					} else {
						q.opt_dbmeta = append(q.opt_dbmeta, fmt.Sprintf("-%s%s %s", group, id, desc))
					}
					if Cfg.Maxspodlines < 1 || zinnen <= Cfg.Maxspodlines {
						q.opt_dbspod = append(q.opt_dbspod, fmt.Sprintf("%s%s %s", group, id, desc))
						q.spodprefixes[id] = true
					} else {
						q.opt_dbspod = append(q.opt_dbspod, fmt.Sprintf("-%s%s %s", group, id, desc))
					}
				}
				q.desc[id] = desc
				q.lines[id] = zinnen
				q.words[id] = woorden
				q.shared[id] = shared
				q.params[id] = params
				q.infos[id] = info
				q.infops[id] = infop
				q.protected[id] = protected > 0
				if hasmeta > 0 {
					q.hasmeta[id] = true
				}
				if hasud > 0 {
					q.hasud[id] = true
				}
				if q.auth && owner == q.user {
					q.myprefixes[id] = true
				}
				q.dates[id] = date
			}

			if r.Method == "OPTIONS" && options.OptionsMethodHandler != nil {
				options.OptionsMethodHandler(q)
				return
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
				if err != nil && options.NeedForm {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					logerr(err)
					return
				}
				if err == nil {
					q.form, err = reader.ReadForm(10 * 1024 * 1024)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						logerr(err)
						return
					}
				}
			default:
				http.Error(w, "Methode "+r.Method+" is niet toegestaan", http.StatusMethodNotAllowed)
				return
			}

			// Update login-cookies
			setcookie(q)

			handler(q)
		})
}

// Laat niet meer dan een deel van een e-mailadres zien
func displayEmail(s string) string {
	p1 := strings.Index(s, "@")
	p2 := strings.LastIndex(s, ".")
	if p1 < 0 || p2 < 0 {
		return s
	}
	return s[0:p1+1] + ".." + s[p2:len(s)]
}

func datum(t time.Time) string {
	return fmt.Sprintf("%d&nbsp;%s&nbsp;%d", t.Day(), ranges.Maanden[t.Month()], t.Year())
}
