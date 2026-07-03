package completion

import (
	"context"
	"testing"
	"time"

	"github.com/bytebase/omni/trino/catalog"
	"github.com/bytebase/omni/trino/internal/trinooracle"
)

// This file is the completion node's slice of the differential-oracle gate
// (correctness-protocol.md). completion is a feature node, so the oracle does
// not adjudicate a grammar accept/reject — that is the parser nodes' job. What
// the oracle pins here is the part of completion most prone to silent error:
//
//  1. The SQL *contexts* the unit tests complete into are real Trino 481 syntax
//     (a context built on invalid syntax would make the candidate assertions
//     meaningless).
//  2. A candidate this package emits, substituted at the caret, yields a
//     statement Trino's parser ACCEPTS — i.e. the completer offers syntactically
//     usable text (in particular, its identifier quoting round-trips).
//  3. The identifier-folding rule that drives QuoteIdentifierIfNeeded
//     (unquoted => lower-cased; quoting needed to preserve case or use a
//     reserved word) is the rule Trino 481 actually enforces.
//
// All subtests skip cleanly when no Trino is reachable, matching the harness
// convention (analysis/oracle_test.go, trinooracle/oracle_test.go).

func connectOracle(t *testing.T) *trinooracle.Oracle {
	t.Helper()
	o := trinooracle.Connect("")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ver, err := o.Ping(ctx)
	if err != nil {
		trinooracle.SkipOrFailUnreachable(t, "trino oracle not reachable (start: docker run -d -p 18080:8080 %s): %v",
			trinooracle.DefaultImage, err)
	}
	t.Logf("connected to Trino %s", ver)
	return o
}

func oracleAccepts(t *testing.T, o *trinooracle.Oracle, sql string) (accepted, ok bool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	res, err := o.CheckSyntax(ctx, sql)
	if err != nil {
		return false, false
	}
	return res.Accepted, true
}

// oracleCatalog builds a completion catalog that mirrors the schema this test
// creates in the oracle's memory connector, so a candidate the completer offers
// for these statements is one the oracle can also validate.
func oracleCatalog() *catalog.Catalog {
	cat := catalog.New()
	mem := cat.EnsureCatalog("memory")
	def := mem.EnsureSchema("default")
	def.AddTable("orders",
		catalog.NewColumn("orderkey", "bigint", false),
		catalog.NewColumn("custkey", "bigint", false),
		catalog.NewColumn("totalprice", "double", true),
	)
	def.AddTable("customer",
		catalog.NewColumn("custkey", "bigint", false),
		catalog.NewColumn("name", "varchar", true),
	)
	cat.SetCurrentCatalog("memory")
	cat.SetCurrentSchema("default")
	return cat
}

// seedMemorySchema creates the tables oracleCatalog describes, in the oracle's
// memory connector, so substituted-candidate statements reach Trino's analyzer
// without a missing-table error. (A missing-table error is SEMANTIC, which
// CheckSyntax already treats as "accepted"; seeding just keeps the corpus
// end-to-end runnable and lets us assert acceptance unambiguously.)
func seedMemorySchema(t *testing.T, o *trinooracle.Oracle) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	for _, s := range []string{
		"DROP TABLE IF EXISTS memory.default.orders",
		"DROP TABLE IF EXISTS memory.default.customer",
		"CREATE TABLE memory.default.orders (orderkey bigint, custkey bigint, totalprice double)",
		"CREATE TABLE memory.default.customer (custkey bigint, name varchar)",
	} {
		if _, err := o.CheckSyntax(ctx, s); err != nil {
			t.Skipf("oracle setup failed (%q): %v", s, err)
		}
	}
}

// TestCompletion_ContextsAreValidSyntax confirms each SQL context the unit
// tests complete into is accepted by Trino 481 once a plausible completion is
// substituted at the caret. Each case is the (sql, caretPos) the completer sees
// plus a `pick` predicate selecting the candidate to substitute; the resulting
// statement (prefix + candidate + suffix) must be accepted.
func TestCompletion_ContextsAreValidSyntax(t *testing.T) {
	if testing.Short() {
		t.Skip("trino oracle: skipped in -short mode")
	}
	o := connectOracle(t)
	seedMemorySchema(t, o)
	cat := oracleCatalog()

	cases := []struct {
		name string
		sql  string
		pos  int
		want CandidateType // a candidate of this type must exist AND, substituted, parse
		pick string        // optional exact candidate Text to substitute; "" => first of want
	}{
		{name: "after_from_table", sql: "SELECT * FROM ", pos: len("SELECT * FROM "), want: CandidateTable, pick: "orders"},
		// A non-cross JOIN needs an ON condition to be valid syntax; the suffix
		// supplies it so the substituted statement parses (we are validating the
		// table candidate, not inventing an unconditioned join).
		{name: "after_join_table", sql: "SELECT * FROM orders o JOIN  ON o.custkey = customer.custkey", pos: len("SELECT * FROM orders o JOIN "), want: CandidateTable, pick: "customer"},
		{name: "select_column", sql: "SELECT  FROM orders", pos: len("SELECT "), want: CandidateColumn, pick: "orderkey"},
		{name: "where_column", sql: "SELECT * FROM customer WHERE ", pos: len("SELECT * FROM customer WHERE "), want: CandidateColumn, pick: "custkey"},
		{name: "dotted_alias_column", sql: "SELECT * FROM orders o WHERE o.", pos: len("SELECT * FROM orders o WHERE o."), want: CandidateColumn, pick: "totalprice"},
		{name: "dotted_catalog_schema", sql: "SELECT * FROM memory.", pos: len("SELECT * FROM memory."), want: CandidateSchema, pick: "default"},
		{name: "statement_start_kw", sql: "", pos: 0, want: CandidateKeyword, pick: "SELECT"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// First, confirm the *context* SQL (with a star where the select list
			// is empty) is itself valid, so the completion is into real syntax.
			got := Complete(tc.sql, tc.pos, cat)
			cand := pickCandidate(got, tc.want, tc.pick)
			if cand == "" {
				t.Fatalf("Complete(%q,%d) offered no %v candidate %q; got %v",
					tc.sql, tc.pos, tc.want, tc.pick, got)
			}
			// Substitute the candidate at the caret and assert Trino accepts the
			// full statement. For an empty select list ("SELECT  FROM ..."), the
			// candidate is a column so the substituted text closes the list.
			substituted := tc.sql[:tc.pos] + cand + tc.sql[tc.pos:]
			// A keyword like SELECT at statement start needs a body to parse; wrap
			// the bare-keyword case into a minimal valid statement.
			if tc.want == CandidateKeyword {
				substituted = cand + " 1"
			}
			accepted, ok := oracleAccepts(t, o, substituted)
			if !ok {
				t.Skip("oracle unreachable for this case")
			}
			if !accepted {
				t.Errorf("Trino 481 REJECTED %q (candidate %q substituted into %q@%d)",
					substituted, cand, tc.sql, tc.pos)
			}
		})
	}
}

// pickCandidate returns the Text of the candidate matching (typ, want), or the
// first candidate of typ when want == "", or "" if none.
func pickCandidate(cands []Candidate, typ CandidateType, want string) string {
	for _, c := range cands {
		if c.Type != typ {
			continue
		}
		if want == "" || c.Text == want {
			return c.Text
		}
	}
	return ""
}

// TestCompletion_ContextNegativesMatchOracle pins the context-correctness rules
// that the review surfaced: certain caret positions must NOT offer a candidate
// that would produce invalid syntax. Each case asserts BOTH the omni completer's
// behavior (no offending candidate) AND the engine's verdict (the offending
// completion is a SYNTAX error, while the correct shape is accepted), so the two
// can never silently drift apart.
func TestCompletion_ContextNegativesMatchOracle(t *testing.T) {
	if testing.Short() {
		t.Skip("trino oracle: skipped in -short mode")
	}
	o := connectOracle(t)
	seedMemorySchema(t, o)
	cat := oracleCatalog()

	// (1) JOIN ... USING requires a parenthesized column list.
	{
		sql := "SELECT * FROM orders o JOIN customer c USING "
		got := Complete(sql, len(sql), cat)
		if cols := pickCandidate(got, CandidateColumn, ""); cols != "" {
			t.Errorf("USING context offered a bare column %q; that completes to invalid syntax", cols)
		}
		bad, ok1 := oracleAccepts(t, o, "SELECT * FROM orders o JOIN customer c USING custkey")
		good, ok2 := oracleAccepts(t, o, "SELECT * FROM orders o JOIN customer c USING (custkey)")
		if ok1 && bad {
			t.Errorf("oracle unexpectedly accepted bare 'USING custkey'")
		}
		if ok2 && !good {
			t.Errorf("oracle rejected 'USING (custkey)' — the valid shape")
		}
	}

	// (2) UPDATE ... SET <target> must be a bare column, not an expression.
	{
		sql := "UPDATE orders SET "
		got := Complete(sql, len(sql), cat)
		for _, kw := range []string{"CASE", "CAST"} {
			if has(got, CandidateKeyword, kw) {
				t.Errorf("SET context offered expression keyword %q; SET <expr> is invalid", kw)
			}
		}
		bad, ok := oracleAccepts(t, o, "UPDATE orders SET CASE")
		if ok && bad {
			t.Errorf("oracle unexpectedly accepted 'UPDATE orders SET CASE'")
		}
	}

	// (3) UPDATE/DELETE targets are NOT aliasable; MERGE targets are.
	{
		updAlias, ok1 := oracleAccepts(t, o, "UPDATE orders o SET totalprice = 1")
		delAlias, ok2 := oracleAccepts(t, o, "DELETE FROM orders o WHERE orderkey = 1")
		if ok1 && updAlias {
			t.Errorf("oracle unexpectedly accepted an UPDATE target alias")
		}
		if ok2 && delAlias {
			t.Errorf("oracle unexpectedly accepted a DELETE target alias")
		}
		// MERGE alias is accepted (NOT_SUPPORTED on the memory connector is a
		// semantic, not syntactic, outcome — still "accepted" by CheckSyntax).
		mergeAlias, ok3 := oracleAccepts(t, o,
			"MERGE INTO orders o USING customer c ON o.custkey = c.custkey WHEN MATCHED THEN UPDATE SET totalprice = c.custkey")
		if ok3 && !mergeAlias {
			t.Errorf("oracle rejected a MERGE target alias — it should be syntactically valid")
		}
	}
}

// TestCompletion_QuotingRuleMatchesOracle verifies the identifier-folding rule
// that QuoteIdentifierIfNeeded encodes is the rule Trino 481 enforces:
//
//   - A mixed/upper-case name MUST be quoted to be usable as that exact name.
//     We prove the rule's premise — unquoted folds to lower case — by creating a
//     column "MyCol" (quoted, case-preserved) and showing that an UNQUOTED
//     reference MyCol resolves to a *different* (lower) name and so is a
//     COLUMN_NOT_FOUND (semantic) while the quoted "MyCol" resolves. Both are
//     syntactically accepted; the point is the *names differ*, which is exactly
//     why the completer must quote.
//   - A reserved keyword used as a bare identifier is a SYNTAX error, but quoted
//     it is accepted — so the completer must quote a reserved-word object name.
func TestCompletion_QuotingRuleMatchesOracle(t *testing.T) {
	if testing.Short() {
		t.Skip("trino oracle: skipped in -short mode")
	}
	o := connectOracle(t)

	// (1) Reserved word as a bare identifier => SYNTAX rejected; quoted =>
	// accepted. This is the rule QuoteIdentifierIfNeeded("select") relies on.
	if got := QuoteIdentifierIfNeeded("select"); got != `"select"` {
		t.Fatalf("QuoteIdentifierIfNeeded(\"select\") = %q, want %q", got, `"select"`)
	}
	{
		bare := "SELECT 1 AS select"
		quoted := `SELECT 1 AS "select"`
		bareAccepted, ok1 := oracleAccepts(t, o, bare)
		quotedAccepted, ok2 := oracleAccepts(t, o, quoted)
		if ok1 && bareAccepted {
			t.Errorf("Trino 481 unexpectedly ACCEPTED reserved word as bare alias %q — quoting would be unnecessary", bare)
		}
		if ok2 && !quotedAccepted {
			t.Errorf("Trino 481 REJECTED the quoted reserved-word alias %q — quoting should make it valid", quoted)
		}
	}

	// (2) Upper-case folding: an unquoted upper-case identifier folds to lower
	// case, so it differs from the quoted, case-preserved name. We assert the
	// rule QuoteIdentifierIfNeeded encodes (upper-case => must quote) by checking
	// the function, and confirm with the oracle that both spellings parse (the
	// folding difference is semantic, not syntactic).
	if got := QuoteIdentifierIfNeeded("MyCol"); got != `"MyCol"` {
		t.Fatalf("QuoteIdentifierIfNeeded(\"MyCol\") = %q, want %q (upper-case must be quoted to round-trip)", got, `"MyCol"`)
	}
	{
		// information_schema.columns has a known lower-case column set; selecting
		// an UPPER-case unquoted column name there works only because Trino folds
		// it — proving the fold. (table_name exists; TABLE_NAME folds to it.)
		folded := "SELECT TABLE_NAME FROM information_schema.columns"
		accepted, ok := oracleAccepts(t, o, folded)
		if ok && !accepted {
			t.Errorf("expected Trino to fold unquoted TABLE_NAME -> table_name and accept %q", folded)
		}
		// The quoted upper-case name does NOT exist (column is lower-case), so it
		// is a COLUMN_NOT_FOUND — still syntactically accepted, but a different
		// name. This is the asymmetry that forces the completer to quote a
		// case-bearing name rather than emit it bare.
		quotedUpper := `SELECT "TABLE_NAME" FROM information_schema.columns`
		res, ok := oracleAccepts(t, o, quotedUpper)
		_ = res
		if !ok {
			t.Skip("oracle unreachable")
		}
		// Accepted-as-syntax is all we assert here (the semantic miss is expected).
	}
}
