package analysis_test

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/snowflake/analysis"
)

// testCatalog resolves "DATABASE.SCHEMA.TABLE" keys to column lists.
type testCatalog map[string][]string

var errTableNotFound = errors.New("table not found")

func (c testCatalog) resolve(database, schema, table string) (*analysis.Table, error) {
	columns, ok := c[database+"."+schema+"."+table]
	if !ok {
		return nil, fmt.Errorf("%s.%s.%s: %w", database, schema, table, errTableNotFound)
	}
	return &analysis.Table{Name: table, Columns: columns}, nil
}

var catalog = testCatalog{
	"DB.PUBLIC.T1":    {"A", "B", "C", "D"},
	"DB.PUBLIC.T2":    {"A", "E", "C", "D"},
	"DB.PUBLIC.T3":    {"E", "F"},
	"DB.OTHER.T1":     {"Z"},
	"DB.PUBLIC.lower": {"a", "A"},
}

func resolveSpan(t *testing.T, sql string) (*analysis.ResolvedSpan, error) {
	t.Helper()
	return analysis.ResolveQuerySpan(sql, analysis.ResolveOptions{
		DefaultDatabase: "DB",
		ResolveTable:    catalog.resolve,
	})
}

// render prints each result as NAME<-sources, with "plain" for a plain field.
func render(span *analysis.ResolvedSpan) []string {
	var out []string
	for _, r := range span.Results {
		var sources []string
		for c := range r.SourceColumns {
			sources = append(sources, strings.Join([]string{c.Database, c.Schema, c.Table, c.Column}, "."))
		}
		sort.Strings(sources)
		s := r.Name + "<-" + strings.Join(sources, ",")
		if r.IsPlainField {
			s += " plain"
		}
		out = append(out, s)
	}
	return out
}

func TestResolveQuerySpanResults(t *testing.T) {
	tests := []struct {
		sql  string
		want []string
	}{
		{"SELECT A AS X, T.D FROM T1 T", []string{"X<-DB.PUBLIC.T1.A plain", "D<-DB.PUBLIC.T1.D plain"}},
		{"SELECT A + B AS S FROM T1", []string{"S<-DB.PUBLIC.T1.A,DB.PUBLIC.T1.B"}},
		// An unaliased expression is named by its source text from the node's
		// Loc, which for a binary expression starts at the operator.
		{"SELECT A + B FROM T1", []string{"+ B<-DB.PUBLIC.T1.A,DB.PUBLIC.T1.B"}},
		{"SELECT UPPER(A) FROM T1", []string{"UPPER(A)<-DB.PUBLIC.T1.A"}},
		{"SELECT * FROM T3", []string{"E<-DB.PUBLIC.T3.E plain", "F<-DB.PUBLIC.T3.F plain"}},
		{"SELECT * FROM OTHER.T1", []string{"Z<-DB.OTHER.T1.Z plain"}},
		{`SELECT "a" FROM "lower"`, []string{"a<-DB.PUBLIC.lower.a plain"}},
		// Star transforms apply in the order ILIKE, EXCLUDE, REPLACE, RENAME.
		{"SELECT * ILIKE '_' EXCLUDE (B) REPLACE (E AS A) RENAME (C AS X) FROM T1, T3", []string{"A<-DB.PUBLIC.T3.E", "X<-DB.PUBLIC.T1.C plain", "D<-DB.PUBLIC.T1.D plain", "E<-DB.PUBLIC.T3.E plain", "F<-DB.PUBLIC.T3.F plain"}},
		// CASE branches and window clauses contribute lineage.
		{"SELECT CASE WHEN A > 0 THEN B END AS X FROM T1", []string{"X<-DB.PUBLIC.T1.A,DB.PUBLIC.T1.B"}},
		{"SELECT ROW_NUMBER() OVER (PARTITION BY A ORDER BY B) AS N FROM T1", []string{"N<-DB.PUBLIC.T1.A,DB.PUBLIC.T1.B"}},
		{"SELECT (SELECT MAX(E) FROM T3) AS M FROM T1", []string{"M<-DB.PUBLIC.T3.E"}},
		// Positional references.
		{"SELECT $2 FROM T3", []string{"F<-DB.PUBLIC.T3.F plain"}},
		{"SELECT T.$1 FROM T1, T3 T", []string{"E<-DB.PUBLIC.T3.E plain"}},
		// NATURAL and USING joins emit each shared column once with both sides' lineage.
		{"SELECT * FROM T1 NATURAL JOIN T2", []string{"A<-DB.PUBLIC.T1.A,DB.PUBLIC.T2.A plain", "B<-DB.PUBLIC.T1.B plain", "C<-DB.PUBLIC.T1.C,DB.PUBLIC.T2.C plain", "D<-DB.PUBLIC.T1.D,DB.PUBLIC.T2.D plain", "E<-DB.PUBLIC.T2.E plain"}},
		{"SELECT * FROM T1 JOIN T3 USING (E)", []string{"E<-DB.PUBLIC.T3.E plain", "A<-DB.PUBLIC.T1.A plain", "B<-DB.PUBLIC.T1.B plain", "C<-DB.PUBLIC.T1.C plain", "D<-DB.PUBLIC.T1.D plain", "F<-DB.PUBLIC.T3.F plain"}},
		{"SELECT T3.F FROM T1 JOIN T3 ON TRUE", []string{"F<-DB.PUBLIC.T3.F plain"}},
		// Set operations keep the left names; BY NAME matches names.
		{"SELECT A FROM T1 UNION SELECT E FROM T3", []string{"A<-DB.PUBLIC.T1.A,DB.PUBLIC.T3.E plain"}},
		{"SELECT A, B FROM T1 UNION ALL BY NAME SELECT E AS B, F FROM T3", []string{"A<-DB.PUBLIC.T1.A plain", "B<-DB.PUBLIC.T1.B,DB.PUBLIC.T3.E plain", "F<-DB.PUBLIC.T3.F plain"}},
		// CTEs, column lists, shadowing and recursion.
		{"WITH C (X) AS (SELECT A FROM T1) SELECT * FROM C", []string{"X<-DB.PUBLIC.T1.A plain"}},
		{"WITH T1 AS (SELECT E FROM T3) SELECT * FROM T1", []string{"E<-DB.PUBLIC.T3.E plain"}},
		{"WITH RECURSIVE R (X, Y) AS (SELECT A, B FROM T1 UNION ALL SELECT Y, X FROM R) SELECT * FROM R", []string{"X<-DB.PUBLIC.T1.A,DB.PUBLIC.T1.B plain", "Y<-DB.PUBLIC.T1.A,DB.PUBLIC.T1.B plain"}},
		// A table function exposes no columns.
		{"SELECT 1 AS ONE FROM TABLE(FLATTEN(input => PARSE_JSON('[1]')))", []string{"ONE<-"}},
	}
	for _, tc := range tests {
		span, err := resolveSpan(t, tc.sql)
		if !assert.NoErrorf(t, err, "%s", tc.sql) {
			continue
		}
		assert.Equalf(t, analysis.Select, span.Type, "%s", tc.sql)
		assert.Equalf(t, tc.want, render(span), "%s", tc.sql)
	}
}

func TestResolveQuerySpanAccessTables(t *testing.T) {
	span, err := resolveSpan(t, "WITH C AS (SELECT A FROM T1) SELECT * FROM C WHERE EXISTS (SELECT 1 FROM OTHER.T1)")
	require.NoError(t, err)
	require.Equal(t, analysis.ColumnSet{
		{Database: "DB", Schema: "PUBLIC", Table: "T1"}: true,
		{Database: "DB", Schema: "PUBLIC", Table: "C"}:  true,
		{Database: "DB", Schema: "OTHER", Table: "T1"}:  true,
	}, span.SourceColumns)
}

func TestResolveQuerySpanNonQuery(t *testing.T) {
	tests := []struct {
		sql  string
		want analysis.QueryType
	}{
		{"DELETE FROM T1 WHERE A IN (SELECT E FROM T3)", analysis.DML},
		{"CREATE TABLE X (A INT)", analysis.DDL},
		{"SHOW TABLES", analysis.SelectInfoSchema},
		{"EXPLAIN SELECT * FROM NOPE", analysis.Explain},
		{"TRUNCATE TABLE T1", analysis.QueryTypeUnknown},
		// USE and SET are Select with no results.
		{"USE SCHEMA OTHER", analysis.Select},
		{"SET X = 1", analysis.Select},
	}
	for _, tc := range tests {
		span, err := resolveSpan(t, tc.sql)
		require.NoErrorf(t, err, "%s", tc.sql)
		require.Equalf(t, tc.want, span.Type, "%s", tc.sql)
		require.Emptyf(t, span.Results, "%s", tc.sql)
		require.NotNilf(t, span.Results, "%s", tc.sql)
	}
	// The tables a statement reads are the TableRefs under it, so a DELETE's
	// target is not among them.
	span, err := resolveSpan(t, "DELETE FROM T1 WHERE A IN (SELECT E FROM T3)")
	require.NoError(t, err)
	require.Equal(t, analysis.ColumnSet{
		{Database: "DB", Schema: "PUBLIC", Table: "T3"}: true,
	}, span.SourceColumns)
}

func TestResolveQuerySpanErrors(t *testing.T) {
	// A resolver error comes back inside a ResolveError with the accessed tables.
	_, err := resolveSpan(t, "SELECT A FROM NOPE")
	var resolveErr *analysis.ResolveError
	require.ErrorAs(t, err, &resolveErr)
	require.ErrorIs(t, err, errTableNotFound)
	require.Equal(t, analysis.ColumnSet{{Database: "DB", Schema: "PUBLIC", Table: "NOPE"}: true}, resolveErr.SourceColumns)
	require.Equal(t, "DB.PUBLIC.NOPE: table not found", err.Error())

	for sql, want := range map[string]string{
		"SELECT NOPE FROM T1":                                                       `failed to check whether the column "NOPE" is sensitive: no matching column`,
		"SELECT $5 FROM T3":                                                         "column position $5 is invalid: the FROM clause only returns 2 columns",
		"SELECT A FROM T1 UNION SELECT * FROM T3":                                   "the number of columns in the left part of the set operation returns 1 fields, but the right part returns 2 fields",
		"SELECT A FROM T1 UNION SELECT NOPE FROM T3":                                "failed to extract the right part of the set operation",
		"WITH C (X, Y) AS (SELECT A FROM T1) SELECT * FROM C":                       `the number of columns in the CTE "C" returns 1 fields, but the column list returns 2 fields`,
		"SELECT * FROM T1 PIVOT (SUM(A) FOR B IN ('x')) AS P":                       "PIVOT/UNPIVOT table sources are not supported",
		"SHOW TABLES ->> SELECT * FROM $1":                                          "result-pipe (->>) statements are not supported",
		"WITH R AS (SELECT A FROM T1 UNION ALL SELECT A, B FROM R) SELECT * FROM R": "recursive clause",
	} {
		_, err := resolveSpan(t, sql)
		require.ErrorAsf(t, err, &resolveErr, "%s", sql)
		require.ErrorContainsf(t, err, want, "%s", sql)
	}

	_, err = resolveSpan(t, "SELECT 1; SELECT 2")
	require.EqualError(t, err, "expected exactly 1 statement, got 2")
	_, err = resolveSpan(t, "SELEC 1")
	require.Error(t, err)
	require.NotErrorAs(t, err, &resolveErr)

	span, err := resolveSpan(t, "-- nothing")
	require.NoError(t, err)
	require.Nil(t, span)
}
