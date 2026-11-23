package analysis

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest/span"
	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/trino/internal/testcontainer"
)

// TestQuerySpanFidelity holds ResolveQuerySpan to Trino: for each query of
// testdata/span/queries.sql, the output column names, each column's lineage
// and the tables read, as omni resolves them against the schema Trino stored
// for the fixture, and as Trino states them. testdata/README.md says where
// each of Trino's answers comes from.
func TestQuerySpanFidelity(t *testing.T) {
	queries, err := readSpanCorpus(filepath.Join("testdata", "span", "queries.sql"))
	if err != nil {
		t.Fatal(err)
	}
	// The schema Trino stored for each fixture, which omni resolves against.
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
			o := testcontainer.OpenSpans(t)
			built := ""
			return func(fx span.Fixture, session, statement string) (span.Answer, error) {
				if built != fx.Name {
					// The memory catalog holds one fixture at a time, so a
					// fixture is built again when the corpus returns to it.
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
				return engineSpan(stored[fx.Name], a)
			}
		},
	})
}

// omniSpan is ResolveQuerySpan's answer, resolved as Bytebase resolves it
// for Trino: the session's catalog is the database and its schema the
// default schema, against the metadata a sync stores, case-sensitively
// (Bytebase's store.IsObjectCaseSensitive holds for Trino).
func omniSpan(fx testcontainer.SpanFixture, session, statement string) span.Answer {
	catalog, schema, _ := strings.Cut(session, ".")
	meta := spanMetadata(fx)
	resolved, err := ResolveQuerySpan(statement, ResolveOptions{
		DefaultDatabase: catalog,
		DefaultSchema:   schema,
		ListDatabases:   func() ([]string, error) { return []string{meta.Name}, nil },
		GetDatabaseMetadata: func(database string) (*metadatapb.DatabaseSchemaMetadata, error) {
			if database == meta.Name {
				return meta, nil
			}
			return nil, nil
		},
	})
	if err != nil {
		return span.Answer{Error: "omni"}
	}
	var a span.Answer
	tables := map[string]bool{}
	for c := range resolved.SourceColumns {
		tables[Column{Database: c.Database, Schema: c.Schema, Table: c.Table}.String()] = true
	}
	for name := range tables {
		a.Tables = append(a.Tables, name)
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
	meta := &metadatapb.DatabaseSchemaMetadata{Name: "memory"}
	for _, s := range fx.Schemas {
		schema := &metadatapb.SchemaMetadata{Name: s.Name}
		for _, t := range s.Tables {
			schema.Tables = append(schema.Tables, &metadatapb.TableMetadata{Name: t.Name, Columns: columns(t.Columns)})
		}
		for _, v := range s.Views {
			schema.Views = append(schema.Views, &metadatapb.ViewMetadata{Name: v.Name, Definition: v.Definition, Columns: columns(v.Columns)})
		}
		meta.Schemas = append(meta.Schemas, schema)
	}
	return meta
}

// engineSpan is Trino's answer. Trino's lineage stops at a view's column;
// this follows it through the view's definition, whose lineage the fixture
// recorded, so a column read through a view names every column on the way
// to the tables, as omni's lineage does.
func engineSpan(fx testcontainer.SpanFixture, a testcontainer.SpanAnswer) (span.Answer, error) {
	if a.Error != "" {
		return span.Answer{Error: a.Error}, nil
	}
	if a.LineageError != "" {
		return span.Answer{}, fmt.Errorf("trino refused the CREATE TABLE AS that reads the lineage with %s; the corpus needs a query it takes", a.LineageError)
	}
	viewColumns := map[string][]string{}
	for _, s := range fx.Schemas {
		for _, v := range s.Views {
			for _, c := range v.Columns {
				viewColumns["memory."+s.Name+"."+v.Name] = append(viewColumns["memory."+s.Name+"."+v.Name], c.Name)
			}
		}
	}
	var through func(col string, into map[string]bool)
	through = func(col string, into map[string]bool) {
		if into[col] {
			return
		}
		into[col] = true
		i := strings.LastIndexByte(col, '.')
		view, name := col[:i], col[i+1:]
		for k, c := range viewColumns[view] {
			if c == name && k < len(fx.Views[view]) {
				for _, src := range fx.Views[view][k] {
					through(strings.ToLower(src), into)
				}
			}
		}
	}
	// Trino names relations and columns case-insensitively and stores them in
	// lower case, but its lineage spells a view's column as the view's
	// definition wrote it (VipName): fold it to the stored name.
	out := span.Answer{Results: a.Columns}
	for _, t := range a.Tables {
		out.Tables = append(out.Tables, strings.ToLower(t.Name))
	}
	for _, sources := range a.Lineage {
		set := map[string]bool{}
		for _, src := range sources {
			through(strings.ToLower(src), set)
		}
		var list []string
		for c := range set {
			list = append(list, c)
		}
		out.Lineage = append(out.Lineage, list)
	}
	return out, nil
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
