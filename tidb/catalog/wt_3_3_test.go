package catalog

import "testing"

// --- PRIMARY KEY ---

func TestWalkThrough_3_3_PrimaryKeyCreatesIndex(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT, name VARCHAR(50), PRIMARY KEY (id))")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	found := false
	for _, idx := range tbl.Indexes {
		if idx.Primary {
			found = true
			if idx.Name != "PRIMARY" {
				t.Errorf("expected PK index name 'PRIMARY', got %q", idx.Name)
			}
			if !idx.Unique {
				t.Error("PK index should be Unique=true")
			}
			if len(idx.Columns) != 1 || idx.Columns[0].Name != "id" {
				t.Errorf("PK index columns mismatch: %+v", idx.Columns)
			}
		}
	}
	if !found {
		t.Error("no index with Primary=true found")
	}
}

func TestWalkThrough_3_3_PrimaryKeyColumnsNotNull(t *testing.T) {
	c := wtSetup(t)
	// Do NOT specify NOT NULL — PK should auto-mark columns NOT NULL.
	wtExec(t, c, "CREATE TABLE t (id INT, PRIMARY KEY (id))")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	col := tbl.GetColumn("id")
	if col == nil {
		t.Fatal("column id not found")
	}
	if col.Nullable {
		t.Error("PK column should be NOT NULL automatically")
	}
}

// --- UNIQUE KEY ---

func TestWalkThrough_3_3_UniqueKeyCreatesIndexAndConstraint(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, email VARCHAR(100), UNIQUE KEY uk_email (email))")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}

	// Check index.
	var uqIdx *Index
	for _, idx := range tbl.Indexes {
		if idx.Name == "uk_email" {
			uqIdx = idx
			break
		}
	}
	if uqIdx == nil {
		t.Fatal("unique index uk_email not found")
	}
	if !uqIdx.Unique {
		t.Error("expected Unique=true")
	}
	if len(uqIdx.Columns) != 1 || uqIdx.Columns[0].Name != "email" {
		t.Errorf("unique index columns mismatch: %+v", uqIdx.Columns)
	}

	// Check constraint.
	var uqCon *Constraint
	for _, con := range tbl.Constraints {
		if con.Name == "uk_email" && con.Type == ConUniqueKey {
			uqCon = con
			break
		}
	}
	if uqCon == nil {
		t.Fatal("unique constraint uk_email not found")
	}
	if len(uqCon.Columns) != 1 || uqCon.Columns[0] != "email" {
		t.Errorf("unique constraint columns mismatch: %v", uqCon.Columns)
	}
}

// --- Regular INDEX ---

func TestWalkThrough_3_3_RegularIndex(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(100), INDEX idx_name (name))")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	var idx *Index
	for _, i := range tbl.Indexes {
		if i.Name == "idx_name" {
			idx = i
			break
		}
	}
	if idx == nil {
		t.Fatal("index idx_name not found")
	}
	if idx.Unique {
		t.Error("regular index should not be unique")
	}
	if idx.Primary {
		t.Error("regular index should not be primary")
	}
	if len(idx.Columns) != 1 || idx.Columns[0].Name != "name" {
		t.Errorf("index columns mismatch: %+v", idx.Columns)
	}
}

// --- Multi-column index ---

func TestWalkThrough_3_3_MultiColumnIndexOrder(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b INT, c INT, INDEX idx_abc (a, b, c))")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	var idx *Index
	for _, i := range tbl.Indexes {
		if i.Name == "idx_abc" {
			idx = i
			break
		}
	}
	if idx == nil {
		t.Fatal("index idx_abc not found")
	}
	if len(idx.Columns) != 3 {
		t.Fatalf("expected 3 columns, got %d", len(idx.Columns))
	}
	expected := []string{"a", "b", "c"}
	for i, exp := range expected {
		if idx.Columns[i].Name != exp {
			t.Errorf("column %d: expected %q, got %q", i, exp, idx.Columns[i].Name)
		}
	}
}

// --- FULLTEXT index ---

func TestWalkThrough_3_3_FulltextIndex(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, body TEXT, FULLTEXT INDEX ft_body (body))")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	var idx *Index
	for _, i := range tbl.Indexes {
		if i.Name == "ft_body" {
			idx = i
			break
		}
	}
	if idx == nil {
		t.Fatal("fulltext index ft_body not found")
	}
	if !idx.Fulltext {
		t.Error("expected Fulltext=true")
	}
}

// --- SPATIAL index ---

func TestWalkThrough_3_3_SpatialIndex(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, geo GEOMETRY NOT NULL SRID 0, SPATIAL INDEX sp_geo (geo))")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	var idx *Index
	for _, i := range tbl.Indexes {
		if i.Name == "sp_geo" {
			idx = i
			break
		}
	}
	if idx == nil {
		t.Fatal("spatial index sp_geo not found")
	}
	if !idx.Spatial {
		t.Error("expected Spatial=true")
	}
}

// --- Index COMMENT ---

func TestWalkThrough_3_3_IndexComment(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(100), INDEX idx_name (name) COMMENT 'name lookup')")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	var idx *Index
	for _, i := range tbl.Indexes {
		if i.Name == "idx_name" {
			idx = i
			break
		}
	}
	if idx == nil {
		t.Fatal("index idx_name not found")
	}
	if idx.Comment != "name lookup" {
		t.Errorf("expected comment 'name lookup', got %q", idx.Comment)
	}
}

// --- Index INVISIBLE ---

func TestWalkThrough_3_3_IndexInvisible(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(100), INDEX idx_name (name) INVISIBLE)")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	var idx *Index
	for _, i := range tbl.Indexes {
		if i.Name == "idx_name" {
			idx = i
			break
		}
	}
	if idx == nil {
		t.Fatal("index idx_name not found")
	}
	if idx.Visible {
		t.Error("expected Visible=false for INVISIBLE index")
	}
}

// --- FOREIGN KEY ---

func TestWalkThrough_3_3_ForeignKeyConstraint(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE parent (id INT PRIMARY KEY)")
	wtExec(t, c, "CREATE TABLE child (id INT PRIMARY KEY, pid INT, CONSTRAINT fk_pid FOREIGN KEY (pid) REFERENCES parent(id) ON DELETE CASCADE ON UPDATE SET NULL)")
	tbl := c.GetDatabase("testdb").GetTable("child")
	if tbl == nil {
		t.Fatal("table child not found")
	}

	var fk *Constraint
	for _, con := range tbl.Constraints {
		if con.Type == ConForeignKey && con.Name == "fk_pid" {
			fk = con
			break
		}
	}
	if fk == nil {
		t.Fatal("FK constraint fk_pid not found")
	}
	if fk.RefTable != "parent" {
		t.Errorf("expected RefTable 'parent', got %q", fk.RefTable)
	}
	if len(fk.RefColumns) != 1 || fk.RefColumns[0] != "id" {
		t.Errorf("expected RefColumns [id], got %v", fk.RefColumns)
	}
	if fk.OnDelete != "CASCADE" {
		t.Errorf("expected OnDelete 'CASCADE', got %q", fk.OnDelete)
	}
	if fk.OnUpdate != "SET NULL" {
		t.Errorf("expected OnUpdate 'SET NULL', got %q", fk.OnUpdate)
	}
	if len(fk.Columns) != 1 || fk.Columns[0] != "pid" {
		t.Errorf("expected Columns [pid], got %v", fk.Columns)
	}
}

// --- FOREIGN KEY implicit backing index ---

func TestWalkThrough_3_3_ForeignKeyBackingIndex(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE parent (id INT PRIMARY KEY)")
	wtExec(t, c, "CREATE TABLE child (id INT PRIMARY KEY, pid INT, CONSTRAINT fk_pid FOREIGN KEY (pid) REFERENCES parent(id))")
	tbl := c.GetDatabase("testdb").GetTable("child")
	if tbl == nil {
		t.Fatal("table child not found")
	}

	// FK should create an implicit backing index.
	found := false
	for _, idx := range tbl.Indexes {
		if idx.Primary {
			continue
		}
		for _, col := range idx.Columns {
			if col.Name == "pid" {
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	if !found {
		t.Error("expected implicit backing index for FK on column pid")
	}
}

// --- CHECK constraint ---

func TestWalkThrough_3_3_CheckConstraint(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, age INT, CONSTRAINT chk_age CHECK (age >= 0))")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}

	var chk *Constraint
	for _, con := range tbl.Constraints {
		if con.Type == ConCheck && con.Name == "chk_age" {
			chk = con
			break
		}
	}
	if chk == nil {
		t.Fatal("CHECK constraint chk_age not found")
	}
	if chk.CheckExpr == "" {
		t.Error("CHECK expression should not be empty")
	}
	if chk.NotEnforced {
		t.Error("CHECK should be enforced by default")
	}
}

func TestWalkThrough_3_3_CheckConstraintNotEnforced(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, age INT, CONSTRAINT chk_age CHECK (age >= 0) NOT ENFORCED)")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}

	var chk *Constraint
	for _, con := range tbl.Constraints {
		if con.Type == ConCheck && con.Name == "chk_age" {
			chk = con
			break
		}
	}
	if chk == nil {
		t.Fatal("CHECK constraint chk_age not found")
	}
	if !chk.NotEnforced {
		t.Error("CHECK should be NOT ENFORCED")
	}
}

// --- Named constraints ---

func TestWalkThrough_3_3_NamedConstraints(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT, CONSTRAINT my_pk PRIMARY KEY (id))")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	// PK constraint name in MySQL is always "PRIMARY", regardless of user-specified name.
	// But non-PK constraints should preserve names. Test with UNIQUE.
	c2 := wtSetup(t)
	wtExec(t, c2, "CREATE TABLE t (id INT PRIMARY KEY, email VARCHAR(100), CONSTRAINT my_unique UNIQUE KEY (email))")
	tbl2 := c2.GetDatabase("testdb").GetTable("t")
	if tbl2 == nil {
		t.Fatal("table not found")
	}
	found := false
	for _, con := range tbl2.Constraints {
		if con.Type == ConUniqueKey && con.Name == "my_unique" {
			found = true
			break
		}
	}
	if !found {
		t.Error("named unique constraint 'my_unique' not found")
	}
}

// --- Unnamed constraints auto-generated names ---

func TestWalkThrough_3_3_UnnamedConstraintAutoName(t *testing.T) {
	c := wtSetup(t)
	// FK without explicit name should get auto-generated name like t_ibfk_1.
	wtExec(t, c, "CREATE TABLE parent (id INT PRIMARY KEY)")
	wtExec(t, c, "CREATE TABLE child (id INT PRIMARY KEY, pid INT, FOREIGN KEY (pid) REFERENCES parent(id))")
	tbl := c.GetDatabase("testdb").GetTable("child")
	if tbl == nil {
		t.Fatal("table child not found")
	}
	var fk *Constraint
	for _, con := range tbl.Constraints {
		if con.Type == ConForeignKey {
			fk = con
			break
		}
	}
	if fk == nil {
		t.Fatal("FK constraint not found")
	}
	if fk.Name == "" {
		t.Error("unnamed FK constraint should get auto-generated name")
	}
	// MySQL auto-names FKs as tableName_ibfk_N.
	expected := "child_ibfk_1"
	if fk.Name != expected {
		t.Errorf("expected auto-generated FK name %q, got %q", expected, fk.Name)
	}

	// Also test unnamed CHECK: should get tableName_chk_N.
	c2 := wtSetup(t)
	wtExec(t, c2, "CREATE TABLE t (id INT PRIMARY KEY, age INT, CHECK (age >= 0))")
	tbl2 := c2.GetDatabase("testdb").GetTable("t")
	if tbl2 == nil {
		t.Fatal("table not found")
	}
	var chk *Constraint
	for _, con := range tbl2.Constraints {
		if con.Type == ConCheck {
			chk = con
			break
		}
	}
	if chk == nil {
		t.Fatal("CHECK constraint not found")
	}
	if chk.Name == "" {
		t.Error("unnamed CHECK constraint should get auto-generated name")
	}
	expectedChk := "t_chk_1"
	if chk.Name != expectedChk {
		t.Errorf("expected auto-generated CHECK name %q, got %q", expectedChk, chk.Name)
	}
}

// Bug A (CREATE path): auto-generated FK name counter increments per unnamed
// FK, starting from 0, ignoring user-named FKs.
//
// MySQL reference: sql/sql_table.cc:9252 initializes the counter to 0 for
// create_table_impl; sql/sql_table.cc:5912 generate_fk_name uses ++counter.
// This means user-named FKs do NOT seed the counter during CREATE TABLE.
//
// Example: CREATE TABLE child (a INT, CONSTRAINT child_ibfk_5 FK, b INT, FK)
// Real MySQL: unnamed FK gets "child_ibfk_1" (first auto-named, counter 0 → 1).
// Spot-check confirmed with real MySQL 8.0 container.
//
// NOTE: ALTER TABLE ADD FK uses a different rule (max+1) — see the test below.
func TestBugFix_FKCounterCreateTable(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE parent (id INT PRIMARY KEY)")
	wtExec(t, c, `CREATE TABLE child (
		a INT,
		CONSTRAINT child_ibfk_5 FOREIGN KEY (a) REFERENCES parent(id),
		b INT,
		FOREIGN KEY (b) REFERENCES parent(id)
	)`)
	tbl := c.GetDatabase("testdb").GetTable("child")
	if tbl == nil {
		t.Fatal("table child not found")
	}
	var autoGenName string
	for _, con := range tbl.Constraints {
		if con.Type != ConForeignKey {
			continue
		}
		if con.Name == "child_ibfk_5" {
			continue
		}
		autoGenName = con.Name
		break
	}
	if autoGenName == "" {
		t.Fatal("expected a second (auto-named) FK constraint, found none")
	}
	// Real MySQL 8.0 produces "child_ibfk_1" here (verified by spot-check).
	// The user-named "child_ibfk_5" does NOT seed the counter during CREATE.
	if autoGenName != "child_ibfk_1" {
		t.Errorf("expected child_ibfk_1 (first unnamed FK, ignoring user-named _5), got %s", autoGenName)
	}
}

// Bug A (ALTER path): ALTER TABLE ADD FOREIGN KEY uses max(existing)+1 logic.
// MySQL reference: sql/sql_table.cc:14345 (ALTER TABLE) initializes the
// counter via get_fk_max_generated_name_number(), which scans the existing
// table definition for the max generated-name counter.
func TestBugFix_FKCounterAlterTable(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE parent (id INT PRIMARY KEY)")
	wtExec(t, c, `CREATE TABLE child (
		a INT,
		b INT,
		CONSTRAINT child_ibfk_20 FOREIGN KEY (a) REFERENCES parent(id)
	)`)
	wtExec(t, c, "ALTER TABLE child ADD FOREIGN KEY (b) REFERENCES parent(id)")
	tbl := c.GetDatabase("testdb").GetTable("child")
	if tbl == nil {
		t.Fatal("table child not found")
	}
	var autoGenName string
	for _, con := range tbl.Constraints {
		if con.Type != ConForeignKey {
			continue
		}
		if con.Name == "child_ibfk_20" {
			continue
		}
		autoGenName = con.Name
		break
	}
	if autoGenName != "child_ibfk_21" {
		t.Errorf("expected child_ibfk_21 (max 20 + 1), got %s", autoGenName)
	}
}

// Bug B: TIMESTAMP first column must NOT be auto-promoted under MySQL 8.0
// defaults. In 8.0, explicit_defaults_for_timestamp = ON by default, which
// disables promote_first_timestamp_column() (sql/sql_table.cc:10148). omni
// catalog matches this default — it never promotes. This test locks in the
// absence of a stale TIMESTAMP-promotion bug.
func TestBugFix_TimestampNoAutoPromotion(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (ts TIMESTAMP NOT NULL)")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table t not found")
	}
	col := tbl.GetColumn("ts")
	if col == nil {
		t.Fatal("column ts not found")
	}
	if col.Default != nil {
		t.Errorf("expected no auto-promoted DEFAULT (8.0 default behavior), got Default=%q", *col.Default)
	}
	if col.OnUpdate != "" {
		t.Errorf("expected no auto-promoted ON UPDATE (8.0 default behavior), got OnUpdate=%q", col.OnUpdate)
	}
}

// PS1: CHECK constraint counter (CREATE path) follows the same rule as FK
// counter: it's a local counter starting at 0, incrementing per unnamed
// CHECK, IGNORING user-named _chk_N constraints.
//
// MySQL source: sql/sql_table.cc:19073 declares `uint cc_max_generated_number = 0`
// as a fresh local counter. Uses ++cc_max_generated_number per unnamed CHECK.
// If the generated name collides with a user-named one, MySQL errors with
// ER_CHECK_CONSTRAINT_DUP_NAME at sql/sql_table.cc:19595.
//
// Example: CREATE TABLE t (a INT, CONSTRAINT t_chk_1 CHECK(a>0), b INT, CHECK(b<100))
// Real MySQL: the second unnamed CHECK gets t_chk_1 (counter 0 → 1), which
// collides with user-named t_chk_1 → ER_CHECK_CONSTRAINT_DUP_NAME.
// omni currently does not error on collision (PS7 tracking), but at minimum
// it should assign the correct counter sequence ignoring user-named entries.
func TestBugFix_CheckCounterCreateTable(t *testing.T) {
	c := wtSetup(t)
	// Use user-named t_chk_5 so omni's unnamed-CHECK counter (starting 1)
	// doesn't collide.
	wtExec(t, c, `CREATE TABLE t (
		a INT,
		CONSTRAINT t_chk_5 CHECK (a > 0),
		b INT,
		CHECK (b < 100)
	)`)
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table t not found")
	}
	var autoGenName string
	for _, con := range tbl.Constraints {
		if con.Type != ConCheck {
			continue
		}
		if con.Name == "t_chk_5" {
			continue
		}
		autoGenName = con.Name
		break
	}
	if autoGenName == "" {
		t.Fatal("expected a second (auto-named) CHECK constraint, found none")
	}
	// Real MySQL 8.0 produces t_chk_1 here — the user-named _5 is NOT seeded
	// into the counter during CREATE (verified by source code analysis;
	// sql/sql_table.cc:19073 starts cc_max_generated_number at 0).
	if autoGenName != "t_chk_1" {
		t.Errorf("expected t_chk_1 (first unnamed CHECK, ignoring user-named _5), got %s", autoGenName)
	}
}
