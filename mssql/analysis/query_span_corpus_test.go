package analysis

import (
	"fmt"
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

// The corpora under testdata/query-span are Bytebase's SQL Server query-span
// fixtures, recorded from the resolver masking ran on before this package
// existed. They are the masking-parity bar for this package and must be
// reproduced, never re-recorded against this implementation.

type spanCorpusCase struct {
	Description         string          `yaml:"description"`
	Statement           string          `yaml:"statement"`
	DefaultDatabase     string          `yaml:"defaultDatabase"`
	IgnoreCaseSensitive bool            `yaml:"ignoreCaseSensitive"`
	Metadata            []string        `yaml:"metadata"`
	QuerySpan           *spanCorpusSpan `yaml:"querySpan"`
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

func catalogFromMetadata(databases []*metadatapb.DatabaseSchemaMetadata) Catalog {
	return Catalog{
		ListDatabaseNames: func() ([]string, error) {
			names := make([]string, 0, len(databases))
			for _, db := range databases {
				names = append(names, db.GetName())
			}
			return names, nil
		},
		GetDatabase: func(name string) (Database, error) {
			for _, db := range databases {
				if db.GetName() != name {
					continue
				}
				d := &testDatabase{schemas: make(map[string]*metadatapb.SchemaMetadata, len(db.GetSchemas()))}
				for _, schema := range db.GetSchemas() {
					d.schemas[schema.GetName()] = schema
					d.order = append(d.order, schema.GetName())
				}
				return d, nil
			}
			return nil, fmt.Errorf("database %q not found", name)
		},
	}
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
				var metas []*metadatapb.DatabaseSchemaMetadata
				for _, raw := range tc.Metadata {
					m := &metadatapb.DatabaseSchemaMetadata{}
					require.NoError(t, unmarshaler.Unmarshal([]byte(raw), m))
					metas = append(metas, m)
				}
				span, err := ResolveQuerySpan(tc.Statement, ResolveOptions{
					DefaultDatabase:     tc.DefaultDatabase,
					DefaultSchema:       "dbo",
					IgnoreCaseSensitive: tc.IgnoreCaseSensitive,
					Catalog:             catalogFromMetadata(metas),
					TempTables:          make(map[string]*TempTable),
				})
				require.NoError(t, err, "statement: %s", tc.Statement)
				require.Equal(t, tc.QuerySpan, renderSpan(span), "statement: %s", tc.Statement)
			})
		}
	}
}

func TestResolveQuerySpanUnsyncedDatabase(t *testing.T) {
	synced := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "dbo",
			Tables: []*metadatapb.TableMetadata{{Name: "t", Columns: columns("a")}},
		}},
	}
	syncedCat := catalogFromMetadata([]*metadatapb.DatabaseSchemaMetadata{synced})
	cat := Catalog{
		ListDatabaseNames: func() ([]string, error) {
			return []string{"db", "appdb"}, nil
		},
		GetDatabase: func(name string) (Database, error) {
			if name == "appdb" {
				return nil, nil
			}
			return syncedCat.GetDatabase(name)
		},
	}
	const notSynced = `database metadata for database "appdb" not found (database not synced)`

	tests := []struct {
		name            string
		statement       string
		defaultDatabase string
		wantErr         string
		wantNotFound    bool
	}{
		{name: "table in the unsynced connected database", statement: "SELECT * FROM t", defaultDatabase: "appdb", wantErr: notSynced},
		{name: "cross-database reference to an unsynced database", statement: "SELECT * FROM appdb.dbo.t", defaultDatabase: "db", wantErr: notSynced},
		{name: "statement that reads no table", statement: "SELECT 1", defaultDatabase: "appdb"},
		{name: "table missing from a synced database", statement: "SELECT * FROM missing", defaultDatabase: "db", wantNotFound: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			span, err := ResolveQuerySpan(tc.statement, ResolveOptions{
				DefaultDatabase:     tc.defaultDatabase,
				DefaultSchema:       "dbo",
				IgnoreCaseSensitive: true,
				Catalog:             cat,
				TempTables:          make(map[string]*TempTable),
			})
			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.wantNotFound, span.NotFoundError != nil)
		})
	}
}

func TestResolveQuerySpanCyclicTwoViews(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "dbo",
			Views: []*metadatapb.ViewMetadata{
				{Name: "v1", Definition: "CREATE VIEW [dbo].[v1] AS SELECT * FROM v2"},
				{Name: "v2", Definition: "CREATE VIEW [dbo].[v2] AS SELECT * FROM v1"},
			},
		}},
	}
	_, err := ResolveQuerySpan("SELECT * FROM v1", ResolveOptions{
		DefaultDatabase:     "db",
		DefaultSchema:       "dbo",
		IgnoreCaseSensitive: true,
		Catalog:             catalogFromMetadata([]*metadatapb.DatabaseSchemaMetadata{meta}),
		TempTables:          make(map[string]*TempTable),
	})
	require.ErrorContains(t, err, "cyclic view reference")
}

func TestResolveQuerySpanMarksPlanOutput(t *testing.T) {
	cat := catalogFromMetadata([]*metadatapb.DatabaseSchemaMetadata{{Name: "db"}})
	planOutput := map[string]bool{}
	cases := []struct {
		statement      string
		queryType      QueryType
		explainAnalyze bool
	}{
		{"SELECT 1", Select, false},
		{"SET /* x */ STATISTICS XML ON -- y", Select, true},
		{"SELECT 2", Select, true},
		{"SET NOCOUNT ON", Select, true},
		{"DELETE FROM t", DML, true},
		{"SET STATISTICS PROFILE OFF", Select, true},
		{"SELECT 8", Select, true},
		{"SET STATISTICS XML OFF", Select, true},
		{"SELECT 3", Select, false},
		{"SET STATISTICS IO ON", Select, false},
		{"SELECT 4", Select, false},
		{"set statistics io, profile on", Select, true},
		{"SELECT 5", Select, true},
		{"SET STATISTICS PROFILE OFF", Select, true},
		{"SET SHOWPLAN_TEXT ON -- z", Explain, false},
		{"SET NOCOUNT OFF", Explain, false},
		{"SET SHOWPLAN_ALL OFF", Explain, false},
		{"SELECT 6", Explain, false},
		{"SET SHOWPLAN_TEXT OFF", Explain, false},
		{"SELECT 7", Select, false},
	}
	for _, c := range cases {
		span, err := ResolveQuerySpan(c.statement, ResolveOptions{
			DefaultDatabase:     "db",
			DefaultSchema:       "dbo",
			IgnoreCaseSensitive: true,
			Catalog:             cat,
			PlanOutput:          planOutput,
		})
		require.NoError(t, err, c.statement)
		require.Equal(t, c.queryType, span.Type, c.statement)
		require.Equal(t, c.explainAnalyze, span.ExplainAnalyze, c.statement)
	}
}
