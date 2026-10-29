package parser

import nodes "github.com/bytebase/omni/pg/ast"

// maxNestingDepth bounds how deeply the parser may descend into nested
// constructs: expressions, queries, types, and the other productions that
// recurse. Without a bound, a few hundred kilobytes of '(' overflowed the Go
// stack, and a stack overflow kills the process instead of returning an
// error. PostgreSQL 17 itself rejects 10,000 nested parentheses ("memory
// exhausted") and 3,000 nested subqueries ("stack depth limit exceeded").
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

// nestedTooDeeply is the error for a statement past maxNestingDepth or
// maxTreeHeight, at byte offset pos. 54001 is PostgreSQL's
// statement_too_complex, which it raises for "stack depth limit exceeded".
func nestedTooDeeply(pos int) *ParseError {
	return &ParseError{Code: "54001", Message: "statement is nested too deeply", Position: pos}
}

// maxTreeHeight bounds how deep the finished tree of one statement may be. A
// left-deep chain (1 + 1 + ..., SELECT 1 UNION ALL SELECT 1 ...) is parsed
// by a loop, so it never reaches maxNestingDepth, yet each term adds a level
// to the tree, and the code that walks a tree recursively (ast.Walk, the
// catalog, deparse, review) overflowed the stack on chains of a few hundred
// thousand terms. PostgreSQL 17 rejects 10,000-term sums and UNION ALL chains
// ("stack depth limit exceeded").
const maxTreeHeight = 100000

// treeTooHigh reports whether stmt, parsed from textLen bytes, is more than
// maxTreeHeight nodes deep. No construct adds more than 16 levels of tree
// per byte of text (TestTreeHeightPerByte checks this), so a shorter
// statement cannot be too deep and is not walked. The walk keeps its own
// stack.
func treeTooHigh(stmt nodes.Node, textLen int) bool {
	if textLen < maxTreeHeight/16 {
		return false
	}
	return treeHeightExceeds(stmt, maxTreeHeight)
}

// treeHeightExceeds reports whether the tree under root is more than limit
// nodes deep, root counting as one.
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
		kids.entered, kids.children = false, kids.children[:0]
		nodes.Walk(&kids, top.node)
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
	children []nodes.Node
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
