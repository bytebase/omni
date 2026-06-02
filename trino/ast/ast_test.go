package ast

import (
	"reflect"
	"testing"
)

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
		{&Identifier{Value: "t"}, T_Identifier},
		{&QualifiedName{}, T_QualifiedName},
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
		{T_Identifier, "Identifier"},
		{T_QualifiedName, "QualifiedName"},
		{NodeTag(9999), "Unknown"},
	}
	for _, tt := range tests {
		if got := tt.tag.String(); got != tt.want {
			t.Errorf("NodeTag(%d).String() = %q, want %q", tt.tag, got, tt.want)
		}
	}
}

// TestNodeTagStringTotal guarantees String() has a non-"Unknown" arm for
// every declared tag, so a newly added tag without a String() case is caught.
func TestNodeTagStringTotal(t *testing.T) {
	for tag := T_File; tag <= T_QualifiedName; tag++ {
		if got := tag.String(); got == "Unknown" {
			t.Errorf("NodeTag(%d).String() = %q; declared tag missing a String() case", tag, got)
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
	if outer.Contains(NoLoc()) {
		t.Error("nothing should contain NoLoc")
	}
	// Inclusive on both ends: a Loc contains itself.
	if !inner.Contains(inner) {
		t.Error("a Loc should contain itself")
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
	// Both invalid -> NoLoc.
	if got := NoLoc().Merge(NoLoc()); got != NoLoc() {
		t.Errorf("NoLoc.Merge(NoLoc) = %v, want NoLoc", got)
	}
}

// TestIdentifierNormalize pins Trino's identifier-normalization semantics
// (mirrors legacy NormalizeTrinoIdentifier): unquoted names fold to lower
// case; quoted names ("..." or `...`) keep their exact case.
func TestIdentifierNormalize(t *testing.T) {
	tests := []struct {
		name string
		id   *Identifier
		want string
	}{
		{"unquoted lowercases", &Identifier{Value: "MyTable"}, "mytable"},
		{"unquoted already lower", &Identifier{Value: "orders"}, "orders"},
		{"double-quoted preserves case", &Identifier{Value: "MyTable", Quoted: true, QuoteRune: '"'}, "MyTable"},
		{"backtick-quoted preserves case", &Identifier{Value: "MyTable", Quoted: true, QuoteRune: '`'}, "MyTable"},
		{"digit-leading unquoted lowercases", &Identifier{Value: "1Col"}, "1col"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.id.Normalize(); got != tt.want {
				t.Errorf("Normalize() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestIdentifierString pins source-faithful rendering: unquoted identifiers
// are returned verbatim (case preserved); quoted identifiers re-add their
// delimiter and double any embedded delimiter. Mirrors snowflake/ast.Ident.
func TestIdentifierString(t *testing.T) {
	tests := []struct {
		name string
		id   *Identifier
		want string
	}{
		{"unquoted verbatim", &Identifier{Value: "MyTable"}, "MyTable"},
		{"double-quoted re-quotes", &Identifier{Value: "MyTable", Quoted: true, QuoteRune: '"'}, `"MyTable"`},
		{"double-quoted escapes inner quote", &Identifier{Value: `a"b`, Quoted: true, QuoteRune: '"'}, `"a""b"`},
		{"backtick re-quotes", &Identifier{Value: "My Col", Quoted: true, QuoteRune: '`'}, "`My Col`"},
		{"backtick escapes inner backtick", &Identifier{Value: "a`b", Quoted: true, QuoteRune: '`'}, "`a``b`"},
		{"quoted with no recorded delimiter defaults to double-quote", &Identifier{Value: "x", Quoted: true}, `"x"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.id.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestQualifiedNameNormalize pins the canonical (case-folded) dotted form
// used for name comparison: unquoted parts lowercase, quoted parts preserve
// case, delimiters dropped.
func TestQualifiedNameNormalize(t *testing.T) {
	tests := []struct {
		name  string
		parts []*Identifier
		want  string
	}{
		{
			name:  "single unquoted lowercases",
			parts: []*Identifier{{Value: "Orders"}},
			want:  "orders",
		},
		{
			name:  "schema.table",
			parts: []*Identifier{{Value: "Sales"}, {Value: "Orders"}},
			want:  "sales.orders",
		},
		{
			name:  "catalog.schema.table",
			parts: []*Identifier{{Value: "Hive"}, {Value: "Sales"}, {Value: "Orders"}},
			want:  "hive.sales.orders",
		},
		{
			name: "mixed quoting preserves quoted part only",
			parts: []*Identifier{
				{Value: "Hive"},
				{Value: "MixedCase", Quoted: true, QuoteRune: '"'},
				{Value: "Orders"},
			},
			want: "hive.MixedCase.orders",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qn := &QualifiedName{Parts: tt.parts}
			if got := qn.Normalize(); got != tt.want {
				t.Errorf("QualifiedName.Normalize() = %q, want %q", got, tt.want)
			}
		})
	}

	// nil receiver is safe.
	var nilQN *QualifiedName
	if got := nilQN.Normalize(); got != "" {
		t.Errorf("nil.Normalize() = %q, want empty", got)
	}
}

// TestQualifiedNameString pins source-faithful dotted rendering: case is
// preserved and originally-quoted parts are re-quoted (with embedded
// delimiters escaped). Used for deparse / error messages.
func TestQualifiedNameString(t *testing.T) {
	tests := []struct {
		name  string
		parts []*Identifier
		want  string
	}{
		{
			name:  "single unquoted preserves case",
			parts: []*Identifier{{Value: "Orders"}},
			want:  "Orders",
		},
		{
			name:  "schema.table preserves case",
			parts: []*Identifier{{Value: "Sales"}, {Value: "Orders"}},
			want:  "Sales.Orders",
		},
		{
			name: "quoted part re-quoted, unquoted preserved",
			parts: []*Identifier{
				{Value: "Sales"},
				{Value: "Order Items", Quoted: true, QuoteRune: '"'},
			},
			want: `Sales."Order Items"`,
		},
		{
			name: "embedded quote escaped",
			parts: []*Identifier{
				{Value: `we"ird`, Quoted: true, QuoteRune: '"'},
			},
			want: `"we""ird"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qn := &QualifiedName{Parts: tt.parts}
			if got := qn.String(); got != tt.want {
				t.Errorf("QualifiedName.String() = %q, want %q", got, tt.want)
			}
		})
	}

	// nil receiver is safe.
	var nilQN *QualifiedName
	if got := nilQN.String(); got != "" {
		t.Errorf("nil.String() = %q, want empty", got)
	}

	// nil component parts are skipped.
	qn := &QualifiedName{Parts: []*Identifier{{Value: "a"}, nil, {Value: "b"}}}
	if got := qn.String(); got != "a.b" {
		t.Errorf("String() with nil part = %q, want %q", got, "a.b")
	}
}

func TestQualifiedNameNormalizedParts(t *testing.T) {
	qn := &QualifiedName{Parts: []*Identifier{
		{Value: "Hive"},
		nil, // skipped
		{Value: "Orders", Quoted: true, QuoteRune: '"'},
	}}
	got := qn.NormalizedParts()
	want := []string{"hive", "Orders"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("NormalizedParts() = %v, want %v", got, want)
	}

	// nil receiver is safe.
	var nilQN *QualifiedName
	if got := nilQN.NormalizedParts(); got != nil {
		t.Errorf("nil.NormalizedParts() = %v, want nil", got)
	}
}

func TestWalkFile(t *testing.T) {
	// Build a small tree: File with two QualifiedName children, each with
	// one Identifier part.
	child1 := &QualifiedName{Parts: []*Identifier{{Value: "a", Loc: Loc{0, 1}}}, Loc: Loc{0, 1}}
	child2 := &QualifiedName{Parts: []*Identifier{{Value: "b", Loc: Loc{2, 3}}}, Loc: Loc{2, 3}}
	file := &File{Stmts: []Node{child1, child2}, Loc: Loc{0, 3}}

	var visited []NodeTag
	Inspect(file, func(n Node) bool {
		if n != nil {
			visited = append(visited, n.Tag())
		}
		return true
	})

	want := []NodeTag{
		T_File,
		T_QualifiedName, T_Identifier,
		T_QualifiedName, T_Identifier,
	}
	if !reflect.DeepEqual(visited, want) {
		t.Fatalf("visited = %v, want %v", visited, want)
	}
}

func TestWalkNilNode(t *testing.T) {
	// Walk(nil) should not panic.
	Walk(inspector(func(Node) bool { return true }), nil)
}

// TestWalkSkipsNilParts ensures the walker does not hand a typed-nil
// *Identifier to the visitor for a nil component of a QualifiedName — a
// realistic visitor that type-asserts and dereferences would otherwise
// panic. This keeps the walker consistent with NormalizedParts/String,
// which also skip nil parts.
func TestWalkSkipsNilParts(t *testing.T) {
	qn := &QualifiedName{Parts: []*Identifier{
		{Value: "a"},
		nil,
		{Value: "b"},
	}}

	var values []string
	Inspect(qn, func(n Node) bool {
		// A realistic visitor: type-assert and dereference.
		if id, ok := n.(*Identifier); ok {
			values = append(values, id.Value)
		}
		return true
	})

	want := []string{"a", "b"}
	if !reflect.DeepEqual(values, want) {
		t.Errorf("visited identifier values = %v, want %v", values, want)
	}
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
	child := &Identifier{Value: "a"}
	qn := &QualifiedName{Parts: []*Identifier{child}}
	file := &File{Stmts: []Node{qn}}

	var events []string
	Walk(recorder{&events}, file)

	// Pre-order on the way down, "post" on the way back up.
	want := []string{
		"File",
		"QualifiedName",
		"Identifier",
		"post", // Identifier post
		"post", // QualifiedName post
		"post", // File post
	}
	if !reflect.DeepEqual(events, want) {
		t.Fatalf("events = %v, want %v", events, want)
	}
}

func TestWalkAbort(t *testing.T) {
	// Returning false from the inspector should prevent descending into
	// children.
	child := &Identifier{Value: "a"}
	qn := &QualifiedName{Parts: []*Identifier{child}}
	file := &File{Stmts: []Node{qn}}

	var visited []string
	Inspect(file, func(n Node) bool {
		if n != nil {
			visited = append(visited, n.Tag().String())
		}
		// Return false for everything -> visit only the root.
		return false
	})

	want := []string{"File"}
	if !reflect.DeepEqual(visited, want) {
		t.Fatalf("visited = %v, want %v", visited, want)
	}
}

func TestNodeLoc(t *testing.T) {
	file := &File{Loc: Loc{0, 100}}
	id := &Identifier{Value: "t", Loc: Loc{5, 10}}
	qn := &QualifiedName{Loc: Loc{5, 20}}

	if got := NodeLoc(file); got != file.Loc {
		t.Errorf("NodeLoc(File) = %v, want %v", got, file.Loc)
	}
	if got := NodeLoc(id); got != id.Loc {
		t.Errorf("NodeLoc(Identifier) = %v, want %v", got, id.Loc)
	}
	if got := NodeLoc(qn); got != qn.Loc {
		t.Errorf("NodeLoc(QualifiedName) = %v, want %v", got, qn.Loc)
	}
	if got := NodeLoc(nil); got != NoLoc() {
		t.Errorf("NodeLoc(nil) = %v, want NoLoc", got)
	}
}

// TestNodeLocTypedNil guards against the Go typed-nil pitfall: a (*T)(nil)
// boxed in a non-nil Node interface must yield NoLoc(), not a panic. This
// happens when an optional child field is left unset but still passed as a
// Node.
func TestNodeLocTypedNil(t *testing.T) {
	cases := []struct {
		name string
		node Node
	}{
		{"typed-nil *File", (*File)(nil)},
		{"typed-nil *Identifier", (*Identifier)(nil)},
		{"typed-nil *QualifiedName", (*QualifiedName)(nil)},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := NodeLoc(c.node); got != NoLoc() {
				t.Errorf("NodeLoc(%s) = %v, want NoLoc", c.name, got)
			}
		})
	}

	// SpanNodes must likewise survive a typed-nil entry (it routes through
	// NodeLoc) and treat it as having no valid loc.
	valid := &Identifier{Value: "a", Loc: Loc{10, 20}}
	if got := SpanNodes((*Identifier)(nil), valid); got != valid.Loc {
		t.Errorf("SpanNodes(typed-nil, valid) = %v, want %v", got, valid.Loc)
	}
}

func TestSpanNodes(t *testing.T) {
	a := &Identifier{Value: "a", Loc: Loc{10, 20}}
	b := &Identifier{Value: "b", Loc: Loc{5, 30}}

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

	// A node whose Loc is fully invalid (NoLoc) is skipped.
	invalid := &Identifier{Value: "x", Loc: NoLoc()}
	if got := SpanNodes(invalid, a); got != a.Loc {
		t.Errorf("SpanNodes(invalid, a) = %v, want %v", got, a.Loc)
	}

	// A node whose Loc is HALF-invalid (e.g. End == -1 from an upstream
	// slip) must also be skipped, not leaked into the result.
	halfInvalid := &Identifier{Value: "y", Loc: Loc{Start: 5, End: -1}}
	if got := SpanNodes(halfInvalid); got != NoLoc() {
		t.Errorf("SpanNodes(halfInvalid) = %v, want NoLoc", got)
	}
	if got := SpanNodes(halfInvalid, a); got != a.Loc {
		t.Errorf("SpanNodes(halfInvalid, a) = %v, want %v", got, a.Loc)
	}
}
