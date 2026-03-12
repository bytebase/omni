package ddl

import (
	"errors"
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

func makeCreateSchemaStmt(name string, ifNotExists bool) *nodes.CreateSchemaStmt {
	return &nodes.CreateSchemaStmt{
		Schemaname:  name,
		IfNotExists: ifNotExists,
	}
}

func makeDropSchemaStmt(name string, ifExists bool, cascade bool) *nodes.DropStmt {
	behavior := int(nodes.DROP_RESTRICT)
	if cascade {
		behavior = int(nodes.DROP_CASCADE)
	}
	return &nodes.DropStmt{
		Objects:    &nodes.List{Items: []nodes.Node{&nodes.String{Str: name}}},
		RemoveType: int(nodes.OBJECT_SCHEMA),
		Behavior:   behavior,
		Missing_ok: ifExists,
	}
}

func TestCreateSchema(t *testing.T) {
	c := New()

	if err := c.CreateSchemaCommand(makeCreateSchemaStmt("myschema", false)); err != nil {
		t.Fatal(err)
	}
	if s := c.GetSchema("myschema"); s == nil {
		t.Fatal("schema not found after CREATE")
	}
}

func TestCreateSchemaDuplicate(t *testing.T) {
	c := New()

	c.CreateSchemaCommand(makeCreateSchemaStmt("dup", false))
	err := c.CreateSchemaCommand(makeCreateSchemaStmt("dup", false))
	assertErrorCode(t, err, CodeDuplicateSchema)
}

func TestCreateSchemaIfNotExists(t *testing.T) {
	c := New()

	c.CreateSchemaCommand(makeCreateSchemaStmt("dup", false))
	if err := c.CreateSchemaCommand(makeCreateSchemaStmt("dup", true)); err != nil {
		t.Fatalf("IF NOT EXISTS should not error, got: %v", err)
	}
}

func TestDropSchema(t *testing.T) {
	c := New()

	c.CreateSchemaCommand(makeCreateSchemaStmt("todrop", false))
	if err := c.RemoveSchemas(makeDropSchemaStmt("todrop", false, false)); err != nil {
		t.Fatal(err)
	}
	if s := c.GetSchema("todrop"); s != nil {
		t.Fatal("schema still exists after DROP")
	}
}

func TestDropSchemaNotEmpty(t *testing.T) {
	c := New()

	c.CreateSchemaCommand(makeCreateSchemaStmt("notempty", false))
	c.DefineRelation(makeCreateTableStmt("notempty", "t", []ColumnDef{{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}}, nil, false), 'r')

	err := c.RemoveSchemas(makeDropSchemaStmt("notempty", false, false))
	assertErrorCode(t, err, CodeSchemaNotEmpty)
}

func TestDropSchemaCascade(t *testing.T) {
	c := New()

	c.CreateSchemaCommand(makeCreateSchemaStmt("casc", false))
	c.DefineRelation(makeCreateTableStmt("casc", "t", []ColumnDef{{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}}, nil, false), 'r')

	if err := c.RemoveSchemas(makeDropSchemaStmt("casc", false, true)); err != nil {
		t.Fatal(err)
	}
	if s := c.GetSchema("casc"); s != nil {
		t.Fatal("schema still exists after DROP CASCADE")
	}
}

func TestDropBuiltinSchema(t *testing.T) {
	c := New()

	err := c.RemoveSchemas(makeDropSchemaStmt("pg_catalog", false, false))
	if err == nil {
		t.Fatal("expected error dropping pg_catalog")
	}
}

func TestDropSchemaIfExists(t *testing.T) {
	c := New()

	if err := c.RemoveSchemas(makeDropSchemaStmt("nosuch", true, false)); err != nil {
		t.Fatalf("IF EXISTS should not error, got: %v", err)
	}
}

func TestDropSchemaNonExistent(t *testing.T) {
	c := New()

	err := c.RemoveSchemas(makeDropSchemaStmt("nosuch", false, false))
	assertErrorCode(t, err, CodeUndefinedSchema)
}

func TestDropSchemaCascadeDropsEnum(t *testing.T) {
	c := New()
	c.CreateSchemaCommand(makeCreateSchemaStmt("s", false))
	c.DefineEnum(makeCreateEnumStmt("s", "mood", []string{"happy", "sad"}))

	// Schema not empty.
	err := c.RemoveSchemas(makeDropSchemaStmt("s", false, false))
	assertErrorCode(t, err, CodeSchemaNotEmpty)

	// CASCADE should drop enum.
	err = c.RemoveSchemas(makeDropSchemaStmt("s", false, true))
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = c.ResolveType(TypeName{Schema: "s", Name: "mood", TypeMod: -1})
	if err == nil {
		t.Error("enum type should be gone after schema cascade")
	}
}

func TestDropSchemaCascadeDropsDomain(t *testing.T) {
	c := New()
	c.CreateSchemaCommand(makeCreateSchemaStmt("s", false))
	c.DefineDomain(makeCreateDomainStmt("s", "posint", TypeName{Name: "integer", TypeMod: -1}, false, "", "", ""))

	err := c.RemoveSchemas(makeDropSchemaStmt("s", false, false))
	assertErrorCode(t, err, CodeSchemaNotEmpty)

	err = c.RemoveSchemas(makeDropSchemaStmt("s", false, true))
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = c.ResolveType(TypeName{Schema: "s", Name: "posint", TypeMod: -1})
	if err == nil {
		t.Error("domain type should be gone after schema cascade")
	}
}

func TestDropSchemaCascadeDropsFunction(t *testing.T) {
	c := New()
	c.CreateSchemaCommand(makeCreateSchemaStmt("s", false))
	c.CreateFunctionStmt(makeCreateFuncStmt("s", "myfunc",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		TypeName{Name: "int4", TypeMod: -1},
		"sql", "SELECT $1", false))

	err := c.RemoveSchemas(makeDropSchemaStmt("s", false, false))
	assertErrorCode(t, err, CodeSchemaNotEmpty)

	err = c.RemoveSchemas(makeDropSchemaStmt("s", false, true))
	if err != nil {
		t.Fatal(err)
	}

	// Function should no longer be resolvable as a user proc.
	for _, up := range c.userProcs {
		if up.Name == "myfunc" {
			t.Error("function should be gone after schema cascade")
		}
	}
}

func TestDropSchemaCascadeDropsSequence(t *testing.T) {
	c := New()
	c.CreateSchemaCommand(makeCreateSchemaStmt("s", false))
	c.DefineSequence(makeCreateSeqStmt("s", "myseq", false))

	err := c.RemoveSchemas(makeDropSchemaStmt("s", false, false))
	assertErrorCode(t, err, CodeSchemaNotEmpty)

	err = c.RemoveSchemas(makeDropSchemaStmt("s", false, true))
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.findSequence("s", "myseq")
	if err == nil {
		t.Error("sequence should be gone after schema cascade")
	}
}

func TestDropSchemaCascadeDropsAll(t *testing.T) {
	c := New()
	c.CreateSchemaCommand(makeCreateSchemaStmt("s", false))

	// Create one of each object type.
	c.DefineEnum(makeCreateEnumStmt("s", "mood", []string{"ok"}))
	c.DefineDomain(makeCreateDomainStmt("s", "posint", TypeName{Name: "int4", TypeMod: -1}, false, "", "", ""))
	c.DefineSequence(makeCreateSeqStmt("s", "myseq", false))
	c.DefineRelation(makeCreateTableStmt("s", "t", []ColumnDef{{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}}, nil, false), 'r')
	c.CreateFunctionStmt(makeCreateFuncStmt("s", "trig_fn",
		nil, TypeName{Name: "trigger", TypeMod: -1},
		"plpgsql", "BEGIN RETURN NEW; END;", false))
	c.CreateTriggerStmt(makeCreateTrigStmt("s", "t", "my_trig", TriggerBefore, TriggerEventInsert, false, "", "trig_fn"))

	err := c.RemoveSchemas(makeDropSchemaStmt("s", false, true))
	if err != nil {
		t.Fatal(err)
	}

	if c.GetSchema("s") != nil {
		t.Error("schema should be gone")
	}
	if len(c.triggers) != 0 {
		t.Errorf("triggers remaining: %d", len(c.triggers))
	}
}

// --- Table AST helpers ---

func makeTypeNameNode(tn TypeName) *nodes.TypeName {
	var items []nodes.Node
	if tn.Schema != "" {
		items = append(items, &nodes.String{Str: tn.Schema})
	}
	items = append(items, &nodes.String{Str: tn.Name})
	result := &nodes.TypeName{
		Names:   &nodes.List{Items: items},
		Typemod: tn.TypeMod,
	}
	if tn.IsArray {
		result.ArrayBounds = &nodes.List{Items: []nodes.Node{&nodes.Integer{Ival: -1}}}
	}
	return result
}

func makeColumnDefNode(cd ColumnDef) *nodes.ColumnDef {
	result := &nodes.ColumnDef{
		Colname:   cd.Name,
		IsNotNull: cd.NotNull,
	}
	if cd.IsSerial != 0 {
		var sname string
		switch cd.IsSerial {
		case 2:
			sname = "smallserial"
		case 4:
			sname = "serial"
		default:
			sname = "bigserial"
		}
		result.TypeName = &nodes.TypeName{
			Names:   &nodes.List{Items: []nodes.Node{&nodes.String{Str: sname}}},
			Typemod: -1,
		}
	} else {
		result.TypeName = makeTypeNameNode(cd.Type)
	}
	if cd.Default != "" {
		var cons []nodes.Node
		if result.Constraints != nil {
			cons = result.Constraints.Items
		}
		cons = append(cons, &nodes.Constraint{
			Contype:    nodes.CONSTR_DEFAULT,
			CookedExpr: cd.Default,
		})
		result.Constraints = &nodes.List{Items: cons}
	}
	return result
}

func makeStringListNode(ss []string) *nodes.List {
	if len(ss) == 0 {
		return nil
	}
	items := make([]nodes.Node, len(ss))
	for i, s := range ss {
		items[i] = &nodes.String{Str: s}
	}
	return &nodes.List{Items: items}
}

func makeConstraintNode(cd ConstraintDef) *nodes.Constraint {
	con := &nodes.Constraint{Conname: cd.Name}
	switch cd.Type {
	case ConstraintPK:
		con.Contype = nodes.CONSTR_PRIMARY
		con.Keys = makeStringListNode(cd.Columns)
	case ConstraintUnique:
		con.Contype = nodes.CONSTR_UNIQUE
		con.Keys = makeStringListNode(cd.Columns)
	case ConstraintFK:
		con.Contype = nodes.CONSTR_FOREIGN
		con.FkAttrs = makeStringListNode(cd.Columns)
		con.PkAttrs = makeStringListNode(cd.RefColumns)
		if cd.RefTable != "" {
			con.Pktable = &nodes.RangeVar{
				Schemaname: cd.RefSchema,
				Relname:    cd.RefTable,
			}
		}
	case ConstraintCheck:
		con.Contype = nodes.CONSTR_CHECK
		con.CookedExpr = cd.CheckExpr
	}
	return con
}

func makeCreateTableStmt(schema, name string, cols []ColumnDef, cons []ConstraintDef, ifNotExists bool) *nodes.CreateStmt {
	stmt := &nodes.CreateStmt{
		Relation:    &nodes.RangeVar{Schemaname: schema, Relname: name},
		IfNotExists: ifNotExists,
	}
	var elts []nodes.Node
	for _, cd := range cols {
		elts = append(elts, makeColumnDefNode(cd))
	}
	if len(elts) > 0 {
		stmt.TableElts = &nodes.List{Items: elts}
	}
	var conNodes []nodes.Node
	for _, cd := range cons {
		conNodes = append(conNodes, makeConstraintNode(cd))
	}
	if len(conNodes) > 0 {
		stmt.Constraints = &nodes.List{Items: conNodes}
	}
	return stmt
}

func makeDropTableStmt(schema, name string, ifExists, cascade bool) *nodes.DropStmt {
	behavior := int(nodes.DROP_RESTRICT)
	if cascade {
		behavior = int(nodes.DROP_CASCADE)
	}
	var items []nodes.Node
	if schema != "" {
		items = []nodes.Node{&nodes.String{Str: schema}, &nodes.String{Str: name}}
	} else {
		items = []nodes.Node{&nodes.String{Str: name}}
	}
	return &nodes.DropStmt{
		Objects:    &nodes.List{Items: []nodes.Node{&nodes.List{Items: items}}},
		RemoveType: int(nodes.OBJECT_TABLE),
		Behavior:   behavior,
		Missing_ok: ifExists,
	}
}

// assertErrorCode checks that the error is a catalog.Error with the expected code.
func assertErrorCode(t *testing.T, err error, code string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error with code %s, got nil", code)
	}
	var catErr *Error
	if !errors.As(err, &catErr) {
		t.Fatalf("expected *catalog.Error, got %T: %v", err, err)
	}
	if catErr.Code != code {
		t.Errorf("error code: got %s, want %s (msg: %s)", catErr.Code, code, catErr.Message)
	}
}
