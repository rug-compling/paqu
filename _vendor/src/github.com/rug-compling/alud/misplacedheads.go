//
// GENERATED FILE -- DO NOT EDIT
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

	seen := make(map[[2]int]int)
	counter := 0

START:
	for true {
		for _, n1 := range q.varallnodes {
			// FIND op varallnodes niet mogelijk omdat twee keer naar $node wordt verwezen, en dat moet dezelfde node zijn
			for _, n2 := range find(q, /*
				$n1[@rel=("hd","ld") and
				      @index and
				      (@pt or @cat) and
				      ancestor::node[@rel="cnj"] and
				      ancestor::node[@cat="conj"]/node[@rel="cnj" and
				                                       descendant-or-self::node[@rel=("hd","ld") and
				                                                                @index=$n1/@index and
				                                                                not(@cat or @pt) and
				                                                                ( @begin        = ..//node[@cat or @pt]/@end or
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
																			DATA: []interface{}{"hd", "ld"},
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
					$q.varallnodes[@rel=("hd","ld","vc") and @index and not(@pt or @cat) and
					                 ancestor::node[@rel="cnj"]  and
					                                    ( @begin        = ..//node[@cat or @pt]/@end or
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
						pair := [2]int{node2.Id, node3.Id}
						if i, ok := seen[pair]; ok {
							if i == 1 {
								panic(fmt.Sprintf("Loop detected in fixMisplacedHeadsInCoordination: %d %d", node2.Id, node3.Id))
							}
							seen[pair]++
							continue
						}
						seen[pair] = 1
						counter++
						// kopieer inhoud van node2 (niet leeg) naar node3 (leeg)
						id, rel := node3.Id, node3.Rel
						*node3 = *node2
						node3.Id, node3.Rel = id, rel
						// maak node2 leeg
						*node2 = nodeType{
							Begin: node2.Begin,
							End:   node2.End,
							Id:    node2.Id,
							Index: node2.Index,
							Rel:   node2.Rel,
							Node:  []*nodeType{},
						}
						// opnieuw beginnen
						inspect(q)
						continue START
					}
				}
			}
		}
		break
	}
	if counter > 0 {
		q.debugs = append(q.debugs, fmt.Sprintf("fixMisplacedHeadsInCoordination: %d swaps", counter))
	}
}
