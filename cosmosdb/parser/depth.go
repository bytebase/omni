package parser

import (
	"reflect"

	nodes "github.com/bytebase/omni/cosmosdb/ast"
)

// maxNestingDepth bounds how deeply the parser may descend into nested
// constructs: expressions and subqueries. Without a bound, about 880KB of '('
// overflowed the Go stack, and a stack overflow kills the process instead of
// returning an error.
const maxNestingDepth = 10000

const nestingErrorMessage = "statement is nested too deeply"

// enterNesting is called on entry to the productions that recurse, placed so
// that every recursion cycle in the parser passes through one. It fails when
// the statement nests deeper than maxNestingDepth; otherwise the caller
// defers leaveNesting.
func (p *Parser) enterNesting() error {
	if p.depth >= maxNestingDepth {
		return &ParseError{Message: nestingErrorMessage, Position: p.cur.Loc}
	}
	p.depth++
	return nil
}

func (p *Parser) leaveNesting() { p.depth-- }

// maxTreeHeight bounds how deep the finished tree of a statement may be. A
// left-deep chain (1 + 1 + ..., a OR b OR ..., c.a.a..., c[0][0]...) is
// parsed by a loop, so it never reaches maxNestingDepth, yet each term adds a
// level to the tree, and the code that walks a tree recursively (analysis,
// ast.NodeToString) overflowed the stack on chains of a million terms.
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
	var children []nodes.Node
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

var nodeInterface = reflect.TypeOf((*nodes.Node)(nil)).Elem()

// appendChildren appends the nodes that n's fields hold directly. The ast
// package has no walker, so the fields are read through reflection: every
// node type is a struct whose node-valued fields are nodes or slices of
// nodes, and reading all of them reaches every node a consumer can recurse
// into.
func appendChildren(dst []nodes.Node, n nodes.Node) []nodes.Node {
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
func appendNode(dst []nodes.Node, v reflect.Value) []nodes.Node {
	if v.Kind() == reflect.Interface {
		if v.IsNil() {
			return dst
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Pointer || v.IsNil() || !v.Type().Implements(nodeInterface) {
		return dst
	}
	return append(dst, v.Interface().(nodes.Node))
}
