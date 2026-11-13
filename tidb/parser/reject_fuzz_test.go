package parser

import (
	"math/rand"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest/compat"
	"github.com/bytebase/omni/enginetest/grammar"
)

// The reject fuzz mutates statements TiDB accepts, at the token and byte
// level, and holds omni to TiDB on each mutant. It reaches malformed strings
// the grammar's near-misses never will. The mutants are rows of compat.tsv
// labelled rejectFuzzLabel, judged once in a recording like any other row, so
// an over-acceptance (omni accepts, TiDB rejects) is a disagreeing row and a
// change that adds one fails TestCompatibility by name.
//
// TestRejectFuzzRows holds those rows to the generator: rejectFuzzSize
// mutants from a fixed seed over a fixed pool. The pool is the statements TiDB
// accepts among the grammar's alternatives and the corpora (rejectFuzzPool),
// read from the expected column, so closing a gap in omni does not move it.
const (
	rejectFuzzSeed  int64 = 0x7E1DB_EEF
	rejectFuzzSize        = 300
	rejectFuzzLabel       = "suite:reject-fuzz-sweep"
)

// rejectFuzzSuites are the corpora whose accepted statements seed the sweep,
// beside the grammar's alternatives.
var rejectFuzzSuites = map[string]bool{
	"suite:pingcap-parser":       true,
	"suite:pingcap-parser-rest":  true,
	"suite:tidb-integrationtest": true,
	"suite:bytebase":             true,
}

// rejectFuzzPool returns the distinct statements TiDB accepts among the
// grammar's alternatives (not its near-misses) and rejectFuzzSuites, sorted. A
// mutant of a statement TiDB accepts is the input where "omni accepts, TiDB
// rejects" is a real over-acceptance rather than noise about an invalid seed.
func rejectFuzzPool(rows []compat.Row) []string {
	seen := map[string]bool{}
	var pool []string
	for _, r := range rows {
		if r.Expected != compat.Accept || r.Statement == "" || seen[r.Statement] {
			continue
		}
		generated := !strings.Contains(r.Rule, ":") && !grammar.IsNearMiss(r.Rule)
		if !generated && !rejectFuzzSuites[r.Rule] {
			continue
		}
		seen[r.Statement] = true
		pool = append(pool, r.Statement)
	}
	sort.Strings(pool)
	return pool
}

// rejectFuzzMutants returns the sweep's distinct mutants, in the order the
// generator first makes them.
func rejectFuzzMutants(pool []string) []string {
	rng := rand.New(rand.NewSource(rejectFuzzSeed))
	seen := map[string]bool{}
	var out []string
	for i := 0; i < rejectFuzzSize; i++ {
		seed := pool[rng.Intn(len(pool))]
		mutant, _ := mutate(rng, seed)
		if mutant == "" || mutant == seed || seen[mutant] {
			continue
		}
		seen[mutant] = true
		out = append(out, mutant)
	}
	return out
}

// TestRejectFuzzRows fails when compat.tsv's rejectFuzzLabel rows are not the
// sweep's mutants. -update-compat rewrites them, keeping the verdicts of
// mutants it had; a new mutant comes in unjudged, for a recording of
// TestCompatibility to judge.
func TestRejectFuzzRows(t *testing.T) {
	path := filepath.Join("testdata", "compat.tsv")
	f, err := compat.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	pool := rejectFuzzPool(f.Rows)
	if len(pool) == 0 {
		t.Fatalf("%s holds no statement TiDB accepts to seed the sweep", path)
	}
	want := rejectFuzzMutants(pool)

	var kept []compat.Row
	had := map[string]compat.Row{}
	var got []string
	for _, r := range f.Rows {
		if r.Rule == rejectFuzzLabel {
			had[r.Statement] = r
			got = append(got, r.Statement)
			continue
		}
		kept = append(kept, r)
	}
	if strings.Join(got, "\n") == strings.Join(want, "\n") {
		return
	}
	if !compat.Updating() {
		t.Fatalf("%s holds %d %s rows, but the sweep makes %d other ones; rerun with -update-compat, then judge the new rows:\n  OMNI_RECORD=1 go test ./tidb/parser/ -run '^TestCompatibility$'",
			path, len(got), rejectFuzzLabel, len(want))
	}
	for _, m := range want {
		r, ok := had[m]
		if !ok {
			r = compat.Row{Omni: compat.Unjudged, Expected: compat.Unjudged, Source: "-", Rule: rejectFuzzLabel, Statement: m}
		}
		kept = append(kept, r)
	}
	f.Rows = kept
	if err := f.Write(path); err != nil {
		t.Fatal(err)
	}
	t.Logf("rewrote the %d %s rows of %s", len(want), rejectFuzzLabel, path)
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
