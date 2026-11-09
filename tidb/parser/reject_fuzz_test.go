package parser

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// TestTiDBRejectFuzz is the differential mutation half of the §3.6 rejection
// suite. Where the grammar-guided negatives corpus enumerates one-step-off
// near-misses of each grammar alternative, this sweep mutates real valid SQL
// — the §3.3 completeness corpus and the §3.4 real-SQL samples — at the token
// and byte level, runs omni and live TiDB v8.5.5 on each mutant, and flags
// every over-acceptance (omni accepts, TiDB rejects). It reaches malformed
// strings enumeration never will.
//
// It mirrors pg/parser's TestParenOracleFuzz: a deterministic seed so a CI
// rerun produces the same corpus, a size dialled by TIDB_REJECT_FUZZ_N, a
// known-mismatches allowlist enforced by set-equality, and a TIDB_REJECT_FUZZ_DEFER
// mode that downgrades the gate to log-and-persist for the larger CI batch.
//
// The gating direction is over-acceptance. A mutant omni rejects and TiDB
// accepts is a §3.3 completeness gap, out of scope here; it is logged, not
// gated. A confirmed over-acceptance is minimised to the shortest disagreeing
// input and pinned into known-mismatches.txt (and, once fixed, promoted into
// parser_y_negatives_corpus.tsv as a permanent regression row).
const (
	rejectFuzzSeed        int64 = 0x7E1DB_EEF
	rejectFuzzSizeDefault       = 300

	// TIDB_REJECT_FUZZ_N overrides rejectFuzzSizeDefault so CI can dial the
	// per-PR batch up without recompiling. Unset or <=0 falls back to the default.
	envRejectFuzzN = "TIDB_REJECT_FUZZ_N"

	// TIDB_REJECT_FUZZ_DEFER switches the strict known-mismatches gate into a
	// non-blocking log-and-persist mode for the CI fuzz job, so a fresh finding
	// at large N does not auto-fail the build before a human triages it.
	envRejectFuzzDefer = "TIDB_REJECT_FUZZ_DEFER"

	rejectFuzzDir = "testdata/reject-fuzz-corpus"

	// rejectFuzzMinimizeCap bounds the live-TiDB probes one minimisation spends,
	// so a pathological seed cannot stall the sweep.
	rejectFuzzMinimizeCap = 400
)

func rejectFuzzSize() int {
	if v := strings.TrimSpace(os.Getenv(envRejectFuzzN)); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return rejectFuzzSizeDefault
}

func rejectFuzzDeferMode() bool {
	return strings.TrimSpace(os.Getenv(envRejectFuzzDefer)) != ""
}

// loadFuzzSeeds returns the distinct valid (accept/accept) statements of the
// §3.3 corpus and the §3.4 samples, the pool mutations draw from. A mutant of a
// statement both sides accept is the only input where "omni accepts, TiDB
// rejects" is a real over-acceptance rather than noise about an already-invalid
// seed.
func loadFuzzSeeds(t testing.TB) []string {
	t.Helper()
	paths := []string{
		filepath.Join("testdata", "parser_y_corpus.tsv"),
		filepath.Join("testdata", "pingcap_parser_corpus.tsv"),
		filepath.Join("testdata", "pingcap_parser_rest_corpus.tsv"),
		filepath.Join("testdata", "tidb_integrationtest_corpus.tsv"),
		filepath.Join("testdata", "bytebase_tidb_corpus.tsv"),
	}
	seen := map[string]bool{}
	var seeds []string
	for _, p := range paths {
		b, err := os.ReadFile(p)
		if err != nil {
			t.Fatalf("read %s: %v", p, err)
		}
		for _, raw := range strings.Split(string(b), "\n") {
			if raw == "" || strings.HasPrefix(raw, "#") {
				continue
			}
			fields := strings.Split(raw, "\t")
			if len(fields) < 3 {
				continue
			}
			// omni verdict, engine verdict, [label], statement.
			if fields[0] != "accept" || fields[1] != "accept" {
				continue
			}
			sql := fields[len(fields)-1]
			if sql == "" || seen[sql] {
				continue
			}
			seen[sql] = true
			seeds = append(seeds, sql)
		}
	}
	sort.Strings(seeds) // deterministic order so RNG picks are stable
	if len(seeds) == 0 {
		t.Fatal("no valid seeds loaded")
	}
	return seeds
}

// tokenSpans returns the byte spans of the real tokens in sql (those with a
// positive-width source range), in order. Mutations operate on these spans so a
// dropped or duplicated token is a lexically clean edit.
func tokenSpans(sql string) [][2]int {
	var spans [][2]int
	for _, tok := range Tokenize(sql) {
		if tok.Loc >= 0 && tok.End > tok.Loc && tok.End <= len(sql) {
			spans = append(spans, [2]int{tok.Loc, tok.End})
		}
	}
	return spans
}

// renderTokens joins the given token texts with single spaces. Mutations build a
// token-text slice and render it, so spacing never carries hidden structure.
func renderTokens(parts []string) string {
	return strings.Join(parts, " ")
}

// mutate applies one random token- or byte-level mutation to a valid seed and
// returns the mutant and a tag naming the mutation family.
func mutate(rng *rand.Rand, sql string) (string, string) {
	// 60% token mutations (structural near-misses), 40% byte mutations (lexical).
	if rng.Intn(5) < 3 {
		return mutateToken(rng, sql)
	}
	return mutateByte(rng, sql)
}

var fuzzStrayTokens = []string{"zzz", ",", "(", ")", "1", "'s'", "SELECT", "FROM", "NULL"}

func mutateToken(rng *rand.Rand, sql string) (string, string) {
	spans := tokenSpans(sql)
	if len(spans) < 2 {
		return mutateByte(rng, sql)
	}
	parts := make([]string, len(spans))
	for i, s := range spans {
		parts[i] = sql[s[0]:s[1]]
	}
	switch rng.Intn(5) {
	case 0: // drop a token
		i := rng.Intn(len(parts))
		return renderTokens(append(append([]string{}, parts[:i]...), parts[i+1:]...)), "drop"
	case 1: // duplicate a token
		i := rng.Intn(len(parts))
		out := append([]string{}, parts[:i+1]...)
		out = append(out, parts[i])
		out = append(out, parts[i+1:]...)
		return renderTokens(out), "dup"
	case 2: // swap two adjacent tokens
		i := rng.Intn(len(parts) - 1)
		out := append([]string{}, parts...)
		out[i], out[i+1] = out[i+1], out[i]
		return renderTokens(out), "swap"
	case 3: // insert a stray token
		i := rng.Intn(len(parts) + 1)
		stray := fuzzStrayTokens[rng.Intn(len(fuzzStrayTokens))]
		out := append([]string{}, parts[:i]...)
		out = append(out, stray)
		out = append(out, parts[i:]...)
		return renderTokens(out), "insert"
	default: // truncate at a token boundary
		i := 1 + rng.Intn(len(parts)-1)
		return renderTokens(parts[:i]), "truncate"
	}
}

func mutateByte(rng *rand.Rand, sql string) (string, string) {
	b := []byte(sql)
	if len(b) == 0 {
		return sql, "byte-noop"
	}
	switch rng.Intn(4) {
	case 0: // delete a byte
		i := rng.Intn(len(b))
		return string(append(append([]byte{}, b[:i]...), b[i+1:]...)), "byte-del"
	case 1: // duplicate a byte
		i := rng.Intn(len(b))
		return string(append(append(append([]byte{}, b[:i+1]...), b[i]), b[i+1:]...)), "byte-dup"
	case 2: // replace a byte with a punctuation byte
		i := rng.Intn(len(b))
		repl := []byte{'(', ')', ',', ';', '.', '\'', '`', '*', '='}[rng.Intn(9)]
		out := append([]byte{}, b...)
		out[i] = repl
		return string(out), "byte-repl"
	default: // truncate at a random byte
		i := 1 + rng.Intn(len(b))
		return string(b[:i]), "byte-truncate"
	}
}

// overAccepts reports whether omni accepts sql while TiDB rejects it: the
// over-acceptance the sweep gates on. A missing TiDB verdict (err) is not an
// over-acceptance.
func overAccepts(t *testing.T, tc *tidbContainer, sql, prev string) (bool, error) {
	if _, err := Parse(sql); err != nil {
		return false, nil // omni rejects: not an over-acceptance
	}
	engineAccepts, _, err := tc.syntaxVerdict(t, sql, prev)
	if err != nil {
		return false, err
	}
	return !engineAccepts, nil
}

// minimizeOverAcceptance shrinks an over-accepting mutant to a shorter one that
// still over-accepts, by repeatedly dropping each token while the property
// holds. It is a token-level ddmin bounded by rejectFuzzMinimizeCap live probes.
func minimizeOverAcceptance(t *testing.T, tc *tidbContainer, sql string) string {
	budget := rejectFuzzMinimizeCap
	changed := true
	for changed && budget > 0 {
		changed = false
		spans := tokenSpans(sql)
		if len(spans) < 2 {
			break
		}
		parts := make([]string, len(spans))
		for i, s := range spans {
			parts[i] = sql[s[0]:s[1]]
		}
		for i := range parts {
			if budget <= 0 {
				break
			}
			budget--
			cand := renderTokens(append(append([]string{}, parts[:i]...), parts[i+1:]...))
			if cand == "" || cand == sql {
				continue
			}
			over, err := overAccepts(t, tc, cand, sql)
			if err != nil {
				continue
			}
			if over {
				sql = cand
				changed = true
				break
			}
		}
	}
	return sql
}

// A fuzzFinding is one minimised over-acceptance.
type fuzzFinding struct {
	sql    string
	family string
}

func (f fuzzFinding) key() string { return f.sql }

func loadRejectFuzzKnown(t *testing.T) map[string]string {
	t.Helper()
	path := filepath.Join(rejectFuzzDir, "known-mismatches.txt")
	out := map[string]string{}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return out
		}
		t.Fatalf("open %s: %v", path, err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 1<<20)
	for sc.Scan() {
		line := sc.Text()
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}
		// family<TAB>sql ; sql is the key (it can contain spaces but no tab).
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			t.Fatalf("%s: malformed line (need family<TAB>sql): %q", path, line)
		}
		out[parts[1]] = parts[0]
	}
	return out
}

func writeRejectFuzzFindings(t *testing.T, findings []fuzzFinding) {
	t.Helper()
	if err := os.MkdirAll(rejectFuzzDir, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", rejectFuzzDir, err)
	}
	path := filepath.Join(rejectFuzzDir, "mismatches.txt")
	sort.Slice(findings, func(i, j int) bool { return findings[i].sql < findings[j].sql })
	var b strings.Builder
	fmt.Fprintf(&b, "# TiDB reject-fuzz over-acceptances (seed=0x%x, N=%d).\n", rejectFuzzSeed, rejectFuzzSize())
	fmt.Fprintf(&b, "# omni accepts, TiDB v8.5.5 rejects. Minimised. %d finding(s).\n", len(findings))
	fmt.Fprintf(&b, "# Triage: fix across the forks and pin into parser_y_negatives_corpus.tsv,\n")
	fmt.Fprintf(&b, "# or, if intended, move the family<TAB>sql line into known-mismatches.txt.\n\n")
	for _, f := range findings {
		fmt.Fprintf(&b, "%s\t%s\n", f.family, f.sql)
	}
	if err := os.WriteFile(path, []byte(b.String()), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
	t.Logf("wrote %d finding(s) to %s", len(findings), path)
}

func TestTiDBRejectFuzz(t *testing.T) {
	if want := "pingcap/tidb:" + tidbCorpusVersion; tidbImage != want {
		t.Fatalf("reject-fuzz is recorded against %s, but this package pins %s", want, tidbImage)
	}
	tc := startTiDB(t)
	assertTiDBVersion(t, tc)

	seeds := loadFuzzSeeds(t)
	rng := rand.New(rand.NewSource(rejectFuzzSeed))
	size := rejectFuzzSize()

	var (
		findings   []fuzzFinding
		byFamily   = map[string]int{}
		gaps       int // omni rejects, TiDB accepts (§3.3, logged not gated)
		noVerdict  int
		unanswered []string // mutants TiDB never gave a verdict for
		seenSQL    = map[string]bool{}
	)
	for i := 0; i < size; i++ {
		seed := seeds[rng.Intn(len(seeds))]
		mutant, family := mutate(rng, seed)
		if mutant == "" || mutant == seed {
			continue
		}
		byFamily[family]++

		omniAccepts := func() bool { _, err := Parse(mutant); return err == nil }()
		engineAccepts, _, err := tc.syntaxVerdict(t, mutant, seed)
		if err != nil {
			noVerdict++
			unanswered = append(unanswered, fmt.Sprintf("%q: %v", mutant, err))
			if tc.lost() {
				t.Errorf("no server left after %d mutants", i)
				break
			}
			if len(unanswered) > 10 {
				t.Errorf("gave up after %d mutants with no verdict", len(unanswered))
				break
			}
			continue
		}
		switch {
		case omniAccepts && !engineAccepts:
			min := minimizeOverAcceptance(t, tc, mutant)
			if seenSQL[min] {
				continue
			}
			seenSQL[min] = true
			findings = append(findings, fuzzFinding{sql: min, family: family})
		case !omniAccepts && engineAccepts:
			gaps++
		}
	}

	t.Logf("reject-fuzz: %d mutants, %d over-acceptances, %d gaps, %d no-verdict", size, len(findings), gaps, noVerdict)
	fams := make([]string, 0, len(byFamily))
	for f := range byFamily {
		fams = append(fams, f)
	}
	sort.Strings(fams)
	for _, f := range fams {
		t.Logf("  family %-12s %d", f, byFamily[f])
	}
	writeRejectFuzzFindings(t, findings)
	for i, f := range findings {
		if i >= 10 {
			t.Logf("  ... and %d more (see mismatches.txt)", len(findings)-10)
			break
		}
		t.Logf("  over-accept[%s]: %q", f.family, f.sql)
	}

	// An unmeasured mutant fails the run independently of the mismatch gate and
	// of defer mode: a sweep that could not get a verdict for every mutant is
	// not a clean run, and must not read as success just because the findings it
	// did measure match the baseline. This mirrors TestTiDBNegativesCorpusEngineVerdicts.
	for _, u := range unanswered {
		t.Errorf("reject-fuzz: no verdict: %s", u)
	}

	// Known-mismatches gate: strict set-equality, like pg's paren fuzz.
	known := loadRejectFuzzKnown(t)
	discovered := map[string]bool{}
	for _, f := range findings {
		discovered[f.key()] = true
	}
	var unexpected, disappeared []string
	for k := range discovered {
		if _, ok := known[k]; !ok {
			unexpected = append(unexpected, k)
		}
	}
	for k := range known {
		if !discovered[k] {
			disappeared = append(disappeared, k)
		}
	}
	sort.Strings(unexpected)
	sort.Strings(disappeared)

	if rejectFuzzDeferMode() {
		if len(unexpected) > 0 {
			t.Logf("reject-fuzz-defer: %d new over-acceptance(s) for triage (see mismatches.txt)", len(unexpected))
		}
		return
	}
	if len(unexpected) > 0 {
		for _, k := range unexpected {
			t.Errorf("reject-fuzz: NEW over-acceptance not in known-mismatches.txt: %q", k)
		}
	}
	if len(disappeared) > 0 {
		for _, k := range disappeared {
			t.Errorf("reject-fuzz: known over-acceptance no longer reproduces (re-baseline known-mismatches.txt after the fix): %q", k)
		}
	}
	if len(unexpected) > 0 || len(disappeared) > 0 {
		t.Fatalf("reject-fuzz mismatch set differs from known-mismatches.txt")
	}
}
