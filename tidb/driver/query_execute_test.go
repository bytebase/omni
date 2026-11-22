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

// TestExecuteDelimiterAndDDLFlows pins DELIMITER script splitting, CREATE INDEX
// inside an explicit transaction block, and PREPARE/EXECUTE/DEALLOCATE DDL flows.
func TestExecuteDelimiterAndDDLFlows(t *testing.T) {
	ctx := context.Background()
	const dbName = "omni_qe_ddl_flows"
	d := openTestDriver(t, dbName)

	// 1. DELIMITER groups statements without sending the directive to the server.
	if _, err := d.Execute(ctx, `
		DELIMITER //
		CREATE TABLE delimiter_one (id INT)//
		CREATE TABLE delimiter_two (id INT)//
		DELIMITER ;
	`, driver.ExecuteOptions{}); err != nil {
		t.Fatalf("DELIMITER Execute failed: %v", err)
	}
	for _, table := range []string{"delimiter_one", "delimiter_two"} {
		var count int
		if err := d.db.QueryRowContext(ctx, fmt.Sprintf(
			"SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = '%s' AND table_name = '%s'", dbName, table,
		)).Scan(&count); err != nil || count != 1 {
			t.Fatalf("table %s count = (%d, %v), want 1", table, count, err)
		}
	}

	// 2. Explicit BEGIN/CREATE INDEX/COMMIT inside Execute.
	if _, err := d.Execute(ctx, fmt.Sprintf(`
		USE %[1]s;
		CREATE TABLE %[1]s.execute_create_index_in_transaction (id INT);
		BEGIN;
		CREATE INDEX idx_execute_create_index_in_transaction ON %[1]s.execute_create_index_in_transaction(id);
		COMMIT;
	`, dbName), driver.ExecuteOptions{}); err != nil {
		t.Fatalf("CREATE INDEX in tx failed: %v", err)
	}
	var idxCount int
	if err := d.db.QueryRowContext(ctx, fmt.Sprintf(`
		SELECT COUNT(*) FROM information_schema.tidb_indexes
		WHERE table_schema = '%s' AND table_name = 'execute_create_index_in_transaction' AND key_name = 'idx_execute_create_index_in_transaction'
	`, dbName)).Scan(&idxCount); err != nil || idxCount != 1 {
		t.Fatalf("tx index count = (%d, %v), want 1", idxCount, err)
	}

	// 3. PREPARE / EXECUTE / DEALLOCATE dynamic DDL flow.
	if _, err := d.Execute(ctx, fmt.Sprintf(`
		USE %[1]s;
		CREATE TABLE %[1]s.prepare_statement_flow (id INT);
		SET @sql := 'CREATE INDEX idx_prepare_statement_flow ON %[1]s.prepare_statement_flow(id)';
		PREPARE stmt FROM @sql;
		EXECUTE stmt;
		DEALLOCATE PREPARE stmt;
	`, dbName), driver.ExecuteOptions{}); err != nil {
		t.Fatalf("prepared DDL flow failed: %v", err)
	}
	if err := d.db.QueryRowContext(ctx, fmt.Sprintf(`
		SELECT COUNT(*) FROM information_schema.tidb_indexes
		WHERE table_schema = '%s' AND table_name = 'prepare_statement_flow' AND key_name = 'idx_prepare_statement_flow'
	`, dbName)).Scan(&idxCount); err != nil || idxCount != 1 {
		t.Fatalf("prepared index count = (%d, %v), want 1", idxCount, err)
	}
}

// TestQueryExplainAndSessionState pins session variable persistence across
// Query calls, EXPLAIN JSON/TEXT rewriting and plan detection, and parser
// error formatting.
func TestQueryExplainAndSessionState(t *testing.T) {
	ctx := context.Background()
	const dbName = "omni_qe_explain_session"
	d := openTestDriver(t, dbName)
	if _, err := d.db.ExecContext(ctx, "CREATE TABLE plan_target (id INT PRIMARY KEY, note VARCHAR(64))"); err != nil {
		t.Fatal(err)
	}

	sess, err := d.Session(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close()

	// Session variables persist across calls on the same session.
	if res, err := sess.Query(ctx, "SET @omni_session_probe = 42", driver.QueryOptions{}); err != nil || len(res) != 1 || res[0].Err != nil {
		t.Fatalf("SET session var = %+v, %v", res, err)
	}
	probe, err := sess.Query(ctx, "SELECT @omni_session_probe", driver.QueryOptions{})
	if err != nil || len(probe) != 1 || probe[0].Err != nil || len(probe[0].Rows) != 1 || probe[0].Rows[0][0].Int64 != 42 {
		t.Fatalf("SELECT session var = %+v, %v", probe, err)
	}

	const stmt = "SELECT note FROM plan_target WHERE id > 100"
	explain := func(s string, opts driver.QueryOptions) *driver.Result {
		t.Helper()
		res, err := sess.Query(ctx, s, opts)
		if err != nil || len(res) != 1 || res[0].Err != nil {
			t.Fatalf("Query(%q) = %+v, %v", s, res, err)
		}
		return res[0]
	}
	requireJSONRootID := func(res *driver.Result) {
		t.Helper()
		if len(res.Columns) != 1 || len(res.Rows) != 1 {
			t.Fatalf("expected 1x1 JSON plan, got %d cols x %d rows", len(res.Columns), len(res.Rows))
		}
		var plan []struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal([]byte(res.Rows[0][0].String), &plan); err != nil || len(plan) == 0 || plan[0].ID == "" {
			t.Fatalf("tidb_json plan = %q, err = %v", res.Rows[0][0].String, err)
		}
	}

	// JSON explain rewrite.
	jsonPlan := explain(stmt, driver.QueryOptions{Explain: true, ExplainFormat: "JSON"})
	if jsonPlan.Statement != "EXPLAIN FORMAT=tidb_json "+stmt {
		t.Errorf("json explain statement = %q", jsonPlan.Statement)
	}
	requireJSONRootID(jsonPlan)
	if jsonPlan.Plan == nil || jsonPlan.Plan.Format != "JSON" || jsonPlan.Plan.Executed {
		t.Errorf("json explain plan = %+v", jsonPlan.Plan)
	}

	// Default TEXT explain stays tabular.
	textPlan := explain(stmt, driver.QueryOptions{Explain: true})
	if textPlan.Statement != "EXPLAIN "+stmt || len(textPlan.Columns) <= 1 {
		t.Errorf("text explain = stmt %q, %d cols", textPlan.Statement, len(textPlan.Columns))
	}

	// Existing EXPLAIN is replanned in the requested format rather than wrapped.
	replanned := explain("EXPLAIN FORMAT=brief "+stmt, driver.QueryOptions{Explain: true, ExplainFormat: "JSON"})
	if replanned.Statement != "EXPLAIN FORMAT=tidb_json "+stmt {
		t.Errorf("replanned statement = %q", replanned.Statement)
	}
	requireJSONRootID(replanned)

	// User-written EXPLAIN FORMAT=tidb_json and plain EXPLAIN report typed Plan.
	ranJSON := explain("EXPLAIN FORMAT=tidb_json "+stmt, driver.QueryOptions{})
	if ranJSON.Plan == nil || ranJSON.Plan.Format != "JSON" || ranJSON.Plan.Executed {
		t.Errorf("ran JSON plan = %+v", ranJSON.Plan)
	}
	requireJSONRootID(ranJSON)

	ranText := explain("EXPLAIN "+stmt, driver.QueryOptions{})
	if ranText.Plan == nil || ranText.Plan.Format != "TEXT" || ranText.Plan.Executed {
		t.Errorf("ran TEXT plan = %+v", ranText.Plan)
	}

	// Parser error retains message, position, and related text.
	_, err = sess.Query(ctx, "SELECT FROM WHERE", driver.QueryOptions{})
	if err == nil || !strings.Contains(err.Error(), `syntax error at or near "FROM" (line 1, column 8)`) || !strings.Contains(err.Error(), "related text: SELECT FROM WHERE") {
		t.Errorf("parse error = %v", err)
	}
}

// TestProbeAndDryRun pins TiDB's dry-run and session rollback behavior: Explain
// refuses invalid statements, BATCH ... DRY RUN fails inside a transaction and
// succeeds in autocommit mode, and transactional session writes roll back.
func TestProbeAndDryRun(t *testing.T) {
	ctx := context.Background()
	const dbName = "omni_qe_probe"
	d := openTestDriver(t, dbName)
	if _, err := d.db.ExecContext(ctx, "CREATE TABLE t (id INT PRIMARY KEY, note VARCHAR(64))"); err != nil {
		t.Fatal(err)
	}

	if err := d.Explain(ctx, "DELETE FROM t WHERE id = 1"); err != nil {
		t.Fatalf("Explain valid DELETE failed: %v", err)
	}
	if err := d.Explain(ctx, "DELETE FROM t WHERE missing_column = 1"); err == nil {
		t.Fatal("Explain on missing_column succeeded, want error")
	}

	const batchDryRun = "BATCH ON id LIMIT 2 DRY RUN DELETE FROM t WHERE id > 0"
	if err := d.DryRunAutoCommit(ctx, batchDryRun); err != nil {
		t.Fatalf("DryRunAutoCommit BATCH DRY RUN failed: %v", err)
	}

	sess, err := d.Session(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close()

	// BATCH ... DRY RUN is refused inside an explicit transaction.
	if res, err := sess.Query(ctx, "START TRANSACTION", driver.QueryOptions{}); err != nil || len(res) != 1 || res[0].Err != nil {
		t.Fatalf("START TRANSACTION = %+v, %v", res, err)
	}
	batchRes, err := sess.Query(ctx, batchDryRun, driver.QueryOptions{})
	if err == nil && (len(batchRes) == 0 || batchRes[0].Err == nil) {
		t.Fatalf("BATCH DRY RUN inside transaction succeeded: %+v", batchRes)
	}
	if res, err := sess.Query(ctx, "ROLLBACK", driver.QueryOptions{}); err != nil || len(res) != 1 || res[0].Err != nil {
		t.Fatalf("ROLLBACK after BATCH DRY RUN = %+v, %v", res, err)
	}

	// Pre-execution writes inside a probe transaction are visible to the probed
	// query and gone after ROLLBACK.
	for _, stmt := range []string{
		"START TRANSACTION",
		"INSERT INTO t VALUES (42, 'probe')",
	} {
		if res, err := sess.Query(ctx, stmt, driver.QueryOptions{}); err != nil || len(res) != 1 || res[0].Err != nil {
			t.Fatalf("%s = %+v, %v", stmt, res, err)
		}
	}
	inTx, err := sess.Query(ctx, "SELECT COUNT(*) FROM t", driver.QueryOptions{})
	if err != nil || len(inTx) != 1 || inTx[0].Err != nil || len(inTx[0].Rows) != 1 || inTx[0].Rows[0][0].Int64 != 1 {
		t.Fatalf("in-transaction count = %+v, %v", inTx, err)
	}
	if res, err := sess.Query(ctx, "ROLLBACK", driver.QueryOptions{}); err != nil || len(res) != 1 || res[0].Err != nil {
		t.Fatalf("ROLLBACK = %+v, %v", res, err)
	}
	afterRollback, err := sess.Query(ctx, "SELECT COUNT(*) FROM t", driver.QueryOptions{})
	if err != nil || len(afterRollback) != 1 || afterRollback[0].Err != nil || len(afterRollback[0].Rows) != 1 || afterRollback[0].Rows[0][0].Int64 != 0 {
		t.Fatalf("post-rollback count = %+v, %v", afterRollback, err)
	}
}
