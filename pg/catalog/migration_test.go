package catalog

import (
	"strings"
	"testing"
)

func TestMigrationPlan(t *testing.T) {
	t.Run("types compile", func(t *testing.T) {
		// If this test runs, the types compiled successfully.
		var _ MigrationOpType = OpCreateTable
		var _ MigrationOp
		var _ MigrationPlan
	})

	t.Run("empty diff returns empty plan", func(t *testing.T) {
		from, err := LoadSQL("")
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL("")
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := GenerateMigration(from, to, diff)
		if len(plan.Ops) != 0 {
			t.Fatalf("expected 0 ops, got %d", len(plan.Ops))
		}
	})

	t.Run("identical catalogs returns empty plan", func(t *testing.T) {
		sql := `CREATE TABLE t (id int);`
		from, err := LoadSQL(sql)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(sql)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := GenerateMigration(from, to, diff)
		if len(plan.Ops) != 0 {
			t.Fatalf("expected 0 ops for identical catalogs, got %d", len(plan.Ops))
		}
	})

	t.Run("SQL joins with semicolons and newlines", func(t *testing.T) {
		plan := &MigrationPlan{
			Ops: []MigrationOp{
				{Type: OpCreateTable, SQL: "CREATE TABLE a (id int)"},
				{Type: OpCreateTable, SQL: "CREATE TABLE b (id int)"},
				{Type: OpDropTable, SQL: "DROP TABLE c"},
			},
		}
		got := plan.SQL()
		want := "CREATE TABLE a (id int);\nCREATE TABLE b (id int);\nDROP TABLE c"
		if got != want {
			t.Errorf("SQL() =\n%s\nwant:\n%s", got, want)
		}

		// Empty plan returns empty string.
		empty := &MigrationPlan{}
		if s := empty.SQL(); s != "" {
			t.Errorf("empty plan SQL() = %q, want empty", s)
		}
	})

	t.Run("Summary groups and counts", func(t *testing.T) {
		plan := &MigrationPlan{
			Ops: []MigrationOp{
				{Type: OpCreateTable, SQL: "CREATE TABLE a (id int)"},
				{Type: OpCreateTable, SQL: "CREATE TABLE b (id int)"},
				{Type: OpDropTable, SQL: "DROP TABLE c"},
				{Type: OpAlterColumn, SQL: "ALTER TABLE a ALTER COLUMN id TYPE bigint"},
			},
		}
		summary := plan.Summary()
		// Should contain total count and breakdown.
		if !strings.Contains(summary, "4 operation(s)") {
			t.Errorf("Summary missing total count: %s", summary)
		}
		if !strings.Contains(summary, "2 create") {
			t.Errorf("Summary missing create count: %s", summary)
		}
		if !strings.Contains(summary, "1 drop") {
			t.Errorf("Summary missing drop count: %s", summary)
		}
		if !strings.Contains(summary, "1 alter") {
			t.Errorf("Summary missing alter count: %s", summary)
		}
		if !strings.Contains(summary, "CreateTable: 2") {
			t.Errorf("Summary missing CreateTable breakdown: %s", summary)
		}
		if !strings.Contains(summary, "DropTable: 1") {
			t.Errorf("Summary missing DropTable breakdown: %s", summary)
		}

		// Empty plan.
		empty := &MigrationPlan{}
		if s := empty.Summary(); s != "No changes" {
			t.Errorf("empty Summary() = %q, want %q", s, "No changes")
		}
	})

	t.Run("Filter returns subset", func(t *testing.T) {
		plan := &MigrationPlan{
			Ops: []MigrationOp{
				{Type: OpCreateTable, SchemaName: "public", SQL: "CREATE TABLE a (id int)"},
				{Type: OpDropTable, SchemaName: "public", SQL: "DROP TABLE b"},
				{Type: OpCreateIndex, SchemaName: "app", SQL: "CREATE INDEX idx ON c (id)"},
			},
		}
		creates := plan.Filter(func(op MigrationOp) bool {
			return strings.HasPrefix(string(op.Type), "Create")
		})
		if len(creates.Ops) != 2 {
			t.Fatalf("expected 2 create ops, got %d", len(creates.Ops))
		}
		if creates.Ops[0].Type != OpCreateTable {
			t.Errorf("first op type = %s, want CreateTable", creates.Ops[0].Type)
		}
		if creates.Ops[1].Type != OpCreateIndex {
			t.Errorf("second op type = %s, want CreateIndex", creates.Ops[1].Type)
		}

		// Filter that matches nothing.
		none := plan.Filter(func(op MigrationOp) bool { return false })
		if len(none.Ops) != 0 {
			t.Errorf("expected 0 ops, got %d", len(none.Ops))
		}
	})

	t.Run("HasWarnings detects warnings", func(t *testing.T) {
		noWarn := &MigrationPlan{
			Ops: []MigrationOp{
				{Type: OpCreateTable, SQL: "CREATE TABLE a (id int)"},
			},
		}
		if noWarn.HasWarnings() {
			t.Error("HasWarnings() = true for plan with no warnings")
		}

		withWarn := &MigrationPlan{
			Ops: []MigrationOp{
				{Type: OpCreateTable, SQL: "CREATE TABLE a (id int)"},
				{Type: OpDropColumn, SQL: "ALTER TABLE a DROP COLUMN x", Warning: "data loss"},
			},
		}
		if !withWarn.HasWarnings() {
			t.Error("HasWarnings() = false for plan with warnings")
		}

		// Empty plan.
		empty := &MigrationPlan{}
		if empty.HasWarnings() {
			t.Error("HasWarnings() = true for empty plan")
		}
	})

	t.Run("Warnings returns warning ops", func(t *testing.T) {
		plan := &MigrationPlan{
			Ops: []MigrationOp{
				{Type: OpCreateTable, SQL: "CREATE TABLE a (id int)"},
				{Type: OpDropColumn, SQL: "ALTER TABLE a DROP COLUMN x", Warning: "data loss"},
				{Type: OpAlterColumn, SQL: "ALTER TABLE a ALTER COLUMN y TYPE text", Warning: "possible truncation"},
				{Type: OpCreateIndex, SQL: "CREATE INDEX idx ON a (id)"},
			},
		}
		warns := plan.Warnings()
		if len(warns) != 2 {
			t.Fatalf("expected 2 warnings, got %d", len(warns))
		}
		if warns[0].Warning != "data loss" {
			t.Errorf("first warning = %q, want %q", warns[0].Warning, "data loss")
		}
		if warns[1].Warning != "possible truncation" {
			t.Errorf("second warning = %q, want %q", warns[1].Warning, "possible truncation")
		}

		// No warnings.
		clean := &MigrationPlan{
			Ops: []MigrationOp{{Type: OpCreateTable, SQL: "CREATE TABLE a (id int)"}},
		}
		if w := clean.Warnings(); len(w) != 0 {
			t.Errorf("expected 0 warnings, got %d", len(w))
		}
	})
}
