package ddl

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// =============================================================================
// Phase 1: tablecmds.go validations
// =============================================================================

func TestOnCommitNonTempError(t *testing.T) {
	c := New()
	stmt := makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false)
	stmt.OnCommit = nodes.ONCOMMIT_DELETE_ROWS
	err := c.DefineRelation(stmt, 'r')
	assertErrorCode(t, err, CodeInvalidTableDefinition)
}

func TestOnCommitTempOK(t *testing.T) {
	c := New()
	stmt := makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false)
	stmt.OnCommit = nodes.ONCOMMIT_DELETE_ROWS
	stmt.Relation.Relpersistence = 't'
	err := c.DefineRelation(stmt, 'r')
	if err != nil {
		t.Fatal(err)
	}
}

func TestDuplicateParentInherits(t *testing.T) {
	c := New()
	// Create parent.
	c.DefineRelation(makeCreateTableStmt("", "p", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	// Create child inheriting from same parent twice.
	child := makeCreateTableStmt("", "ch", nil, nil, false)
	child.InhRelations = &nodes.List{Items: []nodes.Node{
		&nodes.RangeVar{Relname: "p"},
		&nodes.RangeVar{Relname: "p"},
	}}
	err := c.DefineRelation(child, 'r')
	assertErrorCode(t, err, CodeDuplicateObject)
}

func TestListPartitionMultiColumn(t *testing.T) {
	c := New()

	stmt := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "t"},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}}),
			makeColumnDefNode(ColumnDef{Name: "b", Type: TypeName{Name: "int4", TypeMod: -1}}),
		}},
		Partspec: &nodes.PartitionSpec{
			Strategy: "list",
			PartParams: &nodes.List{Items: []nodes.Node{
				&nodes.PartitionElem{Name: "a"},
				&nodes.PartitionElem{Name: "b"},
			}},
		},
	}
	err := c.DefineRelation(stmt, 'r')
	assertErrorCode(t, err, CodeInvalidObjectDefinition)
}

func TestListPartitionSingleColumnOK(t *testing.T) {
	c := New()

	stmt := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "t"},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}}),
		}},
		Partspec: &nodes.PartitionSpec{
			Strategy: "list",
			PartParams: &nodes.List{Items: []nodes.Node{
				&nodes.PartitionElem{Name: "a"},
			}},
		},
	}
	err := c.DefineRelation(stmt, 'r')
	if err != nil {
		t.Fatal(err)
	}
}

// =============================================================================
// Phase 2: typecmds.go — AlterEnum RENAME VALUE
// =============================================================================

func TestAlterEnumRenameValue(t *testing.T) {
	c := New()
	c.DefineEnum(&nodes.CreateEnumStmt{
		TypeName: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "color"}}},
		Vals:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: "red"}, &nodes.String{Str: "green"}}},
	})

	// Rename "red" → "blue".
	err := c.AlterEnumStmt(&nodes.AlterEnumStmt{
		Typname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "color"}}},
		Oldval:   "red",
		Newval:   "blue",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Verify "blue" exists and "red" doesn't.
	et := c.enumTypes[c.typeByName[typeKey{ns: PublicNamespace, name: "color"}].OID]
	if _, ok := et.labelMap["blue"]; !ok {
		t.Error("expected 'blue' in enum labels")
	}
	if _, ok := et.labelMap["red"]; ok {
		t.Error("expected 'red' to be removed from enum labels")
	}
}

func TestAlterEnumRenameValueDuplicate(t *testing.T) {
	c := New()
	c.DefineEnum(&nodes.CreateEnumStmt{
		TypeName: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "color"}}},
		Vals:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: "red"}, &nodes.String{Str: "green"}}},
	})

	err := c.AlterEnumStmt(&nodes.AlterEnumStmt{
		Typname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "color"}}},
		Oldval:   "red",
		Newval:   "green", // already exists
	})
	assertErrorCode(t, err, CodeDuplicateObject)
}

func TestAlterEnumRenameValueNotFound(t *testing.T) {
	c := New()
	c.DefineEnum(&nodes.CreateEnumStmt{
		TypeName: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "color"}}},
		Vals:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: "red"}}},
	})

	err := c.AlterEnumStmt(&nodes.AlterEnumStmt{
		Typname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "color"}}},
		Oldval:   "nosuch",
		Newval:   "blue",
	})
	assertErrorCode(t, err, CodeInvalidParameterValue)
}

// =============================================================================
// Phase 2: typecmds.go — DefineDomain multiple CHECK + constraint type rejection
// =============================================================================

func TestDefineDomainMultipleChecks(t *testing.T) {
	c := New()
	err := c.DefineDomain(&nodes.CreateDomainStmt{
		Domainname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "posint"}}},
		Typname:    makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
		Constraints: &nodes.List{Items: []nodes.Node{
			&nodes.Constraint{Contype: nodes.CONSTR_CHECK, Conname: "check_pos", CookedExpr: "VALUE > 0"},
			&nodes.Constraint{Contype: nodes.CONSTR_CHECK, Conname: "check_lt100", CookedExpr: "VALUE < 100"},
		}},
	})
	if err != nil {
		t.Fatal(err)
	}

	dt := c.domainTypes[c.typeByName[typeKey{ns: PublicNamespace, name: "posint"}].OID]
	if len(dt.Constraints) != 2 {
		t.Fatalf("expected 2 constraints, got %d", len(dt.Constraints))
	}
	if dt.Constraints[0].Name != "check_pos" || dt.Constraints[1].Name != "check_lt100" {
		t.Errorf("constraint names mismatch: %q, %q", dt.Constraints[0].Name, dt.Constraints[1].Name)
	}
}

func TestDefineDomainRejectsUnique(t *testing.T) {
	c := New()
	err := c.DefineDomain(&nodes.CreateDomainStmt{
		Domainname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "d"}}},
		Typname:    makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
		Constraints: &nodes.List{Items: []nodes.Node{
			&nodes.Constraint{Contype: nodes.CONSTR_UNIQUE},
		}},
	})
	assertErrorCode(t, err, CodeSyntaxError)
}

func TestDefineDomainRejectsPrimary(t *testing.T) {
	c := New()
	err := c.DefineDomain(&nodes.CreateDomainStmt{
		Domainname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "d"}}},
		Typname:    makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
		Constraints: &nodes.List{Items: []nodes.Node{
			&nodes.Constraint{Contype: nodes.CONSTR_PRIMARY},
		}},
	})
	assertErrorCode(t, err, CodeSyntaxError)
}

func TestAlterDomainAddNonCheckError(t *testing.T) {
	c := New()
	c.DefineDomain(&nodes.CreateDomainStmt{
		Domainname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "d"}}},
		Typname:    makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
	})

	err := c.AlterDomainStmt(&nodes.AlterDomainStmt{
		Typname:  &nodes.List{Items: []nodes.Node{&nodes.String{Str: "d"}}},
		Subtype:  'C',
		Def:      &nodes.Constraint{Contype: nodes.CONSTR_UNIQUE},
	})
	assertErrorCode(t, err, CodeSyntaxError)
}

// =============================================================================
// Phase 3: Trigger validation
// =============================================================================

func TestTriggerRowTruncateRejected(t *testing.T) {
	c := New()
	createTestTable(t, c)
	createTriggerFunc(t, c)

	err := c.CreateTriggerStmt(makeCreateTrigStmt("", "t", "my_trig", TriggerBefore, TriggerEventTruncate, true, "", "trig_fn"))
	assertCode(t, err, CodeInvalidObjectDefinition)
}

func TestTriggerInsteadOfMustBeRow(t *testing.T) {
	c := New()
	createTestTable(t, c)
	createTriggerFunc(t, c)

	c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			From:       []*FromItem{{Kind: FromTable, Table: "t"}},
			TargetList: []*ResTarget{{Val: &Expr{Kind: ExprStar}}},
		}, nil))

	// INSTEAD OF with forEachRow=false should fail.
	err := c.CreateTriggerStmt(makeCreateTrigStmt("", "v", "my_trig", TriggerInsteadOf, TriggerEventInsert, false, "", "trig_fn"))
	assertCode(t, err, CodeInvalidObjectDefinition)
}

func TestTriggerInsteadOfNoWhen(t *testing.T) {
	c := New()
	createTestTable(t, c)
	createTriggerFunc(t, c)

	c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			From:       []*FromItem{{Kind: FromTable, Table: "t"}},
			TargetList: []*ResTarget{{Val: &Expr{Kind: ExprStar}}},
		}, nil))

	// Build INSTEAD OF with WHEN clause.
	stmt := makeCreateTrigStmt("", "v", "my_trig", TriggerInsteadOf, TriggerEventInsert, true, "", "trig_fn")
	stmt.WhenClause = &nodes.String{Str: "true"}
	err := c.CreateTriggerStmt(stmt)
	assertCode(t, err, CodeInvalidObjectDefinition)
}

func TestTriggerInsteadOfNoColumns(t *testing.T) {
	c := New()
	createTestTable(t, c)
	createTriggerFunc(t, c)

	c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			From:       []*FromItem{{Kind: FromTable, Table: "t"}},
			TargetList: []*ResTarget{{Val: &Expr{Kind: ExprStar}}},
		}, nil))

	// Build INSTEAD OF with column list.
	stmt := makeCreateTrigStmtWithColumns("", "v", "my_trig", TriggerInsteadOf,
		TriggerEventUpdate, true, "", "trig_fn", []string{"x"})
	err := c.CreateTriggerStmt(stmt)
	assertCode(t, err, CodeInvalidObjectDefinition)
}

func TestTriggerTransitionTableAfterOnly(t *testing.T) {
	c := New()
	createTestTable(t, c)
	createTriggerFunc(t, c)

	stmt := makeCreateTrigStmt("", "t", "my_trig", TriggerBefore, TriggerEventInsert, false, "", "trig_fn")
	stmt.TransitionRels = &nodes.List{Items: []nodes.Node{
		&nodes.TriggerTransition{IsNew: true, Name: "new_t"},
	}}
	err := c.CreateTriggerStmt(stmt)
	assertCode(t, err, CodeInvalidObjectDefinition)
}

func TestTriggerTransitionTableMultiEventError(t *testing.T) {
	c := New()
	createTestTable(t, c)
	createTriggerFunc(t, c)

	stmt := makeCreateTrigStmt("", "t", "my_trig", TriggerAfter, TriggerEventInsert|TriggerEventUpdate, false, "", "trig_fn")
	stmt.TransitionRels = &nodes.List{Items: []nodes.Node{
		&nodes.TriggerTransition{IsNew: true, Name: "new_t"},
	}}
	err := c.CreateTriggerStmt(stmt)
	assertCode(t, err, CodeInvalidObjectDefinition)
}

func TestTriggerTransitionTableOK(t *testing.T) {
	c := New()
	createTestTable(t, c)
	createTriggerFunc(t, c)

	stmt := makeCreateTrigStmt("", "t", "my_trig", TriggerAfter, TriggerEventInsert, false, "", "trig_fn")
	stmt.TransitionRels = &nodes.List{Items: []nodes.Node{
		&nodes.TriggerTransition{IsNew: true, Name: "new_t"},
	}}
	err := c.CreateTriggerStmt(stmt)
	if err != nil {
		t.Fatal(err)
	}

	// Verify transition table name stored.
	trigs := c.TriggersOf(c.GetRelation("", "t").OID)
	if len(trigs) == 0 {
		t.Fatal("expected trigger")
	}
	if trigs[0].NewTransitionName != "new_t" {
		t.Errorf("new transition name: got %q, want %q", trigs[0].NewTransitionName, "new_t")
	}
}

func TestTriggerStatementBeforeOnViewOK(t *testing.T) {
	c := New()
	createTestTable(t, c)
	createTriggerFunc(t, c)

	c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			From:       []*FromItem{{Kind: FromTable, Table: "t"}},
			TargetList: []*ResTarget{{Val: &Expr{Kind: ExprStar}}},
		}, nil))

	// STATEMENT-level BEFORE triggers are OK on views.
	err := c.CreateTriggerStmt(makeCreateTrigStmt("", "v", "my_trig", TriggerBefore, TriggerEventInsert, false, "", "trig_fn"))
	if err != nil {
		t.Fatal(err)
	}
}

// =============================================================================
// Phase 4: Policy relkind check
// =============================================================================

func TestCreatePolicyOnViewError(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			From:       []*FromItem{{Kind: FromTable, Table: "t"}},
			TargetList: []*ResTarget{{Val: &Expr{Kind: ExprStar}}},
		}, nil))

	err := c.CreatePolicy(&nodes.CreatePolicyStmt{
		Table:      &nodes.RangeVar{Relname: "v"},
		PolicyName: "pol",
	})
	assertErrorCode(t, err, CodeWrongObjectType)
}

func TestCreatePolicyOnTableOK(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	err := c.CreatePolicy(&nodes.CreatePolicyStmt{
		Table:      &nodes.RangeVar{Relname: "t"},
		PolicyName: "pol",
		Permissive: true,
	})
	if err != nil {
		t.Fatal(err)
	}
}

// =============================================================================
// Phase 4: DROP POLICY
// =============================================================================

func TestDropPolicy(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	c.CreatePolicy(&nodes.CreatePolicyStmt{
		Table: &nodes.RangeVar{Relname: "t"}, PolicyName: "pol", Permissive: true,
	})

	err := c.RemoveObjects(&nodes.DropStmt{
		RemoveType: int(nodes.OBJECT_POLICY),
		Objects: &nodes.List{Items: []nodes.Node{
			&nodes.List{Items: []nodes.Node{&nodes.String{Str: "t"}, &nodes.String{Str: "pol"}}},
		}},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Verify policy is gone.
	rel := c.GetRelation("", "t")
	if len(c.policiesByRel[rel.OID]) != 0 {
		t.Error("expected policy to be removed")
	}
}

func TestDropPolicyIfExistsNotFound(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	err := c.RemoveObjects(&nodes.DropStmt{
		RemoveType: int(nodes.OBJECT_POLICY),
		Missing_ok: true,
		Objects: &nodes.List{Items: []nodes.Node{
			&nodes.List{Items: []nodes.Node{&nodes.String{Str: "t"}, &nodes.String{Str: "nope"}}},
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDropPolicyNotFoundError(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	err := c.RemoveObjects(&nodes.DropStmt{
		RemoveType: int(nodes.OBJECT_POLICY),
		Objects: &nodes.List{Items: []nodes.Node{
			&nodes.List{Items: []nodes.Node{&nodes.String{Str: "t"}, &nodes.String{Str: "nope"}}},
		}},
	})
	assertErrorCode(t, err, CodeUndefinedObject)
}

// =============================================================================
// Phase 4: CREATE SCHEMA with schemaElts
// =============================================================================

func TestCreateSchemaWithElements(t *testing.T) {
	c := New()

	err := c.ProcessUtility(&nodes.CreateSchemaStmt{
		Schemaname: "myschema",
		SchemaElts: &nodes.List{Items: []nodes.Node{
			makeCreateTableStmt("", "t1", []ColumnDef{
				{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
			}, nil, false),
		}},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Verify table exists in myschema.
	r := c.GetRelation("myschema", "t1")
	if r == nil {
		t.Fatal("expected table myschema.t1 to exist")
	}
}

// =============================================================================
// Phase 4: Utility no-ops
// =============================================================================

func TestCreateExtensionNoOp(t *testing.T) {
	c := New()
	err := c.ProcessUtility(&nodes.CreateExtensionStmt{Extname: "pg_trgm"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestReindexNoOp(t *testing.T) {
	c := New()
	err := c.ProcessUtility(&nodes.ReindexStmt{Name: "foo"})
	if err != nil {
		t.Fatal(err)
	}
}

// =============================================================================
// Phase 5: ALTER TABLE objtype relkind
// =============================================================================

func TestAlterTableOnForeignTable(t *testing.T) {
	c := New()
	// Create foreign table.
	c.DefineRelation(makeCreateTableStmt("", "ft", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'f')

	// ALTER TABLE with OBJECT_TABLE objtype should work on foreign tables.
	stmt := makeAlterTableStmt("", "ft",
		makeATSetNotNull("id"),
	)
	stmt.ObjType = int(nodes.OBJECT_TABLE)
	err := c.AlterTableStmt(stmt)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAlterViewOnViewOK(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			From:       []*FromItem{{Kind: FromTable, Table: "t"}},
			TargetList: []*ResTarget{{Val: &Expr{Kind: ExprStar}}},
		}, nil))

	// ALTER VIEW on a view should work.
	stmt := makeAlterTableStmt("", "v",
		makeATSetNotNull("id"),
	)
	stmt.ObjType = int(nodes.OBJECT_VIEW)
	err := c.AlterTableStmt(stmt)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAlterTableObjTypeMismatch(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	// ALTER VIEW on a regular table should fail.
	stmt := makeAlterTableStmt("", "t",
		makeATSetNotNull("id"),
	)
	stmt.ObjType = int(nodes.OBJECT_VIEW)
	err := c.AlterTableStmt(stmt)
	assertErrorCode(t, err, CodeWrongObjectType)
}

// =============================================================================
// Phase 5: RENAME FOREIGN TABLE
// =============================================================================

func TestRenameForeignTable(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "ft", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'f')

	err := c.ExecRenameStmt(&nodes.RenameStmt{
		RenameType: nodes.OBJECT_FOREIGN_TABLE,
		Relation:   &nodes.RangeVar{Relname: "ft"},
		Newname:    "ft2",
	})
	if err != nil {
		t.Fatal(err)
	}

	if r := c.GetRelation("", "ft2"); r == nil {
		t.Error("expected foreign table ft2 to exist")
	}
}

// =============================================================================
// Phase 5: FK persistence cross-reference
// =============================================================================

func TestFKPermanentToTempError(t *testing.T) {
	c := New()
	// Create a permanent table.
	c.DefineRelation(makeCreateTableStmt("", "perm", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, []ConstraintDef{
		{Type: ConstraintPK, Columns: []string{"id"}},
	}, false), 'r')

	// Create a temp table.
	tempStmt := makeCreateTableStmt("", "tmp", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "ref", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false)
	tempStmt.Relation.Relpersistence = 't'
	c.DefineRelation(tempStmt, 'r')

	// Try to add FK from permanent → temp. Should fail.
	stmt := makeAlterTableStmt("", "perm",
		makeATAddConstraint(ConstraintDef{
			Type:       ConstraintFK,
			Columns:    []string{"id"},
			RefTable:   "tmp",
			RefColumns: []string{"id"},
		}),
	)
	err := c.AlterTableStmt(stmt)
	assertErrorCode(t, err, CodeInvalidFK)
}

// =============================================================================
// Phase 6: Sequence validation
// =============================================================================

func TestSequenceCacheZeroError(t *testing.T) {
	c := New()
	err := c.DefineSequence(makeCreateSeqStmt("", "s", false,
		&nodes.DefElem{Defname: "cache", Arg: &nodes.Integer{Ival: 0}},
	))
	assertCode(t, err, CodeInvalidParameterValue)
}

func TestSequenceCacheNegativeError(t *testing.T) {
	c := New()
	err := c.DefineSequence(makeCreateSeqStmt("", "s", false,
		&nodes.DefElem{Defname: "cache", Arg: &nodes.Integer{Ival: -1}},
	))
	assertCode(t, err, CodeInvalidParameterValue)
}

func TestAlterSequenceRestart(t *testing.T) {
	c := New()
	c.DefineSequence(makeCreateSeqStmt("", "s", false))

	err := c.AlterSequenceStmt(makeAlterSeqStmt("", "s",
		&nodes.DefElem{Defname: "restart", Arg: &nodes.Integer{Ival: 42}},
	))
	if err != nil {
		t.Fatal(err)
	}

	seq, _ := c.findSequence("", "s")
	if seq.Start != 42 {
		t.Errorf("sequence start: got %d, want 42", seq.Start)
	}
}

func TestAlterSequenceCacheZeroError(t *testing.T) {
	c := New()
	c.DefineSequence(makeCreateSeqStmt("", "s", false))

	err := c.AlterSequenceStmt(makeAlterSeqStmt("", "s",
		&nodes.DefElem{Defname: "cache", Arg: &nodes.Integer{Ival: 0}},
	))
	assertCode(t, err, CodeInvalidParameterValue)
}

// =============================================================================
// ALTER TABLE SET SCHEMA for foreign table
// =============================================================================

func TestAlterForeignTableSchema(t *testing.T) {
	c := New()
	c.CreateSchemaCommand(&nodes.CreateSchemaStmt{Schemaname: "s2"})
	c.DefineRelation(makeCreateTableStmt("", "ft", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'f')

	err := c.ExecAlterObjectSchemaStmt(&nodes.AlterObjectSchemaStmt{
		ObjectType: nodes.OBJECT_FOREIGN_TABLE,
		Relation:   &nodes.RangeVar{Relname: "ft"},
		Newschema:  "s2",
	})
	if err != nil {
		t.Fatal(err)
	}

	if r := c.GetRelation("s2", "ft"); r == nil {
		t.Error("expected foreign table in s2")
	}
}
