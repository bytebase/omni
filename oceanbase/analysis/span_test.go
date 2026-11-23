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
	"github.com/bytebase/omni/enginetest/tape"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oceanbase/internal/testcontainer"
	"github.com/bytebase/omni/oceanbase/parser"
)

// TestMain runs the tests through the tape harness, and removes the server a
// recording started.
func TestMain(m *testing.M) {
	code := tape.Main(m)
	testcontainer.OceanBase.Terminate()
	os.Exit(code)
}

// TestQuerySpanFidelity holds ResolveQuerySpan to OceanBase: for each query of
// testdata/span/queries.sql, the output column names, each column's lineage
// and the tables read, as omni resolves them against the schema OceanBase
// stored for the fixture, and as OceanBase states them. testdata/README.md says
// where each of OceanBase's answers comes from.
func TestQuerySpanFidelity(t *testing.T) {
	queries, err := readSpanCorpus(filepath.Join("testdata", "span", "queries.sql"))
	if err != nil {
		t.Fatal(err)
	}
	// The schema OceanBase stored for each fixture, which omni resolves against.
	stored := map[string]testcontainer.SpanFixture{}
	span.Check(t, span.Suite{
		Path:     filepath.Join("testdata", "span.tsv"),
		Pins:     testcontainer.Pins,
		PinRole:  "default",
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
// for OceanBase: the session's database is the default database, against the
// metadata a sync stores, ignoring case, as on a server with
// lower_case_table_names 1, the pinned server's.
func omniSpan(fx testcontainer.SpanFixture, session, statement string) span.Answer {
	resolved, err := ResolveQuerySpan(statement, ResolveOptions{
		DefaultDatabase:     session,
		IgnoreCaseSensitive: true,
		Catalog:             spanCatalog(fx),
	})
	var pe *parser.ParseError
	switch {
	case errors.As(err, &pe):
		return span.Answer{Error: "1064"}
	case err != nil:
		return span.Answer{Error: "omni"}
	case resolved.NotFoundError != nil:
		// The error OceanBase gives for what is missing.
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

// engineSpan is OceanBase's answer:
//
//   - results: each column definition's name in the prepared statement's
//     result set.
//   - lineage: the column a definition's org_table and org_name name, read
//     through a view's stored definition (span_resolved_test.go). OceanBase
//     names the relation the query reads, so a column of a derived table or
//     CTE, and an expression, which has no origin, leave it unstated.
//   - tables: the tables and views OceanBase checks the SELECT privilege on.
//     A view is checked as the view, never its tables, so a query that reads
//     a view states its tables with Unstated.
func engineSpan(fx testcontainer.SpanFixture, a testcontainer.SpanAnswer) (span.Answer, error) {
	if a.Error != "" {
		return span.Answer{Error: a.Error}, nil
	}
	var out span.Answer
	r := newResolvedReader(fx)
	names := newFixtureNames(fx)
	for _, c := range a.Columns {
		out.Results = append(out.Results, c.Name)
		key := strings.ToLower(c.Schema + "." + c.OrgTable)
		var l lineage
		switch {
		case c.Schema == "" || c.OrgTable == "":
			l = unknown()
		case names.views[key] != "":
			l = r.tableColumn(names.views[key], c.OrgName)
		case names.tables[key] != "":
			l = unknown()
			if col, ok := names.columns[key+"."+strings.ToLower(c.OrgName)]; ok {
				l = lineage{sources: map[string]bool{col: true}}
			}
		default:
			l = unknown()
		}
		out.Lineage = append(out.Lineage, l.list())
	}
	for _, t := range a.Tables {
		key := strings.ToLower(t)
		switch {
		case names.views[key] != "":
			out.Tables = append(out.Tables, names.views[key], span.Unstated)
		case names.tables[key] != "":
			out.Tables = append(out.Tables, names.tables[key])
		default:
			return out, fmt.Errorf("OceanBase checks %s, which the fixture does not have", t)
		}
	}
	return out, nil
}

// fixtureNames maps each of a fixture's tables and views, and the tables'
// columns, by lower-case name to the name it stored.
type fixtureNames struct {
	tables  map[string]string // db.table
	views   map[string]string // db.view
	columns map[string]string // db.table.column
}

func newFixtureNames(fx testcontainer.SpanFixture) fixtureNames {
	n := fixtureNames{tables: map[string]string{}, views: map[string]string{}, columns: map[string]string{}}
	for _, d := range fx.Databases {
		for _, t := range d.Tables {
			key := strings.ToLower(d.Name + "." + t.Name)
			n.tables[key] = d.Name + "." + t.Name
			for _, c := range t.Columns {
				n.columns[key+"."+strings.ToLower(c)] = d.Name + "." + t.Name + "." + c
			}
		}
		for _, v := range d.Views {
			n.views[strings.ToLower(d.Name+"."+v.Name)] = d.Name + "." + v.Name
		}
	}
	return n
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
