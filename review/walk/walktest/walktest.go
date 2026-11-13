// Package walktest holds an engine's WALK_THROUGH to the engine itself.
// The engine's tests supply an Engine, set up for one change on a real
// server; Judge runs the change there statement by statement and fails
// the test where the walk reports what the server does not show.
//
// The method, for each statement the server rejects: the walk over the
// statements it accepted before it, then that statement, may report
// nothing or that statement. What it reports must be true: either the
// server rejects the statement with an error code the message stands for,
// or the server's catalog, as it was right before the statement, shows
// the facts the finding claims (the statement then fails for a reason the
// server checks first, and for this one too). The walk over the
// statements the server accepts must report nothing.
//
// An engine shows its catalog one of two ways. A Lister lists the server's
// objects before every statement, with queries that are the same whatever
// the walk reports, and Judge reads the claimed facts off the list: what
// the oracle sends then depends on the change and the server alone, so a
// tape (enginetest/tape) can replay it. A Holder is asked about the facts
// a finding claims, so what it sends follows the walk's findings, and it
// runs on a live server only.
package walktest

import (
	"cmp"
	"fmt"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/bytebase/omni/review/walk"
)

// Outcome is how the server took one statement of a change.
type Outcome struct {
	// Code is the server's error code, "" when it accepted the statement.
	Code    string
	Message string
}

// Finding is the walk's report on a change: the statement it names and
// what it says.
type Finding struct {
	Statement int
	Message   string
}

// Engine is one engine's side of the oracle for one change: the target
// synced, and the server as the change starts. An Engine is also a Lister
// or a Holder, which is how Judge confirms a finding the server's error
// code does not.
type Engine interface {
	// Statements splits a change into the statements the server runs.
	Statements(change string) []string
	// Exec runs one statement. A statement the server rejects must leave
	// the database as it was, undone by the engine, as with a savepoint,
	// or by the server's own atomicity.
	Exec(t testing.TB, stmt string) Outcome
	// Report reviews a change on the target and returns its WALK_THROUGH
	// finding, or nil.
	Report(t testing.TB, change string) *Finding
	// Violation returns the walk's violation for a change Report reports
	// on.
	Violation(t testing.TB, change string) *walk.Violation
	// Codes lists the error codes a finding's message stands for.
	Codes(message string) []string
}

// Object is one object the server holds: its key in the walk's spaces, with
// each name as the server stores it, and its kind where the walk claims a
// kind the server can confirm ("" elsewhere).
type Object struct {
	Key  walk.Key
	Kind string
}

// Lister is an Engine that lists the server's objects, so Judge can read a
// claim off the list and tell whether a rejected statement changed the
// database. Judge asks for the list before every statement and after each
// one the server rejects, whatever the walk reports.
type Lister interface {
	// Inventory lists every object the walk may claim a fact of, as the
	// server holds them now, the session's temporary tables among them.
	// change is the whole change, for an engine that can only look for
	// such a table by a name the change holds. The queries it sends depend
	// on nothing else.
	Inventory(t testing.TB, change string) []Object
	// Same reports whether a key the walk claims names a stored object's
	// key, as the engine folds names.
	Same(claimed, stored walk.Key) bool
}

// Holder is an Engine that checks a claim by asking the server about it.
type Holder interface {
	// Holds reports whether the server's catalog, as it is now, shows the
	// facts a violation claims (see Claims), and what it shows when not.
	Holds(t testing.TB, v *walk.Violation, message string) (bool, string)
}

// Digester is a Holder whose server commits DDL as it runs it, with no
// savepoint to undo a rejected statement. Judge then checks that a
// rejected statement changed nothing, by a digest of the catalog taken
// before and after it. The digest after a rejected statement is the one
// before the next. A Lister's list serves as its digest.
type Digester interface {
	Digest(t testing.TB) string
}

// PartialRejecter is a Digester or a Lister whose server may change the
// database in a statement it rejects: MariaDB's DROP TABLE a, missing drops
// a before it fails on missing. Where the digest shows such a change, Judge
// still counts and judges the statement, but a finding on it must name the
// error code the server gave, since the catalog no longer shows the facts
// as they were before it. Judge then judges nothing more of the change:
// the server's state no longer follows from the statements it accepted.
type PartialRejecter interface {
	RejectsPartially()
}

// Claim is one fact a violation rests on: a key and the fact the walk
// read for it, with the kind it found when that is part of the claim.
type Claim struct {
	Key  walk.Key
	Fact walk.Fact
	Kind string
}

// Claims lists the facts a violation rests on: the checked key and the
// fact found, or the keys an engine's failure read and their facts. A
// failure with no keys rests on the statement's text and claims no fact.
func Claims(v *walk.Violation) []Claim {
	if v.Check.Key != (walk.Key{}) {
		return []Claim{{Key: v.Check.Key, Fact: v.Got, Kind: v.Kind}}
	}
	claims := make([]Claim, len(v.Keys))
	for i, k := range v.Keys {
		claims[i] = Claim{Key: k, Fact: v.Facts[i]}
	}
	return claims
}

// Judge runs a change on the server through e and holds the walk to it,
// as the package comment says. With complete, every statement the server
// rejects must be found. It returns how many rejected statements the walk
// found and how many there were.
func Judge(t testing.TB, e Engine, change string, complete bool) (found, rejected int) {
	t.Helper()
	lister, _ := e.(Lister)
	holder, _ := e.(Holder)
	if lister == nil && holder == nil {
		t.Fatalf("walktest: the engine neither lists the server's objects (Lister) nor checks a claim on it (Holder)")
	}
	digester, _ := e.(Digester)
	_, partial := e.(PartialRejecter)
	// read takes the catalog as it is now: a Lister's objects, or a
	// Digester's digest. ok is false for an engine that shows neither.
	read := func() (objects []Object, digest string, ok bool) {
		switch {
		case lister != nil:
			objects = sorted(lister.Inventory(t, change))
			return objects, "", true
		case digester != nil:
			return nil, digester.Digest(t), true
		}
		return nil, "", false
	}
	var accepted []string
	var objects []Object
	var digest string
	unchanged := false
	for _, stmt := range e.Statements(change) {
		// After a rejected statement that changed nothing, the catalog read
		// after it is the one before this statement.
		watched := unchanged
		if !unchanged {
			objects, digest, watched = read()
		}
		unchanged = false
		o := e.Exec(t, stmt)
		if o.Code == "" {
			accepted = append(accepted, stmt)
			continue
		}
		changed := false
		if watched {
			if after, now, _ := read(); now != digest || !slices.Equal(after, objects) {
				if !partial {
					t.Fatalf("the server rejected %q (%s: %s) but changed the database%s", stmt, o.Code, o.Message, difference(objects, after))
				}
				changed = true
			}
			unchanged = true
		}
		rejected++
		if judgeRejected(t, e, accepted, stmt, o, complete, changed, objects) {
			found++
		}
		if changed {
			break
		}
	}
	if f := e.Report(t, strings.Join(accepted, "\n")); f != nil {
		t.Errorf("reported %q on statement %d of a change the server runs:\n%s", f.Message, f.Statement, strings.Join(accepted, "\n"))
	}
	return found, rejected
}

// judgeRejected judges the walk over the accepted statements and one the
// server rejected, and reports whether the walk found it. objects is a
// Lister's list from before the statement. When the server changed the
// database in rejecting it, only the error code can confirm a finding.
func judgeRejected(t testing.TB, e Engine, accepted []string, stmt string, o Outcome, complete, changed bool, objects []Object) bool {
	t.Helper()
	text := strings.Join(append(slices.Clone(accepted), stmt), "\n")
	f := e.Report(t, text)
	switch {
	case f == nil:
		if complete {
			t.Errorf("missed the rejection of %q (%s: %s)", stmt, o.Code, o.Message)
		}
		return false
	case f.Statement != len(accepted):
		t.Errorf("reported %q on statement %d, which the server accepts; it rejects %q (%s: %s)", f.Message, f.Statement, stmt, o.Code, o.Message)
		return false
	}
	if slices.Contains(e.Codes(f.Message), o.Code) {
		return true
	}
	if changed {
		t.Errorf("reported %q on %q, which the server rejects with %s (%s) after changing the database, so the claim cannot be checked", f.Message, stmt, o.Code, o.Message)
		return true
	}
	v := e.Violation(t, text)
	var ok bool
	var why string
	if lister, lists := e.(Lister); lists {
		ok, why = listed(Claims(v), objects, lister.Same)
	} else {
		ok, why = e.(Holder).Holds(t, v, f.Message)
	}
	if !ok {
		t.Errorf("reported %q on %q, which the server rejects with %s (%s), and the claim does not hold: %s", f.Message, stmt, o.Code, o.Message, why)
	}
	return true
}

// listed reports whether a list of the server's objects shows the facts
// claimed, and what it shows when not: an object under each key claimed
// present, of the kind claimed where the list gives one, and none under a
// key claimed absent. A claim of an unknown fact holds by itself.
func listed(claims []Claim, objects []Object, same func(claimed, stored walk.Key) bool) (bool, string) {
	for _, c := range claims {
		if c.Fact == walk.Unknown {
			continue
		}
		var kinds []string
		present := false
		for _, o := range objects {
			if same(c.Key, o.Key) {
				present = true
				kinds = append(kinds, o.Kind)
			}
		}
		switch {
		case present != (c.Fact == walk.Present):
			return false, fmt.Sprintf("%s %q in %q is present: %v", c.Key.Space, c.Key.Name, walk.Names(c.Key.Parent), present)
		case present && c.Kind != "" && !slices.Contains(kinds, c.Kind) && !slices.Contains(kinds, ""):
			return false, fmt.Sprintf("%s %q in %q is a %s, not a %s", c.Key.Space, c.Key.Name, walk.Names(c.Key.Parent), strings.Join(kinds, " and a "), c.Kind)
		}
	}
	return true, ""
}

// sorted returns the objects in one order, so two lists of the same
// objects are equal.
func sorted(objects []Object) []Object {
	out := slices.Clone(objects)
	slices.SortFunc(out, func(a, b Object) int {
		return cmp.Or(
			strings.Compare(a.Key.Space, b.Key.Space),
			strings.Compare(a.Key.Parent, b.Key.Parent),
			strings.Compare(a.Key.Name, b.Key.Name),
			strings.Compare(a.Kind, b.Kind),
		)
	})
	return out
}

// describe writes an object in one line: its space, the names that hold
// it, its name, and its kind.
func describe(o Object) string {
	s := o.Key.Space + " " + strings.Join(append(walk.Names(o.Key.Parent), o.Key.Name), ".")
	if o.Kind != "" {
		s += " (" + o.Kind + ")"
	}
	return s
}

// difference names the objects only one of two sorted lists holds, for a
// message; "" for a difference a Digester's digest showed.
func difference(before, after []Object) string {
	var b strings.Builder
	for _, o := range before {
		if !slices.Contains(after, o) {
			b.WriteString("\n- " + describe(o))
		}
	}
	for _, o := range after {
		if !slices.Contains(before, o) {
			b.WriteString("\n+ " + describe(o))
		}
	}
	if b.Len() == 0 {
		return ""
	}
	return ":" + b.String()
}

// SweepParams reads a sweep's size, seed, and server image from
// WALK_SWEEP_N, WALK_SWEEP_SEED, and WALK_SWEEP_IMAGE, over the defaults
// given.
func SweepParams(t testing.TB, n int, seed int64, image string) (int, int64, string) {
	t.Helper()
	if v := os.Getenv("WALK_SWEEP_N"); v != "" {
		var err error
		if n, err = strconv.Atoi(v); err != nil {
			t.Fatalf("WALK_SWEEP_N: %v", err)
		}
	}
	if v := os.Getenv("WALK_SWEEP_SEED"); v != "" {
		var err error
		if seed, err = strconv.ParseInt(v, 10, 64); err != nil {
			t.Fatalf("WALK_SWEEP_SEED: %v", err)
		}
	}
	if v := os.Getenv("WALK_SWEEP_IMAGE"); v != "" {
		image = v
	}
	return n, seed, image
}

// Shards is how many servers a sweep of n changes spreads over. A sweep
// runs one statement at a time against one server and waits for the
// answer, so a single-server sweep leaves the host idle: one change
// measured at 0.21 s costs 2 of this host's 20 cores. The changes do not
// depend on each other -- each starts from a reset schema -- so they can
// be dealt over servers of their own and judged at once.
//
// One shard per 100 changes, up to 4. A sweep of the default size keeps
// to the one server a developer has warm; the larger ones CI runs spread
// out.
func Shards(n int) int {
	return max(1, min(4, n/100))
}

// SweepChange is one change of a sweep: the base that sets the server up,
// and the change judged on it.
type SweepChange struct {
	Base   string
	Change string
}

// SweepChanges generates a sweep's n changes from seed, ahead of running
// any, so which changes a sweep covers does not depend on how many
// shards it deals them over.
func SweepChanges(n int, seed int64, gen func(rng *rand.Rand) (base, change string)) []SweepChange {
	rng := rand.New(rand.NewSource(seed))
	changes := make([]SweepChange, n)
	for i := range changes {
		base, change := gen(rng)
		changes[i] = SweepChange{Base: base, Change: change}
	}
	return changes
}

// Form is one statement form of an engine's statement table: SQL, a
// change of one or more statements on the engine's coverage schema, and
// Last when its final statement ends the walk after its checks by design.
type Form struct {
	SQL  string
	Last bool
}

// Coverage checks that the walk lowers every form, so a statement the
// table lists never ends the walk unseen, as when the parser gives it a
// shape the lowering does not expect. lowers walks a form and returns how
// many statements it has, how many the walk lowered before it stopped,
// and whether the last lowered one ends the walk.
func Coverage(t *testing.T, forms []Form, lowers func(t *testing.T, sql string) (n, lowered int, last bool)) {
	t.Helper()
	for _, f := range forms {
		n, lowered, last := lowers(t, f.SQL)
		switch {
		case lowered < n:
			t.Errorf("the walk stops at statement %d of %q", lowered, f.SQL)
		case last != f.Last:
			t.Errorf("%q ends the walk: %v, want %v", f.SQL, last, f.Last)
		}
	}
}
