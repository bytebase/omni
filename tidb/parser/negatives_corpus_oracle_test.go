package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The negatives corpus is the grammar-guided half of the §3.6 rejection suite:
// for each reachable grammar alternative, near-misses that break the production
// in one place (drop a required token, truncate, insert a stray token, swap a
// keyword). tidb/parser/cmd/grammarcorpus -mode negatives generates the
// statements and the Rule#alternative/mutation labels; these two tests fill the
// verdict columns, sharing the pingcap corpus's engine harness and tally.
//
// The expected verdict is reject/reject. The gating direction is
// over-acceptance: a row omni accepts and the server rejects is a bug to fix
// across the mysql/mariadb forks like a gap. omni rejecting what the server
// accepts is a §3.3 completeness gap, out of scope here.
var negativesCorpusPath = filepath.Join("testdata", "parser_y_negatives_corpus.tsv")

// negativeCase is one recorded negatives row: the verdict each side gave, the
// alternative and mutation it exercises (Rule#index/mutation), and the statement.
type negativeCase struct {
	line   int
	omni   bool
	engine bool
	label  string
	sql    string
}

func loadNegativesCorpus(t testing.TB) (header []string, cases []negativeCase) {
	t.Helper()
	b, err := os.ReadFile(negativesCorpusPath)
	if err != nil {
		t.Fatalf("read %s: %v", negativesCorpusPath, err)
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
			t.Fatalf("%s:%d: want 4 tab-separated fields, got %d", negativesCorpusPath, i+1, len(fields))
		}
		omni, err := recordedVerdict(fields[0])
		if err != nil {
			t.Fatalf("%s:%d: omni column: %v", negativesCorpusPath, i+1, err)
		}
		engine, err := recordedVerdict(fields[1])
		if err != nil {
			t.Fatalf("%s:%d: engine column: %v", negativesCorpusPath, i+1, err)
		}
		if fields[2] == "" {
			t.Fatalf("%s:%d: empty label", negativesCorpusPath, i+1)
		}
		if seen[fields[2]] {
			t.Fatalf("%s:%d: duplicate label %q", negativesCorpusPath, i+1, fields[2])
		}
		seen[fields[2]] = true
		cases = append(cases, negativeCase{line: i + 1, omni: omni, engine: engine, label: fields[2], sql: fields[3]})
	}
	if len(cases) == 0 {
		t.Fatalf("%s holds no cases", negativesCorpusPath)
	}
	return header, cases
}

func writeNegativesCorpus(t *testing.T, header []string, cases []negativeCase) {
	t.Helper()
	var b strings.Builder
	for _, line := range header {
		b.WriteString(line)
		b.WriteByte('\n')
	}
	for _, c := range cases {
		fmt.Fprintf(&b, "%s\t%s\t%s\t%s\n", verdictName(c.omni), verdictName(c.engine), c.label, c.sql)
	}
	if err := os.WriteFile(negativesCorpusPath, []byte(b.String()), 0o644); err != nil {
		t.Fatalf("write %s: %v", negativesCorpusPath, err)
	}
	t.Logf("wrote %s (%d cases)", negativesCorpusPath, len(cases))
}

// reportOverAcceptances lists the rows omni accepts while the corpus records the
// server rejecting them — the over-acceptance work list for the fork fixes. It
// is capped like the drift report so closing a family at once stays readable.
func reportOverAcceptances(t *testing.T, cases []negativeCase) {
	t.Helper()
	var over []negativeCase
	for _, c := range cases {
		if c.omni && !c.engine {
			over = append(over, c)
		}
	}
	if len(over) == 0 {
		t.Logf("no over-acceptances recorded: omni rejects every negative the server rejects")
		return
	}
	t.Logf("%d over-acceptances recorded (omni accepts, TiDB rejects) — fix across mysql/mariadb/oceanbase:", len(over))
	for i, c := range over {
		if i >= corpusMaxReported {
			t.Logf("  ... and %d more", len(over)-corpusMaxReported)
			break
		}
		t.Logf("  %s:%d %s %q", negativesCorpusPath, c.line, c.label, c.sql)
	}
}

// TestTiDBNegativesCorpusOmniVerdicts holds tidb/parser to the omni column and
// reports the over-acceptances. It needs no container, so it is the gate a
// grammar change re-baselines: a parser edit that starts accepting a recorded
// reject shows here as drift. Once the recorded over-acceptances are closed, a
// new one appearing is a regression this test fails on.
func TestTiDBNegativesCorpusOmniVerdicts(t *testing.T) {
	header, cases := loadNegativesCorpus(t)

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
					negativesCorpusPath, c.line, verdictName(accepts), c.label, c.sql, verdictName(c.omni), err)
			}
			c.omni = accepts
		}
		tally.add(accepts, c.engine)
	}
	if drift > corpusMaxReported {
		t.Errorf("%d more rows drifted from the recorded omni verdict; rerun with -update and read the diff", drift-corpusMaxReported)
	}
	tally.report(t, negativesCorpusPath, "as recorded")
	reportOverAcceptances(t, cases)

	if *update {
		writeNegativesCorpus(t, header, cases)
	}
}

// TestTiDBNegativesCorpusEngineVerdicts measures omni and TiDB v8.5.5 against
// each other over the negatives corpus in one run, the same way the completeness
// corpus's engine test does, and holds the engine column to what the server
// says so a moved pin surfaces as drift to read.
func TestTiDBNegativesCorpusEngineVerdicts(t *testing.T) {
	if want := "pingcap/tidb:" + tidbCorpusVersion; tidbImage != want {
		t.Fatalf("the corpus engine column was recorded against %s, but this package pins %s; rerun with -update and read the drift", want, tidbImage)
	}
	tc := startTiDB(t)
	assertTiDBVersion(t, tc)

	header, cases := loadNegativesCorpus(t)

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
			unanswered = append(unanswered, fmt.Sprintf("%s:%d: %v", negativesCorpusPath, c.line, err))
			if tc.lost() {
				t.Errorf("no server left after %s:%d; stopping with %d of %d cases measured", negativesCorpusPath, c.line, i, len(cases))
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
					negativesCorpusPath, c.line, tidbCorpusVersion, verdictName(engineAccepts), c.label, c.sql, verdictName(c.engine))
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
		t.Errorf("%d rows drifted from the recorded omni verdict; TestTiDBNegativesCorpusOmniVerdicts names them", omniDrift)
	}
	tally.report(t, negativesCorpusPath, "live, via PREPARE")
	reportOverAcceptances(t, cases)
	for _, u := range unanswered {
		t.Errorf("no verdict: %s", u)
	}

	if *update && len(unanswered) == 0 {
		writeNegativesCorpus(t, header, cases)
	}
}
