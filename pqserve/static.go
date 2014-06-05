package main

import (
	"fmt"
	"strings"
)

func static_busy_gif(q *Context) {
	q.w.Header().Set("Content-type", "image/gif")
	fmt.Fprint(q.w, file__busy__gif)
}

func static_clarinnl_png(q *Context) {
	q.w.Header().Set("Content-type", "image/png")
	fmt.Fprint(q.w, file__clarinnl__png)
}

func static_favicon_ico(q *Context) {
	q.w.Header().Set("Content-type", "image/x-icon")
	fmt.Fprint(q.w, file__favicon__ico)
}

func static_info_html(q *Context) {
	writeHead(q, "Info", 3)
	i := strings.Index(file__info__html, "<body>")
	fmt.Fprint(q.w, file__info__html[i+6:])
}

func static_jquery_js(q *Context) {
	q.w.Header().Set("Content-type", "application/javascript")
	fmt.Fprint(q.w, file__jquery__js)
}

func static_paqu_css(q *Context) {
	q.w.Header().Set("Content-type", "text/css")
	fmt.Fprint(q.w, file__paqu__css)
}

func static_paqu_png(q *Context) {
	q.w.Header().Set("Content-type", "image/png")
	fmt.Fprint(q.w, file__paqu__png)
}

func static_robots_txt(q *Context) {
	q.w.Header().Set("Content-type", "text/plain")
	fmt.Fprint(q.w, file__robots__txt)
}

func static_tooltip_css(q *Context) {
	q.w.Header().Set("Content-type", "text/css")
	fmt.Fprint(q.w, file__tooltip__css)
}

func static_tooltip_js(q *Context) {
	q.w.Header().Set("Content-type", "application/javascript")
	fmt.Fprint(q.w, file__tooltip__js)
}
