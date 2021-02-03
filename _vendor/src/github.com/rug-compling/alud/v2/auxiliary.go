//
// // THIS IS A GENERATED FILE. DO NOT EDIT.
//

package alud

import (
	"fmt"
)

/*
func auxiliary(nodes []*nodeType, q *context) string {
	if len(nodes) != 1 { // TODO: in script staat: = 0
		return "ERROR_AUXILIARY_FUNCTION_TAKES_EXACTLY_ONE_ARG"
	}
	return auxiliary1(nodes[0], q)
}
*/

func auxiliary1(node *nodeType, q *context) (aux string, err error) {

	if node.Pt != "ww" {
		return "", fmt.Errorf("ERROR_NO_VERB")
	}
	if node.Rel != "hd" {
		return "verb", nil
	}

	if test(q, /* $node[not(../node[@rel=("obj1","se","vc")]) and
			        (: ud documentation suggests 1 cop per lg, van Eynde suggests much more, compromise: the traditional ones :)
			        (: @lemma=("zijn","lijken","blijken","blijven","schijnen","heten","voorkomen","worden","dunken") and :)
			        @lemma="zijn" and
		             (:    ( contains(@sc,'copula') or
		                   contains(@sc,'pred')   or
		                   contains(@sc,'cleft')  or :)
		                   ../node[@rel="predc"]
		                  ] */&xPath{
			arg1: &dSort{
				arg1: &dFilter{
					arg1: &dVariable{
						VAR: node,
					},
					arg2: &dSort{
						arg1: &dAnd{
							arg1: &dAnd{
								arg1: &dFunction{
									ARG: function__not__1__args,
									arg1: &dArg{
										arg1: &dSort{
											arg1: &dCollect{
												ARG: collect__child__node,
												arg1: &dCollect{
													ARG:  collect__parent__type__node,
													arg1: &dNode{},
												},
												arg2: &dPredicate{
													arg1: &dEqual{
														ARG: equal__is,
														arg1: &dCollect{
															ARG:  collect__attributes__rel,
															arg1: &dNode{},
														},
														arg2: &dElem{
															DATA: []interface{}{"obj1", "se", "vc"},
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
								},
								arg2: &dEqual{
									ARG: equal__is,
									arg1: &dCollect{
										ARG:  collect__attributes__lemma,
										arg1: &dNode{},
									},
									arg2: &dElem{
										DATA: []interface{}{"zijn"},
										arg1: &dCollect{
											ARG:  collect__attributes__lemma,
											arg1: &dNode{},
										},
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
									arg1: &dEqual{
										ARG: equal__is,
										arg1: &dCollect{
											ARG:  collect__attributes__rel,
											arg1: &dNode{},
										},
										arg2: &dElem{
											DATA: []interface{}{"predc"},
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
				},
			},
		}) {
		return "cop", nil
	}

	if test(q, /* $node[@lemma=("zijn","worden") and
		   ( @sc="passive"  or
		   	 ( ../node[@rel="vc"] and
		         ( ../node[@rel="su"]/@index = ../node[@rel="vc"]/node[@rel="obj1"]/@index or
		           ../node[@rel="su"]/@index = ../node[@rel="vc"]/node[@rel="cnj"]/node[@rel="obj1"]/@index or
		           ../node[@rel="vc" and not(@pt or @cat)]/@index =
		               $q.varindexnodes[@rel="vc" and node[@rel="obj1"]/@index = $node/../node[@rel="su"]/@index]/@index
		          or not(../node[@rel="su"])
		          )
		      )
		   ) ] */&xPath{
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
									ARG:  collect__attributes__lemma,
									arg1: &dNode{},
								},
								arg2: &dElem{
									DATA: []interface{}{"zijn", "worden"},
									arg1: &dCollect{
										ARG:  collect__attributes__lemma,
										arg1: &dNode{},
									},
								},
							},
							arg2: &dSort{
								arg1: &dOr{
									arg1: &dEqual{
										ARG: equal__is,
										arg1: &dCollect{
											ARG:  collect__attributes__sc,
											arg1: &dNode{},
										},
										arg2: &dElem{
											DATA: []interface{}{"passive"},
											arg1: &dCollect{
												ARG:  collect__attributes__sc,
												arg1: &dNode{},
											},
										},
									},
									arg2: &dSort{
										arg1: &dAnd{
											arg1: &dCollect{
												ARG: collect__child__node,
												arg1: &dCollect{
													ARG:  collect__parent__type__node,
													arg1: &dNode{},
												},
												arg2: &dPredicate{
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
												},
											},
											arg2: &dSort{
												arg1: &dOr{
													arg1: &dOr{
														arg1: &dOr{
															arg1: &dEqual{
																ARG: equal__is,
																arg1: &dCollect{
																	ARG: collect__attributes__index,
																	arg1: &dCollect{
																		ARG: collect__child__node,
																		arg1: &dCollect{
																			ARG:  collect__parent__type__node,
																			arg1: &dNode{},
																		},
																		arg2: &dPredicate{
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
																		},
																	},
																},
																arg2: &dCollect{
																	ARG: collect__attributes__index,
																	arg1: &dCollect{
																		ARG: collect__child__node,
																		arg1: &dCollect{
																			ARG: collect__child__node,
																			arg1: &dCollect{
																				ARG:  collect__parent__type__node,
																				arg1: &dNode{},
																			},
																			arg2: &dPredicate{
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
																					DATA: []interface{}{"obj1"},
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
															arg2: &dEqual{
																ARG: equal__is,
																arg1: &dCollect{
																	ARG: collect__attributes__index,
																	arg1: &dCollect{
																		ARG: collect__child__node,
																		arg1: &dCollect{
																			ARG:  collect__parent__type__node,
																			arg1: &dNode{},
																		},
																		arg2: &dPredicate{
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
																		},
																	},
																},
																arg2: &dCollect{
																	ARG: collect__attributes__index,
																	arg1: &dCollect{
																		ARG: collect__child__node,
																		arg1: &dCollect{
																			ARG: collect__child__node,
																			arg1: &dCollect{
																				ARG: collect__child__node,
																				arg1: &dCollect{
																					ARG:  collect__parent__type__node,
																					arg1: &dNode{},
																				},
																				arg2: &dPredicate{
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
																					DATA: []interface{}{"vc"},
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
															arg2: &dCollect{
																ARG: collect__attributes__index,
																arg1: &dFilter{
																	arg1: &dVariable{
																		VAR: q.varindexnodes,
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
																					ARG: collect__attributes__index,
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
																									DATA: []interface{}{"obj1"},
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
																					ARG: collect__attributes__index,
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
													arg2: &dFunction{
														ARG: function__not__1__args,
														arg1: &dArg{
															arg1: &dSort{
																arg1: &dCollect{
																	ARG: collect__child__node,
																	arg1: &dCollect{
																		ARG:  collect__parent__type__node,
																		arg1: &dNode{},
																	},
																	arg2: &dPredicate{
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
		return "aux:pass", nil
	}

	// krijgen passive with iobj control
	if test(q, /* $node[@lemma="krijgen" and
			              ( ../node[@rel="su"]/@index = ../node[@rel="vc"]/node[@rel="obj2"]/@index or
		                  ../node[@rel="su"]/@index = ../node[@rel="vc"]/node[@rel="cnj"]/node[@rel="obj2"]/@index
		                )] */&xPath{
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
									ARG:  collect__attributes__lemma,
									arg1: &dNode{},
								},
								arg2: &dElem{
									DATA: []interface{}{"krijgen"},
									arg1: &dCollect{
										ARG:  collect__attributes__lemma,
										arg1: &dNode{},
									},
								},
							},
							arg2: &dSort{
								arg1: &dOr{
									arg1: &dEqual{
										ARG: equal__is,
										arg1: &dCollect{
											ARG: collect__attributes__index,
											arg1: &dCollect{
												ARG: collect__child__node,
												arg1: &dCollect{
													ARG:  collect__parent__type__node,
													arg1: &dNode{},
												},
												arg2: &dPredicate{
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
												},
											},
										},
										arg2: &dCollect{
											ARG: collect__attributes__index,
											arg1: &dCollect{
												ARG: collect__child__node,
												arg1: &dCollect{
													ARG: collect__child__node,
													arg1: &dCollect{
														ARG:  collect__parent__type__node,
														arg1: &dNode{},
													},
													arg2: &dPredicate{
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
															DATA: []interface{}{"obj2"},
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
									arg2: &dEqual{
										ARG: equal__is,
										arg1: &dCollect{
											ARG: collect__attributes__index,
											arg1: &dCollect{
												ARG: collect__child__node,
												arg1: &dCollect{
													ARG:  collect__parent__type__node,
													arg1: &dNode{},
												},
												arg2: &dPredicate{
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
												},
											},
										},
										arg2: &dCollect{
											ARG: collect__attributes__index,
											arg1: &dCollect{
												ARG: collect__child__node,
												arg1: &dCollect{
													ARG: collect__child__node,
													arg1: &dCollect{
														ARG: collect__child__node,
														arg1: &dCollect{
															ARG:  collect__parent__type__node,
															arg1: &dNode{},
														},
														arg2: &dPredicate{
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
													arg1: &dEqual{
														ARG: equal__is,
														arg1: &dCollect{
															ARG:  collect__attributes__rel,
															arg1: &dNode{},
														},
														arg2: &dElem{
															DATA: []interface{}{"obj2"},
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
								},
							},
						},
					},
				},
			},
		}) {
		return "aux:pass", nil
	}

	// alpino has no principled distinction between AUX and VERB, should be TAME verbs semantically, we follow ENGLISH
	// blijken and hoeven removed from list
	if test(q, /* $node[not(../node[@rel="predc"]) and  (: hij heeft als opdracht stammen uit elkaar te houden  :)
		   ( starts-with(@sc,'aux') or
		     ( ../node[@rel="vc"  and
		                ( @cat=("ppart","inf","ti") or
		                  ( @cat="conj" and node[@rel="cnj" and @cat=("ppart","inf","ti")] ) or
		                  ( @index and not(@pt or @cat))  (: dangling aux in gapped coordination :)
		                )
		              ]   and
		       @lemma=("hebben","kunnen","moeten","mogen","zijn","zullen")
		     )
		   )
		 ] */&xPath{
			arg1: &dSort{
				arg1: &dFilter{
					arg1: &dVariable{
						VAR: node,
					},
					arg2: &dSort{
						arg1: &dAnd{
							arg1: &dFunction{
								ARG: function__not__1__args,
								arg1: &dArg{
									arg1: &dSort{
										arg1: &dCollect{
											ARG: collect__child__node,
											arg1: &dCollect{
												ARG:  collect__parent__type__node,
												arg1: &dNode{},
											},
											arg2: &dPredicate{
												arg1: &dEqual{
													ARG: equal__is,
													arg1: &dCollect{
														ARG:  collect__attributes__rel,
														arg1: &dNode{},
													},
													arg2: &dElem{
														DATA: []interface{}{"predc"},
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
							},
							arg2: &dSort{
								arg1: &dOr{
									arg1: &dFunction{
										ARG: function__starts__with__2__args,
										arg1: &dArg{
											arg1: &dArg{
												arg1: &dSort{
													arg1: &dCollect{
														ARG:  collect__attributes__sc,
														arg1: &dNode{},
													},
												},
											},
											arg2: &dElem{
												DATA: []interface{}{"aux"},
											},
										},
									},
									arg2: &dSort{
										arg1: &dAnd{
											arg1: &dCollect{
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
																DATA: []interface{}{"vc"},
																arg1: &dCollect{
																	ARG:  collect__attributes__rel,
																	arg1: &dNode{},
																},
															},
														},
														arg2: &dSort{
															arg1: &dOr{
																arg1: &dOr{
																	arg1: &dEqual{
																		ARG: equal__is,
																		arg1: &dCollect{
																			ARG:  collect__attributes__cat,
																			arg1: &dNode{},
																		},
																		arg2: &dElem{
																			DATA: []interface{}{"ppart", "inf", "ti"},
																			arg1: &dCollect{
																				ARG:  collect__attributes__cat,
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
																					DATA: []interface{}{"conj"},
																					arg1: &dCollect{
																						ARG:  collect__attributes__cat,
																						arg1: &dNode{},
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
																								ARG:  collect__attributes__cat,
																								arg1: &dNode{},
																							},
																							arg2: &dElem{
																								DATA: []interface{}{"ppart", "inf", "ti"},
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
																arg2: &dSort{
																	arg1: &dAnd{
																		arg1: &dCollect{
																			ARG:  collect__attributes__index,
																			arg1: &dNode{},
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
													},
												},
											},
											arg2: &dEqual{
												ARG: equal__is,
												arg1: &dCollect{
													ARG:  collect__attributes__lemma,
													arg1: &dNode{},
												},
												arg2: &dElem{
													DATA: []interface{}{"hebben", "kunnen", "moeten", "mogen", "zijn", "zullen"},
													arg1: &dCollect{
														ARG:  collect__attributes__lemma,
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
		}) {
		return "aux", nil
	}

	return "verb", nil
}
