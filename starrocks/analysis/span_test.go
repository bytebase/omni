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
	"github.com/bytebase/omni/starrocks/internal/testcontainer"
	"github.com/bytebase/omni/starrocks/parser"
)

// TestMain removes the server a recording started.
func TestMain(m *testing.M) {
	code := tape.Main(m)
	testcontainer.Terminate()
	os.Exit(code)
}

// TestQuerySpanFidelity holds ResolveQuerySpan to StarRocks: for each query of
// testdata/span/queries.sql, the output column names, each column's lineage
// and the tables read, as omni resolves them against the schema StarRocks
// stored for the fixture, and as StarRocks states them. testdata/README.md says
// where each of StarRocks's answers comes from.
func TestQuerySpanFidelity(t *testing.T) {
	queries, err := readSpanCorpus(filepath.Join("testdata", "span", "queries.sql"))
	if err != nil {
		t.Fatal(err)
	}
	// The schema StarRocks stored for each fixture, which omni resolves against.
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
				return engineSpan(stored[fx.Name], a)
			}
		},
	})
}

// omniSpan is ResolveQuerySpan's answer: the session's database is the
// default database, against the metadata a sync stores for the fixture,
// case-sensitively, as StarRocks matches table names.
func omniSpan(fx testcontainer.SpanFixture, session, statement string) span.Answer {
	resolved, err := ResolveQuerySpan(statement, ResolveOptions{
		DefaultDatabase: session,
		Catalog:         spanCatalog(fx),
	})
	var pe *parser.ParseError
	switch {
	case errors.As(err, &pe):
		return span.Answer{Error: errSyntax}
	case err != nil:
		return span.Answer{Error: "omni"}
	case resolved.NotFoundError != nil:
		// The cause StarRocks names for what is missing.
		nf := resolved.NotFoundError
		switch {
		case nf.Column != nil:
			return span.Answer{Error: errColumn}
		case nf.Table != nil:
			return span.Answer{Error: errTable}
		default:
			return span.Answer{Error: errDatabase}
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
	for _, d := range fx.Databases {
		db := spanDatabase{
			name:              d.Name,
			tables:            map[string]*metadata.TableMetadata{},
			views:             map[string]*metadata.ViewMetadata{},
			materializedViews: map[string]*metadata.MaterializedViewMetadata{},
		}
		columns := func(names []string) []*metadata.ColumnMetadata {
			var out []*metadata.ColumnMetadata
			for _, n := range names {
				out = append(out, &metadata.ColumnMetadata{Name: n, Nullable: true})
			}
			return out
		}
		for _, t := range d.Tables {
			db.tables[t.Name] = &metadata.TableMetadata{Name: t.Name, Columns: columns(t.Columns)}
		}
		for _, v := range d.Views {
			db.views[v.Name] = &metadata.ViewMetadata{Name: v.Name, Definition: v.Definition, Columns: columns(v.Columns)}
		}
		for _, v := range d.MaterializedViews {
			db.materializedViews[v.Name] = &metadata.MaterializedViewMetadata{Name: v.Name, Definition: v.Definition}
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

type spanDatabase struct {
	name              string
	tables            map[string]*metadata.TableMetadata
	views             map[string]*metadata.ViewMetadata
	materializedViews map[string]*metadata.MaterializedViewMetadata
}

func (d spanDatabase) Name() string                                 { return d.name }
func (d spanDatabase) GetTable(name string) *metadata.TableMetadata { return d.tables[name] }
func (d spanDatabase) GetView(name string) *metadata.ViewMetadata   { return d.views[name] }
func (d spanDatabase) GetMaterializedView(name string) *metadata.MaterializedViewMetadata {
	return d.materializedViews[name]
}

// The classes of StarRocks's refusals, named by the cause its message
// states: its error numbers do not tell an unknown column from a syntax
// error.
const (
	errSyntax    = "syntax"
	errColumn    = "unknown-column"
	errTable     = "unknown-table"
	errDatabase  = "unknown-database"
	errAmbiguous = "ambiguous-column"
)

// refusal names the cause of a StarRocks error, or keeps the error.
func refusal(err string) string {
	switch {
	case strings.Contains(err, "syntax error"):
		return errSyntax
	case strings.Contains(err, "cannot be resolved"):
		return errColumn
	case strings.Contains(err, "Unknown table"):
		return errTable
	case strings.Contains(err, "Unknown database"):
		return errDatabase
	case strings.Contains(err, "is ambiguous"):
		return errAmbiguous
	}
	return err
}

// engineSpan is StarRocks's answer:
//
//   - results: each column definition's name in the query's result set.
//   - lineage: each output slot of the logical plan (EXPLAIN LOGICAL),
//     traced to the scanned columns it reads (span_plan_test.go). A slot
//     that is a scanned column, or a rename of one, states it exactly; an
//     expression states the columns it reads, plus Unstated. The plan
//     inlines views, so it states no view's column: a query that reads a
//     view states every column's sources plus Unstated.
//   - tables: the tables and views the query dump says the analyzer bound,
//     before the optimizer dropped any.
func engineSpan(fx testcontainer.SpanFixture, a testcontainer.SpanAnswer) (span.Answer, error) {
	if a.Error != "" {
		return span.Answer{Error: refusal(a.Error)}, nil
	}
	out := span.Answer{Results: a.Columns}
	lineage, err := planLineage(a.Plan, a.Tables)
	if err != nil {
		return out, fmt.Errorf("read the plan: %w\n    %s", err, strings.Join(a.Plan, "\n    "))
	}
	if len(lineage) != len(a.Columns) {
		return out, fmt.Errorf("the plan has %d columns and the result set %d:\n    %s", len(lineage), len(a.Columns), strings.Join(a.Plan, "\n    "))
	}
	if len(a.Views) > 0 {
		for i := range lineage {
			lineage[i] = append(lineage[i], span.Unstated)
		}
	}
	out.Lineage = lineage
	out.Tables = append(append([]string{}, a.Tables...), a.Views...)
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
