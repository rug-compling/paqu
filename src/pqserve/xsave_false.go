// +build nodbxml

package main

func xsavez(q *Context) {
	http.NotFound(q.w, q.r)
}

func xsavez2(q *Context) {
	http.NotFound(q.w, q.r)
}
