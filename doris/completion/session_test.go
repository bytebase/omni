package completion

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
)

func col(name, typ string) *metadatapb.ColumnMetadata {
	return &metadatapb.ColumnMetadata{Name: name, Type: typ}
}

// testDatabases are the fixture databases: db is the session's, otherdb one
// a qualified name reaches.
func testDatabases() map[string]*metadatapb.DatabaseSchemaMetadata {
	return map[string]*metadatapb.DatabaseSchemaMetadata{
		"db": {Name: "db", Schemas: []*metadatapb.SchemaMetadata{{
			Tables: []*metadatapb.TableMetadata{
				{Name: "orders", Columns: []*metadatapb.ColumnMetadata{col("id", "BIGINT"), col("customer_id", "BIGINT"), col("amount", "DECIMAL(10,2)")}},
				{Name: "customers", Columns: []*metadatapb.ColumnMetadata{col("id", "BIGINT"), col("name", "VARCHAR(64)")}},
				{Name: "order items", Columns: []*metadatapb.ColumnMetadata{col("qty", "INT"), col("select", "INT")}},
			},
			Views:             []*metadatapb.ViewMetadata{{Name: "v_orders", Columns: []*metadatapb.ColumnMetadata{col("id", "BIGINT")}}},
			MaterializedViews: []*metadatapb.MaterializedViewMetadata{{Name: "mv_totals", Columns: []*metadatapb.ColumnMetadata{col("total", "DECIMAL(10,2)")}}},
		}}},
		"otherdb": {Name: "otherdb", Schemas: []*metadatapb.SchemaMetadata{{
			Tables: []*metadatapb.TableMetadata{{Name: "other_t", Columns: []*metadatapb.ColumnMetadata{col("other_col", "INT")}}},
		}}},
	}
}

func testConfig() SessionConfig {
	dbs := testDatabases()
	return SessionConfig{
		DefaultDatabase: "db",
		LoadDatabase: func(_ context.Context, name string) (*metadatapb.DatabaseSchemaMetadata, error) {
			if db, ok := dbs[name]; ok {
				return db, nil
			}
			return nil, fmt.Errorf("database %q not found", name)
		},
		ListDatabaseNames: func(context.Context) ([]string, error) {
			return []string{"db", "otherdb"}, nil
		},
	}
}

// completeAt completes text, whose caret is the first |.
func completeAt(t *testing.T, cfg SessionConfig, text string) []Candidate {
	t.Helper()
	i := strings.Index(text, "|")
	if i < 0 {
		t.Fatalf("no caret in %q", text)
	}
	stmt := text[:i] + text[i+1:]
	line := 1 + strings.Count(stmt[:i], "\n")
	col := i - (strings.LastIndex(stmt[:i], "\n") + 1)
	return CompleteSession(context.Background(), cfg, stmt, line, col)
}

func labels(got []Candidate) []string {
	out := make([]string, len(got))
	for i, c := range got {
		out[i] = fmt.Sprintf("%v:%s", c.Type, c.Text)
	}
	return out
}

func kw(words ...string) []string {
	out := make([]string, len(words))
	for i, w := range words {
		out[i] = fmt.Sprintf("%v:%s", CandidateKeyword, w)
	}
	return out
}

func typed(typ CandidateType, names ...string) []string {
	out := make([]string, len(names))
	for i, n := range names {
		out[i] = fmt.Sprintf("%v:%s", typ, n)
	}
	return out
}

func join(parts ...[]string) []string {
	return slices.Concat(parts...)
}

var (
	relationsOfDB = join(
		typed(CandidateDatabase, "db", "otherdb"),
		typed(CandidateTable, "`order items`", "customers", "orders"),
		typed(CandidateView, "v_orders"),
		typed(CandidateMaterializedView, "mv_totals"),
	)
	ordersColumns = typed(CandidateColumn, "amount", "customer_id", "id")
)

func TestCompleteSession(t *testing.T) {
	tests := []struct {
		name string
		text string
		want []string // the whole result, sorted by type and text
		has  []string // a subset the result holds
		not  []string // labels the result must not hold
	}{
		{name: "statement start", text: "SEL|", want: kw("SELECT")},
		{name: "statement start keywords", text: "|", has: kw("ALTER", "CREATE", "SELECT", "SHOW", "WITH"), not: kw("FROM", "JOIN")},
		{name: "keywords after a statement", text: "SELECT 1; |", has: kw("SELECT", "INSERT")},
		{name: "insert", text: "INSERT |", want: kw("INTO", "OVERWRITE")},
		{name: "create table if", text: "CREATE TABLE |", want: kw("IF")},
		{name: "group by", text: "SELECT a FROM orders GROUP |", want: kw("BY")},
		{name: "join kinds", text: "SELECT * FROM orders LEFT |", want: kw("ANTI", "JOIN", "OUTER", "SEMI")},
		{name: "lookahead keyword", text: "SELECT * FROM orders o |", has: kw("CROSS", "JOIN", "LEFT", "WHERE"), not: kw("SET", "VALUES")},
		{name: "doris table clauses", text: "SELECT * FROM orders |", has: kw("PARTITION", "TABLET", "TABLESAMPLE")},
		{name: "after a qualifier no keywords", text: "SELECT * FROM db.|", want: join(typed(CandidateTable, "`order items`", "customers", "orders"), typed(CandidateView, "v_orders"), typed(CandidateMaterializedView, "mv_totals"))},
		{name: "doris create table tail", text: "CREATE TABLE t (a INT) DUPLICATE KEY(a) |", has: kw("DISTRIBUTED", "PARTITION", "PROPERTIES"), not: kw("FROM", "WHERE")},

		{name: "from", text: "SELECT * FROM |", has: relationsOfDB, not: typed(CandidateColumn, "id")},
		{name: "from typed", text: "SELECT * FROM ord|", want: typed(CandidateTable, "`order items`", "orders")},
		{name: "join", text: "SELECT * FROM orders JOIN cu|", want: typed(CandidateTable, "customers")},
		{name: "comma join", text: "SELECT * FROM orders, cu|", want: typed(CandidateTable, "customers")},
		{name: "other database", text: "SELECT * FROM otherdb.|", want: typed(CandidateTable, "other_t")},
		{name: "insert into", text: "INSERT INTO ord|", want: typed(CandidateTable, "`order items`", "orders")},
		{name: "update", text: "UPDATE ord|", want: typed(CandidateTable, "`order items`", "orders")},
		{name: "delete", text: "DELETE FROM ord|", want: typed(CandidateTable, "`order items`", "orders")},
		{name: "describe", text: "DESC ord|", want: typed(CandidateTable, "`order items`", "orders")},
		{name: "no tables for a new table", text: "CREATE TABLE ord|", want: nil},
		{name: "create table as", text: "CREATE TABLE t AS SELECT * FROM ord|", want: typed(CandidateTable, "`order items`", "orders")},
		{name: "cte as table", text: "WITH recent AS (SELECT 1 AS x) SELECT * FROM rec|", want: typed(CandidateTable, "recent")},

		{name: "select list", text: "SELECT | FROM orders", has: ordersColumns, not: typed(CandidateColumn, "name")},
		{name: "where", text: "SELECT * FROM orders WHERE am|", want: typed(CandidateColumn, "amount")},
		{name: "two tables", text: "SELECT * FROM orders o JOIN customers c ON o.customer_id = c.id WHERE na|", want: typed(CandidateColumn, "name")},
		{name: "alias qualifier", text: "SELECT c.| FROM orders o JOIN customers c ON o.customer_id = c.id", want: typed(CandidateColumn, "id", "name")},
		{name: "table qualifier", text: "SELECT orders.| FROM orders", want: ordersColumns},
		{name: "alias hides table name", text: "SELECT orders.| FROM orders o", want: nil},
		{name: "database qualified table", text: "SELECT otherdb.other_t.| FROM otherdb.other_t", want: typed(CandidateColumn, "other_col")},
		{name: "view columns", text: "SELECT | FROM v_orders", has: typed(CandidateColumn, "id")},
		{name: "materialized view columns", text: "SELECT | FROM mv_totals", has: typed(CandidateColumn, "total")},
		{name: "incomplete where", text: "SELECT * FROM orders WHERE amount > 1 AND cust|", want: typed(CandidateColumn, "customer_id")},
		{name: "join on", text: "SELECT * FROM orders o LEFT JOIN customers c ON c.|", want: typed(CandidateColumn, "id", "name")},
		{name: "update set", text: "UPDATE orders SET am|", want: typed(CandidateColumn, "amount")},
		{name: "delete where", text: "DELETE FROM orders WHERE am|", want: typed(CandidateColumn, "amount")},
		{name: "order by", text: "SELECT * FROM orders ORDER BY am|", want: typed(CandidateColumn, "amount")},
		{name: "subquery", text: "SELECT * FROM customers WHERE id IN (SELECT cust| FROM orders)", want: typed(CandidateColumn, "customer_id")},
		{name: "cte body out of scope", text: "WITH recent AS (SELECT id FROM orders) SELECT am| FROM customers", want: nil},
		{name: "other union arm out of scope", text: "SELECT id FROM orders UNION SELECT am| FROM customers", want: nil},
		{name: "union arm in scope", text: "SELECT id FROM orders UNION SELECT na| FROM customers", want: typed(CandidateColumn, "name")},
		{name: "unfinished nested expression", text: "SELECT * FROM orders WHERE ((am|", want: typed(CandidateColumn, "amount")},
		{name: "catalog-qualified table", text: "SELECT o.am| FROM internal.db.orders o", want: typed(CandidateColumn, "amount")},
		{name: "unfinished clause after the caret", text: "SELECT am| FROM orders WHERE", want: typed(CandidateColumn, "amount")},
		{name: "create table as select list", text: "CREATE TABLE t AS SELECT am| FROM orders", want: typed(CandidateColumn, "amount")},
		{name: "create table as subquery", text: "CREATE TABLE t AS SELECT * FROM orders WHERE id IN (SELECT na| FROM customers)", want: typed(CandidateColumn, "name")},
		{name: "lexical scope after an unaliased table", text: "SELECT * FROM orders JOIN customers ON ((na|", want: typed(CandidateColumn, "name")},
		{name: "correlated outer alias", text: "SELECT * FROM orders o WHERE EXISTS (SELECT 1 FROM customers c WHERE o.am|)", want: typed(CandidateColumn, "amount")},
		{name: "inner alias first", text: "SELECT * FROM orders c WHERE EXISTS (SELECT 1 FROM customers c WHERE c.na|)", want: typed(CandidateColumn, "name")},
		{name: "outer cte in subquery", text: "WITH recent AS (SELECT 1 AS x) SELECT * FROM orders WHERE id IN (SELECT x FROM rec|)", want: typed(CandidateTable, "recent")},
		{name: "spaced qualifier", text: "SELECT * FROM otherdb . oth|", want: typed(CandidateTable, "other_t")},
		{name: "insert columns", text: "INSERT INTO orders (am|) VALUES (1)", want: typed(CandidateColumn, "amount")},
		{name: "insert columns second", text: "INSERT INTO orders (id, cu|", want: typed(CandidateColumn, "customer_id")},
		{name: "catalog-qualified relation", text: "SELECT * FROM internal.db.ord|", want: typed(CandidateTable, "`order items`", "orders")},
		{name: "catalog-qualified column", text: "SELECT internal.db.orders.am| FROM internal.db.orders", want: typed(CandidateColumn, "amount")},
		{name: "lexical comma list", text: "SELECT * FROM orders, customers WHERE ((na|", want: typed(CandidateColumn, "name")},
		{name: "nested cte not visible", text: "WITH recent AS (WITH orders AS (SELECT 1) SELECT * FROM orders) SELECT am| FROM orders", want: typed(CandidateColumn, "amount")},
		{name: "cte of own block visible", text: "WITH recent AS (WITH orders AS (SELECT 1 AS x) SELECT am| FROM orders) SELECT 1", want: nil},
		{name: "join using", text: "SELECT * FROM orders JOIN customers USING (i|)", want: typed(CandidateColumn, "id")},
		{name: "join using shared only", text: "SELECT * FROM orders JOIN customers USING (cu|)", want: nil},
		{name: "drop table if exists", text: "DROP TABLE IF EXISTS ord|", want: typed(CandidateTable, "`order items`", "orders")},
		{name: "alias in order by", text: "SELECT amount AS total FROM orders ORDER BY tot|", want: typed(CandidateColumn, "total")},
		{name: "alias in group by", text: "SELECT amount AS total FROM orders GROUP BY tot|", want: typed(CandidateColumn, "total")},
		{name: "alias not in where", text: "SELECT amount AS total FROM orders WHERE tot|", want: nil},
		{name: "fallback reads after the caret", text: "SELECT ((am| FROM orders", want: typed(CandidateColumn, "amount")},
		{name: "cte name case-insensitive", text: "WITH Orders AS (SELECT 1 AS x) SELECT am| FROM orders", want: nil},
		{name: "unclosed subquery keywords", text: "SELECT * FROM (SELECT id |", has: kw("FROM"), not: kw("ON")},
		{name: "derived table", text: "SELECT * FROM (SELECT * FROM cu|", want: typed(CandidateTable, "customers")},

		{name: "quoted name", text: "SELECT * FROM `ord|", want: typed(CandidateTable, "order items", "orders")},
		{name: "reserved column quoted", text: "SELECT sel| FROM `order items`", want: typed(CandidateColumn, "`select`")},
		{name: "string", text: "SELECT 'ord|", want: nil},
		{name: "comment", text: "SELECT /* ord| */ 1", want: nil},
		{name: "line comment", text: "SELECT 1 -- ord|", want: nil},
	}
	cfg := testConfig()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := labels(completeAt(t, cfg, tc.text))
			if tc.has == nil && tc.not == nil {
				if !slices.Equal(got, tc.want) {
					t.Errorf("got  %v\nwant %v", got, tc.want)
				}
				return
			}
			for _, w := range tc.has {
				if !slices.Contains(got, w) {
					t.Errorf("missing %s in %v", w, got)
				}
			}
			for _, w := range tc.not {
				if slices.Contains(got, w) {
					t.Errorf("unexpected %s in %v", w, got)
				}
			}
		})
	}
}

// TestCompleteSessionWithoutMetadata offers keywords and no names.
func TestCompleteSessionWithoutMetadata(t *testing.T) {
	got := labels(completeAt(t, SessionConfig{}, "SELECT * FROM orders WH|"))
	if want := kw("WHERE"); !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
	failing := SessionConfig{
		DefaultDatabase: "db",
		LoadDatabase: func(context.Context, string) (*metadatapb.DatabaseSchemaMetadata, error) {
			return nil, fmt.Errorf("unavailable")
		},
	}
	if got := completeAt(t, failing, "SELECT * FROM |"); slices.ContainsFunc(got, func(c Candidate) bool { return c.Type != CandidateKeyword }) {
		t.Errorf("names without metadata: %v", labels(got))
	}
}

// TestCompleteSessionLoadsFewDatabases bounds the databases one completion
// loads for qualifiers.
func TestCompleteSessionLoadsFewDatabases(t *testing.T) {
	cfg := testConfig()
	load := cfg.LoadDatabase
	var loaded []string
	cfg.LoadDatabase = func(ctx context.Context, name string) (*metadatapb.DatabaseSchemaMetadata, error) {
		loaded = append(loaded, name)
		return load(ctx, name)
	}
	completeAt(t, cfg, "SELECT a.x.|, b.x.y, c.x.y, d.x.y FROM a.x, b.x, c.x, d.x")
	if len(loaded) > 1+maxDatabaseLoads {
		t.Errorf("loaded %v", loaded)
	}
}

func TestQuoteIdentifierIfNeeded(t *testing.T) {
	for _, tc := range []struct {
		name, inQuotes, bare string
	}{
		{"orders", "orders", "orders"},
		{"Orders_1$", "Orders_1$", "Orders_1$"},
		{"order items", "order items", "`order items`"},
		{"select", "select", "`select`"},
		{"comment", "comment", "comment"}, // non-reserved
		{"1abc", "1abc", "`1abc`"},
		{"a`b", "a``b", "`a``b`"},
	} {
		if got := QuoteIdentifierIfNeeded(tc.name, true); got != tc.inQuotes {
			t.Errorf("QuoteIdentifierIfNeeded(%q, true) = %q, want %q", tc.name, got, tc.inQuotes)
		}
		if got := QuoteIdentifierIfNeeded(tc.name, false); got != tc.bare {
			t.Errorf("QuoteIdentifierIfNeeded(%q, false) = %q, want %q", tc.name, got, tc.bare)
		}
	}
}
