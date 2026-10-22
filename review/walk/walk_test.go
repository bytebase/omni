package walk

import (
	"context"
	"errors"
	"testing"
)

// fakeSnapshot answers from a fixed map; every other key is Absent when
// its parent is listed in complete, Unknown otherwise.
type fakeSnapshot struct {
	present  map[Key]string
	complete map[string]bool
}

func (s fakeSnapshot) Lookup(k Key) (Fact, string) {
	if kind, ok := s.present[k]; ok {
		return Present, kind
	}
	if s.complete[k.Parent] {
		return Absent, ""
	}
	return Unknown, ""
}

func rel(schema, name string) Key { return Key{Space: "relation", Parent: Path(schema), Name: name} }
func col(schema, table, name string) Key {
	return Key{Space: "column", Parent: Path(schema, table), Name: name}
}

func snap() fakeSnapshot {
	return fakeSnapshot{
		present: map[Key]string{
			rel("public", "t"):      "table",
			rel("public", "v"):      "view",
			col("public", "t", "a"): "column",
		},
		complete: map[string]bool{Path("public"): true, Path("public", "t"): true},
	}
}

type want struct {
	fact Fact
	kind string
}

func expect(t *testing.T, l *Ledger, k Key, w want) {
	t.Helper()
	f, kind := l.Lookup(k)
	if f != w.fact || kind != w.kind {
		t.Errorf("Lookup(%q, %q, %q) = %v %q, want %v %q", k.Space, k.Parent, k.Name, f, kind, w.fact, w.kind)
	}
}

func TestSnapshotFacts(t *testing.T) {
	l := NewLedger(snap())
	expect(t, l, rel("public", "t"), want{Present, "table"})
	expect(t, l, rel("public", "x"), want{Absent, ""})
	expect(t, l, rel("other", "t"), want{Unknown, ""})
	expect(t, NewLedger(nil), rel("public", "t"), want{Unknown, ""})
}

func TestSetOverridesSnapshotAndEarlierSets(t *testing.T) {
	l := NewLedger(snap())
	l.Set(rel("public", "t"), Absent, "ignored")
	expect(t, l, rel("public", "t"), want{Absent, ""})
	l.Set(rel("public", "t"), Present, "view")
	expect(t, l, rel("public", "t"), want{Present, "view"})
	l.Set(rel("public", "t"), Unknown, "")
	expect(t, l, rel("public", "t"), want{Unknown, ""})
	expect(t, l, rel("public", "v"), want{Present, "view"})
}

func TestMayAddKeepsPresentAndLosesAbsent(t *testing.T) {
	l := NewLedger(snap())
	l.Set(rel("public", "made"), Absent, "")
	l.MayAdd(Scope{Space: "relation", Parent: Path("public"), Prefix: "m"})
	expect(t, l, rel("public", "made"), want{Unknown, ""})
	expect(t, l, rel("public", "x"), want{Absent, ""}) // outside the prefix
	expect(t, l, rel("public", "t"), want{Present, "table"})
	l.MayAdd(Scope{Space: "relation", Parent: Path("public")})
	expect(t, l, rel("public", "x"), want{Unknown, ""})
	expect(t, l, rel("public", "t"), want{Present, "table"})
	expect(t, l, col("public", "t", "b"), want{Absent, ""}) // another space
	l.Set(rel("public", "x"), Absent, "")
	expect(t, l, rel("public", "x"), want{Absent, ""}) // set after the scope stands
}

func TestMayRemoveKeepsAbsentAndLosesPresent(t *testing.T) {
	l := NewLedger(snap())
	l.MayRemove(Scope{Deep: true})
	expect(t, l, rel("public", "t"), want{Unknown, ""})
	expect(t, l, col("public", "t", "a"), want{Unknown, ""})
	expect(t, l, rel("public", "x"), want{Absent, ""})
	l.Set(rel("public", "t"), Present, "table")
	expect(t, l, rel("public", "t"), want{Present, "table"})
}

func TestMayRemoveOfKinds(t *testing.T) {
	l := NewLedger(snap())
	l.Set(rel("public", "k"), Present, "")
	l.MayRemove(Scope{Space: "relation", Parent: Path("public"), Kinds: []string{"view"}})
	expect(t, l, rel("public", "t"), want{Present, "table"})
	expect(t, l, rel("public", "v"), want{Unknown, ""})
	expect(t, l, rel("public", "k"), want{Unknown, ""}) // unknown kind may be a view
	expect(t, l, rel("public", "x"), want{Absent, ""})
}

func TestForget(t *testing.T) {
	l := NewLedger(snap())
	l.Forget(Scope{Space: "column", Parent: Path("public", "t")})
	expect(t, l, col("public", "t", "a"), want{Unknown, ""})
	expect(t, l, col("public", "t", "b"), want{Unknown, ""})
	expect(t, l, rel("public", "t"), want{Present, "table"})
}

func TestClear(t *testing.T) {
	l := NewLedger(snap())
	l.Clear(Scope{Space: "column", Parent: Path("other", "n")})
	expect(t, l, col("other", "n", "z"), want{Absent, ""})
	expect(t, l, col("other", "m", "z"), want{Unknown, ""})
	l.Clear(Scope{Parent: Path("public"), Deep: true})
	expect(t, l, rel("public", "t"), want{Absent, ""})
	expect(t, l, col("public", "t", "a"), want{Absent, ""})
	l.Set(col("public", "t", "a"), Present, "column")
	expect(t, l, col("public", "t", "a"), want{Present, "column"})
	// A scope before the clear no longer matters.
	l2 := NewLedger(snap())
	l2.MayAdd(Scope{Space: "column", Parent: Path("public", "t")})
	l2.Clear(Scope{Space: "column", Parent: Path("public", "t")})
	expect(t, l2, col("public", "t", "q"), want{Absent, ""})
}

func TestMoveCarriesChildren(t *testing.T) {
	l := NewLedger(snap())
	l.Set(col("public", "t", "b"), Present, "column")
	l.Move(Path("public", "t"), Path("public", "u"))
	// The snapshot and the change's facts follow the move.
	expect(t, l, col("public", "u", "a"), want{Present, "column"})
	expect(t, l, col("public", "u", "b"), want{Present, "column"})
	expect(t, l, col("public", "u", "c"), want{Absent, ""})
	// Nothing is left under the old path.
	expect(t, l, col("public", "t", "a"), want{Absent, ""})
	// The moved object's own key is the engine's to set.
	expect(t, l, rel("public", "t"), want{Present, "table"})
	// Scopes after the move apply under the new path.
	l.MayAdd(Scope{Space: "column", Parent: Path("public", "u")})
	expect(t, l, col("public", "u", "c"), want{Unknown, ""})
	expect(t, l, col("public", "u", "a"), want{Present, "column"})
}

func TestMoveCarriesEarlierScopes(t *testing.T) {
	l := NewLedger(snap())
	l.MayRemove(Scope{Space: "column", Parent: Path("public", "t")})
	l.Move(Path("public", "t"), Path("public", "u"))
	expect(t, l, col("public", "u", "a"), want{Unknown, ""})
	expect(t, l, col("public", "u", "c"), want{Absent, ""})
}

func TestMoveOfASchemaMovesTablesBelowIt(t *testing.T) {
	l := NewLedger(snap())
	l.Move(Path("public"), Path("renamed"))
	expect(t, l, rel("renamed", "t"), want{Present, "table"})
	expect(t, l, col("renamed", "t", "a"), want{Present, "column"})
	expect(t, l, rel("public", "t"), want{Absent, ""})
	// Moves chain.
	l.Move(Path("renamed", "t"), Path("renamed", "t2"))
	expect(t, l, col("renamed", "t2", "a"), want{Present, "column"})
}

func TestMoveFromTheRootForgetsEverything(t *testing.T) {
	l := NewLedger(snap())
	l.Move("", Path("x"))
	expect(t, l, rel("public", "t"), want{Unknown, ""})
	expect(t, l, rel("public", "x"), want{Unknown, ""})
}

func TestViolates(t *testing.T) {
	l := NewLedger(snap())
	tests := []struct {
		name  string
		check Check
		bad   bool
	}{
		{"present and wanted", Check{Key: rel("public", "t"), Want: Present}, false},
		{"absent but wanted", Check{Key: rel("public", "x"), Want: Present}, true},
		{"unknown but wanted", Check{Key: rel("other", "x"), Want: Present}, false},
		{"present of an accepted kind", Check{Key: rel("public", "t"), Want: Present, Kinds: []string{"table", "view"}}, false},
		{"present of another kind", Check{Key: rel("public", "v"), Want: Present, Kinds: []string{"table"}}, true},
		{"absent and wanted absent", Check{Key: rel("public", "x"), Want: Absent}, false},
		{"present but wanted absent", Check{Key: rel("public", "t"), Want: Absent}, true},
		{"unknown and wanted absent", Check{Key: rel("other", "t"), Want: Absent}, false},
	}
	for _, tt := range tests {
		v, bad := l.Violates(tt.check)
		if bad != tt.bad {
			t.Errorf("%s: Violates = %v, want %v", tt.name, bad, tt.bad)
		}
		if bad && v.Check.Key != tt.check.Key {
			t.Errorf("%s: violation names %v", tt.name, v.Check.Key)
		}
	}
	// A Present fact of unknown kind satisfies any kinds.
	l.Set(rel("public", "k"), Present, "")
	if _, bad := l.Violates(Check{Key: rel("public", "k"), Want: Present, Kinds: []string{"table"}}); bad {
		t.Error("a Present fact of unknown kind violated a kind check")
	}
	v, _ := l.Violates(Check{Key: rel("public", "v"), Want: Present, Kinds: []string{"table"}})
	if v.Got != Present || v.Kind != "view" {
		t.Errorf("wrong-kind violation = %v %q, want present view", v.Got, v.Kind)
	}
}

func TestRun(t *testing.T) {
	ctx := context.Background()
	create := func(name string) *Step {
		k := rel("public", name)
		return &Step{
			Checks: []Check{{Key: k, Want: Absent, Note: name}},
			Apply:  func(l *Ledger) { l.Set(k, Present, "table") },
		}
	}
	steps := []*Step{create("a"), create("b"), create("a"), create("c")}
	v, err := Run(ctx, len(steps), snap(), func(i int, _ *Ledger) *Step { return steps[i] })
	if err != nil || v == nil {
		t.Fatalf("Run = %v, %v; want a violation", v, err)
	}
	if v.Statement != 2 || v.Check.Note != "a" || v.Got != Present || v.Kind != "table" {
		t.Errorf("violation = %+v, want statement 2 on a", v)
	}

	// A nil step ends the walk before its statement; Last after its checks.
	stop := []*Step{create("a"), nil, create("a")}
	if v, _ := Run(ctx, len(stop), snap(), func(i int, _ *Ledger) *Step { return stop[i] }); v != nil {
		t.Errorf("walk went past a nil step: %+v", v)
	}
	last := []*Step{create("a"), {Last: true, Apply: func(l *Ledger) { t.Error("Apply of a Last step ran") }}, create("a")}
	if v, _ := Run(ctx, len(last), snap(), func(i int, _ *Ledger) *Step { return last[i] }); v != nil {
		t.Errorf("walk went past a Last step: %+v", v)
	}
	lastBad := []*Step{create("a"), {Last: true, Checks: []Check{{Key: rel("public", "a"), Want: Absent}}}}
	if v, _ := Run(ctx, len(lastBad), snap(), func(i int, _ *Ledger) *Step { return lastBad[i] }); v == nil || v.Statement != 1 {
		t.Errorf("a Last step's checks were not judged: %+v", v)
	}

	// A violation the engine established is reported before the checks.
	failed := []*Step{create("a"), {
		Fail:   &Violation{Check: Check{Note: "resolved nowhere"}, Got: Absent},
		Checks: []Check{{Key: rel("public", "a"), Want: Absent, Note: "a"}},
	}}
	if v, _ := Run(ctx, len(failed), snap(), func(i int, _ *Ledger) *Step { return failed[i] }); v == nil || v.Statement != 1 || v.Check.Note != "resolved nowhere" {
		t.Errorf("Fail = %+v, want statement 1 resolved nowhere", v)
	}

	// A failed check stops before Apply.
	applied := false
	failing := []*Step{{
		Checks: []Check{{Key: rel("public", "t"), Want: Absent}},
		Apply:  func(*Ledger) { applied = true },
	}}
	if v, _ := Run(ctx, 1, snap(), func(i int, _ *Ledger) *Step { return failing[i] }); v == nil || applied {
		t.Errorf("violation %+v, applied %v", v, applied)
	}

	// Lowering sees the ledger the earlier statements left.
	seen := Unknown
	reader := []*Step{create("x"), nil}
	Run(ctx, 2, snap(), func(i int, l *Ledger) *Step {
		if i == 1 {
			seen, _ = l.Lookup(rel("public", "x"))
		}
		return reader[i]
	})
	if seen != Present {
		t.Errorf("lowering saw %v for a relation the first statement created", seen)
	}
}

func TestRunStopsOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := Run(ctx, 1, nil, func(int, *Ledger) *Step { return &Step{} })
	if !errors.Is(err, context.Canceled) {
		t.Errorf("Run = %v, want context.Canceled", err)
	}
}

func TestLookupStaysFastOnLongChanges(t *testing.T) {
	// Every table of a long change gets columns, a cleared column space,
	// and generated names; looking up any one must not read them all.
	l := NewLedger(snap())
	for i := range 20000 {
		name := string(rune('a'+i%26)) + string(rune('a'+(i/26)%26)) + string(rune('a'+(i/676)%26))
		l.Clear(Scope{Space: "column", Parent: Path("public", name)})
		l.Set(col("public", name, "id"), Present, "column")
		l.MayAdd(Scope{Space: "relation", Parent: Path("public"), Prefix: name})
		l.Set(rel("public", name), Present, "table")
	}
	for range 20000 {
		if f, _ := l.Lookup(col("public", "abc", "id")); f != Present {
			t.Fatalf("Lookup = %v", f)
		}
		if f, _ := l.Lookup(rel("public", "zzzz")); f != Unknown {
			t.Fatalf("Lookup = %v", f) // zzz's generated names may start with zzzz
		}
		if f, _ := l.Lookup(rel("public", "0")); f != Absent {
			t.Fatalf("Lookup = %v", f)
		}
	}
}
