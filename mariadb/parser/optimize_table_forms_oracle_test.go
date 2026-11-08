package parser

import "testing"

// TestOptimizeTableOrTables holds OPTIMIZE's TABLE/TABLES keyword to MariaDB
// 11.8. table_or_tables is TABLE or TABLES and is required, so `OPTIMIZE TABLES
// t` parses and `OPTIMIZE t` with neither is a 1064; each is measured against
// the server and omni must agree.
func TestOptimizeTableOrTables(t *testing.T) {
	o := startMariaDB(t)
	for _, sql := range []string{
		"OPTIMIZE TABLE t1",
		"OPTIMIZE TABLES t1, t2",
		"OPTIMIZE NO_WRITE_TO_BINLOG TABLES t1",
		"OPTIMIZE LOCAL TABLE t1",
		"OPTIMIZE t1",
		"OPTIMIZE",
	} {
		t.Run(sql, func(t *testing.T) {
			mariaAccepts := o.prepareAccepts(t, sql)
			omniAccepts, perr := parsesAsOneStatement(sql)
			if omniAccepts != mariaAccepts {
				t.Errorf("DIVERGENCE omni accepts=%v, MariaDB accepts=%v for %q: %v",
					omniAccepts, mariaAccepts, sql, perr)
			}
		})
	}
}
