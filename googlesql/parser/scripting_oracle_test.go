//go:build googlesql_oracle

// Differential test for the `parser-scripting` node against the live Cloud
// Spanner emulator oracle. Run with:
//
//	SPANNER_EMULATOR_HOST=localhost:9010 \
//	  go test -tags googlesql_oracle ./googlesql/parser/ -run TestScriptingDifferential
//
// It is the PROVE gate for googlesql/parser-scripting (correctness-protocol.md):
// for every AUTHORITATIVE fixture it (1) feeds the statement to the emulator via
// the harness/googlesql-spanner CLI and reads the accept/reject verdict, then
// (2) parses the same statement with omni's Parse, and asserts the two agree —
// BOTH polarities. A harness `error` verdict fails loudly; it is never folded
// into accept or reject. The harness plumbing (newSelectHarness / verdict /
// close) is shared with select_oracle_test.go.
//
// DIALECT (oracle.md). Scripting is BigQuery-ONLY at the GoogleSQL union level,
// with three oracle-authoritative LEADING forms — the emulator PARSES them
// (verdict accept, "Statement not supported: …"):
//   - BEGIN…END block            → "BeginStmt"
//   - SET                        → "SingleAssignment" / "AssignmentFromStruct" /
//                                  "ParameterAssignment" / "SystemVariableAssignment" /
//                                  "SetTransactionStatement"
//   - EXECUTE IMMEDIATE          → "ExecuteImmediateStatement"
// plus genuine SHARED syntax rejects (a bare `SET`, `SET x` with no `= expr`).
// These are diffed here.
//
// EXCLUDED — NON-authoritative on Spanner (covered by scripting_test.go + the
// divergence ledger, NOT here, because the emulator's verdict would manufacture
// a false divergence):
//   - The control-flow statements (IF / CASE / WHILE / LOOP / REPEAT / FOR-IN /
//     DECLARE / BREAK / CONTINUE / RETURN / RAISE / labels). The emulator
//     SYNTAX-rejects them ("Unexpected keyword IF/WHILE/…"); they are
//     BigQuery-valid (the .g4 procedural grammar + truth1 SCRIPT-001…019) and
//     omni accepts them. A Spanner reject here would be a false divergence.
//     These are asserted omni-side in scriptingBigQueryOnly below.
//   - The INTERIOR statement_list grammar of a BEGIN…END block. The emulator's
//     BEGIN…END recognizer is SHALLOW — it accepts ANY token run between BEGIN
//     and END (`BEGIN SELECT 1 SELECT 2 END`, `BEGIN SELECT 1 END` with no
//     trailing `;`, `BEGIN; END`), while the .g4 requires `;`-separated
//     statements and a trailing `;`. omni follows the .g4, so those shallow
//     over-accepts are deliberate divergences (ledger), excluded from this diff
//     and asserted omni-side in scriptingBlockInteriorReject.

package parser

import (
	"os"
	"testing"
)

// scriptingFixture is one statement fed to BOTH the oracle and omni's Parse.
// wantParse is the expected omni outcome and MUST equal the emulator's grammar
// verdict.
type scriptingFixture struct {
	sql       string
	wantParse bool
}

// scriptingAuthoritative are the forms whose Spanner verdict is authoritative
// (the emulator parses the leading form, or genuinely syntax-rejects) AND on
// which omni and the oracle agree.
var scriptingAuthoritative = []scriptingFixture{
	// ===================== BEGIN…END — accept (envelope) ================
	{"BEGIN END", true},
	{"BEGIN SELECT 1; END", true},
	{"BEGIN DECLARE x INT64; SET x = 1; SELECT x; END", true},
	{"BEGIN SELECT 1; EXCEPTION WHEN ERROR THEN SELECT 2; END", true},
	{"BEGIN BEGIN SELECT 1; END; SELECT 2; END", true},
	// A labeled BEGIN…END is parsed by the emulator's shallow recognizer as a
	// BeginStmt (the label is swallowed) — accept on both sides.
	{"label_1: BEGIN SELECT 1; BREAK label_1; END", true},

	// ===================== SET — accept (parsed) ========================
	{"SET x = 5", true},
	{"SET x = 1 + 2 * 3", true},
	{"SET (a, b, c) = (1, 'foo', false)", true},
	{"SET (a) = (1)", true},
	{"SET @p = 5", true},
	{"SET @@x.y = 5", true},
	{"SET TRANSACTION READ ONLY", true},
	{"SET TRANSACTION READ WRITE", true},
	{"SET TRANSACTION ISOLATION LEVEL SERIALIZABLE", true},

	// ===================== EXECUTE IMMEDIATE — accept (parsed) ===========
	{`EXECUTE IMMEDIATE "SELECT 1"`, true},
	{`EXECUTE IMMEDIATE "SELECT ? * (? + 2)" INTO y USING 1, 3`, true},
	{`EXECUTE IMMEDIATE "SELECT @a" INTO y USING 1 AS a`, true},
	{`EXECUTE IMMEDIATE "SELECT 1, 2" INTO a, b`, true},
	{`EXECUTE IMMEDIATE CONCAT("SELECT ", "1")`, true},

	// ===================== SET — reject (shared syntax) =================
	{"SET", false},   // missing target ("Unexpected end of statement")
	{"SET x", false}, // missing '=' / ',' ("Expected ',' or '='")

	// ===================== EXECUTE — reject (shared syntax) =============
	{"EXECUTE", false}, // missing IMMEDIATE
}

func TestScriptingDifferential(t *testing.T) {
	if os.Getenv("SPANNER_EMULATOR_HOST") == "" {
		t.Skip("SPANNER_EMULATOR_HOST not set; skipping live differential")
	}
	h := newSelectHarness(t)
	defer h.close()

	for _, fx := range scriptingAuthoritative {
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
			// non-authoritative on Spanner and should move to the omni-only sets).
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

// scriptingBigQueryOnly are control-flow forms the Spanner emulator
// syntax-rejects (BigQuery-only) but which omni must ACCEPT on the authority of
// the legacy .g4 + the BigQuery truth1 corpus. They are NOT diffed against the
// oracle (its reject is non-authoritative); they assert omni's accept and, when
// the emulator is reachable, that the emulator indeed REJECTS (so the BigQuery-
// only classification stays honest — if a future emulator started parsing these,
// this test flags that the form should move to the authoritative set).
var scriptingBigQueryOnly = []string{
	"DECLARE x INT64",
	"DECLARE d DATE DEFAULT CURRENT_DATE()",
	"DECLARE x, y, z INT64 DEFAULT 0",
	"DECLARE item DEFAULT (SELECT col FROM s.products LIMIT 1)",
	"IF x > 0 THEN SELECT 1; END IF",
	"IF c1 THEN SELECT 1; ELSEIF c2 THEN SELECT 2; ELSE SELECT 3; END IF",
	"CASE WHEN c THEN SELECT 1; ELSE SELECT 2; END CASE",
	"CASE x WHEN 1 THEN SELECT 'one'; ELSE SELECT 'other'; END CASE",
	"WHILE x < 10 DO SET x = x + 1; END WHILE",
	"LOOP SET x = x + 1; IF x >= 10 THEN LEAVE; END IF; END LOOP",
	"REPEAT SET x = x + 1; UNTIL x >= 3 END REPEAT",
	"FOR rec IN (SELECT word FROM s.shakespeare LIMIT 5) DO SELECT rec.word; END FOR",
	"BREAK",
	"LEAVE my_label",
	"CONTINUE",
	"ITERATE outer_loop",
	"RETURN",
	"RAISE",
	"RAISE USING MESSAGE = 'error'",
	// A labeled LOOP/WHILE/REPEAT/FOR is BigQuery-only: the emulator
	// syntax-rejects it ("Unexpected <label>"). (A labeled BEGIN…END, by contrast,
	// the emulator's shallow recognizer ACCEPTS as a BeginStmt — the label is
	// swallowed — so it lives in scriptingAuthoritative, not here.)
	"label_1: LOOP SELECT 1; END LOOP label_1",
	"lbl: WHILE x DO SELECT 1; END WHILE",
}

func TestScriptingBigQueryOnlyAcceptedByOmni(t *testing.T) {
	// omni must accept every BigQuery-only control-flow form.
	for _, sql := range scriptingBigQueryOnly {
		sql := sql
		t.Run(sql, func(t *testing.T) {
			_, errs := Parse(sql)
			if len(errs) != 0 {
				t.Errorf("omni rejected BigQuery-only form %q: %v", sql, errs)
			}
		})
	}

	// When the emulator is reachable, confirm it SYNTAX-rejects these (the basis
	// for treating them as BigQuery-only / non-authoritative). A standalone
	// labeled top-level form may be reported differently; we only require that the
	// emulator does NOT cleanly accept it as a Spanner statement.
	if os.Getenv("SPANNER_EMULATOR_HOST") == "" {
		return
	}
	h := newSelectHarness(t)
	defer h.close()
	for _, sql := range scriptingBigQueryOnly {
		sql := sql
		t.Run("oracle-rejects/"+sql, func(t *testing.T) {
			v := h.verdict(t, sql)
			if v.Verdict == "accept" {
				t.Errorf("emulator now ACCEPTS %q (%s) — it is no longer Spanner-non-authoritative; "+
					"move it to scriptingAuthoritative and diff it", sql, v.Message)
			}
		})
	}
}

// scriptingBlockInteriorReject documents the deliberate divergence: the Spanner
// emulator's BEGIN…END recognizer is SHALLOW and ACCEPTS these, but the .g4
// statement_list grammar (which bytebase consumes) requires `;`-separated
// statements and a trailing `;`, so omni REJECTS them. This asserts omni's
// reject and (when reachable) records that the emulator over-accepts.
var scriptingBlockInteriorReject = []string{
	"BEGIN SELECT 1 END",           // no trailing ';' after the single statement
	"BEGIN SELECT 1 SELECT 2 END",  // no ';' between statements
	"BEGIN SELECT 1; SELECT 2 END", // no trailing ';' after the last statement
}

func TestScriptingBlockInteriorRejectedByOmni(t *testing.T) {
	for _, sql := range scriptingBlockInteriorReject {
		sql := sql
		t.Run(sql, func(t *testing.T) {
			_, errs := Parse(sql)
			if len(errs) == 0 {
				t.Errorf("omni accepted %q, but the .g4 statement_list requires `;` "+
					"termination — omni must reject (Spanner's shallow recognizer over-accepts)", sql)
			}
		})
	}

	// Confirm the emulator's shallow over-accept (the basis for the divergence).
	if os.Getenv("SPANNER_EMULATOR_HOST") == "" {
		return
	}
	h := newSelectHarness(t)
	defer h.close()
	for _, sql := range scriptingBlockInteriorReject {
		sql := sql
		t.Run("oracle-overaccepts/"+sql, func(t *testing.T) {
			v := h.verdict(t, sql)
			if v.Verdict != "accept" {
				t.Logf("note: emulator verdict for %q is %q (%s) — the documented shallow "+
					"over-accept may have changed", sql, v.Verdict, v.Message)
			}
		})
	}
}
