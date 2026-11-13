// Package compat holds an engine's parser to its compat.tsv: one row per
// statement, with omni's verdict and the judge's. A replay reads the file
// and needs no server; a recording (OMNI_RECORD=1) asks the judge for every
// row not judged yet and writes the answers.
package compat

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/grammar"
)

var updateCompat = flag.Bool("update-compat", false, "write omni cells that now equal their row's expected verdict")

// Updating reports whether the run was given -update-compat. enginetest/simulate
// reads it too: one flag commits the closed gaps of either file.
func Updating() bool { return *updateCompat }

// maxReported caps the rows one failure prints. Closing a grammar family
// moves dozens of rows at once, and a thousand like lines is not a report.
const maxReported = 25

// maxUnanswered is how many statements a recording lets the judge leave
// unanswered before it stops asking. A server that stopped answering will
// not start again.
const maxUnanswered = 10

// A Suite is one engine's compat.tsv and what checks it.
type Suite struct {
	// Path is the compat.tsv.
	Path string
	// Pins is the engine's pins.tsv; the default row is the pin Expected
	// was recorded at.
	Pins string
	// Grammar is the vendor's grammar, bison (.y) or ANTLR 4 (.g4), and
	// Tokens the engine's tokens.tsv. The generated rows of Path must be
	// what they generate.
	Grammar, Tokens string
	// Defined names the conditions Grammar is built with, where one file
	// holds more than one grammar between %ifdef lines, as MariaDB's does
	// (grammar.Conditional).
	Defined []string
	// Omni returns omni's verdict on one statement.
	Omni func(sql string) string
	// Source names the judge: SourceEngine, SourceGrammar or SourceDocs.
	Source string
	// Judge is called only in a recording, and returns the function that
	// asks the judge about one statement. It starts the server where there
	// is one, and fails t unless that server is pin. An error from the
	// returned function means no verdict, never a rejection.
	//
	// A SourceGrammar or SourceDocs suite has no Judge: a recognizer over
	// Grammar judges it (see recognizer), and since that needs no server,
	// every run asks it again and fails a row whose cell it no longer gives.
	Judge func(t testing.TB, pin enginetest.Pin) func(sql string) (string, error)
}

// Check holds omni to the suite's compat.tsv. A row passes when omni still
// gives the verdict its omni cell records: an agreement, or a known gap. A
// row fails when omni's verdict moved, and when it has not been judged. A
// row that agreed and no longer does is a regression. Any other move is
// committed with -update-compat: a gap that closed, and a known gap that is
// now another one, as when a parser that accepted what the judge rejects
// comes to reject it somewhere else (reject@<byte>).
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

	var g *grammar.Grammar
	if s.Grammar != "" {
		var generated []grammar.Row
		var uncovered []string
		g, generated, uncovered, err = generate(s.Grammar, s.Tokens, s.Defined)
		if err != nil {
			t.Fatalf("generate rows from %s: %v", s.Grammar, err)
		}
		if len(uncovered) > 0 {
			t.Logf("%d alternatives have no generated row and need a hand-added one: %s", len(uncovered), strings.Join(uncovered, ", "))
		}
		merged, added, removed := merge(f.Rows, generated)
		if len(added)+len(removed) > 0 {
			if !recording {
				t.Errorf("%s: the generated rows are not what %s generates: %d to add, %d to remove. Record them:\n  %s%s",
					s.Path, s.Grammar, len(added), len(removed), record, sample(added, removed))
			} else {
				t.Logf("%s: %d generated rows added, %d removed", s.Path, len(added), len(removed))
				f.Rows = merged
				changed = true
			}
		}
	}

	if f.Pin != pin.String() {
		if !recording {
			t.Errorf("%s was recorded at pin %q, but %s pins %q. Record it:\n  %s", s.Path, f.Pin, s.Pins, pin, record)
		} else {
			// Every expected cell describes the old server: ask again.
			for i := range f.Rows {
				f.Rows[i].Expected, f.Rows[i].Source = Unjudged, noSource
			}
			f.Pin = pin.String()
			changed = true
		}
	}

	var recognize func(Row) (string, error)
	if recognized(s.Source) {
		if s.Judge != nil || g == nil {
			t.Fatalf("a %s suite is judged by a recognizer over its Grammar: give the Grammar and no Judge", s.Source)
		}
		recognize = recognizer(g, s.Source)
	}

	if recording {
		n, err := judge(t, s, pin, f.Rows, recognize)
		if n > 0 {
			changed = true
		}
		if err != nil {
			t.Errorf("%v", err)
		}
	}

	// A recognizer needs no server, so its cells are checked on every run:
	// a grammar fix moves them, and a cell written by hand fails here.
	var moved []string
	if recognize != nil {
		rejudged := 0
		for i := range f.Rows {
			r := &f.Rows[i]
			if r.Source != s.Source {
				continue
			}
			v, err := recognize(*r)
			if err != nil {
				moved = append(moved, fmt.Sprintf("%s: %s: %v:\n    %s\n", s.Path, r.Rule, err, r.Statement))
				continue
			}
			if v == r.Expected {
				continue
			}
			if recording {
				r.Expected = v
				rejudged++
				continue
			}
			moved = append(moved, fmt.Sprintf("%s: the %s judge now says %s, the row records %s:\n    %s\n", describeRule(s.Path, *r), s.Source, v, r.Expected, r.Statement))
		}
		if rejudged > 0 {
			t.Logf("%s: the %s judge moved %d rows", s.Path, s.Source, rejudged)
			changed = true
		}
	}

	var tally tally
	var unjudged, closed, shifted, regressed []string
	for i := range f.Rows {
		r := &f.Rows[i]
		now := s.Omni(r.Statement)
		switch {
		case r.Expected == Unjudged || r.Omni == Unjudged:
			unjudged = append(unjudged, describe(s.Path, *r, now))
			tally.add(*r, now)
			continue
		case r.Expected == Undocumented:
			// The docs say nothing either way: whatever omni says passes,
			// and -update-compat keeps the cell current.
			if now != r.Omni && *updateCompat {
				r.Omni = now
				changed = true
			}
		case now == r.Omni:
		case r.Omni == r.Expected:
			regressed = append(regressed, describe(s.Path, *r, now)+fmt.Sprintf("    the row records omni %s\n", r.Omni))
		case *updateCompat:
			r.Omni = now
			changed = true
		case now == r.Expected:
			closed = append(closed, describe(s.Path, *r, now))
		default:
			shifted = append(shifted, describe(s.Path, *r, now)+fmt.Sprintf("    the row records omni %s\n", r.Omni))
		}
		tally.add(*r, now)
	}
	report(t, "not judged yet. Record them:\n  "+record, unjudged)
	report(t, "the judge's cell is not what its recognizer gives: the grammar changed, or the cell was edited. Record them:\n  "+record, moved)
	report(t, "a gap closed: omni now says what the judge says. Commit it:\n  "+run+" -update-compat", closed)
	report(t, "a known gap is now another one. Commit it:\n  "+run+" -update-compat", shifted)
	report(t, "a regression: omni agreed with the judge and no longer does", regressed)
	tally.log(t, s.Path, f.Pin)

	if changed {
		if err := f.Write(s.Path); err != nil {
			t.Fatalf("%v", err)
		}
		t.Logf("wrote %s", s.Path)
	}
}

// generate reads the grammar, by its extension, and returns its rows.
// defined is for a bison grammar's %ifdef conditions.
func generate(grammarPath, tokensPath string, defined []string) (*grammar.Grammar, []grammar.Row, []string, error) {
	tokens, err := grammar.ReadTokens(tokensPath)
	if err != nil {
		return nil, nil, nil, err
	}
	src, err := os.ReadFile(grammarPath)
	if err != nil {
		return nil, nil, nil, err
	}
	var g *grammar.Grammar
	if strings.HasSuffix(grammarPath, ".g4") {
		if len(defined) > 0 {
			return nil, nil, nil, fmt.Errorf("%s: Defined names %%ifdef conditions, which only a bison grammar has", grammarPath)
		}
		g, err = grammar.ReadANTLR(string(src), tokens)
	} else {
		g, err = grammar.ReadBison(string(src), tokens, defined...)
	}
	if err != nil {
		return nil, nil, nil, err
	}
	rows, uncovered := g.Rows()
	return g, rows, uncovered, nil
}

// recognized reports whether a source's judge is a recognizer over the
// grammar rather than a server.
func recognized(source string) bool { return source == SourceGrammar || source == SourceDocs }

// recognizer returns the judge of a grammar or docs suite. A statement the
// grammar derives is accepted. One it does not is rejected by a grammar,
// and undocumented by docs, which say what is valid and not what is not.
// A generated alternative the grammar does not derive has no verdict: the
// recognizer is reading the grammar otherwise than the generator.
func recognizer(g *grammar.Grammar, source string) func(Row) (string, error) {
	rec := g.Recognizer()
	return func(r Row) (string, error) {
		if rec.Accepts(r.Statement) {
			return Accept, nil
		}
		if r.generatedAlternative() {
			return "", fmt.Errorf("the recognizer does not derive the grammar's own statement")
		}
		if source == SourceDocs {
			return Undocumented, nil
		}
		return Reject, nil
	}
}

// merge returns the file's rows with its generated rows replaced by
// generated, in generated's order, followed by the hand-added rows in their
// order. A row the file already holds keeps its verdicts; a new one is
// unjudged. added and removed are the generated rows that differ.
func merge(rows []Row, generated []grammar.Row) (merged []Row, added, removed []Row) {
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
		r := Row{Omni: Unjudged, Expected: Unjudged, Source: noSource, Rule: g.Rule, Statement: g.Statement}
		if o, ok := old[r.key()]; ok {
			r = o
			delete(old, r.key())
		} else {
			added = append(added, r)
		}
		merged = append(merged, r)
	}
	for _, r := range rows {
		if _, gone := old[r.key()]; gone {
			removed = append(removed, r)
		}
	}
	return append(merged, hand...), added, removed
}

// sample names a few of the rows merge added and removed.
func sample(added, removed []Row) string {
	var b strings.Builder
	for i, r := range added {
		if i == 3 {
			break
		}
		fmt.Fprintf(&b, "\n  + %s\t%s", r.Rule, r.Statement)
	}
	for i, r := range removed {
		if i == 3 {
			break
		}
		fmt.Fprintf(&b, "\n  - %s\t%s", r.Rule, r.Statement)
	}
	return b.String()
}

// judge asks the judge about every unjudged row and fills in omni's cell
// where that is unjudged too. It returns how many rows it filled. A row the
// judge gave no verdict for stays unjudged, so a partial recording is safe
// to write and a rerun resumes it. recognize, when set, is the judge.
func judge(t testing.TB, s Suite, pin enginetest.Pin, rows []Row, recognize func(Row) (string, error)) (int, error) {
	t.Helper()
	var ask func(Row) (string, error)
	if recognize != nil {
		ask = recognize
	}
	filled, judged := 0, 0
	var unanswered []string
	todo := 0
	for _, r := range rows {
		if r.Expected == Unjudged {
			todo++
		}
	}
	for i := range rows {
		r := &rows[i]
		if r.Omni == Unjudged {
			r.Omni = s.Omni(r.Statement)
			filled++
		}
		if r.Expected != Unjudged || len(unanswered) >= maxUnanswered {
			continue
		}
		if ask == nil {
			if s.Judge == nil {
				return filled, fmt.Errorf("%s has %d unjudged rows and the suite has no judge", s.Path, todo)
			}
			server := s.Judge(t, pin)
			ask = func(r Row) (string, error) { return server(r.Statement) }
		}
		v, err := ask(*r)
		if err == nil {
			err = validVerdict(v)
		}
		if err == nil && (v == Unjudged || v == Undocumented && s.Source != SourceDocs) {
			err = fmt.Errorf("the judge answered %q", v)
		}
		if err != nil {
			unanswered = append(unanswered, fmt.Sprintf("%s %q: %v", r.Rule, r.Statement, err))
			continue
		}
		r.Expected, r.Source = v, s.Source
		filled++
		if judged++; judged%1000 == 0 {
			fmt.Fprintf(os.Stderr, "compat: %s: judged %d of %d rows\n", s.Path, judged, todo)
		}
	}
	if len(unanswered) > 0 {
		return filled, fmt.Errorf("the judge gave no verdict for %d statements (stopping at %d); they stay unjudged:\n  %s",
			len(unanswered), maxUnanswered, strings.Join(unanswered, "\n  "))
	}
	return filled, nil
}

// describeRule names a row by where it is and its rule.
func describeRule(path string, r Row) string {
	if r.line > 0 {
		return fmt.Sprintf("%s:%d: %s", path, r.line, r.Rule)
	}
	return fmt.Sprintf("%s: %s", path, r.Rule)
}

// describe renders one row and omni's verdict on it now.
func describe(path string, r Row, now string) string {
	where := path
	if r.line > 0 {
		where = fmt.Sprintf("%s:%d", path, r.line)
	}
	if r.Expected == Unjudged {
		return fmt.Sprintf("%s: %s: omni now says %s, and no judge has answered:\n    %s\n", where, r.Rule, now, r.Statement)
	}
	return fmt.Sprintf("%s: %s: omni now says %s, %s says %s:\n    %s\n", where, r.Rule, now, r.Source, r.Expected, r.Statement)
}

// report fails t with up to maxReported of rows, under what.
func report(t testing.TB, what string, rows []string) {
	t.Helper()
	if len(rows) == 0 {
		return
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d rows: %s\n", len(rows), what)
	for i, r := range rows {
		if i == maxReported {
			fmt.Fprintf(&b, "... and %d more\n", len(rows)-maxReported)
			break
		}
		b.WriteString(r)
	}
	t.Errorf("%s", b.String())
}

// tally counts each kind of row by how omni and the judge answer it.
type tally struct {
	kinds  []string
	byKind map[string]*counts
}

type counts struct {
	rows, agree, judgeOnly, omniOnly, other, undocumented, unjudged int
}

func (c *tally) add(r Row, now string) {
	if c.byKind == nil {
		c.byKind = map[string]*counts{}
	}
	k := r.kind()
	n, ok := c.byKind[k]
	if !ok {
		n = &counts{}
		c.byKind[k] = n
		c.kinds = append(c.kinds, k)
	}
	n.rows++
	switch {
	case r.Expected == Unjudged:
		n.unjudged++
	case r.Expected == Undocumented:
		n.undocumented++
	case now == r.Expected:
		n.agree++
	case accepts(r.Expected) && rejects(now):
		n.judgeOnly++
	case rejects(r.Expected) && accepts(now):
		n.omniOnly++
	default:
		n.other++
	}
}

// log prints the counts in one block, so a run answers how far the engine
// is from fully compatible without anything else being read.
func (c *tally) log(t testing.TB, path, pin string) {
	t.Helper()
	slices.Sort(c.kinds)
	var b strings.Builder
	fmt.Fprintf(&b, "%s at pin %s\n", path, pin)
	fmt.Fprintf(&b, "  %-28s %7s %7s %12s %10s %6s %12s %9s\n", "rows", "total", "agree", "judge only", "omni only", "other", "undocumented", "unjudged")
	var sum counts
	line := func(k string, n counts) string {
		return fmt.Sprintf("  %-28s %7d %7d %12d %10d %6d %12d %9d", k, n.rows, n.agree, n.judgeOnly, n.omniOnly, n.other, n.undocumented, n.unjudged)
	}
	for _, k := range c.kinds {
		n := c.byKind[k]
		b.WriteString(line(k, *n) + "\n")
		sum.rows += n.rows
		sum.agree += n.agree
		sum.judgeOnly += n.judgeOnly
		sum.omniOnly += n.omniOnly
		sum.other += n.other
		sum.undocumented += n.undocumented
		sum.unjudged += n.unjudged
	}
	b.WriteString(line("all", sum))
	t.Logf("%s", b.String())
}
