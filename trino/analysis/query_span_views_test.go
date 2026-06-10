package analysis

import (
	"testing"

	"github.com/bytebase/omni/trino/catalog"
)

// lineageCatalog builds the catalog used by the view-lineage tests:
//
//	catalog1.public.customer(id, phone, name, phones)
//	catalog1.public.orders(id, k)
//	catalog1.other.secrets(token)
//	views in catalog1.public:
//	  customer_v   AS SELECT id, phone, name FROM customer
//	  cust_masked  AS SELECT phone AS ph FROM customer
//	  v2           AS SELECT phone FROM customer_v            (view over view)
//	  renamed_v(ph) AS SELECT phone FROM customer             (explicit column list)
//	  comp_v       AS SELECT phone || name AS token FROM customer
//	  cte_v        AS WITH c AS (SELECT phone AS pp FROM customer) SELECT pp FROM c
//	  star_v       AS SELECT * FROM customer                  (star-bodied)
//	  cyc_a        AS SELECT x FROM cyc_b                     (definition cycle)
//	  cyc_b        AS SELECT x FROM cyc_a
//	  dual         — BOTH a table and a view of this name exist (precedence)
//	view in catalog1.other:
//	  other_v      AS SELECT token FROM secrets               (cross-schema context)
//
// Session context: catalog1.public.
func lineageCatalog() *catalog.Catalog {
	cat := catalog.New()
	db := cat.EnsureCatalog("catalog1")
	pub := db.EnsureSchema("public")
	pub.AddTable("customer",
		catalog.NewColumn("id", "integer", false),
		catalog.NewColumn("phone", "varchar", true),
		catalog.NewColumn("name", "varchar", true),
		catalog.NewColumn("phones", "array(varchar)", true),
	)
	pub.AddTable("orders",
		catalog.NewColumn("id", "integer", false),
		catalog.NewColumn("k", "integer", true),
	)
	pub.AddTable("dual", catalog.NewColumn("tcol", "integer", false))

	addView := func(sch *catalog.Schema, name, def string, cols ...*catalog.Column) {
		v := sch.AddView(name, cols...)
		v.Definition = def
	}
	addView(pub, "customer_v", "SELECT id, phone, name FROM customer",
		catalog.NewColumn("id", "integer", false),
		catalog.NewColumn("phone", "varchar", true),
		catalog.NewColumn("name", "varchar", true),
	)
	addView(pub, "cust_masked", "SELECT phone AS ph FROM customer",
		catalog.NewColumn("ph", "varchar", true),
	)
	addView(pub, "v2", "SELECT phone FROM customer_v",
		catalog.NewColumn("phone", "varchar", true),
	)
	addView(pub, "renamed_v", "SELECT phone FROM customer",
		catalog.NewColumn("ph", "varchar", true),
	)
	addView(pub, "comp_v", "SELECT phone || name AS token FROM customer",
		catalog.NewColumn("token", "varchar", true),
	)
	addView(pub, "cte_v", "WITH c AS (SELECT phone AS pp FROM customer) SELECT pp FROM c",
		catalog.NewColumn("pp", "varchar", true),
	)
	addView(pub, "star_v", "SELECT * FROM customer",
		catalog.NewColumn("id", "integer", false),
		catalog.NewColumn("phone", "varchar", true),
		catalog.NewColumn("name", "varchar", true),
		catalog.NewColumn("phones", "array(varchar)", true),
	)
	addView(pub, "cyc_a", "SELECT x FROM cyc_b", catalog.NewColumn("x", "integer", true))
	addView(pub, "cyc_b", "SELECT x FROM cyc_a", catalog.NewColumn("x", "integer", true))
	// stale_v: metadata says two columns but the definition projects one — a
	// stale-metadata disagreement that must make the view opaque.
	addView(pub, "stale_v", "SELECT id FROM customer",
		catalog.NewColumn("id", "integer", false),
		catalog.NewColumn("phone", "varchar", true),
	)
	// alias_v: the definition renames customer's columns via relation column
	// aliases; p positionally is customer.phone.
	addView(pub, "alias_v", "SELECT p FROM customer AS c(i, p, n, ps)",
		catalog.NewColumn("p", "varchar", true),
	)
	addView(pub, "dual", "SELECT phone AS vcol FROM customer", catalog.NewColumn("vcol", "varchar", true))

	other := db.EnsureSchema("other")
	other.AddTable("secrets", catalog.NewColumn("token", "varchar", true))
	addView(other, "other_v", "SELECT token FROM secrets", catalog.NewColumn("token", "varchar", true))

	cat.SetCurrentCatalog("catalog1")
	cat.SetCurrentSchema("public")
	return cat
}

func viewSpan(t *testing.T, sql string) *QuerySpan {
	t.Helper()
	span, err := GetQuerySpanWithCatalog(sql, lineageCatalog())
	if err != nil {
		t.Fatalf("GetQuerySpanWithCatalog returned error: %v", err)
	}
	return span
}

func hasTable(span *QuerySpan, catalogName, schema, table string) bool {
	for _, t := range span.AccessTables {
		if t.Catalog == catalogName && t.Schema == schema && t.Table == table {
			return true
		}
	}
	return false
}

// TestViewLineage_ColumnResolvesToBase covers BYT-9679: a column selected
// through a view resolves to the underlying base column, and the view's base
// table joins AccessTables (qualified) so the consumer can expand it.
func TestViewLineage_ColumnResolvesToBase(t *testing.T) {
	span := viewSpan(t, "SELECT phone FROM customer_v")
	r, ok := resultByName(span, "phone")
	if !ok {
		t.Fatalf("Results = %+v, want a column named phone", span.Results)
	}
	if !hasSource(r.SourceColumns, ColumnRef{Column: "phone"}) {
		t.Errorf("phone.SourceColumns = %+v, want {Column:phone} (through customer_v)", r.SourceColumns)
	}
	if !hasTable(span, "", "", "customer_v") {
		t.Errorf("AccessTables = %+v, want the view customer_v itself (access-control parity)", span.AccessTables)
	}
	if !hasTable(span, "catalog1", "public", "customer") {
		t.Errorf("AccessTables = %+v, want qualified base table catalog1.public.customer", span.AccessTables)
	}
}

// TestViewLineage_StarOverViewExpands expands SELECT * over a view to the
// view's exact projection with base lineage.
func TestViewLineage_StarOverViewExpands(t *testing.T) {
	span := viewSpan(t, "SELECT * FROM customer_v")
	if got := resultNames(span); !sameNames(got, "id", "phone", "name") {
		t.Fatalf("Results = %v, want [id phone name]", got)
	}
	if !hasSource(span.Results[1].SourceColumns, ColumnRef{Column: "phone"}) {
		t.Errorf("Results[1].SourceColumns = %+v, want {Column:phone}", span.Results[1].SourceColumns)
	}
}

// TestViewLineage_RenamingViewResolves resolves a view whose definition renames
// the projected column (phone AS ph).
func TestViewLineage_RenamingViewResolves(t *testing.T) {
	span := viewSpan(t, "SELECT ph FROM cust_masked")
	r, ok := resultByName(span, "ph")
	if !ok {
		t.Fatalf("Results = %+v, want a column named ph", span.Results)
	}
	if !hasSource(r.SourceColumns, ColumnRef{Column: "phone"}) {
		t.Errorf("ph.SourceColumns = %+v, want {Column:phone}", r.SourceColumns)
	}
}

// TestViewLineage_ViewOverViewResolves resolves transitively through a view
// defined over another view.
func TestViewLineage_ViewOverViewResolves(t *testing.T) {
	span := viewSpan(t, "SELECT phone FROM v2")
	r, ok := resultByName(span, "phone")
	if !ok {
		t.Fatalf("Results = %+v, want a column named phone", span.Results)
	}
	if !hasSource(r.SourceColumns, ColumnRef{Column: "phone"}) {
		t.Errorf("phone.SourceColumns = %+v, want {Column:phone} (v2 -> customer_v -> customer)", r.SourceColumns)
	}
	if !hasTable(span, "catalog1", "public", "customer") {
		t.Errorf("AccessTables = %+v, want catalog1.public.customer through both views", span.AccessTables)
	}
}

// TestViewLineage_ExplicitColumnListResolves covers a view whose metadata
// column list renames the definition's output (CREATE VIEW renamed_v (ph) AS
// SELECT phone ...): the outer ph maps positionally to the base column.
func TestViewLineage_ExplicitColumnListResolves(t *testing.T) {
	span := viewSpan(t, "SELECT ph FROM renamed_v")
	r, ok := resultByName(span, "ph")
	if !ok {
		t.Fatalf("Results = %+v, want a column named ph", span.Results)
	}
	if !hasSource(r.SourceColumns, ColumnRef{Column: "phone"}) {
		t.Errorf("ph.SourceColumns = %+v, want {Column:phone} (positional metadata column list)", r.SourceColumns)
	}
}

// TestViewLineage_ComposedColumnBothSources resolves a view column composed of
// two base columns to both.
func TestViewLineage_ComposedColumnBothSources(t *testing.T) {
	span := viewSpan(t, "SELECT token FROM comp_v")
	r, ok := resultByName(span, "token")
	if !ok {
		t.Fatalf("Results = %+v, want a column named token", span.Results)
	}
	if !hasSource(r.SourceColumns, ColumnRef{Column: "phone"}) || !hasSource(r.SourceColumns, ColumnRef{Column: "name"}) {
		t.Errorf("token.SourceColumns = %+v, want both {Column:phone} and {Column:name}", r.SourceColumns)
	}
}

// TestViewLineage_CTEDefinitionResolves resolves a view whose definition itself
// uses a CTE — the case the consumer-side post-pass could not reach with a
// shallow analyzer.
func TestViewLineage_CTEDefinitionResolves(t *testing.T) {
	span := viewSpan(t, "SELECT pp FROM cte_v")
	r, ok := resultByName(span, "pp")
	if !ok {
		t.Fatalf("Results = %+v, want a column named pp", span.Results)
	}
	if !hasSource(r.SourceColumns, ColumnRef{Column: "phone"}) {
		t.Errorf("pp.SourceColumns = %+v, want {Column:phone} (through the view's CTE)", r.SourceColumns)
	}
}

// TestViewLineage_StarBodiedViewResolves resolves a view defined as SELECT *
// over a catalog-known base table: the definition's star expands against the
// catalog, so a named reference through the view reaches the base column.
func TestViewLineage_StarBodiedViewResolves(t *testing.T) {
	span := viewSpan(t, "SELECT phone FROM star_v")
	r, ok := resultByName(span, "phone")
	if !ok {
		t.Fatalf("Results = %+v, want a column named phone", span.Results)
	}
	if !hasSource(r.SourceColumns, ColumnRef{Catalog: "catalog1", Schema: "public", Table: "customer", Column: "phone"}) {
		t.Errorf("phone.SourceColumns = %+v, want qualified catalog1.public.customer.phone (def star expanded via catalog)", r.SourceColumns)
	}
}

// TestViewLineage_CycleTerminatesOpaque guards that a definition cycle (a -> b
// -> a, malformed metadata) terminates and leaves the views opaque rather than
// hanging or panicking.
func TestViewLineage_CycleTerminatesOpaque(t *testing.T) {
	span := viewSpan(t, "SELECT x FROM cyc_a")
	r, ok := resultByName(span, "x")
	if !ok {
		t.Fatalf("Results = %+v, want a column named x", span.Results)
	}
	// The cycle prevents resolution; the written ref is all that remains.
	if !hasSource(r.SourceColumns, ColumnRef{Column: "x"}) {
		t.Errorf("x.SourceColumns = %+v, want the original {Column:x} retained", r.SourceColumns)
	}
}

// TestViewLineage_CrossSchemaContext resolves a view in another schema: its
// definition's unqualified table reference resolves in the VIEW's schema, not
// the session schema, and surfaces qualified.
func TestViewLineage_CrossSchemaContext(t *testing.T) {
	span := viewSpan(t, "SELECT token FROM other.other_v")
	r, ok := resultByName(span, "token")
	if !ok {
		t.Fatalf("Results = %+v, want a column named token", span.Results)
	}
	if !hasSource(r.SourceColumns, ColumnRef{Column: "token"}) {
		t.Errorf("token.SourceColumns = %+v, want {Column:token}", r.SourceColumns)
	}
	if !hasTable(span, "catalog1", "other", "secrets") {
		t.Errorf("AccessTables = %+v, want catalog1.other.secrets (resolved in the view's schema)", span.AccessTables)
	}
}

// TestViewLineage_TablePrecedenceOverView guards catalog.ResolveRelation
// parity: a table shadows a same-named view, so the reference binds to the
// table (star expands to the TABLE's columns, not the view's projection).
func TestViewLineage_TablePrecedenceOverView(t *testing.T) {
	span := viewSpan(t, "SELECT * FROM dual")
	if got := resultNames(span); !sameNames(got, "tcol") {
		t.Fatalf("Results = %v, want [tcol] (table shadows the same-named view)", got)
	}
}

// TestViewLineage_BaseTableStarExpandsViaCatalog covers the catalog-table star:
// SELECT * over a catalog-known base table expands to its exact projection with
// fully-qualified sources.
func TestViewLineage_BaseTableStarExpandsViaCatalog(t *testing.T) {
	span := viewSpan(t, "SELECT * FROM customer")
	if got := resultNames(span); !sameNames(got, "id", "phone", "name", "phones") {
		t.Fatalf("Results = %v, want [id phone name phones]", got)
	}
	want := ColumnRef{Catalog: "catalog1", Schema: "public", Table: "customer", Column: "phone"}
	if !hasSource(span.Results[1].SourceColumns, want) {
		t.Errorf("Results[1].SourceColumns = %+v, want %+v", span.Results[1].SourceColumns, want)
	}
}

// TestViewLineage_MixedBaseDerivedStarExpands covers the formerly-opaque mixed
// case: with a catalog, a star over a base table joined with a derived relation
// expands at the true total width.
func TestViewLineage_MixedBaseDerivedStarExpands(t *testing.T) {
	span := viewSpan(t, "SELECT * FROM orders JOIN (SELECT phone FROM customer) d ON true")
	if got := resultNames(span); !sameNames(got, "id", "k", "phone") {
		t.Fatalf("Results = %v, want [id k phone]", got)
	}
	if !hasSource(span.Results[2].SourceColumns, ColumnRef{Column: "phone"}) {
		t.Errorf("Results[2].SourceColumns = %+v, want {Column:phone}", span.Results[2].SourceColumns)
	}
}

// TestViewLineage_UnknownNameParity guards that a name the catalog cannot
// resolve binds exactly as without a catalog (written-form ref, opaque star).
func TestViewLineage_UnknownNameParity(t *testing.T) {
	span := viewSpan(t, "SELECT * FROM mystery")
	if got := resultNames(span); !sameNames(got, "*") {
		t.Fatalf("Results = %v, want [*] (unknown relation stays opaque)", got)
	}
}

// TestViewLineage_TableColumnAliasResolves covers relation column aliases over
// a catalog table: FROM customer AS c(i, p, n, ps) renames the columns
// positionally, so p (bare or c.p) must reach customer.phone.
func TestViewLineage_TableColumnAliasResolves(t *testing.T) {
	want := ColumnRef{Catalog: "catalog1", Schema: "public", Table: "customer", Column: "phone"}
	span := viewSpan(t, "SELECT p FROM customer AS c(i, p, n, ps)")
	r, ok := resultByName(span, "p")
	if !ok {
		t.Fatalf("Results = %+v, want a column named p", span.Results)
	}
	if !hasSource(r.SourceColumns, want) {
		t.Errorf("p.SourceColumns = %+v, want %+v (positional relation column alias)", r.SourceColumns, want)
	}

	span = viewSpan(t, "SELECT c.p FROM customer AS c(i, p, n, ps)")
	r, ok = resultByName(span, "p")
	if !ok {
		t.Fatalf("Results = %+v, want a column named p", span.Results)
	}
	if !hasSource(r.SourceColumns, want) {
		t.Errorf("c.p SourceColumns = %+v, want %+v", r.SourceColumns, want)
	}
}

// TestViewLineage_ViewOverAliasedTableResolves resolves a view whose definition
// renames the base table's columns via relation column aliases.
func TestViewLineage_ViewOverAliasedTableResolves(t *testing.T) {
	span := viewSpan(t, "SELECT p FROM alias_v")
	r, ok := resultByName(span, "p")
	if !ok {
		t.Fatalf("Results = %+v, want a column named p", span.Results)
	}
	want := ColumnRef{Catalog: "catalog1", Schema: "public", Table: "customer", Column: "phone"}
	if !hasSource(r.SourceColumns, want) {
		t.Errorf("p.SourceColumns = %+v, want %+v (through alias_v's renamed projection)", r.SourceColumns, want)
	}
}

// TestViewLineage_StaleMetadataCountMismatchOpaque guards that a view whose
// metadata column count disagrees with its analyzed definition (stale metadata)
// is unresolvable: the true output shape is unknown, so expanding on either
// width could misalign the positional masker.
func TestViewLineage_StaleMetadataCountMismatchOpaque(t *testing.T) {
	span := viewSpan(t, "SELECT * FROM stale_v")
	if got := resultNames(span); !sameNames(got, "*") {
		t.Fatalf("Results = %v, want exactly [*] (metadata/definition disagree; must stay opaque)", got)
	}
	span = viewSpan(t, "SELECT phone FROM stale_v")
	r, ok := resultByName(span, "phone")
	if !ok {
		t.Fatalf("Results = %+v, want a column named phone", span.Results)
	}
	if len(r.SourceColumns) != 1 || r.SourceColumns[0] != (ColumnRef{Column: "phone"}) {
		t.Errorf("phone.SourceColumns = %+v, want exactly the written [{Column:phone}] (no projection trusted)", r.SourceColumns)
	}
}

// TestViewLineage_CycleStarStaysOpaque strengthens the cycle guard: a star over
// a cyclic view must not expand from a partial (mid-cycle) projection.
func TestViewLineage_CycleStarStaysOpaque(t *testing.T) {
	span := viewSpan(t, "SELECT * FROM cyc_a")
	if got := resultNames(span); !sameNames(got, "*") {
		t.Fatalf("Results = %v, want exactly [*] (cyclic views are fully opaque)", got)
	}
	// Both orders: referencing cyc_b first must not seed a partial memo either.
	span = viewSpan(t, "SELECT * FROM cyc_b JOIN cyc_a ON true")
	if got := resultNames(span); !sameNames(got, "*") {
		t.Fatalf("Results = %v, want exactly [*] (no partial projection memoized)", got)
	}
}

// TestViewLineage_BaseColumnRefRetained guards additivity for a plain column
// reference against a catalog-known base table: the written ref is always
// retained, with the catalog-qualified form added beside it (the addition is
// what resolves relation column aliases; for an unaliased table it merely
// restates the same column).
func TestViewLineage_BaseColumnRefRetained(t *testing.T) {
	span := viewSpan(t, "SELECT phone FROM customer")
	r, ok := resultByName(span, "phone")
	if !ok {
		t.Fatalf("Results = %+v, want a column named phone", span.Results)
	}
	if !hasSource(r.SourceColumns, ColumnRef{Column: "phone"}) {
		t.Errorf("phone.SourceColumns = %+v, want the written {Column:phone} retained", r.SourceColumns)
	}
	if !hasSource(r.SourceColumns, ColumnRef{Catalog: "catalog1", Schema: "public", Table: "customer", Column: "phone"}) {
		t.Errorf("phone.SourceColumns = %+v, want the qualified catalog form added", r.SourceColumns)
	}
}
