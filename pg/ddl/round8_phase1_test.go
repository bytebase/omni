package ddl

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// =============================================================================
// Round 8 Phase 1: DefineRelation & MergeAttributes hardening tests
// =============================================================================

// -----------------------------------------------------------------------------
// PARTITION_MAX_KEYS limit
// -----------------------------------------------------------------------------

func TestPartitionMaxKeysLimit(t *testing.T) {
	c := New()
	// Create a table with many columns.
	var colDefs []ColumnDef
	for i := 0; i < 40; i++ {
		colDefs = append(colDefs, ColumnDef{
			Name: "c" + string(rune('a'+i%26)) + string(rune('0'+i/26)),
			Type: TypeName{Name: "int4", TypeMod: -1},
		})
	}
	stmt := makeCreateTableStmt("", "too_many_keys", colDefs, nil, false)

	// Build a partition spec with > 32 columns.
	var partElts []nodes.Node
	for i := 0; i < 33; i++ {
		partElts = append(partElts, &nodes.PartitionElem{
			Name: colDefs[i].Name,
		})
	}
	stmt.Partspec = &nodes.PartitionSpec{
		Strategy:   "r",
		PartParams: &nodes.List{Items: partElts},
	}

	err := c.DefineRelation(stmt, 'r')
	assertErrorCode(t, err, CodeProgramLimitExceeded)
	if !strings.Contains(err.Error(), "cannot partition using more than") {
		t.Errorf("unexpected message: %s", err)
	}
}

func TestPartitionWithinKeysLimit(t *testing.T) {
	c := New()
	colDefs := []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "b", Type: TypeName{Name: "int4", TypeMod: -1}},
	}
	stmt := makeCreateTableStmt("", "ok_keys", colDefs, nil, false)
	stmt.Partspec = &nodes.PartitionSpec{
		Strategy: "r",
		PartParams: &nodes.List{Items: []nodes.Node{
			&nodes.PartitionElem{Name: "a"},
			&nodes.PartitionElem{Name: "b"},
		}},
	}
	err := c.DefineRelation(stmt, 'r')
	if err != nil {
		t.Fatal(err)
	}
}

// -----------------------------------------------------------------------------
// Array dimension limit (max 6)
// -----------------------------------------------------------------------------

func TestArrayDimensionLimit(t *testing.T) {
	c := New()
	// 7 dimensions — should fail.
	var bounds []nodes.Node
	for i := 0; i < 7; i++ {
		bounds = append(bounds, &nodes.Integer{Ival: -1})
	}
	stmt := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "arr_tbl"},
		TableElts: &nodes.List{Items: []nodes.Node{
			&nodes.ColumnDef{
				Colname: "arr",
				TypeName: &nodes.TypeName{
					Names:       &nodes.List{Items: []nodes.Node{&nodes.String{Str: "int4"}}},
					ArrayBounds: &nodes.List{Items: bounds},
				},
			},
		}},
	}
	err := c.DefineRelation(stmt, 'r')
	assertErrorCode(t, err, CodeProgramLimitExceeded)
	if !strings.Contains(err.Error(), "number of array dimensions") {
		t.Errorf("unexpected message: %s", err)
	}
}

func TestArrayDimensionWithinLimit(t *testing.T) {
	c := New()
	// 6 dimensions — should succeed.
	var bounds []nodes.Node
	for i := 0; i < 6; i++ {
		bounds = append(bounds, &nodes.Integer{Ival: -1})
	}
	stmt := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "arr_tbl6"},
		TableElts: &nodes.List{Items: []nodes.Node{
			&nodes.ColumnDef{
				Colname: "arr",
				TypeName: &nodes.TypeName{
					Names:       &nodes.List{Items: []nodes.Node{&nodes.String{Str: "int4"}}},
					ArrayBounds: &nodes.List{Items: bounds},
				},
			},
		}},
	}
	err := c.DefineRelation(stmt, 'r')
	if err != nil {
		t.Fatal(err)
	}
}

// -----------------------------------------------------------------------------
// Foreign table parent in INHERITS
// -----------------------------------------------------------------------------

func TestInheritFromForeignTable(t *testing.T) {
	c := New()
	// Create a foreign table.
	ftStmt := makeCreateTableStmt("", "ft_parent", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false)
	if err := c.DefineRelation(ftStmt, 'f'); err != nil {
		t.Fatal(err)
	}

	// Inherit from foreign table — should succeed per PG.
	childStmt := makeCreateTableStmt("", "ft_child", []ColumnDef{}, nil, false)
	childStmt.InhRelations = &nodes.List{Items: []nodes.Node{
		&nodes.RangeVar{Relname: "ft_parent"},
	}}
	if err := c.DefineRelation(childStmt, 'r'); err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "ft_child")
	if len(rel.Columns) != 1 || rel.Columns[0].Name != "id" {
		t.Errorf("expected inherited column id, got %v", rel.Columns)
	}
}

// -----------------------------------------------------------------------------
// Typed table column validation (OF TYPE)
// -----------------------------------------------------------------------------

func TestTypedTableUnknownColumn(t *testing.T) {
	c := New()
	// Create composite type with column 'x'.
	c.ProcessUtility(&nodes.CompositeTypeStmt{
		Typevar: &nodes.RangeVar{Relname: "mytype"},
		Coldeflist: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "x", Type: TypeName{Name: "int4", TypeMod: -1}}),
		}},
	})

	// Specify column 'y' which doesn't exist in the type.
	stmt := &nodes.CreateStmt{
		Relation:   &nodes.RangeVar{Relname: "typed_tbl"},
		OfTypename: &nodes.TypeName{Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "mytype"}}}},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "y", Type: TypeName{Name: "text", TypeMod: -1}}),
		}},
	}
	err := c.DefineRelation(stmt, 'r')
	assertErrorCode(t, err, CodeUndefinedColumn)
}

func TestTypedTableMergeNotNull(t *testing.T) {
	c := New()
	c.ProcessUtility(&nodes.CompositeTypeStmt{
		Typevar: &nodes.RangeVar{Relname: "mytype2"},
		Coldeflist: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}}),
			makeColumnDefNode(ColumnDef{Name: "b", Type: TypeName{Name: "text", TypeMod: -1}}),
		}},
	})

	// Merge NOT NULL into type column 'a'.
	stmt := &nodes.CreateStmt{
		Relation:   &nodes.RangeVar{Relname: "typed_tbl2"},
		OfTypename: &nodes.TypeName{Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "mytype2"}}}},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}, NotNull: true}),
		}},
	}
	if err := c.DefineRelation(stmt, 'r'); err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "typed_tbl2")
	if len(rel.Columns) != 2 {
		t.Fatalf("expected 2 columns from type, got %d", len(rel.Columns))
	}
	if !rel.Columns[0].NotNull {
		t.Error("expected column a to be NOT NULL after merge")
	}
	if rel.Columns[1].NotNull {
		t.Error("expected column b to NOT be NOT NULL")
	}
}

func TestTypedTableMergeDefault(t *testing.T) {
	c := New()
	c.ProcessUtility(&nodes.CompositeTypeStmt{
		Typevar: &nodes.RangeVar{Relname: "dtype"},
		Coldeflist: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "val", Type: TypeName{Name: "int4", TypeMod: -1}}),
		}},
	})

	// Merge default into type column 'val'.
	stmt := &nodes.CreateStmt{
		Relation:   &nodes.RangeVar{Relname: "dtbl"},
		OfTypename: &nodes.TypeName{Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "dtype"}}}},
		TableElts: &nodes.List{Items: []nodes.Node{
			&nodes.ColumnDef{
				Colname: "val",
				TypeName: &nodes.TypeName{
					Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "int4"}}},
				},
				RawDefault: &nodes.Integer{Ival: 42},
			},
		}},
	}
	if err := c.DefineRelation(stmt, 'r'); err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "dtbl")
	if !rel.Columns[0].HasDefault {
		t.Error("expected column val to have a default after merge")
	}
}

func TestTypedTableNoExtraCols(t *testing.T) {
	c := New()
	c.ProcessUtility(&nodes.CompositeTypeStmt{
		Typevar: &nodes.RangeVar{Relname: "notype"},
		Coldeflist: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}}),
		}},
	})

	// No extra column specs — should succeed with just type columns.
	stmt := &nodes.CreateStmt{
		Relation:   &nodes.RangeVar{Relname: "notbl"},
		OfTypename: &nodes.TypeName{Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "notype"}}}},
	}
	if err := c.DefineRelation(stmt, 'r'); err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "notbl")
	if len(rel.Columns) != 1 {
		t.Fatalf("expected 1 column from type, got %d", len(rel.Columns))
	}
	if rel.OfTypeOID == 0 {
		t.Error("expected OfTypeOID to be set")
	}
}

// -----------------------------------------------------------------------------
// Partition bound strategy validation
// -----------------------------------------------------------------------------

func TestPartitionBoundDefaultDuplicate(t *testing.T) {
	c := New()
	// Create partitioned parent.
	parentStmt := makeCreateTableStmt("", "pt", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false)
	parentStmt.Partspec = &nodes.PartitionSpec{
		Strategy: "l",
		PartParams: &nodes.List{Items: []nodes.Node{
			&nodes.PartitionElem{Name: "id"},
		}},
	}
	if err := c.DefineRelation(parentStmt, 'r'); err != nil {
		t.Fatal(err)
	}

	// First default partition — OK.
	p1Stmt := makeCreateTableStmt("", "pt_def1", []ColumnDef{}, nil, false)
	p1Stmt.InhRelations = &nodes.List{Items: []nodes.Node{
		&nodes.RangeVar{Relname: "pt"},
	}}
	p1Stmt.Partbound = &nodes.PartitionBoundSpec{IsDefault: true}
	if err := c.DefineRelation(p1Stmt, 'r'); err != nil {
		t.Fatal(err)
	}

	// Second default partition — should fail.
	p2Stmt := makeCreateTableStmt("", "pt_def2", []ColumnDef{}, nil, false)
	p2Stmt.InhRelations = &nodes.List{Items: []nodes.Node{
		&nodes.RangeVar{Relname: "pt"},
	}}
	p2Stmt.Partbound = &nodes.PartitionBoundSpec{IsDefault: true}
	err := c.DefineRelation(p2Stmt, 'r')
	assertErrorCode(t, err, CodeInvalidObjectDefinition)
}
