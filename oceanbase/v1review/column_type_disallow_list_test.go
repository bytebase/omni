package v1review

import (
	"testing"

	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/v1review"
)

func TestColumnTypeDisallowList(t *testing.T) {
	rule := func(types ...string) func([]ast.Node) []v1review.Finding {
		return func(stmts []ast.Node) []v1review.Finding { return ColumnTypeDisallowList(stmts, types) }
	}
	runColumnCases(t, rule("json", "BLOB", "VARCHAR(20)"), []columnCase{
		{name: "allowed types, a length in the list matches nothing", sql: "CREATE TABLE t(a INT, b VARCHAR(20));"},
		{
			name: "create and alter, any case",
			sql:  "CREATE TABLE db.t(a Json);\nALTER TABLE t ADD COLUMN (b blob), MODIFY c JSON, CHANGE d e BLOB;",
			want: []string{
				"column-type-disallowed t.a JSON @0: a Json",
				"column-type-disallowed t.b BLOB @0: b blob",
				"column-type-disallowed t.c JSON @0: c JSON",
				"column-type-disallowed t.e BLOB @0: e BLOB",
			},
		},
	})
	runColumnCases(t, rule(), []columnCase{
		{name: "an empty list allows everything", sql: "CREATE TABLE t(a JSON);"},
	})
}
