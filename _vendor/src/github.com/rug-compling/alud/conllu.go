package alud

import (
	"bytes"
	"fmt"
	"strings"
)

func conll(q *context, options int) string {

	var buf bytes.Buffer

	if options&OPT_NO_COMMENTS == 0 {
		fmt.Fprintf(&buf, `# source = %s
# sent_id = %s
# text = %s
# auto = %s
`,
			q.filename,
			strings.Replace(q.sentid, "/", "\\", -1), // het teken / is gereserveerd
			q.sentence,
			versionID)

		if options&OPT_DEBUG != 0 {
			for i, d := range q.debugs {
				fmt.Fprintf(&buf, "# debug_%d = %s\n", i+1, d)
			}
		}
	}

	u := func(s string) string {
		if s == "" {
			return "_"
		}
		return s
	}
	misc := func(node *nodeType) string {
		ss := make([]string, 0, 2)
		if node.udCopiedFrom > 0 {
			ss = append(ss, "CopiedFrom="+number(node.udCopiedFrom))
		}
		if node.udNoSpaceAfter {
			ss = append(ss, "SpaceAfter=No")
		}
		return strings.Join(ss, "|")
	}
	postag := func(s string) string {
		return strings.Join(strings.FieldsFunc(s, func(r rune) bool {
			return r == '(' || r == ')' || r == ','
		}), "|")
	}

	for _, node := range q.ptnodes {
		fmt.Fprintf(&buf, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			number(node.End),            // ID
			node.Word,                   // FORM
			node.Lemma,                  // LEMMA
			u(node.udPos),               // UPOS
			u(postag(node.Postag)),      // XPOS
			u(featuresToString(node)),   // FEATS
			number(node.udHeadPosition), // HEAD
			u(node.udRelation),          // DEPREL
			u(node.udEnhanced),          // DEPS
			u(misc(node)))               // MISC
	}

	fmt.Fprintln(&buf)

	return buf.String()
}

func featuresToString(node *nodeType) string {
	features := make([]string, 0)
	for _, f := range [][2]string{
		[2]string{node.udAbbr, "Abbr"},
		[2]string{node.udCase, "Case"},
		[2]string{node.udDefinite, "Definite"},
		[2]string{node.udDegree, "Degree"},
		[2]string{node.udForeign, "Foreign"},
		[2]string{node.udGender, "Gender"},
		[2]string{node.udNumber, "Number"},
		[2]string{node.udPerson, "Person"},
		[2]string{node.udPoss, "Poss"},
		[2]string{node.udPronType, "PronType"},
		[2]string{node.udReflex, "Reflex"},
		[2]string{node.udTense, "Tense"},
		[2]string{node.udVerbForm, "VerbForm"},
	} {
		if f[0] != "" {
			features = append(features, f[1]+"="+f[0])
		}
	}
	return strings.Join(features, "|")
}

func number(n int) string {
	if n < 0 {
		switch n {
		case underscore:
			return "_"
		case empty_head:
			return "empty head"
		case error_no_head:
			panic("No head")
		default:
			panic("Missing case")
		}
	}
	i, j := n/1000, n%1000
	if j == 0 {
		return fmt.Sprint(i)
	}
	return fmt.Sprintf("%d.%d", i, j)
}
