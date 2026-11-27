package analysis

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/clickhouse/internal/testcontainer"
	"github.com/bytebase/omni/clickhouse/parser"
	"github.com/bytebase/omni/enginetest/span"
	"github.com/bytebase/omni/metadata"
)

// TestQuerySpanFidelity holds ResolveQuerySpan to ClickHouse: for each query
// of testdata/span/queries.sql, the output column names, each column's
// lineage and the tables read, as omni resolves them against the schema
// ClickHouse stored for the fixture, and as ClickHouse states them.
// testdata/README.md says where each of ClickHouse's answers comes from.
func TestQuerySpanFidelity(t *testing.T) {
	queries, err := readSpanCorpus(filepath.Join("testdata", "span", "queries.sql"))
	if err != nil {
		t.Fatal(err)
	}
	// The schema ClickHouse stored for each fixture, which omni resolves
	// against.
	stored := map[string]testcontainer.SpanFixture{}
	span.Check(t, span.Suite{
		Path:     filepath.Join("testdata", "span.tsv"),
		Pins:     testcontainer.Pins,
		Fixtures: filepath.Join("testdata", "span", "fixtures"),
		Queries:  queries,
		Omni: func(fx span.Fixture, session, statement string) span.Answer {
			return omniSpan(stored[fx.Name], session, statement)
		},
		Engine: func(t testing.TB) func(span.Fixture, string, string) (span.Answer, error) {
			s := testcontainer.OpenSpans(t)
			built := ""
			return func(fx span.Fixture, session, statement string) (span.Answer, error) {
				if built != fx.Name {
					// The server holds one fixture at a time, so a fixture
					// is built again when the corpus returns to it.
					f, err := s.Fixture(fx.SQL)
					if err != nil {
						return span.Answer{}, fmt.Errorf("fixture %s: %w", fx.Name, err)
					}
					stored[fx.Name], built = f, fx.Name
				}
				a, err := s.Span(session, statement)
				if err != nil {
					return span.Answer{}, err
				}
				return engineSpan(a), nil
			}
		},
	})
}

// omniSpan is ResolveQuerySpan's answer, with the session's database as the
// default database, against the metadata a sync stores for the fixture.
func omniSpan(fx testcontainer.SpanFixture, session, statement string) span.Answer {
	resolved, err := ResolveQuerySpan(statement, ResolveOptions{
		DefaultDatabase: session,
		Catalog:         spanCatalog(fx),
	})
	var pe *parser.ParseError
	switch {
	case errors.As(err, &pe):
		return span.Answer{Error: "SYNTAX_ERROR"}
	case errors.Is(err, ErrAmbiguousColumn):
		return span.Answer{Error: "AMBIGUOUS_IDENTIFIER"}
	case err != nil:
		return span.Answer{Error: "omni: " + err.Error()}
	case resolved.NotFoundError != nil:
		// The error ClickHouse names for what is missing.
		nf := resolved.NotFoundError
		switch {
		case nf.Column != nil:
			return span.Answer{Error: "UNKNOWN_IDENTIFIER"}
		case nf.Table != nil:
			return span.Answer{Error: "UNKNOWN_TABLE"}
		default:
			return span.Answer{Error: "UNKNOWN_DATABASE"}
		}
	}
	var a span.Answer
	for c := range resolved.SourceColumns {
		a.Tables = append(a.Tables, c.Database+"."+c.Table)
	}
	for _, r := range resolved.Results {
		a.Results = append(a.Results, r.Name)
		var sources []string
		for c := range r.SourceColumns {
			sources = append(sources, c.Database+"."+c.Table+"."+c.Column)
		}
		a.Lineage = append(a.Lineage, sources)
	}
	return a
}

// spanCatalog is the stored schema as a sync records it: each database's
// tables and views with their columns, and each view's create_table_query
// as its definition.
func spanCatalog(fx testcontainer.SpanFixture) Catalog {
	databases := map[string]spanDatabase{}
	for _, d := range fx.Databases {
		db := spanDatabase{tables: map[string]*metadata.TableMetadata{}, views: map[string]*metadata.ViewMetadata{}}
		columns := func(names []string) []*metadata.ColumnMetadata {
			var out []*metadata.ColumnMetadata
			for _, n := range names {
				out = append(out, &metadata.ColumnMetadata{Name: n})
			}
			return out
		}
		for _, t := range d.Tables {
			table := &metadata.TableMetadata{Name: t.Name, Columns: columns(t.Columns)}
			for _, c := range table.Columns {
				if computed, ok := t.Computed[c.Name]; ok {
					c.Generation = generation(computed)
				}
			}
			db.tables[t.Name] = table
		}
		for _, v := range d.Views {
			db.views[v.Name] = &metadata.ViewMetadata{Name: v.Name, Definition: v.Definition, Columns: columns(v.Columns)}
		}
		databases[d.Name] = db
	}
	return Catalog{
		GetDatabase: func(name string) (Database, error) {
			if db, ok := databases[name]; ok {
				return db, nil
			}
			return nil, nil
		},
	}
}

// generation is an ALIAS or MATERIALIZED column as a sync records it.
func generation(c testcontainer.SpanComputed) *metadata.GenerationMetadata {
	g := &metadata.GenerationMetadata{Type: metadata.GenerationMetadata_TYPE_STORED, Expression: c.Expression}
	if c.Kind == "ALIAS" {
		g.Type = metadata.GenerationMetadata_TYPE_VIRTUAL
	}
	return g
}

type spanDatabase struct {
	tables map[string]*metadata.TableMetadata
	views  map[string]*metadata.ViewMetadata
}

func (d spanDatabase) GetTable(name string) *metadata.TableMetadata { return d.tables[name] }
func (d spanDatabase) GetView(name string) *metadata.ViewMetadata   { return d.views[name] }

// TableNames lists the tables and views, as a sync's metadata does.
func (d spanDatabase) TableNames() []string {
	var out []string
	for name := range d.tables {
		out = append(out, name)
	}
	for name := range d.views {
		out = append(out, name)
	}
	slices.Sort(out)
	return out
}

// engineSpan is ClickHouse's answer: the result set's column names, the
// analyzer's lineage of each, and query_log's tables.
func engineSpan(a testcontainer.SpanAnswer) span.Answer {
	if a.Error != "" {
		return span.Answer{Error: a.Error}
	}
	return span.Answer{Results: a.Columns, Lineage: a.Lineage, Tables: a.Tables}
}

// readSpanCorpus reads queries.sql: one query per line, in the session (the
// current database) and on the fixture the last "-- session:" and
// "-- fixture:" lines name. Other lines starting "--" are comments.
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
