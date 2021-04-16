// GENERATED FILE. DO NOT EDIT.

package main

var spod2xpath = map[string]*xPath{
	"smain": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dEqual{
						ARG: equal__is,
						arg1: &dCollect{
							ARG:  collect__attributes__cat,
							arg1: &dNode{},
						},
						arg2: &dElem{
							DATA: []interface{}{"smain"},
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
	"whq": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dEqual{
						ARG: equal__is,
						arg1: &dCollect{
							ARG:  collect__attributes__cat,
							arg1: &dNode{},
						},
						arg2: &dElem{
							DATA: []interface{}{"whq"},
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
	"janee": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dAnd{
						arg1: &dAnd{
							arg1: &dAnd{
								arg1: &dEqual{
									ARG: equal__is,
									arg1: &dCollect{
										ARG:  collect__attributes__cat,
										arg1: &dNode{},
									},
									arg2: &dElem{
										DATA: []interface{}{"sv1"},
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
							arg2: &dFunction{
								ARG: function__not__1__args,
								arg1: &dArg{
									arg1: &dSort{
										arg1: &dCollect{
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
														arg2: &dCollect{
															ARG:  collect__attributes__stype,
															arg1: &dNode{},
														},
													},
													arg2: &dFunction{
														ARG: function__not__1__args,
														arg1: &dArg{
															arg1: &dSort{
																arg1: &dEqual{
																	ARG: equal__is,
																	arg1: &dCollect{
																		ARG:  collect__attributes__stype,
																		arg1: &dNode{},
																	},
																	arg2: &dElem{
																		DATA: []interface{}{"ynquestion"},
																		arg1: &dCollect{
																			ARG:  collect__attributes__stype,
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
						arg2: &dCmp{
							ARG: cmp__lt,
							arg1: &dCollect{
								ARG:  collect__attributes__end,
								arg1: &dNode{},
							},
							arg2: &dCollect{
								ARG: collect__attributes__end,
								arg1: &dCollect{
									ARG: collect__child__node,
									arg1: &dCollect{
										ARG:  collect__descendant__or__self__type__node,
										arg1: &dRoot{},
									},
									arg2: &dPredicate{
										arg1: &dEqual{
											ARG: equal__is,
											arg1: &dCollect{
												ARG:  collect__attributes__word,
												arg1: &dNode{},
											},
											arg2: &dElem{
												DATA: []interface{}{"?"},
												arg1: &dCollect{
													ARG:  collect__attributes__word,
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
	"imp": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dSort{
						arg1: &dAnd{
							arg1: &dAnd{
								arg1: &dAnd{
									arg1: &dEqual{
										ARG: equal__is,
										arg1: &dCollect{
											ARG:  collect__attributes__cat,
											arg1: &dNode{},
										},
										arg2: &dElem{
											DATA: []interface{}{"sv1"},
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
								arg2: &dFunction{
									ARG: function__not__1__args,
									arg1: &dArg{
										arg1: &dSort{
											arg1: &dCollect{
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
															arg2: &dCollect{
																ARG:  collect__attributes__stype,
																arg1: &dNode{},
															},
														},
														arg2: &dFunction{
															ARG: function__not__1__args,
															arg1: &dArg{
																arg1: &dSort{
																	arg1: &dEqual{
																		ARG: equal__is,
																		arg1: &dCollect{
																			ARG:  collect__attributes__stype,
																			arg1: &dNode{},
																		},
																		arg2: &dElem{
																			DATA: []interface{}{"imparative"},
																			arg1: &dCollect{
																				ARG:  collect__attributes__stype,
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
															arg1: &dOr{
																arg1: &dOr{
																	arg1: &dEqual{
																		ARG: equal__is,
																		arg1: &dCollect{
																			ARG:  collect__attributes__tense,
																			arg1: &dNode{},
																		},
																		arg2: &dElem{
																			DATA: []interface{}{"past"},
																			arg1: &dCollect{
																				ARG:  collect__attributes__tense,
																				arg1: &dNode{},
																			},
																		},
																	},
																	arg2: &dEqual{
																		ARG: equal__is,
																		arg1: &dCollect{
																			ARG:  collect__attributes__pvagr,
																			arg1: &dNode{},
																		},
																		arg2: &dElem{
																			DATA: []interface{}{"mv", "met-t"},
																			arg1: &dCollect{
																				ARG:  collect__attributes__pvagr,
																				arg1: &dNode{},
																			},
																		},
																	},
																},
																arg2: &dEqual{
																	ARG: equal__is,
																	arg1: &dCollect{
																		ARG:  collect__attributes__pvtijd,
																		arg1: &dNode{},
																	},
																	arg2: &dElem{
																		DATA: []interface{}{"verl"},
																		arg1: &dCollect{
																			ARG:  collect__attributes__pvtijd,
																			arg1: &dNode{},
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
																	DATA: []interface{}{"zijn", "kunnen", "willen", "moeten", "mogen", "zullen", "denken", "geloven", "vinden", "hebben"},
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
						},
					},
				},
			},
		},
	},
	"whsub": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
								DATA: []interface{}{"whsub"},
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
											DATA: []interface{}{"body"},
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
														DATA: []interface{}{"ssub"},
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
																			DATA: []interface{}{"ssub"},
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
													arg2: &dCollect{
														ARG:  collect__child__node,
														arg1: &dNode{},
														arg2: &dPredicate{
															arg1: &dEqual{
																ARG: equal__is,
																arg1: &dCollect{
																	ARG:  collect__attributes__cat,
																	arg1: &dNode{},
																},
																arg2: &dElem{
																	DATA: []interface{}{"ssub"},
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
	},
	"ssub": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
								DATA: []interface{}{"cp"},
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
											DATA: []interface{}{"body"},
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
											DATA: []interface{}{"ssub"},
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
	"ssubdat": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dAnd{
						arg1: &dAnd{
							arg1: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG:  collect__attributes__cat,
									arg1: &dNode{},
								},
								arg2: &dElem{
									DATA: []interface{}{"cp"},
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
												DATA: []interface{}{"cmp"},
												arg1: &dCollect{
													ARG:  collect__attributes__rel,
													arg1: &dNode{},
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
												DATA: []interface{}{"dat"},
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
											DATA: []interface{}{"body"},
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
											DATA: []interface{}{"ssub"},
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
	"ssubof": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dAnd{
						arg1: &dAnd{
							arg1: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG:  collect__attributes__cat,
									arg1: &dNode{},
								},
								arg2: &dElem{
									DATA: []interface{}{"cp"},
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
												DATA: []interface{}{"cmp"},
												arg1: &dCollect{
													ARG:  collect__attributes__rel,
													arg1: &dNode{},
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
												DATA: []interface{}{"of"},
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
											DATA: []interface{}{"body"},
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
											DATA: []interface{}{"ssub"},
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
	"ssubcmp": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dAnd{
						arg1: &dAnd{
							arg1: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG:  collect__attributes__cat,
									arg1: &dNode{},
								},
								arg2: &dElem{
									DATA: []interface{}{"cp"},
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
												DATA: []interface{}{"cmp"},
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
															ARG:  collect__attributes__lemma,
															arg1: &dNode{},
														},
														arg2: &dElem{
															DATA: []interface{}{"of", "dat"},
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
											DATA: []interface{}{"body"},
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
											DATA: []interface{}{"ssub"},
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
	"oti": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dEqual{
						ARG: equal__is,
						arg1: &dCollect{
							ARG:  collect__attributes__cat,
							arg1: &dNode{},
						},
						arg2: &dElem{
							DATA: []interface{}{"oti"},
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
	"otivc": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
								DATA: []interface{}{"oti"},
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
		},
	},
	"otimod": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
								DATA: []interface{}{"oti"},
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
								DATA: []interface{}{"mod"},
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
	"otiww": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
								DATA: []interface{}{"oti"},
								arg1: &dCollect{
									ARG:  collect__attributes__cat,
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
											DATA: []interface{}{"ww"},
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
	},
	"otin": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
								DATA: []interface{}{"oti"},
								arg1: &dCollect{
									ARG:  collect__attributes__cat,
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
											DATA: []interface{}{"n", "vnw"},
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
	},
	"otisu": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
								DATA: []interface{}{"oti"},
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
	"otipred": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
								DATA: []interface{}{"oti"},
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
	"otiobc": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
								DATA: []interface{}{"oti"},
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
								DATA: []interface{}{"obcomp"},
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
	"tite": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
								DATA: []interface{}{"ti"},
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
											ARG: collect__attributes__cat,
											arg1: &dCollect{
												ARG:  collect__parent__type__node,
												arg1: &dNode{},
											},
										},
										arg2: &dElem{
											DATA: []interface{}{"oti", "cp"},
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
				},
			},
		},
	},
	"ti": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
								DATA: []interface{}{"cp"},
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
											DATA: []interface{}{"body"},
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
											DATA: []interface{}{"ti"},
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
	"relssub": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
								DATA: []interface{}{"rel"},
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
											DATA: []interface{}{"body"},
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
														DATA: []interface{}{"ssub"},
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
																			DATA: []interface{}{"ssub"},
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
													arg2: &dCollect{
														ARG:  collect__child__node,
														arg1: &dNode{},
														arg2: &dPredicate{
															arg1: &dEqual{
																ARG: equal__is,
																arg1: &dCollect{
																	ARG:  collect__attributes__cat,
																	arg1: &dNode{},
																},
																arg2: &dElem{
																	DATA: []interface{}{"ssub"},
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
	},
	"whrel": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dAnd{
						arg1: &dAnd{
							arg1: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG:  collect__attributes__cat,
									arg1: &dNode{},
								},
								arg2: &dElem{
									DATA: []interface{}{"whrel"},
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
											DATA: []interface{}{"body"},
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
														DATA: []interface{}{"ssub"},
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
																			DATA: []interface{}{"ssub"},
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
													arg2: &dCollect{
														ARG:  collect__child__node,
														arg1: &dNode{},
														arg2: &dPredicate{
															arg1: &dEqual{
																ARG: equal__is,
																arg1: &dCollect{
																	ARG:  collect__attributes__cat,
																	arg1: &dNode{},
																},
																arg2: &dElem{
																	DATA: []interface{}{"ssub"},
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
	},
	"corc": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dAnd{
						arg1: &dSort{
							arg1: &dAnd{
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
								arg2: &dCmp{
									ARG: cmp__gt,
									arg1: &dFunction{
										ARG: function__count__1__args,
										arg1: &dArg{
											arg1: &dCollect{
												ARG: collect__child__node,
												arg1: &dCollect{
													ARG:  collect__descendant__or__self__type__node,
													arg1: &dNode{},
												},
												arg2: &dPredicate{
													arg1: &dSort{
														arg1: &dAnd{
															arg1: &dCollect{
																ARG:  collect__child__node,
																arg1: &dNode{},
																arg2: &dPredicate{
																	arg1: &dOr{
																		arg1: &dEqual{
																			ARG: equal__is,
																			arg1: &dCollect{
																				ARG:  collect__attributes__lemma,
																				arg1: &dNode{},
																			},
																			arg2: &dElem{
																				DATA: []interface{}{"hoe", "deste"},
																				arg1: &dCollect{
																					ARG:  collect__attributes__lemma,
																					arg1: &dNode{},
																				},
																			},
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
																								ARG:  collect__attributes__lemma,
																								arg1: &dNode{},
																							},
																							arg2: &dElem{
																								DATA: []interface{}{"des"},
																								arg1: &dCollect{
																									ARG:  collect__attributes__lemma,
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
																						arg1: &dEqual{
																							ARG: equal__is,
																							arg1: &dCollect{
																								ARG:  collect__attributes__lemma,
																								arg1: &dNode{},
																							},
																							arg2: &dElem{
																								DATA: []interface{}{"te"},
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
															arg2: &dCollect{
																ARG:  collect__child__node,
																arg1: &dNode{},
																arg2: &dPredicate{
																	arg1: &dEqual{
																		ARG: equal__is,
																		arg1: &dCollect{
																			ARG:  collect__attributes__graad,
																			arg1: &dNode{},
																		},
																		arg2: &dElem{
																			DATA: []interface{}{"comp"},
																			arg1: &dCollect{
																				ARG:  collect__attributes__graad,
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
									arg2: &dElem{
										DATA: []interface{}{1},
										arg1: &dFunction{
											ARG: function__count__1__args,
											arg1: &dArg{
												arg1: &dCollect{
													ARG: collect__child__node,
													arg1: &dCollect{
														ARG:  collect__descendant__or__self__type__node,
														arg1: &dNode{},
													},
													arg2: &dPredicate{
														arg1: &dSort{
															arg1: &dAnd{
																arg1: &dCollect{
																	ARG:  collect__child__node,
																	arg1: &dNode{},
																	arg2: &dPredicate{
																		arg1: &dOr{
																			arg1: &dEqual{
																				ARG: equal__is,
																				arg1: &dCollect{
																					ARG:  collect__attributes__lemma,
																					arg1: &dNode{},
																				},
																				arg2: &dElem{
																					DATA: []interface{}{"hoe", "deste"},
																					arg1: &dCollect{
																						ARG:  collect__attributes__lemma,
																						arg1: &dNode{},
																					},
																				},
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
																									ARG:  collect__attributes__lemma,
																									arg1: &dNode{},
																								},
																								arg2: &dElem{
																									DATA: []interface{}{"des"},
																									arg1: &dCollect{
																										ARG:  collect__attributes__lemma,
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
																							arg1: &dEqual{
																								ARG: equal__is,
																								arg1: &dCollect{
																									ARG:  collect__attributes__lemma,
																									arg1: &dNode{},
																								},
																								arg2: &dElem{
																									DATA: []interface{}{"te"},
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
																arg2: &dCollect{
																	ARG:  collect__child__node,
																	arg1: &dNode{},
																	arg2: &dPredicate{
																		arg1: &dEqual{
																			ARG: equal__is,
																			arg1: &dCollect{
																				ARG:  collect__attributes__graad,
																				arg1: &dNode{},
																			},
																			arg2: &dElem{
																				DATA: []interface{}{"comp"},
																				arg1: &dCollect{
																					ARG:  collect__attributes__graad,
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
						arg2: &dFunction{
							ARG: function__not__1__args,
							arg1: &dArg{
								arg1: &dSort{
									arg1: &dCollect{
										ARG: collect__child__node,
										arg1: &dCollect{
											ARG:  collect__descendant__or__self__type__node,
											arg1: &dNode{},
										},
										arg2: &dPredicate{
											arg1: &dSort{
												arg1: &dAnd{
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
													arg2: &dCmp{
														ARG: cmp__gt,
														arg1: &dFunction{
															ARG: function__count__1__args,
															arg1: &dArg{
																arg1: &dCollect{
																	ARG: collect__child__node,
																	arg1: &dCollect{
																		ARG:  collect__descendant__or__self__type__node,
																		arg1: &dNode{},
																	},
																	arg2: &dPredicate{
																		arg1: &dSort{
																			arg1: &dAnd{
																				arg1: &dCollect{
																					ARG:  collect__child__node,
																					arg1: &dNode{},
																					arg2: &dPredicate{
																						arg1: &dOr{
																							arg1: &dEqual{
																								ARG: equal__is,
																								arg1: &dCollect{
																									ARG:  collect__attributes__lemma,
																									arg1: &dNode{},
																								},
																								arg2: &dElem{
																									DATA: []interface{}{"hoe", "deste"},
																									arg1: &dCollect{
																										ARG:  collect__attributes__lemma,
																										arg1: &dNode{},
																									},
																								},
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
																													ARG:  collect__attributes__lemma,
																													arg1: &dNode{},
																												},
																												arg2: &dElem{
																													DATA: []interface{}{"des"},
																													arg1: &dCollect{
																														ARG:  collect__attributes__lemma,
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
																											arg1: &dEqual{
																												ARG: equal__is,
																												arg1: &dCollect{
																													ARG:  collect__attributes__lemma,
																													arg1: &dNode{},
																												},
																												arg2: &dElem{
																													DATA: []interface{}{"te"},
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
																				arg2: &dCollect{
																					ARG:  collect__child__node,
																					arg1: &dNode{},
																					arg2: &dPredicate{
																						arg1: &dEqual{
																							ARG: equal__is,
																							arg1: &dCollect{
																								ARG:  collect__attributes__graad,
																								arg1: &dNode{},
																							},
																							arg2: &dElem{
																								DATA: []interface{}{"comp"},
																								arg1: &dCollect{
																									ARG:  collect__attributes__graad,
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
														arg2: &dElem{
															DATA: []interface{}{1},
															arg1: &dFunction{
																ARG: function__count__1__args,
																arg1: &dArg{
																	arg1: &dCollect{
																		ARG: collect__child__node,
																		arg1: &dCollect{
																			ARG:  collect__descendant__or__self__type__node,
																			arg1: &dNode{},
																		},
																		arg2: &dPredicate{
																			arg1: &dSort{
																				arg1: &dAnd{
																					arg1: &dCollect{
																						ARG:  collect__child__node,
																						arg1: &dNode{},
																						arg2: &dPredicate{
																							arg1: &dOr{
																								arg1: &dEqual{
																									ARG: equal__is,
																									arg1: &dCollect{
																										ARG:  collect__attributes__lemma,
																										arg1: &dNode{},
																									},
																									arg2: &dElem{
																										DATA: []interface{}{"hoe", "deste"},
																										arg1: &dCollect{
																											ARG:  collect__attributes__lemma,
																											arg1: &dNode{},
																										},
																									},
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
																														ARG:  collect__attributes__lemma,
																														arg1: &dNode{},
																													},
																													arg2: &dElem{
																														DATA: []interface{}{"des"},
																														arg1: &dCollect{
																															ARG:  collect__attributes__lemma,
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
																												arg1: &dEqual{
																													ARG: equal__is,
																													arg1: &dCollect{
																														ARG:  collect__attributes__lemma,
																														arg1: &dNode{},
																													},
																													arg2: &dElem{
																														DATA: []interface{}{"te"},
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
																					arg2: &dCollect{
																						ARG:  collect__child__node,
																						arg1: &dNode{},
																						arg2: &dPredicate{
																							arg1: &dEqual{
																								ARG: equal__is,
																								arg1: &dCollect{
																									ARG:  collect__attributes__graad,
																									arg1: &dNode{},
																								},
																								arg2: &dElem{
																									DATA: []interface{}{"comp"},
																									arg1: &dCollect{
																										ARG:  collect__attributes__graad,
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
							},
						},
					},
				},
			},
		},
	},
	"cc": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
										DATA: []interface{}{"obcomp"},
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
	"cczo": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
											DATA: []interface{}{"obcomp"},
											arg1: &dCollect{
												ARG:  collect__attributes__rel,
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
								ARG:  collect__attributes__lemma,
								arg1: &dNode{},
							},
							arg2: &dElem{
								DATA: []interface{}{"zo"},
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
	"cceven": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
											DATA: []interface{}{"obcomp"},
											arg1: &dCollect{
												ARG:  collect__attributes__rel,
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
								ARG:  collect__attributes__lemma,
								arg1: &dNode{},
							},
							arg2: &dElem{
								DATA: []interface{}{"even"},
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
	"ccca": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
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
										DATA: []interface{}{"hd"},
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
												DATA: []interface{}{"obcomp"},
												arg1: &dCollect{
													ARG:  collect__attributes__rel,
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
									ARG:  collect__attributes__pt,
									arg1: &dNode{},
								},
								arg2: &dElem{
									DATA: []interface{}{"adj"},
									arg1: &dCollect{
										ARG:  collect__attributes__pt,
										arg1: &dNode{},
									},
								},
							},
						},
						arg2: &dEqual{
							ARG: equal__is,
							arg1: &dCollect{
								ARG:  collect__attributes__graad,
								arg1: &dNode{},
							},
							arg2: &dElem{
								DATA: []interface{}{"comp"},
								arg1: &dCollect{
									ARG:  collect__attributes__graad,
									arg1: &dNode{},
								},
							},
						},
					},
				},
			},
		},
	},
	"ccdannp": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
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
										DATA: []interface{}{"hd"},
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
													DATA: []interface{}{"obcomp"},
													arg1: &dCollect{
														ARG:  collect__attributes__rel,
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
																DATA: []interface{}{"body"},
																arg1: &dCollect{
																	ARG:  collect__attributes__rel,
																	arg1: &dNode{},
																},
															},
														},
														arg2: &dSort{
															arg1: &dOr{
																arg1: &dSort{
																	arg1: &dOr{
																		arg1: &dOr{
																			arg1: &dOr{
																				arg1: &dOr{
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
																					arg2: &dSort{
																						arg1: &dAnd{
																							arg1: &dEqual{
																								ARG: equal__is,
																								arg1: &dCollect{
																									ARG:  collect__attributes__lcat,
																									arg1: &dNode{},
																								},
																								arg2: &dElem{
																									DATA: []interface{}{"np"},
																									arg1: &dCollect{
																										ARG:  collect__attributes__lcat,
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
																												ARG:  collect__attributes__rel,
																												arg1: &dNode{},
																											},
																											arg2: &dElem{
																												DATA: []interface{}{"hd", "mwp"},
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
																					arg1: &dAnd{
																						arg1: &dEqual{
																							ARG: equal__is,
																							arg1: &dCollect{
																								ARG:  collect__attributes__pt,
																								arg1: &dNode{},
																							},
																							arg2: &dElem{
																								DATA: []interface{}{"n"},
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
																			arg2: &dSort{
																				arg1: &dAnd{
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
																						arg2: &dEqual{
																							ARG: equal__is,
																							arg1: &dCollect{
																								ARG:  collect__attributes__pdtype,
																								arg1: &dNode{},
																							},
																							arg2: &dElem{
																								DATA: []interface{}{"pron"},
																								arg1: &dCollect{
																									ARG:  collect__attributes__pdtype,
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
																		arg2: &dSort{
																			arg1: &dAnd{
																				arg1: &dAnd{
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
																					arg2: &dFunction{
																						ARG: function__not__1__args,
																						arg1: &dArg{
																							arg1: &dSort{
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
																				},
																				arg2: &dEqual{
																					ARG: equal__is,
																					arg1: &dCollect{
																						ARG:  collect__attributes__rel,
																						arg1: &dNode{},
																					},
																					arg2: &dElem{
																						DATA: []interface{}{"su", "obj1", "obj2", "app"},
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
																				arg1: &dSort{
																					arg1: &dOr{
																						arg1: &dOr{
																							arg1: &dOr{
																								arg1: &dOr{
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
																									arg2: &dSort{
																										arg1: &dAnd{
																											arg1: &dEqual{
																												ARG: equal__is,
																												arg1: &dCollect{
																													ARG:  collect__attributes__lcat,
																													arg1: &dNode{},
																												},
																												arg2: &dElem{
																													DATA: []interface{}{"np"},
																													arg1: &dCollect{
																														ARG:  collect__attributes__lcat,
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
																																ARG:  collect__attributes__rel,
																																arg1: &dNode{},
																															},
																															arg2: &dElem{
																																DATA: []interface{}{"hd", "mwp"},
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
																									arg1: &dAnd{
																										arg1: &dEqual{
																											ARG: equal__is,
																											arg1: &dCollect{
																												ARG:  collect__attributes__pt,
																												arg1: &dNode{},
																											},
																											arg2: &dElem{
																												DATA: []interface{}{"n"},
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
																							arg2: &dSort{
																								arg1: &dAnd{
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
																										arg2: &dEqual{
																											ARG: equal__is,
																											arg1: &dCollect{
																												ARG:  collect__attributes__pdtype,
																												arg1: &dNode{},
																											},
																											arg2: &dElem{
																												DATA: []interface{}{"pron"},
																												arg1: &dCollect{
																													ARG:  collect__attributes__pdtype,
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
																						arg2: &dSort{
																							arg1: &dAnd{
																								arg1: &dAnd{
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
																									arg2: &dFunction{
																										ARG: function__not__1__args,
																										arg1: &dArg{
																											arg1: &dSort{
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
																								},
																								arg2: &dEqual{
																									ARG: equal__is,
																									arg1: &dCollect{
																										ARG:  collect__attributes__rel,
																										arg1: &dNode{},
																									},
																									arg2: &dElem{
																										DATA: []interface{}{"su", "obj1", "obj2", "app"},
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
								},
							},
							arg2: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG:  collect__attributes__pt,
									arg1: &dNode{},
								},
								arg2: &dElem{
									DATA: []interface{}{"adj"},
									arg1: &dCollect{
										ARG:  collect__attributes__pt,
										arg1: &dNode{},
									},
								},
							},
						},
						arg2: &dEqual{
							ARG: equal__is,
							arg1: &dCollect{
								ARG:  collect__attributes__graad,
								arg1: &dNode{},
							},
							arg2: &dElem{
								DATA: []interface{}{"comp"},
								arg1: &dCollect{
									ARG:  collect__attributes__graad,
									arg1: &dNode{},
								},
							},
						},
					},
				},
			},
		},
	},
	"ccdanvs": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
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
										DATA: []interface{}{"hd"},
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
													DATA: []interface{}{"obcomp"},
													arg1: &dCollect{
														ARG:  collect__attributes__rel,
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
																DATA: []interface{}{"body"},
																arg1: &dCollect{
																	ARG:  collect__attributes__rel,
																	arg1: &dNode{},
																},
															},
														},
														arg2: &dSort{
															arg1: &dOr{
																arg1: &dOr{
																	arg1: &dSort{
																		arg1: &dOr{
																			arg1: &dEqual{
																				ARG: equal__is,
																				arg1: &dCollect{
																					ARG:  collect__attributes__cat,
																					arg1: &dNode{},
																				},
																				arg2: &dElem{
																					DATA: []interface{}{"inf", "ti", "ssub", "oti", "ppart"},
																					arg1: &dCollect{
																						ARG:  collect__attributes__cat,
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
																					DATA: []interface{}{"smain", "sv1"},
																					arg1: &dCollect{
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
																			ARG:  collect__attributes__cat,
																			arg1: &dNode{},
																		},
																		arg2: &dElem{
																			DATA: []interface{}{"cp", "ssub"},
																			arg1: &dCollect{
																				ARG:  collect__attributes__cat,
																				arg1: &dNode{},
																			},
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
																		DATA: []interface{}{"ww"},
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
								},
							},
							arg2: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG:  collect__attributes__pt,
									arg1: &dNode{},
								},
								arg2: &dElem{
									DATA: []interface{}{"adj"},
									arg1: &dCollect{
										ARG:  collect__attributes__pt,
										arg1: &dNode{},
									},
								},
							},
						},
						arg2: &dEqual{
							ARG: equal__is,
							arg1: &dCollect{
								ARG:  collect__attributes__graad,
								arg1: &dNode{},
							},
							arg2: &dElem{
								DATA: []interface{}{"comp"},
								arg1: &dCollect{
									ARG:  collect__attributes__graad,
									arg1: &dNode{},
								},
							},
						},
					},
				},
			},
		},
	},
	"ccdanpp": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
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
										DATA: []interface{}{"hd"},
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
													DATA: []interface{}{"obcomp"},
													arg1: &dCollect{
														ARG:  collect__attributes__rel,
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
																DATA: []interface{}{"body"},
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
										},
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
									DATA: []interface{}{"adj"},
									arg1: &dCollect{
										ARG:  collect__attributes__pt,
										arg1: &dNode{},
									},
								},
							},
						},
						arg2: &dEqual{
							ARG: equal__is,
							arg1: &dCollect{
								ARG:  collect__attributes__graad,
								arg1: &dNode{},
							},
							arg2: &dElem{
								DATA: []interface{}{"comp"},
								arg1: &dCollect{
									ARG:  collect__attributes__graad,
									arg1: &dNode{},
								},
							},
						},
					},
				},
			},
		},
	},
	"ccdanav": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
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
										DATA: []interface{}{"hd"},
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
													DATA: []interface{}{"obcomp"},
													arg1: &dCollect{
														ARG:  collect__attributes__rel,
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
																DATA: []interface{}{"body"},
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
																		ARG:  collect__attributes__cat,
																		arg1: &dNode{},
																	},
																	arg2: &dElem{
																		DATA: []interface{}{"advp", "ap"},
																		arg1: &dCollect{
																			ARG:  collect__attributes__cat,
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
																		DATA: []interface{}{"adj", "bw"},
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
								},
							},
							arg2: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG:  collect__attributes__pt,
									arg1: &dNode{},
								},
								arg2: &dElem{
									DATA: []interface{}{"adj"},
									arg1: &dCollect{
										ARG:  collect__attributes__pt,
										arg1: &dNode{},
									},
								},
							},
						},
						arg2: &dEqual{
							ARG: equal__is,
							arg1: &dCollect{
								ARG:  collect__attributes__graad,
								arg1: &dNode{},
							},
							arg2: &dElem{
								DATA: []interface{}{"comp"},
								arg1: &dCollect{
									ARG:  collect__attributes__graad,
									arg1: &dNode{},
								},
							},
						},
					},
				},
			},
		},
	},
	"ccmm": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
											DATA: []interface{}{"obcomp"},
											arg1: &dCollect{
												ARG:  collect__attributes__rel,
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
								ARG:  collect__attributes__lemma,
								arg1: &dNode{},
							},
							arg2: &dElem{
								DATA: []interface{}{"veel", "minder", "weinig"},
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
	"ccmdnp": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
												DATA: []interface{}{"obcomp"},
												arg1: &dCollect{
													ARG:  collect__attributes__rel,
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
															DATA: []interface{}{"body"},
															arg1: &dCollect{
																ARG:  collect__attributes__rel,
																arg1: &dNode{},
															},
														},
													},
													arg2: &dSort{
														arg1: &dOr{
															arg1: &dSort{
																arg1: &dOr{
																	arg1: &dOr{
																		arg1: &dOr{
																			arg1: &dOr{
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
																				arg2: &dSort{
																					arg1: &dAnd{
																						arg1: &dEqual{
																							ARG: equal__is,
																							arg1: &dCollect{
																								ARG:  collect__attributes__lcat,
																								arg1: &dNode{},
																							},
																							arg2: &dElem{
																								DATA: []interface{}{"np"},
																								arg1: &dCollect{
																									ARG:  collect__attributes__lcat,
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
																											ARG:  collect__attributes__rel,
																											arg1: &dNode{},
																										},
																										arg2: &dElem{
																											DATA: []interface{}{"hd", "mwp"},
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
																				arg1: &dAnd{
																					arg1: &dEqual{
																						ARG: equal__is,
																						arg1: &dCollect{
																							ARG:  collect__attributes__pt,
																							arg1: &dNode{},
																						},
																						arg2: &dElem{
																							DATA: []interface{}{"n"},
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
																		arg2: &dSort{
																			arg1: &dAnd{
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
																					arg2: &dEqual{
																						ARG: equal__is,
																						arg1: &dCollect{
																							ARG:  collect__attributes__pdtype,
																							arg1: &dNode{},
																						},
																						arg2: &dElem{
																							DATA: []interface{}{"pron"},
																							arg1: &dCollect{
																								ARG:  collect__attributes__pdtype,
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
																	arg2: &dSort{
																		arg1: &dAnd{
																			arg1: &dAnd{
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
																				arg2: &dFunction{
																					ARG: function__not__1__args,
																					arg1: &dArg{
																						arg1: &dSort{
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
																			},
																			arg2: &dEqual{
																				ARG: equal__is,
																				arg1: &dCollect{
																					ARG:  collect__attributes__rel,
																					arg1: &dNode{},
																				},
																				arg2: &dElem{
																					DATA: []interface{}{"su", "obj1", "obj2", "app"},
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
																			arg1: &dSort{
																				arg1: &dOr{
																					arg1: &dOr{
																						arg1: &dOr{
																							arg1: &dOr{
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
																								arg2: &dSort{
																									arg1: &dAnd{
																										arg1: &dEqual{
																											ARG: equal__is,
																											arg1: &dCollect{
																												ARG:  collect__attributes__lcat,
																												arg1: &dNode{},
																											},
																											arg2: &dElem{
																												DATA: []interface{}{"np"},
																												arg1: &dCollect{
																													ARG:  collect__attributes__lcat,
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
																															ARG:  collect__attributes__rel,
																															arg1: &dNode{},
																														},
																														arg2: &dElem{
																															DATA: []interface{}{"hd", "mwp"},
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
																								arg1: &dAnd{
																									arg1: &dEqual{
																										ARG: equal__is,
																										arg1: &dCollect{
																											ARG:  collect__attributes__pt,
																											arg1: &dNode{},
																										},
																										arg2: &dElem{
																											DATA: []interface{}{"n"},
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
																						arg2: &dSort{
																							arg1: &dAnd{
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
																									arg2: &dEqual{
																										ARG: equal__is,
																										arg1: &dCollect{
																											ARG:  collect__attributes__pdtype,
																											arg1: &dNode{},
																										},
																										arg2: &dElem{
																											DATA: []interface{}{"pron"},
																											arg1: &dCollect{
																												ARG:  collect__attributes__pdtype,
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
																					arg2: &dSort{
																						arg1: &dAnd{
																							arg1: &dAnd{
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
																								arg2: &dFunction{
																									ARG: function__not__1__args,
																									arg1: &dArg{
																										arg1: &dSort{
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
																							},
																							arg2: &dEqual{
																								ARG: equal__is,
																								arg1: &dCollect{
																									ARG:  collect__attributes__rel,
																									arg1: &dNode{},
																								},
																								arg2: &dElem{
																									DATA: []interface{}{"su", "obj1", "obj2", "app"},
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
							},
						},
						arg2: &dEqual{
							ARG: equal__is,
							arg1: &dCollect{
								ARG:  collect__attributes__lemma,
								arg1: &dNode{},
							},
							arg2: &dElem{
								DATA: []interface{}{"veel", "minder", "weinig"},
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
	"ccmdvs": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
												DATA: []interface{}{"obcomp"},
												arg1: &dCollect{
													ARG:  collect__attributes__rel,
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
															DATA: []interface{}{"body"},
															arg1: &dCollect{
																ARG:  collect__attributes__rel,
																arg1: &dNode{},
															},
														},
													},
													arg2: &dSort{
														arg1: &dOr{
															arg1: &dOr{
																arg1: &dSort{
																	arg1: &dOr{
																		arg1: &dEqual{
																			ARG: equal__is,
																			arg1: &dCollect{
																				ARG:  collect__attributes__cat,
																				arg1: &dNode{},
																			},
																			arg2: &dElem{
																				DATA: []interface{}{"inf", "ti", "ssub", "oti", "ppart"},
																				arg1: &dCollect{
																					ARG:  collect__attributes__cat,
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
																				DATA: []interface{}{"smain", "sv1"},
																				arg1: &dCollect{
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
																		ARG:  collect__attributes__cat,
																		arg1: &dNode{},
																	},
																	arg2: &dElem{
																		DATA: []interface{}{"cp", "ssub"},
																		arg1: &dCollect{
																			ARG:  collect__attributes__cat,
																			arg1: &dNode{},
																		},
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
																	DATA: []interface{}{"ww"},
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
							},
						},
						arg2: &dEqual{
							ARG: equal__is,
							arg1: &dCollect{
								ARG:  collect__attributes__lemma,
								arg1: &dNode{},
							},
							arg2: &dElem{
								DATA: []interface{}{"veel", "minder", "weinig"},
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
	"ccmdpp": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
												DATA: []interface{}{"obcomp"},
												arg1: &dCollect{
													ARG:  collect__attributes__rel,
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
															DATA: []interface{}{"body"},
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
								DATA: []interface{}{"veel", "minder", "weinig"},
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
	"ccmdav": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
												DATA: []interface{}{"obcomp"},
												arg1: &dCollect{
													ARG:  collect__attributes__rel,
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
															DATA: []interface{}{"body"},
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
																	ARG:  collect__attributes__cat,
																	arg1: &dNode{},
																},
																arg2: &dElem{
																	DATA: []interface{}{"advp", "ap"},
																	arg1: &dCollect{
																		ARG:  collect__attributes__cat,
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
																	DATA: []interface{}{"adj", "bw"},
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
							},
						},
						arg2: &dEqual{
							ARG: equal__is,
							arg1: &dCollect{
								ARG:  collect__attributes__lemma,
								arg1: &dNode{},
							},
							arg2: &dElem{
								DATA: []interface{}{"veel", "minder", "weinig"},
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
	"ccnn": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
											DATA: []interface{}{"obcomp"},
											arg1: &dCollect{
												ARG:  collect__attributes__rel,
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
								ARG:  collect__attributes__lemma,
								arg1: &dNode{},
							},
							arg2: &dElem{
								DATA: []interface{}{"niet", "niets", "ander", "anders"},
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
	"ccndnp": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
												DATA: []interface{}{"obcomp"},
												arg1: &dCollect{
													ARG:  collect__attributes__rel,
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
															DATA: []interface{}{"body"},
															arg1: &dCollect{
																ARG:  collect__attributes__rel,
																arg1: &dNode{},
															},
														},
													},
													arg2: &dSort{
														arg1: &dOr{
															arg1: &dSort{
																arg1: &dOr{
																	arg1: &dOr{
																		arg1: &dOr{
																			arg1: &dOr{
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
																				arg2: &dSort{
																					arg1: &dAnd{
																						arg1: &dEqual{
																							ARG: equal__is,
																							arg1: &dCollect{
																								ARG:  collect__attributes__lcat,
																								arg1: &dNode{},
																							},
																							arg2: &dElem{
																								DATA: []interface{}{"np"},
																								arg1: &dCollect{
																									ARG:  collect__attributes__lcat,
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
																											ARG:  collect__attributes__rel,
																											arg1: &dNode{},
																										},
																										arg2: &dElem{
																											DATA: []interface{}{"hd", "mwp"},
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
																				arg1: &dAnd{
																					arg1: &dEqual{
																						ARG: equal__is,
																						arg1: &dCollect{
																							ARG:  collect__attributes__pt,
																							arg1: &dNode{},
																						},
																						arg2: &dElem{
																							DATA: []interface{}{"n"},
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
																		arg2: &dSort{
																			arg1: &dAnd{
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
																					arg2: &dEqual{
																						ARG: equal__is,
																						arg1: &dCollect{
																							ARG:  collect__attributes__pdtype,
																							arg1: &dNode{},
																						},
																						arg2: &dElem{
																							DATA: []interface{}{"pron"},
																							arg1: &dCollect{
																								ARG:  collect__attributes__pdtype,
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
																	arg2: &dSort{
																		arg1: &dAnd{
																			arg1: &dAnd{
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
																				arg2: &dFunction{
																					ARG: function__not__1__args,
																					arg1: &dArg{
																						arg1: &dSort{
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
																			},
																			arg2: &dEqual{
																				ARG: equal__is,
																				arg1: &dCollect{
																					ARG:  collect__attributes__rel,
																					arg1: &dNode{},
																				},
																				arg2: &dElem{
																					DATA: []interface{}{"su", "obj1", "obj2", "app"},
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
																			arg1: &dSort{
																				arg1: &dOr{
																					arg1: &dOr{
																						arg1: &dOr{
																							arg1: &dOr{
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
																								arg2: &dSort{
																									arg1: &dAnd{
																										arg1: &dEqual{
																											ARG: equal__is,
																											arg1: &dCollect{
																												ARG:  collect__attributes__lcat,
																												arg1: &dNode{},
																											},
																											arg2: &dElem{
																												DATA: []interface{}{"np"},
																												arg1: &dCollect{
																													ARG:  collect__attributes__lcat,
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
																															ARG:  collect__attributes__rel,
																															arg1: &dNode{},
																														},
																														arg2: &dElem{
																															DATA: []interface{}{"hd", "mwp"},
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
																								arg1: &dAnd{
																									arg1: &dEqual{
																										ARG: equal__is,
																										arg1: &dCollect{
																											ARG:  collect__attributes__pt,
																											arg1: &dNode{},
																										},
																										arg2: &dElem{
																											DATA: []interface{}{"n"},
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
																						arg2: &dSort{
																							arg1: &dAnd{
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
																									arg2: &dEqual{
																										ARG: equal__is,
																										arg1: &dCollect{
																											ARG:  collect__attributes__pdtype,
																											arg1: &dNode{},
																										},
																										arg2: &dElem{
																											DATA: []interface{}{"pron"},
																											arg1: &dCollect{
																												ARG:  collect__attributes__pdtype,
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
																					arg2: &dSort{
																						arg1: &dAnd{
																							arg1: &dAnd{
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
																								arg2: &dFunction{
																									ARG: function__not__1__args,
																									arg1: &dArg{
																										arg1: &dSort{
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
																							},
																							arg2: &dEqual{
																								ARG: equal__is,
																								arg1: &dCollect{
																									ARG:  collect__attributes__rel,
																									arg1: &dNode{},
																								},
																								arg2: &dElem{
																									DATA: []interface{}{"su", "obj1", "obj2", "app"},
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
							},
						},
						arg2: &dEqual{
							ARG: equal__is,
							arg1: &dCollect{
								ARG:  collect__attributes__lemma,
								arg1: &dNode{},
							},
							arg2: &dElem{
								DATA: []interface{}{"niet", "niets", "ander", "anders"},
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
	"ccndvs": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
												DATA: []interface{}{"obcomp"},
												arg1: &dCollect{
													ARG:  collect__attributes__rel,
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
															DATA: []interface{}{"body"},
															arg1: &dCollect{
																ARG:  collect__attributes__rel,
																arg1: &dNode{},
															},
														},
													},
													arg2: &dSort{
														arg1: &dOr{
															arg1: &dOr{
																arg1: &dSort{
																	arg1: &dOr{
																		arg1: &dEqual{
																			ARG: equal__is,
																			arg1: &dCollect{
																				ARG:  collect__attributes__cat,
																				arg1: &dNode{},
																			},
																			arg2: &dElem{
																				DATA: []interface{}{"inf", "ti", "ssub", "oti", "ppart"},
																				arg1: &dCollect{
																					ARG:  collect__attributes__cat,
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
																				DATA: []interface{}{"smain", "sv1"},
																				arg1: &dCollect{
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
																		ARG:  collect__attributes__cat,
																		arg1: &dNode{},
																	},
																	arg2: &dElem{
																		DATA: []interface{}{"cp", "ssub"},
																		arg1: &dCollect{
																			ARG:  collect__attributes__cat,
																			arg1: &dNode{},
																		},
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
																	DATA: []interface{}{"ww"},
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
							},
						},
						arg2: &dEqual{
							ARG: equal__is,
							arg1: &dCollect{
								ARG:  collect__attributes__lemma,
								arg1: &dNode{},
							},
							arg2: &dElem{
								DATA: []interface{}{"niet", "niets", "ander", "anders"},
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
	"ccndpp": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
												DATA: []interface{}{"obcomp"},
												arg1: &dCollect{
													ARG:  collect__attributes__rel,
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
															DATA: []interface{}{"body"},
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
								DATA: []interface{}{"niet", "niets", "ander", "anders"},
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
	"ccndav": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
												DATA: []interface{}{"obcomp"},
												arg1: &dCollect{
													ARG:  collect__attributes__rel,
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
															DATA: []interface{}{"body"},
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
																	ARG:  collect__attributes__cat,
																	arg1: &dNode{},
																},
																arg2: &dElem{
																	DATA: []interface{}{"advp", "ap"},
																	arg1: &dCollect{
																		ARG:  collect__attributes__cat,
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
																	DATA: []interface{}{"adj", "bw"},
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
							},
						},
						arg2: &dEqual{
							ARG: equal__is,
							arg1: &dCollect{
								ARG:  collect__attributes__lemma,
								arg1: &dNode{},
							},
							arg2: &dElem{
								DATA: []interface{}{"niet", "niets", "ander", "anders"},
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
	"conj": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
	"crd0": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
							},
						},
					},
				},
			},
		},
	},
	"crd1": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
						arg2: &dEqual{
							ARG: equal__is,
							arg1: &dFunction{
								ARG: function__count__1__args,
								arg1: &dArg{
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
							},
							arg2: &dElem{
								DATA: []interface{}{1},
								arg1: &dFunction{
									ARG: function__count__1__args,
									arg1: &dArg{
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
								},
							},
						},
					},
				},
			},
		},
	},
	"crd1en": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dAnd{
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
							arg2: &dEqual{
								ARG: equal__is,
								arg1: &dFunction{
									ARG: function__count__1__args,
									arg1: &dArg{
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
								},
								arg2: &dElem{
									DATA: []interface{}{1},
									arg1: &dFunction{
										ARG: function__count__1__args,
										arg1: &dArg{
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
											DATA: []interface{}{"crd"},
											arg1: &dCollect{
												ARG:  collect__attributes__rel,
												arg1: &dNode{},
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
											DATA: []interface{}{"en"},
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
	"crd1of": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dAnd{
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
							arg2: &dEqual{
								ARG: equal__is,
								arg1: &dFunction{
									ARG: function__count__1__args,
									arg1: &dArg{
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
								},
								arg2: &dElem{
									DATA: []interface{}{1},
									arg1: &dFunction{
										ARG: function__count__1__args,
										arg1: &dArg{
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
											DATA: []interface{}{"crd"},
											arg1: &dCollect{
												ARG:  collect__attributes__rel,
												arg1: &dNode{},
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
											DATA: []interface{}{"of"},
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
	"crd1maa": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dAnd{
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
							arg2: &dEqual{
								ARG: equal__is,
								arg1: &dFunction{
									ARG: function__count__1__args,
									arg1: &dArg{
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
								},
								arg2: &dElem{
									DATA: []interface{}{1},
									arg1: &dFunction{
										ARG: function__count__1__args,
										arg1: &dArg{
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
											DATA: []interface{}{"crd"},
											arg1: &dCollect{
												ARG:  collect__attributes__rel,
												arg1: &dNode{},
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
											DATA: []interface{}{"maar"},
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
	"crd1enz": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dAnd{
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
							arg2: &dEqual{
								ARG: equal__is,
								arg1: &dFunction{
									ARG: function__count__1__args,
									arg1: &dArg{
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
								},
								arg2: &dElem{
									DATA: []interface{}{1},
									arg1: &dFunction{
										ARG: function__count__1__args,
										arg1: &dArg{
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
									},
								},
							},
						},
						arg2: &dEqual{
							ARG: equal__is,
							arg1: &dCollect{
								ARG: collect__attributes__end,
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
							arg2: &dCollect{
								ARG:  collect__attributes__end,
								arg1: &dNode{},
							},
						},
					},
				},
			},
		},
	},
	"crd2": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
						arg2: &dEqual{
							ARG: equal__is,
							arg1: &dFunction{
								ARG: function__count__1__args,
								arg1: &dArg{
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
							},
							arg2: &dElem{
								DATA: []interface{}{2},
								arg1: &dFunction{
									ARG: function__count__1__args,
									arg1: &dArg{
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
								},
							},
						},
					},
				},
			},
		},
	},
	"crd22": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dAnd{
						arg1: &dAnd{
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
								arg2: &dCmp{
									ARG: cmp__gt,
									arg1: &dFunction{
										ARG: function__count__1__args,
										arg1: &dArg{
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
									},
									arg2: &dElem{
										DATA: []interface{}{1},
										arg1: &dFunction{
											ARG: function__count__1__args,
											arg1: &dArg{
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
										},
									},
								},
							},
							arg2: &dCollect{
								ARG:  collect__child__node,
								arg1: &dNode{},
								arg2: &dPredicate{
									arg1: &dPredicate{
										arg1: &dFunction{
											ARG: function__first__0__args,
										},
									},
									arg2: &dEqual{
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
						arg2: &dCollect{
							ARG:  collect__child__node,
							arg1: &dNode{},
							arg2: &dPredicate{
								arg1: &dPredicate{
									arg1: &dElem{
										DATA: []interface{}{3},
									},
								},
								arg2: &dEqual{
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
				},
			},
		},
	},
	"crd2p": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
						arg2: &dCmp{
							ARG: cmp__gt,
							arg1: &dFunction{
								ARG: function__count__1__args,
								arg1: &dArg{
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
							},
							arg2: &dElem{
								DATA: []interface{}{2},
								arg1: &dFunction{
									ARG: function__count__1__args,
									arg1: &dArg{
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
								},
							},
						},
					},
				},
			},
		},
	},
	"cnj1": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
						arg2: &dEqual{
							ARG: equal__is,
							arg1: &dFunction{
								ARG: function__count__1__args,
								arg1: &dArg{
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
							},
							arg2: &dElem{
								DATA: []interface{}{1},
								arg1: &dFunction{
									ARG: function__count__1__args,
									arg1: &dArg{
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
								},
							},
						},
					},
				},
			},
		},
	},
	"cnj2": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
						arg2: &dEqual{
							ARG: equal__is,
							arg1: &dFunction{
								ARG: function__count__1__args,
								arg1: &dArg{
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
							},
							arg2: &dElem{
								DATA: []interface{}{2},
								arg1: &dFunction{
									ARG: function__count__1__args,
									arg1: &dArg{
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
								},
							},
						},
					},
				},
			},
		},
	},
	"cnj3": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
						arg2: &dEqual{
							ARG: equal__is,
							arg1: &dFunction{
								ARG: function__count__1__args,
								arg1: &dArg{
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
							},
							arg2: &dElem{
								DATA: []interface{}{3},
								arg1: &dFunction{
									ARG: function__count__1__args,
									arg1: &dArg{
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
								},
							},
						},
					},
				},
			},
		},
	},
	"cnj4": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
						arg2: &dEqual{
							ARG: equal__is,
							arg1: &dFunction{
								ARG: function__count__1__args,
								arg1: &dArg{
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
							},
							arg2: &dElem{
								DATA: []interface{}{4},
								arg1: &dFunction{
									ARG: function__count__1__args,
									arg1: &dArg{
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
								},
							},
						},
					},
				},
			},
		},
	},
	"cnj5": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
						arg2: &dEqual{
							ARG: equal__is,
							arg1: &dFunction{
								ARG: function__count__1__args,
								arg1: &dArg{
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
							},
							arg2: &dElem{
								DATA: []interface{}{5},
								arg1: &dFunction{
									ARG: function__count__1__args,
									arg1: &dArg{
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
								},
							},
						},
					},
				},
			},
		},
	},
	"cnj6": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
						arg2: &dEqual{
							ARG: equal__is,
							arg1: &dFunction{
								ARG: function__count__1__args,
								arg1: &dArg{
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
							},
							arg2: &dElem{
								DATA: []interface{}{6},
								arg1: &dFunction{
									ARG: function__count__1__args,
									arg1: &dArg{
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
								},
							},
						},
					},
				},
			},
		},
	},
	"cnj6p": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
						arg2: &dCmp{
							ARG: cmp__gt,
							arg1: &dFunction{
								ARG: function__count__1__args,
								arg1: &dArg{
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
							},
							arg2: &dElem{
								DATA: []interface{}{6},
								arg1: &dFunction{
									ARG: function__count__1__args,
									arg1: &dArg{
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
								},
							},
						},
					},
				},
			},
		},
	},
	"cnjnp": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dAnd{
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
										arg2: &dSort{
											arg1: &dOr{
												arg1: &dSort{
													arg1: &dOr{
														arg1: &dOr{
															arg1: &dOr{
																arg1: &dOr{
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
																	arg2: &dSort{
																		arg1: &dAnd{
																			arg1: &dEqual{
																				ARG: equal__is,
																				arg1: &dCollect{
																					ARG:  collect__attributes__lcat,
																					arg1: &dNode{},
																				},
																				arg2: &dElem{
																					DATA: []interface{}{"np"},
																					arg1: &dCollect{
																						ARG:  collect__attributes__lcat,
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
																								ARG:  collect__attributes__rel,
																								arg1: &dNode{},
																							},
																							arg2: &dElem{
																								DATA: []interface{}{"hd", "mwp"},
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
																	arg1: &dAnd{
																		arg1: &dEqual{
																			ARG: equal__is,
																			arg1: &dCollect{
																				ARG:  collect__attributes__pt,
																				arg1: &dNode{},
																			},
																			arg2: &dElem{
																				DATA: []interface{}{"n"},
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
															arg2: &dSort{
																arg1: &dAnd{
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
																		arg2: &dEqual{
																			ARG: equal__is,
																			arg1: &dCollect{
																				ARG:  collect__attributes__pdtype,
																				arg1: &dNode{},
																			},
																			arg2: &dElem{
																				DATA: []interface{}{"pron"},
																				arg1: &dCollect{
																					ARG:  collect__attributes__pdtype,
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
														arg2: &dSort{
															arg1: &dAnd{
																arg1: &dAnd{
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
																	arg2: &dFunction{
																		ARG: function__not__1__args,
																		arg1: &dArg{
																			arg1: &dSort{
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
																},
																arg2: &dEqual{
																	ARG: equal__is,
																	arg1: &dCollect{
																		ARG:  collect__attributes__rel,
																		arg1: &dNode{},
																	},
																	arg2: &dElem{
																		DATA: []interface{}{"su", "obj1", "obj2", "app"},
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
																arg1: &dSort{
																	arg1: &dOr{
																		arg1: &dOr{
																			arg1: &dOr{
																				arg1: &dOr{
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
																					arg2: &dSort{
																						arg1: &dAnd{
																							arg1: &dEqual{
																								ARG: equal__is,
																								arg1: &dCollect{
																									ARG:  collect__attributes__lcat,
																									arg1: &dNode{},
																								},
																								arg2: &dElem{
																									DATA: []interface{}{"np"},
																									arg1: &dCollect{
																										ARG:  collect__attributes__lcat,
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
																												ARG:  collect__attributes__rel,
																												arg1: &dNode{},
																											},
																											arg2: &dElem{
																												DATA: []interface{}{"hd", "mwp"},
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
																					arg1: &dAnd{
																						arg1: &dEqual{
																							ARG: equal__is,
																							arg1: &dCollect{
																								ARG:  collect__attributes__pt,
																								arg1: &dNode{},
																							},
																							arg2: &dElem{
																								DATA: []interface{}{"n"},
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
																			arg2: &dSort{
																				arg1: &dAnd{
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
																						arg2: &dEqual{
																							ARG: equal__is,
																							arg1: &dCollect{
																								ARG:  collect__attributes__pdtype,
																								arg1: &dNode{},
																							},
																							arg2: &dElem{
																								DATA: []interface{}{"pron"},
																								arg1: &dCollect{
																									ARG:  collect__attributes__pdtype,
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
																		arg2: &dSort{
																			arg1: &dAnd{
																				arg1: &dAnd{
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
																					arg2: &dFunction{
																						ARG: function__not__1__args,
																						arg1: &dArg{
																							arg1: &dSort{
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
																				},
																				arg2: &dEqual{
																					ARG: equal__is,
																					arg1: &dCollect{
																						ARG:  collect__attributes__rel,
																						arg1: &dNode{},
																					},
																					arg2: &dElem{
																						DATA: []interface{}{"su", "obj1", "obj2", "app"},
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
						arg2: &dFunction{
							ARG: function__not__1__args,
							arg1: &dArg{
								arg1: &dSort{
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
														DATA: []interface{}{"cnj"},
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
															arg1: &dSort{
																arg1: &dOr{
																	arg1: &dSort{
																		arg1: &dOr{
																			arg1: &dOr{
																				arg1: &dOr{
																					arg1: &dOr{
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
																						arg2: &dSort{
																							arg1: &dAnd{
																								arg1: &dEqual{
																									ARG: equal__is,
																									arg1: &dCollect{
																										ARG:  collect__attributes__lcat,
																										arg1: &dNode{},
																									},
																									arg2: &dElem{
																										DATA: []interface{}{"np"},
																										arg1: &dCollect{
																											ARG:  collect__attributes__lcat,
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
																													ARG:  collect__attributes__rel,
																													arg1: &dNode{},
																												},
																												arg2: &dElem{
																													DATA: []interface{}{"hd", "mwp"},
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
																						arg1: &dAnd{
																							arg1: &dEqual{
																								ARG: equal__is,
																								arg1: &dCollect{
																									ARG:  collect__attributes__pt,
																									arg1: &dNode{},
																								},
																								arg2: &dElem{
																									DATA: []interface{}{"n"},
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
																				arg2: &dSort{
																					arg1: &dAnd{
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
																							arg2: &dEqual{
																								ARG: equal__is,
																								arg1: &dCollect{
																									ARG:  collect__attributes__pdtype,
																									arg1: &dNode{},
																								},
																								arg2: &dElem{
																									DATA: []interface{}{"pron"},
																									arg1: &dCollect{
																										ARG:  collect__attributes__pdtype,
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
																			arg2: &dSort{
																				arg1: &dAnd{
																					arg1: &dAnd{
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
																						arg2: &dFunction{
																							ARG: function__not__1__args,
																							arg1: &dArg{
																								arg1: &dSort{
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
																					},
																					arg2: &dEqual{
																						ARG: equal__is,
																						arg1: &dCollect{
																							ARG:  collect__attributes__rel,
																							arg1: &dNode{},
																						},
																						arg2: &dElem{
																							DATA: []interface{}{"su", "obj1", "obj2", "app"},
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
																					arg1: &dSort{
																						arg1: &dOr{
																							arg1: &dOr{
																								arg1: &dOr{
																									arg1: &dOr{
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
																										arg2: &dSort{
																											arg1: &dAnd{
																												arg1: &dEqual{
																													ARG: equal__is,
																													arg1: &dCollect{
																														ARG:  collect__attributes__lcat,
																														arg1: &dNode{},
																													},
																													arg2: &dElem{
																														DATA: []interface{}{"np"},
																														arg1: &dCollect{
																															ARG:  collect__attributes__lcat,
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
																																	ARG:  collect__attributes__rel,
																																	arg1: &dNode{},
																																},
																																arg2: &dElem{
																																	DATA: []interface{}{"hd", "mwp"},
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
																										arg1: &dAnd{
																											arg1: &dEqual{
																												ARG: equal__is,
																												arg1: &dCollect{
																													ARG:  collect__attributes__pt,
																													arg1: &dNode{},
																												},
																												arg2: &dElem{
																													DATA: []interface{}{"n"},
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
																								arg2: &dSort{
																									arg1: &dAnd{
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
																											arg2: &dEqual{
																												ARG: equal__is,
																												arg1: &dCollect{
																													ARG:  collect__attributes__pdtype,
																													arg1: &dNode{},
																												},
																												arg2: &dElem{
																													DATA: []interface{}{"pron"},
																													arg1: &dCollect{
																														ARG:  collect__attributes__pdtype,
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
																							arg2: &dSort{
																								arg1: &dAnd{
																									arg1: &dAnd{
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
																										arg2: &dFunction{
																											ARG: function__not__1__args,
																											arg1: &dArg{
																												arg1: &dSort{
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
																									},
																									arg2: &dEqual{
																										ARG: equal__is,
																										arg1: &dCollect{
																											ARG:  collect__attributes__rel,
																											arg1: &dNode{},
																										},
																										arg2: &dElem{
																											DATA: []interface{}{"su", "obj1", "obj2", "app"},
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
									},
								},
							},
						},
					},
				},
			},
		},
	},
	"cnjpp": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dAnd{
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
						},
						arg2: &dFunction{
							ARG: function__not__1__args,
							arg1: &dArg{
								arg1: &dSort{
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
														DATA: []interface{}{"cnj"},
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
	"cnjmain": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dAnd{
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
												DATA: []interface{}{"smain"},
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
						arg2: &dFunction{
							ARG: function__not__1__args,
							arg1: &dArg{
								arg1: &dSort{
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
														DATA: []interface{}{"cnj"},
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
																	ARG:  collect__attributes__cat,
																	arg1: &dNode{},
																},
																arg2: &dElem{
																	DATA: []interface{}{"smain"},
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
	},
	"cnjvp": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dAnd{
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
										arg2: &dSort{
											arg1: &dOr{
												arg1: &dEqual{
													ARG: equal__is,
													arg1: &dCollect{
														ARG:  collect__attributes__cat,
														arg1: &dNode{},
													},
													arg2: &dElem{
														DATA: []interface{}{"ssub", "ti", "ppart", "inf"},
														arg1: &dCollect{
															ARG:  collect__attributes__cat,
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
														DATA: []interface{}{"ww"},
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
						arg2: &dFunction{
							ARG: function__not__1__args,
							arg1: &dArg{
								arg1: &dSort{
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
														DATA: []interface{}{"cnj"},
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
																arg1: &dEqual{
																	ARG: equal__is,
																	arg1: &dCollect{
																		ARG:  collect__attributes__cat,
																		arg1: &dNode{},
																	},
																	arg2: &dElem{
																		DATA: []interface{}{"ssub", "ti", "ppart", "inf"},
																		arg1: &dCollect{
																			ARG:  collect__attributes__cat,
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
																		DATA: []interface{}{"ww"},
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
								},
							},
						},
					},
				},
			},
		},
	},
	"cnjcp": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dAnd{
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
												DATA: []interface{}{"cp"},
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
						arg2: &dFunction{
							ARG: function__not__1__args,
							arg1: &dArg{
								arg1: &dSort{
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
														DATA: []interface{}{"cnj"},
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
																	ARG:  collect__attributes__cat,
																	arg1: &dNode{},
																},
																arg2: &dElem{
																	DATA: []interface{}{"cp"},
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
	},
	"np": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dSort{
						arg1: &dOr{
							arg1: &dSort{
								arg1: &dOr{
									arg1: &dOr{
										arg1: &dOr{
											arg1: &dOr{
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
												arg2: &dSort{
													arg1: &dAnd{
														arg1: &dEqual{
															ARG: equal__is,
															arg1: &dCollect{
																ARG:  collect__attributes__lcat,
																arg1: &dNode{},
															},
															arg2: &dElem{
																DATA: []interface{}{"np"},
																arg1: &dCollect{
																	ARG:  collect__attributes__lcat,
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
																			ARG:  collect__attributes__rel,
																			arg1: &dNode{},
																		},
																		arg2: &dElem{
																			DATA: []interface{}{"hd", "mwp"},
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
												arg1: &dAnd{
													arg1: &dEqual{
														ARG: equal__is,
														arg1: &dCollect{
															ARG:  collect__attributes__pt,
															arg1: &dNode{},
														},
														arg2: &dElem{
															DATA: []interface{}{"n"},
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
										arg2: &dSort{
											arg1: &dAnd{
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
													arg2: &dEqual{
														ARG: equal__is,
														arg1: &dCollect{
															ARG:  collect__attributes__pdtype,
															arg1: &dNode{},
														},
														arg2: &dElem{
															DATA: []interface{}{"pron"},
															arg1: &dCollect{
																ARG:  collect__attributes__pdtype,
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
									arg2: &dSort{
										arg1: &dAnd{
											arg1: &dAnd{
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
												arg2: &dFunction{
													ARG: function__not__1__args,
													arg1: &dArg{
														arg1: &dSort{
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
											},
											arg2: &dEqual{
												ARG: equal__is,
												arg1: &dCollect{
													ARG:  collect__attributes__rel,
													arg1: &dNode{},
												},
												arg2: &dElem{
													DATA: []interface{}{"su", "obj1", "obj2", "app"},
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
											arg1: &dSort{
												arg1: &dOr{
													arg1: &dOr{
														arg1: &dOr{
															arg1: &dOr{
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
																arg2: &dSort{
																	arg1: &dAnd{
																		arg1: &dEqual{
																			ARG: equal__is,
																			arg1: &dCollect{
																				ARG:  collect__attributes__lcat,
																				arg1: &dNode{},
																			},
																			arg2: &dElem{
																				DATA: []interface{}{"np"},
																				arg1: &dCollect{
																					ARG:  collect__attributes__lcat,
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
																							ARG:  collect__attributes__rel,
																							arg1: &dNode{},
																						},
																						arg2: &dElem{
																							DATA: []interface{}{"hd", "mwp"},
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
																arg1: &dAnd{
																	arg1: &dEqual{
																		ARG: equal__is,
																		arg1: &dCollect{
																			ARG:  collect__attributes__pt,
																			arg1: &dNode{},
																		},
																		arg2: &dElem{
																			DATA: []interface{}{"n"},
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
														arg2: &dSort{
															arg1: &dAnd{
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
																	arg2: &dEqual{
																		ARG: equal__is,
																		arg1: &dCollect{
																			ARG:  collect__attributes__pdtype,
																			arg1: &dNode{},
																		},
																		arg2: &dElem{
																			DATA: []interface{}{"pron"},
																			arg1: &dCollect{
																				ARG:  collect__attributes__pdtype,
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
													arg2: &dSort{
														arg1: &dAnd{
															arg1: &dAnd{
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
																arg2: &dFunction{
																	ARG: function__not__1__args,
																	arg1: &dArg{
																		arg1: &dSort{
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
															},
															arg2: &dEqual{
																ARG: equal__is,
																arg1: &dCollect{
																	ARG:  collect__attributes__rel,
																	arg1: &dNode{},
																},
																arg2: &dElem{
																	DATA: []interface{}{"su", "obj1", "obj2", "app"},
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
	"pp": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
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
	},
	"ap": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dOr{
						arg1: &dEqual{
							ARG: equal__is,
							arg1: &dCollect{
								ARG:  collect__attributes__cat,
								arg1: &dNode{},
							},
							arg2: &dElem{
								DATA: []interface{}{"ap"},
								arg1: &dCollect{
									ARG:  collect__attributes__cat,
									arg1: &dNode{},
								},
							},
						},
						arg2: &dAnd{
							arg1: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG:  collect__attributes__pt,
									arg1: &dNode{},
								},
								arg2: &dElem{
									DATA: []interface{}{"adj"},
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
	"advp": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dOr{
						arg1: &dEqual{
							ARG: equal__is,
							arg1: &dCollect{
								ARG:  collect__attributes__cat,
								arg1: &dNode{},
							},
							arg2: &dElem{
								DATA: []interface{}{"advp"},
								arg1: &dCollect{
									ARG:  collect__attributes__cat,
									arg1: &dNode{},
								},
							},
						},
						arg2: &dAnd{
							arg1: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG:  collect__attributes__pt,
									arg1: &dNode{},
								},
								arg2: &dElem{
									DATA: []interface{}{"bw"},
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
	"ppnp": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dAnd{
						arg1: &dAnd{
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
							arg2: &dEqual{
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
								DATA: []interface{}{"np"},
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
	},
	"ppap": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dAnd{
						arg1: &dAnd{
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
							arg2: &dEqual{
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
										DATA: []interface{}{"ap"},
										arg1: &dCollect{
											ARG: collect__attributes__cat,
											arg1: &dCollect{
												ARG:  collect__parent__type__node,
												arg1: &dNode{},
											},
										},
									},
								},
								arg2: &dSort{
									arg1: &dAnd{
										arg1: &dAnd{
											arg1: &dSort{
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
															DATA: []interface{}{"ppart"},
															arg1: &dCollect{
																ARG: collect__attributes__cat,
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
																ARG:  collect__parent__type__node,
																arg1: &dNode{},
															},
														},
														arg2: &dElem{
															DATA: []interface{}{"ppres"},
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
											arg2: &dFunction{
												ARG: function__not__1__args,
												arg1: &dArg{
													arg1: &dSort{
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
																DATA: []interface{}{"vc"},
																arg1: &dCollect{
																	ARG: collect__attributes__rel,
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
										arg2: &dFunction{
											ARG: function__not__1__args,
											arg1: &dArg{
												arg1: &dSort{
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
																DATA: []interface{}{"cnj"},
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
																DATA: []interface{}{"vc"},
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
	"pppc": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
								DATA: []interface{}{"pp"},
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
								DATA: []interface{}{"pc"},
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
	"ppld": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
								DATA: []interface{}{"pp"},
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
								DATA: []interface{}{"ld"},
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
	"pppredc": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
								DATA: []interface{}{"pp"},
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
	"ppbep": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dAnd{
						arg1: &dAnd{
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
							arg2: &dEqual{
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
										DATA: []interface{}{"smain", "sv1", "whq", "ssub", "inf"},
										arg1: &dCollect{
											ARG: collect__attributes__cat,
											arg1: &dCollect{
												ARG:  collect__parent__type__node,
												arg1: &dNode{},
											},
										},
									},
								},
								arg2: &dSort{
									arg1: &dAnd{
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
												DATA: []interface{}{"ppart"},
												arg1: &dCollect{
													ARG: collect__attributes__cat,
													arg1: &dCollect{
														ARG:  collect__parent__type__node,
														arg1: &dNode{},
													},
												},
											},
										},
										arg2: &dSort{
											arg1: &dOr{
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
														DATA: []interface{}{"vc"},
														arg1: &dCollect{
															ARG: collect__attributes__rel,
															arg1: &dCollect{
																ARG:  collect__parent__type__node,
																arg1: &dNode{},
															},
														},
													},
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
																DATA: []interface{}{"cnj"},
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
																DATA: []interface{}{"vc"},
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
	"ppinp": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dAnd{
						arg1: &dAnd{
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
							arg2: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG: collect__attributes__begin,
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
								arg2: &dCollect{
									ARG:  collect__attributes__begin,
									arg1: &dNode{},
								},
							},
						},
						arg2: &dEqual{
							ARG: equal__is,
							arg1: &dCollect{
								ARG: collect__attributes__end,
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
								ARG:  collect__attributes__end,
								arg1: &dNode{},
							},
						},
					},
				},
			},
		},
	},
	"ppirp": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dAnd{
						arg1: &dAnd{
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
							arg2: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG: collect__attributes__end,
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
								arg2: &dCollect{
									ARG:  collect__attributes__end,
									arg1: &dNode{},
								},
							},
						},
						arg2: &dEqual{
							ARG: equal__is,
							arg1: &dCollect{
								ARG: collect__attributes__begin,
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
											arg2: &dEqual{
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
										},
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
	"ppimwu": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
								DATA: []interface{}{"pp"},
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
	"vwuit": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
								DATA: []interface{}{"svp"},
								arg1: &dCollect{
									ARG:  collect__attributes__rel,
									arg1: &dNode{},
								},
							},
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
	"groen": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
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
										ARG:  collect__attributes__wvorm,
										arg1: &dNode{},
									},
									arg2: &dElem{
										DATA: []interface{}{"vd"},
										arg1: &dCollect{
											ARG:  collect__attributes__wvorm,
											arg1: &dNode{},
										},
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
									arg1: &dCollect{
										ARG: collect__child__node,
										arg1: &dCollect{
											ARG: collect__parent__type__node,
											arg1: &dCollect{
												ARG:  collect__parent__node,
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
											DATA: []interface{}{"smain", "sv1"},
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
	},
	"rood": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
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
										ARG:  collect__attributes__wvorm,
										arg1: &dNode{},
									},
									arg2: &dElem{
										DATA: []interface{}{"vd"},
										arg1: &dCollect{
											ARG:  collect__attributes__wvorm,
											arg1: &dNode{},
										},
									},
								},
							},
							arg2: &dCmp{
								ARG: cmp__gt,
								arg1: &dCollect{
									ARG:  collect__attributes__begin,
									arg1: &dNode{},
								},
								arg2: &dCollect{
									ARG: collect__attributes__begin,
									arg1: &dCollect{
										ARG: collect__child__node,
										arg1: &dCollect{
											ARG: collect__parent__type__node,
											arg1: &dCollect{
												ARG:  collect__parent__type__node,
												arg1: &dNode{},
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
														ARG:  collect__attributes__pt,
														arg1: &dNode{},
													},
													arg2: &dElem{
														DATA: []interface{}{"ww"},
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
											DATA: []interface{}{"smain", "sv1"},
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
	},
	"wwclus": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dSort{
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
										DATA: []interface{}{"ti", "inf", "ppart"},
										arg1: &dCollect{
											ARG:  collect__attributes__cat,
											arg1: &dNode{},
										},
									},
								},
							},
							arg2: &dCmp{
								ARG: cmp__lt,
								arg1: &dCollect{
									ARG: collect__attributes__begin,
									arg1: &dCollect{
										ARG:  collect__child__node,
										arg1: &dNode{},
									},
								},
								arg2: &dCollect{
									ARG: collect__attributes__begin,
									arg1: &dCollect{
										ARG: collect__child__node,
										arg1: &dCollect{
											ARG:  collect__parent__type__node,
											arg1: &dNode{},
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
													arg2: &dEqual{
														ARG: equal__is,
														arg1: &dCollect{
															ARG:  collect__attributes__pt,
															arg1: &dNode{},
														},
														arg2: &dElem{
															DATA: []interface{}{"ww"},
															arg1: &dCollect{
																ARG:  collect__attributes__pt,
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
																		ARG:  collect__parent__type__node,
																		arg1: &dNode{},
																	},
																},
																arg2: &dElem{
																	DATA: []interface{}{"smain", "sv1"},
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
	"accinf": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dAnd{
						arg1: &dAnd{
							arg1: &dSort{
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
												DATA: []interface{}{"ti", "inf", "ppart"},
												arg1: &dCollect{
													ARG:  collect__attributes__cat,
													arg1: &dNode{},
												},
											},
										},
									},
									arg2: &dCmp{
										ARG: cmp__lt,
										arg1: &dCollect{
											ARG: collect__attributes__begin,
											arg1: &dCollect{
												ARG:  collect__child__node,
												arg1: &dNode{},
											},
										},
										arg2: &dCollect{
											ARG: collect__attributes__begin,
											arg1: &dCollect{
												ARG: collect__child__node,
												arg1: &dCollect{
													ARG:  collect__parent__type__node,
													arg1: &dNode{},
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
															arg2: &dEqual{
																ARG: equal__is,
																arg1: &dCollect{
																	ARG:  collect__attributes__pt,
																	arg1: &dNode{},
																},
																arg2: &dElem{
																	DATA: []interface{}{"ww"},
																	arg1: &dCollect{
																		ARG:  collect__attributes__pt,
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
																				ARG:  collect__parent__type__node,
																				arg1: &dNode{},
																			},
																		},
																		arg2: &dElem{
																			DATA: []interface{}{"smain", "sv1"},
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
												},
											},
										},
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
									DATA: []interface{}{"inf"},
									arg1: &dCollect{
										ARG:  collect__attributes__cat,
										arg1: &dNode{},
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
	"passive": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dSort{
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
									arg1: &dSort{
										arg1: &dAnd{
											arg1: &dEqual{
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
															arg1: &dCollect{
																ARG:  collect__parent__type__node,
																arg1: &dNode{},
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
									arg2: &dSort{
										arg1: &dAnd{
											arg1: &dEqual{
												ARG: equal__is,
												arg1: &dCollect{
													ARG:  collect__attributes__cat,
													arg1: &dNode{},
												},
												arg2: &dElem{
													DATA: []interface{}{"ti"},
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
												arg2: &dCollect{
													ARG: collect__attributes__index,
													arg1: &dCollect{
														ARG: collect__child__node,
														arg1: &dCollect{
															ARG: collect__parent__type__node,
															arg1: &dCollect{
																ARG:  collect__parent__type__node,
																arg1: &dNode{},
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
				},
			},
		},
	},
	"nppas": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dSort{
						arg1: &dAnd{
							arg1: &dSort{
								arg1: &dOr{
									arg1: &dSort{
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
																	DATA: []interface{}{"ww"},
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
											arg2: &dSort{
												arg1: &dOr{
													arg1: &dSort{
														arg1: &dAnd{
															arg1: &dEqual{
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
															arg2: &dFunction{
																ARG: function__not__1__args,
																arg1: &dArg{
																	arg1: &dSort{
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
																						DATA: []interface{}{"obj1", "su", "vc", "predc"},
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
													arg2: &dSort{
														arg1: &dAnd{
															arg1: &dEqual{
																ARG: equal__is,
																arg1: &dCollect{
																	ARG:  collect__attributes__cat,
																	arg1: &dNode{},
																},
																arg2: &dElem{
																	DATA: []interface{}{"ti"},
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
																							DATA: []interface{}{"body"},
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
																						DATA: []interface{}{"obj1", "su", "vc", "predc"},
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
																ARG:  collect__attributes__sc,
																arg1: &dNode{},
															},
															arg2: &dElem{
																DATA: []interface{}{"passive", "te_passive"},
																arg1: &dCollect{
																	ARG:  collect__attributes__sc,
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
								ARG:  collect__parent__node,
								arg1: &dNode{},
								arg2: &dPredicate{
									arg1: &dAnd{
										arg1: &dFunction{
											ARG: function__not__1__args,
											arg1: &dArg{
												arg1: &dSort{
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
										arg2: &dSort{
											arg1: &dOr{
												arg1: &dEqual{
													ARG: equal__is,
													arg1: &dCollect{
														ARG:  collect__attributes__cat,
														arg1: &dNode{},
													},
													arg2: &dElem{
														DATA: []interface{}{"smain", "ssub"},
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
																DATA: []interface{}{"sv1"},
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
																	arg1: &dSort{
																		arg1: &dAnd{
																			arg1: &dAnd{
																				arg1: &dAnd{
																					arg1: &dEqual{
																						ARG: equal__is,
																						arg1: &dCollect{
																							ARG:  collect__attributes__cat,
																							arg1: &dNode{},
																						},
																						arg2: &dElem{
																							DATA: []interface{}{"sv1"},
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
																				arg2: &dFunction{
																					ARG: function__not__1__args,
																					arg1: &dArg{
																						arg1: &dSort{
																							arg1: &dCollect{
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
																											arg2: &dCollect{
																												ARG:  collect__attributes__stype,
																												arg1: &dNode{},
																											},
																										},
																										arg2: &dFunction{
																											ARG: function__not__1__args,
																											arg1: &dArg{
																												arg1: &dSort{
																													arg1: &dEqual{
																														ARG: equal__is,
																														arg1: &dCollect{
																															ARG:  collect__attributes__stype,
																															arg1: &dNode{},
																														},
																														arg2: &dElem{
																															DATA: []interface{}{"imparative"},
																															arg1: &dCollect{
																																ARG:  collect__attributes__stype,
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
																											arg1: &dOr{
																												arg1: &dOr{
																													arg1: &dEqual{
																														ARG: equal__is,
																														arg1: &dCollect{
																															ARG:  collect__attributes__tense,
																															arg1: &dNode{},
																														},
																														arg2: &dElem{
																															DATA: []interface{}{"past"},
																															arg1: &dCollect{
																																ARG:  collect__attributes__tense,
																																arg1: &dNode{},
																															},
																														},
																													},
																													arg2: &dEqual{
																														ARG: equal__is,
																														arg1: &dCollect{
																															ARG:  collect__attributes__pvagr,
																															arg1: &dNode{},
																														},
																														arg2: &dElem{
																															DATA: []interface{}{"mv", "met-t"},
																															arg1: &dCollect{
																																ARG:  collect__attributes__pvagr,
																																arg1: &dNode{},
																															},
																														},
																													},
																												},
																												arg2: &dEqual{
																													ARG: equal__is,
																													arg1: &dCollect{
																														ARG:  collect__attributes__pvtijd,
																														arg1: &dNode{},
																													},
																													arg2: &dElem{
																														DATA: []interface{}{"verl"},
																														arg1: &dCollect{
																															ARG:  collect__attributes__pvtijd,
																															arg1: &dNode{},
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
																													DATA: []interface{}{"zijn", "kunnen", "willen", "moeten", "mogen", "zullen", "denken", "geloven", "vinden", "hebben"},
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
	"vpart": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
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
								DATA: []interface{}{"part_"},
							},
						},
					},
				},
			},
		},
	},
	"vpartex": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dAnd{
						arg1: &dAnd{
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
										DATA: []interface{}{"part_"},
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
											DATA: []interface{}{"svp"},
											arg1: &dCollect{
												ARG:  collect__attributes__rel,
												arg1: &dNode{},
											},
										},
									},
									arg2: &dFunction{
										ARG: function__starts__with__2__args,
										arg1: &dArg{
											arg1: &dArg{
												arg1: &dSort{
													arg1: &dCollect{
														ARG:  collect__attributes__frame,
														arg1: &dNode{},
													},
												},
											},
											arg2: &dElem{
												DATA: []interface{}{"particle"},
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
	"vpartin": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dAnd{
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
									DATA: []interface{}{"part_"},
								},
							},
						},
						arg2: &dSort{
							arg1: &dOr{
								arg1: &dFunction{
									ARG: function__not__1__args,
									arg1: &dArg{
										arg1: &dSort{
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
								arg2: &dAnd{
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
																	DATA: []interface{}{"svp"},
																	arg1: &dCollect{
																		ARG:  collect__attributes__rel,
																		arg1: &dNode{},
																	},
																},
															},
															arg2: &dFunction{
																ARG: function__starts__with__2__args,
																arg1: &dArg{
																	arg1: &dArg{
																		arg1: &dSort{
																			arg1: &dCollect{
																				ARG:  collect__attributes__frame,
																				arg1: &dNode{},
																			},
																		},
																	},
																	arg2: &dElem{
																		DATA: []interface{}{"particle"},
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
	"vprtn": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dAnd{
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
									DATA: []interface{}{"part_"},
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
											ARG:  collect__attributes__wvorm,
											arg1: &dNode{},
										},
										arg2: &dElem{
											DATA: []interface{}{"pv"},
											arg1: &dCollect{
												ARG:  collect__attributes__wvorm,
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
	"vprtnex": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dAnd{
						arg1: &dAnd{
							arg1: &dAnd{
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
											DATA: []interface{}{"part_"},
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
													ARG:  collect__attributes__wvorm,
													arg1: &dNode{},
												},
												arg2: &dElem{
													DATA: []interface{}{"pv"},
													arg1: &dCollect{
														ARG:  collect__attributes__wvorm,
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
											DATA: []interface{}{"svp"},
											arg1: &dCollect{
												ARG:  collect__attributes__rel,
												arg1: &dNode{},
											},
										},
									},
									arg2: &dFunction{
										ARG: function__starts__with__2__args,
										arg1: &dArg{
											arg1: &dArg{
												arg1: &dSort{
													arg1: &dCollect{
														ARG:  collect__attributes__frame,
														arg1: &dNode{},
													},
												},
											},
											arg2: &dElem{
												DATA: []interface{}{"particle"},
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
	"vprtnin": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dAnd{
						arg1: &dAnd{
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
										DATA: []interface{}{"part_"},
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
												ARG:  collect__attributes__wvorm,
												arg1: &dNode{},
											},
											arg2: &dElem{
												DATA: []interface{}{"pv"},
												arg1: &dCollect{
													ARG:  collect__attributes__wvorm,
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
								arg1: &dFunction{
									ARG: function__not__1__args,
									arg1: &dArg{
										arg1: &dSort{
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
								arg2: &dAnd{
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
																	DATA: []interface{}{"svp"},
																	arg1: &dCollect{
																		ARG:  collect__attributes__rel,
																		arg1: &dNode{},
																	},
																},
															},
															arg2: &dFunction{
																ARG: function__starts__with__2__args,
																arg1: &dArg{
																	arg1: &dArg{
																		arg1: &dSort{
																			arg1: &dCollect{
																				ARG:  collect__attributes__frame,
																				arg1: &dNode{},
																			},
																		},
																	},
																	arg2: &dElem{
																		DATA: []interface{}{"particle"},
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
	"inb0": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dSort{
						arg1: &dAnd{
							arg1: &dAnd{
								arg1: &dEqual{
									ARG: equal__is,
									arg1: &dCollect{
										ARG:  collect__attributes__cat,
										arg1: &dNode{},
									},
									arg2: &dElem{
										DATA: []interface{}{"smain", "sv1", "ssub"},
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
															DATA: []interface{}{"smain", "sv1", "ssub"},
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
							arg2: &dFunction{
								ARG: function__not__1__args,
								arg1: &dArg{
									arg1: &dSort{
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
														DATA: []interface{}{"smain", "sv1", "ssub"},
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
	"inb1": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dSort{
						arg1: &dAnd{
							arg1: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG:  collect__attributes__cat,
									arg1: &dNode{},
								},
								arg2: &dElem{
									DATA: []interface{}{"smain", "sv1", "ssub"},
									arg1: &dCollect{
										ARG:  collect__attributes__cat,
										arg1: &dNode{},
									},
								},
							},
							arg2: &dCollect{
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
											DATA: []interface{}{"smain", "sv1", "ssub"},
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
	"inb2": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dSort{
						arg1: &dAnd{
							arg1: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG:  collect__attributes__cat,
									arg1: &dNode{},
								},
								arg2: &dElem{
									DATA: []interface{}{"smain", "sv1", "ssub"},
									arg1: &dCollect{
										ARG:  collect__attributes__cat,
										arg1: &dNode{},
									},
								},
							},
							arg2: &dCollect{
								ARG: collect__child__node,
								arg1: &dCollect{
									ARG:  collect__descendant__or__self__type__node,
									arg1: &dNode{},
								},
								arg2: &dPredicate{
									arg1: &dSort{
										arg1: &dAnd{
											arg1: &dEqual{
												ARG: equal__is,
												arg1: &dCollect{
													ARG:  collect__attributes__cat,
													arg1: &dNode{},
												},
												arg2: &dElem{
													DATA: []interface{}{"smain", "sv1", "ssub"},
													arg1: &dCollect{
														ARG:  collect__attributes__cat,
														arg1: &dNode{},
													},
												},
											},
											arg2: &dCollect{
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
															DATA: []interface{}{"smain", "sv1", "ssub"},
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
	"inb3": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dSort{
						arg1: &dAnd{
							arg1: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG:  collect__attributes__cat,
									arg1: &dNode{},
								},
								arg2: &dElem{
									DATA: []interface{}{"smain", "sv1", "ssub"},
									arg1: &dCollect{
										ARG:  collect__attributes__cat,
										arg1: &dNode{},
									},
								},
							},
							arg2: &dCollect{
								ARG: collect__child__node,
								arg1: &dCollect{
									ARG:  collect__descendant__or__self__type__node,
									arg1: &dNode{},
								},
								arg2: &dPredicate{
									arg1: &dSort{
										arg1: &dAnd{
											arg1: &dEqual{
												ARG: equal__is,
												arg1: &dCollect{
													ARG:  collect__attributes__cat,
													arg1: &dNode{},
												},
												arg2: &dElem{
													DATA: []interface{}{"smain", "sv1", "ssub"},
													arg1: &dCollect{
														ARG:  collect__attributes__cat,
														arg1: &dNode{},
													},
												},
											},
											arg2: &dCollect{
												ARG: collect__child__node,
												arg1: &dCollect{
													ARG:  collect__descendant__or__self__type__node,
													arg1: &dNode{},
												},
												arg2: &dPredicate{
													arg1: &dSort{
														arg1: &dAnd{
															arg1: &dEqual{
																ARG: equal__is,
																arg1: &dCollect{
																	ARG:  collect__attributes__cat,
																	arg1: &dNode{},
																},
																arg2: &dElem{
																	DATA: []interface{}{"smain", "sv1", "ssub"},
																	arg1: &dCollect{
																		ARG:  collect__attributes__cat,
																		arg1: &dNode{},
																	},
																},
															},
															arg2: &dCollect{
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
																			DATA: []interface{}{"smain", "sv1", "ssub"},
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
			},
		},
	},
	"inb4": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dSort{
						arg1: &dAnd{
							arg1: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG:  collect__attributes__cat,
									arg1: &dNode{},
								},
								arg2: &dElem{
									DATA: []interface{}{"smain", "sv1", "ssub"},
									arg1: &dCollect{
										ARG:  collect__attributes__cat,
										arg1: &dNode{},
									},
								},
							},
							arg2: &dCollect{
								ARG: collect__child__node,
								arg1: &dCollect{
									ARG:  collect__descendant__or__self__type__node,
									arg1: &dNode{},
								},
								arg2: &dPredicate{
									arg1: &dSort{
										arg1: &dAnd{
											arg1: &dEqual{
												ARG: equal__is,
												arg1: &dCollect{
													ARG:  collect__attributes__cat,
													arg1: &dNode{},
												},
												arg2: &dElem{
													DATA: []interface{}{"smain", "sv1", "ssub"},
													arg1: &dCollect{
														ARG:  collect__attributes__cat,
														arg1: &dNode{},
													},
												},
											},
											arg2: &dCollect{
												ARG: collect__child__node,
												arg1: &dCollect{
													ARG:  collect__descendant__or__self__type__node,
													arg1: &dNode{},
												},
												arg2: &dPredicate{
													arg1: &dSort{
														arg1: &dAnd{
															arg1: &dEqual{
																ARG: equal__is,
																arg1: &dCollect{
																	ARG:  collect__attributes__cat,
																	arg1: &dNode{},
																},
																arg2: &dElem{
																	DATA: []interface{}{"smain", "sv1", "ssub"},
																	arg1: &dCollect{
																		ARG:  collect__attributes__cat,
																		arg1: &dNode{},
																	},
																},
															},
															arg2: &dCollect{
																ARG: collect__child__node,
																arg1: &dCollect{
																	ARG:  collect__descendant__or__self__type__node,
																	arg1: &dNode{},
																},
																arg2: &dPredicate{
																	arg1: &dSort{
																		arg1: &dAnd{
																			arg1: &dEqual{
																				ARG: equal__is,
																				arg1: &dCollect{
																					ARG:  collect__attributes__cat,
																					arg1: &dNode{},
																				},
																				arg2: &dElem{
																					DATA: []interface{}{"smain", "sv1", "ssub"},
																					arg1: &dCollect{
																						ARG:  collect__attributes__cat,
																						arg1: &dNode{},
																					},
																				},
																			},
																			arg2: &dCollect{
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
																							DATA: []interface{}{"smain", "sv1", "ssub"},
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
							},
						},
					},
				},
			},
		},
	},
	"inb5": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dSort{
						arg1: &dAnd{
							arg1: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG:  collect__attributes__cat,
									arg1: &dNode{},
								},
								arg2: &dElem{
									DATA: []interface{}{"smain", "sv1", "ssub"},
									arg1: &dCollect{
										ARG:  collect__attributes__cat,
										arg1: &dNode{},
									},
								},
							},
							arg2: &dCollect{
								ARG: collect__child__node,
								arg1: &dCollect{
									ARG:  collect__descendant__or__self__type__node,
									arg1: &dNode{},
								},
								arg2: &dPredicate{
									arg1: &dSort{
										arg1: &dAnd{
											arg1: &dEqual{
												ARG: equal__is,
												arg1: &dCollect{
													ARG:  collect__attributes__cat,
													arg1: &dNode{},
												},
												arg2: &dElem{
													DATA: []interface{}{"smain", "sv1", "ssub"},
													arg1: &dCollect{
														ARG:  collect__attributes__cat,
														arg1: &dNode{},
													},
												},
											},
											arg2: &dCollect{
												ARG: collect__child__node,
												arg1: &dCollect{
													ARG:  collect__descendant__or__self__type__node,
													arg1: &dNode{},
												},
												arg2: &dPredicate{
													arg1: &dSort{
														arg1: &dAnd{
															arg1: &dEqual{
																ARG: equal__is,
																arg1: &dCollect{
																	ARG:  collect__attributes__cat,
																	arg1: &dNode{},
																},
																arg2: &dElem{
																	DATA: []interface{}{"smain", "sv1", "ssub"},
																	arg1: &dCollect{
																		ARG:  collect__attributes__cat,
																		arg1: &dNode{},
																	},
																},
															},
															arg2: &dCollect{
																ARG: collect__child__node,
																arg1: &dCollect{
																	ARG:  collect__descendant__or__self__type__node,
																	arg1: &dNode{},
																},
																arg2: &dPredicate{
																	arg1: &dSort{
																		arg1: &dAnd{
																			arg1: &dEqual{
																				ARG: equal__is,
																				arg1: &dCollect{
																					ARG:  collect__attributes__cat,
																					arg1: &dNode{},
																				},
																				arg2: &dElem{
																					DATA: []interface{}{"smain", "sv1", "ssub"},
																					arg1: &dCollect{
																						ARG:  collect__attributes__cat,
																						arg1: &dNode{},
																					},
																				},
																			},
																			arg2: &dCollect{
																				ARG: collect__child__node,
																				arg1: &dCollect{
																					ARG:  collect__descendant__or__self__type__node,
																					arg1: &dNode{},
																				},
																				arg2: &dPredicate{
																					arg1: &dSort{
																						arg1: &dAnd{
																							arg1: &dEqual{
																								ARG: equal__is,
																								arg1: &dCollect{
																									ARG:  collect__attributes__cat,
																									arg1: &dNode{},
																								},
																								arg2: &dElem{
																									DATA: []interface{}{"smain", "sv1", "ssub"},
																									arg1: &dCollect{
																										ARG:  collect__attributes__cat,
																										arg1: &dNode{},
																									},
																								},
																							},
																							arg2: &dCollect{
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
																											DATA: []interface{}{"smain", "sv1", "ssub"},
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
	"inb6": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dSort{
						arg1: &dAnd{
							arg1: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG:  collect__attributes__cat,
									arg1: &dNode{},
								},
								arg2: &dElem{
									DATA: []interface{}{"smain", "sv1", "ssub"},
									arg1: &dCollect{
										ARG:  collect__attributes__cat,
										arg1: &dNode{},
									},
								},
							},
							arg2: &dCollect{
								ARG: collect__child__node,
								arg1: &dCollect{
									ARG:  collect__descendant__or__self__type__node,
									arg1: &dNode{},
								},
								arg2: &dPredicate{
									arg1: &dSort{
										arg1: &dAnd{
											arg1: &dEqual{
												ARG: equal__is,
												arg1: &dCollect{
													ARG:  collect__attributes__cat,
													arg1: &dNode{},
												},
												arg2: &dElem{
													DATA: []interface{}{"smain", "sv1", "ssub"},
													arg1: &dCollect{
														ARG:  collect__attributes__cat,
														arg1: &dNode{},
													},
												},
											},
											arg2: &dCollect{
												ARG: collect__child__node,
												arg1: &dCollect{
													ARG:  collect__descendant__or__self__type__node,
													arg1: &dNode{},
												},
												arg2: &dPredicate{
													arg1: &dSort{
														arg1: &dAnd{
															arg1: &dEqual{
																ARG: equal__is,
																arg1: &dCollect{
																	ARG:  collect__attributes__cat,
																	arg1: &dNode{},
																},
																arg2: &dElem{
																	DATA: []interface{}{"smain", "sv1", "ssub"},
																	arg1: &dCollect{
																		ARG:  collect__attributes__cat,
																		arg1: &dNode{},
																	},
																},
															},
															arg2: &dCollect{
																ARG: collect__child__node,
																arg1: &dCollect{
																	ARG:  collect__descendant__or__self__type__node,
																	arg1: &dNode{},
																},
																arg2: &dPredicate{
																	arg1: &dSort{
																		arg1: &dAnd{
																			arg1: &dEqual{
																				ARG: equal__is,
																				arg1: &dCollect{
																					ARG:  collect__attributes__cat,
																					arg1: &dNode{},
																				},
																				arg2: &dElem{
																					DATA: []interface{}{"smain", "sv1", "ssub"},
																					arg1: &dCollect{
																						ARG:  collect__attributes__cat,
																						arg1: &dNode{},
																					},
																				},
																			},
																			arg2: &dCollect{
																				ARG: collect__child__node,
																				arg1: &dCollect{
																					ARG:  collect__descendant__or__self__type__node,
																					arg1: &dNode{},
																				},
																				arg2: &dPredicate{
																					arg1: &dSort{
																						arg1: &dAnd{
																							arg1: &dEqual{
																								ARG: equal__is,
																								arg1: &dCollect{
																									ARG:  collect__attributes__cat,
																									arg1: &dNode{},
																								},
																								arg2: &dElem{
																									DATA: []interface{}{"smain", "sv1", "ssub"},
																									arg1: &dCollect{
																										ARG:  collect__attributes__cat,
																										arg1: &dNode{},
																									},
																								},
																							},
																							arg2: &dCollect{
																								ARG: collect__child__node,
																								arg1: &dCollect{
																									ARG:  collect__descendant__or__self__type__node,
																									arg1: &dNode{},
																								},
																								arg2: &dPredicate{
																									arg1: &dSort{
																										arg1: &dAnd{
																											arg1: &dEqual{
																												ARG: equal__is,
																												arg1: &dCollect{
																													ARG:  collect__attributes__cat,
																													arg1: &dNode{},
																												},
																												arg2: &dElem{
																													DATA: []interface{}{"smain", "sv1", "ssub"},
																													arg1: &dCollect{
																														ARG:  collect__attributes__cat,
																														arg1: &dNode{},
																													},
																												},
																											},
																											arg2: &dCollect{
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
																															DATA: []interface{}{"smain", "sv1", "ssub"},
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
	"inb7": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dSort{
						arg1: &dAnd{
							arg1: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG:  collect__attributes__cat,
									arg1: &dNode{},
								},
								arg2: &dElem{
									DATA: []interface{}{"smain", "sv1", "ssub"},
									arg1: &dCollect{
										ARG:  collect__attributes__cat,
										arg1: &dNode{},
									},
								},
							},
							arg2: &dCollect{
								ARG: collect__child__node,
								arg1: &dCollect{
									ARG:  collect__descendant__or__self__type__node,
									arg1: &dNode{},
								},
								arg2: &dPredicate{
									arg1: &dSort{
										arg1: &dAnd{
											arg1: &dEqual{
												ARG: equal__is,
												arg1: &dCollect{
													ARG:  collect__attributes__cat,
													arg1: &dNode{},
												},
												arg2: &dElem{
													DATA: []interface{}{"smain", "sv1", "ssub"},
													arg1: &dCollect{
														ARG:  collect__attributes__cat,
														arg1: &dNode{},
													},
												},
											},
											arg2: &dCollect{
												ARG: collect__child__node,
												arg1: &dCollect{
													ARG:  collect__descendant__or__self__type__node,
													arg1: &dNode{},
												},
												arg2: &dPredicate{
													arg1: &dSort{
														arg1: &dAnd{
															arg1: &dEqual{
																ARG: equal__is,
																arg1: &dCollect{
																	ARG:  collect__attributes__cat,
																	arg1: &dNode{},
																},
																arg2: &dElem{
																	DATA: []interface{}{"smain", "sv1", "ssub"},
																	arg1: &dCollect{
																		ARG:  collect__attributes__cat,
																		arg1: &dNode{},
																	},
																},
															},
															arg2: &dCollect{
																ARG: collect__child__node,
																arg1: &dCollect{
																	ARG:  collect__descendant__or__self__type__node,
																	arg1: &dNode{},
																},
																arg2: &dPredicate{
																	arg1: &dSort{
																		arg1: &dAnd{
																			arg1: &dEqual{
																				ARG: equal__is,
																				arg1: &dCollect{
																					ARG:  collect__attributes__cat,
																					arg1: &dNode{},
																				},
																				arg2: &dElem{
																					DATA: []interface{}{"smain", "sv1", "ssub"},
																					arg1: &dCollect{
																						ARG:  collect__attributes__cat,
																						arg1: &dNode{},
																					},
																				},
																			},
																			arg2: &dCollect{
																				ARG: collect__child__node,
																				arg1: &dCollect{
																					ARG:  collect__descendant__or__self__type__node,
																					arg1: &dNode{},
																				},
																				arg2: &dPredicate{
																					arg1: &dSort{
																						arg1: &dAnd{
																							arg1: &dEqual{
																								ARG: equal__is,
																								arg1: &dCollect{
																									ARG:  collect__attributes__cat,
																									arg1: &dNode{},
																								},
																								arg2: &dElem{
																									DATA: []interface{}{"smain", "sv1", "ssub"},
																									arg1: &dCollect{
																										ARG:  collect__attributes__cat,
																										arg1: &dNode{},
																									},
																								},
																							},
																							arg2: &dCollect{
																								ARG: collect__child__node,
																								arg1: &dCollect{
																									ARG:  collect__descendant__or__self__type__node,
																									arg1: &dNode{},
																								},
																								arg2: &dPredicate{
																									arg1: &dSort{
																										arg1: &dAnd{
																											arg1: &dEqual{
																												ARG: equal__is,
																												arg1: &dCollect{
																													ARG:  collect__attributes__cat,
																													arg1: &dNode{},
																												},
																												arg2: &dElem{
																													DATA: []interface{}{"smain", "sv1", "ssub"},
																													arg1: &dCollect{
																														ARG:  collect__attributes__cat,
																														arg1: &dNode{},
																													},
																												},
																											},
																											arg2: &dCollect{
																												ARG: collect__child__node,
																												arg1: &dCollect{
																													ARG:  collect__descendant__or__self__type__node,
																													arg1: &dNode{},
																												},
																												arg2: &dPredicate{
																													arg1: &dSort{
																														arg1: &dAnd{
																															arg1: &dEqual{
																																ARG: equal__is,
																																arg1: &dCollect{
																																	ARG:  collect__attributes__cat,
																																	arg1: &dNode{},
																																},
																																arg2: &dElem{
																																	DATA: []interface{}{"smain", "sv1", "ssub"},
																																	arg1: &dCollect{
																																		ARG:  collect__attributes__cat,
																																		arg1: &dNode{},
																																	},
																																},
																															},
																															arg2: &dCollect{
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
																																			DATA: []interface{}{"smain", "sv1", "ssub"},
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
	},
	"inb8": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dSort{
						arg1: &dAnd{
							arg1: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG:  collect__attributes__cat,
									arg1: &dNode{},
								},
								arg2: &dElem{
									DATA: []interface{}{"smain", "sv1", "ssub"},
									arg1: &dCollect{
										ARG:  collect__attributes__cat,
										arg1: &dNode{},
									},
								},
							},
							arg2: &dCollect{
								ARG: collect__child__node,
								arg1: &dCollect{
									ARG:  collect__descendant__or__self__type__node,
									arg1: &dNode{},
								},
								arg2: &dPredicate{
									arg1: &dSort{
										arg1: &dAnd{
											arg1: &dEqual{
												ARG: equal__is,
												arg1: &dCollect{
													ARG:  collect__attributes__cat,
													arg1: &dNode{},
												},
												arg2: &dElem{
													DATA: []interface{}{"smain", "sv1", "ssub"},
													arg1: &dCollect{
														ARG:  collect__attributes__cat,
														arg1: &dNode{},
													},
												},
											},
											arg2: &dCollect{
												ARG: collect__child__node,
												arg1: &dCollect{
													ARG:  collect__descendant__or__self__type__node,
													arg1: &dNode{},
												},
												arg2: &dPredicate{
													arg1: &dSort{
														arg1: &dAnd{
															arg1: &dEqual{
																ARG: equal__is,
																arg1: &dCollect{
																	ARG:  collect__attributes__cat,
																	arg1: &dNode{},
																},
																arg2: &dElem{
																	DATA: []interface{}{"smain", "sv1", "ssub"},
																	arg1: &dCollect{
																		ARG:  collect__attributes__cat,
																		arg1: &dNode{},
																	},
																},
															},
															arg2: &dCollect{
																ARG: collect__child__node,
																arg1: &dCollect{
																	ARG:  collect__descendant__or__self__type__node,
																	arg1: &dNode{},
																},
																arg2: &dPredicate{
																	arg1: &dSort{
																		arg1: &dAnd{
																			arg1: &dEqual{
																				ARG: equal__is,
																				arg1: &dCollect{
																					ARG:  collect__attributes__cat,
																					arg1: &dNode{},
																				},
																				arg2: &dElem{
																					DATA: []interface{}{"smain", "sv1", "ssub"},
																					arg1: &dCollect{
																						ARG:  collect__attributes__cat,
																						arg1: &dNode{},
																					},
																				},
																			},
																			arg2: &dCollect{
																				ARG: collect__child__node,
																				arg1: &dCollect{
																					ARG:  collect__descendant__or__self__type__node,
																					arg1: &dNode{},
																				},
																				arg2: &dPredicate{
																					arg1: &dSort{
																						arg1: &dAnd{
																							arg1: &dEqual{
																								ARG: equal__is,
																								arg1: &dCollect{
																									ARG:  collect__attributes__cat,
																									arg1: &dNode{},
																								},
																								arg2: &dElem{
																									DATA: []interface{}{"smain", "sv1", "ssub"},
																									arg1: &dCollect{
																										ARG:  collect__attributes__cat,
																										arg1: &dNode{},
																									},
																								},
																							},
																							arg2: &dCollect{
																								ARG: collect__child__node,
																								arg1: &dCollect{
																									ARG:  collect__descendant__or__self__type__node,
																									arg1: &dNode{},
																								},
																								arg2: &dPredicate{
																									arg1: &dSort{
																										arg1: &dAnd{
																											arg1: &dEqual{
																												ARG: equal__is,
																												arg1: &dCollect{
																													ARG:  collect__attributes__cat,
																													arg1: &dNode{},
																												},
																												arg2: &dElem{
																													DATA: []interface{}{"smain", "sv1", "ssub"},
																													arg1: &dCollect{
																														ARG:  collect__attributes__cat,
																														arg1: &dNode{},
																													},
																												},
																											},
																											arg2: &dCollect{
																												ARG: collect__child__node,
																												arg1: &dCollect{
																													ARG:  collect__descendant__or__self__type__node,
																													arg1: &dNode{},
																												},
																												arg2: &dPredicate{
																													arg1: &dSort{
																														arg1: &dAnd{
																															arg1: &dEqual{
																																ARG: equal__is,
																																arg1: &dCollect{
																																	ARG:  collect__attributes__cat,
																																	arg1: &dNode{},
																																},
																																arg2: &dElem{
																																	DATA: []interface{}{"smain", "sv1", "ssub"},
																																	arg1: &dCollect{
																																		ARG:  collect__attributes__cat,
																																		arg1: &dNode{},
																																	},
																																},
																															},
																															arg2: &dCollect{
																																ARG: collect__child__node,
																																arg1: &dCollect{
																																	ARG:  collect__descendant__or__self__type__node,
																																	arg1: &dNode{},
																																},
																																arg2: &dPredicate{
																																	arg1: &dSort{
																																		arg1: &dAnd{
																																			arg1: &dEqual{
																																				ARG: equal__is,
																																				arg1: &dCollect{
																																					ARG:  collect__attributes__cat,
																																					arg1: &dNode{},
																																				},
																																				arg2: &dElem{
																																					DATA: []interface{}{"smain", "sv1", "ssub"},
																																					arg1: &dCollect{
																																						ARG:  collect__attributes__cat,
																																						arg1: &dNode{},
																																					},
																																				},
																																			},
																																			arg2: &dCollect{
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
																																							DATA: []interface{}{"smain", "sv1", "ssub"},
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
					},
				},
			},
		},
	},
	"nptsub": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dAnd{
						arg1: &dAnd{
							arg1: &dSort{
								arg1: &dAnd{
									arg1: &dSort{
										arg1: &dOr{
											arg1: &dSort{
												arg1: &dAnd{
													arg1: &dSort{
														arg1: &dOr{
															arg1: &dCmp{
																ARG: cmp__gt,
																arg1: &dCollect{
																	ARG: collect__attributes__begin,
																	arg1: &dCollect{
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
																						DATA: []interface{}{"smain"},
																						arg1: &dCollect{
																							ARG:  collect__attributes__cat,
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
																arg2: &dCollect{
																	ARG: collect__attributes__begin,
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
																					DATA: []interface{}{"hd", "cmp", "mwp", "crd", "rhd", "whd", "nucl", "dp"},
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
																arg1: &dAnd{
																	arg1: &dCmp{
																		ARG: cmp__gt,
																		arg1: &dCollect{
																			ARG: collect__attributes__begin,
																			arg1: &dCollect{
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
																								DATA: []interface{}{"smain"},
																								arg1: &dCollect{
																									ARG:  collect__attributes__cat,
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
																		arg2: &dCollect{
																			ARG:  collect__attributes__begin,
																			arg1: &dNode{},
																		},
																	},
																	arg2: &dFunction{
																		ARG: function__not__1__args,
																		arg1: &dArg{
																			arg1: &dSort{
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
																								DATA: []interface{}{"hd", "cmp", "mwp", "crd", "rhd", "whd", "nucl", "dp"},
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
													arg2: &dFunction{
														ARG: function__not__1__args,
														arg1: &dArg{
															arg1: &dSort{
																arg1: &dCollect{
																	ARG:  collect__parent__node,
																	arg1: &dNode{},
																	arg2: &dPredicate{
																		arg1: &dSort{
																			arg1: &dOr{
																				arg1: &dCmp{
																					ARG: cmp__gt,
																					arg1: &dCollect{
																						ARG: collect__attributes__begin,
																						arg1: &dCollect{
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
																											DATA: []interface{}{"smain"},
																											arg1: &dCollect{
																												ARG:  collect__attributes__cat,
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
																					arg2: &dCollect{
																						ARG: collect__attributes__begin,
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
																										DATA: []interface{}{"hd", "cmp", "mwp", "crd", "rhd", "whd", "nucl", "dp"},
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
																					arg1: &dAnd{
																						arg1: &dCmp{
																							ARG: cmp__gt,
																							arg1: &dCollect{
																								ARG: collect__attributes__begin,
																								arg1: &dCollect{
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
																													DATA: []interface{}{"smain"},
																													arg1: &dCollect{
																														ARG:  collect__attributes__cat,
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
																							arg2: &dCollect{
																								ARG:  collect__attributes__begin,
																								arg1: &dNode{},
																							},
																						},
																						arg2: &dFunction{
																							ARG: function__not__1__args,
																							arg1: &dArg{
																								arg1: &dSort{
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
																													DATA: []interface{}{"hd", "cmp", "mwp", "crd", "rhd", "whd", "nucl", "dp"},
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
											},
											arg2: &dSort{
												arg1: &dAnd{
													arg1: &dSort{
														arg1: &dOr{
															arg1: &dCmp{
																ARG: cmp__gt,
																arg1: &dCollect{
																	ARG: collect__attributes__begin,
																	arg1: &dCollect{
																		ARG: collect__child__node,
																		arg1: &dCollect{
																			ARG: collect__child__node,
																			arg1: &dCollect{
																				ARG: collect__parent__node,
																				arg1: &dCollect{
																					ARG:  collect__ancestors__or__self__node,
																					arg1: &dNode{},
																					arg2: &dPredicate{
																						arg1: &dEqual{
																							ARG: equal__is,
																							arg1: &dCollect{
																								ARG:  collect__attributes__rel,
																								arg1: &dNode{},
																							},
																							arg2: &dElem{
																								DATA: []interface{}{"whd"},
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
																							ARG:  collect__attributes__cat,
																							arg1: &dNode{},
																						},
																						arg2: &dElem{
																							DATA: []interface{}{"whq"},
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
																							DATA: []interface{}{"body"},
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
																							DATA: []interface{}{"sv1"},
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
																arg2: &dCollect{
																	ARG: collect__attributes__begin,
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
																					DATA: []interface{}{"hd", "cmp", "mwp", "crd", "rhd", "whd", "nucl", "dp"},
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
																arg1: &dAnd{
																	arg1: &dCmp{
																		ARG: cmp__gt,
																		arg1: &dCollect{
																			ARG: collect__attributes__begin,
																			arg1: &dCollect{
																				ARG: collect__child__node,
																				arg1: &dCollect{
																					ARG: collect__child__node,
																					arg1: &dCollect{
																						ARG: collect__parent__node,
																						arg1: &dCollect{
																							ARG:  collect__ancestors__or__self__node,
																							arg1: &dNode{},
																							arg2: &dPredicate{
																								arg1: &dEqual{
																									ARG: equal__is,
																									arg1: &dCollect{
																										ARG:  collect__attributes__rel,
																										arg1: &dNode{},
																									},
																									arg2: &dElem{
																										DATA: []interface{}{"whd"},
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
																									ARG:  collect__attributes__cat,
																									arg1: &dNode{},
																								},
																								arg2: &dElem{
																									DATA: []interface{}{"whq"},
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
																									DATA: []interface{}{"body"},
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
																									DATA: []interface{}{"sv1"},
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
																		arg2: &dCollect{
																			ARG:  collect__attributes__begin,
																			arg1: &dNode{},
																		},
																	},
																	arg2: &dFunction{
																		ARG: function__not__1__args,
																		arg1: &dArg{
																			arg1: &dSort{
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
																								DATA: []interface{}{"hd", "cmp", "mwp", "crd", "rhd", "whd", "nucl", "dp"},
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
													arg2: &dFunction{
														ARG: function__not__1__args,
														arg1: &dArg{
															arg1: &dSort{
																arg1: &dCollect{
																	ARG:  collect__parent__node,
																	arg1: &dNode{},
																	arg2: &dPredicate{
																		arg1: &dSort{
																			arg1: &dOr{
																				arg1: &dCmp{
																					ARG: cmp__gt,
																					arg1: &dCollect{
																						ARG: collect__attributes__begin,
																						arg1: &dCollect{
																							ARG: collect__child__node,
																							arg1: &dCollect{
																								ARG: collect__child__node,
																								arg1: &dCollect{
																									ARG: collect__parent__node,
																									arg1: &dCollect{
																										ARG:  collect__ancestors__or__self__node,
																										arg1: &dNode{},
																										arg2: &dPredicate{
																											arg1: &dEqual{
																												ARG: equal__is,
																												arg1: &dCollect{
																													ARG:  collect__attributes__rel,
																													arg1: &dNode{},
																												},
																												arg2: &dElem{
																													DATA: []interface{}{"whd"},
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
																												ARG:  collect__attributes__cat,
																												arg1: &dNode{},
																											},
																											arg2: &dElem{
																												DATA: []interface{}{"whq"},
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
																												DATA: []interface{}{"body"},
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
																												DATA: []interface{}{"sv1"},
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
																					arg2: &dCollect{
																						ARG: collect__attributes__begin,
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
																										DATA: []interface{}{"hd", "cmp", "mwp", "crd", "rhd", "whd", "nucl", "dp"},
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
																					arg1: &dAnd{
																						arg1: &dCmp{
																							ARG: cmp__gt,
																							arg1: &dCollect{
																								ARG: collect__attributes__begin,
																								arg1: &dCollect{
																									ARG: collect__child__node,
																									arg1: &dCollect{
																										ARG: collect__child__node,
																										arg1: &dCollect{
																											ARG: collect__parent__node,
																											arg1: &dCollect{
																												ARG:  collect__ancestors__or__self__node,
																												arg1: &dNode{},
																												arg2: &dPredicate{
																													arg1: &dEqual{
																														ARG: equal__is,
																														arg1: &dCollect{
																															ARG:  collect__attributes__rel,
																															arg1: &dNode{},
																														},
																														arg2: &dElem{
																															DATA: []interface{}{"whd"},
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
																														ARG:  collect__attributes__cat,
																														arg1: &dNode{},
																													},
																													arg2: &dElem{
																														DATA: []interface{}{"whq"},
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
																														DATA: []interface{}{"body"},
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
																														DATA: []interface{}{"sv1"},
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
																							arg2: &dCollect{
																								ARG:  collect__attributes__begin,
																								arg1: &dNode{},
																							},
																						},
																						arg2: &dFunction{
																							ARG: function__not__1__args,
																							arg1: &dArg{
																								arg1: &dSort{
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
																													DATA: []interface{}{"hd", "cmp", "mwp", "crd", "rhd", "whd", "nucl", "dp"},
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
							arg2: &dEqual{
								ARG: equal__is,
								arg1: &dCollect{
									ARG:  collect__attributes__rel,
									arg1: &dNode{},
								},
								arg2: &dElem{
									DATA: []interface{}{"su", "sup"},
									arg1: &dCollect{
										ARG:  collect__attributes__rel,
										arg1: &dNode{},
									},
								},
							},
						},
						arg2: &dSort{
							arg1: &dOr{
								arg1: &dSort{
									arg1: &dOr{
										arg1: &dOr{
											arg1: &dOr{
												arg1: &dOr{
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
													arg2: &dSort{
														arg1: &dAnd{
															arg1: &dEqual{
																ARG: equal__is,
																arg1: &dCollect{
																	ARG:  collect__attributes__lcat,
																	arg1: &dNode{},
																},
																arg2: &dElem{
																	DATA: []interface{}{"np"},
																	arg1: &dCollect{
																		ARG:  collect__attributes__lcat,
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
																				ARG:  collect__attributes__rel,
																				arg1: &dNode{},
																			},
																			arg2: &dElem{
																				DATA: []interface{}{"hd", "mwp"},
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
													arg1: &dAnd{
														arg1: &dEqual{
															ARG: equal__is,
															arg1: &dCollect{
																ARG:  collect__attributes__pt,
																arg1: &dNode{},
															},
															arg2: &dElem{
																DATA: []interface{}{"n"},
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
											arg2: &dSort{
												arg1: &dAnd{
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
														arg2: &dEqual{
															ARG: equal__is,
															arg1: &dCollect{
																ARG:  collect__attributes__pdtype,
																arg1: &dNode{},
															},
															arg2: &dElem{
																DATA: []interface{}{"pron"},
																arg1: &dCollect{
																	ARG:  collect__attributes__pdtype,
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
										arg2: &dSort{
											arg1: &dAnd{
												arg1: &dAnd{
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
													arg2: &dFunction{
														ARG: function__not__1__args,
														arg1: &dArg{
															arg1: &dSort{
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
												},
												arg2: &dEqual{
													ARG: equal__is,
													arg1: &dCollect{
														ARG:  collect__attributes__rel,
														arg1: &dNode{},
													},
													arg2: &dElem{
														DATA: []interface{}{"su", "obj1", "obj2", "app"},
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
												arg1: &dSort{
													arg1: &dOr{
														arg1: &dOr{
															arg1: &dOr{
																arg1: &dOr{
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
																	arg2: &dSort{
																		arg1: &dAnd{
																			arg1: &dEqual{
																				ARG: equal__is,
																				arg1: &dCollect{
																					ARG:  collect__attributes__lcat,
																					arg1: &dNode{},
																				},
																				arg2: &dElem{
																					DATA: []interface{}{"np"},
																					arg1: &dCollect{
																						ARG:  collect__attributes__lcat,
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
																								ARG:  collect__attributes__rel,
																								arg1: &dNode{},
																							},
																							arg2: &dElem{
																								DATA: []interface{}{"hd", "mwp"},
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
																	arg1: &dAnd{
																		arg1: &dEqual{
																			ARG: equal__is,
																			arg1: &dCollect{
																				ARG:  collect__attributes__pt,
																				arg1: &dNode{},
																			},
																			arg2: &dElem{
																				DATA: []interface{}{"n"},
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
															arg2: &dSort{
																arg1: &dAnd{
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
																		arg2: &dEqual{
																			ARG: equal__is,
																			arg1: &dCollect{
																				ARG:  collect__attributes__pdtype,
																				arg1: &dNode{},
																			},
																			arg2: &dElem{
																				DATA: []interface{}{"pron"},
																				arg1: &dCollect{
																					ARG:  collect__attributes__pdtype,
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
														arg2: &dSort{
															arg1: &dAnd{
																arg1: &dAnd{
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
																	arg2: &dFunction{
																		ARG: function__not__1__args,
																		arg1: &dArg{
																			arg1: &dSort{
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
																},
																arg2: &dEqual{
																	ARG: equal__is,
																	arg1: &dCollect{
																		ARG:  collect__attributes__rel,
																		arg1: &dNode{},
																	},
																	arg2: &dElem{
																		DATA: []interface{}{"su", "obj1", "obj2", "app"},
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
	"nptnsub": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dAnd{
						arg1: &dAnd{
							arg1: &dSort{
								arg1: &dAnd{
									arg1: &dSort{
										arg1: &dOr{
											arg1: &dSort{
												arg1: &dAnd{
													arg1: &dSort{
														arg1: &dOr{
															arg1: &dCmp{
																ARG: cmp__gt,
																arg1: &dCollect{
																	ARG: collect__attributes__begin,
																	arg1: &dCollect{
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
																						DATA: []interface{}{"smain"},
																						arg1: &dCollect{
																							ARG:  collect__attributes__cat,
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
																arg2: &dCollect{
																	ARG: collect__attributes__begin,
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
																					DATA: []interface{}{"hd", "cmp", "mwp", "crd", "rhd", "whd", "nucl", "dp"},
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
																arg1: &dAnd{
																	arg1: &dCmp{
																		ARG: cmp__gt,
																		arg1: &dCollect{
																			ARG: collect__attributes__begin,
																			arg1: &dCollect{
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
																								DATA: []interface{}{"smain"},
																								arg1: &dCollect{
																									ARG:  collect__attributes__cat,
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
																		arg2: &dCollect{
																			ARG:  collect__attributes__begin,
																			arg1: &dNode{},
																		},
																	},
																	arg2: &dFunction{
																		ARG: function__not__1__args,
																		arg1: &dArg{
																			arg1: &dSort{
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
																								DATA: []interface{}{"hd", "cmp", "mwp", "crd", "rhd", "whd", "nucl", "dp"},
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
													arg2: &dFunction{
														ARG: function__not__1__args,
														arg1: &dArg{
															arg1: &dSort{
																arg1: &dCollect{
																	ARG:  collect__parent__node,
																	arg1: &dNode{},
																	arg2: &dPredicate{
																		arg1: &dSort{
																			arg1: &dOr{
																				arg1: &dCmp{
																					ARG: cmp__gt,
																					arg1: &dCollect{
																						ARG: collect__attributes__begin,
																						arg1: &dCollect{
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
																											DATA: []interface{}{"smain"},
																											arg1: &dCollect{
																												ARG:  collect__attributes__cat,
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
																					arg2: &dCollect{
																						ARG: collect__attributes__begin,
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
																										DATA: []interface{}{"hd", "cmp", "mwp", "crd", "rhd", "whd", "nucl", "dp"},
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
																					arg1: &dAnd{
																						arg1: &dCmp{
																							ARG: cmp__gt,
																							arg1: &dCollect{
																								ARG: collect__attributes__begin,
																								arg1: &dCollect{
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
																													DATA: []interface{}{"smain"},
																													arg1: &dCollect{
																														ARG:  collect__attributes__cat,
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
																							arg2: &dCollect{
																								ARG:  collect__attributes__begin,
																								arg1: &dNode{},
																							},
																						},
																						arg2: &dFunction{
																							ARG: function__not__1__args,
																							arg1: &dArg{
																								arg1: &dSort{
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
																													DATA: []interface{}{"hd", "cmp", "mwp", "crd", "rhd", "whd", "nucl", "dp"},
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
											},
											arg2: &dSort{
												arg1: &dAnd{
													arg1: &dSort{
														arg1: &dOr{
															arg1: &dCmp{
																ARG: cmp__gt,
																arg1: &dCollect{
																	ARG: collect__attributes__begin,
																	arg1: &dCollect{
																		ARG: collect__child__node,
																		arg1: &dCollect{
																			ARG: collect__child__node,
																			arg1: &dCollect{
																				ARG: collect__parent__node,
																				arg1: &dCollect{
																					ARG:  collect__ancestors__or__self__node,
																					arg1: &dNode{},
																					arg2: &dPredicate{
																						arg1: &dEqual{
																							ARG: equal__is,
																							arg1: &dCollect{
																								ARG:  collect__attributes__rel,
																								arg1: &dNode{},
																							},
																							arg2: &dElem{
																								DATA: []interface{}{"whd"},
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
																							ARG:  collect__attributes__cat,
																							arg1: &dNode{},
																						},
																						arg2: &dElem{
																							DATA: []interface{}{"whq"},
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
																							DATA: []interface{}{"body"},
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
																							DATA: []interface{}{"sv1"},
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
																arg2: &dCollect{
																	ARG: collect__attributes__begin,
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
																					DATA: []interface{}{"hd", "cmp", "mwp", "crd", "rhd", "whd", "nucl", "dp"},
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
																arg1: &dAnd{
																	arg1: &dCmp{
																		ARG: cmp__gt,
																		arg1: &dCollect{
																			ARG: collect__attributes__begin,
																			arg1: &dCollect{
																				ARG: collect__child__node,
																				arg1: &dCollect{
																					ARG: collect__child__node,
																					arg1: &dCollect{
																						ARG: collect__parent__node,
																						arg1: &dCollect{
																							ARG:  collect__ancestors__or__self__node,
																							arg1: &dNode{},
																							arg2: &dPredicate{
																								arg1: &dEqual{
																									ARG: equal__is,
																									arg1: &dCollect{
																										ARG:  collect__attributes__rel,
																										arg1: &dNode{},
																									},
																									arg2: &dElem{
																										DATA: []interface{}{"whd"},
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
																									ARG:  collect__attributes__cat,
																									arg1: &dNode{},
																								},
																								arg2: &dElem{
																									DATA: []interface{}{"whq"},
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
																									DATA: []interface{}{"body"},
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
																									DATA: []interface{}{"sv1"},
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
																		arg2: &dCollect{
																			ARG:  collect__attributes__begin,
																			arg1: &dNode{},
																		},
																	},
																	arg2: &dFunction{
																		ARG: function__not__1__args,
																		arg1: &dArg{
																			arg1: &dSort{
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
																								DATA: []interface{}{"hd", "cmp", "mwp", "crd", "rhd", "whd", "nucl", "dp"},
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
													arg2: &dFunction{
														ARG: function__not__1__args,
														arg1: &dArg{
															arg1: &dSort{
																arg1: &dCollect{
																	ARG:  collect__parent__node,
																	arg1: &dNode{},
																	arg2: &dPredicate{
																		arg1: &dSort{
																			arg1: &dOr{
																				arg1: &dCmp{
																					ARG: cmp__gt,
																					arg1: &dCollect{
																						ARG: collect__attributes__begin,
																						arg1: &dCollect{
																							ARG: collect__child__node,
																							arg1: &dCollect{
																								ARG: collect__child__node,
																								arg1: &dCollect{
																									ARG: collect__parent__node,
																									arg1: &dCollect{
																										ARG:  collect__ancestors__or__self__node,
																										arg1: &dNode{},
																										arg2: &dPredicate{
																											arg1: &dEqual{
																												ARG: equal__is,
																												arg1: &dCollect{
																													ARG:  collect__attributes__rel,
																													arg1: &dNode{},
																												},
																												arg2: &dElem{
																													DATA: []interface{}{"whd"},
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
																												ARG:  collect__attributes__cat,
																												arg1: &dNode{},
																											},
																											arg2: &dElem{
																												DATA: []interface{}{"whq"},
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
																												DATA: []interface{}{"body"},
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
																												DATA: []interface{}{"sv1"},
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
																					arg2: &dCollect{
																						ARG: collect__attributes__begin,
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
																										DATA: []interface{}{"hd", "cmp", "mwp", "crd", "rhd", "whd", "nucl", "dp"},
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
																					arg1: &dAnd{
																						arg1: &dCmp{
																							ARG: cmp__gt,
																							arg1: &dCollect{
																								ARG: collect__attributes__begin,
																								arg1: &dCollect{
																									ARG: collect__child__node,
																									arg1: &dCollect{
																										ARG: collect__child__node,
																										arg1: &dCollect{
																											ARG: collect__parent__node,
																											arg1: &dCollect{
																												ARG:  collect__ancestors__or__self__node,
																												arg1: &dNode{},
																												arg2: &dPredicate{
																													arg1: &dEqual{
																														ARG: equal__is,
																														arg1: &dCollect{
																															ARG:  collect__attributes__rel,
																															arg1: &dNode{},
																														},
																														arg2: &dElem{
																															DATA: []interface{}{"whd"},
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
																														ARG:  collect__attributes__cat,
																														arg1: &dNode{},
																													},
																													arg2: &dElem{
																														DATA: []interface{}{"whq"},
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
																														DATA: []interface{}{"body"},
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
																														DATA: []interface{}{"sv1"},
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
																							arg2: &dCollect{
																								ARG:  collect__attributes__begin,
																								arg1: &dNode{},
																							},
																						},
																						arg2: &dFunction{
																							ARG: function__not__1__args,
																							arg1: &dArg{
																								arg1: &dSort{
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
																													DATA: []interface{}{"hd", "cmp", "mwp", "crd", "rhd", "whd", "nucl", "dp"},
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
							arg2: &dFunction{
								ARG: function__not__1__args,
								arg1: &dArg{
									arg1: &dSort{
										arg1: &dEqual{
											ARG: equal__is,
											arg1: &dCollect{
												ARG:  collect__attributes__rel,
												arg1: &dNode{},
											},
											arg2: &dElem{
												DATA: []interface{}{"su", "sup"},
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
							arg1: &dOr{
								arg1: &dSort{
									arg1: &dOr{
										arg1: &dOr{
											arg1: &dOr{
												arg1: &dOr{
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
													arg2: &dSort{
														arg1: &dAnd{
															arg1: &dEqual{
																ARG: equal__is,
																arg1: &dCollect{
																	ARG:  collect__attributes__lcat,
																	arg1: &dNode{},
																},
																arg2: &dElem{
																	DATA: []interface{}{"np"},
																	arg1: &dCollect{
																		ARG:  collect__attributes__lcat,
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
																				ARG:  collect__attributes__rel,
																				arg1: &dNode{},
																			},
																			arg2: &dElem{
																				DATA: []interface{}{"hd", "mwp"},
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
													arg1: &dAnd{
														arg1: &dEqual{
															ARG: equal__is,
															arg1: &dCollect{
																ARG:  collect__attributes__pt,
																arg1: &dNode{},
															},
															arg2: &dElem{
																DATA: []interface{}{"n"},
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
											arg2: &dSort{
												arg1: &dAnd{
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
														arg2: &dEqual{
															ARG: equal__is,
															arg1: &dCollect{
																ARG:  collect__attributes__pdtype,
																arg1: &dNode{},
															},
															arg2: &dElem{
																DATA: []interface{}{"pron"},
																arg1: &dCollect{
																	ARG:  collect__attributes__pdtype,
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
										arg2: &dSort{
											arg1: &dAnd{
												arg1: &dAnd{
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
													arg2: &dFunction{
														ARG: function__not__1__args,
														arg1: &dArg{
															arg1: &dSort{
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
												},
												arg2: &dEqual{
													ARG: equal__is,
													arg1: &dCollect{
														ARG:  collect__attributes__rel,
														arg1: &dNode{},
													},
													arg2: &dElem{
														DATA: []interface{}{"su", "obj1", "obj2", "app"},
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
												arg1: &dSort{
													arg1: &dOr{
														arg1: &dOr{
															arg1: &dOr{
																arg1: &dOr{
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
																	arg2: &dSort{
																		arg1: &dAnd{
																			arg1: &dEqual{
																				ARG: equal__is,
																				arg1: &dCollect{
																					ARG:  collect__attributes__lcat,
																					arg1: &dNode{},
																				},
																				arg2: &dElem{
																					DATA: []interface{}{"np"},
																					arg1: &dCollect{
																						ARG:  collect__attributes__lcat,
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
																								ARG:  collect__attributes__rel,
																								arg1: &dNode{},
																							},
																							arg2: &dElem{
																								DATA: []interface{}{"hd", "mwp"},
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
																	arg1: &dAnd{
																		arg1: &dEqual{
																			ARG: equal__is,
																			arg1: &dCollect{
																				ARG:  collect__attributes__pt,
																				arg1: &dNode{},
																			},
																			arg2: &dElem{
																				DATA: []interface{}{"n"},
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
															arg2: &dSort{
																arg1: &dAnd{
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
																		arg2: &dEqual{
																			ARG: equal__is,
																			arg1: &dCollect{
																				ARG:  collect__attributes__pdtype,
																				arg1: &dNode{},
																			},
																			arg2: &dElem{
																				DATA: []interface{}{"pron"},
																				arg1: &dCollect{
																					ARG:  collect__attributes__pdtype,
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
														arg2: &dSort{
															arg1: &dAnd{
																arg1: &dAnd{
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
																	arg2: &dFunction{
																		ARG: function__not__1__args,
																		arg1: &dArg{
																			arg1: &dSort{
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
																},
																arg2: &dEqual{
																	ARG: equal__is,
																	arg1: &dCollect{
																		ARG:  collect__attributes__rel,
																		arg1: &dNode{},
																	},
																	arg2: &dElem{
																		DATA: []interface{}{"su", "obj1", "obj2", "app"},
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
	"tnonloc": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
								DATA: []interface{}{"ssub"},
								arg1: &dCollect{
									ARG:  collect__attributes__cat,
									arg1: &dNode{},
								},
							},
						},
						arg2: &dSort{
							arg1: &dAnd{
								arg1: &dSort{
									arg1: &dAnd{
										arg1: &dEqual{
											ARG: equal__is,
											arg1: &dCollect{
												ARG:  collect__attributes__cat,
												arg1: &dNode{},
											},
											arg2: &dElem{
												DATA: []interface{}{"ssub", "smain"},
												arg1: &dCollect{
													ARG:  collect__attributes__cat,
													arg1: &dNode{},
												},
											},
										},
										arg2: &dCollect{
											ARG: collect__child__node,
											arg1: &dCollect{
												ARG:  collect__descendant__or__self__type__node,
												arg1: &dNode{},
											},
											arg2: &dPredicate{
												arg1: &dSort{
													arg1: &dAnd{
														arg1: &dSort{
															arg1: &dOr{
																arg1: &dSort{
																	arg1: &dAnd{
																		arg1: &dSort{
																			arg1: &dOr{
																				arg1: &dCmp{
																					ARG: cmp__gt,
																					arg1: &dCollect{
																						ARG: collect__attributes__begin,
																						arg1: &dCollect{
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
																											DATA: []interface{}{"smain"},
																											arg1: &dCollect{
																												ARG:  collect__attributes__cat,
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
																					arg2: &dCollect{
																						ARG: collect__attributes__begin,
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
																										DATA: []interface{}{"hd", "cmp", "mwp", "crd", "rhd", "whd", "nucl", "dp"},
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
																					arg1: &dAnd{
																						arg1: &dCmp{
																							ARG: cmp__gt,
																							arg1: &dCollect{
																								ARG: collect__attributes__begin,
																								arg1: &dCollect{
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
																													DATA: []interface{}{"smain"},
																													arg1: &dCollect{
																														ARG:  collect__attributes__cat,
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
																							arg2: &dCollect{
																								ARG:  collect__attributes__begin,
																								arg1: &dNode{},
																							},
																						},
																						arg2: &dFunction{
																							ARG: function__not__1__args,
																							arg1: &dArg{
																								arg1: &dSort{
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
																													DATA: []interface{}{"hd", "cmp", "mwp", "crd", "rhd", "whd", "nucl", "dp"},
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
																		arg2: &dFunction{
																			ARG: function__not__1__args,
																			arg1: &dArg{
																				arg1: &dSort{
																					arg1: &dCollect{
																						ARG:  collect__parent__node,
																						arg1: &dNode{},
																						arg2: &dPredicate{
																							arg1: &dSort{
																								arg1: &dOr{
																									arg1: &dCmp{
																										ARG: cmp__gt,
																										arg1: &dCollect{
																											ARG: collect__attributes__begin,
																											arg1: &dCollect{
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
																																DATA: []interface{}{"smain"},
																																arg1: &dCollect{
																																	ARG:  collect__attributes__cat,
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
																										arg2: &dCollect{
																											ARG: collect__attributes__begin,
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
																															DATA: []interface{}{"hd", "cmp", "mwp", "crd", "rhd", "whd", "nucl", "dp"},
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
																										arg1: &dAnd{
																											arg1: &dCmp{
																												ARG: cmp__gt,
																												arg1: &dCollect{
																													ARG: collect__attributes__begin,
																													arg1: &dCollect{
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
																																		DATA: []interface{}{"smain"},
																																		arg1: &dCollect{
																																			ARG:  collect__attributes__cat,
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
																												arg2: &dCollect{
																													ARG:  collect__attributes__begin,
																													arg1: &dNode{},
																												},
																											},
																											arg2: &dFunction{
																												ARG: function__not__1__args,
																												arg1: &dArg{
																													arg1: &dSort{
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
																																		DATA: []interface{}{"hd", "cmp", "mwp", "crd", "rhd", "whd", "nucl", "dp"},
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
																},
																arg2: &dSort{
																	arg1: &dAnd{
																		arg1: &dSort{
																			arg1: &dOr{
																				arg1: &dCmp{
																					ARG: cmp__gt,
																					arg1: &dCollect{
																						ARG: collect__attributes__begin,
																						arg1: &dCollect{
																							ARG: collect__child__node,
																							arg1: &dCollect{
																								ARG: collect__child__node,
																								arg1: &dCollect{
																									ARG: collect__parent__node,
																									arg1: &dCollect{
																										ARG:  collect__ancestors__or__self__node,
																										arg1: &dNode{},
																										arg2: &dPredicate{
																											arg1: &dEqual{
																												ARG: equal__is,
																												arg1: &dCollect{
																													ARG:  collect__attributes__rel,
																													arg1: &dNode{},
																												},
																												arg2: &dElem{
																													DATA: []interface{}{"whd"},
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
																												ARG:  collect__attributes__cat,
																												arg1: &dNode{},
																											},
																											arg2: &dElem{
																												DATA: []interface{}{"whq"},
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
																												DATA: []interface{}{"body"},
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
																												DATA: []interface{}{"sv1"},
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
																					arg2: &dCollect{
																						ARG: collect__attributes__begin,
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
																										DATA: []interface{}{"hd", "cmp", "mwp", "crd", "rhd", "whd", "nucl", "dp"},
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
																					arg1: &dAnd{
																						arg1: &dCmp{
																							ARG: cmp__gt,
																							arg1: &dCollect{
																								ARG: collect__attributes__begin,
																								arg1: &dCollect{
																									ARG: collect__child__node,
																									arg1: &dCollect{
																										ARG: collect__child__node,
																										arg1: &dCollect{
																											ARG: collect__parent__node,
																											arg1: &dCollect{
																												ARG:  collect__ancestors__or__self__node,
																												arg1: &dNode{},
																												arg2: &dPredicate{
																													arg1: &dEqual{
																														ARG: equal__is,
																														arg1: &dCollect{
																															ARG:  collect__attributes__rel,
																															arg1: &dNode{},
																														},
																														arg2: &dElem{
																															DATA: []interface{}{"whd"},
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
																														ARG:  collect__attributes__cat,
																														arg1: &dNode{},
																													},
																													arg2: &dElem{
																														DATA: []interface{}{"whq"},
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
																														DATA: []interface{}{"body"},
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
																														DATA: []interface{}{"sv1"},
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
																							arg2: &dCollect{
																								ARG:  collect__attributes__begin,
																								arg1: &dNode{},
																							},
																						},
																						arg2: &dFunction{
																							ARG: function__not__1__args,
																							arg1: &dArg{
																								arg1: &dSort{
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
																													DATA: []interface{}{"hd", "cmp", "mwp", "crd", "rhd", "whd", "nucl", "dp"},
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
																		arg2: &dFunction{
																			ARG: function__not__1__args,
																			arg1: &dArg{
																				arg1: &dSort{
																					arg1: &dCollect{
																						ARG:  collect__parent__node,
																						arg1: &dNode{},
																						arg2: &dPredicate{
																							arg1: &dSort{
																								arg1: &dOr{
																									arg1: &dCmp{
																										ARG: cmp__gt,
																										arg1: &dCollect{
																											ARG: collect__attributes__begin,
																											arg1: &dCollect{
																												ARG: collect__child__node,
																												arg1: &dCollect{
																													ARG: collect__child__node,
																													arg1: &dCollect{
																														ARG: collect__parent__node,
																														arg1: &dCollect{
																															ARG:  collect__ancestors__or__self__node,
																															arg1: &dNode{},
																															arg2: &dPredicate{
																																arg1: &dEqual{
																																	ARG: equal__is,
																																	arg1: &dCollect{
																																		ARG:  collect__attributes__rel,
																																		arg1: &dNode{},
																																	},
																																	arg2: &dElem{
																																		DATA: []interface{}{"whd"},
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
																																	ARG:  collect__attributes__cat,
																																	arg1: &dNode{},
																																},
																																arg2: &dElem{
																																	DATA: []interface{}{"whq"},
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
																																	DATA: []interface{}{"body"},
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
																																	DATA: []interface{}{"sv1"},
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
																										arg2: &dCollect{
																											ARG: collect__attributes__begin,
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
																															DATA: []interface{}{"hd", "cmp", "mwp", "crd", "rhd", "whd", "nucl", "dp"},
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
																										arg1: &dAnd{
																											arg1: &dCmp{
																												ARG: cmp__gt,
																												arg1: &dCollect{
																													ARG: collect__attributes__begin,
																													arg1: &dCollect{
																														ARG: collect__child__node,
																														arg1: &dCollect{
																															ARG: collect__child__node,
																															arg1: &dCollect{
																																ARG: collect__parent__node,
																																arg1: &dCollect{
																																	ARG:  collect__ancestors__or__self__node,
																																	arg1: &dNode{},
																																	arg2: &dPredicate{
																																		arg1: &dEqual{
																																			ARG: equal__is,
																																			arg1: &dCollect{
																																				ARG:  collect__attributes__rel,
																																				arg1: &dNode{},
																																			},
																																			arg2: &dElem{
																																				DATA: []interface{}{"whd"},
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
																																			ARG:  collect__attributes__cat,
																																			arg1: &dNode{},
																																		},
																																		arg2: &dElem{
																																			DATA: []interface{}{"whq"},
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
																																			DATA: []interface{}{"body"},
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
																																			DATA: []interface{}{"sv1"},
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
																												arg2: &dCollect{
																													ARG:  collect__attributes__begin,
																													arg1: &dNode{},
																												},
																											},
																											arg2: &dFunction{
																												ARG: function__not__1__args,
																												arg1: &dArg{
																													arg1: &dSort{
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
																																		DATA: []interface{}{"hd", "cmp", "mwp", "crd", "rhd", "whd", "nucl", "dp"},
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
									},
								},
								arg2: &dFunction{
									ARG: function__not__1__args,
									arg1: &dArg{
										arg1: &dSort{
											arg1: &dCollect{
												ARG: collect__child__node,
												arg1: &dCollect{
													ARG:  collect__descendant__or__self__type__node,
													arg1: &dNode{},
												},
												arg2: &dPredicate{
													arg1: &dSort{
														arg1: &dAnd{
															arg1: &dEqual{
																ARG: equal__is,
																arg1: &dCollect{
																	ARG:  collect__attributes__cat,
																	arg1: &dNode{},
																},
																arg2: &dElem{
																	DATA: []interface{}{"ssub", "smain"},
																	arg1: &dCollect{
																		ARG:  collect__attributes__cat,
																		arg1: &dNode{},
																	},
																},
															},
															arg2: &dCollect{
																ARG: collect__child__node,
																arg1: &dCollect{
																	ARG:  collect__descendant__or__self__type__node,
																	arg1: &dNode{},
																},
																arg2: &dPredicate{
																	arg1: &dSort{
																		arg1: &dAnd{
																			arg1: &dSort{
																				arg1: &dOr{
																					arg1: &dSort{
																						arg1: &dAnd{
																							arg1: &dSort{
																								arg1: &dOr{
																									arg1: &dCmp{
																										ARG: cmp__gt,
																										arg1: &dCollect{
																											ARG: collect__attributes__begin,
																											arg1: &dCollect{
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
																																DATA: []interface{}{"smain"},
																																arg1: &dCollect{
																																	ARG:  collect__attributes__cat,
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
																										arg2: &dCollect{
																											ARG: collect__attributes__begin,
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
																															DATA: []interface{}{"hd", "cmp", "mwp", "crd", "rhd", "whd", "nucl", "dp"},
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
																										arg1: &dAnd{
																											arg1: &dCmp{
																												ARG: cmp__gt,
																												arg1: &dCollect{
																													ARG: collect__attributes__begin,
																													arg1: &dCollect{
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
																																		DATA: []interface{}{"smain"},
																																		arg1: &dCollect{
																																			ARG:  collect__attributes__cat,
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
																												arg2: &dCollect{
																													ARG:  collect__attributes__begin,
																													arg1: &dNode{},
																												},
																											},
																											arg2: &dFunction{
																												ARG: function__not__1__args,
																												arg1: &dArg{
																													arg1: &dSort{
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
																																		DATA: []interface{}{"hd", "cmp", "mwp", "crd", "rhd", "whd", "nucl", "dp"},
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
																							arg2: &dFunction{
																								ARG: function__not__1__args,
																								arg1: &dArg{
																									arg1: &dSort{
																										arg1: &dCollect{
																											ARG:  collect__parent__node,
																											arg1: &dNode{},
																											arg2: &dPredicate{
																												arg1: &dSort{
																													arg1: &dOr{
																														arg1: &dCmp{
																															ARG: cmp__gt,
																															arg1: &dCollect{
																																ARG: collect__attributes__begin,
																																arg1: &dCollect{
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
																																					DATA: []interface{}{"smain"},
																																					arg1: &dCollect{
																																						ARG:  collect__attributes__cat,
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
																															arg2: &dCollect{
																																ARG: collect__attributes__begin,
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
																																				DATA: []interface{}{"hd", "cmp", "mwp", "crd", "rhd", "whd", "nucl", "dp"},
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
																															arg1: &dAnd{
																																arg1: &dCmp{
																																	ARG: cmp__gt,
																																	arg1: &dCollect{
																																		ARG: collect__attributes__begin,
																																		arg1: &dCollect{
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
																																							DATA: []interface{}{"smain"},
																																							arg1: &dCollect{
																																								ARG:  collect__attributes__cat,
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
																																	arg2: &dCollect{
																																		ARG:  collect__attributes__begin,
																																		arg1: &dNode{},
																																	},
																																},
																																arg2: &dFunction{
																																	ARG: function__not__1__args,
																																	arg1: &dArg{
																																		arg1: &dSort{
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
																																							DATA: []interface{}{"hd", "cmp", "mwp", "crd", "rhd", "whd", "nucl", "dp"},
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
																					},
																					arg2: &dSort{
																						arg1: &dAnd{
																							arg1: &dSort{
																								arg1: &dOr{
																									arg1: &dCmp{
																										ARG: cmp__gt,
																										arg1: &dCollect{
																											ARG: collect__attributes__begin,
																											arg1: &dCollect{
																												ARG: collect__child__node,
																												arg1: &dCollect{
																													ARG: collect__child__node,
																													arg1: &dCollect{
																														ARG: collect__parent__node,
																														arg1: &dCollect{
																															ARG:  collect__ancestors__or__self__node,
																															arg1: &dNode{},
																															arg2: &dPredicate{
																																arg1: &dEqual{
																																	ARG: equal__is,
																																	arg1: &dCollect{
																																		ARG:  collect__attributes__rel,
																																		arg1: &dNode{},
																																	},
																																	arg2: &dElem{
																																		DATA: []interface{}{"whd"},
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
																																	ARG:  collect__attributes__cat,
																																	arg1: &dNode{},
																																},
																																arg2: &dElem{
																																	DATA: []interface{}{"whq"},
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
																																	DATA: []interface{}{"body"},
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
																																	DATA: []interface{}{"sv1"},
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
																										arg2: &dCollect{
																											ARG: collect__attributes__begin,
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
																															DATA: []interface{}{"hd", "cmp", "mwp", "crd", "rhd", "whd", "nucl", "dp"},
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
																										arg1: &dAnd{
																											arg1: &dCmp{
																												ARG: cmp__gt,
																												arg1: &dCollect{
																													ARG: collect__attributes__begin,
																													arg1: &dCollect{
																														ARG: collect__child__node,
																														arg1: &dCollect{
																															ARG: collect__child__node,
																															arg1: &dCollect{
																																ARG: collect__parent__node,
																																arg1: &dCollect{
																																	ARG:  collect__ancestors__or__self__node,
																																	arg1: &dNode{},
																																	arg2: &dPredicate{
																																		arg1: &dEqual{
																																			ARG: equal__is,
																																			arg1: &dCollect{
																																				ARG:  collect__attributes__rel,
																																				arg1: &dNode{},
																																			},
																																			arg2: &dElem{
																																				DATA: []interface{}{"whd"},
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
																																			ARG:  collect__attributes__cat,
																																			arg1: &dNode{},
																																		},
																																		arg2: &dElem{
																																			DATA: []interface{}{"whq"},
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
																																			DATA: []interface{}{"body"},
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
																																			DATA: []interface{}{"sv1"},
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
																												arg2: &dCollect{
																													ARG:  collect__attributes__begin,
																													arg1: &dNode{},
																												},
																											},
																											arg2: &dFunction{
																												ARG: function__not__1__args,
																												arg1: &dArg{
																													arg1: &dSort{
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
																																		DATA: []interface{}{"hd", "cmp", "mwp", "crd", "rhd", "whd", "nucl", "dp"},
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
																							arg2: &dFunction{
																								ARG: function__not__1__args,
																								arg1: &dArg{
																									arg1: &dSort{
																										arg1: &dCollect{
																											ARG:  collect__parent__node,
																											arg1: &dNode{},
																											arg2: &dPredicate{
																												arg1: &dSort{
																													arg1: &dOr{
																														arg1: &dCmp{
																															ARG: cmp__gt,
																															arg1: &dCollect{
																																ARG: collect__attributes__begin,
																																arg1: &dCollect{
																																	ARG: collect__child__node,
																																	arg1: &dCollect{
																																		ARG: collect__child__node,
																																		arg1: &dCollect{
																																			ARG: collect__parent__node,
																																			arg1: &dCollect{
																																				ARG:  collect__ancestors__or__self__node,
																																				arg1: &dNode{},
																																				arg2: &dPredicate{
																																					arg1: &dEqual{
																																						ARG: equal__is,
																																						arg1: &dCollect{
																																							ARG:  collect__attributes__rel,
																																							arg1: &dNode{},
																																						},
																																						arg2: &dElem{
																																							DATA: []interface{}{"whd"},
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
																																						ARG:  collect__attributes__cat,
																																						arg1: &dNode{},
																																					},
																																					arg2: &dElem{
																																						DATA: []interface{}{"whq"},
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
																																						DATA: []interface{}{"body"},
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
																																						DATA: []interface{}{"sv1"},
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
																															arg2: &dCollect{
																																ARG: collect__attributes__begin,
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
																																				DATA: []interface{}{"hd", "cmp", "mwp", "crd", "rhd", "whd", "nucl", "dp"},
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
																															arg1: &dAnd{
																																arg1: &dCmp{
																																	ARG: cmp__gt,
																																	arg1: &dCollect{
																																		ARG: collect__attributes__begin,
																																		arg1: &dCollect{
																																			ARG: collect__child__node,
																																			arg1: &dCollect{
																																				ARG: collect__child__node,
																																				arg1: &dCollect{
																																					ARG: collect__parent__node,
																																					arg1: &dCollect{
																																						ARG:  collect__ancestors__or__self__node,
																																						arg1: &dNode{},
																																						arg2: &dPredicate{
																																							arg1: &dEqual{
																																								ARG: equal__is,
																																								arg1: &dCollect{
																																									ARG:  collect__attributes__rel,
																																									arg1: &dNode{},
																																								},
																																								arg2: &dElem{
																																									DATA: []interface{}{"whd"},
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
																																								ARG:  collect__attributes__cat,
																																								arg1: &dNode{},
																																							},
																																							arg2: &dElem{
																																								DATA: []interface{}{"whq"},
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
																																								DATA: []interface{}{"body"},
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
																																								DATA: []interface{}{"sv1"},
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
																																	arg2: &dCollect{
																																		ARG:  collect__attributes__begin,
																																		arg1: &dNode{},
																																	},
																																},
																																arg2: &dFunction{
																																	ARG: function__not__1__args,
																																	arg1: &dArg{
																																		arg1: &dSort{
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
																																							DATA: []interface{}{"hd", "cmp", "mwp", "crd", "rhd", "whd", "nucl", "dp"},
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
	"locext": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
								DATA: []interface{}{"whd", "rhd"},
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
									arg1: &dSort{
										arg1: &dEqual{
											ARG: equal__is,
											arg1: &dCollect{
												ARG:  collect__attributes__index,
												arg1: &dNode{},
											},
											arg2: &dCollect{
												ARG: collect__attributes__index,
												arg1: &dCollect{
													ARG: collect__descendant__node,
													arg1: &dCollect{
														ARG: collect__child__node,
														arg1: &dCollect{
															ARG: collect__descendant__or__self__type__node,
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
																	arg1: &dEqual{
																		ARG: equal__is,
																		arg1: &dCollect{
																			ARG:  collect__attributes__cat,
																			arg1: &dNode{},
																		},
																		arg2: &dElem{
																			DATA: []interface{}{"ssub"},
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
																		ARG:  collect__attributes__cat,
																		arg1: &dNode{},
																	},
																	arg2: &dElem{
																		DATA: []interface{}{"ssub"},
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
																					ARG: collect__attributes__rel,
																					arg1: &dCollect{
																						ARG:  collect__parent__type__node,
																						arg1: &dNode{},
																					},
																				},
																				arg2: &dElem{
																					DATA: []interface{}{"obcomp"},
																					arg1: &dCollect{
																						ARG: collect__attributes__rel,
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
	"nlocext": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
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
								DATA: []interface{}{"whd", "rhd"},
								arg1: &dCollect{
									ARG:  collect__attributes__rel,
									arg1: &dNode{},
								},
							},
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
										ARG: collect__descendant__node,
										arg1: &dCollect{
											ARG: collect__child__node,
											arg1: &dCollect{
												ARG: collect__descendant__or__self__type__node,
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
														arg1: &dEqual{
															ARG: equal__is,
															arg1: &dCollect{
																ARG:  collect__attributes__cat,
																arg1: &dNode{},
															},
															arg2: &dElem{
																DATA: []interface{}{"ssub"},
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
															ARG:  collect__attributes__cat,
															arg1: &dNode{},
														},
														arg2: &dElem{
															DATA: []interface{}{"ssub"},
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
																		ARG: collect__attributes__rel,
																		arg1: &dCollect{
																			ARG:  collect__parent__type__node,
																			arg1: &dNode{},
																		},
																	},
																	arg2: &dElem{
																		DATA: []interface{}{"obcomp"},
																		arg1: &dCollect{
																			ARG: collect__attributes__rel,
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
							},
						},
					},
				},
			},
		},
	},
	"his": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dCollect{
						ARG:  collect__attributes__his,
						arg1: &dNode{},
					},
				},
			},
		},
	},
	"normal": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dEqual{
						ARG: equal__is,
						arg1: &dCollect{
							ARG:  collect__attributes__his,
							arg1: &dNode{},
						},
						arg2: &dElem{
							DATA: []interface{}{"normal"},
							arg1: &dCollect{
								ARG:  collect__attributes__his,
								arg1: &dNode{},
							},
						},
					},
				},
			},
		},
	},
	"onbeken": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dAnd{
						arg1: &dCollect{
							ARG:  collect__attributes__his,
							arg1: &dNode{},
						},
						arg2: &dFunction{
							ARG: function__not__1__args,
							arg1: &dArg{
								arg1: &dSort{
									arg1: &dEqual{
										ARG: equal__is,
										arg1: &dCollect{
											ARG:  collect__attributes__his,
											arg1: &dNode{},
										},
										arg2: &dElem{
											DATA: []interface{}{"normal", "robust_skip", "skip"},
											arg1: &dCollect{
												ARG:  collect__attributes__his,
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
	"compoun": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dEqual{
						ARG: equal__is,
						arg1: &dCollect{
							ARG:  collect__attributes__his,
							arg1: &dNode{},
						},
						arg2: &dElem{
							DATA: []interface{}{"compound"},
							arg1: &dCollect{
								ARG:  collect__attributes__his,
								arg1: &dNode{},
							},
						},
					},
				},
			},
		},
	},
	"name": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dEqual{
						ARG: equal__is,
						arg1: &dCollect{
							ARG:  collect__attributes__his,
							arg1: &dNode{},
						},
						arg2: &dElem{
							DATA: []interface{}{"name"},
							arg1: &dCollect{
								ARG:  collect__attributes__his,
								arg1: &dNode{},
							},
						},
					},
				},
			},
		},
	},
	"noun": &xPath{
		arg1: &dSort{
			arg1: &dCollect{
				ARG: collect__child__node,
				arg1: &dCollect{
					ARG:  collect__descendant__or__self__type__node,
					arg1: &dRoot{},
				},
				arg2: &dPredicate{
					arg1: &dAnd{
						arg1: &dCollect{
							ARG:  collect__attributes__his,
							arg1: &dNode{},
						},
						arg2: &dFunction{
							ARG: function__not__1__args,
							arg1: &dArg{
								arg1: &dSort{
									arg1: &dEqual{
										ARG: equal__is,
										arg1: &dCollect{
											ARG:  collect__attributes__his,
											arg1: &dNode{},
										},
										arg2: &dElem{
											DATA: []interface{}{"normal", "compound", "name", "robust_skip", "skip"},
											arg1: &dCollect{
												ARG:  collect__attributes__his,
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
}
