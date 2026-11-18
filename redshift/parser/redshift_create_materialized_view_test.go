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
			sql:  "CREATE MATERIALIZED VIEW auto_refresh_interval_mv AUTO REFRESH YES INTERVAL 30 AS SELECT 1 AS id;",
			name: "auto_refresh_interval_mv",
			check: func(t *testing.T, opts *nodes.List) {
				assertDefElemBool(t, opts, "auto_refresh", true)
				assertDefElemInt(t, opts, "auto_refresh_interval", 30)
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

// TestRedshiftCreateMaterializedViewIcebergParse pins the Iceberg form of
// r_CREATE_MATERIALIZED_VIEW: USING ICEBERG takes LOCATION, PARTITIONED BY and
// TABLE PROPERTIES, each optional, before AS.
func TestRedshiftCreateMaterializedViewIcebergParse(t *testing.T) {
	stmt := firstCreateMaterializedViewStmt(t, "CREATE MATERIALIZED VIEW spectrum_sales USING ICEBERG LOCATION 's3://bucket/sales/' PARTITIONED BY (sale_date, YEAR(sale_date), BUCKET(16, customer_id), TRUNCATE(100, order_id)) TABLE PROPERTIES ('format' = 'parquet', 'write_target_block_size_mb' = '512') AS SELECT 1 AS id;")
	if stmt.Into.AccessMethod != "iceberg" {
		t.Fatalf("access method = %q, want iceberg", stmt.Into.AccessMethod)
	}
	assertDefElemString(t, stmt.Into.Options, "location", "s3://bucket/sales/")

	transforms := findDefElem(stmt.Into.Options, "partitioned_by")
	if transforms == nil {
		t.Fatalf("expected partitioned_by in %#v", stmt.Into.Options)
	}
	list, ok := transforms.Arg.(*nodes.List)
	if !ok || len(list.Items) != 4 {
		t.Fatalf("partitioned_by = %#v, want four transforms", transforms.Arg)
	}
	column, ok := list.Items[0].(*nodes.ColumnRef)
	if !ok || len(column.Fields.Items) != 1 || column.Fields.Items[0].(*nodes.String).Str != "sale_date" {
		t.Fatalf("first transform = %#v, want column sale_date", list.Items[0])
	}
	transform, ok := list.Items[1].(*nodes.FuncCall)
	if !ok {
		t.Fatalf("second transform = %T, want FuncCall", list.Items[1])
	}
	if got := transform.Funcname.Items[0].(*nodes.String).Str; got != "year" {
		t.Fatalf("second transform = %q, want year", got)
	}
	bucket, ok := list.Items[2].(*nodes.FuncCall)
	if !ok || len(bucket.Args.Items) != 2 {
		t.Fatalf("third transform = %#v, want BUCKET(number, column)", list.Items[2])
	}
	if got := bucket.Funcname.Items[0].(*nodes.String).Str; got != "bucket" {
		t.Fatalf("third transform = %q, want bucket", got)
	}
	if got := list.Items[3].(*nodes.FuncCall).Funcname.Items[0].(*nodes.String).Str; got != "truncate" {
		t.Fatalf("fourth transform = %q, want truncate", got)
	}

	properties := findDefElem(stmt.Into.Options, "table_properties")
	if properties == nil {
		t.Fatalf("expected table_properties in %#v", stmt.Into.Options)
	}
	props, ok := properties.Arg.(*nodes.List)
	if !ok || len(props.Items) != 2 {
		t.Fatalf("table_properties = %#v, want two pairs", properties.Arg)
	}
	first, ok := props.Items[0].(*nodes.DefElem)
	if !ok || first.Defname != "format" || first.Arg.(*nodes.String).Str != "parquet" {
		t.Fatalf("first property = %#v, want format = parquet", props.Items[0])
	}
}

// TestRedshiftCreateMaterializedViewIcebergRejectsOtherTransforms holds the
// partition_transform list to what the docs' Parameters section names: a
// column, YEAR, MONTH, DAY, HOUR, BUCKET and TRUNCATE of one, and nothing else.
func TestRedshiftCreateMaterializedViewIcebergRejectsOtherTransforms(t *testing.T) {
	for _, sql := range []string{
		"CREATE MATERIALIZED VIEW a USING ICEBERG PARTITIONED BY(YES(b)) AS SELECT 1",
		"CREATE MATERIALIZED VIEW a USING ICEBERG PARTITIONED BY(MSK(b)) AS SELECT 1",
		"CREATE MATERIALIZED VIEW a USING ICEBERG PARTITIONED BY(YEAR()) AS SELECT 1",
		"CREATE MATERIALIZED VIEW a USING ICEBERG PARTITIONED BY(BUCKET(b, b)) AS SELECT 1",
		"CREATE MATERIALIZED VIEW a USING ICEBERG TABLE PROPERTIES('s') AS SELECT 1",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) succeeded, want a syntax error", sql)
		}
	}
	// A column named like a transform is a column when no '(' follows it.
	if _, err := Parse("CREATE MATERIALIZED VIEW a USING ICEBERG PARTITIONED BY(bucket, truncate) AS SELECT 1"); err != nil {
		t.Errorf("Parse of bare column transforms: %v", err)
	}
}
