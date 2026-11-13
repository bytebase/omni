package catalog

import (
	"context"
	"fmt"
	"github.com/bytebase/omni/mysql/internal/testcontainer"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// Version-dispatch tests for the MySQL SDL catalog (node omni:version-dispatch).
//
// The catalog records a server version, and Diff/GenerateMigration canonicalize against that
// version's stored form. MySQL80 is the only version: MySQL 5.7 is end-of-life and its form is
// no longer reproduced. These lock in what the dispatch still owes on 8.0: a synced dump that
// spells the default collation out (`DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci`, as
// bytebase MetadataToSDL emits) and the user's bare `DEFAULT CHARSET=utf8mb4` do not
// phantom-diff, and the generated DDL names no collation for a bare charset.

// loadOneTableWithVersion loads a CREATE TABLE wrapped in a database with the given server
// charset and records the target version on the catalog (the bytebase wiring path:
// LoadSDLWithVersion). It returns the catalog + the named table.
func loadTableVersioned(t *testing.T, serverCharset, createSQL, table string, v Version) (*Catalog, *Table) {
	t.Helper()
	wrapped := fmt.Sprintf("CREATE DATABASE vdb DEFAULT CHARSET=%s;\nUSE vdb;\n%s", serverCharset, createSQL)
	cat, err := LoadSDLWithVersion(wrapped, v)
	if err != nil {
		t.Fatalf("LoadSDLWithVersion(%v) failed for %q: %v", v, createSQL, err)
	}
	var tbl *Table
	for _, db := range cat.Databases() {
		if tt := db.GetTable(table); tt != nil {
			tbl = tt
			break
		}
	}
	if tbl == nil {
		t.Fatalf("table %q not found after load of %q", table, createSQL)
	}
	return cat, tbl
}

// --- Unit proof (no live engine): the exact phantom-diff scenario --------------------

// TestVersionDispatch_BareCharsetNoPhantom is the headline unit proof. It takes the two
// input forms and asserts the diff is EMPTY:
//   - synced/stored form: bare CHARSET=utf8mb4 with the stored default COLLATE spelled out
//     (as MetadataToSDL emits).
//   - target/user form: bare CHARSET=utf8mb4 (no COLLATE) — loader bakes the 8.0 default.
func TestVersionDispatch_BareCharsetNoPhantom(t *testing.T) {
	const userTarget = "CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(50)) DEFAULT CHARSET=utf8mb4"

	cases := []struct {
		version    Version
		syncedDump string // what MetadataToSDL emits for this version's synced table
	}{
		{
			version:    MySQL80,
			syncedDump: "CREATE TABLE t (id INT NOT NULL, name VARCHAR(50) DEFAULT NULL, PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci",
		},
	}

	for _, c := range cases {
		v := c.version
		t.Run(versionName(v), func(t *testing.T) {
			n := NormalizerFor(v)
			// serverCharset must match the version box default so inheritance resolves the
			// same way on both sides — but the tables here declare their own utf8mb4 charset,
			// so it is moot; use the box default.
			sc := serverCharsetFor(v)

			fromCat, fromTbl := loadTableVersioned(t, sc, c.syncedDump, "t", v)
			toCat, toTbl := loadTableVersioned(t, sc, userTarget, "t", v)

			// effectiveTableCollation must collapse both forms onto the version default.
			wantColl := n.defaultCollationFor("utf8mb4")
			if got := n.effectiveTableCollation(fromTbl); got != wantColl {
				t.Errorf("[%s] synced-dump effectiveTableCollation = %q, want %q", versionName(v), got, wantColl)
			}
			if got := n.effectiveTableCollation(toTbl); got != wantColl {
				t.Errorf("[%s] user-target effectiveTableCollation = %q, want %q", versionName(v), got, wantColl)
			}

			// The full diff (both directions) must be EMPTY — no phantom table-option ALTER
			// and no cascaded column MODIFY.
			if d := DiffWithNormalizer(fromCat, toCat, n); !d.IsEmpty() {
				t.Errorf("[%s] PHANTOM DIFF synced->target not empty:\n  %s", versionName(v), describeDiff(d))
			}
			if d := DiffWithNormalizer(toCat, fromCat, n); !d.IsEmpty() {
				t.Errorf("[%s] PHANTOM DIFF target->synced not empty:\n  %s", versionName(v), describeDiff(d))
			}

			// And via the version-on-catalog entry point (Diff reads to.session.Version),
			// which is exactly how bytebase calls it.
			if d := Diff(fromCat, toCat); !d.IsEmpty() {
				t.Errorf("[%s] PHANTOM DIFF via Diff() (catalog version) not empty:\n  %s", versionName(v), describeDiff(d))
			}
		})
	}
}

// TestVersionDispatch_DefaultVersionIs80 guards the contract that a catalog with no version
// set canonicalizes as 8.0.
func TestVersionDispatch_DefaultVersionIs80(t *testing.T) {
	c := New()
	if c.Version() != MySQL80 {
		t.Fatalf("New() catalog Version = %v, want MySQL80 (back-compat)", c.Version())
	}
	cat, err := LoadSDL("CREATE DATABASE d; USE d; CREATE TABLE t (id INT PRIMARY KEY) DEFAULT CHARSET=utf8mb4;")
	if err != nil {
		t.Fatalf("LoadSDL: %v", err)
	}
	if cat.Version() != MySQL80 {
		t.Fatalf("LoadSDL catalog Version = %v, want MySQL80", cat.Version())
	}
	// defaultNormalizer(to) must therefore be 8.0.
	if got := defaultNormalizer(cat).Version; got != MySQL80 {
		t.Fatalf("defaultNormalizer version = %v, want MySQL80", got)
	}
}

// TestVersionDispatch_BareCharsetEmitsNoCollate proves the generate path names no collation
// for a bare-charset table/column: the rendered CREATE carries no COLLATE for a
// default-collation charset.
func TestVersionDispatch_BareCharsetEmitsNoCollate(t *testing.T) {
	// from: empty; to: a bare CHARSET=utf8mb4 table (CREATE TABLE rendering).
	for _, v := range []Version{MySQL80} {
		t.Run(versionName(v), func(t *testing.T) {
			n := NormalizerFor(v)
			sc := serverCharsetFor(v)
			toCat, _ := loadTableVersioned(t, sc,
				"CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(50)) DEFAULT CHARSET=utf8mb4", "t", v)
			fromCat := New()
			fromCat.SetVersion(v)

			diff := DiffWithNormalizer(fromCat, toCat, n)
			plan := GenerateMigrationWithNormalizer(fromCat, toCat, diff, n)
			sql := plan.SQL()
			// A bare CHARSET table renders DEFAULT CHARSET=utf8mb4 but NO COLLATE clause.
			if strings.Contains(strings.ToUpper(sql), "COLLATE") {
				t.Errorf("[%s] bare-charset CREATE should emit no COLLATE, got:\n%s", versionName(v), sql)
			}
		})
	}
}

func versionName(Version) string { return "8.0" }

// TestVersionDispatch_BareTimestampNoPhantom: LoadSDLWithVersion seeds the session
// explicit_defaults_for_timestamp to the server default (ON), so a bare TIMESTAMP stays
// nullable and the user's bare `TIMESTAMP` target does not phantom-diff against the stored
// form (nullable, no implicit default).
func TestVersionDispatch_BareTimestampNoPhantom(t *testing.T) {
	synced := "CREATE TABLE t (id INT NOT NULL, ts TIMESTAMP NULL DEFAULT NULL, PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci"
	target := "CREATE TABLE t (id INT PRIMARY KEY, ts TIMESTAMP) DEFAULT CHARSET=utf8mb4"

	from, err := LoadSDLWithVersion("CREATE DATABASE d; USE d;\n"+synced, MySQL80)
	if err != nil {
		t.Fatalf("load synced: %v", err)
	}
	to, err := LoadSDLWithVersion("CREATE DATABASE d; USE d;\n"+target, MySQL80)
	if err != nil {
		t.Fatalf("load target: %v", err)
	}
	if d := Diff(from, to); !d.IsEmpty() {
		t.Errorf("bare-TIMESTAMP no-op diff not empty:\n  %s", describeDiff(d))
	}
}

// TestVersionDispatch_EDFTExplicitSetWins proves EDFT remains a SESSION variable, not a
// version trait: LoadSDLWithVersion seeds the session to the server default (ON), but an
// explicit `SET explicit_defaults_for_timestamp` inside the SDL (processed after the seed)
// still wins. A schema that ran with EDFT explicitly OFF must therefore give its first bare
// TIMESTAMP the implicit NOT NULL.
func TestVersionDispatch_EDFTExplicitSetWins(t *testing.T) {
	sdl := "CREATE DATABASE d; USE d; SET SESSION explicit_defaults_for_timestamp=0;\n" +
		"CREATE TABLE t (id INT PRIMARY KEY, ts TIMESTAMP)"
	cat, err := LoadSDLWithVersion(sdl, MySQL80)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cat.session.ExplicitDefaultsForTimestamp {
		t.Fatalf("explicit SET ...=0 must win over the server-default seed; session EDFT = true")
	}
	var tbl *Table
	for _, db := range cat.Databases() {
		if tt := db.GetTable("t"); tt != nil {
			tbl = tt
		}
	}
	if tbl == nil {
		t.Fatal("table t not found")
	}
	ts := tbl.GetColumn("ts")
	if ts == nil {
		t.Fatal("column ts not found")
	}
	// The diff-time normalizer reads the session EDFT and agrees the column is NOT NULL.
	n := defaultNormalizer(cat)
	if n.ExplicitDefaultsForTimestamp {
		t.Errorf("defaultNormalizer EDFT = true, want false (session override)")
	}
	if !n.CanonicalNotNull(tbl, ts) {
		t.Errorf("CanonicalNotNull(bare TIMESTAMP) = false under EDFT OFF, want true (NOT NULL)")
	}
}

// TestVersionDispatch_ColumnDefaultUnderNonDefaultTable locks in the COLLATE-omission
// correctness fix (Codex blocking finding): when a table carries a NON-default collation
// and a column's resolved collation is the charset DEFAULT, the generated column MUST
// render an explicit COLLATE — omitting it would make the column wrongly inherit the
// table's non-default collation on apply (non-convergent). The omission is keyed off the
// table's effective collation (inheritance), not the charset default.
func TestVersionDispatch_ColumnDefaultUnderNonDefaultTable(t *testing.T) {
	for _, v := range []Version{MySQL80} {
		t.Run(versionName(v), func(t *testing.T) {
			n := NormalizerFor(v)
			defColl := n.defaultCollationFor("utf8mb4")
			// Table COLLATE=utf8mb4_unicode_ci (non-default); column at the charset default.
			create := fmt.Sprintf(
				"CREATE TABLE t (id INT PRIMARY KEY, a VARCHAR(20) CHARACTER SET utf8mb4 COLLATE %s) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci",
				defColl)
			toCat, _ := loadTableVersioned(t, serverCharsetFor(v), create, "t", v)
			fromCat := New()
			fromCat.SetVersion(v)
			diff := DiffWithNormalizer(fromCat, toCat, n)
			plan := GenerateMigrationWithNormalizer(fromCat, toCat, diff, n)
			sql := plan.SQL()
			// The column line must carry an explicit COLLATE for the charset default, because
			// the table collation differs from it.
			if !strings.Contains(sql, "COLLATE "+defColl) && !strings.Contains(sql, "COLLATE="+defColl) {
				t.Errorf("[%s] column at charset default under non-default table must emit explicit COLLATE %s, got:\n%s",
					versionName(v), defColl, sql)
			}
		})
	}
}

// --- Oracle proof (live engine): apply-clean on 8.0 -----------------------------------

// TestOracle_VersionDispatchBareCharset is the live spine for this node. For each version it:
//  1. applies the user target to a real DB, reads back the engine's stored form, and dumps a
//     "synced" SDL form (bare charset + the version's explicit stored default collation —
//     mirroring what bytebase MetadataToSDL produces);
//  2. asserts Diff(synced, target) is EMPTY under the version normalizer (idempotence);
//  3. generates the change DDL for a real modification, applies it to the real DB, and
//     asserts convergence (re-diff empty).
func TestOracle_VersionDispatchBareCharset(t *testing.T) {
	const userTarget = "CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(50)) DEFAULT CHARSET=utf8mb4"

	for _, v := range oracleVersions() {
		o := connectOracle(t, v)
		t.Run(o.name, func(t *testing.T) {
			// (1) engine's authentic stored form of the bare-charset target.
			readback, ok := o.showCreate(t, "vd_bare", userTarget, "t")
			if !ok {
				t.Fatalf("[%s] no readback for bare-charset target", o.name)
			}
			n := NormalizerFor(v)
			sc := serverCharsetFor(v)

			// The "synced" SDL is the engine's own SHOW CREATE — the most authentic stored
			// form, equivalent to (stricter than) MetadataToSDL's dump.
			fromCat, _ := loadTableVersioned(t, sc, readback, "t", v)
			toCat, _ := loadTableVersioned(t, sc, userTarget, "t", v)

			// (2) idempotence: no-op diff empty, both directions, and via catalog-version Diff.
			if d := DiffWithNormalizer(fromCat, toCat, n); !d.IsEmpty() {
				t.Errorf("[%s] no-op diff (stored vs target) not empty:\n  stored: %s\n  diff: %s",
					o.name, strings.TrimSpace(readback), describeDiff(d))
			}
			if d := Diff(fromCat, toCat); !d.IsEmpty() {
				t.Errorf("[%s] no-op Diff() (catalog version) not empty:\n  diff: %s", o.name, describeDiff(d))
			}

			// (3) apply-correctness of a real change: add an isbn column. Generate, apply to a
			// real DB seeded in the stored `from` state, and confirm it applies (no Error 1273)
			// and converges.
			targetPlus := "CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(50), isbn VARCHAR(20) NULL) DEFAULT CHARSET=utf8mb4"
			toPlusCat, _ := loadTableVersioned(t, sc, targetPlus, "t", v)
			diff := DiffWithNormalizer(fromCat, toPlusCat, n)
			if diff.IsEmpty() {
				t.Fatalf("[%s] expected a non-empty change diff for the added column", o.name)
			}
			plan := GenerateMigrationWithNormalizer(fromCat, toPlusCat, diff, n)
			planSQL := plan.SQL()

			ctx := context.Background()
			// A connection of its own: the test's other one stays free for the helpers.
			conn := testcontainer.OpenRoot(t, testcontainer.Options{MultiStatements: true})
			const applyDB = "vdb" // loadTableVersioned wraps in database `vdb`
			setup := []string{
				"DROP DATABASE IF EXISTS " + applyDB,
				fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s", applyDB, sc),
				"USE " + applyDB,
				userTarget,
			}
			for _, s := range setup {
				if _, err := conn.ExecContext(ctx, s); err != nil {
					t.Fatalf("[%s] setup failed: %q: %v", o.name, s, err)
				}
			}
			for _, op := range plan.Ops {
				if _, err := conn.ExecContext(ctx, op.SQL); err != nil {
					t.Fatalf("[%s] APPLY FAILED:\n  stmt: %s\n  err: %v\n  plan:\n%s",
						o.name, op.SQL, err, planSQL)
				}
			}
			// Converge: read back, reload, re-diff against target → empty.
			var nm, ddl string
			if err := conn.QueryRowContext(ctx, fmt.Sprintf("SHOW CREATE TABLE %s.t", applyDB)).Scan(&nm, &ddl); err != nil {
				t.Fatalf("[%s] readback after apply: %v", o.name, err)
			}
			resultCat, _ := loadTableVersioned(t, sc, ddl, "t", v)
			if d := DiffWithNormalizer(resultCat, toPlusCat, n); !d.IsEmpty() {
				t.Errorf("[%s] apply did not converge:\n  result: %s\n  diff: %s", o.name, strings.TrimSpace(ddl), describeDiff(d))
			}
		})
	}
}
