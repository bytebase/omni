package catalog

import (
	"strings"
	"testing"
)

// ADD PARTITION PARTITIONS n appends n partitions and names them the way
// CREATE TABLE does, continuing the pN run. TiDB v8.5.5 answers p0..p4 to
// `PARTITION BY HASH(id) PARTITIONS 3` followed by `ADD PARTITION PARTITIONS
// 2`, which is also what MySQL 8.0 does.
//
// TestScenario_C6/6_16_add_partition_auto_naming puts the same statement to a
// live server; this test covers the behaviour without a container:
// the parser only began accepting this form with the partition-maintenance
// modifiers, and the handler ignored cmd.Number until then, which made the
// statement a silent no-op on the catalog.
func TestAlterAddPartitionCount(t *testing.T) {
	for _, tc := range []struct {
		name  string
		alter string
		want  []string
	}{
		{
			name:  "PARTITIONS n continues the pN run",
			alter: "ALTER TABLE t ADD PARTITION PARTITIONS 2",
			want:  []string{"p0", "p1", "p2", "p3", "p4"},
		},
		{
			name:  "the modifier does not displace the count",
			alter: "ALTER TABLE t ADD PARTITION LOCAL PARTITIONS 1",
			want:  []string{"p0", "p1", "p2", "p3"},
		},
		{
			// v8.5.5 runs this on a HASH table and adds nothing.
			name:  "no count and no definitions adds nothing",
			alter: "ALTER TABLE t ADD PARTITION",
			want:  []string{"p0", "p1", "p2"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := wtSetup(t)
			wtExec(t, c, "CREATE TABLE t (id INT) PARTITION BY HASH(id) PARTITIONS 3")
			wtExec(t, c, tc.alter)

			tbl := c.GetDatabase("testdb").GetTable("t")
			if tbl == nil || tbl.Partitioning == nil {
				t.Fatal("table t has no partitioning")
			}
			names := make([]string, len(tbl.Partitioning.Partitions))
			for i, p := range tbl.Partitioning.Partitions {
				names[i] = p.Name
			}
			if got := strings.Join(names, ","); got != strings.Join(tc.want, ",") {
				t.Errorf("partition names = %q, want %q", got, strings.Join(tc.want, ","))
			}
			if tbl.Partitioning.NumParts != len(tc.want) {
				t.Errorf("NumParts = %d, want %d", tbl.Partitioning.NumParts, len(tc.want))
			}
		})
	}
}
