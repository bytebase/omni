package ddl

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// =============================================================================
// SETOF Column Rejection Tests
// =============================================================================

func TestSetofColumnRejected(t *testing.T) {
	c := New()

	cd := makeColumnDefNode(ColumnDef{Name: "vals", Type: TypeName{Name: "int4", TypeMod: -1}})
	cd.TypeName.Setof = true

	stmt := &nodes.CreateStmt{
		Relation:  &nodes.RangeVar{Relname: "t"},
		TableElts: &nodes.List{Items: []nodes.Node{cd}},
	}
	err := c.DefineRelation(stmt, 'r')
	assertErrorCode(t, err, CodeInvalidTableDefinition)
}

func TestNonSetofColumnOK(t *testing.T) {
	c := New()

	stmt := makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false)
	if err := c.DefineRelation(stmt, 'r'); err != nil {
		t.Fatal(err)
	}
}

// =============================================================================
// Column Count Limit Tests
// =============================================================================

func TestColumnCountAtLimit(t *testing.T) {
	c := New()

	// Build a table with exactly MaxHeapAttributeNumber columns.
	cols := make([]ColumnDef, MaxHeapAttributeNumber)
	for i := range cols {
		cols[i] = ColumnDef{Name: colName(i), Type: TypeName{Name: "int4", TypeMod: -1}}
	}
	stmt := makeCreateTableStmt("", "t", cols, nil, false)
	if err := c.DefineRelation(stmt, 'r'); err != nil {
		t.Fatal(err)
	}
}

func TestColumnCountOverLimit(t *testing.T) {
	c := New()

	cols := make([]ColumnDef, MaxHeapAttributeNumber+1)
	for i := range cols {
		cols[i] = ColumnDef{Name: colName(i), Type: TypeName{Name: "int4", TypeMod: -1}}
	}
	stmt := makeCreateTableStmt("", "t", cols, nil, false)
	err := c.DefineRelation(stmt, 'r')
	assertErrorCode(t, err, CodeTooManyColumns)
}

func colName(i int) string {
	return "c" + string(rune('a'+i/26)) + string(rune('a'+i%26))
}

// =============================================================================
// Persistence Conflict Tests
// =============================================================================

func TestTempPartitionOfPermanentParent(t *testing.T) {
	c := New()

	// Permanent partitioned parent.
	parent := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "parent"},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}),
		}},
		Partspec: &nodes.PartitionSpec{
			Strategy:   "l",
			PartParams: &nodes.List{Items: []nodes.Node{&nodes.PartitionElem{Name: "id"}}},
		},
	}
	if err := c.DefineRelation(parent, 'r'); err != nil {
		t.Fatal(err)
	}

	// Temp partition — should fail.
	child := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "child", Relpersistence: 't'},
		InhRelations: &nodes.List{Items: []nodes.Node{
			&nodes.RangeVar{Relname: "parent", Relpersistence: 't'},
		}},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}),
		}},
		Partbound: &nodes.PartitionBoundSpec{
			Strategy:   'l',
			Listdatums: &nodes.List{Items: []nodes.Node{&nodes.A_Const{Val: &nodes.Integer{Ival: 1}}}},
		},
	}
	err := c.DefineRelation(child, 'r')
	assertErrorCode(t, err, CodeWrongObjectType)
}

func TestPermanentChildOfTempParent(t *testing.T) {
	c := New()

	// Temp parent.
	parent := makeCreateTableStmt("", "tmp_parent", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false)
	parent.Relation.Relpersistence = 't'
	if err := c.DefineRelation(parent, 'r'); err != nil {
		t.Fatal(err)
	}

	// Permanent child inheriting from temp parent — should fail.
	child := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "perm_child"},
		InhRelations: &nodes.List{Items: []nodes.Node{
			&nodes.RangeVar{Relname: "tmp_parent"},
		}},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}),
		}},
	}
	err := c.DefineRelation(child, 'r')
	assertErrorCode(t, err, CodeWrongObjectType)
}

func TestPermanentPartitionOfPermanentOK(t *testing.T) {
	c := New()

	parent := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "pp"},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}),
		}},
		Partspec: &nodes.PartitionSpec{
			Strategy:   "l",
			PartParams: &nodes.List{Items: []nodes.Node{&nodes.PartitionElem{Name: "id"}}},
		},
	}
	if err := c.DefineRelation(parent, 'r'); err != nil {
		t.Fatal(err)
	}

	child := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "pp_child"},
		InhRelations: &nodes.List{Items: []nodes.Node{
			&nodes.RangeVar{Relname: "pp"},
		}},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}),
		}},
		Partbound: &nodes.PartitionBoundSpec{
			Strategy:   'l',
			Listdatums: &nodes.List{Items: []nodes.Node{&nodes.A_Const{Val: &nodes.Integer{Ival: 1}}}},
		},
	}
	if err := c.DefineRelation(child, 'r'); err != nil {
		t.Fatal(err)
	}
}

func TestTempChildOfTempParentOK(t *testing.T) {
	c := New()

	parent := makeCreateTableStmt("", "tp", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false)
	parent.Relation.Relpersistence = 't'
	if err := c.DefineRelation(parent, 'r'); err != nil {
		t.Fatal(err)
	}

	child := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "tc", Relpersistence: 't'},
		InhRelations: &nodes.List{Items: []nodes.Node{
			&nodes.RangeVar{Relname: "tp", Relpersistence: 't'},
		}},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}),
		}},
	}
	if err := c.DefineRelation(child, 'r'); err != nil {
		t.Fatal(err)
	}
}

// =============================================================================
// Identity Inheritance Tests
// =============================================================================

func TestPartitionInheritsIdentity(t *testing.T) {
	c := New()

	// Parent with identity column.
	parent := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "idp"},
		TableElts: &nodes.List{Items: []nodes.Node{
			func() *nodes.ColumnDef {
				cd := makeColumnDefNode(ColumnDef{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}})
				cd.Identity = 'a'
				return cd
			}(),
		}},
		Partspec: &nodes.PartitionSpec{
			Strategy:   "l",
			PartParams: &nodes.List{Items: []nodes.Node{&nodes.PartitionElem{Name: "id"}}},
		},
	}
	if err := c.DefineRelation(parent, 'r'); err != nil {
		t.Fatal(err)
	}

	// Partition should inherit identity.
	child := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "idp_child"},
		InhRelations: &nodes.List{Items: []nodes.Node{
			&nodes.RangeVar{Relname: "idp"},
		}},
		Partbound: &nodes.PartitionBoundSpec{
			Strategy:   'l',
			Listdatums: &nodes.List{Items: []nodes.Node{&nodes.A_Const{Val: &nodes.Integer{Ival: 1}}}},
		},
	}
	if err := c.DefineRelation(child, 'r'); err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "idp_child")
	if rel == nil {
		t.Fatal("partition not found")
	}
	if rel.Columns[0].Identity != 'a' {
		t.Errorf("identity: got %c, want 'a'", rel.Columns[0].Identity)
	}
}

func TestRegularInheritDoesNotInheritIdentity(t *testing.T) {
	c := New()

	// Parent with identity column (non-partitioned).
	parent := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "id_base"},
		TableElts: &nodes.List{Items: []nodes.Node{
			func() *nodes.ColumnDef {
				cd := makeColumnDefNode(ColumnDef{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}})
				cd.Identity = 'a'
				return cd
			}(),
		}},
	}
	if err := c.DefineRelation(parent, 'r'); err != nil {
		t.Fatal(err)
	}

	// Regular inheriting child should NOT get identity.
	child := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "id_child"},
		InhRelations: &nodes.List{Items: []nodes.Node{
			&nodes.RangeVar{Relname: "id_base"},
		}},
	}
	if err := c.DefineRelation(child, 'r'); err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "id_child")
	if rel.Columns[0].Identity != 0 {
		t.Errorf("regular inherit should not copy identity, got %c", rel.Columns[0].Identity)
	}
}

// =============================================================================
// Default Conflict Tests
// =============================================================================

func TestConflictingDefaultsError(t *testing.T) {
	c := New()

	// Two parents with same column but different defaults.
	c.DefineRelation(makeCreateTableStmt("", "p1", []ColumnDef{
		{Name: "val", Type: TypeName{Name: "int4", TypeMod: -1}, Default: "1"},
	}, nil, false), 'r')
	c.DefineRelation(makeCreateTableStmt("", "p2", []ColumnDef{
		{Name: "val", Type: TypeName{Name: "int4", TypeMod: -1}, Default: "2"},
	}, nil, false), 'r')

	// Child inherits from both — conflicting defaults.
	child := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "ch"},
		InhRelations: &nodes.List{Items: []nodes.Node{
			&nodes.RangeVar{Relname: "p1"},
			&nodes.RangeVar{Relname: "p2"},
		}},
	}
	err := c.DefineRelation(child, 'r')
	assertErrorCode(t, err, CodeInvalidColumnDefinition)
}

func TestConflictingDefaultsChildOverrideOK(t *testing.T) {
	c := New()

	c.DefineRelation(makeCreateTableStmt("", "cp1", []ColumnDef{
		{Name: "val", Type: TypeName{Name: "int4", TypeMod: -1}, Default: "1"},
	}, nil, false), 'r')
	c.DefineRelation(makeCreateTableStmt("", "cp2", []ColumnDef{
		{Name: "val", Type: TypeName{Name: "int4", TypeMod: -1}, Default: "2"},
	}, nil, false), 'r')

	// Child locally defines the same column — overrides conflict.
	child := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "cch"},
		InhRelations: &nodes.List{Items: []nodes.Node{
			&nodes.RangeVar{Relname: "cp1"},
			&nodes.RangeVar{Relname: "cp2"},
		}},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "val", Type: TypeName{Name: "int4", TypeMod: -1}, Default: "3"}),
		}},
	}
	if err := c.DefineRelation(child, 'r'); err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "cch")
	if rel.Columns[0].Default != "3" {
		t.Errorf("default: got %q, want %q", rel.Columns[0].Default, "3")
	}
}

func TestSameDefaultsOK(t *testing.T) {
	c := New()

	c.DefineRelation(makeCreateTableStmt("", "sp1", []ColumnDef{
		{Name: "val", Type: TypeName{Name: "int4", TypeMod: -1}, Default: "42"},
	}, nil, false), 'r')
	c.DefineRelation(makeCreateTableStmt("", "sp2", []ColumnDef{
		{Name: "val", Type: TypeName{Name: "int4", TypeMod: -1}, Default: "42"},
	}, nil, false), 'r')

	// Same default from both parents — should not error.
	child := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "sch"},
		InhRelations: &nodes.List{Items: []nodes.Node{
			&nodes.RangeVar{Relname: "sp1"},
			&nodes.RangeVar{Relname: "sp2"},
		}},
	}
	if err := c.DefineRelation(child, 'r'); err != nil {
		t.Fatal(err)
	}
}

// =============================================================================
// OF TYPE (Typed Table) Tests
// =============================================================================

func TestOfTypeBasic(t *testing.T) {
	c := New()

	// Create a composite type.
	c.ProcessUtility(&nodes.CompositeTypeStmt{
		Typevar: &nodes.RangeVar{Relname: "mytype"},
		Coldeflist: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "x", Type: TypeName{Name: "int4", TypeMod: -1}}),
			makeColumnDefNode(ColumnDef{Name: "y", Type: TypeName{Name: "text", TypeMod: -1}}),
		}},
	})

	// Create table OF TYPE.
	stmt := &nodes.CreateStmt{
		Relation:   &nodes.RangeVar{Relname: "typed_tbl"},
		OfTypename: &nodes.TypeName{Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "mytype"}}}},
	}
	if err := c.DefineRelation(stmt, 'r'); err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "typed_tbl")
	if rel == nil {
		t.Fatal("table not found")
	}
	if len(rel.Columns) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(rel.Columns))
	}
	if rel.Columns[0].Name != "x" || rel.Columns[1].Name != "y" {
		t.Errorf("column names: got %q, %q", rel.Columns[0].Name, rel.Columns[1].Name)
	}
	if rel.OfTypeOID == 0 {
		t.Error("expected OfTypeOID to be set")
	}
}

func TestOfTypeNonCompositeError(t *testing.T) {
	c := New()

	// Try to create table OF a non-composite type (e.g. int4).
	stmt := &nodes.CreateStmt{
		Relation:   &nodes.RangeVar{Relname: "bad_typed"},
		OfTypename: &nodes.TypeName{Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "int4"}}}},
	}
	err := c.DefineRelation(stmt, 'r')
	assertErrorCode(t, err, CodeWrongObjectType)
}

func TestOfTypeWithExtraColumns(t *testing.T) {
	c := New()

	// Create composite type.
	c.ProcessUtility(&nodes.CompositeTypeStmt{
		Typevar: &nodes.RangeVar{Relname: "extype"},
		Coldeflist: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}}),
		}},
	})

	// In PG, typed tables only allow column options (NOT NULL, DEFAULT) for
	// existing type columns, not new column definitions. Specifying a column
	// that doesn't exist in the type should fail.
	// pg: src/backend/commands/tablecmds.c — MergeAttributes (line 2516-2527)
	stmt := &nodes.CreateStmt{
		Relation:   &nodes.RangeVar{Relname: "ext_tbl"},
		OfTypename: &nodes.TypeName{Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "extype"}}}},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "b", Type: TypeName{Name: "text", TypeMod: -1}}),
		}},
	}
	err := c.DefineRelation(stmt, 'r')
	assertErrorCode(t, err, CodeUndefinedColumn)

	// Specifying options for an existing type column should succeed.
	stmt2 := &nodes.CreateStmt{
		Relation:   &nodes.RangeVar{Relname: "ext_tbl2"},
		OfTypename: &nodes.TypeName{Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "extype"}}}},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}, NotNull: true}),
		}},
	}
	if err := c.DefineRelation(stmt2, 'r'); err != nil {
		t.Fatal(err)
	}
	rel := c.GetRelation("", "ext_tbl2")
	if len(rel.Columns) != 1 {
		t.Fatalf("expected 1 column, got %d", len(rel.Columns))
	}
	if !rel.Columns[0].NotNull {
		t.Error("expected column a to be NOT NULL after merge")
	}
}

// =============================================================================
// Constraint Fields Tests
// =============================================================================

func TestConstraintConIsLocal(t *testing.T) {
	c := New()

	c.DefineRelation(makeCreateTableStmt("", "cfl", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, []ConstraintDef{
		{Name: "cfl_pk", Type: ConstraintPK, Columns: []string{"id"}},
	}, false), 'r')

	rel := c.GetRelation("", "cfl")
	cons := c.ConstraintsOf(rel.OID)
	if len(cons) != 1 {
		t.Fatalf("expected 1 constraint, got %d", len(cons))
	}
	if !cons[0].ConIsLocal {
		t.Error("ConIsLocal should be true for locally-defined constraint")
	}
}

func TestConstraintQueryPgConstraintFields(t *testing.T) {
	c := New()

	c.DefineRelation(makeCreateTableStmt("", "qpc", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, []ConstraintDef{
		{Name: "qpc_pk", Type: ConstraintPK, Columns: []string{"id"}},
	}, false), 'r')

	rel := c.GetRelation("", "qpc")
	rows := c.QueryPgConstraint(rel.OID)
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	if !rows[0].ConIsLocal {
		t.Error("ConIsLocal should be true")
	}
	if rows[0].ConInhCount != 0 {
		t.Errorf("ConInhCount: got %d, want 0", rows[0].ConInhCount)
	}
}

// =============================================================================
// FK Operator Array Tests
// =============================================================================

func TestFKOperatorArrays(t *testing.T) {
	c := New()

	// Create referenced table with PK.
	c.DefineRelation(makeCreateTableStmt("", "fk_ref", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, []ConstraintDef{
		{Name: "fk_ref_pk", Type: ConstraintPK, Columns: []string{"id"}},
	}, false), 'r')

	// Create table with FK.
	c.DefineRelation(makeCreateTableStmt("", "fk_src", []ColumnDef{
		{Name: "ref_id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, []ConstraintDef{
		{Name: "fk_src_fkey", Type: ConstraintFK, Columns: []string{"ref_id"}, RefTable: "fk_ref"},
	}, false), 'r')

	rel := c.GetRelation("", "fk_src")
	cons := c.ConstraintsOf(rel.OID)
	var fk *Constraint
	for _, cc := range cons {
		if cc.Type == ConstraintFK {
			fk = cc
			break
		}
	}
	if fk == nil {
		t.Fatal("FK constraint not found")
	}
	if len(fk.PFEqOp) != 1 {
		t.Fatalf("PFEqOp: expected 1 element, got %d", len(fk.PFEqOp))
	}
	if fk.PFEqOp[0] == 0 {
		t.Error("PFEqOp[0] should not be 0")
	}
}

func TestFKOperatorArraysInQueryPgConstraint(t *testing.T) {
	c := New()

	c.DefineRelation(makeCreateTableStmt("", "fkq_ref", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, []ConstraintDef{
		{Name: "fkq_ref_pk", Type: ConstraintPK, Columns: []string{"id"}},
	}, false), 'r')

	c.DefineRelation(makeCreateTableStmt("", "fkq_src", []ColumnDef{
		{Name: "ref_id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, []ConstraintDef{
		{Name: "fkq_fkey", Type: ConstraintFK, Columns: []string{"ref_id"}, RefTable: "fkq_ref"},
	}, false), 'r')

	rel := c.GetRelation("", "fkq_src")
	rows := c.QueryPgConstraint(rel.OID)
	var fkRow *PgConstraintRow
	for i := range rows {
		if rows[i].ConType == 'f' {
			fkRow = &rows[i]
			break
		}
	}
	if fkRow == nil {
		t.Fatal("FK row not found")
	}
	if len(fkRow.PFEqOp) == 0 {
		t.Error("PFEqOp should be populated")
	}
	if len(fkRow.PPEqOp) == 0 {
		t.Error("PPEqOp should be populated")
	}
	if len(fkRow.FFEqOp) == 0 {
		t.Error("FFEqOp should be populated")
	}
}

// =============================================================================
// Partition Dependency Types Tests
// =============================================================================

func TestPartitionDepTypes(t *testing.T) {
	c := New()

	// The partition dep type constants exist.
	if DepPartitionPri != 'P' {
		t.Errorf("DepPartitionPri: got %c, want 'P'", DepPartitionPri)
	}
	if DepPartitionSec != 'S' {
		t.Errorf("DepPartitionSec: got %c, want 'S'", DepPartitionSec)
	}

	// Ensure partition creates DepAuto dependency (existing behavior).
	parent := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "dp"},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}),
		}},
		Partspec: &nodes.PartitionSpec{
			Strategy:   "l",
			PartParams: &nodes.List{Items: []nodes.Node{&nodes.PartitionElem{Name: "id"}}},
		},
	}
	c.DefineRelation(parent, 'r')

	child := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "dp_child"},
		InhRelations: &nodes.List{Items: []nodes.Node{
			&nodes.RangeVar{Relname: "dp"},
		}},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}),
		}},
		Partbound: &nodes.PartitionBoundSpec{
			Strategy:   'l',
			Listdatums: &nodes.List{Items: []nodes.Node{&nodes.A_Const{Val: &nodes.Integer{Ival: 1}}}},
		},
	}
	c.DefineRelation(child, 'r')

	parentRel := c.GetRelation("", "dp")
	childRel := c.GetRelation("", "dp_child")
	rows := c.QueryPgDepend()

	found := false
	for _, d := range rows {
		if d.ObjID == childRel.OID && d.RefObjID == parentRel.OID && d.DepType == byte(DepAuto) {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected DepAuto from partition to parent")
	}
}

// =============================================================================
// Policy Validation Tests
// =============================================================================

func TestPolicySelectWithUsingOK(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "pol_t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	stmt := &nodes.CreatePolicyStmt{
		PolicyName: "sel_pol",
		Table:      &nodes.RangeVar{Relname: "pol_t"},
		CmdName:    "select",
		Qual:       &nodes.A_Const{Val: &nodes.String{Str: "true"}},
		Permissive: true,
	}
	if err := c.CreatePolicy(stmt); err != nil {
		t.Fatal(err)
	}
}

func TestPolicySelectWithCheckError(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "pol_t2", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	stmt := &nodes.CreatePolicyStmt{
		PolicyName: "sel_pol2",
		Table:      &nodes.RangeVar{Relname: "pol_t2"},
		CmdName:    "select",
		WithCheck:  &nodes.A_Const{Val: &nodes.String{Str: "true"}},
		Permissive: true,
	}
	err := c.CreatePolicy(stmt)
	assertErrorCode(t, err, CodeSyntaxError)
}

func TestPolicyInsertWithUsingError(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "pol_t3", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	stmt := &nodes.CreatePolicyStmt{
		PolicyName: "ins_pol",
		Table:      &nodes.RangeVar{Relname: "pol_t3"},
		CmdName:    "insert",
		Qual:       &nodes.A_Const{Val: &nodes.String{Str: "true"}},
		Permissive: true,
	}
	err := c.CreatePolicy(stmt)
	assertErrorCode(t, err, CodeSyntaxError)
}

func TestPolicyInsertWithCheckOK(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "pol_t4", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	stmt := &nodes.CreatePolicyStmt{
		PolicyName: "ins_pol2",
		Table:      &nodes.RangeVar{Relname: "pol_t4"},
		CmdName:    "insert",
		WithCheck:  &nodes.A_Const{Val: &nodes.String{Str: "true"}},
		Permissive: true,
	}
	if err := c.CreatePolicy(stmt); err != nil {
		t.Fatal(err)
	}
}

func TestPolicyDeleteWithCheckError(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "pol_t5", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	stmt := &nodes.CreatePolicyStmt{
		PolicyName: "del_pol",
		Table:      &nodes.RangeVar{Relname: "pol_t5"},
		CmdName:    "delete",
		WithCheck:  &nodes.A_Const{Val: &nodes.String{Str: "true"}},
		Permissive: true,
	}
	err := c.CreatePolicy(stmt)
	assertErrorCode(t, err, CodeSyntaxError)
}

func TestPolicyEmptyRolesIsPublic(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "pol_t6", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	// Policy with no roles specified.
	stmt := &nodes.CreatePolicyStmt{
		PolicyName: "pub_pol",
		Table:      &nodes.RangeVar{Relname: "pol_t6"},
		CmdName:    "all",
		Permissive: true,
	}
	if err := c.CreatePolicy(stmt); err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "pol_t6")
	policies := c.policiesByRel[rel.OID]
	if len(policies) != 1 {
		t.Fatalf("expected 1 policy, got %d", len(policies))
	}
	if len(policies[0].Roles) != 1 || policies[0].Roles[0] != "public" {
		t.Errorf("roles: got %v, want [\"public\"]", policies[0].Roles)
	}
}

// =============================================================================
// Comment Enhancement Tests
// =============================================================================

func TestCommentOnMatview(t *testing.T) {
	c := New()

	// Create a matview via DefineRelation with relkind='m'.
	stmt := makeCreateTableStmt("", "mv_c", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false)
	c.DefineRelation(stmt, 'm')

	// Comment on MATERIALIZED VIEW.
	err := c.CommentObject(&nodes.CommentStmt{
		Objtype: nodes.OBJECT_MATVIEW,
		Object:  &nodes.List{Items: []nodes.Node{&nodes.String{Str: "mv_c"}}},
		Comment: "my matview",
	})
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "mv_c")
	comment, ok := c.GetComment('r', rel.OID, 0)
	if !ok || comment != "my matview" {
		t.Errorf("comment: got %q, want %q", comment, "my matview")
	}
}

func TestCommentOnForeignTable(t *testing.T) {
	c := New()

	// Create a foreign table.
	c.ProcessUtility(&nodes.CreateForeignTableStmt{
		Base: nodes.CreateStmt{
			Relation: &nodes.RangeVar{Relname: "ft_c"},
			TableElts: &nodes.List{Items: []nodes.Node{
				makeColumnDefNode(ColumnDef{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}),
			}},
		},
		Servername: "srv",
	})

	err := c.CommentObject(&nodes.CommentStmt{
		Objtype: nodes.OBJECT_FOREIGN_TABLE,
		Object:  &nodes.List{Items: []nodes.Node{&nodes.String{Str: "ft_c"}}},
		Comment: "my foreign table",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCommentOnConstraint(t *testing.T) {
	c := New()

	c.DefineRelation(makeCreateTableStmt("", "cc_t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, []ConstraintDef{
		{Name: "cc_t_pk", Type: ConstraintPK, Columns: []string{"id"}},
	}, false), 'r')

	err := c.CommentObject(&nodes.CommentStmt{
		Objtype: nodes.OBJECT_TABCONSTRAINT,
		Object: &nodes.List{Items: []nodes.Node{
			&nodes.String{Str: "cc_t"},
			&nodes.String{Str: "cc_t_pk"},
		}},
		Comment: "pk constraint",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCommentOnPolicy(t *testing.T) {
	c := New()

	c.DefineRelation(makeCreateTableStmt("", "cp_t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	c.CreatePolicy(&nodes.CreatePolicyStmt{
		PolicyName: "cp_pol",
		Table:      &nodes.RangeVar{Relname: "cp_t"},
		CmdName:    "all",
		Permissive: true,
	})

	err := c.CommentObject(&nodes.CommentStmt{
		Objtype: nodes.OBJECT_POLICY,
		Object: &nodes.List{Items: []nodes.Node{
			&nodes.String{Str: "cp_t"},
			&nodes.String{Str: "cp_pol"},
		}},
		Comment: "my policy",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCommentOnColumnRelkindValidation(t *testing.T) {
	c := New()

	// Comment on a column of a regular table — should work.
	c.DefineRelation(makeCreateTableStmt("", "cr_t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	err := c.CommentObject(&nodes.CommentStmt{
		Objtype: nodes.OBJECT_COLUMN,
		Object: &nodes.List{Items: []nodes.Node{
			&nodes.String{Str: "cr_t"},
			&nodes.String{Str: "id"},
		}},
		Comment: "col comment",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCommentOnMatviewWrongType(t *testing.T) {
	c := New()

	// Create a regular table, then try COMMENT ON MATERIALIZED VIEW — should fail.
	c.DefineRelation(makeCreateTableStmt("", "not_mv", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	err := c.CommentObject(&nodes.CommentStmt{
		Objtype: nodes.OBJECT_MATVIEW,
		Object:  &nodes.List{Items: []nodes.Node{&nodes.String{Str: "not_mv"}}},
		Comment: "not a matview",
	})
	assertErrorCode(t, err, CodeWrongObjectType)
}

// =============================================================================
// Trigger Transition Table Tests
// =============================================================================

func TestTriggerOldTransitionTable(t *testing.T) {
	c := New()

	c.DefineRelation(makeCreateTableStmt("", "tr_t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	// Create trigger function.
	c.CreateFunctionStmt(makeCreateFuncStmt("", "tr_fn",
		nil, TypeName{Name: "trigger", TypeMod: -1}, "plpgsql", "BEGIN END;", false))

	stmt := makeCreateTrigStmt("", "tr_t", "my_trig", TriggerAfter, TriggerEventDelete, true, "", "tr_fn")
	stmt.TransitionRels = &nodes.List{Items: []nodes.Node{
		&nodes.TriggerTransition{Name: "old_rows", IsNew: false},
	}}
	if err := c.CreateTriggerStmt(stmt); err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "tr_t")
	trigs := c.TriggersOf(rel.OID)
	if len(trigs) != 1 {
		t.Fatalf("expected 1 trigger, got %d", len(trigs))
	}
	if trigs[0].OldTransitionName != "old_rows" {
		t.Errorf("OldTransitionName: got %q, want %q", trigs[0].OldTransitionName, "old_rows")
	}
	if trigs[0].NewTransitionName != "" {
		t.Errorf("NewTransitionName should be empty, got %q", trigs[0].NewTransitionName)
	}
}

func TestTriggerNewTransitionTable(t *testing.T) {
	c := New()

	c.DefineRelation(makeCreateTableStmt("", "tr_t2", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	c.CreateFunctionStmt(makeCreateFuncStmt("", "tr_fn2",
		nil, TypeName{Name: "trigger", TypeMod: -1}, "plpgsql", "BEGIN END;", false))

	stmt := makeCreateTrigStmt("", "tr_t2", "my_trig2", TriggerAfter, TriggerEventInsert, true, "", "tr_fn2")
	stmt.TransitionRels = &nodes.List{Items: []nodes.Node{
		&nodes.TriggerTransition{Name: "new_rows", IsNew: true},
	}}
	if err := c.CreateTriggerStmt(stmt); err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "tr_t2")
	trigs := c.TriggersOf(rel.OID)
	if len(trigs) != 1 {
		t.Fatalf("expected 1 trigger, got %d", len(trigs))
	}
	if trigs[0].NewTransitionName != "new_rows" {
		t.Errorf("NewTransitionName: got %q, want %q", trigs[0].NewTransitionName, "new_rows")
	}
}

// =============================================================================
// Error Code Tests
// =============================================================================

func TestErrorCodeInvalidTableDefinition(t *testing.T) {
	if CodeInvalidTableDefinition != "42P16" {
		t.Errorf("CodeInvalidTableDefinition: got %q, want %q", CodeInvalidTableDefinition, "42P16")
	}
}

func TestErrorCodeInvalidColumnDefinition(t *testing.T) {
	if CodeInvalidColumnDefinition != "42611" {
		t.Errorf("CodeInvalidColumnDefinition: got %q, want %q", CodeInvalidColumnDefinition, "42611")
	}
}

func TestErrorCodeSyntaxError(t *testing.T) {
	if CodeSyntaxError != "42601" {
		t.Errorf("CodeSyntaxError: got %q, want %q", CodeSyntaxError, "42601")
	}
}

func TestErrorCodeTooManyColumns(t *testing.T) {
	if CodeTooManyColumns != "54011" {
		t.Errorf("CodeTooManyColumns: got %q, want %q", CodeTooManyColumns, "54011")
	}
}

// =============================================================================
// ALTER TABLE Column Limit Tests
// =============================================================================

func TestATAddColumnOverLimit(t *testing.T) {
	c := New()

	// Create table with exactly MaxHeapAttributeNumber columns.
	cols := make([]ColumnDef, MaxHeapAttributeNumber)
	for i := range cols {
		cols[i] = ColumnDef{Name: colName(i), Type: TypeName{Name: "int4", TypeMod: -1}}
	}
	stmt := makeCreateTableStmt("", "big_t", cols, nil, false)
	if err := c.DefineRelation(stmt, 'r'); err != nil {
		t.Fatal(err)
	}

	// Try to add one more column — should fail.
	err := c.AlterTableStmt(makeAlterTableStmt("", "big_t",
		makeATAddColumn(ColumnDef{Name: "extra", Type: TypeName{Name: "int4", TypeMod: -1}})))
	assertErrorCode(t, err, CodeTooManyColumns)
}
