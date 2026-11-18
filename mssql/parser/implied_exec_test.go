package parser

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// TestImpliedExecStmt holds a procedure called by name without EXEC as a
// batch's first statement: optSimpleExecute (TSql160.g:567), which batch
// (TSql160.g:158) takes first and the same text elsewhere in a batch is not
// a statement at all.
func TestImpliedExecStmt(t *testing.T) {
	tests := []struct {
		sql    string
		server string // TableRef.Server
		schema string // TableRef.Schema
		object string // TableRef.Object
		module string // ExecStmt.ModuleVar
		args   int
	}{
		{sql: "my_proc", object: "my_proc"},
		{sql: "my_proc 1, 2", object: "my_proc", args: 2},
		{sql: "my_proc;", object: "my_proc"},
		{sql: "my_proc 1, @p = 'a'", object: "my_proc", args: 2},
		{sql: "dbo.my_proc 1", schema: "dbo", object: "my_proc", args: 1},
		{sql: "db.dbo.my_proc", schema: "dbo", object: "my_proc"},
		// identifierList (TSql160.g:30325) starts with an empty part too.
		{sql: ".my_proc 1", object: "my_proc", args: 1},
		{sql: "..my_proc", object: "my_proc"},
		// procNumOpt (TSql160.g:26505): a procedure number.
		{sql: "my_proc; 1", object: "my_proc"},
		{sql: "my_proc; 1 'a'", object: "my_proc", args: 1},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			list := ParseAndCheck(t, tt.sql)
			if list.Len() != 1 {
				t.Fatalf("Parse(%q): %d statements, want 1", tt.sql, list.Len())
			}
			stmt, ok := list.Items[0].(*nodes.ExecStmt)
			if !ok {
				t.Fatalf("Parse(%q): %T, want *asta.ExecStmt", tt.sql, list.Items[0])
			}
			if !stmt.Implied {
				t.Errorf("Parse(%q): Implied = false, want true", tt.sql)
			}
			if stmt.ModuleVar != tt.module {
				t.Errorf("Parse(%q): ModuleVar = %q, want %q", tt.sql, stmt.ModuleVar, tt.module)
			}
			if stmt.Name != nil {
				if got := [3]string{stmt.Name.Server, stmt.Name.Schema, stmt.Name.Object}; got != [3]string{tt.server, tt.schema, tt.object} {
					t.Errorf("Parse(%q): name = server %q schema %q object %q, want %q", tt.sql, got[0], got[1], got[2], [3]string{tt.server, tt.schema, tt.object})
				}
			} else if tt.object != "" {
				t.Errorf("Parse(%q): Name = nil, want object %q", tt.sql, tt.object)
			}
			// The name ends at its last part: the arguments and the
			// procedure number are none of it, as in EXEC.
			if stmt.Name != nil {
				want := strings.IndexAny(tt.sql, " ;")
				if want < 0 {
					want = len(tt.sql)
				}
				if stmt.Name.Loc.End != want {
					t.Errorf("Parse(%q): Name.Loc.End = %d, want %d", tt.sql, stmt.Name.Loc.End, want)
				}
			}
			if n := argsLen(stmt); n != tt.args {
				t.Errorf("Parse(%q): %d arguments, want %d", tt.sql, n, tt.args)
			}
			// The statement covers its text, its ; ending an empty statement.
			if stmt.Loc.Start != 0 {
				t.Errorf("Parse(%q): Loc.Start = %d, want 0", tt.sql, stmt.Loc.Start)
			}
			if end := len(tt.sql); stmt.Loc.End != end && stmt.Loc.End != end-1 {
				t.Errorf("Parse(%q): Loc.End = %d, want %d", tt.sql, stmt.Loc.End, end)
			}
		})
	}
}

// argsLen counts an EXEC's arguments.
func argsLen(stmt *nodes.ExecStmt) int {
	if stmt.Args == nil {
		return 0
	}
	return len(stmt.Args.Items)
}

// TestImpliedExecStmtBatchStart holds the call to the batch's first statement:
// a GO starts a batch, a semicolon before the name takes its place away
// (Msg 102), and the words optSimpleExecute's predicate keeps start their own
// statements. A statement after the call is read as SQL Server reads it,
// except that a variable it uses must be one the batch declares first (Msg
// 137, measured): the call is where a declaration cannot precede.
func TestImpliedExecStmtBatchStart(t *testing.T) {
	// A GO starts a batch, so the call is its first statement again, and
	// the statements after a later GO are statements.
	list := ParseAndCheck(t, "SELECT 1;\nGO\nmy_proc 1;\nGO\nSELECT 2")
	if list.Len() != 5 {
		t.Fatalf("%d statements, want 5: %s", list.Len(), nodes.NodeToString(list))
	}
	stmt, ok := list.Items[2].(*nodes.ExecStmt)
	if !ok || !stmt.Implied {
		t.Fatalf("statement after GO = %s, want an *ast.ExecStmt with Implied", nodes.NodeToString(list.Items[2]))
	}

	// Statements after the call, and declarations its batch may hold.
	for _, sql := range []string{
		"my_proc; SELECT 1",
		"my_proc DECLARE @a INT",
		"my_proc; DECLARE @a INT; SELECT @a",
		"my_proc @a = 1",
		"my_proc @a = 1, @b = 2",
		"my_proc; EXEC p @x = 1",
		"my_proc; SELECT @@ROWCOUNT",
		"my_proc PRINT '@a'",
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v, want accepted", sql, err)
		}
	}

	for _, sql := range []string{
		// Not the batch's first statement: inside a batch a name is no
		// statement.
		"SELECT 1; my_proc",
		// A semicolon before it: SQL Server refuses `; my_proc`.
		"; my_proc",
		// A variable the batch does not declare (Msg 137), in the call's
		// arguments or in a statement after it.
		"my_proc @a",
		"my_proc @a = @b",
		"my_proc @a OUTPUT",
		"my_proc SET @a = 1",
		"my_proc; SELECT @a",
		"my_proc FETCH c INTO @a",
		"my_proc; SET @a = 1; DECLARE @a INT",
		// More than four parts are no procedure name.
		"a.b.c.d.e",
		"a.b.c.d.e 1",
		// The words optSimpleExecute's predicate keeps (TSql160.g:572).
		"GET 1",
		"ENABLE 1",
		"DISABLE 1",
		"MOVE 1",
		"SEND 1",
		"RECEIVE 1",
		"THROW 1",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q): accepted, want rejected", sql)
		} else if _, ok := err.(*ParseError); !ok {
			t.Errorf("Parse(%q): %T, want *ParseError", sql, err)
		}
	}
}

// TestImpliedExecCompletion holds completion after a call by name to
// completion after the same call with EXEC: the call is the batch's first
// statement, and Collect reads it as one.
func TestImpliedExecCompletion(t *testing.T) {
	for _, tt := range []struct{ exec, implied string }{
		{"EXEC foo @p = ", "foo @p = "},
		{"EXEC foo 1, @p = ", "foo 1, @p = "},
		{"EXEC foo ", "foo "},
	} {
		exec, implied := Collect(tt.exec, len(tt.exec)), Collect(tt.implied, len(tt.implied))
		if len(implied.Rules) != len(exec.Rules) || len(implied.Tokens) != len(exec.Tokens) {
			t.Errorf("Collect(%q) = %d tokens, %d rules; Collect(%q) = %d tokens, %d rules",
				tt.implied, len(implied.Tokens), len(implied.Rules),
				tt.exec, len(exec.Tokens), len(exec.Rules))
		}
	}
}
