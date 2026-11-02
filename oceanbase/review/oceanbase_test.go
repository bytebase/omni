package review

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/review"
)

// TestReviewOceanBase checks the OceanBase entry: the MySQL rules, but
// prior backup reported as one finding for the request.
func TestReviewOceanBase(t *testing.T) {
	targets := []review.Target{{Schema: shop(), BackupDatabaseExists: false}}
	for _, tt := range []struct {
		name   string
		rule   review.Rule
		sql    string
		change review.Change
		want   []finding
	}{
		{
			name:   "prior backup is the request itself",
			rule:   review.PriorBackup,
			sql:    "UPDATE t SET n = 1 WHERE id = 1;\nDELETE FROM t WHERE id = 2;",
			change: review.Change{PriorBackup: true, MaxBackupSize: 10},
			want:   []finding{{-1, "", "the change enables prior backup, which Bytebase does not take on OceanBase: the change runs without a backup", []int{0}}},
		},
		{
			name: "prior backup not asked for",
			rule: review.PriorBackup,
			sql:  "UPDATE t SET n = 1 WHERE id = 1;",
		},
		{
			name: "a column of the primary key dropped, which OceanBase refuses",
			rule: review.RequirePrimaryKey,
			sql:  "ALTER TABLE plain DROP COLUMN id;",
		},
		{
			name: "the primary key dropped",
			rule: review.RequirePrimaryKey,
			sql:  "ALTER TABLE plain DROP PRIMARY KEY;",
			want: []finding{{0, "DROP PRIMARY KEY", "removes the primary key of plain, and the change adds none back", []int{0}}},
		},
		{
			name: "the MySQL rules",
			rule: review.RequireWhere,
			sql:  "DELETE FROM t;",
			want: []finding{{0, "DELETE FROM t", "DELETE FROM t has no WHERE clause", []int{0}}},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ReviewOceanBase(context.Background(), tt.sql, review.Options{Rules: []review.Rule{tt.rule}, Change: tt.change}, targets)
			if err != nil {
				t.Fatal(err)
			}
			var got []finding
			for _, f := range result.Findings {
				if f.Rule == tt.rule {
					got = append(got, finding{f.Statement, tt.sql[f.Range.Start:f.Range.End], f.Message, f.Targets})
				}
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
