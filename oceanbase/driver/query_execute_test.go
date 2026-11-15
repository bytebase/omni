package driver

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/driver"
	oceanbaseanalysis "github.com/bytebase/omni/oceanbase/analysis"
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

	// 1. Execute: transactional commit with ReadCommitted, rollback on error,
	// auto-commit partial commit, and DELIMITER procedure.
	txLog := &stepLogger{}
	const setupScript = "CREATE TABLE items (id INT AUTO_INCREMENT PRIMARY KEY, note VARCHAR(64));\nINSERT INTO items VALUES (1, 'a'), (2, 'b'), (3, 'c');\nUPDATE items SET note = 'z' WHERE id <= 2;"
	aff, err := d.Execute(ctx, setupScript, driver.ExecuteOptions{
		Transaction: driver.Transactional,
		Isolation:   driver.ReadCommitted,
		Log:         txLog,
	})
	if err != nil || aff != 5 {
		t.Fatalf("Execute tx = (%d, %v)", aff, err)
	}
	wantSteps := []string{
		"tx:BEGIN",
		"exec[0:73]", "resp:0:[0]",
		"exec[73:128]", "resp:3:[3]",
		"exec[128:171]", "resp:2:[2]",
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

	if _, err := d.Execute(ctx, "DELIMITER //\nCREATE PROCEDURE sp_items_count(OUT n INT)\nBEGIN\n  SELECT COUNT(*) INTO n FROM items;\nEND//\nDELIMITER ;\n", driver.ExecuteOptions{
		Transaction: driver.AutoCommit,
	}); err != nil {
		t.Fatalf("DELIMITER Execute failed: %v", err)
	}

	// 2. CountAffectedRows, Explain, ExplainJSON, and ColumnExtra.
	est, err := d.CountAffectedRows(ctx, "DELETE FROM items WHERE id <= 2;")
	if err != nil || est < 1 {
		t.Fatalf("CountAffectedRows = (%d, %v)", est, err)
	}
	if err := d.Explain(ctx, "UPDATE items SET note = 'x' WHERE id = 1"); err != nil {
		t.Fatalf("Explain dry run failed: %v", err)
	}
	planJSON, err := d.ExplainJSON(ctx, "DELETE FROM items WHERE id = 1")
	if err != nil || !strings.Contains(planJSON, "EST.ROWS") {
		t.Fatalf("ExplainJSON = %q, %v", planJSON, err)
	}
	extra, err := d.ColumnExtra(ctx, "omni_qe_test", "items", "id")
	if err != nil || !strings.Contains(strings.ToLower(extra), "auto_increment") {
		t.Fatalf("ColumnExtra = %q, %v", extra, err)
	}

	// 3. Session.Query: types, limit rewrite, EXPLAIN FORMAT=JSON, truncation,
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
		CAST(0xdeadbeef AS BINARY(4)) AS bin,
		NULL AS nul;`, driver.QueryOptions{})
	if err != nil || len(res) != 1 || res[0].Err != nil {
		t.Fatalf("types query = %+v, %v", res, err)
	}
	row := res[0].Rows[0]
	if row[0].Int64 != 42 || row[2].String != "hello" || row[3].Kind != driver.Timestamp || !bytes.Equal(row[4].Bytes, []byte{0xde, 0xad, 0xbe, 0xef}) || row[5].Kind != driver.Null {
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
	var obPlan map[string]any
	var planBody strings.Builder
	for _, r := range jsonPlan[0].Rows {
		planBody.WriteString(r[0].String)
	}
	if err := json.Unmarshal([]byte(planBody.String()), &obPlan); err != nil || len(obPlan) == 0 {
		t.Fatalf("OceanBase json plan = %q, err = %v", planBody.String(), err)
	}

	if _, err := sess.Query(ctx, "SHOW TABLES;", driver.QueryOptions{Explain: true}); !errors.Is(err, oceanbaseanalysis.ErrNotExplainable) {
		t.Errorf("SHOW TABLES explain err = %v, want ErrNotExplainable", err)
	}
	if _, err := sess.Query(ctx, "EXPLAIN ANALYZE SELECT 1;", driver.QueryOptions{Explain: true}); err == nil {
		t.Error("EXPLAIN ANALYZE succeeded")
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
