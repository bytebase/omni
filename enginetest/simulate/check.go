// Package simulate holds an engine's catalog to its simulate.tsv: one row per
// DDL statement on a starting schema, with omni's outcome and the engine's.
// An outcome is the error code the statement is refused with, or ok: and a
// hash of the schema stored after it.
//
// The engine's side of every row is a conversation: build the fixture, run
// the statement, read what is stored. The engine's test holds it on a tape
// (enginetest/tape), so a replay reads the answers and starts no server, and
// a recording (OMNI_RECORD=1) asks the pinned server and writes both the
// tape and the file's expected column.
//
// The design is §3 of
// docs/specs/2026-11-11-engine-compatibility-without-servers-design.md.
package simulate

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/compat"
)

var show = flag.String("simulate-show", "", "print the engine's stored schema and omni's difference for the simulate.tsv rows whose fixture, rule or statement holds this text")

// maxReported caps the rows one failure prints.
const maxReported = 25

// A Fixture is a starting schema: <name>.sql in the suite's fixture
// directory, the DDL that builds it in an empty database.
type Fixture struct {
	Name string
	SQL  string
}

// A Generated row is one the file must hold: a statement on a fixture, known
// by the rule it came from.
type Generated struct {
	Fixture, Rule, Statement string
}

// FromCompat returns the rows a compat.tsv gives a simulate.tsv: each
// statement the judge accepts and pick takes, once, on each fixture. A
// statement several rules reach is known by the first. The rows of a fixture
// stay together, in the compat.tsv's order.
func FromCompat(path string, fixtures []string, pick func(statement string) bool) ([]Generated, error) {
	f, err := compat.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var statements []Generated
	seen := map[string]bool{}
	for _, r := range f.Rows {
		if r.Expected != compat.Accept || seen[r.Statement] || !pick(r.Statement) {
			continue
		}
		seen[r.Statement] = true
		statements = append(statements, Generated{Rule: r.Rule, Statement: r.Statement})
	}
	var rows []Generated
	for _, fx := range fixtures {
		for _, s := range statements {
			rows = append(rows, Generated{Fixture: fx, Rule: s.Rule, Statement: s.Statement})
		}
	}
	return rows, nil
}

// Omni is what omni's catalog does with one statement on one fixture.
type Omni struct {
	// Code is the error code omni refuses the statement with, in decimal, or
	// NoCode for a refusal that carries none. Empty when omni takes it.
	Code string
	// Same reports that omni took the statement and stores the schema the
	// engine stores. It is read only when both took the statement.
	Same bool
	// Why says how omni's schema differs from the engine's, or what its
	// refusal says, for a reader.
	Why string
}

// A Suite is one engine's simulate.tsv and what checks it.
type Suite struct {
	// Path is the simulate.tsv.
	Path string
	// Pins is the engine's pins.tsv; the default row is the pin Expected
	// was recorded at.
	Pins string
	// Fixtures is the directory of starting schemas.
	Fixtures string
	// Rows are the generated rows the file must hold, in order. Rows a
	// person added, labelled like issue:BYT-1, follow them.
	Rows []Generated
	// MadeUp matches the names the engine makes up, which differ from one
	// recording to the next. Nil for an engine that makes up none.
	MadeUp *regexp.Regexp
	// Omni runs one statement on one fixture in omni's catalog. engine is
	// the engine's answer, whose stored schema omni's is compared with.
	Omni func(fx Fixture, statement string, engine Answer) Omni
	// Engine returns the function that asks the engine about one statement
	// on one fixture. It runs in a replay too, where a tape answers: what it
	// sends must depend on the fixture and the statement alone. An error
	// means no answer, never a refusal.
	Engine func(t testing.TB) func(fx Fixture, statement string) (Answer, error)
}

// Check holds omni to the suite's simulate.tsv, as compat.Check holds a
// parser to its compat.tsv. A row passes when omni still gives the outcome
// its omni cell records: an agreement, or a known difference. It fails when
// omni's outcome moved. A row that agreed and no longer does is a
// regression. Any other move is committed with -update-compat: a gap that
// closed, and a known difference that is now another one, as when a parser
// fix turns a 1064 into the catalog's own answer. The test also fails when
// the engine's answer is not the one the row records.
func Check(t *testing.T, s Suite) {
	t.Helper()
	check(t, s)
}

func check(t testing.TB, s Suite) {
	t.Helper()
	pin, err := enginetest.LoadPin(s.Pins, "default")
	if err != nil {
		t.Fatalf("%v", err)
	}
	recording := enginetest.Recording(t)
	run := enginetest.RunCommand(t)
	record := enginetest.RecordCommand(t)
	changed := false
	f, err := ReadFile(s.Path)
	switch {
	case errors.Is(err, os.ErrNotExist) && recording:
		// A recording starts an engine's file.
		f, changed = &File{}, true
	case errors.Is(err, os.ErrNotExist):
		t.Fatalf("no %s. Record it:\n  %s", s.Path, record)
	case err != nil:
		t.Fatalf("%v", err)
	}

	merged, added, removed := merge(f.Rows, s.Rows)
	if added+removed > 0 {
		if !recording {
			t.Fatalf("%s: its generated rows are not the suite's: %d to add, %d to remove. Record them:\n  %s", s.Path, added, removed, record)
		}
		t.Logf("%s: %d generated rows added, %d removed", s.Path, added, removed)
		f.Rows, changed = merged, true
	}

	if f.Pin != pin.String() {
		if !recording {
			t.Fatalf("%s was recorded at pin %q, but %s pins %q. Record it:\n  %s", s.Path, f.Pin, s.Pins, pin, record)
		}
		// Every expected cell describes the old server: ask again.
		for i := range f.Rows {
			f.Rows[i].Expected, f.Rows[i].Source = Unjudged, noSource
		}
		f.Pin, changed = pin.String(), true
	}

	if !recording {
		var unjudged []string
		for _, r := range f.Rows {
			if r.Expected == Unjudged || r.Omni == Unjudged {
				unjudged = append(unjudged, fmt.Sprintf("%s: no engine has answered:\n    %s\n", r.where(s.Path), r.Statement))
			}
		}
		if len(unjudged) > 0 {
			// The tape holds no answer for them, and none for what follows.
			t.Fatalf("%s", summary("not judged yet. Record them:\n  "+record, unjudged))
		}
	}

	fixtures := map[string]Fixture{}
	for _, r := range f.Rows {
		if _, ok := fixtures[r.Fixture]; ok {
			continue
		}
		path := filepath.Join(s.Fixtures, r.Fixture+".sql")
		b, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("%s names fixture %s: %v", r.where(s.Path), r.Fixture, err)
		}
		fixtures[r.Fixture] = Fixture{Name: r.Fixture, SQL: strings.TrimSpace(string(b))}
	}

	ask := s.Engine(t)
	var tally tally
	var moved, closed, shifted, regressed []string
	for i := range f.Rows {
		r := &f.Rows[i]
		fx := fixtures[r.Fixture]
		answer, err := ask(fx, r.Statement)
		if err != nil {
			t.Fatalf("%s: the engine gave no answer: %v\n    %s", r.where(s.Path), err, r.Statement)
		}
		got := answer.Outcome(s.MadeUp)
		switch {
		case got == r.Expected:
		case recording && r.Expected == Unjudged:
			r.Expected, r.Source, changed = got, compat.SourceEngine, true
		case recording:
			// The pin is the one the row was recorded at, so this is an
			// answer that changes between recordings. It is not written.
			moved = append(moved, fmt.Sprintf("%s: the engine now answers %s, and the row records %s at the same pin:\n    %s\n", r.where(s.Path), got, r.Expected, r.Statement))
		default:
			moved = append(moved, fmt.Sprintf("%s: the tape answers %s, and the row records %s:\n    %s\n", r.where(s.Path), got, r.Expected, r.Statement))
		}

		omni := s.Omni(fx, r.Statement, answer)
		now := omni.outcome(answer, got)
		if recording && r.Omni == Unjudged {
			r.Omni, changed = now, true
		}
		was := fmt.Sprintf("    the row records omni %s\n", r.Omni)
		switch {
		case now == r.Omni:
		case r.Omni == r.Expected:
			regressed = append(regressed, r.describe(s.Path, now, omni)+was)
		case compat.Updating():
			r.Omni, changed = now, true
		case now == r.Expected:
			closed = append(closed, r.describe(s.Path, now, omni))
		default:
			shifted = append(shifted, r.describe(s.Path, now, omni)+was)
		}
		tally.add(*r, now)
		if *show != "" && strings.Contains(r.Fixture+"\t"+r.Rule+"\t"+r.Statement, *show) {
			t.Logf("%s%s", r.describe(s.Path, now, omni), indent(fx.SQL+"\n"+answer.Text()))
		}
	}
	if recording {
		report(t, "an answer does not reproduce, so it cannot be replayed", moved)
	} else {
		report(t, "the tape and the file disagree. Record them:\n  "+record, moved)
	}
	report(t, "a gap closed: omni now does what the engine does. Commit it:\n  "+run+" -update-compat", closed)
	report(t, "a known difference is now another one. Commit it:\n  "+run+" -update-compat", shifted)
	report(t, "a regression: omni agreed with the engine and no longer does", regressed)
	tally.log(t, s.Path, f.Pin)

	if changed {
		if err := f.Write(s.Path); err != nil {
			t.Fatalf("%v", err)
		}
		t.Logf("wrote %s", s.Path)
	}
}

// outcome is the cell omni's result is recorded as. want is the engine's.
func (o Omni) outcome(engine Answer, want string) string {
	switch {
	case o.Code != "":
		return o.Code
	case o.Same && !engine.Refused():
		return want
	default:
		return Took
	}
}

// merge returns the file's rows with its generated rows replaced by
// generated, in generated's order, followed by the hand-added rows in their
// order. A row the file already holds keeps its outcomes; a new one is
// unjudged. added and removed count the generated rows that differ.
func merge(rows []Row, generated []Generated) (merged []Row, added, removed int) {
	old := map[rowKey]Row{}
	var hand []Row
	for _, r := range rows {
		if r.hand() {
			hand = append(hand, r)
		} else {
			old[r.key()] = r
		}
	}
	for _, g := range generated {
		r := Row{Omni: Unjudged, Expected: Unjudged, Source: noSource, Fixture: g.Fixture, Rule: g.Rule, Statement: g.Statement}
		if o, ok := old[r.key()]; ok {
			r = o
			delete(old, r.key())
		} else {
			added++
		}
		merged = append(merged, r)
	}
	return append(merged, hand...), added, len(old)
}

func (r Row) where(path string) string {
	if r.line > 0 {
		return fmt.Sprintf("%s:%d: %s on %s", path, r.line, r.Rule, r.Fixture)
	}
	return fmt.Sprintf("%s: %s on %s", path, r.Rule, r.Fixture)
}

// describe renders one row and omni's outcome on it now.
func (r Row) describe(path, now string, omni Omni) string {
	s := fmt.Sprintf("%s: omni now says %s, the engine says %s:\n    %s\n", r.where(path), now, r.Expected, r.Statement)
	if omni.Why != "" {
		s += fmt.Sprintf("    omni: %s\n", omni.Why)
	}
	return s
}

func indent(text string) string {
	return "      " + strings.ReplaceAll(strings.TrimRight(text, "\n"), "\n", "\n      ") + "\n"
}

// summary renders up to maxReported of rows, under what.
func summary(what string, rows []string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d rows: %s\n", len(rows), what)
	for i, r := range rows {
		if i == maxReported {
			fmt.Fprintf(&b, "... and %d more\n", len(rows)-maxReported)
			break
		}
		b.WriteString(r)
	}
	return b.String()
}

func report(t testing.TB, what string, rows []string) {
	t.Helper()
	if len(rows) > 0 {
		t.Errorf("%s", summary(what, rows))
	}
}

// tally counts each fixture's rows by how omni and the engine answer them.
type tally struct {
	fixtures  []string
	byFixture map[string]*counts
}

type counts struct {
	rows, agree, engineOnly, omniOnly, schema, code int
}

func (c *tally) add(r Row, now string) {
	if c.byFixture == nil {
		c.byFixture = map[string]*counts{}
	}
	n, ok := c.byFixture[r.Fixture]
	if !ok {
		n = &counts{}
		c.byFixture[r.Fixture] = n
		c.fixtures = append(c.fixtures, r.Fixture)
	}
	n.rows++
	switch {
	case now == r.Expected:
		n.agree++
	case takes(r.Expected) && takes(now):
		n.schema++
	case takes(r.Expected):
		n.engineOnly++
	case takes(now):
		n.omniOnly++
	default:
		n.code++
	}
}

// log prints the counts in one block: how far the catalog is from doing
// what the engine does, with nothing else to read.
func (c *tally) log(t testing.TB, path, pin string) {
	t.Helper()
	slices.Sort(c.fixtures)
	var b strings.Builder
	fmt.Fprintf(&b, "%s at pin %s\n", path, pin)
	fmt.Fprintf(&b, "  %-22s %6s %6s %18s %16s %14s %12s\n", "fixture", "rows", "agree", "only engine takes", "only omni takes", "other schema", "other code")
	var sum counts
	line := func(name string, n counts) {
		fmt.Fprintf(&b, "  %-22s %6d %6d %18d %16d %14d %12d\n", name, n.rows, n.agree, n.engineOnly, n.omniOnly, n.schema, n.code)
	}
	for _, k := range c.fixtures {
		n := *c.byFixture[k]
		line(k, n)
		sum.rows += n.rows
		sum.agree += n.agree
		sum.engineOnly += n.engineOnly
		sum.omniOnly += n.omniOnly
		sum.schema += n.schema
		sum.code += n.code
	}
	line("all", sum)
	t.Logf("%s", strings.TrimRight(b.String(), "\n"))
}
