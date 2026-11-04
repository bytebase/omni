package review

import (
	"testing"

	"github.com/bytebase/omni/review"
)

// TestWalkThroughTiDBAdmin pins what the walk does with TiDB's ADMIN family.
//
// Two things matter here and neither shows up in a parser test. The first is
// that the walk follows these statements at all: a statement the walk does not
// lower ends the walk, so one ADMIN statement in a change would have hidden
// every later finding. The second is that a kind naming relations needs them,
// which is the only thing in the family the walk can predict -- v8.5.5 answers
// 1146 for each form below, and answers nothing the walk claims for a missing
// index, a temporary table, or a server not in repair mode.
func TestWalkThroughTiDBAdmin(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}

	// A kind that names relations needs each one to exist (1146).
	oneFinding(t, "ADMIN CHECK TABLE nosuch;", 0, "`nosuch` doesn't exist.", shop)
	oneFinding(t, "ADMIN CHECK TABLE t, nosuch;", 0, "`nosuch` doesn't exist.", shop)
	oneFinding(t, "ADMIN CHECKSUM TABLE nosuch;", 0, "`nosuch` doesn't exist.", shop)
	oneFinding(t, "ADMIN SHOW nosuch NEXT_ROW_ID;", 0, "`nosuch` doesn't exist.", shop)
	oneFinding(t, "ADMIN CLEANUP TABLE LOCK nosuch;", 0, "`nosuch` doesn't exist.", shop)
	oneFinding(t, "ADMIN CHECK INDEX nosuch i;", 0, "`nosuch` doesn't exist.", shop)
	oneFinding(t, "ADMIN CHECK INDEX nosuch i (1, 2);", 0, "`nosuch` doesn't exist.", shop)
	oneFinding(t, "ADMIN RECOVER INDEX nosuch i;", 0, "`nosuch` doesn't exist.", shop)
	oneFinding(t, "ADMIN CLEANUP INDEX nosuch i;", 0, "`nosuch` doesn't exist.", shop)
	// The relation may be of any kind: ADMIN CHECK TABLE and ADMIN CHECKSUM
	// TABLE both succeed on a view and on a sequence on v8.5.5.
	noFinding(t, "ADMIN CHECK TABLE t, v, sq;", shop)
	noFinding(t, "ADMIN CHECKSUM TABLE v, sq;", shop)
	noFinding(t, "ADMIN SHOW v NEXT_ROW_ID;", shop)
	// A table the change made, and a table it dropped.
	noFinding(t, "CREATE TABLE y (a int);\nADMIN CHECK TABLE y;\nADMIN CHECKSUM TABLE y;", shop)
	oneFinding(t, "DROP TABLE t;\nADMIN CHECK TABLE t;", 1, "`t` doesn't exist.", shop)

	// A missing index is a runtime answer (1105 for CHECK INDEX, RECOVER
	// INDEX, and CLEANUP INDEX alike), not a fact the walk claims.
	noFinding(t, "ADMIN CHECK INDEX t nosuchindex;", shop)
	noFinding(t, "ADMIN RECOVER INDEX t nosuchindex;", shop)

	// A kind that names nothing needs nothing, and the walk keeps going
	// after it: without that, one ADMIN statement would hide the ALTER.
	oneFinding(t, "ADMIN SHOW DDL JOBS;\nALTER TABLE nosuch ADD COLUMN c int;", 1, "Table `nosuch` doesn't exist.", shop)
	oneFinding(t, "ADMIN FLUSH BINDINGS;\nADMIN SHOW SLOW TOP 3;\nADMIN RELOAD BINDINGS;\nDROP TABLE nosuch;", 3, "Table `nosuch` doesn't exist.", shop)

	// ADMIN REPAIR TABLE carries a table's whole definition and creates no
	// table: outside repair mode v8.5.5 answers 8215, and in repair mode it
	// rewrites the metadata of a table that is already there. So the name
	// the definition gives is not taken afterwards, and is not there to
	// alter either.
	noFinding(t, "ADMIN REPAIR TABLE t CREATE TABLE x (a int);\nCREATE TABLE x (b int);", shop)
	oneFinding(t, "ADMIN REPAIR TABLE t CREATE TABLE t2 (a int);\nALTER TABLE t2 ADD COLUMN b int;", 1, "Table `t2` doesn't exist.", shop)
	// Nor does the repaired table have to exist: 8215 comes first, so the
	// walk claims nothing about the name.
	noFinding(t, "ADMIN REPAIR TABLE nosuch CREATE TABLE nosuch (a int);", shop)
}

// TestAdminRepairTableIsNotACreateTable pins the rule side of the same
// decision: a rule that judges a CREATE TABLE reads the statement's top-level
// node, and ADMIN REPAIR TABLE is not one. Reporting "Creates table t without
// a primary key." for a statement that creates no table would be wrong twice
// -- the statement repairs rather than creates, and outside repair mode it
// does not run at all.
func TestAdminRepairTableIsNotACreateTable(t *testing.T) {
	runCases(t, review.RequirePrimaryKey, []ruleCase{
		{
			name: "ADMIN REPAIR TABLE's definition is not a CREATE TABLE",
			sql:  "ADMIN REPAIR TABLE t CREATE TABLE nokeys (id INT);\nCREATE TABLE plainly (id INT);",
			want: []finding{{1, "CREATE TABLE plainly (id INT)", "Creates table `plainly` without a primary key.", []int{0}}},
		},
	})
}
