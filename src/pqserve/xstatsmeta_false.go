// +build nodbxml

package main

import "net/http"

func xstatsmeta(q *Context) {
	http.NotFound(q.w, q.r)
}
