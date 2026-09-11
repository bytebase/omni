package catalog

import (
	"context"
	"slices"
	"strconv"
	"strings"
	"testing"

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
			// A body that does not resolve names nothing of its own.
			{Name: "unresolved", Definition: "SELECT * FROM missing", Columns: []*metadata.ColumnMetadata{{Name: "a"}, {Name: "b"}}},
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
	for name, want := range map[string][]string{"same": {"id"}, "star": {"id"}, "star_qualified": {"id"}, "expr": {"COUNT(*)"}, "unresolved": {"a", "b"}} {
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
