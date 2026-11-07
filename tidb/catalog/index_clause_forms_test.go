package catalog

import "testing"

// The index type an index DDL statement declares has to reach the catalog, and
// a DROP INDEX has to drop the index it names. A statement that parses and then
// leaves the schema wrong is worse than one that is rejected, because Bytebase
// would review later migrations against a schema that never existed, and the
// corpus gate in tidb/parser cannot see any of it: it records only whether a
// statement parses.

// indexNamed returns the index of the given name on testdb.t, or nil.
func indexNamed(c *Catalog, table, name string) *Index {
	db := c.GetDatabase("testdb")
	if db == nil {
		return nil
	}
	tbl := db.Tables[table]
	if tbl == nil {
		return nil
	}
	for _, idx := range tbl.Indexes {
		if idx.Name == name {
			return idx
		}
	}
	return nil
}

// TestIndexTypeReachesTheCatalog runs each statement that declares an index
// type through Exec and reads the type back off the index it made.
func TestIndexTypeReachesTheCatalog(t *testing.T) {
	for _, tc := range []struct {
		name  string
		ddl   []string
		table string
		index string
		want  string
	}{
		{
			name:  "CREATE INDEX with the type after the column list",
			ddl:   []string{"CREATE INDEX i ON t (a) USING BTREE"},
			table: "t",
			index: "i",
			want:  "BTREE",
		},
		{
			name:  "CREATE INDEX with the type spelled TYPE",
			ddl:   []string{"CREATE INDEX i ON t (a) TYPE BTREE"},
			table: "t",
			index: "i",
			want:  "BTREE",
		},
		{
			name:  "CREATE INDEX with the type in front of ON",
			ddl:   []string{"CREATE INDEX i TYPE BTREE ON t (a)"},
			table: "t",
			index: "i",
			want:  "BTREE",
		},
		{
			name:  "CREATE INDEX naming the type twice keeps the last",
			ddl:   []string{"CREATE INDEX i USING HASH ON t (a) USING BTREE"},
			table: "t",
			index: "i",
			want:  "BTREE",
		},
		{
			name:  "a key definition with the type spelled TYPE",
			ddl:   []string{"CREATE TABLE kt (a INT, KEY k TYPE BTREE (a))"},
			table: "kt",
			index: "k",
			want:  "BTREE",
		},
		{
			name:  "a key definition naming the type twice keeps the last",
			ddl:   []string{"CREATE TABLE kt (a INT, KEY k USING HASH (a) USING BTREE)"},
			table: "kt",
			index: "k",
			want:  "BTREE",
		},
		{
			name:  "ALTER TABLE ADD KEY with the type spelled TYPE",
			ddl:   []string{"ALTER TABLE t ADD KEY k (a) TYPE BTREE"},
			table: "t",
			index: "k",
			want:  "BTREE",
		},
		{
			name:  "ALTER TABLE ADD KEY with the type in front of the column list",
			ddl:   []string{"ALTER TABLE t ADD KEY k TYPE BTREE (a)"},
			table: "t",
			index: "k",
			want:  "BTREE",
		},
		{
			name:  "TiDB normalizes TYPE HASH even with ENGINE MEMORY",
			ddl:   []string{"CREATE TABLE mt (a INT) ENGINE=MEMORY", "CREATE INDEX i ON mt (a) TYPE HASH"},
			table: "mt",
			index: "i",
			want:  "BTREE",
		},
		{
			// CREATE INDEX itself refuses HNSW on an index that is no vector index: its
			// grammar action answers 1149. A key definition takes it.
			name:  "key definition with TiDB's own index type",
			ddl:   []string{"CREATE TABLE kh (a INT, KEY i (a) USING HNSW)"},
			table: "kh",
			index: "i",
			want:  "HNSW",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := wtSetup(t)
			wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT)")
			for _, ddl := range tc.ddl {
				wtExec(t, c, ddl)
			}
			idx := indexNamed(c, tc.table, tc.index)
			if idx == nil {
				t.Fatalf("index %q on %q is not in the catalog", tc.index, tc.table)
			}
			if idx.IndexType != tc.want {
				t.Errorf("IndexType = %q, want %q", idx.IndexType, tc.want)
			}
		})
	}
}

// TestDropIndexIfExists covers the three outcomes IF EXISTS has. The index is
// dropped when it is there; a missing index leaves the table as it was, with no
// error, where the bare statement is a 1091; and a missing table is still an
// error, because IF EXISTS speaks for the index and not for the table.
func TestDropIndexIfExists(t *testing.T) {
	t.Run("drops the index it names", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY idx (a))")
		wtExec(t, c, "DROP INDEX IF EXISTS idx ON t")
		if idx := indexNamed(c, "t", "idx"); idx != nil {
			t.Error("idx is still in the catalog after DROP INDEX IF EXISTS")
		}
	})

	t.Run("a missing index is no error and no change", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY idx (a))")
		results, err := c.Exec("DROP INDEX IF EXISTS nosuchindex ON t", nil)
		if err != nil {
			t.Fatalf("parse error: %v", err)
		}
		for _, r := range results {
			assertNoError(t, r.Error)
		}
		if idx := indexNamed(c, "t", "idx"); idx == nil {
			t.Error("idx is gone after a DROP INDEX IF EXISTS that named another index")
		}
	})

	t.Run("without IF EXISTS a missing index is a 1091", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY idx (a))")
		results, err := c.Exec("DROP INDEX nosuchindex ON t", nil)
		if err != nil {
			t.Fatalf("parse error: %v", err)
		}
		if len(results) != 1 {
			t.Fatalf("Exec read %d statements, want 1", len(results))
		}
		assertError(t, results[0].Error, ErrCantDropKey)
	})

	t.Run("a missing table is still an error", func(t *testing.T) {
		c := wtSetup(t)
		results, err := c.Exec("DROP INDEX IF EXISTS idx ON nosuchtable", nil)
		if err != nil {
			t.Fatalf("parse error: %v", err)
		}
		if len(results) != 1 {
			t.Fatalf("Exec read %d statements, want 1", len(results))
		}
		assertError(t, results[0].Error, ErrNoSuchTable)
	})
}

// TestDropIndexWithAlgorithmAndLock drops the index the statement names, with
// the clauses that close it taken as clauses rather than as the name. A
// statement whose ALGORITHM = DEFAULT were read as the index name would parse
// and then drop nothing, or the wrong index.
func TestDropIndexWithAlgorithmAndLock(t *testing.T) {
	for _, drop := range []string{
		"DROP INDEX idx ON t ALGORITHM = DEFAULT",
		"DROP INDEX idx ON t ALGORITHM DEFAULT",
		"DROP INDEX idx ON t LOCK = DEFAULT",
		"DROP INDEX idx ON t LOCK DEFAULT",
		"DROP INDEX idx ON t ALGORITHM = DEFAULT LOCK = DEFAULT",
		"DROP INDEX idx ON t LOCK = DEFAULT ALGORITHM = DEFAULT",
		"DROP INDEX idx ON t ALGORITHM = INPLACE LOCK = NONE",
		"DROP INDEX IF EXISTS idx ON t ALGORITHM = DEFAULT",
	} {
		t.Run(drop, func(t *testing.T) {
			c := wtSetup(t)
			wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY idx (a), KEY keep (b))")
			wtExec(t, c, drop)
			if idx := indexNamed(c, "t", "idx"); idx != nil {
				t.Error("idx is still in the catalog")
			}
			if idx := indexNamed(c, "t", "keep"); idx == nil {
				t.Error("keep is gone: the statement dropped more than it named")
			}
		})
	}
}

// TestStatsLockLeavesTheSchemaAlone runs LOCK STATS and UNLOCK STATS through
// Exec. They change no schema, so the catalog has nothing to do with them, and
// what matters is that they reach Exec without an error and take nothing with
// them: before the grammar landed, `LOCK STATS t` parsed as a LOCK TABLES of a
// table called `t`, which the catalog also ignored, so the two cannot be told
// apart downstream without asserting the schema.
func TestStatsLockLeavesTheSchemaAlone(t *testing.T) {
	for _, stmt := range []string{
		"LOCK STATS t",
		"UNLOCK STATS t",
		"LOCK STATS t, t2",
		"LOCK STATS t PARTITION (p0, p1)",
		"LOCK STATS t PARTITION p0",
		"UNLOCK STATS test.t",
		// A table the catalog has never heard of is no error either: the
		// statement names nothing the catalog keeps.
		"LOCK STATS nosuchtable",
	} {
		t.Run(stmt, func(t *testing.T) {
			c := wtSetup(t)
			wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY idx (a))")
			wtExec(t, c, "CREATE TABLE t2 (id INT PRIMARY KEY)")
			wtExec(t, c, stmt)

			db := c.GetDatabase("testdb")
			if db == nil {
				t.Fatal("database testdb is gone")
			}
			if len(db.Tables) != 2 {
				t.Errorf("tables = %d, want 2", len(db.Tables))
			}
			if idx := indexNamed(c, "t", "idx"); idx == nil {
				t.Error("idx is gone")
			}
			if tbl := db.Tables["t"]; tbl == nil {
				t.Error("t is gone")
			} else if len(tbl.Columns) != 2 {
				t.Errorf("t has %d columns, want 2", len(tbl.Columns))
			}
		})
	}
}
