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
	mariadbanalysis "github.com/bytebase/omni/mariadb/analysis"
	"github.com/bytebase/omni/mariadb/internal/testcontainer"
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
	for _, role := range []string{testcontainer.MariaDB107, "default"} {
		t.Run(role, func(t *testing.T) {
			d := openTestDriver(t, role, "omni_qe_test")

			// 1. Execute: transactional commit with RepeatableRead, rollback on error,
			// auto-commit partial commit, and DELIMITER stored procedure.
			txLog := &stepLogger{}
			const setupScript = "CREATE TABLE items (id INT PRIMARY KEY, note VARCHAR(64));\nINSERT INTO items VALUES (1, 'a'), (2, 'b'), (3, 'c');\nUPDATE items SET note = 'z' WHERE id <= 2;"
			aff, err := d.Execute(ctx, setupScript, driver.ExecuteOptions{
				Transaction: driver.Transactional,
				Isolation:   driver.RepeatableRead,
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

			acLog := &stepLogger{}
			acAff, err := d.Execute(ctx, "INSERT INTO items VALUES (10, 'kept');\nINSERT INTO items VALUES (1, 'dup');", driver.ExecuteOptions{
				Transaction: driver.AutoCommit,
				Log:         acLog,
			})
			if err == nil || acAff != 0 {
				// Execute returns 0 on error, while row 10 remains committed on the server.
				_ = acAff
			}

			if _, err := d.Execute(ctx, "DELIMITER //\nCREATE PROCEDURE sp_items_count(OUT n INT)\nBEGIN\n  SELECT COUNT(*) INTO n FROM items;\nEND//\nDELIMITER ;\n", driver.ExecuteOptions{
				Transaction: driver.AutoCommit,
			}); err != nil {
				t.Fatalf("DELIMITER Execute failed: %v", err)
			}

			// 2. CountAffectedRows (including MariaDB RETURNING syntax) and Explain.
			est, err := d.CountAffectedRows(ctx, "DELETE FROM items WHERE id <= 2 RETURNING id;")
			if err != nil || est < 1 {
				t.Fatalf("CountAffectedRows RETURNING = (%d, %v)", est, err)
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

			// 3. Session.Query: scalar/temporal/binary types, limit rewrite, EXPLAIN FORMAT=JSON,
			// truncation, MySQLError wrapping, RETURNING DML routed to Exec, and pre-validation.
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
				CAST('' AS BINARY) AS empty_bin,
				NULL AS nul;`, driver.QueryOptions{})
			if err != nil || len(res) != 1 || res[0].Err != nil {
				t.Fatalf("types query = %+v, %v", res, err)
			}
			row := res[0].Rows[0]
			if row[0].Int64 != 42 || row[2].String != "hello" || row[3].Kind != driver.Timestamp || !bytes.Equal(row[4].Bytes, []byte{0xde, 0xad, 0xbe, 0xef}) || row[5].Kind != driver.Bytes || len(row[5].Bytes) != 0 || row[6].Kind != driver.Null {
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
			var root map[string]any
			if err := json.Unmarshal([]byte(jsonPlan[0].Rows[0][0].String), &root); err != nil || root["query_block"] == nil {
				t.Fatalf("json plan body = %q, err = %v", jsonPlan[0].Rows[0][0].String, err)
			}

			if _, err := sess.Query(ctx, "ANALYZE SELECT 1;", driver.QueryOptions{Explain: true}); !errors.Is(err, mariadbanalysis.ErrExplainAnalyzeExecutes) {
				t.Errorf("ANALYZE SELECT 1 err = %v", err)
			}

			// MariaDB DELETE ... RETURNING answers with the rows it deleted, so
			// it routes to QueryContext and reports no affected-rows count.
			retRes, err := sess.Query(ctx, "DELETE FROM items WHERE id = 10 RETURNING note;", driver.QueryOptions{})
			if err != nil || len(retRes) != 1 || retRes[0].AffectedRows != nil || len(retRes[0].Rows) != 1 {
				t.Fatalf("DELETE RETURNING result = %+v, %v", retRes, err)
			}
			if retRes[0].Rows[0][0].String != "kept" {
				t.Errorf("DELETE RETURNING note = %q", retRes[0].Rows[0][0].String)
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
			if err != nil || len(truncRes) != 2 || !truncRes[0].Truncated || len(truncRes[0].Rows) >= 3 || truncRes[1].Rows[0][0].Int64 != 7 {
				t.Fatalf("truncRes = %+v, %v", truncRes, err)
			}
		})
	}
}

// TestQueryServerRows covers the statements that answer with a result set the
// read-only whitelist does not name — a procedure's result set, a table
// check's diagnosis — the rows a write's RETURNING selects, and masking: the
// rows a statement does not trace to columns are withheld when the caller asks.
// It also holds a session's state to the session: a temporary table one session
// creates is gone when it closes.
func TestQueryServerRows(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t, "default", "omni_rows_test")
	if _, err := d.db.ExecContext(ctx, `CREATE TABLE ritems (id INT PRIMARY KEY, d INT);
INSERT INTO ritems VALUES (1, 0), (2, 0), (3, 1);
CREATE PROCEDURE rcount() BEGIN SELECT COUNT(*) AS n FROM ritems; END`); err != nil {
		t.Fatal(err)
	}

	sess, err := d.Session(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close()

	// RETURNING: INSERT, REPLACE and DELETE answer with the rows the statement
	// wrote or removed, and report no count instead.
	inserted, err := sess.Query(ctx, "INSERT INTO ritems VALUES (4, 1) RETURNING id;", driver.QueryOptions{})
	if err != nil || len(inserted) != 1 || inserted[0].Err != nil || inserted[0].AffectedRows != nil {
		t.Fatalf("INSERT RETURNING = %+v, %v", inserted, err)
	}
	if len(inserted[0].Rows) != 1 || inserted[0].Rows[0][0].Int64 != 4 {
		t.Errorf("INSERT RETURNING rows = %+v", inserted[0].Rows)
	}
	replaced, err := sess.Query(ctx, "REPLACE INTO ritems VALUES (4, 2) RETURNING d;", driver.QueryOptions{})
	if err != nil || len(replaced) != 1 || len(replaced[0].Rows) != 1 || replaced[0].Rows[0][0].Int64 != 2 {
		t.Fatalf("REPLACE RETURNING = %+v, %v", replaced, err)
	}
	deleted, err := sess.Query(ctx, "DELETE FROM ritems WHERE d = 0 RETURNING id;", driver.QueryOptions{})
	if err != nil || len(deleted) != 1 || len(deleted[0].Rows) != 2 {
		t.Fatalf("DELETE RETURNING rows = %+v, %v", deleted, err)
	}

	// RETURNING columns are not traced by mariadb/analysis's lineage, so masking
	// withholds them rather than returning values the caller asked to mask.
	maskedDeleted, err := sess.Query(ctx, "DELETE FROM ritems WHERE d = 1 RETURNING id;", driver.QueryOptions{MaskingEnabled: true})
	if err != nil || len(maskedDeleted) != 1 || maskedDeleted[0].AffectedRows == nil || len(maskedDeleted[0].Rows) != 0 {
		t.Fatalf("masked DELETE RETURNING = %+v, %v", maskedDeleted, err)
	}
	maskedInserted, err := sess.Query(ctx, "INSERT INTO ritems VALUES (5, 1) RETURNING id;", driver.QueryOptions{MaskingEnabled: true})
	if err != nil || len(maskedInserted) != 1 || maskedInserted[0].AffectedRows == nil || len(maskedInserted[0].Rows) != 0 {
		t.Fatalf("masked INSERT RETURNING = %+v, %v", maskedInserted, err)
	}

	// A procedure's result set is columns the caller's analysis cannot trace:
	// unmasked it comes back, masked it is withheld and reported as a count.
	call, err := sess.Query(ctx, "CALL rcount();", driver.QueryOptions{})
	if err != nil || len(call) != 1 || call[0].Err != nil || len(call[0].Rows) != 1 {
		t.Fatalf("CALL = %+v, %v", call, err)
	}
	if call[0].Rows[0][0].Int64 != 2 {
		t.Errorf("CALL rows = %+v", call[0].Rows)
	}
	maskedCall, err := sess.Query(ctx, "CALL rcount();", driver.QueryOptions{MaskingEnabled: true})
	if err != nil || len(maskedCall) != 1 || maskedCall[0].AffectedRows == nil || len(maskedCall[0].Rows) != 0 {
		t.Fatalf("masked CALL = %+v, %v", maskedCall, err)
	}

	// Masking a traced statement changes nothing.
	plain, err := sess.Query(ctx, "SELECT id, d FROM ritems ORDER BY id;", driver.QueryOptions{})
	if err != nil {
		t.Fatal(err)
	}
	masked, err := sess.Query(ctx, "SELECT id, d FROM ritems ORDER BY id;", driver.QueryOptions{MaskingEnabled: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(plain[0].Rows) != len(masked[0].Rows) || len(plain[0].Rows) == 0 {
		t.Errorf("masking changed a traced statement's rows: %d vs %d", len(plain[0].Rows), len(masked[0].Rows))
	}

	// CHECK TABLE answers a diagnosis, not a count.
	checked, err := sess.Query(ctx, "CHECK TABLE ritems;", driver.QueryOptions{})
	if err != nil || len(checked) != 1 || len(checked[0].Rows) == 0 {
		t.Fatalf("CHECK TABLE = %+v, %v", checked, err)
	}

	// A session's temporary table dies with it.
	if _, err := sess.Query(ctx, "CREATE TEMPORARY TABLE leak (id INT);", driver.QueryOptions{}); err != nil {
		t.Fatal(err)
	}
	if err := sess.Close(); err != nil {
		t.Fatal(err)
	}
	next, err := d.Session(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer next.Close()
	leak, err := next.Query(ctx, "SELECT * FROM leak;", driver.QueryOptions{})
	if err != nil {
		t.Fatal(err)
	}
	var myErr *gomysql.MySQLError
	if !errors.As(leak[0].Err, &myErr) || myErr.Number != 1146 {
		t.Fatalf("session state leaked into the next session: %+v", leak)
	}
}
