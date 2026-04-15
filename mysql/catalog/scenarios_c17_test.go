package catalog

import (
	"strings"
	"testing"
)

// TestScenario_C17 covers section C17 (String function charset / collation
// propagation) from SCENARIOS-mysql-implicit-behavior.md. 8 scenarios, each
// asserted on both a MySQL 8.0 container and the omni catalog. Every C17
// scenario is expected to currently fail on omni because analyze_expr.go /
// function_types.go do not track charset, collation, or derivation
// coercibility. Each omni failure is documented in scenarios_bug_queue/c17.md
// and reported via t.Errorf as KNOWN BUG so proof stays compile-clean.
func TestScenario_C17(t *testing.T) {
	scenariosSkipIfShort(t)
	scenariosSkipIfNoDocker(t)

	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	// c17OracleViewCol fetches (CHARACTER_SET_NAME, COLLATION_NAME) for a
	// view column, tolerating NULL as "".
	c17OracleViewCol := func(t *testing.T, view, col string) (charset, collation string, ok bool) {
		t.Helper()
		var cs, co any
		row := mc.db.QueryRowContext(mc.ctx,
			`SELECT CHARACTER_SET_NAME, COLLATION_NAME
			   FROM information_schema.COLUMNS
			  WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME=? AND COLUMN_NAME=?`,
			view, col)
		if err := row.Scan(&cs, &co); err != nil {
			t.Errorf("oracle view col (%s.%s): %v", view, col, err)
			return "", "", false
		}
		return asString(cs), asString(co), true
	}

	// c17OracleViewMaxLen returns CHARACTER_MAXIMUM_LENGTH for a view column.
	c17OracleViewMaxLen := func(t *testing.T, view, col string) (int, bool) {
		t.Helper()
		var v any
		row := mc.db.QueryRowContext(mc.ctx,
			`SELECT CHARACTER_MAXIMUM_LENGTH FROM information_schema.COLUMNS
			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME=? AND COLUMN_NAME=?`,
			view, col)
		if err := row.Scan(&v); err != nil {
			t.Errorf("oracle max_len (%s.%s): %v", view, col, err)
			return 0, false
		}
		switch x := v.(type) {
		case nil:
			return 0, false
		case int64:
			return int(x), true
		case []byte:
			n := 0
			for _, b := range x {
				if b >= '0' && b <= '9' {
					n = n*10 + int(b-'0')
				}
			}
			return n, true
		}
		return 0, false
	}

	// c17RunExpectError runs a DDL expected to fail and returns both errors.
	c17RunExpectError := func(t *testing.T, c *Catalog, ddl string) (oracleErr, omniErr error) {
		t.Helper()
		_, oracleErr = mc.db.ExecContext(mc.ctx, ddl)
		results, parseErr := c.Exec(ddl, nil)
		if parseErr != nil {
			omniErr = parseErr
			return
		}
		for _, r := range results {
			if r.Error != nil {
				omniErr = r.Error
				return
			}
		}
		return
	}

	// c17OmniViewExists checks whether omni's catalog has a view by name.
	c17OmniViewExists := func(c *Catalog, view string) bool {
		db := c.GetDatabase("testdb")
		if db == nil {
			return false
		}
		_, ok := db.Views[toLower(view)]
		return ok
	}

	// ---- 17.1 CONCAT identical charset/collation --------------------------
	t.Run("17_1_CONCAT_same_charset", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (
			a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
			b VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci
		)`)
		runOnBoth(t, mc, c, `CREATE VIEW v1 AS SELECT CONCAT(a, b) AS c FROM t`)

		cs, co, ok := c17OracleViewCol(t, "v1", "c")
		if ok {
			assertStringEq(t, "oracle v1.c CHARACTER_SET_NAME", cs, "utf8mb4")
			assertStringEq(t, "oracle v1.c COLLATION_NAME", co, "utf8mb4_0900_ai_ci")
		}
		if ml, ok := c17OracleViewMaxLen(t, "v1", "c"); ok {
			assertIntEq(t, "oracle v1.c CHARACTER_MAXIMUM_LENGTH", ml, 20)
		}

		// omni: no charset metadata on the view target list (KNOWN GAP).
		if !c17OmniViewExists(c, "v1") {
			t.Errorf("omni: view v1 not created")
		} else {
			t.Errorf("omni: KNOWN GAP — CONCAT result carries no charset/collation metadata (17.1), see scenarios_bug_queue/c17.md")
		}
	})

	// ---- 17.2 CONCAT mixing latin1 + utf8mb4 ------------------------------
	t.Run("17_2_CONCAT_latin1_utf8mb4_superset", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (
			a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
			b VARCHAR(10) CHARACTER SET latin1  COLLATE latin1_swedish_ci
		)`)
		runOnBoth(t, mc, c, `CREATE VIEW v2 AS SELECT CONCAT(a, b) AS c FROM t`)

		cs, co, ok := c17OracleViewCol(t, "v2", "c")
		if ok {
			assertStringEq(t, "oracle v2.c CHARACTER_SET_NAME", cs, "utf8mb4")
			assertStringEq(t, "oracle v2.c COLLATION_NAME", co, "utf8mb4_0900_ai_ci")
		}

		if c17OmniViewExists(c, "v2") {
			t.Errorf("omni: KNOWN GAP — CONCAT superset widening (latin1→utf8mb4) not tracked (17.2)")
		}
	})

	// ---- 17.3 CONCAT incompatible collations should error -----------------
	t.Run("17_3_CONCAT_incompatible_collations", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (
			a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
			b VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_as_cs
		)`)
		if _, err := mc.db.ExecContext(mc.ctx, `INSERT INTO t VALUES ('x','y')`); err != nil {
			t.Errorf("oracle INSERT: %v", err)
		}

		// ORACLE FINDING: MySQL 8.0.45 does NOT raise 1267 for
		// CONCAT(utf8mb4_0900_ai_ci, utf8mb4_0900_as_cs) — the CONCAT path
		// silently widens via a newer pad-space-compat rule. The canonical
		// illegal-mix trigger for this pair is the `=` comparison path
		// (Item_bool_func2::fix_length_and_dec), which we use here as the
		// stable probe that forces DTCollation::aggregate to fail.
		_, oracleCmpErr := mc.db.ExecContext(mc.ctx, `SELECT 1 FROM t WHERE a = b`)
		if oracleCmpErr == nil {
			t.Errorf("oracle: comparison of two incompatible IMPLICIT collations unexpectedly accepted — expected ER_CANT_AGGREGATE_2COLLATIONS (1267)")
		} else if !strings.Contains(oracleCmpErr.Error(), "1267") &&
			!strings.Contains(strings.ToLower(oracleCmpErr.Error()), "illegal mix") {
			t.Errorf("oracle: unexpected error: %v", oracleCmpErr)
		}

		// omni: SELECT ... WHERE a=b must reject (aggregation gap).
		results, parseErr := c.Exec(`CREATE VIEW v3 AS SELECT a FROM t WHERE a = b`, nil)
		omniAccepted := parseErr == nil
		if parseErr == nil {
			for _, r := range results {
				if r.Error != nil {
					omniAccepted = false
					break
				}
			}
		}
		if omniAccepted {
			t.Errorf("omni: KNOWN BUG — soft-accept of illegal-mix comparison (17.3); should error 1267. See scenarios_bug_queue/c17.md")
		}
	})

	// ---- 17.4 CONCAT_WS NULL skipping + separator aggregation -------------
	t.Run("17_4_CONCAT_WS_nullskip", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (
			a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
			b VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci
		)`)
		runOnBoth(t, mc, c, `CREATE VIEW v4 AS SELECT CONCAT_WS(',', a, b, NULL) AS c FROM t`)

		if cs, co, ok := c17OracleViewCol(t, "v4", "c"); ok {
			assertStringEq(t, "oracle v4.c CHARACTER_SET_NAME", cs, "utf8mb4")
			assertStringEq(t, "oracle v4.c COLLATION_NAME", co, "utf8mb4_0900_ai_ci")
		}

		// Runtime: CONCAT(NULL,'x') is NULL; CONCAT_WS(',',NULL,'x') is 'x'.
		if _, err := mc.db.ExecContext(mc.ctx, `INSERT INTO t VALUES (NULL, 'x')`); err != nil {
			t.Errorf("oracle INSERT: %v", err)
		}
		var concatRes, cwsRes any
		row := mc.db.QueryRowContext(mc.ctx,
			`SELECT CONCAT(a,b), CONCAT_WS(',',a,b) FROM t LIMIT 1`)
		if err := row.Scan(&concatRes, &cwsRes); err != nil {
			t.Errorf("oracle runtime scan: %v", err)
		} else {
			if concatRes != nil {
				t.Errorf("oracle CONCAT(NULL,'x'): got %v, want NULL", concatRes)
			}
			if s := asString(cwsRes); s != "x" {
				t.Errorf("oracle CONCAT_WS(',',NULL,'x'): got %q, want \"x\"", s)
			}
		}

		if c17OmniViewExists(c, "v4") {
			t.Errorf("omni: KNOWN GAP — CONCAT_WS NULL-skip semantics + charset aggregation not tracked (17.4)")
		}
	})

	// ---- 17.5 _utf8mb4'x' introducer is still COERCIBLE -------------------
	t.Run("17_5_introducer_still_coercible", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		// MySQL session: force a known session charset.
		if _, err := mc.db.ExecContext(mc.ctx, `SET NAMES utf8mb4`); err != nil {
			t.Errorf("oracle SET NAMES: %v", err)
		}

		runOnBoth(t, mc, c,
			`CREATE TABLE t (a VARCHAR(10) CHARACTER SET latin1 COLLATE latin1_swedish_ci)`)
		runOnBoth(t, mc, c, `CREATE VIEW v5a AS SELECT CONCAT(a, 'x') AS c FROM t`)
		runOnBoth(t, mc, c, `CREATE VIEW v5b AS SELECT CONCAT(a, _utf8mb4'x') AS c FROM t`)

		for _, view := range []string{"v5a", "v5b"} {
			if cs, co, ok := c17OracleViewCol(t, view, "c"); ok {
				assertStringEq(t, "oracle "+view+".c CHARACTER_SET_NAME", cs, "latin1")
				assertStringEq(t, "oracle "+view+".c COLLATION_NAME", co, "latin1_swedish_ci")
			}
		}

		if c17OmniViewExists(c, "v5a") || c17OmniViewExists(c, "v5b") {
			t.Errorf("omni: KNOWN GAP — literal/introducer coercibility (COERCIBLE vs IMPLICIT) not tracked (17.5)")
		}
	})

	// ---- 17.6 REPEAT / LPAD / RPAD pin to first-arg charset ---------------
	t.Run("17_6_first_arg_pins_charset", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (
			a VARCHAR(10) CHARACTER SET latin1  COLLATE latin1_swedish_ci,
			b VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci
		)`)
		runOnBoth(t, mc, c, `CREATE VIEW v6a AS SELECT REPEAT(a, 3) AS c FROM t`)
		runOnBoth(t, mc, c, `CREATE VIEW v6b AS SELECT LPAD(a, 20, b) AS c FROM t`)
		runOnBoth(t, mc, c, `CREATE VIEW v6c AS SELECT RPAD(b, 20, a) AS c FROM t`)

		type want struct{ cs, co string }
		cases := map[string]want{
			"v6a": {"latin1", "latin1_swedish_ci"},
			"v6b": {"latin1", "latin1_swedish_ci"},
			"v6c": {"utf8mb4", "utf8mb4_0900_ai_ci"},
		}
		for view, w := range cases {
			if cs, co, ok := c17OracleViewCol(t, view, "c"); ok {
				assertStringEq(t, "oracle "+view+".c CHARACTER_SET_NAME", cs, w.cs)
				assertStringEq(t, "oracle "+view+".c COLLATION_NAME", co, w.co)
			}
		}

		anyExists := c17OmniViewExists(c, "v6a") || c17OmniViewExists(c, "v6b") || c17OmniViewExists(c, "v6c")
		if anyExists {
			t.Errorf("omni: KNOWN GAP — REPEAT/LPAD/RPAD first-arg-pins rule missing (17.6). Fix in function_types.go")
		}
	})

	// ---- 17.7 CONVERT(x USING cs) forces charset IMPLICIT -----------------
	t.Run("17_7_CONVERT_USING_pins_charset", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (
			a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_as_cs,
			b VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci
		)`)
		// Without CONVERT this would be 17.3 (ER 1267). CONVERT rescues it.
		ddl := `CREATE VIEW v7 AS SELECT CONCAT(a, CONVERT(b USING utf8mb4)) AS c FROM t`
		_, oracleErr := mc.db.ExecContext(mc.ctx, ddl)
		if oracleErr != nil {
			// If MySQL rejects even one-sided CONVERT, the scenario has no
			// oracle ground truth to compare omni against. Record this so the
			// scenario can be refined and return — do NOT assert omni-side
			// behavior against a missing oracle.
			t.Skipf("17.7 oracle rejected one-sided CONVERT; scenario needs two-sided wrap. err=%v", oracleErr)
			return
		}
		if cs, co, ok := c17OracleViewCol(t, "v7", "c"); ok {
			assertStringEq(t, "oracle v7.c CHARACTER_SET_NAME", cs, "utf8mb4")
			// The default collation of utf8mb4 is server-configurable; accept
			// any utf8mb4_* collation. What matters is the charset pin.
			if !strings.HasPrefix(co, "utf8mb4_") {
				t.Errorf("oracle v7.c COLLATION_NAME: got %q, want utf8mb4_* (some utf8mb4 collation)", co)
			}
		}

		// omni side: parse the same DDL. Only reached when oracle accepted,
		// so the KNOWN GAP comparison is meaningful.
		results, parseErr := c.Exec(ddl, nil)
		omniAccepted := parseErr == nil
		if parseErr == nil {
			for _, r := range results {
				if r.Error != nil {
					omniAccepted = false
					break
				}
			}
		}
		if omniAccepted {
			t.Errorf("omni: KNOWN GAP — CONVERT ... USING cs accepted but charset not pinned on result (17.7), see scenarios_bug_queue/c17.md")
		}
	})

	// ---- 17.8 COLLATE clause is EXPLICIT — highest precedence -------------
	t.Run("17_8_COLLATE_explicit", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (
			a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
			b VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_as_cs
		)`)

		// Case A: one side EXPLICIT wins.
		ddlA := `CREATE VIEW v8a AS SELECT CONCAT(a, b COLLATE utf8mb4_bin) AS c FROM t`
		if _, err := mc.db.ExecContext(mc.ctx, ddlA); err != nil {
			t.Errorf("oracle v8a unexpected error: %v", err)
		} else {
			if cs, co, ok := c17OracleViewCol(t, "v8a", "c"); ok {
				assertStringEq(t, "oracle v8a.c CHARACTER_SET_NAME", cs, "utf8mb4")
				assertStringEq(t, "oracle v8a.c COLLATION_NAME", co, "utf8mb4_bin")
			}
		}

		// omni: same DDL should analyze — gap is that EXPLICIT derivation
		// isn't tracked, so downstream charset metadata is missing.
		resultsA, parseErrA := c.Exec(ddlA, nil)
		omniAcceptedA := parseErrA == nil
		if parseErrA == nil {
			for _, r := range resultsA {
				if r.Error != nil {
					omniAcceptedA = false
					break
				}
			}
		}
		if omniAcceptedA {
			t.Errorf("omni: KNOWN GAP — COLLATE EXPLICIT derivation not tracked on v8a (17.8)")
		}

		// Case B: two EXPLICIT sides with different collations must error.
		scenarioReset(t, mc)
		cB := scenarioNewCatalog(t)
		runOnBoth(t, mc, cB, `CREATE TABLE t (
			a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
			b VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_as_cs
		)`)
		ddlB := `CREATE VIEW v8b AS SELECT CONCAT(a COLLATE utf8mb4_0900_ai_ci, b COLLATE utf8mb4_bin) AS c FROM t`
		oracleErr, omniErr := c17RunExpectError(t, cB, ddlB)
		if oracleErr == nil {
			t.Errorf("oracle: %q unexpectedly accepted — expected illegal-mix error (1267/1270)", ddlB)
		} else if !strings.Contains(oracleErr.Error(), "1267") &&
			!strings.Contains(oracleErr.Error(), "1270") &&
			!strings.Contains(strings.ToLower(oracleErr.Error()), "illegal mix") {
			t.Errorf("oracle v8b unexpected error: %v", oracleErr)
		}
		if omniErr == nil {
			t.Errorf("omni: KNOWN BUG — two EXPLICIT COLLATE sides silently accepted (17.8), see scenarios_bug_queue/c17.md")
		}
	})
}
