package catalog

import (
	"strings"
	"testing"
)

// TestScenario_C16 covers section C16 (Date/time function precision defaults)
// from SCENARIOS-mysql-implicit-behavior.md. 12 scenarios, each asserted on
// both MySQL 8.0 container and the omni catalog. Failures in omni assertions
// are documented in scenarios_bug_queue/c16.md (NOT proof failures).
func TestScenario_C16(t *testing.T) {
	scenariosSkipIfShort(t)
	scenariosSkipIfNoDocker(t)

	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	// Helper: fetch IS.COLUMNS.DATETIME_PRECISION for testdb.t.<col>.
	c16OracleDatetimePrecision := func(t *testing.T, col string) (int, bool) {
		t.Helper()
		var v any
		row := mc.db.QueryRowContext(mc.ctx,
			`SELECT DATETIME_PRECISION FROM information_schema.COLUMNS
			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME=?`, col)
		if err := row.Scan(&v); err != nil {
			t.Errorf("oracle DATETIME_PRECISION(%q): %v", col, err)
			return 0, false
		}
		if v == nil {
			return 0, false // NULL (e.g. DATE columns)
		}
		switch x := v.(type) {
		case int64:
			return int(x), true
		case []byte:
			var n int
			for _, c := range x {
				if c >= '0' && c <= '9' {
					n = n*10 + int(c-'0')
				}
			}
			return n, true
		}
		return 0, false
	}

	// Helper: run a DDL that we expect to fail on MySQL. Returns the errors
	// from each side so the caller can assert.
	c16RunExpectError := func(t *testing.T, c *Catalog, ddl string) (oracleErr, omniErr error) {
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

	// ---- 16.1 NOW()/CURRENT_TIMESTAMP precision defaults to 0 -------------
	t.Run("16_1_NOW_precision_default_0", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (ts DATETIME DEFAULT NOW())`)

		// Oracle: SHOW CREATE TABLE should render DEFAULT CURRENT_TIMESTAMP
		// without any (n) suffix.
		create := oracleShow(t, mc, "SHOW CREATE TABLE t")
		lo := strings.ToLower(create)
		if !strings.Contains(lo, "default current_timestamp") {
			t.Errorf("oracle SHOW CREATE TABLE missing DEFAULT CURRENT_TIMESTAMP: %s", create)
		}
		if strings.Contains(lo, "current_timestamp(") {
			t.Errorf("oracle SHOW CREATE TABLE unexpectedly has fsp suffix: %s", create)
		}

		// Oracle: DATETIME_PRECISION = 0 for plain DATETIME col.
		if p, ok := c16OracleDatetimePrecision(t, "ts"); !ok {
			t.Errorf("oracle DATETIME_PRECISION NULL for plain DATETIME")
		} else if p != 0 {
			t.Errorf("oracle DATETIME_PRECISION: got %d, want 0", p)
		}

		// omni: column must exist with a default referencing CURRENT_TIMESTAMP
		// (any rendering: "now()", "CURRENT_TIMESTAMP", etc).
		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t not found")
			return
		}
		col := tbl.GetColumn("ts")
		if col == nil {
			t.Errorf("omni: column ts not found")
			return
		}
		if col.Default == nil {
			t.Errorf("omni: DEFAULT missing for ts")
		} else {
			lo := strings.ToLower(*col.Default)
			if !strings.Contains(lo, "current_timestamp") && !strings.Contains(lo, "now") {
				t.Errorf("omni: DEFAULT = %q, expected CURRENT_TIMESTAMP/NOW form", *col.Default)
			}
			if strings.Contains(lo, "(") && strings.Contains(lo, ")") && !strings.Contains(lo, "()") {
				// Contains some fsp number — unexpected for plain DATETIME.
				t.Errorf("omni: DEFAULT = %q, expected no fsp suffix", *col.Default)
			}
		}
	})

	// ---- 16.2 NOW(N) range 0..6; NOW(7) rejected --------------------------
	t.Run("16_2_NOW_explicit_precision_range", func(t *testing.T) {
		scenarioReset(t, mc)

		// 0..6 should all work at runtime via SELECT LENGTH(NOW(n)).
		wantLen := map[int]int{0: 19, 1: 21, 2: 22, 3: 23, 4: 24, 5: 25, 6: 26}
		for n, want := range wantLen {
			var got int
			row := mc.db.QueryRowContext(mc.ctx,
				`SELECT LENGTH(NOW(`+itoaC16(n)+`))`)
			if err := row.Scan(&got); err != nil {
				t.Errorf("oracle LENGTH(NOW(%d)): %v", n, err)
				continue
			}
			if got != want {
				t.Errorf("oracle LENGTH(NOW(%d)): got %d, want %d", n, got, want)
			}
		}

		// NOW(7) must be rejected with ER_TOO_BIG_PRECISION (1426).
		_, err := mc.db.ExecContext(mc.ctx, `DO NOW(7)`)
		if err == nil {
			t.Errorf("oracle: NOW(7) unexpectedly accepted")
		} else if !strings.Contains(err.Error(), "1426") &&
			!strings.Contains(strings.ToLower(err.Error()), "precision") {
			t.Errorf("oracle: NOW(7) unexpected error: %v", err)
		}

		// omni: NOW(7) as DEFAULT should be rejected (strictness gap if it isn't).
		c := scenarioNewCatalog(t)
		ddl := `CREATE TABLE t (a DATETIME(7))`
		results, parseErr := c.Exec(ddl, nil)
		var omniErr error
		if parseErr != nil {
			omniErr = parseErr
		} else {
			for _, r := range results {
				if r.Error != nil {
					omniErr = r.Error
					break
				}
			}
		}
		if omniErr == nil {
			t.Errorf("omni: KNOWN BUG — DATETIME(7) should be rejected (fsp > 6), see scenarios_bug_queue/c16.md")
		}
	})

	// ---- 16.3 CURDATE / CURRENT_DATE / UTC_DATE take no precision arg ----
	t.Run("16_3_CURDATE_no_precision", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		// Oracle: CURDATE() works, CURDATE(6) is a parse error.
		if _, err := mc.db.ExecContext(mc.ctx, `DO CURDATE()`); err != nil {
			t.Errorf("oracle CURDATE() failed: %v", err)
		}
		if _, err := mc.db.ExecContext(mc.ctx, `DO CURDATE(6)`); err == nil {
			t.Errorf("oracle: CURDATE(6) unexpectedly accepted")
		}

		// Table with default CURDATE() works and has NULL DATETIME_PRECISION.
		runOnBoth(t, mc, c, `CREATE TABLE t (d DATE DEFAULT (CURDATE()))`)
		if _, ok := c16OracleDatetimePrecision(t, "d"); ok {
			t.Errorf("oracle: DATE column should have NULL DATETIME_PRECISION")
		}

		// omni: CURDATE(6) as DEFAULT should be rejected by parser.
		c2 := scenarioNewCatalog(t)
		results, parseErr := c2.Exec(`SELECT CURDATE(6)`, nil)
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
			t.Errorf("omni: KNOWN BUG — CURDATE(6) should fail parse, see scenarios_bug_queue/c16.md")
		}
	})

	// ---- 16.4 CURTIME / CURRENT_TIME / UTC_TIME precision defaults to 0 --
	t.Run("16_4_CURTIME_precision_default_0", func(t *testing.T) {
		scenarioReset(t, mc)

		var l0, l6 int
		row := mc.db.QueryRowContext(mc.ctx,
			`SELECT LENGTH(CURTIME()), LENGTH(CURTIME(6))`)
		if err := row.Scan(&l0, &l6); err != nil {
			t.Errorf("oracle CURTIME length scan: %v", err)
		}
		if l0 != 8 {
			t.Errorf("oracle LENGTH(CURTIME())=%d, want 8", l0)
		}
		if l6 != 15 {
			t.Errorf("oracle LENGTH(CURTIME(6))=%d, want 15", l6)
		}

		// CURTIME(7) rejected
		if _, err := mc.db.ExecContext(mc.ctx, `DO CURTIME(7)`); err == nil {
			t.Errorf("oracle: CURTIME(7) unexpectedly accepted")
		}

		// omni: parses a SELECT CURTIME() successfully. Rejects CURTIME(7)?
		c := scenarioNewCatalog(t)
		results, parseErr := c.Exec(`SELECT CURTIME(7)`, nil)
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
			t.Errorf("omni: KNOWN BUG — CURTIME(7) should be rejected (ER_TOO_BIG_PRECISION), see scenarios_bug_queue/c16.md")
		}
	})

	// ---- 16.5 SYSDATE cannot be used as DEFAULT --------------------------
	t.Run("16_5_SYSDATE_not_allowed_as_default", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		// SYSDATE is not a valid bareword in DEFAULT (parser rejects as syntax
		// error) and SYSDATE() also fails add_field (not Item_func_now). Both
		// forms must be rejected by MySQL.
		ddl := `CREATE TABLE t (a DATETIME DEFAULT SYSDATE())`
		oracleErr, omniErr := c16RunExpectError(t, c, ddl)
		if oracleErr == nil {
			t.Errorf("oracle: DATETIME DEFAULT SYSDATE() unexpectedly accepted")
		}
		if omniErr == nil {
			t.Errorf("omni: KNOWN BUG — DATETIME DEFAULT SYSDATE should error (not a NOW_FUNC), see scenarios_bug_queue/c16.md")
		}
	})

	// ---- 16.6 UTC_TIMESTAMP precision defaults to 0; not valid as DEFAULT -
	t.Run("16_6_UTC_TIMESTAMP_precision_and_default", func(t *testing.T) {
		scenarioReset(t, mc)

		var l0, l3 int
		row := mc.db.QueryRowContext(mc.ctx,
			`SELECT LENGTH(UTC_TIMESTAMP()), LENGTH(UTC_TIMESTAMP(3))`)
		if err := row.Scan(&l0, &l3); err != nil {
			t.Errorf("oracle UTC_TIMESTAMP length scan: %v", err)
		}
		if l0 != 19 {
			t.Errorf("oracle LENGTH(UTC_TIMESTAMP())=%d, want 19", l0)
		}
		if l3 != 23 {
			t.Errorf("oracle LENGTH(UTC_TIMESTAMP(3))=%d, want 23", l3)
		}

		// UTC_TIMESTAMP as DEFAULT → ER_INVALID_DEFAULT on oracle.
		c := scenarioNewCatalog(t)
		ddl := `CREATE TABLE t (a DATETIME DEFAULT UTC_TIMESTAMP)`
		oracleErr, omniErr := c16RunExpectError(t, c, ddl)
		if oracleErr == nil {
			t.Errorf("oracle: DATETIME DEFAULT UTC_TIMESTAMP unexpectedly accepted")
		}
		if omniErr == nil {
			t.Errorf("omni: KNOWN BUG — DATETIME DEFAULT UTC_TIMESTAMP should error, see scenarios_bug_queue/c16.md")
		}
	})

	// ---- 16.7 UNIX_TIMESTAMP return type depends on arg ------------------
	t.Run("16_7_UNIX_TIMESTAMP_return_type", func(t *testing.T) {
		scenarioReset(t, mc)

		// Oracle: zero-arg UNIX_TIMESTAMP() returns BIGINT UNSIGNED;
		// UNIX_TIMESTAMP(NOW(6)) returns DECIMAL with scale 6. Observe via
		// the data type MySQL assigns to a view column.
		if _, err := mc.db.ExecContext(mc.ctx,
			`CREATE VIEW v AS SELECT UNIX_TIMESTAMP() AS u0, UNIX_TIMESTAMP(NOW(6)) AS u6`); err != nil {
			t.Errorf("oracle CREATE VIEW v: %v", err)
			return
		}
		var t0, t6 string
		var s6 any
		row := mc.db.QueryRowContext(mc.ctx,
			`SELECT DATA_TYPE, NUMERIC_SCALE FROM information_schema.COLUMNS
			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v' AND COLUMN_NAME='u6'`)
		if err := row.Scan(&t6, &s6); err != nil {
			t.Errorf("oracle u6 type scan: %v", err)
		} else {
			if strings.ToLower(t6) != "decimal" {
				t.Errorf("oracle u6 DATA_TYPE: got %q, want decimal", t6)
			}
			// NUMERIC_SCALE may be int64 or []byte.
			scale := 0
			switch x := s6.(type) {
			case int64:
				scale = int(x)
			case []byte:
				for _, c := range x {
					if c >= '0' && c <= '9' {
						scale = scale*10 + int(c-'0')
					}
				}
			}
			if scale != 6 {
				t.Errorf("oracle u6 NUMERIC_SCALE: got %d, want 6", scale)
			}
		}
		row = mc.db.QueryRowContext(mc.ctx,
			`SELECT DATA_TYPE FROM information_schema.COLUMNS
			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v' AND COLUMN_NAME='u0'`)
		if err := row.Scan(&t0); err != nil {
			t.Errorf("oracle u0 type scan: %v", err)
		} else if strings.ToLower(t0) != "bigint" {
			t.Errorf("oracle u0 DATA_TYPE: got %q, want bigint", t0)
		}

		// omni side: parse the same view. If omni cannot derive a type for
		// UNIX_TIMESTAMP at all, that is a gap documented in c16.md.
		c := scenarioNewCatalog(t)
		results, parseErr := c.Exec(
			`CREATE VIEW v AS SELECT UNIX_TIMESTAMP() AS u0, UNIX_TIMESTAMP(NOW(6)) AS u6`, nil)
		var omniErr error
		if parseErr != nil {
			omniErr = parseErr
		} else {
			for _, r := range results {
				if r.Error != nil {
					omniErr = r.Error
					break
				}
			}
		}
		if omniErr != nil {
			t.Errorf("omni: KNOWN GAP — UNIX_TIMESTAMP CREATE VIEW failed: %v", omniErr)
		}
	})

	// ---- 16.8 DATETIME(N) DEFAULT NOW() fsp mismatch → ER_INVALID_DEFAULT -
	t.Run("16_8_Datetime_fsp_mismatch_DEFAULT", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		cases := []string{
			`CREATE TABLE t (a DATETIME(6) DEFAULT NOW())`,
			`CREATE TABLE t (a DATETIME(3) DEFAULT NOW(6))`,
		}
		for _, ddl := range cases {
			scenarioReset(t, mc)
			c2 := scenarioNewCatalog(t)
			oracleErr, omniErr := c16RunExpectError(t, c2, ddl)
			if oracleErr == nil {
				t.Errorf("oracle: %q unexpectedly accepted", ddl)
			} else if !strings.Contains(oracleErr.Error(), "1067") &&
				!strings.Contains(strings.ToLower(oracleErr.Error()), "invalid default") {
				t.Errorf("oracle: %q unexpected error: %v", ddl, oracleErr)
			}
			if omniErr == nil {
				t.Errorf("omni: KNOWN BUG — %q should error (fsp mismatch), see scenarios_bug_queue/c16.md", ddl)
			}
		}

		// And the valid pair must succeed on both.
		scenarioReset(t, mc)
		c3 := scenarioNewCatalog(t)
		runOnBoth(t, mc, c3, `CREATE TABLE t (a DATETIME(6) DEFAULT NOW(6))`)
		_ = c
	})

	// ---- 16.9 ON UPDATE NOW(N) must match column fsp ---------------------
	t.Run("16_9_On_update_fsp_mismatch", func(t *testing.T) {
		// DATETIME(6) ON UPDATE NOW() → fsp mismatch (0 vs 6).
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)
		ddl := `CREATE TABLE t (a DATETIME(6) DEFAULT NOW(6) ON UPDATE NOW())`
		oracleErr, omniErr := c16RunExpectError(t, c, ddl)
		if oracleErr == nil {
			t.Errorf("oracle: %q unexpectedly accepted", ddl)
		}
		if omniErr == nil {
			t.Errorf("omni: KNOWN BUG — ON UPDATE NOW() fsp 0 on DATETIME(6) should error, see scenarios_bug_queue/c16.md")
		}

		// DATE ON UPDATE NOW() → not allowed (ON UPDATE only on TIMESTAMP/DATETIME).
		scenarioReset(t, mc)
		c2 := scenarioNewCatalog(t)
		ddl2 := `CREATE TABLE t (a DATE ON UPDATE NOW())`
		oErr, mErr := c16RunExpectError(t, c2, ddl2)
		if oErr == nil {
			t.Errorf("oracle: %q unexpectedly accepted", ddl2)
		}
		if mErr == nil {
			t.Errorf("omni: KNOWN BUG — DATE ON UPDATE NOW() should error, see scenarios_bug_queue/c16.md")
		}
	})

	// ---- 16.10 DATETIME storage & SHOW CREATE omits (0) ------------------
	t.Run("16_10_Datetime_fsp_round_trip", func(t *testing.T) {
		// Verify SHOW CREATE renders DATETIME (no suffix) for fsp=0, DATETIME(3)
		// for fsp=3, and the omni catalog preserves the same column type.
		type tc struct {
			decl     string
			wantShow string
			wantPrec int
		}
		cases := []tc{
			{"DATETIME", "datetime", 0},
			{"DATETIME(0)", "datetime", 0},
			{"DATETIME(3)", "datetime(3)", 3},
			{"DATETIME(6)", "datetime(6)", 6},
			{"TIMESTAMP(6)", "timestamp(6)", 6},
		}
		for _, k := range cases {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			// Use explicit_defaults_for_timestamp so TIMESTAMP doesn't pick up promotion.
			if _, err := mc.db.ExecContext(mc.ctx, `SET SESSION explicit_defaults_for_timestamp=1`); err != nil {
				t.Errorf("oracle SET: %v", err)
			}
			ddl := `CREATE TABLE t (a ` + k.decl + ` NULL)`
			runOnBoth(t, mc, c, ddl)

			show := strings.ToLower(oracleShow(t, mc, "SHOW CREATE TABLE t"))
			if !strings.Contains(show, "`a` "+k.wantShow) {
				t.Errorf("oracle SHOW CREATE for %q: missing %q in %s", k.decl, k.wantShow, show)
			}
			// DATETIME(0) must NOT appear.
			if strings.Contains(show, "datetime(0)") || strings.Contains(show, "timestamp(0)") {
				t.Errorf("oracle SHOW CREATE for %q: (0) suffix not elided: %s", k.decl, show)
			}

			if p, ok := c16OracleDatetimePrecision(t, "a"); !ok {
				t.Errorf("oracle DATETIME_PRECISION NULL for %q", k.decl)
			} else if p != k.wantPrec {
				t.Errorf("oracle DATETIME_PRECISION for %q: got %d, want %d", k.decl, p, k.wantPrec)
			}

			// omni: column type should match the rendered form.
			tbl := c.GetDatabase("testdb").GetTable("t")
			if tbl == nil {
				t.Errorf("omni: table missing for %q", k.decl)
				continue
			}
			col := tbl.GetColumn("a")
			if col == nil {
				t.Errorf("omni: column a missing for %q", k.decl)
				continue
			}
			omniType := strings.ToLower(col.ColumnType)
			if omniType == "" {
				omniType = strings.ToLower(col.DataType)
			}
			if !strings.Contains(omniType, k.wantShow) {
				t.Errorf("omni: ColumnType=%q for %q, want containing %q", omniType, k.decl, k.wantShow)
			}
			if strings.Contains(omniType, "(0)") {
				t.Errorf("omni: KNOWN BUG — %q kept (0) suffix: %q", k.decl, omniType)
			}
		}
	})

	// ---- 16.11 YEAR(N) deprecation — only YEAR(4) accepted ----------------
	t.Run("16_11_YEAR_normalization", func(t *testing.T) {
		// YEAR(2), YEAR(3), YEAR(5) → ER_INVALID_YEAR_COLUMN_LENGTH (1818).
		for _, n := range []string{"2", "3", "5"} {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			ddl := `CREATE TABLE t (y YEAR(` + n + `))`
			oErr, mErr := c16RunExpectError(t, c, ddl)
			if oErr == nil {
				t.Errorf("oracle: YEAR(%s) unexpectedly accepted", n)
			}
			if mErr == nil {
				t.Errorf("omni: KNOWN BUG — YEAR(%s) should error (ER_INVALID_YEAR_COLUMN_LENGTH), see scenarios_bug_queue/c16.md", n)
			}
		}

		// YEAR(4) normalized to YEAR in SHOW CREATE.
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)
		runOnBoth(t, mc, c, `CREATE TABLE t (y YEAR(4))`)
		show := strings.ToLower(oracleShow(t, mc, "SHOW CREATE TABLE t"))
		if !strings.Contains(show, "`y` year") || strings.Contains(show, "year(4)") {
			t.Errorf("oracle SHOW CREATE: YEAR(4) not normalized: %s", show)
		}
		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl != nil {
			if col := tbl.GetColumn("y"); col != nil {
				omniType := strings.ToLower(col.ColumnType)
				if omniType == "" {
					omniType = strings.ToLower(col.DataType)
				}
				if strings.Contains(omniType, "year(4)") || strings.Contains(omniType, "year(2)") {
					t.Errorf("omni: KNOWN BUG — YEAR(4) not normalized: %q", omniType)
				}
			}
		}
	})

	// ---- 16.12 TIMESTAMP first-column promotion carries column fsp -------
	//
	// NOTE: asymmetric scenario. The session variable
	// `explicit_defaults_for_timestamp=0` only affects the MySQL oracle —
	// omni has no session-variable model today. The omni-side assertions
	// below are tagged "KNOWN GAP" and are expected to fail in either
	// direction: omni's promotion path either doesn't honor the oracle's
	// session state at all (today's behavior) or eventually will honor
	// session vars and then match automatically. If omni starts tracking
	// session vars, revisit this test to mirror the SET on both sides.
	t.Run("16_12_Timestamp_promotion_fsp", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)
		// Default explicit_defaults_for_timestamp=ON on MySQL 8.0, so turn it
		// off to trigger implicit promotion on the oracle side.
		if _, err := mc.db.ExecContext(mc.ctx, `SET SESSION explicit_defaults_for_timestamp=0`); err != nil {
			t.Errorf("oracle SET: %v", err)
		}
		runOnBoth(t, mc, c, `CREATE TABLE t (ts TIMESTAMP(3) NOT NULL)`)

		show := strings.ToLower(oracleShow(t, mc, "SHOW CREATE TABLE t"))
		if !strings.Contains(show, "timestamp(3)") {
			t.Errorf("oracle SHOW CREATE: missing timestamp(3): %s", show)
		}
		if !strings.Contains(show, "default current_timestamp(3)") {
			t.Errorf("oracle SHOW CREATE: missing DEFAULT CURRENT_TIMESTAMP(3): %s", show)
		}
		if !strings.Contains(show, "on update current_timestamp(3)") {
			t.Errorf("oracle SHOW CREATE: missing ON UPDATE CURRENT_TIMESTAMP(3): %s", show)
		}

		// omni side: check promotion produced matching fsp on DEFAULT and ON UPDATE.
		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
			return
		}
		col := tbl.GetColumn("ts")
		if col == nil {
			t.Errorf("omni: column ts missing")
			return
		}
		if col.Default == nil {
			t.Errorf("omni: KNOWN GAP — expected DEFAULT CURRENT_TIMESTAMP(3) after first-col TIMESTAMP promotion, see scenarios_bug_queue/c16.md")
		} else {
			lo := strings.ToLower(*col.Default)
			if !strings.Contains(lo, "current_timestamp(3)") && !strings.Contains(lo, "now(3)") {
				t.Errorf("omni: KNOWN BUG — DEFAULT = %q, expected CURRENT_TIMESTAMP(3), see scenarios_bug_queue/c16.md", *col.Default)
			}
		}
		if col.OnUpdate == "" {
			t.Errorf("omni: KNOWN GAP — expected ON UPDATE CURRENT_TIMESTAMP(3) after promotion, see scenarios_bug_queue/c16.md")
		} else if !strings.Contains(strings.ToLower(col.OnUpdate), "current_timestamp(3)") &&
			!strings.Contains(strings.ToLower(col.OnUpdate), "now(3)") {
			t.Errorf("omni: KNOWN BUG — ON UPDATE = %q, expected CURRENT_TIMESTAMP(3)", col.OnUpdate)
		}
	})
}

// itoaC16 is a tiny local helper that avoids importing strconv just for
// TestScenario_C16. Only used with small non-negative integers.
func itoaC16(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [8]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
