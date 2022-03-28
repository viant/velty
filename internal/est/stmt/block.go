package stmt

import (
	"github.com/viant/velty/internal/est"
	"unsafe"
)

type Block struct {
	Stmt []est.Compute
}

func (s *Block) compute(state *est.State) unsafe.Pointer {
	var result unsafe.Pointer
	for i := 0; i < len(s.Stmt); i++ {
		result = s.Stmt[i](state)
	}
	return result
}

type stmt1 struct {
	est.Compute
}

func newStmt1(args []est.Compute) *stmt1 {
	return &stmt1{Compute: args[0]}
}
func (s *stmt1) compute(state *est.State) unsafe.Pointer {
	return s.Compute(state)
}

type stmt struct {
	stmt1
	est.Compute
}

func newAstmt(args []est.Compute) *stmt {
	return &stmt{Compute: args[1],
		stmt1: stmt1{
			Compute: args[0],
		},
	}
}
func (s *stmt) compute(state *est.State) unsafe.Pointer {
	s.stmt1.compute(state)
	return s.Compute(state)
}

type estmt struct {
	stmt
	est.Compute
}

func newStmt(args []est.Compute) *estmt {
	return &estmt{Compute: args[2],
		stmt: stmt{
			Compute: args[1],
			stmt1: stmt1{
				Compute: args[0],
			},
		},
	}
}
func (s *estmt) compute(state *est.State) unsafe.Pointer {
	s.stmt.compute(state)
	return s.Compute(state)
}

type stmt4 struct {
	estmt
	est.Compute
}

func newStmt4(args []est.Compute) *stmt4 {
	return &stmt4{Compute: args[3],
		estmt: estmt{
			Compute: args[2],
			stmt: stmt{
				Compute: args[1],
				stmt1: stmt1{
					Compute: args[0],
				},
			},
		},
	}
}
func (s *stmt4) compute(state *est.State) unsafe.Pointer {
	s.estmt.compute(state)
	return s.Compute(state)
}

type stmt5 struct {
	stmt4
	est.Compute
}

func newStmt5(args []est.Compute) *stmt5 {
	return &stmt5{Compute: args[4],
		stmt4: stmt4{
			Compute: args[3],
			estmt: estmt{
				Compute: args[2],
				stmt: stmt{
					Compute: args[1],
					stmt1: stmt1{
						Compute: args[0],
					},
				},
			},
		},
	}
}
func (s *stmt5) compute(state *est.State) unsafe.Pointer {
	s.stmt4.compute(state)
	return s.Compute(state)
}

type stmt6 struct {
	stmt5
	est.Compute
}

func newStmt6(args []est.Compute) *stmt6 {
	return &stmt6{Compute: args[5],
		stmt5: stmt5{
			Compute: args[4],
			stmt4: stmt4{
				Compute: args[3],
				estmt: estmt{
					Compute: args[2],
					stmt: stmt{
						Compute: args[1],
						stmt1: stmt1{
							Compute: args[0],
						},
					},
				},
			},
		},
	}
}
func (s *stmt6) compute(state *est.State) unsafe.Pointer {
	s.stmt5.compute(state)
	return s.Compute(state)
}

type stmt7 struct {
	stmt6
	est.Compute
}

func newStmt7(args []est.Compute) *stmt7 {
	return &stmt7{Compute: args[6],
		stmt6: stmt6{
			Compute: args[5],
			stmt5: stmt5{
				Compute: args[4],
				stmt4: stmt4{
					Compute: args[3],
					estmt: estmt{
						Compute: args[2],
						stmt: stmt{
							Compute: args[1],
							stmt1: stmt1{
								Compute: args[0],
							},
						},
					},
				},
			},
		},
	}
}
func (s *stmt7) compute(state *est.State) unsafe.Pointer {
	s.stmt6.compute(state)
	return s.Compute(state)
}

type stmt8 struct {
	stmt7
	est.Compute
}

func newStmt8(args []est.Compute) *stmt8 {
	return &stmt8{Compute: args[7],
		stmt7: stmt7{
			Compute: args[6],
			stmt6: stmt6{
				Compute: args[5],
				stmt5: stmt5{
					Compute: args[4],
					stmt4: stmt4{
						Compute: args[3],
						estmt: estmt{
							Compute: args[2],
							stmt: stmt{
								Compute: args[1],
								stmt1: stmt1{
									Compute: args[0],
								},
							},
						},
					},
				},
			},
		},
	}
}
func (s *stmt8) compute(state *est.State) unsafe.Pointer {
	s.stmt7.compute(state)
	return s.Compute(state)
}

type stmt9 struct {
	stmt8
	est.Compute
}

func newStmt9(args []est.Compute) *stmt9 {
	return &stmt9{Compute: args[8],
		stmt8: stmt8{
			Compute: args[7],
			stmt7: stmt7{
				Compute: args[6],
				stmt6: stmt6{
					Compute: args[5],
					stmt5: stmt5{
						Compute: args[4],
						stmt4: stmt4{
							Compute: args[3],
							estmt: estmt{
								Compute: args[2],
								stmt: stmt{
									Compute: args[1],
									stmt1: stmt1{
										Compute: args[0],
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
func (s *stmt9) compute(state *est.State) unsafe.Pointer {
	s.stmt8.compute(state)
	return s.Compute(state)
}

type stmt10 struct {
	stmt9
	est.Compute
}

func newStmt10(args []est.Compute) *stmt10 {
	return &stmt10{Compute: args[9],
		stmt9: stmt9{
			Compute: args[8],
			stmt8: stmt8{
				Compute: args[7],
				stmt7: stmt7{
					Compute: args[6],
					stmt6: stmt6{
						Compute: args[5],
						stmt5: stmt5{
							Compute: args[4],
							stmt4: stmt4{
								Compute: args[3],
								estmt: estmt{
									Compute: args[2],
									stmt: stmt{
										Compute: args[1],
										stmt1: stmt1{
											Compute: args[0],
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
}
func (s *stmt10) compute(state *est.State) unsafe.Pointer {
	s.stmt9.compute(state)
	return s.Compute(state)
}

type stmt11 struct {
	stmt10
	est.Compute
}

func newStmt11(args []est.Compute) *stmt11 {
	return &stmt11{Compute: args[10],
		stmt10: stmt10{
			Compute: args[9],
			stmt9: stmt9{
				Compute: args[8],
				stmt8: stmt8{
					Compute: args[7],
					stmt7: stmt7{
						Compute: args[6],
						stmt6: stmt6{
							Compute: args[5],
							stmt5: stmt5{
								Compute: args[4],
								stmt4: stmt4{
									Compute: args[3],
									estmt: estmt{
										Compute: args[2],
										stmt: stmt{
											Compute: args[1],
											stmt1: stmt1{
												Compute: args[0],
											},
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
}
func (s *stmt11) compute(state *est.State) unsafe.Pointer {
	s.stmt10.compute(state)
	return s.Compute(state)
}

type stmt12 struct {
	stmt11
	est.Compute
}

func newStmt12(args []est.Compute) *stmt12 {
	return &stmt12{Compute: args[11],
		stmt11: stmt11{
			Compute: args[10],
			stmt10: stmt10{
				Compute: args[9],
				stmt9: stmt9{
					Compute: args[8],
					stmt8: stmt8{
						Compute: args[7],
						stmt7: stmt7{
							Compute: args[6],
							stmt6: stmt6{
								Compute: args[5],
								stmt5: stmt5{
									Compute: args[4],
									stmt4: stmt4{
										Compute: args[3],
										estmt: estmt{
											Compute: args[2],
											stmt: stmt{
												Compute: args[1],
												stmt1: stmt1{
													Compute: args[0],
												},
											},
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
}
func (s *stmt12) compute(state *est.State) unsafe.Pointer {
	s.stmt11.compute(state)
	return s.Compute(state)
}

type stmt13 struct {
	stmt12
	est.Compute
}

func newStmt13(args []est.Compute) *stmt13 {
	return &stmt13{Compute: args[12],
		stmt12: stmt12{
			Compute: args[11],
			stmt11: stmt11{
				Compute: args[10],
				stmt10: stmt10{
					Compute: args[9],
					stmt9: stmt9{
						Compute: args[8],
						stmt8: stmt8{
							Compute: args[7],
							stmt7: stmt7{
								Compute: args[6],
								stmt6: stmt6{
									Compute: args[5],
									stmt5: stmt5{
										Compute: args[4],
										stmt4: stmt4{
											Compute: args[3],
											estmt: estmt{
												Compute: args[2],
												stmt: stmt{
													Compute: args[1],
													stmt1: stmt1{
														Compute: args[0],
													},
												},
											},
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
}
func (s *stmt13) compute(state *est.State) unsafe.Pointer {
	s.stmt12.compute(state)
	return s.Compute(state)
}

type stmt14 struct {
	stmt13
	est.Compute
}

func newStmt14(args []est.Compute) *stmt14 {
	return &stmt14{Compute: args[13],
		stmt13: stmt13{
			Compute: args[12],
			stmt12: stmt12{
				Compute: args[11],
				stmt11: stmt11{
					Compute: args[10],
					stmt10: stmt10{
						Compute: args[9],
						stmt9: stmt9{
							Compute: args[8],
							stmt8: stmt8{
								Compute: args[7],
								stmt7: stmt7{
									Compute: args[6],
									stmt6: stmt6{
										Compute: args[5],
										stmt5: stmt5{
											Compute: args[4],
											stmt4: stmt4{
												Compute: args[3],
												estmt: estmt{
													Compute: args[2],
													stmt: stmt{
														Compute: args[1],
														stmt1: stmt1{
															Compute: args[0],
														},
													},
												},
											},
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
}
func (s *stmt14) compute(state *est.State) unsafe.Pointer {
	s.stmt13.compute(state)
	return s.Compute(state)
}

type stmt15 struct {
	stmt14
	est.Compute
}

func newStmt15(args []est.Compute) *stmt15 {
	return &stmt15{Compute: args[14],
		stmt14: stmt14{
			Compute: args[13],
			stmt13: stmt13{
				Compute: args[12],
				stmt12: stmt12{
					Compute: args[11],
					stmt11: stmt11{
						Compute: args[10],
						stmt10: stmt10{
							Compute: args[9],
							stmt9: stmt9{
								Compute: args[8],
								stmt8: stmt8{
									Compute: args[7],
									stmt7: stmt7{
										Compute: args[6],
										stmt6: stmt6{
											Compute: args[5],
											stmt5: stmt5{
												Compute: args[4],
												stmt4: stmt4{
													Compute: args[3],
													estmt: estmt{
														Compute: args[2],
														stmt: stmt{
															Compute: args[1],
															stmt1: stmt1{
																Compute: args[0],
															},
														},
													},
												},
											},
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
}
func (s *stmt15) compute(state *est.State) unsafe.Pointer {
	s.stmt14.compute(state)
	return s.Compute(state)
}

type stmt16 struct {
	stmt15
	est.Compute
}

func newStmt16(args []est.Compute) *stmt16 {
	return &stmt16{Compute: args[15],
		stmt15: stmt15{
			Compute: args[14],
			stmt14: stmt14{
				Compute: args[13],
				stmt13: stmt13{
					Compute: args[12],
					stmt12: stmt12{
						Compute: args[11],
						stmt11: stmt11{
							Compute: args[10],
							stmt10: stmt10{
								Compute: args[9],
								stmt9: stmt9{
									Compute: args[8],
									stmt8: stmt8{
										Compute: args[7],
										stmt7: stmt7{
											Compute: args[6],
											stmt6: stmt6{
												Compute: args[5],
												stmt5: stmt5{
													Compute: args[4],
													stmt4: stmt4{
														Compute: args[3],
														estmt: estmt{
															Compute: args[2],
															stmt: stmt{
																Compute: args[1],
																stmt1: stmt1{
																	Compute: args[0],
																},
															},
														},
													},
												},
											},
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
}
func (s *stmt16) compute(state *est.State) unsafe.Pointer {
	s.stmt15.compute(state)
	return s.Compute(state)
}

func nop(_ *est.State) unsafe.Pointer {
	return nil
}

func NewBlock(stmtsNew []est.New) est.New {
	computers := est.Computers(stmtsNew)
	return func(control est.Control) (est.Compute, error) {
		stmts, err := computers.New(control)
		if err != nil {
			return nil, err
		}
		switch len(stmts) {
		case 0:
			return nop, nil
		case 1:
			return newStmt1(stmts).compute, nil
		case 2:
			return newAstmt(stmts).compute, nil
		case 3:
			return newStmt(stmts).compute, nil
		case 4:
			return newStmt4(stmts).compute, nil
		case 5:
			return newStmt5(stmts).compute, nil
		case 6:
			return newStmt6(stmts).compute, nil
		case 7:
			return newStmt7(stmts).compute, nil
		case 8:
			return newStmt8(stmts).compute, nil
		case 9:
			return newStmt9(stmts).compute, nil
		case 10:
			return newStmt10(stmts).compute, nil
		case 11:
			return newStmt11(stmts).compute, nil
		case 12:
			return newStmt12(stmts).compute, nil
		case 13:
			return newStmt13(stmts).compute, nil
		case 14:
			return newStmt14(stmts).compute, nil
		case 15:
			return newStmt15(stmts).compute, nil
		case 16:
			return newStmt16(stmts).compute, nil

		default:
			b := &Block{Stmt: stmts}
			return b.compute, nil
		}
	}
}
