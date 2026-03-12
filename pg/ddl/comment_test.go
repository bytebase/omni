package ddl

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// makeCommentStmtTable builds a COMMENT ON TABLE/VIEW statement.
func makeCommentStmtTable(objtype nodes.ObjectType, schema, name string, comment *string) *nodes.CommentStmt {
	stmt := &nodes.CommentStmt{
		Objtype: objtype,
		Object:  makeQualNameList(schema, name),
	}
	if comment != nil {
		stmt.Comment = *comment
	}
	return stmt
}

// makeCommentStmtColumn builds a COMMENT ON COLUMN statement.
func makeCommentStmtColumn(schema, table, column string, comment *string) *nodes.CommentStmt {
	var items []nodes.Node
	if schema != "" {
		items = append(items, &nodes.String{Str: schema})
	}
	items = append(items, &nodes.String{Str: table})
	items = append(items, &nodes.String{Str: column})
	stmt := &nodes.CommentStmt{
		Objtype: nodes.OBJECT_COLUMN,
		Object:  &nodes.List{Items: items},
	}
	if comment != nil {
		stmt.Comment = *comment
	}
	return stmt
}

// makeCommentStmtSchema builds a COMMENT ON SCHEMA statement.
func makeCommentStmtSchema(name string, comment *string) *nodes.CommentStmt {
	stmt := &nodes.CommentStmt{
		Objtype: nodes.OBJECT_SCHEMA,
		Object:  &nodes.String{Str: name},
	}
	if comment != nil {
		stmt.Comment = *comment
	}
	return stmt
}

// makeCommentStmtIndex builds a COMMENT ON INDEX statement.
func makeCommentStmtIndex(schema, name string, comment *string) *nodes.CommentStmt {
	return makeCommentStmtTable(nodes.OBJECT_INDEX, schema, name, comment)
}

// makeCommentStmtSequence builds a COMMENT ON SEQUENCE statement.
func makeCommentStmtSequence(schema, name string, comment *string) *nodes.CommentStmt {
	return makeCommentStmtTable(nodes.OBJECT_SEQUENCE, schema, name, comment)
}

// makeCommentStmtType builds a COMMENT ON TYPE statement.
func makeCommentStmtType(schema, name string, comment *string) *nodes.CommentStmt {
	return makeCommentStmtTable(nodes.OBJECT_TYPE, schema, name, comment)
}

// makeCommentStmtFunction builds a COMMENT ON FUNCTION statement.
func makeCommentStmtFunction(schema, name string, argTypes []TypeName, comment *string) *nodes.CommentStmt {
	var nameItems []nodes.Node
	if schema != "" {
		nameItems = append(nameItems, &nodes.String{Str: schema})
	}
	nameItems = append(nameItems, &nodes.String{Str: name})
	owa := &nodes.ObjectWithArgs{Objname: &nodes.List{Items: nameItems}}
	if len(argTypes) > 0 {
		var args []nodes.Node
		for _, at := range argTypes {
			args = append(args, makeTypeNameNode(at))
		}
		owa.Objargs = &nodes.List{Items: args}
	} else {
		owa.ArgsUnspecified = true
	}
	stmt := &nodes.CommentStmt{
		Objtype: nodes.OBJECT_FUNCTION,
		Object:  owa,
	}
	if comment != nil {
		stmt.Comment = *comment
	}
	return stmt
}

// makeCommentStmtTrigger builds a COMMENT ON TRIGGER statement.
func makeCommentStmtTrigger(schema, table, name string, comment *string) *nodes.CommentStmt {
	var items []nodes.Node
	if schema != "" {
		items = append(items, &nodes.String{Str: schema})
	}
	items = append(items, &nodes.String{Str: table})
	items = append(items, &nodes.String{Str: name})
	stmt := &nodes.CommentStmt{
		Objtype: nodes.OBJECT_TRIGGER,
		Object:  &nodes.List{Items: items},
	}
	if comment != nil {
		stmt.Comment = *comment
	}
	return stmt
}

func makeQualNameList(schema, name string) *nodes.List {
	var items []nodes.Node
	if schema != "" {
		items = append(items, &nodes.String{Str: schema})
	}
	items = append(items, &nodes.String{Str: name})
	return &nodes.List{Items: items}
}

func TestCommentOnTable(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{{Name: "x", Type: TypeName{Name: "int4", TypeMod: -1}}}, nil, false), 'r')

	comment := "my table"
	err := c.CommentObject(makeCommentStmtTable(nodes.OBJECT_TABLE, "", "t", &comment))
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	got, ok := c.GetComment('r', rel.OID, 0)
	if !ok || got != "my table" {
		t.Errorf("comment: got %q, ok=%v", got, ok)
	}
}

func TestCommentOnColumn(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "b", Type: TypeName{Name: "text", TypeMod: -1}},
	}, nil, false), 'r')

	comment := "second col"
	err := c.CommentObject(makeCommentStmtColumn("", "t", "b", &comment))
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	got, ok := c.GetComment('r', rel.OID, 2)
	if !ok || got != "second col" {
		t.Errorf("comment: got %q, ok=%v", got, ok)
	}
}

func TestCommentOnSchema(t *testing.T) {
	c := New()
	comment := "the public schema"
	err := c.CommentObject(makeCommentStmtSchema("public", &comment))
	if err != nil {
		t.Fatal(err)
	}

	got, ok := c.GetComment('n', PublicNamespace, 0)
	if !ok || got != "the public schema" {
		t.Errorf("comment: got %q, ok=%v", got, ok)
	}
}

func TestCommentOnIndex(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{{Name: "x", Type: TypeName{Name: "int4", TypeMod: -1}}}, nil, false), 'r')
	c.DefineIndex(makeIndexStmt("", "t", "t_x_idx", []string{"x"}, false, false))

	comment := "my index"
	err := c.CommentObject(makeCommentStmtIndex("", "t_x_idx", &comment))
	if err != nil {
		t.Fatal(err)
	}

	schema := c.GetSchema("public")
	idx := schema.Indexes["t_x_idx"]
	got, ok := c.GetComment('i', idx.OID, 0)
	if !ok || got != "my index" {
		t.Errorf("comment: got %q, ok=%v", got, ok)
	}
}

func TestCommentOnSequence(t *testing.T) {
	c := New()
	c.DefineSequence(makeCreateSeqStmt("", "s", false))

	comment := "my seq"
	err := c.CommentObject(makeCommentStmtSequence("", "s", &comment))
	if err != nil {
		t.Fatal(err)
	}

	seq, _ := c.findSequence("", "s")
	got, ok := c.GetComment('s', seq.OID, 0)
	if !ok || got != "my seq" {
		t.Errorf("comment: got %q, ok=%v", got, ok)
	}
}

func TestCommentOnFunction(t *testing.T) {
	c := New()
	c.CreateFunctionStmt(makeCreateFuncStmt("", "f",
		[]TypeName{{Name: "int4", TypeMod: -1}},
		TypeName{Name: "int4", TypeMod: -1},
		"sql", "SELECT $1", false))

	comment := "my func"
	err := c.CommentObject(makeCommentStmtFunction("", "f",
		[]TypeName{{Name: "int4", TypeMod: -1}}, &comment))
	if err != nil {
		t.Fatal(err)
	}
}

func TestCommentOnType(t *testing.T) {
	c := New()
	c.DefineEnum(makeCreateEnumStmt("", "mood", []string{"happy"}))

	comment := "feelings"
	err := c.CommentObject(makeCommentStmtType("", "mood", &comment))
	if err != nil {
		t.Fatal(err)
	}

	oid, _, _ := c.ResolveType(TypeName{Name: "mood", TypeMod: -1})
	got, ok := c.GetComment('t', oid, 0)
	if !ok || got != "feelings" {
		t.Errorf("comment: got %q, ok=%v", got, ok)
	}
}

func TestCommentRemove(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{{Name: "x", Type: TypeName{Name: "int4", TypeMod: -1}}}, nil, false), 'r')

	comment := "temporary"
	c.CommentObject(makeCommentStmtTable(nodes.OBJECT_TABLE, "", "t", &comment))

	// Remove comment by passing nil (empty string in CommentStmt).
	err := c.CommentObject(makeCommentStmtTable(nodes.OBJECT_TABLE, "", "t", nil))
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	_, ok := c.GetComment('r', rel.OID, 0)
	if ok {
		t.Error("comment should be removed")
	}
}

func TestCommentNonexistentTarget(t *testing.T) {
	c := New()
	comment := "nope"
	err := c.CommentObject(makeCommentStmtTable(nodes.OBJECT_TABLE, "", "nope", &comment))
	assertCode(t, err, CodeUndefinedTable)
}

func TestCommentOnView(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{{Name: "x", Type: TypeName{Name: "int4", TypeMod: -1}}}, nil, false), 'r')
	c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			From:       []*FromItem{{Kind: FromTable, Table: "t"}},
			TargetList: []*ResTarget{{Val: &Expr{Kind: ExprStar}}},
		}, nil))

	comment := "my view"
	err := c.CommentObject(makeCommentStmtTable(nodes.OBJECT_VIEW, "", "v", &comment))
	if err != nil {
		t.Fatal(err)
	}
}

func TestCommentCleanedOnDropTable(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "x", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	// Add comments on table and column.
	tblComment := "the table"
	c.CommentObject(makeCommentStmtTable(nodes.OBJECT_TABLE, "", "t", &tblComment))
	colComment := "the column"
	c.CommentObject(makeCommentStmtColumn("", "t", "x", &colComment))

	c.RemoveRelations(makeDropTableStmt("", "t", false, false))

	if len(c.comments) != 0 {
		t.Errorf("comments should be empty after drop, got %d", len(c.comments))
	}
}

func TestCommentCleanedOnDropSchema(t *testing.T) {
	c := New()
	c.CreateSchemaCommand(makeCreateSchemaStmt("s", false))

	comment := "my schema"
	c.CommentObject(makeCommentStmtSchema("s", &comment))

	c.RemoveSchemas(makeDropSchemaStmt("s", false, false))

	if len(c.comments) != 0 {
		t.Errorf("comments should be empty after drop, got %d", len(c.comments))
	}
}

func TestCommentOnTrigger(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{{Name: "x", Type: TypeName{Name: "int4", TypeMod: -1}}}, nil, false), 'r')
	c.CreateFunctionStmt(makeCreateFuncStmt("", "trig_fn",
		nil, TypeName{Name: "trigger", TypeMod: -1},
		"plpgsql", "BEGIN RETURN NEW; END;", false))
	c.CreateTriggerStmt(makeCreateTrigStmt("", "t", "my_trig", TriggerBefore, TriggerEventInsert, false, "", "trig_fn"))

	comment := "my trigger"
	err := c.CommentObject(makeCommentStmtTrigger("", "t", "my_trig", &comment))
	if err != nil {
		t.Fatal(err)
	}
}
