package v1review

import (
	"testing"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

func TestColumnAutoIncrementMustUnsigned(t *testing.T) {
	bad := func(table, name, text string) placed {
		return placed{Code: v1review.ColumnAutoIncrementNotUnsigned, Name: name, Table: table, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []placed
	}{
		{name: "unsigned and zerofill pass", sql: "CREATE TABLE t (a INT UNSIGNED AUTO_INCREMENT, b INT ZEROFILL AUTO_INCREMENT, c INT);"},
		{
			name: "signed auto increment",
			sql:  "CREATE TABLE t (a INT AUTO_INCREMENT, b DOUBLE AUTO_INCREMENT);\nALTER TABLE t MODIFY a BIGINT AUTO_INCREMENT;",
			want: []placed{
				bad("t", "a", "a INT AUTO_INCREMENT"),
				bad("t", "b", "b DOUBLE AUTO_INCREMENT"),
				bad("t", "a", "ALTER TABLE t MODIFY a BIGINT AUTO_INCREMENT"),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			checkPlaced(t, tt.sql, ColumnAutoIncrementMustUnsigned, tt.want)
		})
	}
	t.Run("no type is reported", func(t *testing.T) {
		stmts := []ast.Node{&ast.CreateTableStmt{Table: &ast.TableRef{Name: "t"}, Columns: []*ast.ColumnDef{{Name: "a", AutoIncrement: true}}}}
		if got := ColumnAutoIncrementMustUnsigned(stmts); len(got) != 1 {
			t.Errorf("got %v, want one finding", got)
		}
	})
}
