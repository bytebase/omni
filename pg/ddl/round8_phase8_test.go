package ddl

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// =============================================================================
// Round 8 Phase 8: Miscellaneous Cross-File Gaps
// =============================================================================

// -----------------------------------------------------------------------------
// Gap 1: TRUNCATE CASCADE
// pg: src/backend/commands/tablecmds.c — ExecuteTruncateGuts
// -----------------------------------------------------------------------------

func TestTruncateCascadeOK(t *testing.T) {
	c := New()
	// Create two tables with FK relationship.
	if err := c.DefineRelation(makeCreateTableStmt("", "parent", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, []ConstraintDef{
		{Type: ConstraintPK, Columns: []string{"id"}},
	}, false), 'r'); err != nil {
		t.Fatal(err)
	}
	if err := c.DefineRelation(makeCreateTableStmt("", "child", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "pid", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, []ConstraintDef{
		{Type: ConstraintFK, Columns: []string{"pid"}, RefTable: "parent", RefColumns: []string{"id"}},
	}, false), 'r'); err != nil {
		t.Fatal(err)
	}

	// TRUNCATE parent CASCADE — should succeed.
	stmt := &nodes.TruncateStmt{
		Relations: &nodes.List{Items: []nodes.Node{
			&nodes.RangeVar{Relname: "parent"},
		}},
		Behavior: nodes.DROP_CASCADE,
	}
	err := c.ExecuteTruncate(stmt)
	if err != nil {
		t.Fatal(err)
	}
}

// -----------------------------------------------------------------------------
// Gap 6: TRUNCATE on foreign table rejection (verify existing behavior)
// pg: src/backend/commands/tablecmds.c — ExecuteTruncate
// -----------------------------------------------------------------------------

func TestTruncateForeignTableRejected(t *testing.T) {
	c := New()
	// Create a foreign table.
	c.ProcessUtility(&nodes.CreateForeignTableStmt{
		Base: nodes.CreateStmt{
			Relation: &nodes.RangeVar{Relname: "ft"},
			TableElts: &nodes.List{Items: []nodes.Node{
				makeColumnDefNode(ColumnDef{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}),
			}},
		},
	})

	// TRUNCATE ft — should fail because foreign tables cannot be truncated.
	stmt := &nodes.TruncateStmt{
		Relations: &nodes.List{Items: []nodes.Node{
			&nodes.RangeVar{Relname: "ft"},
		}},
	}
	err := c.ExecuteTruncate(stmt)
	assertCode(t, err, CodeWrongObjectType)
}

// -----------------------------------------------------------------------------
// Gap 10: TRUNCATE on partitioned table
// pg: src/backend/commands/tablecmds.c — ExecuteTruncate
// -----------------------------------------------------------------------------

func TestTruncatePartitionedTableOK(t *testing.T) {
	c := New()
	// Create a partitioned table.
	if err := c.DefineRelation(&nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "pt"},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}),
		}},
		Partspec: &nodes.PartitionSpec{
			Strategy: "range",
			PartParams: &nodes.List{Items: []nodes.Node{
				&nodes.PartitionElem{Name: "id"},
			}},
		},
	}, 'p'); err != nil {
		t.Fatal(err)
	}

	// TRUNCATE pt — should succeed for partitioned tables.
	stmt := &nodes.TruncateStmt{
		Relations: &nodes.List{Items: []nodes.Node{
			&nodes.RangeVar{Relname: "pt", Inh: true},
		}},
	}
	err := c.ExecuteTruncate(stmt)
	if err != nil {
		t.Fatal(err)
	}
}

// -----------------------------------------------------------------------------
// Gap 3/8: Primary key columns set NOT NULL
// pg: src/backend/commands/indexcmds.c — index_check_primary_key
// -----------------------------------------------------------------------------

func TestPrimaryKeyColumnsSetNotNull(t *testing.T) {
	c := New()
	// Create a table with a nullable column and add PK constraint.
	if err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "name", Type: TypeName{Name: "text", TypeMod: -1}},
	}, []ConstraintDef{
		{Type: ConstraintPK, Columns: []string{"id"}},
	}, false), 'r'); err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	if rel == nil {
		t.Fatal("table not found")
	}
	// PK column "id" should be NOT NULL.
	if !rel.Columns[0].NotNull {
		t.Error("PK column 'id' should be NOT NULL")
	}
}

func TestPrimaryKeyAlreadyNotNull(t *testing.T) {
	c := New()
	// Create a table where the PK column is already NOT NULL.
	if err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}, NotNull: true},
	}, []ConstraintDef{
		{Type: ConstraintPK, Columns: []string{"id"}},
	}, false), 'r'); err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	if rel == nil {
		t.Fatal("table not found")
	}
	if !rel.Columns[0].NotNull {
		t.Error("PK column 'id' should remain NOT NULL")
	}
}

// -----------------------------------------------------------------------------
// Gap 4: Privilege type validation
// pg: src/backend/catalog/aclchk.c — ExecGrant_Relation, ExecGrant_Sequence, etc.
// -----------------------------------------------------------------------------

func TestGrantExecuteOnTableError(t *testing.T) {
	c := New()
	if err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r'); err != nil {
		t.Fatal(err)
	}

	// GRANT EXECUTE ON TABLE t TO PUBLIC — EXECUTE is invalid for tables.
	stmt := &nodes.GrantStmt{
		IsGrant: true,
		Objtype: nodes.OBJECT_TABLE,
		Objects: &nodes.List{Items: []nodes.Node{
			&nodes.RangeVar{Relname: "t"},
		}},
		Privileges: &nodes.List{Items: []nodes.Node{
			&nodes.AccessPriv{PrivName: "EXECUTE"},
		}},
		Grantees: &nodes.List{Items: []nodes.Node{
			&nodes.RoleSpec{Roletype: int(nodes.ROLESPEC_PUBLIC)},
		}},
	}
	err := c.ExecGrantStmt(stmt)
	assertCode(t, err, CodeInvalidGrantOperation)
	if !strings.Contains(err.Error(), "EXECUTE") {
		t.Errorf("expected EXECUTE in error message, got: %s", err)
	}
}

func TestGrantSelectOnFunctionError(t *testing.T) {
	c := New()
	// Create a function.
	funcStmt := &nodes.CreateFunctionStmt{
		Funcname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "myfunc"}}},
		ReturnType: makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
		Options: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{Defname: "language", Arg: &nodes.String{Str: "sql"}},
			&nodes.DefElem{Defname: "as", Arg: &nodes.List{Items: []nodes.Node{
				&nodes.String{Str: "SELECT 1"},
			}}},
		}},
	}
	if err := c.CreateFunctionStmt(funcStmt); err != nil {
		t.Fatal(err)
	}

	// GRANT SELECT ON FUNCTION myfunc TO PUBLIC — SELECT is invalid for functions.
	grantStmt := &nodes.GrantStmt{
		IsGrant: true,
		Objtype: nodes.OBJECT_FUNCTION,
		Objects: &nodes.List{Items: []nodes.Node{
			&nodes.ObjectWithArgs{
				Objname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "myfunc"}}},
				ArgsUnspecified: true,
			},
		}},
		Privileges: &nodes.List{Items: []nodes.Node{
			&nodes.AccessPriv{PrivName: "SELECT"},
		}},
		Grantees: &nodes.List{Items: []nodes.Node{
			&nodes.RoleSpec{Roletype: int(nodes.ROLESPEC_PUBLIC)},
		}},
	}
	err := c.ExecGrantStmt(grantStmt)
	assertCode(t, err, CodeInvalidGrantOperation)
	if !strings.Contains(err.Error(), "SELECT") {
		t.Errorf("expected SELECT in error message, got: %s", err)
	}
}

func TestGrantUsageOnSequenceOK(t *testing.T) {
	c := New()
	// Create a sequence.
	seqStmt := &nodes.CreateSeqStmt{
		Sequence: &nodes.RangeVar{Relname: "myseq"},
	}
	if err := c.DefineSequence(seqStmt); err != nil {
		t.Fatal(err)
	}

	// GRANT USAGE ON SEQUENCE myseq TO PUBLIC — should succeed.
	grantStmt := &nodes.GrantStmt{
		IsGrant: true,
		Objtype: nodes.OBJECT_SEQUENCE,
		Objects: &nodes.List{Items: []nodes.Node{
			&nodes.RangeVar{Relname: "myseq"},
		}},
		Privileges: &nodes.List{Items: []nodes.Node{
			&nodes.AccessPriv{PrivName: "USAGE"},
		}},
		Grantees: &nodes.List{Items: []nodes.Node{
			&nodes.RoleSpec{Roletype: int(nodes.ROLESPEC_PUBLIC)},
		}},
	}
	err := c.ExecGrantStmt(grantStmt)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGrantAllOnSchemaOK(t *testing.T) {
	c := New()
	// GRANT ALL PRIVILEGES ON SCHEMA public TO PUBLIC — should succeed.
	grantStmt := &nodes.GrantStmt{
		IsGrant: true,
		Objtype: nodes.OBJECT_SCHEMA,
		Objects: &nodes.List{Items: []nodes.Node{
			&nodes.String{Str: "public"},
		}},
		// nil Privileges means ALL PRIVILEGES.
		Grantees: &nodes.List{Items: []nodes.Node{
			&nodes.RoleSpec{Roletype: int(nodes.ROLESPEC_PUBLIC)},
		}},
	}
	err := c.ExecGrantStmt(grantStmt)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPrivilegeValidationInsertOnSequence(t *testing.T) {
	c := New()
	// Create a sequence.
	seqStmt := &nodes.CreateSeqStmt{
		Sequence: &nodes.RangeVar{Relname: "myseq"},
	}
	if err := c.DefineSequence(seqStmt); err != nil {
		t.Fatal(err)
	}

	// GRANT INSERT ON SEQUENCE myseq TO PUBLIC — INSERT is invalid for sequences.
	grantStmt := &nodes.GrantStmt{
		IsGrant: true,
		Objtype: nodes.OBJECT_SEQUENCE,
		Objects: &nodes.List{Items: []nodes.Node{
			&nodes.RangeVar{Relname: "myseq"},
		}},
		Privileges: &nodes.List{Items: []nodes.Node{
			&nodes.AccessPriv{PrivName: "INSERT"},
		}},
		Grantees: &nodes.List{Items: []nodes.Node{
			&nodes.RoleSpec{Roletype: int(nodes.ROLESPEC_PUBLIC)},
		}},
	}
	err := c.ExecGrantStmt(grantStmt)
	assertCode(t, err, CodeInvalidGrantOperation)
	if !strings.Contains(err.Error(), "INSERT") {
		t.Errorf("expected INSERT in error message, got: %s", err)
	}
}

// -----------------------------------------------------------------------------
// Gap 5: DROP FUNCTION on aggregate must hint DROP AGGREGATE
// pg: src/backend/commands/dropcmds.c — RemoveObjects (line 91-99)
// -----------------------------------------------------------------------------

func TestDropFunctionOnAggregateHint(t *testing.T) {
	c := New()
	// Create a function and manually set its Kind to 'a' (aggregate).
	funcStmt := &nodes.CreateFunctionStmt{
		Funcname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "myagg"}}},
		ReturnType: makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
		Options: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{Defname: "language", Arg: &nodes.String{Str: "sql"}},
			&nodes.DefElem{Defname: "as", Arg: &nodes.List{Items: []nodes.Node{
				&nodes.String{Str: "SELECT 1"},
			}}},
		}},
	}
	if err := c.CreateFunctionStmt(funcStmt); err != nil {
		t.Fatal(err)
	}

	// Find the proc and change its kind to 'a' (aggregate).
	schema, _ := c.resolveTargetSchema("")
	candidates := c.findUserProcsByName(schema, "myagg")
	if len(candidates) == 0 {
		t.Fatal("function not found")
	}
	bp := c.procByOID[candidates[0].OID]
	bp.Kind = 'a'

	// DROP FUNCTION myagg — should error with hint about DROP AGGREGATE.
	dropStmt := &nodes.DropStmt{
		RemoveType: int(nodes.OBJECT_FUNCTION),
		Objects: &nodes.List{Items: []nodes.Node{
			&nodes.ObjectWithArgs{
				Objname:         &nodes.List{Items: []nodes.Node{&nodes.String{Str: "myagg"}}},
				ArgsUnspecified: true,
			},
		}},
	}
	err := c.removeFunctionObjects(dropStmt)
	assertCode(t, err, CodeWrongObjectType)
	if !strings.Contains(err.Error(), "aggregate") {
		t.Errorf("expected 'aggregate' in error message, got: %s", err)
	}
	if !strings.Contains(err.Error(), "DROP AGGREGATE") {
		t.Errorf("expected 'DROP AGGREGATE' hint in error message, got: %s", err)
	}
}

// -----------------------------------------------------------------------------
// Gap 2/7: Exclusion constraint AM validation
// pg: src/backend/commands/indexcmds.c — DefineIndex (line 875-879)
// -----------------------------------------------------------------------------

func TestExcludeConstraintAMGistOK(t *testing.T) {
	c := New()
	stmt := makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false)
	// Add EXCLUDE constraint with gist AM directly.
	exclCon := &nodes.Constraint{
		Contype:      nodes.CONSTR_EXCLUSION,
		Conname:      "excl_a",
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

	// Verify the constraint exists.
	rel := c.GetRelation("", "t")
	if rel == nil {
		t.Fatal("table not found")
	}
	cons := c.consByRel[rel.OID]
	found := false
	for _, con := range cons {
		if con.Type == ConstraintExclude {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected EXCLUDE constraint to be created")
	}
}

func TestExcludeConstraintAMGinRejected(t *testing.T) {
	c := New()
	stmt := makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false)
	// Add EXCLUDE constraint with gin AM — should be rejected.
	exclCon := &nodes.Constraint{
		Contype:      nodes.CONSTR_EXCLUSION,
		Conname:      "excl_a",
		AccessMethod: "gin",
		Exclusions: &nodes.List{Items: []nodes.Node{
			&nodes.IndexElem{Name: "a"},
			&nodes.List{Items: []nodes.Node{&nodes.String{Str: "="}}},
		}},
	}
	stmt.Constraints = &nodes.List{Items: []nodes.Node{exclCon}}

	err := c.DefineRelation(stmt, 'r')
	assertCode(t, err, CodeFeatureNotSupported)
	if !strings.Contains(err.Error(), "gin") {
		t.Errorf("expected 'gin' in error message, got: %s", err)
	}
	if !strings.Contains(err.Error(), "exclusion") {
		t.Errorf("expected 'exclusion' in error message, got: %s", err)
	}
}
