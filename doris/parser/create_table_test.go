package parser

import (
	"testing"

	"github.com/bytebase/omni/doris/ast"
)

// helper: parse a single statement and expect no errors, returning the node.
func parseCreateTableStmt(t *testing.T, sql string) *ast.CreateTableStmt {
	t.Helper()
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse(%q) errors: %v", sql, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("Parse(%q): got %d stmts, want 1", sql, len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.CreateTableStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateTableStmt, got %T", file.Stmts[0])
	}
	return stmt
}

// ---------------------------------------------------------------------------
// Basic CREATE TABLE
// ---------------------------------------------------------------------------

func TestCreateTable_Basic(t *testing.T) {
	stmt := parseCreateTableStmt(t, "CREATE TABLE t (id INT, name VARCHAR(50))")
	if stmt.Name.String() != "t" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "t")
	}
	if len(stmt.Columns) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(stmt.Columns))
	}
	if stmt.Columns[0].Name != "id" {
		t.Errorf("Columns[0].Name = %q, want %q", stmt.Columns[0].Name, "id")
	}
	if stmt.Columns[0].Type.Name != "INT" {
		t.Errorf("Columns[0].Type.Name = %q, want %q", stmt.Columns[0].Type.Name, "INT")
	}
	if stmt.Columns[1].Name != "name" {
		t.Errorf("Columns[1].Name = %q, want %q", stmt.Columns[1].Name, "name")
	}
	if stmt.Columns[1].Type.Name != "VARCHAR" {
		t.Errorf("Columns[1].Type.Name = %q, want %q", stmt.Columns[1].Type.Name, "VARCHAR")
	}
	if len(stmt.Columns[1].Type.Params) != 1 || stmt.Columns[1].Type.Params[0] != 50 {
		t.Errorf("Columns[1].Type.Params = %v, want [50]", stmt.Columns[1].Type.Params)
	}
}

func TestCreateTable_WithNotNullDefault(t *testing.T) {
	stmt := parseCreateTableStmt(t, "CREATE TABLE t (id INT NOT NULL DEFAULT 0, name VARCHAR(50) NULL)")
	if len(stmt.Columns) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(stmt.Columns))
	}
	// id: NOT NULL, DEFAULT 0
	col0 := stmt.Columns[0]
	if col0.Nullable == nil || *col0.Nullable != false {
		t.Errorf("Columns[0].Nullable = %v, want false", col0.Nullable)
	}
	if col0.Default == nil {
		t.Fatal("Columns[0].Default is nil, want literal 0")
	}
	lit, ok := col0.Default.(*ast.Literal)
	if !ok {
		t.Fatalf("Columns[0].Default is %T, want *ast.Literal", col0.Default)
	}
	if lit.Value != "0" {
		t.Errorf("Columns[0].Default = %q, want %q", lit.Value, "0")
	}
	// name: NULL
	col1 := stmt.Columns[1]
	if col1.Nullable == nil || *col1.Nullable != true {
		t.Errorf("Columns[1].Nullable = %v, want true", col1.Nullable)
	}
}

func TestCreateTable_WithComment(t *testing.T) {
	stmt := parseCreateTableStmt(t, `CREATE TABLE t (id INT COMMENT 'primary key') COMMENT 'my table'`)
	if stmt.Columns[0].Comment != "primary key" {
		t.Errorf("Columns[0].Comment = %q, want %q", stmt.Columns[0].Comment, "primary key")
	}
	if stmt.Comment != "my table" {
		t.Errorf("Comment = %q, want %q", stmt.Comment, "my table")
	}
}

func TestCreateTable_AggregateKey(t *testing.T) {
	stmt := parseCreateTableStmt(t, "CREATE TABLE t (k1 INT, v1 INT) AGGREGATE KEY(k1)")
	if stmt.KeyDesc == nil {
		t.Fatal("KeyDesc is nil")
	}
	if stmt.KeyDesc.Type != "AGGREGATE" {
		t.Errorf("KeyDesc.Type = %q, want %q", stmt.KeyDesc.Type, "AGGREGATE")
	}
	if len(stmt.KeyDesc.Columns) != 1 || stmt.KeyDesc.Columns[0] != "k1" {
		t.Errorf("KeyDesc.Columns = %v, want [k1]", stmt.KeyDesc.Columns)
	}
}

func TestCreateTable_UniqueKey(t *testing.T) {
	stmt := parseCreateTableStmt(t, "CREATE TABLE t (k1 INT, v1 INT) UNIQUE KEY(k1)")
	if stmt.KeyDesc == nil {
		t.Fatal("KeyDesc is nil")
	}
	if stmt.KeyDesc.Type != "UNIQUE" {
		t.Errorf("KeyDesc.Type = %q, want %q", stmt.KeyDesc.Type, "UNIQUE")
	}
}

func TestCreateTable_DuplicateKey(t *testing.T) {
	stmt := parseCreateTableStmt(t, "CREATE TABLE t (k1 INT, v1 INT) DUPLICATE KEY(k1)")
	if stmt.KeyDesc == nil {
		t.Fatal("KeyDesc is nil")
	}
	if stmt.KeyDesc.Type != "DUPLICATE" {
		t.Errorf("KeyDesc.Type = %q, want %q", stmt.KeyDesc.Type, "DUPLICATE")
	}
}

func TestCreateTable_DistributedByHash(t *testing.T) {
	stmt := parseCreateTableStmt(t, "CREATE TABLE t (id INT) DISTRIBUTED BY HASH(id) BUCKETS 10")
	if stmt.DistributedBy == nil {
		t.Fatal("DistributedBy is nil")
	}
	if stmt.DistributedBy.Type != "HASH" {
		t.Errorf("DistributedBy.Type = %q, want %q", stmt.DistributedBy.Type, "HASH")
	}
	if len(stmt.DistributedBy.Columns) != 1 || stmt.DistributedBy.Columns[0] != "id" {
		t.Errorf("DistributedBy.Columns = %v, want [id]", stmt.DistributedBy.Columns)
	}
	if stmt.DistributedBy.Buckets != 10 {
		t.Errorf("DistributedBy.Buckets = %d, want 10", stmt.DistributedBy.Buckets)
	}
}

func TestCreateTable_DistributedByRandom(t *testing.T) {
	stmt := parseCreateTableStmt(t, "CREATE TABLE t (id INT) DISTRIBUTED BY RANDOM BUCKETS AUTO")
	if stmt.DistributedBy == nil {
		t.Fatal("DistributedBy is nil")
	}
	if stmt.DistributedBy.Type != "RANDOM" {
		t.Errorf("DistributedBy.Type = %q, want %q", stmt.DistributedBy.Type, "RANDOM")
	}
	if !stmt.DistributedBy.Auto {
		t.Error("DistributedBy.Auto should be true")
	}
}

func TestCreateTable_PartitionByRange(t *testing.T) {
	stmt := parseCreateTableStmt(t,
		`CREATE TABLE t (id INT, dt DATE) PARTITION BY RANGE(dt) (PARTITION p1 VALUES LESS THAN ('2024-01-01'))`)
	if stmt.PartitionBy == nil {
		t.Fatal("PartitionBy is nil")
	}
	if stmt.PartitionBy.Type != "RANGE" {
		t.Errorf("PartitionBy.Type = %q, want %q", stmt.PartitionBy.Type, "RANGE")
	}
	if len(stmt.PartitionBy.Columns) != 1 || stmt.PartitionBy.Columns[0] != "dt" {
		t.Errorf("PartitionBy.Columns = %v, want [dt]", stmt.PartitionBy.Columns)
	}
	if len(stmt.PartitionBy.Partitions) != 1 {
		t.Fatalf("expected 1 partition, got %d", len(stmt.PartitionBy.Partitions))
	}
	p := stmt.PartitionBy.Partitions[0]
	if p.Name != "p1" {
		t.Errorf("Partitions[0].Name = %q, want %q", p.Name, "p1")
	}
	if len(p.Values) != 1 || p.Values[0] != "2024-01-01" {
		t.Errorf("Partitions[0].Values = %v, want [2024-01-01]", p.Values)
	}
}

func TestCreateTable_PartitionByList(t *testing.T) {
	stmt := parseCreateTableStmt(t,
		`CREATE TABLE t (id INT, dt DATE NOT NULL) PARTITION BY LIST(dt) (PARTITION p1 VALUES IN (('2020-01-01'), ('2020-01-02')))`)
	if stmt.PartitionBy == nil {
		t.Fatal("PartitionBy is nil")
	}
	if stmt.PartitionBy.Type != "LIST" {
		t.Errorf("PartitionBy.Type = %q, want %q", stmt.PartitionBy.Type, "LIST")
	}
	if len(stmt.PartitionBy.Partitions) != 1 {
		t.Fatalf("expected 1 partition, got %d", len(stmt.PartitionBy.Partitions))
	}
	p := stmt.PartitionBy.Partitions[0]
	if p.Name != "p1" {
		t.Errorf("Partitions[0].Name = %q, want %q", p.Name, "p1")
	}
	if len(p.InValues) != 2 {
		t.Fatalf("expected 2 IN value groups, got %d", len(p.InValues))
	}
}

func TestCreateTable_WithEngine(t *testing.T) {
	stmt := parseCreateTableStmt(t, "CREATE TABLE t (id INT) ENGINE = olap")
	if stmt.Engine != "olap" {
		t.Errorf("Engine = %q, want %q", stmt.Engine, "olap")
	}
}

func TestCreateTable_WithProperties(t *testing.T) {
	stmt := parseCreateTableStmt(t, `CREATE TABLE t (id INT) PROPERTIES("replication_num"="3")`)
	if len(stmt.Properties) != 1 {
		t.Fatalf("expected 1 property, got %d", len(stmt.Properties))
	}
	if stmt.Properties[0].Key != "replication_num" {
		t.Errorf("Properties[0].Key = %q, want %q", stmt.Properties[0].Key, "replication_num")
	}
	if stmt.Properties[0].Value != "3" {
		t.Errorf("Properties[0].Value = %q, want %q", stmt.Properties[0].Value, "3")
	}
}

func TestCreateTable_IfNotExists(t *testing.T) {
	stmt := parseCreateTableStmt(t, "CREATE TABLE IF NOT EXISTS t (id INT)")
	if !stmt.IfNotExists {
		t.Error("IfNotExists should be true")
	}
	if stmt.Name.String() != "t" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "t")
	}
}

func TestCreateTable_External(t *testing.T) {
	stmt := parseCreateTableStmt(t, "CREATE EXTERNAL TABLE t (id INT)")
	if !stmt.External {
		t.Error("External should be true")
	}
}

func TestCreateTable_Temporary(t *testing.T) {
	stmt := parseCreateTableStmt(t, "CREATE TEMPORARY TABLE t (id INT)")
	if !stmt.Temporary {
		t.Error("Temporary should be true")
	}
}

func TestCreateTable_CTAS(t *testing.T) {
	stmt := parseCreateTableStmt(t, `CREATE TABLE t PROPERTIES ('replication_num' = '1') AS SELECT * FROM t1`)
	if stmt.AsSelect == nil {
		t.Fatal("AsSelect is nil")
	}
	if stmt.AsSelect.RawText != "SELECT * FROM t1" {
		t.Errorf("AsSelect.RawText = %q, want %q", stmt.AsSelect.RawText, "SELECT * FROM t1")
	}
	if len(stmt.Properties) != 1 {
		t.Fatalf("expected 1 property, got %d", len(stmt.Properties))
	}
}

func TestCreateTable_Like(t *testing.T) {
	stmt := parseCreateTableStmt(t, "CREATE TABLE t11 LIKE t10")
	if stmt.Like == nil {
		t.Fatal("Like is nil")
	}
	if stmt.Like.String() != "t10" {
		t.Errorf("Like = %q, want %q", stmt.Like.String(), "t10")
	}
	if stmt.Name.String() != "t11" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "t11")
	}
}

func TestCreateTable_QualifiedName(t *testing.T) {
	stmt := parseCreateTableStmt(t, "CREATE TABLE example_db.my_table (id INT)")
	if stmt.Name.String() != "example_db.my_table" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "example_db.my_table")
	}
}

// ---------------------------------------------------------------------------
// Column definition details
// ---------------------------------------------------------------------------

func TestCreateTable_ColumnAggType(t *testing.T) {
	stmt := parseCreateTableStmt(t, "CREATE TABLE t (c1 INT, c2 INT MAX) AGGREGATE KEY(c1)")
	if stmt.Columns[1].AggType != "MAX" {
		t.Errorf("Columns[1].AggType = %q, want %q", stmt.Columns[1].AggType, "MAX")
	}
}

func TestCreateTable_ColumnReplaceAggType(t *testing.T) {
	stmt := parseCreateTableStmt(t, "CREATE TABLE t (k1 TINYINT, v1 CHAR(10) REPLACE) AGGREGATE KEY(k1)")
	if stmt.Columns[1].AggType != "REPLACE" {
		t.Errorf("Columns[1].AggType = %q, want %q", stmt.Columns[1].AggType, "REPLACE")
	}
}

func TestCreateTable_ColumnDefaultString(t *testing.T) {
	stmt := parseCreateTableStmt(t, `CREATE TABLE t (v1 VARCHAR(2048), v2 SMALLINT DEFAULT "10")`)
	if stmt.Columns[1].Default == nil {
		t.Fatal("Columns[1].Default is nil")
	}
	lit := stmt.Columns[1].Default.(*ast.Literal)
	if lit.Value != "10" {
		t.Errorf("Columns[1].Default = %q, want %q", lit.Value, "10")
	}
}

func TestCreateTable_ColumnDefaultTimestamp(t *testing.T) {
	stmt := parseCreateTableStmt(t, `CREATE TABLE t (v1 DATETIME DEFAULT "2014-02-04 15:36:00")`)
	if stmt.Columns[0].Default == nil {
		t.Fatal("Default is nil")
	}
	lit := stmt.Columns[0].Default.(*ast.Literal)
	if lit.Value != "2014-02-04 15:36:00" {
		t.Errorf("Default = %q, want %q", lit.Value, "2014-02-04 15:36:00")
	}
}

func TestCreateTable_ColumnDefaultInt(t *testing.T) {
	stmt := parseCreateTableStmt(t, "CREATE TABLE t (c1 INT, c2 INT DEFAULT 10)")
	if stmt.Columns[1].Default == nil {
		t.Fatal("Columns[1].Default is nil")
	}
	lit := stmt.Columns[1].Default.(*ast.Literal)
	if lit.Value != "10" {
		t.Errorf("Columns[1].Default = %q, want %q", lit.Value, "10")
	}
}

func TestCreateTable_ColumnGenerated(t *testing.T) {
	stmt := parseCreateTableStmt(t, "CREATE TABLE t (c1 INT, c2 INT GENERATED ALWAYS AS (c1 + 1))")
	if stmt.Columns[1].Generated == nil {
		t.Fatal("Columns[1].Generated is nil")
	}
}

func TestCreateTable_EmptyColumnComment(t *testing.T) {
	stmt := parseCreateTableStmt(t, `CREATE TABLE t (id INT COMMENT "")`)
	if stmt.Columns[0].Comment != "" {
		t.Errorf("Columns[0].Comment = %q, want empty string", stmt.Columns[0].Comment)
	}
}

// ---------------------------------------------------------------------------
// Inline INDEX definitions
// ---------------------------------------------------------------------------

func TestCreateTable_InlineIndex(t *testing.T) {
	stmt := parseCreateTableStmt(t,
		`CREATE TABLE t (k1 TINYINT, k2 DECIMAL(10, 2) DEFAULT "10.5", v1 CHAR(10) REPLACE, v2 INT SUM, INDEX k1_idx (k1) USING INVERTED COMMENT 'my first index') AGGREGATE KEY(k1, k2)`)
	if len(stmt.Indexes) != 1 {
		t.Fatalf("expected 1 index, got %d", len(stmt.Indexes))
	}
	idx := stmt.Indexes[0]
	if idx.Name != "k1_idx" {
		t.Errorf("Index.Name = %q, want %q", idx.Name, "k1_idx")
	}
	if len(idx.Columns) != 1 || idx.Columns[0] != "k1" {
		t.Errorf("Index.Columns = %v, want [k1]", idx.Columns)
	}
	if idx.IndexType != "inverted" {
		t.Errorf("Index.IndexType = %q, want %q", idx.IndexType, "inverted")
	}
	if idx.Comment != "my first index" {
		t.Errorf("Index.Comment = %q, want %q", idx.Comment, "my first index")
	}
}

// ---------------------------------------------------------------------------
// Distribution variations
// ---------------------------------------------------------------------------

func TestCreateTable_DistributedByHashMultiCol(t *testing.T) {
	stmt := parseCreateTableStmt(t,
		"CREATE TABLE t (k1 BIGINT, k2 LARGEINT) UNIQUE KEY(k1, k2) DISTRIBUTED BY HASH(k1, k2) BUCKETS 32")
	if stmt.DistributedBy == nil {
		t.Fatal("DistributedBy is nil")
	}
	if len(stmt.DistributedBy.Columns) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(stmt.DistributedBy.Columns))
	}
	if stmt.DistributedBy.Columns[0] != "k1" || stmt.DistributedBy.Columns[1] != "k2" {
		t.Errorf("DistributedBy.Columns = %v, want [k1 k2]", stmt.DistributedBy.Columns)
	}
	if stmt.DistributedBy.Buckets != 32 {
		t.Errorf("DistributedBy.Buckets = %d, want 32", stmt.DistributedBy.Buckets)
	}
}

func TestCreateTable_DistributedByRandomNoBuckets(t *testing.T) {
	stmt := parseCreateTableStmt(t,
		"CREATE TABLE t (c1 INT, c2 INT) DUPLICATE KEY(c1) DISTRIBUTED BY RANDOM")
	if stmt.DistributedBy == nil {
		t.Fatal("DistributedBy is nil")
	}
	if stmt.DistributedBy.Type != "RANDOM" {
		t.Errorf("DistributedBy.Type = %q, want %q", stmt.DistributedBy.Type, "RANDOM")
	}
	if stmt.DistributedBy.Buckets != 0 {
		t.Errorf("DistributedBy.Buckets = %d, want 0", stmt.DistributedBy.Buckets)
	}
	if stmt.DistributedBy.Auto {
		t.Error("DistributedBy.Auto should be false")
	}
}

// ---------------------------------------------------------------------------
// Full combo tests from legacy corpus
// ---------------------------------------------------------------------------

func TestCreateTable_Legacy_t1(t *testing.T) {
	stmt := parseCreateTableStmt(t,
		"CREATE TABLE t1(c1 INT, c2 STRING) DUPLICATE KEY(c1) DISTRIBUTED BY HASH(c1) PROPERTIES ('replication_num' = '1')")
	if stmt.Name.String() != "t1" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "t1")
	}
	if len(stmt.Columns) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(stmt.Columns))
	}
	if stmt.KeyDesc == nil || stmt.KeyDesc.Type != "DUPLICATE" {
		t.Error("expected DUPLICATE KEY")
	}
	if stmt.DistributedBy == nil || stmt.DistributedBy.Type != "HASH" {
		t.Error("expected DISTRIBUTED BY HASH")
	}
	if len(stmt.Properties) != 1 {
		t.Errorf("expected 1 property, got %d", len(stmt.Properties))
	}
}

func TestCreateTable_Legacy_t2(t *testing.T) {
	stmt := parseCreateTableStmt(t,
		"CREATE TABLE t2(c1 INT, c2 INT MAX) AGGREGATE KEY(c1) DISTRIBUTED BY HASH(c1) PROPERTIES ('replication_num' = '1')")
	if stmt.Columns[1].AggType != "MAX" {
		t.Errorf("c2 AggType = %q, want %q", stmt.Columns[1].AggType, "MAX")
	}
	if stmt.KeyDesc.Type != "AGGREGATE" {
		t.Errorf("KeyDesc.Type = %q, want %q", stmt.KeyDesc.Type, "AGGREGATE")
	}
}

func TestCreateTable_Legacy_t3(t *testing.T) {
	stmt := parseCreateTableStmt(t,
		"CREATE TABLE t3(c1 INT, c2 INT) UNIQUE KEY(c1) DISTRIBUTED BY HASH(c1) PROPERTIES ('replication_num' = '1')")
	if stmt.KeyDesc.Type != "UNIQUE" {
		t.Errorf("KeyDesc.Type = %q, want %q", stmt.KeyDesc.Type, "UNIQUE")
	}
}

func TestCreateTable_Legacy_t4_Generated(t *testing.T) {
	stmt := parseCreateTableStmt(t,
		"CREATE TABLE t4(c1 INT, c2 INT GENERATED ALWAYS AS (c1 + 1)) DUPLICATE KEY(c1) DISTRIBUTED BY HASH(c1) PROPERTIES ('replication_num' = '1')")
	if stmt.Columns[1].Generated == nil {
		t.Error("expected generated column")
	}
}

func TestCreateTable_Legacy_t5_Default(t *testing.T) {
	stmt := parseCreateTableStmt(t,
		"CREATE TABLE t5(c1 INT, c2 INT DEFAULT 10) DUPLICATE KEY(c1) DISTRIBUTED BY HASH(c1) PROPERTIES ('replication_num' = '1')")
	if stmt.Columns[1].Default == nil {
		t.Fatal("expected DEFAULT value")
	}
	lit := stmt.Columns[1].Default.(*ast.Literal)
	if lit.Value != "10" {
		t.Errorf("Default = %q, want %q", lit.Value, "10")
	}
}

func TestCreateTable_Legacy_t6_RandomDist(t *testing.T) {
	stmt := parseCreateTableStmt(t,
		"CREATE TABLE t6(c1 INT, c2 INT) DUPLICATE KEY(c1) DISTRIBUTED BY RANDOM PROPERTIES ('replication_num' = '1')")
	if stmt.DistributedBy.Type != "RANDOM" {
		t.Errorf("DistributedBy.Type = %q, want %q", stmt.DistributedBy.Type, "RANDOM")
	}
}

func TestCreateTable_Legacy_t10_CTAS(t *testing.T) {
	stmt := parseCreateTableStmt(t,
		"CREATE TABLE t10 PROPERTIES ('replication_num' = '1') AS SELECT * FROM t1")
	if stmt.AsSelect == nil {
		t.Fatal("AsSelect is nil")
	}
	if stmt.AsSelect.RawText != "SELECT * FROM t1" {
		t.Errorf("AsSelect.RawText = %q, want %q", stmt.AsSelect.RawText, "SELECT * FROM t1")
	}
}

func TestCreateTable_Legacy_t11_Like(t *testing.T) {
	stmt := parseCreateTableStmt(t, "CREATE TABLE t11 LIKE t10")
	if stmt.Like == nil {
		t.Fatal("Like is nil")
	}
	if stmt.Like.String() != "t10" {
		t.Errorf("Like = %q, want %q", stmt.Like.String(), "t10")
	}
}

func TestCreateTable_Legacy_FullCombo(t *testing.T) {
	sql := `CREATE TABLE example_db.table_hash(k1 BIGINT, k2 LARGEINT, v1 VARCHAR(2048), v2 SMALLINT DEFAULT "10") UNIQUE KEY(k1, k2) DISTRIBUTED BY HASH(k1, k2) BUCKETS 32 PROPERTIES("storage_medium" = "SSD", "storage_cooldown_time" = "2015-06-04 00:00:00")`
	stmt := parseCreateTableStmt(t, sql)
	if stmt.Name.String() != "example_db.table_hash" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "example_db.table_hash")
	}
	if len(stmt.Columns) != 4 {
		t.Fatalf("expected 4 columns, got %d", len(stmt.Columns))
	}
	if stmt.KeyDesc.Type != "UNIQUE" {
		t.Errorf("KeyDesc.Type = %q, want %q", stmt.KeyDesc.Type, "UNIQUE")
	}
	if stmt.DistributedBy.Buckets != 32 {
		t.Errorf("Buckets = %d, want 32", stmt.DistributedBy.Buckets)
	}
	if len(stmt.Properties) != 2 {
		t.Fatalf("expected 2 properties, got %d", len(stmt.Properties))
	}
}

func TestCreateTable_Legacy_IfNotExists(t *testing.T) {
	sql := `CREATE TABLE IF NOT EXISTS create_table_use_created_policy (k1 BIGINT, k2 LARGEINT, v1 VARCHAR(2048)) UNIQUE KEY(k1) DISTRIBUTED BY HASH(k1) BUCKETS 3 PROPERTIES("storage_policy" = "test_create_table_use_policy", "replication_num" = "1")`
	stmt := parseCreateTableStmt(t, sql)
	if !stmt.IfNotExists {
		t.Error("IfNotExists should be true")
	}
	if len(stmt.Columns) != 3 {
		t.Fatalf("expected 3 columns, got %d", len(stmt.Columns))
	}
	if stmt.DistributedBy.Buckets != 3 {
		t.Errorf("Buckets = %d, want 3", stmt.DistributedBy.Buckets)
	}
}

func TestCreateTable_Legacy_ColocateGroup(t *testing.T) {
	sql := `CREATE TABLE t1 (id INT COMMENT "", value VARCHAR(8) COMMENT "") DUPLICATE KEY(id) DISTRIBUTED BY HASH(id) BUCKETS 10 PROPERTIES ("colocate_with" = "group1")`
	stmt := parseCreateTableStmt(t, sql)
	if len(stmt.Properties) != 1 {
		t.Fatalf("expected 1 property, got %d", len(stmt.Properties))
	}
	if stmt.Properties[0].Key != "colocate_with" {
		t.Errorf("Properties[0].Key = %q, want %q", stmt.Properties[0].Key, "colocate_with")
	}
}

func TestCreateTable_Legacy_InlineIndexWithAggTypes(t *testing.T) {
	sql := `CREATE TABLE example_db.table_hash(k1 TINYINT, k2 DECIMAL(10, 2) DEFAULT "10.5", v1 CHAR(10) REPLACE, v2 INT SUM, INDEX k1_idx (k1) USING INVERTED COMMENT 'my first index') AGGREGATE KEY(k1, k2) DISTRIBUTED BY HASH(k1) BUCKETS 32 PROPERTIES ("bloom_filter_columns" = "k2")`
	stmt := parseCreateTableStmt(t, sql)
	if len(stmt.Columns) != 4 {
		t.Fatalf("expected 4 columns, got %d", len(stmt.Columns))
	}
	if stmt.Columns[2].AggType != "REPLACE" {
		t.Errorf("v1 AggType = %q, want %q", stmt.Columns[2].AggType, "REPLACE")
	}
	if stmt.Columns[3].AggType != "SUM" {
		t.Errorf("v2 AggType = %q, want %q", stmt.Columns[3].AggType, "SUM")
	}
	if len(stmt.Indexes) != 1 {
		t.Fatalf("expected 1 index, got %d", len(stmt.Indexes))
	}
}

func TestCreateTable_Legacy_NoKeyClause(t *testing.T) {
	sql := `CREATE TABLE example_db.table_hash(k1 TINYINT, k2 DECIMAL(10, 2) DEFAULT "10.5") DISTRIBUTED BY HASH(k1) BUCKETS 32 PROPERTIES ("replication_allocation" = "tag.location.group_a:1, tag.location.group_b:2")`
	stmt := parseCreateTableStmt(t, sql)
	if stmt.KeyDesc != nil {
		t.Error("expected nil KeyDesc")
	}
	if stmt.DistributedBy.Buckets != 32 {
		t.Errorf("Buckets = %d, want 32", stmt.DistributedBy.Buckets)
	}
}

func TestCreateTable_Legacy_DynamicPartition(t *testing.T) {
	sql := `CREATE TABLE example_db.dynamic_partition(k1 DATE, k2 INT, k3 SMALLINT, v1 VARCHAR(2048), v2 DATETIME DEFAULT "2014-02-04 15:36:00") DUPLICATE KEY(k1, k2, k3) PARTITION BY RANGE(k1) () DISTRIBUTED BY HASH(k2) BUCKETS 32 PROPERTIES("dynamic_partition.time_unit" = "DAY", "dynamic_partition.start" = "-3", "dynamic_partition.end" = "3", "dynamic_partition.prefix" = "p", "dynamic_partition.buckets" = "32")`
	stmt := parseCreateTableStmt(t, sql)
	if stmt.PartitionBy == nil {
		t.Fatal("PartitionBy is nil")
	}
	if stmt.PartitionBy.Type != "RANGE" {
		t.Errorf("PartitionBy.Type = %q, want %q", stmt.PartitionBy.Type, "RANGE")
	}
	// Empty partition list is valid for dynamic partition
	if len(stmt.PartitionBy.Partitions) != 0 {
		t.Errorf("expected 0 partitions, got %d", len(stmt.PartitionBy.Partitions))
	}
	if len(stmt.Properties) != 5 {
		t.Fatalf("expected 5 properties, got %d", len(stmt.Properties))
	}
}

func TestCreateTable_Legacy_StepPartition(t *testing.T) {
	sql := `CREATE TABLE t8(c1 INT, c2 DATETIME NOT NULL) DUPLICATE KEY(c1) PARTITION BY RANGE(c2) (FROM ('2020-01-01') TO ('2020-01-10') INTERVAL 1 DAY) DISTRIBUTED BY RANDOM PROPERTIES ('replication_num' = '1')`
	stmt := parseCreateTableStmt(t, sql)
	if stmt.PartitionBy == nil {
		t.Fatal("PartitionBy is nil")
	}
	if len(stmt.PartitionBy.Partitions) != 1 {
		t.Fatalf("expected 1 partition item, got %d", len(stmt.PartitionBy.Partitions))
	}
	p := stmt.PartitionBy.Partitions[0]
	if !p.IsStep {
		t.Error("expected step partition")
	}
	if len(p.FromValues) != 1 || p.FromValues[0] != "2020-01-01" {
		t.Errorf("FromValues = %v, want [2020-01-01]", p.FromValues)
	}
	if len(p.ToValues) != 1 || p.ToValues[0] != "2020-01-10" {
		t.Errorf("ToValues = %v, want [2020-01-10]", p.ToValues)
	}
	if p.Interval != "1" {
		t.Errorf("Interval = %q, want %q", p.Interval, "1")
	}
	if p.IntervalUnit != "DAY" {
		t.Errorf("IntervalUnit = %q, want %q", p.IntervalUnit, "DAY")
	}
}

func TestCreateTable_Legacy_ListPartition(t *testing.T) {
	sql := `CREATE TABLE t9(c1 INT, c2 DATE NOT NULL) DUPLICATE KEY(c1) PARTITION BY LIST(c2) (PARTITION p1 VALUES IN (('2020-01-01'), ('2020-01-02'))) DISTRIBUTED BY RANDOM PROPERTIES ('replication_num' = '1')`
	stmt := parseCreateTableStmt(t, sql)
	if stmt.PartitionBy == nil {
		t.Fatal("PartitionBy is nil")
	}
	if stmt.PartitionBy.Type != "LIST" {
		t.Errorf("PartitionBy.Type = %q, want %q", stmt.PartitionBy.Type, "LIST")
	}
	if len(stmt.PartitionBy.Partitions) != 1 {
		t.Fatalf("expected 1 partition, got %d", len(stmt.PartitionBy.Partitions))
	}
}

// ---------------------------------------------------------------------------
// Auto partition
// ---------------------------------------------------------------------------

func TestCreateTable_AutoPartition(t *testing.T) {
	sql := `CREATE TABLE t7(c1 INT, c2 DATETIME NOT NULL) DUPLICATE KEY(c1) AUTO PARTITION BY RANGE(date_trunc(c2, 'day')) () DISTRIBUTED BY RANDOM PROPERTIES ('replication_num' = '1')`
	stmt := parseCreateTableStmt(t, sql)
	if stmt.PartitionBy == nil {
		t.Fatal("PartitionBy is nil")
	}
	if !stmt.PartitionBy.Auto {
		t.Error("PartitionBy.Auto should be true")
	}
	if stmt.PartitionBy.Type != "RANGE" {
		t.Errorf("PartitionBy.Type = %q, want %q", stmt.PartitionBy.Type, "RANGE")
	}
	// Should have captured the function expression
	if len(stmt.PartitionBy.FuncExprs) != 1 {
		t.Fatalf("expected 1 funcExpr, got %d (cols=%v)", len(stmt.PartitionBy.FuncExprs), stmt.PartitionBy.Columns)
	}
}

// ---------------------------------------------------------------------------
// NodeTag
// ---------------------------------------------------------------------------

func TestCreateTableStmt_Tag(t *testing.T) {
	stmt := &ast.CreateTableStmt{}
	if stmt.Tag() != ast.T_CreateTableStmt {
		t.Errorf("Tag() = %v, want T_CreateTableStmt", stmt.Tag())
	}
}

func TestColumnDef_Tag(t *testing.T) {
	col := &ast.ColumnDef{}
	if col.Tag() != ast.T_ColumnDef {
		t.Errorf("Tag() = %v, want T_ColumnDef", col.Tag())
	}
}

func TestKeyDesc_Tag(t *testing.T) {
	kd := &ast.KeyDesc{}
	if kd.Tag() != ast.T_KeyDesc {
		t.Errorf("Tag() = %v, want T_KeyDesc", kd.Tag())
	}
}

func TestDistributionDesc_Tag(t *testing.T) {
	dd := &ast.DistributionDesc{}
	if dd.Tag() != ast.T_DistributionDesc {
		t.Errorf("Tag() = %v, want T_DistributionDesc", dd.Tag())
	}
}

func TestPartitionDesc_Tag(t *testing.T) {
	pd := &ast.PartitionDesc{}
	if pd.Tag() != ast.T_PartitionDesc {
		t.Errorf("Tag() = %v, want T_PartitionDesc", pd.Tag())
	}
}

// ---------------------------------------------------------------------------
// Loc tracking
// ---------------------------------------------------------------------------

func TestCreateTable_LocIsValid(t *testing.T) {
	stmt := parseCreateTableStmt(t, "CREATE TABLE t (id INT)")
	if !stmt.Loc.IsValid() {
		t.Errorf("Loc = %v, want valid", stmt.Loc)
	}
	if stmt.Loc.Start != 0 {
		t.Errorf("Loc.Start = %d, want 0", stmt.Loc.Start)
	}
}

// ---------------------------------------------------------------------------
// Walk
// ---------------------------------------------------------------------------

func TestCreateTable_Walk(t *testing.T) {
	file, errs := Parse("CREATE TABLE t (id INT, name VARCHAR(50)) DISTRIBUTED BY HASH(id) BUCKETS 10")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}

	var visited []ast.NodeTag
	ast.Inspect(file, func(n ast.Node) bool {
		if n != nil {
			visited = append(visited, n.Tag())
		}
		return true
	})

	// Should visit: File, CreateTableStmt, ObjectName(t), ColumnDef(id), TypeName(INT),
	// ColumnDef(name), TypeName(VARCHAR), DistributionDesc
	if len(visited) < 5 {
		t.Errorf("expected at least 5 visited nodes, got %d: %v", len(visited), visited)
	}

	// First should be File
	if visited[0] != ast.T_File {
		t.Errorf("visited[0] = %v, want File", visited[0])
	}
	// Second should be CreateTableStmt
	if visited[1] != ast.T_CreateTableStmt {
		t.Errorf("visited[1] = %v, want CreateTableStmt", visited[1])
	}
}

// ---------------------------------------------------------------------------
// Full legacy corpus parse test
// ---------------------------------------------------------------------------

func TestCreateTable_LegacyCorpus(t *testing.T) {
	// All statements from the legacy test corpus should parse without errors.
	// Some may have features we partially handle; the test just verifies no panics or errors.
	corpus := []string{
		"CREATE TABLE t1(c1 INT, c2 STRING) DUPLICATE KEY(c1) DISTRIBUTED BY HASH(c1) PROPERTIES ('replication_num' = '1')",
		"CREATE TABLE t2(c1 INT, c2 INT MAX) AGGREGATE KEY(c1) DISTRIBUTED BY HASH(c1) PROPERTIES ('replication_num' = '1')",
		"CREATE TABLE t3(c1 INT, c2 INT) UNIQUE KEY(c1) DISTRIBUTED BY HASH(c1) PROPERTIES ('replication_num' = '1')",
		"CREATE TABLE t4(c1 INT, c2 INT GENERATED ALWAYS AS (c1 + 1)) DUPLICATE KEY(c1) DISTRIBUTED BY HASH(c1) PROPERTIES ('replication_num' = '1')",
		"CREATE TABLE t5(c1 INT, c2 INT DEFAULT 10) DUPLICATE KEY(c1) DISTRIBUTED BY HASH(c1) PROPERTIES ('replication_num' = '1')",
		"CREATE TABLE t6(c1 INT, c2 INT) DUPLICATE KEY(c1) DISTRIBUTED BY RANDOM PROPERTIES ('replication_num' = '1')",
		`CREATE TABLE t7(c1 INT, c2 DATETIME NOT NULL) DUPLICATE KEY(c1) AUTO PARTITION BY RANGE(date_trunc(c2, 'day')) () DISTRIBUTED BY RANDOM PROPERTIES ('replication_num' = '1')`,
		`CREATE TABLE t8(c1 INT, c2 DATETIME NOT NULL) DUPLICATE KEY(c1) PARTITION BY RANGE(c2) (FROM ('2020-01-01') TO ('2020-01-10') INTERVAL 1 DAY) DISTRIBUTED BY RANDOM PROPERTIES ('replication_num' = '1')`,
		`CREATE TABLE t9(c1 INT, c2 DATE NOT NULL) DUPLICATE KEY(c1) PARTITION BY LIST(c2) (PARTITION p1 VALUES IN (('2020-01-01'), ('2020-01-02'))) DISTRIBUTED BY RANDOM PROPERTIES ('replication_num' = '1')`,
		`CREATE TABLE example_db.table_hash(k1 BIGINT, k2 LARGEINT, v1 VARCHAR(2048), v2 SMALLINT DEFAULT "10") UNIQUE KEY(k1, k2) DISTRIBUTED BY HASH(k1, k2) BUCKETS 32 PROPERTIES("storage_medium" = "SSD", "storage_cooldown_time" = "2015-06-04 00:00:00")`,
		`CREATE TABLE IF NOT EXISTS create_table_use_created_policy (k1 BIGINT, k2 LARGEINT, v1 VARCHAR(2048)) UNIQUE KEY(k1) DISTRIBUTED BY HASH(k1) BUCKETS 3 PROPERTIES("storage_policy" = "test_create_table_use_policy", "replication_num" = "1")`,
		`CREATE TABLE t1 (id INT COMMENT "", value VARCHAR(8) COMMENT "") DUPLICATE KEY(id) DISTRIBUTED BY HASH(id) BUCKETS 10 PROPERTIES ("colocate_with" = "group1")`,
		`CREATE TABLE t2 (id INT COMMENT "", value1 VARCHAR(8) COMMENT "", value2 VARCHAR(8) COMMENT "") DUPLICATE KEY(id) DISTRIBUTED BY HASH(id) BUCKETS 10 PROPERTIES ("colocate_with" = "group1")`,
		`CREATE TABLE example_db.table_hash(k1 TINYINT, k2 DECIMAL(10, 2) DEFAULT "10.5", v1 CHAR(10) REPLACE, v2 INT SUM, INDEX k1_idx (k1) USING INVERTED COMMENT 'my first index') AGGREGATE KEY(k1, k2) DISTRIBUTED BY HASH(k1) BUCKETS 32 PROPERTIES ("bloom_filter_columns" = "k2")`,
		`CREATE TABLE example_db.table_hash(k1 TINYINT, k2 DECIMAL(10, 2) DEFAULT "10.5") DISTRIBUTED BY HASH(k1) BUCKETS 32 PROPERTIES ("replication_allocation" = "tag.location.group_a:1, tag.location.group_b:2")`,
		`CREATE TABLE example_db.dynamic_partition(k1 DATE, k2 INT, k3 SMALLINT, v1 VARCHAR(2048), v2 DATETIME DEFAULT "2014-02-04 15:36:00") DUPLICATE KEY(k1, k2, k3) PARTITION BY RANGE(k1) () DISTRIBUTED BY HASH(k2) BUCKETS 32 PROPERTIES("dynamic_partition.time_unit" = "DAY", "dynamic_partition.start" = "-3", "dynamic_partition.end" = "3", "dynamic_partition.prefix" = "p", "dynamic_partition.buckets" = "32")`,
		"CREATE TABLE t10 PROPERTIES ('replication_num' = '1') AS SELECT * FROM t1",
		"CREATE TABLE t11 LIKE t10",
	}

	for i, sql := range corpus {
		file, errs := Parse(sql)
		if len(errs) != 0 {
			t.Errorf("corpus[%d] Parse errors: %v\nSQL: %s", i, errs, sql)
			continue
		}
		if len(file.Stmts) != 1 {
			t.Errorf("corpus[%d] expected 1 stmt, got %d", i, len(file.Stmts))
			continue
		}
		if _, ok := file.Stmts[0].(*ast.CreateTableStmt); !ok {
			t.Errorf("corpus[%d] expected *ast.CreateTableStmt, got %T", i, file.Stmts[0])
		}
	}
}
