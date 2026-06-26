package catalog

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// The oracle proof for the trigger differ + generator (correctness-protocol.md gates 1 & 2),
// against the LIVE MySQL engines (5.7 :13307 ssl-disabled, 8.0 :13306). Triggers are database-
// level objects whose stored form is read from information_schema.triggers (the exact path the
// bytebase MySQL driver syncs them with), so the table-centric harness in migration_oracle_test.go
// does not fit — this file carries a trigger-shaped harness that:
//
//   - applies a schema (a base table plus its triggers) to a throwaway database;
//   - reads the triggers back from information_schema.triggers (TRIGGER_NAME / ACTION_TIMING /
//     EVENT_MANIPULATION / EVENT_OBJECT_TABLE / ACTION_STATEMENT) and reconstructs CREATE TRIGGER
//     statements — the engine's authentic stored form — then loads them via LoadSDL;
//   - proves both gates on that stored-form catalog.
//
// GATE 1 — IDEMPOTENCE (the spine): for a schema in its real stored form, Diff(c,c) is empty and
// the no-op plan SQL() == "". Additionally the user-written form diffed against its engine readback
// is empty (the canonicalization property: DEFINER/whitespace differences must not phantom-diff).
//
// GATE 2 — APPLY-CORRECTNESS: for representative (from, to) pairs, GenerateMigration(from,to).SQL()
// applied to a real `from` database yields a trigger set whose stored form equals `to`'s.
//
// The harness skips cleanly when the engines are unreachable (go test -short skips it).

// triggerOracleDB is the database the trigger harness applies generated plans into (mirrors the
// table harness's diffdb convention but kept distinct so parallel oracle tests don't collide).
const triggerOracleDB = "trigdb"

// triggerSchemaSDL wraps a schema DDL list (base tables + CREATE TRIGGER statements) in a
// CREATE DATABASE + USE under dbName, returning a loadable SDL snippet. Used to build the
// user-written form for the canonicalization check, under the same database name AND server
// charset as the engine's reconstructed stored form (serverCharsetFor(version) — latin1 on 5.7,
// utf8mb4 on 8.0) so only trigger differences can surface, not a table-charset mismatch (the same
// wrapping convention loadOneTable uses for the table harness).
func triggerSchemaSDL(dbName string, version Version, schema []string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "CREATE DATABASE %s DEFAULT CHARSET=%s;\nUSE %s;\n", dbName, serverCharsetFor(version), dbName)
	for _, s := range schema {
		b.WriteString(s)
		b.WriteString(";\n")
	}
	return b.String()
}

// reconstructTriggersSQL applies the given schema DDL to a throwaway database (probeDB), reads the
// triggers back from information_schema, and returns an SDL snippet (CREATE DATABASE + USE + every
// base table's SHOW CREATE + every reconstructed CREATE TRIGGER) representing the engine's stored
// form — but with the database name in the returned SDL rewritten to emitDB. Separating probeDB
// from emitDB lets the apply-correctness harness load BOTH the `from` and `to` catalogs under the
// SAME database name (the apply target), so the table differ doesn't see a spurious DROP+CREATE of
// the carried base table just because the two probes ran in differently-named scratch databases —
// exactly the loadOneTable-wraps-everything-in-diffdb convention the table harness uses.
// tables lists the base tables to carry over (their SHOW CREATE TABLE), so the reloaded catalog has
// the trigger's target table present.
func (o *oracleConn) reconstructTriggersSQL(t *testing.T, probeDB, emitDB string, schemaDDL []string, tables []string) (string, bool) {
	t.Helper()
	ctx := context.Background()
	stmts := append([]string{
		"DROP DATABASE IF EXISTS " + probeDB,
		fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s", probeDB, serverCharsetFor(o.version)),
		"USE " + probeDB,
	}, schemaDDL...)
	for _, s := range stmts {
		if _, err := o.db.ExecContext(ctx, s); err != nil {
			t.Logf("[%s] trigger setup failed (may be expected): %q: %v", o.name, s, err)
			return "", false
		}
	}

	var b strings.Builder
	fmt.Fprintf(&b, "CREATE DATABASE %s DEFAULT CHARSET=%s;\nUSE %s;\n", emitDB, serverCharsetFor(o.version), emitDB)
	// Carry the base tables (so the trigger's ON table exists when the SDL reloads).
	for _, tbl := range tables {
		var name, ddl string
		row := o.db.QueryRowContext(ctx, fmt.Sprintf("SHOW CREATE TABLE %s.%s", probeDB, tbl))
		if err := row.Scan(&name, &ddl); err != nil {
			t.Logf("[%s] SHOW CREATE TABLE %s failed: %v", o.name, tbl, err)
			return "", false
		}
		b.WriteString(ddl)
		b.WriteString(";\n")
	}
	// Reconstruct every trigger from information_schema, ordered by ACTION_ORDER for determinism.
	rows, err := o.db.QueryContext(ctx,
		"SELECT TRIGGER_NAME, ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_STATEMENT "+
			"FROM information_schema.triggers WHERE TRIGGER_SCHEMA = ? ORDER BY EVENT_OBJECT_TABLE, ACTION_TIMING, EVENT_MANIPULATION, ACTION_ORDER", probeDB)
	if err != nil {
		t.Logf("[%s] query information_schema.triggers failed: %v", o.name, err)
		return "", false
	}
	defer rows.Close()
	for rows.Next() {
		var name, timing, event, table, stmt string
		if err := rows.Scan(&name, &timing, &event, &table, &stmt); err != nil {
			t.Logf("[%s] scan trigger row failed: %v", o.name, err)
			return "", false
		}
		fmt.Fprintf(&b, "CREATE TRIGGER `%s` %s %s ON `%s` FOR EACH ROW %s;\n", name, timing, event, table, stmt)
	}
	if err := rows.Err(); err != nil {
		t.Logf("[%s] iterate trigger rows failed: %v", o.name, err)
		return "", false
	}
	return b.String(), true
}

// triggerStoredCatalog reconstructs and loads the engine's stored form of a schema into a catalog.
// The schema is applied in probeDB but the loaded catalog's database is named emitDB.
func (o *oracleConn) triggerStoredCatalog(t *testing.T, probeDB, emitDB string, schemaDDL []string, tables []string) (*Catalog, bool) {
	t.Helper()
	sdl, ok := o.reconstructTriggersSQL(t, probeDB, emitDB, schemaDDL, tables)
	if !ok {
		return nil, false
	}
	cat, err := LoadSDLWithVersion(sdl, o.version)
	if err != nil {
		t.Fatalf("[%s] reload of reconstructed trigger SDL failed: %v\n%s", o.name, err, sdl)
	}
	return cat, true
}

// triggerIdemProbe is one idempotence case: a base-table set + the trigger DDL applied on top.
type triggerIdemProbe struct {
	id       string
	tables   []string // base tables (SHOW CREATE carried into the reload)
	schema   []string // full schema DDL (tables + CREATE TRIGGER ...) applied to the engine
	versions []Version
}

// triggerIdempotenceProbes enumerates the trigger FORMS that must round-trip empty: every
// timing×event, multi-trigger-per-table, BEGIN…END bodies, OLD/NEW references, and a trigger with
// FOLLOWS (whose readback drops the FOLLOWS — the order-not-modelled flag).
func triggerIdempotenceProbes() []triggerIdemProbe {
	tbl := "CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)"
	mk := func(id, trig string, vs []Version) triggerIdemProbe {
		return triggerIdemProbe{
			id:       id,
			tables:   []string{"t"},
			schema:   []string{tbl, trig},
			versions: vs,
		}
	}
	return []triggerIdemProbe{
		// ---- every timing × event ----
		mk("before-insert", "CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1", both()),
		mk("after-insert", "CREATE TRIGGER tr AFTER INSERT ON t FOR EACH ROW SET @x = NEW.a", both()),
		mk("before-update", "CREATE TRIGGER tr BEFORE UPDATE ON t FOR EACH ROW SET NEW.a = OLD.a + 1", both()),
		mk("after-update", "CREATE TRIGGER tr AFTER UPDATE ON t FOR EACH ROW SET @x = NEW.b", both()),
		mk("before-delete", "CREATE TRIGGER tr BEFORE DELETE ON t FOR EACH ROW SET @x = OLD.a", both()),
		mk("after-delete", "CREATE TRIGGER tr AFTER DELETE ON t FOR EACH ROW SET @x = OLD.b", both()),

		// ---- body shapes ----
		// BEGIN…END multi-statement body (verbatim storage incl. newlines/indentation).
		mk("begin-end", "CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW BEGIN\n  SET NEW.a = 1;\n  IF NEW.b IS NULL THEN\n    SET NEW.b = 0;\n  END IF;\nEND", both()),
		// OLD/NEW + a string literal in the body (quoting must survive).
		mk("string-literal", "CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.c = (SELECT 'literal value')", both()),
		// Lower-case keywords in the user form (MySQL preserves body case verbatim).
		mk("lowercase-body", "CREATE TRIGGER tr before insert on t for each row set new.a = new.b + 1", both()),

		// ---- multiple triggers, distinct events on one table ----
		// (handled by the multi-trigger probe below, which needs >1 CREATE)
	}
}

// TestOracle_TriggerIdempotence proves gate 1: a schema in its real stored form self-diffs empty
// and the no-op plan is empty, and the user form vs its engine readback is empty (canonicalization
// — DEFINER, whitespace and FOLLOWS must not phantom-diff). Probed on every supported version.
func TestOracle_TriggerIdempotence(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		n := NormalizerFor(version)
		for _, probe := range triggerIdempotenceProbes() {
			if !containsVersion(probe.versions, version) {
				continue
			}
			t.Run(fmt.Sprintf("%s/%s", o.name, probe.id), func(t *testing.T) {
				dbName := "trig_idem_" + strings.ReplaceAll(probe.id, "-", "_")
				storedCat, ok := o.triggerStoredCatalog(t, dbName, dbName, probe.schema, probe.tables)
				if !ok {
					t.Skipf("[%s] could not obtain stored form for %s", o.name, probe.id)
				}

				// Self-diff + self-plan must be empty (idempotence spine).
				if d := DiffWithNormalizer(storedCat, storedCat, n); !d.IsEmpty() {
					t.Errorf("[%s] IDEMPOTENCE: self-diff not empty for %s: %s", o.name, probe.id, describeTriggerDiff(d))
				}
				plan := GenerateMigrationWithNormalizer(storedCat, storedCat, DiffWithNormalizer(storedCat, storedCat, n), n)
				if plan.SQL() != "" {
					t.Errorf("[%s] NON-EMPTY NO-OP PLAN for %s:\n%s", o.name, probe.id, plan.SQL())
				}

				// User form vs engine stored form must collapse to empty (the canonicalization
				// property: the user wrote no DEFINER / maybe a FOLLOWS / odd spacing; the readback
				// carries the engine's form — they must not phantom-diff). The user form is loaded
				// under the SAME database name as the stored form so only trigger differences (not a
				// database-name mismatch) can surface.
				userSDL := triggerSchemaSDL(dbName, version, probe.schema)
				userCat, err := LoadSDLWithVersion(userSDL, version)
				if err != nil {
					t.Fatalf("[%s] user-form load failed for %s: %v", o.name, probe.id, err)
				}
				if d := DiffWithNormalizer(userCat, storedCat, n); !d.IsEmpty() {
					t.Errorf("[%s] CANONICALIZATION: user vs stored not empty for %s: %s", o.name, probe.id, describeTriggerDiff(d))
				}
				if d := DiffWithNormalizer(storedCat, userCat, n); !d.IsEmpty() {
					t.Errorf("[%s] CANONICALIZATION (reverse): stored vs user not empty for %s: %s", o.name, probe.id, describeTriggerDiff(d))
				}
			})
		}
	}
}

// TestOracle_TriggerMultiPerTableIdempotence proves a table carrying MULTIPLE triggers (across
// several timing/event slots, including two on the same timing+event linked by FOLLOWS) round-trips
// empty. The FOLLOWS is dropped by the readback (order-not-modelled flag): the stored form and the
// user form must still collapse to empty.
func TestOracle_TriggerMultiPerTableIdempotence(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		n := NormalizerFor(version)
		t.Run(o.name+"/multi", func(t *testing.T) {
			tbl := "CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)"
			schema := []string{
				tbl,
				"CREATE TRIGGER bi1 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1",
				"CREATE TRIGGER bi2 BEFORE INSERT ON t FOR EACH ROW FOLLOWS bi1 SET NEW.b = 2",
				"CREATE TRIGGER au AFTER UPDATE ON t FOR EACH ROW SET @x = NEW.c",
				"CREATE TRIGGER bd BEFORE DELETE ON t FOR EACH ROW SET @y = OLD.a",
			}
			storedCat, ok := o.triggerStoredCatalog(t, "trig_multi", "trig_multi", schema, []string{"t"})
			if !ok {
				t.Skipf("[%s] could not obtain stored form", o.name)
			}
			if d := DiffWithNormalizer(storedCat, storedCat, n); !d.IsEmpty() {
				t.Errorf("[%s] multi-trigger self-diff not empty: %s", o.name, describeTriggerDiff(d))
			}
			if plan := GenerateMigrationWithNormalizer(storedCat, storedCat, DiffWithNormalizer(storedCat, storedCat, n), n); plan.SQL() != "" {
				t.Errorf("[%s] multi-trigger non-empty no-op plan:\n%s", o.name, plan.SQL())
			}

			// User form (with the FOLLOWS) vs stored form (FOLLOWS dropped) must be empty.
			userCat, err := LoadSDLWithVersion(triggerSchemaSDL("trig_multi", version, schema), version)
			if err != nil {
				t.Fatalf("[%s] user-form load failed: %v", o.name, err)
			}
			if d := DiffWithNormalizer(userCat, storedCat, n); !d.IsEmpty() {
				t.Errorf("[%s] multi-trigger user vs stored not empty (FOLLOWS must not diff): %s", o.name, describeTriggerDiff(d))
			}
		})
	}
}

// triggerMigProbe is one apply-correctness case: transform a from-schema into a to-schema. Both
// are full schema DDL lists (base table + triggers); an empty trigger set means "no trigger".
//
// tables is the base-table set carried into the from/to/result readbacks (so the trigger's ON table
// is present when the reconstructed SDL reloads). When the from and to schemas carry DIFFERENT base
// tables (a trigger relocated to a new table while the old table is dropped), fromTables and toTables
// override tables for the from and to+result readbacks respectively.
type triggerMigProbe struct {
	id         string
	tables     []string
	fromTables []string // overrides tables for the `from` readback when set
	toTables   []string // overrides tables for the `to` (and post-apply result) readback when set
	fromDDL    []string
	toDDL      []string
	versions   []Version
}

// fromTableSet / toTableSet resolve the effective base-table list for each side.
func (p triggerMigProbe) fromTableSet() []string {
	if p.fromTables != nil {
		return p.fromTables
	}
	return p.tables
}

func (p triggerMigProbe) toTableSet() []string {
	if p.toTables != nil {
		return p.toTables
	}
	return p.tables
}

// triggerMigrationProbes enumerates the CREATE / DROP / CHANGE(=DROP+CREATE) forms the generator
// covers, each proven by applying the generated plan to a real `from` database and comparing the
// trigger readback to `to`.
func triggerMigrationProbes() []triggerMigProbe {
	tbl := "CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)"
	return []triggerMigProbe{
		// ---- CREATE (add a trigger to an existing table) ----
		{"create-before-insert", []string{"t"}, nil, nil, []string{tbl},
			[]string{tbl, "CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1"}, both()},
		{"create-after-update", []string{"t"}, nil, nil, []string{tbl},
			[]string{tbl, "CREATE TRIGGER tr AFTER UPDATE ON t FOR EACH ROW SET @x = NEW.b"}, both()},
		{"create-before-delete", []string{"t"}, nil, nil, []string{tbl},
			[]string{tbl, "CREATE TRIGGER tr BEFORE DELETE ON t FOR EACH ROW SET @x = OLD.a"}, both()},
		{"create-begin-end", []string{"t"}, nil, nil, []string{tbl},
			[]string{tbl, "CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW BEGIN\n  SET NEW.a = 1;\n  SET NEW.b = 2;\nEND"}, both()},

		// ---- DROP (remove a trigger, table stays) ----
		{"drop-trigger", []string{"t"}, nil, nil,
			[]string{tbl, "CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1"},
			[]string{tbl}, both()},

		// ---- CHANGE (= DROP + CREATE) ----
		{"change-body", []string{"t"}, nil, nil,
			[]string{tbl, "CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1"},
			[]string{tbl, "CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 2"}, both()},
		{"change-timing", []string{"t"}, nil, nil,
			[]string{tbl, "CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET @x = NEW.a"},
			[]string{tbl, "CREATE TRIGGER tr AFTER INSERT ON t FOR EACH ROW SET @x = NEW.a"}, both()},
		{"change-event", []string{"t"}, nil, nil,
			[]string{tbl, "CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1"},
			[]string{tbl, "CREATE TRIGGER tr BEFORE UPDATE ON t FOR EACH ROW SET NEW.a = 1"}, both()},

		// ---- multi-trigger plan: add one, drop one, keep one ----
		{"multi-add-drop", []string{"t"}, nil, nil,
			[]string{tbl,
				"CREATE TRIGGER keep AFTER UPDATE ON t FOR EACH ROW SET @x = NEW.a",
				"CREATE TRIGGER goes BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1"},
			[]string{tbl,
				"CREATE TRIGGER keep AFTER UPDATE ON t FOR EACH ROW SET @x = NEW.a",
				"CREATE TRIGGER comes BEFORE DELETE ON t FOR EACH ROW SET @y = OLD.a"}, both()},

		// ---- regression: trigger RELOCATED to a new table while the OLD table is dropped ----
		// Identity is (database, name), so moving `tr` from old_t to new_t is a MODIFY (DROP+CREATE).
		// The DROP half targets old_t, which is dropped in the same plan — MySQL cascades the trigger
		// when DROP TABLE old_t runs, so an explicit DROP TRIGGER would fail (errno 1360). The
		// suppression in generateTriggerDDL's DiffModify branch must skip that DROP; only the CREATE
		// on the surviving new_t is emitted. (Without the fix this probe fails at apply.)
		{"modify-table-move-drop-old", nil,
			[]string{"old_t", "new_t"}, []string{"new_t"},
			[]string{
				"CREATE TABLE old_t (id INT PRIMARY KEY, a INT)",
				"CREATE TABLE new_t (id INT PRIMARY KEY, a INT)",
				"CREATE TRIGGER tr BEFORE INSERT ON old_t FOR EACH ROW SET NEW.a = 1"},
			[]string{
				"CREATE TABLE new_t (id INT PRIMARY KEY, a INT)",
				"CREATE TRIGGER tr BEFORE INSERT ON new_t FOR EACH ROW SET NEW.a = 1"}, both()},
	}
}

// TestOracle_TriggerMigrationApplyCorrectness proves gate 2: the generated DDL transforms a real
// `from` database's trigger set into a `to`-equal one (compared via information_schema readback,
// reloaded and diffed).
func TestOracle_TriggerMigrationApplyCorrectness(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		n := NormalizerFor(version)
		for _, probe := range triggerMigrationProbes() {
			if !containsVersion(probe.versions, version) {
				continue
			}
			t.Run(fmt.Sprintf("%s/%s", o.name, probe.id), func(t *testing.T) {
				assertTriggerApplyCorrect(t, o, n, probe)
			})
		}
	}
}

// assertTriggerApplyCorrect builds from/to catalogs from the engine's own readbacks, generates the
// plan, applies it to a real from-state database, reads the result's triggers back, and asserts the
// result canonicalizes equal to `to`.
func assertTriggerApplyCorrect(t *testing.T, o *oracleConn, n *Normalizer, p triggerMigProbe) {
	t.Helper()
	slug := strings.ReplaceAll(p.id, "-", "_")

	// Both probes reconstruct under DIFFERENT scratch databases but EMIT under triggerOracleDB —
	// the apply target — so the table differ compares the carried base table against itself (no
	// spurious DROP+CREATE table) and only the trigger set differs.
	fromCat, ok := o.triggerStoredCatalog(t, "trig_from_"+slug, triggerOracleDB, p.fromDDL, p.fromTableSet())
	if !ok {
		t.Skipf("[%s] could not obtain `from` stored form for %s", o.name, p.id)
	}
	toCat, ok := o.triggerStoredCatalog(t, "trig_to_"+slug, triggerOracleDB, p.toDDL, p.toTableSet())
	if !ok {
		t.Skipf("[%s] could not obtain `to` stored form for %s", o.name, p.id)
	}

	diff := DiffWithNormalizer(fromCat, toCat, n)
	plan := GenerateMigrationWithNormalizer(fromCat, toCat, diff, n)

	// Build a real database in state `from` on a dedicated connection (so a stray USE stays on the
	// same pooled connection as the following statements), then apply the plan.
	ctx := context.Background()
	conn, err := o.db.Conn(ctx)
	if err != nil {
		t.Fatalf("[%s] grab conn: %v", o.name, err)
	}
	defer func() { _ = conn.Close() }()

	setup := append([]string{
		"DROP DATABASE IF EXISTS " + triggerOracleDB,
		fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s", triggerOracleDB, serverCharsetFor(o.version)),
		"USE " + triggerOracleDB,
	}, p.fromDDL...)
	for _, s := range setup {
		if _, err := conn.ExecContext(ctx, s); err != nil {
			t.Skipf("[%s] could not set up `from` state for %s: %q: %v", o.name, p.id, s, err)
		}
	}

	for _, op := range plan.Ops {
		if _, err := conn.ExecContext(ctx, op.SQL); err != nil {
			t.Fatalf("[%s] APPLY FAILED for %s:\n  stmt: %s\n  err: %v\n  full plan:\n%s", o.name, p.id, op.SQL, err, plan.SQL())
		}
	}

	// Read back the resulting triggers and reload as a catalog (in the applied database, on the
	// SAME connection the plan was applied on). After apply the schema equals `to`, so the result
	// readback carries the `to` base-table set.
	resultSDL, ok := o.reconstructTriggersFromConn(t, conn, triggerOracleDB, p.toTableSet())
	if !ok {
		t.Fatalf("[%s] could not read back result for %s after apply:\n%s", o.name, p.id, plan.SQL())
	}
	resultCat, err := LoadSDLWithVersion(resultSDL, o.version)
	if err != nil {
		t.Fatalf("[%s] reload of result SDL failed for %s: %v\n%s", o.name, p.id, err, resultSDL)
	}

	if d := DiffWithNormalizer(resultCat, toCat, n); !d.IsEmpty() {
		t.Errorf("[%s] APPLY-CORRECTNESS FAILED for %s: result != to\n  plan:\n%s\n  diff: %s", o.name, p.id, plan.SQL(), describeTriggerDiff(d))
	}
	if d := DiffWithNormalizer(toCat, resultCat, n); !d.IsEmpty() {
		t.Errorf("[%s] APPLY-CORRECTNESS (reverse) FAILED for %s: %s", o.name, p.id, describeTriggerDiff(d))
	}
}

// reconstructTriggersFromConn reads the triggers (and carries the base tables' SHOW CREATE) of an
// already-applied database on a specific *sql.Conn, returning the stored-form SDL. Unlike
// reconstructTriggersSQL it does NOT re-create the database — it inspects the just-applied state on
// the same connection the plan ran on.
func (o *oracleConn) reconstructTriggersFromConn(t *testing.T, conn *sql.Conn, dbName string, tables []string) (string, bool) {
	t.Helper()
	ctx := context.Background()
	var b strings.Builder
	fmt.Fprintf(&b, "CREATE DATABASE %s DEFAULT CHARSET=%s;\nUSE %s;\n", dbName, serverCharsetFor(o.version), dbName)
	for _, tbl := range tables {
		var name, ddl string
		row := conn.QueryRowContext(ctx, fmt.Sprintf("SHOW CREATE TABLE %s.%s", dbName, tbl))
		if err := row.Scan(&name, &ddl); err != nil {
			t.Logf("[%s] SHOW CREATE TABLE %s failed: %v", o.name, tbl, err)
			return "", false
		}
		b.WriteString(ddl)
		b.WriteString(";\n")
	}
	rows, err := conn.QueryContext(ctx,
		"SELECT TRIGGER_NAME, ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_STATEMENT "+
			"FROM information_schema.triggers WHERE TRIGGER_SCHEMA = ? ORDER BY EVENT_OBJECT_TABLE, ACTION_TIMING, EVENT_MANIPULATION, ACTION_ORDER", dbName)
	if err != nil {
		t.Logf("[%s] query information_schema.triggers failed: %v", o.name, err)
		return "", false
	}
	defer rows.Close()
	for rows.Next() {
		var name, timing, event, table, stmt string
		if err := rows.Scan(&name, &timing, &event, &table, &stmt); err != nil {
			t.Logf("[%s] scan trigger row failed: %v", o.name, err)
			return "", false
		}
		fmt.Fprintf(&b, "CREATE TRIGGER `%s` %s %s ON `%s` FOR EACH ROW %s;\n", name, timing, event, table, stmt)
	}
	if err := rows.Err(); err != nil {
		t.Logf("[%s] iterate trigger rows failed: %v", o.name, err)
		return "", false
	}
	return b.String(), true
}
