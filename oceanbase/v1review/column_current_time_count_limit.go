package v1review

import (
	"slices"
	"strconv"
	"strings"

	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/v1review"
)

// The most columns of one table that ColumnCurrentTimeCountLimit allows to
// default to, and to update to, the current time.
const (
	ColumnCurrentTimeDefaultLimit  = 2
	ColumnCurrentTimeOnUpdateLimit = 1
)

// ColumnCurrentTimeCountLimit judges the statements together, counting per
// table the DATETIME and TIMESTAMP columns (type name in any case) whose
// DEFAULT, and separately whose ON UPDATE, is a call to NOW,
// CURRENT_TIMESTAMP, LOCALTIME or LOCALTIMESTAMP (name in any case). It
// counts every column a CREATE TABLE defines, and in an ALTER TABLE the
// columns of ADD COLUMN in both its forms, MODIFY COLUMN and CHANGE COLUMN.
// Counts only grow: a table is keyed by its name without the database
// qualifier, a CREATE TABLE does not reset it, a DROP TABLE does not forget
// it, and a column modified twice counts twice. A statement whose table has
// no name is not judged, and neither are statements nested in a routine
// body.
//
// A table whose default count exceeds ColumnCurrentTimeDefaultLimit yields a
// ColumnDefaultCurrentTimeCount finding, and one whose on-update count
// exceeds ColumnCurrentTimeOnUpdateLimit a ColumnOnUpdateCurrentTimeCount
// finding, the first before the second; Value holds the count. Both anchor
// the last column counted for the table, in the statement Statement names.
// Tables are in the order of those anchors, statement first and offset
// second. V1 sorted tables by the anchor's line with an unstable sort over a
// map, so two tables anchored on the same line, which only statements
// sharing a line can produce, came out in either order; this order puts the
// earlier statement first. Ported from Bytebase's V1 rule
// COLUMN_CURRENT_TIME_COUNT_LIMIT.
func ColumnCurrentTimeCountLimit(stmts []ast.Node) []v1review.Finding {
	tables := make(map[string]*currentTimeCount)
	var order []*currentTimeCount
	for i, stmt := range stmts {
		table, cols := definedColumns(stmt)
		for _, col := range cols {
			if col.TypeName == nil || !isTimeType(col.TypeName) {
				continue
			}
			def, upd := isCurrentTimeCall(col.DefaultValue), isCurrentTimeCall(col.OnUpdate)
			if !def && !upd {
				continue
			}
			count, ok := tables[table]
			if !ok {
				count = &currentTimeCount{table: table}
				tables[table] = count
				order = append(order, count)
			}
			if def {
				count.defaults++
			}
			if upd {
				count.onUpdates++
			}
			count.anchor = columnAnchor{statement: i, loc: col.Loc}
		}
	}

	slices.SortStableFunc(order, func(a, b *currentTimeCount) int {
		if a.anchor.statement != b.anchor.statement {
			return a.anchor.statement - b.anchor.statement
		}
		return a.anchor.loc.Start - b.anchor.loc.Start
	})
	var findings []v1review.Finding
	for _, count := range order {
		start, end := span(count.anchor.loc)
		finding := v1review.Finding{Table: count.table, Start: start, End: end, Statement: count.anchor.statement}
		if count.defaults > ColumnCurrentTimeDefaultLimit {
			finding.Code = v1review.ColumnDefaultCurrentTimeCount
			finding.Value = strconv.Itoa(count.defaults)
			findings = append(findings, finding)
		}
		if count.onUpdates > ColumnCurrentTimeOnUpdateLimit {
			finding.Code = v1review.ColumnOnUpdateCurrentTimeCount
			finding.Value = strconv.Itoa(count.onUpdates)
			findings = append(findings, finding)
		}
	}
	return findings
}

type currentTimeCount struct {
	table     string
	defaults  int
	onUpdates int
	anchor    columnAnchor
}

func isTimeType(dt *ast.DataType) bool {
	switch strings.ToUpper(dt.Name) {
	case "DATETIME", "TIMESTAMP":
		return true
	default:
		return false
	}
}

func isCurrentTimeCall(expr ast.ExprNode) bool {
	fc, ok := expr.(*ast.FuncCallExpr)
	if !ok {
		return false
	}
	switch strings.ToUpper(fc.Name) {
	case "NOW", "CURRENT_TIMESTAMP", "LOCALTIME", "LOCALTIMESTAMP":
		return true
	default:
		return false
	}
}
