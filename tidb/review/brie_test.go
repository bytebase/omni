package review

import (
	"testing"

	"github.com/bytebase/omni/review"
)

// TiDB's backup, restore, recover and flashback statements in the review.
//
// None of them is lowered, so each ends the walk, which is where LOCK TABLES
// and LOCK STATS already are and for the same two reasons. RESTORE, RECOVER
// TABLE and the bring-back FLASHBACK forms each put a relation or a schema
// back without saying what is in it, and FLASHBACK CLUSTER ... TO TIMESTAMP
// rolls the whole cluster back, so the walk cannot keep its names accurate
// past any of them. BACKUP changes nothing at all, but it does fail on a
// table that is not there -- BACKUP TABLE nosuch answers 1146 on
// pingcap/tidb:v8.5.5, read on 2026-11-04 -- so listing it as a statement
// with no effect would drop a rejection the review ought to report rather
// than add one it should not.
//
// What ending the walk costs is the statements below one, and that is the
// trade these tests record rather than hide.
//
// They are not in walk_through_oracle_test.go's curated cases: every
// statement here needs a TiKV-backed cluster that the harness's single
// tidb-server is not, so the server answers 1105 ("BACKUP requires tikv
// store, not unistore", and "can not get 'tikv_gc_safe_point'" for
// FLASHBACK) to all of them, which says nothing about the walk.
//
// These tests need no container and run unconditionally.

// brieChanges is one of each form, enough that a reader of the family would
// have to lower all of it to change any of these expectations.
var brieChanges = []string{
	`BACKUP DATABASE * TO 'noop://'`,
	`BACKUP DATABASE shop TO 'noop://' snapshot 5 minute ago`,
	`BACKUP TABLE t TO 'noop://' concurrency 4 checksum off`,
	`BACKUP LOGS TO 'noop://' start_ts = '20220304'`,
	`RESTORE DATABASE * FROM 'noop://'`,
	`RESTORE TABLE t FROM 'noop://' online 1`,
	`RESTORE POINT FROM 'noop://log_backup' full_backup_storage = 'noop://full'`,
	`RECOVER TABLE BY JOB 11`,
	`RECOVER TABLE t`,
	`RECOVER TABLE t 100`,
	`FLASHBACK TABLE t`,
	`FLASHBACK TABLE t TO t1`,
	`FLASHBACK DATABASE shop`,
	`FLASHBACK DATABASE shop TO other`,
	`FLASHBACK SCHEMA shop TO other`,
	`FLASHBACK CLUSTER TO TIMESTAMP '2021-05-26 16:45:26'`,
	`FLASHBACK CLUSTER TO TSO 445494955052105721`,
	`FLASHBACK TABLE t TO TIMESTAMP '2021-05-26 16:45:26'`,
	`FLASHBACK DATABASE shop TO TSO 445494955052105724`,
}

// TestWalkThroughEndsAtBRIEAndFlashback pins that each statement is read
// without a finding and that the walk claims nothing after it.
func TestWalkThroughEndsAtBRIEAndFlashback(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	for _, sql := range brieChanges {
		t.Run(sql, func(t *testing.T) {
			// On its own it reports nothing. Before this grammar
			// landed it was a syntax error, which the review reports
			// as a finding at its position.
			noFinding(t, sql+";", shop)

			// A name that is not there is not reported either: the
			// walk claims nothing about these statements, so it
			// claims nothing about their operands.
			noFinding(t, "BACKUP TABLE nosuch TO 'noop://';", shop)

			// And the walk ends there, so the error below it goes
			// unreported. That is the cost of claiming nothing, and
			// it is the same cost LOCK TABLES has always had.
			noFinding(t, sql+";\nALTER TABLE t DROP COLUMN zz;", shop)

			// A statement above one is still judged, so ending the
			// walk here loses nothing before it.
			oneFinding(t, "ALTER TABLE t DROP COLUMN zz;\n"+sql+";", 0,
				"Column `zz` doesn't exist in `t`.", shop)
		})
	}
}

// TestDisallowRenameIgnoresFlashback pins that FLASHBACK TABLE ... TO ... is
// not reported as a rename, which is the one shape in the family that looks
// like one. It is not a rename: the table it names has to be gone already for
// the statement to succeed, and the new name is created, so nothing a
// consumer of the old name depends on moves. Its sibling makes that plain --
// v8.5.5 answers 1007, "Can't create database 'd'; database exists", for
// FLASHBACK DATABASE d when d is still there.
//
// The rule's own subject is pinned alongside it, so a reader who broadened
// the rule to cover FLASHBACK would see both halves at once.
func TestDisallowRenameIgnoresFlashback(t *testing.T) {
	runCases(t, review.DisallowRename, []ruleCase{
		{name: "FLASHBACK TABLE to a new name", sql: "FLASHBACK TABLE t TO t1;"},
		{name: "FLASHBACK TABLE of a qualified name", sql: "FLASHBACK TABLE shop.t TO t1;"},
		{name: "FLASHBACK DATABASE to a new name", sql: "FLASHBACK DATABASE shop TO other;"},
		{name: "FLASHBACK SCHEMA to a new name", sql: "FLASHBACK SCHEMA shop TO other;"},
		{name: "FLASHBACK TABLE with no new name", sql: "FLASHBACK TABLE t;"},
		{name: "RECOVER TABLE", sql: "RECOVER TABLE t;"},
		{name: "RECOVER TABLE BY JOB", sql: "RECOVER TABLE BY JOB 11;"},
		{name: "BACKUP and RESTORE", sql: "BACKUP TABLE t TO 'noop://';\nRESTORE TABLE t FROM 'noop://';"},
		// A real rename still reports, so the absence above is about
		// these statements and not about the rule being off.
		{
			name: "RENAME TABLE still reports",
			sql:  "RENAME TABLE t TO t1;",
			want: []finding{{0, "t TO t1", "Renames table `t` to `t1`.", []int{0}}},
		},
		{
			name: "ALTER TABLE ... RENAME still reports",
			sql:  "ALTER TABLE t RENAME TO t1;",
			want: []finding{{0, "RENAME TO t1", "Renames table `t` to `t1`.", []int{0}}},
		},
	})
}
