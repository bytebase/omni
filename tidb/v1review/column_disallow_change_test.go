package v1review

import (
	"testing"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

func TestColumnDisallowChange(t *testing.T) {
	for _, tt := range []struct {
		name string
		sql  string
		want []placed
	}{
		{name: "other actions pass", sql: "ALTER TABLE t MODIFY a INT, RENAME COLUMN b TO c, ADD d INT; CREATE TABLE u (a INT);"},
		{
			name: "once per statement",
			sql:  "ALTER TABLE t CHANGE a b INT, CHANGE COLUMN c d INT;\nALTER TABLE u ADD x INT, CHANGE e f INT;",
			want: []placed{
				{Code: v1review.ColumnChanged, Text: "ALTER TABLE t CHANGE a b INT, CHANGE COLUMN c d INT"},
				{Code: v1review.ColumnChanged, Text: "ALTER TABLE u ADD x INT, CHANGE e f INT"},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			checkPlaced(t, tt.sql, ColumnDisallowChange, tt.want)
		})
	}
	t.Run("no table and nil actions", func(t *testing.T) {
		stmts := []ast.Node{&ast.AlterTableStmt{Commands: []*ast.AlterTableCmd{nil, {Type: ast.ATChangeColumn}}}}
		if got := ColumnDisallowChange(stmts); len(got) != 1 {
			t.Errorf("got %v, want one finding", got)
		}
	})
}
