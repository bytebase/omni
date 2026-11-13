package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/mariadb/parser"
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
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	var got []systemFinding
	for _, f := range rule(list.Items) {
		got = append(got, systemFinding{f.Code, f.Name, f.Statement, sql[f.Start:f.End]})
	}
	return got
}

func TestSystemCharsetAllowlist(t *testing.T) {
	bad := func(name, text string) systemFinding {
		return systemFinding{Code: v1review.CharsetNotAllowed, Name: name, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []systemFinding
	}{
		{name: "allowed in any case", sql: "CREATE TABLE t (a VARCHAR(9) CHARACTER SET UTF8MB4) CHARSET = Utf8mb4;"},
		{name: "no charset", sql: "CREATE TABLE t (a INT);"},
		{
			name: "create database",
			sql:  "CREATE DATABASE d CHARACTER SET Latin1;",
			want: []systemFinding{bad("latin1", "CREATE DATABASE d CHARACTER SET Latin1")},
		},
		{
			name: "create database judges a repeated charset once",
			sql:  "CREATE DATABASE d CHARSET latin1 CHARSET latin1;",
			want: []systemFinding{bad("latin1", "CREATE DATABASE d CHARSET latin1 CHARSET latin1")},
		},
		{
			name: "alter database judges a repeated charset each time",
			sql:  "ALTER DATABASE d CHARSET latin1 CHARSET latin1;",
			want: []systemFinding{
				bad("latin1", "ALTER DATABASE d CHARSET latin1 CHARSET latin1"),
				bad("latin1", "ALTER DATABASE d CHARSET latin1 CHARSET latin1"),
			},
		},
		{
			name: "create table judges its first table charset once, then each column",
			sql:  "CREATE TABLE t (a VARCHAR(9) CHARACTER SET LATIN1, b TEXT CHARSET ascii, c INT) DEFAULT CHARSET = big5 CHARSET = big5;",
			want: []systemFinding{
				bad("big5", "CREATE TABLE t (a VARCHAR(9) CHARACTER SET LATIN1, b TEXT CHARSET ascii, c INT) DEFAULT CHARSET = big5 CHARSET = big5"),
				bad("latin1", "CREATE TABLE t (a VARCHAR(9) CHARACTER SET LATIN1, b TEXT CHARSET ascii, c INT) DEFAULT CHARSET = big5 CHARSET = big5"),
				bad("ascii", "CREATE TABLE t (a VARCHAR(9) CHARACTER SET LATIN1, b TEXT CHARSET ascii, c INT) DEFAULT CHARSET = big5 CHARSET = big5"),
			},
		},
		{
			name: "alter table columns and options",
			sql:  "ALTER TABLE t ADD COLUMN a TEXT CHARSET ascii, MODIFY b TEXT CHARACTER SET big5, CHANGE c d TEXT CHARSET gbk, DEFAULT CHARSET = latin1;",
			want: []systemFinding{
				bad("ascii", "ALTER TABLE t ADD COLUMN a TEXT CHARSET ascii, MODIFY b TEXT CHARACTER SET big5, CHANGE c d TEXT CHARSET gbk, DEFAULT CHARSET = latin1"),
				bad("big5", "ALTER TABLE t ADD COLUMN a TEXT CHARSET ascii, MODIFY b TEXT CHARACTER SET big5, CHANGE c d TEXT CHARSET gbk, DEFAULT CHARSET = latin1"),
				bad("gbk", "ALTER TABLE t ADD COLUMN a TEXT CHARSET ascii, MODIFY b TEXT CHARACTER SET big5, CHANGE c d TEXT CHARSET gbk, DEFAULT CHARSET = latin1"),
				bad("latin1", "ALTER TABLE t ADD COLUMN a TEXT CHARSET ascii, MODIFY b TEXT CHARACTER SET big5, CHANGE c d TEXT CHARSET gbk, DEFAULT CHARSET = latin1"),
			},
		},
		{
			name: "parenthesized add",
			sql:  "ALTER TABLE t ADD COLUMN (a TEXT CHARSET ascii, b INT);",
			want: []systemFinding{bad("ascii", "ALTER TABLE t ADD COLUMN (a TEXT CHARSET ascii, b INT)")},
		},
		{name: "convert to is not judged", sql: "ALTER TABLE t CONVERT TO CHARACTER SET latin2;"},
		{name: "a routine body is not judged", sql: "CREATE PROCEDURE p() BEGIN CREATE TABLE t (a TEXT CHARSET ascii); END;"},
		{name: "other statements are not judged", sql: "SET NAMES latin1; SELECT CONVERT('a' USING latin1);"},
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
	if got := judgeSystem(t, "CREATE DATABASE d CHARSET latin1;", func(stmts []ast.Node) []v1review.Finding {
		return SystemCharsetAllowlist(stmts, nil)
	}); len(got) != 1 {
		t.Errorf("empty allow-list: got %v, want one finding", got)
	}
}
