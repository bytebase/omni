package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/mariadb/parser"
	"github.com/bytebase/omni/v1review"
)

func TestStatementInsertRowLimit(t *testing.T) {
	list, err := parser.Parse("INSERT INTO t VALUES (1), (2), (3); INSERT INTO t VALUES (1); INSERT INTO t SELECT * FROM s; UPDATE t SET a = 1; INSERT INTO t SELECT * FROM u; INSERT INTO t SELECT * FROM v; INSERT INTO t SELECT * FROM w")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if diff := cmp.Diff([]int{2, 4, 5, 6}, StatementInsertRowLimitProbes(list.Items)); diff != "" {
		t.Errorf("probes (-want +got):\n%s", diff)
	}
	got := StatementInsertRowLimit(list.Items, map[int]v1review.Probe{
		2: {Rows: 5},
		4: {Error: "explain failed"},
		5: {CountError: "no estimate"},
		6: {Rows: 2},
	}, 2)
	want := []v1review.Finding{
		{Code: v1review.RowLimitExceeded, Count: 3, Statement: 0},
		{Code: v1review.RowLimitExceeded, Count: 5, Statement: 2},
		{Code: v1review.DryRunFailed, Detail: "explain failed", Statement: 4},
		{Code: v1review.RowCountFailed, Detail: "no estimate", Statement: 5},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("findings (-want +got):\n%s", diff)
	}
	if got := StatementInsertRowLimit(list.Items, nil, 2); len(got) != 1 || got[0].Statement != 0 {
		t.Errorf("unprobed: got %+v, want only the VALUES insert", got)
	}
}

func TestStatementAffectedRowLimit(t *testing.T) {
	list, err := parser.Parse("UPDATE t SET a = 1; DELETE FROM t; INSERT INTO t VALUES (1); DELETE FROM u; UPDATE v SET a = 1")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if diff := cmp.Diff([]int{0, 1, 3, 4}, StatementAffectedRowLimitProbes(list.Items)); diff != "" {
		t.Errorf("probes (-want +got):\n%s", diff)
	}
	got := StatementAffectedRowLimit(list.Items, map[int]v1review.Probe{
		0: {Rows: 10},
		1: {Rows: 10, Error: "explain failed"},
		3: {Rows: 2},
	}, 2)
	want := []v1review.Finding{
		{Code: v1review.RowLimitExceeded, Count: 10, Statement: 0},
		{Code: v1review.DryRunFailed, Detail: "explain failed", Statement: 1},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("findings (-want +got):\n%s", diff)
	}
}
