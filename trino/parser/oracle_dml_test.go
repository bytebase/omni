package parser

import (
	"testing"
)

// This file is the parser-dml node's differential-oracle gate
// (correctness-protocol.md): for every statement in the corpus below, omni's
// Parse accept/reject verdict must equal Trino 481's verdict as reported by the
// live oracle (trinooracle.CheckSyntax, which classifies a SYNTAX_ERROR as
// reject and every other outcome — success or a semantic error like
// TABLE_NOT_FOUND / COLUMN_NOT_FOUND / NOT_SUPPORTED — as accept).
//
// The corpus is NOT split into hardcoded accept/reject lists: the oracle is the
// source of truth, so each case is run through both omni and Trino and the two
// verdicts compared. It covers every legacy grammar example (truth2,
// examples/{insert_into,delete,update,truncate_table,whereless_update,merge}.sql
// with their assertions stripped), every documented form (truth1
// query-dml.md), the oracle-pinned grammar facts (D-INS*/D-DEL*/D-UPD*/D-MRG*),
// and a set of negative forms.
//
// EXCLUDED — the `@ branch_name` form (divergence ledger #5). Trino 481 accepts
// `INSERT/DELETE/UPDATE/MERGE INTO tbl @ branch …`, but omni REJECTS it because
// the merged lexer emits no '@' token (an "unrecognized character" lex error),
// a lexer-node concern outside this node's writes scope. Those cases are listed
// in dmlBranchDivergenceCorpus below and asserted as a KNOWN, flagged divergence
// (omni rejects, Trino accepts) rather than mixed into the parity gate, so the
// gate stays a true equality check.
//
// Skipped cleanly when no Trino oracle is reachable.

// dmlOracleCorpus is the full DML statement corpus for the differential gate.
// Each string is fed to both omni and the oracle; the verdicts must match.
var dmlOracleCorpus = []string{
	// ---------------------------------------------------------------------
	// INSERT (truth1 insert.html + truth2 insert_into.sql)
	// ---------------------------------------------------------------------
	"INSERT INTO orders SELECT * FROM new_orders",
	"INSERT INTO cities VALUES (1, 'San Francisco')",
	"INSERT INTO cities VALUES (2, 'San Jose'), (3, 'Oakland')",
	"INSERT INTO nation (nationkey, name, regionkey, comment) SELECT 1, 'a', 2, 'c'",
	"INSERT INTO nation (nationkey, name, regionkey) VALUES (1, 'a', 2)",
	`INSERT INTO a."b/c".d SELECT * FROM t`,
	`INSERT INTO a."b/c".d (c1, c2) SELECT * FROM t`,
	"INSERT INTO t (SELECT 1)", // parenthesized source query
	"INSERT INTO t VALUES (1), (2), (3)",
	"INSERT INTO t SELECT a FROM b UNION SELECT c FROM d",    // set-op source
	"INSERT INTO t SELECT * FROM a ORDER BY 1 LIMIT 5",       // query tail
	"INSERT INTO t TABLE other",                              // TABLE query source
	"INSERT INTO t WITH cte AS (SELECT 1) SELECT * FROM cte", // WITH source
	"INSERT INTO cat.sch.t VALUES (1)",                       // 3-part target
	// INSERT negatives
	"INSERT INTO t",                  // D-INS1: source query mandatory
	"INSERT INTO t (a, b)",           // D-INS1: column list without source
	"INSERT INTO a.b.c.d VALUES (1)", // D-INS2: 4-part target
	"INSERT t VALUES (1)",            // missing INTO
	"INSERT INTO",                    // missing target
	"INSERT INTO t VALUES",           // VALUES with no rows

	// ---------------------------------------------------------------------
	// TRUNCATE (truth2 truncate_table.sql)
	// ---------------------------------------------------------------------
	"TRUNCATE TABLE orders",
	"TRUNCATE TABLE a",
	"TRUNCATE TABLE a.b",
	"TRUNCATE TABLE a.b.c",
	`TRUNCATE TABLE "My Table"`,
	// TRUNCATE negatives
	"TRUNCATE orders",             // D-TR1: TABLE keyword mandatory
	"TRUNCATE TABLE a.b.c.d",      // 4-part target
	"TRUNCATE TABLE orders extra", // trailing token
	"TRUNCATE TABLE",              // missing target

	// ---------------------------------------------------------------------
	// DELETE (truth1 delete.html + truth2 delete.sql)
	// ---------------------------------------------------------------------
	"DELETE FROM orders",
	"DELETE FROM lineitem WHERE shipmode = 'AIR'",
	"DELETE FROM lineitem WHERE orderkey IN (SELECT orderkey FROM orders WHERE totalprice > 100)",
	"DELETE FROM t",
	`DELETE FROM "awesome table"`,
	"DELETE FROM t WHERE a = b",
	"DELETE FROM cat.sch.orders WHERE x = 1",
	"DELETE FROM orders WHERE a = 1 AND b > 2 OR c IS NULL",
	// DELETE negatives
	"DELETE orders",              // D-DEL1: FROM mandatory
	"DELETE FROM orders WHERE",   // dangling WHERE
	"DELETE FROM a.b.c.d",        // 4-part target
	"DELETE FROM orders garbage", // trailing token
	"DELETE FROM",                // missing target

	// ---------------------------------------------------------------------
	// UPDATE (truth1 update.html + truth2 update.sql / whereless_update.sql)
	// ---------------------------------------------------------------------
	"UPDATE purchases SET status = 'OVERDUE' WHERE ship_date IS NULL",
	"UPDATE customers SET account_manager = 'John Henry', assign_date = now()",
	"UPDATE new_hires SET manager = (SELECT e.name FROM employees e WHERE e.employee_id = new_hires.manager_id)",
	"UPDATE foo_tablen SET bar = 23, baz = 3.1415E0, bletch = 'barf' WHERE (nothing = 'fun')",
	"UPDATE foo_tablen SET bar = 23",
	`UPDATE t SET "col one" = a + b * 2, x = CASE WHEN y THEN 1 ELSE 2 END`,
	"UPDATE cat.sch.t SET a = 1",
	// UPDATE negatives
	"UPDATE t a = 1",               // D-UPD2: SET mandatory
	"UPDATE t SET",                 // D-UPD2: at least one assignment
	"UPDATE t SET (a, b) = (1, 2)", // D-UPD1: no row assignment
	"UPDATE t SET a = 1 garbage",   // trailing token
	"UPDATE SET a = 1",             // missing target
	"UPDATE t SET a",               // assignment missing = expr

	// ---------------------------------------------------------------------
	// MERGE (truth1 merge.html + truth2 merge.sql)
	// ---------------------------------------------------------------------
	"MERGE INTO accounts t USING monthly_accounts_update s ON t.customer = s.customer WHEN MATCHED THEN DELETE",
	"MERGE INTO accounts t USING monthly_accounts_update s ON (t.customer = s.customer) " +
		"WHEN MATCHED THEN UPDATE SET purchases = s.purchases + t.purchases " +
		"WHEN NOT MATCHED THEN INSERT (customer, purchases, address) VALUES (s.customer, s.purchases, s.address)",
	"MERGE INTO accounts t USING monthly_accounts_update s ON (t.customer = s.customer) " +
		"WHEN MATCHED AND s.address = 'Centreville' THEN DELETE " +
		"WHEN MATCHED THEN UPDATE SET purchases = s.purchases + t.purchases, address = s.address " +
		"WHEN NOT MATCHED THEN INSERT (customer, purchases, address) VALUES (s.customer, s.purchases, s.address)",
	"MERGE INTO inventory AS i USING changes AS c ON i.part = c.part " +
		"WHEN MATCHED AND c.action = 'mod' THEN UPDATE SET qty = qty + c.qty, ts = CURRENT_TIMESTAMP " +
		"WHEN MATCHED AND c.action = 'del' THEN DELETE " +
		"WHEN NOT MATCHED AND c.action = 'new' THEN INSERT (part, qty) VALUES (c.part, c.qty)",
	"MERGE INTO accounts USING src s ON accounts.c = s.c WHEN MATCHED THEN DELETE",                               // no target alias
	"MERGE INTO accounts AS t USING src s ON t.c = s.c WHEN MATCHED THEN DELETE",                                 // AS alias
	"MERGE INTO acc data USING src s ON data.c = s.c WHEN MATCHED THEN DELETE",                                   // non-reserved kw alias
	"MERGE INTO acc t USING (SELECT * FROM src) s ON t.c = s.c WHEN MATCHED THEN DELETE",                         // subquery source
	"MERGE INTO acc t USING (SELECT * FROM src) AS s ON t.c = s.c WHEN MATCHED THEN DELETE",                      // subquery AS source
	"MERGE INTO cat.sch.acc t USING src s ON t.c = s.c WHEN MATCHED THEN DELETE",                                 // qualified target
	"MERGE INTO acc t USING cat.sch.src s ON t.c = s.c WHEN MATCHED THEN DELETE",                                 // qualified source
	"MERGE INTO acc t USING src s ON t.c = s.c WHEN NOT MATCHED THEN INSERT VALUES (s.c, s.p)",                   // INSERT no column list
	"MERGE INTO acc t USING src s ON t.c = s.c WHEN MATCHED AND s.x > 0 THEN DELETE",                             // matched AND guard
	"MERGE INTO acc t USING src s ON t.c = s.c WHEN NOT MATCHED AND s.x > 0 THEN INSERT VALUES (1)",              // not-matched AND guard
	"MERGE INTO acc t USING src s ON t.c = s.c WHEN MATCHED THEN DELETE WHEN NOT MATCHED THEN INSERT VALUES (1)", // multi-clause
	"MERGE INTO acc t USING src CROSS JOIN other ON t.c = 1 WHEN MATCHED THEN DELETE",                            // D-MRG4: join source
	// MERGE negatives
	"MERGE INTO accounts t USING src s ON t.c = s.c",                                                 // D-MRG1: no WHEN
	"MERGE INTO accounts t USING src s ON t.c = s.c WHEN MATCHED THEN UPDATE SET (p = s.p)",          // D-MRG2: parenthesized SET
	"MERGE INTO accounts t USING src s ON t.c = s.c WHEN MATCHED THEN UPDATE SET (p = s.p, a = s.a)", // D-MRG2
	"MERGE INTO acc t USING src s ON t.c = s.c WHEN NOT MATCHED THEN INSERT VALUES ()",               // D-MRG3: empty VALUES
	"MERGE INTO a.b.c.d t USING s ON t.c = s.c WHEN MATCHED THEN DELETE",                             // D-MRG4: 4-part target
	"MERGE INTO acc t USING a JOIN b ON a.k = b.k s ON t.c = a.c WHEN MATCHED THEN DELETE",           // alias after join
	"MERGE INTO acc t USING src s WHEN MATCHED THEN DELETE",                                          // missing ON
	"MERGE acc t USING src s ON t.c = s.c WHEN MATCHED THEN DELETE",                                  // missing INTO
	"MERGE INTO acc t USING src s ON t.c = s.c WHEN MATCHED THEN INSERT VALUES (1)",                  // INSERT under MATCHED
	"MERGE INTO acc t USING src s ON t.c = s.c WHEN NOT MATCHED THEN DELETE",                         // DELETE under NOT MATCHED
	"MERGE INTO acc t USING src s ON t.c = s.c WHEN NOT MATCHED THEN UPDATE SET p = 1",               // UPDATE under NOT MATCHED
	"MERGE INTO acc t USING src s ON t.c = s.c WHEN NOT MATCHED THEN DELETE VALUES (1)",              // NOT MATCHED action must be INSERT (not DELETE)
	"MERGE INTO acc t USING src s ON t.c = s.c WHEN NOT MATCHED THEN UPDATE VALUES (1)",              // NOT MATCHED action must be INSERT (not UPDATE)
	"MERGE INTO acc t USING src s ON t.c = s.c WHEN NOT MATCHED THEN foo VALUES (1)",                 // NOT MATCHED action must be the INSERT keyword
}

// TestDML_OracleDifferential is the authoritative accept/reject gate: omni's
// Parse verdict must equal Trino 481's verdict for every corpus statement.
func TestDML_OracleDifferential(t *testing.T) {
	if testing.Short() {
		t.Skip("trino oracle: skipped in -short mode")
	}
	o := connectOracle(t)
	for _, sql := range dmlOracleCorpus {
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

// dmlBranchDivergenceCorpus holds the `@ branch_name` DML forms (divergence
// ledger #5). Trino 481 ACCEPTS each; omni REJECTS each because the lexer emits
// no '@' token. This documents and pins the known divergence so a future lexer
// change that adds the '@' token (and the corresponding DML support) is detected
// here as a behavior change.
var dmlBranchDivergenceCorpus = []string{
	"INSERT INTO cities @ audit VALUES (1, 'San Francisco')",
	"DELETE FROM orders @ audit",
	"DELETE FROM orders @ audit WHERE x = 1",
	"UPDATE purchases @ audit SET status = 'OVERDUE' WHERE ship_date IS NULL",
	"MERGE INTO accounts @ audit t USING src s ON t.c = s.c WHEN MATCHED THEN DELETE",
}

// TestDML_BranchDivergence asserts the flagged divergence #5 is still in effect:
// Trino accepts `@ branch`, omni rejects it. If omni ever starts accepting these
// (e.g. after the lexer learns '@'), this test FAILS loudly so the divergence
// ledger and this corpus are revisited and the cases promoted into the parity
// gate.
func TestDML_BranchDivergence(t *testing.T) {
	if testing.Short() {
		t.Skip("trino oracle: skipped in -short mode")
	}
	o := connectOracle(t)
	for _, sql := range dmlBranchDivergenceCorpus {
		sql := sql
		t.Run(truncateName(sql), func(t *testing.T) {
			_, errs := Parse(sql)
			omniAccepts := len(errs) == 0
			if omniAccepts {
				t.Errorf("DIVERGENCE RESOLVED: omni now ACCEPTS %q — Trino 481 accepts it too; "+
					"promote this case into dmlOracleCorpus and update divergence ledger #5", sql)
			}

			trinoAccepts, ok := oracleAccepts(t, o, sql)
			if !ok {
				t.Skip("oracle unreachable for this case")
			}
			if !trinoAccepts {
				t.Errorf("ORACLE CHANGED: Trino 481 now REJECTS %q — revisit divergence ledger #5", sql)
			}
		})
	}
}
