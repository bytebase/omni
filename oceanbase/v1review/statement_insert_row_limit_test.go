package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/oceanbase/parser"
	"github.com/bytebase/omni/v1review"
)

// rowFinding is a row-limit finding as the tests compare it.
type rowFinding struct {
	Code      v1review.Code
	Value     string
	Statement int
	Text      string
}

func rowFindings(sql string, findings []v1review.Finding) []rowFinding {
	var got []rowFinding
	for _, f := range findings {
		got = append(got, rowFinding{f.Code, f.Value, f.Statement, sql[f.Start:f.End]})
	}
	return got
}

func TestStatementInsertRowLimit(t *testing.T) {
	sql := "INSERT INTO t VALUES (1), (2), (3); INSERT INTO t SELECT * FROM s; INSERT INTO t SELECT * FROM u; INSERT INTO t SELECT * FROM v"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if diff := cmp.Diff([]int{1, 2, 3}, StatementInsertRowLimitProbes(list.Items)); diff != "" {
		t.Errorf("probes (-want +got):\n%s", diff)
	}
	got := rowFindings(sql, StatementInsertRowLimit(list.Items, map[int]v1review.Probe{
		1: {Rows: 5},
		2: {CountError: "no EST.ROWS"},
	}, 2))
	want := []rowFinding{
		{v1review.RowLimitExceeded, "3", 0, "INSERT INTO t VALUES (1), (2), (3)"},
		{v1review.RowLimitExceeded, "5", 1, "INSERT INTO t SELECT * FROM s"},
		{v1review.RowCountFailed, "no EST.ROWS", 2, "INSERT INTO t SELECT * FROM u"},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("findings (-want +got):\n%s", diff)
	}
}

func TestStatementAffectedRowLimit(t *testing.T) {
	sql := "UPDATE t SET a = 1; INSERT INTO t VALUES (1); DELETE FROM u; DELETE FROM v"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if diff := cmp.Diff([]int{0, 2, 3}, StatementAffectedRowLimitProbes(list.Items)); diff != "" {
		t.Errorf("probes (-want +got):\n%s", diff)
	}
	got := rowFindings(sql, StatementAffectedRowLimit(list.Items, map[int]v1review.Probe{
		0: {Error: "explain failed"},
		2: {Rows: 3},
		3: {Rows: 1},
	}, 2))
	want := []rowFinding{
		{v1review.DryRunFailed, "explain failed", 0, "UPDATE t SET a = 1"},
		{v1review.RowLimitExceeded, "3", 2, "DELETE FROM u"},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("findings (-want +got):\n%s", diff)
	}
}
