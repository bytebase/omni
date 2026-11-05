package parser

import (
	"errors"
	"strings"
	"testing"

	"github.com/bytebase/omni/tidb/ast"
)

// parseSetTransaction parses one SET TRANSACTION statement.
func parseSetTransaction(t *testing.T, sql string) *ast.SetTransactionStmt {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("Parse(%q): %d statements, want 1", sql, len(list.Items))
	}
	stmt, ok := list.Items[0].(*ast.SetTransactionStmt)
	if !ok {
		t.Fatalf("Parse(%q): %T, want *ast.SetTransactionStmt", sql, list.Items[0])
	}
	return stmt
}

// TestSetTransactionFields holds every field to the text it came from. A statement that
// names a characteristic twice keeps the last, and READ WRITE leaves an earlier AS OF.
func TestSetTransactionFields(t *testing.T) {
	for _, tc := range []struct {
		sql    string
		scope  string
		level  string
		access string
		asOf   string
	}{
		{sql: "set transaction read only", access: "READ ONLY"},
		{sql: "set transaction read write", access: "READ WRITE"},
		{sql: "set transaction isolation level repeatable read", level: "REPEATABLE READ"},
		{sql: "set transaction isolation level read committed", level: "READ COMMITTED"},
		{sql: "set transaction isolation level read uncommitted", level: "READ UNCOMMITTED"},
		{sql: "set transaction isolation level serializable", level: "SERIALIZABLE"},
		{sql: "set transaction read only as of timestamp 7", access: "READ ONLY", asOf: "7"},
		{sql: "set session transaction read only as of timestamp 7", scope: "SESSION", access: "READ ONLY", asOf: "7"},
		{sql: "set global transaction read only as of timestamp 7", scope: "GLOBAL", access: "READ ONLY", asOf: "7"},
		{sql: "set transaction isolation level serializable, read only as of timestamp 7", level: "SERIALIZABLE", access: "READ ONLY", asOf: "7"},
		{sql: "set transaction read only as of timestamp 7, isolation level serializable", level: "SERIALIZABLE", access: "READ ONLY", asOf: "7"},
		{sql: "set transaction read only as of timestamp 1, read only as of timestamp 2", access: "READ ONLY", asOf: "2"},
		{sql: "set transaction read only as of timestamp 7, read write", access: "READ WRITE", asOf: "7"},
		{sql: "set transaction read only as of timestamp 7, read only", access: "READ ONLY", asOf: "7"},
		{sql: "set transaction read only, read write", access: "READ WRITE"},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			stmt := parseSetTransaction(t, tc.sql)
			if stmt.Scope != tc.scope {
				t.Errorf("Scope = %q, want %q", stmt.Scope, tc.scope)
			}
			if stmt.IsolationLevel != tc.level {
				t.Errorf("IsolationLevel = %q, want %q", stmt.IsolationLevel, tc.level)
			}
			if stmt.AccessMode != tc.access {
				t.Errorf("AccessMode = %q, want %q", stmt.AccessMode, tc.access)
			}
			if got := exprText(t, stmt.AsOf); got != tc.asOf {
				t.Errorf("AsOf = %q, want %q", got, tc.asOf)
			}
			if got := tc.sql[stmt.Loc.Start:stmt.Loc.End]; got != tc.sql {
				t.Errorf("Loc %+v covers %q, want the statement", stmt.Loc, got)
			}
		})
	}
}

// TestSetTransactionWalk checks a consumer that walks the tree reaches the expression,
// and that the node dump shows it.
func TestSetTransactionWalk(t *testing.T) {
	list, err := Parse("set transaction read only as of timestamp (select 1)")
	if err != nil {
		t.Fatal(err)
	}
	var subqueries int
	ast.Inspect(list.Items[0], func(n ast.Node) bool {
		if _, ok := n.(*ast.SubqueryExpr); ok {
			subqueries++
		}
		return true
	})
	if subqueries != 1 {
		t.Errorf("walk reached %d subqueries, want the expression's one", subqueries)
	}
	if dump := ast.NodeToString(list.Items[0]); !strings.Contains(dump, ":as_of {SUBQUERY") {
		t.Errorf("NodeToString has no expression: %s", dump)
	}
}

// TestSetTransactionErrorPosition checks each form cut short reports the token it
// stopped at, an empty string meaning the end of the text.
func TestSetTransactionErrorPosition(t *testing.T) {
	for _, tc := range []struct{ sql, at string }{
		{"set transaction", ""},
		{"set transaction foo", "foo"},
		{"set transaction read", ""},
		{"set transaction read committed", "committed"},
		{"set transaction read only,", ""},
		{"set transaction isolation serializable", "serializable"},
		{"set transaction isolation level", ""},
		{"set transaction isolation level repeatable", ""},
		{"set transaction isolation level read", ""},
		{"set transaction isolation level read only", "only"},
		{"set transaction isolation level snapshot", "snapshot"},
		{"set transaction read only as of", ""},
		{"set transaction read only as now()", "now()"},
		{"set transaction read write as of timestamp 1", "as of timestamp 1"},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			_, err := Parse(tc.sql)
			var pe *ParseError
			if !errors.As(err, &pe) {
				t.Fatalf("err = %v, want a ParseError", err)
			}
			if got := tc.sql[pe.Position:]; got != tc.at {
				t.Errorf("stopped at %q, want %q", got, tc.at)
			}
		})
	}
}

// TestSetTransactionTruncated cuts the longest form at every byte. Each prefix either
// parses or returns a ParseError holding an offset inside the text, and none panics.
func TestSetTransactionTruncated(t *testing.T) {
	const sql = "set session transaction isolation level read committed, read only as of timestamp now() - interval 1 second, read write"
	for i := 0; i <= len(sql); i++ {
		_, err := Parse(sql[:i])
		if err == nil {
			continue
		}
		var pe *ParseError
		if !errors.As(err, &pe) {
			t.Fatalf("Parse(%q): %T is not a ParseError", sql[:i], err)
		}
		if pe.Position < 0 || pe.Position > i {
			t.Errorf("Parse(%q): error position %d is outside the text", sql[:i], pe.Position)
		}
	}
}

// TestSetTransactionTreeHeight checks the expression counts toward the height the parser
// measures, and a nested one toward the depth.
func TestSetTransactionTreeHeight(t *testing.T) {
	chain := strings.Repeat("1+", maxTreeHeight+1) + "1"
	_, err := Parse("set transaction read only as of timestamp " + chain)
	var pe *ParseError
	if !errors.As(err, &pe) || pe.Message != "statement is nested too deeply" {
		t.Errorf("err = %v, want the nesting error", briefErr(err))
	}
	nested := strings.Repeat("(", maxNestingDepth+1) + "1" + strings.Repeat(")", maxNestingDepth+1)
	if _, err := Parse("set transaction read only as of timestamp " + nested); err == nil {
		t.Error("an expression nested past maxNestingDepth parsed")
	}
}
