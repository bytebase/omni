package completion

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/trino/catalog"
)

// buildCatalog constructs a small Trino catalog mirroring the oracle's shape:
// a "tpch" catalog with a "sf1" schema holding "orders" and "customer" tables
// and a "v_orders" view, and a "memory" catalog with a "default" schema. The
// session defaults to tpch.sf1 so unqualified names resolve there.
func buildCatalog() *catalog.Catalog {
	cat := catalog.New()

	tpch := cat.EnsureCatalog("tpch")
	sf1 := tpch.EnsureSchema("sf1")
	sf1.AddTable("orders",
		catalog.NewColumn("orderkey", "bigint", false),
		catalog.NewColumn("custkey", "bigint", false),
		catalog.NewColumn("orderstatus", "varchar(1)", true),
		catalog.NewColumn("totalprice", "double", true),
	)
	sf1.AddTable("customer",
		catalog.NewColumn("custkey", "bigint", false),
		catalog.NewColumn("name", "varchar(25)", true),
		catalog.NewColumn("nationkey", "bigint", true),
	)
	sf1.AddView("v_orders",
		catalog.NewColumn("orderkey", "bigint", false),
		catalog.NewColumn("custname", "varchar(25)", true),
	)

	mem := cat.EnsureCatalog("memory")
	mem.EnsureSchema("default").AddTable("kv",
		catalog.NewColumn("k", "varchar", true),
		catalog.NewColumn("v", "varchar", true),
	)

	cat.SetCurrentCatalog("tpch")
	cat.SetCurrentSchema("sf1")
	return cat
}

// texts extracts the Text of every candidate of the given type.
func texts(cands []Candidate, typ CandidateType) []string {
	var out []string
	for _, c := range cands {
		if c.Type == typ {
			out = append(out, c.Text)
		}
	}
	return out
}

// has reports whether any candidate of the given type has the given text.
func has(cands []Candidate, typ CandidateType, text string) bool {
	for _, c := range cands {
		if c.Type == typ && c.Text == text {
			return true
		}
	}
	return false
}

// hasText reports whether any candidate (any type) has the given text.
func hasText(cands []Candidate, text string) bool {
	for _, c := range cands {
		if c.Text == text {
			return true
		}
	}
	return false
}

func TestComplete_StatementStart(t *testing.T) {
	cat := buildCatalog()
	for _, tc := range []struct {
		name string
		sql  string
		pos  int
	}{
		{"empty", "", 0},
		{"after_semicolon", "SELECT 1; ", len("SELECT 1; ")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := Complete(tc.sql, tc.pos, cat)
			for _, kw := range []string{"SELECT", "INSERT", "CREATE", "WITH", "EXPLAIN"} {
				if !has(got, CandidateKeyword, kw) {
					t.Errorf("statement-start completion missing keyword %q; got keywords %v", kw, texts(got, CandidateKeyword))
				}
			}
			// No object candidates at a bare statement start.
			if cols := texts(got, CandidateColumn); len(cols) != 0 {
				t.Errorf("statement start offered columns %v, want none", cols)
			}
		})
	}
}

func TestComplete_StatementStartPrefix(t *testing.T) {
	got := Complete("SEL", 3, buildCatalog())
	if !has(got, CandidateKeyword, "SELECT") {
		t.Fatalf("Complete(\"SEL\") should offer SELECT; got %v", got)
	}
	for _, c := range got {
		if !strings.HasPrefix(strings.ToUpper(c.Text), "SEL") {
			t.Errorf("candidate %q does not match prefix SEL", c.Text)
		}
	}
}

func TestComplete_AfterFrom(t *testing.T) {
	cat := buildCatalog()
	sql := "SELECT * FROM "
	got := Complete(sql, len(sql), cat)

	// Tables and the view in the session schema.
	for _, name := range []string{"orders", "customer"} {
		if !has(got, CandidateTable, name) {
			t.Errorf("after FROM missing table %q; tables=%v", name, texts(got, CandidateTable))
		}
	}
	if !has(got, CandidateView, "v_orders") {
		t.Errorf("after FROM missing view v_orders; views=%v", texts(got, CandidateView))
	}
	// Catalog names (drill-down) and the session schema name.
	if !has(got, CandidateCatalog, "tpch") || !has(got, CandidateCatalog, "memory") {
		t.Errorf("after FROM missing catalog candidates; catalogs=%v", texts(got, CandidateCatalog))
	}
	if !has(got, CandidateSchema, "sf1") {
		t.Errorf("after FROM missing schema sf1; schemas=%v", texts(got, CandidateSchema))
	}
}

func TestComplete_AfterFromPrefix(t *testing.T) {
	cat := buildCatalog()
	sql := "SELECT * FROM cus"
	got := Complete(sql, len(sql), cat)
	if !has(got, CandidateTable, "customer") {
		t.Fatalf("FROM cus should complete customer; got %v", got)
	}
	if has(got, CandidateTable, "orders") {
		t.Errorf("FROM cus should not offer orders (prefix mismatch); got %v", texts(got, CandidateTable))
	}
}

func TestComplete_AfterJoin(t *testing.T) {
	cat := buildCatalog()
	sql := "SELECT * FROM orders o JOIN "
	got := Complete(sql, len(sql), cat)
	if !has(got, CandidateTable, "customer") {
		t.Errorf("after JOIN missing table customer; tables=%v", texts(got, CandidateTable))
	}
}

func TestComplete_SelectColumns(t *testing.T) {
	cat := buildCatalog()
	// Caret right after SELECT, FROM clause to the right.
	sql := "SELECT  FROM orders"
	got := Complete(sql, len("SELECT "), cat)
	for _, col := range []string{"orderkey", "custkey", "orderstatus", "totalprice"} {
		if !has(got, CandidateColumn, col) {
			t.Errorf("SELECT col context missing column %q; columns=%v", col, texts(got, CandidateColumn))
		}
	}
	// Expression keywords are also offered here.
	if !has(got, CandidateKeyword, "CASE") {
		t.Errorf("SELECT col context should offer expression keyword CASE; keywords=%v", texts(got, CandidateKeyword))
	}
}

func TestComplete_WhereColumns(t *testing.T) {
	cat := buildCatalog()
	sql := "SELECT * FROM customer WHERE "
	got := Complete(sql, len(sql), cat)
	for _, col := range []string{"custkey", "name", "nationkey"} {
		if !has(got, CandidateColumn, col) {
			t.Errorf("WHERE context missing column %q; columns=%v", col, texts(got, CandidateColumn))
		}
	}
	// Columns of an out-of-scope table must NOT appear.
	if has(got, CandidateColumn, "totalprice") {
		t.Errorf("WHERE on customer offered orders column totalprice; columns=%v", texts(got, CandidateColumn))
	}
}

func TestComplete_JoinOnColumns(t *testing.T) {
	cat := buildCatalog()
	sql := "SELECT * FROM orders o JOIN customer c ON "
	got := Complete(sql, len(sql), cat)
	// Both tables are in scope: columns from each appear.
	if !has(got, CandidateColumn, "orderkey") {
		t.Errorf("ON context missing orders column orderkey; columns=%v", texts(got, CandidateColumn))
	}
	if !has(got, CandidateColumn, "name") {
		t.Errorf("ON context missing customer column name; columns=%v", texts(got, CandidateColumn))
	}
}

func TestComplete_DottedAliasColumns(t *testing.T) {
	cat := buildCatalog()
	// Alias-qualified: o.<caret> must offer ONLY orders' columns. The select
	// list is complete (*) so the only incomplete clause is the predicate at the
	// caret, which the placeholder fills.
	sql := "SELECT * FROM orders o WHERE o."
	got := Complete(sql, len(sql), cat)
	if !has(got, CandidateColumn, "orderkey") || !has(got, CandidateColumn, "totalprice") {
		t.Errorf("o. should offer orders columns; columns=%v", texts(got, CandidateColumn))
	}
	if has(got, CandidateColumn, "name") {
		t.Errorf("o. (orders alias) should not offer customer column name; columns=%v", texts(got, CandidateColumn))
	}
}

func TestComplete_DottedTableNameColumns(t *testing.T) {
	cat := buildCatalog()
	// Table-name-qualified (no alias): orders.<caret>.
	sql := "SELECT * FROM orders WHERE orders."
	got := Complete(sql, len(sql), cat)
	if !has(got, CandidateColumn, "orderkey") {
		t.Errorf("orders. should offer orders columns; columns=%v", texts(got, CandidateColumn))
	}
}

func TestComplete_DottedCatalogSchemas(t *testing.T) {
	cat := buildCatalog()
	// "tpch." → schemas of the tpch catalog.
	sql := "SELECT * FROM tpch."
	got := Complete(sql, len(sql), cat)
	if !has(got, CandidateSchema, "sf1") {
		t.Errorf("tpch. should offer schema sf1; schemas=%v", texts(got, CandidateSchema))
	}
}

func TestComplete_DottedSchemaRelations(t *testing.T) {
	cat := buildCatalog()
	// "sf1." (schema in the session catalog) → relations of sf1.
	sql := "SELECT * FROM sf1."
	got := Complete(sql, len(sql), cat)
	if !has(got, CandidateTable, "orders") || !has(got, CandidateTable, "customer") {
		t.Errorf("sf1. should offer its tables; tables=%v", texts(got, CandidateTable))
	}
	if !has(got, CandidateView, "v_orders") {
		t.Errorf("sf1. should offer its view; views=%v", texts(got, CandidateView))
	}
}

func TestComplete_DottedCatalogSchemaRelations(t *testing.T) {
	cat := buildCatalog()
	// "tpch.sf1." → relations of tpch.sf1.
	sql := "SELECT * FROM tpch.sf1."
	got := Complete(sql, len(sql), cat)
	if !has(got, CandidateTable, "orders") {
		t.Errorf("tpch.sf1. should offer table orders; tables=%v", texts(got, CandidateTable))
	}
}

func TestComplete_DottedQualifiedColumns(t *testing.T) {
	cat := buildCatalog()
	// "tpch.sf1.orders." appearing in a WHERE → columns of orders.
	sql := "SELECT * FROM tpch.sf1.orders WHERE tpch.sf1.orders."
	got := Complete(sql, len(sql), cat)
	if !has(got, CandidateColumn, "orderkey") {
		t.Errorf("tpch.sf1.orders. should offer column orderkey; columns=%v", texts(got, CandidateColumn))
	}
}

func TestComplete_CTEAsRelation(t *testing.T) {
	cat := buildCatalog()
	// A CTE defined earlier is select-able after FROM.
	sql := "WITH recent AS (SELECT * FROM orders) SELECT * FROM "
	got := Complete(sql, len(sql), cat)
	if !has(got, CandidateTable, "recent") {
		t.Errorf("FROM after a WITH should offer the CTE name 'recent'; tables=%v", texts(got, CandidateTable))
	}
}

func TestComplete_NilCatalog(t *testing.T) {
	// With no catalog, keyword completion still works and nothing panics.
	got := Complete("SELECT * FROM ", len("SELECT * FROM "), nil)
	// CTE names would still be offered, but there are none here; the result may
	// be empty. The contract is "no panic, no object candidates".
	for _, c := range got {
		if c.Type == CandidateTable || c.Type == CandidateColumn || c.Type == CandidateSchema || c.Type == CandidateCatalog || c.Type == CandidateView {
			t.Errorf("nil catalog produced object candidate %v %q", c.Type, c.Text)
		}
	}
	// And statement-start keywords still work.
	if !has(Complete("", 0, nil), CandidateKeyword, "SELECT") {
		t.Errorf("nil catalog: statement-start keywords should still be offered")
	}
}

func TestComplete_NilCatalogCTE(t *testing.T) {
	// CTE names are derived from the statement text, so they are offered even
	// without a catalog.
	sql := "WITH recent AS (SELECT 1) SELECT * FROM "
	got := Complete(sql, len(sql), nil)
	if !has(got, CandidateTable, "recent") {
		t.Errorf("nil catalog should still offer CTE 'recent' from the statement; got %v", got)
	}
}

func TestComplete_OutOfRangeOffsets(t *testing.T) {
	cat := buildCatalog()
	cases := []struct {
		name string
		sql  string
		pos  int
	}{
		{"far_past_end", "SELECT * FROM ", 1000},
		{"one_past_end", "SELECT * FROM ", len("SELECT * FROM ") + 1},
		{"negative", "SELECT", -5},
		{"huge", "SELECT", 1 << 30},
		{"empty_past_end", "", 9},
		{"empty_negative", "", -1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("Complete(%q, %d) panicked: %v", tc.sql, tc.pos, r)
				}
			}()
			_ = Complete(tc.sql, tc.pos, cat)
		})
	}

	// A caret past the end of "SELECT * FROM " behaves like a caret at the end.
	atEnd := Complete("SELECT * FROM ", len("SELECT * FROM "), cat)
	pastEnd := Complete("SELECT * FROM ", 1000, cat)
	if len(atEnd) != len(pastEnd) {
		t.Errorf("past-end completion (%d cands) differs from at-end (%d cands)", len(pastEnd), len(atEnd))
	}
}

func TestComplete_ResultsSortedAndDeduped(t *testing.T) {
	cat := buildCatalog()
	got := Complete("SELECT * FROM ", len("SELECT * FROM "), cat)
	// Sorted by (Type, Text): verify non-decreasing.
	for i := 1; i < len(got); i++ {
		prev, cur := got[i-1], got[i]
		if cur.Type < prev.Type {
			t.Errorf("candidates not sorted by type at %d: %v then %v", i, prev, cur)
		}
		if cur.Type == prev.Type && cur.Text < prev.Text {
			t.Errorf("candidates not sorted by text within type at %d: %q then %q", i, prev.Text, cur.Text)
		}
	}
	// No duplicate (Text, Type).
	seen := map[string]bool{}
	for _, c := range got {
		k := c.Type.String() + "\x00" + strings.ToLower(c.Text)
		if seen[k] {
			t.Errorf("duplicate candidate %v %q", c.Type, c.Text)
		}
		seen[k] = true
	}
}

func TestComplete_MultiStatementCaretInSecond(t *testing.T) {
	cat := buildCatalog()
	// The caret is in the second statement; the first must not pollute scope.
	sql := "SELECT * FROM customer; SELECT * FROM "
	got := Complete(sql, len(sql), cat)
	// FROM in the second statement still offers tables (context independent of
	// statement index), and the first statement's customer scope is irrelevant.
	if !has(got, CandidateTable, "orders") {
		t.Errorf("second-statement FROM should offer tables; tables=%v", texts(got, CandidateTable))
	}
}

func TestQuoteIdentifierIfNeeded(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"orders", "orders"},           // clean lower-case
		{"MyTable", `"MyTable"`},       // mixed case (from a quoted source) → re-quote to preserve
		{"with space", `"with space"`}, // space
		{"1leading", `"1leading"`},     // leading digit
		{"weird-name", `"weird-name"`}, // hyphen
		{"", `""`},                     // empty
		{"select", `"select"`},         // reserved keyword
		{"data", "data"},               // non-reserved keyword → no quoting
		{`a"b`, `"a""b"`},              // embedded quote → doubled
		{"_underscore", "_underscore"}, // leading underscore is fine
		{"col1", "col1"},               // trailing digit fine
	}
	for _, tc := range cases {
		if got := QuoteIdentifierIfNeeded(tc.in); got != tc.want {
			t.Errorf("QuoteIdentifierIfNeeded(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestComplete_QuotedCandidateForReservedName(t *testing.T) {
	cat := catalog.New()
	tpch := cat.EnsureCatalog("tpch")
	// A table whose normalized name is a reserved keyword must be offered quoted
	// so the completion is itself valid syntax.
	tpch.EnsureSchema("sf1").AddTable("select",
		catalog.NewColumn("x", "bigint", false),
	)
	cat.SetCurrentCatalog("tpch")
	cat.SetCurrentSchema("sf1")

	got := Complete("SELECT * FROM ", len("SELECT * FROM "), cat)
	if !hasText(got, `"select"`) {
		t.Errorf("a table named 'select' must be offered quoted; got %v", texts(got, CandidateTable))
	}
}
