package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// Forms of v8.5.5 this parser did not read, and ones it read and v8.5.5 does not, all of them
// measured: a refusal is 1064 or a code raised while parsing (prepareVerdict). Every verdict was
// read from pingcap/tidb:v8.5.5 on 2026-11-05.

type readGapCase struct {
	sql    string
	accept bool
}

var readGapCases = []readGapCase{
	// ILIKE is PredicateExpr's second LIKE, and v8.5.5 reserves the word.
	{"SELECT 'a' ILIKE 'b'", true},
	{"SELECT 'a' NOT ILIKE 'b'", true},
	{"SELECT 'a' ! ILIKE 'b'", false},
	{"SELECT 'a' ILIKE 'b' ESCAPE 'x'", true},
	{"SELECT 'a' ILIKE 'b' ESCAPE 'xy'", false},
	{"SELECT 'a' ILIKE 'b' ESCAPE ''", true},
	{"SELECT 'a' ILIKE 'b' ESCAPE concat('c')", false},
	{"SELECT 'a' ilike 'b' ilike 'c'", false},
	{"SELECT 'a' ILIKE 'b' LIKE 'c'", false},
	{"SELECT 'a' ILIKE 'b' + 1", false},
	{"SELECT 'a' ILIKE ('b')", true},
	{"SELECT 'a' ILIKE", false},
	{"SELECT ILIKE 'a'", false},
	{"SELECT 'a' ILIKE ALL (SELECT 1)", false},
	{"SELECT 'a' ILIKE 'b' IS NULL", true},
	{"SELECT NOT 'a' ILIKE 'b'", true},
	{"SELECT 1 FROM t WHERE a ILIKE '%x%' OR b NOT ILIKE 'y'", true},
	{"SELECT a ILIKE b FROM t", true},
	{"SELECT ilike FROM t", false},
	{"SELECT `ilike` FROM t", true},
	{"SELECT t.ilike FROM t", true},
	{"SELECT ilike(1)", false},
	{"SELECT `ilike`(1)", true},
	{"CREATE TABLE ilike (a INT)", false},
	{"CREATE TABLE `ilike` (a INT)", true},
	{"CREATE TABLE t (ilike INT)", false},
	{"SELECT 1 AS ilike", false},
	{"SELECT 1 ilike", false},
	{"SELECT 1 AS `ilike`", true},
	{"SHOW TABLES ILIKE 'a'", false},
	{"SHOW DATABASES ILIKE 'a'", false},
	// INTERVAL belongs to RANGE and RANGE COLUMNS alone: FIRST and LAST come together, NULL PARTITION comes before MAXVALUE PARTITION, and the server asks for FIRST and LAST only when it runs the statement (8200).
	{"CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1)", true},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN (1) LAST PARTITION LESS THAN (5)", true},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN (1) LAST PARTITION LESS THAN (5) NULL PARTITION", true},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN (1) LAST PARTITION LESS THAN (5) MAXVALUE PARTITION", true},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN (1) LAST PARTITION LESS THAN (5) NULL PARTITION MAXVALUE PARTITION", true},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN (1) LAST PARTITION LESS THAN (5) MAXVALUE PARTITION NULL PARTITION", false},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) NULL PARTITION", true},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) MAXVALUE PARTITION", true},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN (1)", false},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) LAST PARTITION LESS THAN (5)", false},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) LAST PARTITION LESS THAN (5) FIRST PARTITION LESS THAN (1)", false},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN 1 LAST PARTITION LESS THAN 5", false},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN (1,2) LAST PARTITION LESS THAN (5)", false},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN (MAXVALUE) LAST PARTITION LESS THAN (5)", false},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL ()", false},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL 1", false},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL", false},
	{"CREATE TABLE t (a DATETIME) PARTITION BY RANGE (TO_DAYS(a)) INTERVAL (1 DAY) FIRST PARTITION LESS THAN (1) LAST PARTITION LESS THAN (5)", true},
	{"CREATE TABLE t (a DATETIME) PARTITION BY RANGE COLUMNS (a) INTERVAL (1 MONTH) FIRST PARTITION LESS THAN ('2020-01-01') LAST PARTITION LESS THAN ('2021-01-01')", true},
	{"CREATE TABLE t (a DATETIME) PARTITION BY RANGE COLUMNS (a) INTERVAL (1 FOO)", false},
	{"CREATE TABLE t (a DATETIME) PARTITION BY RANGE COLUMNS (a) INTERVAL (1 YEAR_MONTH)", true},
	{"CREATE TABLE t (a DATETIME) PARTITION BY RANGE COLUMNS (a) INTERVAL (1 SQL_TSI_DAY)", true},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE FIELDS (a) INTERVAL (1)", true},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1+1)", true},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1,2)", false},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) PARTITIONS 3", true},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) (PARTITION p VALUES LESS THAN (1))", true},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN (1) LAST PARTITION LESS THAN (5) (PARTITION p VALUES LESS THAN (1))", true},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) INTERVAL (1) SUBPARTITION BY HASH (b)", true},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN (1) LAST PARTITION LESS THAN (5) SUBPARTITION BY HASH (b)", true},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) UPDATE INDEXES (i GLOBAL)", true},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN (1) LAST PARTITION LESS THAN (5) UPDATE INDEXES (i GLOBAL)", true},
	{"CREATE TABLE t (a INT) PARTITION BY LIST (a) INTERVAL (1)", false},
	{"CREATE TABLE t (a INT) PARTITION BY LIST COLUMNS (a) INTERVAL (1)", false},
	{"CREATE TABLE t (a INT) PARTITION BY HASH (a) INTERVAL (1)", false},
	{"CREATE TABLE t (a INT) PARTITION BY KEY (a) INTERVAL (1)", false},
	{"ALTER TABLE t PARTITION BY RANGE (a) INTERVAL (1)", true},
	{"ALTER TABLE t PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN (1) LAST PARTITION LESS THAN (5)", true},
	{"ALTER TABLE t PARTITION BY RANGE COLUMNS (a) INTERVAL (1 MONTH) FIRST PARTITION LESS THAN ('2020-01-01') LAST PARTITION LESS THAN ('2021-01-01') NULL PARTITION MAXVALUE PARTITION", true},
	// SYSTEM_TIME, with HISTORY and CURRENT as the bounds only it takes: PartitionOptions.Validate wants two partitions (4128), 1480 for other bounds, 4113 for HISTORY elsewhere. ADD and REORGANIZE PARTITION have no method to hold them to, so they read them.
	{"CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME", false},
	{"CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x HISTORY)", false},
	{"CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x HISTORY, PARTITION y CURRENT)", true},
	{"CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x CURRENT, PARTITION y HISTORY)", true},
	{"CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x HISTORY, PARTITION y HISTORY, PARTITION z CURRENT)", true},
	{"CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x, PARTITION y CURRENT)", false},
	{"CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x HISTORY, PARTITION y)", false},
	{"CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x VALUES LESS THAN (1), PARTITION y CURRENT)", false},
	{"CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x VALUES IN (1), PARTITION y CURRENT)", false},
	{"CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x DEFAULT, PARTITION y CURRENT)", false},
	{"CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME INTERVAL 1 DAY (PARTITION x HISTORY, PARTITION y CURRENT)", true},
	{"CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME INTERVAL 1+1 DAY (PARTITION x HISTORY, PARTITION y CURRENT)", true},
	{"CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME INTERVAL 1 FOO (PARTITION x HISTORY, PARTITION y CURRENT)", false},
	{"CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME INTERVAL 1 (PARTITION x HISTORY, PARTITION y CURRENT)", false},
	{"CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME INTERVAL (PARTITION x HISTORY, PARTITION y CURRENT)", false},
	{"CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME INTERVAL (1) DAY (PARTITION x HISTORY, PARTITION y CURRENT)", true},
	{"CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME LIMIT 5 (PARTITION x HISTORY, PARTITION y CURRENT)", true},
	{"CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME LIMIT 0 (PARTITION x HISTORY, PARTITION y CURRENT)", true},
	{"CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME LIMIT (PARTITION x HISTORY, PARTITION y CURRENT)", false},
	{"CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME LIMIT 'a' (PARTITION x HISTORY, PARTITION y CURRENT)", false},
	{"CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME LIMIT 5, 6 (PARTITION x HISTORY, PARTITION y CURRENT)", false},
	{"CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME PARTITIONS 2 (PARTITION x HISTORY, PARTITION y CURRENT)", true},
	{"CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME PARTITIONS 3 (PARTITION x HISTORY, PARTITION y CURRENT)", false},
	{"CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME PARTITIONS 2", false},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY SYSTEM_TIME SUBPARTITION BY HASH (b) (PARTITION x HISTORY, PARTITION y CURRENT)", true},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY SYSTEM_TIME (PARTITION x HISTORY (SUBPARTITION s), PARTITION y CURRENT (SUBPARTITION t))", false},
	{"CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x HISTORY COMMENT 'c', PARTITION y CURRENT)", true},
	{"CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x HISTORY, PARTITION y CURRENT) UPDATE INDEXES (i GLOBAL)", true},
	{"CREATE TABLE t (a INT) PARTITION BY LINEAR SYSTEM_TIME (PARTITION x HISTORY, PARTITION y CURRENT)", false},
	{"ALTER TABLE t PARTITION BY SYSTEM_TIME (PARTITION x HISTORY, PARTITION y CURRENT)", true},
	{"ALTER TABLE t PARTITION BY SYSTEM_TIME LIMIT 5 (PARTITION x HISTORY, PARTITION y CURRENT)", true},
	{"ALTER TABLE t ADD PARTITION (PARTITION x HISTORY)", true},
	{"ALTER TABLE t ADD PARTITION (PARTITION x CURRENT)", true},
	{"ALTER TABLE t REORGANIZE PARTITION a INTO (PARTITION x HISTORY)", true},
	{"CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION x HISTORY)", false},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION x CURRENT)", false},
	{"CREATE TABLE t (a INT) PARTITION BY LIST (a) (PARTITION x HISTORY)", false},
	{"CREATE TABLE t (a INT) PARTITION BY KEY (a) (PARTITION x CURRENT)", false},
	// v8.5.5 has no CREATE, ALTER or DROP for TABLESPACE, LOGFILE GROUP or SERVER, and no FLUSH ... FOR EXPORT, though it keeps a table's TABLESPACE option and the privilege's name.
	{"CREATE TABLESPACE ts ADD DATAFILE 'a'", false},
	{"CREATE TABLESPACE ts", false},
	{"CREATE UNDO TABLESPACE ts ADD DATAFILE 'a.ibu'", false},
	{"ALTER TABLESPACE ts ADD DATAFILE 'a'", false},
	{"ALTER TABLESPACE ts RENAME TO b", false},
	{"DROP TABLESPACE ts", false},
	{"DROP UNDO TABLESPACE ts", false},
	{"CREATE LOGFILE GROUP lg ADD UNDOFILE 'u'", false},
	{"ALTER LOGFILE GROUP lg ADD UNDOFILE 'u'", false},
	{"DROP LOGFILE GROUP lg", false},
	{"CREATE SERVER s FOREIGN DATA WRAPPER mysql OPTIONS (USER 'r')", false},
	{"ALTER SERVER s OPTIONS (USER 'r')", false},
	{"DROP SERVER s", false},
	{"CREATE TABLE t (a INT) TABLESPACE ts", true},
	{"CREATE TABLE t (a INT) TABLESPACE = ts", true},
	{"CREATE TABLE t (a INT) TABLESPACE ts STORAGE DISK", true},
	{"CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p TABLESPACE ts)", true},
	{"ALTER TABLE t TABLESPACE ts", true},
	{"ALTER TABLE t DISCARD TABLESPACE", true},
	{"ALTER TABLE t IMPORT TABLESPACE", true},
	{"GRANT CREATE TABLESPACE ON *.* TO u", true},
	{"GRANT CREATE TABLESPACE ON d.* TO u", true},
	{"SHOW CREATE TABLESPACE ts", false},
	{"FLUSH TABLES t FOR EXPORT", false},
	{"SELECT tablespace FROM t", true},
	{"CREATE TABLE tablespace (a INT)", true},
	{"SELECT * FROM information_schema.tablespaces", true},
}

// TestReadGapFormsOracle holds each case to v8.5.5 and to omni.
func TestReadGapFormsOracle(t *testing.T) {
	tc := startTiDB(t)
	assertTiDBVersion(t, tc)
	for i, c := range readGapCases {
		t.Run(c.sql, func(t *testing.T) {
			previous := "<start of cases>"
			if i > 0 {
				previous = readGapCases[i-1].sql
			}
			accepts, _, err := tc.syntaxVerdict(t, c.sql, previous)
			if err != nil {
				t.Fatalf("no verdict: %v", err)
			}
			if accepts != c.accept {
				t.Fatalf("TiDB %s accepts=%v, want %v for %q", tidbCorpusVersion, accepts, c.accept, c.sql)
			}
		})
	}
}

// TestReadGapForms holds the parser to the recorded verdicts without the server.
func TestReadGapForms(t *testing.T) {
	for _, c := range readGapCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := parsesAsOneStatement(c.sql)
			if got != c.accept {
				t.Errorf("omni accepts=%v, the record says %v for %q: %v", got, c.accept, c.sql, err)
			}
		})
	}
}

// TestReadGapTree: the clauses reach the tree.
func TestReadGapTree(t *testing.T) {
	selectExpr := func(sql string) *nodes.LikeExpr {
		t.Helper()
		list, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", sql, err)
		}
		like, ok := list.Items[0].(*nodes.SelectStmt).TargetList[0].(*nodes.LikeExpr)
		if !ok {
			t.Fatalf("%q: no LIKE at the top", sql)
		}
		return like
	}
	if like := selectExpr("SELECT 'a' NOT ILIKE 'b' ESCAPE 'x'"); !like.ILike || !like.Not || like.Escape == nil {
		t.Errorf("NOT ILIKE ... ESCAPE: %+v", like)
	}
	if like := selectExpr("SELECT 'a' LIKE 'b'"); like.ILike {
		t.Errorf("LIKE read as ILIKE: %+v", like)
	}
	partitions := func(sql string) *nodes.PartitionClause {
		t.Helper()
		list, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", sql, err)
		}
		return list.Items[0].(*nodes.CreateTableStmt).Partitions
	}
	iv := partitions("CREATE TABLE t (a DATE) PARTITION BY RANGE COLUMNS (a) INTERVAL (1 SQL_TSI_DAY) FIRST PARTITION LESS THAN ('2020-01-01') LAST PARTITION LESS THAN ('2020-01-03') NULL PARTITION MAXVALUE PARTITION").Interval
	if iv == nil || iv.Unit != "DAY" || iv.First == nil || iv.Last == nil || !iv.NullPart || !iv.MaxValPart {
		t.Errorf("the interval: %+v", iv)
	}
	if iv := partitions("CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1)").Interval; iv == nil || iv.Unit != "" || iv.First != nil || iv.NullPart || iv.MaxValPart {
		t.Errorf("the bare interval: %+v", iv)
	}
	byTime := partitions("CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME INTERVAL 1 DAY (PARTITION x HISTORY, PARTITION y CURRENT)")
	if byTime.Type != nodes.PartitionSystemTime || byTime.SystemTimeInterval == nil || byTime.SystemTimeUnit != "DAY" || len(byTime.Partitions) != 2 {
		t.Errorf("SYSTEM_TIME INTERVAL: %+v", byTime)
	}
	if v, ok := byTime.Partitions[0].Values.(*nodes.String); !ok || v.Str != "HISTORY" {
		t.Errorf("HISTORY: %+v", byTime.Partitions[0].Values)
	}
	if v, ok := byTime.Partitions[1].Values.(*nodes.String); !ok || v.Str != "CURRENT" {
		t.Errorf("CURRENT: %+v", byTime.Partitions[1].Values)
	}
	limited := partitions("CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME LIMIT 5 (PARTITION x HISTORY, PARTITION y CURRENT)")
	if !limited.HasSystemTimeLimit || limited.SystemTimeLimit != 5 {
		t.Errorf("SYSTEM_TIME LIMIT: %+v", limited)
	}
}
