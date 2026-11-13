package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/pg/parser"
)

func TestTableDisallowPartition(t *testing.T) {
	for _, tt := range []struct {
		name string
		sql  string
		want []string
	}{
		{name: "plain table", sql: "CREATE TABLE t (a int);"},
		{
			name: "partition by",
			sql:  "CREATE TABLE t (a int) PARTITION BY RANGE (a);",
			want: []string{"CREATE TABLE t (a int) PARTITION BY RANGE (a)"},
		},
		{name: "partition of without partition by", sql: "CREATE TABLE p PARTITION OF t FOR VALUES FROM (1) TO (2);"},
		{
			name: "partition of with partition by",
			sql:  "CREATE TABLE p PARTITION OF t FOR VALUES FROM (1) TO (2) PARTITION BY LIST (a);",
			want: []string{"CREATE TABLE p PARTITION OF t FOR VALUES FROM (1) TO (2) PARTITION BY LIST (a)"},
		},
		{
			name: "attach partition, once per statement",
			sql:  "ALTER TABLE t ATTACH PARTITION p FOR VALUES FROM (1) TO (2), ATTACH PARTITION q DEFAULT;",
			want: []string{"ALTER TABLE t ATTACH PARTITION p FOR VALUES FROM (1) TO (2), ATTACH PARTITION q DEFAULT"},
		},
		{
			name: "alter index attach partition",
			sql:  "ALTER INDEX i ATTACH PARTITION j;",
			want: []string{"ALTER INDEX i ATTACH PARTITION j"},
		},
		{name: "detach partition", sql: "ALTER TABLE t DETACH PARTITION p;"},
		{name: "other alter", sql: "ALTER TABLE t ADD COLUMN b int;"},
		{name: "explain not judged", sql: "EXPLAIN CREATE TABLE t (a int) PARTITION BY RANGE (a);"},
		{name: "routine body not judged", sql: "CREATE FUNCTION f() RETURNS void AS $$ CREATE TABLE t (a int) PARTITION BY RANGE (a) $$ LANGUAGE sql;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			list, err := parser.Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			var got []string
			for _, f := range TableDisallowPartition(list.Items) {
				if f.Code != "table-partitioned" {
					t.Errorf("code %q", f.Code)
				}
				got = append(got, tt.sql[f.Start:f.End])
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
