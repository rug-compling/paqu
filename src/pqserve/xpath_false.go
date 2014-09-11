// +build nodbxml

package main

import "net/http"

func xpath(q *Context) {
	http.NotFound(q.w, q.r)
}
