package driver

import (
	"context"
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/driver"
)

type commandResponse struct {
	affected int64
	per      []int64
	err      error
}

type responseLogger struct {
	responses []commandResponse
}

func (*responseLogger) CommandExecute(driver.Command) {}
func (l *responseLogger) CommandResponse(affected int64, per []int64, err error) {
	l.responses = append(l.responses, commandResponse{affected, per, err})
}
func (*responseLogger) TransactionControl(driver.TransactionOp, error) {}
func (*responseLogger) Retry(error, int, int)                          {}

func counts(results []*driver.Result) []int64 {
	out := make([]int64, len(results))
	for i, result := range results {
		if result.AffectedRows != nil {
			out[i] = *result.AffectedRows
		} else {
			out[i] = -1
		}
	}
	return out
}

// TestQueryBatchAffectedRows pins a multi-statement batch's per-statement
// counts, a batch that fails after some statements ran, and NOCOUNT, which
// reports no count at all rather than zero.
func TestQueryBatchAffectedRows(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t)
	execCorpus(t, d.db, "CREATE TABLE r_items (id INT PRIMARY KEY, note NVARCHAR(20));")
	sess, err := d.Session(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close()

	results, err := sess.Query(ctx, "INSERT INTO r_items VALUES (1, N'a'), (2, N'b'), (3, N'c');\nUPDATE r_items SET note = N'u' WHERE id <= 2;\nDELETE FROM r_items WHERE id = 3;", driver.QueryOptions{})
	if err != nil || len(results) != 3 {
		t.Fatalf("three statements = %+v, %v", results, err)
	}
	if got := counts(results); !slices.Equal(got, []int64{3, 2, 1}) {
		t.Errorf("counts = %v, want [3 2 1]", got)
	}

	results, err = sess.Query(ctx, "INSERT INTO r_items VALUES (10, N'x'), (11, N'y');\nUPDATE r_items SET note = N'z' WHERE id = 10;\nINSERT INTO r_items VALUES (10, N'duplicate');\nDELETE FROM r_items WHERE id = 11;", driver.QueryOptions{})
	if err != nil || len(results) != 3 {
		t.Fatalf("failing batch = %+v, %v", results, err)
	}
	if got := counts(results); !slices.Equal(got, []int64{2, 1, -1}) {
		t.Errorf("counts before the failure = %v, want [2 1 -1]", got)
	}
	if results[2].Err == nil || !strings.Contains(results[2].Err.Error(), "PRIMARY KEY") {
		t.Errorf("failing statement = %v, want the duplicate key error", results[2].Err)
	}

	results, err = sess.Query(ctx, "SET NOCOUNT ON;\nINSERT INTO r_items VALUES (20, N'n'), (21, N'm');\nUPDATE r_items SET note = N'q' WHERE id >= 20;", driver.QueryOptions{})
	if err != nil || len(results) != 3 {
		t.Fatalf("NOCOUNT batch = %+v, %v", results, err)
	}
	if got := counts(results); !slices.Equal(got, []int64{-1, -1, -1}) {
		t.Errorf("NOCOUNT counts = %v, want no count for every statement", got)
	}
}

// TestExecuteBatchAffectedRows pins the logger's per-statement counts and that
// the counts before a failure are kept beside its error.
func TestExecuteBatchAffectedRows(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t)
	execCorpus(t, d.db, "CREATE TABLE x_items (id INT PRIMARY KEY, note NVARCHAR(20));")

	log := &responseLogger{}
	affected, err := d.Execute(ctx, "INSERT INTO x_items VALUES (1, N'a'), (2, N'b'), (3, N'c');\nUPDATE x_items SET note = N'u' WHERE id <= 2;\nDELETE FROM x_items WHERE id = 3;", driver.ExecuteOptions{Transaction: driver.Transactional, Log: log})
	if err != nil || affected != 6 {
		t.Fatalf("execute = %d, %v", affected, err)
	}
	if len(log.responses) != 1 || log.responses[0].affected != 6 || !slices.Equal(log.responses[0].per, []int64{3, 2, 1}) || log.responses[0].err != nil {
		t.Fatalf("response = %+v, want 6 over [3 2 1]", log.responses)
	}

	log = &responseLogger{}
	affected, err = d.Execute(ctx, "INSERT INTO x_items VALUES (10, N'x'), (11, N'y');\nUPDATE x_items SET note = N'z' WHERE id = 10;\nINSERT INTO x_items VALUES (10, N'duplicate');", driver.ExecuteOptions{Transaction: driver.Transactional, Log: log})
	if err == nil || affected != 0 {
		t.Fatalf("failing execute = %d, %v", affected, err)
	}
	if len(log.responses) != 1 || log.responses[0].affected != 3 || !slices.Equal(log.responses[0].per, []int64{2, 1}) || log.responses[0].err == nil {
		t.Fatalf("failing response = %+v, want 3 over [2 1] and an error", log.responses)
	}
	if !strings.Contains(log.responses[0].err.Error(), "PRIMARY KEY") {
		t.Errorf("logged error = %v, want the duplicate key error", log.responses[0].err)
	}
}
