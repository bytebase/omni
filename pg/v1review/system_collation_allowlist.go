package v1review

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// SystemCollationAllowlist judges the collations column definitions name
// against the allowed list, compared exactly: letter case counts, and the
// list is not lower-cased. Judged are the columns a top-level CREATE TABLE
// defines and those an ALTER TABLE gives with ADD COLUMN or ALTER COLUMN
// TYPE, in the order written. Only the last part of a qualified collation
// name is judged, as the parser stores it: an unquoted name is folded to
// lower case. A COLLATE in an expression, a table constraint, a CREATE
// DATABASE, a domain, an index, or a statement under EXPLAIN or in a routine
// body is not judged.
//
// Name is the collation's last name part. Every finding anchors the whole
// statement. Ported from Bytebase's V1 rule SYSTEM_COLLATION_ALLOWLIST.
func SystemCollationAllowlist(stmts []ast.Node, allowed []string) []v1review.Finding {
	allowedSet := make(map[string]bool, len(allowed))
	for _, collation := range allowed {
		allowedSet[collation] = true
	}
	var findings []v1review.Finding
	check := func(loc ast.Loc, col *ast.ColumnDef) {
		if col == nil || col.CollClause == nil {
			return
		}
		parts := names(col.CollClause.Collname)
		if len(parts) == 0 {
			return
		}
		if name := parts[len(parts)-1]; name != "" && !allowedSet[name] {
			start, end := span(loc)
			findings = append(findings, v1review.Finding{Code: v1review.CollationNotAllowed, Name: name, Start: start, End: end})
		}
	}
	for _, stmt := range stmts {
		if raw, ok := stmt.(*ast.RawStmt); ok {
			stmt = raw.Stmt
		}
		switch n := stmt.(type) {
		case *ast.CreateStmt:
			if n.TableElts == nil {
				continue
			}
			for _, item := range n.TableElts.Items {
				if col, ok := item.(*ast.ColumnDef); ok {
					check(n.Loc, col)
				}
			}
		case *ast.AlterTableStmt:
			if n.Cmds == nil {
				continue
			}
			for _, item := range n.Cmds.Items {
				cmd, ok := item.(*ast.AlterTableCmd)
				if !ok {
					continue
				}
				switch ast.AlterTableType(cmd.Subtype) {
				case ast.AT_AddColumn, ast.AT_AlterColumnType:
					if col, ok := cmd.Def.(*ast.ColumnDef); ok {
						check(n.Loc, col)
					}
				}
			}
		}
	}
	return findings
}
