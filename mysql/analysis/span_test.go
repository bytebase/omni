package analysis

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest/span"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mysql/internal/testcontainer"
	"github.com/bytebase/omni/mysql/parser"
)

// TestQuerySpanFidelity holds ResolveQuerySpan to MySQL: for each query of
// testdata/span/queries.sql, the output column names, each column's lineage
// and the tables read, as omni resolves them against the schema MySQL
// stored for the fixture, and as MySQL states them. testdata/README.md says
// where each of MySQL's answers comes from.
func TestQuerySpanFidelity(t *testing.T) {
	queries, err := readSpanCorpus(filepath.Join("testdata", "span", "queries.sql"))
	if err != nil {
		t.Fatal(err)
	}
	// The schema MySQL stored for each fixture, which omni resolves against.
	stored := map[string]testcontainer.SpanFixture{}
	span.Check(t, span.Suite{
		Path:     filepath.Join("testdata", "span.tsv"),
		Pins:     testcontainer.Pins,
		PinRole:  testcontainer.SpanRole,
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
				return engineSpan(stored[fx.Name], a)
			}
		},
	})
}

// omniSpan is ResolveQuerySpan's answer, resolved as Bytebase resolves it
// for MySQL: the session's database is the default database, against the
// metadata a sync stores, case-sensitively, as on a server with
// lower_case_table_names 0, the pinned server's.
func omniSpan(fx testcontainer.SpanFixture, session, statement string) span.Answer {
	resolved, err := ResolveQuerySpan(statement, ResolveOptions{
		DefaultDatabase: session,
		Catalog:         spanCatalog(fx),
	})
	var pe *parser.ParseError
	switch {
	case errors.As(err, &pe):
		return span.Answer{Error: "1064"}
	case err != nil:
		return span.Answer{Error: "omni"}
	case resolved.NotFoundError != nil:
		// The error MySQL gives for what is missing.
		nf := resolved.NotFoundError
		switch {
		case nf.Column != nil:
			return span.Answer{Error: "1054"}
		case nf.Table != nil:
			return span.Answer{Error: "1146"}
		default:
			return span.Answer{Error: "1049"}
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

// spanCatalog is the stored schema as Bytebase's sync records it: each
// database's tables and views with their columns, and each view's
// definition as information_schema holds it.
func spanCatalog(fx testcontainer.SpanFixture) Catalog {
	databases := map[string]spanDatabase{}
	var names []string
	for _, d := range fx.Databases {
		db := spanDatabase{name: d.Name, tables: map[string]*metadata.TableMetadata{}, views: map[string]*metadata.ViewMetadata{}}
		columns := func(names []string) []*metadata.ColumnMetadata {
			var out []*metadata.ColumnMetadata
			for _, n := range names {
				out = append(out, &metadata.ColumnMetadata{Name: n, Nullable: true})
			}
			return out
		}
		for _, t := range d.Tables {
			db.tableNames = append(db.tableNames, t.Name)
			db.tables[t.Name] = &metadata.TableMetadata{Name: t.Name, Columns: columns(t.Columns)}
		}
		for _, v := range d.Views {
			db.viewNames = append(db.viewNames, v.Name)
			db.views[v.Name] = &metadata.ViewMetadata{Name: v.Name, Definition: v.Definition, Columns: columns(v.Columns)}
		}
		databases[d.Name] = db
		names = append(names, d.Name)
	}
	return Catalog{
		ListDatabaseNames: func() ([]string, error) { return names, nil },
		GetDatabase: func(name string) (Database, error) {
			if db, ok := databases[name]; ok {
				return db, nil
			}
			return nil, nil
		},
	}
}

type spanDatabase struct {
	name                  string
	tableNames, viewNames []string
	tables                map[string]*metadata.TableMetadata
	views                 map[string]*metadata.ViewMetadata
}

func (d spanDatabase) Name() string                                 { return d.name }
func (d spanDatabase) HasSchema() bool                              { return true }
func (d spanDatabase) ListTableNames() []string                     { return d.tableNames }
func (d spanDatabase) GetTable(name string) *metadata.TableMetadata { return d.tables[name] }
func (d spanDatabase) ListViewNames() []string                      { return d.viewNames }
func (d spanDatabase) GetView(name string) *metadata.ViewMetadata   { return d.views[name] }

// engineSpan is MySQL's answer:
//
//   - results: each column definition's name in the prepared statement's
//     result set.
//   - lineage: read from the query MySQL resolved (span_resolved_test.go),
//     where a view's column reads through the view's definition.
//   - tables: the tables the plan accesses. The plan names a view's tables
//     and never the view, and leaves out a table it read while optimizing,
//     so a query that reads a view or whose plan hides a table states its
//     tables with Unstated.
func engineSpan(fx testcontainer.SpanFixture, a testcontainer.SpanAnswer) (span.Answer, error) {
	if a.Error != "" {
		return span.Answer{Error: a.Error}, nil
	}
	var out span.Answer
	for _, c := range a.Columns {
		out.Results = append(out.Results, c.Name)
	}
	lineage, err := resolvedLineage(fx, a.Resolved)
	if err != nil {
		return out, fmt.Errorf("read the resolved query: %w\n    %s", err, a.Resolved)
	}
	if len(lineage) != len(out.Results) {
		return out, fmt.Errorf("the resolved query has %d columns and the result set %d:\n    %s", len(lineage), len(out.Results), a.Resolved)
	}
	out.Lineage = lineage
	out.Tables = a.Tables
	if a.ReadsView || a.Hidden {
		out.Tables = append(out.Tables, span.Unstated)
	}
	return out, nil
}

// readSpanCorpus reads queries.sql: one query per line, in the session (the
// selected database) and on the fixture the last "-- session:" and
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
