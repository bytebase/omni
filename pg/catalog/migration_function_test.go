package catalog

import (
	"strings"
	"testing"
)

func TestMigrationFunction(t *testing.T) {
	tests := []struct {
		name    string
		fromSQL string
		toSQL   string
		check   func(t *testing.T, plan *MigrationPlan)
	}{
		{
			name:    "CREATE FUNCTION with full signature, return type, language, body",
			fromSQL: "",
			toSQL: `
				CREATE FUNCTION public.add_numbers(a integer, b integer)
				RETURNS integer
				LANGUAGE plpgsql
				AS $$ BEGIN RETURN a + b; END; $$;
			`,
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := filterOps(plan, OpCreateFunction)
				if len(ops) != 1 {
					t.Fatalf("expected 1 CreateFunction op, got %d", len(ops))
				}
				sql := ops[0].SQL
				if !strings.Contains(sql, "CREATE FUNCTION") {
					t.Errorf("missing CREATE FUNCTION: %s", sql)
				}
				if !strings.Contains(sql, `"public"."add_numbers"`) {
					t.Errorf("missing qualified name: %s", sql)
				}
				if !strings.Contains(sql, "integer") {
					t.Errorf("missing arg types: %s", sql)
				}
				if !strings.Contains(sql, "RETURNS integer") {
					t.Errorf("missing RETURNS: %s", sql)
				}
				if !strings.Contains(sql, "LANGUAGE plpgsql") {
					t.Errorf("missing LANGUAGE: %s", sql)
				}
				if !strings.Contains(sql, "RETURN a + b") {
					t.Errorf("missing body: %s", sql)
				}
			},
		},
		{
			name:    "CREATE FUNCTION with volatility, strictness, security, parallel, leakproof",
			fromSQL: "",
			toSQL: `
				CREATE FUNCTION public.safe_fn(x integer)
				RETURNS integer
				LANGUAGE sql
				IMMUTABLE
				STRICT
				SECURITY DEFINER
				LEAKPROOF
				PARALLEL SAFE
				AS $$ SELECT x; $$;
			`,
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := filterOps(plan, OpCreateFunction)
				if len(ops) != 1 {
					t.Fatalf("expected 1 CreateFunction op, got %d", len(ops))
				}
				sql := ops[0].SQL
				if !strings.Contains(sql, "IMMUTABLE") {
					t.Errorf("missing IMMUTABLE: %s", sql)
				}
				if !strings.Contains(sql, "STRICT") {
					t.Errorf("missing STRICT: %s", sql)
				}
				if !strings.Contains(sql, "SECURITY DEFINER") {
					t.Errorf("missing SECURITY DEFINER: %s", sql)
				}
				if !strings.Contains(sql, "LEAKPROOF") {
					t.Errorf("missing LEAKPROOF: %s", sql)
				}
				if !strings.Contains(sql, "PARALLEL SAFE") {
					t.Errorf("missing PARALLEL SAFE: %s", sql)
				}
			},
		},
		{
			name:    "Dollar-quoting for function body",
			fromSQL: "",
			toSQL: `
				CREATE FUNCTION public.greet(name text)
				RETURNS text
				LANGUAGE plpgsql
				AS $fn$
				BEGIN
					RETURN 'Hello, ' || name || '$$';
				END;
				$fn$;
			`,
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := filterOps(plan, OpCreateFunction)
				if len(ops) != 1 {
					t.Fatalf("expected 1 CreateFunction op, got %d", len(ops))
				}
				sql := ops[0].SQL
				// Body contains $$, so output must use $fn$...$fn$.
				if !strings.Contains(sql, "$fn$") {
					t.Errorf("expected $fn$ quoting when body contains $$: %s", sql)
				}
			},
		},
		{
			name: "DROP FUNCTION with argument types in signature",
			fromSQL: `
				CREATE FUNCTION public.add_numbers(a integer, b text)
				RETURNS integer
				LANGUAGE sql
				AS $$ SELECT 1; $$;
			`,
			toSQL: "",
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := filterOps(plan, OpDropFunction)
				if len(ops) != 1 {
					t.Fatalf("expected 1 DropFunction op, got %d", len(ops))
				}
				sql := ops[0].SQL
				if !strings.Contains(sql, "DROP FUNCTION") {
					t.Errorf("missing DROP FUNCTION: %s", sql)
				}
				if !strings.Contains(sql, `"public"."add_numbers"`) {
					t.Errorf("missing qualified name: %s", sql)
				}
				// Must include arg types in the signature.
				if !strings.Contains(sql, "integer") {
					t.Errorf("missing integer arg type in DROP signature: %s", sql)
				}
				if !strings.Contains(sql, "text") {
					t.Errorf("missing text arg type in DROP signature: %s", sql)
				}
			},
		},
		{
			name: "CREATE OR REPLACE FUNCTION for modified function",
			fromSQL: `
				CREATE FUNCTION public.get_val()
				RETURNS integer
				LANGUAGE sql
				AS $$ SELECT 1; $$;
			`,
			toSQL: `
				CREATE FUNCTION public.get_val()
				RETURNS integer
				LANGUAGE sql
				IMMUTABLE
				AS $$ SELECT 42; $$;
			`,
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := filterOps(plan, OpAlterFunction)
				if len(ops) != 1 {
					t.Fatalf("expected 1 AlterFunction op, got %d; all ops: %v", len(ops), plan.Ops)
				}
				sql := ops[0].SQL
				if !strings.Contains(sql, "CREATE OR REPLACE FUNCTION") {
					t.Errorf("missing CREATE OR REPLACE FUNCTION: %s", sql)
				}
				if !strings.Contains(sql, "SELECT 42") {
					t.Errorf("missing new body: %s", sql)
				}
				if !strings.Contains(sql, "IMMUTABLE") {
					t.Errorf("missing IMMUTABLE attribute: %s", sql)
				}
			},
		},
		{
			name:    "CREATE PROCEDURE for added procedure",
			fromSQL: "",
			toSQL: `
				CREATE PROCEDURE public.do_insert(val integer)
				LANGUAGE plpgsql
				AS $$ BEGIN INSERT INTO t VALUES (val); END; $$;
			`,
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := filterOps(plan, OpCreateFunction)
				if len(ops) != 1 {
					t.Fatalf("expected 1 CreateFunction op, got %d", len(ops))
				}
				sql := ops[0].SQL
				if !strings.Contains(sql, "CREATE PROCEDURE") {
					t.Errorf("missing CREATE PROCEDURE: %s", sql)
				}
				if !strings.Contains(sql, `"public"."do_insert"`) {
					t.Errorf("missing qualified name: %s", sql)
				}
				if !strings.Contains(sql, "integer") {
					t.Errorf("missing arg type: %s", sql)
				}
				if !strings.Contains(sql, "LANGUAGE plpgsql") {
					t.Errorf("missing LANGUAGE: %s", sql)
				}
				// Procedure should NOT have RETURNS clause.
				if strings.Contains(sql, "RETURNS") {
					t.Errorf("procedure should not have RETURNS: %s", sql)
				}
			},
		},
		{
			name: "DROP PROCEDURE with argument types",
			fromSQL: `
				CREATE PROCEDURE public.do_insert(val integer)
				LANGUAGE plpgsql
				AS $$ BEGIN NULL; END; $$;
			`,
			toSQL: "",
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := filterOps(plan, OpDropFunction)
				if len(ops) != 1 {
					t.Fatalf("expected 1 DropFunction op, got %d", len(ops))
				}
				sql := ops[0].SQL
				if !strings.Contains(sql, "DROP PROCEDURE") {
					t.Errorf("missing DROP PROCEDURE: %s", sql)
				}
				if !strings.Contains(sql, "integer") {
					t.Errorf("missing arg type in DROP signature: %s", sql)
				}
			},
		},
		{
			name: "Modified procedure body-only as CREATE OR REPLACE PROCEDURE",
			fromSQL: `
				CREATE PROCEDURE public.do_work(x integer)
				LANGUAGE plpgsql
				AS $$ BEGIN NULL; END; $$;
			`,
			toSQL: `
				CREATE PROCEDURE public.do_work(x integer)
				LANGUAGE plpgsql
				AS $$ BEGIN RAISE NOTICE 'done'; END; $$;
			`,
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := filterOps(plan, OpAlterFunction)
				if len(ops) != 1 {
					t.Fatalf("expected 1 AlterFunction op for modified procedure, got %d; all ops: %v", len(ops), plan.Ops)
				}
				sql := ops[0].SQL
				if !strings.Contains(sql, "CREATE OR REPLACE PROCEDURE") {
					t.Errorf("missing CREATE OR REPLACE PROCEDURE: %s", sql)
				}
				if !strings.Contains(sql, "RAISE NOTICE") {
					t.Errorf("missing new body: %s", sql)
				}
				// Should NOT have RETURNS clause.
				if strings.Contains(sql, "RETURNS") {
					t.Errorf("procedure should not have RETURNS: %s", sql)
				}
			},
		},
		{
			name: "Modified procedure with signature change as DROP + CREATE",
			fromSQL: `
				CREATE PROCEDURE public.do_work(x integer)
				LANGUAGE plpgsql
				AS $$ BEGIN NULL; END; $$;
			`,
			toSQL: `
				CREATE PROCEDURE public.do_work(x integer, y integer)
				LANGUAGE plpgsql
				AS $$ BEGIN RAISE NOTICE 'done'; END; $$;
			`,
			check: func(t *testing.T, plan *MigrationPlan) {
				dropOps := filterOps(plan, OpDropFunction)
				createOps := filterOps(plan, OpCreateFunction)
				if len(dropOps) != 1 {
					t.Fatalf("expected 1 DropFunction op for signature-changed procedure, got %d", len(dropOps))
				}
				if len(createOps) != 1 {
					t.Fatalf("expected 1 CreateFunction op for signature-changed procedure, got %d", len(createOps))
				}
				if !strings.Contains(dropOps[0].SQL, "DROP PROCEDURE") {
					t.Errorf("missing DROP PROCEDURE: %s", dropOps[0].SQL)
				}
				if !strings.Contains(createOps[0].SQL, "CREATE PROCEDURE") {
					t.Errorf("missing CREATE PROCEDURE: %s", createOps[0].SQL)
				}
			},
		},
		{
			name:    "Overloaded functions produce distinct DDL",
			fromSQL: "",
			toSQL: `
				CREATE FUNCTION public.calc(a integer)
				RETURNS integer LANGUAGE sql AS $$ SELECT a; $$;

				CREATE FUNCTION public.calc(a integer, b integer)
				RETURNS integer LANGUAGE sql AS $$ SELECT a + b; $$;
			`,
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := filterOps(plan, OpCreateFunction)
				if len(ops) != 2 {
					t.Fatalf("expected 2 CreateFunction ops for overloaded functions, got %d", len(ops))
				}
				// Both should have distinct signatures.
				if ops[0].SQL == ops[1].SQL {
					t.Errorf("overloaded functions produced identical DDL")
				}
				// They should have different ObjectName (Identity).
				if ops[0].ObjectName == ops[1].ObjectName {
					t.Errorf("overloaded functions should have distinct identities, both got: %s", ops[0].ObjectName)
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

