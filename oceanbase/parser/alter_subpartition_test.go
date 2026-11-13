package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// TestAlterSubpartitionActions pins the command each subpartition action and
// name list parses into, and keeps subpartition, an unreserved word, a column
// name where DROP and ALTER take a column. Whether each form parses is held to
// OceanBase by the alter_partition_option rows of testdata/compat.tsv.
func TestAlterSubpartitionActions(t *testing.T) {
	for _, tc := range []struct {
		sql   string
		typ   nodes.AlterTableCmdType
		names []string
	}{
		{"ALTER TABLE t DROP PARTITION p0 p1, p2", nodes.ATDropPartition, []string{"p0", "p1", "p2"}},
		{"ALTER TABLE t DROP SUBPARTITION (s0, s1)", nodes.ATDropSubpartition, []string{"s0", "s1"}},
		{"ALTER TABLE t TRUNCATE SUBPARTITION s0", nodes.ATTruncateSubpartition, []string{"s0"}},
		{"ALTER TABLE t ALTER SUBPARTITION s0 STORAGE_CACHE_POLICY = 'hot'", nodes.ATAlterPartitionStorageCachePolicy, []string{"s0"}},
		{"ALTER TABLE t MODIFY PARTITION p0 ADD SUBPARTITION (SUBPARTITION s9 VALUES LESS THAN (10))", nodes.ATAddSubpartition, nil},
		{"ALTER TABLE t DROP subpartition", nodes.ATDropColumn, nil},
		{"ALTER TABLE t DROP subpartition RESTRICT", nodes.ATDropColumn, nil},
		{"ALTER TABLE t DROP subpartition CASCADE, ADD c INT", nodes.ATDropColumn, nil},
		{"ALTER TABLE t ALTER subpartition SET DEFAULT 1", nodes.ATAlterColumnDefault, nil},
		{"ALTER TABLE t ALTER subpartition DROP DEFAULT", nodes.ATAlterColumnDefault, nil},
	} {
		list, err := Parse(tc.sql)
		if err != nil {
			t.Errorf("Parse(%q): %v", tc.sql, err)
			continue
		}
		cmd := list.Items[0].(*nodes.AlterTableStmt).Commands[0]
		if cmd.Type != tc.typ {
			t.Errorf("Parse(%q) command type = %v, want %v", tc.sql, cmd.Type, tc.typ)
		}
		if tc.names != nil && !equalStrings(cmd.PartitionNames, tc.names) {
			t.Errorf("Parse(%q) names = %q, want %q", tc.sql, cmd.PartitionNames, tc.names)
		}
	}
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
