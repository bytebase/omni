package catalog

import (
	"context"
	"database/sql"
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bytebase/omni/metadata"
)

const snapshotDB = "testdb"

func loadTiDBSnapshot(t *testing.T, meta *metadata.DatabaseSchemaMetadata) (*Catalog, *LoadMetadataReport) {
	t.Helper()
	c := New()
	report, err := c.LoadMetadata(context.Background(), meta)
	if err != nil {
		t.Fatalf("LoadMetadata: %v", err)
	}
	return c, report
}

// requireReport fails unless exactly the given objects were replaced by
// stand-ins and nothing is missing.
func requireReport(t *testing.T, report *LoadMetadataReport, degraded ...string) {
	t.Helper()
	var got []string
	for k := range report.Degraded {
		got = append(got, k)
	}
	slices.Sort(got)
	slices.Sort(degraded)
	if !slices.Equal(got, degraded) || len(report.Missing) != 0 {
		t.Errorf("degraded = %v, missing = %v; want degraded %v and nothing missing", report.Degraded, report.Missing, degraded)
	}
}

func requireTable(t *testing.T, c *Catalog, name string) *Table {
	t.Helper()
	db := c.GetDatabase(snapshotDB)
	if db == nil || db.Tables[strings.ToLower(name)] == nil {
		t.Fatalf("table %q not in catalog", name)
	}
	return db.Tables[strings.ToLower(name)]
}

func requireView(t *testing.T, c *Catalog, name string) *View {
	t.Helper()
	db := c.GetDatabase(snapshotDB)
	if db == nil || db.Views[strings.ToLower(name)] == nil {
		t.Fatalf("view %q not in catalog", name)
	}
	return db.Views[strings.ToLower(name)]
}

func snapshot(s *metadata.SchemaMetadata) *metadata.DatabaseSchemaMetadata {
	return &metadata.DatabaseSchemaMetadata{Name: snapshotDB, Schemas: []*metadata.SchemaMetadata{s}}
}

func TestLoadMetadataInstallsTables(t *testing.T) {
	c := New()
	c.SetForeignKeyChecks(true)
	meta := snapshot(&metadata.SchemaMetadata{Tables: []*metadata.TableMetadata{
		{
			// Sorts first and references z_parent.
			Name:           "a_child",
			Charset:        "utf8mb4",
			Comment:        "children",
			PrimaryKeyType: "CLUSTERED",
			Columns: []*metadata.ColumnMetadata{
				{Name: "id", Type: "bigint unsigned", Default: autoIncrementDefault},
				{Name: "parent_id", Type: "bigint unsigned"},
				{Name: "email", Type: "varchar(255)", Nullable: true, CharacterSet: "latin1"},
				{Name: "created_at", Type: "timestamp", Default: "CURRENT_TIMESTAMP", OnUpdate: "CURRENT_TIMESTAMP"},
				{Name: "total", Type: "int", Nullable: true, Generation: &metadata.GenerationMetadata{
					Type: metadata.GenerationMetadata_TYPE_STORED, Expression: "parent_id + 1",
				}},
				// Sync reports DEFAULT NULL for every nullable column, including
				// types that show no default.
				{Name: "payload", Type: "json", Nullable: true, Default: "NULL"},
				{Name: "notes", Type: "text", Nullable: true, Default: "NULL"},
				{Name: "name", Type: "varchar(255)", Nullable: true, Default: "NULL"},
				// Sync quotes a literal default, so it loads as a literal.
				{Name: "label", Type: "varchar(32)", Nullable: true, Default: "'hello world'"},
			},
			Indexes: []*metadata.IndexMetadata{
				{Name: "PRIMARY", Type: "BTREE", Primary: true, Unique: true, Expressions: []string{"id"}},
				{Name: "idx_parent_id", Type: "BTREE", Expressions: []string{"parent_id"}},
				{Name: "idx_lower_name", Type: "BTREE", Expressions: []string{"(lower(`name`))"}},
				{Name: "idx_prefix", Type: "BTREE", Expressions: []string{"name"}, KeyLength: []int64{10}, Descending: []bool{true}},
				{Name: "idx_shown", Type: "BTREE", Expressions: []string{"parent_id"}, Visible: true},
				{Name: "idx_hidden", Type: "BTREE", Expressions: []string{"parent_id"}, Comment: "rarely used"},
			},
			ForeignKeys: []*metadata.ForeignKeyMetadata{{
				Name:              "fk_child_parent",
				Columns:           []string{"parent_id"},
				ReferencedTable:   "z_parent",
				ReferencedColumns: []string{"id"},
				OnDelete:          "CASCADE",
			}},
		},
		{
			Name:           "z_parent",
			Columns:        []*metadata.ColumnMetadata{{Name: "id", Type: "bigint unsigned", Default: autoIncrementDefault}},
			Indexes:        []*metadata.IndexMetadata{{Name: "PRIMARY", Type: "BTREE", Primary: true, Unique: true, Expressions: []string{"id"}}},
			PrimaryKeyType: "NONCLUSTERED",
		},
	}})
	meta.CharacterSet, meta.Collation = "latin1", "latin1_bin"
	report, err := c.LoadMetadata(context.Background(), meta)
	if err != nil {
		t.Fatalf("LoadMetadata: %v", err)
	}
	requireReport(t, report)
	if c.CurrentDatabase() != snapshotDB || !c.ForeignKeyChecks() {
		t.Errorf("current database %q, foreign key checks %v; want %s and checks back on", c.CurrentDatabase(), c.ForeignKeyChecks(), snapshotDB)
	}

	if db := c.GetDatabase(snapshotDB); db.Charset != "latin1" || db.Collation != "latin1_bin" {
		t.Errorf("database charset %q, collation %q; want the snapshot's latin1, latin1_bin", db.Charset, db.Collation)
	}

	// Loading into an existing database adds to it and takes the snapshot's
	// defaults, so a table that overrides neither resolves the same as it would
	// in an empty catalog.
	second := snapshot(&metadata.SchemaMetadata{Tables: []*metadata.TableMetadata{
		{Name: "u", Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int"}}},
	}})
	second.CharacterSet, second.Collation = "utf8mb4", "utf8mb4_bin"
	if _, err := c.LoadMetadata(context.Background(), second); err != nil {
		t.Fatalf("second LoadMetadata: %v", err)
	}
	requireTable(t, c, "a_child")
	requireTable(t, c, "u")
	if db := c.GetDatabase(snapshotDB); db.Charset != "utf8mb4" || db.Collation != "utf8mb4_bin" {
		t.Errorf("database charset %q, collation %q; want the second snapshot's utf8mb4, utf8mb4_bin", db.Charset, db.Collation)
	}

	child := requireTable(t, c, "a_child")
	if child.Comment != "children" || child.Charset != "utf8mb4" {
		t.Errorf("a_child comment %q, charset %q", child.Comment, child.Charset)
	}
	if id := child.GetColumn("id"); !id.AutoIncrement || id.Nullable {
		t.Errorf("id = auto increment %v, nullable %v", id.AutoIncrement, id.Nullable)
	}
	if email := child.GetColumn("email"); email.Charset != "latin1" {
		t.Errorf("email charset = %q, want latin1", email.Charset)
	}
	if created := child.GetColumn("created_at"); created.Default == nil || created.OnUpdate == "" {
		t.Errorf("created_at lost its default or ON UPDATE: %+v", created)
	}
	if gen := child.GetColumn("total").Generated; gen == nil || !gen.Stored {
		t.Errorf("total generated = %+v, want stored", gen)
	}
	// The loaded expression is the one a CREATE TABLE gives, carrying none of
	// the parentheses the parse probe adds.
	if _, err := c.Exec("CREATE TABLE ddl (parent_id BIGINT UNSIGNED, total INT AS (parent_id + 1) STORED);", nil); err != nil {
		t.Fatalf("CREATE TABLE: %v", err)
	}
	if loaded, fromDDL := child.GetColumn("total").Generated, requireTable(t, c, "ddl").GetColumn("total").Generated; loaded == nil || fromDDL == nil || loaded.Expr != fromDDL.Expr {
		t.Errorf("generated expression = %+v, want the DDL's %+v", loaded, fromDDL)
	}
	for _, name := range []string{"payload", "notes"} {
		if col := child.GetColumn(name); col.Default != nil {
			t.Errorf("%s default = %q, want none", name, *col.Default)
		}
	}
	if name := child.GetColumn("name"); name.Default == nil {
		t.Error("name lost its DEFAULT NULL")
	}
	if label := child.GetColumn("label"); label.Default == nil || !strings.Contains(*label.Default, "hello world") {
		t.Errorf("label default = %v, want the quoted literal", label.Default)
	}
	indexes := make(map[string]*Index)
	for _, idx := range child.Indexes {
		indexes[idx.Name] = idx
	}
	if idx := indexes["idx_lower_name"]; idx == nil || len(idx.Columns) != 1 || !strings.Contains(strings.ToLower(idx.Columns[0].Expr), "lower(") {
		t.Errorf("idx_lower_name = %+v, want the lower() expression", idx)
	}
	if idx := indexes["idx_prefix"]; idx == nil || len(idx.Columns) != 1 || idx.Columns[0].Length != 10 || !idx.Columns[0].Descending {
		t.Errorf("idx_prefix = %+v, want name(10) DESC", idx)
	}
	// The primary key is visible though the snapshot does not say so.
	for name, visible := range map[string]bool{"PRIMARY": true, "idx_shown": true, "idx_hidden": false} {
		if idx := indexes[name]; idx == nil || idx.Visible != visible {
			t.Errorf("%s = %+v, want visible %v", name, idx, visible)
		}
	}
	if idx := indexes["idx_hidden"]; idx == nil || idx.Comment != "rarely used" {
		t.Errorf("idx_hidden = %+v, want its comment", idx)
	}
	// The access method a key without USING gets stays unwritten.
	if idx := indexes["idx_parent_id"]; idx == nil || idx.IndexType != "" {
		t.Errorf("idx_parent_id = %+v, want no index type", idx)
	}
	var fk *Constraint
	for _, con := range child.Constraints {
		if con.Type == ConForeignKey && con.Name == "fk_child_parent" {
			fk = con
		}
	}
	if fk == nil || !strings.EqualFold(fk.OnDelete, "CASCADE") {
		t.Errorf("fk_child_parent = %+v, want ON DELETE CASCADE", fk)
	}
	clustered := func(table string) string {
		for _, con := range requireTable(t, c, table).Constraints {
			if con.Type == ConPrimaryKey && con.Clustered != nil {
				return strconv.FormatBool(*con.Clustered)
			}
		}
		return "unset"
	}
	if got := clustered("a_child") + " " + clustered("z_parent"); got != "true false" {
		t.Errorf("primary keys clustered = %s, want true false", got)
	}
}

func TestLoadMetadataInstallsViewsInAnyOrder(t *testing.T) {
	// Each view reads the next, which installs after it, and m_view's columns
	// come only from analyzing its body.
	c, report := loadTiDBSnapshot(t, snapshot(&metadata.SchemaMetadata{
		Tables: []*metadata.TableMetadata{{
			Name:    "users",
			Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int"}, {Name: "name", Type: "varchar(255)", Nullable: true}},
		}},
		Views: []*metadata.ViewMetadata{
			{Name: "a_view", Definition: "SELECT id FROM m_view"},
			{Name: "m_view", Definition: "SELECT * FROM z_view"},
			{Name: "z_view", Definition: "SELECT id, name FROM users"},
		},
	}))
	requireReport(t, report)
	// The columns come from analyzing the parsed body.
	if v := requireView(t, c, "z_view"); !slices.Equal(v.Columns, []string{"id", "name"}) {
		t.Errorf("z_view columns = %v, want id, name", v.Columns)
	}
	for _, name := range []string{"a_view", "m_view"} {
		if v := requireView(t, c, name); v.AnalyzedQuery == nil {
			t.Errorf("%s was not analyzed once the view it reads installed", name)
		}
	}
}

func TestLoadMetadataKeepsViewColumnListsOnlyWhereNeeded(t *testing.T) {
	c, report := loadTiDBSnapshot(t, snapshot(&metadata.SchemaMetadata{
		Tables: []*metadata.TableMetadata{{Name: "users", Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int"}}}},
		Views: []*metadata.ViewMetadata{
			// Created as CREATE VIEW renamed (user_id) AS SELECT id FROM users.
			{Name: "renamed", Definition: "SELECT id FROM users", Columns: []*metadata.ColumnMetadata{{Name: "user_id"}}},
			{Name: "same", Definition: "SELECT id FROM users", Columns: []*metadata.ColumnMetadata{{Name: "ID"}}},
			{Name: "star", Definition: "SELECT * FROM users", Columns: []*metadata.ColumnMetadata{{Name: "id"}}},
			{Name: "star_qualified", Definition: "SELECT u.* FROM users u", Columns: []*metadata.ColumnMetadata{{Name: "id"}}},
			{Name: "expr", Definition: "SELECT count(*) FROM users", Columns: []*metadata.ColumnMetadata{{Name: "count(*)"}}},
			// A body that does not resolve names nothing of its own. A qualified
			// star leaves one unnamed column behind rather than none.
			{Name: "unresolved", Definition: "SELECT * FROM missing", Columns: []*metadata.ColumnMetadata{{Name: "a"}, {Name: "b"}}},
			{Name: "unresolved_qualified", Definition: "SELECT m.* FROM missing m", Columns: []*metadata.ColumnMetadata{{Name: "a"}}},
			// Created as CREATE VIEW star_renamed (user_id) AS SELECT * FROM users.
			{Name: "star_renamed", Definition: "SELECT * FROM users", Columns: []*metadata.ColumnMetadata{{Name: "user_id"}}},
		},
	}))
	requireReport(t, report)
	for _, name := range []string{"renamed", "star_renamed"} {
		if v := requireView(t, c, name); !v.ExplicitColumns || !slices.Equal(v.Columns, []string{"user_id"}) {
			t.Errorf("%s = explicit %v, columns %v; want the snapshot's user_id", name, v.ExplicitColumns, v.Columns)
		}
	}
	// A body naming no target of its own keeps the columns it resolves to, and
	// the snapshot names what a body that does not resolve leaves unnamed.
	for name, want := range map[string][]string{"same": {"id"}, "star": {"id"}, "star_qualified": {"id"}, "expr": {"COUNT(*)"}, "unresolved": {"a", "b"}, "unresolved_qualified": {"a"}} {
		if v := requireView(t, c, name); v.ExplicitColumns || !slices.Equal(v.Columns, want) {
			t.Errorf("%s = explicit %v, columns %v; want %v inferred", name, v.ExplicitColumns, v.Columns, want)
		}
	}
}

func TestLoadMetadataStandsInForFailedObjects(t *testing.T) {
	c, report := loadTiDBSnapshot(t, snapshot(&metadata.SchemaMetadata{
		Tables: []*metadata.TableMetadata{{
			Name:    "broken",
			Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int"}, {Name: "bad", Type: "varchar((("}},
		}},
		Views: []*metadata.ViewMetadata{
			{Name: "bad_view", Definition: "SELEC nonsense", Columns: []*metadata.ColumnMetadata{{Name: "x"}}},
		},
	}))
	requireReport(t, report, "table:broken", "view:bad_view")
	if cols := requireView(t, c, "bad_view").Columns; !slices.Equal(cols, []string{"x"}) {
		t.Errorf("bad_view stand-in columns = %v, want x", cols)
	}
	cols := requireTable(t, c, "broken").Columns
	if len(cols) != 2 {
		t.Fatalf("stand-in has %d columns, want 2", len(cols))
	}
	for _, col := range cols {
		if !strings.EqualFold(col.DataType, "text") {
			t.Errorf("stand-in column %s is %s, want text", col.Name, col.DataType)
		}
	}
}

func TestStandInStatements(t *testing.T) {
	c := New()
	if _, err := c.LoadMetadata(context.Background(), snapshot(&metadata.SchemaMetadata{})); err != nil {
		t.Fatalf("LoadMetadata: %v", err)
	}
	for name, err := range map[string]error{
		"empty": c.DefineTable(standInTableStmt(&metadata.TableMetadata{Name: "empty"})),
		"pv":    c.DefineView(standInViewStmt(&metadata.ViewMetadata{Name: "pv", Columns: []*metadata.ColumnMetadata{{Name: "x"}, {Name: "y"}}})),
		"dv":    c.DefineView(standInViewStmt(&metadata.ViewMetadata{Name: "dv", DependencyColumns: []*metadata.DependencyColumn{{Column: "id"}}})),
	} {
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
	}
	if cols := requireTable(t, c, "empty").Columns; len(cols) != 1 || cols[0].Name != standInPlaceholder {
		t.Errorf("empty stand-in columns = %+v, want the placeholder", cols)
	}
	if cols := requireView(t, c, "pv").Columns; len(cols) != 2 {
		t.Errorf("pv stand-in has %d columns, want 2", len(cols))
	}
	if cols := requireView(t, c, "dv").Columns; !slices.Equal(cols, []string{"id"}) {
		t.Errorf("dv stand-in columns = %v, want the column it reads", cols)
	}
}

func TestLoadMetadataStopsWhenContextIsDone(t *testing.T) {
	for _, tt := range []struct {
		name   string
		checks int
		meta   *metadata.DatabaseSchemaMetadata
	}{
		{name: "before planning", checks: 0, meta: snapshot(&metadata.SchemaMetadata{})},
		{name: "before an install", checks: 1, meta: snapshot(&metadata.SchemaMetadata{Tables: []*metadata.TableMetadata{{Name: "t"}}})},
		{name: "before a view reinstall", checks: 3, meta: snapshot(&metadata.SchemaMetadata{Views: []*metadata.ViewMetadata{
			{Name: "a", Definition: "SELECT id FROM z"},
			{Name: "z", Definition: "SELECT 1 AS id"},
		}})},
	} {
		if _, err := New().LoadMetadata(&doneAfter{Context: context.Background(), checks: tt.checks}, tt.meta); err != context.Canceled {
			t.Errorf("%s: err = %v, want context.Canceled", tt.name, err)
		}
	}
}

// doneAfter is a context that is done once Err has been called checks times.
type doneAfter struct {
	context.Context
	checks int
}

func (c *doneAfter) Err() error {
	if c.checks == 0 {
		return context.Canceled
	}
	c.checks--
	return nil
}

func TestParseTypeName(t *testing.T) {
	for typ, want := range map[string]string{
		"int":                          "INT",
		"bigint(20) unsigned zerofill": "BIGINT",
		"varchar(255)":                 "VARCHAR",
		"decimal(10,2)":                "DECIMAL",
		"datetime(6)":                  "DATETIME",
		"json":                         "JSON",
		"enum('a','b','c')":            "ENUM",
		"set('x','y')":                 "SET",
	} {
		dt, err := parseTypeName(typ)
		if err != nil || !strings.EqualFold(dt.Name, want) {
			t.Errorf("parseTypeName(%q) = %v, %v; want %s", typ, dt, err, want)
		}
	}
	for _, typ := range []string{"", "not a real type ((("} {
		if _, err := parseTypeName(typ); err == nil {
			t.Errorf("parseTypeName(%q) succeeded, want an error", typ)
		}
	}
}

// requireIncomplete fails unless exactly the given objects were installed
// without part of their definition.
func requireIncomplete(t *testing.T, report *LoadMetadataReport, incomplete ...string) {
	t.Helper()
	got := slices.Sorted(maps.Keys(report.Incomplete))
	slices.Sort(incomplete)
	if !slices.Equal(got, incomplete) {
		t.Errorf("incomplete = %v, want %v (errors %v)", got, incomplete, report.Incomplete)
	}
}

func TestLoadMetadataInstallsPartitioning(t *testing.T) {
	c, report := loadTiDBSnapshot(t, snapshot(&metadata.SchemaMetadata{Tables: []*metadata.TableMetadata{
		{
			Name:    "sales",
			Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int"}, {Name: "sold_on", Type: "date"}},
			Partitions: []*metadata.TablePartitionMetadata{
				{Name: "p2020", Type: metadata.TablePartitionMetadata_RANGE, Expression: "YEAR(`sold_on`)", Value: "2021"},
				{Name: "p_future", Type: metadata.TablePartitionMetadata_RANGE, Expression: "YEAR(`sold_on`)", Value: "MAXVALUE"},
			},
		},
		{
			// Sync records each partition of PARTITIONS 2.
			Name:    "hashed",
			Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int"}},
			Partitions: []*metadata.TablePartitionMetadata{
				{Name: "p0", Type: metadata.TablePartitionMetadata_HASH, Expression: "`id`"},
				{Name: "p1", Type: metadata.TablePartitionMetadata_HASH, Expression: "`id`"},
			},
		},
	}}))
	requireReport(t, report)
	requireIncomplete(t, report)
	p := requireTable(t, c, "sales").Partitioning
	if p == nil || p.Type != "RANGE" || len(p.Partitions) != 2 || p.Partitions[1].Name != "p_future" {
		t.Fatalf("sales partitioning = %+v, want two RANGE partitions", p)
	}
	if got := c.ShowCreateTable(snapshotDB, "hashed"); !strings.HasSuffix(got, "PARTITION BY HASH (`id`) PARTITIONS 2") {
		t.Errorf("hashed:\n%s\nwant PARTITION BY HASH (`id`) PARTITIONS 2", got)
	}
}

func TestLoadMetadataInstallsAutoRandomAndShardRowIDBits(t *testing.T) {
	pk := []*metadata.IndexMetadata{{Name: "PRIMARY", Type: "BTREE", Primary: true, Unique: true, Expressions: []string{"id"}}}
	c, report := loadTiDBSnapshot(t, snapshot(&metadata.SchemaMetadata{Tables: []*metadata.TableMetadata{
		{
			Name: "bits", PrimaryKeyType: "CLUSTERED", ShardingInfo: "PK_AUTO_RANDOM_BITS=5", Indexes: pk,
			Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "bigint", Default: "AUTO_RANDOM(5)"}},
		},
		{
			Name: "ranged", PrimaryKeyType: "CLUSTERED", ShardingInfo: "PK_AUTO_RANDOM_BITS=4, RANGE BITS=54", Indexes: pk,
			Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "bigint", Default: "AUTO_RANDOM(4, 54)"}},
		},
		{
			Name: "sharded", PrimaryKeyType: "NONCLUSTERED", ShardingInfo: "SHARD_BITS=4",
			Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int", Nullable: true, Default: "NULL"}},
		},
		{
			Name: "unsharded", PrimaryKeyType: "NONCLUSTERED", ShardingInfo: "NOT_SHARDED",
			Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int", Nullable: true, Default: "NULL"}},
		},
	}}))
	requireReport(t, report)
	requireIncomplete(t, report)
	for _, tc := range []struct {
		table       string
		shard, rnge int
	}{{"bits", 5, 0}, {"ranged", 4, 54}} {
		col := requireTable(t, c, tc.table).Columns[0]
		if !col.AutoRandom || col.AutoRandomShardBits != tc.shard || col.AutoRandomRangeBits != tc.rnge || col.Default != nil {
			t.Errorf("%s.id = %+v, want AUTO_RANDOM(%d, %d) and no default", tc.table, col, tc.shard, tc.rnge)
		}
	}
	if bits := requireTable(t, c, "sharded").ShardRowIDBits; bits != 4 {
		t.Errorf("sharded SHARD_ROW_ID_BITS = %d, want 4", bits)
	}
	if bits := requireTable(t, c, "unsharded").ShardRowIDBits; bits != 0 {
		t.Errorf("unsharded SHARD_ROW_ID_BITS = %d, want 0", bits)
	}
}

func TestLoadMetadataInstallsSyncedKeyExpressions(t *testing.T) {
	// TIDB_INDEXES.EXPRESSION is SQL as written, with doubled quotes and
	// escaped backslashes; sync wraps it in parentheses.
	c, report := loadTiDBSnapshot(t, snapshot(&metadata.SchemaMetadata{Tables: []*metadata.TableMetadata{{
		Name:    "t",
		Columns: []*metadata.ColumnMetadata{{Name: "j", Type: "json", Nullable: true}},
		Indexes: []*metadata.IndexMetadata{
			{Name: "ix_quote", Type: "BTREE", Expressions: []string{"(cast(json_unquote(json_extract(`j`, _utf8mb4'$.\"it''s\"')) as char(10)))"}, KeyLength: []int64{-1}, Visible: true},
			{Name: "ix_backslash", Type: "BTREE", Expressions: []string{"(cast(json_extract(`j`, _utf8mb4'$.\"a\\\\b\"') as char(10)))"}, KeyLength: []int64{-1}, Visible: true},
		},
	}}}))
	requireReport(t, report)
	requireIncomplete(t, report)
	exprs := map[string]string{}
	for _, idx := range requireTable(t, c, "t").Indexes {
		exprs[idx.Name] = idx.Columns[0].Expr
	}
	if got := exprs["ix_quote"]; !strings.Contains(got, `it\'s`) {
		t.Errorf("ix_quote = %q, want the path it's", got)
	}
	if got := exprs["ix_backslash"]; !strings.Contains(got, `a\\b`) {
		t.Errorf("ix_backslash = %q, want the path a\\b, its backslash escaped", got)
	}
}

func TestLoadMetadataIntoAnotherDatabase(t *testing.T) {
	meta := snapshot(&metadata.SchemaMetadata{
		Tables: []*metadata.TableMetadata{
			{Name: "parent", Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int"}, {Name: "k", Type: "int"}},
				Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Type: "BTREE", Primary: true, Unique: true, Expressions: []string{"id"}}}},
			{Name: "child", Columns: []*metadata.ColumnMetadata{{Name: "pid", Type: "int"}, {Name: "oid", Type: "int"}},
				ForeignKeys: []*metadata.ForeignKeyMetadata{
					{Name: "fk_same", Columns: []string{"pid"}, ReferencedSchema: snapshotDB, ReferencedTable: "parent", ReferencedColumns: []string{"id"}, OnDelete: "NO ACTION", OnUpdate: "RESTRICT"},
					{Name: "fk_other", Columns: []string{"oid"}, ReferencedSchema: "elsewhere", ReferencedTable: "parent", ReferencedColumns: []string{"id"}},
				}},
		},
		Views: []*metadata.ViewMetadata{
			// TiDB qualifies every table with the database; testdb is also an
			// alias here, and its column qualifier stays.
			{Name: "v", Definition: "SELECT `testdb`.`id` AS `id`,`testdb`.`p`.`k` AS `k` FROM `testdb`.`parent` AS `testdb` JOIN `testdb`.`parent` AS `p` ON `testdb`.`id`=`p`.`id`"},
		},
	})
	c := New()
	report, err := c.LoadMetadataInto(context.Background(), meta, "other")
	if err != nil {
		t.Fatalf("LoadMetadataInto: %v", err)
	}
	requireReport(t, report)
	if c.CurrentDatabase() != "other" || c.GetDatabase(snapshotDB) != nil {
		t.Fatalf("current database = %q, want only other", c.CurrentDatabase())
	}
	v := c.GetDatabase("other").Views["v"]
	if v == nil || v.AnalyzedQuery == nil {
		t.Fatalf("v = %+v, want an analyzed body", v)
	}
	if want := "SELECT `testdb`.`id` AS `id`,`p`.`k` AS `k` FROM `parent` AS `testdb` JOIN `parent` AS `p` ON `testdb`.`id`=`p`.`id`"; stripViewBodyDatabaseQualifier(meta.Schemas[0].Views[0].Definition, snapshotDB, "other") != want {
		t.Errorf("stripped body = %q, want %q", stripViewBodyDatabaseQualifier(meta.Schemas[0].Views[0].Definition, snapshotDB, "other"), want)
	}
	refs := map[string]*Constraint{}
	for _, con := range c.GetDatabase("other").Tables["child"].Constraints {
		if con.Type == ConForeignKey {
			refs[con.Name] = con
		}
	}
	// An unqualified reference is to the table's own database.
	if refs["fk_same"].RefDatabase != "" || refs["fk_other"].RefDatabase != "elsewhere" {
		t.Errorf("referenced databases: fk_same %q, fk_other %q; want fk_same unqualified and fk_other in elsewhere", refs["fk_same"].RefDatabase, refs["fk_other"].RefDatabase)
	}
	if got := c.ShowCreateTable("other", "child"); !strings.Contains(got, "REFERENCES `parent` (`id`) ON UPDATE RESTRICT,") {
		t.Errorf("child:\n%s\nwant fk_same with ON UPDATE RESTRICT and no NO ACTION", got)
	}
}

func TestStripViewBodyDatabaseQualifier(t *testing.T) {
	for _, tc := range []struct{ body, want string }{
		{"SELECT 1", "SELECT 1"},
		// A body that does not parse stays as it is.
		{"SELECT `testdb`.`t`.`a` FROM `testdb`.`t` WHERE (", "SELECT `testdb`.`t`.`a` FROM `testdb`.`t` WHERE ("},
		// Another database's qualifier stays.
		{"SELECT `a` FROM `other`.`t` JOIN `testdb`.`u`", "SELECT `a` FROM `other`.`t` JOIN `u`"},
		{"SELECT `testdb`.`f`(1) AS `x`", "SELECT `f`(1) AS `x`"},
		{"SELECT `a` FROM `testdb`.`t` WITH CASCADED CHECK OPTION", "SELECT `a` FROM `t` WITH CASCADED CHECK OPTION"},
		// A table the body also names as a CTE keeps a qualifier, the target's.
		{"WITH `t` AS (SELECT 1 AS `id`) SELECT `testdb`.`t`.`id` AS `id` FROM `testdb`.`t` JOIN `t` `c` JOIN `testdb`.`u`",
			"WITH `t` AS (SELECT 1 AS `id`) SELECT `other`.`t`.`id` AS `id` FROM `other`.`t` JOIN `t` `c` JOIN `u`"},
		// So does one another database's table also has.
		{"SELECT `testdb`.`t`.`id` AS `id` FROM `testdb`.`t` JOIN `ext`.`t` JOIN `testdb`.`u`",
			"SELECT `other`.`t`.`id` AS `id` FROM `other`.`t` JOIN `ext`.`t` JOIN `u`"},
	} {
		if got := stripViewBodyDatabaseQualifier(tc.body, snapshotDB, "other"); got != tc.want {
			t.Errorf("stripViewBodyDatabaseQualifier(%q) = %q, want %q", tc.body, got, tc.want)
		}
	}
}

func TestLoadMetadataReportsWhatATableLeftOut(t *testing.T) {
	c, report := loadTiDBSnapshot(t, snapshot(&metadata.SchemaMetadata{Tables: []*metadata.TableMetadata{
		{Name: "bad_default", Columns: []*metadata.ColumnMetadata{{Name: "a", Type: "int", Default: "((("}}},
		{Name: "bad_on_update", Columns: []*metadata.ColumnMetadata{{Name: "a", Type: "timestamp", OnUpdate: "((("}}},
		{Name: "bad_generation", Columns: []*metadata.ColumnMetadata{{Name: "a", Type: "int", Generation: &metadata.GenerationMetadata{Expression: "a +"}}}},
		{Name: "bad_auto_random", Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "bigint", Default: "AUTO_RANDOM(x)"}}},
		{Name: "bad_check", Columns: []*metadata.ColumnMetadata{{Name: "a", Type: "int"}},
			CheckConstraints: []*metadata.CheckConstraintMetadata{{Name: "chk", Expression: "a >>> 0"}}},
		{Name: "bad_index", Columns: []*metadata.ColumnMetadata{{Name: "a", Type: "int"}},
			Indexes: []*metadata.IndexMetadata{{Name: "idx", Type: "BTREE", Expressions: []string{"(lower(((a)"}}}},
		{Name: "bad_partitions", Columns: []*metadata.ColumnMetadata{{Name: "a", Type: "int"}},
			Partitions: []*metadata.TablePartitionMetadata{{Name: "p0", Type: metadata.TablePartitionMetadata_RANGE, Expression: "a +", Value: "1"}}},
	}}))
	requireReport(t, report)
	requireIncomplete(t, report, "table:bad_default", "table:bad_on_update", "table:bad_generation", "table:bad_auto_random",
		"table:bad_check", "table:bad_index", "table:bad_partitions")
	// The rest of the table installs.
	if tbl := requireTable(t, c, "bad_index"); len(tbl.Columns) != 1 || len(tbl.Indexes) != 0 {
		t.Errorf("bad_index = %+v, want column a without idx", tbl)
	}
}

func TestParseExpr(t *testing.T) {
	for _, expr := range []string{`'hello'`, `CURRENT_TIMESTAMP`, `NULL`, `a + b`, `CONCAT('x', 'y')`} {
		if _, err := parseExpr(expr); err != nil {
			t.Errorf("parseExpr(%q): %v", expr, err)
		}
	}
	for _, expr := range []string{"", "(((", "this is not sql at all ###"} {
		if _, err := parseExpr(expr); err == nil {
			t.Errorf("parseExpr(%q) succeeded, want an error", expr)
		}
	}
}

// TestLoadMetadataKeepsCheckConstraints: a snapshot's CHECK constraints exist on
// the server it came from, so loading keeps them while tidb_enable_check_constraint
// is off for the session, and leaves the session's setting as it was.
func TestLoadMetadataKeepsCheckConstraints(t *testing.T) {
	c, report := loadTiDBSnapshot(t, snapshot(&metadata.SchemaMetadata{
		Tables: []*metadata.TableMetadata{{
			Name:             "t",
			Columns:          []*metadata.ColumnMetadata{{Name: "a", Type: "int", Nullable: true}},
			CheckConstraints: []*metadata.CheckConstraintMetadata{{Name: "ck_a", Expression: "(`a` > 0)"}},
		}},
	}))
	requireReport(t, report)
	var found bool
	for _, con := range requireTable(t, c, "t").Constraints {
		if con.Type == ConCheck && con.Name == "ck_a" {
			found = true
		}
	}
	if !found {
		t.Errorf("the snapshot's CHECK ck_a was dropped")
	}
	if c.CheckConstraints() {
		t.Errorf("loading left tidb_enable_check_constraint on")
	}
}

// metadataLoadOracleCases are TiDB scripts whose tables and views
// TestContainerLoadMetadataMatchesTiDB syncs back and loads. The database they
// run in is metadataLoadSourceDB, and v_alias reads a table aliased like it.
var metadataLoadOracleCases = []struct{ name, script string }{
	{"partitions", `
CREATE TABLE p_range (id int, d date) PARTITION BY RANGE (YEAR(d)) (PARTITION z1 VALUES LESS THAN (2000), PARTITION a2 VALUES LESS THAN (2010), PARTITION pmax VALUES LESS THAN MAXVALUE);
CREATE TABLE p_range_columns (a int, b varchar(10)) PARTITION BY RANGE COLUMNS (a, b) (PARTITION p0 VALUES LESS THAN (10, 'm'), PARTITION p1 VALUES LESS THAN (MAXVALUE, MAXVALUE));
CREATE TABLE p_list (id int, r int) PARTITION BY LIST (r) (PARTITION pa VALUES IN (1, 2), PARTITION pb VALUES IN (3, NULL), PARTITION pd VALUES IN (DEFAULT));
CREATE TABLE p_list_columns (id int, c varchar(10)) PARTITION BY LIST COLUMNS (c) (PARTITION pa VALUES IN ('a', 'b'), PARTITION pb VALUES IN ('c'));
CREATE TABLE p_hash (id int) PARTITION BY HASH (id) PARTITIONS 4;
CREATE TABLE p_hash_named (id int) PARTITION BY HASH (id) (PARTITION x0, PARTITION x1);
CREATE TABLE p_key (id int PRIMARY KEY) PARTITION BY KEY (id) PARTITIONS 3;
`},
	{"auto random and row id sharding", `
CREATE TABLE ar_bits (id bigint PRIMARY KEY AUTO_RANDOM(5), v int);
CREATE TABLE ar_range (id bigint AUTO_RANDOM(4, 54), v int, PRIMARY KEY (id) CLUSTERED);
CREATE TABLE ar_default (id bigint AUTO_RANDOM, PRIMARY KEY (id));
CREATE TABLE sharded (id int, v int) SHARD_ROW_ID_BITS=4;
CREATE TABLE sharded_pk (id int, v int, PRIMARY KEY (id) NONCLUSTERED) SHARD_ROW_ID_BITS=3;
CREATE TABLE nonclustered (a int, b varchar(10), PRIMARY KEY (a, b) NONCLUSTERED);
`},
	{"expression indexes", `
CREATE TABLE ex (id int PRIMARY KEY, j json, name varchar(20),
  KEY ix_quote ((cast(json_unquote(json_extract(j, '$."it''s"')) as char(10)))),
  KEY ix_backslash ((cast(json_extract(j, '$."a\\b"') as char(10)))),
  KEY ix_mixed ((lower(name)), id),
  KEY ix_prefix (name(4), id),
  KEY ix_array ((cast(j->'$.tags' as unsigned array))));
`},
	{"views and foreign keys", `
CREATE TABLE parent (id int PRIMARY KEY, k int, a int, b int, UNIQUE KEY uk (k), UNIQUE KEY uab (a, b));
CREATE TABLE child (id int PRIMARY KEY, pid int, k int, a int, b int,
  CONSTRAINT fk_cascade FOREIGN KEY (pid) REFERENCES parent (id) ON DELETE CASCADE,
  CONSTRAINT fk_qualified FOREIGN KEY (k) REFERENCES omni_mdload.parent (k) ON UPDATE SET NULL,
  CONSTRAINT fk_pair FOREIGN KEY (a, b) REFERENCES parent (a, b));
CREATE VIEW v_join AS SELECT c.id, p.k AS pk, omni_mdload.c.pid FROM omni_mdload.child c JOIN parent p ON c.pid = p.id;
CREATE VIEW v_alias AS SELECT omni_mdload.id FROM parent omni_mdload WHERE omni_mdload.k > 1;
CREATE VIEW v_nested AS SELECT id, count(*) AS n FROM v_join GROUP BY id;
`},
	{"checks, defaults and generated columns", `
CREATE TABLE ck (id int PRIMARY KEY, c varchar(10), g int, CONSTRAINT ck_named CHECK (g > 0 AND c <> 'it''s'), CHECK (g < 100));
CREATE TABLE defaults (id int PRIMARY KEY,
  a timestamp(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  b datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  c varchar(10) DEFAULT 'it''s',
  e varchar(36) DEFAULT (uuid()),
  f date DEFAULT (curdate()),
  g int DEFAULT '5',
  h enum('x','y') NOT NULL DEFAULT 'y',
  i decimal(10,2) DEFAULT '1.50' COMMENT 'price');
CREATE TABLE gen_cols (a int, b int GENERATED ALWAYS AS (a * 2) VIRTUAL, c varchar(10) GENERATED ALWAYS AS (concat('it''s', a)) STORED);
`},
}

const (
	metadataLoadSourceDB = "omni_mdload"
	metadataLoadTargetDB = "omni_loaded"
)

// metadataLoadGaps lists the tables whose SHOW CREATE TABLE the catalog writes
// differently from TiDB however they were defined, and why. A listed table must
// still load as the catalog defines it from TiDB's own SHOW CREATE TABLE, and
// one that matches TiDB fails, so a fix removes its entry.
var metadataLoadGaps = map[string]string{
	"table:ex":       "key expressions deparse as MySQL writes them: no space after a comma, CAST AS CHAR with its charset, 'it\\'s'",
	"table:ck":       "CHECK expressions deparse as MySQL writes them: each comparison parenthesized, <> for !=, 'it\\'s'",
	"table:defaults": "a string default with a quote is written 'it\\'s', not 'it''s'",
	"table:gen_cols": "a generation expression writes a quote in a string unescaped: 'it's'",
}

// TestContainerLoadMetadataMatchesTiDB runs each script on TiDB, builds the
// snapshot Bytebase's TiDB sync would (syncTiDBSnapshot), loads it under
// another database name, and compares the catalog's SHOW CREATE TABLE and SHOW
// CREATE VIEW with TiDB's, modulo the database name and the order of keys and
// constraints, which sync records by name.
func TestContainerLoadMetadataMatchesTiDB(t *testing.T) {
	ctr, done := startContainer(t)
	defer done()
	exec := func(t *testing.T, stmt string) {
		t.Helper()
		if _, err := ctr.db.ExecContext(ctr.ctx, stmt); err != nil {
			t.Fatalf("TiDB: %q: %v", stmt, err)
		}
	}
	// TiDB ignores CHECK constraints unless the global switch is on.
	exec(t, "SET GLOBAL tidb_enable_check_constraint = ON")
	defer exec(t, "SET GLOBAL tidb_enable_check_constraint = OFF")
	for _, tc := range metadataLoadOracleCases {
		t.Run(tc.name, func(t *testing.T) {
			exec(t, "DROP DATABASE IF EXISTS "+metadataLoadSourceDB)
			exec(t, "CREATE DATABASE "+metadataLoadSourceDB)
			exec(t, "USE "+metadataLoadSourceDB)
			defer exec(t, "USE test")
			for _, stmt := range splitStatements(tc.script) {
				if stmt = strings.TrimSpace(stmt); stmt != "" {
					exec(t, stmt)
				}
			}
			meta := syncTiDBSnapshot(t, ctr, metadataLoadSourceDB)
			c := New()
			report, err := c.LoadMetadataInto(context.Background(), meta, metadataLoadTargetDB)
			if err != nil {
				t.Fatal(err)
			}
			if len(report.Degraded)+len(report.Missing)+len(report.Incomplete) > 0 {
				t.Errorf("degraded %v, missing %v, incomplete %v", report.Degraded, report.Missing, report.Incomplete)
			}
			for _, obj := range syncedObjects(meta) {
				stored := showCreateOnTiDB(t, ctr, obj)
				want := sortKeyLines(columnDatabaseQualifier.ReplaceAllString(stored, "$1"))
				got := showCreateInCatalog(c, obj)
				gap, listed := metadataLoadGaps[obj]
				switch {
				case got != want && !listed:
					t.Errorf("%s: catalog\n%s\nTiDB\n%s", obj, got, want)
				case got == want && listed:
					t.Errorf("%s: listed as a gap (%s) but matches: remove the entry", obj, gap)
				case listed:
					defined := New()
					defined.SetCheckConstraints(true)
					defined.SetForeignKeyChecks(false)
					results, err := defined.Exec("CREATE DATABASE "+metadataLoadTargetDB+"; USE "+metadataLoadTargetDB+";\n"+stored, nil)
					if err != nil {
						t.Fatal(err)
					}
					for _, r := range results {
						if r.Error != nil {
							t.Fatalf("%s: the catalog refuses TiDB's SHOW CREATE TABLE: %v", obj, r.Error)
						}
					}
					if written := showCreateInCatalog(defined, obj); got != written {
						t.Errorf("%s: loaded\n%s\ndefined from TiDB's SHOW CREATE TABLE\n%s", obj, got, written)
					}
				}
			}
		})
	}
}

// showCreateInCatalog writes obj as SHOW CREATE does in metadataLoadTargetDB,
// with that database's qualifier spelled as metadataLoadSourceDB.
func showCreateInCatalog(c *Catalog, obj string) string {
	kind, name, _ := strings.Cut(obj, ":")
	var text string
	if kind == "view" {
		text = c.ShowCreateView(metadataLoadTargetDB, name)
	} else {
		text = c.ShowCreateTable(metadataLoadTargetDB, name)
	}
	return sortKeyLines(strings.ReplaceAll(text, quoteIdent(metadataLoadTargetDB)+".", quoteIdent(metadataLoadSourceDB)+"."))
}

// sortKeyLines sorts the key and constraint lines of a CREATE TABLE.
func sortKeyLines(text string) string {
	lines := strings.Split(text, "\n")
	end := slices.IndexFunc(lines, func(l string) bool { return strings.HasPrefix(l, ")") })
	if end < 0 {
		return text
	}
	first := slices.IndexFunc(lines[:end], func(l string) bool { return strings.HasPrefix(l, "  ") && !strings.HasPrefix(l, "  `") })
	if first < 0 {
		return text
	}
	keys := lines[first:end]
	for i := range keys {
		keys[i] = strings.TrimSuffix(keys[i], ",")
	}
	slices.Sort(keys)
	for i := range keys[:len(keys)-1] {
		keys[i] += ","
	}
	return strings.Join(lines, "\n")
}

// columnDatabaseQualifier matches a column reference qualified with
// metadataLoadSourceDB, which the loader drops.
var columnDatabaseQualifier = regexp.MustCompile("`" + metadataLoadSourceDB + "`\\.(`[^`]+`\\.`)")

func syncedObjects(meta *metadata.DatabaseSchemaMetadata) []string {
	var out []string
	for _, s := range meta.GetSchemas() {
		for _, tbl := range s.GetTables() {
			out = append(out, "table:"+tbl.GetName())
		}
		for _, v := range s.GetViews() {
			out = append(out, "view:"+v.GetName())
		}
	}
	return out
}

func showCreateOnTiDB(t *testing.T, ctr *mysqlContainer, obj string) string {
	t.Helper()
	kind, name, _ := strings.Cut(obj, ":")
	query := "SHOW CREATE TABLE " + quoteIdent(name)
	if kind == "view" {
		query = "SHOW CREATE VIEW " + quoteIdent(name)
	}
	rows, err := ctr.db.QueryContext(ctr.ctx, query)
	if err != nil {
		t.Fatalf("%s: %v", query, err)
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		t.Fatal(err)
	}
	if !rows.Next() {
		t.Fatalf("%s: no row", query)
	}
	vals := make([]sql.NullString, len(cols))
	ptrs := make([]any, len(cols))
	for i := range vals {
		ptrs[i] = &vals[i]
	}
	if err := rows.Scan(ptrs...); err != nil {
		t.Fatal(err)
	}
	return vals[1].String
}

// syncTiDBSnapshot reads database name back as Bytebase's TiDB sync does
// (backend/plugin/db/tidb/sync.go), text transformations included. It departs
// from that sync where the sync is wrong about TiDB v8.5.5, so the snapshot
// describes the database: it orders index key parts by SEQ_IN_INDEX and
// partitions by PARTITION_ORDINAL_POSITION (sync orders neither, and orders
// partitions by name), reads CHECK constraints from TIDB_CHECK_CONSTRAINTS
// (TABLE_CONSTRAINTS, which sync joins, lists none), and records a foreign
// key's REFERENCED_TABLE_SCHEMA (sync leaves it empty).
func syncTiDBSnapshot(t *testing.T, ctr *mysqlContainer, name string) *metadata.DatabaseSchemaMetadata {
	t.Helper()
	query := func(q string, scan func(rows *sql.Rows) error) {
		t.Helper()
		rows, err := ctr.db.QueryContext(ctr.ctx, q, name)
		if err != nil {
			t.Fatalf("%s: %v", q, err)
		}
		defer rows.Close()
		for rows.Next() {
			if err := scan(rows); err != nil {
				t.Fatal(err)
			}
		}
		if err := rows.Err(); err != nil {
			t.Fatal(err)
		}
	}

	indexes := map[string][]*metadata.IndexMetadata{}
	query(`SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX`, func(rows *sql.Rows) error {
		var table, key, comment string
		var column, expression sql.NullString
		var nonUnique bool
		var subPart sql.NullInt64
		if err := rows.Scan(&table, &key, &column, &nonUnique, &comment, &subPart, &expression); err != nil {
			return err
		}
		var idx *metadata.IndexMetadata
		if list := indexes[table]; len(list) > 0 && list[len(list)-1].Name == key {
			idx = list[len(list)-1]
		} else {
			idx = &metadata.IndexMetadata{Name: key, Type: "BTREE", Unique: !nonUnique, Primary: key == "PRIMARY", Visible: true, Comment: comment}
			indexes[table] = append(indexes[table], idx)
		}
		switch {
		case expression.Valid:
			idx.Expressions = append(idx.Expressions, "("+expression.String+")")
		case column.Valid:
			idx.Expressions = append(idx.Expressions, column.String)
		}
		length := int64(-1)
		if subPart.Valid {
			length = subPart.Int64
		}
		idx.KeyLength = append(idx.KeyLength, length)
		return nil
	})

	columns := map[string][]*metadata.ColumnMetadata{}
	query(`SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION`, func(rows *sql.Rows) error {
		col := &metadata.ColumnMetadata{}
		var table, nullable, extra string
		var def, generation sql.NullString
		if err := rows.Scan(&table, &col.Name, &col.Position, &def, &nullable, &col.Type, &col.CharacterSet, &col.Collation, &col.Comment, &generation, &extra); err != nil {
			return err
		}
		col.Comment = strings.TrimSuffix(strings.TrimPrefix(col.Comment, "'"), "'")
		col.Nullable = nullable == "YES"
		syncColumnDefault(col, def, extra)
		if generation.String != "" {
			switch {
			case strings.Contains(strings.ToUpper(extra), "VIRTUAL GENERATED"):
				col.Generation = &metadata.GenerationMetadata{Type: metadata.GenerationMetadata_TYPE_VIRTUAL, Expression: generation.String}
			case strings.Contains(strings.ToUpper(extra), "STORED GENERATED"):
				col.Generation = &metadata.GenerationMetadata{Type: metadata.GenerationMetadata_TYPE_STORED, Expression: generation.String}
			}
		}
		columns[table] = append(columns[table], col)
		return nil
	})

	var views []*metadata.ViewMetadata
	query(`SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME`, func(rows *sql.Rows) error {
		v := &metadata.ViewMetadata{}
		if err := rows.Scan(&v.Name, &v.Definition); err != nil {
			return err
		}
		v.Columns = columns[v.Name]
		views = append(views, v)
		return nil
	})

	foreignKeys := map[string][]*metadata.ForeignKeyMetadata{}
	query(`SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION`, func(rows *sql.Rows) error {
		var table, column, refColumn string
		fk := &metadata.ForeignKeyMetadata{}
		if err := rows.Scan(&table, &fk.Name, &column, &fk.ReferencedSchema, &fk.ReferencedTable, &refColumn, &fk.OnDelete, &fk.OnUpdate, &fk.MatchType); err != nil {
			return err
		}
		if list := foreignKeys[table]; len(list) > 0 && list[len(list)-1].Name == fk.Name {
			fk = list[len(list)-1]
		} else {
			foreignKeys[table] = append(foreignKeys[table], fk)
		}
		fk.Columns = append(fk.Columns, column)
		fk.ReferencedColumns = append(fk.ReferencedColumns, refColumn)
		return nil
	})

	checks := map[string][]*metadata.CheckConstraintMetadata{}
	query(`SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME`, func(rows *sql.Rows) error {
		var table string
		check := &metadata.CheckConstraintMetadata{}
		if err := rows.Scan(&table, &check.Name, &check.Expression); err != nil {
			return err
		}
		checks[table] = append(checks[table], check)
		return nil
	})

	partitions := map[string][]*metadata.TablePartitionMetadata{}
	query(`SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION`, func(rows *sql.Rows) error {
		var table, partition, method string
		var expression, description sql.NullString
		if err := rows.Scan(&table, &partition, &method, &expression, &description); err != nil {
			return err
		}
		typ, ok := metadata.TablePartitionMetadata_Type_value[strings.ReplaceAll(strings.ToUpper(method), " ", "_")]
		if !ok {
			return fmt.Errorf("partition method %q", method)
		}
		partitions[table] = append(partitions[table], &metadata.TablePartitionMetadata{
			Name: partition, Type: metadata.TablePartitionMetadata_Type(typ), Expression: expression.String, Value: description.String,
		})
		return nil
	})

	schema := &metadata.SchemaMetadata{Views: views}
	query(`SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, '')
		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME`, func(rows *sql.Rows) error {
		var table, tableType, engine, collation, comment, shardingInfo, pkType string
		if err := rows.Scan(&table, &tableType, &engine, &collation, &comment, &shardingInfo, &pkType); err != nil {
			return err
		}
		if tableType != "BASE TABLE" {
			return nil
		}
		cols := columns[table]
		// An AUTO_RANDOM primary key column carries the attribute as its
		// default, rebuilt from PK_AUTO_RANDOM_BITS=n[, RANGE BITS=m].
		if bits, ok := strings.CutPrefix(shardingInfo, "PK_AUTO_RANDOM_BITS="); ok {
			text := "AUTO_RANDOM"
			shard, rangeBits, hasRange := strings.Cut(bits, ", RANGE BITS=")
			if hasRange {
				text += "(" + shard + ", " + rangeBits + ")"
			} else {
				text += "(" + shard + ")"
			}
			for _, idx := range indexes[table] {
				if idx.Primary && len(idx.Expressions) > 0 {
					for _, col := range cols {
						if col.Name == idx.Expressions[0] {
							col.Default = text
						}
					}
				}
			}
		}
		schema.Tables = append(schema.Tables, &metadata.TableMetadata{
			Name:             table,
			Columns:          cols,
			Indexes:          indexes[table],
			ForeignKeys:      foreignKeys[table],
			CheckConstraints: checks[table],
			Partitions:       partitions[table],
			Engine:           engine,
			Collation:        collation,
			Charset:          strings.Split(collation, "_")[0],
			Comment:          strings.TrimSuffix(strings.TrimPrefix(comment, "'"), "'"),
			ShardingInfo:     shardingInfo,
			PrimaryKeyType:   pkType,
		})
		return nil
	})

	meta := &metadata.DatabaseSchemaMetadata{Name: name, Schemas: []*metadata.SchemaMetadata{schema}}
	query(`SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?`, func(rows *sql.Rows) error {
		return rows.Scan(&meta.CharacterSet, &meta.Collation)
	})
	return meta
}

// syncColumnDefault is Bytebase's setColumnMetadataDefault for TiDB: def is
// QUOTE(COLUMN_DEFAULT).
func syncColumnDefault(col *metadata.ColumnMetadata, def sql.NullString, extra string) {
	switch {
	case def.Valid:
		unquoted := unquoteSyncedString(def.String)
		upper := strings.ToUpper(unquoted)
		switch {
		case strings.HasPrefix(upper, "CURRENT_TIMESTAMP"), strings.HasPrefix(upper, "CURRENT_DATE"):
			col.Default = unquoted
		case strings.Contains(extra, "DEFAULT_GENERATED"):
			if _, err := time.Parse(time.DateTime, def.String); err == nil || def.String == "0000-00-00 00:00:00" {
				col.Default = def.String
			} else {
				col.Default = "(" + strings.ReplaceAll(strings.ReplaceAll(unquoted, `\'`, `'`), `\\`, `\`) + ")"
			}
		default:
			col.Default = def.String
		}
	case strings.Contains(strings.ToUpper(extra), "AUTO_INCREMENT"):
		col.Default = autoIncrementDefault
	case col.Nullable:
		col.Default = "NULL"
	}
	if strings.Contains(extra, "on update CURRENT_TIMESTAMP") {
		col.OnUpdate = "CURRENT_TIMESTAMP"
		if i := strings.Index(extra, "CURRENT_TIMESTAMP("); i >= 0 {
			col.OnUpdate = extra[i:]
			if j := strings.Index(col.OnUpdate, ")"); j >= 0 {
				col.OnUpdate = col.OnUpdate[:j+1]
			}
		}
	}
}

// unquoteSyncedString is Bytebase's UnquoteMySQLString.
func unquoteSyncedString(s string) string {
	if len(s) >= 2 && s[0] == '\'' && s[len(s)-1] == '\'' {
		s = s[1 : len(s)-1]
	}
	escapes := map[byte]byte{'\\': '\\', '\'': '\'', '0': 0, 'Z': 0x1A, 'n': '\n', 't': '\t', 'r': '\r', 'b': '\b'}
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) {
			if r, ok := escapes[s[i+1]]; ok {
				b.WriteByte(r)
				i++
				continue
			}
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

func TestLoadMetadataKeepsTableWhosePartitioningIsRefused(t *testing.T) {
	// Sync lists partitions by name, which can put MAXVALUE before another
	// partition; the catalog refuses that order.
	c, report := loadTiDBSnapshot(t, snapshot(&metadata.SchemaMetadata{Tables: []*metadata.TableMetadata{{
		Name:    "sales",
		Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int"}, {Name: "sold_on", Type: "date"}},
		Indexes: []*metadata.IndexMetadata{{Name: "k", Expressions: []string{"id"}, Visible: true}},
		Partitions: []*metadata.TablePartitionMetadata{
			{Name: "a2", Type: metadata.TablePartitionMetadata_RANGE, Expression: "YEAR(`sold_on`)", Value: "2022"},
			{Name: "pmax", Type: metadata.TablePartitionMetadata_RANGE, Expression: "YEAR(`sold_on`)", Value: "MAXVALUE"},
			{Name: "z1", Type: metadata.TablePartitionMetadata_RANGE, Expression: "YEAR(`sold_on`)", Value: "2021"},
		},
	}}}))
	requireReport(t, report)
	requireIncomplete(t, report, "table:sales")
	tbl := requireTable(t, c, "sales")
	if tbl.Partitioning != nil || len(tbl.Indexes) != 1 || tbl.Columns[1].DataType != "date" {
		t.Errorf("sales = partitioning %+v, %d indexes, sold_on %s; want no partitioning, the index and a date column", tbl.Partitioning, len(tbl.Indexes), tbl.Columns[1].DataType)
	}
}
