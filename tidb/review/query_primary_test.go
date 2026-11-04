package review

import (
	"testing"

	"github.com/bytebase/omni/review"
)

// TestWalkThroughTiDBQueryPrimary pins what the walk does with a TABLE or
// VALUES query primary: it follows it and records nothing, as it does a
// SELECT. These statements reach the walk through the shared query grammar,
// and a statement the walk does not lower ends the walk before itself, so one
// in the middle of a change would leave every statement after it unjudged --
// which is what the second half of each pair here would miss.
func TestWalkThroughTiDBQueryPrimary(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}

	// A TABLE or VALUES query expression reads; it makes and drops nothing,
	// whatever name it holds.
	noFinding(t, "TABLE nosuch;", shop)
	noFinding(t, "TABLE t;", shop)
	noFinding(t, "TABLE t ORDER BY a LIMIT 1;", shop)
	noFinding(t, "TABLE nosuch UNION TABLE t;", shop)
	noFinding(t, "VALUES ROW(1,2), ROW(3,4);", shop)
	noFinding(t, "VALUES ROW(1) UNION SELECT 1;", shop)
	noFinding(t, "SELECT * FROM t FOR UPDATE OF t WAIT 5;", shop)
	noFinding(t, "TABLE t FOR UPDATE NOWAIT;", shop)
	noFinding(t, "EXPLAIN TABLE t UNION TABLE t;", shop)

	// The walk carries on past one, so the statement after it is judged.
	oneFinding(t, "TABLE t;\nCREATE TABLE t (a int);", 1, "Table `t` already exists.", shop)
	oneFinding(t, "VALUES ROW(1);\nTABLE t;\nALTER TABLE nosuch ADD COLUMN c int;", 2,
		"Table `nosuch` doesn't exist.", shop)
}
