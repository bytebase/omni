package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// systemFinding is a finding as the system rule tests compare it: its code and
// name, the statement it is in, and the text its range anchors.
type systemFinding struct {
	Code      v1review.Code
	Name      string
	Statement int
	Text      string
}

// judgeSystem parses sql and returns what rule finds in its statements.
func judgeSystem(t *testing.T, sql string, rule func([]ast.Node) []v1review.Finding) []systemFinding {
	t.Helper()
	var got []systemFinding
	for _, f := range rule(parse(t, sql)) {
		got = append(got, systemFinding{f.Code, f.Name, f.Statement, sql[f.Start:f.End]})
	}
	return got
}

func TestSystemCharsetAllowlist(t *testing.T) {
	bad := func(name, text string) systemFinding {
		return systemFinding{Code: v1review.CharsetNotAllowed, Name: name, Text: text}
	}
	const alter = "ALTER TABLE t ADD COLUMN (a TEXT CHARSET ascii, b TEXT CHARSET big5), MODIFY c TEXT CHARACTER SET gbk, CHARSET = latin1, CHANGE d e TEXT CHARSET utf8mb4"
	for _, tt := range []struct {
		name string
		sql  string
		want []systemFinding
	}{
		{name: "allowed in any case", sql: "CREATE TABLE t (a VARCHAR(9) CHARACTER SET UTF8MB4) CHARSET = Utf8mb4;"},
		{name: "no charset", sql: "CREATE TABLE t (a INT); ALTER TABLE t ADD COLUMN b INT;"},
		{
			name: "create database",
			sql:  "CREATE DATABASE d DEFAULT CHARACTER SET Latin1;",
			want: []systemFinding{bad("latin1", "CREATE DATABASE d DEFAULT CHARACTER SET Latin1")},
		},
		{
			name: "create database judges its first charset only",
			sql:  "CREATE DATABASE d CHARSET utf8mb4 CHARSET latin1;",
		},
		{
			name: "alter database judges its first charset only",
			sql:  "ALTER DATABASE d CHARSET latin1 CHARSET ascii;",
			want: []systemFinding{bad("latin1", "ALTER DATABASE d CHARSET latin1 CHARSET ascii")},
		},
		{
			name: "create table stops at a bad table charset",
			sql:  "CREATE TABLE t (a TEXT CHARSET ascii) CHARSET = big5 CHARACTER SET gbk;",
			want: []systemFinding{bad("big5", "CREATE TABLE t (a TEXT CHARSET ascii) CHARSET = big5 CHARACTER SET gbk")},
		},
		{
			name: "create table judges its first table charset only",
			sql:  "CREATE TABLE t (a INT) CHARSET = utf8mb4 CHARACTER SET gbk;",
		},
		{
			name: "create table reports its first bad column, anchoring it",
			sql:  "CREATE TABLE t (a TEXT CHARSET utf8mb4,\n  b TEXT CHARACTER SET LATIN1,\n  c TEXT CHARSET ascii) CHARSET = utf8mb4;",
			want: []systemFinding{bad("latin1", "b TEXT CHARACTER SET LATIN1")},
		},
		{
			name: "alter table reports its last bad charset once",
			sql:  alter + ";",
			want: []systemFinding{bad("latin1", alter)},
		},
		{
			name: "alter table add column reports its first bad column",
			sql:  "ALTER TABLE t ADD COLUMN (a TEXT CHARSET ascii, b TEXT CHARSET big5);",
			want: []systemFinding{bad("ascii", "ALTER TABLE t ADD COLUMN (a TEXT CHARSET ascii, b TEXT CHARSET big5)")},
		},
		{
			name: "alter table convert",
			sql:  "ALTER TABLE t CONVERT TO CHARACTER SET gbk COLLATE gbk_bin;",
			want: []systemFinding{bad("gbk", "ALTER TABLE t CONVERT TO CHARACTER SET gbk COLLATE gbk_bin")},
		},
		{
			name: "alter table convert to default is judged as a name",
			sql:  "ALTER TABLE t CONVERT TO CHARSET DEFAULT;",
			want: []systemFinding{bad("default", "ALTER TABLE t CONVERT TO CHARSET DEFAULT")},
		},
		{name: "collate options are not judged", sql: "ALTER TABLE t COLLATE = latin1_bin; CREATE DATABASE d COLLATE latin1_bin;"},
		{name: "a routine body is not judged", sql: "CREATE PROCEDURE p() BEGIN CREATE TABLE t (a TEXT CHARSET ascii); END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := judgeSystem(t, tt.sql, func(stmts []ast.Node) []v1review.Finding {
				return SystemCharsetAllowlist(stmts, []string{"UTF8MB4"})
			})
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestSystemCharsetAllowlistEmptyList(t *testing.T) {
	got := judgeSystem(t, "CREATE TABLE t (a INT) CHARSET = utf8mb4;", func(stmts []ast.Node) []v1review.Finding {
		return SystemCharsetAllowlist(stmts, nil)
	})
	want := []systemFinding{{Code: v1review.CharsetNotAllowed, Name: "utf8mb4", Text: "CREATE TABLE t (a INT) CHARSET = utf8mb4"}}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("findings mismatch (-want +got):\n%s", diff)
	}
}
