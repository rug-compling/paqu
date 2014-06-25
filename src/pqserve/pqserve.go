package main

//. Imports

import (
	"github.com/BurntSushi/toml"
	"github.com/pebbe/util"

	"crypto/tls"
	"expvar"
	"fmt"
	"io"
	"net"
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
	logf("MySQL server-versie: %v (%s)", version, versionstring)

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

	var s string
	if Cfg.Https || Cfg.Httpdual {
		s = "s"
	}
	logf("Server beschikbaar op http%s://127.0.0.1:%v", s, Cfg.Port)

	addr := fmt.Sprint(":", Cfg.Port)
	if Cfg.Https || Cfg.Httpdual {

		if !Cfg.Httpdual {

			// De simpele oplossing: accepteer alleen https.

			logerr(http.ListenAndServeTLS(addr, path.Join(paqudir, "cert.pem"), path.Join(paqudir, "key.pem"), Log(http.DefaultServeMux)))
			return

		}

		// De ingewikkelde oplossing: acepteer zowel http als https.
		// Http wordt omgezet in redirect naar https.

		if tlsConfig.NextProtos == nil {
			tlsConfig.NextProtos = []string{"http/1.1"}
		}
		tlsConfig.Certificates = make([]tls.Certificate, 1)
		tlsConfig.Certificates[0], err = tls.LoadX509KeyPair(path.Join(paqudir, "cert.pem"), path.Join(paqudir, "key.pem"))
		util.CheckErr(err)
		ln, err := net.Listen("tcp", addr)
		util.CheckErr(err)
		logerr(http.Serve(
			&SplitListener{Listener: ln},
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.TLS == nil {
					u := url.URL{
						Scheme:   "https",
						Host:     r.Host,
						Path:     r.URL.Path,
						RawQuery: r.URL.RawQuery,
						Fragment: r.URL.Fragment,
					}
					http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
				} else {
					Log(http.DefaultServeMux).ServeHTTP(w, r)
				}
			})))
	} else {
		logerr(http.ListenAndServe(addr, Log(http.DefaultServeMux)))
	}
}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-chGlobalExit:
			http.Error(w, "Service niet beschikbaar", http.StatusServiceUnavailable)
		default:
			wg.Add(1)
			defer wg.Done()
			if accessView(r.RemoteAddr) {
				logf("[%s] %s %s %s", r.Header.Get("X-Forwarded-For"), r.RemoteAddr, r.Method, r.URL)
				handler.ServeHTTP(w, r)
			} else {
				logf("GEEN TOEGANG: [%s] %s %s %s", r.Header.Get("X-Forwarded-For"), r.RemoteAddr, r.Method, r.URL)
				http.Error(w, "Geen toegang", http.StatusForbidden)
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

type Conn struct {
	net.Conn
	b byte
	e error
	f bool
}

func (c *Conn) Read(b []byte) (int, error) {
	if c.f {
		c.f = false
		b[0] = c.b
		if len(b) > 1 && c.e == nil {
			n, e := c.Conn.Read(b[1:])
			if e != nil {
				c.Conn.Close()
			}
			return n + 1, e
		} else {
			return 1, c.e
		}
	}
	return c.Conn.Read(b)
}

type SplitListener struct {
	net.Listener
}

func (l *SplitListener) Accept() (net.Conn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	if tc, ok := c.(*net.TCPConn); ok {
		tc.SetKeepAlive(true)
		tc.SetKeepAlivePeriod(3 * time.Minute)
	}
	b := make([]byte, 1)
	_, err = c.Read(b)
	if err != nil {
		c.Close()
		if err != io.EOF {
			return nil, err
		}
	}
	con := &Conn{Conn: c, b: b[0], e: err, f: true}
	if b[0] == 22 {
		return tls.Server(con, tlsConfig), nil
	}
	return con, nil
}