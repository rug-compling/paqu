package main

//. Imports

import (
	"github.com/pebbe/util"

	"bytes"
	"crypto/tls"
	"encoding/json"
	"expvar"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type AlpinoInfo struct {
	Limits AlpinoLimits `json:"limits"`
}

type AlpinoLimits struct {
	Jobs int `json:"jobs"`
}

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

	tom := filepath.Join(paquconfigdir, "setup.toml")
	md, err := TomlDecodeFile(tom, &Cfg)
	util.CheckErr(err)
	if un := md.Undecoded(); len(un) > 0 {
		fmt.Fprintln(os.Stderr, "Fout in", tom, ": onbekend :", un)
		return
	}
	if Cfg.Maxjob < 1 {
		Cfg.Maxdup = 1
	}
	if Cfg.Maxspodjob < 1 {
		Cfg.Maxspodjob = 1
	}
	spod_init()

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
		time.Sleep(time.Second)
		close(chLoggerExit)
		wgLogger.Wait()

		os.Exit(0)
	}()

	go foliaclean()

	accessSetup()

	p, err := url.Parse(Cfg.Url)
	util.CheckErr(err)
	cookiepath = p.Path
	if !strings.HasPrefix(cookiepath, "/") {
		cookiepath = "/" + cookiepath
	}

	logf("DbXML: %v", dbxml_version())

	mariadb := false
	db, err := dbopen()
	util.CheckErr(err)
	rows, err := db.Query("SELECT VERSION()")
	util.CheckErr(err)
	if rows.Next() {
		util.CheckErr(rows.Scan(&versionstring))
		rows.Close()
		mariadb = strings.Contains(strings.ToLower(versionstring), "mariadb")
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

	if !mariadb && minversion(5, 7, 7) {
		hasMaxExecutionTime = true
	}
	if mariadb && minversion(10, 1, 2) {
		hasMaxStatementTime = true
	}
	if mariadb {
		logf("MariaDB server-versie: %v (%s)", version, versionstring)
	} else {
		logf("MySQL server-versie: %v (%s)", version, versionstring)
	}

	maxjobs := Cfg.Maxjob
	if Cfg.Alpinoserver != "" {
		buf := bytes.NewBufferString(`{"request":"info"}`)
		resp, err := http.Post(Cfg.Alpinoserver, "application/json", buf)
		util.CheckErr(err)
		data, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		util.CheckErr(err)
		var info AlpinoInfo
		err = json.Unmarshal(data, &info)
		util.CheckErr(err)
		if info.Limits.Jobs > 0 && info.Limits.Jobs < maxjobs {
			maxjobs = info.Limits.Jobs
		}
	}
	semaphore = make(chan struct{}, maxjobs)
	wg.Add(1)
	go func() {
		scheduler()
		wg.Done()
	}()

	go recover()

	go clearMacros()

	handleFunc("", home)
	handleFunc("dl", homedl)
	handleFunc("tree", tree)
	handleFunc("stats", stats)
	handleFunc("statstel", statstel)
	handleFunc("statsmeta", statsmeta)
	handleFunc("statsrel", statsrel)

	handleFunc("logout", logout)
	handleFunc("login", login)
	handleFunc("login1", login1)

	handleFunc("corpuslijst", corpuslijst)
	handleFunc("corsave", corsave)

	handleFunc("browse", browse)
	handleFunc("browserr", browserr)
	handleFunc("browserrud", browserrud)

	handleFunc("xpath", xpath)
	handleFunc("xpathcheck", xpathcheck)
	handleFunc("xstatsmeta", xstatsmeta)
	handleFunc("xpathstats", xpathstats)
	handleFunc("savemacros", savemacros)
	handleFunc("downloadmacros", downloadmacros)
	handleFunc("macroexpand", macroExpand)

	handleFunc("metadata", metadata)
	handleFunc("metadl", metadl)
	handleFunc("meta2", meta2)

	handleFunc("corpora", corpora)
	handleFunc("submitcorpus", submitCorpus)
	handleFunc("share", share)
	handleFunc("share2", share2)
	handleFunc("rename", rename)
	handleFunc("rename2", rename2)
	handleFunc("download", download)
	handleFunc("delete", remove)

	handleFunc("savez", savez)
	handleFunc("savez2", savez2)
	handleFunc("xsavez", xsavez)
	handleFunc("xsavez2", xsavez2)

	handleFunc("info.html", static_info_html)
	handleFunc("relhead.png", static_relhead_png)
	handleFunc("relnone.png", static_relnone_png)
	handleFunc("relother.png", static_relother_png)
	handleFunc("expanded.png", static_expanded_png)

	handleFunc("folia", foliatool)
	handleFunc("foliah0", static_foliahelp0_html)
	handleFunc("foliah1", static_foliahelp1_html)
	handleFunc("foliah2", static_foliahelp2_html)
	handleFunc("foliah3", static_foliahelp3_html)
	handleFunc("foliah4", static_foliahelp4_html)
	handleFunc("foliah5", static_foliahelp5_html)

	handleFunc("spod", spod_main)
	handleFunc("spodform", spod_form)
	handleFunc("spodlist", spod_list)

	handleFunc("ext", extension)

	handleStatic("debug/env", environment)

	handleStatic("busy.gif", static_busy_gif)
	handleStatic("clariah.png", static_clariah_png)
	handleStatic("clarinnl.png", static_clarinnl_png)
	handleStatic("favicon.ico", static_favicon_ico)
	handleStatic("jquery.js", static_jquery_js)
	handleStatic("jquery.textcomplete.js", static_jquery_textcomplete_js)
	handleStatic("leeg.html", static_leeg_html)
	handleStatic("macros.txt", static_macros_txt)
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

	var errserve error
	addr := fmt.Sprint(":", Cfg.Port)
	if !Cfg.Https && !Cfg.Httpdual {
		errserve = http.ListenAndServe(addr, Log(http.DefaultServeMux))
	} else if Cfg.Https && !Cfg.Httpdual {
		errserve = http.ListenAndServeTLS(addr, filepath.Join(paquconfigdir, "cert.pem"), filepath.Join(paquconfigdir, "key.pem"), Log(http.DefaultServeMux))
	} else {

		// De ingewikkelde oplossing: acepteer zowel http als https.
		// Http wordt omgezet in redirect naar https.

		if tlsConfig.NextProtos == nil {
			tlsConfig.NextProtos = []string{"http/1.1"}
		}
		tlsConfig.Certificates = make([]tls.Certificate, 1)
		tlsConfig.Certificates[0], err = tls.LoadX509KeyPair(filepath.Join(paquconfigdir, "cert.pem"), filepath.Join(paquconfigdir, "key.pem"))
		util.CheckErr(err)
		ln, err := net.Listen("tcp", addr)
		util.CheckErr(err)
		errserve = http.Serve(
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
			}))
	}
	logerr(errserve)

	close(chGlobalExit)
	wg.Wait()

	logf("Uptime: %v", time.Now().Sub(started))
	close(chLoggerExit)
	wgLogger.Wait()
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
				if !strings.HasPrefix(r.URL.Path, "/xpathcheck") {
					logf("[%s] %s %s %s", r.Header.Get("X-Forwarded-For"), r.RemoteAddr, r.Method, r.URL)
				}
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
