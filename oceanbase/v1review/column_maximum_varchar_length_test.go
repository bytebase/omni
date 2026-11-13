package v1review

import (
	"testing"

	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/v1review"
)

func TestColumnMaximumVarcharLength(t *testing.T) {
	rule := func(maximum int) func([]ast.Node) []v1review.Finding {
		return func(stmts []ast.Node) []v1review.Finding { return ColumnMaximumVarcharLength(stmts, maximum) }
	}
	runColumnCases(t, rule(100), []columnCase{
		{name: "within the limit, char not judged", sql: "CREATE TABLE t(a VARCHAR(100), b CHAR(200));"},
		{
			name: "create and alter",
			sql:  "CREATE TABLE t(a varchar(101));\nALTER TABLE u MODIFY b VARCHAR(500);",
			want: []string{
				"column-varchar-too-long t.a  @0: a varchar(101)",
				"column-varchar-too-long u.b  @0: b VARCHAR(500)",
			},
		},
	})
	runColumnCases(t, rule(-1), []columnCase{
		{name: "no positive maximum judges nothing", sql: "CREATE TABLE t(a VARCHAR(255));"},
	})
}
