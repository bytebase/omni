package parser

import (
	"testing"
)

// This file is the parser-routines node's differential-oracle gate
// (correctness-protocol.md): for every statement in the corpus below, omni's
// Parse accept/reject verdict must equal Trino 481's verdict as reported by the
// live oracle (trinooracle.CheckSyntax, which classifies a SYNTAX_ERROR as
// reject and every other outcome — success or a semantic error like
// NOT_SUPPORTED / FUNCTION_NOT_FOUND / TYPE_MISMATCH — as accept). Most
// CREATE/DROP FUNCTION forms fail at execution with NOT_SUPPORTED (the disposable
// `memory` connector cannot store a UDF), which the oracle correctly classifies
// as ACCEPTED — the parser+analyzer succeeded.
//
// The corpus is NOT split into hardcoded accept/reject lists: the oracle is the
// source of truth, so each case is run through both omni and Trino and the two
// verdicts are compared. This both proves correctness and documents which
// docs/legacy forms Trino 481 actually accepts. It covers every documented form
// (truth1 create-function / drop-function), every legacy controlStatement and
// routineCharacteristic alternative (truth2), the oracle-confirmed docs-ahead
// WITH (properties) characteristic (the flagged divergence), and a set of
// negative forms that pin the grammar distinctions (RETURN takes a
// valueExpression not a boolean; ITERATE/LEAVE require a label; DECLARE precede
// statements; labels only on LOOP/WHILE/REPEAT; sqlStatementList ';'-termination).
//
// Skipped cleanly when no Trino oracle is reachable.

// routinesOracleCorpus is the full SQL-routine corpus for the differential gate.
// Each string is fed to both omni and the oracle; the verdicts must match.
//
// NOTE: statements carry NO trailing ';'. omni's Parse runs Split first (which
// strips a trailing ';'), while the oracle receives the raw string over the REST
// API and rejects an embedded extra ';'; omitting it keeps the two comparable
// (it is how the real pipeline compares, and matches the merged dcl-tcl / utility
// corpora).
var routinesOracleCorpus = []string{
	// ---------------------------------------------------------------------
	// CREATE FUNCTION — minimal & docs (truth1 create-function)
	// ---------------------------------------------------------------------
	"CREATE FUNCTION meaning_of_life() RETURNS bigint RETURN 42",
	"CREATE FUNCTION example.default.meaning_of_life() RETURNS bigint BEGIN RETURN 42; END",
	"CREATE FUNCTION my.func(x integer) RETURNS integer RETURN x + 1",
	"CREATE OR REPLACE FUNCTION f(x bigint) RETURNS bigint RETURN x * x",

	// ---------------------------------------------------------------------
	// parameterDeclaration (F2): named / type-only / mixed / complex types
	// ---------------------------------------------------------------------
	"CREATE FUNCTION f(bigint) RETURNS bigint RETURN 1",
	"CREATE FUNCTION f(x bigint, varchar, z double) RETURNS bigint RETURN 1",
	"CREATE FUNCTION f(x ROW(a int, b int)) RETURNS int RETURN 1",
	"CREATE FUNCTION f(x ARRAY(int)) RETURNS int RETURN 1",
	"CREATE FUNCTION f(x DECIMAL(10, 2)) RETURNS int RETURN 1",
	"CREATE FUNCTION f(x timestamp(3) with time zone) RETURNS int RETURN 1",
	// Multi-token type-only params: two tokens of lookahead cannot tell these
	// apart from a named param, so they exercise the speculative resolution
	// (parseParameterDeclaration tries unnamed-type-first, then named).
	"CREATE FUNCTION f(double precision) RETURNS int RETURN 1",       // type-only DOUBLE PRECISION
	"CREATE FUNCTION f(x double precision) RETURNS int RETURN 1",     // named, DOUBLE PRECISION type
	"CREATE FUNCTION f(interval day to second) RETURNS int RETURN 1", // type-only INTERVAL DAY TO SECOND
	"CREATE FUNCTION f(comment bigint) RETURNS int RETURN 1",         // non-reserved keyword as a param name
	`CREATE FUNCTION f("my x" bigint) RETURNS int RETURN 1`,          // quoted param name
	"CREATE FUNCTION f() RETURNS varchar(10) RETURN 'x'",             // parameterized RETURNS type
	"CREATE FUNCTION f() RETURNS array(int) RETURN ARRAY[1]",         // complex RETURNS type

	// ---------------------------------------------------------------------
	// routineCharacteristic (truth2): every alternative, order & repetition (F3)
	// ---------------------------------------------------------------------
	"CREATE FUNCTION f(x bigint) RETURNS bigint LANGUAGE SQL DETERMINISTIC RETURNS NULL ON NULL INPUT RETURN x",
	"CREATE FUNCTION f(x bigint) RETURNS bigint NOT DETERMINISTIC CALLED ON NULL INPUT SECURITY DEFINER COMMENT 'doc' RETURN x",
	"CREATE FUNCTION f(x bigint) RETURNS bigint SECURITY INVOKER RETURN x",
	"CREATE FUNCTION f() RETURNS bigint COMMENT 'a' LANGUAGE SQL DETERMINISTIC RETURN 1",
	"CREATE FUNCTION f() RETURNS bigint DETERMINISTIC DETERMINISTIC RETURN 1", // repeated (routineCharacteristic*)
	"CREATE FUNCTION f() RETURNS bigint LANGUAGE PYTHON RETURN 1",
	"CREATE FUNCTION f() RETURNS bigint COMMENT U&'doc' RETURN 1", // unicode string COMMENT (string_ covers U&'…')

	// ---------------------------------------------------------------------
	// WITH (properties) characteristic — docs-ahead-of-legacy (F4, flagged)
	// ---------------------------------------------------------------------
	"CREATE FUNCTION f() RETURNS int LANGUAGE PYTHON WITH (handler = 'myfunc') RETURN 1",
	"CREATE FUNCTION f() RETURNS int WITH (handler = 'x') RETURN 1",
	"CREATE FUNCTION f() RETURNS int WITH (a = 'str', b = 2.5, c = true) RETURN 1",
	"CREATE FUNCTION f() RETURNS int WITH (a = ARRAY[1, 2]) RETURN 1",
	`CREATE FUNCTION f() RETURNS int WITH ("quoted name" = 1) RETURN 1`,
	"CREATE FUNCTION f() RETURNS int COMMENT 'x' WITH (a = 1) RETURN 1",
	"CREATE FUNCTION f() RETURNS int DETERMINISTIC WITH (a = 1) WITH (b = 2) RETURN 1",

	// ---------------------------------------------------------------------
	// controlStatement bodies (truth2) — compound, declarations, assignment
	// ---------------------------------------------------------------------
	"CREATE FUNCTION f(x bigint) RETURNS bigint BEGIN DECLARE y bigint DEFAULT 0; SET y = x + 1; RETURN y; END",
	"CREATE FUNCTION f() RETURNS bigint BEGIN RETURN 1; END",
	"CREATE FUNCTION f() RETURNS bigint BEGIN DECLARE a, b int; SET a = 1; SET b = 2; RETURN a + b; END",
	"CREATE FUNCTION f() RETURNS bigint BEGIN END",                      // empty compound (R6)
	"CREATE FUNCTION f() RETURNS bigint BEGIN BEGIN RETURN 1; END; END", // nested
	"CREATE FUNCTION f() RETURNS bigint BEGIN DECLARE x int; END",       // declares only (R6)
	"CREATE FUNCTION f() RETURNS int BEGIN DECLARE x int DEFAULT 1 + 2 * 3; RETURN x; END",
	"CREATE FUNCTION f() RETURNS int BEGIN DECLARE x int DEFAULT (SELECT 1); RETURN x; END",
	"CREATE FUNCTION f(x int) RETURNS int BEGIN DECLARE y int; SET y = x IN (1, 2); RETURN y; END", // SET takes full expression (R2)

	// IF / ELSEIF / ELSE
	"CREATE FUNCTION f(x bigint) RETURNS bigint BEGIN IF x > 0 THEN RETURN 1; ELSEIF x < 0 THEN RETURN -1; ELSE RETURN 0; END IF; END",
	"CREATE FUNCTION f(x bigint) RETURNS bigint BEGIN IF x > 0 THEN RETURN 1; END IF; RETURN 0; END",
	"CREATE FUNCTION f(x int) RETURNS int BEGIN IF x IN (1, 2, 3) THEN RETURN 1; END IF; RETURN 0; END", // IF cond is full expression (R2)

	// CASE (simple + searched)
	"CREATE FUNCTION f(x bigint) RETURNS varchar BEGIN CASE x WHEN 1 THEN RETURN 'one'; WHEN 2 THEN RETURN 'two'; ELSE RETURN 'other'; END CASE; END",
	"CREATE FUNCTION f(x bigint) RETURNS varchar BEGIN CASE WHEN x > 0 THEN RETURN 'pos'; ELSE RETURN 'neg'; END CASE; END",
	"CREATE FUNCTION f(x int) RETURNS int BEGIN CASE x + 1 WHEN 2 THEN RETURN 1; END CASE; END",     // simple-case subject is an expression
	"CREATE FUNCTION f(x int) RETURNS int BEGIN CASE WHEN x IN (1, 2) THEN RETURN 1; END CASE; END", // searched WHEN is full expression (R2)

	// LOOP / WHILE / REPEAT with labels + ITERATE / LEAVE
	"CREATE FUNCTION f(x bigint) RETURNS bigint BEGIN DECLARE i bigint DEFAULT 0; top: LOOP SET i = i + 1; IF i > 10 THEN LEAVE top; END IF; END LOOP; RETURN i; END",
	"CREATE FUNCTION f(n bigint) RETURNS bigint BEGIN DECLARE i bigint DEFAULT 0; WHILE i < n DO SET i = i + 1; END WHILE; RETURN i; END",
	"CREATE FUNCTION f(n bigint) RETURNS bigint BEGIN DECLARE i bigint DEFAULT 0; REPEAT SET i = i + 1; UNTIL i >= n END REPEAT; RETURN i; END",
	"CREATE FUNCTION f() RETURNS bigint BEGIN DECLARE i bigint DEFAULT 0; abc: WHILE i < 10 DO SET i = i + 1; ITERATE abc; END WHILE; RETURN i; END",
	"CREATE FUNCTION f() RETURNS bigint BEGIN LOOP RETURN 1; END LOOP; END", // unlabeled LOOP
	"CREATE FUNCTION f() RETURNS bigint BEGIN WHILE true DO RETURN 1; END WHILE; RETURN 0; END",
	"CREATE FUNCTION f(n int) RETURNS int BEGIN DECLARE i int DEFAULT 0; REPEAT SET i = i + 1; UNTIL i >= n AND i > 0 END REPEAT; RETURN i; END", // UNTIL is full expression (R2)
	"CREATE FUNCTION f() RETURNS int BEGIN DECLARE i int DEFAULT 0; lbl: REPEAT SET i = i + 1; UNTIL i > 3 END REPEAT; RETURN i; END",            // labeled REPEAT
	// Label names that ARE non-reserved control keywords: the label check must
	// precede the keyword dispatch (Trino accepts `loop:` / `while:` / `set:`).
	"CREATE FUNCTION f() RETURNS int BEGIN loop: LOOP LEAVE loop; END LOOP; RETURN 1; END",
	"CREATE FUNCTION f() RETURNS int BEGIN while: LOOP LEAVE while; END LOOP; RETURN 1; END",
	"CREATE FUNCTION f() RETURNS int BEGIN set: WHILE true DO LEAVE set; END WHILE; RETURN 1; END",

	// ---------------------------------------------------------------------
	// DROP FUNCTION (truth1 drop-function, F5)
	// ---------------------------------------------------------------------
	"DROP FUNCTION f(bigint)",
	"DROP FUNCTION IF EXISTS my.func(integer, varchar)",
	"DROP FUNCTION meaning_of_life()",
	"DROP FUNCTION example.default.f(bigint)",
	"DROP FUNCTION f(x bigint)", // named parameter in drop (allowed)

	// ---------------------------------------------------------------------
	// WITH FUNCTION inline routines (rootQuery, F6)
	// ---------------------------------------------------------------------
	"WITH FUNCTION f(x bigint) RETURNS bigint RETURN x + 1 SELECT f(10)",
	"WITH FUNCTION a(x int) RETURNS int RETURN x, FUNCTION b(y int) RETURNS int RETURN y SELECT a(1), b(2)",
	"WITH FUNCTION f(x int) RETURNS int RETURN x WITH t AS (SELECT 1 AS c) SELECT f(c) FROM t",
	"WITH FUNCTION f(x int) RETURNS int BEGIN RETURN x; END SELECT f(1)",
	// EXPLAIN wraps any statement (it recurses into parseStmt), so it gains
	// CREATE FUNCTION and WITH FUNCTION automatically.
	"EXPLAIN CREATE FUNCTION f() RETURNS int RETURN 1",
	"EXPLAIN WITH FUNCTION f(x int) RETURNS int RETURN x SELECT f(1)",
	// deeply nested control flow (IF within IF within a compound)
	"CREATE FUNCTION f(x int) RETURNS int BEGIN IF x > 0 THEN IF x > 5 THEN RETURN 2; END IF; RETURN 1; END IF; RETURN 0; END",

	// ---------------------------------------------------------------------
	// NEGATIVES — grammar distinctions the parser MUST reject
	// ---------------------------------------------------------------------
	// R1: RETURN / DECLARE DEFAULT take a valueExpression, not a boolean/predicate.
	"CREATE FUNCTION f(x int) RETURNS boolean RETURN x IN (1, 2, 3)",
	"CREATE FUNCTION f(x int) RETURNS boolean RETURN x > 0 AND x < 10",
	"CREATE FUNCTION f(x int) RETURNS boolean RETURN x BETWEEN 1 AND 10",
	"CREATE FUNCTION f() RETURNS int BEGIN DECLARE x boolean DEFAULT 1 IN (1, 2); RETURN 1; END",
	"CREATE FUNCTION f() RETURNS boolean BEGIN DECLARE x boolean DEFAULT true AND false; RETURN x; END",
	// missing pieces
	"CREATE FUNCTION f() RETURNS bigint RETURN",                                     // RETURN needs an expression
	"CREATE FUNCTION f() RETURN 1",                                                  // missing RETURNS type
	"CREATE FUNCTION f() RETURNS bigint",                                            // missing controlStatement body
	"CREATE FUNCTION RETURNS bigint RETURN 1",                                       // missing function name
	"CREATE FUNCTION f() RETURNS bigint BEGIN RETURN 1 END",                         // R7: missing ';' after RETURN
	"CREATE FUNCTION f() RETURNS bigint BEGIN IF true THEN RETURN 1; END; END",      // END IF required, not END
	"CREATE FUNCTION f() RETURNS int BEGIN SET x = 1; DECLARE y int; RETURN y; END", // R6: DECLARE must precede
	// R3: assignment target is a single identifier
	"CREATE FUNCTION f() RETURNS int BEGIN DECLARE x int; SET x.y = 1; RETURN x; END",
	// R5: labels only on LOOP / WHILE / REPEAT, not compound
	"CREATE FUNCTION f() RETURNS int BEGIN lbl: BEGIN RETURN 1; END; END",
	// R4: ITERATE / LEAVE require a label
	"CREATE FUNCTION f() RETURNS bigint BEGIN LOOP LEAVE; END LOOP; RETURN 1; END",
	"CREATE FUNCTION f() RETURNS bigint BEGIN LOOP ITERATE; END LOOP; RETURN 1; END",
	// commentCharacteristic is `COMMENT string_` — a non-string is rejected.
	"CREATE FUNCTION f() RETURNS int COMMENT 42 RETURN 1",
	"CREATE FUNCTION f() RETURNS int COMMENT foo RETURN 1",
	// languageCharacteristic is `LANGUAGE identifier` — a string is rejected.
	"CREATE FUNCTION f() RETURNS int LANGUAGE 'sql' RETURN 1",
	// securityCharacteristic value is the closed set { DEFINER, INVOKER }.
	"CREATE FUNCTION f() RETURNS int SECURITY PUBLIC RETURN 1",
	// no standalone functionSpecification at statement level
	"FUNCTION f() RETURNS int RETURN 1",
	// DROP FUNCTION needs the (possibly empty) parameter parens / a complete IF EXISTS.
	"DROP FUNCTION f",
	"DROP FUNCTION IF f(int)",
	// WITH FUNCTION must be followed by a query; a dangling comma is invalid.
	"WITH FUNCTION f(x int) RETURNS int RETURN x",
	"WITH FUNCTION f(x int) RETURNS int RETURN x, SELECT 1",
	// F4 negative: empty WITH () rejected
	"CREATE FUNCTION f() RETURNS int WITH () RETURN 1",
	// WITH FUNCTION is rootQuery-only. A relation/FROM subquery is parsed via
	// parseQuery (relation.go), which rejects a leading WITH FUNCTION — matching
	// Trino. (A scalar subquery in the SELECT list, e.g.
	// `SELECT (WITH FUNCTION … SELECT …)`, is NOT in this corpus: the expressions
	// node captures expression-embedded subqueries as raw-text placeholders that
	// it does not re-validate (expr.go B1), so omni accepts ANY content there
	// — including syntactically invalid queries like
	// `SELECT (SELECT a FROM t WHERE)`. That over-acceptance is a pre-existing
	// limitation of the expression layer, not of this node, so adjudicating it
	// here would be testing the wrong node.)
	"SELECT * FROM (WITH FUNCTION f(x int) RETURNS int RETURN x SELECT f(1) AS c) t",
}

// TestRoutines_OracleDifferential is the authoritative accept/reject gate: omni's
// Parse verdict must equal Trino 481's verdict for every corpus statement.
func TestRoutines_OracleDifferential(t *testing.T) {
	if testing.Short() {
		t.Skip("trino oracle: skipped in -short mode")
	}
	o := connectOracle(t)
	for _, sql := range routinesOracleCorpus {
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
