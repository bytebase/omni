package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

func TestSystemCommentLength(t *testing.T) {
	long := func(text string) systemFinding {
		return systemFinding{Code: v1review.CommentOnTooLong, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		max  int
		want []systemFinding
	}{
		{name: "within the limit", sql: "COMMENT ON TABLE t IS 'abcde';", max: 5},
		{name: "over the limit", sql: "COMMENT ON TABLE t IS 'abcdef';", max: 5, want: []systemFinding{long("COMMENT ON TABLE t IS 'abcdef'")}},
		{name: "counted in bytes", sql: "COMMENT ON COLUMN t.a IS 'été';", max: 3, want: []systemFinding{long("COMMENT ON COLUMN t.a IS 'été'")}},
		{name: "a limit of 0 requires nothing", sql: "COMMENT ON TABLE t IS 'abcdef';", max: 0},
		{name: "a negative limit requires nothing", sql: "COMMENT ON TABLE t IS 'abcdef';", max: -1},
		{name: "a null comment is empty", sql: "COMMENT ON TABLE t IS NULL;", max: 1},
		{
			name: "any object kind",
			sql:  "COMMENT ON SCHEMA s IS 'abcdef'; COMMENT ON INDEX i IS 'abcdef'; COMMENT ON FUNCTION f() IS 'abcdef';",
			max:  5,
			want: []systemFinding{
				long("COMMENT ON SCHEMA s IS 'abcdef'"),
				long("COMMENT ON INDEX i IS 'abcdef'"),
				long("COMMENT ON FUNCTION f() IS 'abcdef'"),
			},
		},
		{
			name: "a comment outside COMMENT ON is not judged",
			sql:  "CREATE FUNCTION f() RETURNS void LANGUAGE sql AS $$ COMMENT ON TABLE t IS 'abcdef' $$; SELECT 'abcdef';",
			max:  5,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := judgeSystem(t, tt.sql, func(stmts []ast.Node) []v1review.Finding {
				return SystemCommentLength(stmts, tt.max)
			})
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
