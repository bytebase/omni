package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestNamingIdentifierCase(t *testing.T) {
	mismatch := func(name, text string) anchored {
		return anchored{Code: v1review.NameCaseMismatch, Name: name, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []anchored
	}{
		{name: "an unquoted name is never reported", sql: "CREATE TABLE t (id INT, Name VARCHAR, AGE INT);"},
		{name: "a quoted upper-case name", sql: `CREATE TABLE t ("ID" INT, "A_1" INT, "" INT, "数据" INT);`},
		{
			name: "a quoted name with a lower-case letter, anchored at the first column",
			sql:  "CREATE TABLE t(\n  name VARCHAR,\n  \"address\" VARCHAR,\n  \"Zip\" INT\n);",
			want: []anchored{mismatch("address", "name VARCHAR"), mismatch("Zip", "name VARCHAR")},
		},
		{
			name: "a doubled quote is read as one",
			sql:  `CREATE TABLE t ("a""b" INT);`,
			want: []anchored{mismatch(`a"b`, `"a""b" INT`)},
		},
		{
			name: "create table as select",
			sql:  "CREATE TABLE c (\n  name INT,\n  \"address\" INT\n) AS SELECT a, b FROM t;",
			want: []anchored{mismatch("address", "name INT")},
		},
		{
			name: "the anchor is a constraint that comes first",
			sql:  "CREATE TABLE t (\n  PRIMARY KEY (id),\n  id INT,\n  \"address\" INT\n);",
			want: []anchored{mismatch("address", "PRIMARY KEY (id)")},
		},
		{
			name: "the anchor is an inline index that comes first",
			sql:  "CREATE HYBRID TABLE t (\n  INDEX idx (id),\n  id INT PRIMARY KEY,\n  \"address\" INT\n);",
			want: []anchored{mismatch("address", "INDEX idx (id)")},
		},
		{
			name: "rename column judges only the new name, anchored at it",
			sql:  "ALTER TABLE t RENAME COLUMN \"old\" TO b;\nALTER TABLE t RENAME COLUMN a\n  TO \"address\";",
			want: []anchored{mismatch("address", `"address"`)},
		},
		{
			name: "each rename column of one statement is judged",
			sql:  "ALTER TABLE t RENAME COLUMN a TO \"b\",\n  RENAME TO \"tbl\",\n  RENAME COLUMN c TO D,\n  RENAME COLUMN e TO \"f\";",
			want: []anchored{mismatch("b", `"b"`), mismatch("f", `"f"`)},
		},
		{
			name: "table, constraint and index names are not judged",
			sql:  "CREATE TABLE \"tbl\" (id INT, CONSTRAINT \"pk\" PRIMARY KEY (id));\nCREATE HYBRID TABLE t (id INT PRIMARY KEY, INDEX \"idx\" (id));",
		},
		{
			name: "other alter table actions are not judged",
			sql:  "ALTER TABLE t ADD COLUMN \"address\" INT;\nALTER TABLE t DROP COLUMN \"address\";\nALTER TABLE t RENAME TO \"tbl\";\nALTER TABLE t RENAME CONSTRAINT a TO \"pk\";",
		},
		{
			name: "other statements are not judged",
			sql:  "CREATE VIEW v (\"address\") AS SELECT 1;\nSELECT 1 AS \"address\";",
		},
		{name: "a statement under EXPLAIN is not judged", sql: `EXPLAIN CREATE TABLE t ("address" INT);`},
		{name: "a scripting block is not judged", sql: "BEGIN\n  CREATE TABLE t (\"address\" INT);\nEND"},
		{
			name: "a task body is not judged",
			sql:  `CREATE TASK tk WAREHOUSE = wh SCHEDULE = '1 MINUTE' AS CREATE TABLE t ("address" INT);`,
		},
		{
			name: "findings in statement order",
			sql:  "CREATE TABLE a (\"x\" INT);\nALTER TABLE a RENAME COLUMN x TO \"y\";\nCREATE TABLE b (\"z\" INT, \"w\" INT);",
			want: []anchored{mismatch("x", `"x" INT`), mismatch("y", `"y"`), mismatch("z", `"z" INT`), mismatch("w", `"z" INT`)},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, cfg := range []v1review.CaseConfig{{Upper: true}, {Upper: false}} {
				got := anchor(tt.sql, NamingIdentifierCase(parse(t, tt.sql), cfg))
				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Errorf("Upper=%v: findings mismatch (-want +got):\n%s", cfg.Upper, diff)
				}
			}
		})
	}
}
