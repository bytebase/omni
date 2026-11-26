package analysis

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/metadata"
)

func TestQuerySpanScaffold_QueryTypesAndAccessTables(t *testing.T) {
	tests := []struct {
		name          string
		statement     string
		wantType      QueryType
		wantResources []ColumnResource
	}{
		{
			name:      "explain_select_accesses_table",
			statement: "EXPLAIN SELECT * FROM t",
			wantType:  Explain,
			wantResources: []ColumnResource{
				{Database: "db", Table: "t"},
			},
		},
		{
			name:      "set_is_select_with_no_access_tables",
			statement: "SET CHARSET DEFAULT",
			wantType:  Select,
		},
		{
			name:      "show_is_select_info_schema",
			statement: "SHOW DATABASES",
			wantType:  SelectInfoSchema,
		},
		{
			name:      "plain_explain_keeps_the_tables_it_plans",
			statement: "EXPLAIN SELECT * FROM t",
			wantType:  Explain,
			wantResources: []ColumnResource{
				{Database: "db", Table: "t"},
			},
		},
		{
			name:      "describe_is_info_schema",
			statement: "DESCRIBE t",
			wantType:  SelectInfoSchema,
		},
		{
			name:      "desc_tab_is_info_schema",
			statement: "DESC\t t",
			wantType:  SelectInfoSchema,
		},
		{
			name:      "desc_newline_is_info_schema",
			statement: "DESC\n t",
			wantType:  SelectInfoSchema,
		},
		{
			name:          "ddl",
			statement:     "CREATE TABLE t(a INT)",
			wantType:      DDL,
			wantResources: []ColumnResource{{Database: "db", Table: "t"}},
		},
		{
			name:          "dml",
			statement:     "INSERT INTO t VALUES(1)",
			wantType:      DML,
			wantResources: []ColumnResource{{Database: "db", Table: "t"}},
		},
		{
			name:      "system_select",
			statement: "SELECT * FROM mysql.user",
			wantType:  SelectInfoSchema,
		},
		{
			name:      "standard_select_accesses_table",
			statement: "SELECT * FROM t WHERE a > 0",
			wantType:  Select,
			wantResources: []ColumnResource{
				{Database: "db", Table: "t"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			span, err := resolveExtractorSpan(newExtractorTestCatalog(), tc.statement, "db", false)
			require.NoError(t, err)
			require.Equal(t, tc.wantType, span.Type)
			require.Equal(t, sourceColumnSetFromResources(tc.wantResources), span.SourceColumns)
			if tc.wantType != Select {
				require.Empty(t, span.Results)
			}
			require.Empty(t, span.PredicateColumns)
		})
	}
}

func TestQuerySpanPhase1_SimpleSelectResultColumns(t *testing.T) {
	tests := []struct {
		name      string
		statement string
		want      []QuerySpanResult
	}{
		{
			name:      "constant",
			statement: "SELECT 1",
			want: []QuerySpanResult{
				{Name: "1", SourceColumns: ColumnSet{}, IsPlainField: true},
			},
		},
		{
			name:      "bare_column",
			statement: "SELECT a FROM t",
			want: []QuerySpanResult{
				{Name: "a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}}), IsPlainField: true},
			},
		},
		{
			name:      "alias",
			statement: "SELECT a AS x FROM t",
			want: []QuerySpanResult{
				{Name: "x", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}}), IsPlainField: true},
			},
		},
		{
			name:      "qualified_columns",
			statement: "SELECT a, t.b, db.t.c FROM t",
			want: []QuerySpanResult{
				{Name: "a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}}), IsPlainField: true},
				{Name: "b", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "b"}}), IsPlainField: true},
				{Name: "c", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "c"}}), IsPlainField: true},
			},
		},
		{
			name:      "star",
			statement: "SELECT * FROM t",
			want: []QuerySpanResult{
				{Name: "a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}}), IsPlainField: true},
				{Name: "b", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "b"}}), IsPlainField: true},
				{Name: "c", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "c"}}), IsPlainField: true},
			},
		},
		{
			name:      "star_and_column",
			statement: "SELECT *, a FROM t",
			want: []QuerySpanResult{
				{Name: "a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}}), IsPlainField: true},
				{Name: "b", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "b"}}), IsPlainField: true},
				{Name: "c", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "c"}}), IsPlainField: true},
				{Name: "a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}}), IsPlainField: true},
			},
		},
		{
			name:      "table_star",
			statement: "SELECT t.* FROM t",
			want: []QuerySpanResult{
				{Name: "a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}}), IsPlainField: true},
				{Name: "b", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "b"}}), IsPlainField: true},
				{Name: "c", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "c"}}), IsPlainField: true},
			},
		},
	}

	catalog := newExtractorTestCatalog()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			span, err := resolveExtractorSpan(catalog, tc.statement, "db", false)
			require.NoError(t, err)
			require.Equal(t, Select, span.Type)
			require.Equal(t, tc.want, span.Results)
		})
	}
}

func TestQuerySpanPhase2_ExpressionSourceMerging(t *testing.T) {
	span, err := resolveExtractorSpan(newExtractorTestCatalog(),
		"SELECT max(a), a-b AS c1, a=b AS c2, a>b, b in (a, c) FROM t",
		"db",
		false,
	)
	require.NoError(t, err)
	require.Equal(t, []QuerySpanResult{
		{
			Name:          "max(a)",
			SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}}),
			IsPlainField:  false,
		},
		{
			Name: "c1",
			SourceColumns: sourceColumnSetFromResources([]ColumnResource{
				{Database: "db", Table: "t", Column: "a"},
				{Database: "db", Table: "t", Column: "b"},
			}),
			IsPlainField: false,
		},
		{
			Name: "c2",
			SourceColumns: sourceColumnSetFromResources([]ColumnResource{
				{Database: "db", Table: "t", Column: "a"},
				{Database: "db", Table: "t", Column: "b"},
			}),
			IsPlainField: false,
		},
		{
			Name: "a>b",
			SourceColumns: sourceColumnSetFromResources([]ColumnResource{
				{Database: "db", Table: "t", Column: "a"},
				{Database: "db", Table: "t", Column: "b"},
			}),
			IsPlainField: false,
		},
		{
			Name: "b in (a, c)",
			SourceColumns: sourceColumnSetFromResources([]ColumnResource{
				{Database: "db", Table: "t", Column: "a"},
				{Database: "db", Table: "t", Column: "b"},
				{Database: "db", Table: "t", Column: "c"},
			}),
			IsPlainField: false,
		},
	}, span.Results)
}

// TestQuerySpanKeywordArgAndWeightString covers the omni AST nodes
// introduced with the $-token bump: KeywordArg (the unit selector in
// TIMESTAMPDIFF/TIMESTAMPADD/GET_FORMAT — a keyword, not a column, so no
// lineage) and WeightStringExpr (lineage follows the real operand).
// TIMESTAMPDIFF regressed to "unsupported omni MySQL expression" without
// these cases.
func TestQuerySpanKeywordArgAndWeightString(t *testing.T) {
	span, err := resolveExtractorSpan(newExtractorTestCatalog(),
		"SELECT TIMESTAMPDIFF(SECOND, a, b) AS d1, TIMESTAMPADD(MINUTE, 5, a) AS d2, GET_FORMAT(DATE, 'ISO') AS d3, WEIGHT_STRING(c AS CHAR(10)) AS w FROM t",
		"db",
		false,
	)
	require.NoError(t, err)
	require.Equal(t, []QuerySpanResult{
		{
			Name: "d1",
			SourceColumns: sourceColumnSetFromResources([]ColumnResource{
				{Database: "db", Table: "t", Column: "a"},
				{Database: "db", Table: "t", Column: "b"},
			}),
			IsPlainField: false,
		},
		{
			Name: "d2",
			SourceColumns: sourceColumnSetFromResources([]ColumnResource{
				{Database: "db", Table: "t", Column: "a"},
			}),
			IsPlainField: false,
		},
		{
			Name:          "d3",
			SourceColumns: ColumnSet{},
			IsPlainField:  false,
		},
		{
			Name: "w",
			SourceColumns: sourceColumnSetFromResources([]ColumnResource{
				{Database: "db", Table: "t", Column: "c"},
			}),
			IsPlainField: false,
		},
	}, span.Results)
}

func TestQuerySpanPhase3_FromJoinAliasAndScope(t *testing.T) {
	tests := []struct {
		name      string
		statement string
		want      []QuerySpanResult
	}{
		{
			name:      "table_alias",
			statement: "SELECT x.a FROM t AS x",
			want: []QuerySpanResult{
				{Name: "a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}}), IsPlainField: true},
			},
		},
		{
			name:      "join_on_with_qualified_star_and_column",
			statement: "SELECT t1.*, t2.c, 0 FROM t1 JOIN t2 ON 1 = 1",
			want: []QuerySpanResult{
				{Name: "a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t1", Column: "a"}}), IsPlainField: true},
				{Name: "b", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t1", Column: "b"}}), IsPlainField: true},
				{Name: "c", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t1", Column: "c"}}), IsPlainField: true},
				{Name: "c", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t2", Column: "c"}}), IsPlainField: true},
				{Name: "0", SourceColumns: ColumnSet{}, IsPlainField: true},
			},
		},
		{
			name:      "join_using_merges_star_columns",
			statement: "SELECT * FROM t AS t1 JOIN t AS t2 USING(a)",
			want: []QuerySpanResult{
				{Name: "a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}}), IsPlainField: true},
				{Name: "b", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "b"}}), IsPlainField: true},
				{Name: "c", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "c"}}), IsPlainField: true},
				{Name: "b", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "b"}}), IsPlainField: true},
				{Name: "c", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "c"}}), IsPlainField: true},
			},
		},
	}

	catalog := newExtractorTestCatalog()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			span, err := resolveExtractorSpan(catalog, tc.statement, "db", false)
			require.NoError(t, err)
			require.Equal(t, tc.want, span.Results)
		})
	}
}

func TestQuerySpanPhase4_SubqueryExpressions(t *testing.T) {
	catalog := newExtractorTestCatalog()

	t.Run("scalar_subquery_results", func(t *testing.T) {
		span, err := resolveExtractorSpan(
			catalog,
			"SELECT 1 AS col_1, (SELECT(2)) AS col_2, (SELECT AVG(a + b * c) FROM t) AS avg_a_b_c FROM t",
			"db",
			false,
		)
		require.NoError(t, err)
		require.Equal(t, []QuerySpanResult{
			{Name: "col_1", SourceColumns: ColumnSet{}, IsPlainField: true},
			{Name: "col_2", SourceColumns: ColumnSet{}, IsPlainField: false},
			{
				Name: "avg_a_b_c",
				SourceColumns: sourceColumnSetFromResources([]ColumnResource{
					{Database: "db", Table: "t", Column: "a"},
					{Database: "db", Table: "t", Column: "b"},
					{Database: "db", Table: "t", Column: "c"},
				}),
				IsPlainField: false,
			},
		}, span.Results)
		require.Equal(t, sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t"}}), span.SourceColumns)
	})

	t.Run("correlated_subquery_access_tables", func(t *testing.T) {
		span, err := resolveExtractorSpan(
			catalog,
			"SELECT city, (SELECT COUNT(*) FROM paintings p WHERE g.id = p.gallery_id) AS total_paintings FROM galleries g",
			"db",
			false,
		)
		require.NoError(t, err)
		require.Equal(t, []QuerySpanResult{
			{Name: "city", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "galleries", Column: "city"}}), IsPlainField: true},
			{Name: "total_paintings", SourceColumns: ColumnSet{}, IsPlainField: false},
		}, span.Results)
		require.Equal(t, sourceColumnSetFromResources([]ColumnResource{
			{Database: "db", Table: "galleries"},
			{Database: "db", Table: "paintings"},
		}), span.SourceColumns)
	})
}

func TestQuerySpanPhase5And6_CTEAndSetOperations(t *testing.T) {
	tests := []struct {
		name      string
		statement string
		want      []QuerySpanResult
	}{
		{
			name:      "simple_cte",
			statement: "WITH t1 AS (SELECT * FROM t) SELECT * FROM t1",
			want: []QuerySpanResult{
				{Name: "a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}}), IsPlainField: true},
				{Name: "b", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "b"}}), IsPlainField: true},
				{Name: "c", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "c"}}), IsPlainField: true},
			},
		},
		{
			name:      "cte_column_aliases",
			statement: "WITH t1(x, y, z) AS (SELECT * FROM t) SELECT * FROM t1",
			want: []QuerySpanResult{
				{Name: "x", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}}), IsPlainField: true},
				{Name: "y", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "b"}}), IsPlainField: true},
				{Name: "z", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "c"}}), IsPlainField: true},
			},
		},
		{
			name:      "union_merges_positionally",
			statement: "SELECT 1 AS c1 UNION SELECT a FROM t",
			want: []QuerySpanResult{
				{Name: "c1", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}}), IsPlainField: false},
			},
		},
		{
			name:      "recursive_keyword_non_recursive_cte",
			statement: "WITH RECURSIVE t1 AS (SELECT 1 AS c1 UNION SELECT a FROM t) SELECT * FROM t1",
			want: []QuerySpanResult{
				{Name: "c1", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}}), IsPlainField: false},
			},
		},
	}

	catalog := newExtractorTestCatalog()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			span, err := resolveExtractorSpan(catalog, tc.statement, "db", false)
			require.NoError(t, err)
			require.Equal(t, tc.want, span.Results)
		})
	}
}

func TestQuerySpanPhase7_JSONTable(t *testing.T) {
	span, err := resolveExtractorSpan(newExtractorTestCatalog(),
		`SELECT *
FROM products,
JSON_TABLE(product_info, '$' COLUMNS (
  product_id INT PATH '$.id',
  product_name VARCHAR(50) PATH '$.name'
)) AS jt`,
		"db",
		false,
	)
	require.NoError(t, err)
	require.Equal(t, []QuerySpanResult{
		{Name: "product_info", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "products", Column: "product_info"}}), IsPlainField: true},
		{Name: "product_id", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "products", Column: "product_info"}}), IsPlainField: true},
		{Name: "product_name", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "products", Column: "product_info"}}), IsPlainField: true},
	}, span.Results)
}

func TestQuerySpanPhase9_ResourceNotFoundFailOpen(t *testing.T) {
	dbMeta := &metadata.DatabaseSchemaMetadata{Name: "db"}
	span, err := resolveExtractorSpan(newExtractorCatalog([]*metadata.DatabaseSchemaMetadata{dbMeta}), "SELECT * FROM t WHERE a > 0", "db", false)
	require.NoError(t, err)
	require.Equal(t, Select, span.Type)
	require.Equal(t, sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t"}}), span.SourceColumns)
	require.Empty(t, span.Results)
	require.ErrorAs(t, span.NotFoundError, new(*NotFoundError))
}

func TestQuerySpanSystematicMigrationRegressions(t *testing.T) {
	catalog := newExtractorTestCatalog()

	t.Run("derived_table_column_aliases_are_applied", func(t *testing.T) {
		span, err := resolveExtractorSpan(
			catalog,
			"SELECT x.c1 FROM (SELECT a, b FROM t) AS x(c1, c2)",
			"db",
			false,
		)
		require.NoError(t, err)
		require.Equal(t, []QuerySpanResult{
			{Name: "c1", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}}), IsPlainField: true},
		}, span.Results)
	})

	t.Run("in_subquery_sources_are_part_of_result_lineage", func(t *testing.T) {
		span, err := resolveExtractorSpan(
			catalog,
			"SELECT a IN (SELECT t2.c FROM t2) AS matched FROM t1",
			"db",
			false,
		)
		require.NoError(t, err)
		require.Equal(t, []QuerySpanResult{
			{
				Name: "matched",
				SourceColumns: sourceColumnSetFromResources([]ColumnResource{
					{Database: "db", Table: "t1", Column: "a"},
					{Database: "db", Table: "t2", Column: "c"},
				}),
				IsPlainField: false,
			},
		}, span.Results)
	})

	t.Run("explicit_expression_nodes_do_not_drop_lineage", func(t *testing.T) {
		span, err := resolveExtractorSpan(
			catalog,
			`SELECT
  EXISTS (SELECT t2.a FROM t2 WHERE t2.a = t.a) AS exists_field,
  CONVERT(a USING utf8mb4) AS converted,
  a COLLATE utf8mb4_bin AS collated,
  MATCH(a) AGAINST (b) AS matched,
  ROW(a, b) AS row_field,
  a MEMBER OF(c) AS member_field,
  a + INTERVAL b DAY AS interval_field
FROM t`,
			"db",
			false,
		)
		require.NoError(t, err)
		require.Equal(t, []QuerySpanResult{
			{
				Name: "exists_field",
				SourceColumns: sourceColumnSetFromResources([]ColumnResource{
					{Database: "db", Table: "t2", Column: "a"},
				}),
				IsPlainField: false,
			},
			{Name: "converted", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}}), IsPlainField: false},
			{Name: "collated", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}}), IsPlainField: false},
			{Name: "matched", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}, {Database: "db", Table: "t", Column: "b"}}), IsPlainField: false},
			{Name: "row_field", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}, {Database: "db", Table: "t", Column: "b"}}), IsPlainField: false},
			{Name: "member_field", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}, {Database: "db", Table: "t", Column: "c"}}), IsPlainField: false},
			{Name: "interval_field", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}, {Database: "db", Table: "t", Column: "b"}}), IsPlainField: false},
		}, span.Results)
	})

	t.Run("legacy_dml_roots_stay_dml", func(t *testing.T) {
		tests := []string{
			"CALL p()",
			"DO 1",
			"HANDLER t OPEN",
			"HANDLER t READ FIRST",
			"HANDLER t CLOSE",
		}
		for _, statement := range tests {
			span, err := resolveExtractorSpan(catalog, statement, "db", false)
			require.NoError(t, err, statement)
			require.Equal(t, DML, span.Type, statement)
			require.Empty(t, span.Results, statement)
		}
	})

	// OceanBase has no TABLE statement; only VALUES is a query primary here.
	t.Run("values_root_returns_select_results", func(t *testing.T) {
		valuesSpan, err := resolveExtractorSpan(catalog, "VALUES ROW(1, 2 + 3)", "db", false)
		require.NoError(t, err)
		require.Equal(t, Select, valuesSpan.Type)
		require.Equal(t, []QuerySpanResult{
			// MySQL names a VALUES column by its position (span.tsv).
			{Name: "column_0", SourceColumns: ColumnSet{}, IsPlainField: true},
			{Name: "column_1", SourceColumns: ColumnSet{}, IsPlainField: false},
		}, valuesSpan.Results)
	})
}

func TestQuerySpanScenarioBatch1_RootAndQueryTypeCoverage(t *testing.T) {
	catalog := newExtractorTestCatalog()

	t.Run("select_family_roots", func(t *testing.T) {
		tests := []struct {
			name          string
			statement     string
			wantType      QueryType
			wantResults   []QuerySpanResult
			wantResources []ColumnResource
		}{
			{
				name:      "values_default",
				statement: "VALUES ROW(DEFAULT)",
				wantType:  Select,
				wantResults: []QuerySpanResult{
					// OceanBase names a VALUES column by its position (span.tsv).
					{Name: "column_0", SourceColumns: ColumnSet{}, IsPlainField: true},
				},
			},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				span, err := resolveExtractorSpan(catalog, tc.statement, "db", false)
				require.NoError(t, err)
				require.Equal(t, tc.wantType, span.Type)
				require.Equal(t, tc.wantResults, span.Results)
				require.Equal(t, sourceColumnSetFromResources(tc.wantResources), span.SourceColumns)
			})
		}
	})

	t.Run("multiple_statements_are_rejected", func(t *testing.T) {
		_, err := resolveExtractorSpan(catalog, "SELECT 1; SELECT 2;", "db", false)
		require.ErrorContains(t, err, "expected exactly 1 statement")
	})

	t.Run("ddl_bucket_roots", func(t *testing.T) {
		tests := []string{
			"CREATE DATABASE d",
			"CREATE VIEW v AS SELECT a FROM t",
			"ALTER TABLE t ADD COLUMN d INT",
			"DROP TABLE t",
			"RENAME TABLE t TO t_new",
			"TRUNCATE TABLE t",
		}
		for _, statement := range tests {
			span, err := resolveExtractorSpan(catalog, statement, "db", false)
			require.NoError(t, err, statement)
			require.Equal(t, DDL, span.Type, statement)
			require.Empty(t, span.Results, statement)
		}
	})

	t.Run("dml_bucket_roots", func(t *testing.T) {
		tests := []string{
			"UPDATE t SET a = 1",
			"DELETE FROM t WHERE a = 1",
			"BEGIN",
			"COMMIT",
			"ROLLBACK",
			"SAVEPOINT s",
			"LOCK TABLES t READ",
			"UNLOCK TABLES",
			"PREPARE stmt FROM 'SELECT 1'",
			"EXECUTE stmt",
			"DEALLOCATE PREPARE stmt",
		}
		for _, statement := range tests {
			span, err := resolveExtractorSpan(catalog, statement, "db", false)
			require.NoError(t, err, statement)
			require.Equal(t, DML, span.Type, statement)
			require.Empty(t, span.Results, statement)
		}
	})
}

func TestQuerySpanScenarioBatch4_QueryTypeEdges(t *testing.T) {
	catalog := newExtractorTestCatalog()

	t.Run("legacy_query_type_buckets", func(t *testing.T) {
		tests := []struct {
			statement string
			want      QueryType
		}{
			{statement: "SET PASSWORD = 'new_password'", want: QueryTypeUnknown},
			{statement: "CREATE USER u IDENTIFIED BY 'p'", want: DDL},
			{statement: "ALTER USER u IDENTIFIED BY 'p2'", want: DDL},
			{statement: "DROP USER u", want: DDL},
			{statement: "GRANT SELECT ON *.* TO u", want: DDL},
			{statement: "REVOKE SELECT ON *.* FROM u", want: DDL},
			{statement: "IMPORT TABLE FROM '/tmp/t.sdi'", want: DDL},
			{statement: "REPLACE INTO t VALUES(1, 2, 3)", want: DML},
			{statement: "LOAD DATA INFILE '/tmp/t.csv' INTO TABLE t", want: DML},
			{statement: "CHANGE REPLICATION SOURCE TO SOURCE_HOST = '127.0.0.1'", want: DML},
			{statement: "START REPLICA", want: DML},
			{statement: "STOP REPLICA", want: DML},
			{statement: "RESET REPLICA", want: DML},
			{statement: "PURGE BINARY LOGS TO 'mysql-bin.010'", want: DML},
			{statement: "RESET MASTER", want: DML},
			{statement: "START GROUP_REPLICATION", want: DML},
			{statement: "STOP GROUP_REPLICATION", want: DML},
			{statement: "HELP 'contents'", want: QueryTypeUnknown},
		}
		for _, tc := range tests {
			span, err := resolveExtractorSpan(catalog, tc.statement, "db", false)
			require.NoError(t, err, tc.statement)
			require.Equal(t, tc.want, span.Type, tc.statement)
			require.Empty(t, span.Results, tc.statement)
		}
	})
}

func TestQuerySpanScenarioBatch1_ExpressionCoverage(t *testing.T) {
	catalog := newExtractorTestCatalog()

	t.Run("quoted_identifier_names", func(t *testing.T) {
		span, err := resolveExtractorSpan(
			catalog,
			"SELECT `select`, `Camel` AS `ExactCase` FROM keyword_table",
			"db",
			false,
		)
		require.NoError(t, err)
		require.Equal(t, []QuerySpanResult{
			{Name: "select", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "keyword_table", Column: "select"}}), IsPlainField: true},
			{Name: "ExactCase", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "keyword_table", Column: "Camel"}}), IsPlainField: true},
		}, span.Results)
	})

	t.Run("variable_refs_have_empty_lineage", func(t *testing.T) {
		span, err := resolveExtractorSpan(
			catalog,
			"SELECT @x AS user_var, @@version AS system_var",
			"db",
			false,
		)
		require.NoError(t, err)
		require.Equal(t, []QuerySpanResult{
			{Name: "user_var", SourceColumns: ColumnSet{}, IsPlainField: true},
			{Name: "system_var", SourceColumns: ColumnSet{}, IsPlainField: true},
		}, span.Results)
	})

	t.Run("json_extract_operator_uses_object_lineage", func(t *testing.T) {
		span, err := resolveExtractorSpan(
			catalog,
			"SELECT product_info->'$.id' AS id FROM products",
			"db",
			false,
		)
		require.NoError(t, err)
		require.Equal(t, []QuerySpanResult{
			{Name: "id", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "products", Column: "product_info"}}), IsPlainField: false},
		}, span.Results)
	})

	t.Run("group_concat_separator_is_a_string", func(t *testing.T) {
		// opt_separator is SEPARATOR STRING_VALUE (sql_parser_mysql_mode.y:3975),
		// so a separator reads no column, and OceanBase 4.3.5 answers 1064 for
		// SEPARATOR b.
		span, err := resolveExtractorSpan(
			catalog,
			"SELECT GROUP_CONCAT(a SEPARATOR ',') AS gc FROM t",
			"db",
			false,
		)
		require.NoError(t, err)
		require.Equal(t, []QuerySpanResult{
			{Name: "gc", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}}), IsPlainField: false},
		}, span.Results)
		_, err = resolveExtractorSpan(catalog, "SELECT GROUP_CONCAT(a SEPARATOR b) AS gc FROM t", "db", false)
		require.Error(t, err)
	})

	t.Run("group_concat_order_by_contributes_lineage", func(t *testing.T) {
		span, err := resolveExtractorSpan(
			catalog,
			"SELECT GROUP_CONCAT(a ORDER BY b) AS gc FROM t",
			"db",
			false,
		)
		require.NoError(t, err)
		require.Equal(t, []QuerySpanResult{
			{Name: "gc", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}, {Database: "db", Table: "t", Column: "b"}}), IsPlainField: false},
		}, span.Results)
	})

	t.Run("duplicate_output_names_are_preserved", func(t *testing.T) {
		span, err := resolveExtractorSpan(
			catalog,
			"SELECT a AS x, b AS x FROM t",
			"db",
			false,
		)
		require.NoError(t, err)
		require.Equal(t, []QuerySpanResult{
			{Name: "x", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}}), IsPlainField: true},
			{Name: "x", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "b"}}), IsPlainField: true},
		}, span.Results)
	})
}

func TestQuerySpanScenarioBatch2_FromJoinAndScopeCoverage(t *testing.T) {
	catalog := newExtractorTestCatalog()

	t.Run("table_source_variants", func(t *testing.T) {
		tests := []struct {
			name      string
			statement string
			want      []QuerySpanResult
			wantSrc   []ColumnResource
		}{
			{
				name:      "database_qualified_table",
				statement: "SELECT db.t.a FROM db.t",
				want: []QuerySpanResult{
					{Name: "a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}}), IsPlainField: true},
				},
				wantSrc: []ColumnResource{{Database: "db", Table: "t"}},
			},
			{
				name:      "parenthesized_single_table",
				statement: "SELECT a FROM (t)",
				want: []QuerySpanResult{
					{Name: "a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}}), IsPlainField: true},
				},
				wantSrc: []ColumnResource{{Database: "db", Table: "t"}},
			},
			{
				name:      "comma_separated_tables",
				statement: "SELECT t1.a, t2.b FROM t1, t2",
				want: []QuerySpanResult{
					{Name: "a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t1", Column: "a"}}), IsPlainField: true},
					{Name: "b", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t2", Column: "b"}}), IsPlainField: true},
				},
				wantSrc: []ColumnResource{{Database: "db", Table: "t1"}, {Database: "db", Table: "t2"}},
			},
			{
				name:      "dual",
				statement: "SELECT 1 FROM DUAL",
				want: []QuerySpanResult{
					{Name: "1", SourceColumns: ColumnSet{}, IsPlainField: true},
				},
			},
			{
				name:      "dual_join_right",
				statement: "SELECT t.a FROM t JOIN DUAL",
				want: []QuerySpanResult{
					{Name: "a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}}), IsPlainField: true},
				},
				wantSrc: []ColumnResource{{Database: "db", Table: "t"}},
			},
			{
				name:      "dual_join_left",
				statement: "SELECT t.a FROM DUAL JOIN t",
				want: []QuerySpanResult{
					{Name: "a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}}), IsPlainField: true},
				},
				wantSrc: []ColumnResource{{Database: "db", Table: "t"}},
			},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				span, err := resolveExtractorSpan(catalog, tc.statement, "db", false)
				require.NoError(t, err)
				require.Equal(t, tc.want, span.Results)
				require.Equal(t, sourceColumnSetFromResources(tc.wantSrc), span.SourceColumns)
			})
		}
	})

	t.Run("join_variants_expose_columns", func(t *testing.T) {
		tests := []string{
			"SELECT * FROM t1 STRAIGHT_JOIN t2",
			"SELECT * FROM t1 LEFT JOIN t2 ON t1.a = t2.a",
			"SELECT * FROM t1 RIGHT JOIN t2 ON t1.a = t2.a",
			"SELECT * FROM t1 NATURAL LEFT JOIN t2",
			"SELECT * FROM t1 NATURAL RIGHT JOIN t2",
		}
		// A NATURAL join's merged column comes from both tables.
		want := []QuerySpanResult{
			{Name: "a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t1", Column: "a"}, {Database: "db", Table: "t2", Column: "a"}}), IsPlainField: true},
			{Name: "b", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t1", Column: "b"}, {Database: "db", Table: "t2", Column: "b"}}), IsPlainField: true},
			{Name: "c", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t1", Column: "c"}, {Database: "db", Table: "t2", Column: "c"}}), IsPlainField: true},
		}
		wantNonNatural := []QuerySpanResult{
			{Name: "a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t1", Column: "a"}}), IsPlainField: true},
			{Name: "b", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t1", Column: "b"}}), IsPlainField: true},
			{Name: "c", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t1", Column: "c"}}), IsPlainField: true},
		}
		wantNonNatural = append(wantNonNatural,
			QuerySpanResult{Name: "a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t2", Column: "a"}}), IsPlainField: true},
			QuerySpanResult{Name: "b", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t2", Column: "b"}}), IsPlainField: true},
			QuerySpanResult{Name: "c", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t2", Column: "c"}}), IsPlainField: true},
		)
		for i, statement := range tests {
			span, err := resolveExtractorSpan(catalog, statement, "db", false)
			require.NoError(t, err, statement)
			if i < 3 {
				require.Equal(t, wantNonNatural, span.Results, statement)
			} else {
				require.Equal(t, want, span.Results, statement)
			}
			require.Equal(t, sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t1"}, {Database: "db", Table: "t2"}}), span.SourceColumns, statement)
		}
	})

	t.Run("derived_table_edges", func(t *testing.T) {
		_, err := resolveExtractorSpan(
			catalog,
			"SELECT * FROM (SELECT a, b FROM t) AS x(c1)",
			"db",
			false,
		)
		require.ErrorContains(t, err, "derived table column list length")

		span, err := resolveExtractorSpan(
			catalog,
			"SELECT y.c1 FROM (SELECT x.c1 FROM (SELECT a FROM t) AS x(c1)) AS y",
			"db",
			false,
		)
		require.NoError(t, err)
		require.Equal(t, []QuerySpanResult{
			{Name: "c1", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}}), IsPlainField: true},
		}, span.Results)

		span, err = resolveExtractorSpan(
			catalog,
			"SELECT x.a FROM t, LATERAL (SELECT t.a) AS x",
			"db",
			false,
		)
		require.NoError(t, err)
		require.Equal(t, []QuerySpanResult{
			{Name: "a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}}), IsPlainField: true},
		}, span.Results)
	})

	t.Run("join_dependency_edges", func(t *testing.T) {
		tests := []struct {
			name      string
			statement string
			want      []QuerySpanResult
			wantSrc   []ColumnResource
		}{
			{
				name:      "join_on_subquery_contributes_access_table",
				statement: "SELECT t1.a FROM t1 JOIN t2 ON EXISTS (SELECT t.a FROM t WHERE t.a = t2.a)",
				want: []QuerySpanResult{
					{Name: "a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t1", Column: "a"}}), IsPlainField: true},
				},
				wantSrc: []ColumnResource{
					{Database: "db", Table: "t1"},
					{Database: "db", Table: "t2"},
					{Database: "db", Table: "t"},
				},
			},
			{
				name:      "nested_join_tree_preserves_visible_order",
				statement: "SELECT * FROM t1 JOIN (t2 JOIN t ON t2.a = t.a) ON t1.a = t2.a",
				want: []QuerySpanResult{
					{Name: "a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t1", Column: "a"}}), IsPlainField: true},
					{Name: "b", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t1", Column: "b"}}), IsPlainField: true},
					{Name: "c", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t1", Column: "c"}}), IsPlainField: true},
					{Name: "a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t2", Column: "a"}}), IsPlainField: true},
					{Name: "b", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t2", Column: "b"}}), IsPlainField: true},
					{Name: "c", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t2", Column: "c"}}), IsPlainField: true},
					{Name: "a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}}), IsPlainField: true},
					{Name: "b", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "b"}}), IsPlainField: true},
					{Name: "c", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "c"}}), IsPlainField: true},
				},
				wantSrc: []ColumnResource{
					{Database: "db", Table: "t1"},
					{Database: "db", Table: "t2"},
					{Database: "db", Table: "t"},
				},
			},
			{
				name:      "parenthesized_table_reference_list",
				statement: "SELECT t1.a, t2.b FROM (t1, t2)",
				want: []QuerySpanResult{
					{Name: "a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t1", Column: "a"}}), IsPlainField: true},
					{Name: "b", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t2", Column: "b"}}), IsPlainField: true},
				},
				wantSrc: []ColumnResource{
					{Database: "db", Table: "t1"},
					{Database: "db", Table: "t2"},
				},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				span, err := resolveExtractorSpan(catalog, tc.statement, "db", false)
				require.NoError(t, err)
				require.Equal(t, tc.want, span.Results)
				require.Equal(t, sourceColumnSetFromResources(tc.wantSrc), span.SourceColumns)
			})
		}
	})
}

func TestQuerySpanScenarioBatch2_SubqueryCTEAndSetCoverage(t *testing.T) {
	catalog := newExtractorTestCatalog()

	t.Run("subquery_access_locations", func(t *testing.T) {
		tests := []string{
			"SELECT a FROM t WHERE EXISTS (SELECT t2.a FROM t2)",
			"SELECT a FROM t GROUP BY a HAVING EXISTS (SELECT t2.a FROM t2)",
			"SELECT a FROM t ORDER BY (SELECT t2.a FROM t2)",
		}
		for _, statement := range tests {
			span, err := resolveExtractorSpan(catalog, statement, "db", false)
			require.NoError(t, err, statement)
			require.Equal(t, sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t"}, {Database: "db", Table: "t2"}}), span.SourceColumns, statement)
		}
	})

	t.Run("cte_and_set_error_edges", func(t *testing.T) {
		_, err := resolveExtractorSpan(
			catalog,
			"WITH c(x) AS (SELECT a, b FROM t) SELECT * FROM c",
			"db",
			false,
		)
		require.ErrorContains(t, err, "CTE column list")

		_, err = resolveExtractorSpan(
			catalog,
			"WITH RECURSIVE c(x) AS (SELECT a, b FROM t UNION SELECT a, b FROM t2) SELECT * FROM c",
			"db",
			false,
		)
		require.ErrorContains(t, err, "different column counts")

		_, err = resolveExtractorSpan(
			catalog,
			"SELECT a FROM t UNION SELECT a, b FROM t2",
			"db",
			false,
		)
		require.ErrorContains(t, err, "UNION operator left has 1 fields, right has 2 fields")

		span, err := resolveExtractorSpan(
			catalog,
			"WITH c1 AS (SELECT a FROM c2), c2 AS (SELECT a FROM t) SELECT * FROM c1",
			"db",
			false,
		)
		require.NoError(t, err)
		require.Empty(t, span.Results)
		require.ErrorAs(t, span.NotFoundError, new(*NotFoundError))
	})

	t.Run("set_operation_in_derived_and_cte", func(t *testing.T) {
		tests := []string{
			"SELECT x.a FROM (SELECT a FROM t UNION SELECT a FROM t2) AS x",
			"WITH c AS (SELECT a FROM t UNION SELECT a FROM t2) SELECT a FROM c",
		}
		for _, statement := range tests {
			span, err := resolveExtractorSpan(catalog, statement, "db", false)
			require.NoError(t, err, statement)
			require.Equal(t, []QuerySpanResult{
				{Name: "a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}, {Database: "db", Table: "t2", Column: "a"}}), IsPlainField: false},
			}, span.Results, statement)
		}
	})

	t.Run("additional_set_operation_edges", func(t *testing.T) {
		tests := []struct {
			name      string
			statement string
			want      []QuerySpanResult
		}{
			{
				name:      "intersect",
				statement: "SELECT a FROM t INTERSECT SELECT a FROM t2",
				want: []QuerySpanResult{
					{Name: "a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}, {Database: "db", Table: "t2", Column: "a"}}), IsPlainField: false},
				},
			},
			{
				name:      "except",
				statement: "SELECT a FROM t EXCEPT SELECT a FROM t2",
				want: []QuerySpanResult{
					{Name: "a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}, {Database: "db", Table: "t2", Column: "a"}}), IsPlainField: false},
				},
			},
			{
				name:      "parenthesized_grouping",
				statement: "SELECT a FROM t UNION (SELECT a FROM t2 UNION SELECT b FROM t1)",
				want: []QuerySpanResult{
					{Name: "a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}, {Database: "db", Table: "t2", Column: "a"}, {Database: "db", Table: "t1", Column: "b"}}), IsPlainField: false},
				},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				span, err := resolveExtractorSpan(catalog, tc.statement, "db", false)
				require.NoError(t, err)
				require.Equal(t, tc.want, span.Results)
			})
		}
	})

	t.Run("scope_resolution_edges", func(t *testing.T) {
		tests := []struct {
			name      string
			statement string
			want      []QuerySpanResult
		}{
			{
				name:      "nested_subquery_nearest_scope_alias",
				statement: "SELECT (SELECT x.a FROM t2 AS x WHERE x.b = outer_t.b) AS inner_a FROM t AS outer_t",
				want: []QuerySpanResult{
					{Name: "inner_a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t2", Column: "a"}}), IsPlainField: true},
				},
			},
			{
				name:      "unqualified_inner_column_resolves_in_the_subquery_first",
				statement: "SELECT (SELECT a FROM t2) AS inner_a FROM t1",
				want: []QuerySpanResult{
					{Name: "inner_a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t2", Column: "a"}}), IsPlainField: true},
				},
			},
			{
				name:      "inner_alias_shadows_the_outer_alias",
				statement: "SELECT (SELECT x.a FROM t2 AS x) AS inner_a FROM t1 AS x",
				want: []QuerySpanResult{
					{Name: "inner_a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t2", Column: "a"}}), IsPlainField: true},
				},
			},
			{
				name:      "cte_name_shadows_physical_table",
				statement: "WITH t AS (SELECT b FROM t2) SELECT * FROM t",
				want: []QuerySpanResult{
					{Name: "b", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t2", Column: "b"}}), IsPlainField: true},
				},
			},
			{
				name:      "nested_cte_visibility_order",
				statement: "WITH c1 AS (SELECT a FROM t), c2 AS (SELECT a FROM c1 UNION SELECT b FROM t2) SELECT * FROM c2",
				want: []QuerySpanResult{
					{Name: "a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}, {Database: "db", Table: "t2", Column: "b"}}), IsPlainField: false},
				},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				span, err := resolveExtractorSpan(catalog, tc.statement, "db", false)
				require.NoError(t, err)
				require.Equal(t, tc.want, span.Results)
			})
		}
	})

	t.Run("recursive_cte_reaches_stable_source_closure", func(t *testing.T) {
		span, err := resolveExtractorSpan(
			catalog,
			"WITH RECURSIVE c(x) AS (SELECT a FROM t UNION SELECT t2.b FROM t2 JOIN c ON c.x = t2.a) SELECT * FROM c",
			"db",
			false,
		)
		require.NoError(t, err)
		require.Equal(t, []QuerySpanResult{
			{Name: "x", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}, {Database: "db", Table: "t2", Column: "b"}}), IsPlainField: false},
		}, span.Results)
	})

	t.Run("recursive_cte_left_deep_set_op_registers_before_recursive_left_branch", func(t *testing.T) {
		span, err := resolveExtractorSpan(
			catalog,
			"WITH RECURSIVE c(x) AS (SELECT a FROM t UNION SELECT t2.b FROM t2 JOIN c ON c.x = t2.a UNION SELECT t2.c FROM t2 JOIN c ON c.x = t2.b) SELECT * FROM c",
			"db",
			false,
		)
		require.NoError(t, err)
		require.Equal(t, []QuerySpanResult{
			{Name: "x", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}, {Database: "db", Table: "t2", Column: "b"}, {Database: "db", Table: "t2", Column: "c"}}), IsPlainField: false},
		}, span.Results)
	})
}

func TestQuerySpanScenarioBatch4_JSONTableCoverage(t *testing.T) {
	span, err := resolveExtractorSpan(newExtractorTestCatalog(),
		`SELECT *
FROM products,
JSON_TABLE(product_info, '$' COLUMNS (
  product_id INT PATH '$.id',
  NESTED PATH '$.variants[*]' COLUMNS (
    variant_id INT PATH '$.id',
    variant_name VARCHAR(50) PATH '$.name'
  )
)) AS jt`,
		"db",
		false,
	)
	require.NoError(t, err)
	require.Equal(t, []QuerySpanResult{
		{Name: "product_info", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "products", Column: "product_info"}}), IsPlainField: true},
		{Name: "product_id", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "products", Column: "product_info"}}), IsPlainField: true},
		{Name: "variant_id", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "products", Column: "product_info"}}), IsPlainField: true},
		{Name: "variant_name", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "products", Column: "product_info"}}), IsPlainField: true},
	}, span.Results)

	// A JSON_TABLE without an alias is one of the grammar's table factor
	// alternatives; its columns carry the lineage of the expression the factor
	// reads, so they resolve named as written.
	span, err = resolveExtractorSpan(newExtractorTestCatalog(),
		`SELECT *
FROM products,
JSON_TABLE(product_info, '$' COLUMNS (
  product_id INT PATH '$.id'
))`,
		"db",
		false,
	)
	require.NoError(t, err)
	require.Equal(t, []QuerySpanResult{
		{Name: "product_info", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "products", Column: "product_info"}}), IsPlainField: true},
		{Name: "product_id", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "products", Column: "product_info"}}), IsPlainField: true},
	}, span.Results)
}

func TestQuerySpanScenarioBatch4_AccessTableExpressions(t *testing.T) {
	catalog := newExtractorTestCatalog()

	tests := []struct {
		name      string
		statement string
		want      []ColumnResource
	}{
		{
			name:      "function_argument_subquery",
			statement: "SELECT COALESCE((SELECT t2.a FROM t2), a) FROM t",
			want: []ColumnResource{
				{Database: "db", Table: "t"},
				{Database: "db", Table: "t2"},
			},
		},
		{
			name:      "function_order_by_subquery",
			statement: "SELECT GROUP_CONCAT(a ORDER BY (SELECT t2.b FROM t2)) FROM t",
			want: []ColumnResource{
				{Database: "db", Table: "t"},
				{Database: "db", Table: "t2"},
			},
		},
		{
			name:      "json_table_document_subquery",
			statement: "SELECT * FROM JSON_TABLE((SELECT product_info FROM products), '$' COLUMNS (product_id INT PATH '$.id')) AS jt",
			want: []ColumnResource{
				{Database: "db", Table: "products"},
			},
		},
		{
			name:      "values_subquery",
			statement: "VALUES ROW((SELECT a FROM t))",
			want: []ColumnResource{
				{Database: "db", Table: "t"},
			},
		},
		{
			name:      "set_assignment_subquery",
			statement: "SET @x = (SELECT a FROM t)",
			want: []ColumnResource{
				{Database: "db", Table: "t"},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			span, err := resolveExtractorSpan(catalog, tc.statement, "db", false)
			require.NoError(t, err)
			require.Equal(t, sourceColumnSetFromResources(tc.want), span.SourceColumns)
		})
	}
}

func TestQuerySpanScenarioBatch3_ErrorAndMetadataCoverage(t *testing.T) {
	t.Run("error_semantics", func(t *testing.T) {
		catalog := newExtractorTestCatalog()

		_, err := resolveExtractorSpan(catalog, "SELECT a INTO @x FROM t", "db", false)
		require.ErrorContains(t, err, "unsupported select statement with into")

		_, err = resolveExtractorSpan(catalog, "SELECT FROM", "db", false)
		require.Error(t, err)

		_, err = resolveExtractorSpan(catalog, "SELECT * FORM invalid_table;", "db", false)
		require.Error(t, err)

		span, err := resolveExtractorSpan(catalog, "", "db", false)
		require.NoError(t, err)
		require.Equal(t, Select, span.Type)
		require.Empty(t, span.Results)
		require.Empty(t, span.SourceColumns)

		span, err = resolveExtractorSpan(catalog, "EXPLAIN SELECT * FROM missing", "db", false)
		require.NoError(t, err)
		require.Equal(t, Explain, span.Type)
		require.Empty(t, span.Results)
		require.Nil(t, span.NotFoundError)

		_, err = resolveExtractorSpan(catalog, "INSERT INTO t SELECT * FROM mysql.user", "db", false)
		require.ErrorIs(t, err, ErrMixUserSystemTables)

		_, err = resolveExtractorSpan(catalog, "CREATE TABLE new_t AS SELECT * FROM mysql.user", "db", false)
		require.ErrorIs(t, err, ErrMixUserSystemTables)

		_, err = resolveExtractorSpan(catalog, "EXPLAIN SELECT * FROM t, mysql.user", "db", false)
		require.ErrorIs(t, err, ErrMixUserSystemTables)

		_, err = resolveExtractorSpan(catalog, "CALL p((SELECT a FROM t), (SELECT user FROM mysql.user))", "db", false)
		require.ErrorIs(t, err, ErrMixUserSystemTables)

		_, err = resolveExtractorSpan(catalog, "SET @x = (SELECT a FROM t), @y = (SELECT user FROM mysql.user)", "db", false)
		require.ErrorIs(t, err, ErrMixUserSystemTables)

		_, err = resolveExtractorSpan(catalog, "LOCK TABLES t READ, mysql.user READ", "db", false)
		require.ErrorIs(t, err, ErrMixUserSystemTables)

		_, err = resolveExtractorSpan(catalog, "DROP TABLE t, mysql.user", "db", false)
		require.ErrorIs(t, err, ErrMixUserSystemTables)

		_, err = resolveExtractorSpan(catalog, "RENAME TABLE t TO t_new, mysql.user TO mysql.user_new", "db", false)
		require.ErrorIs(t, err, ErrMixUserSystemTables)
	})

	t.Run("view_lineage", func(t *testing.T) {
		catalog := newViewExtractorCatalog()

		span, err := resolveExtractorSpan(catalog, "SELECT * FROM v", "db", false)
		require.NoError(t, err)
		require.Equal(t, []QuerySpanResult{
			{Name: "va", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}}), IsPlainField: true},
			{Name: "b", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "b"}}), IsPlainField: true},
		}, span.Results)

		span, err = resolveExtractorSpan(catalog, "SELECT * FROM v2", "db", false)
		require.NoError(t, err)
		require.Equal(t, []QuerySpanResult{
			{Name: "va", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}}), IsPlainField: true},
		}, span.Results)
	})

	t.Run("missing_view_dependency_fails_open", func(t *testing.T) {
		catalog := newViewExtractorCatalog()
		span, err := resolveExtractorSpan(catalog, "SELECT * FROM broken_v", "db", false)
		require.NoError(t, err)
		require.Equal(t, Select, span.Type)
		require.Empty(t, span.Results)
		require.ErrorAs(t, span.NotFoundError, new(*NotFoundError))
	})

	t.Run("system_schema_detection", func(t *testing.T) {
		tests := []struct {
			name                string
			statement           string
			ignoreCaseSensitive bool
			wantType            QueryType
			wantSource          []ColumnResource
			wantNotFound        bool
		}{
			{
				name:      "information_schema_lowercase",
				statement: "SELECT * FROM information_schema.tables",
				wantType:  SelectInfoSchema,
			},
			{
				name:                "information_schema_uppercase_ignore_case",
				statement:           "SELECT * FROM INFORMATION_SCHEMA.tables",
				ignoreCaseSensitive: true,
				wantType:            SelectInfoSchema,
			},
			{
				name:      "information_schema_uppercase_case_sensitive",
				statement: "SELECT * FROM INFORMATION_SCHEMA.tables",
				wantType:  SelectInfoSchema,
			},
			{
				name:      "performance_schema_lowercase",
				statement: "SELECT * FROM performance_schema.events_statements_current",
				wantType:  SelectInfoSchema,
			},
			{
				name:      "mysql_lowercase",
				statement: "SELECT * FROM mysql.user",
				wantType:  SelectInfoSchema,
			},
			{
				name:                "mysql_uppercase_ignore_case",
				statement:           "SELECT * FROM MYSQL.user",
				ignoreCaseSensitive: true,
				wantType:            SelectInfoSchema,
			},
			{
				name:         "mysql_uppercase_case_sensitive",
				statement:    "SELECT * FROM MYSQL.user",
				wantType:     Select,
				wantNotFound: true,
				wantSource: []ColumnResource{
					{Database: "MYSQL", Table: "user"},
				},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				span, err := resolveExtractorSpan(newExtractorTestCatalog(), tc.statement, "db", tc.ignoreCaseSensitive)
				require.NoError(t, err)
				require.Equal(t, tc.wantType, span.Type)
				require.Equal(t, sourceColumnSetFromResources(tc.wantSource), span.SourceColumns)
				if tc.wantNotFound {
					require.ErrorAs(t, span.NotFoundError, new(*NotFoundError))
				} else {
					require.Nil(t, span.NotFoundError)
				}
			})
		}
	})
}

func TestQuerySpanScenarioBatch4_CaseSensitivityCoverage(t *testing.T) {
	catalog := newExtractorTestCatalog()

	t.Run("table_alias_case_sensitivity", func(t *testing.T) {
		span, err := resolveExtractorSpan(catalog, "SELECT X.a FROM t AS x", "db", false)
		require.NoError(t, err)
		require.Empty(t, span.Results)
		require.ErrorAs(t, span.NotFoundError, new(*NotFoundError))

		span, err = resolveExtractorSpan(catalog, "SELECT X.a FROM t AS x", "db", true)
		require.NoError(t, err)
		require.Equal(t, []QuerySpanResult{
			{Name: "a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}}), IsPlainField: true},
		}, span.Results)
	})

	t.Run("derived_table_alias_case_sensitivity", func(t *testing.T) {
		span, err := resolveExtractorSpan(catalog, "SELECT X.a FROM (SELECT a FROM t) AS x", "db", false)
		require.NoError(t, err)
		require.Empty(t, span.Results)
		require.ErrorAs(t, span.NotFoundError, new(*NotFoundError))

		span, err = resolveExtractorSpan(catalog, "SELECT X.a FROM (SELECT a FROM t) AS x", "db", true)
		require.NoError(t, err)
		require.Equal(t, []QuerySpanResult{
			{Name: "a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}}), IsPlainField: true},
		}, span.Results)
	})

	t.Run("cte_name_case_sensitivity", func(t *testing.T) {
		span, err := resolveExtractorSpan(catalog, "WITH c AS (SELECT a FROM t) SELECT C.a FROM c AS C", "db", false)
		require.NoError(t, err)
		require.Equal(t, []QuerySpanResult{
			{Name: "a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}}), IsPlainField: true},
		}, span.Results)
	})
}

func TestQuerySpanScenarioBatch5_MetadataCoverage(t *testing.T) {
	t.Run("duplicate_database_case_uses_ignore_case_setting", func(t *testing.T) {
		catalog := newCaseCollisionExtractorCatalog()

		span, err := resolveExtractorSpan(catalog, "SELECT DB.t.a FROM DB.t", "db", false)
		require.NoError(t, err)
		require.Equal(t, []QuerySpanResult{
			{Name: "a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "DB", Table: "t", Column: "a"}}), IsPlainField: true},
		}, span.Results)
		require.Equal(t, sourceColumnSetFromResources([]ColumnResource{{Database: "DB", Table: "t"}}), span.SourceColumns)

		span, err = resolveExtractorSpan(catalog, "SELECT DB.t.a FROM DB.t", "db", true)
		require.NoError(t, err)
		require.Equal(t, []QuerySpanResult{
			{Name: "a", SourceColumns: sourceColumnSetFromResources([]ColumnResource{{Database: "db", Table: "t", Column: "a"}}), IsPlainField: true},
		}, span.Results)
		require.Equal(t, sourceColumnSetFromResources([]ColumnResource{{Database: "DB", Table: "t"}}), span.SourceColumns)
	})
}

func sourceColumnSetFromResources(resources []ColumnResource) ColumnSet {
	result := make(ColumnSet)
	for _, resource := range resources {
		result[resource] = true
	}
	return result
}

func newViewExtractorCatalog() Catalog {
	dbMeta := &metadata.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadata.SchemaMetadata{
			{
				Name: "",
				Tables: []*metadata.TableMetadata{
					{
						Name: "t",
						Columns: []*metadata.ColumnMetadata{
							{Name: "a"},
							{Name: "b"},
						},
					},
				},
				Views: []*metadata.ViewMetadata{
					{
						Name:       "v",
						Definition: "SELECT a AS va, b FROM t",
					},
					{
						Name:       "v2",
						Definition: "SELECT va FROM v",
					},
					{
						Name:       "broken_v",
						Definition: "SELECT missing FROM t",
					},
				},
			},
		},
	}
	return newExtractorCatalog([]*metadata.DatabaseSchemaMetadata{dbMeta})
}

func newExtractorTestCatalog() Catalog {
	dbMeta := &metadata.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadata.SchemaMetadata{
			{
				Name: "",
				Tables: []*metadata.TableMetadata{
					{
						Name: "t",
						Columns: []*metadata.ColumnMetadata{
							{Name: "a"},
							{Name: "b"},
							{Name: "c"},
						},
					},
					{
						Name: "t1",
						Columns: []*metadata.ColumnMetadata{
							{Name: "a"},
							{Name: "b"},
							{Name: "c"},
						},
					},
					{
						Name: "t2",
						Columns: []*metadata.ColumnMetadata{
							{Name: "a"},
							{Name: "b"},
							{Name: "c"},
						},
					},
					{
						Name: "galleries",
						Columns: []*metadata.ColumnMetadata{
							{Name: "id"},
							{Name: "city"},
						},
					},
					{
						Name: "paintings",
						Columns: []*metadata.ColumnMetadata{
							{Name: "id"},
							{Name: "gallery_id"},
						},
					},
					{
						Name: "products",
						Columns: []*metadata.ColumnMetadata{
							{Name: "product_info"},
						},
					},
					{
						Name: "keyword_table",
						Columns: []*metadata.ColumnMetadata{
							{Name: "select"},
							{Name: "Camel"},
						},
					},
				},
			},
		},
	}
	return newExtractorCatalog([]*metadata.DatabaseSchemaMetadata{dbMeta})
}

func newCaseCollisionExtractorCatalog() Catalog {
	dbMetadata := &metadata.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadata.SchemaMetadata{
			{
				Name: "",
				Tables: []*metadata.TableMetadata{
					{
						Name: "t",
						Columns: []*metadata.ColumnMetadata{
							{Name: "a"},
						},
					},
				},
			},
		},
	}
	upperMetadata := &metadata.DatabaseSchemaMetadata{
		Name: "DB",
		Schemas: []*metadata.SchemaMetadata{
			{
				Name: "",
				Tables: []*metadata.TableMetadata{
					{
						Name: "t",
						Columns: []*metadata.ColumnMetadata{
							{Name: "a"},
						},
					},
				},
			},
		},
	}
	return newExtractorCatalog([]*metadata.DatabaseSchemaMetadata{dbMetadata, upperMetadata})
}

func resolveExtractorSpan(catalog Catalog, statement, defaultDatabase string, ignoreCaseSensitive bool) (*ResolvedSpan, error) {
	return ResolveQuerySpan(statement, ResolveOptions{
		DefaultDatabase:     defaultDatabase,
		IgnoreCaseSensitive: ignoreCaseSensitive,
		Catalog:             catalog,
	})
}

func newExtractorCatalog(dbs []*metadata.DatabaseSchemaMetadata) Catalog {
	byName := make(map[string]*metadata.DatabaseSchemaMetadata, len(dbs))
	names := make([]string, 0, len(dbs))
	for _, db := range dbs {
		byName[db.GetName()] = db
		names = append(names, db.GetName())
	}
	return Catalog{
		ListDatabaseNames: func() ([]string, error) {
			return names, nil
		},
		GetDatabase: func(name string) (Database, error) {
			if db, ok := byName[name]; ok {
				return &corpusDatabase{proto: db}, nil
			}
			return nil, fmt.Errorf("database %q not found", name)
		},
	}
}
