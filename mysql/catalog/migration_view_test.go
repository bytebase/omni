package catalog

import (
	"strings"
	"testing"
)

// Hermetic unit tests for the view generator: op types / phases / priorities, the rendered
// CREATE OR REPLACE / DROP DDL shape, and the view-on-view CREATE ordering. Apply-correctness
// against the live engines lives in migration_view_oracle_test.go.

// viewPlanFor builds a from→to plan via the public Diff + GenerateMigration entry points at a
// version, wrapping each SDL fragment in the `vt` database (matching loadViewCatalog) so view
// identity is (vt, name).
func viewPlanFor(t *testing.T, version Version, fromSDL, toSDL string) *MigrationPlan {
	t.Helper()
	from := loadViewCatalog(t, version, fromSDL)
	to := loadViewCatalog(t, version, toSDL)
	n := NormalizerFor(version)
	diff := DiffWithNormalizer(from, to, n)
	return GenerateMigrationWithNormalizer(from, to, diff, n)
}

func viewOps(plan *MigrationPlan) []MigrationOp {
	var out []MigrationOp
	for _, op := range plan.Ops {
		if op.Type == OpCreateView || op.Type == OpDropView {
			out = append(out, op)
		}
	}
	return out
}

func TestGenerateView_NoOpEmpty(t *testing.T) {
	for _, version := range []Version{MySQL57, MySQL80} {
		sdl := `CREATE TABLE t (a INT, b INT);
CREATE VIEW v1 AS SELECT a FROM t;`
		plan := viewPlanFor(t, version, sdl, sdl)
		if sql := plan.SQL(); sql != "" {
			t.Errorf("[v%d] no-op plan must be empty, got:\n%s", version, sql)
		}
	}
}

func TestGenerateView_CreateShape(t *testing.T) {
	plan := viewPlanFor(t, MySQL80,
		`CREATE TABLE t (a INT, b INT);`,
		`CREATE TABLE t (a INT, b INT);
CREATE VIEW v1 AS SELECT a FROM t;`)
	ops := viewOps(plan)
	if len(ops) != 1 || ops[0].Type != OpCreateView {
		t.Fatalf("want one OpCreateView, got %+v", ops)
	}
	op := ops[0]
	if op.Phase != PhaseMain || op.Priority != priorityView {
		t.Errorf("create op phase/priority = %v/%d, want PhaseMain/%d", op.Phase, op.Priority, priorityView)
	}
	for _, want := range []string{
		"CREATE OR REPLACE",
		"ALGORITHM=UNDEFINED",
		"SQL SECURITY DEFINER",
		"VIEW `vt`.`v1`",
		" AS select `t`.`a` AS `a` from `t`",
	} {
		if !strings.Contains(op.SQL, want) {
			t.Errorf("create SQL missing %q:\n%s", want, op.SQL)
		}
	}
	// The DEFINER= clause must NOT be emitted (ignore-in-diff + least-privilege apply hazard).
	// (Note: "SQL SECURITY DEFINER" legitimately contains the word DEFINER; check the `DEFINER=`
	// token specifically.)
	if strings.Contains(op.SQL, "DEFINER=") {
		t.Errorf("create SQL must not emit a DEFINER= clause:\n%s", op.SQL)
	}
}

func TestGenerateView_CreateWithOptions(t *testing.T) {
	plan := viewPlanFor(t, MySQL80,
		`CREATE TABLE t (a INT, b INT);`,
		`CREATE TABLE t (a INT, b INT);
CREATE ALGORITHM=MERGE SQL SECURITY INVOKER VIEW v1 (p, q) AS SELECT a, b FROM t WHERE a > 0 WITH CASCADED CHECK OPTION;`)
	ops := viewOps(plan)
	if len(ops) != 1 {
		t.Fatalf("want one view op, got %+v", ops)
	}
	for _, want := range []string{
		"ALGORITHM=MERGE",
		"SQL SECURITY INVOKER",
		"VIEW `vt`.`v1` (`p`,`q`)",
		"WITH CASCADED CHECK OPTION",
	} {
		if !strings.Contains(ops[0].SQL, want) {
			t.Errorf("create SQL missing %q:\n%s", want, ops[0].SQL)
		}
	}
}

func TestGenerateView_ModifyEmitsCreateOrReplace(t *testing.T) {
	plan := viewPlanFor(t, MySQL80,
		`CREATE TABLE t (a INT, b INT);
CREATE VIEW v1 AS SELECT a FROM t;`,
		`CREATE TABLE t (a INT, b INT);
CREATE VIEW v1 AS SELECT b FROM t;`)
	ops := viewOps(plan)
	if len(ops) != 1 || ops[0].Type != OpCreateView {
		t.Fatalf("modify should emit one OpCreateView (CREATE OR REPLACE), got %+v", ops)
	}
	if !strings.HasPrefix(ops[0].SQL, "CREATE OR REPLACE") {
		t.Errorf("modify SQL should start with CREATE OR REPLACE:\n%s", ops[0].SQL)
	}
}

func TestGenerateView_DropShape(t *testing.T) {
	plan := viewPlanFor(t, MySQL80,
		`CREATE TABLE t (a INT);
CREATE VIEW v1 AS SELECT a FROM t;`,
		`CREATE TABLE t (a INT);`)
	ops := viewOps(plan)
	if len(ops) != 1 || ops[0].Type != OpDropView {
		t.Fatalf("want one OpDropView, got %+v", ops)
	}
	op := ops[0]
	if op.Phase != PhasePre {
		t.Errorf("drop op phase = %v, want PhasePre", op.Phase)
	}
	if op.SQL != "DROP VIEW IF EXISTS `vt`.`v1`" {
		t.Errorf("drop SQL = %q", op.SQL)
	}
}

// View-on-view CREATE must order the dependency before the dependent. v_base references table t;
// v_dep references v_base, so v_base's CREATE must precede v_dep's in the plan order.
func TestGenerateView_ViewOnViewOrdering(t *testing.T) {
	plan := viewPlanFor(t, MySQL80,
		`CREATE TABLE t (a INT, b INT);`,
		`CREATE TABLE t (a INT, b INT);
CREATE VIEW v_dep AS SELECT x FROM v_base;
CREATE VIEW v_base (x) AS SELECT a FROM t;`)
	ops := viewOps(plan)
	if len(ops) != 2 {
		t.Fatalf("want 2 view ops, got %+v", ops)
	}
	posBase, posDep := -1, -1
	for i, op := range ops {
		switch strings.ToLower(op.ObjectName) {
		case "v_base":
			posBase = i
		case "v_dep":
			posDep = i
		}
	}
	if posBase < 0 || posDep < 0 {
		t.Fatalf("missing ops: base=%d dep=%d (%+v)", posBase, posDep, ops)
	}
	if posBase > posDep {
		t.Errorf("v_base (dependency) must be created before v_dep (dependent):\n%s", plan.SQL())
	}
}

// A three-level chain (v3 → v2 → v1 → t) must come out strictly in dependency order.
func TestGenerateView_ViewChainOrdering(t *testing.T) {
	plan := viewPlanFor(t, MySQL80,
		`CREATE TABLE t (a INT);`,
		`CREATE TABLE t (a INT);
CREATE VIEW v3 AS SELECT a FROM v2;
CREATE VIEW v1 AS SELECT a FROM t;
CREATE VIEW v2 AS SELECT a FROM v1;`)
	ops := viewOps(plan)
	order := make(map[string]int)
	for i, op := range ops {
		order[strings.ToLower(op.ObjectName)] = i
	}
	if order["v1"] >= order["v2"] || order["v2"] >= order["v3"] {
		t.Errorf("chain order wrong: v1=%d v2=%d v3=%d\n%s",
			order["v1"], order["v2"], order["v3"], plan.SQL())
	}
}

// Regression (review blocker): a column ALIAS that matches a batch view's name must NOT create a
// false dependency edge that reverses CREATE order. View `a` has alias `b` but does NOT reference
// view `b`; view `b` references view `a`. The only real edge is b→a, so `a` must be created first.
// An earlier substring-based scan saw `a`'s alias `\`b\“ as a reference to view `b`, forming a
// false cycle that put `b` before `a` — which fails to apply (a doesn't exist yet).
func TestGenerateView_AliasMatchingViewNameNoFalseEdge(t *testing.T) {
	plan := viewPlanFor(t, MySQL80,
		`CREATE TABLE t (id INT, x INT);`,
		`CREATE TABLE t (id INT, x INT);
CREATE VIEW a AS SELECT id AS b FROM t;
CREATE VIEW b AS SELECT x FROM a;`)
	ops := viewOps(plan)
	posA, posB := -1, -1
	for i, op := range ops {
		switch strings.ToLower(op.ObjectName) {
		case "a":
			posA = i
		case "b":
			posB = i
		}
	}
	if posA < 0 || posB < 0 {
		t.Fatalf("missing ops: a=%d b=%d (%+v)", posA, posB, ops)
	}
	if posA > posB {
		t.Errorf("view a (dependency of b) must be created before b — a column alias named like a view must not reverse order:\n%s", plan.SQL())
	}
}

// Regression: a parenthesized JOIN of two batch views must order both before a view that joins
// them. Exercises relation extraction at the `(` anchor (join arms are parenthesized).
func TestGenerateView_ParenJoinOfViewsOrdering(t *testing.T) {
	plan := viewPlanFor(t, MySQL80,
		`CREATE TABLE t (id INT); CREATE TABLE u (id INT);`,
		`CREATE TABLE t (id INT); CREATE TABLE u (id INT);
CREATE VIEW v1 AS SELECT id FROM t;
CREATE VIEW v2 AS SELECT id FROM u;
CREATE VIEW jv AS SELECT v1.id, v2.id AS id2 FROM v1 JOIN v2 ON v1.id = v2.id;`)
	ops := viewOps(plan)
	order := make(map[string]int)
	for i, op := range ops {
		order[strings.ToLower(op.ObjectName)] = i
	}
	if order["v1"] >= order["jv"] || order["v2"] >= order["jv"] {
		t.Errorf("jv (joins v1,v2) must come after both: v1=%d v2=%d jv=%d\n%s",
			order["v1"], order["v2"], order["jv"], plan.SQL())
	}
}

// Regression (re-review): dependency extraction is AST-based, so a column alias / string literal /
// ON-clause column / CTE name that coincides with a batch view name must NOT distort ordering.
// Here `dep` references base view `b`; a sibling chain has a view literally named like a column the
// other views alias. All real edges must be respected and no false edge may invert order.
func TestGenerateView_DependencyExtractionRobust(t *testing.T) {
	// base `b` (referenced by `dep`); view `c` aliases a column `b` and joins t with a literal that
	// mentions `b`; view `dep` truly depends on `b`. `b` must precede `dep`.
	plan := viewPlanFor(t, MySQL80,
		`CREATE TABLE t (id INT, x INT);`,
		`CREATE TABLE t (id INT, x INT);
CREATE VIEW b AS SELECT id FROM t;
CREATE VIEW c AS SELECT id AS b, 'from `+"`b`"+`' AS note FROM t;
CREATE VIEW dep AS SELECT id FROM b;`)
	ops := viewOps(plan)
	order := make(map[string]int)
	for i, op := range ops {
		order[strings.ToLower(op.ObjectName)] = i
	}
	// dep depends on b → b first. c depends only on t (its alias `b` and literal `b` are NOT view b).
	if order["b"] >= order["dep"] {
		t.Errorf("b (dependency of dep) must precede dep: b=%d dep=%d\n%s", order["b"], order["dep"], plan.SQL())
	}
}

// Regression (re-review): a CTE whose name matches a batch view must not create a false edge. View
// `ctev` has a CTE named `shared` AND there is a sibling view actually named `shared`; `ctev` does
// NOT depend on the view `shared` (its `FROM shared` binds to the CTE), so no ordering constraint
// links them — and no false cycle is formed.
func TestGenerateView_CTENameNotAViewEdge(t *testing.T) {
	plan := viewPlanFor(t, MySQL80,
		`CREATE TABLE t (id INT);`,
		`CREATE TABLE t (id INT);
CREATE VIEW shared AS SELECT id FROM t;
CREATE VIEW ctev AS WITH shared AS (SELECT id FROM t) SELECT id FROM shared;`)
	// Must generate without forming a cycle / panic, and both views must be present.
	ops := viewOps(plan)
	if len(ops) != 2 {
		t.Fatalf("want 2 view ops, got %+v", ops)
	}
	// ctev references the CTE `shared`, not the view `shared`, plus table t — so ctev carries no
	// batch-view edge. The key property is that no false cycle inverts anything (both depth 0).
}

// Regression (re-review): a `AS TABLE <view>` body (MySQL 8.0.19+ TABLE query primary) deparses to
// `table \`v\“ and re-parses to a *TableStmt, not a *SelectStmt. The dependency extractor must
// still find the referenced view so the ordering holds — `tv AS TABLE base` must be created after
// `base`.
func TestGenerateView_TableQueryPrimaryOrdering(t *testing.T) {
	plan := viewPlanFor(t, MySQL80,
		`CREATE TABLE t (id INT);`,
		`CREATE TABLE t (id INT);
CREATE VIEW base AS SELECT id FROM t;
CREATE VIEW tv AS TABLE base;`)
	ops := viewOps(plan)
	order := make(map[string]int)
	for i, op := range ops {
		order[strings.ToLower(op.ObjectName)] = i
	}
	if order["base"] >= order["tv"] {
		t.Errorf("base (dependency) must precede tv (AS TABLE base): base=%d tv=%d\n%s",
			order["base"], order["tv"], plan.SQL())
	}
}

// A view over a freshly created table must be created after the table (CREATE TABLE at
// priorityTable precedes CREATE VIEW at priorityView in the global sort).
func TestGenerateView_ViewAfterTable(t *testing.T) {
	plan := viewPlanFor(t, MySQL80,
		``,
		`CREATE TABLE t (a INT, b INT);
CREATE VIEW v1 AS SELECT a FROM t;`)
	posTable, posView := -1, -1
	for i, op := range plan.Ops {
		if op.Type == OpCreateTable && strings.EqualFold(op.ObjectName, "t") {
			posTable = i
		}
		if op.Type == OpCreateView && strings.EqualFold(op.ObjectName, "v1") {
			posView = i
		}
	}
	if posTable < 0 || posView < 0 {
		t.Fatalf("missing ops: table=%d view=%d", posTable, posView)
	}
	if posTable > posView {
		t.Errorf("CREATE TABLE must precede CREATE VIEW:\n%s", plan.SQL())
	}
}

// opPositions returns the plan index of the first op matching each (type, name) request, -1 when
// absent. Names are compared case-insensitively.
func opPositions(plan *MigrationPlan, wants []struct {
	typ  MigrationOpType
	name string
}) []int {
	pos := make([]int, len(wants))
	for i := range pos {
		pos[i] = -1
	}
	for i, op := range plan.Ops {
		for w, want := range wants {
			if pos[w] < 0 && op.Type == want.typ && strings.EqualFold(op.ObjectName, want.name) {
				pos[w] = i
			}
		}
	}
	return pos
}

// A view whose body calls a stored function created in the SAME plan must be created AFTER the
// function: MySQL eagerly validates function references at CREATE VIEW time (Error 1305 "FUNCTION
// does not exist", live-verified on 8.0.32 and 5.7.25), while a routine body is lazily validated
// (a CREATE FUNCTION referencing a missing function/view/table is accepted, live-verified) — so
// routine creates are hoisted before view creates unconditionally. Regression for the enterprise
// SDL smoke failure (view-before-function plan → Error 1305 on apply).
func TestGenerateView_ViewAfterFunction(t *testing.T) {
	for _, version := range []Version{MySQL57, MySQL80} {
		plan := viewPlanFor(t, version,
			`CREATE TABLE users (id INT PRIMARY KEY);`,
			`CREATE TABLE users (id INT PRIMARY KEY);
CREATE FUNCTION ent_ucount() RETURNS INT READS SQL DATA RETURN (SELECT COUNT(*) FROM users);
CREATE VIEW ent_v_fn AS SELECT ent_ucount() AS n;`)
		pos := opPositions(plan, []struct {
			typ  MigrationOpType
			name string
		}{{OpCreateFunction, "ent_ucount"}, {OpCreateView, "ent_v_fn"}})
		if pos[0] < 0 || pos[1] < 0 {
			t.Fatalf("[v%d] missing ops: function=%d view=%d\n%s", version, pos[0], pos[1], plan.SQL())
		}
		if pos[0] > pos[1] {
			t.Errorf("[v%d] CREATE FUNCTION must precede the CREATE VIEW that calls it:\n%s", version, plan.SQL())
		}
	}
}

// The inverse drop direction: when a view and the function it calls are BOTH dropped in one plan,
// the view is dropped first (drop the dependent before its dependency). MySQL tolerates dropping a
// function a view still references (live-verified — the view merely goes invalid), so this is the
// safe semantic order rather than an apply-success requirement.
func TestGenerateView_DropViewBeforeDropFunction(t *testing.T) {
	plan := viewPlanFor(t, MySQL80,
		`CREATE TABLE users (id INT PRIMARY KEY);
CREATE FUNCTION ent_ucount() RETURNS INT READS SQL DATA RETURN (SELECT COUNT(*) FROM users);
CREATE VIEW ent_v_fn AS SELECT ent_ucount() AS n;`,
		`CREATE TABLE users (id INT PRIMARY KEY);`)
	pos := opPositions(plan, []struct {
		typ  MigrationOpType
		name string
	}{{OpDropView, "ent_v_fn"}, {OpDropFunction, "ent_ucount"}})
	if pos[0] < 0 || pos[1] < 0 {
		t.Fatalf("missing ops: view=%d function=%d\n%s", pos[0], pos[1], plan.SQL())
	}
	if pos[0] > pos[1] {
		t.Errorf("DROP VIEW must precede DROP FUNCTION of a function it calls:\n%s", plan.SQL())
	}
}

// Transitive: function ← view ← view. Combining routine-before-view (priority) with view-on-view
// depth ordering must yield fn < v1 < v2.
func TestGenerateView_FunctionViewViewChainOrdering(t *testing.T) {
	plan := viewPlanFor(t, MySQL80,
		`CREATE TABLE t (a INT);`,
		`CREATE TABLE t (a INT);
CREATE FUNCTION fn_base() RETURNS INT READS SQL DATA RETURN (SELECT COUNT(*) FROM t);
CREATE VIEW v2 AS SELECT n FROM v1;
CREATE VIEW v1 AS SELECT fn_base() AS n;`)
	pos := opPositions(plan, []struct {
		typ  MigrationOpType
		name string
	}{{OpCreateFunction, "fn_base"}, {OpCreateView, "v1"}, {OpCreateView, "v2"}})
	if pos[0] < 0 || pos[1] < 0 || pos[2] < 0 {
		t.Fatalf("missing ops: fn=%d v1=%d v2=%d\n%s", pos[0], pos[1], pos[2], plan.SQL())
	}
	if pos[0] >= pos[1] || pos[1] >= pos[2] {
		t.Errorf("chain order wrong (want fn < v1 < v2): fn=%d v1=%d v2=%d\n%s",
			pos[0], pos[1], pos[2], plan.SQL())
	}
}

// A body-modified function (DROP+CREATE path) with a NEW view calling it: the function's re-CREATE
// (PhaseMain) must still precede the view's CREATE, and its DROP stays in PhasePre.
func TestGenerateView_ViewAfterModifiedFunction(t *testing.T) {
	plan := viewPlanFor(t, MySQL80,
		`CREATE TABLE t (a INT);
CREATE FUNCTION fn_base() RETURNS INT READS SQL DATA RETURN (SELECT COUNT(*) FROM t);`,
		`CREATE TABLE t (a INT);
CREATE FUNCTION fn_base() RETURNS INT READS SQL DATA RETURN (SELECT COUNT(*) + 1 FROM t);
CREATE VIEW v1 AS SELECT fn_base() AS n;`)
	pos := opPositions(plan, []struct {
		typ  MigrationOpType
		name string
	}{{OpDropFunction, "fn_base"}, {OpCreateFunction, "fn_base"}, {OpCreateView, "v1"}})
	if pos[0] < 0 || pos[1] < 0 || pos[2] < 0 {
		t.Fatalf("missing ops: drop=%d create=%d view=%d\n%s", pos[0], pos[1], pos[2], plan.SQL())
	}
	if pos[0] >= pos[1] || pos[1] >= pos[2] {
		t.Errorf("order wrong (want DROP FUNCTION < CREATE FUNCTION < CREATE VIEW): drop=%d create=%d view=%d\n%s",
			pos[0], pos[1], pos[2], plan.SQL())
	}
}

// All drops precede all creates: a view dropped and another created in one plan must not interleave
// such that a create lands before the drops (PhasePre vs PhaseMain).
func TestGenerateView_DropsBeforeCreates(t *testing.T) {
	plan := viewPlanFor(t, MySQL80,
		`CREATE TABLE t (a INT);
CREATE VIEW old_v AS SELECT a FROM t;`,
		`CREATE TABLE t (a INT);
CREATE VIEW new_v AS SELECT a FROM t;`)
	lastDrop, firstCreate := -1, len(plan.Ops)
	for i, op := range plan.Ops {
		if op.Type == OpDropView {
			if i > lastDrop {
				lastDrop = i
			}
		}
		if op.Type == OpCreateView && i < firstCreate {
			firstCreate = i
		}
	}
	if lastDrop >= 0 && lastDrop > firstCreate {
		t.Errorf("a view DROP sorted after a view CREATE:\n%s", plan.SQL())
	}
}
