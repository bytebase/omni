package catalog

import "testing"

// PR3b follow-up walkthrough tests. Each test exercises a code path
// missed in the first cut: DB-level PLACEMENT POLICY, ALTER TABLE ADD
// PRIMARY KEY with CLUSTERED hoisting on both inline-column and
// table-level paths, and deep-copy correctness for Clustered.

// TestWTTiDB_2_1_DatabasePlacementPolicy verifies CREATE DATABASE ...
// PLACEMENT POLICY = p populates Database.PlacementPolicy and ALTER
// DATABASE can change it.
func TestWTTiDB_2_1_DatabasePlacementPolicy(t *testing.T) {
	c := New()
	// Policies must be defined before a database can reference them.
	if _, err := c.Exec("CREATE PLACEMENT POLICY p1 PRIMARY_REGION = 'us'; CREATE PLACEMENT POLICY p2 PRIMARY_REGION = 'eu';", nil); err != nil {
		t.Fatalf("CREATE POLICIES: %v", err)
	}
	if _, err := c.Exec("CREATE DATABASE ddb PLACEMENT POLICY = p1", nil); err != nil {
		t.Fatalf("CREATE DATABASE: %v", err)
	}
	db := c.GetDatabase("ddb")
	if db == nil {
		t.Fatal("database ddb not found")
	}
	if db.PlacementPolicy != "p1" {
		t.Errorf("expected PlacementPolicy=p1, got %q", db.PlacementPolicy)
	}

	if _, err := c.Exec("ALTER DATABASE ddb PLACEMENT POLICY = p2", nil); err != nil {
		t.Fatalf("ALTER DATABASE: %v", err)
	}
	if db.PlacementPolicy != "p2" {
		t.Errorf("after ALTER: expected PlacementPolicy=p2, got %q", db.PlacementPolicy)
	}
}

// TestWTTiDB_2_2_AlterAddPKClustered verifies that ALTER TABLE ADD
// PRIMARY KEY (...) CLUSTERED propagates the flag onto the new
// constraint. Previously dropped silently.
func TestWTTiDB_2_2_AlterAddPKClustered(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT)")
	wtExec(t, c, "ALTER TABLE t ADD PRIMARY KEY (id) CLUSTERED")

	tbl := c.GetDatabase("testdb").GetTable("t")
	for _, con := range tbl.Constraints {
		if con.Type == ConPrimaryKey {
			if con.Clustered == nil || !*con.Clustered {
				t.Error("expected Clustered=true on PK added via ALTER TABLE ADD PRIMARY KEY (...) CLUSTERED")
			}
			return
		}
	}
	t.Error("primary key constraint not found after ALTER")
}

// TestWTTiDB_2_2b_AlterAddPKNonClustered covers the NONCLUSTERED
// variant — a false value is just as important to hoist as true, since
// nil and &false have different semantics.
func TestWTTiDB_2_2b_AlterAddPKNonClustered(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT)")
	wtExec(t, c, "ALTER TABLE t ADD PRIMARY KEY (id) NONCLUSTERED")

	tbl := c.GetDatabase("testdb").GetTable("t")
	for _, con := range tbl.Constraints {
		if con.Type == ConPrimaryKey {
			if con.Clustered == nil {
				t.Error("expected Clustered=&false (set), got nil")
			} else if *con.Clustered {
				t.Error("expected Clustered=false, got true")
			}
			return
		}
	}
	t.Error("primary key constraint not found after ALTER")
}

// TestWTTiDB_2_3_AlterAddColumnInlinePKUnsupported verifies that ALTER TABLE
// ADD COLUMN with an inline PRIMARY KEY constraint is rejected with 8200,
// matching TiDB v8.5.5 (pkg/ddl/add_column.go checkUnsupportedColumnConstraint).
func TestWTTiDB_2_3_AlterAddColumnInlinePKUnsupported(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (x INT)")
	results, err := c.Exec("ALTER TABLE t ADD COLUMN id BIGINT PRIMARY KEY CLUSTERED", nil)
	if err != nil {
		t.Fatalf("unexpected parse err: %v", err)
	}
	if len(results) != 1 || results[0].Error == nil {
		t.Fatal("expected error 8200 for ALTER TABLE ADD COLUMN ... PRIMARY KEY")
	}
	catErr, ok := results[0].Error.(*Error)
	if !ok || catErr.Code != ErrUnsupportedDDL {
		t.Errorf("got err %v, want code %d", results[0].Error, ErrUnsupportedDDL)
	}
}
