package v1review

import (
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/mysql/parser"
	"github.com/bytebase/omni/v1review"
)

func TestStatementMergeAlterTable(t *testing.T) {
	redundant := func(table string, count, stmt int) v1review.Finding {
		return v1review.Finding{Code: v1review.RedundantAlterTable, Table: table, Value: strconv.Itoa(count), Statement: stmt}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []v1review.Finding
	}{
		{name: "two alters", sql: "ALTER TABLE t ADD a INT; ALTER TABLE t ADD b INT;", want: []v1review.Finding{redundant("t", 2, 1)}},
		{name: "create then alter", sql: "CREATE TABLE t (a INT); SELECT 1; ALTER TABLE t ADD b INT;", want: []v1review.Finding{redundant("t", 2, 2)}},
		{name: "create resets the count", sql: "ALTER TABLE t ADD a INT; ALTER TABLE t ADD b INT; CREATE TABLE t (a INT);"},
		{name: "qualifier dropped", sql: "ALTER TABLE db1.t ADD a INT; ALTER TABLE db2.t ADD b INT;", want: []v1review.Finding{redundant("t", 2, 1)}},
		{name: "case sensitive", sql: "ALTER TABLE t ADD a INT; ALTER TABLE T ADD b INT;"},
		{
			name: "ordered by last statement",
			sql:  "ALTER TABLE b ADD a INT; ALTER TABLE a ADD a INT; ALTER TABLE a ADD b INT; ALTER TABLE b ADD b INT;",
			want: []v1review.Finding{redundant("a", 2, 2), redundant("b", 2, 3)},
		},
		{name: "single alter", sql: "CREATE TABLE t (a INT); ALTER TABLE u ADD a INT;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			list, err := parser.Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			if diff := cmp.Diff(tt.want, StatementMergeAlterTable(list.Items)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
