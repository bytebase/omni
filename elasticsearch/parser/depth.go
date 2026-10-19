package parser

import "errors"

// maxNestingDepth bounds how deeply a request body may nest objects and
// arrays. Without a bound, a body nested about 3 million levels deep (16MB)
// overflowed the Go stack, and a stack overflow kills the process instead of
// returning an error. The HJSON decoder that IndentData calls stops at the
// same depth.
const maxNestingDepth = 10000

// errNestedTooDeeply becomes a SyntaxError at the offset where the body
// passes maxNestingDepth, like every other error of a request.
var errNestedTooDeeply = errors.New("statement is nested too deeply")

// enterNesting is called on entry to value, which every recursion cycle of
// the body grammar passes through. It fails when the body nests deeper than
// maxNestingDepth; otherwise the caller defers leaveNesting.
func (s *state) enterNesting() error {
	if s.depth >= maxNestingDepth {
		return errNestedTooDeeply
	}
	s.depth++
	return nil
}

func (s *state) leaveNesting() { s.depth-- }
