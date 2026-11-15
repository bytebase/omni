package driver

import (
	"strings"
	"testing"
	"time"
)

func TestByteSize(t *testing.T) {
	cols := []Column{
		{Name: "id", TypeName: "INT4"},
		{Name: "payload", TypeName: "TEXT"},
	}
	if got := HeaderByteSize(cols); got <= 0 {
		t.Fatalf("HeaderByteSize = %d, want > 0", got)
	}

	zone := time.FixedZone("CST", 8*3600)
	cases := []struct {
		name string
		val  Value
	}{
		{"null", Value{Kind: Null}},
		{"bool false", Value{Kind: Bool, Bool: false}},
		{"bool true", Value{Kind: Bool, Bool: true}},
		{"int32", Value{Kind: Int32, Int64: -1}},
		{"int64", Value{Kind: Int64, Int64: 42}},
		{"uint32", Value{Kind: Uint32, Uint64: 7}},
		{"uint64", Value{Kind: Uint64, Uint64: 1 << 40}},
		{"float32", Value{Kind: Float32, Float64: 1.5}},
		{"float64", Value{Kind: Float64, Float64: 3.14159}},
		{"string", Value{Kind: String, String: "hello"}},
		{"bytes", Value{Kind: Bytes, Bytes: []byte{1, 2, 3, 4}}},
		{"timestamp", Value{Kind: Timestamp, Time: time.Date(2026, 3, 4, 5, 6, 7, 891000, time.UTC), Precision: 6}},
		{"timestamptz", Value{Kind: TimestampTZ, Time: time.Date(2026, 3, 4, 13, 6, 7, 891000, zone), Precision: 3}},
		{"structured", Value{Kind: Structured, Structured: map[string]any{"a": "x", "b": []any{nil, true, 1.5}}}},
	}
	for _, tc := range cases {
		if got := tc.val.ByteSize(); got <= 0 {
			t.Errorf("%s: ByteSize() = %d, want > 0", tc.name, got)
		}
	}

	big := strings.Repeat("x", 4096)
	row := []Value{{Kind: Int64, Int64: 1}, {Kind: String, String: big}}
	if got := RowByteSize(row); got < int64(len(big)) || got > int64(len(big))+32 {
		t.Errorf("RowByteSize(4KB) = %d, want within framing overhead of %d", got, len(big))
	}
}
