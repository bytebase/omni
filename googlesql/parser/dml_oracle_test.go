//go:build googlesql_oracle

// Differential test for the `parser-dml` node against the live Cloud Spanner
// emulator oracle. Run with:
//
//	SPANNER_EMULATOR_HOST=localhost:9010 \
//	  go test -tags googlesql_oracle ./googlesql/parser/ -run TestDMLDifferential
//
// It is a PROVE gate for googlesql/parser-dml (correctness-protocol.md): for
// every fixture it (1) feeds the statement to the emulator via the
// harness/googlesql-spanner CLI and reads the accept/reject verdict, then (2)
// parses the same statement with omni's Parse, and asserts the two agree — BOTH
// polarities (the corpus has accept AND reject cases). A harness `error` verdict
// (oracle could not decide) fails the fixture loudly; it is never folded into
// accept or reject. The harness plumbing (newSelectHarness / verdict / close) is
// shared with select_oracle_test.go (same package).
//
// SCOPE OF THIS DIFFERENTIAL — Spanner is authoritative ONLY for forms its
// grammar actually adjudicates. Every fixture below is one of:
//   - SHARED GoogleSQL DML (INSERT/UPDATE/DELETE incl. VALUES/SELECT/DEFAULT/
//     FROM-update/WITH OFFSET/THEN RETURN/ASSERT_ROWS_MODIFIED/ON CONFLICT) —
//     the Spanner verdict is authoritative (the emulator parses then
//     feature-rejects with "Table not found", which the harness classifies
//     ACCEPT per its message-prefix rule).
//   - Spanner-only upserts (INSERT OR IGNORE/UPDATE/REPLACE, bare IGNORE/REPLACE/
//     UPDATE) — authoritative ACCEPT.
//   - genuine SYNTAX rejects (truncated statements, trailing-comma SET list,
//     ON CONFLICT after a bare query) — authoritative REJECT.
//
// EXCLUDED (NON-authoritative on Spanner — covered in dml_test.go + the
// divergence ledger, NOT here, because the emulator's verdict would manufacture
// a false divergence):
//   - MERGE — the emulator rejects it at statement dispatch ("Statement not
//     supported: MergeStmt") BEFORE parsing the body, so it accepts even a MERGE
//     with zero/garbage WHEN clauses. No authoritative MERGE-syntax signal.
//   - TRUNCATE TABLE — routed to Spanner's DDL path, which has no TRUNCATE and
//     syntax-rejects it; BigQuery-only DML (truth1 BigQuery DML-003).
//   - dashed table paths (`my-project.ds.tbl`) — Spanner syntax-rejects '-' in a
//     table name (divergence #85); BigQuery-valid, omni accepts.
//   - INSERT … TABLE <clause> … ON CONFLICT — the emulator rejects the
//     on-conflict after a TABLE source; the .g4 (union) allows it.

package parser

import (
	"os"
	"testing"
)

// dmlFixture is one DML statement fed to BOTH the oracle and omni's Parse.
// wantParse is the expected omni outcome and MUST equal the emulator's grammar
// verdict.
type dmlFixture struct {
	sql       string
	wantParse bool
}

var dmlFixtures = []dmlFixture{
	// ===================== INSERT — accept (shared) =====================
	{"INSERT INTO Singers (SingerId, FirstName) VALUES (1, 'Marc')", true},
	{"INSERT Singers (SingerId, FirstName) VALUES (1, 'Marc')", true}, // no INTO
	{"INSERT INTO Singers VALUES (1, 'Marc')", true},                  // no column list
	{"INSERT INTO t (a) VALUES (1), (2), (3)", true},                  // multi-row
	{"INSERT INTO t (a, b) VALUES (1, DEFAULT)", true},                // DEFAULT value
	{"INSERT INTO t (a) SELECT x FROM s", true},                       // insert select
	{"INSERT INTO s.t (a) VALUES (1)", true},                          // schema-qualified target
	{"INSERT INTO t (a) (SELECT 1)", true},                            // ( query ) source
	{"INSERT INTO t (a, b) (SELECT 5, 7) UNION ALL SELECT 8, 9", true},
	{"INSERT INTO t (a) ((SELECT 1) UNION ALL (SELECT 2))", true},
	{"INSERT INTO t VALUES ((1, 2), (SELECT 1))", true}, // struct + scalar subquery values
	// A parenthesized query with a query-level trailing ORDER BY / LIMIT is a
	// BARE query (alt-1), accepted WITHOUT ON CONFLICT.
	{"INSERT INTO t (a) (SELECT 1) LIMIT 5", true},
	{"INSERT INTO t (a) (SELECT 1) ORDER BY x", true},
	{"INSERT INTO t (a) (SELECT 1) ORDER BY x LIMIT 5", true},

	// ===================== INSERT — Spanner upserts (accept) ============
	{"INSERT OR IGNORE INTO t (a) VALUES (1)", true},
	{"INSERT OR UPDATE INTO t (a) VALUES (1)", true},
	{"INSERT OR REPLACE INTO t (a) VALUES (1)", true},
	{"INSERT IGNORE INTO t (a) VALUES (1)", true},  // bare IGNORE (no OR)
	{"INSERT REPLACE INTO t (a) VALUES (1)", true}, // bare REPLACE
	{"INSERT UPDATE INTO t (a) VALUES (1)", true},  // bare UPDATE

	// ===================== INSERT — trailers (accept) ===================
	{"INSERT INTO t (a) VALUES (1) THEN RETURN a", true},
	{"INSERT INTO t (a) VALUES (1) THEN RETURN WITH ACTION a", true},
	{"INSERT INTO t (a) VALUES (1) THEN RETURN WITH ACTION AS act *", true},
	{"INSERT INTO t (a) VALUES (1) THEN RETURN a AS x, b", true},
	{"INSERT INTO t (a) VALUES (1) THEN RETURN * EXCEPT (a)", true},
	{"INSERT INTO t (a) VALUES (1) ASSERT_ROWS_MODIFIED 1", true},
	{"INSERT INTO t (a) VALUES (1) ASSERT_ROWS_MODIFIED @p", true},
	{"INSERT INTO t (a) VALUES (1) ASSERT_ROWS_MODIFIED CAST(1 AS INT64)", true},
	{"INSERT INTO t (a) VALUES (1) THEN RETURN WITH(x AS 1, x)", true}, // inline WITH expr item, NOT WITH ACTION
	{"INSERT INTO t (graph) VALUES (1)", true},                         // non-reserved keyword column name

	// ===================== INSERT — ON CONFLICT (accept) ================
	{"INSERT INTO t (a) VALUES (1) ON CONFLICT DO NOTHING", true},
	{"INSERT INTO t (a) VALUES (1) ON CONFLICT (a) DO UPDATE SET a = 2", true},
	{"INSERT INTO t (a) VALUES (1) ON CONFLICT (a, b) DO NOTHING", true},
	{"INSERT INTO t (a) VALUES (1) ON CONFLICT ON UNIQUE CONSTRAINT uc DO NOTHING", true},
	{"INSERT INTO t (a) VALUES (1) ON CONFLICT (a) DO UPDATE SET a = 2 WHERE a > 1", true},
	{"INSERT INTO t (a) (SELECT 1) ON CONFLICT DO NOTHING", true},

	// ===================== UPDATE — accept (shared) =====================
	{"UPDATE Singers SET FirstName = 'A' WHERE SingerId = 1", true},
	{"UPDATE t a SET a.x = 1 WHERE a.id = 2", true},
	{"UPDATE t AS a SET x = 1, y = 2 WHERE id = 3", true},
	{"UPDATE t SET x = DEFAULT WHERE id = 1", true},
	{"UPDATE t a SET x = b.y FROM s b WHERE a.id = b.id", true}, // join-update
	{"UPDATE t SET x = 1", true},                                // no WHERE
	{"UPDATE t SET x = 1 THEN RETURN x", true},
	{"UPDATE t SET x = 1 WHERE id = 1 ASSERT_ROWS_MODIFIED 1", true},
	{"UPDATE t SET arr[OFFSET(0)] = 5 WHERE id = 1", true},              // generalized LHS
	{"UPDATE t SET (DELETE FROM t.arr WHERE x = 1) WHERE id = 1", true}, // nested DML item
	{"UPDATE t SET x = (SELECT c FROM s) WHERE id = 1", true},           // embedded subquery (valid)

	// ===================== DELETE — accept (shared) =====================
	{"DELETE FROM Singers WHERE SingerId = 1", true},
	{"DELETE Singers WHERE SingerId = 1", true}, // no FROM
	{"DELETE FROM t", true},                     // no WHERE
	{"DELETE Singers AS s WHERE s.SingerId = 1 THEN RETURN s.FirstName", true},
	{"DELETE FROM t WHERE id = 1 THEN RETURN *", true},
	{"DELETE FROM t WHERE id = 1 THEN RETURN * EXCEPT (a)", true},

	// ===================== negatives — reject (authoritative) ===========
	{"INSERT INTO", false},
	{"INSERT t", false},
	{"INSERT INTO t (a) VALUES", false},
	{"INSERT INTO t (a) VALUES (1) THEN RETURN", false},
	{"INSERT INTO t (a) SELECT 1 ON CONFLICT DO NOTHING", false},           // ON CONFLICT after bare query
	{"INSERT INTO t (a) (SELECT 1) LIMIT 5 ON CONFLICT DO NOTHING", false}, // ON CONFLICT after a bare (parenthesized + LIMIT) query
	{"UPDATE", false},
	{"UPDATE t SET", false},
	{"UPDATE t SET x = 1, WHERE id = 1", false}, // trailing comma in SET list
	{"UPDATE t SET x = 1 FROM", false},
	{"UPDATE t SET x = (SELECT 1 FROM s a b) WHERE id = 1", false}, // malformed embedded subquery (stray alias)
	{"DELETE", false},
	{"DELETE FROM t ASSERT_ROWS_MODIFIED 1 + 1", false},      // ASSERT operand is not an expression
	{"DELETE FROM t ASSERT_ROWS_MODIFIED 'x'", false},        // ASSERT operand is not a string
	{"DELETE FROM t ASSERT_ROWS_MODIFIED (SELECT 1)", false}, // ASSERT operand is not a subquery
}

func TestDMLDifferential(t *testing.T) {
	if os.Getenv("SPANNER_EMULATOR_HOST") == "" {
		t.Skip("SPANNER_EMULATOR_HOST not set; skipping live differential")
	}
	h := newSelectHarness(t)
	defer h.close()

	for _, fx := range dmlFixtures {
		fx := fx
		t.Run(fx.sql, func(t *testing.T) {
			// 1. Oracle verdict.
			v := h.verdict(t, fx.sql)
			switch v.Verdict {
			case "accept", "reject":
				// proceed
			case "error":
				t.Fatalf("oracle returned ERROR (no verdict) for %q: reason=%s code=%s msg=%s\n"+
					"the oracle could not decide; fix the wrapper/emulator — do NOT treat as accept/reject",
					fx.sql, v.Reason, v.Code, v.Message)
			default:
				t.Fatalf("unexpected harness verdict %q for %q", v.Verdict, fx.sql)
			}
			oracleAccepts := v.Verdict == "accept"

			// Sanity: the asserted wantParse must match the live oracle. A mismatch
			// here means the fixture is mis-tagged (or the form is actually
			// non-authoritative on Spanner and should be moved out of this set).
			if oracleAccepts != fx.wantParse {
				t.Fatalf("fixture wantParse=%v but oracle says %q (%s) for %q",
					fx.wantParse, v.Verdict, v.Message, fx.sql)
			}

			// 2. omni Parse verdict.
			_, errs := Parse(fx.sql)
			omniAccepts := len(errs) == 0

			if omniAccepts != oracleAccepts {
				t.Errorf("DIVERGENCE on %q: omni accepts=%v, oracle accepts=%v (%s: %s)",
					fx.sql, omniAccepts, oracleAccepts, v.Reason, v.Message)
			}
		})
	}
}
