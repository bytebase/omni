package analysis

import (
	"testing"
)

// Unit tests for the pipe syntax's analysis: the classification and statement
// type a pipe statement carries, and the query span it yields. A read-only
// pipeline is a Select; a pipeline that writes is DML or DDL; a suffix-only
// statement's output columns cannot be resolved (the input table is implicit),
// so its span fails closed while still recording the tables its operators name.

func TestClassifyPipe(t *testing.T) {
	tests := []struct {
		sql  string
		want QueryType
	}{
		{"SELECT 1 |> WHERE true", Select},
		{"SELECT x FROM t |> AGGREGATE COUNT(*) GROUP BY x", Select},
		{"|> AGGREGATE", Select},
		{"|> WHERE x > 1 |> DISTINCT", Select},
		{"DESCRIBE t |> AGGREGATE COUNT(*)", Select},
		{"|> INSERT t", DML},
		{"SELECT 1 |> INSERT t", DML},
		{"|> CREATE TABLE t (x INT64)", DDL},
		{"SELECT 1 |> IF true THEN (|> INSERT t) ELSE (|> DISTINCT)", DML},
		{"|> TEE (|> INSERT t), (|> DISTINCT)", DML},
		{"|> EXPORT DATA OPTIONS(uri = 'x')", Unknown},
	}
	for _, tc := range tests {
		for _, dialect := range []Dialect{DialectBigQuery, DialectSpanner} {
			if got := ClassifySQL(tc.sql, dialect); got != tc.want {
				t.Errorf("%v ClassifySQL(%q) = %v, want %v", dialect, tc.sql, got, tc.want)
			}
		}
	}
}

func TestStatementTypeOfPipe(t *testing.T) {
	tests := []struct {
		sql  string
		want StatementType
	}{
		{"SELECT 1 |> WHERE true", StatementTypeUnspecified},
		{"|> AGGREGATE", StatementTypeUnspecified},
		{"|> INSERT t", Insert},
		{"SELECT 1 |> INSERT t", Insert},
		{"|> CREATE TABLE t (x INT64)", CreateTable},
		{"|> IF true THEN (|> INSERT t) ELSE (|> DISTINCT)", Insert},
	}
	for _, tc := range tests {
		file, err := parseFile(tc.sql)
		if err != nil {
			t.Fatalf("parse %q: %v", tc.sql, err)
		}
		if got := StatementTypeOf(file.Stmts[0]); got != tc.want {
			t.Errorf("StatementTypeOf(%q) = %v, want %v", tc.sql, got, tc.want)
		}
	}
}

func TestGetQuerySpanPipe(t *testing.T) {
	t.Run("a query keeps its projection and records pipe tables", func(t *testing.T) {
		span, err := GetQuerySpan("SELECT x FROM t |> JOIN u ON true |> WHERE x > 1", DialectBigQuery)
		if err != nil {
			t.Fatalf("GetQuerySpan: %v", err)
		}
		if got := resultNames(span); len(got) != 1 || got[0] != "x" {
			t.Errorf("results = %v, want [x]", got)
		}
		want := []string{"t", "u"}
		if got := tableNames(span); len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
			t.Errorf("tables = %v, want %v", got, want)
		}
	})

	t.Run("a suffix-only statement fails closed", func(t *testing.T) {
		if _, err := GetQuerySpan("|> JOIN u ON true", DialectBigQuery); err == nil {
			t.Fatalf("GetQuerySpan accepted a suffix-only statement; want fail closed")
		}
	})

	t.Run("a statement followed by a suffix fails closed", func(t *testing.T) {
		if _, err := GetQuerySpan("DESCRIBE t |> AGGREGATE", DialectBigQuery); err == nil {
			t.Fatalf("GetQuerySpan accepted a piped DESCRIBE; want fail closed")
		}
	})

	t.Run("a pipe INSERT is DML", func(t *testing.T) {
		if got := ClassifySQL("|> INSERT t (x)", DialectBigQuery); got != DML {
			t.Errorf("ClassifySQL = %v, want DML", got)
		}
	})
}
