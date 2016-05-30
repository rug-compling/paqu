package main

import (
	"fmt"
	"strings"
)

func contentType(q *Context, ct string) {
	q.w.Header().Set("Content-Type", ct)
}

func cache(q *Context) {
	q.w.Header().Set("Cache-Control", "public, max-age=86400")
}

func nocache(q *Context) {
	q.w.Header().Set("Cache-Control", "no-cache")
	q.w.Header().Add("Pragma", "no-cache")
}

func static_busy_gif(q *Context) {
	contentType(q, "image/gif")
	cache(q)
	fmt.Fprint(q.w, file__busy__gif)
}

func static_clariah_png(q *Context) {
	contentType(q, "image/png")
	cache(q)
	fmt.Fprint(q.w, file__clariah__png)
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
	writeHead(q, "Info", 5)
	s := "<div>\n"
	if !has_dbxml {
		s = "<div class=\"nodbxml\">\n"
	}
	i := strings.Index(file__info__html, "<!--##START-->")
	s += file__info__html[i+15:]
	s = strings.Replace(s, "<!--##CONTACT##-->", string(Cfg.Contact), 1)
	fmt.Fprint(q.w, s)
}

func static_foliahelp(q *Context, html, title string) {
	cache(q)
	writeHead(q, "FoLiA Help -- "+title, 0)
	i := strings.Index(html, "<!--##START-->")
	fmt.Fprint(q.w, html[i+15:])
}

func static_foliahelp1_html(q *Context) {
	static_foliahelp(q, file__foliahelp1__html, "Upload FoLiA-bestanden")
}

func static_foliahelp2_html(q *Context) {
	static_foliahelp(q, file__foliahelp2__html, "Upload metadata")
}

func static_foliahelp3_html(q *Context) {
	static_foliahelp(q, file__foliahelp3__html, "Soort invoer")
}

func static_foliahelp4_html(q *Context) {
	static_foliahelp(q, file__foliahelp4__html, "Label voor metadatabestand")
}

func static_foliahelp5_html(q *Context) {
	static_foliahelp(q, file__foliahelp5__html, "Gebruik van metadata")
}

func static_jquery_js(q *Context) {
	contentType(q, "application/javascript")
	cache(q)
	fmt.Fprint(q.w, file__jquery__js)
}

func static_jquery_textcomplete_js(q *Context) {
	contentType(q, "application/javascript")
	cache(q)
	fmt.Fprint(q.w, file__jquery__textcomplete__js)
}

func static_leeg_html(q *Context) {
	contentType(q, "text/html")
	cache(q)
	fmt.Fprintln(q.w, "<!DOCTYPE html>\n<html><head><title></title></head><body></body></html>")
}

func static_paqu_css(q *Context) {
	contentType(q, "text/css")
	cache(q)
	fmt.Fprint(q.w, file__paqu__css)
}

func static_paqu_png(q *Context) {
	contentType(q, "image/png")
	cache(q)
	fmt.Fprint(q.w, file__paqu__png)
}

func static_relhead_png(q *Context) {
	contentType(q, "image/png")
	cache(q)
	fmt.Fprint(q.w, file__relhead__png)
}

func static_relnone_png(q *Context) {
	contentType(q, "image/png")
	cache(q)
	fmt.Fprint(q.w, file__relnone__png)
}

func static_relother_png(q *Context) {
	contentType(q, "image/png")
	cache(q)
	fmt.Fprint(q.w, file__relother__png)
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
