package catalog

import (
	"strings"
	"testing"
)

func TestWalkThrough_6_1_FKBackingIndexManagement(t *testing.T) {
	// Scenario 1: CREATE TABLE with named FK, no explicit index — implicit index uses constraint name
	t.Run("named_fk_implicit_index_uses_constraint_name", func(t *testing.T) {
		c := setupWithDB(t)
		mustExec(t, c, "CREATE TABLE parent (id INT NOT NULL, PRIMARY KEY (id))")
		mustExec(t, c, `CREATE TABLE child (
			id INT NOT NULL,
			parent_id INT NOT NULL,
			PRIMARY KEY (id),
			CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent(id)
		)`)

		tbl := c.GetDatabase("testdb").GetTable("child")
		if tbl == nil {
			t.Fatal("table child not found")
		}

		// Should have 2 indexes: PRIMARY and fk_parent
		if len(tbl.Indexes) != 2 {
			t.Fatalf("expected 2 indexes, got %d", len(tbl.Indexes))
		}

		var fkIdx *Index
		for _, idx := range tbl.Indexes {
			if idx.Name == "fk_parent" {
				fkIdx = idx
				break
			}
		}
		if fkIdx == nil {
			t.Fatal("expected implicit index named 'fk_parent', not found")
		}
		if len(fkIdx.Columns) != 1 || fkIdx.Columns[0].Name != "parent_id" {
			t.Errorf("expected index on (parent_id), got %v", fkIdx.Columns)
		}
	})

	// Scenario 2: CREATE TABLE with unnamed FK — implicit index uses first column name
	t.Run("unnamed_fk_implicit_index_uses_column_name", func(t *testing.T) {
		c := setupWithDB(t)
		mustExec(t, c, "CREATE TABLE parent (id INT NOT NULL, PRIMARY KEY (id))")
		mustExec(t, c, `CREATE TABLE child (
			id INT NOT NULL,
			parent_id INT NOT NULL,
			PRIMARY KEY (id),
			FOREIGN KEY (parent_id) REFERENCES parent(id)
		)`)

		tbl := c.GetDatabase("testdb").GetTable("child")
		if tbl == nil {
			t.Fatal("table child not found")
		}

		// Should have 2 indexes: PRIMARY and parent_id
		if len(tbl.Indexes) != 2 {
			t.Fatalf("expected 2 indexes, got %d", len(tbl.Indexes))
		}

		var fkIdx *Index
		for _, idx := range tbl.Indexes {
			if idx.Name == "parent_id" {
				fkIdx = idx
				break
			}
		}
		if fkIdx == nil {
			names := make([]string, len(tbl.Indexes))
			for i, idx := range tbl.Indexes {
				names[i] = idx.Name
			}
			t.Fatalf("expected implicit index named 'parent_id', not found; indexes: %v", names)
		}
	})

	// Scenario 3: CREATE TABLE with explicit index on FK columns, then named FK — no duplicate index
	t.Run("explicit_index_before_named_fk_no_duplicate", func(t *testing.T) {
		c := setupWithDB(t)
		mustExec(t, c, "CREATE TABLE parent (id INT NOT NULL, PRIMARY KEY (id))")
		mustExec(t, c, `CREATE TABLE child (
			id INT NOT NULL,
			parent_id INT NOT NULL,
			PRIMARY KEY (id),
			INDEX idx_parent (parent_id),
			CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent(id)
		)`)

		tbl := c.GetDatabase("testdb").GetTable("child")
		if tbl == nil {
			t.Fatal("table child not found")
		}

		// Should have 2 indexes: PRIMARY and idx_parent (no duplicate fk_parent)
		if len(tbl.Indexes) != 2 {
			names := make([]string, len(tbl.Indexes))
			for i, idx := range tbl.Indexes {
				names[i] = idx.Name
			}
			t.Fatalf("expected 2 indexes (no duplicate), got %d: %v", len(tbl.Indexes), names)
		}

		// The existing index should be idx_parent, not fk_parent
		found := false
		for _, idx := range tbl.Indexes {
			if idx.Name == "idx_parent" {
				found = true
			}
		}
		if !found {
			t.Error("expected idx_parent index to exist")
		}
	})

	// Scenario 4: CREATE TABLE with FK on column already in UNIQUE KEY — no duplicate index
	t.Run("fk_on_unique_key_column_no_duplicate", func(t *testing.T) {
		c := setupWithDB(t)
		mustExec(t, c, "CREATE TABLE parent (id INT NOT NULL, PRIMARY KEY (id))")
		mustExec(t, c, `CREATE TABLE child (
			id INT NOT NULL,
			parent_id INT NOT NULL,
			PRIMARY KEY (id),
			UNIQUE KEY uk_parent (parent_id),
			CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent(id)
		)`)

		tbl := c.GetDatabase("testdb").GetTable("child")
		if tbl == nil {
			t.Fatal("table child not found")
		}

		// Should have 2 indexes: PRIMARY and uk_parent (no duplicate fk_parent)
		if len(tbl.Indexes) != 2 {
			names := make([]string, len(tbl.Indexes))
			for i, idx := range tbl.Indexes {
				names[i] = idx.Name
			}
			t.Fatalf("expected 2 indexes (no duplicate), got %d: %v", len(tbl.Indexes), names)
		}
	})

	// Scenario 5: CREATE TABLE with FK on column already in PRIMARY KEY — no duplicate index
	t.Run("fk_on_primary_key_column_no_duplicate", func(t *testing.T) {
		c := setupWithDB(t)
		mustExec(t, c, "CREATE TABLE parent (id INT NOT NULL, PRIMARY KEY (id))")
		mustExec(t, c, `CREATE TABLE child (
			parent_id INT NOT NULL,
			PRIMARY KEY (parent_id),
			CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent(id)
		)`)

		tbl := c.GetDatabase("testdb").GetTable("child")
		if tbl == nil {
			t.Fatal("table child not found")
		}

		// Should have only 1 index: PRIMARY (no duplicate for FK)
		if len(tbl.Indexes) != 1 {
			names := make([]string, len(tbl.Indexes))
			for i, idx := range tbl.Indexes {
				names[i] = idx.Name
			}
			t.Fatalf("expected 1 index (PRIMARY only), got %d: %v", len(tbl.Indexes), names)
		}
		if !tbl.Indexes[0].Primary {
			t.Error("expected the only index to be PRIMARY")
		}
	})

	// Scenario 6: CREATE TABLE with multi-column FK, partial index exists — implicit index still created
	t.Run("multi_column_fk_partial_index_still_creates_implicit", func(t *testing.T) {
		c := setupWithDB(t)
		mustExec(t, c, `CREATE TABLE parent (
			a INT NOT NULL,
			b INT NOT NULL,
			PRIMARY KEY (a, b)
		)`)
		mustExec(t, c, `CREATE TABLE child (
			id INT NOT NULL,
			a INT NOT NULL,
			b INT NOT NULL,
			PRIMARY KEY (id),
			INDEX idx_a (a),
			CONSTRAINT fk_ab FOREIGN KEY (a, b) REFERENCES parent(a, b)
		)`)

		tbl := c.GetDatabase("testdb").GetTable("child")
		if tbl == nil {
			t.Fatal("table child not found")
		}

		// idx_a only covers (a), not (a, b) — so an implicit index fk_ab should be created
		// Expected: PRIMARY, idx_a, fk_ab = 3 indexes
		if len(tbl.Indexes) != 3 {
			names := make([]string, len(tbl.Indexes))
			for i, idx := range tbl.Indexes {
				names[i] = idx.Name
			}
			t.Fatalf("expected 3 indexes, got %d: %v", len(tbl.Indexes), names)
		}

		var fkIdx *Index
		for _, idx := range tbl.Indexes {
			if idx.Name == "fk_ab" {
				fkIdx = idx
				break
			}
		}
		if fkIdx == nil {
			t.Fatal("expected implicit index named 'fk_ab' for multi-column FK")
		}
		if len(fkIdx.Columns) != 2 {
			t.Errorf("expected 2 columns in fk_ab index, got %d", len(fkIdx.Columns))
		}
	})

	// Scenario 7: ALTER TABLE ADD FK when column already has index — no duplicate index
	t.Run("alter_add_fk_column_has_index_no_duplicate", func(t *testing.T) {
		c := setupWithDB(t)
		mustExec(t, c, "CREATE TABLE parent (id INT NOT NULL, PRIMARY KEY (id))")
		mustExec(t, c, `CREATE TABLE child (
			id INT NOT NULL,
			parent_id INT NOT NULL,
			PRIMARY KEY (id),
			INDEX idx_parent (parent_id)
		)`)
		mustExec(t, c, "ALTER TABLE child ADD CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent(id)")

		tbl := c.GetDatabase("testdb").GetTable("child")
		if tbl == nil {
			t.Fatal("table child not found")
		}

		// Should have 2 indexes: PRIMARY and idx_parent (no duplicate fk_parent)
		if len(tbl.Indexes) != 2 {
			names := make([]string, len(tbl.Indexes))
			for i, idx := range tbl.Indexes {
				names[i] = idx.Name
			}
			t.Fatalf("expected 2 indexes (no duplicate), got %d: %v", len(tbl.Indexes), names)
		}

		// FK constraint should exist
		var fk *Constraint
		for _, con := range tbl.Constraints {
			if con.Type == ConForeignKey && con.Name == "fk_parent" {
				fk = con
				break
			}
		}
		if fk == nil {
			t.Fatal("FK constraint fk_parent not found")
		}
	})

	// Scenario 8: ALTER TABLE ADD FK when column has no index — implicit index created
	t.Run("alter_add_fk_no_index_creates_implicit", func(t *testing.T) {
		c := setupWithDB(t)
		mustExec(t, c, "CREATE TABLE parent (id INT NOT NULL, PRIMARY KEY (id))")
		mustExec(t, c, `CREATE TABLE child (
			id INT NOT NULL,
			parent_id INT NOT NULL,
			PRIMARY KEY (id)
		)`)
		mustExec(t, c, "ALTER TABLE child ADD CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent(id)")

		tbl := c.GetDatabase("testdb").GetTable("child")
		if tbl == nil {
			t.Fatal("table child not found")
		}

		// Should have 2 indexes: PRIMARY and fk_parent (implicit)
		if len(tbl.Indexes) != 2 {
			names := make([]string, len(tbl.Indexes))
			for i, idx := range tbl.Indexes {
				names[i] = idx.Name
			}
			t.Fatalf("expected 2 indexes, got %d: %v", len(tbl.Indexes), names)
		}

		var fkIdx *Index
		for _, idx := range tbl.Indexes {
			if idx.Name == "fk_parent" {
				fkIdx = idx
				break
			}
		}
		if fkIdx == nil {
			t.Fatal("expected implicit index named 'fk_parent'")
		}
	})

	// Scenario 9: ALTER TABLE DROP FOREIGN KEY — FK removed but backing index remains
	t.Run("drop_fk_keeps_backing_index", func(t *testing.T) {
		c := setupWithDB(t)
		mustExec(t, c, "CREATE TABLE parent (id INT NOT NULL, PRIMARY KEY (id))")
		mustExec(t, c, `CREATE TABLE child (
			id INT NOT NULL,
			parent_id INT NOT NULL,
			PRIMARY KEY (id),
			CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent(id)
		)`)

		// Verify FK and backing index exist before drop
		tbl := c.GetDatabase("testdb").GetTable("child")
		if len(tbl.Indexes) != 2 {
			t.Fatalf("expected 2 indexes before drop, got %d", len(tbl.Indexes))
		}

		mustExec(t, c, "ALTER TABLE child DROP FOREIGN KEY fk_parent")

		tbl = c.GetDatabase("testdb").GetTable("child")

		// FK constraint should be gone
		for _, con := range tbl.Constraints {
			if con.Type == ConForeignKey && strings.EqualFold(con.Name, "fk_parent") {
				t.Error("FK constraint fk_parent should have been removed")
			}
		}

		// Backing index should remain (MySQL behavior)
		if len(tbl.Indexes) != 2 {
			t.Fatalf("expected 2 indexes after FK drop (backing index remains), got %d", len(tbl.Indexes))
		}

		var fkIdx *Index
		for _, idx := range tbl.Indexes {
			if idx.Name == "fk_parent" {
				fkIdx = idx
				break
			}
		}
		if fkIdx == nil {
			t.Fatal("backing index fk_parent should remain after DROP FOREIGN KEY")
		}
	})

	// Scenario 10: ALTER TABLE DROP FOREIGN KEY, DROP INDEX fk_name — explicit index cleanup after FK drop
	t.Run("drop_fk_then_drop_index_explicit_cleanup", func(t *testing.T) {
		c := setupWithDB(t)
		mustExec(t, c, "CREATE TABLE parent (id INT NOT NULL, PRIMARY KEY (id))")
		mustExec(t, c, `CREATE TABLE child (
			id INT NOT NULL,
			parent_id INT NOT NULL,
			PRIMARY KEY (id),
			CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent(id)
		)`)

		mustExec(t, c, "ALTER TABLE child DROP FOREIGN KEY fk_parent, DROP INDEX fk_parent")

		tbl := c.GetDatabase("testdb").GetTable("child")

		// FK constraint should be gone
		for _, con := range tbl.Constraints {
			if con.Type == ConForeignKey && strings.EqualFold(con.Name, "fk_parent") {
				t.Error("FK constraint fk_parent should have been removed")
			}
		}

		// Backing index should also be gone now
		if len(tbl.Indexes) != 1 {
			names := make([]string, len(tbl.Indexes))
			for i, idx := range tbl.Indexes {
				names[i] = idx.Name
			}
			t.Fatalf("expected 1 index (PRIMARY only) after DROP FK + DROP INDEX, got %d: %v", len(tbl.Indexes), names)
		}
		if !tbl.Indexes[0].Primary {
			t.Error("expected the only remaining index to be PRIMARY")
		}
	})
}
