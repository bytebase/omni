package analysis_test

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/cosmosdb"
	"github.com/bytebase/omni/cosmosdb/analysis"
	"github.com/bytebase/omni/cosmosdb/ast"
	"github.com/bytebase/omni/cosmosdb/internal/testcontainer"
	"github.com/bytebase/omni/enginetest/span"
	"github.com/bytebase/omni/enginetest/tape"
)

// TestQuerySpanFidelity holds Analyze to the Cosmos DB emulator: for each
// query of testdata/span/queries.sql, the properties a result returns and
// the stored properties each one's value comes from, as omni's analysis
// answers them and as the emulator's results state them. testdata/README.md
// says how each of the emulator's answers is read.
func TestQuerySpanFidelity(t *testing.T) {
	queries, err := readSpanCorpus(filepath.Join("testdata", "span", "queries.sql"))
	if err != nil {
		t.Fatal(err)
	}
	// The document each fixture's container returned for SELECT *, which
	// omni's SELECT * is expanded over.
	stored := map[string]*node{}
	span.Check(t, span.Suite{
		Path:     filepath.Join("testdata", "span.tsv"),
		Pins:     testcontainer.Pins,
		PinRole:  testcontainer.Emulator,
		Fixtures: filepath.Join("testdata", "span", "fixtures"),
		Queries:  queries,
		Omni: func(fx span.Fixture, _, statement string) span.Answer {
			return omniSpan(stored[fx.Name], statement)
		},
		Engine: func(t testing.TB) func(span.Fixture, string, string) (span.Answer, error) {
			raw := tape.OpenRaw(t, testcontainer.Server, testcontainer.Emulator, testcontainer.SpanLive(t))
			ask := func(req testcontainer.SpanRequest) (testcontainer.SpanAnswer, error) {
				text, err := json.Marshal(req)
				if err != nil {
					return testcontainer.SpanAnswer{}, err
				}
				answer, err := raw.Do(string(text))
				if err != nil {
					return testcontainer.SpanAnswer{}, err
				}
				var a testcontainer.SpanAnswer
				if err := json.Unmarshal([]byte(answer), &a); err != nil {
					return a, fmt.Errorf("an unreadable answer %s: %w", answer, err)
				}
				return a, nil
			}
			return func(fx span.Fixture, _, statement string) (span.Answer, error) {
				if stored[fx.Name] == nil {
					if _, err := ask(testcontainer.SpanRequest{Op: "load", Fixture: fx.Name, Documents: fx.SQL}); err != nil {
						return span.Answer{}, fmt.Errorf("load fixture %s: %w", fx.Name, err)
					}
					a, err := ask(testcontainer.SpanRequest{Op: "query", Fixture: fx.Name, Statement: "SELECT * FROM c"})
					if err != nil {
						return span.Answer{}, err
					}
					if a.Error != "" || len(a.Items) != 1 {
						return span.Answer{}, fmt.Errorf("fixture %s: SELECT * answered %+v; the fixture is one document", fx.Name, a)
					}
					doc, err := parseNode(a.Items[0])
					if err != nil {
						return span.Answer{}, err
					}
					stored[fx.Name] = doc
				}
				a, err := ask(testcontainer.SpanRequest{Op: "query", Fixture: fx.Name, Statement: statement})
				if err != nil {
					return span.Answer{}, err
				}
				return engineSpan(stored[fx.Name], a)
			}
		},
	})
}

// omniSpan is Analyze's answer, read as Bytebase reads it: a projection's
// name is the property it masks and its paths the stored properties whose
// policies apply, with the container's own name (the first selector)
// dropped. A SELECT * is expanded over the stored document, as Bytebase's
// masker looks each returned property up in the stored schema.
func omniSpan(doc *node, statement string) span.Answer {
	stmts, err := cosmosdb.Parse(statement)
	if err != nil || len(stmts) != 1 {
		return span.Answer{Error: "omni"}
	}
	sel, ok := stmts[0].AST.(*ast.SelectStmt)
	if !ok {
		return span.Answer{Error: "omni"}
	}
	qa := analysis.Analyze(sel)
	var a span.Answer
	if qa.SelectStar {
		for _, p := range qa.StarPaths {
			for _, key := range keysOf(doc.at(spanPath(p))) {
				if !slices.Contains(a.Results, key) {
					a.Results = append(a.Results, key)
					a.Lineage = append(a.Lineage, []string{renderPath(append(spanPath(p), key))})
				}
			}
		}
		return a
	}
	for _, p := range qa.Projections {
		a.Results = append(a.Results, p.Name)
		var sources []string
		for _, fp := range p.SourcePaths {
			sources = append(sources, renderPath(spanPath(fp)))
		}
		a.Lineage = append(a.Lineage, sources)
	}
	return a
}

// spanPath is a field path as the span compares it: the container's name
// dropped, and with no array index, since a masking policy names a property
// and covers every element of an array.
func spanPath(fp analysis.FieldPath) []string {
	var out []string
	for i, s := range fp {
		if i == 0 || s.IsArray() {
			continue
		}
		out = append(out, s.Name)
	}
	return out
}

func renderPath(names []string) string {
	return strings.Join(append([]string{"$"}, names...), ".")
}

// engineSpan is the emulator's answer. The returned properties are the
// members of the items, in the order the gateway sent them; a query whose
// plan says it selects a VALUE returns one bare, unnamed item. A value's
// sources are the stored values it holds: the fixture's every leaf is a
// value no other leaf holds, so a returned value that is a stored subtree,
// or a string holding a stored string, names where it came from. Those are
// the sources of one run on one document, which cannot say that a value has
// no other source, so every column also holds span.Unstated. Cosmos DB
// reads the one container a request names, whatever the FROM clause calls
// it, so the tables read are unstated too.
func engineSpan(doc *node, a testcontainer.SpanAnswer) (span.Answer, error) {
	if a.Error != "" {
		return span.Answer{Error: a.Error}, nil
	}
	if len(a.Items) == 0 {
		return span.Answer{}, fmt.Errorf("the query returned nothing, so the emulator states nothing about it; the corpus needs a query that returns a row")
	}
	idx := newMarkers(doc)
	out := span.Answer{Tables: []string{span.Unstated}}
	sources := map[string]map[string]bool{}
	add := func(name string, v *node) {
		if sources[name] == nil {
			sources[name] = map[string]bool{span.Unstated: true}
			out.Results = append(out.Results, name)
		}
		idx.find(v, sources[name])
	}
	for _, item := range a.Items {
		v, err := parseNode(item)
		if err != nil {
			return span.Answer{}, err
		}
		switch {
		case a.SelectValue:
			add("", v)
		case v.kind != '{':
			return span.Answer{}, fmt.Errorf("the plan states no VALUE, but an item is %s", item)
		default:
			for i, key := range v.members {
				add(key, v.values[i])
			}
		}
	}
	for _, name := range out.Results {
		var list []string
		for s := range sources[name] {
			list = append(list, s)
		}
		out.Lineage = append(out.Lineage, list)
	}
	return out, nil
}

// markers indexes the stored document: each subtree by its encoding, and
// each string leaf, by the path it is stored at. An encoding stored at more
// than one path (the system properties' stand-in) names none of them.
type markers struct {
	subtrees map[string]string
	strings  map[string]string
}

func newMarkers(doc *node) *markers {
	m := &markers{subtrees: map[string]string{}, strings: map[string]string{}}
	seen := map[string]int{}
	var walk func(n *node, path []string)
	walk = func(n *node, path []string) {
		enc := n.encode()
		seen[enc]++
		m.subtrees[enc] = renderPath(path)
		if n.kind == 's' {
			m.strings[n.scalar] = renderPath(path)
		}
		for i, v := range n.values {
			if n.kind == '{' {
				walk(v, append(slices.Clip(path), n.members[i]))
			} else {
				walk(v, path)
			}
		}
	}
	walk(doc, nil)
	for enc, count := range seen {
		if count > 1 {
			delete(m.subtrees, enc)
			var s string
			if json.Unmarshal([]byte(enc), &s) == nil {
				delete(m.strings, s)
			}
		}
	}
	return m
}

// find adds the stored paths a returned value holds to into: the path of a
// stored subtree it equals; else, for a string, every stored string it holds
// in any case (UPPER keeps the marker's letters); else those of its members.
func (m *markers) find(v *node, into map[string]bool) {
	if path, ok := m.subtrees[v.encode()]; ok {
		into[path] = true
		return
	}
	switch v.kind {
	case 's':
		lower := strings.ToLower(v.scalar)
		for s, path := range m.strings {
			if strings.Contains(lower, strings.ToLower(s)) {
				into[path] = true
			}
		}
	case '{', '[':
		for _, member := range v.values {
			m.find(member, into)
		}
	}
}

// node is a JSON value that keeps its members in order, which the property
// names of a result are read in.
type node struct {
	kind    byte // '{', '[', 's' for a string, or 'v' for any other scalar
	members []string
	values  []*node
	scalar  string // a string's text, or another scalar's JSON
}

func parseNode(b []byte) (*node, error) {
	dec := json.NewDecoder(strings.NewReader(string(b)))
	dec.UseNumber()
	n, err := decodeNode(dec)
	if err != nil {
		return nil, fmt.Errorf("unreadable JSON %s: %w", b, err)
	}
	return n, nil
}

func decodeNode(dec *json.Decoder) (*node, error) {
	tok, err := dec.Token()
	if err != nil {
		return nil, err
	}
	switch t := tok.(type) {
	case json.Delim:
		n := &node{kind: byte(t)}
		for dec.More() {
			if t == '{' {
				key, err := dec.Token()
				if err != nil {
					return nil, err
				}
				n.members = append(n.members, key.(string))
			}
			v, err := decodeNode(dec)
			if err != nil {
				return nil, err
			}
			n.values = append(n.values, v)
		}
		if _, err := dec.Token(); err != nil {
			return nil, err
		}
		return n, nil
	case string:
		return &node{kind: 's', scalar: t}, nil
	default:
		b, err := json.Marshal(t)
		return &node{kind: 'v', scalar: string(b)}, err
	}
}

func (n *node) encode() string {
	var b strings.Builder
	n.write(&b)
	return b.String()
}

func (n *node) write(b *strings.Builder) {
	switch n.kind {
	case 's':
		s, _ := json.Marshal(n.scalar)
		b.Write(s)
	case 'v':
		b.WriteString(n.scalar)
	default:
		b.WriteByte(n.kind)
		for i, v := range n.values {
			if i > 0 {
				b.WriteByte(',')
			}
			if n.kind == '{' {
				k, _ := json.Marshal(n.members[i])
				b.Write(k)
				b.WriteByte(':')
			}
			v.write(b)
		}
		if n.kind == '{' {
			b.WriteByte('}')
		} else {
			b.WriteByte(']')
		}
	}
}

// at is the value stored at path, reading through arrays: a property of an
// array is that property of its elements.
func (n *node) at(path []string) []*node {
	cur := []*node{n}
	for _, name := range path {
		var next []*node
		for _, c := range cur {
			next = append(next, c.member(name)...)
		}
		cur = next
	}
	return cur
}

func (n *node) member(name string) []*node {
	switch n.kind {
	case '{':
		if i := slices.Index(n.members, name); i >= 0 {
			return []*node{n.values[i]}
		}
	case '[':
		var out []*node
		for _, v := range n.values {
			out = append(out, v.member(name)...)
		}
		return out
	}
	return nil
}

// keys is the properties of the objects in ns, or of their elements for an
// array, in the order they are stored.
func keysOf(ns []*node) []string {
	var out []string
	for _, n := range ns {
		if n.kind == '[' {
			out = append(out, keysOf(n.values)...)
			continue
		}
		for _, k := range n.members {
			if !slices.Contains(out, k) {
				out = append(out, k)
			}
		}
	}
	return out
}

// readSpanCorpus reads queries.sql: one query per line, on the fixture and in
// the session the last "-- fixture:" and "-- session:" lines name. Other
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
			out = append(out, span.Query{Fixture: fixture, Session: session, Statement: line})
		}
	}
	return out, sc.Err()
}
