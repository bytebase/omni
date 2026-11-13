package v1review

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/oracle/parser"
	"github.com/bytebase/omni/v1review"
)

// columnFinding is a column rule's finding as the tests compare it: the
// text its range anchors in place of the range.
type columnFinding struct {
	Code      v1review.Code
	Name      string
	Table     string
	Value     string
	Statement int
	Text      string
}

// parseStmts parses sql into its top-level statements, unwrapped from
// their *ast.RawStmt as Bytebase passes them.
func parseStmts(t *testing.T, sql string) []ast.Node {
	t.Helper()
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	var stmts []ast.Node
	for _, item := range list.Items {
		if raw, ok := item.(*ast.RawStmt); ok {
			item = raw.Stmt
		}
		stmts = append(stmts, item)
	}
	return stmts
}

func columnFindings(sql string, got []v1review.Finding) []columnFinding {
	var findings []columnFinding
	for _, f := range got {
		findings = append(findings, columnFinding{f.Code, f.Name, f.Table, f.Value, f.Statement, sql[f.Start:f.End]})
	}
	return findings
}

func TestObjectName(t *testing.T) {
	for _, tt := range []struct {
		name *ast.ObjectName
		db   string
		want string
	}{
		{nil, "DB", ""},
		{&ast.ObjectName{Name: "T"}, "DB", "DB.T"},
		{&ast.ObjectName{Name: "T"}, "", "T"},
		{&ast.ObjectName{Schema: "S", Name: "T"}, "DB", "S.T"},
	} {
		if got := objectName(tt.name, tt.db); got != tt.want {
			t.Errorf("objectName(%v, %q) = %q, want %q", tt.name, tt.db, got, tt.want)
		}
	}
}

func TestLastIdentifier(t *testing.T) {
	for key, want := range map[string]string{"DB.T.C": "C", "DB.T.a.b": "b", "C": "C", "DB.T.": ""} {
		if got := lastIdentifier(key); got != want {
			t.Errorf("lastIdentifier(%q) = %q, want %q", key, got, want)
		}
	}
}
