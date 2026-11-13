package catalog

import (
	"strings"

	nodes "github.com/bytebase/omni/oceanbase/ast"
	obparser "github.com/bytebase/omni/oceanbase/parser"
)

// checkTTL refuses a TTL (col + INTERVAL n unit, ...) option whose column tbl lacks
// (10502) or is not DATETIME or TIMESTAMP (10503), entry by entry in order, in
// CREATE TABLE and ALTER TABLE alike (TestOracle_OceanBaseTTL). text is the option's
// value as the parser keeps it, the parenthesized list.
func checkTTL(tbl *Table, text string) error {
	inner := strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(text), "("), ")")
	list, err := obparser.Parse("SELECT " + inner)
	if err != nil || list == nil || len(list.Items) != 1 {
		return nil // the parser took the option; a form this cannot read is left alone
	}
	sel, ok := list.Items[0].(*nodes.SelectStmt)
	if !ok {
		return nil
	}
	for _, target := range sel.TargetList {
		if rt, ok := target.(*nodes.ResTarget); ok {
			target = rt.Val
		}
		names := exprColumnNames(target)
		if len(names) == 0 {
			continue
		}
		col := tbl.GetColumn(names[0])
		if col == nil {
			return &Error{Code: 10502, SQLState: "HY000", Message: "TTL column '" + names[0] + "' not exists"}
		}
		switch strings.ToLower(col.DataType) {
		case "datetime", "timestamp":
		default:
			return &Error{Code: 10503, SQLState: "HY000", Message: "Column type of '" + col.Name + "' is not supported for TTL definition"}
		}
	}
	return nil
}
