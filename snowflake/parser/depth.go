package parser

import (
	"reflect"

	"github.com/bytebase/omni/snowflake/ast"
)

// maxNestingDepth bounds how deeply the parser may descend into nested
// constructs: expressions, queries, FROM sources, types, option groups, and
// the statements that embed other statements. Without a bound, a few hundred
// kilobytes of '(' overflowed the Go stack, and a stack overflow kills the
// process instead of returning an error. Snowflake Scripting blocks keep
// their own, much lower ceiling (maxScriptDepth in scripting.go).
const maxNestingDepth = 10000

// enterNesting is called on entry to the productions that recurse, placed so
// that every recursion cycle in the parser passes through one (the scripting
// statement cycles pass through enterScript instead). It fails when the
// statement nests deeper than maxNestingDepth; otherwise the caller defers
// leaveNesting.
func (p *Parser) enterNesting() error {
	if p.depth >= maxNestingDepth {
		return nestedTooDeeply(p.cur.Loc)
	}
	p.depth++
	return nil
}

func (p *Parser) leaveNesting() { p.depth-- }

// nestedTooDeeply is the error for a statement past maxNestingDepth or
// maxTreeHeight, at the token loc.
func nestedTooDeeply(loc ast.Loc) *ParseError {
	return &ParseError{Position: loc.Start, End: loc.End, Message: "statement is nested too deeply"}
}

// maxTreeHeight bounds how deep the finished tree of one statement may be. A
// left-deep chain (1 + 1 + ..., SELECT 1 UNION ALL SELECT 1 ..., t JOIN t
// ...) is parsed by a loop, so it never reaches maxNestingDepth, yet each
// term adds a level to the tree, and the code that walks a tree recursively
// (ast.Walk and ast.Inspect, analysis, review, the advisor) overflowed the
// stack on chains of about 770,000 terms.
const maxTreeHeight = 100000

// treeTooHigh reports whether stmt, parsed from textLen bytes, is more than
// maxTreeHeight nodes deep. No construct adds more than 16 levels of tree
// per byte of text (TestTreeHeightPerByte checks this), so a shorter
// statement cannot be too deep and is not walked. The walk keeps its own
// stack.
func treeTooHigh(stmt ast.Node, textLen int) bool {
	if textLen < maxTreeHeight/16 || isNilNode(stmt) {
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
		kids.entered, kids.children = false, kids.children[:0]
		ast.Walk(&kids, top.node)
		for _, child := range kids.children {
			stack = append(stack, level{child, top.depth + 1})
		}
	}
	return false
}

// childLister is a Visitor that lets Walk enter only the node it starts at,
// so Walk lists that node's children without descending into them. Walk
// reaches the nodes inside helper structs (CTE, OrderItem, WindowSpec, ...)
// too, so those count as children of the node holding them.
type childLister struct {
	entered  bool
	children []ast.Node
}

func (c *childLister) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}
	if !c.entered {
		c.entered = true
		return c
	}
	if isNilNode(n) {
		return nil
	}
	c.children = append(c.children, n)
	return nil
}

// isNilNode reports whether n is nil or a nil pointer held in the interface,
// which has no children: walking it would dereference nil.
func isNilNode(n ast.Node) bool {
	if n == nil {
		return true
	}
	v := reflect.ValueOf(n)
	return v.Kind() == reflect.Pointer && v.IsNil()
}
