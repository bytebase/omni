package driver

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/driver"
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

	// 1. Execute: transactional commit, rollback on a failing statement, and
	// auto-commit keeping what ran before the failure.
	txLog := &stepLogger{}
	const setupScript = "CREATE TABLE items (id BIGINT NOT NULL, note VARCHAR(64) NULL) ENGINE=OLAP UNIQUE KEY(id) DISTRIBUTED BY HASH(id) BUCKETS 1 PROPERTIES(\"replication_num\" = \"1\");\nINSERT INTO items VALUES (1, 'a'), (2, 'b'), (3, 'c');\nUPDATE items SET note = 'z' WHERE id <= 2;"
	aff, err := d.Execute(ctx, setupScript, driver.ExecuteOptions{
		Transaction: driver.Transactional,
		Isolation:   driver.RepeatableRead,
		Log:         txLog,
	})
	if err != nil {
		t.Fatalf("Execute tx = (%d, %v)", aff, err)
	}
	if aff != 5 {
		t.Fatalf("Execute tx affected %d rows, want 5", aff)
	}
	if txLog.steps[0] != "tx:BEGIN" || txLog.steps[len(txLog.steps)-1] != "tx:COMMIT" {
		t.Fatalf("tx steps = %v, want a BEGIN first and a COMMIT last", txLog.steps)
	}
	if got := strings.Count(strings.Join(txLog.steps, "|"), "exec["); got != 3 {
		t.Fatalf("tx steps = %v, want 3 commands", txLog.steps)
	}
	rbLog := &stepLogger{}
	if _, err := d.Execute(ctx, "INSERT INTO items VALUES (10, 'rb');\nINSERT INTO missing_items VALUES (11, 'rb');", driver.ExecuteOptions{
		Transaction: driver.Transactional,
		Log:         rbLog,
	}); err == nil {
		t.Fatal("expected duplicate key error")
	}
	if !slices.Contains(rbLog.steps, "tx:ROLLBACK") {
		t.Fatalf("missing rollback: %v", rbLog.steps)
	}

	acLog := &stepLogger{}
	if _, err := d.Execute(ctx, "INSERT INTO items VALUES (20, 'kept');\nINSERT INTO missing_items VALUES (21, 'lost');", driver.ExecuteOptions{
		Transaction: driver.AutoCommit,
		Log:         acLog,
	}); err == nil {
		t.Fatal("expected missing table error")
	}

	// 2. Session.Query: scalar, temporal and binary types, limit rewrite,
	// EXPLAIN, truncation, MySQLError wrapping, DML routed to Exec, and a
	// script stopping at its first failing statement.
	sess, err := d.Session(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close()

	countRows := func(query string) int64 {
		t.Helper()
		res, err := sess.Query(ctx, query, driver.QueryOptions{})
		if err != nil || len(res) != 1 || res[0].Err != nil || len(res[0].Rows) == 0 {
			t.Fatalf("%s = %+v, %v", query, res, err)
		}
		return res[0].Rows[0][0].Int64
	}
	// Doris accepts no DDL inside a transaction and commits each data change, so
	// the failed transactional script leaves its row: 3 committed + 1 from it + 1
	// from the failed auto-commit script.
	if got := countRows("SELECT COUNT(*) FROM items;"); got != 5 {
		t.Fatalf("items after the three scripts = %d rows, want 5", got)
	}

	res, err := sess.Query(ctx, `SELECT
		42 AS i,
		3.5 AS f,
		'hello' AS s,
		CAST('2026-03-04 05:06:07.123456' AS DATETIME(6)) AS dt,
		CAST('' AS VARCHAR) AS empty_s,
		NULL AS nul;`, driver.QueryOptions{})
	if err != nil || len(res) != 1 || res[0].Err != nil {
		t.Fatalf("types query = %+v, %v", res, err)
	}
	row := res[0].Rows[0]
	if row[0].Int64 != 42 || row[2].String != "hello" || row[3].Kind != driver.Timestamp || row[5].Kind != driver.Null {
		t.Fatalf("decoded row = %+v", row)
	}
	if !bytes.Equal([]byte(row[3].Time.Format("2006-01-02 15:04:05.000000")), []byte("2026-03-04 05:06:07.123456")) {
		t.Fatalf("decoded timestamp = %v", row[3].Time)
	}

	limRes, err := sess.Query(ctx, "SELECT id FROM items WHERE id > 0 ORDER BY id;", driver.QueryOptions{Limit: 2})
	if err != nil || len(limRes[0].Rows) != 2 {
		t.Fatalf("limit query = %+v, %v", limRes, err)
	}
	if got := limRes[0].Statement; got != "SELECT id FROM items WHERE id > 0 ORDER BY id LIMIT 2;" {
		t.Fatalf("limit statement = %q", got)
	}

	plan, err := sess.Query(ctx, "SELECT note FROM items WHERE id > 1;", driver.QueryOptions{Explain: true})
	if err != nil || len(plan) != 1 || plan[0].Err != nil || len(plan[0].Rows) == 0 {
		t.Fatalf("explain = %+v, %v", plan, err)
	}
	if got := plan[0].Statement; got != "EXPLAIN SELECT note FROM items WHERE id > 1;" {
		t.Fatalf("explain statement = %q", got)
	}
	if plan[0].Plan == nil || plan[0].Plan.Format != "TEXT" || plan[0].Plan.Executed {
		t.Fatalf("plan metadata = %+v, want a TEXT plan that did not run the statement", plan[0].Plan)
	}

	// DML returns an affected-row count rather than rows.
	dmlRes, err := sess.Query(ctx, "UPDATE items SET note = 'x' WHERE id = 1;", driver.QueryOptions{})
	if err != nil || len(dmlRes) != 1 || dmlRes[0].AffectedRows == nil || len(dmlRes[0].Rows) != 0 {
		t.Fatalf("DML result = %+v, %v", dmlRes, err)
	}

	// A rewrite that fails runs no statement at all: the plan of the first
	// statement is not reported, and no request reaches the server.
	if _, err := sess.Query(ctx, "SELECT 1;\nSELECT FROM WHERE;", driver.QueryOptions{Explain: true}); err == nil {
		t.Fatal("expected an explain rewrite error")
	}

	// A runtime error wraps *mysql.MySQLError and stops the script.
	errRes, err := sess.Query(ctx, "SELECT 1;\nSELECT * FROM missing_tbl;\nSELECT 2;", driver.QueryOptions{})
	if err != nil || len(errRes) != 2 {
		t.Fatalf("runtime error script = %+v, %v", errRes, err)
	}
	// Doris reports a missing table as its generic internal error, not MySQL's 1146.
	var myErr *gomysql.MySQLError
	if !errors.As(errRes[1].Err, &myErr) || myErr.Number != 1105 {
		t.Fatalf("errRes[1].Err = %v, want MySQLError 1105", errRes[1].Err)
	}

	// MaxResultBytes truncates and continues to the next statement.
	truncRes, err := sess.Query(ctx, "SELECT id, repeat('x', 80) FROM items ORDER BY id;\nSELECT 7;", driver.QueryOptions{MaxResultBytes: 200})
	if err != nil || len(truncRes) != 2 || !truncRes[0].Truncated || truncRes[1].Rows[0][0].Int64 != 7 {
		t.Fatalf("truncRes = %+v, %v", truncRes, err)
	}
}
