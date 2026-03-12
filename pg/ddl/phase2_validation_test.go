package ddl

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// =============================================================================
// Phase 2: DDL Validation Hardening Tests
// =============================================================================

// TestPhase2_MergeAttributesTypmodMismatch tests that inheriting with a type
// modifier conflict is rejected.
func TestPhase2_MergeAttributesTypmodMismatch(t *testing.T) {
	c := New()
	// Create parent with varchar(10).
	parent := makeCreateTableStmt("", "p", []ColumnDef{
		{Name: "name", Type: TypeName{Name: "varchar", TypeMod: 14}}, // 10+4=14 (VARHDRSZ)
	}, nil, false)
	if err := c.DefineRelation(parent, 'r'); err != nil {
		t.Fatal(err)
	}

	// Create child with same column but different typmod.
	child := makeCreateTableStmt("", "ch", []ColumnDef{
		{Name: "name", Type: TypeName{Name: "varchar", TypeMod: 24}}, // 20+4=24
	}, nil, false)
	child.InhRelations = &nodes.List{Items: []nodes.Node{
		&nodes.RangeVar{Relname: "p"},
	}}
	err := c.DefineRelation(child, 'r')
	assertErrorCode(t, err, CodeDatatypeMismatch)
}

// TestPhase2_MergeAttributesTypmodMatch tests that matching typmods work.
func TestPhase2_MergeAttributesTypmodMatch(t *testing.T) {
	c := New()
	parent := makeCreateTableStmt("", "p", []ColumnDef{
		{Name: "name", Type: TypeName{Name: "varchar", TypeMod: 14}},
	}, nil, false)
	if err := c.DefineRelation(parent, 'r'); err != nil {
		t.Fatal(err)
	}

	child := makeCreateTableStmt("", "ch", []ColumnDef{
		{Name: "name", Type: TypeName{Name: "varchar", TypeMod: 14}},
	}, nil, false)
	child.InhRelations = &nodes.List{Items: []nodes.Node{
		&nodes.RangeVar{Relname: "p"},
	}}
	err := c.DefineRelation(child, 'r')
	if err != nil {
		t.Fatalf("expected no error for matching typmod, got: %v", err)
	}
}

// TestPhase2_MergeAttributesCollationMismatch tests collation conflict detection.
func TestPhase2_MergeAttributesCollationMismatch(t *testing.T) {
	c := New()
	// Use parser to properly set CollClause on the columns.
	stmts := parseStmts(t, `CREATE TABLE p (name text COLLATE "C")`)
	if err := c.ProcessUtility(stmts[0]); err != nil {
		t.Fatal(err)
	}

	stmts = parseStmts(t, `CREATE TABLE ch (name text COLLATE "POSIX") INHERITS (p)`)
	err := c.ProcessUtility(stmts[0])
	assertErrorCode(t, err, CodeDatatypeMismatch)
}

// TestPhase2_IndexMaxKeys tests the INDEX_MAX_KEYS limit.
func TestPhase2_IndexMaxKeys(t *testing.T) {
	c := New()
	// Create table with 33 columns.
	var cols []ColumnDef
	var params []nodes.Node
	for i := 0; i < 33; i++ {
		name := string(rune('a'+i/26)) + string(rune('a'+i%26))
		cols = append(cols, ColumnDef{Name: name, Type: TypeName{Name: "int4", TypeMod: -1}})
		params = append(params, &nodes.IndexElem{Name: name})
	}
	stmt := makeCreateTableStmt("", "t", cols, nil, false)
	if err := c.DefineRelation(stmt, 'r'); err != nil {
		t.Fatal(err)
	}

	idxStmt := &nodes.IndexStmt{
		Relation:    &nodes.RangeVar{Relname: "t"},
		IndexParams: &nodes.List{Items: params},
	}
	err := c.DefineIndex(idxStmt)
	assertErrorCode(t, err, CodeProgramLimitExceeded)
}

// TestPhase2_IndexMaxKeysOK tests that 32 columns is within the limit.
func TestPhase2_IndexMaxKeysOK(t *testing.T) {
	c := New()
	var cols []ColumnDef
	var params []nodes.Node
	for i := 0; i < 32; i++ {
		name := string(rune('a'+i/26)) + string(rune('a'+i%26))
		cols = append(cols, ColumnDef{Name: name, Type: TypeName{Name: "int4", TypeMod: -1}})
		params = append(params, &nodes.IndexElem{Name: name})
	}
	stmt := makeCreateTableStmt("", "t", cols, nil, false)
	if err := c.DefineRelation(stmt, 'r'); err != nil {
		t.Fatal(err)
	}

	idxStmt := &nodes.IndexStmt{
		Idxname:     "t_idx",
		Relation:    &nodes.RangeVar{Relname: "t"},
		IndexParams: &nodes.List{Items: params},
	}
	err := c.DefineIndex(idxStmt)
	if err != nil {
		t.Fatalf("expected no error for 32 columns, got: %v", err)
	}
}

// TestPhase2_NullsNotDistinct tests that the NullsNotDistinct field is stored.
func TestPhase2_NullsNotDistinct(t *testing.T) {
	c := New()
	stmt := makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false)
	if err := c.DefineRelation(stmt, 'r'); err != nil {
		t.Fatal(err)
	}

	idxStmt := &nodes.IndexStmt{
		Idxname:            "t_idx",
		Relation:           &nodes.RangeVar{Relname: "t"},
		IndexParams:        &nodes.List{Items: []nodes.Node{&nodes.IndexElem{Name: "id"}}},
		Unique:             true,
		Nulls_not_distinct: true,
	}
	if err := c.DefineIndex(idxStmt); err != nil {
		t.Fatal(err)
	}

	schema := c.schemaByName["public"]
	idx := schema.Indexes["t_idx"]
	if idx == nil {
		t.Fatal("index not found")
	}
	if !idx.NullsNotDistinct {
		t.Error("expected NullsNotDistinct to be true")
	}
}

// TestPhase2_DropNotNullIdentityColumn tests that dropping NOT NULL on identity
// columns is rejected.
func TestPhase2_DropNotNullIdentityColumn(t *testing.T) {
	c := New()
	stmts := parseStmts(t, "CREATE TABLE t (id int GENERATED ALWAYS AS IDENTITY)")
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatal(err)
		}
	}

	stmts = parseStmts(t, "ALTER TABLE t ALTER COLUMN id DROP NOT NULL")
	err := c.ProcessUtility(stmts[0])
	assertErrorCode(t, err, CodeInvalidObjectDefinition)
}

// TestPhase2_DropColumnPartitionKey tests that dropping a partition key column
// is rejected.
func TestPhase2_DropColumnPartitionKey(t *testing.T) {
	c := New()
	stmt := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "t"},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}),
			makeColumnDefNode(ColumnDef{Name: "val", Type: TypeName{Name: "text"}}),
		}},
		Partspec: &nodes.PartitionSpec{
			Strategy:   "list",
			PartParams: &nodes.List{Items: []nodes.Node{&nodes.PartitionElem{Name: "id"}}},
		},
	}
	if err := c.DefineRelation(stmt, 'r'); err != nil {
		t.Fatal(err)
	}

	alterStmts := parseStmts(t, "ALTER TABLE t DROP COLUMN id")
	err := c.ProcessUtility(alterStmts[0])
	assertErrorCode(t, err, CodeInvalidObjectDefinition)
}

// TestPhase2_TruncateOnlyPartitioned tests TRUNCATE ONLY on partitioned table.
func TestPhase2_TruncateOnlyPartitioned(t *testing.T) {
	c := New()
	stmt := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "t"},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}),
		}},
		Partspec: &nodes.PartitionSpec{
			Strategy:   "list",
			PartParams: &nodes.List{Items: []nodes.Node{&nodes.PartitionElem{Name: "id"}}},
		},
	}
	if err := c.DefineRelation(stmt, 'r'); err != nil {
		t.Fatal(err)
	}

	truncStmt := &nodes.TruncateStmt{
		Relations: &nodes.List{Items: []nodes.Node{
			&nodes.RangeVar{Relname: "t", Inh: false},
		}},
	}
	err := c.ExecuteTruncate(truncStmt)
	assertErrorCode(t, err, CodeWrongObjectType)
}

// TestPhase2_TruncatePartitionedWithInh tests TRUNCATE on partitioned table
// WITH inheritance (default) succeeds.
func TestPhase2_TruncatePartitionedWithInh(t *testing.T) {
	c := New()
	stmt := &nodes.CreateStmt{
		Relation: &nodes.RangeVar{Relname: "t"},
		TableElts: &nodes.List{Items: []nodes.Node{
			makeColumnDefNode(ColumnDef{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}),
		}},
		Partspec: &nodes.PartitionSpec{
			Strategy:   "list",
			PartParams: &nodes.List{Items: []nodes.Node{&nodes.PartitionElem{Name: "id"}}},
		},
	}
	if err := c.DefineRelation(stmt, 'r'); err != nil {
		t.Fatal(err)
	}

	truncStmt := &nodes.TruncateStmt{
		Relations: &nodes.List{Items: []nodes.Node{
			&nodes.RangeVar{Relname: "t", Inh: true},
		}},
	}
	err := c.ExecuteTruncate(truncStmt)
	if err != nil {
		t.Fatalf("TRUNCATE with inheritance should succeed: %v", err)
	}
}

// TestPhase2_TruncateRestrictFK tests TRUNCATE RESTRICT with FK references.
func TestPhase2_TruncateRestrictFK(t *testing.T) {
	c := New()
	// Create referenced table.
	stmts := parseStmts(t, "CREATE TABLE parent (id int PRIMARY KEY)")
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatal(err)
		}
	}

	// Create referencing table.
	stmts = parseStmts(t, "CREATE TABLE child (id int, pid int REFERENCES parent(id))")
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatal(err)
		}
	}

	// TRUNCATE parent with RESTRICT should fail.
	truncStmt := &nodes.TruncateStmt{
		Relations: &nodes.List{Items: []nodes.Node{
			&nodes.RangeVar{Relname: "parent", Inh: true},
		}},
		Behavior: nodes.DROP_RESTRICT,
	}
	err := c.ExecuteTruncate(truncStmt)
	assertErrorCode(t, err, CodeDependentObjects)
}

// TestPhase2_TruncateCascadeFK tests TRUNCATE CASCADE with FK references succeeds.
func TestPhase2_TruncateCascadeFK(t *testing.T) {
	c := New()
	stmts := parseStmts(t, "CREATE TABLE parent (id int PRIMARY KEY)")
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatal(err)
		}
	}
	stmts = parseStmts(t, "CREATE TABLE child (id int, pid int REFERENCES parent(id))")
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatal(err)
		}
	}

	truncStmt := &nodes.TruncateStmt{
		Relations: &nodes.List{Items: []nodes.Node{
			&nodes.RangeVar{Relname: "parent", Inh: true},
		}},
		Behavior: nodes.DROP_CASCADE,
	}
	err := c.ExecuteTruncate(truncStmt)
	if err != nil {
		t.Fatalf("TRUNCATE CASCADE should succeed: %v", err)
	}
}
