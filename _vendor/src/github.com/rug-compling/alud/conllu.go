package alud

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

func conll(q *context) string {

	var buf bytes.Buffer

	// TODO: is dit nodig? ook gesorteerd in fixpunct()
	sort.Slice(q.ptnodes, func(i, j int) bool {
		return q.ptnodes[i].End < q.ptnodes[j].End
	})

	fmt.Fprintf(&buf, `# source = %s
# sent_id = %s
# text = %s
# auto = ALUD %d
`,
		q.filename,
		strings.Replace(q.sentid, "/", "\\", -1), // het teken / is gereserveerd
		q.sentence,
		VersionMajor)

	/*
		for i, d := range q.debugs {
			fmt.Fprintf(&buf, "# debug_%d = %s\n", i+1, d)
		}
	*/
	for i, w := range q.warnings {
		fmt.Fprintf(&buf, "# warning_%d = %s\n", i+1, w)
	}

	u := func(s string) string {
		if s == "" {
			return "_"
		}
		return s
	}
	uc := func(i int) string {
		if i > 0 {
			return "CopiedFrom=" + number(i)
		}
		return "_"
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
			uc(node.udCopiedFrom))       // MISC
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
		case error_EXTERNAL_HEAD_MUST_HAVE_ONE_ARG:
			return "ERROR_EXTERNAL_HEAD_MUST_HAVE_ONE_ARG"
		case error_MORE_THAN_ONE_INTERNAL_HEAD_POSITION_FOUND:
			return "ERROR_MORE_THAN_ONE_INTERNAL_HEAD_POSITION_FOUND"
		case error_NO_EXTERNAL_HEAD:
			return "ERROR_NO_EXTERNAL_HEAD"
		case error_NO_HEAD_FOUND:
			return "ERROR_NO_HEAD_FOUND"
		case error_NO_INTERNAL_HEAD:
			return "ERROR_NO_INTERNAL_HEAD"
		case error_NO_INTERNAL_HEAD_IN_GAPPED_CONSTITUENT:
			return "ERROR_NO_INTERNAL_HEAD_IN_GAPPED_CONSTITUENT"
		case error_NO_INTERNAL_HEAD_POSITION_FOUND:
			return "ERROR_NO_INTERNAL_HEAD_POSITION_FOUND"
		case error_NO_VALUE:
			return "ERROR_NO_VALUE"
		case error_RECURSION_LIMIT:
			return "ERROR_RECURSION_LIMIT"
		case underscore:
			return "_"
		case empty_head:
			return "empty head"
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
