package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/mariadb/parser"
	"github.com/bytebase/omni/v1review"
)

// partitionFinding is a TableDisallowPartition finding as the tests compare
// it, with Text the source its range anchors.
type partitionFinding struct {
	Code  v1review.Code
	Table string
	Text  string
}

func TestTableDisallowPartition(t *testing.T) {
	partitioned := func(table, text string) partitionFinding {
		return partitionFinding{v1review.TablePartitioned, table, text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []partitionFinding
	}{
		{name: "create without partitions", sql: "CREATE TABLE t (a INT);"},
		{
			name: "create with partitions, name without database",
			sql:  "CREATE TABLE db.t (a INT)\nPARTITION BY HASH(a) PARTITIONS 4;",
			want: []partitionFinding{partitioned("t", "CREATE TABLE db.t (a INT)\nPARTITION BY HASH(a) PARTITIONS 4")},
		},
		{
			name: "alter partition by",
			sql:  "ALTER TABLE t PARTITION BY KEY(a) PARTITIONS 2;",
			want: []partitionFinding{partitioned("t", "ALTER TABLE t PARTITION BY KEY(a) PARTITIONS 2")},
		},
		{
			name: "alter with another command first",
			sql:  "ALTER TABLE t ADD COLUMN b INT PARTITION BY KEY(a) PARTITIONS 2;",
			want: []partitionFinding{partitioned("t", "ALTER TABLE t ADD COLUMN b INT PARTITION BY KEY(a) PARTITIONS 2")},
		},
		{name: "partition maintenance not judged", sql: "ALTER TABLE t ADD PARTITION (PARTITION p1 VALUES LESS THAN (20));"},
		{name: "drop partition not judged", sql: "ALTER TABLE t DROP PARTITION p0;"},
		{name: "remove partitioning not judged", sql: "ALTER TABLE t REMOVE PARTITIONING;"},
		{name: "like not judged", sql: "CREATE TABLE t2 LIKE t;"},
		{name: "routine body not judged", sql: "CREATE PROCEDURE p() BEGIN CREATE TABLE t (a INT) PARTITION BY HASH(a); END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			list, err := parser.Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			var got []partitionFinding
			for _, f := range TableDisallowPartition(list.Items) {
				got = append(got, partitionFinding{f.Code, f.Table, tt.sql[f.Start:f.End]})
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
