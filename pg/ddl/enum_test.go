package ddl

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

func makeCreateEnumStmt(schema, name string, values []string) *nodes.CreateEnumStmt {
	var nameItems []nodes.Node
	if schema != "" {
		nameItems = []nodes.Node{&nodes.String{Str: schema}, &nodes.String{Str: name}}
	} else {
		nameItems = []nodes.Node{&nodes.String{Str: name}}
	}
	return &nodes.CreateEnumStmt{
		TypeName: &nodes.List{Items: nameItems},
		Vals:     makeStringListNode(values),
	}
}

func makeAlterEnumStmt(schema, name, newValue, before, after string, ifNotExists bool) *nodes.AlterEnumStmt {
	var nameItems []nodes.Node
	if schema != "" {
		nameItems = []nodes.Node{&nodes.String{Str: schema}, &nodes.String{Str: name}}
	} else {
		nameItems = []nodes.Node{&nodes.String{Str: name}}
	}
	stmt := &nodes.AlterEnumStmt{
		Typname:            &nodes.List{Items: nameItems},
		Newval:             newValue,
		SkipIfNewvalExists: ifNotExists,
	}
	if after != "" {
		stmt.NewvalNeighbor = after
		stmt.NewvalIsAfter = true
	} else if before != "" {
		stmt.NewvalNeighbor = before
		stmt.NewvalIsAfter = false
	}
	return stmt
}

func makeDropTypeStmt(schema, name string, ifExists, cascade bool) *nodes.DropStmt {
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
		RemoveType: int(nodes.OBJECT_TYPE),
		Behavior:   behavior,
		Missing_ok: ifExists,
	}
}

func TestCreateEnum(t *testing.T) {
	c := New()
	err := c.DefineEnum(makeCreateEnumStmt("", "mood", []string{"happy", "sad", "neutral"}))
	if err != nil {
		t.Fatal(err)
	}

	// Resolve the enum type.
	oid, _, err := c.ResolveType(TypeName{Name: "mood", TypeMod: -1})
	if err != nil {
		t.Fatal(err)
	}
	bt := c.typeByOID[oid]
	if bt.Type != 'e' {
		t.Errorf("type kind: got %c, want 'e'", bt.Type)
	}

	// Check values.
	vals := c.EnumValues(oid)
	if len(vals) != 3 {
		t.Fatalf("values: got %d, want 3", len(vals))
	}
	if vals[0] != "happy" || vals[1] != "sad" || vals[2] != "neutral" {
		t.Errorf("values: got %v", vals)
	}
}

func TestCreateEnumDuplicate(t *testing.T) {
	c := New()
	c.DefineEnum(makeCreateEnumStmt("", "mood", []string{"a"}))
	err := c.DefineEnum(makeCreateEnumStmt("", "mood", []string{"b"}))
	assertCode(t, err, CodeDuplicateObject)
}

func TestCreateEnumDuplicateLabel(t *testing.T) {
	c := New()
	err := c.DefineEnum(makeCreateEnumStmt("", "mood", []string{"a", "b", "a"}))
	assertCode(t, err, CodeInvalidParameterValue)
}

func TestEnumAsColumnType(t *testing.T) {
	c := New()
	c.DefineEnum(makeCreateEnumStmt("", "mood", []string{"happy"}))

	err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "m", Type: TypeName{Name: "mood", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	oid, _, _ := c.ResolveType(TypeName{Name: "mood", TypeMod: -1})
	if rel.Columns[0].TypeOID != oid {
		t.Errorf("column type: got %d, want %d", rel.Columns[0].TypeOID, oid)
	}
}

func TestAlterEnumAddValue(t *testing.T) {
	c := New()
	c.DefineEnum(makeCreateEnumStmt("", "mood", []string{"happy", "sad"}))

	err := c.AlterEnumStmt(makeAlterEnumStmt("", "mood", "neutral", "", "", false))
	if err != nil {
		t.Fatal(err)
	}

	oid, _, _ := c.ResolveType(TypeName{Name: "mood", TypeMod: -1})
	vals := c.EnumValues(oid)
	if len(vals) != 3 {
		t.Fatalf("values: got %d, want 3", len(vals))
	}
	if vals[2] != "neutral" {
		t.Errorf("last value: got %q, want %q", vals[2], "neutral")
	}
}

func TestAlterEnumAddValueBefore(t *testing.T) {
	c := New()
	c.DefineEnum(makeCreateEnumStmt("", "mood", []string{"happy", "sad"}))

	err := c.AlterEnumStmt(makeAlterEnumStmt("", "mood", "neutral", "sad", "", false))
	if err != nil {
		t.Fatal(err)
	}

	oid, _, _ := c.ResolveType(TypeName{Name: "mood", TypeMod: -1})
	vals := c.EnumValues(oid)
	if len(vals) != 3 {
		t.Fatalf("values: got %d, want 3", len(vals))
	}
	if vals[1] != "neutral" {
		t.Errorf("values: got %v, want [happy neutral sad]", vals)
	}
}

func TestAlterEnumAddValueAfter(t *testing.T) {
	c := New()
	c.DefineEnum(makeCreateEnumStmt("", "mood", []string{"happy", "sad"}))

	err := c.AlterEnumStmt(makeAlterEnumStmt("", "mood", "neutral", "", "happy", false))
	if err != nil {
		t.Fatal(err)
	}

	oid, _, _ := c.ResolveType(TypeName{Name: "mood", TypeMod: -1})
	vals := c.EnumValues(oid)
	if len(vals) != 3 {
		t.Fatalf("values: got %d, want 3", len(vals))
	}
	if vals[1] != "neutral" {
		t.Errorf("values: got %v, want [happy neutral sad]", vals)
	}
}

func TestAlterEnumDuplicateLabel(t *testing.T) {
	c := New()
	c.DefineEnum(makeCreateEnumStmt("", "mood", []string{"happy"}))
	err := c.AlterEnumStmt(makeAlterEnumStmt("", "mood", "happy", "", "", false))
	assertCode(t, err, CodeDuplicateObject)
}

func TestAlterEnumDuplicateLabelIfNotExists(t *testing.T) {
	c := New()
	c.DefineEnum(makeCreateEnumStmt("", "mood", []string{"happy"}))
	err := c.AlterEnumStmt(makeAlterEnumStmt("", "mood", "happy", "", "", true))
	if err != nil {
		t.Fatal(err)
	}
}

func TestDropTypeEnum(t *testing.T) {
	c := New()
	c.DefineEnum(makeCreateEnumStmt("", "mood", []string{"happy"}))
	err := c.RemoveObjects(makeDropTypeStmt("", "mood", false, false))
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = c.ResolveType(TypeName{Name: "mood", TypeMod: -1})
	assertCode(t, err, CodeUndefinedObject)
}

func TestDropTypeEnumCascade(t *testing.T) {
	c := New()
	c.DefineEnum(makeCreateEnumStmt("", "mood", []string{"happy"}))
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "m", Type: TypeName{Name: "mood", TypeMod: -1}},
	}, nil, false), 'r')

	// Should fail without cascade.
	err := c.RemoveObjects(makeDropTypeStmt("", "mood", false, false))
	assertCode(t, err, CodeDependentObjects)

	// Should succeed with cascade and drop the dependent table.
	err = c.RemoveObjects(makeDropTypeStmt("", "mood", false, true))
	if err != nil {
		t.Fatal(err)
	}
	if c.GetRelation("", "t") != nil {
		t.Error("dependent table should be dropped after CASCADE")
	}
}

func TestDropTypeEnumCascadeDropsFunction(t *testing.T) {
	c := New()
	c.DefineEnum(makeCreateEnumStmt("", "mood", []string{"happy"}))
	c.CreateFunctionStmt(makeCreateFuncStmt("", "mood_func",
		[]TypeName{{Name: "mood", TypeMod: -1}},
		TypeName{Name: "text", TypeMod: -1},
		"sql", "SELECT 'ok'", false))

	// Should fail without cascade.
	err := c.RemoveObjects(makeDropTypeStmt("", "mood", false, false))
	assertCode(t, err, CodeDependentObjects)

	// CASCADE should drop the dependent function.
	err = c.RemoveObjects(makeDropTypeStmt("", "mood", false, true))
	if err != nil {
		t.Fatal(err)
	}
	for _, up := range c.userProcs {
		if up.Name == "mood_func" {
			t.Error("function should be gone after type CASCADE")
		}
	}
}

func TestDropTypeIfExists(t *testing.T) {
	c := New()
	err := c.RemoveObjects(makeDropTypeStmt("", "nope", true, false))
	if err != nil {
		t.Fatal(err)
	}
}

func TestEnumArrayType(t *testing.T) {
	c := New()
	c.DefineEnum(makeCreateEnumStmt("", "mood", []string{"happy"}))

	// Array type should exist.
	oid, _, err := c.ResolveType(TypeName{Name: "mood", TypeMod: -1, IsArray: true})
	if err != nil {
		t.Fatal(err)
	}
	at := c.typeByOID[oid]
	if at.Category != 'A' {
		t.Errorf("array category: got %c, want 'A'", at.Category)
	}
}
