package testcontainer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"

	gocql "github.com/apache/cassandra-gocql-driver/v2"

	"github.com/bytebase/omni/enginetest/tape"
)

// Spans asks the pinned Cassandra what a query reads, for cassandra/analysis's
// span.tsv. Its tape holds one raw exchange per question, and each answer is
// what the server stated, reduced to JSON:
//
//   - Fixture empties the server of every user keyspace, runs a starting
//     schema's statements, and answers what system_schema then holds for the
//     user keyspaces: each table's and materialized view's columns, as
//     cassandra/driver's sync reads system_schema.columns, and the base table
//     of each materialized view (system_schema.views).
//   - Span runs a query in a session keyspace and answers the result
//     metadata Cassandra sent for it: each output column's keyspace, table and
//     name. gocql prepares the query, so this is the prepared statement's
//     result metadata; the fixture's tables are empty, so no row follows it.
type Spans struct {
	raw *tape.Raw
}

// OpenSpans returns t's Spans on the driver pin: a replay of its tape, or the
// pinned server in a recording.
func OpenSpans(t testing.TB) *Spans {
	t.Helper()
	return &Spans{raw: tape.OpenRaw(t, server, Driver, func(request string) (string, error) {
		return liveSpan(request)
	})}
}

// SpanFixture is the schema Cassandra stored for a fixture: its user
// keyspaces, in name order.
type SpanFixture struct {
	Keyspaces []SpanKeyspace `json:"keyspaces"`
}

// SpanKeyspace is one keyspace and its relations, in name order.
type SpanKeyspace struct {
	Name   string         `json:"name"`
	Tables []SpanRelation `json:"tables,omitempty"`
	Views  []SpanRelation `json:"views,omitempty"`
}

// SpanRelation is a table or a materialized view, with its columns as
// system_schema.columns lists them; Base is a materialized view's base table.
type SpanRelation struct {
	Name    string       `json:"name"`
	Base    string       `json:"base,omitempty"`
	Columns []SpanColumn `json:"columns"`
}

// SpanColumn is one row of system_schema.columns.
type SpanColumn struct {
	Name     string `json:"name"`
	Kind     string `json:"kind"`
	Position int    `json:"position"`
	Type     string `json:"type"`
}

// SpanAnswer is Cassandra's answer to one query: the error it refused it with,
// or its result metadata.
type SpanAnswer struct {
	Error   string           `json:"error,omitempty"`
	Columns []SpanColumnSpec `json:"columns,omitempty"`
}

// SpanColumnSpec is one column of a result's metadata: the keyspace and table
// Cassandra names for it, and its name.
type SpanColumnSpec struct {
	Keyspace string `json:"keyspace"`
	Table    string `json:"table"`
	Name     string `json:"name"`
}

const (
	opSpanFixture = "spanFixture"
	opSpan        = "span"
)

type spanRequest struct {
	Op         string   `json:"op"`
	Statements []string `json:"statements,omitempty"`
	Keyspace   string   `json:"keyspace,omitempty"`
	Statement  string   `json:"statement,omitempty"`
}

// Fixture builds a starting schema on an empty server and returns what
// Cassandra stored for it.
func (s *Spans) Fixture(statements []string) (SpanFixture, error) {
	var fx SpanFixture
	err := s.do(spanRequest{Op: opSpanFixture, Statements: statements}, &fx)
	return fx, err
}

// Span returns Cassandra's answer to statement, run in keyspace.
func (s *Spans) Span(keyspace, statement string) (SpanAnswer, error) {
	var a SpanAnswer
	err := s.do(spanRequest{Op: opSpan, Keyspace: keyspace, Statement: statement}, &a)
	return a, err
}

func (s *Spans) do(r spanRequest, out any) error {
	b, err := json.Marshal(r)
	if err != nil {
		return err
	}
	answered, err := s.raw.Do(string(b))
	if err != nil {
		return err
	}
	if err := unmarshalJSON([]byte(answered), out); err != nil {
		return fmt.Errorf("testcontainer: unreadable answer to %s: %w", b, err)
	}
	return nil
}

// liveSpan answers one span request from the pinned server.
func liveSpan(text string) (string, error) {
	var req spanRequest
	if err := unmarshalJSON([]byte(text), &req); err != nil {
		return "", fmt.Errorf("testcontainer: unreadable request %q: %w", text, err)
	}
	var out any
	switch req.Op {
	case opSpanFixture:
		fx, err := liveSpanFixture(req.Statements)
		if err != nil {
			return "", err
		}
		out = fx
	case opSpan:
		a, err := liveSpanQuery(req.Keyspace, req.Statement)
		if err != nil {
			return "", err
		}
		out = a
	default:
		return "", fmt.Errorf("testcontainer: unknown span operation %q", req.Op)
	}
	b, err := json.Marshal(out)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func liveSpanFixture(statements []string) (SpanFixture, error) {
	session, err := liveSession("")
	if err != nil {
		return SpanFixture{}, err
	}
	ctx := context.Background()
	names, err := userKeyspaces(ctx, session)
	if err != nil {
		return SpanFixture{}, err
	}
	for _, name := range names {
		if err := session.Query(`DROP KEYSPACE "` + name + `"`).ExecContext(ctx); err != nil {
			return SpanFixture{}, fmt.Errorf("drop keyspace %s: %w", name, err)
		}
	}
	for _, statement := range statements {
		if err := session.Query(statement).ExecContext(ctx); err != nil {
			return SpanFixture{}, fmt.Errorf("%s:\n%w", statement, err)
		}
	}

	names, err = userKeyspaces(ctx, session)
	if err != nil {
		return SpanFixture{}, err
	}
	user := map[string]bool{}
	for _, name := range names {
		user[name] = true
	}
	type key struct{ keyspace, relation string }
	columns := map[key][]SpanColumn{}
	iter := session.Query("SELECT keyspace_name, table_name, column_name, kind, position, type FROM system_schema.columns").IterContext(ctx)
	var c SpanColumn
	var keyspace, relation string
	for iter.Scan(&keyspace, &relation, &c.Name, &c.Kind, &c.Position, &c.Type) {
		if user[keyspace] {
			columns[key{keyspace, relation}] = append(columns[key{keyspace, relation}], c)
		}
	}
	if err := iter.Close(); err != nil {
		return SpanFixture{}, fmt.Errorf("read system_schema.columns: %w", err)
	}
	relations := func(query string, view bool) (map[string][]SpanRelation, error) {
		out := map[string][]SpanRelation{}
		iter := session.Query(query).IterContext(ctx)
		var keyspace, name, base string
		targets := []any{&keyspace, &name}
		if view {
			targets = append(targets, &base)
		}
		for iter.Scan(targets...) {
			if !user[keyspace] {
				continue
			}
			cols := slices.Clone(columns[key{keyspace, name}])
			slices.SortFunc(cols, func(a, b SpanColumn) int { return strings.Compare(a.Name, b.Name) })
			out[keyspace] = append(out[keyspace], SpanRelation{Name: name, Base: base, Columns: cols})
		}
		if err := iter.Close(); err != nil {
			return nil, fmt.Errorf("%s: %w", query, err)
		}
		return out, nil
	}
	tables, err := relations("SELECT keyspace_name, table_name FROM system_schema.tables", false)
	if err != nil {
		return SpanFixture{}, err
	}
	views, err := relations("SELECT keyspace_name, view_name, base_table_name FROM system_schema.views", true)
	if err != nil {
		return SpanFixture{}, err
	}
	var fx SpanFixture
	for _, name := range names {
		ks := SpanKeyspace{Name: name, Tables: tables[name], Views: views[name]}
		for _, list := range [][]SpanRelation{ks.Tables, ks.Views} {
			slices.SortFunc(list, func(a, b SpanRelation) int { return strings.Compare(a.Name, b.Name) })
		}
		fx.Keyspaces = append(fx.Keyspaces, ks)
	}
	return fx, nil
}

// userKeyspaces lists the server's keyspaces that are not Cassandra's own, in
// name order.
func userKeyspaces(ctx context.Context, session *gocql.Session) ([]string, error) {
	iter := session.Query("SELECT keyspace_name FROM system_schema.keyspaces").IterContext(ctx)
	var name string
	var names []string
	for iter.Scan(&name) {
		if !strings.HasPrefix(name, "system") {
			names = append(names, name)
		}
	}
	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("list keyspaces: %w", err)
	}
	slices.Sort(names)
	return names, nil
}

func liveSpanQuery(keyspace, statement string) (SpanAnswer, error) {
	session, err := liveSession(keyspace)
	if err != nil {
		return SpanAnswer{}, err
	}
	iter := session.Query(statement).IterContext(context.Background())
	columns := iter.Columns()
	if err := iter.Close(); err != nil {
		var requestErr gocql.RequestError
		if errors.As(err, &requestErr) {
			return SpanAnswer{Error: fmt.Sprintf("%#x %s", requestErr.Code(), requestErr.Message())}, nil
		}
		return SpanAnswer{}, err
	}
	var a SpanAnswer
	for _, c := range columns {
		a.Columns = append(a.Columns, SpanColumnSpec{Keyspace: c.Keyspace, Table: c.Table, Name: c.Name})
	}
	return a, nil
}
