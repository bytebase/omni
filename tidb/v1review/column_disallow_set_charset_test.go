package v1review

import (
	"testing"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

func TestColumnDisallowSetCharset(t *testing.T) {
	for _, tt := range []struct {
		name string
		sql  string
		want []placed
	}{
		{name: "no charset or binary passes", sql: "CREATE TABLE t (a VARCHAR(5), b VARCHAR(5) CHARACTER SET BINARY, c VARCHAR(5) COLLATE utf8mb4_bin) CHARSET utf8mb4;"},
		{
			name: "once per statement",
			sql:  "CREATE TABLE t (a VARCHAR(5) CHARSET utf8mb4, b TEXT CHARACTER SET latin1);\nALTER TABLE t ADD (c INT, d CHAR(1) CHARSET ascii);\nALTER TABLE t MODIFY e TEXT CHARSET utf8;\nALTER TABLE t CONVERT TO CHARACTER SET utf8;",
			want: []placed{
				{Code: v1review.ColumnCharsetSet, Text: "CREATE TABLE t (a VARCHAR(5) CHARSET utf8mb4, b TEXT CHARACTER SET latin1)"},
				{Code: v1review.ColumnCharsetSet, Text: "ALTER TABLE t ADD (c INT, d CHAR(1) CHARSET ascii)"},
				{Code: v1review.ColumnCharsetSet, Text: "ALTER TABLE t MODIFY e TEXT CHARSET utf8"},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			checkPlaced(t, tt.sql, ColumnDisallowSetCharset, tt.want)
		})
	}
	t.Run("no table is judged", func(t *testing.T) {
		col := &ast.ColumnDef{Name: "a", TypeName: &ast.DataType{Name: "TEXT", Charset: "latin1"}}
		stmts := []ast.Node{
			&ast.CreateTableStmt{Columns: []*ast.ColumnDef{nil, col}},
			&ast.AlterTableStmt{Commands: []*ast.AlterTableCmd{nil, {Type: ast.ATModifyColumn, Column: col}}},
		}
		if got := ColumnDisallowSetCharset(stmts); len(got) != 2 {
			t.Errorf("got %v, want two findings", got)
		}
	})
}
