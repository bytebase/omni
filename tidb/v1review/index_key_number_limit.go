package v1review

import (
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// IndexKeyNumberLimit reports the keys top-level statements define with more
// than limit key parts: the table-level PRIMARY KEY, UNIQUE, INDEX and
// FOREIGN KEY definitions of a CREATE TABLE, those an ALTER TABLE adds, and
// the index a CREATE INDEX of any kind creates, FULLTEXT, SPATIAL and VECTOR
// included. A key a column definition declares, a FULLTEXT, SPATIAL or VECTOR
// definition, a key inside ALTER TABLE's parenthesized ADD form, and
// statements nested in a routine body are not judged. A limit of 0 or less
// reports nothing.
//
// A PRIMARY KEY, UNIQUE or INDEX definition and a CREATE INDEX count every
// key part, expressions included; a FOREIGN KEY counts the parts of its own
// column list that are plain columns.
//
// Name is the key's name: the CONSTRAINT symbol, or else the name after
// INDEX, KEY, UNIQUE or FOREIGN KEY, or else PRIMARY for a primary key, which
// the AST keeps no other name for; it is empty for any other key with no
// name. Table is the table's name without its database qualifier, and Value
// the kind of key, KeyIndex for every CREATE INDEX.
//
// A CREATE TABLE finding anchors the key definition. Every other finding has
// no range: V1 anchored it at the first character of the statement's text
// that is not white space, which a leading comment moves off the statement.
// Findings are in statement order. Ported from Bytebase's V1 rule
// INDEX_KEY_NUMBER_LIMIT.
func IndexKeyNumberLimit(stmts []ast.Node, limit int) []v1review.Finding {
	if limit <= 0 {
		return nil
	}
	var findings []v1review.Finding
	forEachKey(stmts, func(key keyDefinition) {
		var count int
		if key.index != nil {
			count = len(key.index.Columns)
		} else {
			switch key.constraint.Type {
			case ast.ConstrPrimaryKey, ast.ConstrUnique, ast.ConstrIndex:
				count = len(key.constraint.IndexColumns)
			case ast.ConstrForeignKey:
				count = len(key.constraint.Columns)
			default:
			}
		}
		if count > limit {
			findings = append(findings, v1review.Finding{Code: v1review.IndexKeyCountExceeded, Name: key.name, Table: key.table, Value: key.kind, Start: key.start, End: key.end})
		}
	})
	return findings
}

// keyDefinition is one key a statement defines, as the V1 index rules read
// it.
type keyDefinition struct {
	// start and end are a finding's range: the definition in a CREATE
	// TABLE, none otherwise.
	start, end int
	// table is the table's name without its database qualifier.
	table string
	// name is the name V1 printed for the key.
	name string
	// kind is one of the Key constants, empty for a table-level definition
	// that is not a PRIMARY KEY, UNIQUE, INDEX or FOREIGN KEY.
	kind string
	// constraint is the table-level definition, nil for CREATE INDEX.
	constraint *ast.Constraint
	// index is the CREATE INDEX statement, nil for a table-level
	// definition.
	index *ast.CreateIndexStmt
}

// forEachKey calls visit for each key the top-level statements define at
// table level, in statement order: a CREATE TABLE's constraints, those an
// ALTER TABLE adds with ADD CONSTRAINT or ADD INDEX, and a CREATE INDEX's
// index.
func forEachKey(stmts []ast.Node, visit func(keyDefinition)) {
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.Table == nil {
				continue
			}
			for _, c := range n.Constraints {
				if c == nil {
					continue
				}
				start, end := span(c.Loc)
				visit(keyDefinition{start: start, end: end, table: n.Table.Name, name: constraintKeyName(c), kind: constraintKind(c.Type), constraint: c})
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
					visit(keyDefinition{table: n.Table.Name, name: constraintKeyName(cmd.Constraint), kind: constraintKind(cmd.Constraint.Type), constraint: cmd.Constraint})
				}
			}
		case *ast.CreateIndexStmt:
			if n.Table == nil {
				continue
			}
			visit(keyDefinition{table: n.Table.Name, name: n.IndexName, kind: v1review.KeyIndex, index: n})
		}
	}
}

// constraintKeyName returns the name V1 printed for a table-level
// definition: the CONSTRAINT symbol, else the name the AST keeps beside it,
// else PRIMARY for a primary key.
func constraintKeyName(c *ast.Constraint) string {
	switch {
	case c.Name != "":
		return c.Name
	case c.IndexName != "":
		return c.IndexName
	case c.Type == ast.ConstrPrimaryKey:
		return "PRIMARY"
	default:
		return ""
	}
}

// constraintKind returns the Key constant for a table-level definition,
// empty for any that is not a PRIMARY KEY, UNIQUE, INDEX or FOREIGN KEY.
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
