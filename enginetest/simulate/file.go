package simulate

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/bytebase/omni/enginetest/compat"
)

// Outcomes a cell holds besides an error code, which is the engine's number
// for the refusal, in decimal.
const (
	// Took is omni taking a statement the engine refuses, or storing another
	// schema than the engine does. Only omni's cell holds it: the engine's
	// acceptance always carries the hash of what it stored, "ok:<hash>", and
	// omni's cell holds that same text when it stores the same schema.
	Took = "ok"
	// NoCode is omni refusing a statement with an error that has no code.
	NoCode = "error"
	// Unjudged marks a cell no recording has filled yet.
	Unjudged = compat.Unjudged
)

const (
	noSource      = "-"
	pinPrefix     = "# pin: "
	columnsHeader = "# columns: omni, expected, source, fixture, rule, statement"
)

// A Form is one object as the engine stores it after a statement: what its
// SHOW CREATE, or the engine's like of it, prints.
type Form struct {
	Kind string // TABLE, VIEW, ...
	Name string
	Text string
}

// An Answer is the engine's answer to one statement on one fixture.
type Answer struct {
	// Code is the error the engine refused the statement with, in decimal.
	// It is empty when the engine took the statement.
	Code string
	// Forms is what the fixture's database holds after a statement the
	// engine took, in name order.
	Forms []Form
}

// Refused reports whether the engine refused the statement.
func (a Answer) Refused() bool { return a.Code != "" }

// Text renders the stored forms as one text, for a reader.
func (a Answer) Text() string {
	var b strings.Builder
	for _, f := range a.Forms {
		fmt.Fprintf(&b, "-- %s %s\n%s\n", f.Kind, f.Name, f.Text)
	}
	return b.String()
}

// Outcome is the cell an answer is recorded as: the error code, or ok: and
// the first 12 hex digits of the SHA-256 of the stored forms.
//
// madeUp matches the names the engine makes up, such as a constraint written
// without one, and may be nil. Each distinct name counts from 1 in the order
// the forms hold it, so the hash is of the row's own text and stays the same
// whichever stand-in a tape gave the name.
func (a Answer) Outcome(madeUp *regexp.Regexp) string {
	if a.Refused() {
		return a.Code
	}
	text := a.Text()
	if madeUp != nil {
		seen := map[string]int{}
		text = madeUp.ReplaceAllStringFunc(text, func(name string) string {
			n, ok := seen[name]
			if !ok {
				n = len(seen) + 1
				seen[name] = n
			}
			return fmt.Sprintf("<made-up %d>", n)
		})
	}
	return fmt.Sprintf("%s:%x", Took, sha256.Sum256([]byte(text)))[:len(Took)+1+12]
}

// A Row is one statement on one fixture and the outcomes recorded for it.
type Row struct {
	Omni      string // omni's outcome when the row was last updated
	Expected  string // the engine's outcome
	Source    string // which judge gave Expected: engine, or - for none yet
	Fixture   string // the starting schema, a file of the suite's fixture directory
	Rule      string // the compat.tsv rule the statement came from, or a label such as issue:BYT-1
	Statement string
	line      int // 1-based line in the file it was read from; 0 for a row generated in this run
}

// hand reports whether a person added the row rather than the generator, as
// compat.Row tells them apart: a hand label is a word and a colon.
func (r Row) hand() bool {
	i := strings.IndexAny(r.Rule, ":/")
	return i >= 0 && r.Rule[i] == ':'
}

type rowKey struct{ fixture, rule, statement string }

func (r Row) key() rowKey { return rowKey{r.Fixture, r.Rule, r.Statement} }

// A File is a simulate.tsv.
type File struct {
	Pin      string   // the pin the expected column was recorded at, as enginetest.Pin.String renders it
	Comments []string // the header's other comment lines, kept as written
	Rows     []Row
}

// ReadFile reads and validates a simulate.tsv. Like a compat.tsv it is split
// on '\n' alone and the statement is the last column, so a statement keeps
// every byte it has.
func ReadFile(path string) (*File, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(path, b)
}

// Parse reads and validates the text of a simulate.tsv; name is where it
// came from, for errors.
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
		cols := strings.SplitN(raw, "\t", 6)
		if len(cols) != 6 {
			return nil, fmt.Errorf("%s:%d: want 6 tab-separated columns (omni, expected, source, fixture, rule, statement), got %d", name, line, len(cols))
		}
		r := Row{Omni: cols[0], Expected: cols[1], Source: cols[2], Fixture: cols[3], Rule: cols[4], Statement: cols[5], line: line}
		if err := r.validate(); err != nil {
			return nil, fmt.Errorf("%s:%d: %v", name, line, err)
		}
		if prev, dup := seen[r.key()]; dup {
			return nil, fmt.Errorf("%s:%d: repeats line %d (%s on %s)", name, line, prev, r.Rule, r.Fixture)
		}
		seen[r.key()] = line
		f.Rows = append(f.Rows, r)
	}
	return f, nil
}

func (r Row) validate() error {
	switch {
	case r.Omni == Unjudged, r.Omni == Took, r.Omni == NoCode, isCode(r.Omni), isStored(r.Omni):
	default:
		return fmt.Errorf("omni: %q: want an error code, %s, %s:<hash>, %s, or %s", r.Omni, Took, Took, NoCode, Unjudged)
	}
	switch {
	case r.Expected == Unjudged, isCode(r.Expected), isStored(r.Expected):
	default:
		return fmt.Errorf("expected: %q: want an error code, %s:<hash>, or %s", r.Expected, Took, Unjudged)
	}
	switch r.Source {
	case compat.SourceEngine:
		if r.Expected == Unjudged {
			return fmt.Errorf("source %s with no expected outcome; an unjudged row has source %s", r.Source, noSource)
		}
	case noSource:
		if r.Expected != Unjudged {
			return fmt.Errorf("expected %s with no source", r.Expected)
		}
	default:
		return fmt.Errorf("source %q: want %s or %s", r.Source, compat.SourceEngine, noSource)
	}
	for what, v := range map[string]string{"fixture": r.Fixture, "rule": r.Rule} {
		if v == "" {
			return fmt.Errorf("empty %s", what)
		}
		if strings.ContainsAny(v, " \t") {
			return fmt.Errorf("%s %q holds a space", what, v)
		}
	}
	return nil
}

// isCode reports whether a cell is an error code: a decimal number as the
// engine writes it.
func isCode(v string) bool {
	n, err := strconv.Atoi(v)
	return err == nil && n > 0 && strconv.Itoa(n) == v
}

var stored = regexp.MustCompile(`^ok:[0-9a-f]{12}$`)

// isStored reports whether a cell is an acceptance with the stored schema's hash.
func isStored(v string) bool { return stored.MatchString(v) }

// takes reports whether a cell is an acceptance, whatever was stored.
func takes(v string) bool { return v == Took || isStored(v) }

// Write writes the file: the pin, the columns line, the other comments, and
// the rows in order.
func (f *File) Write(path string) error {
	var b strings.Builder
	b.WriteString(pinPrefix + f.Pin + "\n")
	b.WriteString(columnsHeader + "\n")
	for _, c := range f.Comments {
		b.WriteString(c + "\n")
	}
	for _, r := range f.Rows {
		if strings.Contains(r.Statement, "\n") {
			return fmt.Errorf("%s: statement holds a newline: %q", r.Rule, r.Statement)
		}
		fmt.Fprintf(&b, "%s\t%s\t%s\t%s\t%s\t%s\n", r.Omni, r.Expected, r.Source, r.Fixture, r.Rule, r.Statement)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

// Guard compares a simulate.tsv before and after a change and returns what
// the change may not do. A test run sees only the file it has; these rules
// need the file it replaced.
//
//   - An omni cell that held its row's expected outcome keeps it: a
//     regression is never committed as omni's new outcome. A cell that
//     differed may change to anything, since a fix can move a row from one
//     difference to another, and the test holds the cell to what omni does.
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
			problems = append(problems, fmt.Sprintf("%s on %s: omni changed from %s, the expected outcome, to %s: a regression\n    %s",
				r.Rule, r.Fixture, o.Omni, r.Omni, r.Statement))
		}
		if (r.Expected != o.Expected || r.Source != o.Source) && o.Expected != Unjudged && !repinned {
			problems = append(problems, fmt.Sprintf("%s on %s: expected changed from %s to %s without a pin bump; only a recording at a new pin changes it\n    %s",
				r.Rule, r.Fixture, o.Expected, r.Expected, r.Statement))
		}
	}
	return problems
}

// PinChanged reports whether after was recorded at a different pin than before.
func PinChanged(before, after *File) bool { return before.Pin != after.Pin }
