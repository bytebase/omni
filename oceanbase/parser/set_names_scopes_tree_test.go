package parser

import (
	"errors"
	"slices"
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// parseSetStmtForTest parses one SET statement.
func parseSetStmtForTest(t *testing.T, sql string) *nodes.SetStmt {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	stmt, ok := list.Items[0].(*nodes.SetStmt)
	if !ok {
		t.Fatalf("Parse(%q): %T, want *ast.SetStmt", sql, list.Items[0])
	}
	return stmt
}

// TestSetVariableNameTrees pins what the word turned into, because a verdict cannot tell a
// variable from the statement or the scope the same word can start.
func TestSetVariableNameTrees(t *testing.T) {
	for _, c := range []struct{ sql, scope, table, column string }{
		{"set a = 1", "", "", "a"},
		{"set a to 1", "", "", "a"},
		{"set session a to 1", "SESSION", "", "a"},
		{"set local = 1", "", "", "local"},
		{"set persist = 1", "", "", "persist"},
		{"set global = 1", "", "", "global"},
		{"set global global = 1", "GLOBAL", "", "global"},
		{"set transaction = 1", "", "", "transaction"},
		{"set transaction to 1", "", "", "transaction"},
		{"set names = 1", "", "", "names"},
		{"set role = 1", "", "", "role"},
		{"set charset = 1", "", "", "charset"},
		{"set new = 1", "", "", "new"},
	} {
		t.Run(c.sql, func(t *testing.T) {
			stmt := parseSetStmtForTest(t, c.sql)
			if len(stmt.Assignments) != 1 {
				t.Fatalf("got %d assignments, want 1", len(stmt.Assignments))
			}
			a := stmt.Assignments[0]
			if stmt.Scope != c.scope || a.Scope != c.scope || a.Column.Table != c.table || a.Column.Column != c.column {
				t.Errorf("stmt scope=%q scope=%q table=%q column=%q, want scope=%q table=%q column=%q",
					stmt.Scope, a.Scope, a.Column.Table, a.Column.Column, c.scope, c.table, c.column)
			}
			if stmt.Loc.Start != 0 || stmt.Loc.End != len(c.sql) {
				t.Errorf("statement Loc = %+v, want {0 %d}", stmt.Loc, len(c.sql))
			}
			if got := c.sql[a.Loc.Start:a.Loc.End]; got != c.sql[len("set "):] {
				t.Errorf("assignment Loc covers %q, want %q", got, c.sql[len("set "):])
			}
		})
	}
}

// TestSetAssignmentScopes checks each assignment holds the scope it wrote, and the statement
// the first one's. The words the same scope starts stay themselves.
func TestSetAssignmentScopes(t *testing.T) {
	const sql = "set @a = 1, global b = 2, @@sql_mode = '', session c = 3, global d = 4"
	stmt := parseSetStmtForTest(t, sql)
	var got []string
	for _, a := range stmt.Assignments {
		got = append(got, a.Scope)
	}
	if want := []string{"", "GLOBAL", "", "SESSION", "GLOBAL"}; !slices.Equal(got, want) {
		t.Errorf("scopes = %q, want %q", got, want)
	}
	if stmt.Scope != "" {
		t.Errorf("statement scope = %q, want the first assignment's", stmt.Scope)
	}
	if text := sql[stmt.Assignments[1].Loc.Start:stmt.Assignments[1].Loc.End]; text != "global b = 2" {
		t.Errorf("assignment Loc covers %q", text)
	}
	for _, form := range []string{"set transaction read only", "set session transaction read only"} {
		list, err := Parse(form)
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := list.Items[0].(*nodes.SetTransactionStmt); !ok {
			t.Errorf("%q: got %T, want *ast.SetTransactionStmt", form, list.Items[0])
		}
	}
	if dump := nodes.NodeToString(stmt.Assignments[1]); !containsScope(dump) {
		t.Errorf("NodeToString has no scope: %s", dump)
	}
}

func containsScope(dump string) bool {
	for i := 0; i+len(":scope GLOBAL") <= len(dump); i++ {
		if dump[i:i+len(":scope GLOBAL")] == ":scope GLOBAL" {
			return true
		}
	}
	return false
}

// TestSetNamesScopesErrorPosition checks each refused form reports the token it stopped at,
// an empty string meaning the end of the text.
func TestSetNamesScopesErrorPosition(t *testing.T) {
	for _, tc := range []struct{ sql, at string }{
		{"set session @a = 1", "@a = 1"},
		{"set global @@sql_mode = ''", "@@sql_mode = ''"},
		{"set @a = 1, global @b = 2", "@b = 2"},
		{"set @@global = 1", "@@global = 1"},
		{"set a.b = 1", ".b = 1"},
		{"set global a.b = 1", ".b = 1"},
		{"set @@a.b = 1", "@@a.b = 1"},
		{"set @@persist.a = 1", "@@persist.a = 1"},
		{"set local a = 1", "a = 1"},
		{"set persist a = 1", "a = 1"},
		{"set @a = 1, local a = 2", "a = 2"},
		{"set @a = 1, global", ""},
		{"set @a = 1, password = 2", "password = 2"},
		{"set session transaction = 1", "= 1"},
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

// TestSetNamesScopesTruncated cuts the longest form at every byte. Each prefix either parses
// or returns a ParseError holding an offset inside the text, and none panics.
func TestSetNamesScopesTruncated(t *testing.T) {
	const sql = "set @a = 1, global b = 2, @@sql_mode = '', session c = 3, global d = 4"
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
