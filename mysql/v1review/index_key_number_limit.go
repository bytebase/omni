package v1review

import (
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// IndexKeyNumberLimit reports the keys top-level statements define with more
// than limit columns: the table-level PRIMARY KEY, UNIQUE, INDEX and FOREIGN
// KEY definitions of a CREATE TABLE, those an ALTER TABLE adds, and the index
// a CREATE INDEX creates unless it is FULLTEXT or SPATIAL. A key a column
// definition declares, a FULLTEXT or SPATIAL definition, a key inside ALTER
// TABLE's parenthesized ADD form, and statements nested in a routine body
// are not judged. A limit of 0 or less reports nothing.
//
// A definition's columns are its key parts that are plain columns, so an
// expression key part is not counted; CREATE INDEX also counts a
// parenthesized key part that is one column.
//
// Name is the key's name as parsed, the empty name when it has none, and
// Table the table's name without its database qualifier. Value is the kind
// of key. The AST does not keep a name written after PRIMARY KEY, so that
// key's Name is empty unless CONSTRAINT names it; the V1 advisor then read a
// name out of the statement text, which the caller keeps.
//
// A CREATE TABLE finding anchors the key definition; every other finding
// anchors the whole statement. Findings are in statement order. Ported from
// Bytebase's V1 rule INDEX_KEY_NUMBER_LIMIT.
func IndexKeyNumberLimit(stmts []ast.Node, limit int) []v1review.Finding {
	var findings []v1review.Finding
	forEachKey(stmts, func(key keyDefinition) {
		if key.kind == "" || limit <= 0 || len(key.columns) <= limit {
			return
		}
		start, end := span(key.loc)
		findings = append(findings, v1review.Finding{Code: v1review.IndexKeyCountExceeded, Name: key.name, Table: key.table, Value: key.kind, Start: start, End: end})
	})
	return findings
}

// keyDefinition is one key a statement defines, as the V1 index rules read
// it.
type keyDefinition struct {
	// loc is what a finding about the key anchors: the definition in a
	// CREATE TABLE, the whole statement otherwise.
	loc ast.Loc
	// table is the table's name without its database qualifier.
	table string
	// name is the key's name: the constraint's name, or else the name
	// after FOREIGN KEY.
	name string
	// kind is one of the Key constants, empty for a FULLTEXT or SPATIAL
	// key and for a CHECK constraint.
	kind string
	// columns are the key parts that are plain columns. For CREATE INDEX a
	// parenthesized key part that is one column is among them.
	columns []string
	// constraint is the table-level definition, nil for CREATE INDEX.
	constraint *ast.Constraint
	// index is the CREATE INDEX statement, nil for a table-level
	// definition.
	index *ast.CreateIndexStmt
}

// forEachKey calls visit for each key the top-level statements define at
// table level, in statement order: a CREATE TABLE's constraints, those an
// ALTER TABLE adds with ADD, and a CREATE INDEX's index.
func forEachKey(stmts []ast.Node, visit func(keyDefinition)) {
	constraintKey := func(loc ast.Loc, table string, c *ast.Constraint) keyDefinition {
		name := c.Name
		if name == "" {
			name = c.IndexName
		}
		return keyDefinition{loc: loc, table: table, name: name, kind: constraintKind(c.Type), columns: c.Columns, constraint: c}
	}
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.Table == nil {
				continue
			}
			for _, c := range n.Constraints {
				if c != nil {
					visit(constraintKey(c.Loc, n.Table.Name, c))
				}
			}
		case *ast.AlterTableStmt:
			if n.Table == nil {
				continue
			}
			for _, cmd := range n.Commands {
				if cmd == nil || cmd.Constraint == nil {
					continue
				}
				if cmd.Type == ast.ATAddConstraint || cmd.Type == ast.ATAddIndex {
					visit(constraintKey(n.Loc, n.Table.Name, cmd.Constraint))
				}
			}
		case *ast.CreateIndexStmt:
			if n.Table == nil {
				continue
			}
			key := keyDefinition{loc: n.Loc, table: n.Table.Name, name: n.IndexName, columns: keyPartColumns(n.Columns), index: n}
			switch {
			case n.Fulltext || n.Spatial:
			case n.Unique:
				key.kind = v1review.KeyUniqueIndex
			default:
				key.kind = v1review.KeyIndex
			}
			visit(key)
		}
	}
}

// constraintKind returns the Key constant for a table-level definition,
// empty for a FULLTEXT or SPATIAL key and for a CHECK constraint.
func constraintKind(t ast.ConstraintType) string {
	switch t {
	case ast.ConstrPrimaryKey:
		return v1review.KeyPrimary
	case ast.ConstrUnique:
		return v1review.KeyUnique
	case ast.ConstrIndex:
		return v1review.KeyIndex
	case ast.ConstrForeignKey:
		return v1review.KeyForeign
	default:
		return ""
	}
}
