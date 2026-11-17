package driver

import (
	"context"
	"strings"
	"testing"

	"github.com/bytebase/omni/driver"
)

func TestQuery(t *testing.T) {
	d := openTestDriver(t)
	mustExec(t, d.db, `
CREATE SCHEMA app;
CREATE TABLE app.people (
    id INT8 PRIMARY KEY,
    name STRING NOT NULL,
    active BOOL,
    amount DECIMAL(10, 2),
    ratio FLOAT8,
    born DATE,
    at TIMESTAMP,
    at_tz TIMESTAMPTZ
);
INSERT INTO app.people VALUES
    (1, 'one', true, 10.50, 1.5, '2020-01-02', '2024-01-01 00:00:00', '2024-01-01 00:00:00+00'),
    (2, 'two', false, 20.25, 2.5, '2021-02-03', '2024-01-02 00:00:00', '2024-01-02 00:00:00+00');
CREATE SCHEMA other;
CREATE TABLE other.thing (v INT8);
INSERT INTO other.thing VALUES (7);`)

	ctx := context.Background()
	session, err := d.Session(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	t.Run("rows and types", func(t *testing.T) {
		results, err := session.Query(ctx, "SELECT id, name, active, amount, ratio, born, at, at_tz FROM app.people ORDER BY id;", driver.QueryOptions{})
		if err != nil {
			t.Fatal(err)
		}
		if len(results) != 1 {
			t.Fatalf("results = %d, want 1", len(results))
		}
		res := results[0]
		if res.Err != nil {
			t.Fatal(res.Err)
		}
		if len(res.Columns) != 8 || len(res.Rows) != 2 {
			t.Fatalf("columns/rows = %d/%d, want 8/2", len(res.Columns), len(res.Rows))
		}
		if got := res.Rows[0][0]; got.Kind != driver.Int64 || got.Int64 != 1 {
			t.Errorf("id = %+v, want int64 1", got)
		}
		if got := res.Rows[0][1]; got.Kind != driver.String || got.String != "one" {
			t.Errorf("name = %+v, want string one", got)
		}
		if got := res.Rows[0][2]; got.Kind != driver.Bool || !got.Bool {
			t.Errorf("active = %+v, want bool true", got)
		}
		if got := res.Rows[0][3]; got.Kind != driver.String || got.String != "10.50" {
			t.Errorf("amount = %+v, want string 10.50", got)
		}
		if got := res.Rows[0][4]; got.Kind != driver.Float64 || got.Float64 != 1.5 {
			t.Errorf("ratio = %+v, want double 1.5", got)
		}
		if got := res.Rows[0][5]; got.Kind != driver.String || got.String != "2020-01-02" {
			t.Errorf("born = %+v, want string 2020-01-02", got)
		}
		if got := res.Rows[0][6]; got.Kind != driver.Timestamp {
			t.Errorf("at = %+v, want timestamp", got)
		}
		if got := res.Rows[0][7]; got.Kind != driver.TimestampTZ {
			t.Errorf("at_tz = %+v, want timestamp with time zone", got)
		}
	})

	t.Run("limit", func(t *testing.T) {
		results, err := session.Query(ctx, "SELECT id, name FROM app.people ORDER BY id;", driver.QueryOptions{Limit: 1})
		if err != nil {
			t.Fatal(err)
		}
		if len(results) != 1 || results[0].Err != nil {
			t.Fatalf("results = %+v", results)
		}
		if len(results[0].Rows) != 1 {
			t.Fatalf("rows = %d, want 1", len(results[0].Rows))
		}
		if !strings.Contains(strings.ToUpper(results[0].Statement), "LIMIT 1") {
			t.Errorf("statement = %q, want a LIMIT 1 rewrite", results[0].Statement)
		}
	})

	t.Run("explain json", func(t *testing.T) {
		results, err := session.Query(ctx, "SELECT id FROM app.people;", driver.QueryOptions{Explain: true, ExplainFormat: "JSON"})
		if err != nil {
			t.Fatal(err)
		}
		if len(results) != 1 || results[0].Err != nil {
			t.Fatalf("results = %+v", results)
		}
		res := results[0]
		if len(res.Columns) != 1 || res.Columns[0].Name != "plan" {
			t.Fatalf("columns = %+v, want one plan column", res.Columns)
		}
		if len(res.Rows) != 1 || res.Rows[0][0].Kind != driver.String || !strings.Contains(res.Rows[0][0].String, "processors") {
			t.Fatalf("plan = %+v, want the decoded DISTSQL JSON", res.Rows)
		}
	})

	t.Run("schema", func(t *testing.T) {
		results, err := session.Query(ctx, "SELECT v FROM thing;", driver.QueryOptions{Schema: "other"})
		if err != nil {
			t.Fatal(err)
		}
		if len(results) != 1 || results[0].Err != nil {
			t.Fatalf("results = %+v", results)
		}
		if got := results[0].Rows[0][0]; got.Kind != driver.Int64 || got.Int64 != 7 {
			t.Errorf("v = %+v, want int64 7", got)
		}
	})

	t.Run("bytes", func(t *testing.T) {
		results, err := session.Query(ctx, "SELECT 'abc'::BYTES AS b, ''::BYTES AS empty;", driver.QueryOptions{})
		if err != nil {
			t.Fatal(err)
		}
		if len(results) != 1 || results[0].Err != nil {
			t.Fatalf("results = %+v", results)
		}
		got := results[0].Rows[0]
		if got[0].Kind != driver.Bytes || string(got[0].Bytes) != "abc" {
			t.Errorf("b = %+v, want bytes abc", got[0])
		}
		if got[1].Kind != driver.Bytes || len(got[1].Bytes) != 0 {
			t.Errorf("empty = %+v, want empty bytes", got[1])
		}
	})

	t.Run("affected rows", func(t *testing.T) {
		results, err := session.Query(ctx, "UPDATE app.people SET name = name WHERE id = 1;", driver.QueryOptions{})
		if err != nil {
			t.Fatal(err)
		}
		if len(results) != 1 || results[0].Err != nil {
			t.Fatalf("results = %+v", results)
		}
		if results[0].AffectedRows == nil || *results[0].AffectedRows != 1 {
			t.Errorf("affected rows = %v, want 1", results[0].AffectedRows)
		}
	})

	t.Run("returning", func(t *testing.T) {
		results, err := session.Query(ctx, "INSERT INTO app.people (id, name) VALUES (3, 'three') RETURNING id, name;", driver.QueryOptions{})
		if err != nil {
			t.Fatal(err)
		}
		if len(results) != 1 || results[0].Err != nil {
			t.Fatalf("results = %+v", results)
		}
		res := results[0]
		if len(res.Columns) != 2 || len(res.Rows) != 1 {
			t.Fatalf("columns/rows = %d/%d, want 2/1", len(res.Columns), len(res.Rows))
		}
		if got := res.Rows[0][0]; got.Kind != driver.Int64 || got.Int64 != 3 {
			t.Errorf("id = %+v, want int64 3", got)
		}
		if got := res.Rows[0][1]; got.Kind != driver.String || got.String != "three" {
			t.Errorf("name = %+v, want string three", got)
		}
	})

	t.Run("mutating CTE", func(t *testing.T) {
		results, err := session.Query(ctx, "WITH d AS (DELETE FROM app.people WHERE id = 3 RETURNING id) SELECT count(*) FROM d;", driver.QueryOptions{})
		if err != nil {
			t.Fatal(err)
		}
		if len(results) != 1 || results[0].Err != nil {
			t.Fatalf("results = %+v", results)
		}
		if len(results[0].Columns) != 1 || len(results[0].Rows) != 1 {
			t.Fatalf("columns/rows = %+v/%+v, want one of each", results[0].Columns, results[0].Rows)
		}
		if got := results[0].Rows[0][0]; got.Kind != driver.Int64 || got.Int64 != 1 {
			t.Errorf("count = %+v, want int64 1", got)
		}
	})

	t.Run("error stops the script", func(t *testing.T) {
		results, err := session.Query(ctx, "SELECT 1; SELECT * FROM app.missing; SELECT 3;", driver.QueryOptions{})
		if err != nil {
			t.Fatal(err)
		}
		if len(results) != 2 {
			t.Fatalf("results = %d, want 2", len(results))
		}
		if results[0].Err != nil {
			t.Errorf("first statement: %v", results[0].Err)
		}
		if results[1].Err == nil || !strings.Contains(results[1].Err.Error(), "missing") {
			t.Errorf("second statement error = %v, want the missing table", results[1].Err)
		}
	})

	t.Run("time zone", func(t *testing.T) {
		results, err := session.Query(ctx, "SET TIME ZONE 'America/New_York'; SELECT at_tz FROM app.people WHERE id = 1;", driver.QueryOptions{})
		if err != nil {
			t.Fatal(err)
		}
		if len(results) != 2 || results[1].Err != nil {
			t.Fatalf("results = %+v", results)
		}
		got := results[1].Rows[0][0]
		if got.Kind != driver.TimestampTZ {
			t.Fatalf("at_tz = %+v, want timestamp with time zone", got)
		}
		zone, offset := got.Time.Zone()
		if zone != "EST" && zone != "EDT" {
			t.Errorf("zone = %q, want America/New_York's abbreviation", zone)
		}
		if offset != -5*60*60 && offset != -4*60*60 {
			t.Errorf("offset = %d, want -5h or -4h", offset)
		}
	})
}

func TestPadZeroes(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want string
	}{
		{"00:00:01.5", "00:00:01.500000"},
		{"-00:00:01.5", "-00:00:01.500000"},
		{"00:00:01.5+02", "00:00:01.500000+02"},
		{"-00:00:01.5-02", "-00:00:01.500000-02"},
		{"00:00:01", "00:00:01"},
	} {
		if got := padZeroes(tc.in, 6); got != tc.want {
			t.Errorf("padZeroes(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
