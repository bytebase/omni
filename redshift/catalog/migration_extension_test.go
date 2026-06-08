package catalog

import (
	"strings"
	"testing"
)

func TestMigrationExtension(t *testing.T) {
	// Register lightweight test extensions with empty DDL so tests don't
	// depend on heavy bundled scripts.
	RegisterExtensionSQL("test_ext", "")
	RegisterExtensionSQL("test_ext2", "")

	tests := []struct {
		name    string
		fromSQL string
		toSQL   string
		check   func(t *testing.T, plan *MigrationPlan)
	}{
		{
			name:    "CREATE EXTENSION for added extension",
			fromSQL: "",
			toSQL:   "CREATE EXTENSION test_ext;",
			check: func(t *testing.T, plan *MigrationPlan) {
				extOps := plan.Filter(func(op MigrationOp) bool {
					return op.Type == OpCreateExtension
				})
				if len(extOps.Ops) != 1 {
					t.Fatalf("expected 1 CreateExtension op, got %d", len(extOps.Ops))
				}
				op := extOps.Ops[0]
				if op.ObjectName != "test_ext" {
					t.Errorf("expected ObjectName test_ext, got %s", op.ObjectName)
				}
				if !strings.Contains(op.SQL, "CREATE EXTENSION") {
					t.Errorf("SQL missing CREATE EXTENSION: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, `"test_ext"`) {
					t.Errorf("SQL missing quoted extension name: %s", op.SQL)
				}
			},
		},
		{
			name:    "DROP EXTENSION for removed extension",
			fromSQL: "CREATE EXTENSION test_ext;",
			toSQL:   "",
			check: func(t *testing.T, plan *MigrationPlan) {
				extOps := plan.Filter(func(op MigrationOp) bool {
					return op.Type == OpDropExtension
				})
				if len(extOps.Ops) != 1 {
					t.Fatalf("expected 1 DropExtension op, got %d", len(extOps.Ops))
				}
				op := extOps.Ops[0]
				if op.ObjectName != "test_ext" {
					t.Errorf("expected ObjectName test_ext, got %s", op.ObjectName)
				}
				if !strings.Contains(op.SQL, "DROP EXTENSION") {
					t.Errorf("SQL missing DROP EXTENSION: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, `"test_ext"`) {
					t.Errorf("SQL missing quoted extension name: %s", op.SQL)
				}
			},
		},
		{
			name:    "Extension schema change generates appropriate DDL",
			fromSQL: "CREATE EXTENSION test_ext;",
			toSQL:   "CREATE SCHEMA other; CREATE EXTENSION test_ext SCHEMA other;",
			check: func(t *testing.T, plan *MigrationPlan) {
				alterOps := plan.Filter(func(op MigrationOp) bool {
					return op.Type == OpAlterExtension
				})
				if len(alterOps.Ops) != 1 {
					t.Fatalf("expected 1 AlterExtension op, got %d; all ops: %v", len(alterOps.Ops), opsDebug(plan.Ops))
				}
				op := alterOps.Ops[0]
				if !strings.Contains(op.SQL, "ALTER EXTENSION") {
					t.Errorf("SQL missing ALTER EXTENSION: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, "SET SCHEMA") {
					t.Errorf("SQL missing SET SCHEMA: %s", op.SQL)
				}
				if !strings.Contains(op.SQL, `"other"`) {
					t.Errorf("SQL missing target schema name: %s", op.SQL)
				}
			},
		},
		{
			name:    "Extension operations ordered before types and tables",
			fromSQL: "",
			toSQL:   "CREATE EXTENSION test_ext; CREATE TABLE t1 (id int);",
			check: func(t *testing.T, plan *MigrationPlan) {
				extIdx := -1
				tableIdx := -1
				for i, op := range plan.Ops {
					if op.Type == OpCreateExtension && extIdx == -1 {
						extIdx = i
					}
					if op.Type == OpCreateTable && tableIdx == -1 {
						tableIdx = i
					}
				}
				if extIdx == -1 {
					t.Fatal("expected CreateExtension op")
				}
				if tableIdx == -1 {
					t.Fatal("expected CreateTable op")
				}
				if extIdx >= tableIdx {
					t.Errorf("extension op (idx %d) should come before table op (idx %d)", extIdx, tableIdx)
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

// opsDebug returns a brief description of ops for debugging.
func opsDebug(ops []MigrationOp) []string {
	var result []string
	for _, op := range ops {
		result = append(result, string(op.Type)+": "+op.SQL)
	}
	return result
}
