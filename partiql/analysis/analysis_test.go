package analysis

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/partiql/ast"
	"github.com/bytebase/omni/partiql/parser"
)

// parseSelect is a test helper that parses a SELECT and returns the inner
// *ast.SelectStmt, failing the test if the parse fails or the result isn't
// a SelectStmt.
func parseSelect(t *testing.T, input string) *ast.SelectStmt {
	t.Helper()
	list, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse(%q): %v", input, err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("Parse(%q): got %d items, want 1", input, len(list.Items))
	}
	sel, ok := list.Items[0].(*ast.SelectStmt)
	if !ok {
		t.Fatalf("Parse(%q): got %T, want *ast.SelectStmt", input, list.Items[0])
	}
	return sel
}

// --- ValidateQuery ---

func TestValidateQuery_DQLAllowed(t *testing.T) {
	cases := []string{
		"SELECT * FROM t",
		"SELECT a, b FROM t WHERE x = 1",
		"SELECT * FROM t1 UNION SELECT * FROM t2",
		"SELECT * FROM t1 INTERSECT SELECT * FROM t2",
		"EXPLAIN SELECT * FROM t",
		"EXPLAIN UPDATE t SET a = 1", // EXPLAIN is always read-only
	}
	for _, input := range cases {
		input := input
		t.Run(input[:min(len(input), 30)], func(t *testing.T) {
			if err := ValidateQuery(input); err != nil {
				t.Errorf("ValidateQuery(%q) = %v, want nil", input, err)
			}
		})
	}
}

func TestValidateQuery_DMLRejected(t *testing.T) {
	cases := []struct {
		input string
	}{
		{"INSERT INTO t VALUE 1"},
		{"DELETE FROM t"},
		{"UPDATE t SET a = 1"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.input[:min(len(tc.input), 30)], func(t *testing.T) {
			err := ValidateQuery(tc.input)
			if err == nil {
				t.Fatalf("ValidateQuery(%q) = nil, want DML error", tc.input)
			}
			if !strings.Contains(err.Error(), "DML") {
				t.Errorf("error = %q, want to contain \"DML\"", err.Error())
			}
		})
	}
}

func TestValidateQuery_DDLRejected(t *testing.T) {
	cases := []struct {
		input string
	}{
		{"CREATE TABLE t"},
		{"DROP TABLE t"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.input[:min(len(tc.input), 30)], func(t *testing.T) {
			err := ValidateQuery(tc.input)
			if err == nil {
				t.Fatalf("ValidateQuery(%q) = nil, want DDL error", tc.input)
			}
			if !strings.Contains(err.Error(), "DDL") {
				t.Errorf("error = %q, want to contain \"DDL\"", err.Error())
			}
		})
	}
}

func TestValidateQuery_EXECRejected(t *testing.T) {
	err := ValidateQuery("EXEC myProc")
	if err == nil {
		t.Fatal("ValidateQuery(\"EXEC myProc\") = nil, want EXEC error")
	}
	if !strings.Contains(err.Error(), "EXEC") {
		t.Errorf("error = %q, want to contain \"EXEC\"", err.Error())
	}
}

func TestValidateQuery_ParseError(t *testing.T) {
	// Should propagate a parse error, not panic.
	err := ValidateQuery("SELECT FROM WHERE")
	if err == nil {
		t.Fatal("ValidateQuery(invalid SQL) = nil, want parse error")
	}
}

// --- Analyze ---

func TestAnalyze_SelectStar(t *testing.T) {
	sel := parseSelect(t, "SELECT * FROM t")
	qa := Analyze(sel)
	if !qa.SelectStar {
		t.Error("SelectStar: got false, want true")
	}
	if len(qa.Projections) != 0 {
		t.Errorf("Projections: got %d, want 0", len(qa.Projections))
	}
	if len(qa.Tables) != 1 || qa.Tables[0] != "t" {
		t.Errorf("Tables: got %v, want [t]", qa.Tables)
	}
}

func TestAnalyze_SimpleProjections(t *testing.T) {
	sel := parseSelect(t, "SELECT a, b FROM t")
	qa := Analyze(sel)
	if qa.SelectStar {
		t.Error("SelectStar: got true, want false")
	}
	if len(qa.Projections) != 2 {
		t.Fatalf("Projections: got %d, want 2", len(qa.Projections))
	}
}

func TestAnalyze_AliasedProjection(t *testing.T) {
	sel := parseSelect(t, "SELECT a AS myalias FROM t")
	qa := Analyze(sel)
	if len(qa.Projections) != 1 {
		t.Fatalf("Projections: got %d, want 1", len(qa.Projections))
	}
	if qa.Projections[0].Name != "myalias" {
		t.Errorf("Name: got %q, want %q", qa.Projections[0].Name, "myalias")
	}
}

func TestAnalyze_SelectValue(t *testing.T) {
	sel := parseSelect(t, "SELECT VALUE a FROM t")
	qa := Analyze(sel)
	if qa.SelectStar {
		t.Error("SelectStar: got true, want false")
	}
	if len(qa.Projections) != 1 {
		t.Fatalf("Projections: got %d, want 1", len(qa.Projections))
	}
	if qa.Projections[0].Expr == "" {
		t.Error("Expr: got empty, want non-empty")
	}
}

func TestAnalyze_Tables(t *testing.T) {
	sel := parseSelect(t, "SELECT * FROM t")
	qa := Analyze(sel)
	if len(qa.Tables) != 1 || qa.Tables[0] != "t" {
		t.Errorf("Tables: got %v, want [t]", qa.Tables)
	}
}

func TestAnalyze_JoinTables(t *testing.T) {
	sel := parseSelect(t, "SELECT * FROM t JOIN s ON t.id = s.id")
	qa := Analyze(sel)
	if len(qa.Tables) != 2 {
		t.Fatalf("Tables: got %v, want 2 entries", qa.Tables)
	}
	tableSet := map[string]bool{}
	for _, tbl := range qa.Tables {
		tableSet[tbl] = true
	}
	if !tableSet["t"] || !tableSet["s"] {
		t.Errorf("Tables: got %v, want [t s]", qa.Tables)
	}
}

func TestAnalyze_LiteralProjection(t *testing.T) {
	sel := parseSelect(t, "SELECT 1 FROM t")
	qa := Analyze(sel)
	if len(qa.Projections) != 1 {
		t.Fatalf("Projections: got %d, want 1", len(qa.Projections))
	}
	if qa.Projections[0].Expr == "" {
		t.Error("Expr: got empty, want non-empty")
	}
}

func TestAnalyze_WhereClause(t *testing.T) {
	// WHERE doesn't affect Projections/SelectStar/Tables but should not panic.
	sel := parseSelect(t, "SELECT * FROM t WHERE x = 1")
	qa := Analyze(sel)
	if !qa.SelectStar {
		t.Error("SelectStar: got false, want true")
	}
	if len(qa.Tables) != 1 {
		t.Errorf("Tables: got %v, want [t]", qa.Tables)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
