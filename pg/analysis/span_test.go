package analysis

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest/span"
	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/pg/internal/testcontainer"
)

// TestQuerySpanFidelity holds ResolveQuerySpan to PostgreSQL: for each query
// of testdata/span/queries.sql, the output column names, each column's
// lineage and the relations read, as omni resolves them against the schema
// PostgreSQL stored for the fixture, and as PostgreSQL states them.
// testdata/README.md says where each of PostgreSQL's answers comes from.
func TestQuerySpanFidelity(t *testing.T) {
	queries, err := readSpanCorpus(filepath.Join("testdata", "span", "queries.sql"))
	if err != nil {
		t.Fatal(err)
	}
	// The schema PostgreSQL stored for each fixture, which omni resolves
	// against.
	stored := map[string]testcontainer.SpanFixture{}
	span.Check(t, span.Suite{
		Path:     filepath.Join("testdata", "span.tsv"),
		Pins:     testcontainer.Pins,
		PinRole:  testcontainer.PG18,
		Fixtures: filepath.Join("testdata", "span", "fixtures"),
		Queries:  queries,
		Omni: func(fx span.Fixture, session, statement string) span.Answer {
			return omniSpan(stored[fx.Name], session, statement)
		},
		Engine: func(t testing.TB) func(span.Fixture, string, string) (span.Answer, error) {
			o := testcontainer.OpenSpans(t, testcontainer.PG18)
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
				return engineSpan(a), nil
			}
		},
	})
}

// omniSpan is ResolveQuerySpan's answer, resolved as Bytebase resolves it for
// PostgreSQL: in the fixture's database, under the session's search_path,
// against the metadata a sync stores, with names matched exactly.
func omniSpan(fx testcontainer.SpanFixture, session, statement string) span.Answer {
	meta := spanMetadata(fx)
	resolved, err := ResolveQuerySpan(context.Background(), statement, ResolveOptions{
		DefaultDatabase: meta.GetName(),
		SearchPath:      strings.Split(session, ","),
		Catalog:         testCatalog(meta),
	})
	switch {
	case err != nil:
		return span.Answer{Error: "omni"}
	case resolved.NotFoundError != nil:
		return span.Answer{Error: "omni: not found"}
	}
	var a span.Answer
	for c := range resolved.SourceColumns {
		a.Tables = append(a.Tables, ColumnResource{Database: c.Database, Schema: c.Schema, Table: c.Table}.String())
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

// spanMetadata is the stored schema as Bytebase's sync records it.
func spanMetadata(fx testcontainer.SpanFixture) *metadatapb.DatabaseSchemaMetadata {
	columns := func(cols []testcontainer.SpanColumn) []*metadatapb.ColumnMetadata {
		var out []*metadatapb.ColumnMetadata
		for _, c := range cols {
			out = append(out, &metadatapb.ColumnMetadata{Name: c.Name, Type: c.Type, Nullable: true})
		}
		return out
	}
	meta := &metadatapb.DatabaseSchemaMetadata{Name: fx.Database}
	for _, s := range fx.Schemas {
		schema := &metadatapb.SchemaMetadata{Name: s.Name}
		for _, t := range s.Tables {
			table := &metadatapb.TableMetadata{Name: t.Name, Columns: columns(t.Columns)}
			for _, p := range t.Partitions {
				table.Partitions = append(table.Partitions, &metadatapb.TablePartitionMetadata{Name: p})
			}
			schema.Tables = append(schema.Tables, table)
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

// engineSpan is PostgreSQL's answer. Its RowDescription names the table
// column of an output column that is a plain column reference, followed
// through subqueries, common table expressions and joins, and nothing for
// any other: a column it names nothing for is span.Unstated. A whole-row
// reference names its relation and attribute number 0, which stands for
// every column of the relation.
func engineSpan(a testcontainer.SpanAnswer) span.Answer {
	if a.Error != "" {
		return span.Answer{Error: a.Error}
	}
	out := span.Answer{Tables: a.Tables}
	for _, c := range a.Columns {
		out.Results = append(out.Results, c.Name)
		sources := c.Sources
		if len(sources) == 0 {
			sources = []string{span.Unstated}
		}
		out.Lineage = append(out.Lineage, sources)
	}
	return out
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
