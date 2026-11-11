package catalog

import (
	"fmt"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// Version-dispatch tests for the OceanBase SDL catalog (node omni:version-dispatch).
//
// The invariant these lock in: a bare `CHARSET=utf8mb4` synced from OceanBase is dumped
// (by bytebase MetadataToSDL) as `DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci`
// (OceanBase's stored default), while the user's target DDL `... DEFAULT CHARSET=utf8mb4`
// (no COLLATE) is loaded carrying the same OceanBase default. They must resolve to the same
// collation — no phantom table-collation diff, no cascade into inherited varchar columns —
// and the generated DDL must never emit `utf8mb4_0900_ai_ci`, which OceanBase does not ship.
//
// The default collation is oracle-verified (information_schema.CHARACTER_SETS):
//   utf8mb4 -> utf8mb4_general_ci ; latin1 -> latin1_swedish_ci ; ascii -> ascii_general_ci

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
		syncedDump string // what MetadataToSDL emits for OceanBase's synced table
	}{
		{
			version:    OceanBase,
			syncedDump: "CREATE TABLE t (id INT NOT NULL, name VARCHAR(50) DEFAULT NULL, PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci",
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

// TestVersionDispatch_DefaultVersionIsOceanBase guards that a catalog with no version set
// canonicalizes as OceanBase (its Version zero value).
func TestVersionDispatch_DefaultVersionIsOceanBase(t *testing.T) {
	c := New()
	if c.Version() != OceanBase {
		t.Fatalf("New() catalog Version = %v, want OceanBase (back-compat)", c.Version())
	}
	cat, err := LoadSDL("CREATE DATABASE d; USE d; CREATE TABLE t (id INT PRIMARY KEY) DEFAULT CHARSET=utf8mb4;")
	if err != nil {
		t.Fatalf("LoadSDL: %v", err)
	}
	if cat.Version() != OceanBase {
		t.Fatalf("LoadSDL catalog Version = %v, want OceanBase", cat.Version())
	}
	// defaultNormalizer(to) must therefore be 8.0.
	if got := defaultNormalizer(cat).Version; got != OceanBase {
		t.Fatalf("defaultNormalizer version = %v, want OceanBase", got)
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

func versionName(_ Version) string {
	return "oceanbase"
}

// TestVersionDispatch_BareTimestampNoPhantom locks in the EDFT behavior: OceanBase's
// explicit_defaults_for_timestamp defaults ON, so a bare first TIMESTAMP stays nullable
// (no implicit NOT NULL DEFAULT CURRENT_TIMESTAMP magic) and the engine's stored form
// (nullable, no default) does not phantom-diff against the user's bare `TIMESTAMP`.
func TestVersionDispatch_BareTimestampNoPhantom(t *testing.T) {
	// OceanBase stored form of a bare TIMESTAMP under EDFT=ON: nullable, no default.
	synced := "CREATE TABLE t (id INT NOT NULL, ts TIMESTAMP NULL DEFAULT NULL, PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci"
	target := "CREATE TABLE t (id INT PRIMARY KEY, ts TIMESTAMP) DEFAULT CHARSET=utf8mb4"

	from, err := LoadSDLWithVersion("CREATE DATABASE d; USE d;\n"+synced, OceanBase)
	if err != nil {
		t.Fatalf("load synced: %v", err)
	}
	to, err := LoadSDLWithVersion("CREATE DATABASE d; USE d;\n"+target, OceanBase)
	if err != nil {
		t.Fatalf("load target: %v", err)
	}
	if from.Version() != OceanBase {
		t.Fatalf("from version = %v, want OceanBase", from.Version())
	}
	if d := Diff(from, to); !d.IsEmpty() {
		t.Errorf("bare-TIMESTAMP no-op diff not empty:\n  %s", describeDiff(d))
	}
}

// TestVersionDispatch_EDFTExplicitSetWins proves EDFT remains a SESSION variable, not a
// version trait: LoadSDLWithVersion seeds the session to OceanBase's box default (ON),
// but an explicit `SET explicit_defaults_for_timestamp=0` inside the SDL (processed after
// the seed) still wins. With EDFT OFF a bare first TIMESTAMP gets the implicit NOT NULL /
// DEFAULT CURRENT_TIMESTAMP magic.
func TestVersionDispatch_EDFTExplicitSetWins(t *testing.T) {
	// OceanBase catalog (box default EDFT ON) but the SDL carries SET ...=0 (OFF).
	sdl := "CREATE DATABASE d; USE d; SET SESSION explicit_defaults_for_timestamp=0;\n" +
		"CREATE TABLE t (id INT PRIMARY KEY, ts TIMESTAMP)"
	cat, err := LoadSDLWithVersion(sdl, OceanBase)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cat.session.ExplicitDefaultsForTimestamp {
		t.Fatalf("explicit SET ...=0 must win over the OceanBase box-default seed; session EDFT = true")
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
	// EDFT OFF: bare first TIMESTAMP gets the implicit NOT NULL / DEFAULT magic at load.
	if ts.Nullable {
		t.Errorf("with EDFT explicitly OFF, bare first TIMESTAMP must be NOT NULL, got nullable")
	}
	// And the diff-time normalizer (which reads the session EDFT) agrees it is NOT NULL.
	n := defaultNormalizer(cat)
	if n.ExplicitDefaultsForTimestamp != false {
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
			if strings.Contains(planSQL, "utf8mb4_0900_ai_ci") {
				t.Errorf("[%s] change DDL names utf8mb4_0900_ai_ci (OceanBase has no such collation):\n%s", o.name, planSQL)
			}

			ctx := oracleCtx(t)
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
