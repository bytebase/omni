package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The integrationtest corpus (§3.4 sample 2) holds statements lifted from
// TiDB's own tests/integrationtest/t/**/*.test files: SQL people wrote to
// exercise the server, the closest public stand-in for user SQL. Where the
// pingcap corpus is grammar forms tested one at a time and the rest corpus is
// the forms pingcap's parser tests combine, this one leans on the DML and query
// shapes the other two under-sample.
//
// These two tests fill the verdict columns exactly as the other corpora's pairs
// do, and share the pingcap corpus's engine harness (syntaxVerdict,
// prepareVerdict), tally (corpusTally), verdict plumbing (recordedVerdict,
// verdictName) and the rest corpus's TIDB_CORPUS_N sampler (restCorpusSample).
// The omni-verdict test reads every row; the engine-verdict test measures a
// TIDB_CORPUS_N sample on a plain run and every row under -update, so this large
// corpus stays bounded on CI yet is always recorded whole.
var intCorpusPath = filepath.Join("testdata", "tidb_integrationtest_corpus.tsv")

// intCorpusCase is one recorded row: the verdict each side gave and the
// statement.
type intCorpusCase struct {
	line   int
	omni   bool
	engine bool
	sql    string
}

// loadIntCorpus reads the corpus, keeping the comment header for an -update run
// to write back. Like the other corpora it splits on '\n' and keeps the
// statement as the last field, so a statement's own bytes survive.
func loadIntCorpus(t testing.TB) (header []string, cases []intCorpusCase) {
	t.Helper()
	b, err := os.ReadFile(intCorpusPath)
	if err != nil {
		t.Fatalf("read %s: %v", intCorpusPath, err)
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
		fields := strings.SplitN(raw, "\t", 3)
		if len(fields) != 3 {
			t.Fatalf("%s:%d: want 3 tab-separated fields, got %d", intCorpusPath, i+1, len(fields))
		}
		omni, err := recordedVerdict(fields[0])
		if err != nil {
			t.Fatalf("%s:%d: omni column: %v", intCorpusPath, i+1, err)
		}
		engine, err := recordedVerdict(fields[1])
		if err != nil {
			t.Fatalf("%s:%d: engine column: %v", intCorpusPath, i+1, err)
		}
		if seen[fields[2]] {
			t.Fatalf("%s:%d: duplicate statement %q", intCorpusPath, i+1, fields[2])
		}
		seen[fields[2]] = true
		cases = append(cases, intCorpusCase{line: i + 1, omni: omni, engine: engine, sql: fields[2]})
	}
	if len(cases) == 0 {
		t.Fatalf("%s holds no cases", intCorpusPath)
	}
	return header, cases
}

// writeIntCorpus rewrites the corpus from the verdicts a run measured, keeping
// the header and the row order.
func writeIntCorpus(t *testing.T, header []string, cases []intCorpusCase) {
	t.Helper()
	var b strings.Builder
	for _, line := range header {
		b.WriteString(line)
		b.WriteByte('\n')
	}
	for _, c := range cases {
		fmt.Fprintf(&b, "%s\t%s\t%s\n", verdictName(c.omni), verdictName(c.engine), c.sql)
	}
	if err := os.WriteFile(intCorpusPath, []byte(b.String()), 0o644); err != nil {
		t.Fatalf("write %s: %v", intCorpusPath, err)
	}
	t.Logf("wrote %s (%d cases)", intCorpusPath, len(cases))
}

// TestTiDBIntegrationtestCorpusOmniVerdicts holds tidb/parser to the omni column
// over every row and reports how much of this real-SQL surface it covers. It
// needs no container, so it is the gate a grammar change re-baselines: close a
// gap, rerun with -update, and the diff shows which statements the change
// bought.
func TestTiDBIntegrationtestCorpusOmniVerdicts(t *testing.T) {
	header, cases := loadIntCorpus(t)

	var tally corpusTally
	var drift int
	for i := range cases {
		c := &cases[i]
		_, err := Parse(c.sql)
		accepts := err == nil
		if accepts != c.omni {
			drift++
			if drift <= corpusMaxReported {
				t.Errorf("%s:%d: omni now %ss %q, corpus records %s (err=%v)",
					intCorpusPath, c.line, verdictName(accepts), c.sql, verdictName(c.omni), err)
			}
			c.omni = accepts
		}
		tally.add(accepts, c.engine)
	}
	if drift > corpusMaxReported {
		t.Errorf("%d more rows drifted from the recorded omni verdict; rerun with -update and read the diff", drift-corpusMaxReported)
	}
	tally.report(t, intCorpusPath, "as recorded")

	if *update {
		writeIntCorpus(t, header, cases)
	}
}

// TestTiDBIntegrationtestCorpusEngineVerdicts measures omni and TiDB v8.5.5
// against each other, the same way the other corpora's engine tests do, and
// holds the engine column to what the server says so a moved pin surfaces as
// drift to read. It measures the TIDB_CORPUS_N sample on a plain run and every
// row under -update, so a recorded file is always measured whole.
func TestTiDBIntegrationtestCorpusEngineVerdicts(t *testing.T) {
	if want := "pingcap/tidb:" + tidbCorpusVersion; tidbImage != want {
		t.Fatalf("the corpus engine column was recorded against %s, but this package pins %s; rerun with -update and read the drift", want, tidbImage)
	}
	tc := startTiDB(t)
	assertTiDBVersion(t, tc)

	header, cases := loadIntCorpus(t)
	sample := restCorpusSample(t, len(cases))

	var tally corpusTally
	var omniDrift, engineDrift int
	var unanswered []string
	for n, i := range sample {
		c := &cases[i]

		_, omniErr := Parse(c.sql)
		omniAccepts := omniErr == nil
		if omniAccepts != c.omni {
			omniDrift++
			c.omni = omniAccepts
		}

		previous := "<start of corpus>"
		if n > 0 {
			previous = cases[sample[n-1]].sql
		}
		engineAccepts, recovered, err := tc.syntaxVerdict(t, c.sql, previous)
		if recovered {
			tally.recovered++
		}
		if err != nil {
			tally.addNoVerdict()
			unanswered = append(unanswered, fmt.Sprintf("%s:%d: %v", intCorpusPath, c.line, err))
			if tc.lost() {
				t.Errorf("no server left after %s:%d; stopping with %d of %d cases measured", intCorpusPath, c.line, n, len(sample))
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
				t.Errorf("%s:%d: TiDB %s now %ss %q, corpus records %s",
					intCorpusPath, c.line, tidbCorpusVersion, verdictName(engineAccepts), c.sql, verdictName(c.engine))
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
		t.Errorf("%d rows drifted from the recorded omni verdict; TestTiDBIntegrationtestCorpusOmniVerdicts names them", omniDrift)
	}
	tally.report(t, intCorpusPath, "live, via PREPARE")
	for _, u := range unanswered {
		t.Errorf("no verdict: %s", u)
	}

	// Writing the file needs every row measured; -update forces the full set,
	// so a non-empty unanswered list here means a run that lost its server.
	if *update && len(unanswered) == 0 {
		writeIntCorpus(t, header, cases)
	}
}
