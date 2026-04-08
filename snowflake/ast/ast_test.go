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
