package main

import (
	"fmt"
	"net/http"
)

func static_clarinnl_png(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "image/png")
	fmt.Fprint(w, file__clarinnl__png)
}

func static_favicon_ico(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "image/x-icon")
	fmt.Fprint(w, file__favicon__ico)
}

func static_info_html(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	fmt.Fprint(w, file__info__html)
}

func static_jquery_js(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/javascript")
	fmt.Fprint(w, file__jquery__js)
}

func static_paqu_css(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/css")
	fmt.Fprint(w, file__paqu__css)
}

func static_paqu_png(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "image/png")
	fmt.Fprint(w, file__paqu__png)
}

func static_tooltip_css(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/css")
	fmt.Fprint(w, file__tooltip__css)
}

func static_tooltip_js(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/javascript")
	fmt.Fprint(w, file__tooltip__js)
}
