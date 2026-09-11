package catalog

import (
	"context"
	"slices"
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
	report, err := c.LoadMetadata(context.Background(), snapshot(&metadata.SchemaMetadata{Tables: []*metadata.TableMetadata{
		{
			// Sorts first and references z_parent.
			Name:    "a_child",
			Charset: "utf8mb4",
			Comment: "children",
			Columns: []*metadata.ColumnMetadata{
				{Name: "id", Type: "bigint unsigned", Default: autoIncrementDefault},
				{Name: "parent_id", Type: "bigint unsigned"},
				{Name: "email", Type: "varchar(255)", Nullable: true, CharacterSet: "latin1"},
				{Name: "created_at", Type: "timestamp", Default: "CURRENT_TIMESTAMP", OnUpdate: "CURRENT_TIMESTAMP"},
				{Name: "total", Type: "int", Nullable: true, Generation: &metadata.GenerationMetadata{
					Type: metadata.GenerationMetadata_TYPE_STORED, Expression: "parent_id + 1",
				}},
				// Sync reports DEFAULT NULL for every nullable column, including
				// types that take no DEFAULT clause.
				{Name: "payload", Type: "json", Nullable: true, Default: "NULL"},
				{Name: "name", Type: "varchar(255)", Nullable: true, Default: "NULL"},
			},
			Indexes: []*metadata.IndexMetadata{
				{Name: "PRIMARY", Type: "BTREE", Primary: true, Unique: true, Expressions: []string{"id"}},
				{Name: "idx_parent_id", Type: "BTREE", Expressions: []string{"parent_id"}},
				{Name: "idx_lower_name", Type: "BTREE", Expressions: []string{"(lower(`name`))"}},
				{Name: "idx_prefix", Type: "BTREE", Expressions: []string{"name"}, KeyLength: []int64{10}, Descending: []bool{true}},
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
			Name:    "z_parent",
			Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "bigint unsigned", Default: autoIncrementDefault}},
			Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Type: "BTREE", Primary: true, Unique: true, Expressions: []string{"id"}}},
		},
	}}))
	if err != nil {
		t.Fatalf("LoadMetadata: %v", err)
	}
	requireReport(t, report)
	if c.CurrentDatabase() != snapshotDB || !c.ForeignKeyChecks() {
		t.Errorf("current database %q, foreign key checks %v; want %s and checks back on", c.CurrentDatabase(), c.ForeignKeyChecks(), snapshotDB)
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
	if payload := child.GetColumn("payload"); payload.Default != nil {
		t.Errorf("payload default = %q, want none", *payload.Default)
	}
	if name := child.GetColumn("name"); name.Default == nil {
		t.Error("name lost its DEFAULT NULL")
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
	var fk *Constraint
	for _, con := range child.Constraints {
		if con.Type == ConForeignKey && con.Name == "fk_child_parent" {
			fk = con
		}
	}
	if fk == nil || !strings.EqualFold(fk.OnDelete, "CASCADE") {
		t.Errorf("fk_child_parent = %+v, want ON DELETE CASCADE", fk)
	}
}

func TestLoadMetadataInstallsViewsInAnyOrder(t *testing.T) {
	// a_view reads z_view, which installs after it.
	c, report := loadTiDBSnapshot(t, snapshot(&metadata.SchemaMetadata{
		Tables: []*metadata.TableMetadata{{
			Name:    "users",
			Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int"}, {Name: "name", Type: "varchar(255)", Nullable: true}},
		}},
		Views: []*metadata.ViewMetadata{
			{Name: "a_view", Definition: "SELECT id FROM z_view"},
			{Name: "z_view", Definition: "SELECT id, name FROM users"},
		},
	}))
	requireReport(t, report)
	// The columns come from analyzing the parsed body.
	if v := requireView(t, c, "z_view"); !slices.Equal(v.Columns, []string{"id", "name"}) {
		t.Errorf("z_view columns = %v, want id, name", v.Columns)
	}
	requireView(t, c, "a_view")
}

func TestLoadMetadataStandsInForFailedTables(t *testing.T) {
	c, report := loadTiDBSnapshot(t, snapshot(&metadata.SchemaMetadata{Tables: []*metadata.TableMetadata{{
		Name:    "broken",
		Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int"}, {Name: "bad", Type: "varchar((("}},
	}}}))
	requireReport(t, report, "table:broken")
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
