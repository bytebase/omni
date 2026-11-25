package analysis

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/pgwire/pgcode"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/pgwire/pgerror"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/cockroachdb/internal/testcontainer"
	"github.com/bytebase/omni/enginetest/span"
	metadatapb "github.com/bytebase/omni/metadata"
)

// TestQuerySpanFidelity holds ResolveQuerySpan to CockroachDB: for
// each query of testdata/span/queries.sql, the output column names, each
// column's lineage and the relations read, as omni resolves them against the
// schema CockroachDB stored for the fixture, and as CockroachDB states them.
// testdata/README.md says where each of CockroachDB's answers comes from.
func TestQuerySpanFidelity(t *testing.T) {
	queries, err := readSpanCorpus(filepath.Join("testdata", "span", "queries.sql"))
	if err != nil {
		t.Fatal(err)
	}
	// The schema CockroachDB stored for each fixture, which omni resolves
	// against.
	stored := map[string]testcontainer.SpanFixture{}
	span.Check(t, span.Suite{
		Path:     filepath.Join("testdata", "span.tsv"),
		Pins:     testcontainer.Pins,
		PinRole:  testcontainer.Default,
		Fixtures: filepath.Join("testdata", "span", "fixtures"),
		Queries:  queries,
		Omni: func(fx span.Fixture, session, statement string) span.Answer {
			return omniSpan(stored[fx.Name], session, statement)
		},
		Engine: func(t testing.TB) func(span.Fixture, string, string) (span.Answer, error) {
			o := testcontainer.OpenSpans(t)
			built := ""
			return func(fx span.Fixture, session, statement string) (span.Answer, error) {
				if built != fx.Name {
					// The database holds one fixture at a time, so a fixture
					// is built again when the corpus returns to it.
					s, err := o.Fixture(fx.SQL)
					if err != nil {
						return span.Answer{}, fmt.Errorf("fixture %s: %w", fx.Name, err)
					}
					stored[fx.Name], built = s, fx.Name
				}
				a, err := o.Span(session, statement)
				if err != nil {
					return span.Answer{}, err
				}
				return engineSpan(stored[fx.Name], a), nil
			}
		},
	})
}

// omniSpan is ResolveQuerySpan's answer, resolved as Bytebase resolves it: in
// the fixture's database, under the session's search_path, against the
// metadata a sync stores, with names matched exactly.
func omniSpan(fx testcontainer.SpanFixture, session, statement string) span.Answer {
	meta := spanMetadata(fx)
	resolved, err := ResolveQuerySpan(context.Background(), statement, ResolveOptions{
		DefaultDatabase: meta.GetName(),
		SearchPath:      strings.Split(session, ","),
		Catalog: Catalog{GetDatabase: func(name string) (Database, error) {
			if name != meta.GetName() {
				return nil, errors.New("database not found")
			}
			return protoDatabase{meta}, nil
		}},
	})
	switch {
	case err != nil && pgerror.GetPGCode(err) != pgcode.Uncategorized:
		// A refusal for CockroachDB's reason names its SQLSTATE.
		return span.Answer{Error: pgerror.GetPGCode(err).String()}
	case err != nil:
		return span.Answer{Error: "omni"}
	case resolved.NotFoundError != nil:
		return span.Answer{Error: "omni: not found"}
	}
	var a span.Answer
	for c := range resolved.SourceColumns {
		a.Tables = append(a.Tables, analysis.ColumnResource{Database: c.Database, Schema: c.Schema, Table: c.Table}.String())
	}
	for _, r := range resolved.Results {
		a.Results = append(a.Results, r.Name)
		var sources []string
		for c := range r.SourceColumns {
			sources = append(sources, c.String())
		}
		a.Lineage = append(a.Lineage, sources)
	}
	return a
}

// spanMetadata is the stored schema as cockroachdb/driver's sync records it.
func spanMetadata(fx testcontainer.SpanFixture) *metadatapb.DatabaseSchemaMetadata {
	columns := func(cols []testcontainer.SpanColumn) []*metadatapb.ColumnMetadata {
		var out []*metadatapb.ColumnMetadata
		for _, c := range cols {
			out = append(out, &metadatapb.ColumnMetadata{Name: c.Name, Type: c.Type, Nullable: true, IsInvisible: c.Hidden})
		}
		return out
	}
	meta := &metadatapb.DatabaseSchemaMetadata{Name: fx.Database}
	for _, s := range fx.Schemas {
		schema := &metadatapb.SchemaMetadata{Name: s.Name}
		for _, t := range s.Tables {
			schema.Tables = append(schema.Tables, &metadatapb.TableMetadata{Name: t.Name, Columns: columns(t.Columns)})
		}
		for _, v := range s.Views {
			schema.Views = append(schema.Views, &metadatapb.ViewMetadata{Name: v.Name, Definition: v.Definition, Columns: columns(v.Columns)})
		}
		for _, v := range s.MaterializedViews {
			schema.MaterializedViews = append(schema.MaterializedViews, &metadatapb.MaterializedViewMetadata{Name: v.Name, Definition: v.Definition})
		}
		meta.Schemas = append(meta.Schemas, schema)
	}
	return meta
}

// engineSpan is CockroachDB's answer. Its RowDescription names the table
// column of an output column that is a plain column reference, followed
// through subqueries, joins and views, and nothing for any other: a column
// it names nothing for is span.Unstated. Through a view it names the table
// column behind the view's column and never the view's column, so a column
// whose origin is also the origin of a column of a view the query reads
// names that origin and span.Unstated. The relations of a statement no view
// can hold are span.Unstated.
func engineSpan(fx testcontainer.SpanFixture, a testcontainer.SpanAnswer) span.Answer {
	if a.Error != "" {
		return span.Answer{Error: a.Error}
	}
	throughView := map[string]bool{}
	for _, s := range fx.Schemas {
		for _, v := range s.Views {
			read := false
			for _, table := range a.Tables {
				read = read || table == fx.Database+"."+s.Name+"."+v.Name
			}
			for _, c := range v.Columns {
				if read && c.Origin != "" {
					throughView[c.Origin] = true
				}
			}
		}
	}
	out := span.Answer{Tables: a.Tables}
	if a.TablesUnstated {
		out.Tables = []string{span.Unstated}
	}
	for _, c := range a.Columns {
		out.Results = append(out.Results, c.Name)
		switch {
		case c.Source == "":
			out.Lineage = append(out.Lineage, []string{span.Unstated})
		case throughView[c.Source]:
			out.Lineage = append(out.Lineage, []string{c.Source, span.Unstated})
		default:
			out.Lineage = append(out.Lineage, []string{c.Source})
		}
	}
	return out
}

// protoDatabase answers ResolveQuerySpan's lookups from the metadata, with
// names matched exactly.
type protoDatabase struct {
	meta *metadatapb.DatabaseSchemaMetadata
}

func (d protoDatabase) Metadata() *metadatapb.DatabaseSchemaMetadata { return d.meta }

func (d protoDatabase) schema(name string) *metadatapb.SchemaMetadata {
	for _, s := range d.meta.GetSchemas() {
		if s.GetName() == name {
			return s
		}
	}
	return nil
}

func (d protoDatabase) GetTable(schema, table string) *metadatapb.TableMetadata {
	for _, t := range d.schema(schema).GetTables() {
		if t.GetName() == table {
			return t
		}
	}
	return nil
}

func (d protoDatabase) SearchRelation(searchPath []string, name string) (string, string) {
	for _, schemaName := range searchPath {
		s := d.schema(schemaName)
		if s == nil {
			continue
		}
		for _, t := range s.GetTables() {
			if t.GetName() == name {
				return schemaName, name
			}
		}
		for _, v := range s.GetViews() {
			if v.GetName() == name {
				return schemaName, name
			}
		}
		for _, v := range s.GetMaterializedViews() {
			if v.GetName() == name {
				return schemaName, name
			}
		}
	}
	return "", ""
}

// readSpanCorpus reads queries.sql: one query per line, on the fixture and
// under the search_path the last "-- fixture:" and "-- session:" lines name.
// Other lines starting "--" are comments.
func readSpanCorpus(path string) ([]span.Query, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var (
		out              []span.Query
		fixture, session string
	)
	sc := bufio.NewScanner(f)
	for n := 1; sc.Scan(); n++ {
		line := strings.TrimSpace(sc.Text())
		switch {
		case line == "":
		case strings.HasPrefix(line, "-- fixture:"):
			fixture = strings.TrimSpace(strings.TrimPrefix(line, "-- fixture:"))
		case strings.HasPrefix(line, "-- session:"):
			session = strings.TrimSpace(strings.TrimPrefix(line, "-- session:"))
		case strings.HasPrefix(line, "--"):
		case fixture == "" || session == "":
			return nil, fmt.Errorf("%s:%d: a query before the -- fixture: and -- session: lines", path, n)
		default:
			out = append(out, span.Query{Fixture: fixture, Session: session, Statement: strings.TrimSuffix(line, ";")})
		}
	}
	return out, sc.Err()
}
