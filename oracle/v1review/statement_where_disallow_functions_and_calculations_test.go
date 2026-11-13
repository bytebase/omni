package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/v1review"
)

func TestStatementWhereDisallowFunctionsAndCalculations(t *testing.T) {
	table := func(name string, indexed ...string) *metadata.TableMetadata {
		return &metadata.TableMetadata{
			Name:    name,
			Columns: []*metadata.ColumnMetadata{{Name: "ID"}, {Name: "NAME"}, {Name: "NOTE"}},
			Indexes: []*metadata.IndexMetadata{{Name: "IX", Expressions: indexed}},
		}
	}
	logs := table("LOGS", "ID")
	logs.Partitions = []*metadata.TablePartitionMetadata{{Name: "P1", Subpartitions: []*metadata.TablePartitionMetadata{{Name: "SP1"}}}}
	db := &metadata.DatabaseSchemaMetadata{Name: "SHOP", Schemas: []*metadata.SchemaMetadata{
		{Name: "", Tables: []*metadata.TableMetadata{table("BOOK", "ID", "NAME"), table("Mixed", "Id"), logs}},
		{Name: "HR", Tables: []*metadata.TableMetadata{table("BOOK", "NOTE")}},
	}}
	sensitive := v1review.Schema{Metadata: db, ObjectCaseSensitive: true}
	type whereFinding struct {
		Code   v1review.Code
		Name   string
		Detail string
		Text   string
	}
	fn := func(col, name, text string) whereFinding {
		return whereFinding{Code: v1review.FunctionOnIndexedColumn, Name: col, Detail: name, Text: text}
	}
	calc := func(col, text string) whereFinding {
		return whereFinding{Code: v1review.CalculationOnIndexedColumn, Name: col, Text: text}
	}
	for _, tt := range []struct {
		name     string
		sql      string
		schema   v1review.Schema
		database string
		want     []whereFinding
	}{
		{name: "function, column upper-cased", sql: "SELECT * FROM book WHERE upper(name) = 'a'", schema: sensitive, want: []whereFinding{fn("NAME", "UPPER", "upper(name)")}},
		{name: "qualifier dropped from function name", sql: "SELECT * FROM book WHERE pkg.fn(id) = 1", schema: sensitive, want: []whereFinding{fn("ID", "FN", "pkg.fn(id)")}},
		{name: "unindexed column", sql: "SELECT * FROM book WHERE abs(note) = 1", schema: sensitive},
		{name: "calculation", sql: "SELECT * FROM book WHERE id * 2 = 1", schema: sensitive, want: []whereFinding{calc("ID", "id * 2")}},
		{name: "concatenation is a calculation", sql: "SELECT * FROM book WHERE name || 'x' = 'y'", schema: sensitive, want: []whereFinding{calc("NAME", "name || 'x'")}},
		{name: "unary minus", sql: "SELECT * FROM book WHERE -id = 1", schema: sensitive, want: []whereFinding{calc("ID", "-id")}},
		{name: "comparison alone is not", sql: "SELECT * FROM book WHERE id = 1 AND name <> 'a'", schema: sensitive},
		{name: "column nested in a function argument is not", sql: "SELECT * FROM book WHERE abs(id + 1) = 2", schema: sensitive, want: []whereFinding{calc("ID", "id + 1"), calc("ID", "id + 1")}},
		{name: "nested node reported once per path", sql: "SELECT * FROM book WHERE abs(abs(id)) = 1", schema: sensitive, want: []whereFinding{fn("ID", "ABS", "abs(id)"), fn("ID", "ABS", "abs(id)")}},
		{name: "prior", sql: "SELECT * FROM book CONNECT BY abs(PRIOR id) = id", schema: sensitive, want: []whereFinding{fn("ID", "ABS", "abs(PRIOR id)")}},
		{name: "quoted table keeps its case, columns compare upper-cased", sql: `SELECT * FROM "Mixed" WHERE abs("Id") = 1`, schema: sensitive, want: []whereFinding{fn("ID", "ABS", `abs("Id")`)}},
		{name: "unquoted table is upper-cased", sql: `SELECT * FROM mixed WHERE abs(id) = 1`, schema: sensitive},
		{name: "case-insensitive lower-cases", sql: `SELECT * FROM mixed WHERE abs("ID") = 1`, schema: v1review.Schema{Metadata: db}, want: []whereFinding{fn("id", "ABS", `abs("ID")`)}},
		{name: "partition finds its table", sql: "SELECT * FROM sp1 WHERE abs(id) = 1", schema: sensitive, want: []whereFinding{fn("ID", "ABS", "abs(id)")}},
		{name: "qualifier ignored without current database", sql: "SELECT * FROM hr.book WHERE abs(id) = 1", schema: sensitive, want: []whereFinding{fn("ID", "ABS", "abs(id)")}},
		{name: "own owner resolves", sql: "SELECT * FROM shop.book WHERE abs(id) = 1", schema: sensitive, database: "SHOP", want: []whereFinding{fn("ID", "ABS", "abs(id)")}},
		{name: "other owner does not", sql: "SELECT * FROM hr.book WHERE abs(note) = 1 OR abs(id) = 1", schema: sensitive, database: "SHOP"},
		{name: "database link ignored", sql: "SELECT * FROM book@remote WHERE abs(id) = 1", schema: sensitive, want: []whereFinding{fn("ID", "ABS", "abs(id)")}},
		{name: "cte hides the table", sql: "WITH book AS (SELECT 1 id FROM dual) SELECT * FROM book WHERE abs(id) = 1", schema: sensitive},
		{name: "no snapshot", sql: "SELECT * FROM book WHERE abs(id) = 1"},
		{name: "plsql block judged", sql: "BEGIN DELETE FROM book WHERE abs(id) = 1; END;", schema: sensitive, want: []whereFinding{fn("ID", "ABS", "abs(id)")}},
		{name: "procedure body not judged", sql: "CREATE PROCEDURE p AS BEGIN DELETE FROM book WHERE abs(id) = 1; END;", schema: sensitive},
		{name: "explain not judged", sql: "EXPLAIN PLAN FOR SELECT * FROM book WHERE abs(id) = 1", schema: sensitive},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var stmts []ast.Node
			for _, n := range parse(t, tt.sql) {
				// Bytebase hands the rules each statement unwrapped.
				if raw, ok := n.(*ast.RawStmt); ok {
					n = raw.Stmt
				}
				stmts = append(stmts, n)
			}
			var got []whereFinding
			for _, f := range StatementWhereDisallowFunctionsAndCalculations(stmts, tt.schema, "", tt.database) {
				got = append(got, whereFinding{f.Code, f.Name, f.Detail, tt.sql[f.Start:f.End]})
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
