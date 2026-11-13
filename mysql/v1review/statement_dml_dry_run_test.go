package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/mysql/parser"
	"github.com/bytebase/omni/v1review"
)

func TestStatementDMLDryRun(t *testing.T) {
	list, err := parser.Parse("INSERT INTO t VALUES (1); SELECT 1; UPDATE t SET a = 1; CREATE TABLE u (a INT); DELETE FROM t; CREATE PROCEDURE p() BEGIN DELETE FROM t; END; REPLACE INTO t VALUES (1)")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if diff := cmp.Diff([]int{0, 2, 4, 6}, StatementDMLDryRunProbes(list.Items)); diff != "" {
		t.Errorf("probes (-want +got):\n%s", diff)
	}
	got := StatementDMLDryRun(list.Items, map[int]v1review.Probe{
		0: {Error: "no such table"},
		2: {},
		4: {Error: "bad column"},
		1: {Error: "a SELECT is never probed"},
	})
	want := []v1review.Finding{
		{Code: v1review.DryRunFailed, Detail: "no such table", Statement: 0},
		{Code: v1review.DryRunFailed, Detail: "bad column", Statement: 4},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("findings (-want +got):\n%s", diff)
	}
}
