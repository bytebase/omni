package parser

import (
	"testing"
)

// This file is the parser-dcl-tcl node's differential-oracle gate
// (correctness-protocol.md): for every statement in the corpus below, omni's
// Parse accept/reject verdict must equal Trino 481's verdict as reported by the
// live oracle (trinooracle.CheckSyntax, which classifies a SYNTAX_ERROR as
// reject and every other outcome — success or a semantic error like
// ROLE_NOT_FOUND / NOT_FOUND / CATALOG_NOT_FOUND — as accept).
//
// The corpus is NOT split into "accept" and "reject" lists with hardcoded
// expectations: the oracle is the source of truth, so each case is run through
// both omni and Trino and the two verdicts are compared. This both proves
// correctness and documents (via the captured verdict) which docs/legacy forms
// Trino 481 actually accepts. Every documented form (truth1), every legacy
// grammar alternative (truth2), the oracle-confirmed docs-ahead extensions
// (ON BRANCH, bare ALL, open privilege vocab, DROP ROLE IF EXISTS), and a set
// of negative forms are present.
//
// Skipped cleanly when no Trino oracle is reachable.

// dclTclOracleCorpus is the full DCL/TCL/prepared statement corpus for the
// differential gate. Each string is fed to both omni and the oracle; the
// verdicts must match.
var dclTclOracleCorpus = []string{
	// --- CREATE / DROP ROLE (truth1 + truth2) ---
	"CREATE ROLE admin",
	"CREATE ROLE moderator WITH ADMIN USER bob",
	"CREATE ROLE admin IN hive",
	"CREATE ROLE r WITH ADMIN ROLE other",
	"CREATE ROLE r WITH ADMIN CURRENT_USER",
	"CREATE ROLE r WITH ADMIN CURRENT_ROLE IN hive",
	"DROP ROLE admin",
	"DROP ROLE IF EXISTS analyst", // D4 docs-ahead
	"DROP ROLE IF EXISTS analyst IN hive",
	"DROP ROLE analyst IN hive",
	// negatives
	"CREATE ROLE",                     // missing name
	"DROP ROLE",                       // missing name
	"CREATE ROLE WITH ADMIN USER bob", // missing name before WITH
	"DROP ROLE IF analyst",            // malformed IF (no EXISTS)

	// --- GRANT roles (truth1 + truth2) ---
	"GRANT bar TO USER foo",
	"GRANT bar, foo TO USER baz, ROLE qux WITH ADMIN OPTION",
	"GRANT bar TO foo",
	"GRANT bar TO ROLE qux",
	"GRANT bar TO USER foo GRANTED BY CURRENT_USER",
	"GRANT bar TO USER foo GRANTED BY ROLE admin",
	"GRANT bar TO USER foo GRANTED BY USER carol",
	"GRANT bar TO USER foo IN hive",
	"GRANT bar TO USER foo WITH ADMIN OPTION GRANTED BY CURRENT_USER IN hive",
	"GRANT zone TO USER foo", // non-reserved keyword as role name
	"GRANT ALL TO foo",       // ALL as a role name (followed by TO)
	`GRANT "my role" TO USER foo`,
	// USER / ROLE are non-reserved: a bare USER/ROLE not followed by a name is
	// itself the principal name.
	"GRANT r TO USER",
	"GRANT r TO ROLE",
	"GRANT r TO ROLE, bar",
	"GRANT r TO USER WITH ADMIN OPTION",
	"GRANT r TO foo GRANTED BY USER", // bare USER grantor
	// negatives
	"GRANT SELECT TO foo",                // reserved privilege keyword, no ON
	"GRANT select TO foo",                // reserved word can't be a role
	"GRANT CREATE TO foo",                // reserved CREATE can't be a role
	"GRANT SELECT, INSERT TO foo",        // reserved list, no ON
	"GRANT bar TO foo WITH GRANT OPTION", // WITH GRANT OPTION is privilege-only
	"GRANT bar TO",                       // missing principal

	// --- GRANT privileges (truth1 + truth2 + docs-ahead) ---
	"GRANT INSERT, SELECT ON orders TO alice",
	"GRANT DELETE ON SCHEMA finance TO bob",
	"GRANT SELECT ON nation TO alice WITH GRANT OPTION",
	"GRANT SELECT ON orders TO ROLE PUBLIC",
	"GRANT SELECT ON TABLE orders TO alice",
	"GRANT ALL PRIVILEGES ON test TO alice",
	"GRANT ALL ON test TO alice", // D3 bare ALL
	"GRANT UPDATE ON orders TO USER alice",
	"GRANT SELECT ON cat.sch.tbl TO alice",
	"GRANT INSERT ON BRANCH audit IN orders TO alice",       // D2 docs-ahead
	"GRANT INSERT ON BRANCH audit IN TABLE orders TO alice", // branch + IN entity-kind
	"GRANT INSERT ON BRANCH audit IN cat.s.t TO alice",      // branch + dotted IN target
	"GRANT mypriv ON test TO alice",                         // D1 open vocab
	`GRANT "My Priv" ON test TO alice`,
	"GRANT TO ON t TO alice",   // privilege literally named "TO" (non-reserved)
	"GRANT FROM ON t TO alice", // privilege literally named "FROM"
	// D5 open entity-kind qualifier + branch-vs-name disambiguation
	"GRANT SELECT ON VIEW v TO alice",        // arbitrary entity kind
	"GRANT SELECT ON a b TO alice",           // arbitrary qualifier "a", object "b"
	"GRANT SELECT ON branch TO alice",        // object literally named "branch"
	"GRANT SELECT ON branch orders TO alice", // qualifier "branch", object "orders"
	"GRANT SELECT ON BRANCH orders TO alice", // qualifier "BRANCH", object "orders" (no IN)
	"GRANT SELECT ON x a.b.c TO alice",       // qualifier + dotted object
	// negatives
	"GRANT ALL PRIVILEGES TO foo",                           // privilege spec needs ON
	"GRANT SELECT ON test TO alice IN hive",                 // IN catalog is role-only
	"GRANT SELECT ON test TO alice WITH ADMIN OPTION",       // WITH ADMIN OPTION is role-only
	"GRANT SELECT ON test TO alice GRANTED BY CURRENT_USER", // GRANTED BY is role-only
	"GRANT CREATE TABLE ON SCHEMA s TO alice",               // privileges are single tokens
	"GRANT SELECT ON TO alice",                              // missing object name
	"GRANT SELECT orders TO alice",                          // missing ON
	"GRANT SELECT ON a b c TO alice",                        // at most one entity-kind word
	"GRANT SELECT ON a.b c TO alice",                        // entity-kind qualifier must be single ident
	"GRANT SELECT ON SELECT orders TO alice",                // reserved word can't be entity kind
	"GRANT SELECT ON BRANCH a.b IN orders TO alice",         // branch name must be single ident
	"GRANT SELECT ON BRANCH IN orders TO alice",             // branch name required before IN

	// --- REVOKE roles ---
	"REVOKE bar FROM USER foo",
	"REVOKE ADMIN OPTION FOR bar, foo FROM USER baz, ROLE qux",
	"REVOKE bar FROM foo GRANTED BY CURRENT_ROLE IN hive",
	"REVOKE bar, foo FROM USER baz",
	// negatives
	"REVOKE select FROM foo", // reserved word can't be a role
	"REVOKE bar FROM",        // missing principal

	// --- REVOKE privileges ---
	"REVOKE INSERT, SELECT ON orders FROM alice",
	"REVOKE DELETE ON SCHEMA finance FROM bob",
	"REVOKE GRANT OPTION FOR SELECT ON nation FROM ROLE PUBLIC",
	"REVOKE ALL PRIVILEGES ON test FROM alice",
	"REVOKE INSERT ON BRANCH audit IN orders FROM alice", // D2
	"REVOKE ALL ON test FROM alice",                      // D3 bare ALL
	// negatives
	"REVOKE SELECT ON test FROM alice CASCADE",  // Trino has no CASCADE
	"REVOKE SELECT ON test FROM alice RESTRICT", // nor RESTRICT
	"REVOKE GRANT OPTION FOR bar FROM foo",      // GRANT OPTION FOR with no ON (role-ish)

	// --- DENY ---
	"DENY INSERT, SELECT ON orders TO alice",
	"DENY DELETE ON SCHEMA finance TO bob",
	"DENY SELECT ON orders TO ROLE PUBLIC",
	"DENY ALL PRIVILEGES ON test TO alice",
	"DENY INSERT ON BRANCH audit IN orders TO alice", // D2
	// negatives
	"DENY bar TO foo",                  // DENY has no role form (no ON)
	"DENY SELECT ON orders FROM alice", // DENY uses TO, not FROM

	// --- START TRANSACTION / COMMIT / ROLLBACK ---
	"START TRANSACTION",
	"START TRANSACTION ISOLATION LEVEL REPEATABLE READ",
	"START TRANSACTION READ WRITE",
	"START TRANSACTION READ ONLY",
	"START TRANSACTION ISOLATION LEVEL READ COMMITTED, READ ONLY",
	"START TRANSACTION READ WRITE, ISOLATION LEVEL SERIALIZABLE",
	"START TRANSACTION ISOLATION LEVEL READ UNCOMMITTED",
	"COMMIT",
	"COMMIT WORK",
	"ROLLBACK",
	"ROLLBACK WORK",
	// negatives
	"START TRANSACTION ISOLATION LEVEL BOGUS",
	"START TRANSACTION READ",
	"START TRANSACTION ISOLATION REPEATABLE READ", // missing LEVEL
	"START TRANSACTION ISOLATION LEVEL READ COMMITTED,",
	"COMMIT TRANSACTION", // Trino COMMIT has no TRANSACTION keyword
	"ROLLBACK TO sp",     // Trino ROLLBACK has no savepoint form

	// --- PREPARE / DEALLOCATE / EXECUTE / DESCRIBE INPUT|OUTPUT ---
	"PREPARE my_select1 FROM SELECT * FROM nation",
	"PREPARE my_select2 FROM SELECT name FROM nation WHERE regionkey = ? AND nationkey < ?",
	"PREPARE my_insert FROM INSERT INTO cities VALUES (1, 'San Francisco')",
	"PREPARE my_create FROM CREATE TABLE foo AS SELECT * FROM nation",
	"DEALLOCATE PREPARE my_query",
	"EXECUTE my_select1",
	"EXECUTE my_select2 USING 1, 3",
	"EXECUTE IMMEDIATE 'SELECT name FROM nation'",
	"EXECUTE IMMEDIATE 'SELECT name FROM nation WHERE regionkey = ? AND nationkey < ?' USING 1, 3",
	"EXECUTE IMMEDIATE",         // IMMEDIATE is non-reserved: named-execute of "immediate"
	"EXECUTE IMMEDIATE USING 1", // named-execute "immediate" with USING
	"DESCRIBE INPUT my_select1",
	"DESCRIBE OUTPUT my_select1",
	"DESCRIBE INPUT foo.bar",  // describeInput takes a single identifier...
	"DESCRIBE OUTPUT foo.bar", // ...so a dotted name is a SYNTAX_ERROR
	// negatives
	"PREPARE p FROM",               // missing inner statement
	"PREPARE FROM SELECT 1",        // missing name
	"PREPARE p FROM 1",             // inner is not a statement
	"PREPARE p FROM NOTASTATEMENT", // inner is not a statement
	"DEALLOCATE my_query",          // missing PREPARE keyword
	"EXECUTE IMMEDIATE my_name",    // "immediate" name + trailing token (rejected)
	"EXECUTE my_select USING",      // dangling USING
	// NOTE: bare `DESCRIBE INPUT` / `DESCRIBE INPUT.col` are the SHOW COLUMNS
	// alias `DESCRIBE qualifiedName` (INPUT/OUTPUT are non-reserved table
	// names there), owned by the parser-utility node — out of this node's
	// scope, so they are not in this corpus.
}

// TestDCLTCL_OracleDifferential is the authoritative accept/reject gate: omni's
// Parse verdict must equal Trino 481's verdict for every corpus statement.
func TestDCLTCL_OracleDifferential(t *testing.T) {
	if testing.Short() {
		t.Skip("trino oracle: skipped in -short mode")
	}
	o := connectOracle(t)
	for _, sql := range dclTclOracleCorpus {
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
