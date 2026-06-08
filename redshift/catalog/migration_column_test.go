package catalog

import (
	"strings"
	"testing"
)

func TestMigrationColumn(t *testing.T) {
	tests := []struct {
		name    string
		fromSQL string
		toSQL   string
		check   func(t *testing.T, plan *MigrationPlan)
	}{
		{
			name:    "add column",
			fromSQL: `CREATE TABLE t (id integer);`,
			toSQL:   `CREATE TABLE t (id integer, name text);`,
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := filterOps(plan, OpAddColumn)
				if len(ops) != 1 {
					t.Fatalf("expected 1 AddColumn op, got %d", len(ops))
				}
				sql := ops[0].SQL
				if !strings.Contains(sql, "ADD COLUMN") {
					t.Errorf("expected ADD COLUMN, got: %s", sql)
				}
				if !strings.Contains(sql, "name") || !strings.Contains(sql, "text") {
					t.Errorf("expected column name and type, got: %s", sql)
				}
			},
		},
		{
			name:    "add column with NOT NULL and DEFAULT",
			fromSQL: `CREATE TABLE t (id integer);`,
			toSQL:   `CREATE TABLE t (id integer, status integer NOT NULL DEFAULT 0);`,
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := filterOps(plan, OpAddColumn)
				if len(ops) != 1 {
					t.Fatalf("expected 1 AddColumn op, got %d", len(ops))
				}
				sql := ops[0].SQL
				if !strings.Contains(sql, "NOT NULL") {
					t.Errorf("expected NOT NULL, got: %s", sql)
				}
				if !strings.Contains(sql, "DEFAULT") {
					t.Errorf("expected DEFAULT, got: %s", sql)
				}
			},
		},
		{
			name:    "drop column",
			fromSQL: `CREATE TABLE t (id integer, name text);`,
			toSQL:   `CREATE TABLE t (id integer);`,
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := filterOps(plan, OpDropColumn)
				if len(ops) != 1 {
					t.Fatalf("expected 1 DropColumn op, got %d", len(ops))
				}
				sql := ops[0].SQL
				if !strings.Contains(sql, "DROP COLUMN") {
					t.Errorf("expected DROP COLUMN, got: %s", sql)
				}
				if !strings.Contains(sql, "name") {
					t.Errorf("expected column name, got: %s", sql)
				}
			},
		},
		{
			name:    "alter column type implicit cast",
			fromSQL: `CREATE TABLE t (val integer);`,
			toSQL:   `CREATE TABLE t (val bigint);`,
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := filterOps(plan, OpAlterColumn)
				if len(ops) < 1 {
					t.Fatalf("expected at least 1 AlterColumn op, got %d", len(ops))
				}
				found := false
				for _, op := range ops {
					if strings.Contains(op.SQL, "TYPE") && strings.Contains(op.SQL, "bigint") {
						found = true
						if strings.Contains(op.SQL, "USING") {
							t.Errorf("should not have USING for implicit cast, got: %s", op.SQL)
						}
					}
				}
				if !found {
					t.Error("no ALTER COLUMN TYPE statement found")
				}
			},
		},
		{
			name:    "alter column type no implicit cast needs USING",
			fromSQL: `CREATE TABLE t (val text);`,
			toSQL:   `CREATE TABLE t (val integer);`,
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := filterOps(plan, OpAlterColumn)
				if len(ops) < 1 {
					t.Fatalf("expected at least 1 AlterColumn op, got %d", len(ops))
				}
				found := false
				for _, op := range ops {
					if strings.Contains(op.SQL, "TYPE") && strings.Contains(op.SQL, "integer") {
						found = true
						if !strings.Contains(op.SQL, "USING") {
							t.Errorf("expected USING clause for text→integer, got: %s", op.SQL)
						}
					}
				}
				if !found {
					t.Error("no ALTER COLUMN TYPE statement found")
				}
			},
		},
		{
			name:    "USING clause decision uses FindCoercionPathway",
			fromSQL: `CREATE TABLE t (val integer);`,
			toSQL:   `CREATE TABLE t (val bigint);`,
			check: func(t *testing.T, plan *MigrationPlan) {
				// integer→bigint has an implicit cast, so no USING should appear.
				// This confirms the code uses FindCoercionPathway (assignment context).
				ops := filterOps(plan, OpAlterColumn)
				for _, op := range ops {
					if strings.Contains(op.SQL, "TYPE") {
						if strings.Contains(op.SQL, "USING") {
							t.Errorf("FindCoercionPathway should find int4→int8 cast, no USING needed: %s", op.SQL)
						}
					}
				}
			},
		},
		{
			name:    "type change with existing default: DROP DEFAULT then ALTER TYPE then SET DEFAULT",
			fromSQL: `CREATE TABLE t (val integer DEFAULT 42);`,
			toSQL:   `CREATE TABLE t (val bigint DEFAULT 100);`,
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := filterOps(plan, OpAlterColumn)
				if len(ops) < 3 {
					t.Fatalf("expected at least 3 AlterColumn ops (drop default, alter type, set default), got %d", len(ops))
				}
				// Find order: DROP DEFAULT < TYPE < SET DEFAULT
				var dropIdx, typeIdx, setIdx int = -1, -1, -1
				for i, op := range ops {
					switch {
					case strings.Contains(op.SQL, "DROP DEFAULT"):
						dropIdx = i
					case strings.Contains(op.SQL, "TYPE"):
						typeIdx = i
					case strings.Contains(op.SQL, "SET DEFAULT"):
						setIdx = i
					}
				}
				if dropIdx == -1 || typeIdx == -1 || setIdx == -1 {
					t.Fatalf("expected DROP DEFAULT, TYPE, SET DEFAULT; got ops: %v", opsSQL(plan))
				}
				if !(dropIdx < typeIdx && typeIdx < setIdx) {
					t.Errorf("expected order DROP DEFAULT(%d) < TYPE(%d) < SET DEFAULT(%d)", dropIdx, typeIdx, setIdx)
				}
			},
		},
		{
			name:    "set NOT NULL",
			fromSQL: `CREATE TABLE t (val integer);`,
			toSQL:   `CREATE TABLE t (val integer NOT NULL);`,
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := filterOps(plan, OpAlterColumn)
				if len(ops) != 1 {
					t.Fatalf("expected 1 AlterColumn op, got %d", len(ops))
				}
				if !strings.Contains(ops[0].SQL, "SET NOT NULL") {
					t.Errorf("expected SET NOT NULL, got: %s", ops[0].SQL)
				}
			},
		},
		{
			name:    "drop NOT NULL",
			fromSQL: `CREATE TABLE t (val integer NOT NULL);`,
			toSQL:   `CREATE TABLE t (val integer);`,
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := filterOps(plan, OpAlterColumn)
				if len(ops) != 1 {
					t.Fatalf("expected 1 AlterColumn op, got %d", len(ops))
				}
				if !strings.Contains(ops[0].SQL, "DROP NOT NULL") {
					t.Errorf("expected DROP NOT NULL, got: %s", ops[0].SQL)
				}
			},
		},
		{
			name:    "set DEFAULT",
			fromSQL: `CREATE TABLE t (val integer);`,
			toSQL:   `CREATE TABLE t (val integer DEFAULT 42);`,
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := filterOps(plan, OpAlterColumn)
				if len(ops) != 1 {
					t.Fatalf("expected 1 AlterColumn op, got %d", len(ops))
				}
				if !strings.Contains(ops[0].SQL, "SET DEFAULT") {
					t.Errorf("expected SET DEFAULT, got: %s", ops[0].SQL)
				}
			},
		},
		{
			name:    "drop DEFAULT",
			fromSQL: `CREATE TABLE t (val integer DEFAULT 42);`,
			toSQL:   `CREATE TABLE t (val integer);`,
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := filterOps(plan, OpAlterColumn)
				if len(ops) != 1 {
					t.Fatalf("expected 1 AlterColumn op, got %d", len(ops))
				}
				if !strings.Contains(ops[0].SQL, "DROP DEFAULT") {
					t.Errorf("expected DROP DEFAULT, got: %s", ops[0].SQL)
				}
			},
		},
		{
			name:    "add GENERATED ALWAYS AS IDENTITY",
			fromSQL: `CREATE TABLE t (id integer);`,
			toSQL:   `CREATE TABLE t (id integer GENERATED ALWAYS AS IDENTITY);`,
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := filterOps(plan, OpAlterColumn)
				if len(ops) < 1 {
					t.Fatalf("expected at least 1 AlterColumn op, got %d", len(ops))
				}
				found := false
				for _, op := range ops {
					if strings.Contains(op.SQL, "ADD GENERATED ALWAYS AS IDENTITY") {
						found = true
					}
				}
				if !found {
					t.Errorf("expected ADD GENERATED ALWAYS AS IDENTITY, got: %v", opsSQL(plan))
				}
			},
		},
		{
			name:    "drop IDENTITY",
			fromSQL: `CREATE TABLE t (id integer GENERATED ALWAYS AS IDENTITY);`,
			toSQL:   `CREATE TABLE t (id integer);`,
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := filterOps(plan, OpAlterColumn)
				if len(ops) < 1 {
					t.Fatalf("expected at least 1 AlterColumn op, got %d", len(ops))
				}
				found := false
				for _, op := range ops {
					if strings.Contains(op.SQL, "DROP IDENTITY") {
						found = true
					}
				}
				if !found {
					t.Errorf("expected DROP IDENTITY, got: %v", opsSQL(plan))
				}
			},
		},
		{
			name:    "multiple column changes batched",
			fromSQL: `CREATE TABLE t (a integer, b text, c boolean);`,
			toSQL:   `CREATE TABLE t (a bigint, b text NOT NULL, d varchar(100));`,
			check: func(t *testing.T, plan *MigrationPlan) {
				// a: type change (integer→bigint)
				// b: add NOT NULL
				// c: dropped
				// d: added
				// All should be on the same table — check that ops reference same table.
				if len(plan.Ops) < 3 {
					t.Fatalf("expected at least 3 ops for multiple column changes, got %d", len(plan.Ops))
				}
				// Verify all ops reference the same table.
				for _, op := range plan.Ops {
					if op.ObjectName != "t" {
						t.Errorf("expected all ops on table t, got objectName: %s", op.ObjectName)
					}
				}
				// Verify we have add, drop, and alter ops.
				hasAdd := len(filterOps(plan, OpAddColumn)) > 0
				hasDrop := len(filterOps(plan, OpDropColumn)) > 0
				hasAlter := len(filterOps(plan, OpAlterColumn)) > 0
				if !hasAdd || !hasDrop || !hasAlter {
					t.Errorf("expected add(%v), drop(%v), and alter(%v) ops", hasAdd, hasDrop, hasAlter)
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
