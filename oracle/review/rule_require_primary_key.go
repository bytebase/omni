package review

import (
	"strings"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/review"
)

// RequirePrimaryKey reports two things, each only when the change ends
// that way: a table a CREATE TABLE made without a primary key, and a table
// whose primary key the change drops without adding one back, by DROP
// PRIMARY KEY or by DROP CONSTRAINT of a constraint scanConstraints knows
// is the primary key. A table that had none before the change is never
// reported: a drop is reported on a table the change made, or on one whose
// first primary-key event in the change is the drop, which only succeeds
// on a table that has one. CREATE TABLE ... AS SELECT copies no key from
// its query, so it is judged by what it declares.
//
// Not judged: a GLOBAL or PRIVATE TEMPORARY table; a CREATE TABLE IF NOT
// EXISTS, which may make nothing; an external table, which takes a
// primary key only as RELY DISABLE (ORA-30657 otherwise), never enforced.
//
// Without simulating the change, the rule follows the table by name. A
// later statement that may give it a key, drop it, rename it, restore it
// with FLASHBACK, or create a table of its name settles it, and it is not
// reported.
//
// oracle: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CREATE-TABLE.html

// pendingTable is a table without a primary key that the rule reports at
// the end of the change unless a later statement settles it.
type pendingTable struct {
	key       tableKey
	statement int
	rng       review.Range
	message   string
	settled   bool
}

// judgedTable reports whether RequirePrimaryKey judges a CREATE TABLE.
func judgedTable(v *ast.CreateTableStmt) bool {
	return !v.Global && !v.Private && !v.IfNotExists && !strings.EqualFold(v.Organization, "EXTERNAL")
}

// createsPrimaryKey reports whether a CREATE TABLE declares a primary key,
// on a column or of the table.
func createsPrimaryKey(v *ast.CreateTableStmt) bool {
	for _, c := range tableConstraints(v) {
		if c.kind == ast.CONSTRAINT_PRIMARY {
			return true
		}
	}
	return false
}

// settle stops following every pending table the key may name.
func (s *constraintScan) settle(k tableKey) {
	s.settleWhere(func(p tableKey) bool { return p.mayBe(k) })
}

// settleWhere stops following every pending table whose key matches.
func (s *constraintScan) settleWhere(f func(tableKey) bool) {
	for _, p := range s.pending {
		if !p.settled && f(p.key) {
			p.settled = true
		}
	}
}

// gainedPrimaryKey follows a statement that gives the table a primary key,
// which settles it.
func (s *constraintScan) gainedPrimaryKey(k tableKey) {
	t := s.state(k)
	if t.pkFirst == "" {
		t.pkFirst = "add"
	}
	t.pkDropped = false
	s.settle(k)
}

// lostPrimaryKey follows a statement that drops the table's primary key:
// RequirePrimaryKey follows the table when the change made it, or when the
// drop is the first thing the change does to its key, so the table had
// one before.
func (s *constraintScan) lostPrimaryKey(st *statement, v *ast.AlterTableStmt, cmd *ast.AlterTableCmd, k tableKey) {
	t := s.state(k)
	added := s.anyState(k, func(o *tableState) bool { return o.pkFirst == "add" && !o.created })
	if t.pkFirst == "" && !added {
		t.pkFirst = "drop"
	}
	had := t.created || t.pkFirst == "drop" && !added
	t.pkDropped = true
	if s.on[review.RequirePrimaryKey] && had {
		s.pending = append(s.pending, &pendingTable{key: k, statement: st.index, rng: rangeOf(cmd.Loc), message: "removes the primary key of " + objectName(v.Name) + ", and the change adds none back"})
	}
}

// reportPending reports the tables the change ends without a primary key.
func (s *constraintScan) reportPending() {
	for _, p := range s.pending {
		if !p.settled {
			s.r.report(review.RequirePrimaryKey, p.statement, p.rng, p.message)
		}
	}
}
