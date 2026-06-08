package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

func TestRedshiftCreateMaterializedViewOptionsParse(t *testing.T) {
	tests := []struct {
		sql    string
		schema string
		name   string
		check  func(t *testing.T, opts *nodes.List)
	}{
		{
			sql:  "CREATE MATERIALIZED VIEW backup_mv BACKUP YES AS SELECT 1 AS id;",
			name: "backup_mv",
			check: func(t *testing.T, opts *nodes.List) {
				assertDefElemString(t, opts, "backup", "yes")
			},
		},
		{
			sql:  "CREATE MATERIALIZED VIEW even_dist_mv DISTSTYLE EVEN AS SELECT 1 AS id;",
			name: "even_dist_mv",
			check: func(t *testing.T, opts *nodes.List) {
				assertDefElemString(t, opts, "diststyle", "even")
			},
		},
		{
			sql:  "CREATE MATERIALIZED VIEW customer_dist_mv DISTKEY(customer_id) AS SELECT customer_id FROM orders;",
			name: "customer_dist_mv",
			check: func(t *testing.T, opts *nodes.List) {
				assertDefElemList(t, opts, "distkey")
			},
		},
		{
			sql:  "CREATE MATERIALIZED VIEW compound_sorted_mv COMPOUND SORTKEY(customer_id, order_date) AS SELECT customer_id, order_date FROM orders;",
			name: "compound_sorted_mv",
			check: func(t *testing.T, opts *nodes.List) {
				assertDefElemString(t, opts, "sortstyle", "compound")
				assertDefElemList(t, opts, "sortkey")
			},
		},
		{
			sql:  "CREATE MATERIALIZED VIEW auto_refresh_mv AUTO REFRESH YES AS SELECT 1 AS id;",
			name: "auto_refresh_mv",
			check: func(t *testing.T, opts *nodes.List) {
				assertDefElemBool(t, opts, "auto_refresh", true)
			},
		},
		{
			sql:    "CREATE MATERIALIZED VIEW sales.full_featured_mv BACKUP NO DISTSTYLE KEY DISTKEY(customer_id) SORTKEY(order_date, product_id) AUTO REFRESH NO AS SELECT customer_id, order_date, product_id FROM order_details;",
			schema: "sales",
			name:   "full_featured_mv",
			check: func(t *testing.T, opts *nodes.List) {
				assertDefElemString(t, opts, "backup", "no")
				assertDefElemString(t, opts, "diststyle", "key")
				assertDefElemList(t, opts, "distkey")
				assertDefElemList(t, opts, "sortkey")
				assertDefElemBool(t, opts, "auto_refresh", false)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			stmt := firstCreateMaterializedViewStmt(t, tt.sql)
			if stmt.Objtype != nodes.OBJECT_MATVIEW {
				t.Fatalf("expected materialized view object type, got %v", stmt.Objtype)
			}
			if stmt.Into == nil || stmt.Into.Rel == nil {
				t.Fatalf("expected materialized view target, got %#v", stmt.Into)
			}
			if stmt.Into.Rel.Schemaname != tt.schema || stmt.Into.Rel.Relname != tt.name {
				t.Fatalf("expected target %q.%q, got %#v", tt.schema, tt.name, stmt.Into.Rel)
			}
			tt.check(t, stmt.Into.Options)
		})
	}
}

func firstCreateMaterializedViewStmt(t *testing.T, sql string) *nodes.CreateTableAsStmt {
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
	stmt, ok := raw.Stmt.(*nodes.CreateTableAsStmt)
	if !ok {
		t.Fatalf("expected CreateTableAsStmt, got %T", raw.Stmt)
	}
	return stmt
}
