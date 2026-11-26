package analysis

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/metadata"
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
		if got := ClassifySQL(tc.sql); got != tc.want {
			t.Errorf("ClassifySQL(%q) = %v, want %v", tc.sql, got, tc.want)
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
		span, err := GetQuerySpan("SELECT x FROM t |> WHERE x > 1 |> ORDER BY x |> LIMIT 5")
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
		span, err := GetQuerySpan(sql)
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

	t.Run("a projection-changing suffix over concrete columns resolves", func(t *testing.T) {
		for sql, want := range map[string][]string{
			"SELECT secret FROM t |> EXTEND secret AS copy": {"secret", "copy"},
			"SELECT x FROM t |> SELECT x AS y":              {"y"},
			"SELECT x, y FROM t |> DROP x":                  {"y"},
			"SELECT x FROM t |> AGGREGATE COUNT(*)":         {""},
		} {
			span, err := GetQuerySpan(sql)
			if err != nil {
				t.Errorf("GetQuerySpan(%q): %v", sql, err)
				continue
			}
			if got := resultNames(span); !eqStrings(got, want) {
				t.Errorf("GetQuerySpan(%q) results = %q, want %q", sql, got, want)
			}
		}
	})

	t.Run("a projection-changing suffix it cannot place fails closed", func(t *testing.T) {
		for _, sql := range []string{
			// A base table's columns are unknown without metadata.
			"FROM t |> DROP x",
			"FROM t |> RENAME x AS y",
			"SELECT x FROM t |> JOIN u ON true",
			"SELECT * EXCEPT (x) FROM t |> SELECT y",
			// Operators whose output omni does not derive.
			"SELECT x FROM t |> PIVOT (SUM(x) FOR x IN (1, 2))",
			"SELECT x FROM t |> CALL f()",
			"SELECT x FROM t |> NATURAL JOIN (SELECT 1 AS x)",
			"SELECT x, y FROM t |> SET x = 1, x = 2 |> SELECT x |> RENAME z AS w",
			// Name matching that keeps only the shared columns.
			"SELECT x, y FROM t |> INNER UNION ALL BY NAME (SELECT x FROM u)",
			"SELECT x, y FROM t |> UNION ALL CORRESPONDING (SELECT x FROM u)",
			// A PIVOT source's output is not its input's columns.
			"SELECT 0 AS p |> CROSS JOIN (SELECT s, c FROM t) PIVOT(MAX(s) FOR c IN ('a', 'b'))",
			"FROM (SELECT s, c FROM t) PIVOT(MAX(s) FOR c IN ('a', 'b')) |> SELECT *",
		} {
			if _, err := GetQuerySpan(sql); err == nil {
				t.Errorf("GetQuerySpan(%q) returned a span; want fail closed", sql)
			}
		}
	})

	t.Run("a suffix-only statement fails closed", func(t *testing.T) {
		if _, err := GetQuerySpan("|> JOIN u ON true"); err == nil {
			t.Fatalf("GetQuerySpan accepted a suffix-only statement; want fail closed")
		}
	})

	t.Run("a statement followed by a suffix fails closed", func(t *testing.T) {
		if _, err := GetQuerySpan("DESCRIBE t |> AGGREGATE"); err == nil {
			t.Fatalf("GetQuerySpan accepted a piped DESCRIBE; want fail closed")
		}
	})

	t.Run("a pipe INSERT is DML", func(t *testing.T) {
		if got := ClassifySQL("|> INSERT t (x)"); got != DML {
			t.Errorf("ClassifySQL = %v, want DML", got)
		}
	})

	t.Run("a piped CALL keeps the head's unknown classification", func(t *testing.T) {
		if got := ClassifySQL("CALL p() |> WHERE true"); got != Unknown {
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

// BigQuery's dry run states no column lineage, so span.tsv holds only the
// arity of a pipe's results. These pin the lineage each operator gives, as
// BigQuery evaluates it on inline tables: a range variable reads the original
// row after SET, DROP and RENAME, and SET's values read the input.
func TestResolveQuerySpanPipeLineage(t *testing.T) {
	schemas := pinTables(
		&metadata.TableMetadata{Name: "t", Columns: []*metadata.ColumnMetadata{{Name: "x"}, {Name: "y"}, {Name: "secret"}}},
		&metadata.TableMetadata{Name: "u", Columns: []*metadata.ColumnMetadata{{Name: "x"}, {Name: "z"}}},
	)
	tests := []struct {
		sql  string
		want [][]string // each result's name, then its sources
	}{
		{"FROM t |> SET x = secret, secret = x", [][]string{{"x", "t.secret"}, {"y", "t.y"}, {"secret", "t.x"}}},
		{"FROM t AS a |> SET x = secret |> SELECT a.x, x", [][]string{{"x", "t.x"}, {"x_1", "t.secret"}}},
		{"FROM t AS a |> DROP secret |> SELECT a.secret, *", [][]string{{"secret", "t.secret"}, {"x", "t.x"}, {"y", "t.y"}}},
		{"FROM t AS a |> RENAME x AS w |> SELECT a.x, w", [][]string{{"x", "t.x"}, {"w", "t.x"}}},
		{"FROM t |> SET x = secret |> AS a |> SELECT a.x", [][]string{{"x", "t.secret"}}},
		{"FROM t |> EXTEND secret AS s |> DROP secret |> RENAME s AS secret", [][]string{{"x", "t.x"}, {"y", "t.y"}, {"secret", "t.secret"}}},
		{"FROM t |> JOIN u USING (x)", [][]string{{"x", "t.x", "u.x"}, {"y", "t.y"}, {"secret", "t.secret"}, {"z", "u.z"}}},
		{"FROM t |> JOIN u ON t.x = u.x |> SELECT u.x, t.x", [][]string{{"x", "u.x"}, {"x_1", "t.x"}}},
		{"FROM t |> AGGREGATE MAX(secret) AS m GROUP BY y", [][]string{{"y", "t.y"}, {"m", "t.secret"}}},
		{"FROM t |> AGGREGATE COUNT(*) GROUP BY ROLLUP(y, x)", [][]string{{"y", "t.y"}, {"x", "t.x"}, {"f0_"}}},
		{"FROM t |> SELECT x |> UNION ALL (SELECT z FROM u)", [][]string{{"x", "t.x", "u.z"}}},
		{"FROM t |> SELECT x, y |> UNION ALL BY NAME (SELECT z AS y, x FROM u)", [][]string{{"x", "t.x", "u.x"}, {"y", "t.y", "u.z"}}},
		{"FROM t |> SELECT * REPLACE (secret AS x) |> SELECT x", [][]string{{"x", "t.secret"}}},
		// A FROM-first query's range variables survive parentheses and an
		// unaliased FROM subquery.
		{"(FROM t AS a) |> SELECT a.secret", [][]string{{"secret", "t.secret"}}},
		{"(FROM t AS a |> AS b) |> SELECT b.secret", [][]string{{"secret", "t.secret"}}},
		{"SELECT a.secret FROM (FROM t AS a)", [][]string{{"secret", "t.secret"}}},
		{"FROM (FROM t AS a |> SET secret = x) |> SELECT a.secret, secret", [][]string{{"secret", "t.secret"}, {"secret_1", "t.x"}}},
	}
	for _, tc := range tests {
		span := resolvePin(t, bigqueryResolveOptions(), tc.sql, "ds1", schemas)
		var got [][]string
		for _, r := range span.Results {
			got = append(got, append([]string{r.Name}, sourcesOf(r)...))
		}
		require.Equal(t, tc.want, got, "statement: %s", tc.sql)
	}
}
