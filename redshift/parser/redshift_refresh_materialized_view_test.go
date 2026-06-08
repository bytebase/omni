package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

func TestRedshiftRefreshMaterializedViewBehaviorParse(t *testing.T) {
	tests := []struct {
		sql      string
		schema   string
		name     string
		behavior int
	}{
		{
			sql:      "REFRESH MATERIALIZED VIEW tickets_mv;",
			name:     "tickets_mv",
			behavior: int(nodes.DROP_RESTRICT),
		},
		{
			sql:      "REFRESH MATERIALIZED VIEW myschema.products_mv RESTRICT;",
			schema:   "myschema",
			name:     "products_mv",
			behavior: int(nodes.DROP_RESTRICT),
		},
		{
			sql:      `REFRESH MATERIALIZED VIEW "sales-data"."monthly_metrics_v2" CASCADE;`,
			schema:   "sales-data",
			name:     "monthly_metrics_v2",
			behavior: int(nodes.DROP_CASCADE),
		},
	}

	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			stmt := firstRefreshMaterializedViewStmt(t, tt.sql)
			if stmt.Relation.Schemaname != tt.schema || stmt.Relation.Relname != tt.name {
				t.Fatalf("expected relation %q.%q, got %#v", tt.schema, tt.name, stmt.Relation)
			}
			if stmt.Behavior != tt.behavior {
				t.Fatalf("expected behavior %d, got %d", tt.behavior, stmt.Behavior)
			}
		})
	}
}

func firstRefreshMaterializedViewStmt(t *testing.T, sql string) *nodes.RefreshMatViewStmt {
	t.Helper()
	tree, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if len(tree.Items) != 1 {
		t.Fatalf("expected one statement, got %d", len(tree.Items))
	}
	raw, ok := tree.Items[0].(*nodes.RawStmt)
	if !ok {
		t.Fatalf("expected RawStmt, got %T", tree.Items[0])
	}
	stmt, ok := raw.Stmt.(*nodes.RefreshMatViewStmt)
	if !ok {
		t.Fatalf("expected RefreshMatViewStmt, got %T", raw.Stmt)
	}
	return stmt
}
