package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestTableDisallowPartition(t *testing.T) {
	f := func(text string) finding { return finding{v1review.TablePartitioned, "", text} }
	for _, tt := range []struct {
		name string
		sql  string
		want []finding
	}{
		{
			name: "create with partitions",
			sql:  "CREATE TABLE t (a INT) PARTITION BY HASH(a) PARTITIONS 4;",
			want: []finding{f("CREATE TABLE t (a INT) PARTITION BY HASH(a) PARTITIONS 4")},
		},
		{name: "create without", sql: "CREATE TABLE t (a INT);"},
		{
			name: "alter partition by",
			sql:  "ALTER TABLE t PARTITION BY KEY(a) PARTITIONS 2;",
			want: []finding{f("ALTER TABLE t PARTITION BY KEY(a) PARTITIONS 2")},
		},
		{name: "other partition actions not judged", sql: "ALTER TABLE t ADD PARTITION (PARTITION p3 VALUES LESS THAN (30));"},
		{name: "remove partitioning not judged", sql: "ALTER TABLE t REMOVE PARTITIONING;"},
		{name: "routine body not judged", sql: "CREATE PROCEDURE p() BEGIN CREATE TABLE t (a INT) PARTITION BY HASH(a); END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, tableFindings(t, tt.sql, TableDisallowPartition)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
