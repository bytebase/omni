package parser

import "reflect"

// maxNestingDepth bounds how deeply recursive-descent productions may descend
// so deeply nested inputs return a *ParseError instead of overflowing the Go
// stack.
var maxNestingDepth = 256

// maxTreeHeight bounds the maximum depth of the finished AST for one statement.
// Left-deep operator chains (e.g. `1 + 1 + ...`) are parsed iteratively and do
// not trigger maxNestingDepth, yet build trees that could overflow recursive
// AST consumers.
const maxTreeHeight = 100000

const nestingErrorMessage = "statement is nested too deeply"

func (p *parser) enter() error {
	if p.depth >= maxNestingDepth {
		return p.errAt(p.cur(), nestingErrorMessage)
	}
	p.depth++
	return nil
}

func (p *parser) leave() {
	p.depth--
}

// treeTooHigh reports whether root, parsed from textLen bytes, exceeds
// maxTreeHeight AST levels. No construct adds more than 16 AST levels per byte
// of input, so shorter statements cannot exceed maxTreeHeight and skip the walk.
func treeTooHigh(root Node, textLen int) bool {
	if textLen < maxTreeHeight/16 {
		return false
	}
	return treeHeightExceeds(root, maxTreeHeight)
}

var nodeInterfaceType = reflect.TypeOf((*Node)(nil)).Elem()

// treeHeightExceeds iteratively measures whether the AST rooted at root exceeds
// limit levels of nested AST nodes/structs/slices, walking each pointer target
// at most once.
func treeHeightExceeds(root any, limit int) bool {
	if root == nil {
		return false
	}
	type frame struct {
		val   reflect.Value
		depth int
	}
	visited := make(map[uintptr]bool)
	stack := []frame{{val: reflect.ValueOf(root), depth: 1}}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		v := top.val
		if !v.IsValid() {
			continue
		}
		for v.Kind() == reflect.Interface || v.Kind() == reflect.Pointer {
			if v.IsNil() {
				v = reflect.Value{}
				break
			}
			if v.Kind() == reflect.Pointer {
				ptr := v.Pointer()
				if visited[ptr] {
					v = reflect.Value{}
					break
				}
				visited[ptr] = true
			}
			v = v.Elem()
		}
		if !v.IsValid() {
			continue
		}

		if top.depth > limit {
			return true
		}

		switch v.Kind() {
		case reflect.Struct:
			n := v.NumField()
			for i := 0; i < n; i++ {
				fv := v.Field(i)
				switch fv.Kind() {
				case reflect.Pointer, reflect.Interface:
					if fv.IsNil() {
						continue
					}
					nextDepth := top.depth
					if fv.Type().Implements(nodeInterfaceType) || fv.Elem().Type().Implements(nodeInterfaceType) {
						nextDepth++
					} else if fv.Kind() == reflect.Pointer && fv.Elem().Kind() == reflect.Struct {
						nextDepth++
					}
					stack = append(stack, frame{val: fv, depth: nextDepth})
				case reflect.Slice, reflect.Array:
					if fv.Len() > 0 {
						stack = append(stack, frame{val: fv, depth: top.depth})
					}
				case reflect.Struct:
					if fv.Type() != reflect.TypeOf(span{}) {
						stack = append(stack, frame{val: fv, depth: top.depth + 1})
					}
				}
			}
		case reflect.Slice, reflect.Array:
			elemKind := v.Type().Elem().Kind()
			if elemKind != reflect.Pointer && elemKind != reflect.Interface &&
				elemKind != reflect.Struct && elemKind != reflect.Slice && elemKind != reflect.Array {
				continue
			}
			for i := 0; i < v.Len(); i++ {
				ev := v.Index(i)
				nextDepth := top.depth
				if elemKind == reflect.Pointer || elemKind == reflect.Interface || elemKind == reflect.Struct {
					nextDepth++
				}
				stack = append(stack, frame{val: ev, depth: nextDepth})
			}
		}
	}
	return false
}
