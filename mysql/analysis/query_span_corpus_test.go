package analysis

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v3"

	"github.com/bytebase/omni/metadata"
)

// The corpus under testdata/query-span is Bytebase's MySQL query-span fixture,
// recorded from the resolver masking ran on before this package existed. It is
// the masking-parity bar for this package and must be reproduced, never
// re-recorded against this implementation. An entry changes only where
// testdata/span.tsv shows MySQL answers otherwise (testdata/README.md).
//
// The fixture carries three fields this engine never sets —
// sourcefieldpaths, predicatepaths and selectasterisk — because the shape is
// shared with the engines whose resolvers report them. They are kept so the
// corpus stays byte-identical to the one Bytebase runs; each asserts a value
// this engine must keep producing (empty and false).

type spanCorpusCase struct {
	Description         string `yaml:"description"`
	Statement           string `yaml:"statement"`
	DefaultDatabase     string `yaml:"defaultDatabase"`
	IgnoreCaseSensitive bool   `yaml:"ignoreCaseSensitive"`
	// Metadata is the protojson encoded metadata.DatabaseSchemaMetadata.
	Metadata  string          `yaml:"metadata"`
	QuerySpan *spanCorpusSpan `yaml:"querySpan"`
}

type spanCorpusSpan struct {
	Type             int                `yaml:"type"`
	Results          []spanCorpusResult `yaml:"results"`
	SourceColumns    []spanCorpusColumn `yaml:"sourcecolumns"`
	PredicateColumns []spanCorpusColumn `yaml:"predicatecolumns"`
	PredicatePaths   []string           `yaml:"predicatepaths"`
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

// corpusDatabase serves one case's metadata the way Bytebase's catalog adapter
// does: MySQL's tables and views live in the database's unnamed schema, and a
// database whose metadata has no such schema has no tables.
type corpusDatabase struct {
	proto *metadata.DatabaseSchemaMetadata
}

func (d *corpusDatabase) Name() string { return d.proto.GetName() }

func (d *corpusDatabase) HasSchema() bool { return d.schema() != nil }

func (d *corpusDatabase) ListTableNames() []string {
	var names []string
	for _, table := range d.schema().GetTables() {
		names = append(names, table.GetName())
	}
	return names
}

func (d *corpusDatabase) GetTable(name string) *metadata.TableMetadata {
	for _, table := range d.schema().GetTables() {
		if table.GetName() == name {
			return table
		}
	}
	return nil
}

func (d *corpusDatabase) ListViewNames() []string {
	var names []string
	for _, view := range d.schema().GetViews() {
		names = append(names, view.GetName())
	}
	return names
}

func (d *corpusDatabase) GetView(name string) *metadata.ViewMetadata {
	for _, view := range d.schema().GetViews() {
		if view.GetName() == name {
			return view
		}
	}
	return nil
}

func (d *corpusDatabase) schema() *metadata.SchemaMetadata {
	for _, schema := range d.proto.GetSchemas() {
		if schema.GetName() == "" {
			return schema
		}
	}
	return nil
}

// corpusCatalog serves exactly the case's metadata, keyed by database name. An
// unlisted database is a not-found error, as the caller-side mock the corpus
// was recorded with answers.
func corpusCatalog(dbMetadata *metadata.DatabaseSchemaMetadata) Catalog {
	return Catalog{
		ListDatabaseNames: func() ([]string, error) {
			return []string{dbMetadata.GetName()}, nil
		},
		GetDatabase: func(name string) (Database, error) {
			if name == dbMetadata.GetName() {
				return &corpusDatabase{proto: dbMetadata}, nil
			}
			return nil, &NotFoundError{Database: &name}
		},
	}
}

func renderColumns(set ColumnSet) []spanCorpusColumn {
	out := make([]spanCorpusColumn, 0, len(set))
	for resource := range set {
		out = append(out, spanCorpusColumn{
			Database: resource.Database,
			Schema:   resource.Schema,
			Table:    resource.Table,
			Column:   resource.Column,
		})
	}
	slices.SortFunc(out, func(a, b spanCorpusColumn) int {
		return strings.Compare(columnKey(a), columnKey(b))
	})
	return out
}

// columnKey is the corpus's sort key: the dot-joined non-empty fields, which is
// how the recorded YAML lists each set.
func columnKey(c spanCorpusColumn) string {
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
		SourceColumns:    renderColumns(span.SourceColumns),
		PredicateColumns: renderColumns(span.PredicateColumns),
	}
	for _, result := range span.Results {
		out.Results = append(out.Results, spanCorpusResult{
			Name:             result.Name,
			SourceColumns:    renderColumns(result.SourceColumns),
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
	for _, path := range paths {
		content, err := os.ReadFile(path)
		require.NoError(t, err)
		var cases []spanCorpusCase
		require.NoError(t, yaml.Unmarshal(content, &cases), path)
		for _, tc := range cases {
			t.Run(filepath.Base(path)+"/"+tc.Description, func(t *testing.T) {
				require.NotNil(t, tc.QuerySpan, "case has no querySpan")
				dbMetadata := &metadata.DatabaseSchemaMetadata{}
				require.NoError(t, protojson.Unmarshal([]byte(tc.Metadata), dbMetadata))
				span, err := ResolveQuerySpan(tc.Statement, ResolveOptions{
					DefaultDatabase:     tc.DefaultDatabase,
					IgnoreCaseSensitive: tc.IgnoreCaseSensitive,
					Catalog:             corpusCatalog(dbMetadata),
				})
				require.NoError(t, err)
				require.Equal(t, tc.QuerySpan, renderSpan(span))
			})
		}
	}
}
