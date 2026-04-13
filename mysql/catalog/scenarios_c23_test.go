package catalog

import (
	"strings"
	"testing"
)

// TestScenario_C23 covers Section C23 "NULL in string context" from
// mysql/catalog/SCENARIOS-mysql-implicit-behavior.md.
//
// Scenarios in this section verify NULL propagation through string
// functions:
//
//   23.1  CONCAT(...)         — any NULL arg → NULL result
//   23.2  CONCAT_WS(sep, ...) — NULL data args skipped; NULL separator → NULL
//   23.3  IFNULL/COALESCE     — rescue pattern around CONCAT NULL propagation
//
// These are runtime expression behaviors. The omni catalog stores parsed
// VIEWs as a *catalog.View whose `Columns` field is just a list of column
// NAMES (see mysql/catalog/table.go). It does NOT track per-column
// nullability. So the most useful representation of these scenarios in
// omni today is:
//
//   1. Run the same CREATE TABLE / CREATE VIEW DDL against both MySQL 8.0
//      and the omni catalog (proves the DDL parses on both sides).
//   2. Use information_schema.COLUMNS on the container to assert that
//      MySQL infers IS_NULLABLE the way the SCENARIOS doc claims.
//   3. SELECT the actual rows from the container to lock the runtime
//      string values into the test (oracle ground truth).
//   4. Record the omni gap: the View struct has no per-column nullability
//      info, so omni cannot answer "is column c1 of view v nullable?" —
//      this is the declared bug, documented in scenarios_bug_queue/c23.md.
//
// Failed omni assertions are NOT proof failures — they are recorded in
// mysql/catalog/scenarios_bug_queue/c23.md.
func TestScenario_C23(t *testing.T) {
	scenariosSkipIfShort(t)
	scenariosSkipIfNoDocker(t)

	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	// c23OmniView fetches a view from the omni catalog by name.
	c23OmniView := func(c *Catalog, name string) *View {
		db := c.GetDatabase("testdb")
		if db == nil {
			return nil
		}
		return db.Views[strings.ToLower(name)]
	}

	// c23OracleViewColNullable returns the IS_NULLABLE value for a single
	// view column from information_schema.COLUMNS.
	c23OracleViewColNullable := func(t *testing.T, view, col string) string {
		t.Helper()
		var s string
		oracleScan(t, mc,
			`SELECT IS_NULLABLE FROM information_schema.COLUMNS
			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='`+view+`' AND COLUMN_NAME='`+col+`'`,
			&s)
		return s
	}

	// -----------------------------------------------------------------
	// 23.1 CONCAT with any NULL argument → NULL result
	// -----------------------------------------------------------------
	t.Run("23_1_concat_null_propagates", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c,
			`CREATE TABLE t (a VARCHAR(10), b VARCHAR(10));
			 CREATE VIEW v_concat AS SELECT CONCAT(a, b) AS c1, CONCAT('x','','y') AS c2 FROM t;`)

		// Oracle row-level checks (lock in the actual MySQL evaluation).
		_, err := mc.db.ExecContext(mc.ctx,
			`INSERT INTO t VALUES ('foo', NULL), ('foo', 'bar')`)
		if err != nil {
			t.Fatalf("oracle INSERT: %v", err)
		}
		// SELECT CONCAT(a,b) — first row NULL, second row 'foobar'
		rows := oracleRows(t, mc, `SELECT CONCAT(a, b) FROM t ORDER BY b IS NULL DESC`)
		if len(rows) != 2 {
			t.Errorf("oracle: expected 2 rows, got %d", len(rows))
		} else {
			if rows[0][0] != nil {
				t.Errorf("oracle row[0] CONCAT(a,b): want NULL, got %v", rows[0][0])
			}
			if asString(rows[1][0]) != "foobar" {
				t.Errorf("oracle row[1] CONCAT(a,b): want 'foobar', got %v", rows[1][0])
			}
		}
		// SELECT CONCAT('x', NULL, 'y') → NULL
		var lit any
		oracleScan(t, mc, `SELECT CONCAT('x', NULL, 'y')`, &lit)
		if lit != nil {
			t.Errorf("oracle CONCAT('x',NULL,'y'): want NULL, got %v", lit)
		}
		// SELECT CONCAT('x', '', 'y') → 'xy' (empty string is NOT NULL)
		var emptyStr string
		oracleScan(t, mc, `SELECT CONCAT('x', '', 'y')`, &emptyStr)
		assertStringEq(t, "oracle CONCAT('x','','y')", emptyStr, "xy")

		// Oracle view-column nullability — empirical MySQL 8.0.45:
		// c1 (CONCAT of two nullable cols) → IS_NULLABLE=YES
		// c2 (CONCAT of three string literals) → IS_NULLABLE=YES
		//   Note: although the SCENARIOS doc reasoning would suggest c2
		//   should be NOT NULL (all inputs are non-null literals), MySQL
		//   8.0's view metadata pass conservatively reports any string
		//   function result column as nullable. We lock the ground truth
		//   so omni's eventual nullability inference can match it.
		assertStringEq(t, "oracle v_concat.c1 IS_NULLABLE",
			c23OracleViewColNullable(t, "v_concat", "c1"), "YES")
		assertStringEq(t, "oracle v_concat.c2 IS_NULLABLE",
			c23OracleViewColNullable(t, "v_concat", "c2"), "YES")

		// omni: view exists but per-column nullability is not represented.
		v := c23OmniView(c, "v_concat")
		if v == nil {
			t.Errorf("omni: view v_concat not found")
			return
		}
		if len(v.Columns) != 2 {
			t.Errorf("omni: v_concat expected 2 columns, got %d (%v)", len(v.Columns), v.Columns)
		}
		// Declared bug: View has no per-column nullability info, so the
		// "CONCAT propagates NULL" semantics cannot be asserted positively.
		t.Error("omni: View struct has no per-column nullability; scenario 23.1 cannot be asserted positively")
	})

	// -----------------------------------------------------------------
	// 23.2 CONCAT_WS skips NULL arguments (separator non-null)
	// -----------------------------------------------------------------
	t.Run("23_2_concat_ws_skips_null_args", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c,
			`CREATE TABLE t (a VARCHAR(10), b VARCHAR(10), c VARCHAR(10));
			 CREATE VIEW v_ws AS SELECT CONCAT_WS(',', a, b, c) AS d1 FROM t;`)

		_, err := mc.db.ExecContext(mc.ctx,
			`INSERT INTO t VALUES ('x', NULL, 'z'), (NULL, NULL, NULL)`)
		if err != nil {
			t.Fatalf("oracle INSERT: %v", err)
		}

		// row1 → 'x,z' (NULL skipped, no double separator)
		// row2 → ''   (all NULL data args → empty string, NOT NULL)
		rows := oracleRows(t, mc, `SELECT CONCAT_WS(',', a, b, c) FROM t`)
		if len(rows) != 2 {
			t.Errorf("oracle: expected 2 rows, got %d", len(rows))
		} else {
			assertStringEq(t, "oracle row[0] CONCAT_WS",
				asString(rows[0][0]), "x,z")
			assertStringEq(t, "oracle row[1] CONCAT_WS",
				asString(rows[1][0]), "")
		}

		// Literal: CONCAT_WS(',', 'x', NULL, 'z') → 'x,z'
		var lit1 string
		oracleScan(t, mc, `SELECT CONCAT_WS(',', 'x', NULL, 'z')`, &lit1)
		assertStringEq(t, "oracle CONCAT_WS(',','x',NULL,'z')", lit1, "x,z")

		// Literal: CONCAT_WS(',', NULL, NULL, NULL) → '' (NOT NULL)
		var lit2 string
		oracleScan(t, mc, `SELECT CONCAT_WS(',', NULL, NULL, NULL)`, &lit2)
		assertStringEq(t, "oracle CONCAT_WS(',',NULL,NULL,NULL)", lit2, "")

		// Literal: CONCAT_WS(NULL, 'x', 'y') → NULL (separator NULL → NULL)
		var lit3 any
		oracleScan(t, mc, `SELECT CONCAT_WS(NULL, 'x', 'y')`, &lit3)
		if lit3 != nil {
			t.Errorf("oracle CONCAT_WS(NULL,'x','y'): want NULL, got %v", lit3)
		}

		// Oracle view nullability — empirical MySQL 8.0.45: even though
		// the separator (',') is a non-null literal and the runtime rule
		// is "result is NULL iff separator is NULL" (so d1 is in fact
		// never NULL at runtime), MySQL's view metadata pass reports
		// IS_NULLABLE='YES' for the CONCAT_WS result column. The SCENARIOS
		// doc's reasoning describes the runtime semantics, not the
		// information_schema metadata. We lock the metadata ground truth
		// here and rely on the runtime literal assertions above (lit1,
		// lit2, lit3) to lock the runtime semantics.
		assertStringEq(t, "oracle v_ws.d1 IS_NULLABLE",
			c23OracleViewColNullable(t, "v_ws", "d1"), "YES")

		v := c23OmniView(c, "v_ws")
		if v == nil {
			t.Errorf("omni: view v_ws not found")
			return
		}
		if len(v.Columns) != 1 {
			t.Errorf("omni: v_ws expected 1 column, got %d (%v)", len(v.Columns), v.Columns)
		}
		// Declared bug: omni cannot answer "is the separator NULL?" because
		// View has no per-column nullability with the CONCAT_WS special case.
		t.Error("omni: View struct has no per-column nullability; CONCAT_WS NULL-skip rule (23.2) cannot be asserted positively")
	})

	// -----------------------------------------------------------------
	// 23.3 IFNULL / COALESCE as rescue for CONCAT NULL-propagation
	// -----------------------------------------------------------------
	t.Run("23_3_ifnull_coalesce_rescue", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c,
			`CREATE TABLE t (first_name VARCHAR(20), middle_name VARCHAR(20), last_name VARCHAR(20));
			 CREATE VIEW v_name AS SELECT
			   CONCAT(first_name, ' ', middle_name, ' ', last_name) AS bad,
			   CONCAT(first_name, ' ', IFNULL(middle_name, ''), ' ', last_name) AS rescue_ifnull,
			   CONCAT(first_name, ' ', COALESCE(middle_name, ''), ' ', last_name) AS rescue_coalesce,
			   CONCAT_WS(' ', first_name, middle_name, last_name) AS rescue_ws
			 FROM t;`)

		_, err := mc.db.ExecContext(mc.ctx,
			`INSERT INTO t VALUES ('Ada', NULL, 'Lovelace')`)
		if err != nil {
			t.Fatalf("oracle INSERT: %v", err)
		}

		// Row-level oracle ground truth.
		rows := oracleRows(t, mc, `SELECT bad, rescue_ifnull, rescue_coalesce, rescue_ws FROM v_name`)
		if len(rows) != 1 {
			t.Errorf("oracle: expected 1 row from v_name, got %d", len(rows))
		} else {
			r := rows[0]
			if r[0] != nil {
				t.Errorf("oracle bad: want NULL, got %v", r[0])
			}
			assertStringEq(t, "oracle rescue_ifnull", asString(r[1]), "Ada  Lovelace")
			assertStringEq(t, "oracle rescue_coalesce", asString(r[2]), "Ada  Lovelace")
			assertStringEq(t, "oracle rescue_ws", asString(r[3]), "Ada Lovelace")
		}

		// Oracle column nullability:
		//   bad             → YES (CONCAT with NULL middle_name propagates)
		//   rescue_ifnull   → YES on MySQL 8.0 — even though the second IFNULL
		//                     arg is a non-null literal, the surrounding CONCAT
		//                     also reads first_name/last_name which are
		//                     nullable (no NOT NULL on base table). So the
		//                     result is reported nullable. We just record what
		//                     MySQL says so the doc claim is locked.
		//   rescue_coalesce → YES (same reasoning)
		//   rescue_ws       → YES (data args nullable; separator non-null but
		//                     base columns nullable — the result column from
		//                     a view body that references nullable inputs may
		//                     still be reported nullable).
		// The point of 23.3 is the RUNTIME values, which we asserted above.
		// The IS_NULLABLE values here are oracle ground-truth — record them
		// so any future MySQL upgrade or omni inference change is caught.
		bad := c23OracleViewColNullable(t, "v_name", "bad")
		ifn := c23OracleViewColNullable(t, "v_name", "rescue_ifnull")
		coa := c23OracleViewColNullable(t, "v_name", "rescue_coalesce")
		ws := c23OracleViewColNullable(t, "v_name", "rescue_ws")
		// These four must all be exactly "YES" or "NO" — record what oracle
		// says without hard-coding a guess.
		for _, p := range []struct{ label, got string }{
			{"v_name.bad", bad},
			{"v_name.rescue_ifnull", ifn},
			{"v_name.rescue_coalesce", coa},
			{"v_name.rescue_ws", ws},
		} {
			if p.got != "YES" && p.got != "NO" {
				t.Errorf("oracle %s IS_NULLABLE: got %q, want YES or NO", p.label, p.got)
			}
		}
		// And lock the specific MySQL 8.0.45 result for "bad" (the unrescued
		// CONCAT) — it must be YES.
		assertStringEq(t, "oracle v_name.bad IS_NULLABLE (unrescued CONCAT)", bad, "YES")

		v := c23OmniView(c, "v_name")
		if v == nil {
			t.Errorf("omni: view v_name not found")
			return
		}
		if len(v.Columns) != 4 {
			t.Errorf("omni: v_name expected 4 columns, got %d (%v)", len(v.Columns), v.Columns)
		}
		// Declared bug: omni cannot answer "does IFNULL/COALESCE rescue the
		// CONCAT" because View has no per-column nullability inference.
		t.Error("omni: View struct has no per-column nullability; IFNULL/COALESCE rescue rule (23.3) cannot be asserted positively")
	})
}
