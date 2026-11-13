package v1review

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnDisallowChangeType judges top-level ALTER TABLE statements for
// CHANGE and MODIFY COLUMN actions that give a column a type other than the
// one schema holds for it. The column is the old name of CHANGE COLUMN and
// the defined name of MODIFY COLUMN, looked up in the schema named "" by the
// table's name without a qualifier, column names compared without regard to
// case; one schema lacks is not judged, and neither is an action whose
// definition has no type. An ALTER TABLE with no table is skipped, and so are
// statements nested in a routine body.
//
// The new type is rendered as its lower-cased name, then "(length)" or
// "(length,scale)" when a positive length was given (a scale only with a
// length), then " unsigned"; ZEROFILL, character sets, and ENUM and SET
// values are left out. Both that and the schema's type are compared after
// lower-casing, with a bare TINYINT, SMALLINT, MEDIUMINT, INT or BIGINT,
// signed or unsigned, given its default display width; nothing else is
// normalized, so BOOL, or a bare DECIMAL, against the schema's expanded
// spelling differs.
//
// Each finding names the column and its table and anchors the action.
// Findings are in statement order, then action order. Ported from Bytebase's
// V1 rule COLUMN_DISALLOW_CHANGE_TYPE.
func ColumnDisallowChangeType(stmts []ast.Node, schema v1review.Schema) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		n, ok := stmt.(*ast.AlterTableStmt)
		if !ok || n.Table == nil {
			continue
		}
		table := n.Table.Name
		for _, cmd := range n.Commands {
			column, ok := changedTypeColumn(cmd)
			if !ok {
				continue
			}
			held := schema.TableColumn("", table, column, false)
			if held == nil {
				continue
			}
			if normalizeColumnType(held.GetType()) != normalizeColumnType(renderColumnType(cmd.Column.TypeName)) {
				start, end := span(cmd.Loc)
				findings = append(findings, v1review.Finding{Code: v1review.ColumnTypeChanged, Name: column, Table: table, Start: start, End: end})
			}
		}
	}
	return findings
}

// changedTypeColumn returns the column a CHANGE or MODIFY COLUMN action gives
// a type, and whether the action is one whose type ColumnDisallowChangeType
// judges.
func changedTypeColumn(cmd *ast.AlterTableCmd) (string, bool) {
	if cmd == nil || cmd.Column == nil || cmd.Column.TypeName == nil {
		return "", false
	}
	switch cmd.Type {
	case ast.ATChangeColumn:
		return cmd.Name, true
	case ast.ATModifyColumn:
		return cmd.Column.Name, true
	default:
		return "", false
	}
}

func renderColumnType(dt *ast.DataType) string {
	tp := strings.ToLower(dt.Name)
	if dt.Length > 0 && dt.Scale > 0 {
		tp = fmt.Sprintf("%s(%d,%d)", tp, dt.Length, dt.Scale)
	} else if dt.Length > 0 {
		tp = fmt.Sprintf("%s(%d)", tp, dt.Length)
	}
	if dt.Unsigned {
		tp += " unsigned"
	}
	return tp
}

func normalizeColumnType(tp string) string {
	switch tp = strings.ToLower(tp); tp {
	case "tinyint":
		return "tinyint(4)"
	case "tinyint unsigned":
		return "tinyint(4) unsigned"
	case "smallint":
		return "smallint(6)"
	case "smallint unsigned":
		return "smallint(6) unsigned"
	case "mediumint":
		return "mediumint(9)"
	case "mediumint unsigned":
		return "mediumint(9) unsigned"
	case "int":
		return "int(11)"
	case "int unsigned":
		return "int(11) unsigned"
	case "bigint":
		return "bigint(20)"
	case "bigint unsigned":
		return "bigint(20) unsigned"
	default:
		return tp
	}
}
