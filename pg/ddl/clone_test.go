package ddl

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// buildRichCatalog creates a catalog with a variety of objects for clone testing.
func buildRichCatalog(t *testing.T) *Catalog {
	t.Helper()
	c := New()

	// Schema.
	if err := c.CreateSchemaCommand(&nodes.CreateSchemaStmt{Schemaname: "myschema"}); err != nil {
		t.Fatal(err)
	}

	// Table with columns + PK + CHECK.
	if err := c.DefineRelation(makeCreateTableStmt("", "users", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "integer", TypeMod: -1}, NotNull: true},
		{Name: "name", Type: TypeName{Name: "text", TypeMod: -1}},
		{Name: "email", Type: TypeName{Name: "text", TypeMod: -1}},
	}, []ConstraintDef{
		{Type: ConstraintPK, Columns: []string{"id"}},
		{Type: ConstraintCheck, Columns: []string{"name"}, CheckExpr: "name IS NOT NULL"},
	}, false), 'r'); err != nil {
		t.Fatal(err)
	}

	// Second table with FK to users.
	if err := c.DefineRelation(makeCreateTableStmt("", "orders", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "integer", TypeMod: -1}, NotNull: true},
		{Name: "user_id", Type: TypeName{Name: "integer", TypeMod: -1}},
		{Name: "amount", Type: TypeName{Name: "numeric", TypeMod: -1}},
	}, []ConstraintDef{
		{Type: ConstraintPK, Columns: []string{"id"}},
		{Type: ConstraintFK, Columns: []string{"user_id"}, RefTable: "users", RefColumns: []string{"id"}},
	}, false), 'r'); err != nil {
		t.Fatal(err)
	}

	// Index on orders.
	if err := c.DefineIndex(makeIndexStmt("", "orders", "orders_amount_idx", []string{"amount"}, false, false)); err != nil {
		t.Fatal(err)
	}

	// Sequence.
	if err := c.DefineSequence(&nodes.CreateSeqStmt{
		Sequence: &nodes.RangeVar{Relname: "my_seq"},
	}); err != nil {
		t.Fatal(err)
	}

	// Enum type.
	if err := c.DefineEnum(makeCreateEnumStmt("", "color", []string{"red", "green", "blue"})); err != nil {
		t.Fatal(err)
	}

	// Domain type.
	if err := c.DefineDomain(&nodes.CreateDomainStmt{
		Domainname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "posint"}}},
		Typname:    &nodes.TypeName{Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "integer"}}}},
	}); err != nil {
		t.Fatal(err)
	}

	// Regular function (int -> int).
	if err := c.CreateFunctionStmt(&nodes.CreateFunctionStmt{
		Funcname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "add_one"}}},
		Parameters: &nodes.List{Items: []nodes.Node{
			&nodes.FunctionParameter{
				Name:    "x",
				ArgType: &nodes.TypeName{Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "integer"}}}},
				Mode:    nodes.FUNC_PARAM_IN,
			},
		}},
		ReturnType: &nodes.TypeName{Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "integer"}}}},
		Options: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{Defname: "language", Arg: &nodes.String{Str: "sql"}},
			&nodes.DefElem{Defname: "as", Arg: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "SELECT x + 1"}}}},
		}},
	}); err != nil {
		t.Fatal(err)
	}

	// Trigger function (returns trigger, 0 args).
	if err := c.CreateFunctionStmt(makeTriggerFuncStmt("", "trig_fn", "plpgsql", "BEGIN RETURN NEW; END;")); err != nil {
		t.Fatal(err)
	}

	// Comment on table.
	if err := c.CommentObject(&nodes.CommentStmt{
		Objtype: nodes.OBJECT_TABLE,
		Object:  &nodes.List{Items: []nodes.Node{&nodes.String{Str: "users"}}},
		Comment: "Users table",
	}); err != nil {
		t.Fatal(err)
	}

	// Trigger on users.
	if err := c.CreateTriggerStmt(makeCreateTrigStmt(
		"", "users", "users_trg",
		TriggerBefore, TriggerEventInsert, true,
		"", "trig_fn",
	)); err != nil {
		t.Fatal(err)
	}

	return c
}

func TestClonePreservesState(t *testing.T) {
	orig := buildRichCatalog(t)
	clone := orig.Clone()

	// Verify schemas.
	if s := clone.GetSchema("myschema"); s == nil {
		t.Fatal("clone missing schema myschema")
	}
	if s := clone.GetSchema("public"); s == nil {
		t.Fatal("clone missing schema public")
	}

	// Verify relations.
	users := clone.GetRelation("", "users")
	if users == nil {
		t.Fatal("clone missing users table")
	}
	if len(users.Columns) != 3 {
		t.Errorf("users columns: got %d, want 3", len(users.Columns))
	}

	orders := clone.GetRelation("", "orders")
	if orders == nil {
		t.Fatal("clone missing orders table")
	}

	// Verify constraints.
	usersCons := clone.ConstraintsOf(users.OID)
	if len(usersCons) < 2 {
		t.Errorf("users constraints: got %d, want >= 2", len(usersCons))
	}

	// Verify index.
	orderIdxs := clone.IndexesOf(orders.OID)
	foundAmountIdx := false
	for _, idx := range orderIdxs {
		if idx.Name == "orders_amount_idx" {
			foundAmountIdx = true
		}
	}
	if !foundAmountIdx {
		t.Error("clone missing orders_amount_idx")
	}

	// Verify sequence exists by checking internal map.
	foundSeq := false
	for _, seq := range clone.sequenceByOID {
		if seq.Name == "my_seq" {
			foundSeq = true
		}
	}
	if !foundSeq {
		t.Error("clone missing my_seq sequence")
	}

	// Verify enum type.
	enumRows := clone.QueryPgType("public")
	foundEnum := false
	for _, r := range enumRows {
		if r.TypName == "color" {
			foundEnum = true
		}
	}
	if !foundEnum {
		t.Error("clone missing color enum type")
	}

	// Verify function.
	procRows := clone.QueryPgProc("public")
	foundFunc := false
	for _, r := range procRows {
		if r.ProName == "add_one" {
			foundFunc = true
		}
	}
	if !foundFunc {
		t.Error("clone missing add_one function")
	}

	// Verify trigger.
	trigRows := clone.QueryPgTrigger(users.OID)
	foundTrig := false
	for _, r := range trigRows {
		if r.TgName == "users_trg" {
			foundTrig = true
		}
	}
	if !foundTrig {
		t.Error("clone missing users_trg trigger")
	}

	// Verify comment.
	descRows := clone.QueryPgDescription()
	foundComment := false
	for _, r := range descRows {
		if r.ObjOID == users.OID && r.Description == "Users table" {
			foundComment = true
		}
	}
	if !foundComment {
		t.Error("clone missing comment on users table")
	}
}

func TestCloneIsolation_ModifyClone(t *testing.T) {
	orig := buildRichCatalog(t)
	clone := orig.Clone()

	// Add a table to the clone.
	err := clone.DefineRelation(makeCreateTableStmt("", "extra", []ColumnDef{
		{Name: "x", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	// Clone should have it.
	if r := clone.GetRelation("", "extra"); r == nil {
		t.Fatal("clone should have 'extra' table")
	}

	// Original should NOT have it.
	if r := orig.GetRelation("", "extra"); r != nil {
		t.Fatal("original should NOT have 'extra' table")
	}
}

func TestCloneIsolation_ModifyOriginal(t *testing.T) {
	orig := buildRichCatalog(t)
	clone := orig.Clone()

	// Add a table to the original.
	err := orig.DefineRelation(makeCreateTableStmt("", "extra2", []ColumnDef{
		{Name: "y", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	// Original should have it.
	if r := orig.GetRelation("", "extra2"); r == nil {
		t.Fatal("original should have 'extra2' table")
	}

	// Clone should NOT have it.
	if r := clone.GetRelation("", "extra2"); r != nil {
		t.Fatal("clone should NOT have 'extra2' table")
	}
}

func TestCloneSchemaPointerIntegrity(t *testing.T) {
	orig := buildRichCatalog(t)
	clone := orig.Clone()

	// Verify that Schema pointers on cloned relations point to cloned schemas, not original.
	users := clone.GetRelation("", "users")
	if users == nil {
		t.Fatal("clone missing users")
	}
	clonePublic := clone.GetSchema("public")
	if users.Schema != clonePublic {
		t.Error("cloned relation's Schema pointer does not match cloned schema")
	}

	// Same for indexes.
	for _, idx := range clone.IndexesOf(users.OID) {
		if idx.Schema != clonePublic {
			t.Errorf("cloned index %s Schema pointer mismatch", idx.Name)
		}
	}

	// Same for sequences.
	for _, seq := range clone.sequenceByOID {
		if seq.Schema.OID == PublicNamespace && seq.Schema != clonePublic {
			t.Errorf("cloned sequence %s Schema pointer mismatch", seq.Name)
		}
	}
}

func TestCloneOIDGeneratorIsolation(t *testing.T) {
	orig := buildRichCatalog(t)
	clone := orig.Clone()

	// Get OID generators.
	origNext := orig.oidGen.Next()
	cloneNext := clone.oidGen.Next()

	if origNext != cloneNext {
		t.Errorf("oid generators should start at same value: orig=%d, clone=%d", origNext, cloneNext)
	}

	// Advance clone further.
	clone.oidGen.Next()
	clone.oidGen.Next()

	// Original should not be affected.
	origNext2 := orig.oidGen.Next()
	if origNext2 != origNext+1 {
		t.Errorf("original oid generator affected by clone: got %d, want %d", origNext2, origNext+1)
	}
}

func TestCloneBuiltinTypesShared(t *testing.T) {
	orig := New()
	clone := orig.Clone()

	// Builtin types should be the same pointer.
	origBool := orig.TypeByOID(BOOLOID)
	cloneBool := clone.TypeByOID(BOOLOID)
	if origBool != cloneBool {
		t.Error("builtin type BOOL should be shared between original and clone")
	}
}
