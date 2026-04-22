package catalog

import "testing"

// TiDB-specific walkthrough tests. Each test exercises one TiDB catalog
// feature end-to-end through Exec(), so it catches regressions in both
// parser option handling and catalog option wiring.

func TestWTTiDB_1_1_AutoRandom(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id BIGINT AUTO_RANDOM(5) PRIMARY KEY) SHARD_ROW_ID_BITS = 4")

	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table t not found")
	}
	if tbl.ShardRowIDBits != 4 {
		t.Errorf("expected ShardRowIDBits=4, got %d", tbl.ShardRowIDBits)
	}
	col := tbl.Columns[0]
	if !col.AutoRandom {
		t.Error("expected AutoRandom=true on column id")
	}
	if col.AutoRandomShardBits != 5 {
		t.Errorf("expected AutoRandomShardBits=5, got %d", col.AutoRandomShardBits)
	}
}

func TestWTTiDB_1_2_PlacementPolicy(t *testing.T) {
	c := wtSetup(t)
	// Policy must be defined before a table can reference it — catalog
	// validates refs against the first-class policy map (TiDB error 8237 parity).
	wtExec(t, c, "CREATE PLACEMENT POLICY p1 PRIMARY_REGION = 'us-east'")
	wtExec(t, c, "CREATE TABLE t (id INT) PLACEMENT POLICY = p1")

	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl.PlacementPolicy != "p1" {
		t.Errorf("expected PlacementPolicy=p1, got %q", tbl.PlacementPolicy)
	}
}

func TestWTTiDB_1_3_TTL(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT, created_at DATETIME) TTL = created_at + INTERVAL 1 YEAR TTL_ENABLE = 'ON'")

	tbl := c.GetDatabase("testdb").GetTable("t")
	if !tbl.TTLEnable {
		t.Error("expected TTLEnable=true")
	}
	if tbl.TTLColumn != "created_at" {
		t.Errorf("expected TTLColumn=created_at, got %q", tbl.TTLColumn)
	}
	if tbl.TTLInterval == "" {
		t.Error("expected non-empty TTLInterval")
	}
}

func TestWTTiDB_1_4_TiFlashReplica(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY)")
	wtExec(t, c, "ALTER TABLE t SET TIFLASH REPLICA 2")

	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl.TiFlashReplica != 2 {
		t.Errorf("expected TiFlashReplica=2, got %d", tbl.TiFlashReplica)
	}
}

func TestWTTiDB_1_5_RemoveTTL(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT, created_at DATETIME) TTL = created_at + INTERVAL 1 YEAR TTL_ENABLE = 'ON'")
	wtExec(t, c, "ALTER TABLE t REMOVE TTL")

	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl.TTLEnable {
		t.Error("expected TTLEnable=false after REMOVE TTL")
	}
	if tbl.TTLColumn != "" {
		t.Errorf("expected empty TTLColumn after REMOVE TTL, got %q", tbl.TTLColumn)
	}
}

func TestWTTiDB_1_6_ClusteredPK(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT, PRIMARY KEY (id) CLUSTERED)")

	tbl := c.GetDatabase("testdb").GetTable("t")
	for _, con := range tbl.Constraints {
		if con.Type == ConPrimaryKey {
			if con.Clustered == nil || !*con.Clustered {
				t.Error("expected Clustered=true on table-level PK constraint")
			}
			return
		}
	}
	t.Error("primary key constraint not found")
}

// TestWTTiDB_1_6b_InlineClusteredPK exercises the inline column PK path
// (`id INT PRIMARY KEY CLUSTERED`), where Clustered is set on
// ColumnConstraint and must be hoisted onto the synthesized table-level
// Constraint during createTable.
func TestWTTiDB_1_6b_InlineClusteredPK(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY CLUSTERED)")

	tbl := c.GetDatabase("testdb").GetTable("t")
	for _, con := range tbl.Constraints {
		if con.Type == ConPrimaryKey {
			if con.Clustered == nil || !*con.Clustered {
				t.Error("expected Clustered=true on inline PK constraint (hoisted from ColumnConstraint)")
			}
			return
		}
	}
	t.Error("primary key constraint not found")
}

// TestWTTiDB_1_7_MultiFeatureInteraction combines AUTO_RANDOM + inline
// CLUSTERED PK + PLACEMENT POLICY + TTL + SHARD_ROW_ID_BITS +
// PRE_SPLIT_REGIONS + AUTO_ID_CACHE on one table. Catches bugs where
// option-application order or cross-option state leakage would break a
// single-feature test only by luck.
func TestWTTiDB_1_7_MultiFeatureInteraction(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE PLACEMENT POLICY p1 PRIMARY_REGION = 'us-east'")
	wtExec(t, c, `CREATE TABLE orders (
		id BIGINT AUTO_RANDOM(5) PRIMARY KEY CLUSTERED,
		created_at DATETIME,
		INDEX idx_created (created_at)
	) SHARD_ROW_ID_BITS = 4
	  PRE_SPLIT_REGIONS = 3
	  AUTO_ID_CACHE = 100
	  PLACEMENT POLICY = p1
	  TTL = created_at + INTERVAL 1 YEAR
	  TTL_ENABLE = 'ON'
	  TTL_JOB_INTERVAL = '1h'`)

	tbl := c.GetDatabase("testdb").GetTable("orders")
	if tbl == nil {
		t.Fatal("table orders not found")
	}
	if tbl.ShardRowIDBits != 4 || tbl.PreSplitRegions != 3 || tbl.AutoIDCache != 100 {
		t.Errorf("table options: shard=%d presplit=%d idcache=%d", tbl.ShardRowIDBits, tbl.PreSplitRegions, tbl.AutoIDCache)
	}
	if tbl.PlacementPolicy != "p1" {
		t.Errorf("expected PlacementPolicy=p1, got %q", tbl.PlacementPolicy)
	}
	if tbl.TTLColumn != "created_at" || !tbl.TTLEnable || tbl.TTLJobInterval != "1h" {
		t.Errorf("TTL state: col=%q enable=%v interval=%q", tbl.TTLColumn, tbl.TTLEnable, tbl.TTLJobInterval)
	}
	col := tbl.Columns[0]
	if !col.AutoRandom || col.AutoRandomShardBits != 5 {
		t.Errorf("expected AutoRandom+shard=5, got auto=%v shard=%d", col.AutoRandom, col.AutoRandomShardBits)
	}
	var pkFound bool
	for _, con := range tbl.Constraints {
		if con.Type == ConPrimaryKey {
			pkFound = true
			if con.Clustered == nil || !*con.Clustered {
				t.Error("expected Clustered=true on inline PK")
			}
		}
	}
	if !pkFound {
		t.Error("inline PK not hoisted to table-level constraint")
	}
}

// TestWTTiDB_1_8_TTLDateAdd verifies that DATE_ADD-shaped TTL expressions
// are parsed correctly. This shape was a regression caught in PR2 where
// the lexer stopped at the comma inside the function call.
func TestWTTiDB_1_8_TTLDateAdd(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT, created_at DATETIME) TTL = DATE_ADD(created_at, INTERVAL 1 DAY)")

	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl.TTLColumn != "created_at" {
		t.Errorf("expected TTLColumn=created_at, got %q", tbl.TTLColumn)
	}
	if tbl.TTLInterval == "" {
		t.Error("expected non-empty TTLInterval from DATE_ADD")
	}
}
