package catalog

import (
	"strings"
	"testing"
)

func TestMigrationConstraint(t *testing.T) {
	t.Run("ADD CONSTRAINT PRIMARY KEY", func(t *testing.T) {
		fromSQL := `CREATE TABLE t (id int, name text);`
		toSQL := `CREATE TABLE t (id int PRIMARY KEY, name text);`
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
			if op.Type == OpAddConstraint && strings.Contains(op.SQL, "PRIMARY KEY") {
				found = true
				if !strings.Contains(op.SQL, "ALTER TABLE") {
					t.Errorf("expected ALTER TABLE, got: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "ADD CONSTRAINT") {
					t.Errorf("expected ADD CONSTRAINT, got: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "(\"id\")") && !strings.Contains(op.SQL, "(id)") {
					t.Errorf("expected column id in PRIMARY KEY, got: %s", op.SQL)
				}
			}
		}
		if !found {
			t.Errorf("no ADD CONSTRAINT PRIMARY KEY op found; ops: %v", opsSQL(plan))
		}
	})

	t.Run("ADD CONSTRAINT UNIQUE", func(t *testing.T) {
		fromSQL := `CREATE TABLE t (id int, email text);`
		toSQL := `CREATE TABLE t (id int, email text, CONSTRAINT t_email_key UNIQUE (email));`
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
			if op.Type == OpAddConstraint && strings.Contains(op.SQL, "UNIQUE") {
				found = true
				if !strings.Contains(op.SQL, "ADD CONSTRAINT") {
					t.Errorf("expected ADD CONSTRAINT, got: %s", op.SQL)
				}
			}
		}
		if !found {
			t.Errorf("no ADD CONSTRAINT UNIQUE op found; ops: %v", opsSQL(plan))
		}
	})

	t.Run("ADD CONSTRAINT CHECK", func(t *testing.T) {
		fromSQL := `CREATE TABLE t (id int, age int);`
		toSQL := `CREATE TABLE t (id int, age int, CONSTRAINT t_age_check CHECK (age > 0));`
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
			if op.Type == OpAddConstraint && strings.Contains(op.SQL, "CHECK") {
				found = true
				if !strings.Contains(op.SQL, "age > 0") && !strings.Contains(op.SQL, "(age > 0)") {
					t.Errorf("expected CHECK expression, got: %s", op.SQL)
				}
			}
		}
		if !found {
			t.Errorf("no ADD CONSTRAINT CHECK op found; ops: %v", opsSQL(plan))
		}
	})

	t.Run("ADD CONSTRAINT FOREIGN KEY", func(t *testing.T) {
		fromSQL := `
			CREATE TABLE parent (id int PRIMARY KEY);
			CREATE TABLE child (id int, parent_id int);
		`
		toSQL := `
			CREATE TABLE parent (id int PRIMARY KEY);
			CREATE TABLE child (id int, parent_id int,
				CONSTRAINT child_parent_id_fkey FOREIGN KEY (parent_id) REFERENCES parent (id));
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
			if op.Type == OpAddConstraint && strings.Contains(op.SQL, "FOREIGN KEY") {
				found = true
				if !strings.Contains(op.SQL, "REFERENCES") {
					t.Errorf("expected REFERENCES, got: %s", op.SQL)
				}
			}
		}
		if !found {
			t.Errorf("no ADD CONSTRAINT FOREIGN KEY op found; ops: %v", opsSQL(plan))
		}
	})

	t.Run("FK with ON DELETE/UPDATE actions", func(t *testing.T) {
		fromSQL := `
			CREATE TABLE parent (id int PRIMARY KEY);
			CREATE TABLE child (id int, parent_id int);
		`
		toSQL := `
			CREATE TABLE parent (id int PRIMARY KEY);
			CREATE TABLE child (id int, parent_id int,
				CONSTRAINT child_parent_fk FOREIGN KEY (parent_id) REFERENCES parent (id)
				ON DELETE CASCADE ON UPDATE SET NULL);
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
			if op.Type == OpAddConstraint && strings.Contains(op.SQL, "FOREIGN KEY") {
				found = true
				if !strings.Contains(op.SQL, "ON DELETE CASCADE") {
					t.Errorf("expected ON DELETE CASCADE, got: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "ON UPDATE SET NULL") {
					t.Errorf("expected ON UPDATE SET NULL, got: %s", op.SQL)
				}
			}
		}
		if !found {
			t.Errorf("no FK with actions op found; ops: %v", opsSQL(plan))
		}
	})

	t.Run("FK with DEFERRABLE INITIALLY DEFERRED", func(t *testing.T) {
		fromSQL := `
			CREATE TABLE parent (id int PRIMARY KEY);
			CREATE TABLE child (id int, parent_id int);
		`
		toSQL := `
			CREATE TABLE parent (id int PRIMARY KEY);
			CREATE TABLE child (id int, parent_id int,
				CONSTRAINT child_parent_fk FOREIGN KEY (parent_id) REFERENCES parent (id)
				DEFERRABLE INITIALLY DEFERRED);
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
			if op.Type == OpAddConstraint && strings.Contains(op.SQL, "FOREIGN KEY") {
				found = true
				if !strings.Contains(op.SQL, "DEFERRABLE INITIALLY DEFERRED") {
					t.Errorf("expected DEFERRABLE INITIALLY DEFERRED, got: %s", op.SQL)
				}
			}
		}
		if !found {
			t.Errorf("no FK with DEFERRABLE op found; ops: %v", opsSQL(plan))
		}
	})

	t.Run("Constraint with NOT VALID", func(t *testing.T) {
		fromSQL := `
			CREATE TABLE parent (id int PRIMARY KEY);
			CREATE TABLE child (id int, parent_id int);
		`
		toSQL := `
			CREATE TABLE parent (id int PRIMARY KEY);
			CREATE TABLE child (id int, parent_id int,
				CONSTRAINT child_parent_fk FOREIGN KEY (parent_id) REFERENCES parent (id)
				NOT VALID);
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
			if op.Type == OpAddConstraint && strings.Contains(op.SQL, "FOREIGN KEY") {
				found = true
				if !strings.Contains(op.SQL, "NOT VALID") {
					t.Errorf("expected NOT VALID, got: %s", op.SQL)
				}
			}
		}
		if !found {
			t.Errorf("no NOT VALID constraint op found; ops: %v", opsSQL(plan))
		}
	})

	t.Run("DROP CONSTRAINT for removed constraints", func(t *testing.T) {
		fromSQL := `CREATE TABLE t (id int, age int, CONSTRAINT t_age_check CHECK (age > 0));`
		toSQL := `CREATE TABLE t (id int, age int);`
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
			if op.Type == OpDropConstraint && strings.Contains(op.SQL, "DROP CONSTRAINT") {
				found = true
				if !strings.Contains(op.SQL, "t_age_check") {
					t.Errorf("expected constraint name t_age_check, got: %s", op.SQL)
				}
			}
		}
		if !found {
			t.Errorf("no DROP CONSTRAINT op found; ops: %v", opsSQL(plan))
		}
	})

	t.Run("Modified constraint as DROP + ADD", func(t *testing.T) {
		fromSQL := `CREATE TABLE t (id int, age int, CONSTRAINT t_age_check CHECK (age > 0));`
		toSQL := `CREATE TABLE t (id int, age int, CONSTRAINT t_age_check CHECK (age > 10));`
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
		var dropFound, addFound bool
		for _, op := range plan.Ops {
			if op.Type == OpDropConstraint && strings.Contains(op.SQL, "t_age_check") {
				dropFound = true
			}
			if op.Type == OpAddConstraint && strings.Contains(op.SQL, "t_age_check") && strings.Contains(op.SQL, "age > 10") {
				addFound = true
			}
		}
		if !dropFound {
			t.Errorf("no DROP for modified constraint; ops: %v", opsSQL(plan))
		}
		if !addFound {
			t.Errorf("no ADD for modified constraint; ops: %v", opsSQL(plan))
		}
	})

	t.Run("EXCLUDE constraint ADD/DROP", func(t *testing.T) {
		fromSQL := `CREATE TABLE t (id int, r int4range, CONSTRAINT t_r_excl EXCLUDE USING gist (r WITH &&));`
		toSQL := `CREATE TABLE t (id int, r int4range);`
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
		dropFound := false
		for _, op := range plan.Ops {
			if op.Type == OpDropConstraint && strings.Contains(op.SQL, "t_r_excl") {
				dropFound = true
			}
		}
		if !dropFound {
			t.Errorf("no DROP for EXCLUDE constraint; ops: %v", opsSQL(plan))
		}

		// Now test ADD direction.
		diff2 := Diff(to, from)
		plan2 := GenerateMigration(to, from, diff2)
		addFound := false
		for _, op := range plan2.Ops {
			if op.Type == OpAddConstraint && strings.Contains(op.SQL, "EXCLUDE") {
				addFound = true
				if !strings.Contains(op.SQL, "t_r_excl") {
					t.Errorf("expected constraint name t_r_excl in EXCLUDE, got: %s", op.SQL)
				}
			}
		}
		if !addFound {
			t.Errorf("no ADD for EXCLUDE constraint; ops: %v", opsSQL(plan2))
		}
	})
}

// opsSQL is a test helper that returns all SQL statements from a plan.
func opsSQL(plan *MigrationPlan) []string {
	var ss []string
	for _, op := range plan.Ops {
		ss = append(ss, op.SQL)
	}
	return ss
}
