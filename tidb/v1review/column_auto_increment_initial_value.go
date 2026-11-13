package v1review

import (
	"strconv"
	"strings"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnAutoIncrementInitialValue reports the AUTO_INCREMENT table options
// of top-level CREATE TABLE statements whose value is not value. The option
// name is compared without regard to case. Its value is read as a decimal
// unsigned 64-bit integer, and an option whose value does not read that way
// is not judged; value is converted to uint64 for the comparison, so a
// negative value wraps around. ALTER TABLE ... AUTO_INCREMENT is not judged.
//
// Each option is one finding, anchored at the whole statement, with Table
// the table without a qualifier and Value the option's value as the
// decimal number read. Ported from Bytebase's V1 rule
// COLUMN_AUTO_INCREMENT_INITIAL_VALUE.
func ColumnAutoIncrementInitialValue(stmts []ast.Node, value int) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		create, ok := stmt.(*ast.CreateTableStmt)
		if !ok || create.Table == nil {
			continue
		}
		for _, opt := range create.Options {
			if opt == nil || !strings.EqualFold(opt.Name, "AUTO_INCREMENT") {
				continue
			}
			initial, err := strconv.ParseUint(opt.Value, 10, 0)
			if err != nil || initial == uint64(value) {
				continue
			}
			start, end := span(create.Loc)
			findings = append(findings, v1review.Finding{
				Code:  v1review.ColumnAutoIncrementInitialValue,
				Table: create.Table.Name,
				Value: strconv.FormatUint(initial, 10),
				Start: start,
				End:   end,
			})
		}
	}
	return findings
}
