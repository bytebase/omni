package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestNamingIdentifierNoKeyword(t *testing.T) {
	keyword := func(name, text string) finding { return finding{v1review.NameIsKeyword, name, text} }
	for _, tt := range []struct {
		name string
		sql  string
		want []finding
	}{
		{
			name: "quoted table and unquoted column",
			sql:  `CREATE TABLE "table"(id int, name varchar(255))`,
			want: []finding{keyword("table", `"table"`), keyword("NAME", "name varchar(255)")},
		},
		{name: "no keyword", sql: "CREATE TABLE analyze (id int, title varchar(255))"},
		{
			name: "quoted keyword columns",
			sql:  `CREATE TABLE t_log (id NUMBER, "PARTITION" NUMBER, "Timestamp" DATE)`,
			want: []finding{keyword("PARTITION", `"PARTITION" NUMBER`), keyword("Timestamp", `"Timestamp" DATE`)},
		},
		{
			name: "schema qualifier not judged",
			sql:  `CREATE TABLE "TABLE".order_item (id NUMBER)`,
		},
		{
			name: "alias",
			sql:  `SELECT COUNT(*) AS "DATE" FROM t`,
			want: []finding{keyword("DATE", `COUNT(*) AS "DATE"`)},
		},
		{name: "references not judged", sql: "UPDATE t SET name = 1 WHERE id = COUNT(name)"},
		{name: "dropped constraint not judged", sql: `ALTER TABLE t DROP CONSTRAINT "LEVEL"`},
		{name: "admin create", sql: `CREATE DIRECTORY "LEVEL" AS '/tmp'`, want: []finding{keyword("LEVEL", `"LEVEL"`)}},
		{name: "admin drop not judged", sql: `DROP DIRECTORY "LEVEL"`},
		{name: "rename target", sql: `ALTER TABLE t RENAME COLUMN a TO "level"`, want: []finding{keyword("level", `RENAME COLUMN a TO "level"`)}},
		{
			name: "declaration in a block, on its own line",
			sql:  "DECLARE\n  name VARCHAR2(10);\nBEGIN\n  SELECT COUNT(1) INTO name FROM t;\nEND;",
			want: []finding{keyword("NAME", "name VARCHAR2(10);")},
		},
		{
			name: "built-in references in a routine not judged",
			sql:  "CREATE OR REPLACE PROCEDURE p_demo AS\n  v NUMBER;\n  TYPE t_nos IS TABLE OF NUMBER;\n  coll t_nos;\nBEGIN\n  SELECT COUNT(1) INTO v FROM t;\n  v := coll.COUNT;\n  v := SQL%ROWCOUNT;\nEND;",
		},
		{
			name: "under explain plan",
			sql:  `EXPLAIN PLAN FOR SELECT a AS "type" FROM t`,
			want: []finding{keyword("type", `a AS "type"`)},
		},
		{
			name: "statement order, then walk order",
			sql:  "CREATE TABLE t1 (id NUMBER, CONSTRAINT \"KEY_VALUE\" UNIQUE (id), name NUMBER);\nCREATE TABLE \"SIZE\" (id NUMBER);",
			want: []finding{keyword("NAME", "name NUMBER"), keyword("SIZE", `"SIZE"`)},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := anchored(t, tt.sql, NamingIdentifierNoKeyword(parse(t, tt.sql)))
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
