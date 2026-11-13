package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/snowflake/ast"
	"github.com/bytebase/omni/snowflake/parser"
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
	file, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	var got []stmtFinding
	for _, f := range rule(file.Stmts) {
		got = append(got, stmtFinding{f.Code, f.Name, f.Statement, sql[f.Start:f.End]})
	}
	return got
}

func TestTableNoForeignKey(t *testing.T) {
	fk := func(name string, stmt int, text string) stmtFinding {
		return stmtFinding{v1review.TableForeignKey, name, stmt, text}
	}
	const create = "CREATE TABLE t (a INT FOREIGN KEY REFERENCES u (a))"
	for _, tt := range []struct {
		name string
		sql  string
		want []stmtFinding
	}{
		{name: "no foreign key", sql: "CREATE TABLE t (a INT PRIMARY KEY);"},
		{name: "inline foreign key anchors the statement", sql: create + ";", want: []stmtFinding{fk("t", 0, create)}},
		{
			name: "table constraint anchors, the last one read",
			sql:  "CREATE TABLE t (a INT, CONSTRAINT f1 FOREIGN KEY (a) REFERENCES u (a), FOREIGN KEY (a) REFERENCES v (a));",
			want: []stmtFinding{fk("t", 0, "FOREIGN KEY (a) REFERENCES v (a)")},
		},
		{name: "as select, like and clone not read", sql: "CREATE TABLE t LIKE u; CREATE TABLE t CLONE u; CREATE TABLE t AS SELECT 1 AS a;"},
		{
			name: "alter table adds one",
			sql:  "CREATE TABLE t (a INT); ALTER TABLE public.T ADD CONSTRAINT f FOREIGN KEY (a) REFERENCES u (a);",
			want: []stmtFinding{fk("public.T", 1, "CONSTRAINT f FOREIGN KEY (a) REFERENCES u (a)")},
		},
		{
			name: "drop foreign key counts one fewer",
			sql:  "CREATE TABLE t (a INT FOREIGN KEY REFERENCES u (a), b INT FOREIGN KEY REFERENCES u (b)); ALTER TABLE t DROP FOREIGN KEY (a);",
			want: []stmtFinding{fk("t", 0, "CREATE TABLE t (a INT FOREIGN KEY REFERENCES u (a), b INT FOREIGN KEY REFERENCES u (b))")},
		},
		{name: "drop foreign key to zero", sql: create + "; ALTER TABLE t DROP FOREIGN KEY (a);"},
		{
			name: "drop constraint by name keeps the count",
			sql:  create + "; ALTER TABLE t DROP CONSTRAINT f;",
			want: []stmtFinding{fk("t", 0, create)},
		},
		{
			name: "a constraint action renames, a column action does not",
			sql:  create + "; ALTER TABLE \"T\" ADD COLUMN b INT; ALTER TABLE PUBLIC.T RENAME CONSTRAINT x TO y;",
			want: []stmtFinding{fk("PUBLIC.T", 0, create)},
		},
		{name: "drop table does not forget", sql: create + "; DROP TABLE t;", want: []stmtFinding{fk("t", 0, create)}},
		{name: "a later create starts the table again", sql: "ALTER TABLE t ADD FOREIGN KEY (a) REFERENCES u (a); CREATE TABLE T (a INT);"},
		{name: "quoted lower case is another table", sql: create + "; ALTER TABLE \"t\" DROP FOREIGN KEY (a);", want: []stmtFinding{fk("t", 0, create)}},
		{
			name: "tables in the order first counted",
			sql:  "CREATE TABLE c (a INT); CREATE TABLE b (a INT FOREIGN KEY REFERENCES u (a)); ALTER TABLE a ADD FOREIGN KEY (a) REFERENCES u (a); ALTER TABLE c ADD FOREIGN KEY (a) REFERENCES u (a);",
			want: []stmtFinding{
				fk("c", 3, "FOREIGN KEY (a) REFERENCES u (a)"),
				fk("b", 1, "CREATE TABLE b (a INT FOREIGN KEY REFERENCES u (a))"),
				fk("a", 2, "FOREIGN KEY (a) REFERENCES u (a)"),
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
