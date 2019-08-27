package alud

import (
	"fmt"
)

func addPosTags(q *context) {
	for _, node := range q.ptnodes {
		node.udPos = universalPosTags(node, q)
	}
}

func universalPosTags(node *nodeType, q *context) string {
	pt := node.Pt
	rel := node.Rel

	if pt == "let" {
		if rel == "--" {
			for _, n := range node.parent.Node {
				if n.Pt != "let" || n.Begin < node.Begin {
					return "PUNCT"
				}
			}
		}
		return "SYM"
	}
	if pt == "adj" {
		if rel == "det" {
			return "DET"
		}
		if rel == "hd" && node.parent.Cat == "pp" {
			// vol vertrouwen
			return "ADP"
		}
		if rel == "crd" {
			// respectievelijk
			return "CCONJ"
		}
		return "ADJ" // exceptions forced by 2.4 validation
	}
	if pt == "bw" {
		return "ADV"
	}
	if pt == "lid" {
		return "DET"
	}
	if pt == "n" {
		if node.Ntype == "eigen" {
			return "PROPN"
		}
		return "NOUN"
	}
	if pt == "spec" {
		if node.Spectype == "deeleigen" {
			return "PROPN"
		}
		if node.Spectype == "symb" {
			return "SYM"
		}
		return "X" // afk vreemd afgebr enof meta
	}
	if pt == "tsw" {
		return "INTJ"
	}
	if pt == "tw" {
		if node.Numtype == "rang" {
			return "ADJ"
		}
		return "NUM"
	}
	if pt == "vz" {
		return "ADP" // v2: do not use PART for SVPs and complementizers
	}
	if pt == "vnw" {
		if rel == "det" && node.Vwtype != "bez" {
			return "DET"
		}
		if node.Pdtype == "adv-pron" {
			if rel == "pobj1" {
				return "PRON"
			}
			return "ADV"
		}
		if (rel == "mod" || (rel == "hd" && node.parent.Rel == "mod")) && node.Pdtype == "grad" {
			// veel minder
			return "ADV"
		}
		return "PRON"
	}
	if pt == "vg" {
		if node.Conjtype == "neven" {
			return "CCONJ" // V2: CONJ ==> CCONJ
		}
		return "SCONJ"
	}
	if pt == "ww" {
		aux, err := auxiliary1(node, q)
		if err != nil {
			panic(fmt.Sprintf("No pos found for %s:%s - %v", number(node.End), node.Word, err))
		}
		if aux == "verb" {
			return "VERB"
		}
		return "AUX" // v2: cop and aux:pass --> AUX  (already in place in v1?)
	}
	panic(fmt.Sprintf("No pos found for %s:%s", number(node.End), node.Word))
}
