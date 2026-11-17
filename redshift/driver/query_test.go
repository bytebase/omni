package driver

import (
	"database/sql"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/bytebase/omni/driver"
)

// TestClassifyQuery pins which statements a Session runs as queries, and so
// which it runs as commands that report affected rows. The decision is
// redshift/analysis's; this test holds the wiring to it.
func TestClassifyQuery(t *testing.T) {
	for _, tc := range []struct {
		statement string
		want      bool
	}{
		{statement: "SELECT 1", want: true},
		{statement: "SELECT * FROM t", want: true},
		{statement: "WITH c AS (SELECT 1 AS n) SELECT n FROM c", want: true},
		{statement: "SHOW TABLES", want: true},
		{statement: "EXPLAIN SELECT 1", want: true},
		{statement: "SET search_path TO app", want: false},
		{statement: "INSERT INTO t VALUES (1)", want: false},
		{statement: "UPDATE t SET a = 1", want: false},
		{statement: "DELETE FROM t", want: false},
		{statement: "CREATE TABLE t (a int)", want: false},
		// EXPLAIN ANALYZE plans by running the statement, so it reports no rows
		// of its own and runs as a command.
		{statement: "EXPLAIN ANALYZE SELECT 1", want: false},
		// A CTE that writes is not a read.
		{statement: "WITH c AS (DELETE FROM t RETURNING a) SELECT a FROM c", want: false},
		{statement: "SELECT * INTO t2 FROM t", want: false},
	} {
		t.Run(tc.statement, func(t *testing.T) {
			got, err := classifyQuery(tc.statement)
			if err != nil {
				t.Fatalf("classifyQuery: %v", err)
			}
			if got != tc.want {
				t.Errorf("classifyQuery(%q) = %v, want %v", tc.statement, got, tc.want)
			}
		})
	}
}

// TestClassifyQueryReportsSyntaxErrors checks that a statement the parser
// refuses stops the query instead of being classified from a guess.
func TestClassifyQueryReportsSyntaxErrors(t *testing.T) {
	if _, err := classifyQuery("SELECT FROM WHERE"); err == nil {
		t.Fatal("classifyQuery accepted invalid SQL")
	}
}

// TestMakeValueByTypeName pins the scan destination each of Redshift's type
// names gets. A name missing here falls back to a string destination, which is
// how an unknown type is rendered.
func TestMakeValueByTypeName(t *testing.T) {
	for _, tc := range []struct {
		typeName string
		want     any
	}{
		{typeName: "VARCHAR", want: new(sql.NullString)},
		{typeName: "TEXT", want: new(sql.NullString)},
		{typeName: "UUID", want: new(sql.NullString)},
		{typeName: "BOOL", want: new(sql.NullBool)},
		{typeName: "SMALLINT", want: new(sql.NullInt64)},
		{typeName: "INT4", want: new(sql.NullInt64)},
		{typeName: "BIGINT", want: new(sql.NullInt64)},
		{typeName: "DOUBLE PRECISION", want: new(sql.NullString)},
		{typeName: "FLOAT8", want: new(sql.NullFloat64)},
		{typeName: "DATE", want: new(pgtype.Date)},
		{typeName: "TIMESTAMP", want: new(pgtype.Timestamptz)},
		{typeName: "TIMESTAMPTZ", want: new(pgtype.Timestamptz)},
		{typeName: "VARBIT", want: new([]byte)},
		{typeName: "SUPER", want: new(sql.NullString)},
	} {
		t.Run(tc.typeName, func(t *testing.T) {
			got := makeValueByTypeName(tc.typeName)
			if typeOf(got) != typeOf(tc.want) {
				t.Errorf("makeValueByTypeName(%q) = %T, want %T", tc.typeName, got, tc.want)
			}
		})
	}
}

// typeOf names a scan destination's type, so a test can compare two fresh
// destinations without comparing their pointers.
func typeOf(v any) string {
	switch v.(type) {
	case *sql.NullString:
		return "*sql.NullString"
	case *sql.NullBool:
		return "*sql.NullBool"
	case *sql.NullInt64:
		return "*sql.NullInt64"
	case *sql.NullFloat64:
		return "*sql.NullFloat64"
	case *pgtype.Date:
		return "*pgtype.Date"
	case *pgtype.Timestamptz:
		return "*pgtype.Timestamptz"
	case *[]byte:
		return "*[]byte"
	default:
		return "other"
	}
}

// TestConvertValue pins the contract value each scanned cell becomes, including
// the two Redshift-specific readings: a TIME, TIMETZ or INTERVAL keeps six
// fractional digits, and an empty BIT value is NULL.
func TestConvertValue(t *testing.T) {
	location := time.FixedZone("", 2*60*60)
	textMeta := columnMeta{databaseTypeName: "VARCHAR"}
	timestampMeta := columnMeta{databaseTypeName: "TIMESTAMP", scale: 3, hasScale: true}
	unscaledMeta := columnMeta{databaseTypeName: "TIMESTAMPTZ"}

	for _, tc := range []struct {
		name     string
		typeName string
		meta     columnMeta
		value    any
		want     driver.Value
	}{
		{
			name:     "a string",
			typeName: "VARCHAR",
			meta:     textMeta,
			value:    &sql.NullString{String: "hello", Valid: true},
			want:     driver.Value{Kind: driver.String, String: "hello"},
		},
		{
			name:     "a NULL string",
			typeName: "VARCHAR",
			meta:     textMeta,
			value:    &sql.NullString{},
			want:     driver.Value{Kind: driver.Null},
		},
		{
			name:     "a TIME keeps six fractional digits",
			typeName: "TIME",
			meta:     columnMeta{databaseTypeName: "TIME"},
			value:    &sql.NullString{String: "12:34:56.123", Valid: true},
			want:     driver.Value{Kind: driver.String, String: "12:34:56.123000"},
		},
		{
			name:     "a TIMETZ keeps six and its offset",
			typeName: "TIMETZ",
			meta:     columnMeta{databaseTypeName: timeTzOID},
			value:    &sql.NullString{String: "12:34:56.1+02", Valid: true},
			want:     driver.Value{Kind: driver.String, String: "12:34:56.100000+02"},
		},
		{
			name:     "an INTERVAL keeps six",
			typeName: "INTERVAL",
			meta:     columnMeta{databaseTypeName: "INTERVAL"},
			value:    &sql.NullString{String: "-00:02:45.25", Valid: true},
			want:     driver.Value{Kind: driver.String, String: "-00:02:45.250000"},
		},
		{
			name:     "an integer",
			typeName: "INT8",
			meta:     columnMeta{databaseTypeName: "INT8"},
			value:    &sql.NullInt64{Int64: 42, Valid: true},
			want:     driver.Value{Kind: driver.Int64, Int64: 42},
		},
		{
			name:     "a NULL integer",
			typeName: "INT8",
			meta:     columnMeta{databaseTypeName: "INT8"},
			value:    &sql.NullInt64{},
			want:     driver.Value{Kind: driver.Null},
		},
		{
			name:     "a boolean",
			typeName: "BOOL",
			meta:     columnMeta{databaseTypeName: "BOOL"},
			value:    &sql.NullBool{Bool: true, Valid: true},
			want:     driver.Value{Kind: driver.Bool, Bool: true},
		},
		{
			name:     "a double",
			typeName: "FLOAT8",
			meta:     columnMeta{databaseTypeName: "FLOAT8"},
			value:    &sql.NullFloat64{Float64: 1.5, Valid: true},
			want:     driver.Value{Kind: driver.Float64, Float64: 1.5},
		},
		{
			name:     "a bit value",
			typeName: "VARBIT",
			meta:     columnMeta{databaseTypeName: "VARBIT"},
			value:    &[]byte{0x01, 0x02},
			want:     driver.Value{Kind: driver.Bytes, Bytes: []byte{0x01, 0x02}},
		},
		{
			name:     "an empty bit value is NULL",
			typeName: "VARBIT",
			meta:     columnMeta{databaseTypeName: "VARBIT"},
			value:    &[]byte{},
			want:     driver.Value{Kind: driver.Null},
		},
		{
			name:     "a date",
			typeName: "DATE",
			meta:     columnMeta{databaseTypeName: "DATE"},
			value:    &pgtype.Date{Time: time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC), Valid: true},
			want:     driver.Value{Kind: driver.String, String: "2026-03-04"},
		},
		{
			name:     "an infinite date",
			typeName: "DATE",
			meta:     columnMeta{databaseTypeName: "DATE"},
			value:    &pgtype.Date{InfinityModifier: pgtype.Infinity, Valid: true},
			want:     driver.Value{Kind: driver.String, String: pgtype.Infinity.String()},
		},
		{
			name:     "a timestamp keeps its declared precision",
			typeName: "TIMESTAMP",
			meta:     timestampMeta,
			value:    &pgtype.Timestamptz{Time: time.Date(2026, 3, 4, 5, 6, 7, 891234000, time.UTC), Valid: true},
			want:     driver.Value{Kind: driver.Timestamp, Time: time.Date(2026, 3, 4, 5, 6, 7, 891234000, time.UTC), Precision: 3},
		},
		{
			name:     "a timestamp of unspecified precision is microsecond",
			typeName: "TIMESTAMPTZ",
			meta:     unscaledMeta,
			value:    &pgtype.Timestamptz{Time: time.Date(2026, 3, 4, 5, 6, 7, 0, time.UTC), Valid: true},
			want:     driver.Value{Kind: driver.TimestampTZ, Time: time.Date(2026, 3, 4, 7, 6, 7, 0, location), Precision: 6},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := convertValue(tc.typeName, tc.meta, tc.value, location)
			if got.Kind != tc.want.Kind {
				t.Fatalf("Kind = %q, want %q", got.Kind, tc.want.Kind)
			}
			if got.String != tc.want.String {
				t.Errorf("String = %q, want %q", got.String, tc.want.String)
			}
			if got.Int64 != tc.want.Int64 {
				t.Errorf("Int64 = %d, want %d", got.Int64, tc.want.Int64)
			}
			if got.Bool != tc.want.Bool {
				t.Errorf("Bool = %v, want %v", got.Bool, tc.want.Bool)
			}
			if got.Float64 != tc.want.Float64 {
				t.Errorf("Float64 = %v, want %v", got.Float64, tc.want.Float64)
			}
			if string(got.Bytes) != string(tc.want.Bytes) {
				t.Errorf("Bytes = %v, want %v", got.Bytes, tc.want.Bytes)
			}
			if got.Precision != tc.want.Precision {
				t.Errorf("Precision = %d, want %d", got.Precision, tc.want.Precision)
			}
			if tc.want.Kind == driver.Timestamp || tc.want.Kind == driver.TimestampTZ {
				if !got.Time.Equal(tc.want.Time) {
					t.Errorf("Time = %v, want %v", got.Time, tc.want.Time)
				}
			}
		})
	}
}

// TestPadZeroes pins the fractional-second padding the contract asks for.
func TestPadZeroes(t *testing.T) {
	for _, tc := range []struct {
		raw  string
		acc  int
		want string
	}{
		{raw: "12:34:56", acc: 6, want: "12:34:56"},
		{raw: "12:34:56.123456", acc: 6, want: "12:34:56.123456"},
		{raw: "12:34:56.123", acc: 6, want: "12:34:56.123000"},
		{raw: "12:34:56.123456789", acc: 6, want: "12:34:56.123456789"},
		{raw: "12:34:56.123+02:00", acc: 6, want: "12:34:56.123000+02:00"},
		{raw: "12:34:56.123-02:00", acc: 6, want: "12:34:56.123000-02:00"},
		{raw: "-00:02:45.25", acc: 6, want: "-00:02:45.250000"},
	} {
		if got := padZeroes(tc.raw, tc.acc); got != tc.want {
			t.Errorf("padZeroes(%q, %d) = %q, want %q", tc.raw, tc.acc, got, tc.want)
		}
	}
}

// TestStripDatashareQualifier holds the rewrite a datashare query needs to
// what a shared object is addressed by: the database's name comes off an
// object reference and nothing else changes.
func TestStripDatashareQualifier(t *testing.T) {
	for _, tc := range []struct {
		name      string
		statement string
		want      string
	}{
		{
			name:      "a qualified table loses the database",
			statement: `SELECT * FROM shared.facts`,
			want:      `SELECT * FROM facts`,
		},
		{
			name:      "a quoted database name loses it too",
			statement: `SELECT * FROM "shared".facts JOIN shared.dims ON true`,
			want:      `SELECT * FROM facts JOIN dims ON true`,
		},
		{
			name:      "every reference in a statement is rewritten",
			statement: `SELECT * FROM shared.facts WHERE id IN (SELECT id FROM shared.dims)`,
			want:      `SELECT * FROM facts WHERE id IN (SELECT id FROM dims)`,
		},
		{
			name:      "a literal that spells the same text is left alone",
			statement: `SELECT 'shared.facts' AS name FROM shared.facts`,
			want:      `SELECT 'shared.facts' AS name FROM facts`,
		},
		{
			name:      "a comment that spells the same text is left alone",
			statement: "SELECT * -- shared.facts\nFROM shared.facts",
			want:      "SELECT * -- shared.facts\nFROM facts",
		},
		{
			name:      "a column qualifier is not a database qualifier",
			statement: `SELECT shared.id FROM facts AS shared`,
			want:      `SELECT shared.id FROM facts AS shared`,
		},
		{
			name:      "a schema-qualified name keeps its schema",
			statement: `SELECT * FROM shared.app.facts`,
			want:      `SELECT * FROM app.facts`,
		},
		{
			name:      "a schema named like the database keeps its name",
			statement: `SELECT * FROM shared.shared.facts`,
			want:      `SELECT * FROM shared.facts`,
		},
		{
			name:      "a longer identifier that ends in the database's name is left alone",
			statement: `SELECT * FROM not_shared.facts`,
			want:      `SELECT * FROM not_shared.facts`,
		},
		{
			name:      "a quoted name in another case is another database",
			statement: `SELECT * FROM "SHARED".facts`,
			want:      `SELECT * FROM "SHARED".facts`,
		},
		{
			name:      "an expression's FROM does not open a relation list",
			statement: `SELECT EXTRACT(day FROM shared.ts) FROM app.a AS shared`,
			want:      `SELECT EXTRACT(day FROM shared.ts) FROM app.a AS shared`,
		},
		{
			name:      "a relation inside a subquery is rewritten",
			statement: `SELECT z FROM app.a WHERE id IN (SELECT id FROM shared.b)`,
			want:      `SELECT z FROM app.a WHERE id IN (SELECT id FROM b)`,
		},
		{
			name:      "a function call before a relation does not hide it",
			statement: `SELECT EXTRACT(day FROM shared.ts) FROM shared.facts`,
			want:      `SELECT EXTRACT(day FROM shared.ts) FROM facts`,
		},
		{
			name:      "a relation after a subquery is still in the relation list",
			statement: `SELECT * FROM (SELECT 1) AS x, shared.t`,
			want:      `SELECT * FROM (SELECT 1) AS x, t`,
		},
		{
			name:      "a derived table's alias is not the database's name",
			statement: `SELECT p FROM (SELECT products FROM app.orders) AS shared, shared.products AS p`,
			want:      `SELECT p FROM (SELECT products FROM app.orders) AS shared, shared.products AS p`,
		},
		{
			name:      "a relation alias's SUPER column is not the database's name",
			statement: `SELECT p FROM app.orders AS shared, shared.products AS p`,
			want:      `SELECT p FROM app.orders AS shared, shared.products AS p`,
		},
		{
			// A qualifier that names an alias is left alone: the alias could be
			// a relation's column, so the database's name is not assumed.
			name:      "a bare alias's name is left alone too",
			statement: `SELECT shared.id FROM app.orders shared, shared.facts`,
			want:      `SELECT shared.id FROM app.orders shared, shared.facts`,
		},
		{
			name:      "a function argument is not a relation list",
			statement: `SELECT up.product FROM app.orders AS shared, UNNEST(shared.products) AS up(product)`,
			want:      `SELECT up.product FROM app.orders AS shared, UNNEST(shared.products) AS up(product)`,
		},
		{
			name:      "QUALIFY ends the relation list",
			statement: `SELECT shared.id FROM app.a AS shared, app.b AS b QUALIFY ROW_NUMBER() OVER (PARTITION BY shared.id ORDER BY shared.id) = 1`,
			want:      `SELECT shared.id FROM app.a AS shared, app.b AS b QUALIFY ROW_NUMBER() OVER (PARTITION BY shared.id ORDER BY shared.id) = 1`,
		},
		{
			name:      "a relation after a joined subquery is rewritten",
			statement: `SELECT * FROM shared.a JOIN (SELECT 1 AS n) AS s ON true JOIN shared.b ON true`,
			want:      `SELECT * FROM a JOIN (SELECT 1 AS n) AS s ON true JOIN b ON true`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := stripDatashareQualifier(tc.statement, "shared"); got != tc.want {
				t.Errorf("stripDatashareQualifier(%q) = %q, want %q", tc.statement, got, tc.want)
			}
		})
	}
}

// TestReturnsRowsWithMasking pins that a caller which masks columns is given no
// rows the analysis cannot trace: a procedure call's OUT and INOUT values.
func TestReturnsRowsWithMasking(t *testing.T) {
	for _, tc := range []struct {
		name    string
		stmt    string
		masking bool
		want    bool
	}{
		{name: "a query is read", stmt: "SELECT 1", want: true},
		{name: "a query is read with masking on", stmt: "SELECT 1", masking: true, want: true},
		{name: "a procedure call is read", stmt: "CALL app.load_facts(1)", want: true},
		{name: "a procedure call's rows are withheld with masking on", stmt: "CALL app.load_facts(1)", masking: true, want: false},
		{name: "a command reports rows affected", stmt: "INSERT INTO t VALUES (1)", want: false},
		{name: "a fetch is read", stmt: "FETCH NEXT FROM c", want: true},
		{name: "a fetch's rows are withheld with masking on", stmt: "FETCH NEXT FROM c", masking: true, want: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := returnsRows(tc.stmt, tc.masking)
			if err != nil {
				t.Fatalf("returnsRows(%q): %v", tc.stmt, err)
			}
			if got != tc.want {
				t.Errorf("returnsRows(%q, masking=%v) = %v, want %v", tc.stmt, tc.masking, got, tc.want)
			}
		})
	}
}
