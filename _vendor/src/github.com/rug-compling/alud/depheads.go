//
// GENERATED FILE -- DO NOT EDIT
//

package alud

import (
	"fmt"
	"sort"
)

// recursive
func externalHeadPosition(nodes []interface{}, q *context, tr []trace) int {
	depthCheck(q, "externalHeadPosition")

	if len(nodes) == 0 {
		tr = append(tr, trace{s: "externalHeadPosition"})
		panic(tracer(fmt.Sprint("External head must have one arg, has ", len(nodes)), tr, q))
	}

	node := nodes[0].(*nodeType)
	tr = append(tr, trace{s: "externalHeadPosition", node: node})

	if node.Rel == "hd" && (node.udPos == "ADP" || node.parent.Cat == "pp") {
		// vol vertrouwen
		if n := find(q /* $node/../node[@rel="predc"] */, &xPath{
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
		}); len(n) > 0 {
			// met als titel
			return internalHeadPosition(n[:1], q, tr)
		}
		if obj1_vc_se_me := find(q /* $node/../node[@rel=("obj1","vc","se","me")] */, &xPath{
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
								ARG:  collect__attributes__rel,
								arg1: &dNode{},
							},
							arg2: &dElem{
								DATA: []interface{}{"obj1", "vc", "se", "me"},
								arg1: &dCollect{
									ARG:  collect__attributes__rel,
									arg1: &dNode{},
								},
							},
						},
					},
				},
			},
		}); len(obj1_vc_se_me) > 0 {
			// adding pt/cat enough for gapping cases?
			if test(q /* $obj1_vc_se_me[@pt or @cat] */, &xPath{
				arg1: &dSort{
					arg1: &dFilter{
						arg1: &dVariable{
							VAR: obj1_vc_se_me,
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
				},
			}) {
				return internalHeadPositionWithGapping(if1(obj1_vc_se_me), q, tr)
			}
			if test(q /* $obj1_vc_se_me[@index = ancestor::node/node[@rel=("rhd","whd")]/@index] */, &xPath{
				arg1: &dSort{
					arg1: &dFilter{
						arg1: &dVariable{
							VAR: obj1_vc_se_me,
						},
						arg2: &dSort{
							arg1: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG:  collect__attributes__index,
									arg1: &dNode{},
								},
								arg2: &dCollect{
									ARG: collect__attributes__index,
									arg1: &dCollect{
										ARG: collect__child__node,
										arg1: &dCollect{
											ARG:  collect__ancestors__node,
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
													DATA: []interface{}{"rhd", "whd"},
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
			}) {
				return internalHeadPosition(find(q, /* $node/ancestor::node/node[@rel=("rhd","whd")
					   and @index = $node/../node[@rel=("obj1","vc","se","me")]/@index] */&xPath{
						arg1: &dSort{
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
												DATA: []interface{}{"rhd", "whd"},
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
																DATA: []interface{}{"obj1", "vc", "se", "me"},
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
					}), q, tr)
			}
			if pobj1 := find(q /* $node/../node[@rel="pobj1"] */, &xPath{
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
									ARG:  collect__attributes__rel,
									arg1: &dNode{},
								},
								arg2: &dElem{
									DATA: []interface{}{"pobj1"},
									arg1: &dCollect{
										ARG:  collect__attributes__rel,
										arg1: &dNode{},
									},
								},
							},
						},
					},
				},
			}); len(pobj1) > 0 {
				return internalHeadPosition(if1(pobj1), q, tr)
			}
			// in de eerste rond --> typo in LassySmall/Wiki , binnen en [advp later buiten ]
			return externalHeadPosition(node.axParent, q, tr)
		} else {
			return externalHeadPosition(node.axParent, q, tr)
		}
	}

	aux := auxiliary1(node, q)

	if node.Rel == "hd" && (aux == "aux" || aux == "aux:pass") {
		// aux aux:pass cop
		if vc_predc := find(q /* $node/../node[@rel=("vc","predc")] */, &xPath{
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
								ARG:  collect__attributes__rel,
								arg1: &dNode{},
							},
							arg2: &dElem{
								DATA: []interface{}{"vc", "predc"},
								arg1: &dCollect{
									ARG:  collect__attributes__rel,
									arg1: &dNode{},
								},
							},
						},
					},
				},
			},
		}); len(vc_predc) > 0 {
			if test(q /* $vc_predc[@pt or (@cat and node[@pt or @cat])] */, &xPath{
				arg1: &dSort{
					arg1: &dFilter{
						arg1: &dVariable{
							VAR: vc_predc,
						},
						arg2: &dSort{
							arg1: &dOr{
								arg1: &dCollect{
									ARG:  collect__attributes__pt,
									arg1: &dNode{},
								},
								arg2: &dSort{
									arg1: &dAnd{
										arg1: &dCollect{
											ARG:  collect__attributes__cat,
											arg1: &dNode{},
										},
										arg2: &dCollect{
											ARG:  collect__child__node,
											arg1: &dNode{},
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
								},
							},
						},
					},
				},
			}) {
				// skip vc with just empty nodes
				return internalHeadPositionWithGapping(if1(vc_predc), q, tr)
			}
		}
		// if ($node/../node[@rel="predc"]/@index = $node/../../node[@rel="whd"]/@index)
		//     then local:internal_head_position($node/../../node[@rel="whd"])
		return externalHeadPosition(node.axParent, q, tr) // gapping, but does it ever occur with aux?? with cop: hij was en blijft nog steeds een omstreden figuur
	}

	if node.Rel == "hd" && aux == "cop" {
		predc := find(q /* $node/../node[@rel="predc"] */, &xPath{
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
		})
		if len(predc) > 0 && test(q /* $predc[@pt or @cat] */, &xPath{
			arg1: &dSort{
				arg1: &dFilter{
					arg1: &dVariable{
						VAR: predc,
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
			},
		}) {
			return internalHeadPositionWithGapping(if1(predc), q, tr)
		}
		if test(q /* $node/../node[@rel="predc"]/@index = $node/ancestor::node/node[@rel=("rhd","whd")]/@index */, &xPath{
			arg1: &dSort{
				arg1: &dEqual{
					ARG: equal__is,
					arg1: &dCollect{
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
					arg2: &dCollect{
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
								arg1: &dEqual{
									ARG: equal__is,
									arg1: &dCollect{
										ARG:  collect__attributes__rel,
										arg1: &dNode{},
									},
									arg2: &dElem{
										DATA: []interface{}{"rhd", "whd"},
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
		}) {
			return internalHeadPosition(
				find(q /* $node/ancestor::node/node[@rel=("rhd","whd") and @index = $node/../node[@rel="predc"]/@index] */, &xPath{
					arg1: &dSort{
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
											DATA: []interface{}{"rhd", "whd"},
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
						},
					},
				}),
				q,
				tr)
		}
		return externalHeadPosition(node.axParent, q, tr) // gapping, but could it??
	}

	if node.Rel == "hd" || node.Rel == "nucl" || node.Rel == "body" {
		if n := find(q /* $node/../node[@rel="hd" and @begin < $node/@begin] */, &xPath{
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
							arg2: &dCmp{
								ARG: cmp__lt,
								arg1: &dCollect{
									ARG:  collect__attributes__begin,
									arg1: &dNode{},
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
		}); len(n) > 0 {
			return internalHeadPosition(list(n), q, tr) // dan moet je moet
		}
		return externalHeadPosition(node.axParent, q, tr)
	}

	if node.Rel == "predc" {
		if test(q /* $node[../node[@rel=("obj1","se","vc")] and ../node[@rel="hd" and (@pt or @cat)]] */, &xPath{
			arg1: &dSort{
				arg1: &dFilter{
					arg1: &dVariable{
						VAR: node,
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
											DATA: []interface{}{"obj1", "se", "vc"},
											arg1: &dCollect{
												ARG:  collect__attributes__rel,
												arg1: &dNode{},
											},
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
								},
							},
						},
					},
				},
			},
		}) {
			if test(q /* $node/../node[@rel="hd" and @ud:pos="ADP"] */, &xPath{
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
										DATA: []interface{}{"ADP"},
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
			}) {
				return externalHeadPosition(node.axParent, q, tr) // met als presentator Bruno W , met als gevolg [vc dat ...]
			}
			return internalHeadPosition(find(q /* $node/../node[@rel="hd"] */, &xPath{
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
						},
					},
				},
			}), q, tr)
		}
		if test(q /* $node/parent::node[@cat=("np","ap") and node[@rel="hd" and (@pt or @cat) and not(@ud:pos="AUX") ]  ] */, &xPath{
			arg1: &dSort{
				arg1: &dCollect{
					ARG: collect__parent__node,
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
									DATA: []interface{}{"np", "ap"},
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
								},
							},
						},
					},
				},
			},
		}) {
			//reduced relatives , make sure head is not empty (ellipsis)
			// also : ap with predc: actief als schrijver
			return internalHeadPosition(find(q /* $node/../node[@rel="hd"] */, &xPath{
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
						},
					},
				},
			}), q, tr)
		}
		if test(q /* $node/../node[@rel="hd" and (@pt or @cat) and not(@ud:pos=("AUX","ADP"))] */, &xPath{
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
						arg1: &dAnd{
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
										arg1: &dEqual{
											ARG: equal__is,
											arg1: &dCollect{
												ARG:  collect__attributes__ud_3apos,
												arg1: &dNode{},
											},
											arg2: &dElem{
												DATA: []interface{}{"AUX", "ADP"},
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
		}) { // [met als titel] -- obj1/vc missing
			return internalHeadPosition(find(q /* $node/../node[@rel="hd"] */, &xPath{
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
						},
					},
				},
			}), q, tr)
		}
		return externalHeadPosition(node.axParent, q, tr) // covers gapping as well?
	}

	if test(q /* $node[@rel=("obj1","se","me") and (../@cat="pp" or ../node[@ud:pos="ADP" and @rel="hd"])] */, &xPath{
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
								DATA: []interface{}{"obj1", "se", "me"},
								arg1: &dCollect{
									ARG:  collect__attributes__rel,
									arg1: &dNode{},
								},
							},
						},
						arg2: &dSort{
							arg1: &dOr{
								arg1: &dEqual{
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
													ARG:  collect__attributes__ud_3apos,
													arg1: &dNode{},
												},
												arg2: &dElem{
													DATA: []interface{}{"ADP"},
													arg1: &dCollect{
														ARG:  collect__attributes__ud_3apos,
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
													DATA: []interface{}{"hd"},
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
	}) {
		if predc := find(q /* $node/../node[@rel="predc"] */, &xPath{
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
		}); len(predc) > 0 {
			return internalHeadPosition(predc, q, tr)
		}
		return externalHeadPosition(node.axParent, q, tr)
	}

	if test(q /* $node[@rel="pobj1" and (../@cat="pp" or ../node[@ud:pos="ADP" and @rel="hd"])] */, &xPath{
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
								DATA: []interface{}{"pobj1"},
								arg1: &dCollect{
									ARG:  collect__attributes__rel,
									arg1: &dNode{},
								},
							},
						},
						arg2: &dSort{
							arg1: &dOr{
								arg1: &dEqual{
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
													ARG:  collect__attributes__ud_3apos,
													arg1: &dNode{},
												},
												arg2: &dElem{
													DATA: []interface{}{"ADP"},
													arg1: &dCollect{
														ARG:  collect__attributes__ud_3apos,
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
													DATA: []interface{}{"hd"},
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
	}) {
		if vc := find(q /* $node/../node[@rel="vc"] */, &xPath{
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
			},
		}); len(vc) > 0 {
			return internalHeadPosition(vc, q, tr)
		}
		return externalHeadPosition(node.axParent, q, tr)
	}

	if test(q /* $node[@rel="mod" and not(../node[@rel=("obj1","pobj1","se","me")]) and (../@cat="pp" or ../node[@rel="hd" and @ud:pos="ADP"])] */, &xPath{
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
														DATA: []interface{}{"obj1", "pobj1", "se", "me"},
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
						arg2: &dSort{
							arg1: &dOr{
								arg1: &dEqual{
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
													ARG:  collect__attributes__ud_3apos,
													arg1: &dNode{},
												},
												arg2: &dElem{
													DATA: []interface{}{"ADP"},
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
			},
		},
	}) { // mede op grond hiervan
		// daarom dus
		if test(q /* $node/../node[@rel=("hd","su","obj1","vc") and (@pt or @cat)] */, &xPath{
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
						arg1: &dAnd{
							arg1: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG:  collect__attributes__rel,
									arg1: &dNode{},
								},
								arg2: &dElem{
									DATA: []interface{}{"hd", "su", "obj1", "vc"},
									arg1: &dCollect{
										ARG:  collect__attributes__rel,
										arg1: &dNode{},
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
					},
				},
			},
		}) {
			return internalHeadPositionWithGapping(node.axParent, q, tr)
		}
		return externalHeadPosition(node.axParent, q, tr) // gapping
	}

	if test(q /* $node[@rel=("cnj","dp","mwp")] */, &xPath{
		arg1: &dSort{
			arg1: &dFilter{
				arg1: &dVariable{
					VAR: node,
				},
				arg2: &dSort{
					arg1: &dEqual{
						ARG: equal__is,
						arg1: &dCollect{
							ARG:  collect__attributes__rel,
							arg1: &dNode{},
						},
						arg2: &dElem{
							DATA: []interface{}{"cnj", "dp", "mwp"},
							arg1: &dCollect{
								ARG:  collect__attributes__rel,
								arg1: &dNode{},
							},
						},
					},
				},
			},
		},
	}) {
		if node == nLeft(find(q /* $node/../node[@rel=("cnj","dp","mwp")] */, &xPath{
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
								ARG:  collect__attributes__rel,
								arg1: &dNode{},
							},
							arg2: &dElem{
								DATA: []interface{}{"cnj", "dp", "mwp"},
								arg1: &dCollect{
									ARG:  collect__attributes__rel,
									arg1: &dNode{},
								},
							},
						},
					},
				},
			},
		})) {
			return externalHeadPosition(node.axParent, q, tr)
		}
		if node.Rel == "cnj" {
			return headPositionOfConjunction(node, q, tr)
		}
		return internalHeadPositionWithGapping(node.axParent, q, tr)
	}

	if test(q /* $node[@rel="cmp" and ../node[@rel="body"]] */, &xPath{
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
								DATA: []interface{}{"cmp"},
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
								arg1: &dEqual{
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
			},
		},
	}) {
		return internalHeadPositionWithGapping(find(q /* $node/../node[@rel="body"][1] */, &xPath{
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
						arg1: &dPredicate{
							arg1: &dEqual{
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
						arg2: &dFunction{
							ARG: function__first__0__args,
						},
					},
				},
			},
		}), q, tr)
	}

	if node.Rel == "--" && node.Cat != "" {
		if node.Cat == "mwu" {
			if test(q /* $node/../node[@cat and not(@cat="mwu")] */, &xPath{
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
							arg1: &dAnd{
								arg1: &dCollect{
									ARG:  collect__attributes__cat,
									arg1: &dNode{},
								},
								arg2: &dFunction{
									ARG: function__not__1__args,
									arg1: &dArg{
										arg1: &dSort{
											arg1: &dEqual{
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
			}) { // fix for multiword punctuation in Alpino output
				return internalHeadPosition(find(q /* $node/../node[@cat and not(@cat="mwu")][1] */, &xPath{
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
								arg1: &dPredicate{
									arg1: &dAnd{
										arg1: &dCollect{
											ARG:  collect__attributes__cat,
											arg1: &dNode{},
										},
										arg2: &dFunction{
											ARG: function__not__1__args,
											arg1: &dArg{
												arg1: &dSort{
													arg1: &dEqual{
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
								arg2: &dFunction{
									ARG: function__first__0__args,
								},
							},
						},
					},
				}), q, tr)
			}
			return externalHeadPosition(node.axParent, q, tr)
		}
		return externalHeadPosition(node.axParent, q, tr)
	}

	if node.Rel == "--" && node.udPos != "" {
		if test(q, /* $node[@ud:pos = ("PUNCT","SYM","X","CONJ","NOUN","PROPN","NUM","ADP","ADV","DET","PRON")
			   and ../node[@rel="--" and
			               not(@ud:pos=("PUNCT","SYM","X","CONJ","NOUN","PROPN","NUM","ADP","ADV","DET","PRON")) ]
			  ] */&xPath{
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
										ARG:  collect__attributes__ud_3apos,
										arg1: &dNode{},
									},
									arg2: &dElem{
										DATA: []interface{}{"PUNCT", "SYM", "X", "CONJ", "NOUN", "PROPN", "NUM", "ADP", "ADV", "DET", "PRON"},
										arg1: &dCollect{
											ARG:  collect__attributes__ud_3apos,
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
													DATA: []interface{}{"--"},
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
																DATA: []interface{}{"PUNCT", "SYM", "X", "CONJ", "NOUN", "PROPN", "NUM", "ADP", "ADV", "DET", "PRON"},
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
						},
					},
				},
			}) {
			return internalHeadPositionWithGapping(
				find(q /* $node/../node[@rel="--" and not(@ud:pos=("PUNCT","SYM","X","CONJ","NOUN","ADP","ADV","DET","PROPN","NUM","PRON"))][1] */, &xPath{
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
								arg1: &dPredicate{
									arg1: &dAnd{
										arg1: &dEqual{
											ARG: equal__is,
											arg1: &dCollect{
												ARG:  collect__attributes__rel,
												arg1: &dNode{},
											},
											arg2: &dElem{
												DATA: []interface{}{"--"},
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
															DATA: []interface{}{"PUNCT", "SYM", "X", "CONJ", "NOUN", "ADP", "ADV", "DET", "PROPN", "NUM", "PRON"},
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
								arg2: &dFunction{
									ARG: function__first__0__args,
								},
							},
						},
					},
				}),
				q,
				tr)
		}
		if n := find(q /* $node/../node[@cat][1] */, &xPath{
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
						arg1: &dPredicate{
							arg1: &dCollect{
								ARG:  collect__attributes__cat,
								arg1: &dNode{},
							},
						},
						arg2: &dFunction{
							ARG: function__first__0__args,
						},
					},
				},
			},
		}); len(n) > 0 {
			return internalHeadPosition(n, q, tr)
		}
		if test(q /* $node[@ud:pos="PUNCT" and count(../node) > 1] */, &xPath{
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
									ARG:  collect__attributes__ud_3apos,
									arg1: &dNode{},
								},
								arg2: &dElem{
									DATA: []interface{}{"PUNCT"},
									arg1: &dCollect{
										ARG:  collect__attributes__ud_3apos,
										arg1: &dNode{},
									},
								},
							},
							arg2: &dCmp{
								ARG: cmp__gt,
								arg1: &dFunction{
									ARG: function__count__1__args,
									arg1: &dArg{
										arg1: &dCollect{
											ARG: collect__child__node,
											arg1: &dCollect{
												ARG:  collect__parent__type__node,
												arg1: &dNode{},
											},
										},
									},
								},
								arg2: &dElem{
									DATA: []interface{}{1},
									arg1: &dFunction{
										ARG: function__count__1__args,
										arg1: &dArg{
											arg1: &dCollect{
												ARG: collect__child__node,
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
			if n := find(q /* $node/../node[not(@ud:pos="PUNCT")][1] */, &xPath{
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
							arg1: &dPredicate{
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
													DATA: []interface{}{"PUNCT"},
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
							arg2: &dFunction{
								ARG: function__first__0__args,
							},
						},
					},
				},
			}); len(n) > 0 {
				return internalHeadPosition(n, q, tr)
			}
			if node == nLeft(find(q /* $node/../node[@rel="--" and (@cat or @pt)] */, &xPath{
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
							arg1: &dAnd{
								arg1: &dEqual{
									ARG: equal__is,
									arg1: &dCollect{
										ARG:  collect__attributes__rel,
										arg1: &dNode{},
									},
									arg2: &dElem{
										DATA: []interface{}{"--"},
										arg1: &dCollect{
											ARG:  collect__attributes__rel,
											arg1: &dNode{},
										},
									},
								},
								arg2: &dSort{
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
			})) {
				return externalHeadPosition(node.axParent, q, tr)
			}
			return 1000 // ie end of first punct token
		}
		if node.parent.Begin >= 0 {
			return externalHeadPosition(node.axParent, q, tr)
		}
		panic(tracer("No head found", tr, q))
	}

	if node.Rel == "dlink" || node.Rel == "sat" || node.Rel == "tag" {
		if n := find(q /* $node/../node[@rel="nucl"] */, &xPath{
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
								ARG:  collect__attributes__rel,
								arg1: &dNode{},
							},
							arg2: &dElem{
								DATA: []interface{}{"nucl"},
								arg1: &dCollect{
									ARG:  collect__attributes__rel,
									arg1: &dNode{},
								},
							},
						},
					},
				},
			},
		}); len(n) > 0 {
			return internalHeadPositionWithGapping(n, q, tr)
		}
		panic(tracer("No external head", tr, q))
	}

	if node.Rel == "vc" {
		if test(q, /* $node/../node[@rel="hd" and
			         ( @ud:pos="AUX" or
			           $node/ancestor::node[@rel="top"]//node[@ud:pos="AUX"]/@index = @index
			         )
			     ]
			and not($node/../node[@rel="predc"]) */&xPath{
				arg1: &dSort{
					arg1: &dAnd{
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
											DATA: []interface{}{"hd"},
											arg1: &dCollect{
												ARG:  collect__attributes__rel,
												arg1: &dNode{},
											},
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
													DATA: []interface{}{"AUX"},
													arg1: &dCollect{
														ARG:  collect__attributes__ud_3apos,
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
															ARG: collect__descendant__or__self__type__node,
															arg1: &dCollect{
																ARG: collect__ancestors__node,
																arg1: &dVariable{
																	VAR: node,
																},
																arg2: &dPredicate{
																	arg1: &dEqual{
																		ARG: equal__is,
																		arg1: &dCollect{
																			ARG:  collect__attributes__rel,
																			arg1: &dNode{},
																		},
																		arg2: &dElem{
																			DATA: []interface{}{"top"},
																			arg1: &dCollect{
																				ARG:  collect__attributes__rel,
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
												arg2: &dCollect{
													ARG:  collect__attributes__index,
													arg1: &dNode{},
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
				},
			}) {
			return externalHeadPosition(node.axParent, q, tr)
		}
		if test(q /* $node/../@cat="pp" */, &xPath{
			arg1: &dSort{
				arg1: &dEqual{
					ARG: equal__is,
					arg1: &dCollect{
						ARG: collect__attributes__cat,
						arg1: &dCollect{
							ARG: collect__parent__type__node,
							arg1: &dVariable{
								VAR: node,
							},
						},
					},
					arg2: &dElem{
						DATA: []interface{}{"pp"},
						arg1: &dCollect{
							ARG: collect__attributes__cat,
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
		}) { // eraan dat
			return externalHeadPosition(node.axParent, q, tr)
		}
		if test(q /* $node/../node[@rel=("hd","su") and (@pt or @cat)] */, &xPath{
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
						arg1: &dAnd{
							arg1: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG:  collect__attributes__rel,
									arg1: &dNode{},
								},
								arg2: &dElem{
									DATA: []interface{}{"hd", "su"},
									arg1: &dCollect{
										ARG:  collect__attributes__rel,
										arg1: &dNode{},
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
					},
				},
			},
		}) {
			return internalHeadPositionWithGapping(node.axParent, q, tr)
		}
		return externalHeadPosition(node.axParent, q, tr)
	}

	if node.Rel == "whd" || node.Rel == "rhd" {
		if node.Index > 0 {
			return externalHeadPosition(find(q /* ($node/../node[@rel="body"]//node[@index = $node/@index ])[1] */, &xPath{
				arg1: &dSort{
					arg1: &dFilter{
						arg1: &dSort{
							arg1: &dCollect{
								ARG: collect__child__node,
								arg1: &dCollect{
									ARG: collect__descendant__or__self__type__node,
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
								arg2: &dPredicate{
									arg1: &dEqual{
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
						arg2: &dSort{
							arg1: &dFunction{
								ARG: function__first__0__args,
							},
						},
					},
				},
			}), q, tr)
		}
		return internalHeadPosition(find(q /* $node/../node[@rel="body"] */, &xPath{
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
		}), q, tr)
	}

	/*
		we need to select the original node and not the result of
		following-cnj-sister, as that has no global context
		and global context is needed where the hd is an index node...
		unfortunately, nodes are completely identical in some
		elliptical cases, select last() as brute force solution
	*/
	if node.Rel == "crd" {
		tmp := followingCnjSister(node, q, tr)
		return internalHeadPositionWithGapping(ifZ(find(q, /* $node/../node[@rel="cnj" and
			   @begin=$tmp/@begin and
			   @end=$tmp/@end
			  ] */&xPath{
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
											ARG:  collect__attributes__begin,
											arg1: &dNode{},
										},
										arg2: &dCollect{
											ARG: collect__attributes__begin,
											arg1: &dVariable{
												VAR: tmp,
											},
										},
									},
								},
								arg2: &dEqual{
									ARG: equal__is,
									arg1: &dCollect{
										ARG:  collect__attributes__end,
										arg1: &dNode{},
									},
									arg2: &dCollect{
										ARG: collect__attributes__end,
										arg1: &dVariable{
											VAR: tmp,
										},
									},
								},
							},
						},
					},
				},
			})), q, tr)
	}

	if node.Rel == "su" {
		if test(q, /* $node[../node[@rel="vc"] and ../node[@rel="hd" and
			   ( @ud:pos="AUX" or $node/ancestor::node[@rel="top"]//node[@ud:pos="AUX"]/@index = @index ) ]] */&xPath{
				arg1: &dSort{
					arg1: &dFilter{
						arg1: &dVariable{
							VAR: node,
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
											arg2: &dSort{
												arg1: &dOr{
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
													arg2: &dEqual{
														ARG: equal__is,
														arg1: &dCollect{
															ARG: collect__attributes__index,
															arg1: &dCollect{
																ARG: collect__child__node,
																arg1: &dCollect{
																	ARG: collect__descendant__or__self__type__node,
																	arg1: &dCollect{
																		ARG: collect__ancestors__node,
																		arg1: &dVariable{
																			VAR: node,
																		},
																		arg2: &dPredicate{
																			arg1: &dEqual{
																				ARG: equal__is,
																				arg1: &dCollect{
																					ARG:  collect__attributes__rel,
																					arg1: &dNode{},
																				},
																				arg2: &dElem{
																					DATA: []interface{}{"top"},
																					arg1: &dCollect{
																						ARG:  collect__attributes__rel,
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
														arg2: &dCollect{
															ARG:  collect__attributes__index,
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
			tmp := internalHeadPositionWithGapping(find(q /* $node/../node[@rel="vc"] */, &xPath{
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
				},
			}), q, tr)
			if node.Begin < tmp && tmp <= node.End {
				return externalHeadPosition(node.axParent, q, tr)
			}
			return tmp
		}
		if test(q /* $node/../node[@rel="hd" and (@pt or @cat)] */, &xPath{
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
					},
				},
			},
		}) { // gapping
			return internalHeadPositionWithGapping(node.axParent, q, tr) // ud head could still be a predc
		}
		// only for 1 case where verb is missing -- die eigendom ... (no verb))
		if test(q /* $node[../node[@rel="predc" and (@pt or @cat)] and not(../node[@rel="hd" and (@pt or @cat)])] */, &xPath{
			arg1: &dSort{
				arg1: &dFilter{
					arg1: &dVariable{
						VAR: node,
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
												DATA: []interface{}{"predc"},
												arg1: &dCollect{
													ARG:  collect__attributes__rel,
													arg1: &dNode{},
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
			return internalHeadPositionWithGapping(find(q /* $node/../node[@rel="predc"] */, &xPath{
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
			}), q, tr)
		}
		return externalHeadPosition(node.axParent, q, tr) // this probably does no change anything, as we are still attaching to head of left conjunct
	}

	if node.Rel == "obj1" {
		if test(q /* $node/../node[@rel=("hd","su") and (@pt or @cat)] */, &xPath{
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
						arg1: &dAnd{
							arg1: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG:  collect__attributes__rel,
									arg1: &dNode{},
								},
								arg2: &dElem{
									DATA: []interface{}{"hd", "su"},
									arg1: &dCollect{
										ARG:  collect__attributes__rel,
										arg1: &dNode{},
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
					},
				},
			},
		}) { // gapping, as su but now su could be head as well
			return internalHeadPositionWithGapping(node.axParent, q, tr)
		}
		return externalHeadPosition(node.axParent, q, tr)
	}

	if node.Rel == "pc" {
		if test(q /* $node/../node[@rel=("hd","su","obj1") and (@pt or @cat)] */, &xPath{
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
						arg1: &dAnd{
							arg1: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG:  collect__attributes__rel,
									arg1: &dNode{},
								},
								arg2: &dElem{
									DATA: []interface{}{"hd", "su", "obj1"},
									arg1: &dCollect{
										ARG:  collect__attributes__rel,
										arg1: &dNode{},
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
					},
				},
			},
		}) { // gapping, as su but now su could be head as well
			return internalHeadPositionWithGapping(node.axParent, q, tr)
		}
		return externalHeadPosition(node.axParent, q, tr)
	}

	if node.Rel == "mod" {
		if test(q /* $node/../node[@rel=("hd","su","obj1","pc","predc","body") and (@pt or @cat)] */, &xPath{
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
						arg1: &dAnd{
							arg1: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG:  collect__attributes__rel,
									arg1: &dNode{},
								},
								arg2: &dElem{
									DATA: []interface{}{"hd", "su", "obj1", "pc", "predc", "body"},
									arg1: &dCollect{
										ARG:  collect__attributes__rel,
										arg1: &dNode{},
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
					},
				},
			},
		}) { // gapping, as su but now su or obj1  could be head as well
			return internalHeadPositionWithGapping(node.axParent, q, tr)
		}
		if n := find(q /* $node/../node[@rel="mod" and (@cat or @pt)] */, &xPath{
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
							arg2: &dSort{
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
		}); len(n) > 0 {
			if node == nLeft(n) { // gapping with multiple mods
				return externalHeadPosition(node.axParent, q, tr)
			}
			return internalHeadPositionWithGapping(node.axParent, q, tr)
		}
		if test(q /* $node/../../node[@rel="su" and (@pt or @cat)] */, &xPath{
			arg1: &dSort{
				arg1: &dCollect{
					ARG: collect__child__node,
					arg1: &dCollect{
						ARG: collect__parent__type__node,
						arg1: &dCollect{
							ARG: collect__parent__type__node,
							arg1: &dVariable{
								VAR: node,
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
									DATA: []interface{}{"su"},
									arg1: &dCollect{
										ARG:  collect__attributes__rel,
										arg1: &dNode{},
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
					},
				},
			},
		}) { // an mod in an otherwise empty tree (after fixing heads in conj)
			return internalHeadPosition(find(q /* $node/../../node[@rel="su"] */, &xPath{
				arg1: &dSort{
					arg1: &dCollect{
						ARG: collect__child__node,
						arg1: &dCollect{
							ARG: collect__parent__type__node,
							arg1: &dCollect{
								ARG: collect__parent__type__node,
								arg1: &dVariable{
									VAR: node,
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
			}), q, tr)
		}
		return externalHeadPosition(node.axParent, q, tr) /* an empty mod in an otherwise empty tree
		   -- mod is co-indexed with rhd, rest is elided,
		   LassySmall4/wiki-7064/wiki-7064.p.28.s.3.xml */
	}

	if node.Rel == "app" || node.Rel == "det" {
		if test(q /* $node/../node[@rel=("hd","mod") and (@pt or @cat)] */, &xPath{
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
						arg1: &dAnd{
							arg1: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG:  collect__attributes__rel,
									arg1: &dNode{},
								},
								arg2: &dElem{
									DATA: []interface{}{"hd", "mod"},
									arg1: &dCollect{
										ARG:  collect__attributes__rel,
										arg1: &dNode{},
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
					},
				},
			},
		}) { // gapping with an app (or a det)!
			return internalHeadPositionWithGapping(node.axParent, q, tr)
		}
		return externalHeadPosition(node.axParent, q, tr)
	}

	if node.Rel == "top" {
		return 0
	}

	if node.Rel != "hd" { // TODO: klopt dit?
		return internalHeadPositionWithGapping(node.axParent, q, tr)
	}

	panic(tracer("No external head", tr, q))
}

// recursive
func internalHeadPosition(nodes []interface{}, q *context, tr []trace) int {
	depthCheck(q, "internalHeadPosition")

	if n := len(nodes); n != 1 {
		tr = append(tr, trace{s: "internalHeadPosition"})
		if n == 0 {
			panic(tracer("No internal head position found", tr, q))
		} else if n > 1 {
			panic(tracer("More than one internal head position found", tr, q))
		}
	}
	node := nodes[0]

	if test(q /* $node[@cat="pp"] */, &xPath{
		arg1: &dSort{
			arg1: &dFilter{
				arg1: &dVariable{
					VAR: node,
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
	}) {
		// if ($node/node[@rel="hd" and @pt=("bw","n")] )  ( n --> TEMPORARY HACK to fix error where NP is erroneously tagged as PP )
		// then $node/node[@rel="hd"]/@end
		if n := find(q /* $node/node[@rel=("obj1","pobj1","se")][1] */, &xPath{
			arg1: &dSort{
				arg1: &dCollect{
					ARG: collect__child__node,
					arg1: &dVariable{
						VAR: node,
					},
					arg2: &dPredicate{
						arg1: &dPredicate{
							arg1: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG:  collect__attributes__rel,
									arg1: &dNode{},
								},
								arg2: &dElem{
									DATA: []interface{}{"obj1", "pobj1", "se"},
									arg1: &dCollect{
										ARG:  collect__attributes__rel,
										arg1: &dNode{},
									},
								},
							},
						},
						arg2: &dFunction{
							ARG: function__first__0__args,
						},
					},
				},
			},
		}); len(n) > 0 {
			return internalHeadPosition(n, q, tr)
		}
		if n := find(q /* $node/node[@rel="hd"] */, &xPath{
			arg1: &dSort{
				arg1: &dCollect{
					ARG: collect__child__node,
					arg1: &dVariable{
						VAR: node,
					},
					arg2: &dPredicate{
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
					},
				},
			},
		}); len(n) > 0 {
			// if ($node/@cat="mwu")  ( mede [op grond hiervan] )
			//     local:internal_head_position($node/node[@rel="hd"] )
			return internalHeadPosition(n, q, tr)
		}
		return internalHeadPosition(find(q /* $node/node[1] */, &xPath{
			arg1: &dSort{
				arg1: &dCollect{
					ARG: collect__child__node,
					arg1: &dVariable{
						VAR: node,
					},
					arg2: &dPredicate{
						arg1: &dFunction{
							ARG: function__first__0__args,
						},
					},
				},
			},
		}), q, tr)
	}

	if test(q /* $node[@cat="mwu"] */, &xPath{
		arg1: &dSort{
			arg1: &dFilter{
				arg1: &dVariable{
					VAR: node,
				},
				arg2: &dSort{
					arg1: &dEqual{
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
	}) {
		return find(q /* $node/node[@rel="mwp" and not(../node/@begin < @begin)]/@end */, &xPath{
			arg1: &dSort{
				arg1: &dCollect{
					ARG: collect__attributes__end,
					arg1: &dCollect{
						ARG: collect__child__node,
						arg1: &dVariable{
							VAR: node,
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
										DATA: []interface{}{"mwp"},
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
											arg1: &dCmp{
												ARG: cmp__lt,
												arg1: &dCollect{
													ARG: collect__attributes__begin,
													arg1: &dCollect{
														ARG: collect__child__node,
														arg1: &dCollect{
															ARG:  collect__parent__type__node,
															arg1: &dNode{},
														},
													},
												},
												arg2: &dCollect{
													ARG:  collect__attributes__begin,
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
		})[0].(int)
	}

	if test(q /* $node[@cat="conj"] */, &xPath{
		arg1: &dSort{
			arg1: &dFilter{
				arg1: &dVariable{
					VAR: node,
				},
				arg2: &dSort{
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
	}) {
		return internalHeadPosition(ifLeft(find(q /* $node/node[@rel="cnj"] */, &xPath{
			arg1: &dSort{
				arg1: &dCollect{
					ARG: collect__child__node,
					arg1: &dVariable{
						VAR: node,
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
			},
		})), q, tr)
	}

	if predc := find(q /* $node/node[@rel="predc"] */, &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dVariable{
					VAR: node,
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
	}); len(predc) > 0 {
		if test(q /* $node/node[@rel="hd" and @ud:pos="AUX"] */, &xPath{
			arg1: &dSort{
				arg1: &dCollect{
					ARG: collect__child__node,
					arg1: &dVariable{
						VAR: node,
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
		}) {
			return internalHeadPosition(predc, q, tr)
		}
		hd := find(q /* $node/node[@rel="hd"] */, &xPath{
			arg1: &dSort{
				arg1: &dCollect{
					ARG: collect__child__node,
					arg1: &dVariable{
						VAR: node,
					},
					arg2: &dPredicate{
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
					},
				},
			},
		})
		if len(hd) == 0 { // cases where copula is missing by accident (ungrammatical, not gapping)
			return internalHeadPosition(predc, q, tr)
		}
		return internalHeadPosition(hd, q, tr)
	}

	if test(q, /* $node[node[@rel="vc"] and
		   node[@rel="hd" and
		        ( @ud:pos="AUX" or
		          $node/ancestor::node[@rel="top"]//node[@ud:pos="AUX"]/@index = @index
		         )
		       ]
		  ] */&xPath{
			arg1: &dSort{
				arg1: &dFilter{
					arg1: &dVariable{
						VAR: node,
					},
					arg2: &dSort{
						arg1: &dAnd{
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
											DATA: []interface{}{"vc"},
											arg1: &dCollect{
												ARG:  collect__attributes__rel,
												arg1: &dNode{},
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
												DATA: []interface{}{"hd"},
												arg1: &dCollect{
													ARG:  collect__attributes__rel,
													arg1: &dNode{},
												},
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
														DATA: []interface{}{"AUX"},
														arg1: &dCollect{
															ARG:  collect__attributes__ud_3apos,
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
																ARG: collect__descendant__or__self__type__node,
																arg1: &dCollect{
																	ARG: collect__ancestors__node,
																	arg1: &dVariable{
																		VAR: node,
																	},
																	arg2: &dPredicate{
																		arg1: &dEqual{
																			ARG: equal__is,
																			arg1: &dCollect{
																				ARG:  collect__attributes__rel,
																				arg1: &dNode{},
																			},
																			arg2: &dElem{
																				DATA: []interface{}{"top"},
																				arg1: &dCollect{
																					ARG:  collect__attributes__rel,
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
													arg2: &dCollect{
														ARG:  collect__attributes__index,
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
		return internalHeadPosition(find(q /* $node/node[@rel="vc"] */, &xPath{
			arg1: &dSort{
				arg1: &dCollect{
					ARG: collect__child__node,
					arg1: &dVariable{
						VAR: node,
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
			},
		}), q, tr)
	}

	if n := find(q /* $node/node[@rel="hd"][1] */, &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dVariable{
					VAR: node,
				},
				arg2: &dPredicate{
					arg1: &dPredicate{
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
					},
					arg2: &dFunction{
						ARG: function__first__0__args,
					},
				},
			},
		},
	}); len(n) > 0 {
		return internalHeadPosition(n, q, tr)
	}

	if n := find(q /* $node/node[@rel="body"][1] */, &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dVariable{
					VAR: node,
				},
				arg2: &dPredicate{
					arg1: &dPredicate{
						arg1: &dEqual{
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
					arg2: &dFunction{
						ARG: function__first__0__args,
					},
				},
			},
		},
	}); len(n) > 0 {
		return internalHeadPosition(n, q, tr)
	}

	if n := find(q /* $node/node[@rel="dp"] */, &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dVariable{
					VAR: node,
				},
				arg2: &dPredicate{
					arg1: &dEqual{
						ARG: equal__is,
						arg1: &dCollect{
							ARG:  collect__attributes__rel,
							arg1: &dNode{},
						},
						arg2: &dElem{
							DATA: []interface{}{"dp"},
							arg1: &dCollect{
								ARG:  collect__attributes__rel,
								arg1: &dNode{},
							},
						},
					},
				},
			},
		},
	}); len(n) > 0 {
		return internalHeadPosition(ifLeft(n), q, tr)
		// sometimes co-indexing leads to du's starting at same position ...
	}

	if n := find(q /* $node/node[@rel="nucl"] */, &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dVariable{
					VAR: node,
				},
				arg2: &dPredicate{
					arg1: &dEqual{
						ARG: equal__is,
						arg1: &dCollect{
							ARG:  collect__attributes__rel,
							arg1: &dNode{},
						},
						arg2: &dElem{
							DATA: []interface{}{"nucl"},
							arg1: &dCollect{
								ARG:  collect__attributes__rel,
								arg1: &dNode{},
							},
						},
					},
				},
			},
		},
	}); len(n) > 0 {
		return internalHeadPosition(if1(n), q, tr)
	}

	if n := find(q /* $node/node[@cat="du"] */, &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dVariable{
					VAR: node,
				},
				arg2: &dPredicate{
					arg1: &dEqual{
						ARG: equal__is,
						arg1: &dCollect{
							ARG:  collect__attributes__cat,
							arg1: &dNode{},
						},
						arg2: &dElem{
							DATA: []interface{}{"du"},
							arg1: &dCollect{
								ARG:  collect__attributes__cat,
								arg1: &dNode{},
							},
						},
					},
				},
			},
		},
	}); len(n) > 0 { // is this neccesary at all? , only one referring to cat, and causes problems if applied before @rel=dp case...
		return internalHeadPosition(if1(n), q, tr)
	}

	if test(q /* $node[@word] */, &xPath{
		arg1: &dSort{
			arg1: &dFilter{
				arg1: &dVariable{
					VAR: node,
				},
				arg2: &dSort{
					arg1: &dCollect{
						ARG:  collect__attributes__word,
						arg1: &dNode{},
					},
				},
			},
		},
	}) {
		return i1(find(q /* $node/@end */, &xPath{
			arg1: &dSort{
				arg1: &dCollect{
					ARG: collect__attributes__end,
					arg1: &dVariable{
						VAR: node,
					},
				},
			},
		}))
	}

	/*
		distinguish empty nodes due to gapping/RNR from nonlocal dependencies
		use 'empty head' as string to signal precence of gapping/RNR
	*/
	if test(q /* $node[@index and not(@word or @cat)] */, &xPath{
		arg1: &dSort{
			arg1: &dFilter{
				arg1: &dVariable{
					VAR: node,
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
											ARG:  collect__attributes__word,
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
	}) {
		if test(q /* $node/ancestor::node[@rel="top"]//node[@rel=("whd","rhd") and @index = $node/@index and (@word or @cat)] */, &xPath{
			arg1: &dSort{
				arg1: &dCollect{
					ARG: collect__child__node,
					arg1: &dCollect{
						ARG: collect__descendant__or__self__type__node,
						arg1: &dCollect{
							ARG: collect__ancestors__node,
							arg1: &dVariable{
								VAR: node,
							},
							arg2: &dPredicate{
								arg1: &dEqual{
									ARG: equal__is,
									arg1: &dCollect{
										ARG:  collect__attributes__rel,
										arg1: &dNode{},
									},
									arg2: &dElem{
										DATA: []interface{}{"top"},
										arg1: &dCollect{
											ARG:  collect__attributes__rel,
											arg1: &dNode{},
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
										DATA: []interface{}{"whd", "rhd"},
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
											VAR: node,
										},
									},
								},
							},
							arg2: &dSort{
								arg1: &dOr{
									arg1: &dCollect{
										ARG:  collect__attributes__word,
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
		}) {
			return internalHeadPosition(
				find(q /* $node/ancestor::node[@rel="top"]//node[@index = $node/@index and (@word or @cat)] */, &xPath{
					arg1: &dSort{
						arg1: &dCollect{
							ARG: collect__child__node,
							arg1: &dCollect{
								ARG: collect__descendant__or__self__type__node,
								arg1: &dCollect{
									ARG: collect__ancestors__node,
									arg1: &dVariable{
										VAR: node,
									},
									arg2: &dPredicate{
										arg1: &dEqual{
											ARG: equal__is,
											arg1: &dCollect{
												ARG:  collect__attributes__rel,
												arg1: &dNode{},
											},
											arg2: &dElem{
												DATA: []interface{}{"top"},
												arg1: &dCollect{
													ARG:  collect__attributes__rel,
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
									arg2: &dSort{
										arg1: &dOr{
											arg1: &dCollect{
												ARG:  collect__attributes__word,
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
				}),
				q,
				tr)
		}
		return empty_head
	}

	panic(tracer("No internal head", tr, q))
}

func internalHeadPositionWithGapping(node []interface{}, q *context, tr []trace) int {
	if len(node) == 1 {
		tr = append(tr, trace{s: "internalHeadPositionWithGapping", node: node[0].(*nodeType)})
	} else {
		tr = append(tr, trace{s: fmt.Sprintf("internalHeadPositionWithGapping with %d nodes", len(node))})
	}
	if hdPos := internalHeadPosition(node, q, tr); hdPos == empty_head {
		return internalHeadPositionOfGappedConstituent(node, q, tr)
	} else {
		return hdPos
	}
}

func internalHeadPositionOfGappedConstituent(node []interface{}, q *context, tr []trace) int {
	depthCheck(q, "internalHeadPositionOfGappedConstituent")
	if len(node) == 1 {
	} else {
		tr = append(tr, trace{s: fmt.Sprintf("internalHeadPositionOfGappedConstituent with %d nodes", len(node))})
	}

	if test(q /* $node/node[@rel="hd" and (@pt or @cat) and not(@ud:pos=("AUX","ADP"))] */, &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dVariable{
					VAR: node,
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
									DATA: []interface{}{"hd"},
									arg1: &dCollect{
										ARG:  collect__attributes__rel,
										arg1: &dNode{},
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
									arg1: &dEqual{
										ARG: equal__is,
										arg1: &dCollect{
											ARG:  collect__attributes__ud_3apos,
											arg1: &dNode{},
										},
										arg2: &dElem{
											DATA: []interface{}{"AUX", "ADP"},
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
	}) {
		return internalHeadPositionWithGapping(find(q /* $node/node[@rel="hd"] */, &xPath{
			arg1: &dSort{
				arg1: &dCollect{
					ARG: collect__child__node,
					arg1: &dVariable{
						VAR: node,
					},
					arg2: &dPredicate{
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
					},
				},
			},
		}), q, tr) // aux, prepositions
	}

	if test(q /* $node[ node[@rel="hd" and @ud:pos="AUX"] and node[@rel=("vc","predc")] ] */, &xPath{
		arg1: &dSort{
			arg1: &dFilter{
				arg1: &dVariable{
					VAR: node,
				},
				arg2: &dSort{
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
						arg2: &dCollect{
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
										DATA: []interface{}{"vc", "predc"},
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
		return internalHeadPositionWithGapping(find(q /* $node/node[@rel=("vc","predc")] */, &xPath{
			arg1: &dSort{
				arg1: &dCollect{
					ARG: collect__child__node,
					arg1: &dVariable{
						VAR: node,
					},
					arg2: &dPredicate{
						arg1: &dEqual{
							ARG: equal__is,
							arg1: &dCollect{
								ARG:  collect__attributes__rel,
								arg1: &dNode{},
							},
							arg2: &dElem{
								DATA: []interface{}{"vc", "predc"},
								arg1: &dCollect{
									ARG:  collect__attributes__rel,
									arg1: &dNode{},
								},
							},
						},
					},
				},
			},
		}), q, tr)
	}

	if test(q /* $node/node[@rel="su" and (@pt or @cat)] */, &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dVariable{
					VAR: node,
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
								DATA: []interface{}{"su"},
								arg1: &dCollect{
									ARG:  collect__attributes__rel,
									arg1: &dNode{},
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
				},
			},
		},
	}) {
		return internalHeadPositionWithGapping(find(q /* $node/node[@rel="su"] */, &xPath{
			arg1: &dSort{
				arg1: &dCollect{
					ARG: collect__child__node,
					arg1: &dVariable{
						VAR: node,
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
		}), q, tr) // 44 van 87 in lassysmall
	}

	if test(q /* $node[@rel="vc" and ../node[@rel="su" and (@pt or @cat)]] */, &xPath{
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
								DATA: []interface{}{"vc"},
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
											DATA: []interface{}{"su"},
											arg1: &dCollect{
												ARG:  collect__attributes__rel,
												arg1: &dNode{},
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
							},
						},
					},
				},
			},
		},
	}) {
		// subject realized inside the vc = funny serialization
		return internalHeadPositionWithGapping(find(q /* $node/../node[@rel="su"] */, &xPath{
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
		}), q, tr)
	}

	if test(q /* $node/node[@rel="obj1" and (@pt or @cat)] */, &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dVariable{
					VAR: node,
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
								DATA: []interface{}{"obj1"},
								arg1: &dCollect{
									ARG:  collect__attributes__rel,
									arg1: &dNode{},
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
				},
			},
		},
	}) {
		return internalHeadPositionWithGapping(find(q /* $node/node[@rel="obj1"] */, &xPath{
			arg1: &dSort{
				arg1: &dCollect{
					ARG: collect__child__node,
					arg1: &dVariable{
						VAR: node,
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
		}), q, tr)
	}

	if test(q /* $node/node[@rel="predc" and (@pt or @cat)] */, &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dVariable{
					VAR: node,
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
								DATA: []interface{}{"predc"},
								arg1: &dCollect{
									ARG:  collect__attributes__rel,
									arg1: &dNode{},
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
				},
			},
		},
	}) {
		return internalHeadPositionWithGapping(find(q /* $node/node[@rel="predc"] */, &xPath{
			arg1: &dSort{
				arg1: &dCollect{
					ARG: collect__child__node,
					arg1: &dVariable{
						VAR: node,
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
		}), q, tr)
	}

	if test(q /* $node/node[@rel="vc" and (@pt or @cat)] */, &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dVariable{
					VAR: node,
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
	}) {
		return internalHeadPositionWithGapping(find(q /* $node/node[@rel="vc"][1] */, &xPath{
			arg1: &dSort{
				arg1: &dCollect{
					ARG: collect__child__node,
					arg1: &dVariable{
						VAR: node,
					},
					arg2: &dPredicate{
						arg1: &dPredicate{
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
						arg2: &dFunction{
							ARG: function__first__0__args,
						},
					},
				},
			},
		}), q, tr)
	}

	if test(q /* $node/node[@rel="pc" and (@pt or @cat)] */, &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dVariable{
					VAR: node,
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
								DATA: []interface{}{"pc"},
								arg1: &dCollect{
									ARG:  collect__attributes__rel,
									arg1: &dNode{},
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
				},
			},
		},
	}) {
		return internalHeadPositionWithGapping(find(q /* $node/node[@rel="pc"][1] */, &xPath{
			arg1: &dSort{
				arg1: &dCollect{
					ARG: collect__child__node,
					arg1: &dVariable{
						VAR: node,
					},
					arg2: &dPredicate{
						arg1: &dPredicate{
							arg1: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG:  collect__attributes__rel,
									arg1: &dNode{},
								},
								arg2: &dElem{
									DATA: []interface{}{"pc"},
									arg1: &dCollect{
										ARG:  collect__attributes__rel,
										arg1: &dNode{},
									},
								},
							},
						},
						arg2: &dFunction{
							ARG: function__first__0__args,
						},
					},
				},
			},
		}), q, tr)
	}

	if n := find(q /* $node/node[@rel="mod" and (@pt or @cat)] */, &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dVariable{
					VAR: node,
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
				},
			},
		},
	}); len(n) > 0 {
		return internalHeadPositionWithGapping(if1(n), q, tr)
	}

	if n := find(q /* $node/node[@rel="app" and (@pt or @cat)] */, &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dVariable{
					VAR: node,
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
								DATA: []interface{}{"app"},
								arg1: &dCollect{
									ARG:  collect__attributes__rel,
									arg1: &dNode{},
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
				},
			},
		},
	}); len(n) > 0 {
		return internalHeadPositionWithGapping(if1(n), q, tr)
	}

	if n := find(q /* $node/node[@rel="det" and (@pt or @cat)] */, &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dVariable{
					VAR: node,
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
								DATA: []interface{}{"det"},
								arg1: &dCollect{
									ARG:  collect__attributes__rel,
									arg1: &dNode{},
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
				},
			},
		},
	}); len(n) > 0 {
		return internalHeadPositionWithGapping(if1(n), q, tr)
	}

	if n := find(q /* $node/node[@rel="body" and (@pt or @cat)] */, &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dVariable{
					VAR: node,
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
								DATA: []interface{}{"body"},
								arg1: &dCollect{
									ARG:  collect__attributes__rel,
									arg1: &dNode{},
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
				},
			},
		},
	}); len(n) > 0 {
		return internalHeadPositionWithGapping(if1(n), q, tr)
	}

	if n := find(q /* $node/node[@rel="cnj" and (@pt or @cat)] */, &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dVariable{
					VAR: node,
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
				},
			},
		},
	}); len(n) > 0 {
		return internalHeadPositionWithGapping(if1(n), q, tr)
	}

	if n := find(q /* $node/node[@rel="dp" and (@pt or @cat)] */, &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dVariable{
					VAR: node,
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
								DATA: []interface{}{"dp"},
								arg1: &dCollect{
									ARG:  collect__attributes__rel,
									arg1: &dNode{},
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
				},
			},
		},
	}); len(n) > 0 {
		return internalHeadPositionWithGapping(if1(n), q, tr)
	}
	if n := find(q /* $node/node[@rel="hd" and @ud:pos="ADP"] */, &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dVariable{
					VAR: node,
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
								ARG:  collect__attributes__ud_3apos,
								arg1: &dNode{},
							},
							arg2: &dElem{
								DATA: []interface{}{"ADP"},
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
	}); len(n) > 0 { // in en rond Brussel, case not necessary in xquery code (run-time issue?)
		return internalHeadPositionWithGapping(if1(n), q, tr)
	}

	panic(tracer("No internal head in gapped constituent", tr, q))
}

/*
brute force method to be compliant with conj points to the left rule:
if interhdpos($node) < internalhdpos($node/..) then do something ad hoc
because even fixing misplaced heads fails in cases like
Het front der activisten vertoont dan wel een beeld van lusteloosheid , " aan de basis " is en wordt toch veel werk verzet .
*/
func headPositionOfConjunction(node *nodeType, q *context, tr []trace) int {

	tr = append(tr, trace{s: "headPositionOfConjunction", node: node})

	internal_head := internalHeadPositionWithGapping([]interface{}{node}, q, tr)
	leftmost_conj_daughter := nLeft(find(q /* $node/../node[@rel="cnj"] */, &xPath{
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
	}))
	leftmost_internal_head := internalHeadPositionWithGapping([]interface{}{leftmost_conj_daughter}, q, tr)

	if leftmost_internal_head < internal_head {
		return leftmost_internal_head
	}

	endpos_of_leftmost_conj_constituents := []int{}
	for _, e := range leftmost_conj_daughter.Node {
		if e.End < internal_head {
			endpos_of_leftmost_conj_constituents = append(endpos_of_leftmost_conj_constituents, e.End)
		}
	}
	if len(endpos_of_leftmost_conj_constituents) == 0 {
		return leftmost_conj_daughter.Node[0].End // this should not happen really -- give error msg?
	}
	sort.Ints(endpos_of_leftmost_conj_constituents)
	return endpos_of_leftmost_conj_constituents[len(endpos_of_leftmost_conj_constituents)-1]
}

func followingCnjSister(node *nodeType, q *context, tr []trace) []interface{} {
	tr = append(tr, trace{s: "followingCnjSister", node: node})
	/*
	   declare function local:following-cnj-sister($node as element(node)) as element(node)
	   { let $annotated-sisters :=
	         for $sister in $node/../node[@rel="cnj"]
	         return
	            <begin-node begin="{local:begin-position-of-first-word($sister)}">
	              {$sister}
	            </begin-node>

	     let $sorted-sisters :=
	         for $sister in $annotated-sisters
	         (: where $sister[number(@begin) > $node/number(@begin)] :)
	         order by $sister/number(@begin)
	         return $sister
	     return
	         if  ($sorted-sisters[number(@begin) > $node/number(@begin)] )
	         then ($sorted-sisters[number(@begin) > $node/number(@begin)]/node)[1]
	         else $sorted-sisters[1]/node

	   };
	*/

	// TODO: klopt dit ???

	sisters := []*nodeType{}
	for _, n := range node.parent.Node {
		if n.Rel == "cnj" /* && n.Begin > node.Begin */ {
			b := find(q /* $n/descendant-or-self::node[@word]/@begin */, &xPath{
				arg1: &dSort{
					arg1: &dCollect{
						ARG: collect__attributes__begin,
						arg1: &dCollect{
							ARG: collect__descendant__or__self__node,
							arg1: &dVariable{
								VAR: n,
							},
							arg2: &dPredicate{
								arg1: &dCollect{
									ARG:  collect__attributes__word,
									arg1: &dNode{},
								},
							},
						},
					},
				},
			})
			if len(b) == 0 {
				n.udFirstWordBegin = 1000000
			} else {
				sort.Slice(b, func(i, j int) bool {
					return b[i].(int) < b[j].(int)
				})
				n.udFirstWordBegin = b[0].(int)
			}
			sisters = append(sisters, n)
		}
	}
	sort.Slice(sisters, func(i, j int) bool {
		return sisters[i].udFirstWordBegin < sisters[j].udFirstWordBegin
	})
	for _, n := range sisters {
		if n.udFirstWordBegin > node.Begin {
			return []interface{}{n}
		}
	}
	if len(sisters) > 0 {
		return []interface{}{sisters[0]}
	}
	return []interface{}{}
}
