package parser

import "github.com/bytebase/omni/cassandra/ast"

// maxNestingDepth bounds how deeply the parser may descend into nested
// constructs: expressions, collection literals, function arguments, INITCOND
// literals, and data types. Without a bound, about 880KB of '[' overflowed
// the Go stack, and a stack overflow kills the process instead of returning
// an error.
const maxNestingDepth = 10000

const nestingErrorMessage = "statement is nested too deeply"

// enterNesting is called on entry to the productions that recurse, placed so
// that every recursion cycle in the parser passes through one. It fails when
// the statement nests deeper than maxNestingDepth; otherwise the caller
// defers leaveNesting.
func (p *Parser) enterNesting() error {
	if p.depth >= maxNestingDepth {
		return p.errorf(nestingErrorMessage)
	}
	p.depth++
	return nil
}

func (p *Parser) leaveNesting() { p.depth-- }

// treeTooDeepError is the nesting error for a statement whose tree is past
// maxTreeHeight, at the statement's first token.
func (p *Parser) treeTooDeepError(first Token) *ParseError {
	line, col := offsetToLineCol(p.lexer.lineIdx, first.Loc)
	return &ParseError{
		Message:  nestingErrorMessage,
		Position: first.Loc, End: first.End,
		Line:   line,
		Column: col,
		Near:   first.Str,
	}
}

// maxTreeHeight bounds how deep the finished tree of one statement may be. A
// left-deep chain of index accesses (a[0][0]...) is parsed by a loop, so it
// never reaches maxNestingDepth, yet each term adds a level to the tree, and
// ast.Walk, which recurses, overflowed the stack on three million of them.
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
