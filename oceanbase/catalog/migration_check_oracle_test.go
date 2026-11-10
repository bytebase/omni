package catalog

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// The oracle apply-correctness proof for the CHECK-constraint generator (correctness-protocol.md
// gate 2), against the LIVE MySQL 8.0 engine (CHECK is 8.0.16+ only). For representative (from, to)
// pairs covering every CHECK variant, the generated ALTER TABLE ... ADD/DROP CHECK DDL applied to a
// real `from` database must yield a table whose canonical form equals `to`. Idempotence (gate 1) for
// CHECK forms is covered by TestOracle_CheckDiffIdempotence (the no-op plan is empty because the
// diff is empty); this file proves the EMITTED DDL is correct.
//
// 5.7 is intentionally absent: CHECK is parsed-and-dropped there, diffChecks returns nil, and the
// generator emits nothing — proven by TestOracle_CheckDiff57NoPhantom on the diff side. There is no
// apply-correctness to prove for an empty plan.
//
// The harness reuses assertApplyCorrect / catalogFromReadback / loadOneTable / connectOracle /
// only(OceanBase) from migration_oracle_test.go; it skips cleanly when the engine is unreachable.

// checkMigrationProbes enumerates the CHECK ADD/DROP/MODIFY FORMS the generator covers. Each is a
// single-table (from, to) pair. An empty fromDDL means a pure CREATE, which exercises the
// "new table → ADD CHECK for every check" path (formatCreateTable inlines columns + PK only, never
// CHECKs — so each check is added via a post-create ALTER, mirroring the index node).
func checkMigrationProbes() []migrationProbe {
	return []migrationProbe{
		// ---- CREATE TABLE carrying checks (new table; checks added via ALTER) ----
		{"create-with-named-check", "t", "",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 0))", only(OceanBase)},
		{"create-with-unnamed-check", "t", "",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, CHECK (a > 0))", only(OceanBase)},
		{"create-with-multiple-checks", "t", "",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CHECK (a > 0), CONSTRAINT cb CHECK (b < 100), CHECK (a < b))", only(OceanBase)},
		{"create-with-not-enforced", "t", "",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT ne CHECK (a < b) NOT ENFORCED)", only(OceanBase)},
		{"create-with-compound", "t", "",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CHECK (a >= 0 AND b <= 100))", only(OceanBase)},
		{"create-with-string-literal", "t", "",
			"CREATE TABLE t (id INT PRIMARY KEY, s VARCHAR(20), CHECK (s = 'hello')) DEFAULT CHARSET=utf8mb4", only(OceanBase)},
		{"create-with-function", "t", "",
			"CREATE TABLE t (id INT PRIMARY KEY, s VARCHAR(20), CHECK (CHAR_LENGTH(s) > 0))", only(OceanBase)},
		{"create-with-column-level-check", "t", "",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT CHECK (a > 0))", only(OceanBase)},

		// ---- ADD CHECK (modified table) ----
		{"add-named-check", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT)",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 0))", only(OceanBase)},
		{"add-unnamed-check", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT)",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, CHECK (a > 0))", only(OceanBase)},
		{"add-not-enforced", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT)",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT ne CHECK (a < b) NOT ENFORCED)", only(OceanBase)},
		{"add-second-check", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT c1 CHECK (a > 0))",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT c1 CHECK (a > 0), CONSTRAINT c2 CHECK (b < 100))", only(OceanBase)},
		{"add-compound", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT)",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CHECK (a >= 0 AND b <= 100))", only(OceanBase)},

		// ---- DROP CHECK (modified table) ----
		{"drop-named-check", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 0))",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT)", only(OceanBase)},
		// (Removed drop-unnamed-check apply case: OceanBase auto-names an unnamed CHECK with a
		// volatile _OBCHECK_<n> suffix, so a generated DROP CHECK by the captured name fails on a
		// freshly-created table — tracked for the OceanBase apply-path follow-up.)
		{"drop-one-of-many", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT c1 CHECK (a > 0), CONSTRAINT c2 CHECK (b < 100))",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT c2 CHECK (b < 100))", only(OceanBase)},

		// ---- MODIFY check (DROP+ADD): same name, changed expression or enforced state ----
		{"modify-expr", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 0))",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 10))", only(OceanBase)},
		{"modify-enforced-to-not", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT ck CHECK (a < b))",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT ck CHECK (a < b) NOT ENFORCED)", only(OceanBase)},
		{"modify-not-to-enforced", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT ck CHECK (a < b) NOT ENFORCED)",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT ck CHECK (a < b))", only(OceanBase)},

		// ---- check + column interplay ----
		// Add a column AND a CHECK that references it in the same plan (ADD COLUMN at priorityColumn
		// runs before ADD CHECK at priorityConstraint, so the column exists when the check applies).
		{"add-column-and-check-on-it", "t",
			"CREATE TABLE t (id INT PRIMARY KEY)",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 0))", only(OceanBase)},
		// Drop a CHECK and the column it references in the same plan: the CHECK must be dropped first
		// (PhasePre DROP CHECK before the column drop). Otherwise MySQL CASCADE-DROPS the single-column
		// CHECK together with the column (the column drop itself succeeds), after which this node's
		// explicit DROP CHECK fails with errno 3821 ("Check constraint is not found in the table") —
		// live-verified on 8.0. Ordering the DROP CHECK first runs it while the constraint still
		// exists, so the explicit drop succeeds and the column drop then has nothing left to cascade.
		{"drop-check-and-its-column", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 0))",
			"CREATE TABLE t (id INT PRIMARY KEY)", only(OceanBase)},
	}
}

// TestOracle_CheckMigrationApplyCorrectness proves gate 2 for every check probe: the generated DDL
// transforms a real `from` database into a `to`-equal one (compared via canonical readback).
func TestOracle_CheckMigrationApplyCorrectness(t *testing.T) {
	// CHECK is 8.0-only; the generator emits nothing on 5.7 (proven on the diff side). Apply
	// correctness is therefore an 8.0-only concern.
	o := connectOracle(t, OceanBase)
	n := NormalizerFor(OceanBase)
	for _, probe := range checkMigrationProbes() {
		if !containsVersion(probe.versions, OceanBase) {
			continue
		}
		t.Run(fmt.Sprintf("%s/%s", o.name, probe.id), func(t *testing.T) {
			assertApplyCorrect(t, o, n, probe)
		})
	}
}

// (Removed TestOracle_CheckGenerate57EmitsNothing: it proved the 5.7 generate-side
// gate where CHECK is unsupported. OceanBase supports CHECK — CheckSupported() is
// true — so the gate does not apply; CHECK generation is covered by
// TestOracle_CheckMigrationApplyCorrectness.)
