package main

//. Imports

import (
	"github.com/BurntSushi/toml"
	"github.com/pebbe/util"

	"fmt"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strconv"
)

//. Main

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	for len(os.Args) > 1 {
		if os.Args[1] == "-v" {
			verbose = true
		} else {
			break
		}
		os.Args = append(os.Args[:1], os.Args[2:]...)
	}

	_, err := toml.DecodeFile(os.Args[1], &Cfg)
	util.CheckErr(err)

	go logger()

	db, err := dbopen()
	util.CheckErr(err)
	rows, err := db.Query("SELECT `VARIABLE_VALUE` FROM `information_schema`.`GLOBAL_VARIABLES` WHERE `VARIABLE_NAME` = \"VERSION\"")
	util.CheckErr(err)
	if rows.Next() {
		util.CheckErr(rows.Scan(&versionstring))
		rows.Close()
		r := regexp.MustCompile("[0-9]+").FindAllString(versionstring, 3)
		if r != nil {
			if len(r) > 0 {
				version[0], _ = strconv.Atoi(r[0])
			}
			if len(r) > 1 {
				version[1], _ = strconv.Atoi(r[1])
			}
			if len(r) > 2 {
				version[2], _ = strconv.Atoi(r[2])
			}
		}
	}
	db.Close()
	if minversion(5, 7, 4) {
		hasMaxStatementTime = true
	}
	logf("MySQL server version: %v (%s)", version, versionstring)

	for i := 0; i < Cfg.Maxjob; i++ {
		go func() {
			for task := range chWork {
				work(task)
			}
		}()
	}
	recover()

	handleFunc("", form)
	handleFunc("tree", tree)
	handleFunc("stats", stats)

	handleFunc("logout", logout)
	handleFunc("login", login)
	handleFunc("login1", login1)

	handleFunc("corpora", submit)
	handleFunc("submit", submit2)
	handleFunc("share", share)
	handleFunc("share2", share2)
	handleFunc("rename", rename)
	handleFunc("rename2", rename2)
	handleFunc("delete", remove)
	handleFunc("delete2", remove2)

	handleStatic("info.html", static_info_html)
	handleStatic("busy.gif", static_busy_gif)
	handleStatic("clarinnl.png", static_clarinnl_png)
	handleStatic("favicon.ico", static_favicon_ico)
	handleStatic("jquery.js", static_jquery_js)
	handleStatic("robots.txt", static_robots_txt)
	handleStatic("wordrel.css", static_wordrel_css)
	handleStatic("tooltip.css", static_tooltip_css)
	handleStatic("tooltip.js", static_tooltip_js)

	handleStatic("up", up)

	logf("Serving on %s", Cfg.Url)

	logerr(http.ListenAndServe(fmt.Sprint(":", Cfg.Port), Log(http.DefaultServeMux)))
}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func up(q *Context) {
	q.w.Header().Set("Content-type", "text/plain")
	q.w.Header().Set("Cache-Control", "no-cache")
	q.w.Header().Add("Pragma", "no-cache")
	fmt.Fprintln(q.w, "up")
}
