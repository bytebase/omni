package completion

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	metadatapb "github.com/bytebase/omni/metadata"
)

func col(name, typ string) *metadatapb.ColumnMetadata {
	return &metadatapb.ColumnMetadata{Name: name, Type: typ}
}

// testDatabases are the fixture databases: DB is the session's, OTHERDB one
// a three-part name reaches.
func testDatabases() map[string]*metadatapb.DatabaseSchemaMetadata {
	return map[string]*metadatapb.DatabaseSchemaMetadata{
		"DB": {Name: "DB", Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "PUBLIC",
				Tables: []*metadatapb.TableMetadata{
					{Name: "ORDERS", Columns: []*metadatapb.ColumnMetadata{col("ID", "NUMBER"), col("CUSTOMER_ID", "NUMBER"), col("SRC", "VARIANT")}},
					{Name: "CUSTOMERS", Columns: []*metadatapb.ColumnMetadata{col("ID", "NUMBER"), col("NAME", "VARCHAR")}},
					{Name: "lower_t", Columns: []*metadatapb.ColumnMetadata{col("mixedCol", "NUMBER"), col("ORDER", "NUMBER")}},
					{Name: "Order Items", Columns: []*metadatapb.ColumnMetadata{col("QTY", "NUMBER")}},
				},
				ExternalTables:    []*metadatapb.ExternalTableMetadata{{Name: "EXT_EVENTS", Columns: []*metadatapb.ColumnMetadata{col("PAYLOAD", "VARIANT")}}},
				Views:             []*metadatapb.ViewMetadata{{Name: "V_ORDERS", Columns: []*metadatapb.ColumnMetadata{col("ID", "NUMBER")}}},
				MaterializedViews: []*metadatapb.MaterializedViewMetadata{{Name: "MV_TOTALS"}},
				Functions:         []*metadatapb.FunctionMetadata{{Name: "MY_UDF", Signature: "MY_UDF(X NUMBER)"}},
			},
			{
				Name: "SALES",
				Tables: []*metadatapb.TableMetadata{
					{Name: "DEALS", Columns: []*metadatapb.ColumnMetadata{col("AMOUNT", "NUMBER")}},
					{Name: "CUSTOMERS", Columns: []*metadatapb.ColumnMetadata{col("SALES_REP", "VARCHAR")}},
				},
				Functions: []*metadatapb.FunctionMetadata{{Name: "SALES_UDF"}},
			},
		}},
		"OTHERDB": {Name: "OTHERDB", Schemas: []*metadatapb.SchemaMetadata{
			{Name: "PUBLIC", Tables: []*metadatapb.TableMetadata{{Name: "OTHER_T", Columns: []*metadatapb.ColumnMetadata{col("OTHER_COL", "NUMBER")}}}},
		}},
	}
}

func testConfig() SessionConfig {
	dbs := testDatabases()
	return SessionConfig{
		DefaultDatabase: "DB",
		LoadDatabase: func(_ context.Context, name string) (*metadatapb.DatabaseSchemaMetadata, error) {
			if db, ok := dbs[name]; ok {
				return db, nil
			}
			return nil, fmt.Errorf("database %q not found", name)
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

func TestCompleteSession(t *testing.T) {
	tests := []struct {
		name   string
		text   string
		want   []string
		absent []string
	}{
		// Keywords come from Snowflake's grammar.
		{"statement start", "|", []string{"keyword:SELECT", "keyword:UNDROP", "keyword:EXECUTE"}, []string{"keyword:FROM"}},
		{"clauses after a FROM item", "SELECT * FROM orders |", []string{"keyword:QUALIFY", "keyword:MATCH_RECOGNIZE", "keyword:ASOF", "keyword:SAMPLE", "keyword:WHERE"}, []string{"keyword:SELECT"}},
		{"typed keyword", "SELECT * FROM orders QUA|", []string{"keyword:QUALIFY"}, []string{"keyword:WHERE"}},
		{"lower-case typed keyword", "sel|", []string{"keyword:SELECT"}, nil},
		{"after ORDER", "SELECT * FROM orders ORDER |", []string{"keyword:BY"}, []string{"keyword:WHERE"}},
		{"a word the parser matches by text", "CREATE NETWORK |", []string{"keyword:RULE", "keyword:POLICY"}, nil},
		{"data types", "SELECT CAST(id AS |) FROM orders", []string{"keyword:NUMBER", "keyword:VARCHAR", "keyword:VARIANT"}, nil},
		{"no keyword after a qualifier", "SELECT o.| FROM orders o", []string{"column:ID"}, []string{"keyword:FROM"}},

		// Relations: the session's schema bare, others qualified.
		{"relations", "SELECT * FROM |", []string{"table:ORDERS", "table:EXT_EVENTS", "view:V_ORDERS", "materialized_view:MV_TOTALS", "schema:SALES", "schema:PUBLIC"}, []string{"table:DEALS", "column:ID"}},
		{"names needing quotes", "SELECT * FROM |", []string{`table:"lower_t"`, `table:"Order Items"`}, nil},
		{"typed relation", "SELECT * FROM ord|", []string{"table:ORDERS"}, []string{"table:CUSTOMERS"}},
		{"schema qualified", "SELECT * FROM sales.|", []string{"table:DEALS"}, []string{"table:ORDERS", "keyword:WHERE"}},
		{"database qualified", "SELECT * FROM db.sales.|", []string{"table:DEALS"}, []string{"table:ORDERS"}},
		{"another database", "SELECT * FROM otherdb.public.|", []string{"table:OTHER_T"}, []string{"table:ORDERS"}},
		{"another database's schemas", "SELECT * FROM otherdb.|", []string{"schema:PUBLIC"}, []string{"table:ORDERS"}},
		{"session schema", "SELECT * FROM public.|", []string{"table:ORDERS"}, []string{"table:DEALS"}},
		{"CTE name", "WITH recent AS (SELECT 1 AS n) SELECT * FROM |", []string{"table:RECENT", "table:ORDERS"}, nil},
		{"join", "SELECT * FROM orders o JOIN |", []string{"table:CUSTOMERS"}, nil},
		{"INSERT target", "INSERT INTO |", []string{"table:ORDERS"}, nil},
		{"DROP TABLE names tables", "DROP TABLE |", []string{"table:ORDERS"}, []string{"view:V_ORDERS"}},
		{"DROP VIEW names views", "DROP VIEW |", []string{"view:V_ORDERS"}, []string{"table:ORDERS"}},
		{"incomplete ALTER TABLE", "ALTER TABLE public.|", []string{"table:ORDERS"}, nil},

		// Columns in scope.
		{"select list", "SELECT | FROM orders", []string{"column:ID", "column:CUSTOMER_ID", "table:orders", "function:COALESCE()", "function:MY_UDF()"}, []string{"column:NAME", "function:SALES_UDF()"}},
		{"alias qualified", "SELECT o.| FROM orders o JOIN customers c ON o.id = c.id", []string{"column:CUSTOMER_ID"}, []string{"column:NAME", "function:COALESCE()"}},
		{"table qualified", "SELECT orders.| FROM orders", []string{"column:CUSTOMER_ID"}, nil},
		{"schema and table qualified", "SELECT public.orders.| FROM orders", []string{"column:CUSTOMER_ID"}, nil},
		{"alias hides the table name", "SELECT orders.| FROM orders o", nil, []string{"column:CUSTOMER_ID"}},
		{"WHERE", "SELECT * FROM orders WHERE |", []string{"column:SRC", "keyword:NOT"}, nil},
		{"QUALIFY sees select aliases", "SELECT id AS ident FROM orders QUALIFY |", []string{"column:IDENT", "column:ID"}, nil},
		{"ORDER BY sees select aliases", "SELECT id AS ident FROM orders ORDER BY |", []string{"column:IDENT"}, nil},
		{"caret before an incomplete tail", "SELECT | FROM orders WHERE", []string{"column:CUSTOMER_ID"}, nil},
		{"unclosed call", "SELECT coalesce(|", []string{"function:COALESCE()"}, nil},
		{"quoted table's columns", `SELECT | FROM "lower_t"`, []string{`column:"mixedCol"`, `column:"ORDER"`}, nil},
		{"an unquoted name folds to upper case", "SELECT | FROM lower_t", nil, []string{`column:"mixedCol"`}},
		{"view columns", "SELECT | FROM v_orders", []string{"column:ID"}, []string{"column:CUSTOMER_ID"}},
		{"another database's columns", "SELECT | FROM otherdb.public.other_t", []string{"column:OTHER_COL"}, nil},

		// Scope rules.
		{"ON sees its join's operands", "SELECT * FROM orders a, customers b JOIN orders c ON |", []string{"table:b", "table:c"}, []string{"table:a"}},
		{"ON does not see later joins", "SELECT * FROM orders a JOIN customers b ON | JOIN orders c ON c.id = a.id", []string{"table:a", "table:b"}, []string{"table:c"}},
		{"USING lists both operands' columns", "SELECT * FROM orders a JOIN customers b USING (|)", []string{"column:CUSTOMER_ID", "column:NAME"}, []string{"keyword:SELECT"}},
		{"LATERAL subquery sees earlier items", "SELECT * FROM orders a, LATERAL (SELECT * FROM customers WHERE id = a.|)", []string{"column:CUSTOMER_ID"}, nil},
		{"derived table does not see siblings", "SELECT * FROM orders a, (SELECT * FROM customers WHERE id = a.|)", nil, []string{"column:CUSTOMER_ID"}},
		{"FLATTEN sees earlier items", "SELECT * FROM orders o, LATERAL FLATTEN(input => o.|)", []string{"column:SRC"}, nil},
		{"FLATTEN columns", "SELECT f.| FROM orders o, LATERAL FLATTEN(input => o.src) f", []string{"column:VALUE", "column:KEY", "column:PATH", "column:INDEX", "column:SEQ", "column:THIS"}, nil},
		{"TABLE(FLATTEN) columns", "SELECT f.| FROM orders o, TABLE(FLATTEN(o.src)) f", []string{"column:VALUE"}, nil},
		{"inner alias hides outer", "SELECT * FROM orders o WHERE EXISTS (SELECT 1 FROM customers o WHERE o.|)", []string{"column:NAME"}, []string{"column:CUSTOMER_ID"}},
		{"correlated outer alias", "SELECT * FROM orders o WHERE EXISTS (SELECT 1 FROM customers c WHERE c.id = o.|)", []string{"column:CUSTOMER_ID"}, nil},
		{"derived table outputs", "SELECT x.| FROM (SELECT id, name AS nm FROM customers) x", []string{"column:ID", "column:NM"}, []string{"column:NAME"}},
		{"star with EXCLUDE", "SELECT x.| FROM (SELECT * EXCLUDE (name) FROM customers) x", []string{"column:ID"}, []string{"column:NAME"}},
		{"CTE column list", "WITH a(p, q) AS (SELECT 1, 2), b AS (SELECT * FROM a) SELECT b.| FROM b", []string{"column:P", "column:Q"}, nil},
		{"recursive CTE", "WITH RECURSIVE r AS (SELECT 1 AS n UNION ALL SELECT n + 1 FROM r WHERE r.| < 3) SELECT * FROM r", []string{"column:N"}, nil},
		{"a CTE body does not see later CTEs", "WITH a AS (SELECT * FROM |), b AS (SELECT 1 AS z) SELECT 1", []string{"table:ORDERS"}, []string{"table:B", "table:A"}},
		{"recursive CTE without RECURSIVE", "WITH r(n) AS (SELECT 1 UNION ALL SELECT r.| + 1 FROM r WHERE n < 3) SELECT * FROM r", []string{"column:N"}, nil},
		{"long UNION chain", "WITH x AS (" + strings.Repeat("SELECT id FROM orders UNION ALL ", 12) + "SELECT id FROM orders) SELECT x.| FROM x", []string{"column:ID"}, nil},
		{"chained RENAME", "SELECT x.| FROM (SELECT * RENAME (id AS customer_id, customer_id AS old_customer_id) FROM orders) x", []string{"column:CUSTOMER_ID", "column:OLD_CUSTOMER_ID", "column:SRC"}, []string{"column:ID"}},
		{"database qualifier in an expression", "SELECT otherdb.|", []string{"schema:PUBLIC"}, nil},
		{"unaliased expression outputs", "SELECT x.| FROM (SELECT COUNT(*), src:a, coalesce(id, 'x') FROM orders) x", []string{`column:"COUNT(*)"`, `column:"SRC:A"`, `column:"COALESCE(ID, 'x')"`}, nil},
		{"multi-table INSERT VALUES sees the source", "INSERT ALL INTO orders (id) VALUES (|) SELECT id, name FROM customers", []string{"column:ID", "column:NAME"}, nil},
		{"multi-table INSERT WHEN sees the source", "INSERT FIRST WHEN | > 1 THEN INTO orders SELECT id, name FROM customers", []string{"column:NAME"}, nil},
		{"earlier select-list alias", "SELECT id AS ident, iden| FROM orders", []string{"column:IDENT"}, nil},
		{"later select-list alias is not visible", "SELECT iden|, id AS ident FROM orders", nil, []string{"column:IDENT"}},
		{"quoted alias as qualifier", `SELECT My| FROM orders AS "My Alias"`, []string{`table:"My Alias"`}, nil},
		{"quoted alias inside quotes", `SELECT "My|" FROM orders AS "My Alias"`, []string{"table:My Alias"}, nil},
		{"BY NAME operand sees the local WITH", "WITH x AS (WITH a AS (SELECT 1 AS n) SELECT 2 AS m UNION ALL BY NAME SELECT * FROM a) SELECT x.| FROM x", []string{"column:M", "column:N"}, nil},
		{"PIVOT aggregate alias suffix", "SELECT p.| FROM orders PIVOT(SUM(id) AS total FOR customer_id IN (1 AS one, 2)) p", []string{"column:ONE_TOTAL", `column:"2_TOTAL"`}, []string{"column:ONE"}},
		{"table function before its parentheses", "SELECT * FROM TABLE(flat|)", []string{"function:FLATTEN()"}, nil},
		{"MATCH_RECOGNIZE symbol qualifier", "SELECT * FROM orders MATCH_RECOGNIZE (ORDER BY id MEASURES FINAL FIRST(a.|) AS f PATTERN (a+ b*) DEFINE a AS id > 0)", []string{"column:CUSTOMER_ID"}, nil},
		{"VALUES columns", "SELECT v.| FROM (VALUES (1, 2)) v", []string{"column:COLUMN1", "column:COLUMN2"}, nil},
		{"PIVOT sees its table", "SELECT * FROM orders PIVOT (SUM(|) FOR id IN (1, 2))", []string{"column:CUSTOMER_ID"}, nil},
		{"later set-operation operand sees the WITH", "WITH a AS (SELECT 1 AS n) SELECT * FROM a UNION ALL SELECT a.| FROM a", []string{"column:N"}, nil},
		{"later set-operation operand's FROM", "WITH a AS (SELECT 1 AS n) SELECT * FROM a UNION ALL SELECT * FROM |", []string{"table:A"}, nil},
		{"PIVOT alias and columns", "SELECT p.| FROM (SELECT 1 AS empid, 2 AS amount, 3 AS quarter) PIVOT(SUM(amount) FOR quarter IN (3 AS q3)) AS p", []string{"column:EMPID", "column:Q3"}, []string{"column:AMOUNT", "column:QUARTER"}},
		{"PIVOT value columns", "SELECT p.| FROM orders PIVOT(SUM(id) FOR customer_id IN (1, 'a')) AS p", []string{`column:"1"`, `column:"'a'"`, "column:SRC"}, []string{"column:ID"}},
		{"UNION BY NAME outputs", "WITH x AS (SELECT 1 AS a UNION ALL BY NAME SELECT 2 AS b) SELECT x.| FROM x", []string{"column:A", "column:B"}, nil},
		{"UNION outputs are the left's", "WITH x AS (SELECT 1 AS a UNION ALL SELECT 2 AS b) SELECT x.| FROM x", []string{"column:A"}, []string{"column:B"}},
		{"star ILIKE", "SELECT x.| FROM (SELECT * ILIKE '%id' FROM orders) x", []string{"column:ID", "column:CUSTOMER_ID"}, []string{"column:SRC"}},
		{"star EXCLUDE list", "SELECT * EXCLUDE (|) FROM orders", []string{"column:CUSTOMER_ID"}, []string{"function:COALESCE()"}},
		{"qualified star RENAME", "SELECT o.* RENAME (| AS x) FROM orders o, customers c", []string{"column:CUSTOMER_ID"}, []string{"column:NAME"}},
		{"UNPIVOT IN list", "SELECT * FROM orders UNPIVOT(val FOR col IN (|))", []string{"column:CUSTOMER_ID"}, nil},
		{"UNPIVOT alias and columns", "SELECT u.| FROM orders UNPIVOT(val FOR col IN (id, customer_id)) u", []string{"column:COL", "column:VAL", "column:SRC"}, []string{"column:ID"}},
		{"MATCH_RECOGNIZE alias and columns", "SELECT m.| FROM orders MATCH_RECOGNIZE(PARTITION BY customer_id ORDER BY id MEASURES COUNT(*) AS cnt PATTERN (a+) DEFINE a AS id > 0) AS m", []string{"column:CNT", "column:CUSTOMER_ID"}, []string{"column:SRC"}},
		{"MATCH_RECOGNIZE sees its table", "SELECT * FROM orders MATCH_RECOGNIZE(PARTITION BY | ORDER BY id PATTERN (a) DEFINE a AS true)", []string{"column:CUSTOMER_ID"}, nil},

		// DML.
		{"INSERT column list", "INSERT INTO orders (|", []string{"column:CUSTOMER_ID"}, []string{"function:COALESCE()"}},
		{"multi-table INSERT column list", "INSERT ALL INTO orders (|) INTO customers (id) SELECT id FROM customers", []string{"column:CUSTOMER_ID"}, []string{"column:NAME"}},
		{"UPDATE SET", "UPDATE orders SET |", []string{"column:CUSTOMER_ID"}, nil},
		{"UPDATE FROM", "UPDATE orders o SET id = 1 FROM customers c WHERE |", []string{"table:o", "table:c", "column:NAME", "column:SRC"}, nil},
		{"DELETE USING", "DELETE FROM orders USING customers c WHERE c.|", []string{"column:NAME"}, []string{"column:SRC"}},
		{"MERGE ON", "MERGE INTO orders t USING customers s ON t.id = s.|", []string{"column:NAME"}, []string{"column:SRC"}},

		// Functions.
		{"call parentheses", "SELECT coal|", []string{"function:COALESCE()"}, nil},
		{"no parentheses before one", "SELECT coal|(1, 2)", []string{"function:COALESCE"}, []string{"function:COALESCE()"}},
		{"a keyword built-in is offered once", "SELECT |", []string{"keyword:CAST", "function:LEFT()"}, []string{"function:CAST()", "function:CASE()"}},
		{"schema functions", "SELECT sales.|", []string{"function:SALES_UDF()", "table:DEALS"}, []string{"function:COALESCE()"}},

		// Quoting.
		{"caret in a quoted name", `SELECT * FROM "low|`, []string{"table:lower_t"}, []string{`table:"lower_t"`, "keyword:WHERE"}},
		{"reserved column name", `SELECT | FROM "lower_t"`, []string{`column:"ORDER"`}, []string{"column:ORDER"}},

		// Statements of a script.
		{"second statement", "SELECT 1; SELECT | FROM customers", []string{"column:NAME"}, []string{"column:SRC"}},
		{"after a semicolon", "SELECT * FROM orders WHERE id = 1;|", []string{"keyword:SELECT"}, []string{"column:ID"}},
		{"broken earlier statement", "SELECT FROM FROM; SELECT | FROM customers", []string{"column:NAME"}, nil},
	}
	cfg := testConfig()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := labels(completeAt(t, cfg, tt.text))
			for _, w := range tt.want {
				if !slices.Contains(got, w) {
					t.Errorf("%q: missing %s; got %v", tt.text, w, got)
				}
			}
			for _, a := range tt.absent {
				if slices.Contains(got, a) {
					t.Errorf("%q: unexpected %s", tt.text, a)
				}
			}
		})
	}
}

// TestCompleteSessionSearchPath: a query resolves an unqualified name in the
// current schema, then PUBLIC; DDL and DML in the current schema alone.
func TestCompleteSessionSearchPath(t *testing.T) {
	cfg := testConfig()
	cfg.DefaultSchema = "SALES"
	tests := []struct {
		text   string
		want   []string
		absent []string
	}{
		{"SELECT * FROM |", []string{"table:DEALS", "table:ORDERS"}, nil},
		{"INSERT INTO |", []string{"table:DEALS"}, []string{"table:ORDERS"}},
		{"SELECT | FROM orders", []string{"column:CUSTOMER_ID"}, nil},
		{"SELECT | FROM customers", []string{"column:SALES_REP"}, []string{"column:NAME"}},
		{"SELECT public.customers.| FROM public.customers", []string{"column:NAME"}, nil},
		{"SELECT | FROM orders o", []string{"function:SALES_UDF()", "function:MY_UDF()"}, nil},
		{"UPDATE orders SET |", nil, []string{"column:CUSTOMER_ID"}},
		{"DELETE FROM deals WHERE |", []string{"column:AMOUNT", "function:SALES_UDF()"}, []string{"function:MY_UDF()"}},
	}
	for _, tt := range tests {
		got := labels(completeAt(t, cfg, tt.text))
		for _, w := range tt.want {
			if !slices.Contains(got, w) {
				t.Errorf("%q: missing %s; got %v", tt.text, w, got)
			}
		}
		for _, a := range tt.absent {
			if slices.Contains(got, a) {
				t.Errorf("%q: unexpected %s", tt.text, a)
			}
		}
	}
}

func TestCompleteSessionNothing(t *testing.T) {
	// Nothing completes in a semi-structured path, a stage, a variable, a
	// string or a comment.
	for _, text := range []string{
		"SELECT src:| FROM orders",
		"SELECT src:customer.| FROM orders",
		"SELECT src['a'].| FROM orders",
		"SELECT $1 FROM @|",
		"SELECT $1 FROM @mystage/|",
		"SELECT $1|",
		"SELECT 'abc|",
		"SELECT * FROM orders -- |",
		"SELECT * FROM orders /* |",
		"SELECT * FROM orders o AS |",
	} {
		if got := completeAt(t, testConfig(), text); len(got) > 0 {
			t.Errorf("%q: got %v", text, labels(got))
		}
	}
}

func TestCompleteSessionWithoutMetadata(t *testing.T) {
	for _, cfg := range []SessionConfig{
		{DefaultDatabase: "DB"},
		{DefaultDatabase: "DB", LoadDatabase: func(context.Context, string) (*metadatapb.DatabaseSchemaMetadata, error) {
			return nil, errors.New("metadata unavailable")
		}},
	} {
		got := completeAt(t, cfg, "SELECT | FROM orders")
		for _, c := range got {
			switch c.Type {
			case CandidateKeyword, CandidateFunction:
			case CandidateTable:
				if c.Text != "orders" {
					t.Errorf("unexpected %v", c)
				}
			default:
				t.Errorf("unexpected %v", c)
			}
		}
		if !slices.Contains(labels(got), "function:COALESCE()") {
			t.Errorf("built-in functions missing: %v", labels(got))
		}
	}
}

func TestCompleteSessionBoundsDatabaseLoads(t *testing.T) {
	var loaded []string
	cfg := testConfig()
	load := cfg.LoadDatabase
	cfg.LoadDatabase = func(ctx context.Context, name string) (*metadatapb.DatabaseSchemaMetadata, error) {
		loaded = append(loaded, name)
		return load(ctx, name)
	}
	completeAt(t, cfg, "SELECT | FROM a.s.t, b.s.t, c.s.t, d.s.t, otherdb.public.other_t, a.s.u")
	if len(loaded) > 1+maxDatabaseLoads {
		t.Errorf("loaded %v, want at most %d databases", loaded, 1+maxDatabaseLoads)
	}
	for _, name := range loaded {
		if n := countOf(loaded, name); n > 1 {
			t.Errorf("loaded %s %d times", name, n)
		}
	}
}

func countOf(list []string, s string) int {
	n := 0
	for _, x := range list {
		if x == s {
			n++
		}
	}
	return n
}

func TestCompleteSessionCaretPosition(t *testing.T) {
	// The caret is a UTF-16 column on its line: 😀 takes two code units.
	stmt := "SELECT '😀', \nNAME, | FROM customers"
	got := completeAt(t, testConfig(), stmt)
	if !slices.Contains(labels(got), "column:NAME") {
		t.Fatalf("got %v", labels(got))
	}
	text := "SELECT '😀', id FROM customers WHERE "
	got = CompleteSession(context.Background(), testConfig(), text, 1, len([]rune(text))+1)
	if !slices.Contains(labels(got), "column:NAME") {
		t.Fatalf("got %v", labels(got))
	}
}

func TestStatementAtCaret(t *testing.T) {
	tests := []struct {
		script string
		pos    int
		stmt   string
		at     int
	}{
		{"SELECT 1; SELECT 2", 3, "SELECT 1", 3},
		{"SELECT 1; SELECT 2", 8, "SELECT 1", 8},
		{"SELECT 1; SELECT 2", 9, " SELECT 2", 0},
		{"SELECT 1;", 9, "", 0},
		{"SELECT 1;  ", 11, "  ", 2},
		{"", 0, "", 0},
	}
	for _, tt := range tests {
		stmt, at := statementAtCaret(tt.script, tt.pos)
		if stmt != tt.stmt || at != tt.at {
			t.Errorf("statementAtCaret(%q, %d) = %q, %d; want %q, %d", tt.script, tt.pos, stmt, at, tt.stmt, tt.at)
		}
	}
}

func TestQuoteIdentifierIfNeeded(t *testing.T) {
	tests := []struct {
		name     string
		inQuotes bool
		want     string
	}{
		{"ORDERS", false, "ORDERS"},
		{"ORDER_2$", false, "ORDER_2$"},
		{"orders", false, `"orders"`},
		{"Order Items", false, `"Order Items"`},
		{"SELECT", false, `"SELECT"`},
		{"1ST", false, `"1ST"`},
		{`A"B`, false, `"A""B"`},
		{`a"b`, true, `a""b`},
	}
	for _, tt := range tests {
		if got := QuoteIdentifierIfNeeded(tt.name, tt.inQuotes); got != tt.want {
			t.Errorf("QuoteIdentifierIfNeeded(%q, %v) = %q, want %q", tt.name, tt.inQuotes, got, tt.want)
		}
	}
}

// TestCompleteSessionLatency holds completion of a 2 KB statement to its
// budget, with a catalog of a few thousand tables.
func TestCompleteSessionLatency(t *testing.T) {
	dbs := testDatabases()
	public := dbs["DB"].Schemas[0]
	for i := range 3000 {
		public.Tables = append(public.Tables, &metadatapb.TableMetadata{Name: fmt.Sprintf("T%d", i), Columns: []*metadatapb.ColumnMetadata{col("A", "NUMBER"), col("B", "NUMBER")}})
	}
	cfg := SessionConfig{DefaultDatabase: "DB", LoadDatabase: func(_ context.Context, name string) (*metadatapb.DatabaseSchemaMetadata, error) {
		return dbs[name], nil
	}}
	var b strings.Builder
	b.WriteString("SELECT ")
	for i := 0; b.Len() < 1500; i++ {
		fmt.Fprintf(&b, "o.id + coalesce(c.id, %d) AS a%d, ", i, i)
	}
	b.WriteString("o.src FROM orders o JOIN customers c ON o.customer_id = c.id, LATERAL FLATTEN(input => o.src) f WHERE ")
	for b.Len() < 2048 {
		b.WriteString("c.id > 1 AND ")
	}
	for _, text := range []string{
		b.String() + "|",
		b.String() + "o.|",
		b.String()[:900] + "|" + b.String()[900:],
		b.String() + "| ((( SELECT FROM WHERE",
	} {
		var took []time.Duration
		for range 20 {
			start := time.Now()
			completeAt(t, cfg, text)
			took = append(took, time.Since(start))
		}
		slices.Sort(took)
		if p95 := took[len(took)*95/100]; p95 > 100*time.Millisecond {
			t.Errorf("p95 %v over 100ms at caret %d of %d bytes", p95, strings.Index(text, "|"), len(text)-1)
		}
	}
}
