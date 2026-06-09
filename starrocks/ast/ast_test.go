package ast

import "testing"

// stubNode is a test-only Node implementation used to verify that NodeLoc
// returns NoLoc() for unrecognized types.
type stubNode struct{}

func (stubNode) Tag() NodeTag { return T_Invalid }

func TestNodeLocUnknownType(t *testing.T) {
	if got := NodeLoc(stubNode{}); got != NoLoc() {
		t.Errorf("NodeLoc(stubNode) = %v, want NoLoc", got)
	}
}

func TestNodeTags(t *testing.T) {
	tests := []struct {
		node Node
		want NodeTag
	}{
		{&File{}, T_File},
		{&ObjectName{Parts: []string{"t"}}, T_ObjectName},
	}
	for _, tt := range tests {
		if got := tt.node.Tag(); got != tt.want {
			t.Errorf("Tag() = %v, want %v", got, tt.want)
		}
	}
}

func TestNodeTagString(t *testing.T) {
	tests := []struct {
		tag  NodeTag
		want string
	}{
		{T_Invalid, "Invalid"},
		{T_File, "File"},
		{T_ObjectName, "ObjectName"},
		{NodeTag(9999), "Unknown"},
	}
	for _, tt := range tests {
		if got := tt.tag.String(); got != tt.want {
			t.Errorf("NodeTag(%d).String() = %q, want %q", tt.tag, got, tt.want)
		}
	}
}

func TestLocIsValid(t *testing.T) {
	tests := []struct {
		loc  Loc
		want bool
	}{
		{Loc{0, 10}, true},
		{Loc{-1, 10}, false},
		{Loc{0, -1}, false},
		{NoLoc(), false},
	}
	for _, tt := range tests {
		if got := tt.loc.IsValid(); got != tt.want {
			t.Errorf("Loc{%d,%d}.IsValid() = %v, want %v", tt.loc.Start, tt.loc.End, got, tt.want)
		}
	}
}

func TestLocContains(t *testing.T) {
	outer := Loc{0, 100}
	inner := Loc{10, 50}
	disjoint := Loc{200, 300}

	if !outer.Contains(inner) {
		t.Error("outer should contain inner")
	}
	if inner.Contains(outer) {
		t.Error("inner should not contain outer")
	}
	if outer.Contains(disjoint) {
		t.Error("outer should not contain disjoint")
	}
	if NoLoc().Contains(inner) {
		t.Error("NoLoc should not contain anything")
	}
}

func TestLocMerge(t *testing.T) {
	a := Loc{10, 20}
	b := Loc{5, 30}
	merged := a.Merge(b)
	if merged.Start != 5 || merged.End != 30 {
		t.Errorf("Merge = %v, want {5, 30}", merged)
	}

	// Merge with NoLoc returns the valid side.
	if got := a.Merge(NoLoc()); got != a {
		t.Errorf("Merge(NoLoc) = %v, want %v", got, a)
	}
	if got := NoLoc().Merge(b); got != b {
		t.Errorf("NoLoc.Merge = %v, want %v", got, b)
	}
}

func TestObjectNameString(t *testing.T) {
	tests := []struct {
		parts []string
		want  string
	}{
		{[]string{"t"}, "t"},
		{[]string{"db", "t"}, "db.t"},
		{[]string{"catalog", "db", "t"}, "catalog.db.t"},
	}
	for _, tt := range tests {
		n := &ObjectName{Parts: tt.parts}
		if got := n.String(); got != tt.want {
			t.Errorf("ObjectName(%v).String() = %q, want %q", tt.parts, got, tt.want)
		}
	}
}

func TestWalkFile(t *testing.T) {
	// Build a small tree: File with two ObjectName children.
	child1 := &ObjectName{Parts: []string{"a"}, Loc: Loc{0, 1}}
	child2 := &ObjectName{Parts: []string{"b"}, Loc: Loc{2, 3}}
	file := &File{Stmts: []Node{child1, child2}, Loc: Loc{0, 3}}

	var visited []NodeTag
	Inspect(file, func(n Node) bool {
		if n != nil {
			visited = append(visited, n.Tag())
		}
		return true
	})

	want := []NodeTag{T_File, T_ObjectName, T_ObjectName}
	if len(visited) != len(want) {
		t.Fatalf("visited %d nodes, want %d", len(visited), len(want))
	}
	for i, tag := range visited {
		if tag != want[i] {
			t.Errorf("visited[%d] = %v, want %v", i, tag, want[i])
		}
	}
}

func TestWalkNilNode(t *testing.T) {
	// Walk(nil) should not panic.
	Walk(inspector(func(Node) bool { return true }), nil)
}

// recorder is a Visitor that logs pre-order and post-order events.
type recorder struct {
	events *[]string
}

func (r recorder) Visit(node Node) Visitor {
	if node == nil {
		*r.events = append(*r.events, "post")
		return nil
	}
	*r.events = append(*r.events, node.Tag().String())
	return r
}

func TestWalkPostOrder(t *testing.T) {
	child := &ObjectName{Parts: []string{"a"}}
	file := &File{Stmts: []Node{child}}

	var events []string
	Walk(recorder{&events}, file)

	// Expect: File (pre), ObjectName (pre), ObjectName (post), File (post)
	want := []string{"File", "ObjectName", "post", "post"}
	if len(events) != len(want) {
		t.Fatalf("events = %v, want %v", events, want)
	}
	for i, e := range events {
		if e != want[i] {
			t.Errorf("events[%d] = %q, want %q", i, e, want[i])
		}
	}
}

func TestWalkAbort(t *testing.T) {
	// Returning nil from Visit should prevent descending into children.
	child := &ObjectName{Parts: []string{"a"}}
	file := &File{Stmts: []Node{child}}

	var visited []string
	Inspect(file, func(n Node) bool {
		if n != nil {
			visited = append(visited, n.Tag().String())
		}
		// Return false for File -> should not visit children.
		return false
	})

	want := []string{"File"}
	if len(visited) != len(want) {
		t.Fatalf("visited = %v, want %v", visited, want)
	}
	if visited[0] != "File" {
		t.Errorf("visited[0] = %q, want %q", visited[0], "File")
	}
}

func TestNodeLoc(t *testing.T) {
	file := &File{Loc: Loc{0, 100}}
	obj := &ObjectName{Parts: []string{"t"}, Loc: Loc{5, 10}}

	if got := NodeLoc(file); got != file.Loc {
		t.Errorf("NodeLoc(File) = %v, want %v", got, file.Loc)
	}
	if got := NodeLoc(obj); got != obj.Loc {
		t.Errorf("NodeLoc(ObjectName) = %v, want %v", got, obj.Loc)
	}
	if got := NodeLoc(nil); got != NoLoc() {
		t.Errorf("NodeLoc(nil) = %v, want NoLoc", got)
	}
}

func TestSpanNodes(t *testing.T) {
	a := &ObjectName{Parts: []string{"a"}, Loc: Loc{10, 20}}
	b := &ObjectName{Parts: []string{"b"}, Loc: Loc{5, 30}}

	span := SpanNodes(a, b)
	if span.Start != 5 || span.End != 30 {
		t.Errorf("SpanNodes = %v, want {5, 30}", span)
	}

	// With nil entries.
	span2 := SpanNodes(nil, a, nil)
	if span2 != a.Loc {
		t.Errorf("SpanNodes(nil, a, nil) = %v, want %v", span2, a.Loc)
	}

	// Empty args.
	if got := SpanNodes(); got != NoLoc() {
		t.Errorf("SpanNodes() = %v, want NoLoc", got)
	}
}
