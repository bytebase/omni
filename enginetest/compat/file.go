package compat

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bytebase/omni/enginetest/grammar"
)

// Verdicts a cell holds. A judge that reports where it stopped writes
// reject@<byte> instead of Reject.
const (
	Accept = "accept"
	Reject = "reject"
	// Undocumented is a docs judge's answer to a near-miss: the docs say
	// what is valid, not what is not. It counts as neither a gap nor a
	// regression, whatever omni says.
	Undocumented = "undocumented"
	// Unjudged marks a cell no recording has filled yet.
	Unjudged = "?"
)

// Sources of the expected verdict, strongest first. A row no judge has
// answered has source "-".
const (
	SourceEngine  = "engine"
	SourceGrammar = "grammar"
	SourceDocs    = "docs"
	noSource      = "-"
)

const (
	pinPrefix     = "# pin: "
	columnsHeader = "# columns: omni, expected, source, rule, statement"
)

// A Row is one statement and the verdicts recorded for it.
type Row struct {
	Omni      string // omni's verdict when the row was last updated
	Expected  string // the judge's verdict
	Source    string // which judge gave Expected
	Rule      string // the grammar alternative the row came from, or a label such as issue:BYT-1 or suite:pingcap
	Statement string
	line      int // 1-based line in the file it was read from; 0 for a row generated in this run
}

// hand reports whether a person added the row rather than the generator.
// A hand label starts with a word and a colon (issue:, suite:); a generated
// one with a rule name and a slash, and no rule name holds a colon, though a
// symbol after the slash can.
func (r Row) hand() bool {
	i := strings.IndexAny(r.Rule, ":/")
	return i >= 0 && r.Rule[i] == ':'
}

// generatedAlternative reports whether the generator made the row for an
// alternative, as opposed to a near-miss or a hand-added row.
func (r Row) generatedAlternative() bool { return !r.hand() && !grammar.IsNearMiss(r.Rule) }

// kind groups rows for the progress report.
func (r Row) kind() string {
	switch {
	case r.hand():
		label, _, _ := strings.Cut(r.Rule, ":")
		if label == "suite" {
			// A suite is a corpus with its own name: report each one.
			return r.Rule
		}
		return label
	case grammar.IsNearMiss(r.Rule):
		return "near-miss"
	default:
		return "alternative"
	}
}

type rowKey struct{ rule, statement string }

func (r Row) key() rowKey { return rowKey{r.Rule, r.Statement} }

// A File is a compat.tsv.
type File struct {
	Pin      string   // the pin the expected column was recorded at, as enginetest.Pin.String renders it
	Comments []string // the header's other comment lines, kept as written
	Rows     []Row
}

// ReadFile reads and validates a compat.tsv. The file is split on '\n'
// alone, and the statement is the last column, so a statement keeps every
// byte it has: trailing spaces, a CR inside a literal, a tab, invalid UTF-8.
func ReadFile(path string) (*File, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(path, b)
}

// Parse reads and validates the text of a compat.tsv; name is where it came
// from, for errors.
func Parse(name string, b []byte) (*File, error) {
	path := name
	f := &File{}
	seen := map[rowKey]int{}
	for i, raw := range strings.Split(string(b), "\n") {
		line := i + 1
		switch {
		case raw == "":
			continue
		case strings.HasPrefix(raw, "#"):
			if len(f.Rows) > 0 {
				return nil, fmt.Errorf("%s:%d: a comment after the first row; comments belong in the header", path, line)
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
		cols := strings.SplitN(raw, "\t", 5)
		if len(cols) != 5 {
			return nil, fmt.Errorf("%s:%d: want 5 tab-separated columns (omni, expected, source, rule, statement), got %d", path, line, len(cols))
		}
		r := Row{Omni: cols[0], Expected: cols[1], Source: cols[2], Rule: cols[3], Statement: cols[4], line: line}
		if err := r.validate(); err != nil {
			return nil, fmt.Errorf("%s:%d: %v", path, line, err)
		}
		if prev, dup := seen[r.key()]; dup {
			return nil, fmt.Errorf("%s:%d: repeats line %d (%s)", path, line, prev, r.Rule)
		}
		seen[r.key()] = line
		f.Rows = append(f.Rows, r)
	}
	return f, nil
}

func (r Row) validate() error {
	if err := validVerdict(r.Omni); err != nil {
		return fmt.Errorf("omni: %v", err)
	}
	if r.Omni == Undocumented {
		return fmt.Errorf("omni: %q is a judge's answer, not omni's", Undocumented)
	}
	if err := validVerdict(r.Expected); err != nil {
		return fmt.Errorf("expected: %v", err)
	}
	switch r.Source {
	case SourceEngine, SourceGrammar, SourceDocs:
		if r.Expected == Unjudged {
			return fmt.Errorf("source %s with no expected verdict; an unjudged row has source %s", r.Source, noSource)
		}
	case noSource:
		if r.Expected != Unjudged {
			return fmt.Errorf("expected %s with no source", r.Expected)
		}
	default:
		return fmt.Errorf("source %q: want engine, grammar, docs, or %s", r.Source, noSource)
	}
	if r.Rule == "" {
		return fmt.Errorf("empty rule")
	}
	if strings.ContainsAny(r.Rule, " \t") {
		return fmt.Errorf("rule %q holds a space", r.Rule)
	}
	return nil
}

func validVerdict(v string) error {
	switch v {
	case Accept, Reject, Undocumented, Unjudged:
		return nil
	}
	if at, ok := strings.CutPrefix(v, Reject+"@"); ok {
		if n, err := strconv.Atoi(at); err == nil && n >= 0 && strconv.Itoa(n) == at {
			return nil
		}
	}
	return fmt.Errorf("%q: want accept, reject, reject@<byte>, %s, or %s", v, Undocumented, Unjudged)
}

// accepts reports whether a verdict is an acceptance.
func accepts(v string) bool { return v == Accept }

// rejects reports whether a verdict is a rejection, with or without a position.
func rejects(v string) bool { return v == Reject || strings.HasPrefix(v, Reject+"@") }

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
		if err := r.validate(); err != nil {
			return fmt.Errorf("%s: %v", r.Rule, err)
		}
		fmt.Fprintf(&b, "%s\t%s\t%s\t%s\t%s\n", r.Omni, r.Expected, r.Source, r.Rule, r.Statement)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}
