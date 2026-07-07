package main

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"
)

func fixtureRows() []Row {
	rows := []Row{
		{SQL: "SELECT 1", Expected: VerdictAccept, OmniVerdict: VerdictAccept, Family: "SELECT"},
		{SQL: "SELECT * FROM ((SELECT 1)) x", Expected: VerdictAccept, OmniVerdict: VerdictReject, OmniError: "syntax error at or near \"x\" (line 1, column 27)\nrelated text: SELECT * FROM ((SELECT 1)) x", Family: "SELECT"},
		{SQL: "SELECT * FROM ((SELECT 2)) y", Expected: VerdictAccept, OmniVerdict: VerdictReject, OmniError: "syntax error at or near \"y\" (line 1, column 27)\nrelated text: SELECT * FROM ((SELECT 2)) y", Family: "SELECT"},
		{SQL: "CREATE FUNCTION f() RETURNS INT RETURN 1", Expected: VerdictReject, OmniVerdict: VerdictAccept, Family: "CREATE OTHER"},
	}
	for i := range rows {
		rows[i].Engine, rows[i].Lane = "tidb", "upstream"
		rows[i].StmtHash = stmtHash(rows[i].SQL)
		classify(&rows[i])
	}
	return rows
}

func TestScoreboardDeterministic(t *testing.T) {
	meta := RunMeta{Engine: "tidb", EngineVersion: "v8.5.5", OmniSHA: "abc123", CorpusTag: "v8.5.5", ClassifierVersion: classifierVersion}
	a := renderScoreboard(meta, fixtureRows())
	b := renderScoreboard(meta, fixtureRows())
	if a != b {
		t.Error("scoreboard not deterministic")
	}
	if !strings.Contains(a, "| GAP | 2 |") {
		t.Errorf("expected 2 GAP statements in counts:\n%s", a)
	}
	// the two GAP rows share a normalized error -> ONE cluster
	if !strings.Contains(a, "GAP clusters: 1") {
		t.Errorf("expected 1 GAP cluster:\n%s", a)
	}
	// counts-table naming matches the wire fields
	if !strings.Contains(a, "| duplicates_dropped | 0 |") {
		t.Errorf("expected duplicates_dropped line in counts:\n%s", a)
	}
	if !strings.Contains(a, "| duplicate_label_conflicts | 0 |") {
		t.Errorf("expected duplicate_label_conflicts line in counts:\n%s", a)
	}
}

// TestRenderDeterministicUnderPermutation: rendering must not depend on row
// order — the exemplar tie-break must be total (length, then lexicographic).
// fixtureRows contains two equal-length SQLs in one GAP cluster.
func TestRenderDeterministicUnderPermutation(t *testing.T) {
	meta := RunMeta{Engine: "tidb", ClassifierVersion: classifierVersion}
	a := renderScoreboard(meta, fixtureRows())
	rev := fixtureRows()
	for i, j := 0, len(rev)-1; i < j; i, j = i+1, j-1 {
		rev[i], rev[j] = rev[j], rev[i]
	}
	b := renderScoreboard(meta, rev)
	if a != b {
		t.Errorf("render depends on row order:\n--- forward ---\n%s\n--- reversed ---\n%s", a, b)
	}
	// the tie (two 28-byte SQLs) must resolve to the lexicographically smallest
	if !strings.Contains(a, "SELECT * FROM ((SELECT 1)) x") {
		t.Errorf("expected lexicographically smallest tied exemplar:\n%s", a)
	}
}

// TestClusterExemplarEmptySQLTotal: an empty SQL is a legitimate corpus
// statement (the corpus has one) and the shortest possible exemplar; the
// tie-break must stay total (an empty-string sentinel would let later rows
// steal the slot depending on order).
func TestClusterExemplarEmptySQLTotal(t *testing.T) {
	forward := []Row{
		{Lane: "upstream", Class: ClassGap, Family: "UNKNOWN", DivergenceKey: "k2", SQL: "", SourcePath: "a_test.go", Line: 1},
		{Lane: "upstream", Class: ClassGap, Family: "UNKNOWN", DivergenceKey: "k2", SQL: "SELECT 1", SourcePath: "b_test.go", Line: 2},
	}
	reversed := []Row{forward[1], forward[0]}
	cf := clusterRows(forward, ClassGap)
	cr := clusterRows(reversed, ClassGap)
	if len(cf) != 1 || len(cr) != 1 {
		t.Fatalf("want 1 cluster each, got %d and %d", len(cf), len(cr))
	}
	if cf[0].Exemplar != "" || cf[0].ExemplarSrc != "a_test.go:1" {
		t.Errorf("forward exemplar = %q (%s), want empty SQL from a_test.go:1", cf[0].Exemplar, cf[0].ExemplarSrc)
	}
	if cr[0].Exemplar != cf[0].Exemplar || cr[0].ExemplarSrc != cf[0].ExemplarSrc {
		t.Errorf("exemplar order-dependent: forward %q (%s) vs reversed %q (%s)",
			cf[0].Exemplar, cf[0].ExemplarSrc, cr[0].Exemplar, cr[0].ExemplarSrc)
	}
}

// TestNonUpstreamLaneExcluded: non-upstream rows appear in neither the
// headline counts nor the cluster tables.
func TestNonUpstreamLaneExcluded(t *testing.T) {
	rows := fixtureRows()
	gen := Row{Engine: "tidb", Lane: "generated", SQL: "GENERATED LANE ONLY",
		Expected: VerdictAccept, OmniVerdict: VerdictReject,
		OmniError: "generated-lane-only failure mode", Family: "SELECT"}
	gen.StmtHash = stmtHash(gen.SQL)
	classify(&gen)
	rows = append(rows, gen)
	meta := RunMeta{Engine: "tidb", ClassifierVersion: classifierVersion}
	a := renderScoreboard(meta, rows)
	if !strings.Contains(a, "| GAP | 2 |") {
		t.Errorf("generated-lane row leaked into counts:\n%s", a)
	}
	if !strings.Contains(a, "GAP clusters: 1") {
		t.Errorf("generated-lane row leaked into clusters:\n%s", a)
	}
	if strings.Contains(a, "GENERATED LANE ONLY") {
		t.Errorf("generated-lane exemplar rendered:\n%s", a)
	}
}

// TestTruncateRuneSafe: truncation backs off to a rune boundary instead of
// emitting invalid UTF-8.
func TestTruncateRuneSafe(t *testing.T) {
	s := "aaaébbb" // é is bytes 3-4
	if got := truncate(s, 4); got != "aaa..." {
		t.Errorf("truncate mid-rune = %q, want %q", got, "aaa...")
	}
	if got := truncate(s, 5); got != "aaaé..." {
		t.Errorf("truncate at rune boundary = %q, want %q", got, "aaaé...")
	}
	if got := truncate(s, 8); got != s {
		t.Errorf("truncate no-op = %q, want %q", got, s)
	}
	if got := truncate(s, 4); !utf8.ValidString(got) {
		t.Errorf("truncated string is not valid UTF-8: %q", got)
	}
}

// TestWriteJSONLMetaFirstRowsSorted: the meta line is line 1; rows follow
// sorted by stmt_hash so run files diff cleanly across harvests.
func TestWriteJSONLMetaFirstRowsSorted(t *testing.T) {
	rows := []Row{
		{Engine: "tidb", Lane: "upstream", SQL: "SELECT 3", OmniVerdict: VerdictAccept},
		{Engine: "tidb", Lane: "upstream", SQL: "SELECT 1", OmniVerdict: VerdictAccept},
		{Engine: "tidb", Lane: "upstream", SQL: "SELECT 2", OmniVerdict: VerdictAccept},
	}
	for i := range rows {
		rows[i].StmtHash = stmtHash(rows[i].SQL)
	}
	meta := RunMeta{Engine: "tidb", ClassifierVersion: classifierVersion, DuplicatesDropped: 7}
	path := filepath.Join(t.TempDir(), "run.jsonl")
	if err := writeJSONL(path, meta, rows); err != nil {
		t.Fatal(err)
	}
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	if !sc.Scan() {
		t.Fatal("empty JSONL file")
	}
	var metaLine struct {
		Meta *RunMeta `json:"meta"`
	}
	if err := json.Unmarshal(sc.Bytes(), &metaLine); err != nil || metaLine.Meta == nil {
		t.Fatalf("first line is not a meta line: %s", sc.Text())
	}
	if metaLine.Meta.DuplicatesDropped != 7 {
		t.Errorf("meta duplicates_dropped = %d, want 7", metaLine.Meta.DuplicatesDropped)
	}
	var got []Row
	for sc.Scan() {
		var r Row
		if err := json.Unmarshal(sc.Bytes(), &r); err != nil {
			t.Fatalf("bad row line %q: %v", sc.Text(), err)
		}
		got = append(got, r)
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	if len(got) != len(rows) {
		t.Fatalf("read back %d rows, want %d", len(got), len(rows))
	}
	for i := 1; i < len(got); i++ {
		if got[i-1].StmtHash >= got[i].StmtHash {
			t.Errorf("rows not sorted by stmt_hash: %q before %q", got[i-1].StmtHash, got[i].StmtHash)
		}
	}
	seen := map[string]bool{}
	for _, r := range got {
		seen[r.SQL] = true
	}
	for _, s := range []string{"SELECT 1", "SELECT 2", "SELECT 3"} {
		if !seen[s] {
			t.Errorf("row %q lost in round-trip", s)
		}
	}
}

// TestClusterExemplarShortestSQL: the exemplar is the shortest SQL in the
// cluster, and the source pointer follows the exemplar.
func TestClusterExemplarShortestSQL(t *testing.T) {
	rows := []Row{
		{Lane: "upstream", Class: ClassGap, Family: "SELECT", DivergenceKey: "k", SQL: "SELECT 1 FROM longer_table", SourcePath: "a_test.go", Line: 10},
		{Lane: "upstream", Class: ClassGap, Family: "SELECT", DivergenceKey: "k", SQL: "SELECT 1", SourcePath: "b_test.go", Line: 20},
		{Lane: "upstream", Class: ClassGap, Family: "SELECT", DivergenceKey: "k", SQL: "SELECT 1 FROM t", SourcePath: "c_test.go", Line: 30},
	}
	cs := clusterRows(rows, ClassGap)
	if len(cs) != 1 {
		t.Fatalf("want 1 cluster, got %d: %+v", len(cs), cs)
	}
	if cs[0].Count != 3 {
		t.Errorf("cluster count = %d, want 3", cs[0].Count)
	}
	if cs[0].Exemplar != "SELECT 1" {
		t.Errorf("exemplar = %q, want shortest SQL %q", cs[0].Exemplar, "SELECT 1")
	}
	if cs[0].ExemplarSrc != "b_test.go:20" {
		t.Errorf("exemplar source = %q, want %q", cs[0].ExemplarSrc, "b_test.go:20")
	}
}

// TestOverClusterKeyedByLeadingTokens: pre-adjudication OVER rows have no
// engine error message, so they cluster on normalized leading tokens.
func TestOverClusterKeyedByLeadingTokens(t *testing.T) {
	rows := []Row{
		{SQL: "CREATE FUNCTION f1() RETURNS INT RETURN 1", Expected: VerdictReject, OmniVerdict: VerdictAccept, Family: "CREATE OTHER", Lane: "upstream"},
		{SQL: "CREATE FUNCTION f2() RETURNS INT RETURN 2", Expected: VerdictReject, OmniVerdict: VerdictAccept, Family: "CREATE OTHER", Lane: "upstream"},
	}
	for i := range rows {
		classify(&rows[i])
	}
	cs := clusterRows(rows, ClassOver)
	if len(cs) != 1 {
		t.Fatalf("want 1 OVER cluster (digit-normalized leading tokens), got %d: %+v", len(cs), cs)
	}
	if cs[0].Count != 2 {
		t.Errorf("cluster count = %d, want 2", cs[0].Count)
	}
	if want := leadingTokens(rows[0].SQL, 4); cs[0].Key != want {
		t.Errorf("OVER cluster key = %q, want leading tokens %q", cs[0].Key, want)
	}
}

// TestBuildRowsSkipPassesThroughUntouched: skip entries become Class=SKIP rows
// without omni evaluation and without classification (classify would
// overwrite Class).
func TestBuildRowsSkipPassesThroughUntouched(t *testing.T) {
	entries := []CorpusEntry{
		{SourcePath: "corpus/x_test.go", Line: 5, TestName: "TestX", SkipReason: "non_literal"},
	}
	rows, stats := buildRows(entries)
	if stats.duplicatesDropped != 0 {
		t.Errorf("duplicates dropped = %d, want 0", stats.duplicatesDropped)
	}
	if len(rows) != 1 {
		t.Fatalf("rows = %d, want 1", len(rows))
	}
	r := rows[0]
	if r.Class != ClassSkip {
		t.Errorf("class = %q, want %q", r.Class, ClassSkip)
	}
	if r.OmniVerdict != VerdictNone || r.OmniError != "" {
		t.Errorf("skip row must not be omni-evaluated: verdict=%q err=%q", r.OmniVerdict, r.OmniError)
	}
	if r.ClassifierReason != "" || r.DivergenceKey != "" {
		t.Errorf("skip row must not be classified: %+v", r)
	}
	if r.SkipReason != "non_literal" {
		t.Errorf("skip reason = %q, want preserved", r.SkipReason)
	}
}

// TestBuildRowsDedupByStmtHash: non-skip duplicates are dropped and counted;
// the first occurrence keeps its provenance; skip rows are never deduped.
func TestBuildRowsDedupByStmtHash(t *testing.T) {
	entries := []CorpusEntry{
		{SQL: "SELECT 1", Expected: VerdictAccept, SourcePath: "a_test.go", Line: 1, TestName: "TA"},
		{SQL: "SELECT 2", Expected: VerdictAccept, SourcePath: "a_test.go", Line: 2, TestName: "TA"},
		{SQL: "SELECT 1", Expected: VerdictAccept, SourcePath: "z_test.go", Line: 99, TestName: "TZ"},
		// Two skip entries share stmtHash("") — both must survive.
		{SkipReason: "non_literal", SourcePath: "a_test.go", Line: 3},
		{SkipReason: "non_composite_element", SourcePath: "a_test.go", Line: 4},
	}
	rows, stats := buildRows(entries)
	if stats.duplicatesDropped != 1 {
		t.Errorf("duplicates dropped = %d, want 1", stats.duplicatesDropped)
	}
	if stats.duplicateLabelConflicts != 0 {
		t.Errorf("label conflicts = %d, want 0 (labels agree)", stats.duplicateLabelConflicts)
	}
	if len(rows) != 4 {
		t.Fatalf("rows = %d, want 4 (2 distinct + 2 skip)", len(rows))
	}
	var sel1 *Row
	var skips int
	for i := range rows {
		if rows[i].SQL == "SELECT 1" {
			sel1 = &rows[i]
		}
		if rows[i].Class == ClassSkip {
			skips++
		}
	}
	if sel1 == nil {
		t.Fatal("SELECT 1 row missing")
	}
	if sel1.SourcePath != "a_test.go" || sel1.Line != 1 || sel1.TestName != "TA" {
		t.Errorf("kept row must carry first-occurrence provenance, got %s:%d %s", sel1.SourcePath, sel1.Line, sel1.TestName)
	}
	if skips != 2 {
		t.Errorf("skip rows = %d, want 2 (skips are never deduped)", skips)
	}
}

// TestBuildRowsDuplicateLabelConflict: a dropped duplicate whose upstream
// label conflicts with the kept row's flips the kept row to INDETERMINATE —
// the corpus context (SQL mode, runner flags) evidently matters for that
// statement, so neither label is trustworthy ground truth.
func TestBuildRowsDuplicateLabelConflict(t *testing.T) {
	entries := []CorpusEntry{
		{SQL: "SELECT 1", Expected: VerdictAccept, SourcePath: "a_test.go", Line: 1},
		{SQL: "SELECT 1", Expected: VerdictReject, SourcePath: "b_test.go", Line: 2}, // conflicting dup
		{SQL: "SELECT 2", Expected: VerdictAccept, SourcePath: "a_test.go", Line: 3},
		{SQL: "SELECT 2", Expected: VerdictAccept, SourcePath: "b_test.go", Line: 4}, // agreeing dup
	}
	rows, stats := buildRows(entries)
	if stats.duplicatesDropped != 2 {
		t.Errorf("duplicatesDropped = %d, want 2", stats.duplicatesDropped)
	}
	if stats.duplicateLabelConflicts != 1 {
		t.Errorf("duplicateLabelConflicts = %d, want 1", stats.duplicateLabelConflicts)
	}
	if len(rows) != 2 {
		t.Fatalf("rows = %d, want 2", len(rows))
	}
	bySQL := map[string]Row{}
	for _, r := range rows {
		bySQL[r.SQL] = r
	}
	c := bySQL["SELECT 1"]
	if c.Class != ClassIndeterminate || c.ClassifierReason != "duplicate_label_conflict" {
		t.Errorf("conflicted row: class %q reason %q, want INDETERMINATE/duplicate_label_conflict", c.Class, c.ClassifierReason)
	}
	if c.DivergenceKey != "" {
		t.Errorf("conflicted row divergence_key = %q, want cleared (INDETERMINATE rows are not clustered)", c.DivergenceKey)
	}
	agree := bySQL["SELECT 2"]
	if agree.Class != ClassAgreeAccept {
		t.Errorf("agreeing-dup row class = %q, want AGREE_ACCEPT", agree.Class)
	}
}
