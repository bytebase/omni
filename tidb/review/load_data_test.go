package review

import (
	"testing"

	"github.com/bytebase/omni/review"
)

// IMPORT INTO in the WALK_THROUGH walk.
//
// It creates, drops and renames nothing a walk keeps, so it is lowered the
// way LOAD DATA is: followed, claiming only that its target is there. What
// that has to earn is the statement after it. A node the walk does not lower
// returns nil from lower(), which ends the walk there, so an IMPORT INTO in
// the middle of a change script would hide every finding below it -- a change
// that adds a column twice would go unreported because the walk never
// reached the second statement.
//
// The one thing it does claim is the target: v8.5.5 answers 1146 for `IMPORT
// INTO nosuch FROM '/f.csv'` (read on 2026-11-04), so the walk reports a
// missing table rather than passing a script that would fail.
//
// These tests need no container and run unconditionally.

// importIntoWalkStatements are the forms the walk must follow, written
// against walkSchema's table t.
var importIntoWalkStatements = []string{
	`IMPORT INTO t FROM '/f.csv'`,
	`IMPORT INTO t (a) FROM '/f.csv' FORMAT 'csv'`,
	`IMPORT INTO t (a, @1) SET b = @1 FROM '/f.csv' WITH thread=2, detached`,
	`IMPORT INTO t FROM SELECT a FROM t`,
	`IMPORT INTO t FROM TABLE t`,
	`LOAD DATA INFILE '/f.csv' FORMAT 'csv' INTO TABLE t WITH detached`,
	`LOAD DATA INFILE '/f.csv' INTO TABLE t FIELDS TERMINATED BY X'6B6B' DEFINED NULL BY 'nil'`,
}

func TestWalkThroughFollowsImportInto(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	for _, sql := range importIntoWalkStatements {
		t.Run(sql, func(t *testing.T) {
			// On its own it reports nothing: its table is there and it
			// changes no name.
			noFinding(t, sql+";", shop)

			// The walk carries on past it, so the error below is still
			// found, and it is reported against statement 2 -- the
			// statement's own index, not the bulk load's.
			oneFinding(t, sql+";\nALTER TABLE t ADD COLUMN c int;\nALTER TABLE t ADD COLUMN c int;",
				2, "Column `c` already exists in `t`.", shop)

			// And the names a statement above it set are still there
			// below it, so the walk did not lose its state either.
			noFinding(t, "ALTER TABLE t ADD COLUMN c int;\n"+sql+";\nALTER TABLE t DROP COLUMN c;", shop)
		})
	}
}

// TestWalkThroughImportIntoNeedsItsTable pins the one claim IMPORT INTO
// makes: the table it writes into has to exist, and a rename above it counts.
func TestWalkThroughImportIntoNeedsItsTable(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	oneFinding(t, "IMPORT INTO nosuch FROM '/f.csv';", 0, "`nosuch` doesn't exist.", shop)
	oneFinding(t, "ALTER TABLE t RENAME TO u;\nIMPORT INTO t FROM '/f.csv';", 1, "`t` doesn't exist.", shop)
	noFinding(t, "ALTER TABLE t RENAME TO u;\nIMPORT INTO u FROM '/f.csv';", shop)
}
