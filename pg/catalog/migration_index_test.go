package catalog

import (
	"strings"
	"testing"
)

func TestMigrationIndex(t *testing.T) {
	tests := []struct {
		name    string
		fromSQL string
		toSQL   string
		check   func(t *testing.T, plan *MigrationPlan)
	}{
		{
			name:    "CREATE INDEX for added standalone index",
			fromSQL: `CREATE TABLE t (id int, name text);`,
			toSQL: `CREATE TABLE t (id int, name text);
CREATE INDEX idx_t_name ON t (name);`,
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := filterOps(plan, OpCreateIndex)
				if len(ops) != 1 {
					t.Fatalf("expected 1 CreateIndex op, got %d", len(ops))
				}
				sql := ops[0].SQL
				if !strings.Contains(sql, "CREATE INDEX") {
					t.Errorf("expected CREATE INDEX, got %s", sql)
				}
				if !strings.Contains(sql, "idx_t_name") {
					t.Errorf("expected index name idx_t_name, got %s", sql)
				}
				if !strings.Contains(sql, "(name") {
					t.Errorf("expected column name, got %s", sql)
				}
			},
		},
		{
			name:    "CREATE UNIQUE INDEX for unique indexes",
			fromSQL: `CREATE TABLE t (id int, email text);`,
			toSQL: `CREATE TABLE t (id int, email text);
CREATE UNIQUE INDEX idx_t_email ON t (email);`,
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := filterOps(plan, OpCreateIndex)
				if len(ops) != 1 {
					t.Fatalf("expected 1 CreateIndex op, got %d", len(ops))
				}
				if !strings.Contains(ops[0].SQL, "CREATE UNIQUE INDEX") {
					t.Errorf("expected CREATE UNIQUE INDEX, got %s", ops[0].SQL)
				}
			},
		},
		{
			name:    "CREATE INDEX USING hash/gist/gin for non-btree",
			fromSQL: `CREATE TABLE t (id int, data jsonb);`,
			toSQL: `CREATE TABLE t (id int, data jsonb);
CREATE INDEX idx_t_data ON t USING gin (data);`,
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := filterOps(plan, OpCreateIndex)
				if len(ops) != 1 {
					t.Fatalf("expected 1 CreateIndex op, got %d", len(ops))
				}
				sql := ops[0].SQL
				if !strings.Contains(sql, "USING gin") {
					t.Errorf("expected USING gin, got %s", sql)
				}
				// btree should not appear
				if strings.Contains(sql, "USING btree") {
					t.Errorf("btree is default and should be omitted, got %s", sql)
				}
			},
		},
		{
			name:    "CREATE INDEX with WHERE clause (partial index)",
			fromSQL: `CREATE TABLE t (id int, active boolean);`,
			toSQL: `CREATE TABLE t (id int, active boolean);
CREATE INDEX idx_t_active ON t (id) WHERE active;`,
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := filterOps(plan, OpCreateIndex)
				if len(ops) != 1 {
					t.Fatalf("expected 1 CreateIndex op, got %d", len(ops))
				}
				sql := ops[0].SQL
				if !strings.Contains(sql, "WHERE") {
					t.Errorf("expected WHERE clause, got %s", sql)
				}
			},
		},
		{
			name:    "CREATE INDEX with INCLUDE columns",
			fromSQL: `CREATE TABLE t (id int, name text, email text);`,
			toSQL: `CREATE TABLE t (id int, name text, email text);
CREATE INDEX idx_t_name ON t (name) INCLUDE (email);`,
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := filterOps(plan, OpCreateIndex)
				if len(ops) != 1 {
					t.Fatalf("expected 1 CreateIndex op, got %d", len(ops))
				}
				sql := ops[0].SQL
				if !strings.Contains(sql, "INCLUDE") {
					t.Errorf("expected INCLUDE clause, got %s", sql)
				}
				// Check that email appears after INCLUDE
				includeIdx := strings.Index(sql, "INCLUDE")
				if includeIdx < 0 || !strings.Contains(sql[includeIdx:], "email") {
					t.Errorf("expected email in INCLUDE clause, got %s", sql)
				}
			},
		},
		{
			name:    "CREATE INDEX with expression columns",
			fromSQL: `CREATE TABLE t (id int, name text);`,
			toSQL: `CREATE TABLE t (id int, name text);
CREATE INDEX idx_t_lower ON t (lower(name));`,
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := filterOps(plan, OpCreateIndex)
				if len(ops) != 1 {
					t.Fatalf("expected 1 CreateIndex op, got %d", len(ops))
				}
				sql := ops[0].SQL
				if !strings.Contains(sql, "lower") {
					t.Errorf("expected expression with lower(), got %s", sql)
				}
			},
		},
		{
			name:    "CREATE INDEX with DESC/NULLS FIRST options",
			fromSQL: `CREATE TABLE t (id int, created_at timestamp);`,
			toSQL: `CREATE TABLE t (id int, created_at timestamp);
CREATE INDEX idx_t_created ON t (created_at DESC NULLS FIRST);`,
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := filterOps(plan, OpCreateIndex)
				if len(ops) != 1 {
					t.Fatalf("expected 1 CreateIndex op, got %d", len(ops))
				}
				sql := ops[0].SQL
				if !strings.Contains(sql, "DESC") {
					t.Errorf("expected DESC, got %s", sql)
				}
				if !strings.Contains(sql, "NULLS FIRST") {
					t.Errorf("expected NULLS FIRST, got %s", sql)
				}
			},
		},
		{
			name: "DROP INDEX for removed index",
			fromSQL: `CREATE TABLE t (id int, name text);
CREATE INDEX idx_t_name ON t (name);`,
			toSQL: `CREATE TABLE t (id int, name text);`,
			check: func(t *testing.T, plan *MigrationPlan) {
				ops := filterOps(plan, OpDropIndex)
				if len(ops) != 1 {
					t.Fatalf("expected 1 DropIndex op, got %d", len(ops))
				}
				sql := ops[0].SQL
				if !strings.Contains(sql, "DROP INDEX") {
					t.Errorf("expected DROP INDEX, got %s", sql)
				}
				if !strings.Contains(sql, "idx_t_name") {
					t.Errorf("expected index name idx_t_name, got %s", sql)
				}
			},
		},
		{
			name: "Modified index generated as DROP + CREATE",
			fromSQL: `CREATE TABLE t (id int, name text, email text);
CREATE INDEX idx_t_name ON t (name);`,
			toSQL: `CREATE TABLE t (id int, name text, email text);
CREATE INDEX idx_t_name ON t (name, email);`,
			check: func(t *testing.T, plan *MigrationPlan) {
				drops := filterOps(plan, OpDropIndex)
				creates := filterOps(plan, OpCreateIndex)
				if len(drops) != 1 {
					t.Fatalf("expected 1 DropIndex op, got %d", len(drops))
				}
				if len(creates) != 1 {
					t.Fatalf("expected 1 CreateIndex op, got %d", len(creates))
				}
				// DROP should come before CREATE in the plan.
				dropIdx := -1
				createIdx := -1
				for i, op := range plan.Ops {
					if op.Type == OpDropIndex && dropIdx == -1 {
						dropIdx = i
					}
					if op.Type == OpCreateIndex && createIdx == -1 {
						createIdx = i
					}
				}
				if dropIdx > createIdx {
					t.Errorf("DROP should come before CREATE, drop at %d, create at %d", dropIdx, createIdx)
				}
			},
		},
		{
			name:    "PK/UNIQUE backing indexes NOT generated",
			fromSQL: `CREATE TABLE t (id int);`,
			toSQL:   `CREATE TABLE t (id int PRIMARY KEY);`,
			check: func(t *testing.T, plan *MigrationPlan) {
				indexOps := filterOps(plan, OpCreateIndex)
				dropIndexOps := filterOps(plan, OpDropIndex)
				if len(indexOps) != 0 {
					t.Errorf("expected no CreateIndex ops for PK backing index, got %d: %s", len(indexOps), indexOps[0].SQL)
				}
				if len(dropIndexOps) != 0 {
					t.Errorf("expected no DropIndex ops for PK backing index, got %d", len(dropIndexOps))
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

func filterOps(plan *MigrationPlan, typ MigrationOpType) []MigrationOp {
	var result []MigrationOp
	for _, op := range plan.Ops {
		if op.Type == typ {
			result = append(result, op)
		}
	}
	return result
}
