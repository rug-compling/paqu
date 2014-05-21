package main

import (
	_ "github.com/go-sql-driver/mysql"

	"bytes"
	"database/sql"
	"fmt"
	"html"
	"mime/multipart"
	"net/http"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

func doErr(q *Context, err error) bool {
	if err == nil {
		return false
	}

	s := err.Error()

	var s1 string
	_, filename, lineno, ok := runtime.Caller(1)
	if ok {
		s1 = fmt.Sprintf("ERROR: %v:%v: %v", path.Base(filename), lineno, s)
	} else {
		s1 = "ERROR: " + s
	}
	chLog <- s1

	fmt.Fprintln(q.w, "<pre>")
	fmt.Fprintln(q.w, html.EscapeString(s))
	fmt.Fprintln(q.w, "</pre>\n</body>\n</html>")

	return true
}

func writeHead(q *Context, title string) {
	q.w.Header().Set("Content-type", "text/html; charset=utf-8")
	q.w.Header().Set("Cache-Control", "no-cache")
	q.w.Header().Add("Pragma", "no-cache")
	fmt.Fprintf(q.w, `<!DOCTYPE html>
<html>
<head>
<title>%s</title>
<link rel="stylesheet" type="text/css" href="wordrel.css">
<meta name="robots" content="noindex,nofollow">
</head>
<body>
`, title)
}

func writeHtml(q *Context, title, msg, verder string) {
	writeHead(q, title)
	fmt.Fprintf(q.w, `
<h1>%s</h1>
%s
<div class="next">
<a href="%s">Verder</a>
</div>
</body>
</html>
`, title, msg, verder)
}

func first(r *http.Request, opt string) string {
	if len(r.Form[opt]) > 0 {
		return strings.TrimSpace(r.Form[opt][0])
	}
	return ""
}

func firstf(form *multipart.Form, opt string) string {
	if len(form.Value[opt]) > 0 {
		return strings.TrimSpace(form.Value[opt][0])
	}
	return ""
}

func urlencode(s string) string {
	var buf bytes.Buffer
	for _, b := range []byte(s) {
		if b >= 'a' && b <= 'z' ||
			b >= 'A' && b <= 'Z' ||
			b >= '0' && b <= '0' {
			buf.WriteByte(b)
		} else {
			buf.WriteString(fmt.Sprintf("%%%02x", b))
		}
	}
	return buf.String()
}

func iformat(i int) string {
	s1 := fmt.Sprint(i)
	s2 := ""
	for n := len(s1); n > 3; n = len(s1) {
		s2 = "." + s1[n-3:n] + s2
		s1 = s1[0 : n-3]
	}
	return s1 + s2
}

func shell(format string, a ...interface{}) *exec.Cmd {
	cmd := exec.Command(Cfg.Sh, "-c", fmt.Sprintf(format, a...))
	cmd.Env = []string{
		"ALPINO_HOME=" + Cfg.Alpino,
		"PATH=" + Cfg.Path,
		"LANG=en_US.utf8",
		"LANGUAGE=en_US.utf8",
		"LC_ALL=en_US.utf8",
	}
	return cmd
}

func minversion(major, minor, patch int) bool {
	if version[0] > major {
		return true
	}
	if version[0] == major {
		if version[1] > minor {
			return true
		}
		if version[1] == minor {
			if version[2] >= patch {
				return true
			}
		}
	}
	return false
}

func dbopen() (*sql.DB, error) {
	return sql.Open("mysql", Cfg.Login+"?charset=utf8mb4,utf8&parseTime=true&loc=Europe%2FAmsterdam")
}
