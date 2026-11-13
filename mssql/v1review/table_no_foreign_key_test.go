package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/mssql/parser"
	"github.com/bytebase/omni/v1review"
)

// stmtFinding is a cross-statement finding as the tests compare it: the
// code, the name, the anchoring statement's index, and the text its range
// anchors.
type stmtFinding struct {
	Code v1review.Code
	Name string
	Stmt int
	Text string
}

func runStmtFindings(t *testing.T, sql string, rule func([]ast.Node) []v1review.Finding) []stmtFinding {
	t.Helper()
	stmts, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	var got []stmtFinding
	for _, f := range rule(stmts.Items) {
		got = append(got, stmtFinding{f.Code, f.Name, f.Statement, sql[f.Start:f.End]})
	}
	return got
}

func TestTableNoForeignKey(t *testing.T) {
	fk := func(name string, stmt int, text string) stmtFinding {
		return stmtFinding{v1review.TableForeignKey, name, stmt, text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []stmtFinding
	}{
		{name: "no foreign key", sql: "CREATE TABLE t (a INT PRIMARY KEY);"},
		{name: "column REFERENCES", sql: "CREATE TABLE t (a INT REFERENCES u (a));", want: []stmtFinding{fk("t", 0, "REFERENCES u (a)")}},
		{
			name: "first foreign key anchors, column constraints first",
			sql:  "CREATE TABLE dbo.t (a INT, CONSTRAINT f1 FOREIGN KEY (a) REFERENCES u (a), b INT REFERENCES v (b), CONSTRAINT f2 FOREIGN KEY (b) REFERENCES v (b));",
			want: []stmtFinding{fk("dbo.t", 0, "REFERENCES v (b)")},
		},
		{
			name: "alter table adds one, anchoring its constraint",
			sql:  "CREATE TABLE t (a INT); ALTER TABLE [dbo].[t] ADD CONSTRAINT f FOREIGN KEY (a) REFERENCES u (a);",
			want: []stmtFinding{fk("t", 1, "CONSTRAINT f FOREIGN KEY (a) REFERENCES u (a)")},
		},
		{
			name: "alter of an uncreated table names it as the alter wrote it",
			sql:  "ALTER TABLE [dbo].[T] ADD CONSTRAINT f FOREIGN KEY (a) REFERENCES u (a);",
			want: []stmtFinding{fk("dbo.T", 0, "CONSTRAINT f FOREIGN KEY (a) REFERENCES u (a)")},
		},
		{
			name: "a later create starts the table again",
			sql:  "ALTER TABLE t ADD CONSTRAINT f FOREIGN KEY (a) REFERENCES u (a); CREATE TABLE dbo.t (a INT);",
		},
		{
			name: "dropping does not remove one",
			sql:  "CREATE TABLE t (a INT REFERENCES u (a)); ALTER TABLE t DROP CONSTRAINT f; DROP TABLE t;",
			want: []stmtFinding{fk("t", 0, "REFERENCES u (a)")},
		},
		{name: "added column's REFERENCES not judged", sql: "ALTER TABLE t ADD b INT REFERENCES u (b);"},
		{name: "case-sensitive names are different tables", sql: "CREATE TABLE T (a INT); ALTER TABLE t ADD CONSTRAINT f FOREIGN KEY (a) REFERENCES u (a);", want: []stmtFinding{fk("t", 1, "CONSTRAINT f FOREIGN KEY (a) REFERENCES u (a)")}},
		{name: "nested not read", sql: "IF 1 = 1 CREATE TABLE t (a INT REFERENCES u (a));"},
		{
			name: "tables in the order first tracked",
			sql:  "CREATE TABLE c (a INT); CREATE TABLE b (a INT REFERENCES u (a)); ALTER TABLE a ADD CONSTRAINT f FOREIGN KEY (a) REFERENCES u (a); ALTER TABLE c ADD CONSTRAINT f FOREIGN KEY (a) REFERENCES u (a);",
			want: []stmtFinding{
				fk("c", 3, "CONSTRAINT f FOREIGN KEY (a) REFERENCES u (a)"),
				fk("b", 1, "REFERENCES u (a)"),
				fk("a", 2, "CONSTRAINT f FOREIGN KEY (a) REFERENCES u (a)"),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := runStmtFindings(t, tt.sql, TableNoForeignKey)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
