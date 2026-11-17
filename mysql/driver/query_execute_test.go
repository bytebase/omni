package driver

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/driver"
	mysqlanalysis "github.com/bytebase/omni/mysql/analysis"
	"github.com/bytebase/omni/mysql/internal/testcontainer"
	mysqlparser "github.com/bytebase/omni/mysql/parser"
)

// queryRoles are the pinned servers the query and execute corpora run over:
// mysql:8.0.46 with lower_case_table_names=1, and the 8.4 LTS release.
var queryRoles = []string{testcontainer.Folded, testcontainer.MySQL84}

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

// openQueryTestDriver is openTestDriver with the application name a caller
// sets, which comments the statements this driver sends. The connection-id
// probe of the kill-on-cancel path stays off: the server answers the id of a
// connection, which a tape could only hold by pinning what ran before the test.
func openQueryTestDriver(t *testing.T, role, dbName string) *Driver {
	t.Helper()
	d := openTestDriver(t, role, dbName)
	d.appName = "bytebase"
	return d
}

// execSpans is where the step logger sees statements. Each command begins
// where the previous one ended and ends past its semicolon, which is how the
// splitter hands a script's statements over: Bytebase's Statement.Text keeps
// the separators before a statement as this does.
func execSpans(script string, statements ...string) []string {
	spans := make([]string, 0, len(statements))
	start := 0
	for _, statement := range statements {
		index := strings.Index(script[start:], statement)
		if index < 0 {
			spans = append(spans, "exec[not in script: "+statement+"]")
			continue
		}
		index += start
		end := index + len(statement) + 1
		spans = append(spans, fmt.Sprintf("exec[%d:%d]", start, end))
		start = end
	}
	return spans
}

const querySetup = `CREATE TABLE qitems (
  id BIGINT PRIMARY KEY,
  note VARCHAR(64),
  amount DECIMAL(10,2),
  created DATETIME(6),
  flag TINYINT(1),
  payload JSON,
  raw VARBINARY(8),
  bigu BIGINT UNSIGNED
);
INSERT INTO qitems VALUES
  (1, 'a', 1.25, '2026-03-04 05:06:07.123456', 1, '{"a": 1}', 0x00ff, 18446744073709551615),
  (2, 'b', 2.50, '2026-03-04 05:06:08', 0, '{"b": 2}', 0x10, 10),
  (3, 'c', NULL, NULL, NULL, NULL, NULL, NULL);
CREATE PROCEDURE qcount()
BEGIN
  SELECT COUNT(*) AS n FROM qitems;
END;`

// TestQueryCorpus runs the query path over the pinned servers: what it sends
// for each statement, the rows and column types it decodes, the limit rewrite,
// the explain formats, and per-statement errors.
func TestQueryCorpus(t *testing.T) {
	ctx := context.Background()
	for _, role := range queryRoles {
		t.Run(role, func(t *testing.T) {
			d := openQueryTestDriver(t, role, "omni_qe_query")
			sess, err := d.Session(ctx)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = sess.Close() }()

			// 1. A script of statements that return no rows reports one affected
			// count each, and the statement text it sent.
			setup, err := sess.Query(ctx, querySetup, driver.QueryOptions{})
			if err != nil || len(setup) != 3 {
				t.Fatalf("setup = %+v, %v", setup, err)
			}
			if setup[0].AffectedRows == nil || *setup[0].AffectedRows != 0 || len(setup[0].Rows) != 0 {
				t.Errorf("CREATE TABLE result = %+v", setup[0])
			}
			if !strings.HasPrefix(strings.TrimSpace(setup[1].Statement), "INSERT INTO qitems") {
				t.Errorf("INSERT statement = %q", setup[1].Statement)
			}
			if setup[1].AffectedRows == nil || *setup[1].AffectedRows != 3 {
				t.Errorf("INSERT affected rows = %v, want 3", setup[1].AffectedRows)
			}
			if setup[2].AffectedRows == nil || *setup[2].AffectedRows != 0 {
				t.Errorf("CREATE PROCEDURE affected rows = %v, want 0", setup[2].AffectedRows)
			}

			// 2. Column names, engine type names, and every value kind a MySQL
			// column decodes to. An unsigned BIGINT has no scan case in the
			// shared value maker, so it arrives as a string: Bytebase reads it
			// the same way.
			all, err := sess.Query(ctx, "SELECT * FROM qitems ORDER BY id;", driver.QueryOptions{})
			if err != nil || len(all) != 1 || all[0].Err != nil {
				t.Fatalf("SELECT * = %+v, %v", all, err)
			}
			got := all[0]
			var names, types []string
			for _, c := range got.Columns {
				names = append(names, c.Name)
				types = append(types, c.TypeName)
			}
			wantNames := []string{"id", "note", "amount", "created", "flag", "payload", "raw", "bigu"}
			wantTypes := []string{"BIGINT", "VARCHAR", "DECIMAL", "DATETIME", "TINYINT", "JSON", "VARBINARY", "UNSIGNED BIGINT"}
			if !slices.Equal(names, wantNames) || !slices.Equal(types, wantTypes) {
				t.Fatalf("columns = %v %v, want %v %v", names, types, wantNames, wantTypes)
			}
			if len(got.Rows) != 3 {
				t.Fatalf("rows = %d, want 3", len(got.Rows))
			}
			first := got.Rows[0]
			if first[0].Kind != driver.Int64 || first[0].Int64 != 1 ||
				first[1].Kind != driver.String || first[1].String != "a" ||
				first[2].Kind != driver.String || first[2].String != "1.25" ||
				first[3].Kind != driver.Timestamp || !first[3].Time.Equal(time.Date(2026, 3, 4, 5, 6, 7, 123456000, time.UTC)) || first[3].Precision != 6 ||
				first[4].Kind != driver.Int64 || first[4].Int64 != 1 ||
				first[5].Kind != driver.String || first[5].String != `{"a": 1}` ||
				first[6].Kind != driver.Bytes || !bytes.Equal(first[6].Bytes, []byte{0x00, 0xff}) ||
				first[7].Kind != driver.String || first[7].String != "18446744073709551615" {
				t.Errorf("row 0 = %+v", first)
			}
			last := got.Rows[2]
			if last[0].Kind != driver.Int64 || last[0].Int64 != 3 ||
				last[1].Kind != driver.String || last[1].String != "c" {
				t.Errorf("row 2 = %+v", last)
			}
			for i, v := range last[2:] {
				if v.Kind != driver.Null {
					t.Errorf("row 2 column %d = %+v, want NULL", i+2, v)
				}
			}

			// 3. The limit rewrite: an absent LIMIT is added, a larger one is
			// lowered, and a statement that returns no rows is sent unchanged.
			limited, err := sess.Query(ctx, "SELECT id FROM qitems ORDER BY id;", driver.QueryOptions{Limit: 2})
			if err != nil || len(limited) != 1 || limited[0].Err != nil {
				t.Fatalf("limited = %+v, %v", limited, err)
			}
			if want := "SELECT id FROM qitems ORDER BY id LIMIT 2;"; limited[0].Statement != want {
				t.Errorf("limited statement = %q, want %q", limited[0].Statement, want)
			}
			if len(limited[0].Rows) != 2 {
				t.Errorf("limited rows = %d, want 2", len(limited[0].Rows))
			}
			lowered, err := sess.Query(ctx, "SELECT id FROM qitems ORDER BY id LIMIT 10;", driver.QueryOptions{Limit: 2})
			if err != nil || len(lowered) != 1 || len(lowered[0].Rows) != 2 {
				t.Fatalf("lowered = %+v, %v", lowered, err)
			}
			if want := "SELECT id FROM qitems ORDER BY id LIMIT 2;"; lowered[0].Statement != want {
				t.Errorf("lowered statement = %q, want %q", lowered[0].Statement, want)
			}
			unwrapped, err := sess.Query(ctx, "SHOW TABLES;", driver.QueryOptions{Limit: 1})
			if err != nil || len(unwrapped) != 1 || unwrapped[0].Statement != "SHOW TABLES;" {
				t.Fatalf("SHOW TABLES with a limit = %+v, %v", unwrapped, err)
			}

			// 4. Explaining: the JSON plan is one row of JSON with a query
			// block, the default plan states no format because MySQL's follows
			// the session's explain_format, and a statement EXPLAIN cannot plan
			// or one that ANALYZE would run fail the call before anything runs.
			jsonPlan, err := sess.Query(ctx, "SELECT note FROM qitems WHERE id > 1;", driver.QueryOptions{
				Explain:       true,
				ExplainFormat: "JSON",
			})
			if err != nil || len(jsonPlan) != 1 || jsonPlan[0].Err != nil {
				t.Fatalf("JSON explain = %+v, %v", jsonPlan, err)
			}
			if jsonPlan[0].Plan == nil || jsonPlan[0].Plan.Format != "JSON" || jsonPlan[0].Plan.Executed {
				t.Errorf("JSON plan = %+v", jsonPlan[0].Plan)
			}
			if jsonPlan[0].Statement != "EXPLAIN FORMAT=JSON SELECT note FROM qitems WHERE id > 1;" {
				t.Errorf("JSON plan statement = %q", jsonPlan[0].Statement)
			}
			var body strings.Builder
			for _, row := range jsonPlan[0].Rows {
				body.WriteString(row[0].String)
			}
			var plan map[string]any
			if err := json.Unmarshal([]byte(body.String()), &plan); err != nil || plan["query_block"] == nil {
				t.Errorf("JSON plan body = %q, err = %v", body.String(), err)
			}

			textPlan, err := sess.Query(ctx, "SELECT note FROM qitems WHERE id > 1;", driver.QueryOptions{Explain: true})
			if err != nil || len(textPlan) != 1 || textPlan[0].Err != nil {
				t.Fatalf("text explain = %+v, %v", textPlan, err)
			}
			if textPlan[0].Plan == nil || textPlan[0].Plan.Format != "" {
				t.Errorf("text plan = %+v, want a plan whose format is unstated", textPlan[0].Plan)
			}
			if len(textPlan[0].Rows) == 0 {
				t.Error("text plan returned no rows")
			}

			if _, err := sess.Query(ctx, "SHOW TABLES;", driver.QueryOptions{Explain: true}); !errors.Is(err, mysqlanalysis.ErrNotExplainable) {
				t.Errorf("SHOW TABLES explain err = %v, want ErrNotExplainable", err)
			}
			if _, err := sess.Query(ctx, "EXPLAIN ANALYZE SELECT 1;", driver.QueryOptions{Explain: true}); !errors.Is(err, mysqlanalysis.ErrExplainAnalyzeExecutes) {
				t.Errorf("EXPLAIN ANALYZE err = %v, want ErrExplainAnalyzeExecutes", err)
			}

			// 5. A statement that fails at the server carries its error, ends
			// the script, and leaves the earlier results in place.
			failed, err := sess.Query(ctx, "SELECT 1;\nSELECT * FROM missing_tbl;\nSELECT 2;", driver.QueryOptions{})
			if err != nil || len(failed) != 2 {
				t.Fatalf("runtime error script = %+v, %v", failed, err)
			}
			var myErr *gomysql.MySQLError
			if !errors.As(failed[1].Err, &myErr) || myErr.Number != 1146 {
				t.Fatalf("failed[1].Err = %v, want MySQLError 1146", failed[1].Err)
			}
			if statement := strings.TrimSpace(failed[1].Statement); statement != "SELECT * FROM missing_tbl;" {
				t.Errorf("failed statement = %q", statement)
			}

			// 6. A script the parser cannot read runs nothing: every statement
			// is read before the first is sent.
			const invalid = "INSERT INTO qitems (id) VALUES (99);\nSELECT FROM WHERE;"
			_, err = sess.Query(ctx, invalid, driver.QueryOptions{})
			if err == nil {
				t.Fatal("expected a syntax error")
			}
			var parseErr *mysqlparser.ParseError
			if !errors.As(err, &parseErr) {
				t.Fatalf("syntax error = %v, want *parser.ParseError", err)
			}
			if want := strings.Index(invalid, "FROM"); parseErr.Position != want {
				t.Errorf("parse error position = %d, want the script's %d (%q)", parseErr.Position, want, invalid[parseErr.Position:])
			}
			if parseErr.Line != 2 || parseErr.RelatedText != "SELECT FROM WHERE;" {
				t.Errorf("parse error = line %d, related %q", parseErr.Line, parseErr.RelatedText)
			}
			after, err := sess.Query(ctx, "SELECT COUNT(*) FROM qitems WHERE id = 99;", driver.QueryOptions{})
			if err != nil || after[0].Rows[0][0].Int64 != 0 {
				t.Fatalf("statement before the syntax error ran: %+v, %v", after, err)
			}

			// 7. The result size cap truncates a statement's rows and the next
			// statement still runs.
			truncated, err := sess.Query(ctx, "SELECT id, repeat('x', 80) FROM qitems ORDER BY id;\nSELECT 7;", driver.QueryOptions{MaxResultBytes: 120})
			if err != nil || len(truncated) != 2 {
				t.Fatalf("truncated = %+v, %v", truncated, err)
			}
			if !truncated[0].Truncated || len(truncated[0].Rows) >= 3 {
				t.Errorf("truncated[0] = truncated %v with %d rows", truncated[0].Truncated, len(truncated[0].Rows))
			}
			if truncated[1].Truncated || truncated[1].Rows[0][0].Int64 != 7 {
				t.Errorf("truncated[1] = %+v", truncated[1])
			}

			// 8. A statement the read-only whitelist does not name still answers
			// with rows — a procedure's result set, a table check's diagnosis —
			// and masking withholds them, because the caller's analysis cannot
			// trace the columns.
			call, err := sess.Query(ctx, "CALL qcount();", driver.QueryOptions{})
			if err != nil || len(call) != 1 || call[0].Err != nil {
				t.Fatalf("CALL = %+v, %v", call, err)
			}
			if len(call[0].Rows) != 1 || call[0].Rows[0][0].Int64 != 3 {
				t.Errorf("CALL rows = %+v", call[0].Rows)
			}
			maskedCall, err := sess.Query(ctx, "CALL qcount();", driver.QueryOptions{MaskingEnabled: true})
			if err != nil || len(maskedCall) != 1 || len(maskedCall[0].Rows) != 0 || maskedCall[0].AffectedRows == nil {
				t.Fatalf("masked CALL = %+v, %v", maskedCall, err)
			}
			checked, err := sess.Query(ctx, "CHECK TABLE qitems;", driver.QueryOptions{})
			if err != nil || len(checked) != 1 || len(checked[0].Rows) == 0 {
				t.Fatalf("CHECK TABLE = %+v, %v", checked, err)
			}

			// An executable comment that holds no statement is nothing to run,
			// not a panic.
			if none, err := sess.Query(ctx, "/*! */;", driver.QueryOptions{}); err != nil || len(none) != 0 {
				t.Fatalf("empty executable comment = %+v, %v", none, err)
			}

			// 9. A write that returns no rows reports its count, and masking
			// asks this driver for nothing: the same rows come back either way.
			updated, err := sess.Query(ctx, "UPDATE qitems SET note = 'z' WHERE id <= 2;", driver.QueryOptions{})
			if err != nil || len(updated) != 1 || updated[0].AffectedRows == nil || *updated[0].AffectedRows != 2 || len(updated[0].Rows) != 0 {
				t.Fatalf("UPDATE result = %+v, %v", updated, err)
			}
			plain, err := sess.Query(ctx, "SELECT id, note FROM qitems ORDER BY id;", driver.QueryOptions{})
			if err != nil {
				t.Fatal(err)
			}
			masked, err := sess.Query(ctx, "SELECT id, note FROM qitems ORDER BY id;", driver.QueryOptions{MaskingEnabled: true})
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(plain[0].Rows, masked[0].Rows) {
				t.Errorf("masking changed the rows: %+v vs %+v", plain[0].Rows, masked[0].Rows)
			}
		})
	}
}

const executeSetup = `CREATE TABLE eitems (id INT PRIMARY KEY, note VARCHAR(64));
INSERT INTO eitems VALUES (1, 'a'), (2, 'b'), (3, 'c');
UPDATE eitems SET note = 'z' WHERE id <= 2;`

// TestExecuteCorpus runs change scripts over the pinned servers: a transaction
// that commits, one that rolls back a failure mid-script, auto-commit leaving
// what ran before the failure, a DDL/DML mix, and a stored program.
func TestExecuteCorpus(t *testing.T) {
	ctx := context.Background()
	for _, role := range queryRoles {
		t.Run(role, func(t *testing.T) {
			d := openQueryTestDriver(t, role, "omni_qe_execute")

			// 1. One transaction around the script, at a stated isolation level.
			log := &stepLogger{}
			affected, err := d.Execute(ctx, executeSetup, driver.ExecuteOptions{
				Transaction: driver.Transactional,
				Isolation:   driver.ReadCommitted,
				Log:         log,
			})
			if err != nil || affected != 5 {
				t.Fatalf("Execute = (%d, %v)", affected, err)
			}
			spans := execSpans(executeSetup,
				"CREATE TABLE eitems (id INT PRIMARY KEY, note VARCHAR(64))",
				"INSERT INTO eitems VALUES (1, 'a'), (2, 'b'), (3, 'c')",
				"UPDATE eitems SET note = 'z' WHERE id <= 2",
			)
			want := []string{
				"tx:BEGIN",
				spans[0], "resp:0:[0]",
				spans[1], "resp:3:[3]",
				spans[2], "resp:2:[2]",
				"tx:COMMIT",
			}
			if !slices.Equal(log.steps, want) {
				t.Fatalf("steps = %v, want %v", log.steps, want)
			}

			// 2. A failure mid-script rolls the transaction back, and the
			// statement that ran before it leaves nothing behind.
			rollbackLog := &stepLogger{}
			if _, err := d.Execute(ctx, "INSERT INTO eitems VALUES (10, 'rb');\nINSERT INTO eitems VALUES (1, 'dup');", driver.ExecuteOptions{
				Transaction: driver.Transactional,
				Isolation:   driver.RepeatableRead,
				Log:         rollbackLog,
			}); err == nil {
				t.Fatal("expected a duplicate key error")
			}
			if got := rollbackLog.steps[len(rollbackLog.steps)-1]; got != "tx:ROLLBACK" {
				t.Fatalf("steps = %v, want a rollback", rollbackLog.steps)
			}
			if !slices.Contains(rollbackLog.steps, "resp:err") {
				t.Errorf("steps = %v, want the failing response", rollbackLog.steps)
			}

			// 3. Auto-commit keeps what ran before the failure.
			autoLog := &stepLogger{}
			if _, err := d.Execute(ctx, "INSERT INTO eitems VALUES (20, 'kept');\nINSERT INTO eitems VALUES (1, 'dup');", driver.ExecuteOptions{
				Transaction: driver.AutoCommit,
				Log:         autoLog,
			}); err == nil {
				t.Fatal("expected a duplicate key error in auto-commit")
			}
			for _, step := range autoLog.steps {
				if strings.HasPrefix(step, "tx:") {
					t.Fatalf("auto-commit logged a transaction step: %v", autoLog.steps)
				}
			}

			// 4. A DDL/DML mix commits the DDL even when a later statement
			// fails: MySQL commits implicitly on every DDL statement.
			mixedLog := &stepLogger{}
			if _, err := d.Execute(ctx, "CREATE TABLE emixed (id INT PRIMARY KEY);\nINSERT INTO emixed VALUES (1);\nINSERT INTO emixed VALUES (1);", driver.ExecuteOptions{
				Transaction: driver.Transactional,
				Log:         mixedLog,
			}); err == nil {
				t.Fatal("expected a duplicate key error in the DDL/DML mix")
			}

			// 5. A database-creation step runs outside any transaction.
			createLog := &stepLogger{}
			if _, err := d.Execute(ctx, "CREATE TABLE ecreated (id INT PRIMARY KEY);\nINSERT INTO ecreated VALUES (1);", driver.ExecuteOptions{
				Transaction:    driver.Transactional,
				CreateDatabase: true,
				Log:            createLog,
			}); err != nil {
				t.Fatalf("CreateDatabase Execute = %v", err)
			}
			for _, step := range createLog.steps {
				if strings.HasPrefix(step, "tx:") {
					t.Fatalf("CreateDatabase logged a transaction step: %v", createLog.steps)
				}
			}

			// 6. A stored program is one command after the directives, which
			// the splitter keeps out of the text it sends.
			procedureLog := &stepLogger{}
			const procedure = "DELIMITER //\nCREATE PROCEDURE ecounter(OUT n INT)\nBEGIN\n  SELECT COUNT(*) INTO n FROM eitems;\nEND//\nDELIMITER ;\n"
			if _, err := d.Execute(ctx, procedure, driver.ExecuteOptions{
				Transaction: driver.AutoCommit,
				Log:         procedureLog,
			}); err != nil {
				t.Fatalf("DELIMITER Execute = %v", err)
			}
			procedureStart := strings.Index(procedure, "CREATE PROCEDURE")
			procedureEnd := strings.Index(procedure, "END//") + len("END")
			wantProcedure := []string{fmt.Sprintf("exec[%d:%d]", procedureStart, procedureEnd), "resp:0:[0]"}
			if !slices.Equal(procedureLog.steps, wantProcedure) {
				t.Fatalf("procedure steps = %v, want %v", procedureLog.steps, wantProcedure)
			}

			// 7. An isolation level MySQL does not have is refused before any
			// connection is made.
			if _, err := d.Execute(ctx, "SELECT 1;", driver.ExecuteOptions{
				Transaction: driver.Transactional,
				Isolation:   "BOGUS",
			}); err == nil || !strings.Contains(err.Error(), "invalid isolation level for MySQL") {
				t.Fatalf("bogus isolation err = %v", err)
			}

			// 8. What the server holds afterwards: the transaction's rows, the
			// auto-commit statement's row, neither rolled-back row, and the
			// DDL that committed beside its failed DML.
			sess, err := d.Session(ctx)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = sess.Close() }()
			state, err := sess.Query(ctx, `SELECT
				(SELECT COUNT(*) FROM eitems),
				(SELECT COUNT(*) FROM eitems WHERE id = 10),
				(SELECT COUNT(*) FROM eitems WHERE id = 20),
				(SELECT COUNT(*) FROM emixed),
				(SELECT COUNT(*) FROM ecreated),
				(SELECT COUNT(*) FROM information_schema.routines WHERE routine_schema = DATABASE() AND routine_name = 'ecounter');`, driver.QueryOptions{})
			if err != nil || len(state) != 1 || state[0].Err != nil {
				t.Fatalf("state = %+v, %v", state, err)
			}
			var got []int64
			for _, v := range state[0].Rows[0] {
				got = append(got, v.Int64)
			}
			if want := []int64{4, 0, 1, 1, 1, 1}; !slices.Equal(got, want) {
				t.Fatalf("final state = %v, want %v", got, want)
			}
		})
	}
}

// TestCountAffectedRowsAndExplain runs the live probes Bytebase calls on the
// MySQL driver: the planner's estimate, and an EXPLAIN dry run.
func TestCountAffectedRowsAndExplain(t *testing.T) {
	ctx := context.Background()
	for _, role := range queryRoles {
		t.Run(role, func(t *testing.T) {
			d := openQueryTestDriver(t, role, "omni_qe_stats")
			if _, err := d.Execute(ctx, `CREATE TABLE sitems (id INT PRIMARY KEY, note VARCHAR(64));
INSERT INTO sitems VALUES (1, 'a'), (2, 'b'), (3, 'c'), (4, 'd');`, driver.ExecuteOptions{Transaction: driver.Transactional}); err != nil {
				t.Fatal(err)
			}

			estimate, err := d.CountAffectedRows(ctx, "DELETE FROM sitems WHERE id <= 2;")
			if err != nil || estimate < 1 {
				t.Fatalf("CountAffectedRows = (%d, %v)", estimate, err)
			}
			estimate, err = d.CountAffectedRows(ctx, "UPDATE sitems SET note = 'x' WHERE id = 2;")
			if err != nil || estimate != 1 {
				t.Fatalf("CountAffectedRows of a key lookup = (%d, %v)", estimate, err)
			}
			// A statement that changes nothing is refused before any plan is
			// read, so a SELECT of a table cannot report the rows it reads as
			// rows it modifies.
			for _, statement := range []string{"SELECT 1;", "SELECT * FROM sitems;"} {
				if _, err := d.CountAffectedRows(ctx, statement); err == nil || !strings.Contains(err.Error(), "expected an INSERT, UPDATE or DELETE") {
					t.Errorf("CountAffectedRows(%q) err = %v", statement, err)
				}
			}
			if _, err := d.CountAffectedRows(ctx, "DELETE FROM sitems WHERE id <= 2; SELECT 1;"); err == nil || !strings.Contains(err.Error(), "exactly one statement") {
				t.Errorf("CountAffectedRows of a script err = %v", err)
			}

			plan, err := d.ExplainJSON(ctx, "SELECT * FROM sitems")
			if err != nil || !strings.Contains(plan, "query_block") {
				t.Fatalf("ExplainJSON = %q, %v", plan, err)
			}
			if err := d.Explain(ctx, "UPDATE sitems SET note = 'x' WHERE id = 1"); err != nil {
				t.Fatalf("Explain = %v", err)
			}
			// The probes plan the one statement they were given: a second one
			// would run, and a DDL would commit despite the rollback.
			if err := d.Explain(ctx, "SELECT 1; DROP TABLE sitems;"); err == nil || !strings.Contains(err.Error(), "exactly one statement") {
				t.Errorf("Explain of a script err = %v", err)
			}
			if _, err := d.ExplainJSON(ctx, "SELECT 1; DROP TABLE sitems;"); err == nil || !strings.Contains(err.Error(), "exactly one statement") {
				t.Errorf("ExplainJSON of a script err = %v", err)
			}
			if _, _, ok := d.ExplainTabularEstimate(ctx, "SELECT 1; DROP TABLE sitems;"); ok {
				t.Error("ExplainTabularEstimate planned a script")
			}
			kept, err := d.ExplainJSON(ctx, "SELECT COUNT(*) FROM sitems")
			if err != nil || !strings.Contains(kept, "query_block") {
				t.Errorf("ExplainJSON after the refused probes = %q, %v", kept, err)
			}
		})
	}
}
