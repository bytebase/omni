package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/mssql/parser"
	"github.com/bytebase/omni/v1review"
)

func TestTableDisallowDML(t *testing.T) {
	disallowed := func(name, text string) finding { return finding{v1review.TableDisallowed, name, text} }
	list := []string{"t", "dbo.t", "x"}
	for _, tt := range []struct {
		name string
		sql  string
		want []finding
	}{
		{name: "insert", sql: "INSERT INTO t VALUES (1);", want: []finding{disallowed("t", "INSERT INTO t VALUES (1)")}},
		{name: "insert without INTO", sql: "INSERT [dbo].[t] VALUES (1);", want: []finding{disallowed("dbo.t", "INSERT [dbo].[t] VALUES (1)")}},
		{name: "update", sql: "UPDATE t SET a = 1;", want: []finding{disallowed("t", "UPDATE t SET a = 1")}},
		{name: "delete", sql: "DELETE FROM t;", want: []finding{disallowed("t", "DELETE FROM t")}},
		{
			name: "merge",
			sql:  "MERGE INTO t USING u ON t.a = u.a WHEN MATCHED THEN DELETE;",
			want: []finding{disallowed("t", "MERGE INTO t USING u ON t.a = u.a WHEN MATCHED THEN DELETE")},
		},
		{name: "select into", sql: "SELECT * INTO t FROM u;", want: []finding{disallowed("t", "SELECT * INTO t FROM u")}},
		{name: "update through an alias judges the alias", sql: "UPDATE x SET a = 1 FROM t AS x;", want: []finding{disallowed("x", "UPDATE x SET a = 1 FROM t AS x")}},
		{name: "tables read not judged", sql: "INSERT INTO u SELECT * FROM t;"},
		{name: "select not judged", sql: "SELECT * FROM t;"},
		{name: "table variable not judged", sql: "INSERT INTO @t VALUES (1);"},
		{name: "case-sensitive", sql: "DELETE FROM T;"},
		{name: "DDL not judged", sql: "TRUNCATE TABLE t;"},
		{
			name: "statement with a CTE",
			sql:  "WITH c AS (SELECT 1 AS a) INSERT INTO t SELECT a FROM c;",
			want: []finding{disallowed("t", "WITH c AS (SELECT 1 AS a) INSERT INTO t SELECT a FROM c")},
		},
		{name: "nested not judged", sql: "IF 1 = 1 INSERT INTO t VALUES (1);"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			stmts, err := parser.Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			var got []finding
			for _, f := range TableDisallowDML(stmts.Items, list) {
				got = append(got, finding{f.Code, f.Name, tt.sql[f.Start:f.End]})
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
