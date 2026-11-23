package span

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/compat"
)

// The aspects of a query's span a row compares, in the order a query's rows
// are written.
const (
	// Results is the query's output column names, in order.
	Results = "results"
	// Lineage is, for each output column in order, the columns its values
	// come from. A column read through a view names the view's column and
	// every column the view's definition reads for it.
	Lineage = "lineage"
	// Tables is the tables and views the query reads.
	Tables = "tables"
)

// Aspects lists every aspect, in the order a query's rows are written.
var Aspects = []string{Results, Lineage, Tables}

// Unjudged marks a cell no recording has filled yet.
const Unjudged = compat.Unjudged

const (
	noSource      = "-"
	pinPrefix     = "# pin: "
	columnsHeader = "# columns: omni, expected, source, fixture, session, aspect, statement"
)

// A Row is one aspect of one query's span, in one session on one fixture,
// as omni and the engine answer it.
type Row struct {
	Omni      string // omni's answer when the row was last updated
	Expected  string // the engine's answer
	Source    string // engine, or - before a recording fills Expected
	Fixture   string // the starting schema, a file of the suite's fixture directory
	Session   string // the session's catalog and schema, which unqualified names resolve in
	Aspect    string // one of Aspects
	Statement string
	line      int // 1-based line in the file it was read from; 0 for a row made in this run
}

type rowKey struct{ fixture, session, aspect, statement string }

func (r Row) key() rowKey { return rowKey{r.Fixture, r.Session, r.Aspect, r.Statement} }

// Line reports the 1-based line of the row in the file it was read from, or
// 0 for a row made in memory.
func (r Row) Line() int { return r.line }

// A File is a span.tsv.
type File struct {
	Pin      string   // the pin the expected column was recorded at, as enginetest.Pin.String renders it
	Comments []string // the header's other comment lines, kept as written
	Rows     []Row
}

// ReadFile reads and validates a span.tsv.
func ReadFile(path string) (*File, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(path, b)
}

// Parse reads and validates the text of a span.tsv; name is where it came
// from, for errors. As in a compat.tsv, a row is split on '\n' alone and the
// statement is the last column, so it keeps every byte it has.
func Parse(name string, b []byte) (*File, error) {
	f := &File{}
	seen := map[rowKey]int{}
	for i, raw := range strings.Split(string(b), "\n") {
		line := i + 1
		switch {
		case raw == "":
			continue
		case strings.HasPrefix(raw, "#"):
			if len(f.Rows) > 0 {
				return nil, fmt.Errorf("%s:%d: a comment after the first row; comments belong in the header", name, line)
			}
			switch {
			case strings.HasPrefix(raw, pinPrefix):
				f.Pin = strings.TrimPrefix(raw, pinPrefix)
			case raw == columnsHeader:
			default:
				f.Comments = append(f.Comments, raw)
			}
			continue
		}
		cols := strings.SplitN(raw, "\t", 7)
		if len(cols) != 7 {
			return nil, fmt.Errorf("%s:%d: want 7 tab-separated columns (omni, expected, source, fixture, session, aspect, statement), got %d", name, line, len(cols))
		}
		r := Row{Omni: cols[0], Expected: cols[1], Source: cols[2], Fixture: cols[3], Session: cols[4], Aspect: cols[5], Statement: cols[6], line: line}
		if err := r.validate(); err != nil {
			return nil, fmt.Errorf("%s:%d: %v", name, line, err)
		}
		if prev, dup := seen[r.key()]; dup {
			return nil, fmt.Errorf("%s:%d: repeats line %d (%s of a query in %s on %s)", name, line, prev, r.Aspect, r.Session, r.Fixture)
		}
		seen[r.key()] = line
		f.Rows = append(f.Rows, r)
	}
	return f, nil
}

func (r Row) validate() error {
	if r.Omni == "" || r.Expected == "" {
		return fmt.Errorf("an empty cell; an unanswered one holds %s", Unjudged)
	}
	switch r.Source {
	case compat.SourceEngine:
		if r.Expected == Unjudged {
			return fmt.Errorf("source %s with no expected answer; an unjudged row has source %s", r.Source, noSource)
		}
	case noSource:
		if r.Expected != Unjudged {
			return fmt.Errorf("expected %q with no source", r.Expected)
		}
	default:
		return fmt.Errorf("source %q: want %s or %s", r.Source, compat.SourceEngine, noSource)
	}
	if !slices.Contains(Aspects, r.Aspect) {
		return fmt.Errorf("aspect %q: want one of %s", r.Aspect, strings.Join(Aspects, ", "))
	}
	for what, v := range map[string]string{"fixture": r.Fixture, "session": r.Session} {
		if v == "" {
			return fmt.Errorf("empty %s", what)
		}
		if strings.ContainsAny(v, " \t") {
			return fmt.Errorf("%s %q holds a space", what, v)
		}
	}
	if strings.TrimSpace(r.Statement) == "" {
		return fmt.Errorf("empty statement")
	}
	return nil
}

// Write writes the file: the pin, the columns line, the other comments, and
// the rows in order. It refuses to write in verification mode
// (OMNI_RECORD=verify).
func (f *File) Write(path string) error {
	mode, err := enginetest.RunMode()
	if err != nil {
		return err
	}
	if mode == enginetest.ModeVerify {
		return fmt.Errorf("span: cannot write %s with %s=verify: verification writes nothing", path, enginetest.RecordEnv)
	}
	var b strings.Builder
	b.WriteString(pinPrefix + f.Pin + "\n")
	b.WriteString(columnsHeader + "\n")
	for _, c := range f.Comments {
		b.WriteString(c + "\n")
	}
	for _, r := range f.Rows {
		for _, v := range []string{r.Omni, r.Expected, r.Statement} {
			if strings.ContainsAny(v, "\t\n") {
				return fmt.Errorf("%s of %q: a cell holds a tab or a newline: %q", r.Aspect, r.Statement, v)
			}
		}
		fmt.Fprintf(&b, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", r.Omni, r.Expected, r.Source, r.Fixture, r.Session, r.Aspect, r.Statement)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

// Guard compares a span.tsv before and after a change and returns what the
// change may not do, as simulate.Guard does for a simulate.tsv:
//
//   - An omni cell that held its row's expected answer keeps it: a regression
//     is never committed as omni's new answer. A cell that differed may
//     change to anything, since a fix can move a row from one difference to
//     another, and the test holds the cell to what omni does.
//   - An expected cell (and its source) changes only from ?, or when the
//     file was recorded at a new pin.
//
// Rows added or removed are not compared.
func Guard(before, after *File) []string {
	repinned := PinChanged(before, after)
	old := map[rowKey]Row{}
	for _, r := range before.Rows {
		old[r.key()] = r
	}
	var problems []string
	for _, r := range after.Rows {
		o, ok := old[r.key()]
		if !ok {
			continue
		}
		if r.Omni != o.Omni && o.Omni == o.Expected && o.Omni != Unjudged && r.Omni != r.Expected {
			problems = append(problems, fmt.Sprintf("%s in %s on %s: omni changed from %s, the expected answer, to %s: a regression\n    %s",
				r.Aspect, r.Session, r.Fixture, o.Omni, r.Omni, r.Statement))
		}
		if (r.Expected != o.Expected || r.Source != o.Source) && o.Expected != Unjudged && !repinned {
			problems = append(problems, fmt.Sprintf("%s in %s on %s: expected changed from %s to %s without a pin bump; only a recording at a new pin changes it\n    %s",
				r.Aspect, r.Session, r.Fixture, o.Expected, r.Expected, r.Statement))
		}
	}
	return problems
}

// PinChanged reports whether after was recorded at a different pin than before.
func PinChanged(before, after *File) bool { return before.Pin != after.Pin }
