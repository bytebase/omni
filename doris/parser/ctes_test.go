package parser

import (
	"testing"

	"github.com/bytebase/omni/doris/ast"
)

// mustParseWithSelect parses input and returns the first statement as
// *ast.SelectStmt, asserting that it has a non-nil With field.
func mustParseWithSelect(t *testing.T, input string) *ast.SelectStmt {
	t.Helper()
	file, errs := parseForTest(input)
	if len(errs) > 0 {
		t.Fatalf("parseForTest(%q) errors: %v", input, errs)
	}
	if len(file.Stmts) == 0 {
		t.Fatalf("parseForTest(%q) returned no statements", input)
	}
	stmt, ok := file.Stmts[0].(*ast.SelectStmt)
	if !ok {
		t.Fatalf("parseForTest(%q) got %T, want *ast.SelectStmt", input, file.Stmts[0])
	}
	if stmt.With == nil {
		t.Fatalf("parseForTest(%q) stmt.With = nil, want non-nil", input)
	}
	return stmt
}

// ---------------------------------------------------------------------------
// Basic single CTE
// ---------------------------------------------------------------------------

func TestCTEBasic(t *testing.T) {
	stmt := mustParseWithSelect(t, "WITH c AS (SELECT 1) SELECT * FROM c")
	with := stmt.With

	if with.Recursive {
		t.Error("Recursive = true, want false")
	}
	if len(with.CTEs) != 1 {
		t.Fatalf("CTEs = %d, want 1", len(with.CTEs))
	}

	cte := with.CTEs[0]
	if cte.Name != "c" {
		t.Errorf("CTE[0].Name = %q, want %q", cte.Name, "c")
	}
	if len(cte.Columns) != 0 {
		t.Errorf("CTE[0].Columns = %v, want empty", cte.Columns)
	}
	if cte.Query == nil {
		t.Fatal("CTE[0].Query = nil, want non-nil")
	}
	innerStmt, ok := cte.Query.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("CTE[0].Query = %T, want *ast.SelectStmt", cte.Query)
	}
	if len(innerStmt.Items) != 1 {
		t.Fatalf("inner SELECT items = %d, want 1", len(innerStmt.Items))
	}

	// Outer SELECT should reference FROM c
	if len(stmt.From) != 1 {
		t.Fatalf("outer FROM = %d, want 1", len(stmt.From))
	}
	tbl, ok := stmt.From[0].(*ast.TableRef)
	if !ok {
		t.Fatalf("outer FROM[0] = %T, want *ast.TableRef", stmt.From[0])
	}
	if tbl.Name.Parts[0] != "c" {
		t.Errorf("outer FROM[0].Name = %q, want %q", tbl.Name.Parts[0], "c")
	}
}

// ---------------------------------------------------------------------------
// CTE with single column alias
// ---------------------------------------------------------------------------

func TestCTEWithSingleColumnAlias(t *testing.T) {
	stmt := mustParseWithSelect(t, "WITH c(x) AS (SELECT 1) SELECT * FROM c")
	cte := stmt.With.CTEs[0]

	if cte.Name != "c" {
		t.Errorf("CTE.Name = %q, want %q", cte.Name, "c")
	}
	if len(cte.Columns) != 1 {
		t.Fatalf("CTE.Columns = %v, want [x]", cte.Columns)
	}
	if cte.Columns[0] != "x" {
		t.Errorf("CTE.Columns[0] = %q, want %q", cte.Columns[0], "x")
	}
}

// ---------------------------------------------------------------------------
// CTE with multiple column aliases
// ---------------------------------------------------------------------------

func TestCTEWithMultipleColumnAliases(t *testing.T) {
	stmt := mustParseWithSelect(t, "WITH c(x, y) AS (SELECT 1, 2) SELECT * FROM c")
	cte := stmt.With.CTEs[0]

	if len(cte.Columns) != 2 {
		t.Fatalf("CTE.Columns = %v, want [x y]", cte.Columns)
	}
	if cte.Columns[0] != "x" {
		t.Errorf("CTE.Columns[0] = %q, want %q", cte.Columns[0], "x")
	}
	if cte.Columns[1] != "y" {
		t.Errorf("CTE.Columns[1] = %q, want %q", cte.Columns[1], "y")
	}
}

// ---------------------------------------------------------------------------
// Multiple CTEs
// ---------------------------------------------------------------------------

func TestCTEMultiple(t *testing.T) {
	stmt := mustParseWithSelect(t, "WITH c1 AS (SELECT 1), c2 AS (SELECT 2) SELECT * FROM c1, c2")
	with := stmt.With

	if len(with.CTEs) != 2 {
		t.Fatalf("CTEs = %d, want 2", len(with.CTEs))
	}
	if with.CTEs[0].Name != "c1" {
		t.Errorf("CTE[0].Name = %q, want %q", with.CTEs[0].Name, "c1")
	}
	if with.CTEs[1].Name != "c2" {
		t.Errorf("CTE[1].Name = %q, want %q", with.CTEs[1].Name, "c2")
	}

	// Outer SELECT should have two FROM items
	if len(stmt.From) != 2 {
		t.Fatalf("outer FROM = %d, want 2", len(stmt.From))
	}
}

// ---------------------------------------------------------------------------
// RECURSIVE keyword
// ---------------------------------------------------------------------------

func TestCTERecursive(t *testing.T) {
	// DorisParser.g4 has no RECURSIVE at the pin: the 3.1.4 server refuses
	// WITH RECURSIVE c AS (...), and reads WITH recursive AS (...) as a CTE
	// named recursive (compat.tsv's suite:conformance rows).
	if _, errs := parseForTest("WITH RECURSIVE c AS (SELECT 1) SELECT * FROM c"); len(errs) == 0 {
		t.Error("WITH RECURSIVE c parsed, want an error")
	}
	stmt := mustParseWithSelect(t, "WITH recursive AS (SELECT 1) SELECT * FROM recursive")
	if len(stmt.With.CTEs) != 1 || stmt.With.CTEs[0].Name != "recursive" {
		t.Errorf("CTEs = %+v, want one named recursive", stmt.With.CTEs)
	}
}

// ---------------------------------------------------------------------------
// Complex CTE
// ---------------------------------------------------------------------------

func TestCTEComplex(t *testing.T) {
	sql := `WITH sales AS (SELECT product, SUM(qty) FROM orders GROUP BY product) SELECT * FROM sales WHERE product > 10`
	stmt := mustParseWithSelect(t, sql)

	with := stmt.With
	if len(with.CTEs) != 1 {
		t.Fatalf("CTEs = %d, want 1", len(with.CTEs))
	}

	cte := with.CTEs[0]
	if cte.Name != "sales" {
		t.Errorf("CTE.Name = %q, want %q", cte.Name, "sales")
	}

	innerStmt, ok := cte.Query.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("CTE.Query = %T, want *ast.SelectStmt", cte.Query)
	}
	// Inner SELECT should have 2 items: product, SUM(qty)
	if len(innerStmt.Items) != 2 {
		t.Fatalf("inner SELECT items = %d, want 2", len(innerStmt.Items))
	}
	// Inner SELECT should have GROUP BY
	if len(innerStmt.GroupBy) != 1 {
		t.Fatalf("inner GROUP BY = %d, want 1", len(innerStmt.GroupBy))
	}

	// Outer SELECT: FROM sales WHERE product > 10
	if len(stmt.From) != 1 {
		t.Fatalf("outer FROM = %d, want 1", len(stmt.From))
	}
	if stmt.Where == nil {
		t.Fatal("outer WHERE = nil, want non-nil")
	}
}

// ---------------------------------------------------------------------------
// Node tag verification
// ---------------------------------------------------------------------------

func TestCTENodeTags(t *testing.T) {
	stmt := mustParseWithSelect(t, "WITH c AS (SELECT 1) SELECT * FROM c")

	if stmt.With.Tag() != ast.T_WithClause {
		t.Errorf("WithClause.Tag() = %v, want T_WithClause", stmt.With.Tag())
	}
	if stmt.With.CTEs[0].Tag() != ast.T_CTE {
		t.Errorf("CTE.Tag() = %v, want T_CTE", stmt.With.CTEs[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// Loc verification
// ---------------------------------------------------------------------------

func TestCTELoc(t *testing.T) {
	input := "WITH c AS (SELECT 1) SELECT * FROM c"
	stmt := mustParseWithSelect(t, input)

	// The SelectStmt Loc should start at the WITH token (offset 0).
	if stmt.Loc.Start != 0 {
		t.Errorf("SelectStmt.Loc.Start = %d, want 0", stmt.Loc.Start)
	}
	if stmt.Loc.End != len(input) {
		t.Errorf("SelectStmt.Loc.End = %d, want %d", stmt.Loc.End, len(input))
	}

	// The WithClause Loc should also start at 0.
	if stmt.With.Loc.Start != 0 {
		t.Errorf("WithClause.Loc.Start = %d, want 0", stmt.With.Loc.Start)
	}
}

// ---------------------------------------------------------------------------
// Walk visits CTE children
// ---------------------------------------------------------------------------

func TestCTEWalk(t *testing.T) {
	stmt := mustParseWithSelect(t, "WITH c AS (SELECT 1) SELECT * FROM c")

	var visited []ast.NodeTag
	ast.Inspect(stmt, func(n ast.Node) bool {
		if n != nil {
			visited = append(visited, n.Tag())
		}
		return true
	})

	found := func(tag ast.NodeTag) bool {
		for _, v := range visited {
			if v == tag {
				return true
			}
		}
		return false
	}

	if !found(ast.T_WithClause) {
		t.Error("Walk did not visit T_WithClause")
	}
	if !found(ast.T_CTE) {
		t.Error("Walk did not visit T_CTE")
	}
	if !found(ast.T_SelectStmt) {
		t.Error("Walk did not visit T_SelectStmt")
	}
}

// DorisParser.g4:138/145/149 put a cte before INSERT, UPDATE and DELETE, not
// only before a query.
func TestCTEPrefixedDML(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input string
		check func(t *testing.T, node ast.Node, with *ast.WithClause)
	}{
		{
			name:  "insert",
			input: "WITH c AS (SELECT 1) INSERT INTO t SELECT * FROM c",
			check: func(t *testing.T, node ast.Node, with *ast.WithClause) {
				stmt, ok := node.(*ast.InsertStmt)
				if !ok {
					t.Fatalf("got %T, want *ast.InsertStmt", node)
				}
				if stmt.With != with {
					t.Error("stmt.With is not the parsed clause")
				}
				if stmt.Query == nil {
					t.Error("Query = nil, want the SELECT source")
				}
			},
		},
		{
			name:  "update",
			input: "WITH c AS (SELECT 1 AS a) UPDATE t SET a = 2 WHERE a IN (SELECT a FROM c)",
			check: func(t *testing.T, node ast.Node, with *ast.WithClause) {
				stmt, ok := node.(*ast.UpdateStmt)
				if !ok {
					t.Fatalf("got %T, want *ast.UpdateStmt", node)
				}
				if stmt.With != with {
					t.Error("stmt.With is not the parsed clause")
				}
				if len(stmt.Assignments) != 1 {
					t.Errorf("Assignments = %d, want 1", len(stmt.Assignments))
				}
			},
		},
		{
			name:  "delete",
			input: "WITH c AS (SELECT 1 AS a) DELETE FROM t WHERE a IN (SELECT a FROM c)",
			check: func(t *testing.T, node ast.Node, with *ast.WithClause) {
				stmt, ok := node.(*ast.DeleteStmt)
				if !ok {
					t.Fatalf("got %T, want *ast.DeleteStmt", node)
				}
				if stmt.With != with {
					t.Error("stmt.With is not the parsed clause")
				}
				if stmt.Where == nil {
					t.Error("Where = nil, want the predicate")
				}
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			file, errs := parseForTest(tc.input)
			if len(errs) > 0 {
				t.Fatalf("parseForTest(%q) errors: %v", tc.input, errs)
			}
			if len(file.Stmts) != 1 {
				t.Fatalf("parseForTest(%q) returned %d statements, want 1", tc.input, len(file.Stmts))
			}
			node := file.Stmts[0]
			var with *ast.WithClause
			var loc ast.Loc
			switch n := node.(type) {
			case *ast.InsertStmt:
				with, loc = n.With, n.Loc
			case *ast.UpdateStmt:
				with, loc = n.With, n.Loc
			case *ast.DeleteStmt:
				with, loc = n.With, n.Loc
			}
			if with == nil {
				t.Fatalf("%T has no With clause", node)
			}
			if len(with.CTEs) != 1 || with.CTEs[0].Name != "c" {
				t.Errorf("CTEs = %+v, want one named c", with.CTEs)
			}
			// The statement's range covers the WITH clause it starts with.
			if loc.Start != 0 || loc.End != len(tc.input) {
				t.Errorf("Loc = %+v, want [0,%d)", loc, len(tc.input))
			}
			tc.check(t, node, with)
		})
	}

	// EXPLAIN plans a WITH-prefixed DML statement too.
	file, errs := parseForTest("EXPLAIN WITH c AS (SELECT 1 AS a) DELETE FROM t")
	if len(errs) > 0 {
		t.Fatalf("EXPLAIN WITH ... DELETE errors: %v", errs)
	}
	explain, ok := file.Stmts[0].(*ast.ExplainStmt)
	if !ok {
		t.Fatalf("got %T, want *ast.ExplainStmt", file.Stmts[0])
	}
	if _, ok := explain.Query.(*ast.DeleteStmt); !ok {
		t.Errorf("EXPLAIN's query = %T, want *ast.DeleteStmt", explain.Query)
	}
}
