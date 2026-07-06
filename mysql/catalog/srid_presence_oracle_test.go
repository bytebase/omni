package catalog

import (
	"fmt"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// The oracle proof for column SRID presence (correctness-protocol.md gates 1 & 2), against
// the LIVE MySQL 8.0 engine. The spatial column `SRID` clause is 8.0-only (MySQL 5.7 has no
// column-level SRID), so every probe here is gated to 8.0; 5.7 is skipped, not proven.
//
// The bug this locks in: the catalog used to model SRID as a plain int where 0 doubled as
// "not set", gating every emit/diff site on `SRID != 0`. But MySQL keeps an explicit
// `SRID 0` as present-and-distinct from a column with no SRID clause: SHOW CREATE echoes
// `/*!80003 SRID 0 */`, and ST_GEOMETRY_COLUMNS.SRS_ID is 0 vs NULL. Modeling presence
// (Column.HasSRID) makes `SRID 0` round-trip and makes add/remove/change of an explicit
// `SRID 0` a real MODIFY instead of a silent no-op.
//
// Two properties, both mechanically proven:
//   - Idempotence (gate 1): user form vs the engine's SHOW CREATE readback diffs EMPTY, and
//     the stored form self-diffs / self-plans empty, for every SRID form (0, nonzero, none,
//     and with a spatial index).
//   - Presence distinctness + apply-correctness (gate 2): SRID 0 / none / 4326 produce
//     distinct canonical keys, and add/remove/change of the SRID clause generates a MODIFY
//     that, applied to a real from-state database, yields a to-equal table.
//
// The harness reuses connectOracle / showCreate / assertDiffEmptyAgainstReadback /
// assertApplyCorrect / loadColumn / serverCharsetFor / loadOneTable / both() / only() from
// the sibling oracle tests. It skips cleanly when the engine is unreachable (go test -short
// skips it entirely), so the unit suite stays hermetic.

// sridIdempotenceProbes enumerates the spatial-column FORMS that must round-trip EMPTY: an
// explicit `SRID 0`, a nonzero SRID, a no-SRID spatial column, and a spatial column with a
// SPATIAL index (SRID matters for spatial indexes). Each is 8.0-only.
func sridIdempotenceProbes() []diffProbe {
	return []diffProbe{
		// Explicit SRID 0 — the headline case. Must not collapse to no-SRID.
		{"srid-zero", "CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0)", "t", only(MySQL80)},
		// Nonzero SRID — the previously-working case must stay working.
		{"srid-4326", "CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 4326)", "t", only(MySQL80)},
		// No-SRID spatial column — must stay no-SRID (no phantom SRID 0 emitted).
		{"srid-none", "CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL)", "t", only(MySQL80)},
		// Nullable spatial column with an explicit SRID 0 (nullable geometry has no spatial
		// index requirement; proves presence is independent of nullability).
		{"srid-zero-nullable", "CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY SRID 0)", "t", only(MySQL80)},
		// Multiple spatial types carrying SRID 0 (POINT/LINESTRING/POLYGON) — the clause is a
		// column attribute, not specific to GEOMETRY.
		{"srid-zero-multi-type",
			"CREATE TABLE t (id INT PRIMARY KEY, p POINT NOT NULL SRID 0, l LINESTRING NOT NULL SRID 0, y POLYGON NOT NULL SRID 0)",
			"t", only(MySQL80)},
		// Spatial column WITH a SPATIAL index + SRID (SRID is meaningful for SPATIAL indexes,
		// which require a NOT NULL, SRID-constrained column). Must still round-trip.
		{"srid-zero-spatial-index",
			"CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0, SPATIAL INDEX sp_g (g))",
			"t", only(MySQL80)},
		{"srid-4326-spatial-index",
			"CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 4326, SPATIAL INDEX sp_g (g))",
			"t", only(MySQL80)},
		// Mixed table: a no-SRID spatial column alongside an explicit SRID 0 one — the two must
		// not smear into each other (one stays absent, one stays present-0).
		{"srid-mixed-none-and-zero",
			"CREATE TABLE t (id INT PRIMARY KEY, a GEOMETRY, b GEOMETRY NOT NULL SRID 0)",
			"t", only(MySQL80)},
	}
}

// TestOracle_SRIDPresenceIdempotence proves gate 1 for every SRID form: user DDL vs its
// engine readback diffs EMPTY, and the stored form self-diffs empty and self-plans empty,
// on MySQL 8.0.
func TestOracle_SRIDPresenceIdempotence(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range only(MySQL80) {
		o := connectOracle(t, version)
		n := NormalizerFor(version)
		for _, probe := range sridIdempotenceProbes() {
			if !containsVersion(probe.versions, version) {
				continue
			}
			t.Run(fmt.Sprintf("%s/%s", o.name, probe.id), func(t *testing.T) {
				db := "srid_" + strings.ReplaceAll(probe.id, "-", "_")
				// Gate 1a + 2 (canonicalization): user form vs stored readback empty, both directions.
				assertDiffEmptyAgainstReadback(t, o, db, probe.create, probe.table)

				// Gate 1b (the spine, plan level): the stored form's self-plan is EMPTY, and the
				// user-form-vs-stored round-trip plan is EMPTY (no phantom MODIFY re-rendering the
				// column without SRID).
				cat, _, rb, ok := o.catalogFromReadback(t, db+"_idem", probe.create, probe.table)
				if !ok {
					t.Skipf("[%s] could not obtain readback for %s", o.name, probe.id)
				}
				if plan := GenerateMigrationWithNormalizer(cat, cat, DiffWithNormalizer(cat, cat, n), n); plan.SQL() != "" {
					t.Errorf("[%s] NON-EMPTY NO-OP PLAN for %s (SRID presence normalization bug):\n  stored: %s\n  plan:\n%s",
						o.name, probe.id, strings.TrimSpace(rb), plan.SQL())
				}
				userCat, _ := loadOneTable(t, serverCharsetFor(o.version), probe.create, probe.table)
				if rtPlan := GenerateMigrationWithNormalizer(userCat, cat, DiffWithNormalizer(userCat, cat, n), n); rtPlan.SQL() != "" {
					t.Errorf("[%s] NON-EMPTY ROUND-TRIP PLAN (user vs stored) for %s:\n  user:   %s\n  stored: %s\n  plan:\n%s",
						o.name, probe.id, strings.TrimSpace(probe.create), strings.TrimSpace(rb), rtPlan.SQL())
				}
			})
		}
	}
}

// TestOracle_SRIDPresenceDistinctness is the dual of idempotence: SRID 0, no-SRID, and a
// nonzero SRID are three genuinely different MySQL schemas and MUST produce three distinct
// canonical keys — otherwise add/remove/change of an explicit SRID 0 is an invisible no-op.
// Each pair is loaded from the real engine's SHOW CREATE so the stored forms are authentic.
func TestOracle_SRIDPresenceDistinctness(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range only(MySQL80) {
		o := connectOracle(t, version)
		sc := serverCharsetFor(o.version)
		n := NormalizerFor(o.version)

		// Apply all three forms in one database and read them back, so each column's stored
		// form is authentic and comparable (same table charset, same engine).
		ddl := "CREATE TABLE t (id INT PRIMARY KEY, gz GEOMETRY NOT NULL SRID 0, gn GEOMETRY NOT NULL, gf GEOMETRY NOT NULL SRID 4326)"
		rb, ok := o.showCreate(t, "srid_distinct", ddl, "t")
		if !ok {
			t.Skipf("[%s] could not obtain readback", o.name)
		}
		tbl, zCol := loadColumn(t, sc, rb, "t", "gz")
		_, nCol := loadColumn(t, sc, rb, "t", "gn")
		_, fCol := loadColumn(t, sc, rb, "t", "gf")
		if zCol == nil || nCol == nil || fCol == nil {
			t.Fatalf("[%s] missing columns (z=%v n=%v f=%v)", o.name, zCol != nil, nCol != nil, fCol != nil)
		}

		kz := n.CanonicalColumn(tbl, zCol)
		kn := n.CanonicalColumn(tbl, nCol)
		kf := n.CanonicalColumn(tbl, fCol)

		t.Run(o.name+"/srid-zero-vs-none", func(t *testing.T) {
			// The headline distinctness: SRID 0 (present, value 0) must differ from no-SRID.
			if kz == kn {
				t.Errorf("[%s] SRID 0 and no-SRID must produce DIFFERENT keys (else add/remove SRID 0 is a no-op):\n  srid0: %s\n  none:  %s",
					o.name, kz, kn)
			}
		})
		t.Run(o.name+"/srid-zero-vs-4326", func(t *testing.T) {
			if kz == kf {
				t.Errorf("[%s] SRID 0 and SRID 4326 must produce DIFFERENT keys:\n  srid0:    %s\n  srid4326: %s", o.name, kz, kf)
			}
		})
		t.Run(o.name+"/srid-none-vs-4326", func(t *testing.T) {
			if kn == kf {
				t.Errorf("[%s] no-SRID and SRID 4326 must produce DIFFERENT keys:\n  none:     %s\n  srid4326: %s", o.name, kn, kf)
			}
		})
	}
}

// sridMigrationProbes enumerates the add / remove / change transitions of the SRID clause,
// each of which must generate a MODIFY that applies to a real from-state database and yields
// a to-equal table. 8.0-only. These reuse the generate-core apply-correctness harness
// (assertApplyCorrect), which loads from/to from the engine's own readbacks, generates the
// plan, applies it, and diffs the result against `to` in both directions.
func sridMigrationProbes() []migrationProbe {
	return []migrationProbe{
		// ---- ADD SRID (none -> value) ----
		{"add-srid-none-to-zero", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL)",
			"CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0)", only(MySQL80)},
		{"add-srid-none-to-4326", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL)",
			"CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 4326)", only(MySQL80)},
		// ---- REMOVE SRID (value -> none) ----
		{"remove-srid-zero-to-none", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0)",
			"CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL)", only(MySQL80)},
		{"remove-srid-4326-to-none", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 4326)",
			"CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL)", only(MySQL80)},
		// ---- CHANGE SRID (value -> different value, incl. to/from 0) ----
		{"change-srid-4326-to-zero", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 4326)",
			"CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0)", only(MySQL80)},
		{"change-srid-zero-to-4326", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0)",
			"CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 4326)", only(MySQL80)},
		// ---- ADD spatial column carrying SRID 0 (from empty table) ----
		{"add-column-with-srid-zero", "t",
			"CREATE TABLE t (id INT PRIMARY KEY)",
			"CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0)", only(MySQL80)},
		// ---- CREATE table with a SRID-0 column + spatial index (from empty) ----
		{"create-srid-zero-spatial-index", "t", "",
			"CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0, SPATIAL INDEX sp_g (g))", only(MySQL80)},
	}
}

// TestOracle_SRIDPresenceApplyCorrectness proves gate 2 (apply-correctness) for every SRID
// add/remove/change: the generated DDL transforms a real from-state database into a to-equal
// one, on MySQL 8.0. A non-empty result diff (or a failed apply) is a bug.
func TestOracle_SRIDPresenceApplyCorrectness(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range only(MySQL80) {
		o := connectOracle(t, version)
		n := NormalizerFor(version)
		for _, probe := range sridMigrationProbes() {
			if !containsVersion(probe.versions, version) {
				continue
			}
			t.Run(fmt.Sprintf("%s/%s", o.name, probe.id), func(t *testing.T) {
				assertApplyCorrect(t, o, n, probe)
			})
		}
	}
}

// TestOracle_SRIDPresenceEmitsClause is the dedicated readback proof: for an add/change TO a
// specific SRID, the generated plan text carries the exact `SRID <n>` clause (including
// `SRID 0`), and for a remove the plan drops it. This complements apply-correctness (which
// could in principle be satisfied by two equally-wrong sides) by asserting the emitted DDL
// itself, and directly pins the "MODIFY re-renders WITHOUT SRID 0" regression: a remove must
// emit no SRID clause, an add of SRID 0 must emit `SRID 0`.
func TestOracle_SRIDPresenceEmitsClause(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range only(MySQL80) {
		o := connectOracle(t, version)
		n := NormalizerFor(version)
		sc := serverCharsetFor(o.version)

		cases := []struct {
			id          string
			fromDDL     string
			toDDL       string
			wantInPlan  string // substring the plan MUST contain ("" = skip)
			wantNotPlan string // substring the plan MUST NOT contain ("" = skip)
		}{
			// none -> SRID 0 : the plan must render the explicit SRID 0 clause.
			{"add-zero", "CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL)",
				"CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0)", "SRID 0", ""},
			// none -> SRID 4326.
			{"add-4326", "CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL)",
				"CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 4326)", "SRID 4326", ""},
			// SRID 0 -> none : the plan must NOT carry any SRID clause (the regression case).
			{"remove-zero", "CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0)",
				"CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL)", "", "SRID"},
			// SRID 4326 -> SRID 0 : must emit SRID 0, not drop to no-clause.
			{"change-to-zero", "CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 4326)",
				"CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0)", "SRID 0", ""},
		}
		for _, c := range cases {
			t.Run(o.name+"/"+c.id, func(t *testing.T) {
				// Load from/to from the engine's own readbacks so the catalog is authentic.
				rbFrom, ok1 := o.showCreate(t, "srid_emit_f", c.fromDDL, "t")
				rbTo, ok2 := o.showCreate(t, "srid_emit_t", c.toDDL, "t")
				if !ok1 || !ok2 {
					t.Skipf("[%s] could not obtain readbacks", o.name)
				}
				fromCat, _ := loadOneTable(t, sc, rbFrom, "t")
				toCat, _ := loadOneTable(t, sc, rbTo, "t")
				diff := DiffWithNormalizer(fromCat, toCat, n)
				plan := GenerateMigrationWithNormalizer(fromCat, toCat, diff, n)
				sql := plan.SQL()
				if sql == "" {
					t.Fatalf("[%s] expected a MODIFY plan for %s, got empty (SRID change went undetected)", o.name, c.id)
				}
				if c.wantInPlan != "" && !strings.Contains(sql, c.wantInPlan) {
					t.Errorf("[%s] %s: plan must contain %q:\n%s", o.name, c.id, c.wantInPlan, sql)
				}
				if c.wantNotPlan != "" && strings.Contains(sql, c.wantNotPlan) {
					t.Errorf("[%s] %s: plan must NOT contain %q (SRID leaked on a remove):\n%s", o.name, c.id, c.wantNotPlan, sql)
				}
			})
		}
	}
}
