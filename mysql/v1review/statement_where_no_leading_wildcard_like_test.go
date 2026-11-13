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
		{name: "leading wildcard", sql: "SELECT * FROM t WHERE a LIKE '%x';", want: []stmtFinding{like("a LIKE '%x'")}},
		{name: "not like", sql: "DELETE FROM t WHERE a NOT LIKE '%x';", want: []stmtFinding{like("a NOT LIKE '%x'")}},
		{name: "select list judged too", sql: "SELECT a LIKE '%x' FROM t;", want: []stmtFinding{like("a LIKE '%x'")}},
		{name: "trailing wildcard", sql: "SELECT * FROM t WHERE a LIKE 'x%';"},
		{name: "pattern not a literal", sql: "SELECT * FROM t WHERE a LIKE CONCAT('%', b);"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runStatement(t, tt.sql, StatementWhereNoLeadingWildcardLike)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
