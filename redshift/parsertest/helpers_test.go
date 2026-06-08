package parsertest

import (
	nodes "github.com/bytebase/omni/redshift/ast"
	"github.com/bytebase/omni/redshift/parser"
)

// parse calls the recursive descent parser and unwraps RawStmt wrappers
// so that tests can assert on the inner statement types directly.
func parse(sql string) (*nodes.List, error) {
	result, err := parser.Parse(sql)
	if err != nil || result == nil {
		return result, err
	}
	for i, item := range result.Items {
		if raw, ok := item.(*nodes.RawStmt); ok {
			result.Items[i] = raw.Stmt
		}
	}
	return result, nil
}
