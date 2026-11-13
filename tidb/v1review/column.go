package v1review

import (
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// alteredColumns returns the column definitions an ALTER TABLE action
// writes: the columns ADD COLUMN adds, and the new definition of CHANGE
// COLUMN and MODIFY COLUMN. Every other action writes none.
func alteredColumns(cmd *ast.AlterTableCmd) []*ast.ColumnDef {
	switch cmd.Type {
	case ast.ATAddColumn:
		return addedColumns(cmd)
	case ast.ATChangeColumn, ast.ATModifyColumn:
		if cmd.Column != nil {
			return []*ast.ColumnDef{cmd.Column}
		}
	default:
	}
	return nil
}

// columnFinding is a finding naming a column of a table, anchored at loc.
func columnFinding(code v1review.Code, table, column string, loc ast.Loc) v1review.Finding {
	start, end := span(loc)
	return v1review.Finding{Code: code, Name: column, Table: table, Start: start, End: end}
}

// columnRule reports the columns top-level CREATE TABLE and ALTER TABLE
// statements define that bad holds for: every column a CREATE TABLE
// defines, anchored at the column definition, and the columns ADD, CHANGE
// and MODIFY COLUMN write, anchored at the whole ALTER TABLE statement. A
// statement naming no table is skipped.
func columnRule(stmts []ast.Node, code v1review.Code, bad func(*ast.ColumnDef) bool) []v1review.Finding {
	return columnValueRule(stmts, code, func(col *ast.ColumnDef) (string, bool) {
		return "", bad(col)
	})
}

// columnValueRule is columnRule for a rule whose findings carry a Value,
// which judge returns with its verdict.
func columnValueRule(stmts []ast.Node, code v1review.Code, judge func(*ast.ColumnDef) (string, bool)) []v1review.Finding {
	var findings []v1review.Finding
	walkColumns(stmts, func(table string, col *ast.ColumnDef, loc ast.Loc) {
		if value, bad := judge(col); bad {
			finding := columnFinding(code, table, col.Name, loc)
			finding.Value = value
			findings = append(findings, finding)
		}
	})
	return findings
}

// walkColumns calls visit, in order, for every non-nil column a top-level
// CREATE TABLE defines, with the column definition's range, and every one
// ALTER TABLE's ADD, CHANGE and MODIFY COLUMN write, with the whole
// statement's range. A statement naming no table is skipped.
func walkColumns(stmts []ast.Node, visit func(table string, col *ast.ColumnDef, loc ast.Loc)) {
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.Table == nil {
				continue
			}
			for _, col := range n.Columns {
				if col != nil {
					visit(n.Table.Name, col, col.Loc)
				}
			}
		case *ast.AlterTableStmt:
			if n.Table == nil {
				continue
			}
			for _, cmd := range n.Commands {
				if cmd == nil {
					continue
				}
				for _, col := range alteredColumns(cmd) {
					if col != nil {
						visit(n.Table.Name, col, n.Loc)
					}
				}
			}
		}
	}
}

// hasColumnConstraint reports whether a column carries a constraint of the
// type.
func hasColumnConstraint(col *ast.ColumnDef, typ ast.ColumnConstraintType) bool {
	for _, c := range col.Constraints {
		if c != nil && c.Type == typ {
			return true
		}
	}
	return false
}
