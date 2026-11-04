package parser

import "testing"

// The SHOW target table, held to the OceanBase verdicts recorded below. One
// case per target this parser reads and per malformed variant around it,
// because the forms a permissive SHOW swallows are exactly the ones no test
// suite writes: a target keyword with nothing after it, a missing FROM|IN, a
// singular where the grammar has only a plural, two filters at once, a LIKE
// whose operand is not a string, a FOR CHANNEL named as an identifier, a
// whole statement after the target, and the other engines' SHOW statements.
//
// The list is the same text in all four engine trees. What differs is the
// engine column, which is where the forks part, so a row that reads
// differently here than in a sibling tree is a dialect delta rather than an
// oversight.
//
// The engine column is what 5.7.25-OceanBase_CE-v4.3.5.6 answered on
// 2026-11-04, read with scripts/oceanbase-oracle.sh. OceanBase needs about
// 8 GB and CI has no room for it, so TestShowTargetsOceanBaseOracle reads that
// column back from a live server only where the DSN is set; TestShowTargets
// holds the parser to the record unconditionally, which is the gate CI runs.
//
// overAccepted marks the 23 rows where this parser takes a statement
// OceanBase answers 1064 on. All of them are MySQL replication, binary-log,
// relay-log, event and profiling statements this tree inherited from mysql
// and OceanBase CE does not have; closing them is a change of its own, and
// the flag is here so the measurement is on record rather than lost.
type showTargetCase struct {
	sql string
	// engine is whether OceanBase CE 4.3.5.6 parses the statement.
	engine bool
	// overAccepted marks a statement this parser takes that OceanBase does
	// not. Removing one means clearing the flag in the same change.
	overAccepted bool
}

var showTargetCases = []showTargetCase{
	// A SHOW whose target keyword does not complete. Every one of these
	// left an empty SHOW behind, which then stood in front of whatever
	// followed it.
	{sql: "SHOW", engine: false},
	{sql: "SHOW FULL", engine: false},
	{sql: "SHOW EXTENDED", engine: false},
	{sql: "SHOW CREATE", engine: false},
	{sql: "SHOW TABLE", engine: false},
	{sql: "SHOW GLOBAL", engine: false},
	{sql: "SHOW SESSION", engine: false},
	{sql: "SHOW MASTER", engine: false},
	{sql: "SHOW SLAVE", engine: false},
	{sql: "SHOW REPLICA", engine: false},
	{sql: "SHOW BINARY", engine: false},
	{sql: "SHOW BINLOG", engine: false},
	{sql: "SHOW RELAYLOG", engine: false},
	{sql: "SHOW STORAGE", engine: false},
	{sql: "SHOW OPEN", engine: false},
	{sql: "SHOW CHARACTER", engine: false},
	{sql: "SHOW ENGINE", engine: false},
	{sql: "SHOW PROCEDURE", engine: false},
	{sql: "SHOW FUNCTION", engine: false},
	{sql: "SHOW COUNT", engine: false},
	{sql: "SHOW COUNT(*)", engine: false},
	{sql: "SHOW COUNT ERRORS", engine: false},
	{sql: "SHOW COUNT(*) ERRORS", engine: true},
	{sql: "SHOW COUNT(*) WARNINGS", engine: true},
	{sql: "SHOW COUNT(*) TABLES", engine: false},

	// The database and table listings.
	{sql: "SHOW DATABASES", engine: true},
	{sql: "SHOW SCHEMAS", engine: true},
	{sql: "SHOW TABLES", engine: true},
	{sql: "SHOW TABLES FROM test", engine: true},
	{sql: "SHOW TABLES IN test", engine: true},
	{sql: "SHOW FULL TABLES", engine: true},
	{sql: "SHOW FULL COLUMNS FROM t", engine: true},
	{sql: "SHOW FULL FIELDS FROM t", engine: true},
	{sql: "SHOW FULL PROCESSLIST", engine: true},
	{sql: "SHOW FULL STATUS", engine: false},

	// The column listings, and the FROM|IN target they require.
	{sql: "SHOW COLUMNS FROM t", engine: true},
	{sql: "SHOW COLUMNS FROM t FROM test", engine: true},
	{sql: "SHOW COLUMNS", engine: false},
	{sql: "SHOW COLUMNS t", engine: false},
	{sql: "SHOW FIELDS FROM t", engine: true},

	// The EXTENDED prefix, and what it may stand in front of.
	{sql: "SHOW EXTENDED COLUMNS FROM t", engine: true},
	{sql: "SHOW EXTENDED FIELDS FROM t", engine: true},
	{sql: "SHOW EXTENDED FULL COLUMNS FROM t", engine: true},
	{sql: "SHOW EXTENDED FULL FIELDS FROM t", engine: true},
	{sql: "SHOW EXTENDED TABLES", engine: true},
	{sql: "SHOW EXTENDED FULL TABLES", engine: true},
	{sql: "SHOW EXTENDED INDEX FROM t", engine: true},
	{sql: "SHOW EXTENDED KEYS FROM t", engine: true},

	// SHOW CREATE, one line per object this engine has one for.
	{sql: "SHOW CREATE TABLE t", engine: true},
	{sql: "SHOW CREATE DATABASE d", engine: true},
	{sql: "SHOW CREATE DATABASE IF NOT EXISTS d", engine: true},
	{sql: "SHOW CREATE DATABASE IF d", engine: false},
	{sql: "SHOW CREATE SCHEMA d", engine: true},
	{sql: "SHOW CREATE VIEW v", engine: true},
	{sql: "SHOW CREATE PROCEDURE pr", engine: true},
	{sql: "SHOW CREATE FUNCTION fn", engine: true},
	{sql: "SHOW CREATE TRIGGER tg", engine: true},
	{sql: "SHOW CREATE EVENT ev", engine: false, overAccepted: true},
	{sql: "SHOW CREATE USER u", engine: true},
	{sql: "SHOW CREATE USER", engine: false},
	{sql: "SHOW CREATE SEQUENCE seq", engine: false},
	{sql: "SHOW CREATE SEQUENCE test.seq", engine: false},
	{sql: "SHOW CREATE PLACEMENT POLICY x", engine: false},
	{sql: "SHOW CREATE PLACEMENT POLICY", engine: false},
	{sql: "SHOW CREATE RESOURCE GROUP rg", engine: false},
	{sql: "SHOW CREATE RESOURCE GROUP default", engine: false},
	{sql: "SHOW CREATE RESOURCE GROUP", engine: false},
	{sql: "SHOW CREATE RESOURCE rg", engine: false},

	// The index listings, whose three spellings are INDEX, INDEXES and
	// KEYS.
	{sql: "SHOW INDEX FROM t", engine: true},
	{sql: "SHOW INDEX IN t IN test", engine: true},
	{sql: "SHOW INDEX", engine: false},
	{sql: "SHOW INDEX t", engine: false},
	{sql: "SHOW KEY FROM t", engine: false},
	{sql: "SHOW KEYS FROM t", engine: true},
	{sql: "SHOW INDEXES FROM t", engine: true},

	// The scope prefix, and what may follow it.
	{sql: "SHOW GLOBAL VARIABLES", engine: true},
	{sql: "SHOW GLOBAL STATUS", engine: true},
	{sql: "SHOW SESSION VARIABLES", engine: true},
	{sql: "SHOW SESSION STATUS", engine: true},
	{sql: "SHOW GLOBAL BINDINGS", engine: false},
	{sql: "SHOW SESSION BINDINGS", engine: false},
	{sql: "SHOW BINDINGS", engine: false},
	{sql: "SHOW BINDING_CACHE STATUS", engine: false},
	{sql: "SHOW BINDING_CACHE", engine: false},

	// Variables, status, warnings and errors.
	{sql: "SHOW VARIABLES", engine: true},
	{sql: "SHOW VARIABLES LIKE 'x%'", engine: true},
	{sql: "SHOW STATUS", engine: true},
	{sql: "SHOW WARNINGS", engine: true},
	{sql: "SHOW WARNINGS LIMIT 10", engine: true},
	{sql: "SHOW WARNINGS LIMIT 5, 10", engine: true},
	{sql: "SHOW ERRORS", engine: true},
	{sql: "SHOW ERRORS LIMIT 10", engine: true},

	// The storage-engine reports.
	{sql: "SHOW ENGINES", engine: true},
	{sql: "SHOW STORAGE ENGINES", engine: true},
	{sql: "SHOW ENGINE INNODB STATUS", engine: true},
	{sql: "SHOW ENGINE INNODB MUTEX", engine: true},
	{sql: "SHOW ENGINE INNODB LOGS", engine: false},
	{sql: "SHOW ENGINE INNODB", engine: false},
	{sql: "SHOW PLUGINS", engine: true},

	// The replication reports, and the FOR CHANNEL whose name is a
	// string literal rather than an identifier.
	{sql: "SHOW MASTER STATUS", engine: false, overAccepted: true},
	{sql: "SHOW MASTER LOGS", engine: false, overAccepted: true},
	{sql: "SHOW BINARY LOGS", engine: false, overAccepted: true},
	{sql: "SHOW BINARY LOG STATUS", engine: false},
	{sql: "SHOW BINARY LOG", engine: false},
	{sql: "SHOW BINLOG STATUS", engine: false},
	{sql: "SHOW SLAVE STATUS", engine: false, overAccepted: true},
	{sql: "SHOW SLAVE HOSTS", engine: false, overAccepted: true},
	{sql: "SHOW REPLICA STATUS", engine: false, overAccepted: true},
	{sql: "SHOW REPLICA HOSTS", engine: false},
	{sql: "SHOW REPLICAS", engine: false, overAccepted: true},
	{sql: "SHOW SLAVE STATUS FOR CHANNEL 'c'", engine: false, overAccepted: true},
	{sql: "SHOW SLAVE STATUS FOR CHANNEL c", engine: false},
	{sql: "SHOW SLAVE STATUS FOR", engine: false},
	{sql: "SHOW SLAVE STATUS FOR CHANNEL", engine: false},
	{sql: "SHOW REPLICA STATUS FOR CHANNEL 'c'", engine: false, overAccepted: true},
	{sql: "SHOW REPLICA STATUS FOR CHANNEL c", engine: false},
	{sql: "SHOW BINLOG EVENT", engine: false},
	{sql: "SHOW BINLOG EVENTS", engine: false, overAccepted: true},
	{sql: "SHOW BINLOG EVENTS IN 'f'", engine: false, overAccepted: true},
	{sql: "SHOW BINLOG EVENTS FROM 4", engine: false, overAccepted: true},
	{sql: "SHOW BINLOG EVENTS LIMIT 5, 10", engine: false, overAccepted: true},
	{sql: "SHOW RELAYLOG EVENT", engine: false},
	{sql: "SHOW RELAYLOG EVENTS", engine: false, overAccepted: true},
	{sql: "SHOW RELAYLOG EVENTS FOR CHANNEL 'c'", engine: false, overAccepted: true},
	{sql: "SHOW RELAYLOG EVENTS FOR CHANNEL c", engine: false},

	// The remaining MySQL targets, whose plurals are the only spellings
	// the grammar has.
	{sql: "SHOW TABLE STATUS", engine: true},
	{sql: "SHOW TABLE STATUS FROM test", engine: true},
	{sql: "SHOW TABLE STATUS LIKE 't'", engine: true},
	{sql: "SHOW TRIGGER", engine: false},
	{sql: "SHOW TRIGGERS", engine: true},
	{sql: "SHOW EVENT", engine: false},
	{sql: "SHOW EVENTS", engine: false, overAccepted: true},
	{sql: "SHOW PROCEDURE STATUS", engine: true},
	{sql: "SHOW PROCEDURE CODE pr", engine: true},
	{sql: "SHOW FUNCTION STATUS", engine: true},
	{sql: "SHOW FUNCTION CODE fn", engine: true},
	{sql: "SHOW OPEN TABLES", engine: true},
	{sql: "SHOW PRIVILEGES", engine: true},
	{sql: "SHOW PROFILES", engine: true},

	// SHOW PROFILE, whose arguments OceanBase does not read at all.
	{sql: "SHOW PROFILE", engine: true},
	{sql: "SHOW PROFILE CPU", engine: false, overAccepted: true},
	{sql: "SHOW PROFILE CPU, MEMORY FOR QUERY 2", engine: false, overAccepted: true},
	{sql: "SHOW PROFILE BLOCK IO", engine: false, overAccepted: true},
	{sql: "SHOW PROFILE FOR QUERY 2", engine: false, overAccepted: true},
	{sql: "SHOW PROFILE CPU FOR QUERY 2 LIMIT 1", engine: false, overAccepted: true},
	{sql: "SHOW PROFILE CPU FOR QUERY 2 LIMIT 1,1", engine: false},
	{sql: "SHOW PROFILE CPU FOR QUERY 2 LIMIT 1 OFFSET 2", engine: false, overAccepted: true},
	{sql: "SHOW PROFILE CPU LIMIT 1,1 OFFSET 2", engine: false},
	{sql: "SHOW PROFILE CPU LIMIT 1,", engine: false},
	{sql: "SHOW PROFILE CPU OFFSET 2", engine: false},

	// Character sets, collations, grants and the process list.
	{sql: "SHOW CHARACTER SET", engine: true},
	{sql: "SHOW CHARSET", engine: true},
	{sql: "SHOW COLLATION", engine: true},
	{sql: "SHOW GRANTS", engine: true},
	{sql: "SHOW GRANTS FOR u", engine: true},
	{sql: "SHOW PROCESSLIST", engine: true},

	// The filter. LIKE and WHERE read different things, and only one
	// filter may follow a target. The operand after LIKE is one string
	// literal here, so the COLLATE suffix TiDB reads is a 1064.
	{sql: "SHOW DATABASES LIKE 't'", engine: true},
	{sql: "SHOW DATABASES LIKE @v", engine: false},
	{sql: "SHOW DATABASES LIKE concat('a','%')", engine: false},
	{sql: "SHOW DATABASES LIKE 1", engine: false},
	{sql: "SHOW DATABASES LIKE -1", engine: false},
	{sql: "SHOW DATABASES LIKE NULL", engine: false},
	{sql: "SHOW DATABASES LIKE ('t')", engine: false},
	{sql: "SHOW DATABASES LIKE 't%' COLLATE utf8mb4_bin", engine: false},
	{sql: "SHOW DATABASES LIKE 't%' COLLATE 'utf8mb4_bin'", engine: false},
	{sql: "SHOW DATABASES LIKE 't%' COLLATE binary", engine: false},
	{sql: "SHOW DATABASES LIKE 't%' COLLATE", engine: false},
	{sql: "SHOW DATABASES LIKE 't%' COLLATE 1", engine: false},
	{sql: "SHOW TABLES LIKE 't%' COLLATE utf8mb4_bin", engine: false},
	{sql: "SHOW DATABASES LIKE 'a' || 'b'", engine: false},
	{sql: "SHOW DATABASES LIKE 'a' + 1", engine: false},
	{sql: "SHOW DATABASES LIKE 'a' -> '$.x'", engine: false},
	{sql: "SHOW DATABASES LIKE", engine: false},
	{sql: "SHOW DATABASES LIKE 't' LIKE 'u'", engine: false},
	{sql: "SHOW DATABASES LIKE 't' AND 1", engine: false},
	{sql: "SHOW DATABASES LIKE 't' WHERE 1=1", engine: false},
	{sql: "SHOW DATABASES WHERE 1=1", engine: true},
	{sql: "SHOW DATABASES WHERE", engine: false},
	{sql: "SHOW DATABASES WHERE 1=1 LIKE 't'", engine: true},
	{sql: "SHOW TABLES LIKE 't' WHERE 1=1", engine: false},

	// TiDB's placement statements.
	{sql: "SHOW PLACEMENT", engine: false},
	{sql: "SHOW PLACEMENT LIKE 'POLICY foo%'", engine: false},
	{sql: "SHOW PLACEMENT WHERE Target='TABLE test.t1'", engine: false},
	{sql: "SHOW PLACEMENT LABELS", engine: false},
	{sql: "SHOW PLACEMENT LABELS LIKE '%zone%'", engine: false},
	{sql: "SHOW PLACEMENT FOR DATABASE db1", engine: false},
	{sql: "SHOW PLACEMENT FOR SCHEMA db1", engine: false},
	{sql: "SHOW PLACEMENT FOR TABLE tb1", engine: false},
	{sql: "SHOW PLACEMENT FOR TABLE db1.tb1", engine: false},
	{sql: "SHOW PLACEMENT FOR TABLE tb1 PARTITION p1", engine: false},
	{sql: "SHOW PLACEMENT FOR", engine: false},
	{sql: "SHOW PLACEMENT FOR DB db1", engine: false},
	{sql: "SHOW PLACEMENT FOR PARTITION p1", engine: false},
	{sql: "SHOW PLACEMENT FOR TABLE tb1 PARTITION", engine: false},
	{sql: "SHOW PLACEMENT FOR TABLE tb1 LIKE 't'", engine: false},
	{sql: "SHOW PLACEMENT FOR DATABASE db1 LIKE 't'", engine: false},

	// TiDB's SHOW TABLE forms, which share a table name and a
	// partition list.
	{sql: "SHOW TABLE t1 REGIONS", engine: false},
	{sql: "SHOW TABLE t1 REGIONS WHERE a=1", engine: false},
	{sql: "SHOW TABLE t1 REGIONS LIKE 'x'", engine: false},
	{sql: "SHOW TABLE test.t1 REGIONS", engine: false},
	{sql: "SHOW TABLE t1 PARTITION (p0) REGIONS", engine: false},
	{sql: "SHOW TABLE t1 PARTITION (p0,p1) REGIONS WHERE a=1", engine: false},
	{sql: "SHOW TABLE t1 PARTITION p0 REGIONS", engine: false},
	{sql: "SHOW TABLE t1 PARTITION () REGIONS", engine: false},
	{sql: "SHOW TABLE t1 INDEX idx1 REGIONS", engine: false},
	{sql: "SHOW TABLE t1 PARTITION (p0) INDEX idx1 REGIONS WHERE a=2", engine: false},
	{sql: "SHOW TABLE t1 INDEX REGIONS", engine: false},
	{sql: "SHOW TABLE t1 INDEX idx1", engine: false},
	{sql: "SHOW TABLE t1 DISTRIBUTIONS", engine: false},
	{sql: "SHOW TABLE t1 DISTRIBUTIONS WHERE a=1", engine: false},
	{sql: "SHOW TABLE t1 PARTITION (p0,p1) DISTRIBUTIONS", engine: false},
	{sql: "SHOW TABLE t1 NEXT_ROW_ID", engine: false},
	{sql: "SHOW TABLE test.t1 NEXT_ROW_ID", engine: false},
	{sql: "SHOW TABLE t1 PARTITION (p0) NEXT_ROW_ID", engine: false},
	{sql: "SHOW TABLE t1 NEXT_ROW_ID WHERE a=1", engine: false},
	{sql: "SHOW TABLE t1", engine: false},
	{sql: "SHOW TABLE t1 REGIONS REGIONS", engine: false},

	// TiDB's distribution jobs.
	{sql: "SHOW DISTRIBUTION JOBS", engine: false},
	{sql: "SHOW DISTRIBUTION JOBS WHERE id > 0", engine: false},
	{sql: "SHOW DISTRIBUTION JOB 1", engine: false},
	{sql: "SHOW DISTRIBUTION JOB 'x'", engine: false},
	{sql: "SHOW DISTRIBUTION JOB", engine: false},
	{sql: "SHOW DISTRIBUTION JOB 1 WHERE id > 0", engine: false},
	{sql: "SHOW DISTRIBUTION JOBS 1", engine: false},
	{sql: "SHOW DISTRIBUTION", engine: false},

	// TiDB's statistics, session-state, binding, backup and config
	// targets.
	{sql: "SHOW SESSION_STATES", engine: false},
	{sql: "SHOW STATS_EXTENDED", engine: false},
	{sql: "SHOW STATS_META", engine: false},
	{sql: "SHOW STATS_META WHERE table_name = 't'", engine: false},
	{sql: "SHOW STATS_META LIKE 't'", engine: false},
	{sql: "SHOW STATS_HISTOGRAMS", engine: false},
	{sql: "SHOW STATS_TOPN", engine: false},
	{sql: "SHOW STATS_BUCKETS", engine: false},
	{sql: "SHOW STATS_HEALTHY", engine: false},
	{sql: "SHOW STATS_LOCKED", engine: false},
	{sql: "SHOW HISTOGRAMS_IN_FLIGHT", engine: false},
	{sql: "SHOW COLUMN_STATS_USAGE", engine: false},
	{sql: "SHOW ANALYZE STATUS", engine: false},
	{sql: "SHOW ANALYZE STATUS WHERE table_name = 't'", engine: false},
	{sql: "SHOW ANALYZE", engine: false},
	{sql: "SHOW BUILTINS", engine: false},
	{sql: "SHOW BUILTINS LIKE 'x'", engine: false},
	{sql: "SHOW BACKUPS", engine: false},
	{sql: "SHOW BACKUPS WHERE start_time > now()", engine: false},
	{sql: "SHOW RESTORES", engine: false},
	{sql: "SHOW RESTORES LIKE 'r0001'", engine: false},
	{sql: "SHOW CONFIG", engine: false},
	{sql: "SHOW CONFIG WHERE type='tidb'", engine: false},

	// A statement after the target. Split has already cut on the
	// delimiter, so what is left in the segment is a second statement, and
	// a SHOW that returned without reading to the end let the loop take it
	// as one.
	{sql: "SHOW DATABASES TABLE t", engine: false},
	{sql: "SHOW TABLES TABLE t", engine: false},
	{sql: "SHOW CREATE TABLE t TABLE t2", engine: false},
	{sql: "SHOW TABLE STATUS TABLE t", engine: false},
	{sql: "SHOW VARIABLES LIKE 'x%' TABLE t", engine: false},
	{sql: "SHOW PLACEMENT TABLE t", engine: false},
	{sql: "SHOW PLACEMENT FOR DATABASE db1 TABLE tb1", engine: false},
	{sql: "SHOW PLACEMENT FOR SCHEMA db1 TABLE tb1", engine: false},
	{sql: "SHOW PLACEMENT FOR TABLE tb1 TABLE tb2", engine: false},
	{sql: "SHOW PLACEMENT FOR TABLE tb1 DATABASE db1", engine: false},
	{sql: "SHOW PLACEMENT FOR TABLE tb1 PARTITION p1 PARTITION p2", engine: false},
	{sql: "SHOW TABLE t1 REGIONS TABLE t2", engine: false},
	{sql: "SHOW TABLE t1 NEXT_ROW_ID TABLE t2", engine: false},
	{sql: "SHOW BUILTINS TABLE t", engine: false},
	{sql: "SHOW DISTRIBUTION JOB 1 TABLE t", engine: false},
	{sql: "SHOW STATS_META TABLE t", engine: false},
}

// TestShowTargets holds the parser to the recorded OceanBase verdicts. It
// needs no server, so it is the gate CI runs.
func TestShowTargets(t *testing.T) {
	for _, c := range showTargetCases {
		t.Run(c.sql, func(t *testing.T) {
			want := c.engine || c.overAccepted
			_, err := Parse(c.sql)
			if got := err == nil; got != want {
				t.Errorf("parser accepts=%v, want %v for %q (OceanBase accepts=%v, overAccepted=%v): %v",
					got, want, c.sql, c.engine, c.overAccepted, err)
			}
		})
	}
}
