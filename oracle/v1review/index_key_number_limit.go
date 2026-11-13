package v1review

import (
	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/v1review"
)

// IndexKeyNumberLimit reports the indexes and keys top-level statements
// define with more than limit columns: the index a CREATE INDEX creates,
// the table-level PRIMARY KEY and UNIQUE constraints of a CREATE TABLE, and
// a PRIMARY KEY or UNIQUE constraint any ALTER TABLE action carries. A
// column's own constraint, a FOREIGN KEY or CHECK constraint, and
// statements nested in a PL/SQL block are not judged. A limit of 0 or less
// reports nothing.
//
// A CREATE INDEX counts every item of its column list, expressions
// included; a constraint counts the columns it lists. An ALTER TABLE
// action that names a constraint without listing its columns, such as
// MODIFY PRIMARY KEY, counts none.
//
// IndexKeyCountExceeded is the code, and no name is set: the V1 advisor
// printed only the limit. A CREATE INDEX finding anchors the statement and
// a constraint finding anchors the constraint. Findings are in statement
// order, a statement's constraints in the order written. A statement may
// be passed bare or wrapped in an *ast.RawStmt. Ported from Bytebase's V1
// rule INDEX_KEY_NUMBER_LIMIT.
func IndexKeyNumberLimit(stmts []ast.Node, limit int) []v1review.Finding {
	if limit <= 0 {
		return nil
	}
	var findings []v1review.Finding
	check := func(columns *ast.List, loc ast.Loc) {
		if columns != nil && len(columns.Items) > limit {
			start, end := span(loc)
			findings = append(findings, v1review.Finding{Code: v1review.IndexKeyCountExceeded, Start: start, End: end})
		}
	}
	checkConstraint := func(c *ast.TableConstraint) {
		if c != nil && (c.Type == ast.CONSTRAINT_PRIMARY || c.Type == ast.CONSTRAINT_UNIQUE) {
			check(c.Columns, c.Loc)
		}
	}
	for _, stmt := range stmts {
		if raw, ok := stmt.(*ast.RawStmt); ok {
			stmt = raw.Stmt
		}
		switch n := stmt.(type) {
		case *ast.CreateIndexStmt:
			check(n.Columns, n.Loc)
		case *ast.CreateTableStmt:
			for _, c := range tableConstraints(n.Constraints) {
				checkConstraint(c)
			}
		case *ast.AlterTableStmt:
			if n.Actions == nil {
				continue
			}
			for _, item := range n.Actions.Items {
				if cmd, ok := item.(*ast.AlterTableCmd); ok && cmd.Constraint != nil {
					checkConstraint(cmd.Constraint)
				}
			}
		}
	}
	return findings
}
