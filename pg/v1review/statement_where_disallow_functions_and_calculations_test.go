package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/v1review"
)

func TestStatementWhereDisallowFunctionsAndCalculations(t *testing.T) {
	table := func(name string, indexed ...string) *metadata.TableMetadata {
		return &metadata.TableMetadata{
			Name:    name,
			Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "total"}, {Name: "Note"}},
			Indexes: []*metadata.IndexMetadata{{Name: "ix", Expressions: indexed}},
		}
	}
	items := &metadata.TableMetadata{Name: "items", Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "qty"}}}
	logs := table("logs", "id")
	logs.Partitions = []*metadata.TablePartitionMetadata{{Name: "logs_2024"}}
	schema := v1review.Schema{Metadata: &metadata.DatabaseSchemaMetadata{Name: "shop", Schemas: []*metadata.SchemaMetadata{
		{Name: "public", Tables: []*metadata.TableMetadata{table("orders", "id"), items, logs}},
		{Name: "app", Tables: []*metadata.TableMetadata{table("orders", "total")}},
		{Name: "sales", Tables: []*metadata.TableMetadata{table("Orders", "Note"), table("orders", "total")}},
	}}}
	sensitive := schema
	sensitive.ObjectCaseSensitive = true
	type whereFinding struct {
		Code               v1review.Code
		Name, Detail, Text string
	}
	fn := func(col, name, text string) whereFinding {
		return whereFinding{Code: v1review.FunctionOnIndexedColumn, Name: col, Detail: name, Text: text}
	}
	calc := func(col, text string) whereFinding {
		return whereFinding{Code: v1review.CalculationOnIndexedColumn, Name: col, Text: text}
	}
	for _, tt := range []struct {
		name       string
		sql        string
		schema     v1review.Schema
		searchPath []string
		want       []whereFinding
	}{
		{name: "first search path schema wins", sql: "SELECT * FROM orders WHERE abs(id) = 1 AND total + 1 = 2", schema: schema, searchPath: []string{"app", "public"}, want: []whereFinding{calc("total", "total + 1")}},
		{name: "later search path schema", sql: "SELECT * FROM logs WHERE upper(id) = 'x'", schema: schema, searchPath: []string{"app", "public"}, want: []whereFinding{fn("id", "UPPER", "upper(id)")}},
		{name: "empty search path resolves only qualified", sql: "SELECT * FROM orders WHERE abs(id) = 1; SELECT * FROM public.orders WHERE -id = 1", schema: schema, want: []whereFinding{calc("id", "-id")}},
		{name: "partition name finds its table", sql: "SELECT * FROM logs_2024 WHERE id * 2 = 1", schema: schema, searchPath: []string{"public"}, want: []whereFinding{calc("id", "id * 2")}},
		{name: "case-insensitive folds names, last table wins", sql: `SELECT * FROM sales."Orders" WHERE upper("Note") = 'x' OR total - 1 = 0`, schema: schema, want: []whereFinding{calc("total", "total - 1")}},
		{name: "case-sensitive keeps names", sql: `SELECT * FROM sales."Orders" WHERE pg_catalog.upper("Note") = 'x' OR total - 1 = 0`, schema: sensitive, want: []whereFinding{fn("Note", "UPPER", `pg_catalog.upper("Note")`)}},
		{name: "calculation inside function is a calculation", sql: "SELECT * FROM orders WHERE abs(id + 1) = 2 AND id::text = '1' AND abs(id::bigint) = 1", schema: schema, searchPath: []string{"public"}, want: []whereFinding{calc("id", "id + 1"), fn("id", "ABS", "abs(id::bigint)")}},
		{name: "local column shadows outer", sql: "SELECT * FROM orders o WHERE EXISTS (SELECT 1 FROM items WHERE abs(id) = 1 AND abs(o.id) = qty)", schema: schema, searchPath: []string{"public"}, want: []whereFinding{fn("id", "ABS", "abs(o.id)")}},
		{name: "non-recursive CTE body sees the table", sql: "WITH orders AS (SELECT * FROM orders WHERE abs(id) = 1) SELECT * FROM orders WHERE abs(id) = 2", schema: schema, searchPath: []string{"public"}, want: []whereFinding{fn("id", "ABS", "abs(id)")}},
		{name: "recursive CTE body sees itself", sql: "WITH RECURSIVE orders AS (SELECT * FROM orders WHERE abs(id) = 1) SELECT 1", schema: schema, searchPath: []string{"public"}},
		{name: "join condition judged, having not", sql: "SELECT id FROM orders o JOIN items i ON abs(o.id) = i.qty GROUP BY id HAVING abs(id) > 1", schema: schema, searchPath: []string{"public"}, want: []whereFinding{fn("id", "ABS", "abs(o.id)")}},
		{name: "update delete merge", sql: "UPDATE orders SET total = 1 WHERE -id = 1; DELETE FROM orders WHERE id % 2 = 0; MERGE INTO orders t USING items s ON abs(t.id) = s.qty WHEN MATCHED THEN DELETE", schema: schema, searchPath: []string{"public"}, want: []whereFinding{calc("id", "-id"), calc("id", "id % 2"), fn("id", "ABS", "abs(t.id)")}},
		{name: "explain and routine bodies not judged", sql: "EXPLAIN SELECT * FROM orders WHERE abs(id) = 1; CREATE FUNCTION f() RETURNS int LANGUAGE sql AS $$ SELECT id FROM orders WHERE abs(id) = 1 $$", schema: schema, searchPath: []string{"public"}},
		{name: "column printed folded", sql: `SELECT * FROM public.orders WHERE abs("ID") = 1`, schema: schema, want: []whereFinding{fn("id", "ABS", `abs("ID")`)}},
		{name: "no snapshot", sql: "SELECT * FROM orders WHERE abs(id) = 1", searchPath: []string{"public"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			stmts := parseStatements(t, tt.sql)
			findings := StatementWhereDisallowFunctionsAndCalculations(stmts, tt.schema, tt.searchPath)
			var got []whereFinding
			for _, f := range findings {
				got = append(got, whereFinding{f.Code, f.Name, f.Detail, tt.sql[f.Start:f.End]})
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
