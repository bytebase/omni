package catalog

import (
	"context"
	"fmt"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// The oracle apply-correctness proof for the foreign-key generator (correctness-protocol.md gate
// 2), against the LIVE MySQL engines. For representative (from, to) pairs covering every FK
// variant, the generated ALTER TABLE ... ADD/DROP FOREIGN KEY DDL applied to a real `from`
// database must yield a child table whose canonical form equals `to`. This file ALSO proves the
// index↔FK interaction explicitly (the contract with the merged index node):
//   - ADD FK reusing a user index (no duplicate backing index);
//   - ADD FK auto-creating a backing index;
//   - DROP FK leaving / dropping the backing index (the FK node owns the leftover);
//   - FK column change (DROP+ADD with the backing index handled).
//
// FKs always span a parent and a child, so the harness is multi-table: it sets up parent `p` +
// child `c` in state `from`, runs the FULL plan (table/column/index DDL from the other nodes is
// inert here since only the child's FK differs, but the FK ops include the PhasePost ordering and
// backing-index drops), reads the child back, and compares. The harness reuses connectOracle /
// serverCharsetFor / both / only / containsVersion from the shared oracle tests.

// fkMigrationCase is one apply-correctness case: transform child `c` (referencing parent `p`) from
// fromChild to toChild. parentDDL is shared by both states. An empty fromChild/toChild means the
// child does not exist in that state (pure CREATE / DROP of the whole child).
type fkMigrationCase struct {
	id        string
	parentDDL string
	fromChild string
	toChild   string
	versions  []Version
}

// fkMigrationCases enumerate the FK ADD/DROP/MODIFY forms the generator covers, each proven as a
// (from, to) child transition applied to a real database.
func fkMigrationCases() []fkMigrationCase {
	return []fkMigrationCase{
		// ---- ADD FK (modified child) ----
		// ADD FK with no covering index → engine auto-creates the backing index (named `fk`).
		{"add-fk-autoindex", fkParentSingle,
			"CREATE TABLE c (id INT PRIMARY KEY, pid INT)",
			"CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))", both()},
		// ADD FK reusing an existing user index (distinct name) → NO duplicate backing index. The
		// index is already present in `from`, so the FK add in PhasePost reuses it.
		{"add-fk-reuse-user-index", fkParentSingle,
			"CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY my_idx (pid))",
			"CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY my_idx (pid), CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))", both()},
		// ADD FK with explicit actions.
		{"add-fk-cascade", fkParentSingle,
			"CREATE TABLE c (id INT PRIMARY KEY, pid INT)",
			"CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id) ON DELETE CASCADE ON UPDATE SET NULL)", both()},
		// ADD unnamed FK (auto-named constraint + backing index).
		{"add-fk-unnamed", fkParentSingle,
			"CREATE TABLE c (id INT PRIMARY KEY, pid INT)",
			"CREATE TABLE c (id INT PRIMARY KEY, pid INT, FOREIGN KEY (pid) REFERENCES p(id))", both()},
		// ADD composite FK (auto-creates a composite backing index).
		{"add-fk-composite", fkParentComposite,
			"CREATE TABLE c (id INT PRIMARY KEY, x INT, y INT)",
			"CREATE TABLE c (id INT PRIMARY KEY, x INT, y INT, CONSTRAINT fk FOREIGN KEY (x,y) REFERENCES p(a,b))", both()},

		// ---- DROP FK (modified child) ----
		// DROP FK whose backing index was auto-created → the FK node drops the leftover index too
		// (target child has neither FK nor index).
		{"drop-fk-and-backing-index", fkParentSingle,
			"CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))",
			"CREATE TABLE c (id INT PRIMARY KEY, pid INT)", both()},
		// DROP FK while a distinctly-named user index covers the columns → the FK node drops ONLY
		// the FK; the user index `my_idx` is kept (index node's domain) and must remain.
		{"drop-fk-keep-user-index", fkParentSingle,
			"CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY my_idx (pid), CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))",
			"CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY my_idx (pid))", both()},
		// DROP unnamed FK (backing index named after the column `pid`) → drop both.
		{"drop-fk-unnamed", fkParentSingle,
			"CREATE TABLE c (id INT PRIMARY KEY, pid INT, FOREIGN KEY (pid) REFERENCES p(id))",
			"CREATE TABLE c (id INT PRIMARY KEY, pid INT)", both()},

		// ---- MODIFY FK (DROP+ADD) ----
		// Action change CASCADE→SET NULL: DROP the old FK, ADD the new one (same name).
		{"modify-fk-action", fkParentSingle,
			"CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id) ON DELETE CASCADE)",
			"CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id) ON DELETE SET NULL)", both()},
		// Referenced-column change a→b: DROP+ADD. The backing index rides along; reused.
		{"modify-fk-refcol", fkParentSingle,
			"CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(a))",
			"CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(b))", both()},
		// FK referencing-column change x→y: the backing index changes columns with it (the FK node
		// owns it; the index node emits nothing). DROP+ADD.
		{"modify-fk-column", fkParentSingle,
			"CREATE TABLE c (id INT PRIMARY KEY, x INT, y INT, CONSTRAINT fk FOREIGN KEY (x) REFERENCES p(a))",
			"CREATE TABLE c (id INT PRIMARY KEY, x INT, y INT, CONSTRAINT fk FOREIGN KEY (y) REFERENCES p(b))", both()},

		// ---- self-referential + multiple ----
		// Self-referential FK (child references itself). Parent `p` exists but is unused by `c`.
		{"self-referential", fkParentSingle,
			"CREATE TABLE c (id INT PRIMARY KEY, parent_id INT)",
			"CREATE TABLE c (id INT PRIMARY KEY, parent_id INT, CONSTRAINT fk_self FOREIGN KEY (parent_id) REFERENCES c(id))", both()},
		// Add a second FK alongside an existing one (only the new FK is added; the existing one is
		// untouched).
		{"add-second-fk", fkParentSingle,
			"CREATE TABLE c (id INT PRIMARY KEY, pa INT, pb INT, CONSTRAINT fka FOREIGN KEY (pa) REFERENCES p(a))",
			"CREATE TABLE c (id INT PRIMARY KEY, pa INT, pb INT, CONSTRAINT fka FOREIGN KEY (pa) REFERENCES p(a), CONSTRAINT fkb FOREIGN KEY (pb) REFERENCES p(b))", both()},
	}
}

// TestOracle_ForeignKeyMigrationApplyCorrectness proves gate 2 for every FK case: the generated
// DDL transforms a real `from` database into a `to`-equal one (child compared via canonical
// readback), including the PhasePost ordering and backing-index handling.
func TestOracle_ForeignKeyMigrationApplyCorrectness(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		n := NormalizerFor(version)
		for _, fc := range fkMigrationCases() {
			if !containsVersion(fc.versions, version) {
				continue
			}
			t.Run(fmt.Sprintf("%s/%s", o.name, fc.id), func(t *testing.T) {
				assertFKApplyCorrect(t, o, n, fc)
			})
		}
	}
}

// assertFKApplyCorrect is the multi-table apply-correctness assertion. It loads from/to catalogs
// (parent + child) from the engine's own readbacks (authentic stored form), generates the plan,
// builds a real database in state `from`, applies the plan one statement at a time, reads the
// child back, and asserts the child canonicalizes equal to `to` (both directions).
func assertFKApplyCorrect(t *testing.T, o *oracleConn, n *Normalizer, fc fkMigrationCase) {
	t.Helper()
	slug := strings.ReplaceAll(fc.id, "-", "_")
	sc := serverCharsetFor(o.version)
	ctx := context.Background()

	// Build the from/to catalogs from engine readbacks, wrapped in a FIXED logical db so the FK
	// references resolve and the plan qualifies tables consistently with the apply db.
	const logicalDB = "fkapply"
	fromCat := loadFKSchema(t, o, logicalDB, "fkapply_f_"+slug, fc.parentDDL, fc.fromChild)
	toCat := loadFKSchema(t, o, logicalDB, "fkapply_t_"+slug, fc.parentDDL, fc.toChild)

	diff := DiffWithNormalizer(fromCat, toCat, n)
	plan := GenerateMigrationWithNormalizer(fromCat, toCat, diff, n)

	// Build a real database in state `from` on ONE dedicated connection (USE must stick).
	conn, err := o.db.Conn(ctx)
	if err != nil {
		t.Fatalf("[%s] grab conn: %v", o.name, err)
	}
	defer func() { _ = conn.Close() }()

	setup := []string{
		"DROP DATABASE IF EXISTS " + logicalDB,
		fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s", logicalDB, sc),
		"USE " + logicalDB,
		fc.parentDDL,
		fc.fromChild,
	}
	for _, s := range setup {
		if _, err := conn.ExecContext(ctx, s); err != nil {
			t.Skipf("[%s] could not set up `from` state for %s: %q: %v", o.name, fc.id, s, err)
		}
	}

	// Apply the migration statements one at a time on the same connection.
	for _, op := range plan.Ops {
		if _, err := conn.ExecContext(ctx, op.SQL); err != nil {
			t.Fatalf("[%s] APPLY FAILED for %s:\n  stmt: %s\n  err: %v\n  full plan:\n%s",
				o.name, fc.id, op.SQL, err, plan.SQL())
		}
	}

	// Read the child back and compare to `to`'s child. Wrap both in the same logical db (the diff
	// shares it) so the comparison is FK-resolution-consistent.
	var nm, resultDDL string
	row := conn.QueryRowContext(ctx, fmt.Sprintf("SHOW CREATE TABLE %s.c", logicalDB))
	if err := row.Scan(&nm, &resultDDL); err != nil {
		t.Fatalf("[%s] %s: result child missing after apply:\n%s", o.name, fc.id, plan.SQL())
	}
	// Build a result catalog that mirrors the `to` catalog shape: parent readback + result child,
	// so FK references resolve identically on both sides.
	resultCat := loadResultChildSchema(t, o, logicalDB, sc, resultDDL)

	if d := DiffWithNormalizer(resultCat, toCat, n); !d.IsEmpty() {
		t.Errorf("[%s] APPLY-CORRECTNESS FAILED for %s: result != to\n  plan:\n%s\n  result child: %s\n  diff: %s",
			o.name, fc.id, plan.SQL(), strings.TrimSpace(resultDDL), describeForeignKeyDiff(d))
	}
	if d := DiffWithNormalizer(toCat, resultCat, n); !d.IsEmpty() {
		t.Errorf("[%s] APPLY-CORRECTNESS (reverse) FAILED for %s: %s", o.name, fc.id, describeForeignKeyDiff(d))
	}
}

// loadResultChildSchema reads the parent `p` back from the live apply db and combines it with the
// already-read result child DDL into a single-catalog schema (wrapped in the logical db), so the
// result catalog has the same parent+child shape as `to` and FK references resolve on both sides.
func loadResultChildSchema(t *testing.T, o *oracleConn, logicalDB, sc, childDDL string) *Catalog {
	t.Helper()
	ctx := context.Background()
	var nm, parentDDL string
	row := o.db.QueryRowContext(ctx, fmt.Sprintf("SHOW CREATE TABLE %s.p", logicalDB))
	if err := row.Scan(&nm, &parentDDL); err != nil {
		t.Fatalf("[%s] SHOW CREATE parent after apply: %v", o.name, err)
	}
	wrapped := fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s;\nUSE %s;\n%s;\n%s;",
		logicalDB, sc, logicalDB, parentDDL, childDDL)
	cat, err := LoadSQL(wrapped)
	if err != nil {
		t.Fatalf("[%s] reload result schema failed: %v\n%s", o.name, err, wrapped)
	}
	return cat
}

// TestOracle_ForeignKeyMigrationIdempotence proves gate 1 (the spine) for the generator: for an FK
// schema in its real stored form, the generated no-op plan is EMPTY — including the FK ops. A
// non-empty no-op plan is a normalization/ordering bug.
func TestOracle_ForeignKeyMigrationIdempotence(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		n := NormalizerFor(version)
		for _, fc := range fkDiffCases() {
			if !containsVersion(fc.versions, version) {
				continue
			}
			t.Run(fmt.Sprintf("%s/%s", o.name, fc.id), func(t *testing.T) {
				slug := strings.ReplaceAll(fc.id, "-", "_")
				cat := loadFKSchema(t, o, "fkidem", "fkidem_"+slug, fc.parentDDL, fc.childDDL)
				diff := DiffWithNormalizer(cat, cat, n)
				plan := GenerateMigrationWithNormalizer(cat, cat, diff, n)
				if plan.SQL() != "" {
					t.Errorf("[%s] NON-EMPTY NO-OP PLAN for %s (FK normalization/ordering bug):\n%s",
						o.name, fc.id, plan.SQL())
				}
			})
		}
	}
}

// TestOracle_ForeignKeyAddReusesExistingIndex proves the index↔FK ADD interaction at the plan
// level: when `from` already has a user index covering the FK columns, the plan's FK ADD does NOT
// also emit an ADD INDEX (the index is reused), and applying the plan does not create a duplicate
// index. This is the generate-side guard against errno 1061 (duplicate key name).
func TestOracle_ForeignKeyAddReusesExistingIndex(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		n := NormalizerFor(version)
		t.Run(o.name+"/add-fk-no-extra-index-op", func(t *testing.T) {
			from := loadFKSchema(t, o, "fkreuse", "fkreuse_f", fkParentSingle,
				"CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY my_idx (pid))")
			to := loadFKSchema(t, o, "fkreuse", "fkreuse_t", fkParentSingle,
				"CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY my_idx (pid), CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))")
			diff := DiffWithNormalizer(from, to, n)
			plan := GenerateMigrationWithNormalizer(from, to, diff, n)
			// The plan must contain exactly one ADD FOREIGN KEY and NO ADD INDEX (the user index is
			// already present and reused).
			var addFK, addIdx int
			for _, op := range plan.Ops {
				switch op.Type {
				case OpAddForeignKey:
					addFK++
				case OpAddIndex:
					addIdx++
				}
			}
			if addFK != 1 {
				t.Errorf("[%s] expected exactly 1 ADD FOREIGN KEY, got %d:\n%s", o.name, addFK, plan.SQL())
			}
			if addIdx != 0 {
				t.Errorf("[%s] expected NO ADD INDEX (user index reused), got %d:\n%s", o.name, addIdx, plan.SQL())
			}
		})
	}
}
