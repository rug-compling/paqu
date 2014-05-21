package main

import (
	_ "github.com/go-sql-driver/mysql"

	"bytes"
	"database/sql"
	"fmt"
	"html"
	"net/http"
	"runtime"
)

func writeErrHtml(html bool, w http.ResponseWriter, err error, msg ...interface{}) {
	if err == nil {
		return
	}
	if html {
		fmt.Fprintln(w, "<pre>")
	}
	_, filename, lineno, ok := runtime.Caller(1)
	if ok {
		fmt.Fprintf(w, "%v:%v: %v", escape(html, filename), lineno, escape(html, err.Error()))
	} else {
		fmt.Fprint(w, escape(html, err.Error()))
	}
	if len(msg) > 0 {
		fmt.Fprint(w, ",")
		for _, m := range msg {
			fmt.Fprint(w, escape(html, fmt.Sprintf(" %v", m)))
		}
	}
	if html {
		fmt.Fprintln(w, "\n</pre>")
	}
}

func writeErr(w http.ResponseWriter, err error, msg ...interface{}) {
	if err == nil {
		return
	}
	fmt.Fprintln(w, "<pre>")
	_, filename, lineno, ok := runtime.Caller(1)
	if ok {
		fmt.Fprintf(w, "%v:%v: %v", html.EscapeString(filename), lineno, html.EscapeString(err.Error()))
	} else {
		fmt.Fprint(w, html.EscapeString(err.Error()))
	}
	if len(msg) > 0 {
		fmt.Fprint(w, ",")
		for _, m := range msg {
			fmt.Fprint(w, html.EscapeString(fmt.Sprintf(" %v", m)))
		}
	}
	fmt.Fprintln(w, "\n</pre>")
}

func escape(is_html bool, s string) string {
	if is_html {
		return html.EscapeString(s)
	}
	return s
}

func connect() (*sql.DB, error) {
	return sql.Open("mysql", login+"?charset=utf8mb4,utf8&parseTime=true&loc=Europe%2FAmsterdam")
}

func first(r *http.Request, opt string) string {
	if len(r.Form[opt]) > 0 {
		return r.Form[opt][0]
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

func (c Corpora) Len() int {
	return len(c)
}

func (c Corpora) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c Corpora) Less(i, j int) bool {
	return c[i].Title < c[j].Title
}
