package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/tidb/parser"
	"github.com/bytebase/omni/v1review"
)

func TestStatementDMLDryRun(t *testing.T) {
	list, err := parser.Parse("INSERT INTO t VALUES (1); SELECT 1; BATCH ON id LIMIT 10 DELETE FROM t WHERE a = 1; EXPLAIN DELETE FROM t; UPDATE t SET a = 1; REPLACE INTO t VALUES (1)")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if diff := cmp.Diff([]int{0, 2, 4, 5}, StatementDMLDryRunProbes(list.Items)); diff != "" {
		t.Errorf("probes (-want +got):\n%s", diff)
	}
	got := StatementDMLDryRun(list.Items, map[int]v1review.Probe{2: {Error: "failed to extract inner DML"}, 4: {}})
	want := []v1review.Finding{{Code: v1review.DryRunFailed, Value: "failed to extract inner DML", Statement: 2}}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("findings (-want +got):\n%s", diff)
	}
}
