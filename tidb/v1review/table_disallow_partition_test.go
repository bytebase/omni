package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestTableDisallowPartition(t *testing.T) {
	partitioned := func(table, text string) located {
		return located{Code: v1review.TablePartitioned, Table: table, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []located
	}{
		{name: "create without partitions", sql: "CREATE TABLE t (a INT);"},
		{
			name: "create with partitions, name without database",
			sql:  "CREATE TABLE db.t (a INT)\nPARTITION BY HASH(a) PARTITIONS 4;",
			want: []located{partitioned("t", "CREATE TABLE db.t (a INT)\nPARTITION BY HASH(a) PARTITIONS 4")},
		},
		{
			name: "alter partition by",
			sql:  "ALTER TABLE t PARTITION BY KEY(a) PARTITIONS 2;",
			want: []located{partitioned("t", "ALTER TABLE t PARTITION BY KEY(a) PARTITIONS 2")},
		},
		{name: "partition maintenance not judged", sql: "ALTER TABLE t ADD PARTITION (PARTITION p1 VALUES LESS THAN (20));"},
		{name: "drop partition not judged", sql: "ALTER TABLE t DROP PARTITION p0;"},
		{name: "remove partitioning not judged", sql: "ALTER TABLE t REMOVE PARTITIONING;"},
		{name: "like not judged", sql: "CREATE TABLE t2 LIKE t;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := locate(tt.sql, TableDisallowPartition(parse(t, tt.sql)))
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
