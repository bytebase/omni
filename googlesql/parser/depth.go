package parser

import "github.com/bytebase/omni/googlesql/ast"

// maxNestingDepth bounds how deeply the parser may descend into nested
// constructs: expressions, queries, join groups, types, script blocks, and
// the other productions that recurse. Without a bound, a few hundred
// kilobytes of '(' overflowed the Go stack, and a stack overflow kills the
// process instead of returning an error. The Spanner emulator (ZetaSQL)
// itself gives up on 3,000 nested subqueries and on 10,000 nested unary
// minuses ("Out of stack space due to deeply nested query expression").
// It is a variable only so tests that nest deeper on purpose can raise it.
var maxNestingDepth = 10000

// nestingErrorMessage is the message of the error for a statement past
// maxNestingDepth or maxTreeHeight.
const nestingErrorMessage = "statement is nested too deeply"

// enterNesting is called on entry to the productions that recurse, placed so
// that every recursion cycle in the parser passes through one. It fails when
// the statement nests deeper than maxNestingDepth; otherwise the caller
// defers leaveNesting.
func (p *Parser) enterNesting() error {
	if p.depth >= maxNestingDepth {
		return nestedTooDeeply(p.cur.Loc)
	}
	p.depth++
	return nil
}

func (p *Parser) leaveNesting() { p.depth-- }

// nestedTooDeeply is the error for a statement past maxNestingDepth or
// maxTreeHeight, at the token at loc.
func nestedTooDeeply(loc ast.Loc) *ParseError {
	return &ParseError{Message: nestingErrorMessage, Position: loc.Start, End: loc.End}
}

// isNestingError reports whether err is the error enterNesting returns.
func isNestingError(err error) bool {
	pe, ok := err.(*ParseError)
	return ok && pe.Message == nestingErrorMessage
}

// maxTreeHeight bounds how deep the finished tree of one statement may be. A
// left-deep chain (1 + 1 + ..., a OR b OR ..., SELECT 1 UNION ALL SELECT 1
// ..., t JOIN t ...) is parsed by a loop, so it never reaches
// maxNestingDepth, yet each term adds a level to the tree, and the code that
// walks a tree recursively (ast.Walk, the query span analysis) overflowed the
// stack: Parse itself did, filling subqueries, at 800,000 terms of 1+1+...
// The Spanner emulator rejects a 10,000-term sum itself ("Out of stack
// space"), but it keeps AND and OR chains flat and runs 100,001 terms of
// either, which this limit rejects.
const maxTreeHeight = 100000

// treeTooHigh reports whether stmt, parsed from textLen bytes, is more than
// maxTreeHeight nodes deep. No construct adds more than 16 levels of tree
// per byte of text (TestTreeHeightPerByte checks this), so a shorter
// statement cannot be too deep and is not walked. The walk keeps its own
// stack.
func treeTooHigh(stmt ast.Node, textLen int) bool {
	if textLen < maxTreeHeight/16 {
		return false
	}
	return treeHeightExceeds(stmt, maxTreeHeight)
}

// treeHeightExceeds reports whether the tree under root is more than limit
// nodes deep, root counting as one.
func treeHeightExceeds(root ast.Node, limit int) bool {
	type level struct {
		node  ast.Node
		depth int
	}
	stack := []level{{root, 1}}
	var kids childLister
	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if top.depth > limit {
			return true
		}
		kids.list(top.node)
		for _, child := range kids.children {
			stack = append(stack, level{child, top.depth + 1})
		}
	}
	return false
}

// childLister is a Visitor that lets Walk enter only the node it starts at,
// so Walk lists that node's children without descending into them.
type childLister struct {
	entered  bool
	children []ast.Node
}

// list sets children to n's direct children, in the order Walk visits them.
func (c *childLister) list(n ast.Node) {
	c.entered, c.children = false, c.children[:0]
	ast.Walk(c, n)
}

func (c *childLister) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}
	if !c.entered {
		c.entered = true
		return c
	}
	c.children = append(c.children, n)
	return nil
}
