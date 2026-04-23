package catalog

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

func TestPartitionChildAddPrimaryKeyReusesPartitionConstraint(t *testing.T) {
	c := New()

	stmts := parseStmts(t, `
CREATE TABLE partition_parent (
    id integer NOT NULL,
    bucket integer NOT NULL,
    PRIMARY KEY (id, bucket)
) PARTITION BY RANGE (bucket);
CREATE TABLE partition_child
    PARTITION OF partition_parent FOR VALUES FROM (0) TO (10);
ALTER TABLE ONLY partition_child
    ADD CONSTRAINT partition_child_pkey PRIMARY KEY (id, bucket);
`)
	for _, stmt := range stmts {
		if err := c.ProcessUtility(stmt); err != nil {
			t.Fatalf("ProcessUtility: %v", err)
		}
	}

	_, child, err := c.findRelation("", "partition_child")
	if err != nil {
		t.Fatal(err)
	}

	var pkCount int
	var pk *Constraint
	for _, con := range c.consByRel[child.OID] {
		if con.Type == ConstraintPK {
			pkCount++
			pk = con
		}
	}
	if pkCount != 1 {
		t.Fatalf("partition child PK count: got %d, want 1", pkCount)
	}
	if pk.Name != "partition_child_pkey" {
		t.Fatalf("partition child PK name: got %q, want %q", pk.Name, "partition_child_pkey")
	}
	if pk.ConParentID == 0 {
		t.Fatal("partition child PK should retain parent constraint link")
	}
}

func TestAlterTableAddPrimaryKeyUsingExistingIndex(t *testing.T) {
	c := New()

	stmts := parseStmts(t, `
CREATE TABLE attach_index_table (
    id integer NOT NULL,
    bucket integer NOT NULL
);
CREATE UNIQUE INDEX attach_index_table_idx ON attach_index_table (id, bucket);
ALTER TABLE attach_index_table
    ADD CONSTRAINT attach_index_table_pkey PRIMARY KEY USING INDEX attach_index_table_idx;
`)
	for _, stmt := range stmts {
		if err := c.ProcessUtility(stmt); err != nil {
			t.Fatalf("ProcessUtility: %v", err)
		}
	}

	_, rel, err := c.findRelation("", "attach_index_table")
	if err != nil {
		t.Fatal(err)
	}

	pk := requireSingleConstraintOfType(t, c, rel, ConstraintPK)
	requireConstraintColumns(t, pk, []int16{1, 2})
	if pk.Name != "attach_index_table_pkey" {
		t.Fatalf("constraint name: got %q, want %q", pk.Name, "attach_index_table_pkey")
	}

	indexes := c.indexesByRel[rel.OID]
	if len(indexes) != 1 {
		t.Fatalf("index count: got %d, want 1", len(indexes))
	}
	if indexes[0].Name != "attach_index_table_pkey" {
		t.Fatalf("attached index name: got %q, want %q", indexes[0].Name, "attach_index_table_pkey")
	}
	if indexes[0].ConstraintOID != pk.OID {
		t.Fatalf("attached index constraint OID: got %d, want %d", indexes[0].ConstraintOID, pk.OID)
	}
}

func TestAlterTableAddIndexConstraintUsesExistingIndexOID(t *testing.T) {
	c := New()

	execSQL(t, c, `
CREATE TABLE attach_index_oid_table (
    id integer NOT NULL,
    bucket integer NOT NULL
);
CREATE UNIQUE INDEX attach_index_oid_table_idx ON attach_index_oid_table (id, bucket);
`)

	schema := c.schemaByName["public"]
	idx := schema.Indexes["attach_index_oid_table_idx"]
	if idx == nil {
		t.Fatal("setup index not found")
	}

	atCmd := &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_AddIndexConstraint),
		Def: &nodes.IndexStmt{
			Idxname:  "attach_index_oid_table_pkey",
			Primary:  true,
			IndexOid: nodes.Oid(idx.OID),
		},
	}
	if err := c.AlterTableStmt(makeAlterTableStmt("", "attach_index_oid_table", atCmd)); err != nil {
		t.Fatalf("AlterTableStmt: %v", err)
	}

	_, rel, err := c.findRelation("", "attach_index_oid_table")
	if err != nil {
		t.Fatal(err)
	}
	pk := requireSingleConstraintOfType(t, c, rel, ConstraintPK)
	requireConstraintColumns(t, pk, []int16{1, 2})
	if pk.IndexOID != idx.OID {
		t.Fatalf("constraint index OID: got %d, want %d", pk.IndexOID, idx.OID)
	}
	if idx.Name != "attach_index_oid_table_pkey" {
		t.Fatalf("attached index name: got %q, want %q", idx.Name, "attach_index_oid_table_pkey")
	}
}

func TestPartitionChildAddUniqueReusesPartitionConstraint(t *testing.T) {
	c := New()

	stmts := parseStmts(t, `
CREATE TABLE partition_unique_parent (
    id integer NOT NULL,
    bucket integer NOT NULL,
    CONSTRAINT partition_unique_parent_key UNIQUE (id, bucket)
) PARTITION BY RANGE (bucket);
CREATE TABLE partition_unique_child
    PARTITION OF partition_unique_parent FOR VALUES FROM (0) TO (10);
ALTER TABLE ONLY partition_unique_child
    ADD CONSTRAINT partition_unique_child_key UNIQUE (id, bucket);
`)
	for _, stmt := range stmts {
		if err := c.ProcessUtility(stmt); err != nil {
			t.Fatalf("ProcessUtility: %v", err)
		}
	}

	_, child, err := c.findRelation("", "partition_unique_child")
	if err != nil {
		t.Fatal(err)
	}

	unique := requireSingleConstraintOfType(t, c, child, ConstraintUnique)
	requireConstraintColumns(t, unique, []int16{1, 2})
	if unique.Name != "partition_unique_child_key" {
		t.Fatalf("partition child UNIQUE name: got %q, want %q", unique.Name, "partition_unique_child_key")
	}
	if unique.ConParentID == 0 {
		t.Fatal("partition child UNIQUE should retain parent constraint link")
	}
}

func TestPartitionChildAddPrimaryKeyUsingIndexReusesPartitionConstraint(t *testing.T) {
	c := New()

	stmts := parseStmts(t, `
CREATE TABLE partition_pk_index_parent (
    id integer NOT NULL,
    bucket integer NOT NULL,
    CONSTRAINT partition_pk_index_parent_pkey PRIMARY KEY (id, bucket)
) PARTITION BY RANGE (bucket);
CREATE TABLE partition_pk_index_child
    PARTITION OF partition_pk_index_parent FOR VALUES FROM (0) TO (10);
ALTER TABLE ONLY partition_pk_index_child
    ADD CONSTRAINT partition_pk_index_child_pkey PRIMARY KEY USING INDEX partition_pk_index_child_pkey;
`)
	for _, stmt := range stmts {
		if err := c.ProcessUtility(stmt); err != nil {
			t.Fatalf("ProcessUtility: %v", err)
		}
	}

	_, child, err := c.findRelation("", "partition_pk_index_child")
	if err != nil {
		t.Fatal(err)
	}

	pk := requireSingleConstraintOfType(t, c, child, ConstraintPK)
	requireConstraintColumns(t, pk, []int16{1, 2})
	if pk.Name != "partition_pk_index_child_pkey" {
		t.Fatalf("partition child PK name: got %q, want %q", pk.Name, "partition_pk_index_child_pkey")
	}
	if pk.ConParentID == 0 {
		t.Fatal("partition child PK should retain parent constraint link")
	}
}

func TestPartitionChildAddUniqueUsingIndexReusesPartitionConstraint(t *testing.T) {
	c := New()

	stmts := parseStmts(t, `
CREATE TABLE pui_parent (
    id integer NOT NULL,
    bucket integer NOT NULL,
    CONSTRAINT pui_parent_key UNIQUE (id, bucket)
) PARTITION BY RANGE (bucket);
CREATE TABLE pui_child
    PARTITION OF pui_parent FOR VALUES FROM (0) TO (10);
ALTER TABLE ONLY pui_child
    ADD CONSTRAINT pui_child_key UNIQUE USING INDEX pui_child_pui_parent_key_idx;
`)
	for _, stmt := range stmts {
		if err := c.ProcessUtility(stmt); err != nil {
			t.Fatalf("ProcessUtility: %v", err)
		}
	}

	_, child, err := c.findRelation("", "pui_child")
	if err != nil {
		t.Fatal(err)
	}

	unique := requireSingleConstraintOfType(t, c, child, ConstraintUnique)
	requireConstraintColumns(t, unique, []int16{1, 2})
	if unique.Name != "pui_child_key" {
		t.Fatalf("partition child UNIQUE name: got %q, want %q", unique.Name, "pui_child_key")
	}
	if unique.ConParentID == 0 {
		t.Fatal("partition child UNIQUE should retain parent constraint link")
	}
	if idx := c.indexes[unique.IndexOID]; idx == nil || idx.Name != "pui_child_key" {
		t.Fatalf("partition child UNIQUE index was not renamed to the constraint name")
	}
}

func requireSingleConstraintOfType(t *testing.T, c *Catalog, rel *Relation, typ ConstraintType) *Constraint {
	t.Helper()

	var matches []*Constraint
	for _, con := range c.consByRel[rel.OID] {
		if con.Type == typ {
			matches = append(matches, con)
		}
	}
	if len(matches) != 1 {
		t.Fatalf("constraint count for type %q on %s: got %d, want 1", typ, rel.Name, len(matches))
	}
	return matches[0]
}

func requireConstraintColumns(t *testing.T, con *Constraint, want []int16) {
	t.Helper()

	if len(con.Columns) != len(want) {
		t.Fatalf("constraint columns: got %v, want %v", con.Columns, want)
	}
	for i := range want {
		if con.Columns[i] != want[i] {
			t.Fatalf("constraint columns: got %v, want %v", con.Columns, want)
		}
	}
}
