package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/tidb/ast"
)

// What these statements leave in the tree, which a parse-only verdict cannot
// see. The corpus records only whether a statement parses, and the shapes
// here are the ones that would read alike if it did not: `RECOVER TABLE t1
// 100` must not hold 100 as a second table name, `FLASHBACK TABLE t TO t1`
// must keep t1, and the rename form and the point-in-time form must not land
// in the same node.

// TestBRIEInTree asserts the kind, the target and every operand of a BACKUP
// or RESTORE statement.
func TestBRIEInTree(t *testing.T) {
	for _, tc := range []struct {
		sql     string
		kind    ast.BRIEKind
		target  ast.BRIETarget
		schemas []string
		tables  []string
		storage string
	}{
		{
			sql:     "BACKUP DATABASE * TO 'noop://'",
			kind:    ast.BRIEBackup,
			target:  ast.BRIETargetAllSchemas,
			storage: "noop://",
		},
		{
			sql:     "BACKUP SCHEMA * TO 'noop://'",
			kind:    ast.BRIEBackup,
			target:  ast.BRIETargetAllSchemas,
			storage: "noop://",
		},
		{
			// `DATABASE `*`` is a one-element list naming a database
			// called `*`, not the every-database form.
			sql:     "BACKUP DATABASE `*` TO 'noop://'",
			kind:    ast.BRIEBackup,
			target:  ast.BRIETargetSchemas,
			schemas: []string{"*"},
			storage: "noop://",
		},
		{
			sql:     "BACKUP DATABASE a, b, c TO 'noop://'",
			kind:    ast.BRIEBackup,
			target:  ast.BRIETargetSchemas,
			schemas: []string{"a", "b", "c"},
			storage: "noop://",
		},
		{
			sql:     "BACKUP TABLE a.b, c.d, e TO 'local:///tmp/archive01/'",
			kind:    ast.BRIEBackup,
			target:  ast.BRIETargetTables,
			tables:  []string{"a.b", "c.d", "e"},
			storage: "local:///tmp/archive01/",
		},
		{
			sql:     "BACKUP LOGS TO 'noop://'",
			kind:    ast.BRIEBackupLogs,
			target:  ast.BRIETargetNone,
			storage: "noop://",
		},
		{
			sql:     "RESTORE DATABASE * FROM 's3://bucket/path/'",
			kind:    ast.BRIERestore,
			target:  ast.BRIETargetAllSchemas,
			storage: "s3://bucket/path/",
		},
		{
			sql:     "RESTORE SCHEMA a FROM 'noop://'",
			kind:    ast.BRIERestore,
			target:  ast.BRIETargetSchemas,
			schemas: []string{"a"},
			storage: "noop://",
		},
		{
			sql:     "RESTORE TABLE a FROM 'noop://'",
			kind:    ast.BRIERestore,
			target:  ast.BRIETargetTables,
			tables:  []string{"a"},
			storage: "noop://",
		},
		{
			sql:     "RESTORE POINT FROM 'noop://log_backup'",
			kind:    ast.BRIERestorePoint,
			target:  ast.BRIETargetNone,
			storage: "noop://log_backup",
		},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			stmt, ok := parseOne(t, tc.sql).(*ast.BRIEStmt)
			if !ok {
				t.Fatalf("Parse(%q) read a %T, want a *ast.BRIEStmt", tc.sql, parseOne(t, tc.sql))
			}
			if stmt.Kind != tc.kind {
				t.Errorf("Kind = %d, want %d", stmt.Kind, tc.kind)
			}
			if stmt.Target != tc.target {
				t.Errorf("Target = %d, want %d", stmt.Target, tc.target)
			}
			if got := strings.Join(stmt.Schemas, ","); got != strings.Join(tc.schemas, ",") {
				t.Errorf("Schemas = %q, want %q", stmt.Schemas, tc.schemas)
			}
			names := make([]string, 0, len(stmt.Tables))
			for _, ref := range stmt.Tables {
				names = append(names, qualified(ref))
			}
			if got := strings.Join(names, ","); got != strings.Join(tc.tables, ",") {
				t.Errorf("Tables = %q, want %q", names, tc.tables)
			}
			if stmt.Storage != tc.storage {
				t.Errorf("Storage = %q, want %q", stmt.Storage, tc.storage)
			}
			if len(stmt.Options) != 0 {
				t.Errorf("Options = %d, want none", len(stmt.Options))
			}
			// The statement's Loc covers all of it, and the storage
			// literal's covers the literal with its quotes.
			if stmt.Loc.Start != 0 || stmt.Loc.End != len(tc.sql) {
				t.Errorf("Loc = %v, want {0 %d}", stmt.Loc, len(tc.sql))
			}
			quoted := "'" + tc.storage + "'"
			if got := tc.sql[stmt.StorageLoc.Start:stmt.StorageLoc.End]; got != quoted {
				t.Errorf("StorageLoc spans %q, want %q", got, quoted)
			}
		})
	}
}

// TestBRIEOptionsInTree asserts each option's name, the form its value was
// written in, and the value itself. One permissive reader would take every
// combination and keep none of it apart, so the value form is asserted as
// much as the value: SNAPSHOT alone is written three ways.
func TestBRIEOptionsInTree(t *testing.T) {
	type want struct {
		name  string
		value ast.BRIEOptionValue
		str   string
		num   uint64
		unit  string
	}
	for _, tc := range []struct {
		sql  string
		opts []want
	}{
		{
			sql:  "BACKUP TABLE a TO 'noop://' concurrency 4",
			opts: []want{{name: "CONCURRENCY", value: ast.BRIEValueInt, num: 4}},
		},
		{
			sql:  "BACKUP TABLE a TO 'noop://' CONCURRENCY = 4",
			opts: []want{{name: "CONCURRENCY", value: ast.BRIEValueInt, num: 4}},
		},
		{
			// A repeat is kept twice rather than folded.
			sql: "BACKUP TABLE a TO 'noop://' concurrency 4 concurrency 8",
			opts: []want{
				{name: "CONCURRENCY", value: ast.BRIEValueInt, num: 4},
				{name: "CONCURRENCY", value: ast.BRIEValueInt, num: 8},
			},
		},
		{
			sql: "BACKUP TABLE a TO 'noop://' checksum_concurrency 4 compression_level 4 ignore_stats 1 compression_type 'lz4'",
			opts: []want{
				{name: "CHECKSUM_CONCURRENCY", value: ast.BRIEValueInt, num: 4},
				{name: "COMPRESSION_LEVEL", value: ast.BRIEValueInt, num: 4},
				{name: "IGNORE_STATS", value: ast.BRIEValueInt, num: 1},
				{name: "COMPRESSION_TYPE", value: ast.BRIEValueString, str: "lz4"},
			},
		},
		{
			// TRUE and FALSE come back as 1 and 0, and a number is
			// kept as written: 2 is not folded to 1.
			sql: "BACKUP TABLE a TO 'noop://' online true checkpoint FALSE load_stats 2",
			opts: []want{
				{name: "ONLINE", value: ast.BRIEValueInt, num: 1},
				{name: "CHECKPOINT", value: ast.BRIEValueInt, num: 0},
				{name: "LOAD_STATS", value: ast.BRIEValueInt, num: 2},
			},
		},
		{
			sql: "RESTORE TABLE a FROM 'noop://' checksum off analyze required",
			opts: []want{
				{name: "CHECKSUM", value: ast.BRIEValueLevel, str: "OFF"},
				{name: "ANALYZE", value: ast.BRIEValueLevel, str: "REQUIRED"},
			},
		},
		{
			// CHECKSUM reads a boolean and a level into different
			// forms, so `checksum 0` and `checksum off` stay apart.
			sql:  "RESTORE TABLE a FROM 'noop://' checksum 0",
			opts: []want{{name: "CHECKSUM", value: ast.BRIEValueInt, num: 0}},
		},
		{
			sql:  "BACKUP DATABASE * TO 'noop://' snapshot = '2020-03-18 18:13:54'",
			opts: []want{{name: "SNAPSHOT", value: ast.BRIEValueString, str: "2020-03-18 18:13:54"}},
		},
		{
			sql:  "BACKUP DATABASE * TO 'noop://' snapshot = 1234567890",
			opts: []want{{name: "SNAPSHOT", value: ast.BRIEValueInt, num: 1234567890}},
		},
		{
			// The count and the unit are both kept, unmultiplied.
			sql:  "BACKUP DATABASE * TO 'noop://' snapshot 5 minute ago",
			opts: []want{{name: "SNAPSHOT", value: ast.BRIEValueTimeAgo, num: 5, unit: "MINUTE"}},
		},
		{
			sql:  "BACKUP DATABASE * TO 'noop://' snapshot 5 SQL_TSI_WEEK ago",
			opts: []want{{name: "SNAPSHOT", value: ast.BRIEValueTimeAgo, num: 5, unit: "SQL_TSI_WEEK"}},
		},
		{
			sql:  "BACKUP DATABASE * TO 'noop://' LAST_BACKUP = '2020-02-02 14:14:14'",
			opts: []want{{name: "LAST_BACKUP", value: ast.BRIEValueString, str: "2020-02-02 14:14:14"}},
		},
		{
			sql:  "BACKUP DATABASE * TO 'noop://' LAST_BACKUP = 1234567890",
			opts: []want{{name: "LAST_BACKUP", value: ast.BRIEValueInt, num: 1234567890}},
		},
		{
			// 500 is the rate as written; upstream multiplies it out
			// to bytes, and the form says what the unit was.
			sql:  "backup database * to 'noop://' rate_limit 500 MB/second",
			opts: []want{{name: "RATE_LIMIT", value: ast.BRIEValueRateLimit, num: 500}},
		},
		{
			sql: "BACKUP TABLE a TO 'noop://' csv_header = fields csv_header = 3",
			opts: []want{
				{name: "CSV_HEADER", value: ast.BRIEValueCSVHeaderWords, str: "FIELDS"},
				{name: "CSV_HEADER", value: ast.BRIEValueInt, num: 3},
			},
		},
		{
			// ON DUPLICATE and ON_DUPLICATE are one option, so both
			// arrive under the one name.
			sql: "BACKUP TABLE a TO 'noop://' ON DUPLICATE = ignore on_duplicate = replace backend = `local`",
			opts: []want{
				{name: "ON_DUPLICATE", value: ast.BRIEValueString, str: "ignore"},
				{name: "ON_DUPLICATE", value: ast.BRIEValueString, str: "replace"},
				{name: "BACKEND", value: ast.BRIEValueString, str: "local"},
			},
		},
		{
			sql: "restore point from 'noop://log_backup' full_backup_storage='noop://full_log' start_ts='20230101' restored_ts='20230123'",
			opts: []want{
				{name: "FULL_BACKUP_STORAGE", value: ast.BRIEValueString, str: "noop://full_log"},
				{name: "START_TS", value: ast.BRIEValueString, str: "20230101"},
				{name: "RESTORED_TS", value: ast.BRIEValueString, str: "20230123"},
			},
		},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			node := parseOne(t, tc.sql)
			stmt, ok := node.(*ast.BRIEStmt)
			if !ok {
				t.Fatalf("Parse(%q) read a %T, want a *ast.BRIEStmt", tc.sql, node)
			}
			if len(stmt.Options) != len(tc.opts) {
				t.Fatalf("read %d options, want %d", len(stmt.Options), len(tc.opts))
			}
			for i, w := range tc.opts {
				got := stmt.Options[i]
				if got.Name != w.name {
					t.Errorf("option %d: Name = %q, want %q", i, got.Name, w.name)
				}
				if got.Value != w.value {
					t.Errorf("option %d (%s): Value = %d, want %d", i, w.name, got.Value, w.value)
				}
				if got.StrValue != w.str {
					t.Errorf("option %d (%s): StrValue = %q, want %q", i, w.name, got.StrValue, w.str)
				}
				if got.IntValue != w.num {
					t.Errorf("option %d (%s): IntValue = %d, want %d", i, w.name, got.IntValue, w.num)
				}
				if got.Unit != w.unit {
					t.Errorf("option %d (%s): Unit = %q, want %q", i, w.name, got.Unit, w.unit)
				}
				// Each option's Loc covers its own words only.
				if got.Loc.Start < 0 || got.Loc.End > len(tc.sql) || got.Loc.Start >= got.Loc.End {
					t.Errorf("option %d (%s): Loc = %v, outside %q", i, w.name, got.Loc, tc.sql)
				}
			}
		})
	}
}

// TestRecoverTableInTree asserts which of the two productions a RECOVER
// TABLE read, and that the number of `RECOVER TABLE t1 100` is a number and
// not a second table.
func TestRecoverTableInTree(t *testing.T) {
	for _, tc := range []struct {
		sql       string
		byJob     bool
		jobID     int64
		table     string
		jobNum    int64
		hasJobNum bool
	}{
		{sql: "RECOVER TABLE BY JOB 11", byJob: true, jobID: 11},
		{sql: "recover table by job 0", byJob: true, jobID: 0},
		{sql: "RECOVER TABLE BY JOB 9223372036854775807", byJob: true, jobID: 9223372036854775807},
		{sql: "recover table t1", table: "t1"},
		{sql: "recover table t1 100", table: "t1", jobNum: 100, hasJobNum: true},
		{sql: "recover table t1 0", table: "t1", jobNum: 0, hasJobNum: true},
		{sql: "recover table a.b 100", table: "a.b", jobNum: 100, hasJobNum: true},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			node := parseOne(t, tc.sql)
			stmt, ok := node.(*ast.RecoverTableStmt)
			if !ok {
				t.Fatalf("Parse(%q) read a %T, want a *ast.RecoverTableStmt", tc.sql, node)
			}
			if stmt.ByJob != tc.byJob {
				t.Fatalf("ByJob = %v, want %v", stmt.ByJob, tc.byJob)
			}
			if stmt.JobID != tc.jobID {
				t.Errorf("JobID = %d, want %d", stmt.JobID, tc.jobID)
			}
			if tc.byJob {
				if stmt.Table != nil {
					t.Errorf("Table = %s, want none in the BY JOB form", qualified(stmt.Table))
				}
				if stmt.HasJobNum {
					t.Errorf("HasJobNum = true, want false in the BY JOB form")
				}
				return
			}
			if got := qualified(stmt.Table); got != tc.table {
				t.Errorf("Table = %q, want %q", got, tc.table)
			}
			if stmt.HasJobNum != tc.hasJobNum {
				t.Errorf("HasJobNum = %v, want %v", stmt.HasJobNum, tc.hasJobNum)
			}
			if stmt.JobNum != tc.jobNum {
				t.Errorf("JobNum = %d, want %d", stmt.JobNum, tc.jobNum)
			}
		})
	}
}

// TestFlashbackInTree asserts that the bring-back forms and the
// point-in-time forms land in different nodes, and that each keeps its own
// operands: the new name of a rename, the scope and the point of a flashback.
func TestFlashbackInTree(t *testing.T) {
	// The bring-back forms.
	for _, tc := range []struct {
		sql     string
		table   string // set for FLASHBACK TABLE
		db      string // set for FLASHBACK DATABASE / SCHEMA
		newName string
	}{
		{sql: "flashback table t", table: "t"},
		{sql: "flashback table t TO t1", table: "t", newName: "t1"},
		{sql: "flashback table a.b to c", table: "a.b", newName: "c"},
		{sql: "flashback table `t` to `t1`", table: "t", newName: "t1"},
		// TO timestamp with no literal after it is a rename to a table
		// called timestamp, which is why it is here and not below.
		{sql: "flashback table t TO timestamp", table: "t", newName: "timestamp"},
		{sql: "flashback table t TO tso", table: "t", newName: "tso"},
		{sql: "flashback database db1", db: "db1"},
		{sql: "flashback schema db1", db: "db1"},
		{sql: "flashback database db1 to db2", db: "db1", newName: "db2"},
		{sql: "flashback schema db1 to db2", db: "db1", newName: "db2"},
		{sql: "flashback database d TO timestamp", db: "d", newName: "timestamp"},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			node := parseOne(t, tc.sql)
			switch n := node.(type) {
			case *ast.FlashbackTableStmt:
				if tc.table == "" {
					t.Fatalf("Parse(%q) read a FlashbackTableStmt, want a FlashbackDatabaseStmt", tc.sql)
				}
				if got := qualified(n.Table); got != tc.table {
					t.Errorf("Table = %q, want %q", got, tc.table)
				}
				if n.NewName != tc.newName {
					t.Errorf("NewName = %q, want %q", n.NewName, tc.newName)
				}
				if (n.NewName == "") != (n.NewNameLoc == ast.NoLoc()) {
					t.Errorf("NewName = %q with NewNameLoc = %v", n.NewName, n.NewNameLoc)
				}
			case *ast.FlashbackDatabaseStmt:
				if tc.db == "" {
					t.Fatalf("Parse(%q) read a FlashbackDatabaseStmt, want a FlashbackTableStmt", tc.sql)
				}
				if n.Name != tc.db {
					t.Errorf("Name = %q, want %q", n.Name, tc.db)
				}
				if n.NewName != tc.newName {
					t.Errorf("NewName = %q, want %q", n.NewName, tc.newName)
				}
				if (n.NewName == "") != (n.NewNameLoc == ast.NoLoc()) {
					t.Errorf("NewName = %q with NewNameLoc = %v", n.NewName, n.NewNameLoc)
				}
			default:
				t.Fatalf("Parse(%q) read a %T, want a bring-back FLASHBACK statement", tc.sql, node)
			}
		})
	}

	// The point-in-time forms.
	for _, tc := range []struct {
		sql    string
		scope  ast.FlashbackScope
		tables []string
		db     string
		isTSO  bool
		ts     string
		tso    uint64
	}{
		{
			sql:   "flashback cluster to timestamp '2021-05-26 16:45:26'",
			scope: ast.FlashbackCluster,
			ts:    "2021-05-26 16:45:26",
		},
		{
			sql:   "flashback cluster to tso 445494955052105721",
			scope: ast.FlashbackCluster,
			isTSO: true,
			tso:   445494955052105721,
		},
		{
			sql:    "flashback table t to timestamp '2021-05-26 16:45:26'",
			scope:  ast.FlashbackTables,
			tables: []string{"t"},
			ts:     "2021-05-26 16:45:26",
		},
		{
			sql:    "flashback table t,t1 to tso 445494955052105723",
			scope:  ast.FlashbackTables,
			tables: []string{"t", "t1"},
			isTSO:  true,
			tso:    445494955052105723,
		},
		{
			sql:   "flashback schema test to timestamp '2021-05-26 16:45:26'",
			scope: ast.FlashbackDatabase,
			db:    "test",
			ts:    "2021-05-26 16:45:26",
		},
		{
			sql:   "flashback database test to tso 445494955052105724",
			scope: ast.FlashbackDatabase,
			db:    "test",
			isTSO: true,
			tso:   445494955052105724,
		},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			node := parseOne(t, tc.sql)
			stmt, ok := node.(*ast.FlashbackToTimestampStmt)
			if !ok {
				t.Fatalf("Parse(%q) read a %T, want a *ast.FlashbackToTimestampStmt", tc.sql, node)
			}
			if stmt.Scope != tc.scope {
				t.Errorf("Scope = %d, want %d", stmt.Scope, tc.scope)
			}
			names := make([]string, 0, len(stmt.Tables))
			for _, ref := range stmt.Tables {
				names = append(names, qualified(ref))
			}
			if strings.Join(names, ",") != strings.Join(tc.tables, ",") {
				t.Errorf("Tables = %q, want %q", names, tc.tables)
			}
			if stmt.Database != tc.db {
				t.Errorf("Database = %q, want %q", stmt.Database, tc.db)
			}
			if stmt.IsTSO != tc.isTSO {
				t.Fatalf("IsTSO = %v, want %v", stmt.IsTSO, tc.isTSO)
			}
			if stmt.Timestamp != tc.ts {
				t.Errorf("Timestamp = %q, want %q", stmt.Timestamp, tc.ts)
			}
			if stmt.TSO != tc.tso {
				t.Errorf("TSO = %d, want %d", stmt.TSO, tc.tso)
			}
			if stmt.Loc.Start != 0 || stmt.Loc.End != len(tc.sql) {
				t.Errorf("Loc = %v, want {0 %d}", stmt.Loc, len(tc.sql))
			}
		})
	}
}

// TestFlashbackKeywordStaysIdentifier pins that registering FLASHBACK in the
// keyword table left it an identifier everywhere else, which is what v8.5.5
// does with it: it is TiDBKeyword upstream, not a reserved word. BACKUP,
// RESTORE and RECOVER were already registered and are here for the same
// reason, since this family gave each of them a statement to start.
func TestFlashbackKeywordStaysIdentifier(t *testing.T) {
	for _, sql := range []string{
		"CREATE TABLE flashback (flashback int)",
		"SELECT flashback FROM flashback",
		"SELECT * FROM t flashback",
		"SELECT * FROM t AS flashback",
		"DELETE FROM t flashback WHERE flashback.a = 1",
		"UPDATE t SET flashback = 1",
		"CREATE INDEX flashback ON t (a)",
		"CREATE DATABASE flashback",
		"SELECT 1 AS flashback",
		"SELECT * FROM flashback.t",
		"CREATE TABLE t (backup int, restore int, recover int)",
		"SELECT * FROM t backup",
		"SELECT * FROM t restore",
		"SELECT * FROM t recover",
		"DELETE FROM t backup WHERE backup.a = 1",
	} {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Errorf("Parse(%q): %v", sql, err)
			}
		})
	}
	if IsReservedKeyword("flashback") {
		t.Error("IsReservedKeyword(\"flashback\") = true, want false")
	}
}
