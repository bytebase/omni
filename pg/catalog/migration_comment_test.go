package catalog

import (
	"strings"
	"testing"
)

func TestMigrationComment(t *testing.T) {
	t.Run("COMMENT ON TABLE for added comment", func(t *testing.T) {
		fromSQL := `CREATE TABLE t (id int);`
		toSQL := `
			CREATE TABLE t (id int);
			COMMENT ON TABLE t IS 'This is a table';
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
			if op.Type == OpComment {
				found = true
				if !strings.Contains(op.SQL, "COMMENT ON TABLE") {
					t.Errorf("expected COMMENT ON TABLE, got: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "This is a table") {
					t.Errorf("expected comment text, got: %s", op.SQL)
				}
			}
		}
		if !found {
			t.Errorf("no COMMENT op found; ops: %v", opsSQL(plan))
		}
	})

	t.Run("COMMENT ON TABLE IS NULL for removed comment", func(t *testing.T) {
		fromSQL := `
			CREATE TABLE t (id int);
			COMMENT ON TABLE t IS 'This is a table';
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
			if op.Type == OpComment {
				found = true
				if !strings.Contains(op.SQL, "IS NULL") {
					t.Errorf("expected IS NULL for removed comment, got: %s", op.SQL)
				}
			}
		}
		if !found {
			t.Errorf("no COMMENT op found for removed comment; ops: %v", opsSQL(plan))
		}
	})

	t.Run("COMMENT ON COLUMN for column comments", func(t *testing.T) {
		fromSQL := `CREATE TABLE t (id int, name text);`
		toSQL := `
			CREATE TABLE t (id int, name text);
			COMMENT ON COLUMN t.name IS 'The name column';
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
			if op.Type == OpComment && strings.Contains(op.SQL, "COLUMN") {
				found = true
				if !strings.Contains(op.SQL, "The name column") {
					t.Errorf("expected column comment text, got: %s", op.SQL)
				}
			}
		}
		if !found {
			t.Errorf("no COMMENT ON COLUMN op found; ops: %v", opsSQL(plan))
		}
	})

	t.Run("COMMENT ON INDEX FUNCTION SCHEMA TYPE SEQUENCE", func(t *testing.T) {
		fromSQL := `
			CREATE TABLE t (id int);
			CREATE INDEX idx_t ON t (id);
			CREATE FUNCTION fn() RETURNS void LANGUAGE sql AS $$ SELECT 1 $$;
			CREATE TYPE mood AS ENUM ('happy', 'sad');
			CREATE SEQUENCE myseq;
		`
		toSQL := `
			CREATE TABLE t (id int);
			CREATE INDEX idx_t ON t (id);
			COMMENT ON INDEX idx_t IS 'index comment';
			CREATE FUNCTION fn() RETURNS void LANGUAGE sql AS $$ SELECT 1 $$;
			COMMENT ON FUNCTION fn() IS 'function comment';
			COMMENT ON SCHEMA public IS 'default schema';
			CREATE TYPE mood AS ENUM ('happy', 'sad');
			COMMENT ON TYPE mood IS 'mood enum';
			CREATE SEQUENCE myseq;
			COMMENT ON SEQUENCE myseq IS 'a sequence';
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
		foundIndex := false
		foundFunction := false
		foundSchema := false
		foundType := false
		foundSequence := false
		for _, op := range plan.Ops {
			if op.Type == OpComment {
				if strings.Contains(op.SQL, "COMMENT ON INDEX") {
					foundIndex = true
				}
				if strings.Contains(op.SQL, "COMMENT ON FUNCTION") {
					foundFunction = true
				}
				if strings.Contains(op.SQL, "COMMENT ON SCHEMA") {
					foundSchema = true
				}
				if strings.Contains(op.SQL, "COMMENT ON TYPE") {
					foundType = true
				}
				if strings.Contains(op.SQL, "COMMENT ON SEQUENCE") {
					foundSequence = true
				}
			}
		}
		if !foundIndex {
			t.Errorf("no COMMENT ON INDEX found; ops: %v", opsSQL(plan))
		}
		if !foundFunction {
			t.Errorf("no COMMENT ON FUNCTION found; ops: %v", opsSQL(plan))
		}
		if !foundSchema {
			t.Errorf("no COMMENT ON SCHEMA found; ops: %v", opsSQL(plan))
		}
		if !foundType {
			t.Errorf("no COMMENT ON TYPE found; ops: %v", opsSQL(plan))
		}
		if !foundSequence {
			t.Errorf("no COMMENT ON SEQUENCE found; ops: %v", opsSQL(plan))
		}
	})

	t.Run("comments generated after object creation", func(t *testing.T) {
		fromSQL := ``
		toSQL := `
			CREATE TABLE t (id int);
			COMMENT ON TABLE t IS 'new table';
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
		// The comment op must appear after the CREATE TABLE op.
		createIdx := -1
		commentIdx := -1
		for i, op := range plan.Ops {
			if op.Type == OpCreateTable {
				createIdx = i
			}
			if op.Type == OpComment {
				commentIdx = i
			}
		}
		if createIdx < 0 {
			t.Fatalf("no CREATE TABLE op found; ops: %v", opsSQL(plan))
		}
		if commentIdx < 0 {
			t.Fatalf("no COMMENT op found; ops: %v", opsSQL(plan))
		}
		if commentIdx < createIdx {
			t.Errorf("COMMENT op (idx %d) should appear after CREATE TABLE (idx %d)", commentIdx, createIdx)
		}
	})

	t.Run("COMMENT ON VIEW uses VIEW not TABLE", func(t *testing.T) {
		fromSQL := `CREATE VIEW v AS SELECT 1 AS id;`
		toSQL := `
			CREATE VIEW v AS SELECT 1 AS id;
			COMMENT ON VIEW v IS 'This is a view';
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
			if op.Type == OpComment {
				found = true
				if !strings.Contains(op.SQL, "COMMENT ON VIEW") {
					t.Errorf("expected COMMENT ON VIEW, got: %s", op.SQL)
				}
				if strings.Contains(op.SQL, "COMMENT ON TABLE") {
					t.Errorf("should not use TABLE for view comment: %s", op.SQL)
				}
			}
		}
		if !found {
			t.Errorf("no COMMENT op found; ops: %v", opsSQL(plan))
		}
	})

	t.Run("COMMENT ON MATERIALIZED VIEW uses MATERIALIZED VIEW", func(t *testing.T) {
		fromSQL := `CREATE MATERIALIZED VIEW mv AS SELECT 1 AS id;`
		toSQL := `
			CREATE MATERIALIZED VIEW mv AS SELECT 1 AS id;
			COMMENT ON MATERIALIZED VIEW mv IS 'This is a matview';
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
			if op.Type == OpComment {
				found = true
				if !strings.Contains(op.SQL, "COMMENT ON MATERIALIZED VIEW") {
					t.Errorf("expected COMMENT ON MATERIALIZED VIEW, got: %s", op.SQL)
				}
			}
		}
		if !found {
			t.Errorf("no COMMENT op found; ops: %v", opsSQL(plan))
		}
	})

	t.Run("single quotes in a comment are doubled", func(t *testing.T) {
		fromSQL := `CREATE TABLE t (id int);`
		toSQL := `
			CREATE TABLE t (id int);
			COMMENT ON TABLE t IS 'User''s table';
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
		ops := plan.Filter(func(op MigrationOp) bool {
			return op.Type == OpComment
		}).Ops
		if len(ops) != 1 {
			t.Fatalf("expected 1 Comment op, got %d; ops: %v", len(ops), opsSQL(plan))
		}
		// The catalog holds the comment unescaped. Rendering it back with a bare
		// quote closes the literal early and the statement no longer parses, so
		// the whole literal has to be checked, not just the text around it.
		if !strings.Contains(ops[0].SQL, `IS 'User''s table'`) {
			t.Errorf("expected the quote doubled in the emitted literal, got: %s", ops[0].SQL)
		}
	})

	t.Run("dropping a procedure leaves an unrelated function's comment alone", func(t *testing.T) {
		const survivor = `
			CREATE FUNCTION audit_log() RETURNS void
			LANGUAGE plpgsql AS $$ BEGIN NULL; END; $$;
			COMMENT ON FUNCTION audit_log() IS 'Writes the audit log';
		`
		fromSQL := survivor + `
			CREATE PROCEDURE purge_rows() LANGUAGE plpgsql AS $$ BEGIN NULL; END; $$;
			COMMENT ON PROCEDURE purge_rows() IS 'Purges old rows';
		`
		toSQL := survivor
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
		// Procedures and functions share one catalog, so a comment differ that
		// keys on the bare name rather than the full signature reports the
		// dropped procedure's comment as a change to the surviving function.
		if ops := filterOps(plan, OpDropFunction); len(ops) != 1 {
			t.Fatalf("expected 1 DropFunction op for the procedure, got %d; ops: %v", len(ops), opsSQL(plan))
		}
		for _, op := range plan.Ops {
			if op.Type == OpComment && strings.Contains(op.SQL, "audit_log") {
				t.Errorf("the surviving function's comment must not change: %s", op.SQL)
			}
		}
	})

	t.Run("COMMENT ON PROCEDURE uses PROCEDURE not FUNCTION", func(t *testing.T) {
		fromSQL := `
			CREATE PROCEDURE do_work(x integer)
			LANGUAGE plpgsql
			AS $$ BEGIN NULL; END; $$;
		`
		toSQL := `
			CREATE PROCEDURE do_work(x integer)
			LANGUAGE plpgsql
			AS $$ BEGIN NULL; END; $$;
			COMMENT ON PROCEDURE do_work(integer) IS 'A procedure';
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
			if op.Type == OpComment {
				found = true
				if !strings.Contains(op.SQL, "COMMENT ON PROCEDURE") {
					t.Errorf("expected COMMENT ON PROCEDURE, got: %s", op.SQL)
				}
				if strings.Contains(op.SQL, "COMMENT ON FUNCTION") {
					t.Errorf("should not use FUNCTION for procedure comment: %s", op.SQL)
				}
			}
		}
		if !found {
			t.Errorf("no COMMENT op found; ops: %v", opsSQL(plan))
		}
	})
}
