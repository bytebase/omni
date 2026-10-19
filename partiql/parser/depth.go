package parser

import (
	"errors"
	"reflect"

	"github.com/bytebase/omni/partiql/ast"
)

// maxNestingDepth bounds how deeply the parser may descend into nested
// constructs: expressions, queries, table references, and graph patterns.
// Without a bound, about 300KB of '(' overflowed the Go stack, and a stack
// overflow kills the process instead of returning an error.
const maxNestingDepth = 10000

const nestingErrorMessage = "statement is nested too deeply"

// enterNesting is called on entry to the productions that recurse, placed so
// that every recursion cycle in the parser passes through one. It fails when
// the statement nests deeper than maxNestingDepth; otherwise the caller
// defers leaveNesting.
func (p *Parser) enterNesting() error {
	if p.depth >= maxNestingDepth {
		return nestedTooDeeply(p.cur)
	}
	p.depth++
	return nil
}

func (p *Parser) leaveNesting() { p.depth-- }

// nestedTooDeeply is the error for a statement past maxNestingDepth or
// maxTreeHeight, at token tok.
func nestedTooDeeply(tok Token) *ParseError {
	return &ParseError{Message: nestingErrorMessage, Position: tok.Loc.Start, End: tok.Loc.End}
}

// isNestingError reports whether err is the error enterNesting returns.
func isNestingError(err error) bool {
	var pe *ParseError
	return errors.As(err, &pe) && pe.Message == nestingErrorMessage
}

// maxTreeHeight bounds how deep the finished tree of one statement may be. A
// left-deep chain (1 + 1 + ..., a OR b OR ..., SELECT ... UNION ALL
// SELECT ..., t CROSS JOIN t ...) is parsed by a loop, so it never reaches
// maxNestingDepth, yet each term adds a level to the tree, and the code that
// walks a tree recursively (ast.NodeToString, analysis) overflowed the stack
// on chains of a million terms.
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
	var children []ast.Node
	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if top.depth > limit {
			return true
		}
		children = appendChildren(children[:0], top.node)
		for _, child := range children {
			stack = append(stack, level{child, top.depth + 1})
		}
	}
	return false
}

var nodeInterface = reflect.TypeOf((*ast.Node)(nil)).Elem()

// appendChildren appends the nodes that n's fields hold directly. The ast
// package has no walker, so the fields are read through reflection: every
// node type is a struct whose node-valued fields are nodes or slices of
// nodes, and reading all of them reaches every node a consumer can recurse
// into.
func appendChildren(dst []ast.Node, n ast.Node) []ast.Node {
	v := reflect.ValueOf(n)
	if v.Kind() != reflect.Pointer || v.IsNil() || v.Elem().Kind() != reflect.Struct {
		return dst
	}
	v = v.Elem()
	for i := 0; i < v.NumField(); i++ {
		if !v.Type().Field(i).IsExported() {
			continue
		}
		f := v.Field(i)
		if f.Kind() == reflect.Slice {
			for j := 0; j < f.Len(); j++ {
				dst = appendNode(dst, f.Index(j))
			}
			continue
		}
		dst = appendNode(dst, f)
	}
	return dst
}

// appendNode appends v when it holds a non-nil node.
func appendNode(dst []ast.Node, v reflect.Value) []ast.Node {
	if v.Kind() == reflect.Interface {
		if v.IsNil() {
			return dst
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Pointer || v.IsNil() || !v.Type().Implements(nodeInterface) {
		return dst
	}
	return append(dst, v.Interface().(ast.Node))
}
