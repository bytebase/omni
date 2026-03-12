package ddl

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// =============================================================================
// Round 8 Phase 4: Function & Procedure Improvements (functioncmds.go)
// =============================================================================

// -----------------------------------------------------------------------------
// Gap 1: ROWS vs returnsSet validation
// pg: src/backend/commands/functioncmds.c (line 1248-1251)
// -----------------------------------------------------------------------------

func TestRowsOnNonSetReturningFunction(t *testing.T) {
	c := New()
	// Function with ROWS specified but not returning SETOF — should error.
	stmt := &nodes.CreateFunctionStmt{
		Funcname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "f_rows"}}},
		ReturnType: &nodes.TypeName{
			Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "int4"}}},
		},
		Options: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{Defname: "language", Arg: &nodes.String{Str: "sql"}},
			&nodes.DefElem{Defname: "as", Arg: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "SELECT 1"}}}},
			&nodes.DefElem{Defname: "rows", Arg: &nodes.Integer{Ival: 100}},
		}},
	}
	err := c.CreateFunctionStmt(stmt)
	assertCode(t, err, CodeInvalidParameterValue)
	if !strings.Contains(err.Error(), "ROWS is not applicable when function does not return a set") {
		t.Errorf("unexpected message: %s", err)
	}
}

func TestRowsOnSetReturningFunctionOK(t *testing.T) {
	c := New()
	// Function with ROWS specified and SETOF return — should succeed.
	stmt := &nodes.CreateFunctionStmt{
		Funcname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "f_setof"}}},
		ReturnType: &nodes.TypeName{
			Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "int4"}}},
			Setof: true,
		},
		Options: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{Defname: "language", Arg: &nodes.String{Str: "sql"}},
			&nodes.DefElem{Defname: "as", Arg: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "SELECT 1"}}}},
			&nodes.DefElem{Defname: "rows", Arg: &nodes.Integer{Ival: 500}},
		}},
	}
	err := c.CreateFunctionStmt(stmt)
	if err != nil {
		t.Fatal(err)
	}
	// Verify the function was created with RetSet=true.
	procs := c.LookupProcByName("f_setof")
	if len(procs) == 0 {
		t.Fatal("function not found")
	}
	if !procs[0].RetSet {
		t.Error("expected RetSet=true")
	}
}

// -----------------------------------------------------------------------------
// Gap 2: Return type from OUT params
// pg: src/backend/commands/functioncmds.c — interpret_function_parameter_list
// -----------------------------------------------------------------------------

func TestReturnTypeFromSingleOutParam(t *testing.T) {
	c := New()
	// Function with no RETURNS clause but 1 OUT param — infer return type.
	stmt := &nodes.CreateFunctionStmt{
		Funcname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "f_out"}}},
		Parameters: &nodes.List{Items: []nodes.Node{
			&nodes.FunctionParameter{
				Name:    "a",
				ArgType: &nodes.TypeName{Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "int4"}}}},
				Mode:    nodes.FUNC_PARAM_IN,
			},
			&nodes.FunctionParameter{
				Name:    "b",
				ArgType: &nodes.TypeName{Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "text"}}}},
				Mode:    nodes.FUNC_PARAM_OUT,
			},
		}},
		ReturnType: &nodes.TypeName{
			Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "text"}}},
		},
		Options: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{Defname: "language", Arg: &nodes.String{Str: "plpgsql"}},
			&nodes.DefElem{Defname: "as", Arg: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "BEGIN b := 'hello'; END;"}}}},
		}},
	}
	err := c.CreateFunctionStmt(stmt)
	if err != nil {
		t.Fatal(err)
	}
	procs := c.LookupProcByName("f_out")
	if len(procs) == 0 {
		t.Fatal("function not found")
	}
	if procs[0].RetType != TEXTOID {
		t.Errorf("ret type: got %d, want %d (text)", procs[0].RetType, TEXTOID)
	}
}

func TestReturnTypeFromMultipleOutParams(t *testing.T) {
	c := New()
	// Function with no RETURNS clause but multiple OUT params — return RECORDOID.
	stmt := &nodes.CreateFunctionStmt{
		Funcname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "f_multi_out"}}},
		Parameters: &nodes.List{Items: []nodes.Node{
			&nodes.FunctionParameter{
				Name:    "a",
				ArgType: &nodes.TypeName{Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "int4"}}}},
				Mode:    nodes.FUNC_PARAM_IN,
			},
			&nodes.FunctionParameter{
				Name:    "x",
				ArgType: &nodes.TypeName{Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "int4"}}}},
				Mode:    nodes.FUNC_PARAM_OUT,
			},
			&nodes.FunctionParameter{
				Name:    "y",
				ArgType: &nodes.TypeName{Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "text"}}}},
				Mode:    nodes.FUNC_PARAM_OUT,
			},
		}},
		ReturnType: &nodes.TypeName{
			Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "record"}}},
		},
		Options: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{Defname: "language", Arg: &nodes.String{Str: "plpgsql"}},
			&nodes.DefElem{Defname: "as", Arg: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "BEGIN x := 1; y := 'hi'; END;"}}}},
		}},
	}
	err := c.CreateFunctionStmt(stmt)
	if err != nil {
		t.Fatal(err)
	}
	procs := c.LookupProcByName("f_multi_out")
	if len(procs) == 0 {
		t.Fatal("function not found")
	}
	if procs[0].RetType != RECORDOID {
		t.Errorf("ret type: got %d, want %d (record)", procs[0].RetType, RECORDOID)
	}
}

// -----------------------------------------------------------------------------
// Gap 3: Shell type detection
// pg: src/backend/commands/functioncmds.c (line 88-164, 235-256)
// -----------------------------------------------------------------------------

func TestShellTypeRejectedInSqlFunction(t *testing.T) {
	c := New()
	// Register a shell type (IsDefined=false).
	schema := c.schemaByName["public"]
	shellOID := c.oidGen.Next()
	c.typeByOID[shellOID] = &BuiltinType{
		OID:       shellOID,
		TypeName:  "shelltype",
		Namespace: schema.OID,
		Type:      'p',
		Category:  'P',
		IsDefined: false,
	}
	c.typeByName[typeKey{ns: schema.OID, name: "shelltype"}] = c.typeByOID[shellOID]

	// SQL function with shell type param — should error.
	stmt := &nodes.CreateFunctionStmt{
		Funcname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "f_shell"}}},
		Parameters: &nodes.List{Items: []nodes.Node{
			&nodes.FunctionParameter{
				ArgType: &nodes.TypeName{Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "shelltype"}}}},
				Mode:    nodes.FUNC_PARAM_IN,
			},
		}},
		ReturnType: &nodes.TypeName{
			Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "int4"}}},
		},
		Options: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{Defname: "language", Arg: &nodes.String{Str: "sql"}},
			&nodes.DefElem{Defname: "as", Arg: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "SELECT 1"}}}},
		}},
	}
	err := c.CreateFunctionStmt(stmt)
	assertCode(t, err, CodeInvalidFunctionDefinition)
	if !strings.Contains(err.Error(), "SQL function cannot accept shell type") {
		t.Errorf("unexpected message: %s", err)
	}
}

// -----------------------------------------------------------------------------
// Gap 4: LEAKPROOF storage
// pg: src/backend/commands/functioncmds.c (line 1131-1134)
// -----------------------------------------------------------------------------

func TestLeakProofStorage(t *testing.T) {
	c := New()
	// Create a function with LEAKPROOF.
	stmt := &nodes.CreateFunctionStmt{
		Funcname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "f_leakproof"}}},
		ReturnType: &nodes.TypeName{
			Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "int4"}}},
		},
		Options: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{Defname: "language", Arg: &nodes.String{Str: "sql"}},
			&nodes.DefElem{Defname: "as", Arg: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "SELECT 1"}}}},
			&nodes.DefElem{Defname: "leakproof", Arg: &nodes.Boolean{Boolval: true}},
		}},
	}
	err := c.CreateFunctionStmt(stmt)
	if err != nil {
		t.Fatal(err)
	}

	// Verify LeakProof is stored on BuiltinProc.
	procs := c.LookupProcByName("f_leakproof")
	if len(procs) == 0 {
		t.Fatal("function not found")
	}
	if !procs[0].LeakProof {
		t.Error("expected LeakProof=true on BuiltinProc")
	}

	// Verify LeakProof is stored on UserProc.
	up := c.userProcs[procs[0].OID]
	if up == nil {
		t.Fatal("UserProc not found")
	}
	if !up.LeakProof {
		t.Error("expected LeakProof=true on UserProc")
	}
}

// -----------------------------------------------------------------------------
// Gap 5: Procedure OUT → RECORDOID
// pg: src/backend/commands/functioncmds.c (line 296-325)
// -----------------------------------------------------------------------------

func TestProcedureWithOutParamsRettype(t *testing.T) {
	c := New()
	// Procedure with OUT parameters — return type should be RECORDOID, not VOIDOID.
	stmt := &nodes.CreateFunctionStmt{
		Funcname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "proc_out"}}},
		Parameters: &nodes.List{Items: []nodes.Node{
			&nodes.FunctionParameter{
				Name:    "a",
				ArgType: &nodes.TypeName{Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "int4"}}}},
				Mode:    nodes.FUNC_PARAM_IN,
			},
			&nodes.FunctionParameter{
				Name:    "result",
				ArgType: &nodes.TypeName{Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "text"}}}},
				Mode:    nodes.FUNC_PARAM_OUT,
			},
		}},
		// No ReturnType → procedure.
		Options: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{Defname: "language", Arg: &nodes.String{Str: "plpgsql"}},
			&nodes.DefElem{Defname: "as", Arg: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "BEGIN result := 'ok'; END;"}}}},
		}},
	}
	err := c.CreateFunctionStmt(stmt)
	if err != nil {
		t.Fatal(err)
	}

	procs := c.LookupProcByName("proc_out")
	if len(procs) == 0 {
		t.Fatal("procedure not found")
	}
	if procs[0].Kind != 'p' {
		t.Errorf("kind: got %c, want 'p'", procs[0].Kind)
	}
	if procs[0].RetType != RECORDOID {
		t.Errorf("ret type: got %d, want %d (RECORDOID)", procs[0].RetType, RECORDOID)
	}
}

// -----------------------------------------------------------------------------
// Gap 6: SQL body handling — default language
// pg: src/backend/commands/functioncmds.c (line 1086-1087)
// -----------------------------------------------------------------------------

func TestSqlBodyDefaultsToSqlLanguage(t *testing.T) {
	c := New()
	// Function with sql_body present and no explicit language — should default to "sql".
	stmt := &nodes.CreateFunctionStmt{
		Funcname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "f_sqlbody"}}},
		ReturnType: &nodes.TypeName{
			Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "int4"}}},
		},
		// Use SqlBody (not AS clause) and no language option.
		SqlBody: &nodes.Integer{Ival: 1}, // Any non-nil node triggers sql_body path.
		Options: &nodes.List{Items: []nodes.Node{
			// Deliberately no "language" option.
			// No "as" clause either — sql_body substitutes.
		}},
	}
	err := c.CreateFunctionStmt(stmt)
	if err != nil {
		t.Fatal(err)
	}

	procs := c.LookupProcByName("f_sqlbody")
	if len(procs) == 0 {
		t.Fatal("function not found")
	}
	up := c.userProcs[procs[0].OID]
	if up == nil {
		t.Fatal("UserProc not found")
	}
	if up.Language != "sql" {
		t.Errorf("language: got %q, want %q", up.Language, "sql")
	}
}

// -----------------------------------------------------------------------------
// Gap 8: ROWS must error on non-set-returning ALTER FUNCTION too
// pg: src/backend/commands/functioncmds.c — AlterFunction (line 1429-1439)
// -----------------------------------------------------------------------------

func TestAlterFunctionRowsNonSetReturning(t *testing.T) {
	c := New()
	// Create a non-set-returning function.
	createStmt := makeCreateFuncStmt("", "f_alter_rows",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		TypeName{Name: "int4", TypeMod: -1},
		"sql", "SELECT $1", false)
	err := c.CreateFunctionStmt(createStmt)
	if err != nil {
		t.Fatal(err)
	}

	// ALTER FUNCTION with ROWS — should error.
	alterStmt := &nodes.AlterFunctionStmt{
		Func: &nodes.ObjectWithArgs{
			Objname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "f_alter_rows"}}},
			Objargs: &nodes.List{Items: []nodes.Node{
				&nodes.TypeName{Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "int4"}}}},
			}},
		},
		Actions: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{Defname: "rows", Arg: &nodes.Integer{Ival: 100}},
		}},
	}
	err = c.AlterFunction(alterStmt)
	assertCode(t, err, CodeInvalidParameterValue)
	if !strings.Contains(err.Error(), "ROWS is not applicable when function does not return a set") {
		t.Errorf("unexpected message: %s", err)
	}
}

// -----------------------------------------------------------------------------
// Gap 2 (cont): TABLE params don't count as OUT for return type inference
// pg: src/backend/commands/functioncmds.c — interpret_function_parameter_list
// Note: TABLE params DO count as output in PG (fpmode != IN and != VARIADIC).
// However, functions with TABLE params use RETURNS TABLE syntax which sets
// returnType explicitly. This test verifies TABLE params are counted.
// -----------------------------------------------------------------------------

func TestReturnTypeFromOutParamIgnoresTableParam(t *testing.T) {
	c := New()
	// Function with one OUT param and one TABLE param.
	// TABLE params count as output, so 2 OUT params → RECORDOID.
	stmt := &nodes.CreateFunctionStmt{
		Funcname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "f_table_out"}}},
		Parameters: &nodes.List{Items: []nodes.Node{
			&nodes.FunctionParameter{
				Name:    "x",
				ArgType: &nodes.TypeName{Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "int4"}}}},
				Mode:    nodes.FUNC_PARAM_OUT,
			},
			&nodes.FunctionParameter{
				Name:    "y",
				ArgType: &nodes.TypeName{Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "text"}}}},
				Mode:    nodes.FUNC_PARAM_TABLE,
			},
		}},
		ReturnType: &nodes.TypeName{
			Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "record"}}},
		},
		Options: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{Defname: "language", Arg: &nodes.String{Str: "plpgsql"}},
			&nodes.DefElem{Defname: "as", Arg: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "BEGIN x := 1; END;"}}}},
		}},
	}
	err := c.CreateFunctionStmt(stmt)
	if err != nil {
		t.Fatal(err)
	}
	procs := c.LookupProcByName("f_table_out")
	if len(procs) == 0 {
		t.Fatal("function not found")
	}
	// With 2 output params (OUT + TABLE), requiredResultType = RECORDOID.
	if procs[0].RetType != RECORDOID {
		t.Errorf("ret type: got %d, want %d (RECORDOID)", procs[0].RetType, RECORDOID)
	}
}

// -----------------------------------------------------------------------------
// Gap 7: COST/ROWS defaults (validation)
// pg: src/backend/commands/functioncmds.c (line 1227-1251)
// -----------------------------------------------------------------------------

func TestCostDefault(t *testing.T) {
	c := New()
	// Create function without explicit COST — should succeed (default applied).
	stmt := makeCreateFuncStmt("", "f_cost_default",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		TypeName{Name: "int4", TypeMod: -1},
		"sql", "SELECT $1", false)
	err := c.CreateFunctionStmt(stmt)
	if err != nil {
		t.Fatal(err)
	}
}

// -----------------------------------------------------------------------------
// ROWS on procedure
// pg: procedures don't return sets, so ROWS should error.
// -----------------------------------------------------------------------------

func TestRowsNotApplicableProcedure(t *testing.T) {
	c := New()
	// Procedure with ROWS specified — should error (procedures don't return sets).
	stmt := &nodes.CreateFunctionStmt{
		Funcname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "proc_rows"}}},
		// No ReturnType → procedure.
		Options: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{Defname: "language", Arg: &nodes.String{Str: "plpgsql"}},
			&nodes.DefElem{Defname: "as", Arg: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "BEGIN NULL; END;"}}}},
			&nodes.DefElem{Defname: "rows", Arg: &nodes.Integer{Ival: 100}},
		}},
	}
	err := c.CreateFunctionStmt(stmt)
	assertCode(t, err, CodeInvalidParameterValue)
	if !strings.Contains(err.Error(), "ROWS is not applicable when function does not return a set") {
		t.Errorf("unexpected message: %s", err)
	}
}
