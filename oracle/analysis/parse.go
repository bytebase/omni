package analysis

import (
	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/oracle/parser"
)

// Parse parses a script statement by statement, as SplitStatements splits
// it, and returns every statement in one list. Locations and a parse error's
// position are byte offsets in the whole script.
func Parse(sql string) (*ast.List, error) {
	list := &ast.List{}
	for _, statement := range SplitStatements(sql) {
		if statement.Empty {
			continue
		}
		parsed, err := parser.ParseRange(sql, statement.Start, statement.End)
		if err != nil {
			return nil, err
		}
		if parsed == nil {
			continue
		}
		list.Items = append(list.Items, parsed.Items...)
	}
	return list, nil
}
