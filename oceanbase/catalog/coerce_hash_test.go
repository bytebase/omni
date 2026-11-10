package catalog

import (
	"strings"
	"testing"
)

// Unit + oracle proof for the InnoDB USING HASH fold (coerceInnoDBHashIndex). InnoDB silently
// coerces USING HASH to BTREE and its SHOW CREATE echo is version-divergent (5.7 keeps USING HASH,
// 8.0 drops it), so a HASH index on InnoDB phantom-diffed. The fold folds InnoDB HASH to the
// no-explicit-type form on load; MEMORY (where HASH is a real index type) is preserved.

// TestCoerceInnoDBHashIndex_Unit checks the load-time fold deterministically, without an engine:
// InnoDB (and engine-less, which defaults to InnoDB) HASH folds to "", BTREE is untouched, and a
// real-HASH engine (MEMORY) keeps HASH.
func TestCoerceInnoDBHashIndex_Unit(t *testing.T) {
	idxType := func(sql, table, index string) string {
		t.Helper()
		cat, err := LoadSQL(sql)
		if err != nil {
			t.Fatalf("load %q: %v", sql, err)
		}
		for _, db := range cat.Databases() {
			if tbl := db.GetTable(table); tbl != nil {
				for _, idx := range tbl.Indexes {
					if strings.EqualFold(idx.Name, index) {
						return idx.IndexType
					}
				}
				t.Fatalf("index %q not found on %q", index, table)
			}
		}
		t.Fatalf("table %q not found", table)
		return ""
	}

	const innodb = "CREATE DATABASE d; USE d; CREATE TABLE t (c INT) ENGINE=InnoDB;"
	const memory = "CREATE DATABASE d; USE d; CREATE TABLE t (c INT) ENGINE=MEMORY;"
	cases := []struct {
		name, sql, index, want string
	}{
		// CREATE TABLE inline index forms.
		{"innodb-unique-hash", "CREATE DATABASE d; USE d; CREATE TABLE t (c INT, UNIQUE KEY uk (c) USING HASH) ENGINE=InnoDB;", "uk", ""},
		{"innodb-key-hash", "CREATE DATABASE d; USE d; CREATE TABLE t (c INT, KEY k (c) USING HASH) ENGINE=InnoDB;", "k", ""},
		{"innodb-noengine-hash", "CREATE DATABASE d; USE d; CREATE TABLE t (c INT, KEY k (c) USING HASH);", "k", ""},
		{"innodb-btree-kept", "CREATE DATABASE d; USE d; CREATE TABLE t (c INT, KEY k (c) USING BTREE) ENGINE=InnoDB;", "k", "BTREE"},
		{"memory-unique-hash-kept", "CREATE DATABASE d; USE d; CREATE TABLE t (c INT, UNIQUE KEY uk (c) USING HASH) ENGINE=MEMORY;", "uk", "HASH"},
		{"memory-key-hash-kept", "CREATE DATABASE d; USE d; CREATE TABLE t (c INT, KEY k (c) USING HASH) ENGINE=MEMORY;", "k", "HASH"},

		// ALTER TABLE ADD UNIQUE / ADD KEY — the previously-uncovered ALTER paths.
		{"alter-add-unique-hash-innodb", innodb + "ALTER TABLE t ADD UNIQUE uk (c) USING HASH;", "uk", ""},
		{"alter-add-key-hash-innodb", innodb + "ALTER TABLE t ADD KEY k (c) USING HASH;", "k", ""},
		{"alter-add-unique-hash-memory", memory + "ALTER TABLE t ADD UNIQUE uk (c) USING HASH;", "uk", "HASH"},
		{"alter-add-key-hash-memory", memory + "ALTER TABLE t ADD KEY k (c) USING HASH;", "k", "HASH"},

		// Standalone CREATE [UNIQUE] INDEX — trailing USING (the common form the parser routes
		// into Options) and leading USING (before ON). Both positions must resolve + fold.
		{"create-index-trailing-hash-innodb", innodb + "CREATE INDEX k ON t (c) USING HASH;", "k", ""},
		{"create-index-leading-hash-innodb", innodb + "CREATE INDEX k USING HASH ON t (c);", "k", ""},
		{"create-unique-index-hash-innodb", innodb + "CREATE UNIQUE INDEX uk ON t (c) USING HASH;", "uk", ""},
		{"create-index-trailing-hash-memory", memory + "CREATE INDEX k ON t (c) USING HASH;", "k", "HASH"},
		{"create-index-leading-hash-memory", memory + "CREATE INDEX k USING HASH ON t (c);", "k", "HASH"},
		{"create-unique-index-hash-memory", memory + "CREATE UNIQUE INDEX uk ON t (c) USING HASH;", "uk", "HASH"},
		{"create-index-btree-innodb-kept", innodb + "CREATE INDEX k ON t (c) USING BTREE;", "k", "BTREE"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := idxType(c.sql, "t", c.index); got != c.want {
				t.Errorf("IndexType = %q, want %q", got, c.want)
			}
		})
	}
}
