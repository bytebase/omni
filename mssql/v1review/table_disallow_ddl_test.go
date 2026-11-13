package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/mssql/parser"
	"github.com/bytebase/omni/v1review"
)

func TestTableDisallowDDL(t *testing.T) {
	disallowed := func(name, text string) finding { return finding{v1review.TableDisallowed, name, text} }
	list := []string{"t", "dbo.t", "db.t", "#tmp", "Bad Name"}
	for _, tt := range []struct {
		name string
		sql  string
		want []finding
	}{
		{name: "create", sql: "CREATE TABLE t (a INT);", want: []finding{disallowed("t", "CREATE TABLE t (a INT)")}},
		{name: "alter", sql: "ALTER TABLE dbo.t ADD b INT;", want: []finding{disallowed("dbo.t", "ALTER TABLE dbo.t ADD b INT")}},
		{name: "truncate", sql: "TRUNCATE TABLE [t];", want: []finding{disallowed("t", "TRUNCATE TABLE [t]")}},
		{name: "brackets removed", sql: "CREATE TABLE [Bad Name] (a INT);", want: []finding{disallowed("Bad Name", "CREATE TABLE [Bad Name] (a INT)")}},
		{name: "temporary table", sql: "CREATE TABLE #tmp (a INT);", want: []finding{disallowed("#tmp", "CREATE TABLE #tmp (a INT)")}},
		{name: "missing schema dropped from the name", sql: "CREATE TABLE db..t (a INT);", want: []finding{disallowed("db.t", "CREATE TABLE db..t (a INT)")}},
		{
			name: "drop judges each table in order",
			sql:  "DROP TABLE IF EXISTS u, dbo.t, t;",
			want: []finding{disallowed("dbo.t", "DROP TABLE IF EXISTS u, dbo.t, t"), disallowed("t", "DROP TABLE IF EXISTS u, dbo.t, t")},
		},
		{name: "case-sensitive", sql: "CREATE TABLE T (a INT);"},
		{name: "qualified name is not its object name", sql: "CREATE TABLE s.t (a INT);"},
		{name: "drop view not judged", sql: "DROP VIEW t;"},
		{name: "create view not judged", sql: "CREATE VIEW t AS SELECT 1 AS a;"},
		{name: "DML not judged", sql: "INSERT INTO t VALUES (1);"},
		{name: "nested not judged", sql: "IF 1 = 1 BEGIN CREATE TABLE t (a INT); END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			stmts, err := parser.Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			var got []finding
			for _, f := range TableDisallowDDL(stmts.Items, list) {
				got = append(got, finding{f.Code, f.Name, tt.sql[f.Start:f.End]})
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}

func TestTableDisallowDDLNilList(t *testing.T) {
	stmts, err := parser.Parse("CREATE TABLE t (a INT);")
	if err != nil {
		t.Fatal(err)
	}
	if got := TableDisallowDDL(stmts.Items, nil); got != nil {
		t.Errorf("nil list: got %v, want none", got)
	}
}
