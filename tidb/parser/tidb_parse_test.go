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
