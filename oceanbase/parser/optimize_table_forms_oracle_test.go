package parser

import "testing"

// TestOptimizeTableOrTables holds OPTIMIZE's TABLE/TABLES keyword to OceanBase (MySQL mode).
// table_or_tables is TABLE or TABLES and is required, so `OPTIMIZE TABLES t`
// parses and `OPTIMIZE t` with neither is a 1064. assertOracleMatch measures
// each against the server.
func TestOptimizeTableOrTables(t *testing.T) {
	o := startParserOracle(t)
	for _, sql := range []string{
		"OPTIMIZE TABLE t1",
		"OPTIMIZE TABLES t1, t2",
		"OPTIMIZE NO_WRITE_TO_BINLOG TABLES t1",
		"OPTIMIZE LOCAL TABLE t1",
		"OPTIMIZE t1",
		"OPTIMIZE",
	} {
		assertOracleMatch(t, o, sql)
	}
}
