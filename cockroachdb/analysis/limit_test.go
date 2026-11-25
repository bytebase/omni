package analysis

import (
	"testing"

	"github.com/bytebase/omni/enginetest/compat"
)

func TestRewriteResultLimit(t *testing.T) {
	for _, tc := range []struct {
		statement string
		want      string
	}{
		{statement: "SELECT * FROM t", want: "SELECT * FROM t LIMIT 1000"},
		{statement: "SELECT * FROM t;", want: "SELECT * FROM t LIMIT 1000;"},
		{statement: "SELECT * FROM t; -- c", want: "SELECT * FROM t LIMIT 1000; -- c"},
		{statement: "SELECT * FROM t LIMIT 5", want: "SELECT * FROM t LIMIT 5"},
		{statement: "SELECT * FROM t LIMIT 0", want: "SELECT * FROM t LIMIT 0"},
		{statement: "SELECT * FROM t LIMIT 5000 OFFSET 2", want: "SELECT * FROM t LIMIT 1000 OFFSET 2"},
		{statement: "SELECT * FROM t LIMIT ALL OFFSET 2;", want: "SELECT * FROM t LIMIT 1000 OFFSET 2;"},
		{statement: "SELECT * FROM t LIMIT NULL", want: "SELECT * FROM t LIMIT 1000"},
		{statement: "SELECT * FROM t FETCH FIRST 5000 ROWS ONLY", want: "SELECT * FROM t FETCH FIRST 1000 ROWS ONLY"},
		{statement: "SELECT * FROM t OFFSET 3", want: "SELECT * FROM t OFFSET 3 LIMIT 1000"},
		{statement: "SELECT * FROM t LIMIT 1+1", want: "SELECT * FROM t LIMIT least(1 + 1, 1000)"},
		{statement: "SELECT * FROM t LIMIT $1", want: "SELECT * FROM t LIMIT least($1, 1000)"},
		// AS OF SYSTEM TIME must stay on the top-level statement, so no rewrite wraps it.
		{statement: "SELECT * FROM t AS OF SYSTEM TIME '-10s'", want: "SELECT * FROM t AS OF SYSTEM TIME '-10s' LIMIT 1000"},
		{statement: "SELECT * FROM t AS OF SYSTEM TIME '-10s' LIMIT ALL", want: "SELECT * FROM t AS OF SYSTEM TIME '-10s' LIMIT 1000"},
		{statement: "SELECT * FROM t FOR UPDATE", want: "SELECT * FROM t LIMIT 1000 FOR UPDATE"},
		{statement: "SELECT * FROM t ORDER BY a FOR SHARE NOWAIT", want: "SELECT * FROM t ORDER BY a LIMIT 1000 FOR SHARE NOWAIT"},
		{statement: "SELECT COLLATION FOR (1)", want: "SELECT COLLATION FOR (1) LIMIT 1000"},
		{statement: "SELECT 1 AS FOR", want: "SELECT 1 AS FOR LIMIT 1000"},
		{statement: "SELECT COLLATION FOR ('a') FROM t FOR UPDATE", want: "SELECT COLLATION FOR ('a') FROM t LIMIT 1000 FOR UPDATE"},
		{statement: "SELECT * FROM t FOR NO KEY UPDATE", want: "SELECT * FROM t LIMIT 1000 FOR NO KEY UPDATE"},
		{statement: "SELECT overlay('abc' placing 'x' from 1 for 1)", want: "SELECT overlay('abc' placing 'x' from 1 for 1) LIMIT 1000"},
		{statement: "SELECT * FROM (SELECT * FROM t LIMIT 9000) s", want: "SELECT * FROM (SELECT * FROM t LIMIT 9000) s LIMIT 1000"},
		{statement: "WITH x AS (SELECT 1 LIMIT 5000) SELECT * FROM x", want: "WITH x AS (SELECT 1 LIMIT 5000) SELECT * FROM x LIMIT 1000"},
		{statement: "SELECT 1 UNION SELECT 2", want: "SELECT 1 UNION SELECT 2 LIMIT 1000"},
		{statement: "(SELECT 1)", want: "(SELECT 1) LIMIT 1000"},
		{statement: "VALUES (1), (2)", want: "VALUES (1), (2) LIMIT 1000"},
		{statement: "SELECT * FROM [SHOW TABLES]", want: "SELECT * FROM [SHOW TABLES] LIMIT 1000"},
		{statement: "SHOW RANGES FROM TABLE t", want: "SHOW RANGES FROM TABLE t"},
		{statement: "UPSERT INTO t VALUES (1)", want: "UPSERT INTO t VALUES (1)"},
		{statement: "DELETE FROM t", want: "DELETE FROM t"},
	} {
		t.Run(tc.statement, func(t *testing.T) {
			got, err := RewriteResultLimit(tc.statement, 1000)
			if err != nil {
				t.Fatalf("RewriteResultLimit() error: %v", err)
			}
			if got != tc.want {
				t.Errorf("RewriteResultLimit() = %q, want %q", got, tc.want)
			}
			if _, err := Parse(got); err != nil {
				t.Errorf("rewritten statement does not parse: %v", err)
			}
		})
	}

	for _, statement := range []string{"", "SELECT 1; SELECT 2", "SELECT * FROM"} {
		if _, err := RewriteResultLimit(statement, 1000); err == nil {
			t.Errorf("RewriteResultLimit(%q) returned no error", statement)
		}
	}
}

func TestRewriteResultLimitKeepsCorpusParsable(t *testing.T) {
	for _, row := range corpus(t) {
		if row.Expected != compat.Accept || !IsSelectStatement(row.Statement) {
			continue
		}
		got, err := RewriteResultLimit(row.Statement, 10)
		if err != nil {
			t.Errorf("RewriteResultLimit(%q): %v", row.Statement, err)
			continue
		}
		if _, err := Parse(got); err != nil {
			t.Errorf("RewriteResultLimit(%q) = %q, which does not parse: %v", row.Statement, got, err)
		}
	}
}
