package ddl

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// --- Phase 2 ALTER TABLE sub-command helpers ---

// makeCompositeTypeStmt builds a CompositeTypeStmt for DefineCompositeType.
func makeCompositeTypeStmt(schema, name string, cols []ColumnDef) *nodes.CompositeTypeStmt {
	var elts []nodes.Node
	for _, cd := range cols {
		elts = append(elts, makeColumnDefNode(cd))
	}
	stmt := &nodes.CompositeTypeStmt{
		Typevar: &nodes.RangeVar{Schemaname: schema, Relname: name},
	}
	if len(elts) > 0 {
		stmt.Coldeflist = &nodes.List{Items: elts}
	}
	return stmt
}

// setupTriggerOnTable creates a trigger function and a trigger on a table.
// The table must already exist. Returns the trigger name.
func setupTriggerOnTable(t *testing.T, c *Catalog, tableName, trigName string) {
	t.Helper()
	err := c.CreateFunctionStmt(makeCreateFuncStmt("", "trig_fn",
		nil, TypeName{Name: "trigger", TypeMod: -1},
		"plpgsql", "BEGIN RETURN NEW; END;", false))
	if err != nil {
		t.Fatalf("create trigger function: %v", err)
	}
	err = c.CreateTriggerStmt(makeCreateTrigStmt("", tableName, trigName, TriggerBefore, TriggerEventInsert, true, "", "trig_fn"))
	if err != nil {
		t.Fatalf("create trigger: %v", err)
	}
}

// --- AT_AddInherit / AT_DropInherit ---

func TestATAddInherit(t *testing.T) {
	c := New()

	// Create parent with columns (a int4, b text).
	err := c.DefineRelation(makeCreateTableStmt("", "parent", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "b", Type: TypeName{Name: "text", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	// Create child with column (a int4) only.
	err = c.DefineRelation(makeCreateTableStmt("", "child", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	// ALTER TABLE child INHERIT parent.
	cmd := &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_AddInherit),
		Def:     &nodes.RangeVar{Relname: "parent"},
	}
	err = c.AlterTableStmt(makeAlterTableStmt("", "child", cmd))
	if err != nil {
		t.Fatal(err)
	}

	child := c.GetRelation("", "child")
	if child == nil {
		t.Fatal("child relation not found")
	}

	// Child should now have 2 columns: "a" (existing, InhCount incremented) and "b" (added).
	if len(child.Columns) != 2 {
		t.Fatalf("child columns: got %d, want 2", len(child.Columns))
	}
	if child.Columns[0].Name != "a" {
		t.Errorf("child column 0: got %q, want %q", child.Columns[0].Name, "a")
	}
	if child.Columns[0].InhCount != 1 {
		t.Errorf("child column 'a' InhCount: got %d, want 1", child.Columns[0].InhCount)
	}
	if child.Columns[1].Name != "b" {
		t.Errorf("child column 1: got %q, want %q", child.Columns[1].Name, "b")
	}
	if child.Columns[1].InhCount != 1 {
		t.Errorf("child column 'b' InhCount: got %d, want 1", child.Columns[1].InhCount)
	}
	if child.Columns[1].IsLocal {
		t.Error("inherited column 'b' should not be IsLocal")
	}

	// InhCount on the relation should be incremented.
	if child.InhCount != 1 {
		t.Errorf("child InhCount: got %d, want 1", child.InhCount)
	}
	if len(child.InhParents) != 1 {
		t.Fatalf("child InhParents: got %d, want 1", len(child.InhParents))
	}
	parent := c.GetRelation("", "parent")
	if child.InhParents[0] != parent.OID {
		t.Errorf("child InhParents[0]: got %d, want %d", child.InhParents[0], parent.OID)
	}
}

func TestATAddInheritTypeMismatch(t *testing.T) {
	c := New()

	err := c.DefineRelation(makeCreateTableStmt("", "parent", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "text", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	err = c.DefineRelation(makeCreateTableStmt("", "child", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	cmd := &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_AddInherit),
		Def:     &nodes.RangeVar{Relname: "parent"},
	}
	err = c.AlterTableStmt(makeAlterTableStmt("", "child", cmd))
	assertErrorCode(t, err, CodeDatatypeMismatch)
}

func TestATDropInherit(t *testing.T) {
	c := New()

	err := c.DefineRelation(makeCreateTableStmt("", "parent", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "b", Type: TypeName{Name: "text", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	err = c.DefineRelation(makeCreateTableStmt("", "child", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	// First, INHERIT.
	addCmd := &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_AddInherit),
		Def:     &nodes.RangeVar{Relname: "parent"},
	}
	err = c.AlterTableStmt(makeAlterTableStmt("", "child", addCmd))
	if err != nil {
		t.Fatal(err)
	}

	child := c.GetRelation("", "child")
	if child.InhCount != 1 {
		t.Fatalf("pre-condition: InhCount=%d, want 1", child.InhCount)
	}

	// Now, DROP INHERIT.
	dropCmd := &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_DropInherit),
		Def:     &nodes.RangeVar{Relname: "parent"},
	}
	err = c.AlterTableStmt(makeAlterTableStmt("", "child", dropCmd))
	if err != nil {
		t.Fatal(err)
	}

	if child.InhCount != 0 {
		t.Errorf("child InhCount after drop: got %d, want 0", child.InhCount)
	}
	if len(child.InhParents) != 0 {
		t.Errorf("child InhParents after drop: got %d, want 0", len(child.InhParents))
	}

	// Column "a" should have InhCount decremented and IsLocal set to true.
	colA := child.Columns[0]
	if colA.InhCount != 0 {
		t.Errorf("column 'a' InhCount after drop: got %d, want 0", colA.InhCount)
	}
	if !colA.IsLocal {
		t.Error("column 'a' should be IsLocal after dropping inheritance")
	}

	// Column "b" (was inherited) should also have InhCount=0 and IsLocal=true.
	colB := child.Columns[1]
	if colB.InhCount != 0 {
		t.Errorf("column 'b' InhCount after drop: got %d, want 0", colB.InhCount)
	}
	if !colB.IsLocal {
		t.Error("column 'b' should be IsLocal after dropping inheritance")
	}
}

func TestATDropInheritNotAParent(t *testing.T) {
	c := New()

	err := c.DefineRelation(makeCreateTableStmt("", "t1", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	err = c.DefineRelation(makeCreateTableStmt("", "t2", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	cmd := &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_DropInherit),
		Def:     &nodes.RangeVar{Relname: "t1"},
	}
	err = c.AlterTableStmt(makeAlterTableStmt("", "t2", cmd))
	assertErrorCode(t, err, CodeInvalidObjectDefinition)
}

// --- AT_SetLogged / AT_SetUnLogged ---

func TestATSetLoggedUnLogged(t *testing.T) {
	c := New()
	err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	// Default persistence should be 'p' (permanent).
	if rel.Persistence != 'p' {
		t.Fatalf("initial persistence: got %c, want p", rel.Persistence)
	}

	// SET UNLOGGED.
	cmd := &nodes.AlterTableCmd{Subtype: int(nodes.AT_SetUnLogged)}
	err = c.AlterTableStmt(makeAlterTableStmt("", "t", cmd))
	if err != nil {
		t.Fatal(err)
	}
	if rel.Persistence != 'u' {
		t.Errorf("after SET UNLOGGED: got %c, want u", rel.Persistence)
	}

	// SET LOGGED.
	cmd = &nodes.AlterTableCmd{Subtype: int(nodes.AT_SetLogged)}
	err = c.AlterTableStmt(makeAlterTableStmt("", "t", cmd))
	if err != nil {
		t.Fatal(err)
	}
	if rel.Persistence != 'p' {
		t.Errorf("after SET LOGGED: got %c, want p", rel.Persistence)
	}
}

// --- AT_ClusterOn / AT_DropCluster ---

func TestATClusterOn(t *testing.T) {
	c := New()
	err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}
	err = c.DefineIndex(makeIndexStmt("", "t", "myidx", []string{"a"}, false, false))
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")

	// CLUSTER ON myidx.
	cmd := &nodes.AlterTableCmd{Subtype: int(nodes.AT_ClusterOn), Name: "myidx"}
	err = c.AlterTableStmt(makeAlterTableStmt("", "t", cmd))
	if err != nil {
		t.Fatal(err)
	}

	indexes := c.IndexesOf(rel.OID)
	if len(indexes) != 1 {
		t.Fatalf("indexes: got %d, want 1", len(indexes))
	}
	if !indexes[0].IsClustered {
		t.Error("index should be clustered after CLUSTER ON")
	}
}

func TestATClusterOnNonexistentIndex(t *testing.T) {
	c := New()
	err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	cmd := &nodes.AlterTableCmd{Subtype: int(nodes.AT_ClusterOn), Name: "nosuch"}
	err = c.AlterTableStmt(makeAlterTableStmt("", "t", cmd))
	assertErrorCode(t, err, CodeUndefinedObject)
}

func TestATDropCluster(t *testing.T) {
	c := New()
	err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}
	err = c.DefineIndex(makeIndexStmt("", "t", "myidx", []string{"a"}, false, false))
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")

	// CLUSTER ON first.
	cmd := &nodes.AlterTableCmd{Subtype: int(nodes.AT_ClusterOn), Name: "myidx"}
	err = c.AlterTableStmt(makeAlterTableStmt("", "t", cmd))
	if err != nil {
		t.Fatal(err)
	}

	// DROP CLUSTER.
	dropCmd := &nodes.AlterTableCmd{Subtype: int(nodes.AT_DropCluster)}
	err = c.AlterTableStmt(makeAlterTableStmt("", "t", dropCmd))
	if err != nil {
		t.Fatal(err)
	}

	indexes := c.IndexesOf(rel.OID)
	for _, idx := range indexes {
		if idx.IsClustered {
			t.Errorf("index %q should not be clustered after DROP CLUSTER", idx.Name)
		}
	}
}

func TestATClusterOnSwitchesIndex(t *testing.T) {
	c := New()
	err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "b", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}
	err = c.DefineIndex(makeIndexStmt("", "t", "idx_a", []string{"a"}, false, false))
	if err != nil {
		t.Fatal(err)
	}
	err = c.DefineIndex(makeIndexStmt("", "t", "idx_b", []string{"b"}, false, false))
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")

	// Cluster on idx_a.
	cmd := &nodes.AlterTableCmd{Subtype: int(nodes.AT_ClusterOn), Name: "idx_a"}
	err = c.AlterTableStmt(makeAlterTableStmt("", "t", cmd))
	if err != nil {
		t.Fatal(err)
	}

	// Switch to idx_b.
	cmd = &nodes.AlterTableCmd{Subtype: int(nodes.AT_ClusterOn), Name: "idx_b"}
	err = c.AlterTableStmt(makeAlterTableStmt("", "t", cmd))
	if err != nil {
		t.Fatal(err)
	}

	// idx_a should no longer be clustered, idx_b should be.
	indexes := c.IndexesOf(rel.OID)
	for _, idx := range indexes {
		if idx.Name == "idx_a" && idx.IsClustered {
			t.Error("idx_a should not be clustered after switching to idx_b")
		}
		if idx.Name == "idx_b" && !idx.IsClustered {
			t.Error("idx_b should be clustered")
		}
	}
}

// --- AT_SetCompression ---

func TestATSetCompression(t *testing.T) {
	c := New()
	err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "col", Type: TypeName{Name: "text", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	cmd := &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_SetCompression),
		Name:    "col",
		Def:     &nodes.String{Str: "pglz"},
	}
	err = c.AlterTableStmt(makeAlterTableStmt("", "t", cmd))
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	if rel.Columns[0].Compression != 'p' {
		t.Errorf("compression: got %c, want 'p'", rel.Columns[0].Compression)
	}
}

func TestATSetCompressionLZ4(t *testing.T) {
	c := New()
	err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "col", Type: TypeName{Name: "text", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	cmd := &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_SetCompression),
		Name:    "col",
		Def:     &nodes.String{Str: "lz4"},
	}
	err = c.AlterTableStmt(makeAlterTableStmt("", "t", cmd))
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	if rel.Columns[0].Compression != 'l' {
		t.Errorf("compression: got %c, want 'l'", rel.Columns[0].Compression)
	}
}

func TestATSetCompressionOnFixedLenColumn(t *testing.T) {
	c := New()
	err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "col", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	cmd := &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_SetCompression),
		Name:    "col",
		Def:     &nodes.String{Str: "pglz"},
	}
	err = c.AlterTableStmt(makeAlterTableStmt("", "t", cmd))
	assertErrorCode(t, err, CodeInvalidObjectDefinition)
}

// --- AT_AddOf / AT_DropOf ---

func TestATAddOf(t *testing.T) {
	c := New()

	// Create composite type with (x int4, y text).
	err := c.DefineCompositeType(makeCompositeTypeStmt("", "mytype", []ColumnDef{
		{Name: "x", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "y", Type: TypeName{Name: "text", TypeMod: -1}},
	}))
	if err != nil {
		t.Fatal(err)
	}

	// Create table with same columns.
	err = c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "x", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "y", Type: TypeName{Name: "text", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	// ALTER TABLE t OF mytype.
	cmd := &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_AddOf),
		Def:     makeTypeNameNode(TypeName{Name: "mytype", TypeMod: -1}),
	}
	err = c.AlterTableStmt(makeAlterTableStmt("", "t", cmd))
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	if rel.OfTypeOID == 0 {
		t.Error("OfTypeOID should be set after OF type")
	}
}

func TestATAddOfColumnMismatch(t *testing.T) {
	c := New()

	err := c.DefineCompositeType(makeCompositeTypeStmt("", "mytype", []ColumnDef{
		{Name: "x", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "y", Type: TypeName{Name: "text", TypeMod: -1}},
	}))
	if err != nil {
		t.Fatal(err)
	}

	// Create table with different columns.
	err = c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	cmd := &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_AddOf),
		Def:     makeTypeNameNode(TypeName{Name: "mytype", TypeMod: -1}),
	}
	err = c.AlterTableStmt(makeAlterTableStmt("", "t", cmd))
	assertErrorCode(t, err, CodeDatatypeMismatch)
}

func TestATDropOf(t *testing.T) {
	c := New()

	err := c.DefineCompositeType(makeCompositeTypeStmt("", "mytype", []ColumnDef{
		{Name: "x", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "y", Type: TypeName{Name: "text", TypeMod: -1}},
	}))
	if err != nil {
		t.Fatal(err)
	}

	err = c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "x", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "y", Type: TypeName{Name: "text", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	// Add OF type.
	addCmd := &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_AddOf),
		Def:     makeTypeNameNode(TypeName{Name: "mytype", TypeMod: -1}),
	}
	err = c.AlterTableStmt(makeAlterTableStmt("", "t", addCmd))
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	if rel.OfTypeOID == 0 {
		t.Fatal("pre-condition: OfTypeOID should be set")
	}

	// NOT OF (drop of).
	dropCmd := &nodes.AlterTableCmd{Subtype: int(nodes.AT_DropOf)}
	err = c.AlterTableStmt(makeAlterTableStmt("", "t", dropCmd))
	if err != nil {
		t.Fatal(err)
	}

	if rel.OfTypeOID != 0 {
		t.Error("OfTypeOID should be 0 after NOT OF")
	}
}

// --- AT_EnableAlwaysTrig / AT_EnableReplicaTrig ---

func TestATEnableAlwaysTrigger(t *testing.T) {
	c := New()
	err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	setupTriggerOnTable(t, c, "t", "mytrig")

	// ENABLE ALWAYS mytrig.
	cmd := &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_EnableAlwaysTrig),
		Name:    "mytrig",
	}
	err = c.AlterTableStmt(makeAlterTableStmt("", "t", cmd))
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	trigs := c.TriggersOf(rel.OID)
	if len(trigs) != 1 {
		t.Fatalf("triggers: got %d, want 1", len(trigs))
	}
	if trigs[0].Enabled != 'A' {
		t.Errorf("trigger Enabled: got %c, want A", trigs[0].Enabled)
	}
}

func TestATEnableReplicaTrigger(t *testing.T) {
	c := New()
	err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	setupTriggerOnTable(t, c, "t", "mytrig")

	// ENABLE REPLICA mytrig.
	cmd := &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_EnableReplicaTrig),
		Name:    "mytrig",
	}
	err = c.AlterTableStmt(makeAlterTableStmt("", "t", cmd))
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	trigs := c.TriggersOf(rel.OID)
	if len(trigs) != 1 {
		t.Fatalf("triggers: got %d, want 1", len(trigs))
	}
	if trigs[0].Enabled != 'R' {
		t.Errorf("trigger Enabled: got %c, want R", trigs[0].Enabled)
	}
}

func TestATEnableTriggerSequence(t *testing.T) {
	c := New()
	err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	setupTriggerOnTable(t, c, "t", "mytrig")
	rel := c.GetRelation("", "t")

	// ENABLE ALWAYS.
	cmd := &nodes.AlterTableCmd{Subtype: int(nodes.AT_EnableAlwaysTrig), Name: "mytrig"}
	err = c.AlterTableStmt(makeAlterTableStmt("", "t", cmd))
	if err != nil {
		t.Fatal(err)
	}
	trigs := c.TriggersOf(rel.OID)
	if trigs[0].Enabled != 'A' {
		t.Fatalf("after ALWAYS: got %c, want A", trigs[0].Enabled)
	}

	// ENABLE REPLICA.
	cmd = &nodes.AlterTableCmd{Subtype: int(nodes.AT_EnableReplicaTrig), Name: "mytrig"}
	err = c.AlterTableStmt(makeAlterTableStmt("", "t", cmd))
	if err != nil {
		t.Fatal(err)
	}
	trigs = c.TriggersOf(rel.OID)
	if trigs[0].Enabled != 'R' {
		t.Fatalf("after REPLICA: got %c, want R", trigs[0].Enabled)
	}

	// DISABLE.
	cmd = &nodes.AlterTableCmd{Subtype: int(nodes.AT_DisableTrig), Name: "mytrig"}
	err = c.AlterTableStmt(makeAlterTableStmt("", "t", cmd))
	if err != nil {
		t.Fatal(err)
	}
	trigs = c.TriggersOf(rel.OID)
	if trigs[0].Enabled != 'D' {
		t.Fatalf("after DISABLE: got %c, want D", trigs[0].Enabled)
	}

	// ENABLE (back to origin).
	cmd = &nodes.AlterTableCmd{Subtype: int(nodes.AT_EnableTrig), Name: "mytrig"}
	err = c.AlterTableStmt(makeAlterTableStmt("", "t", cmd))
	if err != nil {
		t.Fatal(err)
	}
	trigs = c.TriggersOf(rel.OID)
	if trigs[0].Enabled != 'O' {
		t.Fatalf("after ENABLE: got %c, want O", trigs[0].Enabled)
	}
}

func TestATEnableTriggerNonexistent(t *testing.T) {
	c := New()
	err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	cmd := &nodes.AlterTableCmd{Subtype: int(nodes.AT_EnableAlwaysTrig), Name: "nosuch"}
	err = c.AlterTableStmt(makeAlterTableStmt("", "t", cmd))
	assertErrorCode(t, err, CodeUndefinedObject)
}

// --- AT_DropExpression ---

func TestATDropExpression(t *testing.T) {
	c := New()
	err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "b", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	// Manually set column "b" as a generated column.
	rel := c.GetRelation("", "t")
	colIdx := rel.colByName["b"]
	col := rel.Columns[colIdx]
	col.Generated = 's'
	col.GenerationExpr = "a + 1"
	col.HasDefault = true
	col.Default = "a + 1"

	// DROP EXPRESSION on column "b".
	cmd := &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_DropExpression),
		Name:    "b",
	}
	err = c.AlterTableStmt(makeAlterTableStmt("", "t", cmd))
	if err != nil {
		t.Fatal(err)
	}

	if col.Generated != 0 {
		t.Errorf("Generated: got %c, want 0", col.Generated)
	}
	if col.GenerationExpr != "" {
		t.Errorf("GenerationExpr: got %q, want empty", col.GenerationExpr)
	}
	if col.HasDefault {
		t.Error("HasDefault should be false after DROP EXPRESSION")
	}
	if col.Default != "" {
		t.Errorf("Default: got %q, want empty", col.Default)
	}
}

func TestATDropExpressionNotGenerated(t *testing.T) {
	c := New()
	err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	// Column "a" is not generated.
	cmd := &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_DropExpression),
		Name:    "a",
	}
	err = c.AlterTableStmt(makeAlterTableStmt("", "t", cmd))
	assertErrorCode(t, err, CodeInvalidObjectDefinition)
}

func TestATDropExpressionIfExists(t *testing.T) {
	c := New()
	err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	// Column "a" is not generated, but IF EXISTS should not error.
	cmd := &nodes.AlterTableCmd{
		Subtype:    int(nodes.AT_DropExpression),
		Name:       "a",
		Missing_ok: true,
	}
	err = c.AlterTableStmt(makeAlterTableStmt("", "t", cmd))
	if err != nil {
		t.Fatalf("IF EXISTS should not error, got: %v", err)
	}
}

// --- AT_EnableTrigAll / AT_DisableTrigAll ---

func TestATEnableDisableAllTriggers(t *testing.T) {
	c := New()
	err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	// Create trigger function.
	err = c.CreateFunctionStmt(makeCreateFuncStmt("", "trig_fn",
		nil, TypeName{Name: "trigger", TypeMod: -1},
		"plpgsql", "BEGIN RETURN NEW; END;", false))
	if err != nil {
		t.Fatal(err)
	}

	// Create two triggers.
	err = c.CreateTriggerStmt(makeCreateTrigStmt("", "t", "trig1", TriggerBefore, TriggerEventInsert, true, "", "trig_fn"))
	if err != nil {
		t.Fatal(err)
	}
	err = c.CreateTriggerStmt(makeCreateTrigStmt("", "t", "trig2", TriggerAfter, TriggerEventUpdate, true, "", "trig_fn"))
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")

	// DISABLE ALL.
	cmd := &nodes.AlterTableCmd{Subtype: int(nodes.AT_DisableTrigAll)}
	err = c.AlterTableStmt(makeAlterTableStmt("", "t", cmd))
	if err != nil {
		t.Fatal(err)
	}

	trigs := c.TriggersOf(rel.OID)
	for _, trig := range trigs {
		if trig.Enabled != 'D' {
			t.Errorf("trigger %q: Enabled got %c, want D after DISABLE ALL", trig.Name, trig.Enabled)
		}
	}

	// ENABLE ALL.
	cmd = &nodes.AlterTableCmd{Subtype: int(nodes.AT_EnableTrigAll)}
	err = c.AlterTableStmt(makeAlterTableStmt("", "t", cmd))
	if err != nil {
		t.Fatal(err)
	}

	trigs = c.TriggersOf(rel.OID)
	for _, trig := range trigs {
		if trig.Enabled != 'O' {
			t.Errorf("trigger %q: Enabled got %c, want O after ENABLE ALL", trig.Name, trig.Enabled)
		}
	}
}

// --- AT_EnableRowSecurity / AT_DisableRowSecurity ---

func TestATRowSecurity(t *testing.T) {
	c := New()
	err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	if rel.RowSecurity {
		t.Fatal("row security should be disabled by default")
	}

	// ENABLE ROW LEVEL SECURITY.
	cmd := &nodes.AlterTableCmd{Subtype: int(nodes.AT_EnableRowSecurity)}
	err = c.AlterTableStmt(makeAlterTableStmt("", "t", cmd))
	if err != nil {
		t.Fatal(err)
	}
	if !rel.RowSecurity {
		t.Error("row security should be enabled")
	}

	// DISABLE ROW LEVEL SECURITY.
	cmd = &nodes.AlterTableCmd{Subtype: int(nodes.AT_DisableRowSecurity)}
	err = c.AlterTableStmt(makeAlterTableStmt("", "t", cmd))
	if err != nil {
		t.Fatal(err)
	}
	if rel.RowSecurity {
		t.Error("row security should be disabled")
	}
}

func TestATForceRowSecurity(t *testing.T) {
	c := New()
	err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	if rel.ForceRowSecurity {
		t.Fatal("force row security should be disabled by default")
	}

	// FORCE ROW LEVEL SECURITY.
	cmd := &nodes.AlterTableCmd{Subtype: int(nodes.AT_ForceRowSecurity)}
	err = c.AlterTableStmt(makeAlterTableStmt("", "t", cmd))
	if err != nil {
		t.Fatal(err)
	}
	if !rel.ForceRowSecurity {
		t.Error("force row security should be enabled")
	}

	// NO FORCE ROW LEVEL SECURITY.
	cmd = &nodes.AlterTableCmd{Subtype: int(nodes.AT_NoForceRowSecurity)}
	err = c.AlterTableStmt(makeAlterTableStmt("", "t", cmd))
	if err != nil {
		t.Fatal(err)
	}
	if rel.ForceRowSecurity {
		t.Error("force row security should be disabled")
	}
}
