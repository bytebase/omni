package analysis

import (
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
)

func oracleLongTailOptions(t *testing.T) ResolveOptions {
	t.Helper()

	const metadataText = `{
		"name": "PUBLIC",
		"schemas": [{
			"name": "",
			"tables": [{
				"name": "T",
				"columns": [
					{"name": "A"},
					{"name": "B"},
					{"name": "C"},
					{"name": "J"},
					{"name": "X"}
				]
			}, {
				"name": "T2",
				"columns": [
					{"name": "C"},
					{"name": "D"}
				]
			}, {
				"name": "TRADES",
				"columns": [
					{"name": "ACCOUNT_ID"},
					{"name": "TRADE_TIME"},
					{"name": "PRICE"}
				]
			}]
		}]
	}`

	meta := &metadatapb.DatabaseSchemaMetadata{}
	require.NoError(t, protojson.Unmarshal([]byte(metadataText), meta))
	return ResolveOptions{
		Instance:        corpusInstanceIDA,
		DefaultDatabase: "PUBLIC",
		Catalog:         corpusCatalog([]*metadatapb.DatabaseSchemaMetadata{meta}),
	}
}

func columnSetFromList(cols []ColumnResource) ColumnSet {
	result := make(ColumnSet, len(cols))
	for _, c := range cols {
		result[c] = true
	}
	return result
}

func TestOracleLongTailTableSources(t *testing.T) {
	tests := []struct {
		name              string
		statement         string
		wantResult        string
		wantSourceColumns []ColumnResource
		wantQuerySources  []ColumnResource
	}{
		{
			name:              "json table",
			statement:         "SELECT JT.ID FROM T, JSON_TABLE(T.J, '$' COLUMNS (ID NUMBER PATH '$.id')) JT",
			wantResult:        "ID",
			wantSourceColumns: []ColumnResource{{Database: "PUBLIC", Table: "T", Column: "J"}},
			wantQuerySources:  []ColumnResource{{Database: "PUBLIC", Table: "T"}},
		},
		{
			name:              "xml table",
			statement:         "SELECT XT.ID FROM T, XMLTABLE('/root' PASSING T.X COLUMNS ID NUMBER PATH 'id') XT",
			wantResult:        "ID",
			wantSourceColumns: []ColumnResource{{Database: "PUBLIC", Table: "T", Column: "X"}},
			wantQuerySources:  []ColumnResource{{Database: "PUBLIC", Table: "T"}},
		},
		{
			name:              "containers",
			statement:         "SELECT A FROM CONTAINERS(T)",
			wantResult:        "A",
			wantSourceColumns: []ColumnResource{{Database: "PUBLIC", Table: "T", Column: "A"}},
			wantQuerySources:  []ColumnResource{{Database: "PUBLIC", Table: "T"}},
		},
		{
			name:              "inline external table",
			statement:         "SELECT EXT.A FROM EXTERNAL ((A NUMBER, B VARCHAR2(20)) TYPE ORACLE_LOADER DEFAULT DIRECTORY D ACCESS PARAMETERS (FIELDS TERMINATED BY ',') LOCATION ('x.csv')) EXT",
			wantResult:        "A",
			wantSourceColumns: []ColumnResource{},
			wantQuerySources:  []ColumnResource{},
		},
	}

	opts := oracleLongTailOptions(t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			span, err := ResolveQuerySpan(test.statement, opts)
			require.NoError(t, err)
			require.NotNil(t, span)
			require.Len(t, span.Results, 1)
			require.Equal(t, test.wantResult, span.Results[0].Name)
			require.Equal(t, columnSetFromList(test.wantSourceColumns), span.Results[0].SourceColumns)
			require.Equal(t, columnSetFromList(test.wantQuerySources), span.SourceColumns)
		})
	}
}

func TestOracleAnalyticCountStarSources(t *testing.T) {
	span, err := ResolveQuerySpan("SELECT COUNT(*) OVER (PARTITION BY B ORDER BY A) AS C FROM T", oracleLongTailOptions(t))
	require.NoError(t, err)
	require.NotNil(t, span)
	require.Len(t, span.Results, 1)
	require.Equal(t, "C", span.Results[0].Name)
	require.Equal(t, columnSetFromList([]ColumnResource{
		{Database: "PUBLIC", Table: "T", Column: "A"},
		{Database: "PUBLIC", Table: "T", Column: "B"},
	}), span.Results[0].SourceColumns)
}

func TestOracleUnqualifiedStarExpandsAllCommaSources(t *testing.T) {
	tests := []struct {
		name      string
		statement string
		want      []QuerySpanResult
	}{
		{
			name:      "base tables",
			statement: "SELECT * FROM T, T2",
			want: []QuerySpanResult{
				{Name: "A", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "A"}})},
				{Name: "B", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "B"}})},
				{Name: "C", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "C"}})},
				{Name: "J", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "J"}})},
				{Name: "X", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "X"}})},
				{Name: "C", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T2", Column: "C"}})},
				{Name: "D", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T2", Column: "D"}})},
			},
		},
		{
			name:      "json table",
			statement: "SELECT * FROM T, JSON_TABLE(T.J, '$' COLUMNS (ID NUMBER PATH '$.id')) JT",
			want: []QuerySpanResult{
				{Name: "A", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "A"}})},
				{Name: "B", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "B"}})},
				{Name: "C", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "C"}})},
				{Name: "J", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "J"}})},
				{Name: "X", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "X"}})},
				{Name: "ID", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "J"}})},
			},
		},
		{
			name:      "join",
			statement: "SELECT * FROM T JOIN T2 ON T.A = T2.C",
			want: []QuerySpanResult{
				{Name: "A", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "A"}})},
				{Name: "B", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "B"}})},
				{Name: "C", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "C"}})},
				{Name: "J", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "J"}})},
				{Name: "X", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "X"}})},
				{Name: "C", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T2", Column: "C"}})},
				{Name: "D", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T2", Column: "D"}})},
			},
		},
		{
			name:      "natural join",
			statement: "SELECT * FROM T NATURAL JOIN T2",
			want: []QuerySpanResult{
				{Name: "A", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "A"}})},
				{Name: "B", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "B"}})},
				{Name: "C", SourceColumns: columnSetFromList([]ColumnResource{
					{Database: "PUBLIC", Table: "T", Column: "C"},
					{Database: "PUBLIC", Table: "T2", Column: "C"},
				})},
				{Name: "J", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "J"}})},
				{Name: "X", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "X"}})},
				{Name: "D", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T2", Column: "D"}})},
			},
		},
	}

	opts := oracleLongTailOptions(t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			span, err := ResolveQuerySpan(test.statement, opts)
			require.NoError(t, err)
			require.NotNil(t, span)
			require.Len(t, span.Results, len(test.want))
			for i, want := range test.want {
				require.Equal(t, want.Name, span.Results[i].Name)
				require.Equal(t, want.SourceColumns, span.Results[i].SourceColumns)
			}
		})
	}
}

func TestOracleCTEScopeAndSetOperations(t *testing.T) {
	tests := []struct {
		name      string
		statement string
		want      []QuerySpanResult
	}{
		{
			name:      "sibling cte references previous cte",
			statement: "WITH C1 AS (SELECT A FROM T), C2 AS (SELECT A FROM C1 UNION SELECT C FROM T2) SELECT * FROM C2",
			want: []QuerySpanResult{
				{Name: "A", SourceColumns: columnSetFromList([]ColumnResource{
					{Database: "PUBLIC", Table: "T", Column: "A"},
					{Database: "PUBLIC", Table: "T2", Column: "C"},
				})},
			},
		},
		{
			name:      "cte name shadows physical table",
			statement: "WITH T AS (SELECT C FROM T2) SELECT * FROM T",
			want: []QuerySpanResult{
				{Name: "C", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T2", Column: "C"}})},
			},
		},
		{
			name:      "nested cte name shadows outer cte",
			statement: "WITH C AS (WITH C AS (SELECT C FROM T2) SELECT C FROM C) SELECT * FROM C",
			want: []QuerySpanResult{
				{Name: "C", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T2", Column: "C"}})},
			},
		},
		{
			name:      "multiple ctes join",
			statement: "WITH C1 AS (SELECT A FROM T), C2 AS (SELECT C FROM T2) SELECT C1.A, C2.C FROM C1 JOIN C2 ON C1.A = C2.C",
			want: []QuerySpanResult{
				{Name: "A", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "A"}})},
				{Name: "C", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T2", Column: "C"}})},
			},
		},
		{
			name:      "intersect merges source columns positionally",
			statement: "SELECT A FROM T INTERSECT SELECT C FROM T2",
			want: []QuerySpanResult{
				{Name: "A", SourceColumns: columnSetFromList([]ColumnResource{
					{Database: "PUBLIC", Table: "T", Column: "A"},
					{Database: "PUBLIC", Table: "T2", Column: "C"},
				})},
			},
		},
		{
			name:      "minus merges source columns positionally",
			statement: "SELECT A FROM T MINUS SELECT C FROM T2",
			want: []QuerySpanResult{
				{Name: "A", SourceColumns: columnSetFromList([]ColumnResource{
					{Database: "PUBLIC", Table: "T", Column: "A"},
					{Database: "PUBLIC", Table: "T2", Column: "C"},
				})},
			},
		},
		{
			name:      "recursive cte reaches stable source closure",
			statement: "WITH C(X) AS (SELECT A FROM T UNION ALL SELECT T2.C FROM T2 JOIN C ON C.X = T2.C) SELECT * FROM C",
			want: []QuerySpanResult{
				{Name: "X", SourceColumns: columnSetFromList([]ColumnResource{
					{Database: "PUBLIC", Table: "T", Column: "A"},
					{Database: "PUBLIC", Table: "T2", Column: "C"},
				})},
			},
		},
	}

	opts := oracleLongTailOptions(t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			span, err := ResolveQuerySpan(test.statement, opts)
			require.NoError(t, err)
			require.NotNil(t, span)
			require.Equal(t, test.want, span.Results)
		})
	}
}

func TestOracleExpressionAndAccessTableCoverage(t *testing.T) {
	opts := oracleLongTailOptions(t)
	t.Run("expression nodes keep source columns", func(t *testing.T) {
		span, err := ResolveQuerySpan(
			`SELECT
  CASE WHEN A > 0 THEN B ELSE C END AS CASE_VALUE,
  DECODE(A, 1, B, C) AS DECODE_VALUE,
  CAST(A AS NUMBER) AS CAST_VALUE,
  B BETWEEN A AND C AS BETWEEN_VALUE,
  X LIKE J AS LIKE_VALUE
FROM T`,
			opts,
		)
		require.NoError(t, err)
		require.NotNil(t, span)
		want := []QuerySpanResult{
			{Name: "CASE_VALUE", SourceColumns: columnSetFromList([]ColumnResource{
				{Database: "PUBLIC", Table: "T", Column: "A"},
				{Database: "PUBLIC", Table: "T", Column: "B"},
				{Database: "PUBLIC", Table: "T", Column: "C"},
			})},
			{Name: "DECODE_VALUE", SourceColumns: columnSetFromList([]ColumnResource{
				{Database: "PUBLIC", Table: "T", Column: "A"},
				{Database: "PUBLIC", Table: "T", Column: "B"},
				{Database: "PUBLIC", Table: "T", Column: "C"},
			})},
			{Name: "CAST_VALUE", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "A"}})},
			{Name: "BETWEEN_VALUE", SourceColumns: columnSetFromList([]ColumnResource{
				{Database: "PUBLIC", Table: "T", Column: "A"},
				{Database: "PUBLIC", Table: "T", Column: "B"},
				{Database: "PUBLIC", Table: "T", Column: "C"},
			})},
			{Name: "LIKE_VALUE", SourceColumns: columnSetFromList([]ColumnResource{
				{Database: "PUBLIC", Table: "T", Column: "J"},
				{Database: "PUBLIC", Table: "T", Column: "X"},
			})},
		}
		require.Equal(t, want, span.Results)
	})

	t.Run("subqueries in non result clauses contribute access tables", func(t *testing.T) {
		tests := []string{
			"SELECT A FROM T WHERE EXISTS (SELECT C FROM T2 WHERE T2.C = T.A)",
			"SELECT A FROM T ORDER BY (SELECT C FROM T2)",
			"SELECT COALESCE((SELECT C FROM T2), A) FROM T",
		}
		for _, statement := range tests {
			span, err := ResolveQuerySpan(statement, opts)
			require.NoError(t, err, statement)
			require.NotNil(t, span, statement)
			require.Equal(t, columnSetFromList([]ColumnResource{
				{Database: "PUBLIC", Table: "T"},
				{Database: "PUBLIC", Table: "T2"},
			}), span.SourceColumns, statement)
		}
	})

	t.Run("cursor expression uses subquery scope", func(t *testing.T) {
		span, err := ResolveQuerySpan("SELECT CURSOR(SELECT C FROM T2) AS CUR FROM T", opts)
		require.NoError(t, err)
		require.NotNil(t, span)
		require.Equal(t, []QuerySpanResult{
			{Name: "CUR", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T2", Column: "C"}})},
		}, span.Results)
	})
}

func TestOracleResourceNotFoundPropagation(t *testing.T) {
	tests := []struct {
		name      string
		statement string
	}{
		{name: "top level table", statement: "SELECT A FROM MISSING_TABLE"},
		{name: "predicate subquery", statement: "SELECT A FROM T WHERE EXISTS (SELECT 1 FROM MISSING_TABLE)"},
		{name: "target subquery", statement: "SELECT (SELECT A FROM MISSING_TABLE) AS A FROM T"},
		{name: "recursive cte arm", statement: "WITH C(X) AS (SELECT A FROM T UNION ALL SELECT MISSING_TABLE.A FROM MISSING_TABLE JOIN C ON C.X = MISSING_TABLE.A) SELECT * FROM C"},
	}

	opts := oracleLongTailOptions(t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			span, err := ResolveQuerySpan(test.statement, opts)
			require.NoError(t, err)
			require.NotNil(t, span)
			require.Empty(t, span.Results)
			require.NotNil(t, span.NotFoundError)
		})
	}
}

func TestOracleTypedLongTailTableSources(t *testing.T) {
	tests := []struct {
		name      string
		statement string
		want      []QuerySpanResult
	}{
		{
			name:      "table collection column aliases",
			statement: "SELECT * FROM TABLE(pkg.values()) X(A, B)",
			want: []QuerySpanResult{
				{Name: "A", SourceColumns: columnSetFromList(nil)},
				{Name: "B", SourceColumns: columnSetFromList(nil)},
			},
		},
		{
			name:      "pivot typed output",
			statement: "SELECT * FROM (SELECT A, B FROM T) PIVOT (COUNT(*) AS CNT FOR B IN (1 AS ONE, 2 AS TWO))",
			want: []QuerySpanResult{
				{Name: "A", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "A"}})},
				{Name: "ONE_CNT", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "B"}})},
				{Name: "TWO_CNT", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "B"}})},
			},
		},
		{
			name:      "pivot trailing alias qualifies references",
			statement: "SELECT P.A, P.ONE_CNT FROM (SELECT A, B FROM T) PIVOT (COUNT(*) AS CNT FOR B IN (1 AS ONE, 2 AS TWO)) P ORDER BY P.A",
			want: []QuerySpanResult{
				{Name: "A", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "A"}})},
				{Name: "ONE_CNT", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "B"}})},
			},
		},
		{
			name:      "pivot over join source resolves merged columns",
			statement: "SELECT C FROM (T JOIN T2 USING (C)) PIVOT (COUNT(*) AS CNT FOR B IN (1 AS ONE))",
			want: []QuerySpanResult{
				{Name: "C", SourceColumns: columnSetFromList([]ColumnResource{
					{Database: "PUBLIC", Table: "T", Column: "C"},
					{Database: "PUBLIC", Table: "T2", Column: "C"},
				})},
			},
		},
		{
			name:      "pivot aggregate over join source resolves merged input",
			statement: "SELECT ONE_S FROM (T JOIN T2 USING (C)) PIVOT (SUM(C) AS S FOR B IN (1 AS ONE))",
			want: []QuerySpanResult{
				{Name: "ONE_S", SourceColumns: columnSetFromList([]ColumnResource{
					{Database: "PUBLIC", Table: "T", Column: "B"},
					{Database: "PUBLIC", Table: "T", Column: "C"},
					{Database: "PUBLIC", Table: "T2", Column: "C"},
				})},
			},
		},
		{
			name:      "pivot alias joins another table",
			statement: "SELECT P.ONE_CNT, T2.C FROM (SELECT A, B FROM T) PIVOT (COUNT(*) AS CNT FOR B IN (1 AS ONE)) P JOIN T2 ON P.A = T2.C",
			want: []QuerySpanResult{
				{Name: "ONE_CNT", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "B"}})},
				{Name: "C", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T2", Column: "C"}})},
			},
		},
		{
			name:      "pivot explicit projection",
			statement: "SELECT A FROM (SELECT A, B FROM T) PIVOT (COUNT(*) AS CNT FOR B IN (1 AS ONE, 2 AS TWO))",
			want: []QuerySpanResult{
				{Name: "A", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "A"}})},
			},
		},
		{
			name:      "pivot aggregate input is consumed",
			statement: "SELECT * FROM (SELECT A, B, C FROM T) PIVOT (SUM(C) AS S FOR B IN (1 AS ONE))",
			want: []QuerySpanResult{
				{Name: "A", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "A"}})},
				{Name: "ONE_S", SourceColumns: columnSetFromList([]ColumnResource{
					{Database: "PUBLIC", Table: "T", Column: "B"},
					{Database: "PUBLIC", Table: "T", Column: "C"},
				})},
			},
		},
		{
			name:      "pivot aggregate projected alias input is consumed",
			statement: "SELECT * FROM (SELECT A, B, C AS D FROM T) PIVOT (SUM(D) AS S FOR B IN (1 AS ONE))",
			want: []QuerySpanResult{
				{Name: "A", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "A"}})},
				{Name: "ONE_S", SourceColumns: columnSetFromList([]ColumnResource{
					{Database: "PUBLIC", Table: "T", Column: "B"},
					{Database: "PUBLIC", Table: "T", Column: "C"},
				})},
			},
		},
		{
			name:      "unpivot typed mappings",
			statement: "SELECT * FROM (SELECT A, B, C FROM T) UNPIVOT (VAL FOR COL IN (B AS 'B', C AS 'C'))",
			want: []QuerySpanResult{
				// Oracle puts the FOR column before the measure.
				{Name: "A", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "A"}})},
				{Name: "COL", SourceColumns: columnSetFromList(nil)},
				{Name: "VAL", SourceColumns: columnSetFromList([]ColumnResource{
					{Database: "PUBLIC", Table: "T", Column: "B"},
					{Database: "PUBLIC", Table: "T", Column: "C"},
				})},
			},
		},
		{
			name:      "unpivot explicit projection",
			statement: "SELECT VAL FROM (SELECT A, B, C FROM T) UNPIVOT (VAL FOR COL IN (B AS 'B', C AS 'C'))",
			want: []QuerySpanResult{
				{Name: "VAL", SourceColumns: columnSetFromList([]ColumnResource{
					{Database: "PUBLIC", Table: "T", Column: "B"},
					{Database: "PUBLIC", Table: "T", Column: "C"},
				})},
			},
		},
		{
			name:      "model dimension and measures",
			statement: "SELECT * FROM T MODEL DIMENSION BY (A) MEASURES (B) RULES (B[1] = 1)",
			want: []QuerySpanResult{
				{Name: "A", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "A"}})},
				{Name: "B", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "B"}})},
			},
		},
		{
			name:      "model explicit projection",
			statement: "SELECT B FROM T MODEL DIMENSION BY (A) MEASURES (B) RULES (B[1] = 1)",
			want: []QuerySpanResult{
				{Name: "B", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "B"}})},
			},
		},
		{
			name:      "model partition dimension and measures",
			statement: "SELECT * FROM T MODEL PARTITION BY (C) DIMENSION BY (A) MEASURES (B) RULES (B[1] = 1)",
			want: []QuerySpanResult{
				{Name: "C", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "C"}})},
				{Name: "A", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "A"}})},
				{Name: "B", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "B"}})},
			},
		},
		{
			name:      "model partition explicit projection",
			statement: "SELECT C FROM T MODEL PARTITION BY (C) DIMENSION BY (A) MEASURES (B) RULES (B[1] = 1)",
			want: []QuerySpanResult{
				{Name: "C", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "T", Column: "C"}})},
			},
		},
		{
			name: "match recognize measures",
			statement: `SELECT * FROM TRADES MATCH_RECOGNIZE (
  PARTITION BY ACCOUNT_ID
  ORDER BY TRADE_TIME
  MEASURES FIRST(PRICE) AS FIRST_PRICE, LAST(PRICE) AS LAST_PRICE
  ONE ROW PER MATCH
  PATTERN (A B+)
  DEFINE B AS B.PRICE > A.PRICE
) MR`,
			want: []QuerySpanResult{
				{Name: "ACCOUNT_ID", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "TRADES", Column: "ACCOUNT_ID"}})},
				{Name: "FIRST_PRICE", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "TRADES", Column: "PRICE"}})},
				{Name: "LAST_PRICE", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "TRADES", Column: "PRICE"}})},
			},
		},
		{
			name: "match recognize partition projection",
			statement: `SELECT ACCOUNT_ID FROM TRADES MATCH_RECOGNIZE (
  PARTITION BY ACCOUNT_ID
  ORDER BY TRADE_TIME
  MEASURES FIRST(PRICE) AS FIRST_PRICE
  ONE ROW PER MATCH
  PATTERN (A B+)
  DEFINE B AS B.PRICE > A.PRICE
) MR`,
			want: []QuerySpanResult{
				{Name: "ACCOUNT_ID", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "TRADES", Column: "ACCOUNT_ID"}})},
			},
		},
		{
			name: "match recognize pattern variable measures",
			statement: `SELECT * FROM TRADES MATCH_RECOGNIZE (
  PARTITION BY ACCOUNT_ID
  ORDER BY TRADE_TIME
  MEASURES FIRST(A.PRICE) AS FIRST_PRICE, LAST(B.PRICE) AS LAST_PRICE
  ONE ROW PER MATCH
  PATTERN (A B+)
  DEFINE B AS B.PRICE > A.PRICE
) MR`,
			want: []QuerySpanResult{
				{Name: "ACCOUNT_ID", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "TRADES", Column: "ACCOUNT_ID"}})},
				{Name: "FIRST_PRICE", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "TRADES", Column: "PRICE"}})},
				{Name: "LAST_PRICE", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "TRADES", Column: "PRICE"}})},
			},
		},
		{
			name: "match recognize all rows keeps input columns",
			statement: `SELECT * FROM TRADES MATCH_RECOGNIZE (
  PARTITION BY ACCOUNT_ID
  ORDER BY TRADE_TIME
  MEASURES FIRST(PRICE) AS FIRST_PRICE
  ALL ROWS PER MATCH
  PATTERN (A B+)
  DEFINE B AS B.PRICE > A.PRICE
) MR`,
			want: []QuerySpanResult{
				{Name: "ACCOUNT_ID", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "TRADES", Column: "ACCOUNT_ID"}}), IsPlainField: true},
				{Name: "TRADE_TIME", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "TRADES", Column: "TRADE_TIME"}}), IsPlainField: true},
				{Name: "PRICE", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "TRADES", Column: "PRICE"}}), IsPlainField: true},
				{Name: "FIRST_PRICE", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "TRADES", Column: "PRICE"}})},
			},
		},
		{
			name: "match recognize all rows pattern variable measure",
			statement: `SELECT * FROM TRADES MATCH_RECOGNIZE (
  PARTITION BY ACCOUNT_ID
  ORDER BY TRADE_TIME
  MEASURES FIRST(A.PRICE) AS FIRST_PRICE
  ALL ROWS PER MATCH
  PATTERN (A B+)
  DEFINE B AS B.PRICE > A.PRICE
) MR`,
			want: []QuerySpanResult{
				{Name: "ACCOUNT_ID", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "TRADES", Column: "ACCOUNT_ID"}}), IsPlainField: true},
				{Name: "TRADE_TIME", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "TRADES", Column: "TRADE_TIME"}}), IsPlainField: true},
				{Name: "PRICE", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "TRADES", Column: "PRICE"}}), IsPlainField: true},
				{Name: "FIRST_PRICE", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "TRADES", Column: "PRICE"}})},
			},
		},
		{
			name: "match recognize all rows input projection",
			statement: `SELECT PRICE FROM TRADES MATCH_RECOGNIZE (
  PARTITION BY ACCOUNT_ID
  ORDER BY TRADE_TIME
  MEASURES FIRST(PRICE) AS FIRST_PRICE
  ALL ROWS PER MATCH
  PATTERN (A B+)
  DEFINE B AS B.PRICE > A.PRICE
) MR`,
			want: []QuerySpanResult{
				{Name: "PRICE", SourceColumns: columnSetFromList([]ColumnResource{{Database: "PUBLIC", Table: "TRADES", Column: "PRICE"}})},
			},
		},
	}

	opts := oracleLongTailOptions(t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			span, err := ResolveQuerySpan(test.statement, opts)
			require.NoError(t, err)
			require.NotNil(t, span)
			require.Equal(t, test.want, span.Results)
		})
	}
}
