package catalog

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// The oracle apply-correctness proof for CROSS-OBJECT operation ordering (correctness-protocol.md
// gate 2), against the LIVE MySQL engines (5.7 :13307 ssl-disabled, 8.0 :13306). Per-object-kind
// apply-correctness lives with each kind's own oracle file; THIS file proves the plan-level
// ordering rules that span kinds — currently view ↔ stored-function:
//
//   - CREATE direction: a view whose body calls a stored function created in the SAME plan must be
//     created AFTER the function. MySQL eagerly validates function references at CREATE VIEW time —
//     Error 1305 "FUNCTION does not exist" when the function is missing (live-verified on 8.0.32
//     and 5.7.25) — so the pre-fix plan (view first) fails to apply. Routine bodies are lazily
//     validated (a CREATE FUNCTION referencing a missing function/view/table is accepted on both
//     versions), so hoisting every routine create above every view create is unconditionally safe:
//     routine CREATE/ALTER runs at priorityRoutine < priorityView (migration.go).
//   - DROP direction (inverse): when a view and a function it calls are BOTH dropped, the view is
//     dropped first (view DROP at priorityView < routine DROP at priorityRoutineDrop in PhasePre).
//     MySQL tolerates dropping a function a view still references (live-verified — the view merely
//     goes invalid), so this is the safe semantic order rather than an apply-success requirement.
//   - Transitive: function ← view ← view composes the routine-before-view priority with the view
//     node's own view-on-view depth ordering.
//
// Each probe builds from/to catalogs from the ENGINE'S OWN readbacks, generates the plan, applies
// it statement-by-statement to a real from-state database (the pre-fix failure point), reads the
// result back, and asserts (a) the result canonicalizes equal to `to` in both directions, and
// (b) CONVERGENCE — re-diffing the applied result against `to` generates an EMPTY follow-up plan.
// Scratch databases use the ordsdl_ prefix and are dropped after each probe.

// orderingProbe is one cross-object apply-correctness case: transform a database from the `from`
// schema to the `to` schema. tables/views/functions name the objects present in `to` (for readback
// assembly); the `from` side's object names are extracted from its CREATE statements.
type orderingProbe struct {
	id        string
	from      []string
	to        []string
	tables    []string
	views     []string
	functions []string
	versions  []Version
}

// orderingProbes enumerates the cross-object ordering FORMS: function+view created together (the
// enterprise-SDL-smoke repro), both dropped together, the transitive function←view←view chain, a
// function created together with a view its body SELECTs from (the lazy-direction teeth), and a
// body-modified function (DROP+CREATE path) combined with a new dependent view. Functions declare
// READS SQL DATA so they create on the 8.0 box (log_bin_trust_function_creators off → error 1418
// for a function with no binlog-safe characteristic).
func orderingProbes() []orderingProbe {
	base := func(ss ...string) []string { return ss }
	return []orderingProbe{
		// The work-order repro: one release CREATEs a function AND a view calling it. Pre-fix the
		// plan put the view first → Error 1305 on apply.
		{"create-function-and-view",
			base("CREATE TABLE users (id INT PRIMARY KEY)"),
			base("CREATE TABLE users (id INT PRIMARY KEY)",
				"CREATE FUNCTION ent_ucount() RETURNS INT READS SQL DATA RETURN (SELECT COUNT(*) FROM users)",
				"CREATE VIEW ent_v_fn AS SELECT ent_ucount() AS n"),
			[]string{"users"}, []string{"ent_v_fn"}, []string{"ent_ucount"}, both()},
		// Inverse: drop both. View drops first (PhasePre priorityView < priorityRoutineDrop).
		{"drop-function-and-view",
			base("CREATE TABLE users (id INT PRIMARY KEY)",
				"CREATE FUNCTION ent_ucount() RETURNS INT READS SQL DATA RETURN (SELECT COUNT(*) FROM users)",
				"CREATE VIEW ent_v_fn AS SELECT ent_ucount() AS n"),
			base("CREATE TABLE users (id INT PRIMARY KEY)"),
			[]string{"users"}, nil, nil, both()},
		// Transitive: v2 → v1 → fn_base. The plan must order fn_base < v1 < v2.
		{"create-function-view-view-chain",
			base("CREATE TABLE t (a INT)"),
			base("CREATE TABLE t (a INT)",
				"CREATE FUNCTION fn_base() RETURNS INT READS SQL DATA RETURN (SELECT COUNT(*) FROM t)",
				"CREATE VIEW v1 AS SELECT fn_base() AS n",
				"CREATE VIEW v2 AS SELECT n FROM v1"),
			[]string{"t"}, []string{"v1", "v2"}, []string{"fn_base"}, both()},
		// The lazy-direction teeth: the function's body SELECTs from a view created in the SAME
		// plan. The global hoist deliberately creates the function FIRST (before the view exists),
		// so this apply succeeds only because routine bodies are lazily validated — the
		// load-bearing safety claim behind hoisting routines above views. An engine that eagerly
		// validated routine bodies would fail this probe.
		{"create-function-selecting-view",
			base("CREATE TABLE t (a INT)"),
			base("CREATE TABLE t (a INT)",
				"CREATE VIEW v_src AS SELECT a FROM t",
				"CREATE FUNCTION fn_over_view() RETURNS INT READS SQL DATA RETURN (SELECT COUNT(*) FROM v_src)"),
			[]string{"t"}, []string{"v_src"}, []string{"fn_over_view"}, both()},
		// A body change forces the routine DROP+CREATE path; a NEW view calling the function lands
		// in the same plan. Order must be DROP FUNCTION (PhasePre) → CREATE FUNCTION → CREATE VIEW.
		{"modify-function-add-view",
			base("CREATE TABLE t (a INT)",
				"CREATE FUNCTION fn_base() RETURNS INT READS SQL DATA RETURN (SELECT COUNT(*) FROM t)"),
			base("CREATE TABLE t (a INT)",
				"CREATE FUNCTION fn_base() RETURNS INT READS SQL DATA RETURN (SELECT COUNT(*) + 1 FROM t)",
				"CREATE VIEW v1 AS SELECT fn_base() AS n"),
			[]string{"t"}, []string{"v1"}, []string{"fn_base"}, both()},
	}
}

// TestOracle_CrossObjectOrderingApplyCorrectness proves gate 2 for every ordering probe on both
// engines: the generated plan applies cleanly in emitted order and converges.
func TestOracle_CrossObjectOrderingApplyCorrectness(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		n := NormalizerFor(version)
		for _, probe := range orderingProbes() {
			if !containsVersion(probe.versions, version) {
				continue
			}
			t.Run(fmt.Sprintf("%s/%s", o.name, probe.id), func(t *testing.T) {
				assertOrderingApplyCorrect(t, o, n, probe)
			})
		}
	}
}

// assertOrderingApplyCorrect is the per-probe gate-2 driver (mirrors assertViewApplyCorrect, plus
// functions in the readback set and an explicit convergence check).
func assertOrderingApplyCorrect(t *testing.T, o *oracleConn, n *Normalizer, p orderingProbe) {
	t.Helper()
	slug := strings.ReplaceAll(p.id, "-", "_")
	ctx := context.Background()

	// All three scratch databases share the ordsdl_ prefix and are dropped when the probe ends.
	applyDB := "ordsdl_t"
	toDB := "ordsdl_to_" + slug
	fromDB := "ordsdl_from_" + slug
	t.Cleanup(func() {
		for _, db := range []string{applyDB, toDB, fromDB} {
			_, _ = o.db.ExecContext(context.Background(), "DROP DATABASE IF EXISTS "+db)
		}
	})

	// Build both catalogs from the engine's own readbacks, re-homed under the apply database so
	// the generated DDL qualifies objects with it.
	fromTables, fromViews, fromFunctions := orderingObjectNames(p.from)
	toCat := loadOrderingSchemaFromEngine(t, o, toDB, applyDB, p.to, p.tables, p.views, p.functions)
	if toCat == nil {
		t.Skipf("[%s] could not obtain `to` readback for %s", o.name, p.id)
	}
	fromCat := loadOrderingSchemaFromEngine(t, o, fromDB, applyDB, p.from, fromTables, fromViews, fromFunctions)
	if fromCat == nil {
		t.Skipf("[%s] could not obtain `from` readback for %s", o.name, p.id)
	}

	diff := DiffWithNormalizer(fromCat, toCat, n)
	plan := GenerateMigrationWithNormalizer(fromCat, toCat, diff, n)

	// Build a real database in `from` state and apply the plan statement-by-statement, in emitted
	// order, on one connection. A mis-ordered plan fails HERE (Error 1305 for the CREATE probes).
	conn, err := o.db.Conn(ctx)
	if err != nil {
		t.Fatalf("[%s] grab conn: %v", o.name, err)
	}
	defer func() { _ = conn.Close() }()

	setup := append([]string{
		"DROP DATABASE IF EXISTS " + applyDB,
		fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s", applyDB, serverCharsetFor(o.version)),
		"USE " + applyDB,
	}, p.from...)
	for _, s := range setup {
		if _, err := conn.ExecContext(ctx, s); err != nil {
			t.Skipf("[%s] could not set up `from` state for %s: %q: %v", o.name, p.id, s, err)
		}
	}
	for _, op := range plan.Ops {
		if _, err := conn.ExecContext(ctx, op.SQL); err != nil {
			t.Fatalf("[%s] APPLY FAILED for %s:\n  stmt: %s\n  err: %v\n  full plan:\n%s",
				o.name, p.id, op.SQL, err, plan.SQL())
		}
	}

	// Read the applied schema back and assert it canonicalizes equal to `to`, both directions.
	resultSDL, ok := o.dumpOrderingSchemaSDL(t, conn, applyDB, applyDB, p.tables, p.views, p.functions)
	if !ok {
		t.Fatalf("[%s] %s: could not read back result", o.name, p.id)
	}
	resultCat, err := LoadSDLWithVersion(resultSDL, o.version)
	if err != nil {
		t.Fatalf("[%s] %s: reload of result failed: %v\n%s", o.name, p.id, err, resultSDL)
	}
	resultDiff := DiffWithNormalizer(resultCat, toCat, n)
	if !resultDiff.IsEmpty() {
		t.Errorf("[%s] APPLY-CORRECTNESS FAILED for %s: result != to\n  plan:\n%s\n  result SDL:\n%s\n  diff: %s",
			o.name, p.id, plan.SQL(), resultSDL, describeOrderingDiff(resultDiff))
	}
	if d := DiffWithNormalizer(toCat, resultCat, n); !d.IsEmpty() {
		t.Errorf("[%s] APPLY-CORRECTNESS (reverse) FAILED for %s: %s", o.name, p.id, describeOrderingDiff(d))
	}

	// CONVERGENCE: the follow-up plan from the applied result to `to` must be EMPTY.
	followUp := GenerateMigrationWithNormalizer(resultCat, toCat, resultDiff, n)
	if sql := followUp.SQL(); sql != "" {
		t.Errorf("[%s] CONVERGENCE FAILED for %s: follow-up plan not empty:\n%s", o.name, p.id, sql)
	}

	// DROP correctness: any view/function in `from` but not in `to` must not survive.
	assertOrderingObjectsGone(t, o, conn, applyDB, p, fromViews, fromFunctions)
}

// assertOrderingObjectsGone asserts the views/functions dropped by the plan no longer exist.
func assertOrderingObjectsGone(t *testing.T, o *oracleConn, conn *sql.Conn, applyDB string, p orderingProbe, fromViews, fromFunctions []string) {
	t.Helper()
	ctx := context.Background()
	inTo := func(names []string, n string) bool {
		for _, x := range names {
			if strings.EqualFold(x, n) {
				return true
			}
		}
		return false
	}
	for _, v := range fromViews {
		if inTo(p.views, v) {
			continue
		}
		var c1, c2, c3, c4 string
		row := conn.QueryRowContext(ctx, fmt.Sprintf("SHOW CREATE VIEW %s.%s", applyDB, v))
		if err := row.Scan(&c1, &c2, &c3, &c4); err == nil {
			t.Errorf("[%s] %s: dropped view %s still exists", o.name, p.id, v)
		}
	}
	for _, f := range fromFunctions {
		if inTo(p.functions, f) {
			continue
		}
		var c1, c2, c3, c4, c5, c6 sql.NullString
		row := conn.QueryRowContext(ctx, fmt.Sprintf("SHOW CREATE FUNCTION %s.%s", applyDB, f))
		if err := row.Scan(&c1, &c2, &c3, &c4, &c5, &c6); err == nil {
			t.Errorf("[%s] %s: dropped function %s still exists", o.name, p.id, f)
		}
	}
}

// dumpOrderingSchemaSDL reads back SHOW CREATE for the named tables, views, AND functions on conn
// and assembles a reloadable SDL string under the logical database. Tables and views route through
// dumpViewSchemaSDL, which re-homes the physical-db qualifier; function readbacks are appended
// as-is — no rewrite is needed because SHOW CREATE FUNCTION returns an unqualified name and a
// verbatim-stored body (the probes never db-qualify), live-verified on both versions, unlike
// SHOW CREATE VIEW whose stored body IS db-qualified.
func (o *oracleConn) dumpOrderingSchemaSDL(t *testing.T, conn *sql.Conn, physDB, logicalDB string, tables, views, functions []string) (string, bool) {
	t.Helper()
	sdl, ok := o.dumpViewSchemaSDL(t, conn, physDB, logicalDB, tables, views)
	if !ok {
		return "", false
	}
	var b strings.Builder
	b.WriteString(sdl)
	for _, fn := range functions {
		rb, ok := o.showCreateRoutine(t, physDB, "FUNCTION", fn)
		if !ok {
			return "", false
		}
		b.WriteString(rb)
		b.WriteString(";\n")
	}
	return b.String(), true
}

// loadOrderingSchemaFromEngine applies a schema to a throwaway physical db, reads back the named
// tables + views + functions, and loads them into a catalog re-homed under the logical db
// (version-aware). Returns nil on any apply/dump failure (the caller skips). An empty setup yields
// an empty catalog.
func loadOrderingSchemaFromEngine(t *testing.T, o *oracleConn, physDB, logicalDB string, setup, tables, views, functions []string) *Catalog {
	t.Helper()
	if len(setup) == 0 {
		return New()
	}
	conn, ok := o.applyViewSchema(t, physDB, setup)
	if !ok {
		return nil
	}
	defer func() { _ = conn.Close() }()
	sdl, ok := o.dumpOrderingSchemaSDL(t, conn, physDB, logicalDB, tables, views, functions)
	if !ok {
		return nil
	}
	cat, err := LoadSDLWithVersion(sdl, o.version)
	if err != nil {
		t.Fatalf("[%s] load ordering-schema-from-engine failed: %v\n%s", o.name, err, sdl)
	}
	return cat
}

// orderingObjectNames extracts the table, view, and function names declared by a list of CREATE
// statements (the shapes the ordering probes use), so the `from` side needs no hand-maintained
// name list.
func orderingObjectNames(setup []string) (tables, views, functions []string) {
	tables, views = schemaObjectNames(setup)
	for _, s := range setup {
		fields := strings.Fields(s)
		if len(fields) >= 3 && strings.EqualFold(fields[0], "CREATE") && strings.EqualFold(fields[1], "FUNCTION") {
			name := fields[2]
			if i := strings.IndexByte(name, '('); i >= 0 {
				name = name[:i]
			}
			functions = append(functions, strings.Trim(name, "`"))
		}
	}
	return tables, views, functions
}

// describeOrderingDiff renders a compact description of a SchemaDiff's view/function/table changes
// for failure output.
func describeOrderingDiff(d *SchemaDiff) string {
	var b strings.Builder
	for _, fe := range d.Functions {
		fmt.Fprintf(&b, "[function %s.%s %s]", fe.Database, fe.Name, fe.Action)
	}
	for _, ve := range d.Views {
		fmt.Fprintf(&b, "[view %s.%s %s]", ve.Database, ve.Name, ve.Action)
	}
	for _, te := range d.Tables {
		fmt.Fprintf(&b, "[table %s %s]", te.Name, te.Action)
	}
	return b.String()
}
