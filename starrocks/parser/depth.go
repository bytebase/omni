package parser

import (
	"reflect"

	"github.com/bytebase/omni/starrocks/ast"
)

// maxNestingDepth bounds how deeply the parser may descend into nested
// constructs: expressions, queries, FROM sources, types, and the other
// productions that recurse. Without a bound, a few hundred kilobytes of '('
// overflowed the Go stack, and a stack overflow kills the process instead of
// returning an error. It is a variable only so tests that nest deeper on
// purpose can raise it.
var maxNestingDepth = 10000

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
// maxTreeHeight, at the token loc.
func nestedTooDeeply(loc ast.Loc) *ParseError {
	return &ParseError{Message: nestingErrorMessage, Position: loc.Start, End: loc.End}
}

// isNestingError reports whether err is the error enterNesting returns. The
// typed collection literals (array<...>[...], map<...>{...}) speculatively
// parse a type and backtrack when it fails; failing at the limit, they return
// the error instead. Backtracking reread the too-deep type from each of its
// levels in turn, up to the limit every time, so a long one never finished.
func isNestingError(err error) bool {
	pe, ok := err.(*ParseError)
	return ok && pe.Message == nestingErrorMessage
}

// maxTreeHeight bounds how deep the finished tree of one statement may be. A
// left-deep chain (1 + 1 + ..., SELECT 1 UNION ALL SELECT 1 ...) is parsed
// by a loop, so it never reaches maxNestingDepth, yet each term adds a level
// to the tree, and the code that walks a tree recursively (ast.Walk, query
// span analysis) overflowed the stack on chains of a few million terms.
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
		kids.entered, kids.children = false, kids.children[:0]
		ast.Walk(&kids, top.node)
		if sub, ok := top.node.(*ast.SubqueryExpr); ok && sub.Query != nil {
			// Walk treats a subquery placeholder as a leaf, but analysis
			// walks the body the strict parse built in place from inside
			// the statement's tree, so the body counts below it.
			kids.children = append(kids.children, sub.Query)
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
	children []ast.Node
}

func (c *childLister) Visit(n ast.Node) ast.Visitor {
	// A nil pointer stored as a Node has no children, and Walk would
	// dereference it to list them.
	if n == nil {
		return nil
	}
	if v := reflect.ValueOf(n); v.Kind() == reflect.Pointer && v.IsNil() {
		return nil
	}
	if !c.entered {
		c.entered = true
		return c
	}
	c.children = append(c.children, n)
	return nil
}
