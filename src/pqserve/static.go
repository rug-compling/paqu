package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
)

func contentType(q *Context, ct string) {
	q.w.Header().Set("Content-Type", ct)
}
func cache(q *Context) {
	q.w.Header().Set("Cache-Control", "public, max-age=86400")
}

func static_busy_gif(q *Context) {
	contentType(q, "image/gif")
	cache(q)
	fmt.Fprint(q.w, file__busy__gif)
}

func static_clarinnl_png(q *Context) {
	contentType(q, "image/png")
	cache(q)
	fmt.Fprint(q.w, file__clarinnl__png)
}

func static_favicon_ico(q *Context) {
	contentType(q, "image/x-icon")
	cache(q)
	fmt.Fprint(q.w, file__favicon__ico)
}

func static_info_html(q *Context) {
	writeHead(q, "Info", 4)
	i := strings.Index(file__info__html, "<body>")
	s := file__info__html[i+6:]
	data, err := ioutil.ReadFile(path.Join(paqudir, "contact.html"))
	if err == nil {
		s = strings.Replace(s, "<!--##CONTACT##-->", string(data), 1)
	}
	fmt.Fprint(q.w, s)
}

func static_jquery_js(q *Context) {
	contentType(q, "application/javascript")
	cache(q)
	fmt.Fprint(q.w, file__jquery__js)
}

func static_paqu_css(q *Context) {
	contentType(q, "text/css")
	cache(q)
	fmt.Fprint(q.w, file__paqu__css)
}

func static_leeg_html(q *Context) {
	contentType(q, "text/html")
	cache(q)
	fmt.Fprintln(q.w, "<!DOCTYPE html>\n<html><head><title></title></head><body></body></html>")
}

func static_paqu_png(q *Context) {
	contentType(q, "image/png")
	cache(q)
	fmt.Fprint(q.w, file__paqu__png)
}

func static_robots_txt(q *Context) {
	contentType(q, "text/plain")
	cache(q)
	fmt.Fprint(q.w, file__robots__txt)
}

func static_tooltip_css(q *Context) {
	contentType(q, "text/css")
	cache(q)
	fmt.Fprint(q.w, file__tooltip__css)
}

func static_tooltip_js(q *Context) {
	contentType(q, "application/javascript")
	cache(q)
	fmt.Fprint(q.w, file__tooltip__js)
}
