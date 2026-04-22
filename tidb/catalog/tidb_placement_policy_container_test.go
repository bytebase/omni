package catalog

import "testing"

// TestTiDBPlacementPolicyContainer cross-validates that omni's
// PLACEMENT POLICY DDL matches real TiDB v8.5.5 acceptance. Reuses the
// startTiDBForCatalog helper from tidb_container_test.go.
//
// Note: TiDB requires at least some option to be present on CREATE;
// the grammar allows zero-option policies but TiDB rejects them at
// semantic time ("Could not create placement policy with empty info").
// We therefore always send at least one option.
func TestTiDBPlacementPolicyContainer(t *testing.T) {
	tc := startTiDBForCatalog(t)

	mustExecTiDB(t, tc, "CREATE DATABASE IF NOT EXISTS omni_pp_test")
	mustExecTiDB(t, tc, "USE omni_pp_test")
	t.Cleanup(func() { _, _ = tc.db.ExecContext(tc.ctx, "DROP DATABASE IF EXISTS omni_pp_test") })

	cases := []struct {
		name    string
		setup   string
		sql     string
		cleanup string
	}{
		{
			name:    "create_primary_region",
			sql:     "CREATE PLACEMENT POLICY pp_cpr PRIMARY_REGION = 'us-east-1' REGIONS = 'us-east-1'",
			cleanup: "DROP PLACEMENT POLICY IF EXISTS pp_cpr",
		},
		{
			name:    "create_followers",
			sql:     "CREATE PLACEMENT POLICY pp_cf FOLLOWERS = 2",
			cleanup: "DROP PLACEMENT POLICY IF EXISTS pp_cf",
		},
		{
			name:    "create_constraints_list",
			sql:     `CREATE PLACEMENT POLICY pp_ccl CONSTRAINTS='[+region=us-east-1]'`,
			cleanup: "DROP PLACEMENT POLICY IF EXISTS pp_ccl",
		},
		{
			name:    "create_or_replace",
			setup:   "CREATE PLACEMENT POLICY pp_or FOLLOWERS = 2",
			sql:     "CREATE OR REPLACE PLACEMENT POLICY pp_or FOLLOWERS = 3",
			cleanup: "DROP PLACEMENT POLICY IF EXISTS pp_or",
		},
		{
			name:    "alter_placement_policy",
			setup:   "CREATE PLACEMENT POLICY pp_alter FOLLOWERS = 2",
			sql:     "ALTER PLACEMENT POLICY pp_alter PRIMARY_REGION = 'us-east-1' REGIONS = 'us-east-1'",
			cleanup: "DROP PLACEMENT POLICY IF EXISTS pp_alter",
		},
		{
			name:    "drop_if_exists_missing",
			sql:     "DROP PLACEMENT POLICY IF EXISTS never_defined_pp",
			cleanup: "",
		},
		{
			// The end-to-end case that PR3b couldn't run: define a
			// policy, then reference it from CREATE TABLE. Without
			// Tier 1, catalog.Exec rejected the CREATE PLACEMENT
			// POLICY line and the container test had to skip.
			name:    "policy_then_table_reference",
			setup:   "CREATE PLACEMENT POLICY pp_tbl PRIMARY_REGION='us-east-1' REGIONS='us-east-1'",
			sql:     "CREATE TABLE t_pp (id INT) PLACEMENT POLICY = pp_tbl",
			cleanup: "DROP TABLE IF EXISTS t_pp; DROP PLACEMENT POLICY IF EXISTS pp_tbl",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.setup != "" {
				mustExecTiDB(t, tc, c.setup)
			}
			// 1. TiDB must accept.
			if _, err := tc.db.ExecContext(tc.ctx, c.sql); err != nil {
				t.Fatalf("TiDB rejected %q: %v", c.sql, err)
			}
			// 2. Catalog must accept through its own parser+exec.
			cat := New()
			catSetup := "CREATE DATABASE omni_pp_test; USE omni_pp_test;"
			if _, err := cat.Exec(catSetup, nil); err != nil {
				t.Fatalf("catalog setup: %v", err)
			}
			if c.setup != "" {
				if _, err := cat.Exec(c.setup, nil); err != nil {
					t.Fatalf("catalog rejected setup %q: %v", c.setup, err)
				}
			}
			results, err := cat.Exec(c.sql, nil)
			if err != nil {
				t.Fatalf("catalog rejected %q: %v", c.sql, err)
			}
			for _, r := range results {
				if r.Error != nil {
					t.Fatalf("catalog exec error on %q: %v", c.sql, r.Error)
				}
			}
			if c.cleanup != "" {
				_, _ = tc.db.ExecContext(tc.ctx, c.cleanup)
			}
		})
	}
}
