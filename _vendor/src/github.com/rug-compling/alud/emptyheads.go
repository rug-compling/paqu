//
// // THIS IS A GENERATED FILE. DO NOT EDIT.
//

package alud

import (
	"sort"
)

func reconstructEmptyHead(q *context) bool {
	seen := make(map[int]bool)
	found := false
	for _, n := range q.varindexnodes {
		node := n.(*nodeType)

		if node.Rel != "hd" || node.Pt != "" || node.Cat != "" {
			continue
		}

		antecedent := find(q /* $q.varindexnodes[(@pt or @cat) and @index = $node/@index ] */, &xPath{
			arg1: &dSort{
				arg1: &dFilter{
					arg1: &dVariable{
						VAR: q.varindexnodes,
					},
					arg2: &dSort{
						arg1: &dAnd{
							arg1: &dSort{
								arg1: &dOr{
									arg1: &dCollect{
										ARG:  collect__attributes__pt,
										arg1: &dNode{},
									},
									arg2: &dCollect{
										ARG:  collect__attributes__cat,
										arg1: &dNode{},
									},
								},
							},
							arg2: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG:  collect__attributes__index,
									arg1: &dNode{},
								},
								arg2: &dCollect{
									ARG: collect__attributes__index,
									arg1: &dVariable{
										VAR: node,
									},
								},
							},
						},
					},
				},
			},
		})
		if !test(q, /* $antecedent/@word (: onder andere as hd... :)
			   (: and not(local:auxiliary($antecedent) = ("aux","aux:pass","cop")) skip auxiliaries and copulas, prepositions as well? :)
			*/&xPath{
				arg1: &dSort{
					arg1: &dCollect{
						ARG: collect__attributes__word,
						arg1: &dVariable{
							VAR: antecedent,
						},
					},
				},
			}) {
			continue
		}
		found = true

		others := find(q /* $node/../node[@pt or @cat] */, &xPath{
			arg1: &dSort{
				arg1: &dCollect{
					ARG: collect__child__node,
					arg1: &dCollect{
						ARG: collect__parent__type__node,
						arg1: &dVariable{
							VAR: node,
						},
					},
					arg2: &dPredicate{
						arg1: &dOr{
							arg1: &dCollect{
								ARG:  collect__attributes__pt,
								arg1: &dNode{},
							},
							arg2: &dCollect{
								ARG:  collect__attributes__cat,
								arg1: &dNode{},
							},
						},
					},
				},
			},
		})
		var end int
		if len(others) > 0 {
			if test(q /* $node/../node[@pt or @cat]/@begin = $node/../@begin */, &xPath{
				arg1: &dSort{
					arg1: &dEqual{
						ARG: equal__is,
						arg1: &dCollect{
							ARG: collect__attributes__begin,
							arg1: &dCollect{
								ARG: collect__child__node,
								arg1: &dCollect{
									ARG: collect__parent__type__node,
									arg1: &dVariable{
										VAR: node,
									},
								},
								arg2: &dPredicate{
									arg1: &dOr{
										arg1: &dCollect{
											ARG:  collect__attributes__pt,
											arg1: &dNode{},
										},
										arg2: &dCollect{
											ARG:  collect__attributes__cat,
											arg1: &dNode{},
										},
									},
								},
							},
						},
						arg2: &dCollect{
							ARG: collect__attributes__begin,
							arg1: &dCollect{
								ARG: collect__parent__type__node,
								arg1: &dVariable{
									VAR: node,
								},
							},
						},
					},
				},
			}) {
				end = nZ(others).End + 1 // + 0.1
			} else {
				end = leftEdge(n1(others), q) + 1 // + 0.1
			}
		} else {
			end = i1(find(q /* $node/../@end */, &xPath{
				arg1: &dSort{
					arg1: &dCollect{
						ARG: collect__attributes__end,
						arg1: &dCollect{
							ARG: collect__parent__type__node,
							arg1: &dVariable{
								VAR: node,
							},
						},
					},
				},
			})) - 999 // - 0.9 // covers cases where there is no sister with content
		}
		for seen[end] {
			end++
		}
		seen[end] = true
		var copied int
		antenode := n1(antecedent)
		if antenode.udCopiedFrom > 0 {
			copied = antenode.udCopiedFrom
		} else {
			copied = antenode.End
		}

		node.udOldState = &nodeType{
			Begin:  node.Begin,
			End:    node.End,
			Word:   node.Word,
			Lemma:  node.Lemma,
			Postag: node.Postag,
			Pt:     node.Pt,
		}

		node.Begin = end - 1
		node.End = end
		node.Word = antenode.Word
		node.Lemma = antenode.Lemma
		node.Postag = antenode.Postag
		node.Pt = antenode.Pt
		node.udRelation = "_"
		node.udHeadPosition = underscore
		node.udCopiedFrom = copied

		// kopieer verder alle ud-attributen
		node.udAbbr = antenode.udAbbr
		node.udCase = antenode.udCase
		//niet: node.udCopiedFrom = antenode.udCopiedFrom
		node.udDefinite = antenode.udDefinite
		node.udDegree = antenode.udDegree
		node.udEnhanced = antenode.udEnhanced
		node.udForeign = antenode.udForeign
		node.udGender = antenode.udGender
		//niet: node.udHeadPosition = antenode.udHeadPosition
		node.udNumber = antenode.udNumber
		node.udPerson = antenode.udPerson
		node.udPos = antenode.udPos
		node.udPoss = antenode.udPoss
		node.udPronType = antenode.udPronType
		node.udReflex = antenode.udReflex
		//niet: node.udRelation = antenode.udRelation
		node.udTense = antenode.udTense
		node.udVerbForm = antenode.udVerbForm
		node.udFirstWordBegin = antenode.udFirstWordBegin
		node.udERelation = antenode.udERelation
		node.udEHeadPosition = antenode.udEHeadPosition

		q.ptnodes = append(q.ptnodes, node)
		q.varptnodes = append(q.varptnodes, node)
	}
	if found {
		sort.Slice(q.ptnodes, func(i, j int) bool {
			return q.ptnodes[i].End < q.ptnodes[j].End
		})
		sort.Slice(q.varptnodes, func(i, j int) bool {
			return q.varptnodes[i].(*nodeType).End < q.varptnodes[j].(*nodeType).End
		})
	}
	return found
}

func leftEdge(node *nodeType, q *context) int {
	left := 1000000
	for _, n := range find(q /* $node/descendant-or-self::node[@pt] */, &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__descendant__or__self__node,
				arg1: &dVariable{
					VAR: node,
				},
				arg2: &dPredicate{
					arg1: &dCollect{
						ARG:  collect__attributes__pt,
						arg1: &dNode{},
					},
				},
			},
		},
	}) {
		if begin := n.(*nodeType).Begin; begin < left {
			left = begin
		}
	}
	return left
}
