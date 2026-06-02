package parser

import (
	"context"
	"testing"
	"time"

	"github.com/bytebase/omni/trino/internal/trinooracle"
)

// This file is the parser-foundation node's slice of the differential-oracle
// gate (correctness-protocol.md). The foundation ships the parser entry point,
// the statement splitter, identifier/qualifiedName parsing, and diagnostics —
// but every *statement body* is still stubbed (parseStmt routes to
// unsupported), so a naive "omni Parse accept == Trino accept" differential is
// not yet meaningful: omni rejects every valid statement with "not yet
// supported" until later DAG nodes implement the bodies.
//
// What the foundation CAN adjudicate against the live oracle, and does here:
//
//  1. Identifier / qualifiedName accept-reject. This is the only grammar the
//     foundation actually decides. For each candidate name we compare
//     ParseQualifiedName's verdict to Trino's verdict for the name used in an
//     unambiguous qualifiedName position (`SELECT * FROM <name>`). Both
//     positive and negative (reserved/digit-leading) cases are covered, as the
//     protocol requires.
//  2. Splitter round-trip. For multi-statement inputs, each segment Split
//     produces must be a statement Trino accepts on its own — proving Split
//     cut on real statement boundaries and never mid-statement (including
//     inside string/comment/quoted-identifier and a CREATE FUNCTION routine
//     body).
//  3. Robustness. Parse/Diagnose must terminate without panic on every
//     accepted corpus statement and never invent a lex error there.
//
// All oracle-backed subtests skip cleanly when no Trino is reachable, matching
// the harness convention; the robustness sweep runs oracle-free too.

// connectOracle dials the live Trino oracle, skipping the test when it is
// unreachable. Shared by the oracle-backed subtests below.
func connectOracle(t *testing.T) *trinooracle.Oracle {
	t.Helper()
	o := trinooracle.Connect("")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ver, err := o.Ping(ctx)
	if err != nil {
		t.Skipf("trino oracle not reachable (start: docker run -d -p 18080:8080 %s): %v",
			trinooracle.DefaultImage, err)
	}
	t.Logf("connected to Trino %s", ver)
	return o
}

// oracleAccepts asks the oracle whether Trino syntactically accepts sql,
// treating an oracle transport error as a skip signal (returns ok=false).
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

// qualifiedNameCandidates spans the foundation's identifier/qualifiedName
// grammar surface: unquoted/quoted/back-quoted, non-reserved keywords as
// identifiers, multi-part names, reserved-word rejection, and the
// digit-leading rejection adjudicated by the oracle. Each is a name only;
// the differential wraps it as `SELECT * FROM <name>`.
var qualifiedNameCandidates = []string{
	// accepted
	"t",
	"orders",
	"public.orders",
	"tpch.sf1.nation",
	`"My Table"`,
	`"public"."order"`,
	`cat."My Schema".tbl`,
	"zone",      // non-reserved keyword as identifier
	"t.zone",    // non-reserved after dot
	"including", // non-reserved keyword
	// rejected
	"from",         // reserved bare
	"select",       // reserved bare
	"s.from",       // reserved after dot
	"public.order", // reserved after dot
	"1abc",         // digit-leading (oracle rejects despite legacy grammar)
	"t.1abc",       // digit-leading after dot
	"`bt table`",   // backtick quoting (oracle rejects despite legacy grammar)
	"t.",           // trailing dot
}

// TestFoundation_QualifiedNameDifferential cross-checks ParseQualifiedName
// against the live Trino oracle for every candidate name, using the
// `SELECT * FROM <name>` context to make the name a qualifiedName. omni's
// accept/reject of the standalone name must equal Trino's accept/reject of the
// wrapped statement.
func TestFoundation_QualifiedNameDifferential(t *testing.T) {
	if testing.Short() {
		t.Skip("trino oracle: skipped in -short mode")
	}
	o := connectOracle(t)
	for _, name := range qualifiedNameCandidates {
		name := name
		t.Run(truncateName(name), func(t *testing.T) {
			_, errs := ParseQualifiedName(name)
			omniAccepts := len(errs) == 0

			trinoAccepts, ok := oracleAccepts(t, o, "SELECT * FROM "+name)
			if !ok {
				t.Skip("oracle unreachable for this case")
			}
			if omniAccepts != trinoAccepts {
				t.Errorf("MISMATCH name=%q: omni accepts=%v (errs=%v), Trino accepts=%v",
					name, omniAccepts, errs, trinoAccepts)
			}
		})
	}
}

// multiStatementCorpus exercises the splitter on inputs with several top-level
// statements, including ';' hidden inside strings, comments, quoted
// identifiers, and a CREATE FUNCTION routine body. Each entry lists the input
// and the statements Split must produce (each of which must independently be
// accepted by Trino).
var multiStatementCorpus = []struct {
	name      string
	input     string
	wantStmts []string
}{
	{
		name:      "two_selects",
		input:     "SELECT 1; SELECT 2",
		wantStmts: []string{"SELECT 1", " SELECT 2"},
	},
	{
		name:      "trailing_semicolon",
		input:     "SELECT 1;",
		wantStmts: []string{"SELECT 1"},
	},
	{
		name:      "semicolon_in_string",
		input:     "SELECT 'a;b' AS c; SELECT 2",
		wantStmts: []string{"SELECT 'a;b' AS c", " SELECT 2"},
	},
	{
		name:      "semicolon_in_quoted_ident",
		input:     `SELECT 1 AS "a;b"; SELECT 2`,
		wantStmts: []string{`SELECT 1 AS "a;b"`, " SELECT 2"},
	},
	{
		name:      "semicolon_in_line_comment",
		input:     "SELECT 1 -- a;b\n; SELECT 2",
		wantStmts: []string{"SELECT 1 -- a;b\n", " SELECT 2"},
	},
	{
		name:      "semicolon_in_block_comment",
		input:     "SELECT 1 /* a;b */; SELECT 2",
		wantStmts: []string{"SELECT 1 /* a;b */", " SELECT 2"},
	},
	{
		name:      "ddl_then_dml",
		input:     "CREATE SCHEMA s; DROP SCHEMA s",
		wantStmts: []string{"CREATE SCHEMA s", " DROP SCHEMA s"},
	},
	{
		name: "create_function_routine_body_not_split",
		// The ';' separating the routine's control statements must NOT split
		// the CREATE FUNCTION into pieces; it is one statement.
		input:     "CREATE FUNCTION test.default.f(x integer) RETURNS integer BEGIN RETURN x * 2; END",
		wantStmts: []string{"CREATE FUNCTION test.default.f(x integer) RETURNS integer BEGIN RETURN x * 2; END"},
	},
	{
		name: "create_function_nested_if_not_split",
		// A nested IF (with its own internal ';') inside the BEGIN block keeps
		// the whole CREATE FUNCTION as a single segment.
		input:     "CREATE FUNCTION test.default.g(x integer) RETURNS integer BEGIN IF x > 0 THEN RETURN 1; END IF; RETURN 0; END",
		wantStmts: []string{"CREATE FUNCTION test.default.g(x integer) RETURNS integer BEGIN IF x > 0 THEN RETURN 1; END IF; RETURN 0; END"},
	},
	{
		name: "inline_routine_bare_return_then_select",
		// WITH FUNCTION with a bare RETURN body (no block): the inline-routine
		// query is one statement, and a trailing ';' splits the next.
		input:     "WITH FUNCTION f(x integer) RETURNS integer RETURN x * 2 SELECT f(1); SELECT 2",
		wantStmts: []string{"WITH FUNCTION f(x integer) RETURNS integer RETURN x * 2 SELECT f(1)", " SELECT 2"},
	},
	{
		name: "function_typed_end_closer_then_split",
		// A typed END IF closer inside the routine body must keep depth
		// balanced so the ';' ending the CREATE FUNCTION splits the trailing
		// SELECT off (both segments independently Trino-accepted).
		input:     "CREATE FUNCTION test.default.g(x integer) RETURNS integer BEGIN IF x > 0 THEN RETURN 1; END IF; RETURN 0; END; SELECT 2",
		wantStmts: []string{"CREATE FUNCTION test.default.g(x integer) RETURNS integer BEGIN IF x > 0 THEN RETURN 1; END IF; RETURN 0; END", " SELECT 2"},
	},
}

// TestFoundation_SplitMatchesCount verifies (oracle-free) that Split produces
// the expected segment text for each multi-statement input.
func TestFoundation_SplitMatchesCount(t *testing.T) {
	for _, tc := range multiStatementCorpus {
		t.Run(tc.name, func(t *testing.T) {
			segs := Split(tc.input)
			if len(segs) != len(tc.wantStmts) {
				t.Fatalf("Split(%q): got %d segments, want %d: %+v",
					tc.input, len(segs), len(tc.wantStmts), segs)
			}
			for i, seg := range segs {
				if seg.Text != tc.wantStmts[i] {
					t.Errorf("segment %d Text = %q, want %q", i, seg.Text, tc.wantStmts[i])
				}
				// The byte range must address exactly the segment text.
				if got := tc.input[seg.ByteStart:seg.ByteEnd]; got != seg.Text {
					t.Errorf("segment %d range [%d,%d) -> %q, want %q",
						i, seg.ByteStart, seg.ByteEnd, got, seg.Text)
				}
			}
		})
	}
}

// TestFoundation_SplitRoundTripOracle verifies each statement Split produces is
// independently accepted by Trino — proving Split cut on real boundaries and
// never mid-statement. Skipped without an oracle.
func TestFoundation_SplitRoundTripOracle(t *testing.T) {
	if testing.Short() {
		t.Skip("trino oracle: skipped in -short mode")
	}
	o := connectOracle(t)
	for _, tc := range multiStatementCorpus {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			for _, seg := range Split(tc.input) {
				accepted, ok := oracleAccepts(t, o, seg.Text)
				if !ok {
					t.Skip("oracle unreachable for this case")
				}
				if !accepted {
					t.Errorf("Split produced a segment Trino rejects: %q (from input %q)", seg.Text, tc.input)
				}
			}
		})
	}
}

// TestFoundation_ParseRobustnessOnCorpus runs Parse and Diagnose over the
// shared oracleCorpus (real Trino 481 SQL) and asserts they always terminate
// without panic and never invent a lex error on accepted SQL. It is oracle-free
// (the corpus is curated all-accept) and complements the oracle differential.
func TestFoundation_ParseRobustnessOnCorpus(t *testing.T) {
	for _, sql := range oracleCorpus {
		// Parse must not panic and must always return a non-nil File.
		file, _ := Parse(sql)
		if file == nil {
			t.Errorf("Parse(%q) returned nil File", sql)
		}
		// Diagnose must not panic.
		_ = Diagnose(sql)
		// On accepted corpus SQL the lexer must not invent a hard error; any
		// diagnostic must be a "not yet supported" stub (parseStmt), never a
		// lex error. (lexerHasHardError lives in lexer_oracle_test.go.)
		if bad, errs := lexerHasHardError(sql); bad {
			t.Errorf("foundation: spurious lex error on accepted corpus SQL %q: %v", sql, errs)
		}
	}
}

// TestFoundation_NegativeInputsRejected feeds malformed SQL the engine rejects
// and asserts the foundation reports a diagnostic for each (never silently
// accepts). Required negative coverage per the correctness protocol.
func TestFoundation_NegativeInputsRejected(t *testing.T) {
	negatives := []string{
		"SELECT 'unterminated",   // lex error: unterminated string
		`SELECT "unterminated`,   // lex error: unterminated quoted ident
		"SELECT X'00",            // lex error: unterminated binary literal
		"SELECT /* unterminated", // lex error: unterminated comment
		"FROBNICATE foo",         // unknown statement keyword
		"\x00",                   // unrecognized character
	}
	for _, sql := range negatives {
		diags := Diagnose(sql)
		if len(diags) == 0 {
			t.Errorf("Diagnose(%q): want at least one diagnostic, got none", sql)
		}
	}
}
