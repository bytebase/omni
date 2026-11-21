package analysis

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v3"
)

const (
	corpusInstanceIDA = "INSTANCE_ID_A"
	corpusInstanceIDB = "INSTANCE_ID_B"
)

// The corpus under testdata/query-span is Bytebase's Oracle query-span fixture,
// recorded from the resolver masking ran on before this package existed. It is
// the masking-parity bar for this package and must be reproduced, never
// re-recorded against this implementation.

type spanCorpusCase struct {
	Description           string          `yaml:"description"`
	Statement             string          `yaml:"statement"`
	DefaultDatabase       string          `yaml:"defaultDatabase"`
	Metadata              string          `yaml:"metadata"`
	CrossDatabaseMetadata string          `yaml:"crossDatabaseMetadata"`
	QuerySpan             *spanCorpusSpan `yaml:"querySpan"`
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
	Instance string `yaml:"instance,omitempty"`
	Server   string `yaml:"server"`
	Database string `yaml:"database"`
	Schema   string `yaml:"schema"`
	Table    string `yaml:"table"`
	Column   string `yaml:"column"`
}

type corpusProtoDatabase struct {
	proto *metadatapb.DatabaseSchemaMetadata
}

func (d *corpusProtoDatabase) schema() *metadatapb.SchemaMetadata {
	for _, s := range d.proto.GetSchemas() {
		if s.GetName() == "" {
			return s
		}
	}
	return nil
}

func (d *corpusProtoDatabase) HasSchema() bool { return d.schema() != nil }

func (d *corpusProtoDatabase) GetTable(name string) *metadatapb.TableMetadata {
	for _, t := range d.schema().GetTables() {
		if t.GetName() == name {
			return t
		}
	}
	return nil
}

func (d *corpusProtoDatabase) GetView(name string) *metadatapb.ViewMetadata {
	for _, v := range d.schema().GetViews() {
		if v.GetName() == name {
			return v
		}
	}
	return nil
}

func (d *corpusProtoDatabase) GetMaterializedView(name string) *metadatapb.MaterializedViewMetadata {
	for _, mv := range d.schema().GetMaterializedViews() {
		if mv.GetName() == name {
			return mv
		}
	}
	return nil
}

func (d *corpusProtoDatabase) GetExternalTable(name string) *metadatapb.ExternalTableMetadata {
	for _, et := range d.schema().GetExternalTables() {
		if et.GetName() == name {
			return et
		}
	}
	return nil
}

func linkedCorpusMetadata() []*metadatapb.DatabaseSchemaMetadata {
	return []*metadatapb.DatabaseSchemaMetadata{{
		Name: "SCHEMA1",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "",
			Tables: []*metadatapb.TableMetadata{
				{
					Name:    "LT1",
					Columns: []*metadatapb.ColumnMetadata{{Name: "LC1", Type: "int"}},
				},
				{
					Name: "LT2",
					Columns: []*metadatapb.ColumnMetadata{
						{Name: "LC1", Type: "int"},
						{Name: "LC2", Type: "int"},
					},
				},
			},
			Views: []*metadatapb.ViewMetadata{{
				Name:       "LV1",
				Definition: "SELECT LC1, LC2\nFROM LT2\n",
			}},
		}},
	}}
}

func corpusCatalog(defaultMetadata []*metadatapb.DatabaseSchemaMetadata) Catalog {
	return Catalog{
		GetDatabase: func(instanceID, databaseName string) (Database, error) {
			list := defaultMetadata
			if instanceID == corpusInstanceIDB {
				list = linkedCorpusMetadata()
			}
			for _, m := range list {
				if m.GetName() == databaseName {
					return &corpusProtoDatabase{proto: m}, nil
				}
			}
			return nil, fmt.Errorf("database %q not found", databaseName)
		},
		GetLinkedDatabase: func(_, linkedDatabaseName, schemaName string) (string, string, Database, error) {
			var linkedDBInfo *metadatapb.LinkedDatabaseMetadata
			for _, m := range defaultMetadata {
				for _, l := range m.GetLinkedDatabases() {
					if l.GetName() == linkedDatabaseName {
						linkedDBInfo = l
						break
					}
				}
				if linkedDBInfo != nil {
					break
				}
			}
			if linkedDBInfo == nil {
				return "", "", nil, nil
			}
			databaseName := linkedDBInfo.GetUsername()
			if schemaName != "" {
				databaseName = schemaName
			}
			candidates, instanceID := linkedCorpusMetadata(), corpusInstanceIDB
			if linkedDBInfo.GetName() == "LOOPBACK" {
				candidates, instanceID = defaultMetadata, corpusInstanceIDA
			}
			for _, m := range candidates {
				if m.GetName() == databaseName {
					return instanceID, m.GetName(), &corpusProtoDatabase{proto: m}, nil
				}
			}
			return "", "", nil, nil
		},
	}
}

func renderSpanColumns(set ColumnSet) []spanCorpusColumn {
	out := make([]spanCorpusColumn, 0, len(set))
	for resource := range set {
		out = append(out, spanCorpusColumn{
			Instance: resource.Instance,
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
	content, err := os.ReadFile("testdata/query-span/query_span.yaml")
	require.NoError(t, err)
	var cases []spanCorpusCase
	require.NoError(t, yaml.Unmarshal(content, &cases))
	require.NotEmpty(t, cases)

	unmarshaler := protojson.UnmarshalOptions{DiscardUnknown: true}
	for i, tc := range cases {
		name := tc.Description
		if name == "" {
			name = tc.Statement
		}
		t.Run(name, func(t *testing.T) {
			require.NotNil(t, tc.QuerySpan, "case %d has no querySpan", i)
			dbMetadata := &metadatapb.DatabaseSchemaMetadata{}
			require.NoError(t, unmarshaler.Unmarshal([]byte(tc.Metadata), dbMetadata))
			list := []*metadatapb.DatabaseSchemaMetadata{dbMetadata}
			if tc.CrossDatabaseMetadata != "" {
				crossDB := &metadatapb.DatabaseSchemaMetadata{}
				require.NoError(t, unmarshaler.Unmarshal([]byte(tc.CrossDatabaseMetadata), crossDB))
				list = append(list, crossDB)
			}
			span, err := ResolveQuerySpan(tc.Statement, ResolveOptions{
				Instance:        corpusInstanceIDA,
				DefaultDatabase: tc.DefaultDatabase,
				Catalog:         corpusCatalog(list),
			})
			require.NoError(t, err, "statement: %s", tc.Statement)
			require.NotNil(t, span)
			require.Equal(t, tc.QuerySpan, renderSpan(span), "statement: %s", tc.Statement)
		})
	}
}

func TestResolveQuerySpanRefusesSystemTableWithLink(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name:            "PUBLIC",
		Schemas:         []*metadatapb.SchemaMetadata{{Name: "", Tables: []*metadatapb.TableMetadata{{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "A"}}}}}},
		LinkedDatabases: []*metadatapb.LinkedDatabaseMetadata{{Name: "REMOTE", Username: "SCHEMA1"}},
	}
	opts := ResolveOptions{
		Instance:        corpusInstanceIDA,
		DefaultDatabase: "PUBLIC",
		Catalog:         corpusCatalog([]*metadatapb.DatabaseSchemaMetadata{meta}),
	}
	_, err := ResolveQuerySpan("SELECT * FROM SYS.DBA_USERS, T@REMOTE;", opts)
	require.ErrorIs(t, err, ErrMixUserSystemTables)

	span, err := ResolveQuerySpan("SELECT * FROM SYS.DBA_USERS@REMOTE;", opts)
	require.NoError(t, err)
	require.Equal(t, Select, span.Type)
}

func TestResolveQuerySpanParenthesizedQuery(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name:    "PUBLIC",
		Schemas: []*metadatapb.SchemaMetadata{{Name: "", Tables: []*metadatapb.TableMetadata{{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "A"}, {Name: "B"}}}}}},
	}
	opts := ResolveOptions{
		Instance:        corpusInstanceIDA,
		DefaultDatabase: "PUBLIC",
		Catalog:         corpusCatalog([]*metadatapb.DatabaseSchemaMetadata{meta}),
	}
	columnA := ColumnResource{Database: "PUBLIC", Table: "T", Column: "A"}
	columnB := ColumnResource{Database: "PUBLIC", Table: "T", Column: "B"}
	tests := []struct {
		statement string
		want      ColumnSet
	}{
		{statement: "(SELECT a FROM t)", want: ColumnSet{columnA: true}},
		{statement: "(SELECT a FROM t) ORDER BY a", want: ColumnSet{columnA: true}},
		{statement: "((SELECT a FROM t)) ORDER BY 1", want: ColumnSet{columnA: true}},
		{statement: "(SELECT a FROM t) FETCH FIRST 5 ROWS ONLY", want: ColumnSet{columnA: true}},
		{statement: "WITH w AS (SELECT a FROM t) (SELECT a FROM w)", want: ColumnSet{columnA: true}},
		{statement: "(SELECT a FROM t) UNION (SELECT b FROM t) ORDER BY 1", want: ColumnSet{columnA: true, columnB: true}},
		{statement: "SELECT a FROM t EXCEPT SELECT b FROM t", want: ColumnSet{columnA: true, columnB: true}},
	}
	for _, tc := range tests {
		t.Run(tc.statement, func(t *testing.T) {
			span, err := ResolveQuerySpan(tc.statement, opts)
			require.NoError(t, err)
			require.Len(t, span.Results, 1)
			require.Equal(t, "A", span.Results[0].Name)
			require.Equal(t, tc.want, span.Results[0].SourceColumns)
		})
	}
}
