package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// The "rest" corpus holds the statements in pingcap's own parser tests that the
// single-statement pingcap_parser_corpus.tsv did not already carry: the src of
// the []testCase / []testErrMsgCase elements it skipped, and the statement each
// ast/*_test.go restore case parses (a fragment inside a full-statement
// template). Where the pingcap corpus is forms tested one at a time, these are
// forms combined the way people write them, so they catch a reading that is
// right for a form alone but wrong for the form in context.
//
// These two tests fill the verdict columns exactly as the pingcap corpus's pair
// does, and share that pair's engine harness (syntaxVerdict, prepareVerdict),
// tally (corpusTally), and verdict plumbing (recordedVerdict, verdictName). The
// omni-verdict test reads every row; the engine-verdict test can run a sample,
// sized by TIDB_CORPUS_N, so the real-SQL corpora (this one and the larger
// §3.4 samples to come) stay bounded on every CI run.
var restCorpusPath = filepath.Join("testdata", "pingcap_parser_rest_corpus.tsv")

// restCorpusSampleEnv names the knob that caps how many rows the live engine
// test measures. Unset or non-positive means every row, as -update always does.
const restCorpusSampleEnv = "TIDB_CORPUS_N"

// restCorpusCase is one recorded row: the verdict each side gave and the
// statement.
type restCorpusCase struct {
	line   int
	omni   bool
	engine bool
	sql    string
}

// loadRestCorpus reads the corpus, keeping the comment header for an -update run
// to write back. Like the pingcap corpus it splits on '\n' and keeps the
// statement as the last field, so a statement's own bytes survive.
func loadRestCorpus(t testing.TB) (header []string, cases []restCorpusCase) {
	t.Helper()
	b, err := os.ReadFile(restCorpusPath)
	if err != nil {
		t.Fatalf("read %s: %v", restCorpusPath, err)
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
			t.Fatalf("%s:%d: want 3 tab-separated fields, got %d", restCorpusPath, i+1, len(fields))
		}
		omni, err := recordedVerdict(fields[0])
		if err != nil {
			t.Fatalf("%s:%d: omni column: %v", restCorpusPath, i+1, err)
		}
		engine, err := recordedVerdict(fields[1])
		if err != nil {
			t.Fatalf("%s:%d: engine column: %v", restCorpusPath, i+1, err)
		}
		if seen[fields[2]] {
			t.Fatalf("%s:%d: duplicate statement %q", restCorpusPath, i+1, fields[2])
		}
		seen[fields[2]] = true
		cases = append(cases, restCorpusCase{line: i + 1, omni: omni, engine: engine, sql: fields[2]})
	}
	if len(cases) == 0 {
		t.Fatalf("%s holds no cases", restCorpusPath)
	}
	return header, cases
}

// writeRestCorpus rewrites the corpus from the verdicts a run measured, keeping
// the header and the row order.
func writeRestCorpus(t *testing.T, header []string, cases []restCorpusCase) {
	t.Helper()
	var b strings.Builder
	for _, line := range header {
		b.WriteString(line)
		b.WriteByte('\n')
	}
	for _, c := range cases {
		fmt.Fprintf(&b, "%s\t%s\t%s\n", verdictName(c.omni), verdictName(c.engine), c.sql)
	}
	if err := os.WriteFile(restCorpusPath, []byte(b.String()), 0o644); err != nil {
		t.Fatalf("write %s: %v", restCorpusPath, err)
	}
	t.Logf("wrote %s (%d cases)", restCorpusPath, len(cases))
}

// restCorpusSample returns the indices the live engine test measures. With no
// TIDB_CORPUS_N, with an -update run, or with a cap at or above the row count,
// that is every row. Otherwise it is an evenly spread subset of n rows, chosen
// deterministically so the sample spans the whole file -- the corpus groups
// statements by source, so a sample that stopped short would never measure the
// later families -- and a rerun measures the same rows.
func restCorpusSample(t *testing.T, total int) []int {
	all := make([]int, total)
	for i := range all {
		all[i] = i
	}
	if *update {
		return all
	}
	raw := os.Getenv(restCorpusSampleEnv)
	if raw == "" {
		return all
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n <= 0 {
		t.Fatalf("%s=%q: want a positive integer", restCorpusSampleEnv, raw)
	}
	if n >= total {
		return all
	}
	// i*total/n, not i*(total/n): dividing first truncates the stride to a whole
	// number and packs the sample into the first rows, dropping the tail. Because
	// total > n here, the gap total/n exceeds 1, so the indices are strictly
	// increasing and span 0 .. total-total/n with no repeats.
	idx := make([]int, 0, n)
	for i := 0; i < n; i++ {
		idx = append(idx, i*total/n)
	}
	t.Logf("%s=%d: measuring %d of %d rows", restCorpusSampleEnv, n, len(idx), total)
	return idx
}

// TestTiDBRestCorpusOmniVerdicts holds tidb/parser to the omni column over every
// row and reports how much of this real-SQL surface it covers. It needs no
// container, so it is the gate a grammar change re-baselines: close a gap,
// rerun with -update, and the diff shows which statements the change bought.
func TestTiDBRestCorpusOmniVerdicts(t *testing.T) {
	header, cases := loadRestCorpus(t)

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
					restCorpusPath, c.line, verdictName(accepts), c.sql, verdictName(c.omni), err)
			}
			c.omni = accepts
		}
		tally.add(accepts, c.engine)
	}
	if drift > corpusMaxReported {
		t.Errorf("%d more rows drifted from the recorded omni verdict; rerun with -update and read the diff", drift-corpusMaxReported)
	}
	tally.report(t, restCorpusPath, "as recorded")

	if *update {
		writeRestCorpus(t, header, cases)
	}
}

// TestTiDBRestCorpusEngineVerdicts measures omni and TiDB v8.5.5 against each
// other, the same way the pingcap corpus's engine test does, and holds the
// engine column to what the server says so a moved pin surfaces as drift to
// read. It measures the TIDB_CORPUS_N sample on a plain run and every row under
// -update, so a recorded file is always measured whole.
func TestTiDBRestCorpusEngineVerdicts(t *testing.T) {
	if want := "pingcap/tidb:" + tidbCorpusVersion; tidbImage != want {
		t.Fatalf("the corpus engine column was recorded against %s, but this package pins %s; rerun with -update and read the drift", want, tidbImage)
	}
	tc := startTiDB(t)
	assertTiDBVersion(t, tc)

	header, cases := loadRestCorpus(t)
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
			unanswered = append(unanswered, fmt.Sprintf("%s:%d: %v", restCorpusPath, c.line, err))
			if tc.lost() {
				t.Errorf("no server left after %s:%d; stopping with %d of %d cases measured", restCorpusPath, c.line, n, len(sample))
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
					restCorpusPath, c.line, tidbCorpusVersion, verdictName(engineAccepts), c.sql, verdictName(c.engine))
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
		t.Errorf("%d rows drifted from the recorded omni verdict; TestTiDBRestCorpusOmniVerdicts names them", omniDrift)
	}
	tally.report(t, restCorpusPath, "live, via PREPARE")
	for _, u := range unanswered {
		t.Errorf("no verdict: %s", u)
	}

	// Writing the file needs every row measured; -update forces the full set,
	// so a non-empty unanswered list here means a run that lost its server.
	if *update && len(unanswered) == 0 {
		writeRestCorpus(t, header, cases)
	}
}
