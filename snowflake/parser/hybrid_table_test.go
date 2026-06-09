package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// CREATE HYBRID TABLE — Hybrid flag + inline INDEX element (gap-hybrid-table)
// ---------------------------------------------------------------------------

// The HYBRID modifier sets CreateTableStmt.Hybrid and does not disturb the
// reuse of the standard column-def grammar.
func TestCreateHybridTable_Flag(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE HYBRID TABLE t (id INT)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Hybrid {
		t.Error("expected Hybrid=true")
	}
	if stmt.Transient || stmt.Temporary || stmt.Volatile {
		t.Error("temp modifiers must be false for a plain HYBRID TABLE")
	}
	if len(stmt.Columns) != 1 || stmt.Columns[0].Name.Normalize() != "ID" {
		t.Fatalf("unexpected columns: %+v", stmt.Columns)
	}
}

// A regular (non-HYBRID) CREATE TABLE must be unaffected: Hybrid=false and no
// spurious indexes. Guards against the HYBRID/INDEX additions regressing the
// common path.
func TestCreateHybridTable_NonHybridUnaffected(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT, name VARCHAR(100))")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Hybrid {
		t.Error("expected Hybrid=false for a non-HYBRID CREATE TABLE")
	}
	if len(stmt.Indexes) != 0 {
		t.Errorf("expected no indexes, got %d", len(stmt.Indexes))
	}
	if len(stmt.Columns) != 2 {
		t.Fatalf("columns = %d, want 2", len(stmt.Columns))
	}
}

func TestCreateHybridTable_OrReplace(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE OR REPLACE HYBRID TABLE t (id INT)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Hybrid {
		t.Error("expected Hybrid=true")
	}
	if !stmt.OrReplace {
		t.Error("expected OrReplace=true")
	}
}

func TestCreateHybridTable_IfNotExists(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE HYBRID TABLE IF NOT EXISTS t (id INT)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Hybrid {
		t.Error("expected Hybrid=true")
	}
	if !stmt.IfNotExists {
		t.Error("expected IfNotExists=true")
	}
}

// ---------------------------------------------------------------------------
// Inline INDEX element
// ---------------------------------------------------------------------------

func TestCreateHybridTable_IndexSingleColumn(t *testing.T) {
	stmt, errs := testParseCreateTable(
		"CREATE HYBRID TABLE t (id INT, full_name VARCHAR(255), INDEX index_full_name (full_name))")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.Columns) != 2 {
		t.Fatalf("columns = %d, want 2", len(stmt.Columns))
	}
	if len(stmt.Indexes) != 1 {
		t.Fatalf("indexes = %d, want 1", len(stmt.Indexes))
	}
	idx := stmt.Indexes[0]
	if idx.Name.Normalize() != "INDEX_FULL_NAME" {
		t.Errorf("index name = %q, want INDEX_FULL_NAME", idx.Name.Normalize())
	}
	if len(idx.Columns) != 1 || idx.Columns[0].Normalize() != "FULL_NAME" {
		t.Errorf("index columns = %+v, want [FULL_NAME]", idx.Columns)
	}
}

func TestCreateHybridTable_IndexMultiColumn(t *testing.T) {
	stmt, errs := testParseCreateTable(
		"CREATE HYBRID TABLE t (a INT, b INT, c INT, INDEX idx_abc (a, b, c))")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.Indexes) != 1 {
		t.Fatalf("indexes = %d, want 1", len(stmt.Indexes))
	}
	idx := stmt.Indexes[0]
	if idx.Name.Normalize() != "IDX_ABC" {
		t.Errorf("index name = %q, want IDX_ABC", idx.Name.Normalize())
	}
	got := make([]string, len(idx.Columns))
	for i, c := range idx.Columns {
		got[i] = c.Normalize()
	}
	want := []string{"A", "B", "C"}
	if len(got) != len(want) {
		t.Fatalf("index columns = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index column[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

// Loc on a parsed INDEX element must span "INDEX ... )".
func TestCreateHybridTable_IndexLoc(t *testing.T) {
	src := "CREATE HYBRID TABLE t (a INT, INDEX idx_a (a))"
	stmt, errs := testParseCreateTable(src)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.Indexes) != 1 {
		t.Fatalf("indexes = %d, want 1", len(stmt.Indexes))
	}
	idx := stmt.Indexes[0]
	got := src[idx.Loc.Start:idx.Loc.End]
	want := "INDEX idx_a (a)"
	if got != want {
		t.Errorf("index Loc slice = %q, want %q", got, want)
	}
}

// ---------------------------------------------------------------------------
// Corpus shapes (official/create-hybrid-table)
// ---------------------------------------------------------------------------

// example_01: AUTOINCREMENT PRIMARY KEY, UNIQUE, VARIANT, inline INDEX.
func TestCreateHybridTable_CorpusExample01(t *testing.T) {
	stmt, errs := testParseCreateTable(`CREATE HYBRID TABLE mytable (
  customer_id INT AUTOINCREMENT PRIMARY KEY,
  full_name VARCHAR(255),
  email VARCHAR(255) UNIQUE,
  extended_customer_info VARIANT,
  INDEX index_full_name (full_name)
)`)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Hybrid {
		t.Error("expected Hybrid=true")
	}
	if len(stmt.Columns) != 4 {
		t.Fatalf("columns = %d, want 4", len(stmt.Columns))
	}
	// customer_id INT AUTOINCREMENT PRIMARY KEY
	c0 := stmt.Columns[0]
	if c0.Identity == nil {
		t.Error("customer_id: expected AUTOINCREMENT (Identity != nil)")
	}
	if c0.InlineConstraint == nil || c0.InlineConstraint.Type != ast.ConstrPrimaryKey {
		t.Error("customer_id: expected inline PRIMARY KEY")
	}
	// email VARCHAR(255) UNIQUE
	c2 := stmt.Columns[2]
	if c2.InlineConstraint == nil || c2.InlineConstraint.Type != ast.ConstrUnique {
		t.Error("email: expected inline UNIQUE")
	}
	// extended_customer_info VARIANT
	c3 := stmt.Columns[3]
	if c3.DataType == nil || c3.DataType.Kind != ast.TypeVariant {
		t.Errorf("extended_customer_info: type = %v, want TypeVariant", c3.DataType)
	}
	if len(stmt.Indexes) != 1 || stmt.Indexes[0].Name.Normalize() != "INDEX_FULL_NAME" {
		t.Errorf("expected one INDEX index_full_name, got %+v", stmt.Indexes)
	}
}

// example_07: OR REPLACE + out-of-line CONSTRAINT ... PRIMARY KEY (multi-col).
func TestCreateHybridTable_CorpusExample07(t *testing.T) {
	stmt, errs := testParseCreateTable(`CREATE OR REPLACE HYBRID TABLE ht2pk (
  col1 INTEGER NOT NULL,
  col2 INTEGER NOT NULL,
  col3 VARCHAR,
  CONSTRAINT pkey_1 PRIMARY KEY (col1, col2)
)`)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Hybrid || !stmt.OrReplace {
		t.Error("expected Hybrid=true and OrReplace=true")
	}
	if len(stmt.Columns) != 3 {
		t.Fatalf("columns = %d, want 3", len(stmt.Columns))
	}
	if len(stmt.Constraints) != 1 {
		t.Fatalf("constraints = %d, want 1", len(stmt.Constraints))
	}
	con := stmt.Constraints[0]
	if con.Type != ast.ConstrPrimaryKey || con.Name.Normalize() != "PKEY_1" {
		t.Errorf("constraint = %v/%q, want PRIMARY KEY/PKEY_1", con.Type, con.Name.Normalize())
	}
	if len(con.Columns) != 2 {
		t.Errorf("PK columns = %d, want 2", len(con.Columns))
	}
}

// example_08: out-of-line FOREIGN KEY (col) REFERENCES team(team_id).
func TestCreateHybridTable_CorpusExample08_ForeignKey(t *testing.T) {
	stmt, errs := testParseCreateTable(`CREATE OR REPLACE HYBRID TABLE player
  (player_id INT PRIMARY KEY,
  first_name VARCHAR(40),
  team_id INT,
  FOREIGN KEY (team_id) REFERENCES team(team_id))`)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.Constraints) != 1 {
		t.Fatalf("constraints = %d, want 1", len(stmt.Constraints))
	}
	con := stmt.Constraints[0]
	if con.Type != ast.ConstrForeignKey {
		t.Errorf("constraint type = %v, want FOREIGN KEY", con.Type)
	}
	if con.References == nil || con.References.Table.Normalize() != "TEAM" {
		t.Errorf("FK references = %+v, want team", con.References)
	}
}

// example_13: COLLATE 'de' column + INDEX + (the DESCRIBE TABLE is a separate
// statement, exercised by the corpus test).
func TestCreateHybridTable_CorpusExample13_Collate(t *testing.T) {
	stmt, errs := testParseCreateTable(
		"CREATE OR REPLACE HYBRID TABLE ht1 (c1 INT PRIMARY KEY, c2 VARCHAR(10) COLLATE 'de')")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Columns[1].Collate != "de" {
		t.Errorf("c2 collate = %q, want de", stmt.Columns[1].Collate)
	}
}

// example_15: INDEX element followed by a DEFAULT_DDL_COLLATION table property.
func TestCreateHybridTable_CorpusExample15_IndexPlusTableProperty(t *testing.T) {
	stmt, errs := testParseCreateTable(
		"CREATE OR REPLACE HYBRID TABLE ht2 (c1 INT PRIMARY KEY, c2 VARCHAR(10), INDEX idx_c2 (c2)) DEFAULT_DDL_COLLATION = ''")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.Indexes) != 1 || stmt.Indexes[0].Name.Normalize() != "IDX_C2" {
		t.Errorf("expected INDEX idx_c2, got %+v", stmt.Indexes)
	}
}

// example_16: NUMBER(38,0) NOT NULL COMMENT 'text' + out-of-line PK.
func TestCreateHybridTable_CorpusExample16_NumberCommentPK(t *testing.T) {
	stmt, errs := testParseCreateTable(`CREATE OR REPLACE HYBRID TABLE ht1pk
  (COL1 NUMBER(38,0) NOT NULL COMMENT 'Primary key',
  COL2 NUMBER(38,0) NOT NULL,
  CONSTRAINT PKEY_1 PRIMARY KEY (COL1))`)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	c0 := stmt.Columns[0]
	if c0.DataType == nil || c0.DataType.Kind != ast.TypeNumber {
		t.Errorf("COL1 type = %v, want TypeNumber", c0.DataType)
	}
	if len(c0.DataType.Params) != 2 || c0.DataType.Params[0] != 38 || c0.DataType.Params[1] != 0 {
		t.Errorf("COL1 params = %v, want [38 0]", c0.DataType.Params)
	}
	if !c0.NotNull {
		t.Error("COL1: expected NOT NULL")
	}
	if c0.Comment == nil || *c0.Comment != "Primary key" {
		t.Errorf("COL1 comment = %v, want 'Primary key'", c0.Comment)
	}
}

// ---------------------------------------------------------------------------
// Disambiguation: INDEX is a non-reserved keyword, so a column literally named
// "index" must still parse as a column, not an INDEX element.
// ---------------------------------------------------------------------------

func TestCreateHybridTable_ColumnNamedIndex(t *testing.T) {
	for _, tc := range []struct {
		name string
		sql  string
	}{
		{"plain type", "CREATE TABLE t (index INT)"},
		{"parameterized type", "CREATE TABLE t (index NUMBER(38,0))"},
		{"varchar param", "CREATE TABLE t (index VARCHAR(10), other INT)"},
		{"hybrid table column named index", "CREATE HYBRID TABLE t (index INT, other VARCHAR)"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			stmt, errs := testParseCreateTable(tc.sql)
			if len(errs) > 0 {
				t.Fatalf("unexpected errors: %v", errs)
			}
			if len(stmt.Indexes) != 0 {
				t.Errorf("expected no INDEX elements, got %d", len(stmt.Indexes))
			}
			if len(stmt.Columns) == 0 || stmt.Columns[0].Name.Normalize() != "INDEX" {
				t.Errorf("expected first column named INDEX, got %+v", stmt.Columns)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Negatives
// ---------------------------------------------------------------------------

// INDEX with no name followed by a column list is rejected (it is not a valid
// column def either — a column named "index" cannot be followed by '(').
func TestCreateHybridTable_IndexWithoutNameRejected(t *testing.T) {
	result := ParseBestEffort("CREATE HYBRID TABLE t (a INT, INDEX (a))")
	if len(result.Errors) == 0 {
		t.Fatal("expected a parse error for INDEX (a) with no name, got none")
	}
}

// HYBRID pairs only with TABLE. CREATE HYBRID VIEW must be rejected (HYBRID is
// only consumed when immediately followed by TABLE, so HYBRID is left as an
// unexpected object keyword).
func TestCreateHybridTable_HybridViewRejected(t *testing.T) {
	result := ParseBestEffort("CREATE HYBRID VIEW v AS SELECT 1")
	if len(result.Errors) == 0 {
		t.Fatal("expected a parse error for CREATE HYBRID VIEW, got none")
	}
}

// HYBRID remains usable as an ordinary identifier elsewhere (it is not a
// reserved keyword). A table named HYBRID must parse with Hybrid=false.
func TestCreateHybridTable_HybridAsIdentifier(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE hybrid (id INT)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Hybrid {
		t.Error("expected Hybrid=false for a table named HYBRID")
	}
	if stmt.Name.Normalize() != "HYBRID" {
		t.Errorf("name = %q, want HYBRID", stmt.Name.Normalize())
	}
}
