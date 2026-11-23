// Package span holds an engine's query span analysis to its span.tsv: for
// each query of a corpus, in a session on a starting schema, what omni's
// analysis answers and what the engine states about the same query. A query
// has a row for each aspect a span answers (Aspects): its output column
// names, each output column's lineage, and the tables it reads.
//
// The engine's side of every row is a conversation: build the fixture, run
// the query, read what the engine says it read. The engine's test holds it on
// a tape (enginetest/tape), so a replay reads the answers and starts no
// server, a recording (OMNI_RECORD=1) asks the pinned server and writes both
// the tape and the file's expected column, and a verification
// (OMNI_RECORD=verify) asks the pinned server again, fails on any difference
// from the committed files, and writes nothing. It is simulate's design
// (docs/specs/2026-11-11-engine-compatibility-without-servers-design.md §3)
// applied to the analysis layer.
package span

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/compat"
)

var show = flag.String("span-show", "", "print omni's and the engine's whole answers for the span.tsv rows whose fixture, session or statement holds this text")

// maxReported caps the rows one failure prints.
const maxReported = 25

// A Fixture is a starting schema: <name>.sql in the suite's fixture
// directory, the DDL that builds it in an empty catalog.
type Fixture struct {
	Name string
	SQL  string
}

// A Query is one statement of the corpus, run in a session on a fixture.
type Query struct {
	Fixture   string
	Session   string
	Statement string
}

// An Answer is one side's span of a query.
type Answer struct {
	// Error names why there is no span: the engine's error name for a query
	// it refuses, or what omni's analysis failed with. Empty for a span.
	Error string
	// Results is the output column names, in order.
	Results []string
	// Lineage is, for each output column, its source columns as
	// catalog.schema.table.column, in any order.
	Lineage [][]string
	// Tables is the tables and views read, as catalog.schema.table, in any
	// order.
	Tables []string
}

// Cell renders one aspect of the answer as a span.tsv cell.
func (a Answer) Cell(aspect string) string {
	if a.Error != "" {
		return "error:" + oneLine(a.Error)
	}
	switch aspect {
	case Results:
		names := make([]string, len(a.Results))
		for i, n := range a.Results {
			if n == "" {
				n = `""`
			}
			names[i] = oneLine(n)
		}
		return orNone(strings.Join(names, ", "))
	case Lineage:
		cols := make([]string, len(a.Lineage))
		for i, sources := range a.Lineage {
			cols[i] = orNone(strings.Join(sorted(sources), " "))
		}
		return orNone(strings.Join(cols, " | "))
	case Tables:
		return orNone(strings.Join(sorted(a.Tables), ", "))
	}
	panic("span: unknown aspect " + aspect)
}

// sorted returns the distinct values of list in order.
func sorted(list []string) []string {
	out := slices.Clone(list)
	slices.Sort(out)
	return slices.Compact(out)
}

func orNone(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

// oneLine keeps a cell on its line: a tab or a newline in a name would split
// the row.
func oneLine(s string) string {
	return strings.NewReplacer("\t", `\t`, "\n", `\n`, "\r", `\r`).Replace(s)
}

// A Suite is one engine's span.tsv and what checks it.
type Suite struct {
	// Path is the span.tsv.
	Path string
	// Pins is the engine's pins.tsv; the PinRole row (default "default") is
	// the pin Expected was recorded at.
	Pins    string
	PinRole string
	// Fixtures is the directory of starting schemas.
	Fixtures string
	// Queries are the queries the file must hold rows for, in order.
	Queries []Query
	// Omni is omni's span of a query. It runs after Engine has answered the
	// same query, so an engine's test can hand omni the schema the engine
	// stored for the fixture.
	Omni func(fx Fixture, session, statement string) Answer
	// Engine returns the function that asks the engine about one query. It
	// runs in a replay too, where a tape answers: what it sends must depend
	// on the fixture, the session and the statement alone. An error means no
	// answer, never a refusal.
	Engine func(t testing.TB) func(fx Fixture, session, statement string) (Answer, error)
}

// Check holds omni to the suite's span.tsv, as simulate.Check holds a
// catalog to its simulate.tsv. A row passes when omni still gives the answer
// its omni cell records: an agreement, or a known difference. It fails when
// omni's answer moved: a row that agreed and no longer does is a regression,
// and any other move is committed with -update-compat. The test also fails
// when the engine's answer is not the one the row records.
func Check(t *testing.T, s Suite) {
	t.Helper()
	check(t, s)
}

func check(t testing.TB, s Suite) {
	t.Helper()
	mode, err := enginetest.RunMode()
	if err != nil {
		t.Fatalf("%v", err)
	}
	if mode == enginetest.ModeVerify && compat.Updating() {
		t.Fatalf("-update-compat cannot be combined with %s=verify: verification writes nothing", enginetest.RecordEnv)
	}
	role := s.PinRole
	if role == "" {
		role = "default"
	}
	pin, err := enginetest.LoadPin(s.Pins, role)
	if err != nil {
		t.Fatalf("%v", err)
	}
	recording := mode == enginetest.ModeRecord
	verifying := mode == enginetest.ModeVerify
	run := enginetest.RunCommand(t)
	record := enginetest.RecordCommand(t)
	changed := false
	f, err := ReadFile(s.Path)
	switch {
	case errors.Is(err, os.ErrNotExist) && recording:
		f, changed = &File{}, true
	case errors.Is(err, os.ErrNotExist):
		t.Fatalf("no %s. Record it:\n  %s", s.Path, record)
	case err != nil:
		t.Fatalf("%v", err)
	}

	merged, added, removed := merge(f.Rows, s.Queries)
	if added+removed > 0 {
		if !recording {
			t.Fatalf("%s: its rows are not the suite's queries: %d to add, %d to remove. Record them:\n  %s", s.Path, added, removed, record)
		}
		t.Logf("%s: %d rows added, %d removed", s.Path, added, removed)
		f.Rows, changed = merged, true
	}

	if f.Pin != pin.String() {
		if !recording {
			t.Fatalf("%s was recorded at pin %q, but %s pins %q. Record it:\n  %s", s.Path, f.Pin, s.Pins, pin, record)
		}
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
	tally := map[string]*counts{}
	var moved, closed, shifted, regressed []string
	for i := 0; i < len(f.Rows); {
		// The rows of one query are consecutive, one per aspect.
		first := f.Rows[i]
		fx := fixtures[first.Fixture]
		engine, err := ask(fx, first.Session, first.Statement)
		if err != nil {
			t.Fatalf("%s: the engine gave no answer: %v\n    %s", first.where(s.Path), err, first.Statement)
		}
		omni := s.Omni(fx, first.Session, first.Statement)
		if *show != "" && strings.Contains(first.Fixture+"\t"+first.Session+"\t"+first.Statement, *show) {
			t.Logf("%s\n  omni:   %+v\n  engine: %+v", first.Statement, omni, engine)
		}
		for ; i < len(f.Rows) && sameQuery(f.Rows[i], first); i++ {
			r := &f.Rows[i]
			got := engine.Cell(r.Aspect)
			switch {
			case got == r.Expected:
			case recording && r.Expected == Unjudged:
				r.Expected, r.Source, changed = got, compat.SourceEngine, true
			case recording:
				moved = append(moved, fmt.Sprintf("%s: the engine now answers %s, and the row records %s at the same pin:\n    %s\n", r.where(s.Path), got, r.Expected, r.Statement))
			case verifying:
				moved = append(moved, fmt.Sprintf("%s: the engine now answers %s, and the row records %s:\n    %s\n", r.where(s.Path), got, r.Expected, r.Statement))
			default:
				moved = append(moved, fmt.Sprintf("%s: the tape answers %s, and the row records %s:\n    %s\n", r.where(s.Path), got, r.Expected, r.Statement))
			}

			now := omni.Cell(r.Aspect)
			if recording && r.Omni == Unjudged {
				r.Omni, changed = now, true
			}
			was := fmt.Sprintf("    the row records omni %s\n", r.Omni)
			switch {
			case now == r.Omni:
			case r.Omni == r.Expected:
				regressed = append(regressed, r.describe(s.Path, now)+was)
			case compat.Updating():
				r.Omni, changed = now, true
			case now == r.Expected:
				closed = append(closed, r.describe(s.Path, now))
			default:
				shifted = append(shifted, r.describe(s.Path, now)+was)
			}
			n := tally[r.Aspect]
			if n == nil {
				n = &counts{}
				tally[r.Aspect] = n
			}
			n.rows++
			if now == r.Expected {
				n.agree++
			}
		}
	}
	switch mode {
	case enginetest.ModeRecord:
		report(t, "an answer does not reproduce, so it cannot be replayed", moved)
	case enginetest.ModeVerify:
		report(t, "the engine's answer differs from the committed expected cell", moved)
	default:
		report(t, "the tape and the file disagree. Record them:\n  "+record, moved)
	}
	report(t, "a gap closed: omni now answers what the engine does. Commit it:\n  "+run+" -update-compat", closed)
	report(t, "a known difference is now another one. Commit it:\n  "+run+" -update-compat", shifted)
	report(t, "a regression: omni agreed with the engine and no longer does", regressed)
	logTally(t, s.Path, f.Pin, tally)

	if changed {
		if verifying {
			t.Fatalf("span: internal error: verify mode attempted to write %s", s.Path)
		}
		if err := f.Write(s.Path); err != nil {
			t.Fatalf("%v", err)
		}
		t.Logf("wrote %s", s.Path)
	}
}

func sameQuery(a, b Row) bool {
	return a.Fixture == b.Fixture && a.Session == b.Session && a.Statement == b.Statement
}

// merge returns the rows the suite's queries need, one per aspect, in the
// queries' order. A row the file already holds keeps its answers; a new one
// is unjudged. added and removed count the rows that differ.
func merge(rows []Row, queries []Query) (merged []Row, added, removed int) {
	old := map[rowKey]Row{}
	for _, r := range rows {
		old[r.key()] = r
	}
	for _, q := range queries {
		for _, aspect := range Aspects {
			r := Row{Omni: Unjudged, Expected: Unjudged, Source: noSource, Fixture: q.Fixture, Session: q.Session, Aspect: aspect, Statement: q.Statement}
			if o, ok := old[r.key()]; ok {
				r = o
				delete(old, r.key())
			} else {
				added++
			}
			merged = append(merged, r)
		}
	}
	return merged, added, len(old)
}

func (r Row) where(path string) string {
	if r.line > 0 {
		return fmt.Sprintf("%s:%d: %s in %s on %s", path, r.line, r.Aspect, r.Session, r.Fixture)
	}
	return fmt.Sprintf("%s: %s in %s on %s", path, r.Aspect, r.Session, r.Fixture)
}

// describe renders one row and omni's answer on it now.
func (r Row) describe(path, now string) string {
	return fmt.Sprintf("%s: omni now answers %s, the engine %s:\n    %s\n", r.where(path), now, r.Expected, r.Statement)
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

type counts struct{ rows, agree int }

// logTally prints how many rows of each aspect agree, in one block.
func logTally(t testing.TB, path, pin string, tally map[string]*counts) {
	t.Helper()
	var b strings.Builder
	fmt.Fprintf(&b, "%s at pin %s\n", path, pin)
	fmt.Fprintf(&b, "  %-8s %6s %6s %6s\n", "aspect", "rows", "agree", "differ")
	var sum counts
	for _, aspect := range Aspects {
		n := tally[aspect]
		if n == nil {
			continue
		}
		fmt.Fprintf(&b, "  %-8s %6d %6d %6d\n", aspect, n.rows, n.agree, n.rows-n.agree)
		sum.rows += n.rows
		sum.agree += n.agree
	}
	fmt.Fprintf(&b, "  %-8s %6d %6d %6d", "all", sum.rows, sum.agree, sum.rows-sum.agree)
	t.Logf("%s", b.String())
}
