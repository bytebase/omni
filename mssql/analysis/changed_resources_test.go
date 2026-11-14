package analysis

import (
	"reflect"
	"testing"

	"github.com/bytebase/omni/mssql"
)

func TestExtractChangedResources(t *testing.T) {
	tests := []struct {
		statement string
		want      StatementChanges
	}{
		{
			statement: "CREATE TABLE t (a INT)",
			want:      StatementChanges{Resources: []ChangedResource{{Database: "db", Schema: "dbo", Table: "t"}}},
		},
		{
			statement: "DROP TABLE d2.s.t, u",
			want: StatementChanges{Resources: []ChangedResource{
				{Database: "d2", Schema: "s", Table: "t", AffectsData: true},
				{Database: "db", Schema: "dbo", Table: "u", AffectsData: true},
			}},
		},
		{
			statement: "INSERT INTO t VALUES (1), (2)",
			want:      StatementChanges{Resources: []ChangedResource{{Database: "db", Schema: "dbo", Table: "t"}}, InsertCount: 2},
		},
		{
			statement: "INSERT INTO t EXEC p",
			want:      StatementChanges{Resources: []ChangedResource{{Database: "db", Schema: "dbo", Table: "t"}}, DMLCount: 1},
		},
		{
			// An unqualified target binds to the FROM table it aliases.
			statement: "UPDATE x SET a = 1 FROM s.t AS x;  ",
			want: StatementChanges{
				Resources:     []ChangedResource{{Database: "db", Schema: "s", Table: "t"}},
				DMLStatements: []string{"UPDATE x SET a = 1 FROM s.t AS x;"},
				DMLCount:      1,
			},
		},
		{
			// A CTE target stays as written.
			statement: "WITH c AS (SELECT a FROM t) DELETE c FROM c",
			want: StatementChanges{
				Resources:     []ChangedResource{{Database: "db", Schema: "dbo", Table: "c"}},
				DMLStatements: []string{"WITH c AS (SELECT a FROM t) DELETE c FROM c;"},
				DMLCount:      1,
			},
		},
		{
			statement: "SELECT a INTO new_t FROM x UNION SELECT b FROM y",
			want:      StatementChanges{Resources: []ChangedResource{{Database: "db", Schema: "dbo", Table: "new_t"}}},
		},
		{statement: "SELECT a INTO #t FROM x"},
		{
			statement: "CREATE VIEW other.dbo.v AS SELECT 1 AS a",
			want:      StatementChanges{Resources: []ChangedResource{{Database: "other", DatabaseOnly: true}}},
		},
		{statement: "CREATE VIEW dbo.v AS SELECT 1 AS a"},
	}
	for _, tc := range tests {
		stmts, err := mssql.Parse(tc.statement)
		if err != nil || len(stmts) != 1 {
			t.Fatalf("%q: %v", tc.statement, err)
		}
		got := ExtractChangedResources(stmts[0].AST, "db", "dbo", tc.statement)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%q:\ngot  %+v\nwant %+v", tc.statement, got, tc.want)
		}
	}
}
