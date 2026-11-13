package v1review

import (
	"testing"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

func TestColumnRequired(t *testing.T) {
	required := []string{"id", "created_ts"}
	miss := func(table, name string, statement int, text string) placed {
		return placed{Code: v1review.ColumnMissing, Name: name, Table: table, Statement: statement, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []placed
	}{
		{name: "all present", sql: "CREATE TABLE t (id INT, created_ts INT);"},
		{
			name: "tables and columns in order, anchored at create",
			sql:  "CREATE TABLE db.b (x INT);\nCREATE TABLE a (id INT);",
			want: []placed{
				miss("a", "created_ts", 1, "CREATE TABLE a (id INT)"),
				miss("b", "created_ts", 0, "CREATE TABLE db.b (x INT)"),
				miss("b", "id", 0, "CREATE TABLE db.b (x INT)"),
			},
		},
		{name: "alter of an unseen table has them", sql: "ALTER TABLE t ADD x INT, MODIFY id BIGINT;"},
		{
			name: "drop column of an unseen table",
			sql:  "ALTER TABLE t DROP COLUMN id;",
			want: []placed{miss("t", "id", 0, "ALTER TABLE t DROP COLUMN id")},
		},
		{
			name: "rename anchors whatever it names",
			sql:  "CREATE TABLE t (id INT);\nALTER TABLE t RENAME COLUMN x TO y;",
			want: []placed{miss("t", "created_ts", 1, "ALTER TABLE t RENAME COLUMN x TO y")},
		},
		{
			name: "change anchors only from a required column",
			sql:  "CREATE TABLE t (x INT, created_ts INT);\nALTER TABLE t CHANGE x id INT;\nALTER TABLE t CHANGE created_ts y INT;\nALTER TABLE t CHANGE y z INT;",
			want: []placed{miss("t", "created_ts", 2, "ALTER TABLE t CHANGE created_ts y INT")},
		},
		{
			name: "add restores, drop table forgets",
			sql:  "CREATE TABLE t (x INT);\nALTER TABLE t ADD id INT, ADD COLUMN (created_ts INT);\nCREATE TABLE u (x INT); DROP TABLE u;",
		},
		{
			name: "names compared bytewise",
			sql:  "CREATE TABLE t (ID INT, created_ts INT);",
			want: []placed{miss("t", "id", 0, "CREATE TABLE t (ID INT, created_ts INT)")},
		},
		{name: "routine bodies are not judged", sql: "CREATE PROCEDURE p() BEGIN CREATE TABLE t (x INT); END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			checkPlaced(t, tt.sql, func(stmts []ast.Node) []v1review.Finding {
				return ColumnRequired(stmts, required)
			}, tt.want)
		})
	}
}
