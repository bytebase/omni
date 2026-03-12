package ddl

import (
	"math"
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// =============================================================================
// Phase 1: Schema reserved name validation
// =============================================================================

func TestCreateSchemaReservedPrefix(t *testing.T) {
	c := New()
	err := c.CreateSchemaCommand(&nodes.CreateSchemaStmt{Schemaname: "pg_test"})
	assertErrorCode(t, err, CodeReservedName)
}

func TestCreateSchemaPgCatalog(t *testing.T) {
	c := New()
	err := c.CreateSchemaCommand(&nodes.CreateSchemaStmt{Schemaname: "pg_catalog"})
	assertErrorCode(t, err, CodeReservedName)
}

func TestCreateSchemaNormalName(t *testing.T) {
	c := New()
	err := c.CreateSchemaCommand(&nodes.CreateSchemaStmt{Schemaname: "my_schema"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateSchemaPgPrefixLong(t *testing.T) {
	c := New()
	err := c.CreateSchemaCommand(&nodes.CreateSchemaStmt{Schemaname: "pg_temp_123"})
	assertErrorCode(t, err, CodeReservedName)
}

// =============================================================================
// Phase 2: View validations
// =============================================================================

func TestDefineViewUnlogged(t *testing.T) {
	c := New()
	createTestTable(t, c)

	stmt := makeViewStmt("", "v",
		&SelectStmt{
			From:       []*FromItem{{Kind: FromTable, Table: "t"}},
			TargetList: []*ResTarget{{Val: &Expr{Kind: ExprStar}}},
		}, nil)
	stmt.View.Relpersistence = 'u'
	err := c.DefineView(stmt)
	assertErrorCode(t, err, CodeFeatureNotSupported)
}

func TestDefineViewNormalPersistence(t *testing.T) {
	c := New()
	createTestTable(t, c)

	stmt := makeViewStmt("", "v",
		&SelectStmt{
			From:       []*FromItem{{Kind: FromTable, Table: "t"}},
			TargetList: []*ResTarget{{Val: &Expr{Kind: ExprStar}}},
		}, nil)
	err := c.DefineView(stmt)
	if err != nil {
		t.Fatal(err)
	}
}

// =============================================================================
// Phase 3: Sequence bounds validation
// =============================================================================

func TestSequenceSmallintMaxOutOfRange(t *testing.T) {
	c := New()
	// AS smallint with MAXVALUE beyond int16 range.
	err := c.DefineSequence(makeCreateSeqStmt("", "s", false,
		seqOptStr("as", "smallint"),
		seqOptInt("maxvalue", math.MaxInt32),
	))
	assertErrorCode(t, err, CodeInvalidParameterValue)
}

func TestSequenceSmallintMinOutOfRange(t *testing.T) {
	c := New()
	err := c.DefineSequence(makeCreateSeqStmt("", "s", false,
		seqOptStr("as", "smallint"),
		seqOptInt("minvalue", math.MinInt32),
	))
	assertErrorCode(t, err, CodeInvalidParameterValue)
}

func TestSequenceIntegerMaxOutOfRange(t *testing.T) {
	c := New()
	err := c.DefineSequence(makeCreateSeqStmt("", "s", false,
		seqOptStr("as", "integer"),
		seqOptInt("maxvalue", math.MaxInt64),
	))
	assertErrorCode(t, err, CodeInvalidParameterValue)
}

func TestSequenceBigintMaxOK(t *testing.T) {
	c := New()
	// Bigint allows full range — no error.
	err := c.DefineSequence(makeCreateSeqStmt("", "s", false,
		seqOptStr("as", "bigint"),
		seqOptInt("maxvalue", math.MaxInt64),
	))
	if err != nil {
		t.Fatal(err)
	}
}

func TestSequenceSmallintBoundsOK(t *testing.T) {
	c := New()
	err := c.DefineSequence(makeCreateSeqStmt("", "s", false,
		seqOptStr("as", "smallint"),
		seqOptInt("minvalue", 1),
		seqOptInt("maxvalue", 100),
	))
	if err != nil {
		t.Fatal(err)
	}
}

func TestAlterSequenceBoundsOutOfRange(t *testing.T) {
	c := New()
	// Create smallint sequence first.
	c.DefineSequence(makeCreateSeqStmt("", "s", false,
		seqOptStr("as", "smallint"),
	))
	// Try to set maxvalue beyond range.
	err := c.AlterSequenceStmt(makeAlterSeqStmt("", "s",
		seqOptInt("maxvalue", math.MaxInt32),
	))
	assertErrorCode(t, err, CodeInvalidParameterValue)
}

func TestAlterSequenceMinGteMax(t *testing.T) {
	c := New()
	c.DefineSequence(makeCreateSeqStmt("", "s", false))
	err := c.AlterSequenceStmt(makeAlterSeqStmt("", "s",
		seqOptInt("minvalue", 100),
		seqOptInt("maxvalue", 50),
	))
	assertErrorCode(t, err, CodeInvalidParameterValue)
}

// =============================================================================
// Phase 4: Domain base type validation
// =============================================================================

func TestDomainBaseTypeOK(t *testing.T) {
	c := New()
	// Domain over int4 (base type) should work.
	err := c.DefineDomain(&nodes.CreateDomainStmt{
		Domainname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "posint"}}},
		Typname:    makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDomainOverEnumOK(t *testing.T) {
	c := New()
	// Create enum first.
	c.DefineEnum(&nodes.CreateEnumStmt{
		TypeName: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "color"}}},
		Vals:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: "red"}}},
	})
	// Domain over enum should work.
	err := c.DefineDomain(&nodes.CreateDomainStmt{
		Domainname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "mycolor"}}},
		Typname:    makeTypeNameNode(TypeName{Name: "color", TypeMod: -1}),
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDomainOverDomainOK(t *testing.T) {
	c := New()
	// Create base domain.
	c.DefineDomain(&nodes.CreateDomainStmt{
		Domainname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "posint"}}},
		Typname:    makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
	})
	// Domain over domain should work.
	err := c.DefineDomain(&nodes.CreateDomainStmt{
		Domainname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "posint2"}}},
		Typname:    makeTypeNameNode(TypeName{Name: "posint", TypeMod: -1}),
	})
	if err != nil {
		t.Fatal(err)
	}
}

// =============================================================================
// Phase 5: Enum neighbor label validation
// =============================================================================

func TestEnumAddValueBeforeNonexistent(t *testing.T) {
	c := New()
	c.DefineEnum(&nodes.CreateEnumStmt{
		TypeName: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "color"}}},
		Vals:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: "red"}, &nodes.String{Str: "green"}}},
	})
	err := c.AlterEnumStmt(&nodes.AlterEnumStmt{
		Typname:         &nodes.List{Items: []nodes.Node{&nodes.String{Str: "color"}}},
		Newval:          "blue",
		NewvalNeighbor:  "yellow",
		NewvalIsAfter:   false,
	})
	assertErrorCode(t, err, CodeInvalidParameterValue)
}

func TestEnumAddValueAfterNonexistent(t *testing.T) {
	c := New()
	c.DefineEnum(&nodes.CreateEnumStmt{
		TypeName: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "color"}}},
		Vals:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: "red"}, &nodes.String{Str: "green"}}},
	})
	err := c.AlterEnumStmt(&nodes.AlterEnumStmt{
		Typname:         &nodes.List{Items: []nodes.Node{&nodes.String{Str: "color"}}},
		Newval:          "blue",
		NewvalNeighbor:  "yellow",
		NewvalIsAfter:   true,
	})
	assertErrorCode(t, err, CodeInvalidParameterValue)
}

func TestEnumAddValueBeforeExistingOK(t *testing.T) {
	c := New()
	c.DefineEnum(&nodes.CreateEnumStmt{
		TypeName: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "color"}}},
		Vals:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: "red"}, &nodes.String{Str: "green"}}},
	})
	err := c.AlterEnumStmt(&nodes.AlterEnumStmt{
		Typname:         &nodes.List{Items: []nodes.Node{&nodes.String{Str: "color"}}},
		Newval:          "blue",
		NewvalNeighbor:  "green",
		NewvalIsAfter:   false,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestEnumAddValueAfterExistingOK(t *testing.T) {
	c := New()
	c.DefineEnum(&nodes.CreateEnumStmt{
		TypeName: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "color"}}},
		Vals:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: "red"}, &nodes.String{Str: "green"}}},
	})
	err := c.AlterEnumStmt(&nodes.AlterEnumStmt{
		Typname:         &nodes.List{Items: []nodes.Node{&nodes.String{Str: "color"}}},
		Newval:          "blue",
		NewvalNeighbor:  "red",
		NewvalIsAfter:   true,
	})
	if err != nil {
		t.Fatal(err)
	}
}

// =============================================================================
// Phase 6: Trigger validation — TRUNCATE on views
// =============================================================================

func TestTruncateTriggerOnView(t *testing.T) {
	c := New()
	createTestTable(t, c)
	createTriggerFunc(t, c)

	c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			From:       []*FromItem{{Kind: FromTable, Table: "t"}},
			TargetList: []*ResTarget{{Val: &Expr{Kind: ExprStar}}},
		}, nil))

	// TRUNCATE trigger on a view should be rejected.
	err := c.CreateTriggerStmt(makeCreateTrigStmt("", "v", "trig", TriggerAfter, TriggerEventTruncate, false, "", "trig_fn"))
	assertCode(t, err, CodeInvalidObjectDefinition)
}

func TestTruncateTriggerOnTableOK(t *testing.T) {
	c := New()
	createTestTable(t, c)
	createTriggerFunc(t, c)

	// TRUNCATE trigger on table (STATEMENT level) should work.
	err := c.CreateTriggerStmt(makeCreateTrigStmt("", "t", "trig", TriggerAfter, TriggerEventTruncate, false, "", "trig_fn"))
	if err != nil {
		t.Fatal(err)
	}
}

// =============================================================================
// Phase 7: Function parameter mode filtering
// =============================================================================

func TestFunctionOutParamNotInSignature(t *testing.T) {
	c := New()

	// Create function with IN + OUT params: f(IN int4, OUT text).
	// The signature should only include int4.
	stmt := &nodes.CreateFunctionStmt{
		Funcname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "f"}}},
		Parameters: &nodes.List{Items: []nodes.Node{
			&nodes.FunctionParameter{
				Name:    "a",
				ArgType: makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
				Mode:    nodes.FUNC_PARAM_IN,
			},
			&nodes.FunctionParameter{
				Name:    "b",
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

	// Now create another function f(int4, text) with IN, IN — different signature.
	stmt2 := &nodes.CreateFunctionStmt{
		Funcname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "f"}}},
		Parameters: &nodes.List{Items: []nodes.Node{
			&nodes.FunctionParameter{
				Name:    "a",
				ArgType: makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
				Mode:    nodes.FUNC_PARAM_IN,
			},
			&nodes.FunctionParameter{
				Name:    "b",
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
	err = c.CreateFunctionStmt(stmt2)
	if err != nil {
		t.Fatal("should not conflict: OUT param excluded from signature")
	}

	// But creating f(int4) IN should conflict with f(IN int4, OUT text).
	stmt3 := &nodes.CreateFunctionStmt{
		Funcname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "f"}}},
		Parameters: &nodes.List{Items: []nodes.Node{
			&nodes.FunctionParameter{
				Name:    "a",
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
	err = c.CreateFunctionStmt(stmt3)
	assertErrorCode(t, err, CodeDuplicateFunction)
}

func TestFunctionVariadicIncludedInSignature(t *testing.T) {
	c := New()

	// VARIADIC param IS included in signature.
	stmt := &nodes.CreateFunctionStmt{
		Funcname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "g"}}},
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

// =============================================================================
// Helper
// =============================================================================

func seqOptStr(name, val string) *nodes.DefElem {
	return &nodes.DefElem{Defname: name, Arg: &nodes.String{Str: val}}
}
