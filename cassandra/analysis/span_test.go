package analysis

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/cassandra"
	"github.com/bytebase/omni/cassandra/internal/testcontainer"
	"github.com/bytebase/omni/enginetest/span"
	"github.com/bytebase/omni/enginetest/tape"
	"github.com/bytebase/omni/metadata"
)

// TestMain accounts for the package's tapes and removes the Cassandra server a
// recording started.
func TestMain(m *testing.M) {
	code := tape.Main(m)
	testcontainer.Terminate()
	os.Exit(code)
}

// TestQuerySpanFidelity holds ResolveQuerySpan to Cassandra: for each query of
// testdata/span/queries.sql, the output column names, each column's lineage
// and the relations read, as omni resolves them against the schema Cassandra
// stored for the fixture, and as Cassandra's result metadata states them.
// testdata/README.md says where each of Cassandra's answers comes from.
func TestQuerySpanFidelity(t *testing.T) {
	queries, err := readSpanCorpus(filepath.Join("testdata", "span", "queries.sql"))
	if err != nil {
		t.Fatal(err)
	}
	// The schema Cassandra stored for each fixture, which omni resolves
	// against.
	stored := map[string]testcontainer.SpanFixture{}
	span.Check(t, span.Suite{
		Path:     filepath.Join("testdata", "span.tsv"),
		Pins:     testcontainer.Pins,
		PinRole:  testcontainer.Driver,
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
					// The server holds one fixture at a time, so a fixture is
					// built again when the corpus returns to it.
					var statements []string
					for _, seg := range cassandra.Split(fx.SQL) {
						if !seg.Empty {
							statements = append(statements, strings.TrimSuffix(strings.TrimSpace(seg.Text), ";"))
						}
					}
					s, err := o.Fixture(statements)
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
// the session keyspace, against the metadata cassandra/driver's sync stores
// for each keyspace.
func omniSpan(fx testcontainer.SpanFixture, session, statement string) span.Answer {
	resolved, err := ResolveQuerySpan(statement, ResolveOptions{
		DefaultKeyspace: session,
		GetDatabaseMetadata: func(keyspace string) (*metadata.DatabaseSchemaMetadata, error) {
			return spanMetadata(fx, keyspace), nil
		},
	})
	if err != nil {
		return span.Answer{Error: "omni"}
	}
	var a span.Answer
	for c := range resolved.SourceColumns {
		a.Tables = append(a.Tables, c.String())
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

// spanMetadata is the keyspace's schema as cassandra/driver's SyncDatabase
// records it: one unnamed schema holding each table of system_schema.tables
// and each materialized view of system_schema.views, with its columns in
// SELECT * order, and a view's columns depending on the base-table columns of
// the same name. A keyspace the fixture does not hold has no metadata.
func spanMetadata(fx testcontainer.SpanFixture, keyspace string) *metadata.DatabaseSchemaMetadata {
	for _, ks := range fx.Keyspaces {
		if ks.Name != keyspace {
			continue
		}
		schema := &metadata.SchemaMetadata{}
		for _, t := range ks.Tables {
			schema.Tables = append(schema.Tables, &metadata.TableMetadata{Name: t.Name, Columns: starColumns(t.Columns)})
		}
		for _, v := range ks.Views {
			view := &metadata.MaterializedViewMetadata{Name: v.Name, Columns: starColumns(v.Columns)}
			for _, c := range view.Columns {
				view.DependencyColumns = append(view.DependencyColumns, &metadata.DependencyColumn{Table: v.Base, Column: c.GetName()})
			}
			schema.MaterializedViews = append(schema.MaterializedViews, view)
		}
		return &metadata.DatabaseSchemaMetadata{Name: keyspace, Schemas: []*metadata.SchemaMetadata{schema}}
	}
	return nil
}

// starColumns orders a relation's columns as SELECT * returns them: the
// partition key and the clustering columns by position, then the static
// columns and the regular ones, each by name.
func starColumns(cols []testcontainer.SpanColumn) []*metadata.ColumnMetadata {
	rank := map[string]int{"partition_key": 0, "clustering": 1, "static": 2, "regular": 3}
	sorted := slices.Clone(cols)
	slices.SortStableFunc(sorted, func(a, b testcontainer.SpanColumn) int {
		if ra, rb := rank[a.Kind], rank[b.Kind]; ra != rb {
			return ra - rb
		}
		if rank[a.Kind] < 2 {
			return a.Position - b.Position
		}
		return strings.Compare(a.Name, b.Name)
	})
	var out []*metadata.ColumnMetadata
	for _, c := range sorted {
		out = append(out, &metadata.ColumnMetadata{Name: c.Name, Type: c.Type, Nullable: true})
	}
	return out
}

// engineSpan is Cassandra's answer, read from the result metadata: each
// output column's name, and the keyspace and table Cassandra names for it.
//
// The metadata names a column's keyspace and table for every output column,
// an expression's too, and gives it a name: the column's for a plain column
// reference, the alias where there is one, and the selector's text for any
// other. So a column whose name is a column of the relation it names is that
// column; any other is span.Unstated, since the metadata does not say which
// columns an expression reads. A materialized view is named as itself, and
// the metadata says nothing of the base table its rows come from, so a
// column or a relation that is a materialized view is followed by
// span.Unstated.
func engineSpan(fx testcontainer.SpanFixture, a testcontainer.SpanAnswer) span.Answer {
	if a.Error != "" {
		return span.Answer{Error: a.Error}
	}
	type relation struct {
		columns map[string]bool
		view    bool
	}
	relations := map[string]relation{}
	for _, ks := range fx.Keyspaces {
		for _, list := range []struct {
			rels []testcontainer.SpanRelation
			view bool
		}{{ks.Tables, false}, {ks.Views, true}} {
			for _, r := range list.rels {
				cols := map[string]bool{}
				for _, c := range r.Columns {
					cols[c.Name] = true
				}
				relations[ks.Name+"."+r.Name] = relation{columns: cols, view: list.view}
			}
		}
	}
	var out span.Answer
	tables := map[string]bool{}
	for _, c := range a.Columns {
		name := c.Keyspace + "." + c.Table
		rel := relations[name]
		out.Results = append(out.Results, c.Name)
		var sources []string
		if rel.columns[c.Name] {
			sources = append(sources, name+"."+c.Name)
		}
		if !rel.columns[c.Name] || rel.view {
			sources = append(sources, span.Unstated)
		}
		out.Lineage = append(out.Lineage, sources)
		if !tables[name] {
			tables[name] = true
			out.Tables = append(out.Tables, name)
			if rel.view {
				out.Tables = append(out.Tables, span.Unstated)
			}
		}
	}
	return out
}

// readSpanCorpus reads queries.sql: one query per line, on the fixture and in
// the keyspace the last "-- fixture:" and "-- session:" lines name. Other
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
