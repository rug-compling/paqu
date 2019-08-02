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

	// q.ptnodes zijn gesorteerd door functie fixpunct
	nodes := make([]*nodeType, 0)
	words1 := make([]string, 0)
	for _, node := range q.ptnodes {
		if node.End%1000 == 0 {
			nodes = append(nodes, node)
			words1 = append(words1, strings.TrimSpace(node.Word)) // TrimSpace uit voorzorg, zou niet nodig moeten zijn
		}
	}

	sent1 := strings.Join(words1, " ")

	sent2 := reUnQ1.ReplaceAllString(" "+sent1+" ", `$1$2$3`)
	sent2 = reUnQ2.ReplaceAllString(sent2, `$1$2$3`)
	sent2 = reUnQ3.ReplaceAllString(sent2, `$1$2$3`)
	sent2 = reUnQ4.ReplaceAllString(sent2, `$1$2$3`)
	sent2 = reUnQ5.ReplaceAllString(sent2, `$1$2$3`)
	sent2 = reUnP1.ReplaceAllString(sent2, `$1`)
	sent2 = reUnP2.ReplaceAllString(sent2, `$1`)
	sent2 = strings.TrimSpace(sent2)
	q.sentence = sent2

	words2 := strings.Fields(sent2)
	j := 0
	for i, word1 := range words1 {
		if word1 != words2[j] {
			nodes[i].udNoSpaceAfter = true
			words2[j] = words2[j][len(word1):]
		} else {
			j++
		}
	}

}
