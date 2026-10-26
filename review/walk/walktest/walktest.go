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
package walktest

import (
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
// synced, and the server as the change starts.
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
	// Holds reports whether the server's catalog, as it is now, shows the
	// facts a violation claims (see Claims), and what it shows when not.
	Holds(t testing.TB, v *walk.Violation, message string) (bool, string)
}

// Digester is an Engine whose server commits DDL as it runs it, with no
// savepoint to undo a rejected statement. Judge then checks that a
// rejected statement changed nothing, by a digest of the catalog taken
// before and after it. The digest after a rejected statement is the one
// before the next.
type Digester interface {
	Digest(t testing.TB) string
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
	digester, _ := e.(Digester)
	var accepted []string
	var before string
	unchanged := false
	for _, stmt := range e.Statements(change) {
		if digester != nil && !unchanged {
			before = digester.Digest(t)
		}
		unchanged = false
		o := e.Exec(t, stmt)
		if o.Code == "" {
			accepted = append(accepted, stmt)
			continue
		}
		if digester != nil {
			if digester.Digest(t) != before {
				t.Fatalf("the server rejected %q (%s: %s) but changed the database", stmt, o.Code, o.Message)
			}
			unchanged = true
		}
		rejected++
		text := strings.Join(append(slices.Clone(accepted), stmt), "\n")
		f := e.Report(t, text)
		switch {
		case f == nil:
			if complete {
				t.Errorf("missed the rejection of %q (%s: %s)", stmt, o.Code, o.Message)
			}
			continue
		case f.Statement != len(accepted):
			t.Errorf("reported %q on statement %d, which the server accepts; it rejects %q (%s: %s)", f.Message, f.Statement, stmt, o.Code, o.Message)
			continue
		}
		found++
		if slices.Contains(e.Codes(f.Message), o.Code) {
			continue
		}
		v := e.Violation(t, text)
		if ok, why := e.Holds(t, v, f.Message); !ok {
			t.Errorf("reported %q on %q, which the server rejects with %s (%s), and the claim does not hold: %s", f.Message, stmt, o.Code, o.Message, why)
		}
	}
	if f := e.Report(t, strings.Join(accepted, "\n")); f != nil {
		t.Errorf("reported %q on statement %d of a change the server runs:\n%s", f.Message, f.Statement, strings.Join(accepted, "\n"))
	}
	return found, rejected
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
