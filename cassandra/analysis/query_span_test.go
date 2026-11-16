package analysis

import (
	"errors"
	"os"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	cassandraparser "github.com/bytebase/omni/cassandra/parser"
	"github.com/bytebase/omni/metadata"
)

// corpusCase is one row of testdata/query_span.yaml. The goldens were recorded
// from the query-span extractor Bytebase shipped for Cassandra — the behavior
// this package replaces — so they are the parity bar, never re-recorded against
// this implementation. The trailing `SELECT * FROM user_summary` rows are the
// exception: they pin the later behavior change that expands a star over a view,
// so they were written against what Cassandra records for a materialized view,
// not recorded from the old extractor.
type corpusCase struct {
	Statement        string           `yaml:"statement"`
	DefaultKeyspace  string           `yaml:"defaultKeyspace"`
	WithMetadata     bool             `yaml:"withMetadata"`
	Error            bool             `yaml:"error"`
	Type             int              `yaml:"type"`
	Results          []corpusResult   `yaml:"results"`
	SourceColumns    []corpusResource `yaml:"sourcecolumns"`
	PredicateColumns []corpusResource `yaml:"predicatecolumns"`
}

type corpusResult struct {
	Name           string           `yaml:"name"`
	SourceColumns  []corpusResource `yaml:"sourcecolumns"`
	IsPlainField   bool             `yaml:"isplainfield"`
	SelectAsterisk bool             `yaml:"selectasterisk"`
}

type corpusResource struct {
	Database string `yaml:"database"`
	Table    string `yaml:"table"`
	Column   string `yaml:"column"`
}

// corpusMetadata is the keyspace fixture the withMetadata rows expand against.
// user_summary is a view: a materialized view reaches a synced Cassandra
// snapshot that way (system_schema.views), with its own column list.
func corpusMetadata() *metadata.DatabaseSchemaMetadata {
	return &metadata.DatabaseSchemaMetadata{
		Name: "test_keyspace",
		Schemas: []*metadata.SchemaMetadata{{
			Name: "",
			Tables: []*metadata.TableMetadata{
				{Name: "users", Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "name"}, {Name: "email"}, {Name: "secret"}}},
				{Name: "products", Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "name"}}},
			},
			Views: []*metadata.ViewMetadata{
				{Name: "user_summary", Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "name"}}},
			},
		}},
	}
}

func corpusOptions(withMetadata bool) ResolveOptions {
	opts := ResolveOptions{DefaultKeyspace: "test_keyspace"}
	if withMetadata {
		opts.GetDatabaseMetadata = func(keyspace string) (*metadata.DatabaseSchemaMetadata, error) {
			if keyspace != "test_keyspace" {
				return nil, nil
			}
			return corpusMetadata(), nil
		}
	}
	return opts
}

func renderResources(set ColumnSet) []corpusResource {
	out := make([]corpusResource, 0, len(set))
	for r := range set {
		out = append(out, corpusResource{Database: r.Database, Table: r.Table, Column: r.Column})
	}
	slices.SortFunc(out, func(a, b corpusResource) int {
		return compareResource(a, b)
	})
	return out
}

func compareResource(a, b corpusResource) int {
	as, bs := a.Database+"\x00"+a.Table+"\x00"+a.Column, b.Database+"\x00"+b.Table+"\x00"+b.Column
	switch {
	case as < bs:
		return -1
	case as > bs:
		return 1
	}
	return 0
}

func renderResults(results []ResolvedQuerySpanResult) []corpusResult {
	if len(results) == 0 {
		return nil
	}
	out := make([]corpusResult, 0, len(results))
	for _, r := range results {
		out = append(out, corpusResult{
			Name:           r.Name,
			SourceColumns:  renderResources(r.SourceColumns),
			IsPlainField:   r.IsPlainField,
			SelectAsterisk: r.SelectAsterisk,
		})
	}
	return out
}

// TestQuerySpanCorpus replays the recorded corpus through the package.
func TestQuerySpanCorpus(t *testing.T) {
	raw, err := os.ReadFile("testdata/query_span.yaml")
	require.NoError(t, err)
	var cases []corpusCase
	require.NoError(t, yaml.Unmarshal(raw, &cases))
	require.NotEmpty(t, cases)

	for i, tc := range cases {
		t.Run(tc.Statement, func(t *testing.T) {
			span, err := ResolveQuerySpan(tc.Statement, corpusOptions(tc.WithMetadata))
			if tc.Error {
				require.Errorf(t, err, "case %d", i+1)
				// The adapter turns the raw parse error into Bytebase's
				// SyntaxError, so the error must carry the parser's position.
				var parseErr *cassandraparser.ParseError
				require.Truef(t, errors.As(err, &parseErr), "case %d: error %T is not a *parser.ParseError", i+1, err)
				return
			}
			require.NoErrorf(t, err, "case %d", i+1)
			require.Equal(t, QueryType(tc.Type), span.Type, "case %d", i+1)
			require.Equal(t, tc.SourceColumns, renderResources(span.SourceColumns), "case %d source columns", i+1)
			require.Equal(t, tc.PredicateColumns, renderResources(span.PredicateColumns), "case %d predicate columns", i+1)
			require.Equal(t, tc.Results, renderResults(span.Results), "case %d results", i+1)
		})
	}
}

// TestQuerySpanMultiStatement: the span covers exactly one statement.
func TestQuerySpanMultiStatement(t *testing.T) {
	_, err := ResolveQuerySpan("SELECT * FROM a; SELECT * FROM b", corpusOptions(false))
	require.EqualError(t, err, "expecting only one statement to get query span, but got 2")
}

// TestQuerySpanMetadataError: a failed metadata lookup fails the span. Answering
// a SelectAsterisk marker instead would read as a statement that touches no
// column after a transient store failure, which no consumer can tell from a
// statement that genuinely reads nothing.
func TestQuerySpanMetadataError(t *testing.T) {
	errStore := errors.New("store unavailable")
	span, err := ResolveQuerySpan("SELECT * FROM users", ResolveOptions{
		DefaultKeyspace: "test_keyspace",
		GetDatabaseMetadata: func(string) (*metadata.DatabaseSchemaMetadata, error) {
			return nil, errStore
		},
	})
	require.Error(t, err)
	require.Nil(t, span)
	require.ErrorContains(t, err, `failed to get database metadata for keyspace "test_keyspace"`)
	require.ErrorIs(t, err, errStore)

	// A keyspace the caller does not hold is not an error: the star stays a
	// SelectAsterisk marker so the read is still recorded.
	span, err = ResolveQuerySpan("SELECT * FROM users", ResolveOptions{
		DefaultKeyspace: "test_keyspace",
		GetDatabaseMetadata: func(string) (*metadata.DatabaseSchemaMetadata, error) {
			return nil, nil
		},
	})
	require.NoError(t, err)
	require.Equal(t, []ResolvedQuerySpanResult{{SourceColumns: ColumnSet{}, SelectAsterisk: true}}, span.Results)
}

// TestQuerySpanStarOverView expands a star over a view the way a table's is
// expanded: from the relation's own ordered column list.
func TestQuerySpanStarOverView(t *testing.T) {
	span, err := ResolveQuerySpan("SELECT * FROM user_summary", corpusOptions(true))
	require.NoError(t, err)
	require.Equal(t, []ResolvedQuerySpanResult{
		{Name: "id", SourceColumns: ColumnSet{{Database: "test_keyspace", Table: "user_summary", Column: "id"}: true}, IsPlainField: true},
		{Name: "name", SourceColumns: ColumnSet{{Database: "test_keyspace", Table: "user_summary", Column: "name"}: true}, IsPlainField: true},
	}, span.Results)
}

// TestQuerySpanStarTableNotFound: a star over a table the keyspace metadata
// does not describe stays unexpanded rather than dropping the read.
func TestQuerySpanStarTableNotFound(t *testing.T) {
	span, err := ResolveQuerySpan("SELECT * FROM missing", corpusOptions(true))
	require.NoError(t, err)
	require.Len(t, span.Results, 1)
	require.True(t, span.Results[0].SelectAsterisk)
	require.Empty(t, span.Results[0].SourceColumns)
}
