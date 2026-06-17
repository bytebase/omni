package catalog

import (
	"strings"
	"testing"
)

// TestScenario_C11 covers Section C11 "Trigger defaults" from
// SCENARIOS-mysql-implicit-behavior.md. Each subtest runs DDL on both a real
// MySQL 8.0 container and the omni catalog, then asserts they agree on
// trigger metadata defaults: DEFINER, SQL SECURITY (no INVOKER option),
// charset/collation snapshot, ACTION_ORDER sequencing, NEW/OLD pseudo-row
// access rules, and trigger-on-partitioned-table survival across partition
// mutations.
//
// Failed omni assertions are NOT proof failures — they are recorded in
// mysql/catalog/scenarios_bug_queue/c11.md.
func TestScenario_C11(t *testing.T) {
	scenariosSkipIfShort(t)
	scenariosSkipIfNoDocker(t)

	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	// c11OmniExec runs a multi-statement DDL on the omni catalog and returns
	// (errored, firstErr). A parse error or any per-statement Error flips
	// the bool. Used by scenarios that expect omni to reject DDL.
	c11OmniExec := func(c *Catalog, ddl string) (bool, error) {
		results, err := c.Exec(ddl, nil)
		if err != nil {
			return true, err
		}
		for _, r := range results {
			if r.Error != nil {
				return true, r.Error
			}
		}
		return false, nil
	}

	// --- 11.1 Trigger DEFINER defaults to current user ---------------------
	t.Run("11_1_trigger_definer_defaults_to_current_user", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (a INT);
CREATE TRIGGER trg BEFORE INSERT ON t FOR EACH ROW SET NEW.a = NEW.a;`)

		// Oracle: DEFINER populated with session user (typically `root`@`%`
		// in the test container).
		var definer string
		oracleScan(t, mc,
			`SELECT DEFINER FROM information_schema.TRIGGERS
             WHERE TRIGGER_SCHEMA='testdb' AND TRIGGER_NAME='trg'`,
			&definer)
		if definer == "" {
			t.Errorf("oracle: DEFINER should be non-empty for trg")
		}

		// omni: trigger stored with non-empty Definer.
		db := c.GetDatabase("testdb")
		if db == nil {
			t.Error("omni: testdb missing")
			return
		}
		trg := db.Triggers[toLower("trg")]
		if trg == nil {
			t.Error("omni: trigger trg missing from Triggers map")
			return
		}
		if trg.Definer == "" {
			t.Errorf("omni: trigger trg Definer should default to a session user, got empty")
		}
	})

	// --- 11.2 Trigger SQL SECURITY always DEFINER; no INVOKER option -------
	t.Run("11_2_trigger_sql_security_always_definer", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		// Valid form: CREATE DEFINER=... TRIGGER must succeed on both sides.
		good := `CREATE TABLE t (a INT);
CREATE DEFINER='root'@'%' TRIGGER trg1 BEFORE INSERT ON t FOR EACH ROW SET NEW.a=1;`
		runOnBoth(t, mc, c, good)

		// Invalid form: SQL SECURITY INVOKER on a trigger is a grammar error.
		bad := `CREATE TRIGGER trg2 SQL SECURITY INVOKER BEFORE INSERT ON t FOR EACH ROW SET NEW.a=1`
		_, oracleErr := mc.db.ExecContext(mc.ctx, bad)
		if oracleErr == nil {
			t.Errorf("oracle: expected ER_PARSE_ERROR for SQL SECURITY INVOKER on trigger, got nil")
		}

		// omni: must also reject. A permissive parse here is a bug.
		omniErrored, _ := c11OmniExec(c, bad+";")
		assertBoolEq(t, "omni rejects SQL SECURITY INVOKER on trigger", omniErrored, true)

		// information_schema.TRIGGERS must NOT expose a SECURITY_TYPE column
		// for triggers (unlike VIEWS / ROUTINES). If MySQL ever added one we
		// would need to rethink this scenario; assert absence empirically.
		var colCount int64
		oracleScan(t, mc,
			`SELECT COUNT(*) FROM information_schema.COLUMNS
             WHERE TABLE_SCHEMA='information_schema' AND TABLE_NAME='TRIGGERS'
             AND COLUMN_NAME='SECURITY_TYPE'`,
			&colCount)
		if colCount != 0 {
			t.Errorf("oracle: information_schema.TRIGGERS has SECURITY_TYPE (unexpected), count=%d", colCount)
		}
	})

	// --- 11.3 charset/collation snapshot at trigger creation time ----------
	t.Run("11_3_charset_collation_snapshot", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `SET NAMES utf8mb4;
CREATE TABLE t (a INT);
CREATE TRIGGER trg BEFORE INSERT ON t FOR EACH ROW SET NEW.a = NEW.a;`)

		var cset, ccoll, dbcoll string
		oracleScan(t, mc,
			`SELECT CHARACTER_SET_CLIENT, COLLATION_CONNECTION, DATABASE_COLLATION
             FROM information_schema.TRIGGERS
             WHERE TRIGGER_SCHEMA='testdb' AND TRIGGER_NAME='trg'`,
			&cset, &ccoll, &dbcoll)
		if cset == "" || ccoll == "" || dbcoll == "" {
			t.Errorf("oracle: expected three non-empty snapshot fields, got (%q,%q,%q)",
				cset, ccoll, dbcoll)
		}
		if !strings.HasPrefix(strings.ToLower(cset), "utf8") {
			t.Errorf("oracle: CHARACTER_SET_CLIENT for trg should be utf8*; got %q", cset)
		}

		db := c.GetDatabase("testdb")
		if db == nil {
			t.Error("omni: testdb missing")
			return
		}
		trg := db.Triggers[toLower("trg")]
		if trg == nil {
			t.Error("omni: trigger trg missing")
			return
		}
		// Assert the core fields omni does track so the subtest has real
		// substance, then flag the charset-snapshot gap explicitly. If a
		// future patch adds CharacterSetClient / CollationConnection /
		// DatabaseCollation fields, tighten the assertion below.
		if strings.ToUpper(trg.Timing) != "BEFORE" {
			t.Errorf("omni 11.3: trg.Timing=%q, want BEFORE", trg.Timing)
		}
		if strings.ToUpper(trg.Event) != "INSERT" {
			t.Errorf("omni 11.3: trg.Event=%q, want INSERT", trg.Event)
		}
		if trg.Table != "t" {
			t.Errorf("omni 11.3: trg.Table=%q, want t", trg.Table)
		}
		if !strings.HasPrefix(strings.ToLower(trg.CharacterSetClient), "utf8") {
			t.Errorf("omni 11.3: CharacterSetClient=%q, want utf8*", trg.CharacterSetClient)
		}
		if trg.CollationConnection == "" {
			t.Errorf("omni 11.3: CollationConnection is empty")
		}
		assertStringEq(t, "omni 11.3 DatabaseCollation",
			strings.ToLower(trg.DatabaseCollation), strings.ToLower(db.Collation))
	})

	// --- 11.4 ACTION_ORDER default sequencing within (table, timing, event) -
	t.Run("11_4_action_order_default_sequencing", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		ddl := `CREATE TABLE t (a INT);
CREATE TRIGGER trg_a BEFORE INSERT ON t FOR EACH ROW SET NEW.a = NEW.a + 1;
CREATE TRIGGER trg_b BEFORE INSERT ON t FOR EACH ROW SET NEW.a = NEW.a + 2;
CREATE TRIGGER trg_c BEFORE INSERT ON t FOR EACH ROW PRECEDES trg_a SET NEW.a = NEW.a + 10;`
		runOnBoth(t, mc, c, ddl)

		// Oracle: expect trg_c=1, trg_a=2, trg_b=3 after the PRECEDES splice.
		rows := oracleRows(t, mc,
			`SELECT TRIGGER_NAME, ACTION_ORDER FROM information_schema.TRIGGERS
             WHERE TRIGGER_SCHEMA='testdb' ORDER BY ACTION_ORDER`)
		if len(rows) != 3 {
			t.Errorf("oracle: expected 3 triggers, got %d", len(rows))
		} else {
			wantOrder := []struct {
				name  string
				order int64
			}{
				{"trg_c", 1},
				{"trg_a", 2},
				{"trg_b", 3},
			}
			for i, w := range wantOrder {
				name := asString(rows[i][0])
				var ord int64
				switch v := rows[i][1].(type) {
				case int64:
					ord = v
				case int32:
					ord = int64(v)
				case int:
					ord = int64(v)
				}
				if name != w.name || ord != w.order {
					t.Errorf("oracle row %d: got (%q, %d), want (%q, %d)",
						i, name, ord, w.name, w.order)
				}
			}
		}

		// omni: Trigger struct has no ActionOrder field. The best we can
		// check today is that all three trigger objects exist and the
		// Order info (FOLLOWS/PRECEDES) was captured for trg_c.
		db := c.GetDatabase("testdb")
		if db == nil {
			t.Error("omni: testdb missing")
			return
		}
		for _, name := range []string{"trg_a", "trg_b", "trg_c"} {
			if db.Triggers[toLower(name)] == nil {
				t.Errorf("omni: trigger %s missing", name)
			}
		}
		if trgC := db.Triggers[toLower("trg_c")]; trgC != nil {
			if trgC.Order == nil {
				t.Errorf("omni: trigger trg_c should have Order info for PRECEDES trg_a")
			} else {
				assertBoolEq(t, "omni trg_c Order.Follows (false means PRECEDES)",
					trgC.Order.Follows, false)
				assertStringEq(t, "omni trg_c Order.TriggerName",
					strings.ToLower(trgC.Order.TriggerName), "trg_a")
			}
		}
	})

	// --- 11.5 NEW/OLD pseudo-row access by event type ----------------------
	t.Run("11_5_new_old_pseudorow_by_event", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		// Legal cases: NEW in BEFORE INSERT, OLD in BEFORE DELETE, both in UPDATE.
		legal := `CREATE TABLE t (a INT);
CREATE TRIGGER t_ins BEFORE INSERT ON t FOR EACH ROW SET NEW.a = NEW.a + 1;
CREATE TRIGGER t_del BEFORE DELETE ON t FOR EACH ROW SET @x = OLD.a;
CREATE TRIGGER t_upd BEFORE UPDATE ON t FOR EACH ROW SET NEW.a = OLD.a + 1;`
		runOnBoth(t, mc, c, legal)

		// Illegal: OLD inside an INSERT trigger -> ER_TRG_NO_SUCH_ROW_IN_TRG.
		badOldInInsert := `CREATE TRIGGER bad1 AFTER INSERT ON t FOR EACH ROW SET @x = OLD.a`
		_, oErr1 := mc.db.ExecContext(mc.ctx, badOldInInsert)
		if oErr1 == nil {
			t.Errorf("oracle: expected rejection of OLD.a in INSERT trigger, got nil")
		}
		omniErr1, _ := c11OmniExec(c, badOldInInsert+";")
		assertBoolEq(t, "omni rejects OLD.* in INSERT trigger body", omniErr1, true)

		// Illegal: NEW inside a DELETE trigger -> ER_TRG_NO_SUCH_ROW_IN_TRG.
		badNewInDelete := `CREATE TRIGGER bad2 AFTER DELETE ON t FOR EACH ROW SET @x = NEW.a`
		_, oErr2 := mc.db.ExecContext(mc.ctx, badNewInDelete)
		if oErr2 == nil {
			t.Errorf("oracle: expected rejection of NEW.a in DELETE trigger, got nil")
		}
		omniErr2, _ := c11OmniExec(c, badNewInDelete+";")
		assertBoolEq(t, "omni rejects NEW.* in DELETE trigger body", omniErr2, true)

		// Illegal: SET NEW.a in AFTER INSERT — NEW is read-only after the row
		// is written.
		badAfterAssign := `CREATE TRIGGER bad3 AFTER INSERT ON t FOR EACH ROW SET NEW.a = 99`
		_, oErr3 := mc.db.ExecContext(mc.ctx, badAfterAssign)
		if oErr3 == nil {
			t.Errorf("oracle: expected rejection of SET NEW.a in AFTER INSERT trigger, got nil")
		}
		omniErr3, _ := c11OmniExec(c, badAfterAssign+";")
		assertBoolEq(t, "omni rejects writes to NEW.* in AFTER trigger", omniErr3, true)

		// Sanity: legal triggers landed in omni and oracle.
		db := c.GetDatabase("testdb")
		if db != nil {
			for _, name := range []string{"t_ins", "t_del", "t_upd"} {
				if db.Triggers[toLower(name)] == nil {
					t.Errorf("omni: legal trigger %s missing", name)
				}
			}
		}
		var legalCount int64
		oracleScan(t, mc,
			`SELECT COUNT(*) FROM information_schema.TRIGGERS
             WHERE TRIGGER_SCHEMA='testdb' AND TRIGGER_NAME IN ('t_ins','t_del','t_upd')`,
			&legalCount)
		if legalCount != 3 {
			t.Errorf("oracle: expected 3 legal triggers, got %d", legalCount)
		}
	})

	// --- 11.6 Trigger on partitioned table survives partition mutation -----
	t.Run("11_6_trigger_on_partitioned_table", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		ddl := `CREATE TABLE t (a INT, b INT) PARTITION BY HASH(a) PARTITIONS 4;
CREATE TRIGGER trg BEFORE INSERT ON t FOR EACH ROW SET NEW.b = NEW.a * 2;`
		runOnBoth(t, mc, c, ddl)

		// Alter the partition layout. Oracle: trigger survives.
		alter := `ALTER TABLE t COALESCE PARTITION 2`
		if _, err := mc.db.ExecContext(mc.ctx, alter); err != nil {
			t.Errorf("oracle: COALESCE PARTITION failed: %v", err)
		}
		// omni: execute the same ALTER; if omni's parser does not support
		// this syntax yet, record it as part of the bug queue but don't
		// abort the subtest.
		if _, err := c.Exec(alter+";", nil); err != nil {
			t.Logf("omni: ALTER TABLE ... COALESCE PARTITION not yet supported: %v", err)
		}

		var trgCount int64
		oracleScan(t, mc,
			`SELECT COUNT(*) FROM information_schema.TRIGGERS
             WHERE TRIGGER_SCHEMA='testdb' AND EVENT_OBJECT_TABLE='t' AND TRIGGER_NAME='trg'`,
			&trgCount)
		if trgCount != 1 {
			t.Errorf("oracle: trigger trg should survive COALESCE PARTITION, count=%d", trgCount)
		}

		// omni: trigger must still be registered against the table, not
		// hanging off any per-partition structure.
		db := c.GetDatabase("testdb")
		if db == nil {
			t.Error("omni: testdb missing")
			return
		}
		trg := db.Triggers[toLower("trg")]
		if trg == nil {
			t.Errorf("omni: trigger trg missing after partition mutation")
			return
		}
		assertStringEq(t, "omni trigger trg.Table", strings.ToLower(trg.Table), "t")
	})
}
