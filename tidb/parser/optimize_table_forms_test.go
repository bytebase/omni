package parser

import "testing"

// TestOptimizeTableOrTables pins that OPTIMIZE takes TABLE or TABLES, and
// requires one. TiDB v8.5.5's grammar is `OPTIMIZE NoWriteToBinLogAliasOpt
// TableOrTables TableNameList` with `TableOrTables: "TABLE" | "TABLES"`, so
// `OPTIMIZE TABLES t` parses (the §3.4 rest-corpus gap) and `OPTIMIZE t` with
// neither answers 1064 (measured against v8.5.5; omni accepted it before).
func TestOptimizeTableOrTables(t *testing.T) {
	for _, c := range []struct {
		sql    string
		accept bool
	}{
		{"OPTIMIZE TABLE t1", true},
		{"OPTIMIZE TABLES t1, t2", true},
		{"OPTIMIZE NO_WRITE_TO_BINLOG TABLES t1", true},
		{"OPTIMIZE LOCAL TABLE t1", true},
		{"OPTIMIZE t1", false},
		{"OPTIMIZE", false},
	} {
		_, err := Parse(c.sql)
		if accepts := err == nil; accepts != c.accept {
			t.Errorf("Parse(%q) accepts=%v, want %v (err=%v)", c.sql, accepts, c.accept, err)
		}
	}
}
