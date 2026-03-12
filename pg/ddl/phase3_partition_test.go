package ddl

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
	pgparser "github.com/bytebase/omni/pg/parser"
)

// =============================================================================
// Phase 3: Partition Completeness Tests
// =============================================================================

// setupPartitioned creates a partitioned parent table and returns the catalog.
func setupPartitioned(t *testing.T, extraSQL ...string) *Catalog {
	t.Helper()
	c := New()
	sql := `CREATE TABLE parent (id int, name text, val int) PARTITION BY LIST (id)`
	list, err := pgparser.Parse(sql)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	for _, item := range list.Items {
		if raw, ok := item.(*nodes.RawStmt); ok {
			item = raw.Stmt
		}
		if err := c.ProcessUtility(item); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}
	for _, s := range extraSQL {
		list, err := pgparser.Parse(s)
		if err != nil {
			t.Fatalf("parse extra: %v", err)
		}
		for _, item := range list.Items {
			if raw, ok := item.(*nodes.RawStmt); ok {
				item = raw.Stmt
			}
			if err := c.ProcessUtility(item); err != nil {
				t.Fatalf("extra setup: %v", err)
			}
		}
	}
	return c
}

// TestPhase3_PKIndexCloned tests that PK index is cloned to partition.
func TestPhase3_PKIndexCloned(t *testing.T) {
	c := setupPartitioned(t,
		"ALTER TABLE parent ADD CONSTRAINT parent_pkey PRIMARY KEY (id)",
		"CREATE TABLE child PARTITION OF parent FOR VALUES IN (1)",
	)

	schema := c.schemaByName["public"]

	// Child should have a cloned PK index.
	var found bool
	for _, idx := range c.indexesByRel[schema.Relations["child"].OID] {
		if idx.IsPrimary {
			found = true
			break
		}
	}
	if !found {
		t.Error("partition should have a cloned primary key index")
	}
}

// TestPhase3_PKNotNullPropagated tests NOT NULL is set on PK columns in partition.
func TestPhase3_PKNotNullPropagated(t *testing.T) {
	c := setupPartitioned(t,
		"ALTER TABLE parent ADD CONSTRAINT parent_pkey PRIMARY KEY (id)",
		"CREATE TABLE child PARTITION OF parent FOR VALUES IN (1)",
	)

	_, rel, err := c.findRelation("", "child")
	if err != nil {
		t.Fatal(err)
	}
	idx, ok := rel.colByName["id"]
	if !ok {
		t.Fatal("id column not found")
	}
	if !rel.Columns[idx].NotNull {
		t.Error("PK column 'id' should be NOT NULL in partition")
	}
}

// TestPhase3_UniqueIndexCloned tests that UNIQUE index is cloned to partition.
func TestPhase3_UniqueIndexCloned(t *testing.T) {
	c := setupPartitioned(t,
		"CREATE UNIQUE INDEX parent_unique_idx ON parent (id, name)",
		"CREATE TABLE child PARTITION OF parent FOR VALUES IN (1)",
	)

	schema := c.schemaByName["public"]
	var found bool
	for _, idx := range c.indexesByRel[schema.Relations["child"].OID] {
		if idx.IsUnique {
			found = true
			break
		}
	}
	if !found {
		t.Error("partition should have a cloned unique index")
	}
}

// TestPhase3_RegularIndexCloned tests that a regular (non-unique) index is cloned.
func TestPhase3_RegularIndexCloned(t *testing.T) {
	c := setupPartitioned(t,
		"CREATE INDEX parent_name_idx ON parent (name)",
		"CREATE TABLE child PARTITION OF parent FOR VALUES IN (1)",
	)

	schema := c.schemaByName["public"]
	childRel := schema.Relations["child"]
	var found bool
	for _, idx := range c.indexesByRel[childRel.OID] {
		if !idx.IsUnique && !idx.IsPrimary {
			found = true
			break
		}
	}
	if !found {
		t.Error("partition should have a cloned regular index")
	}
}

// TestPhase3_TriggerClonedForEachRow tests that FOR EACH ROW trigger is cloned.
func TestPhase3_TriggerClonedForEachRow(t *testing.T) {
	c := setupPartitioned(t,
		"CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END $$",
		"CREATE TRIGGER parent_trg BEFORE INSERT ON parent FOR EACH ROW EXECUTE FUNCTION trig_fn()",
		"CREATE TABLE child PARTITION OF parent FOR VALUES IN (1)",
	)

	_, child, err := c.findRelation("", "child")
	if err != nil {
		t.Fatal(err)
	}

	trigs := c.triggersByRel[child.OID]
	var found bool
	for _, trig := range trigs {
		if trig.Name == "parent_trg" && trig.ForEachRow {
			found = true
			break
		}
	}
	if !found {
		t.Error("partition should have cloned FOR EACH ROW trigger")
	}
}

// TestPhase3_TriggerStatementNotCloned tests that FOR EACH STATEMENT trigger is NOT cloned.
func TestPhase3_TriggerStatementNotCloned(t *testing.T) {
	c := setupPartitioned(t,
		"CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END $$",
		"CREATE TRIGGER parent_stmt_trg AFTER INSERT ON parent FOR EACH STATEMENT EXECUTE FUNCTION trig_fn()",
		"CREATE TABLE child PARTITION OF parent FOR VALUES IN (1)",
	)

	_, child, err := c.findRelation("", "child")
	if err != nil {
		t.Fatal(err)
	}

	trigs := c.triggersByRel[child.OID]
	for _, trig := range trigs {
		if trig.Name == "parent_stmt_trg" {
			t.Error("FOR EACH STATEMENT trigger should NOT be cloned to partition")
		}
	}
}

// TestPhase3_FKCloned tests that FK constraint is cloned to partition.
func TestPhase3_FKCloned(t *testing.T) {
	c := New()
	stmts := parseStmts(t, `
		CREATE TABLE ref (id int PRIMARY KEY);
		CREATE TABLE parent (id int, ref_id int) PARTITION BY LIST (id);
		ALTER TABLE parent ADD CONSTRAINT parent_fk FOREIGN KEY (ref_id) REFERENCES ref(id);
		CREATE TABLE child PARTITION OF parent FOR VALUES IN (1);
	`)
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	_, child, err := c.findRelation("", "child")
	if err != nil {
		t.Fatal(err)
	}

	var found bool
	for _, con := range c.consByRel[child.OID] {
		if con.Type == ConstraintFK {
			found = true
			break
		}
	}
	if !found {
		t.Error("partition should have cloned FK constraint")
	}
}

// TestPhase3_AttachPartitionClonesIndexes tests that ATTACH PARTITION clones indexes.
func TestPhase3_AttachPartitionClonesIndexes(t *testing.T) {
	c := setupPartitioned(t,
		"CREATE INDEX parent_name_idx ON parent (name)",
	)

	// Create a standalone table and attach it.
	stmts := parseStmts(t, "CREATE TABLE standalone (id int, name text, val int)")
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatal(err)
		}
	}

	stmts = parseStmts(t, "ALTER TABLE parent ATTACH PARTITION standalone FOR VALUES IN (42)")
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatal(err)
		}
	}

	schema := c.schemaByName["public"]
	standalonRel := schema.Relations["standalone"]
	idxes := c.indexesByRel[standalonRel.OID]
	if len(idxes) == 0 {
		t.Error("ATTACH PARTITION should clone parent indexes to the attached table")
	}
}

// TestPhase3_AttachPartitionClonesTriggers tests that ATTACH PARTITION clones triggers.
func TestPhase3_AttachPartitionClonesTriggers(t *testing.T) {
	c := setupPartitioned(t,
		"CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END $$",
		"CREATE TRIGGER parent_trg BEFORE INSERT ON parent FOR EACH ROW EXECUTE FUNCTION trig_fn()",
	)

	stmts := parseStmts(t, "CREATE TABLE standalone (id int, name text, val int)")
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatal(err)
		}
	}

	stmts = parseStmts(t, "ALTER TABLE parent ATTACH PARTITION standalone FOR VALUES IN (42)")
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatal(err)
		}
	}

	_, standalone, err := c.findRelation("", "standalone")
	if err != nil {
		t.Fatal(err)
	}

	trigs := c.triggersByRel[standalone.OID]
	var found bool
	for _, trig := range trigs {
		if trig.Name == "parent_trg" {
			found = true
			break
		}
	}
	if !found {
		t.Error("ATTACH PARTITION should clone parent triggers")
	}
}

// TestPhase3_AttachPartitionClonesFK tests ATTACH PARTITION clones FK constraints.
func TestPhase3_AttachPartitionClonesFK(t *testing.T) {
	c := New()
	stmts := parseStmts(t, `
		CREATE TABLE ref (id int PRIMARY KEY);
		CREATE TABLE parent (id int, ref_id int) PARTITION BY LIST (id);
		ALTER TABLE parent ADD CONSTRAINT parent_fk FOREIGN KEY (ref_id) REFERENCES ref(id);
		CREATE TABLE standalone (id int, ref_id int);
		ALTER TABLE parent ATTACH PARTITION standalone FOR VALUES IN (42);
	`)
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	_, standalone, err := c.findRelation("", "standalone")
	if err != nil {
		t.Fatal(err)
	}

	var found bool
	for _, con := range c.consByRel[standalone.OID] {
		if con.Type == ConstraintFK {
			found = true
			break
		}
	}
	if !found {
		t.Error("ATTACH PARTITION should clone FK constraints")
	}
}

// TestPhase3_ConstraintClonedWithParentRef tests that cloned constraint has ConParentID set.
func TestPhase3_ConstraintClonedWithParentRef(t *testing.T) {
	c := setupPartitioned(t,
		"ALTER TABLE parent ADD CONSTRAINT parent_pkey PRIMARY KEY (id)",
		"CREATE TABLE child PARTITION OF parent FOR VALUES IN (1)",
	)

	_, child, err := c.findRelation("", "child")
	if err != nil {
		t.Fatal(err)
	}

	for _, con := range c.consByRel[child.OID] {
		if con.Type == 'p' {
			if con.ConParentID == 0 {
				t.Error("cloned PK constraint should have ConParentID set")
			}
			return
		}
	}
	t.Error("child should have PK constraint")
}
