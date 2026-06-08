package redshift

import "testing"

func TestPhase1SplitScenarios(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want int
	}{
		{
			name: "top-level semicolon statements",
			sql:  "SELECT 1; SELECT 2;",
			want: 2,
		},
		{
			name: "single-quoted string semicolon",
			sql:  "SELECT 'a;b'; SELECT 2;",
			want: 2,
		},
		{
			name: "double-quoted identifier semicolon",
			sql:  `SELECT "a;b"; SELECT 2;`,
			want: 2,
		},
		{
			name: "dollar-quoted string semicolon",
			sql:  "SELECT $$a;b$$; SELECT 2;",
			want: 2,
		},
		{
			name: "block comment semicolon",
			sql:  "SELECT /* ; */ 1; SELECT 2;",
			want: 2,
		},
		{
			name: "redshift procedure body semicolons",
			sql: `CREATE PROCEDURE refresh_sales()
LANGUAGE plpgsql
AS $$
BEGIN
  DELETE FROM sales_stage;
  INSERT INTO sales_stage SELECT * FROM sales;
END;
$$;
SELECT 1;`,
			want: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := nonEmptySegments(Split(tc.sql))
			if len(got) != tc.want {
				t.Fatalf("Split returned %d non-empty segments, want %d: %#v", len(got), tc.want, got)
			}
		})
	}
}

func TestPhase1CoreStatementScenarios(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want StatementType
	}{
		{
			name: "cte select",
			sql:  "WITH c AS (SELECT 1) SELECT * FROM c;",
			want: StatementTypeSelect,
		},
		{
			name: "insert select",
			sql:  "INSERT INTO t SELECT 1;",
			want: StatementTypeDML,
		},
		{
			name: "update",
			sql:  "UPDATE t SET a = 1 WHERE id = 2;",
			want: StatementTypeDML,
		},
		{
			name: "delete",
			sql:  "DELETE FROM t WHERE id = 1;",
			want: StatementTypeDML,
		},
		{
			name: "merge update",
			sql:  "MERGE INTO target USING source ON target.id = source.id WHEN MATCHED THEN UPDATE SET v = source.v;",
			want: StatementTypeDML,
		},
		{
			name: "merge remove duplicates",
			sql:  "MERGE INTO target USING source ON target.id = source.id REMOVE DUPLICATES;",
			want: StatementTypeDML,
		},
		{
			name: "plain explain select",
			sql:  "EXPLAIN SELECT 1;",
			want: StatementTypeExplain,
		},
		{
			name: "explain analyze select",
			sql:  "EXPLAIN ANALYZE SELECT 1;",
			want: StatementTypeExplainAnalyze,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := GetStatementTypes(tc.sql)
			if err != nil {
				t.Fatalf("GetStatementTypes returned error: %v", err)
			}
			if len(got) != 1 || got[0] != tc.want {
				t.Fatalf("GetStatementTypes = %v, want [%s]", got, tc.want)
			}
		})
	}
}

func nonEmptySegments(segments []Segment) []Segment {
	var result []Segment
	for _, segment := range segments {
		if !segment.Empty() {
			result = append(result, segment)
		}
	}
	return result
}
