package walk

import (
	"context"
	"testing"
)

// runOne walks one statement, built by build, on snap, then reports what
// the ledger knows of each key afterwards.
func runOne(t *testing.T, build func(b *Builder)) (*Violation, *Ledger) {
	t.Helper()
	l := NewLedger(snap())
	var b Builder
	build(&b)
	s := b.Step()
	if s.Fail != nil {
		v := *s.Fail
		return &v, l
	}
	for _, c := range s.Checks {
		if v, ok := l.Violates(c); ok {
			return &v, l
		}
	}
	if s.Apply != nil {
		s.Apply(l)
	}
	return nil, l
}

func TestBuilderPassesDropsChecksOnKeysAnotherPartChanges(t *testing.T) {
	// Part 0 needs t absent, which the snapshot contradicts; part 1 drops
	// t. Run in passes, the parts' order is the server's, so the check is
	// left out.
	v, l := runOne(t, func(b *Builder) {
		b.Need(rel("public", "t"), Absent, nil, "t exists")
		b.Next()
		b.Set(rel("public", "t"), Absent, "")
	})
	if v != nil {
		t.Fatalf("violation %+v, want none", v)
	}
	expect(t, l, rel("public", "t"), want{Absent, ""})

	// A check on a key no other part changes stays.
	v, _ = runOne(t, func(b *Builder) {
		b.Need(rel("public", "t"), Absent, nil, "t exists")
		b.Next()
		b.Set(rel("public", "u"), Present, "table")
	})
	if v == nil || v.Check.Note != "t exists" {
		t.Fatalf("violation %+v, want the check on t", v)
	}
}

func TestBuilderLookupsKeepsEveryCheck(t *testing.T) {
	// DROP TABLE t, x: every object is looked up first, so x's absence
	// fails the statement whatever t's drop does.
	v, _ := runOne(t, func(b *Builder) {
		b.Order = Lookups
		b.Need(rel("public", "t"), Present, nil, "t missing")
		b.Set(rel("public", "t"), Absent, "")
		b.MayRemove(Scope{Space: "relation", Parent: Path("public")})
		b.Next()
		b.Need(rel("public", "x"), Present, nil, "x missing")
		b.Set(rel("public", "x"), Absent, "")
	})
	if v == nil || v.Check.Note != "x missing" {
		t.Fatalf("violation %+v, want x missing", v)
	}
}

func TestBuilderFailsSeeTheTouchedRule(t *testing.T) {
	fail := func(order Order) *Violation {
		v, _ := runOne(t, func(b *Builder) {
			b.Order = order
			b.Fail([]Key{rel("public", "t")}, Present, "fails")
			b.Next()
			b.Set(rel("public", "t"), Absent, "")
		})
		return v
	}
	if v := fail(Passes); v != nil {
		t.Errorf("Passes: violation %+v, want the failure left out", v)
	}
	if v := fail(Lookups); v == nil || v.Check.Note != "fails" {
		t.Errorf("Lookups: violation %+v, want the failure", v)
	}
	// A failure without keys rests on the text: no effect reaches it.
	v, _ := runOne(t, func(b *Builder) {
		b.Fail(nil, Present, "text")
		b.Next()
		b.Forget(Scope{Deep: true})
	})
	if v == nil || v.Check.Note != "text" {
		t.Errorf("violation %+v, want the text failure", v)
	}
}

func TestBuilderOptionalPartsWeaken(t *testing.T) {
	v, l := runOne(t, func(b *Builder) {
		// Part 0 certainly happens.
		b.Set(rel("public", "made"), Present, "table")
		b.Next()
		// Part 1 may be skipped: its check and failure are left out, and
		// its exact effects hold either way.
		b.Optional()
		b.Need(rel("public", "t"), Absent, nil, "t exists")
		b.Fail(nil, Present, "fails")
		b.Set(rel("public", "x"), Present, "index")
		b.Clear(Scope{Parent: Path("public", "t"), Deep: true})
		b.Move(Path("public", "v"), Path("public", "w"))
	})
	if v != nil {
		t.Fatalf("violation %+v, want none", v)
	}
	expect(t, l, rel("public", "made"), want{Present, "table"})
	expect(t, l, rel("public", "x"), want{Unknown, ""})
	// The Clear became a MayRemove: Present columns are Unknown, Absent
	// ones stay.
	expect(t, l, col("public", "t", "a"), want{Unknown, ""})
	expect(t, l, col("public", "t", "zz"), want{Absent, ""})
	// The Move forgot both paths.
	expect(t, l, Key{Space: "column", Parent: Path("public", "w"), Name: "a"}, want{Unknown, ""})
}

func TestBuilderSessionAndLast(t *testing.T) {
	ran := false
	var b Builder
	b.Session(func() { ran = true })
	b.Last = true
	s := b.Step()
	if !s.Last {
		t.Error("Last not carried")
	}
	s.Apply(NewLedger(nil))
	if !ran {
		t.Error("session effect did not run")
	}
}

func TestRunReportsAFailuresFacts(t *testing.T) {
	// A name resolves nowhere on a path of two schemas: the failure rests
	// on both schemas' relations being absent and on the first schema
	// existing.
	keys := []Key{{Space: "schema", Name: "public"}, rel("public", "x"), rel("other", "x")}
	s := snap()
	s.present[Key{Space: "schema", Name: "public"}] = ""
	v, err := Run(context.Background(), 1, s, func(int, *Ledger) *Step {
		var b Builder
		b.Fail(keys, Absent, "missing")
		return b.Step()
	})
	if err != nil {
		t.Fatal(err)
	}
	if v == nil || len(v.Keys) != 3 || len(v.Facts) != 3 {
		t.Fatalf("violation %+v, want three keys and facts", v)
	}
	if v.Facts[0] != Present || v.Facts[1] != Absent || v.Facts[2] != Unknown {
		t.Errorf("facts %v, want present absent unknown", v.Facts)
	}
}
