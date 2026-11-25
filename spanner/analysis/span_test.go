package analysis

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest/span"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/spanner/internal/testcontainer"
)

// spanDatabase is the database omni resolves the span suite's queries in.
// The emulator's database has another name in every recording, and neither
// answer names it.
const spanDatabase = "spandb"

// TestQuerySpanFidelity holds ResolveQuerySpan to the pinned Spanner
// emulator: for each query of testdata/span/queries.sql, the output column
// names and the tables read, as omni resolves them against the schema the
// emulator stored for the fixture, and as the emulator states them.
// testdata/README.md says where each of its answers comes from.
func TestQuerySpanFidelity(t *testing.T) {
	queries, err := readSpanCorpus(filepath.Join("testdata", "span", "queries.sql"))
	if err != nil {
		t.Fatal(err)
	}
	// The schema the emulator stored for each fixture, which omni resolves
	// against.
	stored := map[string][]testcontainer.SpanSchema{}
	span.Check(t, span.Suite{
		Path:     filepath.Join("testdata", "span.tsv"),
		Pins:     testcontainer.Pins,
		Fixtures: filepath.Join("testdata", "span", "fixtures"),
		Queries:  queries,
		Omni: func(fx span.Fixture, _, statement string) span.Answer {
			return omniSpan(stored[fx.Name], statement)
		},
		Engine: func(t testing.TB) func(span.Fixture, string, string) (span.Answer, error) {
			o := testcontainer.Open(t)
			built := ""
			return func(fx span.Fixture, _, statement string) (span.Answer, error) {
				if built != fx.Name {
					// The harness holds one fixture at a time, so a fixture
					// is built again when the corpus returns to it.
					a, err := o.SpanFixture(fx.SQL)
					if err != nil {
						return span.Answer{}, fmt.Errorf("fixture %s: %w", fx.Name, err)
					}
					stored[fx.Name], built = a.Schemas, fx.Name
				}
				a, err := o.Span(statement)
				if err != nil {
					return span.Answer{}, err
				}
				return engineSpan(a), nil
			}
		},
	})
}

// omniSpan is ResolveQuerySpan's answer, resolved as Bytebase resolves it
// for Spanner, against the metadata its sync stores.
func omniSpan(schemas []testcontainer.SpanSchema, statement string) span.Answer {
	meta := spanMetadata(schemas)
	opts := spannerResolveOptions()
	opts.DefaultDatabase = spanDatabase
	opts.GetDatabaseMetadata = func(database string) (*metadata.DatabaseSchemaMetadata, error) {
		if database == spanDatabase {
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
		tables[qualifySpan(c.Schema, c.Table)] = true
	}
	for name := range tables {
		a.Tables = append(a.Tables, name)
	}
	for _, r := range resolved.Results {
		a.Results = append(a.Results, r.Name)
		var sources []string
		for c := range r.SourceColumns {
			sources = append(sources, qualifySpan(c.Schema, c.Table)+"."+c.Column)
		}
		a.Lineage = append(a.Lineage, sources)
	}
	return a
}

func qualifySpan(schema, name string) string {
	if schema == "" {
		return name
	}
	return schema + "." + name
}

// spanMetadata is the stored schema as Bytebase's sync records it
// (spanner/driver/sync.go): every schema, the default one named "", with its
// tables, and its views with their definitions and columns.
func spanMetadata(schemas []testcontainer.SpanSchema) *metadata.DatabaseSchemaMetadata {
	columns := func(cols []testcontainer.SpanColumn) []*metadata.ColumnMetadata {
		var out []*metadata.ColumnMetadata
		for _, c := range cols {
			out = append(out, &metadata.ColumnMetadata{Name: c.Name, Type: c.Type, Nullable: true})
		}
		return out
	}
	meta := &metadata.DatabaseSchemaMetadata{Name: spanDatabase}
	for _, s := range schemas {
		schema := &metadata.SchemaMetadata{Name: s.Name}
		for _, t := range s.Tables {
			schema.Tables = append(schema.Tables, &metadata.TableMetadata{Name: t.Name, Columns: columns(t.Columns)})
		}
		for _, v := range s.Views {
			schema.Views = append(schema.Views, &metadata.ViewMetadata{Name: v.Name, Definition: v.Definition, Columns: columns(v.Columns)})
		}
		meta.Schemas = append(meta.Schemas, schema)
	}
	return meta
}

// engineSpan is the emulator's answer. It states no column's sources, so
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
