package v1review

import (
	"testing"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
	"github.com/bytebase/omni/v1review"
)

// ruleFinding is a column rule's finding as the tests compare it: Text is
// the anchored range's text.
type ruleFinding struct {
	Code      v1review.Code
	Table     string
	Name      string
	Value     string
	Search    string
	Statement int
	Text      string
}

func runColumnRule(t *testing.T, sql string, rule func([]ast.Node) []v1review.Finding) []ruleFinding {
	t.Helper()
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	var findings []ruleFinding
	for _, f := range rule(list.Items) {
		if f.Pattern != "" || f.Detail != "" {
			t.Errorf("finding %+v: Pattern and Detail should be empty", f)
		}
		findings = append(findings, ruleFinding{f.Code, f.Table, f.Name, f.Value, f.Search, f.Statement, sql[f.Start:f.End]})
	}
	return findings
}

func TestLastTypeName(t *testing.T) {
	for sql, want := range map[string]string{
		"CREATE TABLE t (a INT);":             "int4",
		"CREATE TABLE t (a CHAR(10));":        "bpchar",
		"CREATE TABLE t (a pg_catalog.text);": "text",
		"CREATE TABLE t (a s.\"MyType\"[]);":  "MyType",
		"CREATE TABLE t OF typ (a NOT NULL);": "",
	} {
		list, err := parser.Parse(sql)
		if err != nil {
			t.Fatalf("Parse %q: %v", sql, err)
		}
		cols := tableColumns(statementNode(list.Items[0]).(*ast.CreateStmt))
		if got := lastTypeName(cols[0].TypeName); got != want {
			t.Errorf("%s: lastTypeName = %q, want %q", sql, got, want)
		}
	}
}
