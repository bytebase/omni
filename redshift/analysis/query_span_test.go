package analysis

import (
	"testing"

	"github.com/bytebase/omni/redshift/catalog"
)

func TestQuerySpanSimpleColumnsAndStar(t *testing.T) {
	c := catalog.New()
	if _, err := c.Exec(`CREATE TABLE users (id int, name text, email text, amount int);`, nil); err != nil {
		t.Fatal(err)
	}

	span := mustQuerySpan(t, c, `SELECT id, users.name, amount AS total, users.* FROM users`)
	assertResult(t, span.Results[0], "id", true, []ColumnResource{{Schema: "public", Table: "users", Column: "id"}})
	assertResult(t, span.Results[1], "name", true, []ColumnResource{{Schema: "public", Table: "users", Column: "name"}})
	assertResult(t, span.Results[2], "total", true, []ColumnResource{{Schema: "public", Table: "users", Column: "amount"}})
	assertResult(t, span.Results[3], "id", true, []ColumnResource{{Schema: "public", Table: "users", Column: "id"}})
	assertResult(t, span.Results[4], "name", true, []ColumnResource{{Schema: "public", Table: "users", Column: "name"}})
	assertResult(t, span.Results[5], "email", true, []ColumnResource{{Schema: "public", Table: "users", Column: "email"}})
}

func TestQuerySpanExpressionsAndPredicates(t *testing.T) {
	c := catalog.New()
	if _, err := c.Exec(`
CREATE TABLE users (id int, name text);
CREATE TABLE payments (id int, user_id int, amount numeric);
`, nil); err != nil {
		t.Fatal(err)
	}

	span := mustQuerySpan(t, c, `
SELECT users.name, payments.amount + 1 AS adjusted
FROM users
JOIN payments ON users.id = payments.user_id
WHERE payments.amount > 10
`)
	assertResult(t, span.Results[0], "name", true, []ColumnResource{{Schema: "public", Table: "users", Column: "name"}})
	assertResult(t, span.Results[1], "adjusted", false, []ColumnResource{{Schema: "public", Table: "payments", Column: "amount"}})
	assertColumnsContain(t, span.SourceColumns,
		ColumnResource{Schema: "public", Table: "users", Column: "name"},
		ColumnResource{Schema: "public", Table: "payments", Column: "amount"},
		ColumnResource{Schema: "public", Table: "users", Column: "id"},
		ColumnResource{Schema: "public", Table: "payments", Column: "user_id"},
	)
	assertColumnsEqual(t, span.PredicateColumns,
		ColumnResource{Schema: "public", Table: "payments", Column: "amount"},
		ColumnResource{Schema: "public", Table: "payments", Column: "user_id"},
		ColumnResource{Schema: "public", Table: "users", Column: "id"},
	)
}

func TestQuerySpanJoinUsing(t *testing.T) {
	c := catalog.New()
	if _, err := c.Exec(`
CREATE TABLE left_t (id int, name text);
CREATE TABLE right_t (id int, status text);
`, nil); err != nil {
		t.Fatal(err)
	}

	span := mustQuerySpan(t, c, `SELECT * FROM left_t JOIN right_t USING(id)`)
	assertColumnsContain(t, span.SourceColumns,
		ColumnResource{Schema: "public", Table: "left_t", Column: "id"},
		ColumnResource{Schema: "public", Table: "right_t", Column: "id"},
	)
	assertColumnsContain(t, span.PredicateColumns,
		ColumnResource{Schema: "public", Table: "left_t", Column: "id"},
		ColumnResource{Schema: "public", Table: "right_t", Column: "id"},
	)
}

func TestQuerySpanCTESubqueryAndSetOperation(t *testing.T) {
	c := catalog.New()
	if _, err := c.Exec(`
CREATE TABLE t1 (a int, b text);
CREATE TABLE t2 (c int, d text);
`, nil); err != nil {
		t.Fatal(err)
	}

	cteSpan := mustQuerySpan(t, c, `WITH q(x, y) AS (SELECT a, b FROM t1) SELECT y FROM q`)
	assertResult(t, cteSpan.Results[0], "y", true, []ColumnResource{{Schema: "public", Table: "t1", Column: "b"}})

	implicitCTESpan := mustQuerySpan(t, c, `WITH q AS (SELECT a, b FROM t1) SELECT b FROM q`)
	assertResult(t, implicitCTESpan.Results[0], "b", true, []ColumnResource{{Schema: "public", Table: "t1", Column: "b"}})

	subquerySpan := mustQuerySpan(t, c, `SELECT x FROM (SELECT a AS x FROM t1) s`)
	assertResult(t, subquerySpan.Results[0], "x", true, []ColumnResource{{Schema: "public", Table: "t1", Column: "a"}})

	setSpan := mustQuerySpan(t, c, `SELECT a FROM t1 UNION ALL SELECT c FROM t2`)
	assertResult(t, setSpan.Results[0], "a", true, []ColumnResource{
		{Schema: "public", Table: "t1", Column: "a"},
		{Schema: "public", Table: "t2", Column: "c"},
	})
}

func TestQuerySpanViewDefinitionLineage(t *testing.T) {
	c := catalog.New()
	if _, err := c.Exec(`
CREATE TABLE base (a int, b text);
CREATE VIEW v AS SELECT a, b FROM base;
`, nil); err != nil {
		t.Fatal(err)
	}

	span := mustQuerySpan(t, c, `SELECT * FROM v`)
	assertResult(t, span.Results[0], "a", true, []ColumnResource{{Schema: "public", Table: "base", Column: "a"}})
	assertResult(t, span.Results[1], "b", true, []ColumnResource{{Schema: "public", Table: "base", Column: "b"}})
}

func TestQuerySpanRedshiftSelectExtensions(t *testing.T) {
	c := catalog.New()
	if _, err := c.Exec(`CREATE TABLE users (id int, password text, amount int);`, nil); err != nil {
		t.Fatal(err)
	}

	excludeSpan := mustQuerySpan(t, c, `SELECT * EXCLUDE (password) FROM users`)
	if len(excludeSpan.Results) != 2 {
		t.Fatalf("EXCLUDE result count = %d, want 2: %#v", len(excludeSpan.Results), excludeSpan.Results)
	}
	assertResult(t, excludeSpan.Results[0], "id", true, []ColumnResource{{Schema: "public", Table: "users", Column: "id"}})
	assertResult(t, excludeSpan.Results[1], "amount", true, []ColumnResource{{Schema: "public", Table: "users", Column: "amount"}})

	qualifySpan := mustQuerySpan(t, c, `SELECT id FROM users QUALIFY amount > 10`)
	assertColumnsEqual(t, qualifySpan.PredicateColumns,
		ColumnResource{Schema: "public", Table: "users", Column: "amount"},
	)
}

func TestQuerySpanRedshiftTopAndMinus(t *testing.T) {
	c := catalog.New()
	if _, err := c.Exec(`
CREATE TABLE all_products (product_id int, name text);
CREATE TABLE discontinued_products (product_id int);
`, nil); err != nil {
		t.Fatal(err)
	}

	topSpan := mustQuerySpan(t, c, `SELECT TOP 10 product_id FROM all_products ORDER BY product_id`)
	assertResult(t, topSpan.Results[0], "product_id", true, []ColumnResource{{Schema: "public", Table: "all_products", Column: "product_id"}})

	minusSpan := mustQuerySpan(t, c, `SELECT product_id FROM all_products MINUS SELECT product_id FROM discontinued_products`)
	assertResult(t, minusSpan.Results[0], "product_id", true, []ColumnResource{
		{Schema: "public", Table: "all_products", Column: "product_id"},
		{Schema: "public", Table: "discontinued_products", Column: "product_id"},
	})
}

func TestQuerySpanRejectsMutatingOrUnsupportedRedshiftSelect(t *testing.T) {
	c := catalog.New()
	if _, err := c.Exec(`CREATE TABLE users (id int, manager_id int);`, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := GetQuerySpan(c, `SELECT id INTO new_users FROM users`); err == nil {
		t.Fatalf("expected SELECT INTO query span to be rejected")
	}
	if _, err := GetQuerySpan(c, `SELECT id FROM users CONNECT BY PRIOR id = manager_id`); err == nil {
		t.Fatalf("expected CONNECT BY query span to be rejected")
	}
}

func mustQuerySpan(t *testing.T, c *catalog.Catalog, sql string) *QuerySpan {
	t.Helper()
	span, err := GetQuerySpan(c, sql)
	if err != nil {
		t.Fatalf("GetQuerySpan(%q) returned error: %v", sql, err)
	}
	return span
}

func assertResult(t *testing.T, got QuerySpanResult, name string, plain bool, sources []ColumnResource) {
	t.Helper()
	if got.Name != name {
		t.Fatalf("result name = %q, want %q", got.Name, name)
	}
	if got.IsPlainField != plain {
		t.Fatalf("result %q plain = %v, want %v", name, got.IsPlainField, plain)
	}
	assertColumnsEqual(t, got.SourceColumns, sources...)
}

func assertColumnsContain(t *testing.T, got []ColumnResource, want ...ColumnResource) {
	t.Helper()
	for _, w := range want {
		found := false
		for _, g := range got {
			if g == w {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("missing column %#v in %#v", w, got)
		}
	}
}

func assertColumnsEqual(t *testing.T, got []ColumnResource, want ...ColumnResource) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("columns = %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("columns[%d] = %#v, want %#v; all = %#v", i, got[i], want[i], got)
		}
	}
}
