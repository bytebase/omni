package analysis

import (
	"errors"
	"testing"
)

func TestExplainStatement(t *testing.T) {
	for _, tc := range []struct {
		name      string
		statement string
		json      bool
		want      string
		wantErr   error
		anyErr    bool
	}{
		{
			name:      "plain select text",
			statement: "SELECT 1",
			want:      "EXPLAIN SELECT 1",
		},
		{
			name:      "plain select json",
			statement: "SELECT 1",
			json:      true,
			want:      "EXPLAIN (DISTSQL) SELECT 1",
		},
		{
			name:      "cockroach native cast json",
			statement: "SELECT 1:::INT8",
			json:      true,
			want:      "EXPLAIN (DISTSQL) SELECT 1:::INT8",
		},
		{
			name:      "unwraps existing explain",
			statement: "EXPLAIN SELECT 1",
			json:      true,
			want:      "EXPLAIN (DISTSQL) SELECT 1",
		},
		{
			name:      "unwraps existing explain text",
			statement: "EXPLAIN (VERBOSE) SELECT 1",
			want:      "EXPLAIN SELECT 1",
		},
		{
			name:      "rejects explain analyze",
			statement: "EXPLAIN ANALYZE SELECT 1",
			wantErr:   ErrExplainAnalyzeExecutes,
		},
		{
			name:      "rejects explain analyze json",
			statement: "EXPLAIN ANALYZE UPDATE t SET v = 1",
			json:      true,
			wantErr:   ErrExplainAnalyzeExecutes,
		},
		{
			name:      "rejects multiple statements",
			statement: "SELECT 1; SELECT 2",
			anyErr:    true,
		},
		{
			name:      "rejects syntax error",
			statement: "SELECT FROM",
			anyErr:    true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ExplainStatement(tc.statement, tc.json)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("ExplainStatement(%q, %v) err = %v, want %v", tc.statement, tc.json, err, tc.wantErr)
				}
				return
			}
			if tc.anyErr {
				if err == nil {
					t.Fatalf("ExplainStatement(%q, %v) succeeded with %q, want error", tc.statement, tc.json, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("ExplainStatement(%q, %v) unexpected err: %v", tc.statement, tc.json, err)
			}
			if got != tc.want {
				t.Errorf("ExplainStatement(%q, %v) = %q, want %q", tc.statement, tc.json, got, tc.want)
			}
		})
	}
}

func TestDescribeExplain(t *testing.T) {
	for _, tc := range []struct {
		statement string
		format    string
		executed  bool
		ok        bool
	}{
		{"SELECT 1", "", false, false},
		{"SELECT 1; SELECT 2", "", false, false},
		{"SELECT FROM", "", false, false},
		{"EXPLAIN SELECT 1", "TEXT", false, true},
		{"EXPLAIN (DISTSQL) SELECT 1", "TEXT", false, true},
		{"EXPLAIN (DISTSQL, JSON) SELECT 1", "JSON", false, true},
		{"EXPLAIN ANALYZE SELECT 1", "TEXT", true, true},
	} {
		t.Run(tc.statement, func(t *testing.T) {
			format, executed, ok := DescribeExplain(tc.statement)
			if format != tc.format || executed != tc.executed || ok != tc.ok {
				t.Errorf("DescribeExplain(%q) = (%q, %v, %v), want (%q, %v, %v)",
					tc.statement, format, executed, ok, tc.format, tc.executed, tc.ok)
			}
		})
	}
}
