package analysis

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/starrocks/ast"
	"github.com/bytebase/omni/starrocks/parser"
)

// QuerySpanResult is one result column of a query and the columns it reads.
type QuerySpanResult = analysis.QuerySpanResult

// NotFoundError reports a database, table or column the catalog does not
// have. A nil field is not part of the report.
type NotFoundError = analysis.ResourceNotFoundError

// Catalog is the synced metadata of the instance a statement runs on.
type Catalog struct {
	// GetDatabase returns a database's metadata, or nil when the database
	// has not been synced. A nil GetDatabase means no metadata is
	// available, and the span resolves no column.
	GetDatabase func(name string) (Database, error)
}

// Database is one database's synced metadata. StarRocks keeps its tables,
// views and materialized views directly in the database.
type Database interface {
	// Name is the database's name as the metadata spells it.
	Name() string
	// GetTable returns the named table, or nil.
	GetTable(name string) *metadata.TableMetadata
	// GetView returns the named view, or nil.
	GetView(name string) *metadata.ViewMetadata
	// GetMaterializedView returns the named materialized view, or nil.
	GetMaterializedView(name string) *metadata.MaterializedViewMetadata
}

// internalCatalog is the catalog that holds StarRocks's own databases; a
// table of any other catalog lives in an external source the metadata
// does not describe.
const internalCatalog = "default_catalog"

// spanColumn is one column a relation or a query produces.
type spanColumn struct {
	name    string
	sources SourceColumnSet
	plain   bool
	// star marks a * the metadata could not expand.
	star bool
}

// relation is a table, view, CTE, derived table or table function in a
// query's FROM clause.
type relation struct {
	// catalog and database are the catalog and database of a table or
	// view, which a column reference may name (catalog.db.t.c); empty for
	// a derived relation.
	catalog  string
	database string
	// name is the alias, or the table's name as written.
	name string
	// qualified is how StarRocks qualifies an analyzed column of a table
	// or view (TableName.toString): its database, then its alias or name,
	// after its catalog when that is not the internal one. Empty for any
	// other relation, whose columns StarRocks qualifies by name.
	qualified string
	columns   []spanColumn
	// open marks a relation whose columns the metadata does not list (a
	// table function, a system table, a table of another catalog): any
	// column name resolves on it, to source, which is nil for a column
	// that reads no table.
	open   bool
	source func(column string) SourceColumnSet
}

// spanScope is the relations one SELECT reads, inside its enclosing
// queries' scopes, which correlated references reach.
type spanScope struct {
	relations []*relation
	// outputs are the SELECT list's columns, which GROUP BY, HAVING and
	// ORDER BY may name by alias.
	outputs []spanColumn
	parent  *spanScope
}

// cteBinding is one CTE in scope.
type cteBinding struct {
	name    string
	columns []spanColumn
	parent  *cteBinding
}

// spanResolver resolves a query's columns through the catalog.
type spanResolver struct {
	opts ResolveOptions
	// database is where an unqualified table name resolves: the session's
	// database, or a view's while its definition resolves.
	database string
	// text is the SQL the tree being resolved was parsed from, which an
	// unnamed result column is named by.
	text     string
	tables   SourceColumnSet
	notFound *NotFoundError
	// err is why a view's stored definition could not be read: the span
	// fails rather than omit what the view reads.
	err   error
	views map[ColumnResource][]spanColumn
	// resolving holds the views whose definitions are being resolved, so
	// a view that reads itself stops.
	resolving map[ColumnResource]bool
	// qualify, while set, makes printed render a column reference as
	// StarRocks prints it once analyzed, resolved in this scope.
	qualify *spanScope
}

func newSpanResolver(opts ResolveOptions, text string) *spanResolver {
	return &spanResolver{
		opts:      opts,
		database:  opts.DefaultDatabase,
		text:      text,
		tables:    SourceColumnSet{},
		views:     map[ColumnResource][]spanColumn{},
		resolving: map[ColumnResource]bool{},
	}
}

// isQuery reports whether node is a query whose columns the resolver reads.
func isQuery(node ast.Node) bool {
	switch node.(type) {
	case *ast.SelectStmt, *ast.SetOpStmt, *ast.ParenSelect:
		return true
	}
	return false
}

func (r *spanResolver) fail(err error) {
	if r.err == nil {
		r.err = err
	}
}

func (r *spanResolver) missing(nf NotFoundError) {
	if r.notFound == nil {
		r.notFound = &nf
	}
}

// query resolves a query node and returns its columns.
func (r *spanResolver) query(node ast.Node, outer *spanScope, ctes *cteBinding) []spanColumn {
	switch n := node.(type) {
	case *ast.SelectStmt:
		return r.selectStmt(n, outer, ctes, false)
	case *ast.SetOpStmt:
		return r.setOp(n, outer, ctes)
	case *ast.ParenSelect:
		ctes = r.with(n.With, outer, ctes)
		cols := r.query(n.Sel, outer, ctes)
		r.trailing(n.OrderBy, n.Limit, n.Offset, &spanScope{outputs: cols, parent: outer}, ctes)
		return cols
	case *ast.SubqueryExpr:
		return r.subquery(n, outer, ctes)
	}
	return nil
}

// subquery resolves a subquery's body: the tree the strict parse built in
// place, or else its text parsed again.
func (r *spanResolver) subquery(sub *ast.SubqueryExpr, outer *spanScope, ctes *cteBinding) []spanColumn {
	if sub.Query != nil {
		return r.query(sub.Query, outer, ctes)
	}
	file, err := parser.Parse(sub.RawText)
	if err != nil || file == nil || len(file.Stmts) != 1 {
		// GetQuerySpan, which ran first, fails the span on such a body.
		return nil
	}
	saved := r.text
	r.text = sub.RawText
	defer func() { r.text = saved }()
	return r.query(file.Stmts[0], outer, ctes)
}

// with binds a WITH clause's CTEs, each resolved with the ones before it
// in scope.
func (r *spanResolver) with(w *ast.WithClause, outer *spanScope, ctes *cteBinding) *cteBinding {
	if w == nil {
		return ctes
	}
	for _, cte := range w.CTEs {
		if cte == nil {
			continue
		}
		cols := r.query(cte.Query, outer, ctes)
		cols = renamed(cols, cte.Columns)
		ctes = &cteBinding{name: cte.Name, columns: cols, parent: ctes}
	}
	return ctes
}

// renamed gives cols the names of a column list, when there is one.
func renamed(cols []spanColumn, names []string) []spanColumn {
	if len(names) == 0 {
		return cols
	}
	out := make([]spanColumn, len(cols))
	copy(out, cols)
	for i := range out {
		if i < len(names) {
			out[i].name = names[i]
		}
	}
	return out
}

// leftmostQuery returns the SELECT a set operation's left spine starts
// with, and its WITH clause, which scopes over every arm.
func leftmostQuery(node ast.Node) (ast.Node, *ast.WithClause) {
	for {
		switch n := node.(type) {
		case *ast.SetOpStmt:
			node = n.Left
		case *ast.SelectStmt:
			return n, n.With
		default:
			return nil, nil
		}
	}
}

func (r *spanResolver) setOp(n *ast.SetOpStmt, outer *spanScope, ctes *cteBinding) []spanColumn {
	// The grammar puts a set operation's WITH on the operation, or on its
	// first SELECT.
	ctes = r.with(n.With, outer, ctes)
	first, with := leftmostQuery(n)
	ctes = r.with(with, outer, ctes)
	var arms []ast.Node
	var flatten func(node ast.Node)
	flatten = func(node ast.Node) {
		if s, ok := node.(*ast.SetOpStmt); ok && s != n && len(s.OrderBy) == 0 && s.Limit == nil && s.Offset == nil {
			flatten(s.Left)
			flatten(s.Right)
			return
		}
		arms = append(arms, node)
	}
	flatten(n.Left)
	flatten(n.Right)
	var cols []spanColumn
	for i, arm := range arms {
		var armCols []spanColumn
		switch a := arm.(type) {
		case *ast.SelectStmt:
			armCols = r.selectStmt(a, outer, ctes, arm == first)
		default:
			armCols = r.query(arm, outer, ctes)
		}
		if i == 0 {
			cols = make([]spanColumn, len(armCols))
			for j, c := range armCols {
				cols[j] = spanColumn{name: c.name, sources: SourceColumnSet{}, star: c.star}
				for k := range c.sources {
					cols[j].sources[k] = true
				}
			}
			continue
		}
		for j := range cols {
			if j < len(armCols) {
				for k := range armCols[j].sources {
					cols[j].sources[k] = true
				}
				cols[j].star = cols[j].star || armCols[j].star
			}
		}
	}
	r.trailing(n.OrderBy, n.Limit, n.Offset, &spanScope{outputs: cols, parent: outer}, ctes)
	return cols
}

// trailing resolves ORDER BY, LIMIT and OFFSET, for the tables and missing
// columns they name.
func (r *spanResolver) trailing(orderBy []*ast.OrderByItem, limit, offset ast.Node, sc *spanScope, ctes *cteBinding) {
	for _, o := range orderBy {
		if o != nil && o.Expr != nil {
			r.expr(o.Expr, sc, ctes)
		}
	}
	if limit != nil {
		r.expr(limit, sc, ctes)
	}
	if offset != nil {
		r.expr(offset, sc, ctes)
	}
}

// selectStmt resolves one SELECT. withBound is set when a set operation
// has bound its WITH clause already.
func (r *spanResolver) selectStmt(stmt *ast.SelectStmt, outer *spanScope, ctes *cteBinding, withBound bool) []spanColumn {
	if !withBound {
		ctes = r.with(stmt.With, outer, ctes)
	}
	sc := &spanScope{parent: outer}
	for _, from := range stmt.From {
		r.from(from, sc, outer, ctes)
	}
	if stmt.Pivot != nil {
		sc = r.pivot(stmt.Pivot, sc, ctes)
	}
	if stmt.Where != nil {
		r.expr(stmt.Where, sc, ctes)
	}
	var cols []spanColumn
	for _, item := range stmt.Items {
		if item == nil {
			continue
		}
		if item.Star {
			cols = append(cols, r.star(item, sc)...)
			continue
		}
		if item.Expr == nil {
			continue
		}
		cols = append(cols, spanColumn{
			name:    r.resultName(item),
			sources: r.expr(item.Expr, sc, ctes),
			plain:   isColumnRef(item.Expr),
		})
	}
	// GROUP BY, HAVING, QUALIFY and ORDER BY may name a result by its
	// alias.
	sc.outputs = cols
	for _, g := range stmt.GroupBy {
		r.expr(g, sc, ctes)
	}
	if stmt.Having != nil {
		r.expr(stmt.Having, sc, ctes)
	}
	if stmt.Qualify != nil {
		r.expr(stmt.Qualify, sc, ctes)
	}
	r.trailing(stmt.OrderBy, stmt.Limit, stmt.Offset, sc, ctes)
	return cols
}

// pivot returns the scope a PIVOT leaves over the FROM relations of sc, as
// StarRocks's QueryAnalyzer.visitPivotRelation builds it. Its columns are
// the FROM columns the PIVOT does not use, which it groups by, each kept in
// its relation, then one column per IN value and aggregate in a relation of
// no name. A column the PIVOT uses is an aggregate's direct column argument
// or a FOR column (PivotRelation.getUsedColumns), matched by name in any
// case. A generated column reads what its aggregate reads, and the FOR
// columns, which StarRocks's rewrite, agg(CASE WHEN for = value THEN arg
// END), compares. A used column is out of scope after the PIVOT.
func (r *spanResolver) pivot(p *ast.PivotClause, sc *spanScope, ctes *cteBinding) *spanScope {
	used := map[string]bool{}
	type aggregate struct {
		name    string
		sources SourceColumnSet
	}
	var aggs []aggregate
	for _, agg := range p.Aggregations {
		if agg == nil || agg.Func == nil {
			continue
		}
		if call, ok := agg.Func.(*ast.FuncCallExpr); ok {
			for _, arg := range call.Args {
				for {
					paren, ok := arg.(*ast.ParenExpr)
					if !ok {
						break
					}
					arg = paren.Expr
				}
				if ref, ok := arg.(*ast.ColumnRef); ok && ref.Name != nil && len(ref.Name.Parts) > 0 {
					used[strings.ToLower(ref.Name.Parts[len(ref.Name.Parts)-1])] = true
				}
			}
		}
		name := agg.Alias
		if name == "" {
			name = r.aggregateName(agg.Func, sc)
		}
		aggs = append(aggs, aggregate{name: name, sources: r.expr(agg.Func, sc, ctes)})
	}
	forSources := SourceColumnSet{}
	for _, c := range p.Columns {
		used[strings.ToLower(c)] = true
		for k := range r.column([]string{c}, sc) {
			forSources[k] = true
		}
	}
	pivoted := &spanScope{parent: sc.parent}
	for _, rel := range sc.relations {
		kept := *rel
		kept.columns = nil
		for _, c := range rel.columns {
			if !used[strings.ToLower(c.name)] {
				kept.columns = append(kept.columns, c)
			}
		}
		pivoted.relations = append(pivoted.relations, &kept)
	}
	generated := &relation{}
	for _, v := range p.Values {
		if v == nil {
			continue
		}
		value := v.Alias
		if value == "" {
			names := make([]string, len(v.Exprs))
			for i, e := range v.Exprs {
				names[i] = r.printed(e)
				// A DATE or DATETIME literal prints as its string.
				if c, ok := e.(*ast.CastExpr); ok {
					names[i] = r.printed(c.Expr)
				}
			}
			value = "{" + strings.Join(names, ",") + "}"
			if len(names) == 1 {
				value = names[0]
			}
		}
		for _, agg := range aggs {
			name := value
			if len(aggs) > 1 {
				name += "_" + agg.name
			}
			sources := clone(agg.sources)
			for k := range forSources {
				sources[k] = true
			}
			generated.columns = append(generated.columns, spanColumn{name: name, sources: sources})
		}
	}
	pivoted.relations = append(pivoted.relations, generated)
	return pivoted
}

// aggregateName is the name StarRocks gives an unaliased PIVOT aggregate
// (AstToStringBuilder.visitFunctionCall): the call as StarRocks prints it
// once analyzed, where each column carries its relation's name (see
// qualifiedColumn).
func (r *spanResolver) aggregateName(node ast.Node, sc *spanScope) string {
	saved := r.qualify
	r.qualify = sc
	defer func() { r.qualify = saved }()
	return r.printed(node)
}

// printedCall renders a call as StarRocks prints an analyzed one
// (AstToStringBuilder.visitFunctionCall): the function's name in lower
// case, then its arguments, each printed, separated by ", ". It reports
// false for a call with a qualified name, an ORDER BY, a SEPARATOR or a
// window, which keeps its text as written.
func (r *spanResolver) printedCall(call *ast.FuncCallExpr) (string, bool) {
	if call.Name == nil || len(call.Name.Parts) != 1 || len(call.OrderBy) > 0 || call.SeparatorExpr != nil || call.Separator != "" || call.Over != nil {
		return "", false
	}
	var b strings.Builder
	b.WriteString(strings.ToLower(call.Name.Parts[0]) + "(")
	if call.Star {
		b.WriteString("*")
	}
	if call.Distinct {
		b.WriteString("DISTINCT ")
	}
	args := make([]string, len(call.Args))
	for i, a := range call.Args {
		args[i] = r.printed(a)
	}
	b.WriteString(strings.Join(args, ", ") + ")")
	return b.String(), true
}

// from adds a FROM item's relations to sc. A derived table sees only the
// enclosing queries' scopes (outer); a table function sees the relations
// before it.
func (r *spanResolver) from(node ast.Node, sc, outer *spanScope, ctes *cteBinding) {
	switch n := node.(type) {
	case *ast.JoinClause:
		before := len(sc.relations)
		r.from(n.Left, sc, outer, ctes)
		middle := len(sc.relations)
		r.from(n.Right, sc, outer, ctes)
		if n.On != nil {
			r.expr(n.On, sc, ctes)
		}
		// A semi or anti join returns one side's rows only: the other side
		// is read, but its columns are out of scope after the join.
		switch n.Type {
		case ast.JoinLeftSemi, ast.JoinLeftAnti:
			sc.relations = sc.relations[:middle]
		case ast.JoinRightSemi, ast.JoinRightAnti:
			sc.relations = append(sc.relations[:before], sc.relations[middle:]...)
		}
	case *ast.TableRef:
		if rel := r.tableRef(n, outer, ctes); rel != nil {
			sc.relations = append(sc.relations, rel)
		}
	case *ast.InlineTable:
		// A VALUES table reads no table; its columns are column_0, ... or
		// its column aliases.
		rel := &relation{name: n.Alias}
		for i, row := range n.Rows {
			for j, e := range row {
				sources := r.expr(e, &spanScope{parent: outer}, ctes)
				if i == 0 {
					rel.columns = append(rel.columns, spanColumn{name: fmt.Sprintf("column_%d", j), sources: sources})
				} else if j < len(rel.columns) {
					for k := range sources {
						rel.columns[j].sources[k] = true
					}
				}
			}
		}
		rel.columns = renamed(rel.columns, n.ColumnAliases)
		sc.relations = append(sc.relations, rel)
	case *ast.TableFunctionRef:
		// A table function reads no table; it reads the columns its
		// arguments name, which a lateral call takes from the relations
		// before it.
		var sources SourceColumnSet
		if n.Call != nil {
			sources = r.expr(n.Call, sc, ctes)
		}
		name := n.Alias
		if name == "" && n.Call != nil && n.Call.Name != nil && len(n.Call.Name.Parts) > 0 {
			name = n.Call.Name.Parts[len(n.Call.Name.Parts)-1]
		}
		rel := &relation{name: name}
		if len(n.ColumnAliases) == 0 {
			rel.open = true
			rel.source = func(string) SourceColumnSet { return clone(sources) }
		}
		for _, c := range n.ColumnAliases {
			rel.columns = append(rel.columns, spanColumn{name: c, sources: clone(sources)})
		}
		sc.relations = append(sc.relations, rel)
	}
}

// tableRef resolves one FROM item to a relation, or nil when the catalog
// lacks what it names.
func (r *spanResolver) tableRef(ref *ast.TableRef, outer *spanScope, ctes *cteBinding) *relation {
	if ref.Subquery != nil {
		cols := r.subquery(ref.Subquery, outer, ctes)
		return derived(ref.Alias, renamed(cols, ref.ColumnAliases))
	}
	if ref.Name == nil || len(ref.Name.Parts) == 0 {
		return nil
	}
	parts := ref.Name.Parts
	alias := func(name string) string {
		if ref.Alias != "" {
			return ref.Alias
		}
		return name
	}
	if len(parts) == 1 {
		for c := ctes; c != nil; c = c.parent {
			if strings.EqualFold(c.name, parts[0]) {
				return derived(alias(parts[0]), renamed(c.columns, ref.ColumnAliases))
			}
		}
	}
	catalog, database, table := internalCatalog, r.database, parts[len(parts)-1]
	switch len(parts) {
	case 2:
		database = parts[0]
	case 3:
		catalog, database = parts[0], parts[1]
	}
	if len(parts) > 3 {
		// The grammar names at most catalog.db.table.
		return nil
	}
	rel := &relation{catalog: catalog, database: database, name: alias(table)}
	rel.qualified = database + "." + rel.name
	if !strings.EqualFold(catalog, internalCatalog) {
		rel.qualified = catalog + "." + rel.qualified
	}
	if !strings.EqualFold(catalog, internalCatalog) {
		// A table of another catalog: the metadata describes none.
		r.tables[ColumnResource{Database: database, Table: table}] = true
		rel.open, rel.source = true, func(string) SourceColumnSet { return nil }
		return rel
	}
	cols, ok := r.relation(database, table)
	if !ok {
		return nil
	}
	if cols == nil {
		// A system table, whose columns the metadata does not list.
		rel.open = true
		rel.source = func(column string) SourceColumnSet {
			if column == "" {
				// A star over it names no column.
				return nil
			}
			return SourceColumnSet{ColumnResource{Database: database, Table: table, Column: column}: true}
		}
		return rel
	}
	rel.columns = renamed(cols, ref.ColumnAliases)
	return rel
}

// derived returns a derived table or CTE of cols. One whose columns hold a
// star the metadata could not expand has columns no one can list, so any
// name resolves on it, to what those stars read.
func derived(name string, cols []spanColumn) *relation {
	rel := &relation{name: name, columns: cols}
	var starred SourceColumnSet
	for _, c := range cols {
		if c.star {
			rel.open = true
			if starred == nil {
				starred = SourceColumnSet{}
			}
			for k := range c.sources {
				starred[k] = true
			}
		}
	}
	if rel.open {
		rel.source = func(string) SourceColumnSet { return clone(starred) }
	}
	return rel
}

// relation returns the columns of a table or view, records it read, and
// reports false when the catalog lacks it. A system table the metadata
// does not hold has nil columns.
func (r *spanResolver) relation(database, table string) ([]spanColumn, bool) {
	db, err := r.opts.Catalog.GetDatabase(database)
	if err != nil || db == nil {
		if IsSystemDatabase(database, r.opts.IgnoreCaseSensitive) {
			r.tables[ColumnResource{Database: database, Table: table}] = true
			return nil, true
		}
		r.missing(NotFoundError{Database: &database})
		return nil, false
	}
	name := db.Name()
	if t := r.getTable(db, table); t != nil {
		r.tables[ColumnResource{Database: name, Table: t.Name}] = true
		cols := make([]spanColumn, 0, len(t.Columns))
		for _, c := range t.Columns {
			cols = append(cols, spanColumn{
				name:    c.Name,
				sources: SourceColumnSet{ColumnResource{Database: name, Table: t.Name, Column: c.Name}: true},
				plain:   true,
			})
		}
		return cols, true
	}
	if v := r.getView(db, table); v != nil {
		r.tables[ColumnResource{Database: name, Table: v.Name}] = true
		return r.view(name, v), true
	}
	if mv := r.getMaterializedView(db, table); mv != nil {
		// A materialized view stores what its definition reads: each
		// column follows its definition's columns, as a view's does, so
		// a masking policy on a base column covers it too.
		r.tables[ColumnResource{Database: name, Table: mv.Name}] = true
		return r.view(name, &metadata.ViewMetadata{Name: mv.Name, Definition: mv.Definition}), true
	}
	if IsSystemDatabase(database, r.opts.IgnoreCaseSensitive) {
		r.tables[ColumnResource{Database: name, Table: table}] = true
		return nil, true
	}
	r.missing(NotFoundError{Database: &name, Table: &table})
	return nil, false
}

func (r *spanResolver) getTable(db Database, name string) *metadata.TableMetadata {
	if t := db.GetTable(name); t != nil || !r.opts.IgnoreCaseSensitive {
		return t
	}
	return db.GetTable(strings.ToLower(name))
}

func (r *spanResolver) getMaterializedView(db Database, name string) *metadata.MaterializedViewMetadata {
	if v := db.GetMaterializedView(name); v != nil || !r.opts.IgnoreCaseSensitive {
		return v
	}
	return db.GetMaterializedView(strings.ToLower(name))
}

func (r *spanResolver) getView(db Database, name string) *metadata.ViewMetadata {
	if v := db.GetView(name); v != nil || !r.opts.IgnoreCaseSensitive {
		return v
	}
	return db.GetView(strings.ToLower(name))
}

// view returns a view's columns: each is the view's own column and the
// columns its definition reads for it, resolved in the view's database.
func (r *spanResolver) view(database string, v *metadata.ViewMetadata) []spanColumn {
	key := ColumnResource{Database: database, Table: v.Name}
	if cols, ok := r.views[key]; ok {
		// Its tables were recorded when it first resolved.
		return cols
	}
	var defined []spanColumn
	columns := v.Columns
	if !r.resolving[key] {
		r.resolving[key] = true
		file, err := parser.Parse(v.Definition)
		var query ast.Node
		if err == nil && file != nil && len(file.Stmts) == 1 {
			query = file.Stmts[0]
			// A materialized view's stored definition is its whole CREATE
			// statement; its query and column list are what it reads.
			if mv, ok := query.(*ast.CreateMTMVStmt); ok {
				query = mv.Query
				if len(columns) == 0 {
					for _, c := range mv.Columns {
						columns = append(columns, &metadata.ColumnMetadata{Name: c.Name})
					}
				}
			}
		}
		switch {
		case err != nil:
			r.fail(fmt.Errorf("read the definition of view %s.%s: %w", database, v.Name, parser.FirstError(err)))
		case query == nil || !isQuery(query):
			r.fail(fmt.Errorf("the definition of view %s.%s is not one query", database, v.Name))
		default:
			savedDatabase, savedText := r.database, r.text
			r.database, r.text = database, v.Definition
			defined = r.query(query, nil, nil)
			r.database, r.text = savedDatabase, savedText
		}
		delete(r.resolving, key)
	}
	var cols []spanColumn
	if len(columns) > 0 {
		for i, c := range columns {
			col := spanColumn{name: c.Name, plain: true}
			if i < len(defined) {
				col.sources = clone(defined[i].sources)
			}
			cols = append(cols, col)
		}
	} else {
		for _, c := range defined {
			cols = append(cols, spanColumn{name: c.name, sources: clone(c.sources), plain: true})
		}
	}
	for i := range cols {
		if cols[i].sources == nil {
			cols[i].sources = SourceColumnSet{}
		}
		cols[i].sources[ColumnResource{Database: database, Table: v.Name, Column: cols[i].name}] = true
	}
	r.views[key] = cols
	return cols
}

// star expands * or t.* to the columns of the relations it names.
func (r *spanResolver) star(item *ast.SelectItem, sc *spanScope) []spanColumn {
	var rels []*relation
	if item.TableName == nil || len(item.TableName.Parts) == 0 {
		rels = sc.relations
	} else {
		rels = matching(sc.relations, item.TableName.Parts, r)
		if len(rels) == 0 {
			table := strings.Join(item.TableName.Parts, ".")
			r.missing(NotFoundError{Table: &table})
			return nil
		}
	}
	except := map[string]bool{}
	for _, c := range item.ExceptColumns {
		except[strings.ToLower(c)] = true
	}
	var cols []spanColumn
	for _, rel := range rels {
		if rel.open && len(rel.columns) == 0 {
			// The metadata does not list its columns. What every column
			// reads, as a table function's arguments, the star reads too.
			cols = append(cols, spanColumn{name: "*", star: true, sources: rel.source("")})
			continue
		}
		for _, c := range rel.columns {
			if except[strings.ToLower(c.name)] {
				continue
			}
			cols = append(cols, spanColumn{name: c.name, sources: clone(c.sources), plain: !c.star, star: c.star})
		}
	}
	return cols
}

// matching returns the relations of rels a qualifier names: table,
// db.table or catalog.db.table. StarRocks matches a table's name and alias
// case-sensitively.
func matching(rels []*relation, qualifier []string, r *spanResolver) []*relation {
	table := qualifier[len(qualifier)-1]
	catalog, database := "", ""
	switch len(qualifier) {
	case 2:
		database = qualifier[0]
	case 3:
		catalog, database = qualifier[0], qualifier[1]
	}
	var out []*relation
	for _, rel := range rels {
		if !r.sameName(rel.name, table) {
			continue
		}
		if catalog != "" {
			// A derived relation is in the session's catalog.
			relCatalog := rel.catalog
			if relCatalog == "" {
				relCatalog = internalCatalog
			}
			if !strings.EqualFold(relCatalog, catalog) {
				continue
			}
		}
		if database != "" && !r.sameName(rel.database, database) && !(rel.database == "" && r.sameName(r.database, database)) {
			continue
		}
		out = append(out, rel)
	}
	return out
}

func (r *spanResolver) sameName(a, b string) bool {
	if r.opts.IgnoreCaseSensitive {
		return strings.EqualFold(a, b)
	}
	return a == b
}

// expr returns the columns an expression reads, resolving each column
// reference in sc and each subquery with sc as its enclosing scope.
func (r *spanResolver) expr(node ast.Node, sc *spanScope, ctes *cteBinding) SourceColumnSet {
	out := SourceColumnSet{}
	bound := map[string]int{}
	var visit func(n ast.Node)
	visit = func(n ast.Node) {
		ast.Inspect(n, func(n ast.Node) bool {
			switch e := n.(type) {
			case *ast.LambdaExpr:
				for _, p := range e.Params {
					bound[strings.ToLower(p)]++
				}
				visit(e.Body)
				for _, p := range e.Params {
					bound[strings.ToLower(p)]--
				}
				return false
			case *ast.ColumnRef:
				if e.Name == nil || len(e.Name.Parts) == 0 {
					return false
				}
				// A lambda parameter, or a field path on one (x -> x.f),
				// shadows any column or table of its name.
				if bound[strings.ToLower(e.Name.Parts[0])] > 0 {
					return false
				}
				for c := range r.column(e.Name.Parts, sc) {
					out[c] = true
				}
				return false
			case *ast.SubqueryExpr:
				for _, c := range r.subquery(e, sc, ctes) {
					for s := range c.sources {
						out[s] = true
					}
				}
				return false
			case *ast.ExistsExpr:
				if e.Subquery != nil {
					for _, c := range r.subquery(e.Subquery, sc, ctes) {
						for s := range c.sources {
							out[s] = true
						}
					}
				}
				return false
			}
			return true
		})
	}
	visit(node)
	return out
}

// column resolves a column reference, innermost scope first. A name more
// than one relation of a scope has (a USING column, or one StarRocks refuses
// as ambiguous) reads all of them.
func (r *spanResolver) column(parts []string, sc *spanScope) SourceColumnSet {
	// The parser holds a struct field path, s.f or t.s.f, as a longer
	// column reference: when the whole name is no column, read a prefix as
	// the column and the rest as its fields.
	for n := len(parts); n >= 1; n-- {
		if found, ok := r.lookup(parts[:n], sc); ok {
			return found
		}
	}
	name := parts[len(parts)-1]
	qualifier := parts[:len(parts)-1]
	column := name
	if len(qualifier) > 0 {
		table := qualifier[len(qualifier)-1]
		r.missing(NotFoundError{Table: &table, Column: &column})
	} else {
		r.missing(NotFoundError{Column: &column})
	}
	return nil
}

// lookup resolves a column reference, innermost scope first, and reports
// whether any relation has it.
func (r *spanResolver) lookup(parts []string, sc *spanScope) (SourceColumnSet, bool) {
	name := parts[len(parts)-1]
	qualifier := parts[:len(parts)-1]
	for s := sc; s != nil; s = s.parent {
		rels := s.relations
		if len(qualifier) > 0 {
			rels = matching(rels, qualifier, r)
		}
		found := SourceColumnSet{}
		hit := false
		for _, rel := range rels {
			for _, c := range rel.columns {
				if strings.EqualFold(c.name, name) {
					hit = true
					for k := range c.sources {
						found[k] = true
					}
				}
			}
		}
		if !hit {
			for _, rel := range rels {
				if rel.open {
					hit = true
					for k := range rel.source(name) {
						found[k] = true
					}
				}
			}
		}
		if !hit && len(qualifier) == 0 {
			for _, c := range s.outputs {
				if strings.EqualFold(c.name, name) {
					hit = true
					for k := range c.sources {
						found[k] = true
					}
				}
			}
		}
		if hit {
			return found, true
		}
		if len(qualifier) > 0 && len(rels) > 0 {
			break
		}
	}
	return nil, false
}

// isColumnRef reports whether expr is a bare column reference, in any
// number of parentheses.
func isColumnRef(expr ast.Node) bool {
	for {
		switch e := expr.(type) {
		case *ast.ParenExpr:
			expr = e.Expr
		case *ast.ColumnRef:
			return true
		default:
			return false
		}
	}
}

// resultName is the name StarRocks gives a SELECT item's column: its
// alias; a column reference's name as written, in any parentheses; or else
// the expression as StarRocks prints it (printed).
func (r *spanResolver) resultName(item *ast.SelectItem) string {
	if item.Aliased {
		return item.Alias
	}
	expr := item.Expr
	for {
		p, ok := expr.(*ast.ParenExpr)
		if !ok {
			break
		}
		expr = p.Expr
	}
	if e, ok := expr.(*ast.ColumnRef); ok && e.Name != nil && len(e.Name.Parts) > 0 {
		return e.Name.Parts[len(e.Name.Parts)-1]
	}
	return r.printed(item.Expr)
}

// printed renders an expression as StarRocks prints an unnamed column: a
// string literal in single quotes, a binary operator between single
// spaces, and a negation as -1 times its operand. Any other expression
// keeps its text as written.
func (r *spanResolver) printed(node ast.Node) string {
	switch e := node.(type) {
	case *ast.Literal:
		if e.Kind == ast.LitString {
			return "'" + strings.ReplaceAll(unquote(e.Value), "'", "\\'") + "'"
		}
	case *ast.BinaryExpr:
		return r.printed(e.Left) + " " + e.Op.String() + " " + r.printed(e.Right)
	case *ast.FuncCallExpr:
		if w := e.Over; w != nil && w.Name == "" && w.Frame == nil {
			return r.printedWindow(e, w)
		}
		if r.qualify != nil {
			if name, ok := r.printedCall(e); ok {
				return name
			}
		}
	case *ast.ParenExpr:
		if r.qualify != nil {
			return r.printed(e.Expr)
		}
	case *ast.ColumnRef:
		if r.qualify != nil {
			if name, ok := r.qualifiedColumn(e, r.qualify); ok {
				return name
			}
		}
	case *ast.UnaryExpr:
		if e.Op == ast.UnaryMinus {
			if lit, ok := e.Expr.(*ast.Literal); ok && (lit.Kind == ast.LitInt || lit.Kind == ast.LitFloat) {
				return "-" + lit.Value
			}
			return "-1 * " + r.printed(e.Expr)
		}
	}
	return r.slice(node)
}

// qualifiedColumn renders a column reference as StarRocks prints it once
// analyzed (AstToStringBuilder.visitSlot): its relation's qualified name
// or else its alias or name, a dot, and the column's name as written. It
// reports false unless exactly one relation of sc lists the column.
func (r *spanResolver) qualifiedColumn(ref *ast.ColumnRef, sc *spanScope) (string, bool) {
	if ref.Name == nil || len(ref.Name.Parts) == 0 || len(ref.Name.Parts) > 4 {
		return "", false
	}
	parts := ref.Name.Parts
	name := parts[len(parts)-1]
	rels := sc.relations
	if len(parts) > 1 {
		rels = matching(rels, parts[:len(parts)-1], r)
	}
	var found *relation
	for _, rel := range rels {
		for _, c := range rel.columns {
			if strings.EqualFold(c.name, name) {
				if found != nil && found != rel {
					return "", false
				}
				found = rel
			}
		}
	}
	if found == nil {
		return "", false
	}
	prefix := found.qualified
	if prefix == "" {
		prefix = found.name
	}
	return prefix + "." + name, true
}

// printedWindow renders an analytic call as StarRocks names it: the call,
// then OVER with PARTITION BY and ORDER BY each followed by a space, and
// every sort key with its direction. A frame or a NULLS order keeps the
// text as written.
func (r *spanResolver) printedWindow(call *ast.FuncCallExpr, w *ast.WindowSpec) string {
	for _, o := range w.OrderBy {
		if o == nil || o.NullsFirst != nil {
			return r.slice(call)
		}
	}
	if call.Loc.Start < 0 || w.Loc.Start > len(r.text) || call.Loc.Start > w.Loc.Start {
		return r.slice(call)
	}
	fn := strings.TrimSpace(r.text[call.Loc.Start:w.Loc.Start])
	if n := len(fn); n >= 4 && strings.EqualFold(fn[n-4:], "OVER") {
		fn = strings.TrimSpace(fn[:n-4])
	}
	var b strings.Builder
	b.WriteString(fn + " OVER (")
	if len(w.PartitionBy) > 0 {
		var parts []string
		for _, p := range w.PartitionBy {
			parts = append(parts, r.printed(p))
		}
		b.WriteString("PARTITION BY " + strings.Join(parts, ", ") + " ")
	}
	if len(w.OrderBy) > 0 {
		var keys []string
		for _, o := range w.OrderBy {
			dir := " ASC"
			if o.Desc {
				dir = " DESC"
			}
			keys = append(keys, r.printed(o.Expr)+dir)
		}
		b.WriteString("ORDER BY " + strings.Join(keys, ", ") + " ")
	}
	b.WriteString(")")
	return b.String()
}

// slice returns the text a node was parsed from.
func (r *spanResolver) slice(node ast.Node) string {
	loc := ast.NodeLoc(node)
	if loc.Start < 0 || loc.End > len(r.text) || loc.Start > loc.End {
		return ""
	}
	return r.text[loc.Start:loc.End]
}

// unquote returns a string literal's value.
func unquote(s string) string {
	if len(s) < 2 || (s[0] != '\'' && s[0] != '"') || s[len(s)-1] != s[0] {
		return s
	}
	q := s[0]
	body := s[1 : len(s)-1]
	var b strings.Builder
	for i := 0; i < len(body); i++ {
		c := body[i]
		switch {
		case c == q && i+1 < len(body) && body[i+1] == q:
			b.WriteByte(q)
			i++
		case c == '\\' && i+1 < len(body):
			i++
			switch body[i] {
			case 'n':
				b.WriteByte('\n')
			case 't':
				b.WriteByte('\t')
			case 'r':
				b.WriteByte('\r')
			case '0':
				b.WriteByte(0)
			default:
				b.WriteByte(body[i])
			}
		default:
			b.WriteByte(c)
		}
	}
	return b.String()
}

// resolvedResults converts resolved columns to the span's results.
func resolvedResults(cols []spanColumn) []QuerySpanResult {
	out := make([]QuerySpanResult, 0, len(cols))
	for _, c := range cols {
		sources := c.sources
		if sources == nil {
			sources = SourceColumnSet{}
		}
		out = append(out, QuerySpanResult{
			Name:           c.name,
			SourceColumns:  sources,
			IsPlainField:   c.plain,
			SelectAsterisk: c.star,
		})
	}
	return out
}

// syntacticResults converts the columns GetQuerySpan names without a
// catalog: each column reference as written, qualified as far as the text
// qualifies it.
func syntacticResults(cols []ColumnInfo) []QuerySpanResult {
	out := make([]QuerySpanResult, 0, len(cols))
	for _, c := range cols {
		sources := SourceColumnSet{}
		for _, ref := range c.SourceColumns {
			sources[ColumnResource{Database: ref.Database, Table: ref.Table, Column: ref.Column}] = true
		}
		out = append(out, QuerySpanResult{Name: c.Name, SourceColumns: sources, SelectAsterisk: c.Star})
	}
	return out
}

// clone copies a column set; nil stays nil.
func clone(s SourceColumnSet) SourceColumnSet {
	if s == nil {
		return nil
	}
	out := make(SourceColumnSet, len(s))
	for k := range s {
		out[k] = true
	}
	return out
}
