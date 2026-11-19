package catalog

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

// TestViewGroupByAllDeparse pins Redshift's GROUP BY ALL through the catalog:
// the clause carries no list of its own, and a stored view keeps it rather than
// deparsing as a query with no grouping at all.
func TestViewGroupByAllDeparse(t *testing.T) {
	c := setupPhase4(t)
	def := viewDef(t, c, "CREATE VIEW v_group_all AS SELECT name, sum(val) FROM t1 GROUP BY ALL")
	if !strings.Contains(def, "GROUP BY ALL") {
		t.Fatalf("view definition = %q, want GROUP BY ALL in it", def)
	}

	// ALL beside a list keeps both: the word and the grouping it modifies.
	quantified := viewDef(t, c, "CREATE VIEW v_group_all_list AS SELECT name, sum(val) FROM t1 GROUP BY ALL name")
	if !strings.Contains(quantified, "GROUP BY ALL name") {
		t.Fatalf("view definition = %q, want GROUP BY ALL name", quantified)
	}

	// A grouped query that names its columns still deparses with the list.
	plain := viewDef(t, c, "CREATE VIEW v_group_list AS SELECT name, sum(val) FROM t1 GROUP BY name")
	if !strings.Contains(plain, "GROUP BY name") || strings.Contains(plain, "ALL") {
		t.Fatalf("view definition = %q, want GROUP BY name", plain)
	}
}

// TestViewConnectByDeparse pins Redshift's START WITH / CONNECT BY through the
// catalog, PRIOR included: the conditions are analyzed, deparsed back in their
// place, and their columns count as the view's dependencies.
func TestViewConnectByDeparse(t *testing.T) {
	c := setupPhase4(t)
	def := viewDef(t, c, "CREATE VIEW v_hier AS SELECT id, name FROM t1 START WITH name IS NULL CONNECT BY PRIOR id = id")
	if !strings.Contains(def, "START WITH") {
		t.Fatalf("view definition = %q, want START WITH in it", def)
	}
	// The deparser parenthesizes the prefix operator, so the condition reads
	// (PRIOR id) = id: the word and its operand are what matter.
	if !strings.Contains(def, "CONNECT BY") || !strings.Contains(def, "PRIOR id") {
		t.Fatalf("view definition = %q, want CONNECT BY PRIOR id in it", def)
	}

	// A CONNECT BY with no START WITH deparses with the one clause it has.
	only := viewDef(t, c, "CREATE VIEW v_hier_connect AS SELECT id FROM t1 CONNECT BY id > 0")
	if strings.Contains(only, "START WITH") || !strings.Contains(only, "CONNECT BY") || !strings.Contains(only, "id > 0") {
		t.Fatalf("view definition = %q, want CONNECT BY alone", only)
	}
}

// TestMigrationMatViewKeepsItsStorage pins a materialized view's storage
// through the catalog and its migration: the USING access method and the option
// list are kept with the relation, and re-creating it writes them back rather
// than converting an Iceberg view into a plain one.
func TestMigrationMatViewKeepsItsStorage(t *testing.T) {
	from, err := LoadSQL(`CREATE TABLE t (id int, d date, customer_id int);`)
	if err != nil {
		t.Fatal(err)
	}
	to, err := LoadSQL(`CREATE TABLE t (id int, d date, customer_id int);
CREATE MATERIALIZED VIEW mv BACKUP NO DISTSTYLE KEY DISTKEY (customer_id) COMPOUND SORTKEY (d) AUTO REFRESH YES INTERVAL 15 AS SELECT id, d, customer_id FROM t;
CREATE MATERIALIZED VIEW ice USING ICEBERG LOCATION 's3://bucket/sales/' PARTITIONED BY (d, YEAR(d), BUCKET(16, customer_id)) TABLE PROPERTIES ('format' = 'parquet') AS SELECT id, d, customer_id FROM t;`)
	if err != nil {
		t.Fatal(err)
	}

	plan := GenerateMigration(from, to, Diff(from, to))
	byName := map[string]string{}
	for _, op := range plan.Ops {
		if op.Type != OpCreateView || !strings.Contains(op.SQL, "CREATE MATERIALIZED VIEW") {
			continue
		}
		for _, name := range []string{`"mv"`, `"ice"`} {
			if strings.Contains(op.SQL, name) {
				byName[name] = op.SQL
			}
		}
	}

	mv := byName[`"mv"`]
	if !strings.Contains(mv, "DISTKEY") {
		t.Errorf("materialized view op = %q, want its DISTKEY in it", mv)
	}
	for _, want := range []string{
		"BACKUP NO",
		"DISTSTYLE KEY",
		`DISTKEY ("customer_id")`,
		`COMPOUND SORTKEY ("d")`,
		"AUTO REFRESH YES",
		"INTERVAL 15",
	} {
		if !strings.Contains(mv, want) {
			t.Errorf("materialized view op = %q, want %q in it", mv, want)
		}
	}

	ice := byName[`"ice"`]
	for _, want := range []string{
		"USING ICEBERG",
		"LOCATION 's3://bucket/sales/'",
		`PARTITIONED BY ("d", YEAR("d"), BUCKET(16, "customer_id"))`,
		"TABLE PROPERTIES ('format' = 'parquet')",
	} {
		if !strings.Contains(ice, want) {
			t.Errorf("Iceberg materialized view op = %q, want %q in it", ice, want)
		}
	}
}

// TestMigrationMatViewBareDistkey pins the bare DISTKEY the attribute parser
// accepts: a migration writes it back without a column list rather than
// failing on it.
func TestMigrationMatViewBareDistkey(t *testing.T) {
	from, err := LoadSQL(`CREATE TABLE t (id int);`)
	if err != nil {
		t.Fatal(err)
	}
	to, err := LoadSQL(`CREATE TABLE t (id int);
CREATE MATERIALIZED VIEW m DISTKEY AS SELECT id FROM t;`)
	if err != nil {
		t.Fatal(err)
	}
	plan := GenerateMigration(from, to, Diff(from, to))
	found := false
	for _, op := range plan.Ops {
		if strings.Contains(op.SQL, "CREATE MATERIALIZED VIEW") && strings.Contains(op.SQL, `"m"`) {
			found = true
			if !strings.Contains(op.SQL, "DISTKEY AS") {
				t.Errorf("materialized view op = %q, want a bare DISTKEY in it", op.SQL)
			}
		}
	}
	if !found {
		t.Fatal("expected a CREATE MATERIALIZED VIEW op for m")
	}
}

// TestMigrationMatViewStorageChange pins the storage in the diff: a
// materialized view whose options changed is modified, and re-creating it keeps
// the new storage.
func TestMigrationMatViewStorageChange(t *testing.T) {
	from, err := LoadSQL(`CREATE TABLE t (id int);
CREATE MATERIALIZED VIEW mv BACKUP NO AS SELECT id FROM t;`)
	if err != nil {
		t.Fatal(err)
	}
	to, err := LoadSQL(`CREATE TABLE t (id int);
CREATE MATERIALIZED VIEW mv BACKUP YES AS SELECT id FROM t;`)
	if err != nil {
		t.Fatal(err)
	}

	diff := Diff(from, to)
	if len(diff.Relations) != 1 || diff.Relations[0].Action != DiffModify {
		t.Fatalf("diff = %#v, want the materialized view modified", diff.Relations)
	}

	plan := GenerateMigration(from, to, diff)
	recreated := ""
	for _, op := range plan.Ops {
		if strings.Contains(op.SQL, "CREATE MATERIALIZED VIEW") {
			recreated = op.SQL
		}
	}
	if !strings.Contains(recreated, "BACKUP YES") {
		t.Fatalf("materialized view op = %q, want BACKUP YES in it", recreated)
	}
}

// TestViewGroupByAllGroupingSets pins ALL beside grouping sets: the catalog
// records the sets, and the word that modified them has to come back too.
func TestViewGroupByAllGroupingSets(t *testing.T) {
	c := setupPhase4(t)
	def := viewDef(t, c, "CREATE VIEW v_group_sets AS SELECT name, sum(val) FROM t1 GROUP BY ALL ROLLUP(name)")
	if !strings.Contains(def, "GROUP BY ALL ROLLUP(name)") {
		t.Fatalf("view definition = %q, want GROUP BY ALL ROLLUP(name)", def)
	}
}

// TestViewDependenciesFromHierarchicalConditions pins the dependency a column
// carries when it appears only in START WITH or CONNECT BY: dropping it then
// reaches the view that reads it, as a column the SELECT list names would.
func TestViewDependenciesFromHierarchicalConditions(t *testing.T) {
	c := New()
	for _, stmt := range parseStmts(t, `CREATE TABLE t (id int, parent_id int);`) {
		if err := c.ProcessUtility(stmt); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}
	for _, stmt := range parseStmts(t, `CREATE VIEW v_hier_dep AS SELECT id FROM t CONNECT BY PRIOR id = parent_id;`) {
		if err := c.ProcessUtility(stmt); err != nil {
			t.Fatalf("create view: %v", err)
		}
	}

	rel := c.GetRelation("", "t")
	if rel == nil {
		t.Fatal("table t is missing")
	}
	parentID := -1
	for i, col := range rel.Columns {
		if col.Name == "parent_id" {
			parentID = i
		}
	}
	if parentID < 0 {
		t.Fatal("column parent_id is missing")
	}
	// The dependency is the column's, not just the table's: dropping parent_id
	// reaches the view.
	if deps := c.relationDependents(rel.OID, int16(parentID+1)); len(deps) == 0 {
		t.Fatalf("view v_hier_dep is no dependent of t.parent_id, though CONNECT BY reads it")
	}
	if deps := c.relationDependents(rel.OID, 1); len(deps) == 0 {
		t.Fatalf("view v_hier_dep is no dependent of t.id, which CONNECT BY and the SELECT list read")
	}
}

// TestViewLevelPseudocolumn pins LEVEL, Redshift's integer pseudocolumn of a
// hierarchical query: a view may name it, its conditions included, and the
// stored definition keeps the word.
func TestViewLevelPseudocolumn(t *testing.T) {
	c := New()
	for _, stmt := range parseStmts(t, `CREATE TABLE t (id int, parent_id int);`) {
		if err := c.ProcessUtility(stmt); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}
	for _, stmt := range parseStmts(t, `CREATE VIEW v_level AS SELECT id, level FROM t WHERE level < 4 START WITH parent_id IS NULL CONNECT BY PRIOR id = parent_id;`) {
		if err := c.ProcessUtility(stmt); err != nil {
			t.Fatalf("create view: %v", err)
		}
	}

	rel := c.GetRelation("", "v_level")
	if rel == nil {
		t.Fatal("view v_level is missing")
	}
	if len(rel.Columns) != 2 || rel.Columns[1].Name != "level" {
		t.Fatalf("columns = %#v, want id and level", rel.Columns)
	}
	def, err := c.GetViewDefinition("", "v_level")
	if err != nil {
		t.Fatalf("GetViewDefinition: %v", err)
	}
	if !strings.Contains(def, "LEVEL") {
		t.Fatalf("view definition = %q, want LEVEL in it", def)
	}
}

// TestMigrationMatViewReplaceKeepsIndexesAndDependents pins the plan for a
// materialized view the catalog replaces: its indexes are created again after
// the new view, and a view that reads it is dropped first and created again
// after it.
func TestMigrationMatViewReplaceKeepsIndexesAndWarns(t *testing.T) {
	from, err := LoadSQL(`CREATE TABLE t (id int);
CREATE MATERIALIZED VIEW mv BACKUP NO AS SELECT id FROM t;
CREATE INDEX idx_mv_id ON mv (id);
CREATE VIEW v_dep AS SELECT id FROM mv;`)
	if err != nil {
		t.Fatal(err)
	}
	to, err := LoadSQL(`CREATE TABLE t (id int);
CREATE MATERIALIZED VIEW mv BACKUP YES AS SELECT id FROM t;
CREATE INDEX idx_mv_id ON mv (id);
CREATE VIEW v_dep AS SELECT id FROM mv;`)
	if err != nil {
		t.Fatal(err)
	}

	plan := GenerateMigration(from, to, Diff(from, to))
	createIdx := 0
	for _, op := range plan.Ops {
		if op.Type == OpCreateIndex && strings.Contains(op.SQL, "idx_mv_id") {
			createIdx++
		}
	}
	if createIdx != 1 {
		t.Fatalf("CREATE INDEX idx_mv_id appears %d times, want once after the re-created view", createIdx)
	}
	// The dependent view is dropped before the view it reads and created again
	// after it, so the drop the plan emits can run.
	var dropDep, dropMV, createMV, createDep = -1, -1, -1, -1
	for i, op := range plan.Ops {
		switch {
		case op.Type == OpDropView && strings.Contains(op.SQL, "DROP VIEW") && strings.Contains(op.SQL, `"v_dep"`):
			dropDep = i
		case op.Type == OpDropView && strings.Contains(op.SQL, "DROP MATERIALIZED VIEW"):
			dropMV = i
		case op.Type == OpCreateView && strings.Contains(op.SQL, "CREATE MATERIALIZED VIEW"):
			createMV = i
		case op.Type == OpCreateView && strings.Contains(op.SQL, "CREATE VIEW") && strings.Contains(op.SQL, `"v_dep"`):
			createDep = i
		}
	}
	if dropDep < 0 || dropMV < 0 || createMV < 0 || createDep < 0 {
		t.Fatalf("plan = %v, want drops and creates of both views", plan.Ops)
	}
	if !(dropDep < dropMV && dropMV < createMV && createMV < createDep) {
		t.Fatalf("order = drop v_dep %d, drop mv %d, create mv %d, create v_dep %d, want them in that order", dropDep, dropMV, createMV, createDep)
	}
}

// TestCloneKeepsMatViewOptionsApart pins the deep copy of a materialized view's
// storage: changing the clone's option leaves the source's alone, so a Diff of
// the two sees the change.
func TestCloneKeepsMatViewOptionsApart(t *testing.T) {
	c := New()
	for _, stmt := range parseStmts(t, `CREATE TABLE t (id int);
CREATE MATERIALIZED VIEW mv BACKUP NO AS SELECT id FROM t;`) {
		if err := c.ProcessUtility(stmt); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}
	clone := c.Clone()
	clonedRel := clone.GetRelation("", "mv")
	if clonedRel == nil || clonedRel.MatViewOptions == nil {
		t.Fatal("the clone holds no materialized view options")
	}
	backup := optionNamed(clonedRel.MatViewOptions, "backup")
	if backup == nil {
		t.Fatalf("options = %#v, want a backup", clonedRel.MatViewOptions)
	}
	backup.Arg.(*nodes.String).Str = "yes"

	source := optionNamed(c.GetRelation("", "mv").MatViewOptions, "backup")
	if source.Arg.(*nodes.String).Str != "no" {
		t.Fatalf("the source's backup followed the clone's")
	}
	if deparseMatViewOptions(clonedRel) == deparseMatViewOptions(c.GetRelation("", "mv")) {
		t.Fatalf("the two catalogs' storage deparses the same after the clone changed")
	}
}

// optionNamed returns the DefElem a materialized view's options hold under name.
func optionNamed(options *nodes.List, name string) *nodes.DefElem {
	if options == nil {
		return nil
	}
	for _, item := range options.Items {
		if elem, ok := item.(*nodes.DefElem); ok && elem.Defname == name {
			return elem
		}
	}
	return nil
}
