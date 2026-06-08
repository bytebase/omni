package parser

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

func TestRedshiftObjectStatementsParse(t *testing.T) {
	tests := []struct {
		sql        string
		command    string
		objectType string
	}{
		{sql: "CREATE DATASHARE s;", command: "create", objectType: "datashare"},
		{sql: "ALTER DATASHARE s ADD TABLE public.t;", command: "alter", objectType: "datashare"},
		{sql: "DROP DATASHARE IF EXISTS s;", command: "drop", objectType: "datashare"},
		{sql: "CREATE EXTERNAL SCHEMA s FROM DATA CATALOG DATABASE 'db' REGION 'us-east-1' IAM_ROLE DEFAULT;", command: "create", objectType: "external schema"},
		{sql: "CREATE EXTERNAL TABLE s.t (id INT) STORED AS PARQUET LOCATION 's3://bucket/';", command: "create", objectType: "external table"},
		{sql: "CREATE EXTERNAL VIEW v AS SELECT 1;", command: "create", objectType: "external view"},
		{sql: "CREATE EXTERNAL PROTECTED VIEW v IF NOT EXISTS AS SELECT 1;", command: "create", objectType: "external protected view"},
		{sql: "CREATE MASKING POLICY p WITH (a INT) USING (a);", command: "create", objectType: "masking policy"},
		{sql: "ATTACH MASKING POLICY p ON t(c) TO ROLE r;", command: "attach", objectType: "masking policy"},
		{sql: "CREATE RLS POLICY p WITH (a INT) USING (a);", command: "create", objectType: "rls policy"},
		{sql: "CREATE IDENTITY PROVIDER p TYPE oauth2 NAMESPACE ns;", command: "create", objectType: "identity provider"},
		{sql: "CREATE MODEL m FROM (SELECT 1 AS y) TARGET y FUNCTION f IAM_ROLE DEFAULT;", command: "create", objectType: "model"},
		{sql: "CREATE LIBRARY f LANGUAGE plpythonu FROM 's3://bucket/lib.zip' IAM_ROLE DEFAULT;", command: "create", objectType: "library"},
	}

	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			stmt := firstRedshiftObjectStmt(t, tt.sql)
			if stmt.Command != tt.command {
				t.Fatalf("expected command %q, got %q", tt.command, stmt.Command)
			}
			if stmt.ObjectType != tt.objectType {
				t.Fatalf("expected object type %q, got %q", tt.objectType, stmt.ObjectType)
			}
			if stmt.Name == nil && tt.command != "attach" {
				t.Fatalf("expected object name")
			}
		})
	}
}

func TestRedshiftCreateExternalFunctionParse(t *testing.T) {
	tests := []struct {
		sql       string
		wantName  string
		wantToken string
	}{
		{
			sql: `CREATE EXTERNAL FUNCTION public.f_distance(float, float, float, float)
RETURNS float
STABLE
LAMBDA 'lambda_distance'
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftLambdaRole';`,
			wantName:  "public.f_distance",
			wantToken: "lambda_distance",
		},
		{
			sql: `CREATE EXTERNAL FUNCTION get_current_time()
RETURNS timestamp
IMMUTABLE
LAMBDA 'lambda_get_time'
IAM_ROLE DEFAULT;`,
			wantName:  "get_current_time",
			wantToken: "default",
		},
		{
			sql: `CREATE OR REPLACE EXTERNAL FUNCTION calculate_score(integer, integer, integer)
RETURNS decimal(10,2)
STABLE
LAMBDA 'lambda_score_calculator'
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftLambdaRole'
RETRY_TIMEOUT 45000
MAX_BATCH_ROWS 10000
MAX_BATCH_SIZE 256 KB;`,
			wantName:  "calculate_score",
			wantToken: "max_batch_size",
		},
	}

	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			stmt := firstRedshiftObjectStmt(t, tt.sql)
			if stmt.ObjectType != "external function" {
				t.Fatalf("ObjectType = %q, want external function", stmt.ObjectType)
			}
			if got := redshiftObjectNameString(stmt.Name); got != tt.wantName {
				t.Fatalf("Name = %q, want %q", got, tt.wantName)
			}
			if !redshiftStringListContains(stmt.Options, tt.wantToken) {
				t.Fatalf("Options missing %q: %#v", tt.wantToken, stmt.Options)
			}
		})
	}
}

func TestRedshiftCreateExternalModelParse(t *testing.T) {
	tests := []struct {
		sql       string
		wantName  string
		wantToken string
	}{
		{
			sql: `CREATE EXTERNAL MODEL customer_review_model
FUNCTION customer_review_llm
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftBedrockRole'
MODEL_TYPE BEDROCK
SETTINGS (
    MODEL_ID 'anthropic.claude-3-sonnet-20240229-v1:0',
    REQUEST_TYPE UNIFIED,
    RESPONSE_TYPE SUPER
);`,
			wantName:  "customer_review_model",
			wantToken: "bedrock",
		},
		{
			sql: `CREATE EXTERNAL MODEL sentiment_model
FUNCTION sentiment_analysis_llm
IAM_ROLE DEFAULT
MODEL_TYPE BEDROCK
SETTINGS (
    MODEL_ID 'anthropic.claude-v2:1'
);`,
			wantName:  "sentiment_model",
			wantToken: "default",
		},
		{
			sql: `CREATE EXTERNAL MODEL sales_forecast_model
FROM 's3://mybucket/models/sales_forecast.tar.gz'
FUNCTION_NAME predict_sales
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftMLRole';`,
			wantName:  "sales_forecast_model",
			wantToken: "function_name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			stmt := firstRedshiftObjectStmt(t, tt.sql)
			if stmt.ObjectType != "external model" {
				t.Fatalf("ObjectType = %q, want external model", stmt.ObjectType)
			}
			if got := redshiftObjectNameString(stmt.Name); got != tt.wantName {
				t.Fatalf("Name = %q, want %q", got, tt.wantName)
			}
			if !redshiftStringListContains(stmt.Options, tt.wantToken) {
				t.Fatalf("Options missing %q: %#v", tt.wantToken, stmt.Options)
			}
		})
	}
}

func TestRedshiftCancelParse(t *testing.T) {
	tests := []struct {
		sql       string
		wantToken string
	}{
		{
			sql:       "CANCEL 18764;",
			wantToken: "18764",
		},
		{
			sql:       "CANCEL 18764 'Long-running query terminated by administrator';",
			wantToken: "Long-running query terminated by administrator",
		},
	}

	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			stmt := firstRedshiftObjectStmt(t, tt.sql)
			if stmt.Command != "cancel" {
				t.Fatalf("Command = %q, want cancel", stmt.Command)
			}
			if stmt.ObjectType != "query" {
				t.Fatalf("ObjectType = %q, want query", stmt.ObjectType)
			}
			if !redshiftStringListContains(stmt.Options, tt.wantToken) {
				t.Fatalf("Options missing %q: %#v", tt.wantToken, stmt.Options)
			}
		})
	}
}

func firstRedshiftObjectStmt(t *testing.T, sql string) *nodes.RedshiftObjectStmt {
	t.Helper()
	tree, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if len(tree.Items) != 1 {
		t.Fatalf("expected one statement, got %d", len(tree.Items))
	}
	raw, ok := tree.Items[0].(*nodes.RawStmt)
	if !ok {
		t.Fatalf("expected RawStmt, got %T", tree.Items[0])
	}
	stmt, ok := raw.Stmt.(*nodes.RedshiftObjectStmt)
	if !ok {
		t.Fatalf("expected RedshiftObjectStmt, got %T", raw.Stmt)
	}
	return stmt
}

func redshiftObjectNameString(list *nodes.List) string {
	if list == nil {
		return ""
	}
	parts := make([]string, 0, len(list.Items))
	for _, item := range list.Items {
		if s, ok := item.(*nodes.String); ok {
			parts = append(parts, s.Str)
		}
	}
	return strings.Join(parts, ".")
}

func redshiftStringListContains(list *nodes.List, value string) bool {
	if list == nil {
		return false
	}
	for _, item := range list.Items {
		if s, ok := item.(*nodes.String); ok && s.Str == value {
			return true
		}
	}
	return false
}
