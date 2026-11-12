package main

import (
	"flag"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var update = flag.Bool("update", false, "rewrite the checked-in coverage report from the current sources")

// Paths relative to this package directory. The grammar and the TiDB source
// tree sit above it; the report lives in this package's testdata.
const (
	testGrammar = "../../../testdata/parser.y"
	testTree    = "../../.." // the tidb/ top-level directory
	testReport  = "testdata/parser_y_coverage.md"
)

// TestRefCitationsResolve is the gate: every "// Ref: ... parser.y:N" citation in
// omni's TiDB productions must name a line that falls inside a real v8.5.5 rule.
// A line number past the grammar section, or landing in the gap between rules,
// is a stale or mistyped citation and fails here. This needs no server.
func TestRefCitationsResolve(t *testing.T) {
	a := mustAudit(t)
	if len(a.Broken) == 0 {
		t.Logf("%d citations across %d rules all resolve to a v8.5.5 rule", a.NumCits, len(a.Rules))
		return
	}
	for _, b := range a.Broken {
		t.Errorf("%s:%d: %s (range %s)", b.File, b.Line, b.Reason, b.RawRange)
	}
}

// TestCoverageReport keeps the checked-in coverage report in step with the
// sources. The report lists the rules no production cites, which guides the
// grammar work; -update rewrites it. It does not gate on the coverage level —
// an uncited rule is a lead, not a failure (design doc section 3.6).
func TestCoverageReport(t *testing.T) {
	a := mustAudit(t)
	got := renderReport(a)

	if *update {
		if err := os.WriteFile(testReport, []byte(got), 0o644); err != nil {
			t.Fatalf("write %s: %v", testReport, err)
		}
		t.Logf("wrote %s (%d uncited rules)", testReport, len(a.uncitedRules()))
		return
	}

	want, err := os.ReadFile(testReport)
	if err != nil {
		t.Fatalf("read %s: %v (run with -update to create it)", testReport, err)
	}
	if string(want) != got {
		t.Errorf("%s is stale; rerun: go test ./tidb/parser/cmd/refaudit -run TestCoverageReport -update", testReport)
	}
}

func mustAudit(t *testing.T) Audit {
	t.Helper()
	a, err := audit(testGrammar, testTree)
	if err != nil {
		t.Fatalf("audit: %v", err)
	}
	if len(a.Rules) == 0 {
		t.Fatal("no grammar rules parsed")
	}
	return a
}

// TestParseGrammarRules pins the rule-span parsing on a tiny grammar, so the
// line-to-rule resolution the audit depends on is checked without the 16k-line
// file. The two %% markers are lines 1 and 6; the grammar section is 2-5.
func TestParseGrammarRules(t *testing.T) {
	src := "%%\n" + // line 1
		"Start:\n" + // line 2
		"\t\"A\"\n" + // line 3
		"Other:\n" + // line 4
		"\t\"B\"\n" + // line 5
		"%%\n" // line 6
	rules, s, e, err := parseGrammarRules(src)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if s != 2 || e != 5 {
		t.Errorf("section bounds = %d-%d, want 2-5", s, e)
	}
	want := []Rule{{Name: "Start", Start: 2, End: 3}, {Name: "Other", Start: 4, End: 5}}
	if !reflect.DeepEqual(rules, want) {
		t.Errorf("rules = %+v, want %+v", rules, want)
	}
	for _, c := range []struct {
		line int
		rule string
	}{{2, "Start"}, {3, "Start"}, {4, "Other"}, {5, "Other"}} {
		if r := resolveLine(rules, c.line); r == nil || r.Name != c.rule {
			t.Errorf("line %d resolved to %v, want %s", c.line, r, c.rule)
		}
	}
	if r := resolveLine(rules, 99); r != nil {
		t.Errorf("line 99 resolved to %v, want nil", r)
	}
}

// TestExtractCitations pins the comment parsing on the awkward real forms: a
// single range, several ranges under one anchor, a range trailed by a
// parenthetical with its own comma, and lines that are not citations.
func TestExtractCitations(t *testing.T) {
	cases := []struct {
		text string
		want []struct{ s, e int }
	}{
		{"// Ref: TiDB v8.5.5 parser.y:11204-11501 (AdminStmt).", []struct{ s, e int }{{11204, 11501}}},
		{"// Ref: TiDB v8.5.5 parser.y:7573 (Literal), lexer.go", []struct{ s, e int }{{7573, 7573}}},
		{"// Ref: TiDB v8.5.5 parser.y:13764-13769, 2068-2084.", []struct{ s, e int }{{13764, 13769}, {2068, 2084}}},
		{"// Ref: parser.y:5623-5747 (BRIEStmt, BRIETables), 5759-6039 (the", []struct{ s, e int }{{5623, 5747}, {5759, 6039}}},
		{"// a plain comment, no citation", nil},
		{"// Ref: AlterTableSpec in pkg/parser/parser.y: the DROP arms", nil}, // no line number
		{"// Ref: https://dev.mysql.com/doc/ (not parser.y)", nil},
	}
	for _, c := range cases {
		got := extractCitations("f.go", 1, c.text)
		if len(got) != len(c.want) {
			t.Errorf("%q: got %d citations, want %d (%+v)", c.text, len(got), len(c.want), got)
			continue
		}
		for i, w := range c.want {
			if got[i].Start != w.s || got[i].End != w.e {
				t.Errorf("%q: citation %d = %d-%d, want %d-%d", c.text, i, got[i].Start, got[i].End, w.s, w.e)
			}
		}
	}
}

// TestNameOnlyReferenceCredited pins that a rule cited only by name — "Ref:
// AnalyzeOption in ... parser.y." with no line number — counts as cited. Before
// the fix these were dropped, so the coverage report listed them as uncovered.
func TestNameOnlyReferenceCredited(t *testing.T) {
	rules := []Rule{{Name: "AnalyzeOption", Start: 10, End: 12}, {Name: "Lonely", Start: 13, End: 15}}
	cits, names := scanFile("f.go", "package p\n// Ref: AnalyzeOption in pingcap/tidb pkg/parser/parser.y.\n")
	if len(cits) != 0 {
		t.Errorf("name-only reference yielded %d line citations, want 0", len(cits))
	}
	named := map[string]bool{}
	for _, n := range names {
		named[n] = true
	}
	a := runAudit(rules, nil, named)
	if !a.Cited["AnalyzeOption"] {
		t.Errorf("AnalyzeOption not credited by its name-only reference")
	}
	if a.Cited["Lonely"] {
		t.Errorf("Lonely credited though nothing references it")
	}

	// End to end: the four rules Codex named must leave the uncited list.
	full := mustAudit(t)
	for _, name := range []string{"AnalyzeOption", "AttributesOpt", "StatsOptionsOpt", "PartitionIntervalOpt"} {
		if !full.Cited[name] {
			t.Errorf("%s still reads as uncited after name crediting", name)
		}
	}
}

// TestScanFileMultilineBlock pins that a Ref comment wrapping across several //
// lines is read whole: a range or anchor on a continuation line (no "Ref:" of
// its own) still takes part in validation and coverage. Before the fix these
// were dropped, so a stale continuation range slipped the gate and a rule cited
// only on a continuation line read as uncited.
func TestScanFileMultilineBlock(t *testing.T) {
	// A continuation range under an earlier anchor, separated by a comma and by
	// "and", each range attributed to the line it appears on.
	src := "package p\n" +
		"// Ref: TiDB v8.5.5 parser.y:15733 (CREATE), 15743 (ALTER),\n" + // line 2
		"// 16445 (CALIBRATE), 16543 and 16631 (QUERY WATCH)\n" + // line 3
		"func f() {}\n" +
		"// Ref: parser.y:5903-5916 (A),\n" + // line 5
		"// parser.y:8768-8836 (B).\n" + // line 6 — anchor on a continuation line
		"func g() {}\n"
	cits, named := scanFile("f.go", src)

	wantStarts := map[int]int{15733: 2, 15743: 2, 16445: 3, 16543: 3, 16631: 3, 5903: 5, 8768: 6}
	if len(cits) != len(wantStarts) {
		t.Fatalf("got %d citations, want %d: %+v", len(cits), len(wantStarts), cits)
	}
	for _, c := range cits {
		if ln, ok := wantStarts[c.Start]; !ok || c.Line != ln {
			t.Errorf("citation %d-%d on line %d, want start on line %d", c.Start, c.End, c.Line, wantStarts[c.Start])
		}
	}
	// Rule-like words from across the block are collected for coverage.
	hasQueryWatch := false
	for _, n := range named {
		if n == "CREATE" {
			hasQueryWatch = true
		}
	}
	if !hasQueryWatch {
		t.Errorf("named words not collected from the block: %v", named)
	}
}

// TestReportPathExists guards against a refactor that moves the report out from
// under the test's relative path.
func TestReportPathExists(t *testing.T) {
	if _, err := os.Stat(filepath.Dir(testReport)); err != nil {
		t.Fatalf("report directory missing: %v", err)
	}
}
