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
	if _, present := diffableIndexMap(child, nil)["fk"]; present {
		t.Errorf("FK-implicit backing index `fk` must be excluded from the diff-able set")
	}
}

// TestDiffIndex_RicherIndexOnFKColumnsKept asserts the boundary of FK-implicit detection: an
// index that covers the FK columns but is STRUCTURALLY RICHER than a synthesized backing index
// stays user-managed and diff-able. Two forms:
//   - a UNIQUE index on the FK columns (MySQL reuses it to back the FK, but it is genuinely the
//     user's unique constraint, not a plain backing index);
//   - a WIDER composite index whose left prefix covers the FK columns (its column count differs
//     from the FK's, so it is not a backing-index shape).
//
// A PLAIN index whose columns exactly equal the FK columns is, by contrast, byte-for-byte what
// ensureFKBackingIndex would synthesize and is treated as FK-implicit (the FK node owns it) —
// proven empty against the real engine in TestOracle_IndexDiffFKImplicit/fk-explicit-covering.
func TestDiffIndex_RicherIndexOnFKColumnsKept(t *testing.T) {
	t.Run("unique-on-fk-columns", func(t *testing.T) {
		sdl := dbDDL + "CREATE TABLE p (id INT NOT NULL PRIMARY KEY);\n" +
			"CREATE TABLE c (id INT PRIMARY KEY, pid INT, UNIQUE KEY uq_pid (pid), CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id));"
		cat := loadCat(t, sdl)
		child := mustGetTable(t, cat, "c")
		if _, present := diffableIndexMap(child, nil)["uq_pid"]; !present {
			t.Errorf("UNIQUE index on FK columns must remain diff-able (not a plain backing-index shape)")
		}
	})

	t.Run("wider-composite-covering-fk", func(t *testing.T) {
		sdl := dbDDL + "CREATE TABLE p (a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (a,b));\n" +
			"CREATE TABLE c (id INT PRIMARY KEY, x INT, y INT, z INT, KEY wide (x,y,z), CONSTRAINT fk FOREIGN KEY (x,y) REFERENCES p(a,b));"
		cat := loadCat(t, sdl)
		child := mustGetTable(t, cat, "c")
		if _, present := diffableIndexMap(child, nil)["wide"]; !present {
			t.Errorf("wider composite index covering the FK by left-prefix must remain diff-able")
		}
	})
}

// mustGetTable returns the named table from a catalog or fails the test.
func mustGetTable(t *testing.T, cat *Catalog, name string) *Table {
	t.Helper()
	for _, db := range cat.Databases() {
		if tt := db.GetTable(name); tt != nil {
			return tt
		}
	}
	t.Fatalf("table %s not loaded", name)
	return nil
}

// TestDiffIndex_FKDroppedExplicitIndexKept asserts the asymmetric-exclusion firewall: when a
// foreign key is dropped but the plain user index that backed it is KEPT, the index must NOT be
// reported as added (it exists on both sides). Without the cross-side exclusion check, the index
// would be excluded on the FK side (from) but present on the no-FK side (to), producing a phantom
// ADD that fails on apply with errno 1061 (duplicate key name).
func TestDiffIndex_FKDroppedExplicitIndexKept(t *testing.T) {
	from := loadCat(t, dbDDL+"CREATE TABLE p (id INT NOT NULL PRIMARY KEY);\n"+
		"CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY my_idx (pid), CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id));")
	to := loadCat(t, dbDDL+"CREATE TABLE p (id INT NOT NULL PRIMARY KEY);\n"+
		"CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY my_idx (pid));")
	d := Diff(from, to)
	if te := findTable(d, "c"); te != nil {
		for _, ie := range te.Indexes {
			t.Errorf("no index change expected on c (my_idx kept), got %s %s", ie.Action, ie.Name)
		}
	}
	// And the reverse (FK added, index already present): must not DROP my_idx.
	dRev := Diff(to, from)
	if te := findTable(dRev, "c"); te != nil {
		for _, ie := range te.Indexes {
			t.Errorf("no index change expected on c (FK added, my_idx already present), got %s %s", ie.Action, ie.Name)
		}
	}
}

// TestDiffIndex_ExplicitlyNamedIndexBackingFKEmitted asserts that a user index named DIFFERENTLY
// from the FK's auto-backing name (constraint name / first column) is NOT treated as FK-implicit:
// the user deliberately named it, so it must be diff-able and emitted (otherwise, when the FK is
// added, the FK node would auto-create a differently-named backing index and the user's chosen
// name would be lost). The pure-implicit FK (no explicit KEY) still emits nothing.
func TestDiffIndex_ExplicitlyNamedIndexBackingFKEmitted(t *testing.T) {
	// User explicitly names the backing index `idx_pid` (≠ constraint `fk_c_p`, ≠ column `pid`).
	from := loadCat(t, dbDDL+"CREATE TABLE p (id INT PRIMARY KEY);\nCREATE TABLE c (pid INT);")
	to := loadCat(t, dbDDL+"CREATE TABLE p (id INT PRIMARY KEY);\n"+
		"CREATE TABLE c (pid INT, KEY idx_pid (pid), CONSTRAINT fk_c_p FOREIGN KEY (pid) REFERENCES p(id));")
	d := Diff(from, to)
	te := findTable(d, "c")
	if te == nil || findIndexDiff(te, "idx_pid") == nil || findIndexDiff(te, "idx_pid").Action != DiffAdd {
		t.Fatalf("expected ADD of user-named index idx_pid, got %s", describeDiff(d))
	}

	// The user index `idx_pid` is diff-able even with the FK present (it is not the auto name).
	toChild := mustGetTable(t, to, "c")
	if _, present := diffableIndexMap(toChild, nil)["idx_pid"]; !present {
		t.Errorf("user-named idx_pid must remain diff-able alongside the FK")
	}

	// Counter-check: an index named like the auto-backing name (the column `pid`) IS excluded.
	toAuto := loadCat(t, dbDDL+"CREATE TABLE p (id INT PRIMARY KEY);\n"+
		"CREATE TABLE c (pid INT, KEY pid (pid), CONSTRAINT fk_c_p FOREIGN KEY (pid) REFERENCES p(id));")
	if _, present := diffableIndexMap(mustGetTable(t, toAuto, "c"), nil)["pid"]; present {
		t.Errorf("index named `pid` (the auto-backing name) must be excluded as FK-implicit")
	}
}

// TestDiffIndex_FKImplicitNameCollision asserts that an unnamed FK whose backing-index name was
// suffixed by allocIndexName on collision (here `pid_2`, because the user wrote `KEY pid (a)`) is
// still detected as FK-implicit and excluded — structural detection does not depend on the
// synthesized name. The unrelated user index `KEY pid (a)` stays diff-able.
func TestDiffIndex_FKImplicitNameCollision(t *testing.T) {
	sdl := dbDDL + "CREATE TABLE p (id INT NOT NULL PRIMARY KEY);\n" +
		"CREATE TABLE c (id INT PRIMARY KEY, a INT, pid INT, KEY pid (a), FOREIGN KEY (pid) REFERENCES p(id));"
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
	// Precondition: the loader synthesized the FK backing index as `pid_2` (collided with the
	// user's `KEY pid (a)`).
	if findIndex(child, "pid_2") == nil {
		t.Fatal("precondition: expected synthesized backing index `pid_2`")
	}
	dm := diffableIndexMap(child, nil)
	if _, present := dm["pid_2"]; present {
		t.Errorf("FK-implicit backing index `pid_2` must be excluded despite the collision-suffixed name")
	}
	if _, present := dm["pid"]; !present {
		t.Errorf("unrelated user index `KEY pid (a)` must remain diff-able")
	}
	if d := Diff(cat, cat); !d.IsEmpty() {
		t.Errorf("self-diff must be empty, got: %s", describeDiff(d))
	}
}

// TestDiffIndex_FKImplicitUserNamedConstraint asserts that an FK whose user-chosen constraint name
// happens to match the auto-generated `<table>_ibfk_N` shape still has its backing index detected
// and excluded — structural detection does not misclassify based on the constraint name.
func TestDiffIndex_FKImplicitUserNamedConstraint(t *testing.T) {
	sdl := dbDDL + "CREATE TABLE p (id INT NOT NULL PRIMARY KEY);\n" +
		"CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT c_ibfk_1 FOREIGN KEY (pid) REFERENCES p(id));"
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
	dm := diffableIndexMap(child, nil)
	if _, present := dm["c_ibfk_1"]; present {
		t.Errorf("FK backing index `c_ibfk_1` must be excluded even though the FK name matches the auto shape")
	}
	if d := Diff(cat, cat); !d.IsEmpty() {
		t.Errorf("self-diff must be empty, got: %s", describeDiff(d))
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
