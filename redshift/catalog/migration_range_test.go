package catalog

import (
	"strings"
	"testing"
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
