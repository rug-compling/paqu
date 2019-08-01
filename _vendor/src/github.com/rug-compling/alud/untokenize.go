package alud

import (
	"regexp"
	"strings"
)

var (
	reUnQ1 = regexp.MustCompile(`( (?:,,|'')) (.*?) ('' )`)
	reUnQ2 = regexp.MustCompile(`( ") (.*?) (" )`)
	reUnQ3 = regexp.MustCompile("( [`']) (.*?) (' )")
	reUnQ4 = regexp.MustCompile(`( [‘’]) (.*?) (’ )`)
	reUnQ5 = regexp.MustCompile(`( [“„”]) (.*?) (” )`)
	reUnP1 = regexp.MustCompile(`([\[({]) `)
	reUnP2 = regexp.MustCompile(` ([\])}:;.,!?])`)
)

func untokenize(q *context) {
	if t := q.sentence; t != "" {
		t = reUnQ1.ReplaceAllString(" "+t+" ", `$1$2$3`)
		t = reUnQ2.ReplaceAllString(t, `$1$2$3`)
		t = reUnQ3.ReplaceAllString(t, `$1$2$3`)
		t = reUnQ4.ReplaceAllString(t, `$1$2$3`)
		t = reUnQ5.ReplaceAllString(t, `$1$2$3`)
		t = reUnP1.ReplaceAllString(t, `$1`)
		t = reUnP2.ReplaceAllString(t, `$1`)
		q.sentence = strings.TrimSpace(t)
	}
}
