package ddl

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// --- ALTER TABLE AST helpers ---

func makeAlterTableStmt(schema, name string, cmds ...*nodes.AlterTableCmd) *nodes.AlterTableStmt {
	items := make([]nodes.Node, len(cmds))
	for i, cmd := range cmds {
		items[i] = cmd
	}
	return &nodes.AlterTableStmt{
		Relation: &nodes.RangeVar{Schemaname: schema, Relname: name},
		Cmds:     &nodes.List{Items: items},
	}
}

func makeATAddColumn(cd ColumnDef) *nodes.AlterTableCmd {
	return &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_AddColumn),
		Def:     makeColumnDefNode(cd),
	}
}

func makeATDropColumn(name string, cascade, ifExists bool) *nodes.AlterTableCmd {
	behavior := int(nodes.DROP_RESTRICT)
	if cascade {
		behavior = int(nodes.DROP_CASCADE)
	}
	return &nodes.AlterTableCmd{
		Subtype:    int(nodes.AT_DropColumn),
		Name:       name,
		Behavior:   behavior,
		Missing_ok: ifExists,
	}
}

func makeATAddConstraint(cd ConstraintDef) *nodes.AlterTableCmd {
	return &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_AddConstraint),
		Def:     makeConstraintNode(cd),
	}
}

func makeATDropConstraint(name string, cascade, ifExists bool) *nodes.AlterTableCmd {
	behavior := int(nodes.DROP_RESTRICT)
	if cascade {
		behavior = int(nodes.DROP_CASCADE)
	}
	return &nodes.AlterTableCmd{
		Subtype:    int(nodes.AT_DropConstraint),
		Name:       name,
		Behavior:   behavior,
		Missing_ok: ifExists,
	}
}

func makeATSetNotNull(col string) *nodes.AlterTableCmd {
	return &nodes.AlterTableCmd{Subtype: int(nodes.AT_SetNotNull), Name: col}
}

func makeATDropNotNull(col string) *nodes.AlterTableCmd {
	return &nodes.AlterTableCmd{Subtype: int(nodes.AT_DropNotNull), Name: col}
}

func makeATSetDefault(col, defaultExpr string) *nodes.AlterTableCmd {
	return &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_ColumnDefault),
		Name:    col,
		Def:     &nodes.Constraint{Contype: nodes.CONSTR_DEFAULT, CookedExpr: defaultExpr},
	}
}

func makeATDropDefault(col string) *nodes.AlterTableCmd {
	return &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_ColumnDefault),
		Name:    col,
		// Def is nil -> means drop default
	}
}

func makeATAlterColumnType(col string, tn TypeName) *nodes.AlterTableCmd {
	return &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_AlterColumnType),
		Name:    col,
		Def:     &nodes.ColumnDef{TypeName: makeTypeNameNode(tn)},
	}
}

func makeRenameColumnStmt(schema, table, oldName, newName string) *nodes.RenameStmt {
	return &nodes.RenameStmt{
		RenameType: nodes.OBJECT_COLUMN,
		Relation:   &nodes.RangeVar{Schemaname: schema, Relname: table},
		Subname:    oldName,
		Newname:    newName,
	}
}

func makeRenameTableStmt(schema, oldName, newName string) *nodes.RenameStmt {
	return &nodes.RenameStmt{
		RenameType: nodes.OBJECT_TABLE,
		Relation:   &nodes.RangeVar{Schemaname: schema, Relname: oldName},
		Newname:    newName,
	}
}

// --- Tests ---

func TestAlterTableAddColumn(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}}, nil, false), 'r')

	err := c.AlterTableStmt(makeAlterTableStmt("", "t",
		makeATAddColumn(ColumnDef{Name: "name", Type: TypeName{Name: "text", TypeMod: -1}}),
	))
	if err != nil {
		t.Fatal(err)
	}

	r := c.GetRelation("", "t")
	if len(r.Columns) != 2 {
		t.Fatalf("columns: got %d, want 2", len(r.Columns))
	}
	if r.Columns[1].Name != "name" {
		t.Errorf("column name: got %q, want %q", r.Columns[1].Name, "name")
	}
	if r.Columns[1].TypeOID != TEXTOID {
		t.Errorf("column type: got %d, want %d", r.Columns[1].TypeOID, TEXTOID)
	}
	if r.Columns[1].AttNum != 2 {
		t.Errorf("attnum: got %d, want 2", r.Columns[1].AttNum)
	}
}

func TestAlterTableAddColumnDuplicate(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}}, nil, false), 'r')

	err := c.AlterTableStmt(makeAlterTableStmt("", "t",
		makeATAddColumn(ColumnDef{Name: "id", Type: TypeName{Name: "text", TypeMod: -1}}),
	))
	assertErrorCode(t, err, CodeDuplicateColumn)
}

func TestAlterTableDropColumn(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "b", Type: TypeName{Name: "text", TypeMod: -1}},
		{Name: "c", Type: TypeName{Name: "int8", TypeMod: -1}},
	}, nil, false), 'r')

	err := c.AlterTableStmt(makeAlterTableStmt("", "t",
		makeATDropColumn("b", false, false),
	))
	if err != nil {
		t.Fatal(err)
	}

	r := c.GetRelation("", "t")
	if len(r.Columns) != 2 {
		t.Fatalf("columns: got %d, want 2", len(r.Columns))
	}
	if r.Columns[0].Name != "a" || r.Columns[1].Name != "c" {
		t.Errorf("remaining columns: %q, %q", r.Columns[0].Name, r.Columns[1].Name)
	}
	// Attnums should be renumbered.
	if r.Columns[0].AttNum != 1 || r.Columns[1].AttNum != 2 {
		t.Errorf("attnums: %d, %d — want 1, 2", r.Columns[0].AttNum, r.Columns[1].AttNum)
	}
	// colByName should be updated.
	if _, ok := r.colByName["b"]; ok {
		t.Error("dropped column still in colByName")
	}
}

func TestAlterTableDropColumnIfExists(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}}, nil, false), 'r')

	err := c.AlterTableStmt(makeAlterTableStmt("", "t",
		makeATDropColumn("nosuch", false, true),
	))
	if err != nil {
		t.Fatalf("IF EXISTS should not error, got: %v", err)
	}
}

func TestAlterTableDropColumnNonexistent(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}}, nil, false), 'r')

	err := c.AlterTableStmt(makeAlterTableStmt("", "t",
		makeATDropColumn("nosuch", false, false),
	))
	assertErrorCode(t, err, CodeUndefinedColumn)
}

func TestAlterTableAddConstraint(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	err := c.AlterTableStmt(makeAlterTableStmt("", "t",
		makeATAddConstraint(ConstraintDef{Type: ConstraintPK, Columns: []string{"id"}}),
	))
	if err != nil {
		t.Fatal(err)
	}

	r := c.GetRelation("", "t")
	cons := c.ConstraintsOf(r.OID)
	if len(cons) != 1 {
		t.Fatalf("constraints: got %d, want 1", len(cons))
	}
	if cons[0].Type != ConstraintPK {
		t.Errorf("constraint type: got %c, want %c", cons[0].Type, ConstraintPK)
	}
}

func TestAlterTableDropConstraint(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, []ConstraintDef{{Type: ConstraintPK, Columns: []string{"id"}}}, false), 'r')

	r := c.GetRelation("", "t")
	if len(c.ConstraintsOf(r.OID)) != 1 {
		t.Fatal("expected 1 constraint before DROP")
	}

	err := c.AlterTableStmt(makeAlterTableStmt("", "t",
		makeATDropConstraint("t_pkey", false, false),
	))
	if err != nil {
		t.Fatal(err)
	}

	if len(c.ConstraintsOf(r.OID)) != 0 {
		t.Error("constraint still exists after DROP")
	}
	// Backing index should also be removed.
	if len(c.IndexesOf(r.OID)) != 0 {
		t.Error("backing index still exists after DROP CONSTRAINT")
	}
}

func TestAlterTableDropConstraintIfExists(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}}, nil, false), 'r')

	err := c.AlterTableStmt(makeAlterTableStmt("", "t",
		makeATDropConstraint("nosuch", false, true),
	))
	if err != nil {
		t.Fatalf("IF EXISTS should not error, got: %v", err)
	}
}

func TestAlterTableRenameColumn(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "old_name", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	err := c.ExecRenameStmt(makeRenameColumnStmt("", "t", "old_name", "new_name"))
	if err != nil {
		t.Fatal(err)
	}

	r := c.GetRelation("", "t")
	if r.Columns[0].Name != "new_name" {
		t.Errorf("column name: got %q, want %q", r.Columns[0].Name, "new_name")
	}
	if _, ok := r.colByName["old_name"]; ok {
		t.Error("old name still in colByName")
	}
	if _, ok := r.colByName["new_name"]; !ok {
		t.Error("new name not in colByName")
	}
}

func TestAlterTableRenameColumnNonexistent(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}}, nil, false), 'r')

	err := c.ExecRenameStmt(makeRenameColumnStmt("", "t", "nosuch", "x"))
	assertErrorCode(t, err, CodeUndefinedColumn)
}

func TestAlterTableRenameTable(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "old_t", []ColumnDef{{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}}, nil, false), 'r')

	err := c.ExecRenameStmt(makeRenameTableStmt("", "old_t", "new_t"))
	if err != nil {
		t.Fatal(err)
	}

	if r := c.GetRelation("", "old_t"); r != nil {
		t.Error("old name still resolves")
	}
	r := c.GetRelation("", "new_t")
	if r == nil {
		t.Fatal("new name does not resolve")
	}

	// Row type name should be updated.
	oid, _, err := c.ResolveType(TypeName{Name: "new_t", TypeMod: -1})
	if err != nil {
		t.Fatalf("row type not resolvable: %v", err)
	}
	typ := c.TypeByOID(oid)
	if typ.TypeName != "new_t" {
		t.Errorf("row type name: got %q, want %q", typ.TypeName, "new_t")
	}

	// Old row type should not resolve.
	_, _, err = c.ResolveType(TypeName{Name: "old_t", TypeMod: -1})
	if err == nil {
		t.Error("old row type name should not resolve")
	}
}

func TestAlterTableRenameTableDuplicate(t *testing.T) {
	c := New()
	cols := []ColumnDef{{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}}
	c.DefineRelation(makeCreateTableStmt("", "t1", cols, nil, false), 'r')
	c.DefineRelation(makeCreateTableStmt("", "t2", cols, nil, false), 'r')

	err := c.ExecRenameStmt(makeRenameTableStmt("", "t1", "t2"))
	assertErrorCode(t, err, CodeDuplicateTable)
}

func TestAlterTableSetNotNull(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}}, nil, false), 'r')

	err := c.AlterTableStmt(makeAlterTableStmt("", "t",
		makeATSetNotNull("id"),
	))
	if err != nil {
		t.Fatal(err)
	}

	r := c.GetRelation("", "t")
	if !r.Columns[0].NotNull {
		t.Error("column should be NOT NULL")
	}
}

func TestAlterTableDropNotNull(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}, NotNull: true}}, nil, false), 'r')

	err := c.AlterTableStmt(makeAlterTableStmt("", "t",
		makeATDropNotNull("id"),
	))
	if err != nil {
		t.Fatal(err)
	}

	r := c.GetRelation("", "t")
	if r.Columns[0].NotNull {
		t.Error("column should not be NOT NULL")
	}
}

func TestAlterTableDropNotNullOnPK(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}}, []ConstraintDef{{Type: ConstraintPK, Columns: []string{"id"}}}, false), 'r')

	err := c.AlterTableStmt(makeAlterTableStmt("", "t",
		makeATDropNotNull("id"),
	))
	assertErrorCode(t, err, CodeDependentObjects)
}

func TestAlterTableSetDefault(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}}, nil, false), 'r')

	err := c.AlterTableStmt(makeAlterTableStmt("", "t",
		makeATSetDefault("id", "42"),
	))
	if err != nil {
		t.Fatal(err)
	}

	r := c.GetRelation("", "t")
	if !r.Columns[0].HasDefault {
		t.Error("column should have default")
	}
	if r.Columns[0].Default != "42" {
		t.Errorf("default: got %q, want %q", r.Columns[0].Default, "42")
	}
}

func TestAlterTableDropDefault(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}, Default: "1"},
	}, nil, false), 'r')

	err := c.AlterTableStmt(makeAlterTableStmt("", "t",
		makeATDropDefault("id"),
	))
	if err != nil {
		t.Fatal(err)
	}

	r := c.GetRelation("", "t")
	if r.Columns[0].HasDefault {
		t.Error("column should not have default")
	}
}

func TestAlterTableAlterColumnType(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}}, nil, false), 'r')

	err := c.AlterTableStmt(makeAlterTableStmt("", "t",
		makeATAlterColumnType("id", TypeName{Name: "int8", TypeMod: -1}),
	))
	if err != nil {
		t.Fatal(err)
	}

	r := c.GetRelation("", "t")
	if r.Columns[0].TypeOID != INT8OID {
		t.Errorf("column type: got %d, want %d", r.Columns[0].TypeOID, INT8OID)
	}
}

func TestAlterTableAlterColumnTypeIncompatible(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{{Name: "ts", Type: TypeName{Name: "timestamp", TypeMod: -1}}}, nil, false), 'r')

	err := c.AlterTableStmt(makeAlterTableStmt("", "t",
		makeATAlterColumnType("ts", TypeName{Name: "int4", TypeMod: -1}),
	))
	assertErrorCode(t, err, CodeDatatypeMismatch)
}

func TestAlterTableOnView(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}}, nil, false), 'r')
	c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			TargetList: []*ResTarget{{Val: &Expr{Kind: ExprStar}}},
			From:       []*FromItem{{Kind: FromTable, Table: "t"}},
		}, nil))

	err := c.AlterTableStmt(makeAlterTableStmt("", "v",
		makeATSetNotNull("id"),
	))
	assertErrorCode(t, err, CodeWrongObjectType)
}

// --- ALTER TABLE pass ordering tests ---

func TestAlterTablePassOrdering(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "b", Type: TypeName{Name: "text", TypeMod: -1}},
	}, nil, false), 'r')

	// Issue DROP COLUMN b and ADD COLUMN c in statement order.
	// Pass ordering should execute DROP first (pass 0), then ADD (pass 2).
	err := c.AlterTableStmt(makeAlterTableStmt("", "t",
		makeATAddColumn(ColumnDef{Name: "c", Type: TypeName{Name: "int8", TypeMod: -1}}),
		makeATDropColumn("b", false, false),
	))
	if err != nil {
		t.Fatal(err)
	}

	r := c.GetRelation("", "t")
	if len(r.Columns) != 2 {
		t.Fatalf("columns: got %d, want 2", len(r.Columns))
	}
	if r.Columns[0].Name != "a" || r.Columns[1].Name != "c" {
		t.Errorf("columns: got %q, %q — want a, c", r.Columns[0].Name, r.Columns[1].Name)
	}
}

// --- Rename tests for expanded ExecRenameStmt ---

func makeRenameViewStmt(schema, oldName, newName string) *nodes.RenameStmt {
	return &nodes.RenameStmt{
		RenameType: nodes.OBJECT_VIEW,
		Relation:   &nodes.RangeVar{Schemaname: schema, Relname: oldName},
		Newname:    newName,
	}
}

func makeRenameIndexStmt(schema, oldName, newName string) *nodes.RenameStmt {
	return &nodes.RenameStmt{
		RenameType: nodes.OBJECT_INDEX,
		Relation:   &nodes.RangeVar{Schemaname: schema, Relname: oldName},
		Newname:    newName,
	}
}

func makeRenameSequenceStmt(schema, oldName, newName string) *nodes.RenameStmt {
	return &nodes.RenameStmt{
		RenameType: nodes.OBJECT_SEQUENCE,
		Relation:   &nodes.RangeVar{Schemaname: schema, Relname: oldName},
		Newname:    newName,
	}
}

func makeRenameSchemaStmt(oldName, newName string) *nodes.RenameStmt {
	return &nodes.RenameStmt{
		RenameType: nodes.OBJECT_SCHEMA,
		Object:     &nodes.String{Str: oldName},
		Newname:    newName,
	}
}

func makeRenameTypeStmt(schema, oldName, newName string) *nodes.RenameStmt {
	var items []nodes.Node
	if schema != "" {
		items = append(items, &nodes.String{Str: schema})
	}
	items = append(items, &nodes.String{Str: oldName})
	return &nodes.RenameStmt{
		RenameType: nodes.OBJECT_TYPE,
		Object:     &nodes.List{Items: items},
		Newname:    newName,
	}
}

func makeRenameTriggerStmt(schema, table, oldName, newName string) *nodes.RenameStmt {
	return &nodes.RenameStmt{
		RenameType: nodes.OBJECT_TRIGGER,
		Relation:   &nodes.RangeVar{Schemaname: schema, Relname: table},
		Subname:    oldName,
		Newname:    newName,
	}
}

func makeRenameConstraintStmt(schema, table, oldName, newName string) *nodes.RenameStmt {
	return &nodes.RenameStmt{
		RenameType: nodes.OBJECT_TABCONSTRAINT,
		Relation:   &nodes.RangeVar{Schemaname: schema, Relname: table},
		Subname:    oldName,
		Newname:    newName,
	}
}

func TestRenameView(t *testing.T) {
	c := newTestCatalogWithTable()
	c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			TargetList: []*ResTarget{{Val: &Expr{Kind: ExprStar}}},
			From:       []*FromItem{{Kind: FromTable, Table: "t"}},
		}, nil))

	err := c.ExecRenameStmt(makeRenameViewStmt("", "v", "v2"))
	if err != nil {
		t.Fatal(err)
	}
	if c.GetRelation("", "v") != nil {
		t.Error("old name still resolves")
	}
	r := c.GetRelation("", "v2")
	if r == nil {
		t.Fatal("new name does not resolve")
	}
	if r.RelKind != 'v' {
		t.Errorf("relkind: got %c, want v", r.RelKind)
	}
}

func TestRenameIndex(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	c.DefineIndex(makeIndexStmt("", "t", "myidx", []string{"a"}, false, false))

	err := c.ExecRenameStmt(makeRenameIndexStmt("", "myidx", "new_idx"))
	if err != nil {
		t.Fatal(err)
	}

	schema := c.GetSchema("public")
	if schema.Indexes["myidx"] != nil {
		t.Error("old index name still in schema")
	}
	if schema.Indexes["new_idx"] == nil {
		t.Error("new index name not in schema")
	}
}

func TestRenameSequence(t *testing.T) {
	c := New()
	c.DefineSequence(makeCreateSeqStmt("", "myseq", false))

	err := c.ExecRenameStmt(makeRenameSequenceStmt("", "myseq", "new_seq"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = c.findSequence("", "myseq")
	if err == nil {
		t.Error("old sequence name still resolves")
	}
	_, err = c.findSequence("", "new_seq")
	if err != nil {
		t.Fatalf("new sequence name not found: %v", err)
	}
}

func TestRenameSchema(t *testing.T) {
	c := New()
	c.CreateSchemaCommand(makeCreateSchemaStmt("old_s", false))

	err := c.ExecRenameStmt(makeRenameSchemaStmt("old_s", "new_s"))
	if err != nil {
		t.Fatal(err)
	}

	if c.GetSchema("old_s") != nil {
		t.Error("old schema name still resolves")
	}
	if c.GetSchema("new_s") == nil {
		t.Error("new schema name not found")
	}
}

func TestRenameSchemaDuplicate(t *testing.T) {
	c := New()
	c.CreateSchemaCommand(makeCreateSchemaStmt("s1", false))
	c.CreateSchemaCommand(makeCreateSchemaStmt("s2", false))

	err := c.ExecRenameStmt(makeRenameSchemaStmt("s1", "s2"))
	assertErrorCode(t, err, CodeDuplicateSchema)
}

func TestRenameSchemaBuiltinFails(t *testing.T) {
	c := New()
	err := c.ExecRenameStmt(makeRenameSchemaStmt("pg_catalog", "new_pg"))
	if err == nil {
		t.Fatal("expected error renaming pg_catalog")
	}
}

func TestRenameType(t *testing.T) {
	c := New()
	c.DefineEnum(makeCreateEnumStmt("", "mood", []string{"happy", "sad"}))

	err := c.ExecRenameStmt(makeRenameTypeStmt("", "mood", "feeling"))
	if err != nil {
		t.Fatal(err)
	}

	// Old name should not resolve.
	_, _, err = c.ResolveType(TypeName{Name: "mood", TypeMod: -1})
	if err == nil {
		t.Error("old type name still resolves")
	}
	// New name should resolve.
	oid, _, err := c.ResolveType(TypeName{Name: "feeling", TypeMod: -1})
	if err != nil {
		t.Fatalf("new type name not found: %v", err)
	}
	bt := c.TypeByOID(oid)
	if bt.TypeName != "feeling" {
		t.Errorf("type name: got %q, want %q", bt.TypeName, "feeling")
	}
	// Array type should also be updated.
	_, _, err = c.ResolveType(TypeName{Name: "_feeling", TypeMod: -1})
	if err != nil {
		t.Fatalf("array type not found: %v", err)
	}
}

func TestRenameTrigger(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}}, nil, false), 'r')
	c.CreateFunctionStmt(makeCreateFuncStmt("", "trig_fn",
		nil, TypeName{Name: "trigger", TypeMod: -1},
		"plpgsql", "BEGIN RETURN NEW; END;", false))
	c.CreateTriggerStmt(makeCreateTrigStmt("", "t", "my_trig", TriggerBefore, TriggerEventInsert, false, "", "trig_fn"))

	err := c.ExecRenameStmt(makeRenameTriggerStmt("", "t", "my_trig", "new_trig"))
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	trigs := c.TriggersOf(rel.OID)
	if len(trigs) != 1 {
		t.Fatalf("triggers: got %d, want 1", len(trigs))
	}
	if trigs[0].Name != "new_trig" {
		t.Errorf("trigger name: got %q, want %q", trigs[0].Name, "new_trig")
	}
}

func TestRenameTriggerNonexistent(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}}, nil, false), 'r')

	err := c.ExecRenameStmt(makeRenameTriggerStmt("", "t", "nosuch", "x"))
	assertErrorCode(t, err, CodeUndefinedObject)
}

func TestRenameConstraint(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, []ConstraintDef{{Type: ConstraintPK, Columns: []string{"id"}}}, false), 'r')

	err := c.ExecRenameStmt(makeRenameConstraintStmt("", "t", "t_pkey", "t_pk_new"))
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	cons := c.ConstraintsOf(rel.OID)
	if len(cons) != 1 {
		t.Fatalf("constraints: got %d, want 1", len(cons))
	}
	if cons[0].Name != "t_pk_new" {
		t.Errorf("constraint name: got %q, want %q", cons[0].Name, "t_pk_new")
	}
}

func TestRenameConstraintDuplicate(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "name", Type: TypeName{Name: "text", TypeMod: -1}},
	}, []ConstraintDef{
		{Type: ConstraintPK, Columns: []string{"id"}},
		{Type: ConstraintUnique, Columns: []string{"name"}},
	}, false), 'r')

	rel := c.GetRelation("", "t")
	cons := c.ConstraintsOf(rel.OID)
	if len(cons) < 2 {
		t.Fatalf("expected >= 2 constraints, got %d", len(cons))
	}

	// Try to rename PK to the unique constraint's name.
	err := c.ExecRenameStmt(makeRenameConstraintStmt("", "t", cons[0].Name, cons[1].Name))
	assertErrorCode(t, err, CodeDuplicateObject)
}

// --- Phase 1 new DDL statement tests ---

// --- AlterFunction AST helpers ---

func makeAlterFunctionStmt(schema, name string, argTypes []TypeName, actions ...*nodes.DefElem) *nodes.AlterFunctionStmt {
	var argItems []nodes.Node
	for _, at := range argTypes {
		argItems = append(argItems, makeTypeNameNode(at))
	}
	var objargs *nodes.List
	if len(argItems) > 0 {
		objargs = &nodes.List{Items: argItems}
	}
	items := make([]nodes.Node, len(actions))
	for i, a := range actions {
		items[i] = a
	}
	return &nodes.AlterFunctionStmt{
		Func: &nodes.ObjectWithArgs{
			Objname: makeQualNameList(schema, name),
			Objargs: objargs,
		},
		Actions: &nodes.List{Items: items},
	}
}

// --- AlterFunction tests ---

func TestAlterFunctionVolatility(t *testing.T) {
	c := New()
	c.CreateFunctionStmt(makeCreateFuncStmt("", "myfunc",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		TypeName{Name: "int4", TypeMod: -1},
		"sql", "SELECT $1", false))

	err := c.AlterFunction(makeAlterFunctionStmt("", "myfunc",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		&nodes.DefElem{Defname: "volatility", Arg: &nodes.String{Str: "immutable"}},
	))
	if err != nil {
		t.Fatal(err)
	}

	procs := c.procByName["myfunc"]
	if len(procs) == 0 {
		t.Fatal("function not found in procByName")
	}
	bp := procs[0]
	if bp.Volatile != 'i' {
		t.Errorf("volatile: got %c, want 'i'", bp.Volatile)
	}

	// Also check UserProc.
	up := c.userProcs[bp.OID]
	if up == nil {
		t.Fatal("user proc not found")
	}
	if up.Volatile != 'i' {
		t.Errorf("user proc volatile: got %c, want 'i'", up.Volatile)
	}
}

func TestAlterFunctionStrict(t *testing.T) {
	c := New()
	c.CreateFunctionStmt(makeCreateFuncStmt("", "myfunc",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		TypeName{Name: "int4", TypeMod: -1},
		"sql", "SELECT $1", false))

	err := c.AlterFunction(makeAlterFunctionStmt("", "myfunc",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		&nodes.DefElem{Defname: "strict", Arg: nil}, // nil Arg means true
	))
	if err != nil {
		t.Fatal(err)
	}

	procs := c.procByName["myfunc"]
	if len(procs) == 0 {
		t.Fatal("function not found")
	}
	if !procs[0].IsStrict {
		t.Error("expected IsStrict=true")
	}
	up := c.userProcs[procs[0].OID]
	if up == nil {
		t.Fatal("user proc not found")
	}
	if !up.IsStrict {
		t.Error("user proc: expected IsStrict=true")
	}
}

func TestAlterFunctionSecurityDefiner(t *testing.T) {
	c := New()
	c.CreateFunctionStmt(makeCreateFuncStmt("", "myfunc",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		TypeName{Name: "int4", TypeMod: -1},
		"sql", "SELECT $1", false))

	err := c.AlterFunction(makeAlterFunctionStmt("", "myfunc",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		&nodes.DefElem{Defname: "security", Arg: &nodes.String{Str: "definer"}},
	))
	if err != nil {
		t.Fatal(err)
	}

	procs := c.procByName["myfunc"]
	if len(procs) == 0 {
		t.Fatal("function not found")
	}
	if !procs[0].SecDef {
		t.Error("expected SecDef=true")
	}
	up := c.userProcs[procs[0].OID]
	if up == nil {
		t.Fatal("user proc not found")
	}
	if !up.SecDef {
		t.Error("user proc: expected SecDef=true")
	}
}

func TestAlterFunctionParallelSafe(t *testing.T) {
	c := New()
	c.CreateFunctionStmt(makeCreateFuncStmt("", "myfunc",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		TypeName{Name: "int4", TypeMod: -1},
		"sql", "SELECT $1", false))

	err := c.AlterFunction(makeAlterFunctionStmt("", "myfunc",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		&nodes.DefElem{Defname: "parallel", Arg: &nodes.String{Str: "safe"}},
	))
	if err != nil {
		t.Fatal(err)
	}

	procs := c.procByName["myfunc"]
	if len(procs) == 0 {
		t.Fatal("function not found")
	}
	if procs[0].Parallel != 's' {
		t.Errorf("parallel: got %c, want 's'", procs[0].Parallel)
	}
	up := c.userProcs[procs[0].OID]
	if up == nil {
		t.Fatal("user proc not found")
	}
	if up.Parallel != 's' {
		t.Errorf("user proc parallel: got %c, want 's'", up.Parallel)
	}
}

func TestAlterFunctionMultipleActions(t *testing.T) {
	c := New()
	c.CreateFunctionStmt(makeCreateFuncStmt("", "myfunc",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		TypeName{Name: "int4", TypeMod: -1},
		"sql", "SELECT $1", false))

	err := c.AlterFunction(makeAlterFunctionStmt("", "myfunc",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		&nodes.DefElem{Defname: "volatility", Arg: &nodes.String{Str: "stable"}},
		&nodes.DefElem{Defname: "strict", Arg: nil},
		&nodes.DefElem{Defname: "parallel", Arg: &nodes.String{Str: "restricted"}},
	))
	if err != nil {
		t.Fatal(err)
	}

	procs := c.procByName["myfunc"]
	if len(procs) == 0 {
		t.Fatal("function not found")
	}
	bp := procs[0]
	if bp.Volatile != 's' {
		t.Errorf("volatile: got %c, want 's'", bp.Volatile)
	}
	if !bp.IsStrict {
		t.Error("expected IsStrict=true")
	}
	if bp.Parallel != 'r' {
		t.Errorf("parallel: got %c, want 'r'", bp.Parallel)
	}
}

func TestAlterFunctionNonexistent(t *testing.T) {
	c := New()
	err := c.AlterFunction(makeAlterFunctionStmt("", "nosuch",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		&nodes.DefElem{Defname: "volatility", Arg: &nodes.String{Str: "immutable"}},
	))
	assertErrorCode(t, err, CodeUndefinedFunction)
}

// --- ExecAlterObjectSchemaStmt (SET SCHEMA) helpers ---

func makeAlterObjectSchemaStmt(objType nodes.ObjectType, schema, name, newSchema string) *nodes.AlterObjectSchemaStmt {
	stmt := &nodes.AlterObjectSchemaStmt{
		ObjectType: objType,
		Newschema:  newSchema,
	}
	if objType == nodes.OBJECT_TABLE || objType == nodes.OBJECT_VIEW || objType == nodes.OBJECT_MATVIEW || objType == nodes.OBJECT_SEQUENCE {
		stmt.Relation = &nodes.RangeVar{Schemaname: schema, Relname: name}
	} else {
		var items []nodes.Node
		if schema != "" {
			items = append(items, &nodes.String{Str: schema})
		}
		items = append(items, &nodes.String{Str: name})
		stmt.Object = &nodes.List{Items: items}
	}
	return stmt
}

// --- ExecAlterObjectSchemaStmt (SET SCHEMA) tests ---

func TestSetSchemaTable(t *testing.T) {
	c := New()
	c.CreateSchemaCommand(makeCreateSchemaStmt("s2", false))
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	err := c.ExecAlterObjectSchemaStmt(makeAlterObjectSchemaStmt(
		nodes.OBJECT_TABLE, "", "t", "s2"))
	if err != nil {
		t.Fatal(err)
	}

	// Old schema should not have the table.
	if r := c.GetRelation("public", "t"); r != nil {
		t.Error("table still in old schema")
	}
	// New schema should have the table.
	r := c.GetRelation("s2", "t")
	if r == nil {
		t.Fatal("table not found in new schema")
	}
	if r.Schema.Name != "s2" {
		t.Errorf("schema: got %q, want %q", r.Schema.Name, "s2")
	}

	// Row type should resolve in new namespace.
	oid, _, err := c.ResolveType(TypeName{Schema: "s2", Name: "t", TypeMod: -1})
	if err != nil {
		t.Fatalf("row type not resolvable in new schema: %v", err)
	}
	typ := c.TypeByOID(oid)
	if typ == nil {
		t.Fatal("row type not found by OID")
	}
}

func TestSetSchemaTableDuplicate(t *testing.T) {
	c := New()
	c.CreateSchemaCommand(makeCreateSchemaStmt("s2", false))
	// Create tables with the same name in different schemas.
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	c.DefineRelation(makeCreateTableStmt("s2", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	err := c.ExecAlterObjectSchemaStmt(makeAlterObjectSchemaStmt(
		nodes.OBJECT_TABLE, "", "t", "s2"))
	assertErrorCode(t, err, CodeDuplicateTable)
}

func TestSetSchemaEnum(t *testing.T) {
	c := New()
	c.CreateSchemaCommand(makeCreateSchemaStmt("s2", false))
	c.DefineEnum(makeCreateEnumStmt("", "mood", []string{"happy", "sad"}))

	err := c.ExecAlterObjectSchemaStmt(makeAlterObjectSchemaStmt(
		nodes.OBJECT_TYPE, "", "mood", "s2"))
	if err != nil {
		t.Fatal(err)
	}

	// Should resolve in new schema.
	oid, _, err := c.ResolveType(TypeName{Schema: "s2", Name: "mood", TypeMod: -1})
	if err != nil {
		t.Fatalf("type not found in new schema: %v", err)
	}
	bt := c.TypeByOID(oid)
	if bt == nil || bt.TypeName != "mood" {
		t.Error("type metadata mismatch after SET SCHEMA")
	}

	// Array type should also be moved.
	_, _, err = c.ResolveType(TypeName{Schema: "s2", Name: "_mood", TypeMod: -1})
	if err != nil {
		t.Fatalf("array type not found in new schema: %v", err)
	}
}

func TestSetSchemaTargetNotFound(t *testing.T) {
	c := New()
	// Target schema does not exist.
	err := c.ExecAlterObjectSchemaStmt(makeAlterObjectSchemaStmt(
		nodes.OBJECT_TABLE, "", "t", "nosuch"))
	assertErrorCode(t, err, CodeUndefinedSchema)
}

// --- ExecuteTruncate helpers ---

func makeTruncateStmt(tables ...*nodes.RangeVar) *nodes.TruncateStmt {
	items := make([]nodes.Node, len(tables))
	for i, rv := range tables {
		items[i] = rv
	}
	return &nodes.TruncateStmt{
		Relations: &nodes.List{Items: items},
	}
}

// --- ExecuteTruncate tests ---

func TestTruncateTable(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	err := c.ExecuteTruncate(makeTruncateStmt(
		&nodes.RangeVar{Relname: "t"},
	))
	if err != nil {
		t.Fatal(err)
	}
}

func TestTruncateView(t *testing.T) {
	c := newTestCatalogWithTable()
	c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			TargetList: []*ResTarget{{Val: &Expr{Kind: ExprStar}}},
			From:       []*FromItem{{Kind: FromTable, Table: "t"}},
		}, nil))

	err := c.ExecuteTruncate(makeTruncateStmt(
		&nodes.RangeVar{Relname: "v"},
	))
	assertErrorCode(t, err, CodeWrongObjectType)
}

func TestTruncateNonexistent(t *testing.T) {
	c := New()
	err := c.ExecuteTruncate(makeTruncateStmt(
		&nodes.RangeVar{Relname: "nosuch"},
	))
	assertErrorCode(t, err, CodeUndefinedTable)
}

func TestTruncateMultipleTables(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t1", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	c.DefineRelation(makeCreateTableStmt("", "t2", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	err := c.ExecuteTruncate(makeTruncateStmt(
		&nodes.RangeVar{Relname: "t1"},
		&nodes.RangeVar{Relname: "t2"},
	))
	if err != nil {
		t.Fatal(err)
	}
}

// --- Rename function tests ---

func makeRenameFunctionStmt(schema, oldName, newName string) *nodes.RenameStmt {
	var items []nodes.Node
	if schema != "" {
		items = append(items, &nodes.String{Str: schema})
	}
	items = append(items, &nodes.String{Str: oldName})
	return &nodes.RenameStmt{
		RenameType: nodes.OBJECT_FUNCTION,
		Object:     &nodes.List{Items: items},
		Newname:    newName,
	}
}

func TestRenameFunction(t *testing.T) {
	c := New()
	c.CreateFunctionStmt(makeCreateFuncStmt("", "oldfn",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		TypeName{Name: "int4", TypeMod: -1},
		"sql", "SELECT $1", false))

	err := c.ExecRenameStmt(makeRenameFunctionStmt("", "oldfn", "newfn"))
	if err != nil {
		t.Fatal(err)
	}

	// Old name should not resolve.
	oldProcs := c.LookupProcByName("oldfn")
	for _, p := range oldProcs {
		if c.userProcs[p.OID] != nil {
			t.Error("old function name still has user procs")
		}
	}

	// New name should resolve.
	newProcs := c.LookupProcByName("newfn")
	found := false
	for _, p := range newProcs {
		if c.userProcs[p.OID] != nil {
			found = true
		}
	}
	if !found {
		t.Error("renamed function not found under new name")
	}
}

func TestRenameFunctionNonexistent(t *testing.T) {
	c := New()
	err := c.ExecRenameStmt(makeRenameFunctionStmt("", "nosuch", "newfn"))
	assertErrorCode(t, err, CodeUndefinedFunction)
}

// --- Rename policy tests ---

func makeRenamePolicyStmt(schema, table, oldName, newName string) *nodes.RenameStmt {
	return &nodes.RenameStmt{
		RenameType: nodes.OBJECT_POLICY,
		Relation:   &nodes.RangeVar{Schemaname: schema, Relname: table},
		Subname:    oldName,
		Newname:    newName,
	}
}

func TestRenamePolicy(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	c.CreatePolicy(&nodes.CreatePolicyStmt{
		PolicyName: "mypol",
		Table:      &nodes.RangeVar{Relname: "t"},
		CmdName:    "all",
	})

	err := c.ExecRenameStmt(makeRenamePolicyStmt("", "t", "mypol", "newpol"))
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	policies := c.policiesByRel[rel.OID]
	if len(policies) != 1 {
		t.Fatalf("policies: got %d, want 1", len(policies))
	}
	if policies[0].Name != "newpol" {
		t.Errorf("policy name: got %q, want %q", policies[0].Name, "newpol")
	}
}

func TestRenamePolicyNonexistent(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	err := c.ExecRenameStmt(makeRenamePolicyStmt("", "t", "nosuch", "newpol"))
	assertErrorCode(t, err, CodeUndefinedObject)
}

func TestRenamePolicyDuplicate(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	c.CreatePolicy(&nodes.CreatePolicyStmt{
		PolicyName: "pol1",
		Table:      &nodes.RangeVar{Relname: "t"},
		CmdName:    "all",
	})
	c.CreatePolicy(&nodes.CreatePolicyStmt{
		PolicyName: "pol2",
		Table:      &nodes.RangeVar{Relname: "t"},
		CmdName:    "select",
	})

	err := c.ExecRenameStmt(makeRenamePolicyStmt("", "t", "pol1", "pol2"))
	assertErrorCode(t, err, CodeDuplicateObject)
}

// --- Rename attribute (composite type) tests ---

func makeRenameAttributeStmt(schema, typeName, oldAttr, newAttr string) *nodes.RenameStmt {
	return &nodes.RenameStmt{
		RenameType: nodes.OBJECT_ATTRIBUTE,
		Relation:   &nodes.RangeVar{Schemaname: schema, Relname: typeName},
		Subname:    oldAttr,
		Newname:    newAttr,
	}
}

func TestRenameAttribute(t *testing.T) {
	c := New()
	c.DefineCompositeType(&nodes.CompositeTypeStmt{
		Typevar: &nodes.RangeVar{Relname: "mytype"},
		Coldeflist: &nodes.List{Items: []nodes.Node{
			&nodes.ColumnDef{
				Colname:  "x",
				TypeName: makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
			},
			&nodes.ColumnDef{
				Colname:  "y",
				TypeName: makeTypeNameNode(TypeName{Name: "text", TypeMod: -1}),
			},
		}},
	})

	err := c.ExecRenameStmt(makeRenameAttributeStmt("", "mytype", "x", "z"))
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "mytype")
	if rel == nil {
		t.Fatal("composite type relation not found")
	}
	if _, ok := rel.colByName["z"]; !ok {
		t.Error("renamed attribute 'z' not found")
	}
	if _, ok := rel.colByName["x"]; ok {
		t.Error("old attribute 'x' still exists")
	}
	if rel.Columns[0].Name != "z" {
		t.Errorf("column name: got %q, want %q", rel.Columns[0].Name, "z")
	}
}

func TestRenameAttributeOnNonComposite(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	err := c.ExecRenameStmt(makeRenameAttributeStmt("", "t", "id", "new_id"))
	assertErrorCode(t, err, CodeWrongObjectType)
}

func TestRenameAttributeNonexistent(t *testing.T) {
	c := New()
	c.DefineCompositeType(&nodes.CompositeTypeStmt{
		Typevar: &nodes.RangeVar{Relname: "mytype"},
		Coldeflist: &nodes.List{Items: []nodes.Node{
			&nodes.ColumnDef{
				Colname:  "x",
				TypeName: makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
			},
		}},
	})

	err := c.ExecRenameStmt(makeRenameAttributeStmt("", "mytype", "nosuch", "z"))
	assertErrorCode(t, err, CodeUndefinedColumn)
}
