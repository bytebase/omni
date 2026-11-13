package v1review

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/v1review"
)

// stmtFinding is a statement rule's finding as the tests compare it: the
// finding with its range replaced by the text the range anchors.
type stmtFinding struct {
	Code v1review.Code
	Name string
	Text string
}

// runStatement parses sql and runs a statement rule over its statements.
func runStatement(t *testing.T, sql string, rule func([]ast.Node) []v1review.Finding) []stmtFinding {
	t.Helper()
	var got []stmtFinding
	for _, f := range rule(parse(t, sql)) {
		got = append(got, stmtFinding{f.Code, f.Name, sql[f.Start:f.End]})
	}
	return got
}
