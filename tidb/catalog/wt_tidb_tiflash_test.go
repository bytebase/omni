package catalog

import "testing"

// Walkthrough tests for Tier 3 additions: LOCATION LABELS on ALTER TABLE
// SET TIFLASH REPLICA, and DB-level SET TIFLASH REPLICA via CREATE/ALTER
// DATABASE option.

func TestWTTiDBTiFlash_TableLocationLabels(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY)")
	wtExec(t, c, "ALTER TABLE t SET TIFLASH REPLICA 3 LOCATION LABELS 'zone', 'rack'")

	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl.TiFlashReplica != 3 {
		t.Errorf("TiFlashReplica = %d, want 3", tbl.TiFlashReplica)
	}
	want := []string{"zone", "rack"}
	if len(tbl.TiFlashLocationLabels) != 2 || tbl.TiFlashLocationLabels[0] != want[0] || tbl.TiFlashLocationLabels[1] != want[1] {
		t.Errorf("TiFlashLocationLabels = %v, want %v", tbl.TiFlashLocationLabels, want)
	}
}

// TestWTTiDBTiFlash_LabelsClearedOnReset exercises the semantics of
// SET TIFLASH REPLICA without a LOCATION LABELS clause after the labels
// were previously set. Upstream re-applies the clause state verbatim
// (labels absent = labels cleared), not a merge. Omni copies that.
func TestWTTiDBTiFlash_LabelsClearedOnReset(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY)")
	wtExec(t, c, "ALTER TABLE t SET TIFLASH REPLICA 3 LOCATION LABELS 'zone', 'rack'")
	wtExec(t, c, "ALTER TABLE t SET TIFLASH REPLICA 0")

	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl.TiFlashReplica != 0 {
		t.Errorf("TiFlashReplica = %d, want 0", tbl.TiFlashReplica)
	}
	if len(tbl.TiFlashLocationLabels) != 0 {
		t.Errorf("labels should be cleared on reset, got %v", tbl.TiFlashLocationLabels)
	}
}

func TestWTTiDBTiFlash_DatabaseReplicaCreate(t *testing.T) {
	c := New()
	if _, err := c.Exec("CREATE DATABASE ddb SET TIFLASH REPLICA 2 LOCATION LABELS 'us-east'", nil); err != nil {
		t.Fatalf("Exec: %v", err)
	}
	db := c.GetDatabase("ddb")
	if db == nil {
		t.Fatal("database ddb not found")
	}
	if db.TiFlashReplica != 2 {
		t.Errorf("TiFlashReplica = %d, want 2", db.TiFlashReplica)
	}
	if len(db.TiFlashLocationLabels) != 1 || db.TiFlashLocationLabels[0] != "us-east" {
		t.Errorf("TiFlashLocationLabels = %v, want [us-east]", db.TiFlashLocationLabels)
	}
}

func TestWTTiDBTiFlash_DatabaseReplicaAlter(t *testing.T) {
	c := New()
	if _, err := c.Exec("CREATE DATABASE ddb", nil); err != nil {
		t.Fatalf("CREATE: %v", err)
	}
	if _, err := c.Exec("ALTER DATABASE ddb SET TIFLASH REPLICA 4", nil); err != nil {
		t.Fatalf("ALTER: %v", err)
	}
	db := c.GetDatabase("ddb")
	if db.TiFlashReplica != 4 {
		t.Errorf("TiFlashReplica = %d, want 4", db.TiFlashReplica)
	}
	if len(db.TiFlashLocationLabels) != 0 {
		t.Errorf("labels = %v, want empty", db.TiFlashLocationLabels)
	}
}

// TestWTTiDBTiFlash_MixedWithPlacementPolicy verifies that a CREATE
// DATABASE statement can carry both PLACEMENT POLICY and SET TIFLASH
// REPLICA options. The upstream DatabaseOption list is one rule shared
// by both; order independence is part of the contract.
func TestWTTiDBTiFlash_MixedWithPlacementPolicy(t *testing.T) {
	c := New()
	if _, err := c.Exec("CREATE PLACEMENT POLICY p1 PRIMARY_REGION = 'us'", nil); err != nil {
		t.Fatalf("CREATE POLICY: %v", err)
	}
	if _, err := c.Exec("CREATE DATABASE ddb PLACEMENT POLICY = p1 SET TIFLASH REPLICA 2 LOCATION LABELS 'a'", nil); err != nil {
		t.Fatalf("CREATE DB: %v", err)
	}
	db := c.GetDatabase("ddb")
	if db.PlacementPolicy != "p1" {
		t.Errorf("PlacementPolicy = %q, want p1", db.PlacementPolicy)
	}
	if db.TiFlashReplica != 2 {
		t.Errorf("TiFlashReplica = %d, want 2", db.TiFlashReplica)
	}
	if len(db.TiFlashLocationLabels) != 1 || db.TiFlashLocationLabels[0] != "a" {
		t.Errorf("labels = %v, want [a]", db.TiFlashLocationLabels)
	}
}
