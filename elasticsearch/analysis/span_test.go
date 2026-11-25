package analysis_test

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/elasticsearch"
	"github.com/bytebase/omni/elasticsearch/analysis"
	"github.com/bytebase/omni/elasticsearch/internal/testcontainer"
	"github.com/bytebase/omni/enginetest/span"
)

// session is the one session every row runs in: Elasticsearch has no
// session state a request's names resolve in.
const session = "elastic"

// TestQuerySpanFidelity holds the query span Bytebase reads for masking to
// Elasticsearch: for each request of testdata/span/queries.http, the
// document fields it returns, the index each comes from, and the indices it
// reads, as omni's analysis answers them against the indices a sync stores,
// and as Elasticsearch states them. testdata/README.md says where each of
// Elasticsearch's answers comes from.
func TestQuerySpanFidelity(t *testing.T) {
	queries, err := readSpanCorpus(filepath.Join("testdata", "span", "queries.http"))
	if err != nil {
		t.Fatal(err)
	}
	// What a sync stores for each fixture, which omni resolves against.
	stored := map[string]syncedIndices{}
	span.Check(t, span.Suite{
		Path:     filepath.Join("testdata", "span.tsv"),
		Pins:     testcontainer.Pins,
		PinRole:  testcontainer.Driver,
		Fixtures: filepath.Join("testdata", "span", "fixtures"),
		Queries:  queries,
		Omni: func(fx span.Fixture, _, statement string) span.Answer {
			return omniSpan(stored[fx.Name], statement)
		},
		Engine: func(t testing.TB) func(span.Fixture, string, string) (span.Answer, error) {
			testcontainer.Reset(t)
			c := testcontainer.Open(t)
			return func(fx span.Fixture, _, statement string) (span.Answer, error) {
				if _, ok := stored[fx.Name]; !ok {
					s, err := buildFixture(c, fx.SQL)
					if err != nil {
						return span.Answer{}, fmt.Errorf("fixture %s: %w", fx.Name, err)
					}
					stored[fx.Name] = s
				}
				return engineSpan(c, statement)
			}
		},
	})
}

// readSpanCorpus reads queries.http: one request per line under the last
// "# fixture:" line.
func readSpanCorpus(path string) ([]span.Query, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var queries []span.Query
	fixture := ""
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		switch {
		case line == "":
		case strings.HasPrefix(line, "# fixture:"):
			fixture = strings.TrimSpace(strings.TrimPrefix(line, "# fixture:"))
		case strings.HasPrefix(line, "#"):
		default:
			if fixture == "" {
				return nil, fmt.Errorf("%s: %q comes before any fixture", path, line)
			}
			queries = append(queries, span.Query{Fixture: fixture, Session: session, Statement: line})
		}
	}
	return queries, scanner.Err()
}

// splitStatement reads one corpus line: its method, its route and its body,
// whose lines a multi-search separates with " \n ".
func splitStatement(statement string) (method, route, body string) {
	method, rest, _ := strings.Cut(statement, " ")
	route, body, _ = strings.Cut(rest, " ")
	if strings.Contains(body, ` \n `) {
		// A multi-search body is lines, the last one ended too.
		body = strings.ReplaceAll(body, ` \n `, "\n") + "\n"
	}
	return method, route, body
}

// consoleText is the statement as the console holds it, which omni parses.
func consoleText(statement string) string {
	method, route, body := splitStatement(statement)
	if body == "" {
		return method + " " + route
	}
	return method + " " + route + "\n" + body
}

// --- the fixture and what a sync stores ---

// syncedIndices is what Bytebase holds of a fixture: the indices
// elasticsearch/driver's sync lists, and the fields each index's mapping
// defines, which an ObjectSchema describes.
type syncedIndices struct {
	names  []string
	fields map[string][]string
}

// buildFixture runs the fixture's requests, separated by blank lines, and
// reads back what a sync stores.
func buildFixture(c *testcontainer.Client, text string) (syncedIndices, error) {
	for chunk := range strings.SplitSeq(text, "\n\n") {
		chunk = strings.TrimSpace(chunk)
		if chunk == "" {
			continue
		}
		head, body, _ := strings.Cut(chunk, "\n")
		method, route, _ := strings.Cut(head, " ")
		if strings.Contains(route, "_bulk") {
			// A bulk body ends with a newline.
			body += "\n"
		}
		status, answer, err := do(c, method, route, body)
		if err != nil {
			return syncedIndices{}, err
		}
		if status >= 400 || bytes.Contains(answer, []byte(`"errors":true`)) {
			return syncedIndices{}, fmt.Errorf("%s: %d %s", head, status, answer)
		}
	}
	// The sync's own two requests: every index the cat API lists, less
	// those marked hidden (elasticsearch/driver/sync.go).
	_, answer, err := do(c, "GET", "/*/_settings?filter_path=**.hidden&expand_wildcards=open,hidden", "")
	if err != nil {
		return syncedIndices{}, err
	}
	var settings map[string]struct {
		Settings struct {
			Index struct {
				Hidden string `json:"hidden"`
			} `json:"index"`
		} `json:"settings"`
	}
	if err := json.Unmarshal(answer, &settings); err != nil {
		return syncedIndices{}, fmt.Errorf("settings: %w", err)
	}
	_, answer, err = do(c, "GET", "/_cat/indices?format=json&pretty", "")
	if err != nil {
		return syncedIndices{}, err
	}
	var cat []struct {
		Index string `json:"index"`
	}
	if err := json.Unmarshal(answer, &cat); err != nil {
		return syncedIndices{}, fmt.Errorf("cat indices: %w", err)
	}
	s := syncedIndices{fields: map[string][]string{}}
	for _, e := range cat {
		if settings[e.Index].Settings.Index.Hidden == "true" {
			continue
		}
		s.names = append(s.names, e.Index)
		_, answer, err := do(c, "GET", "/"+url.PathEscape(e.Index)+"/_mapping", "")
		if err != nil {
			return syncedIndices{}, err
		}
		var mapping map[string]struct {
			Mappings mappingNode `json:"mappings"`
		}
		if err := json.Unmarshal(answer, &mapping); err != nil {
			return syncedIndices{}, fmt.Errorf("mapping of %s: %w", e.Index, err)
		}
		s.fields[e.Index] = mapping[e.Index].Mappings.leaves("")
	}
	slices.Sort(s.names)
	return s, nil
}

// mappingNode is an object's mapping: its properties, each a field or an
// object of its own. A field's multi-fields (fields) hold no _source value.
type mappingNode struct {
	Properties map[string]mappingNode `json:"properties"`
}

func (m mappingNode) leaves(prefix string) []string {
	var out []string
	for name, child := range m.Properties {
		if len(child.Properties) > 0 {
			out = append(out, child.leaves(prefix+name+".")...)
			continue
		}
		out = append(out, prefix+name)
	}
	slices.Sort(out)
	return out
}

// do sends one request and reads its answer.
func do(c *testcontainer.Client, method, route, body string) (int, []byte, error) {
	resp, err := c.Do(context.Background(), method, route, []byte(body))
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	return resp.StatusCode, b, err
}

// --- omni's answer ---

// omniSpan is the analysis Bytebase reads, resolved as Bytebase would
// resolve it: the indices the request names expanded against the synced
// indices, and the fields each one's mapping defines filtered by the
// request's _source filtering.
func omniSpan(s syncedIndices, statement string) span.Answer {
	qs, err := elasticsearch.GetQuerySpan(consoleText(statement))
	if err != nil {
		return span.Answer{Error: "omni"}
	}
	a := qs.ElasticsearchAnalysis
	if a == nil || a.API == analysis.APIUnsupported || a.API == analysis.APIBlocked {
		return span.Answer{Error: "omni: not a maskable read"}
	}
	tables := analysis.ExpandIndices(a.Indices, s.names)
	lineage := map[string][]string{}
	for _, index := range tables {
		if a.SourceDisabled {
			break
		}
		for _, field := range s.fields[index] {
			if sourceKeeps(field, a.SourceFields, a.SourceExcludes) {
				lineage[field] = append(lineage[field], index+"."+field)
			}
		}
	}
	return answerOf(tables, lineage)
}

// sourceKeeps reports whether _source filtering keeps a leaf field: an
// include or exclude pattern names the field or an object holding it, and *
// matches any run of characters, dots included.
func sourceKeeps(field string, includes, excludes []string) bool {
	names := func(patterns []string) bool {
		for _, p := range patterns {
			re := regexp.MustCompile("^" + strings.ReplaceAll(regexp.QuoteMeta(p), `\*`, ".*") + "$")
			for prefix := field; ; {
				if re.MatchString(prefix) {
					return true
				}
				i := strings.LastIndexByte(prefix, '.')
				if i < 0 {
					break
				}
				prefix = prefix[:i]
			}
		}
		return false
	}
	if len(includes) > 0 && !names(includes) {
		return false
	}
	return !names(excludes)
}

func answerOf(tables []string, lineage map[string][]string) span.Answer {
	a := span.Answer{Tables: tables}
	for field := range lineage {
		a.Results = append(a.Results, field)
	}
	slices.Sort(a.Results)
	for _, field := range a.Results {
		a.Lineage = append(a.Lineage, lineage[field])
	}
	return a
}

// --- Elasticsearch's answer ---

// backingIndex is a data stream's backing index, .ds-<stream>-<date>-<generation>.
// A request reads it through its stream, the name Bytebase sees.
var backingIndex = regexp.MustCompile(`^\.ds-(.+)-\d{4}\.\d{2}\.\d{2}-\d{6}$`)

func streamOf(index string) string {
	if m := backingIndex.FindStringSubmatch(index); m != nil {
		return m[1]
	}
	return index
}

// engineSpan is Elasticsearch's answer: the request's own answer, whose
// documents name the index each came from, and the indices _resolve/index
// resolves the request's targets to.
func engineSpan(c *testcontainer.Client, statement string) (span.Answer, error) {
	method, route, body := splitStatement(statement)
	status, answer, err := do(c, method, route, body)
	if err != nil {
		return span.Answer{}, err
	}
	if status >= 400 {
		var failure struct {
			Error struct {
				Type string `json:"type"`
			} `json:"error"`
		}
		if err := json.Unmarshal(answer, &failure); err != nil || failure.Error.Type == "" {
			return span.Answer{}, fmt.Errorf("%d %s", status, answer)
		}
		return span.Answer{Error: failure.Error.Type}, nil
	}
	path, _, _ := strings.Cut(route, "?")
	segments := strings.Split(strings.TrimPrefix(path, "/"), "/")
	target := ""
	if !strings.HasPrefix(segments[0], "_") {
		target = segments[0]
	}

	tables, err := resolveTargets(c, target, body, segments)
	if err != nil {
		return span.Answer{}, err
	}

	lineage := map[string][]string{}
	add := func(index string, source any) {
		for _, field := range sourceLeaves("", source) {
			lineage[field] = append(lineage[field], streamOf(index)+"."+field)
		}
	}
	var document map[string]any
	if err := json.Unmarshal(answer, &document); err != nil {
		return span.Answer{}, fmt.Errorf("unreadable answer %s: %w", answer, err)
	}
	hits := func(v any) {
		h, _ := v.(map[string]any)
		list, _ := h["hits"].([]any)
		for _, item := range list {
			hit, _ := item.(map[string]any)
			index, _ := hit["_index"].(string)
			add(index, hit["_source"])
		}
	}
	switch {
	case slices.Contains(segments, "_source"):
		// The source API answers the document itself; the index is the
		// one the target resolves to.
		if len(tables) != 1 {
			return span.Answer{}, fmt.Errorf("%s resolves to %v, not one index", target, tables)
		}
		add(tables[0], document)
	case document["hits"] != nil:
		hits(document["hits"])
	case document["responses"] != nil:
		list, _ := document["responses"].([]any)
		for _, item := range list {
			r, _ := item.(map[string]any)
			if e, ok := r["error"].(map[string]any); ok {
				return span.Answer{Error: fmt.Sprint(e["type"])}, nil
			}
			hits(r["hits"])
		}
	case document["docs"] != nil:
		list, _ := document["docs"].([]any)
		for _, item := range list {
			doc, _ := item.(map[string]any)
			index, _ := doc["_index"].(string)
			add(index, doc["_source"])
		}
	case document["_index"] != nil:
		index, _ := document["_index"].(string)
		add(index, document["_source"])
	default:
		return span.Answer{}, fmt.Errorf("no documents in %s", answer)
	}
	a := answerOf(tables, lineage)
	// One run shows where the documents it returned came from, never that
	// no other index could supply a field.
	for i := range a.Lineage {
		a.Lineage[i] = append(a.Lineage[i], span.Unstated)
	}
	return a, nil
}

// sourceLeaves is the dotted path of every value a document's _source holds;
// an array of objects is read through, an array of values is one leaf.
func sourceLeaves(prefix string, v any) []string {
	switch v := v.(type) {
	case map[string]any:
		var out []string
		for k, child := range v {
			out = append(out, sourceLeaves(prefix+k+".", child)...)
		}
		return out
	case []any:
		var out []string
		objects := false
		for _, item := range v {
			if _, ok := item.(map[string]any); ok {
				objects = true
				out = append(out, sourceLeaves(prefix, item)...)
			}
		}
		if objects {
			return out
		}
	case nil:
		if prefix == "" {
			return nil
		}
	}
	return []string{strings.TrimSuffix(prefix, ".")}
}

// resolveTargets is every concrete index the request's targets resolve to,
// a data stream named for itself: the path's target (every index when it
// names none) and the indices a multi-get's docs or a multi-search's headers
// name, each through _resolve/index with a search's default
// expand_wildcards=open.
func resolveTargets(c *testcontainer.Client, target, body string, segments []string) ([]string, error) {
	var targets []string
	if target != "" {
		targets = append(targets, target)
	}
	switch {
	case slices.Contains(segments, "_mget"):
		var parsed struct {
			Docs []struct {
				Index string `json:"_index"`
			} `json:"docs"`
		}
		if err := json.Unmarshal([]byte(body), &parsed); err != nil {
			return nil, err
		}
		for _, d := range parsed.Docs {
			if d.Index != "" {
				targets = append(targets, url.PathEscape(d.Index))
			}
		}
	case slices.Contains(segments, "_msearch"):
		lines := strings.Split(strings.TrimSuffix(body, "\n"), "\n")
		for i := 0; i < len(lines); i += 2 {
			var header struct {
				Index string `json:"index"`
			}
			if err := json.Unmarshal([]byte(lines[i]), &header); err != nil {
				return nil, err
			}
			if header.Index != "" {
				targets = append(targets, url.PathEscape(header.Index))
			}
		}
	}
	if len(targets) == 0 {
		targets = []string{"*"}
	}
	var out []string
	for _, t := range targets {
		status, answer, err := do(c, "GET", "/_resolve/index/"+t+"?expand_wildcards=open", "")
		if err != nil {
			return nil, err
		}
		if status >= 400 {
			return nil, fmt.Errorf("resolve %s: %d %s", t, status, answer)
		}
		var resolved struct {
			Indices []struct {
				Name       string `json:"name"`
				DataStream string `json:"data_stream"`
			} `json:"indices"`
			Aliases []struct {
				Indices []string `json:"indices"`
			} `json:"aliases"`
			DataStreams []struct {
				Name string `json:"name"`
			} `json:"data_streams"`
		}
		if err := json.Unmarshal(answer, &resolved); err != nil {
			return nil, fmt.Errorf("resolve %s: %w", t, err)
		}
		for _, i := range resolved.Indices {
			if i.DataStream != "" {
				out = append(out, i.DataStream)
				continue
			}
			out = append(out, i.Name)
		}
		for _, a := range resolved.Aliases {
			for _, i := range a.Indices {
				out = append(out, streamOf(i))
			}
		}
		for _, d := range resolved.DataStreams {
			out = append(out, d.Name)
		}
	}
	slices.Sort(out)
	return slices.Compact(out), nil
}
