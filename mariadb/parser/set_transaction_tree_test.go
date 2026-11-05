package parser

import (
	"errors"
	"testing"

	"github.com/bytebase/omni/mariadb/ast"
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

// TestSetTransactionFields holds every field to the text it came from.
func TestSetTransactionFields(t *testing.T) {
	for _, tc := range []struct {
		sql    string
		scope  string
		level  string
		access string
	}{
		{sql: "set transaction read only", access: "READ ONLY"},
		{sql: "set transaction read write", access: "READ WRITE"},
		{sql: "set transaction isolation level repeatable read", level: "REPEATABLE READ"},
		{sql: "set transaction isolation level read committed", level: "READ COMMITTED"},
		{sql: "set transaction isolation level read uncommitted", level: "READ UNCOMMITTED"},
		{sql: "set transaction isolation level serializable", level: "SERIALIZABLE"},
		{sql: "set transaction read only, isolation level serializable", level: "SERIALIZABLE", access: "READ ONLY"},
		{sql: "set transaction isolation level read committed, read write", level: "READ COMMITTED", access: "READ WRITE"},
		{sql: "set session transaction read only", scope: "SESSION", access: "READ ONLY"},
		{sql: "set global transaction isolation level serializable", scope: "GLOBAL", level: "SERIALIZABLE"},
		{sql: "set local transaction read write", scope: "LOCAL", access: "READ WRITE"},
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
			if got := tc.sql[stmt.Loc.Start:stmt.Loc.End]; got != tc.sql {
				t.Errorf("Loc %+v covers %q, want the statement", stmt.Loc, got)
			}
		})
	}
}

// TestSetTransactionErrorPosition checks each refused form reports the token it stopped
// at, an empty string meaning the end of the text.
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
		{"set transaction read only, read write", "read write"},
		{"set transaction isolation level serializable, isolation level read committed", "isolation level read committed"},
		{"set transaction read only, isolation level serializable, read write", ", read write"},
		{"set persist transaction read only", "transaction read only"},
		{"set persist_only transaction read only", "transaction read only"},
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
	const sql = "set session transaction isolation level read committed, read only"
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
