package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestStatementDisallowCrossDBQueries(t *testing.T) {
	cross := func(db, current string, statement int, text string) statementFinding {
		return statementFinding{Code: v1review.CrossDatabaseQuery, Name: db, Value: current, Statement: statement, Text: text}
	}
	for _, tt := range []struct {
		name    string
		sql     string
		current string
		want    []statementFinding
	}{
		{name: "other database", sql: "SELECT * FROM other.dbo.t", current: "Main", want: []statementFinding{cross("other", "Main", 0, "other.dbo.t")}},
		{name: "current database case-insensitively", sql: "SELECT * FROM MAIN.dbo.t", current: "main"},
		{name: "unqualified not judged", sql: "SELECT * FROM dbo.t", current: "main"},
		{
			name:    "use lower-cases and carries over",
			sql:     "USE Other;\nSELECT * FROM OTHER.dbo.t;\nSELECT * FROM Main.dbo.t;",
			current: "Main",
			want:    []statementFinding{cross("Main", "other", 2, "Main.dbo.t")},
		},
		{
			name:    "join and subquery",
			sql:     "SELECT * FROM a.dbo.t JOIN b.dbo.u ON t.id = u.id WHERE EXISTS (SELECT 1 FROM c.dbo.v)",
			current: "main",
			want:    []statementFinding{cross("a", "main", 0, "a.dbo.t"), cross("b", "main", 0, "b.dbo.u"), cross("c", "main", 0, "c.dbo.v")},
		},
		{name: "use nested in a block does not change it", sql: "BEGIN USE other; END;\nSELECT * FROM other.dbo.t;", current: "main", want: []statementFinding{cross("other", "main", 1, "other.dbo.t")}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := stmtAnchored(tt.sql, StatementDisallowCrossDBQueries(parseStatements(t, tt.sql), tt.current))
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
