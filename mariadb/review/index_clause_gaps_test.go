package review

import (
	"testing"

	"github.com/bytebase/omni/review"
)

// TestWalkThroughMariaDBAddIfNotExists: what the clause finds is a note, each read from
// mariadb:11.8.8, and the walk goes on from the table as it was.
func TestWalkThroughMariaDBAddIfNotExists(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	noFinding(t, "ALTER TABLE t ADD INDEX IF NOT EXISTS i (a);", shop)
	noFinding(t, "ALTER TABLE t ADD UNIQUE IF NOT EXISTS i (a);", shop)
	noFinding(t, "ALTER TABLE t ADD PRIMARY KEY IF NOT EXISTS (a);", shop)
	noFinding(t, "ALTER TABLE t ADD FOREIGN KEY IF NOT EXISTS fk1 (a) REFERENCES r (id);", shop)
	noFinding(t, "ALTER TABLE t ADD CONSTRAINT fk1 FOREIGN KEY IF NOT EXISTS (a) REFERENCES r (id);", shop)
	noFinding(t, "ALTER TABLE t ADD INDEX IF NOT EXISTS n (a), ADD INDEX IF NOT EXISTS n (a);", shop)
	noFinding(t, "ALTER TABLE t ADD (INDEX IF NOT EXISTS n (a), INDEX IF NOT EXISTS n (a));", shop)
	noFinding(t, "ALTER TABLE t ADD PRIMARY KEY IF NOT EXISTS (a), ADD PRIMARY KEY IF NOT EXISTS (a);", shop)
	oneFinding(t, "ALTER TABLE nokey ADD PRIMARY KEY IF NOT EXISTS (a), ADD PRIMARY KEY IF NOT EXISTS (a);", 0, "`nokey` already has a primary key.", shop)
	oneFinding(t, "ALTER TABLE t ADD INDEX IF NOT EXISTS n (a), ADD INDEX n (a);", 0, "This statement adds index `n` twice.", shop)
	noFinding(t, "ALTER TABLE t ADD INDEX n (a), ADD INDEX IF NOT EXISTS n (a);", shop)
	// What it adds is there for what comes after.
	oneFinding(t, "ALTER TABLE t ADD INDEX IF NOT EXISTS n (a);\nALTER TABLE t ADD INDEX n (a);", 1, "Index `n` already exists in `t`.", shop)
	noFinding(t, "ALTER TABLE t ADD INDEX IF NOT EXISTS n (a);\nALTER TABLE t DROP INDEX n;", shop)
	oneFinding(t, "ALTER TABLE nokey ADD PRIMARY KEY IF NOT EXISTS (a);\nALTER TABLE nokey ADD PRIMARY KEY (a);", 1, "`nokey` already has a primary key.", shop)
	// The same text without the clause still fails.
	oneFinding(t, "ALTER TABLE t ADD INDEX i (a);", 0, "Index `i` already exists in `t`.", shop)
	oneFinding(t, "ALTER TABLE t ADD PRIMARY KEY (a);", 0, "`t` already has a primary key.", shop)
	// WAIT changes nothing about what the walk follows.
	oneFinding(t, "CREATE INDEX n ON t (a) WAIT 3;\nCREATE INDEX n ON t (a) NOWAIT;", 1, "Index `n` already exists in `t`.", shop)
	oneFinding(t, "DROP INDEX i ON t WAIT 3;\nDROP INDEX i ON t NOWAIT;", 1, "Index `i` doesn't exist in `t`.", shop)
}
