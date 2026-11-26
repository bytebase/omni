package analysis

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
)

func names(set SourceColumnSet) []string {
	var out []string
	for c := range set {
		out = append(out, c.String())
	}
	slices.Sort(out)
	return out
}

// Without a catalog every table's columns are unknown: a column resolves to
// each table it could belong to, and nothing is reported missing.
func TestResolveQuerySpanWithoutCatalog(t *testing.T) {
	got, err := ResolveQuerySpan("SELECT o.amount, note FROM shop.orders AS o JOIN notes ON notes.id = o.id", ResolveOptions{DefaultDatabase: "shop"})
	if err != nil {
		t.Fatal(err)
	}
	if got.NotFoundError != nil {
		t.Fatalf("NotFoundError %v", got.NotFoundError)
	}
	if want := []string{"shop.notes", "shop.orders"}; !slices.Equal(names(got.SourceColumns), want) {
		t.Errorf("tables %v, want %v", names(got.SourceColumns), want)
	}
	want := [][]string{{"shop.orders.amount"}, {"shop.notes.note", "shop.orders.note"}}
	for i, r := range got.Results {
		if !slices.Equal(names(r.SourceColumns), want[i]) {
			t.Errorf("result %d (%s) reads %v, want %v", i, r.Name, names(r.SourceColumns), want[i])
		}
	}
}

type oneDatabase struct{ view *metadata.ViewMetadata }

func (oneDatabase) GetTable(string) *metadata.TableMetadata { return nil }
func (d oneDatabase) GetView(string) *metadata.ViewMetadata { return d.view }

// A view whose stored definition does not parse fails the span rather than
// omit what the view reads.
func TestResolveQuerySpanUnreadableView(t *testing.T) {
	db := oneDatabase{view: &metadata.ViewMetadata{Name: "v", Definition: "CREATE VIEW shop.v AS SELEC", Columns: []*metadata.ColumnMetadata{{Name: "a"}}}}
	_, err := ResolveQuerySpan("SELECT a FROM v", ResolveOptions{
		DefaultDatabase: "shop",
		Catalog:         Catalog{GetDatabase: func(string) (Database, error) { return db, nil }},
	})
	if err == nil || !strings.Contains(err.Error(), "view shop.v") {
		t.Fatalf("err %v, want the view's parse error", err)
	}
}

// A statement that is not a query has its type and no span.
func TestResolveQuerySpanNotAQuery(t *testing.T) {
	got, err := ResolveQuerySpan("INSERT INTO t VALUES (1)", ResolveOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if got.Type != analysis.DML || len(got.Results) != 0 || len(got.SourceColumns) != 0 {
		t.Fatalf("got %+v", got)
	}
}

// A window that inherits itself, which ClickHouse refuses, still ends.
func TestResolveQuerySpanWindowCycle(t *testing.T) {
	if _, err := ResolveQuerySpan("SELECT count() OVER w FROM t WINDOW w AS (w)", ResolveOptions{}); err != nil {
		t.Fatal(err)
	}
}

// A quoted name holding a dot is one name.
func TestResolveQuerySpanQuotedDottedName(t *testing.T) {
	got, err := ResolveQuerySpan("SELECT `a.b` FROM (SELECT 1 AS `a.b`)", ResolveOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Results) != 1 || got.Results[0].Name != "a.b" {
		t.Fatalf("results %+v", got.Results)
	}
}

// Without a catalog a subquery's * leaves its columns unknown, and a name
// read from it reads that column of the tables the * read.
func TestResolveQuerySpanStarWithoutCatalog(t *testing.T) {
	for _, q := range []string{
		"SELECT email FROM (SELECT * FROM shop.customers)",
		"WITH c AS (SELECT * FROM shop.customers) SELECT email FROM c",
	} {
		got, err := ResolveQuerySpan(q, ResolveOptions{DefaultDatabase: "shop"})
		if err != nil {
			t.Fatal(err)
		}
		if got.NotFoundError != nil || len(got.Results) != 1 {
			t.Fatalf("%s: %+v", q, got)
		}
		if want := []string{"shop.customers.email"}; !slices.Equal(names(got.Results[0].SourceColumns), want) {
			t.Errorf("%s: reads %v, want %v", q, names(got.Results[0].SourceColumns), want)
		}
	}
}

// Aliases built on aliases resolve in linear time.
func TestResolveQuerySpanNestedAliases(t *testing.T) {
	var b strings.Builder
	b.WriteString("WITH email AS a0")
	for i := 1; i <= 200; i++ {
		fmt.Fprintf(&b, ", concat(a%d, a%d) AS a%d", i-1, i-1, i)
	}
	b.WriteString(" SELECT a200 FROM shop.customers")
	got, err := ResolveQuerySpan(b.String(), ResolveOptions{DefaultDatabase: "shop"})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"shop.customers.email"}; !slices.Equal(names(got.Results[0].SourceColumns), want) {
		t.Errorf("reads %v, want %v", names(got.Results[0].SourceColumns), want)
	}
}

// A deep expression is named without writing it out once per level.
func TestResolveQuerySpanDeepExpression(t *testing.T) {
	q := "SELECT 1" + strings.Repeat("+1", 20000)
	got, err := ResolveQuerySpan(q, ResolveOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Results) != 1 || got.Results[0].Name != q[len("SELECT "):] {
		t.Fatalf("results %d, name of %d bytes", len(got.Results), len(got.Results[0].Name))
	}
}

// A catalog that fails fails the span, also where x IN t asks for t.
func TestResolveQuerySpanCatalogError(t *testing.T) {
	broken := errors.New("catalog unavailable")
	_, err := ResolveQuerySpan("SELECT 1 IN allowed_ids", ResolveOptions{
		DefaultDatabase: "shop",
		Catalog:         Catalog{GetDatabase: func(string) (Database, error) { return nil, broken }},
	})
	if !errors.Is(err, broken) {
		t.Fatalf("err %v, want the catalog's", err)
	}
}

// Without a catalog, a * over a subquery that also names a column keeps the
// subquery's unexpanded * beside that column.
func TestResolveQuerySpanStarBesideColumnWithoutCatalog(t *testing.T) {
	got, err := ResolveQuerySpan("SELECT * FROM (SELECT *, email AS masked FROM shop.customers)", ResolveOptions{DefaultDatabase: "shop"})
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Results) != 2 || !got.Results[0].SelectAsterisk || got.Results[1].Name != "masked" {
		t.Fatalf("results %+v", got.Results)
	}
	if want := []string{"shop.customers"}; !slices.Equal(names(got.Results[0].SourceColumns), want) {
		t.Errorf("* reads %v, want %v", names(got.Results[0].SourceColumns), want)
	}
}

// A column computed in a subquery or by an alias is no plain field where
// it is read by name.
func TestResolveQuerySpanPlainField(t *testing.T) {
	for q, plain := range map[string]bool{
		"SELECT x FROM (SELECT lower(email) AS x FROM shop.customers)":      false,
		"WITH lower(email) AS x SELECT x FROM shop.customers":               false,
		"SELECT x FROM (SELECT email AS x FROM shop.customers)":             true,
		"SELECT email FROM shop.customers":                                  true,
		"SELECT x AS y FROM (SELECT lower(email) AS x FROM shop.customers)": false,
		"WITH lower(email) AS x, x AS y SELECT y FROM shop.customers":       false,
		"SELECT x AS y FROM (SELECT email AS x FROM shop.customers)":        true,
	} {
		got, err := ResolveQuerySpan(q, ResolveOptions{DefaultDatabase: "shop"})
		if err != nil {
			t.Fatal(err)
		}
		if got.Results[0].IsPlainField != plain {
			t.Errorf("%s: IsPlainField %v, want %v", q, got.Results[0].IsPlainField, plain)
		}
	}
}

// A recursive CTE that moves a value one column a pass converges however
// many columns it moves it through.
func TestResolveQuerySpanRecursiveConvergence(t *testing.T) {
	const n = 40
	var seed, step []string
	seed = append(seed, "email AS c0")
	step = append(step, "r.c0 AS c0")
	for i := 1; i <= n; i++ {
		seed = append(seed, fmt.Sprintf("'' AS c%d", i))
		step = append(step, fmt.Sprintf("r.c%d AS c%d", i-1, i))
	}
	q := fmt.Sprintf("WITH RECURSIVE r AS (SELECT %s FROM shop.customers UNION ALL SELECT %s FROM r) SELECT c%d FROM r",
		strings.Join(seed, ", "), strings.Join(step, ", "), n)
	got, err := ResolveQuerySpan(q, ResolveOptions{DefaultDatabase: "shop"})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"shop.customers.email"}; !slices.Equal(names(got.Results[0].SourceColumns), want) {
		t.Errorf("c%d reads %v, want %v", n, names(got.Results[0].SourceColumns), want)
	}
}

// x IN t.m where t is an alias's tuple reads the member, not a table t.m.
func TestResolveQuerySpanInAliasMember(t *testing.T) {
	got, err := ResolveQuerySpan("WITH CAST(tuple([email]), 'Tuple(vals Array(String))') AS t SELECT 'a' IN t.vals FROM shop.customers", ResolveOptions{DefaultDatabase: "shop"})
	if err != nil {
		t.Fatal(err)
	}
	if got.NotFoundError != nil {
		t.Fatalf("NotFoundError %v", got.NotFoundError)
	}
	if want := []string{"shop.customers.email"}; !slices.Equal(names(got.Results[0].SourceColumns), want) {
		t.Errorf("reads %v, want %v", names(got.Results[0].SourceColumns), want)
	}
}

// A chain of ALIAS columns, each reading the one before, resolves to the
// column at its root in linear time.
func TestResolveQuerySpanAliasColumnChain(t *testing.T) {
	const n = 60
	cols := []*metadata.ColumnMetadata{{Name: "email"}}
	for i := 1; i <= n; i++ {
		prev := "email"
		if i > 1 {
			prev = fmt.Sprintf("a%d", i-1)
		}
		cols = append(cols, &metadata.ColumnMetadata{Name: fmt.Sprintf("a%d", i), Generation: &metadata.GenerationMetadata{Type: metadata.GenerationMetadata_TYPE_VIRTUAL, Expression: "concat(" + prev + ", " + prev + ")"}})
	}
	db := tableDatabase{table: &metadata.TableMetadata{Name: "t", Columns: cols}}
	got, err := ResolveQuerySpan(fmt.Sprintf("SELECT a%d FROM t", n), ResolveOptions{
		DefaultDatabase: "shop",
		Catalog:         Catalog{GetDatabase: func(string) (Database, error) { return db, nil }},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !got.Results[0].SourceColumns[ColumnResource{Database: "shop", Table: "t", Column: "email"}] {
		t.Errorf("a%d reads %v, want shop.t.email among them", n, names(got.Results[0].SourceColumns))
	}
}

type tableDatabase struct{ table *metadata.TableMetadata }

func (d tableDatabase) GetTable(name string) *metadata.TableMetadata {
	if name == d.table.Name {
		return d.table
	}
	return nil
}
func (tableDatabase) GetView(string) *metadata.ViewMetadata { return nil }

// A self-named alias (email AS email reads the column) does not stop the
// aliases built on it from being resolved once.
func TestResolveQuerySpanSelfNamedAlias(t *testing.T) {
	var b strings.Builder
	b.WriteString("WITH concat(email, email) AS email, email AS a0")
	for i := 1; i <= 200; i++ {
		fmt.Fprintf(&b, ", concat(a%d, a%d) AS a%d", i-1, i-1, i)
	}
	b.WriteString(" SELECT a200 FROM shop.customers")
	got, err := ResolveQuerySpan(b.String(), ResolveOptions{DefaultDatabase: "shop"})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"shop.customers.email"}; !slices.Equal(names(got.Results[0].SourceColumns), want) {
		t.Errorf("reads %v, want %v", names(got.Results[0].SourceColumns), want)
	}
}
