package catalog

import (
	"slices"
	"testing"

	"github.com/bytebase/omni/metadata"
)

func redshiftSnapshot() *metadata.DatabaseSchemaMetadata {
	return &metadata.DatabaseSchemaMetadata{
		Name: "dev",
		Schemas: []*metadata.SchemaMetadata{{
			Name: "sales",
			Tables: []*metadata.TableMetadata{
				{
					Name: "orders",
					Columns: []*metadata.ColumnMetadata{
						{Name: "id", Type: "bigint"},
						{Name: "note", Type: "character varying(256)"},
						{Name: "geo", Type: "geometry"},
						{Name: "ttl", Type: "interval"},
					},
				},
				{Name: "empty"},
			},
			ExternalTables: []*metadata.ExternalTableMetadata{{
				Name:    "clicks",
				Columns: []*metadata.ColumnMetadata{{Name: "at", Type: "timestamp without time zone"}},
			}},
			Views: []*metadata.ViewMetadata{
				{Name: "recent", Definition: "SELECT id FROM orders", Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "bigint"}}},
				// No definition: it resolves as a table of its columns.
				{Name: "bare", Columns: []*metadata.ColumnMetadata{{Name: "x", Type: "integer"}}},
				// A definition that does not parse resolves the same way.
				{Name: "odd", Definition: "SELEC nonsense", Columns: []*metadata.ColumnMetadata{{Name: "x", Type: "integer"}}},
				// A definition that parses but does not install falls back too.
				{Name: "unanalyzable", Definition: "SELECT no_such_fn(id) AS x FROM orders", Columns: []*metadata.ColumnMetadata{{Name: "x", Type: "integer"}}},
			},
			MaterializedViews: []*metadata.MaterializedViewMetadata{
				{Name: "totals", Definition: "SELECT count(*) AS n FROM orders"},
				{Name: "broken_mv", Definition: "SELEC nonsense"},
				{Name: "stale_mv", Definition: "SELECT id FROM nowhere"},
				{Name: "undefined_mv"},
			},
			Sequences: []*metadata.SequenceMetadata{{Name: "order_ids"}},
			Functions: []*metadata.FunctionMetadata{{
				Name:       "double_it",
				Definition: "CREATE FUNCTION sales.double_it(integer) RETURNS integer STABLE AS $$ SELECT $1 * 2 $$ LANGUAGE sql;",
			}},
		}},
	}
}

func TestLoadMetadataInstallsSnapshot(t *testing.T) {
	c := New()
	c.SetSearchPath([]string{"public"})
	c.LoadMetadata(redshiftSnapshot())

	for _, name := range []string{"orders", "empty", "clicks", "recent", "bare", "totals", "broken_mv", "stale_mv", "undefined_mv"} {
		if c.GetRelation("sales", name) == nil {
			t.Fatalf("relation sales.%s not installed", name)
		}
	}
	var types []string
	for _, col := range c.GetRelation("sales", "orders").Columns {
		types = append(types, c.FormatType(col.TypeOID, col.TypeMod))
	}
	if want := []string{"bigint", "text", "text", "interval"}; !slices.Equal(types, want) {
		t.Errorf("orders column types = %v, want the coarse %v", types, want)
	}
	if recent := c.GetRelation("sales", "recent"); recent.RelKind != 'v' || len(recent.Columns) != 1 || recent.Columns[0].Name != "id" {
		t.Errorf("recent = kind %c, columns %v; want a view of id", recent.RelKind, recent.Columns)
	}
	// A definition that runs installs as written; one that does not falls back.
	if totals := c.GetRelation("sales", "totals"); totals.RelKind != 'm' || len(totals.Columns) != 1 || totals.Columns[0].Name != "n" {
		t.Errorf("totals = kind %c, columns %v; want the materialized view of n", totals.RelKind, totals.Columns)
	}
	for _, name := range []string{"broken_mv", "stale_mv", "undefined_mv"} {
		if mv := c.GetRelation("sales", name); mv.RelKind != 'm' || len(mv.Columns) != 1 || mv.Columns[0].Name != metadataPlaceholderColumn {
			t.Errorf("%s = kind %c, columns %v; want a materialized view stand-in", name, mv.RelKind, mv.Columns)
		}
	}
	if c.GetSchema("sales").Sequences["order_ids"] == nil {
		t.Error("sequence sales.order_ids not installed")
	}
	if !slices.Equal(c.searchPath, []string{"public"}) {
		t.Errorf("search path = %v, want it restored to public", c.searchPath)
	}
}

func TestLoadMetadataInstallsMaterializedViewsAfterWhatTheyName(t *testing.T) {
	c := New()
	c.LoadMetadata(&metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{
		{Name: "reports", MaterializedViews: []*metadata.MaterializedViewMetadata{
			// a_totals names z_counts, which names a table of a later schema.
			{Name: "a_totals", Definition: "SELECT n FROM reports.z_counts"},
			{Name: "z_counts", Definition: "SELECT count(*) AS n FROM sales.orders"},
			// Calls a function of a later schema.
			{Name: "doubled", Definition: "SELECT sales.double_it(id) AS n FROM sales.orders"},
			// needs_broken names one that can never install.
			{Name: "broken_mv", Definition: "SELEC nonsense"},
			{Name: "needs_broken", Definition: "SELECT count(*) AS n FROM reports.broken_mv"},
		}},
		{
			Name:      "sales",
			Tables:    []*metadata.TableMetadata{{Name: "orders", Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "bigint"}}}},
			Functions: []*metadata.FunctionMetadata{{Name: "double_it", Definition: "CREATE FUNCTION sales.double_it(bigint) RETURNS bigint STABLE AS $$ SELECT $1 * 2 $$ LANGUAGE sql;"}},
		},
	}})
	for _, name := range []string{"a_totals", "z_counts", "doubled", "needs_broken"} {
		if mv := c.GetRelation("reports", name); mv == nil || len(mv.Columns) != 1 || mv.Columns[0].Name != "n" {
			t.Errorf("%s did not install from its definition", name)
		}
	}
	if mv := c.GetRelation("reports", "broken_mv"); mv == nil || len(mv.Columns) != 1 || mv.Columns[0].Name != metadataPlaceholderColumn {
		t.Errorf("broken_mv = %+v, want the stand-in", mv)
	}
}

func TestUseMetadataDegradesAViewThatDoesNotInstall(t *testing.T) {
	c := New()
	c.SetSearchPath([]string{"sales", "public"})
	c.UseMetadata(redshiftSnapshot())
	if _, err := c.Exec("CREATE VIEW public.u AS SELECT x FROM unanalyzable;", nil); err != nil {
		t.Fatalf("Exec: %v", err)
	}
	if c.GetRelation("public", "u") == nil {
		t.Error("a view whose definition does not install did not fall back to its columns")
	}
}

func TestUseMetadataResolvesRelationsWhenNamed(t *testing.T) {
	c := New()
	// $user expands to the session user.
	c.SetSessionUser("sales")
	c.SetSearchPath([]string{"$user", "public"})
	c.UseMetadata(redshiftSnapshot())
	if c.GetRelation("sales", "orders") != nil {
		t.Fatal("a relation must resolve only when a statement names it")
	}
	if _, err := c.Exec("CREATE VIEW public.report AS SELECT r.id, t.n FROM recent r, totals t;", nil); err != nil {
		t.Fatalf("CREATE VIEW over resolved relations: %v", err)
	}
	for _, name := range []string{"recent", "totals", "orders"} {
		if c.GetRelation("sales", name) == nil {
			t.Errorf("sales.%s not resolved", name)
		}
	}
	if procs := c.LookupProcByName("double_it"); len(procs) != 1 {
		t.Errorf("double_it installed %d times, want once", len(procs))
	}
}

func TestUseMetadataResolvesEarlierSchemasFirst(t *testing.T) {
	c := New()
	c.SetSearchPath([]string{"tenant", "public"})
	c.UseMetadata(&metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{
		{Name: "tenant", Tables: []*metadata.TableMetadata{{Name: "accounts", Columns: []*metadata.ColumnMetadata{{Name: "tenant_id", Type: "integer"}}}}},
		{Name: "public", Tables: []*metadata.TableMetadata{{Name: "accounts", Columns: []*metadata.ColumnMetadata{{Name: "legacy_id", Type: "integer"}}}}},
	}})
	// public.accounts is in the catalog before the unqualified name resolves.
	if _, err := c.Exec("CREATE VIEW public.a AS SELECT legacy_id FROM public.accounts; CREATE VIEW public.b AS SELECT tenant_id FROM accounts;", nil); err != nil {
		t.Fatalf("Exec: %v", err)
	}
	if c.GetRelation("public", "a") == nil {
		t.Fatal("public.a was not created")
	}
	if c.GetRelation("public", "b") == nil {
		t.Error("accounts bound to public.accounts, not tenant.accounts earlier on the search path")
	}
}

func TestMetadataResolverResolvesSpecs(t *testing.T) {
	r := newMetadataResolver(redshiftSnapshot())
	resolve := func(schema, name string, searchPath ...string) *RelationSpec {
		t.Helper()
		spec, err := r.ResolveRelation(schema, name, searchPath)
		if err != nil {
			t.Fatalf("ResolveRelation(%q, %q): %v", schema, name, err)
		}
		return spec
	}
	columns := func(spec *RelationSpec) []RelationColumnSpec {
		if spec == nil {
			return nil
		}
		return spec.Columns
	}

	if spec := resolve("", "recent", "sales", "public"); spec == nil || spec.SchemaName != "sales" || spec.Kind != RelationKindView || spec.Definition != `SELECT id FROM "sales".orders` {
		t.Errorf("recent = %+v, want a view with its relation qualified", spec)
	}
	if got, want := columns(resolve("sales", "orders")), []RelationColumnSpec{{"id", "bigint"}, {"note", "text"}, {"geo", "text"}, {"ttl", "interval"}}; !slices.Equal(got, want) {
		t.Errorf("orders columns = %v, want %v", got, want)
	}
	for _, name := range []string{"bare", "odd"} {
		if spec := resolve("", name, "sales"); spec == nil || spec.Kind != RelationKindTable || !slices.Equal(spec.Columns, []RelationColumnSpec{{"x", "integer"}}) {
			t.Errorf("%s = %+v, want a table of its columns", name, spec)
		}
	}
	if got := columns(resolve("", "order_ids", "sales")); len(got) != 3 || got[0].Name != "last_value" {
		t.Errorf("order_ids columns = %v, want the sequence's", got)
	}
	if got := columns(resolve("", "clicks", "public", "sales")); !slices.Equal(got, []RelationColumnSpec{{"at", "timestamp"}}) {
		t.Errorf("clicks columns = %v, want the external table's", got)
	}
	if got := columns(resolve("", "undefined_mv", "sales")); !slices.Equal(got, []RelationColumnSpec{{metadataPlaceholderColumn, "text"}}) {
		t.Errorf("undefined_mv columns = %v, want the placeholder", got)
	}
	for _, tt := range []struct {
		schema, name string
		searchPath   []string
	}{
		{"", "orders", []string{"public"}},
		{"sales", "Orders", nil},
		{"missing", "orders", nil},
	} {
		if spec := resolve(tt.schema, tt.name, tt.searchPath...); spec != nil {
			t.Errorf("ResolveRelation(%q, %q, %v) = %+v, want nil", tt.schema, tt.name, tt.searchPath, spec)
		}
	}
}

func TestMetadataResolverQualifiesViewDefinitions(t *testing.T) {
	r := newMetadataResolver(redshiftSnapshot())
	for definition, want := range map[string]string{
		"SELECT id FROM orders": `SELECT id FROM "sales".orders`,
		"SELECT id FROM orders WHERE id IN (SELECT id FROM recent)":                                       `SELECT id FROM "sales".orders WHERE id IN (SELECT id FROM "sales".recent)`,
		"WITH orders AS (SELECT 1 AS id) SELECT id FROM orders":                                           "WITH orders AS (SELECT 1 AS id) SELECT id FROM orders",
		"WITH o AS (SELECT id FROM orders) SELECT id FROM o":                                              `WITH o AS (SELECT id FROM "sales".orders) SELECT id FROM o`,
		`WITH "Orders" AS (SELECT 1 AS id) SELECT id FROM orders`:                                         `WITH "Orders" AS (SELECT 1 AS id) SELECT id FROM "sales".orders`,
		"WITH orders AS (SELECT id FROM orders) SELECT id FROM orders":                                    `WITH orders AS (SELECT id FROM "sales".orders) SELECT id FROM orders`,
		"WITH orders AS (SELECT 1 AS id), later AS (SELECT id FROM orders) SELECT id FROM later":          "WITH orders AS (SELECT 1 AS id), later AS (SELECT id FROM orders) SELECT id FROM later",
		"WITH RECURSIVE orders AS (SELECT 1 AS id UNION ALL SELECT id FROM orders) SELECT id FROM orders": "WITH RECURSIVE orders AS (SELECT 1 AS id UNION ALL SELECT id FROM orders) SELECT id FROM orders",
		"CREATE MATERIALIZED VIEW sales.m AS SELECT id FROM orders":                                       `CREATE MATERIALIZED VIEW sales.m AS SELECT id FROM "sales".orders`,
		"SELECT id FROM sales.orders":                                                                     "SELECT id FROM sales.orders",
		"SELECT id FROM unknown":                                                                          "SELECT id FROM unknown",
	} {
		if got, _ := r.qualify(definition, "sales"); got != want {
			t.Errorf("qualify(%q) = %q, want %q", definition, got, want)
		}
	}
	if _, ok := r.qualify("SELEC nonsense", "sales"); ok {
		t.Error("qualify accepted a definition that does not parse")
	}
}

func TestCoarseType(t *testing.T) {
	for typ, want := range map[string]string{
		"":                            "text",
		"bigint":                      "bigint",
		"smallint":                    "integer",
		"boolean":                     "boolean",
		"character varying(256)":      "text",
		"numeric(10,2)":               "numeric",
		"double precision":            "double precision",
		"real":                        "real",
		"date":                        "date",
		"timestamp without time zone": "timestamp",
		"timestamp with time zone":    "timestamp",
		"time":                        "time",
		"time without time zone":      "time",
		"timetz":                      "time",
		"super":                       "text",
		"int8":                        "bigint",
		"float8":                      "double precision",
		"float":                       "double precision",
		"float4":                      "real",
		"interval":                    "interval",
		"interval year to month":      "interval",
	} {
		if got := coarseType(typ); got != want {
			t.Errorf("coarseType(%q) = %q, want %q", typ, got, want)
		}
	}
}
