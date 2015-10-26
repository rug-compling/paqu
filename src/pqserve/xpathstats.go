package main

import (
	"fmt"
	"html"
	"net/http"
)

func updateText(q *Context, s string) {
	fmt.Fprintf(q.w, `<script type="text/javascript">
f(%q);
</script>
`, s)
	if ff, ok := q.w.(http.Flusher); ok {
		ff.Flush()
	}
}

func updateError(q *Context, err error, is_html bool) {
	s := err.Error()
	if is_html {
		updateText(q, "Interne fout: "+html.EscapeString(s))
	} else {
		fmt.Fprintln(q.w, "Interne fout:", s)
	}
}

func init() {
	for _, tag := range NodeTags {
		keyTags[tag] = true
	}
}
