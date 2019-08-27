//
// // THIS IS A GENERATED FILE. DO NOT EDIT.
//

package alud

import (
	"fmt"
)

// voorkwam dat LPF opnieuw of SGP voor het eerst in de regering zou komen  -- gapped LD
func fixMisplacedHeadsInCoordination(q *context) {

	if len(q.varindexnodes) == 0 {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			panic(trace(r, "fixMisplacedHeadsInCoordination", q))
		}
	}()

	seen := make(map[[2]int]bool)

START:
	for true {
		for _, n1 := range q.varallnodes {
			// FIND op varallnodes niet mogelijk omdat twee keer naar $node wordt verwezen, en dat moet dezelfde node zijn
			for _, n2 := range find(q, /*
				$n1[(@rel=("hd","ld") or (@rel="obj1" and ../node[@rel="hd" and @pt="vz"])) and
				      @index and
				      (@pt or @cat) and
				      ancestor::node[@rel="cnj"] and
				      ancestor::node[@cat="conj"]/node[@rel="cnj" and
				                                       descendant-or-self::node[@rel=("hd","ld","obj1") and
				                                                                @index=$n1/@index and
				                                                                not(@cat or @pt) and
				                                                                ( @begin        = ..//node[@cat or @pt]/@end or
				                                                                  @begin        = ../..//node[@cat or @pt]/@end or
				                                                                  @begin - 1000 = ..//node[@cat or @pt]/@end
				                                                                )
				                                                               ]
				                                       ]] */&xPath{
					arg1: &dSort{
						arg1: &dFilter{
							arg1: &dVariable{
								VAR: n1,
							},
							arg2: &dSort{
								arg1: &dAnd{
									arg1: &dAnd{
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
																DATA: []interface{}{"hd", "ld"},
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
																	ARG: collect__child__node,
																	arg1: &dCollect{
																		ARG:  collect__parent__type__node,
																		arg1: &dNode{},
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
																					ARG:  collect__attributes__pt,
																					arg1: &dNode{},
																				},
																				arg2: &dElem{
																					DATA: []interface{}{"vz"},
																					arg1: &dCollect{
																						ARG:  collect__attributes__pt,
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
												arg2: &dCollect{
													ARG:  collect__attributes__index,
													arg1: &dNode{},
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
										arg2: &dCollect{
											ARG:  collect__ancestors__node,
											arg1: &dNode{},
											arg2: &dPredicate{
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
									},
									arg2: &dCollect{
										ARG: collect__child__node,
										arg1: &dCollect{
											ARG:  collect__ancestors__node,
											arg1: &dNode{},
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
										arg2: &dPredicate{
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
												arg2: &dCollect{
													ARG:  collect__descendant__or__self__node,
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
																			DATA: []interface{}{"hd", "ld", "obj1"},
																			arg1: &dCollect{
																				ARG:  collect__attributes__rel,
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
																				VAR: n1,
																			},
																		},
																	},
																},
																arg2: &dFunction{
																	ARG: function__not__1__args,
																	arg1: &dArg{
																		arg1: &dSort{
																			arg1: &dOr{
																				arg1: &dCollect{
																					ARG:  collect__attributes__cat,
																					arg1: &dNode{},
																				},
																				arg2: &dCollect{
																					ARG:  collect__attributes__pt,
																					arg1: &dNode{},
																				},
																			},
																		},
																	},
																},
															},
															arg2: &dSort{
																arg1: &dOr{
																	arg1: &dOr{
																		arg1: &dEqual{
																			ARG: equal__is,
																			arg1: &dCollect{
																				ARG:  collect__attributes__begin,
																				arg1: &dNode{},
																			},
																			arg2: &dCollect{
																				ARG: collect__attributes__end,
																				arg1: &dCollect{
																					ARG: collect__child__node,
																					arg1: &dCollect{
																						ARG: collect__descendant__or__self__type__node,
																						arg1: &dCollect{
																							ARG:  collect__parent__type__node,
																							arg1: &dNode{},
																						},
																					},
																					arg2: &dPredicate{
																						arg1: &dOr{
																							arg1: &dCollect{
																								ARG:  collect__attributes__cat,
																								arg1: &dNode{},
																							},
																							arg2: &dCollect{
																								ARG:  collect__attributes__pt,
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
																				ARG:  collect__attributes__begin,
																				arg1: &dNode{},
																			},
																			arg2: &dCollect{
																				ARG: collect__attributes__end,
																				arg1: &dCollect{
																					ARG: collect__child__node,
																					arg1: &dCollect{
																						ARG: collect__descendant__or__self__type__node,
																						arg1: &dCollect{
																							ARG: collect__parent__type__node,
																							arg1: &dCollect{
																								ARG:  collect__parent__type__node,
																								arg1: &dNode{},
																							},
																						},
																					},
																					arg2: &dPredicate{
																						arg1: &dOr{
																							arg1: &dCollect{
																								ARG:  collect__attributes__cat,
																								arg1: &dNode{},
																							},
																							arg2: &dCollect{
																								ARG:  collect__attributes__pt,
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
																		arg1: &dPlus{
																			ARG: plus__minus,
																			arg1: &dCollect{
																				ARG:  collect__attributes__begin,
																				arg1: &dNode{},
																			},
																			arg2: &dElem{
																				DATA: []interface{}{1000},
																				arg1: &dCollect{
																					ARG:  collect__attributes__begin,
																					arg1: &dNode{},
																				},
																			},
																		},
																		arg2: &dCollect{
																			ARG: collect__attributes__end,
																			arg1: &dCollect{
																				ARG: collect__child__node,
																				arg1: &dCollect{
																					ARG: collect__descendant__or__self__type__node,
																					arg1: &dCollect{
																						ARG:  collect__parent__type__node,
																						arg1: &dNode{},
																					},
																				},
																				arg2: &dPredicate{
																					arg1: &dOr{
																						arg1: &dCollect{
																							ARG:  collect__attributes__cat,
																							arg1: &dNode{},
																						},
																						arg2: &dCollect{
																							ARG:  collect__attributes__pt,
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
						},
					},
				}) {
				node2 := n2.(*nodeType)
				for _, n3 := range find(q, /*
					$q.varallnodes[(@rel=("hd","ld","vc") or (@rel="obj1" and ../node[@rel="hd" and @pt="vz"])) and @index and not(@pt or @cat) and
					                 ancestor::node[@rel="cnj"]  and
					                                    ( @begin        = ..//node[@cat or @pt]/@end or
					                                      @begin        = ../..//node[@cat or @pt]/@end or
					                                      @begin - 1000 = ..//node[@cat or @pt]/@end
					                                     )] */&xPath{
						arg1: &dSort{
							arg1: &dFilter{
								arg1: &dVariable{
									VAR: q.varallnodes,
								},
								arg2: &dSort{
									arg1: &dAnd{
										arg1: &dAnd{
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
																	DATA: []interface{}{"hd", "ld", "vc"},
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
																		ARG: collect__child__node,
																		arg1: &dCollect{
																			ARG:  collect__parent__type__node,
																			arg1: &dNode{},
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
																						ARG:  collect__attributes__pt,
																						arg1: &dNode{},
																					},
																					arg2: &dElem{
																						DATA: []interface{}{"vz"},
																						arg1: &dCollect{
																							ARG:  collect__attributes__pt,
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
													arg2: &dCollect{
														ARG:  collect__attributes__index,
														arg1: &dNode{},
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
											arg2: &dCollect{
												ARG:  collect__ancestors__node,
												arg1: &dNode{},
												arg2: &dPredicate{
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
										},
										arg2: &dSort{
											arg1: &dOr{
												arg1: &dOr{
													arg1: &dEqual{
														ARG: equal__is,
														arg1: &dCollect{
															ARG:  collect__attributes__begin,
															arg1: &dNode{},
														},
														arg2: &dCollect{
															ARG: collect__attributes__end,
															arg1: &dCollect{
																ARG: collect__child__node,
																arg1: &dCollect{
																	ARG: collect__descendant__or__self__type__node,
																	arg1: &dCollect{
																		ARG:  collect__parent__type__node,
																		arg1: &dNode{},
																	},
																},
																arg2: &dPredicate{
																	arg1: &dOr{
																		arg1: &dCollect{
																			ARG:  collect__attributes__cat,
																			arg1: &dNode{},
																		},
																		arg2: &dCollect{
																			ARG:  collect__attributes__pt,
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
															ARG:  collect__attributes__begin,
															arg1: &dNode{},
														},
														arg2: &dCollect{
															ARG: collect__attributes__end,
															arg1: &dCollect{
																ARG: collect__child__node,
																arg1: &dCollect{
																	ARG: collect__descendant__or__self__type__node,
																	arg1: &dCollect{
																		ARG: collect__parent__type__node,
																		arg1: &dCollect{
																			ARG:  collect__parent__type__node,
																			arg1: &dNode{},
																		},
																	},
																},
																arg2: &dPredicate{
																	arg1: &dOr{
																		arg1: &dCollect{
																			ARG:  collect__attributes__cat,
																			arg1: &dNode{},
																		},
																		arg2: &dCollect{
																			ARG:  collect__attributes__pt,
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
													arg1: &dPlus{
														ARG: plus__minus,
														arg1: &dCollect{
															ARG:  collect__attributes__begin,
															arg1: &dNode{},
														},
														arg2: &dElem{
															DATA: []interface{}{1000},
															arg1: &dCollect{
																ARG:  collect__attributes__begin,
																arg1: &dNode{},
															},
														},
													},
													arg2: &dCollect{
														ARG: collect__attributes__end,
														arg1: &dCollect{
															ARG: collect__child__node,
															arg1: &dCollect{
																ARG: collect__descendant__or__self__type__node,
																arg1: &dCollect{
																	ARG:  collect__parent__type__node,
																	arg1: &dNode{},
																},
															},
															arg2: &dPredicate{
																arg1: &dOr{
																	arg1: &dCollect{
																		ARG:  collect__attributes__cat,
																		arg1: &dNode{},
																	},
																	arg2: &dCollect{
																		ARG:  collect__attributes__pt,
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
					}) {
					node3 := n3.(*nodeType)
					if node2.Index == node3.Index {

						q.debugs = append(q.debugs, fmt.Sprintf(
							"fixMisplacedHeadsInCoordination: %d -> %d ?", node2.Id, node3.Id))
						min, max := minmaxword(node2)
						score2 := wordwalk(q.alpino.Node, node2, node3, min, max)
						score3 := wordwalk(q.alpino.Node, node3, node2, min, max)
						if score2 <= score3 {
							q.debugs = append(q.debugs, fmt.Sprintf("  -> reject, scores: %d %d", score2, score3))
							if score2 == score3 {
								q.debugs = append(q.debugs, "  -> warning: equal score")
							}
							continue
						}
						q.debugs = append(q.debugs, fmt.Sprintf("  -> accept, scores: %d %d", score2, score3))

						pair := [2]int{node2.Id, node3.Id}
						if seen[pair] {
							panic(fmt.Sprintf("Loop detected in fixMisplacedHeadsInCoordination: %d -> %d", node2.Id, node3.Id))
						}
						seen[pair] = true
						// kopieer inhoud van node2 (niet leeg) naar node3 (leeg)
						swap(node2, node3)
						q.swapped = append(q.swapped, [2]*nodeType{node2, node3})
						// opnieuw beginnen
						inspect(q)
						continue START
					}
				}
			}
		}
		break
	}
}

func swap(nietLeeg, leeg *nodeType) {
	// kopieer nietLeeg naar leeg
	id, rel := leeg.Id, leeg.Rel
	*leeg = *nietLeeg
	leeg.Id, leeg.Rel = id, rel
	// maak nietLeeg leeg
	*nietLeeg = nodeType{
		Begin: nietLeeg.Begin,
		End:   nietLeeg.End,
		Id:    nietLeeg.Id,
		Index: nietLeeg.Index,
		Rel:   nietLeeg.Rel,
		Node:  []*nodeType{},
	}
}

func minmaxword(node *nodeType) (min, max int) {
	min = 9999000
	max = 0
	var r func(*nodeType)
	r = func(n *nodeType) {
		if n.Word != "" {
			if n.End > max {
				max = n.End
			}
			if n.End < min {
				min = n.End
			}
		} else if n.Node != nil {
			for _, n1 := range n.Node {
				r(n1)
			}
		}
	}
	r(node)

	return min, max
}

func wordwalk(root, node, skip *nodeType, min, max int) (score int) {

	left := true
	e1 := node.Begin + 1000
	e2 := node.End
	var r func(*nodeType)
	r = func(n *nodeType) {
		if n == node {
			left = false
			return
		}
		if n == skip {
			return
		}
		if n.Word != "" {
			if left {
				if n.End > min {
					score++
				}
			} else {
				if n.End < max {
					score++
				}
			}
			return
		}
		if n.Node == nil {
			return
		}
		if left && n.End < e1 || !left && n.Begin+1000 > e2 {
			return
		}
		for _, n2 := range n.Node {
			r(n2)
		}
	}
	r(root)
	return
}
