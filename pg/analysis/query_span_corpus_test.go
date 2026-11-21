package analysis

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v3"
)

// The corpora under testdata/query-span are Bytebase's PostgreSQL query-span
// and query-type fixtures, recorded from the resolver masking ran on before
// this package existed. They are the masking-parity bar for this package and
// must be reproduced, never re-recorded against this implementation.

type spanCorpusCase struct {
	Description     string          `yaml:"description"`
	Statement       string          `yaml:"statement"`
	DefaultDatabase string          `yaml:"defaultDatabase"`
	Metadata        string          `yaml:"metadata"`
	QuerySpan       *spanCorpusSpan `yaml:"querySpan"`
}

type spanCorpusSpan struct {
	Type             int                `yaml:"type"`
	Results          []spanCorpusResult `yaml:"results"`
	SourceColumns    []spanCorpusColumn `yaml:"sourcecolumns"`
	PredicateColumns []spanCorpusColumn `yaml:"predicatecolumns"`
}

type spanCorpusResult struct {
	Name             string                `yaml:"name"`
	SourceColumns    []spanCorpusColumn    `yaml:"sourcecolumns"`
	IsPlainField     bool                  `yaml:"isplainfield"`
	SourceFieldPaths []spanCorpusFieldPath `yaml:"sourcefieldpaths"`
	SelectAsterisk   bool                  `yaml:"selectasterisk"`
}

type spanCorpusFieldPath struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

type spanCorpusColumn struct {
	Server   string `yaml:"server"`
	Database string `yaml:"database"`
	Schema   string `yaml:"schema"`
	Table    string `yaml:"table"`
	Column   string `yaml:"column"`
}

func renderSpanColumns(set ColumnSet) []spanCorpusColumn {
	out := make([]spanCorpusColumn, 0, len(set))
	for resource := range set {
		out = append(out, spanCorpusColumn{
			Server:   resource.Server,
			Database: resource.Database,
			Schema:   resource.Schema,
			Table:    resource.Table,
			Column:   resource.Column,
		})
	}
	slices.SortFunc(out, func(a, b spanCorpusColumn) int {
		return strings.Compare(spanColumnKey(a), spanColumnKey(b))
	})
	return out
}

func spanColumnKey(c spanCorpusColumn) string {
	var parts []string
	for _, p := range []string{c.Server, c.Database, c.Schema, c.Table, c.Column} {
		if p != "" {
			parts = append(parts, p)
		}
	}
	return strings.Join(parts, ".")
}

func renderSpan(span *ResolvedSpan) *spanCorpusSpan {
	out := &spanCorpusSpan{
		Type:             int(span.Type),
		Results:          make([]spanCorpusResult, 0, len(span.Results)),
		SourceColumns:    renderSpanColumns(span.SourceColumns),
		PredicateColumns: renderSpanColumns(span.PredicateColumns),
	}
	for _, result := range span.Results {
		out.Results = append(out.Results, spanCorpusResult{
			Name:             result.Name,
			SourceColumns:    renderSpanColumns(result.SourceColumns),
			IsPlainField:     result.IsPlainField,
			SourceFieldPaths: []spanCorpusFieldPath{},
		})
	}
	return out
}

func TestQuerySpanCorpus(t *testing.T) {
	paths, err := filepath.Glob(filepath.Join("testdata", "query-span", "*.yaml"))
	require.NoError(t, err)
	require.NotEmpty(t, paths)

	unmarshaler := protojson.UnmarshalOptions{DiscardUnknown: true}
	for _, path := range paths {
		content, err := os.ReadFile(path)
		require.NoError(t, err)
		var cases []spanCorpusCase
		require.NoError(t, yaml.Unmarshal(content, &cases), path)
		for i, tc := range cases {
			name := tc.Description
			if name == "" {
				name = tc.Statement
			}
			t.Run(filepath.Base(path)+"/"+name, func(t *testing.T) {
				require.NotNil(t, tc.QuerySpan, "case %d has no querySpan", i)
				dbMetadata := &metadatapb.DatabaseSchemaMetadata{}
				require.NoError(t, unmarshaler.Unmarshal([]byte(tc.Metadata), dbMetadata))
				span, err := ResolveQuerySpan(context.Background(), tc.Statement, ResolveOptions{
					DefaultDatabase: tc.DefaultDatabase,
					Catalog:         testCatalog(dbMetadata),
				})
				require.NoError(t, err, "idx: %d statement: %s", i, tc.Statement)
				require.Equal(t, tc.QuerySpan, renderSpan(span), "idx: %d statement: %s", i, tc.Statement)
			})
		}
	}
}

func mustHaveExactSource(t *testing.T, result QuerySpanResult, table, col string) {
	t.Helper()
	want := ColumnResource{
		Database: testDB,
		Schema:   testSchema,
		Table:    table,
		Column:   col,
	}
	if !result.SourceColumns[want] {
		t.Errorf("result %q missing source %+v; have %+v", result.Name, want, result.SourceColumns)
	}
}

func TestQuerySpanSurvivesABadQuotedIdentifier(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: testDB,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: testSchema,
			Tables: []*metadatapb.TableMetadata{
				{
					Name:    "'weird'table",
					Columns: []*metadatapb.ColumnMetadata{column("id", "int4")},
				},
				{
					Name:    "accounts",
					Columns: []*metadatapb.ColumnMetadata{column("id", "int4"), column("email", "text")},
				},
			},
		}},
	}

	span := spanFor(t, `SELECT id, email FROM accounts`, meta)
	require.Len(t, span.Results, 2)
	for _, r := range span.Results {
		require.NotEmpty(t, r.SourceColumns, "result %q has empty sources", r.Name)
	}
}

func TestQuerySpanWithSelectedSchemaFallsBackToPublic(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: testDB,
		Schemas: []*metadatapb.SchemaMetadata{
			{Name: "app"},
			{
				Name: testSchema,
				Tables: []*metadatapb.TableMetadata{{
					Name:    "customer",
					Columns: []*metadatapb.ColumnMetadata{column("ssn", "text")},
				}},
			},
		},
	}
	span, err := ResolveQuerySpan(context.Background(), `SELECT ssn FROM customer`, ResolveOptions{
		DefaultDatabase: testDB,
		SearchPath:      SearchPathForSelectedSchema("app", []string{"public"}),
		Catalog:         testCatalog(meta),
	})
	require.NoError(t, err)
	require.Len(t, span.Results, 1)
	mustHaveExactSource(t, span.Results[0], "customer", "ssn")
	require.True(t, span.SourceColumns[ColumnResource{Database: testDB, Schema: testSchema, Table: "customer"}])
}

func TestQuerySpanSelectedSchemaTakesPrecedenceOverPublic(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: testDB,
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "app",
				Tables: []*metadatapb.TableMetadata{{
					Name:    "customer",
					Columns: []*metadatapb.ColumnMetadata{column("email", "text")},
				}},
			},
			{
				Name: testSchema,
				Tables: []*metadatapb.TableMetadata{{
					Name:    "customer",
					Columns: []*metadatapb.ColumnMetadata{column("ssn", "text")},
				}},
			},
		},
	}
	span, err := ResolveQuerySpan(context.Background(), `SELECT email FROM customer`, ResolveOptions{
		DefaultDatabase: testDB,
		SearchPath:      SearchPathForSelectedSchema("app", []string{"public"}),
		Catalog:         testCatalog(meta),
	})
	require.NoError(t, err)
	require.Len(t, span.Results, 1)
	want := ColumnResource{Database: testDB, Schema: "app", Table: "customer", Column: "email"}
	require.True(t, span.Results[0].SourceColumns[want])
}

func TestQuerySpanSurvivesABrokenEnum(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: testDB,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: testSchema,
			Tables: []*metadatapb.TableMetadata{{
				Name: "tasks",
				Columns: []*metadatapb.ColumnMetadata{
					column("id", "int4"),
					column("status", "public.nonexistent_enum"),
					column("title", "text"),
				},
			}},
		}},
	}

	span := spanFor(t, `SELECT id, status, title FROM tasks`, meta)
	require.Len(t, span.Results, 3)
	for i, want := range []string{"id", "status", "title"} {
		require.Equal(t, want, span.Results[i].Name)
		require.NotEmpty(t, span.Results[i].SourceColumns)
	}
}

func TestQuerySpanKeepsAHealthyNeighborOfABrokenTable(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: testDB,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: testSchema,
			Tables: []*metadatapb.TableMetadata{
				{
					Name:    "broken",
					Columns: []*metadatapb.ColumnMetadata{column("id", "int4"), column("bad", "public.nonexistent_type")},
				},
				{
					Name:    "healthy",
					Columns: []*metadatapb.ColumnMetadata{column("id", "int4"), column("label", "text")},
				},
			},
		}},
	}

	span := spanFor(t, `SELECT id, label FROM healthy`, meta)
	require.Len(t, span.Results, 2)
	for _, r := range span.Results {
		require.NotEmpty(t, r.SourceColumns)
	}
}

func TestQuerySpanResolvesAViewOverABrokenTable(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: testDB,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: testSchema,
			Tables: []*metadatapb.TableMetadata{{
				Name: "records",
				Columns: []*metadatapb.ColumnMetadata{
					column("id", "int4"),
					column("status", "public.nonexistent_enum"),
					column("title", "text"),
				},
			}},
			Views: []*metadatapb.ViewMetadata{{
				Name:       "records_view",
				Definition: "SELECT id, title, status FROM records",
				DependencyColumns: []*metadatapb.DependencyColumn{
					{Schema: "public", Table: "records", Column: "id"},
					{Schema: "public", Table: "records", Column: "title"},
					{Schema: "public", Table: "records", Column: "status"},
				},
			}},
		}},
	}

	span := spanFor(t, `SELECT id, title, status FROM records_view`, meta)
	require.Len(t, span.Results, 3)
	for _, r := range span.Results {
		require.NotEmpty(t, r.SourceColumns)
	}
}

func TestQuerySpanResolvesADeclaredEnum(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: testDB,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: testSchema,
			EnumTypes: []*metadatapb.EnumTypeMetadata{{
				Name:   "task_status",
				Values: []string{"pending", "running", "done"},
			}},
			Tables: []*metadatapb.TableMetadata{{
				Name: "tasks",
				Columns: []*metadatapb.ColumnMetadata{
					column("id", "int4"),
					column("status", "public.task_status"),
				},
			}},
		}},
	}

	span := spanFor(t, `SELECT id, status FROM tasks WHERE status = 'pending'`, meta)
	require.Len(t, span.Results, 2)
	mustHaveExactSource(t, span.Results[0], "tasks", "id")
	mustHaveExactSource(t, span.Results[1], "tasks", "status")
}

func TestQuerySpanOnACorrelatedRangeFunctionSubquery(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: testDB,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: testSchema,
			Tables: []*metadatapb.TableMetadata{
				{
					Name: "compliance_case_record_audits",
					Columns: []*metadatapb.ColumnMetadata{
						column("case_id", "text"),
						column("entity_id", "text"),
						column("reason", "text"),
						column("remark", "text"),
						column("case_record", "jsonb"),
						column("reviewer_by", "text"),
						column("deleted_at", "timestamptz"),
					},
				},
				{
					Name: "compliance_cases",
					Columns: []*metadatapb.ColumnMetadata{
						column("case_id", "text"),
						column("case_info", "jsonb"),
						column("deleted_at", "timestamptz"),
					},
				},
			},
		}},
	}
	sql := `
select
  a.case_id,
  a.entity_id as uid,
  a.reason,
  a.remark,
  (
    select elem->>'matchedDateTimeValue'
    from jsonb_array_elements(a.case_record::jsonb->'secondaryFieldResults') elem
    where elem->>'typeId' = '******'
    limit 1
  ) as wc_dob,
  (
    select elem->>'matchedValue'
    from jsonb_array_elements(a.case_record::jsonb->'secondaryFieldResults') elem
    where elem->>'typeId' = '******'
    limit 1
  ) as wc_citizenship,
  c.case_info->>'birth_date' as kyc_dob,
  c.case_info->>'nationality' as kyc_citizenship,
  a.reviewer_by as review_by
from compliance_case_record_audits a
inner join compliance_cases c
  on c.case_id = a.case_id and c.deleted_at is null
where a.reviewer_by in ('****** ', '******')
  and a.deleted_at is null;
`
	span := spanFor(t, sql, meta)
	require.Len(t, span.Results, 9)
	mustHaveExactSource(t, span.Results[4], "compliance_case_record_audits", "case_record")
	mustHaveExactSource(t, span.Results[5], "compliance_case_record_audits", "case_record")
	mustHaveExactSource(t, span.Results[6], "compliance_cases", "case_info")
	mustHaveExactSource(t, span.Results[7], "compliance_cases", "case_info")
}

func TestQuerySpanKeepsJSONBLineageThroughCTEsAndALateralJoin(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: testDB,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: testSchema,
			Tables: []*metadatapb.TableMetadata{{
				Name:    "ai_conversation",
				Columns: []*metadatapb.ColumnMetadata{column("id", "text"), column("parts", "jsonb")},
			}},
		}},
	}

	tests := []struct {
		name string
		sql  string
	}{
		{
			name: "single_cte",
			sql: `
WITH convs AS (
  SELECT ac.id, ac.parts AS conv_parts
  FROM ai_conversation ac
  LIMIT 1
)
SELECT
  LEFT(
    COALESCE(
      (SELECT string_agg(v::text, ' ')
       FROM jsonb_array_elements_text(
         jsonb_path_query_array(t.part->'bodyData', 'strict $.**.text', '{}'::jsonb, true)
       ) AS v),
      '<no text>'
    ),
    2000
  ) AS content
FROM convs c
CROSS JOIN LATERAL jsonb_array_elements(c.conv_parts) WITH ORDINALITY AS t(part, part_ord)
WHERE t.part->>'type' IN ('prompt', 'text');
`,
		},
		{
			name: "multiple_ctes",
			sql: `
WITH cte1 AS (
  SELECT ac.id, ac.parts AS conv_parts
  FROM ai_conversation ac
  LIMIT 1
),
cte2 AS (
  SELECT id, conv_parts
  FROM cte1
)
SELECT
  LEFT(
    COALESCE(
      (SELECT string_agg(v::text, ' ')
       FROM jsonb_array_elements_text(
         jsonb_path_query_array(t.part->'bodyData', 'strict $.**.text', '{}'::jsonb, true)
       ) AS v),
      '<no text>'
    ),
    2000
  ) AS content
FROM cte2 c
CROSS JOIN LATERAL jsonb_array_elements(c.conv_parts) WITH ORDINALITY AS t(part, part_ord)
WHERE t.part->>'type' IN ('prompt', 'text');
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			span := spanFor(t, tt.sql, meta)
			require.Len(t, span.Results, 1)
			mustHaveExactSource(t, span.Results[0], "ai_conversation", "parts")
		})
	}
}
