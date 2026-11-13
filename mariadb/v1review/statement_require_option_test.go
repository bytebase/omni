package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/v1review"
)

func TestStatementRequireOption(t *testing.T) {
	for _, tt := range []struct {
		name string
		rule func([]ast.Node) []v1review.Finding
		sql  string
		want []stmtFinding
	}{
		{name: "algorithm missing", rule: StatementRequireAlgorithmOption, sql: "ALTER TABLE t ADD a INT, LOCK=NONE;", want: []stmtFinding{{Code: v1review.OptionMissing, Name: "ALGORITHM"}}},
		{name: "algorithm given", rule: StatementRequireAlgorithmOption, sql: "ALTER TABLE t ADD a INT, ALGORITHM=INPLACE;"},
		{name: "lock missing", rule: StatementRequireLockOption, sql: "ALTER TABLE t ALGORITHM=INSTANT, ADD a INT;", want: []stmtFinding{{Code: v1review.OptionMissing, Name: "LOCK"}}},
		{name: "lock given", rule: StatementRequireLockOption, sql: "ALTER TABLE t ADD a INT, LOCK=DEFAULT;"},
		{name: "other statements not judged", rule: StatementRequireLockOption, sql: "CREATE TABLE t (a INT); CREATE PROCEDURE p() BEGIN ALTER TABLE t ADD a INT; END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runStatement(t, tt.sql, tt.rule)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
