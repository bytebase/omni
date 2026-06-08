package catalog

import (
	"strings"
	"testing"
)

func TestMigrationTable(t *testing.T) {
	t.Run("CREATE TABLE with columns and inline PK/UNIQUE/CHECK", func(t *testing.T) {
		from, err := LoadSQL("")
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(`
			CREATE TABLE t (
				id integer PRIMARY KEY,
				name text UNIQUE,
				val integer CHECK (val > 0)
			);
		`)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := GenerateMigration(from, to, diff)

		// Should have exactly 1 CreateTable op.
		creates := plan.Filter(func(op MigrationOp) bool {
			return op.Type == OpCreateTable
		})
		if len(creates.Ops) != 1 {
			t.Fatalf("expected 1 CreateTable op, got %d; all ops: %v", len(creates.Ops), plan.Ops)
		}
		sql := creates.Ops[0].SQL
		if !strings.Contains(sql, "CREATE TABLE") {
			t.Errorf("expected CREATE TABLE, got: %s", sql)
		}
		if !strings.Contains(sql, "PRIMARY KEY") {
			t.Errorf("expected PRIMARY KEY constraint, got: %s", sql)
		}
		if !strings.Contains(sql, "UNIQUE") {
			t.Errorf("expected UNIQUE constraint, got: %s", sql)
		}
		if !strings.Contains(sql, "CHECK") {
			t.Errorf("expected CHECK constraint, got: %s", sql)
		}
		// Should NOT contain FOREIGN KEY.
		if strings.Contains(sql, "FOREIGN KEY") {
			t.Errorf("should not contain FOREIGN KEY in CREATE TABLE, got: %s", sql)
		}
	})

	t.Run("CREATE TABLE with column defaults", func(t *testing.T) {
		from, err := LoadSQL("")
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(`
			CREATE TABLE t (
				id integer,
				status text DEFAULT 'active',
				created_at timestamp DEFAULT now()
			);
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
			t.Fatalf("expected 1 CreateTable op, got %d", len(creates.Ops))
		}
		sql := creates.Ops[0].SQL
		if !strings.Contains(sql, "DEFAULT") {
			t.Errorf("expected DEFAULT in DDL, got: %s", sql)
		}
	})

	t.Run("CREATE TABLE with NOT NULL columns", func(t *testing.T) {
		from, err := LoadSQL("")
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(`
			CREATE TABLE t (
				id integer NOT NULL,
				name text NOT NULL,
				optional text
			);
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
			t.Fatalf("expected 1 CreateTable op, got %d", len(creates.Ops))
		}
		sql := creates.Ops[0].SQL
		// Count NOT NULL occurrences — should be exactly 2 (for id and name, not optional).
		count := strings.Count(sql, "NOT NULL")
		if count != 2 {
			t.Errorf("expected 2 NOT NULL in DDL, got %d; SQL: %s", count, sql)
		}
	})

	t.Run("CREATE TABLE with column identity GENERATED ALWAYS AS IDENTITY", func(t *testing.T) {
		from, err := LoadSQL("")
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(`
			CREATE TABLE t (
				id integer GENERATED ALWAYS AS IDENTITY,
				name text
			);
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
			t.Fatalf("expected 1 CreateTable op, got %d", len(creates.Ops))
		}
		sql := creates.Ops[0].SQL
		if !strings.Contains(sql, "GENERATED ALWAYS AS IDENTITY") {
			t.Errorf("expected GENERATED ALWAYS AS IDENTITY, got: %s", sql)
		}
	})

	t.Run("CREATE TABLE with generated column GENERATED ALWAYS AS ... STORED", func(t *testing.T) {
		from, err := LoadSQL("")
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(`
			CREATE TABLE t (
				price numeric,
				tax numeric,
				total numeric GENERATED ALWAYS AS (price + tax) STORED
			);
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
			t.Fatalf("expected 1 CreateTable op, got %d", len(creates.Ops))
		}
		sql := creates.Ops[0].SQL
		if !strings.Contains(sql, "GENERATED ALWAYS AS") {
			t.Errorf("expected GENERATED ALWAYS AS, got: %s", sql)
		}
		if !strings.Contains(sql, "STORED") {
			t.Errorf("expected STORED, got: %s", sql)
		}
	})

	t.Run("CREATE UNLOGGED TABLE for unlogged tables", func(t *testing.T) {
		from, err := LoadSQL("")
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(`CREATE UNLOGGED TABLE t (id integer);`)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := GenerateMigration(from, to, diff)
		creates := plan.Filter(func(op MigrationOp) bool {
			return op.Type == OpCreateTable
		})
		if len(creates.Ops) != 1 {
			t.Fatalf("expected 1 CreateTable op, got %d", len(creates.Ops))
		}
		sql := creates.Ops[0].SQL
		if !strings.Contains(sql, "CREATE UNLOGGED TABLE") {
			t.Errorf("expected CREATE UNLOGGED TABLE, got: %s", sql)
		}
	})

	t.Run("DROP TABLE CASCADE for dropped tables", func(t *testing.T) {
		from, err := LoadSQL(`CREATE TABLE t (id integer);`)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL("")
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := GenerateMigration(from, to, diff)
		drops := plan.Filter(func(op MigrationOp) bool {
			return op.Type == OpDropTable
		})
		if len(drops.Ops) != 1 {
			t.Fatalf("expected 1 DropTable op, got %d", len(drops.Ops))
		}
		sql := drops.Ops[0].SQL
		if !strings.Contains(sql, "DROP TABLE") {
			t.Errorf("expected DROP TABLE, got: %s", sql)
		}
		if !strings.Contains(sql, "CASCADE") {
			t.Errorf("expected CASCADE, got: %s", sql)
		}
	})

	t.Run("all identifiers double-quoted in generated DDL", func(t *testing.T) {
		from, err := LoadSQL("")
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(`CREATE TABLE t (id integer, name text);`)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := GenerateMigration(from, to, diff)
		creates := plan.Filter(func(op MigrationOp) bool {
			return op.Type == OpCreateTable
		})
		if len(creates.Ops) != 1 {
			t.Fatalf("expected 1 CreateTable op, got %d", len(creates.Ops))
		}
		sql := creates.Ops[0].SQL
		// Table name should be quoted.
		if !strings.Contains(sql, `"t"`) {
			t.Errorf("expected table name to be double-quoted, got: %s", sql)
		}
		// Column names should be quoted.
		if !strings.Contains(sql, `"id"`) {
			t.Errorf("expected column name 'id' to be double-quoted, got: %s", sql)
		}
		if !strings.Contains(sql, `"name"`) {
			t.Errorf("expected column name 'name' to be double-quoted, got: %s", sql)
		}
	})

	t.Run("schema-qualified table names in DDL", func(t *testing.T) {
		from, err := LoadSQL("")
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(`
			CREATE SCHEMA myschema;
			CREATE TABLE myschema.t (id integer);
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
			t.Fatalf("expected 1 CreateTable op, got %d", len(creates.Ops))
		}
		sql := creates.Ops[0].SQL
		// Should contain schema-qualified name.
		if !strings.Contains(sql, `"myschema"."t"`) {
			t.Errorf("expected schema-qualified name \"myschema\".\"t\", got: %s", sql)
		}
	})
}
