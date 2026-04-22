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

// TestWTTiDBPolicy_AlterReplacesWholesale verifies that ALTER
// replaces the option list rather than merging. The container test
// asserts TiDB accepts the ALTER; this test assertions the catalog's
// in-memory state actually reflects the replace, pinning the behavior
// claimed in catalog/placement_policy.go's alterPlacementPolicy comment.
func TestWTTiDBPolicy_AlterReplacesWholesale(t *testing.T) {
	c := New()
	wtExecRaw(t, c, "CREATE PLACEMENT POLICY p PRIMARY_REGION = 'us' FOLLOWERS = 2")
	wtExecRaw(t, c, "ALTER PLACEMENT POLICY p PRIMARY_REGION = 'eu'")

	p := c.GetPlacementPolicy("p")
	if p == nil {
		t.Fatal("policy missing")
	}
	if len(p.Options) != 1 {
		t.Fatalf("ALTER should have replaced option list (want 1 entry), got %d: %+v", len(p.Options), p.Options)
	}
	if p.Options[0].Name != "PRIMARY_REGION" || p.Options[0].Value != "eu" {
		t.Errorf("option content: %+v", p.Options[0])
	}
	// FOLLOWERS = 2 must be gone.
	for _, o := range p.Options {
		if o.Name == "FOLLOWERS" {
			t.Error("FOLLOWERS=2 from CREATE leaked into ALTER result (merge, not replace)")
		}
	}
}

// TestWTTiDBPolicy_RefValidation verifies that CREATE TABLE and
// CREATE DATABASE with an unknown policy name are rejected (TiDB
// error 8237 parity).
func TestWTTiDBPolicy_RefValidation(t *testing.T) {
	c := New()
	// No policy defined yet; table reference must fail.
	if _, err := c.Exec("CREATE DATABASE testdb; USE testdb;", nil); err != nil {
		t.Fatalf("setup: %v", err)
	}
	results, err := c.Exec("CREATE TABLE t (id INT) PLACEMENT POLICY = ghost", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if results[0].Error == nil {
		t.Error("expected unknown-policy error on CREATE TABLE referencing undefined policy")
	}
	// Database reference to unknown policy.
	results, err = c.Exec("CREATE DATABASE d2 PLACEMENT POLICY = ghost", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if results[0].Error == nil {
		t.Error("expected unknown-policy error on CREATE DATABASE referencing undefined policy")
	}
}

// TestWTTiDBPolicy_DropInUseRejected verifies that DROP PLACEMENT
// POLICY is rejected when any table or database still references it.
// Mirrors TiDB error 8240.
func TestWTTiDBPolicy_DropInUseRejected(t *testing.T) {
	c := New()
	wtExecRaw(t, c, "CREATE PLACEMENT POLICY p1 PRIMARY_REGION = 'us'")
	wtExecRaw(t, c, "CREATE DATABASE testdb; USE testdb;")
	wtExecRaw(t, c, "CREATE TABLE t (id INT) PLACEMENT POLICY = p1")

	// Drop must fail: table references policy.
	results, err := c.Exec("DROP PLACEMENT POLICY p1", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if results[0].Error == nil {
		t.Error("expected in-use error on DROP PLACEMENT POLICY with table reference")
	}

	// Policy still in catalog (DROP didn't partially succeed).
	if c.GetPlacementPolicy("p1") == nil {
		t.Error("failed DROP left policy removed from catalog")
	}

	// Remove the reference, DROP succeeds.
	wtExecRaw(t, c, "DROP TABLE t")
	if _, err := c.Exec("DROP PLACEMENT POLICY p1", nil); err != nil {
		t.Fatalf("DROP after removing ref: %v", err)
	}
	if c.GetPlacementPolicy("p1") != nil {
		t.Error("DROP did not remove policy from catalog")
	}
}

// TestWTTiDBPolicy_DropInUseByDatabase is the DB-reference equivalent
// of TestWTTiDBPolicy_DropInUseRejected.
func TestWTTiDBPolicy_DropInUseByDatabase(t *testing.T) {
	c := New()
	wtExecRaw(t, c, "CREATE PLACEMENT POLICY p1 PRIMARY_REGION = 'us'")
	wtExecRaw(t, c, "CREATE DATABASE ddb PLACEMENT POLICY = p1")

	results, err := c.Exec("DROP PLACEMENT POLICY p1", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if results[0].Error == nil {
		t.Error("expected in-use error on DROP PLACEMENT POLICY with database reference")
	}
}

// TestWTTiDBPolicy_DefaultSentinel verifies the special-cased "default"
// policy name: it must bypass ref validation and clear the assigned
// policy to the empty string (matching TiDB's short-circuit at
// pkg/ddl/placement_policy.go's defaultPlacementPolicyName check).
// All three input forms — `default`, 'default', DEFAULT — collapse to
// the same catalog result.
func TestWTTiDBPolicy_DefaultSentinel(t *testing.T) {
	inputs := []struct {
		label string
		sql   string
	}{
		{"identifier default", "CREATE TABLE t (id INT) PLACEMENT POLICY = default"},
		{"quoted 'default'", "CREATE TABLE t (id INT) PLACEMENT POLICY = 'default'"},
		// Upper-case DEFAULT is also accepted; catalog normalizes via
		// toLower in isDefaultPolicyName.
		{"keyword DEFAULT", "CREATE TABLE t (id INT) PLACEMENT POLICY = DEFAULT"},
	}
	for _, in := range inputs {
		t.Run(in.label, func(t *testing.T) {
			c := wtSetup(t)
			// No policy defined; default must still work.
			wtExec(t, c, in.sql)
			tbl := c.GetDatabase("testdb").GetTable("t")
			if tbl == nil {
				t.Fatal("table not created")
			}
			if tbl.PlacementPolicy != "" {
				t.Errorf("default sentinel should clear policy, got %q", tbl.PlacementPolicy)
			}
		})
	}
}

// TestWTTiDBPolicy_DefaultOnDatabase verifies the sentinel works on
// CREATE DATABASE too (same code path via validatePolicyRef +
// resolvePolicyRef in dbcmds.go).
func TestWTTiDBPolicy_DefaultOnDatabase(t *testing.T) {
	c := New()
	wtExecRaw(t, c, "CREATE DATABASE ddb PLACEMENT POLICY = default")
	db := c.GetDatabase("ddb")
	if db == nil {
		t.Fatal("database not created")
	}
	if db.PlacementPolicy != "" {
		t.Errorf("default sentinel should clear DB policy, got %q", db.PlacementPolicy)
	}
}

// TestWTTiDBPolicy_DefaultClearsExisting verifies that ALTER with
// `default` clears a previously-set policy rather than setting the
// field to the literal string "default".
func TestWTTiDBPolicy_DefaultClearsExisting(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE PLACEMENT POLICY p1 PRIMARY_REGION = 'us'")
	wtExec(t, c, "CREATE TABLE t (id INT) PLACEMENT POLICY = p1")
	// Sanity.
	if tbl := c.GetDatabase("testdb").GetTable("t"); tbl.PlacementPolicy != "p1" {
		t.Fatalf("setup failed: %q", tbl.PlacementPolicy)
	}
	wtExec(t, c, "ALTER TABLE t PLACEMENT POLICY = default")
	if tbl := c.GetDatabase("testdb").GetTable("t"); tbl.PlacementPolicy != "" {
		t.Errorf("ALTER TABLE ... = default should clear policy, got %q", tbl.PlacementPolicy)
	}
}

// TestWTTiDBPolicy_ExplicitDefaultPolicyShadows documents a quirk: if
// a user literally names a policy `default` (via backtick-quoting), it
// still won't actually be referenced by unquoted/string `default`
// because the sentinel short-circuit fires FIRST. This matches TiDB's
// upstream check order at placement_policy.go:~478: the defaultName
// compare happens before PolicyByName lookup.
func TestWTTiDBPolicy_ExplicitDefaultPolicyShadows(t *testing.T) {
	c := wtSetup(t)
	// Users CAN create a backtick-quoted `default` policy in TiDB, but
	// it becomes unreachable via PLACEMENT POLICY = default — the
	// sentinel always wins. Omni must match.
	wtExec(t, c, "CREATE PLACEMENT POLICY `default` PRIMARY_REGION = 'us-east'")
	wtExec(t, c, "CREATE TABLE t (id INT) PLACEMENT POLICY = default")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl.PlacementPolicy != "" {
		t.Errorf("sentinel should win over literally-named `default` policy; got %q", tbl.PlacementPolicy)
	}
}

// wtExecRaw executes SQL on an arbitrary catalog (no testdb setup),
// fataling on any parse or exec error. Suited for tests that manage
// their own database-creation order.
func wtExecRaw(t *testing.T, c *Catalog, sql string) {
	t.Helper()
	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatalf("Exec parse error on %q: %v", sql, err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatalf("Exec runtime error on %q (stmt %d): %v", sql, r.Index, r.Error)
		}
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
