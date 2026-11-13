package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestStatementWhereNoLeadingWildcardLike(t *testing.T) {
	like := func(text string) stmtFinding { return stmtFinding{Code: v1review.LeadingWildcardLike, Text: text} }
	for _, tt := range []struct {
		name string
		sql  string
		want []stmtFinding
	}{
		{name: "like", sql: "SELECT a FROM t WHERE b LIKE '%x';", want: []stmtFinding{like("b LIKE '%x'")}},
		{name: "not like", sql: "DELETE FROM t WHERE b NOT LIKE '%x';", want: []stmtFinding{like("b NOT LIKE '%x'")}},
		{name: "select list", sql: "SELECT a LIKE '%x' FROM t;", want: []stmtFinding{like("a LIKE '%x'")}},
		{name: "routine body", sql: "CREATE PROCEDURE p() BEGIN SELECT a FROM t WHERE b LIKE '%x'; END;", want: []stmtFinding{like("b LIKE '%x'")}},
		{name: "trailing wildcard", sql: "SELECT a FROM t WHERE b LIKE 'x%';"},
		{name: "pattern not a literal", sql: "SELECT a FROM t WHERE b LIKE CONCAT('%', 'x');"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runStatement(t, tt.sql, StatementWhereNoLeadingWildcardLike)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
