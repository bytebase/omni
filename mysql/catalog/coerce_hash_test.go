package catalog

import (
	"fmt"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
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

	cases := []struct {
		name, sql, index, want string
	}{
		{"innodb-unique-hash", "CREATE DATABASE d; USE d; CREATE TABLE t (c INT, UNIQUE KEY uk (c) USING HASH) ENGINE=InnoDB;", "uk", ""},
		{"innodb-key-hash", "CREATE DATABASE d; USE d; CREATE TABLE t (c INT, KEY k (c) USING HASH) ENGINE=InnoDB;", "k", ""},
		{"innodb-noengine-hash", "CREATE DATABASE d; USE d; CREATE TABLE t (c INT, KEY k (c) USING HASH);", "k", ""},
		{"innodb-btree-kept", "CREATE DATABASE d; USE d; CREATE TABLE t (c INT, KEY k (c) USING BTREE) ENGINE=InnoDB;", "k", "BTREE"},
		{"memory-unique-hash-kept", "CREATE DATABASE d; USE d; CREATE TABLE t (c INT, UNIQUE KEY uk (c) USING HASH) ENGINE=MEMORY;", "uk", "HASH"},
		{"memory-key-hash-kept", "CREATE DATABASE d; USE d; CREATE TABLE t (c INT, KEY k (c) USING HASH) ENGINE=MEMORY;", "k", "HASH"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := idxType(c.sql, "t", c.index); got != c.want {
				t.Errorf("IndexType = %q, want %q", got, c.want)
			}
		})
	}
}

// hashProbe is one CREATE-TABLE form whose index uses USING HASH; it must round-trip empty against
// the engine readback (folded away on InnoDB, preserved on MEMORY).
type hashProbe struct {
	id        string
	table     string
	createSQL string
	versions  []Version
}

func hashRoundTripProbes() []hashProbe {
	return []hashProbe{
		// InnoDB: USING HASH must fold so the round-trip is empty on BOTH versions (5.7 keeps the
		// echo, 8.0 drops it — both must collapse to the no-USING form).
		{"innodb-unique-hash", "t",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20), UNIQUE KEY uk (name) USING HASH) ENGINE=InnoDB", both()},
		{"innodb-key-hash", "t",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20), KEY k (name) USING HASH) ENGINE=InnoDB", both()},
		// InnoDB with no explicit ENGINE (defaults to InnoDB): USING HASH still folds.
		{"innodb-default-engine-hash", "t",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20), KEY k (name) USING HASH)", both()},
		// InnoDB USING BTREE must round-trip too (echoed verbatim on both versions; NOT folded).
		{"innodb-key-btree", "t",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20), KEY k (name) USING BTREE) ENGINE=InnoDB", both()},
		// MEMORY: HASH is the real index type and is echoed on both versions, so it round-trips as
		// HASH (the fold must NOT fire for MEMORY).
		{"memory-unique-hash", "t",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20), UNIQUE KEY uk (name) USING HASH) ENGINE=MEMORY", both()},
		{"memory-key-hash", "t",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20), KEY k (name) USING HASH) ENGINE=MEMORY", both()},
	}
}

// TestOracle_InnoDBHashRoundTrip proves the USING HASH fold against the LIVE engines (5.7 :13307,
// 8.0 :13306): every probe's user form vs its engine readback diffs empty on both versions.
func TestOracle_InnoDBHashRoundTrip(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		for _, p := range hashRoundTripProbes() {
			if !containsVersion(p.versions, version) {
				continue
			}
			t.Run(fmt.Sprintf("%s/%s", o.name, p.id), func(t *testing.T) {
				db := "hash_" + strings.ReplaceAll(p.id, "-", "_")
				assertDiffEmptyAgainstReadback(t, o, db, p.createSQL, p.table)
			})
		}
	}
}
