package driver

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/driver"
)

// querySetup holds the types a tape can hold. Decimal, UUID, Array, Tuple and
// Map arrive as the client's own Go types, which the tape keeps no answer for;
// the live parity harness covers them.
const querySetup = `
CREATE TABLE q (
  id Int32,
  i8 Int8,
  i16 Int16,
  i32 Int32,
  i64 Int64,
  u8 UInt8,
  u16 UInt16,
  u32 UInt32,
  u64 UInt64,
  f32 Float32,
  f64 Float64,
  s String,
  b Bool,
  d Date,
  dt DateTime('UTC'),
  dt64 DateTime64(3, 'UTC'),
  n Nullable(Int32),
  fs FixedString(4)
) ENGINE = Memory;

INSERT INTO q VALUES
  (1, -8, -16, -32, -64, 8, 16, 32, 64, 1.5, 2.5, 'a', true,
   '2026-03-04', '2026-03-04 05:06:07', '2026-03-04 05:06:07.891', 42, 'abcd'),
  (2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, '', false,
   '1970-01-01', '1970-01-01 00:00:00', '1970-01-01 00:00:00', NULL, 'wxyz');
`

func query(t *testing.T, d *Driver, script string, opts driver.QueryOptions) []*driver.Result {
	t.Helper()
	ctx := context.Background()
	session, err := d.Session(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = session.Close() }()
	results, err := session.Query(ctx, script, opts)
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	return results
}

func setupQuery(t *testing.T) *Driver {
	t.Helper()
	d := openTestDriver(t)
	execCorpus(t, d.db, querySetup)
	return d
}

func TestQuery(t *testing.T) {
	d := setupQuery(t)
	results := query(t, d, "SELECT * FROM q ORDER BY id", driver.QueryOptions{})
	if len(results) != 1 {
		t.Fatalf("results = %d, want 1", len(results))
	}
	result := results[0]
	if result.Err != nil {
		t.Fatalf("query failed: %v", result.Err)
	}
	wantColumns := []driver.Column{
		{Name: "id", TypeName: "INT32"},
		{Name: "i8", TypeName: "INT8"},
		{Name: "i16", TypeName: "INT16"},
		{Name: "i32", TypeName: "INT32"},
		{Name: "i64", TypeName: "INT64"},
		{Name: "u8", TypeName: "UINT8"},
		{Name: "u16", TypeName: "UINT16"},
		{Name: "u32", TypeName: "UINT32"},
		{Name: "u64", TypeName: "UINT64"},
		{Name: "f32", TypeName: "FLOAT32"},
		{Name: "f64", TypeName: "FLOAT64"},
		{Name: "s", TypeName: "STRING"},
		{Name: "b", TypeName: "BOOL"},
		{Name: "d", TypeName: "DATE"},
		{Name: "dt", TypeName: "DATETIME('UTC')"},
		{Name: "dt64", TypeName: "DATETIME64(3, 'UTC')"},
		{Name: "n", TypeName: "NULLABLE(INT32)"},
		{Name: "fs", TypeName: "FIXEDSTRING(4)"},
	}
	if diff := cmp.Diff(wantColumns, result.Columns); diff != "" {
		t.Errorf("columns (-want +got):\n%s", diff)
	}
	if len(result.Rows) != 2 {
		t.Fatalf("rows = %d, want 2", len(result.Rows))
	}
	want := []driver.Value{
		{Kind: driver.Int32, Int64: 1},
		{Kind: driver.Int32, Int64: -8},
		{Kind: driver.Int32, Int64: -16},
		{Kind: driver.Int32, Int64: -32},
		{Kind: driver.Int64, Int64: -64},
		{Kind: driver.Uint32, Uint64: 8},
		{Kind: driver.Uint32, Uint64: 16},
		{Kind: driver.Uint32, Uint64: 32},
		{Kind: driver.Uint64, Uint64: 64},
		{Kind: driver.Float32, Float64: 1.5},
		{Kind: driver.Float64, Float64: 2.5},
		{Kind: driver.String, String: "a"},
		{Kind: driver.Bool, Bool: true},
		{Kind: driver.Timestamp, Time: time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)},
		{Kind: driver.TimestampTZ, Time: time.Date(2026, 3, 4, 5, 6, 7, 0, time.UTC)},
		{Kind: driver.TimestampTZ, Time: time.Date(2026, 3, 4, 5, 6, 7, 891000000, time.UTC), Precision: 3},
		{Kind: driver.Int32, Int64: 42},
		{Kind: driver.String, String: "abcd"},
	}
	if diff := cmp.Diff(want, result.Rows[0]); diff != "" {
		t.Errorf("first row (-want +got):\n%s", diff)
	}
	if result.Rows[1][16].Kind != driver.Null {
		t.Errorf("an unset nullable = %+v, want NULL", result.Rows[1][16])
	}
	if got := result.Rows[1][17]; got.Kind != driver.String || got.String != "wxyz" {
		t.Errorf("FixedString = %+v, want wxyz", got)
	}
}

func TestQueryLimit(t *testing.T) {
	d := setupQuery(t)
	results := query(t, d, "SELECT id FROM q ORDER BY id", driver.QueryOptions{Limit: 1})
	if len(results) != 1 {
		t.Fatalf("results = %d, want 1", len(results))
	}
	if got, want := results[0].Statement, "WITH result AS (SELECT id FROM q ORDER BY id) SELECT * FROM result LIMIT 1;"; got != want {
		t.Errorf("Statement = %q, want %q", got, want)
	}
	if len(results[0].Rows) != 1 || results[0].Rows[0][0].Int64 != 1 {
		t.Errorf("rows = %+v, want the first row", results[0].Rows)
	}
}

func TestQueryMetadataStatements(t *testing.T) {
	d := setupQuery(t)
	results := query(t, d, "SHOW TABLES;\nDESCRIBE TABLE q;\nEXISTS TABLE q", driver.QueryOptions{})
	if len(results) != 3 {
		t.Fatalf("results = %d, want 3", len(results))
	}
	for i, result := range results {
		if result.Err != nil {
			t.Errorf("result %d: %v", i, result.Err)
			continue
		}
		if result.AffectedRows != nil {
			t.Errorf("result %d reported affected rows instead of a result set", i)
		}
		if len(result.Rows) == 0 {
			t.Errorf("result %d returned no rows", i)
		}
	}
	if len(results[1].Rows) != 18 {
		t.Errorf("DESCRIBE returned %d rows, want the table's 18 columns", len(results[1].Rows))
	}
	if got := results[2].Rows[0][0]; got.Kind != driver.Uint32 || got.Uint64 != 1 {
		t.Errorf("EXISTS = %+v, want 1", got)
	}
}

func TestQueryLimitWithTrailingComment(t *testing.T) {
	d := setupQuery(t)
	results := query(t, d, "SELECT id FROM q -- note", driver.QueryOptions{Limit: 1})
	if len(results) != 1 || results[0].Err != nil {
		t.Fatalf("results = %+v", results)
	}
	if got, want := results[0].Statement, "WITH result AS (SELECT id FROM q) SELECT * FROM result LIMIT 1;"; got != want {
		t.Errorf("Statement = %q, want %q", got, want)
	}
	if len(results[0].Rows) != 1 || results[0].Rows[0][0].Int64 != 1 {
		t.Errorf("rows = %+v", results[0].Rows)
	}
}

func TestQueryExplain(t *testing.T) {
	d := setupQuery(t)
	results := query(t, d, "SELECT id FROM q", driver.QueryOptions{Explain: true})
	if len(results) != 1 {
		t.Fatalf("results = %d, want 1", len(results))
	}
	if got, want := results[0].Statement, "EXPLAIN SELECT id FROM q"; got != want {
		t.Errorf("Statement = %q, want %q", got, want)
	}
	if results[0].Err != nil {
		t.Fatalf("explain failed: %v", results[0].Err)
	}
	if results[0].Plan == nil || results[0].Plan.Format != "TEXT" || results[0].Plan.Executed {
		t.Errorf("Plan = %+v, want a TEXT plan that did not run", results[0].Plan)
	}
	if len(results[0].Rows) == 0 {
		t.Error("explain returned no rows")
	}
}

func TestQueryExplainJSON(t *testing.T) {
	d := setupQuery(t)
	results := query(t, d, "SELECT id FROM q", driver.QueryOptions{Explain: true, ExplainFormat: "JSON"})
	if len(results) != 1 {
		t.Fatalf("results = %d, want 1", len(results))
	}
	if got, want := results[0].Statement, "EXPLAIN json = 1, indexes = 1 SELECT id FROM q"; got != want {
		t.Errorf("Statement = %q, want %q", got, want)
	}
	if results[0].Err != nil {
		t.Fatalf("explain failed: %v", results[0].Err)
	}
	if results[0].Plan == nil || results[0].Plan.Format != "JSON" {
		t.Errorf("Plan = %+v, want JSON", results[0].Plan)
	}
	if len(results[0].Rows) == 0 {
		t.Error("explain returned no rows")
	}
}

func TestQueryExplainPreserved(t *testing.T) {
	d := setupQuery(t)
	results := query(t, d, "EXPLAIN SELECT id FROM q;\nEXPLAIN json = 1, indexes = 1 SELECT id FROM q;\nEXPLAIN PLAN json = 1 SELECT id FROM q;\nEXPLAIN json = 0 SELECT id FROM q",
		driver.QueryOptions{Explain: true, ExplainFormat: "JSON"})
	if len(results) != 4 {
		t.Fatalf("results = %d, want 4", len(results))
	}
	if results[0].Statement != "EXPLAIN SELECT id FROM q;" {
		t.Errorf("text explain was rewritten: %q", results[0].Statement)
	}
	if results[0].Plan == nil || results[0].Plan.Format != "TEXT" {
		t.Errorf("text plan = %+v, want TEXT", results[0].Plan)
	}
	if results[1].Plan == nil || results[1].Plan.Format != "JSON" {
		t.Errorf("json plan = %+v, want JSON", results[1].Plan)
	}
	if results[2].Plan == nil || results[2].Plan.Format != "JSON" {
		t.Errorf("PLAN with json = 1 = %+v, want JSON", results[2].Plan)
	}
	if results[3].Plan == nil || results[3].Plan.Format != "TEXT" {
		t.Errorf("json = 0 = %+v, want TEXT", results[3].Plan)
	}
}

func TestQueryStopsAtTheFailingStatement(t *testing.T) {
	d := setupQuery(t)
	results := query(t, d, "SELECT 1;\nSELECT * FROM no_such_table;\nSELECT 3", driver.QueryOptions{})
	if len(results) != 2 {
		t.Fatalf("results = %d, want 2", len(results))
	}
	if results[0].Err != nil {
		t.Fatalf("first statement failed: %v", results[0].Err)
	}
	if results[1].Err == nil || !strings.Contains(results[1].Err.Error(), "no_such_table") {
		t.Fatalf("second statement error = %v", results[1].Err)
	}
}

func TestQueryError(t *testing.T) {
	d := setupQuery(t)
	results := query(t, d, "SELECT * FROM no_such_table", driver.QueryOptions{})
	if len(results) != 1 || results[0].Err == nil {
		t.Fatalf("results = %+v", results)
	}
	if !strings.Contains(results[0].Err.Error(), "no_such_table") {
		t.Errorf("error = %v", results[0].Err)
	}
}

func TestQueryLimitInsert(t *testing.T) {
	d := setupQuery(t)
	results := query(t, d, "INSERT INTO q (id) VALUES (4)", driver.QueryOptions{Limit: 1})
	if len(results) != 1 || results[0].Err != nil {
		t.Fatalf("results = %+v", results)
	}
	if results[0].AffectedRows == nil {
		t.Errorf("Statement = %q, want the insert run as written", results[0].Statement)
	}
}

func TestQueryValidationErrorKeepsEarlierResults(t *testing.T) {
	d := setupQuery(t)
	results := query(t, d, "INSERT INTO q (id) VALUES (5);\nSELECT 'unterminated", driver.QueryOptions{})
	if len(results) != 2 {
		t.Fatalf("results = %d, want the insert and the refused statement", len(results))
	}
	if results[0].Err != nil || results[0].AffectedRows == nil {
		t.Errorf("first result = %+v, want the insert's affected row", results[0])
	}
	if results[1].Err == nil || !strings.Contains(results[1].Err.Error(), "invalid string") {
		t.Errorf("second error = %v, want the unterminated string", results[1].Err)
	}
	if got := countRows(t, d, "q"); got != 3 {
		t.Errorf("rows = %d, want the insert to have taken effect", got)
	}
}

func TestQueryMaxResultBytes(t *testing.T) {
	d := setupQuery(t)
	results := query(t, d, "SELECT * FROM q ORDER BY id", driver.QueryOptions{MaxResultBytes: 100})
	if len(results) != 1 {
		t.Fatalf("results = %d, want 1", len(results))
	}
	if !results[0].Truncated {
		t.Error("reading stopped at MaxResultBytes but the result does not say so")
	}
	if len(results[0].Rows) != 1 {
		t.Errorf("rows = %d, want the cap to stop at the first row", len(results[0].Rows))
	}
}

func TestQueryMultipleStatements(t *testing.T) {
	d := setupQuery(t)
	results := query(t, d, "SELECT 1;\n-- a comment\nSELECT 2", driver.QueryOptions{})
	if len(results) != 2 {
		t.Fatalf("results = %d, want 2", len(results))
	}
	for i, want := range []int64{1, 2} {
		if results[i].Err != nil {
			t.Errorf("result %d: %v", i, results[i].Err)
			continue
		}
		value := results[i].Rows[0][0]
		if value.Kind != driver.Uint32 || value.Uint64 != uint64(want) {
			t.Errorf("result %d = %+v, want the uint32 %d", i, value, want)
		}
	}
}

func TestQueryAffectedRows(t *testing.T) {
	d := setupQuery(t)
	results := query(t, d, "INSERT INTO q (id) VALUES (3)", driver.QueryOptions{})
	if len(results) != 1 {
		t.Fatalf("results = %d, want 1", len(results))
	}
	if results[0].Err != nil {
		t.Fatalf("insert failed: %v", results[0].Err)
	}
	if results[0].AffectedRows == nil {
		t.Fatal("a statement that returned no rows reports no affected count")
	}
	if diff := cmp.Diff(int64(0), *results[0].AffectedRows); diff != "" {
		t.Errorf("affected rows (-want +got):\n%s", diff)
	}
}
