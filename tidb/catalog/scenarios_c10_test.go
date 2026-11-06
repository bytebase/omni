package catalog

import (
	"strings"
	"testing"
)

// TestScenario_C10 covers Section C10 "View metadata defaults" from
// mysql/catalog/SCENARIOS-mysql-implicit-behavior.md. Each subtest runs DDL against both
// TiDB v8.5.5 and the omni catalog, then asserts that both agree on the
// effective default for a given view-metadata behavior. Where TiDB answers
// differently from MySQL -- CHECK_OPTION, IS_UPDATABLE, ALGORITHM, column
// nullability -- the subtest asserts TiDB's answer.
func TestScenario_C10(t *testing.T) {
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

		// Oracle: SECURITY_TYPE=DEFINER, SHOW CREATE has ALGORITHM=UNDEFINED.
		// TiDB reports CHECK_OPTION=CASCADED for a view written without the
		// clause (MySQL reports NONE): its view model has no "none", and
		// SHOW CREATE VIEW never writes the clause either way.
		_, secType, checkOpt, _ := c10OracleViewRow(t, "v")
		if secType != "DEFINER" {
			t.Errorf("oracle SECURITY_TYPE: got %q, want %q", secType, "DEFINER")
		}
		if checkOpt != "CASCADED" {
			t.Errorf("oracle CHECK_OPTION: got %q, want %q", checkOpt, "CASCADED")
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
		// CheckOption records the clause as written: "" when there is none.
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

	// --- 10.3 CHECK OPTION is CASCADED or LOCAL as written ---
	// TiDB has no bare WITH CHECK OPTION (1064); refused_forms_test.go pins that.
	t.Run("10_3_check_option_cascaded_or_local", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c,
			`CREATE TABLE t (a INT);
			 CREATE VIEW v1 AS SELECT a FROM t WHERE a > 0 WITH CASCADED CHECK OPTION;
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

		// omni must distinguish three states (NONE/LOCAL/CASCADED).
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

	// --- 10.4 ALGORITHM is stored as declared ---
	t.Run("10_4_algorithm_undefined_resolution", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c,
			`CREATE TABLE t (a INT);
			 CREATE ALGORITHM=UNDEFINED VIEW v_agg AS SELECT COUNT(*) FROM t;
			 CREATE ALGORITHM=MERGE VIEW v_merge_bad AS SELECT DISTINCT a FROM t;`)

		// Oracle: TiDB stores the algorithm declared. MySQL downgrades MERGE
		// on a view that cannot merge (DISTINCT) to UNDEFINED; TiDB keeps
		// ALGORITHM=MERGE.
		showAgg := oracleShow(t, mc, "SHOW CREATE VIEW v_agg")
		if !strings.Contains(strings.ToUpper(showAgg), "ALGORITHM=UNDEFINED") {
			t.Errorf("oracle SHOW CREATE VIEW v_agg: got %q, want contains ALGORITHM=UNDEFINED", showAgg)
		}
		showMergeBad := oracleShow(t, mc, "SHOW CREATE VIEW v_merge_bad")
		if !strings.Contains(strings.ToUpper(showMergeBad), "ALGORITHM=MERGE") {
			t.Errorf("oracle SHOW CREATE VIEW v_merge_bad: got %q, want contains ALGORITHM=MERGE", showMergeBad)
		}
		for name, want := range map[string]string{"v_agg": showAgg, "v_merge_bad": showMergeBad} {
			if got := c.ShowCreateView("testdb", name); got != want {
				t.Errorf("omni SHOW CREATE VIEW %s:\n got: %s\nwant: %s", name, got, want)
			}
		}

		// omni: v_agg Algorithm=UNDEFINED; v_merge_bad records MERGE as declared.
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

		// TiDB reports IS_UPDATABLE=NO for every view, whatever its shape
		// (MySQL derives YES for v_ok): it does not derive updatability, so
		// the catalog has none to model either.
		for _, tc := range []struct {
			name string
			want string
		}{
			{"v_ok", "NO"},
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

			v := c10OmniView(c, tc.name)
			if v == nil {
				t.Errorf("omni: view %s not found", tc.name)
			}
		}
	})

	// --- 10.8 View column nullability widened vs base columns -------------
	t.Run("10_8_outer_join_nullability_widening", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c,
			`CREATE TABLE t1 (id INT NOT NULL, a INT NOT NULL);
			 CREATE TABLE t2 (id INT NOT NULL, b INT NOT NULL);
			 CREATE VIEW v AS SELECT t1.a, t2.b FROM t1 LEFT JOIN t2 ON t1.id = t2.id;`)

		// Oracle: TiDB reports every view column nullable, the preserved
		// side of the LEFT JOIN included (MySQL keeps t1.a NOT NULL). A view
		// column carries no nullability in TiDB, so the catalog's View, whose
		// columns are names only, has none to model.
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
		if gotNullable["a"] != "YES" {
			t.Errorf("oracle view column a IS_NULLABLE: got %q, want YES", gotNullable["a"])
		}
		if gotNullable["b"] != "YES" {
			t.Errorf("oracle view column b IS_NULLABLE: got %q, want YES (right/optional side)", gotNullable["b"])
		}

		v := c10OmniView(c, "v")
		if v == nil {
			t.Error("omni: view v not found")
			return
		}
		if strings.Join(v.Columns, ",") != "a,b" {
			t.Errorf("omni v Columns: got %v, want [a b]", v.Columns)
		}
	})
}
