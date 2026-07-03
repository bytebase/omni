package catalog

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

// TestOracle_SysSchemaDogfood is the north-star gate for the temporal-unit
// keyword-argument deparse fix: dogfood-load the STOCK MySQL sys schema.
//
// From a live readback of every sys view (100 views on 8.0.32, including
// innodb_lock_waits and the x$ variants whose bodies carry
// timestampdiff(SECOND,...)):
//
//  1. LoadSDL the canonical dump (verbatim SHOW CREATE VIEW text) — catA;
//     catA must self-diff empty.
//  2. GenerateMigration(empty → catA) — one CREATE OR REPLACE VIEW per view.
//  3. APPLY the full plan to a scratch database on the live engine (each
//     statement rehomed to the scratch db by view name only; body references
//     keep pointing at the real sys/performance_schema objects, which is
//     exactly what the plan's deparse fidelity is being tested against).
//     Before the fix this leg failed with error 1064: the plan rendered
//     timestampdiff(`SECOND`,...) with a quoted unit.
//  4. Re-dump the scratch database and LoadSDL it — catB.
//  5. catA vs catB must diff EMPTY in both directions (the readback →
//     generate → apply → re-dump loop is a fixed point).
//
// 8.0 only: the 5.7 box ships sys 1.5.1 whose view bodies exercise an
// unrelated surface; the enterprise-baseline harness leg this test
// replicates (TestSDLEnterpriseBaseline/mysql80/sys) targets 8.0.
func TestOracle_SysSchemaDogfood(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	o := connectOracle(t, MySQL80)
	n := NormalizerFor(MySQL80)
	ctx := context.Background()

	// Enumerate and read back every sys view, verbatim.
	rows, err := o.db.QueryContext(ctx,
		"SELECT table_name FROM information_schema.views WHERE table_schema='sys' ORDER BY table_name")
	if err != nil {
		t.Skipf("[%s] cannot enumerate sys views: %v", o.name, err)
	}
	var viewNames []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("[%s] scan view name: %v", o.name, err)
		}
		viewNames = append(viewNames, name)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("[%s] enumerate sys views: %v", o.name, err)
	}
	_ = rows.Close()
	if len(viewNames) == 0 {
		t.Skipf("[%s] sys schema has no views (not installed?)", o.name)
	}

	var dump strings.Builder
	dump.WriteString("CREATE DATABASE sys;\nUSE sys;\n")
	for _, name := range viewNames {
		var vn, ddl, cs, coll string
		row := o.db.QueryRowContext(ctx, fmt.Sprintf("SHOW CREATE VIEW sys.`%s`", name))
		if err := row.Scan(&vn, &ddl, &cs, &coll); err != nil {
			t.Fatalf("[%s] SHOW CREATE VIEW sys.%s: %v", o.name, name, err)
		}
		dump.WriteString(ddl)
		dump.WriteString(";\n")
	}

	// Gate 1: the canonical dump loads and self-diffs empty.
	catA, err := LoadSDLWithVersion(dump.String(), MySQL80)
	if err != nil {
		t.Fatalf("[%s] LoadSDL(sys dump) failed: %v", o.name, err)
	}
	dbA := catA.GetDatabase("sys")
	if dbA == nil || len(dbA.Views) != len(viewNames) {
		t.Fatalf("[%s] loaded catalog has %d sys views, want %d", o.name, len(dbA.Views), len(viewNames))
	}
	if d := DiffWithNormalizer(catA, catA, n); !d.IsEmpty() {
		t.Fatalf("[%s] sys dump self-diff not empty: %s", o.name, describeViewDiff(d))
	}

	// Gate 2: generate the full create plan from empty.
	empty := New()
	diff := DiffWithNormalizer(empty, catA, n)
	plan := GenerateMigrationWithNormalizer(empty, catA, diff, n)
	var viewOps int
	for _, op := range plan.Ops {
		if strings.Contains(op.SQL, " VIEW ") {
			viewOps++
		}
	}
	if viewOps != len(viewNames) {
		t.Fatalf("[%s] plan has %d view statements, want %d\n%s", o.name, viewOps, len(viewNames), plan.SQL())
	}

	// Apply the plan to a scratch database. Each statement is rehomed by
	// view NAME only — body references stay on the real sys /
	// performance_schema objects, so every generated body is parsed and
	// validated by the live engine exactly as rendered.
	const scratch = "tudsdl_sysdf"
	conn, err := o.db.Conn(ctx)
	if err != nil {
		t.Fatalf("[%s] grab conn: %v", o.name, err)
	}
	defer func() { _ = conn.Close() }()
	defer func() {
		_, _ = conn.ExecContext(ctx, "DROP DATABASE IF EXISTS "+scratch)
	}()
	// The session default database is `sys`, not the scratch db: stored view
	// bodies keep FUNCTION references unqualified (sys_config getters,
	// format_path, ...), and the engine resolves those against the default
	// database at CREATE VIEW time. The enterprise harness this replicates
	// applies into a database that carries the sys routines; pointing
	// resolution at the real sys schema is the in-omni equivalent. Every
	// CREATE statement still explicitly targets the scratch db.
	for _, s := range []string{
		"DROP DATABASE IF EXISTS " + scratch,
		"CREATE DATABASE " + scratch,
		"USE sys",
	} {
		if _, err := conn.ExecContext(ctx, s); err != nil {
			t.Fatalf("[%s] scratch setup %q: %v", o.name, s, err)
		}
	}
	for _, op := range plan.Ops {
		if !strings.Contains(op.SQL, " VIEW ") {
			continue // the CREATE DATABASE op; the scratch db replaces it
		}
		stmt := strings.Replace(op.SQL, "VIEW `sys`.`", "VIEW `"+scratch+"`.`", 1)
		if stmt == op.SQL {
			t.Fatalf("[%s] could not rehome plan statement: %s", o.name, op.SQL)
		}
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("[%s] APPLY FAILED (the dogfood leg):\n  stmt: %.300s\n  err: %v", o.name, stmt, err)
		}
	}

	// Gate 3: re-dump the scratch database and require a fixed point.
	var redump strings.Builder
	redump.WriteString("CREATE DATABASE sys;\nUSE sys;\n")
	for _, name := range viewNames {
		var vn, ddl, cs, coll string
		row := conn.QueryRowContext(ctx, fmt.Sprintf("SHOW CREATE VIEW %s.`%s`", scratch, name))
		if err := row.Scan(&vn, &ddl, &cs, &coll); err != nil {
			t.Fatalf("[%s] re-dump SHOW CREATE VIEW %s.%s: %v", o.name, scratch, name, err)
		}
		rehomed := strings.Replace(ddl, "VIEW `"+scratch+"`.`", "VIEW `sys`.`", 1)
		if rehomed == ddl {
			t.Fatalf("[%s] could not rehome re-dumped view %s: %.200s", o.name, name, ddl)
		}
		redump.WriteString(rehomed)
		redump.WriteString(";\n")
	}
	catB, err := LoadSDLWithVersion(redump.String(), MySQL80)
	if err != nil {
		t.Fatalf("[%s] LoadSDL(re-dump) failed: %v", o.name, err)
	}
	if d := DiffWithNormalizer(catA, catB, n); !d.IsEmpty() {
		reportSysViewDrift(t, catA, catB)
		t.Fatalf("[%s] NORTH STAR: dump vs re-dump not empty: %s", o.name, describeViewDiff(d))
	}
	if d := DiffWithNormalizer(catB, catA, n); !d.IsEmpty() {
		reportSysViewDrift(t, catA, catB)
		t.Fatalf("[%s] NORTH STAR (reverse): re-dump vs dump not empty: %s", o.name, describeViewDiff(d))
	}
}

// reportSysViewDrift prints per-view Definition differences for failure triage.
func reportSysViewDrift(t *testing.T, catA, catB *Catalog) {
	t.Helper()
	dbA := catA.GetDatabase("sys")
	dbB := catB.GetDatabase("sys")
	if dbA == nil || dbB == nil {
		return
	}
	for name, va := range dbA.Views {
		vb := dbB.Views[name]
		if vb == nil {
			t.Logf("view %s missing after re-dump", name)
			continue
		}
		if va.Definition != vb.Definition {
			t.Logf("view %s drifted:\n  dump:    %s\n  re-dump: %s", name, va.Definition, vb.Definition)
		}
	}
}
