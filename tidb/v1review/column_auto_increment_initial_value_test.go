package v1review

import (
	"testing"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

func TestColumnAutoIncrementInitialValue(t *testing.T) {
	for _, tt := range []struct {
		name  string
		sql   string
		value int
		want  []placed
	}{
		{name: "equal passes", sql: "CREATE TABLE t (a INT) AUTO_INCREMENT = 20;", value: 20},
		{
			name:  "differs, anchored at the statement",
			sql:   "\n-- c\nCREATE TABLE db.t (a INT)\n  auto_increment = 7;",
			value: 20,
			want:  []placed{{Code: v1review.ColumnAutoIncrementInitialValue, Table: "t", Value: "7", Text: "CREATE TABLE db.t (a INT)\n  auto_increment = 7"}},
		},
		{
			name:  "every option is judged",
			sql:   "CREATE TABLE t (a INT) AUTO_INCREMENT 5, AUTO_INCREMENT = 6;",
			value: 6,
			want:  []placed{{Code: v1review.ColumnAutoIncrementInitialValue, Table: "t", Value: "5", Text: "CREATE TABLE t (a INT) AUTO_INCREMENT 5, AUTO_INCREMENT = 6"}},
		},
		{name: "a negative value wraps", sql: "CREATE TABLE t (a INT) AUTO_INCREMENT = 18446744073709551615;", value: -1},
		{name: "alter is not judged", sql: "ALTER TABLE t AUTO_INCREMENT = 5;", value: 20},
		{name: "other options are not judged", sql: "CREATE TABLE t (a BIGINT AUTO_RANDOM) AUTO_RANDOM_BASE = 5, AUTO_ID_CACHE = 1;", value: 20},
	} {
		t.Run(tt.name, func(t *testing.T) {
			checkPlaced(t, tt.sql, func(stmts []ast.Node) []v1review.Finding {
				return ColumnAutoIncrementInitialValue(stmts, tt.value)
			}, tt.want)
		})
	}
	t.Run("unreadable values are skipped", func(t *testing.T) {
		stmts := []ast.Node{&ast.CreateTableStmt{Table: &ast.TableRef{Name: "t"}, Options: []*ast.TableOption{nil, {Name: "AUTO_INCREMENT", Value: "0x10"}, {Name: "AUTO_INCREMENT", Value: "-1"}}}}
		if got := ColumnAutoIncrementInitialValue(stmts, 1); len(got) != 0 {
			t.Errorf("got %v, want none", got)
		}
	})
}
