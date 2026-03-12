package ddl

import (
	"math"
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// ---------------------------------------------------------------------------
// Test helpers for sequence statements
// ---------------------------------------------------------------------------

func makeCreateSeqStmt(schema, name string, ifNotExists bool, opts ...*nodes.DefElem) *nodes.CreateSeqStmt {
	stmt := &nodes.CreateSeqStmt{
		Sequence:    &nodes.RangeVar{Schemaname: schema, Relname: name},
		IfNotExists: ifNotExists,
	}
	if len(opts) > 0 {
		items := make([]nodes.Node, len(opts))
		for i, o := range opts {
			items[i] = o
		}
		stmt.Options = &nodes.List{Items: items}
	}
	return stmt
}

func seqOptInt(name string, val int64) *nodes.DefElem {
	return &nodes.DefElem{Defname: name, Arg: &nodes.Integer{Ival: val}}
}

func seqOptType(typeName string) *nodes.DefElem {
	return &nodes.DefElem{
		Defname: "as",
		Arg:     &nodes.TypeName{Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: typeName}}}},
	}
}

func seqOptCycle() *nodes.DefElem {
	return &nodes.DefElem{Defname: "cycle", Arg: &nodes.Boolean{Boolval: true}}
}

func seqOptOwnedBy(ref string) *nodes.DefElem {
	parts := strings.SplitN(ref, ".", 2)
	items := make([]nodes.Node, len(parts))
	for i, p := range parts {
		items[i] = &nodes.String{Str: p}
	}
	return &nodes.DefElem{Defname: "owned_by", Arg: &nodes.List{Items: items}}
}

func makeAlterSeqStmt(schema, name string, opts ...*nodes.DefElem) *nodes.AlterSeqStmt {
	stmt := &nodes.AlterSeqStmt{
		Sequence: &nodes.RangeVar{Schemaname: schema, Relname: name},
	}
	if len(opts) > 0 {
		items := make([]nodes.Node, len(opts))
		for i, o := range opts {
			items[i] = o
		}
		stmt.Options = &nodes.List{Items: items}
	}
	return stmt
}

func makeDropSeqStmt(schema, name string, ifExists, cascade bool) *nodes.DropStmt {
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
		RemoveType: int(nodes.OBJECT_SEQUENCE),
		Behavior:   behavior,
		Missing_ok: ifExists,
	}
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestCreateSequenceBasic(t *testing.T) {
	c := New()
	err := c.DefineSequence(makeCreateSeqStmt("", "my_seq", false))
	if err != nil {
		t.Fatal(err)
	}

	seq, err := c.findSequence("", "my_seq")
	if err != nil {
		t.Fatal(err)
	}
	if seq.TypeOID != INT8OID {
		t.Errorf("type: got %d, want %d", seq.TypeOID, INT8OID)
	}
	if seq.Increment != 1 {
		t.Errorf("increment: got %d, want 1", seq.Increment)
	}
	if seq.MinValue != 1 {
		t.Errorf("min: got %d, want 1", seq.MinValue)
	}
	if seq.MaxValue != math.MaxInt64 {
		t.Errorf("max: got %d, want MaxInt64", seq.MaxValue)
	}
	if seq.Start != 1 {
		t.Errorf("start: got %d, want 1", seq.Start)
	}
	if seq.CacheValue != 1 {
		t.Errorf("cache: got %d, want 1", seq.CacheValue)
	}
}

func TestCreateSequenceSmallint(t *testing.T) {
	c := New()
	err := c.DefineSequence(makeCreateSeqStmt("", "s", false, seqOptType("smallint")))
	if err != nil {
		t.Fatal(err)
	}
	seq, _ := c.findSequence("", "s")
	if seq.TypeOID != INT2OID {
		t.Errorf("type: got %d, want %d", seq.TypeOID, INT2OID)
	}
	if seq.MaxValue != math.MaxInt16 {
		t.Errorf("max: got %d, want %d", seq.MaxValue, math.MaxInt16)
	}
}

func TestCreateSequenceInteger(t *testing.T) {
	c := New()
	err := c.DefineSequence(makeCreateSeqStmt("", "s", false, seqOptType("integer")))
	if err != nil {
		t.Fatal(err)
	}
	seq, _ := c.findSequence("", "s")
	if seq.TypeOID != INT4OID {
		t.Errorf("type: got %d, want %d", seq.TypeOID, INT4OID)
	}
	if seq.MaxValue != math.MaxInt32 {
		t.Errorf("max: got %d, want %d", seq.MaxValue, math.MaxInt32)
	}
}

func TestCreateSequenceDescending(t *testing.T) {
	c := New()
	err := c.DefineSequence(makeCreateSeqStmt("", "s", false, seqOptInt("increment", -1)))
	if err != nil {
		t.Fatal(err)
	}
	seq, _ := c.findSequence("", "s")
	if seq.MinValue != math.MinInt64 {
		t.Errorf("min: got %d, want MinInt64", seq.MinValue)
	}
	if seq.MaxValue != -1 {
		t.Errorf("max: got %d, want -1", seq.MaxValue)
	}
	if seq.Start != -1 {
		t.Errorf("start: got %d, want -1", seq.Start)
	}
}

func TestCreateSequenceDuplicate(t *testing.T) {
	c := New()
	c.DefineSequence(makeCreateSeqStmt("", "s", false))
	err := c.DefineSequence(makeCreateSeqStmt("", "s", false))
	assertCode(t, err, CodeDuplicateObject)
}

func TestCreateSequenceIfNotExists(t *testing.T) {
	c := New()
	c.DefineSequence(makeCreateSeqStmt("", "s", false))
	err := c.DefineSequence(makeCreateSeqStmt("", "s", true))
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateSequenceZeroIncrement(t *testing.T) {
	c := New()
	err := c.DefineSequence(makeCreateSeqStmt("", "s", false, seqOptInt("increment", 0)))
	assertCode(t, err, CodeInvalidParameterValue)
}

func TestDropSequence(t *testing.T) {
	c := New()
	c.DefineSequence(makeCreateSeqStmt("", "s", false))
	err := c.RemoveObjects(makeDropSeqStmt("", "s", false, false))
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.findSequence("", "s")
	assertCode(t, err, CodeUndefinedObject)
}

func TestDropSequenceIfExists(t *testing.T) {
	c := New()
	err := c.RemoveObjects(makeDropSeqStmt("", "nope", true, false))
	if err != nil {
		t.Fatal(err)
	}
}

func TestDropSequenceNonexistent(t *testing.T) {
	c := New()
	err := c.RemoveObjects(makeDropSeqStmt("", "nope", false, false))
	assertCode(t, err, CodeUndefinedObject)
}

func TestAlterSequence(t *testing.T) {
	c := New()
	c.DefineSequence(makeCreateSeqStmt("", "s", false))

	err := c.AlterSequenceStmt(makeAlterSeqStmt("", "s",
		seqOptInt("increment", 5),
		seqOptCycle(),
	))
	if err != nil {
		t.Fatal(err)
	}

	seq, _ := c.findSequence("", "s")
	if seq.Increment != 5 {
		t.Errorf("increment: got %d, want 5", seq.Increment)
	}
	if !seq.Cycle {
		t.Error("cycle should be true")
	}
}

func TestSerialColumn(t *testing.T) {
	c := New()
	err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", IsSerial: 4},
		{Name: "name", Type: TypeName{Name: "text", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	if rel == nil {
		t.Fatal("table not found")
	}

	// Column should be int4, NOT NULL, with default.
	idCol := rel.Columns[0]
	if idCol.TypeOID != INT4OID {
		t.Errorf("type: got %d, want %d", idCol.TypeOID, INT4OID)
	}
	if !idCol.NotNull {
		t.Error("serial column should be NOT NULL")
	}
	if !idCol.HasDefault {
		t.Error("serial column should have default")
	}

	// Sequence should exist.
	seq, err := c.findSequence("", "t_id_seq")
	if err != nil {
		t.Fatal(err)
	}
	if seq.TypeOID != INT4OID {
		t.Errorf("seq type: got %d, want %d", seq.TypeOID, INT4OID)
	}
	if seq.OwnerRelOID != rel.OID {
		t.Error("sequence should be owned by the table")
	}
}

func TestBigserialColumn(t *testing.T) {
	c := New()
	err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", IsSerial: 8},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	if rel.Columns[0].TypeOID != INT8OID {
		t.Errorf("type: got %d, want %d", rel.Columns[0].TypeOID, INT8OID)
	}

	seq, err := c.findSequence("", "t_id_seq")
	if err != nil {
		t.Fatal(err)
	}
	if seq.TypeOID != INT8OID {
		t.Errorf("seq type: got %d, want %d", seq.TypeOID, INT8OID)
	}
}

func TestSmallserialColumn(t *testing.T) {
	c := New()
	err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", IsSerial: 2},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	if rel.Columns[0].TypeOID != INT2OID {
		t.Errorf("type: got %d, want %d", rel.Columns[0].TypeOID, INT2OID)
	}
}

func TestDropTableCascadeDropsOwnedSequence(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", IsSerial: 4},
	}, nil, false), 'r')
	err := c.RemoveRelations(makeDropTableStmt("", "t", false, true))
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.findSequence("", "t_id_seq")
	assertCode(t, err, CodeUndefinedObject)
}

func TestSequenceConflictsWithRelation(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	err := c.DefineSequence(makeCreateSeqStmt("", "t", false))
	assertCode(t, err, CodeDuplicateTable)
}

func TestSequenceOwnedBy(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')
	err := c.DefineSequence(makeCreateSeqStmt("", "s", false, seqOptOwnedBy("t.id")))
	if err != nil {
		t.Fatal(err)
	}

	seq, _ := c.findSequence("", "s")
	rel := c.GetRelation("", "t")
	if seq.OwnerRelOID != rel.OID {
		t.Error("sequence should be owned by the table")
	}
}

func assertCode(t *testing.T, err error, code string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error with code %s, got nil", code)
	}
	pgErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("expected *Error, got %T: %v", err, err)
	}
	if pgErr.Code != code {
		t.Errorf("code: got %s, want %s (msg: %s)", pgErr.Code, code, pgErr.Message)
	}
}
