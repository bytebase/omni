package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/v1review"
)

func TestSystemCollationAllowlist(t *testing.T) {
	bad := func(name, text string) systemFinding {
		return systemFinding{Code: v1review.CollationNotAllowed, Name: name, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []systemFinding
	}{
		{name: "allowed in any case", sql: "CREATE TABLE t (a VARCHAR(9) COLLATE UTF8MB4_BIN) COLLATE = Utf8mb4_bin;"},
		{name: "no collation", sql: "CREATE TABLE t (a INT);"},
		{
			name: "create database judges every collation at the statement",
			sql:  "CREATE DATABASE d COLLATE latin1_bin COLLATE latin1_bin;",
			want: []systemFinding{
				bad("latin1_bin", "CREATE DATABASE d COLLATE latin1_bin COLLATE latin1_bin"),
				bad("latin1_bin", "CREATE DATABASE d COLLATE latin1_bin COLLATE latin1_bin"),
			},
		},
		{
			name: "alter database",
			sql:  "ALTER DATABASE d COLLATE Latin1_Bin;",
			want: []systemFinding{bad("latin1_bin", "ALTER DATABASE d COLLATE Latin1_Bin")},
		},
		{
			name: "create table: options at the option, columns at the column",
			sql:  "CREATE TABLE t (a VARCHAR(9) COLLATE latin1_bin, b TEXT NOT NULL COLLATE ascii_bin, c INT) COLLATE = big5_bin;",
			want: []systemFinding{
				bad("big5_bin", "COLLATE = big5_bin"),
				bad("latin1_bin", "a VARCHAR(9) COLLATE latin1_bin"),
				bad("ascii_bin", "b TEXT NOT NULL COLLATE ascii_bin"),
			},
		},
		{
			name: "alter table: columns at the action, options at the option",
			sql:  "ALTER TABLE t ADD COLUMN a TEXT COLLATE ascii_bin, MODIFY b TEXT COLLATE big5_bin, CHANGE c d TEXT COLLATE gbk_bin, COLLATE = latin1_bin;",
			want: []systemFinding{
				bad("ascii_bin", "ADD COLUMN a TEXT COLLATE ascii_bin"),
				bad("big5_bin", "MODIFY b TEXT COLLATE big5_bin"),
				bad("gbk_bin", "CHANGE c d TEXT COLLATE gbk_bin"),
				bad("latin1_bin", "COLLATE = latin1_bin"),
			},
		},
		{name: "convert to is not judged", sql: "ALTER TABLE t CONVERT TO CHARACTER SET latin1 COLLATE latin1_bin;"},
		{name: "a routine body is not judged", sql: "CREATE PROCEDURE p() BEGIN CREATE TABLE t (a TEXT COLLATE ascii_bin); END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := judgeSystem(t, tt.sql, func(stmts []ast.Node) []v1review.Finding {
				return SystemCollationAllowlist(stmts, []string{"UTF8MB4_bin"})
			})
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
