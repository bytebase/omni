package v1review

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// IndexKeyNumberLimit reports the keys top-level statements define with more
// than limit columns: the index a CREATE INDEX creates, and the table-level
// PRIMARY KEY, UNIQUE and FOREIGN KEY constraints a CREATE TABLE defines or
// an ALTER TABLE ... ADD adds. A key a column definition declares, an
// EXCLUDE or CHECK constraint, and a statement in a CREATE SCHEMA or in a
// routine body are not judged. A limit of 0 or less reports nothing.
//
// A CREATE INDEX's columns are its plain column elements, so an expression
// element is not counted, nor is an INCLUDE column; a constraint's are the
// columns it lists, a FOREIGN KEY's its referencing columns, and a key built
// USING INDEX has none.
//
// Name is the index's or the constraint's name as the parser stores it, the
// empty name when it has none, and Table the table's name without its
// schema. Value is the kind of key: a CREATE INDEX is KeyIndex, unique or
// not. Every finding anchors the whole statement; Bytebase's V1 advisor
// reported the statement text's first non-blank line. Findings are in
// statement order. Ported from Bytebase's V1 rule INDEX_KEY_NUMBER_LIMIT.
func IndexKeyNumberLimit(stmts []ast.Node, limit int) []v1review.Finding {
	var findings []v1review.Finding
	forEachKey(stmts, func(key keyDefinition) {
		if limit <= 0 || len(key.columns) <= limit {
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
	// loc is the whole statement's.
	loc ast.Loc
	// table is the table's name without its schema.
	table string
	// name is the index's or the constraint's name.
	name string
	// kind is one of the Key constants.
	kind string
	// columns are the plain columns a CREATE INDEX lists, or the columns a
	// constraint lists.
	columns []string
}

// forEachKey calls visit for each key the top-level statements define, in
// statement order: a CREATE INDEX's index, and the PRIMARY KEY, UNIQUE and
// FOREIGN KEY table constraints of a CREATE TABLE and of an ALTER TABLE's
// ADD actions.
func forEachKey(stmts []ast.Node, visit func(keyDefinition)) {
	constraintKey := func(loc ast.Loc, table string, c *ast.Constraint) {
		key := keyDefinition{loc: loc, table: table, name: c.Conname}
		switch c.Contype {
		case ast.CONSTR_PRIMARY:
			key.kind, key.columns = v1review.KeyPrimary, names(c.Keys)
		case ast.CONSTR_UNIQUE:
			key.kind, key.columns = v1review.KeyUnique, names(c.Keys)
		case ast.CONSTR_FOREIGN:
			key.kind, key.columns = v1review.KeyForeign, names(c.FkAttrs)
		default:
			return
		}
		visit(key)
	}
	for _, stmt := range stmts {
		if raw, ok := stmt.(*ast.RawStmt); ok {
			stmt = raw.Stmt
		}
		switch n := stmt.(type) {
		case *ast.IndexStmt:
			visit(keyDefinition{loc: n.Loc, table: relname(n.Relation), name: n.Idxname, kind: v1review.KeyIndex, columns: indexColumns(n)})
		case *ast.CreateStmt:
			if n.TableElts == nil {
				continue
			}
			for _, item := range n.TableElts.Items {
				if c, ok := item.(*ast.Constraint); ok {
					constraintKey(n.Loc, relname(n.Relation), c)
				}
			}
		case *ast.AlterTableStmt:
			if n.Cmds == nil {
				continue
			}
			for _, item := range n.Cmds.Items {
				cmd, ok := item.(*ast.AlterTableCmd)
				if !ok || ast.AlterTableType(cmd.Subtype) != ast.AT_AddConstraint {
					continue
				}
				if c, ok := cmd.Def.(*ast.Constraint); ok {
					constraintKey(n.Loc, relname(n.Relation), c)
				}
			}
		}
	}
}
