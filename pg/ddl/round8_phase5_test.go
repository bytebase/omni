package ddl

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// =============================================================================
// Round 8 Phase 5: Trigger Hardening
// =============================================================================

// -----------------------------------------------------------------------------
// Constraint trigger creation (success)
// -----------------------------------------------------------------------------

func TestConstraintTriggerCreation(t *testing.T) {
	c := New()
	createTestTable(t, c)
	createTriggerFunc(t, c)

	// Create a second table as the constraint's referenced relation.
	if err := c.DefineRelation(makeCreateTableStmt("", "ref_t",
		[]ColumnDef{{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}},
		nil, false), 'r'); err != nil {
		t.Fatal(err)
	}

	stmt := makeCreateTrigStmt("", "t", "constr_trig", TriggerAfter, TriggerEventInsert, true, "", "trig_fn")
	stmt.IsConstraint = true
	stmt.Constrrel = &nodes.RangeVar{Relname: "ref_t"}
	stmt.Deferrable = true
	stmt.Initdeferred = true

	err := c.CreateTriggerStmt(stmt)
	if err != nil {
		t.Fatal(err)
	}

	// Verify trigger was created with constraint fields.
	rel := c.GetRelation("", "t")
	trigs := c.TriggersOf(rel.OID)
	if len(trigs) != 1 {
		t.Fatalf("triggers: got %d, want 1", len(trigs))
	}
	if !trigs[0].IsConstraint {
		t.Error("expected IsConstraint to be true")
	}
	if trigs[0].ConstraintRelOID == 0 {
		t.Error("expected ConstraintRelOID to be set")
	}
}

func TestConstraintTriggerNoConstrrel(t *testing.T) {
	c := New()
	createTestTable(t, c)
	createTriggerFunc(t, c)

	// Constraint trigger without Constrrel — should succeed with ConstraintRelOID=0.
	stmt := makeCreateTrigStmt("", "t", "constr_trig2", TriggerAfter, TriggerEventInsert, true, "", "trig_fn")
	stmt.IsConstraint = true

	err := c.CreateTriggerStmt(stmt)
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	trigs := c.TriggersOf(rel.OID)
	if len(trigs) != 1 {
		t.Fatalf("triggers: got %d, want 1", len(trigs))
	}
	if !trigs[0].IsConstraint {
		t.Error("expected IsConstraint to be true")
	}
	if trigs[0].ConstraintRelOID != 0 {
		t.Errorf("expected ConstraintRelOID=0, got %d", trigs[0].ConstraintRelOID)
	}
}

// -----------------------------------------------------------------------------
// OR REPLACE on constraint trigger (error)
// -----------------------------------------------------------------------------

func TestOrReplaceConstraintTriggerError(t *testing.T) {
	c := New()
	createTestTable(t, c)
	createTriggerFunc(t, c)

	// Create a constraint trigger first.
	stmt := makeCreateTrigStmt("", "t", "constr_trig", TriggerAfter, TriggerEventInsert, true, "", "trig_fn")
	stmt.IsConstraint = true
	if err := c.CreateTriggerStmt(stmt); err != nil {
		t.Fatal(err)
	}

	// Try to OR REPLACE it — should fail.
	replaceStmt := makeCreateTrigStmtReplace("", "t", "constr_trig", TriggerAfter, TriggerEventUpdate, true, "", "trig_fn")
	err := c.CreateTriggerStmt(replaceStmt)
	assertCode(t, err, CodeInvalidObjectDefinition)
	if !strings.Contains(err.Error(), "cannot replace trigger") {
		t.Errorf("unexpected message: %s", err)
	}
	if !strings.Contains(err.Error(), "constraint trigger") {
		t.Errorf("expected 'constraint trigger' in message: %s", err)
	}
}

func TestOrReplaceNonConstraintTriggerOK(t *testing.T) {
	c := New()
	createTestTable(t, c)
	createTriggerFunc(t, c)

	// Create a regular (non-constraint) trigger.
	stmt := makeCreateTrigStmt("", "t", "reg_trig", TriggerBefore, TriggerEventInsert, true, "", "trig_fn")
	if err := c.CreateTriggerStmt(stmt); err != nil {
		t.Fatal(err)
	}

	// OR REPLACE on non-constraint trigger — should succeed.
	replaceStmt := makeCreateTrigStmtReplace("", "t", "reg_trig", TriggerAfter, TriggerEventUpdate, true, "", "trig_fn")
	err := c.CreateTriggerStmt(replaceStmt)
	if err != nil {
		t.Fatal(err)
	}
}

// -----------------------------------------------------------------------------
// Trigger with arguments (success)
// -----------------------------------------------------------------------------

func TestTriggerWithArgs(t *testing.T) {
	c := New()
	createTestTable(t, c)
	createTriggerFunc(t, c)

	stmt := makeCreateTrigStmt("", "t", "arg_trig", TriggerBefore, TriggerEventInsert, true, "", "trig_fn")
	stmt.Args = &nodes.List{Items: []nodes.Node{
		&nodes.String{Str: "arg1"},
		&nodes.String{Str: "arg2"},
		&nodes.String{Str: "arg3"},
	}}

	err := c.CreateTriggerStmt(stmt)
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	trigs := c.TriggersOf(rel.OID)
	if len(trigs) != 1 {
		t.Fatalf("triggers: got %d, want 1", len(trigs))
	}
	if len(trigs[0].Args) != 3 {
		t.Fatalf("trigger args: got %d, want 3", len(trigs[0].Args))
	}
	if trigs[0].Args[0] != "arg1" || trigs[0].Args[1] != "arg2" || trigs[0].Args[2] != "arg3" {
		t.Errorf("trigger args: got %v, want [arg1 arg2 arg3]", trigs[0].Args)
	}
}

func TestTriggerNoArgs(t *testing.T) {
	c := New()
	createTestTable(t, c)
	createTriggerFunc(t, c)

	stmt := makeCreateTrigStmt("", "t", "noarg_trig", TriggerBefore, TriggerEventInsert, true, "", "trig_fn")
	err := c.CreateTriggerStmt(stmt)
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	trigs := c.TriggersOf(rel.OID)
	if len(trigs[0].Args) != 0 {
		t.Errorf("trigger args: got %d, want 0", len(trigs[0].Args))
	}
}

// -----------------------------------------------------------------------------
// Inheritance child + ROW + transition tables (error)
// -----------------------------------------------------------------------------

func TestInheritanceChildRowTransitionTablesError(t *testing.T) {
	c := New()
	createTriggerFunc(t, c)

	// Create parent table.
	parent := makeCreateTableStmt("", "parent", []ColumnDef{
		{Name: "x", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false)
	if err := c.DefineRelation(parent, 'r'); err != nil {
		t.Fatal(err)
	}

	// Create child table that inherits from parent (non-partition).
	child := makeCreateTableStmt("", "child", []ColumnDef{
		{Name: "x", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false)
	child.InhRelations = &nodes.List{Items: []nodes.Node{
		&nodes.RangeVar{Relname: "parent"},
	}}
	if err := c.DefineRelation(child, 'r'); err != nil {
		t.Fatal(err)
	}

	// Verify child has InhCount > 0 and is not a partition.
	childRel := c.GetRelation("", "child")
	if childRel.InhCount == 0 {
		t.Fatal("expected child InhCount > 0")
	}
	if childRel.PartitionOf != 0 {
		t.Fatal("expected child to not be a partition")
	}

	// Create AFTER ROW trigger with transition tables on child — should fail.
	stmt := makeCreateTrigStmt("", "child", "trig1", TriggerAfter, TriggerEventInsert, true, "", "trig_fn")
	stmt.TransitionRels = &nodes.List{Items: []nodes.Node{
		&nodes.TriggerTransition{Name: "new_tbl", IsNew: true},
	}}
	err := c.CreateTriggerStmt(stmt)
	assertCode(t, err, CodeFeatureNotSupported)
	if !strings.Contains(err.Error(), "inheritance children") {
		t.Errorf("unexpected message: %s", err)
	}
}

// -----------------------------------------------------------------------------
// Non-inheritance table + ROW + transition tables (success)
// -----------------------------------------------------------------------------

func TestNonInheritanceTableRowTransitionTablesOK(t *testing.T) {
	c := New()
	createTestTable(t, c)
	createTriggerFunc(t, c)

	// Regular table (no inheritance) with ROW trigger + transition tables — should succeed.
	stmt := makeCreateTrigStmt("", "t", "trig1", TriggerAfter, TriggerEventInsert, true, "", "trig_fn")
	stmt.TransitionRels = &nodes.List{Items: []nodes.Node{
		&nodes.TriggerTransition{Name: "new_tbl", IsNew: true},
	}}
	err := c.CreateTriggerStmt(stmt)
	if err != nil {
		t.Fatal(err)
	}
}

// -----------------------------------------------------------------------------
// Statement-level trigger with transition tables on inheritance child (OK)
// The inheritance check only applies to ROW triggers.
// -----------------------------------------------------------------------------

func TestInheritanceChildStatementTransitionTablesOK(t *testing.T) {
	c := New()
	createTriggerFunc(t, c)

	// Create parent and child.
	parent := makeCreateTableStmt("", "parent2", []ColumnDef{
		{Name: "x", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false)
	if err := c.DefineRelation(parent, 'r'); err != nil {
		t.Fatal(err)
	}

	child := makeCreateTableStmt("", "child2", []ColumnDef{
		{Name: "x", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false)
	child.InhRelations = &nodes.List{Items: []nodes.Node{
		&nodes.RangeVar{Relname: "parent2"},
	}}
	if err := c.DefineRelation(child, 'r'); err != nil {
		t.Fatal(err)
	}

	// Statement-level AFTER trigger with transition tables on child — should succeed.
	// (The inheritance child check only applies to FOR EACH ROW triggers.)
	stmt := makeCreateTrigStmt("", "child2", "trig1", TriggerAfter, TriggerEventInsert, false, "", "trig_fn")
	stmt.TransitionRels = &nodes.List{Items: []nodes.Node{
		&nodes.TriggerTransition{Name: "new_tbl", IsNew: true},
	}}
	err := c.CreateTriggerStmt(stmt)
	if err != nil {
		t.Fatal(err)
	}
}

// -----------------------------------------------------------------------------
// Trigger function return type validation
// (Already tested in TestCreateTriggerWrongReturnType, verify message matches PG)
// -----------------------------------------------------------------------------

func TestTriggerFuncReturnTypeMessage(t *testing.T) {
	c := New()
	createTestTable(t, c)

	// Create a function that returns int4 (not trigger).
	c.CreateFunctionStmt(makeNonTriggerFuncStmt("", "bad_fn2", TypeName{Name: "int4", TypeMod: -1}, "sql", "SELECT 1"))

	err := c.CreateTriggerStmt(makeCreateTrigStmt("", "t", "trig1", TriggerBefore, TriggerEventInsert, false, "", "bad_fn2"))
	assertCode(t, err, CodeInvalidObjectDefinition)
	if !strings.Contains(err.Error(), "must return type trigger") {
		t.Errorf("expected 'must return type trigger' in message: %s", err)
	}
}
