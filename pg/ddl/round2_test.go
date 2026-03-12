package ddl

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// =============================================================================
// EXCLUDE Constraint Tests
// =============================================================================

func TestExcludeConstraintBasic(t *testing.T) {
	c := New()

	// Create a table with an EXCLUDE constraint.
	stmt := makeCreateTableStmt("", "reservations", []ColumnDef{
		{Name: "room", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "during", Type: TypeName{Name: "tsrange", TypeMod: -1}},
	}, nil, false)

	// Add EXCLUDE constraint via Constraints list.
	exclCon := &nodes.Constraint{
		Contype:      nodes.CONSTR_EXCLUSION,
		Conname:      "no_overlap",
		AccessMethod: "gist",
		Exclusions: &nodes.List{Items: []nodes.Node{
			&nodes.IndexElem{Name: "room"},
			&nodes.List{Items: []nodes.Node{&nodes.String{Str: "="}}},
			&nodes.IndexElem{Name: "during"},
			&nodes.List{Items: []nodes.Node{&nodes.String{Str: "&&"}}},
		}},
	}
	stmt.Constraints = &nodes.List{Items: []nodes.Node{exclCon}}

	err := c.DefineRelation(stmt, 'r')
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "reservations")
	if rel == nil {
		t.Fatal("table not found")
	}

	cons := c.ConstraintsOf(rel.OID)
	if len(cons) != 1 {
		t.Fatalf("expected 1 constraint, got %d", len(cons))
	}
	if cons[0].Type != ConstraintExclude {
		t.Errorf("expected EXCLUDE constraint, got %c", cons[0].Type)
	}
	if cons[0].Name != "no_overlap" {
		t.Errorf("constraint name: got %q, want %q", cons[0].Name, "no_overlap")
	}
	if len(cons[0].ExclOps) != 2 {
		t.Errorf("expected 2 exclusion operators, got %d", len(cons[0].ExclOps))
	} else {
		if cons[0].ExclOps[0] != "=" {
			t.Errorf("excl op 0: got %q, want %q", cons[0].ExclOps[0], "=")
		}
		if cons[0].ExclOps[1] != "&&" {
			t.Errorf("excl op 1: got %q, want %q", cons[0].ExclOps[1], "&&")
		}
	}

	// Backing index should exist with GiST access method.
	idxs := c.IndexesOf(rel.OID)
	if len(idxs) != 1 {
		t.Fatalf("expected 1 index, got %d", len(idxs))
	}
	if idxs[0].AccessMethod != "gist" {
		t.Errorf("access method: got %q, want %q", idxs[0].AccessMethod, "gist")
	}
	if idxs[0].IsUnique {
		t.Error("EXCLUDE index should not be unique")
	}
}

func TestExcludeConstraintAutoName(t *testing.T) {
	c := New()

	stmt := makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false)

	exclCon := &nodes.Constraint{
		Contype:      nodes.CONSTR_EXCLUSION,
		AccessMethod: "gist",
		Exclusions: &nodes.List{Items: []nodes.Node{
			&nodes.IndexElem{Name: "a"},
			&nodes.List{Items: []nodes.Node{&nodes.String{Str: "="}}},
		}},
	}
	stmt.Constraints = &nodes.List{Items: []nodes.Node{exclCon}}

	if err := c.DefineRelation(stmt, 'r'); err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	cons := c.ConstraintsOf(rel.OID)
	if len(cons) != 1 {
		t.Fatalf("expected 1 constraint, got %d", len(cons))
	}
	if cons[0].Name != "t_a_excl" {
		t.Errorf("auto name: got %q, want %q", cons[0].Name, "t_a_excl")
	}
}

func TestExcludeConstraintWithAccessMethod(t *testing.T) {
	c := New()

	stmt := makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "val", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false)

	exclCon := &nodes.Constraint{
		Contype:      nodes.CONSTR_EXCLUSION,
		Conname:      "t_val_excl",
		AccessMethod: "btree",
		Exclusions: &nodes.List{Items: []nodes.Node{
			&nodes.IndexElem{Name: "val"},
			&nodes.List{Items: []nodes.Node{&nodes.String{Str: "="}}},
		}},
	}
	stmt.Constraints = &nodes.List{Items: []nodes.Node{exclCon}}

	if err := c.DefineRelation(stmt, 'r'); err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	idxs := c.IndexesOf(rel.OID)
	if len(idxs) != 1 {
		t.Fatalf("expected 1 index, got %d", len(idxs))
	}
	if idxs[0].AccessMethod != "btree" {
		t.Errorf("access method: got %q, want %q", idxs[0].AccessMethod, "btree")
	}
}

// =============================================================================
// Partition Validation Tests
// =============================================================================

func TestPartitionKeyMissingColumn(t *testing.T) {
	c := New()

	stmt := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "t"},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}),
		}},
		Partspec: &nodes.PartitionSpec{
			Strategy: "l",
			PartParams: &nodes.List{Items: []nodes.Node{
				&nodes.PartitionElem{Name: "nonexistent"},
			}},
		},
	}

	err := c.DefineRelation(stmt, 'r')
	assertErrorCode(t, err, CodeUndefinedColumn)
}

func TestPartitionBoundListOverlap(t *testing.T) {
	c := New()

	// Create partitioned parent.
	parent := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "sales"},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "region", Type: TypeName{Name: "text", TypeMod: -1}}),
			makeColumnDefNode(ColumnDef{Name: "amount", Type: TypeName{Name: "int4", TypeMod: -1}}),
		}},
		Partspec: &nodes.PartitionSpec{
			Strategy: "l",
			PartParams: &nodes.List{Items: []nodes.Node{
				&nodes.PartitionElem{Name: "region"},
			}},
		},
	}
	if err := c.DefineRelation(parent, 'r'); err != nil {
		t.Fatal(err)
	}

	// First partition: 'us', 'ca'.
	p1 := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "sales_americas"},
		InhRelations: &nodes.List{Items: []nodes.Node{
			&nodes.RangeVar{Relname: "sales"},
		}},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "region", Type: TypeName{Name: "text", TypeMod: -1}}),
			makeColumnDefNode(ColumnDef{Name: "amount", Type: TypeName{Name: "int4", TypeMod: -1}}),
		}},
		Partbound: &nodes.PartitionBoundSpec{
			Strategy: 'l',
			Listdatums: &nodes.List{Items: []nodes.Node{
				&nodes.A_Const{Val: &nodes.String{Str: "us"}},
				&nodes.A_Const{Val: &nodes.String{Str: "ca"}},
			}},
		},
	}
	if err := c.DefineRelation(p1, 'r'); err != nil {
		t.Fatal(err)
	}

	// Second partition overlapping with 'us'.
	p2 := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "sales_overlap"},
		InhRelations: &nodes.List{Items: []nodes.Node{
			&nodes.RangeVar{Relname: "sales"},
		}},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "region", Type: TypeName{Name: "text", TypeMod: -1}}),
			makeColumnDefNode(ColumnDef{Name: "amount", Type: TypeName{Name: "int4", TypeMod: -1}}),
		}},
		Partbound: &nodes.PartitionBoundSpec{
			Strategy: 'l',
			Listdatums: &nodes.List{Items: []nodes.Node{
				&nodes.A_Const{Val: &nodes.String{Str: "us"}},
				&nodes.A_Const{Val: &nodes.String{Str: "uk"}},
			}},
		},
	}
	err := c.DefineRelation(p2, 'r')
	assertErrorCode(t, err, CodeInvalidObjectDefinition)
}

func TestPartitionBoundMultipleDefault(t *testing.T) {
	c := New()

	parent := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "logs"},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "level", Type: TypeName{Name: "int4", TypeMod: -1}}),
		}},
		Partspec: &nodes.PartitionSpec{
			Strategy: "l",
			PartParams: &nodes.List{Items: []nodes.Node{
				&nodes.PartitionElem{Name: "level"},
			}},
		},
	}
	if err := c.DefineRelation(parent, 'r'); err != nil {
		t.Fatal(err)
	}

	// First default partition.
	d1 := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "logs_default1"},
		InhRelations: &nodes.List{Items: []nodes.Node{
			&nodes.RangeVar{Relname: "logs"},
		}},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "level", Type: TypeName{Name: "int4", TypeMod: -1}}),
		}},
		Partbound: &nodes.PartitionBoundSpec{
			Strategy:  'l',
			IsDefault: true,
		},
	}
	if err := c.DefineRelation(d1, 'r'); err != nil {
		t.Fatal(err)
	}

	// Second default partition — should fail.
	d2 := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "logs_default2"},
		InhRelations: &nodes.List{Items: []nodes.Node{
			&nodes.RangeVar{Relname: "logs"},
		}},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "level", Type: TypeName{Name: "int4", TypeMod: -1}}),
		}},
		Partbound: &nodes.PartitionBoundSpec{
			Strategy:  'l',
			IsDefault: true,
		},
	}
	err := c.DefineRelation(d2, 'r')
	assertErrorCode(t, err, CodeInvalidObjectDefinition)
}

func TestPartitionBoundHashConflict(t *testing.T) {
	c := New()

	parent := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "ht"},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}),
		}},
		Partspec: &nodes.PartitionSpec{
			Strategy: "h",
			PartParams: &nodes.List{Items: []nodes.Node{
				&nodes.PartitionElem{Name: "id"},
			}},
		},
	}
	if err := c.DefineRelation(parent, 'r'); err != nil {
		t.Fatal(err)
	}

	// First hash partition: modulus 4, remainder 0.
	h1 := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "ht_0"},
		InhRelations: &nodes.List{Items: []nodes.Node{
			&nodes.RangeVar{Relname: "ht"},
		}},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}),
		}},
		Partbound: &nodes.PartitionBoundSpec{
			Strategy:  'h',
			Modulus:   4,
			Remainder: 0,
		},
	}
	if err := c.DefineRelation(h1, 'r'); err != nil {
		t.Fatal(err)
	}

	// Same hash partition: modulus 4, remainder 0 — should fail.
	h2 := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "ht_conflict"},
		InhRelations: &nodes.List{Items: []nodes.Node{
			&nodes.RangeVar{Relname: "ht"},
		}},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}),
		}},
		Partbound: &nodes.PartitionBoundSpec{
			Strategy:  'h',
			Modulus:   4,
			Remainder: 0,
		},
	}
	err := c.DefineRelation(h2, 'r')
	assertErrorCode(t, err, CodeInvalidObjectDefinition)
}

func TestPartitionBoundRangeOverlap(t *testing.T) {
	c := New()

	parent := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "rt"},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "val", Type: TypeName{Name: "int4", TypeMod: -1}}),
		}},
		Partspec: &nodes.PartitionSpec{
			Strategy: "r",
			PartParams: &nodes.List{Items: []nodes.Node{
				&nodes.PartitionElem{Name: "val"},
			}},
		},
	}
	if err := c.DefineRelation(parent, 'r'); err != nil {
		t.Fatal(err)
	}

	// First range: [1, 100).
	r1 := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "rt_1"},
		InhRelations: &nodes.List{Items: []nodes.Node{
			&nodes.RangeVar{Relname: "rt"},
		}},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "val", Type: TypeName{Name: "int4", TypeMod: -1}}),
		}},
		Partbound: &nodes.PartitionBoundSpec{
			Strategy: 'r',
			Lowerdatums: &nodes.List{Items: []nodes.Node{
				&nodes.A_Const{Val: &nodes.Integer{Ival: 1}},
			}},
			Upperdatums: &nodes.List{Items: []nodes.Node{
				&nodes.A_Const{Val: &nodes.Integer{Ival: 100}},
			}},
		},
	}
	if err := c.DefineRelation(r1, 'r'); err != nil {
		t.Fatal(err)
	}

	// Overlapping range: [50, 200).
	r2 := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "rt_overlap"},
		InhRelations: &nodes.List{Items: []nodes.Node{
			&nodes.RangeVar{Relname: "rt"},
		}},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "val", Type: TypeName{Name: "int4", TypeMod: -1}}),
		}},
		Partbound: &nodes.PartitionBoundSpec{
			Strategy: 'r',
			Lowerdatums: &nodes.List{Items: []nodes.Node{
				&nodes.A_Const{Val: &nodes.Integer{Ival: 50}},
			}},
			Upperdatums: &nodes.List{Items: []nodes.Node{
				&nodes.A_Const{Val: &nodes.Integer{Ival: 200}},
			}},
		},
	}
	err := c.DefineRelation(r2, 'r')
	assertErrorCode(t, err, CodeInvalidObjectDefinition)
}

// =============================================================================
// Index Validation Tests
// =============================================================================

func TestIndexOnViewBlocked(t *testing.T) {
	c := New()

	// Create a simple view.
	c.DefineRelation(makeCreateTableStmt("", "src", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	c.DefineView(makeViewStmt("", "v", &SelectStmt{
		TargetList: []*ResTarget{{Name: "id", Val: &Expr{Kind: ExprColumnRef, ColumnName: "id"}}},
		From:       []*FromItem{{Kind: FromTable, Table: "src"}},
	}, nil))

	// Try to create index on view — should fail.
	err := c.DefineIndex(makeIndexStmt("", "v", "v_idx", []string{"id"}, false, false))
	assertErrorCode(t, err, CodeWrongObjectType)
}

func TestIndexOnMatviewAllowed(t *testing.T) {
	c := New()

	// Create a matview via DefineRelation with relkind='m'.
	stmt := makeCreateTableStmt("", "mv", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false)
	if err := c.DefineRelation(stmt, 'm'); err != nil {
		t.Fatal(err)
	}

	// Create index on matview — should succeed.
	if err := c.DefineIndex(makeIndexStmt("", "mv", "mv_idx", []string{"id"}, false, false)); err != nil {
		t.Fatal(err)
	}
}

func TestATAddIndexConstraint(t *testing.T) {
	c := New()

	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "name", Type: TypeName{Name: "text", TypeMod: -1}},
	}, nil, false), 'r')

	// AT_AddIndexConstraint with an IndexStmt for a PK.
	idxStmt := &nodes.IndexStmt{
		Idxname: "t_pkey",
		Primary: true,
		IndexParams: &nodes.List{Items: []nodes.Node{
			&nodes.IndexElem{Name: "id"},
		}},
	}
	atCmd := &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_AddIndexConstraint),
		Def:     idxStmt,
	}
	err := c.AlterTableStmt(makeAlterTableStmt("", "t", atCmd))
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	cons := c.ConstraintsOf(rel.OID)
	if len(cons) != 1 {
		t.Fatalf("expected 1 constraint, got %d", len(cons))
	}
	if cons[0].Type != ConstraintPK {
		t.Errorf("expected PK, got %c", cons[0].Type)
	}
	if cons[0].Name != "t_pkey" {
		t.Errorf("name: got %q, want %q", cons[0].Name, "t_pkey")
	}
}

// =============================================================================
// LIKE Clause Tests
// =============================================================================

func TestLikeBasic(t *testing.T) {
	c := New()

	// Create source table.
	c.DefineRelation(makeCreateTableStmt("", "src", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}, NotNull: true},
		{Name: "name", Type: TypeName{Name: "text", TypeMod: -1}},
	}, nil, false), 'r')

	// Create table with LIKE.
	likeStmt := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "dst"},
		TableElts: &nodes.List{Items: []nodes.Node{
			&nodes.TableLikeClause{
				Relation: &nodes.RangeVar{Relname: "src"},
				Options:  0, // no options = columns only
			},
		}},
	}
	if err := c.DefineRelation(likeStmt, 'r'); err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "dst")
	if rel == nil {
		t.Fatal("table not found")
	}
	if len(rel.Columns) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(rel.Columns))
	}
	if rel.Columns[0].Name != "id" {
		t.Errorf("col 0 name: got %q, want %q", rel.Columns[0].Name, "id")
	}
	if rel.Columns[0].TypeOID != INT4OID {
		t.Errorf("col 0 type: got %d, want %d", rel.Columns[0].TypeOID, INT4OID)
	}
}

func TestLikeIncludingDefaults(t *testing.T) {
	c := New()

	// Source with default.
	c.DefineRelation(makeCreateTableStmt("", "src", []ColumnDef{
		{Name: "status", Type: TypeName{Name: "text", TypeMod: -1}, Default: "'active'"},
	}, nil, false), 'r')

	// LIKE with INCLUDING DEFAULTS.
	likeStmt := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "dst"},
		TableElts: &nodes.List{Items: []nodes.Node{
			&nodes.TableLikeClause{
				Relation: &nodes.RangeVar{Relname: "src"},
				Options:  tableLikeDefaults,
			},
		}},
	}
	if err := c.DefineRelation(likeStmt, 'r'); err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "dst")
	if !rel.Columns[0].HasDefault {
		t.Error("expected default to be copied")
	}
	if rel.Columns[0].Default != "'active'" {
		t.Errorf("default: got %q, want %q", rel.Columns[0].Default, "'active'")
	}
}

func TestLikeIncludingConstraints(t *testing.T) {
	c := New()

	// Source with CHECK constraint.
	c.DefineRelation(makeCreateTableStmt("", "src", []ColumnDef{
		{Name: "val", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, []ConstraintDef{
		{Name: "val_positive", Type: ConstraintCheck, CheckExpr: "(val > 0)"},
	}, false), 'r')

	// LIKE with INCLUDING CONSTRAINTS.
	likeStmt := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "dst"},
		TableElts: &nodes.List{Items: []nodes.Node{
			&nodes.TableLikeClause{
				Relation: &nodes.RangeVar{Relname: "src"},
				Options:  tableLikeConstraints,
			},
		}},
	}
	if err := c.DefineRelation(likeStmt, 'r'); err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "dst")
	cons := c.ConstraintsOf(rel.OID)
	if len(cons) != 1 {
		t.Fatalf("expected 1 constraint, got %d", len(cons))
	}
	if cons[0].Type != ConstraintCheck {
		t.Errorf("expected CHECK, got %c", cons[0].Type)
	}
}

func TestLikeNonexistentTable(t *testing.T) {
	c := New()

	likeStmt := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "dst"},
		TableElts: &nodes.List{Items: []nodes.Node{
			&nodes.TableLikeClause{
				Relation: &nodes.RangeVar{Relname: "nosuch"},
				Options:  0,
			},
		}},
	}
	err := c.DefineRelation(likeStmt, 'r')
	assertErrorCode(t, err, CodeUndefinedTable)
}

func TestLikePlusOwnColumns(t *testing.T) {
	c := New()

	c.DefineRelation(makeCreateTableStmt("", "src", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	// LIKE + additional column.
	likeStmt := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "dst"},
		TableElts: &nodes.List{Items: []nodes.Node{
			&nodes.TableLikeClause{
				Relation: &nodes.RangeVar{Relname: "src"},
				Options:  0,
			},
			makeColumnDefNode(ColumnDef{Name: "extra", Type: TypeName{Name: "text", TypeMod: -1}}),
		}},
	}
	if err := c.DefineRelation(likeStmt, 'r'); err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "dst")
	if len(rel.Columns) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(rel.Columns))
	}
	if rel.Columns[1].Name != "extra" {
		t.Errorf("col 1 name: got %q, want %q", rel.Columns[1].Name, "extra")
	}
}

func TestLikeMultipleSources(t *testing.T) {
	c := New()

	c.DefineRelation(makeCreateTableStmt("", "a", []ColumnDef{
		{Name: "x", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	c.DefineRelation(makeCreateTableStmt("", "b", []ColumnDef{
		{Name: "y", Type: TypeName{Name: "text", TypeMod: -1}},
	}, nil, false), 'r')

	likeStmt := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "dst"},
		TableElts: &nodes.List{Items: []nodes.Node{
			&nodes.TableLikeClause{Relation: &nodes.RangeVar{Relname: "a"}, Options: 0},
			&nodes.TableLikeClause{Relation: &nodes.RangeVar{Relname: "b"}, Options: 0},
		}},
	}
	if err := c.DefineRelation(likeStmt, 'r'); err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "dst")
	if len(rel.Columns) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(rel.Columns))
	}
}

func TestLikeIncludingAll(t *testing.T) {
	c := New()

	c.DefineRelation(makeCreateTableStmt("", "src", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}, NotNull: true, Default: "42"},
	}, []ConstraintDef{
		{Name: "src_check", Type: ConstraintCheck, CheckExpr: "(id > 0)"},
	}, false), 'r')

	// INCLUDING ALL = 0x7FFFFFFF.
	likeStmt := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "dst"},
		TableElts: &nodes.List{Items: []nodes.Node{
			&nodes.TableLikeClause{
				Relation: &nodes.RangeVar{Relname: "src"},
				Options:  0x7FFFFFFF,
			},
		}},
	}
	if err := c.DefineRelation(likeStmt, 'r'); err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "dst")
	if !rel.Columns[0].HasDefault {
		t.Error("expected default copied with INCLUDING ALL")
	}
	cons := c.ConstraintsOf(rel.OID)
	if len(cons) != 1 {
		t.Fatalf("expected 1 constraint, got %d", len(cons))
	}
}

// =============================================================================
// View CHECK OPTION Tests
// =============================================================================

func TestViewCheckOptionLocal(t *testing.T) {
	c := New()

	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	vstmt := makeViewStmt("", "v", &SelectStmt{
		TargetList: []*ResTarget{{Name: "id", Val: &Expr{Kind: ExprColumnRef, ColumnName: "id"}}},
		From:       []*FromItem{{Kind: FromTable, Table: "t"}},
	}, nil)
	vstmt.WithCheckOption = 1 // LOCAL
	if err := c.DefineView(vstmt); err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "v")
	if rel.CheckOption != 'l' {
		t.Errorf("CheckOption: got %c, want 'l'", rel.CheckOption)
	}
}

func TestViewCheckOptionCascaded(t *testing.T) {
	c := New()

	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	vstmt := makeViewStmt("", "v", &SelectStmt{
		TargetList: []*ResTarget{{Name: "id", Val: &Expr{Kind: ExprColumnRef, ColumnName: "id"}}},
		From:       []*FromItem{{Kind: FromTable, Table: "t"}},
	}, nil)
	vstmt.WithCheckOption = 2 // CASCADED
	if err := c.DefineView(vstmt); err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "v")
	if rel.CheckOption != 'c' {
		t.Errorf("CheckOption: got %c, want 'c'", rel.CheckOption)
	}
}

// =============================================================================
// ON COMMIT Tests
// =============================================================================

func TestOnCommitDeleteRows(t *testing.T) {
	c := New()

	stmt := makeCreateTableStmt("", "temp_t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false)
	stmt.Relation.Relpersistence = 't' // temp table
	stmt.OnCommit = nodes.ONCOMMIT_DELETE_ROWS

	if err := c.DefineRelation(stmt, 'r'); err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "temp_t")
	if rel.OnCommit != 'd' {
		t.Errorf("OnCommit: got %c, want 'd'", rel.OnCommit)
	}
	if rel.Persistence != 't' {
		t.Errorf("Persistence: got %c, want 't'", rel.Persistence)
	}
}

func TestOnCommitDrop(t *testing.T) {
	c := New()

	stmt := makeCreateTableStmt("", "temp_t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false)
	stmt.Relation.Relpersistence = 't'
	stmt.OnCommit = nodes.ONCOMMIT_DROP

	if err := c.DefineRelation(stmt, 'r'); err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "temp_t")
	if rel.OnCommit != 'D' {
		t.Errorf("OnCommit: got %c, want 'D'", rel.OnCommit)
	}
}

// =============================================================================
// Column Collation Tests
// =============================================================================

func TestColumnCollationExplicit(t *testing.T) {
	c := New()

	cd := makeColumnDefNode(ColumnDef{Name: "name", Type: TypeName{Name: "text", TypeMod: -1}})
	cd.CollClause = &nodes.CollateClause{
		Collname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "C"}}},
	}

	stmt := &nodes.CreateStmt{
		Relation:  &nodes.RangeVar{Relname: "t"},
		TableElts: &nodes.List{Items: []nodes.Node{cd}},
	}
	if err := c.DefineRelation(stmt, 'r'); err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	if rel.Columns[0].CollationName != "C" {
		t.Errorf("CollationName: got %q, want %q", rel.Columns[0].CollationName, "C")
	}
}

func TestColumnCollationInherited(t *testing.T) {
	c := New()

	// Without explicit COLLATE, CollationName should be empty (type default).
	stmt := makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "name", Type: TypeName{Name: "text", TypeMod: -1}},
	}, nil, false)
	if err := c.DefineRelation(stmt, 'r'); err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	if rel.Columns[0].CollationName != "" {
		t.Errorf("CollationName: got %q, want empty", rel.Columns[0].CollationName)
	}
}

// =============================================================================
// AlterConstraint Tests
// =============================================================================

func TestAlterConstraintDeferrable(t *testing.T) {
	c := New()

	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, []ConstraintDef{
		{Name: "t_pk", Type: ConstraintPK, Columns: []string{"id"}},
	}, false), 'r')

	// ALTER TABLE t ALTER CONSTRAINT t_pk DEFERRABLE INITIALLY DEFERRED.
	atCmd := &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_AlterConstraint),
		Def: &nodes.Constraint{
			Conname:       "t_pk",
			Deferrable:    true,
			Initdeferred:  true,
		},
	}
	if err := c.AlterTableStmt(makeAlterTableStmt("", "t", atCmd)); err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	cons := c.ConstraintsOf(rel.OID)
	if !cons[0].Deferrable {
		t.Error("expected Deferrable=true")
	}
	if !cons[0].Deferred {
		t.Error("expected Deferred=true")
	}
}

func TestAlterConstraintRevert(t *testing.T) {
	c := New()

	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, []ConstraintDef{
		{Name: "t_pk", Type: ConstraintPK, Columns: []string{"id"}},
	}, false), 'r')

	// Set deferrable.
	atCmd := &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_AlterConstraint),
		Def: &nodes.Constraint{
			Conname:    "t_pk",
			Deferrable: true,
		},
	}
	c.AlterTableStmt(makeAlterTableStmt("", "t", atCmd))

	// Revert to not deferrable.
	atCmd2 := &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_AlterConstraint),
		Def: &nodes.Constraint{
			Conname:    "t_pk",
			Deferrable: false,
		},
	}
	if err := c.AlterTableStmt(makeAlterTableStmt("", "t", atCmd2)); err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	cons := c.ConstraintsOf(rel.OID)
	if cons[0].Deferrable {
		t.Error("expected Deferrable=false after revert")
	}
}

// =============================================================================
// New DDL No-op Tests
// =============================================================================

func TestAlterTypeStmtNoOp(t *testing.T) {
	c := New()
	err := c.ProcessUtility(&nodes.AlterTypeStmt{})
	if err != nil {
		t.Fatalf("ALTER TYPE should be no-op, got: %v", err)
	}
}

func TestDefineStmtNoOp(t *testing.T) {
	c := New()
	err := c.ProcessUtility(&nodes.DefineStmt{})
	if err != nil {
		t.Fatalf("DefineStmt should be no-op, got: %v", err)
	}
}

func TestRuleStmtNoOp(t *testing.T) {
	c := New()
	err := c.ProcessUtility(&nodes.RuleStmt{})
	if err != nil {
		t.Fatalf("RuleStmt should be no-op, got: %v", err)
	}
}

func TestCreateCastStmt_EmptyErrors(t *testing.T) {
	c := New()
	// Empty CreateCastStmt (no source/target types) should error.
	err := c.ProcessUtility(&nodes.CreateCastStmt{})
	if err == nil {
		t.Fatal("expected error for empty CreateCastStmt")
	}
}

func TestCreateStatsStmtNoOp(t *testing.T) {
	c := New()
	err := c.ProcessUtility(&nodes.CreateStatsStmt{})
	if err != nil {
		t.Fatalf("CreateStatsStmt should be no-op, got: %v", err)
	}
}

func TestAlterStatsStmtNoOp(t *testing.T) {
	c := New()
	err := c.ProcessUtility(&nodes.AlterStatsStmt{})
	if err != nil {
		t.Fatalf("AlterStatsStmt should be no-op, got: %v", err)
	}
}

func TestAlterDefaultPrivilegesStmtNoOp(t *testing.T) {
	c := New()
	err := c.ProcessUtility(&nodes.AlterDefaultPrivilegesStmt{})
	if err != nil {
		t.Fatalf("AlterDefaultPrivilegesStmt should be no-op, got: %v", err)
	}
}

func TestAlterTSDictionaryStmtNoOp(t *testing.T) {
	c := New()
	err := c.ProcessUtility(&nodes.AlterTSDictionaryStmt{})
	if err != nil {
		t.Fatalf("AlterTSDictionaryStmt should be no-op, got: %v", err)
	}
}

func TestAlterTSConfigurationStmtNoOp(t *testing.T) {
	c := New()
	err := c.ProcessUtility(&nodes.AlterTSConfigurationStmt{})
	if err != nil {
		t.Fatalf("AlterTSConfigurationStmt should be no-op, got: %v", err)
	}
}

func TestCreateForeignTableStmt(t *testing.T) {
	c := New()

	stmt := &nodes.CreateForeignTableStmt{
		Base: nodes.CreateStmt{
			Relation: &nodes.RangeVar{Relname: "ft"},
			TableElts: &nodes.List{Items: []nodes.Node{
				makeColumnDefNode(ColumnDef{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}),
			}},
		},
		Servername: "foreign_srv",
	}
	if err := c.ProcessUtility(stmt); err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "ft")
	if rel == nil {
		t.Fatal("foreign table not found")
	}
	if rel.RelKind != 'f' {
		t.Errorf("RelKind: got %c, want 'f'", rel.RelKind)
	}
}
