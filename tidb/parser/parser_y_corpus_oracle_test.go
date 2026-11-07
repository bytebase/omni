package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The parser.y completeness corpus holds one statement per alternative of
// every rule in TiDB's grammar at v8.5.5: the shortest statement that reaches
// that alternative from the start symbol. Where the pingcap corpus is a sample
// of the statements pingcap happened to test, this one is derived from the
// grammar itself, so a clean run is evidence that omni covers every form the
// grammar has, not just the ones someone wrote a test for.
//
// harness/tidb-grammar generates the statements and the Rule#alternative
// labels; these two tests fill the verdict columns, exactly as the pingcap
// corpus's pair does, and share that pair's engine harness (syntaxVerdict,
// prepareVerdict) and tally. Regenerate the statements when the pin moves; the
// diff is the list of forms the new release added or removed.
var parserYCorpusPath = filepath.Join("testdata", "parser_y_corpus.tsv")

// parserYCase is one row: the verdicts each side gave, the grammar alternative
// it exercises (Rule#index), and the statement.
type parserYCase struct {
	line   int
	omni   bool
	engine bool
	label  string
	sql    string
}

// loadParserYCorpus reads the four-column corpus, keeping the comment header
// for an -update run to write back. Like the pingcap corpus it splits on '\n'
// so a statement's own bytes survive, and the statement stays the last field
// so a space or tab in it, were there one, would not be mistaken for a column.
func loadParserYCorpus(t testing.TB) (header []string, cases []parserYCase) {
	t.Helper()
	b, err := os.ReadFile(parserYCorpusPath)
	if err != nil {
		t.Fatalf("read %s: %v", parserYCorpusPath, err)
	}
	seen := map[string]bool{}
	for i, raw := range strings.Split(string(b), "\n") {
		switch {
		case raw == "":
			continue
		case strings.HasPrefix(raw, "#"):
			header = append(header, raw)
			continue
		}
		fields := strings.SplitN(raw, "\t", 4)
		if len(fields) != 4 {
			t.Fatalf("%s:%d: want 4 tab-separated fields, got %d", parserYCorpusPath, i+1, len(fields))
		}
		omni, err := recordedVerdict(fields[0])
		if err != nil {
			t.Fatalf("%s:%d: omni column: %v", parserYCorpusPath, i+1, err)
		}
		engine, err := recordedVerdict(fields[1])
		if err != nil {
			t.Fatalf("%s:%d: engine column: %v", parserYCorpusPath, i+1, err)
		}
		if fields[2] == "" {
			t.Fatalf("%s:%d: empty rule label", parserYCorpusPath, i+1)
		}
		if seen[fields[2]] {
			t.Fatalf("%s:%d: duplicate rule label %q", parserYCorpusPath, i+1, fields[2])
		}
		seen[fields[2]] = true
		cases = append(cases, parserYCase{line: i + 1, omni: omni, engine: engine, label: fields[2], sql: fields[3]})
	}
	if len(cases) == 0 {
		t.Fatalf("%s holds no cases", parserYCorpusPath)
	}
	return header, cases
}

// writeParserYCorpus rewrites the corpus from the verdicts a run measured,
// keeping the header and the generator's row order.
func writeParserYCorpus(t *testing.T, header []string, cases []parserYCase) {
	t.Helper()
	var b strings.Builder
	for _, line := range header {
		b.WriteString(line)
		b.WriteByte('\n')
	}
	for _, c := range cases {
		fmt.Fprintf(&b, "%s\t%s\t%s\t%s\n", verdictName(c.omni), verdictName(c.engine), c.label, c.sql)
	}
	if err := os.WriteFile(parserYCorpusPath, []byte(b.String()), 0o644); err != nil {
		t.Fatalf("write %s: %v", parserYCorpusPath, err)
	}
	t.Logf("wrote %s (%d cases)", parserYCorpusPath, len(cases))
}

// TestTiDBParserYCorpusOmniVerdicts holds tidb/parser to the omni column and
// reports how much of the grammar it covers. It needs no container, so it is
// the gate a grammar change re-baselines: change the parser, rerun with
// -update, and the diff shows which alternatives the change moved.
func TestTiDBParserYCorpusOmniVerdicts(t *testing.T) {
	header, cases := loadParserYCorpus(t)

	var tally corpusTally
	var drift int
	for i := range cases {
		c := &cases[i]
		_, err := Parse(c.sql)
		accepts := err == nil
		if accepts != c.omni {
			drift++
			if drift <= corpusMaxReported {
				t.Errorf("%s:%d: omni now %ss %s %q, corpus records %s (err=%v)",
					parserYCorpusPath, c.line, verdictName(accepts), c.label, c.sql, verdictName(c.omni), err)
			}
			c.omni = accepts
		}
		tally.add(accepts, c.engine)
	}
	if drift > corpusMaxReported {
		t.Errorf("%d more rows drifted from the recorded omni verdict; rerun with -update and read the diff", drift-corpusMaxReported)
	}
	tally.report(t, parserYCorpusPath, "as recorded")

	if *update {
		writeParserYCorpus(t, header, cases)
	}
}

// TestTiDBParserYCorpusEngineVerdicts measures omni and TiDB v8.5.5 against
// each other over the whole corpus in one run, the same way the pingcap
// corpus's engine test does, and holds the engine column to what the server
// says so a moved pin surfaces as drift to read.
func TestTiDBParserYCorpusEngineVerdicts(t *testing.T) {
	if want := "pingcap/tidb:" + tidbCorpusVersion; tidbImage != want {
		t.Fatalf("the corpus engine column was recorded against %s, but this package pins %s; rerun with -update and read the drift", want, tidbImage)
	}
	tc := startTiDB(t)
	assertTiDBVersion(t, tc)

	header, cases := loadParserYCorpus(t)

	var tally corpusTally
	var omniDrift, engineDrift int
	var unanswered []string
	for i := range cases {
		c := &cases[i]

		_, omniErr := Parse(c.sql)
		omniAccepts := omniErr == nil
		if omniAccepts != c.omni {
			omniDrift++
			c.omni = omniAccepts
		}

		previous := "<start of corpus>"
		if i > 0 {
			previous = cases[i-1].sql
		}
		engineAccepts, recovered, err := tc.syntaxVerdict(t, c.sql, previous)
		if recovered {
			tally.recovered++
		}
		if err != nil {
			tally.addNoVerdict()
			unanswered = append(unanswered, fmt.Sprintf("%s:%d: %v", parserYCorpusPath, c.line, err))
			if tc.lost() {
				t.Errorf("no server left after %s:%d; stopping with %d of %d cases measured", parserYCorpusPath, c.line, i, len(cases))
				break
			}
			if len(unanswered) > 10 {
				t.Errorf("gave up after %d statements with no verdict", len(unanswered))
				break
			}
			continue
		}

		if engineAccepts != c.engine {
			engineDrift++
			if engineDrift <= corpusMaxReported {
				t.Errorf("%s:%d: TiDB %s now %ss %s %q, corpus records %s",
					parserYCorpusPath, c.line, tidbCorpusVersion, verdictName(engineAccepts), c.label, c.sql, verdictName(c.engine))
			}
			c.engine = engineAccepts
		}
		tally.add(omniAccepts, engineAccepts)
	}

	if engineDrift > corpusMaxReported {
		t.Errorf("%d more rows drifted from the recorded engine verdict; the pinned image decides this column, so rerun with -update only after reading the drift",
			engineDrift-corpusMaxReported)
	}
	if omniDrift > 0 {
		t.Errorf("%d rows drifted from the recorded omni verdict; TestTiDBParserYCorpusOmniVerdicts names them", omniDrift)
	}
	tally.report(t, parserYCorpusPath, "live, via PREPARE")
	for _, u := range unanswered {
		t.Errorf("no verdict: %s", u)
	}

	if *update && len(unanswered) == 0 {
		writeParserYCorpus(t, header, cases)
	}
}
