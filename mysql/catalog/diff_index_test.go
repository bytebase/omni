package catalog

import (
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

// TestDiffIndex_FKImplicitExcluded asserts the headline FK-implicit rule structurally, without a
// live engine: a table whose only index is the one synthesized to back a foreign key reports NO
// index diff against the same table, and the index is absent from the diff-able set. (The loader
// synthesizes the backing index via ensureFKBackingIndex, so it is present in tbl.Indexes — the
// differ must still exclude it.)
func TestDiffIndex_FKImplicitExcluded(t *testing.T) {
	sdl := dbDDL + "CREATE TABLE p (id INT NOT NULL PRIMARY KEY);\n" +
		"CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id));"
	cat := loadCat(t, sdl)

	// Self-diff is empty (no phantom on the FK backing index).
	if d := Diff(cat, cat); !d.IsEmpty() {
		t.Fatalf("self-diff with FK backing index must be empty, got: %s", describeDiff(d))
	}

	// The backing index `fk` exists in the model but is excluded from the diff-able set.
	var child *Table
	for _, db := range cat.Databases() {
		if tt := db.GetTable("c"); tt != nil {
			child = tt
		}
	}
	if child == nil {
		t.Fatal("table c not loaded")
	}
	if findIndex(child, "fk") == nil {
		t.Fatal("precondition: loader should have synthesized backing index `fk`")
	}
	if _, present := diffableIndexMap(child)["fk"]; present {
		t.Errorf("FK-implicit backing index `fk` must be excluded from the diff-able set")
	}
}

// TestDiffIndex_ExplicitIndexOnFKColumnsKept asserts the inverse: a user index that merely covers
// the FK columns but carries a user-only attribute (here a distinct name + it is the explicit
// covering index, so MySQL creates no separate backing index) stays user-managed and IS diff-able.
func TestDiffIndex_ExplicitIndexOnFKColumnsKept(t *testing.T) {
	sdl := dbDDL + "CREATE TABLE p (id INT NOT NULL PRIMARY KEY);\n" +
		"CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY explicit_idx (pid), CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id));"
	cat := loadCat(t, sdl)
	var child *Table
	for _, db := range cat.Databases() {
		if tt := db.GetTable("c"); tt != nil {
			child = tt
		}
	}
	if child == nil {
		t.Fatal("table c not loaded")
	}
	if _, present := diffableIndexMap(child)["explicit_idx"]; !present {
		t.Errorf("explicit user index `explicit_idx` covering FK columns must remain diff-able")
	}
}

// TestCanonicalIndex_VersionFlaggedAttrs asserts the version-divergent normalization decisions in
// canonicalIndex directly: DESC and index KEY_BLOCK_SIZE are significant on 8.0/5.7 respectively
// and dropped on the other, matching the live-engine behavior proven in the oracle tests.
func TestCanonicalIndex_VersionFlaggedAttrs(t *testing.T) {
	n80 := NormalizerFor(MySQL80)
	n57 := NormalizerFor(MySQL57)

	asc := &Index{Name: "k", Visible: true, Columns: []*IndexColumn{{Name: "a"}, {Name: "b"}}}
	desc := &Index{Name: "k", Visible: true, Columns: []*IndexColumn{{Name: "a"}, {Name: "b", Descending: true}}}

	// 8.0: DESC distinguishes; 5.7: DESC is dropped (stored ascending), so the keys collapse.
	if canonicalIndex(asc, n80) == canonicalIndex(desc, n80) {
		t.Errorf("8.0: ASC and DESC index keys must differ")
	}
	if canonicalIndex(asc, n57) != canonicalIndex(desc, n57) {
		t.Errorf("5.7: DESC must be dropped from the key (stored ascending), keys must match")
	}

	plain := &Index{Name: "k", Visible: true, Columns: []*IndexColumn{{Name: "a"}}}
	kbs := &Index{Name: "k", Visible: true, KeyBlockSize: 4, Columns: []*IndexColumn{{Name: "a"}}}

	// 5.7: index KEY_BLOCK_SIZE is echoed → significant; 8.0: absorbed → dropped from the key.
	if canonicalIndex(plain, n57) == canonicalIndex(kbs, n57) {
		t.Errorf("5.7: index KEY_BLOCK_SIZE must contribute to the key")
	}
	if canonicalIndex(plain, n80) != canonicalIndex(kbs, n80) {
		t.Errorf("8.0: index KEY_BLOCK_SIZE must be ignored in the key")
	}
}
