package ddl

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// makeCreateTrigStmt builds a CreateTrigStmt from test parameters.
func makeCreateTrigStmt(schema, table, name string, timing TriggerTiming, events TriggerEvent, forEachRow bool, funcSchema, funcName string) *nodes.CreateTrigStmt {
	var tm int16
	switch timing {
	case TriggerBefore:
		tm = nodes.TRIGGER_TYPE_BEFORE
	case TriggerInsteadOf:
		tm = nodes.TRIGGER_TYPE_INSTEAD
	default:
		tm = nodes.TRIGGER_TYPE_AFTER
	}

	var ev int16
	if events&TriggerEventInsert != 0 {
		ev |= nodes.TRIGGER_TYPE_INSERT
	}
	if events&TriggerEventUpdate != 0 {
		ev |= nodes.TRIGGER_TYPE_UPDATE
	}
	if events&TriggerEventDelete != 0 {
		ev |= nodes.TRIGGER_TYPE_DELETE
	}
	if events&TriggerEventTruncate != 0 {
		ev |= nodes.TRIGGER_TYPE_TRUNCATE
	}

	var fnItems []nodes.Node
	if funcSchema != "" {
		fnItems = append(fnItems, &nodes.String{Str: funcSchema})
	}
	fnItems = append(fnItems, &nodes.String{Str: funcName})

	return &nodes.CreateTrigStmt{
		Trigname: name,
		Relation: &nodes.RangeVar{Schemaname: schema, Relname: table},
		Funcname: &nodes.List{Items: fnItems},
		Row:      forEachRow,
		Timing:   tm,
		Events:   ev,
	}
}

// makeCreateTrigStmtReplace builds a CreateTrigStmt with Replace=true.
func makeCreateTrigStmtReplace(schema, table, name string, timing TriggerTiming, events TriggerEvent, forEachRow bool, funcSchema, funcName string) *nodes.CreateTrigStmt {
	stmt := makeCreateTrigStmt(schema, table, name, timing, events, forEachRow, funcSchema, funcName)
	stmt.Replace = true
	return stmt
}

// makeCreateTrigStmtWithColumns builds a CreateTrigStmt with UPDATE OF columns.
func makeCreateTrigStmtWithColumns(schema, table, name string, timing TriggerTiming, events TriggerEvent, forEachRow bool, funcSchema, funcName string, columns []string) *nodes.CreateTrigStmt {
	stmt := makeCreateTrigStmt(schema, table, name, timing, events, forEachRow, funcSchema, funcName)
	if len(columns) > 0 {
		var items []nodes.Node
		for _, col := range columns {
			items = append(items, &nodes.String{Str: col})
		}
		stmt.Columns = &nodes.List{Items: items}
	}
	return stmt
}

// makeDropTrigStmt builds a DropStmt for DROP TRIGGER.
func makeDropTrigStmt(schema, table, name string, ifExists, cascade bool) *nodes.DropStmt {
	behavior := int(nodes.DROP_RESTRICT)
	if cascade {
		behavior = int(nodes.DROP_CASCADE)
	}

	var items []nodes.Node
	if schema != "" {
		items = []nodes.Node{&nodes.String{Str: schema}, &nodes.String{Str: table}, &nodes.String{Str: name}}
	} else {
		items = []nodes.Node{&nodes.String{Str: table}, &nodes.String{Str: name}}
	}

	return &nodes.DropStmt{
		Objects:    &nodes.List{Items: []nodes.Node{&nodes.List{Items: items}}},
		RemoveType: int(nodes.OBJECT_TRIGGER),
		Behavior:   behavior,
		Missing_ok: ifExists,
	}
}

// makeTriggerFuncStmt builds a CreateFunctionStmt for a trigger function.
func makeTriggerFuncStmt(schema, name string, language, body string) *nodes.CreateFunctionStmt {
	var nameItems []nodes.Node
	if schema != "" {
		nameItems = append(nameItems, &nodes.String{Str: schema})
	}
	nameItems = append(nameItems, &nodes.String{Str: name})

	var opts []nodes.Node
	opts = append(opts, &nodes.DefElem{Defname: "language", Arg: &nodes.String{Str: language}})
	opts = append(opts, &nodes.DefElem{Defname: "as", Arg: &nodes.List{Items: []nodes.Node{&nodes.String{Str: body}}}})

	return &nodes.CreateFunctionStmt{
		Funcname:   &nodes.List{Items: nameItems},
		ReturnType: makeTypeNameNode(TypeName{Name: "trigger", TypeMod: -1}),
		Options:    &nodes.List{Items: opts},
	}
}

// makeNonTriggerFuncStmt builds a CreateFunctionStmt for a function with a custom return type.
func makeNonTriggerFuncStmt(schema, name string, retType TypeName, language, body string) *nodes.CreateFunctionStmt {
	var nameItems []nodes.Node
	if schema != "" {
		nameItems = append(nameItems, &nodes.String{Str: schema})
	}
	nameItems = append(nameItems, &nodes.String{Str: name})

	var opts []nodes.Node
	opts = append(opts, &nodes.DefElem{Defname: "language", Arg: &nodes.String{Str: language}})
	opts = append(opts, &nodes.DefElem{Defname: "as", Arg: &nodes.List{Items: []nodes.Node{&nodes.String{Str: body}}}})

	return &nodes.CreateFunctionStmt{
		Funcname:   &nodes.List{Items: nameItems},
		ReturnType: makeTypeNameNode(retType),
		Options:    &nodes.List{Items: opts},
	}
}

func createTriggerFunc(t *testing.T, c *Catalog) {
	t.Helper()
	err := c.CreateFunctionStmt(makeTriggerFuncStmt("", "trig_fn", "plpgsql", "BEGIN RETURN NEW; END;"))
	if err != nil {
		t.Fatal(err)
	}
}

func createTestTable(t *testing.T, c *Catalog) {
	t.Helper()
	err := c.DefineRelation(makeCreateTableStmt("", "t",
		[]ColumnDef{{Name: "x", Type: TypeName{Name: "int4", TypeMod: -1}}},
		nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateTriggerBeforeInsert(t *testing.T) {
	c := New()
	createTestTable(t, c)
	createTriggerFunc(t, c)

	err := c.CreateTriggerStmt(makeCreateTrigStmt("", "t", "my_trig", TriggerBefore, TriggerEventInsert, true, "", "trig_fn"))
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	trigs := c.TriggersOf(rel.OID)
	if len(trigs) != 1 {
		t.Fatalf("triggers: got %d, want 1", len(trigs))
	}
	if trigs[0].Timing != TriggerBefore {
		t.Errorf("timing: got %c, want %c", trigs[0].Timing, TriggerBefore)
	}
	if trigs[0].Events != TriggerEventInsert {
		t.Errorf("events: got %d, want %d", trigs[0].Events, TriggerEventInsert)
	}
	if !trigs[0].ForEachRow {
		t.Error("should be FOR EACH ROW")
	}
}

func TestCreateTriggerAfterUpdate(t *testing.T) {
	c := New()
	createTestTable(t, c)
	createTriggerFunc(t, c)

	err := c.CreateTriggerStmt(makeCreateTrigStmt("", "t", "my_trig", TriggerAfter, TriggerEventUpdate, false, "", "trig_fn"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateTriggerMultipleEvents(t *testing.T) {
	c := New()
	createTestTable(t, c)
	createTriggerFunc(t, c)

	err := c.CreateTriggerStmt(makeCreateTrigStmt("", "t", "my_trig", TriggerBefore, TriggerEventInsert|TriggerEventUpdate|TriggerEventDelete, false, "", "trig_fn"))
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	trigs := c.TriggersOf(rel.OID)
	expected := TriggerEventInsert | TriggerEventUpdate | TriggerEventDelete
	if trigs[0].Events != expected {
		t.Errorf("events: got %d, want %d", trigs[0].Events, expected)
	}
}

func TestCreateTriggerDuplicate(t *testing.T) {
	c := New()
	createTestTable(t, c)
	createTriggerFunc(t, c)

	c.CreateTriggerStmt(makeCreateTrigStmt("", "t", "my_trig", TriggerBefore, TriggerEventInsert, false, "", "trig_fn"))
	err := c.CreateTriggerStmt(makeCreateTrigStmt("", "t", "my_trig", TriggerBefore, TriggerEventInsert, false, "", "trig_fn"))
	assertCode(t, err, CodeDuplicateObject)
}

func TestCreateTriggerOrReplace(t *testing.T) {
	c := New()
	createTestTable(t, c)
	createTriggerFunc(t, c)

	c.CreateTriggerStmt(makeCreateTrigStmt("", "t", "my_trig", TriggerBefore, TriggerEventInsert, false, "", "trig_fn"))
	err := c.CreateTriggerStmt(makeCreateTrigStmtReplace("", "t", "my_trig", TriggerAfter, TriggerEventUpdate, false, "", "trig_fn"))
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	trigs := c.TriggersOf(rel.OID)
	if len(trigs) != 1 {
		t.Fatalf("triggers: got %d, want 1", len(trigs))
	}
	if trigs[0].Timing != TriggerAfter {
		t.Errorf("timing: got %c, want %c", trigs[0].Timing, TriggerAfter)
	}
}

func TestCreateTriggerNoFunction(t *testing.T) {
	c := New()
	createTestTable(t, c)

	err := c.CreateTriggerStmt(makeCreateTrigStmt("", "t", "my_trig", TriggerBefore, TriggerEventInsert, false, "", "nope"))
	assertCode(t, err, CodeUndefinedFunction)
}

func TestCreateTriggerWrongReturnType(t *testing.T) {
	c := New()
	createTestTable(t, c)

	// Create a function that returns int4 (not trigger).
	c.CreateFunctionStmt(makeNonTriggerFuncStmt("", "bad_fn", TypeName{Name: "int4", TypeMod: -1}, "sql", "SELECT 1"))

	err := c.CreateTriggerStmt(makeCreateTrigStmt("", "t", "my_trig", TriggerBefore, TriggerEventInsert, false, "", "bad_fn"))
	assertCode(t, err, CodeInvalidObjectDefinition)
}

func TestCreateTriggerInsteadOfOnTable(t *testing.T) {
	c := New()
	createTestTable(t, c)
	createTriggerFunc(t, c)

	err := c.CreateTriggerStmt(makeCreateTrigStmt("", "t", "my_trig", TriggerInsteadOf, TriggerEventInsert, false, "", "trig_fn"))
	assertCode(t, err, CodeInvalidObjectDefinition)
}

func TestCreateTriggerBeforeOnView(t *testing.T) {
	c := New()
	createTestTable(t, c)
	createTriggerFunc(t, c)

	c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			From:       []*FromItem{{Kind: FromTable, Table: "t"}},
			TargetList: []*ResTarget{{Val: &Expr{Kind: ExprStar}}},
		}, nil))

	// ROW-level BEFORE triggers are rejected on views; STATEMENT-level is OK.
	err := c.CreateTriggerStmt(makeCreateTrigStmt("", "v", "my_trig", TriggerBefore, TriggerEventInsert, true, "", "trig_fn"))
	assertCode(t, err, CodeInvalidObjectDefinition)
}

func TestCreateTriggerInsteadOfOnView(t *testing.T) {
	c := New()
	createTestTable(t, c)
	createTriggerFunc(t, c)

	c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			From:       []*FromItem{{Kind: FromTable, Table: "t"}},
			TargetList: []*ResTarget{{Val: &Expr{Kind: ExprStar}}},
		}, nil))

	// INSTEAD OF triggers must be FOR EACH ROW (PG requirement).
	err := c.CreateTriggerStmt(makeCreateTrigStmt("", "v", "my_trig", TriggerInsteadOf, TriggerEventInsert, true, "", "trig_fn"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestDropTrigger(t *testing.T) {
	c := New()
	createTestTable(t, c)
	createTriggerFunc(t, c)

	c.CreateTriggerStmt(makeCreateTrigStmt("", "t", "my_trig", TriggerBefore, TriggerEventInsert, false, "", "trig_fn"))

	err := c.RemoveObjects(makeDropTrigStmt("", "t", "my_trig", false, false))
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	if len(c.TriggersOf(rel.OID)) != 0 {
		t.Error("trigger should be removed")
	}
}

func TestDropTriggerIfExists(t *testing.T) {
	c := New()
	createTestTable(t, c)

	err := c.RemoveObjects(makeDropTrigStmt("", "t", "nope", true, false))
	if err != nil {
		t.Fatal(err)
	}
}

func TestDropTriggerNonexistent(t *testing.T) {
	c := New()
	createTestTable(t, c)

	err := c.RemoveObjects(makeDropTrigStmt("", "t", "nope", false, false))
	assertCode(t, err, CodeUndefinedObject)
}

func TestDropTableCascadeDropsTriggers(t *testing.T) {
	c := New()
	createTestTable(t, c)
	createTriggerFunc(t, c)

	c.CreateTriggerStmt(makeCreateTrigStmt("", "t", "my_trig", TriggerBefore, TriggerEventInsert, false, "", "trig_fn"))

	err := c.RemoveRelations(makeDropTableStmt("", "t", false, true))
	if err != nil {
		t.Fatal(err)
	}

	// Trigger should no longer be in the catalog.
	if len(c.triggers) != 0 {
		t.Errorf("triggers in catalog: got %d, want 0", len(c.triggers))
	}
}

func TestCreateTriggerUpdateOfColumns(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t",
		[]ColumnDef{
			{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
			{Name: "b", Type: TypeName{Name: "text", TypeMod: -1}},
		},
		nil, false), 'r')
	createTriggerFunc(t, c)

	err := c.CreateTriggerStmt(makeCreateTrigStmtWithColumns("", "t", "my_trig", TriggerBefore, TriggerEventUpdate, false, "", "trig_fn", []string{"b"}))
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	trigs := c.TriggersOf(rel.OID)
	if len(trigs[0].Columns) != 1 {
		t.Fatalf("trigger columns: got %d, want 1", len(trigs[0].Columns))
	}
	if trigs[0].Columns[0] != 2 {
		t.Errorf("trigger column attnum: got %d, want 2", trigs[0].Columns[0])
	}
}
