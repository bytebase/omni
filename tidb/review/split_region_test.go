package review

import (
	"testing"

	"github.com/bytebase/omni/review"
)

// SPLIT ... TABLE and DISTRIBUTE TABLE in the WALK_THROUGH walk.
//
// Neither creates, drops, or renames a name the walk keeps, so each is
// lowered the way ANALYZE TABLE is: followed, recording nothing. What that
// has to earn is the statement after it. A node the walk does not lower
// returns nil from lower(), which ends the walk there, so a SPLIT in the
// middle of a change script would hide every error below it -- a change that
// drops a column twice would go unreported because the walk never reached the
// second statement.
//
// These tests need no container and run unconditionally.

func TestWalkThroughFollowsSplitAndDistribute(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	for _, sql := range []string{
		`SPLIT TABLE t BY ('a'),('b')`,
		`SPLIT TABLE t INDEX i BY ('a')`,
		`SPLIT REGION FOR PARTITION TABLE t PARTITION (p0) INDEX i BETWEEN ('a') AND ('z') REGIONS 10`,
		`SPLIT TABLE t BETWEEN () AND () REGIONS 10`,
		`DISTRIBUTE TABLE t RULE = 'leader-scatter' ENGINE = 'tikv'`,
		`DISTRIBUTE TABLE t PARTITION (p0) RULE = 'peer-scatter' ENGINE = 'tiflash' TIMEOUT = '30m'`,
	} {
		t.Run(sql, func(t *testing.T) {
			// On its own it reports nothing: it changes no name.
			noFinding(t, sql+";", shop)

			// The walk carries on past it, so the error below is still
			// found, and it is reported against statement 2 -- the
			// statement's own index, not the SPLIT's.
			oneFinding(t, sql+";\nALTER TABLE t ADD COLUMN c int;\nALTER TABLE t ADD COLUMN c int;",
				2, "Column `c` already exists in `t`.", shop)

			// And the names a statement above it set are still there
			// below it, so the walk did not lose its state either.
			noFinding(t, "ALTER TABLE t ADD COLUMN c int;\n"+sql+";\nALTER TABLE t DROP COLUMN c;", shop)
		})
	}
}

// TestWalkThroughSplitKeepsTheTableName pins that a SPLIT in the middle does
// not make the walk forget a rename either: the table it names is read as a
// name, not as a new object.
func TestWalkThroughSplitKeepsTheTableName(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	oneFinding(t, "ALTER TABLE t RENAME TO u;\nSPLIT TABLE u BY ('a');\nALTER TABLE t ADD COLUMN c int;",
		2, "Table `t` doesn't exist.", shop)
	oneFinding(t, "ALTER TABLE t RENAME TO u;\nDISTRIBUTE TABLE u RULE = 'x' ENGINE = 'tikv';\nCREATE TABLE u (a int);",
		2, "Table `u` already exists.", shop)
}
