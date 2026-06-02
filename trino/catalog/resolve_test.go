package catalog_test

import (
	"slices"
	"testing"

	"github.com/bytebase/omni/trino/ast"
	"github.com/bytebase/omni/trino/catalog"
)

// buildResolveFixture creates a catalog tpch.sf1 with a "nation" table and an
// "orders_view" view, plus a second catalog hive.default.events table, so the
// resolution tests can exercise 1/2/3-part names and cross-catalog lookups.
func buildResolveFixture() *catalog.Catalog {
	c := catalog.New()
	sf1 := c.EnsureCatalog("tpch").EnsureSchema("sf1")
	sf1.AddTable("nation",
		catalog.NewColumn("nationkey", "bigint", false),
		catalog.NewColumn("name", "varchar(25)", false),
	)
	sf1.AddView("orders_view", catalog.NewColumn("total", "double", true))

	c.EnsureCatalog("hive").EnsureSchema("default").
		AddTable("events", catalog.NewColumn("ts", "timestamp", false))
	return c
}

func TestResolveThreePart(t *testing.T) {
	c := buildResolveFixture()
	// Fully-qualified: independent of session context.
	r := c.ResolveRelation([]string{"tpch", "sf1", "nation"})
	if !r.Found || r.Table == nil {
		t.Fatalf("3-part resolve of tpch.sf1.nation failed: %+v", r)
	}
	if r.Catalog != "tpch" || r.Schema != "sf1" || r.Table.Name != "nation" {
		t.Errorf("resolved to %s.%s.%s, want tpch.sf1.nation", r.Catalog, r.Schema, r.Table.Name)
	}
	if got := r.ColumnNames(); !slices.Equal(got, []string{"nationkey", "name"}) {
		t.Errorf("ColumnNames() = %v, want [nationkey name]", got)
	}
}

func TestResolveThreePartCrossCatalog(t *testing.T) {
	c := buildResolveFixture()
	// Session points at tpch, but a 3-part name reaches hive.
	c.SetCurrentCatalog("tpch")
	c.SetCurrentSchema("sf1")
	r := c.ResolveRelation([]string{"hive", "default", "events"})
	if !r.Found || r.Table == nil || r.Catalog != "hive" {
		t.Fatalf("3-part cross-catalog resolve failed: %+v", r)
	}
}

func TestResolveTwoPartUsesCurrentCatalog(t *testing.T) {
	c := buildResolveFixture()
	c.SetCurrentCatalog("tpch")
	r := c.ResolveRelation([]string{"sf1", "nation"})
	if !r.Found || r.Table == nil {
		t.Fatalf("2-part resolve of sf1.nation (current catalog tpch) failed: %+v", r)
	}
	if r.Catalog != "tpch" || r.Schema != "sf1" {
		t.Errorf("resolved to %s.%s, want tpch.sf1", r.Catalog, r.Schema)
	}
}

func TestResolveTwoPartNoCurrentCatalog(t *testing.T) {
	c := buildResolveFixture()
	// No current catalog set -> a 2-part name cannot resolve.
	r := c.ResolveRelation([]string{"sf1", "nation"})
	if r.Found {
		t.Errorf("2-part resolve should fail with no current catalog, got %+v", r)
	}
	if r.ColumnNames() != nil {
		t.Error("ColumnNames() should be nil for an unresolved relation")
	}
}

func TestResolveOnePartUsesSession(t *testing.T) {
	c := buildResolveFixture()
	c.SetCurrentCatalog("tpch")
	c.SetCurrentSchema("sf1")
	r := c.ResolveRelation([]string{"nation"})
	if !r.Found || r.Table == nil {
		t.Fatalf("1-part resolve of nation (session tpch.sf1) failed: %+v", r)
	}
	if r.Catalog != "tpch" || r.Schema != "sf1" || r.Table.Name != "nation" {
		t.Errorf("resolved to %s.%s.%s, want tpch.sf1.nation", r.Catalog, r.Schema, r.Table.Name)
	}
}

func TestResolveOnePartMissingSchema(t *testing.T) {
	c := buildResolveFixture()
	c.SetCurrentCatalog("tpch") // schema NOT set
	r := c.ResolveRelation([]string{"nation"})
	if r.Found {
		t.Errorf("1-part resolve should fail with no current schema, got %+v", r)
	}
}

func TestResolveViaView(t *testing.T) {
	c := buildResolveFixture()
	c.SetCurrentCatalog("tpch")
	c.SetCurrentSchema("sf1")
	r := c.ResolveRelation([]string{"orders_view"})
	if !r.Found || r.View == nil {
		t.Fatalf("1-part resolve of a view failed: %+v", r)
	}
	if r.Table != nil {
		t.Error("a view resolution must leave Table nil")
	}
	if got := r.ColumnNames(); !slices.Equal(got, []string{"total"}) {
		t.Errorf("view ColumnNames() = %v, want [total]", got)
	}
}

// TestResolveTablePreferredOverView pins the table-over-view precedence when a
// schema holds both names (matches Schema.Relations dedup ordering).
func TestResolveTablePreferredOverView(t *testing.T) {
	c := catalog.New()
	s := c.EnsureCatalog("c").EnsureSchema("s")
	s.AddTable("dup", catalog.NewColumn("tcol", "bigint", false))
	s.AddView("dup", catalog.NewColumn("vcol", "bigint", false))
	c.SetCurrentCatalog("c")
	c.SetCurrentSchema("s")
	r := c.ResolveRelation([]string{"dup"})
	if !r.Found || r.Table == nil || r.View != nil {
		t.Fatalf("table should be preferred over a same-named view: %+v", r)
	}
}

func TestResolveUnknownNames(t *testing.T) {
	c := buildResolveFixture()
	c.SetCurrentCatalog("tpch")
	c.SetCurrentSchema("sf1")
	cases := [][]string{
		{"nope"},                      // unknown table in valid schema
		{"missing_schema", "nation"},  // unknown schema
		{"nope", "sf1", "nation"},     // unknown catalog
		{"tpch", "missing", "nation"}, // unknown schema (3-part)
		{"tpch", "sf1", "nope"},       // unknown table (3-part)
	}
	for _, parts := range cases {
		if r := c.ResolveRelation(parts); r.Found {
			t.Errorf("ResolveRelation(%v) unexpectedly found %+v", parts, r)
		}
	}
}

func TestResolvePartCountBounds(t *testing.T) {
	c := buildResolveFixture()
	c.SetCurrentCatalog("tpch")
	c.SetCurrentSchema("sf1")
	// Zero parts and >3 parts are unsupported.
	if r := c.ResolveRelation(nil); r.Found {
		t.Errorf("ResolveRelation(nil) found %+v", r)
	}
	if r := c.ResolveRelation([]string{}); r.Found {
		t.Errorf("ResolveRelation([]) found %+v", r)
	}
	if r := c.ResolveRelation([]string{"a", "b", "c", "d"}); r.Found {
		t.Errorf("ResolveRelation(4 parts) found %+v", r)
	}
}

// TestResolveQualifiedNameAST drives resolution from real ast.QualifiedName
// nodes, proving the catalog interoperates with parser output: unquoted parts
// fold (so "TPCH"."SF1"."NATION" written unquoted resolves to tpch.sf1.nation),
// and a quoted part is matched case-sensitively.
func TestResolveQualifiedNameAST(t *testing.T) {
	c := buildResolveFixture()

	// Unquoted, mixed-case 3-part name -> folds and resolves.
	qn := &ast.QualifiedName{Parts: []*ast.Identifier{
		{Value: "TPCH"}, {Value: "Sf1"}, {Value: "Nation"},
	}}
	if r := c.ResolveQualifiedName(qn); !r.Found || r.Table == nil {
		t.Fatalf("ResolveQualifiedName(unquoted mixed case) failed: %+v", r)
	}

	// 1-part unquoted name with session context.
	c.SetCurrentCatalog("tpch")
	c.SetCurrentSchema("sf1")
	one := &ast.QualifiedName{Parts: []*ast.Identifier{{Value: "NATION"}}}
	if r := c.ResolveQualifiedName(one); !r.Found || r.Table == nil {
		t.Fatalf("ResolveQualifiedName(1-part) failed: %+v", r)
	}

	// A quoted, case-mismatched component must NOT resolve to the lower-cased
	// stored object (quoted identifiers are case-sensitive).
	quoted := &ast.QualifiedName{Parts: []*ast.Identifier{
		{Value: "TPCH"}, {Value: "SF1", Quoted: true, QuoteRune: '"'}, {Value: "nation"},
	}}
	if r := c.ResolveQualifiedName(quoted); r.Found {
		t.Errorf("quoted SF1 should not match stored sf1, got %+v", r)
	}

	// Nil qualified name.
	if r := c.ResolveQualifiedName(nil); r.Found {
		t.Errorf("ResolveQualifiedName(nil) found %+v", r)
	}
}
