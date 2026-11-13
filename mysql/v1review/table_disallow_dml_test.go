package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

func TestTableDisallowDML(t *testing.T) {
	disallow := []string{"t", "u", "/tmp/out"}
	rule := func(stmts []ast.Node) []v1review.Finding { return TableDisallowDML(stmts, disallow) }
	f := func(name, text string) finding { return finding{v1review.TableDisallowed, name, text} }
	for _, tt := range []struct {
		name string
		sql  string
		want []finding
	}{
		{name: "insert", sql: "INSERT INTO db.t VALUES (1);", want: []finding{f("t", "INSERT INTO db.t VALUES (1)")}},
		{name: "replace", sql: "REPLACE INTO u VALUES (1);", want: []finding{f("u", "REPLACE INTO u VALUES (1)")}},
		{name: "insert select judges only the target", sql: "INSERT INTO x SELECT * FROM t;"},
		{name: "update", sql: "UPDATE t SET a = 1;", want: []finding{f("t", "UPDATE t SET a = 1")}},
		{
			name: "multi-table update in source order",
			sql:  "UPDATE u JOIN x ON u.a = x.a JOIN t ON t.a = x.a SET u.a = 1;",
			want: []finding{f("u", "UPDATE u JOIN x ON u.a = x.a JOIN t ON t.a = x.a SET u.a = 1"), f("t", "UPDATE u JOIN x ON u.a = x.a JOIN t ON t.a = x.a SET u.a = 1")},
		},
		{name: "alias is not the table", sql: "UPDATE x AS t SET t.a = 1;"},
		{name: "where subquery not judged", sql: "UPDATE x SET a = 1 WHERE a IN (SELECT a FROM t);"},
		{
			name: "derived table and join condition subquery judged",
			sql:  "UPDATE x JOIN (SELECT a FROM t) d ON x.a = (SELECT a FROM u) SET x.a = 1;",
			want: []finding{f("t", "UPDATE x JOIN (SELECT a FROM t) d ON x.a = (SELECT a FROM u) SET x.a = 1"), f("u", "UPDATE x JOIN (SELECT a FROM t) d ON x.a = (SELECT a FROM u) SET x.a = 1")},
		},
		{name: "delete", sql: "DELETE FROM t WHERE a = 1;", want: []finding{f("t", "DELETE FROM t WHERE a = 1")}},
		{
			name: "multi-table delete judges the target list",
			sql:  "DELETE t, x FROM t JOIN x ON t.a = x.a;",
			want: []finding{f("t", "DELETE t, x FROM t JOIN x ON t.a = x.a")},
		},
		{name: "multi-table delete's other tables not judged", sql: "DELETE x FROM x JOIN t ON x.a = t.a;"},
		{name: "delete using's other tables not judged", sql: "DELETE FROM x USING x JOIN t ON x.a = t.a;"},
		{name: "outfile judged as a table name", sql: "SELECT * FROM x INTO OUTFILE '/tmp/out';", want: []finding{f("/tmp/out", "SELECT * FROM x INTO OUTFILE '/tmp/out'")}},
		{name: "select not judged", sql: "SELECT * FROM t;"},
		{name: "ddl not judged", sql: "DROP TABLE t;"},
		{name: "routine body not judged", sql: "CREATE PROCEDURE p() BEGIN DELETE FROM t; END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, tableFindings(t, tt.sql, rule)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
