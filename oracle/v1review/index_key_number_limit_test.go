package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/v1review"
)

func TestIndexKeyNumberLimit(t *testing.T) {
	for _, tt := range []struct {
		name  string
		sql   string
		limit int
		want  []string
	}{
		{name: "within the limit", sql: "CREATE TABLE t (a INT, b INT, CONSTRAINT pk PRIMARY KEY (a, b))", limit: 2},
		{
			name:  "create table anchors each primary and unique constraint",
			sql:   "CREATE TABLE t (a INT, b INT, c INT,\n  CONSTRAINT pk PRIMARY KEY (a, b, c),\n  UNIQUE (a, b, c),\n  CONSTRAINT fk FOREIGN KEY (a, b, c) REFERENCES p (a, b, c),\n  CONSTRAINT chk CHECK (a > b))",
			limit: 2,
			want: []string{
				"CONSTRAINT pk PRIMARY KEY (a, b, c)",
				"UNIQUE (a, b, c)",
			},
		},
		{name: "a column's own constraint is not judged", sql: "CREATE TABLE t (a INT PRIMARY KEY, b INT UNIQUE)", limit: 0},
		{
			name:  "create index anchors the statement and counts expressions",
			sql:   "CREATE UNIQUE INDEX idx ON s.t (a, UPPER(b), c DESC)",
			limit: 2,
			want:  []string{"CREATE UNIQUE INDEX idx ON s.t (a, UPPER(b), c DESC)"},
		},
		{name: "create index within the limit", sql: "CREATE INDEX idx ON t (a, b)", limit: 2},
		{
			name:  "alter table anchors each added constraint",
			sql:   "ALTER TABLE t ADD CONSTRAINT uk UNIQUE (a, b, c)",
			limit: 2,
			want:  []string{"CONSTRAINT uk UNIQUE (a, b, c)"},
		},
		{
			name:  "alter table parenthesized add",
			sql:   "ALTER TABLE t ADD (CONSTRAINT pk PRIMARY KEY (a, b), CONSTRAINT uk UNIQUE (a))",
			limit: 1,
			want:  []string{"CONSTRAINT pk PRIMARY KEY (a, b)"},
		},
		{name: "alter table foreign key is not judged", sql: "ALTER TABLE t ADD CONSTRAINT fk FOREIGN KEY (a, b) REFERENCES p (a, b)", limit: 1},
		{name: "modify unique lists no columns", sql: "ALTER TABLE t MODIFY UNIQUE (a, b, c) DISABLE", limit: 1},
		{name: "drop unique lists no columns", sql: "ALTER TABLE t DROP UNIQUE (a, b, c)", limit: 1},
		{name: "limit 0 reports nothing", sql: "CREATE INDEX idx ON t (a, b, c)", limit: 0},
		{name: "negative limit reports nothing", sql: "CREATE INDEX idx ON t (a, b, c)", limit: -1},
		{name: "a table nested in create schema is not judged", sql: "CREATE SCHEMA AUTHORIZATION s CREATE TABLE t (a INT, b INT, PRIMARY KEY (a, b))", limit: 1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var got []string
			for _, f := range IndexKeyNumberLimit(parse(t, tt.sql), tt.limit) {
				if f.Code != v1review.IndexKeyCountExceeded || f.Name != "" || f.Table != "" || f.Value != "" {
					t.Errorf("finding %+v, want only the code and a range", f)
				}
				got = append(got, tt.sql[f.Start:f.End])
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("anchored text (-want +got):\n%s", diff)
			}
		})
	}
}

func TestIndexKeyNumberLimitBareStatement(t *testing.T) {
	sql := "CREATE INDEX idx ON t (a, b)"
	stmts := parse(t, sql)
	raw, ok := stmts[0].(*ast.RawStmt)
	if !ok {
		t.Fatalf("statement is %T, want *ast.RawStmt", stmts[0])
	}
	got := IndexKeyNumberLimit([]ast.Node{raw.Stmt}, 1)
	if len(got) != 1 || sql[got[0].Start:got[0].End] != sql {
		t.Errorf("findings %+v, want one anchoring the statement", got)
	}
}
