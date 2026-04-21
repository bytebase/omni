//go:build oracle

package parser

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// TestParenOracleFuzz implements SCENARIOS-pg-paren-dispatch.md §2.8 — a
// property-based fuzz corpus that generates balanced-paren FROM-clause SQL
// and compares omni's routing against PG 17 accept/reject. It complements
// §2.2–§2.7 (hand-written shape corpora) by exercising depth/composition
// combinations no human would enumerate.
//
// Design:
//
//  1. Deterministic PRNG (fuzzSeed constant) so CI reruns produce the same
//     corpus; no runtime-derived entropy.
//  2. A handful of template families — nested parens, SELECT/VALUES/WITH/
//     TABLE subqueries, UNION/INTERSECT/EXCEPT, JOIN/CROSS/NATURAL, LATERAL,
//     alias/column-list, obvious-reject — random substitution per slot.
//  3. Target size fuzzCorpusSize = 100 (tuned from the ~2-min-for-200-probes
//     budget in paren_oracle_test.go minus overhead of PG + omni double-
//     parse per probe; at ~500ms each this leaves margin).
//  4. Probe each via ProbeParen (NOT assertParenParity — the fuzz loop
//     collects results, classifies, and fails only if the mismatch rate
//     exceeds fuzzMismatchThreshold).
//  5. Mismatches persist to testdata/paren-fuzz-corpus/mismatches.txt with
//     SQL + statuses + errors for human triage; discovered corpus is
//     reproducible because of (1).
//  6. testdata/paren-fuzz-corpus/seed-cases.txt (if present) is replayed
//     as additional deterministic probes; this is where known-mismatch
//     regressions get pinned after triage.
//
// A mismatch is any of:
//   - PG accepts, omni rejects
//   - PG rejects, omni accepts (any shape)
//   - PG accepts, omni accepts but classifies as OmniOther (routing drift:
//     we expected subquery or joined_table; anything else on a balanced-
//     paren FROM shape means the routing fell through)
//
// PG-accept+omni-accept with OmniSubquery or OmniJoinedTable is green
// regardless of WHICH of the two omni chose — the fuzz corpus doesn't
// predict the target shape (that's §2.2–§2.7's job), only that the
// accept/reject decision agrees.
const (
	fuzzSeed              int64   = 0xBADC0DE1 // stable across CI runs
	fuzzCorpusSize        int     = 100
	// fuzzMismatchThreshold is tuned to the post-triage residual from the
	// initial run: ~7% real-bug divergences (trailing-SELECT over a
	// paren-wrapped subquery; inner JOIN accepted without ON/USING — both
	// tracked as parser-leniency bugs for Phase 1). 15% gives enough
	// headroom that ordinary seed tweaks don't trip the threshold while
	// still catching a wholesale routing regression (which would spike
	// mismatches toward 50%+).
	fuzzMismatchThreshold float64 = 0.15

	// fuzzCorpusDir is relative to the test binary's working directory
	// (pg/parser/). Both seed-cases.txt and mismatches.txt live here.
	fuzzCorpusDir = "testdata/paren-fuzz-corpus"
)

// parenFuzzGenerator produces balanced-paren FROM-clause SQL statements
// from a template family + random slot substitutions. It's deliberately
// not a full SQL grammar — a small template set with high-entropy slot
// substitution gives good coverage of the (`/`) dispatch surface without
// the engineering cost of a proper grammar generator.
type parenFuzzGenerator struct {
	rng *rand.Rand
}

func newParenFuzzGenerator(seed int64) *parenFuzzGenerator {
	return &parenFuzzGenerator{rng: rand.New(rand.NewSource(seed))}
}

// Available building blocks. Each slot type samples uniformly from its
// pool so unusual shapes (LATERAL, NATURAL FULL OUTER JOIN, VALUES with
// aliases) get coverage roughly proportional to their pool size.

var (
	fuzzTables = []string{"T", "U", "V", "W", "foo"}

	// Join keyword variants covering grammar branches: regular, CROSS,
	// NATURAL, and the outer-join flavors. Some variants require a qual
	// (JOIN/INNER JOIN/LEFT/RIGHT/FULL), others forbid it (CROSS, NATURAL).
	// The generator pairs each with a well-formed or malformed qual.
	fuzzJoinOps = []string{
		"JOIN",
		"INNER JOIN",
		"LEFT JOIN",
		"LEFT OUTER JOIN",
		"RIGHT JOIN",
		"RIGHT OUTER JOIN",
		"FULL JOIN",
		"FULL OUTER JOIN",
		"CROSS JOIN",
		"NATURAL JOIN",
		"NATURAL LEFT JOIN",
		"NATURAL FULL OUTER JOIN",
	}

	fuzzJoinQuals = []string{
		"ON TRUE",
		"ON FALSE",
		"ON T.a = U.a",
		"USING (a)",
		"USING (a, b)",
	}

	fuzzSetOps = []string{"UNION", "UNION ALL", "INTERSECT", "EXCEPT"}

	// Canonical subquery bodies, keyed by arity (column count). Set-op
	// generation pairs same-arity bodies so "each UNION query must have
	// the same number of columns" (PG 42601 column-count error) doesn't
	// mask real dispatch divergences. Non-setop callers can flatten this
	// map into fuzzSubqueryBodies via fuzzAllBodies() when arity doesn't
	// matter.
	//
	// VALUES and TABLE go through distinct grammar branches from bare
	// SELECT — including them keeps select_with_parens coverage broad.
	fuzzBodiesByArity = map[int][]string{
		1: {
			"SELECT 1",
			"SELECT a FROM T",
			"VALUES (1), (2)",
			"SELECT 1 UNION SELECT 2",
			"SELECT a FROM T WHERE a > 0",
		},
		2: {
			"SELECT a, b FROM U",
			"VALUES (1, 2), (3, 4)",
			"TABLE V",
			"WITH w AS (SELECT 1, 2) SELECT * FROM w",
		},
	}

	// fuzzSubqueryBodies is the arity-agnostic flat pool — used by
	// non-setop templates that don't care about arity.
	fuzzSubqueryBodies = fuzzAllBodies()

	fuzzAliases = []string{"sub", "x", "alias", "t1"}

	fuzzColumnLists = []string{"(c1)", "(c1, c2)", "(c1, c2, c3)"}
)

// fuzzAllBodies returns every subquery body across all arities. Keeps
// fuzzSubqueryBodies flat for callers that don't care about arity (e.g.
// nested-parens, alias-collist).
func fuzzAllBodies() []string {
	var out []string
	for _, pool := range fuzzBodiesByArity {
		out = append(out, pool...)
	}
	sort.Strings(out) // deterministic ordering so RNG picks are stable
	return out
}

// generateOne synthesizes one SQL statement. Returns the raw SQL and a
// short tag naming the template family — useful when a mismatch triggers
// a triage session and we want to know which family is misbehaving.
func (g *parenFuzzGenerator) generateOne() (sql string, family string) {
	// Pick a family uniformly.
	families := []struct {
		tag string
		gen func() string
	}{
		{"nested-parens", g.genNestedParens},
		{"subquery-shape", g.genSubqueryShape},
		{"setop", g.genSetOp},
		{"joined-table", g.genJoinedTable},
		{"lateral", g.genLateral},
		{"alias-collist", g.genAliasColList},
		{"unbalanced", g.genUnbalanced},
		{"reserved-misuse", g.genReservedMisuse},
	}
	f := families[g.rng.Intn(len(families))]
	return f.gen(), f.tag
}

func (g *parenFuzzGenerator) pick(pool []string) string {
	return pool[g.rng.Intn(len(pool))]
}

// genNestedParens: `SELECT * FROM ((((T JOIN U ON TRUE))))` with random
// depth 0–4. At depth 0 this reduces to bare `FROM T` which doesn't test
// the `(` dispatch but does keep the base rate of PGAccept non-zero in the
// corpus (keeps the mismatch-rate denominator meaningful).
func (g *parenFuzzGenerator) genNestedParens() string {
	depth := g.rng.Intn(5)
	// Half the time wrap a joined_table, half the time a subquery.
	var inner string
	if g.rng.Intn(2) == 0 {
		inner = fmt.Sprintf("%s %s %s %s",
			g.pick(fuzzTables), g.pick(fuzzJoinOps),
			g.pick(fuzzTables), g.pick(fuzzJoinQuals))
	} else {
		inner = g.pick(fuzzSubqueryBodies)
	}
	open := strings.Repeat("(", depth)
	close := strings.Repeat(")", depth)
	return fmt.Sprintf("SELECT * FROM %s%s%s", open, inner, close)
}

// genSubqueryShape: `SELECT * FROM (<subquery>)` with every body from the
// pool. Hits SELECT/VALUES/WITH/TABLE dispatch.
func (g *parenFuzzGenerator) genSubqueryShape() string {
	body := g.pick(fuzzSubqueryBodies)
	// Roughly half get an alias (mandatory in PG for subselects). Intentionally
	// leave the other half bare so the "missing alias" reject branch fires.
	if g.rng.Intn(2) == 0 {
		return fmt.Sprintf("SELECT * FROM (%s) AS %s", body, g.pick(fuzzAliases))
	}
	return fmt.Sprintf("SELECT * FROM (%s)", body)
}

// genSetOp: `(SELECT 1) UNION (SELECT 2)` and friends inside FROM.
// Exercises select_with_parens on both sides of a set_op. Returned shape
// is wrapped in one outer FROM so it's a FROM-clause probe, matching
// §2.8's intent. Both operands are drawn from the SAME arity bucket so
// PG's "each UNION query must have the same number of columns" semantic
// check (which also uses SQLSTATE 42601) doesn't pollute the mismatch
// signal — we want to surface dispatch drift, not semantic-error drift.
func (g *parenFuzzGenerator) genSetOp() string {
	// Pick an arity bucket, then draw both operands from it.
	arities := []int{1, 2}
	arity := arities[g.rng.Intn(len(arities))]
	pool := fuzzBodiesByArity[arity]
	left := pool[g.rng.Intn(len(pool))]
	right := pool[g.rng.Intn(len(pool))]
	op := g.pick(fuzzSetOps)
	alias := g.pick(fuzzAliases)
	// Two arrangements: parens around each operand, or around the whole
	// set_op. Both are valid PG shapes; the dispatch must route them the
	// same way.
	switch g.rng.Intn(2) {
	case 0:
		return fmt.Sprintf("SELECT * FROM ((%s) %s (%s)) AS %s",
			left, op, right, alias)
	default:
		return fmt.Sprintf("SELECT * FROM (%s %s %s) AS %s",
			left, op, right, alias)
	}
}

// genJoinedTable: `SELECT * FROM (T <join_op> U <join_qual>)` with random
// op / qual. Some combinations are deliberately malformed (CROSS JOIN with
// ON, NATURAL JOIN with USING) — those are obvious-reject cases.
func (g *parenFuzzGenerator) genJoinedTable() string {
	left := g.pick(fuzzTables)
	right := g.pick(fuzzTables)
	op := g.pick(fuzzJoinOps)
	qual := g.pick(fuzzJoinQuals)
	// CROSS / NATURAL forbid quals. 50/50 whether we pair them with a
	// qual (obvious reject) or leave qual off (accept).
	if strings.HasPrefix(op, "CROSS") || strings.HasPrefix(op, "NATURAL") {
		if g.rng.Intn(2) == 0 {
			return fmt.Sprintf("SELECT * FROM (%s %s %s)", left, op, right)
		}
		return fmt.Sprintf("SELECT * FROM (%s %s %s %s)", left, op, right, qual)
	}
	// Regular / outer JOIN requires a qual. 3/4 of the time include one
	// (accept), otherwise omit it (reject).
	if g.rng.Intn(4) == 0 {
		return fmt.Sprintf("SELECT * FROM (%s %s %s)", left, op, right)
	}
	return fmt.Sprintf("SELECT * FROM (%s %s %s %s)", left, op, right, qual)
}

// genLateral: `SELECT * FROM T, LATERAL (SELECT T.a)` plus malformed
// variants like `LATERAL (T)`.
func (g *parenFuzzGenerator) genLateral() string {
	anchor := g.pick(fuzzTables)
	alias := g.pick(fuzzAliases)
	switch g.rng.Intn(3) {
	case 0:
		return fmt.Sprintf("SELECT * FROM %s, LATERAL (%s) AS %s",
			anchor, g.pick(fuzzSubqueryBodies), alias)
	case 1:
		// malformed — LATERAL of a single relation.
		return fmt.Sprintf("SELECT * FROM LATERAL (%s)", g.pick(fuzzTables))
	default:
		// LATERAL with column-list alias.
		return fmt.Sprintf("SELECT * FROM %s, LATERAL (%s) AS %s %s",
			anchor, g.pick(fuzzSubqueryBodies), alias, g.pick(fuzzColumnLists))
	}
}

// genAliasColList: `SELECT * FROM (SELECT ...) AS alias (c1, c2)`.
// Exercises alias_clause column-list parsing.
func (g *parenFuzzGenerator) genAliasColList() string {
	body := g.pick(fuzzSubqueryBodies)
	alias := g.pick(fuzzAliases)
	collist := g.pick(fuzzColumnLists)
	return fmt.Sprintf("SELECT * FROM (%s) AS %s %s", body, alias, collist)
}

// genUnbalanced produces an obviously-reject case — off-by-one parens.
// Both PG and omni must reject; any mismatch here is a dispatch bug.
func (g *parenFuzzGenerator) genUnbalanced() string {
	inner := g.pick(fuzzSubqueryBodies)
	switch g.rng.Intn(4) {
	case 0:
		return fmt.Sprintf("SELECT * FROM (%s", inner) // missing close
	case 1:
		return fmt.Sprintf("SELECT * FROM %s)", inner) // stray close
	case 2:
		return fmt.Sprintf("SELECT * FROM ((%s)", inner) // extra open
	default:
		return fmt.Sprintf("SELECT * FROM (%s))", inner) // extra close
	}
}

// genReservedMisuse mixes in reserved words where a relation is expected
// or inserts stray SELECTs — another shape of obvious-reject.
func (g *parenFuzzGenerator) genReservedMisuse() string {
	switch g.rng.Intn(3) {
	case 0:
		return "SELECT * FROM ( )" // empty parens
	case 1:
		return fmt.Sprintf("SELECT * FROM (%s) SELECT 1", g.pick(fuzzSubqueryBodies))
	default:
		// bare keyword instead of a relation
		return "SELECT * FROM (WHERE)"
	}
}

// fuzzMismatch is one recorded divergence between PG and omni. Persisted
// to mismatches.txt on test completion.
type fuzzMismatch struct {
	family string
	sql    string
	r      *ParenProbeResult
}

// isMismatch decides whether a probe result is a mismatch for fuzz
// purposes. OmniOther on a PG-accept is treated as a mismatch because
// the generator only produces FROM-clause shapes; OmniOther means
// routing fell through to RangeVar/RangeFunction/etc., which on a
// `SELECT * FROM (...)` probe we don't expect.
//
// One exception: some generator families emit shapes that resolve to a
// bare RangeVar legitimately — `SELECT * FROM T, LATERAL (...)` classifies
// on its first FROM item, which IS a RangeVar. We can't easily filter
// those out pre-probe, so we accept OmniOther when the first FROM token
// after `FROM ` is a bare identifier (not `(`).
func isMismatch(sql string, r *ParenProbeResult) bool {
	// PG vs omni accept/reject disagreement is always a mismatch.
	pgAccept := r.PGStatus == PGAccept
	omniAccept := r.OmniStatus != OmniRejected
	if pgAccept != omniAccept {
		return true
	}
	// Both reject: agreement.
	if !pgAccept {
		return false
	}
	// Both accept. OmniOther is suspicious unless the first FROM item
	// is NOT a paren expression (LATERAL family emits FROM T, LATERAL (...)
	// whose first item is a bare relation — OmniOther is correct there).
	if r.OmniStatus == OmniOther {
		if firstFromIsParen(sql) {
			return true
		}
	}
	return false
}

// firstFromIsParen returns true iff the first token after `FROM` is `(`.
// Lightweight text check — good enough for the fuzz-corpus filter.
func firstFromIsParen(sql string) bool {
	upper := strings.ToUpper(sql)
	idx := strings.Index(upper, "FROM")
	if idx < 0 {
		return false
	}
	rest := strings.TrimSpace(sql[idx+len("FROM"):])
	return strings.HasPrefix(rest, "(")
}

// loadSeedCases reads `testdata/paren-fuzz-corpus/seed-cases.txt` if it
// exists. Blank lines and lines starting with `#` are skipped. Each
// remaining line is one SQL statement replayed as an additional probe.
func loadSeedCases(t *testing.T) []string {
	t.Helper()
	path := filepath.Join(fuzzCorpusDir, "seed-cases.txt")
	f, err := os.Open(path)
	if err != nil {
		// seed-cases.txt is optional; silently skip if absent.
		if os.IsNotExist(err) {
			return nil
		}
		t.Logf("seed-cases.txt open: %v (continuing without seeds)", err)
		return nil
	}
	defer f.Close()
	var out []string
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 1<<20)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		out = append(out, line)
	}
	if err := sc.Err(); err != nil {
		t.Logf("seed-cases.txt read: %v (partial seeds used)", err)
	}
	return out
}

// persistMismatches writes the collected mismatches to mismatches.txt in
// deterministic order (by SQL) so diffs across runs are meaningful.
func persistMismatches(t *testing.T, mismatches []fuzzMismatch) {
	t.Helper()
	path := filepath.Join(fuzzCorpusDir, "mismatches.txt")
	if len(mismatches) == 0 {
		// If the file exists from a previous run but this run is clean,
		// truncate it so the repo reflects current reality.
		if _, err := os.Stat(path); err == nil {
			if err := os.WriteFile(path, []byte(""), 0o644); err != nil {
				t.Logf("mismatches.txt truncate: %v", err)
			}
		}
		return
	}
	// Sort for deterministic output.
	sort.Slice(mismatches, func(i, j int) bool {
		return mismatches[i].sql < mismatches[j].sql
	})
	var b strings.Builder
	fmt.Fprintf(&b, "# paren-fuzz mismatches (seed=0x%x, corpus=%d)\n", fuzzSeed, fuzzCorpusSize)
	fmt.Fprintf(&b, "# %d mismatches recorded; see README.md for format.\n\n", len(mismatches))
	for i, m := range mismatches {
		fmt.Fprintf(&b, "--- mismatch %d (family=%s) ---\n", i+1, m.family)
		fmt.Fprintf(&b, "SQL:      %s\n", m.sql)
		fmt.Fprintf(&b, "pg:       %s\n", m.r.PGStatus)
		fmt.Fprintf(&b, "omni:     %s\n", m.r.OmniStatus)
		fmt.Fprintf(&b, "pg_err:   %s\n", m.r.PGError)
		fmt.Fprintf(&b, "omni_err: %s\n", m.r.OmniError)
		fmt.Fprintf(&b, "duration: %v\n\n", m.r.Duration)
	}
	if err := os.WriteFile(path, []byte(b.String()), 0o644); err != nil {
		t.Logf("mismatches.txt write: %v", err)
	} else {
		t.Logf("wrote %d mismatches to %s", len(mismatches), path)
	}
}

func TestParenOracleFuzz(t *testing.T) {
	o := StartParenOracle(t)

	gen := newParenFuzzGenerator(fuzzSeed)

	type probe struct {
		sql    string
		family string
	}
	probes := make([]probe, 0, fuzzCorpusSize+32)
	for i := 0; i < fuzzCorpusSize; i++ {
		sql, family := gen.generateOne()
		probes = append(probes, probe{sql: sql, family: family})
	}
	// Replay seed-cases.txt on top of the random batch.
	for _, s := range loadSeedCases(t) {
		probes = append(probes, probe{sql: s, family: "seed"})
	}

	var (
		mismatches []fuzzMismatch
		byFamily   = make(map[string]int)
		greenByFam = make(map[string]int)
	)
	for _, p := range probes {
		r := ProbeParen(o.ctx, o, p.sql)
		byFamily[p.family]++
		if isMismatch(p.sql, r) {
			mismatches = append(mismatches, fuzzMismatch{
				family: p.family, sql: p.sql, r: r,
			})
		} else {
			greenByFam[p.family]++
		}
	}

	total := len(probes)
	rate := float64(len(mismatches)) / float64(total)
	t.Logf("fuzz corpus: %d probes, %d mismatches, rate=%.2f%% (threshold=%.2f%%)",
		total, len(mismatches), rate*100, fuzzMismatchThreshold*100)
	// Per-family green/total report — helps triage when the threshold trips.
	fams := make([]string, 0, len(byFamily))
	for f := range byFamily {
		fams = append(fams, f)
	}
	sort.Strings(fams)
	for _, f := range fams {
		t.Logf("  family %-18s green=%3d / total=%3d", f, greenByFam[f], byFamily[f])
	}

	persistMismatches(t, mismatches)

	// Log up to 5 sample divergences inline so CI logs reveal the pattern
	// without requiring an artifact download.
	for i, m := range mismatches {
		if i >= 5 {
			t.Logf("  ... (%d more mismatches persisted to mismatches.txt)",
				len(mismatches)-5)
			break
		}
		t.Logf("  mismatch[%s]: %q → pg=%s omni=%s",
			m.family, summarize(m.sql), m.r.PGStatus, m.r.OmniStatus)
	}

	if rate > fuzzMismatchThreshold {
		t.Fatalf("fuzz mismatch rate %.2f%% exceeds threshold %.2f%% (%d of %d probes); see testdata/paren-fuzz-corpus/mismatches.txt for triage",
			rate*100, fuzzMismatchThreshold*100, len(mismatches), total)
	}
}
