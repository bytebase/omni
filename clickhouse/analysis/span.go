package analysis

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/clickhouse/parser"
	"github.com/bytebase/omni/metadata"
)

// ResolvedSpan is a query's span: the tables it reads and each result
// column's sources.
type ResolvedSpan = analysis.ResolvedSpan

// QuerySpanResult is one result column of a query and the columns it reads.
type QuerySpanResult = analysis.QuerySpanResult

// ColumnResource is a column, or a table when Column is empty. ClickHouse
// has no schemas: Schema is always empty.
type ColumnResource = analysis.ColumnResource

// SourceColumnSet is a set of columns or tables.
type SourceColumnSet = analysis.ColumnSet

// NotFoundError reports a database, table or column the catalog does not
// have. A nil field is not part of the report.
type NotFoundError = analysis.ResourceNotFoundError

// ResolveOptions configures ResolveQuerySpan.
type ResolveOptions struct {
	// DefaultDatabase is where an unqualified table name resolves: the
	// session's current database.
	DefaultDatabase string
	// Catalog is the synced metadata the query's names resolve against.
	Catalog Catalog
}

// Catalog is the synced metadata of the instance a statement runs on.
type Catalog struct {
	// GetDatabase returns a database's metadata, or nil when the database
	// has not been synced. A nil GetDatabase means no metadata is
	// available: every table's columns are unknown, a column resolves to
	// each table it could belong to, and nothing is reported missing.
	GetDatabase func(name string) (Database, error)
}

// Database is one database's synced metadata. A sync stores a ClickHouse
// database as one schema, so its tables and views are the database's.
type Database interface {
	// GetTable returns the named table, or nil.
	GetTable(name string) *metadata.TableMetadata
	// GetView returns the named view, or nil. Its Definition is the view's
	// CREATE VIEW statement (system.tables' create_table_query).
	GetView(name string) *metadata.ViewMetadata
}

// ErrAmbiguousColumn reports an unqualified column more than one joined
// table has, in a query with more than one JOIN: ClickHouse prefers the
// left table only when the query joins once (single_join_prefer_left_table)
// and refuses the name otherwise (AMBIGUOUS_IDENTIFIER).
var ErrAmbiguousColumn = errors.New("ambiguous column")

// systemDatabases hold the server's own tables, which a sync does not
// store: a table there is read with columns the metadata does not list.
var systemDatabases = map[string]bool{
	"system":             true,
	"information_schema": true,
	"INFORMATION_SCHEMA": true,
}

// ResolveQuerySpan resolves a single ClickHouse query: the tables and views it
// reads (a view's own tables included), and for each result column its name
// as ClickHouse names it and the table columns its value is computed from. A
// column read through a view names the view's column and the columns the
// view's definition reads for it.
//
// Names resolve as ClickHouse's analyzer resolves them: a lambda's
// parameter first, then an alias the query defines anywhere (WITH, the
// select list, or inside an expression), then a column of the FROM clause,
// the leftmost table holding it first. A name the catalog does not have is
// reported in NotFoundError, and Results is then empty.
//
// A statement that is not a query returns its Type and no span. A syntax
// error is returned as *parser.ParseError.
func ResolveQuerySpan(statement string, opts ResolveOptions) (*ResolvedSpan, error) {
	seg, ok := singleNonEmptySegment(statement)
	if !ok {
		return nil, errors.New("expected a single ClickHouse statement")
	}
	stmts, err := parser.Parse(seg.Text)
	if err != nil {
		return nil, offsetParseError(err, seg.ByteStart)
	}
	if len(stmts) != 1 {
		return nil, errors.New("expected a single ClickHouse statement")
	}
	out := &ResolvedSpan{Type: classifyParsedStatement(stmts[0]), SourceColumns: SourceColumnSet{}}
	query, ok := stmts[0].(*parser.SelectUnionStmt)
	if !ok {
		return out, nil
	}
	r := newSpanResolver(opts, seg.Text)
	cols := r.query(query, nil, nil)
	if r.err != nil {
		return nil, r.err
	}
	if r.ambiguous != "" {
		return nil, fmt.Errorf("%w %q", ErrAmbiguousColumn, r.ambiguous)
	}
	out.SourceColumns = r.tables
	if r.notFound != nil {
		out.NotFoundError = r.notFound
		return out, nil
	}
	for _, c := range cols {
		out.Results = append(out.Results, QuerySpanResult{
			Name:           c.name,
			SourceColumns:  c.sources,
			IsPlainField:   c.plain,
			SelectAsterisk: c.star,
		})
	}
	return out, nil
}

// spanColumn is one column a relation or a query produces.
type spanColumn struct {
	name    string
	sources SourceColumnSet
	plain   bool
	// star marks a * no metadata could expand.
	star bool
	// hidden marks the right side's USING column, which * skips and only a
	// qualified name reads.
	hidden bool
	// rel is the FROM relation the column was read from, if any.
	rel *relation
}

// relation is a table, view, CTE, subquery or table function in a FROM
// clause.
type relation struct {
	// database and table name a table or view; both are empty for a
	// derived relation.
	database, table string
	// name is what a qualifier matches: the alias, or the table's name.
	name    string
	columns []spanColumn
	// open marks a relation whose columns no metadata lists (a table of a
	// system database, a table function, any table without a catalog): any
	// name resolves on it, to the column of that name of the table, or to
	// nothing for a table function.
	open bool
	// arrayJoin marks the columns an ARRAY JOIN adds, which shadow the FROM
	// clause's.
	arrayJoin bool
}

func (rel *relation) column(name string) (spanColumn, bool) {
	for _, c := range rel.columns {
		if c.name == name {
			c.rel = rel
			return c, true
		}
	}
	if rel.open {
		c := spanColumn{name: name, sources: SourceColumnSet{}, plain: true, rel: rel}
		if rel.table != "" {
			c.sources[ColumnResource{Database: rel.database, Table: rel.table, Column: name}] = true
		}
		return c, true
	}
	return spanColumn{}, false
}

// spanScope is what names resolve against in one SELECT, inside its
// enclosing queries' scopes.
type spanScope struct {
	relations []*relation
	// aliases are the expression aliases the SELECT defines, in its WITH
	// clause and anywhere in its select list.
	aliases map[string]parser.Expr
	// lambda holds a lambda's parameters, which shadow every other name.
	lambda map[string]bool
	// windows are the SELECT's named windows (WINDOW w AS (...)).
	windows map[string]*parser.WindowSpec
	// joins counts the FROM clause's joins.
	joins  int
	parent *spanScope
}

// cteBinding is one CTE in scope.
type cteBinding struct {
	name   string
	query  *parser.SelectUnionStmt
	scope  *spanScope
	ctes   *cteBinding
	parent *cteBinding
	// columns are the CTE's columns once resolved.
	columns  []spanColumn
	resolved bool
	// resolving marks a recursive CTE whose query is being resolved.
	resolving bool
}

type spanResolver struct {
	opts ResolveOptions
	// database is where an unqualified table name resolves: the session's,
	// or a view's while its definition resolves.
	database string
	// text is the SQL being resolved, which names an expression no rule
	// names.
	text     string
	tables   SourceColumnSet
	notFound *NotFoundError
	// ambiguous is the first name more than one joined table has.
	ambiguous string
	// subqueries counts the unnamed subqueries the select lists name
	// _subquery_N.
	subqueries int
	// err is why a view's definition could not be read: the span fails
	// rather than omit what the view reads.
	err   error
	views map[ColumnResource][]spanColumn
	// resolving holds the views and aliases being resolved, so one that
	// reads itself stops.
	resolving map[any]bool
}

func newSpanResolver(opts ResolveOptions, text string) *spanResolver {
	return &spanResolver{
		opts:      opts,
		database:  opts.DefaultDatabase,
		text:      text,
		tables:    SourceColumnSet{},
		views:     map[ColumnResource][]spanColumn{},
		resolving: map[any]bool{},
	}
}

func (r *spanResolver) missing(nf NotFoundError) {
	if r.notFound == nil {
		r.notFound = &nf
	}
}

// query resolves a query's columns: a set operation's columns are named by
// its first query and read what every query reads in that position.
func (r *spanResolver) query(q *parser.SelectUnionStmt, outer *spanScope, ctes *cteBinding) []spanColumn {
	var out []spanColumn
	for i, item := range q.Items {
		cols := r.selectStmt(item, outer, ctes)
		if i == 0 {
			out = cols
			continue
		}
		for j := range out {
			if j < len(cols) {
				out[j].sources = union(out[j].sources, cols[j].sources)
				out[j].plain = false
			}
		}
	}
	return out
}

func (r *spanResolver) selectStmt(s *parser.SelectStmt, outer *spanScope, ctes *cteBinding) []spanColumn {
	if s.Nested != nil {
		return r.query(s.Nested, outer, ctes)
	}
	sc := &spanScope{aliases: map[string]parser.Expr{}, windows: map[string]*parser.WindowSpec{}, parent: outer}
	for _, w := range s.Windows {
		sc.windows[w.Name] = w.Spec
	}
	for _, w := range s.With {
		switch {
		case w.Subquery != nil:
			ctes = &cteBinding{name: w.Alias, query: w.Subquery, scope: outer, parent: ctes}
			ctes.ctes = ctes.parent
			if w.Recursive {
				ctes.ctes = ctes
			}
		case w.Alias != "":
			sc.aliases[w.Alias] = w.Expr
		}
	}
	for _, c := range s.Columns {
		collectAliases(c, sc.aliases)
	}
	if s.From != nil {
		r.from(s.From, sc, outer, ctes)
	}
	if s.ArrayJoin != nil {
		// An ARRAY JOIN column is an element of the array it names: an
		// alias names a new column, and a bare column is the array's own
		// name, read again element by element.
		rel := &relation{arrayJoin: true}
		for _, e := range s.ArrayJoin.Exprs {
			name, inner := "", e
			if a, ok := e.(*parser.AliasExpr); ok {
				name, inner = a.Alias, a.Expr
			} else if id, ok := e.(*parser.IdentExpr); ok {
				name = id.Parts[len(id.Parts)-1]
			}
			sources := r.expr(inner, sc, ctes)
			if name != "" {
				rel.columns = append(rel.columns, spanColumn{name: name, sources: sources})
			}
		}
		sc.relations = append([]*relation{rel}, sc.relations...)
	}

	var out []spanColumn
	for _, item := range s.Columns {
		if star, ok := item.(*parser.StarExpr); ok {
			out = append(out, r.star(star, sc, ctes)...)
			continue
		}
		if re, ok := columnsMatcher(item); ok {
			out = append(out, r.matching(re, sc)...)
			continue
		}
		c := spanColumn{
			name:    r.projectionName(item),
			sources: r.expr(item, sc, ctes),
			plain:   isPlainColumn(item),
		}
		if parts, ok := identParts(item); ok && len(parts) > 1 {
			// A qualified column keeps its qualifier unless its bare name
			// reads the same column.
			if found, ok := lookup(parts, sc); ok == found1 {
				if bare, ok := lookup(parts[len(parts)-1:], sc); ok != found1 || bare.rel != found.rel {
					c.name = strings.Join(parts, ".")
				}
			}
		}
		out = append(out, c)
	}

	for _, e := range []parser.Expr{s.Prewhere, s.Where, s.Having, s.Qualify} {
		r.expr(e, sc, ctes)
	}
	r.exprs(s.DistinctOn, sc, ctes)
	if s.GroupBy != nil {
		r.exprs(s.GroupBy.Exprs, sc, ctes)
	}
	for _, w := range s.Windows {
		r.window(w.Spec, sc, ctes)
	}
	r.orderBy(s.OrderBy, sc, ctes)
	r.exprs(s.Interpolate, sc, ctes)
	if s.LimitBy != nil {
		r.exprs(s.LimitBy.By, sc, ctes)
	}
	return out
}

func (r *spanResolver) exprs(list []parser.Expr, sc *spanScope, ctes *cteBinding) SourceColumnSet {
	out := SourceColumnSet{}
	for _, e := range list {
		out = union(out, r.expr(e, sc, ctes))
	}
	return out
}

func (r *spanResolver) orderBy(items []*parser.OrderByItem, sc *spanScope, ctes *cteBinding) SourceColumnSet {
	out := SourceColumnSet{}
	for _, o := range items {
		out = union(out, r.exprs([]parser.Expr{o.Expr, o.FillFrom, o.FillTo, o.FillStep}, sc, ctes))
	}
	return out
}

func (r *spanResolver) window(w *parser.WindowSpec, sc *spanScope, ctes *cteBinding) SourceColumnSet {
	if w == nil {
		return nil
	}
	return union(r.exprs(w.PartitionBy, sc, ctes), r.orderBy(w.OrderBy, sc, ctes))
}

// collectAliases records every alias an expression defines outside its
// subqueries and lambdas, which ClickHouse makes visible to the whole
// SELECT.
func collectAliases(e parser.Expr, into map[string]parser.Expr) {
	walkExpr(e, func(n parser.Expr) bool {
		switch n := n.(type) {
		case *parser.AliasExpr:
			if _, ok := into[n.Alias]; !ok {
				into[n.Alias] = n.Expr
			}
		case *parser.SubqueryExpr, *parser.LambdaExpr:
			return false
		}
		return true
	})
}

// from adds the relations of a FROM clause to sc, left to right. ON and
// USING resolve against the relations joined so far.
func (r *spanResolver) from(j *parser.JoinExpr, sc, outer *spanScope, ctes *cteBinding) {
	if j == nil {
		return
	}
	if j.Table != nil {
		if rel := r.tableExpr(j.Table, outer, ctes); rel != nil {
			sc.relations = append(sc.relations, rel)
		}
		return
	}
	sc.joins++
	r.from(j.Left, sc, outer, ctes)
	left := len(sc.relations)
	r.from(j.Right, sc, outer, ctes)
	r.exprs(j.On, sc, ctes)
	for _, u := range j.Using {
		name := usingName(u)
		if name == "" {
			r.expr(u, sc, ctes)
			continue
		}
		// A USING column is one column, read from both sides; the right
		// side's column of that name is hidden from *.
		sources := SourceColumnSet{}
		found := false
		for i, rel := range sc.relations {
			c, ok := rel.column(name)
			if !ok {
				continue
			}
			found = true
			sources = union(sources, c.sources)
			if i >= left {
				for k := range rel.columns {
					if rel.columns[k].name == name {
						rel.columns[k].hidden = true
					}
				}
			}
		}
		if !found {
			r.missing(NotFoundError{Column: &name})
			continue
		}
		for _, rel := range sc.relations[:left] {
			for k, c := range rel.columns {
				if c.name == name {
					rel.columns[k].sources = sources
				}
			}
		}
	}
}

func usingName(e parser.Expr) string {
	if a, ok := e.(*parser.AliasExpr); ok {
		e = a.Expr
	}
	if id, ok := e.(*parser.IdentExpr); ok && len(id.Parts) == 1 {
		return id.Parts[0]
	}
	return ""
}

// tableExpr resolves one FROM item.
func (r *spanResolver) tableExpr(t *parser.TableExpr, outer *spanScope, ctes *cteBinding) *relation {
	switch {
	case t.Subquery != nil:
		cols := r.query(t.Subquery, outer, ctes)
		return derived(t.Alias, renamed(cols, t.ColumnAliases))
	case t.Func != nil:
		for _, a := range t.Func.Args {
			if sub, ok := a.(*parser.SubqueryExpr); ok {
				r.query(sub.Query, outer, ctes)
			}
		}
		name := t.Alias
		if name == "" {
			name = t.Func.Name
		}
		if col, ok := tableFunctionColumns[strings.ToLower(t.Func.Name)]; ok {
			return &relation{name: name, columns: []spanColumn{{name: col, sources: SourceColumnSet{}, plain: true}}}
		}
		return &relation{name: name, open: true}
	case t.Table == "":
		// A query parameter names the table: nothing to resolve.
		return &relation{name: t.Alias, open: true}
	}
	if t.Database == "" {
		if cte := ctes.find(t.Table); cte != nil {
			name := t.Alias
			if name == "" {
				name = t.Table
			}
			if cte.resolving {
				// A recursive CTE's own rows: what they read is what the
				// CTE's other arms read.
				return &relation{name: name, open: true}
			}
			return derived(name, renamed(r.cte(cte), t.ColumnAliases))
		}
	}
	database := t.Database
	if database == "" {
		database = r.database
	}
	cols, open, ok := r.relation(database, t.Table)
	if !ok {
		return nil
	}
	name := t.Alias
	if name == "" {
		name = t.Table
	}
	return &relation{database: database, table: t.Table, name: name, columns: renamed(cols, t.ColumnAliases), open: open}
}

// tableFunctionColumns are the one column of the table functions that read
// no table and make their own rows.
var tableFunctionColumns = map[string]string{
	"numbers": "number", "numbers_mt": "number", "zeros": "zero", "zeros_mt": "zero",
}

func (c *cteBinding) find(name string) *cteBinding {
	for ; c != nil; c = c.parent {
		if c.name == name {
			return c
		}
	}
	return nil
}

// cte resolves a CTE's query once, where the CTE is defined.
func (r *spanResolver) cte(c *cteBinding) []spanColumn {
	if !c.resolved {
		c.resolved, c.resolving = true, true
		c.columns = r.query(c.query, c.scope, c.ctes)
		c.resolving = false
	}
	return c.columns
}

func derived(name string, cols []spanColumn) *relation {
	return &relation{name: name, columns: cols}
}

// renamed is cols with the names of a column alias list, AS t(a, b).
func renamed(cols []spanColumn, names []string) []spanColumn {
	out := make([]spanColumn, len(cols))
	copy(out, cols)
	for i := range out {
		if i < len(names) {
			out[i].name = names[i]
		}
	}
	return out
}

// relation reads a table or view of the catalog and records it. ok is false
// when the catalog does not have it; open is true when its columns are not
// known.
func (r *spanResolver) relation(database, table string) (cols []spanColumn, open, ok bool) {
	r.tables[ColumnResource{Database: database, Table: table}] = true
	if systemDatabases[database] || r.opts.Catalog.GetDatabase == nil {
		return nil, true, true
	}
	db, err := r.opts.Catalog.GetDatabase(database)
	if err != nil {
		r.err = errors.Join(r.err, err)
		return nil, false, false
	}
	if db == nil {
		r.missing(NotFoundError{Database: &database})
		return nil, false, false
	}
	if t := db.GetTable(table); t != nil {
		for _, c := range t.Columns {
			cols = append(cols, spanColumn{
				name:    c.Name,
				sources: SourceColumnSet{{Database: database, Table: table, Column: c.Name}: true},
				plain:   true,
			})
		}
		return cols, false, true
	}
	if v := db.GetView(table); v != nil {
		return r.view(database, v), false, true
	}
	r.missing(NotFoundError{Database: &database, Table: &table})
	return nil, false, false
}

// view resolves a view's stored definition in the view's database, once: its
// columns read the view's column and what the definition reads for it.
func (r *spanResolver) view(database string, v *metadata.ViewMetadata) []spanColumn {
	key := ColumnResource{Database: database, Table: v.Name}
	if cols, ok := r.views[key]; ok {
		return cols
	}
	var inner []spanColumn
	if !r.resolving[key] {
		r.resolving[key] = true
		q, err := viewQuery(v.Definition)
		if err != nil {
			r.err = errors.Join(r.err, fmt.Errorf("view %s.%s: %w", database, v.Name, err))
			return nil
		}
		savedDB, savedText := r.database, r.text
		r.database, r.text = database, v.Definition
		inner = r.query(q, nil, nil)
		r.database, r.text = savedDB, savedText
		delete(r.resolving, key)
	}
	var cols []spanColumn
	for i, c := range v.Columns {
		sources := SourceColumnSet{{Database: database, Table: v.Name, Column: c.Name}: true}
		if i < len(inner) {
			sources = union(sources, inner[i].sources)
		}
		cols = append(cols, spanColumn{name: c.Name, sources: sources, plain: true})
	}
	r.views[key] = cols
	return cols
}

// viewQuery is the SELECT of a view's CREATE VIEW statement.
func viewQuery(definition string) (*parser.SelectUnionStmt, error) {
	stmts, err := parser.Parse(definition)
	if err != nil {
		return nil, err
	}
	if len(stmts) == 1 {
		if cv, ok := stmts[0].(*parser.CreateViewStmt); ok && cv.Subquery != nil {
			return cv.Subquery, nil
		}
		if q, ok := stmts[0].(*parser.SelectUnionStmt); ok {
			return q, nil
		}
	}
	return nil, errors.New("the definition is not a CREATE VIEW ... AS SELECT")
}

// star expands *, t.* and their EXCEPT and REPLACE transformers.
func (r *spanResolver) star(s *parser.StarExpr, sc *spanScope, ctes *cteBinding) []spanColumn {
	var out []spanColumn
	for _, rel := range sc.relations {
		if len(s.Qualifier) > 0 && !qualifies(rel, s.Qualifier) {
			continue
		}
		if rel.open && len(rel.columns) == 0 {
			sources := SourceColumnSet{}
			if rel.table != "" {
				sources[ColumnResource{Database: rel.database, Table: rel.table}] = true
			}
			out = append(out, spanColumn{name: "*", sources: sources, star: true})
			continue
		}
		for _, c := range rel.columns {
			if !c.hidden {
				c.rel = rel
				c.name = starName(c, sc)
				out = append(out, c)
			}
		}
	}
	if len(s.Qualifier) > 0 && out == nil {
		q := strings.Join(s.Qualifier, ".")
		r.missing(NotFoundError{Table: &q})
	}
	except := map[string]bool{}
	for _, c := range s.Except {
		except[c] = true
	}
	for _, t := range s.Transformers {
		if t.Kind == "EXCEPT" {
			for _, c := range t.Columns {
				except[c] = true
			}
		}
	}
	var kept []spanColumn
	for _, c := range out {
		if !except[c.name] {
			kept = append(kept, c)
		}
	}
	for _, t := range s.Transformers {
		switch t.Kind {
		case "REPLACE":
			for _, rep := range t.Replaces {
				for i := range kept {
					if kept[i].name == rep.Alias {
						kept[i].sources = r.expr(rep.Expr, sc, ctes)
						kept[i].plain = false
					}
				}
			}
		case "APPLY":
			r.expr(t.ApplyExpr, sc, ctes)
			if t.ApplyFunc == "" {
				continue
			}
			for i := range kept {
				kept[i].name = t.ApplyFunc + "(" + kept[i].name + ")"
				kept[i].plain = false
			}
		}
	}
	return kept
}

func qualifies(rel *relation, q []string) bool {
	switch len(q) {
	case 1:
		return rel.name == q[0]
	case 2:
		return rel.database == q[0] && rel.table == q[1]
	}
	return false
}

// starName is the name of a column a * expands to: its own, unless that
// name does not read it (another joined relation's column has it), when it
// is relation.column, as ClickHouse names it.
func starName(c spanColumn, sc *spanScope) string {
	if bare, ok := lookup([]string{c.name}, sc); ok == found1 && bare.rel == c.rel {
		return c.name
	}
	return c.rel.name + "." + c.name
}

// columnsMatcher is the regular expression of a COLUMNS('re') item.
func columnsMatcher(e parser.Expr) (*regexp.Regexp, bool) {
	f, ok := e.(*parser.FuncCallExpr)
	if !ok || !strings.EqualFold(f.Name, "COLUMNS") || len(f.Args) != 1 {
		return nil, false
	}
	lit, ok := f.Args[0].(*parser.LiteralExpr)
	if !ok || lit.Kind != "STRING" {
		return nil, false
	}
	re, err := regexp.Compile(lit.Val)
	return re, err == nil
}

// matching expands COLUMNS('re'): every column whose name matches.
func (r *spanResolver) matching(re *regexp.Regexp, sc *spanScope) []spanColumn {
	var out []spanColumn
	for _, rel := range sc.relations {
		for _, c := range rel.columns {
			if !c.hidden && re.MatchString(c.name) {
				c.rel = rel
				c.name = starName(c, sc)
				out = append(out, c)
			}
		}
	}
	return out
}

// expr resolves an expression to the columns it reads, recording the tables
// its subqueries read.
func (r *spanResolver) expr(e parser.Expr, sc *spanScope, ctes *cteBinding) SourceColumnSet {
	out := SourceColumnSet{}
	if e == nil {
		return out
	}
	switch n := e.(type) {
	case *parser.IdentExpr:
		return r.ident(n.Parts, sc, ctes)
	case *parser.TupleAccessExpr:
		if parts, ok := identParts(n); ok {
			return r.ident(parts, sc, ctes)
		}
		return r.expr(n.Expr, sc, ctes)
	case *parser.SubqueryExpr:
		for _, c := range r.query(n.Query, sc, ctes) {
			out = union(out, c.sources)
		}
		return out
	case *parser.LambdaExpr:
		inner := &spanScope{lambda: map[string]bool{}, parent: sc}
		for _, p := range n.Params {
			inner.lambda[p] = true
		}
		return r.expr(n.Body, inner, ctes)
	case *parser.StarExpr:
		// count(*) and the like read no column.
		return out
	case *parser.FuncCallExpr:
		out = union(r.exprs(n.Params, sc, ctes), r.exprs(n.Args, sc, ctes))
		over := n.Over
		if n.OverName != "" {
			over = sc.window(n.OverName)
		}
		for ; over != nil; over = sc.window(over.BaseWindow) {
			out = union(out, r.window(over, sc, ctes))
		}
		return out
	}
	for _, c := range children(e) {
		out = union(out, r.expr(c, sc, ctes))
	}
	return out
}

// window is a named window of the innermost SELECT that defines it.
func (sc *spanScope) window(name string) *parser.WindowSpec {
	if name == "" {
		return nil
	}
	for s := sc; s != nil; s = s.parent {
		if w, ok := s.windows[name]; ok {
			return w
		}
	}
	return nil
}

// identParts flattens a.b.`c` (parsed as tuple access on an identifier)
// into the identifier's parts.
func identParts(e parser.Expr) ([]string, bool) {
	switch n := e.(type) {
	case *parser.IdentExpr:
		return n.Parts, true
	case *parser.TupleAccessExpr:
		if n.Member == "" || n.Member[0] >= '0' && n.Member[0] <= '9' {
			return nil, false
		}
		parts, ok := identParts(n.Expr)
		if !ok {
			return nil, false
		}
		return append(append([]string(nil), parts...), n.Member), true
	}
	return nil, false
}

// ident resolves a name: a lambda parameter reads nothing; an alias reads
// what its expression reads; a column, qualified or not, reads its sources.
func (r *spanResolver) ident(parts []string, sc *spanScope, ctes *cteBinding) SourceColumnSet {
	for s := sc; s != nil; s = s.parent {
		if s.lambda[parts[0]] {
			return SourceColumnSet{}
		}
	}
	if len(parts) == 1 {
		for s := sc; s != nil; s = s.parent {
			if e, ok := s.aliases[parts[0]]; ok && !r.resolving[e] {
				r.resolving[e] = true
				out := r.expr(e, s, ctes)
				delete(r.resolving, e)
				return out
			}
		}
	}
	name := strings.Join(parts, ".")
	for s := sc; s != nil; s = s.parent {
		switch c, found := lookup(parts, s); found {
		case found1:
			return c.sources
		case foundMany:
			if r.ambiguous == "" {
				r.ambiguous = name
			}
			return c.sources
		}
	}
	r.missing(NotFoundError{Column: &name})
	return SourceColumnSet{}
}

// How many relations of a scope a name finds.
type found int

const (
	foundNone found = iota
	found1
	// foundMany is an unqualified name more than one joined relation has,
	// which ClickHouse refuses in a query that joins more than once.
	foundMany
)

// lookup finds a column of one scope's relations: an unqualified name in
// the ARRAY JOIN's columns, else in the relation that has it, else in every
// relation whose columns are unknown; t.c and db.t.c in the relation they
// name. A name with more parts than that is a column's subcolumn. When more
// than one joined relation has an unqualified name, the leftmost one's
// column is returned, and the name is ambiguous unless the FROM clause
// joins once.
func lookup(parts []string, sc *spanScope) (spanColumn, found) {
	try := func(rels []*relation, name string) (spanColumn, found) {
		var (
			first spanColumn
			n     int
		)
		for _, rel := range rels {
			if rel.open {
				continue
			}
			c, ok := rel.column(name)
			if !ok || c.hidden {
				continue
			}
			if rel.arrayJoin {
				return c, found1
			}
			if n == 0 {
				first = c
			}
			n++
		}
		switch {
		case n > 1 && sc.joins > 1:
			return first, foundMany
		case n > 0:
			return first, found1
		}
		var open []spanColumn
		for _, rel := range rels {
			if rel.open {
				c, _ := rel.column(name)
				open = append(open, c)
			}
		}
		if len(open) == 0 {
			return spanColumn{}, foundNone
		}
		c := open[0]
		for _, o := range open[1:] {
			c.sources = union(c.sources, o.sources)
		}
		return c, found1
	}
	named := func(q []string) []*relation {
		var out []*relation
		for _, rel := range sc.relations {
			if qualifies(rel, q) {
				out = append(out, rel)
			}
		}
		return out
	}
	if len(parts) >= 3 {
		if rels := named(parts[:2]); len(rels) > 0 {
			return try(rels, parts[2])
		}
	}
	if len(parts) >= 2 {
		if rels := named(parts[:1]); len(rels) > 0 {
			return try(rels, parts[1])
		}
	}
	return try(sc.relations, parts[0])
}

func union(a, b SourceColumnSet) SourceColumnSet {
	out := SourceColumnSet{}
	for k := range a {
		out[k] = true
	}
	for k := range b {
		out[k] = true
	}
	return out
}

func isPlainColumn(e parser.Expr) bool {
	if a, ok := e.(*parser.AliasExpr); ok {
		e = a.Expr
	}
	_, ok := identParts(e)
	return ok
}

// children are an expression's direct subexpressions.
func children(e parser.Expr) []parser.Expr {
	switch n := e.(type) {
	case *parser.AliasExpr:
		return []parser.Expr{n.Expr}
	case *parser.UnaryExpr:
		return []parser.Expr{n.Expr}
	case *parser.BinaryExpr:
		return []parser.Expr{n.Left, n.Right}
	case *parser.TernaryExpr:
		return []parser.Expr{n.Cond, n.Then, n.Else}
	case *parser.FuncCallExpr:
		out := append(append([]parser.Expr(nil), n.Params...), n.Args...)
		if n.Over != nil {
			out = append(out, n.Over.PartitionBy...)
			for _, o := range n.Over.OrderBy {
				out = append(out, o.Expr)
			}
		}
		return out
	case *parser.CastExpr:
		return []parser.Expr{n.Expr}
	case *parser.CaseExpr:
		out := []parser.Expr{n.Operand}
		for _, w := range n.Whens {
			out = append(out, w.Cond, w.Then)
		}
		return append(out, n.Else)
	case *parser.BetweenExpr:
		return []parser.Expr{n.Expr, n.Low, n.High}
	case *parser.InExpr:
		return []parser.Expr{n.Left, n.Right}
	case *parser.LikeExpr:
		return []parser.Expr{n.Left, n.Right}
	case *parser.IsNullExpr:
		return []parser.Expr{n.Expr}
	case *parser.ArrayExpr:
		return n.Items
	case *parser.TupleExpr:
		return n.Items
	case *parser.ParenExpr:
		return []parser.Expr{n.Expr}
	case *parser.IndexExpr:
		return []parser.Expr{n.Expr, n.Index}
	case *parser.TupleAccessExpr:
		return []parser.Expr{n.Expr}
	case *parser.LambdaExpr:
		return []parser.Expr{n.Body}
	case *parser.IntervalExpr:
		return []parser.Expr{n.Value}
	case *parser.ExtractExpr:
		return []parser.Expr{n.From}
	case *parser.TrimExpr:
		return []parser.Expr{n.Chars, n.From}
	case *parser.SubstringExpr:
		return []parser.Expr{n.Expr, n.From, n.For}
	case *parser.StarExpr:
		var out []parser.Expr
		for _, t := range n.Transformers {
			for _, a := range t.Replaces {
				out = append(out, a)
			}
			out = append(out, t.ApplyExpr)
		}
		return out
	}
	return nil
}

// walkExpr visits e and its subexpressions, outside subqueries, while visit
// returns true.
func walkExpr(e parser.Expr, visit func(parser.Expr) bool) {
	if e == nil || !visit(e) {
		return
	}
	for _, c := range children(e) {
		walkExpr(c, visit)
	}
}
