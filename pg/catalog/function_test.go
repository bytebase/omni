package catalog

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// --- Function AST helpers ---

func makeCreateFuncStmt(schema, name string, argTypes []TypeName, retType TypeName, language, body string, orReplace bool) *nodes.CreateFunctionStmt {
	var nameItems []nodes.Node
	if schema != "" {
		nameItems = append(nameItems, &nodes.String{Str: schema})
	}
	nameItems = append(nameItems, &nodes.String{Str: name})

	var params []nodes.Node
	for _, at := range argTypes {
		fp := &nodes.FunctionParameter{
			ArgType: makeTypeNameNode(at),
			Mode:    nodes.FUNC_PARAM_IN,
		}
		params = append(params, fp)
	}

	stmt := &nodes.CreateFunctionStmt{
		Funcname:    &nodes.List{Items: nameItems},
		IsOrReplace: orReplace,
	}
	if retType.Name != "" {
		stmt.ReturnType = makeTypeNameNode(retType)
	}
	if len(params) > 0 {
		stmt.Parameters = &nodes.List{Items: params}
	}

	var opts []nodes.Node
	if language != "" {
		opts = append(opts, &nodes.DefElem{Defname: "language", Arg: &nodes.String{Str: language}})
	}
	if body != "" {
		opts = append(opts, &nodes.DefElem{Defname: "as", Arg: &nodes.List{Items: []nodes.Node{&nodes.String{Str: body}}}})
	}
	if retType.Name == "" {
		opts = append(opts, &nodes.DefElem{Defname: "isProcedure", Arg: &nodes.Integer{Ival: 1}})
	}
	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}

	return stmt
}

func markCreateFunctionStmtAsProcedure(stmt *nodes.CreateFunctionStmt) *nodes.CreateFunctionStmt {
	procDef := &nodes.DefElem{Defname: "isProcedure", Arg: &nodes.Integer{Ival: 1}}
	if stmt.Options == nil {
		stmt.Options = &nodes.List{Items: []nodes.Node{procDef}}
	} else {
		stmt.Options.Items = append(stmt.Options.Items, procDef)
	}
	return stmt
}

func makeDropFuncStmt(schema, name string, argTypes []TypeName, ifExists, cascade bool) *nodes.DropStmt {
	behavior := int(nodes.DROP_RESTRICT)
	if cascade {
		behavior = int(nodes.DROP_CASCADE)
	}

	var nameItems []nodes.Node
	if schema != "" {
		nameItems = append(nameItems, &nodes.String{Str: schema})
	}
	nameItems = append(nameItems, &nodes.String{Str: name})

	owa := &nodes.ObjectWithArgs{
		Objname: &nodes.List{Items: nameItems},
	}
	if len(argTypes) > 0 {
		var args []nodes.Node
		for _, at := range argTypes {
			args = append(args, makeTypeNameNode(at))
		}
		owa.Objargs = &nodes.List{Items: args}
	} else {
		owa.ArgsUnspecified = true
	}

	return &nodes.DropStmt{
		Objects:    &nodes.List{Items: []nodes.Node{owa}},
		RemoveType: int(nodes.OBJECT_FUNCTION),
		Behavior:   behavior,
		Missing_ok: ifExists,
	}
}

func makeDropProcStmt(schema, name string, argTypes []TypeName, ifExists, cascade bool) *nodes.DropStmt {
	stmt := makeDropFuncStmt(schema, name, argTypes, ifExists, cascade)
	stmt.RemoveType = int(nodes.OBJECT_PROCEDURE)
	return stmt
}

// --- Tests ---

func TestCreateFunction(t *testing.T) {
	c := New()
	err := c.CreateFunctionStmt(makeCreateFuncStmt("", "add_one",
		[]TypeName{{Name: "integer", TypeMod: -1}},
		TypeName{Name: "integer", TypeMod: -1},
		"plpgsql", "BEGIN RETURN x + 1; END;", false))
	if err != nil {
		t.Fatal(err)
	}

	// Should be registered as a BuiltinProc.
	procs := c.LookupProcByName("add_one")
	if len(procs) != 1 {
		t.Fatalf("procs: got %d, want 1", len(procs))
	}
	if procs[0].RetType != INT4OID {
		t.Errorf("ret type: got %d, want %d", procs[0].RetType, INT4OID)
	}
}

func TestCreateFunctionOverload(t *testing.T) {
	c := New()
	c.CreateFunctionStmt(makeCreateFuncStmt("", "myfunc",
		[]TypeName{{Name: "integer", TypeMod: -1}},
		TypeName{Name: "integer", TypeMod: -1},
		"sql", "SELECT $1", false))
	err := c.CreateFunctionStmt(makeCreateFuncStmt("", "myfunc",
		[]TypeName{{Name: "text", TypeMod: -1}},
		TypeName{Name: "text", TypeMod: -1},
		"sql", "SELECT $1", false))
	if err != nil {
		t.Fatal(err)
	}

	procs := c.LookupProcByName("myfunc")
	if len(procs) < 2 {
		t.Fatalf("overloads: got %d, want >= 2", len(procs))
	}
}

func TestCreateFunctionDuplicate(t *testing.T) {
	c := New()
	c.CreateFunctionStmt(makeCreateFuncStmt("", "f",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		TypeName{Name: "int4", TypeMod: -1},
		"sql", "SELECT $1", false))
	err := c.CreateFunctionStmt(makeCreateFuncStmt("", "f",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		TypeName{Name: "int4", TypeMod: -1},
		"sql", "SELECT $1", false))
	assertCode(t, err, CodeDuplicateFunction)
}

func TestCreateFunctionOrReplace(t *testing.T) {
	c := New()
	c.CreateFunctionStmt(makeCreateFuncStmt("", "f",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		TypeName{Name: "int4", TypeMod: -1},
		"sql", "SELECT $1", false))
	err := c.CreateFunctionStmt(makeCreateFuncStmt("", "f",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		TypeName{Name: "int4", TypeMod: -1},
		"sql", "SELECT $1 + 1", true))
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateFunctionOrReplaceWrongReturnType(t *testing.T) {
	c := New()
	c.CreateFunctionStmt(makeCreateFuncStmt("", "f",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		TypeName{Name: "int4", TypeMod: -1},
		"sql", "SELECT $1", false))
	err := c.CreateFunctionStmt(makeCreateFuncStmt("", "f",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		TypeName{Name: "text", TypeMod: -1},
		"sql", "SELECT $1::text", true))
	assertCode(t, err, CodeInvalidObjectDefinition)
}

func TestCreateFunctionOrReplaceWrongRetSet(t *testing.T) {
	c := New()
	// Create function returning a single int4.
	c.CreateFunctionStmt(makeCreateFuncStmt("", "f",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		TypeName{Name: "int4", TypeMod: -1},
		"sql", "SELECT $1", false))
	// Try OR REPLACE with SETOF int4 (same type OID but different retset).
	stmt := makeCreateFuncStmt("", "f",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		TypeName{Name: "int4", TypeMod: -1},
		"sql", "SELECT $1", true)
	stmt.ReturnType.Setof = true
	err := c.CreateFunctionStmt(stmt)
	assertCode(t, err, CodeInvalidObjectDefinition)
}

func TestDropFunction(t *testing.T) {
	c := New()
	c.CreateFunctionStmt(makeCreateFuncStmt("", "f",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		TypeName{Name: "int4", TypeMod: -1},
		"sql", "SELECT $1", false))
	err := c.RemoveObjects(makeDropFuncStmt("", "f",
		[]TypeName{{Name: "int4", TypeMod: -1}}, false, false))
	if err != nil {
		t.Fatal(err)
	}

	procs := c.LookupProcByName("f")
	found := false
	for _, p := range procs {
		if c.userProcs[p.OID] != nil {
			found = true
		}
	}
	if found {
		t.Error("function should be removed")
	}
}

func TestDropFunctionNoArgs(t *testing.T) {
	c := New()
	c.CreateFunctionStmt(makeCreateFuncStmt("", "f",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		TypeName{Name: "int4", TypeMod: -1},
		"sql", "SELECT $1", false))
	// Drop without specifying args -- should find the single overload.
	err := c.RemoveObjects(makeDropFuncStmt("", "f", nil, false, false))
	if err != nil {
		t.Fatal(err)
	}
}

func TestDropFunctionAmbiguous(t *testing.T) {
	c := New()
	c.CreateFunctionStmt(makeCreateFuncStmt("", "f",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		TypeName{Name: "int4", TypeMod: -1},
		"sql", "SELECT $1", false))
	c.CreateFunctionStmt(makeCreateFuncStmt("", "f",
		[]TypeName{{Name: "text", TypeMod: -1}},
		TypeName{Name: "text", TypeMod: -1},
		"sql", "SELECT $1", false))
	err := c.RemoveObjects(makeDropFuncStmt("", "f", nil, false, false))
	assertCode(t, err, CodeAmbiguousFunction)
}

func TestDropFunctionIfExists(t *testing.T) {
	c := New()
	err := c.RemoveObjects(makeDropFuncStmt("", "nope", nil, true, false))
	if err != nil {
		t.Fatal(err)
	}
}

func TestDropFunctionNonexistent(t *testing.T) {
	c := New()
	err := c.RemoveObjects(makeDropFuncStmt("", "nope", nil, false, false))
	assertCode(t, err, CodeUndefinedFunction)
}

func TestCreateProcedure(t *testing.T) {
	c := New()
	// Procedure: no ReturnType (pass empty TypeName so stmt.ReturnType == nil).
	err := c.CreateFunctionStmt(makeCreateFuncStmt("", "myproc",
		nil, TypeName{}, "plpgsql", "BEGIN NULL; END;", false))
	if err != nil {
		t.Fatal(err)
	}

	procs := c.LookupProcByName("myproc")
	if len(procs) != 1 {
		t.Fatalf("procs: got %d, want 1", len(procs))
	}
	if procs[0].Kind != 'p' {
		t.Errorf("kind: got %c, want 'p'", procs[0].Kind)
	}
	if procs[0].RetType != VOIDOID {
		t.Errorf("ret type: got %d, want %d (void)", procs[0].RetType, VOIDOID)
	}
}

func TestDropProcedure(t *testing.T) {
	c := New()
	c.CreateFunctionStmt(makeCreateFuncStmt("", "myproc",
		nil, TypeName{}, "plpgsql", "BEGIN NULL; END;", false))
	err := c.RemoveObjects(makeDropProcStmt("", "myproc", nil, false, false))
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserFunctionResolvableInInfer(t *testing.T) {
	c := New()
	// Create a table.
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{{Name: "x", Type: TypeName{Name: "int4", TypeMod: -1}}}, nil, false), 'r')

	// Create a user function.
	c.CreateFunctionStmt(makeCreateFuncStmt("", "double_it",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		TypeName{Name: "int4", TypeMod: -1},
		"sql", "SELECT $1 * 2", false))

	// Use in a view.
	err := c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			From: []*FromItem{{Kind: FromTable, Table: "t"}},
			TargetList: []*ResTarget{
				{
					Name: "doubled",
					Val: &Expr{
						Kind:     ExprFuncCall,
						FuncName: "double_it",
						FuncArgs: []*Expr{
							{Kind: ExprColumnRef, ColumnName: "x"},
						},
					},
				},
			},
		}, nil))
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "v")
	if rel == nil {
		t.Fatal("view not found")
	}
	if rel.Columns[0].TypeOID != INT4OID {
		t.Errorf("type: got %d, want %d", rel.Columns[0].TypeOID, INT4OID)
	}
}

func TestDropFunctionCascadeDropsTrigger(t *testing.T) {
	c := New()
	// Create trigger function.
	c.CreateFunctionStmt(makeCreateFuncStmt("", "trig_fn",
		nil, TypeName{Name: "trigger", TypeMod: -1},
		"plpgsql", "BEGIN RETURN NEW; END;", false))
	// Create table.
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{{Name: "x", Type: TypeName{Name: "int4", TypeMod: -1}}}, nil, false), 'r')
	// Create trigger.
	c.CreateTriggerStmt(makeCreateTrigStmt("", "t", "my_trig", TriggerBefore, TriggerEventInsert, false, "", "trig_fn"))

	// Drop function without cascade -> should fail.
	err := c.RemoveObjects(makeDropFuncStmt("", "trig_fn", nil, false, false))
	assertCode(t, err, CodeDependentObjects)

	// Drop function with cascade -> should succeed and remove trigger.
	err = c.RemoveObjects(makeDropFuncStmt("", "trig_fn", nil, false, true))
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	trigs := c.TriggersOf(rel.OID)
	if len(trigs) != 0 {
		t.Errorf("triggers after cascade: got %d, want 0", len(trigs))
	}
}
