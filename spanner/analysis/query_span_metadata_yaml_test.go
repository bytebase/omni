package analysis

import (
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v3"

	"github.com/bytebase/omni/metadata"
)

// The corpora under testdata/query-span were recorded from the two legacy
// ANTLR resolvers Bytebase shipped for BigQuery and Spanner. They are the
// masking-parity bar for this package and must be reproduced, never re-recorded
// against this implementation. The one departure is where the legacy resolver
// dropped a source column, which masking needs: a scalar subquery in the select
// list reads its subquery's columns, so those cases list them.

type spanCorpusCase struct {
	Description     string `yaml:"description"`
	Statement       string `yaml:"statement"`
	DefaultDatabase string `yaml:"defaultDatabase"`
	// IgnoreCaseSensitive is recorded but unused: both legacy extractors
	// ignored it (schema names match case-sensitively under the spanner model,
	// table names fall back to a case-folded match — the legacy
	// findTableSchema behavior).
	IgnoreCaseSensitive bool `yaml:"ignoreCaseSensitive"`
	// Metadata is the protojson encoded metadata.DatabaseSchemaMetadata.
	Metadata  string          `yaml:"metadata"`
	QuerySpan *spanCorpusSpan `yaml:"querySpan"`
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

// spannerResolveOptions mirrors the flags Bytebase's Spanner engine registers.
func spannerResolveOptions() ResolveOptions {
	return ResolveOptions{
		IgnoreWrittenDatabase:  true,
		SchemaScopedMetadata:   true,
		SystemOnlyEmptySpan:    true,
		SetStatementIsSelect:   true,
		CanonicalBaseFieldName: true,
	}
}

// corpusGetter serves exactly the case's metadata, keyed by database name; an
// unlisted database is an error, matching the caller-side mock the corpora were
// recorded with.
func corpusGetter(dbMetadata *metadata.DatabaseSchemaMetadata) func(string) (*metadata.DatabaseSchemaMetadata, error) {
	return func(name string) (*metadata.DatabaseSchemaMetadata, error) {
		if name == dbMetadata.GetName() {
			return dbMetadata, nil
		}
		return nil, &ResourceNotFoundError{Database: &name}
	}
}

func renderColumns(set ColumnSet) []spanCorpusColumn {
	out := make([]spanCorpusColumn, 0, len(set))
	for r := range set {
		out = append(out, spanCorpusColumn{Database: r.Database, Schema: r.Schema, Table: r.Table, Column: r.Column})
	}
	slices.SortFunc(out, func(a, b spanCorpusColumn) int {
		return strings.Compare(columnKey(a), columnKey(b))
	})
	return out
}

// columnKey is the corpus's sort key: the dot-joined non-empty fields, which is
// how the recorded YAML lists each set.
func columnKey(c spanCorpusColumn) string {
	parts := make([]string, 0, 5)
	for _, p := range []string{c.Server, c.Database, c.Schema, c.Table, c.Column} {
		if p != "" {
			parts = append(parts, p)
		}
	}
	return strings.Join(parts, ".")
}

func renderSpan(span *ResolvedQuerySpan) *spanCorpusSpan {
	out := &spanCorpusSpan{
		Type:             int(span.Type),
		Results:          make([]spanCorpusResult, 0, len(span.Results)),
		SourceColumns:    renderColumns(span.SourceColumns),
		PredicateColumns: renderColumns(span.PredicateColumns),
	}
	for _, r := range span.Results {
		out.Results = append(out.Results, spanCorpusResult{
			Name:          r.Name,
			SourceColumns: renderColumns(r.SourceColumns),
			IsPlainField:  r.IsPlainField,
		})
	}
	return out
}

func runSpanCorpus(t *testing.T, path string, opts ResolveOptions) int {
	t.Helper()
	raw, err := os.ReadFile(path)
	require.NoError(t, err)
	var cases []spanCorpusCase
	require.NoError(t, yaml.Unmarshal(raw, &cases))
	require.NotEmpty(t, cases)

	for i, tc := range cases {
		t.Run(tc.Description, func(t *testing.T) {
			dbMetadata := &metadata.DatabaseSchemaMetadata{}
			require.NoErrorf(t, protojson.Unmarshal([]byte(tc.Metadata), dbMetadata), "case %d", i+1)

			caseOpts := opts
			caseOpts.DefaultDatabase = tc.DefaultDatabase
			caseOpts.GetDatabaseMetadata = corpusGetter(dbMetadata)
			span, err := ResolveQuerySpan(tc.Statement, caseOpts)
			require.NoErrorf(t, err, "case %d %q", i+1, tc.Statement)
			require.Equal(t, tc.QuerySpan, renderSpan(span), "case %d %q", i+1, tc.Statement)
		})
	}
	return len(cases)
}

// TestQuerySpanCorpusSpanner replays the Spanner corpus.
func TestQuerySpanCorpusSpanner(t *testing.T) {
	n := runSpanCorpus(t, "testdata/query-span/spanner.yaml", spannerResolveOptions())
	t.Logf("replayed %d Spanner corpus cases", n)
}
