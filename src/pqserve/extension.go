package main

import (
	"fmt"
)

func extension(q *Context) {
	contentType(q, "text/plain; charset=utf-8")
	nocache(q)
	fmt.Fprint(q.w, `
Dit kan gebruikt worden voor extensies.
Op dit moment zijn er geen extensies gedefinieerd.
`)
}
