package catalog

import "testing"

func TestCreateTable(t *testing.T) {
	c := New()

	err := c.DefineRelation(makeCreateTableStmt("", "users", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "integer", TypeMod: -1}, NotNull: true},
		{Name: "name", Type: TypeName{Name: "text", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	r := c.GetRelation("", "users")
	if r == nil {
		t.Fatal("table not found after CREATE")
	}
	if r.RelKind != 'r' {
		t.Errorf("RelKind: got %c, want 'r'", r.RelKind)
	}
	if len(r.Columns) != 2 {
		t.Errorf("columns: got %d, want 2", len(r.Columns))
	}
	if r.Columns[0].TypeOID != INT4OID {
		t.Errorf("first column type: got %d, want %d", r.Columns[0].TypeOID, INT4OID)
	}
	if !r.Columns[0].NotNull {
		t.Error("first column should be NOT NULL")
	}
}

func TestCreateTableDuplicate(t *testing.T) {
	c := New()

	cols := []ColumnDef{{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}}
	c.DefineRelation(makeCreateTableStmt("", "t", cols, nil, false), 'r')

	err := c.DefineRelation(makeCreateTableStmt("", "t", cols, nil, false), 'r')
	assertErrorCode(t, err, CodeDuplicateTable)
}

func TestCreateTableIfNotExists(t *testing.T) {
	c := New()

	cols := []ColumnDef{{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}}
	c.DefineRelation(makeCreateTableStmt("", "t", cols, nil, false), 'r')

	err := c.DefineRelation(makeCreateTableStmt("", "t", cols, nil, true), 'r')
	if err != nil {
		t.Fatalf("IF NOT EXISTS should not error, got: %v", err)
	}
}

func TestCreateTableDuplicateColumn(t *testing.T) {
	c := New()

	err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	assertErrorCode(t, err, CodeDuplicateColumn)
}

func TestCreateTableUnknownType(t *testing.T) {
	c := New()

	err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "nosuch", TypeMod: -1}},
	}, nil, false), 'r')
	assertErrorCode(t, err, CodeUndefinedObject)
}

func TestCreateTableVarcharTypmod(t *testing.T) {
	c := New()

	err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "v", Type: TypeName{Name: "varchar", TypeMod: 255}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	r := c.GetRelation("", "t")
	if r.Columns[0].TypeMod != 255 {
		t.Errorf("typmod: got %d, want 255", r.Columns[0].TypeMod)
	}
}

func TestCreateTableArrayColumn(t *testing.T) {
	c := New()

	err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "ids", Type: TypeName{Name: "integer", TypeMod: -1, IsArray: true}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	r := c.GetRelation("", "t")
	if r.Columns[0].TypeOID != INT4ARRAYOID {
		t.Errorf("array column type: got %d, want %d", r.Columns[0].TypeOID, INT4ARRAYOID)
	}
}

func TestCreateTableExplicitSchema(t *testing.T) {
	c := New()

	c.CreateSchemaCommand(makeCreateSchemaStmt("myschema", false))
	err := c.DefineRelation(makeCreateTableStmt("myschema", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	// Should be found via explicit schema.
	if r := c.GetRelation("myschema", "t"); r == nil {
		t.Fatal("table not found in explicit schema")
	}
	// Should NOT be found via default search path.
	if r := c.GetRelation("", "t"); r != nil {
		t.Fatal("table should not be on default search path")
	}
}

func TestCreateTableRowType(t *testing.T) {
	c := New()

	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	// Row type should be resolvable.
	oid, _, err := c.ResolveType(TypeName{Name: "t", TypeMod: -1})
	if err != nil {
		t.Fatalf("row type not resolvable: %v", err)
	}
	typ := c.TypeByOID(oid)
	if typ == nil || typ.Type != 'c' {
		t.Fatal("expected composite type")
	}
}

func TestCreateTableConflictsWithType(t *testing.T) {
	c := New()

	// Create a table first (which registers a row type).
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	// Another table with the same name should fail.
	err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "x", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	assertErrorCode(t, err, CodeDuplicateTable)
}

func TestDropTable(t *testing.T) {
	c := New()

	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	if err := c.RemoveRelations(makeDropTableStmt("", "t", false, false)); err != nil {
		t.Fatal(err)
	}
	if r := c.GetRelation("", "t"); r != nil {
		t.Fatal("table still exists after DROP")
	}

	// Row type should also be gone.
	_, _, err := c.ResolveType(TypeName{Name: "t", TypeMod: -1})
	if err == nil {
		t.Fatal("row type should not be resolvable after DROP")
	}
}

func TestDropTableIfExists(t *testing.T) {
	c := New()

	if err := c.RemoveRelations(makeDropTableStmt("", "nosuch", true, false)); err != nil {
		t.Fatalf("IF EXISTS should not error, got: %v", err)
	}
}

func TestDropTableNonExistent(t *testing.T) {
	c := New()

	err := c.RemoveRelations(makeDropTableStmt("", "nosuch", false, false))
	assertErrorCode(t, err, CodeUndefinedTable)
}

func TestCreateTableNoColumns(t *testing.T) {
	c := New()

	err := c.DefineRelation(makeCreateTableStmt("", "t", nil, nil, false), 'r')
	if err != nil {
		t.Fatalf("DefineRelation failed: %v", err)
	}
	rel := c.GetRelation("", "t")
	if rel == nil {
		t.Fatal("expected relation t")
	}
	if got := len(rel.Columns); got != 0 {
		t.Fatalf("columns: got %d, want 0", got)
	}
}
