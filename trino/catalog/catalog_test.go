package catalog_test

import (
	"slices"
	"testing"

	"github.com/bytebase/omni/trino/ast"
	"github.com/bytebase/omni/trino/catalog"
)

// --- normalization boundary --------------------------------------------------

// TestNormalize pins the load-bearing Trino semantic mirrored by this package:
// unquoted identifiers fold to lower case; double-quoted identifiers keep their
// case (and lose the quotes). It must agree with ast.Identifier.Normalize,
// which is what the parser feeds the catalog.
func TestNormalize(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"NATION", "nation"},
		{"nation", "nation"},
		{"MyTable", "mytable"},
		{`"MyTable"`, "MyTable"}, // quoted preserves case
		{`"lower"`, "lower"},     // quoted, already lower
		{`""`, ""},               // empty quoted identifier
		{`"a""b"`, `a""b`},       // strips only outer quotes; embedded "" not collapsed (route quoted source tokens through the AST)
		{"1Col", "1col"},         // digit-leading unquoted folds
	}
	for _, tt := range tests {
		if got := catalog.Normalize(tt.in); got != tt.want {
			t.Errorf("Normalize(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

// TestNormalizeMatchesASTIdentifier proves the package Normalize and the
// EnsureCatalogIdent path agree with ast.Identifier.Normalize for both quoted
// and unquoted identifiers — the interop contract with the parser.
func TestNormalizeMatchesASTIdentifier(t *testing.T) {
	cases := []struct {
		raw string
		id  *ast.Identifier
	}{
		{"MyTable", &ast.Identifier{Value: "MyTable"}},
		{`"MyTable"`, &ast.Identifier{Value: "MyTable", Quoted: true, QuoteRune: '"'}},
		{"orders", &ast.Identifier{Value: "orders"}},
	}
	for _, tc := range cases {
		if catalog.Normalize(tc.raw) != tc.id.Normalize() {
			t.Errorf("Normalize(%q)=%q disagrees with ast.Identifier.Normalize()=%q",
				tc.raw, catalog.Normalize(tc.raw), tc.id.Normalize())
		}
	}
}

// --- top-level Catalog -------------------------------------------------------

func TestNewCatalogEmpty(t *testing.T) {
	c := catalog.New()
	if got := c.Catalogs(); len(got) != 0 {
		t.Errorf("new catalog Catalogs() = %v, want empty", got)
	}
	if c.HasCatalog("tpch") {
		t.Error("HasCatalog returned true on empty catalog")
	}
	if c.GetCatalog("tpch") != nil {
		t.Error("GetCatalog returned non-nil on empty catalog")
	}
	if c.CurrentCatalog() != "" {
		t.Errorf("CurrentCatalog() = %q, want empty", c.CurrentCatalog())
	}
	if c.CurrentSchema() != "" {
		t.Errorf("CurrentSchema() = %q, want empty", c.CurrentSchema())
	}
}

func TestEnsureCatalogIdempotent(t *testing.T) {
	c := catalog.New()
	a := c.EnsureCatalog("tpch")
	b := c.EnsureCatalog("tpch")
	if a != b {
		t.Error("EnsureCatalog returned a different instance on second call")
	}
	if got := c.Catalogs(); !slices.Equal(got, []string{"tpch"}) {
		t.Errorf("Catalogs() = %v, want [tpch]", got)
	}
}

func TestCatalogsSortedAndDeduped(t *testing.T) {
	c := catalog.New()
	c.EnsureCatalog("tpch")
	c.EnsureCatalog("memory")
	c.EnsureCatalog("system")
	c.EnsureCatalog("memory") // dup
	got := c.Catalogs()
	want := []string{"memory", "system", "tpch"}
	if !slices.Equal(got, want) {
		t.Errorf("Catalogs() = %v, want %v", got, want)
	}
}

// TestEnsureCatalogIdent exercises the quote-aware AST boundary: an unquoted
// identifier folds, a quoted one preserves case, and the two are distinct
// catalogs.
func TestEnsureCatalogIdent(t *testing.T) {
	c := catalog.New()
	c.EnsureCatalogIdent(&ast.Identifier{Value: "TPCH"})                               // -> "tpch"
	c.EnsureCatalogIdent(&ast.Identifier{Value: "Hive", Quoted: true, QuoteRune: '"'}) // -> "Hive"
	if !c.HasCatalog("tpch") {
		t.Error("unquoted TPCH did not fold to tpch")
	}
	if !c.HasCatalog("Hive") {
		t.Error("quoted Hive did not preserve case")
	}
	if c.HasCatalog("hive") {
		t.Error("quoted Hive must not match lower-cased hive")
	}
	if got := c.Catalogs(); !slices.Equal(got, []string{"Hive", "tpch"}) {
		t.Errorf("Catalogs() = %v, want [Hive tpch]", got)
	}
}

// --- Database (catalog) -> Schema -------------------------------------------

func TestSchemasSorted(t *testing.T) {
	c := catalog.New()
	db := c.EnsureCatalog("tpch")
	db.EnsureSchema("sf100")
	db.EnsureSchema("sf1")
	db.EnsureSchema("tiny")
	if got := db.Schemas(); !slices.Equal(got, []string{"sf1", "sf100", "tiny"}) {
		t.Errorf("Schemas() = %v, want [sf1 sf100 tiny]", got)
	}
	if !db.HasSchema("sf1") {
		t.Error("HasSchema(sf1) false")
	}
	if db.GetSchema("missing") != nil {
		t.Error("GetSchema(missing) non-nil")
	}
}

func TestEnsureSchemaIdempotent(t *testing.T) {
	db := catalog.New().EnsureCatalog("tpch")
	a := db.EnsureSchema("sf1")
	b := db.EnsureSchema("sf1")
	if a != b {
		t.Error("EnsureSchema returned a different instance on second call")
	}
}

// --- Schema -> Table / View / Column ----------------------------------------

func TestTablesAndColumns(t *testing.T) {
	sch := catalog.New().EnsureCatalog("tpch").EnsureSchema("sf1")
	sch.AddTable("nation",
		catalog.NewColumn("nationkey", "bigint", false),
		catalog.NewColumn("name", "varchar(25)", false),
		catalog.NewColumn("comment", "varchar(152)", true),
	)

	if got := sch.Tables(); !slices.Equal(got, []string{"nation"}) {
		t.Errorf("Tables() = %v, want [nation]", got)
	}
	if !sch.HasTable("nation") {
		t.Error("HasTable(nation) false")
	}

	tbl := sch.GetTable("nation")
	if tbl == nil {
		t.Fatal("GetTable(nation) nil")
	}
	// Column order is ordinal (matters for SELECT *).
	if got := tbl.ColumnNames(); !slices.Equal(got, []string{"nationkey", "name", "comment"}) {
		t.Errorf("ColumnNames() = %v, want [nationkey name comment]", got)
	}
	// Columns() returns the ordinal-ordered column slice.
	if cols := tbl.Columns(); len(cols) != 3 || cols[0].Name != "nationkey" {
		t.Errorf("Columns() = %v, want 3 columns starting with nationkey", cols)
	}
	if !tbl.HasColumn("name") {
		t.Error("HasColumn(name) false")
	}
	if tbl.HasColumn("missing") {
		t.Error("HasColumn(missing) true")
	}
	col := tbl.GetColumn("name")
	if col == nil {
		t.Fatal("GetColumn(name) nil")
	}
	if col.Type != "varchar(25)" {
		t.Errorf("col.Type = %q, want varchar(25)", col.Type)
	}
	if !tbl.GetColumn("comment").Nullable {
		t.Error("comment column should be nullable")
	}
	if tbl.GetColumn("missing") != nil {
		t.Error("GetColumn(missing) non-nil")
	}
	if tbl.Schema() != sch {
		t.Error("Table.Schema() back-reference wrong")
	}
}

// TestNilGuards exercises the defensive nil paths: a nil column passed to
// AddTable is skipped, and EnsureCatalogIdent(nil) keys on the empty name.
func TestNilGuards(t *testing.T) {
	sch := catalog.New().EnsureCatalog("c").EnsureSchema("s")
	tbl := sch.AddTable("t", catalog.NewColumn("a", "bigint", false), nil)
	if got := tbl.ColumnNames(); !slices.Equal(got, []string{"a"}) {
		t.Errorf("ColumnNames() = %v, want [a] (nil column skipped)", got)
	}

	c := catalog.New()
	db := c.EnsureCatalogIdent(nil) // normalizes to ""
	if db == nil {
		t.Fatal("EnsureCatalogIdent(nil) returned nil")
	}
	if !c.HasCatalog("") {
		t.Error(`EnsureCatalogIdent(nil) should register the "" catalog`)
	}
}

func TestAddTableReplaces(t *testing.T) {
	sch := catalog.New().EnsureCatalog("c").EnsureSchema("s")
	sch.AddTable("t", catalog.NewColumn("a", "bigint", false))
	sch.AddTable("t", catalog.NewColumn("b", "varchar", false)) // same name
	if got := sch.Tables(); !slices.Equal(got, []string{"t"}) {
		t.Errorf("Tables() = %v, want [t] (replaced not duplicated)", got)
	}
	tbl := sch.GetTable("t")
	if got := tbl.ColumnNames(); !slices.Equal(got, []string{"b"}) {
		t.Errorf("ColumnNames() = %v, want [b] (last definition wins)", got)
	}
}

func TestViews(t *testing.T) {
	sch := catalog.New().EnsureCatalog("c").EnsureSchema("s")
	sch.AddView("active_users",
		catalog.NewColumn("id", "bigint", false),
		catalog.NewColumn("email", "varchar", true),
	)
	if got := sch.Views(); !slices.Equal(got, []string{"active_users"}) {
		t.Errorf("Views() = %v, want [active_users]", got)
	}
	if !sch.HasView("active_users") {
		t.Error("HasView(active_users) false")
	}
	v := sch.GetView("active_users")
	if v == nil {
		t.Fatal("GetView(active_users) nil")
	}
	if got := v.ColumnNames(); !slices.Equal(got, []string{"id", "email"}) {
		t.Errorf("view ColumnNames() = %v, want [id email]", got)
	}
	if cols := v.Columns(); len(cols) != 2 || cols[0].Name != "id" {
		t.Errorf("view Columns() = %v, want 2 columns starting with id", cols)
	}
	if v.Schema() != sch {
		t.Error("View.Schema() back-reference wrong")
	}
	// A view is not a table.
	if sch.HasTable("active_users") {
		t.Error("HasTable returned true for a view name")
	}
	if sch.GetView("missing") != nil {
		t.Error("GetView(missing) non-nil")
	}
}

func TestRelationsMergesAndDedups(t *testing.T) {
	sch := catalog.New().EnsureCatalog("c").EnsureSchema("s")
	sch.AddTable("orders")
	sch.AddTable("lineitem")
	sch.AddView("orders_summary")
	sch.AddView("orders") // shares a name with the table; Relations dedups
	got := sch.Relations()
	want := []string{"lineitem", "orders", "orders_summary"}
	if !slices.Equal(got, want) {
		t.Errorf("Relations() = %v, want %v", got, want)
	}
}

func TestEmptyTableNoColumns(t *testing.T) {
	sch := catalog.New().EnsureCatalog("c").EnsureSchema("s")
	tbl := sch.AddTable("empty")
	if got := tbl.ColumnNames(); len(got) != 0 {
		t.Errorf("ColumnNames() = %v, want empty", got)
	}
	if len(tbl.Columns()) != 0 {
		t.Error("Columns() not empty for a no-column table")
	}
}

// TestQuotedNameDistinctFromUnquoted proves the catalog keeps a case-sensitive
// quoted name distinct from its lower-cased unquoted sibling — the behavior
// blind internal lowercasing would have destroyed.
func TestQuotedNameDistinctFromUnquoted(t *testing.T) {
	sch := catalog.New().EnsureCatalog("c").EnsureSchema("s")
	// Caller normalizes at the boundary: "MyTable" (unquoted) -> mytable;
	// "\"MyTable\"" (quoted) -> MyTable.
	sch.AddTable(catalog.Normalize("MyTable"))   // "mytable"
	sch.AddTable(catalog.Normalize(`"MyTable"`)) // "MyTable"
	if got := sch.Tables(); !slices.Equal(got, []string{"MyTable", "mytable"}) {
		t.Errorf("Tables() = %v, want [MyTable mytable] (quoted distinct from unquoted)", got)
	}
}
