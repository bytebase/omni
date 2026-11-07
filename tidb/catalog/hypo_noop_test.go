package catalog

import "testing"

// The hypothetical forms the parser accepts (DROP HYPO INDEX, SET HYPO
// TIFLASH REPLICA) are optimizer-only what-if operations. They must not
// mutate the physical catalog, even when a real object shares the name.

// TestDropHypoIndexLeavesRealIndex checks that DROP HYPO INDEX never removes a
// persistent index (or its constraint) of the same name.
func TestDropHypoIndexLeavesRealIndex(t *testing.T) {
	c := New()
	mustExec(t, c, "CREATE DATABASE test")
	c.SetCurrentDatabase("test")
	mustExec(t, c, "CREATE TABLE t (id INT NOT NULL, name VARCHAR(100))")
	mustExec(t, c, "CREATE UNIQUE INDEX i ON t (id)")

	tbl := c.GetDatabase("test").GetTable("t")
	if len(tbl.Indexes) != 1 || len(tbl.Constraints) != 1 {
		t.Fatalf("setup: indexes=%d constraints=%d, want 1 and 1", len(tbl.Indexes), len(tbl.Constraints))
	}

	// Both the plain and the IF EXISTS form are no-ops on the physical index.
	mustExec(t, c, "DROP HYPO INDEX i ON t")
	mustExec(t, c, "DROP HYPO INDEX IF EXISTS i ON t")

	tbl = c.GetDatabase("test").GetTable("t")
	if len(tbl.Indexes) != 1 {
		t.Errorf("DROP HYPO INDEX removed the real index: indexes=%d, want 1", len(tbl.Indexes))
	}
	if len(tbl.Constraints) != 1 {
		t.Errorf("DROP HYPO INDEX removed the real constraint: constraints=%d, want 1", len(tbl.Constraints))
	}

	// A real DROP INDEX still removes it, so the guard is scoped to HYPO.
	mustExec(t, c, "DROP INDEX i ON t")
	tbl = c.GetDatabase("test").GetTable("t")
	if len(tbl.Indexes) != 0 {
		t.Errorf("DROP INDEX left the index: indexes=%d, want 0", len(tbl.Indexes))
	}
}

// TestSetHypoTiFlashReplicaLeavesPhysical checks that SET HYPO TIFLASH REPLICA
// does not overwrite the real replica count or location labels.
func TestSetHypoTiFlashReplicaLeavesPhysical(t *testing.T) {
	c := New()
	mustExec(t, c, "CREATE DATABASE test")
	c.SetCurrentDatabase("test")
	mustExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY)")
	mustExec(t, c, "ALTER TABLE t SET TIFLASH REPLICA 2 LOCATION LABELS 'zone', 'rack'")

	mustExec(t, c, "ALTER TABLE t SET HYPO TIFLASH REPLICA 1")

	tbl := c.GetDatabase("test").GetTable("t")
	if tbl.TiFlashReplica != 2 {
		t.Errorf("SET HYPO TIFLASH REPLICA overwrote the real count: TiFlashReplica=%d, want 2", tbl.TiFlashReplica)
	}
	if len(tbl.TiFlashLocationLabels) != 2 {
		t.Errorf("SET HYPO TIFLASH REPLICA cleared the real labels: %v, want [zone rack]", tbl.TiFlashLocationLabels)
	}

	// A real SET TIFLASH REPLICA still applies, so the guard is scoped to HYPO.
	mustExec(t, c, "ALTER TABLE t SET TIFLASH REPLICA 5")
	tbl = c.GetDatabase("test").GetTable("t")
	if tbl.TiFlashReplica != 5 {
		t.Errorf("SET TIFLASH REPLICA did not apply: TiFlashReplica=%d, want 5", tbl.TiFlashReplica)
	}
}
