package catalog

import (
	"strings"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// checkTTL refuses a TTL (col + INTERVAL n unit, ...) option whose column tbl
// lacks (10502) or is not DATETIME or TIMESTAMP (10503), entry by entry in order,
// in CREATE TABLE and ALTER TABLE alike (TestOracle_OceanBaseTTL). The columns are
// the ones the parser read (TableOption.TTLColumns), so what the option's text
// holds between them, a comment among it, changes nothing.
func checkTTL(tbl *Table, opt *nodes.TableOption) error {
	for _, name := range opt.TTLColumns {
		col := tbl.GetColumn(name)
		if col == nil {
			return &Error{Code: 10502, SQLState: "HY000", Message: "TTL column '" + name + "' not exists"}
		}
		switch strings.ToLower(col.DataType) {
		case "datetime", "timestamp":
		default:
			return &Error{Code: 10503, SQLState: "HY000", Message: "Column type of '" + col.Name + "' is not supported for TTL definition"}
		}
	}
	return nil
}
