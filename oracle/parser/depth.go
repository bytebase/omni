package parser

import (
	"errors"

	nodes "github.com/bytebase/omni/oracle/ast"
)

// maxNestingDepth bounds how deeply the parser may descend into nested
// constructs: expressions, queries, table references, PL/SQL blocks and
// declarations, and the other productions that recurse. Without a bound,
// about 700,000 nested parentheses overflowed the Go stack, and a stack
// overflow kills the process instead of returning an error. Oracle 23ai
// itself rejects 8,001 nested parentheses or function calls (ORA-03095),
// 256 nested subqueries in an expression (ORA-01778), and 99 nested PL/SQL
// blocks (PLS-00103), but takes 10,001 nested derived tables, which this
// limit rejects.
// It is a variable only so tests that nest deeper on purpose can raise it.
var maxNestingDepth = 10000

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

// nestingErrorMessage is the message of the error for a statement past
// maxNestingDepth or maxTreeHeight.
const nestingErrorMessage = "statement is nested too deeply"

// nestedTooDeeply is the error for a statement past maxNestingDepth or
// maxTreeHeight, at byte offset pos.
func nestedTooDeeply(pos int) *ParseError {
	return &ParseError{Message: nestingErrorMessage, Position: pos}
}

// nestingError returns err as the error nestedTooDeeply makes, or nil.
func nestingError(err error) *ParseError {
	var pe *ParseError
	if errors.As(err, &pe) && pe.Message == nestingErrorMessage {
		return pe
	}
	return nil
}

// maxTreeHeight bounds how deep the finished tree of one statement may be.
// A left-deep chain (1 + 1 + ..., SELECT ... UNION ALL SELECT ..., t JOIN t
// ON ... JOIN ...) is parsed by a loop, so it never reaches maxNestingDepth,
// yet each term adds a level to the tree, and the code that walks a tree
// recursively overflowed the stack: ast.NodeToString on chains of a million
// terms, ast.Walk on three million, and this parser's own evaluation of a
// $IF condition on two million.
const maxTreeHeight = 100000

// treeTooHigh reports whether the tree under root, parsed from textLen
// bytes, is more than maxTreeHeight nodes deep. No construct adds more than
// 16 levels of tree per byte of text (TestTreeHeightPerByte checks this), so
// a shorter text cannot be too deep and is not walked.
func treeTooHigh(root nodes.Node, textLen int) bool {
	if textLen < maxTreeHeight/16 {
		return false
	}
	return treeHeightExceeds(root, maxTreeHeight)
}

// treeHeightExceeds reports whether the tree under root is more than limit
// nodes deep, root counting as one. A level is a node ast.Walk visits, a
// List among them, so the height is the depth Walk's recursion reaches. It
// keeps its own stack.
func treeHeightExceeds(root nodes.Node, limit int) bool {
	type level struct {
		node  nodes.Node
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
		for _, child := range kids.list(top.node) {
			stack = append(stack, level{child, top.depth + 1})
		}
	}
	return false
}

// inspect calls f for each node of the tree under root in depth-first
// order, descending into a node's children when f returns true, as
// ast.Inspect does, but keeps its own stack. The parser walks trees before
// their statement is measured against maxTreeHeight: the %CHARSET check of
// a PL/SQL unit, where ast.Inspect overflowed the stack on a sum of two
// million terms, and the select list of a query FOR UPDATE ends.
func inspect(root nodes.Node, f func(nodes.Node) bool) {
	stack := []nodes.Node{root}
	var kids childLister
	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if n == nil || !f(n) {
			continue
		}
		children := kids.list(n)
		for i := len(children) - 1; i >= 0; i-- {
			stack = append(stack, children[i])
		}
	}
}

// childLister is a Visitor that lets ast.Walk enter only the node it starts
// at, so Walk lists that node's children without descending into them.
type childLister struct {
	entered  bool
	children []nodes.Node
}

// list returns the children of n in the order Walk visits them, in a slice
// the next call reuses.
func (c *childLister) list(n nodes.Node) []nodes.Node {
	c.entered, c.children = false, c.children[:0]
	nodes.Walk(c, n)
	return c.children
}

func (c *childLister) Visit(n nodes.Node) nodes.Visitor {
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
