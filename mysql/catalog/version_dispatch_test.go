package catalog

import (
	"context"
	"fmt"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// Version-dispatch tests for the MySQL SDL catalog (node omni:version-dispatch).
//
// The bug these lock in: the catalog had no notion of server version, so the diff/generate
// canonicalizer always assumed 8.0. A bare `CHARSET=utf8mb4` synced from a 5.7 server is
// dumped (by bytebase MetadataToSDL) as `DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci`
// (the 5.7 stored default), while the user's target DDL `... DEFAULT CHARSET=utf8mb4` (no
// COLLATE) is loaded carrying the loader-baked STATIC 8.0 default `utf8mb4_0900_ai_ci`. Under
// an 8.0-only normalizer the two resolve to different collations -> a phantom table-collation
// diff that cascades into inherited varchar columns, and the generated DDL emits
// `COLLATE=utf8mb4_0900_ai_ci` — a collation that DOES NOT EXIST on 5.7 (Error 1273).
//
// The version-correct defaults are oracle-verified (SHOW COLLATION ... Default='Yes'):
//   utf8mb4 -> utf8mb4_general_ci (5.7) vs utf8mb4_0900_ai_ci (8.0)
//   utf8/utf8mb3 -> utf8mb3_general_ci ; latin1 -> latin1_swedish_ci (both stable)

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

// TestVersionDispatch_BareCharsetNoPhantom is the headline unit proof. It reproduces the
// bug's two input forms and asserts the diff is EMPTY under BOTH versions' normalizers:
//   - synced/stored form: bare CHARSET=utf8mb4 with the version's stored default COLLATE
//     spelled out (as MetadataToSDL emits), differing per version.
//   - target/user form: bare CHARSET=utf8mb4 (no COLLATE) — loader bakes the static 8.0
//     default regardless of version.
func TestVersionDispatch_BareCharsetNoPhantom(t *testing.T) {
	const userTarget = "CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(50)) DEFAULT CHARSET=utf8mb4"

	cases := []struct {
		version    Version
		syncedDump string // what MetadataToSDL emits for this version's synced table
	}{
		{
			version:    MySQL57,
			syncedDump: "CREATE TABLE t (id INT NOT NULL, name VARCHAR(50) DEFAULT NULL, PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci",
		},
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
			// same way on both sides (utf8mb4 on the 8.0 box, latin1 on the 5.7 box) — but the
			// tables here declare their own utf8mb4 charset, so it is moot; use the box default.
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

// TestVersionDispatch_DefaultVersionIs80 guards the back-compat contract: a catalog with no
// version set canonicalizes as 8.0 (Version's zero value is MySQL57, so defaultSessionState
// must force MySQL80).
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

// TestVersionDispatch_GeneratedDDLNoMissingCollation proves the generate path never emits the
// 5.7-missing collation for a bare-charset table/column: the rendered CREATE/ALTER carries no
// COLLATE for a default-collation charset, so it is valid on either version.
func TestVersionDispatch_GeneratedDDLNoMissingCollation(t *testing.T) {
	// from: empty; to: a bare CHARSET=utf8mb4 table (CREATE TABLE rendering).
	for _, v := range both() {
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
			if strings.Contains(sql, "utf8mb4_0900_ai_ci") {
				t.Errorf("[%s] generated DDL names utf8mb4_0900_ai_ci (invalid on 5.7):\n%s", versionName(v), sql)
			}
			// A bare CHARSET table renders DEFAULT CHARSET=utf8mb4 but NO COLLATE clause.
			if strings.Contains(strings.ToUpper(sql), "COLLATE") {
				t.Errorf("[%s] bare-charset CREATE should emit no COLLATE, got:\n%s", versionName(v), sql)
			}
		})
	}
}

func versionName(v Version) string {
	if v == MySQL80 {
		return "8.0"
	}
	return "5.7"
}

// TestVersionDispatch_BareTimestampNoPhantom locks in the EDFT half of the 5.7 fix:
// LoadSDLWithVersion seeds the session explicit_defaults_for_timestamp to the version's
// box default (OFF on 5.7), so a bare TIMESTAMP synced from 5.7 (which the server
// materializes to NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP) does
// not phantom-diff against the user's bare `TIMESTAMP` target. Without the EDFT seed the
// 5.7 catalog would keep New()'s 8.0 default (ON) and the bare TIMESTAMP would diff.
func TestVersionDispatch_BareTimestampNoPhantom(t *testing.T) {
	// 5.7 server stored form of a bare first TIMESTAMP under EDFT=OFF.
	synced57 := "CREATE TABLE t (id INT NOT NULL, ts TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci"
	target := "CREATE TABLE t (id INT PRIMARY KEY, ts TIMESTAMP) DEFAULT CHARSET=utf8mb4"

	from, err := LoadSDLWithVersion("CREATE DATABASE d; USE d;\n"+synced57, MySQL57)
	if err != nil {
		t.Fatalf("load synced: %v", err)
	}
	to, err := LoadSDLWithVersion("CREATE DATABASE d; USE d;\n"+target, MySQL57)
	if err != nil {
		t.Fatalf("load target: %v", err)
	}
	if from.Version() != MySQL57 {
		t.Fatalf("from version = %v, want MySQL57", from.Version())
	}
	if d := Diff(from, to); !d.IsEmpty() {
		t.Errorf("5.7 bare-TIMESTAMP no-op diff not empty:\n  %s", describeDiff(d))
	}

	// And on 8.0 a bare TIMESTAMP stays nullable (EDFT ON), so the same user target vs an
	// 8.0 stored form (nullable, no implicit default) is also empty.
	synced80 := "CREATE TABLE t (id INT NOT NULL, ts TIMESTAMP NULL DEFAULT NULL, PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci"
	from80, err := LoadSDLWithVersion("CREATE DATABASE d; USE d;\n"+synced80, MySQL80)
	if err != nil {
		t.Fatalf("load synced80: %v", err)
	}
	to80, err := LoadSDLWithVersion("CREATE DATABASE d; USE d;\n"+target, MySQL80)
	if err != nil {
		t.Fatalf("load target80: %v", err)
	}
	if d := Diff(from80, to80); !d.IsEmpty() {
		t.Errorf("8.0 bare-TIMESTAMP no-op diff not empty:\n  %s", describeDiff(d))
	}
}

// TestVersionDispatch_EDFTExplicitSetWins proves EDFT remains a SESSION variable, not a
// version trait: LoadSDLWithVersion seeds the session to the version box default, but an
// explicit `SET explicit_defaults_for_timestamp` inside the SDL (processed after the seed)
// still wins. A 5.7 schema that ran with EDFT explicitly ON must therefore leave a bare
// TIMESTAMP nullable (no implicit NOT NULL / DEFAULT CURRENT_TIMESTAMP magic).
func TestVersionDispatch_EDFTExplicitSetWins(t *testing.T) {
	// 5.7 catalog (box default EDFT OFF) but the SDL carries SET ...=1 (ON).
	sdl := "CREATE DATABASE d; USE d; SET SESSION explicit_defaults_for_timestamp=1;\n" +
		"CREATE TABLE t (id INT PRIMARY KEY, ts TIMESTAMP)"
	cat, err := LoadSDLWithVersion(sdl, MySQL57)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if !cat.session.ExplicitDefaultsForTimestamp {
		t.Fatalf("explicit SET ...=1 must win over the 5.7 box-default seed; session EDFT = false")
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
	// EDFT ON: bare TIMESTAMP stays nullable, no implicit default materialized at load.
	if !ts.Nullable {
		t.Errorf("with EDFT explicitly ON, bare TIMESTAMP must stay nullable, got NOT NULL")
	}
	if ts.Default != nil {
		t.Errorf("with EDFT explicitly ON, bare TIMESTAMP must have no implicit default, got %q", *ts.Default)
	}
	// And the diff-time normalizer (which reads the session EDFT) agrees it is nullable.
	n := defaultNormalizer(cat)
	if n.ExplicitDefaultsForTimestamp != true {
		t.Errorf("defaultNormalizer EDFT = false, want true (session override)")
	}
	if n.CanonicalNotNull(tbl, ts) {
		t.Errorf("CanonicalNotNull(bare TIMESTAMP) = true under EDFT ON, want false (nullable)")
	}
}

// TestVersionDispatch_ColumnDefaultUnderNonDefaultTable locks in the COLLATE-omission
// correctness fix (Codex blocking finding): when a table carries a NON-default collation
// and a column's resolved collation is the charset DEFAULT, the generated column MUST
// render an explicit COLLATE — omitting it would make the column wrongly inherit the
// table's non-default collation on apply (non-convergent). The omission is keyed off the
// table's effective collation (inheritance), not the charset default.
func TestVersionDispatch_ColumnDefaultUnderNonDefaultTable(t *testing.T) {
	for _, v := range both() {
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

// --- Oracle proof (live engines): apply-clean on 5.7 and 8.0 -------------------------

// TestOracle_VersionDispatchBareCharset is the live spine for this node. For each version it:
//  1. applies the user target to a real DB, reads back the engine's stored form, and dumps a
//     "synced" SDL form (bare charset + the version's explicit stored default collation —
//     mirroring what bytebase MetadataToSDL produces);
//  2. asserts Diff(synced, target) is EMPTY under the version normalizer (idempotence);
//  3. generates the change DDL for a real modification, applies it to the real 5.7/8.0 DB
//     (proving no Error 1273), and asserts convergence (re-diff empty).
func TestOracle_VersionDispatchBareCharset(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	const userTarget = "CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(50)) DEFAULT CHARSET=utf8mb4"

	for _, v := range both() {
		o := connectOracle(t, v)
		t.Run(o.name, func(t *testing.T) {
			// (1) engine's authentic stored form of the bare-charset target.
			readback, ok := o.showCreate(t, "vd_bare", userTarget, "t")
			if !ok {
				t.Skipf("[%s] no readback for bare-charset target", o.name)
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
			if strings.Contains(planSQL, "utf8mb4_0900_ai_ci") && v == MySQL57 {
				t.Errorf("[%s] change DDL names utf8mb4_0900_ai_ci (invalid on 5.7):\n%s", o.name, planSQL)
			}

			ctx := context.Background()
			conn, err := o.db.Conn(ctx)
			if err != nil {
				t.Fatalf("[%s] conn: %v", o.name, err)
			}
			defer func() { _ = conn.Close() }()
			const applyDB = "vdb" // loadTableVersioned wraps in database `vdb`
			setup := []string{
				"DROP DATABASE IF EXISTS " + applyDB,
				fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s", applyDB, sc),
				"USE " + applyDB,
				userTarget,
			}
			for _, s := range setup {
				if _, err := conn.ExecContext(ctx, s); err != nil {
					t.Skipf("[%s] setup failed: %q: %v", o.name, s, err)
				}
			}
			for _, op := range plan.Ops {
				if _, err := conn.ExecContext(ctx, op.SQL); err != nil {
					t.Fatalf("[%s] APPLY FAILED (Error here on 5.7 = the bug):\n  stmt: %s\n  err: %v\n  plan:\n%s",
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
