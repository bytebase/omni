package catalog

import (
	"strings"
	"testing"
)

func TestMigrationTrigger(t *testing.T) {
	t.Run("CREATE TRIGGER with timing events level function", func(t *testing.T) {
		fromSQL := `
			CREATE TABLE t (id int);
			CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;
		`
		toSQL := `
			CREATE TABLE t (id int);
			CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;
			CREATE TRIGGER my_trig BEFORE INSERT ON t FOR EACH ROW EXECUTE FUNCTION trig_fn();
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
			if op.Type == OpCreateTrigger {
				found = true
				if !strings.Contains(op.SQL, "CREATE TRIGGER") {
					t.Errorf("expected CREATE TRIGGER, got: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "BEFORE") {
					t.Errorf("expected BEFORE, got: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "INSERT") {
					t.Errorf("expected INSERT, got: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "FOR EACH ROW") {
					t.Errorf("expected FOR EACH ROW, got: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "EXECUTE FUNCTION") {
					t.Errorf("expected EXECUTE FUNCTION, got: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "trig_fn") {
					t.Errorf("expected trig_fn, got: %s", op.SQL)
				}
			}
		}
		if !found {
			t.Errorf("no CREATE TRIGGER op found; ops: %v", opsSQL(plan))
		}
	})

	t.Run("CREATE TRIGGER with WHEN clause", func(t *testing.T) {
		fromSQL := `
			CREATE TABLE t (id int, active boolean);
			CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;
		`
		toSQL := `
			CREATE TABLE t (id int, active boolean);
			CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;
			CREATE TRIGGER my_trig BEFORE INSERT ON t FOR EACH ROW WHEN (NEW.active IS TRUE) EXECUTE FUNCTION trig_fn();
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
			if op.Type == OpCreateTrigger {
				found = true
				if !strings.Contains(op.SQL, "WHEN (") {
					t.Errorf("expected WHEN clause, got: %s", op.SQL)
				}
			}
		}
		if !found {
			t.Errorf("no CREATE TRIGGER op found; ops: %v", opsSQL(plan))
		}
	})

	t.Run("CREATE TRIGGER with REFERENCING OLD NEW TABLE AS", func(t *testing.T) {
		fromSQL := `
			CREATE TABLE t (id int, val text);
			CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;
		`
		toSQL := `
			CREATE TABLE t (id int, val text);
			CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;
			CREATE TRIGGER my_trig AFTER UPDATE ON t REFERENCING OLD TABLE AS old_t NEW TABLE AS new_t FOR EACH ROW EXECUTE FUNCTION trig_fn();
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
			if op.Type == OpCreateTrigger {
				found = true
				if !strings.Contains(op.SQL, "REFERENCING") {
					t.Errorf("expected REFERENCING, got: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "OLD TABLE AS") {
					t.Errorf("expected OLD TABLE AS, got: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "NEW TABLE AS") {
					t.Errorf("expected NEW TABLE AS, got: %s", op.SQL)
				}
			}
		}
		if !found {
			t.Errorf("no CREATE TRIGGER op found; ops: %v", opsSQL(plan))
		}
	})

	t.Run("DROP TRIGGER ON table for removed trigger", func(t *testing.T) {
		fromSQL := `
			CREATE TABLE t (id int);
			CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;
			CREATE TRIGGER my_trig BEFORE INSERT ON t FOR EACH ROW EXECUTE FUNCTION trig_fn();
		`
		toSQL := `
			CREATE TABLE t (id int);
			CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;
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
			if op.Type == OpDropTrigger {
				found = true
				if !strings.Contains(op.SQL, "DROP TRIGGER") {
					t.Errorf("expected DROP TRIGGER, got: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "ON") {
					t.Errorf("expected ON table clause, got: %s", op.SQL)
				}
			}
		}
		if !found {
			t.Errorf("no DROP TRIGGER op found; ops: %v", opsSQL(plan))
		}
	})

	t.Run("CREATE TRIGGER with UPDATE OF columns", func(t *testing.T) {
		fromSQL := `
			CREATE TABLE t (id int, name text, val text);
			CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;
		`
		toSQL := `
			CREATE TABLE t (id int, name text, val text);
			CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;
			CREATE TRIGGER my_trig BEFORE UPDATE OF name, val ON t FOR EACH ROW EXECUTE FUNCTION trig_fn();
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
			if op.Type == OpCreateTrigger {
				found = true
				if !strings.Contains(op.SQL, "UPDATE OF") {
					t.Errorf("expected UPDATE OF, got: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, `"name"`) {
					t.Errorf("expected column name in UPDATE OF, got: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, `"val"`) {
					t.Errorf("expected column val in UPDATE OF, got: %s", op.SQL)
				}
			}
		}
		if !found {
			t.Errorf("no CREATE TRIGGER op found; ops: %v", opsSQL(plan))
		}
	})

	t.Run("modified trigger as DROP plus CREATE", func(t *testing.T) {
		fromSQL := `
			CREATE TABLE t (id int);
			CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;
			CREATE TRIGGER my_trig BEFORE INSERT ON t FOR EACH ROW EXECUTE FUNCTION trig_fn();
		`
		toSQL := `
			CREATE TABLE t (id int);
			CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;
			CREATE TRIGGER my_trig AFTER INSERT ON t FOR EACH STATEMENT EXECUTE FUNCTION trig_fn();
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
			if op.Type == OpDropTrigger && strings.Contains(op.SQL, "my_trig") {
				foundDrop = true
			}
			if op.Type == OpCreateTrigger && strings.Contains(op.SQL, "my_trig") {
				foundCreate = true
				if !strings.Contains(op.SQL, "AFTER") {
					t.Errorf("expected AFTER in recreated trigger, got: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "FOR EACH STATEMENT") {
					t.Errorf("expected FOR EACH STATEMENT in recreated trigger, got: %s", op.SQL)
				}
			}
		}
		if !foundDrop {
			t.Errorf("no DROP TRIGGER op found for modified trigger; ops: %v", opsSQL(plan))
		}
		if !foundCreate {
			t.Errorf("no CREATE TRIGGER op found for modified trigger; ops: %v", opsSQL(plan))
		}
	})
}
