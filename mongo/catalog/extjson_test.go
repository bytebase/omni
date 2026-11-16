package catalog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRenderExtendedJSON(t *testing.T) {
	for _, tc := range []struct {
		name string
		in   string
		want string
	}{
		{name: "null", in: `null`, want: `null`},
		{name: "bool", in: `true`, want: `true`},
		{name: "string", in: `"a\"b"`, want: `"a\"b"`},
		{name: "int32", in: `{"$numberInt":"18"}`, want: `18`},
		{name: "int64 keeps its width", in: `{"$numberLong":"18"}`, want: `NumberLong("18")`},
		{name: "double", in: `{"$numberDouble":"1.5"}`, want: `1.5`},
		{name: "integral double keeps its kind", in: `{"$numberDouble":"2"}`, want: `2.0`},
		{name: "decimal", in: `{"$numberDecimal":"1.25"}`, want: `NumberDecimal("1.25")`},
		{name: "empty object", in: `{}`, want: `{}`},
		{name: "object keeps the document's order", in: `{"b":1,"a":2}`, want: `{"b": 1, "a": 2}`},
		{name: "array", in: `[{"$numberInt":"1"},"x",null]`, want: `[1, "x", null]`},
		{name: "aggregation operators are document fields", in: `{"$match":{"age":{"$gte":{"$numberInt":"18"}}}}`, want: `{"$match": {"age": {"$gte": 18}}}`},
		{name: "object id", in: `{"$oid":"0123456789abcdef01234567"}`, want: `ObjectId("0123456789abcdef01234567")`},
		{name: "date", in: `{"$date":{"$numberLong":"1767323045000"}}`, want: `ISODate("2026-01-02T03:04:05.000Z")`},
		{name: "timestamp", in: `{"$timestamp":{"t":123,"i":{"$numberInt":"4"}}}`, want: `Timestamp(123, 4)`},
		{name: "regular expression", in: `{"$regularExpression":{"pattern":"^a/b","options":"i"}}`, want: `/^a\/b/i`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := renderExtendedJSON(tc.in)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestRenderExtendedJSONRefusesAKindGomongoCannotBuild(t *testing.T) {
	for _, tc := range []struct {
		name string
		in   string
		want string
	}{
		{name: "binary", in: `{"$binary":{"base64":"AQI=","subType":"00"}}`, want: "no mongosh spelling"},
		{name: "min key", in: `{"$minKey":1}`, want: "no mongosh spelling"},
		{name: "max key", in: `{"$maxKey":1}`, want: "no mongosh spelling"},
		{name: "symbol", in: `{"$symbol":"x"}`, want: "no mongosh spelling"},
		{name: "code", in: `{"$code":"x"}`, want: "no mongosh spelling"},
		{name: "nan", in: `{"$numberDouble":"NaN"}`, want: "no mongosh literal"},
		{name: "infinity", in: `{"$numberDouble":"Infinity"}`, want: "no mongosh literal"},
		{name: "trailing data", in: `{} {}`, want: "trailing data"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := renderExtendedJSON(tc.in)
			require.ErrorContains(t, err, tc.want)
		})
	}
}
