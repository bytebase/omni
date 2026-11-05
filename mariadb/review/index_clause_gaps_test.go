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
	// A CHECK and an index have a name each, and a second index of a name goes as a note only
	// of its own kind.
	noFinding(t, "ALTER TABLE t ADD CONSTRAINT IF NOT EXISTS ck CHECK (a > 0), ADD INDEX ck (a);", shop)
	noFinding(t, "ALTER TABLE t ADD INDEX n (a), ADD KEY IF NOT EXISTS n (a);", shop)
	noFinding(t, "ALTER TABLE t ADD UNIQUE IF NOT EXISTS i (a), ADD FULLTEXT IF NOT EXISTS i (a);", shop)
	oneFinding(t, "ALTER TABLE t ADD INDEX n (a), ADD UNIQUE IF NOT EXISTS n (a);", 0, "This statement adds index `n` twice.", shop)
	oneFinding(t, "ALTER TABLE t ADD INDEX IF NOT EXISTS n (a), ADD UNIQUE IF NOT EXISTS n (a);", 0, "This statement adds index `n` twice.", shop)
	// A foreign key looks for a foreign key of the name on this table: fk1 is another table's
	// here, which is the server's 1005, and an index of the name does not stand in for one.
	if got := walkFindings(t, "ALTER TABLE nokey ADD FOREIGN KEY IF NOT EXISTS fk1 (a) REFERENCES r (id);", shop); len(got) != 1 {
		t.Errorf("a foreign key name another table owns: findings = %+v, want one", got)
	}
	noFinding(t, "ALTER TABLE t ADD FOREIGN KEY IF NOT EXISTS i (a) REFERENCES r (id);\nALTER TABLE t DROP FOREIGN KEY i;", shop)
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
