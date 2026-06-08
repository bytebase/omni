package catalog

import (
	"strings"
	"testing"
)

func TestMigrationDomain(t *testing.T) {
	tests := []struct {
		name    string
		fromSQL string
		toSQL   string
		check   func(t *testing.T, plan *MigrationPlan)
	}{
		{
			name:    "create domain with base type default not null and constraints",
			fromSQL: "",
			toSQL:   `CREATE DOMAIN posint AS integer DEFAULT 0 NOT NULL CONSTRAINT pos CHECK (VALUE > 0);`,
			check: func(t *testing.T, plan *MigrationPlan) {
				if len(plan.Ops) != 1 {
					t.Fatalf("expected 1 op, got %d", len(plan.Ops))
				}
				op := plan.Ops[0]
				if op.Type != OpCreateType {
					t.Errorf("expected OpCreateType, got %s", op.Type)
				}
				sql := op.SQL
				if !strings.Contains(sql, "CREATE DOMAIN") {
					t.Errorf("expected CREATE DOMAIN, got %s", sql)
				}
				if !strings.Contains(sql, "integer") {
					t.Errorf("expected base type integer, got %s", sql)
				}
				if !strings.Contains(sql, "DEFAULT 0") {
					t.Errorf("expected DEFAULT 0, got %s", sql)
				}
				if !strings.Contains(sql, "NOT NULL") {
					t.Errorf("expected NOT NULL, got %s", sql)
				}
				if !strings.Contains(sql, "CHECK") {
					t.Errorf("expected CHECK constraint, got %s", sql)
				}
			},
		},
		{
			name:    "drop domain",
			fromSQL: `CREATE DOMAIN posint AS integer CHECK (VALUE > 0);`,
			toSQL:   "",
			check: func(t *testing.T, plan *MigrationPlan) {
				if len(plan.Ops) != 1 {
					t.Fatalf("expected 1 op, got %d", len(plan.Ops))
				}
				op := plan.Ops[0]
				if op.Type != OpDropType {
					t.Errorf("expected OpDropType, got %s", op.Type)
				}
				if !strings.Contains(op.SQL, "DROP DOMAIN") {
					t.Errorf("expected DROP DOMAIN, got %s", op.SQL)
				}
			},
		},
		{
			name:    "alter domain set default",
			fromSQL: `CREATE DOMAIN mydom AS integer;`,
			toSQL:   `CREATE DOMAIN mydom AS integer DEFAULT 42;`,
			check: func(t *testing.T, plan *MigrationPlan) {
				if len(plan.Ops) != 1 {
					t.Fatalf("expected 1 op, got %d", len(plan.Ops))
				}
				op := plan.Ops[0]
				if op.Type != OpAlterType {
					t.Errorf("expected OpAlterType, got %s", op.Type)
				}
				if !strings.Contains(op.SQL, "SET DEFAULT 42") {
					t.Errorf("expected SET DEFAULT 42, got %s", op.SQL)
				}
			},
		},
		{
			name:    "alter domain drop default",
			fromSQL: `CREATE DOMAIN mydom AS integer DEFAULT 42;`,
			toSQL:   `CREATE DOMAIN mydom AS integer;`,
			check: func(t *testing.T, plan *MigrationPlan) {
				if len(plan.Ops) != 1 {
					t.Fatalf("expected 1 op, got %d", len(plan.Ops))
				}
				op := plan.Ops[0]
				if op.Type != OpAlterType {
					t.Errorf("expected OpAlterType, got %s", op.Type)
				}
				if !strings.Contains(op.SQL, "DROP DEFAULT") {
					t.Errorf("expected DROP DEFAULT, got %s", op.SQL)
				}
			},
		},
		{
			name:    "alter domain set not null",
			fromSQL: `CREATE DOMAIN mydom AS integer;`,
			toSQL:   `CREATE DOMAIN mydom AS integer NOT NULL;`,
			check: func(t *testing.T, plan *MigrationPlan) {
				if len(plan.Ops) != 1 {
					t.Fatalf("expected 1 op, got %d", len(plan.Ops))
				}
				op := plan.Ops[0]
				if op.Type != OpAlterType {
					t.Errorf("expected OpAlterType, got %s", op.Type)
				}
				if !strings.Contains(op.SQL, "SET NOT NULL") {
					t.Errorf("expected SET NOT NULL, got %s", op.SQL)
				}
			},
		},
		{
			name:    "alter domain drop not null",
			fromSQL: `CREATE DOMAIN mydom AS integer NOT NULL;`,
			toSQL:   `CREATE DOMAIN mydom AS integer;`,
			check: func(t *testing.T, plan *MigrationPlan) {
				if len(plan.Ops) != 1 {
					t.Fatalf("expected 1 op, got %d", len(plan.Ops))
				}
				op := plan.Ops[0]
				if op.Type != OpAlterType {
					t.Errorf("expected OpAlterType, got %s", op.Type)
				}
				if !strings.Contains(op.SQL, "DROP NOT NULL") {
					t.Errorf("expected DROP NOT NULL, got %s", op.SQL)
				}
			},
		},
		{
			name:    "alter domain add constraint",
			fromSQL: `CREATE DOMAIN mydom AS integer;`,
			toSQL:   `CREATE DOMAIN mydom AS integer CONSTRAINT pos CHECK (VALUE > 0);`,
			check: func(t *testing.T, plan *MigrationPlan) {
				if len(plan.Ops) != 1 {
					t.Fatalf("expected 1 op, got %d", len(plan.Ops))
				}
				op := plan.Ops[0]
				if op.Type != OpAlterType {
					t.Errorf("expected OpAlterType, got %s", op.Type)
				}
				if !strings.Contains(op.SQL, "ADD CONSTRAINT") {
					t.Errorf("expected ADD CONSTRAINT, got %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "CHECK") {
					t.Errorf("expected CHECK, got %s", op.SQL)
				}
			},
		},
		{
			name:    "alter domain drop constraint",
			fromSQL: `CREATE DOMAIN mydom AS integer CONSTRAINT pos CHECK (VALUE > 0);`,
			toSQL:   `CREATE DOMAIN mydom AS integer;`,
			check: func(t *testing.T, plan *MigrationPlan) {
				if len(plan.Ops) != 1 {
					t.Fatalf("expected 1 op, got %d", len(plan.Ops))
				}
				op := plan.Ops[0]
				if op.Type != OpAlterType {
					t.Errorf("expected OpAlterType, got %s", op.Type)
				}
				if !strings.Contains(op.SQL, "DROP CONSTRAINT") {
					t.Errorf("expected DROP CONSTRAINT, got %s", op.SQL)
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
