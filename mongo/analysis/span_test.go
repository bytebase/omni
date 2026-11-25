package analysis

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest/span"
	"github.com/bytebase/omni/mongo"
	"github.com/bytebase/omni/mongo/internal/testcontainer"
)

// TestQuerySpanFidelity holds the collections a statement reads to the pinned
// MongoDB: for each statement of testdata/span/queries.js, Analyze's Reads,
// with views expanded against the views the server stored for the fixture,
// and the collections the server's plan reads. testdata/README.md says where
// each answer comes from.
func TestQuerySpanFidelity(t *testing.T) {
	queries, err := readSpanCorpus(filepath.Join("testdata", "span", "queries.js"))
	if err != nil {
		t.Fatal(err)
	}
	// The views the server stored for each fixture, which omni expands.
	stored := map[string][]testcontainer.SpanView{}
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
					// The harness holds one fixture at a time, so a fixture
					// is built again when the corpus returns to it.
					a, err := s.SpanFixture(fx.SQL)
					if err != nil {
						return span.Answer{}, fmt.Errorf("fixture %s: %w", fx.Name, err)
					}
					stored[fx.Name], built = a.Views, fx.Name
				}
				a, err := s.Span(session, statement)
				if err != nil {
					return span.Answer{}, err
				}
				if a.Error != "" {
					return span.Answer{Error: a.Error}, nil
				}
				return unstatedFields(span.Answer{Tables: a.Tables}), nil
			}
		},
	})
}

// omniSpan is the collections Analyze says statement reads in database
// session, with each view followed to what it reads, as Bytebase would
// expand them from the views its sync stores.
func omniSpan(views []testcontainer.SpanView, session, statement string) span.Answer {
	statements, err := mongo.Parse(statement)
	if err != nil || len(statements) != 1 {
		return span.Answer{Error: "omni"}
	}
	a := Analyze(statements[0].AST)
	if a == nil {
		return span.Answer{Error: "omni"}
	}
	qualify := func(ns Namespace) Namespace {
		if ns.Database == "" {
			ns.Database = session
		}
		return ns
	}
	definitions := map[Namespace]string{}
	for _, v := range views {
		definitions[Namespace{Database: v.Database, Collection: v.Name}] = v.Definition
	}
	reads, err := ExpandViews(a.Reads, func(ns Namespace) (string, bool) {
		d, ok := definitions[qualify(ns)]
		return d, ok
	})
	if err != nil {
		return span.Answer{Error: "omni"}
	}
	var out span.Answer
	for _, ns := range reads {
		ns = qualify(ns)
		out.Tables = append(out.Tables, ns.Database+"."+ns.Collection)
	}
	return unstatedFields(out)
}

// unstatedFields marks the output fields unstated. A plan does not name the
// fields a query returns, and Analyze does not either: Bytebase masks each
// returned document against its collection's schema.
func unstatedFields(a span.Answer) span.Answer {
	a.Results = []string{span.Unstated}
	a.Lineage = [][]string{{span.Unstated}}
	return a
}

// readSpanCorpus reads queries.js: one statement per line, in the database
// and on the fixture the last "// session:" and "// fixture:" lines name.
// Other lines starting "//" are comments.
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
		case strings.HasPrefix(line, "// fixture:"):
			fixture = strings.TrimSpace(strings.TrimPrefix(line, "// fixture:"))
		case strings.HasPrefix(line, "// session:"):
			session = strings.TrimSpace(strings.TrimPrefix(line, "// session:"))
		case strings.HasPrefix(line, "//"):
		case fixture == "" || session == "":
			return nil, fmt.Errorf("%s:%d: a statement before the // fixture: and // session: lines", path, n)
		default:
			out = append(out, span.Query{Fixture: fixture, Session: session, Statement: line})
		}
	}
	return out, sc.Err()
}
