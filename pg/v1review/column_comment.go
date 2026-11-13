package v1review

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnComment judges the statements together for columns without a
// comment, or with one longer than maxLength bytes.
//
// The columns judged are each one a top-level CREATE TABLE defines and each
// one an ALTER ... ADD COLUMN adds, whatever relation the ALTER names; a
// statement without a relation is skipped. A column's comment is the last
// COMMENT ON COLUMN, before or after the column among all the statements,
// that names its table and column: the relation is matched by name alone,
// its schema ignored, and a COMMENT naming the column with more than three
// parts is not read. COMMENT ... IS NULL leaves an empty comment.
//
// A column whose comment is missing or empty reports ColumnCommentMissing
// when required is set, anchored at the statement defining the column. A
// column whose comment is not empty reports ColumnCommentTooLong when
// maxLength is not negative and the comment is longer, anchored at the
// COMMENT statement; Bytebase passes a negative maxLength for no limit. Each
// column reports at most once, in the order the columns are defined, and a
// column defined twice is judged twice. Name is the column and Table the
// relation's name without its schema. Statement is the index of the
// anchoring statement and Start and End its range: Bytebase's V1 advisor
// reported the line that statement starts on. Ported from Bytebase's V1
// rule COLUMN_COMMENT.
func ColumnComment(stmts []ast.Node, required bool, maxLength int) []v1review.Finding {
	type column struct {
		table, name string
		anchor      columnAnchor
	}
	type comment struct {
		table, name, text string
		anchor            columnAnchor
	}
	var columns []column
	var comments []comment
	for i, stmt := range stmts {
		switch n := statementNode(stmt).(type) {
		case *ast.CreateStmt:
			if n.Relation == nil {
				continue
			}
			for _, col := range tableColumns(n) {
				columns = append(columns, column{n.Relation.Relname, col.Colname, columnAnchor{i, n.Loc}})
			}
		case *ast.AlterTableStmt:
			if n.Relation == nil {
				continue
			}
			for _, cmd := range alterCommands(n) {
				if ast.AlterTableType(cmd.Subtype) != ast.AT_AddColumn {
					continue
				}
				if col, ok := cmd.Def.(*ast.ColumnDef); ok {
					columns = append(columns, column{n.Relation.Relname, col.Colname, columnAnchor{i, n.Loc}})
				}
			}
		case *ast.CommentStmt:
			if n.Objtype != ast.OBJECT_COLUMN {
				continue
			}
			list, ok := n.Object.(*ast.List)
			if !ok || len(list.Items) < 2 || len(list.Items) > 3 {
				continue
			}
			items := list.Items[len(list.Items)-2:]
			comments = append(comments, comment{stringValue(items[0]), stringValue(items[1]), n.Comment, columnAnchor{i, n.Loc}})
		}
	}

	var findings []v1review.Finding
	for _, col := range columns {
		var matched *comment
		for i := range comments {
			if comments[i].table == col.table && comments[i].name == col.name {
				matched = &comments[i]
			}
		}
		switch {
		case matched == nil || matched.text == "":
			if required {
				findings = append(findings, col.anchor.finding(v1review.ColumnCommentMissing, col.table, col.name))
			}
		case maxLength >= 0 && len(matched.text) > maxLength:
			findings = append(findings, matched.anchor.finding(v1review.ColumnCommentTooLong, col.table, col.name))
		}
	}
	return findings
}

// columnAnchor is the statement a cross-statement finding takes its line
// from: its index and its range.
type columnAnchor struct {
	statement int
	loc       ast.Loc
}

func (a columnAnchor) finding(code v1review.Code, table, name string) v1review.Finding {
	start, end := span(a.loc)
	return v1review.Finding{Code: code, Name: name, Table: table, Start: start, End: end, Statement: a.statement}
}
