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
			name:  "CREATE INDEX with the type in front of ON",
			ddl:   []string{"CREATE INDEX i USING BTREE ON t (a)"},
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
			name:  "a key definition naming the type twice keeps the last",
			ddl:   []string{"CREATE TABLE kt (a INT, KEY k USING HASH (a) USING BTREE)"},
			table: "kt",
			index: "k",
			want:  "BTREE",
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
