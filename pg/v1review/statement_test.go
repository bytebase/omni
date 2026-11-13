package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
	"github.com/bytebase/omni/v1review"
)

// stmtFinding is a statement rule's finding as the tests compare it: the
// finding with its range replaced by the text the range anchors.
type stmtFinding struct {
	Code      v1review.Code
	Name      string
	Table     string
	Count     int
	Statement int
	Text      string
}

// parseStatements parses sql into its statement nodes.
func parseStatements(t *testing.T, sql string) []ast.Node {
	t.Helper()
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	return list.Items
}

// toStmtFindings replaces each finding's range by the text of sql it anchors.
func toStmtFindings(sql string, findings []v1review.Finding) []stmtFinding {
	var got []stmtFinding
	for _, f := range findings {
		got = append(got, stmtFinding{f.Code, f.Name, f.Table, f.Count, f.Statement, sql[f.Start:f.End]})
	}
	return got
}

// statementCase is one input of a statement rule and the findings it gives.
type statementCase struct {
	name string
	sql  string
	want []stmtFinding
}

// testStatementRule runs rule over each case's parsed statements.
func testStatementRule(t *testing.T, rule func([]ast.Node) []v1review.Finding, cases []statementCase) {
	t.Helper()
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := toStmtFindings(tt.sql, rule(parseStatements(t, tt.sql)))
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAnchor(t *testing.T) {
	for _, tt := range []struct {
		loc        ast.Loc
		start, end int
	}{
		{ast.Loc{Start: 3, End: 9}, 3, 9},
		{ast.Loc{Start: 3, End: -1}, 3, 3},
		{ast.Loc{Start: 0, End: 0}, 0, 0},
		{ast.Loc{Start: -1, End: 5}, 0, 0},
	} {
		if start, end := anchor(tt.loc); start != tt.start || end != tt.end {
			t.Errorf("anchor(%+v) = %d,%d, want %d,%d", tt.loc, start, end, tt.start, tt.end)
		}
	}
}

func TestUnwrapRaw(t *testing.T) {
	stmts := parseStatements(t, "COMMIT;")
	if _, ok := stmts[0].(*ast.RawStmt); !ok {
		t.Fatalf("parser returned %T, want *ast.RawStmt", stmts[0])
	}
	if _, ok := unwrapRaw(stmts[0]).(*ast.TransactionStmt); !ok {
		t.Errorf("unwrapRaw returned %T, want *ast.TransactionStmt", unwrapRaw(stmts[0]))
	}
}
