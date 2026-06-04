package parser

import (
	"context"
	"testing"
	"time"
)

// This file is the parser-select node's correctness slice for the named-WINDOW
// clause (window.go) plus the node's divergence audit (the docs-ahead-of-legacy
// P1 extensions D1–D3 that omni rejects and Trino 481 accepts, and the one
// over-permissive table-function corner). The oracle differential is
// authoritative; structural tests pin the windowDefinition shape. Helpers live
// in oracle_foundation_test.go and select_test.go.

// windowAcceptCorpus is the querySpecification WINDOW-clause surface Trino 481
// accepts. (The inline OVER window specification is the expressions node's
// concern; here the named WINDOW clause defines names an OVER may reference.)
var windowAcceptCorpus = []string{
	"SELECT count(*) OVER w FROM t WINDOW w AS (PARTITION BY a ORDER BY b)",
	"SELECT 1 FROM t WINDOW w AS ()",
	"SELECT 1 FROM t WINDOW w AS (PARTITION BY a)",
	"SELECT 1 FROM t WINDOW w AS (ORDER BY a)",
	"SELECT 1 FROM t WINDOW w AS (PARTITION BY a ORDER BY b ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW)",
	"SELECT 1 FROM t WINDOW w AS (RANGE UNBOUNDED PRECEDING)",
	"SELECT 1 FROM t WINDOW w1 AS (PARTITION BY a), w2 AS (ORDER BY b)",
	// the legacy window_clause example: a window that references a prior window name
	"SELECT * FROM T WINDOW someWindow AS (PARTITION BY a), otherWindow AS (someWindow ORDER BY b)",
	// a bare WINDOW (non-reserved) with nothing after is the table alias `t AS window`
	"SELECT 1 FROM t WINDOW",
}

// windowRejectCorpus is malformed WINDOW-clause input Trino 481 rejects.
// (Note `SELECT 1 FROM t WINDOW` — a bare WINDOW with nothing after — is NOT a
// reject: WINDOW is non-reserved, so it reads as the table alias `t AS window`.
// It is in windowAcceptCorpus.)
var windowRejectCorpus = []string{
	"SELECT 1 FROM t WINDOW w",                   // WINDOW clause name with no AS (…)
	"SELECT 1 FROM t WINDOW w AS",                // AS with no (…)
	"SELECT 1 FROM t WINDOW w AS PARTITION BY a", // window spec without parens
	"SELECT 1 FROM t WINDOW w AS (PARTITION a)",  // PARTITION without BY
}

func TestWindow_AcceptCorpusParses(t *testing.T) {
	for _, sql := range windowAcceptCorpus {
		t.Run(truncateName(sql), func(t *testing.T) {
			if _, errs := Parse(sql); len(errs) != 0 {
				t.Errorf("Parse(%q) should accept, got: %v", sql, errs)
			}
		})
	}
}

func TestWindow_RejectCorpusRejected(t *testing.T) {
	for _, sql := range windowRejectCorpus {
		t.Run(truncateName(sql), func(t *testing.T) {
			if _, errs := Parse(sql); len(errs) == 0 {
				t.Errorf("Parse(%q) should reject, but accepted", sql)
			}
		})
	}
}

func TestWindow_OracleDifferential(t *testing.T) {
	if testing.Short() {
		t.Skip("trino oracle: skipped in -short mode")
	}
	o := connectOracle(t)
	check := func(t *testing.T, sql string) {
		_, errs := Parse(sql)
		omniAccepts := len(errs) == 0
		trinoAccepts, ok := oracleAccepts(t, o, sql)
		if !ok {
			t.Skip("oracle unreachable for this case")
		}
		if omniAccepts != trinoAccepts {
			t.Errorf("MISMATCH sql=%q: omni accepts=%v (errs=%v), Trino accepts=%v",
				sql, omniAccepts, errs, trinoAccepts)
		}
	}
	for _, sql := range windowAcceptCorpus {
		t.Run("accept/"+truncateName(sql), func(t *testing.T) { check(t, sql) })
	}
	for _, sql := range windowRejectCorpus {
		t.Run("reject/"+truncateName(sql), func(t *testing.T) { check(t, sql) })
	}
}

func TestWindow_StructureDefinitions(t *testing.T) {
	spec := querySpec(t, parseOneQuery(t,
		"SELECT 1 FROM t WINDOW w1 AS (PARTITION BY a), w2 AS (ORDER BY b ROWS UNBOUNDED PRECEDING)"))
	if len(spec.Windows) != 2 {
		t.Fatalf("window defs=%d, want 2", len(spec.Windows))
	}
	if spec.Windows[0].Name.Value != "w1" {
		t.Errorf("window 0 name=%q, want w1", spec.Windows[0].Name.Value)
	}
	if spec.Windows[0].Spec == nil || len(spec.Windows[0].Spec.PartitionBy) != 1 {
		t.Errorf("window 0 spec = %+v, want one PARTITION BY", spec.Windows[0].Spec)
	}
	if spec.Windows[1].Name.Value != "w2" {
		t.Errorf("window 1 name=%q, want w2", spec.Windows[1].Name.Value)
	}
	if spec.Windows[1].Spec == nil || len(spec.Windows[1].Spec.OrderBy) != 1 || spec.Windows[1].Spec.Frame == nil {
		t.Errorf("window 1 spec = %+v, want ORDER BY + frame", spec.Windows[1].Spec)
	}
}

func TestWindow_StructureExistingNameReference(t *testing.T) {
	// `otherWindow AS (someWindow ORDER BY b)` references the base window someWindow.
	spec := querySpec(t, parseOneQuery(t,
		"SELECT * FROM T WINDOW someWindow AS (PARTITION BY a), otherWindow AS (someWindow ORDER BY b)"))
	if len(spec.Windows) != 2 {
		t.Fatalf("window defs=%d, want 2", len(spec.Windows))
	}
	other := spec.Windows[1].Spec
	if other == nil || other.ExistingName == nil || other.ExistingName.Value != "someWindow" {
		t.Errorf("otherWindow ExistingName = %+v, want someWindow", other)
	}
}

// ---------------------------------------------------------------------------
// Divergence audit (review-gate.md concern 2): the node implements the LEGACY
// grammar scope and FLAGS the docs-ahead-of-legacy P1 forms it deliberately does
// not parse. These tests pin those decisions: omni REJECTS the form while the
// Trino 481 oracle ACCEPTS it. If a future node implements the extension, the
// corresponding subtest will start failing — a signal to move the form into the
// accept corpus and clear the divergence.
// ---------------------------------------------------------------------------

// divergenceFlaggedForms are the oracle-accepted forms omni intentionally
// rejects per the migration divergence ledger (analysis §5 target policy:
// implement the legacy grammar; defer docs-ahead-of-legacy items as P1).
var divergenceFlaggedForms = []struct {
	name string
	sql  string
	note string
}{
	{
		name: "JSON_TABLE_relation",
		sql:  "SELECT * FROM JSON_TABLE('[1]', 'lax $' COLUMNS (a integer))",
		note: "divergence #3 (D1): JSON_TABLE relationPrimary is 481 but absent from the legacy grammar; deferred P1",
	},
	{
		name: "WITH_SESSION_prefix",
		sql:  "WITH SESSION foo = 1 SELECT 1",
		note: "divergence #4 (D2): WITH SESSION query prefix is 481 but absent from the legacy grammar; deferred P1",
	},
	{
		name: "set_op_CORRESPONDING",
		sql:  "SELECT 1 UNION CORRESPONDING SELECT 1",
		note: "divergence #6 (D3): set-op CORRESPONDING is 481 but absent from the legacy grammar; deferred P1",
	},
}

// TestSelect_DivergenceFlaggedRejected verifies omni rejects each flagged
// docs-ahead-of-legacy form (the legacy-grammar-scope behavior). Oracle-free.
func TestSelect_DivergenceFlaggedRejected(t *testing.T) {
	for _, d := range divergenceFlaggedForms {
		t.Run(d.name, func(t *testing.T) {
			if _, errs := Parse(d.sql); len(errs) == 0 {
				t.Errorf("Parse(%q) accepted, but this is a deferred P1 form omni should reject (%s)", d.sql, d.note)
			}
		})
	}
}

// TestSelect_DivergenceFlaggedOracleAccepts documents that Trino 481 ACCEPTS each
// flagged form — confirming the divergence is real (omni rejects, 481 accepts).
// This is the evidence half of the divergence packet; skipped without an oracle.
func TestSelect_DivergenceFlaggedOracleAccepts(t *testing.T) {
	if testing.Short() {
		t.Skip("trino oracle: skipped in -short mode")
	}
	o := connectOracle(t)
	for _, d := range divergenceFlaggedForms {
		t.Run(d.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			res, err := o.CheckSyntax(ctx, d.sql)
			if err != nil {
				t.Skip("oracle unreachable for this case")
			}
			if !res.Accepted {
				t.Errorf("Trino 481 REJECTED %q (errorName=%s); the divergence assumes 481 accepts it — re-check the ledger",
					d.sql, res.ErrorName)
			}
		})
	}
}

// TestSelect_CopartitionOverPermissive documents the one over-permissive corner
// (a flagged divergence): omni ACCEPTS a table-function table argument that is
// directly followed by COPARTITION with no intervening PARTITION BY, which Trino
// 481 rejects with a SYNTAX_ERROR. The legacy ANTLR grammar
// (tableFunctionCall: arg* (COPARTITION …)?) accepts it for any argument kind;
// matching 481's narrow rejection would require tracking per-argument partition
// state for a semantically-degenerate form, so the node follows the legacy
// grammar and flags it. Oracle-free assertion of omni's (intentional) accept.
func TestSelect_CopartitionOverPermissive(t *testing.T) {
	const sql = "SELECT * FROM TABLE(f(TABLE(a), TABLE(b) COPARTITION (a, b)))"
	if _, errs := Parse(sql); len(errs) != 0 {
		t.Errorf("Parse(%q): expected the flagged over-permissive ACCEPT, got errors: %v", sql, errs)
	}
}
