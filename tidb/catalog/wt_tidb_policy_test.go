package catalog

import "testing"

// Walkthrough tests for the TiDB PLACEMENT POLICY DDL added in Tier 1.
// Each test goes through Exec() to exercise the full parser→dispatcher
// →catalog path.

func TestWTTiDBPolicy_1_Create(t *testing.T) {
	c := New()
	if _, err := c.Exec("CREATE PLACEMENT POLICY p1 PRIMARY_REGION = 'us-east' REGIONS = 'us-east,us-west'", nil); err != nil {
		t.Fatalf("Exec: %v", err)
	}
	p := c.GetPlacementPolicy("p1")
	if p == nil {
		t.Fatal("policy p1 not found after CREATE")
	}
	if p.Name != "p1" {
		t.Errorf("Name = %q", p.Name)
	}
	if len(p.Options) != 2 {
		t.Fatalf("want 2 options, got %d", len(p.Options))
	}
}

func TestWTTiDBPolicy_2_CreateDuplicateErrors(t *testing.T) {
	c := New()
	if _, err := c.Exec("CREATE PLACEMENT POLICY p1 FOLLOWERS = 2", nil); err != nil {
		t.Fatalf("first CREATE: %v", err)
	}
	results, err := c.Exec("CREATE PLACEMENT POLICY p1 FOLLOWERS = 3", nil)
	if err != nil {
		t.Fatalf("parse error on second CREATE: %v", err)
	}
	if results[0].Error == nil {
		t.Error("expected duplicate-policy error on second CREATE without OR REPLACE / IF NOT EXISTS")
	}
}

func TestWTTiDBPolicy_3_OrReplaceOverwrites(t *testing.T) {
	c := New()
	if _, err := c.Exec("CREATE PLACEMENT POLICY p1 FOLLOWERS = 2", nil); err != nil {
		t.Fatalf("first CREATE: %v", err)
	}
	if _, err := c.Exec("CREATE OR REPLACE PLACEMENT POLICY p1 FOLLOWERS = 5", nil); err != nil {
		t.Fatalf("OR REPLACE: %v", err)
	}
	p := c.GetPlacementPolicy("p1")
	if p == nil || len(p.Options) != 1 || p.Options[0].IntValue != 5 {
		t.Errorf("OR REPLACE did not overwrite options: %+v", p)
	}
}

func TestWTTiDBPolicy_4_IfNotExistsSilences(t *testing.T) {
	c := New()
	if _, err := c.Exec("CREATE PLACEMENT POLICY p1 FOLLOWERS = 2", nil); err != nil {
		t.Fatalf("first CREATE: %v", err)
	}
	results, err := c.Exec("CREATE PLACEMENT POLICY IF NOT EXISTS p1 FOLLOWERS = 9", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if results[0].Error != nil {
		t.Errorf("IF NOT EXISTS should silence duplicate, got: %v", results[0].Error)
	}
	// Should NOT have been overwritten.
	p := c.GetPlacementPolicy("p1")
	if p.Options[0].IntValue != 2 {
		t.Errorf("IF NOT EXISTS overwrote original policy")
	}
}

func TestWTTiDBPolicy_5_AlterReplacesOptions(t *testing.T) {
	c := New()
	if _, err := c.Exec("CREATE PLACEMENT POLICY p1 PRIMARY_REGION = 'us' FOLLOWERS = 2", nil); err != nil {
		t.Fatalf("CREATE: %v", err)
	}
	if _, err := c.Exec("ALTER PLACEMENT POLICY p1 PRIMARY_REGION = 'eu'", nil); err != nil {
		t.Fatalf("ALTER: %v", err)
	}
	p := c.GetPlacementPolicy("p1")
	if len(p.Options) != 1 || p.Options[0].Value != "eu" {
		t.Errorf("ALTER did not replace options: %+v", p.Options)
	}
}

func TestWTTiDBPolicy_6_AlterUnknownErrors(t *testing.T) {
	c := New()
	results, err := c.Exec("ALTER PLACEMENT POLICY ghost PRIMARY_REGION = 'us'", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if results[0].Error == nil {
		t.Error("expected unknown-policy error on ALTER with no IF EXISTS")
	}
}

func TestWTTiDBPolicy_7_AlterIfExistsSilences(t *testing.T) {
	c := New()
	results, err := c.Exec("ALTER PLACEMENT POLICY IF EXISTS ghost PRIMARY_REGION = 'us'", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if results[0].Error != nil {
		t.Errorf("IF EXISTS should silence missing policy, got: %v", results[0].Error)
	}
}

func TestWTTiDBPolicy_8_Drop(t *testing.T) {
	c := New()
	if _, err := c.Exec("CREATE PLACEMENT POLICY p1 FOLLOWERS = 2; DROP PLACEMENT POLICY p1", nil); err != nil {
		t.Fatalf("Exec: %v", err)
	}
	if c.GetPlacementPolicy("p1") != nil {
		t.Error("policy should be gone after DROP")
	}
}

func TestWTTiDBPolicy_9_DropUnknownErrors(t *testing.T) {
	c := New()
	results, err := c.Exec("DROP PLACEMENT POLICY ghost", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if results[0].Error == nil {
		t.Error("expected unknown-policy error on DROP without IF EXISTS")
	}
}

func TestWTTiDBPolicy_10_DropIfExistsSilences(t *testing.T) {
	c := New()
	results, err := c.Exec("DROP PLACEMENT POLICY IF EXISTS ghost", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if results[0].Error != nil {
		t.Errorf("IF EXISTS should silence missing policy on DROP, got: %v", results[0].Error)
	}
}

// TestWTTiDBPolicy_11_EndToEndRoundTrip exercises the scenario that
// PR3b left broken: defining a policy and then referencing it from a
// CREATE TABLE. Before Tier 1, the CREATE PLACEMENT POLICY statement
// failed to parse and this whole flow was impossible.
func TestWTTiDBPolicy_11_EndToEndRoundTrip(t *testing.T) {
	c := New()
	sql := `
		CREATE PLACEMENT POLICY p1 PRIMARY_REGION = 'us-east' REGIONS = 'us-east,us-west';
		CREATE DATABASE ddb PLACEMENT POLICY = p1;
		USE ddb;
		CREATE TABLE t (id INT) PLACEMENT POLICY = p1;
	`
	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	for i, r := range results {
		if r.Error != nil {
			t.Fatalf("stmt %d (%q) failed: %v", i, r.SQL, r.Error)
		}
	}
	if c.GetPlacementPolicy("p1") == nil {
		t.Fatal("policy p1 not in catalog after round-trip")
	}
	if c.GetDatabase("ddb").PlacementPolicy != "p1" {
		t.Errorf("database ddb did not record policy reference")
	}
	tbl := c.GetDatabase("ddb").GetTable("t")
	if tbl.PlacementPolicy != "p1" {
		t.Errorf("table t did not record policy reference")
	}
}
