package ddl

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// =============================================================================
// Round 8 Phase 6: View & CreateAs Improvements
// =============================================================================

// -----------------------------------------------------------------------------
// view.go: hasModifyingCTE check — views must not contain data-modifying CTEs
// -----------------------------------------------------------------------------

func TestViewModifyingCTEInsertError(t *testing.T) {
	c := New()
	// Create target table for the CTE.
	if err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r'); err != nil {
		t.Fatal(err)
	}
	// CREATE VIEW v AS WITH ins AS (INSERT INTO t VALUES (1) RETURNING id) SELECT * FROM ins
	viewStmt := &nodes.ViewStmt{
		View: &nodes.RangeVar{Relname: "v"},
		Query: &nodes.SelectStmt{
			WithClause: &nodes.WithClause{
				Ctes: &nodes.List{Items: []nodes.Node{
					&nodes.CommonTableExpr{
						Ctename: "ins",
						Ctequery: &nodes.InsertStmt{
							Relation: &nodes.RangeVar{Relname: "t"},
						},
					},
				}},
			},
			TargetList: &nodes.List{Items: []nodes.Node{
				&nodes.ResTarget{
					Val: &nodes.ColumnRef{Fields: &nodes.List{Items: []nodes.Node{&nodes.A_Star{}}}},
				},
			}},
			FromClause: &nodes.List{Items: []nodes.Node{
				&nodes.RangeVar{Relname: "ins"},
			}},
		},
	}
	err := c.DefineView(viewStmt)
	assertCode(t, err, CodeFeatureNotSupported)
	if !strings.Contains(err.Error(), "data-modifying statements in WITH") {
		t.Errorf("unexpected message: %s", err)
	}
}

func TestViewModifyingCTEUpdateError(t *testing.T) {
	c := New()
	if err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r'); err != nil {
		t.Fatal(err)
	}
	viewStmt := &nodes.ViewStmt{
		View: &nodes.RangeVar{Relname: "v"},
		Query: &nodes.SelectStmt{
			WithClause: &nodes.WithClause{
				Ctes: &nodes.List{Items: []nodes.Node{
					&nodes.CommonTableExpr{
						Ctename: "upd",
						Ctequery: &nodes.UpdateStmt{
							Relation: &nodes.RangeVar{Relname: "t"},
						},
					},
				}},
			},
			TargetList: &nodes.List{Items: []nodes.Node{
				&nodes.ResTarget{
					Val: &nodes.ColumnRef{Fields: &nodes.List{Items: []nodes.Node{&nodes.A_Star{}}}},
				},
			}},
			FromClause: &nodes.List{Items: []nodes.Node{
				&nodes.RangeVar{Relname: "upd"},
			}},
		},
	}
	err := c.DefineView(viewStmt)
	assertCode(t, err, CodeFeatureNotSupported)
	if !strings.Contains(err.Error(), "data-modifying statements in WITH") {
		t.Errorf("unexpected message: %s", err)
	}
}

func TestViewModifyingCTEDeleteError(t *testing.T) {
	c := New()
	if err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r'); err != nil {
		t.Fatal(err)
	}
	viewStmt := &nodes.ViewStmt{
		View: &nodes.RangeVar{Relname: "v"},
		Query: &nodes.SelectStmt{
			WithClause: &nodes.WithClause{
				Ctes: &nodes.List{Items: []nodes.Node{
					&nodes.CommonTableExpr{
						Ctename: "del",
						Ctequery: &nodes.DeleteStmt{
							Relation: &nodes.RangeVar{Relname: "t"},
						},
					},
				}},
			},
			TargetList: &nodes.List{Items: []nodes.Node{
				&nodes.ResTarget{
					Val: &nodes.ColumnRef{Fields: &nodes.List{Items: []nodes.Node{&nodes.A_Star{}}}},
				},
			}},
			FromClause: &nodes.List{Items: []nodes.Node{
				&nodes.RangeVar{Relname: "del"},
			}},
		},
	}
	err := c.DefineView(viewStmt)
	assertCode(t, err, CodeFeatureNotSupported)
	if !strings.Contains(err.Error(), "data-modifying statements in WITH") {
		t.Errorf("unexpected message: %s", err)
	}
}

func TestViewNonModifyingCTEOK(t *testing.T) {
	c := New()
	if err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r'); err != nil {
		t.Fatal(err)
	}
	// Create a view with a non-modifying CTE. Since the catalog inference engine
	// doesn't resolve CTEs as FROM sources, use the base table in FROM instead.
	// The key check is that the CTE's SELECT query passes the modifying-CTE check.
	// CREATE VIEW v AS WITH cte AS (SELECT id FROM t) SELECT id FROM t
	viewStmt := &nodes.ViewStmt{
		View: &nodes.RangeVar{Relname: "v"},
		Query: &nodes.SelectStmt{
			WithClause: &nodes.WithClause{
				Ctes: &nodes.List{Items: []nodes.Node{
					&nodes.CommonTableExpr{
						Ctename: "cte",
						Ctequery: &nodes.SelectStmt{
							TargetList: &nodes.List{Items: []nodes.Node{
								&nodes.ResTarget{
									Val: &nodes.ColumnRef{Fields: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "id"}}}},
								},
							}},
							FromClause: &nodes.List{Items: []nodes.Node{
								&nodes.RangeVar{Relname: "t"},
							}},
						},
					},
				}},
			},
			TargetList: &nodes.List{Items: []nodes.Node{
				&nodes.ResTarget{
					Val: &nodes.ColumnRef{Fields: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "id"}}}},
				},
			}},
			FromClause: &nodes.List{Items: []nodes.Node{
				&nodes.RangeVar{Relname: "t"},
			}},
		},
	}
	err := c.DefineView(viewStmt)
	if err != nil {
		t.Fatal(err)
	}
}

func TestViewNoCTEOK(t *testing.T) {
	c := New()
	if err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r'); err != nil {
		t.Fatal(err)
	}
	// CREATE VIEW v AS SELECT id FROM t — no CTE at all
	viewStmt := &nodes.ViewStmt{
		View: &nodes.RangeVar{Relname: "v"},
		Query: &nodes.SelectStmt{
			TargetList: &nodes.List{Items: []nodes.Node{
				&nodes.ResTarget{
					Val: &nodes.ColumnRef{Fields: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "id"}}}},
				},
			}},
			FromClause: &nodes.List{Items: []nodes.Node{
				&nodes.RangeVar{Relname: "t"},
			}},
		},
	}
	err := c.DefineView(viewStmt)
	if err != nil {
		t.Fatal(err)
	}
}

// -----------------------------------------------------------------------------
// createas.go: Column count error code is CodeSyntaxError
// -----------------------------------------------------------------------------

func TestCreateAsColumnCountErrorCode(t *testing.T) {
	c := New()
	if err := c.DefineRelation(makeCreateTableStmt("", "src", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "b", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r'); err != nil {
		t.Fatal(err)
	}
	// CREATE MATERIALIZED VIEW mv (x, y, z) AS SELECT a, b FROM src
	// 3 column names but query produces 2 columns — should get CodeSyntaxError.
	stmt := &nodes.CreateTableAsStmt{
		Objtype: nodes.OBJECT_MATVIEW,
		Into: &nodes.IntoClause{
			Rel: &nodes.RangeVar{Relname: "mv"},
			ColNames: &nodes.List{Items: []nodes.Node{
				&nodes.String{Str: "x"},
				&nodes.String{Str: "y"},
				&nodes.String{Str: "z"},
			}},
		},
		Query: &nodes.SelectStmt{
			TargetList: &nodes.List{Items: []nodes.Node{
				&nodes.ResTarget{
					Val: &nodes.ColumnRef{Fields: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "a"}}}},
				},
				&nodes.ResTarget{
					Val: &nodes.ColumnRef{Fields: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "b"}}}},
				},
			}},
			FromClause: &nodes.List{Items: []nodes.Node{
				&nodes.RangeVar{Relname: "src"},
			}},
		},
	}
	err := c.ExecCreateTableAs(stmt)
	assertCode(t, err, CodeSyntaxError)
	if !strings.Contains(err.Error(), "specifies 3 column names") {
		t.Errorf("unexpected message: %s", err)
	}
}

func TestCreateAsColumnCountMatchOK(t *testing.T) {
	c := New()
	if err := c.DefineRelation(makeCreateTableStmt("", "src", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "b", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r'); err != nil {
		t.Fatal(err)
	}
	// CREATE MATERIALIZED VIEW mv (x, y) AS SELECT a, b FROM src — 2 = 2, OK.
	stmt := &nodes.CreateTableAsStmt{
		Objtype: nodes.OBJECT_MATVIEW,
		Into: &nodes.IntoClause{
			Rel: &nodes.RangeVar{Relname: "mv"},
			ColNames: &nodes.List{Items: []nodes.Node{
				&nodes.String{Str: "x"},
				&nodes.String{Str: "y"},
			}},
		},
		Query: &nodes.SelectStmt{
			TargetList: &nodes.List{Items: []nodes.Node{
				&nodes.ResTarget{
					Val: &nodes.ColumnRef{Fields: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "a"}}}},
				},
				&nodes.ResTarget{
					Val: &nodes.ColumnRef{Fields: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "b"}}}},
				},
			}},
			FromClause: &nodes.List{Items: []nodes.Node{
				&nodes.RangeVar{Relname: "src"},
			}},
		},
	}
	err := c.ExecCreateTableAs(stmt)
	if err != nil {
		t.Fatal(err)
	}
}

// -----------------------------------------------------------------------------
// view.go: Modifying CTE with multiple CTEs (only one is modifying)
// -----------------------------------------------------------------------------

func TestViewMultipleCTEsOneModifyingError(t *testing.T) {
	c := New()
	if err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r'); err != nil {
		t.Fatal(err)
	}
	// WITH cte1 AS (SELECT id FROM t), cte2 AS (INSERT INTO t VALUES (1) RETURNING id)
	// SELECT * FROM cte1
	viewStmt := &nodes.ViewStmt{
		View: &nodes.RangeVar{Relname: "v"},
		Query: &nodes.SelectStmt{
			WithClause: &nodes.WithClause{
				Ctes: &nodes.List{Items: []nodes.Node{
					&nodes.CommonTableExpr{
						Ctename: "cte1",
						Ctequery: &nodes.SelectStmt{
							TargetList: &nodes.List{Items: []nodes.Node{
								&nodes.ResTarget{
									Val: &nodes.ColumnRef{Fields: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "id"}}}},
								},
							}},
							FromClause: &nodes.List{Items: []nodes.Node{
								&nodes.RangeVar{Relname: "t"},
							}},
						},
					},
					&nodes.CommonTableExpr{
						Ctename: "cte2",
						Ctequery: &nodes.InsertStmt{
							Relation: &nodes.RangeVar{Relname: "t"},
						},
					},
				}},
			},
			TargetList: &nodes.List{Items: []nodes.Node{
				&nodes.ResTarget{
					Val: &nodes.ColumnRef{Fields: &nodes.List{Items: []nodes.Node{&nodes.A_Star{}}}},
				},
			}},
			FromClause: &nodes.List{Items: []nodes.Node{
				&nodes.RangeVar{Relname: "cte1"},
			}},
		},
	}
	err := c.DefineView(viewStmt)
	assertCode(t, err, CodeFeatureNotSupported)
	if !strings.Contains(err.Error(), "data-modifying statements in WITH") {
		t.Errorf("unexpected message: %s", err)
	}
}
