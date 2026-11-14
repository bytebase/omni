package catalog

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/pg/internal/testcontainer"
)

// TestContainerMetadataMigration checks the migration between two synced
// snapshots on PostgreSQL: one database starts as the old snapshot and runs
// LoadMetadata, Diff, and GenerateMigration's plan, another is the new
// snapshot as DatabaseDefinition renders it, and their catalogs must match.
func TestContainerMetadataMigration(t *testing.T) {
	public := func(s *metadata.SchemaMetadata) *metadata.DatabaseSchemaMetadata {
		s.Name = "public"
		return &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{s}}
	}
	col := func(name, typ string) *metadata.ColumnMetadata {
		return &metadata.ColumnMetadata{Name: name, Type: typ}
	}
	orders := func(partitions ...*metadata.TablePartitionMetadata) *metadata.DatabaseSchemaMetadata {
		return public(&metadata.SchemaMetadata{Tables: []*metadata.TableMetadata{{
			Name:       "orders",
			Columns:    []*metadata.ColumnMetadata{col("id", "integer"), col("created_at", "date")},
			Partitions: partitions,
		}}})
	}
	partition := func(name, bound string) *metadata.TablePartitionMetadata {
		return &metadata.TablePartitionMetadata{Name: name, Type: metadata.TablePartitionMetadata_RANGE, Expression: "RANGE (created_at)", Value: bound}
	}
	// sales is partitioned with a key and indexes, which each partition
	// holds attached to the parent's, as sync reports them.
	index := func(table, name, columns string, primary bool, parent string) *metadata.IndexMetadata {
		def := "CREATE INDEX " + name + " ON public." + table + " USING btree (" + columns + ");"
		if primary {
			def = "CREATE UNIQUE INDEX " + name + " ON public." + table + " USING btree (" + columns + ");"
		}
		idx := &metadata.IndexMetadata{
			Name: name, Expressions: strings.Split(columns, ", "), Type: "btree", Definition: def,
			Unique: primary, Primary: primary, IsConstraint: primary,
		}
		if parent != "" {
			idx.ParentIndexSchema, idx.ParentIndexName = "public", parent
		}
		return idx
	}
	salesPartition := func(name, bound string) *metadata.TablePartitionMetadata {
		return &metadata.TablePartitionMetadata{
			Name: name, Type: metadata.TablePartitionMetadata_RANGE, Expression: "RANGE (d)", Value: bound,
			Indexes: []*metadata.IndexMetadata{
				index(name, name+"_pkey", "id, d", true, "sales_pkey"),
				index(name, name+"_region_idx", "region", false, "sales_region_idx"),
				index(name, name+"_lower_idx", "lower(region)", false, "sales_lower_idx"),
			},
		}
	}
	sales := func(partitions ...*metadata.TablePartitionMetadata) *metadata.DatabaseSchemaMetadata {
		return public(&metadata.SchemaMetadata{Tables: []*metadata.TableMetadata{{
			Name:    "sales",
			Columns: []*metadata.ColumnMetadata{col("id", "integer"), col("region", "text"), col("d", "date")},
			Indexes: []*metadata.IndexMetadata{
				index("sales", "sales_pkey", "id, d", true, ""),
				index("sales", "sales_region_idx", "region", false, ""),
				index("sales", "sales_lower_idx", "lower(region)", false, ""),
			},
			Partitions: partitions,
		}}})
	}
	orderSeq := func(owner string) *metadata.SequenceMetadata {
		seq := &metadata.SequenceMetadata{Name: "order_seq", DataType: "integer", Start: "1", MinValue: "1", MaxValue: "2147483647", Increment: "1", CacheSize: "1"}
		if owner != "" {
			seq.OwnerTable, seq.OwnerColumn = owner, "id"
		}
		return seq
	}
	trgm := func(version string) *metadata.DatabaseSchemaMetadata {
		m := public(&metadata.SchemaMetadata{})
		if version != "" {
			m.Extensions = []*metadata.ExtensionMetadata{{Name: "pg_trgm", Schema: "public", Version: version}}
		}
		return m
	}
	addr := func(comment string, attrs ...*metadata.CompositeTypeAttribute) *metadata.DatabaseSchemaMetadata {
		return public(&metadata.SchemaMetadata{CompositeTypes: []*metadata.CompositeTypeMetadata{{Name: "addr", Comment: comment, Attributes: attrs}}})
	}
	// plan and issue are cut from Bytebase's own metadata schema, as sync
	// reports them: a column inserted ahead of the key columns renumbers them.
	plan := func(editor bool) *metadata.DatabaseSchemaMetadata {
		columns := []*metadata.ColumnMetadata{col("id", "bigint"), col("creator", "text")}
		if editor {
			columns = append(columns, &metadata.ColumnMetadata{Name: "last_plan_editor", Type: "text", Nullable: true})
		}
		columns = append(columns, col("project", "text"), col("config", "jsonb"))
		return public(&metadata.SchemaMetadata{Tables: []*metadata.TableMetadata{
			{
				Name:    "plan",
				Columns: columns,
				Indexes: []*metadata.IndexMetadata{
					{
						Name: "plan_pkey", Expressions: []string{"project", "id"}, Descending: []bool{false, false}, Type: "btree", Unique: true, Primary: true, IsConstraint: true,
						Definition:   "CREATE UNIQUE INDEX plan_pkey ON public.plan USING btree (project, id);",
						OpclassNames: []string{"text_ops", "int8_ops"}, OpclassDefaults: []bool{true, true},
					},
					{
						Name: "idx_plan_project", Expressions: []string{"project"}, Descending: []bool{false}, Type: "btree",
						Definition:   "CREATE INDEX idx_plan_project ON public.plan USING btree (project);",
						OpclassNames: []string{"text_ops"}, OpclassDefaults: []bool{true},
					},
					{
						Name: "idx_plan_config_has_rollout", Expressions: []string{"(config ->> 'hasRollout'::text)"}, Descending: []bool{false}, Type: "btree",
						Definition:   "CREATE INDEX idx_plan_config_has_rollout ON public.plan USING btree (((config ->> 'hasRollout'::text)));",
						OpclassNames: []string{"text_ops"}, OpclassDefaults: []bool{true},
					},
					{
						Name: "idx_plan_config", Expressions: []string{"config"}, Descending: []bool{false}, Type: "gin",
						Definition:   "CREATE INDEX idx_plan_config ON public.plan USING gin (config jsonb_path_ops);",
						OpclassNames: []string{"jsonb_path_ops"}, OpclassDefaults: []bool{false},
					},
				},
			},
			{
				Name:    "issue",
				Columns: []*metadata.ColumnMetadata{col("id", "bigint"), col("project", "text"), {Name: "plan_id", Type: "bigint", Nullable: true}},
				ForeignKeys: []*metadata.ForeignKeyMetadata{{
					Name: "issue_project_plan_id_fkey", Columns: []string{"project", "plan_id"},
					ReferencedSchema: "public", ReferencedTable: "plan", ReferencedColumns: []string{"project", "id"},
					OnDelete: "NO ACTION", OnUpdate: "NO ACTION", MatchType: "SIMPLE",
				}},
			},
		}})
	}

	tests := []struct {
		name     string
		from, to *metadata.DatabaseSchemaMetadata
	}{
		{
			name: "partition",
			from: orders(partition("orders_2023", "FOR VALUES FROM ('2023-01-01') TO ('2024-01-01')")),
			to:   orders(partition("orders_2024", "FOR VALUES FROM ('2024-01-01') TO ('2025-01-01')"), partition("orders_other", "DEFAULT")),
		},
		{
			name: "partition_with_indexes",
			from: sales(salesPartition("sales_2024", "FOR VALUES FROM ('2024-01-01') TO ('2025-01-01')")),
			to: sales(salesPartition("sales_2024", "FOR VALUES FROM ('2024-01-01') TO ('2025-01-01')"),
				salesPartition("sales_2025", "FOR VALUES FROM ('2025-01-01') TO ('2026-01-01')")),
		},
		{
			name: "owned_sequence",
			from: public(&metadata.SchemaMetadata{Sequences: []*metadata.SequenceMetadata{orderSeq("")}}),
			to: public(&metadata.SchemaMetadata{
				Tables:    []*metadata.TableMetadata{{Name: "orders", Columns: []*metadata.ColumnMetadata{col("id", "bigint")}}},
				Sequences: []*metadata.SequenceMetadata{orderSeq("orders")},
			}),
		},
		{
			name: "released_sequence",
			from: public(&metadata.SchemaMetadata{
				Tables:    []*metadata.TableMetadata{{Name: "orders", Columns: []*metadata.ColumnMetadata{col("id", "bigint")}}},
				Sequences: []*metadata.SequenceMetadata{orderSeq("orders")},
			}),
			to: public(&metadata.SchemaMetadata{Sequences: []*metadata.SequenceMetadata{orderSeq("")}}),
		},
		{name: "create_extension", from: trgm(""), to: trgm("1.6")},
		{name: "update_extension", from: trgm("1.5"), to: trgm("1.6")},
		{name: "drop_extension", from: trgm("1.6"), to: trgm("")},
		{
			name: "extension_member",
			from: public(&metadata.SchemaMetadata{CompositeTypes: []*metadata.CompositeTypeMetadata{{
				Name: "ext_owned", SkipDump: true, Attributes: []*metadata.CompositeTypeAttribute{{Name: "x", Type: "integer"}},
			}}}),
			to: public(&metadata.SchemaMetadata{}),
		},
		{
			name: "composite_collation",
			from: addr("", &metadata.CompositeTypeAttribute{Name: "street", Type: "text"}),
			to: addr("", &metadata.CompositeTypeAttribute{Name: "street", Type: "text", Collation: `"C"`},
				&metadata.CompositeTypeAttribute{Name: "zip", Type: "text", Collation: `"C"`}),
		},
		{
			name: "composite_comments",
			from: addr("old type comment", &metadata.CompositeTypeAttribute{Name: "street", Type: "text", Comment: "old attr comment"}),
			to:   addr("new type comment", &metadata.CompositeTypeAttribute{Name: "street", Type: "text", Comment: "new attr comment"}),
		},
		{name: "create_with_expression_index", from: public(&metadata.SchemaMetadata{}), to: plan(false)},
		{name: "column_ahead_of_key", from: plan(false), to: plan(true)},
		{name: "column_ahead_of_key_dropped", from: plan(true), to: plan(false)},
	}

	ctx := context.Background()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			migrated := testcontainer.OpenDatabase(t, testcontainer.Default, "metadata_"+tc.name+"_migrated")
			expected := testcontainer.OpenDatabase(t, testcontainer.Default, "metadata_"+tc.name+"_expected")
			execDefinition(ctx, t, migrated, tc.from)
			execDefinition(ctx, t, expected, tc.to)

			load := func(m *metadata.DatabaseSchemaMetadata) *Catalog {
				c := New()
				report, err := c.LoadMetadata(ctx, m, LoadMetadataOptions{Full: true})
				if err != nil {
					t.Fatalf("LoadMetadata: %v", err)
				}
				requireCleanReport(t, report)
				return c
			}
			from, to := load(tc.from), load(tc.to)
			migration := GenerateMigration(from, to, Diff(from, to)).SQL()
			if migration != "" {
				if _, err := migrated.ExecContext(ctx, migration); err != nil {
					t.Fatalf("migration failed: %v\n%s", err, migration)
				}
			}
			if got, want := dumpUserCatalog(ctx, t, migrated), dumpUserCatalog(ctx, t, expected); got != want {
				t.Errorf("migrated database differs from the target\nmigration:\n%s\ngot:\n%s\nwant:\n%s", migration, got, want)
			}
		})
	}
}

func execDefinition(ctx context.Context, t *testing.T, db *sql.DB, meta *metadata.DatabaseSchemaMetadata) {
	t.Helper()
	ddl, err := DatabaseDefinition(meta, DefinitionOptions{})
	if err != nil {
		t.Fatalf("DatabaseDefinition: %v", err)
	}
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		t.Fatalf("restore snapshot: %v\n%s", err, ddl)
	}
}

// dumpUserCatalog lists what a migration can change outside extensions, one
// line per object, sorted: columns by name, since ADD COLUMN appends.
func dumpUserCatalog(ctx context.Context, t *testing.T, db *sql.DB) string {
	t.Helper()
	rows, err := db.QueryContext(ctx, `
WITH ns AS (
  SELECT oid, nspname FROM pg_namespace
  WHERE nspname NOT IN ('pg_catalog', 'information_schema') AND nspname NOT LIKE 'pg_toast%' AND nspname NOT LIKE 'pg_temp%'
), member AS (
  SELECT objid FROM pg_depend WHERE deptype = 'e'
)
SELECT format('schema %s comment=%s', nspname, obj_description(oid, 'pg_namespace')) FROM ns
UNION ALL
SELECT format('extension %s %s in %s', e.extname, e.extversion, n.nspname)
  FROM pg_extension e JOIN pg_namespace n ON n.oid = e.extnamespace WHERE e.extname <> 'plpgsql'
UNION ALL
SELECT format('relation %s.%s kind=%s key=%s bound=%s index=%s comment=%s', ns.nspname, c.relname, c.relkind,
    pg_get_partkeydef(c.oid), pg_get_expr(c.relpartbound, c.oid),
    CASE WHEN c.relkind IN ('i', 'I') THEN pg_get_indexdef(c.oid) END, obj_description(c.oid, 'pg_class'))
  FROM pg_class c JOIN ns ON ns.oid = c.relnamespace
  WHERE c.relkind IN ('r', 'p', 'v', 'm', 'S', 'c', 'i', 'I') AND c.oid NOT IN (SELECT objid FROM member)
UNION ALL
SELECT format('column %s.%s.%s %s collation=%s notnull=%s default=%s comment=%s', ns.nspname, c.relname, a.attname,
    format_type(a.atttypid, a.atttypmod), co.collname, a.attnotnull, pg_get_expr(d.adbin, d.adrelid), col_description(c.oid, a.attnum))
  FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid JOIN ns ON ns.oid = c.relnamespace
  JOIN pg_type ty ON ty.oid = a.atttypid
  LEFT JOIN pg_collation co ON co.oid = a.attcollation AND a.attcollation <> ty.typcollation
  LEFT JOIN pg_attrdef d ON d.adrelid = a.attrelid AND d.adnum = a.attnum
  WHERE a.attnum > 0 AND NOT a.attisdropped AND c.relkind IN ('r', 'p', 'v', 'm', 'c') AND c.oid NOT IN (SELECT objid FROM member)
UNION ALL
SELECT format('constraint %s.%s.%s %s comment=%s', ns.nspname, c.relname, k.conname, pg_get_constraintdef(k.oid), obj_description(k.oid, 'pg_constraint'))
  FROM pg_constraint k JOIN pg_class c ON c.oid = k.conrelid JOIN ns ON ns.oid = c.relnamespace
UNION ALL
SELECT format('sequence %s.%s %s start=%s increment=%s min=%s max=%s cache=%s cycle=%s owner=%s', ns.nspname, c.relname,
    format_type(s.seqtypid, NULL), s.seqstart, s.seqincrement, s.seqmin, s.seqmax, s.seqcache, s.seqcycle,
    (SELECT format('%s.%s', o.relname, a.attname) FROM pg_depend d
       JOIN pg_class o ON o.oid = d.refobjid JOIN pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
       WHERE d.classid = 'pg_class'::regclass AND d.objid = c.oid AND d.refclassid = 'pg_class'::regclass AND d.deptype IN ('a', 'i')))
  FROM pg_sequence s JOIN pg_class c ON c.oid = s.seqrelid JOIN ns ON ns.oid = c.relnamespace
UNION ALL
SELECT format('type %s.%s %s comment=%s', ns.nspname, ty.typname, ty.typtype, obj_description(ty.oid, 'pg_type'))
  FROM pg_type ty JOIN ns ON ns.oid = ty.typnamespace
  WHERE (ty.typtype IN ('e', 'd') OR ty.typtype = 'c' AND (SELECT relkind FROM pg_class WHERE oid = ty.typrelid) = 'c')
    AND ty.oid NOT IN (SELECT objid FROM member)
ORDER BY 1`)
	if err != nil {
		t.Fatalf("dump catalog: %v", err)
	}
	defer rows.Close()
	var lines []string
	for rows.Next() {
		var line string
		if err := rows.Scan(&line); err != nil {
			t.Fatalf("dump catalog: %v", err)
		}
		lines = append(lines, line)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("dump catalog: %v", err)
	}
	return strings.Join(lines, "\n")
}

// TestContainerMigrationFromSQL checks plans between catalogs loaded from SQL
// on PostgreSQL, for what the metadata snapshots above do not reach.
func TestContainerMigrationFromSQL(t *testing.T) {
	ctr := startPGContainer(t)
	partitioned := `CREATE TABLE p (id int) PARTITION BY RANGE (id);
CREATE TABLE child PARTITION OF p FOR VALUES FROM (0) TO (10);
CREATE INDEX pi ON ONLY p (id);
CREATE INDEX ci ON child (id);`
	tests := []struct {
		name, before, after string
	}{
		{
			name:   "attach_partition_index",
			before: partitioned,
			after:  partitioned + "\nALTER INDEX pi ATTACH PARTITION ci;",
		},
		{
			name:   "collation_under_view",
			before: `CREATE TABLE t (x text COLLATE "C"); CREATE VIEW v AS SELECT x FROM t;`,
			after:  `CREATE TABLE t (x text COLLATE "POSIX"); CREATE VIEW v AS SELECT x FROM t;`,
		},
		{
			name:   "owner_takes_its_sequence",
			before: `CREATE SEQUENCE s; CREATE TABLE t (id int DEFAULT nextval('s')); ALTER SEQUENCE s OWNED BY t.id;`,
			after:  ``,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertOracleRoundtrip(t, ctr, tt.before, tt.after)
		})
	}
}
