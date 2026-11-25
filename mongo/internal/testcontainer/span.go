package testcontainer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/bytebase/gomongo"

	"github.com/bytebase/omni/enginetest/tape"
	omnimongo "github.com/bytebase/omni/mongo"
	"github.com/bytebase/omni/mongo/ast"
)

// Spans asks the pinned MongoDB what a query reads, for mongo/analysis's
// span.tsv. Its tape holds one raw exchange per question, and each answer is
// what the server stated, reduced to JSON:
//
//   - SpanFixture empties the databases a fixture names, runs its statements,
//     and answers the views each database then holds, with their options as
//     listCollections reports them and mongo/driver's sync stores them.
//   - Span plans a query with explain at queryPlanner verbosity, which does
//     not run it, and answers the collections the plan reads.
type Spans struct {
	raw *tape.Raw
}

// OpenSpans returns t's Spans: a replay of its tape, or the pinned MongoDB in
// a recording.
func OpenSpans(t testing.TB) *Spans {
	t.Helper()
	return &Spans{raw: tape.OpenRaw(t, server, Default, func(request string) (string, error) {
		client := live()
		if client == nil {
			t.Fatalf("start MongoDB: %v", liveServer.err)
		}
		emptyOnce(t, client)
		return spanLive(context.Background(), client, request)
	})}
}

// SpanView is a view a fixture made: its database, its name, and its options
// as canonical Extended JSON.
type SpanView struct {
	Database   string `json:"database"`
	Name       string `json:"name"`
	Definition string `json:"definition"`
}

// SpanFixture is what MongoDB stored for a fixture that omni's span reads.
type SpanFixture struct {
	Views []SpanView `json:"views"`
}

// SpanAnswer is MongoDB's answer about one query: the collections its plan
// reads, as database.collection, or the name of the error it refused the
// query with.
type SpanAnswer struct {
	Error  string   `json:"error,omitempty"`
	Tables []string `json:"tables,omitempty"`
}

type spanRequest struct {
	Op        string `json:"op"`
	Database  string `json:"database,omitempty"`
	Statement string `json:"statement"`
}

const (
	opSpanFixture = "SpanFixture"
	opSpan        = "Span"
)

// SpanFixture builds a fixture. Its script holds one mongosh statement per
// line, run in the database the last "// database: <name>" line names; other
// lines starting with // are comments.
func (s *Spans) SpanFixture(script string) (SpanFixture, error) {
	var out SpanFixture
	err := s.do(spanRequest{Op: opSpanFixture, Statement: script}, &out)
	return out, err
}

// Span asks what statement reads when run in database.
func (s *Spans) Span(database, statement string) (SpanAnswer, error) {
	var out SpanAnswer
	err := s.do(spanRequest{Op: opSpan, Database: database, Statement: statement}, &out)
	return out, err
}

func (s *Spans) do(req spanRequest, out any) error {
	b, err := json.Marshal(req)
	if err != nil {
		return err
	}
	answered, err := s.raw.Do(string(b))
	if err != nil {
		return err
	}
	return unmarshalJSON([]byte(answered), out)
}

func spanLive(ctx context.Context, client *mongo.Client, text string) (string, error) {
	var req spanRequest
	if err := unmarshalJSON([]byte(text), &req); err != nil {
		return "", fmt.Errorf("testcontainer: unreadable request %q: %w", text, err)
	}
	var out any
	var err error
	switch req.Op {
	case opSpanFixture:
		out, err = liveSpanFixture(ctx, client, req.Statement)
	case opSpan:
		out, err = liveSpan(ctx, client, req.Database, req.Statement)
	default:
		return "", fmt.Errorf("testcontainer: unknown span operation %q", req.Op)
	}
	if err != nil {
		return "", err
	}
	b, err := json.Marshal(out)
	return string(b), err
}

func liveSpanFixture(ctx context.Context, client *mongo.Client, script string) (SpanFixture, error) {
	type step struct{ database, statement string }
	var steps []step
	var databases []string
	database := ""
	for _, line := range strings.Split(script, "\n") {
		line = strings.TrimSpace(line)
		switch {
		case line == "":
		case strings.HasPrefix(line, "// database:"):
			database = strings.TrimSpace(strings.TrimPrefix(line, "// database:"))
			databases = append(databases, database)
		case strings.HasPrefix(line, "//"):
		case database == "":
			return SpanFixture{}, fmt.Errorf("a statement before the first // database: line: %s", line)
		default:
			steps = append(steps, step{database, line})
		}
	}
	for _, name := range databases {
		if err := client.Database(name).Drop(ctx); err != nil {
			return SpanFixture{}, fmt.Errorf("drop %s: %w", name, err)
		}
	}
	shell := gomongo.NewClient(client)
	for _, s := range steps {
		if _, err := shell.Execute(ctx, s.database, s.statement); err != nil {
			return SpanFixture{}, fmt.Errorf("in %s: %s: %w", s.database, s.statement, err)
		}
	}
	out := SpanFixture{Views: []SpanView{}}
	seen := map[string]bool{}
	for _, name := range databases {
		if seen[name] {
			continue
		}
		seen[name] = true
		infos, err := liveListCollections(ctx, client, name)
		if err != nil {
			return SpanFixture{}, err
		}
		for _, info := range infos {
			if kind, _ := info["type"].(string); kind != "view" {
				continue
			}
			view, _ := info["name"].(string)
			definition, err := bson.MarshalExtJSON(ordered(info["options"]), true, false)
			if err != nil {
				return SpanFixture{}, err
			}
			out.Views = append(out.Views, SpanView{Database: name, Name: view, Definition: string(definition)})
		}
	}
	return out, nil
}

// liveSpan plans statement and reads the collections the plan names. A plan
// names the collections it scans (queryPlanner.namespace, one per pipeline a
// $unionWith or a lowered $lookup runs), and states each stage that joins
// another collection with that collection's name: $lookup and $graphLookup's
// from, $unionWith's coll, and an EQ_LOOKUP's foreignCollection. A view the
// query targets is already resolved in the plan; one a stage names is not,
// so it is planned in turn, and the answer names the collections under it.
func liveSpan(ctx context.Context, client *mongo.Client, database, statement string) (SpanAnswer, error) {
	database, planned, err := explainable(database, statement)
	if err != nil {
		return SpanAnswer{}, err
	}
	shell := gomongo.NewClient(client)
	kinds := map[string]map[string]string{}
	kind := func(ns string) (string, error) {
		db, coll, _ := strings.Cut(ns, ".")
		if kinds[db] == nil {
			infos, err := liveListCollections(ctx, client, db)
			if err != nil {
				return "", err
			}
			kinds[db] = map[string]string{}
			for _, info := range infos {
				name, _ := info["name"].(string)
				k, _ := info["type"].(string)
				kinds[db][name] = k
			}
		}
		return kinds[db][coll], nil
	}
	tables := map[string]bool{}
	expanded := map[string]bool{}
	type pending struct{ database, statement string }
	queue := []pending{{database, planned}}
	first := true
	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		result, err := shell.Execute(ctx, p.database, p.statement)
		if err != nil {
			var ce mongo.CommandError
			if first && errors.As(err, &ce) && ce.Name != "" {
				return SpanAnswer{Error: ce.Name}, nil
			}
			return SpanAnswer{}, fmt.Errorf("explain %s in %s: %w", p.statement, p.database, err)
		}
		first = false
		if len(result.Value) != 1 {
			return SpanAnswer{}, fmt.Errorf("explain %s: %d values", p.statement, len(result.Value))
		}
		var names []string
		planNamespaces(result.Value[0], p.database, &names)
		for _, ns := range names {
			k, err := kind(ns)
			if err != nil {
				return SpanAnswer{}, err
			}
			if k != "view" {
				tables[ns] = true
				continue
			}
			if expanded[ns] {
				continue
			}
			expanded[ns] = true
			db, coll, _ := strings.Cut(ns, ".")
			queue = append(queue, pending{db, fmt.Sprintf("db.getCollection(%q).aggregate([]).explain(%q)", coll, "queryPlanner")})
		}
	}
	out := SpanAnswer{}
	for ns := range tables {
		out.Tables = append(out.Tables, ns)
	}
	sort.Strings(out.Tables)
	return out, nil
}

// explainable returns the database a statement reads in and the statement
// that plans it. gomongo, as Bytebase runs it, takes a database rather than
// db.getSiblingDB, and does not plan findOne, which MongoDB plans as a find.
func explainable(database, statement string) (string, string, error) {
	statements, err := omnimongo.Parse(statement)
	if err != nil {
		return "", "", err
	}
	if len(statements) != 1 {
		return "", "", fmt.Errorf("want one statement, got %d", len(statements))
	}
	text := statement
	if s, ok := statements[0].AST.(*ast.DatabaseStatement); ok && s.Method == "getSiblingDB" && len(s.Args) == 1 && len(s.ChainedMethods) > 0 {
		name, ok := s.Args[0].(*ast.StringLiteral)
		if !ok {
			return "", "", fmt.Errorf("getSiblingDB of a non-string: %s", statement)
		}
		database = name.Value
		text = "db." + statement[s.ChainedMethods[0].Loc.Start:s.Loc.End]
		if statements, err = omnimongo.Parse(text); err != nil {
			return "", "", err
		}
	}
	s, ok := statements[0].AST.(*ast.CollectionStatement)
	if !ok {
		return "", "", fmt.Errorf("not a collection statement: %s", statement)
	}
	text = text[s.Loc.Start:s.Loc.End]
	explain := `explain("queryPlanner")`
	planned := false
	for _, cursor := range s.CursorMethods {
		if cursor.Method == "explain" {
			// A statement that asks for its own plan is planned without
			// running it, whatever verbosity it names.
			text = text[:cursor.Loc.Start-s.Loc.Start] + explain + text[cursor.Loc.End-s.Loc.Start:]
			planned = true
			break
		}
	}
	if s.Method == "findOne" {
		i := strings.Index(text[s.CollectionLoc.End-s.Loc.Start:], ".findOne(")
		if i < 0 {
			return "", "", fmt.Errorf("no findOne call in %s", statement)
		}
		i += s.CollectionLoc.End - s.Loc.Start
		text = text[:i] + ".find(" + text[i+len(".findOne("):]
	}
	if !planned {
		text += "." + explain
	}
	return database, text, nil
}

// planNamespaces appends the collections plan names, as database.collection,
// with a stage's bare collection name read in database. It reads only the
// plan's structure: the query plan's stages and the pipeline's, never a
// filter or a stage's own arguments, where a document is the query's data.
func planNamespaces(plan any, database string, out *[]string) {
	type pending struct {
		value    any
		pipeline bool // an array of pipeline stages, not plan stages
	}
	stack := []pending{{value: plan}}
	for len(stack) > 0 {
		p := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if p.pipeline {
			for _, stage := range asArray(p.value) {
				doc := asDoc(stage)
				if len(doc) != 1 {
					continue
				}
				spec := asDoc(doc[0].Value)
				switch doc[0].Key {
				case "$cursor":
					stack = append(stack, pending{value: spec})
				case "$lookup", "$graphLookup":
					switch from := lookup(spec, "from").(type) {
					case string:
						*out = append(*out, database+"."+from)
					case bson.D:
						db, _ := lookup(from, "db").(string)
						coll, _ := lookup(from, "coll").(string)
						*out = append(*out, db+"."+coll)
					}
					if sub := lookup(spec, "pipeline"); sub != nil {
						stack = append(stack, pending{value: sub, pipeline: true})
					}
				case "$unionWith":
					if coll, ok := lookup(spec, "coll").(string); ok {
						*out = append(*out, database+"."+coll)
					}
					if sub := lookup(spec, "pipeline"); sub != nil {
						stack = append(stack, pending{value: sub, pipeline: true})
					}
				case "$facet":
					for _, branch := range spec {
						stack = append(stack, pending{value: branch.Value, pipeline: true})
					}
				}
			}
			continue
		}
		for _, e := range asDoc(p.value) {
			switch e.Key {
			case "namespace", "foreignCollection":
				if ns, ok := e.Value.(string); ok {
					*out = append(*out, ns)
				}
			case "stages":
				stack = append(stack, pending{value: e.Value, pipeline: true})
			case "queryPlanner", "winningPlan", "queryPlan", "inputStage", "outerStage", "innerStage", "thenStage", "elseStage":
				stack = append(stack, pending{value: e.Value})
			case "rejectedPlans", "inputStages":
				for _, child := range asArray(e.Value) {
					stack = append(stack, pending{value: child})
				}
			}
		}
	}
}

// asDoc returns v as an ordered document, or nil.
func asDoc(v any) bson.D {
	switch x := v.(type) {
	case bson.D:
		return x
	case bson.M:
		d, _ := ordered(x).(bson.D)
		return d
	}
	return nil
}

// asArray returns v as an array, or nil.
func asArray(v any) []any {
	switch x := v.(type) {
	case bson.A:
		return x
	case []any:
		return x
	}
	return nil
}

func lookup(doc bson.D, key string) any {
	for _, e := range doc {
		if e.Key == key {
			return e.Value
		}
	}
	return nil
}
