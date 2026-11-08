package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The "bytebase" corpus is §3.4 sample 3 of the drop-pingcap design: the SQL
// Bytebase already promises to handle, extracted from the input side of its own
// 64 TiDB test files. Where the pingcap corpora are pingcap's grammar cases,
// these are the forms Bytebase's users and tests actually drive through the
// TiDB path, so an over-acceptance here is a form Bytebase ships against.
//
// These two tests fill the verdict columns exactly as the pingcap corpus's pair
// does, and share that pair's engine harness (syntaxVerdict, prepareVerdict),
// tally (corpusTally), verdict plumbing (recordedVerdict, verdictName), and the
// TIDB_CORPUS_N sampler (restCorpusSample). The omni-verdict test reads every
// row; the engine-verdict test can run a sample so the real-SQL corpora stay
// bounded on every CI run.
var bbCorpusPath = filepath.Join("testdata", "bytebase_tidb_corpus.tsv")

// bbCorpusCase is one recorded row: the verdict each side gave and the
// statement.
type bbCorpusCase struct {
	line   int
	omni   bool
	engine bool
	sql    string
}

// loadBBCorpus reads the corpus, keeping the comment header for an -update run
// to write back. Like the other corpora it splits on '\n' and keeps the
// statement as the last field, so a statement's own bytes survive.
func loadBBCorpus(t testing.TB) (header []string, cases []bbCorpusCase) {
	t.Helper()
	b, err := os.ReadFile(bbCorpusPath)
	if err != nil {
		t.Fatalf("read %s: %v", bbCorpusPath, err)
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
			t.Fatalf("%s:%d: want 3 tab-separated fields, got %d", bbCorpusPath, i+1, len(fields))
		}
		omni, err := recordedVerdict(fields[0])
		if err != nil {
			t.Fatalf("%s:%d: omni column: %v", bbCorpusPath, i+1, err)
		}
		engine, err := recordedVerdict(fields[1])
		if err != nil {
			t.Fatalf("%s:%d: engine column: %v", bbCorpusPath, i+1, err)
		}
		if seen[fields[2]] {
			t.Fatalf("%s:%d: duplicate statement %q", bbCorpusPath, i+1, fields[2])
		}
		seen[fields[2]] = true
		cases = append(cases, bbCorpusCase{line: i + 1, omni: omni, engine: engine, sql: fields[2]})
	}
	if len(cases) == 0 {
		t.Fatalf("%s holds no cases", bbCorpusPath)
	}
	return header, cases
}

// writeBBCorpus rewrites the corpus from the verdicts a run measured, keeping
// the header and the row order.
func writeBBCorpus(t *testing.T, header []string, cases []bbCorpusCase) {
	t.Helper()
	var b strings.Builder
	for _, line := range header {
		b.WriteString(line)
		b.WriteByte('\n')
	}
	for _, c := range cases {
		fmt.Fprintf(&b, "%s\t%s\t%s\n", verdictName(c.omni), verdictName(c.engine), c.sql)
	}
	if err := os.WriteFile(bbCorpusPath, []byte(b.String()), 0o644); err != nil {
		t.Fatalf("write %s: %v", bbCorpusPath, err)
	}
	t.Logf("wrote %s (%d cases)", bbCorpusPath, len(cases))
}

// TestTiDBBytebaseCorpusOmniVerdicts holds tidb/parser to the omni column over
// every row and reports how much of this real-SQL surface it covers. It needs
// no container, so it is the gate a grammar change re-baselines: close a gap,
// rerun with -update, and the diff shows which statements the change bought.
func TestTiDBBytebaseCorpusOmniVerdicts(t *testing.T) {
	header, cases := loadBBCorpus(t)

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
					bbCorpusPath, c.line, verdictName(accepts), c.sql, verdictName(c.omni), err)
			}
			c.omni = accepts
		}
		tally.add(accepts, c.engine)
	}
	if drift > corpusMaxReported {
		t.Errorf("%d more rows drifted from the recorded omni verdict; rerun with -update and read the diff", drift-corpusMaxReported)
	}
	tally.report(t, bbCorpusPath, "as recorded")

	if *update {
		writeBBCorpus(t, header, cases)
	}
}

// TestTiDBBytebaseCorpusEngineVerdicts measures omni and TiDB v8.5.5 against
// each other, the same way the pingcap corpus's engine test does, and holds the
// engine column to what the server says so a moved pin surfaces as drift to
// read. It measures the TIDB_CORPUS_N sample on a plain run and every row under
// -update, so a recorded file is always measured whole.
func TestTiDBBytebaseCorpusEngineVerdicts(t *testing.T) {
	if want := "pingcap/tidb:" + tidbCorpusVersion; tidbImage != want {
		t.Fatalf("the corpus engine column was recorded against %s, but this package pins %s; rerun with -update and read the drift", want, tidbImage)
	}
	tc := startTiDB(t)
	assertTiDBVersion(t, tc)

	header, cases := loadBBCorpus(t)
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
			unanswered = append(unanswered, fmt.Sprintf("%s:%d: %v", bbCorpusPath, c.line, err))
			if tc.lost() {
				t.Errorf("no server left after %s:%d; stopping with %d of %d cases measured", bbCorpusPath, c.line, n, len(sample))
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
					bbCorpusPath, c.line, tidbCorpusVersion, verdictName(engineAccepts), c.sql, verdictName(c.engine))
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
		t.Errorf("%d rows drifted from the recorded omni verdict; TestTiDBBytebaseCorpusOmniVerdicts names them", omniDrift)
	}
	tally.report(t, bbCorpusPath, "live, via PREPARE")
	for _, u := range unanswered {
		t.Errorf("no verdict: %s", u)
	}

	// Writing the file needs every row measured; -update forces the full set,
	// so a non-empty unanswered list here means a run that lost its server.
	if *update && len(unanswered) == 0 {
		writeBBCorpus(t, header, cases)
	}
}
