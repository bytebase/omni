package parser

import "testing"

// The SHOW target table, held to a live MariaDB 11.8.8. One case per target
// this parser reads and per malformed variant around it, because the forms a
// permissive SHOW swallows are exactly the ones no test suite writes: a
// target keyword with nothing after it, a missing FROM|IN, a singular where
// the grammar has only a plural, two filters at once, a LIKE whose operand is
// not a string, a FOR CHANNEL named as an identifier, a whole statement
// after the target, and the other engines' SHOW statements.
//
// The list is the same text in all four engine trees. What differs is this
// column, which is where the forks part, so a row that reads differently here
// than in a sibling tree is a dialect delta rather than an oversight.
//
// Measured on 2026-11-04. The expectation is what the server answered, so an
// image upgrade that changes a verdict fails here instead of passing
// silently.

// showTargetCase is a statement and whether the engine's parser accepts it.
type showTargetCase struct {
	sql    string
	accept bool
}

var showTargetCases = []showTargetCase{
	// A SHOW whose target keyword does not complete. Every one of these
	// left an empty SHOW behind, which then stood in front of whatever
	// followed it.
	{"SHOW", false},
	{"SHOW FULL", false},
	{"SHOW EXTENDED", false},
	{"SHOW CREATE", false},
	{"SHOW TABLE", false},
	{"SHOW GLOBAL", false},
	{"SHOW SESSION", false},
	{"SHOW MASTER", false},
	{"SHOW SLAVE", false},
	{"SHOW REPLICA", false},
	{"SHOW BINARY", false},
	{"SHOW BINLOG", false},
	{"SHOW RELAYLOG", false},
	{"SHOW STORAGE", false},
	{"SHOW OPEN", false},
	{"SHOW CHARACTER", false},
	{"SHOW ENGINE", false},
	{"SHOW PROCEDURE", false},
	{"SHOW FUNCTION", false},
	{"SHOW COUNT", false},
	{"SHOW COUNT(*)", false},
	{"SHOW COUNT ERRORS", false},
	{"SHOW COUNT(*) ERRORS", true},
	{"SHOW COUNT(*) WARNINGS", true},
	{"SHOW COUNT(*) TABLES", false},

	// The database and table listings.
	{"SHOW DATABASES", true},
	{"SHOW SCHEMAS", true},
	{"SHOW TABLES", true},
	{"SHOW TABLES FROM test", true},
	{"SHOW TABLES IN test", true},
	{"SHOW FULL TABLES", true},
	{"SHOW FULL COLUMNS FROM t", true},
	{"SHOW FULL FIELDS FROM t", true},
	{"SHOW FULL PROCESSLIST", true},
	{"SHOW FULL STATUS", false},

	// The column listings, and the FROM|IN target they require.
	{"SHOW COLUMNS FROM t", true},
	{"SHOW COLUMNS FROM t FROM test", true},
	{"SHOW COLUMNS", false},
	{"SHOW COLUMNS t", false},
	{"SHOW FIELDS FROM t", true},

	// The EXTENDED prefix, and what it may stand in front of.
	{"SHOW EXTENDED COLUMNS FROM t", false},
	{"SHOW EXTENDED FIELDS FROM t", false},
	{"SHOW EXTENDED FULL COLUMNS FROM t", false},
	{"SHOW EXTENDED FULL FIELDS FROM t", false},
	{"SHOW EXTENDED TABLES", false},
	{"SHOW EXTENDED FULL TABLES", false},
	{"SHOW EXTENDED INDEX FROM t", false},
	{"SHOW EXTENDED KEYS FROM t", false},

	// SHOW CREATE, one line per object this engine has one for.
	{"SHOW CREATE TABLE t", true},
	{"SHOW CREATE DATABASE d", true},
	{"SHOW CREATE DATABASE IF NOT EXISTS d", true},
	{"SHOW CREATE DATABASE IF d", false},
	{"SHOW CREATE SCHEMA d", true},
	{"SHOW CREATE VIEW v", true},
	{"SHOW CREATE PROCEDURE pr", true},
	{"SHOW CREATE FUNCTION fn", true},
	{"SHOW CREATE TRIGGER tg", true},
	{"SHOW CREATE EVENT ev", true},
	{"SHOW CREATE USER u", true},
	{"SHOW CREATE USER", true},
	{"SHOW CREATE SEQUENCE seq", true},
	{"SHOW CREATE SEQUENCE test.seq", true},
	{"SHOW CREATE PLACEMENT POLICY x", false},
	{"SHOW CREATE PLACEMENT POLICY", false},
	{"SHOW CREATE RESOURCE GROUP rg", false},
	{"SHOW CREATE RESOURCE GROUP default", false},
	{"SHOW CREATE RESOURCE GROUP", false},
	{"SHOW CREATE RESOURCE rg", false},

	// The index listings, whose three spellings are INDEX, INDEXES and
	// KEYS.
	{"SHOW INDEX FROM t", true},
	{"SHOW INDEX IN t IN test", true},
	{"SHOW INDEX", false},
	{"SHOW INDEX t", false},
	{"SHOW KEY FROM t", false},
	{"SHOW KEYS FROM t", true},
	{"SHOW INDEXES FROM t", true},

	// The scope prefix, and what may follow it.
	{"SHOW GLOBAL VARIABLES", true},
	{"SHOW GLOBAL STATUS", true},
	{"SHOW SESSION VARIABLES", true},
	{"SHOW SESSION STATUS", true},
	{"SHOW GLOBAL BINDINGS", false},
	{"SHOW SESSION BINDINGS", false},
	{"SHOW BINDINGS", false},
	{"SHOW BINDING_CACHE STATUS", false},
	{"SHOW BINDING_CACHE", false},

	// Variables, status, warnings and errors.
	{"SHOW VARIABLES", true},
	{"SHOW VARIABLES LIKE 'x%'", true},
	{"SHOW STATUS", true},
	{"SHOW WARNINGS", true},
	{"SHOW WARNINGS LIMIT 10", true},
	{"SHOW WARNINGS LIMIT 5, 10", true},
	{"SHOW ERRORS", true},
	{"SHOW ERRORS LIMIT 10", true},

	// The storage-engine reports.
	{"SHOW ENGINES", true},
	{"SHOW STORAGE ENGINES", true},
	{"SHOW ENGINE INNODB STATUS", true},
	{"SHOW ENGINE INNODB MUTEX", true},
	{"SHOW ENGINE INNODB LOGS", true},
	{"SHOW ENGINE INNODB", false},
	{"SHOW PLUGINS", true},

	// The replication reports, and the FOR CHANNEL whose name is a
	// string literal rather than an identifier.
	{"SHOW MASTER STATUS", true},
	{"SHOW MASTER LOGS", true},
	{"SHOW BINARY LOGS", true},
	{"SHOW BINARY LOG STATUS", false},
	{"SHOW BINARY LOG", false},
	{"SHOW BINLOG STATUS", true},
	{"SHOW SLAVE STATUS", true},
	{"SHOW SLAVE HOSTS", true},
	{"SHOW REPLICA STATUS", true},
	{"SHOW REPLICA HOSTS", true},
	{"SHOW REPLICAS", false},
	{"SHOW SLAVE STATUS FOR CHANNEL 'c'", true},
	{"SHOW SLAVE STATUS FOR CHANNEL c", false},
	{"SHOW SLAVE STATUS FOR", false},
	{"SHOW SLAVE STATUS FOR CHANNEL", false},
	{"SHOW REPLICA STATUS FOR CHANNEL 'c'", true},
	{"SHOW REPLICA STATUS FOR CHANNEL c", false},
	{"SHOW BINLOG EVENT", false},
	{"SHOW BINLOG EVENTS", true},
	{"SHOW BINLOG EVENTS IN 'f'", true},
	{"SHOW BINLOG EVENTS FROM 4", true},
	{"SHOW BINLOG EVENTS LIMIT 5, 10", true},
	{"SHOW RELAYLOG EVENT", false},
	{"SHOW RELAYLOG EVENTS", true},
	{"SHOW RELAYLOG EVENTS FOR CHANNEL 'c'", true},
	{"SHOW RELAYLOG EVENTS FOR CHANNEL c", false},

	// The remaining MySQL targets, whose plurals are the only spellings
	// the grammar has.
	{"SHOW TABLE STATUS", true},
	{"SHOW TABLE STATUS FROM test", true},
	{"SHOW TABLE STATUS LIKE 't'", true},
	{"SHOW TRIGGER", false},
	{"SHOW TRIGGERS", true},
	{"SHOW EVENT", false},
	{"SHOW EVENTS", true},
	{"SHOW PROCEDURE STATUS", true},
	{"SHOW PROCEDURE CODE pr", true},
	{"SHOW FUNCTION STATUS", true},
	{"SHOW FUNCTION CODE fn", true},
	{"SHOW OPEN TABLES", true},
	{"SHOW PRIVILEGES", true},
	{"SHOW PROFILES", true},

	// SHOW PROFILE, whose limit clause reads the comma form as well as
	// the OFFSET one, and not both at once.
	{"SHOW PROFILE", true},
	{"SHOW PROFILE CPU", true},
	{"SHOW PROFILE CPU, MEMORY FOR QUERY 2", true},
	{"SHOW PROFILE BLOCK IO", true},
	{"SHOW PROFILE FOR QUERY 2", true},
	{"SHOW PROFILE CPU FOR QUERY 2 LIMIT 1", true},
	{"SHOW PROFILE CPU FOR QUERY 2 LIMIT 1,1", true},
	{"SHOW PROFILE CPU FOR QUERY 2 LIMIT 1 OFFSET 2", true},
	{"SHOW PROFILE CPU LIMIT 1,1 OFFSET 2", false},
	{"SHOW PROFILE CPU LIMIT 1,", false},
	{"SHOW PROFILE CPU OFFSET 2", false},

	// Character sets, collations, grants and the process list.
	{"SHOW CHARACTER SET", true},
	{"SHOW CHARSET", true},
	{"SHOW COLLATION", true},
	{"SHOW GRANTS", true},
	{"SHOW GRANTS FOR u", true},
	{"SHOW PROCESSLIST", true},

	// The filter. LIKE and WHERE read different things, and only one
	// filter may follow a target. TiDB reads a SimpleExpr after LIKE, so a
	// COLLATE suffix is part of the pattern there and a 1064 on the other
	// three, where the operand is one string literal.
	{"SHOW DATABASES LIKE 't'", true},
	{"SHOW DATABASES LIKE @v", false},
	{"SHOW DATABASES LIKE concat('a','%')", false},
	{"SHOW DATABASES LIKE 1", false},
	{"SHOW DATABASES LIKE -1", false},
	{"SHOW DATABASES LIKE NULL", false},
	{"SHOW DATABASES LIKE ('t')", false},
	{"SHOW DATABASES LIKE 't%' COLLATE utf8mb4_bin", false},
	{"SHOW DATABASES LIKE 't%' COLLATE 'utf8mb4_bin'", false},
	{"SHOW DATABASES LIKE 't%' COLLATE binary", false},
	{"SHOW DATABASES LIKE 't%' COLLATE", false},
	{"SHOW DATABASES LIKE 't%' COLLATE 1", false},
	{"SHOW TABLES LIKE 't%' COLLATE utf8mb4_bin", false},
	{"SHOW DATABASES LIKE 'a' || 'b'", false},
	{"SHOW DATABASES LIKE 'a' + 1", false},
	{"SHOW DATABASES LIKE 'a' -> '$.x'", false},
	{"SHOW DATABASES LIKE", false},
	{"SHOW DATABASES LIKE 't' LIKE 'u'", false},
	{"SHOW DATABASES LIKE 't' AND 1", false},
	{"SHOW DATABASES LIKE 't' WHERE 1=1", false},
	{"SHOW DATABASES WHERE 1=1", true},
	{"SHOW DATABASES WHERE", false},
	{"SHOW DATABASES WHERE 1=1 LIKE 't'", true},
	{"SHOW TABLES LIKE 't' WHERE 1=1", false},

	// TiDB's placement statements.
	{"SHOW PLACEMENT", false},
	{"SHOW PLACEMENT LIKE 'POLICY foo%'", false},
	{"SHOW PLACEMENT WHERE Target='TABLE test.t1'", false},
	{"SHOW PLACEMENT LABELS", false},
	{"SHOW PLACEMENT LABELS LIKE '%zone%'", false},
	{"SHOW PLACEMENT FOR DATABASE db1", false},
	{"SHOW PLACEMENT FOR SCHEMA db1", false},
	{"SHOW PLACEMENT FOR TABLE tb1", false},
	{"SHOW PLACEMENT FOR TABLE db1.tb1", false},
	{"SHOW PLACEMENT FOR TABLE tb1 PARTITION p1", false},
	{"SHOW PLACEMENT FOR", false},
	{"SHOW PLACEMENT FOR DB db1", false},
	{"SHOW PLACEMENT FOR PARTITION p1", false},
	{"SHOW PLACEMENT FOR TABLE tb1 PARTITION", false},
	{"SHOW PLACEMENT FOR TABLE tb1 LIKE 't'", false},
	{"SHOW PLACEMENT FOR DATABASE db1 LIKE 't'", false},

	// TiDB's SHOW TABLE forms, which share a table name and a
	// partition list.
	{"SHOW TABLE t1 REGIONS", false},
	{"SHOW TABLE t1 REGIONS WHERE a=1", false},
	{"SHOW TABLE t1 REGIONS LIKE 'x'", false},
	{"SHOW TABLE test.t1 REGIONS", false},
	{"SHOW TABLE t1 PARTITION (p0) REGIONS", false},
	{"SHOW TABLE t1 PARTITION (p0,p1) REGIONS WHERE a=1", false},
	{"SHOW TABLE t1 PARTITION p0 REGIONS", false},
	{"SHOW TABLE t1 PARTITION () REGIONS", false},
	{"SHOW TABLE t1 INDEX idx1 REGIONS", false},
	{"SHOW TABLE t1 PARTITION (p0) INDEX idx1 REGIONS WHERE a=2", false},
	{"SHOW TABLE t1 INDEX REGIONS", false},
	{"SHOW TABLE t1 INDEX idx1", false},
	{"SHOW TABLE t1 DISTRIBUTIONS", false},
	{"SHOW TABLE t1 DISTRIBUTIONS WHERE a=1", false},
	{"SHOW TABLE t1 PARTITION (p0,p1) DISTRIBUTIONS", false},
	{"SHOW TABLE t1 NEXT_ROW_ID", false},
	{"SHOW TABLE test.t1 NEXT_ROW_ID", false},
	{"SHOW TABLE t1 PARTITION (p0) NEXT_ROW_ID", false},
	{"SHOW TABLE t1 NEXT_ROW_ID WHERE a=1", false},
	{"SHOW TABLE t1", false},
	{"SHOW TABLE t1 REGIONS REGIONS", false},

	// TiDB's distribution jobs.
	{"SHOW DISTRIBUTION JOBS", false},
	{"SHOW DISTRIBUTION JOBS WHERE id > 0", false},
	{"SHOW DISTRIBUTION JOB 1", false},
	{"SHOW DISTRIBUTION JOB 'x'", false},
	{"SHOW DISTRIBUTION JOB", false},
	{"SHOW DISTRIBUTION JOB 1 WHERE id > 0", false},
	{"SHOW DISTRIBUTION JOBS 1", false},
	{"SHOW DISTRIBUTION", false},

	// TiDB's statistics, session-state, binding, backup and config
	// targets.
	{"SHOW SESSION_STATES", false},
	{"SHOW STATS_EXTENDED", false},
	{"SHOW STATS_META", false},
	{"SHOW STATS_META WHERE table_name = 't'", false},
	{"SHOW STATS_META LIKE 't'", false},
	{"SHOW STATS_HISTOGRAMS", false},
	{"SHOW STATS_TOPN", false},
	{"SHOW STATS_BUCKETS", false},
	{"SHOW STATS_HEALTHY", false},
	{"SHOW STATS_LOCKED", false},
	{"SHOW HISTOGRAMS_IN_FLIGHT", false},
	{"SHOW COLUMN_STATS_USAGE", false},
	{"SHOW ANALYZE STATUS", false},
	{"SHOW ANALYZE STATUS WHERE table_name = 't'", false},
	{"SHOW ANALYZE", false},
	{"SHOW BUILTINS", false},
	{"SHOW BUILTINS LIKE 'x'", false},
	{"SHOW BACKUPS", false},
	{"SHOW BACKUPS WHERE start_time > now()", false},
	{"SHOW RESTORES", false},
	{"SHOW RESTORES LIKE 'r0001'", false},
	{"SHOW CONFIG", false},
	{"SHOW CONFIG WHERE type='tidb'", false},

	// A statement after the target. Split has already cut on the
	// delimiter, so what is left in the segment is a second statement, and
	// a SHOW that returned without reading to the end let the loop take it
	// as one.
	{"SHOW DATABASES TABLE t", false},
	{"SHOW TABLES TABLE t", false},
	{"SHOW CREATE TABLE t TABLE t2", false},
	{"SHOW TABLE STATUS TABLE t", false},
	{"SHOW VARIABLES LIKE 'x%' TABLE t", false},
	{"SHOW PLACEMENT TABLE t", false},
	{"SHOW PLACEMENT FOR DATABASE db1 TABLE tb1", false},
	{"SHOW PLACEMENT FOR SCHEMA db1 TABLE tb1", false},
	{"SHOW PLACEMENT FOR TABLE tb1 TABLE tb2", false},
	{"SHOW PLACEMENT FOR TABLE tb1 DATABASE db1", false},
	{"SHOW PLACEMENT FOR TABLE tb1 PARTITION p1 PARTITION p2", false},
	{"SHOW TABLE t1 REGIONS TABLE t2", false},
	{"SHOW TABLE t1 NEXT_ROW_ID TABLE t2", false},
	{"SHOW BUILTINS TABLE t", false},
	{"SHOW DISTRIBUTION JOB 1 TABLE t", false},
	{"SHOW STATS_META TABLE t", false},
}

// TestShowTargetsOracle asserts omni agrees with a live MariaDB 11.8.8 on
// showTargetCases (container test). The verdict comes from PREPARE, which
// parses the statement without running it: 1064 is the parser's rejection,
// and any other server error means the text parsed and failed later.
func TestShowTargetsOracle(t *testing.T) {
	o := startMariaDB(t)
	for _, tc := range showTargetCases {
		t.Run(tc.sql, func(t *testing.T) {
			if got := o.prepareAccepts(t, tc.sql); got != tc.accept {
				t.Fatalf("MariaDB 11.8.8 accepts=%v, want %v for %q", got, tc.accept, tc.sql)
			}
			_, err := Parse(tc.sql)
			if omniAccepts := err == nil; omniAccepts != tc.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, MariaDB 11.8.8 accepts=%v for %q: %v",
					omniAccepts, tc.accept, tc.sql, err)
			}
		})
	}
}

// TestShowTargets holds the parser to the recorded verdicts and needs no
// server, so it is the gate that runs without Docker.
func TestShowTargets(t *testing.T) {
	for _, tc := range showTargetCases {
		t.Run(tc.sql, func(t *testing.T) {
			_, err := Parse(tc.sql)
			if got := err == nil; got != tc.accept {
				t.Errorf("parser accepts=%v, want %v for %q: %v", got, tc.accept, tc.sql, err)
			}
		})
	}
}
