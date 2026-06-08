package catalog

import (
	"strings"
	"testing"
)

func TestMigrationGrant(t *testing.T) {
	t.Run("GRANT privilege ON object TO role for added grant", func(t *testing.T) {
		fromSQL := `CREATE TABLE t (id int);`
		toSQL := `
			CREATE TABLE t (id int);
			GRANT SELECT ON t TO some_role;
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
			if op.Type == OpGrant {
				found = true
				if !strings.Contains(op.SQL, "GRANT SELECT") {
					t.Errorf("expected GRANT SELECT, got: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "TO") {
					t.Errorf("expected TO clause, got: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "some_role") {
					t.Errorf("expected role name, got: %s", op.SQL)
				}
			}
		}
		if !found {
			t.Errorf("no GRANT op found; ops: %v", opsSQL(plan))
		}
	})

	t.Run("REVOKE privilege ON object FROM role for removed grant", func(t *testing.T) {
		fromSQL := `
			CREATE TABLE t (id int);
			GRANT SELECT ON t TO some_role;
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
			if op.Type == OpRevoke {
				found = true
				if !strings.Contains(op.SQL, "REVOKE SELECT") {
					t.Errorf("expected REVOKE SELECT, got: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "FROM") {
					t.Errorf("expected FROM clause, got: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "some_role") {
					t.Errorf("expected role name, got: %s", op.SQL)
				}
			}
		}
		if !found {
			t.Errorf("no REVOKE op found; ops: %v", opsSQL(plan))
		}
	})

	t.Run("GRANT WITH GRANT OPTION", func(t *testing.T) {
		fromSQL := `CREATE TABLE t (id int);`
		toSQL := `
			CREATE TABLE t (id int);
			GRANT SELECT ON t TO some_role WITH GRANT OPTION;
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
			if op.Type == OpGrant {
				found = true
				if !strings.Contains(op.SQL, "WITH GRANT OPTION") {
					t.Errorf("expected WITH GRANT OPTION, got: %s", op.SQL)
				}
			}
		}
		if !found {
			t.Errorf("no GRANT op found; ops: %v", opsSQL(plan))
		}
	})

	t.Run("modified grant as REVOKE plus GRANT", func(t *testing.T) {
		fromSQL := `
			CREATE TABLE t (id int);
			GRANT SELECT ON t TO some_role;
		`
		toSQL := `
			CREATE TABLE t (id int);
			GRANT SELECT ON t TO some_role WITH GRANT OPTION;
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
		foundRevoke := false
		foundGrant := false
		for _, op := range plan.Ops {
			if op.Type == OpRevoke && strings.Contains(op.SQL, "REVOKE SELECT") {
				foundRevoke = true
			}
			if op.Type == OpGrant && strings.Contains(op.SQL, "GRANT SELECT") {
				foundGrant = true
				if !strings.Contains(op.SQL, "WITH GRANT OPTION") {
					t.Errorf("expected WITH GRANT OPTION in new grant, got: %s", op.SQL)
				}
			}
		}
		if !foundRevoke {
			t.Errorf("no REVOKE op found for modified grant; ops: %v", opsSQL(plan))
		}
		if !foundGrant {
			t.Errorf("no GRANT op found for modified grant; ops: %v", opsSQL(plan))
		}
	})

	t.Run("grants generated after object creation", func(t *testing.T) {
		fromSQL := ``
		toSQL := `
			CREATE TABLE t (id int);
			GRANT SELECT ON t TO some_role;
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
		createIdx := -1
		grantIdx := -1
		for i, op := range plan.Ops {
			if op.Type == OpCreateTable {
				createIdx = i
			}
			if op.Type == OpGrant {
				grantIdx = i
			}
		}
		if createIdx < 0 {
			t.Fatalf("no CREATE TABLE op found; ops: %v", opsSQL(plan))
		}
		if grantIdx < 0 {
			t.Fatalf("no GRANT op found; ops: %v", opsSQL(plan))
		}
		if grantIdx < createIdx {
			t.Errorf("GRANT op (idx %d) should appear after CREATE TABLE (idx %d)", grantIdx, createIdx)
		}
	})
}
