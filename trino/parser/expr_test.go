package parser

import (
	"strings"
	"testing"
)

// This file is the correctness gate for the `expressions` DAG node
// (expr.go + function.go + predicate.go). Following correctness-protocol.md:
//
//   - Layer 1 (accept/reject differential): exprAcceptCorpus and exprRejectCorpus
//     span every primaryExpression / valueExpression / booleanExpression /
//     predicate_ alternative of TrinoParser.g4 plus oracle-discovered negatives.
//     TestExpr_OracleDifferential is the authoritative gate: ParseExpression's
//     accept/reject must equal Trino 481's accept/reject of `SELECT <expr>`.
//   - Layer 2 (structural): TestExpr_Precedence pins the parse-tree shape of the
//     oracle-confirmed precedence facts (single non-associative predicate;
//     ||/arithmetic/unary/AT-TIME-ZONE layering), and TestExpr_RoundTrip checks a
//     parse → render → re-parse is stable.
//   - Completeness: every grammar alternative appears in the corpus with its
//     oracle-derived verdict; the per-alternative checklist is the corpus comment
//     blocks below.
//
// The oracle helpers (connectOracle / oracleAccepts / truncateName) live in
// oracle_foundation_test.go and lexer_oracle_test.go.

// wrapExpr wraps a bare expression in a SELECT so the live Trino oracle (which
// speaks the statement protocol, not the standalone-expression entry) can judge
// it. `SELECT <expr>` reaches Trino's expression parser; a COLUMN_NOT_FOUND or
// other semantic error still counts as syntactically ACCEPTED (the oracle keys
// only on SYNTAX_ERROR), so bare column references like `a`, `b` are fine.
func wrapExpr(expr string) string { return "SELECT " + expr }

// exprAcceptCorpus is every expression form omni must ACCEPT, organised by the
// grammar alternative it exercises. Each is a bare expression; the oracle
// differential wraps it as `SELECT <expr>`.
var exprAcceptCorpus = []string{
	// --- literals (nullLiteral/booleanLiteral/stringLiteral/numericLiteral/binaryLiteral) ---
	"NULL",
	"TRUE",
	"FALSE",
	"'a string'",
	"'it''s escaped'",
	"U&'\\0041'",
	"123",
	"1.5",
	"1.5e10",
	"-3",
	"+4",
	"X'00ff'",
	"?", // parameter

	// --- typeConstructor (identifier string_ / DOUBLE PRECISION string_) ---
	"DATE '2020-01-01'",
	"TIMESTAMP '2012-10-31 01:00:00'",
	"BIGINT '7'",
	"json '[1,2]'",
	"DOUBLE PRECISION '1.5'",

	// --- intervalLiteral (interval rule) ---
	"INTERVAL '1' DAY",
	"INTERVAL '1-2' YEAR TO MONTH",
	"INTERVAL '3' HOUR",
	"INTERVAL -'5' MINUTE",
	"INTERVAL '1 02:03:04' DAY TO SECOND",

	// --- column reference / dereference / subscript ---
	"a",
	"a.b.c",
	"a[1]",
	"m['k'][1].f",
	"a.b[1]",

	// --- arithmetic / concat / unary (valueExpression) ---
	"a + b - c * d / e % f",
	"a || b || c",
	"- - 1",
	"-a * b",

	// --- AT TIME ZONE (atTimeZone) ---
	"ts AT TIME ZONE 'UTC'",
	"ts AT TIME ZONE INTERVAL '1' HOUR",
	"a * b AT TIME ZONE 'UTC'",

	// --- boolean (and/or/logicalNot/predicated) ---
	"TRUE AND FALSE OR NULL",
	"NOT a",
	"NOT NOT a",
	"NOT a = b",
	"NOT a AND b",

	// --- predicates (predicate_) ---
	"a <> b",
	"a != b",
	"a <= b",
	"a >= b",
	"a < b",
	"a > b",
	"a = b",
	"a = ANY (SELECT 1)",         // quantifiedComparison
	"a < ALL (SELECT 1)",         // quantifiedComparison
	"a > SOME (SELECT x FROM t)", // quantifiedComparison
	"a BETWEEN 1 AND 2",          // between
	"a NOT BETWEEN 1 AND 2",      // between (negated)
	"a IN (1, 2, 3)",             // inList
	"a NOT IN (1, 2, 3)",         // inList (negated)
	"a IN (SELECT 1)",            // inSubquery
	"a NOT IN (SELECT 1)",        // inSubquery (negated)
	"a LIKE 'x%'",                // like
	"a NOT LIKE 'x%'",            // like (negated)
	"a LIKE 'x%' ESCAPE '!'",     // like + escape
	"a IS NULL",                  // nullPredicate
	"a IS NOT NULL",              // nullPredicate
	"a IS DISTINCT FROM b",       // distinctFrom
	"a IS NOT DISTINCT FROM b",   // distinctFrom (negated)
	"a || b = c",                 // || binds tighter than predicate (P2)

	// --- constructors (rowConstructor/arrayConstructor) ---
	"(1, 2)",
	"(1, 2, 3)",
	"ROW(1, 'a')",
	"ROW(1)",
	"ARRAY[1, 2, 3]",
	"ARRAY[]",
	"(1)", // parenthesizedExpression

	// --- subquery / exists (subqueryExpression/exists) ---
	"(SELECT 1)",
	"(SELECT max(x) FROM (SELECT a AS x FROM t) s)", // nested-paren placeholder
	"EXISTS (SELECT 1)",
	"EXISTS (SELECT * FROM t WHERE a > 0)",
	"EXISTS (TABLE t)",
	"EXISTS (VALUES 1)",
	"a IN (SELECT b FROM t WHERE c IN (1, 2))", // nested IN, depth-balanced

	// --- CASE (simpleCase/searchedCase) ---
	"CASE WHEN a > 1 THEN 'hi' ELSE 'lo' END",
	"CASE a WHEN 1 THEN 'one' WHEN 2 THEN 'two' END",
	"CASE WHEN a THEN 1 END",

	// --- CAST / TRY_CAST (cast) ---
	"CAST(a AS varchar)",
	"TRY_CAST(b AS bigint)",
	"CAST(x AS ARRAY(INTEGER))",
	"CAST(x AS ROW(a bigint, b varchar))",
	"CAST(x AS MAP(VARCHAR, ARRAY(INTEGER)))",
	"CAST(x AS INTERVAL DAY TO SECOND)",
	"CAST(ts AS timestamp) AT TIME ZONE 'UTC'",

	// --- function calls (functionCall + decorators) ---
	"count(*)",
	"count(a)",
	"count(DISTINCT a)",
	"my.catalog.func(1)",
	"\"f\"(1)",
	"coalesce(a, b, 0)",
	"nullif(a, b)",
	"if(a > 0, 1, 2)",
	"count(*) FILTER (WHERE a > 0)",
	"count(t.*)",
	"f(t.*)",
	"array_agg(a ORDER BY b)",
	"array_agg(DISTINCT a ORDER BY a)",
	"array_agg(a ORDER BY b DESC NULLS LAST)",
	"sum(x) OVER (PARTITION BY a ORDER BY b)",
	"sum(x) OVER ()",
	"sum(x) OVER w",
	"row_number() OVER ()",
	"lag(x, 1) IGNORE NULLS OVER (ORDER BY y)",
	"lag(x) RESPECT NULLS OVER (ORDER BY y)",
	"count(DISTINCT a ORDER BY a) FILTER (WHERE a > 0) OVER (PARTITION BY b)",
	"FINAL first(x)", // processingMode

	// --- window frames (frameExtent/frameBound) ---
	"sum(x) OVER (PARTITION BY a ORDER BY b ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW)",
	"sum(x) OVER (ORDER BY b RANGE 5 PRECEDING)",
	"sum(x) OVER (ORDER BY b GROUPS BETWEEN 1 PRECEDING AND 1 FOLLOWING)",
	"sum(x) OVER (ORDER BY b ROWS UNBOUNDED PRECEDING)",
	"sum(x) OVER (ROWS CURRENT ROW)",
	"sum(x) OVER (ORDER BY b RANGE BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED FOLLOWING)",
	"sum(x) OVER (ORDER BY b ROWS BETWEEN CURRENT ROW AND UNBOUNDED FOLLOWING)",
	// named base window (windowSpecification existingWindowName)
	"sum(x) OVER (w ORDER BY b)",
	"sum(x) OVER (w PARTITION BY a)",
	"sum(x) OVER (w ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW)",

	// --- measure (identifier over) ---
	"a OVER (ORDER BY b)",

	// --- lambda ---
	"x -> x + 1",
	"() -> 1",
	"(x) -> x",
	"(k, v) -> k + v",
	"filter(arr, x -> x > 0)",
	"reduce(a, 0, (s, x) -> s + x, s -> s)",

	// --- special built-ins (extract/substring/trim/normalize/position/grouping/listagg) ---
	"EXTRACT(YEAR FROM a)",
	"SUBSTRING('abc' FROM 1 FOR 2)",
	"SUBSTRING('abc' FROM 1)",
	"TRIM(BOTH 'x' FROM 'xax')",
	"TRIM(LEADING FROM '  x')",
	"TRIM('y' FROM '  x  ')",
	"TRIM('  x  ')",
	"TRIM('abc', 'a')",
	"NORMALIZE('abc')",
	"NORMALIZE('abc', NFC)",
	"POSITION('a' IN 'abc')",
	// SUBSTRING / POSITION are non-reserved: bare column ref + comma-arg call
	// forms are also valid (oracle-confirmed), in addition to the special syntax.
	"substring",
	"position",
	"substring('abc', 1, 2)",
	"substring('abc', 1)",
	"position('a', 'abc')",
	"GROUPING(a, b)",
	"GROUPING(a)",
	"listagg(x, ',') WITHIN GROUP (ORDER BY x)",
	"listagg(DISTINCT x) WITHIN GROUP (ORDER BY x)",
	"listagg(x ON OVERFLOW ERROR) WITHIN GROUP (ORDER BY x)",
	"listagg(x ON OVERFLOW TRUNCATE '...' WITH COUNT) WITHIN GROUP (ORDER BY x)",

	// --- specialDateTimeFunction + current* set ---
	"CURRENT_DATE",
	"CURRENT_TIME",
	"CURRENT_TIME(3)",
	"CURRENT_TIMESTAMP",
	"CURRENT_TIMESTAMP(6)",
	"LOCALTIME",
	"LOCALTIMESTAMP",
	"CURRENT_USER",
	"CURRENT_CATALOG",
	"CURRENT_SCHEMA",
	"CURRENT_PATH",
}

// exprRejectCorpus is every expression form omni must REJECT, with the reason
// each is a SYNTAX_ERROR in Trino 481 (confirmed via the oracle).
var exprRejectCorpus = []string{
	// incomplete operators / clauses
	"1 +",                  // dangling binary operator
	"a AT TIME ZONE",       // missing zone
	"CAST(1 AS)",           // missing target type
	"a BETWEEN 1",          // missing AND upper
	"EXTRACT(FROM a)",      // missing field
	"TRIM(LEADING FROM)",   // missing trim source
	"TRIM(FROM 'x')",       // bare FROM (no spec/char) — Trino rejects despite legacy grammar
	"NORMALIZE()",          // missing source
	"LISTAGG(x)",           // missing WITHIN GROUP
	"ARRAY[1, 2,",          // unterminated array
	"CASE END",             // CASE with no WHEN
	"CASE WHEN THEN 1 END", // WHEN with no condition

	// non-associative predicate chains (oracle-confirmed P1)
	"a = b = c",
	"a < b < c",
	"a IN (1) IN (2)",
	"a IS NULL IS NULL",
	"a BETWEEN 1 AND 2 BETWEEN 3 AND 4",

	// IS variants Trino does not have
	"a IS TRUE",
	"a IS NOT FALSE",
	"a IS UNKNOWN",

	// reserved special-function keywords reject ordinary readings
	"trim",              // reserved: not a bare column reference
	"extract",           // reserved: not a bare column reference
	"extract('a', 'b')", // EXTRACT has only the `field FROM source` form
	"listagg('a')",      // LISTAGG requires WITHIN GROUP (ORDER BY …)
	"grouping('a')",     // GROUPING arguments are qualifiedNames, not strings

	// EXISTS requires a real subquery — no empty/expression form
	"EXISTS ()",
	"EXISTS (1)",
	"EXISTS (1 + 2)",

	// empty paren / subquery in scalar & IN positions
	"()",
	"a IN ()",

	// SQL/JSON functions are the expr-json node (this node rejects them as
	// "not yet supported" — a deliberate, node-scoped non-acceptance that the
	// oracle differential treats specially; see TestExpr_JSONDeferred).
}

// TestExpr_AcceptCorpusParses verifies omni accepts every form in
// exprAcceptCorpus (oracle-free completeness smoke test). The oracle differential
// proves the verdicts actually match Trino.
func TestExpr_AcceptCorpusParses(t *testing.T) {
	for _, expr := range exprAcceptCorpus {
		t.Run(truncateName(expr), func(t *testing.T) {
			node, errs := ParseExpression(expr)
			if len(errs) != 0 {
				t.Errorf("ParseExpression(%q) should accept, got errors: %v", expr, errs)
			}
			if node == nil {
				t.Errorf("ParseExpression(%q) returned nil node", expr)
			}
		})
	}
}

// TestExpr_RejectCorpusRejected verifies omni rejects every form in
// exprRejectCorpus (required negative coverage).
func TestExpr_RejectCorpusRejected(t *testing.T) {
	for _, expr := range exprRejectCorpus {
		t.Run(truncateName(expr), func(t *testing.T) {
			_, errs := ParseExpression(expr)
			if len(errs) == 0 {
				t.Errorf("ParseExpression(%q) should reject, but accepted", expr)
			}
		})
	}
}

// TestExpr_OracleDifferential is the authoritative accept/reject gate: for every
// form in both corpora, omni's ParseExpression accept/reject must equal Trino
// 481's accept/reject of `SELECT <expr>`. Skipped when no oracle is reachable.
func TestExpr_OracleDifferential(t *testing.T) {
	if testing.Short() {
		t.Skip("trino oracle: skipped in -short mode")
	}
	o := connectOracle(t)

	check := func(t *testing.T, expr string) {
		_, errs := ParseExpression(expr)
		omniAccepts := len(errs) == 0

		trinoAccepts, ok := oracleAccepts(t, o, wrapExpr(expr))
		if !ok {
			t.Skip("oracle unreachable for this case")
		}
		if omniAccepts != trinoAccepts {
			t.Errorf("MISMATCH expr=%q: omni accepts=%v (errs=%v), Trino accepts=%v",
				expr, omniAccepts, errs, trinoAccepts)
		}
	}

	t.Run("accept", func(t *testing.T) {
		for _, expr := range exprAcceptCorpus {
			expr := expr
			t.Run(truncateName(expr), func(t *testing.T) { check(t, expr) })
		}
	})
	t.Run("reject", func(t *testing.T) {
		for _, expr := range exprRejectCorpus {
			expr := expr
			t.Run(truncateName(expr), func(t *testing.T) { check(t, expr) })
		}
	})
}

// TestExpr_JSONDeferred documents the node-scope boundary: the SQL/JSON
// functions are accepted by Trino but deferred to the expr-json node, so this
// node rejects them with a "not yet supported" diagnostic. This is a deliberate
// non-acceptance recorded in the divergence ledger; it is NOT a Trino-disagreement
// bug. When expr-json lands, these move into the accept corpus.
func TestExpr_JSONDeferred(t *testing.T) {
	jsonForms := []string{
		"JSON_EXISTS(a, 'lax $.b')",
		"JSON_VALUE(a, 'lax $.b')",
		"JSON_QUERY(a, 'lax $.b')",
		"JSON_OBJECT('k' VALUE 1)",
		"JSON_ARRAY(1, 2)",
	}
	for _, expr := range jsonForms {
		t.Run(truncateName(expr), func(t *testing.T) {
			_, errs := ParseExpression(expr)
			if len(errs) == 0 {
				t.Errorf("expected %q to be rejected as not-yet-supported, but accepted", expr)
				return
			}
			if !strings.Contains(errs[0].Msg, "not yet supported") {
				t.Errorf("expected a 'not yet supported' diagnostic for %q, got: %v", expr, errs[0].Msg)
			}
		})
	}
}
