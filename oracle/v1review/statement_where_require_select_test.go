package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestStatementWhereRequireSelect(t *testing.T) {
	missing := func(text string) stmtFinding { return stmtFinding{Code: v1review.WhereMissing, Text: text} }
	for _, tt := range []struct {
		name string
		sql  string
		want []stmtFinding
	}{
		{name: "no where", sql: "SELECT a FROM t", want: []stmtFinding{missing("SELECT a FROM t")}},
		{name: "where", sql: "SELECT a FROM t WHERE b = 1"},
		{name: "dual not judged", sql: "SELECT 1 FROM Dual"},
		{name: "dual joined", sql: "SELECT 1 FROM dual, t", want: []stmtFinding{missing("SELECT 1 FROM dual, t")}},
		{name: "subquery", sql: "SELECT a FROM t WHERE b IN (SELECT b FROM u)", want: []stmtFinding{missing("SELECT b FROM u")}},
		{name: "update not judged", sql: "UPDATE t SET a = 1"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runStatement(t, tt.sql, StatementWhereRequireSelect)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
