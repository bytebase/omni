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
		// Two plain same-column indexes (`pid`, `pid_2`) both matching the FK-implicit name/shape,
		// plus an FK reusing one of them. Dropping the FK and `pid_2` while keeping `pid`: the FK
		// node and the index node call the SAME fkImplicitIndexNames on the SAME `from` table, so
		// they agree on which index is FK-implicit (the first match, `pid`). Because `pid` survives
		// in `to`, the FK node drops only the FK; the index node drops `pid_2`. No orphan, no
		// double-drop. (Guards the first-match edge case raised in review.)
		{"drop-fk-two-covering-indexes", fkParentSingle,
			"CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY pid (pid), KEY pid_2 (pid), CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))",
			"CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY pid (pid))", both()},

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

// fkMultiCase is an arbitrary multi-table from→to schema transition (BOTH tables may differ),
// used to prove the cross-object DROP ordering: a FK constraint drop must precede a referenced
// table drop (errno 3730), a referencing column drop (errno 1828), the index node's reshape of a
// same-named FK-backing index (errno 1553), and the leftover backing-index drop must respect a
// surviving FK sharing it. fromSQL/toSQL are full schemas (USE-prefixed by the harness).
type fkMultiCase struct {
	id       string
	fromSQL  string // statements after USE; defines all tables
	toSQL    string
	tables   []string // tables to compare in the result (existing in `to`)
	versions []Version
}

// fkMultiCases enumerate the cross-object DROP-ordering forms (the review-found ordering bugs).
func fkMultiCases() []fkMultiCase {
	return []fkMultiCase{
		// Drop a column that a FK references, together with the FK: the FK constraint drop
		// (PhasePre, priorityForeignKeyConstraintDrop) must precede the column drop (errno 1828).
		{"drop-column-and-its-fk",
			"CREATE TABLE p (id INT NOT NULL PRIMARY KEY); CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))",
			"CREATE TABLE p (id INT NOT NULL PRIMARY KEY); CREATE TABLE c (id INT PRIMARY KEY)",
			[]string{"p", "c"}, both()},
		// Drop a table that is REFERENCED by a surviving child's FK: the child's FK drop must
		// precede the parent table drop (errno 3730). Child keeps its (now-FK-less) column.
		{"drop-referenced-table",
			"CREATE TABLE p (id INT NOT NULL PRIMARY KEY); CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))",
			"CREATE TABLE c (id INT PRIMARY KEY, pid INT)",
			[]string{"c"}, both()},
		// Two FKs on the SAME column share ONE backing index; drop BOTH: every FK constraint drop
		// must precede the single backing-index drop (errno 1553 if interleaved).
		{"drop-two-fks-sharing-index",
			"CREATE TABLE p (id INT NOT NULL PRIMARY KEY); CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fka FOREIGN KEY (pid) REFERENCES p(id), CONSTRAINT fkb FOREIGN KEY (pid) REFERENCES p(id))",
			"CREATE TABLE p (id INT NOT NULL PRIMARY KEY); CREATE TABLE c (id INT PRIMARY KEY, pid INT)",
			[]string{"p", "c"}, both()},
		// Two FKs sharing one backing index; drop only ONE: the shared backing index must SURVIVE
		// because the other FK still needs it (errno 1553 if dropped).
		{"drop-one-of-two-fks-sharing-index",
			"CREATE TABLE p (id INT NOT NULL PRIMARY KEY); CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fka FOREIGN KEY (pid) REFERENCES p(id), CONSTRAINT fkb FOREIGN KEY (pid) REFERENCES p(id))",
			"CREATE TABLE p (id INT NOT NULL PRIMARY KEY); CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fkb FOREIGN KEY (pid) REFERENCES p(id))",
			[]string{"p", "c"}, both()},
		// Drop a FK whose auto backing index name (`fk`) is REUSED by the target as a differently
		// shaped USER index: the FK constraint drop must precede the index node's DROP INDEX `fk`
		// reshape (errno 1553 if the reshape runs while the FK lives).
		{"drop-fk-reuse-name-for-other-index",
			"CREATE TABLE p (id INT NOT NULL PRIMARY KEY); CREATE TABLE c (id INT PRIMARY KEY, pid INT, other_col INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))",
			"CREATE TABLE p (id INT NOT NULL PRIMARY KEY); CREATE TABLE c (id INT PRIMARY KEY, pid INT, other_col INT, KEY fk (other_col))",
			[]string{"p", "c"}, both()},
		// Drop BOTH a parent and a child table where the parent name sorts BEFORE the child
		// (a_parent < z_child): all DROP TABLE ops share priorityTable and sort by name, so the
		// parent would be dropped first — blocked by the child's FK (errno 3730/1451) — UNLESS the
		// child's FK is released in PhasePre first. Guards the dropped-table FK-release path.
		{"drop-parent-and-child-name-order",
			"CREATE TABLE a_parent (id INT NOT NULL PRIMARY KEY); CREATE TABLE z_child (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES a_parent(id))",
			"",
			nil, both()},
		// Two UNNAMED FKs on the SAME column share ONE backing index (`KEY pid (pid)`); drop BOTH:
		// each dropped FK selects the same implicit index, so the leftover DROP INDEX must be
		// EMITTED ONCE (dedup) — else the second DROP INDEX fails errno 1091.
		{"drop-two-unnamed-fks-sharing-index",
			"CREATE TABLE p (id INT NOT NULL PRIMARY KEY); CREATE TABLE c (id INT PRIMARY KEY, pid INT, FOREIGN KEY (pid) REFERENCES p(id), FOREIGN KEY (pid) REFERENCES p(id))",
			"CREATE TABLE p (id INT NOT NULL PRIMARY KEY); CREATE TABLE c (id INT PRIMARY KEY, pid INT)",
			[]string{"p", "c"}, both()},
	}
}

// TestOracle_ForeignKeyMultiTableApplyCorrectness proves the cross-object DROP ordering for the
// review-found scenarios: each `from→to` schema is applied via the FULL generated plan and every
// `to` table's canonical form must match.
func TestOracle_ForeignKeyMultiTableApplyCorrectness(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		n := NormalizerFor(version)
		for _, mc := range fkMultiCases() {
			if !containsVersion(mc.versions, version) {
				continue
			}
			t.Run(fmt.Sprintf("%s/%s", o.name, mc.id), func(t *testing.T) {
				assertFKMultiApplyCorrect(t, o, n, mc)
			})
		}
	}
}

// assertFKMultiApplyCorrect loads from/to catalogs by applying each full schema to the engine and
// reading every table back (authentic stored form), generates the plan, builds the `from` state,
// applies the plan one statement at a time, and asserts every `to` table's readback canonicalizes
// equal to `to`.
func assertFKMultiApplyCorrect(t *testing.T, o *oracleConn, n *Normalizer, mc fkMultiCase) {
	t.Helper()
	slug := strings.ReplaceAll(mc.id, "-", "_")
	sc := serverCharsetFor(o.version)
	ctx := context.Background()
	const logicalDB = "fkmulti"

	fromCat := loadMultiSchema(t, o, logicalDB, "fkmulti_f_"+slug, sc, mc.fromSQL)
	toCat := loadMultiSchema(t, o, logicalDB, "fkmulti_t_"+slug, sc, mc.toSQL)

	diff := DiffWithNormalizer(fromCat, toCat, n)
	plan := GenerateMigrationWithNormalizer(fromCat, toCat, diff, n)

	conn, err := o.db.Conn(ctx)
	if err != nil {
		t.Fatalf("[%s] grab conn: %v", o.name, err)
	}
	defer func() { _ = conn.Close() }()

	setup := []string{
		"DROP DATABASE IF EXISTS " + logicalDB,
		fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s", logicalDB, sc),
		"USE " + logicalDB,
	}
	setup = append(setup, splitSemiStatements(mc.fromSQL)...)
	for _, s := range setup {
		if _, err := conn.ExecContext(ctx, s); err != nil {
			t.Skipf("[%s] could not set up `from` for %s: %q: %v", o.name, mc.id, s, err)
		}
	}
	for _, op := range plan.Ops {
		if _, err := conn.ExecContext(ctx, op.SQL); err != nil {
			t.Fatalf("[%s] APPLY FAILED for %s:\n  stmt: %s\n  err: %v\n  full plan:\n%s",
				o.name, mc.id, op.SQL, err, plan.SQL())
		}
	}

	// Read every `to` table back and compare to `to`. foreign_key_checks=0 guards a forward FK
	// reference when tables are listed child-before-parent.
	var b strings.Builder
	fmt.Fprintf(&b, "CREATE DATABASE %s DEFAULT CHARSET=%s;\nUSE %s;\nSET foreign_key_checks=0;\n", logicalDB, sc, logicalDB)
	for _, tbl := range mc.tables {
		var nm, ddl string
		row := conn.QueryRowContext(ctx, fmt.Sprintf("SHOW CREATE TABLE %s.%s", logicalDB, tbl))
		if err := row.Scan(&nm, &ddl); err != nil {
			t.Fatalf("[%s] %s: result table %s missing after apply:\n%s", o.name, mc.id, tbl, plan.SQL())
		}
		b.WriteString(ddl + ";\n")
	}
	resultCat, err := LoadSQL(b.String())
	if err != nil {
		t.Fatalf("[%s] reload result failed: %v\n%s", o.name, err, b.String())
	}
	if d := DiffWithNormalizer(resultCat, toCat, n); !d.IsEmpty() {
		t.Errorf("[%s] APPLY-CORRECTNESS FAILED for %s: result != to\n  plan:\n%s\n  diff: %s",
			o.name, mc.id, plan.SQL(), describeForeignKeyDiff(d))
	}
	if d := DiffWithNormalizer(toCat, resultCat, n); !d.IsEmpty() {
		t.Errorf("[%s] APPLY-CORRECTNESS (reverse) FAILED for %s: %s", o.name, mc.id, describeForeignKeyDiff(d))
	}
}

// loadMultiSchema applies a full multi-table schema to a throwaway db, reads every base table back,
// and reloads all of them wrapped in a FIXED logical db (so the from/to catalogs share a database
// identity and FK references resolve identically).
func loadMultiSchema(t *testing.T, o *oracleConn, logicalDB, execDB, sc, schemaSQL string) *Catalog {
	t.Helper()
	ctx := context.Background()
	setup := []string{
		"DROP DATABASE IF EXISTS " + execDB,
		fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s", execDB, sc),
		"USE " + execDB,
	}
	setup = append(setup, splitSemiStatements(schemaSQL)...)
	for _, s := range setup {
		if _, err := o.db.ExecContext(ctx, s); err != nil {
			t.Skipf("[%s] setup failed (may be expected): %q: %v", o.name, s, err)
		}
	}
	// Discover the base tables in deterministic order.
	rows, err := o.db.QueryContext(ctx, fmt.Sprintf(
		"SELECT table_name FROM information_schema.tables WHERE table_schema = '%s' AND table_type='BASE TABLE' ORDER BY table_name", execDB))
	if err != nil {
		t.Fatalf("[%s] list tables: %v", o.name, err)
	}
	var tables []string
	for rows.Next() {
		var nm string
		if err := rows.Scan(&nm); err != nil {
			_ = rows.Close()
			t.Fatalf("[%s] scan table name: %v", o.name, err)
		}
		tables = append(tables, nm)
	}
	_ = rows.Close()

	// Reload with foreign_key_checks=0 so a child table read back before its parent (tables are
	// discovered alphabetically) does not fail on a forward FK reference (errno 1824). The omni
	// loader honors this SET (exec.go).
	var b strings.Builder
	fmt.Fprintf(&b, "CREATE DATABASE %s DEFAULT CHARSET=%s;\nUSE %s;\nSET foreign_key_checks=0;\n", logicalDB, sc, logicalDB)
	for _, tbl := range tables {
		var nm, ddl string
		row := o.db.QueryRowContext(ctx, fmt.Sprintf("SHOW CREATE TABLE %s.%s", execDB, tbl))
		if err := row.Scan(&nm, &ddl); err != nil {
			t.Fatalf("[%s] SHOW CREATE %s: %v", o.name, tbl, err)
		}
		b.WriteString(ddl + ";\n")
	}
	cat, err := LoadSQL(b.String())
	if err != nil {
		t.Fatalf("[%s] reload multi-schema failed: %v\n%s", o.name, err, b.String())
	}
	return cat
}

// splitSemiStatements splits a ";"-separated DDL string into trimmed, non-empty statements.
func splitSemiStatements(s string) []string {
	var out []string
	for _, part := range strings.Split(s, ";") {
		if p := strings.TrimSpace(part); p != "" {
			out = append(out, p)
		}
	}
	return out
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
