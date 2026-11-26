package analysis

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/bigquery/internal/testcontainer"
	"github.com/bytebase/omni/enginetest/span"
	"github.com/bytebase/omni/metadata"
)

// TestQuerySpanFidelity holds ResolveQuerySpan to BigQuery: for each query of
// testdata/span/queries.sql, the output column names and the tables read, as
// omni resolves them against the relations BigQuery stored for the fixture,
// and as a dry run of the query states them. testdata/README.md says where
// each of BigQuery's answers comes from.
func TestQuerySpanFidelity(t *testing.T) {
	queries, err := readSpanCorpus(filepath.Join("testdata", "span", "queries.sql"))
	if err != nil {
		t.Fatal(err)
	}
	// The relations BigQuery stored for each fixture, which omni resolves
	// against.
	stored := map[string][]testcontainer.SpanRelation{}
	span.Check(t, span.Suite{
		Path:     filepath.Join("testdata", "span.tsv"),
		Pins:     testcontainer.Pins,
		PinRole:  testcontainer.Span,
		Fixtures: filepath.Join("testdata", "span", "fixtures"),
		Queries:  queries,
		Omni: func(fx span.Fixture, session, statement string) span.Answer {
			return omniSpan(stored[fx.Name], session, statement)
		},
		Engine: func(t testing.TB) func(span.Fixture, string, string) (span.Answer, error) {
			o := testcontainer.OpenSpan(t)
			built := ""
			return func(fx span.Fixture, _, statement string) (span.Answer, error) {
				if built != fx.Name {
					// The dataset holds one fixture at a time, so a fixture
					// is built again when the corpus returns to it.
					a, err := o.Fixture(fx.SQL)
					if err != nil {
						return span.Answer{}, fmt.Errorf("fixture %s: %w", fx.Name, err)
					}
					stored[fx.Name], built = a.Relations, fx.Name
				}
				a, err := o.Query(statement)
				if err != nil {
					return span.Answer{}, err
				}
				return engineSpan(a), nil
			}
		},
	})
}

// omniSpan is ResolveQuerySpan's answer, resolved as Bytebase resolves it for
// BigQuery, with the session's dataset as the default one, against the
// metadata its sync stores.
func omniSpan(relations []testcontainer.SpanRelation, session, statement string) span.Answer {
	meta := spanMetadata(session, relations)
	opts := bigqueryResolveOptions()
	opts.DefaultDatabase = session
	opts.GetDatabaseMetadata = func(database string) (*metadata.DatabaseSchemaMetadata, error) {
		if database == session {
			return meta, nil
		}
		return nil, &ResourceNotFoundError{Database: &database}
	}
	resolved, err := ResolveQuerySpan(statement, opts)
	if err != nil {
		return span.Answer{Error: "omni"}
	}
	var a span.Answer
	tables := map[string]bool{}
	for c := range resolved.SourceColumns {
		tables[qualifySpan(c)] = true
	}
	for name := range tables {
		a.Tables = append(a.Tables, name)
	}
	for _, r := range resolved.Results {
		a.Results = append(a.Results, r.Name)
		var sources []string
		for c := range r.SourceColumns {
			sources = append(sources, qualifySpan(c)+"."+c.Column)
		}
		a.Lineage = append(a.Lineage, sources)
	}
	return a
}

// qualifySpan names a table as BigQuery's dry run does: dataset.table, or
// dataset.INFORMATION_SCHEMA.view for a system view.
func qualifySpan(c ColumnResource) string {
	if c.Schema != "" {
		return c.Database + "." + c.Schema + "." + c.Table
	}
	return c.Database + "." + c.Table
}

// spanMetadata is the stored dataset as Bytebase's sync records it
// (bigquery/driver/sync.go): one unnamed schema holding every table and view
// of the dataset as a table, with its columns in order.
func spanMetadata(dataset string, relations []testcontainer.SpanRelation) *metadata.DatabaseSchemaMetadata {
	schema := &metadata.SchemaMetadata{}
	for _, r := range relations {
		table := &metadata.TableMetadata{Name: r.Name}
		for i, c := range r.Columns {
			table.Columns = append(table.Columns, &metadata.ColumnMetadata{Name: c, Position: int32(i + 1), Nullable: true})
		}
		schema.Tables = append(schema.Tables, table)
	}
	return &metadata.DatabaseSchemaMetadata{Name: dataset, Schemas: []*metadata.SchemaMetadata{schema}}
}

// engineSpan is BigQuery's answer. A dry run states no column's sources, so
// every column's lineage is unstated.
func engineSpan(a testcontainer.SpanAnswer) span.Answer {
	if a.Error != "" {
		return span.Answer{Error: a.Error}
	}
	out := span.Answer{Results: a.Columns, Tables: a.Tables}
	for range a.Columns {
		out.Lineage = append(out.Lineage, []string{span.Unstated})
	}
	return out
}

// readSpanCorpus reads queries.sql: one query per line, in the session and
// on the fixture the last "-- session:" and "-- fixture:" lines name. Other
// lines starting "--" are comments.
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
