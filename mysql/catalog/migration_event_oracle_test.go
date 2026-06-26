package catalog

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// The oracle proof for the event differ + generator (correctness-protocol.md gates 1 & 2),
// against the LIVE MySQL engines (5.7 :13307 ssl-disabled, 8.0 :13306). Events are database-level
// objects, so this harness is event-shaped (SHOW CREATE EVENT, a per-probe apply database) rather
// than reusing the table harness. Two properties are proven on EVERY supported version:
//
//  1. IDEMPOTENCE (the spine). For each event form, the user-form CREATE EVENT applied to a real
//     database and read back via SHOW CREATE EVENT must diff EMPTY against the user form — proving
//     the canonicalization (STARTS-stripping, interval re-rendering, status / ON COMPLETION /
//     DEFINER normalization) matches what the engine stores. The stored form also self-diffs empty
//     and the no-op plan is "".
//
//  2. APPLY-CORRECTNESS. For each (from, to) event transition, the generated plan applied to a
//     real `from` database yields an event whose canonical form equals `to` (both directions).
//     Covers CREATE (add), DROP, and ALTER (every mutable aspect: schedule, status, ON COMPLETION,
//     comment, body). When nothing changed the plan is empty.
//
// The harness reuses connectOracle / serverCharsetFor / both / only / containsVersion /
// NormalizerFor from the existing oracle tests; it skips cleanly when the engines are unreachable,
// so the unit suite stays hermetic (go test -short skips it). Event DDL applies regardless of
// event_scheduler state, so these run with the scheduler at the box default.

// loadOneEvent loads a CREATE EVENT (wrapped in a database whose default charset matches the
// oracle box) and returns the named event from the catalog. Mirrors loadOneTable.
func loadOneEvent(t *testing.T, serverCharset, createSQL, event string) (*Catalog, *Event) {
	t.Helper()
	wrapped := fmt.Sprintf("CREATE DATABASE evdb DEFAULT CHARSET=%s;\nUSE evdb;\n%s", serverCharset, createSQL)
	cat, err := LoadSQL(wrapped)
	if err != nil {
		t.Fatalf("LoadSQL failed for %q: %v", createSQL, err)
	}
	for _, db := range cat.Databases() {
		if e := db.Events[toLower(event)]; e != nil {
			return cat, e
		}
	}
	t.Fatalf("event %q not found after load of %q", event, createSQL)
	return nil, nil
}

// eventReadback applies createSQL in a throwaway database (on the given pinned connection) and
// returns the SHOW CREATE EVENT readback (the engine's canonical stored form) as a reloadable
// CREATE EVENT statement. Returns ok=false (and logs) when the version rejects the input — the
// caller decides whether that is expected (e.g. 5.7 rejecting a far-future AT value).
func eventReadback(t *testing.T, conn *sql.Conn, name, sc, dbName, createSQL, event string) (string, bool) {
	t.Helper()
	ctx := context.Background()
	for _, s := range []string{
		"DROP DATABASE IF EXISTS " + dbName,
		fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s", dbName, sc),
		"USE " + dbName,
		createSQL,
	} {
		if _, err := conn.ExecContext(ctx, s); err != nil {
			t.Logf("[%s] event readback setup failed (may be expected): %q: %v", name, s, err)
			return "", false
		}
	}
	ddl, ok := scanShowCreateEvent(conn, dbName, event)
	if !ok {
		t.Logf("[%s] SHOW CREATE EVENT %s.%s failed", name, dbName, event)
	}
	return ddl, ok
}

// scanShowCreateEvent runs SHOW CREATE EVENT and returns the "Create Event" column. The result set
// has seven columns: Event, sql_mode, time_zone, Create Event, character_set_client,
// collation_connection, Database Collation.
func scanShowCreateEvent(conn *sql.Conn, dbName, event string) (string, bool) {
	var name, sqlMode, tz, ddl, csClient, collConn, dbColl string
	row := conn.QueryRowContext(context.Background(), fmt.Sprintf("SHOW CREATE EVENT %s.%s", dbName, event))
	if err := row.Scan(&name, &sqlMode, &tz, &ddl, &csClient, &collConn, &dbColl); err != nil {
		return "", false
	}
	return ddl, true
}

// eventIdempotenceProbes enumerates representative event FORMS that must round-trip empty: user
// DDL whose engine-stored form differs (auto-injected STARTS, re-rendered intervals, explicit
// default status/completion) but must canonicalize equal. Recurring + one-shot, varied
// schedule/status/options.
func eventIdempotenceProbes() []struct {
	id, event, create string
	versions          []Version
} {
	return []struct {
		id, event, create string
		versions          []Version
	}{
		{"every-default", "e", "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x = 1", both()},
		{"every-explicit-enable", "e", "CREATE EVENT e ON SCHEDULE EVERY 1 SECOND ENABLE DO SET @x = 1", both()},
		{"every-disable", "e", "CREATE EVENT e ON SCHEDULE EVERY 5 MINUTE DISABLE DO SET @x = 1", both()},
		{"every-disable-on-slave", "e", "CREATE EVENT e ON SCHEDULE EVERY 30 MINUTE ON COMPLETION NOT PRESERVE DISABLE ON SLAVE DO SET @x = 1", both()},
		{"every-preserve", "e", "CREATE EVENT e ON SCHEDULE EVERY 1 DAY ON COMPLETION PRESERVE DO SET @x = 1", both()},
		{"every-explicit-starts", "e", "CREATE EVENT e ON SCHEDULE EVERY 5 MINUTE STARTS '2025-03-01 12:00:00' DO SET @x = 1", both()},
		{"every-starts-ends", "e", "CREATE EVENT e ON SCHEDULE EVERY 1 DAY STARTS '2025-01-01 00:00:00' ENDS '2030-01-01 00:00:00' ON COMPLETION PRESERVE DISABLE COMMENT 'c' DO SET @q = 1", both()},
		{"every-only-ends", "e", "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR ENDS '2030-01-01 00:00:00' DO SET @x = 1", both()},
		{"every-interval-quoted-hm", "e", "CREATE EVENT e ON SCHEDULE EVERY '1:30' HOUR_MINUTE DO SET @x = 1", both()},
		{"every-interval-quoted-leadingzero", "e", "CREATE EVENT e ON SCHEDULE EVERY '02:03:04' HOUR_SECOND DO SET @x = 1", both()},
		{"every-interval-quoted-daysecond", "e", "CREATE EVENT e ON SCHEDULE EVERY '1 12:30:00' DAY_SECOND DO SET @x = 1", both()},
		{"every-interval-single-quoted", "e", "CREATE EVENT e ON SCHEDULE EVERY '05' MINUTE DO SET @x = 1", both()},
		{"every-comment", "e", "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR COMMENT 'it''s a test' DO SET @x = 1", both()},
		// One-shot AT: near-future literal works on both versions (5.7 rejects far-future AT values).
		{"at-literal", "e", "CREATE EVENT e ON SCHEDULE AT '2030-01-01 00:00:00' ON COMPLETION PRESERVE DO SET @y = 2", both()},
		{"at-preserve-disable", "e", "CREATE EVENT e ON SCHEDULE AT '2030-06-01 00:00:00' ON COMPLETION PRESERVE DISABLE DO SET @y = 2", both()},
	}
}

// TestOracle_EventIdempotence proves gate 1: each event form's user DDL vs its engine readback
// diffs EMPTY, the stored form self-diffs empty, and the no-op plan is "", on every supported
// version.
func TestOracle_EventIdempotence(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		sc := serverCharsetFor(version)
		n := NormalizerFor(version)
		ctx := context.Background()
		conn, err := o.db.Conn(ctx)
		if err != nil {
			t.Fatalf("[%s] grab conn: %v", o.name, err)
		}

		for _, p := range eventIdempotenceProbes() {
			if !containsVersion(p.versions, version) {
				continue
			}
			t.Run(fmt.Sprintf("%s/%s", o.name, p.id), func(t *testing.T) {
				slug := strings.ReplaceAll(p.id, "-", "_")
				rb, ok := eventReadback(t, conn, o.name, sc, "ev_idem_"+slug, p.create, p.event)
				if !ok {
					t.Skipf("[%s] could not obtain readback for %s", o.name, p.id)
				}

				userCat, _ := loadOneEvent(t, sc, p.create, p.event)
				storedCat, _ := loadOneEvent(t, sc, rb, p.event)

				if d := DiffWithNormalizer(storedCat, storedCat, n); !d.IsEmpty() {
					t.Errorf("[%s] IDEMPOTENCE: stored self-diff not empty for %s:\n  stored: %s",
						o.name, p.id, strings.TrimSpace(rb))
				}
				if d := DiffWithNormalizer(userCat, storedCat, n); !d.IsEmpty() {
					t.Errorf("[%s] CANONICALIZATION: user vs stored not empty for %s:\n  user:   %s\n  stored: %s\n  events: %s",
						o.name, p.id, strings.TrimSpace(p.create), strings.TrimSpace(rb), describeEventDiff(d))
				}
				if d := DiffWithNormalizer(storedCat, userCat, n); !d.IsEmpty() {
					t.Errorf("[%s] CANONICALIZATION (reverse): stored vs user not empty for %s:\n  events: %s",
						o.name, p.id, describeEventDiff(d))
				}
				plan := GenerateMigrationWithNormalizer(storedCat, storedCat, DiffWithNormalizer(storedCat, storedCat, n), n)
				if plan.SQL() != "" {
					t.Errorf("[%s] NON-EMPTY NO-OP PLAN for %s:\n%s", o.name, p.id, plan.SQL())
				}
			})
		}
		_ = conn.Close()
	}
}

// eventMigrationProbes enumerates the (from, to) event transitions the generator covers: CREATE
// (add), DROP, and ALTER for every mutable aspect (schedule, status, ON COMPLETION, comment, body)
// and EVERY↔AT.
func eventMigrationProbes() []struct {
	id, event, fromCreate, toCreate string
	versions                        []Version
} {
	const base = "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x = 1"
	return []struct {
		id, event, fromCreate, toCreate string
		versions                        []Version
	}{
		// CREATE (add).
		{"create-every", "e", "", base, both()},
		{"create-at", "e", "", "CREATE EVENT e ON SCHEDULE AT '2030-01-01 00:00:00' ON COMPLETION PRESERVE DO SET @y = 2", both()},
		{"create-full", "e", "", "CREATE EVENT e ON SCHEDULE EVERY 1 DAY STARTS '2025-01-01 00:00:00' ENDS '2030-01-01 00:00:00' ON COMPLETION PRESERVE DISABLE COMMENT 'c' DO SET @q = 1", both()},
		// DROP.
		{"drop", "e", base, "", both()},
		// ALTER — schedule interval.
		{"alter-schedule-interval", "e", base, "CREATE EVENT e ON SCHEDULE EVERY 2 HOUR DO SET @x = 1", both()},
		// ALTER — schedule unit.
		{"alter-schedule-unit", "e", base, "CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SET @x = 1", both()},
		// ALTER — EVERY -> AT (recurring to one-shot).
		{"alter-every-to-at", "e", base, "CREATE EVENT e ON SCHEDULE AT '2030-01-01 00:00:00' ON COMPLETION PRESERVE DO SET @x = 1", both()},
		// ALTER — status enable -> disable.
		{"alter-status-disable", "e", "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR ENABLE DO SET @x = 1", "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DISABLE DO SET @x = 1", both()},
		// ALTER — status disable -> disable on slave.
		{"alter-status-slave", "e", "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DISABLE DO SET @x = 1", "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DISABLE ON SLAVE DO SET @x = 1", both()},
		// ALTER — ON COMPLETION not preserve -> preserve.
		{"alter-completion", "e", "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR ON COMPLETION NOT PRESERVE DO SET @x = 1", "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR ON COMPLETION PRESERVE DO SET @x = 1", both()},
		// ALTER — comment.
		{"alter-comment", "e", "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR COMMENT 'a' DO SET @x = 1", "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR COMMENT 'b' DO SET @x = 1", both()},
		// ALTER — add a comment where none existed.
		{"alter-add-comment", "e", base, "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR COMMENT 'new' DO SET @x = 1", both()},
		// ALTER — body.
		{"alter-body", "e", base, "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x = 2", both()},
		// ALTER — add ENDS.
		{"alter-add-ends", "e", base, "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR ENDS '2030-01-01 00:00:00' DO SET @x = 1", both()},

		// ---- field-CLEARING transitions (an ALTER that omits a clause must RESET it, not keep the
		// old value). These are the class that hides the comment-removal bug. ----
		{"alter-remove-comment", "e", "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR COMMENT 'old' DO SET @x = 1", base, both()},
		{"alter-remove-ends", "e", "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR ENDS '2030-01-01 00:00:00' DO SET @x = 1", base, both()},
		{"alter-status-reenable", "e", "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DISABLE DO SET @x = 1", "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR ENABLE DO SET @x = 1", both()},
		{"alter-completion-reset", "e", "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR ON COMPLETION PRESERVE DO SET @x = 1", "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR ON COMPLETION NOT PRESERVE DO SET @x = 1", both()},
		// AT one-shot -> recurring EVERY (reverse of alter-every-to-at).
		{"alter-at-to-every", "e", "CREATE EVENT e ON SCHEDULE AT '2030-01-01 00:00:00' ON COMPLETION PRESERVE DO SET @x = 1", base, both()},
	}
}

// TestOracle_EventApplyCorrectness proves gate 2: the generated plan transforms a real `from`
// event-database into a `to`-equal one, for CREATE / DROP / ALTER across every form.
func TestOracle_EventApplyCorrectness(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		sc := serverCharsetFor(version)
		n := NormalizerFor(version)
		for _, p := range eventMigrationProbes() {
			if !containsVersion(p.versions, version) {
				continue
			}
			t.Run(fmt.Sprintf("%s/%s", o.name, p.id), func(t *testing.T) {
				assertEventApplyCorrect(t, o, n, sc, p.id, p.event, p.fromCreate, p.toCreate)
			})
		}
	}
}

// assertEventApplyCorrect loads from/to catalogs from the engine's own readbacks (authentic stored
// forms), generates the plan, applies it to a from-state database, reads the result back, and
// asserts the result canonicalizes equal to `to` (both directions). Each probe uses its own apply
// database on ONE pinned connection, so the test is hermetic.
func assertEventApplyCorrect(t *testing.T, o *oracleConn, n *Normalizer, sc, id, event, fromCreate, toCreate string) {
	t.Helper()
	slug := strings.ReplaceAll(id, "-", "_")
	applyDB := "eva_" + slug
	ctx := context.Background()
	conn, err := o.db.Conn(ctx)
	if err != nil {
		t.Fatalf("[%s] grab conn: %v", o.name, err)
	}
	defer func() { _ = conn.Close() }()

	eventInTo := strings.TrimSpace(toCreate) != ""
	eventInFrom := strings.TrimSpace(fromCreate) != ""

	var toCat, fromCat *Catalog
	if eventInTo {
		rb, ok := eventReadback(t, conn, o.name, sc, "evto_"+slug, toCreate, event)
		if !ok {
			t.Skipf("[%s] could not obtain `to` readback for %s", o.name, id)
		}
		toCat, _ = loadOneEvent(t, sc, rb, event)
	} else {
		toCat = New()
	}
	if eventInFrom {
		rb, ok := eventReadback(t, conn, o.name, sc, "evfrom_"+slug, fromCreate, event)
		if !ok {
			t.Skipf("[%s] could not obtain `from` readback for %s", o.name, id)
		}
		fromCat, _ = loadOneEvent(t, sc, rb, event)
	} else {
		fromCat = New()
	}

	diff := DiffWithNormalizer(fromCat, toCat, n)
	plan := GenerateMigrationWithNormalizer(fromCat, toCat, diff, n)

	setup := []string{
		"DROP DATABASE IF EXISTS " + applyDB,
		fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s", applyDB, sc),
		"USE " + applyDB,
	}
	if eventInFrom {
		setup = append(setup, fromCreate)
	}
	for _, s := range setup {
		if _, err := conn.ExecContext(ctx, s); err != nil {
			t.Skipf("[%s] could not set up `from` state for %s: %q: %v", o.name, id, s, err)
		}
	}
	for _, op := range plan.Ops {
		if _, err := conn.ExecContext(ctx, op.SQL); err != nil {
			t.Fatalf("[%s] APPLY FAILED for %s:\n  stmt: %s\n  err: %v\n  full plan:\n%s",
				o.name, id, op.SQL, err, plan.SQL())
		}
	}

	if !eventInTo {
		if _, ok := scanShowCreateEvent(conn, applyDB, event); ok {
			t.Errorf("[%s] %s: event %s still exists after DROP plan:\n%s", o.name, id, event, plan.SQL())
		}
		return
	}

	resultRB, ok := scanShowCreateEvent(conn, applyDB, event)
	if !ok {
		t.Fatalf("[%s] %s: result event %s missing after apply:\n%s", o.name, id, event, plan.SQL())
	}
	resultCat, _ := loadOneEvent(t, sc, resultRB, event)

	if d := DiffWithNormalizer(resultCat, toCat, n); !d.IsEmpty() {
		t.Errorf("[%s] APPLY-CORRECTNESS FAILED for %s: result != to\n  plan:\n%s\n  result: %s\n  events: %s",
			o.name, id, plan.SQL(), strings.TrimSpace(resultRB), describeEventDiff(d))
	}
	if d := DiffWithNormalizer(toCat, resultCat, n); !d.IsEmpty() {
		t.Errorf("[%s] APPLY-CORRECTNESS (reverse) FAILED for %s: %s", o.name, id, describeEventDiff(d))
	}
}

// describeEventDiff renders a compact human description of a SchemaDiff's event changes.
func describeEventDiff(d *SchemaDiff) string {
	var b strings.Builder
	for _, e := range d.Events {
		fmt.Fprintf(&b, "[event %s %s]", e.Name, e.Action)
	}
	if b.Len() == 0 {
		return "(no event changes)"
	}
	return b.String()
}
