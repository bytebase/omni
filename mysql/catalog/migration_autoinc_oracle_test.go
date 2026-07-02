package catalog

import (
	"fmt"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// The oracle proof for the AUTO_INCREMENT / backing-key ALTER grouping (migration_autoinc.go),
// against the LIVE MySQL engines (5.7.25 and 8.0.32 — identical errno-1075 behavior, probed
// 2026-07). MySQL enforces, at the END OF EVERY STATEMENT, that a table has at most one
// AUTO_INCREMENT column and that this column is the FIRST column of some index (errno 1075;
// MyISAM alone also accepts a non-first position). A plan that adds an AUTO_INCREMENT column in
// one ALTER and its backing key in the next therefore fails on the FIRST statement — the column
// and its key must land in the SAME statement. These probes cover every grouping form the merge
// pass emits (correctness-protocol.md gate 2: each generated plan, applied to a real `from`
// database, must yield `to`), plus the no-merge control cases that must keep working ungrouped.
//
// The idempotence gate (gate 1) for the same forms is TestOracle_AutoIncGroupingIdempotence:
// every `to` schema, read back from the engine, must self-plan to an empty migration.
//
// The harness reuses assertApplyCorrect / connectOracle / both / only from
// migration_oracle_test.go and skips cleanly when the engines are unreachable.
func autoIncMigrationProbes() []migrationProbe {
	return []migrationProbe{
		// ---- ADD direction: AUTO_INCREMENT column arrives without a pre-existing key ----

		// The headline bug: adding an AUTO_INCREMENT column plus its UNIQUE backing key to an
		// existing table. Ungrouped, the ADD COLUMN statement fails errno 1075.
		{"aig-add-col-uk", "t",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT)",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT, seq BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_seq (seq))", both()},
		// A plain (non-unique) KEY also satisfies errno 1075.
		{"aig-add-col-plain-key", "t",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT)",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT, seq BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, KEY k_seq (seq))", both()},
		// AUTO_INCREMENT column + new PRIMARY KEY on a table that had none.
		{"aig-add-col-new-pk", "t",
			"CREATE TABLE t (id INT NOT NULL, a INT)",
			"CREATE TABLE t (id INT NOT NULL, a INT, seq BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, PRIMARY KEY (seq))", both()},
		// Surrogate-key promotion: the new AUTO_INCREMENT column REPLACES an existing natural PK.
		// The column add must join the combined DROP PRIMARY KEY / ADD PRIMARY KEY statement.
		{"aig-add-col-promote-pk", "t",
			"CREATE TABLE t (a INT NOT NULL PRIMARY KEY, b INT)",
			"CREATE TABLE t (a INT NOT NULL, b INT, seq BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, PRIMARY KEY (seq))", both()},
		// Surrogate-key promotion where the old PK member is ALSO dropped (forces the split PK
		// path): the column add merges with the standalone ADD PRIMARY KEY instead.
		{"aig-add-col-promote-pk-drop-old", "t",
			"CREATE TABLE t (a INT NOT NULL PRIMARY KEY, b INT)",
			"CREATE TABLE t (b INT, seq BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, PRIMARY KEY (seq))", both()},
		// An EXISTING column gains AUTO_INCREMENT while its backing key arrives in the same plan.
		{"aig-modify-gain-ai-new-key", "t",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL)",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_c (c))", both()},
		// CONTROL: gaining AUTO_INCREMENT on a column that is ALREADY keyed needs no grouping and
		// must keep working as a single MODIFY.
		{"aig-modify-gain-ai-existing-key", "t",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL, UNIQUE KEY uk_c (c))",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_c (c))", both()},
		// Composite backing key whose OTHER member is also added in this plan: the second column's
		// add is pulled into the grouped statement so the key can reference it.
		{"aig-add-col-composite-backing", "t",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY)",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, seq BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, d INT NOT NULL, UNIQUE KEY uk_sd (seq, d))", both()},
		// AUTO_INCREMENT migrates between EXISTING columns (both already keyed): the old column's
		// de-AUTO_INCREMENT MODIFY must share the statement that makes the new column
		// AUTO_INCREMENT (errno 1075's "only one auto column" half fires between two statements).
		{"aig-ai-migrate-columns", "t",
			"CREATE TABLE t (a INT NOT NULL AUTO_INCREMENT, b INT NOT NULL, UNIQUE KEY uk_a (a), UNIQUE KEY uk_b (b))",
			"CREATE TABLE t (a INT NOT NULL, b INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_a (a), UNIQUE KEY uk_b (b))", both()},
		// AUTO_INCREMENT migrates from an existing column to a NEW column (new backing key too).
		{"aig-ai-migrate-to-new-col", "t",
			"CREATE TABLE t (a INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_a (a))",
			"CREATE TABLE t (a INT NOT NULL, b INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_a (a), UNIQUE KEY uk_b (b))", both()},
		// AUTO_INCREMENT column + backing key grouped while UNRELATED changes ride along in the
		// same plan (they must stay separate statements, unharmed).
		{"aig-add-col-with-other-changes", "t",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT, e VARCHAR(10))",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT, e VARCHAR(20), f INT, seq BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_seq (seq))", both()},
		// MyISAM accepts a NON-first position in a multi-column key as the backing (probed: OK on
		// MyISAM, errno 1075 on InnoDB/MEMORY) — the grouping must still fire for it.
		{"aig-myisam-add-nonleftmost", "t",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT NOT NULL) ENGINE=MyISAM",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, seq INT NOT NULL AUTO_INCREMENT, KEY k_as (a, seq)) ENGINE=MyISAM", both()},
		// CONTROL: MyISAM column gaining AUTO_INCREMENT backed by an EXISTING non-first-position
		// key needs no grouping (probed OK on both versions).
		{"aig-myisam-gain-ai-existing-nonleftmost", "t",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, c INT NOT NULL, KEY k (a, c)) ENGINE=MyISAM",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, c INT NOT NULL AUTO_INCREMENT, KEY k (a, c)) ENGINE=MyISAM", both()},
		// 8.0-only backing-key shapes that still satisfy errno 1075 (probed): DESC first column,
		// INVISIBLE index, and a trailing functional part behind the plain first column.
		{"aig-add-col-desc-key", "t",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY)",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, seq INT NOT NULL AUTO_INCREMENT, KEY k_seq (seq DESC))", only(MySQL80)},
		{"aig-add-col-invisible-key", "t",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY)",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, seq INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_seq (seq) INVISIBLE)", only(MySQL80)},
		{"aig-add-col-trailing-expr-key", "t",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT)",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT, seq INT NOT NULL AUTO_INCREMENT, KEY k_se (seq, (a + 1)))", only(MySQL80)},

		// ---- DROP direction: the last backing key leaves while the column is AUTO_INCREMENT ----

		// The AUTO_INCREMENT column and its backing key are dropped together: the DROP INDEX must
		// share the DROP COLUMN statement (ungrouped, the PhasePre DROP INDEX fails errno 1075).
		{"aig-drop-key-and-column", "t",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_c (c))",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY)", both()},
		// The column loses AUTO_INCREMENT and its backing key is dropped: the DROP INDEX must
		// share the de-AUTO_INCREMENT MODIFY statement.
		{"aig-drop-key-deai", "t",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_c (c))",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL)", both()},
		// The backing key is RESHAPED (same name) while the column stays AUTO_INCREMENT: the
		// DROP INDEX must share the re-ADD statement — the secondary-key analog of the combined
		// PRIMARY KEY change.
		{"aig-replace-backing-key", "t",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, UNIQUE KEY uk_c (c))",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, UNIQUE KEY uk_c (c, a))", both()},
		// The backing key is SWAPPED for a differently named one while the column stays
		// AUTO_INCREMENT.
		{"aig-swap-backing-key", "t",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_c (c))",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, KEY k_c (c))", both()},
		// Several covering keys all drop while a new one arrives: every covering drop joins the
		// covering add's statement.
		{"aig-multi-drop-plus-add", "t",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk1 (c), KEY k2 (c))",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk3 (c))", both()},
		// CONTROL: dropping ONE of several covering keys needs no grouping — another key still
		// covers the column at every statement boundary.
		{"aig-drop-covered-elsewhere", "t",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk1 (c), KEY k2 (c))",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, KEY k2 (c))", both()},
		// The PRIMARY KEY backs the AUTO_INCREMENT column and both leave: DROP PRIMARY KEY joins
		// the DROP COLUMN statement.
		{"aig-drop-pk-and-column", "t",
			"CREATE TABLE t (c INT NOT NULL AUTO_INCREMENT PRIMARY KEY, a INT NOT NULL)",
			"CREATE TABLE t (a INT NOT NULL)", both()},
		// The PRIMARY KEY leaves and the column loses AUTO_INCREMENT: DROP PRIMARY KEY joins the
		// MODIFY statement.
		{"aig-drop-pk-deai", "t",
			"CREATE TABLE t (c INT NOT NULL AUTO_INCREMENT PRIMARY KEY, a INT)",
			"CREATE TABLE t (c INT NOT NULL, a INT)", both()},
		// A COMPOSITE PK (AUTO_INCREMENT column first) leaves together with the column.
		{"aig-drop-composite-pk-and-column", "t",
			"CREATE TABLE t (c INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, PRIMARY KEY (c, a))",
			"CREATE TABLE t (a INT NOT NULL)", both()},
		// The PK moves ONTO the AUTO_INCREMENT column while its old secondary backing key drops:
		// the covering drop joins the combined DROP+ADD PRIMARY KEY statement.
		{"aig-pk-swap-drop-secondary", "t",
			"CREATE TABLE t (x INT NOT NULL, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_c (c), PRIMARY KEY (x))",
			"CREATE TABLE t (x INT NOT NULL, c INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (c))", both()},
		// The AUTO_INCREMENT column keeps its attribute but is REDEFINED (type widen) while the
		// backing key is reshaped in the same plan: the widen MODIFY stays its own statement
		// (still covered by the old key) and the reshape is grouped.
		{"aig-widen-ai-key-replaced", "t",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, UNIQUE KEY uk_c (c))",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c BIGINT NOT NULL AUTO_INCREMENT, a INT NOT NULL, UNIQUE KEY uk_c (c, a))", both()},
	}
}

// TestOracle_AutoIncGroupingApplyCorrectness proves gate 2 for every grouping form: the
// generated plan, applied statement-by-statement to a real `from` database, must yield `to`.
// Before the merge pass, the hazard probes fail on the first ungrouped statement (errno 1075).
func TestOracle_AutoIncGroupingApplyCorrectness(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		n := NormalizerFor(version)
		for _, probe := range autoIncMigrationProbes() {
			if !containsVersion(probe.versions, version) {
				continue
			}
			t.Run(fmt.Sprintf("%s/%s", o.name, probe.id), func(t *testing.T) {
				assertApplyCorrect(t, o, n, probe)
			})
		}
	}
}

// TestOracle_AutoIncGroupingIdempotence proves gate 1 for the same forms: each `to` schema,
// loaded from the engine's own readback, must self-plan to an EMPTY migration (the merge pass
// must never manufacture ops on a no-op release), and the user-form round-trip must be empty.
func TestOracle_AutoIncGroupingIdempotence(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		n := NormalizerFor(version)
		for _, probe := range autoIncMigrationProbes() {
			if !containsVersion(probe.versions, version) || strings.TrimSpace(probe.toDDL) == "" {
				continue
			}
			t.Run(fmt.Sprintf("%s/%s", o.name, probe.id), func(t *testing.T) {
				dbName := "aig_idem_" + strings.ReplaceAll(probe.id, "-", "_")
				cat, _, rb, ok := o.catalogFromReadback(t, dbName, probe.toDDL, probe.table)
				if !ok {
					t.Skipf("[%s] could not obtain readback for %s", o.name, probe.id)
				}
				diff := DiffWithNormalizer(cat, cat, n)
				plan := GenerateMigrationWithNormalizer(cat, cat, diff, n)
				if plan.SQL() != "" {
					t.Errorf("[%s] NON-EMPTY NO-OP PLAN for %s:\n  stored: %s\n  plan:\n%s",
						o.name, probe.id, strings.TrimSpace(rb), plan.SQL())
				}
				userCat, _ := loadOneTable(t, serverCharsetFor(o.version), probe.toDDL, probe.table)
				rtDiff := DiffWithNormalizer(userCat, cat, n)
				rtPlan := GenerateMigrationWithNormalizer(userCat, cat, rtDiff, n)
				if rtPlan.SQL() != "" {
					t.Errorf("[%s] NON-EMPTY ROUND-TRIP PLAN for %s:\n  user:   %s\n  stored: %s\n  plan:\n%s",
						o.name, probe.id, strings.TrimSpace(probe.toDDL), strings.TrimSpace(rb), rtPlan.SQL())
				}
			})
		}
	}
}

// TestOracle_AutoIncGroupingFKImplicitBacking proves the FK-implicit last-resort form on the
// ALTER path (the analog of formatCreateTable's last-resort inline): an AUTO_INCREMENT column
// added together with a FOREIGN KEY on it — and no user key — must have the FK's backing index
// synthesized into the grouped ADD COLUMN statement; the deferred PhasePost FK then reuses that
// index instead of auto-creating a duplicate (probed live: both versions accept the sequence).
func TestOracle_AutoIncGroupingFKImplicitBacking(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	mc := fkMultiCase{
		id: "aig-fk-implicit-backing",
		fromSQL: "CREATE TABLE p (id INT NOT NULL PRIMARY KEY); " +
			"CREATE TABLE c (id INT NOT NULL PRIMARY KEY, x INT)",
		toSQL: "CREATE TABLE p (id INT NOT NULL PRIMARY KEY); " +
			"CREATE TABLE c (id INT NOT NULL PRIMARY KEY, x INT, seq INT NOT NULL AUTO_INCREMENT, " +
			"CONSTRAINT fk_seq FOREIGN KEY (seq) REFERENCES p (id))",
		tables:   []string{"p", "c"},
		versions: both(),
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		n := NormalizerFor(version)
		t.Run(fmt.Sprintf("%s/%s", o.name, mc.id), func(t *testing.T) {
			assertFKMultiApplyCorrect(t, o, n, mc)
		})
	}
}
