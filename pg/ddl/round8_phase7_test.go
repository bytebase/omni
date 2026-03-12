package ddl

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// =============================================================================
// Round 8 Phase 7: ALTER TABLE Expansion
// =============================================================================

// -----------------------------------------------------------------------------
// Gap 1: Multiple ALTER TYPE on the same column in one statement
// pg: src/backend/commands/tablecmds.c:13197-13203
// -----------------------------------------------------------------------------

func TestAlterTypeSameColumnTwiceError(t *testing.T) {
	c := New()
	if err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r'); err != nil {
		t.Fatal(err)
	}

	// ALTER TABLE t ALTER COLUMN a TYPE int8, ALTER COLUMN a TYPE text — same column twice.
	stmt := makeAlterTableStmt("", "t",
		makeATAlterColumnType("a", TypeName{Name: "int8", TypeMod: -1}),
		makeATAlterColumnType("a", TypeName{Name: "text", TypeMod: -1}),
	)
	err := c.AlterTableStmt(stmt)
	assertErrorCode(t, err, CodeInvalidObjectDefinition)
	if !strings.Contains(err.Error(), "cannot alter type of column \"a\" twice") {
		t.Errorf("unexpected message: %s", err)
	}
}

func TestAlterTypeDifferentColumnsOK(t *testing.T) {
	c := New()
	if err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "b", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r'); err != nil {
		t.Fatal(err)
	}

	// ALTER TABLE t ALTER COLUMN a TYPE int8, ALTER COLUMN b TYPE int8 — different columns, OK.
	stmt := makeAlterTableStmt("", "t",
		makeATAlterColumnType("a", TypeName{Name: "int8", TypeMod: -1}),
		makeATAlterColumnType("b", TypeName{Name: "int8", TypeMod: -1}),
	)
	err := c.AlterTableStmt(stmt)
	if err != nil {
		t.Fatal(err)
	}

	r := c.GetRelation("", "t")
	if r.Columns[0].TypeOID != INT8OID {
		t.Errorf("column a type: got %d, want %d", r.Columns[0].TypeOID, INT8OID)
	}
	if r.Columns[1].TypeOID != INT8OID {
		t.Errorf("column b type: got %d, want %d", r.Columns[1].TypeOID, INT8OID)
	}
}

// -----------------------------------------------------------------------------
// Gap 2: ALTER TYPE with existing default — coercion check
// -----------------------------------------------------------------------------

func TestAlterTypeWithExistingDefault(t *testing.T) {
	c := New()
	if err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}, Default: "42"},
	}, nil, false), 'r'); err != nil {
		t.Fatal(err)
	}

	// int4 -> int8 is a valid coercion, should succeed even with default.
	stmt := makeAlterTableStmt("", "t",
		makeATAlterColumnType("a", TypeName{Name: "int8", TypeMod: -1}),
	)
	err := c.AlterTableStmt(stmt)
	if err != nil {
		t.Fatal(err)
	}

	r := c.GetRelation("", "t")
	if r.Columns[0].TypeOID != INT8OID {
		t.Errorf("column type: got %d, want %d", r.Columns[0].TypeOID, INT8OID)
	}
}

func TestAlterTypeWithIncompatibleDefault(t *testing.T) {
	c := New()
	if err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "ts", Type: TypeName{Name: "timestamp", TypeMod: -1}, Default: "now()"},
	}, nil, false), 'r'); err != nil {
		t.Fatal(err)
	}

	// timestamp -> int4 is not a valid coercion.
	stmt := makeAlterTableStmt("", "t",
		makeATAlterColumnType("ts", TypeName{Name: "int4", TypeMod: -1}),
	)
	err := c.AlterTableStmt(stmt)
	assertErrorCode(t, err, CodeDatatypeMismatch)
}

// -----------------------------------------------------------------------------
// Gap 4: Recursive child handling for inheritance
// pg: src/backend/commands/tablecmds.c — ATSimpleRecursion
// -----------------------------------------------------------------------------

func TestRecursiveSetNotNull(t *testing.T) {
	c := New()
	// Create parent with nullable column.
	if err := c.DefineRelation(makeCreateTableStmt("", "parent", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r'); err != nil {
		t.Fatal(err)
	}
	// Create child inheriting from parent.
	childStmt := makeCreateTableStmt("", "child", nil, nil, false)
	childStmt.InhRelations = &nodes.List{Items: []nodes.Node{
		&nodes.RangeVar{Relname: "parent"},
	}}
	if err := c.DefineRelation(childStmt, 'r'); err != nil {
		t.Fatal(err)
	}

	// SET NOT NULL on parent should propagate to child.
	err := c.AlterTableStmt(makeAlterTableStmt("", "parent",
		makeATSetNotNull("a"),
	))
	if err != nil {
		t.Fatal(err)
	}

	parentRel := c.GetRelation("", "parent")
	if !parentRel.Columns[0].NotNull {
		t.Error("parent column 'a' should be NOT NULL")
	}
	childRel := c.GetRelation("", "child")
	if !childRel.Columns[0].NotNull {
		t.Error("child column 'a' should be NOT NULL (propagated from parent)")
	}
}

func TestRecursiveDropNotNull(t *testing.T) {
	c := New()
	// Create parent with NOT NULL column.
	if err := c.DefineRelation(makeCreateTableStmt("", "parent", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}, NotNull: true},
	}, nil, false), 'r'); err != nil {
		t.Fatal(err)
	}
	// Create child inheriting from parent.
	childStmt := makeCreateTableStmt("", "child", nil, nil, false)
	childStmt.InhRelations = &nodes.List{Items: []nodes.Node{
		&nodes.RangeVar{Relname: "parent"},
	}}
	if err := c.DefineRelation(childStmt, 'r'); err != nil {
		t.Fatal(err)
	}

	// DROP NOT NULL on parent should propagate to child.
	err := c.AlterTableStmt(makeAlterTableStmt("", "parent",
		makeATDropNotNull("a"),
	))
	if err != nil {
		t.Fatal(err)
	}

	parentRel := c.GetRelation("", "parent")
	if parentRel.Columns[0].NotNull {
		t.Error("parent column 'a' should not be NOT NULL")
	}
	childRel := c.GetRelation("", "child")
	if childRel.Columns[0].NotNull {
		t.Error("child column 'a' should not be NOT NULL (propagated from parent)")
	}
}

func TestRecursiveAddColumn(t *testing.T) {
	c := New()
	// Create parent.
	if err := c.DefineRelation(makeCreateTableStmt("", "parent", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r'); err != nil {
		t.Fatal(err)
	}
	// Create child inheriting from parent.
	childStmt := makeCreateTableStmt("", "child", nil, nil, false)
	childStmt.InhRelations = &nodes.List{Items: []nodes.Node{
		&nodes.RangeVar{Relname: "parent"},
	}}
	if err := c.DefineRelation(childStmt, 'r'); err != nil {
		t.Fatal(err)
	}

	// ADD COLUMN on parent should propagate to child.
	err := c.AlterTableStmt(makeAlterTableStmt("", "parent",
		makeATAddColumn(ColumnDef{Name: "b", Type: TypeName{Name: "text", TypeMod: -1}}),
	))
	if err != nil {
		t.Fatal(err)
	}

	parentRel := c.GetRelation("", "parent")
	if len(parentRel.Columns) != 2 {
		t.Fatalf("parent columns: got %d, want 2", len(parentRel.Columns))
	}
	childRel := c.GetRelation("", "child")
	if len(childRel.Columns) != 2 {
		t.Fatalf("child columns: got %d, want 2 (propagated from parent)", len(childRel.Columns))
	}
	if childRel.Columns[1].Name != "b" {
		t.Errorf("child column name: got %q, want %q", childRel.Columns[1].Name, "b")
	}
}

func TestRecursiveDropColumn(t *testing.T) {
	c := New()
	// Create parent with two columns.
	if err := c.DefineRelation(makeCreateTableStmt("", "parent", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "b", Type: TypeName{Name: "text", TypeMod: -1}},
	}, nil, false), 'r'); err != nil {
		t.Fatal(err)
	}
	// Create child inheriting from parent.
	childStmt := makeCreateTableStmt("", "child", nil, nil, false)
	childStmt.InhRelations = &nodes.List{Items: []nodes.Node{
		&nodes.RangeVar{Relname: "parent"},
	}}
	if err := c.DefineRelation(childStmt, 'r'); err != nil {
		t.Fatal(err)
	}

	// Verify child has both inherited columns.
	childRel := c.GetRelation("", "child")
	if len(childRel.Columns) != 2 {
		t.Fatalf("child columns before drop: got %d, want 2", len(childRel.Columns))
	}

	// DROP COLUMN b CASCADE on parent should propagate to child.
	// CASCADE is needed because the child relation depends on the parent.
	err := c.AlterTableStmt(makeAlterTableStmt("", "parent",
		makeATDropColumn("b", true, false),
	))
	if err != nil {
		t.Fatal(err)
	}

	parentRel := c.GetRelation("", "parent")
	if len(parentRel.Columns) != 1 {
		t.Fatalf("parent columns after drop: got %d, want 1", len(parentRel.Columns))
	}
}

// -----------------------------------------------------------------------------
// Gap 5: Expanded pass ordering
// pg: src/backend/commands/tablecmds.c — AlterTablePass
// -----------------------------------------------------------------------------

func TestExpandedPassOrdering(t *testing.T) {
	c := New()
	if err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "b", Type: TypeName{Name: "text", TypeMod: -1}},
	}, nil, false), 'r'); err != nil {
		t.Fatal(err)
	}

	// Issue commands in reverse order: ADD CONSTRAINT, ALTER TYPE, DROP.
	// Expected execution order: DROP (pass 0), ALTER TYPE (pass 1), ADD CONSTRAINT (pass 5).
	// The DROP removes column 'b', then ALTER TYPE changes 'a' to int8,
	// then ADD CONSTRAINT adds a CHECK on column 'a'.
	stmt := makeAlterTableStmt("", "t",
		makeATAddConstraint(ConstraintDef{Type: ConstraintCheck, Name: "chk", CheckExpr: "a > 0"}),
		makeATAlterColumnType("a", TypeName{Name: "int8", TypeMod: -1}),
		makeATDropColumn("b", false, false),
	)
	err := c.AlterTableStmt(stmt)
	if err != nil {
		t.Fatal(err)
	}

	r := c.GetRelation("", "t")
	if len(r.Columns) != 1 {
		t.Fatalf("columns: got %d, want 1", len(r.Columns))
	}
	if r.Columns[0].TypeOID != INT8OID {
		t.Errorf("column type: got %d, want %d", r.Columns[0].TypeOID, INT8OID)
	}
	cons := c.ConstraintsOf(r.OID)
	found := false
	for _, con := range cons {
		if con.Name == "chk" {
			found = true
		}
	}
	if !found {
		t.Error("CHECK constraint 'chk' not found after pass-ordered execution")
	}
}

func TestPassOrderDropBeforeAdd(t *testing.T) {
	c := New()
	if err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "b", Type: TypeName{Name: "text", TypeMod: -1}},
	}, nil, false), 'r'); err != nil {
		t.Fatal(err)
	}

	// ADD COLUMN 'c', then DROP COLUMN 'b' in statement order.
	// Pass ordering: DROP (pass 0) before ADD (pass 3).
	stmt := makeAlterTableStmt("", "t",
		makeATAddColumn(ColumnDef{Name: "c", Type: TypeName{Name: "int8", TypeMod: -1}}),
		makeATDropColumn("b", false, false),
	)
	err := c.AlterTableStmt(stmt)
	if err != nil {
		t.Fatal(err)
	}

	r := c.GetRelation("", "t")
	if len(r.Columns) != 2 {
		t.Fatalf("columns: got %d, want 2", len(r.Columns))
	}
	if r.Columns[0].Name != "a" || r.Columns[1].Name != "c" {
		t.Errorf("columns: got %q, %q — want a, c", r.Columns[0].Name, r.Columns[1].Name)
	}
}

func TestPassOrderTypeBeforeConstraint(t *testing.T) {
	c := New()
	if err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r'); err != nil {
		t.Fatal(err)
	}

	// ADD CONSTRAINT first in statement, ALTER TYPE second.
	// Pass ordering: ALTER TYPE (pass 1) before ADD CONSTRAINT (pass 5).
	stmt := makeAlterTableStmt("", "t",
		makeATAddConstraint(ConstraintDef{Type: ConstraintCheck, Name: "chk", CheckExpr: "a > 0"}),
		makeATAlterColumnType("a", TypeName{Name: "int8", TypeMod: -1}),
	)
	err := c.AlterTableStmt(stmt)
	if err != nil {
		t.Fatal(err)
	}

	r := c.GetRelation("", "t")
	if r.Columns[0].TypeOID != INT8OID {
		t.Errorf("column type: got %d, want %d", r.Columns[0].TypeOID, INT8OID)
	}
}

// -----------------------------------------------------------------------------
// Gap 4: findChildRelations helper
// -----------------------------------------------------------------------------

func TestFindChildRelations(t *testing.T) {
	c := New()
	// Create parent.
	if err := c.DefineRelation(makeCreateTableStmt("", "parent", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r'); err != nil {
		t.Fatal(err)
	}
	parentRel := c.GetRelation("", "parent")

	// Create two children inheriting from parent.
	child1Stmt := makeCreateTableStmt("", "child1", nil, nil, false)
	child1Stmt.InhRelations = &nodes.List{Items: []nodes.Node{
		&nodes.RangeVar{Relname: "parent"},
	}}
	if err := c.DefineRelation(child1Stmt, 'r'); err != nil {
		t.Fatal(err)
	}
	child2Stmt := makeCreateTableStmt("", "child2", nil, nil, false)
	child2Stmt.InhRelations = &nodes.List{Items: []nodes.Node{
		&nodes.RangeVar{Relname: "parent"},
	}}
	if err := c.DefineRelation(child2Stmt, 'r'); err != nil {
		t.Fatal(err)
	}

	children := c.findChildRelations(parentRel.OID)
	if len(children) != 2 {
		t.Fatalf("children: got %d, want 2", len(children))
	}

	child1 := c.GetRelation("", "child1")
	child2 := c.GetRelation("", "child2")
	childOIDs := map[uint32]bool{child1.OID: true, child2.OID: true}
	for _, oid := range children {
		if !childOIDs[oid] {
			t.Errorf("unexpected child OID: %d", oid)
		}
	}
}

func TestAlterTablePartitionedRecurse(t *testing.T) {
	c := New()
	// Create partitioned table.
	pt := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "pt"},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}}),
			makeColumnDefNode(ColumnDef{Name: "b", Type: TypeName{Name: "text", TypeMod: -1}}),
		}},
		Partspec: &nodes.PartitionSpec{
			Strategy: "range",
			PartParams: &nodes.List{Items: []nodes.Node{
				&nodes.PartitionElem{Name: "a"},
			}},
		},
	}
	if err := c.DefineRelation(pt, 'r'); err != nil {
		t.Fatal(err)
	}
	// Create partition child.
	partChild := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "pt_child"},
		Partbound: &nodes.PartitionBoundSpec{
			Strategy:    'r',
			Lowerdatums: &nodes.List{Items: []nodes.Node{&nodes.Integer{Ival: 1}}},
			Upperdatums: &nodes.List{Items: []nodes.Node{&nodes.Integer{Ival: 100}}},
		},
		InhRelations: &nodes.List{Items: []nodes.Node{
			&nodes.RangeVar{Relname: "pt"},
		}},
	}
	if err := c.DefineRelation(partChild, 'r'); err != nil {
		t.Fatal(err)
	}

	// SET NOT NULL on partitioned table should propagate to partition child.
	err := c.AlterTableStmt(makeAlterTableStmt("", "pt",
		makeATSetNotNull("b"),
	))
	if err != nil {
		t.Fatal(err)
	}

	ptRel := c.GetRelation("", "pt")
	if !ptRel.Columns[1].NotNull {
		t.Error("partitioned table column 'b' should be NOT NULL")
	}
	childRel := c.GetRelation("", "pt_child")
	if !childRel.Columns[1].NotNull {
		t.Error("partition child column 'b' should be NOT NULL (propagated)")
	}
}

func TestAlterTableNoRecurseForNonApplicable(t *testing.T) {
	c := New()
	// Create parent.
	if err := c.DefineRelation(makeCreateTableStmt("", "parent", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r'); err != nil {
		t.Fatal(err)
	}
	// Create child inheriting from parent.
	childStmt := makeCreateTableStmt("", "child", nil, nil, false)
	childStmt.InhRelations = &nodes.List{Items: []nodes.Node{
		&nodes.RangeVar{Relname: "parent"},
	}}
	if err := c.DefineRelation(childStmt, 'r'); err != nil {
		t.Fatal(err)
	}

	// SET STORAGE is a no-op and should not recurse.
	// Just verify it doesn't error.
	stmt := makeAlterTableStmt("", "parent",
		&nodes.AlterTableCmd{Subtype: int(nodes.AT_SetStorage), Name: "a"},
	)
	err := c.AlterTableStmt(stmt)
	if err != nil {
		t.Fatal(err)
	}

	// Verify shouldRecurse returns false for AT_SetStorage.
	if shouldRecurse(nodes.AT_SetStorage) {
		t.Error("AT_SetStorage should not recurse")
	}
}

// -----------------------------------------------------------------------------
// Gap 6: AT_SetNotNull on primary key columns (verification)
// -----------------------------------------------------------------------------

func TestAlterTableInheritedColumnNotNull(t *testing.T) {
	c := New()
	// Create parent with NOT NULL column.
	if err := c.DefineRelation(makeCreateTableStmt("", "parent", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}, NotNull: true},
	}, nil, false), 'r'); err != nil {
		t.Fatal(err)
	}
	// Create child inheriting from parent.
	childStmt := makeCreateTableStmt("", "child", nil, nil, false)
	childStmt.InhRelations = &nodes.List{Items: []nodes.Node{
		&nodes.RangeVar{Relname: "parent"},
	}}
	if err := c.DefineRelation(childStmt, 'r'); err != nil {
		t.Fatal(err)
	}

	// The inherited column should already have NOT NULL.
	childRel := c.GetRelation("", "child")
	if !childRel.Columns[0].NotNull {
		t.Error("inherited column should be NOT NULL from parent")
	}

	// SET NOT NULL on already NOT NULL inherited column should be idempotent.
	err := c.AlterTableStmt(makeAlterTableStmt("", "child",
		makeATSetNotNull("id"),
	))
	if err != nil {
		t.Fatal(err)
	}
	if !childRel.Columns[0].NotNull {
		t.Error("column should still be NOT NULL")
	}
}
