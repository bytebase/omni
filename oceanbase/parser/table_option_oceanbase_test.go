package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// TestOceanBaseOwnTableOptions pins the name and value each of OceanBase's own
// table options parses into. Whether each parses is held to OceanBase by its
// rows in testdata/compat.tsv (rule table_option and the rules under it).
func TestOceanBaseOwnTableOptions(t *testing.T) {
	for _, tc := range []struct{ option, name, value string }{
		{"TABLEGROUP = tg1", "TABLEGROUP", "tg1"},
		{"TABLEGROUP 'tg1'", "TABLEGROUP", "'tg1'"},
		{"READ ONLY", "READ ONLY", ""},
		{"PARALLEL 4", "PARALLEL", "4"},
		{"NOPARALLEL", "NOPARALLEL", ""},
		{"TABLE_MODE = 'QUEUING'", "TABLE_MODE", "'QUEUING'"},
		{"TTL (c + INTERVAL 1 DAY, d.c + INTERVAL 2 HOUR)", "TTL", "(c + INTERVAL 1 DAY, d.c + INTERVAL 2 HOUR)"},
		{"FORMAT = (TYPE = 'CSV' FIELD_DELIMITER = ',', SKIP_HEADER = 1)", "FORMAT", "(TYPE = 'CSV' FIELD_DELIMITER = ',', SKIP_HEADER = 1)"},
		{"LOCATION = 'oss://b/d'", "LOCATION", "'oss://b/d'"},
		{"STORAGE_CACHE_POLICY (BOUNDARY_COLUMN = c, HOT_RETENTION = 1 DAY)", "STORAGE_CACHE_POLICY", "(BOUNDARY_COLUMN = c, HOT_RETENTION = 1 DAY)"},
		{"MERGE_ENGINE = PARTIAL_UPDATE", "MERGE_ENGINE", "PARTIAL_UPDATE"},
		{"JSON(c) STORE AS (CHUNK 1 CHUNK '2K')", "JSON", "(c) STORE AS (CHUNK 1 CHUNK '2K')"},
	} {
		for _, sql := range []string{"CREATE TABLE t (c INT) " + tc.option, "ALTER TABLE t " + tc.option} {
			list, err := Parse(sql)
			if err != nil {
				t.Errorf("Parse(%q): %v", sql, err)
				continue
			}
			var opt *nodes.TableOption
			switch stmt := list.Items[0].(type) {
			case *nodes.CreateTableStmt:
				opt = stmt.Options[0]
			case *nodes.AlterTableStmt:
				opt = stmt.Commands[0].Option
			}
			if opt == nil || opt.Name != tc.name || opt.Value != tc.value {
				t.Errorf("Parse(%q): option = %+v, want %s = %q", sql, opt, tc.name, tc.value)
			}
		}
	}
	for _, sql := range []string{
		"ALTER TABLE t TABLE_MODE = 1",       // a string
		"ALTER TABLE t PARALLEL",             // a degree
		"ALTER TABLE t AUTO_REFRESH = DAILY", // OFF, IMMEDIATE or INTERVAL
		"ALTER TABLE t TTL (c + INTERVAL 1 WEEK)",
		"CREATE TABLE t NOPARALLEL", // no definitions and no query
		"CREATE TABLE t",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) was accepted", sql)
		}
	}
}
