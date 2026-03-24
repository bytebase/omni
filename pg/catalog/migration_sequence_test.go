package catalog

import (
	"strings"
	"testing"
)

func TestMigrationSequence(t *testing.T) {
	tests := []struct {
		name    string
		fromSQL string
		toSQL   string
		check   func(t *testing.T, plan *MigrationPlan)
	}{
		{
			name:    "CREATE SEQUENCE with all options",
			fromSQL: "",
			toSQL: `CREATE SEQUENCE my_seq
				INCREMENT BY 5
				MINVALUE 10
				MAXVALUE 1000
				START WITH 10
				CACHE 20
				CYCLE;`,
			check: func(t *testing.T, plan *MigrationPlan) {
				seqOps := plan.Filter(func(op MigrationOp) bool {
					return op.Type == OpCreateSequence
				})
				if len(seqOps.Ops) != 1 {
					t.Fatalf("expected 1 CreateSequence op, got %d", len(seqOps.Ops))
				}
				op := seqOps.Ops[0]
				if op.SchemaName != "public" {
					t.Errorf("expected schema public, got %s", op.SchemaName)
				}
				if op.ObjectName != "my_seq" {
					t.Errorf("expected object name my_seq, got %s", op.ObjectName)
				}
				if !strings.Contains(op.SQL, "CREATE SEQUENCE") {
					t.Errorf("SQL missing CREATE SEQUENCE: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "INCREMENT BY 5") {
					t.Errorf("SQL missing INCREMENT BY 5: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "MINVALUE 10") {
					t.Errorf("SQL missing MINVALUE 10: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "MAXVALUE 1000") {
					t.Errorf("SQL missing MAXVALUE 1000: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "START WITH 10") {
					t.Errorf("SQL missing START WITH 10: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "CACHE 20") {
					t.Errorf("SQL missing CACHE 20: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "CYCLE") {
					t.Errorf("SQL missing CYCLE: %s", op.SQL)
				}
			},
		},
		{
			name: "DROP SEQUENCE for removed sequence",
			fromSQL: `CREATE SEQUENCE old_seq
				INCREMENT BY 1
				MINVALUE 1
				MAXVALUE 9223372036854775807
				START WITH 1
				CACHE 1
				NO CYCLE;`,
			toSQL: "",
			check: func(t *testing.T, plan *MigrationPlan) {
				seqOps := plan.Filter(func(op MigrationOp) bool {
					return op.Type == OpDropSequence
				})
				if len(seqOps.Ops) != 1 {
					t.Fatalf("expected 1 DropSequence op, got %d", len(seqOps.Ops))
				}
				op := seqOps.Ops[0]
				if !strings.Contains(op.SQL, "DROP SEQUENCE") {
					t.Errorf("SQL missing DROP SEQUENCE: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "\"old_seq\"") {
					t.Errorf("SQL missing quoted sequence name: %s", op.SQL)
				}
			},
		},
		{
			name: "ALTER SEQUENCE for modified properties",
			fromSQL: `CREATE SEQUENCE counter_seq
				INCREMENT BY 1
				MINVALUE 1
				MAXVALUE 9223372036854775807
				START WITH 1
				CACHE 1
				NO CYCLE;`,
			toSQL: `CREATE SEQUENCE counter_seq
				INCREMENT BY 10
				MINVALUE 100
				MAXVALUE 9223372036854775807
				START WITH 100
				CACHE 5
				CYCLE;`,
			check: func(t *testing.T, plan *MigrationPlan) {
				seqOps := plan.Filter(func(op MigrationOp) bool {
					return op.Type == OpAlterSequence
				})
				if len(seqOps.Ops) != 1 {
					t.Fatalf("expected 1 AlterSequence op, got %d", len(seqOps.Ops))
				}
				op := seqOps.Ops[0]
				if !strings.Contains(op.SQL, "ALTER SEQUENCE") {
					t.Errorf("SQL missing ALTER SEQUENCE: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "INCREMENT BY 10") {
					t.Errorf("SQL missing INCREMENT BY 10: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "MINVALUE 100") {
					t.Errorf("SQL missing MINVALUE 100: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "START WITH 100") {
					t.Errorf("SQL missing START WITH 100: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "CACHE 5") {
					t.Errorf("SQL missing CACHE 5: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "CYCLE") {
					t.Errorf("SQL missing CYCLE: %s", op.SQL)
				}
				// Should NOT contain MAXVALUE since it didn't change
				if strings.Contains(op.SQL, "MAXVALUE") {
					t.Errorf("SQL should not contain MAXVALUE (unchanged): %s", op.SQL)
				}
			},
		},
		{
			name: "SERIAL-owned sequences not generated",
			fromSQL: "",
			toSQL: `CREATE TABLE orders (
				id SERIAL PRIMARY KEY,
				name text
			);`,
			check: func(t *testing.T, plan *MigrationPlan) {
				seqOps := plan.Filter(func(op MigrationOp) bool {
					return op.Type == OpCreateSequence || op.Type == OpDropSequence || op.Type == OpAlterSequence
				})
				if len(seqOps.Ops) != 0 {
					t.Errorf("expected 0 sequence ops for SERIAL column, got %d", len(seqOps.Ops))
					for _, op := range seqOps.Ops {
						t.Logf("  unexpected op: %s %s", op.Type, op.SQL)
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
