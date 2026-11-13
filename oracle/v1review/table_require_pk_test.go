package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/v1review"
)

func TestTableRequirePK(t *testing.T) {
	noPK := func(name string, statement int, text string) tableFinding {
		return tableFinding{v1review.TablePrimaryKeyMissing, name, statement, text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []tableFinding
	}{
		{name: "column primary key", sql: "CREATE TABLE t (a INT PRIMARY KEY);"},
		{name: "table primary key", sql: "CREATE TABLE t (a INT, CONSTRAINT pk PRIMARY KEY (a));"},
		{
			name: "no primary key",
			sql:  "CREATE TABLE \"x\".t (a INT);",
			want: []tableFinding{noPK("x.T", 0, "CREATE TABLE \"x\".t (a INT)")},
		},
		{name: "added later", sql: "CREATE TABLE t (a INT);\nALTER TABLE t ADD CONSTRAINT pk PRIMARY KEY (a);"},
		{name: "dropped table forgotten", sql: "CREATE TABLE t (a INT);\nDROP TABLE t;"},
		{name: "a drop of any kind forgets", sql: "CREATE TABLE t (a INT);\nDROP VIEW t;"},
		{
			name: "drop primary key anchors the command",
			sql:  "CREATE TABLE t (a INT PRIMARY KEY);\nALTER TABLE t DROP PRIMARY KEY;",
			want: []tableFinding{noPK("S.T", 1, "DROP PRIMARY KEY")},
		},
		{
			name: "drop primary key of a table never created",
			sql:  "ALTER TABLE t DROP PRIMARY KEY;",
			want: []tableFinding{noPK("S.T", 0, "DROP PRIMARY KEY")},
		},
		{name: "drop another constraint keeps the key", sql: "CREATE TABLE t (a INT PRIMARY KEY);\nALTER TABLE t DROP CONSTRAINT c;"},
		{name: "drop constraint named primary keeps the key", sql: "CREATE TABLE t (a INT PRIMARY KEY);\nALTER TABLE t DROP CONSTRAINT \"primary\";\nALTER TABLE t DROP CONSTRAINT primary;"},
		{
			name: "ordered by first record, anchored on the last create",
			sql:  "CREATE TABLE b (a INT);\nCREATE TABLE a (a INT);\nDROP TABLE b;\nCREATE TABLE b (x INT);",
			want: []tableFinding{noPK("S.B", 3, "CREATE TABLE b (x INT)"), noPK("S.A", 1, "CREATE TABLE a (a INT)")},
		},
		{
			name: "the anchor survives a drop and an added key",
			sql:  "CREATE TABLE t (a INT);\nDROP TABLE t;\nALTER TABLE t ADD CONSTRAINT pk PRIMARY KEY (a);\nALTER TABLE t ADD c INT;",
		},
		{name: "plsql block not judged", sql: "BEGIN EXECUTE IMMEDIATE 'CREATE TABLE t (a INT)'; END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := runTableRule(t, tt.sql, func(stmts []ast.Node) []v1review.Finding { return TableRequirePK(stmts, "S") })
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
