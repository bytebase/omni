package analysis_test

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v3"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/snowflake/analysis"
)

type spanCorpusCase struct {
	Description         string          `yaml:"description"`
	Statement           string          `yaml:"statement"`
	DefaultDatabase     string          `yaml:"defaultDatabase"`
	IgnoreCaseSensitive bool            `yaml:"ignoreCaseSensitive"`
	Metadata            string          `yaml:"metadata"`
	QuerySpan           *spanCorpusSpan `yaml:"querySpan"`
}

type spanCorpusSpan struct {
	Type             int                `yaml:"type"`
	Results          []spanCorpusResult `yaml:"results"`
	SourceColumns    []spanCorpusColumn `yaml:"sourcecolumns"`
	PredicateColumns []spanCorpusColumn `yaml:"predicatecolumns"`
}

type spanCorpusResult struct {
	Name           string             `yaml:"name"`
	SourceColumns  []spanCorpusColumn `yaml:"sourcecolumns"`
	IsPlainField   bool               `yaml:"isplainfield"`
	SelectAsterisk bool               `yaml:"selectasterisk"`
}

type spanCorpusColumn struct {
	Server   string `yaml:"server"`
	Database string `yaml:"database"`
	Schema   string `yaml:"schema"`
	Table    string `yaml:"table"`
	Column   string `yaml:"column"`
}

func resolveFromMetadata(db *metadata.DatabaseSchemaMetadata) func(database, schema, table string) (*analysis.Table, error) {
	return func(database, schema, table string) (*analysis.Table, error) {
		if database != "" && database != db.GetName() {
			return nil, fmt.Errorf("table %s.%s.%s is not found", database, schema, table)
		}
		for _, s := range db.GetSchemas() {
			if schema != "" && schema != s.GetName() {
				continue
			}
			for _, tbl := range s.GetTables() {
				if tbl.GetName() != table {
					continue
				}
				cols := make([]string, 0, len(tbl.GetColumns()))
				for _, c := range tbl.GetColumns() {
					cols = append(cols, c.GetName())
				}
				return &analysis.Table{Name: tbl.GetName(), Columns: cols}, nil
			}
		}
		return nil, fmt.Errorf("table %s.%s.%s is not found", database, schema, table)
	}
}

func renderSpanColumns(set analysis.ColumnSet) []spanCorpusColumn {
	out := make([]spanCorpusColumn, 0, len(set))
	for r := range set {
		out = append(out, spanCorpusColumn{
			Database: r.Database,
			Schema:   r.Schema,
			Table:    r.Table,
			Column:   r.Column,
		})
	}
	slices.SortFunc(out, func(a, b spanCorpusColumn) int {
		return strings.Compare(spanColumnKey(a), spanColumnKey(b))
	})
	return out
}

func spanColumnKey(c spanCorpusColumn) string {
	parts := make([]string, 0, 5)
	for _, p := range []string{c.Server, c.Database, c.Schema, c.Table, c.Column} {
		if p != "" {
			parts = append(parts, p)
		}
	}
	return strings.Join(parts, ".")
}

func renderCorpusSpan(span *analysis.ResolvedSpan) *spanCorpusSpan {
	out := &spanCorpusSpan{
		Type:             int(span.Type),
		Results:          make([]spanCorpusResult, 0, len(span.Results)),
		SourceColumns:    renderSpanColumns(span.SourceColumns),
		PredicateColumns: []spanCorpusColumn{},
	}
	for _, r := range span.Results {
		out.Results = append(out.Results, spanCorpusResult{
			Name:          r.Name,
			SourceColumns: renderSpanColumns(r.SourceColumns),
			IsPlainField:  r.IsPlainField,
		})
	}
	return out
}

func TestResolveQuerySpanCorpus(t *testing.T) {
	paths := []string{
		"testdata/query-span/standard.yaml",
		"testdata/query-span/subquery.yaml",
		"testdata/query-span/set-operator.yaml",
		"testdata/query-span/cte.yaml",
	}

	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			raw, err := os.ReadFile(path)
			require.NoError(t, err)
			var cases []spanCorpusCase
			require.NoError(t, yaml.Unmarshal(raw, &cases))
			require.NotEmpty(t, cases)

			for i, tc := range cases {
				t.Run(fmt.Sprintf("%d_%s", i+1, tc.Description), func(t *testing.T) {
					dbMetadata := &metadata.DatabaseSchemaMetadata{}
					require.NoErrorf(t, protojson.Unmarshal([]byte(tc.Metadata), dbMetadata), "case %d", i+1)

					span, err := analysis.ResolveQuerySpan(tc.Statement, analysis.ResolveOptions{
						DefaultDatabase:     tc.DefaultDatabase,
						DefaultSchema:       "PUBLIC",
						IgnoreCaseSensitive: tc.IgnoreCaseSensitive,
						ResolveTable:        resolveFromMetadata(dbMetadata),
					})
					require.NoErrorf(t, err, "statement: %s", tc.Statement)
					require.Equalf(t, tc.QuerySpan, renderCorpusSpan(span), "statement: %s", tc.Statement)
				})
			}
		})
	}
}

// TestResolveQuerySpanPivotCorpusFailsClosed preserves testdata/query-span/pivot.yaml
// for future pivot projection support while pinning the current fail-closed contract:
// ResolveQuerySpan must return an explicit PIVOT/UNPIVOT error rather than mis-attributing.
func TestResolveQuerySpanPivotCorpusFailsClosed(t *testing.T) {
	raw, err := os.ReadFile("testdata/query-span/pivot.yaml")
	require.NoError(t, err)
	var cases []spanCorpusCase
	require.NoError(t, yaml.Unmarshal(raw, &cases))
	require.NotEmpty(t, cases)

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d_%s", i+1, tc.Description), func(t *testing.T) {
			dbMetadata := &metadata.DatabaseSchemaMetadata{}
			require.NoError(t, protojson.Unmarshal([]byte(tc.Metadata), dbMetadata))

			_, err := analysis.ResolveQuerySpan(tc.Statement, analysis.ResolveOptions{
				DefaultDatabase:     tc.DefaultDatabase,
				DefaultSchema:       "PUBLIC",
				IgnoreCaseSensitive: tc.IgnoreCaseSensitive,
				ResolveTable:        resolveFromMetadata(dbMetadata),
			})
			require.Errorf(t, err, "expected fail-closed error for %s", tc.Statement)
			require.Containsf(t, err.Error(), "PIVOT", "statement: %s", tc.Statement)
		})
	}
}
