package catalog

import (
	"strings"
	"testing"
)

func TestMigrationPartition(t *testing.T) {
	t.Run("CREATE TABLE PARTITION BY RANGE/LIST/HASH for partitioned table", func(t *testing.T) {
		from, err := LoadSQL("")
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(`
			CREATE TABLE orders_list (id int, region text) PARTITION BY LIST (region);
			CREATE TABLE orders_range (id int, created_at date) PARTITION BY RANGE (created_at);
			CREATE TABLE orders_hash (id int, val text) PARTITION BY HASH (id);
		`)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := GenerateMigration(from, to, diff)
		creates := plan.Filter(func(op MigrationOp) bool {
			return op.Type == OpCreateTable
		})
		if len(creates.Ops) != 3 {
			t.Fatalf("expected 3 CreateTable ops, got %d; ops: %v", len(creates.Ops), plan.Ops)
		}

		// Check each partitioned table has the correct PARTITION BY clause.
		found := map[string]bool{"LIST": false, "RANGE": false, "HASH": false}
		for _, op := range creates.Ops {
			if strings.Contains(op.SQL, "PARTITION BY LIST") {
				found["LIST"] = true
			}
			if strings.Contains(op.SQL, "PARTITION BY RANGE") {
				found["RANGE"] = true
			}
			if strings.Contains(op.SQL, "PARTITION BY HASH") {
				found["HASH"] = true
			}
		}
		for strategy, ok := range found {
			if !ok {
				t.Errorf("expected PARTITION BY %s in one of the ops", strategy)
			}
		}
	})

	t.Run("CREATE TABLE PARTITION OF range partition", func(t *testing.T) {
		from, err := LoadSQL(`
			CREATE TABLE orders (id int, created_at date) PARTITION BY RANGE (created_at);
		`)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(`
			CREATE TABLE orders (id int, created_at date) PARTITION BY RANGE (created_at);
			CREATE TABLE orders_2024 PARTITION OF orders FOR VALUES FROM ('2024-01-01') TO ('2025-01-01');
		`)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := GenerateMigration(from, to, diff)
		creates := plan.Filter(func(op MigrationOp) bool {
			return op.Type == OpCreateTable
		})
		if len(creates.Ops) != 1 {
			t.Fatalf("expected 1 CreateTable op for partition child, got %d; all ops: %v", len(creates.Ops), plan.Ops)
		}
		sql := creates.Ops[0].SQL
		if !strings.Contains(sql, "PARTITION OF") {
			t.Errorf("expected PARTITION OF, got: %s", sql)
		}
		if !strings.Contains(sql, "FOR VALUES FROM") {
			t.Errorf("expected FOR VALUES FROM, got: %s", sql)
		}
		if !strings.Contains(sql, "TO") {
			t.Errorf("expected TO in range bound, got: %s", sql)
		}
		if creates.Ops[0].ObjectName != "orders_2024" {
			t.Errorf("expected object name orders_2024, got: %s", creates.Ops[0].ObjectName)
		}
	})

	t.Run("CREATE TABLE PARTITION OF list partition", func(t *testing.T) {
		from, err := LoadSQL(`
			CREATE TABLE orders (id int, region text) PARTITION BY LIST (region);
		`)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(`
			CREATE TABLE orders (id int, region text) PARTITION BY LIST (region);
			CREATE TABLE orders_us PARTITION OF orders FOR VALUES IN ('US', 'CA');
		`)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := GenerateMigration(from, to, diff)
		creates := plan.Filter(func(op MigrationOp) bool {
			return op.Type == OpCreateTable
		})
		if len(creates.Ops) != 1 {
			t.Fatalf("expected 1 CreateTable op for list partition, got %d; all ops: %v", len(creates.Ops), plan.Ops)
		}
		sql := creates.Ops[0].SQL
		if !strings.Contains(sql, "PARTITION OF") {
			t.Errorf("expected PARTITION OF, got: %s", sql)
		}
		if !strings.Contains(sql, "FOR VALUES IN") {
			t.Errorf("expected FOR VALUES IN, got: %s", sql)
		}
	})

	t.Run("CREATE TABLE PARTITION OF default partition", func(t *testing.T) {
		from, err := LoadSQL(`
			CREATE TABLE orders (id int, region text) PARTITION BY LIST (region);
		`)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(`
			CREATE TABLE orders (id int, region text) PARTITION BY LIST (region);
			CREATE TABLE orders_default PARTITION OF orders DEFAULT;
		`)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := GenerateMigration(from, to, diff)
		creates := plan.Filter(func(op MigrationOp) bool {
			return op.Type == OpCreateTable
		})
		if len(creates.Ops) != 1 {
			t.Fatalf("expected 1 CreateTable op for default partition, got %d; all ops: %v", len(creates.Ops), plan.Ops)
		}
		sql := creates.Ops[0].SQL
		if !strings.Contains(sql, "PARTITION OF") {
			t.Errorf("expected PARTITION OF, got: %s", sql)
		}
		if !strings.Contains(sql, "DEFAULT") {
			t.Errorf("expected DEFAULT, got: %s", sql)
		}
	})

	t.Run("ALTER TABLE REPLICA IDENTITY for replica identity change", func(t *testing.T) {
		from, err := LoadSQL(`CREATE TABLE t1 (id int);`)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(`CREATE TABLE t1 (id int); ALTER TABLE t1 REPLICA IDENTITY FULL;`)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := GenerateMigration(from, to, diff)
		alters := plan.Filter(func(op MigrationOp) bool {
			return op.Type == OpAlterTable
		})
		if len(alters.Ops) != 1 {
			t.Fatalf("expected 1 AlterTable op for REPLICA IDENTITY, got %d; all ops: %v", len(alters.Ops), plan.Ops)
		}
		sql := alters.Ops[0].SQL
		if !strings.Contains(sql, "REPLICA IDENTITY") {
			t.Errorf("expected REPLICA IDENTITY in DDL, got: %s", sql)
		}
		if !strings.Contains(sql, "ALTER TABLE") {
			t.Errorf("expected ALTER TABLE, got: %s", sql)
		}
	})

	t.Run("CREATE VIEW WITH CHECK OPTION", func(t *testing.T) {
		from, err := LoadSQL("")
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(`
			CREATE TABLE t1 (id int, val int);
			CREATE VIEW v1 AS SELECT id, val FROM t1 WHERE val > 0
				WITH CASCADED CHECK OPTION;
		`)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := GenerateMigration(from, to, diff)
		viewOps := plan.Filter(func(op MigrationOp) bool {
			return op.Type == OpCreateView
		})
		if len(viewOps.Ops) != 1 {
			t.Fatalf("expected 1 CreateView op, got %d; all ops: %v", len(viewOps.Ops), plan.Ops)
		}
		sql := viewOps.Ops[0].SQL
		if !strings.Contains(sql, "CREATE VIEW") {
			t.Errorf("expected CREATE VIEW, got: %s", sql)
		}
		if !strings.Contains(sql, "CHECK OPTION") {
			t.Errorf("expected CHECK OPTION in DDL, got: %s", sql)
		}
	})

	t.Run("table inheritance INHERITS clause in CREATE TABLE", func(t *testing.T) {
		from, err := LoadSQL("")
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(`
			CREATE TABLE parent (id int, name text);
			CREATE TABLE child (extra int) INHERITS (parent);
		`)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := GenerateMigration(from, to, diff)
		creates := plan.Filter(func(op MigrationOp) bool {
			return op.Type == OpCreateTable
		})
		// Should have 2 CreateTable ops: parent and child.
		if len(creates.Ops) < 2 {
			t.Fatalf("expected at least 2 CreateTable ops, got %d; all ops: %v", len(creates.Ops), plan.Ops)
		}
		// Find the child table op.
		var childSQL string
		for _, op := range creates.Ops {
			if op.ObjectName == "child" {
				childSQL = op.SQL
				break
			}
		}
		if childSQL == "" {
			t.Fatalf("expected CreateTable op for child table")
		}
		if !strings.Contains(childSQL, "INHERITS") {
			t.Errorf("expected INHERITS clause in child DDL, got: %s", childSQL)
		}
		if !strings.Contains(childSQL, `"parent"`) {
			t.Errorf("expected parent table reference in INHERITS, got: %s", childSQL)
		}
	})
}
