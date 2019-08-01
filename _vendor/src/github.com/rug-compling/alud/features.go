package alud

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
		node.udGender = "ERROR_IRREGULAR_GENDER"

	}
	switch node.Getal {
	case "ev":
		node.udNumber = "Sing"
	case "mv":
		node.udNumber = "Plur"
	case "":
		node.udNumber = ""
	default:
		node.udNumber = "ERROR_IRREGULAR_NUMBER"
	}
	switch node.Graad {
	case "dim":
		node.udDegree = ""
	case "basis":
		node.udDegree = ""
	case "":
		node.udDegree = ""
	default:
		node.udDegree = "ERROR_IRREGULAR_DEGREE"
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
		node.udDegree = "ERROR_IRREGULAR_DEGREE"
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
		node.udPronType = "ERROR_IRREGULAR_PRONTYPE"
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
		node.udPerson = "ERROR_IRREGULAR_PERSON"
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
		node.udCase = "ERROR_IRREGULAR_CASE"

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
		node.udVerbForm = "ERROR_IRREGULAR_VERBFORM"
	}

	switch node.Pvtijd {
	case "verl":
		node.udTense = "Past"
	case "tgw", "conj":
		node.udTense = "Pres"
	case "":
		node.udTense = ""
	default:
		node.udTense = "ERROR_IRREGULAR_TENSE"
	}

	switch node.Pvagr {
	case "ev", "met-t":
		node.udNumber = "Sing"
	case "mv":
		node.udNumber = "Plur"
	case "":
		node.udNumber = ""
	default:
		node.udNumber = "ERROR_IRREGULAR_NUMBER"
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
		node.udDefinite = "ERROR_IRREGULAR_DEFINITE"
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
