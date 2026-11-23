package testcontainer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/tape"
)

// Spans asks the pinned Trino what a query reads, for trino/analysis's
// span.tsv. Its tape holds one raw exchange per question, and each answer is
// what the server stated, reduced to JSON:
//
//   - Fixture builds a starting schema in the memory catalog and answers the
//     schema Trino stored, read as trino/driver syncs it, and the lineage of
//     each view's definition.
//   - Span runs a query in a session and answers its output column names
//     (the statement protocol's columns), the tables and views it read and
//     whether each is a view (the query's QueryInfo referencedTables, from
//     Trino's analyzer), and each output column's source columns (the
//     QueryInfo output columns of a CREATE TABLE AS over the query, Trino's
//     own column lineage).
type Spans struct {
	raw *tape.Raw
}

// OpenSpans returns t's Spans: a replay of its tape, or the pinned Trino in
// a recording.
func OpenSpans(t testing.TB) *Spans {
	t.Helper()
	var live func(string) (string, error)
	mode, err := enginetest.RunMode()
	if err != nil {
		t.Fatalf("%v", err)
	}
	if mode != enginetest.ModeReplay {
		pin, err := enginetest.LoadPin(Pins, Default)
		if err != nil {
			t.Fatalf("%v", err)
		}
		c := newClient(started(t, pin))
		emptyOnce(t, c)
		live = func(req string) (string, error) { return c.spanLive(context.Background(), req) }
	}
	return &Spans{raw: tape.OpenRaw(t, server, Default, live)}
}

// SpanFixture is the schema Trino stored for a fixture.
type SpanFixture struct {
	Schemas []SpanSchema `json:"schemas"`
	// Views maps each view, catalog.schema.view, to its definition's
	// lineage: for each of its columns, the columns its values come from.
	Views map[string][][]string `json:"views"`
}

// SpanSchema is one schema of the memory catalog.
type SpanSchema struct {
	Name   string         `json:"name"`
	Tables []SpanRelation `json:"tables"`
	Views  []SpanRelation `json:"views"`
}

// SpanRelation is a table or a view and its columns in order.
type SpanRelation struct {
	Name       string       `json:"name"`
	Columns    []SpanColumn `json:"columns"`
	Definition string       `json:"definition,omitempty"`
}

// SpanColumn is a column and its type, as system.jdbc.columns names it.
type SpanColumn struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// SpanAnswer is what Trino states about one query.
type SpanAnswer struct {
	// Error is Trino's error name for a query it refused; the rest is empty.
	Error string `json:"error,omitempty"`
	// Columns is the output column names, in order.
	Columns []string `json:"columns,omitempty"`
	// Tables is every table and view the query read, as catalog.schema.name.
	Tables []SpanTable `json:"tables,omitempty"`
	// Lineage is each output column's source columns, as
	// catalog.schema.table.column. A column read from a view names the
	// view's column.
	Lineage [][]string `json:"lineage,omitempty"`
	// LineageError is Trino's error name when the CREATE TABLE AS that
	// reads the lineage was refused; Lineage is then empty.
	LineageError string `json:"lineageError,omitempty"`
}

// SpanTable is one relation a query read.
type SpanTable struct {
	Name string `json:"name"`
	View bool   `json:"view,omitempty"`
}

// Fixture builds the starting schema sql, statements ending in ";" at the
// end of a line, in an emptied memory catalog, and returns what Trino stored.
// The statements run in memory.default until a USE catalog.schema names
// another session.
func (s *Spans) Fixture(sql string) (SpanFixture, error) {
	var fx SpanFixture
	ans, err := s.raw.Do("fixture\n" + sql)
	if err != nil {
		return fx, err
	}
	return fx, json.Unmarshal([]byte(ans), &fx)
}

// Span runs statement in session, catalog.schema, and returns what Trino
// states it read.
func (s *Spans) Span(session, statement string) (SpanAnswer, error) {
	var a SpanAnswer
	ans, err := s.raw.Do("span " + session + "\n" + statement)
	if err != nil {
		return a, err
	}
	return a, json.Unmarshal([]byte(ans), &a)
}

// lineageSchema holds the tables Span creates to read lineage. Fixture makes
// it after reading the stored schema, so the schema a fixture stores never
// lists it.
const lineageSchema = "span_lineage"

func (c *client) spanLive(ctx context.Context, req string) (string, error) {
	head, body, _ := strings.Cut(req, "\n")
	var (
		v   any
		err error
	)
	switch kind, sess, _ := strings.Cut(head, " "); kind {
	case "fixture":
		v, err = c.spanFixture(ctx, body)
	case "span":
		catalog, schema, ok := strings.Cut(sess, ".")
		if !ok {
			return "", fmt.Errorf("testcontainer: span session %q is not catalog.schema", sess)
		}
		v, err = c.span(ctx, session{catalog, schema}, body)
	default:
		return "", fmt.Errorf("testcontainer: unknown span request %q", head)
	}
	if err != nil {
		return "", err
	}
	b, err := json.Marshal(v)
	return string(b), err
}

// session is the catalog and schema a statement runs in.
type session struct{ catalog, schema string }

var memoryDefault = session{"memory", "default"}

func (c *client) spanFixture(ctx context.Context, sql string) (SpanFixture, error) {
	var fx SpanFixture
	if err := c.emptyMemory(ctx); err != nil {
		return fx, err
	}
	in := memoryDefault
	for _, stmt := range strings.Split(sql, ";\n") {
		if stmt = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(stmt), ";")); stmt == "" {
			continue
		}
		// USE is the client's to apply: the server only answers the session
		// it asks the client to set. A view stores the session it was made
		// in, which resolves its definition's names.
		if name, ok := strings.CutPrefix(stmt, "USE "); ok {
			catalog, schema, ok := strings.Cut(name, ".")
			if !ok {
				return fx, fmt.Errorf("fixture statement %q: name the catalog and the schema", stmt)
			}
			in = session{catalog, schema}
			continue
		}
		if _, err := c.statement(ctx, in, stmt); err != nil {
			return fx, fmt.Errorf("fixture statement %q: %w", stmt, err)
		}
	}
	// What follows reads the stored schema as trino/driver's sync does.
	schemas, err := c.query(ctx, "SELECT table_schem FROM system.jdbc.schemas WHERE table_catalog = 'memory' AND table_schem <> 'information_schema' ORDER BY table_schem")
	if err != nil {
		return fx, err
	}
	bySchema := map[string]*SpanSchema{}
	for _, r := range schemas {
		fx.Schemas = append(fx.Schemas, SpanSchema{Name: fmt.Sprint(r[0])})
	}
	for i := range fx.Schemas {
		bySchema[fx.Schemas[i].Name] = &fx.Schemas[i]
	}
	columns, err := c.query(ctx, "SELECT table_schem, table_name, column_name, type_name FROM system.jdbc.columns WHERE table_cat = 'memory' ORDER BY table_schem, table_name, ordinal_position")
	if err != nil {
		return fx, err
	}
	cols := map[[2]string][]SpanColumn{}
	for _, r := range columns {
		k := [2]string{fmt.Sprint(r[0]), fmt.Sprint(r[1])}
		cols[k] = append(cols[k], SpanColumn{Name: fmt.Sprint(r[2]), Type: fmt.Sprint(r[3])})
	}
	tables, err := c.query(ctx, "SELECT table_schem, table_name FROM system.jdbc.tables WHERE table_cat = 'memory' AND table_type = 'TABLE' ORDER BY table_schem, table_name")
	if err != nil {
		return fx, err
	}
	for _, r := range tables {
		s, n := fmt.Sprint(r[0]), fmt.Sprint(r[1])
		if sc := bySchema[s]; sc != nil {
			sc.Tables = append(sc.Tables, SpanRelation{Name: n, Columns: cols[[2]string{s, n}]})
		}
	}
	views, err := c.query(ctx, "SELECT table_schema, table_name, view_definition FROM memory.information_schema.views ORDER BY table_schema, table_name")
	if err != nil {
		return fx, err
	}
	// The stored schema is read: make the schema lineage tables go in.
	if _, err := c.statement(ctx, memoryDefault, "CREATE SCHEMA memory."+lineageSchema); err != nil {
		return fx, err
	}
	fx.Views = map[string][][]string{}
	for _, r := range views {
		s, n, def := fmt.Sprint(r[0]), fmt.Sprint(r[1]), fmt.Sprint(r[2])
		sc := bySchema[s]
		if sc == nil {
			continue
		}
		vcols := cols[[2]string{s, n}]
		sc.Views = append(sc.Views, SpanRelation{Name: n, Columns: vcols, Definition: def})
		// A view's definition runs in the view's own catalog and schema.
		lineage, errName, err := c.lineage(ctx, session{"memory", s}, def, len(vcols))
		if err != nil {
			return fx, err
		}
		if errName != "" {
			return fx, fmt.Errorf("the lineage of view %s.%s: %s", s, n, errName)
		}
		fx.Views["memory."+s+"."+n] = lineage
	}
	return fx, nil
}

func (c *client) span(ctx context.Context, in session, sql string) (SpanAnswer, error) {
	var a SpanAnswer
	res, err := c.statement(ctx, in, sql)
	var qe *QueryError
	switch {
	case errors.As(err, &qe):
		a.Error = qe.Name
		return a, nil
	case err != nil:
		return a, err
	}
	a.Columns = res.columns
	info, err := c.queryInfo(ctx, res.id)
	if err != nil {
		return a, err
	}
	for _, t := range info.ReferencedTables {
		a.Tables = append(a.Tables, SpanTable{Name: t.Catalog + "." + t.Schema + "." + t.Table, View: t.ViewText != nil})
	}
	a.Lineage, a.LineageError, err = c.lineage(ctx, in, sql, len(res.columns))
	return a, err
}

// lineage returns the source columns of each of query's n output columns,
// as Trino's analyzer records them for a CREATE TABLE AS over it, or the
// error name Trino refused that statement with. The query is wrapped in a
// derived table whose column list names every column, so an unnamed or
// repeated output column still makes a table.
func (c *client) lineage(ctx context.Context, in session, query string, n int) ([][]string, string, error) {
	names := make([]string, n)
	for i := range names {
		names[i] = fmt.Sprintf("c%d", i+1)
	}
	target := "memory." + lineageSchema + ".q"
	ctas := fmt.Sprintf("CREATE TABLE %s AS SELECT * FROM (\n%s\n) AS q(%s)", target, query, strings.Join(names, ", "))
	res, err := c.statement(ctx, in, ctas)
	var qe *QueryError
	switch {
	case errors.As(err, &qe):
		return nil, qe.Name, nil
	case err != nil:
		return nil, "", err
	}
	info, err := c.queryInfo(ctx, res.id)
	if err != nil {
		return nil, "", err
	}
	if _, err := c.statement(ctx, memoryDefault, "DROP TABLE "+target); err != nil {
		return nil, "", err
	}
	if info.Output == nil || len(info.Output.Columns) != n {
		return nil, "", fmt.Errorf("testcontainer: CREATE TABLE AS reported no lineage for %d columns", n)
	}
	lineage := make([][]string, n)
	for i, col := range info.Output.Columns {
		lineage[i] = []string{}
		for _, src := range col.SourceColumns {
			lineage[i] = append(lineage[i], src.TableName+"."+src.ColumnName)
		}
	}
	return lineage, "", nil
}

// result is what the statement protocol says about a statement run to its end.
type result struct {
	id      string
	columns []string
}

// statement runs sql in session to its end.
func (c *client) statement(ctx context.Context, in session, sql string) (result, error) {
	resp, err := c.postIn(ctx, sql, in)
	if err != nil {
		return result{}, err
	}
	res := result{id: resp.ID}
	for i := 0; i < 10000; i++ {
		if qe := resp.queryError(); qe != nil {
			return result{}, qe
		}
		if res.columns == nil && resp.Columns != nil {
			res.columns = make([]string, len(resp.Columns))
			for i, col := range resp.Columns {
				res.columns[i] = col.Name
			}
		}
		if resp.NextURI == "" {
			return res, nil
		}
		select {
		case <-ctx.Done():
			return result{}, ctx.Err()
		case <-time.After(20 * time.Millisecond):
		}
		if resp, err = c.get(ctx, resp.NextURI); err != nil {
			return result{}, err
		}
	}
	return result{}, fmt.Errorf("trino: statement %q did not settle", truncate(sql, 60))
}

// queryInfo is the part of /v1/query/{id} the span answers read.
type queryInfo struct {
	ReferencedTables []struct {
		Catalog  string  `json:"catalog"`
		Schema   string  `json:"schema"`
		Table    string  `json:"table"`
		ViewText *string `json:"viewText"`
	} `json:"referencedTables"`
	Output *struct {
		Columns []struct {
			SourceColumns []struct {
				TableName  string `json:"tableName"`
				ColumnName string `json:"columnName"`
			} `json:"sourceColumns"`
		} `json:"columns"`
	} `json:"output"`
}

func (c *client) queryInfo(ctx context.Context, id string) (queryInfo, error) {
	var info queryInfo
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/v1/query/"+id, nil)
	if err != nil {
		return info, err
	}
	req.Header.Set("X-Trino-User", "omni-oracle")
	resp, err := c.http.Do(req)
	if err != nil {
		return info, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return info, fmt.Errorf("trino /v1/query/%s: status %d", id, resp.StatusCode)
	}
	return info, json.NewDecoder(resp.Body).Decode(&info)
}
