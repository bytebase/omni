package driver

import (
	"testing"
	"time"

	"github.com/bytebase/omni/driver"
)

// TestClassify pins which statements a Session runs as queries, and so which it
// runs as commands that report affected rows. The decision is
// snowflake/analysis's; this test holds the wiring to it.
func TestClassify(t *testing.T) {
	for _, tc := range []struct {
		statement string
		want      bool
	}{
		{statement: "SELECT 1", want: true},
		{statement: "SELECT * FROM t", want: true},
		{statement: "WITH c AS (SELECT 1 AS n) SELECT n FROM c", want: true},
		{statement: "SHOW TABLES", want: true},
		{statement: "DESCRIBE TABLE t", want: true},
		{statement: "EXPLAIN SELECT 1", want: true},
		{statement: "SET TIMEZONE = 'UTC'", want: false},
		{statement: "INSERT INTO t VALUES (1)", want: false},
		{statement: "UPDATE t SET a = 1", want: false},
		{statement: "DELETE FROM t", want: false},
		{statement: "MERGE INTO t USING s ON t.a = s.a WHEN MATCHED THEN UPDATE SET a = s.a", want: false},
		{statement: "CREATE TABLE t (a INT)", want: false},
		// A procedure call answers a value or a table, which the analysis
		// counts as a statement that returns nothing.
		{statement: "CALL a_procedure(1)", want: true},
		{statement: "BEGIN SELECT 1; END", want: false},
		// A statement omni cannot parse runs as a query: the answer is a hint
		// about how to send it, and the server judges the statement.
		{statement: "SELECT $1 FROM @s", want: true},
	} {
		t.Run(tc.statement, func(t *testing.T) {
			got, err := classify(tc.statement, false)
			if err != nil {
				t.Fatalf("classify(%q): %v", tc.statement, err)
			}
			if got != tc.want {
				t.Errorf("classify(%q) = %v, want %v", tc.statement, got, tc.want)
			}
		})
	}

	// A caller that masks columns traces them through the statement, so it
	// cannot run a statement the analysis could not read, and a procedure
	// call's rows are withheld.
	if _, err := classify("SELECT $1 FROM @s", true); err == nil {
		t.Error("a statement the analysis cannot read was accepted with masking on")
	}
	if got, err := classify("CALL a_procedure(1)", true); err != nil || got {
		t.Errorf("classify of a procedure call with masking on = %v, %v; want no rows and no error", got, err)
	}
}

// TestTrimStatement pins what a statement is wrapped as when the analysis cannot
// splice a LIMIT into it.
func TestTrimStatement(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{in: "SELECT 1", want: "SELECT 1"},
		{in: "SELECT 1;", want: "SELECT 1"},
		{in: " SELECT 1; ;", want: "SELECT 1"},
		{in: "\nSELECT 1\n", want: "SELECT 1"},
	} {
		if got := trimStatement(tc.in); got != tc.want {
			t.Errorf("trimStatement(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

// TestFractionalDigits pins where the client states a temporal column's
// fractional-second digits: a TIME column's arrive as its precision, and a
// TIMESTAMP column's are Snowflake's default of 9, because the client reports
// none for them.
func TestFractionalDigits(t *testing.T) {
	for _, tc := range []struct {
		meta columnMeta
		want int32
	}{
		{meta: columnMeta{typeName: "TIME", precision: 3, hasPrecision: true}, want: 3},
		{meta: columnMeta{typeName: "TIME", precision: 0, hasPrecision: true}, want: 0},
		{meta: columnMeta{typeName: "TIME"}, want: 9},
		{meta: columnMeta{typeName: "TIMESTAMP_NTZ"}, want: 9},
		{meta: columnMeta{typeName: "TIMESTAMP_LTZ"}, want: 9},
		{meta: columnMeta{typeName: "TIMESTAMP_TZ"}, want: 9},
	} {
		if got := tc.meta.fractionalDigits(); got != tc.want {
			t.Errorf("fractionalDigits(%+v) = %d, want %d", tc.meta, got, tc.want)
		}
	}
}

// TestConvertValue pins the contract value each scanned cell becomes.
func TestConvertValue(t *testing.T) {
	zone := time.FixedZone("+02:00", 2*60*60)
	for _, tc := range []struct {
		name  string
		meta  columnMeta
		value any
		want  driver.Value
	}{
		{
			name:  "NULL",
			meta:  columnMeta{typeName: "TEXT"},
			value: nil,
			want:  driver.Value{Kind: driver.Null},
		},
		{
			name:  "text",
			meta:  columnMeta{typeName: "TEXT"},
			value: "hello",
			want:  driver.Value{Kind: driver.String, String: "hello"},
		},
		{
			name:  "variant keeps its document as text",
			meta:  columnMeta{typeName: "VARIANT"},
			value: `{"a": 1}`,
			want:  driver.Value{Kind: driver.String, String: `{"a": 1}`},
		},
		{
			name:  "an object the client left as text stays text",
			meta:  columnMeta{typeName: "OBJECT"},
			value: `{"a": 1}`,
			want:  driver.Value{Kind: driver.String, String: `{"a": 1}`},
		},
		{
			name:  "an integer arrives as text and becomes an int64",
			meta:  columnMeta{typeName: "FIXED", precision: 10, hasPrecision: true},
			value: "42",
			want:  driver.Value{Kind: driver.Int64, Int64: 42},
		},
		{
			name:  "a FIXED past int64 keeps its digits",
			meta:  columnMeta{typeName: "FIXED", precision: 38},
			value: "123456789012345678901234567890",
			want:  driver.Value{Kind: driver.String, String: "123456789012345678901234567890"},
		},
		{
			name:  "a decimal keeps its exact digits",
			meta:  columnMeta{typeName: "FIXED", precision: 12, scale: 2, hasPrecision: true},
			value: "12.34",
			want:  driver.Value{Kind: driver.String, String: "12.34"},
		},
		{
			name:  "a real arrives as text and becomes a float",
			meta:  columnMeta{typeName: "REAL"},
			value: "1.5",
			want:  driver.Value{Kind: driver.Float64, Float64: 1.5},
		},
		{
			name:  "a boolean arrives as text and becomes a bool",
			meta:  columnMeta{typeName: "BOOLEAN"},
			value: "true",
			want:  driver.Value{Kind: driver.Bool, Bool: true},
		},
		{
			name:  "a value the client already typed is kept",
			meta:  columnMeta{typeName: "FIXED"},
			value: int64(7),
			want:  driver.Value{Kind: driver.Int64, Int64: 7},
		},
		{
			name:  "binary",
			meta:  columnMeta{typeName: "BINARY"},
			value: []byte{0x01, 0x02},
			want:  driver.Value{Kind: driver.Bytes, Bytes: []byte{0x01, 0x02}},
		},
		{
			name:  "an empty binary is an empty value, not NULL",
			meta:  columnMeta{typeName: "BINARY"},
			value: []byte{},
			want:  driver.Value{Kind: driver.Bytes, Bytes: []byte{}},
		},
		{
			name:  "a date is a calendar day",
			meta:  columnMeta{typeName: "DATE"},
			value: time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
			want:  driver.Value{Kind: driver.String, String: "2026-01-02"},
		},
		{
			name:  "a time keeps the digits the column declares",
			meta:  columnMeta{typeName: "TIME", precision: 3, hasPrecision: true},
			value: time.Date(0, 1, 1, 12, 34, 56, 789000000, time.UTC),
			want:  driver.Value{Kind: driver.String, String: "12:34:56.789"},
		},
		{
			name:  "a time without digits",
			meta:  columnMeta{typeName: "TIME", precision: 0, hasPrecision: true},
			value: time.Date(0, 1, 1, 12, 34, 56, 0, time.UTC),
			want:  driver.Value{Kind: driver.String, String: "12:34:56"},
		},
		{
			name:  "a timestamp without a zone is a wall clock",
			meta:  columnMeta{typeName: "TIMESTAMP_NTZ"},
			value: time.Date(2026, 1, 2, 3, 4, 5, 123456000, time.UTC),
			want:  driver.Value{Kind: driver.Timestamp, Time: time.Date(2026, 1, 2, 3, 4, 5, 123456000, time.UTC), Precision: 9},
		},
		{
			name:  "a local timestamp is an instant in the session's zone",
			meta:  columnMeta{typeName: "TIMESTAMP_LTZ"},
			value: time.Date(2026, 1, 2, 3, 4, 5, 0, zone),
			want:  driver.Value{Kind: driver.TimestampTZ, Time: time.Date(2026, 1, 2, 3, 4, 5, 0, zone), Precision: 9},
		},
		{
			name:  "a zoned timestamp keeps its own offset",
			meta:  columnMeta{typeName: "TIMESTAMP_TZ"},
			value: time.Date(2026, 1, 2, 3, 4, 5, 0, zone),
			want:  driver.Value{Kind: driver.TimestampTZ, Time: time.Date(2026, 1, 2, 3, 4, 5, 0, zone), Precision: 9},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := convertValue(tc.meta, tc.value)
			if !sameValue(got, tc.want) {
				t.Errorf("convertValue(%+v, %v) = %+v, want %+v", tc.meta, tc.value, got, tc.want)
			}
		})
	}
}

// sameValue compares two contract values, comparing their instants rather than
// their time.Time representations.
func sameValue(got, want driver.Value) bool {
	if got.Kind != want.Kind || got.Precision != want.Precision {
		return false
	}
	switch got.Kind {
	case driver.Null:
		return true
	case driver.Bool:
		return got.Bool == want.Bool
	case driver.Bytes:
		return string(got.Bytes) == string(want.Bytes)
	case driver.Float32, driver.Float64:
		return got.Float64 == want.Float64
	case driver.Int32, driver.Int64:
		return got.Int64 == want.Int64
	case driver.Uint32, driver.Uint64:
		return got.Uint64 == want.Uint64
	case driver.String:
		return got.String == want.String
	case driver.Timestamp, driver.TimestampTZ:
		return got.Time.Equal(want.Time) && zoneOffset(got.Time) == zoneOffset(want.Time)
	default:
		return false
	}
}

// zoneOffset is a time's offset from UTC in seconds.
func zoneOffset(value time.Time) int {
	_, offset := value.Zone()
	return offset
}
