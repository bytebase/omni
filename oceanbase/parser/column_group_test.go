package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// TestColumnGroup pins where WITH COLUMN GROUP and ALTER TABLE's column group
// actions land. Whether each form parses is held to OceanBase by the
// with_column_group and alter_column_group_option rows of testdata/compat.tsv.
func TestColumnGroup(t *testing.T) {
	one := func(sql string) nodes.Node {
		t.Helper()
		list, err := Parse(sql)
		if err != nil || len(list.Items) != 1 {
			t.Fatalf("Parse(%q) = %v, %v", sql, list, err)
		}
		return list.Items[0]
	}
	if s := one("CREATE TABLE t (a INT, b INT, KEY k (a) WITH COLUMN GROUP (each column)) WITH COLUMN GROUP (all columns, g (a, b))").(*nodes.CreateTableStmt); s.ColumnGroup != "all columns, g (a, b)" || s.Constraints[0].ColumnGroup != "each column" {
		t.Errorf("CREATE TABLE column groups = %q, %q", s.ColumnGroup, s.Constraints[0].ColumnGroup)
	}
	if s := one("CREATE INDEX i ON t (a) WITH COLUMN GROUP (ALL COLUMNS)").(*nodes.CreateIndexStmt); s.ColumnGroup != "ALL COLUMNS" {
		t.Errorf("CREATE INDEX column group = %q", s.ColumnGroup)
	}
	cmd := one("ALTER TABLE t DROP COLUMN GROUP (EACH COLUMN) DELAYED").(*nodes.AlterTableStmt).Commands[0]
	if cmd.Type != nodes.ATDropColumnGroup || cmd.Name != "EACH COLUMN" || cmd.Option == nil || cmd.Option.Name != "DELAYED" {
		t.Errorf("DROP COLUMN GROUP = %+v", cmd)
	}
	for _, sql := range []string{
		"CREATE TABLE t (a INT) WITH COLUMN GROUP ()",
		"CREATE TABLE t (a INT) WITH COLUMN GROUP (ALL)",
		"CREATE TABLE t (a INT, PRIMARY KEY (a) WITH COLUMN GROUP (ALL COLUMNS))", // a primary key takes none
		"ALTER TABLE t ADD COLUMN GROUP (g ())",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) was accepted", sql)
		}
	}
}
