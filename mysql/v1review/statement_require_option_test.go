package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestStatementRequireOption(t *testing.T) {
	algorithm := []stmtFinding{{Code: v1review.OptionMissing, Name: "ALGORITHM"}}
	lock := []stmtFinding{{Code: v1review.OptionMissing, Name: "LOCK"}}
	for _, tt := range []struct {
		name               string
		sql                string
		wantAlgo, wantLock []stmtFinding
	}{
		{name: "neither", sql: "ALTER TABLE t ADD a INT;", wantAlgo: algorithm, wantLock: lock},
		{name: "algorithm", sql: "ALTER TABLE t ADD a INT, ALGORITHM = INPLACE;", wantLock: lock},
		{name: "lock", sql: "ALTER TABLE t ADD a INT, LOCK = NONE;", wantAlgo: algorithm},
		{name: "both", sql: "ALTER TABLE t ADD a INT, ALGORITHM = INSTANT, LOCK = DEFAULT;"},
		{name: "create index not judged", sql: "CREATE INDEX i ON t (a);"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.wantAlgo, runStatement(t, tt.sql, StatementRequireAlgorithmOption)); diff != "" {
				t.Errorf("ALGORITHM findings (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantLock, runStatement(t, tt.sql, StatementRequireLockOption)); diff != "" {
				t.Errorf("LOCK findings (-want +got):\n%s", diff)
			}
		})
	}
}
