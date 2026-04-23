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

// TestParseAlterTiFlashReplicaLocationLabels covers the LOCATION LABELS
// suffix on ALTER TABLE SET TIFLASH REPLICA. Replaces a pinned negative
// placeholder when Tier 3 added feature support.
// Ref: TiDB v8.5.5 parser.y:2193 + 2176-2183.
func TestParseAlterTiFlashReplicaLocationLabels(t *testing.T) {
	tests := []struct {
		name       string
		sql        string
		wantLabels []string
	}{
		{"single label", "ALTER TABLE t SET TIFLASH REPLICA 2 LOCATION LABELS 'zone'", []string{"zone"}},
		{"multi label", "ALTER TABLE t SET TIFLASH REPLICA 2 LOCATION LABELS 'zone', 'rack'", []string{"zone", "rack"}},
		{"zero replica with labels", "ALTER TABLE t SET TIFLASH REPLICA 0 LOCATION LABELS 'z'", []string{"z"}},
		{"no labels clause", "ALTER TABLE t SET TIFLASH REPLICA 2", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			stmt := list.Items[0].(*nodes.AlterTableStmt)
			cmd := stmt.Commands[0]
			if cmd.Type != nodes.ATSetTiFlashReplica {
				t.Errorf("Type = %v, want ATSetTiFlashReplica", cmd.Type)
			}
			if len(cmd.TiFlashLocationLabels) != len(tt.wantLabels) {
				t.Fatalf("TiFlashLocationLabels len = %d, want %d: %v", len(cmd.TiFlashLocationLabels), len(tt.wantLabels), cmd.TiFlashLocationLabels)
			}
			for i, got := range cmd.TiFlashLocationLabels {
				if got != tt.wantLabels[i] {
					t.Errorf("TiFlashLocationLabels[%d] = %q, want %q", i, got, tt.wantLabels[i])
				}
			}
		})
	}
}

// TestTiDBKeywordsUsableAsIdentifiers pins that the new unreserved
// keywords added for TIFLASH / replication work (LOCATION, LABELS,
// MASTER_LOG_FILE) can still be used as bare identifiers in CREATE
// TABLE column names. Upstream categorizes all three as non-reserved
// (parser.y:7071-7072 for LOCATION/LABELS, and MASTER_LOG_FILE is in
// the UnReservedKeyword block at line 6780+). Miscategorizing them
// as reserved would break legitimate MySQL-compat DDL with `location`
// or `labels` columns — exactly the kind of silent regression #104
// had a precedent for.
func TestTiDBKeywordsUsableAsIdentifiers(t *testing.T) {
	sql := "CREATE TABLE t (location VARCHAR(50), labels VARCHAR(50), master_log_file VARCHAR(50))"
	if _, err := Parse(sql); err != nil {
		t.Fatalf("new unreserved keywords should parse as identifiers; Parse failed: %v", err)
	}
}

// TestParseAlterTiFlashReplicaLocationLabelsNegatives covers inputs the
// upstream grammar rejects. LocationLabelList requires LABELS after
// LOCATION, at least one string literal when the clause is present, no
// trailing comma, and string literals (not identifiers).
func TestParseAlterTiFlashReplicaLocationLabelsNegatives(t *testing.T) {
	negatives := []string{
		"ALTER TABLE t SET TIFLASH REPLICA 2 LOCATION",             // missing LABELS
		"ALTER TABLE t SET TIFLASH REPLICA 2 LOCATION LABELS",      // empty list
		"ALTER TABLE t SET TIFLASH REPLICA 2 LOCATION LABELS 'a',", // trailing comma
		"ALTER TABLE t SET TIFLASH REPLICA 2 LOCATION LABELS zone", // bare ident not string
	}
	for _, sql := range negatives {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err == nil {
				t.Errorf("expected parse error for %q, got nil", sql)
			}
		})
	}
}

// TestParseAlterDatabaseSetTiFlashReplica covers SET TIFLASH REPLICA as
// a DatabaseOption — both ALTER DATABASE and CREATE DATABASE paths.
// Ref: parser.y:4482 DatabaseOption arm.
func TestParseAlterDatabaseSetTiFlashReplica(t *testing.T) {
	tests := []struct {
		name       string
		sql        string
		wantCount  int
		wantLabels []string
	}{
		{"alter no labels", "ALTER DATABASE ddb SET TIFLASH REPLICA 2", 2, nil},
		{"alter with labels", "ALTER DATABASE ddb SET TIFLASH REPLICA 2 LOCATION LABELS 'zone'", 2, []string{"zone"}},
		{"alter zero replica", "ALTER DATABASE ddb SET TIFLASH REPLICA 0", 0, nil},
		{"create with labels", "CREATE DATABASE ddb SET TIFLASH REPLICA 1 LOCATION LABELS 'a', 'b'", 1, []string{"a", "b"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			var opts []*nodes.DatabaseOption
			switch stmt := list.Items[0].(type) {
			case *nodes.AlterDatabaseStmt:
				opts = stmt.Options
			case *nodes.CreateDatabaseStmt:
				opts = stmt.Options
			default:
				t.Fatalf("unexpected stmt type %T", list.Items[0])
			}
			if len(opts) != 1 {
				t.Fatalf("Options len = %d, want 1", len(opts))
			}
			opt := opts[0]
			if opt.Name != "TIFLASH REPLICA" {
				t.Errorf("Name = %q, want TIFLASH REPLICA", opt.Name)
			}
			if opt.TiFlashReplica != tt.wantCount {
				t.Errorf("TiFlashReplica = %d, want %d", opt.TiFlashReplica, tt.wantCount)
			}
			if len(opt.TiFlashLocationLabels) != len(tt.wantLabels) {
				t.Fatalf("labels len = %d, want %d: %v", len(opt.TiFlashLocationLabels), len(tt.wantLabels), opt.TiFlashLocationLabels)
			}
			for i, got := range opt.TiFlashLocationLabels {
				if got != tt.wantLabels[i] {
					t.Errorf("label[%d] = %q, want %q", i, got, tt.wantLabels[i])
				}
			}
		})
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
