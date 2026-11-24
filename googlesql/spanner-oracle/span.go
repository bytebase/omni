package main

// The span mode answers what googlesql/analysis's span suite asks
// (googlesql/analysis/testdata/README.md): a fixture's stored schema, and a
// query's output column names and the tables and views it reads.
//
// The emulator states neither a query plan (AnalyzeQuery answers "query plan
// unavailable") nor a privilege check (it ignores fine-grained access
// control), so the tables a query reads are found by name resolution: for
// each table and view of the fixture, a database where that object alone is
// gone, renamed away (a table) or dropped (a view), with the views that
// read it dropped first. The query reads the object exactly when it fails
// there.

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"cloud.google.com/go/spanner"
	"cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// spanFixturePrefix starts a request to build a fixture: the rest is its
	// DDL, statements ending with a semicolon at the end of a line.
	spanFixturePrefix = "-- omni:span-fixture\n"
	// spanQueryPrefix starts a request for a query's span on the last
	// fixture built.
	spanQueryPrefix = "-- omni:span-query\n"
)

// spanAnswer is the harness's answer to a span request: a fixture's stored
// schema, or a query's span.
type spanAnswer struct {
	// Error is why there is no answer: for a query the emulator refuses,
	// its gRPC code; for a fixture, what failed.
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`

	Schemas []spanSchema `json:"schemas,omitempty"`

	// Columns is the query's output column names, from its result set
	// metadata; an anonymous column is "".
	Columns []string `json:"columns,omitempty"`
	// Tables is the tables and views the query reads, as schema.name, or
	// name in the default schema.
	Tables []string `json:"tables,omitempty"`
}

type spanSchema struct {
	Name   string         `json:"name"`
	Tables []spanRelation `json:"tables,omitempty"`
	Views  []spanRelation `json:"views,omitempty"`
}

type spanRelation struct {
	Name       string       `json:"name"`
	Definition string       `json:"definition,omitempty"`
	Columns    []spanColumn `json:"columns"`
}

type spanColumn struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// spanFixture is the fixture the queries run on: its database, and for
// each table and view one database without it.
type spanFixture struct {
	full    *spanner.Client
	objects []string
	without map[string]*spanner.Client
	dbs     []string
	clients []*spanner.Client
}

// span answers a span request, or reports false for any other statement.
func (o *oracle) span(ctx context.Context, req string) (any, bool) {
	switch {
	case strings.HasPrefix(req, spanFixturePrefix):
		ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
		defer cancel()
		a, err := o.buildSpanFixture(ctx, strings.TrimPrefix(req, spanFixturePrefix))
		if err != nil {
			return spanAnswer{Error: "fixture", Message: oneLine(err.Error())}, true
		}
		return a, true
	case strings.HasPrefix(req, spanQueryPrefix):
		ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
		defer cancel()
		return o.querySpan(ctx, strings.TrimPrefix(req, spanQueryPrefix)), true
	}
	return nil, false
}

func (o *oracle) dropSpanFixture() {
	f := o.spanFx
	if f == nil {
		return
	}
	for _, c := range f.clients {
		c.Close()
	}
	for _, db := range f.dbs {
		_ = o.dbAdmin.DropDatabase(context.Background(), &databasepb.DropDatabaseRequest{Database: db})
	}
	o.spanFx = nil
}

func splitDDL(sql string) []string {
	var out []string
	for _, s := range regexp.MustCompile(`;[ \t]*(\n|$)`).Split(sql, -1) {
		if s = strings.TrimSpace(s); s != "" {
			out = append(out, s)
		}
	}
	return out
}

func (o *oracle) buildSpanFixture(ctx context.Context, sql string) (spanAnswer, error) {
	o.dropSpanFixture()
	f := &spanFixture{without: map[string]*spanner.Client{}}
	o.spanFx = f
	ddl := splitDDL(sql)
	full, err := o.spanDatabase(ctx, f, ddl)
	if err != nil {
		return spanAnswer{}, err
	}
	schemas, kinds, err := readSchema(ctx, full)
	if err != nil {
		return spanAnswer{}, err
	}
	for _, s := range schemas {
		for _, t := range s.Tables {
			f.objects = append(f.objects, qualify(s.Name, t.Name))
		}
		for _, v := range s.Views {
			f.objects = append(f.objects, qualify(s.Name, v.Name))
		}
	}
	for _, obj := range f.objects {
		c, err := o.spanDatabase(ctx, f, ddl)
		if err != nil {
			return spanAnswer{}, err
		}
		db := c.DatabaseName()
		if err := o.remove(ctx, db, obj, kinds, map[string]bool{}); err != nil {
			return spanAnswer{}, fmt.Errorf("remove %s: %w", obj, err)
		}
		f.without[obj] = c
	}
	// Queries run only on a fixture built whole.
	f.full = full
	return spanAnswer{Schemas: schemas}, nil
}

// spanDatabase makes a database holding the fixture.
func (o *oracle) spanDatabase(ctx context.Context, f *spanFixture, ddl []string) (*spanner.Client, error) {
	o.spanN++
	id := fmt.Sprintf("span_%d_%d", os.Getpid(), o.spanN)
	op, err := o.dbAdmin.CreateDatabase(ctx, &databasepb.CreateDatabaseRequest{
		Parent:          parent,
		CreateStatement: "CREATE DATABASE `" + id + "`",
		ExtraStatements: ddl,
	})
	if err != nil {
		return nil, err
	}
	path := parent + "/databases/" + id
	f.dbs = append(f.dbs, path)
	if _, err := op.Wait(ctx); err != nil {
		return nil, err
	}
	c, err := spanner.NewClient(ctx, path)
	if err != nil {
		return nil, err
	}
	f.clients = append(f.clients, c)
	return c, nil
}

// dependentView reads the view the emulator names when it refuses to
// rename a table a view reads, and dependentViews the views it lists when
// it refuses to drop a view another view reads.
var (
	dependentView  = regexp.MustCompile("definition of VIEW `([^`]+)`")
	dependentViews = regexp.MustCompile(`dependent views: ([^.]+(?:\.[^ ,.]+)?)`)
)

// dependent returns the first view a refusal names.
func dependent(msg string) string {
	if m := dependentView.FindStringSubmatch(msg); m != nil {
		return m[1]
	}
	if m := dependentViews.FindStringSubmatch(msg); m != nil {
		first, _, _ := strings.Cut(m[1], ",")
		return strings.TrimSpace(first)
	}
	return ""
}

// remove makes obj unreadable in db: a table is renamed away, a view
// dropped. A view that reads it, which the emulator names in its refusal,
// is dropped first.
func (o *oracle) remove(ctx context.Context, db, obj string, kinds map[string]string, gone map[string]bool) error {
	if gone[obj] {
		return nil
	}
	stmt := "DROP VIEW " + quote(obj)
	if kinds[obj] == "BASE TABLE" {
		hidden := obj + "_omni_hidden"
		stmt = "ALTER TABLE " + quote(obj) + " RENAME TO " + quote(hidden)
	}
	for range 100 {
		op, err := o.dbAdmin.UpdateDatabaseDdl(ctx, &databasepb.UpdateDatabaseDdlRequest{Database: db, Statements: []string{stmt}})
		if err == nil {
			err = op.Wait(ctx)
		}
		if err == nil {
			gone[obj] = true
			return nil
		}
		view := dependent(status.Convert(err).Message())
		if view == "" || view == obj || gone[view] {
			return fmt.Errorf("%s: %w", stmt, err)
		}
		if err := o.remove(ctx, db, view, kinds, gone); err != nil {
			return err
		}
	}
	return fmt.Errorf("%s: too many dependent views", stmt)
}

func qualify(schema, name string) string {
	if schema == "" {
		return name
	}
	return schema + "." + name
}

func quote(name string) string {
	parts := strings.Split(name, ".")
	for i, p := range parts {
		parts[i] = "`" + p + "`"
	}
	return strings.Join(parts, ".")
}

// readSchema reads the stored schema as Bytebase's sync reads it
// (spanner/driver/sync.go): tables, views with their definitions, and the
// columns of both in order. kinds holds each object's TABLE_TYPE.
func readSchema(ctx context.Context, c *spanner.Client) ([]spanSchema, map[string]string, error) {
	tx := c.ReadOnlyTransaction()
	defer tx.Close()
	rows := func(sql string, each func(r *spanner.Row) error) error {
		it := tx.Query(ctx, spanner.Statement{SQL: sql})
		defer it.Stop()
		for {
			r, err := it.Next()
			if errors.Is(err, iterator.Done) {
				return nil
			}
			if err != nil {
				return err
			}
			if err := each(r); err != nil {
				return err
			}
		}
	}
	const user = "TABLE_SCHEMA NOT IN ('INFORMATION_SCHEMA', 'SPANNER_SYS')"
	columns := map[string][]spanColumn{}
	if err := rows("SELECT TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, SPANNER_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE "+user+" ORDER BY TABLE_SCHEMA, TABLE_NAME, ORDINAL_POSITION", func(r *spanner.Row) error {
		var schema, table, name string
		var typ spanner.NullString
		if err := r.Columns(&schema, &table, &name, &typ); err != nil {
			return err
		}
		columns[qualify(schema, table)] = append(columns[qualify(schema, table)], spanColumn{Name: name, Type: typ.StringVal})
		return nil
	}); err != nil {
		return nil, nil, err
	}
	definitions := map[string]string{}
	if err := rows("SELECT TABLE_SCHEMA, TABLE_NAME, VIEW_DEFINITION FROM INFORMATION_SCHEMA.VIEWS WHERE "+user, func(r *spanner.Row) error {
		var schema, name string
		var def spanner.NullString
		if err := r.Columns(&schema, &name, &def); err != nil {
			return err
		}
		definitions[qualify(schema, name)] = def.StringVal
		return nil
	}); err != nil {
		return nil, nil, err
	}
	var out []spanSchema
	kinds := map[string]string{}
	if err := rows("SELECT TABLE_SCHEMA, TABLE_NAME, TABLE_TYPE FROM INFORMATION_SCHEMA.TABLES WHERE "+user+" ORDER BY TABLE_SCHEMA, TABLE_NAME", func(r *spanner.Row) error {
		var schema, name, kind string
		if err := r.Columns(&schema, &name, &kind); err != nil {
			return err
		}
		if len(out) == 0 || out[len(out)-1].Name != schema {
			out = append(out, spanSchema{Name: schema})
		}
		s := &out[len(out)-1]
		key := qualify(schema, name)
		kinds[key] = kind
		rel := spanRelation{Name: name, Columns: columns[key]}
		if kind == "VIEW" {
			rel.Definition = definitions[key]
			s.Views = append(s.Views, rel)
		} else {
			s.Tables = append(s.Tables, rel)
		}
		return nil
	}); err != nil {
		return nil, nil, err
	}
	return out, kinds, nil
}

func (o *oracle) querySpan(ctx context.Context, sql string) spanAnswer {
	f := o.spanFx
	if f == nil || f.full == nil {
		return spanAnswer{Error: "fixture", Message: "no fixture built"}
	}
	columns, err := runQuery(ctx, f.full, sql)
	if err != nil {
		return spanAnswer{Error: status.Code(err).String(), Message: oneLine(status.Convert(err).Message())}
	}
	a := spanAnswer{Columns: columns}
	for _, obj := range f.objects {
		_, err := runQuery(ctx, f.without[obj], sql)
		switch {
		case err == nil:
		case status.Code(err) == codes.InvalidArgument:
			a.Tables = append(a.Tables, obj)
		default:
			// Not name resolution: the harness cannot tell.
			return spanAnswer{Error: "infra", Message: oneLine(obj + ": " + err.Error())}
		}
	}
	return a
}

// runQuery runs sql and returns its output column names. The fixture holds
// no rows, so the query reads none.
func runQuery(ctx context.Context, c *spanner.Client, sql string) ([]string, error) {
	it := c.Single().Query(ctx, spanner.Statement{SQL: sql})
	defer it.Stop()
	_, err := it.Next()
	if err != nil && !errors.Is(err, iterator.Done) {
		return nil, err
	}
	if it.Metadata == nil || it.Metadata.RowType == nil {
		return nil, errors.New("no result set metadata")
	}
	names := []string{}
	for _, f := range it.Metadata.RowType.Fields {
		names = append(names, f.Name)
	}
	return names, nil
}
