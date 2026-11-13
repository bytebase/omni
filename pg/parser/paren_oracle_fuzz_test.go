package parser

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest/compat"
)

// The paren fuzz generates balanced-paren FROM-clause statements from a
// handful of template families (nested parens, SELECT/VALUES/WITH/TABLE
// subqueries, set operations, joins, LATERAL, alias lists, obvious rejects)
// and holds omni's routing of `(` to PostgreSQL on each. It complements the
// hand-written shape tests (paren_oracle_from_*_test.go) with depths and
// compositions no one would enumerate.
//
// The statements are rows of compat.tsv labelled parenFuzzLabel, judged once
// in a recording like any other row. So accept against reject is
// TestCompatibility's to hold: a statement only one side accepts is a
// disagreeing row, and a change that adds one fails by name.
// TestParenFuzzRows holds the rows to the generator, parenFuzzSize statements
// from a fixed seed, and holds the one thing a verdict does not say: a
// statement both accept routes its first FROM item, when that opens with a
// paren, to a subquery or a joined table.
const (
	fuzzSeed       int64 = 0xBADC0DE1
	parenFuzzSize        = 10000
	parenFuzzLabel       = "suite:paren-fuzz"
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
// depth 1–4. Depth 0 is excluded because it reduces to bare `FROM T`
// (no `(` dispatch); including it would inflate the denominator of any
// mismatch-rate statistic with cases that don't exercise the surface
// this fuzz corpus is built to stress.
func (g *parenFuzzGenerator) genNestedParens() string {
	// depth ∈ [1,4]
	depth := 1 + g.rng.Intn(4)
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

// genReservedMisuse mixes in reserved words where a relation is
// expected — another shape of obvious-reject.
//
// Trailing-SELECT patterns like `FROM (<subquery>) SELECT 1` were
// removed (see PAREN_KNOWN_BUGS.md PAREN-KB-2): they measure omni's
// top-level statement-list handling (Parse() silently truncating to
// the first RawStmt), not the `(` dispatch question this corpus is
// scoped to. Counting them as mismatches here was scope leakage.
func (g *parenFuzzGenerator) genReservedMisuse() string {
	switch g.rng.Intn(2) {
	case 0:
		return "SELECT * FROM ( )" // empty parens
	default:
		// bare keyword instead of a relation
		return "SELECT * FROM (WHERE)"
	}
}

// parenFuzzStatements returns the sweep's distinct statements, in the order
// the generator first makes them.
func parenFuzzStatements() []string {
	gen := newParenFuzzGenerator(fuzzSeed)
	seen := map[string]bool{}
	var out []string
	for i := 0; i < parenFuzzSize; i++ {
		sql, _ := gen.generateOne()
		if seen[sql] {
			continue
		}
		seen[sql] = true
		out = append(out, sql)
	}
	return out
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

// TestParenFuzzRows fails when compat.tsv's parenFuzzLabel rows are not the
// sweep's statements. -update-compat rewrites them, keeping the verdicts of
// statements it had; a new one comes in unjudged, for a recording of
// TestCompatibility to judge.
//
// It then checks the routing of every row PostgreSQL accepts and omni
// parses. The generator only makes FROM-clause shapes, so a first FROM item
// that opens with a paren and is neither a subquery nor a joined table
// means the dispatch fell through. The LATERAL family writes FROM T,
// LATERAL (...), whose first item is a plain relation; those are left out.
func TestParenFuzzRows(t *testing.T) {
	path := filepath.Join("testdata", "compat.tsv")
	f, err := compat.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	want := parenFuzzStatements()

	var kept []compat.Row
	had := map[string]compat.Row{}
	var got []string
	for _, r := range f.Rows {
		if r.Rule == parenFuzzLabel {
			had[r.Statement] = r
			got = append(got, r.Statement)
			continue
		}
		kept = append(kept, r)
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		if !compat.Updating() {
			t.Fatalf("%s holds %d %s rows, but the sweep makes %d other ones; rerun with -update-compat, then judge the new rows:\n  OMNI_RECORD=1 go test ./pg/parser/ -run '^TestCompatibility$'",
				path, len(got), parenFuzzLabel, len(want))
		}
		for _, sql := range want {
			r, ok := had[sql]
			if !ok {
				r = compat.Row{Omni: compat.Unjudged, Expected: compat.Unjudged, Source: "-", Rule: parenFuzzLabel, Statement: sql}
			}
			kept = append(kept, r)
		}
		f.Rows = kept
		if err := f.Write(path); err != nil {
			t.Fatal(err)
		}
		t.Logf("rewrote the %d %s rows of %s", len(want), parenFuzzLabel, path)
		return
	}

	routed := 0
	for _, r := range f.Rows {
		if r.Rule != parenFuzzLabel || r.Expected != compat.Accept || !firstFromIsParen(r.Statement) {
			continue
		}
		switch status, _ := classifyOmni(r.Statement); status {
		case OmniSubquery, OmniJoinedTable:
			routed++
		case OmniOther:
			t.Errorf("omni routes the first FROM item to neither a subquery nor a joined table:\n    %s", r.Statement)
		}
	}
	t.Logf("%d %s rows; %d that PostgreSQL accepts open FROM with a paren and route to a subquery or a joined table", len(want), parenFuzzLabel, routed)
}
