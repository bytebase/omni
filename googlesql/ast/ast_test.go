package ast

import (
	"reflect"
	"testing"
)

// Compile-time assertion: *File satisfies Node.
var _ Node = (*File)(nil)

// stubNode is a private test-only Node implementation that the production
// NodeLoc switch deliberately does not handle. It exists solely to verify
// that NodeLoc returns NoLoc() for unrecognized concrete types.
type stubNode struct{}

func (stubNode) Tag() NodeTag { return T_Invalid }

// -----------------------------------------------------------------------
// Loc tests
// -----------------------------------------------------------------------

func TestNoLoc(t *testing.T) {
	l := NoLoc()
	if l.Start != -1 || l.End != -1 {
		t.Errorf("NoLoc() = %+v, want {-1, -1}", l)
	}
	if l.IsValid() {
		t.Error("NoLoc().IsValid() should be false")
	}
}

func TestLocIsValid(t *testing.T) {
	cases := []struct {
		name string
		loc  Loc
		want bool
	}{
		{"both valid", Loc{Start: 0, End: 5}, true},
		{"zero range", Loc{Start: 3, End: 3}, true},
		{"start unknown", Loc{Start: -1, End: 5}, false},
		{"end unknown", Loc{Start: 0, End: -1}, false},
		{"both unknown", Loc{Start: -1, End: -1}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.loc.IsValid(); got != c.want {
				t.Errorf("Loc{%d,%d}.IsValid() = %v, want %v", c.loc.Start, c.loc.End, got, c.want)
			}
		})
	}
}

func TestLocContains(t *testing.T) {
	cases := []struct {
		name string
		l    Loc
		o    Loc
		want bool
	}{
		{"proper containment", Loc{0, 10}, Loc{2, 8}, true},
		{"equal endpoints", Loc{0, 10}, Loc{0, 10}, true},
		{"equal start, smaller end", Loc{0, 10}, Loc{0, 5}, true},
		{"smaller start, equal end", Loc{0, 10}, Loc{5, 10}, true},
		{"disjoint after", Loc{0, 5}, Loc{6, 10}, false},
		{"disjoint before", Loc{6, 10}, Loc{0, 5}, false},
		{"partial overlap", Loc{0, 5}, Loc{3, 8}, false},
		{"l invalid", NoLoc(), Loc{0, 5}, false},
		{"o invalid", Loc{0, 5}, NoLoc(), false},
		{"both invalid", NoLoc(), NoLoc(), false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.l.Contains(c.o); got != c.want {
				t.Errorf("Contains: %+v.Contains(%+v) = %v, want %v", c.l, c.o, got, c.want)
			}
		})
	}
}

func TestLocMerge(t *testing.T) {
	cases := []struct {
		name string
		l    Loc
		o    Loc
		want Loc
	}{
		{"both valid disjoint", Loc{0, 5}, Loc{10, 20}, Loc{0, 20}},
		{"both valid overlap", Loc{0, 10}, Loc{5, 15}, Loc{0, 15}},
		{"both valid nested", Loc{0, 100}, Loc{20, 30}, Loc{0, 100}},
		{"both valid equal", Loc{5, 5}, Loc{5, 5}, Loc{5, 5}},
		{"extends start only", Loc{10, 20}, Loc{0, 15}, Loc{0, 20}},
		{"extends end only", Loc{0, 10}, Loc{5, 25}, Loc{0, 25}},
		{"l invalid", NoLoc(), Loc{0, 5}, Loc{0, 5}},
		{"o invalid", Loc{0, 5}, NoLoc(), Loc{0, 5}},
		{"both invalid", NoLoc(), NoLoc(), NoLoc()},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.l.Merge(c.o); got != c.want {
				t.Errorf("Merge: %+v.Merge(%+v) = %+v, want %+v", c.l, c.o, got, c.want)
			}
		})
	}
}

// -----------------------------------------------------------------------
// NodeLoc / SpanNodes tests
// -----------------------------------------------------------------------

func TestNodeLocNil(t *testing.T) {
	if got := NodeLoc(nil); got != NoLoc() {
		t.Errorf("NodeLoc(nil) = %+v, want NoLoc()", got)
	}
}

func TestNodeLocFile(t *testing.T) {
	f := &File{Loc: Loc{Start: 7, End: 42}}
	if got := NodeLoc(f); got != (Loc{Start: 7, End: 42}) {
		t.Errorf("NodeLoc(File) = %+v, want {7, 42}", got)
	}
}

func TestNodeLocUnknown(t *testing.T) {
	if got := NodeLoc(stubNode{}); got != NoLoc() {
		t.Errorf("NodeLoc(stubNode) = %+v, want NoLoc()", got)
	}
}

func TestSpanNodes(t *testing.T) {
	a := &File{Loc: Loc{Start: 0, End: 5}}
	b := &File{Loc: Loc{Start: 10, End: 20}}
	c := &File{Loc: NoLoc()}

	t.Run("empty", func(t *testing.T) {
		if got := SpanNodes(); got != NoLoc() {
			t.Errorf("SpanNodes() = %+v, want NoLoc()", got)
		}
	})
	t.Run("single valid", func(t *testing.T) {
		if got := SpanNodes(a); got != (Loc{0, 5}) {
			t.Errorf("SpanNodes(a) = %+v, want {0, 5}", got)
		}
	})
	t.Run("two valid", func(t *testing.T) {
		if got := SpanNodes(a, b); got != (Loc{0, 20}) {
			t.Errorf("SpanNodes(a, b) = %+v, want {0, 20}", got)
		}
	})
	t.Run("mixed valid and invalid", func(t *testing.T) {
		if got := SpanNodes(a, c, b); got != (Loc{0, 20}) {
			t.Errorf("SpanNodes(a, c, b) = %+v, want {0, 20}", got)
		}
	})
	t.Run("nil entry skipped", func(t *testing.T) {
		if got := SpanNodes(nil, a, nil, b, nil); got != (Loc{0, 20}) {
			t.Errorf("SpanNodes(nil,a,nil,b,nil) = %+v, want {0, 20}", got)
		}
	})
	t.Run("all invalid", func(t *testing.T) {
		if got := SpanNodes(c, c); got != NoLoc() {
			t.Errorf("SpanNodes(c,c) = %+v, want NoLoc()", got)
		}
	})
}

// -----------------------------------------------------------------------
// Walker tests
// -----------------------------------------------------------------------

// recordEvent captures one Visit call: tag is the visited node's tag,
// post is true when the call is the post-order Visit(nil).
type recordEvent struct {
	tag  NodeTag
	post bool
}

// recorder is a Visitor that appends every Visit call to events. It always
// recurses (returns itself) so the full traversal is captured.
type recorder struct {
	events *[]recordEvent
}

func (r *recorder) Visit(node Node) Visitor {
	if node == nil {
		*r.events = append(*r.events, recordEvent{post: true})
		return nil
	}
	*r.events = append(*r.events, recordEvent{tag: node.Tag()})
	return r
}

func TestWalkNil(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Walk(nil) panicked: %v", r)
		}
	}()
	var events []recordEvent
	Walk(&recorder{events: &events}, nil)
	if len(events) != 0 {
		t.Errorf("Walk(nil) called Visit %d times: %+v", len(events), events)
	}
}

// nestedFile builds a 3-level *File tree:
//
//	outer { inner { leaf } }
func nestedFile() (outer, inner, leaf *File) {
	leaf = &File{Loc: Loc{Start: 5, End: 10}}
	inner = &File{Stmts: []Node{leaf}, Loc: Loc{Start: 0, End: 15}}
	outer = &File{Stmts: []Node{inner}, Loc: Loc{Start: 0, End: 20}}
	return
}

func TestWalkVisitOrder(t *testing.T) {
	outer, _, _ := nestedFile()

	var events []recordEvent
	Walk(&recorder{events: &events}, outer)

	want := []recordEvent{
		{tag: T_File}, // pre-order outer
		{tag: T_File}, // pre-order inner
		{tag: T_File}, // pre-order leaf
		{post: true},  // post-order leaf
		{post: true},  // post-order inner
		{post: true},  // post-order outer
	}
	if !reflect.DeepEqual(events, want) {
		t.Errorf("visit order mismatch:\n got: %+v\nwant: %+v", events, want)
	}
}

// abortAtInner returns nil from Visit when it sees the second pre-order File
// (i.e. the inner one), pruning the leaf.
type abortAtInner struct {
	preCount int
	events   *[]recordEvent
}

func (a *abortAtInner) Visit(node Node) Visitor {
	if node == nil {
		*a.events = append(*a.events, recordEvent{post: true})
		return nil
	}
	*a.events = append(*a.events, recordEvent{tag: node.Tag()})
	a.preCount++
	if a.preCount == 2 {
		return nil // prune children of the inner File
	}
	return a
}

func TestWalkAbort(t *testing.T) {
	outer, _, _ := nestedFile()

	var events []recordEvent
	Walk(&abortAtInner{events: &events}, outer)

	// Expected: outer pre, inner pre (returned nil → no recursion, no post),
	// outer post.
	want := []recordEvent{
		{tag: T_File},
		{tag: T_File},
		{post: true},
	}
	if !reflect.DeepEqual(events, want) {
		t.Errorf("abort sequence mismatch:\n got: %+v\nwant: %+v", events, want)
	}
}

func TestInspectShortCircuit(t *testing.T) {
	outer, inner, leaf := nestedFile()

	var visited []*File
	preCount := 0
	Inspect(outer, func(n Node) bool {
		if n == nil {
			return false
		}
		f, ok := n.(*File)
		if !ok {
			return true
		}
		visited = append(visited, f)
		preCount++
		// Return false at the second pre-order visit (inner) to prune the leaf.
		return preCount < 2
	})

	if len(visited) != 2 {
		t.Fatalf("Inspect visited %d nodes, want 2 (outer + inner)", len(visited))
	}
	if visited[0] != outer {
		t.Errorf("first visited = %p, want outer %p", visited[0], outer)
	}
	if visited[1] != inner {
		t.Errorf("second visited = %p, want inner %p", visited[1], inner)
	}
	for _, v := range visited {
		if v == leaf {
			t.Error("leaf should not have been visited (pruned)")
		}
	}
}

// TestInspectFullTraversal verifies that returning true from the Inspect
// callback for every node visits the entire tree exactly once (pre-order).
func TestInspectFullTraversal(t *testing.T) {
	outer, inner, leaf := nestedFile()

	var visited []*File
	Inspect(outer, func(n Node) bool {
		if f, ok := n.(*File); ok {
			visited = append(visited, f)
		}
		return true
	})

	want := []*File{outer, inner, leaf}
	if !reflect.DeepEqual(visited, want) {
		t.Errorf("Inspect full traversal = %v, want %v", visited, want)
	}
}

// TestInspectNoNilCallback locks in the documented contract that Inspect, unlike
// go/ast.Inspect, never invokes f with a nil argument for the post-order signal.
func TestInspectNoNilCallback(t *testing.T) {
	outer, _, _ := nestedFile()

	nilCalls := 0
	Inspect(outer, func(n Node) bool {
		if n == nil {
			nilCalls++
		}
		return true
	})

	if nilCalls != 0 {
		t.Errorf("Inspect called f with nil %d times, want 0 (post-order is not forwarded)", nilCalls)
	}
}

// TestWalkSliceOfPointerNodeFields is the regression test for the genwalker
// []*<NodeStruct> traversal gap (PR #195 review, finding 2). GrantStmt and
// RevokeStmt carry []*Privilege and []*Grantee fields — full Node types — but
// before the fix genwalker's isChildType recognized only Node, []Node, and
// *<NodeStruct>, so it emitted no walk case for slice-of-pointer-to-Node fields.
// As a result Walk/Inspect over a GrantStmt visited only the GrantStmt itself,
// reaching 0 of its Privilege and Grantee children. This is latent for DCL but
// expressions/parser-select define []*Expr / []*SelectItem fields and query-span
// walks the AST, so the generator had to descend into []*T node fields. The test
// asserts the children are now reached for both statement kinds.
func TestWalkSliceOfPointerNodeFields(t *testing.T) {
	mk := func() (privs []*Privilege, grantees []*Grantee) {
		privs = []*Privilege{
			{Name: "select", Loc: Loc{0, 6}},
			{Name: "insert", Loc: Loc{8, 14}},
		}
		grantees = []*Grantee{
			{Kind: GranteeString, Value: "x", Loc: Loc{20, 23}},
			{Kind: GranteeNamedParameter, Value: "p", Loc: Loc{25, 27}},
		}
		return
	}

	cases := []struct {
		name string
		node Node
	}{
		{"GrantStmt", func() Node { p, g := mk(); return &GrantStmt{Privileges: p, Grantees: g} }()},
		{"RevokeStmt", func() Node { p, g := mk(); return &RevokeStmt{Privileges: p, Grantees: g} }()},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			seen := map[NodeTag]int{}
			Inspect(c.node, func(n Node) bool {
				if n != nil {
					seen[n.Tag()]++
				}
				return true
			})
			if seen[T_Privilege] != 2 {
				t.Errorf("Inspect reached %d Privilege nodes, want 2 (slice-of-pointer field not walked)", seen[T_Privilege])
			}
			if seen[T_Grantee] != 2 {
				t.Errorf("Inspect reached %d Grantee nodes, want 2 (slice-of-pointer field not walked)", seen[T_Grantee])
			}
		})
	}
}

// TestWalkValueStructNodeFields is the regression test for the genwalker
// value-struct traversal gap. WindowFrame holds `Start WindowBound` and
// `End WindowBound` — value (non-pointer, non-Node) struct fields — and
// WindowBound carries an `Offset Node` (the `expr` in `5 PRECEDING` /
// `@p FOLLOWING`). Before the fix genwalker's isChildType recognized only Node,
// []Node, *<NodeStruct> and []*<NodeStruct>, so it emitted NO case for
// WindowFrame at all and the bound offset expressions were unreachable by any
// AST pass (query-span / lineage would miss frame-bound subexpressions). The
// fix teaches the generator to descend into Node fields nested inside value
// struct fields. This test asserts both bound offsets are now visited by Walk
// (with the post-order Visit(nil)) and Inspect.
func TestWalkValueStructNodeFields(t *testing.T) {
	// A BETWEEN frame `ROWS BETWEEN <startOffset> PRECEDING AND <endOffset>
	// FOLLOWING` with a distinct leaf Literal at each bound offset.
	startOffset := &Literal{Kind: LitInt, Ival: 5, Loc: Loc{0, 1}}
	endOffset := &Literal{Kind: LitInt, Ival: 10, Loc: Loc{2, 4}}
	frame := &WindowFrame{
		Kind:    FrameRows,
		Between: true,
		Start:   WindowBound{Kind: BoundPreceding, Offset: startOffset},
		End:     WindowBound{Kind: BoundFollowing, Offset: endOffset},
	}

	// Walk must reach the frame, both bound offset literals, and the matching
	// post-order Visit(nil) for each visited node.
	var events []recordEvent
	Walk(&recorder{events: &events}, frame)
	want := []recordEvent{
		{tag: T_WindowFrame}, // pre-order frame
		{tag: T_Literal},     // pre-order Start.Offset
		{post: true},         // post-order Start.Offset
		{tag: T_Literal},     // pre-order End.Offset
		{post: true},         // post-order End.Offset
		{post: true},         // post-order frame
	}
	if !reflect.DeepEqual(events, want) {
		t.Errorf("Walk over WindowFrame visit order mismatch:\n got: %+v\nwant: %+v", events, want)
	}

	// Inspect must reach the frame and both offset literals (no post-order).
	var visited []Node
	Inspect(frame, func(n Node) bool {
		visited = append(visited, n)
		return true
	})
	if len(visited) != 3 {
		t.Fatalf("Inspect visited %d nodes, want 3 (frame + 2 bound offsets); value-struct field not walked", len(visited))
	}
	if visited[0] != Node(frame) {
		t.Errorf("Inspect[0] = %v, want the WindowFrame", visited[0])
	}
	// Identity check: the exact offset Node pointers must be the ones reached,
	// proving Walk descended into n.Start.Offset / n.End.Offset specifically.
	if visited[1] != Node(startOffset) {
		t.Errorf("Inspect[1] = %v, want Start.Offset literal %p", visited[1], startOffset)
	}
	if visited[2] != Node(endOffset) {
		t.Errorf("Inspect[2] = %v, want End.Offset literal %p", visited[2], endOffset)
	}
}

// TestWalkValueStructNilOffsets verifies the value-struct walk is nil-safe: a
// frame whose bounds carry no offset (UNBOUNDED PRECEDING / CURRENT ROW leave
// WindowBound.Offset nil) must not panic and must visit only the frame itself.
func TestWalkValueStructNilOffsets(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Walk over WindowFrame with nil bound offsets panicked: %v", r)
		}
	}()
	frame := &WindowFrame{
		Kind:    FrameRows,
		Between: true,
		Start:   WindowBound{Kind: BoundUnboundedPreceding}, // Offset nil
		End:     WindowBound{Kind: BoundCurrentRow},         // Offset nil
	}
	var visited []NodeTag
	Inspect(frame, func(n Node) bool {
		visited = append(visited, n.Tag())
		return true
	})
	if len(visited) != 1 || visited[0] != T_WindowFrame {
		t.Errorf("Inspect over nil-offset frame visited %v, want [T_WindowFrame] only", visited)
	}
}

// -----------------------------------------------------------------------
// NodeTag tests
// -----------------------------------------------------------------------

func TestNodeTagString(t *testing.T) {
	cases := []struct {
		tag  NodeTag
		want string
	}{
		{T_Invalid, "Invalid"},
		{T_File, "File"},
		{NodeTag(9999), "Unknown"},
	}
	for _, c := range cases {
		if got := c.tag.String(); got != c.want {
			t.Errorf("NodeTag(%d).String() = %q, want %q", c.tag, got, c.want)
		}
	}
}

// TestFileTag verifies the File root container reports the correct tag.
func TestFileTag(t *testing.T) {
	f := &File{}
	if got := f.Tag(); got != T_File {
		t.Errorf("(*File).Tag() = %v, want T_File", got)
	}
}
