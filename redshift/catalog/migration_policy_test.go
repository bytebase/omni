package catalog

import (
	"strings"
	"testing"
)

func TestMigrationPolicy(t *testing.T) {
	t.Run("CREATE POLICY on table for added policy", func(t *testing.T) {
		fromSQL := `
			CREATE TABLE t (id int, owner text);
			ALTER TABLE t ENABLE ROW LEVEL SECURITY;
		`
		toSQL := `
			CREATE TABLE t (id int, owner text);
			ALTER TABLE t ENABLE ROW LEVEL SECURITY;
			CREATE POLICY read_own ON t FOR SELECT USING (owner = current_user);
		`
		from, err := LoadSQL(fromSQL)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(toSQL)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := GenerateMigration(from, to, diff)
		found := false
		for _, op := range plan.Ops {
			if op.Type == OpCreatePolicy {
				found = true
				if !strings.Contains(op.SQL, "CREATE POLICY") {
					t.Errorf("expected CREATE POLICY, got: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "read_own") {
					t.Errorf("expected policy name read_own, got: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "FOR SELECT") {
					t.Errorf("expected FOR SELECT, got: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "USING (") {
					t.Errorf("expected USING clause, got: %s", op.SQL)
				}
			}
		}
		if !found {
			t.Errorf("no CREATE POLICY op found; ops: %v", opsSQL(plan))
		}
	})

	t.Run("DROP POLICY on table for removed policy", func(t *testing.T) {
		fromSQL := `
			CREATE TABLE t (id int, owner text);
			ALTER TABLE t ENABLE ROW LEVEL SECURITY;
			CREATE POLICY read_own ON t FOR SELECT USING (owner = current_user);
		`
		toSQL := `
			CREATE TABLE t (id int, owner text);
			ALTER TABLE t ENABLE ROW LEVEL SECURITY;
		`
		from, err := LoadSQL(fromSQL)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(toSQL)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := GenerateMigration(from, to, diff)
		found := false
		for _, op := range plan.Ops {
			if op.Type == OpDropPolicy {
				found = true
				if !strings.Contains(op.SQL, "DROP POLICY") {
					t.Errorf("expected DROP POLICY, got: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "read_own") {
					t.Errorf("expected policy name read_own, got: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, " ON ") {
					t.Errorf("expected ON table clause, got: %s", op.SQL)
				}
			}
		}
		if !found {
			t.Errorf("no DROP POLICY op found; ops: %v", opsSQL(plan))
		}
	})

	t.Run("ALTER POLICY for simple changes roles USING WITH CHECK", func(t *testing.T) {
		fromSQL := `
			CREATE TABLE t (id int, owner text);
			ALTER TABLE t ENABLE ROW LEVEL SECURITY;
			CREATE POLICY read_own ON t FOR SELECT USING (owner = current_user);
		`
		toSQL := `
			CREATE TABLE t (id int, owner text);
			ALTER TABLE t ENABLE ROW LEVEL SECURITY;
			CREATE POLICY read_own ON t FOR SELECT USING (id > 0);
		`
		from, err := LoadSQL(fromSQL)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(toSQL)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := GenerateMigration(from, to, diff)
		found := false
		for _, op := range plan.Ops {
			if op.Type == OpAlterPolicy {
				found = true
				if !strings.Contains(op.SQL, "ALTER POLICY") {
					t.Errorf("expected ALTER POLICY, got: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "USING (") {
					t.Errorf("expected USING clause, got: %s", op.SQL)
				}
			}
		}
		if !found {
			t.Errorf("no ALTER POLICY op found; ops: %v", opsSQL(plan))
		}
	})

	t.Run("complex policy change as DROP plus CREATE", func(t *testing.T) {
		fromSQL := `
			CREATE TABLE t (id int, owner text);
			ALTER TABLE t ENABLE ROW LEVEL SECURITY;
			CREATE POLICY read_own ON t FOR SELECT USING (owner = current_user);
		`
		toSQL := `
			CREATE TABLE t (id int, owner text);
			ALTER TABLE t ENABLE ROW LEVEL SECURITY;
			CREATE POLICY read_own ON t FOR UPDATE USING (owner = current_user);
		`
		from, err := LoadSQL(fromSQL)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(toSQL)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := GenerateMigration(from, to, diff)
		foundDrop := false
		foundCreate := false
		for _, op := range plan.Ops {
			if op.Type == OpDropPolicy && strings.Contains(op.SQL, "read_own") {
				foundDrop = true
			}
			if op.Type == OpCreatePolicy && strings.Contains(op.SQL, "read_own") {
				foundCreate = true
				if !strings.Contains(op.SQL, "FOR UPDATE") {
					t.Errorf("expected FOR UPDATE in recreated policy, got: %s", op.SQL)
				}
			}
		}
		if !foundDrop {
			t.Errorf("no DROP POLICY op found for complex change; ops: %v", opsSQL(plan))
		}
		if !foundCreate {
			t.Errorf("no CREATE POLICY op found for complex change; ops: %v", opsSQL(plan))
		}
	})

	t.Run("ALTER TABLE ENABLE ROW LEVEL SECURITY", func(t *testing.T) {
		fromSQL := `CREATE TABLE t (id int);`
		toSQL := `
			CREATE TABLE t (id int);
			ALTER TABLE t ENABLE ROW LEVEL SECURITY;
		`
		from, err := LoadSQL(fromSQL)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(toSQL)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := GenerateMigration(from, to, diff)
		found := false
		for _, op := range plan.Ops {
			if op.Type == OpAlterTable && strings.Contains(op.SQL, "ENABLE ROW LEVEL SECURITY") {
				found = true
			}
		}
		if !found {
			t.Errorf("no ENABLE ROW LEVEL SECURITY op found; ops: %v", opsSQL(plan))
		}
	})

	t.Run("ALTER TABLE DISABLE ROW LEVEL SECURITY", func(t *testing.T) {
		fromSQL := `
			CREATE TABLE t (id int);
			ALTER TABLE t ENABLE ROW LEVEL SECURITY;
		`
		toSQL := `CREATE TABLE t (id int);`
		from, err := LoadSQL(fromSQL)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(toSQL)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := GenerateMigration(from, to, diff)
		found := false
		for _, op := range plan.Ops {
			if op.Type == OpAlterTable && strings.Contains(op.SQL, "DISABLE ROW LEVEL SECURITY") {
				found = true
			}
		}
		if !found {
			t.Errorf("no DISABLE ROW LEVEL SECURITY op found; ops: %v", opsSQL(plan))
		}
	})

	t.Run("ALTER TABLE FORCE ROW LEVEL SECURITY", func(t *testing.T) {
		fromSQL := `CREATE TABLE t (id int);`
		toSQL := `
			CREATE TABLE t (id int);
			ALTER TABLE t FORCE ROW LEVEL SECURITY;
		`
		from, err := LoadSQL(fromSQL)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(toSQL)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := GenerateMigration(from, to, diff)
		found := false
		for _, op := range plan.Ops {
			if op.Type == OpAlterTable && strings.Contains(op.SQL, "FORCE ROW LEVEL SECURITY") {
				found = true
				if strings.Contains(op.SQL, "NO FORCE") {
					t.Errorf("should not be NO FORCE, got: %s", op.SQL)
				}
			}
		}
		if !found {
			t.Errorf("no FORCE ROW LEVEL SECURITY op found; ops: %v", opsSQL(plan))
		}
	})

	t.Run("ALTER TABLE NO FORCE ROW LEVEL SECURITY", func(t *testing.T) {
		fromSQL := `
			CREATE TABLE t (id int);
			ALTER TABLE t FORCE ROW LEVEL SECURITY;
		`
		toSQL := `CREATE TABLE t (id int);`
		from, err := LoadSQL(fromSQL)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(toSQL)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := GenerateMigration(from, to, diff)
		found := false
		for _, op := range plan.Ops {
			if op.Type == OpAlterTable && strings.Contains(op.SQL, "NO FORCE ROW LEVEL SECURITY") {
				found = true
			}
		}
		if !found {
			t.Errorf("no NO FORCE ROW LEVEL SECURITY op found; ops: %v", opsSQL(plan))
		}
	})
}
