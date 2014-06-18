package main

//. Imports

import (
	"github.com/BurntSushi/toml"
	"github.com/pebbe/util"

	"expvar"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func init() {
	expvar.Publish("tasks", ProcessMap(processes))
	expvar.Publish("info", expvar.Func(GetInfo))
}

//. Main

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	for _, arg := range os.Args[1:] {
		if arg == "-v" {
			verbose = true
		}
	}

	paqudir = os.Getenv("PAQU")
	if paqudir == "" {
		paqudir = path.Join(os.Getenv("HOME"), ".paqu")
	}
	_, err := toml.DecodeFile(path.Join(paqudir, "setup.toml"), &Cfg)
	util.CheckErr(err)

	go func() {
		wgLogger.Add(1)
		logger()
		wgLogger.Done()
	}()

	go func() {
		chSignal := make(chan os.Signal, 1)
		signal.Notify(chSignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		sig := <-chSignal
		logf("Signal: %v", sig)

		close(chGlobalExit)
		wg.Wait()

		logf("Uptime: %v", time.Now().Sub(started))
		close(chLoggerExit)
		wgLogger.Wait()

		os.Exit(0)
	}()

	accessSetup()

	p, err := url.Parse(Cfg.Url)
	util.CheckErr(err)
	cookiepath = p.Path
	if !strings.HasPrefix(cookiepath, "/") {
		cookiepath = "/" + cookiepath
	}

	logf("Met DbXML: %v", has_dbxml)

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
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {

				// prioriteit voor chGlobalExit
				select {
				case <-chGlobalExit:
					return
				default:
				}

				select {
				case <-chGlobalExit:
					return
				case task := <-chWork:
					work(task)
				}

			}
		}()
	}
	go recover()

	handleFunc("", home)
	handleFunc("tree", tree)
	handleFunc("stats", stats)
	handleFunc("statsrel", statsrel)

	handleFunc("logout", logout)
	handleFunc("login", login)
	handleFunc("login1", login1)

	handleFunc("corpuslijst", corpuslijst)
	handleFunc("corsave", corsave)

	handleFunc("corpora", corpora)
	handleFunc("submitcorpus", submitCorpus)
	handleFunc("share", share)
	handleFunc("share2", share2)
	handleFunc("rename", rename)
	handleFunc("rename2", rename2)
	handleFunc("download", download)
	handleFunc("delete", remove)

	handleFunc("info.html", static_info_html)

	handleStatic("debug/env", environment)

	handleStatic("busy.gif", static_busy_gif)
	handleStatic("clarinnl.png", static_clarinnl_png)
	handleStatic("favicon.ico", static_favicon_ico)
	handleStatic("jquery.js", static_jquery_js)
	handleStatic("paqu.css", static_paqu_css)
	handleStatic("paqu.png", static_paqu_png)
	handleStatic("robots.txt", static_robots_txt)
	handleStatic("tooltip.css", static_tooltip_css)
	handleStatic("tooltip.js", static_tooltip_js)

	handleStatic("up", up)

	logf("Serving on %s", Cfg.Url)

	logerr(http.ListenAndServe(fmt.Sprint(":", Cfg.Port), Log(http.DefaultServeMux)))
}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-chGlobalExit:
			http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		default:
			wg.Add(1)
			defer wg.Done()
			if accessView(r.RemoteAddr) {
				logf("[%s] %s %s %s", r.Header.Get("X-Forwarded-For"), r.RemoteAddr, r.Method, r.URL)
				handler.ServeHTTP(w, r)
			} else {
				logf("ACCESS DENIED: [%s] %s %s %s", r.Header.Get("X-Forwarded-For"), r.RemoteAddr, r.Method, r.URL)
				http.Error(w, "Access denied", http.StatusForbidden)
			}
		}
	})
}

func up(q *Context) {
	q.w.Header().Set("Content-type", "text/plain")
	q.w.Header().Set("Cache-Control", "no-cache")
	q.w.Header().Add("Pragma", "no-cache")
	fmt.Fprintln(q.w, "up")
}
