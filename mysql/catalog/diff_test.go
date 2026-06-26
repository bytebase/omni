package catalog

import (
	"strings"
	"testing"
)

// In-process, hermetic structural tests for diff-core. They assert the SchemaDiff
// content (right DiffAction per object) for representative table/column change forms,
// loading both sides through the omni catalog (LoadSDL/LoadSQL) so the model state is
// authentic. The oracle-backed idempotence + canonicalization-empty proofs live in
// diff_oracle_test.go (correctness-protocol.md gate 1 & 2); these lock in gate 3 (the
// non-empty diffs are structurally correct) and run without a live engine.

// loadCat parses SDL into a catalog or fails the test.
func loadCat(t *testing.T, sql string) *Catalog {
	t.Helper()
	c, err := LoadSDL(sql)
	if err != nil {
		t.Fatalf("LoadSDL failed for %q: %v", sql, err)
	}
	return c
}

// findTable locates a table entry in a SchemaDiff by name (lower-cased).
func findTable(d *SchemaDiff, name string) *TableDiffEntry {
	for i := range d.Tables {
		if strings.EqualFold(d.Tables[i].Name, name) {
			return &d.Tables[i]
		}
	}
	return nil
}

// findColumn locates a column entry within a table diff by name (lower-cased).
func findColumn(e *TableDiffEntry, name string) *ColumnDiffEntry {
	if e == nil {
		return nil
	}
	for i := range e.Columns {
		if strings.EqualFold(e.Columns[i].Name, name) {
			return &e.Columns[i]
		}
	}
	return nil
}

const dbDDL = "CREATE DATABASE app DEFAULT CHARSET=utf8mb4;\nUSE app;\n"

// ---- self-diff empty (the idempotence spine, in-process) -------------------

func TestDiff_SelfEmpty(t *testing.T) {
	schemas := []string{
		dbDDL + "CREATE TABLE t (id INT NOT NULL, name VARCHAR(50), PRIMARY KEY (id));",
		dbDDL + "CREATE TABLE a (id INT NOT NULL PRIMARY KEY); CREATE TABLE b (a_id INT, v DECIMAL(10,2) DEFAULT 0);",
		dbDDL + "CREATE TABLE g (a INT, b INT GENERATED ALWAYS AS (a+1) STORED, c VARCHAR(10) DEFAULT 'x');",
		dbDDL + "CREATE TABLE ts (id INT PRIMARY KEY, created TIMESTAMP DEFAULT CURRENT_TIMESTAMP, n INT DEFAULT 0);",
	}
	for _, s := range schemas {
		c := loadCat(t, s)
		for _, n := range []*Normalizer{NormalizerFor(MySQL80), NormalizerFor(MySQL57)} {
			d := DiffWithNormalizer(c, c, n)
			if !d.IsEmpty() {
				t.Errorf("self-diff not empty (version=%d) for %q:\n%+v", n.Version, s, d.Tables)
			}
		}
		// Parameterless Diff must also be empty.
		if d := Diff(c, c); !d.IsEmpty() {
			t.Errorf("parameterless self-diff not empty for %q:\n%+v", s, d.Tables)
		}
	}
}

func TestDiff_EmptyCatalogs(t *testing.T) {
	if d := Diff(New(), New()); !d.IsEmpty() {
		t.Errorf("two empty catalogs must diff empty, got %+v", d)
	}
}

// ---- table-level change forms ----------------------------------------------

func TestDiff_AddTable(t *testing.T) {
	from := loadCat(t, dbDDL+"CREATE TABLE a (id INT PRIMARY KEY);")
	to := loadCat(t, dbDDL+"CREATE TABLE a (id INT PRIMARY KEY); CREATE TABLE b (id INT PRIMARY KEY);")
	d := Diff(from, to)
	if len(d.Tables) != 1 {
		t.Fatalf("want 1 table diff, got %d: %+v", len(d.Tables), d.Tables)
	}
	e := findTable(d, "b")
	if e == nil || e.Action != DiffAdd {
		t.Fatalf("want ADD table b, got %+v", d.Tables)
	}
	if e.To == nil || e.From != nil {
		t.Errorf("ADD entry must carry To and not From, got From=%v To=%v", e.From, e.To)
	}
}

func TestDiff_DropTable(t *testing.T) {
	from := loadCat(t, dbDDL+"CREATE TABLE a (id INT PRIMARY KEY); CREATE TABLE b (id INT PRIMARY KEY);")
	to := loadCat(t, dbDDL+"CREATE TABLE a (id INT PRIMARY KEY);")
	d := Diff(from, to)
	e := findTable(d, "b")
	if e == nil || e.Action != DiffDrop {
		t.Fatalf("want DROP table b, got %+v", d.Tables)
	}
	if e.From == nil || e.To != nil {
		t.Errorf("DROP entry must carry From and not To, got From=%v To=%v", e.From, e.To)
	}
}

func TestDiff_TableEngineChange(t *testing.T) {
	from := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY) ENGINE=InnoDB;")
	to := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY) ENGINE=MyISAM;")
	d := Diff(from, to)
	e := findTable(d, "t")
	if e == nil || e.Action != DiffModify {
		t.Fatalf("want MODIFY table t on engine change, got %+v", d.Tables)
	}
}

func TestDiff_TableCommentChange(t *testing.T) {
	from := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY) COMMENT='old';")
	to := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY) COMMENT='new';")
	d := Diff(from, to)
	if e := findTable(d, "t"); e == nil || e.Action != DiffModify {
		t.Fatalf("want MODIFY table t on comment change, got %+v", d.Tables)
	}
}

func TestDiff_TableCharsetChange(t *testing.T) {
	from := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY) DEFAULT CHARSET=utf8mb4;")
	to := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY) DEFAULT CHARSET=latin1;")
	d := Diff(from, to)
	if e := findTable(d, "t"); e == nil || e.Action != DiffModify {
		t.Fatalf("want MODIFY table t on charset change, got %+v", d.Tables)
	}
}

// AUTO_INCREMENT=N counter must be ignored (it is a live counter, not schema).
func TestDiff_TableAutoIncrementCounterIgnored(t *testing.T) {
	from := loadCat(t, dbDDL+"CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY) AUTO_INCREMENT=1;")
	to := loadCat(t, dbDDL+"CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY) AUTO_INCREMENT=5000;")
	if d := Diff(from, to); !d.IsEmpty() {
		t.Errorf("AUTO_INCREMENT counter change must not diff, got %+v", d.Tables)
	}
}

// A default ROW_FORMAT must not diff against an unspecified one.
func TestDiff_RowFormatDefaultIgnored(t *testing.T) {
	from := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY);")
	to := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY) ROW_FORMAT=DEFAULT;")
	if d := Diff(from, to); !d.IsEmpty() {
		t.Errorf("default ROW_FORMAT must not diff, got %+v", d.Tables)
	}
}

// A genuine non-default ROW_FORMAT change MUST be detected. Regression for the review
// finding that ROW_FORMAT was not compared at all. Oracle-grounded: SHOW CREATE echoes
// ROW_FORMAT=X iff the user explicitly declared a non-DEFAULT value, so an explicit
// format vs an unspecified side is a real declared-option change (the user adds/removes a
// recorded create-option) — not an environment artifact.
func TestDiff_RowFormatExplicitChange(t *testing.T) {
	from := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY) ROW_FORMAT=COMPRESSED;")
	to := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY) ROW_FORMAT=DYNAMIC;")
	if e := findTable(d_(from, to), "t"); e == nil || e.Action != DiffModify {
		t.Errorf("explicit ROW_FORMAT change (COMPRESSED->DYNAMIC) must diff, got %+v", d_(from, to).Tables)
	}
	// Explicit non-default format vs unspecified IS a real change (declared-option
	// added/removed). The oracle confirms a bare table reads back with no ROW_FORMAT clause.
	bare := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY);")
	if e := findTable(d_(from, bare), "t"); e == nil || e.Action != DiffModify {
		t.Errorf("explicit COMPRESSED vs unspecified must diff (declared option removed), got %+v", d_(from, bare).Tables)
	}
}

// ROW_FORMAT=DEFAULT and an unspecified ROW_FORMAT must NOT diff: the oracle confirms
// MySQL strips ROW_FORMAT=DEFAULT (a DEFAULT-declared table reads back with no clause), so
// IgnoreRowFormat treats both as "" and they canonicalize equal. This is the idempotence
// guard for the ignorable case.
func TestDiff_RowFormatDefaultVsBareEmpty(t *testing.T) {
	withDefault := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY) ROW_FORMAT=DEFAULT;")
	bare := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY);")
	if d := Diff(withDefault, bare); !d.IsEmpty() {
		t.Errorf("ROW_FORMAT=DEFAULT vs unspecified must not diff, got %+v", d.Tables)
	}
	if d := Diff(bare, withDefault); !d.IsEmpty() {
		t.Errorf("unspecified vs ROW_FORMAT=DEFAULT must not diff (reverse), got %+v", d.Tables)
	}
}

// d_ is a Diff shorthand for table-only assertions.
func d_(from, to *Catalog) *SchemaDiff { return Diff(from, to) }

// The generated invisible primary key (my_row_id) is engine-synthesized and must NEVER
// produce a column diff. A schema synced with sql_generate_invisible_primary_key=ON
// carries my_row_id; the same schema declared without it must diff EMPTY. Regression for
// the review finding that VisibleColumns() leaks the GIPK (it is Invisible but not
// Hidden).
func TestDiff_GIPKExcluded(t *testing.T) {
	withGIPK := loadCat(t, "SET sql_generate_invisible_primary_key=ON;\n"+dbDDL+"CREATE TABLE t (a INT, b VARCHAR(10));")
	withoutGIPK := loadCat(t, dbDDL+"CREATE TABLE t (a INT, b VARCHAR(10));")

	// Confirm the GIPK column actually exists on the synced side (else the test is moot).
	var tbl *Table
	for _, db := range withGIPK.Databases() {
		if x := db.GetTable("t"); x != nil {
			tbl = x
		}
	}
	if tbl == nil {
		t.Fatal("table t not found")
	}
	hasGIPK := false
	for _, c := range tbl.Columns {
		if c.GeneratedInvisiblePrimaryKey {
			hasGIPK = true
		}
	}
	if !hasGIPK {
		t.Skip("GIPK not generated in this build; nothing to assert")
	}

	if d := Diff(withGIPK, withoutGIPK); !d.IsEmpty() {
		t.Errorf("GIPK my_row_id must not produce a diff (synced ON vs declared without), got %+v", d.Tables)
	}
	if d := Diff(withoutGIPK, withGIPK); !d.IsEmpty() {
		t.Errorf("GIPK my_row_id must not produce a diff (reverse direction), got %+v", d.Tables)
	}
	// Self-diff of the GIPK table must be empty too.
	if d := Diff(withGIPK, withGIPK); !d.IsEmpty() {
		t.Errorf("self-diff of a GIPK table must be empty, got %+v", d.Tables)
	}
}

// A user-declared INVISIBLE column is REAL schema and must be diffed (added/dropped),
// unlike the engine-synthesized GIPK. This guards against an over-broad exclusion.
func TestDiff_UserInvisibleColumnDiffed(t *testing.T) {
	from := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY, a INT);")
	to := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY, a INT, secret INT INVISIBLE);")
	e := findTable(d_(from, to), "t")
	c := findColumn(e, "secret")
	if c == nil || c.Action != DiffAdd {
		t.Fatalf("user-declared INVISIBLE column must be diffed as ADD, got %+v", e)
	}
}

// ---- column-level change forms ---------------------------------------------

func TestDiff_AddColumn(t *testing.T) {
	from := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY);")
	to := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(20));")
	d := Diff(from, to)
	e := findTable(d, "t")
	if e == nil || e.Action != DiffModify {
		t.Fatalf("want MODIFY table t, got %+v", d.Tables)
	}
	c := findColumn(e, "name")
	if c == nil || c.Action != DiffAdd {
		t.Fatalf("want ADD column name, got %+v", e.Columns)
	}
	if c.To == nil || c.From != nil {
		t.Errorf("ADD column must carry To not From, got From=%v To=%v", c.From, c.To)
	}
}

func TestDiff_DropColumn(t *testing.T) {
	from := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(20));")
	to := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY);")
	d := Diff(from, to)
	e := findTable(d, "t")
	c := findColumn(e, "name")
	if c == nil || c.Action != DiffDrop {
		t.Fatalf("want DROP column name, got %+v", e.Columns)
	}
}

func TestDiff_ModifyColumnType(t *testing.T) {
	from := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY, v INT);")
	to := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY, v BIGINT);")
	d := Diff(from, to)
	c := findColumn(findTable(d, "t"), "v")
	if c == nil || c.Action != DiffModify {
		t.Fatalf("want MODIFY column v (int->bigint), got %+v", d.Tables)
	}
}

func TestDiff_ModifyColumnNullability(t *testing.T) {
	from := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY, v INT NULL);")
	to := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY, v INT NOT NULL);")
	d := Diff(from, to)
	c := findColumn(findTable(d, "t"), "v")
	if c == nil || c.Action != DiffModify {
		t.Fatalf("want MODIFY column v on nullability change, got %+v", d.Tables)
	}
}

func TestDiff_ModifyColumnDefault(t *testing.T) {
	from := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY, v INT DEFAULT 0);")
	to := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY, v INT DEFAULT 1);")
	d := Diff(from, to)
	c := findColumn(findTable(d, "t"), "v")
	if c == nil || c.Action != DiffModify {
		t.Fatalf("want MODIFY column v on default change, got %+v", d.Tables)
	}
}

func TestDiff_ModifyColumnCharset(t *testing.T) {
	from := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY, v VARCHAR(10) CHARACTER SET utf8mb4);")
	to := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY, v VARCHAR(10) CHARACTER SET latin1);")
	d := Diff(from, to)
	c := findColumn(findTable(d, "t"), "v")
	if c == nil || c.Action != DiffModify {
		t.Fatalf("want MODIFY column v on charset change, got %+v", d.Tables)
	}
}

func TestDiff_ModifyColumnComment(t *testing.T) {
	from := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY, v INT COMMENT 'a');")
	to := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY, v INT COMMENT 'b');")
	d := Diff(from, to)
	c := findColumn(findTable(d, "t"), "v")
	if c == nil || c.Action != DiffModify {
		t.Fatalf("want MODIFY column v on comment change, got %+v", d.Tables)
	}
}

func TestDiff_ModifyGeneratedExpr(t *testing.T) {
	from := loadCat(t, dbDDL+"CREATE TABLE t (a INT, g INT GENERATED ALWAYS AS (a+1) STORED);")
	to := loadCat(t, dbDDL+"CREATE TABLE t (a INT, g INT GENERATED ALWAYS AS (a+2) STORED);")
	d := Diff(from, to)
	c := findColumn(findTable(d, "t"), "g")
	if c == nil || c.Action != DiffModify {
		t.Fatalf("want MODIFY column g on generated-expr change, got %+v", d.Tables)
	}
}

func TestDiff_ModifyColumnAutoIncrementAttr(t *testing.T) {
	from := loadCat(t, dbDDL+"CREATE TABLE t (id INT NOT NULL PRIMARY KEY);")
	to := loadCat(t, dbDDL+"CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY);")
	d := Diff(from, to)
	c := findColumn(findTable(d, "t"), "id")
	if c == nil || c.Action != DiffModify {
		t.Fatalf("want MODIFY column id on AUTO_INCREMENT attribute change, got %+v", d.Tables)
	}
}

// ---- canonicalization-driven empties (in-process, no engine) ----------------
// A schema written in a user form must diff EMPTY against the same schema written in
// a form the engine would normalize to — proving column equality routes through
// CanonicalColumn. (The authoritative version is the oracle test loading the real
// SHOW CREATE; these are the fast, engine-free siblings.)

func TestDiff_CanonicalEmpty_BooleanVsTinyint1(t *testing.T) {
	// BOOLEAN canonicalizes to tinyint(1) on both versions.
	from := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY, flag BOOLEAN);")
	to := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY, flag TINYINT(1));")
	for _, n := range []*Normalizer{NormalizerFor(MySQL80), NormalizerFor(MySQL57)} {
		if d := DiffWithNormalizer(from, to, n); !d.IsEmpty() {
			t.Errorf("BOOLEAN vs tinyint(1) must be empty (version=%d), got %+v", n.Version, d.Tables)
		}
	}
}

func TestDiff_CanonicalEmpty_NumericDefaultQuoting(t *testing.T) {
	// DEFAULT 0 and DEFAULT '0' are value-equal.
	from := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY, v INT DEFAULT 0);")
	to := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY, v INT DEFAULT '0');")
	if d := Diff(from, to); !d.IsEmpty() {
		t.Errorf("DEFAULT 0 vs DEFAULT '0' must be empty, got %+v", d.Tables)
	}
}

func TestDiff_CanonicalEmpty_IntDisplayWidth80(t *testing.T) {
	// On 8.0, INT(11) and INT collapse to the same canonical type.
	from := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY, v INT(11));")
	to := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY, v INT);")
	if d := DiffWithNormalizer(from, to, NormalizerFor(MySQL80)); !d.IsEmpty() {
		t.Errorf("INT(11) vs INT must be empty on 8.0, got %+v", d.Tables)
	}
}

// FLAG (normalize-core / loader): on 5.7, a NON-DEFAULT integer display width such as
// INT(5) is invisible to the diff. The shared omni loader (formatColumnType in
// tablecmds.go) bakes MySQL 8.0's width-less stored form into Column.ColumnType at load
// time — it strips non-zerofill int display widths unconditionally — so INT(5) and INT
// both load as ColumnType="int" before normalize-core sees them. normalize-core's
// CanonicalColumnType then re-injects 5.7's *default* width (int(11)) for both, so they
// compare equal. This is NOT a diff-core defect: diff-core routes type comparison
// through CanonicalColumn correctly; the width is discarded upstream in the loader.
//
// Idempotence is unaffected (the release path runs target and current through the same
// loader, so both collapse to int identically — no phantom diff). The only loss is
// detecting a width-ONLY change on 5.7, where display width is cosmetic. To distinguish
// non-default int widths on 5.7, the loader must preserve the user width and
// CanonicalColumnType must honor it — owned by normalize-core/loader, flagged not patched.
//
// This test pins the CURRENT diff-core contract: a default-width int change is a no-op
// on 5.7. A genuine type change (int -> bigint) is still detected on 5.7.
func TestDiff_IntDisplayWidth57_LoaderNormalizesAway(t *testing.T) {
	from := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY, v INT(5));")
	to := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY, v INT(11));")
	if d := DiffWithNormalizer(from, to, NormalizerFor(MySQL57)); !d.IsEmpty() {
		t.Errorf("loader strips int width: INT(5) vs INT(11) currently collapse on 5.7; got %+v", d.Tables)
	}
	// But a real base-type change is detected on 5.7.
	to2 := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY, v BIGINT);")
	if d := DiffWithNormalizer(from, to2, NormalizerFor(MySQL57)); d.IsEmpty() {
		t.Errorf("INT vs BIGINT must diff on 5.7")
	}
}

func TestDiff_CanonicalEmpty_DecimalAlias(t *testing.T) {
	// NUMERIC -> decimal, default precision/scale filled.
	from := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY, v NUMERIC);")
	to := loadCat(t, dbDDL+"CREATE TABLE t (id INT PRIMARY KEY, v DECIMAL(10,0));")
	if d := Diff(from, to); !d.IsEmpty() {
		t.Errorf("NUMERIC vs DECIMAL(10,0) must be empty, got %+v", d.Tables)
	}
}

// ---- hidden / system columns are not diffed --------------------------------

// A functional index creates a hidden system column on the synced side; it must not
// manufacture a phantom column diff against a target that declares only the index.
func TestDiff_FunctionalIndexHiddenColumnIgnored(t *testing.T) {
	ddl := dbDDL + "CREATE TABLE t (id INT PRIMARY KEY, a VARCHAR(20), INDEX idx ((UPPER(a))));"
	c := loadCat(t, ddl)
	// Confirm the hidden column exists, so the test is meaningful.
	var tbl *Table
	for _, db := range c.Databases() {
		if tt := db.GetTable("t"); tt != nil {
			tbl = tt
		}
	}
	if tbl == nil {
		t.Fatal("table t not found")
	}
	if len(tbl.HiddenColumns()) == 0 {
		t.Skip("functional index did not create a hidden column in this build; nothing to assert")
	}
	if d := Diff(c, c); !d.IsEmpty() {
		t.Errorf("self-diff of a table with a functional-index hidden column must be empty, got %+v", d.Tables)
	}
}
