package v1review

import (
	"slices"
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
	"github.com/bytebase/omni/v1review"
)

// TableNoForeignKey judges each table as the top-level statements leave it,
// reporting after reading them all. Tables are told apart by their name
// normalized as Snowflake resolves it (unquoted parts upper-cased), a
// missing schema read as PUBLIC and a missing database as none.
//
// A plain CREATE TABLE gives its table no foreign keys, then counts each
// column's inline FOREIGN KEY constraint and each FOREIGN KEY table
// constraint; CREATE TABLE ... AS SELECT, LIKE and CLONE are not read. In
// an ALTER TABLE, each ADD of a FOREIGN KEY constraint counts one more, and
// each DROP FOREIGN KEY one fewer, never below zero, for a table counted
// before; DROP CONSTRAINT by name does not, since it does not say the
// constraint's kind. A column added with an inline foreign key counts
// none, and DROP TABLE removes nothing.
//
// Each table with a count above zero is one TableForeignKey. Name is the
// table as written, quotes kept, in the last CREATE TABLE or ALTER TABLE
// read that adds, drops or renames a constraint. The finding anchors the
// last FOREIGN KEY table constraint or added foreign key read for the
// table, else its CREATE TABLE; an inline foreign key does not move it.
// Statement is the anchoring statement's index.
//
// The V1 advisor reported tables in Go map order, which varies from run to
// run; the findings come in the order the tables were first counted.
// Ported from Bytebase's V1 rule TABLE_NO_FOREIGN_KEY.
func TableNoForeignKey(stmts []ast.Node) []v1review.Finding {
	tables := newTableStates()
	for i, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.AsSelect != nil || n.Like != nil || n.Clone != nil {
				continue
			}
			key := tableStateKey(n.Name)
			t := tables.get(key)
			t.count = 0
			tables.names[key] = n.Name.String()
			t.anchor(i, n.Loc)
			for _, column := range n.Columns {
				if column.InlineConstraint != nil && column.InlineConstraint.Type == ast.ConstrForeignKey {
					t.count++
				}
			}
			for _, constraint := range n.Constraints {
				if constraint.Type == ast.ConstrForeignKey {
					t.count++
					t.anchor(i, constraint.Loc)
				}
			}
		case *ast.AlterTableStmt:
			if !alterTableHasConstraintAction(n) {
				continue
			}
			key := tableStateKey(n.Name)
			tables.names[key] = n.Name.String()
			for _, action := range n.Actions {
				switch action.Kind {
				case ast.AlterTableAddConstraint:
					if action.Constraint != nil && action.Constraint.Type == ast.ConstrForeignKey {
						t := tables.get(key)
						t.count++
						t.anchor(i, action.Constraint.Loc)
					}
				case ast.AlterTableDropConstraint:
					if t, ok := tables.byKey[key]; ok && action.DropForeignKey && t.count > 0 {
						t.count--
					}
				default:
				}
			}
		default:
		}
	}
	var findings []v1review.Finding
	for _, t := range tables.list {
		if t.count > 0 {
			findings = append(findings, tables.finding(t, v1review.TableForeignKey))
		}
	}
	return findings
}

func alterTableHasConstraintAction(stmt *ast.AlterTableStmt) bool {
	for _, action := range stmt.Actions {
		switch action.Kind {
		case ast.AlterTableAddConstraint, ast.AlterTableDropConstraint, ast.AlterTableRenameConstraint:
			return true
		default:
		}
	}
	return false
}

// tableStateKey is the "database.schema.table" key TableNoForeignKey and
// TableRequirePK tell tables apart by; a nil name is "".
func tableStateKey(name *ast.ObjectName) string {
	if name == nil {
		return ""
	}
	schema := "PUBLIC"
	if s := name.Schema.Normalize(); s != "" {
		schema = s
	}
	parts := []string{name.Database.Normalize(), schema}
	if o := name.Name.Normalize(); o != "" {
		parts = append(parts, o)
	}
	return strings.Join(parts, ".")
}

// tableStates holds TableNoForeignKey's and TableRequirePK's state per
// table, in the order the tables were first tracked, and the name each was
// last written with, which is kept for tables not tracked too.
type tableStates struct {
	list  []*tableState
	byKey map[string]*tableState
	names map[string]string
}

func newTableStates() *tableStates {
	return &tableStates{byKey: map[string]*tableState{}, names: map[string]string{}}
}

type tableState struct {
	key        string
	count      int  // foreign keys
	hasPK      bool // primary key
	stmt       int
	start, end int
}

// get returns the table's state, tracking it if it is new.
func (s *tableStates) get(key string) *tableState {
	if t, ok := s.byKey[key]; ok {
		return t
	}
	t := &tableState{key: key}
	s.byKey[key] = t
	s.list = append(s.list, t)
	return t
}

// remove forgets the table and the name it was written with.
func (s *tableStates) remove(key string) {
	if t, ok := s.byKey[key]; ok {
		s.list = slices.DeleteFunc(s.list, func(u *tableState) bool { return u == t })
	}
	delete(s.byKey, key)
	delete(s.names, key)
}

func (s *tableStates) finding(t *tableState, code v1review.Code) v1review.Finding {
	return v1review.Finding{Code: code, Name: s.names[t.key], Start: t.start, End: t.end, Statement: t.stmt}
}

// anchor records the statement index and the start of loc; a missing end
// collapses the range to its start, and a missing start gives 0,0.
func (t *tableState) anchor(stmt int, loc ast.Loc) {
	t.stmt = stmt
	t.start, t.end = 0, 0
	if loc.Start >= 0 {
		t.start, t.end = loc.Start, max(loc.End, loc.Start)
	}
}
