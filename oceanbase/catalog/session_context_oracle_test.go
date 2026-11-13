package catalog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// End-to-end apply-correctness proof for BYT-9832 against OceanBase: a
// routine / event whose BODY changes must be re-emitted under its ORIGINAL session context.
// Each case creates the object under a distinctive sql_mode (and, for events, time_zone),
// syncs that context into the source catalog via ApplySessionContext exactly as bytebase
// will, generates the plan, APPLIES it to the real server, and reads back
// information_schema to prove the original context survived the recreate/ALTER. Without the
// framing these assertions fail (the object comes back under the deploy session's mode) —
// which is the whole point of the change.

const scOracleDB = "byt9832_sc_apply"

// sessionExecer is the part of oracleSession the setup helpers use.
type sessionExecer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

// execUnderSession runs the setup DDL on a dedicated connection whose session is first set
// to the given sql_mode / time_zone, so the created object captures that context — mirroring
// how a real routine/event acquires its stored sql_mode.
//
// Every engine this runs on creates the seed, so a failure, a dropped connection included,
// fails the test: skipping would pass the proof without running it.
func execUnderSession(t testing.TB, conn sessionExecer, ctx context.Context, sqlMode, timeZone string, stmts ...string) {
	t.Helper()
	pre := []string{"SET SESSION sql_mode = " + quoteStringLiteral(sqlMode)}
	if timeZone != "" {
		pre = append(pre, "SET SESSION time_zone = "+quoteStringLiteral(timeZone))
	}
	for _, s := range append(pre, stmts...) {
		if _, err := conn.ExecContext(ctx, s); err != nil {
			t.Fatalf("setup under sql_mode %q: %q: %v", sqlMode, s, err)
		}
	}
}

// TestOracle_RoutineRecreatePreservesSQLMode proves the routine DROP+CREATE recreate framing
// preserves sql_mode on the live engine.
func TestOracle_RoutineRecreatePreservesSQLMode(t *testing.T) {
	const origMode = "PIPES_AS_CONCAT,NO_BACKSLASH_ESCAPES"
	for _, version := range oracleVersions() {
		o := oceanBaseOracle(t)
		n := NormalizerFor(version)
		t.Run(o.name, gapped(func(t testing.TB) {
			ctx := context.Background()
			conn, err := o.session(ctx)
			if err != nil {
				t.Fatalf("grab conn: %v", err)
			}
			defer func() { _ = conn.Close() }()

			db := scOracleDB + "_rtn"
			// Deploy session runs under a DIFFERENT (default-ish) mode so a bare recreate would
			// lose PIPES_AS_CONCAT.
			deployMode := "STRICT_TRANS_TABLES,NO_ENGINE_SUBSTITUTION"

			setup := []string{
				"DROP DATABASE IF EXISTS " + db,
				fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s", db, serverCharsetFor(version)),
				"USE " + db,
			}
			for _, s := range setup {
				if _, err := conn.ExecContext(ctx, s); err != nil {
					t.Fatalf("[%s] setup: %q: %v", o.name, s, err)
				}
			}
			// Create the function under the ORIGINAL mode.
			execUnderSession(t, conn, ctx, origMode, "",
				"CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC RETURN a + 1")

			// Build source (from) catalog = current state, and stamp the synced context onto it.
			from := mustLoadUnderDB(t, db, "CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC RETURN a + 1", version)
			from.ApplySessionContext(SessionContextMap{Functions: map[string]SessionContext{
				"f": {SQLMode: origMode, CharacterSetClient: serverCharsetFor(version), CollationConnection: defaultCollationForServerCharset(version)},
			}})
			// Target (to) catalog = desired body change.
			to := mustLoadUnderDB(t, db, "CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC RETURN a + 2", version)

			plan := GenerateMigrationWithNormalizer(from, to, DiffWithNormalizer(from, to, n), n)
			if plan == nil || len(plan.Ops) == 0 {
				t.Fatalf("[%s] expected a non-empty recreate plan", o.name)
			}

			// Apply under the DEPLOY session mode; the framing must override it per-statement.
			if _, err := conn.ExecContext(ctx, "SET SESSION sql_mode = "+quoteStringLiteral(deployMode)); err != nil {
				t.Fatalf("[%s] set deploy mode: %v", o.name, err)
			}
			if !o.applyPlan(t, conn, plan) {
				return
			}

			got := routineSQLMode(t, conn, ctx, db, "f")
			if got != origMode {
				t.Fatalf("[%s] routine sql_mode not preserved across recreate: got %q want %q\nplan:\n%s", o.name, got, origMode, plan.SQL())
			}
		}))
	}
}

// TestOracle_EventAlterPreservesSQLMode proves the ALTER EVENT framing preserves sql_mode on
// the live engine — the case that a bare ALTER would silently re-stamp.
func TestOracle_EventAlterPreservesSQLMode(t *testing.T) {
	const origMode = "PIPES_AS_CONCAT"
	const origTZ = "+08:00"
	for _, version := range oracleVersions() {
		o := oceanBaseOracle(t)
		n := NormalizerFor(version)
		t.Run(o.name, gapped(func(t testing.TB) {
			ctx := context.Background()
			conn, err := o.session(ctx)
			if err != nil {
				t.Fatalf("grab conn: %v", err)
			}
			defer func() { _ = conn.Close() }()

			db := scOracleDB + "_evt"
			deployMode := "STRICT_TRANS_TABLES,NO_ENGINE_SUBSTITUTION"

			setup := []string{
				"DROP DATABASE IF EXISTS " + db,
				fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s", db, serverCharsetFor(version)),
				"USE " + db,
				"SET GLOBAL event_scheduler = OFF",
			}
			for _, s := range setup {
				if _, err := conn.ExecContext(ctx, s); err != nil {
					// A user without SUPER cannot toggle event_scheduler, and the event is
					// created DISABLE anyway, so the server's refusal of that one is fine.
					// Anything else, and any answer that is not the server's, fails.
					var myErr *gomysql.MySQLError
					if !strings.Contains(strings.ToLower(s), "event_scheduler") || !errors.As(err, &myErr) {
						t.Fatalf("[%s] setup: %q: %v", o.name, s, err)
					}
				}
			}
			execUnderSession(t, conn, ctx, origMode, origTZ,
				"CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DISABLE DO SET @x = 1")

			from := mustLoadUnderDB(t, db, "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DISABLE DO SET @x = 1", version)
			from.ApplySessionContext(SessionContextMap{Events: map[string]SessionContext{
				"e": {SQLMode: origMode, CharacterSetClient: serverCharsetFor(version), CollationConnection: defaultCollationForServerCharset(version), TimeZone: origTZ},
			}})
			to := mustLoadUnderDB(t, db, "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DISABLE DO SET @x = 2", version)

			d := DiffWithNormalizer(from, to, n)
			if len(d.Events) != 1 || d.Events[0].Action != DiffModify {
				t.Fatalf("[%s] want one event MODIFY, got %+v", o.name, d.Events)
			}
			// This is the ALTER path (body-only change), the case a bare ALTER would re-stamp.
			if eventModifyNeedsRecreate(d.Events[0].From, d.Events[0].To) {
				t.Fatalf("[%s] expected the ALTER path for a body-only change", o.name)
			}
			plan := GenerateMigrationWithNormalizer(from, to, d, n)

			if _, err := conn.ExecContext(ctx, "SET SESSION sql_mode = "+quoteStringLiteral(deployMode)); err != nil {
				t.Fatalf("[%s] set deploy mode: %v", o.name, err)
			}
			if !o.applyPlan(t, conn, plan) {
				return
			}

			gotMode, gotTZ := eventSQLModeTZ(t, conn, ctx, db, "e")
			if gotMode != origMode {
				t.Fatalf("[%s] event sql_mode not preserved across ALTER: got %q want %q\nplan:\n%s", o.name, gotMode, origMode, plan.SQL())
			}
			// time_zone is captured at CREATE and not re-stamped by ALTER, so it should still be origTZ.
			if gotTZ != origTZ {
				t.Fatalf("[%s] event time_zone changed unexpectedly: got %q want %q", o.name, gotTZ, origTZ)
			}
		}))
	}
}

// mustLoadUnderDB loads a single bare CREATE under db name `db` via LoadSQL, matching the
// identity keying the oracle uses; it fails the test on a load error.
func mustLoadUnderDB(t testing.TB, db, ddl string, version Version) *Catalog {
	t.Helper()
	c, err := LoadSQLWithVersionHelper(fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s;\nUSE %s;\n%s", db, serverCharsetFor(version), db, ddl), version)
	if err != nil {
		t.Fatalf("load under db %s: %v", db, err)
	}
	return c
}

// LoadSQLWithVersionHelper loads imperative SQL and fixes the catalog version (LoadSQL has no
// version-taking variant; the version only affects canonicalization, not loading).
func LoadSQLWithVersionHelper(sqlText string, version Version) (*Catalog, error) {
	c, err := LoadSQL(sqlText)
	if err != nil {
		return nil, err
	}
	c.SetVersion(version)
	return c, nil
}

// routineSQLMode reads information_schema.ROUTINES.SQL_MODE for a routine.
func routineSQLMode(t testing.TB, conn oracleSession, ctx context.Context, db, name string) string {
	t.Helper()
	var mode sql.NullString
	row := conn.QueryRowContext(ctx,
		"SELECT SQL_MODE FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA=? AND ROUTINE_NAME=?", db, name)
	if err := row.Scan(&mode); err != nil {
		t.Fatalf("read routine sql_mode: %v", err)
	}
	return mode.String
}

// eventSQLModeTZ reads information_schema.EVENTS.SQL_MODE and TIME_ZONE for an event.
func eventSQLModeTZ(t testing.TB, conn oracleSession, ctx context.Context, db, name string) (string, string) {
	t.Helper()
	var mode, tz sql.NullString
	row := conn.QueryRowContext(ctx,
		"SELECT SQL_MODE, TIME_ZONE FROM information_schema.EVENTS WHERE EVENT_SCHEMA=? AND EVENT_NAME=?", db, name)
	if err := row.Scan(&mode, &tz); err != nil {
		t.Fatalf("read event sql_mode/time_zone: %v", err)
	}
	return mode.String, tz.String
}

// defaultCollationForServerCharset returns a plausible default collation for the server
// charset of a version, used only to populate a realistic synced context in these tests.
func defaultCollationForServerCharset(Version) string { return "utf8mb4_0900_ai_ci" }
