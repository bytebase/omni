package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

// TestLockTableClausesStrict checks the lock mode, WAIT interval, and
// partition forms of LOCK TABLE against what Oracle 23ai accepts (reference
// rows ref_140 to ref_146).
func TestLockTableClausesStrict(t *testing.T) {
	for sql, want := range map[string]string{
		"LOCK TABLE t IN ROW SHARE MODE":           "ROW SHARE",
		"LOCK TABLE t IN ROW EXCLUSIVE MODE":       "ROW EXCLUSIVE",
		"LOCK TABLE t IN SHARE UPDATE MODE":        "SHARE UPDATE",
		"LOCK TABLE t IN SHARE MODE":               "SHARE",
		"LOCK TABLE t IN SHARE ROW EXCLUSIVE MODE": "SHARE ROW EXCLUSIVE",
		"LOCK TABLE t IN EXCLUSIVE MODE WAIT 5":    "EXCLUSIVE",
	} {
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.LockTableStmt)
		if stmt.LockMode != want {
			t.Errorf("Parse(%q) LockMode = %q, want %q", sql, stmt.LockMode, want)
		}
	}

	stmt := ParseAndCheck(t, "LOCK TABLE t PARTITION FOR (2023, 'x') IN SHARE MODE").Items[0].(*ast.RawStmt).Stmt.(*ast.LockTableStmt)
	if got := stmt.Tables[0].PartitionName; got != "2023, 'x'" {
		t.Errorf("PARTITION FOR key values = %q, want %q", got, "2023, 'x'")
	}

	for _, bad := range []string{
		"LOCK TABLE t IN FOO MODE",
		"LOCK TABLE t IN ROW MODE",
		"LOCK TABLE t IN MODE",
		"LOCK TABLE t IN SHARE EXCLUSIVE MODE",
		"LOCK TABLE t IN EXCLUSIVE ROW MODE",
		"LOCK TABLE t IN ROW SHARE UPDATE MODE",
		"LOCK TABLE t IN EXCLUSIVE MODE WAIT",
		"LOCK TABLE t IN EXCLUSIVE MODE WAIT -1",
		"LOCK TABLE t IN EXCLUSIVE MODE WAIT 2.5",
		"LOCK TABLE t IN EXCLUSIVE MODE WAIT 1 + 1",
		"LOCK TABLE t PARTITION () IN EXCLUSIVE MODE",
		"LOCK TABLE t PARTITION (p1, p2) IN EXCLUSIVE MODE",
		"LOCK TABLE t PARTITION FOR () IN EXCLUSIVE MODE",
		"BEGIN LOCK TABLE t IN FOO MODE; END;",
		"DECLARE c CONSTANT PLS_INTEGER := 5; BEGIN LOCK TABLE t IN EXCLUSIVE MODE WAIT c; END;",
	} {
		ParseShouldFail(t, bad)
	}
}
