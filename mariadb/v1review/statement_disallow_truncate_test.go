package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestStatementDisallowTruncate(t *testing.T) {
	truncate := func(table, text string) stmtFinding {
		return stmtFinding{Code: v1review.Truncate, Table: table, Text: text}
	}
	partition := func(table, name, text string) stmtFinding {
		return stmtFinding{Code: v1review.TruncatePartition, Name: name, Table: table, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []stmtFinding
	}{
		{name: "truncate", sql: "TRUNCATE TABLE t;", want: []stmtFinding{truncate("t", "TRUNCATE TABLE t")}},
		{name: "truncate without TABLE", sql: "TRUNCATE db.t;", want: []stmtFinding{truncate("db.t", "TRUNCATE db.t")}},
		{
			name: "truncate partitions",
			sql:  "ALTER TABLE db.t ADD COLUMN c INT,\n  TRUNCATE PARTITION p0, p1;",
			want: []stmtFinding{partition("db.t", "p0", "TRUNCATE PARTITION p0, p1"), partition("db.t", "p1", "TRUNCATE PARTITION p0, p1")},
		},
		{
			name: "truncate all partitions",
			sql:  "ALTER TABLE t TRUNCATE PARTITION ALL;",
			want: []stmtFinding{partition("t", "ALL PARTITIONS", "TRUNCATE PARTITION ALL")},
		},
		{name: "drop partition not judged", sql: "ALTER TABLE t DROP PARTITION p0;"},
		{name: "delete not judged", sql: "DELETE FROM t;"},
		{name: "truncate in routine body not judged", sql: "CREATE PROCEDURE p() BEGIN TRUNCATE TABLE t; END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runStatement(t, tt.sql, StatementDisallowTruncate)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
