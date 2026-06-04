package parser

import (
	"testing"
)

// This file is the parser-dcl-tcl node's correctness gate for the TCL
// statements (START TRANSACTION / COMMIT / ROLLBACK), structural layer. The
// authoritative accept/reject differential against the live Trino 481 oracle
// lives in oracle_dcl_tcl_test.go alongside the DCL/prepared corpora.

// dclParseOne parses exactly one statement via the public Parse entry point and
// returns it, failing the test if parsing errored or did not yield exactly one
// statement. It is the shared structural-test helper for the parser-dcl-tcl
// node (used by the transaction / grant-revoke / prepared structural tests).
func dclParseOne(t *testing.T, sql string) interface{} {
	t.Helper()
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse(%q): unexpected errors: %v", sql, errs)
	}
	if file == nil || len(file.Stmts) != 1 {
		got := 0
		if file != nil {
			got = len(file.Stmts)
		}
		t.Fatalf("Parse(%q): got %d statements, want 1", sql, got)
	}
	return file.Stmts[0]
}

// dclParseErr asserts that Parse reports at least one error for sql (a
// rejected-by-the-grammar input). It is the structural-layer negative-case
// helper; the oracle differential confirms these rejections match Trino.
func dclParseErr(t *testing.T, sql string) {
	t.Helper()
	_, errs := Parse(sql)
	if len(errs) == 0 {
		t.Errorf("Parse(%q): want at least one error, got none", sql)
	}
}

func TestTransaction_Structure(t *testing.T) {
	t.Run("start_bare", func(t *testing.T) {
		stmt, ok := dclParseOne(t, "START TRANSACTION").(*StartTransactionStmt)
		if !ok {
			t.Fatalf("got %T, want *StartTransactionStmt", stmt)
		}
		if len(stmt.Modes) != 0 {
			t.Errorf("got %d modes, want 0", len(stmt.Modes))
		}
	})

	t.Run("start_isolation_repeatable_read", func(t *testing.T) {
		stmt := dclParseOne(t, "START TRANSACTION ISOLATION LEVEL REPEATABLE READ").(*StartTransactionStmt)
		if len(stmt.Modes) != 1 {
			t.Fatalf("got %d modes, want 1", len(stmt.Modes))
		}
		m := stmt.Modes[0]
		if m.IsAccessMode {
			t.Errorf("mode IsAccessMode = true, want false (isolation)")
		}
		if m.Isolation != IsolationRepeatableRead {
			t.Errorf("isolation = %v, want RepeatableRead", m.Isolation)
		}
	})

	t.Run("start_read_write", func(t *testing.T) {
		stmt := dclParseOne(t, "START TRANSACTION READ WRITE").(*StartTransactionStmt)
		if len(stmt.Modes) != 1 {
			t.Fatalf("got %d modes, want 1", len(stmt.Modes))
		}
		m := stmt.Modes[0]
		if !m.IsAccessMode || m.ReadOnly {
			t.Errorf("got IsAccessMode=%v ReadOnly=%v, want access-mode READ WRITE", m.IsAccessMode, m.ReadOnly)
		}
	})

	t.Run("start_two_modes", func(t *testing.T) {
		stmt := dclParseOne(t, "START TRANSACTION ISOLATION LEVEL READ COMMITTED, READ ONLY").(*StartTransactionStmt)
		if len(stmt.Modes) != 2 {
			t.Fatalf("got %d modes, want 2", len(stmt.Modes))
		}
		if stmt.Modes[0].IsAccessMode || stmt.Modes[0].Isolation != IsolationReadCommitted {
			t.Errorf("mode[0] = %+v, want isolation READ COMMITTED", stmt.Modes[0])
		}
		if !stmt.Modes[1].IsAccessMode || !stmt.Modes[1].ReadOnly {
			t.Errorf("mode[1] = %+v, want access READ ONLY", stmt.Modes[1])
		}
	})

	t.Run("start_serializable", func(t *testing.T) {
		stmt := dclParseOne(t, "START TRANSACTION ISOLATION LEVEL SERIALIZABLE").(*StartTransactionStmt)
		if stmt.Modes[0].Isolation != IsolationSerializable {
			t.Errorf("isolation = %v, want Serializable", stmt.Modes[0].Isolation)
		}
	})

	t.Run("start_read_uncommitted", func(t *testing.T) {
		stmt := dclParseOne(t, "START TRANSACTION ISOLATION LEVEL READ UNCOMMITTED").(*StartTransactionStmt)
		if stmt.Modes[0].Isolation != IsolationReadUncommitted {
			t.Errorf("isolation = %v, want ReadUncommitted", stmt.Modes[0].Isolation)
		}
	})

	t.Run("commit_bare", func(t *testing.T) {
		stmt, ok := dclParseOne(t, "COMMIT").(*CommitStmt)
		if !ok {
			t.Fatalf("got %T, want *CommitStmt", stmt)
		}
		if stmt.Work {
			t.Errorf("Work = true, want false")
		}
	})

	t.Run("commit_work", func(t *testing.T) {
		stmt := dclParseOne(t, "COMMIT WORK").(*CommitStmt)
		if !stmt.Work {
			t.Errorf("Work = false, want true")
		}
	})

	t.Run("rollback_bare", func(t *testing.T) {
		stmt, ok := dclParseOne(t, "ROLLBACK").(*RollbackStmt)
		if !ok {
			t.Fatalf("got %T, want *RollbackStmt", stmt)
		}
		if stmt.Work {
			t.Errorf("Work = true, want false")
		}
	})

	t.Run("rollback_work", func(t *testing.T) {
		stmt := dclParseOne(t, "ROLLBACK WORK").(*RollbackStmt)
		if !stmt.Work {
			t.Errorf("Work = false, want true")
		}
	})
}

func TestTransaction_Negative(t *testing.T) {
	// Forms Trino rejects: a bogus isolation level, READ with no access mode,
	// a trailing comma, and ISOLATION without LEVEL. The oracle differential
	// confirms these match Trino 481.
	dclParseErr(t, "START TRANSACTION ISOLATION LEVEL BOGUS")
	dclParseErr(t, "START TRANSACTION READ")
	dclParseErr(t, "START TRANSACTION ISOLATION REPEATABLE READ")
	dclParseErr(t, "START TRANSACTION ISOLATION LEVEL READ COMMITTED,")
}
