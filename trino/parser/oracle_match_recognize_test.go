package parser

import "testing"

// This file is the parser-match-recognize node's differential-oracle gate
// (correctness-protocol.md): for every statement in the corpus below, omni's
// Parse accept/reject verdict must equal Trino 481's verdict as reported by the
// live oracle (trinooracle.CheckSyntax, which classifies a SYNTAX_ERROR as
// reject and every other outcome — success or a semantic error such as
// TABLE_HAS_NO_COLUMNS / COLUMN_NOT_FOUND / INVALID_LABEL / NOT_SUPPORTED — as
// accept).
//
// The corpus is NOT split into hardcoded accept/reject lists: the oracle is the
// source of truth, so each case is run through both omni and Trino and the two
// verdicts are compared. Every legacy ANTLR alternative of the patternRecognition
// subsystem (truth2) and every documented MATCH_RECOGNIZE form (truth1) appears,
// alongside the precedence/structure cases (R1–R5, M1–M6) and a set of negatives
// the engine rejects. Skipped cleanly when no Trino oracle is reachable.
//
// Helpers (connectOracle / oracleAccepts / truncateName) live in
// oracle_foundation_test.go.

// matchRecognizeOracleCorpus is the full MATCH_RECOGNIZE / row-pattern corpus.
// Each string is a complete statement (the differential needs a whole statement;
// table/column references resolve to semantic errors, never syntax rejections).
var matchRecognizeOracleCorpus = []string{
	// --- minimal & mandatory clauses (M1) ---
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A B) DEFINE A AS true, B AS false)",
	"SELECT * FROM t MATCH_RECOGNIZE ()",                                 // empty body — rejected
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A))",                      // no DEFINE — rejected
	"SELECT * FROM t MATCH_RECOGNIZE (DEFINE A AS true)",                 // no PATTERN — rejected
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN A DEFINE A AS true)",       // PATTERN body parens required
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN () DEFINE A AS true)",      // empty PATTERN body — rejected (rowPattern required; cf. PATTERN(()) below)
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A) DEFINE)",               // DEFINE needs >=1 def
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A) SUBSET DEFINE A AS x)", // SUBSET needs >=1 def

	// --- full clause stack, in grammar order ---
	"SELECT * FROM t MATCH_RECOGNIZE (PARTITION BY a ORDER BY b MEASURES c AS m ONE ROW PER MATCH AFTER MATCH SKIP TO NEXT ROW PATTERN (A B+ C) DEFINE B AS B.x > 0)",
	"SELECT * FROM t MATCH_RECOGNIZE (PARTITION BY a, b ORDER BY c DESC, d ASC NULLS LAST PATTERN (A) DEFINE A AS true)",

	// --- clause ordering is fixed (M2): reordering is rejected ---
	"SELECT * FROM t MATCH_RECOGNIZE (ORDER BY b PARTITION BY a PATTERN (A) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (MEASURES x AS m PARTITION BY a PATTERN (A) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A) MEASURES x AS m DEFINE A AS true)", // MEASURES after PATTERN — rejected
	"SELECT * FROM t MATCH_RECOGNIZE (DEFINE A AS true PATTERN (A))",                 // DEFINE before PATTERN — rejected

	// --- MEASURES (with RUNNING / FINAL semantics) ---
	"SELECT * FROM t MATCH_RECOGNIZE (MEASURES match_number() AS mn PATTERN (A) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (MEASURES RUNNING last(A.x) AS m PATTERN (A) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (MEASURES FINAL last(A.x) AS m PATTERN (A) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (MEASURES x AS m, y AS n PATTERN (A) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (MEASURES x m PATTERN (A) DEFINE A AS true)", // measure AS required — rejected

	// --- rowsPerMatch + emptyMatchHandling ---
	"SELECT * FROM t MATCH_RECOGNIZE (ONE ROW PER MATCH PATTERN (A) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (ALL ROWS PER MATCH PATTERN (A) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (ALL ROWS PER MATCH SHOW EMPTY MATCHES PATTERN (A) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (ALL ROWS PER MATCH OMIT EMPTY MATCHES PATTERN (A) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (ALL ROWS PER MATCH WITH UNMATCHED ROWS PATTERN (A) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (ONE ROW PER MATCH WITH UNMATCHED ROWS PATTERN (A) DEFINE A AS true)", // handling only on ALL ROWS — rejected
	"SELECT * FROM t MATCH_RECOGNIZE (ONE ROWS PER MATCH PATTERN (A) DEFINE A AS true)",                    // ONE ROWS — rejected
	"SELECT * FROM t MATCH_RECOGNIZE (ALL ROW PER MATCH PATTERN (A) DEFINE A AS true)",                     // ALL ROW — rejected

	// --- AFTER MATCH SKIP variants (M5) ---
	"SELECT * FROM t MATCH_RECOGNIZE (AFTER MATCH SKIP PAST LAST ROW PATTERN (A) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (AFTER MATCH SKIP TO NEXT ROW PATTERN (A) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (AFTER MATCH SKIP TO FIRST A PATTERN (A) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (AFTER MATCH SKIP TO LAST A PATTERN (A) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (AFTER MATCH SKIP TO A PATTERN (A) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (AFTER MATCH SKIP TO PATTERN (A) DEFINE A AS true)", // SKIP TO needs a target — rejected
	"SELECT * FROM t MATCH_RECOGNIZE (AFTER MATCH SKIP PATTERN (A) DEFINE A AS true)",    // SKIP needs TO/PAST — rejected
	// NEXT/FIRST/LAST are non-reserved: each may itself be the skip variable name
	// when not in its special form (NEXT not followed by ROW, FIRST/LAST not
	// followed by an identifier).
	"SELECT * FROM t MATCH_RECOGNIZE (AFTER MATCH SKIP TO next PATTERN (A) DEFINE A AS true, next AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (AFTER MATCH SKIP TO NEXT PATTERN (A) DEFINE A AS true)",            // NEXT as variable (no ROW)
	"SELECT * FROM t MATCH_RECOGNIZE (AFTER MATCH SKIP TO FIRST PATTERN (A) DEFINE A AS true)",           // FIRST as variable (no ident)
	"SELECT * FROM t MATCH_RECOGNIZE (AFTER MATCH SKIP TO LAST PATTERN (A) DEFINE A AS true)",            // LAST as variable (no ident)
	"SELECT * FROM t MATCH_RECOGNIZE (AFTER MATCH SKIP TO FIRST first PATTERN (A) DEFINE first AS true)", // FIRST keyword + variable "first"

	// --- INITIAL / SEEK (M5) ---
	"SELECT * FROM t MATCH_RECOGNIZE (INITIAL PATTERN (A) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (SEEK PATTERN (A) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (INITIAL SEEK PATTERN (A) DEFINE A AS true)", // only one of INITIAL/SEEK — rejected

	// --- SUBSET ---
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A B) SUBSET U = (A) DEFINE A AS true, B AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A B) SUBSET U = (A, B) DEFINE A AS true, B AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A B) SUBSET U = (A), V = (B) DEFINE A AS true, B AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A) SUBSET U = () DEFINE A AS true)", // empty union — rejected
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A) SUBSET U DEFINE A AS true)",      // SUBSET needs = (…) — rejected

	// --- trailing alias + column aliases (M3) ---
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A) DEFINE A AS true) AS m",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A) DEFINE A AS true) m",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A) DEFINE A AS true) AS m (x, y)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A) DEFINE A AS true) m (x, y)",
	"SELECT * FROM t AS r MATCH_RECOGNIZE (PATTERN (A) DEFINE A AS true) AS m", // double alias
	"SELECT * FROM t r MATCH_RECOGNIZE (PATTERN (A) DEFINE A AS true) m",       // double alias, no AS
	// the trailing alias may itself be the non-reserved keyword MATCH_RECOGNIZE,
	// even followed by a column list — no second clause can follow, so it is an alias.
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A) DEFINE A AS true) MATCH_RECOGNIZE (x)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A) DEFINE A AS true) AS MATCH_RECOGNIZE (x)",

	// --- MATCH_RECOGNIZE binds before TABLESAMPLE (M4) ---
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A) DEFINE A AS true) TABLESAMPLE BERNOULLI (10)",
	"SELECT * FROM t TABLESAMPLE BERNOULLI (10) MATCH_RECOGNIZE (PATTERN (A) DEFINE A AS true)", // wrong order — rejected

	// --- input-relation shapes the clause attaches to ---
	"SELECT * FROM (SELECT 1 x) t MATCH_RECOGNIZE (PATTERN (A) DEFINE A AS true)",
	"SELECT * FROM (t JOIN s ON true) MATCH_RECOGNIZE (PATTERN (A) DEFINE A AS true)",

	// --- rowPattern: concatenation / alternation precedence (R1) ---
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A B C) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A | B) DEFINE A AS true, B AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A | B | C) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A B | C D) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN ((A B) | C) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN ((A | B)+) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A |) DEFINE A AS true)", // dangling alternation — rejected
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (| A) DEFINE A AS true)", // leading alternation — rejected

	// --- rowPattern: quantifiers (R2) ---
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A*) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A+) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A?) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A*?) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A+?) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A??) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A{3}) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A{1,3}) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A{,3}) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A{2,}) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A{2,}?) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A{,}) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN ((A)*) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN ((A*)*) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A* B* C*) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A**) DEFINE A AS true)",     // adjacent quantifiers — rejected
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A{2}{3}) DEFINE A AS true)", // adjacent quantifiers — rejected
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A{}) DEFINE A AS true)",     // empty braces — rejected

	// --- rowPattern: anchors (R5) ---
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (^ A+ $) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (^) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN ($) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (^+) DEFINE A AS true)", // anchor + quantifier parses (semantic INVALID_LABEL)

	// --- rowPattern: PERMUTE ---
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (PERMUTE(A, B)) DEFINE A AS true, B AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (PERMUTE(A)) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (PERMUTE(A B, C)) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (PERMUTE(A | B, C)) DEFINE A AS true)",
	// empty PERMUTE() is ACCEPTED by the live Trino 481 server even though the
	// published SqlBase.g4 requires >=1 operand (divergence #92, oracle decides).
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (PERMUTE()) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (PERMUTE(,)) DEFINE A AS true)",  // leading comma — rejected
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (PERMUTE(A,)) DEFINE A AS true)", // trailing comma — rejected
	// PERMUTE is non-reserved: a bare `permute` not followed by '(' is a label.
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (permute) DEFINE permute AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (permute A) DEFINE permute AS true, A AS true)",

	// --- rowPattern: excluded pattern ---
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN ({- A -} B) DEFINE A AS true, B AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN ({- A B -} C) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN ({- A | B -}) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN ({- -}) DEFINE A AS true)", // empty exclusion — rejected

	// --- rowPattern: empty pattern (R4) ---
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (()) DEFINE A AS true)",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (() A) DEFINE A AS true)",

	// --- quoted / non-reserved pattern variable labels ---
	`SELECT * FROM t MATCH_RECOGNIZE (PATTERN ("quoted var") DEFINE "quoted var" AS true)`,
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (zone) DEFINE zone AS true)", // non-reserved keyword as a label

	// --- DEFINE condition expressions (PREV / NEXT / CLASSIFIER / FIRST / LAST) ---
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A) DEFINE A AS A.price > PREV(A.price))",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A) DEFINE A AS CLASSIFIER() = 'X')",
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A B) DEFINE B AS B.value > LAST(A.value))",

	// --- window-frame pattern recognition (M6) ---
	"SELECT count(*) OVER (PARTITION BY a ORDER BY b MEASURES x AS m ROWS BETWEEN CURRENT ROW AND UNBOUNDED FOLLOWING PATTERN (A) DEFINE A AS true) FROM t",
	"SELECT count(*) OVER (ORDER BY b ROWS BETWEEN CURRENT ROW AND UNBOUNDED FOLLOWING PATTERN (A) DEFINE A AS true) FROM t",
	"SELECT count(*) OVER (MEASURES x AS m ROWS CURRENT ROW PATTERN (A) DEFINE A AS true) FROM t",
	"SELECT count(*) OVER (ROWS CURRENT ROW AFTER MATCH SKIP PAST LAST ROW PATTERN (A) DEFINE A AS true) FROM t",
	"SELECT count(*) OVER (ROWS CURRENT ROW INITIAL PATTERN (A) DEFINE A AS true) FROM t",
	"SELECT count(*) OVER (ROWS CURRENT ROW PATTERN (A B) SUBSET U = (A, B) DEFINE A AS true) FROM t",
	// In a window frame the skipTo target may be followed by NOTHING (frame end),
	// or by SUBSET/DEFINE with no PATTERN — every clause after frameExtent is
	// optional (M6). The FIRST/LAST skip-keyword disambiguation must allow these.
	"SELECT count(*) OVER (ROWS CURRENT ROW AFTER MATCH SKIP TO FIRST A) FROM t",
	"SELECT count(*) OVER (ROWS CURRENT ROW AFTER MATCH SKIP TO FIRST A DEFINE A AS true) FROM t",
	"SELECT count(*) OVER (ROWS CURRENT ROW AFTER MATCH SKIP TO LAST A SUBSET U = (A) DEFINE A AS true) FROM t",
	"SELECT count(*) OVER (MEASURES x AS m ROWS CURRENT ROW) FROM t",                                                                       // MEASURES, no PATTERN/DEFINE — accepted (M6)
	"SELECT count(*) OVER (ROWS CURRENT ROW PATTERN (A)) FROM t",                                                                           // PATTERN, no DEFINE — accepted (M6)
	"SELECT count(*) OVER w FROM t WINDOW w AS (ORDER BY b ROWS BETWEEN CURRENT ROW AND UNBOUNDED FOLLOWING PATTERN (A) DEFINE A AS true)", // named window
	"SELECT count(*) OVER (PATTERN (A) DEFINE A AS true) FROM t",                                                                           // PATTERN with no frameExtent — rejected (frameExtent mandatory)
	"SELECT count(*) OVER (MEASURES x AS m PATTERN (A) DEFINE A AS true) FROM t",                                                           // MEASURES + PATTERN, no frameExtent — rejected
	"SELECT count(*) OVER (MEASURES x AS m) FROM t",                                                                                        // bare MEASURES, no frameExtent — rejected
	"SELECT count(*) OVER (ORDER BY a PATTERN (A) DEFINE A AS true) FROM t",                                                                // PATTERN after ORDER BY, no frameExtent — rejected
	// MEASURES vs existing-window-name disambiguation (the measure expr starts with
	// a non-reserved keyword that also begins a frame type — needs >1-token lookahead).
	"SELECT count(*) OVER (MEASURES rows AS m ROWS CURRENT ROW PATTERN (A) DEFINE A AS true) FROM t",  // MEASURES clause, measure col "rows"
	"SELECT count(*) OVER (MEASURES range AS m ROWS CURRENT ROW PATTERN (A) DEFINE A AS true) FROM t", // measure col "range"
	// a window named "measures" (MEASURES followed by a window-spec part, not a measure expr)
	"SELECT count(*) OVER (measures ORDER BY b ROWS CURRENT ROW) FROM t WINDOW measures AS (PARTITION BY a)",
	"SELECT count(*) OVER (MEASURES ROWS CURRENT ROW) FROM t WINDOW MEASURES AS (PARTITION BY a)", // MEASURES is the window name + ROWS frame

	// --- negatives that look almost right ---
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A) DEFINE A AS true) FROM",     // trailing FROM is not an alias — rejected
	"SELECT * FROM t MATCH_RECOGNIZE PATTERN (A) DEFINE A AS true",            // missing outer parens — rejected
	"SELECT * FROM t MATCH_RECOGNIZE (PATTERN (A) DEFINE A AS true",           // unterminated — rejected
	"SELECT * FROM t MATCH_RECOGNIZE (PARTITION a PATTERN (A) DEFINE A AS x)", // PARTITION without BY — rejected
}

// TestMatchRecognize_OracleDifferential is the authoritative accept/reject gate:
// omni's Parse verdict must equal Trino 481's verdict for every corpus statement.
func TestMatchRecognize_OracleDifferential(t *testing.T) {
	if testing.Short() {
		t.Skip("trino oracle: skipped in -short mode")
	}
	o := connectOracle(t)
	for _, sql := range matchRecognizeOracleCorpus {
		sql := sql
		t.Run(truncateName(sql), func(t *testing.T) {
			_, errs := Parse(sql)
			omniAccepts := len(errs) == 0

			trinoAccepts, ok := oracleAccepts(t, o, sql)
			if !ok {
				t.Skip("oracle unreachable for this case")
			}
			if omniAccepts != trinoAccepts {
				t.Errorf("MISMATCH %q: omni accepts=%v (errs=%v), Trino accepts=%v",
					sql, omniAccepts, errs, trinoAccepts)
			}
		})
	}
}
