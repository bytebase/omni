package catalog

import (
	"strings"
	"testing"
)

// TestScenario_C10 covers Section C10 "View metadata defaults" from
// SCENARIOS-mysql-implicit-behavior.md. Each subtest runs DDL against both
// a real MySQL 8.0 container and the omni catalog, then asserts that both
// agree on the effective default for a given view-metadata behavior.
//
// Failed omni assertions are NOT proof failures — they are recorded in
// mysql/catalog/scenarios_bug_queue/c10.md.
func TestScenario_C10(t *testing.T) {
	scenariosSkipIfShort(t)
	scenariosSkipIfNoDocker(t)

	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	// c10OmniView fetches a view from the omni catalog by name.
	c10OmniView := func(c *Catalog, name string) *View {
		db := c.GetDatabase("testdb")
		if db == nil {
			return nil
		}
		return db.Views[strings.ToLower(name)]
	}

	// c10OracleViewRow returns one row from information_schema.VIEWS for (testdb, name).
	c10OracleViewRow := func(t *testing.T, name string) (definer, securityType, checkOption, isUpdatable string) {
		t.Helper()
		oracleScan(t, mc,
			`SELECT DEFINER, SECURITY_TYPE, CHECK_OPTION, IS_UPDATABLE
             FROM information_schema.VIEWS
             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='`+name+`'`,
			&definer, &securityType, &checkOption, &isUpdatable)
		return
	}

	// --- 10.1 ALGORITHM defaults to UNDEFINED ------------------------------
	t.Run("10_1_algorithm_defaults_undefined", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c,
			`CREATE TABLE t (a INT);
			 CREATE VIEW v AS SELECT a FROM t;`)

		// Oracle: CHECK_OPTION=NONE, SECURITY_TYPE=DEFINER, SHOW CREATE has ALGORITHM=UNDEFINED.
		_, secType, checkOpt, _ := c10OracleViewRow(t, "v")
		if secType != "DEFINER" {
			t.Errorf("oracle SECURITY_TYPE: got %q, want %q", secType, "DEFINER")
		}
		if checkOpt != "NONE" {
			t.Errorf("oracle CHECK_OPTION: got %q, want %q", checkOpt, "NONE")
		}
		showCreate := oracleShow(t, mc, "SHOW CREATE VIEW v")
		if !strings.Contains(strings.ToUpper(showCreate), "ALGORITHM=UNDEFINED") {
			t.Errorf("oracle SHOW CREATE VIEW v: got %q, want contains ALGORITHM=UNDEFINED", showCreate)
		}

		// omni: view object reports Algorithm=UNDEFINED, SqlSecurity=DEFINER, CheckOption=NONE.
		v := c10OmniView(c, "v")
		if v == nil {
			t.Error("omni: view v not found")
			return
		}
		if strings.ToUpper(v.Algorithm) != "UNDEFINED" {
			t.Errorf("omni Algorithm: got %q, want UNDEFINED", v.Algorithm)
		}
		if strings.ToUpper(v.SqlSecurity) != "DEFINER" {
			t.Errorf("omni SqlSecurity: got %q, want DEFINER", v.SqlSecurity)
		}
		// CheckOption may be represented as "" or "NONE" depending on omni; "NONE" semantics
		// means "no WITH CHECK OPTION clause". Accept either.
		if v.CheckOption != "" && strings.ToUpper(v.CheckOption) != "NONE" {
			t.Errorf("omni CheckOption: got %q, want \"\" or NONE", v.CheckOption)
		}
	})

	// --- 10.2 DEFINER defaults to current user -----------------------------
	t.Run("10_2_definer_defaults_current_user", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c,
			`CREATE TABLE t (a INT);
			 CREATE VIEW v AS SELECT a FROM t;`)

		definer, _, _, _ := c10OracleViewRow(t, "v")
		// The container uses root; DEFINER comes back as something like "root@%"
		// (without backticks in I_S but with them in SHOW CREATE).
		if !strings.Contains(strings.ToLower(definer), "root") {
			t.Errorf("oracle DEFINER: got %q, want contains 'root'", definer)
		}

		v := c10OmniView(c, "v")
		if v == nil {
			t.Error("omni: view v not found")
			return
		}
		if v.Definer == "" {
			t.Error("omni Definer: got empty, want non-empty (e.g. `root`@`%`)")
		}
	})

	// --- 10.3 CHECK OPTION default is CASCADED when WITH CHECK OPTION lacks a qualifier ---
	t.Run("10_3_check_option_default_cascaded", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c,
			`CREATE TABLE t (a INT);
			 CREATE VIEW v1 AS SELECT a FROM t WHERE a > 0 WITH CHECK OPTION;
			 CREATE VIEW v2 AS SELECT a FROM t WHERE a > 0 WITH LOCAL CHECK OPTION;`)

		// Oracle v1 → CASCADED, v2 → LOCAL.
		_, _, v1CheckOpt, _ := c10OracleViewRow(t, "v1")
		_, _, v2CheckOpt, _ := c10OracleViewRow(t, "v2")
		if v1CheckOpt != "CASCADED" {
			t.Errorf("oracle v1 CHECK_OPTION: got %q, want CASCADED", v1CheckOpt)
		}
		if v2CheckOpt != "LOCAL" {
			t.Errorf("oracle v2 CHECK_OPTION: got %q, want LOCAL", v2CheckOpt)
		}

		// omni must distinguish three states (NONE/LOCAL/CASCADED) and normalize
		// bare WITH CHECK OPTION → CASCADED.
		v1 := c10OmniView(c, "v1")
		v2 := c10OmniView(c, "v2")
		if v1 == nil || v2 == nil {
			t.Error("omni: v1 or v2 not found")
			return
		}
		if strings.ToUpper(v1.CheckOption) != "CASCADED" {
			t.Errorf("omni v1 CheckOption: got %q, want CASCADED", v1.CheckOption)
		}
		if strings.ToUpper(v2.CheckOption) != "LOCAL" {
			t.Errorf("omni v2 CheckOption: got %q, want LOCAL", v2.CheckOption)
		}
	})

	// --- 10.4 ALGORITHM=UNDEFINED persists at CREATE; MERGE downgrades on non-mergeable ---
	t.Run("10_4_algorithm_undefined_resolution", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c,
			`CREATE TABLE t (a INT);
			 CREATE ALGORITHM=UNDEFINED VIEW v_agg AS SELECT COUNT(*) FROM t;
			 CREATE ALGORITHM=MERGE VIEW v_merge_bad AS SELECT DISTINCT a FROM t;`)

		// Oracle: both views stored with ALGORITHM=UNDEFINED (v_merge_bad gets downgraded).
		showAgg := oracleShow(t, mc, "SHOW CREATE VIEW v_agg")
		if !strings.Contains(strings.ToUpper(showAgg), "ALGORITHM=UNDEFINED") {
			t.Errorf("oracle SHOW CREATE VIEW v_agg: got %q, want contains ALGORITHM=UNDEFINED", showAgg)
		}
		showMergeBad := oracleShow(t, mc, "SHOW CREATE VIEW v_merge_bad")
		if !strings.Contains(strings.ToUpper(showMergeBad), "ALGORITHM=UNDEFINED") {
			t.Errorf("oracle SHOW CREATE VIEW v_merge_bad (post-downgrade): got %q, want contains ALGORITHM=UNDEFINED", showMergeBad)
		}

		// omni: v_agg Algorithm=UNDEFINED; v_merge_bad should record the user-declared
		// value (MERGE) verbatim per SCENARIOS guidance — MySQL silently downgrades
		// but the catalog representation must preserve the pre-downgrade value.
		vAgg := c10OmniView(c, "v_agg")
		if vAgg == nil {
			t.Error("omni: v_agg not found")
		} else if strings.ToUpper(vAgg.Algorithm) != "UNDEFINED" {
			t.Errorf("omni v_agg Algorithm: got %q, want UNDEFINED", vAgg.Algorithm)
		}
		vMergeBad := c10OmniView(c, "v_merge_bad")
		if vMergeBad == nil {
			t.Error("omni: v_merge_bad not found")
		} else if strings.ToUpper(vMergeBad.Algorithm) != "MERGE" {
			// Per scenario: omni must record the declared algorithm, not silently downgrade.
			t.Errorf("omni v_merge_bad Algorithm: got %q, want MERGE (as declared)", vMergeBad.Algorithm)
		}
	})

	// --- 10.5 SQL SECURITY defaults to DEFINER -----------------------------
	t.Run("10_5_sql_security_defaults_definer", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c,
			`CREATE TABLE t (a INT);
			 CREATE VIEW v AS SELECT a FROM t;`)

		var secType string
		oracleScan(t, mc,
			`SELECT SECURITY_TYPE FROM information_schema.VIEWS
             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v'`,
			&secType)
		if secType != "DEFINER" {
			t.Errorf("oracle SECURITY_TYPE: got %q, want DEFINER", secType)
		}

		v := c10OmniView(c, "v")
		if v == nil {
			t.Error("omni: view v not found")
			return
		}
		if strings.ToUpper(v.SqlSecurity) != "DEFINER" {
			t.Errorf("omni SqlSecurity: got %q, want DEFINER (must be defaulted, not empty)", v.SqlSecurity)
		}
	})

	// --- 10.6 View column names default to SELECT expression spelling ------
	t.Run("10_6_view_column_name_derivation", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c,
			`CREATE TABLE t (a INT);
			 CREATE VIEW v_auto AS SELECT a, a+1, COUNT(*) FROM t GROUP BY a;
			 CREATE VIEW v_list (x,y,z) AS SELECT a, a+1, COUNT(*) FROM t GROUP BY a;`)

		// Oracle: v_auto columns are ['a', 'a+1', 'COUNT(*)'] exactly.
		rows := oracleRows(t, mc,
			`SELECT COLUMN_NAME FROM information_schema.COLUMNS
             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v_auto'
             ORDER BY ORDINAL_POSITION`)
		oracleCols := make([]string, 0, len(rows))
		for _, r := range rows {
			if len(r) > 0 {
				oracleCols = append(oracleCols, asString(r[0]))
			}
		}
		wantAuto := []string{"a", "a+1", "COUNT(*)"}
		if strings.Join(oracleCols, ",") != strings.Join(wantAuto, ",") {
			t.Errorf("oracle v_auto columns: got %v, want %v", oracleCols, wantAuto)
		}

		// v_list: explicit column list wins.
		rows2 := oracleRows(t, mc,
			`SELECT COLUMN_NAME FROM information_schema.COLUMNS
             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v_list'
             ORDER BY ORDINAL_POSITION`)
		oracleCols2 := make([]string, 0, len(rows2))
		for _, r := range rows2 {
			if len(r) > 0 {
				oracleCols2 = append(oracleCols2, asString(r[0]))
			}
		}
		wantList := []string{"x", "y", "z"}
		if strings.Join(oracleCols2, ",") != strings.Join(wantList, ",") {
			t.Errorf("oracle v_list columns: got %v, want %v", oracleCols2, wantList)
		}

		// omni: same expectations.
		vAuto := c10OmniView(c, "v_auto")
		if vAuto == nil {
			t.Error("omni: v_auto not found")
		} else {
			if strings.Join(vAuto.Columns, ",") != strings.Join(wantAuto, ",") {
				t.Errorf("omni v_auto Columns: got %v, want %v", vAuto.Columns, wantAuto)
			}
		}
		vList := c10OmniView(c, "v_list")
		if vList == nil {
			t.Error("omni: v_list not found")
		} else {
			if strings.Join(vList.Columns, ",") != strings.Join(wantList, ",") {
				t.Errorf("omni v_list Columns: got %v, want %v", vList.Columns, wantList)
			}
			if !vList.ExplicitColumns {
				t.Error("omni v_list ExplicitColumns: got false, want true")
			}
		}
	})

	// --- 10.7 View updatability is derived from SELECT shape ---------------
	t.Run("10_7_is_updatable_derivation", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c,
			`CREATE TABLE t (a INT);
			 CREATE VIEW v_ok AS SELECT a FROM t;
			 CREATE VIEW v_distinct AS SELECT DISTINCT a FROM t;
			 CREATE ALGORITHM=TEMPTABLE VIEW v_temp AS SELECT a FROM t;`)

		for _, tc := range []struct {
			name string
			want string
		}{
			{"v_ok", "YES"},
			{"v_distinct", "NO"},
			{"v_temp", "NO"},
		} {
			var isUpd string
			oracleScan(t, mc,
				`SELECT IS_UPDATABLE FROM information_schema.VIEWS
                 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='`+tc.name+`'`,
				&isUpd)
			if isUpd != tc.want {
				t.Errorf("oracle %s IS_UPDATABLE: got %q, want %q", tc.name, isUpd, tc.want)
			}

			// omni: the View struct has no IsUpdatable field today — see bug queue.
			// This stanza documents the absence by asserting the view at least
			// exists in the catalog.
			v := c10OmniView(c, tc.name)
			if v == nil {
				t.Errorf("omni: view %s not found", tc.name)
			}
		}
		// Explicit assertion: omni View has no IsUpdatable representation.
		// This is the "declared bug": we want omni to carry IsUpdatable per scenario.
		t.Error("omni: View struct is missing an IsUpdatable field (scenario 10.7 cannot be asserted positively)")
	})

	// --- 10.8 View column nullability widened vs base columns -------------
	t.Run("10_8_outer_join_nullability_widening", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c,
			`CREATE TABLE t1 (id INT NOT NULL, a INT NOT NULL);
			 CREATE TABLE t2 (id INT NOT NULL, b INT NOT NULL);
			 CREATE VIEW v AS SELECT t1.a, t2.b FROM t1 LEFT JOIN t2 ON t1.id = t2.id;`)

		// Oracle (empirical, MySQL 8.0.45): t1.a stays NOT NULL (left/preserved
		// side of LEFT JOIN), t2.b widens to nullable (right/optional side).
		// The SCENARIOS doc text claims `a → YES` but that is incorrect; real
		// MySQL only widens the optional side. We assert against the oracle
		// ground truth.
		rows := oracleRows(t, mc,
			`SELECT COLUMN_NAME, IS_NULLABLE FROM information_schema.COLUMNS
             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v'
             ORDER BY ORDINAL_POSITION`)
		gotNullable := map[string]string{}
		for _, r := range rows {
			if len(r) < 2 {
				continue
			}
			gotNullable[asString(r[0])] = asString(r[1])
		}
		if gotNullable["a"] != "NO" {
			t.Errorf("oracle view column a IS_NULLABLE: got %q, want NO (left/preserved side)", gotNullable["a"])
		}
		if gotNullable["b"] != "YES" {
			t.Errorf("oracle view column b IS_NULLABLE: got %q, want YES (right/optional side)", gotNullable["b"])
		}

		// omni: omni's View struct does not carry per-column nullability info.
		// The column list (v.Columns) is just names. This is the "declared bug":
		// omni view column resolver must track outer-join nullability.
		v := c10OmniView(c, "v")
		if v == nil {
			t.Error("omni: view v not found")
			return
		}
		t.Error("omni: View struct has no per-column nullability; scenario 10.8 cannot be asserted positively")
	})
}
