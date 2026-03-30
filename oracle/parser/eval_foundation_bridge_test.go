package parser

import "github.com/bytebase/omni/oracle/ast"

func init() {
	evalNoLoc = ast.NoLoc
	evalNodeLoc = ast.NodeLoc
	evalListSpan = ast.ListSpan
}
