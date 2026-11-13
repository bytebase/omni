package v1review

import (
	"testing"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/v1review"
)

func TestColumnMaximumCharacterLength(t *testing.T) {
	rule := func(maximum int) func([]ast.Node) []v1review.Finding {
		return func(stmts []ast.Node) []v1review.Finding { return ColumnMaximumCharacterLength(stmts, maximum) }
	}
	runColumnCases(t, rule(20), []columnCase{
		{name: "within the limit, varchar not judged", sql: "CREATE TABLE t(a CHAR(20), b VARCHAR(100));"},
		{
			name: "create, add, modify and change",
			sql:  "CREATE TABLE db.t(a char(21));\nALTER TABLE t ADD COLUMN (b CHAR(30)), MODIFY c CHAR(40), CHANGE d e CHAR(50);",
			want: []string{
				"column-char-too-long t.a  @0: a char(21)",
				"column-char-too-long t.b  @0: b CHAR(30)",
				"column-char-too-long t.c  @0: c CHAR(40)",
				"column-char-too-long t.e  @0: e CHAR(50)",
			},
		},
	})
	runColumnCases(t, rule(0), []columnCase{
		{name: "no positive maximum judges nothing", sql: "CREATE TABLE t(a CHAR(255));"},
	})
}
