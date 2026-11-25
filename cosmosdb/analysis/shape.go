package analysis

import (
	"strconv"
	"strings"

	"github.com/bytebase/omni/cosmosdb/ast"
)

// shapeKind says what is known of a value's structure.
type shapeKind int

const (
	// leafShape is a value whose members are not known: a stored property
	// or a computed value. Its sources are its paths.
	leafShape shapeKind = iota
	// objectShape is an object of known members: an object literal, or the
	// item of a SELECT that names its items.
	objectShape
	// arrayShape is an array of known elements: an array literal.
	arrayShape
	// elementsShape is an array whose every element has one shape: the
	// result of ARRAY(SELECT ...).
	elementsShape
	// unionShape is a value that is one of several (c.a ?? c.b, a ? b : c):
	// a property of it is that property of each.
	unionShape
	// opaqueShape is a value computed from several by an operator or a
	// function (s + t, ARRAY_CONCAT(a, b)): its sources are theirs, but its
	// structure is not, so a property of it is that property of every source.
	opaqueShape
)

// maxPathLen bounds a path: the container's name and the 128 levels Cosmos
// DB nests a document to
// (https://learn.microsoft.com/en-us/azure/cosmos-db/concepts-limits#per-item-limits).
// A longer path names no stored property; its first maxPathLen selectors
// name the deepest one that can hold it, which covers it for masking, and a
// chain of selectors applied at each of many levels stays linear.
const maxPathLen = 129

// A shape is what is known of a value: where it comes from and, for a value
// the query built, its members. A property of an object literal is then the
// member it selects ({"k": c.name}.k is c.name) rather than a property of
// the member's source (c.name.k), however many subqueries and FROM bindings
// the value passes through. Shapes are shared, never changed once built.
type shape struct {
	kind    shapeKind
	leaf    []FieldPath // leafShape: the value's sources
	keys    []string    // objectShape: member names, in order
	members []*shape    // objectShape: member shapes; arrayShape: elements; elementsShape: the one element; unionShape, opaqueShape: operands
}

// union is the shape of a value that is one of members: the one member
// when they are all the same shape, or a union of the distinct ones.
func union(members []*shape) *shape {
	d := distinct(members)
	switch len(d) {
	case 0:
		return leaf(nil)
	case 1:
		return d[0]
	}
	return &shape{kind: unionShape, members: d}
}

// opaque is the shape of a value computed from members by an operator or a
// function, whose structure is not theirs.
func opaque(members []*shape) *shape {
	d := distinct(members)
	switch {
	case len(d) == 0:
		return leaf(nil)
	case len(d) == 1 && (d[0].kind == leafShape || d[0].kind == opaqueShape):
		return d[0]
	}
	return &shape{kind: opaqueShape, members: d}
}

// distinct is members without repeats, in order.
func distinct(members []*shape) []*shape {
	seen := make(map[*shape]bool, len(members))
	var out []*shape
	for _, m := range members {
		if !seen[m] {
			seen[m] = true
			out = append(out, m)
		}
	}
	return out
}

func leaf(paths []FieldPath) *shape { return &shape{kind: leafShape, leaf: paths} }

// paths appends every source of the value to dst. A shape reached twice
// (a binding named twice in one literal) is walked once, so a value built
// by doubling a binding at every level stays linear.
func (s *shape) paths(dst []FieldPath) []FieldPath {
	if s.kind == leafShape {
		return append(dst, s.leaf...)
	}
	seen := map[*shape]bool{}
	stack := []*shape{s}
	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if seen[n] {
			continue
		}
		seen[n] = true
		if n.kind == leafShape {
			dst = append(dst, n.leaf...)
			continue
		}
		for i := len(n.members) - 1; i >= 0; i-- {
			stack = append(stack, n.members[i])
		}
	}
	return dst
}

// member is the shape of property name of the value; nil when the value
// is known not to have it, which is undefined and returns nothing.
func (s *shape) member(name string) *shape {
	found := -1
	for i, k := range s.keys {
		if k == name {
			found = i // the last of repeated keys wins
		}
	}
	if found < 0 {
		return nil
	}
	return s.members[found]
}

// element is the shape of each element of the value, which IN iterates.
// It is memoized per shape, so iterating a shared value walks it once.
func (a *analyzer) element(s *shape) *shape {
	if e, ok := a.elements[s]; ok {
		return e
	}
	var e *shape
	switch s.kind {
	case elementsShape:
		e = s.members[0]
	case arrayShape:
		// Each element keeps its own shape: [{"k": c.a}, {"k": c.b}]
		// iterates values whose k is c.a or c.b.
		e = union(s.members)
	case unionShape:
		out := make([]*shape, len(s.members))
		for i, m := range s.members {
			out[i] = a.element(m)
		}
		e = union(out)
	default:
		// A stored array's elements are read at its own path: array
		// indexes name no property a policy can be set on.
		e = s
	}
	if a.elements == nil {
		a.elements = map[*shape]*shape{}
	}
	a.elements[s] = e
	return e
}

// A chain is the selectors of one property path, with an id for each of
// its suffixes: equal suffixes of any two chains have the same id, so a
// selection made through either is made once.
type chain struct {
	sels []Selector
	ids  []int // ids[i] names sels[i:]; ids[len(sels)] is 0, no selector
}

// suffix interns a selector followed by the suffix next.
type suffix struct {
	sel  Selector
	next int
}

// chainOf interns the suffixes of sels.
func (a *analyzer) chainOf(sels []Selector) chain {
	if a.suffixes == nil {
		a.suffixes = map[suffix]int{}
	}
	ids := make([]int, len(sels)+1)
	for i := len(sels) - 1; i >= 0; i-- {
		key := suffix{sels[i], ids[i+1]}
		id, ok := a.suffixes[key]
		if !ok {
			id = len(a.suffixes) + 1
			a.suffixes[key] = id
		}
		ids[i] = id
	}
	return chain{sels: sels, ids: ids}
}

// selection identifies a suffix of selectors applied to a shape, by its
// structure or, blind, to every source.
type selection struct {
	shape  *shape
	suffix int
	blind  bool
}

// selectFrom applies c's selectors from the i-th on to the value s. It is
// memoized per shape and suffix, so selecting through a shape that several
// unions or accesses share builds each result once.
func (a *analyzer) selectFrom(s *shape, c chain, i int) *shape {
	if i == len(c.sels) {
		return s
	}
	key := selection{s, c.ids[i], false}
	if r, ok := a.selections[key]; ok {
		return r
	}
	r := a.selectOnce(s, c, i)
	if a.selections == nil {
		a.selections = map[selection]*shape{}
	}
	a.selections[key] = r
	return r
}

// selectOnce applies the i-th of c's selectors to s and the rest, through
// selectFrom, to what it selects, so every step is memoized.
func (a *analyzer) selectOnce(s *shape, c chain, i int) *shape {
	sel := c.sels[i]
	switch {
	case s.kind == objectShape && !sel.IsArray():
		m := s.member(sel.Name)
		if m == nil {
			return leaf(nil)
		}
		return a.selectFrom(m, c, i+1)
	case s.kind == arrayShape && sel.IsArray():
		if sel.ArrayIndex >= len(s.members) {
			return leaf(nil)
		}
		return a.selectFrom(s.members[sel.ArrayIndex], c, i+1)
	case s.kind == elementsShape && sel.IsArray():
		return a.selectFrom(s.members[0], c, i+1)
	case s.kind == unionShape:
		// The selectors select from each operand.
		out := make([]*shape, len(s.members))
		for k, m := range s.members {
			out[k] = a.selectFrom(m, c, i)
		}
		return union(out)
	}
	// A stored property, a computed value, or a selector a built value does
	// not answer by structure (an index of an object, a name of an array):
	// the selectors select that property of every source.
	return a.selectAll(s, c, i)
}

// selectAll applies c's selectors from the i-th on to every source of s,
// ignoring its structure. It is memoized as selectFrom is.
func (a *analyzer) selectAll(s *shape, c chain, i int) *shape {
	if i == len(c.sels) {
		return s
	}
	key := selection{s, c.ids[i], true}
	if r, ok := a.selections[key]; ok {
		return r
	}
	var r *shape
	if s.kind == leafShape {
		out := make([]FieldPath, len(s.leaf))
		for k, p := range s.leaf {
			out[k] = extend(p, c.sels[i:])
		}
		r = a.internLeaf(out)
	} else {
		out := make([]*shape, len(s.members))
		for k, m := range s.members {
			out[k] = a.selectAll(m, c, i)
		}
		r = union(out)
	}
	if a.selections == nil {
		a.selections = map[selection]*shape{}
	}
	a.selections[key] = r
	return r
}

// internLeaf is the leaf of paths, one shape for each distinct list of
// paths however it was reached (s.k.k, or s.k and then .k), so unions of
// equal selections merge.
func (a *analyzer) internLeaf(paths []FieldPath) *shape {
	var b strings.Builder
	for _, p := range paths {
		for _, sel := range p {
			b.WriteString(strconv.Itoa(sel.ArrayIndex))
			b.WriteByte(':')
			b.WriteString(strconv.Itoa(len(sel.Name)))
			b.WriteByte(':')
			b.WriteString(sel.Name)
		}
		b.WriteByte('|')
	}
	key := b.String()
	if l, ok := a.leaves[key]; ok {
		return l
	}
	l := leaf(paths)
	if a.leaves == nil {
		a.leaves = map[string]*shape{}
	}
	a.leaves[key] = l
	return l
}

// extend returns a new path: p followed by sels, at most maxPathLen long.
func extend(p FieldPath, sels []Selector) FieldPath {
	if len(p) >= maxPathLen {
		return p
	}
	if room := maxPathLen - len(p); len(sels) > room {
		sels = sels[:room]
	}
	out := make(FieldPath, 0, len(p)+len(sels))
	out = append(out, p...)
	return append(out, sels...)
}

// shapeOf resolves an expression to the shape of its value.
func (a *analyzer) shapeOf(expr ast.ExprNode, sc *scope) *shape {
	switch e := expr.(type) {
	case *ast.ColumnRef:
		if b := sc.lookup(e.Name); b != nil {
			return b
		}
		return leaf([]FieldPath{{ItemSelector(e.Name)}})
	case *ast.DotAccessExpr, *ast.BracketAccessExpr:
		return a.chainShape(e, sc)
	case *ast.CreateObjectExpr:
		s := &shape{kind: objectShape}
		for _, f := range e.Fields {
			s.keys = append(s.keys, f.Key)
			s.members = append(s.members, a.shapeOf(f.Value, sc))
		}
		return s
	case *ast.CreateArrayExpr:
		s := &shape{kind: arrayShape}
		for _, el := range e.Elements {
			s.members = append(s.members, a.shapeOf(el, sc))
		}
		return s
	case *ast.SubLink:
		return a.selectStmt(e.Select, sc).item()
	case *ast.SubqueryFuncExpr:
		return a.selectStmt(e.Select, sc).item()
	case *ast.ArrayExpr:
		return &shape{kind: elementsShape, members: []*shape{a.selectStmt(e.Select, sc).item()}}
	case *ast.BinaryExpr:
		if e.Op == "??" {
			// The value is one operand or the other.
			return union([]*shape{a.shapeOf(e.Left, sc), a.shapeOf(e.Right, sc)})
		}
	case *ast.TernaryExpr:
		// The value is one branch or the other; the condition is read, and
		// stays a source of it.
		return union([]*shape{opaque([]*shape{a.shapeOf(e.Cond, sc)}), a.shapeOf(e.Then, sc), a.shapeOf(e.Else, sc)})
	case *ast.ExistsExpr:
		// EXISTS returns a boolean, which carries no value of the subquery;
		// the subquery's predicates are still collected.
		a.selectStmt(e.Select, sc)
		return leaf(nil)
	}
	return opaque(a.exprShapes(nil, expr, sc))
}

// chainShape resolves a property path: the value it is rooted at, then
// each selector. A property of a value the query built is the member it
// selects; a property of any other value ((c.a ?? c.b).city, a stored
// property) is that property of each of the value's sources.
func (a *analyzer) chainShape(expr ast.ExprNode, sc *scope) *shape {
	// The selectors are found from the last to the first.
	var sels []Selector
	root := expr
peel:
	for {
		switch e := root.(type) {
		case *ast.DotAccessExpr:
			sels = append(sels, ItemSelector(e.Property))
			root = e.Expr
		case *ast.BracketAccessExpr:
			sels = append(sels, bracketSelector(e.Index))
			root = e.Expr
		default:
			break peel
		}
	}
	for i, j := 0, len(sels)-1; i < j; i, j = i+1, j-1 {
		sels[i], sels[j] = sels[j], sels[i]
	}
	return a.selectFrom(a.shapeOf(root, sc), a.chainOf(sels), 0)
}
