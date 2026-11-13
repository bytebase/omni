package compat

import (
	"fmt"
	"os"
	"strings"
)

// An Override is one row of an engine's grammar-overrides.tsv: a statement a
// recognizer over the grammar judges wrongly, the verdict the engine gives,
// and what says so (§2.3 of
// docs/specs/2026-11-11-engine-compatibility-without-servers-design.md).
// Its source says which of two things that is:
//
//   - SourceGrammar: the grammar decides it where the recognizer does not
//     read, in the lexer's code or in a grammar action. The cite is the
//     grammar's line.
//   - SourceDocs: the vendor's docs take a statement the grammar does not
//     derive, as Cosmos DB's take an alias with no AS and DynamoDB's take
//     forms the PartiQL reference grammar lacks. The cite is the page. Docs
//     say what is valid and not what is not, so such a row only accepts.
type Override struct {
	Verdict   string // Accept or Reject
	Source    string // SourceGrammar or SourceDocs
	Cite      string // the grammar line, lexer rule or docs page that decides it
	Statement string
	line      int
}

// ReadOverrides reads a grammar-overrides.tsv. Each row is verdict, source,
// cite and statement, tab-separated; the statement is last, so it keeps
// every byte. A row needs a cite: an override nothing vouches for is a cell
// edited by hand.
func ReadOverrides(path string) ([]Override, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var out []Override
	seen := map[string]int{}
	for i, raw := range strings.Split(string(b), "\n") {
		line := i + 1
		if raw == "" || strings.HasPrefix(raw, "#") {
			continue
		}
		cols := strings.SplitN(raw, "\t", 4)
		if len(cols) != 4 {
			return nil, fmt.Errorf("%s:%d: want 4 tab-separated columns (verdict, source, cite, statement), got %d", path, line, len(cols))
		}
		o := Override{Verdict: cols[0], Source: cols[1], Cite: cols[2], Statement: cols[3], line: line}
		if o.Verdict != Accept && o.Verdict != Reject {
			return nil, fmt.Errorf("%s:%d: verdict %q: want %s or %s", path, line, o.Verdict, Accept, Reject)
		}
		switch o.Source {
		case SourceGrammar:
		case SourceDocs:
			if o.Verdict != Accept {
				return nil, fmt.Errorf("%s:%d: docs say what is valid, not what is not: a %s override only accepts", path, line, SourceDocs)
			}
		default:
			return nil, fmt.Errorf("%s:%d: source %q: want %s or %s", path, line, o.Source, SourceGrammar, SourceDocs)
		}
		if strings.TrimSpace(o.Cite) == "" {
			return nil, fmt.Errorf("%s:%d: no cite: name the grammar line, lexer rule or docs page that decides the statement", path, line)
		}
		if prev, dup := seen[o.Statement]; dup {
			return nil, fmt.Errorf("%s:%d: repeats line %d", path, line, prev)
		}
		seen[o.Statement] = line
		out = append(out, o)
	}
	return out, nil
}

// overridden wraps a recognizer with the overrides: a statement one names
// gets its verdict and source. It returns the judge and the problems that
// make the file wrong for these rows: an override no row holds, which
// decides nothing, and one the recognizer already agrees with, which is no
// miss.
func overridden(recognize recognizeFunc, path string, overrides []Override, rows []Row) (recognizeFunc, []string) {
	byStatement := map[string]Override{}
	for _, o := range overrides {
		byStatement[o.Statement] = o
	}
	held := map[string]Row{}
	for _, r := range rows {
		held[r.Statement] = r
	}
	var problems []string
	for _, o := range overrides {
		r, ok := held[o.Statement]
		if !ok {
			problems = append(problems, fmt.Sprintf("%s:%d: no row of the file holds the statement, so the override decides nothing; delete it, or add the row:\n    %s\n", path, o.line, o.Statement))
			continue
		}
		if v, _, err := recognize(r); err == nil && v == o.Verdict {
			problems = append(problems, fmt.Sprintf("%s:%d: the recognizer already says %s, so this is no miss of its; delete the override:\n    %s\n", path, o.line, v, o.Statement))
		}
	}
	return func(r Row) (string, string, error) {
		if o, ok := byStatement[r.Statement]; ok {
			return o.Verdict, o.Source, nil
		}
		return recognize(r)
	}, problems
}
