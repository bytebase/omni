package catalog

import (
	"context"
	"fmt"
	"hash/fnv"
	"os"
	"strconv"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// The oracle proof for generate-core (correctness-protocol.md gates 1 & 2), against the
// LIVE MySQL engines (5.7 and 8.0). This is the first node that EMITS DDL, so it proves both
// oracle gates mechanically:
//
//  1. APPLY-CORRECTNESS. For a representative (from, to) pair: build a real database in
//     state `from` (apply from's DDL), run GenerateMigration(from, to).SQL() against it, then
//     SHOW CREATE the result and assert its canonical form equals `to`'s canonical form
//     (DiffWithNormalizer(result, to).IsEmpty()). The from/to catalogs are loaded from the
//     ENGINE'S OWN readbacks so they are in authentic stored form, exactly as the release
//     path sees them.
//
//  2. IDEMPOTENCE (the spine). When from == to the generated plan is EMPTY (SQL() == ""), and
//     the full round-trip — apply `to`, dump SHOW CREATE, reload, diff against `to` → empty,
//     generate → empty — produces no DDL. A non-empty no-op plan is a bug in ordering or
//     rendering (or an upstream normalization gap → flagged, not patched here).
//
// The harness reuses connectOracle / showCreate / serverCharsetFor / loadOneTable / both()
// from the diff + normalize oracle tests. It skips cleanly when the engines are unreachable,
// so the unit suite stays hermetic (go test -short skips it).

// readbackTable returns the SHOW CREATE TABLE for a table in an existing database.
func (o *oracleConn) readbackTable(t *testing.T, dbName, table string) (string, bool) {
	t.Helper()
	var name, ddl string
	row := o.db.QueryRowContext(context.Background(), fmt.Sprintf("SHOW CREATE TABLE %s.%s", dbName, table))
	if err := row.Scan(&name, &ddl); err != nil {
		t.Logf("[%s] SHOW CREATE %s.%s failed: %v", o.name, dbName, table, err)
		return "", false
	}
	return ddl, true
}

// catalogFromReadback applies createSQL in a throwaway db, reads back the table's SHOW
// CREATE, and loads it into a catalog (wrapped with the box's server charset so table-charset
// inheritance resolves identically to loadOneTable). It returns the catalog + the table's
// authentic stored form.
func (o *oracleConn) catalogFromReadback(t *testing.T, dbName, createSQL, table string) (*Catalog, *Table, string, bool) {
	t.Helper()
	rb, ok := o.showCreate(t, dbName, createSQL, table)
	if !ok {
		return nil, nil, "", false
	}
	cat, tbl := loadOneTable(t, serverCharsetFor(o.version), rb, table)
	return cat, tbl, rb, true
}

// migrationProbe is one apply-correctness case: transform a single table from `fromDDL` to
// `toDDL`. An empty fromDDL means "table does not exist yet" (a pure CREATE).
type migrationProbe struct {
	id       string
	table    string
	fromDDL  string // "" → table absent in from
	toDDL    string
	versions []Version
}

// migrationProbes enumerates the table + column op FORMS this node covers. Each is proven on
// the listed versions: the generated DDL, applied to a real `from` database, must yield a
// table whose canonical form equals `to`.
func migrationProbes() []migrationProbe {
	return []migrationProbe{
		// ---- CREATE TABLE (from empty) — every column rendering surface ----
		{"create-simple", "t", "",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20))", both()},
		{"create-int-widths", "t", "",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT(11), b BIGINT(20), p TINYINT, h INT UNSIGNED, z INT ZEROFILL)", both()},
		{"create-boolean", "t", "",
			"CREATE TABLE t (id INT PRIMARY KEY, m BOOLEAN, d TINYINT(1))", both()},
		{"create-numeric", "t", "",
			"CREATE TABLE t (id INT PRIMARY KEY, a DECIMAL, b DECIMAL(10,2), j REAL, f FLOAT, g DOUBLE(15,4))", both()},
		{"create-char-bit-year", "t", "",
			"CREATE TABLE t (id INT PRIMARY KEY, a CHAR, d BINARY, c BIT, y YEAR, v VARCHAR(40), bb VARBINARY(8))", both()},
		{"create-enum-set", "t", "",
			`CREATE TABLE t (id INT PRIMARY KEY, a ENUM('x','y','z'), b SET('a','b','c'))`, both()},
		{"create-defaults", "t", "",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT DEFAULT 0, b INT DEFAULT '7', f DECIMAL(10,2) DEFAULT 0, c VARCHAR(20) DEFAULT 'x', j BOOLEAN DEFAULT TRUE)", both()},
		{"create-nullability", "t", "",
			"CREATE TABLE t (id INT, a INT NULL, b INT DEFAULT NULL, d VARCHAR(10) NOT NULL, PRIMARY KEY (id))", both()},
		{"create-comment", "t", "",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT COMMENT 'has ''quote'' inside') COMMENT='tbl ''c'''", both()},
		{"create-autoinc", "t", "",
			"CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY, a INT) AUTO_INCREMENT=100", both()},
		{"create-generated", "t", "",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT GENERATED ALWAYS AS (a+1) VIRTUAL, c INT GENERATED ALWAYS AS (a*2) STORED)", both()},
		{"create-timestamp", "t", "",
			"CREATE TABLE t (id INT PRIMARY KEY, a TIMESTAMP, b TIMESTAMP NULL, d TIMESTAMP DEFAULT CURRENT_TIMESTAMP)", both()},
		{"create-datetime-default", "t", "",
			"CREATE TABLE t (id INT PRIMARY KEY, g DATETIME DEFAULT CURRENT_TIMESTAMP, h DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP)", both()},
		{"create-charset-explicit", "t", "",
			"CREATE TABLE t (id INT PRIMARY KEY, a VARCHAR(10) CHARACTER SET utf8mb4, b VARCHAR(10)) DEFAULT CHARSET=utf8mb4", both()},
		{"create-engine-myisam", "t", "",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT) ENGINE=MyISAM", both()},
		{"create-multi-column", "t", "",
			"CREATE TABLE t (id INT NOT NULL, c1 VARCHAR(10), c2 INT DEFAULT 0, c3 DATE, c4 DECIMAL(8,3), PRIMARY KEY (id))", both()},
		{"create-functional-default", "t", "",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT DEFAULT (1+1), b VARCHAR(36) DEFAULT (UUID()))", only(MySQL80)},
		{"create-row-format", "t", "",
			"CREATE TABLE t (id INT PRIMARY KEY) ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8", both()},

		// ---- DROP TABLE ----
		{"drop-table", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT)",
			"", both()},

		// ---- ADD COLUMN ----
		{"add-column", "t",
			"CREATE TABLE t (id INT PRIMARY KEY)",
			"CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(20))", both()},
		{"add-column-notnull-default", "t",
			"CREATE TABLE t (id INT PRIMARY KEY)",
			"CREATE TABLE t (id INT PRIMARY KEY, n INT NOT NULL DEFAULT 5)", both()},
		{"add-generated-column", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT)",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT GENERATED ALWAYS AS (a+1) STORED)", both()},
		{"add-timestamp-column", "t",
			"CREATE TABLE t (id INT PRIMARY KEY)",
			"CREATE TABLE t (id INT PRIMARY KEY, ts TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP)", both()},

		// ---- DROP COLUMN ----
		{"drop-column", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(20))",
			"CREATE TABLE t (id INT PRIMARY KEY)", both()},

		// ---- MODIFY COLUMN (type / nullability / default / charset / comment / autoinc) ----
		{"modify-type", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, v INT)",
			"CREATE TABLE t (id INT PRIMARY KEY, v BIGINT)", both()},
		{"modify-type-width-57-vs-80", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, v SMALLINT)",
			"CREATE TABLE t (id INT PRIMARY KEY, v BIGINT)", both()},
		{"modify-nullability", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, v INT NULL)",
			"CREATE TABLE t (id INT PRIMARY KEY, v INT NOT NULL)", both()},
		{"modify-default", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, v INT DEFAULT 0)",
			"CREATE TABLE t (id INT PRIMARY KEY, v INT DEFAULT 1)", both()},
		// bug: changing TO the BIGINT UNSIGNED max default must generate a MODIFY that applies
		// to the EXACT value 18446744073709551615 (not int64-clamped 9223372036854775807). The
		// apply-correctness diff catches a corrupted value; TestOracle_BigintUnsignedExactValue
		// additionally asserts the literal survives byte-for-byte.
		{"modify-default-bigint-unsigned-max", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, big BIGINT UNSIGNED NOT NULL DEFAULT 0)",
			"CREATE TABLE t (id INT PRIMARY KEY, big BIGINT UNSIGNED NOT NULL DEFAULT 18446744073709551615)", both()},
		// bug: changing TO a BIT default written as hex (0x05) must apply as b'101'.
		{"modify-default-bit-hex", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, f BIT(8) NOT NULL DEFAULT b'0')",
			"CREATE TABLE t (id INT PRIMARY KEY, f BIT(8) NOT NULL DEFAULT 0x05)", both()},
		// bug: changing TO a numeric YEAR default (2000) must apply and round-trip vs '2000'.
		{"modify-default-year", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, y YEAR NOT NULL DEFAULT 1999)",
			"CREATE TABLE t (id INT PRIMARY KEY, y YEAR NOT NULL DEFAULT 2000)", both()},
		{"modify-charset", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, v VARCHAR(10) CHARACTER SET latin1) DEFAULT CHARSET=utf8mb4",
			"CREATE TABLE t (id INT PRIMARY KEY, v VARCHAR(10) CHARACTER SET utf8mb4) DEFAULT CHARSET=utf8mb4", both()},
		// bug (Roundcube et al.): the legacy `BINARY` string-column modifier resolves to the
		// charset's `_bin` collation. The generated DDL for a `varchar BINARY` target must
		// apply to a real engine and yield a column whose stored form (COLLATE <cs>_bin) equals
		// the target — i.e. ADD/MODIFY emits the resolved _bin collation, not the table default.
		// (entry char-binary-attribute)
		{"add-column-binary-modifier", "t",
			"CREATE TABLE t (id INT PRIMARY KEY) DEFAULT CHARSET=utf8mb4",
			"CREATE TABLE t (id INT PRIMARY KEY, k VARCHAR(128) BINARY NOT NULL) DEFAULT CHARSET=utf8mb4", both()},
		// MODIFY a plain column TO its BINARY (_bin) form — a genuine collation change that must
		// be emitted and applied (proves the fix doesn't make plain↔BINARY an invisible no-op).
		{"modify-column-to-binary-modifier", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, v VARCHAR(10)) DEFAULT CHARSET=utf8mb4",
			"CREATE TABLE t (id INT PRIMARY KEY, v VARCHAR(10) BINARY) DEFAULT CHARSET=utf8mb4", both()},
		{"modify-comment", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, v INT COMMENT 'a')",
			"CREATE TABLE t (id INT PRIMARY KEY, v INT COMMENT 'b')", both()},
		{"modify-add-autoinc", "t",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY)",
			"CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY)", both()},
		{"modify-generated-expr", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, g INT GENERATED ALWAYS AS (a+1) STORED)",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, g INT GENERATED ALWAYS AS (a+2) STORED)", both()},

		// ---- table-option ALTER (engine / charset / comment / row_format) ----
		{"alter-engine", "t",
			"CREATE TABLE t (id INT PRIMARY KEY) ENGINE=InnoDB",
			"CREATE TABLE t (id INT PRIMARY KEY) ENGINE=MyISAM", both()},
		{"alter-comment", "t",
			"CREATE TABLE t (id INT PRIMARY KEY) COMMENT='old'",
			"CREATE TABLE t (id INT PRIMARY KEY) COMMENT='new'", both()},
		{"alter-charset", "t",
			"CREATE TABLE t (id INT PRIMARY KEY) DEFAULT CHARSET=latin1",
			"CREATE TABLE t (id INT PRIMARY KEY) DEFAULT CHARSET=utf8mb4", both()},
		{"alter-rowformat-add", "t",
			"CREATE TABLE t (id INT PRIMARY KEY)",
			"CREATE TABLE t (id INT PRIMARY KEY) ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8", both()},
		{"alter-rowformat-change", "t",
			"CREATE TABLE t (id INT PRIMARY KEY) ROW_FORMAT=COMPACT",
			"CREATE TABLE t (id INT PRIMARY KEY) ROW_FORMAT=DYNAMIC", both()},

		// ---- review-regression apply proofs (blockers #1, #2) ----
		// Generated storage-mode flip VIRTUAL→STORED: must be DROP+ADD (in-place MODIFY is
		// rejected by MySQL, error 3106). 8.0 only — stored generated cols need 5.7.x patches the
		// box may lack; proven on 8.0.
		{"gen-storage-flip", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, g INT GENERATED ALWAYS AS (a+1) VIRTUAL)",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, g INT GENERATED ALWAYS AS (a+1) STORED)", only(MySQL80)},
		// Dropping a base column AND a generated column that references it: the generated column
		// must be dropped first (error 3108 otherwise).
		{"gen-drop-with-base", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, z INT GENERATED ALWAYS AS (a+1) STORED)",
			"CREATE TABLE t (id INT PRIMARY KEY)", both()},
		// Re-round N2: a generated column `g` (refs base `z`, sorts before it) flips storage
		// (DROP+ADD) in the same plan as a type-widening MODIFY of `z`. The generated ADD must
		// land after the base MODIFY (else the re-add references a column mid-change). 8.0 only
		// (stored generated columns).
		{"gen-flip-with-base-modify", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, g INT GENERATED ALWAYS AS (z+1) VIRTUAL, z SMALLINT)",
			"CREATE TABLE t (id INT PRIMARY KEY, g INT GENERATED ALWAYS AS (z+1) STORED, z BIGINT)", only(MySQL80)},

		// ---- coverage: INVISIBLE column, COLLATE-only divergence (8.0 features) ----
		{"create-invisible-column", "t", "",
			"CREATE TABLE t (id INT PRIMARY KEY, secret INT INVISIBLE)", only(MySQL80)},
		{"add-invisible-column", "t",
			"CREATE TABLE t (id INT PRIMARY KEY)",
			"CREATE TABLE t (id INT PRIMARY KEY, secret INT INVISIBLE)", only(MySQL80)},
		{"create-collate-only-column", "t", "",
			"CREATE TABLE t (id INT PRIMARY KEY, a VARCHAR(10) COLLATE utf8mb4_unicode_ci) DEFAULT CHARSET=utf8mb4", both()},
	}
}

// TestOracle_MigrationApplyCorrectness proves gate 2 (apply-correctness) for every probe:
// the generated DDL transforms a real `from` database into a `to`-equal one.
func TestOracle_MigrationApplyCorrectness(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		n := NormalizerFor(version)
		for _, probe := range migrationProbes() {
			if !containsVersion(probe.versions, version) {
				continue
			}
			t.Run(fmt.Sprintf("%s/%s", o.name, probe.id), func(t *testing.T) {
				assertApplyCorrect(t, o, n, probe)
			})
		}
	}
}

// assertApplyCorrect is the core gate-2 assertion. It loads from/to catalogs from the
// engine's own readbacks, generates the plan, applies it to a from-state database, reads the
// result back, and asserts the result canonicalizes equal to `to` (in both directions).
func assertApplyCorrect(t *testing.T, o *oracleConn, n *Normalizer, p migrationProbe) {
	t.Helper()

	// Each probe gets its OWN apply database (and its own throwaway readback databases) so the
	// suite is hermetic: the shared apply harness used to contend on one fixed `diffdb` across the
	// generate-core, index, and check apply suites — a pre-existing pooled-`USE` race that surfaced
	// as `Unknown database 'diffdb'` when several apply suites ran in one process (or two
	// concurrent `go test` processes). applyDBName derives a per-probe-unique, identifier-safe name
	// (mirrors the partition harness's pa_<slug> convention; see assertPartitionApplyCorrect).
	applyDB := applyDBName(t, p.id)
	sc := serverCharsetFor(o.version)
	ctx := context.Background()

	// One pinned connection for ALL of this probe's statements (readbacks + setup + apply), so the
	// pool never lands a `USE` on a different connection than the next statement.
	conn, err := o.db.Conn(ctx)
	if err != nil {
		t.Fatalf("[%s] grab conn: %v", o.name, err)
	}
	defer func() { _ = conn.Close() }()

	// readbackVia applies a CREATE in a throwaway db on this connection and returns its SHOW
	// CREATE (the engine's authentic stored form). The throwaway db name is derived per-probe so
	// it never collides with another suite's readback db.
	readbackVia := func(throwaway, createSQL, table string) (string, bool) {
		for _, s := range []string{
			"DROP DATABASE IF EXISTS " + throwaway,
			fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s", throwaway, sc),
			"USE " + throwaway,
			createSQL,
		} {
			if _, err := conn.ExecContext(ctx, s); err != nil {
				t.Logf("[%s] %s readback setup failed (may be expected): %q: %v", o.name, p.id, s, err)
				return "", false
			}
		}
		var name, ddl string
		row := conn.QueryRowContext(ctx, fmt.Sprintf("SHOW CREATE TABLE %s.%s", throwaway, table))
		if err := row.Scan(&name, &ddl); err != nil {
			t.Logf("[%s] %s SHOW CREATE failed: %v", o.name, p.id, err)
			return "", false
		}
		return ddl, true
	}

	// Build `to` and `from` catalogs from the engine's own readbacks (authentic stored form),
	// each wrapped in applyDB so the generated plan qualifies tables as `applyDB`.`t` — the same
	// database we set up and apply in below.
	var toCat *Catalog
	tableInTo := strings.TrimSpace(p.toDDL) != ""
	if tableInTo {
		rb, ok := readbackVia(applyDB+"_to", p.toDDL, p.table)
		if !ok {
			t.Skipf("[%s] could not obtain `to` readback for %s", o.name, p.id)
		}
		toCat, _ = loadOnePartTable(t, applyDB, sc, rb, p.table)
	} else {
		toCat = New()
	}

	var fromCat *Catalog
	tableInFrom := strings.TrimSpace(p.fromDDL) != ""
	if tableInFrom {
		rb, ok := readbackVia(applyDB+"_from", p.fromDDL, p.table)
		if !ok {
			t.Skipf("[%s] could not obtain `from` readback for %s", o.name, p.id)
		}
		fromCat, _ = loadOnePartTable(t, applyDB, sc, rb, p.table)
	} else {
		fromCat = New()
	}

	// Generate the plan (version-aware).
	diff := DiffWithNormalizer(fromCat, toCat, n)
	plan := GenerateMigrationWithNormalizer(fromCat, toCat, diff, n)

	// Build a real database in state `from` and apply the plan, on the same pinned connection.
	setup := []string{
		"DROP DATABASE IF EXISTS " + applyDB,
		fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s", applyDB, sc),
		"USE " + applyDB,
	}
	if tableInFrom {
		setup = append(setup, p.fromDDL)
	}
	for _, s := range setup {
		if _, err := conn.ExecContext(ctx, s); err != nil {
			t.Skipf("[%s] could not set up `from` state for %s: %q: %v", o.name, p.id, s, err)
		}
	}

	// Apply the migration statements one at a time on the same connection.
	for _, op := range plan.Ops {
		if _, err := conn.ExecContext(ctx, op.SQL); err != nil {
			t.Fatalf("[%s] APPLY FAILED for %s:\n  stmt: %s\n  err: %v\n  full plan:\n%s",
				o.name, p.id, op.SQL, err, plan.SQL())
		}
	}

	// Read back the result (on the same connection) and compare to `to`.
	readback := func(table string) (string, bool) {
		var name, ddl string
		row := conn.QueryRowContext(ctx, fmt.Sprintf("SHOW CREATE TABLE %s.%s", applyDB, table))
		if err := row.Scan(&name, &ddl); err != nil {
			return "", false
		}
		return ddl, true
	}

	if !tableInTo {
		// The table should be gone (DROP). Assert SHOW CREATE now fails.
		if _, ok := readback(p.table); ok {
			t.Errorf("[%s] %s: table %s still exists after DROP plan:\n%s", o.name, p.id, p.table, plan.SQL())
		}
		return
	}

	resultRB, ok := readback(p.table)
	if !ok {
		t.Fatalf("[%s] %s: result table %s missing after apply:\n%s", o.name, p.id, p.table, plan.SQL())
	}
	// Load the result into the SAME applyDB so it shares the (database, table) identity with toCat
	// — loading into a different db name would phantom a whole-table DROP+ADD in the diff.
	resultCat, _ := loadOnePartTable(t, applyDB, sc, resultRB, p.table)

	if d := DiffWithNormalizer(resultCat, toCat, n); !d.IsEmpty() {
		toRB, _ := o.readbackTable(t, applyDB+"_to", p.table)
		t.Errorf("[%s] APPLY-CORRECTNESS FAILED for %s: result != to\n  plan:\n%s\n  result: %s\n  want:   %s\n  diff: %s",
			o.name, p.id, plan.SQL(), strings.TrimSpace(resultRB), strings.TrimSpace(toRB), describeDiff(d))
	}
	// Symmetry: to vs result must also be empty.
	if d := DiffWithNormalizer(toCat, resultCat, n); !d.IsEmpty() {
		t.Errorf("[%s] APPLY-CORRECTNESS (reverse) FAILED for %s: %s", o.name, p.id, describeDiff(d))
	}
}

// applyDBName builds a per-probe, identifier-safe apply-database name for the shared apply harness
// so each probe (and each suite that reuses assertApplyCorrect) is isolated instead of contending
// on one fixed `diffdb`. It combines the slugged probe id with a short hash of (process pid +
// full test name): the test name (t.Name()) carries the suite path
// (e.g. TestOracle_IndexMigrationApplyCorrectness/8.0/...) so two probes sharing an id across suites
// get distinct databases within a process, and the pid salts the hash so two CONCURRENT `go test`
// processes targeting the same engine never compute the same DROP/CREATE DATABASE name (the
// remaining cross-process race the fixed `diffdb` had). The result stays well under MySQL's 64-char
// identifier limit even with the harness's `_to` / `_from` throwaway suffixes appended. The `pa_`
// prefix mirrors the partition harness's convention (assertPartitionApplyCorrect).
func applyDBName(t *testing.T, probeID string) string {
	t.Helper()
	slug := strings.ReplaceAll(probeID, "-", "_")
	if len(slug) > 40 {
		slug = slug[:40]
	}
	h := fnv.New32a()
	_, _ = fmt.Fprintf(h, "%d:%s", os.Getpid(), t.Name())
	return "pa_" + slug + "_" + strconv.FormatUint(uint64(h.Sum32()), 16)
}

// TestOracle_MigrationIdempotence proves gate 1 (the spine): for a schema in its real stored
// form, the generated no-op plan is EMPTY — and the full round-trip (apply → dump → reload →
// diff → generate) emits nothing. A non-empty no-op plan is a bug.
func TestOracle_MigrationIdempotence(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		n := NormalizerFor(version)
		// Reuse the diff idempotence corpus (the same high-risk forms), but assert the PLAN is
		// empty rather than just the diff.
		for _, probe := range diffIdempotenceProbes() {
			if !containsVersion(probe.versions, version) {
				continue
			}
			t.Run(fmt.Sprintf("%s/%s", o.name, probe.id), func(t *testing.T) {
				// Load the engine's stored form of the schema.
				cat, _, rb, ok := o.catalogFromReadback(t, "gen_idem_"+strings.ReplaceAll(probe.id, "-", "_"), probe.create, probe.table)
				if !ok {
					t.Skipf("[%s] could not obtain readback for %s", o.name, probe.id)
				}
				// Self-plan must be empty (idempotence spine).
				diff := DiffWithNormalizer(cat, cat, n)
				plan := GenerateMigrationWithNormalizer(cat, cat, diff, n)
				if plan.SQL() != "" {
					t.Errorf("[%s] NON-EMPTY NO-OP PLAN for %s (normalization/ordering bug):\n  stored: %s\n  plan:\n%s",
						o.name, probe.id, strings.TrimSpace(rb), plan.SQL())
				}

				// Round-trip: apply the user form, dump, reload, diff against the user-form
				// catalog, generate — must be empty. This catches a rendering form that the
				// engine would re-normalize differently from the loader.
				userCat, _ := loadOneTable(t, serverCharsetFor(o.version), probe.create, probe.table)
				rtDiff := DiffWithNormalizer(userCat, cat, n)
				rtPlan := GenerateMigrationWithNormalizer(userCat, cat, rtDiff, n)
				if rtPlan.SQL() != "" {
					t.Errorf("[%s] NON-EMPTY ROUND-TRIP PLAN (user vs stored) for %s:\n  user:   %s\n  stored: %s\n  plan:\n%s",
						o.name, probe.id, strings.TrimSpace(probe.create), strings.TrimSpace(rb), rtPlan.SQL())
				}
			})
		}
	}
}

// TestOracle_BigintUnsignedExactValue is the dedicated bug-1 proof: a generated migration to
// a BIGINT UNSIGNED default of the uint64 maximum (18446744073709551615) must, when applied to
// a real database, store the EXACT value — not the int64-clamped 9223372036854775807. It
// asserts the literal survives byte-for-byte in the readback (the apply-correctness diff alone
// could be satisfied by both sides being equally corrupted; this checks the absolute value).
// It also covers the other integer boundaries that overflow int64 on the way in.
func TestOracle_BigintUnsignedExactValue(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	cases := []struct {
		id      string
		colType string
		def     string // exact decimal literal the user writes (and that must be stored)
	}{
		{"uint64-max", "BIGINT UNSIGNED", "18446744073709551615"},
		{"int64-max-plus-1", "BIGINT UNSIGNED", "9223372036854775808"},
		{"int64-max", "BIGINT", "9223372036854775807"},
		{"int64-min", "BIGINT", "-9223372036854775808"},
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		n := NormalizerFor(version)
		sc := serverCharsetFor(o.version)
		ctx := context.Background()
		for _, c := range cases {
			t.Run(o.name+"/"+c.id, func(t *testing.T) {
				applyDB := applyDBName(t, "bigexact_"+c.id)
				fromDDL := fmt.Sprintf("CREATE TABLE t (id INT PRIMARY KEY, big %s NOT NULL DEFAULT 0)", c.colType)
				toDDL := fmt.Sprintf("CREATE TABLE t (id INT PRIMARY KEY, big %s NOT NULL DEFAULT %s)", c.colType, c.def)

				// One pinned connection for all statements so a pooled USE never lands on a
				// different connection than the next statement.
				conn, err := o.db.Conn(ctx)
				if err != nil {
					t.Fatalf("[%s] grab conn: %v", o.name, err)
				}
				defer func() { _ = conn.Close() }()

				readbackVia := func(throwaway, createSQL string) (string, bool) {
					for _, s := range []string{
						"DROP DATABASE IF EXISTS " + throwaway,
						fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s", throwaway, sc),
						"USE " + throwaway,
						createSQL,
					} {
						if _, err := conn.ExecContext(ctx, s); err != nil {
							t.Logf("[%s] readback setup failed: %q: %v", o.name, s, err)
							return "", false
						}
					}
					var name, ddl string
					row := conn.QueryRowContext(ctx, fmt.Sprintf("SHOW CREATE TABLE %s.t", throwaway))
					if err := row.Scan(&name, &ddl); err != nil {
						return "", false
					}
					return ddl, true
				}

				// Load from/to into applyDB-named catalogs (authentic stored form) so the plan
				// qualifies tables as `applyDB`.`t` — the db we set up and apply in below.
				fromRB, ok := readbackVia(applyDB+"_from", fromDDL)
				if !ok {
					t.Skipf("[%s] could not obtain `from` readback", o.name)
				}
				toRB, ok := readbackVia(applyDB+"_to", toDDL)
				if !ok {
					t.Skipf("[%s] could not obtain `to` readback", o.name)
				}
				fromCat, _ := loadOnePartTable(t, applyDB, sc, fromRB, "t")
				toCat, _ := loadOnePartTable(t, applyDB, sc, toRB, "t")

				diff := DiffWithNormalizer(fromCat, toCat, n)
				plan := GenerateMigrationWithNormalizer(fromCat, toCat, diff, n)
				if plan.SQL() == "" {
					t.Fatalf("[%s] expected a MODIFY plan changing the default to %s, got empty", o.name, c.def)
				}
				// The generated DDL itself must carry the exact literal, never the clamped value.
				if !strings.Contains(plan.SQL(), c.def) {
					t.Errorf("[%s] generated plan must carry exact default %s:\n%s", o.name, c.def, plan.SQL())
				}
				if c.def != "9223372036854775807" && strings.Contains(plan.SQL(), "9223372036854775807") {
					t.Errorf("[%s] generated plan leaked int64-clamped value:\n%s", o.name, plan.SQL())
				}

				// Build a real from-state database and apply the plan on the pinned connection.
				for _, s := range []string{
					"DROP DATABASE IF EXISTS " + applyDB,
					fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s", applyDB, sc),
					"USE " + applyDB,
					fromDDL,
				} {
					if _, err := conn.ExecContext(ctx, s); err != nil {
						t.Skipf("[%s] could not set up from-state: %q: %v", o.name, s, err)
					}
				}
				for _, op := range plan.Ops {
					if _, err := conn.ExecContext(ctx, op.SQL); err != nil {
						t.Fatalf("[%s] APPLY FAILED: %s: %v", o.name, op.SQL, err)
					}
				}
				// COLUMN_DEFAULT is NOT NULL for these columns; COALESCE guards against a driver
				// surfacing it as NULL so a plain string scan is safe.
				var storedDefault string
				row := conn.QueryRowContext(ctx,
					"SELECT COALESCE(COLUMN_DEFAULT, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=? AND TABLE_NAME='t' AND COLUMN_NAME='big'", applyDB)
				if err := row.Scan(&storedDefault); err != nil {
					t.Fatalf("[%s] reading stored default failed: %v", o.name, err)
				}
				if storedDefault != c.def {
					t.Errorf("[%s] EXACT-VALUE CORRUPTION: stored default = %q, want %q (plan:\n%s)",
						o.name, storedDefault, c.def, plan.SQL())
				}
			})
		}
	}
}

// fkDropTableEntry builds a DiffModify TableDiffEntry whose `from` table has a single foreign key
// (named fkName, on column col) backed by a plain auto-created index of the SAME name, and whose
// `to` table keeps neither — so dropForeignKeyOps emits both the DROP FOREIGN KEY and the
// leftover-backing-index DROP INDEX. It is the minimal fixture that drives leftoverBackingIndexToDrop
// down its positive branch.
func fkDropTableEntry(dbName, table, fkName, col string) TableDiffEntry {
	db := &Database{Name: dbName}
	fk := &Constraint{Name: fkName, Type: ConForeignKey, Columns: []string{col}, RefTable: "ref", RefColumns: []string{"id"}}
	backing := &Index{Name: fkName, Columns: []*IndexColumn{{Name: col}}, Visible: true}
	from := &Table{Name: table, Database: db, Constraints: []*Constraint{fk}, Indexes: []*Index{backing}}
	to := &Table{Name: table, Database: db} // FK gone, no surviving index of that name
	return TableDiffEntry{
		Action:      DiffModify,
		Database:    dbName,
		Name:        table,
		From:        from,
		To:          to,
		ForeignKeys: []ForeignKeyDiffEntry{{Action: DiffDrop, Name: fkName, From: fk}},
	}
}

// TestForeignKeyBackingDropDedupDottedIdentifiers is the regression for the FK dedup-key collision.
// The leftover-backing-index DROP is deduplicated per (database, table, index). The old key joined
// those three with '.', so two DISTINCT triples whose names contain a literal '.' could collapse to
// the same string and wrongly suppress a legitimate second DROP INDEX:
//
//	(db "d", table "a",   index "b.c")  -> "d.a.b.c"
//	(db "d", table "a.b", index "c")    -> "d.a.b.c"   // collision!
//
// Both tables legitimately need their backing index dropped, so a correct plan emits TWO DROP INDEX
// ops. encodeKeyFields length-prefixes each field, so the two triples encode distinctly and both
// drops survive. This test fails (only one DROP INDEX) under the dotted key and passes under the
// collision-free key. It is a pure unit test (no engine) and ordinary-identifier behavior is
// unaffected.
func TestForeignKeyBackingDropDedupDottedIdentifiers(t *testing.T) {
	diff := &SchemaDiff{Tables: []TableDiffEntry{
		fkDropTableEntry("d", "a", "b.c", "x"),
		fkDropTableEntry("d", "a.b", "c", "y"),
	}}

	ops := generateForeignKeyDDL(nil, nil, diff, NormalizerFor(MySQL80))

	var dropIndex []string
	for _, op := range ops {
		if strings.Contains(op.SQL, "DROP INDEX") {
			dropIndex = append(dropIndex, op.SQL)
		}
	}
	if len(dropIndex) != 2 {
		t.Fatalf("dotted-identifier dedup collision: want 2 DROP INDEX ops (one per table), got %d:\n%s",
			len(dropIndex), strings.Join(dropIndex, "\n"))
	}
	// Both tables' leftover indexes must be present — the collision would drop exactly one of these.
	joined := strings.Join(dropIndex, "\n")
	if !strings.Contains(joined, "`d`.`a`") || !strings.Contains(joined, "`d`.`a.b`") {
		t.Errorf("expected a DROP INDEX for both `d`.`a` and `d`.`a.b`, got:\n%s", joined)
	}
}
