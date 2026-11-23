package analysis

import (
	"errors"
	"strings"
	"testing"

	"github.com/bytebase/omni/doris/parser"
)

func spanResources(span *ResolvedQuerySpan) map[string]bool {
	resources := make(map[string]bool, len(span.SourceColumns))
	for resource := range span.SourceColumns {
		resources[resource.Database+"."+resource.Table] = true
	}
	return resources
}

func TestResolveQuerySpanClassification(t *testing.T) {
	tests := []struct {
		statement string
		want      QueryType
	}{
		{"SELECT * FROM users", QueryTypeSelect},
		{"INSERT INTO users (id, name) VALUES (1, 'test')", QueryTypeDML},
		{"UPDATE users SET name = 'test' WHERE id = 1", QueryTypeDML},
		{"DELETE FROM users WHERE id = 1", QueryTypeDML},
		{"CREATE TABLE test (id INT)", QueryTypeDDL},
		{"SHOW DATABASES", QueryTypeSelectInfoSchema},
		{"SHOW TABLES", QueryTypeSelectInfoSchema},
		{"SHOW DATA", QueryTypeSelectInfoSchema},
		{"SHOW DATA FROM db1", QueryTypeSelectInfoSchema},
		{"SHOW TABLETS FROM table1", QueryTypeSelectInfoSchema},
		{"SHOW VARIABLES", QueryTypeSelectInfoSchema},
		{"SHOW CREATE TABLE users", QueryTypeSelectInfoSchema},
		{"(SELECT * FROM users)", QueryTypeSelect},
		{"(SELECT 1 LIMIT 1) LIMIT 2", QueryTypeSelect},
		{"(SELECT * FROM users) ORDER BY 1 LIMIT 5", QueryTypeSelect},
		// EXPLAIN defers to the statement it plans.
		{"EXPLAIN SELECT * FROM users", QueryTypeSelect},
		{"EXPLAIN INSERT INTO users (id) VALUES (1)", QueryTypeDML},
		{"EXPLAIN DROP TABLE users", QueryTypeDDL},
		// USE is Unknown so a permission check denies it rather than
		// permitting it under the DDL grant.
		{"USE db1", QueryTypeUnknown},
	}
	for _, tc := range tests {
		t.Run(tc.statement, func(t *testing.T) {
			span, err := ResolveQuerySpan(tc.statement, ResolveOptions{DefaultDatabase: "testdb"})
			if err != nil {
				t.Fatalf("ResolveQuerySpan: %v", err)
			}
			if span.Type != tc.want {
				t.Errorf("Type = %s, want %s", span.Type, tc.want)
			}
		})
	}
}

func TestResolveQuerySpanTables(t *testing.T) {
	tests := []struct {
		name      string
		statement string
		want      []string
	}{
		{
			name:      "default database fills a bare name",
			statement: "SELECT a FROM t",
			want:      []string{"testdb.t"},
		},
		{
			name:      "qualified name keeps its database",
			statement: "SELECT * FROM db.t",
			want:      []string{"db.t"},
		},
		{
			name:      "catalog qualifier is dropped",
			statement: "SELECT * FROM catalog1.db1.t1",
			want:      []string{"db1.t1"},
		},
		{
			name: "CTE bodies and joins",
			statement: `WITH summary AS (
				SELECT user_id, SUM(amount) AS total_amount
				FROM orders
				WHERE created_at > '2024-01-01'
				GROUP BY user_id
			)
			SELECT u.id, u.name, s.total_amount
			FROM summary s
			LEFT JOIN users u ON s.user_id = u.id`,
			want: []string{"testdb.orders", "testdb.users"},
		},
		{
			name: "cross-database CTEs",
			statement: `WITH
				cte1 AS (SELECT * FROM db1.table1),
				cte2 AS (SELECT * FROM db2.table2 WHERE active = 1)
			SELECT c1.id, c2.name FROM cte1 c1 INNER JOIN cte2 c2 ON c1.ref_id = c2.id`,
			want: []string{"db1.table1", "db2.table2"},
		},
		{
			name:      "CTE shadows a table",
			statement: "WITH products AS (SELECT * FROM items) SELECT * FROM products",
			want:      []string{"testdb.items"},
		},
		{
			name:      "simple CTE",
			statement: "WITH cte AS (SELECT * FROM users)\nSELECT * FROM cte",
			want:      []string{"testdb.users"},
		},
		{
			name:      "multiple CTEs",
			statement: "WITH\ncte1 AS (SELECT * FROM table1),\ncte2 AS (SELECT * FROM table2)\nSELECT * FROM cte1 JOIN cte2",
			want:      []string{"testdb.table1", "testdb.table2"},
		},
		{
			name:      "nested CTEs",
			statement: "WITH outer_cte AS (\n  WITH inner_cte AS (SELECT * FROM base_table)\n  SELECT * FROM inner_cte\n)\nSELECT * FROM outer_cte",
			want:      []string{"testdb.base_table"},
		},
		{
			name: "CTE referencing another CTE",
			statement: `WITH
				base_data AS (SELECT id, name, department_id FROM employees WHERE active = true),
				dept_summary AS (SELECT d.name as dept_name, COUNT(b.id) as emp_count FROM base_data b JOIN departments d ON b.department_id = d.id GROUP BY d.name)
			SELECT * FROM dept_summary WHERE emp_count > 10`,
			want: []string{"testdb.employees", "testdb.departments"},
		},
		{
			name:      "EXPLAIN exposes the planned statement's tables",
			statement: "EXPLAIN SELECT * FROM users",
			want:      []string{"testdb.users"},
		},
		{
			name:      "EXPLAIN VERBOSE exposes tables",
			statement: "EXPLAIN VERBOSE SELECT id FROM t1 JOIN t2 ON t1.id = t2.id",
			want:      []string{"testdb.t1", "testdb.t2"},
		},
		{
			name:      "subquery in WHERE",
			statement: "SELECT a FROM t WHERE b IN (SELECT b FROM other)",
			want:      []string{"testdb.t", "testdb.other"},
		},
		{
			name:      "DML subquery sources",
			statement: "UPDATE target SET x = (SELECT x FROM secret)",
			want:      []string{"testdb.secret"},
		},
		{
			name:      "table function is not a table",
			statement: "SELECT * FROM BACKENDS()",
			want:      nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			span, err := ResolveQuerySpan(tc.statement, ResolveOptions{DefaultDatabase: "testdb"})
			if err != nil {
				t.Fatalf("ResolveQuerySpan: %v", err)
			}
			got := spanResources(span)
			if len(got) != len(tc.want) {
				t.Fatalf("SourceColumns = %v, want %v", got, tc.want)
			}
			for _, want := range tc.want {
				if !got[want] {
					t.Errorf("SourceColumns = %v, missing %s", got, want)
				}
			}
		})
	}
}

func TestResolveQuerySpanSystemTables(t *testing.T) {
	tests := []struct {
		statement string
		wantType  QueryType
		want      string
	}{
		{"SELECT * FROM information_schema.tables", QueryTypeSelectInfoSchema, "information_schema.tables"},
		{"SELECT table_name FROM information_schema.columns WHERE table_schema = 'testdb'", QueryTypeSelectInfoSchema, "information_schema.columns"},
		{"SELECT * FROM mysql.user", QueryTypeSelectInfoSchema, "mysql.user"},
		{"SELECT * FROM __internal_schema.audit_log", QueryTypeSelectInfoSchema, "__internal_schema.audit_log"},
		{"SELECT * FROM _statistics_.table_statistics", QueryTypeSelect, "_statistics_.table_statistics"},
		{"SELECT * FROM users", QueryTypeSelect, "testdb.users"},
	}
	for _, tc := range tests {
		t.Run(tc.statement, func(t *testing.T) {
			span, err := ResolveQuerySpan(tc.statement, ResolveOptions{DefaultDatabase: "testdb"})
			if err != nil {
				t.Fatalf("ResolveQuerySpan: %v", err)
			}
			if span.Type != tc.wantType {
				t.Errorf("Type = %s, want %s", span.Type, tc.wantType)
			}
			if got := spanResources(span); !got[tc.want] {
				t.Errorf("SourceColumns = %v, missing %s", got, tc.want)
			}
		})
	}
}

func TestResolveQuerySpanMixedTables(t *testing.T) {
	_, err := ResolveQuerySpan("SELECT * FROM users JOIN information_schema.tables", ResolveOptions{DefaultDatabase: "testdb"})
	if !errors.Is(err, ErrMixUserSystemTables) {
		t.Fatalf("err = %v, want ErrMixUserSystemTables", err)
	}
}

func TestResolveQuerySpanCaseSensitivity(t *testing.T) {
	// With case-folded names, MYSQL.user is a system table; without, it is a
	// user table.
	folded, err := ResolveQuerySpan("SELECT * FROM MYSQL.user", ResolveOptions{DefaultDatabase: "testdb", IgnoreCaseSensitive: true})
	if err != nil {
		t.Fatalf("ResolveQuerySpan: %v", err)
	}
	if folded.Type != QueryTypeSelectInfoSchema {
		t.Errorf("case-folded Type = %s, want %s", folded.Type, QueryTypeSelectInfoSchema)
	}
	exact, err := ResolveQuerySpan("SELECT * FROM MYSQL.user", ResolveOptions{DefaultDatabase: "testdb"})
	if err != nil {
		t.Fatalf("ResolveQuerySpan: %v", err)
	}
	if exact.Type != QueryTypeSelect {
		t.Errorf("case-sensitive Type = %s, want %s", exact.Type, QueryTypeSelect)
	}
}

func TestResolveQuerySpanEmptyAndMulti(t *testing.T) {
	span, err := ResolveQuerySpan("", ResolveOptions{DefaultDatabase: "testdb"})
	if err != nil {
		t.Fatalf("empty input: %v", err)
	}
	if span.Type != QueryTypeUnknown || len(span.SourceColumns) != 0 || len(span.Results) != 0 {
		t.Errorf("empty input span = %+v, want zero span", span)
	}

	for _, comment := range []string{"-- only a comment", "/* block */"} {
		span, err := ResolveQuerySpan(comment, ResolveOptions{DefaultDatabase: "testdb"})
		if err != nil {
			t.Fatalf("%q: %v", comment, err)
		}
		if span.Type != QueryTypeUnknown || len(span.SourceColumns) != 0 {
			t.Errorf("%q span = %+v, want zero span", comment, span)
		}
	}

	_, err = ResolveQuerySpan("SELECT 1; SELECT 2", ResolveOptions{DefaultDatabase: "testdb"})
	if err == nil || err.Error() != "expecting only one statement to get query span, but got 2" {
		t.Errorf("multi-statement error = %v", err)
	}
}

func TestResolveQuerySpanFailsClosed(t *testing.T) {
	for _, statement := range []string{
		"SELECT a FROM secret xx yy zz",
		"SELECT (SELECT 1 */ 2 FROM t2) x FROM t1",
		"SELECT EXISTS (DELETE FROM secret) FROM public",
		"SELECT EXISTS () FROM public",
	} {
		if _, err := ResolveQuerySpan(statement, ResolveOptions{DefaultDatabase: "db1"}); err == nil {
			t.Errorf("ResolveQuerySpan(%q) accepted invalid SQL", statement)
		}
	}

	// Parenthesized and repeated-clause queries still report their reads.
	for _, statement := range []string{
		"(SELECT a FROM secret) LIMIT 1",
		"(SELECT a FROM secret ORDER BY 1) ORDER BY 1",
	} {
		span, err := ResolveQuerySpan(statement, ResolveOptions{DefaultDatabase: "db1"})
		if err != nil {
			t.Fatalf("ResolveQuerySpan(%q): %v", statement, err)
		}
		if !span.SourceColumns[ColumnResource{Database: "db1", Table: "secret"}] {
			t.Errorf("%q SourceColumns = %v, missing db1.secret", statement, spanResources(span))
		}
	}
}

// A script with an empty leading segment reports its parse error in the
// script's coordinates, not the segment's.
func TestResolveQuerySpanReportsScriptPositions(t *testing.T) {
	const script = ";SELECT FROM"
	_, err := ResolveQuerySpan(script, ResolveOptions{DefaultDatabase: "testdb"})
	if err == nil {
		t.Fatal("ResolveQuerySpan accepted a syntax error")
	}
	errs := parser.AllErrors(err)
	if len(errs) == 0 {
		t.Fatalf("error = %v, want a ParseError", err)
	}
	if want := strings.Index(script, "FROM"); errs[0].Position != want {
		t.Errorf("error position = %d, want %d", errs[0].Position, want)
	}
}

func TestResolveQuerySpanFailsClosedOnLexError(t *testing.T) {
	// SplitStatements drops a chunk that holds no statement, and a dropped
	// chunk takes its lex error with it, so a buffer such as "/*" would
	// otherwise answer an empty span for text the engine rejects.
	for _, script := range []string{
		"/*",
		"SELECT * FROM secret; /*",
	} {
		span, err := ResolveQuerySpan(script, ResolveOptions{DefaultDatabase: "testdb"})
		if err == nil {
			t.Errorf("ResolveQuerySpan(%q) = %+v, want the lex error", script, span)
			continue
		}
		var parseErr *parser.ParseError
		if !errors.As(err, &parseErr) {
			t.Fatalf("ResolveQuerySpan(%q) error = %v, want a *parser.ParseError", script, err)
		}
		if want := strings.Index(script, "/*"); parseErr.Position != want {
			t.Errorf("ResolveQuerySpan(%q) error position = %d, want %d", script, parseErr.Position, want)
		}
	}

	// A well-formed comment-only buffer stays an empty span, not an error.
	for _, script := range []string{"-- note", "/* note */"} {
		span, err := ResolveQuerySpan(script, ResolveOptions{DefaultDatabase: "testdb"})
		if err != nil {
			t.Errorf("ResolveQuerySpan(%q): %v", script, err)
			continue
		}
		if span.Type != QueryTypeUnknown || len(span.SourceColumns) != 0 {
			t.Errorf("ResolveQuerySpan(%q) = %+v, want an empty span", script, span)
		}
	}
}

func TestResolveQuerySpanResultsLineage(t *testing.T) {
	// Without a catalog the span carries per-column lineage as written: each
	// output column lists the column references that feed it, unresolved.
	// Bytebase's adapter keeps Results empty so its masking sees what it saw
	// before this package moved; span_test.go holds the catalog's lineage.
	span, err := ResolveQuerySpan("SELECT a, b + 1 AS c FROM t", ResolveOptions{DefaultDatabase: "testdb"})
	if err != nil {
		t.Fatalf("ResolveQuerySpan: %v", err)
	}
	if len(span.Results) != 2 {
		t.Fatalf("Results = %+v, want 2 columns", span.Results)
	}
	for i, tc := range []struct {
		name   string
		source string
	}{
		{"a", "a"},
		{"c", "b"},
	} {
		if span.Results[i].Name != tc.name {
			t.Errorf("Results[%d].Name = %q, want %q", i, span.Results[i].Name, tc.name)
		}
		if len(span.Results[i].SourceColumns) != 1 || !span.Results[i].SourceColumns[ColumnResource{Column: tc.source}] {
			t.Errorf("Results[%d].SourceColumns = %+v, want column %q", i, span.Results[i].SourceColumns, tc.source)
		}
	}
}

func TestIsSystemDatabase(t *testing.T) {
	// The pinned 3.1.4 server answers SHOW DATABASES with information_schema,
	// mysql and __internal_schema (doris/testdata/pins.tsv). `_statistics_` is
	// not a Doris database, and `sys` is StarRocks's.
	for name, want := range map[string]bool{
		"information_schema": true,
		"mysql":              true,
		"__internal_schema":  true,
		"_statistics_":       false,
		"sys":                false,
		"app":                false,
	} {
		if got := IsSystemDatabase(name, false); got != want {
			t.Errorf("IsSystemDatabase(%q) = %v, want %v", name, got, want)
		}
	}
	if !IsSystemDatabase("INFORMATION_SCHEMA", true) {
		t.Error("IsSystemDatabase(\"INFORMATION_SCHEMA\", true) = false, want true")
	}
	if IsSystemDatabase("INFORMATION_SCHEMA", false) {
		t.Error("IsSystemDatabase(\"INFORMATION_SCHEMA\", false) = true, want false")
	}
}
