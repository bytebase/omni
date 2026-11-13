package v1review

import (
	"slices"
	"strconv"
	"strings"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// The most columns of one table that may default to, or update to, the
// current time.
const (
	maxDefaultCurrentTimeColumns  = 2
	maxOnUpdateCurrentTimeColumns = 1
)

// ColumnCurrentTimeCountLimit judges the statements together, counting per
// table the DATETIME and TIMESTAMP columns whose DEFAULT, and separately
// whose ON UPDATE, is CURRENT_TIMESTAMP, NOW, LOCALTIME or LOCALTIMESTAMP
// called as a function. A table with more than 2 of the first is reported
// ColumnDefaultCurrentTimeCount, and one with more than 1 of the second
// ColumnOnUpdateCurrentTimeCount; Value is the count.
//
// The columns counted are those of top-level statements: every column a
// CREATE TABLE defines and the columns ALTER TABLE's ADD, CHANGE and MODIFY
// COLUMN write. Counts only grow: a CREATE TABLE does not restart its
// table's, DROP TABLE does not clear them, and a column changed twice is
// counted twice. Tables are keyed by name without a qualifier, compared
// bytewise.
//
// A table's findings are anchored at the whole statement that held the last
// column counted for it, and Statement says which. Tables are ordered by
// that anchor, the default finding before the on-update one. Two tables
// anchored in different statements on one line came out in either order
// from V1, which sorted on the line alone; here the earlier statement comes
// first. Ported from Bytebase's V1 rule COLUMN_CURRENT_TIME_COUNT_LIMIT.
func ColumnCurrentTimeCountLimit(stmts []ast.Node) []v1review.Finding {
	type tableCount struct {
		name              string
		defaults, updates int
		statement         int
		loc               ast.Loc
	}
	tables := make(map[string]*tableCount)
	count := func(i int, loc ast.Loc, table string, col *ast.ColumnDef) {
		if col == nil || !isTimeType(col.TypeName) {
			return
		}
		isDefault, isUpdate := isCurrentTimeCall(col.DefaultValue), isCurrentTimeCall(col.OnUpdate)
		if !isDefault && !isUpdate {
			return
		}
		t, ok := tables[table]
		if !ok {
			t = &tableCount{name: table}
			tables[table] = t
		}
		if isDefault {
			t.defaults++
		}
		if isUpdate {
			t.updates++
		}
		t.statement, t.loc = i, loc
	}
	for i, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.Table == nil {
				continue
			}
			for _, col := range n.Columns {
				count(i, n.Loc, n.Table.Name, col)
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
					count(i, n.Loc, n.Table.Name, col)
				}
			}
		}
	}

	var ordered []*tableCount
	for _, t := range tables {
		ordered = append(ordered, t)
	}
	slices.SortFunc(ordered, func(a, b *tableCount) int {
		if a.statement != b.statement {
			return a.statement - b.statement
		}
		return a.loc.Start - b.loc.Start
	})
	var findings []v1review.Finding
	for _, t := range ordered {
		start, end := span(t.loc)
		finding := v1review.Finding{Table: t.name, Start: start, End: end, Statement: t.statement}
		if t.defaults > maxDefaultCurrentTimeColumns {
			finding.Code, finding.Value = v1review.ColumnDefaultCurrentTimeCount, strconv.Itoa(t.defaults)
			findings = append(findings, finding)
		}
		if t.updates > maxOnUpdateCurrentTimeColumns {
			finding.Code, finding.Value = v1review.ColumnOnUpdateCurrentTimeCount, strconv.Itoa(t.updates)
			findings = append(findings, finding)
		}
	}
	return findings
}

func isTimeType(dt *ast.DataType) bool {
	if dt == nil {
		return false
	}
	switch strings.ToUpper(dt.Name) {
	case "DATETIME", "TIMESTAMP":
		return true
	default:
		return false
	}
}

// isCurrentTimeCall reports whether expr calls CURRENT_TIMESTAMP or one of
// its synonyms, compared without regard to case.
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
