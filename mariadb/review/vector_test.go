package review

import (
	"testing"

	"github.com/bytebase/omni/review"
)

// TestWalkThroughMariaDBVector: a vector key in a change is an index the walk goes on
// from, in every spelling.
func TestWalkThroughMariaDBVector(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	noFinding(t, "CREATE TABLE x (id int, v VECTOR(3) NOT NULL, VECTOR KEY vi (v));\nALTER TABLE x DROP INDEX vi;\nCREATE VECTOR INDEX vj ON x (v) M=6;\nDROP INDEX vj ON x;\nALTER TABLE x DROP COLUMN v;", shop)
	oneFinding(t, "CREATE VECTOR INDEX vi ON t (a) M=6;\nCREATE VECTOR INDEX VI ON t (a);", 1, "Index `VI` already exists in `t`.", shop)
	oneFinding(t, "CREATE VECTOR INDEX vi ON t (a);\nDROP INDEX vi ON t;\nDROP INDEX vi ON t;", 2, "Index `vi` doesn't exist in `t`.", shop)
	oneFinding(t, "CREATE VECTOR INDEX I ON t (a);", 0, "Index `I` already exists in `t`.", shop)
	noFinding(t, "CREATE VECTOR INDEX IF NOT EXISTS I ON t (a);", shop)
	noFinding(t, "CREATE OR REPLACE VECTOR INDEX I ON t (a);", shop)
	oneFinding(t, "CREATE VECTOR INDEX vi ON nosuch (a);", 0, "Table `nosuch` doesn't exist.", shop)
	oneFinding(t, "CREATE VECTOR INDEX vi ON t (a);\nALTER TABLE nosuch ADD COLUMN c int;", 1, "Table `nosuch` doesn't exist.", shop)
	oneFinding(t, "ALTER TABLE t ADD VECTOR vi (a);\nALTER TABLE t ADD VECTOR KEY vi (a);", 1, "Index `vi` already exists in `t`.", shop)
	oneFinding(t, "ALTER TABLE t ADD VECTOR INDEX I (a);", 0, "Index `I` already exists in `t`.", shop)
	noFinding(t, "ALTER TABLE t ADD VECTOR IF NOT EXISTS I (a);\nALTER TABLE t ADD VECTOR IF NOT EXISTS vi (a);\nALTER TABLE t DROP INDEX vi;", shop)
	oneFinding(t, "CREATE TABLE x (v VECTOR(3) NOT NULL, VECTOR vi (v));\nCREATE INDEX vi ON x (v);", 1, "Index `vi` already exists in `x`.", shop)
	// IF NOT EXISTS within one statement, each as 11.8.8 answers it: a repeat of a vector key
	// is a note, and beside another kind of key or before a plain ADD it is 1061.
	noFinding(t, "ALTER TABLE t ADD VECTOR IF NOT EXISTS I (a), ADD VECTOR IF NOT EXISTS I (a);", shop)
	noFinding(t, "ALTER TABLE t ADD VECTOR IF NOT EXISTS vi (a), ADD VECTOR IF NOT EXISTS vi (a);", shop)
	noFinding(t, "ALTER TABLE t ADD VECTOR vi (a), ADD VECTOR IF NOT EXISTS vi (a);", shop)
	if got := walkFindings(t, "ALTER TABLE t ADD VECTOR IF NOT EXISTS vi (a), ADD VECTOR vi (a);", shop); len(got) != 1 {
		t.Errorf("IF NOT EXISTS then a plain ADD: findings = %+v, want one", got)
	}
	if got := walkFindings(t, "ALTER TABLE t ADD INDEX vi (a), ADD VECTOR IF NOT EXISTS vi (a);", shop); len(got) != 1 {
		t.Errorf("another kind of key then IF NOT EXISTS: findings = %+v, want one", got)
	}
	if got := walkFindings(t, "ALTER TABLE t ADD VECTOR IF NOT EXISTS vi (a), ADD INDEX vi (a);", shop); len(got) != 1 {
		t.Errorf("IF NOT EXISTS then another kind of key: findings = %+v, want one", got)
	}
	oneFinding(t, "ALTER TABLE t ADD COLUMN c VECTOR(3);\nALTER TABLE t ADD COLUMN C VECTOR(3);", 1, "Column `C` already exists in `t`.", shop)
}
