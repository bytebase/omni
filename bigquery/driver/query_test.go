package driver

import (
	"math/big"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/driver"
)

// TestStatements covers the BigQuery splitter the driver wires: a ';' inside a
// literal, a quoted identifier or a comment is not a boundary, and a segment that
// holds no statement asks BigQuery nothing.
func TestStatements(t *testing.T) {
	for _, tc := range []struct {
		name   string
		script string
		want   []string
	}{
		{
			name:   "one statement",
			script: "SELECT 1;",
			want:   []string{"SELECT 1"},
		},
		{
			name:   "two statements are trimmed",
			script: "SELECT 1;\n  SELECT 2  ;\n",
			want:   []string{"SELECT 1", "SELECT 2"},
		},
		{
			name:   "a ';' inside a literal is not a boundary",
			script: "SELECT ';' AS a; SELECT 2",
			want:   []string{"SELECT ';' AS a", "SELECT 2"},
		},
		{
			name:   "a ';' inside a triple-quoted literal is not a boundary",
			script: "SELECT \"\"\";\"\"\" AS a; SELECT 2",
			want:   []string{"SELECT \"\"\";\"\"\" AS a", "SELECT 2"},
		},
		{
			name:   "a ';' inside a comment is not a boundary",
			script: "SELECT 1 # a comment holding ;\n; SELECT 2",
			want:   []string{"SELECT 1 # a comment holding ;", "SELECT 2"},
		},
		{
			name:   "a trailing comment after the last terminator asks nothing",
			script: "SELECT 1; SELECT 2; -- nothing after this",
			want:   []string{"SELECT 1", "SELECT 2"},
		},
		{
			name:   "an empty script holds no statement",
			script: "  \n\t",
			want:   nil,
		},
		{
			name:   "lone terminators ask nothing",
			script: ";;  ;",
			want:   nil,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.want == nil {
				require.Empty(t, statements(tc.script))
				return
			}
			require.Equal(t, tc.want, statements(tc.script))
		})
	}
}

// TestReturnsRows covers the classification the driver wires its query path to:
// a statement bigquery/analysis calls a query returns rows, and every other
// statement is run as a job that reports rows affected.
func TestReturnsRows(t *testing.T) {
	for _, tc := range []struct {
		name      string
		statement string
		want      bool
	}{
		{name: "a SELECT", statement: "SELECT 1", want: true},
		{name: "a query with a leading comment", statement: "-- a note\nSELECT 1", want: true},
		{name: "a parenthesized query", statement: "(SELECT 1)", want: true},
		{name: "a CTE query", statement: "WITH x AS (SELECT 1) SELECT * FROM x", want: true},
		{
			name:      "a query that reads only INFORMATION_SCHEMA",
			statement: "SELECT * FROM `example-project.analytics.INFORMATION_SCHEMA.COLUMNS`",
			want:      true,
		},
		{name: "a UNION", statement: "SELECT 1 UNION ALL SELECT 2", want: true},
		{name: "an INSERT", statement: "INSERT INTO t (a) VALUES (1)"},
		{name: "an UPDATE", statement: "UPDATE t SET a = 1 WHERE true"},
		{name: "a DELETE", statement: "DELETE FROM t WHERE true"},
		{name: "a MERGE", statement: "MERGE t USING s ON false WHEN NOT MATCHED THEN INSERT (a) VALUES (1)"},
		{name: "a CTE-prefixed DML", statement: "WITH x AS (SELECT 1) INSERT INTO t SELECT * FROM x"},
		{name: "a CREATE TABLE", statement: "CREATE TABLE t (a INT64)"},
		{name: "an ALTER TABLE", statement: "ALTER TABLE t ADD COLUMN b INT64"},
		{name: "a DROP TABLE", statement: "DROP TABLE t"},
		{name: "an EXPORT DATA", statement: "EXPORT DATA OPTIONS(uri='gs://bucket/*') AS SELECT 1"},
		{name: "a SET", statement: "SET x = 1"},
		{name: "a DECLARE", statement: "DECLARE x INT64"},
		{name: "an empty statement", statement: ""},
		{name: "text that does not parse", statement: "SELEC 1"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, returnsRows(tc.statement))
		})
	}
}

// TestLimitStatement covers the rewrite the driver applies to a statement it
// classified as a query: BigQuery's CTE wrap, which is
// bigquery/analysis.StatementWithResultLimit.
func TestLimitStatement(t *testing.T) {
	require.Equal(t,
		"WITH result AS (SELECT 1) SELECT * FROM result LIMIT 10;",
		limitStatement("SELECT 1", 10),
	)
	require.Equal(t,
		"WITH result AS (SELECT 1) SELECT * FROM result LIMIT 10;",
		limitStatement("SELECT 1;", 10),
	)
	// A trailing comment would comment out the closing parenthesis, so the wrap
	// starts its own SQL on a new line.
	require.Equal(t,
		"WITH result AS (SELECT 1 -- a note\n) SELECT * FROM result LIMIT 10;",
		limitStatement("SELECT 1 -- a note", 10),
	)
	// A limit of zero is no cap, and the wrapper then caps nothing.
	require.Equal(t,
		"WITH result AS (SELECT 1) SELECT * FROM result;",
		limitStatement("SELECT 1", 0),
	)
}

// TestConvertValue covers every kind the driver maps and every type it falls back
// on. A cell whose declared type has no kind is Go's own rendering of the value,
// which is what Bytebase's driver returned.
func TestConvertValue(t *testing.T) {
	instant := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)

	for _, tc := range []struct {
		name      string
		value     bigquery.Value
		fieldType bigquery.FieldType
		want      driver.Value
	}{
		{
			name:      "a NULL",
			value:     nil,
			fieldType: bigquery.StringFieldType,
			want:      driver.Value{Kind: driver.Null},
		},
		{
			name:      "a STRING",
			value:     "hello",
			fieldType: bigquery.StringFieldType,
			want:      driver.Value{Kind: driver.String, String: "hello"},
		},
		{
			name:      "BYTES",
			value:     []byte{0x01, 0x02},
			fieldType: bigquery.BytesFieldType,
			want:      driver.Value{Kind: driver.Bytes, Bytes: []byte{0x01, 0x02}},
		},
		{
			name:      "an INTEGER",
			value:     int64(7),
			fieldType: bigquery.IntegerFieldType,
			want:      driver.Value{Kind: driver.Int64, Int64: 7},
		},
		{
			name:      "a FLOAT",
			value:     3.5,
			fieldType: bigquery.FloatFieldType,
			want:      driver.Value{Kind: driver.Float64, Float64: 3.5},
		},
		{
			name:      "a BOOLEAN",
			value:     true,
			fieldType: bigquery.BooleanFieldType,
			want:      driver.Value{Kind: driver.Bool, Bool: true},
		},
		{
			name:      "a TIMESTAMP is an instant",
			value:     instant,
			fieldType: bigquery.TimestampFieldType,
			want:      driver.Value{Kind: driver.TimestampTZ, Time: instant, Precision: bigqueryFractionalDigits},
		},
		{
			// The client states a DATETIME as a wall clock with no zone, which is
			// what the contract's Timestamp holds.
			name:      "a DATETIME is a wall clock",
			value:     civil.DateTime{Date: civil.Date{Year: 2026, Month: 1, Day: 2}, Time: civil.Time{Hour: 3, Minute: 4, Second: 5, Nanosecond: 123000000}},
			fieldType: bigquery.DateTimeFieldType,
			want: driver.Value{
				Kind:      driver.Timestamp,
				Time:      time.Date(2026, 1, 2, 3, 4, 5, 123000000, time.UTC),
				Precision: bigqueryFractionalDigits,
			},
		},
		{
			// The client hands a NUMERIC over as a *big.Rat, whose own rendering
			// is a ratio rather than BigQuery's literal.
			name:      "a NUMERIC is its exact decimal text",
			value:     big.NewRat(12345, 100),
			fieldType: bigquery.NumericFieldType,
			want:      driver.Value{Kind: driver.String, String: "123.45"},
		},
		{
			name:      "a BIGNUMERIC is its exact decimal text, without trailing zeros",
			value:     big.NewRat(50, 4),
			fieldType: bigquery.BigNumericFieldType,
			want:      driver.Value{Kind: driver.String, String: "12.5"},
		},
		{
			name:      "a whole NUMERIC states no point",
			value:     big.NewRat(12, 1),
			fieldType: bigquery.NumericFieldType,
			want:      driver.Value{Kind: driver.String, String: "12"},
		},
		{
			name:      "an ARRAY is its elements in order, each by its own rule",
			value:     []bigquery.Value{"a", int64(1)},
			fieldType: bigquery.StringFieldType,
			want:      driver.Value{Kind: driver.String, String: "[a, 1]"},
		},
		{
			name:      "an empty ARRAY states no elements",
			value:     []bigquery.Value{},
			fieldType: bigquery.StringFieldType,
			want:      driver.Value{Kind: driver.String, String: "[]"},
		},
		{
			name:      "an ARRAY of NUMERIC states each decimal, not a ratio",
			value:     []bigquery.Value{big.NewRat(3, 2), big.NewRat(9, 4)},
			fieldType: bigquery.NumericFieldType,
			want:      driver.Value{Kind: driver.String, String: "[1.5, 2.25]"},
		},
		{
			// The client hands over a record's values in the field's order and
			// the names live in the field's nested schema, which the conversion is
			// not given: a record reads as a list, which VERIFY.md pins.
			name:      "a RECORD reads as its values in order",
			value:     []bigquery.Value{int64(1), "x"},
			fieldType: bigquery.RecordFieldType,
			want:      driver.Value{Kind: driver.String, String: "[1, x]"},
		},
		{
			// Go's own rendering of the client's []bigquery.Value prints the
			// element as an address, so an array of ranges is rendered here.
			name: "an ARRAY of RANGE is each element's literal",
			value: []bigquery.Value{
				&bigquery.RangeValue{Start: civil.Date{Year: 2026, Month: 1, Day: 1}, End: civil.Date{Year: 2026, Month: 2, Day: 1}},
				&bigquery.RangeValue{},
			},
			fieldType: bigquery.RangeFieldType,
			want:      driver.Value{Kind: driver.String, String: "[[2026-01-01, 2026-02-01), [UNBOUNDED, UNBOUNDED)]"},
		},
		{
			name:      "a RANGE is the half-open literal the service states",
			value:     &bigquery.RangeValue{Start: civil.Date{Year: 2026, Month: 1, Day: 1}, End: civil.Date{Year: 2026, Month: 2, Day: 1}},
			fieldType: bigquery.RangeFieldType,
			want:      driver.Value{Kind: driver.String, String: "[2026-01-01, 2026-02-01)"},
		},
		{
			name:      "a RANGE with no bounds is UNBOUNDED on both sides",
			value:     &bigquery.RangeValue{},
			fieldType: bigquery.RangeFieldType,
			want:      driver.Value{Kind: driver.String, String: "[UNBOUNDED, UNBOUNDED)"},
		},
		{
			name:      "a DATETIME bound is written with a T, and its microseconds kept",
			value:     &bigquery.RangeValue{Start: civil.DateTime{Date: civil.Date{Year: 2026, Month: 1, Day: 2}, Time: civil.Time{Hour: 3, Minute: 4, Second: 5, Nanosecond: 123456000}}},
			fieldType: bigquery.RangeFieldType,
			want:      driver.Value{Kind: driver.String, String: "[2026-01-02T03:04:05.123456, UNBOUNDED)"},
		},
		{
			// The service's own rendering of the same range drops the fraction;
			// this keeps it, which the test pins so the choice stays deliberate.
			name:      "a TIMESTAMP bound is a UTC wall clock, and its microseconds kept",
			value:     &bigquery.RangeValue{Start: time.Date(2026, 1, 2, 3, 4, 5, 123456000, time.UTC)},
			fieldType: bigquery.RangeFieldType,
			want:      driver.Value{Kind: driver.String, String: "[2026-01-02 03:04:05.123456, UNBOUNDED)"},
		},
		{
			name:      "a value whose Go type is not the declared one falls back",
			value:     int64(7),
			fieldType: bigquery.StringFieldType,
			want:      driver.Value{Kind: driver.String, String: "7"},
		},
		{
			name:      "an unknown type falls back",
			value:     "x",
			fieldType: bigquery.FieldType("SOMETHING_NEW"),
			want:      driver.Value{Kind: driver.String, String: "x"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, convertValue(tc.value, tc.fieldType))
		})
	}
}

// TestJobReturnsRows covers the answer the command path reads a statement's rows
// from: a plain query states its statement type; DDL and DML are commands whatever
// their result schema says; and a statement the analysis could read nothing about is
// the service's to answer, by the schema it states.
func TestJobReturnsRows(t *testing.T) {
	for _, tc := range []struct {
		name            string
		statementType   string
		hasResultSchema bool
		unclassified    bool
		want            bool
	}{
		{name: "a plain SELECT, which states no schema and no rows", statementType: selectStatementType, want: true},
		{name: "a plain SELECT that returned rows", statementType: selectStatementType, hasResultSchema: true, want: true},
		{name: "a DDL, which states a schema and no rows", statementType: "CREATE_TABLE", hasResultSchema: true},
		{name: "a DML, whose rows are the count it affected", statementType: "INSERT"},
		{name: "a pipe query with rows", unclassified: true, hasResultSchema: true, want: true},
		{name: "a pipe query or a CALL that returned no rows", unclassified: true, hasResultSchema: true, want: true},
		{name: "an unclassified statement that answered no result set", unclassified: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, jobReturnsRows(tc.statementType, tc.hasResultSchema, tc.unclassified))
		})
	}
}

// TestColumnsOf covers the column list a result reports, including a nil field the
// client may hold.
func TestColumnsOf(t *testing.T) {
	columns := columnsOf(bigquery.Schema{
		{Name: "id", Type: bigquery.IntegerFieldType},
		nil,
		{Name: "label", Type: bigquery.StringFieldType},
	})
	require.Equal(t, []driver.Column{
		{Name: "id", TypeName: "INTEGER"},
		{Name: "label", TypeName: "STRING"},
	}, columns)

	require.Empty(t, columnsOf(nil))
}
