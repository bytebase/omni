package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestStatementDisallowTruncate(t *testing.T) {
	truncate := func(table, text string) statementFinding {
		return statementFinding{Code: v1review.Truncate, Table: table, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []statementFinding
	}{
		{name: "table", sql: "TRUNCATE TABLE t", want: []statementFinding{truncate("t", "TRUNCATE TABLE t")}},
		{name: "schema kept", sql: "TRUNCATE TABLE dbo.t", want: []statementFinding{truncate("dbo.t", "TRUNCATE TABLE dbo.t")}},
		{name: "database dropped", sql: "TRUNCATE TABLE db.dbo.t", want: []statementFinding{truncate("dbo.t", "TRUNCATE TABLE db.dbo.t")}},
		{name: "bracketed name unquoted", sql: "TRUNCATE TABLE [my t]", want: []statementFinding{truncate("my t", "TRUNCATE TABLE [my t]")}},
		{name: "local temp table not judged", sql: "TRUNCATE TABLE #t"},
		{name: "global temp table not judged", sql: "TRUNCATE TABLE ##t"},
		{name: "in a block not judged", sql: "BEGIN TRUNCATE TABLE t; END"},
		{name: "delete not judged", sql: "DELETE FROM t"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runStatement(t, tt.sql, StatementDisallowTruncate)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
