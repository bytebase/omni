package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/v1review"
)

func TestStatementDMLDryRun(t *testing.T) {
	var stmts []ast.Node
	for _, n := range parse(t, "INSERT INTO t VALUES (1); SELECT 1 FROM dual; UPDATE t SET a = 1; MERGE INTO t USING s ON (t.id = s.id) WHEN MATCHED THEN UPDATE SET t.a = 1; DELETE FROM t; CREATE TABLE u (a INT)") {
		// Bytebase hands the rules each statement unwrapped.
		if raw, ok := n.(*ast.RawStmt); ok {
			n = raw.Stmt
		}
		stmts = append(stmts, n)
	}
	if diff := cmp.Diff([]int{0, 2, 3, 4}, StatementDMLDryRunProbes(stmts)); diff != "" {
		t.Errorf("probes (-want +got):\n%s", diff)
	}
	got := StatementDMLDryRun(stmts, map[int]v1review.Probe{0: {}, 3: {Error: "ORA-00942"}, 1: {Error: "not probed"}})
	want := []v1review.Finding{{Code: v1review.DryRunFailed, Value: "ORA-00942", Statement: 3}}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("findings (-want +got):\n%s", diff)
	}
}
