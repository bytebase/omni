package driver

import (
	"math/big"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
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
// a statement googlesql/analysis calls a query returns rows, and every other
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
// googlesql/analysis.StatementWithResultLimit's BigQuery half.
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
			name:      "a TIMESTAMP reads as Go's rendering of the value",
			value:     instant,
			fieldType: bigquery.TimestampFieldType,
			want:      driver.Value{Kind: driver.String, String: "2026-01-02 03:04:05 +0000 UTC"},
		},
		{
			// The client hands a NUMERIC over as a *big.Rat, and Go renders it as
			// a ratio rather than as BigQuery's literal. VERIFY.md records the gap.
			name:      "a NUMERIC reads as Go's rendering of the value",
			value:     big.NewRat(12345, 100),
			fieldType: bigquery.NumericFieldType,
			want:      driver.Value{Kind: driver.String, String: "2469/20"},
		},
		{
			name:      "an ARRAY reads as Go's rendering of the slice",
			value:     []bigquery.Value{"a", int64(1)},
			fieldType: bigquery.StringFieldType,
			want:      driver.Value{Kind: driver.String, String: "[a 1]"},
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
