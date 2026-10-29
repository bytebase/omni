package parser

import (
	"reflect"
	"sync"

	"github.com/bytebase/omni/trino/ast"
)

// maxNestingDepth bounds how deeply the parser may descend into nested
// constructs: expressions, queries, relations, types, row patterns, routine
// statements, and the statements that hold another statement. Without a
// bound, a few hundred kilobytes of '(' overflowed the Go stack, and a stack
// overflow kills the process instead of returning an error. Trino 482 itself
// rejects a few thousand nested parentheses and a few hundred nested
// subqueries ("statement is too large (stack overflow while parsing)"). It is
// a variable only so tests that nest deeper on purpose can raise it.
var maxNestingDepth = 10000

// nestingTooDeep is the message of the error a statement fails with when it
// nests deeper than maxNestingDepth.
const nestingTooDeep = "statement is nested too deeply"

// enterNesting is called on entry to the productions that recurse, placed so
// that every recursion cycle in the parser passes through one. It fails when
// the statement nests deeper than maxNestingDepth; otherwise the caller
// defers leaveNesting.
//
// The first failure is kept, and every later call fails with it: the parse
// continues only to unwind, and the statement fails with that error (see
// parseSingle) whatever reading a speculative parse falls back to. The
// readings the parser speculates on, and the failures it records by offset
// (parenQueryFail, parenRelationFail), assume that a parse depends only on
// where it starts, but a nesting failure depends on how deep the parse
// started too. Failing the whole statement keeps that assumption sound, and
// a fallback reading cannot descend to the limit again and again.
func (p *Parser) enterNesting() error {
	if p.nestingErr == nil && p.depth >= maxNestingDepth {
		p.nestingErr = &ParseError{Position: p.cur.Loc.Start, End: p.cur.Loc.End, Message: nestingTooDeep}
	}
	if p.nestingErr != nil {
		e := *p.nestingErr
		return &e
	}
	p.depth++
	return nil
}

func (p *Parser) leaveNesting() { p.depth-- }

// maxTreeHeight bounds how deep the finished tree of one statement may be. A
// left-deep chain (1 + 1 + ..., SELECT 1 UNION ALL SELECT 1 ..., t JOIN t ON
// true JOIN ...) is parsed by a loop, so it never reaches maxNestingDepth,
// yet each term adds a level to the tree, and the code that walks a tree
// recursively (trino/analysis, DataType.String) overflowed the stack on
// chains of a few million terms. It is a variable only so tests that nest
// deeper than maxNestingDepth on purpose can raise it along with that.
var maxTreeHeight = 100000

// treeTooHigh reports whether root, a statement, expression, or type parsed
// from textLen bytes, is more than maxTreeHeight nodes deep. No construct
// adds more than 16 levels of tree per byte of text (TestTreeHeightPerByte
// checks this), so a shorter one cannot be too deep and is not walked.
func (p *Parser) treeTooHigh(root any, textLen int) bool {
	if textLen < maxTreeHeight/16 {
		return false
	}
	return treeHeight(root, maxTreeHeight) > maxTreeHeight
}

// treeHeight returns how many nodes deep the tree under root is, root
// counting as one, or limit+1 as soon as it is deeper than limit. The walk
// keeps its own stack, so the tree's height does not matter to the Go stack.
//
// The nodes are the structs of this package and trino/ast (the expressions,
// queries, relations, types, and the rest live here, and ast.Walk does not
// reach them), and a node's children are the nodes its exported fields hold:
// through a pointer, an interface, a struct value, or slices of any of them,
// nested to any depth (GROUPING SETS hold a [][]Expr).
// A subquery placeholder's body counts below it through SubqueryExpr.Query,
// the tree the strict parse built in place, which trino/analysis walks from
// inside the statement's.
func treeHeight(root any, limit int) int {
	type level struct {
		node  reflect.Value // a node struct
		depth int
	}
	var stack []level
	// push recurses only through the containers a field's type nests, never
	// down the tree.
	var push func(v reflect.Value, depth int)
	push = func(v reflect.Value, depth int) {
		for v.Kind() == reflect.Interface || v.Kind() == reflect.Pointer {
			if v.IsNil() {
				return
			}
			v = v.Elem()
		}
		switch {
		case v.Kind() == reflect.Slice || v.Kind() == reflect.Array:
			for j := 0; j < v.Len(); j++ {
				push(v.Index(j), depth)
			}
		case v.Kind() == reflect.Struct && isTreeNode(v.Type()):
			stack = append(stack, level{v, depth})
		}
	}
	height := 0
	push(reflect.ValueOf(root), 1)
	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if top.depth > height {
			height = top.depth
			if height > limit {
				return limit + 1
			}
		}
		for _, i := range treeChildFields(top.node.Type()) {
			push(top.node.Field(i), top.depth+1)
		}
	}
	return height
}

var (
	parserPkg        = reflect.TypeOf(Parser{}).PkgPath()
	astPkg           = reflect.TypeOf(ast.File{}).PkgPath()
	locType          = reflect.TypeOf(ast.Loc{})
	childFieldsCache sync.Map // reflect.Type -> []int
)

// isTreeNode reports whether t, a struct type, is a tree node.
func isTreeNode(t reflect.Type) bool {
	return (t.PkgPath() == parserPkg || t.PkgPath() == astPkg) && t != locType
}

// treeChildFields lists the exported fields of the node struct t that can
// hold a node.
func treeChildFields(t reflect.Type) []int {
	if f, ok := childFieldsCache.Load(t); ok {
		return f.([]int)
	}
	var fields []int
	for i := 0; i < t.NumField(); i++ {
		if f := t.Field(i); f.IsExported() && canHoldTreeNode(f.Type) {
			fields = append(fields, i)
		}
	}
	childFieldsCache.Store(t, fields)
	return fields
}

// canHoldTreeNode reports whether a field of type t can hold a node.
func canHoldTreeNode(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Pointer, reflect.Slice, reflect.Array:
		return canHoldTreeNode(t.Elem())
	case reflect.Interface:
		return t.PkgPath() == parserPkg || t.PkgPath() == astPkg
	case reflect.Struct:
		return isTreeNode(t)
	default:
		return false
	}
}
