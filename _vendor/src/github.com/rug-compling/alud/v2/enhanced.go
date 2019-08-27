//
// // THIS IS A GENERATED FILE. DO NOT EDIT.
//

package alud

import (
	"fmt"
	"sort"
	"strings"
)

type depT struct {
	head int
	dep  string
}

func enhancedDependencies(q *context) {

	var node *nodeType

	defer func() {
		if r := recover(); r != nil {
			panic(trace(r, "enhancedDependencies", q, node))
		}
	}()

	changed := reconstructEmptyHead(q)

	// add_Edependency_relations
	for _, node = range q.ptnodes {
		// Edependency_relation
		if changed {
			q.depth = 0
			node.udERelation = dependencyLabel(node, q)
			q.depth = 0
			node.udEHeadPosition = externalHeadPosition(list(node), q)
		} else {
			node.udERelation = node.udRelation
			node.udEHeadPosition = node.udHeadPosition
		}
	}

	for _, node = range q.ptnodes {
		enhancedDependencies1(node, q)
	}
}

func enhancedDependencies1(node *nodeType, q *context) {

	defer func() {
		if r := recover(); r != nil {
			panic(trace(r, "enhancedDependencies1", q, node))
		}
	}()

	// iobj2 control : de commissie kan de raad aanbevelen/adviseren/ X te doen
	// rhd : een levend visje dat doorgeslikt moet worden
	q.depth = 0
	var enhanced []depT
	for {

		if test(q /* $node[@ud:ERelation=("nsubj","obj","iobj","nsubj:pass")] */, &xPath{
			arg1: &dSort{
				arg1: &dFilter{
					arg1: &dVariable{
						VAR: node,
					},
					arg2: &dSort{
						arg1: &dEqual{
							ARG: equal__is,
							arg1: &dCollect{
								ARG:  collect__attributes__ud_3aERelation,
								arg1: &dNode{},
							},
							arg2: &dElem{
								DATA: []interface{}{"nsubj", "obj", "iobj", "nsubj:pass"},
								arg1: &dCollect{
									ARG:  collect__attributes__ud_3aERelation,
									arg1: &dNode{},
								},
							},
						},
					},
				},
			},
		}) { // TODO: klopt dit? exists binnen [ ]
			so := find(q,
				/* $node/ancestor::node/node[@rel=("su","obj1","obj2") and local:internal_head_position(.) = $node/@end ]/@index */ &xPath{
					arg1: &dSort{
						arg1: &dCollect{
							ARG: collect__attributes__index,
							arg1: &dCollect{
								ARG: collect__child__node,
								arg1: &dCollect{
									ARG: collect__ancestors__node,
									arg1: &dVariable{
										VAR: node,
									},
								},
								arg2: &dPredicate{
									arg1: &dAnd{
										arg1: &dEqual{
											ARG: equal__is,
											arg1: &dCollect{
												ARG:  collect__attributes__rel,
												arg1: &dNode{},
											},
											arg2: &dElem{
												DATA: []interface{}{"su", "obj1", "obj2"},
												arg1: &dCollect{
													ARG:  collect__attributes__rel,
													arg1: &dNode{},
												},
											},
										},
										arg2: &dEqual{
											ARG: equal__is,
											arg1: &dFunction{
												ARG: function__local__internal__head__position__1__args,
												arg1: &dArg{
													arg1: &dSort{
														arg1: &dNode{},
													},
												},
											},
											arg2: &dCollect{
												ARG: collect__attributes__end,
												arg1: &dVariable{
													VAR: node,
												},
											},
										},
									},
								},
							},
						},
					},
				})
			if len(so) > 0 {
				soIndex := i1(so)
				//NP
				enhanced = []depT{depT{head: node.udEHeadPosition, dep: enhanceDependencyLabel(node, q)}} // self
				enhanced = append(enhanced, anaphoricRelpronoun(node, q)...)                              // self
				//NP
				enhanced = append(enhanced, distributeConjuncts(node, q)...) // self
				//NP
				enhanced = append(enhanced, distributeDependents(node, q)...) // self
				//NP
				enhanced = append(enhanced, xcompControl(node, q, soIndex)...)
				//NP
				enhanced = append(enhanced, upstairsControl(node, q, soIndex)...)
				//NP
				enhanced = append(enhanced, passiveVpControl(node, q, soIndex)...)
				break
			}
		}

		rhd := find(q,
			/* $node/ancestor::node/node[@rel="rhd" and local:internal_head_position(.) = $node/@end ]/@index */ &xPath{
				arg1: &dSort{
					arg1: &dCollect{
						ARG: collect__attributes__index,
						arg1: &dCollect{
							ARG: collect__child__node,
							arg1: &dCollect{
								ARG: collect__ancestors__node,
								arg1: &dVariable{
									VAR: node,
								},
							},
							arg2: &dPredicate{
								arg1: &dAnd{
									arg1: &dEqual{
										ARG: equal__is,
										arg1: &dCollect{
											ARG:  collect__attributes__rel,
											arg1: &dNode{},
										},
										arg2: &dElem{
											DATA: []interface{}{"rhd"},
											arg1: &dCollect{
												ARG:  collect__attributes__rel,
												arg1: &dNode{},
											},
										},
									},
									arg2: &dEqual{
										ARG: equal__is,
										arg1: &dFunction{
											ARG: function__local__internal__head__position__1__args,
											arg1: &dArg{
												arg1: &dSort{
													arg1: &dNode{},
												},
											},
										},
										arg2: &dCollect{
											ARG: collect__attributes__end,
											arg1: &dVariable{
												VAR: node,
											},
										},
									},
								},
							},
						},
					},
				},
			})
		if len(rhd) > 0 {
			rhdIndex := i1(rhd)
			rhdNp := find(q /* $node/ancestor::node[@cat="np" and node[@rel="mod"]/node[@rel="rhd"]/@index = $rhdIndex] */, &xPath{
				arg1: &dSort{
					arg1: &dCollect{
						ARG: collect__ancestors__node,
						arg1: &dVariable{
							VAR: node,
						},
						arg2: &dPredicate{
							arg1: &dAnd{
								arg1: &dEqual{
									ARG: equal__is,
									arg1: &dCollect{
										ARG:  collect__attributes__cat,
										arg1: &dNode{},
									},
									arg2: &dElem{
										DATA: []interface{}{"np"},
										arg1: &dCollect{
											ARG:  collect__attributes__cat,
											arg1: &dNode{},
										},
									},
								},
								arg2: &dEqual{
									ARG: equal__is,
									arg1: &dCollect{
										ARG: collect__attributes__index,
										arg1: &dCollect{
											ARG: collect__child__node,
											arg1: &dCollect{
												ARG:  collect__child__node,
												arg1: &dNode{},
												arg2: &dPredicate{
													arg1: &dEqual{
														ARG: equal__is,
														arg1: &dCollect{
															ARG:  collect__attributes__rel,
															arg1: &dNode{},
														},
														arg2: &dElem{
															DATA: []interface{}{"mod"},
															arg1: &dCollect{
																ARG:  collect__attributes__rel,
																arg1: &dNode{},
															},
														},
													},
												},
											},
											arg2: &dPredicate{
												arg1: &dEqual{
													ARG: equal__is,
													arg1: &dCollect{
														ARG:  collect__attributes__rel,
														arg1: &dNode{},
													},
													arg2: &dElem{
														DATA: []interface{}{"rhd"},
														arg1: &dCollect{
															ARG:  collect__attributes__rel,
															arg1: &dNode{},
														},
													},
												},
											},
										},
									},
									arg2: &dVariable{
										VAR: rhdIndex,
									},
								},
							},
						},
					},
				},
			})
			// de enige _i die voldoet aan de eisen -- make sure empty heads are covered as well
			if len(rhdNp) > 0 {
				//NP
				enhanced = []depT{depT{head: internalHeadPositionWithGapping(rhdNp, q), dep: "ref"}} // rhdref
				//NP
				enhanced = append(enhanced, xcompControl(node, q, rhdIndex)...)
				//NP
				enhanced = append(enhanced, passiveVpControl(node, q, rhdIndex)...)
				break
			}
			// if there is no antecedent, lets keep the basic relation
			//NP
			enhanced = []depT{depT{head: node.udEHeadPosition, dep: enhanceDependencyLabel(node, q)}} // self
			enhanced = append(enhanced, anaphoricRelpronoun(node, q)...)                              // self
			//NP
			enhanced = append(enhanced, distributeConjuncts(node, q)...) // self
			//NP
			enhanced = append(enhanced, distributeDependents(node, q)...) // self
			//NP
			enhanced = append(enhanced, xcompControl(node, q, rhdIndex)...)
			//NP
			enhanced = append(enhanced, passiveVpControl(node, q, rhdIndex)...)
			break
		}

		relSister := find(q /* ($node/../node[@rel="mod" and @cat="rel"]/node[@rel="rhd"]/@index)[1] */, &xPath{
			arg1: &dSort{
				arg1: &dFilter{
					arg1: &dSort{
						arg1: &dCollect{
							ARG: collect__attributes__index,
							arg1: &dCollect{
								ARG: collect__child__node,
								arg1: &dCollect{
									ARG: collect__child__node,
									arg1: &dCollect{
										ARG: collect__parent__type__node,
										arg1: &dVariable{
											VAR: node,
										},
									},
									arg2: &dPredicate{
										arg1: &dAnd{
											arg1: &dEqual{
												ARG: equal__is,
												arg1: &dCollect{
													ARG:  collect__attributes__rel,
													arg1: &dNode{},
												},
												arg2: &dElem{
													DATA: []interface{}{"mod"},
													arg1: &dCollect{
														ARG:  collect__attributes__rel,
														arg1: &dNode{},
													},
												},
											},
											arg2: &dEqual{
												ARG: equal__is,
												arg1: &dCollect{
													ARG:  collect__attributes__cat,
													arg1: &dNode{},
												},
												arg2: &dElem{
													DATA: []interface{}{"rel"},
													arg1: &dCollect{
														ARG:  collect__attributes__cat,
														arg1: &dNode{},
													},
												},
											},
										},
									},
								},
								arg2: &dPredicate{
									arg1: &dEqual{
										ARG: equal__is,
										arg1: &dCollect{
											ARG:  collect__attributes__rel,
											arg1: &dNode{},
										},
										arg2: &dElem{
											DATA: []interface{}{"rhd"},
											arg1: &dCollect{
												ARG:  collect__attributes__rel,
												arg1: &dNode{},
											},
										},
									},
								},
							},
						},
					},
					arg2: &dSort{
						arg1: &dFunction{
							ARG: function__first__0__args,
						},
					},
				},
			},
		})
		if len(relSister) > 0 {
			relSisterIndex := i1(relSister)
			//NP
			enhanced = []depT{depT{head: node.udEHeadPosition, dep: enhanceDependencyLabel(node, q)}} // self
			enhanced = append(enhanced, anaphoricRelpronoun(node, q)...)                              // self
			//NP
			enhanced = append(enhanced, distributeConjuncts(node, q)...) // self
			//NP
			enhanced = append(enhanced, distributeDependents(node, q)...) // self
			//NP
			enhanced = append(enhanced, xcompControl(node, q, relSisterIndex)...)
			//NP
			enhanced = append(enhanced, passiveVpControl(node, q, relSisterIndex)...)
			break
		}

		// underscore is resultaat van reconstructEmptyHead()
		if node.udHeadPosition >= 0 || node.udHeadPosition == underscore {
			//NP
			enhanced = []depT{depT{head: node.udEHeadPosition, dep: enhanceDependencyLabel(node, q)}} // self
			enhanced = append(enhanced, anaphoricRelpronoun(node, q)...)                              // self
			//NP
			enhanced = append(enhanced, distributeConjuncts(node, q)...) // self
			//NP
			enhanced = append(enhanced, distributeDependents(node, q)...) // self
			break
		}

		//NP
		enhanced = []depT{depT{head: node.udEHeadPosition, dep: enhanceDependencyLabel(node, q)}}
		break
	}

	sort.Slice(enhanced, func(i, j int) bool {
		if enhanced[i].head != enhanced[j].head {
			return enhanced[i].head < enhanced[j].head
		}
		return enhanced[i].dep < enhanced[j].dep
	})
	for i := 1; i < len(enhanced); i++ {
		if enhanced[i].head == enhanced[i-1].head && enhanced[i].dep == enhanced[i-1].dep {
			enhanced = append(enhanced[:i], enhanced[1+i:]...)
			i--
		}
	}
	ss := make([]string, 0, len(enhanced))
	for _, e := range enhanced {
		if e.head == 0 && e.dep != "root" ||
			e.head != 0 && e.dep == "root" ||
			e.dep == "orphan" {
			panic(fmt.Sprintf("Invalid EUD %s:%s", number(e.head), e.dep))
		}
		if e.dep != "" {
			ss = append(ss, number(e.head)+":"+e.dep)
		}
	}
	node.udEnhanced = strings.Join(ss, "|")

}

func join(a, b string) string {
	if b == "" {
		return a
	}
	return a + ":" + b
}

func enhanceDependencyLabel(node *nodeType, q *context) string {

	defer func() {
		if r := recover(); r != nil {
			panic(trace(r, "enhanceDependencyLabel", q, node))
		}
	}()

	label := node.udERelation
	if label == "conj" {
		if crd := n1(find(q, /* ($node/ancestor::node[@cat="conj" and
			   not(.//node[@cat="conj"]//node/@begin = $node/@begin)]/node[@rel="crd"])[1] */&xPath{
				arg1: &dSort{
					arg1: &dFilter{
						arg1: &dSort{
							arg1: &dCollect{
								ARG: collect__child__node,
								arg1: &dCollect{
									ARG: collect__ancestors__node,
									arg1: &dVariable{
										VAR: node,
									},
									arg2: &dPredicate{
										arg1: &dAnd{
											arg1: &dEqual{
												ARG: equal__is,
												arg1: &dCollect{
													ARG:  collect__attributes__cat,
													arg1: &dNode{},
												},
												arg2: &dElem{
													DATA: []interface{}{"conj"},
													arg1: &dCollect{
														ARG:  collect__attributes__cat,
														arg1: &dNode{},
													},
												},
											},
											arg2: &dFunction{
												ARG: function__not__1__args,
												arg1: &dArg{
													arg1: &dSort{
														arg1: &dEqual{
															ARG: equal__is,
															arg1: &dCollect{
																ARG: collect__attributes__begin,
																arg1: &dCollect{
																	ARG: collect__descendant__node,
																	arg1: &dCollect{
																		ARG: collect__child__node,
																		arg1: &dCollect{
																			ARG:  collect__descendant__or__self__type__node,
																			arg1: &dNode{},
																		},
																		arg2: &dPredicate{
																			arg1: &dEqual{
																				ARG: equal__is,
																				arg1: &dCollect{
																					ARG:  collect__attributes__cat,
																					arg1: &dNode{},
																				},
																				arg2: &dElem{
																					DATA: []interface{}{"conj"},
																					arg1: &dCollect{
																						ARG:  collect__attributes__cat,
																						arg1: &dNode{},
																					},
																				},
																			},
																		},
																	},
																},
															},
															arg2: &dCollect{
																ARG: collect__attributes__begin,
																arg1: &dVariable{
																	VAR: node,
																},
															},
														},
													},
												},
											},
										},
									},
								},
								arg2: &dPredicate{
									arg1: &dEqual{
										ARG: equal__is,
										arg1: &dCollect{
											ARG:  collect__attributes__rel,
											arg1: &dNode{},
										},
										arg2: &dElem{
											DATA: []interface{}{"crd"},
											arg1: &dCollect{
												ARG:  collect__attributes__rel,
												arg1: &dNode{},
											},
										},
									},
								},
							},
						},
						arg2: &dSort{
							arg1: &dFunction{
								ARG: function__first__0__args,
							},
						},
					},
				},
			})); crd != noNode {
			if crd.Lemma != "" {
				return join(label, enhancedLemmaString1(crd, q))
			}
			if crd.Cat == "mwu" {
				return join(label, enhancedLemmaString1(n1(find(q /* ($crd/node[@rel="mwp"])[1] */, &xPath{
					arg1: &dSort{
						arg1: &dFilter{
							arg1: &dSort{
								arg1: &dCollect{
									ARG: collect__child__node,
									arg1: &dVariable{
										VAR: crd,
									},
									arg2: &dPredicate{
										arg1: &dEqual{
											ARG: equal__is,
											arg1: &dCollect{
												ARG:  collect__attributes__rel,
												arg1: &dNode{},
											},
											arg2: &dElem{
												DATA: []interface{}{"mwp"},
												arg1: &dCollect{
													ARG:  collect__attributes__rel,
													arg1: &dNode{},
												},
											},
										},
									},
								},
							},
							arg2: &dSort{
								arg1: &dFunction{
									ARG: function__first__0__args,
								},
							},
						},
					},
				})), q))
			}
			panic("Empty EUD label")
		}
	}

	if label == "nmod" || label == "obl" {
		if casee := n1(find(q /* ($q.varptnodes[@ud:ERelation="case" and @ud:EHeadPosition=$node/@end])[1] */, &xPath{
			arg1: &dSort{
				arg1: &dFilter{
					arg1: &dSort{
						arg1: &dFilter{
							arg1: &dVariable{
								VAR: q.varptnodes,
							},
							arg2: &dSort{
								arg1: &dAnd{
									arg1: &dEqual{
										ARG: equal__is,
										arg1: &dCollect{
											ARG:  collect__attributes__ud_3aERelation,
											arg1: &dNode{},
										},
										arg2: &dElem{
											DATA: []interface{}{"case"},
											arg1: &dCollect{
												ARG:  collect__attributes__ud_3aERelation,
												arg1: &dNode{},
											},
										},
									},
									arg2: &dEqual{
										ARG: equal__is,
										arg1: &dCollect{
											ARG:  collect__attributes__ud_3aEHeadPosition,
											arg1: &dNode{},
										},
										arg2: &dCollect{
											ARG: collect__attributes__end,
											arg1: &dVariable{
												VAR: node,
											},
										},
									},
								},
							},
						},
					},
					arg2: &dSort{
						arg1: &dFunction{
							ARG: function__first__0__args,
						},
					},
				},
			},
		})); casee != noNode {
			return join(label, enhancedLemmaString1(casee, q))
		}
	}

	if label == "advcl" || label == "acl" {
		if mark := n1(find(q /* ($q.varptnodes[@ud:ERelation=("mark","case") and @ud:EHeadPosition=$node/@end])[1] */, &xPath{
			arg1: &dSort{
				arg1: &dFilter{
					arg1: &dSort{
						arg1: &dFilter{
							arg1: &dVariable{
								VAR: q.varptnodes,
							},
							arg2: &dSort{
								arg1: &dAnd{
									arg1: &dEqual{
										ARG: equal__is,
										arg1: &dCollect{
											ARG:  collect__attributes__ud_3aERelation,
											arg1: &dNode{},
										},
										arg2: &dElem{
											DATA: []interface{}{"mark", "case"},
											arg1: &dCollect{
												ARG:  collect__attributes__ud_3aERelation,
												arg1: &dNode{},
											},
										},
									},
									arg2: &dEqual{
										ARG: equal__is,
										arg1: &dCollect{
											ARG:  collect__attributes__ud_3aEHeadPosition,
											arg1: &dNode{},
										},
										arg2: &dCollect{
											ARG: collect__attributes__end,
											arg1: &dVariable{
												VAR: node,
											},
										},
									},
								},
							},
						},
					},
					arg2: &dSort{
						arg1: &dFunction{
							ARG: function__first__0__args,
						},
					},
				},
			},
		})); mark != noNode {
			return join(label, enhancedLemmaString1(mark, q))
		}
	}

	if label != "" {
		return label
	}

	panic("Empty EUD label")
}

func anaphoricRelpronoun(node *nodeType, q *context) []depT {
	// works voor waar, and last() picks waar in 'daar waar' cases
	// dont add anything for hij werd voorzitter, wat hij nog steeds is (otherwise self-reference)
	// for loop ensures correct result if N has 2 acl:relcl dependents
	result := []depT{}
	for _, a := range find(q, /* $node/ancestor::node[@cat="np" and local:internal_head_position(.) = $node/@end]/
		   node[@rel="mod"]/node[@rel="rhd"]/descendant-or-self::node[@pt="vnw" and not(@ud:HeadPosition = $node/@end)][last()] */&xPath{
			arg1: &dSort{
				arg1: &dCollect{
					ARG: collect__descendant__or__self__node,
					arg1: &dCollect{
						ARG: collect__child__node,
						arg1: &dCollect{
							ARG: collect__child__node,
							arg1: &dCollect{
								ARG: collect__ancestors__node,
								arg1: &dVariable{
									VAR: node,
								},
								arg2: &dPredicate{
									arg1: &dAnd{
										arg1: &dEqual{
											ARG: equal__is,
											arg1: &dCollect{
												ARG:  collect__attributes__cat,
												arg1: &dNode{},
											},
											arg2: &dElem{
												DATA: []interface{}{"np"},
												arg1: &dCollect{
													ARG:  collect__attributes__cat,
													arg1: &dNode{},
												},
											},
										},
										arg2: &dEqual{
											ARG: equal__is,
											arg1: &dFunction{
												ARG: function__local__internal__head__position__1__args,
												arg1: &dArg{
													arg1: &dSort{
														arg1: &dNode{},
													},
												},
											},
											arg2: &dCollect{
												ARG: collect__attributes__end,
												arg1: &dVariable{
													VAR: node,
												},
											},
										},
									},
								},
							},
							arg2: &dPredicate{
								arg1: &dEqual{
									ARG: equal__is,
									arg1: &dCollect{
										ARG:  collect__attributes__rel,
										arg1: &dNode{},
									},
									arg2: &dElem{
										DATA: []interface{}{"mod"},
										arg1: &dCollect{
											ARG:  collect__attributes__rel,
											arg1: &dNode{},
										},
									},
								},
							},
						},
						arg2: &dPredicate{
							arg1: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG:  collect__attributes__rel,
									arg1: &dNode{},
								},
								arg2: &dElem{
									DATA: []interface{}{"rhd"},
									arg1: &dCollect{
										ARG:  collect__attributes__rel,
										arg1: &dNode{},
									},
								},
							},
						},
					},
					arg2: &dPredicate{
						arg1: &dPredicate{
							arg1: &dAnd{
								arg1: &dEqual{
									ARG: equal__is,
									arg1: &dCollect{
										ARG:  collect__attributes__pt,
										arg1: &dNode{},
									},
									arg2: &dElem{
										DATA: []interface{}{"vnw"},
										arg1: &dCollect{
											ARG:  collect__attributes__pt,
											arg1: &dNode{},
										},
									},
								},
								arg2: &dFunction{
									ARG: function__not__1__args,
									arg1: &dArg{
										arg1: &dSort{
											arg1: &dEqual{
												ARG: equal__is,
												arg1: &dCollect{
													ARG:  collect__attributes__ud_3aHeadPosition,
													arg1: &dNode{},
												},
												arg2: &dCollect{
													ARG: collect__attributes__end,
													arg1: &dVariable{
														VAR: node,
													},
												},
											},
										},
									},
								},
							},
						},
						arg2: &dFunction{
							ARG: function__last__0__args,
						},
					},
				},
			},
		}) {
		anrel := a.(*nodeType)
		var label string
		if r := anrel.udRelation; r == "nsubj" || r == "nsubj:pass" {
			label = r + ":relsubj"
		} else if r == "obj" || r == "obl" {
			label = r + ":relobj"
		} else {
			label = r
		}
		result = append(result, depT{head: anrel.udHeadPosition, dep: label})
	}
	return result
}

// Glastra en Terlouw verzonnen een list --> nsubj(verzonnen,Glastra) nsubj(verzonnen,Terlouw)
func distributeConjuncts(node *nodeType, q *context) []depT {

	defer func() {
		if r := recover(); r != nil {
			panic(trace(r, "distributeConjuncts", q, node))
		}
	}()

	if node.udRelation == "conj" {
		coordHead := n1(find(q, /* $q.varallnodes[@end = $node/@ud:HeadPosition
			   and @ud:Relation=("amod","appos","nmod","nsubj","nsubj:pass","nummod","obj","iobj","obl","obl:agent","advcl")] */&xPath{
				arg1: &dSort{
					arg1: &dFilter{
						arg1: &dVariable{
							VAR: q.varallnodes,
						},
						arg2: &dSort{
							arg1: &dAnd{
								arg1: &dEqual{
									ARG: equal__is,
									arg1: &dCollect{
										ARG:  collect__attributes__end,
										arg1: &dNode{},
									},
									arg2: &dCollect{
										ARG: collect__attributes__ud_3aHeadPosition,
										arg1: &dVariable{
											VAR: node,
										},
									},
								},
								arg2: &dEqual{
									ARG: equal__is,
									arg1: &dCollect{
										ARG:  collect__attributes__ud_3aRelation,
										arg1: &dNode{},
									},
									arg2: &dElem{
										DATA: []interface{}{"amod", "appos", "nmod", "nsubj", "nsubj:pass", "nummod", "obj", "iobj", "obl", "obl:agent", "advcl"},
										arg1: &dCollect{
											ARG:  collect__attributes__ud_3aRelation,
											arg1: &dNode{},
										},
									},
								},
							},
						},
					},
				},
			}))
		if coordHead != noNode {
			// in A en B vs in A en naast B --> use enh_dep_label($node) in the latter case...
			depLabel := enhanceDependencyLabel(coordHead, q)
			return []depT{depT{head: coordHead.udHeadPosition, dep: depLabel}}
		}
	}
	return []depT{}
}

// de onrust kan een reis vertragen of frustreren  --> obj(vertragen,reis) obj(frustreren,reis)
// todo: passives ze werd ontmanteld en verkocht  su coindexed with two obj1
// done: phrases [np_i [een scoutskameraad] werd .. en _i zocht hem op]
// idem: de hond was gebaseerd op Lassy en verscheen onder de naam Wirel nsubj:pass in conj1, nsubj in conj 2
func distributeDependents(node *nodeType, q *context) []depT {

	defer func() {
		if r := recover(); r != nil {
			panic(trace(r, "distributeDependents", q, node))
		}
	}()

	var phrase *nodeType
	if node.Rel == "hd" {
		if test(q /* $node/../../@cat="pp" */, &xPath{
			arg1: &dSort{
				arg1: &dEqual{
					ARG: equal__is,
					arg1: &dCollect{
						ARG: collect__attributes__cat,
						arg1: &dCollect{
							ARG: collect__parent__type__node,
							arg1: &dCollect{
								ARG: collect__parent__type__node,
								arg1: &dVariable{
									VAR: node,
								},
							},
						},
					},
					arg2: &dElem{
						DATA: []interface{}{"pp"},
						arg1: &dCollect{
							ARG: collect__attributes__cat,
							arg1: &dCollect{
								ARG: collect__parent__type__node,
								arg1: &dCollect{
									ARG: collect__parent__type__node,
									arg1: &dVariable{
										VAR: node,
									},
								},
							},
						},
					},
				},
			},
		}) { // door het schilderij
			phrase = node.parent.parent
		} else {
			phrase = node.parent
		}
	} else {
		if test(q /* $node[@rel="mwp" and @begin = ../@begin] */, &xPath{
			arg1: &dSort{
				arg1: &dFilter{
					arg1: &dVariable{
						VAR: node,
					},
					arg2: &dSort{
						arg1: &dAnd{
							arg1: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG:  collect__attributes__rel,
									arg1: &dNode{},
								},
								arg2: &dElem{
									DATA: []interface{}{"mwp"},
									arg1: &dCollect{
										ARG:  collect__attributes__rel,
										arg1: &dNode{},
									},
								},
							},
							arg2: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG:  collect__attributes__begin,
									arg1: &dNode{},
								},
								arg2: &dCollect{
									ARG: collect__attributes__begin,
									arg1: &dCollect{
										ARG:  collect__parent__type__node,
										arg1: &dNode{},
									},
								},
							},
						},
					},
				},
			},
		}) {
			if test(q /* $node[ ../@rel="obj1" and ../../@cat="pp"] */, &xPath{
				arg1: &dSort{
					arg1: &dFilter{
						arg1: &dVariable{
							VAR: node,
						},
						arg2: &dSort{
							arg1: &dAnd{
								arg1: &dEqual{
									ARG: equal__is,
									arg1: &dCollect{
										ARG: collect__attributes__rel,
										arg1: &dCollect{
											ARG:  collect__parent__type__node,
											arg1: &dNode{},
										},
									},
									arg2: &dElem{
										DATA: []interface{}{"obj1"},
										arg1: &dCollect{
											ARG: collect__attributes__rel,
											arg1: &dCollect{
												ARG:  collect__parent__type__node,
												arg1: &dNode{},
											},
										},
									},
								},
								arg2: &dEqual{
									ARG: equal__is,
									arg1: &dCollect{
										ARG: collect__attributes__cat,
										arg1: &dCollect{
											ARG: collect__parent__type__node,
											arg1: &dCollect{
												ARG:  collect__parent__type__node,
												arg1: &dNode{},
											},
										},
									},
									arg2: &dElem{
										DATA: []interface{}{"pp"},
										arg1: &dCollect{
											ARG: collect__attributes__cat,
											arg1: &dCollect{
												ARG: collect__parent__type__node,
												arg1: &dCollect{
													ARG:  collect__parent__type__node,
													arg1: &dNode{},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			}) {
				phrase = node.parent.parent
			} else {
				if test(q /* $node[../@rel="hd" and ../../@rel="obj1" and ../../../@cat="pp"] */, &xPath{
					arg1: &dSort{
						arg1: &dFilter{
							arg1: &dVariable{
								VAR: node,
							},
							arg2: &dSort{
								arg1: &dAnd{
									arg1: &dAnd{
										arg1: &dEqual{
											ARG: equal__is,
											arg1: &dCollect{
												ARG: collect__attributes__rel,
												arg1: &dCollect{
													ARG:  collect__parent__type__node,
													arg1: &dNode{},
												},
											},
											arg2: &dElem{
												DATA: []interface{}{"hd"},
												arg1: &dCollect{
													ARG: collect__attributes__rel,
													arg1: &dCollect{
														ARG:  collect__parent__type__node,
														arg1: &dNode{},
													},
												},
											},
										},
										arg2: &dEqual{
											ARG: equal__is,
											arg1: &dCollect{
												ARG: collect__attributes__rel,
												arg1: &dCollect{
													ARG: collect__parent__type__node,
													arg1: &dCollect{
														ARG:  collect__parent__type__node,
														arg1: &dNode{},
													},
												},
											},
											arg2: &dElem{
												DATA: []interface{}{"obj1"},
												arg1: &dCollect{
													ARG: collect__attributes__rel,
													arg1: &dCollect{
														ARG: collect__parent__type__node,
														arg1: &dCollect{
															ARG:  collect__parent__type__node,
															arg1: &dNode{},
														},
													},
												},
											},
										},
									},
									arg2: &dEqual{
										ARG: equal__is,
										arg1: &dCollect{
											ARG: collect__attributes__cat,
											arg1: &dCollect{
												ARG: collect__parent__type__node,
												arg1: &dCollect{
													ARG: collect__parent__type__node,
													arg1: &dCollect{
														ARG:  collect__parent__type__node,
														arg1: &dNode{},
													},
												},
											},
										},
										arg2: &dElem{
											DATA: []interface{}{"pp"},
											arg1: &dCollect{
												ARG: collect__attributes__cat,
												arg1: &dCollect{
													ARG: collect__parent__type__node,
													arg1: &dCollect{
														ARG: collect__parent__type__node,
														arg1: &dCollect{
															ARG:  collect__parent__type__node,
															arg1: &dNode{},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				}) {
					phrase = node.parent.parent.parent // in en rond het Hoofstedelijk Gewest --> do not distribute Hoofdstedelijk
				} else {
					if test(q /* $node[ ../@rel="hd" and not( ../../@cat="pp") ] */, &xPath{
						arg1: &dSort{
							arg1: &dFilter{
								arg1: &dVariable{
									VAR: node,
								},
								arg2: &dSort{
									arg1: &dAnd{
										arg1: &dEqual{
											ARG: equal__is,
											arg1: &dCollect{
												ARG: collect__attributes__rel,
												arg1: &dCollect{
													ARG:  collect__parent__type__node,
													arg1: &dNode{},
												},
											},
											arg2: &dElem{
												DATA: []interface{}{"hd"},
												arg1: &dCollect{
													ARG: collect__attributes__rel,
													arg1: &dCollect{
														ARG:  collect__parent__type__node,
														arg1: &dNode{},
													},
												},
											},
										},
										arg2: &dFunction{
											ARG: function__not__1__args,
											arg1: &dArg{
												arg1: &dSort{
													arg1: &dEqual{
														ARG: equal__is,
														arg1: &dCollect{
															ARG: collect__attributes__cat,
															arg1: &dCollect{
																ARG: collect__parent__type__node,
																arg1: &dCollect{
																	ARG:  collect__parent__type__node,
																	arg1: &dNode{},
																},
															},
														},
														arg2: &dElem{
															DATA: []interface{}{"pp"},
															arg1: &dCollect{
																ARG: collect__attributes__cat,
																arg1: &dCollect{
																	ARG: collect__parent__type__node,
																	arg1: &dCollect{
																		ARG:  collect__parent__type__node,
																		arg1: &dNode{},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					}) { // mwu as head, but not complex P
						phrase = node.parent.parent
					} else {
						phrase = node.parent
					}
				}
			}
		} else {
			if test(q /* $node[@rel="obj1" and ../@cat="pp"] */, &xPath{
				arg1: &dSort{
					arg1: &dFilter{
						arg1: &dVariable{
							VAR: node,
						},
						arg2: &dSort{
							arg1: &dAnd{
								arg1: &dEqual{
									ARG: equal__is,
									arg1: &dCollect{
										ARG:  collect__attributes__rel,
										arg1: &dNode{},
									},
									arg2: &dElem{
										DATA: []interface{}{"obj1"},
										arg1: &dCollect{
											ARG:  collect__attributes__rel,
											arg1: &dNode{},
										},
									},
								},
								arg2: &dEqual{
									ARG: equal__is,
									arg1: &dCollect{
										ARG: collect__attributes__cat,
										arg1: &dCollect{
											ARG:  collect__parent__type__node,
											arg1: &dNode{},
										},
									},
									arg2: &dElem{
										DATA: []interface{}{"pp"},
										arg1: &dCollect{
											ARG: collect__attributes__cat,
											arg1: &dCollect{
												ARG:  collect__parent__type__node,
												arg1: &dNode{},
											},
										},
									},
								},
							},
						},
					},
				},
			}) {
				phrase = node.parent
			} else {
				phrase = node
				// do not apply to prepositions and auxiliaries, ever. Too strict?
			}
		}
	}

	if !test(q /* $phrase[@rel=("obj1","su","mod","pc","det") and @index] */, &xPath{
		arg1: &dSort{
			arg1: &dFilter{
				arg1: &dVariable{
					VAR: phrase,
				},
				arg2: &dSort{
					arg1: &dAnd{
						arg1: &dEqual{
							ARG: equal__is,
							arg1: &dCollect{
								ARG:  collect__attributes__rel,
								arg1: &dNode{},
							},
							arg2: &dElem{
								DATA: []interface{}{"obj1", "su", "mod", "pc", "det"},
								arg1: &dCollect{
									ARG:  collect__attributes__rel,
									arg1: &dNode{},
								},
							},
						},
						arg2: &dCollect{
							ARG:  collect__attributes__index,
							arg1: &dNode{},
						},
					},
				},
			},
		},
	}) {
		return []depT{}
	}

	// TODO: dit xpath kan efficiënter?
	conj_heads := find(q, /* $node[not(@ud:pos=("ADP","AUX"))]/ancestor::node//node[@rel="cnj"
			     and node[
			    (: @rel=$phrase/@rel
				and -- this constraint is too strict for coord of passives:)
				not(@pt or @cat)]/@index = $phrase/@index
			     and node[@rel=("hd","predc") and not(@ud:pos="AUX") and (@pt or @cat) and	 (: bekende cafes zijn A en B :)
				(: not(@ud:pos=("ADP","AUX")) and not(@cat="mwu") :)
				not(local:internal_head_position(..) = @end and (@ud:pos=("ADP","AUX") or @cat="mwu") )
				]
		      ]
			      (: not coordination of AUX or (complex) Ps :) */&xPath{
			arg1: &dSort{
				arg1: &dCollect{
					ARG: collect__child__node,
					arg1: &dCollect{
						ARG: collect__descendant__or__self__type__node,
						arg1: &dCollect{
							ARG: collect__ancestors__node,
							arg1: &dFilter{
								arg1: &dVariable{
									VAR: node,
								},
								arg2: &dSort{
									arg1: &dFunction{
										ARG: function__not__1__args,
										arg1: &dArg{
											arg1: &dSort{
												arg1: &dEqual{
													ARG: equal__is,
													arg1: &dCollect{
														ARG:  collect__attributes__ud_3apos,
														arg1: &dNode{},
													},
													arg2: &dElem{
														DATA: []interface{}{"ADP", "AUX"},
														arg1: &dCollect{
															ARG:  collect__attributes__ud_3apos,
															arg1: &dNode{},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					arg2: &dPredicate{
						arg1: &dAnd{
							arg1: &dAnd{
								arg1: &dEqual{
									ARG: equal__is,
									arg1: &dCollect{
										ARG:  collect__attributes__rel,
										arg1: &dNode{},
									},
									arg2: &dElem{
										DATA: []interface{}{"cnj"},
										arg1: &dCollect{
											ARG:  collect__attributes__rel,
											arg1: &dNode{},
										},
									},
								},
								arg2: &dEqual{
									ARG: equal__is,
									arg1: &dCollect{
										ARG: collect__attributes__index,
										arg1: &dCollect{
											ARG:  collect__child__node,
											arg1: &dNode{},
											arg2: &dPredicate{
												arg1: &dFunction{
													ARG: function__not__1__args,
													arg1: &dArg{
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
													},
												},
											},
										},
									},
									arg2: &dCollect{
										ARG: collect__attributes__index,
										arg1: &dVariable{
											VAR: phrase,
										},
									},
								},
							},
							arg2: &dCollect{
								ARG:  collect__child__node,
								arg1: &dNode{},
								arg2: &dPredicate{
									arg1: &dAnd{
										arg1: &dAnd{
											arg1: &dAnd{
												arg1: &dEqual{
													ARG: equal__is,
													arg1: &dCollect{
														ARG:  collect__attributes__rel,
														arg1: &dNode{},
													},
													arg2: &dElem{
														DATA: []interface{}{"hd", "predc"},
														arg1: &dCollect{
															ARG:  collect__attributes__rel,
															arg1: &dNode{},
														},
													},
												},
												arg2: &dFunction{
													ARG: function__not__1__args,
													arg1: &dArg{
														arg1: &dSort{
															arg1: &dEqual{
																ARG: equal__is,
																arg1: &dCollect{
																	ARG:  collect__attributes__ud_3apos,
																	arg1: &dNode{},
																},
																arg2: &dElem{
																	DATA: []interface{}{"AUX"},
																	arg1: &dCollect{
																		ARG:  collect__attributes__ud_3apos,
																		arg1: &dNode{},
																	},
																},
															},
														},
													},
												},
											},
											arg2: &dSort{
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
										arg2: &dFunction{
											ARG: function__not__1__args,
											arg1: &dArg{
												arg1: &dSort{
													arg1: &dAnd{
														arg1: &dEqual{
															ARG: equal__is,
															arg1: &dFunction{
																ARG: function__local__internal__head__position__1__args,
																arg1: &dArg{
																	arg1: &dSort{
																		arg1: &dCollect{
																			ARG:  collect__parent__type__node,
																			arg1: &dNode{},
																		},
																	},
																},
															},
															arg2: &dCollect{
																ARG:  collect__attributes__end,
																arg1: &dNode{},
															},
														},
														arg2: &dSort{
															arg1: &dOr{
																arg1: &dEqual{
																	ARG: equal__is,
																	arg1: &dCollect{
																		ARG:  collect__attributes__ud_3apos,
																		arg1: &dNode{},
																	},
																	arg2: &dElem{
																		DATA: []interface{}{"ADP", "AUX"},
																		arg1: &dCollect{
																			ARG:  collect__attributes__ud_3apos,
																			arg1: &dNode{},
																		},
																	},
																},
																arg2: &dEqual{
																	ARG: equal__is,
																	arg1: &dCollect{
																		ARG:  collect__attributes__cat,
																		arg1: &dNode{},
																	},
																	arg2: &dElem{
																		DATA: []interface{}{"mwu"},
																		arg1: &dCollect{
																			ARG:  collect__attributes__cat,
																			arg1: &dNode{},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		})
	if len(conj_heads) == 0 {
		return []depT{}
	}

	//NP
	udRelation := nonLocalDependencyLabel(phrase, n1(find(q, /* ($q.varallnodes[@rel="cnj"]/
		   			    node[
		   			    (: @rel=$phrase/@rel and :)
						not(@pt or @cat) and @index=$phrase/@index])[1] */&xPath{
			arg1: &dSort{
				arg1: &dFilter{
					arg1: &dSort{
						arg1: &dCollect{
							ARG: collect__child__node,
							arg1: &dFilter{
								arg1: &dVariable{
									VAR: q.varallnodes,
								},
								arg2: &dSort{
									arg1: &dEqual{
										ARG: equal__is,
										arg1: &dCollect{
											ARG:  collect__attributes__rel,
											arg1: &dNode{},
										},
										arg2: &dElem{
											DATA: []interface{}{"cnj"},
											arg1: &dCollect{
												ARG:  collect__attributes__rel,
												arg1: &dNode{},
											},
										},
									},
								},
							},
							arg2: &dPredicate{
								arg1: &dAnd{
									arg1: &dFunction{
										ARG: function__not__1__args,
										arg1: &dArg{
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
												VAR: phrase,
											},
										},
									},
								},
							},
						},
					},
					arg2: &dSort{
						arg1: &dFunction{
							ARG: function__first__0__args,
						},
					},
				},
			},
		})), q)

	EudRelation := udRelation
	if test(q /* $udRelation = ("nmod","obl") and $phrase[@cat="pp"]//node[@ud:Relation="case" and @ud:HeadPosition=$node/@end] */, &xPath{
		arg1: &dSort{
			arg1: &dAnd{
				arg1: &dEqual{
					ARG: equal__is,
					arg1: &dVariable{
						VAR: udRelation,
					},
					arg2: &dElem{
						DATA: []interface{}{"nmod", "obl"},
						arg1: &dVariable{
							VAR: udRelation,
						},
					},
				},
				arg2: &dCollect{
					ARG: collect__child__node,
					arg1: &dCollect{
						ARG: collect__descendant__or__self__type__node,
						arg1: &dFilter{
							arg1: &dVariable{
								VAR: phrase,
							},
							arg2: &dSort{
								arg1: &dEqual{
									ARG: equal__is,
									arg1: &dCollect{
										ARG:  collect__attributes__cat,
										arg1: &dNode{},
									},
									arg2: &dElem{
										DATA: []interface{}{"pp"},
										arg1: &dCollect{
											ARG:  collect__attributes__cat,
											arg1: &dNode{},
										},
									},
								},
							},
						},
					},
					arg2: &dPredicate{
						arg1: &dAnd{
							arg1: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG:  collect__attributes__ud_3aRelation,
									arg1: &dNode{},
								},
								arg2: &dElem{
									DATA: []interface{}{"case"},
									arg1: &dCollect{
										ARG:  collect__attributes__ud_3aRelation,
										arg1: &dNode{},
									},
								},
							},
							arg2: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG:  collect__attributes__ud_3aHeadPosition,
									arg1: &dNode{},
								},
								arg2: &dCollect{
									ARG: collect__attributes__end,
									arg1: &dVariable{
										VAR: node,
									},
								},
							},
						},
					},
				},
			},
		},
	}) {
		//NP
		EudRelation = udRelation + ":" + enhancedLemmaString(find(q /* $phrase//node[@ud:Relation="case" and @ud:HeadPosition=$node/@end] */, &xPath{
			arg1: &dSort{
				arg1: &dCollect{
					ARG: collect__child__node,
					arg1: &dCollect{
						ARG: collect__descendant__or__self__type__node,
						arg1: &dVariable{
							VAR: phrase,
						},
					},
					arg2: &dPredicate{
						arg1: &dAnd{
							arg1: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG:  collect__attributes__ud_3aRelation,
									arg1: &dNode{},
								},
								arg2: &dElem{
									DATA: []interface{}{"case"},
									arg1: &dCollect{
										ARG:  collect__attributes__ud_3aRelation,
										arg1: &dNode{},
									},
								},
							},
							arg2: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG:  collect__attributes__ud_3aHeadPosition,
									arg1: &dNode{},
								},
								arg2: &dCollect{
									ARG: collect__attributes__end,
									arg1: &dVariable{
										VAR: node,
									},
								},
							},
						},
					},
				},
			},
		}), q)
	}

	result := []depT{}
	for _, conj_head := range conj_heads {
		//NP
		result = append(result, depT{head: internalHeadPosition([]interface{}{conj_head.(*nodeType)}, q), dep: EudRelation})

	}
	return result
}

// should work in coordinations like te laten reizen en te laten beleven,
// and recursive cases: Andras blijft ontkennen sexuele relaties met Timea te hebben gehad ,
//    .. of hij ook voor hen wilde komen tekenen :)
func xcompControl(node *nodeType, q *context, so_index int) []depT {

	defer func() {
		if r := recover(); r != nil {
			panic(trace(r, "xcompControl", q, node))
		}
	}()

	result := []depT{}
	for _, xcomp := range find(q, /* $node[not(@ud:PronType="Rel")]/ancestor::node//node[(@rel="vc" or (@cat="inf" and @rel="body")) (: covers inf ti oti :)
		   and node[@rel=("hd","predc") and @ud:Relation="xcomp"]  (: vrouwen moeten vertegenwoordigd zijn :)
		   and node[@rel="su" and @index]/@index = $so_index
		  ] */&xPath{
			arg1: &dSort{
				arg1: &dCollect{
					ARG: collect__child__node,
					arg1: &dCollect{
						ARG: collect__descendant__or__self__type__node,
						arg1: &dCollect{
							ARG: collect__ancestors__node,
							arg1: &dFilter{
								arg1: &dVariable{
									VAR: node,
								},
								arg2: &dSort{
									arg1: &dFunction{
										ARG: function__not__1__args,
										arg1: &dArg{
											arg1: &dSort{
												arg1: &dEqual{
													ARG: equal__is,
													arg1: &dCollect{
														ARG:  collect__attributes__ud_3aPronType,
														arg1: &dNode{},
													},
													arg2: &dElem{
														DATA: []interface{}{"Rel"},
														arg1: &dCollect{
															ARG:  collect__attributes__ud_3aPronType,
															arg1: &dNode{},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					arg2: &dPredicate{
						arg1: &dAnd{
							arg1: &dAnd{
								arg1: &dSort{
									arg1: &dOr{
										arg1: &dEqual{
											ARG: equal__is,
											arg1: &dCollect{
												ARG:  collect__attributes__rel,
												arg1: &dNode{},
											},
											arg2: &dElem{
												DATA: []interface{}{"vc"},
												arg1: &dCollect{
													ARG:  collect__attributes__rel,
													arg1: &dNode{},
												},
											},
										},
										arg2: &dSort{
											arg1: &dAnd{
												arg1: &dEqual{
													ARG: equal__is,
													arg1: &dCollect{
														ARG:  collect__attributes__cat,
														arg1: &dNode{},
													},
													arg2: &dElem{
														DATA: []interface{}{"inf"},
														arg1: &dCollect{
															ARG:  collect__attributes__cat,
															arg1: &dNode{},
														},
													},
												},
												arg2: &dEqual{
													ARG: equal__is,
													arg1: &dCollect{
														ARG:  collect__attributes__rel,
														arg1: &dNode{},
													},
													arg2: &dElem{
														DATA: []interface{}{"body"},
														arg1: &dCollect{
															ARG:  collect__attributes__rel,
															arg1: &dNode{},
														},
													},
												},
											},
										},
									},
								},
								arg2: &dCollect{
									ARG:  collect__child__node,
									arg1: &dNode{},
									arg2: &dPredicate{
										arg1: &dAnd{
											arg1: &dEqual{
												ARG: equal__is,
												arg1: &dCollect{
													ARG:  collect__attributes__rel,
													arg1: &dNode{},
												},
												arg2: &dElem{
													DATA: []interface{}{"hd", "predc"},
													arg1: &dCollect{
														ARG:  collect__attributes__rel,
														arg1: &dNode{},
													},
												},
											},
											arg2: &dEqual{
												ARG: equal__is,
												arg1: &dCollect{
													ARG:  collect__attributes__ud_3aRelation,
													arg1: &dNode{},
												},
												arg2: &dElem{
													DATA: []interface{}{"xcomp"},
													arg1: &dCollect{
														ARG:  collect__attributes__ud_3aRelation,
														arg1: &dNode{},
													},
												},
											},
										},
									},
								},
							},
							arg2: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG: collect__attributes__index,
									arg1: &dCollect{
										ARG:  collect__child__node,
										arg1: &dNode{},
										arg2: &dPredicate{
											arg1: &dAnd{
												arg1: &dEqual{
													ARG: equal__is,
													arg1: &dCollect{
														ARG:  collect__attributes__rel,
														arg1: &dNode{},
													},
													arg2: &dElem{
														DATA: []interface{}{"su"},
														arg1: &dCollect{
															ARG:  collect__attributes__rel,
															arg1: &dNode{},
														},
													},
												},
												arg2: &dCollect{
													ARG:  collect__attributes__index,
													arg1: &dNode{},
												},
											},
										},
									},
								},
								arg2: &dVariable{
									VAR: so_index,
								},
							},
						},
					},
				},
			},
		}) {
		//NP
		result = append(result, depT{head: internalHeadPosition([]interface{}{xcomp.(*nodeType)}, q), dep: "nsubj:xsubj"})
	}
	return result
}

// alpino NF specific case, controllers with extraposed content are realized downstairs
func upstairsControl(node *nodeType, q *context, so_index int) []depT {

	defer func() {
		if r := recover(); r != nil {
			panic(trace(r, "upstairsControl", q, node))
		}
	}()

	result := []depT{}
	for _, upstairs := range find(q, /* $node/ancestor::node[node[@rel="hd" and @ud:pos="VERB"]
		 and node[@rel=("su","obj1","obj2") and not(@pt or @cat)]/@index = $so_index
		] */&xPath{
			arg1: &dSort{
				arg1: &dCollect{
					ARG: collect__ancestors__node,
					arg1: &dVariable{
						VAR: node,
					},
					arg2: &dPredicate{
						arg1: &dAnd{
							arg1: &dCollect{
								ARG:  collect__child__node,
								arg1: &dNode{},
								arg2: &dPredicate{
									arg1: &dAnd{
										arg1: &dEqual{
											ARG: equal__is,
											arg1: &dCollect{
												ARG:  collect__attributes__rel,
												arg1: &dNode{},
											},
											arg2: &dElem{
												DATA: []interface{}{"hd"},
												arg1: &dCollect{
													ARG:  collect__attributes__rel,
													arg1: &dNode{},
												},
											},
										},
										arg2: &dEqual{
											ARG: equal__is,
											arg1: &dCollect{
												ARG:  collect__attributes__ud_3apos,
												arg1: &dNode{},
											},
											arg2: &dElem{
												DATA: []interface{}{"VERB"},
												arg1: &dCollect{
													ARG:  collect__attributes__ud_3apos,
													arg1: &dNode{},
												},
											},
										},
									},
								},
							},
							arg2: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG: collect__attributes__index,
									arg1: &dCollect{
										ARG:  collect__child__node,
										arg1: &dNode{},
										arg2: &dPredicate{
											arg1: &dAnd{
												arg1: &dEqual{
													ARG: equal__is,
													arg1: &dCollect{
														ARG:  collect__attributes__rel,
														arg1: &dNode{},
													},
													arg2: &dElem{
														DATA: []interface{}{"su", "obj1", "obj2"},
														arg1: &dCollect{
															ARG:  collect__attributes__rel,
															arg1: &dNode{},
														},
													},
												},
												arg2: &dFunction{
													ARG: function__not__1__args,
													arg1: &dArg{
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
													},
												},
											},
										},
									},
								},
								arg2: &dVariable{
									VAR: so_index,
								},
							},
						},
					},
				},
			},
		}) {
		//NP
		result = append(result, depT{head: internalHeadPosition([]interface{}{upstairs.(*nodeType)}, q), dep: "nsubj:xsubj"})
	}
	return result

}

// een koers waarin de Alsemberg moet worden beklommen
func passiveVpControl(node *nodeType, q *context, so_index int) []depT {

	defer func() {
		if r := recover(); r != nil {
			panic(trace(r, "passiveVpControl", q, node))
		}
	}()

	result := []depT{}
	for _, passive_vp := range find(q, /* $q.varallnodes[@rel="vc" and @cat="ppart"
		   and node[@rel="hd" and @ud:Relation="xcomp"]
		   and node[@rel="obj1" and @index]/@index = $so_index ] */&xPath{
			arg1: &dSort{
				arg1: &dFilter{
					arg1: &dVariable{
						VAR: q.varallnodes,
					},
					arg2: &dSort{
						arg1: &dAnd{
							arg1: &dAnd{
								arg1: &dAnd{
									arg1: &dEqual{
										ARG: equal__is,
										arg1: &dCollect{
											ARG:  collect__attributes__rel,
											arg1: &dNode{},
										},
										arg2: &dElem{
											DATA: []interface{}{"vc"},
											arg1: &dCollect{
												ARG:  collect__attributes__rel,
												arg1: &dNode{},
											},
										},
									},
									arg2: &dEqual{
										ARG: equal__is,
										arg1: &dCollect{
											ARG:  collect__attributes__cat,
											arg1: &dNode{},
										},
										arg2: &dElem{
											DATA: []interface{}{"ppart"},
											arg1: &dCollect{
												ARG:  collect__attributes__cat,
												arg1: &dNode{},
											},
										},
									},
								},
								arg2: &dCollect{
									ARG:  collect__child__node,
									arg1: &dNode{},
									arg2: &dPredicate{
										arg1: &dAnd{
											arg1: &dEqual{
												ARG: equal__is,
												arg1: &dCollect{
													ARG:  collect__attributes__rel,
													arg1: &dNode{},
												},
												arg2: &dElem{
													DATA: []interface{}{"hd"},
													arg1: &dCollect{
														ARG:  collect__attributes__rel,
														arg1: &dNode{},
													},
												},
											},
											arg2: &dEqual{
												ARG: equal__is,
												arg1: &dCollect{
													ARG:  collect__attributes__ud_3aRelation,
													arg1: &dNode{},
												},
												arg2: &dElem{
													DATA: []interface{}{"xcomp"},
													arg1: &dCollect{
														ARG:  collect__attributes__ud_3aRelation,
														arg1: &dNode{},
													},
												},
											},
										},
									},
								},
							},
							arg2: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG: collect__attributes__index,
									arg1: &dCollect{
										ARG:  collect__child__node,
										arg1: &dNode{},
										arg2: &dPredicate{
											arg1: &dAnd{
												arg1: &dEqual{
													ARG: equal__is,
													arg1: &dCollect{
														ARG:  collect__attributes__rel,
														arg1: &dNode{},
													},
													arg2: &dElem{
														DATA: []interface{}{"obj1"},
														arg1: &dCollect{
															ARG:  collect__attributes__rel,
															arg1: &dNode{},
														},
													},
												},
												arg2: &dCollect{
													ARG:  collect__attributes__index,
													arg1: &dNode{},
												},
											},
										},
									},
								},
								arg2: &dVariable{
									VAR: so_index,
								},
							},
						},
					},
				},
			},
		}) {
		//NP
		result = append(result, depT{head: internalHeadPosition([]interface{}{passive_vp.(*nodeType)}, q), dep: "nsubj:pass:xsubj"})
	}

	return result
}

func enhancedLemmaString(nodes []interface{}, q *context) string {
	sort.Slice(nodes, func(i, j int) bool {
		// TODO: nodeType heeft geen head
		return nodes[i].(*nodeType).udEHeadPosition < nodes[j].(*nodeType).udEHeadPosition
	})
	lemmas := make([]string, len(nodes))
	for i, node := range nodes {
		lemmas[i] = enhancedLemmaString1(node.(*nodeType), q)
	}
	return strings.Join(lemmas, "_")
}

func enhancedLemmaString1(node *nodeType, q *context) string {
	var lemma string
	switch node.Lemma {
	case "a.k.a":
		lemma = "also_known_as"
	case "c.q.":
		lemma = "casu_quo"
	case "dwz.", "d.w.z.":
		lemma = "dat_wil_zeggen"
	case "e.d.":
		lemma = "en_dergelijke"
	case "en/of":
		lemma = "en_of"
	case "enz.":
		lemma = "enzovoort"
	case "etc.":
		lemma = "etcetera"
	case "m.a.w.":
		lemma = "met_andere_woorden"
	case "nl.":
		lemma = "namelijk"
	case "resp.":
		lemma = "respectievelijk"
	case "t/m":
		lemma = "tot_en_met"
	case "t.a.v.":
		lemma = "ten_aanzien_van"
	case "t.g.v.":
		lemma = "ten_gunste_van"
	case "t.n.v.":
		lemma = "ten_name_van"
	case "t.o.v.":
		lemma = "ten_opzichte_van"
	default:
		lemma = node.Lemma
	}
	fixed := find(q /* $node/../node[@ud:ERelation="fixed"] */, &xPath{
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
					arg1: &dEqual{
						ARG: equal__is,
						arg1: &dCollect{
							ARG:  collect__attributes__ud_3aERelation,
							arg1: &dNode{},
						},
						arg2: &dElem{
							DATA: []interface{}{"fixed"},
							arg1: &dCollect{
								ARG:  collect__attributes__ud_3aERelation,
								arg1: &dNode{},
							},
						},
					},
				},
			},
		},
	})
	if len(fixed) > 0 {
		sort.Slice(fixed, func(i, j int) bool {
			fi := fixed[i].(*nodeType)
			fj := fixed[j].(*nodeType)
			ei := fi.End
			ej := fj.End
			if fi.udCopiedFrom > 0 {
				ei = fi.udCopiedFrom
			}
			if fj.udCopiedFrom > 0 {
				ej = fj.udCopiedFrom
			}
			return ei < ej
		})
		for _, f := range fixed {
			lemma += "_" + f.(*nodeType).Lemma
		}
	}
	lemma = strings.Replace(lemma, "/", "schuine_streep", -1)
	lemma = strings.Replace(lemma, "-", "_", -1)
	return strings.ToLower(lemma)
}
