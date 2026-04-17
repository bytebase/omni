package parser

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
)

func TestParseAutoRandom(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantShard int
		wantRange int
	}{
		{"basic", "CREATE TABLE t (id BIGINT AUTO_RANDOM PRIMARY KEY)", 0, 0},
		{"shard bits", "CREATE TABLE t (id BIGINT AUTO_RANDOM(5) PRIMARY KEY)", 5, 0},
		{"shard and range", "CREATE TABLE t (id BIGINT AUTO_RANDOM(5, 64) PRIMARY KEY)", 5, 64},
		{"in tidb comment", "CREATE TABLE t (id BIGINT /*T![auto_rand] AUTO_RANDOM */ PRIMARY KEY)", 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			stmt := list.Items[0].(*nodes.CreateTableStmt)
			col := stmt.Columns[0]
			if !col.AutoRandom {
				t.Error("expected AutoRandom=true")
			}
			if col.AutoRandomShardBits != tt.wantShard {
				t.Errorf("shard bits: got %d, want %d", col.AutoRandomShardBits, tt.wantShard)
			}
			if col.AutoRandomRangeBits != tt.wantRange {
				t.Errorf("range bits: got %d, want %d", col.AutoRandomRangeBits, tt.wantRange)
			}
		})
	}
}

func TestAutoRandomOutfuncs(t *testing.T) {
	list, err := Parse("CREATE TABLE t (id BIGINT AUTO_RANDOM(5) PRIMARY KEY)")
	if err != nil {
		t.Fatal(err)
	}
	out := nodes.NodeToString(list.Items[0])
	if !strings.Contains(out, "auto_random") {
		t.Errorf("outfuncs missing auto_random: %s", out)
	}
}

func TestParseTiDBTableOptions(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		optName  string
		optValue string // empty = just verify parse success
	}{
		{"shard_row_id_bits", "CREATE TABLE t (id INT) SHARD_ROW_ID_BITS = 4", "SHARD_ROW_ID_BITS", "4"},
		{"pre_split_regions", "CREATE TABLE t (id INT) PRE_SPLIT_REGIONS = 3", "PRE_SPLIT_REGIONS", "3"},
		{"auto_id_cache", "CREATE TABLE t (id INT) AUTO_ID_CACHE = 100", "AUTO_ID_CACHE", "100"},
		{"auto_random_base", "CREATE TABLE t (id BIGINT AUTO_RANDOM PRIMARY KEY) AUTO_RANDOM_BASE = 1000", "AUTO_RANDOM_BASE", "1000"},
		{"ttl_enable", "CREATE TABLE t (id INT) TTL_ENABLE = 'OFF'", "TTL_ENABLE", ""},
		{"ttl_job_interval", "CREATE TABLE t (id INT) TTL_JOB_INTERVAL = '1h'", "TTL_JOB_INTERVAL", ""},
		{"placement_policy", "CREATE TABLE t (id INT) PLACEMENT POLICY = p1", "PLACEMENT POLICY", "p1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			stmt := list.Items[0].(*nodes.CreateTableStmt)
			found := false
			for _, opt := range stmt.Options {
				if strings.EqualFold(opt.Name, tt.optName) {
					found = true
					if tt.optValue != "" && opt.Value != tt.optValue {
						t.Errorf("option %s: got %q, want %q", tt.optName, opt.Value, tt.optValue)
					}
				}
			}
			if !found {
				t.Errorf("option %s not found in parsed options", tt.optName)
			}
		})
	}
}

func TestParseTTLExpression(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		wantVal  string // substring expected in TTL Value
	}{
		{
			"simple interval",
			"CREATE TABLE t (id INT, created_at DATETIME) TTL = created_at + INTERVAL 1 YEAR",
			"created_at + INTERVAL 1 YEAR",
		},
		{
			"date_add function",
			"CREATE TABLE t (id INT, created_at DATETIME) TTL = DATE_ADD(created_at, INTERVAL 1 DAY)",
			"DATE_ADD(created_at, INTERVAL 1 DAY)",
		},
		{
			"with other options after",
			"CREATE TABLE t (id INT, created_at DATETIME) TTL = created_at + INTERVAL 1 MONTH TTL_ENABLE = 'ON'",
			"created_at + INTERVAL 1 MONTH",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			stmt := list.Items[0].(*nodes.CreateTableStmt)
			found := false
			for _, opt := range stmt.Options {
				if opt.Name == "TTL" {
					found = true
					if !strings.Contains(opt.Value, tt.wantVal) {
						t.Errorf("TTL Value = %q, want substring %q", opt.Value, tt.wantVal)
					}
				}
			}
			if !found {
				t.Error("TTL option not found")
			}
		})
	}
}

func TestParseClusteredPK(t *testing.T) {
	trueVal := true
	falseVal := false
	tests := []struct {
		name      string
		sql       string
		clustered *bool
	}{
		{"clustered", "CREATE TABLE t (id INT, PRIMARY KEY (id) CLUSTERED)", &trueVal},
		{"nonclustered", "CREATE TABLE t (id INT, PRIMARY KEY (id) NONCLUSTERED)", &falseVal},
		{"default", "CREATE TABLE t (id INT, PRIMARY KEY (id))", nil},
		{"inline clustered", "CREATE TABLE t (id INT PRIMARY KEY CLUSTERED)", &trueVal},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			stmt := list.Items[0].(*nodes.CreateTableStmt)
			// Check table-level constraints
			for _, c := range stmt.Constraints {
				if c.Type == nodes.ConstrPrimaryKey {
					if tt.clustered == nil && c.Clustered != nil {
						t.Errorf("expected nil Clustered, got %v", *c.Clustered)
					}
					if tt.clustered != nil && (c.Clustered == nil || *c.Clustered != *tt.clustered) {
						t.Errorf("Clustered mismatch")
					}
					return
				}
			}
			// Check column-level constraints (for inline PK)
			for _, col := range stmt.Columns {
				for _, cc := range col.Constraints {
					if cc.Type == nodes.ColConstrPrimaryKey {
						if tt.clustered == nil && cc.Clustered != nil {
							t.Errorf("expected nil Clustered, got %v", *cc.Clustered)
						}
						if tt.clustered != nil && (cc.Clustered == nil || *cc.Clustered != *tt.clustered) {
							t.Errorf("Clustered mismatch")
						}
						return
					}
				}
			}
			t.Error("no PRIMARY KEY found")
		})
	}
}

func TestParseAlterTiFlashReplica(t *testing.T) {
	tests := []struct {
		name  string
		sql   string
		count int
	}{
		{"set 2", "ALTER TABLE t SET TIFLASH REPLICA 2", 2},
		{"set 0", "ALTER TABLE t SET TIFLASH REPLICA 0", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			stmt := list.Items[0].(*nodes.AlterTableStmt)
			if len(stmt.Commands) == 0 {
				t.Fatal("no commands")
			}
			cmd := stmt.Commands[0]
			if cmd.Type != nodes.ATSetTiFlashReplica {
				t.Errorf("expected ATSetTiFlashReplica, got %v", cmd.Type)
			}
			if cmd.TiFlashReplica != tt.count {
				t.Errorf("replica count: got %d, want %d", cmd.TiFlashReplica, tt.count)
			}
		})
	}
}

func TestParseAlterTiFlashReplicaLocationLabelsDeferred(t *testing.T) {
	// Pinned negative test: LOCATION LABELS is a valid TiDB syntax extension
	// that is not yet supported. This test documents the gap and will fail
	// (in a good way) once support is added.
	sql := "ALTER TABLE t SET TIFLASH REPLICA 2 LOCATION LABELS 'zone', 'rack'"
	_, err := Parse(sql)
	if err == nil {
		t.Fatal("expected parse error for unsupported LOCATION LABELS syntax — if this passes, the feature was added; update the test")
	}
}

func TestParseAlterRemoveTTL(t *testing.T) {
	sql := "ALTER TABLE t REMOVE TTL"
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	stmt := list.Items[0].(*nodes.AlterTableStmt)
	if stmt.Commands[0].Type != nodes.ATRemoveTTL {
		t.Errorf("expected ATRemoveTTL, got %v", stmt.Commands[0].Type)
	}
}

func TestParseDatabasePlacementPolicy(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{"create", "CREATE DATABASE d PLACEMENT POLICY = p1"},
		{"alter", "ALTER DATABASE d PLACEMENT POLICY = p1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
		})
	}
}
