package v1review

import (
	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/v1review"
)

// TableNoForeignKey judges each table as the top-level statements leave it,
// reporting after reading them all; statements nested in a routine body or
// an IF or BEGIN block are not read. Tables are told apart by their name as
// written, without brackets, a missing schema read as dbo, case-sensitively.
//
// A CREATE TABLE gives its table no foreign key, then one if a column's
// constraint or a table constraint is a FOREIGN KEY; a REFERENCES clause
// written on a column counts. Each ADD CONSTRAINT ... FOREIGN KEY action of
// an ALTER TABLE gives its table one; ADD of a column with REFERENCES does
// not. Nothing removes one, neither dropping the constraint nor dropping
// the table.
//
// Each table with a foreign key is one TableForeignKey. Name is the table as
// written without brackets, with the qualifiers it was written with, in its
// last CREATE TABLE, else in the first ALTER TABLE that added it a foreign
// key. The finding anchors the foreign key that last gave the table one: a
// CREATE TABLE's first FOREIGN KEY constraint, column constraints before
// table constraints, or an ALTER TABLE's added constraint; Statement is
// that statement's index.
//
// The V1 advisor reported tables in Go map order, which varies from run to
// run; the findings come in the order the tables were first created or
// given a foreign key. Ported from Bytebase's V1 rule TABLE_NO_FOREIGN_KEY.
func TableNoForeignKey(stmts []ast.Node) []v1review.Finding {
	tables := &tableStates{}
	for i, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.Name == nil {
				continue
			}
			t := tables.get(joinTableRef(n.Name, "dbo"))
			t.flag, t.name = false, tableRefText(n.Name)
			t.anchor(i, n.Name.Loc)
			if loc, ok := createTableForeignKey(n); ok {
				t.flag = true
				t.anchor(i, loc)
			}
		case *ast.AlterTableStmt:
			if n.Name == nil || n.Actions == nil {
				continue
			}
			key := joinTableRef(n.Name, "dbo")
			for _, item := range n.Actions.Items {
				action, ok := item.(*ast.AlterTableAction)
				if !ok || action.Type != ast.ATAddConstraint || action.Constraint == nil || action.Constraint.Type != ast.ConstraintForeignKey {
					continue
				}
				t, existed := tables.lookup(key)
				if !existed {
					t = tables.get(key)
					t.name = tableRefText(n.Name)
				}
				t.flag = true
				t.anchor(i, action.Constraint.Loc)
			}
		default:
		}
	}
	var findings []v1review.Finding
	for _, t := range tables.list {
		if t.flag {
			findings = append(findings, t.finding(v1review.TableForeignKey))
		}
	}
	return findings
}

// createTableForeignKey returns the first FOREIGN KEY constraint of a
// CREATE TABLE, column constraints before table constraints.
func createTableForeignKey(n *ast.CreateTableStmt) (ast.Loc, bool) {
	if n.Columns != nil {
		for _, item := range n.Columns.Items {
			col, ok := item.(*ast.ColumnDef)
			if !ok || col.Constraints == nil {
				continue
			}
			for _, c := range col.Constraints.Items {
				if cd, ok := c.(*ast.ConstraintDef); ok && cd.Type == ast.ConstraintForeignKey {
					return cd.Loc, true
				}
			}
		}
	}
	if n.Constraints != nil {
		for _, item := range n.Constraints.Items {
			if cd, ok := item.(*ast.ConstraintDef); ok && cd.Type == ast.ConstraintForeignKey {
				return cd.Loc, true
			}
		}
	}
	return ast.Loc{}, false
}

// tableStates holds TableNoForeignKey's and TableRequirePK's state per
// table, in the order the tables were first tracked.
type tableStates struct {
	list  []*tableState
	byKey map[string]*tableState
}

type tableState struct {
	flag       bool // has a foreign key, or a primary key
	name       string
	stmt       int
	start, end int
}

func (s *tableStates) lookup(key string) (*tableState, bool) {
	t, ok := s.byKey[key]
	return t, ok
}

// get returns the table's state, tracking it if it is new.
func (s *tableStates) get(key string) *tableState {
	if t, ok := s.byKey[key]; ok {
		return t
	}
	if s.byKey == nil {
		s.byKey = map[string]*tableState{}
	}
	t := &tableState{}
	s.byKey[key] = t
	s.list = append(s.list, t)
	return t
}

func (t *tableState) anchor(stmt int, loc ast.Loc) {
	t.stmt = stmt
	t.start, t.end = span(loc)
}

func (t *tableState) finding(code v1review.Code) v1review.Finding {
	return v1review.Finding{Code: code, Name: t.name, Start: t.start, End: t.end, Statement: t.stmt}
}

// tableRefText is a table reference's object name prefixed by its schema
// and database names where written, joined with dots and without brackets;
// db..t is db.t.
func tableRefText(ref *ast.TableRef) string {
	text := ref.Object
	if ref.Schema != "" {
		text = ref.Schema + "." + text
	}
	if ref.Database != "" {
		text = ref.Database + "." + text
	}
	return text
}
