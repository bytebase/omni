package parser

import (
	"errors"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// maxNestingDepth bounds how deeply the parser may descend into nested
// constructs: expressions, queries, table references, stored-program blocks,
// and the other productions that recurse. Without a bound, a few hundred
// kilobytes of '(' overflowed the Go stack, and a stack overflow kills the
// process instead of returning an error. MariaDB 11.8 has no single nesting
// limit: it rejects 64 nested subqueries (1473, too high level of nesting
// for select) and 1,000 nested function calls (1436, thread stack overrun),
// while its parser stack takes 20,000 redundant parentheses and gives out
// ("memory exhausted", 1064) before 50,000.
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

const nestingErrorMessage = "statement is nested too deeply"

// nestedTooDeeply is the error for a statement past maxNestingDepth or
// maxTreeHeight, at byte offset pos.
func nestedTooDeeply(pos int) *ParseError {
	return &ParseError{Message: nestingErrorMessage, Position: pos}
}

// isNestingError reports whether err is nestedTooDeeply's error.
func isNestingError(err error) bool {
	var pe *ParseError
	return errors.As(err, &pe) && pe.Message == nestingErrorMessage
}

// maxTreeHeight bounds how deep the finished tree of one statement may be. A
// left-deep chain (1 + 1 + ..., SELECT 1 UNION ALL SELECT 1 ..., a OR b OR
// ..., t JOIN t ...) is parsed by a loop, so it never reaches
// maxNestingDepth, yet each term adds a level to the tree, and the code that
// walks a tree recursively (ast.Walk, NodeToString, deparse, the catalog,
// review) overflows the stack on a long enough chain.
// MariaDB 11.8 overruns its thread stack (1436) on a 10,000-term sum, but it
// keeps AND and OR chains flat and accepts 200,000 terms of either, while
// this parser nests them one level per term, so an AND or OR chain longer
// than 100,000 terms is rejected here.
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
		// Walk skips the rows of INSERT ... VALUES and of a VALUES
		// statement, which NodeToString and deparse recurse into.
		var rows [][]nodes.ExprNode
		switch n := top.node.(type) {
		case *nodes.InsertStmt:
			rows = n.Values
		case *nodes.ValuesStmt:
			rows = n.Rows
		}
		for _, row := range rows {
			for _, e := range row {
				if e != nil {
					kids.children = append(kids.children, e)
				}
			}
		}
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
