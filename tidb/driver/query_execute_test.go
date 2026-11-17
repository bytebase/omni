package driver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/driver"
	tidbanalysis "github.com/bytebase/omni/tidb/analysis"
)

type stepLogger struct {
	steps []string
}

func (l *stepLogger) CommandExecute(cmd driver.Command) {
	l.steps = append(l.steps, fmt.Sprintf("exec[%d:%d]", cmd.ByteStart, cmd.ByteEnd))
}

func (l *stepLogger) CommandResponse(affected int64, perResult []int64, err error) {
	if err != nil {
		l.steps = append(l.steps, "resp:err")
		return
	}
	l.steps = append(l.steps, fmt.Sprintf("resp:%d:%v", affected, perResult))
}

func (l *stepLogger) TransactionControl(op driver.TransactionOp, err error) {
	if err != nil {
		l.steps = append(l.steps, fmt.Sprintf("tx:%s:err", op))
		return
	}
	l.steps = append(l.steps, fmt.Sprintf("tx:%s", op))
}

func (l *stepLogger) Retry(_ error, attempt, max int) {
	l.steps = append(l.steps, fmt.Sprintf("retry:%d/%d", attempt, max))
}

func TestQueryExecuteAndStats(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t, "omni_qe_test")

	// 1. Execute: transactional commit, rollback on error, and auto-commit partial commit.
	txLog := &stepLogger{}
	const setupScript = "CREATE TABLE items (id INT PRIMARY KEY, note VARCHAR(64));\nINSERT INTO items VALUES (1, 'a'), (2, 'b'), (3, 'c');\nUPDATE items SET note = 'z' WHERE id <= 2;"
	aff, err := d.Execute(ctx, setupScript, driver.ExecuteOptions{
		Transaction: driver.Transactional,
		Log:         txLog,
	})
	if err != nil || aff != 5 {
		t.Fatalf("Execute tx = (%d, %v)", aff, err)
	}
	wantSteps := []string{
		"tx:BEGIN",
		"exec[0:58]", "resp:0:[0]",
		"exec[58:113]", "resp:3:[3]",
		"exec[113:156]", "resp:2:[2]",
		"tx:COMMIT",
	}
	if !slices.Equal(txLog.steps, wantSteps) {
		t.Fatalf("tx steps = %v, want %v", txLog.steps, wantSteps)
	}

	rbLog := &stepLogger{}
	if _, err := d.Execute(ctx, "INSERT INTO items VALUES (10, 'rb');\nINSERT INTO items VALUES (1, 'dup');", driver.ExecuteOptions{
		Transaction: driver.Transactional,
		Log:         rbLog,
	}); err == nil {
		t.Fatal("expected duplicate PK error")
	}
	if !slices.Contains(rbLog.steps, "tx:ROLLBACK") {
		t.Fatalf("missing rollback: %v", rbLog.steps)
	}

	if _, err := d.Execute(ctx, "INSERT INTO items VALUES (10, 'kept');\nINSERT INTO items VALUES (1, 'dup');", driver.ExecuteOptions{
		Transaction: driver.AutoCommit,
	}); err == nil {
		t.Fatal("expected duplicate PK error in auto-commit")
	}

	// 2. CountAffectedRows, Explain, and DryRunAutoCommit (TiDB BATCH DRY RUN).
	est, err := d.CountAffectedRows(ctx, "DELETE FROM items WHERE id <= 2;")
	if err != nil || est < 1 {
		t.Fatalf("CountAffectedRows = (%d, %v)", est, err)
	}
	for _, statement := range []string{"SELECT 1;", "SELECT * FROM items;"} {
		if _, err := d.CountAffectedRows(ctx, statement); err == nil || !strings.Contains(err.Error(), "expected an INSERT, UPDATE or DELETE") {
			t.Errorf("CountAffectedRows(%q) err = %v", statement, err)
		}
	}
	if _, err := d.CountAffectedRows(ctx, "DELETE FROM items WHERE id <= 2; SELECT 1;"); err == nil || !strings.Contains(err.Error(), "exactly one statement") {
		t.Errorf("CountAffectedRows of a script err = %v", err)
	}
	if err := d.Explain(ctx, "UPDATE items SET note = 'x' WHERE id = 1"); err != nil {
		t.Fatalf("Explain dry run failed: %v", err)
	}
	if err := d.DryRunAutoCommit(ctx, "BATCH ON id LIMIT 2 DRY RUN DELETE FROM items WHERE id > 0;"); err != nil {
		t.Fatalf("DryRunAutoCommit BATCH DRY RUN failed: %v", err)
	}

	// 3. Session.Query: types, limit rewrite, EXPLAIN FORMAT=tidb_json, truncation,
	// pre-validation, and runtime error wrapping.
	sess, err := d.Session(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close()

	res, err := sess.Query(ctx, `SELECT
		42 AS i,
		3.5 AS f,
		'hello' AS s,
		CAST('2026-03-04 05:06:07.123456' AS DATETIME(6)) AS dt,
		NULL AS nul;`, driver.QueryOptions{})
	if err != nil || len(res) != 1 || res[0].Err != nil {
		t.Fatalf("types query = %+v, %v", res, err)
	}
	emptyBytes := []byte{}
	if got := convertValue("BIT", nil, &emptyBytes); got.Kind != driver.Bytes || len(got.Bytes) != 0 {
		t.Fatalf("convertValue empty []byte = %+v, want empty Bytes", got)
	}
	var nilBytes []byte
	if got := convertValue("BIT", nil, &nilBytes); got.Kind != driver.Null {
		t.Fatalf("convertValue nil []byte = %+v, want Null", got)
	}
	row := res[0].Rows[0]
	if row[0].Int64 != 42 || row[2].String != "hello" || row[3].Kind != driver.Timestamp || row[4].Kind != driver.Null {
		t.Fatalf("decoded row = %+v", row)
	}

	limRes, err := sess.Query(ctx, "SELECT id FROM items ORDER BY id;", driver.QueryOptions{Limit: 2})
	if err != nil || len(limRes[0].Rows) != 2 {
		t.Fatalf("limit query = %+v, %v", limRes, err)
	}

	jsonPlan, err := sess.Query(ctx, "SELECT note FROM items WHERE id > 1", driver.QueryOptions{
		Explain:       true,
		ExplainFormat: "JSON",
	})
	if err != nil || len(jsonPlan) != 1 || jsonPlan[0].Plan == nil || jsonPlan[0].Plan.Format != "JSON" {
		t.Fatalf("json explain = %+v, %v", jsonPlan, err)
	}
	var tidbPlan []map[string]any
	if err := json.Unmarshal([]byte(jsonPlan[0].Rows[0][0].String), &tidbPlan); err != nil || len(tidbPlan) == 0 {
		t.Fatalf("tidb_json body = %q, err = %v", jsonPlan[0].Rows[0][0].String, err)
	}

	if _, err := sess.Query(ctx, "EXPLAIN ANALYZE SELECT 1;", driver.QueryOptions{Explain: true}); !errors.Is(err, tidbanalysis.ErrExplainAnalyzeExecutes) {
		t.Errorf("EXPLAIN ANALYZE err = %v", err)
	}

	// Pre-validation rejects syntax errors before running earlier statements.
	if _, err := sess.Query(ctx, "INSERT INTO items VALUES (99, 'never'); SELECT FROM WHERE;", driver.QueryOptions{}); err == nil {
		t.Fatal("expected syntax error")
	}
	chkRes, err := sess.Query(ctx, "SELECT COUNT(*) FROM items WHERE id = 99;", driver.QueryOptions{})
	if err != nil || chkRes[0].Rows[0][0].Int64 != 0 {
		t.Fatalf("row 99 ran before syntax error: %+v, %v", chkRes, err)
	}

	// Runtime error wraps *mysql.MySQLError and stops the script.
	errRes, err := sess.Query(ctx, "SELECT 1;\nSELECT * FROM missing_tbl;\nSELECT 2;", driver.QueryOptions{})
	if err != nil || len(errRes) != 2 {
		t.Fatalf("runtime error script = %+v, %v", errRes, err)
	}
	var myErr *gomysql.MySQLError
	if !errors.As(errRes[1].Err, &myErr) || myErr.Number != 1146 {
		t.Fatalf("errRes[1].Err = %v, want MySQLError 1146", errRes[1].Err)
	}

	// MaxResultBytes truncates and continues to the next statement.
	truncRes, err := sess.Query(ctx, "SELECT id, repeat('x', 80) FROM items ORDER BY id;\nSELECT 7;", driver.QueryOptions{MaxResultBytes: 120})
	if err != nil || len(truncRes) != 2 || !truncRes[0].Truncated || len(truncRes[0].Rows) >= 4 || truncRes[1].Rows[0][0].Int64 != 7 {
		t.Fatalf("truncRes = %+v, %v", truncRes, err)
	}
}
