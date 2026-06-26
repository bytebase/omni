package catalog

import (
	"context"
	"fmt"
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
		{"modify-charset", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, v VARCHAR(10) CHARACTER SET latin1) DEFAULT CHARSET=utf8mb4",
			"CREATE TABLE t (id INT PRIMARY KEY, v VARCHAR(10) CHARACTER SET utf8mb4) DEFAULT CHARSET=utf8mb4", both()},
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

	// Hyphens are illegal in unquoted MySQL identifiers; slug the probe id for db names.
	slug := strings.ReplaceAll(p.id, "-", "_")

	// Build the `to` catalog from the engine's readback (authentic stored form).
	var toCat *Catalog
	tableInTo := strings.TrimSpace(p.toDDL) != ""
	if tableInTo {
		c, _, _, ok := o.catalogFromReadback(t, "gen_to_"+slug, p.toDDL, p.table)
		if !ok {
			t.Skipf("[%s] could not obtain `to` readback for %s", o.name, p.id)
		}
		toCat = c
	} else {
		toCat = New()
	}

	// Build the `from` catalog likewise.
	var fromCat *Catalog
	tableInFrom := strings.TrimSpace(p.fromDDL) != ""
	if tableInFrom {
		c, _, _, ok := o.catalogFromReadback(t, "gen_from_"+slug, p.fromDDL, p.table)
		if !ok {
			t.Skipf("[%s] could not obtain `from` readback for %s", o.name, p.id)
		}
		fromCat = c
	} else {
		fromCat = New()
	}

	// Generate the plan (version-aware).
	diff := DiffWithNormalizer(fromCat, toCat, n)
	plan := GenerateMigrationWithNormalizer(fromCat, toCat, diff, n)

	// Build a real database in state `from` and apply the plan, all on ONE dedicated
	// connection (database/sql pools connections, so a stray `USE` on the pool can land on a
	// different connection than the next statement). The catalogs were loaded via loadOneTable,
	// which wraps DDL in database `diffdb`, so the plan qualifies tables as `diffdb`.`t` — we
	// must set up and apply in that same `diffdb`.
	ctx := context.Background()
	conn, err := o.db.Conn(ctx)
	if err != nil {
		t.Fatalf("[%s] grab conn: %v", o.name, err)
	}
	defer func() { _ = conn.Close() }()

	const applyDB = "diffdb"
	setup := []string{
		"DROP DATABASE IF EXISTS " + applyDB,
		fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s", applyDB, serverCharsetFor(o.version)),
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
	resultCat, _ := loadOneTable(t, serverCharsetFor(o.version), resultRB, p.table)

	if d := DiffWithNormalizer(resultCat, toCat, n); !d.IsEmpty() {
		toRB, _ := o.readbackTable(t, "gen_to_"+slug, p.table)
		t.Errorf("[%s] APPLY-CORRECTNESS FAILED for %s: result != to\n  plan:\n%s\n  result: %s\n  want:   %s\n  diff: %s",
			o.name, p.id, plan.SQL(), strings.TrimSpace(resultRB), strings.TrimSpace(toRB), describeDiff(d))
	}
	// Symmetry: to vs result must also be empty.
	if d := DiffWithNormalizer(toCat, resultCat, n); !d.IsEmpty() {
		t.Errorf("[%s] APPLY-CORRECTNESS (reverse) FAILED for %s: %s", o.name, p.id, describeDiff(d))
	}
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
