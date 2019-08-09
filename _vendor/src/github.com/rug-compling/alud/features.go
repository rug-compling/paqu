package alud

import (
	"fmt"
)

func addFeatures(q *context) {
	for _, node := range q.ptnodes {
		switch node.udPos {
		case "NOUN", "PROPN":
			nominalFeatues(node, q)
		case "ADJ":
			adjectiveFeatures(node, q)
		case "PRON":
			pronounFeatures(node, q)
		case "VERB", "AUX":
			verbalFeatures(node, q)
		case "DET":
			determinerFeatures(node, q)
		case "X":
			specialFeatures(node, q)
		}
	}
}

func nominalFeatues(node *nodeType, q *context) {
	switch node.Genus {
	case "zijd":
		node.udGender = "Com"
	case "onz":
		node.udGender = "Neut"
	case "genus":
		node.udGender = "Com,Neut"
	case "":
		node.udGender = ""
	default:
		panic(fmt.Sprintf("Irregular gender for %s:%s : %s", number(node.End), node.Word, node.Genus))
	}
	switch node.Getal {
	case "ev":
		node.udNumber = "Sing"
	case "mv":
		node.udNumber = "Plur"
	case "":
		node.udNumber = ""
	default:
		panic(fmt.Sprintf("Irregular number for %s:%s : %s", number(node.End), node.Word, node.Getal))
	}
	switch node.Graad {
	case "dim":
		node.udDegree = ""
	case "basis":
		node.udDegree = ""
	case "":
		node.udDegree = ""
	default:
		panic(fmt.Sprintf("Irregular degree for %s:%s : %s", number(node.End), node.Word, node.Graad))
	}
}

func adjectiveFeatures(node *nodeType, q *context) {
	switch node.Graad {
	case "basis":
		node.udDegree = "Pos"
	case "comp":
		node.udDegree = "Cmp"
	case "sup":
		node.udDegree = "Sup"
	case "dim": // netjes
		node.udDegree = "Pos"
	case "":
		node.udDegree = ""
	default:
		panic(fmt.Sprintf("Irregular degree for %s:%s : %s", number(node.End), node.Word, node.Graad))
	}
}

func pronounFeatures(node *nodeType, q *context) {
	switch node.Vwtype {
	case "refl":
		node.udPronType = "Prs"
		node.udReflex = "Yes"
	case "bez":
		node.udPronType = "Prs"
		node.udPoss = "Yes"
	case "pers", "pr":
		node.udPronType = "Prs"
	case "recip":
		node.udPronType = "Rcp"
	case "vb":
		node.udPronType = "Int"
	case "aanw":
		node.udPronType = "Dem"
	case "onbep":
		node.udPronType = "Ind"
	case "betr":
		node.udPronType = "Rel"
	case "excl": // occurs only once
		node.udPronType = ""
	case "":
		node.udPronType = ""
	default:
		panic(fmt.Sprintf("Irregular prontype for %s:%s : %s", number(node.End), node.Word, node.Vwtype))
	}

	switch node.Persoon {
	case "1":
		node.udPerson = "1"
	case "2", "2b", "2v":
		node.udPerson = "2"
	case "3", "3o", "3v", "3p", "3m":
		node.udPerson = "3"
	case "persoon":
		node.udPerson = ""
	case "":
		node.udPerson = ""
	default:
		panic(fmt.Sprintf("Irregular person for %s:%s : %s", number(node.End), node.Word, node.Persoon))
	}

	switch node.Naamval {
	case "nomin":
		node.udCase = "Nom"
	case "obl":
		node.udCase = "Acc"
	case "gen":
		node.udCase = "Gen"
	case "dat": // van dien aard
		node.udCase = "Dat"
	case "stan":
		node.udCase = ""
	case "":
		node.udCase = ""
	default:
		panic(fmt.Sprintf("Irregular case for %s:%s : %s", number(node.End), node.Word, node.Naamval))
	}
}

func verbalFeatures(node *nodeType, q *context) {

	switch node.Wvorm {
	case "pv":
		node.udVerbForm = "Fin"
	case "inf":
		node.udVerbForm = "Inf"
	case "vd", "od":
		node.udVerbForm = "Part"
	case "":
		node.udVerbForm = ""
	default:
		panic(fmt.Sprintf("Irregular verbform for %s:%s : %s", number(node.End), node.Word, node.Wvorm))
	}

	switch node.Pvtijd {
	case "verl":
		node.udTense = "Past"
	case "tgw", "conj":
		node.udTense = "Pres"
	case "":
		node.udTense = ""
	default:
		panic(fmt.Sprintf("Irregular tense for %s:%s : %s", number(node.End), node.Word, node.Pvtijd))
	}

	switch node.Pvagr {
	case "ev", "met-t":
		node.udNumber = "Sing"
	case "mv":
		node.udNumber = "Plur"
	case "":
		node.udNumber = ""
	default:
		panic(fmt.Sprintf("Irregular number for %s:%s : %s", number(node.End), node.Word, node.Pvagr))
	}
}

func determinerFeatures(node *nodeType, q *context) {
	switch node.Lwtype {
	case "bep":
		node.udDefinite = "Def"
	case "onbep":
		node.udDefinite = "Ind"
	case "":
		node.udDefinite = ""
	default:
		panic(fmt.Sprintf("Irregular definite for %s:%s : %s", number(node.End), node.Word, node.Lwtype))
	}
}

func specialFeatures(node *nodeType, q *context) {
	switch node.Spectype {
	case "vreemd":
		node.udForeign = "Yes"
	case "afk":
		node.udAbbr = "Yes"
	}
}
