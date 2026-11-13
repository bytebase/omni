package v1review

import (
	"strconv"
	"strings"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnAutoIncrementInitialValue judges top-level CREATE TABLE and ALTER
// TABLE statements for an AUTO_INCREMENT table option whose value is not
// value: every such option of a CREATE TABLE, and every table-option action
// of an ALTER TABLE, the option name matched in any case. An option value
// that does not parse as an unsigned decimal integer is not judged, and the
// comparison is against value converted to uint64, so a negative value
// matches no option. A statement whose table has no name is not judged, and
// neither are statements nested in a routine body.
//
// Each finding names the table, without its database qualifier, and holds in
// Value the option's value as parsed, in decimal. It anchors the whole
// statement. Findings are in statement order, then option order. Ported from
// Bytebase's V1 rule COLUMN_AUTO_INCREMENT_INITIAL_VALUE.
func ColumnAutoIncrementInitialValue(stmts []ast.Node, value int) []v1review.Finding {
	var findings []v1review.Finding
	check := func(loc ast.Loc, table, option string) {
		parsed, err := strconv.ParseUint(option, 10, 0)
		if err != nil || parsed == uint64(value) {
			return
		}
		start, end := span(loc)
		findings = append(findings, v1review.Finding{
			Code:  v1review.ColumnAutoIncrementInitialValue,
			Table: table,
			Value: strconv.FormatUint(parsed, 10),
			Start: start,
			End:   end,
		})
	}
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			table := refName(n.Table)
			if table == "" {
				continue
			}
			for _, opt := range n.Options {
				if strings.EqualFold(opt.Name, "AUTO_INCREMENT") {
					check(n.Loc, table, opt.Value)
				}
			}
		case *ast.AlterTableStmt:
			table := refName(n.Table)
			if table == "" {
				continue
			}
			for _, cmd := range n.Commands {
				if cmd.Type == ast.ATTableOption && cmd.Option != nil && strings.EqualFold(cmd.Option.Name, "AUTO_INCREMENT") {
					check(n.Loc, table, cmd.Option.Value)
				}
			}
		}
	}
	return findings
}
