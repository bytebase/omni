package driver

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/driver"
)

// TestQuery runs scripts through a session and holds what each statement
// returned: its columns, rows, affected rows and error.
func TestQuery(t *testing.T) {
	ctx := context.Background()

	t.Run("types", func(t *testing.T) {
		d := openDriver(t)
		session := openSession(t, d)
		results, err := session.Query(ctx, `SELECT
			CAST(42 AS INTEGER) AS c_integer,
			CAST(42 AS BIGINT) AS c_bigint,
			true AS c_boolean,
			CAST('hi' AS VARCHAR) AS c_varchar,
			CAST('ab' AS CHAR(5)) AS c_char,
			CAST(1.25 AS DECIMAL(5,2)) AS c_decimal,
			CAST(1.5 AS REAL) AS c_real,
			CAST(2.5 AS DOUBLE) AS c_double,
			DATE '2026-01-02' AS c_date,
			TIME '12:34:56.789' AS c_time,
			TIME '12:34:56 +02:00' AS c_time_tz,
			TIMESTAMP '2026-01-02 12:34:56' AS c_timestamp,
			CAST('2026-01-02 12:34:56.123 +02:00' AS TIMESTAMP(3) WITH TIME ZONE) AS c_timestamp_tz,
			X'0102' AS c_varbinary,
			UUID '12151fd2-7586-11e9-8f9e-2a86e4085a59' AS c_uuid,
			CAST('10.0.0.1' AS IPADDRESS) AS c_ipaddress,
			INTERVAL '1' DAY AS c_interval,
			CAST('{"a":1}' AS JSON) AS c_json,
			CAST(NULL AS VARCHAR) AS c_null;`, driver.QueryOptions{})
		if err != nil {
			t.Fatal(err)
		}
		if len(results) != 1 || results[0].Err != nil {
			t.Fatalf("results = %+v", results)
		}
		result := results[0]
		wantColumns := []driver.Column{
			{Name: "c_integer", TypeName: "INTEGER"},
			{Name: "c_bigint", TypeName: "BIGINT"},
			{Name: "c_boolean", TypeName: "BOOLEAN"},
			{Name: "c_varchar", TypeName: "VARCHAR"},
			{Name: "c_char", TypeName: "CHAR"},
			{Name: "c_decimal", TypeName: "DECIMAL"},
			{Name: "c_real", TypeName: "REAL"},
			{Name: "c_double", TypeName: "DOUBLE"},
			{Name: "c_date", TypeName: "DATE"},
			{Name: "c_time", TypeName: "TIME"},
			{Name: "c_time_tz", TypeName: "TIME WITH TIME ZONE"},
			{Name: "c_timestamp", TypeName: "TIMESTAMP"},
			{Name: "c_timestamp_tz", TypeName: "TIMESTAMP WITH TIME ZONE"},
			{Name: "c_varbinary", TypeName: "VARBINARY"},
			{Name: "c_uuid", TypeName: "UUID"},
			{Name: "c_ipaddress", TypeName: "IPADDRESS"},
			{Name: "c_interval", TypeName: "INTERVAL DAY TO SECOND"},
			{Name: "c_json", TypeName: "JSON"},
			{Name: "c_null", TypeName: "VARCHAR"},
		}
		if diff := cmp.Diff(wantColumns, result.Columns); diff != "" {
			t.Errorf("columns (-want +got):\n%s", diff)
		}
		if len(result.Rows) != 1 {
			t.Fatalf("rows = %d, want 1", len(result.Rows))
		}
		wantRow := []driver.Value{
			{Kind: driver.Int64, Int64: 42},
			{Kind: driver.Int64, Int64: 42},
			{Kind: driver.Bool, Bool: true},
			{Kind: driver.String, String: "hi"},
			{Kind: driver.String, String: "ab   "},
			{Kind: driver.Float64, Float64: 1.25},
			{Kind: driver.Float64, Float64: 1.5},
			{Kind: driver.Float64, Float64: 2.5},
			{Kind: driver.String, String: "2026-01-02"},
			{Kind: driver.String, String: "12:34:56.789"},
			{Kind: driver.String, String: "12:34:56+02:00"},
			{Kind: driver.Timestamp, Time: time.Date(2026, 1, 2, 12, 34, 56, 0, time.UTC), Precision: 6},
			{Kind: driver.TimestampTZ, Time: time.Date(2026, 1, 2, 12, 34, 56, 123000000, time.FixedZone("", 2*60*60)), Precision: 6},
			{Kind: driver.Bytes, Bytes: []byte{1, 2}},
			{Kind: driver.String, String: "12151fd2-7586-11e9-8f9e-2a86e4085a59"},
			{Kind: driver.String, String: "10.0.0.1"},
			{Kind: driver.String, String: "1 00:00:00.000"},
			{Kind: driver.String, String: `"{\"a\":1}"`},
			{Kind: driver.Null},
		}
		if diff := cmp.Diff(wantRow, result.Rows[0]); diff != "" {
			t.Errorf("row (-want +got):\n%s", diff)
		}
		if result.AffectedRows != nil {
			t.Errorf("AffectedRows = %d, want nil for a result set", *result.AffectedRows)
		}
		if !strings.HasPrefix(result.Statement, "SELECT") {
			t.Errorf("Statement = %q", result.Statement)
		}
	})

	t.Run("change and query statements", func(t *testing.T) {
		d := openDriver(t)
		mustExec(t, d, "CREATE SCHEMA memory.omni_query;\nCREATE TABLE memory.omni_query.items (id integer, note varchar);")
		session := openSession(t, d)
		results, err := session.Query(ctx, "INSERT INTO memory.omni_query.items VALUES (1, 'a');\nINSERT INTO memory.omni_query.items VALUES (2, 'b');\nSELECT id, note FROM memory.omni_query.items ORDER BY id;", driver.QueryOptions{})
		if err != nil {
			t.Fatal(err)
		}
		if len(results) != 3 {
			t.Fatalf("results = %d, want 3", len(results))
		}
		for i, result := range results[:2] {
			if result.Err != nil {
				t.Fatalf("statement %d: %v", i, result.Err)
			}
			if result.AffectedRows == nil || *result.AffectedRows != 1 {
				t.Errorf("statement %d affected %v, want 1", i, result.AffectedRows)
			}
			if len(result.Columns) != 0 || len(result.Rows) != 0 {
				t.Errorf("statement %d returned a result set: %+v", i, result)
			}
		}
		last := results[2]
		if last.Err != nil {
			t.Fatal(last.Err)
		}
		if last.AffectedRows != nil {
			t.Errorf("SELECT reported affected rows %d", *last.AffectedRows)
		}
		want := [][]driver.Value{
			{{Kind: driver.Int64, Int64: 1}, {Kind: driver.String, String: "a"}},
			{{Kind: driver.Int64, Int64: 2}, {Kind: driver.String, String: "b"}},
		}
		if diff := cmp.Diff(want, last.Rows); diff != "" {
			t.Errorf("rows (-want +got):\n%s", diff)
		}
	})

	t.Run("failure stops the script", func(t *testing.T) {
		d := openDriver(t)
		session := openSession(t, d)
		results, err := session.Query(ctx, "SELECT 1;\nSELECT * FROM memory.default.no_such_table;\nSELECT 2;", driver.QueryOptions{})
		if err != nil {
			t.Fatalf("Query returned a script error: %v", err)
		}
		if len(results) != 2 {
			t.Fatalf("results = %d, want 2 (the script stops at the failure)", len(results))
		}
		if results[0].Err != nil {
			t.Fatalf("first statement failed: %v", results[0].Err)
		}
		if results[1].Err == nil {
			t.Fatal("second statement reported no error")
		}
		if !strings.Contains(results[1].Err.Error(), "does not exist") {
			t.Errorf("error = %v, want the server's table-not-found error", results[1].Err)
		}
	})

	t.Run("limit", func(t *testing.T) {
		d := openDriver(t)
		session := openSession(t, d)
		results, err := session.Query(ctx, "SELECT x FROM UNNEST(SEQUENCE(1, 5)) AS t(x);\nSELECT x FROM UNNEST(SEQUENCE(1, 5)) AS t(x) LIMIT 3;", driver.QueryOptions{Limit: 2})
		if err != nil {
			t.Fatal(err)
		}
		if len(results) != 2 {
			t.Fatalf("results = %d", len(results))
		}
		if !strings.HasSuffix(results[0].Statement, " LIMIT 2") {
			t.Errorf("Statement = %q, want the limit appended", results[0].Statement)
		}
		if len(results[0].Rows) != 2 {
			t.Errorf("rows = %d, want 2", len(results[0].Rows))
		}
		if strings.Contains(results[1].Statement, "LIMIT 2") {
			t.Errorf("Statement = %q, want the caller's own LIMIT kept", results[1].Statement)
		}
		if len(results[1].Rows) != 3 {
			t.Errorf("rows = %d, want 3", len(results[1].Rows))
		}
	})

	t.Run("limit goes before a trailing comment", func(t *testing.T) {
		d := openDriver(t)
		session := openSession(t, d)
		results, err := session.Query(ctx, "SELECT x FROM UNNEST(SEQUENCE(1, 5)) AS t(x) -- note", driver.QueryOptions{Limit: 2})
		if err != nil {
			t.Fatal(err)
		}
		if len(results) != 1 || results[0].Err != nil {
			t.Fatalf("results = %+v", results)
		}
		if !strings.Contains(results[0].Statement, " LIMIT 2 -- note") {
			t.Errorf("Statement = %q, want the limit before the trailing comment", results[0].Statement)
		}
		if len(results[0].Rows) != 2 {
			t.Errorf("rows = %d, want 2", len(results[0].Rows))
		}
	})

	t.Run("explain", func(t *testing.T) {
		d := openDriver(t)
		session := openSession(t, d)
		for _, tc := range []struct {
			name, format, wantPrefix string
		}{
			{"text", "TEXT", "EXPLAIN SELECT"},
			{"json", "JSON", "EXPLAIN (FORMAT JSON) SELECT"},
		} {
			t.Run(tc.name, func(t *testing.T) {
				results, err := session.Query(ctx, "SELECT 1", driver.QueryOptions{Explain: true, ExplainFormat: tc.format, Limit: 10})
				if err != nil {
					t.Fatal(err)
				}
				if len(results) != 1 || results[0].Err != nil {
					t.Fatalf("results = %+v", results)
				}
				if !strings.HasPrefix(results[0].Statement, tc.wantPrefix) {
					t.Errorf("Statement = %q, want prefix %q", results[0].Statement, tc.wantPrefix)
				}
				if strings.Contains(results[0].Statement, "LIMIT") {
					t.Errorf("Statement = %q, want no limit rewrite on a plan", results[0].Statement)
				}
				if len(results[0].Rows) == 0 {
					t.Error("the plan returned no rows")
				}
			})
		}
	})

	t.Run("explain analyze is refused", func(t *testing.T) {
		d := openDriver(t)
		session := openSession(t, d)
		results, err := session.Query(ctx, "EXPLAIN ANALYZE SELECT 1", driver.QueryOptions{Explain: true})
		if err != nil {
			t.Fatal(err)
		}
		if len(results) != 1 || results[0].Err == nil {
			t.Fatalf("results = %+v, want one failed result", results)
		}
		if !strings.Contains(results[0].Err.Error(), "EXPLAIN ANALYZE executes the statement it plans") {
			t.Errorf("error = %v", results[0].Err)
		}
	})

	t.Run("schema", func(t *testing.T) {
		d := openDriver(t)
		mustExec(t, d, "CREATE SCHEMA memory.omni_query_scoped;\nCREATE TABLE memory.omni_query_scoped.scoped (id integer);\nINSERT INTO memory.omni_query_scoped.scoped VALUES (7);")
		session := openSession(t, d)
		results, err := session.Query(ctx, "SELECT id FROM scoped", driver.QueryOptions{Schema: "omni_query_scoped"})
		if err != nil {
			t.Fatal(err)
		}
		if len(results) != 1 || results[0].Err != nil {
			t.Fatalf("results = %+v", results)
		}
		if diff := cmp.Diff([][]driver.Value{{{Kind: driver.Int64, Int64: 7}}}, results[0].Rows); diff != "" {
			t.Errorf("rows (-want +got):\n%s", diff)
		}
	})

	t.Run("a leading comment does not change the statement kind", func(t *testing.T) {
		d := openDriver(t)
		session := openSession(t, d)
		results, err := session.Query(ctx, "-- a comment\n/* another */ SELECT 1", driver.QueryOptions{})
		if err != nil {
			t.Fatal(err)
		}
		if len(results) != 1 || results[0].Err != nil {
			t.Fatalf("results = %+v", results)
		}
		if results[0].AffectedRows != nil {
			t.Fatalf("the statement ran as a change script: %+v", results[0])
		}
		if len(results[0].Rows) != 1 {
			t.Fatalf("rows = %d, want 1", len(results[0].Rows))
		}
	})

	// TestQuerysWithAqueryThatStartsWithWith pins a defect carried verbatim:
	// Bytebase decides a statement returns rows by its leading keyword, so a
	// WITH query runs on the exec path and its rows are discarded.
	t.Run("a WITH query loses its rows", func(t *testing.T) {
		d := openDriver(t)
		session := openSession(t, d)
		results, err := session.Query(ctx, "WITH t AS (SELECT 1 AS x) SELECT * FROM t", driver.QueryOptions{})
		if err != nil {
			t.Fatal(err)
		}
		if len(results) != 1 || results[0].Err != nil {
			t.Fatalf("results = %+v", results)
		}
		if results[0].AffectedRows == nil || *results[0].AffectedRows != 0 {
			t.Fatalf("result = %+v, want the exec path's affected-rows result", results[0])
		}
		if len(results[0].Rows) != 0 {
			t.Fatalf("a WITH query returned rows; if it now does, fix(trino) landed: %+v", results[0])
		}
	})

	// TestQueryShowWithALimitFails pins a defect carried verbatim: the limit
	// rewrite appends LIMIT to a statement that does not take one.
	t.Run("SHOW does not take a limit", func(t *testing.T) {
		d := openDriver(t)
		session := openSession(t, d)
		results, err := session.Query(ctx, "SHOW SCHEMAS", driver.QueryOptions{Limit: 10})
		if err != nil {
			t.Fatal(err)
		}
		if len(results) != 1 || results[0].Err == nil {
			t.Fatalf("results = %+v, want one failed result", results)
		}
		if !strings.Contains(results[0].Err.Error(), "mismatched input 'LIMIT'") {
			t.Fatalf("error = %v; if SHOW now takes the limit, fix(trino) landed", results[0].Err)
		}
	})

	t.Run("result size cap", func(t *testing.T) {
		d := openDriver(t)
		session := openSession(t, d)
		results, err := session.Query(ctx, "SELECT x FROM UNNEST(SEQUENCE(1, 100)) AS t(x)", driver.QueryOptions{MaxResultBytes: 100})
		if err != nil {
			t.Fatal(err)
		}
		if len(results) != 1 || results[0].Err != nil {
			t.Fatalf("results = %+v", results)
		}
		if !results[0].Truncated {
			t.Fatal("Truncated = false, want the cap to stop the read")
		}
		if len(results[0].Rows) == 0 || len(results[0].Rows) >= 100 {
			t.Errorf("rows = %d, want some but not all", len(results[0].Rows))
		}
	})

	t.Run("ping", func(t *testing.T) {
		d := openDriver(t)
		if err := d.Ping(ctx); err != nil {
			t.Fatal(err)
		}
	})
}

// TestExecute runs change scripts and holds what each step reported.
func TestExecute(t *testing.T) {
	ctx := context.Background()

	t.Run("counts and steps", func(t *testing.T) {
		d := openDriver(t)
		mustExec(t, d, "CREATE SCHEMA memory.omni_execute;\nCREATE TABLE memory.omni_execute.items (id integer, note varchar);")
		log := &stepLogger{}
		script := "INSERT INTO memory.omni_execute.items VALUES (1, 'a');\n\nINSERT INTO memory.omni_execute.items VALUES (2, 'b');"
		affected, err := d.Execute(ctx, script, driver.ExecuteOptions{Log: log})
		if err != nil {
			t.Fatal(err)
		}
		if affected != 2 {
			t.Errorf("affected = %d, want 2", affected)
		}
		first, second := "INSERT INTO memory.omni_execute.items VALUES (1, 'a');", "INSERT INTO memory.omni_execute.items VALUES (2, 'b');"
		secondStart := len(first) + 2
		wantSteps := []string{
			fmt.Sprintf("exec[0:%d]", len(first)-1), "resp:1:[1]",
			fmt.Sprintf("exec[%d:%d]", secondStart, secondStart+len(second)-1), "resp:1:[1]",
		}
		if !slices.Equal(log.steps, wantSteps) {
			t.Errorf("steps = %v, want %v", log.steps, wantSteps)
		}
		for i, command := range log.commands {
			if script[command.ByteStart:command.ByteEnd] != command.Text {
				t.Errorf("command %d range [%d:%d] does not hold its text", i, command.ByteStart, command.ByteEnd)
			}
		}
	})

	t.Run("failure keeps the rows before it", func(t *testing.T) {
		d := openDriver(t)
		mustExec(t, d, "CREATE SCHEMA memory.omni_execute_failure;\nCREATE TABLE memory.omni_execute_failure.items (id integer);")
		log := &stepLogger{}
		affected, err := d.Execute(ctx, "INSERT INTO memory.omni_execute_failure.items VALUES (1);\nINSERT INTO memory.omni_execute_failure.missing VALUES (2);", driver.ExecuteOptions{Log: log})
		if err == nil || !strings.Contains(err.Error(), "failed to execute statement") {
			t.Fatalf("err = %v, want the wrapped execute failure", err)
		}
		if affected != 1 {
			t.Errorf("affected = %d, want the first statement's row", affected)
		}
		if len(log.responses) != 2 || log.responses[1].err == nil {
			t.Errorf("responses = %+v, want the second carrying the error", log.responses)
		}
		if !strings.Contains(log.responses[1].err.Error(), "does not exist") {
			t.Errorf("logged error = %v, want the server's error", log.responses[1].err)
		}
	})

	t.Run("transaction mode is ignored", func(t *testing.T) {
		d := openDriver(t)
		mustExec(t, d, "CREATE SCHEMA memory.omni_execute_tx;\nCREATE TABLE memory.omni_execute_tx.items (id integer);")
		log := &stepLogger{}
		if _, err := d.Execute(ctx, "INSERT INTO memory.omni_execute_tx.items VALUES (1)", driver.ExecuteOptions{
			Transaction: driver.Transactional,
			Isolation:   driver.Serializable,
			Log:         log,
		}); err != nil {
			t.Fatal(err)
		}
		for _, step := range log.steps {
			if strings.HasPrefix(step, "tx:") {
				t.Errorf("steps = %v, want no transaction control", log.steps)
			}
		}
	})

	t.Run("empty script", func(t *testing.T) {
		d := openDriver(t)
		log := &stepLogger{}
		affected, err := d.Execute(ctx, ";;\n-- nothing\n;", driver.ExecuteOptions{Log: log})
		if err != nil || affected != 0 {
			t.Fatalf("Execute = (%d, %v)", affected, err)
		}
		if len(log.steps) != 0 {
			t.Errorf("steps = %v, want none", log.steps)
		}
	})
}

func openSession(t *testing.T, d *Driver) driver.Session {
	t.Helper()
	session, err := d.Session(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { session.Close() })
	return session
}

// stepLogger records what Execute reports, in order.
type stepLogger struct {
	steps     []string
	commands  []driver.Command
	responses []stepResponse
}

type stepResponse struct {
	affected  int64
	perResult []int64
	err       error
}

func (l *stepLogger) CommandExecute(command driver.Command) {
	l.commands = append(l.commands, command)
	l.steps = append(l.steps, fmt.Sprintf("exec[%d:%d]", command.ByteStart, command.ByteEnd))
}

func (l *stepLogger) CommandResponse(affected int64, perResult []int64, err error) {
	l.responses = append(l.responses, stepResponse{affected: affected, perResult: perResult, err: err})
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

func (l *stepLogger) Retry(err error, attempt, max int) {
	l.steps = append(l.steps, fmt.Sprintf("retry:%d/%d", attempt, max))
}
