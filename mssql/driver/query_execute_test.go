package driver

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"

	gomssqldb "github.com/microsoft/go-mssqldb"

	"github.com/bytebase/omni/driver"
	omnimssql "github.com/bytebase/omni/mssql"
)

func quoteIdent(name string) string {
	return "[" + strings.ReplaceAll(name, "]", "]]") + "]"
}

type stepLogger struct {
	steps []string
}

func (l *stepLogger) CommandExecute(command driver.Command) {
	l.steps = append(l.steps, fmt.Sprintf("exec[%d:%d]", command.ByteStart, command.ByteEnd))
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

func TestQuery(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t)
	execCorpus(t, d.db, `
CREATE TABLE q_items (id INT PRIMARY KEY, note NVARCHAR(100), bin VARBINARY(16), flag BIT, amount DECIMAL(10,2), uid UNIQUEIDENTIFIER, created DATETIME2(3), offset_at DATETIMEOFFSET(0));
INSERT INTO q_items VALUES (1, N'a', 0x01, 1, 1.50, 'A0EEBC99-9C0B-4EF8-BB6D-6BB9BD380A11', '2026-03-04 05:06:07.891', '2026-03-04 05:06:07.891 +02:00');
INSERT INTO q_items VALUES (2, N'b', NULL, 0, 2.50, NULL, NULL, NULL);
INSERT INTO q_items VALUES (3, N'c', 0x, 1, 3.50, NULL, NULL, NULL);
GO
CREATE PROCEDURE q_sp_items AS
BEGIN
	SELECT id, note FROM q_items WHERE id = 1;
	SELECT COUNT(*) AS n FROM q_items;
END;
`)

	sess, err := d.Session(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close()

	// One statement per type, then a row read back.
	results, err := sess.Query(ctx, `SELECT
		42 AS i,
		3.5 AS f,
		N'hello' AS s,
		CAST(1 AS BIT) AS b,
		CAST(0 AS BIT) AS b0,
		CAST(12.3400 AS DECIMAL(10,4)) AS num,
		CAST('A0EEBC99-9C0B-4EF8-BB6D-6BB9BD380A11' AS UNIQUEIDENTIFIER) AS uid,
		CAST(0xDEADBEEF AS VARBINARY(4)) AS bin,
		CAST(0x AS VARBINARY(4)) AS empty_bin,
		CAST('12:34:56' AS TIME) AS t,
		CAST('2026-03-04' AS DATE) AS dt,
		CAST('2026-03-04 05:06:07.891' AS DATETIME2(3)) AS ts,
		CAST('2026-03-04 05:06:07 +02:00' AS DATETIMEOFFSET(0)) AS tstz,
		NULL AS nul;`, driver.QueryOptions{})
	if err != nil || len(results) != 1 || results[0].Err != nil {
		var detail error
		if len(results) == 1 {
			detail = results[0].Err
		}
		t.Fatalf("types query = %v, %v", detail, err)
	}
	row := results[0].Rows[0]
	if row[0].Int64 != 42 || row[2].String != "hello" {
		t.Errorf("scalars = %+v", row[:3])
	}
	// BIT is shown as 1/0.
	if row[3].Int64 != 1 || row[4].Int64 != 0 {
		t.Errorf("bits = %v %v", row[3], row[4])
	}
	if row[5].String != "12.3400" {
		t.Errorf("decimal = %q", row[5].String)
	}
	if !strings.EqualFold(row[6].String, "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11") {
		t.Errorf("uniqueidentifier = %q", row[6].String)
	}
	if !bytes.Equal(row[7].Bytes, []byte{0xde, 0xad, 0xbe, 0xef}) {
		t.Errorf("varbinary = %x", row[7].Bytes)
	}
	if row[8].Kind != driver.Bytes || len(row[8].Bytes) != 0 {
		t.Errorf("empty varbinary = %+v", row[8])
	}
	if row[9].String != "12:34:56" || row[10].String != "2026-03-04" {
		t.Errorf("time/date = %q %q", row[9].String, row[10].String)
	}
	if row[11].Kind != driver.Timestamp {
		t.Errorf("datetime2 = %+v", row[11])
	}
	if row[12].Kind != driver.TimestampTZ {
		t.Errorf("datetimeoffset = %+v", row[12])
	}
	if row[13].Kind != driver.Null {
		t.Errorf("null = %+v", row[13])
	}

	// The limit rewrite caps a SELECT and leaves a batch's statements in order.
	limited, err := sess.Query(ctx, "SELECT id FROM q_items ORDER BY id;", driver.QueryOptions{Limit: 2})
	if err != nil || len(limited) != 1 || len(limited[0].Rows) != 2 {
		t.Fatalf("limit query = %+v, %v", limited, err)
	}

	// A row-count statement reports affected rows rather than a result set.
	inserted, err := sess.Query(ctx, "INSERT INTO q_items (id, note, flag) VALUES (4, N'd', 1);", driver.QueryOptions{})
	if err != nil || len(inserted) != 1 || inserted[0].AffectedRows == nil || *inserted[0].AffectedRows != 1 {
		t.Fatalf("insert result = %+v, %v", inserted, err)
	}

	// A procedure's result sets are returned, in order, after the statements'.
	procedure, err := sess.Query(ctx, "EXEC q_sp_items;", driver.QueryOptions{})
	if err != nil || len(procedure) != 2 {
		t.Fatalf("procedure results = %+v, %v", procedure, err)
	}

	// Every statement of a batch keeps its text, including one that reports
	// neither rows nor a count: a caller shows each result beside its statement.
	texts, err := sess.Query(ctx, "SET NOCOUNT OFF;\nSELECT 1 AS one;", driver.QueryOptions{})
	if err != nil || len(texts) != 2 || texts[0].Err != nil || texts[1].Err != nil {
		t.Fatalf("statement texts = %+v, %v", texts, err)
	}
	if texts[0].Statement != "SET NOCOUNT OFF;" || texts[1].Statement != "\nSELECT 1 AS one;" {
		t.Errorf("statement texts = %q, %q", texts[0].Statement, texts[1].Statement)
	}

	// A procedure's further result sets come after every statement's result of
	// the script, not after those of the procedure's own batch, so a caller
	// pairing results with the request's statements keeps its order.
	acrossBatches, err := sess.Query(ctx, "EXEC q_sp_items;\nGO\nSELECT 1 AS one;", driver.QueryOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(acrossBatches) != 3 {
		t.Fatalf("across batches = %d results, want 3", len(acrossBatches))
	}
	if acrossBatches[1].Statement != "SELECT 1 AS one;" || acrossBatches[2].Statement != "EXEC q_sp_items;" {
		t.Errorf("across batches = %q, %q, %q", acrossBatches[0].Statement, acrossBatches[1].Statement, acrossBatches[2].Statement)
	}
	if len(acrossBatches[2].Rows) != 1 || len(acrossBatches[2].Rows[0]) != 1 {
		t.Errorf("the procedure's further result set = %+v, want its COUNT", acrossBatches[2].Rows)
	}

	// A limit caps a statement the rewrite cannot reach, such as a procedure
	// call's result sets.
	limitedExec, err := sess.Query(ctx, "EXEC('SELECT TOP 200 name FROM sys.all_objects');", driver.QueryOptions{Limit: 2})
	if err != nil || len(limitedExec) != 1 || limitedExec[0].Err != nil {
		t.Fatalf("limited procedure = %+v, %v", limitedExec, err)
	}
	if len(limitedExec[0].Rows) != 2 {
		t.Errorf("limited procedure kept %d rows, want 2", len(limitedExec[0].Rows))
	}
	// A capped result set is read to its end, so the next one still arrives.
	limitedSets, err := sess.Query(ctx, "EXEC('SELECT TOP 200 name FROM sys.all_objects; SELECT 7 AS n;');", driver.QueryOptions{Limit: 2})
	if err != nil || len(limitedSets) != 2 {
		t.Fatalf("limited procedure sets = %+v, %v", limitedSets, err)
	}
	if len(limitedSets[0].Rows) != 2 || len(limitedSets[1].Rows) != 1 {
		t.Errorf("limited procedure sets kept %d and %d rows, want 2 and 1", len(limitedSets[0].Rows), len(limitedSets[1].Rows))
	}

	// With masking on, rows the statement does not trace to columns are withheld.
	withheld, err := sess.Query(ctx, "EXEC q_sp_items;", driver.QueryOptions{MaskingEnabled: true})
	if err != nil || len(withheld) != 2 {
		t.Fatalf("masked procedure results = %+v, %v", withheld, err)
	}
	for _, result := range withheld {
		if len(result.Rows) != 0 || result.Err == nil {
			t.Errorf("masked result kept rows: %+v", result)
		}
	}

	// Withholding a procedure's rows is not a statement failure: the batches
	// after it still run.
	maskedAcrossBatches, err := sess.Query(ctx, "EXEC q_sp_items;\nGO\nSELECT 1 AS one;", driver.QueryOptions{MaskingEnabled: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(maskedAcrossBatches) != 3 {
		t.Fatalf("masked across batches = %d results, want 3", len(maskedAcrossBatches))
	}
	if maskedAcrossBatches[1].Err != nil || maskedAcrossBatches[1].Rows == nil {
		t.Errorf("the statement after the masked procedure = %+v, want its rows", maskedAcrossBatches[1])
	}
	if maskedAcrossBatches[2].Err == nil || maskedAcrossBatches[2].Rows != nil {
		t.Errorf("the procedure's further result set = %+v, want it withheld", maskedAcrossBatches[2])
	}

	// The plan the server prints is the result of an explain request.
	plan, err := sess.Query(ctx, "SELECT note FROM q_items WHERE id > 1;", driver.QueryOptions{Explain: true})
	if err != nil || len(plan) != 1 || plan[0].Err != nil || len(plan[0].Columns) == 0 {
		t.Fatalf("explain = %+v, %v", plan, err)
	}
	physicalOpIndex := -1
	for i, column := range plan[0].Columns {
		if column.Name == "PhysicalOp" {
			physicalOpIndex = i
		}
	}
	if physicalOpIndex < 0 {
		t.Fatalf("explain plan has no PhysicalOp column: %+v", plan[0].Columns)
	}
	if plan[0].Plan == nil || plan[0].Plan.Format != "TEXT" || plan[0].Plan.Executed || plan[0].Plan.Content != "" {
		t.Errorf("tabular plan = %+v, want a TEXT plan in the rows that did not run the statement", plan[0].Plan)
	}
	var physicalOps []string
	for _, row := range plan[0].Rows {
		if op := row[physicalOpIndex].String; op != "" {
			physicalOps = append(physicalOps, op)
		}
	}
	if len(physicalOps) == 0 {
		t.Errorf("explain plan names no operator: %+v", plan[0].Rows)
	}
	xmlPlan, err := sess.Query(ctx, "SELECT note FROM q_items WHERE id > 1;", driver.QueryOptions{Explain: true, ExplainFormat: "XML"})
	if err != nil || len(xmlPlan) != 1 || xmlPlan[0].Err != nil || len(xmlPlan[0].Rows) != 1 {
		t.Fatalf("xml explain = %+v, %v", xmlPlan, err)
	}
	if !strings.Contains(xmlPlan[0].Rows[0][0].String, "ShowPlanXML") {
		t.Errorf("xml plan = %q", xmlPlan[0].Rows[0][0].String)
	}
	if xmlPlan[0].Plan == nil || xmlPlan[0].Plan.Format != "XML" {
		t.Errorf("xml plan metadata = %+v, want XML", xmlPlan[0].Plan)
	}

	// A batch keeps what SQL Server scopes to it: a DECLARE'd variable holds
	// across the batch's statements, which one request carries.
	scoped, err := sess.Query(ctx, "DECLARE @id int = 1;\nSELECT @id AS id;", driver.QueryOptions{})
	if err != nil || len(scoped) != 2 || scoped[1].Err != nil || scoped[1].Rows[0][0].Int64 != 1 {
		t.Fatalf("batch-scoped variable = %+v, %v", scoped, err)
	}

	// A lone failing statement reports the error as its result, with every
	// message the server sent.
	errorResults, err := sess.Query(ctx, "SELECT * FROM q_missing_table;", driver.QueryOptions{})
	if err != nil || len(errorResults) != 1 || errorResults[0].Err == nil {
		t.Fatalf("error result = %+v, %v", errorResults, err)
	}
	if text := errorResults[0].Err.Error(); !strings.Contains(text, "Msg 208") || !strings.Contains(text, "q_missing_table") {
		t.Errorf("error text = %q, want SQL Server's 208 and the object's name", text)
	}

	// A failing statement ends the script: the batch after it does not run.
	failScript, err := sess.Query(ctx, "SELECT * FROM q_missing_table;\nGO\nDELETE FROM q_items WHERE id = 1;", driver.QueryOptions{})
	if err != nil || len(failScript) != 1 || failScript[0].Err == nil {
		t.Fatalf("failing script = %+v, %v", failScript, err)
	}
	var serverError gomssqldb.Error
	if !errors.As(failScript[0].Err, &serverError) || serverError.Number != 208 {
		t.Fatalf("statement error = %v, want SQL Server's 208", failScript[0].Err)
	}
	counted, err := sess.Query(ctx, "SELECT COUNT(*) FROM q_items WHERE id = 1;", driver.QueryOptions{})
	if err != nil || len(counted) != 1 || counted[0].Err != nil || counted[0].Rows[0][0].Int64 != 1 {
		t.Fatalf("the batch after the failure ran: %+v, %v", counted, err)
	}

	// A batch of statements that report no rows answers with columnless results,
	// which still leaves each statement its own result.
	batchy, err := sess.Query(ctx, "CREATE TABLE q_scratch (id INT);\nINSERT INTO q_scratch VALUES (1);", driver.QueryOptions{})
	if err != nil || len(batchy) != 2 || batchy[0].Err != nil || batchy[1].Err != nil {
		t.Fatalf("columnless batch = %+v, %v", batchy, err)
	}

	// A statement that reports no count is not a write of zero rows: DDL, SET and
	// SELECT @v = ... report nothing, so the result carries no count at all.
	created, err := sess.Query(ctx, "CREATE TABLE q_archive (id INT, note NVARCHAR(100));", driver.QueryOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(created) != 1 || created[0].AffectedRows != nil {
		t.Fatalf("DDL result = %+v, want no affected count", created)
	}

	// A limit bounds a result set, so it leaves a SELECT ... INTO materialising
	// every row as it would.
	into, err := sess.Query(ctx, "SELECT id, note INTO q_archive2 FROM q_items;", driver.QueryOptions{Limit: 1})
	if err != nil || len(into) != 1 || into[0].Err != nil {
		t.Fatalf("SELECT INTO = %+v, %v", into, err)
	}
	counted2, err := sess.Query(ctx, "SELECT COUNT(*) FROM q_archive2;", driver.QueryOptions{})
	if err != nil || counted2[0].Rows[0][0].Int64 < 2 {
		t.Fatalf("SELECT INTO materialised %+v, %v; want every row", counted2, err)
	}

	// The size cap truncates the statement and the next one still runs.
	truncated, err := sess.Query(ctx, "SELECT id, REPLICATE(N'x', 200) AS payload FROM q_items ORDER BY id;\nSELECT 7 AS after_cap;", driver.QueryOptions{MaxResultBytes: 200})
	if err != nil || len(truncated) != 2 || !truncated[0].Truncated || truncated[1].Rows[0][0].Int64 != 7 {
		t.Fatalf("truncated = %+v, %v", truncated, err)
	}

	// A result set cut off at the size cap does not swallow the next one: the
	// client learns of it only as it reads past the cut set's end.
	cut, err := sess.Query(ctx, "SELECT id FROM q_items ORDER BY id;\nSELECT 9 AS after_cut;", driver.QueryOptions{MaxResultBytes: 1})
	if err != nil || len(cut) != 2 || !cut[0].Truncated || cut[1].Err != nil {
		t.Fatalf("cut result sets = %+v, %v", cut, err)
	}
	if len(cut[1].Rows) != 1 || cut[1].Rows[0][0].Int64 != 9 {
		t.Errorf("the result set after a cut one = %+v", cut[1].Rows)
	}

	// A batch separator splits batches, and a failed statement ends its own.
	batches, err := sess.Query(ctx, "SELECT 1 AS a;\nGO\nSELECT 2 AS b;", driver.QueryOptions{})
	if err != nil || len(batches) != 2 || batches[1].Rows[0][0].Int64 != 2 {
		t.Fatalf("GO batches = %+v, %v", batches, err)
	}

	// A batch of more than one statement yields an empty result for a statement
	// that reports only a count: SQL Server states those in the ReturnMessage
	// message loop, which a tape cannot carry, having no message argument and no
	// ordered counts beside its result sets.
	inserted, err = sess.Query(ctx, "INSERT INTO q_items (id, note, flag) VALUES (5, N'e', 1);\nSELECT 1 AS after_insert;", driver.QueryOptions{})
	if err != nil || len(inserted) != 2 || inserted[0].AffectedRows != nil || inserted[1].Rows[0][0].Int64 != 1 {
		t.Fatalf("batched insert = %+v, %v", inserted, err)
	}

	// A statement the parser cannot read fails the batch before it runs.
	if _, err := sess.Query(ctx, "SELECT FROM WHERE;", driver.QueryOptions{}); err == nil {
		t.Fatal("an unparseable batch should fail")
	}
}

func TestExecute(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t)
	execCorpus(t, d.db, "CREATE TABLE e_items (id INT PRIMARY KEY, note NVARCHAR(64));")

	// Transactional: both batches commit, and every step is reported.
	log := &stepLogger{}
	affected, err := d.Execute(ctx, "INSERT INTO e_items VALUES (1, N'a'), (2, N'b');\nGO\nUPDATE e_items SET note = N'z' WHERE id <= 1;", driver.ExecuteOptions{
		Transaction: driver.Transactional,
		Log:         log,
	})
	if err != nil || affected != 3 {
		t.Fatalf("execute = %d, %v", affected, err)
	}
	wantSteps := []string{
		"tx:BEGIN",
		"exec[0:48]", "resp:2:[2]",
		"exec[52:97]", "resp:1:[1]",
		"tx:COMMIT",
	}
	if !slices.Equal(log.steps, wantSteps) {
		t.Fatalf("steps = %v, want %v", log.steps, wantSteps)
	}

	// A failure rolls the batch back.
	rollback := &stepLogger{}
	if _, err := d.Execute(ctx, "INSERT INTO e_items VALUES (10, N'x');\nGO\nINSERT INTO e_items VALUES (1, N'duplicate');", driver.ExecuteOptions{
		Transaction: driver.Transactional,
		Log:         rollback,
	}); err == nil {
		t.Fatal("expected a duplicate key error")
	}
	if !slices.Contains(rollback.steps, "tx:ROLLBACK") {
		t.Fatalf("steps = %v, want a rollback", rollback.steps)
	}
	var count int
	if err := d.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM e_items WHERE id = 10").Scan(&count); err != nil || count != 0 {
		t.Fatalf("rolled-back row count = %d, %v", count, err)
	}

	// Auto-commit keeps what ran before the failure.
	partial, err := d.Execute(ctx, "INSERT INTO e_items VALUES (20, N'kept');\nGO\nINSERT INTO e_items VALUES (1, N'duplicate');", driver.ExecuteOptions{
		Transaction: driver.AutoCommit,
	})
	if err == nil || partial != 1 {
		t.Fatalf("auto-commit = %d, %v", partial, err)
	}
	if err := d.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM e_items WHERE id = 20").Scan(&count); err != nil || count != 1 {
		t.Fatalf("auto-commit row count = %d, %v", count, err)
	}

	// GO <count> runs the batch that many times.
	repeated, err := d.Execute(ctx, "UPDATE e_items SET note = N'r' WHERE id = 2;\nGO 3", driver.ExecuteOptions{Transaction: driver.AutoCommit})
	if err != nil || repeated != 3 {
		t.Fatalf("GO 3 = %d, %v", repeated, err)
	}

	// CreateDatabase runs the statement as given.
	newDatabase := "omni_created_" + strings.ToLower(strings.ReplaceAll(t.Name(), "/", "_"))
	defer func() {
		_, _ = d.db.ExecContext(context.Background(), "DROP DATABASE IF EXISTS "+quoteIdent(newDatabase))
	}()
	if _, err := d.Execute(ctx, "CREATE DATABASE "+quoteIdent(newDatabase), driver.ExecuteOptions{CreateDatabase: true}); err != nil {
		t.Fatal(err)
	}
	if _, err := d.SyncInstance(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestCountAffectedRows(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t)
	execCorpus(t, d.db, `
CREATE TABLE c_big (id INT PRIMARY KEY, s_id INT NOT NULL, v INT);
INSERT INTO c_big SELECT TOP 1000 ROW_NUMBER() OVER (ORDER BY (SELECT NULL)), (ROW_NUMBER() OVER (ORDER BY (SELECT NULL)) % 100) + 1, 0 FROM sys.all_objects;
CREATE TABLE c_archive (id INT PRIMARY KEY, s_id INT NOT NULL, v INT);
UPDATE STATISTICS c_big;
`)

	estimated, err := d.CountAffectedRows(ctx, "DELETE FROM c_big WHERE s_id <= 50;")
	if err != nil {
		t.Fatal(err)
	}
	if estimated <= 0 {
		t.Fatalf("estimate = %d, want more than 0", estimated)
	}
	if _, err := d.CountAffectedRows(ctx, "SELECT * FROM app.missing_table;"); err == nil {
		t.Fatal("a missing table should fail its plan")
	}
}

// TestClassifyReceiveAndReadText pins that the statements whose rows masking
// cannot trace are marked, so masking withholds them.
func TestClassifyReceiveAndReadText(t *testing.T) {
	for _, tc := range []struct {
		sql  string
		want stmtType
	}{
		{"RECEIVE conversation_handle FROM dbo.q;", stmtTypeResultSetGenerating | stmtTypeRowCountGenerating | stmtTypeUntraced},
		{"RECEIVE @h = conversation_handle FROM dbo.q;", stmtTypeRowCountGenerating},
		{"READTEXT t.c @ptr 0 100;", stmtTypeResultSetGenerating | stmtTypeRowCountGenerating | stmtTypeUntraced},
	} {
		stmts, err := omnimssql.Parse(tc.sql)
		if err != nil || len(stmts) != 1 {
			t.Fatalf("%q: %v", tc.sql, err)
		}
		got, err := getStmtType(stmts[0].AST)
		if err != nil {
			t.Fatalf("%q: %v", tc.sql, err)
		}
		if got != tc.want {
			t.Errorf("%q: type = %d, want %d", tc.sql, got, tc.want)
		}
	}
}
