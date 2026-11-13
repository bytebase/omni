package catalog

import (
	"strings"
	"testing"
)

// In-process, hermetic structural tests for the index differ (diff_index.go). They assert the
// per-index SchemaDiff content (right DiffAction per index) for representative change forms,
// loading both sides through the omni catalog so the model state is authentic. The oracle-backed
// idempotence + canonicalization-empty + apply-correctness proofs live in diff_index_oracle_test.go
// and migration_index_oracle_test.go (correctness-protocol.md gates 1 & 2); these lock in gate 3
// (the non-empty index diffs are structurally correct) and run without a live engine.

// findIndexDiff locates an index diff entry within a table diff by name (case-insensitive).
func findIndexDiff(e *TableDiffEntry, name string) *IndexDiffEntry {
	if e == nil {
		return nil
	}
	for i := range e.Indexes {
		if toLower(e.Indexes[i].Name) == toLower(name) {
			return &e.Indexes[i]
		}
	}
	return nil
}

func TestDiffIndex_AddDropModify(t *testing.T) {
	cases := []struct {
		name       string
		from       string
		to         string
		table      string
		idx        string
		wantAction DiffAction
	}{
		{
			"add-secondary",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT)",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a))",
			"t", "ka", DiffAdd,
		},
		{
			"drop-secondary",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a))",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT)",
			"t", "ka", DiffDrop,
		},
		{
			"add-unique",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT)",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, UNIQUE KEY ua (a))",
			"t", "ua", DiffAdd,
		},
		{
			// Same name, columns change a -> a,b: a MODIFY (DROP+ADD at generate time).
			"modify-columns",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY k (a))",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY k (a,b))",
			"t", "k", DiffModify,
		},
		{
			// Non-unique -> unique with the same name is a kind change → MODIFY.
			"modify-kind-unique",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a))",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, UNIQUE KEY k (a))",
			"t", "k", DiffModify,
		},
		{
			"add-primary-key",
			"CREATE TABLE t (id INT NOT NULL, a INT)",
			"CREATE TABLE t (id INT NOT NULL, a INT, PRIMARY KEY (id))",
			"t", "PRIMARY", DiffAdd,
		},
		{
			"drop-primary-key",
			"CREATE TABLE t (id INT NOT NULL, a INT, PRIMARY KEY (id))",
			"CREATE TABLE t (id INT NOT NULL, a INT)",
			"t", "PRIMARY", DiffDrop,
		},
		{
			// COMMENT change on a same-name index → MODIFY.
			"modify-comment",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a) COMMENT 'x')",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a) COMMENT 'y')",
			"t", "k", DiffModify,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			from := loadCat(t, dbDDL+tc.from)
			to := loadCat(t, dbDDL+tc.to)
			d := Diff(from, to)
			te := findTable(d, tc.table)
			if te == nil {
				t.Fatalf("table %s not in diff", tc.table)
			}
			ie := findIndexDiff(te, tc.idx)
			if ie == nil {
				t.Fatalf("index %s not in diff of table %s; entries=%+v", tc.idx, tc.table, te.Indexes)
			}
			if ie.Action != tc.wantAction {
				t.Errorf("index %s: want action %s, got %s", tc.idx, tc.wantAction, ie.Action)
			}
		})
	}
}

// TestDiffIndex_NoPhantomOnReorder asserts that reordering indexes (declaration order) is not a
// change: index identity is the name, and the column-content key is order-independent across
// indexes (only column order WITHIN an index matters). Two tables with the same indexes in a
// different textual order must diff empty.
func TestDiffIndex_NoPhantomOnReorder(t *testing.T) {
	a := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY, x INT, y INT, KEY kx (x), KEY ky (y))")
	b := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY, x INT, y INT, KEY ky (y), KEY kx (x))")
	if d := Diff(a, b); !d.IsEmpty() {
		t.Errorf("reorder of indexes must be empty, got: %s", describeDiff(d))
	}
}

// TestDiffIndex_NoIndexBehindAForeignKey asserts that the loader gives a foreign key no
// index, named or not: OceanBase makes none, where MySQL adds one when no index serves the
// key. TestOracle_OceanBaseForeignKeyIndex asks the engine.
func TestDiffIndex_NoIndexBehindAForeignKey(t *testing.T) {
	cat := loadCat(t, dbDDL+"CREATE TABLE p (id INT NOT NULL PRIMARY KEY);\n"+
		"CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id), FOREIGN KEY (pid) REFERENCES p(id));")
	for _, idx := range mustGetTable(t, cat, "c").Indexes {
		if !idx.Primary {
			t.Errorf("table c declares no index and has %s", idx.Name)
		}
	}
	if d := Diff(cat, cat); !d.IsEmpty() {
		t.Errorf("self-diff must be empty, got: %s", describeDiff(d))
	}
}

// TestDiffIndex_IndexOverAForeignKeyIsTheUsers asserts that an index over a foreign key's
// columns is diffed like any other, whatever its name: none of them is the foreign key's
// own, so each one written is added and each one left out is dropped, and dropping the
// foreign key leaves it.
func TestDiffIndex_IndexOverAForeignKeyIsTheUsers(t *testing.T) {
	const parent = "CREATE TABLE p (id INT NOT NULL PRIMARY KEY);\n"
	bare := loadCat(t, dbDDL+parent+"CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id));")
	// The names MySQL gives the index it adds (the constraint's, the column's, the
	// column's with a number), and one of the user's.
	for _, name := range []string{"fk", "pid", "pid_2", "my_idx"} {
		t.Run(name, func(t *testing.T) {
			indexed := loadCat(t, dbDDL+parent+"CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY "+name+" (pid), CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id));")
			for _, tc := range []struct {
				from, to *Catalog
				want     DiffAction
			}{{bare, indexed, DiffAdd}, {indexed, bare, DiffDrop}} {
				te := findTable(Diff(tc.from, tc.to), "c")
				if te == nil || len(te.Indexes) != 1 || te.Indexes[0].Action != tc.want || te.Indexes[0].Name != name {
					t.Errorf("want one %v of index %s, got %s", tc.want, name, describeDiff(Diff(tc.from, tc.to)))
				}
				if te != nil && len(te.ForeignKeys) != 0 {
					t.Errorf("the foreign key did not change, got %s", describeDiff(Diff(tc.from, tc.to)))
				}
			}
			// Dropping the foreign key leaves the index.
			alone := loadCat(t, dbDDL+parent+"CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY "+name+" (pid));")
			te := findTable(Diff(indexed, alone), "c")
			if te == nil || len(te.Indexes) != 0 || len(te.ForeignKeys) != 1 || te.ForeignKeys[0].Action != DiffDrop {
				t.Errorf("want the foreign key dropped and no index change, got %s", describeDiff(Diff(indexed, alone)))
			}
		})
	}
}

// mustGetTable returns the named table from a catalog or fails the test.
func mustGetTable(t testing.TB, cat *Catalog, name string) *Table {
	t.Helper()
	for _, db := range cat.Databases() {
		if tt := db.GetTable(name); tt != nil {
			return tt
		}
	}
	t.Fatalf("table %s not loaded", name)
	return nil
}

// TestCanonicalIndex_WhatOceanBaseDoesNotStore asserts what canonicalIndex leaves out because
// OceanBase does not keep it: a column's DESC, the USING type, and an index KEY_BLOCK_SIZE.
// TestOracle_OceanBaseIndexOrderAndType asks the engine.
func TestCanonicalIndex_WhatOceanBaseDoesNotStore(t *testing.T) {
	n := NormalizerFor(MySQL80)
	plain := &Index{Name: "k", Visible: true, Columns: []*IndexColumn{{Name: "a"}, {Name: "b"}}}
	for name, other := range map[string]*Index{
		"DESC":           {Name: "k", Visible: true, Columns: []*IndexColumn{{Name: "a"}, {Name: "b", Descending: true}}},
		"USING BTREE":    {Name: "k", Visible: true, IndexType: "BTREE", Columns: []*IndexColumn{{Name: "a"}, {Name: "b"}}},
		"USING HASH":     {Name: "k", Visible: true, IndexType: "HASH", Columns: []*IndexColumn{{Name: "a"}, {Name: "b"}}},
		"KEY_BLOCK_SIZE": {Name: "k", Visible: true, KeyBlockSize: 4, Columns: []*IndexColumn{{Name: "a"}, {Name: "b"}}},
	} {
		if canonicalIndex(plain, n) != canonicalIndex(other, n) {
			t.Errorf("%s is not stored, and must not be part of the key", name)
		}
	}
	// What it does keep still tells indexes apart.
	for name, other := range map[string]*Index{
		"another column": {Name: "k", Visible: true, Columns: []*IndexColumn{{Name: "a"}, {Name: "c"}}},
		"a prefix":       {Name: "k", Visible: true, Columns: []*IndexColumn{{Name: "a"}, {Name: "b", Length: 4}}},
		"UNIQUE":         {Name: "k", Visible: true, Unique: true, Columns: []*IndexColumn{{Name: "a"}, {Name: "b"}}},
	} {
		if canonicalIndex(plain, n) == canonicalIndex(other, n) {
			t.Errorf("%s must change the key", name)
		}
	}
}

// TestDiffIndex_ScopeChange pins the index-scope matrix. OceanBase prints LOCAL on every
// secondary index whether or not the DDL asked for it, so an unstated scope must compare
// equal to an echoed LOCAL or every table phantom-diffs against its own readback. GLOBAL
// is only ever written deliberately, so a change into or out of it must still diff.
func TestDiffIndex_ScopeChange(t *testing.T) {
	tbl := func(scope string) string {
		return dbDDL + "CREATE TABLE t (id INT, a INT, PRIMARY KEY (id), KEY ka (a) " + scope +
			") PARTITION BY KEY(id) PARTITIONS 2"
	}
	for _, tc := range []struct {
		name      string
		from, to  string
		wantEmpty bool
	}{
		{"unstated against the echoed LOCAL", "", "LOCAL", true},
		{"LOCAL against LOCAL", "LOCAL", "LOCAL", true},
		{"GLOBAL against GLOBAL", "GLOBAL", "GLOBAL", true},
		{"LOCAL to GLOBAL", "LOCAL", "GLOBAL", false},
		{"GLOBAL to LOCAL", "GLOBAL", "LOCAL", false},
		{"unstated to GLOBAL", "", "GLOBAL", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := Diff(loadCat(t, tbl(tc.from)), loadCat(t, tbl(tc.to)))
			if got := d.IsEmpty(); got != tc.wantEmpty {
				t.Errorf("empty = %v, want %v (%s)", got, tc.wantEmpty, describeDiff(d))
			}
		})
	}

	// The scope still has to reach the catalog and SHOW CREATE TABLE, folded or not.
	cat := loadCat(t, tbl("GLOBAL"))
	table := cat.GetDatabase("app").GetTable("t")
	var idx *Index
	for _, i := range table.Indexes {
		if i.Name == "ka" {
			idx = i
		}
	}
	if idx == nil {
		t.Fatal("index ka not found")
	}
	if idx.Scope != "GLOBAL" {
		t.Errorf("Scope = %q, want GLOBAL", idx.Scope)
	}
	if got := cat.ShowCreateTable("app", "t"); !strings.Contains(got, "GLOBAL") {
		t.Errorf("SHOW CREATE TABLE dropped the scope:\n%s", got)
	}
}

// TestDiffIndex_ScopeMigrationRoundTrip asserts the generated migration actually effects a
// scope change. The diff reporting one is not enough: if the ADD INDEX omits the scope, the
// index comes back at the default and the next diff reports the same change forever.
func TestDiffIndex_ScopeMigrationRoundTrip(t *testing.T) {
	tbl := func(scope string) string {
		return "CREATE TABLE t (id INT, a INT, PRIMARY KEY (id), KEY ka (a) " + scope +
			") PARTITION BY KEY(id) PARTITIONS 2"
	}
	for _, tc := range []struct{ name, from, to string }{
		{"LOCAL to GLOBAL", "LOCAL", "GLOBAL"},
		{"unstated to GLOBAL", "", "GLOBAL"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			n := NormalizerFor(MySQL80)
			from := loadCat(t, dbDDL+tbl(tc.from))
			to := loadCat(t, dbDDL+tbl(tc.to))
			plan := GenerateMigrationWithNormalizer(from, to, DiffWithNormalizer(from, to, n), n)
			sql := plan.SQL()
			if !strings.Contains(strings.ToUpper(sql), "GLOBAL") {
				t.Fatalf("the migration does not mention the target scope:\n%s", sql)
			}
			// Applying it must reach the target, or the change repeats on every run.
			applied := loadCat(t, dbDDL+tbl(tc.from))
			mustExec(t, applied, sql)
			if d := DiffWithNormalizer(applied, to, n); !d.IsEmpty() {
				t.Errorf("the migration did not land the scope, so it would repeat: %s\n%s",
					describeDiff(d), sql)
			}
		})
	}
}
