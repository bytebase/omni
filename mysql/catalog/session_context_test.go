package catalog

import (
	"strings"
	"testing"
)

// Unit coverage for BYT-9832: on a MySQL declarative RECREATE (DROP+CREATE) of a
// routine/trigger/event, the emitted CREATE must be wrapped in a save/restore of the OLD
// object's session context (sql_mode / charset / collation / — for events — time_zone),
// sourced out-of-band via ApplySessionContext. A brand-new object emits a bare CREATE, and
// the session context is NOT part of the diff identity (a mode-only difference never
// triggers a change).
//
// These tests build catalogs the same way the other unit tests do (LoadSQL under a fixed
// database), stamp the source objects with ApplySessionContext, then assert the generated
// plan's framing exactly.

const scDB = "d"

func loadSCCatalog(t *testing.T, ddl string) *Catalog {
	t.Helper()
	c, err := LoadSQL("CREATE DATABASE " + scDB + " DEFAULT CHARSET=utf8mb4;\nUSE " + scDB + ";\n" + ddl)
	if err != nil {
		t.Fatalf("load %q: %v", ddl, err)
	}
	return c
}

// planSQL diffs from→to and returns the full migration SQL.
func planSQL(t *testing.T, from, to *Catalog, n *Normalizer) string {
	t.Helper()
	d := DiffWithNormalizer(from, to, n)
	return GenerateMigrationWithNormalizer(from, to, d, n).SQL()
}

// TestRecreate_RoutineCarriesSQLMode covers requirement (1): a routine authored under a
// non-default sql_mode whose body is edited is DROPped and re-CREATEd wrapped in an exact
// SET sql_mode='<mode>' … restore.
func TestRecreate_RoutineCarriesSQLMode(t *testing.T) {
	n := NormalizerFor(MySQL80)
	from := loadSCCatalog(t, "CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC RETURN a + 1")
	to := loadSCCatalog(t, "CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC RETURN a + 2")

	from.ApplySessionContext(SessionContextMap{
		Functions: map[string]SessionContext{
			"f": {
				SQLMode:             "PIPES_AS_CONCAT,NO_BACKSLASH_ESCAPES",
				CharacterSetClient:  "utf8mb4",
				CollationConnection: "utf8mb4_0900_ai_ci",
			},
		},
	})

	sql := planSQL(t, from, to, n)

	// It must be a DROP then a framed CREATE.
	if !strings.Contains(sql, "DROP FUNCTION") {
		t.Fatalf("expected DROP FUNCTION in plan:\n%s", sql)
	}
	// Exact save/restore framing around the CREATE.
	mustContainInOrder(t, sql,
		"SET @saved_sql_mode = @@sql_mode",
		"SET sql_mode = 'PIPES_AS_CONCAT,NO_BACKSLASH_ESCAPES'",
		"SET @saved_cs_client = @@character_set_client",
		"SET character_set_client = utf8mb4",
		"SET @saved_coll_conn = @@collation_connection",
		"SET collation_connection = utf8mb4_0900_ai_ci",
		"CREATE FUNCTION `d`.`f`",
		"RETURN a + 2",
		"SET collation_connection = @saved_coll_conn",
		"SET character_set_client = @saved_cs_client",
		"SET sql_mode = @saved_sql_mode",
	)
}

// TestNewRoutine_BareCreate covers requirement (2): a brand-new routine (no From) emits a
// bare CREATE with no session framing.
func TestNewRoutine_BareCreate(t *testing.T) {
	n := NormalizerFor(MySQL80)
	from := loadSCCatalog(t, "CREATE FUNCTION other() RETURNS INT DETERMINISTIC RETURN 1")
	to := loadSCCatalog(t, "CREATE FUNCTION other() RETURNS INT DETERMINISTIC RETURN 1;\nCREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC RETURN a")

	// Even if a stray context map is supplied, an ADD has no From, so no framing is emitted.
	from.ApplySessionContext(SessionContextMap{
		Functions: map[string]SessionContext{
			"f": {SQLMode: "PIPES_AS_CONCAT", CharacterSetClient: "utf8mb4"},
		},
	})

	sql := planSQL(t, from, to, n)
	if !strings.Contains(sql, "CREATE FUNCTION `d`.`f`") {
		t.Fatalf("expected CREATE FUNCTION f:\n%s", sql)
	}
	if strings.Contains(sql, "sql_mode") || strings.Contains(sql, "@saved_") {
		t.Fatalf("brand-new routine must emit a BARE CREATE (no session framing):\n%s", sql)
	}
	if strings.Contains(sql, "DROP FUNCTION `f`") {
		t.Fatalf("brand-new routine must not DROP:\n%s", sql)
	}
}

// TestUnchangedRoutine_NoDiff covers requirement (3): an unchanged routine (same body, only
// the carried session context differs) produces NO diff entry — the identity is
// context-agnostic, so there is no phantom churn.
func TestUnchangedRoutine_NoDiff(t *testing.T) {
	n := NormalizerFor(MySQL80)
	body := "CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC RETURN a"
	from := loadSCCatalog(t, body)
	to := loadSCCatalog(t, body)

	// Only the source carries context; a mode difference alone must not trigger a change.
	from.ApplySessionContext(SessionContextMap{
		Functions: map[string]SessionContext{
			"f": {SQLMode: "PIPES_AS_CONCAT,NO_BACKSLASH_ESCAPES", CharacterSetClient: "latin1", CollationConnection: "latin1_swedish_ci"},
		},
	})

	d := DiffWithNormalizer(from, to, n)
	if !d.IsEmpty() {
		t.Fatalf("carried session context must not create a diff, got: %s", describeRoutineDiff(d))
	}
}

// TestRecreate_RoutineCharsetCollationOnly covers requirement (4): charset/collation are
// carried through a routine recreate even when sql_mode is empty. sql_mode is still framed
// (empty modes are authored, not the server default).
func TestRecreate_RoutineCharsetCollationOnly(t *testing.T) {
	n := NormalizerFor(MySQL80)
	from := loadSCCatalog(t, "CREATE PROCEDURE p(IN a INT) BEGIN SET @v = a; END")
	to := loadSCCatalog(t, "CREATE PROCEDURE p(IN a INT) BEGIN SET @v = a + 1; END")

	from.ApplySessionContext(SessionContextMap{
		Procedures: map[string]SessionContext{
			"p": {
				SQLMode:             "", // authored empty modes must round-trip as SET sql_mode=''
				CharacterSetClient:  "latin1",
				CollationConnection: "latin1_swedish_ci",
			},
		},
	})

	sql := planSQL(t, from, to, n)
	mustContainInOrder(t, sql,
		"SET sql_mode = ''",
		"SET character_set_client = latin1",
		"SET collation_connection = latin1_swedish_ci",
		"CREATE PROCEDURE `d`.`p`",
		"SET collation_connection = @saved_coll_conn",
		"SET character_set_client = @saved_cs_client",
		"SET sql_mode = @saved_sql_mode",
	)
}

// TestRecreate_EventCarriesTimeZoneAndSQLMode covers requirement (5): an event recreated via
// DROP+CREATE (an ENDS-removal change, which forces the recreate path) carries BOTH its
// non-UTC time_zone and its sql_mode.
func TestRecreate_EventCarriesTimeZoneAndSQLMode(t *testing.T) {
	n := NormalizerFor(MySQL80)
	// From has an ENDS bound; To drops it → eventModifyNeedsRecreate → DROP+CREATE.
	from := loadSCCatalog(t, "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR ENDS '2030-01-01 00:00:00' DO SET @x = 1")
	to := loadSCCatalog(t, "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x = 1")

	from.ApplySessionContext(SessionContextMap{
		Events: map[string]SessionContext{
			"e": {
				SQLMode:             "STRICT_ALL_TABLES",
				CharacterSetClient:  "utf8mb4",
				CollationConnection: "utf8mb4_0900_ai_ci",
				TimeZone:            "+08:00",
			},
		},
	})

	// Confirm this actually takes the recreate path.
	d := DiffWithNormalizer(from, to, n)
	if len(d.Events) != 1 || d.Events[0].Action != DiffModify {
		t.Fatalf("want one event MODIFY, got %+v", d.Events)
	}
	if !eventModifyNeedsRecreate(d.Events[0].From, d.Events[0].To) {
		t.Fatalf("ENDS removal should force DROP+CREATE recreate path")
	}

	sql := GenerateMigrationWithNormalizer(from, to, d, n).SQL()
	if !strings.Contains(sql, "DROP EVENT IF EXISTS") {
		t.Fatalf("expected DROP EVENT in recreate plan:\n%s", sql)
	}
	mustContainInOrder(t, sql,
		"SET @saved_sql_mode = @@sql_mode",
		"SET sql_mode = 'STRICT_ALL_TABLES'",
		"SET @saved_time_zone = @@time_zone",
		"SET time_zone = '+08:00'",
		"CREATE EVENT `e`",
		"SET time_zone = @saved_time_zone",
		"SET sql_mode = @saved_sql_mode",
	)
}

// TestAlterEvent_CarriesSQLMode covers the common event-modify path: a body change goes
// through ALTER EVENT (NOT DROP+CREATE), and ALTER EVENT … DO re-stamps sql_mode/charset/
// collation from the executing session (verified on live 8.0.32), so the ALTER must be
// framed with the OLD event's context too.
func TestAlterEvent_CarriesSQLMode(t *testing.T) {
	n := NormalizerFor(MySQL80)
	// Body-only change → ALTER path (no ENDS to remove).
	from := loadSCCatalog(t, "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x = 1")
	to := loadSCCatalog(t, "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x = 2")

	from.ApplySessionContext(SessionContextMap{
		Events: map[string]SessionContext{
			"e": {
				SQLMode:             "PIPES_AS_CONCAT",
				CharacterSetClient:  "latin1",
				CollationConnection: "latin1_swedish_ci",
				TimeZone:            "+08:00",
			},
		},
	})

	d := DiffWithNormalizer(from, to, n)
	if len(d.Events) != 1 || d.Events[0].Action != DiffModify {
		t.Fatalf("want one event MODIFY, got %+v", d.Events)
	}
	// Confirm this takes the ALTER path (not DROP+CREATE).
	if eventModifyNeedsRecreate(d.Events[0].From, d.Events[0].To) {
		t.Fatalf("a body-only change should take the ALTER path, not DROP+CREATE")
	}

	sql := GenerateMigrationWithNormalizer(from, to, d, n).SQL()
	if strings.Contains(sql, "DROP EVENT") {
		t.Fatalf("body-only event change must not DROP:\n%s", sql)
	}
	mustContainInOrder(t, sql,
		"SET @saved_sql_mode = @@sql_mode",
		"SET sql_mode = 'PIPES_AS_CONCAT'",
		"SET character_set_client = latin1",
		"SET collation_connection = latin1_swedish_ci",
		"ALTER EVENT `e`",
		"SET collation_connection = @saved_coll_conn",
		"SET character_set_client = @saved_cs_client",
		"SET sql_mode = @saved_sql_mode",
	)
}

// TestAlterEvent_NoContextIsBare confirms a body-only event change with NO applied context
// stays a bare ALTER EVENT (framing is opt-in via ApplySessionContext).
func TestAlterEvent_NoContextIsBare(t *testing.T) {
	n := NormalizerFor(MySQL80)
	from := loadSCCatalog(t, "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x = 1")
	to := loadSCCatalog(t, "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x = 2")
	// No ApplySessionContext.

	sql := planSQL(t, from, to, n)
	if !strings.Contains(sql, "ALTER EVENT `e`") {
		t.Fatalf("expected ALTER EVENT:\n%s", sql)
	}
	if strings.Contains(sql, "@saved_") || strings.Contains(sql, "sql_mode") {
		t.Fatalf("without applied context the ALTER must be bare (no framing):\n%s", sql)
	}
}

// TestAlterRoutine_NoFraming confirms a characteristic-only routine modify (ALTER FUNCTION/
// PROCEDURE) is emitted BARE even with context applied: ALTER FUNCTION does NOT re-stamp
// sql_mode (verified on live 8.0.32), so framing it is unnecessary and would be noise.
func TestAlterRoutine_NoFraming(t *testing.T) {
	n := NormalizerFor(MySQL80)
	from := loadSCCatalog(t, "CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC COMMENT 'old' RETURN a")
	to := loadSCCatalog(t, "CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC COMMENT 'new' RETURN a")

	from.ApplySessionContext(SessionContextMap{
		Functions: map[string]SessionContext{
			"f": {SQLMode: "PIPES_AS_CONCAT", CharacterSetClient: "utf8mb4", CollationConnection: "utf8mb4_0900_ai_ci"},
		},
	})

	d := DiffWithNormalizer(from, to, n)
	if len(d.Functions) != 1 || d.Functions[0].Action != DiffModify {
		t.Fatalf("want one function MODIFY, got %s", describeRoutineDiff(d))
	}
	if !routineAlterSuffices(d.Functions[0].From, d.Functions[0].To, n) {
		t.Fatalf("comment-only change should be ALTER-able")
	}
	sql := GenerateMigrationWithNormalizer(from, to, d, n).SQL()
	if !strings.HasPrefix(sql, "ALTER FUNCTION") {
		t.Fatalf("want a single ALTER FUNCTION, got:\n%s", sql)
	}
	if strings.Contains(sql, "@saved_") || strings.Contains(sql, "sql_mode") {
		t.Fatalf("characteristic-only ALTER FUNCTION must be bare (ALTER does not re-stamp sql_mode):\n%s", sql)
	}
}

// TestRecreate_TriggerCarriesSQLMode exercises the trigger recreate path (a body change is a
// DROP+CREATE — MySQL has no ALTER TRIGGER). Triggers already model charset/collation; this
// adds sql_mode and asserts the full framing.
func TestRecreate_TriggerCarriesSQLMode(t *testing.T) {
	n := NormalizerFor(MySQL80)
	from := loadSCCatalog(t, "CREATE TABLE t (id INT);\nCREATE TRIGGER trg BEFORE INSERT ON t FOR EACH ROW SET @x = 1")
	to := loadSCCatalog(t, "CREATE TABLE t (id INT);\nCREATE TRIGGER trg BEFORE INSERT ON t FOR EACH ROW SET @x = 2")

	from.ApplySessionContext(SessionContextMap{
		Triggers: map[string]SessionContext{
			"trg": {
				SQLMode:             "NO_ENGINE_SUBSTITUTION",
				CharacterSetClient:  "utf8mb4",
				CollationConnection: "utf8mb4_0900_ai_ci",
			},
		},
	})

	sql := planSQL(t, from, to, n)
	if !strings.Contains(sql, "DROP TRIGGER") {
		t.Fatalf("expected DROP TRIGGER in recreate plan:\n%s", sql)
	}
	mustContainInOrder(t, sql,
		"SET @saved_sql_mode = @@sql_mode",
		"SET sql_mode = 'NO_ENGINE_SUBSTITUTION'",
		"CREATE TRIGGER `d`.`trg`",
		"SET @x = 2",
		"SET sql_mode = @saved_sql_mode",
	)
	// A trigger has no time_zone axis.
	if strings.Contains(sql, "time_zone") {
		t.Fatalf("trigger framing must not touch time_zone:\n%s", sql)
	}
}

// TestRecreate_NoContextIsBareRecreate confirms that a MODIFY whose From carries NO applied
// session context stays a bare DROP+CREATE — the framing is opt-in via ApplySessionContext.
func TestRecreate_NoContextIsBareRecreate(t *testing.T) {
	n := NormalizerFor(MySQL80)
	from := loadSCCatalog(t, "CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC RETURN a + 1")
	to := loadSCCatalog(t, "CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC RETURN a + 2")
	// No ApplySessionContext.

	sql := planSQL(t, from, to, n)
	if !strings.Contains(sql, "DROP FUNCTION") || !strings.Contains(sql, "CREATE FUNCTION `d`.`f`") {
		t.Fatalf("expected bare DROP+CREATE:\n%s", sql)
	}
	if strings.Contains(sql, "@saved_") || strings.Contains(sql, "sql_mode") {
		t.Fatalf("without applied context the recreate must be bare (no framing):\n%s", sql)
	}
}

// mustContainInOrder asserts each substring appears in s, and in the given relative order.
func mustContainInOrder(t *testing.T, s string, subs ...string) {
	t.Helper()
	idx := 0
	for _, sub := range subs {
		pos := strings.Index(s[idx:], sub)
		if pos < 0 {
			t.Fatalf("missing or out-of-order substring %q\nfull SQL:\n%s", sub, s)
		}
		idx += pos + len(sub)
	}
}
