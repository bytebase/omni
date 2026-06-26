package catalog

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// The oracle proof for the view differ + generator (correctness-protocol.md gates 1 & 2), against
// the LIVE MySQL engines (5.7 :13307 ssl-disabled, 8.0 :13306). Two properties are proven
// mechanically across every view FORM the node covers, on BOTH versions:
//
//  1. IDEMPOTENCE (the spine). A schema (tables + views) is applied to a real database; each
//     view's SHOW CREATE VIEW and each table's SHOW CREATE TABLE are read back and reassembled
//     into SDL. Loading that engine-stored SDL and diffing it against itself is empty, and —
//     the headline property — diffing the USER-form SDL against the ENGINE-stored SDL is also
//     empty. A non-empty diff on a no-op means the view-body canonicalization (db-qualifier
//     stripping) disagrees with what the engine stores → a normalization bug.
//
//  2. APPLY-CORRECTNESS. For representative (from, to) schema pairs, the generated CREATE OR
//     REPLACE / DROP VIEW DDL applied to a real `from` database yields a schema whose views
//     canonicalize equal to `to` — including view-on-table and view-on-view dependency ordering.
//
// The harness reuses connectOracle / serverCharsetFor / both() / only() / containsVersion from the
// normalize + diff oracle tests, and skips cleanly when the engines are unreachable (go test
// -short skips it).
//
// 5.7-vs-8.0 view-body divergences this file exercises (oracle-verified, see diff_view.go):
//   - column qualification: 8.0 qualifies same-db refs with the database, 5.7 (for view-on-view)
//     does not — folded by canonicalViewBody on both.
//   - explicit column list: kept on 8.0, rewritten into renamed SELECT aliases on 5.7. The 5.7
//     rewrite is NOT reproduced by the omni loader, so explicit-column-list views are proven on
//     8.0 and FLAGGED as a known 5.7 limitation (they are version-gated to 8.0 here).

// viewSchemaProbe is one schema (a set of CREATE statements) plus the table + view names to read
// back. `setup` are the CREATE TABLE / CREATE VIEW statements in dependency order (LoadSDL is
// order-tolerant, but the engine apply is not, so they are listed creatable-order).
type viewSchemaProbe struct {
	id       string
	setup    []string
	tables   []string
	views    []string
	versions []Version
}

// viewIdempotenceProbes enumerates representative view FORMS whose user-declared body differs from
// the engine's stored SHOW CREATE VIEW form but must canonicalize equal. Each is proven on the
// listed versions.
func viewIdempotenceProbes() []viewSchemaProbe {
	return []viewSchemaProbe{
		{"simple", []string{
			"CREATE TABLE t (a INT, b VARCHAR(20), c INT)",
			"CREATE VIEW v AS SELECT a, b FROM t",
		}, []string{"t"}, []string{"v"}, both()},
		{"star-expanded", []string{
			"CREATE TABLE t (a INT, b VARCHAR(20), c INT)",
			"CREATE VIEW v AS SELECT * FROM t",
		}, []string{"t"}, []string{"v"}, both()},
		{"where-expr", []string{
			"CREATE TABLE t (a INT, b VARCHAR(20))",
			"CREATE VIEW v AS SELECT a FROM t WHERE a > 0",
		}, []string{"t"}, []string{"v"}, both()},
		{"expr-and-func", []string{
			"CREATE TABLE t (a INT, b VARCHAR(20))",
			"CREATE VIEW v AS SELECT a + 1 AS s, UPPER(b) AS ub FROM t",
		}, []string{"t"}, []string{"v"}, both()},
		{"join", []string{
			"CREATE TABLE t (a INT, b VARCHAR(20), c INT)",
			"CREATE TABLE u (x INT, y INT)",
			"CREATE VIEW v AS SELECT t.a, u.y FROM t JOIN u ON t.c = u.x",
		}, []string{"t", "u"}, []string{"v"}, both()},
		{"left-join", []string{
			"CREATE TABLE t (a INT, c INT)",
			"CREATE TABLE u (x INT, y INT)",
			"CREATE VIEW v AS SELECT t.a, u.y FROM t LEFT JOIN u ON t.c = u.x",
		}, []string{"t", "u"}, []string{"v"}, both()},
		{"algorithm-merge", []string{
			"CREATE TABLE t (a INT)",
			"CREATE ALGORITHM=MERGE VIEW v AS SELECT a FROM t",
		}, []string{"t"}, []string{"v"}, both()},
		{"algorithm-temptable", []string{
			"CREATE TABLE t (a INT)",
			"CREATE ALGORITHM=TEMPTABLE VIEW v AS SELECT a FROM t",
		}, []string{"t"}, []string{"v"}, both()},
		{"sql-security-invoker", []string{
			"CREATE TABLE t (a INT)",
			"CREATE SQL SECURITY INVOKER VIEW v AS SELECT a FROM t",
		}, []string{"t"}, []string{"v"}, both()},
		{"check-cascaded", []string{
			"CREATE TABLE t (a INT)",
			"CREATE VIEW v AS SELECT a FROM t WHERE a > 0 WITH CASCADED CHECK OPTION",
		}, []string{"t"}, []string{"v"}, both()},
		{"check-local", []string{
			"CREATE TABLE t (a INT)",
			"CREATE VIEW v AS SELECT a FROM t WHERE a > 0 WITH LOCAL CHECK OPTION",
		}, []string{"t"}, []string{"v"}, both()},
		{"view-on-view", []string{
			"CREATE TABLE t (a INT, b INT)",
			"CREATE VIEW base AS SELECT a, b FROM t",
			"CREATE VIEW v AS SELECT a FROM base",
		}, []string{"t"}, []string{"base", "v"}, both()},
		{"cte", []string{
			"CREATE TABLE t (a INT, b INT)",
			"CREATE VIEW v AS WITH cte AS (SELECT a, b FROM t) SELECT a FROM cte",
		}, []string{"t"}, []string{"v"}, only(MySQL80)},
		{"union", []string{
			"CREATE TABLE t (a INT)",
			"CREATE TABLE u (a INT)",
			"CREATE VIEW v AS SELECT a FROM t UNION SELECT a FROM u",
		}, []string{"t", "u"}, []string{"v"}, both()},
		{"group-by-agg", []string{
			"CREATE TABLE t (a INT, b INT)",
			"CREATE VIEW v AS SELECT a, COUNT(*) AS n FROM t GROUP BY a",
		}, []string{"t"}, []string{"v"}, both()},
		// Explicit column list: 8.0 keeps the list + original aliases (round-trips); 5.7 rewrites
		// into renamed aliases (a FLAGGED loader limitation), so proven on 8.0 only.
		{"explicit-columns", []string{
			"CREATE TABLE t (a INT, b VARCHAR(20))",
			"CREATE VIEW v (p, q) AS SELECT a, b FROM t",
		}, []string{"t"}, []string{"v"}, only(MySQL80)},
		{"view-on-view-explicit-cols", []string{
			"CREATE TABLE t (a INT, b INT)",
			"CREATE VIEW base (x, y) AS SELECT a, b FROM t",
			"CREATE VIEW v AS SELECT x FROM base",
		}, []string{"t"}, []string{"base", "v"}, only(MySQL80)},
	}
}

// applyViewSchema drops + recreates dbName with the box's server charset and applies every setup
// statement on a single dedicated connection (database/sql pools connections, so USE must share
// the connection with the following statements). Returns the connection (caller closes) or skips.
func (o *oracleConn) applyViewSchema(t *testing.T, dbName string, setup []string) (*sql.Conn, bool) {
	t.Helper()
	ctx := context.Background()
	conn, err := o.db.Conn(ctx)
	if err != nil {
		t.Fatalf("[%s] grab conn: %v", o.name, err)
	}
	stmts := append([]string{
		"DROP DATABASE IF EXISTS " + dbName,
		fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s", dbName, serverCharsetFor(o.version)),
		"USE " + dbName,
	}, setup...)
	for _, s := range stmts {
		if _, err := conn.ExecContext(ctx, s); err != nil {
			t.Logf("[%s] view schema setup failed (may be expected): %q: %v", o.name, s, err)
			_ = conn.Close()
			return nil, false
		}
	}
	return conn, true
}

// dumpViewSchemaSDL reads back SHOW CREATE for the named tables and views on conn (queried in the
// physical database `physDB`) and assembles a reloadable SDL string wrapped in the LOGICAL database
// `logicalDB`. SHOW CREATE VIEW returns a db-qualified view name (`physDB`.`v`); to make the
// reloaded catalog's identity independent of the throwaway physical db, that qualifier is rewritten
// to `logicalDB`, so both the `from` and `to` catalogs (and the apply database) share one database
// name and the generated DDL qualifies objects under it. Views are dumped in the given order
// (dependency order); LoadSDL re-sorts anyway.
func (o *oracleConn) dumpViewSchemaSDL(t *testing.T, conn *sql.Conn, physDB, logicalDB string, tables, views []string) (string, bool) {
	t.Helper()
	ctx := context.Background()
	var b strings.Builder
	fmt.Fprintf(&b, "CREATE DATABASE %s DEFAULT CHARSET=%s;\nUSE %s;\n", logicalDB, serverCharsetFor(o.version), logicalDB)
	for _, tbl := range tables {
		var name, ddl string
		row := conn.QueryRowContext(ctx, fmt.Sprintf("SHOW CREATE TABLE %s.%s", physDB, tbl))
		if err := row.Scan(&name, &ddl); err != nil {
			t.Logf("[%s] SHOW CREATE TABLE %s failed: %v", o.name, tbl, err)
			return "", false
		}
		// SHOW CREATE TABLE returns an unqualified `CREATE TABLE \`t\` ...`; the USE above scopes it.
		b.WriteString(ddl)
		b.WriteString(";\n")
	}
	physPrefix := "`" + physDB + "`."
	logicalPrefix := "`" + logicalDB + "`."
	for _, vw := range views {
		// SHOW CREATE VIEW returns 4 columns: View, Create View, character_set_client, collation.
		var name, ddl, csClient, coll string
		row := conn.QueryRowContext(ctx, fmt.Sprintf("SHOW CREATE VIEW %s.%s", physDB, vw))
		if err := row.Scan(&name, &ddl, &csClient, &coll); err != nil {
			t.Logf("[%s] SHOW CREATE VIEW %s failed: %v", o.name, vw, err)
			return "", false
		}
		// Rewrite the physical-db qualifier (on the view name and any same-db references in the
		// body) to the logical db so the reloaded catalog is db-name-stable.
		b.WriteString(strings.ReplaceAll(ddl, physPrefix, logicalPrefix))
		b.WriteString(";\n")
	}
	return b.String(), true
}

// userFormSDL wraps a probe's setup statements into a reloadable SDL string under dbName.
func userFormSDL(o *oracleConn, dbName string, setup []string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "CREATE DATABASE %s DEFAULT CHARSET=%s;\nUSE %s;\n", dbName, serverCharsetFor(o.version), dbName)
	for _, s := range setup {
		b.WriteString(s)
		b.WriteString(";\n")
	}
	return b.String()
}

// describeViewDiff renders a compact description of a SchemaDiff's view changes for failure output.
func describeViewDiff(d *SchemaDiff) string {
	var b strings.Builder
	for _, ve := range d.Views {
		fmt.Fprintf(&b, "[view %s.%s %s]", ve.Database, ve.Name, ve.Action)
	}
	for _, te := range d.Tables {
		fmt.Fprintf(&b, "[table %s %s]", te.Name, te.Action)
	}
	return b.String()
}

// TestOracle_ViewIdempotence proves gates 1 & 2 for every view form: the user-declared schema and
// its engine-stored readback diff EMPTY (both directions), and each side self-diffs empty.
func TestOracle_ViewIdempotence(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		n := NormalizerFor(version)
		for _, probe := range viewIdempotenceProbes() {
			if !containsVersion(probe.versions, version) {
				continue
			}
			t.Run(fmt.Sprintf("%s/%s", o.name, probe.id), func(t *testing.T) {
				dbName := "videm_" + strings.ReplaceAll(probe.id, "-", "_")
				conn, ok := o.applyViewSchema(t, dbName, probe.setup)
				if !ok {
					t.Skipf("[%s] could not apply schema for %s", o.name, probe.id)
				}
				defer func() { _ = conn.Close() }()

				storedSDL, ok := o.dumpViewSchemaSDL(t, conn, dbName, dbName, probe.tables, probe.views)
				if !ok {
					t.Skipf("[%s] could not dump readback for %s", o.name, probe.id)
				}
				userSDL := userFormSDL(o, dbName, probe.setup)

				storedCat, err := LoadSDLWithVersion(storedSDL, version)
				if err != nil {
					t.Fatalf("[%s] load stored SDL failed: %v\n%s", o.name, err, storedSDL)
				}
				userCat, err := LoadSDLWithVersion(userSDL, version)
				if err != nil {
					t.Fatalf("[%s] load user SDL failed: %v\n%s", o.name, err, userSDL)
				}

				if d := DiffWithNormalizer(storedCat, storedCat, n); !d.IsEmpty() {
					t.Errorf("[%s] IDEMPOTENCE: stored self-diff not empty for %s: %s\n  stored SDL:\n%s",
						o.name, probe.id, describeViewDiff(d), storedSDL)
				}
				if d := DiffWithNormalizer(userCat, userCat, n); !d.IsEmpty() {
					t.Errorf("[%s] user self-diff not empty for %s: %s", o.name, probe.id, describeViewDiff(d))
				}
				if d := DiffWithNormalizer(userCat, storedCat, n); !d.IsEmpty() {
					t.Errorf("[%s] CANONICALIZATION: user vs stored not empty for %s: %s\n  user SDL:\n%s\n  stored SDL:\n%s",
						o.name, probe.id, describeViewDiff(d), userSDL, storedSDL)
				}
				if d := DiffWithNormalizer(storedCat, userCat, n); !d.IsEmpty() {
					t.Errorf("[%s] CANONICALIZATION (reverse): stored vs user not empty for %s: %s",
						o.name, probe.id, describeViewDiff(d))
				}
			})
		}
	}
}

// viewMigrationProbe is one apply-correctness case: transform a database from the `from` schema to
// the `to` schema. Empty `from`/`to` slices mean "no objects" on that side.
type viewMigrationProbe struct {
	id       string
	from     []string
	to       []string
	tables   []string // table names present in `to` (for readback assembly)
	views    []string // view names present in `to`
	versions []Version
}

// viewMigrationProbes enumerates the CREATE / REPLACE / DROP view FORMS the generator covers,
// including view-on-table and view-on-view dependency ordering.
func viewMigrationProbes() []viewMigrationProbe {
	base := func(ss ...string) []string { return ss }
	return []viewMigrationProbe{
		// ---- CREATE a view (from a schema that has only the table) ----
		{"create-simple",
			base("CREATE TABLE t (a INT, b INT)"),
			base("CREATE TABLE t (a INT, b INT)", "CREATE VIEW v AS SELECT a FROM t"),
			[]string{"t"}, []string{"v"}, both()},
		{"create-with-options",
			base("CREATE TABLE t (a INT, b INT)"),
			base("CREATE TABLE t (a INT, b INT)", "CREATE ALGORITHM=MERGE SQL SECURITY INVOKER VIEW v AS SELECT a FROM t WHERE a > 0 WITH CASCADED CHECK OPTION"),
			[]string{"t"}, []string{"v"}, both()},
		// view-on-table created together with its table in one plan (table first).
		{"create-view-and-table",
			base(),
			base("CREATE TABLE t (a INT, b INT)", "CREATE VIEW v AS SELECT a FROM t"),
			[]string{"t"}, []string{"v"}, both()},
		// ---- view-on-view created in one plan: dependency ordering ----
		{"create-view-on-view",
			base("CREATE TABLE t (a INT, b INT)"),
			base("CREATE TABLE t (a INT, b INT)", "CREATE VIEW base AS SELECT a, b FROM t", "CREATE VIEW v AS SELECT a FROM base"),
			[]string{"t"}, []string{"base", "v"}, both()},
		{"create-view-chain",
			base("CREATE TABLE t (a INT)"),
			base("CREATE TABLE t (a INT)", "CREATE VIEW v1 AS SELECT a FROM t", "CREATE VIEW v2 AS SELECT a FROM v1", "CREATE VIEW v3 AS SELECT a FROM v2"),
			[]string{"t"}, []string{"v1", "v2", "v3"}, both()},
		// Regression (review blocker): base view `a` carries a column alias named `b`, and view `b`
		// references `a` (selecting the column `b` that `a` exposes). The alias must not create a
		// false `a->b` edge that reverses CREATE order (which would fail to apply, since `a` would be
		// created after `b`). Only the real edge b->a exists, so `a` must be created first.
		{"create-view-alias-name-clash",
			base("CREATE TABLE t (id INT)"),
			base("CREATE TABLE t (id INT)", "CREATE VIEW a AS SELECT id AS b FROM t", "CREATE VIEW b AS SELECT b FROM a"),
			[]string{"t"}, []string{"a", "b"}, both()},
		// ---- REPLACE (modify body / options) ----
		{"replace-body",
			base("CREATE TABLE t (a INT, b INT)", "CREATE VIEW v AS SELECT a FROM t"),
			base("CREATE TABLE t (a INT, b INT)", "CREATE VIEW v AS SELECT b FROM t"),
			[]string{"t"}, []string{"v"}, both()},
		{"replace-algorithm",
			base("CREATE TABLE t (a INT)", "CREATE ALGORITHM=UNDEFINED VIEW v AS SELECT a FROM t"),
			base("CREATE TABLE t (a INT)", "CREATE ALGORITHM=MERGE VIEW v AS SELECT a FROM t"),
			[]string{"t"}, []string{"v"}, both()},
		{"replace-sql-security",
			base("CREATE TABLE t (a INT)", "CREATE SQL SECURITY DEFINER VIEW v AS SELECT a FROM t"),
			base("CREATE TABLE t (a INT)", "CREATE SQL SECURITY INVOKER VIEW v AS SELECT a FROM t"),
			[]string{"t"}, []string{"v"}, both()},
		{"replace-add-check-option",
			base("CREATE TABLE t (a INT)", "CREATE VIEW v AS SELECT a FROM t WHERE a > 0"),
			base("CREATE TABLE t (a INT)", "CREATE VIEW v AS SELECT a FROM t WHERE a > 0 WITH CASCADED CHECK OPTION"),
			[]string{"t"}, []string{"v"}, both()},
		// ---- DROP ----
		{"drop-view",
			base("CREATE TABLE t (a INT)", "CREATE VIEW v AS SELECT a FROM t"),
			base("CREATE TABLE t (a INT)"),
			[]string{"t"}, nil, both()},
		// Drop a dependent view and its base in one plan (both dropped; order irrelevant for DROP).
		{"drop-view-on-view",
			base("CREATE TABLE t (a INT)", "CREATE VIEW base AS SELECT a FROM t", "CREATE VIEW v AS SELECT a FROM base"),
			base("CREATE TABLE t (a INT)"),
			[]string{"t"}, nil, both()},
		// ---- explicit column list (8.0 only — 5.7 loader limitation flagged) ----
		{"create-explicit-columns",
			base("CREATE TABLE t (a INT, b VARCHAR(20))"),
			base("CREATE TABLE t (a INT, b VARCHAR(20))", "CREATE VIEW v (p, q) AS SELECT a, b FROM t"),
			[]string{"t"}, []string{"v"}, only(MySQL80)},
		// Regression (re-review): a `AS TABLE <view>` body (TABLE query primary, 8.0.19+) deparses to
		// `table \`base\``, which re-parses to a *TableStmt; the dependency extractor must still see
		// the referenced view so `tv` is created after `base`. 8.0 only (TABLE primary is 8.0.19+).
		{"create-view-on-view-table-primary",
			base("CREATE TABLE t (id INT)"),
			base("CREATE TABLE t (id INT)", "CREATE VIEW base AS SELECT id FROM t", "CREATE VIEW tv AS TABLE base"),
			[]string{"t"}, []string{"base", "tv"}, only(MySQL80)},
	}
}

// TestOracle_ViewMigrationApplyCorrectness proves gate 2 for every view migration probe: the
// generated DDL transforms a real `from` database into a `to`-equal one (compared via canonical
// readback of the full schema).
func TestOracle_ViewMigrationApplyCorrectness(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		n := NormalizerFor(version)
		for _, probe := range viewMigrationProbes() {
			if !containsVersion(probe.versions, version) {
				continue
			}
			t.Run(fmt.Sprintf("%s/%s", o.name, probe.id), func(t *testing.T) {
				assertViewApplyCorrect(t, o, n, probe)
			})
		}
	}
}

// assertViewApplyCorrect loads from/to catalogs from the engine's own readbacks, generates the
// plan, applies it to a real `from`-state database, reads the result back, and asserts the result
// canonicalizes equal to `to`.
func assertViewApplyCorrect(t *testing.T, o *oracleConn, n *Normalizer, p viewMigrationProbe) {
	t.Helper()
	slug := strings.ReplaceAll(p.id, "-", "_")

	// Both catalogs load under ONE logical database (`vt`, the apply database) so the generated
	// DDL qualifies objects with `vt` and applies against the real `vt` database below.
	const applyDB = "vt"

	// Build the `to` catalog from the engine's readback of the `to` schema, re-homed to `vt`.
	toCat := loadSchemaFromEngine(t, o, "vgen_to_"+slug, applyDB, p.to, p.tables, p.views)
	if toCat == nil {
		t.Skipf("[%s] could not obtain `to` readback for %s", o.name, p.id)
	}
	// Build the `from` catalog likewise (its objects determine the readback set).
	fromTables, fromViews := schemaObjectNames(p.from)
	fromCat := loadSchemaFromEngine(t, o, "vgen_from_"+slug, applyDB, p.from, fromTables, fromViews)
	if fromCat == nil {
		t.Skipf("[%s] could not obtain `from` readback for %s", o.name, p.id)
	}

	diff := DiffWithNormalizer(fromCat, toCat, n)
	plan := GenerateMigrationWithNormalizer(fromCat, toCat, diff, n)

	// Build a real database in state `from` (wrapped in `vt`, the database both catalogs load
	// under via dumpViewSchemaSDL/userFormSDL → matches the plan's `vt`.`obj` qualification).
	ctx := context.Background()
	conn, err := o.db.Conn(ctx)
	if err != nil {
		t.Fatalf("[%s] grab conn: %v", o.name, err)
	}
	defer func() { _ = conn.Close() }()

	setup := []string{
		"DROP DATABASE IF EXISTS " + applyDB,
		fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s", applyDB, serverCharsetFor(o.version)),
		"USE " + applyDB,
	}
	setup = append(setup, p.from...)
	for _, s := range setup {
		if _, err := conn.ExecContext(ctx, s); err != nil {
			t.Skipf("[%s] could not set up `from` state for %s: %q: %v", o.name, p.id, s, err)
		}
	}

	// Apply the migration one statement at a time on the same connection.
	for _, op := range plan.Ops {
		if _, err := conn.ExecContext(ctx, op.SQL); err != nil {
			t.Fatalf("[%s] APPLY FAILED for %s:\n  stmt: %s\n  err: %v\n  full plan:\n%s",
				o.name, p.id, op.SQL, err, plan.SQL())
		}
	}

	// Read back the resulting schema and compare to `to`.
	resultSDL, ok := o.dumpViewSchemaSDL(t, conn, applyDB, applyDB, p.tables, p.views)
	if !ok {
		t.Fatalf("[%s] %s: could not read back result", o.name, p.id)
	}
	resultCat, err := LoadSDLWithVersion(resultSDL, o.version)
	if err != nil {
		t.Fatalf("[%s] %s: reload of result failed: %v\n%s", o.name, p.id, err, resultSDL)
	}

	if d := DiffWithNormalizer(resultCat, toCat, n); !d.IsEmpty() {
		t.Errorf("[%s] APPLY-CORRECTNESS FAILED for %s: result != to\n  plan:\n%s\n  result SDL:\n%s\n  diff: %s",
			o.name, p.id, plan.SQL(), resultSDL, describeViewDiff(d))
	}
	if d := DiffWithNormalizer(toCat, resultCat, n); !d.IsEmpty() {
		t.Errorf("[%s] APPLY-CORRECTNESS (reverse) FAILED for %s: %s", o.name, p.id, describeViewDiff(d))
	}

	// Also assert the dropped views are actually gone (DROP correctness): any view in `from` but
	// not in `to` must not survive.
	_, fromV := schemaObjectNames(p.from)
	toSet := make(map[string]bool)
	for _, v := range p.views {
		toSet[strings.ToLower(v)] = true
	}
	for _, v := range fromV {
		if toSet[strings.ToLower(v)] {
			continue
		}
		var n2, ddl, cs, co string
		row := conn.QueryRowContext(ctx, fmt.Sprintf("SHOW CREATE VIEW %s.%s", applyDB, v))
		if err := row.Scan(&n2, &ddl, &cs, &co); err == nil {
			t.Errorf("[%s] %s: dropped view %s still exists after plan:\n%s", o.name, p.id, v, plan.SQL())
		}
	}
}

// loadSchemaFromEngine applies a schema to a throwaway PHYSICAL db, reads back the named tables +
// views, and loads them into a catalog re-homed under the LOGICAL db name (version-aware). Returns
// nil (and the caller skips) on any apply/dump failure — an input a version rejects is not a
// generator failure. An empty setup yields an empty catalog (the `from` of a pure CREATE).
func loadSchemaFromEngine(t *testing.T, o *oracleConn, physDB, logicalDB string, setup, tables, views []string) *Catalog {
	t.Helper()
	if len(setup) == 0 {
		return New()
	}
	conn, ok := o.applyViewSchema(t, physDB, setup)
	if !ok {
		return nil
	}
	defer func() { _ = conn.Close() }()
	sdl, ok := o.dumpViewSchemaSDL(t, conn, physDB, logicalDB, tables, views)
	if !ok {
		return nil
	}
	cat, err := LoadSDLWithVersion(sdl, o.version)
	if err != nil {
		t.Fatalf("[%s] load schema-from-engine failed: %v\n%s", o.name, err, sdl)
	}
	return cat
}

// schemaObjectNames extracts the table and view names declared by a list of CREATE statements, so
// the `from` side can be read back without a hand-maintained name list. It recognises
// "CREATE [OR REPLACE] [ALGORITHM=..] [DEFINER=..] [SQL SECURITY ..] VIEW <name>" and
// "CREATE TABLE <name>" (the only DDL the probes use), taking the unqualified name.
func schemaObjectNames(setup []string) (tables, views []string) {
	for _, s := range setup {
		up := strings.ToUpper(s)
		if idx := strings.Index(up, " VIEW "); idx >= 0 && strings.HasPrefix(strings.TrimSpace(up), "CREATE") {
			views = append(views, firstIdentAfter(s, idx+len(" VIEW ")))
			continue
		}
		if idx := strings.Index(up, "CREATE TABLE "); idx == 0 {
			tables = append(tables, firstIdentAfter(s, idx+len("CREATE TABLE ")))
		}
	}
	return tables, views
}

// firstIdentAfter returns the first identifier starting at or after pos in s, stripping optional
// backticks, a schema qualifier, and stopping at the first space / '(' / '`'.
func firstIdentAfter(s string, pos int) string {
	for pos < len(s) && (s[pos] == ' ' || s[pos] == '\t') {
		pos++
	}
	start := pos
	for pos < len(s) {
		c := s[pos]
		if c == ' ' || c == '(' || c == '\t' || c == '\n' {
			break
		}
		pos++
	}
	ident := s[start:pos]
	ident = strings.ReplaceAll(ident, "`", "")
	if dot := strings.LastIndex(ident, "."); dot >= 0 {
		ident = ident[dot+1:]
	}
	return ident
}
