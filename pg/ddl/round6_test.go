package ddl

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// =============================================================================
// Phase 1: Function/Procedure Validations (functioncmds.go)
// =============================================================================

func TestFuncSetofParamRejected(t *testing.T) {
	c := New()
	stmt := makeCreateFuncStmt("", "f",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		TypeName{Name: "int4", TypeMod: -1},
		"sql", "SELECT $1", false)
	// Mark first param as SETOF.
	fp := stmt.Parameters.Items[0].(*nodes.FunctionParameter)
	fp.ArgType.Setof = true
	err := c.CreateFunctionStmt(stmt)
	assertCode(t, err, CodeInvalidFunctionDefinition)
}

func TestProcSetofParamRejected(t *testing.T) {
	c := New()
	// Procedure with SETOF param.
	stmt := &nodes.CreateFunctionStmt{
		Funcname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "p"}}},
		Parameters: &nodes.List{Items: []nodes.Node{
			&nodes.FunctionParameter{
				ArgType: &nodes.TypeName{Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "int4"}}}, Setof: true},
				Mode:    nodes.FUNC_PARAM_IN,
			},
		}},
		Options: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{Defname: "language", Arg: &nodes.String{Str: "sql"}},
			&nodes.DefElem{Defname: "as", Arg: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "NULL"}}}},
		}},
	}
	err := c.CreateFunctionStmt(stmt)
	assertCode(t, err, CodeInvalidFunctionDefinition)
}

func TestFuncVariadicMustBeLast(t *testing.T) {
	c := New()
	// VARIADIC param followed by another IN param.
	stmt := &nodes.CreateFunctionStmt{
		Funcname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "f"}}},
		Parameters: &nodes.List{Items: []nodes.Node{
			&nodes.FunctionParameter{
				Name:    "a",
				ArgType: makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1, IsArray: true}),
				Mode:    nodes.FUNC_PARAM_VARIADIC,
			},
			&nodes.FunctionParameter{
				Name:    "b",
				ArgType: makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
				Mode:    nodes.FUNC_PARAM_IN,
			},
		}},
		ReturnType: makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
		Options: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{Defname: "language", Arg: &nodes.String{Str: "sql"}},
			&nodes.DefElem{Defname: "as", Arg: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "SELECT 1"}}}},
		}},
	}
	err := c.CreateFunctionStmt(stmt)
	assertCode(t, err, CodeInvalidFunctionDefinition)
}

func TestFuncVariadicMustBeArray(t *testing.T) {
	c := New()
	// VARIADIC param with non-array type.
	stmt := &nodes.CreateFunctionStmt{
		Funcname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "f"}}},
		Parameters: &nodes.List{Items: []nodes.Node{
			&nodes.FunctionParameter{
				Name:    "a",
				ArgType: makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}), // not an array
				Mode:    nodes.FUNC_PARAM_VARIADIC,
			},
		}},
		ReturnType: makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
		Options: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{Defname: "language", Arg: &nodes.String{Str: "sql"}},
			&nodes.DefElem{Defname: "as", Arg: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "SELECT 1"}}}},
		}},
	}
	err := c.CreateFunctionStmt(stmt)
	assertCode(t, err, CodeInvalidFunctionDefinition)
}

func TestFuncVariadicArrayOK(t *testing.T) {
	c := New()
	// VARIADIC param with array type should succeed.
	stmt := &nodes.CreateFunctionStmt{
		Funcname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "f"}}},
		Parameters: &nodes.List{Items: []nodes.Node{
			&nodes.FunctionParameter{
				Name:    "a",
				ArgType: makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1, IsArray: true}),
				Mode:    nodes.FUNC_PARAM_VARIADIC,
			},
		}},
		ReturnType: makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
		Options: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{Defname: "language", Arg: &nodes.String{Str: "sql"}},
			&nodes.DefElem{Defname: "as", Arg: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "SELECT 1"}}}},
		}},
	}
	err := c.CreateFunctionStmt(stmt)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFuncDuplicateParamName(t *testing.T) {
	c := New()
	stmt := &nodes.CreateFunctionStmt{
		Funcname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "f"}}},
		Parameters: &nodes.List{Items: []nodes.Node{
			&nodes.FunctionParameter{
				Name:    "x",
				ArgType: makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
				Mode:    nodes.FUNC_PARAM_IN,
			},
			&nodes.FunctionParameter{
				Name:    "x",
				ArgType: makeTypeNameNode(TypeName{Name: "text", TypeMod: -1}),
				Mode:    nodes.FUNC_PARAM_IN,
			},
		}},
		ReturnType: makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
		Options: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{Defname: "language", Arg: &nodes.String{Str: "sql"}},
			&nodes.DefElem{Defname: "as", Arg: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "SELECT 1"}}}},
		}},
	}
	err := c.CreateFunctionStmt(stmt)
	assertCode(t, err, CodeInvalidFunctionDefinition)
}

func TestFuncDuplicateNameINvsOUT_OK(t *testing.T) {
	c := New()
	// Same name for IN and OUT params is allowed.
	// With a single OUT param of type text, the return type must be text (not record).
	// pg: src/backend/commands/functioncmds.c — single OUT param determines return type
	stmt := &nodes.CreateFunctionStmt{
		Funcname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "f"}}},
		Parameters: &nodes.List{Items: []nodes.Node{
			&nodes.FunctionParameter{
				Name:    "x",
				ArgType: makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
				Mode:    nodes.FUNC_PARAM_IN,
			},
			&nodes.FunctionParameter{
				Name:    "x",
				ArgType: makeTypeNameNode(TypeName{Name: "text", TypeMod: -1}),
				Mode:    nodes.FUNC_PARAM_OUT,
			},
		}},
		ReturnType: makeTypeNameNode(TypeName{Name: "text", TypeMod: -1}),
		Options: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{Defname: "language", Arg: &nodes.String{Str: "sql"}},
			&nodes.DefElem{Defname: "as", Arg: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "SELECT 'x'"}}}},
		}},
	}
	err := c.CreateFunctionStmt(stmt)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFuncDefaultParamOrdering(t *testing.T) {
	c := New()
	// Param with default followed by param without default.
	stmt := &nodes.CreateFunctionStmt{
		Funcname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "f"}}},
		Parameters: &nodes.List{Items: []nodes.Node{
			&nodes.FunctionParameter{
				Name:    "a",
				ArgType: makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
				Mode:    nodes.FUNC_PARAM_IN,
				Defexpr: &nodes.Integer{Ival: 42},
			},
			&nodes.FunctionParameter{
				Name:    "b",
				ArgType: makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
				Mode:    nodes.FUNC_PARAM_IN,
				// No default
			},
		}},
		ReturnType: makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
		Options: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{Defname: "language", Arg: &nodes.String{Str: "sql"}},
			&nodes.DefElem{Defname: "as", Arg: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "SELECT 1"}}}},
		}},
	}
	err := c.CreateFunctionStmt(stmt)
	assertCode(t, err, CodeInvalidFunctionDefinition)
}

func TestFuncDefaultParamOrderingOK(t *testing.T) {
	c := New()
	// Both params with defaults — OK.
	stmt := &nodes.CreateFunctionStmt{
		Funcname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "f"}}},
		Parameters: &nodes.List{Items: []nodes.Node{
			&nodes.FunctionParameter{
				Name:    "a",
				ArgType: makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
				Mode:    nodes.FUNC_PARAM_IN,
				Defexpr: &nodes.Integer{Ival: 42},
			},
			&nodes.FunctionParameter{
				Name:    "b",
				ArgType: makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
				Mode:    nodes.FUNC_PARAM_IN,
				Defexpr: &nodes.Integer{Ival: 0},
			},
		}},
		ReturnType: makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
		Options: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{Defname: "language", Arg: &nodes.String{Str: "sql"}},
			&nodes.DefElem{Defname: "as", Arg: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "SELECT 1"}}}},
		}},
	}
	err := c.CreateFunctionStmt(stmt)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFuncNoBodyRejected(t *testing.T) {
	c := New()
	// No AS clause and no SqlBody.
	stmt := &nodes.CreateFunctionStmt{
		Funcname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "f"}}},
		Parameters: &nodes.List{Items: []nodes.Node{
			&nodes.FunctionParameter{
				ArgType: makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
				Mode:    nodes.FUNC_PARAM_IN,
			},
		}},
		ReturnType: makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
		Options: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{Defname: "language", Arg: &nodes.String{Str: "sql"}},
		}},
	}
	err := c.CreateFunctionStmt(stmt)
	assertCode(t, err, CodeInvalidFunctionDefinition)
}

func TestProcVolatilityRejected(t *testing.T) {
	c := New()
	// Procedure with IMMUTABLE should be rejected.
	stmt := &nodes.CreateFunctionStmt{
		Funcname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "p"}}},
		// No ReturnType → procedure
		Options: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{Defname: "language", Arg: &nodes.String{Str: "plpgsql"}},
			&nodes.DefElem{Defname: "as", Arg: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "BEGIN NULL; END;"}}}},
			&nodes.DefElem{Defname: "volatility", Arg: &nodes.String{Str: "immutable"}},
		}},
	}
	err := c.CreateFunctionStmt(stmt)
	assertCode(t, err, CodeInvalidFunctionDefinition)
}

func TestProcStrictRejected(t *testing.T) {
	c := New()
	stmt := &nodes.CreateFunctionStmt{
		Funcname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "p"}}},
		Options: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{Defname: "language", Arg: &nodes.String{Str: "plpgsql"}},
			&nodes.DefElem{Defname: "as", Arg: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "BEGIN NULL; END;"}}}},
			&nodes.DefElem{Defname: "strict", Arg: &nodes.Integer{Ival: 1}},
		}},
	}
	err := c.CreateFunctionStmt(stmt)
	assertCode(t, err, CodeInvalidFunctionDefinition)
}

func TestProcParallelRejected(t *testing.T) {
	c := New()
	stmt := &nodes.CreateFunctionStmt{
		Funcname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "p"}}},
		Options: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{Defname: "language", Arg: &nodes.String{Str: "plpgsql"}},
			&nodes.DefElem{Defname: "as", Arg: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "BEGIN NULL; END;"}}}},
			&nodes.DefElem{Defname: "parallel", Arg: &nodes.String{Str: "safe"}},
		}},
	}
	err := c.CreateFunctionStmt(stmt)
	assertCode(t, err, CodeInvalidFunctionDefinition)
}

func TestProcWindowRejected(t *testing.T) {
	c := New()
	stmt := &nodes.CreateFunctionStmt{
		Funcname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "p"}}},
		Options: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{Defname: "language", Arg: &nodes.String{Str: "plpgsql"}},
			&nodes.DefElem{Defname: "as", Arg: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "BEGIN NULL; END;"}}}},
			&nodes.DefElem{Defname: "window", Arg: &nodes.Integer{Ival: 1}},
		}},
	}
	err := c.CreateFunctionStmt(stmt)
	assertCode(t, err, CodeInvalidFunctionDefinition)
}

// =============================================================================
// Phase 2: Index AM Validations (indexcmds.go)
// =============================================================================

func TestIndexUniqueOnHashOK(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	stmt := makeIndexStmt("", "t", "idx1", []string{"a"}, true, false)
	stmt.AccessMethod = "hash"
	err := c.DefineIndex(stmt)
	if err != nil {
		t.Fatal(err)
	}
}

func TestIndexUniqueOnBtreeOK(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	stmt := makeIndexStmt("", "t", "idx1", []string{"a"}, true, false)
	stmt.AccessMethod = "btree"
	err := c.DefineIndex(stmt)
	if err != nil {
		t.Fatal(err)
	}
}

func TestIndexUniqueOnGistRejected(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	stmt := makeIndexStmt("", "t", "idx1", []string{"a"}, true, false)
	stmt.AccessMethod = "gist"
	err := c.DefineIndex(stmt)
	assertCode(t, err, CodeFeatureNotSupported)
}

func TestIndexUniqueOnGinRejected(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	stmt := makeIndexStmt("", "t", "idx1", []string{"a"}, true, false)
	stmt.AccessMethod = "gin"
	err := c.DefineIndex(stmt)
	assertCode(t, err, CodeFeatureNotSupported)
}

func TestIndexIncludeOnGinRejected(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "b", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	stmt := makeIndexStmt("", "t", "idx1", []string{"a"}, false, false)
	stmt.AccessMethod = "gin"
	stmt.IndexIncludingParams = &nodes.List{Items: []nodes.Node{
		&nodes.IndexElem{Name: "b"},
	}}
	err := c.DefineIndex(stmt)
	assertCode(t, err, CodeFeatureNotSupported)
}

func TestIndexIncludeOnBtreeOK(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "b", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	stmt := makeIndexStmt("", "t", "idx1", []string{"a"}, false, false)
	stmt.AccessMethod = "btree"
	stmt.IndexIncludingParams = &nodes.List{Items: []nodes.Node{
		&nodes.IndexElem{Name: "b"},
	}}
	err := c.DefineIndex(stmt)
	if err != nil {
		t.Fatal(err)
	}
}

func TestIndexIncludeOnGistOK(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "b", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	stmt := makeIndexStmt("", "t", "idx1", []string{"a"}, false, false)
	stmt.AccessMethod = "gist"
	stmt.IndexIncludingParams = &nodes.List{Items: []nodes.Node{
		&nodes.IndexElem{Name: "b"},
	}}
	err := c.DefineIndex(stmt)
	if err != nil {
		t.Fatal(err)
	}
}

func TestIndexHashMultiColumnRejected(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "b", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	stmt := makeIndexStmt("", "t", "idx1", []string{"a", "b"}, false, false)
	stmt.AccessMethod = "hash"
	err := c.DefineIndex(stmt)
	assertCode(t, err, CodeFeatureNotSupported)
}

func TestIndexHashSingleColumnOK(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	stmt := makeIndexStmt("", "t", "idx1", []string{"a"}, false, false)
	stmt.AccessMethod = "hash"
	err := c.DefineIndex(stmt)
	if err != nil {
		t.Fatal(err)
	}
}

// =============================================================================
// Phase 3: View SELECT INTO Validation (view.go)
// =============================================================================

func TestViewSelectIntoRejected(t *testing.T) {
	c := New()
	createTestTable(t, c)

	stmt := makeViewStmt("", "v",
		&SelectStmt{
			From:       []*FromItem{{Kind: FromTable, Table: "t"}},
			TargetList: []*ResTarget{{Val: &Expr{Kind: ExprStar}}},
		}, nil)

	// Simulate SELECT INTO by setting IntoClause on the underlying SelectStmt.
	selStmt := stmt.Query.(*nodes.SelectStmt)
	selStmt.IntoClause = &nodes.IntoClause{
		Rel: &nodes.RangeVar{Relname: "into_table"},
	}
	err := c.DefineView(stmt)
	assertCode(t, err, CodeFeatureNotSupported)
}

func TestViewNormalSelectOK(t *testing.T) {
	c := New()
	createTestTable(t, c)

	err := c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			From:       []*FromItem{{Kind: FromTable, Table: "t"}},
			TargetList: []*ResTarget{{Val: &Expr{Kind: ExprStar}}},
		}, nil))
	if err != nil {
		t.Fatal(err)
	}
}

// =============================================================================
// Phase 4: Comment on Domain Constraint (comment.go)
// =============================================================================

func TestCommentOnDomainConstraint(t *testing.T) {
	c := New()
	// Create a domain with a CHECK constraint.
	err := c.DefineDomain(&nodes.CreateDomainStmt{
		Domainname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "posint"}}},
		Typname:    makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
		Constraints: &nodes.List{Items: []nodes.Node{
			&nodes.Constraint{
				Contype:  nodes.CONSTR_CHECK,
				Conname:  "posint_check",
				RawExpr:  &nodes.A_Const{Val: &nodes.String{Str: "VALUE > 0"}},
			},
		}},
	})
	if err != nil {
		t.Fatal(err)
	}

	// COMMENT ON CONSTRAINT posint_check ON DOMAIN posint
	err = c.CommentObject(&nodes.CommentStmt{
		Objtype: nodes.OBJECT_DOMCONSTRAINT,
		Object: &nodes.List{Items: []nodes.Node{
			&nodes.String{Str: "posint"},
			&nodes.String{Str: "posint_check"},
		}},
		Comment: "this is a domain constraint comment",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCommentOnDomainConstraintNotFound(t *testing.T) {
	c := New()
	c.DefineDomain(&nodes.CreateDomainStmt{
		Domainname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "posint"}}},
		Typname:    makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
	})

	err := c.CommentObject(&nodes.CommentStmt{
		Objtype: nodes.OBJECT_DOMCONSTRAINT,
		Object: &nodes.List{Items: []nodes.Node{
			&nodes.String{Str: "posint"},
			&nodes.String{Str: "no_such_constraint"},
		}},
		Comment: "comment",
	})
	assertCode(t, err, CodeUndefinedObject)
}

func TestCommentOnDomainConstraintNotDomain(t *testing.T) {
	c := New()
	// Create an enum type (not a domain).
	c.DefineEnum(&nodes.CreateEnumStmt{
		TypeName: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "color"}}},
		Vals:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: "red"}}},
	})

	err := c.CommentObject(&nodes.CommentStmt{
		Objtype: nodes.OBJECT_DOMCONSTRAINT,
		Object: &nodes.List{Items: []nodes.Node{
			&nodes.String{Str: "color"},
			&nodes.String{Str: "some_constraint"},
		}},
		Comment: "comment",
	})
	assertCode(t, err, CodeWrongObjectType)
}

// =============================================================================
// Phase 5: DROP AGGREGATE (dropcmds.go)
// =============================================================================

func TestDropAggregateAsFunction(t *testing.T) {
	c := New()
	// Create a function to act as our "aggregate".
	c.CreateFunctionStmt(makeCreateFuncStmt("", "myagg",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		TypeName{Name: "int4", TypeMod: -1},
		"sql", "SELECT $1", false))

	// DROP AGGREGATE dispatches to removeFunctionObjects.
	err := c.RemoveObjects(&nodes.DropStmt{
		RemoveType: int(nodes.OBJECT_AGGREGATE),
		Objects: &nodes.List{Items: []nodes.Node{
			&nodes.ObjectWithArgs{
				Objname:        &nodes.List{Items: []nodes.Node{&nodes.String{Str: "myagg"}}},
				ArgsUnspecified: true,
			},
		}},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Verify it was removed.
	procs := c.LookupProcByName("myagg")
	found := false
	for _, p := range procs {
		if c.userProcs[p.OID] != nil {
			found = true
		}
	}
	if found {
		t.Error("aggregate/function should be removed")
	}
}

// =============================================================================
// Phase 6: Sequence Duplicate DefElem (sequence.go)
// =============================================================================

func TestSeqDuplicateIncrementRejected(t *testing.T) {
	c := New()
	err := c.DefineSequence(makeCreateSeqStmt("", "s", false,
		seqOptInt("increment", 1),
		seqOptInt("increment", 2),
	))
	assertCode(t, err, CodeSyntaxError)
}

func TestSeqDuplicateMaxvalueRejected(t *testing.T) {
	c := New()
	err := c.DefineSequence(makeCreateSeqStmt("", "s", false,
		seqOptInt("maxvalue", 100),
		seqOptInt("maxvalue", 200),
	))
	assertCode(t, err, CodeSyntaxError)
}

func TestSeqDuplicateMinvalueRejected(t *testing.T) {
	c := New()
	err := c.DefineSequence(makeCreateSeqStmt("", "s", false,
		seqOptInt("minvalue", 1),
		seqOptInt("minvalue", 2),
	))
	assertCode(t, err, CodeSyntaxError)
}

func TestSeqDuplicateCacheRejected(t *testing.T) {
	c := New()
	err := c.DefineSequence(makeCreateSeqStmt("", "s", false,
		seqOptInt("cache", 1),
		seqOptInt("cache", 2),
	))
	assertCode(t, err, CodeSyntaxError)
}

func TestSeqDuplicateStartRejected(t *testing.T) {
	c := New()
	err := c.DefineSequence(makeCreateSeqStmt("", "s", false,
		seqOptInt("start", 1),
		seqOptInt("start", 2),
	))
	assertCode(t, err, CodeSyntaxError)
}

func TestSeqDuplicateAsRejected(t *testing.T) {
	c := New()
	err := c.DefineSequence(makeCreateSeqStmt("", "s", false,
		seqOptStr("as", "smallint"),
		seqOptStr("as", "integer"),
	))
	assertCode(t, err, CodeSyntaxError)
}

// =============================================================================
// Phase 7: ALTER SEQUENCE Validations (sequence.go)
// =============================================================================

func TestAlterSeqDuplicateOptionRejected(t *testing.T) {
	c := New()
	c.DefineSequence(makeCreateSeqStmt("", "s", false))

	err := c.AlterSequenceStmt(makeAlterSeqStmt("", "s",
		seqOptInt("increment", 1),
		seqOptInt("increment", 2),
	))
	assertCode(t, err, CodeSyntaxError)
}

func TestAlterSeqRestartDuplicate(t *testing.T) {
	c := New()
	c.DefineSequence(makeCreateSeqStmt("", "s", false))

	err := c.AlterSequenceStmt(makeAlterSeqStmt("", "s",
		seqOptInt("restart", 1),
		seqOptInt("restart", 2),
	))
	assertCode(t, err, CodeSyntaxError)
}

func TestAlterSeqStartOutOfRange(t *testing.T) {
	c := New()
	c.DefineSequence(makeCreateSeqStmt("", "s", false,
		seqOptInt("minvalue", 1),
		seqOptInt("maxvalue", 100),
	))

	// Set start beyond maxvalue.
	err := c.AlterSequenceStmt(makeAlterSeqStmt("", "s",
		seqOptInt("start", 200),
	))
	assertCode(t, err, CodeInvalidParameterValue)
}

func TestAlterSeqRestartOutOfRange(t *testing.T) {
	c := New()
	c.DefineSequence(makeCreateSeqStmt("", "s", false,
		seqOptInt("minvalue", 1),
		seqOptInt("maxvalue", 100),
	))

	// Restart beyond maxvalue.
	err := c.AlterSequenceStmt(makeAlterSeqStmt("", "s",
		seqOptInt("restart", 200),
	))
	assertCode(t, err, CodeInvalidParameterValue)
}

// =============================================================================
// Phase 8: Sequence OWNED BY Validations (sequence.go)
// =============================================================================

func TestSeqOwnedByViewAllowed(t *testing.T) {
	c := New()
	createTestTable(t, c)
	c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			From:       []*FromItem{{Kind: FromTable, Table: "t"}},
			TargetList: []*ResTarget{{Val: &Expr{Kind: ExprStar}}},
		}, nil))

	err := c.DefineSequence(makeCreateSeqStmt("", "s", false,
		seqOptOwnedBy("v.x"),
	))
	// Views can own sequences (PG allows this).
	if err != nil {
		t.Fatal(err)
	}
}

func TestSeqOwnedByTableOK(t *testing.T) {
	c := New()
	createTestTable(t, c)
	err := c.DefineSequence(makeCreateSeqStmt("", "s", false,
		seqOptOwnedBy("t.x"),
	))
	if err != nil {
		t.Fatal(err)
	}
}

// =============================================================================
// Phase 9: ProcessUtility no-op DDL dispatches (utility.go)
// =============================================================================

func TestProcessUtilityNoOpDDL(t *testing.T) {
	c := New()

	noOpStmts := []nodes.Node{
		&nodes.CreateConversionStmt{},
		&nodes.CreatePublicationStmt{},
		&nodes.AlterPublicationStmt{},
		&nodes.CreateSubscriptionStmt{},
		&nodes.AlterSubscriptionStmt{},
		&nodes.DropSubscriptionStmt{},
		&nodes.CreateTransformStmt{},
		&nodes.CreateEventTrigStmt{},
		&nodes.AlterEventTrigStmt{},
	}

	for _, stmt := range noOpStmts {
		err := c.ProcessUtility(stmt)
		if err != nil {
			t.Errorf("%T: unexpected error: %v", stmt, err)
		}
	}
}

// =============================================================================
// Phase 10: Index on wrong relkind (indexcmds.go — existing)
// =============================================================================

func TestIndexOnViewRejected(t *testing.T) {
	c := New()
	createTestTable(t, c)
	c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			From:       []*FromItem{{Kind: FromTable, Table: "t"}},
			TargetList: []*ResTarget{{Val: &Expr{Kind: ExprStar}}},
		}, nil))

	err := c.DefineIndex(makeIndexStmt("", "v", "vidx", []string{"id"}, false, false))
	assertCode(t, err, CodeWrongObjectType)
}

func TestIndexOnMatviewOK(t *testing.T) {
	c := New()
	createTestTable(t, c)

	// Create a matview.
	c.ExecCreateTableAs(&nodes.CreateTableAsStmt{
		Query: &nodes.SelectStmt{
			TargetList: &nodes.List{Items: []nodes.Node{
				&nodes.ResTarget{Val: &nodes.ColumnRef{Fields: &nodes.List{Items: []nodes.Node{&nodes.A_Star{}}}}},
			}},
			FromClause: &nodes.List{Items: []nodes.Node{
				&nodes.RangeVar{Relname: "t"},
			}},
		},
		Into: &nodes.IntoClause{
			Rel: &nodes.RangeVar{Relname: "mv"},
		},
		Objtype: nodes.OBJECT_MATVIEW,
	})

	err := c.DefineIndex(makeIndexStmt("", "mv", "mvidx", []string{"x"}, false, false))
	if err != nil {
		t.Fatal(err)
	}
}

// seqOptOwnedBy and seqOptStr helpers defined in sequence_test.go / round5_test.go
