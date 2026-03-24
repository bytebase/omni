package catalog

import (
	"strings"
	"testing"
)

func TestMigrationView(t *testing.T) {
	tests := []struct {
		name    string
		fromSQL string
		toSQL   string
		check   func(t *testing.T, plan *MigrationPlan)
	}{
		{
			name:    "CREATE VIEW AS for added view",
			fromSQL: `CREATE TABLE t (id int, name text, active boolean);`,
			toSQL: `CREATE TABLE t (id int, name text, active boolean);
CREATE VIEW v AS SELECT id, name FROM t WHERE active;`,
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := filterOps(plan, OpCreateView)
				if len(ops) != 1 {
					t.Fatalf("expected 1 CreateView op, got %d", len(ops))
				}
				sql := ops[0].SQL
				if !strings.Contains(sql, "CREATE VIEW") {
					t.Errorf("expected CREATE VIEW, got %s", sql)
				}
				if !strings.Contains(sql, `"v"`) {
					t.Errorf("expected view name quoted, got %s", sql)
				}
				if !strings.Contains(sql, "AS") {
					t.Errorf("expected AS keyword, got %s", sql)
				}
				if !strings.Contains(sql, "SELECT") {
					t.Errorf("expected SELECT in definition, got %s", sql)
				}
			},
		},
		{
			name: "DROP VIEW for removed view",
			fromSQL: `CREATE TABLE t (id int, name text);
CREATE VIEW v AS SELECT id, name FROM t;`,
			toSQL: `CREATE TABLE t (id int, name text);`,
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := filterOps(plan, OpDropView)
				if len(ops) != 1 {
					t.Fatalf("expected 1 DropView op, got %d", len(ops))
				}
				sql := ops[0].SQL
				if !strings.Contains(sql, "DROP VIEW") {
					t.Errorf("expected DROP VIEW, got %s", sql)
				}
				if !strings.Contains(sql, `"v"`) {
					t.Errorf("expected view name quoted, got %s", sql)
				}
			},
		},
		{
			name: "CREATE OR REPLACE VIEW for modified view",
			fromSQL: `CREATE TABLE t (id int, name text, email text);
CREATE VIEW v AS SELECT id, name FROM t;`,
			toSQL: `CREATE TABLE t (id int, name text, email text);
CREATE VIEW v AS SELECT id, name, email FROM t;`,
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := filterOps(plan, OpAlterView)
				found := false
				for _, op := range ops {
					if strings.Contains(op.SQL, "CREATE OR REPLACE VIEW") {
						found = true
						if !strings.Contains(op.SQL, `"v"`) {
							t.Errorf("expected view name quoted, got %s", op.SQL)
						}
						if !strings.Contains(op.SQL, "email") {
							t.Errorf("expected new definition with email, got %s", op.SQL)
						}
						break
					}
				}
				if !found {
					t.Errorf("expected CREATE OR REPLACE VIEW op, found none in ops: %v", ops)
				}
			},
		},
		{
			name:    "CREATE MATERIALIZED VIEW AS for added matview",
			fromSQL: `CREATE TABLE t (id int, name text);`,
			toSQL: `CREATE TABLE t (id int, name text);
CREATE MATERIALIZED VIEW mv AS SELECT id, name FROM t;`,
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := filterOps(plan, OpCreateView)
				found := false
				for _, op := range ops {
					if strings.Contains(op.SQL, "CREATE MATERIALIZED VIEW") {
						found = true
						if !strings.Contains(op.SQL, `"mv"`) {
							t.Errorf("expected matview name quoted, got %s", op.SQL)
						}
						if !strings.Contains(op.SQL, "SELECT") {
							t.Errorf("expected SELECT in definition, got %s", op.SQL)
						}
						break
					}
				}
				if !found {
					t.Errorf("expected CREATE MATERIALIZED VIEW op, found none")
				}
			},
		},
		{
			name: "DROP MATERIALIZED VIEW for removed matview",
			fromSQL: `CREATE TABLE t (id int, name text);
CREATE MATERIALIZED VIEW mv AS SELECT id, name FROM t;`,
			toSQL: `CREATE TABLE t (id int, name text);`,
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := filterOps(plan, OpDropView)
				if len(ops) != 1 {
					t.Fatalf("expected 1 DropView op, got %d", len(ops))
				}
				sql := ops[0].SQL
				if !strings.Contains(sql, "DROP MATERIALIZED VIEW") {
					t.Errorf("expected DROP MATERIALIZED VIEW, got %s", sql)
				}
				if !strings.Contains(sql, `"mv"`) {
					t.Errorf("expected matview name quoted, got %s", sql)
				}
			},
		},
		{
			name: "Modified matview as DROP plus CREATE",
			fromSQL: `CREATE TABLE t (id int, name text, email text);
CREATE MATERIALIZED VIEW mv AS SELECT id, name FROM t;`,
			toSQL: `CREATE TABLE t (id int, name text, email text);
CREATE MATERIALIZED VIEW mv AS SELECT id, name, email FROM t;`,
			check: func(t *testing.T, plan *MigrationPlan) {
				// Should have a DROP MATERIALIZED VIEW and a CREATE MATERIALIZED VIEW.
				dropOps := filterOps(plan, OpDropView)
				createOps := filterOps(plan, OpCreateView)
				foundDrop := false
				for _, op := range dropOps {
					if strings.Contains(op.SQL, "DROP MATERIALIZED VIEW") && strings.Contains(op.SQL, `"mv"`) {
						foundDrop = true
						break
					}
				}
				if !foundDrop {
					t.Errorf("expected DROP MATERIALIZED VIEW for mv, not found")
				}
				foundCreate := false
				for _, op := range createOps {
					if strings.Contains(op.SQL, "CREATE MATERIALIZED VIEW") && strings.Contains(op.SQL, `"mv"`) {
						foundCreate = true
						if !strings.Contains(op.SQL, "email") {
							t.Errorf("expected new definition with email, got %s", op.SQL)
						}
						break
					}
				}
				if !foundCreate {
					t.Errorf("expected CREATE MATERIALIZED VIEW for mv, not found")
				}

				// DROP should come before CREATE in plan order.
				dropIdx := -1
				createIdx := -1
				for i, op := range plan.Ops {
					if strings.Contains(op.SQL, "DROP MATERIALIZED VIEW") && strings.Contains(op.SQL, `"mv"`) {
						dropIdx = i
					}
					if strings.Contains(op.SQL, "CREATE MATERIALIZED VIEW") && strings.Contains(op.SQL, `"mv"`) {
						createIdx = i
					}
				}
				if dropIdx >= 0 && createIdx >= 0 && dropIdx > createIdx {
					t.Errorf("DROP should come before CREATE, drop at %d, create at %d", dropIdx, createIdx)
				}
			},
		},
		{
			name: "View dependency chain generates warning",
			fromSQL: `CREATE TABLE t (id int, name text, email text);
CREATE VIEW v1 AS SELECT id, name FROM t;
CREATE VIEW v2 AS SELECT id FROM v1;`,
			toSQL: `CREATE TABLE t (id int, name text, email text);
CREATE VIEW v1 AS SELECT id, name, email FROM t;
CREATE VIEW v2 AS SELECT id FROM v1;`,
			check: func(t *testing.T, plan *MigrationPlan) {
				// v1 is modified and v2 depends on v1.
				// Should have a warning about dependent views.
				if !plan.HasWarnings() {
					t.Errorf("expected warnings about dependent views")
				}
				warnings := plan.Warnings()
				foundDep := false
				for _, w := range warnings {
					if strings.Contains(w.Warning, "dependent") && strings.Contains(w.Warning, "v2") {
						foundDep = true
						break
					}
				}
				if !foundDep {
					t.Errorf("expected warning mentioning dependent view v2, got warnings: %v", warnings)
				}

				// Should also have a CREATE OR REPLACE VIEW for v1.
				alterOps := filterOps(plan, OpAlterView)
				foundReplace := false
				for _, op := range alterOps {
					if strings.Contains(op.SQL, "CREATE OR REPLACE VIEW") {
						foundReplace = true
						break
					}
				}
				if !foundReplace {
					t.Errorf("expected CREATE OR REPLACE VIEW for v1")
				}
			},
		},
		{
			name: "Matview indexes generated after matview creation",
			fromSQL: `CREATE TABLE t (id int, name text);`,
			toSQL: `CREATE TABLE t (id int, name text);
CREATE MATERIALIZED VIEW mv AS SELECT id, name FROM t;
CREATE INDEX idx_mv_id ON mv (id);`,
			check: func(t *testing.T, plan *MigrationPlan) {
				// Should have CREATE MATERIALIZED VIEW and CREATE INDEX.
				createViewOps := filterOps(plan, OpCreateView)
				createIdxOps := filterOps(plan, OpCreateIndex)

				foundMatView := false
				for _, op := range createViewOps {
					if strings.Contains(op.SQL, "CREATE MATERIALIZED VIEW") && strings.Contains(op.SQL, `"mv"`) {
						foundMatView = true
						break
					}
				}
				if !foundMatView {
					t.Fatalf("expected CREATE MATERIALIZED VIEW for mv")
				}

				if len(createIdxOps) < 1 {
					t.Fatalf("expected at least 1 CreateIndex op, got %d", len(createIdxOps))
				}
				foundIdx := false
				for _, op := range createIdxOps {
					if strings.Contains(op.SQL, "idx_mv_id") {
						foundIdx = true
						break
					}
				}
				if !foundIdx {
					t.Errorf("expected CREATE INDEX idx_mv_id, not found")
				}

				// Index ops should come after matview creation in the plan.
				matviewIdx := -1
				indexIdx := -1
				for i, op := range plan.Ops {
					if strings.Contains(op.SQL, "CREATE MATERIALIZED VIEW") && strings.Contains(op.SQL, `"mv"`) {
						matviewIdx = i
					}
					if strings.Contains(op.SQL, "idx_mv_id") {
						indexIdx = i
					}
				}
				if matviewIdx >= 0 && indexIdx >= 0 && matviewIdx > indexIdx {
					t.Errorf("matview creation (pos %d) should come before index creation (pos %d)", matviewIdx, indexIdx)
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
