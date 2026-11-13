package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestStatementNonTransactional(t *testing.T) {
	f := []stmtFinding{{Code: v1review.NonTransactional}}
	for _, tt := range []struct {
		name string
		sql  string
		want []stmtFinding
	}{
		{name: "drop database", sql: "DROP DATABASE d;", want: f},
		{name: "drop database lower case", sql: "drop database if exists d", want: f},
		{name: "drop database with two spaces", sql: "DROP  DATABASE d;"},
		{name: "create index concurrently", sql: "CREATE INDEX CONCURRENTLY i ON t (a);", want: f},
		{name: "create unique index concurrently", sql: "CREATE UNIQUE\nINDEX CONCURRENTLY i ON t (a);", want: f},
		{name: "create index", sql: "CREATE INDEX i ON t (a);"},
		{name: "drop index concurrently", sql: "DROP INDEX CONCURRENTLY i;", want: f},
		{name: "drop index", sql: "DROP INDEX i;"},
		{name: "vacuum", sql: "VACUUM t;", want: f},
		{name: "vacuum after white space", sql: "\n  vacuum full t;\n", want: f},
		{name: "vacuum after a comment", sql: "-- c\nVACUUM t;"},
		{name: "analyze", sql: "ANALYZE t;"},
		{name: "pattern in a quoted name", sql: "DROP TABLE \"DROP DATABASE\";", want: f},
		{name: "pattern in a comment", sql: "/* DROP INDEX CONCURRENTLY */ DROP TABLE t;", want: f},
		{name: "pattern on another kind", sql: "SELECT 'DROP DATABASE';"},
		{name: "reindex concurrently not judged", sql: "REINDEX INDEX CONCURRENTLY i;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := toStmtFindings(tt.sql, StatementNonTransactional(parseStatements(t, tt.sql), tt.sql))
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
