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
		// A CALL argument that names a table is an access, so a system table
		// beside it does not make the query a system-only read.
		{"SELECT * FROM INFORMATION_SCHEMA.TABLES |> CALL f(TABLE secret)", Select},
		{"SELECT * FROM INFORMATION_SCHEMA.TABLES |> UNION ALL TABLE f(TABLE secret)", Select},
		{"|> AGGREGATE", Select},
		{"|> WHERE x > 1 |> DISTINCT", Select},
		// A piped statement keeps its head's classification: DESCRIBE and CALL
		// are Unknown, and a reading suffix does not turn them into a read.
		{"DESCRIBE t |> AGGREGATE COUNT(*)", Unknown},
		{"|> INSERT t", DML},
		{"SELECT 1 |> INSERT t", DML},
		{"|> CREATE TABLE t (x INT64)", DDL},
		{"SELECT 1 |> IF true THEN (|> INSERT t) ELSE (|> DISTINCT)", DML},
		{"|> TEE (|> INSERT t), (|> DISTINCT)", DML},
		// Conflicting write types on different runtime paths are Unknown.
		{"|> IF true THEN (|> INSERT t) ELSE (|> CREATE TABLE u (x INT64))", Unknown},
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
		{"|> IF true THEN (|> INSERT t) ELSE (|> CREATE TABLE u (x INT64))", StatementTypeUnspecified},
		{"|> IF true THEN (|> INSERT t) ELSEIF false THEN (|> CREATE TABLE u (x INT64)) ELSE (|> INSERT v)", StatementTypeUnspecified},
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
	t.Run("a row-only suffix keeps the projection", func(t *testing.T) {
		span, err := GetQuerySpan("SELECT x FROM t |> WHERE x > 1 |> ORDER BY x |> LIMIT 5", DialectBigQuery)
		if err != nil {
			t.Fatalf("GetQuerySpan: %v", err)
		}
		if got := resultNames(span); len(got) != 1 || got[0] != "x" {
			t.Errorf("results = %v, want [x]", got)
		}
		if got := tableNames(span); len(got) != 1 || got[0] != "t" {
			t.Errorf("tables = %v, want [t]", got)
		}
	})

	t.Run("a sample's embedded query is walked", func(t *testing.T) {
		sql := "SELECT * FROM INFORMATION_SCHEMA.TABLES |> TABLESAMPLE SYSTEM (10 PERCENT PARTITION BY (SELECT x FROM secret))"
		span, err := GetQuerySpan(sql, DialectBigQuery)
		if err != nil {
			t.Fatalf("GetQuerySpan: %v", err)
		}
		found := false
		for _, name := range tableNames(span) {
			if name == "secret" {
				found = true
			}
		}
		if !found {
			t.Errorf("tables = %v, want secret", tableNames(span))
		}
	})

	t.Run("a projection-changing suffix fails closed", func(t *testing.T) {
		for _, sql := range []string{
			"SELECT secret FROM t |> EXTEND secret AS copy",
			"SELECT x FROM t |> JOIN u ON true",
			"SELECT x FROM t |> SELECT x AS y",
			"SELECT x FROM t |> DROP x",
			"SELECT x FROM t |> AGGREGATE COUNT(*)",
		} {
			if _, err := GetQuerySpan(sql, DialectBigQuery); err == nil {
				t.Errorf("GetQuerySpan(%q) returned a span; want fail closed", sql)
			}
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

	t.Run("a piped CALL keeps the head's unknown classification", func(t *testing.T) {
		if got := ClassifySQL("CALL p() |> WHERE true", DialectBigQuery); got != Unknown {
			t.Errorf("ClassifySQL = %v, want Unknown", got)
		}
		file, err := parseFile("CALL p() |> WHERE true")
		if err != nil {
			t.Fatalf("parse: %v", err)
		}
		if got := StatementTypeOf(file.Stmts[0]); got != Execute {
			t.Errorf("StatementTypeOf = %v, want Execute", got)
		}
	})
}
