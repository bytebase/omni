package catalog

import (
	"strings"
	"testing"
)

// TestJSONValueClausesRoundTrip: the clauses change what JSON_VALUE answers, so the analyzed
// query keeps them, and what it deparses to reads back to the same text.
func TestJSONValueClausesRoundTrip(t *testing.T) {
	c := New()
	if _, err := c.Exec("CREATE DATABASE d; USE d", nil); err != nil {
		t.Fatal(err)
	}
	cases := []struct{ sql, want string }{
		{"SELECT JSON_VALUE('{}', '$.a')", "json_value('{}', '$.a')"},
		{"SELECT JSON_VALUE('{}', '$.a' RETURNING SIGNED DEFAULT 7 ON EMPTY)", "returning signed default 7 on empty)"},
		{"SELECT JSON_VALUE('{}', '$.a' DEFAULT -7 ON EMPTY DEFAULT -1.5 ON ERROR)", "default -7 on empty default -1.5 on error)"},
		{"SELECT JSON_VALUE('{}', '$.a' NULL ON EMPTY ERROR ON ERROR)", "null on empty error on error)"},
		{"SELECT JSON_VALUE('{}', '$.a' DEFAULT 'x' ON ERROR)", "default 'x' on error)"},
		{"SELECT JSON_VALUE('{}', '$.a' DEFAULT N'x' ON ERROR)", "default 'x' on error)"},
		{"SELECT JSON_VALUE('{}', '$.a' RETURNING DECIMAL(5,2))", "returning decimal(5,2))"},
		{"SELECT JSON_VALUE('{}', '$.a' RETURNING DATETIME(3))", "returning datetime(3))"},
		{"SELECT JSON_VALUE('{}', '$.a' RETURNING CHAR(3))", "returning char(3)"},
		{"SELECT JSON_VALUE('{}', '$.a' RETURNING UNSIGNED)", "returning unsigned)"},
		{"SELECT JSON_VALUE('{}', '$.a' RETURNING JSON)", "returning json)"},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			_, first, second := roundTrip(t, c, tc.sql)
			if !strings.Contains(first, tc.want) {
				t.Errorf("DeparseQuery(%q) = %q, want it to hold %q", tc.sql, first, tc.want)
			}
			if again := DeparseQuery(second); again != first {
				t.Errorf("deparsed %q, which reads back as %q", first, again)
			}
		})
	}
}

// TestJSONValueReturningIsTheResultType: RETURNING says what the call answers.
func TestJSONValueReturningIsTheResultType(t *testing.T) {
	c := New()
	if _, err := c.Exec("CREATE DATABASE d; USE d", nil); err != nil {
		t.Fatal(err)
	}
	q, _, _ := roundTrip(t, c, "SELECT JSON_VALUE('{}', '$.a' RETURNING SIGNED) AS x")
	call, ok := q.TargetList[0].Expr.(*FuncCallExprQ)
	if !ok {
		t.Fatalf("target = %T", q.TargetList[0].Expr)
	}
	if call.ResultType == nil || call.ResultType.BaseType != BaseTypeBigInt {
		t.Errorf("ResultType = %+v, want a BIGINT", call.ResultType)
	}
}
