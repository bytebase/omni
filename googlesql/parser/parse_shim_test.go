package parser

import "github.com/bytebase/omni/googlesql/ast"

// parseForTest is the pre-ParseErrors shape of Parse, kept for the test
// corpus that inspects the error list directly. Parse itself returns error.
func parseForTest(input string) (*ast.File, []ParseError) {
	file, err := Parse(input)
	return file, AllErrors(err)
}
