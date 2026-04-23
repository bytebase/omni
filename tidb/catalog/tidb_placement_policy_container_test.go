package catalog

import (
	"strings"
	"testing"
)

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

// TestTiDBPlacementPolicyGrammarNegatives pins that omni and TiDB
// agree on rejection of malformed PlacementOptionList inputs. The
// upstream grammar at parser.y:1999-2011 requires >=1 option and
// forbids a trailing comma; this test confirms TiDB enforces both at
// parse time, so omni's parser enforcement is not over-strict.
func TestTiDBPlacementPolicyGrammarNegatives(t *testing.T) {
	tc := startTiDBForCatalog(t)
	t.Cleanup(func() { _, _ = tc.db.ExecContext(tc.ctx, "DROP PLACEMENT POLICY IF EXISTS pp_neg") })

	cases := []struct {
		name string
		sql  string
	}{
		{"create_empty_options", "CREATE PLACEMENT POLICY pp_neg"},
		{"create_trailing_comma", "CREATE PLACEMENT POLICY pp_neg PRIMARY_REGION = 'us',"},
		{"alter_empty_options", "ALTER PLACEMENT POLICY pp_neg"},
		{"alter_trailing_comma", "ALTER PLACEMENT POLICY pp_neg PRIMARY_REGION = 'us',"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := tc.db.ExecContext(tc.ctx, c.sql); err == nil {
				t.Errorf("TiDB unexpectedly accepted malformed SQL %q", c.sql)
			}
			// Omni must also reject. Parse-level check is sufficient;
			// the Exec path routes parse errors as the top-level error.
			cat := New()
			if _, err := cat.Exec(c.sql, nil); err == nil {
				t.Errorf("omni unexpectedly accepted malformed SQL %q (oracle divergence)", c.sql)
			}
		})
	}
}

// TestTiDBPlacementPolicyDefaultSentinel pins the "default"
// special-name behavior against real TiDB. Upstream short-circuits
// at pkg/ddl/placement_policy.go defaultPlacementPolicyName — a
// reference to "default" must bypass policy lookup and be accepted
// even when no policy by that name exists.
func TestTiDBPlacementPolicyDefaultSentinel(t *testing.T) {
	tc := startTiDBForCatalog(t)
	mustExecTiDB(t, tc, "CREATE DATABASE IF NOT EXISTS omni_pp_default")
	mustExecTiDB(t, tc, "USE omni_pp_default")
	t.Cleanup(func() { _, _ = tc.db.ExecContext(tc.ctx, "DROP DATABASE IF EXISTS omni_pp_default") })

	inputs := []struct {
		label string
		sql   string
		drop  string
	}{
		{"identifier", "CREATE TABLE t_def_id (id INT) PLACEMENT POLICY = default", "DROP TABLE IF EXISTS t_def_id"},
		{"quoted_string", "CREATE TABLE t_def_str (id INT) PLACEMENT POLICY = 'default'", "DROP TABLE IF EXISTS t_def_str"},
		{"upper_keyword", "CREATE TABLE t_def_up (id INT) PLACEMENT POLICY = DEFAULT", "DROP TABLE IF EXISTS t_def_up"},
	}
	for _, in := range inputs {
		t.Run(in.label, func(t *testing.T) {
			// 1. TiDB must accept without any CREATE PLACEMENT POLICY default.
			if _, err := tc.db.ExecContext(tc.ctx, in.sql); err != nil {
				t.Fatalf("TiDB rejected sentinel form %q: %v", in.sql, err)
			}
			// 2. Catalog must also accept.
			cat := New()
			if _, err := cat.Exec("CREATE DATABASE omni_pp_default; USE omni_pp_default;", nil); err != nil {
				t.Fatalf("cat setup: %v", err)
			}
			if _, err := cat.Exec(in.sql, nil); err != nil {
				t.Fatalf("catalog rejected sentinel form %q: %v", in.sql, err)
			}
			_, _ = tc.db.ExecContext(tc.ctx, in.drop)
		})
	}
}

// TestTiDBPlacementPolicyAlterReplaceSemantics pins the claim made in
// catalog/placement_policy.go alterPlacementPolicy: ALTER replaces the
// option list wholesale, it does NOT merge. Uses SHOW CREATE PLACEMENT
// POLICY to read back actual TiDB state rather than just asserting
// that the ALTER was accepted.
func TestTiDBPlacementPolicyAlterReplaceSemantics(t *testing.T) {
	tc := startTiDBForCatalog(t)
	t.Cleanup(func() { _, _ = tc.db.ExecContext(tc.ctx, "DROP PLACEMENT POLICY IF EXISTS pp_rb") })

	mustExecTiDB(t, tc, "CREATE PLACEMENT POLICY pp_rb FOLLOWERS = 2")
	mustExecTiDB(t, tc, "ALTER PLACEMENT POLICY pp_rb PRIMARY_REGION = 'us-west-1' REGIONS = 'us-west-1'")

	var name, def string
	row := tc.db.QueryRowContext(tc.ctx, "SHOW CREATE PLACEMENT POLICY pp_rb")
	if err := row.Scan(&name, &def); err != nil {
		t.Fatalf("SHOW CREATE PLACEMENT POLICY: %v", err)
	}
	// After ALTER, FOLLOWERS=2 must be gone. Case-insensitive match
	// since TiDB's SHOW output may uppercase keywords differently
	// across versions.
	low := strings.ToLower(def)
	if strings.Contains(low, "followers") {
		t.Errorf("ALTER did not replace option list; FOLLOWERS leaked: %s", def)
	}
	if !strings.Contains(low, "primary_region") {
		t.Errorf("ALTER lost PRIMARY_REGION: %s", def)
	}

	// Now verify catalog agrees: fresh catalog, same SQL flow.
	cat := New()
	if _, err := cat.Exec("CREATE PLACEMENT POLICY pp_rb FOLLOWERS = 2", nil); err != nil {
		t.Fatalf("cat CREATE: %v", err)
	}
	if _, err := cat.Exec("ALTER PLACEMENT POLICY pp_rb PRIMARY_REGION = 'us-west-1' REGIONS = 'us-west-1'", nil); err != nil {
		t.Fatalf("cat ALTER: %v", err)
	}
	p := cat.GetPlacementPolicy("pp_rb")
	for _, o := range p.Options {
		if o.Name == "FOLLOWERS" {
			t.Errorf("catalog ALTER leaked FOLLOWERS: %+v", p.Options)
		}
	}
	if len(p.Options) != 2 {
		t.Errorf("catalog ALTER result: want 2 options, got %d: %+v", len(p.Options), p.Options)
	}
}
