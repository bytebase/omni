package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

func TestSystemCollationAllowlist(t *testing.T) {
	bad := func(name, text string) systemFinding {
		return systemFinding{Code: v1review.CollationNotAllowed, Name: name, Text: text}
	}
	const alter = "ALTER TABLE t ADD COLUMN (a TEXT COLLATE ascii_bin, b TEXT COLLATE big5_bin), MODIFY c TEXT COLLATE gbk_bin, COLLATE = latin1_bin, CHANGE d e TEXT COLLATE utf8mb4_bin"
	for _, tt := range []struct {
		name string
		sql  string
		want []systemFinding
	}{
		{name: "allowed in any case", sql: "CREATE TABLE t (a VARCHAR(9) COLLATE UTF8MB4_BIN) COLLATE = Utf8mb4_bin;"},
		{
			name: "create database judges its first collation only",
			sql:  "CREATE DATABASE d DEFAULT COLLATE Latin1_Bin COLLATE ascii_bin;",
			want: []systemFinding{bad("latin1_bin", "CREATE DATABASE d DEFAULT COLLATE Latin1_Bin COLLATE ascii_bin")},
		},
		{
			name: "alter database",
			sql:  "ALTER DATABASE d CHARSET latin1 COLLATE latin1_bin;",
			want: []systemFinding{bad("latin1_bin", "ALTER DATABASE d CHARSET latin1 COLLATE latin1_bin")},
		},
		{
			name: "create table stops at a bad table collation",
			sql:  "CREATE TABLE t (a TEXT COLLATE ascii_bin) COLLATE = big5_bin;",
			want: []systemFinding{bad("big5_bin", "CREATE TABLE t (a TEXT COLLATE ascii_bin) COLLATE = big5_bin")},
		},
		{
			name: "create table reports its first bad column, anchoring it",
			sql:  "CREATE TABLE t (a TEXT,\n  b TEXT COLLATE LATIN1_BIN,\n  c TEXT COLLATE ascii_bin);",
			want: []systemFinding{bad("latin1_bin", "b TEXT COLLATE LATIN1_BIN")},
		},
		{
			name: "alter table reports its last bad collation once",
			sql:  alter + ";",
			want: []systemFinding{bad("latin1_bin", alter)},
		},
		{
			name: "alter table convert reads its collate",
			sql:  "ALTER TABLE t CONVERT TO CHARACTER SET gbk COLLATE gbk_bin;",
			want: []systemFinding{bad("gbk_bin", "ALTER TABLE t CONVERT TO CHARACTER SET gbk COLLATE gbk_bin")},
		},
		{name: "alter table convert without a collate", sql: "ALTER TABLE t CONVERT TO CHARACTER SET gbk;"},
		{name: "charset options are not judged", sql: "ALTER TABLE t CHARSET = latin1; CREATE DATABASE d CHARSET latin1;"},
		{name: "a routine body is not judged", sql: "CREATE PROCEDURE p() BEGIN CREATE TABLE t (a TEXT COLLATE ascii_bin); END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := judgeSystem(t, tt.sql, func(stmts []ast.Node) []v1review.Finding {
				return SystemCollationAllowlist(stmts, []string{"UTF8MB4_BIN"})
			})
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
