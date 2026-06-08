package catalog

import (
	"strings"
	"testing"
)

func TestMigrationSchema(t *testing.T) {
	tests := []struct {
		name    string
		fromSQL string
		toSQL   string
		check   func(t *testing.T, plan *MigrationPlan)
	}{
		{
			name:    "CREATE SCHEMA for added schema",
			fromSQL: "",
			toSQL:   "CREATE SCHEMA analytics;",
			check: func(t *testing.T, plan *MigrationPlan) {
				if len(plan.Ops) != 1 {
					t.Fatalf("expected 1 op, got %d", len(plan.Ops))
				}
				op := plan.Ops[0]
				if op.Type != OpCreateSchema {
					t.Errorf("expected OpCreateSchema, got %s", op.Type)
				}
				if !strings.Contains(op.SQL, "CREATE SCHEMA") {
					t.Errorf("SQL missing CREATE SCHEMA: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "analytics") {
					t.Errorf("SQL missing schema name: %s", op.SQL)
				}
			},
		},
		{
			name:    "DROP SCHEMA CASCADE for dropped schema",
			fromSQL: "CREATE SCHEMA analytics;",
			toSQL:   "",
			check: func(t *testing.T, plan *MigrationPlan) {
				if len(plan.Ops) != 1 {
					t.Fatalf("expected 1 op, got %d", len(plan.Ops))
				}
				op := plan.Ops[0]
				if op.Type != OpDropSchema {
					t.Errorf("expected OpDropSchema, got %s", op.Type)
				}
				if !strings.Contains(op.SQL, "DROP SCHEMA") {
					t.Errorf("SQL missing DROP SCHEMA: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "CASCADE") {
					t.Errorf("SQL missing CASCADE: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "analytics") {
					t.Errorf("SQL missing schema name: %s", op.SQL)
				}
			},
		},
		{
			name:    "ALTER SCHEMA OWNER TO for modified schema owner",
			fromSQL: "CREATE SCHEMA sales AUTHORIZATION alice;",
			toSQL:   "CREATE SCHEMA sales AUTHORIZATION bob;",
			check: func(t *testing.T, plan *MigrationPlan) {
				if len(plan.Ops) != 1 {
					t.Fatalf("expected 1 op, got %d", len(plan.Ops))
				}
				op := plan.Ops[0]
				if op.Type != OpAlterSchema {
					t.Errorf("expected OpAlterSchema, got %s", op.Type)
				}
				if !strings.Contains(op.SQL, "ALTER SCHEMA") {
					t.Errorf("SQL missing ALTER SCHEMA: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "OWNER TO") {
					t.Errorf("SQL missing OWNER TO: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "bob") {
					t.Errorf("SQL missing new owner: %s", op.SQL)
				}
			},
		},
		{
			name:    "schema operations ordered before table operations",
			fromSQL: "",
			toSQL:   "CREATE SCHEMA app; CREATE TABLE app.users (id int);",
			check: func(t *testing.T, plan *MigrationPlan) {
				if len(plan.Ops) == 0 {
					t.Fatal("expected at least 1 op")
				}
				// The first operation must be the schema creation.
				firstOp := plan.Ops[0]
				if firstOp.Type != OpCreateSchema {
					t.Errorf("expected first op to be OpCreateSchema, got %s", firstOp.Type)
				}
				// Schema ops should come before any table ops.
				schemaIdx := -1
				tableIdx := -1
				for i, op := range plan.Ops {
					if op.Type == OpCreateSchema && schemaIdx == -1 {
						schemaIdx = i
					}
					if op.Type == OpCreateTable && tableIdx == -1 {
						tableIdx = i
					}
				}
				if schemaIdx >= 0 && tableIdx >= 0 && schemaIdx >= tableIdx {
					t.Errorf("schema op (idx %d) should come before table op (idx %d)", schemaIdx, tableIdx)
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
