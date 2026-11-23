package analysis

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/doris/ast"
	"github.com/bytebase/omni/doris/parser"
	"github.com/bytebase/omni/metadata"
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

// Database is one database's synced metadata. Doris keeps its tables and
// views directly in the database.
type Database interface {
	// Name is the database's name as the metadata spells it.
	Name() string
	// GetTable returns the named table, or nil.
	GetTable(name string) *metadata.TableMetadata
	// GetView returns the named view, or nil.
	GetView(name string) *metadata.ViewMetadata
}

// internalCatalog is the catalog that holds Doris's own databases; a
// table of any other catalog lives in an external source the metadata
// does not describe.
const internalCatalog = "internal"

// spanColumn is one column a relation or a query produces.
type spanColumn struct {
	name    string
	sources SourceColumnSet
	plain   bool
}

// relation is a table, view, CTE, derived table or table function in a
// query's FROM clause.
type relation struct {
	// database is the database of a table or view, which a column
	// reference may name (db.t.c); empty for a derived relation.
	database string
	// name is the alias, or the table's name as written.
	name    string
	columns []spanColumn
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
	case *ast.SelectStmt, *ast.SetOpStmt, *ast.ValuesStmt, *ast.GroupedQuery:
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
	case *ast.ValuesStmt:
		return r.values(n, outer, ctes, false)
	case *ast.GroupedQuery:
		cols := r.query(n.Query, outer, ctes)
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

// leftmostQuery returns the SELECT or VALUES a set operation's left spine
// starts with, and its WITH clause, which scopes over every arm.
func leftmostQuery(node ast.Node) (ast.Node, *ast.WithClause) {
	for {
		switch n := node.(type) {
		case *ast.SetOpStmt:
			node = n.Left
		case *ast.SelectStmt:
			return n, n.With
		case *ast.ValuesStmt:
			return n, n.With
		default:
			return nil, nil
		}
	}
}

func (r *spanResolver) setOp(n *ast.SetOpStmt, outer *spanScope, ctes *cteBinding) []spanColumn {
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
		case *ast.ValuesStmt:
			armCols = r.values(a, outer, ctes, arm == first)
		default:
			armCols = r.query(arm, outer, ctes)
		}
		if i == 0 {
			cols = make([]spanColumn, len(armCols))
			for j, c := range armCols {
				cols[j] = spanColumn{name: c.name, sources: SourceColumnSet{}}
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
			}
		}
	}
	r.trailing(n.OrderBy, n.Limit, n.Offset, &spanScope{outputs: cols, parent: outer}, ctes)
	return cols
}

// values resolves a VALUES query. withBound is set when a set operation
// has bound its WITH clause already.
func (r *spanResolver) values(n *ast.ValuesStmt, outer *spanScope, ctes *cteBinding, withBound bool) []spanColumn {
	if !withBound {
		ctes = r.with(n.With, outer, ctes)
	}
	sc := &spanScope{parent: outer}
	var cols []spanColumn
	for i, row := range n.Rows {
		for j, item := range row {
			if item == nil || item.Expr == nil {
				continue
			}
			sources := r.expr(item.Expr, sc, ctes)
			if i == 0 {
				// A VALUES row names its columns by their text, literals
				// with their quotes.
				name := item.Alias
				if !item.Aliased {
					name = r.slice(item.Expr)
				}
				cols = append(cols, spanColumn{name: name, sources: sources})
			} else if j < len(cols) {
				for k := range sources {
					cols[j].sources[k] = true
				}
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

// from adds a FROM item's relations to sc. A derived table sees only the
// enclosing queries' scopes (outer): Doris has no LATERAL.
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
		rel := r.tableRef(n, outer, ctes)
		if rel == nil {
			return
		}
		sc.relations = append(sc.relations, rel)
		for _, lv := range n.LateralViews {
			if lv == nil {
				continue
			}
			var sources SourceColumnSet
			if lv.Func != nil {
				sources = r.expr(lv.Func, sc, ctes)
			}
			view := &relation{name: lv.TableAlias}
			for _, c := range lv.Columns {
				view.columns = append(view.columns, spanColumn{name: c, sources: clone(sources)})
			}
			sc.relations = append(sc.relations, view)
		}
	}
}

// tableRef resolves one FROM item to a relation, or nil when the catalog
// lacks what it names.
func (r *spanResolver) tableRef(ref *ast.TableRef, outer *spanScope, ctes *cteBinding) *relation {
	if ref.Func != nil {
		// A table function reads no table; its columns are its own.
		r.expr(ref.Func, &spanScope{parent: outer}, ctes)
		name := ref.Alias
		if name == "" && ref.Func.Name != nil && len(ref.Func.Name.Parts) > 0 {
			name = ref.Func.Name.Parts[len(ref.Func.Name.Parts)-1]
		}
		return &relation{name: name, open: true, source: func(string) SourceColumnSet { return nil }}
	}
	if ref.Subquery != nil {
		cols := r.subquery(ref.Subquery, outer, ctes)
		return &relation{name: ref.Alias, columns: renamed(cols, ref.ColumnAliases)}
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
				return &relation{name: alias(parts[0]), columns: renamed(c.columns, ref.ColumnAliases)}
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
	rel := &relation{database: database, name: alias(table)}
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
			return SourceColumnSet{ColumnResource{Database: database, Table: table, Column: column}: true}
		}
		return rel
	}
	rel.columns = renamed(cols, ref.ColumnAliases)
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
	if !r.resolving[key] {
		r.resolving[key] = true
		file, err := parser.Parse(v.Definition)
		switch {
		case err != nil:
			r.fail(fmt.Errorf("read the definition of view %s.%s: %w", database, v.Name, parser.FirstError(err)))
		case file == nil || len(file.Stmts) != 1 || !isQuery(file.Stmts[0]):
			r.fail(fmt.Errorf("the definition of view %s.%s is not one query", database, v.Name))
		default:
			savedDatabase, savedText := r.database, r.text
			r.database, r.text = database, v.Definition
			defined = r.query(file.Stmts[0], nil, nil)
			r.database, r.text = savedDatabase, savedText
		}
		delete(r.resolving, key)
	}
	var cols []spanColumn
	if len(v.Columns) > 0 {
		for i, c := range v.Columns {
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
		if rel.open {
			// The metadata does not list its columns.
			cols = append(cols, spanColumn{name: "*"})
			continue
		}
		for _, c := range rel.columns {
			if except[strings.ToLower(c.name)] {
				continue
			}
			cols = append(cols, spanColumn{name: c.name, sources: clone(c.sources), plain: true})
		}
	}
	return cols
}

// matching returns the relations of rels a qualifier names: table,
// db.table or catalog.db.table. Doris matches a table's name and alias
// case-sensitively.
func matching(rels []*relation, qualifier []string, r *spanResolver) []*relation {
	table := qualifier[len(qualifier)-1]
	database := ""
	switch len(qualifier) {
	case 2:
		database = qualifier[0]
	case 3:
		if !strings.EqualFold(qualifier[0], internalCatalog) {
			return nil
		}
		database = qualifier[1]
	}
	var out []*relation
	for _, rel := range rels {
		if !r.sameName(rel.name, table) {
			continue
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
				if len(e.Name.Parts) == 1 && bound[strings.ToLower(e.Name.Parts[0])] > 0 {
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
// than one relation of a scope has (a USING column, or one Doris refuses
// as ambiguous) reads all of them.
func (r *spanResolver) column(parts []string, sc *spanScope) SourceColumnSet {
	name := parts[len(parts)-1]
	qualifier := parts[:len(parts)-1]
	if len(qualifier) == 3 {
		if !strings.EqualFold(qualifier[0], internalCatalog) {
			return nil
		}
	}
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
			return found
		}
		if len(qualifier) > 0 && len(rels) > 0 {
			break
		}
	}
	column := name
	if len(qualifier) > 0 {
		table := qualifier[len(qualifier)-1]
		r.missing(NotFoundError{Table: &table, Column: &column})
	} else {
		r.missing(NotFoundError{Column: &column})
	}
	return nil
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

// resultName is the name Doris gives a SELECT item's column: its alias; a
// column reference's name as written, in any parentheses; a string
// literal's value; or else the expression's text as written.
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
	switch e := expr.(type) {
	case *ast.ColumnRef:
		if e.Name != nil && len(e.Name.Parts) > 0 {
			return e.Name.Parts[len(e.Name.Parts)-1]
		}
	case *ast.Literal:
		if e.Kind == ast.LitString {
			return unquote(e.Value)
		}
	}
	return r.slice(item.Expr)
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
			SelectAsterisk: c.name == "*" && len(sources) == 0,
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
		out = append(out, QuerySpanResult{Name: c.Name, SourceColumns: sources, SelectAsterisk: c.Name == "*" || strings.HasSuffix(c.Name, ".*")})
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
