package catalog

import "testing"

// TestTiDBTiFlashLocationLabelsContainer cross-validates Tier 3 syntax
// against real TiDB v8.5.5: LOCATION LABELS on ALTER TABLE SET TIFLASH
// REPLICA and SET TIFLASH REPLICA as a DatabaseOption. Each case must
// be accepted by both TiDB and the omni catalog.
//
// Note: TiDB accepts the *syntax* even when no TiFlash servers are
// available (parse-time acceptance), but rejects replica counts > 0
// at execute time without TiFlash nodes. The container test uses 0 or
// uses labels-only-with-replica-0 shapes where TiDB does not require
// actual TiFlash infrastructure. For replica-count > 0 acceptance, we
// rely on the parser container test in tidb/parser/tidb_container_test.go.
func TestTiDBTiFlashLocationLabelsContainer(t *testing.T) {
	tc := startTiDBForCatalog(t)
	mustExecTiDB(t, tc, "CREATE DATABASE IF NOT EXISTS omni_tf_test")
	mustExecTiDB(t, tc, "USE omni_tf_test")
	t.Cleanup(func() { _, _ = tc.db.ExecContext(tc.ctx, "DROP DATABASE IF EXISTS omni_tf_test") })

	cases := []struct {
		name    string
		setup   string
		sql     string
		cleanup string
	}{
		{
			name:    "alter_table_replica_zero_with_labels",
			setup:   "CREATE TABLE IF NOT EXISTS t_tf_l (id INT PRIMARY KEY)",
			sql:     "ALTER TABLE t_tf_l SET TIFLASH REPLICA 0 LOCATION LABELS 'zone'",
			cleanup: "DROP TABLE IF EXISTS t_tf_l",
		},
		{
			name:    "alter_table_multi_label",
			setup:   "CREATE TABLE IF NOT EXISTS t_tf_ml (id INT PRIMARY KEY)",
			sql:     "ALTER TABLE t_tf_ml SET TIFLASH REPLICA 0 LOCATION LABELS 'zone', 'rack'",
			cleanup: "DROP TABLE IF EXISTS t_tf_ml",
		},
		{
			// TiDB rejects ALTER DATABASE SET TIFLASH REPLICA when the
			// database has no tables ("Empty database" error, despite
			// using code 1049). Setup adds a placeholder table so TiDB
			// can propagate the replica setting.
			name:    "alter_database_replica_zero",
			setup:   "CREATE TABLE IF NOT EXISTS t_db_ph (id INT PRIMARY KEY)",
			sql:     "ALTER DATABASE omni_tf_test SET TIFLASH REPLICA 0",
			cleanup: "DROP TABLE IF EXISTS t_db_ph",
		},
		{
			name:    "alter_database_replica_zero_with_labels",
			setup:   "CREATE TABLE IF NOT EXISTS t_db_ph2 (id INT PRIMARY KEY)",
			sql:     "ALTER DATABASE omni_tf_test SET TIFLASH REPLICA 0 LOCATION LABELS 'zone'",
			cleanup: "DROP TABLE IF EXISTS t_db_ph2",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.setup != "" {
				mustExecTiDB(t, tc, c.setup)
			}
			if _, err := tc.db.ExecContext(tc.ctx, c.sql); err != nil {
				t.Fatalf("TiDB rejected %q: %v", c.sql, err)
			}
			cat := New()
			if _, err := cat.Exec("CREATE DATABASE omni_tf_test; USE omni_tf_test;", nil); err != nil {
				t.Fatalf("cat setup: %v", err)
			}
			if c.setup != "" {
				if _, err := cat.Exec(c.setup, nil); err != nil {
					t.Fatalf("cat setup SQL %q: %v", c.setup, err)
				}
			}
			if _, err := cat.Exec(c.sql, nil); err != nil {
				t.Fatalf("catalog rejected %q: %v", c.sql, err)
			}
			if c.cleanup != "" {
				_, _ = tc.db.ExecContext(tc.ctx, c.cleanup)
			}
		})
	}
}

// TestTiDBTiFlashGrammarNegatives — parse-time negatives common to
// both omni and TiDB. Mirrors the parser's unit-level negatives but
// confirms TiDB also rejects each.
func TestTiDBTiFlashGrammarNegatives(t *testing.T) {
	tc := startTiDBForCatalog(t)
	mustExecTiDB(t, tc, "CREATE DATABASE IF NOT EXISTS omni_tfn_test")
	mustExecTiDB(t, tc, "USE omni_tfn_test")
	t.Cleanup(func() { _, _ = tc.db.ExecContext(tc.ctx, "DROP DATABASE IF EXISTS omni_tfn_test") })
	mustExecTiDB(t, tc, "CREATE TABLE t_neg (id INT PRIMARY KEY)")

	cases := []struct {
		name string
		sql  string
	}{
		{"missing_labels_after_location", "ALTER TABLE t_neg SET TIFLASH REPLICA 0 LOCATION"},
		{"empty_label_list", "ALTER TABLE t_neg SET TIFLASH REPLICA 0 LOCATION LABELS"},
		{"trailing_comma", "ALTER TABLE t_neg SET TIFLASH REPLICA 0 LOCATION LABELS 'a',"},
		{"bare_identifier_label", "ALTER TABLE t_neg SET TIFLASH REPLICA 0 LOCATION LABELS zone"},
		// DEFAULT is allowed before CHARSET/COLLATE/ENCRYPTION/PLACEMENT
		// POLICY on DatabaseOption but NOT before SET TIFLASH REPLICA
		// (no DefaultKwdOpt on that grammar arm).
		{"default_create_set_tiflash", "CREATE DATABASE omni_tfn_test_x DEFAULT SET TIFLASH REPLICA 1"},
		{"default_alter_set_tiflash", "ALTER DATABASE omni_tfn_test DEFAULT SET TIFLASH REPLICA 1"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := tc.db.ExecContext(tc.ctx, c.sql); err == nil {
				t.Errorf("TiDB unexpectedly accepted %q", c.sql)
			}
			cat := New()
			_, _ = cat.Exec("CREATE DATABASE omni_tfn_test; USE omni_tfn_test; CREATE TABLE t_neg (id INT PRIMARY KEY);", nil)
			if _, err := cat.Exec(c.sql, nil); err == nil {
				t.Errorf("omni unexpectedly accepted %q (oracle divergence)", c.sql)
			}
		})
	}
}
