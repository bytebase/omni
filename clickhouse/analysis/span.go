package analysis

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
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

// TableLister is a Database that names its tables and views, which
// merge('db', 're') reads. A database that does not list them fails a
// query that calls merge rather than omit what it reads.
type TableLister interface {
	TableNames() []string
}

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
	r.preferColumns = preferColumns(query)
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

// preferColumns reports whether a query's SETTINGS turn on
// prefer_column_name_to_alias.
func preferColumns(q *parser.SelectUnionStmt) bool {
	settings := []*parser.SettingExpr{}
	if q.Output != nil {
		settings = append(settings, q.Output.Settings...)
	}
	for _, s := range q.Items {
		settings = append(settings, s.Settings...)
	}
	for _, st := range settings {
		if !strings.EqualFold(st.Name, "prefer_column_name_to_alias") {
			continue
		}
		if lit, ok := st.Value.(*parser.LiteralExpr); ok {
			return lit.Val == "1" || strings.EqualFold(lit.Val, "true")
		}
	}
	return false
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
	// virtual marks a catalog table, whose virtual columns a name reads.
	virtual bool
	// starTables are the tables an unexpanded * of a derived relation
	// read, whose columns an unknown name of it reads.
	starTables SourceColumnSet
}

func (rel *relation) column(name string) (spanColumn, bool) {
	for _, c := range rel.columns {
		if c.name == name && !c.star {
			c.rel = rel
			return c, true
		}
	}
	if rel.open || rel.virtual && virtualColumns[name] {
		c := spanColumn{name: name, sources: SourceColumnSet{}, plain: true, rel: rel}
		if rel.table != "" {
			c.sources[ColumnResource{Database: rel.database, Table: rel.table, Column: name}] = true
		}
		for t := range rel.starTables {
			c.sources[ColumnResource{Database: t.Database, Table: t.Table, Column: name}] = true
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
	// query marks a SELECT's scope, and with the WITH scope above it: the
	// two are one level of name resolution.
	query, with bool
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
	// working is a recursive CTE's columns so far, which its own rows have.
	working []spanColumn
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
	// preferColumns is the query's prefer_column_name_to_alias: a column
	// name resolves to the column before an alias of that name.
	preferColumns bool
	// aliasMemo holds what an alias read where it was used; skipped counts
	// the aliases passed over because they were being resolved.
	aliasMemo map[aliasUse]SourceColumnSet
	skipped   int
	// err is why a view's definition could not be read: the span fails
	// rather than omit what the view reads.
	err   error
	views map[ColumnResource][]spanColumn
	// resolving holds the views and aliases being resolved, so one that
	// reads itself stops.
	resolving map[any]bool
}

// aliasUse is an alias's expression used in one scope.
type aliasUse struct {
	expr parser.Expr
	sc   *spanScope
}

func newSpanResolver(opts ResolveOptions, text string) *spanResolver {
	return &spanResolver{
		aliasMemo: map[aliasUse]SourceColumnSet{},
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
	var (
		out []spanColumn
		// The first query's WITH clause is every query's
		// (enable_global_with_statement, on by default).
		with *spanScope
	)
	for i, item := range q.Items {
		cols, w, c := r.selectStmt(item, outer, ctes, with)
		if i == 0 {
			out, with, ctes = cols, w, c
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

// selectStmt resolves one SELECT and returns, with its columns, its WITH
// clause's aliases and CTEs, which a set operation's later queries inherit
// (inherited holds the first query's aliases).
func (r *spanResolver) selectStmt(s *parser.SelectStmt, outer *spanScope, ctes *cteBinding, inherited *spanScope) ([]spanColumn, *spanScope, *cteBinding) {
	if s.Nested != nil {
		return r.query(s.Nested, outer, ctes), nil, ctes
	}
	// The WITH clause's expression aliases are visible to its CTEs too, so
	// they are a scope of their own, between the SELECT's and the outer
	// query's, and part of the SELECT's level.
	with := &spanScope{aliases: map[string]parser.Expr{}, with: true, parent: outer}
	if inherited != nil {
		for k, v := range inherited.aliases {
			with.aliases[k] = v
		}
	}
	sc := &spanScope{aliases: map[string]parser.Expr{}, windows: map[string]*parser.WindowSpec{}, query: true, parent: with}
	for _, w := range s.Windows {
		sc.windows[w.Name] = w.Spec
	}
	for _, w := range s.With {
		switch {
		case w.Subquery != nil:
			ctes = &cteBinding{name: w.Alias, query: w.Subquery, scope: with, parent: ctes}
			ctes.ctes = ctes.parent
			if w.Recursive {
				ctes.ctes = ctes
			}
		case w.Alias != "":
			with.aliases[w.Alias] = w.Expr
		}
	}
	// An expression alias is visible to the whole SELECT, wherever it is
	// defined.
	for _, c := range s.Columns {
		collectAliases(c, sc.aliases)
	}
	for _, e := range []parser.Expr{s.Prewhere, s.Where, s.Having, s.Qualify} {
		collectAliases(e, sc.aliases)
	}
	if s.GroupBy != nil {
		for _, e := range s.GroupBy.Exprs {
			collectAliases(e, sc.aliases)
		}
	}
	for _, o := range s.OrderBy {
		collectAliases(o.Expr, sc.aliases)
	}
	if s.From != nil {
		// A FROM subquery sees the WITH clause's aliases, not this
		// SELECT's columns.
		r.from(s.From, sc, with, ctes)
	}
	if s.ArrayJoin != nil {
		// An ARRAY JOIN column is an element of the array it names. An
		// alias names a new column, which * does not expand to; a bare
		// column keeps its name and its place, and reads what the array
		// reads.
		rel := &relation{arrayJoin: true}
		for _, e := range s.ArrayJoin.Exprs {
			if a, ok := e.(*parser.AliasExpr); ok {
				rel.columns = append(rel.columns, spanColumn{name: a.Alias, sources: r.expr(a.Expr, sc, ctes)})
				continue
			}
			if parts, ok := identParts(e); ok && !r.resolvesAsName(parts, sc) && nestedGroup(sc, strings.Join(parts, ".")) {
				continue
			}
			r.expr(e, sc, ctes)
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
		if items, ok := untupled(item); ok {
			// untuple(t) is one result column per element of t, each
			// taken from the whole tuple.
			arg := item.(*parser.FuncCallExpr).Args[0]
			name, sources := (&namer{r: r, sc: sc, ctes: ctes}).name(arg), r.expr(arg, sc, ctes)
			for i := range items {
				out = append(out, spanColumn{name: fmt.Sprintf("tupleElement(%s, '%d')", name, i+1), sources: sources})
			}
			continue
		}
		c := spanColumn{
			name:    r.projectionName(item, sc, ctes),
			sources: r.expr(item, sc, ctes),
			plain:   isPlainColumn(item),
		}
		if parts, ok := identParts(item); ok && r.isAlias(parts, sc) {
			// An alias is a plain field only when it names a column.
			c.plain = isPlainColumn(r.aliasExpr(parts[0], sc))
		} else if ok {
			// A column is named by its own name (a Nested column's is
			// dotted), and a qualified one keeps its qualifier unless the
			// bare name reads the same column.
			found, ok := lookup(parts, sc)
			switch n := nameParts(parts, found.name); {
			case ok == found1 && found.rel != nil && n > 0:
				// A column of a subquery or CTE is as plain as its
				// projection.
				c.name, c.plain = found.name, found.plain
				if len(parts) > n && !bareNames(found, sc) {
					c.name = strings.Join(parts, ".")
				}
			case len(parts) > 1 && (ok == foundNone || n == 0):
				// A member of a tuple (an alias's, or a column's) is named
				// as written.
				c.name = strings.Join(parts, ".")
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
		r.exprs([]parser.Expr{s.LimitBy.Limit, s.LimitBy.Offset}, sc, ctes)
		r.exprs(s.LimitBy.By, sc, ctes)
	}
	if s.Limit != nil {
		r.exprs([]parser.Expr{s.Limit.Limit, s.Limit.Offset}, sc, ctes)
	}
	return out, with, ctes
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
		case *parser.SubqueryExpr, *parser.LambdaExpr, *parser.StarExpr:
			// A star's REPLACE (e AS c) names a column, not an alias.
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
		if strings.EqualFold(t.Func.Name, "merge") {
			return r.merge(t.Func, name)
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
				// A recursive CTE's own rows: the columns of the rows
				// before, which r.cte iterates to a fixed point.
				if cte.working == nil {
					return &relation{name: name, open: true}
				}
				return derived(name, renamed(cte.working, t.ColumnAliases))
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
	_, isView := r.views[ColumnResource{Database: database, Table: t.Table}]
	return &relation{database: database, table: t.Table, name: name, columns: renamed(cols, t.ColumnAliases), open: open, virtual: !open && !isView}
}

// virtualColumns are the MergeTree family's virtual columns, which a
// table's synced columns do not list and only a name reads.
var virtualColumns = map[string]bool{
	"_part": true, "_part_index": true, "_part_starting_offset": true, "_part_offset": true,
	"_part_granule_offset": true, "_part_uuid": true, "_part_data_version": true,
	"_partition_id": true, "_partition_value": true, "_sample_factor": true,
	"_row_exists": true, "_block_number": true, "_block_offset": true, "_disk_name": true,
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
		if c.ctes == c {
			// A recursive CTE's rows feed its next iteration, which may
			// move a value to another column: iterate until no column
			// reads more.
			for i := 0; i < 32; i++ {
				c.working = c.columns
				next := r.query(c.query, c.scope, c.ctes)
				grew := false
				merged := make([]spanColumn, len(c.columns))
				copy(merged, c.columns)
				for j := range merged {
					if j < len(next) && len(union(merged[j].sources, next[j].sources)) > len(merged[j].sources) {
						merged[j].sources, grew = union(merged[j].sources, next[j].sources), true
					}
				}
				c.columns = merged
				if !grew {
					break
				}
			}
		}
		c.resolving = false
	}
	return c.columns
}

// merge resolves merge('db', 're') and merge('re'): the tables of a
// database whose names match, read as one; a column reads that column of
// each of them.
func (r *spanResolver) merge(f *parser.FuncCallExpr, name string) *relation {
	var args []string
	for _, a := range f.Args {
		lit, ok := a.(*parser.LiteralExpr)
		if !ok || lit.Kind != "STRING" {
			r.err = errors.Join(r.err, errors.New("merge() with a database or pattern that is not a string literal"))
			return nil
		}
		args = append(args, lit.Val)
	}
	database := r.database
	switch len(args) {
	case 1:
	case 2:
		database = args[0]
		args = args[1:]
	default:
		r.err = errors.Join(r.err, errors.New("merge() takes a database and a pattern"))
		return nil
	}
	re, err := regexp.Compile(args[0])
	if err != nil {
		r.err = errors.Join(r.err, fmt.Errorf("merge(): %w", err))
		return nil
	}
	var db Database
	if r.opts.Catalog.GetDatabase != nil {
		if db, err = r.opts.Catalog.GetDatabase(database); err != nil {
			r.err = errors.Join(r.err, err)
			return nil
		}
	}
	lister, ok := db.(TableLister)
	if !ok {
		r.err = errors.Join(r.err, fmt.Errorf("merge() reads the tables of %s that match %q, which the catalog does not list", database, args[0]))
		return nil
	}
	rel := &relation{name: name}
	for _, table := range lister.TableNames() {
		if !re.MatchString(table) {
			continue
		}
		cols, _, ok := r.relation(database, table)
		if !ok {
			continue
		}
		for _, c := range cols {
			i := slices.IndexFunc(rel.columns, func(o spanColumn) bool { return o.name == c.name })
			if i < 0 {
				rel.columns = append(rel.columns, spanColumn{name: c.name, sources: c.sources, plain: true})
				continue
			}
			rel.columns[i].sources = union(rel.columns[i].sources, c.sources)
		}
	}
	return rel
}

// derived is a subquery's or CTE's relation. A * no metadata expanded
// leaves its columns unknown: any name then resolves, to that column of
// the tables the * read.
func derived(name string, cols []spanColumn) *relation {
	rel := &relation{name: name, columns: cols}
	for _, c := range cols {
		if !c.star {
			continue
		}
		// The * stays in its place, so a * over the relation keeps it.
		rel.open = true
		if rel.starTables == nil {
			rel.starTables = SourceColumnSet{}
		}
		for t := range c.sources {
			rel.starTables[ColumnResource{Database: t.Database, Table: t.Table}] = true
		}
	}
	return rel
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
		if rel.arrayJoin || len(s.Qualifier) > 0 && !qualifies(rel, s.Qualifier) {
			continue
		}
		if rel.open && len(rel.columns) == 0 {
			sources := union(SourceColumnSet{}, rel.starTables)
			if rel.table != "" {
				sources[ColumnResource{Database: rel.database, Table: rel.table}] = true
			}
			out = append(out, spanColumn{name: "*", sources: sources, star: true})
			continue
		}
		for _, c := range rel.columns {
			switch {
			case c.star:
				out = append(out, c)
			case !c.hidden:
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
	var exceptRe []*regexp.Regexp
	addRe := func(pattern string) {
		if re, err := regexp.Compile(pattern); pattern != "" && err == nil {
			exceptRe = append(exceptRe, re)
		}
	}
	for _, c := range s.Except {
		except[c] = true
	}
	addRe(s.ExceptRe)
	for _, t := range s.Transformers {
		if t.Kind == "EXCEPT" {
			for _, c := range t.Columns {
				except[c] = true
			}
			addRe(t.Regex)
		}
	}
	var kept []spanColumn
	for _, c := range out {
		excluded := except[c.name]
		for _, re := range exceptRe {
			excluded = excluded || re.MatchString(c.name)
		}
		if !excluded {
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
			// Each column is the function, or the lambda, applied to it.
			applied := r.expr(t.ApplyExpr, sc, ctes)
			for i := range kept {
				kept[i].name = (&namer{r: r, sc: sc, ctes: ctes}).apply(t, kept[i].name)
				kept[i].sources = union(kept[i].sources, applied)
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
	if c.rel.name == "" || bareNames(c, sc) {
		return c.name
	}
	return c.rel.name + "." + c.name
}

// bareNames reports whether a column's bare name reads it: no alias has the
// name, and no relation before the column's has a column of that name.
func bareNames(c spanColumn, sc *spanScope) bool {
	for s := sc; s != nil; s = s.parent {
		if _, ok := s.aliases[c.name]; ok {
			return false
		}
	}
	bare, ok := lookup([]string{c.name}, sc)
	return ok == found1 && bare.rel == c.rel
}

// isAlias reports whether a name is an alias the scope resolves before its
// columns.
func (r *spanResolver) isAlias(parts []string, sc *spanScope) bool {
	if len(parts) != 1 || r.preferColumns {
		return false
	}
	for s := sc; s != nil; s = s.parent {
		if _, ok := s.aliases[parts[0]]; ok {
			return true
		}
	}
	return false
}

// nameParts is how many of an identifier's last parts spell a column's
// name: one for a quoted `a.b`, two for a Nested column's a.b; 0 when they
// spell another name (a subcolumn).
func nameParts(parts []string, name string) int {
	for n := 1; n <= len(parts); n++ {
		if strings.Join(parts[len(parts)-n:], ".") == name {
			return n
		}
	}
	return 0
}

// nestedGroup unrolls ARRAY JOIN of a Nested column's name: its arrays
// (name.a, name.b, ...) unroll together, so each of them reads them all. It
// reports whether any relation has the group.
func nestedGroup(sc *spanScope, name string) bool {
	found := false
	for _, rel := range sc.relations {
		group := SourceColumnSet{}
		for _, c := range rel.columns {
			if strings.HasPrefix(c.name, name+".") {
				group = union(group, c.sources)
			}
		}
		if len(group) == 0 {
			continue
		}
		found = true
		for k, c := range rel.columns {
			if strings.HasPrefix(c.name, name+".") {
				rel.columns[k].sources = group
			}
		}
	}
	return found
}

// aliasExpr is the expression an alias in scope stands for.
func (r *spanResolver) aliasExpr(name string, sc *spanScope) parser.Expr {
	for s := sc; s != nil; s = s.parent {
		if e, ok := s.aliases[name]; ok {
			return e
		}
	}
	return nil
}

// untupled is the elements of an untuple((a, b)) item.
func untupled(e parser.Expr) ([]parser.Expr, bool) {
	f, ok := e.(*parser.FuncCallExpr)
	if !ok || !strings.EqualFold(f.Name, "untuple") || len(f.Args) != 1 {
		return nil, false
	}
	switch t := f.Args[0].(type) {
	case *parser.TupleExpr:
		return t.Items, len(t.Items) > 0
	case *parser.FuncCallExpr:
		if strings.EqualFold(t.Name, "tuple") {
			return t.Args, len(t.Args) > 0
		}
	}
	return nil, false
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
		if rel.arrayJoin {
			continue
		}
		for _, c := range rel.columns {
			if !c.hidden && !c.star && re.MatchString(c.name) {
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
		// A star in an expression is every column it expands to.
		for _, c := range r.star(n, sc, ctes) {
			out = union(out, c.sources)
		}
		return out
	case *parser.FuncCallExpr:
		args := n.Args
		if countStar(n) {
			// count(*) reads no column; its window still does.
			args = nil
		}
		if lambda := r.lambdaAlias(n.Name, sc); lambda != nil {
			// A WITH-bound lambda called by its alias reads what its body
			// reads besides its arguments.
			r.resolving[lambda] = true
			out = r.expr(lambda, sc, ctes)
			delete(r.resolving, lambda)
		}
		if re, ok := columnsMatcher(n); ok {
			// COLUMNS('re') in an expression is the columns it matches.
			for _, c := range r.matching(re, sc) {
				out = union(out, c.sources)
			}
			return out
		}
		out = union(out, union(r.exprs(n.Params, sc, ctes), r.exprs(args, sc, ctes)))
		over := n.Over
		if n.OverName != "" {
			over = sc.window(n.OverName)
		}
		// A window may inherit another; a name that inherits itself is
		// read once.
		seen := map[*parser.WindowSpec]bool{}
		for ; over != nil && !seen[over]; over = sc.window(over.BaseWindow) {
			seen[over] = true
			out = union(out, r.window(over, sc, ctes))
		}
		return out
	case *parser.InExpr:
		out = r.expr(n.Left, sc, ctes)
		if parts, ok := identParts(n.Right); ok && len(parts) <= 2 && !r.resolvesAsName(parts, sc) {
			// x IN t reads table t, as x IN (SELECT * FROM t) does.
			return union(out, r.inTable(parts, sc, ctes))
		}
		return union(out, r.expr(n.Right, sc, ctes))
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

// countStar reports whether f is count(*), which reads no column;
// count(DISTINCT *) reads them all.
func countStar(f *parser.FuncCallExpr) bool {
	if len(f.Args) != 1 || f.Distinct || !strings.EqualFold(f.Name, "count") {
		return false
	}
	s, ok := f.Args[0].(*parser.StarExpr)
	return ok && len(s.Qualifier) == 0
}

// lambdaAlias is the lambda an alias in scope binds name to, or nil.
func (r *spanResolver) lambdaAlias(name string, sc *spanScope) *parser.LambdaExpr {
	for s := sc; s != nil; s = s.parent {
		if e, ok := s.aliases[name]; ok {
			if l, ok := e.(*parser.LambdaExpr); ok && !r.resolving[l] {
				return l
			}
			return nil
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
	name := strings.Join(parts, ".")
	// A name resolves in the closest query that has it: there, an alias
	// before a column (unless prefer_column_name_to_alias), then the
	// enclosing query.
	for lv := sc; lv != nil; {
		scopes, next := level(lv)
		if len(parts) == 1 && r.preferColumns {
			for _, s := range scopes {
				if c, found := lookup(parts, s); found == found1 {
					return c.sources
				}
			}
		}
		if len(parts) == 1 {
			if out, ok := r.alias(parts[0], scopes, sc, ctes); ok {
				return out
			}
		}
		for _, s := range scopes {
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
		if len(parts) > 1 {
			// t.e where t is no relation or column: a member of an alias's
			// tuple, which reads what the alias reads.
			if out, ok := r.alias(parts[0], scopes, sc, ctes); ok {
				return out
			}
		}
		lv = next
	}
	r.missing(NotFoundError{Column: &name})
	return SourceColumnSet{}
}

// level is the scopes of one query level from s: lambda scopes, the
// SELECT's, and its WITH scope; next is the enclosing query's.
func level(s *spanScope) (scopes []*spanScope, next *spanScope) {
	for ; s != nil; s = s.parent {
		scopes = append(scopes, s)
		if s.query {
			if p := s.parent; p != nil && p.with {
				return append(scopes, p), p.parent
			}
			return scopes, s.parent
		}
	}
	return scopes, nil
}

// alias resolves name as an alias of one level's scopes, used in sc.
func (r *spanResolver) alias(name string, scopes []*spanScope, sc *spanScope, ctes *cteBinding) (SourceColumnSet, bool) {
	{
		for _, s := range scopes {
			if e, ok := s.aliases[name]; ok && !r.resolving[e] {
				// An alias stands for its expression where it is used. It is
				// resolved once per scope, so aliases built on aliases stay
				// linear.
				key := aliasUse{e, sc}
				if out, ok := r.aliasMemo[key]; ok {
					return out, true
				}
				skipped := r.skipped
				r.resolving[e] = true
				out := r.expr(e, sc, ctes)
				delete(r.resolving, e)
				if r.skipped == skipped {
					// A result that broke a cycle depends on what was being
					// resolved, so only one that did not is kept.
					r.aliasMemo[key] = out
				}
				return out, true
			}
			if _, ok := s.aliases[name]; ok {
				r.skipped++
			}
		}
	}
	return nil, false
}

// resolvesAsName reports whether parts name a lambda parameter, an alias or
// a column in scope, without reporting anything.
func (r *spanResolver) resolvesAsName(parts []string, sc *spanScope) bool {
	for s := sc; s != nil; s = s.parent {
		if s.lambda[parts[0]] {
			return true
		}
		if _, ok := s.aliases[parts[0]]; ok && len(parts) == 1 {
			return true
		}
		if _, found := lookup(parts, s); found != foundNone {
			return true
		}
	}
	return false
}

// inTable resolves the table of x IN t: a CTE, or a table or view of the
// catalog. A name that is neither is a missing identifier, as ClickHouse
// reports it.
func (r *spanResolver) inTable(parts []string, sc *spanScope, ctes *cteBinding) SourceColumnSet {
	t := &parser.TableExpr{Table: parts[len(parts)-1]}
	if len(parts) == 2 {
		t.Database = parts[0]
	} else if cte := ctes.find(t.Table); cte == nil && !r.tableExists(r.database, t.Table) {
		name := strings.Join(parts, ".")
		r.missing(NotFoundError{Column: &name})
		return SourceColumnSet{}
	}
	out := SourceColumnSet{}
	if rel := r.tableExpr(t, sc, ctes); rel != nil {
		for _, c := range rel.columns {
			out = union(out, c.sources)
		}
	}
	return out
}

// tableExists reports whether the catalog has a table or view, or cannot
// say.
func (r *spanResolver) tableExists(database, table string) bool {
	if systemDatabases[database] || r.opts.Catalog.GetDatabase == nil {
		return true
	}
	db, err := r.opts.Catalog.GetDatabase(database)
	if err != nil {
		// The span fails rather than report an unreadable catalog as a
		// missing name.
		r.err = errors.Join(r.err, err)
		return false
	}
	if db == nil {
		return false
	}
	return db.GetTable(table) != nil || db.GetView(table) != nil
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
	try := func(rels []*relation, name string, qualified bool) (spanColumn, found) {
		var (
			first spanColumn
			n     int
		)
		for _, rel := range rels {
			if rel.open {
				continue
			}
			c, ok := rel.column(name)
			if !ok || c.hidden && !qualified {
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
	// rest names a column of rels: a Nested column's dotted name whole,
	// else its first part, whose subcolumn the others are.
	column := func(rels []*relation, rest []string, qualified bool) (spanColumn, found) {
		if len(rest) > 1 {
			if c, f := try(rels, strings.Join(rest, "."), qualified); f != foundNone && !c.rel.open {
				return c, f
			}
		}
		return try(rels, rest[0], qualified)
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
			return column(rels, parts[2:], true)
		}
	}
	if len(parts) >= 2 {
		if rels := named(parts[:1]); len(rels) > 0 {
			return column(rels, parts[1:], true)
		}
	}
	return column(sc.relations, parts, false)
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
