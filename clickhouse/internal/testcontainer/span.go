package testcontainer

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/tape"
)

// Spans asks the pinned ClickHouse what a query reads, for clickhouse/analysis's
// span.tsv. Its tape holds one raw exchange per question, and each answer is
// what the server stated, reduced to JSON:
//
//   - Fixture builds a starting schema on a server emptied of every database
//     the image did not start with, and answers the schema ClickHouse stored:
//     each database's tables and views with their columns in order, and each
//     view's create_table_query, which is what a sync stores as its
//     definition.
//   - Span runs a query in a session's database and answers its result
//     set's column names, each column's lineage as the analyzer resolved it
//     (EXPLAIN QUERY TREE passes = 1, the tree after name resolution and
//     before any optimization), and the tables system.query_log says it
//     read; or the name of the error ClickHouse refused it with.
type Spans struct {
	raw *tape.Raw
}

// OpenSpans returns t's Spans: a replay of its tape, or the pinned ClickHouse
// in a recording.
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
		if err := start(context.Background(), pin); err != nil {
			t.Fatalf("ClickHouse: %v", err)
		}
		l := &spanLive{conns: map[string]*sql.DB{}}
		t.Cleanup(l.close)
		live = l.do
	}
	return &Spans{raw: tape.OpenRaw(t, &tape.Server{Pins: Pins}, Default, live)}
}

// SpanFixture is the schema ClickHouse stored for a fixture.
type SpanFixture struct {
	Databases []SpanDatabase `json:"databases"`
}

// SpanDatabase is one database a fixture made.
type SpanDatabase struct {
	Name   string         `json:"name"`
	Tables []SpanRelation `json:"tables"`
	Views  []SpanRelation `json:"views"`
}

// SpanRelation is a table or a view and its columns in order.
type SpanRelation struct {
	Name    string   `json:"name"`
	Columns []string `json:"columns"`
	// Definition is a view's create_table_query.
	Definition string `json:"definition,omitempty"`
}

// SpanAnswer is what ClickHouse states about one query.
type SpanAnswer struct {
	// Error is the name of the error code ClickHouse refused the query
	// with (UNKNOWN_IDENTIFIER, SYNTAX_ERROR, ...); the rest is empty.
	Error string `json:"error,omitempty"`
	// Columns is the result set's column names, in order.
	Columns []string `json:"columns,omitempty"`
	// Lineage is, for each result column, the columns its value is computed
	// from, as database.table.column, with Unstated ("~") where the tree
	// names a source that is no table (a table function's column).
	Lineage [][]string `json:"lineage,omitempty"`
	// Tables is every table and view query_log says the query read, as
	// database.table, in order.
	Tables []string `json:"tables,omitempty"`
}

// Fixture builds the starting schema sql, one statement per line, on a
// server emptied of every database it was not started with, and returns
// what ClickHouse stored.
func (s *Spans) Fixture(sql string) (SpanFixture, error) {
	var fx SpanFixture
	ans, err := s.raw.Do("fixture\n" + sql)
	if err != nil {
		return fx, err
	}
	return fx, json.Unmarshal([]byte(ans), &fx)
}

// Span runs statement in database and returns what ClickHouse states it
// reads.
func (s *Spans) Span(database, statement string) (SpanAnswer, error) {
	var a SpanAnswer
	ans, err := s.raw.Do("span " + database + "\n" + statement)
	if err != nil {
		return a, err
	}
	return a, json.Unmarshal([]byte(ans), &a)
}

type spanLive struct {
	conns map[string]*sql.DB
	// views maps database.view to its stored SELECT, for following a
	// column read through a view into its definition.
	views map[string]string
	// queries numbers the queries run, for their query_log rows.
	queries int
}

func (l *spanLive) close() {
	for _, db := range l.conns {
		_ = db.Close()
	}
}

func (l *spanLive) conn(database string) *sql.DB {
	db, ok := l.conns[database]
	if !ok {
		db = sql.OpenDB(connector(database))
		db.SetMaxOpenConns(1)
		l.conns[database] = db
	}
	return db
}

func (l *spanLive) do(req string) (string, error) {
	head, body, _ := strings.Cut(req, "\n")
	var (
		v   any
		err error
	)
	switch {
	case head == "fixture":
		v, err = l.fixture(context.Background(), body)
	case strings.HasPrefix(head, "span "):
		v, err = l.span(context.Background(), strings.TrimPrefix(head, "span "), body)
	default:
		return "", fmt.Errorf("spans: an unknown request %q", head)
	}
	if err != nil {
		return "", err
	}
	b, err := json.Marshal(v)
	return string(b), err
}

func (l *spanLive) fixture(ctx context.Context, sqlText string) (SpanFixture, error) {
	var fx SpanFixture
	l.close()
	l.conns = map[string]*sql.DB{}
	if err := empty(ctx, "default"); err != nil {
		return fx, err
	}
	root := l.conn("default")
	for _, stmt := range strings.Split(sqlText, "\n") {
		stmt = strings.TrimSuffix(strings.TrimSpace(stmt), ";")
		if stmt == "" || strings.HasPrefix(stmt, "--") {
			continue
		}
		if _, err := root.ExecContext(ctx, stmt); err != nil {
			return fx, fmt.Errorf("build the fixture: %s: %w", stmt, err)
		}
	}
	found, err := rows(ctx, root, `
		SELECT t.database, t.name, t.engine, t.create_table_query, t.as_select, c.name
		FROM system.tables AS t
		INNER JOIN system.columns AS c ON c.database = t.database AND c.table = t.name
		WHERE t.database NOT IN ('system', 'information_schema', 'INFORMATION_SCHEMA')
		ORDER BY t.database, t.name, c.position`)
	if err != nil {
		return fx, err
	}
	l.views = map[string]string{}
	for _, r := range found {
		database, name, engine, definition, asSelect, column := r[0], r[1], r[2], r[3], r[4], r[5]
		if n := len(fx.Databases); n == 0 || fx.Databases[n-1].Name != database {
			fx.Databases = append(fx.Databases, SpanDatabase{Name: database})
		}
		d := &fx.Databases[len(fx.Databases)-1]
		list, def := &d.Tables, ""
		if engine == "View" {
			list, def = &d.Views, definition
			l.views[database+"."+name] = asSelect
		}
		if n := len(*list); n == 0 || (*list)[n-1].Name != name {
			*list = append(*list, SpanRelation{Name: name, Definition: def})
		}
		rel := &(*list)[len(*list)-1]
		rel.Columns = append(rel.Columns, column)
	}
	return fx, nil
}

func (l *spanLive) span(ctx context.Context, database, statement string) (SpanAnswer, error) {
	var a SpanAnswer
	db := l.conn(database)
	l.queries++
	id := fmt.Sprintf("omni-span-%d", l.queries)
	result, err := db.QueryContext(clickhouse.Context(ctx, clickhouse.WithQueryID(id)), statement)
	if err == nil {
		a.Columns, err = result.Columns()
		for err == nil && result.Next() {
		}
		if err == nil {
			err = result.Err()
		}
		_ = result.Close()
	}
	var refused *clickhouse.Exception
	if errors.As(err, &refused) {
		// The exception's Name is its class (DB::Exception); the code's
		// name is what the server reports it as.
		var name string
		if err := db.QueryRowContext(ctx, "SELECT errorCodeToName(?)", refused.Code).Scan(&name); err != nil {
			return a, err
		}
		return SpanAnswer{Error: name}, nil
	}
	if err != nil {
		return a, err
	}
	tree, err := l.tree(ctx, db, statement)
	if err != nil {
		return a, err
	}
	for i := range a.Columns {
		sources, err := l.lineage(ctx, tree, i)
		if err != nil {
			return a, fmt.Errorf("lineage of column %d: %w", i, err)
		}
		a.Lineage = append(a.Lineage, sources)
	}
	if len(a.Lineage) != len(a.Columns) {
		return a, fmt.Errorf("the query tree has %d projection columns, the result %d", len(a.Lineage), len(a.Columns))
	}
	if _, err := db.ExecContext(ctx, "SYSTEM FLUSH LOGS"); err != nil {
		return a, err
	}
	var tables []string
	if err := db.QueryRowContext(ctx, "SELECT tables FROM system.query_log WHERE query_id = ? AND type = 'QueryFinish'", id).Scan(&tables); err != nil {
		return a, fmt.Errorf("query_log of %s: %w", id, err)
	}
	function := false
	for _, t := range tables {
		// A query without FROM reads system.one, and a recursive CTE's
		// working table is a temporary table with a name made up for the
		// query: neither is a table the query names. A table function is
		// listed as one of _table_function's, which does not say which
		// tables it reads (merge reads some): those are unstated.
		switch {
		case t == "system.one" || strings.HasPrefix(t, temporaryTables):
		case strings.HasPrefix(t, "_table_function."):
			if !function {
				a.Tables, function = append(a.Tables, "~"), true
			}
		default:
			a.Tables = append(a.Tables, t)
		}
	}
	return a, nil
}

// tree is the analyzed query tree of statement.
func (l *spanLive) tree(ctx context.Context, db *sql.DB, statement string) (*queryTree, error) {
	lines, err := rows(ctx, db, "EXPLAIN QUERY TREE passes = 1 "+statement)
	if err != nil {
		return nil, err
	}
	text := make([]string, len(lines))
	for i, r := range lines {
		text[i] = r[0]
	}
	return parseTree(text)
}

// lineage is the sources of the root's i-th projection column.
func (l *spanLive) lineage(ctx context.Context, tree *queryTree, i int) ([]string, error) {
	w := &treeWalker{live: l, ctx: ctx, tree: tree, seen: map[string]bool{}}
	if err := w.projection(tree.root, i); err != nil {
		return nil, err
	}
	slices.Sort(w.out)
	return slices.Compact(w.out), nil
}

// A treeNode is one line of EXPLAIN QUERY TREE and the lines indented under
// it: a node ("COLUMN id: 2, column_name: a, ...") or a section of its
// parent ("PROJECTION", "JOIN TREE", "EXPRESSION").
type treeNode struct {
	kind     string // COLUMN, QUERY, TABLE, ..., or the section's name
	id       int    // -1 for a section
	text     string // the line, without its indentation
	children []*treeNode
}

type queryTree struct {
	root *treeNode
	byID map[int]*treeNode
}

// temporaryTables is the database of temporary tables, which a recursive CTE
// runs on.
const temporaryTables = "_temporary_and_external_tables."

var nodeRe = regexp.MustCompile(`^([A-Z_]+(?: [A-Z_]+)*) id: (\d+)`)

func parseTree(lines []string) (*queryTree, error) {
	t := &queryTree{byID: map[int]*treeNode{}}
	type open struct {
		indent int
		node   *treeNode
	}
	var stack []open
	for _, line := range lines {
		trimmed := strings.TrimLeft(line, " ")
		indent := len(line) - len(trimmed)
		trimmed = strings.TrimRight(trimmed, " ")
		n := &treeNode{text: trimmed, id: -1, kind: trimmed}
		if m := nodeRe.FindStringSubmatch(trimmed); m != nil {
			n.kind = m[1]
			n.id, _ = strconv.Atoi(m[2])
			// A lambda's ARGUMENT lines number from 0 again; a column's
			// source is the first node with its id.
			if _, ok := t.byID[n.id]; !ok {
				t.byID[n.id] = n
			}
		}
		for len(stack) > 0 && stack[len(stack)-1].indent >= indent {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			if t.root != nil {
				return nil, fmt.Errorf("query tree: a second root %q", trimmed)
			}
			t.root = n
		} else {
			p := stack[len(stack)-1].node
			p.children = append(p.children, n)
		}
		stack = append(stack, open{indent, n})
	}
	if t.root == nil {
		return nil, errors.New("query tree: empty")
	}
	return t, nil
}

func (n *treeNode) section(name string) *treeNode {
	for _, c := range n.children {
		if c.id < 0 && c.kind == name {
			return c
		}
	}
	return nil
}

// field returns the value of key in a node's line: what follows "key: " up
// to the next ", " (or the end).
func (n *treeNode) field(key string) string {
	_, rest, ok := strings.Cut(n.text, ", "+key+": ")
	if !ok {
		return ""
	}
	if key == "column_name" {
		// A name may hold ", "; the result type follows it.
		if v, _, ok := strings.Cut(rest, ", result_type: "); ok {
			return v
		}
	}
	v, _, _ := strings.Cut(rest, ", ")
	return v
}

// list is the items of a section's LIST node.
func (n *treeNode) list() []*treeNode {
	if n == nil {
		return nil
	}
	for _, c := range n.children {
		if c.kind == "LIST" {
			return c.children
		}
	}
	return nil
}

// projectionNames is a QUERY node's PROJECTION COLUMNS: each line is the
// name and the type, which holds no space outside its parentheses.
func (n *treeNode) projectionNames() []string {
	sec := n.section("PROJECTION COLUMNS")
	if sec == nil {
		return nil
	}
	var out []string
	for _, c := range sec.children {
		out = append(out, splitNameType(c.text))
	}
	return out
}

func splitNameType(line string) string {
	depth, quoted, cut := 0, false, -1
	for i := 0; i < len(line); i++ {
		switch c := line[i]; {
		case quoted && c == '\\':
			i++
		case c == '\'' || c == '`':
			quoted = !quoted
		case quoted:
		case c == '(':
			depth++
		case c == ')':
			depth--
		case c == ' ' && depth == 0:
			cut = i
		}
	}
	if cut < 0 {
		return line
	}
	return line[:cut]
}

// treeWalker collects the sources of one projection column.
type treeWalker struct {
	live *spanLive
	ctx  context.Context
	tree *queryTree
	out  []string
	// seen holds the projections already walked, so a CTE read twice is
	// walked once.
	seen map[string]bool
}

// projection walks the i-th projection column of a QUERY or UNION node.
func (w *treeWalker) projection(n *treeNode, i int) error {
	key := fmt.Sprintf("%p/%d", n, i)
	if w.seen[key] {
		return nil
	}
	w.seen[key] = true
	switch n.kind {
	case "QUERY":
		items := n.section("PROJECTION").list()
		if i >= len(items) {
			return fmt.Errorf("query %d has %d projection columns, not %d", n.id, len(items), i+1)
		}
		return w.expr(items[i])
	case "UNION":
		for _, arm := range n.section("QUERIES").list() {
			if err := w.projection(arm, i); err != nil {
				return err
			}
		}
		return nil
	}
	return fmt.Errorf("a projection of %s", n.kind)
}

// named walks the projection column of a QUERY or UNION node that a column
// of the enclosing query names.
func (w *treeWalker) named(n *treeNode, name string) error {
	first := n
	for first.kind == "UNION" {
		arms := first.section("QUERIES").list()
		if len(arms) == 0 {
			return fmt.Errorf("union %d has no queries", n.id)
		}
		first = arms[0]
	}
	i := slices.Index(first.projectionNames(), name)
	if i < 0 {
		return fmt.Errorf("query %d has no column %q", n.id, name)
	}
	return w.projection(n, i)
}

func (w *treeWalker) expr(n *treeNode) error {
	switch n.kind {
	case "COLUMN":
		return w.column(n)
	case "QUERY", "UNION":
		// A subquery in an expression yields its own projection.
		names := n
		for names.kind == "UNION" {
			names = names.section("QUERIES").list()[0]
		}
		for i := range names.projectionNames() {
			if err := w.projection(n, i); err != nil {
				return err
			}
		}
		return nil
	}
	for _, c := range n.children {
		if err := w.expr(c); err != nil {
			return err
		}
	}
	return nil
}

func (w *treeWalker) column(n *treeNode) error {
	if e := n.section("EXPRESSION"); e != nil {
		// A column the analyzer computed: a USING column (each side's), an
		// ARRAY JOIN alias, or an ALIAS column.
		return w.expr(e)
	}
	src, err := strconv.Atoi(n.field("source_id"))
	if err != nil {
		return fmt.Errorf("column %d has no source: %q", n.id, n.text)
	}
	source := w.tree.byID[src]
	if source == nil {
		return fmt.Errorf("column %d's source %d is not in the tree", n.id, src)
	}
	name := n.field("column_name")
	switch source.kind {
	case "TABLE":
		table := source.field("table_name")
		if strings.HasPrefix(table, temporaryTables) {
			// A recursive CTE's working table, named afresh for each
			// query: what it holds is unstated here, and the CTE's arms
			// name what they read.
			w.out = append(w.out, "~")
			return nil
		}
		w.out = append(w.out, table+"."+name)
		if def, ok := w.live.views[table]; ok {
			return w.view(table, def, name)
		}
		return nil
	case "QUERY", "UNION":
		return w.named(source, name)
	case "ARRAY_JOIN":
		for _, c := range source.section("JOIN EXPRESSIONS").list() {
			if c.kind == "COLUMN" && (c.field("alias") == name || c.field("column_name") == name) {
				return w.expr(c)
			}
		}
		return fmt.Errorf("array join %d has no column %q", src, name)
	case "ARGUMENTS":
		// A lambda's parameter: its values come from the array arguments,
		// which the enclosing function names.
		return nil
	}
	// A table function's column, or another source the tree names no table
	// for.
	w.out = append(w.out, "~")
	return nil
}

// view walks column name of view (database.view) through the definition
// ClickHouse stored for it.
func (w *treeWalker) view(view, definition, name string) error {
	database, _, _ := strings.Cut(view, ".")
	tree, err := w.live.tree(w.ctx, w.live.conn(database), definition)
	if err != nil {
		return fmt.Errorf("view %s: %w", view, err)
	}
	sub := &treeWalker{live: w.live, ctx: w.ctx, tree: tree, seen: map[string]bool{}}
	if err := sub.named(tree.root, name); err != nil {
		return fmt.Errorf("view %s: %w", view, err)
	}
	w.out = append(w.out, sub.out...)
	return nil
}
