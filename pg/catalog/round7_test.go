package catalog

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// =============================================================================
// Round 7: Tests for gap analysis implementation (validations added by agents)
// =============================================================================

// -----------------------------------------------------------------------------
// typecmds.go: Domain NULL/NOT NULL conflict detection
// -----------------------------------------------------------------------------

func TestDomainConflictNullNotNull(t *testing.T) {
	c := New()
	// CREATE DOMAIN d AS int NOT NULL NULL — conflicting
	stmt := &nodes.CreateDomainStmt{
		Domainname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "d"}}},
		Typname:    makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
		Constraints: &nodes.List{Items: []nodes.Node{
			&nodes.Constraint{Contype: nodes.CONSTR_NOTNULL},
			&nodes.Constraint{Contype: nodes.CONSTR_NULL},
		}},
	}
	err := c.DefineDomain(stmt)
	assertCode(t, err, CodeSyntaxError)
	if !strings.Contains(err.Error(), "conflicting NULL/NOT NULL") {
		t.Errorf("unexpected message: %s", err)
	}
}

func TestDomainConflictNullNotNullReverse(t *testing.T) {
	c := New()
	// CREATE DOMAIN d AS int NULL NOT NULL — conflicting (reversed order)
	stmt := &nodes.CreateDomainStmt{
		Domainname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "d"}}},
		Typname:    makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
		Constraints: &nodes.List{Items: []nodes.Node{
			&nodes.Constraint{Contype: nodes.CONSTR_NULL},
			&nodes.Constraint{Contype: nodes.CONSTR_NOTNULL},
		}},
	}
	err := c.DefineDomain(stmt)
	assertCode(t, err, CodeSyntaxError)
}

func TestDomainNotNullOK(t *testing.T) {
	c := New()
	stmt := &nodes.CreateDomainStmt{
		Domainname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "d"}}},
		Typname:    makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
		Constraints: &nodes.List{Items: []nodes.Node{
			&nodes.Constraint{Contype: nodes.CONSTR_NOTNULL},
		}},
	}
	err := c.DefineDomain(stmt)
	if err != nil {
		t.Fatal(err)
	}
}

// -----------------------------------------------------------------------------
// typecmds.go: Domain multiple DEFAULT rejection
// -----------------------------------------------------------------------------

func TestDomainMultipleDefaultError(t *testing.T) {
	c := New()
	stmt := &nodes.CreateDomainStmt{
		Domainname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "d"}}},
		Typname:    makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
		Constraints: &nodes.List{Items: []nodes.Node{
			&nodes.Constraint{Contype: nodes.CONSTR_DEFAULT, CookedExpr: "1"},
			&nodes.Constraint{Contype: nodes.CONSTR_DEFAULT, CookedExpr: "2"},
		}},
	}
	err := c.DefineDomain(stmt)
	assertCode(t, err, CodeSyntaxError)
	if !strings.Contains(err.Error(), "multiple default") {
		t.Errorf("unexpected message: %s", err)
	}
}

// -----------------------------------------------------------------------------
// typecmds.go: Domain NO INHERIT CHECK rejection
// -----------------------------------------------------------------------------

func TestDomainCheckNoInheritError(t *testing.T) {
	c := New()
	stmt := &nodes.CreateDomainStmt{
		Domainname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "d"}}},
		Typname:    makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
		Constraints: &nodes.List{Items: []nodes.Node{
			&nodes.Constraint{
				Contype:     nodes.CONSTR_CHECK,
				Conname:     "chk",
				CookedExpr:  "VALUE > 0",
				IsNoInherit: true,
			},
		}},
	}
	err := c.DefineDomain(stmt)
	assertCode(t, err, CodeInvalidObjectDefinition)
	if !strings.Contains(err.Error(), "NO INHERIT") {
		t.Errorf("unexpected message: %s", err)
	}
}

// -----------------------------------------------------------------------------
// typecmds.go: Domain deferrability rejection
// -----------------------------------------------------------------------------

func TestDomainDeferrableError(t *testing.T) {
	c := New()
	stmt := &nodes.CreateDomainStmt{
		Domainname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "d"}}},
		Typname:    makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
		Constraints: &nodes.List{Items: []nodes.Node{
			&nodes.Constraint{Contype: nodes.CONSTR_ATTR_DEFERRABLE},
		}},
	}
	err := c.DefineDomain(stmt)
	assertCode(t, err, CodeFeatureNotSupported)
	if !strings.Contains(err.Error(), "deferrability") {
		t.Errorf("unexpected message: %s", err)
	}
}

func TestDomainNotDeferrableError(t *testing.T) {
	c := New()
	stmt := &nodes.CreateDomainStmt{
		Domainname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "d"}}},
		Typname:    makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
		Constraints: &nodes.List{Items: []nodes.Node{
			&nodes.Constraint{Contype: nodes.CONSTR_ATTR_NOT_DEFERRABLE},
		}},
	}
	err := c.DefineDomain(stmt)
	assertCode(t, err, CodeFeatureNotSupported)
}

func TestDomainDeferredError(t *testing.T) {
	c := New()
	stmt := &nodes.CreateDomainStmt{
		Domainname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "d"}}},
		Typname:    makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
		Constraints: &nodes.List{Items: []nodes.Node{
			&nodes.Constraint{Contype: nodes.CONSTR_ATTR_DEFERRED},
		}},
	}
	err := c.DefineDomain(stmt)
	assertCode(t, err, CodeFeatureNotSupported)
}

func TestDomainImmediateError(t *testing.T) {
	c := New()
	stmt := &nodes.CreateDomainStmt{
		Domainname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "d"}}},
		Typname:    makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
		Constraints: &nodes.List{Items: []nodes.Node{
			&nodes.Constraint{Contype: nodes.CONSTR_ATTR_IMMEDIATE},
		}},
	}
	err := c.DefineDomain(stmt)
	assertCode(t, err, CodeFeatureNotSupported)
}

// -----------------------------------------------------------------------------
// typecmds.go: Range pseudo-type rejection
// -----------------------------------------------------------------------------

func TestDefineRangePseudoTypeError(t *testing.T) {
	c := New()
	// "any" is a pseudo-type — cannot be a range subtype.
	stmt := &nodes.CreateRangeStmt{
		TypeName: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "anyrange"}}},
		Params: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{
				Defname: "subtype",
				Arg: &nodes.TypeName{
					Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "anyelement"}}},
				},
			},
		}},
	}
	err := c.DefineRange(stmt)
	assertCode(t, err, CodeDatatypeMismatch)
	if !strings.Contains(err.Error(), "range subtype cannot be") {
		t.Errorf("unexpected message: %s", err)
	}
}

// -----------------------------------------------------------------------------
// tablecmds.go: Cannot inherit from partition
// -----------------------------------------------------------------------------

func TestInheritFromPartitionError(t *testing.T) {
	c := New()
	// Create partitioned parent.
	parent := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "parent"},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}),
		}},
		Partspec: &nodes.PartitionSpec{
			Strategy: "range",
			PartParams: &nodes.List{Items: []nodes.Node{
				&nodes.PartitionElem{Name: "id"},
			}},
		},
	}
	if err := c.DefineRelation(parent, 'r'); err != nil {
		t.Fatal(err)
	}
	// Create partition of parent.
	part := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "part"},
		Partbound: &nodes.PartitionBoundSpec{
			Strategy:    'r',
			Lowerdatums: &nodes.List{Items: []nodes.Node{&nodes.Integer{Ival: 1}}},
			Upperdatums: &nodes.List{Items: []nodes.Node{&nodes.Integer{Ival: 100}}},
		},
		InhRelations: &nodes.List{Items: []nodes.Node{
			&nodes.RangeVar{Relname: "parent"},
		}},
	}
	if err := c.DefineRelation(part, 'r'); err != nil {
		t.Fatal(err)
	}
	// Now try to inherit from the partition using regular INHERITS — should fail.
	child := makeCreateTableStmt("", "child", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false)
	child.InhRelations = &nodes.List{Items: []nodes.Node{
		&nodes.RangeVar{Relname: "part"},
	}}
	err := c.DefineRelation(child, 'r')
	assertCode(t, err, CodeWrongObjectType)
	if !strings.Contains(err.Error(), "cannot inherit from partition") {
		t.Errorf("unexpected message: %s", err)
	}
}

// -----------------------------------------------------------------------------
// tablecmds.go: Cannot inherit from partitioned table via INHERITS
// -----------------------------------------------------------------------------

func TestInheritFromPartitionedTableError(t *testing.T) {
	c := New()
	// Create partitioned table.
	parent := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "ptable"},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}),
		}},
		Partspec: &nodes.PartitionSpec{
			Strategy: "range",
			PartParams: &nodes.List{Items: []nodes.Node{
				&nodes.PartitionElem{Name: "id"},
			}},
		},
	}
	if err := c.DefineRelation(parent, 'r'); err != nil {
		t.Fatal(err)
	}
	// Try to inherit from it using regular INHERITS (not PARTITION OF) — should fail.
	child := makeCreateTableStmt("", "child", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false)
	child.InhRelations = &nodes.List{Items: []nodes.Node{
		&nodes.RangeVar{Relname: "ptable"},
	}}
	err := c.DefineRelation(child, 'r')
	assertCode(t, err, CodeWrongObjectType)
	if !strings.Contains(err.Error(), "cannot inherit from partitioned table") {
		t.Errorf("unexpected message: %s", err)
	}
}

// -----------------------------------------------------------------------------
// trigger.go: Transition table OLD/NEW event-specific validation
// -----------------------------------------------------------------------------

func TestTriggerOldTableOnInsertOnlyError(t *testing.T) {
	c := New()
	createTestTable(t, c)
	createTriggerFunc(t, c)

	stmt := makeCreateTrigStmt("", "t", "trig1", TriggerAfter, TriggerEventInsert, true, "", "trig_fn")
	stmt.TransitionRels = &nodes.List{Items: []nodes.Node{
		&nodes.TriggerTransition{Name: "old_tbl", IsNew: false},
	}}
	err := c.CreateTriggerStmt(stmt)
	assertCode(t, err, CodeInvalidObjectDefinition)
	if !strings.Contains(err.Error(), "OLD TABLE") {
		t.Errorf("unexpected message: %s", err)
	}
}

func TestTriggerNewTableOnDeleteOnlyError(t *testing.T) {
	c := New()
	createTestTable(t, c)
	createTriggerFunc(t, c)

	stmt := makeCreateTrigStmt("", "t", "trig1", TriggerAfter, TriggerEventDelete, true, "", "trig_fn")
	stmt.TransitionRels = &nodes.List{Items: []nodes.Node{
		&nodes.TriggerTransition{Name: "new_tbl", IsNew: true},
	}}
	err := c.CreateTriggerStmt(stmt)
	assertCode(t, err, CodeInvalidObjectDefinition)
	if !strings.Contains(err.Error(), "NEW TABLE") {
		t.Errorf("unexpected message: %s", err)
	}
}

func TestTriggerOldTableOnUpdateOK(t *testing.T) {
	c := New()
	createTestTable(t, c)
	createTriggerFunc(t, c)

	stmt := makeCreateTrigStmt("", "t", "trig1", TriggerAfter, TriggerEventUpdate, true, "", "trig_fn")
	stmt.TransitionRels = &nodes.List{Items: []nodes.Node{
		&nodes.TriggerTransition{Name: "old_tbl", IsNew: false},
	}}
	err := c.CreateTriggerStmt(stmt)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTriggerNewTableOnInsertOK(t *testing.T) {
	c := New()
	createTestTable(t, c)
	createTriggerFunc(t, c)

	stmt := makeCreateTrigStmt("", "t", "trig2", TriggerAfter, TriggerEventInsert, true, "", "trig_fn")
	stmt.TransitionRels = &nodes.List{Items: []nodes.Node{
		&nodes.TriggerTransition{Name: "new_tbl", IsNew: true},
	}}
	err := c.CreateTriggerStmt(stmt)
	if err != nil {
		t.Fatal(err)
	}
}

// -----------------------------------------------------------------------------
// trigger.go: Duplicate transition table name
// -----------------------------------------------------------------------------

func TestTriggerDuplicateTransitionNames(t *testing.T) {
	c := New()
	createTestTable(t, c)
	createTriggerFunc(t, c)

	stmt := makeCreateTrigStmt("", "t", "trig1", TriggerAfter, TriggerEventUpdate, true, "", "trig_fn")
	stmt.TransitionRels = &nodes.List{Items: []nodes.Node{
		&nodes.TriggerTransition{Name: "tbl", IsNew: false},
		&nodes.TriggerTransition{Name: "tbl", IsNew: true},
	}}
	err := c.CreateTriggerStmt(stmt)
	assertCode(t, err, CodeInvalidObjectDefinition)
	if !strings.Contains(err.Error(), "OLD TABLE name and NEW TABLE name cannot be the same") {
		t.Errorf("unexpected message: %s", err)
	}
}

func TestTriggerDifferentTransitionNamesOK(t *testing.T) {
	c := New()
	createTestTable(t, c)
	createTriggerFunc(t, c)

	stmt := makeCreateTrigStmt("", "t", "trig1", TriggerAfter, TriggerEventUpdate, true, "", "trig_fn")
	stmt.TransitionRels = &nodes.List{Items: []nodes.Node{
		&nodes.TriggerTransition{Name: "old_tbl", IsNew: false},
		&nodes.TriggerTransition{Name: "new_tbl", IsNew: true},
	}}
	err := c.CreateTriggerStmt(stmt)
	if err != nil {
		t.Fatal(err)
	}
}

// -----------------------------------------------------------------------------
// trigger.go: Partitioned table ROW trigger with transition tables
// -----------------------------------------------------------------------------

func TestTriggerPartitionedRowTransitionError(t *testing.T) {
	c := New()
	createTriggerFunc(t, c)
	// Create partitioned table.
	pt := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "pt"},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "x", Type: TypeName{Name: "int4", TypeMod: -1}}),
		}},
		Partspec: &nodes.PartitionSpec{
			Strategy: "range",
			PartParams: &nodes.List{Items: []nodes.Node{
				&nodes.PartitionElem{Name: "x"},
			}},
		},
	}
	if err := c.DefineRelation(pt, 'r'); err != nil {
		t.Fatal(err)
	}

	stmt := makeCreateTrigStmt("", "pt", "trig1", TriggerAfter, TriggerEventInsert, true, "", "trig_fn")
	stmt.TransitionRels = &nodes.List{Items: []nodes.Node{
		&nodes.TriggerTransition{Name: "new_tbl", IsNew: true},
	}}
	err := c.CreateTriggerStmt(stmt)
	assertCode(t, err, CodeFeatureNotSupported)
	if !strings.Contains(err.Error(), "partitioned table") {
		t.Errorf("unexpected message: %s", err)
	}
}

// -----------------------------------------------------------------------------
// functioncmds.go: COST/ROWS validation
// -----------------------------------------------------------------------------

func TestFuncCostZeroError(t *testing.T) {
	c := New()
	stmt := makeCreateFuncStmt("", "f",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		TypeName{Name: "int4", TypeMod: -1},
		"sql", "SELECT $1", false)
	stmt.Options.Items = append(stmt.Options.Items,
		&nodes.DefElem{Defname: "cost", Arg: &nodes.Integer{Ival: 0}})
	err := c.CreateFunctionStmt(stmt)
	assertCode(t, err, CodeInvalidParameterValue)
	if !strings.Contains(err.Error(), "COST must be positive") {
		t.Errorf("unexpected message: %s", err)
	}
}

func TestFuncCostNegativeError(t *testing.T) {
	c := New()
	stmt := makeCreateFuncStmt("", "f",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		TypeName{Name: "int4", TypeMod: -1},
		"sql", "SELECT $1", false)
	stmt.Options.Items = append(stmt.Options.Items,
		&nodes.DefElem{Defname: "cost", Arg: &nodes.Integer{Ival: -5}})
	err := c.CreateFunctionStmt(stmt)
	assertCode(t, err, CodeInvalidParameterValue)
}

func TestFuncRowsZeroError(t *testing.T) {
	c := New()
	stmt := makeCreateFuncStmt("", "f",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		TypeName{Name: "int4", TypeMod: -1},
		"sql", "SELECT $1", false)
	stmt.Options.Items = append(stmt.Options.Items,
		&nodes.DefElem{Defname: "rows", Arg: &nodes.Integer{Ival: 0}})
	err := c.CreateFunctionStmt(stmt)
	assertCode(t, err, CodeInvalidParameterValue)
	if !strings.Contains(err.Error(), "ROWS must be positive") {
		t.Errorf("unexpected message: %s", err)
	}
}

func TestFuncCostPositiveOK(t *testing.T) {
	c := New()
	stmt := makeCreateFuncStmt("", "f",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		TypeName{Name: "int4", TypeMod: -1},
		"sql", "SELECT $1", false)
	stmt.Options.Items = append(stmt.Options.Items,
		&nodes.DefElem{Defname: "cost", Arg: &nodes.Integer{Ival: 100}})
	err := c.CreateFunctionStmt(stmt)
	if err != nil {
		t.Fatal(err)
	}
}

// -----------------------------------------------------------------------------
// functioncmds.go: Language validation
// -----------------------------------------------------------------------------

func TestFuncInvalidLanguageError(t *testing.T) {
	c := New()
	stmt := makeCreateFuncStmt("", "f",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		TypeName{Name: "int4", TypeMod: -1},
		"nonexistent_lang", "SELECT $1", false)
	err := c.CreateFunctionStmt(stmt)
	assertCode(t, err, CodeUndefinedObject)
	if !strings.Contains(err.Error(), "language") {
		t.Errorf("unexpected message: %s", err)
	}
}

func TestFuncPlpgsqlLanguageOK(t *testing.T) {
	c := New()
	stmt := makeCreateFuncStmt("", "f",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		TypeName{Name: "int4", TypeMod: -1},
		"plpgsql", "BEGIN RETURN $1; END;", false)
	err := c.CreateFunctionStmt(stmt)
	if err != nil {
		t.Fatal(err)
	}
}

// -----------------------------------------------------------------------------
// functioncmds.go: ALTER FUNCTION COST/ROWS validation
// -----------------------------------------------------------------------------

func TestAlterFuncCostZeroError(t *testing.T) {
	c := New()
	// Create function first.
	stmt := makeCreateFuncStmt("", "f",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		TypeName{Name: "int4", TypeMod: -1},
		"sql", "SELECT $1", false)
	if err := c.CreateFunctionStmt(stmt); err != nil {
		t.Fatal(err)
	}
	// Alter with COST 0.
	alter := &nodes.AlterFunctionStmt{
		Objtype: nodes.OBJECT_FUNCTION,
		Func: &nodes.ObjectWithArgs{
			Objname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "f"}}},
			Objargs: &nodes.List{Items: []nodes.Node{
				makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
			}},
		},
		Actions: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{Defname: "cost", Arg: &nodes.Integer{Ival: 0}},
		}},
	}
	err := c.AlterFunction(alter)
	assertCode(t, err, CodeInvalidParameterValue)
}

func TestAlterFuncRowsNegativeError(t *testing.T) {
	c := New()
	stmt := makeCreateFuncStmt("", "f",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		TypeName{Name: "int4", TypeMod: -1},
		"sql", "SELECT $1", false)
	if err := c.CreateFunctionStmt(stmt); err != nil {
		t.Fatal(err)
	}
	alter := &nodes.AlterFunctionStmt{
		Objtype: nodes.OBJECT_FUNCTION,
		Func: &nodes.ObjectWithArgs{
			Objname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "f"}}},
			Objargs: &nodes.List{Items: []nodes.Node{
				makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
			}},
		},
		Actions: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{Defname: "rows", Arg: &nodes.Integer{Ival: -1}},
		}},
	}
	err := c.AlterFunction(alter)
	assertCode(t, err, CodeInvalidParameterValue)
}

// -----------------------------------------------------------------------------
// dropcmds.go: DROP FUNCTION on procedure / DROP PROCEDURE on function
// -----------------------------------------------------------------------------

func TestDropFunctionOnProcedureError(t *testing.T) {
	c := New()
	// Create procedure (no return type).
	stmt := markCreateFunctionStmtAsProcedure(&nodes.CreateFunctionStmt{
		Funcname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "myproc"}}},
		Options: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{Defname: "language", Arg: &nodes.String{Str: "sql"}},
			&nodes.DefElem{Defname: "as", Arg: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "NULL"}}}},
		}},
	})
	if err := c.CreateFunctionStmt(stmt); err != nil {
		t.Fatal(err)
	}
	// DROP FUNCTION should fail for a procedure.
	drop := &nodes.DropStmt{
		Objects: &nodes.List{Items: []nodes.Node{
			&nodes.ObjectWithArgs{
				Objname:         &nodes.List{Items: []nodes.Node{&nodes.String{Str: "myproc"}}},
				ArgsUnspecified: true,
			},
		}},
		RemoveType: int(nodes.OBJECT_FUNCTION),
		Behavior:   int(nodes.DROP_RESTRICT),
	}
	err := c.RemoveObjects(drop)
	assertCode(t, err, CodeWrongObjectType)
	if !strings.Contains(err.Error(), "not a function") {
		t.Errorf("unexpected message: %s", err)
	}
}

func TestDropProcedureOnFunctionError(t *testing.T) {
	c := New()
	// Create function.
	stmt := makeCreateFuncStmt("", "myfunc",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		TypeName{Name: "int4", TypeMod: -1},
		"sql", "SELECT $1", false)
	if err := c.CreateFunctionStmt(stmt); err != nil {
		t.Fatal(err)
	}
	// DROP PROCEDURE should fail for a function.
	drop := &nodes.DropStmt{
		Objects: &nodes.List{Items: []nodes.Node{
			&nodes.ObjectWithArgs{
				Objname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "myfunc"}}},
				Objargs: &nodes.List{Items: []nodes.Node{
					makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
				}},
			},
		}},
		RemoveType: int(nodes.OBJECT_PROCEDURE),
		Behavior:   int(nodes.DROP_RESTRICT),
	}
	err := c.RemoveObjects(drop)
	assertCode(t, err, CodeWrongObjectType)
	if !strings.Contains(err.Error(), "not a procedure") {
		t.Errorf("unexpected message: %s", err)
	}
}

// -----------------------------------------------------------------------------
// aclchk.go: WITH GRANT OPTION to PUBLIC
// -----------------------------------------------------------------------------

func TestGrantWithGrantOptionToPublicError(t *testing.T) {
	c := New()
	// Create table first.
	if err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r'); err != nil {
		t.Fatal(err)
	}
	// GRANT SELECT ON t TO PUBLIC WITH GRANT OPTION — should fail.
	stmt := &nodes.GrantStmt{
		IsGrant:     true,
		GrantOption: true,
		Objtype:     nodes.OBJECT_TABLE,
		Objects: &nodes.List{Items: []nodes.Node{
			&nodes.RangeVar{Relname: "t"},
		}},
		Privileges: &nodes.List{Items: []nodes.Node{
			&nodes.AccessPriv{PrivName: "select"},
		}},
		Grantees: &nodes.List{Items: []nodes.Node{
			&nodes.RoleSpec{Roletype: int(nodes.ROLESPEC_PUBLIC)},
		}},
	}
	err := c.ExecGrantStmt(stmt)
	assertCode(t, err, CodeInvalidGrantOperation)
	if !strings.Contains(err.Error(), "grant options can only be granted to roles") {
		t.Errorf("unexpected message: %s", err)
	}
}

func TestGrantWithGrantOptionToRoleOK(t *testing.T) {
	c := New()
	if err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r'); err != nil {
		t.Fatal(err)
	}
	// GRANT SELECT ON t TO some_role WITH GRANT OPTION — should succeed.
	stmt := &nodes.GrantStmt{
		IsGrant:     true,
		GrantOption: true,
		Objtype:     nodes.OBJECT_TABLE,
		Objects: &nodes.List{Items: []nodes.Node{
			&nodes.RangeVar{Relname: "t"},
		}},
		Privileges: &nodes.List{Items: []nodes.Node{
			&nodes.AccessPriv{PrivName: "select"},
		}},
		Grantees: &nodes.List{Items: []nodes.Node{
			&nodes.RoleSpec{Roletype: int(nodes.ROLESPEC_CSTRING), Rolename: "some_role"},
		}},
	}
	err := c.ExecGrantStmt(stmt)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGrantToPublicWithoutGrantOptionOK(t *testing.T) {
	c := New()
	if err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r'); err != nil {
		t.Fatal(err)
	}
	// GRANT SELECT ON t TO PUBLIC (no WITH GRANT OPTION) — should succeed.
	stmt := &nodes.GrantStmt{
		IsGrant: true,
		Objtype: nodes.OBJECT_TABLE,
		Objects: &nodes.List{Items: []nodes.Node{
			&nodes.RangeVar{Relname: "t"},
		}},
		Privileges: &nodes.List{Items: []nodes.Node{
			&nodes.AccessPriv{PrivName: "select"},
		}},
		Grantees: &nodes.List{Items: []nodes.Node{
			&nodes.RoleSpec{Roletype: int(nodes.ROLESPEC_PUBLIC)},
		}},
	}
	err := c.ExecGrantStmt(stmt)
	if err != nil {
		t.Fatal(err)
	}
}

// -----------------------------------------------------------------------------
// sequence.go: Owner/schema consistency
// -----------------------------------------------------------------------------

func TestSequenceOwnedByDifferentSchemaError(t *testing.T) {
	c := New()
	// Create schema s1.
	c.CreateSchemaCommand(&nodes.CreateSchemaStmt{Schemaname: "s1"})

	// Create table in s1 (so lookup succeeds).
	tblStmt := makeCreateTableStmt("s1", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false)
	if err := c.DefineRelation(tblStmt, 'r'); err != nil {
		t.Fatal(err)
	}

	// Create sequence in public schema (different from s1).
	seqStmt := &nodes.CreateSeqStmt{
		Sequence: &nodes.RangeVar{Relname: "seq1"},
	}
	if err := c.DefineSequence(seqStmt); err != nil {
		t.Fatal(err)
	}

	// Directly call setSequenceOwner to test the cross-schema check.
	// The sequence is in public, the table is in s1 — should fail.
	seq, _ := c.findSequence("", "seq1")
	s1 := c.schemaByName["s1"]
	err := c.setSequenceOwner(seq, s1, "t.id")
	assertCode(t, err, CodeFeatureNotSupported)
	if !strings.Contains(err.Error(), "same schema") {
		t.Errorf("unexpected message: %s", err)
	}
}

func TestSequenceOwnedBySameSchemaOK(t *testing.T) {
	c := New()
	// Create table and sequence in same schema (public).
	if err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r'); err != nil {
		t.Fatal(err)
	}
	seqStmt := &nodes.CreateSeqStmt{
		Sequence: &nodes.RangeVar{Relname: "seq1"},
	}
	if err := c.DefineSequence(seqStmt); err != nil {
		t.Fatal(err)
	}
	alterStmt := &nodes.AlterSeqStmt{
		Sequence: &nodes.RangeVar{Relname: "seq1"},
		Options: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{
				Defname: "owned_by",
				Arg: &nodes.List{Items: []nodes.Node{
					&nodes.String{Str: "t"},
					&nodes.String{Str: "id"},
				}},
			},
		}},
	}
	err := c.AlterSequenceStmt(alterStmt)
	if err != nil {
		t.Fatal(err)
	}
}

// -----------------------------------------------------------------------------
// view.go: CHECK OPTION on set-op views
// -----------------------------------------------------------------------------

func TestViewCheckOptionOnSetOpError(t *testing.T) {
	c := New()
	// Create base tables.
	if err := c.DefineRelation(makeCreateTableStmt("", "t1", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r'); err != nil {
		t.Fatal(err)
	}
	if err := c.DefineRelation(makeCreateTableStmt("", "t2", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r'); err != nil {
		t.Fatal(err)
	}
	// CREATE VIEW v WITH CASCADED CHECK OPTION AS SELECT id FROM t1 UNION SELECT id FROM t2
	viewStmt := &nodes.ViewStmt{
		View: &nodes.RangeVar{Relname: "v"},
		Query: &nodes.SelectStmt{
			Op: nodes.SETOP_UNION,
			Larg: &nodes.SelectStmt{
				TargetList: &nodes.List{Items: []nodes.Node{
					&nodes.ResTarget{Name: "id", Val: &nodes.ColumnRef{Fields: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "id"}}}}},
				}},
				FromClause: &nodes.List{Items: []nodes.Node{
					&nodes.RangeVar{Relname: "t1"},
				}},
			},
			Rarg: &nodes.SelectStmt{
				TargetList: &nodes.List{Items: []nodes.Node{
					&nodes.ResTarget{Name: "id", Val: &nodes.ColumnRef{Fields: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "id"}}}}},
				}},
				FromClause: &nodes.List{Items: []nodes.Node{
					&nodes.RangeVar{Relname: "t2"},
				}},
			},
		},
		WithCheckOption: 2, // CASCADED
	}
	err := c.DefineView(viewStmt)
	assertCode(t, err, CodeFeatureNotSupported)
	if !strings.Contains(err.Error(), "CHECK OPTION") {
		t.Errorf("unexpected message: %s", err)
	}
}

// -----------------------------------------------------------------------------
// indexcmds.go: Unique index on partitioned table must include partition key
// -----------------------------------------------------------------------------

func TestUniqueIndexPartitionedMissingKeyError(t *testing.T) {
	c := New()
	// Create partitioned table with partition key on 'b'.
	pt := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "pt"},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}}),
			makeColumnDefNode(ColumnDef{Name: "b", Type: TypeName{Name: "int4", TypeMod: -1}}),
		}},
		Partspec: &nodes.PartitionSpec{
			Strategy: "range",
			PartParams: &nodes.List{Items: []nodes.Node{
				&nodes.PartitionElem{Name: "b"},
			}},
		},
	}
	if err := c.DefineRelation(pt, 'r'); err != nil {
		t.Fatal(err)
	}
	// CREATE UNIQUE INDEX on only column 'a' — should fail because 'b' is partition key.
	idxStmt := &nodes.IndexStmt{
		Relation: &nodes.RangeVar{Relname: "pt"},
		Idxname:  "pt_a_idx",
		Unique:   true,
		IndexParams: &nodes.List{Items: []nodes.Node{
			&nodes.IndexElem{Name: "a"},
		}},
	}
	err := c.DefineIndex(idxStmt)
	assertCode(t, err, CodeFeatureNotSupported)
	if !strings.Contains(err.Error(), "partitioning columns") {
		t.Errorf("unexpected message: %s", err)
	}
}

func TestUniqueIndexPartitionedIncludesKeyOK(t *testing.T) {
	c := New()
	pt := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "pt"},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}}),
			makeColumnDefNode(ColumnDef{Name: "b", Type: TypeName{Name: "int4", TypeMod: -1}}),
		}},
		Partspec: &nodes.PartitionSpec{
			Strategy: "range",
			PartParams: &nodes.List{Items: []nodes.Node{
				&nodes.PartitionElem{Name: "b"},
			}},
		},
	}
	if err := c.DefineRelation(pt, 'r'); err != nil {
		t.Fatal(err)
	}
	// CREATE UNIQUE INDEX on (a, b) — includes partition key 'b', should succeed.
	idxStmt := &nodes.IndexStmt{
		Relation: &nodes.RangeVar{Relname: "pt"},
		Idxname:  "pt_ab_idx",
		Unique:   true,
		IndexParams: &nodes.List{Items: []nodes.Node{
			&nodes.IndexElem{Name: "a"},
			&nodes.IndexElem{Name: "b"},
		}},
	}
	err := c.DefineIndex(idxStmt)
	if err != nil {
		t.Fatal(err)
	}
}
