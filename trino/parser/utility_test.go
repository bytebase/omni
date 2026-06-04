package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/trino/ast"
)

// This file is the correctness gate for the `parser-utility` DAG node
// (show.go + session.go + explain_call.go). Following correctness-protocol.md:
//
//   - Layer 1 (accept/reject differential): utilityAcceptCorpus and
//     utilityRejectCorpus span every SHOW / DESCRIBE / USE / SET* / RESET* /
//     EXPLAIN / CALL alternative of TrinoParser.g4, plus the oracle-discovered
//     docs-ahead forms (SHOW CREATE FUNCTION) and oracle-discovered negatives.
//     TestUtility_OracleDifferential is the authoritative gate: this node's
//     accept/reject (via parseUtility) must equal Trino 481's accept/reject of
//     the same statement.
//   - Layer 2 (structural): TestUtility_Structure pins the parse-node shape of
//     representative forms (Kind, Name, scope, options) so a form that "accepts"
//     is also parsed into the right node.
//   - Completeness: every legacy alternative appears in the corpus with its
//     oracle-derived verdict; the per-alternative checklist is the corpus
//     comment blocks below.
//
// The oracle helpers (connectOracle / oracleAccepts / truncateName) live in
// oracle_foundation_test.go and lexer_oracle_test.go.
//
// Scope boundary: DESCRIBE INPUT / DESCRIBE OUTPUT (prepared-statement
// introspection) and START TRANSACTION / COMMIT / ROLLBACK / PREPARE / EXECUTE /
// DEALLOCATE / GRANT / REVOKE / DENY belong to the parser-dcl-tcl node and are
// NOT exercised here — they are still stubbed (unsupported) until that node
// lands. parseUtility asserts the statement it parsed is one of this node's
// kinds, so an accidentally-claimed out-of-scope form would surface.

// utilityVerdict is this node's tri-state classification of a parse result,
// needed because EXPLAIN recurses into the inner statement (E1) whose parser may
// not have landed yet:
//
//   - verdictAccept: parsed into one of this node's node types with no error.
//   - verdictReject: a real syntax rejection (or a nil/empty/foreign result).
//   - verdictPendingInner: the statement IS one of this node's nodes, but its
//     only error is a downstream "not yet supported" stub — exclusively the
//     EXPLAIN-wraps-an-unimplemented-inner-statement case. Trino accepts these;
//     omni will too once the inner DAG node (parser-select / dml / ddl) lands.
//     The differential treats these as pending, not failures, so this node does
//     not falsely depend on a sibling/later node.
type utilityVerdict int

const (
	verdictReject utilityVerdict = iota
	verdictAccept
	verdictPendingInner
)

// isUtilityNode reports whether n is one of the node types this node produces.
func isUtilityNode(n ast.Node) bool {
	switch n.(type) {
	case *ShowStmt, *UseStmt, *SetSessionStmt, *ResetSessionStmt,
		*SetSessionAuthorizationStmt, *ResetSessionAuthorizationStmt,
		*SetRoleStmt, *SetPathStmt, *SetTimeZoneStmt, *ExplainStmt, *CallStmt:
		return true
	default:
		return false
	}
}

// classifyUtility parses a single utility statement and returns its verdict plus
// the parse errors (for diagnostics). See utilityVerdict.
func classifyUtility(sql string) (utilityVerdict, ast.Node, []ParseError) {
	file, errs := Parse(sql)
	var stmt ast.Node
	if file != nil && len(file.Stmts) == 1 {
		stmt = file.Stmts[0]
	}
	if len(errs) == 0 {
		if stmt != nil && isUtilityNode(stmt) {
			return verdictAccept, stmt, nil
		}
		return verdictReject, nil, nil
	}
	// Errors present. The single pending case: an EXPLAIN node was produced and
	// every error is a "not yet supported" inner-statement stub.
	if _, ok := stmt.(*ExplainStmt); ok && allNotYetSupported(errs) {
		return verdictPendingInner, stmt, errs
	}
	return verdictReject, nil, errs
}

// allNotYetSupported reports whether every error is a downstream "not yet
// supported" stub (the foundation's unsupported() sentinel), as opposed to a
// real syntax error.
func allNotYetSupported(errs []ParseError) bool {
	if len(errs) == 0 {
		return false
	}
	for _, e := range errs {
		if !strings.Contains(e.Msg, "not yet supported") {
			return false
		}
	}
	return true
}

// parseUtility reports whether this node cleanly accepts sql (verdictAccept).
// Pending-inner and reject both count as "not a clean accept" for the oracle-
// free smoke tests; the oracle differential uses classifyUtility directly to
// treat pending-inner specially.
func parseUtility(sql string) (ast.Node, bool) {
	v, node, _ := classifyUtility(sql)
	return node, v == verdictAccept
}

// utilityAcceptCorpus is every utility statement this node must ACCEPT,
// organised by the legacy alternative it exercises. Every entry has been
// confirmed ACCEPTED by the live Trino 481 oracle (a non-SYNTAX_ERROR verdict).
var utilityAcceptCorpus = []string{
	// --- use ---
	"USE information_schema",
	"USE hive.finance",
	"USE \"My Catalog\".\"My Schema\"",

	// --- setSession / resetSession ---
	"SET SESSION query_max_run_time = '10m'",
	"SET SESSION example.incremental_refresh_enabled = false",
	"SET SESSION foo = 42",
	"SET SESSION cat.prop = ARRAY[1, 2]",
	// SET/RESET SESSION names are NOT part-limited at parse time (semantic
	// INVALID_SESSION_PROPERTY, not SYNTAX_ERROR) — must not over-restrict.
	"SET SESSION a.b.c.d = 1",
	// AUTHORIZATION as a property NAME (not the auth keyword) when followed by
	// '.' or '=' — Trino 481 accepts these as property assignments/resets.
	"SET SESSION AUTHORIZATION = 1",
	"SET SESSION AUTHORIZATION.foo = 1",
	"RESET SESSION query_max_run_time",
	"RESET SESSION hive.optimized_reader_enabled",
	"RESET SESSION a.b.c.d",
	"RESET SESSION AUTHORIZATION.foo",

	// --- setSessionAuthorization / resetSessionAuthorization ---
	"SET SESSION AUTHORIZATION 'John'",
	"SET SESSION AUTHORIZATION \"John\"",
	"SET SESSION AUTHORIZATION John",
	"RESET SESSION AUTHORIZATION",

	// --- setRole ---
	"SET ROLE admin",
	"SET ROLE ALL",
	"SET ROLE NONE",
	"SET ROLE analyst IN hive",

	// --- setPath ---
	"SET PATH example.system",
	"SET PATH a, b.c, d",

	// --- setTimeZone ---
	"SET TIME ZONE LOCAL",
	"SET TIME ZONE '-08:00'",
	"SET TIME ZONE INTERVAL '10' HOUR",
	"SET TIME ZONE INTERVAL -'08:00' HOUR TO MINUTE",
	"SET TIME ZONE 'America/Los_Angeles'",
	"SET TIME ZONE concat_ws('/', 'America', 'Los_Angeles')",

	// --- showGrants ---
	"SHOW GRANTS",
	"SHOW GRANTS ON TABLE orders",
	"SHOW GRANTS ON orders",

	// --- showCreateTable / Schema / View / MaterializedView / Function ---
	"SHOW CREATE TABLE sf1.orders",
	"SHOW CREATE SCHEMA hive.web",
	"SHOW CREATE VIEW my_view",
	"SHOW CREATE MATERIALIZED VIEW my_mv",
	// U1: docs-ahead-of-legacy — accepted by Trino 481 (NOT_SUPPORTED, semantic).
	"SHOW CREATE FUNCTION example.default.meaning_of_life",
	"SHOW CREATE FUNCTION f",

	// --- showTables ---
	"SHOW TABLES",
	"SHOW TABLES FROM tpch.tiny",
	"SHOW TABLES IN tpch.tiny",
	"SHOW TABLES FROM tpch.tiny LIKE 'p%'",
	"SHOW TABLES LIKE 'p%' ESCAPE '#'",

	// --- showSchemas (FROM/IN takes a single-identifier catalog) ---
	"SHOW SCHEMAS",
	"SHOW SCHEMAS FROM tpch",
	"SHOW SCHEMAS IN tpch",
	"SHOW SCHEMAS FROM tpch LIKE '__3%'",

	// --- showCatalogs ---
	"SHOW CATALOGS",
	"SHOW CATALOGS LIKE 't%'",

	// --- showColumns (FROM/IN + name both mandatory — U2) ---
	"SHOW COLUMNS FROM nation",
	"SHOW COLUMNS IN nation",
	"SHOW COLUMNS FROM tpch.sf1.nation",
	"SHOW COLUMNS FROM nation LIKE '%key'",

	// --- showStats / showStatsForQuery ---
	"SHOW STATS FOR nation",
	"SHOW STATS FOR (SELECT * FROM nation WHERE regionkey = 1)",

	// --- showRoles (CURRENT?, FROM/IN single identifier) ---
	"SHOW ROLES",
	"SHOW ROLES FROM hive",
	"SHOW CURRENT ROLES",
	"SHOW CURRENT ROLES FROM hive",

	// --- showRoleGrants ---
	"SHOW ROLE GRANTS",
	"SHOW ROLE GRANTS FROM hive",

	// --- showFunctions ---
	"SHOW FUNCTIONS",
	"SHOW FUNCTIONS FROM example.default",
	"SHOW FUNCTIONS LIKE 'array%'",

	// --- showSession ---
	"SHOW SESSION",
	"SHOW SESSION LIKE 'query%'",

	// --- describe / desc (== showColumns) ---
	"DESCRIBE nation",
	"DESC nation",
	"DESCRIBE tpch.sf1.nation",
	// INPUT/OUTPUT as a TABLE name (not the prepared form): no trailing name.
	"DESCRIBE input",
	"DESCRIBE output",
	"DESC input",

	// --- call (E3: 1-3 part name, empty/positional/named/mixed args) ---
	"CALL test(123, 'apple')",
	"CALL test(name => 'apple', id => 123)",
	"CALL p()",
	"CALL s.p()",
	"CALL c.s.p()",
	"CALL system.runtime.kill_query(query_id => '20210101_000000_00000_abcde')",
	"CALL f(name => 1, 2)",
}

// utilityRejectCorpus is every malformed utility statement this node must
// REJECT. Every entry has been confirmed a SYNTAX_ERROR by the live Trino 481
// oracle. These are the required negative coverage.
var utilityRejectCorpus = []string{
	// truncated heads
	"SHOW",
	"SHOW TABLES FROM",
	"SHOW COLUMNS",      // U2: name required
	"SHOW COLUMNS FROM", // U2: name required after FROM
	"SHOW STATS FOR",
	"SHOW STATS FOR (nation)",              // parens must hold a query, not a name
	"SHOW STATS FOR (garbage tokens here)", // not a query
	"SHOW STATS FOR ()",                    // empty parens
	"SHOW SESSION LIKE",                    // LIKE needs a pattern
	"SHOW CREATE TABLE",                    // name required
	"SHOW CREATE",                          // object kind required
	"DESCRIBE",
	"USE",
	"SET",
	"SET SESSION = 'x'", // name required
	"SET SESSION foo",   // '=' value required
	"SET TIME ZONE",
	"SET ROLE",
	"SET PATH",
	"RESET",
	"RESET SESSION", // name or AUTHORIZATION required
	"CALL",
	"CALL f", // parens mandatory
	// SHOW ROLE/ROLES disambiguation edge cases
	"SHOW ROLE",                // ROLE needs GRANTS
	"SHOW CURRENT ROLE GRANTS", // CURRENT only precedes ROLES, not ROLE GRANTS
	// dotted-scope where a single identifier is required
	"SHOW SCHEMAS FROM c.x",
	// procedure name too long (E3)
	"CALL c.s.p.x()",
	// qualifiedName part limits (oracle enforces at parse time; legacy grammar
	// is unbounded — confirmed divergence). Object names cap at 3, schema scopes
	// at 2, path elements at 2.
	"SHOW CREATE TABLE a.b.c.d",
	"SHOW CREATE FUNCTION a.b.c.d",
	"SHOW COLUMNS FROM a.b.c.d",
	"SHOW STATS FOR a.b.c.d",
	"SHOW GRANTS ON TABLE a.b.c.d",
	"DESCRIBE a.b.c.d",
	"SHOW TABLES FROM a.b.c",
	"SHOW FUNCTIONS FROM a.b.c",
	"SET PATH a.b.c",
	// DESC has no INPUT/OUTPUT prepared form, and two bare names is invalid
	"DESC INPUT a",
	// trailing tokens after a complete statement (parseSingle trailing-token
	// check — Trino rejects all of these too)
	"SHOW TABLES extra",
	"USE a b",
	"SET ROLE ALL x",
	"SHOW STATS FOR (SELECT 1) x",
	"DESCRIBE a b",
	"SET TIME ZONE LOCAL x",
	"CALL p() x",
	// EXPLAIN structure errors (E1/E2): the inner statement is reached only
	// after EXPLAIN's own prefix parses, so these test EXPLAIN itself.
	"EXPLAIN ()",              // empty option list
	"EXPLAIN (TYPE) SELECT 1", // option value required
	"EXPLAIN (TYPE FOO) SELECT 1",
	"EXPLAIN (FORMAT BAR) SELECT 1",
	"EXPLAIN x",       // inner `x` is not a statement — Trino rejects too
	"EXPLAIN ANALYZE", // inner statement required
}

// explainCleanAcceptCorpus is EXPLAIN forms whose inner statement is already
// implemented (a SHOW / SET / USE / CALL — this node's own statements), so the
// whole thing parses cleanly (verdictAccept). Trino accepts all of these.
var explainCleanAcceptCorpus = []string{
	"EXPLAIN SHOW TABLES",
	"EXPLAIN SHOW SCHEMAS FROM tpch",
	"EXPLAIN ANALYZE SHOW TABLES",
	"EXPLAIN (TYPE LOGICAL) SHOW TABLES",
}

// explainPendingAcceptCorpus is EXPLAIN forms whose inner statement is NOT yet
// implemented (SELECT / INSERT / CREATE … — parser-select / dml / ddl DAG
// nodes). Trino accepts them; this node parses EXPLAIN's own structure and
// recurses, so it currently yields verdictPendingInner (the only error is the
// inner "not yet supported" stub). They become verdictAccept once the inner
// node lands. The differential asserts: Trino accepts AND omni is pending-or-
// accept (never a real syntax reject) — proving EXPLAIN's own grammar is right
// without depending on a sibling/later node.
var explainPendingAcceptCorpus = []string{
	"EXPLAIN SELECT 1",
	"EXPLAIN (TYPE LOGICAL) SELECT regionkey FROM nation",
	"EXPLAIN (TYPE LOGICAL, FORMAT JSON) SELECT 1",
	"EXPLAIN (TYPE DISTRIBUTED) SELECT 1",
	"EXPLAIN (TYPE VALIDATE) SELECT 1",
	"EXPLAIN (TYPE IO, FORMAT JSON) SELECT 1",
	"EXPLAIN (FORMAT TEXT) SELECT 1",
	"EXPLAIN (FORMAT GRAPHVIZ) SELECT 1",
	"EXPLAIN ANALYZE SELECT 1",
	"EXPLAIN ANALYZE VERBOSE SELECT 1",
	"EXPLAIN INSERT INTO t VALUES (1)",
	"EXPLAIN CREATE TABLE t AS SELECT 1 AS x",
}

// TestUtility_AcceptCorpusParses verifies this node parses every form in
// utilityAcceptCorpus into one of its node types with no error (oracle-free
// completeness smoke test). The oracle differential proves the verdicts match
// Trino.
func TestUtility_AcceptCorpusParses(t *testing.T) {
	accept := append(append([]string{}, utilityAcceptCorpus...), explainCleanAcceptCorpus...)
	for _, sql := range accept {
		t.Run(truncateName(sql), func(t *testing.T) {
			node, ok := parseUtility(sql)
			if !ok {
				_, errs := Parse(sql)
				t.Errorf("parseUtility(%q) should accept, got errors: %v", sql, errs)
			}
			if ok && node == nil {
				t.Errorf("parseUtility(%q) accepted but returned nil node", sql)
			}
		})
	}
}

// TestUtility_ExplainPendingInner verifies the EXPLAIN-wraps-an-unimplemented-
// inner-statement cases (E1): EXPLAIN's own grammar parses, and the result is
// pending-inner (the only error is the inner "not yet supported" stub) — never
// a real syntax rejection. This proves EXPLAIN's structure is correct without
// depending on the inner DAG node (parser-select / dml / ddl). Once those land,
// these become verdictAccept (the differential below tolerates both).
func TestUtility_ExplainPendingInner(t *testing.T) {
	for _, sql := range explainPendingAcceptCorpus {
		t.Run(truncateName(sql), func(t *testing.T) {
			v, _, errs := classifyUtility(sql)
			if v == verdictReject {
				t.Errorf("classifyUtility(%q) = reject, want accept or pending-inner (errs=%v)", sql, errs)
			}
		})
	}
}

// TestUtility_RejectCorpusRejected verifies this node rejects every form in
// utilityRejectCorpus (required negative coverage).
func TestUtility_RejectCorpusRejected(t *testing.T) {
	for _, sql := range utilityRejectCorpus {
		t.Run(truncateName(sql), func(t *testing.T) {
			if _, ok := parseUtility(sql); ok {
				t.Errorf("parseUtility(%q) should reject, but accepted", sql)
			}
		})
	}
}

// TestUtility_OracleDifferential is the authoritative accept/reject gate: for
// every form in both corpora, this node's accept/reject must equal Trino 481's
// accept/reject of the same statement. Skipped when no oracle is reachable.
func TestUtility_OracleDifferential(t *testing.T) {
	if testing.Short() {
		t.Skip("trino oracle: skipped in -short mode")
	}
	o := connectOracle(t)

	check := func(t *testing.T, sql string) {
		_, omniAccepts := parseUtility(sql)

		trinoAccepts, ok := oracleAccepts(t, o, sql)
		if !ok {
			t.Skip("oracle unreachable for this case")
		}
		if omniAccepts != trinoAccepts {
			_, errs := Parse(sql)
			t.Errorf("MISMATCH sql=%q: omni accepts=%v (errs=%v), Trino accepts=%v",
				sql, omniAccepts, errs, trinoAccepts)
		}
	}

	t.Run("accept", func(t *testing.T) {
		accept := append(append([]string{}, utilityAcceptCorpus...), explainCleanAcceptCorpus...)
		for _, sql := range accept {
			sql := sql
			t.Run(truncateName(sql), func(t *testing.T) { check(t, sql) })
		}
	})
	t.Run("reject", func(t *testing.T) {
		for _, sql := range utilityRejectCorpus {
			sql := sql
			t.Run(truncateName(sql), func(t *testing.T) { check(t, sql) })
		}
	})
	// EXPLAIN of an unimplemented inner statement: Trino accepts; omni must be
	// pending-or-accept (never a real syntax reject). This is the dependency-
	// aware slice of the differential.
	t.Run("explain_pending", func(t *testing.T) {
		for _, sql := range explainPendingAcceptCorpus {
			sql := sql
			t.Run(truncateName(sql), func(t *testing.T) {
				v, _, errs := classifyUtility(sql)
				trinoAccepts, ok := oracleAccepts(t, o, sql)
				if !ok {
					t.Skip("oracle unreachable for this case")
				}
				if !trinoAccepts {
					t.Errorf("expected Trino to accept %q (pending-inner corpus), but it rejected", sql)
				}
				if v == verdictReject {
					t.Errorf("MISMATCH sql=%q: omni real-rejects (errs=%v), Trino accepts; EXPLAIN's own grammar is wrong", sql, errs)
				}
			})
		}
	})
}
