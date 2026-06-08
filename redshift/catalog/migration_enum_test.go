package catalog

import (
	"strings"
	"testing"
)

func TestMigrationEnum(t *testing.T) {
	tests := []struct {
		name    string
		fromSQL string
		toSQL   string
		check   func(t *testing.T, plan *MigrationPlan)
	}{
		{
			name:    "CREATE TYPE AS ENUM for added enum",
			fromSQL: "",
			toSQL:   "CREATE TYPE public.status AS ENUM ('active', 'inactive');",
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := plan.Filter(func(op MigrationOp) bool { return op.Type == OpCreateType }).Ops
				if len(ops) != 1 {
					t.Fatalf("expected 1 CreateType op, got %d; all ops: %v", len(ops), plan.Ops)
				}
				op := ops[0]
				if !strings.Contains(op.SQL, "CREATE TYPE") {
					t.Errorf("SQL missing CREATE TYPE: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "AS ENUM") {
					t.Errorf("SQL missing AS ENUM: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "'active'") {
					t.Errorf("SQL missing 'active': %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "'inactive'") {
					t.Errorf("SQL missing 'inactive': %s", op.SQL)
				}
				if !op.Transactional {
					t.Errorf("CREATE TYPE should be transactional")
				}
			},
		},
		{
			name:    "DROP TYPE for removed enum",
			fromSQL: "CREATE TYPE public.status AS ENUM ('active', 'inactive');",
			toSQL:   "",
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := plan.Filter(func(op MigrationOp) bool { return op.Type == OpDropType }).Ops
				if len(ops) != 1 {
					t.Fatalf("expected 1 DropType op, got %d; all ops: %v", len(ops), plan.Ops)
				}
				op := ops[0]
				if !strings.Contains(op.SQL, "DROP TYPE") {
					t.Errorf("SQL missing DROP TYPE: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "\"status\"") {
					t.Errorf("SQL missing type name: %s", op.SQL)
				}
				if !op.Transactional {
					t.Errorf("DROP TYPE should be transactional")
				}
			},
		},
		{
			name:    "ALTER TYPE ADD VALUE for appended value",
			fromSQL: "CREATE TYPE public.status AS ENUM ('active', 'inactive');",
			toSQL:   "CREATE TYPE public.status AS ENUM ('active', 'inactive', 'pending');",
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := plan.Filter(func(op MigrationOp) bool { return op.Type == OpAlterType }).Ops
				if len(ops) != 1 {
					t.Fatalf("expected 1 AlterType op, got %d; all ops: %v", len(ops), plan.Ops)
				}
				op := ops[0]
				if !strings.Contains(op.SQL, "ALTER TYPE") {
					t.Errorf("SQL missing ALTER TYPE: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "ADD VALUE") {
					t.Errorf("SQL missing ADD VALUE: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "'pending'") {
					t.Errorf("SQL missing 'pending': %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "AFTER 'inactive'") {
					t.Errorf("SQL missing AFTER 'inactive': %s", op.SQL)
				}
				if op.Transactional {
					t.Errorf("ADD VALUE should NOT be transactional")
				}
			},
		},
		{
			name:    "ALTER TYPE ADD VALUE AFTER existing for positioned value",
			fromSQL: "CREATE TYPE public.status AS ENUM ('active', 'inactive', 'deleted');",
			toSQL:   "CREATE TYPE public.status AS ENUM ('active', 'inactive', 'pending', 'deleted');",
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := plan.Filter(func(op MigrationOp) bool { return op.Type == OpAlterType }).Ops
				if len(ops) != 1 {
					t.Fatalf("expected 1 AlterType op, got %d; all ops: %v", len(ops), plan.Ops)
				}
				op := ops[0]
				if !strings.Contains(op.SQL, "ADD VALUE 'pending'") {
					t.Errorf("SQL missing ADD VALUE 'pending': %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "AFTER 'inactive'") {
					t.Errorf("SQL missing AFTER 'inactive': %s", op.SQL)
				}
				if op.Transactional {
					t.Errorf("ADD VALUE should NOT be transactional")
				}
			},
		},
		{
			name:    "ALTER TYPE ADD VALUE BEFORE existing for positioned value",
			fromSQL: "CREATE TYPE public.status AS ENUM ('active', 'inactive');",
			toSQL:   "CREATE TYPE public.status AS ENUM ('pending', 'active', 'inactive');",
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := plan.Filter(func(op MigrationOp) bool { return op.Type == OpAlterType }).Ops
				if len(ops) != 1 {
					t.Fatalf("expected 1 AlterType op, got %d; all ops: %v", len(ops), plan.Ops)
				}
				op := ops[0]
				if !strings.Contains(op.SQL, "ADD VALUE 'pending'") {
					t.Errorf("SQL missing ADD VALUE 'pending': %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "BEFORE 'active'") {
					t.Errorf("SQL missing BEFORE 'active': %s", op.SQL)
				}
				if op.Transactional {
					t.Errorf("ADD VALUE should NOT be transactional")
				}
			},
		},
		{
			name:    "enum value removal generates Warning",
			fromSQL: "CREATE TYPE public.status AS ENUM ('active', 'inactive', 'deleted');",
			toSQL:   "CREATE TYPE public.status AS ENUM ('active', 'inactive');",
			check: func(t *testing.T, plan *MigrationPlan) {
				if !plan.HasWarnings() {
					t.Fatalf("expected warnings for removed enum value")
				}
				warnings := plan.Warnings()
				if len(warnings) != 1 {
					t.Fatalf("expected 1 warning, got %d", len(warnings))
				}
				w := warnings[0]
				if !strings.Contains(w.Warning, "remov") || !strings.Contains(w.Warning, "'deleted'") {
					t.Errorf("warning should mention inability to remove value: %s", w.Warning)
				}
			},
		},
		{
			name:    "ALTER TYPE ADD VALUE is not transactional",
			fromSQL: "CREATE TYPE public.color AS ENUM ('red');",
			toSQL:   "CREATE TYPE public.color AS ENUM ('red', 'blue', 'green');",
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := plan.Filter(func(op MigrationOp) bool { return op.Type == OpAlterType }).Ops
				if len(ops) != 2 {
					t.Fatalf("expected 2 AlterType ops, got %d; all ops: %v", len(ops), plan.Ops)
				}
				for _, op := range ops {
					if op.Transactional {
						t.Errorf("ADD VALUE op should have Transactional=false: %s", op.SQL)
					}
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
