package analysis

import (
	"testing"

	"github.com/bytebase/omni/databricks/parser"
)

func TestStatementWithResultLimit(t *testing.T) {
	cases := []struct {
		name        string
		statement   string
		limit       int
		want        string
		verifyParse bool
	}{
		{
			name:        "simple select",
			statement:   "SELECT * FROM t",
			limit:       100,
			want:        "SELECT * FROM t LIMIT 100",
			verifyParse: true,
		},
		{
			name:        "select with semicolon",
			statement:   "SELECT * FROM t;",
			limit:       100,
			want:        "SELECT * FROM t LIMIT 100;",
			verifyParse: true,
		},
		{
			name:        "select with trailing comment and semicolon",
			statement:   "SELECT * FROM t -- note\n;",
			limit:       100,
			want:        "SELECT * FROM t LIMIT 100 -- note\n;",
			verifyParse: true,
		},
		{
			name:        "existing stricter limit preserved",
			statement:   "SELECT * FROM t LIMIT 50;",
			limit:       100,
			want:        "SELECT * FROM t LIMIT 50;",
			verifyParse: true,
		},
		{
			name:        "existing equal limit preserved",
			statement:   "SELECT * FROM t LIMIT 100",
			limit:       100,
			want:        "SELECT * FROM t LIMIT 100",
			verifyParse: true,
		},
		{
			name:        "existing larger limit lowered",
			statement:   "SELECT * FROM t LIMIT 500;",
			limit:       100,
			want:        "SELECT * FROM t LIMIT 100;",
			verifyParse: true,
		},
		{
			name:        "limit all lowered",
			statement:   "SELECT * FROM t LIMIT ALL;",
			limit:       100,
			want:        "SELECT * FROM t LIMIT 100;",
			verifyParse: true,
		},
		{
			name:        "offset without limit inserts limit before offset",
			statement:   "SELECT * FROM t OFFSET 5;",
			limit:       100,
			want:        "SELECT * FROM t LIMIT 100 OFFSET 5;",
			verifyParse: true,
		},
		{
			name:        "inner subquery limit is not modified",
			statement:   "SELECT * FROM (SELECT * FROM t LIMIT 500) sub",
			limit:       100,
			want:        "SELECT * FROM (SELECT * FROM t LIMIT 500) sub LIMIT 100",
			verifyParse: true,
		},
		{
			name:        "cte inner limit is not modified",
			statement:   "WITH c AS (SELECT * FROM t LIMIT 500) SELECT * FROM c",
			limit:       100,
			want:        "WITH c AS (SELECT * FROM t LIMIT 500) SELECT * FROM c LIMIT 100",
			verifyParse: true,
		},
		{
			name:        "union query",
			statement:   "SELECT 1 UNION ALL SELECT 2",
			limit:       10,
			want:        "SELECT 1 UNION ALL SELECT 2 LIMIT 10",
			verifyParse: true,
		},
		{
			name:        "from-first branch limit before union gets outer cap",
			statement:   "FROM t SELECT * LIMIT 1 UNION ALL SELECT * FROM u",
			limit:       10,
			want:        "FROM t SELECT * LIMIT 1 UNION ALL SELECT * FROM u LIMIT 10",
			verifyParse: true,
		},
		{
			name:        "from-first branch limit before pipe cross join gets outer pipe limit",
			statement:   "FROM t SELECT * LIMIT 1 |> CROSS JOIN u",
			limit:       10,
			want:        "FROM t SELECT * LIMIT 1 |> CROSS JOIN u |> LIMIT 10",
			verifyParse: true,
		},
		{
			name:        "from-first pipe query",
			statement:   "FROM t |> WHERE x > 1 |> SELECT *",
			limit:       25,
			want:        "FROM t |> WHERE x > 1 |> SELECT * |> LIMIT 25",
			verifyParse: true,
		},
		{
			name:        "pipe ending in ORDER BY adds separate LIMIT stage",
			statement:   "FROM t |> ORDER BY id",
			limit:       10,
			want:        "FROM t |> ORDER BY id |> LIMIT 10",
			verifyParse: true,
		},
		{
			name:        "pipe ending in SORT BY adds separate LIMIT stage",
			statement:   "FROM t |> SORT BY id;",
			limit:       10,
			want:        "FROM t |> SORT BY id |> LIMIT 10;",
			verifyParse: true,
		},
		{
			name:        "pipe ending in CLUSTER BY adds separate LIMIT stage",
			statement:   "FROM t |> CLUSTER BY id",
			limit:       10,
			want:        "FROM t |> CLUSTER BY id |> LIMIT 10",
			verifyParse: true,
		},
		{
			name:        "pipe ending in DISTRIBUTE BY adds separate LIMIT stage",
			statement:   "FROM t |> DISTRIBUTE BY id",
			limit:       10,
			want:        "FROM t |> DISTRIBUTE BY id |> LIMIT 10",
			verifyParse: true,
		},
		{
			name:        "single-bar pipe ending in ORDER BY adds separate LIMIT stage",
			statement:   "FROM t | ORDER BY id",
			limit:       10,
			want:        "FROM t | ORDER BY id |> LIMIT 10",
			verifyParse: true,
		},
		{
			name:        "pipe followed by outer ORDER BY uses regular LIMIT",
			statement:   "FROM t |> SELECT * ORDER BY id",
			limit:       10,
			want:        "FROM t |> SELECT * ORDER BY id LIMIT 10",
			verifyParse: true,
		},
		{
			name:        "pipe SELECT with stage WINDOW uses pipe LIMIT",
			statement:   "FROM t |> SELECT id, ROW_NUMBER() OVER w WINDOW w AS (ORDER BY id)",
			limit:       10,
			want:        "FROM t |> SELECT id, ROW_NUMBER() OVER w WINDOW w AS (ORDER BY id) |> LIMIT 10",
			verifyParse: true,
		},
		{
			name:        "pipe followed by outer WINDOW uses regular LIMIT",
			statement:   "FROM t |> EXTEND x + 1 AS y WINDOW w AS (ORDER BY id)",
			limit:       10,
			want:        "FROM t |> EXTEND x + 1 AS y WINDOW w AS (ORDER BY id) LIMIT 10",
			verifyParse: true,
		},
		{
			name:        "pipe followed by outer SORT BY uses regular LIMIT",
			statement:   "FROM t |> SELECT * SORT BY id",
			limit:       10,
			want:        "FROM t |> SELECT * SORT BY id LIMIT 10",
			verifyParse: true,
		},
		{
			name:        "pipe followed by UNION ALL uses regular LIMIT",
			statement:   "FROM t |> SELECT * UNION ALL SELECT * FROM u",
			limit:       10,
			want:        "FROM t |> SELECT * UNION ALL SELECT * FROM u LIMIT 10",
			verifyParse: true,
		},
		{
			name:        "table shorthand",
			statement:   "TABLE main.app.users",
			limit:       10,
			want:        "TABLE main.app.users LIMIT 10",
			verifyParse: true,
		},
		{
			name:        "values query",
			statement:   "VALUES (1), (2), (3)",
			limit:       2,
			want:        "VALUES (1), (2), (3) LIMIT 2",
			verifyParse: true,
		},
		{
			name:        "non-integer expression limit falls back to wrap",
			statement:   "SELECT * FROM t LIMIT 10 + 20;",
			limit:       5,
			want:        "SELECT * FROM (SELECT * FROM t LIMIT 10 + 20) AS result LIMIT 5;",
			verifyParse: true,
		},
		{
			name:      "non-positive limit leaves statement unchanged",
			statement: "SELECT * FROM t",
			limit:     0,
			want:      "SELECT * FROM t",
		},
		{
			name:      "show unchanged",
			statement: "SHOW TABLES",
			limit:     10,
			want:      "SHOW TABLES",
		},
		{
			name:      "describe unchanged",
			statement: "DESCRIBE t",
			limit:     10,
			want:      "DESCRIBE t",
		},
		{
			name:      "explain unchanged",
			statement: "EXPLAIN SELECT * FROM t",
			limit:     10,
			want:      "EXPLAIN SELECT * FROM t",
		},
		{
			name:      "dml unchanged",
			statement: "INSERT INTO t VALUES (1)",
			limit:     10,
			want:      "INSERT INTO t VALUES (1)",
		},
		{
			name:      "ddl unchanged",
			statement: "CREATE TABLE t (a INT)",
			limit:     10,
			want:      "CREATE TABLE t (a INT)",
		},
		{
			name:      "compound block unchanged",
			statement: "BEGIN SELECT 1; END;",
			limit:     10,
			want:      "BEGIN SELECT 1; END;",
		},
		{
			name:      "multi-statement script unchanged",
			statement: "SELECT 1; SELECT 2;",
			limit:     10,
			want:      "SELECT 1; SELECT 2;",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := StatementWithResultLimit(tc.statement, tc.limit)
			if got != tc.want {
				t.Errorf("StatementWithResultLimit(%q, %d) =\n  %q\nwant:\n  %q", tc.statement, tc.limit, got, tc.want)
			}
			if tc.verifyParse {
				if _, err := parser.Parse(got); err != nil {
					t.Errorf("rewritten query %q failed to parse: %v", got, err)
				}
			}
		})
	}
}

func TestWrapResultLimit(t *testing.T) {
	got := WrapResultLimit("  SELECT 1 FROM t; -- trailing\n/* nested /* c */ */ ", 50)
	want := "SELECT * FROM (SELECT 1 FROM t) AS result LIMIT 50;"
	if got != want {
		t.Errorf("WrapResultLimit = %q, want %q", got, want)
	}
	if _, err := parser.Parse(got); err != nil {
		t.Errorf("WrapResultLimit output %q failed to parse: %v", got, err)
	}
}
