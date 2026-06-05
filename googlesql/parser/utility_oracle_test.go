//go:build googlesql_oracle

// Differential test for the `parser-utility` node against the live Cloud Spanner
// emulator oracle. Run with:
//
//	SPANNER_EMULATOR_HOST=localhost:9010 \
//	  go test -tags googlesql_oracle ./googlesql/parser/ -run TestUtilityDifferential
//
// It is the PROVE gate for googlesql/parser-utility (correctness-protocol.md):
// for every fixture it (1) feeds the statement to the emulator via the
// harness/googlesql-spanner CLI and reads the accept/reject verdict, then (2)
// parses the same statement with omni's Parse, and asserts the two agree — BOTH
// polarities. A harness `error` verdict (oracle could not decide) fails loudly;
// it is never folded into accept or reject. The harness plumbing
// (newSelectHarness / verdict / close) is shared with select_oracle_test.go.
//
// SCOPE — what is AUTHORITATIVE on the Spanner emulator (oracle.md), and so is
// diffed here:
//   - the LEADING-FORM accept of BEGIN / START TRANSACTION / COMMIT / ROLLBACK /
//     START BATCH / RUN BATCH / ABORT BATCH / ASSERT / DESCRIBE — the emulator
//     PARSES these (accepts then feature-rejects "Statement not supported: …").
//   - CALL — parsed in FULL by the emulator: it validates the argument list, so
//     both its accepts and its precise reject cases (trailing comma, `1 => 2`)
//     are authoritative.
//   - RENAME TABLE and bare ANALYZE — go through the emulator's real DDL parser
//     (RENAME TABLE accepts; bare ANALYZE accepts).
//   - genuine shared SYNTAX rejects (CALL without parens, etc.).
//
// EXCLUDED (NON-authoritative on Spanner — covered in transaction_test.go /
// utility_test.go + the divergence ledger, NOT here, because the emulator's
// verdict would manufacture a false divergence):
//   - the PRECISE trailing grammar of BEGIN/COMMIT/ROLLBACK/BATCH/ASSERT/DESCRIBE
//     — the emulator's shallow recognizer SWALLOWS trailing tokens (it accepts
//     `COMMIT WORK`, `START BATCH a b`, `ASSERT 1 AS notstring`, `DESCRIBE @#$`),
//     while the .g4 (the grammar bytebase consumes) rejects them. omni follows
//     the .g4, so a Spanner accept there is a false divergence.
//   - ANALYZE with targets (`ANALYZE t`) — Spanner's ANALYZE is the bare
//     whole-database form and syntax-rejects targets; the .g4 + BigQuery union
//     accept them (omni accepts).
//   - RENAME with a non-TABLE object kind — Spanner narrows RENAME to TABLE; the
//     .g4 allows any object-kind identifier (omni accepts).

package parser

import (
	"os"
	"testing"
)

// utilityFixture is one statement fed to BOTH the oracle and omni's Parse.
// wantParse is the expected omni outcome and MUST equal the emulator's grammar
// verdict.
type utilityFixture struct {
	sql       string
	wantParse bool
}

var utilityFixtures = []utilityFixture{
	// ===================== Transactions — accept (leading form) =========
	{"BEGIN", true},
	{"BEGIN TRANSACTION", true},
	{"START TRANSACTION", true},
	{"COMMIT", true},
	{"COMMIT TRANSACTION", true},
	{"ROLLBACK", true},
	{"ROLLBACK TRANSACTION", true},
	{"BEGIN READ ONLY", true},
	{"BEGIN READ WRITE", true},
	{"BEGIN ISOLATION LEVEL SERIALIZABLE", true},
	{"BEGIN ISOLATION LEVEL REPEATABLE READ", true},
	{"BEGIN TRANSACTION READ ONLY", true},
	{"BEGIN TRANSACTION READ WRITE, ISOLATION LEVEL READ COMMITTED", true},
	{"START TRANSACTION READ ONLY, ISOLATION LEVEL SERIALIZABLE", true},

	// ===================== Batch — accept (leading form) ================
	{"START BATCH", true},
	{"START BATCH ddl", true},
	{"RUN BATCH", true},
	{"ABORT BATCH", true},

	// ===================== ASSERT — accept (leading form) ===============
	{"ASSERT true", true},
	{"ASSERT 1 = 1 AS 'must hold'", true},
	{"ASSERT EXISTS(SELECT 1 FROM t)", true},
	{"ASSERT (SELECT 1) = 1 AS \"simple test\"", true},
	{"ASSERT @param = 1 AS 'param test'", true},
	{"ASSERT 'x' IN ('x', 'y')", true},

	// ===================== DESCRIBE — accept (leading form) =============
	{"DESCRIBE t", true},
	{"DESC t", true},
	{"DESCRIBE namespace.foo", true},
	{"DESCRIBE INDEX myindex", true},
	{"DESCRIBE t FROM s", true},

	// ===================== ANALYZE — bare (accept; DDL path) ============
	{"ANALYZE", true},

	// ===================== RENAME — accept (DDL path) ===================
	{"RENAME TABLE a TO b", true},
	{"RENAME TABLE a.b TO c.d", true},

	// ===================== CALL — accept (full parse) ===================
	{"CALL my_proc()", true},
	{"CALL ns.proc()", true},
	{"CALL my_proc(1, 'x')", true},
	{"CALL my_proc(1 + 2, CAST(NULL AS string))", true},
	{"CALL ns.proc(a => 1)", true},
	{"CALL ns.proc(a => 1, b => 2)", true},
	{"CALL p(MODEL my.model)", true},
	{"CALL p(CONNECTION DEFAULT)", true},
	{"CALL p(CONNECTION my.connection)", true},
	{"CALL p(TABLE my.table)", true},
	{"CALL p(TABLE my.table, (SELECT * FROM s), f(1, 2))", true},
	{"CALL p(DESCRIPTOR(a, b))", true},
	{"CALL p(@param)", true},
	{"CALL p(@@sysvar)", true},
	// Clause keywords used as bare identifier / field-access EXPRESSIONS (the
	// keyword opens its clause only when a clause body follows).
	{"CALL p(model)", true},
	{"CALL p(table)", true},
	{"CALL p(connection)", true},
	{"CALL p(descriptor)", true},
	{"CALL p(model.col)", true},
	{"CALL p(table.x)", true},
	// Clause forms where the body follows the keyword.
	{"CALL p(MODEL x)", true},
	{"CALL p(TABLE x)", true},
	{"CALL p(CONNECTION x)", true},
	// named_argument value may be a lambda.
	{"CALL p(f => x -> x)", true},
	{"CALL p(f => (x) -> x)", true},

	// ===================== CALL — reject (authoritative syntax) =========
	{"CALL", false},                 // missing name
	{"CALL proc", false},            // missing '(' arg list
	{"CALL proc(", false},           // unterminated
	{"CALL proc(,)", false},         // leading comma
	{"CALL proc(1,)", false},        // trailing comma
	{"CALL proc(1 => 2)", false},    // named-arg name must be an identifier
	{"CALL p(descriptor x)", false}, // bare `descriptor` expr then stray `x`
	{"CALL p(model AS x)", false},   // model_clause needs a path, got AS
	{"CALL p(x -> x)", false},       // a bare positional lambda is not a tvf_argument
	{"CALL p(arr, e -> e)", false},  // ditto in a later positional arg
}

func TestUtilityDifferential(t *testing.T) {
	if os.Getenv("SPANNER_EMULATOR_HOST") == "" {
		t.Skip("SPANNER_EMULATOR_HOST not set; skipping live differential")
	}
	h := newSelectHarness(t)
	defer h.close()

	for _, fx := range utilityFixtures {
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
