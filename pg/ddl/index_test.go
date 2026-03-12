package ddl

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// makeIndexStmt builds an IndexStmt for DefineIndex.
func makeIndexStmt(schema, table, name string, columns []string, isUnique, ifNotExists bool) *nodes.IndexStmt {
	var params []nodes.Node
	for _, col := range columns {
		params = append(params, &nodes.IndexElem{Name: col})
	}
	stmt := &nodes.IndexStmt{
		Idxname:     name,
		Relation:    &nodes.RangeVar{Schemaname: schema, Relname: table},
		Unique:      isUnique,
		IfNotExists: ifNotExists,
	}
	if len(params) > 0 {
		stmt.IndexParams = &nodes.List{Items: params}
	}
	return stmt
}

// makeDropIndexStmt builds a DropStmt for RemoveObjects (OBJECT_INDEX).
func makeDropIndexStmt(schema, name string, ifExists bool) *nodes.DropStmt {
	var items []nodes.Node
	if schema != "" {
		items = []nodes.Node{&nodes.String{Str: schema}, &nodes.String{Str: name}}
	} else {
		items = []nodes.Node{&nodes.String{Str: name}}
	}
	return &nodes.DropStmt{
		Objects:    &nodes.List{Items: []nodes.Node{&nodes.List{Items: items}}},
		RemoveType: int(nodes.OBJECT_INDEX),
		Missing_ok: ifExists,
	}
}

func TestCreateIndex(t *testing.T) {
	c := New()

	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "integer", TypeMod: -1}},
		{Name: "b", Type: TypeName{Name: "text", TypeMod: -1}},
	}, nil, false), 'r')

	err := c.DefineIndex(makeIndexStmt("", "t", "", []string{"a", "b"}, false, false))
	if err != nil {
		t.Fatal(err)
	}

	schema := c.GetSchema("public")
	idx := schema.Indexes["t_a_b_idx"]
	if idx == nil {
		t.Fatal("index not found in schema")
	}
	if idx.IsUnique {
		t.Error("index should not be unique")
	}
	if len(idx.Columns) != 2 {
		t.Errorf("expected 2 columns, got %d", len(idx.Columns))
	}
}

func TestCreateIndexDuplicateName(t *testing.T) {
	c := New()

	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "integer", TypeMod: -1}},
	}, nil, false), 'r')

	c.DefineIndex(makeIndexStmt("", "t", "myidx", []string{"a"}, false, false))
	err := c.DefineIndex(makeIndexStmt("", "t", "myidx", []string{"a"}, false, false))
	assertErrorCode(t, err, CodeDuplicateObject)
}

func TestCreateIndexConflictsWithTable(t *testing.T) {
	c := New()

	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "integer", TypeMod: -1}},
	}, nil, false), 'r')

	// Index name same as table name.
	err := c.DefineIndex(makeIndexStmt("", "t", "t", []string{"a"}, false, false))
	assertErrorCode(t, err, CodeDuplicateTable)
}

func TestCreateIndexIfNotExists(t *testing.T) {
	c := New()

	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "integer", TypeMod: -1}},
	}, nil, false), 'r')

	c.DefineIndex(makeIndexStmt("", "t", "myidx", []string{"a"}, false, false))
	err := c.DefineIndex(makeIndexStmt("", "t", "myidx", []string{"a"}, false, true))
	if err != nil {
		t.Fatalf("IF NOT EXISTS should not error, got: %v", err)
	}
}

func TestDropIndexStandalone(t *testing.T) {
	c := New()

	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "integer", TypeMod: -1}},
	}, nil, false), 'r')
	c.DefineIndex(makeIndexStmt("", "t", "myidx", []string{"a"}, false, false))

	err := c.RemoveObjects(makeDropIndexStmt("", "myidx", false))
	if err != nil {
		t.Fatal(err)
	}

	schema := c.GetSchema("public")
	if schema.Indexes["myidx"] != nil {
		t.Error("index should be removed from schema")
	}
}

func TestDropIndexBackingConstraint(t *testing.T) {
	c := New()

	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "integer", TypeMod: -1}},
	}, []ConstraintDef{
		{Type: ConstraintPK, Columns: []string{"a"}},
	}, false), 'r')

	// Try to drop the PK's backing index directly.
	err := c.RemoveObjects(makeDropIndexStmt("", "t_pkey", false))
	assertErrorCode(t, err, CodeDependentObjects)
}

func TestDropIndexIfExists(t *testing.T) {
	c := New()

	err := c.RemoveObjects(makeDropIndexStmt("", "nosuch", true))
	if err != nil {
		t.Fatalf("IF EXISTS should not error, got: %v", err)
	}
}

func TestCreateIndexAccessMethod(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "integer", TypeMod: -1}},
	}, nil, false), 'r')

	stmt := makeIndexStmt("", "t", "myidx", []string{"a"}, false, false)
	stmt.AccessMethod = "hash"
	err := c.DefineIndex(stmt)
	if err != nil {
		t.Fatal(err)
	}

	schema := c.GetSchema("public")
	idx := schema.Indexes["myidx"]
	if idx.AccessMethod != "hash" {
		t.Errorf("access method: got %q, want %q", idx.AccessMethod, "hash")
	}
}

func TestCreateIndexDefaultAccessMethod(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "integer", TypeMod: -1}},
	}, nil, false), 'r')

	err := c.DefineIndex(makeIndexStmt("", "t", "myidx", []string{"a"}, false, false))
	if err != nil {
		t.Fatal(err)
	}

	schema := c.GetSchema("public")
	idx := schema.Indexes["myidx"]
	if idx.AccessMethod != "btree" {
		t.Errorf("access method: got %q, want %q", idx.AccessMethod, "btree")
	}
}

func TestCreateIndexWithInclude(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "integer", TypeMod: -1}},
		{Name: "b", Type: TypeName{Name: "text", TypeMod: -1}},
		{Name: "c", Type: TypeName{Name: "boolean", TypeMod: -1}},
	}, nil, false), 'r')

	stmt := makeIndexStmt("", "t", "myidx", []string{"a"}, true, false)
	// Add INCLUDE columns.
	stmt.IndexIncludingParams = &nodes.List{Items: []nodes.Node{
		&nodes.IndexElem{Name: "b"},
		&nodes.IndexElem{Name: "c"},
	}}
	err := c.DefineIndex(stmt)
	if err != nil {
		t.Fatal(err)
	}

	schema := c.GetSchema("public")
	idx := schema.Indexes["myidx"]
	if len(idx.Columns) != 3 {
		t.Fatalf("columns: got %d, want 3", len(idx.Columns))
	}
	if idx.NKeyColumns != 1 {
		t.Errorf("NKeyColumns: got %d, want 1", idx.NKeyColumns)
	}
}

func TestCreateIndexWithWhere(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "integer", TypeMod: -1}},
	}, nil, false), 'r')

	stmt := makeIndexStmt("", "t", "myidx", []string{"a"}, false, false)
	stmt.WhereClause = &nodes.A_Const{Val: &nodes.String{Str: "a > 0"}}
	err := c.DefineIndex(stmt)
	if err != nil {
		t.Fatal(err)
	}

	schema := c.GetSchema("public")
	idx := schema.Indexes["myidx"]
	if idx.WhereClause == "" {
		t.Error("WhereClause should not be empty")
	}
}

func TestCreateIndexExpressionElement(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "integer", TypeMod: -1}},
	}, nil, false), 'r')

	// Index with an expression element (Name="") and a regular column.
	stmt := &nodes.IndexStmt{
		Idxname:  "expr_idx",
		Relation: &nodes.RangeVar{Relname: "t"},
		IndexParams: &nodes.List{Items: []nodes.Node{
			&nodes.IndexElem{Name: ""},  // expression
			&nodes.IndexElem{Name: "a"}, // column
		}},
	}
	err := c.DefineIndex(stmt)
	if err != nil {
		t.Fatal(err)
	}

	schema := c.GetSchema("public")
	idx := schema.Indexes["expr_idx"]
	if len(idx.Columns) != 2 {
		t.Fatalf("columns: got %d, want 2", len(idx.Columns))
	}
	// Expression column should have attnum 0.
	if idx.Columns[0] != 0 {
		t.Errorf("expression attnum: got %d, want 0", idx.Columns[0])
	}
	if idx.Columns[1] != 1 {
		t.Errorf("column attnum: got %d, want 1", idx.Columns[1])
	}
}

func TestIndexColumnsResolve(t *testing.T) {
	c := New()

	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "integer", TypeMod: -1}},
		{Name: "b", Type: TypeName{Name: "text", TypeMod: -1}},
		{Name: "c", Type: TypeName{Name: "boolean", TypeMod: -1}},
	}, nil, false), 'r')

	c.DefineIndex(makeIndexStmt("", "t", "", []string{"b", "c"}, false, false))

	rel := c.GetRelation("", "t")
	idxs := c.IndexesOf(rel.OID)
	if len(idxs) != 1 {
		t.Fatalf("expected 1 index, got %d", len(idxs))
	}
	// b is attnum 2, c is attnum 3.
	if idxs[0].Columns[0] != 2 || idxs[0].Columns[1] != 3 {
		t.Errorf("columns: got %v, want [2 3]", idxs[0].Columns)
	}
}
