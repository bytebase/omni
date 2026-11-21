package catalog

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMigrationRange(t *testing.T) {
	tests := []struct {
		name    string
		fromSQL string
		toSQL   string
		check   func(t *testing.T, plan *MigrationPlan)
	}{
		{
			name:    "create range type",
			fromSQL: "",
			toSQL:   "CREATE TYPE floatrange AS RANGE (SUBTYPE = float8);",
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := plan.Filter(func(op MigrationOp) bool {
					return op.Type == OpCreateType
				})
				if len(ops.Ops) != 1 {
					t.Fatalf("expected 1 CreateType op, got %d", len(ops.Ops))
				}
				op := ops.Ops[0]
				if !strings.Contains(op.SQL, "CREATE TYPE") {
					t.Errorf("expected CREATE TYPE in SQL, got %q", op.SQL)
				}
				if !strings.Contains(op.SQL, "AS RANGE") {
					t.Errorf("expected AS RANGE in SQL, got %q", op.SQL)
				}
				if !strings.Contains(op.SQL, "SUBTYPE") {
					t.Errorf("expected SUBTYPE in SQL, got %q", op.SQL)
				}
				if !strings.Contains(op.SQL, `"floatrange"`) {
					t.Errorf("expected quoted identifier in SQL, got %q", op.SQL)
				}
				if !strings.Contains(op.SQL, `"public"`) {
					t.Errorf("expected schema-qualified name in SQL, got %q", op.SQL)
				}
				if op.SchemaName != "public" {
					t.Errorf("expected schema 'public', got %q", op.SchemaName)
				}
				if op.ObjectName != "floatrange" {
					t.Errorf("expected object name 'floatrange', got %q", op.ObjectName)
				}
				if op.Warning != "" {
					t.Errorf("unexpected warning: %q", op.Warning)
				}
			},
		},
		{
			name:    "drop range type",
			fromSQL: "CREATE TYPE floatrange AS RANGE (SUBTYPE = float8);",
			toSQL:   "",
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := plan.Filter(func(op MigrationOp) bool {
					return op.Type == OpDropType
				})
				if len(ops.Ops) != 1 {
					t.Fatalf("expected 1 DropType op, got %d", len(ops.Ops))
				}
				op := ops.Ops[0]
				if !strings.Contains(op.SQL, "DROP TYPE") {
					t.Errorf("expected DROP TYPE in SQL, got %q", op.SQL)
				}
				if !strings.Contains(op.SQL, `"floatrange"`) {
					t.Errorf("expected quoted identifier in SQL, got %q", op.SQL)
				}
				if !strings.Contains(op.SQL, `"public"`) {
					t.Errorf("expected schema-qualified name in SQL, got %q", op.SQL)
				}
				if op.SchemaName != "public" {
					t.Errorf("expected schema 'public', got %q", op.SchemaName)
				}
				if op.ObjectName != "floatrange" {
					t.Errorf("expected object name 'floatrange', got %q", op.ObjectName)
				}
				if op.Warning != "" {
					t.Errorf("unexpected warning: %q", op.Warning)
				}
			},
		},
		{
			name:    "range subtype change generates warning",
			fromSQL: "CREATE TYPE myrange AS RANGE (SUBTYPE = int4);",
			toSQL:   "CREATE TYPE myrange AS RANGE (SUBTYPE = int8);",
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := plan.Filter(func(op MigrationOp) bool {
					return op.Type == OpAlterType
				})
				if len(ops.Ops) != 1 {
					t.Fatalf("expected 1 AlterType op, got %d", len(ops.Ops))
				}
				op := ops.Ops[0]
				if op.Warning == "" {
					t.Fatal("expected warning for range subtype change")
				}
				if !strings.Contains(op.Warning, "cannot be ALTERed") {
					t.Errorf("expected warning about ALTER limitation, got %q", op.Warning)
				}
				if !strings.Contains(op.Warning, "DROP + CREATE") {
					t.Errorf("expected warning to mention DROP + CREATE, got %q", op.Warning)
				}
				if !strings.Contains(op.Warning, "dependent") {
					t.Errorf("expected warning to mention dependents, got %q", op.Warning)
				}
				if !plan.HasWarnings() {
					t.Error("expected plan.HasWarnings() to be true")
				}
				if op.SchemaName != "public" {
					t.Errorf("expected schema 'public', got %q", op.SchemaName)
				}
				if op.ObjectName != "myrange" {
					t.Errorf("expected object name 'myrange', got %q", op.ObjectName)
				}
			},
		},
		{
			name:    "create range type with user collation preserves COLLATION",
			fromSQL: "",
			toSQL:   `CREATE COLLATION "my.coll" (LOCALE = 'C'); CREATE TYPE trange AS RANGE (SUBTYPE = text, COLLATION = public."my.coll");`,
			check: func(t *testing.T, plan *MigrationPlan) {
				if !strings.Contains(plan.SQL(), `CREATE TYPE "public"."trange" AS RANGE (SUBTYPE = text, COLLATION = "public"."my.coll")`) {
					t.Fatalf("expected range DDL with COLLATION, got:\n%s", plan.SQL())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			from, err := LoadSQL(tt.fromSQL)
			if err != nil {
				t.Fatal(err)
			}
			to, err := LoadSQL(tt.toSQL)
			if err != nil {
				t.Fatal(err)
			}
			diff := Diff(from, to)
			plan := GenerateMigration(from, to, diff)
			tt.check(t, plan)
		})
	}
}

func TestRangeRenameKeepsIdentity(t *testing.T) {
	plan := func(t *testing.T, base, change string) (string, *Catalog, *Catalog) {
		t.Helper()
		from := New()
		_, err := from.Exec(base, nil)
		require.NoError(t, err)
		to := from.Clone()
		_, err = to.Exec(change, nil)
		require.NoError(t, err)
		sql := GenerateMigration(from, to, Diff(from, to)).SQL()
		applied := from.Clone()
		res, err := applied.Exec(sql, &ExecOptions{ContinueOnError: true})
		require.NoError(t, err)
		for _, r := range res {
			require.NoError(t, r.Error, r.SQL)
		}
		require.True(t, Diff(applied, to).IsEmpty(), sql)
		return sql, applied, to
	}
	base := `
		CREATE SCHEMA s2;
		CREATE ROLE u1;
		CREATE TYPE public.r1 AS RANGE (subtype = int4);
		GRANT USAGE ON TYPE public.r1 TO u1;
		COMMENT ON TYPE public.r1 IS 'kept';
		CREATE FUNCTION public.f(x public.r1) RETURNS int LANGUAGE sql IMMUTABLE AS 'SELECT 1';
		CREATE DOMAIN public.d AS public.r1;
		CREATE TABLE public.t (c public.r1, n int);
		CREATE VIEW public.v AS SELECT c FROM public.t;
	`

	// RENAME TO keeps the range, so its column, function, domain, view, grant, and comment stay.
	sql, _, _ := plan(t, base, `ALTER TYPE public.r1 RENAME TO r2;`)
	require.Equal(t, `ALTER TYPE "public"."r1" RENAME TO "r2"`, sql)

	// SET SCHEMA leaves the multirange behind, as PostgreSQL 17 does.
	sql, applied, _ := plan(t, base, `ALTER TYPE public.r1 SET SCHEMA s2;`)
	require.Equal(t, `ALTER TYPE "public"."r1" SET SCHEMA "s2"`, sql)
	require.NotZero(t, resolveTypeOIDByName(applied, "public", "r1_multirange"))

	// A move and rename of both the range and its multirange.
	sql, _, _ = plan(t, base, `
		ALTER TYPE public.r1 SET SCHEMA s2;
		ALTER TYPE s2.r1 RENAME TO r2;
		ALTER TYPE public.r1_multirange SET SCHEMA s2;
		ALTER TYPE s2.r1_multirange RENAME TO r2_multirange;
	`)
	require.Equal(t, strings.Join([]string{
		`ALTER TYPE "public"."r1" SET SCHEMA "s2"`,
		`ALTER TYPE "s2"."r1" RENAME TO "r2"`,
		`ALTER TYPE "public"."r1_multirange" SET SCHEMA "s2"`,
		`ALTER TYPE "s2"."r1_multirange" RENAME TO "r2_multirange"`,
	}, ";\n"), sql)

	// A move into a schema the migration creates cannot run first, so it stays a rebuild.
	sql, _, _ = plan(t, `CREATE TYPE public.r1 AS RANGE (subtype = int4);`, `
		CREATE SCHEMA s3;
		ALTER TYPE public.r1 SET SCHEMA s3;
		ALTER TYPE public.r1_multirange SET SCHEMA s3;
	`)
	require.NotContains(t, sql, "ALTER TYPE")
	require.Contains(t, sql, `DROP TYPE "public"."r1"`)
}
