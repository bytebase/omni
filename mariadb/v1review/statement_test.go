package v1review

import (
	"testing"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/mariadb/parser"
	"github.com/bytebase/omni/v1review"
)

// stmtFinding is a statement rule's finding as the tests compare it: the
// finding with its range replaced by the text the range anchors.
type stmtFinding struct {
	Code      v1review.Code
	Name      string
	Table     string
	Value     string
	Statement int
	Text      string
}

// runStatement parses sql and runs a statement rule over its statements.
func runStatement(t *testing.T, sql string, rule func([]ast.Node) []v1review.Finding) []stmtFinding {
	t.Helper()
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	var got []stmtFinding
	for _, f := range rule(list.Items) {
		got = append(got, stmtFinding{f.Code, f.Name, f.Table, f.Value, f.Statement, sql[f.Start:f.End]})
	}
	return got
}

func TestStatementAnchor(t *testing.T) {
	for _, tt := range []struct {
		loc        ast.Loc
		start, end int
	}{
		{ast.Loc{Start: 3, End: 9}, 3, 9},
		{ast.Loc{Start: 3, End: -1}, 3, 3},
		{ast.Loc{Start: 0, End: 0}, 0, 0},
		{ast.Loc{Start: -1, End: 5}, 0, 0},
	} {
		if start, end := statementAnchor(tt.loc); start != tt.start || end != tt.end {
			t.Errorf("statementAnchor(%+v) = %d,%d, want %d,%d", tt.loc, start, end, tt.start, tt.end)
		}
	}
}
