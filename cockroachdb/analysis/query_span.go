package analysis

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/pgwire/pgcode"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/pgwire/pgerror"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
)

// ResolveOptions configures ResolveQuerySpan.
type ResolveOptions struct {
	// DefaultDatabase is the database the statement runs in.
	DefaultDatabase string
	// SearchPath resolves unqualified names. Empty means public.
	SearchPath []string
	// Catalog reads the instance's synced metadata. GetDatabase must be set.
	Catalog Catalog
}

// Catalog is the synced metadata of the instance a statement runs on.
type Catalog struct {
	// GetDatabase returns a database's metadata, or a nil Database when the
	// database has not been synced.
	GetDatabase func(name string) (Database, error)
}

// Database is one database's synced metadata. The lookups follow the
// metadata's own name matching, so the caller decides case sensitivity.
type Database interface {
	// Metadata is the whole snapshot, which holds the views' definitions.
	Metadata() *metadata.DatabaseSchemaMetadata
	// GetTable returns the table named table in schema, or nil.
	GetTable(schema, table string) *metadata.TableMetadata
	// SearchRelation returns the schema and name of the first table, view or
	// materialized view named name in searchPath's schemas, or two empty
	// strings.
	SearchRelation(searchPath []string, name string) (schema string, relation string)
}

// SyntaxError reports a statement CockroachDB's parser rejects.
type SyntaxError struct {
	Err error
	// Position is the byte offset of the failure in the statement, or -1.
	Position int
}

func (e *SyntaxError) Error() string { return e.Err.Error() }

func (e *SyntaxError) Unwrap() error { return e.Err }

// ErrMixUserSystemTables reports a statement that reads user and system
// tables together.
var ErrMixUserSystemTables = analysis.ErrMixUserSystemTables

// ResolveQuerySpan returns the column lineage of one CockroachDB statement,
// resolved against the synced metadata on CockroachDB's own AST. A column read
// through a view names the view's column and the columns behind it, as
// CockroachDB's result-set metadata follows a view to its table.
//
// A statement that does not parse returns a *SyntaxError, one that reads user
// and system tables together ErrMixUserSystemTables, and a query that names a
// column its relations lack or name ambiguously CockroachDB's error for it,
// with its SQLSTATE (pgerror.GetPGCode). A relation missing from the metadata
// sets the span's NotFoundError. Statements other than a SELECT of user tables
// return their QueryType and the relations they read, with no Results.
func ResolveQuerySpan(_ context.Context, statement string, opts ResolveOptions) (*analysis.ResolvedSpan, error) {
	stmts, err := Parse(statement)
	if err != nil {
		return nil, &SyntaxError{Err: err, Position: extractParseErrorInfo(statement, err).Position}
	}
	if len(stmts) != 1 {
		return nil, fmt.Errorf("expected exactly 1 statement, got %d", len(stmts))
	}
	searchPath := opts.SearchPath
	if len(searchPath) == 0 {
		searchPath = []string{"public"}
	}
	r := &spanResolver{
		defaultDatabase: opts.DefaultDatabase,
		searchPath:      searchPath,
		catalog:         opts.Catalog,
		databases:       map[string]Database{},
		accesses:        analysis.ColumnSet{},
		predicates:      analysis.ColumnSet{},
		expanding:       map[analysis.ColumnResource]bool{},
	}
	stmt := stmts[0].AST
	explain, explainAnalyze := stmt.(*tree.ExplainAnalyze)
	if explainAnalyze {
		stmt = explain.Statement
	}
	queryType := classifyStatement(stmt)
	var results []spanColumn
	if sel, ok := stmt.(*tree.Select); ok {
		results = r.selectQuery(sel, nil)
	} else {
		r.statementTables(stmt, nil)
	}
	if r.metadataErr != nil {
		return nil, r.metadataErr
	}

	allSystem, mixed := r.systemAccess()
	if mixed {
		return nil, ErrMixUserSystemTables
	}
	if queryType == analysis.Select && allSystem {
		queryType = analysis.SelectInfoSchema
	}
	span := &analysis.ResolvedSpan{
		Type:           queryType,
		ExplainAnalyze: explainAnalyze,
		SourceColumns:  r.accesses,
		Results:        []analysis.QuerySpanResult{},
	}
	if queryType != analysis.Select {
		// The span of any other statement is the relations it reads, which
		// the walk records past a name it cannot resolve, so such a name
		// leaves it whole.
		if explainAnalyze {
			// EXPLAIN ANALYZE runs the statement and reports the rows each
			// condition kept, so its conditions count.
			span.PredicateColumns = r.predicates
			span.UnresolvedColumnsError = r.unresolvedColumns()
		}
		return span, nil
	}
	if r.notFound != nil {
		span.NotFoundError = r.notFound
		return span, nil
	}
	if r.resolveErr != nil {
		return nil, r.resolveErr
	}
	if !explainAnalyze {
		for _, c := range results {
			span.Results = append(span.Results, analysis.QuerySpanResult{
				Name:           c.name,
				SourceColumns:  c.sources,
				IsPlainField:   c.plain,
				SelectAsterisk: c.asterisk,
			})
		}
	}
	span.PredicateColumns = r.predicates
	span.UnresolvedColumnsError = r.unresolvedColumns()
	return span, nil
}

// classifyStatement returns the permission class of stmt. A SELECT whose CTEs
// or statement sources write is DML.
func classifyStatement(stmt tree.Statement) analysis.QueryType {
	switch n := stmt.(type) {
	case *tree.Select, *tree.ParenSelect:
		if len(getMutationTargets(n)) > 0 {
			return analysis.DML
		}
		return analysis.Select
	case *tree.Insert, *tree.Update, *tree.Delete, *tree.CopyFrom, *tree.CopyTo, *tree.Import, *tree.Call:
		return analysis.DML
	case *tree.Explain:
		return analysis.Explain
	case *tree.SetVar:
		return analysis.Select
	case *tree.DoBlock:
		return analysis.DDL
	}
	if strings.HasPrefix(stmt.StatementTag(), "SHOW") {
		return analysis.SelectInfoSchema
	}
	switch stmt.StatementType() {
	case tree.TypeDDL, tree.TypeDCL:
		return analysis.DDL
	}
	return analysis.QueryTypeUnknown
}

// systemAccess reports whether every relation the statement reads is a system
// one, and whether it reads system and user relations together.
func (r *spanResolver) systemAccess() (allSystem, mixed bool) {
	hasSystem, hasUser := r.readsSystem, false
	for resource := range r.accesses {
		if isSystemResource(resource) {
			hasSystem = true
		} else {
			hasUser = true
		}
	}
	return hasSystem && !hasUser, hasSystem && hasUser
}

func isSystemResource(resource analysis.ColumnResource) bool {
	return IsSystemDatabase(resource.Database) || IsSystemSchema(resource.Schema)
}

// unresolvedColumns reports the tables read whose synced metadata has no
// columns, so their lineage cannot be trusted.
func (r *spanResolver) unresolvedColumns() *analysis.UnresolvedColumnsError {
	var relations []analysis.ColumnResource
	for resource := range r.accesses {
		if isSystemResource(resource) {
			continue
		}
		db := r.database(resource.Database)
		if db == nil {
			continue
		}
		if t := db.GetTable(resource.Schema, resource.Table); t != nil && len(t.GetColumns()) == 0 {
			relations = append(relations, resource)
		}
	}
	if len(relations) == 0 {
		return nil
	}
	return &analysis.UnresolvedColumnsError{Relations: relations}
}

// spanResolver resolves one statement's lineage. It records the relations the
// statement names and the columns that decide which rows it returns.
//
// A name it cannot resolve does not stop it: it records the first such error
// in resolveErr and walks on, so every relation the statement names is still
// recorded, and the caller decides what the error means for the statement.
type spanResolver struct {
	defaultDatabase string
	searchPath      []string
	catalog         Catalog
	databases       map[string]Database

	accesses   analysis.ColumnSet
	predicates analysis.ColumnSet
	// readsSystem is set by a SHOW statement source, which reads system
	// metadata without naming a relation.
	readsSystem bool
	notFound    *analysis.ResourceNotFoundError
	resolveErr  error
	// metadataErr is a failure to read the metadata, which fails the
	// statement whatever it is.
	metadataErr error
	// inView counts the view definitions being resolved: what a view reads is
	// not an access of the statement, and its conditions are the view's.
	inView    int
	expanding map[analysis.ColumnResource]bool
}

func (r *spanResolver) fail(err error) {
	if r.resolveErr == nil {
		r.resolveErr = err
	}
}

// spanColumn is one output column of a relation or a query.
type spanColumn struct {
	name    string
	sources analysis.ColumnSet
	// hidden marks a column * leaves out, such as a table's rowid.
	hidden bool
	// plain marks a bare column reference.
	plain bool
	// asterisk marks a * over a relation whose columns are unknown.
	asterisk bool
}

// spanRelation is a relation a qualified name can reach: a table, a view, a
// CTE, or an aliased FROM item.
type spanRelation struct {
	name, schema, database string
	columns                []spanColumn
	// opaque marks a relation whose columns are not all known, such as a set
	// returning function or a table the metadata has no columns for: a name
	// it does not list reads opaqueReads, and the column opaqueTable names.
	opaque      bool
	opaqueReads analysis.ColumnSet
	opaqueTable *analysis.ColumnResource
}

func (rel *spanRelation) opaqueColumn(name string) spanColumn {
	sources := union(rel.opaqueReads)
	if rel.opaqueTable != nil {
		c := *rel.opaqueTable
		c.Column = name
		sources[c] = true
	}
	return spanColumn{name: name, sources: sources}
}

// fromItem is one FROM item: the relations its qualified names reach and its
// output columns, which unqualified names and * read.
type fromItem struct {
	relations []*spanRelation
	columns   []spanColumn
}

func (item *fromItem) opaque() []*spanRelation {
	var out []*spanRelation
	for _, rel := range item.relations {
		if rel.opaque {
			out = append(out, rel)
		}
	}
	return out
}

// spanScope is the names a part of a query sees: its FROM items and CTEs, then
// its parent's.
type spanScope struct {
	parent *spanScope
	items  []*fromItem
	ctes   map[string][]spanColumn
	// from marks the scope of a FROM clause's items, which only a LATERAL
	// item or a function in FROM sees from inside FROM.
	from bool
}

func (sc *spanScope) cte(name string) ([]spanColumn, bool) {
	for s := sc; s != nil; s = s.parent {
		if cols, ok := s.ctes[name]; ok {
			return cols, true
		}
	}
	return nil, false
}

// outside returns the scope a FROM item that is not LATERAL sees: sc without
// the FROM items around it.
func (sc *spanScope) outside() *spanScope {
	for sc != nil && sc.from {
		sc = sc.parent
	}
	return sc
}

// outputScope is a scope in which a name may be an output column of the
// query, as in ORDER BY, before the names of inner.
func outputScope(cols []spanColumn, inner *spanScope) *spanScope {
	return &spanScope{parent: inner, items: []*fromItem{{columns: cols}}}
}

func (r *spanResolver) database(name string) Database {
	if db, ok := r.databases[name]; ok {
		return db
	}
	if r.catalog.GetDatabase == nil {
		r.metadataErr = errors.New("ResolveOptions.Catalog.GetDatabase is not set")
		return nil
	}
	db, err := r.catalog.GetDatabase(name)
	if err != nil {
		if r.metadataErr == nil {
			r.metadataErr = fmt.Errorf("failed to get database metadata for database: %s: %w", name, err)
		}
		return nil
	}
	r.databases[name] = db
	return db
}

func (r *spanResolver) access(resource analysis.ColumnResource) {
	if r.inView == 0 {
		r.accesses[resource] = true
	}
}

func (r *spanResolver) addPredicates(sources analysis.ColumnSet) {
	if r.inView > 0 {
		return
	}
	for c := range sources {
		r.predicates[c] = true
	}
}

// selectQuery resolves a query with its WITH, ORDER BY and LIMIT clauses.
// FOR UPDATE OF names relations of the query's own FROM, so it reads nothing.
func (r *spanResolver) selectQuery(sel *tree.Select, outer *spanScope) []spanColumn {
	sc := outer
	if sel.With != nil {
		sc = r.with(sel.With, outer)
	}
	cols, inner := r.selectStatement(sel.Select, sc)
	r.orderBy(sel.OrderBy, outputScope(cols, inner))
	r.limit(sel.Limit, sc)
	return cols
}

func (r *spanResolver) orderBy(orderBy tree.OrderBy, sc *spanScope) {
	for _, order := range orderBy {
		r.sources(order.Expr, sc)
	}
}

// limit records the columns LIMIT and OFFSET read, which decide how many rows
// come back.
func (r *spanResolver) limit(limit *tree.Limit, sc *spanScope) {
	if limit != nil {
		r.addPredicates(r.sources(limit.Count, sc))
		r.addPredicates(r.sources(limit.Offset, sc))
	}
}

// with resolves a WITH clause into a scope holding its CTEs. Each CTE sees the
// ones before it, and in WITH RECURSIVE itself too.
func (r *spanResolver) with(w *tree.With, outer *spanScope) *spanScope {
	sc := &spanScope{parent: outer, ctes: map[string][]spanColumn{}}
	for _, cte := range w.CTEList {
		name := string(cte.Name.Alias)
		sel, ok := cte.Stmt.(*tree.Select)
		if !ok {
			// A data-modifying CTE makes the statement DML, which has no
			// lineage.
			r.statementTables(cte.Stmt, sc)
			sc.ctes[name] = nil
			continue
		}
		var cols []spanColumn
		if union, ok := sel.Select.(*tree.UnionClause); ok && w.Recursive {
			cols = r.recursiveCTE(name, cte.Name.Cols, union, sc)
		} else {
			cols = r.selectQuery(sel, sc)
		}
		sc.ctes[name] = r.renameColumns(cols, cte.Name.Cols, nil)
	}
	return sc
}

// recursiveCTE resolves a recursive CTE: its first branch, then its second
// with the CTE bound to the columns so far, until the sources stop growing.
func (r *spanResolver) recursiveCTE(name string, aliases tree.ColumnDefList, union *tree.UnionClause, sc *spanScope) []spanColumn {
	cols := r.selectQuery(union.Left, sc)
	for {
		sc.ctes[name] = r.renameColumns(cols, aliases, nil)
		merged := r.mergeBranches(cols, r.selectQuery(union.Right, sc))
		if sameSources(cols, merged) {
			return merged
		}
		cols = merged
	}
}

func sameSources(a, b []spanColumn) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i].sources) != len(b[i].sources) {
			return false
		}
	}
	return true
}

// mergeBranches is the columns of a set operation: the left branch's names
// and every branch's sources.
func (r *spanResolver) mergeBranches(left, right []spanColumn) []spanColumn {
	if len(left) != len(right) {
		r.fail(fmt.Errorf("each branch of a set operation must have the same number of columns: %d and %d", len(left), len(right)))
	}
	out := stripColumns(left)
	for i := range out {
		if i < len(right) {
			out[i].sources = union(out[i].sources, right[i].sources)
		}
	}
	return out
}

func union(sets ...analysis.ColumnSet) analysis.ColumnSet {
	out := analysis.ColumnSet{}
	for _, s := range sets {
		for c := range s {
			out[c] = true
		}
	}
	return out
}

// selectStatement resolves the body of a query. It returns the query's
// columns and the scope its ORDER BY sees beneath the output columns: the
// FROM items of a plain SELECT, or sc.
func (r *spanResolver) selectStatement(stmt tree.SelectStatement, sc *spanScope) ([]spanColumn, *spanScope) {
	switch n := stmt.(type) {
	case *tree.SelectClause:
		return r.selectClause(n, sc)
	case *tree.ParenSelect:
		return r.selectQuery(n.Select, sc), sc
	case *tree.UnionClause:
		left := r.selectQuery(n.Left, sc)
		right := r.selectQuery(n.Right, sc)
		if n.Type == tree.ExceptOp {
			// EXCEPT returns rows of the left branch that the right branch
			// lacks, so only the left supplies values; both decide which rows
			// come back.
			for _, c := range append(left, right...) {
				r.addPredicates(c.sources)
			}
			return stripColumns(left), sc
		}
		// A row INTERSECT returns holds values of both branches, so a masking
		// policy on either applies, though CockroachDB names the left alone.
		if n.Type == tree.IntersectOp {
			for _, c := range append(left, right...) {
				r.addPredicates(c.sources)
			}
		}
		return r.mergeBranches(left, right), sc
	case *tree.ValuesClause:
		var cols []spanColumn
		for _, row := range n.Rows {
			for i, e := range row {
				if i == len(cols) {
					cols = append(cols, spanColumn{name: fmt.Sprintf("column%d", i+1), sources: analysis.ColumnSet{}})
				}
				cols[i].sources = union(cols[i].sources, r.sources(e, sc))
			}
		}
		return cols, sc
	}
	r.fail(fmt.Errorf("unsupported query %T", stmt))
	return nil, sc
}

// stripColumns keeps the names and sources of a query's columns, so a
// relation over the query reads none as a bare reference.
func stripColumns(cols []spanColumn) []spanColumn {
	out := make([]spanColumn, len(cols))
	for i, c := range cols {
		out[i] = spanColumn{name: c.name, sources: c.sources}
	}
	return out
}

// selectClause resolves a SELECT and returns its columns and its scope.
func (r *spanResolver) selectClause(n *tree.SelectClause, outer *spanScope) ([]spanColumn, *spanScope) {
	sc := &spanScope{parent: outer, from: true}
	for _, te := range n.From.Tables {
		sc.items = append(sc.items, r.tableExpr(te, sc))
	}
	sc.from = false
	r.sources(n.From.AsOf.Expr, outer)
	var cols []spanColumn
	for _, target := range n.Exprs {
		cols = append(cols, r.selectTarget(target, sc)...)
	}
	if n.Where != nil {
		r.addPredicates(r.sources(n.Where.Expr, sc))
	}
	if n.Having != nil {
		r.addPredicates(r.sources(n.Having.Expr, sc))
	}
	// Grouping, DISTINCT ON and window definitions add no lineage and no
	// condition, but their subqueries read relations. GROUP BY reads the FROM
	// items before an output column of the name; DISTINCT ON, as ORDER BY,
	// the output column first.
	grouping := &spanScope{parent: outputScope(cols, outer), items: sc.items}
	for _, e := range n.GroupBy {
		r.sources(e, grouping)
	}
	for _, e := range n.DistinctOn {
		r.sources(e, outputScope(cols, sc))
	}
	for _, w := range n.Window {
		r.window(w, sc)
	}
	return cols, sc
}

// selectTarget returns the columns one target of a select list produces: the
// columns a star expands to, or the one column of an expression.
func (r *spanResolver) selectTarget(target tree.SelectExpr, sc *spanScope) []spanColumn {
	switch e := target.Expr.(type) {
	case tree.UnqualifiedStar:
		return visibleColumns(sc.items)
	case *tree.UnresolvedName:
		if e.Star {
			if rel := r.qualifiedRelation(e, sc); rel != nil {
				return relationColumns(rel)
			}
			return nil
		}
	case *tree.TupleStar:
		if rel := r.wholeRow(e.Expr, sc); rel != nil {
			return relationColumns(rel)
		}
	}
	_, plain := target.Expr.(*tree.UnresolvedName)
	return []spanColumn{{name: renderName(target), sources: r.sources(target.Expr, sc), plain: plain}}
}

// wholeRow returns the relation a whole-row reference such as (c).* names,
// or nil.
func (r *spanResolver) wholeRow(e tree.Expr, sc *spanScope) *spanRelation {
	name, ok := tree.StripParens(e).(*tree.UnresolvedName)
	if !ok || name.Star || name.NumParts != 1 {
		return nil
	}
	for s := sc; s != nil; s = s.parent {
		if matches := columnsNamed(s.items, name.Parts[0]); len(matches) > 0 {
			return nil
		}
		if rel := relationNamed(s.items, name.Parts[0], "", ""); rel != nil {
			return rel
		}
	}
	return nil
}

// visibleColumns is what * reads from items.
func visibleColumns(items []*fromItem) []spanColumn {
	var out []spanColumn
	for _, item := range items {
		for _, c := range item.columns {
			if !c.hidden {
				out = append(out, spanColumn{name: c.name, sources: c.sources, plain: true})
			}
		}
		for _, rel := range item.opaque() {
			if len(rel.columns) == 0 {
				out = append(out, spanColumn{name: "*", sources: rel.opaqueColumn("").sources, asterisk: true})
			}
		}
	}
	return out
}

func relationColumns(rel *spanRelation) []spanColumn {
	return visibleColumns([]*fromItem{{relations: []*spanRelation{rel}, columns: rel.columns}})
}

// renderName is the name CockroachDB gives a result column. A name its rules
// cannot derive, such as that of a function called by OID, is ?column?, as
// CockroachDB names an expression it has no better name for; the name is not
// lineage.
func renderName(target tree.SelectExpr) string {
	resolveFunctionNames(target.Expr)
	name, err := tree.GetRenderColName(context.Background(), tree.EmptySearchPath, target, nil)
	if err != nil {
		return "?column?"
	}
	return name
}

// resolveFunctionNames resolves each function reference in e to its name, as
// CockroachDB's naming reads it, without a function registry.
func resolveFunctionNames(e tree.Expr) {
	// The visit function never fails, so neither does the visit.
	_, _ = tree.SimpleVisit(e, func(expr tree.Expr) (bool, tree.Expr, error) {
		if f, ok := expr.(*tree.FuncExpr); ok {
			if name, ok := f.Func.FunctionReference.(*tree.UnresolvedName); ok {
				f.Func.FunctionReference = &tree.ResolvedFunctionDefinition{Name: name.Parts[0]}
			}
		}
		return true, expr, nil
	})
}

// tableExpr resolves one FROM item. sc holds the items before it, which a
// LATERAL item and a function in FROM read.
func (r *spanResolver) tableExpr(te tree.TableExpr, sc *spanScope) *fromItem {
	switch n := te.(type) {
	case *tree.AliasedTableExpr:
		inner := sc
		if n.Lateral {
			inner = &spanScope{parent: sc}
		}
		item := r.tableExpr(n.Expr, inner)
		if n.As.Alias == "" && !n.Ordinality {
			return item
		}
		// An alias or WITH ORDINALITY makes the item one relation.
		rel := &spanRelation{name: string(n.As.Alias), columns: slices.Clone(item.columns)}
		if rel.name == "" && len(item.relations) == 1 {
			inner := item.relations[0]
			rel.name, rel.schema, rel.database = inner.name, inner.schema, inner.database
		}
		if fn, ok := n.Expr.(*tree.RowsFromExpr); ok && len(fn.Items) == 1 && len(item.columns) == 1 && n.As.Alias != "" && len(n.As.Cols) == 0 {
			// A function's one column takes the alias's name.
			rel.columns[0].name = rel.name
		}
		for _, inner := range item.opaque() {
			rel.opaque = true
			rel.opaqueReads = union(rel.opaqueReads, inner.opaqueReads)
			if inner.opaqueTable != nil {
				rel.opaqueTable = inner.opaqueTable
			}
		}
		// An alias beyond the known columns names a column of a relation
		// whose columns are unknown, which reads what the relation reads.
		var extra analysis.ColumnSet
		if rel.opaque {
			extra = rel.opaqueColumn("").sources
		}
		if n.Ordinality {
			// The row's number holds no column's value, but after a relation
			// of unknown columns its position may hold one of them.
			rel.columns = append(rel.columns, spanColumn{name: "ordinality", sources: union(extra)})
		}
		rel.columns = r.renameColumns(rel.columns, n.As.Cols, extra)
		return &fromItem{relations: []*spanRelation{rel}, columns: rel.columns}
	case *tree.TableName:
		rel := r.relation(n, sc)
		return &fromItem{relations: []*spanRelation{rel}, columns: rel.columns}
	case *tree.ParenTableExpr:
		return r.tableExpr(n.Expr, sc)
	case *tree.JoinTableExpr:
		return r.join(n, sc)
	case *tree.Subquery:
		cols, _ := r.selectStatement(n.Select, sc.outside())
		cols = stripColumns(cols)
		return &fromItem{relations: []*spanRelation{{columns: cols}}, columns: cols}
	case *tree.StatementSource:
		if sel, ok := n.Statement.(*tree.Select); ok {
			cols := stripColumns(r.selectQuery(sel, sc.outside()))
			return &fromItem{relations: []*spanRelation{{columns: cols}}, columns: cols}
		}
		if strings.HasPrefix(n.Statement.StatementTag(), "SHOW") && r.inView == 0 {
			r.readsSystem = true
		}
		// A data change in a statement source makes the statement DML, which
		// has no lineage; a SHOW reads system metadata, whose query has none;
		// an EXPLAIN's plan holds no column's value.
		r.statementTables(n.Statement, sc.outside())
		return &fromItem{relations: []*spanRelation{{opaque: true}}}
	case *tree.RowsFromExpr:
		rel := &spanRelation{}
		reads := analysis.ColumnSet{}
		for _, fn := range n.Items {
			cols, known := r.functionColumns(fn, sc)
			for _, c := range cols {
				reads = union(reads, c.sources)
			}
			rel.columns = append(rel.columns, cols...)
			rel.opaque = rel.opaque || !known
		}
		// A function in FROM returns the rows its arguments ask for.
		r.addPredicates(reads)
		if rel.opaque {
			// Past a function of unknown columns no column is where its
			// position says, so each reads every argument.
			rel.opaqueReads = reads
			for i := range rel.columns {
				rel.columns[i].sources = reads
			}
		}
		if len(n.Items) == 1 {
			rel.name = renderName(tree.SelectExpr{Expr: n.Items[0]})
		}
		return &fromItem{relations: []*spanRelation{rel}, columns: rel.columns}
	case *tree.TableRef:
		// A numeric reference names a relation by descriptor ID, which the
		// metadata cannot resolve, so it fails the statement.
		if r.metadataErr == nil {
			r.metadataErr = errors.New("a numeric table reference names no relation the metadata can resolve")
		}
		return &fromItem{relations: []*spanRelation{{opaque: true}}}
	}
	r.fail(fmt.Errorf("unsupported table expression %T", te))
	return &fromItem{relations: []*spanRelation{{opaque: true}}}
}

// functionColumns returns the columns a function in FROM returns, each with
// the arguments it reads, and whether they are known. A function of unknown
// columns, such as a user-defined one or json_to_record, is one column here.
func (r *spanResolver) functionColumns(e tree.Expr, sc *spanScope) ([]spanColumn, bool) {
	name := renderName(tree.SelectExpr{Expr: e})
	fn, ok := e.(*tree.FuncExpr)
	if !ok {
		// A special form such as COALESCE or a cast returns one value.
		return []spanColumn{{name: name, sources: r.sources(e, sc)}}, true
	}
	args := make([]analysis.ColumnSet, len(fn.Exprs))
	all := analysis.ColumnSet{}
	for i, arg := range fn.Exprs {
		args[i] = r.sources(arg, sc)
		all = union(all, args[i])
	}
	switch {
	case name == "unnest" && len(args) > 1:
		// unnest of several arrays returns a column for each.
		cols := make([]spanColumn, len(args))
		for i := range args {
			cols[i] = spanColumn{name: name, sources: args[i]}
		}
		return cols, true
	case keyValueFunctions[name]:
		return []spanColumn{{name: "key", sources: all}, {name: "value", sources: union(all)}}, true
	case oneColumnFunctions[name]:
		return []spanColumn{{name: name, sources: all}}, true
	}
	return []spanColumn{{name: name, sources: all}}, false
}

// keyValueFunctions return a key and a value column; oneColumnFunctions
// return one column, named after the function.
var (
	keyValueFunctions  = setOf("json_each", "jsonb_each", "json_each_text", "jsonb_each_text")
	oneColumnFunctions = setOf(
		"generate_series", "generate_subscripts", "unnest", "json_array_elements", "jsonb_array_elements",
		"json_array_elements_text", "jsonb_array_elements_text", "json_object_keys", "jsonb_object_keys",
		"regexp_split_to_table", "regexp_matches",
	)
)

// renameColumns applies a column alias list to the visible columns in order.
// An alias beyond them names a column that reads extra.
func (r *spanResolver) renameColumns(cols []spanColumn, aliases tree.ColumnDefList, extra analysis.ColumnSet) []spanColumn {
	out := make([]spanColumn, 0, len(cols))
	i := 0
	for _, c := range cols {
		if !c.hidden && i < len(aliases) {
			c.name = string(aliases[i].Name)
			i++
		}
		out = append(out, c)
	}
	if i < len(aliases) && extra == nil {
		r.fail(pgerror.Newf(pgcode.InvalidColumnReference, "%d columns available but %d columns specified", i, len(aliases)))
	}
	for ; i < len(aliases); i++ {
		out = append(out, spanColumn{name: string(aliases[i].Name), sources: union(extra)})
	}
	return out
}

func (r *spanResolver) join(n *tree.JoinTableExpr, sc *spanScope) *fromItem {
	left := r.tableExpr(n.Left, sc)
	// A LATERAL right side reads the left.
	right := r.tableExpr(n.Right, &spanScope{parent: sc, items: []*fromItem{left}, from: true})
	item := &fromItem{relations: append(slices.Clone(left.relations), right.relations...)}
	var using []string
	switch cond := n.Cond.(type) {
	case *tree.OnJoinCond:
		// ON reads the two sides and the query around the FROM clause.
		r.addPredicates(r.sources(cond.Expr, &spanScope{parent: sc.outside(), items: []*fromItem{left, right}}))
	case *tree.UsingJoinCond:
		for _, name := range cond.Cols {
			using = append(using, string(name))
		}
	case tree.NaturalJoinCond:
		for _, c := range left.columns {
			if c.hidden {
				continue
			}
			if _, ok := visibleColumn(right.columns, c.name); ok {
				using = append(using, c.name)
			}
		}
	}
	merged := map[string]bool{}
	for _, name := range using {
		l, lok := visibleColumn(left.columns, name)
		rc, rok := visibleColumn(right.columns, name)
		if !lok || !rok {
			side := "left"
			if lok {
				side = "right"
			}
			r.fail(pgerror.Newf(pgcode.UndefinedColumn, "column %q specified in USING clause does not exist in %s table", name, side))
			continue
		}
		// The merged column is the left side's value, the right side's in a
		// RIGHT JOIN, and either in a FULL JOIN.
		sources := l.sources
		switch n.JoinType {
		case tree.AstRight:
			sources = rc.sources
		case tree.AstFull:
			sources = union(l.sources, rc.sources)
		}
		r.addPredicates(union(l.sources, rc.sources))
		item.columns = append(item.columns, spanColumn{name: name, sources: sources})
		merged[name] = true
	}
	for _, side := range [][]spanColumn{left.columns, right.columns} {
		skipped := map[string]bool{}
		for _, c := range side {
			if merged[c.name] && !c.hidden && !skipped[c.name] {
				skipped[c.name] = true
				continue
			}
			item.columns = append(item.columns, c)
		}
	}
	return item
}

func visibleColumn(cols []spanColumn, name string) (spanColumn, bool) {
	for _, c := range cols {
		if !c.hidden && c.name == name {
			return c, true
		}
	}
	return spanColumn{}, false
}

// relation resolves a table name in FROM: a CTE, a table, a view or a system
// relation.
func (r *spanResolver) relation(tn *tree.TableName, sc *spanScope) *spanRelation {
	name := tn.Table()
	if !tn.ExplicitSchema && !tn.ExplicitCatalog {
		if cols, ok := sc.cte(name); ok {
			if cols == nil {
				return &spanRelation{name: name, opaque: true}
			}
			return &spanRelation{name: name, columns: stripColumns(cols)}
		}
	}
	return r.namedRelation(tn)
}

// namedRelation resolves a table name to a table, a view or a system
// relation, whatever CTE has its name.
func (r *spanResolver) namedRelation(tn *tree.TableName) *spanRelation {
	name := tn.Table()
	resource, db := r.locate(tn)
	rel := &spanRelation{name: name, schema: resource.Schema, database: resource.Database}
	if isSystemResource(resource) {
		r.access(resource)
		rel.opaque, rel.opaqueTable = true, &resource
		return rel
	}
	if db == nil || resource.Schema == "" {
		return r.missing(rel, resource)
	}
	r.access(resource)
	if table := db.GetTable(resource.Schema, name); table != nil {
		for _, c := range table.GetColumns() {
			column := resource
			column.Column = c.GetName()
			rel.columns = append(rel.columns, spanColumn{name: c.GetName(), sources: analysis.ColumnSet{column: true}, hidden: c.GetIsInvisible()})
		}
		if len(rel.columns) == 0 {
			rel.opaque, rel.opaqueTable = true, &resource
		}
		return rel
	}
	for _, s := range db.Metadata().GetSchemas() {
		if s.GetName() != resource.Schema {
			continue
		}
		for _, v := range s.GetViews() {
			if v.GetName() == name {
				return r.view(rel, resource, v.GetDefinition(), v.GetColumns(), true)
			}
		}
		for _, v := range s.GetMaterializedViews() {
			if v.GetName() == name {
				// A materialized view holds its rows, so its columns are
				// its own.
				return r.view(rel, resource, v.GetDefinition(), v.GetColumns(), false)
			}
		}
	}
	return r.missing(rel, resource)
}

// missing records a relation the metadata lacks: an access of the name as
// resolved, and the span's NotFoundError. Its columns are unknown.
func (r *spanResolver) missing(rel *spanRelation, resource analysis.ColumnResource) *spanRelation {
	if resource.Schema == "" {
		resource.Schema = r.firstSchema()
		rel.schema = resource.Schema
	}
	r.access(resource)
	if r.notFound == nil {
		database, schema, table := resource.Database, resource.Schema, resource.Table
		r.notFound = &analysis.ResourceNotFoundError{Database: &database, Schema: &schema, Table: &table}
	}
	rel.opaque, rel.opaqueTable = true, &resource
	return rel
}

func (r *spanResolver) firstSchema() string {
	for _, s := range r.searchPath {
		if !strings.HasPrefix(s, "$") {
			return s
		}
	}
	return "public"
}

// locate resolves where tn lives as CockroachDB does: an unqualified name in
// pg_catalog, then the search path, then pg_extension; a two-part name in a
// schema of the current database, or else in the public schema of a database
// of that name. The schema is empty when no synced relation has the name.
func (r *spanResolver) locate(tn *tree.TableName) (analysis.ColumnResource, Database) {
	name := tn.Table()
	resource := analysis.ColumnResource{Database: r.defaultDatabase, Table: name}
	switch {
	case tn.ExplicitCatalog:
		resource.Database, resource.Schema = tn.Catalog(), tn.Schema()
	case tn.ExplicitSchema:
		resource.Schema = tn.Schema()
		if !IsSystemSchema(resource.Schema) && !r.hasSchema(r.defaultDatabase, resource.Schema) {
			resource.Database, resource.Schema = resource.Schema, "public"
		}
	default:
		if pgCatalogRelations[name] {
			resource.Schema = "pg_catalog"
			return resource, nil
		}
		db := r.database(resource.Database)
		if db != nil {
			resource.Schema, _ = db.SearchRelation(r.searchPath, name)
		}
		if resource.Schema == "" && pgExtensionRelations[name] {
			resource.Schema = "pg_extension"
		}
		return resource, db
	}
	if isSystemResource(resource) {
		return resource, nil
	}
	db := r.database(resource.Database)
	if db == nil {
		return resource, nil
	}
	if schema, _ := db.SearchRelation([]string{resource.Schema}, name); schema == "" {
		return resource, nil
	}
	return resource, db
}

func (r *spanResolver) hasSchema(database, schema string) bool {
	db := r.database(database)
	if db == nil {
		return false
	}
	for _, s := range db.Metadata().GetSchemas() {
		if s.GetName() == schema {
			return true
		}
	}
	return false
}

// view gives rel a view's columns from its definition: named by the synced
// columns where they match, each reading the view's own column and, for a
// view CockroachDB follows, the columns the definition reads for it. A
// definition that cannot be resolved fails the statement, and its columns
// are unknown.
func (r *spanResolver) view(rel *spanRelation, resource analysis.ColumnResource, definition string, synced []*metadata.ColumnMetadata, follow bool) *spanRelation {
	unknown := func(err error) *spanRelation {
		r.fail(fmt.Errorf("the definition of view %s: %w", resource.String(), err))
		rel.opaque, rel.opaqueTable = true, &resource
		return rel
	}
	if r.expanding[resource] {
		return unknown(errors.New("the view refers to itself"))
	}
	stmts, err := Parse(definition)
	if err != nil {
		return unknown(err)
	}
	var sel *tree.Select
	if len(stmts) == 1 {
		sel, _ = stmts[0].AST.(*tree.Select)
	}
	if sel == nil {
		return unknown(errors.New("the definition is not a query"))
	}
	r.expanding[resource] = true
	r.inView++
	before := r.resolveErr
	defCols := r.selectQuery(sel, nil)
	r.inView--
	delete(r.expanding, resource)
	if before == nil && r.resolveErr != nil {
		r.resolveErr = fmt.Errorf("the definition of view %s: %w", resource.String(), r.resolveErr)
	}
	for i, c := range defCols {
		name := c.name
		if len(synced) == len(defCols) {
			name = synced[i].GetName()
		}
		own := resource
		own.Column = name
		sources := analysis.ColumnSet{own: true}
		if follow {
			sources = union(sources, c.sources)
		}
		rel.columns = append(rel.columns, spanColumn{name: name, sources: sources})
	}
	return rel
}

// sources returns the columns e reads. A name it cannot resolve is recorded
// as the resolver's error, and the walk goes on.
func (r *spanResolver) sources(e tree.Expr, sc *spanScope) analysis.ColumnSet {
	v := &exprVisitor{r: r, sc: sc, sources: analysis.ColumnSet{}}
	if e != nil {
		tree.WalkExpr(v, e)
	}
	return v.sources
}

// window resolves a window definition's expressions, which add no lineage
// and no condition, for the relations their subqueries read.
func (r *spanResolver) window(w *tree.WindowDef, sc *spanScope) {
	for _, e := range w.Partitions {
		r.sources(e, sc)
	}
	r.orderBy(w.OrderBy, sc)
	if w.Frame != nil {
		for _, bound := range []*tree.WindowFrameBound{w.Frame.Bounds.StartBound, w.Frame.Bounds.EndBound} {
			if bound != nil {
				r.sources(bound.OffsetExpr, sc)
			}
		}
	}
}

// exprVisitor collects the columns an expression reads. A subquery reads the
// columns of its results; a window function's OVER clause reads none.
type exprVisitor struct {
	r       *spanResolver
	sc      *spanScope
	sources analysis.ColumnSet
}

func (v *exprVisitor) VisitPre(expr tree.Expr) (bool, tree.Expr) {
	switch e := expr.(type) {
	case *tree.UnresolvedName:
		v.add(v.r.nameSources(e, v.sc))
		return false, expr
	case *tree.ColumnAccessExpr:
		// (c).name reads one column of a whole-row reference.
		if rel := v.r.wholeRow(e.Expr, v.sc); rel != nil && !e.ByIndex {
			if c, ok := findColumn(rel.columns, string(e.ColName)); ok {
				v.add(c.sources)
				return false, expr
			}
		}
	case *tree.Subquery:
		cols, _ := v.r.selectStatement(e.Select, v.sc)
		for _, c := range cols {
			v.add(c.sources)
		}
		return false, expr
	case *tree.FuncExpr:
		if e.Filter != nil {
			// An aggregate's FILTER decides which rows it reads.
			v.r.addPredicates(v.r.sources(e.Filter, v.sc))
		}
		if e.WindowDef != nil {
			v.r.window(e.WindowDef, v.sc)
		}
	}
	return true, expr
}

func (v *exprVisitor) VisitPost(expr tree.Expr) tree.Expr { return expr }

func (v *exprVisitor) add(sources analysis.ColumnSet) {
	for c := range sources {
		v.sources[c] = true
	}
}

func findColumn(cols []spanColumn, name string) (spanColumn, bool) {
	for _, c := range cols {
		if c.name == name {
			return c, true
		}
	}
	return spanColumn{}, false
}

// nameSources resolves a name in an expression: a column, a relation's star,
// or a whole-row reference to a relation.
func (r *spanResolver) nameSources(name *tree.UnresolvedName, sc *spanScope) analysis.ColumnSet {
	if name.Star {
		if rel := r.qualifiedRelation(name, sc); rel != nil {
			return columnSources(relationColumns(rel))
		}
		return nil
	}
	column := name.Parts[0]
	if name.NumParts > 1 {
		rel := r.qualifiedRelation(name, sc)
		if rel == nil {
			return nil
		}
		if c, ok := findColumn(rel.columns, column); ok {
			return c.sources
		}
		if rel.opaque {
			return rel.opaqueColumn(column).sources
		}
		r.fail(pgerror.Newf(pgcode.UndefinedColumn, "column %q does not exist", name.String()))
		return nil
	}
	// A relation of unknown columns nearer than the match may be what the
	// name reads, so its sources count too.
	var unknown analysis.ColumnSet
	for s := sc; s != nil; s = s.parent {
		matches := columnsNamed(s.items, column)
		if len(matches) > 1 {
			r.fail(pgerror.Newf(pgcode.AmbiguousColumn, "column reference %q is ambiguous", column))
			return union(columnSources(matches), unknown)
		}
		if len(matches) == 1 {
			return union(matches[0].sources, unknown)
		}
		if rel := relationNamed(s.items, column, "", ""); rel != nil {
			return union(columnSources(relationColumns(rel)), unknown)
		}
		for _, item := range s.items {
			for _, rel := range item.opaque() {
				unknown = union(unknown, rel.opaqueColumn(column).sources)
			}
		}
	}
	if unknown == nil {
		r.fail(pgerror.Newf(pgcode.UndefinedColumn, "column %q does not exist", column))
	}
	return unknown
}

func columnSources(cols []spanColumn) analysis.ColumnSet {
	out := analysis.ColumnSet{}
	for _, c := range cols {
		for s := range c.sources {
			out[s] = true
		}
	}
	return out
}

func columnsNamed(items []*fromItem, name string) []spanColumn {
	var out []spanColumn
	for _, item := range items {
		for _, c := range item.columns {
			if c.name == name {
				out = append(out, c)
			}
		}
	}
	return out
}

// qualifiedRelation resolves the relation that qualifies a column or a star:
// table, schema.table or database.schema.table, or records that none does.
func (r *spanResolver) qualifiedRelation(name *tree.UnresolvedName, sc *spanScope) *spanRelation {
	table, schema, database := name.Parts[1], name.Parts[2], name.Parts[3]
	for s := sc; s != nil; s = s.parent {
		if rel := relationNamed(s.items, table, schema, database); rel != nil {
			return rel
		}
	}
	var prefix []string
	for i := name.NumParts - 1; i >= 1; i-- {
		prefix = append(prefix, name.Parts[i])
	}
	r.fail(pgerror.Newf(pgcode.UndefinedTable, "no data source matches prefix: %s", strings.Join(prefix, ".")))
	return nil
}

func relationNamed(items []*fromItem, table, schema, database string) *spanRelation {
	for _, item := range items {
		for _, rel := range item.relations {
			if rel.name != table || (schema != "" && rel.schema != schema) || (database != "" && rel.database != database) {
				continue
			}
			return rel
		}
	}
	return nil
}

// statementTables records the relations a statement other than a query reads,
// for its access check, walking every clause that holds an expression or a
// relation, and the columns that decide which rows it changes or returns.
// outer is the scope around a statement nested in another, whose CTEs it
// sees. Other DDL names relations but reads no rows of them.
func (r *spanResolver) statementTables(stmt tree.Statement, outer *spanScope) {
	sc := &spanScope{parent: outer}
	// A statement's target is a table whatever CTE shares its name, and it
	// opens the scope of the statement's other clauses.
	target := func(te tree.TableExpr) {
		item := r.tableExpr(te, &spanScope{})
		sc.items = append(sc.items, item)
	}
	from := func(tables tree.TableExprs) {
		sc.from = true
		for _, te := range tables {
			sc.items = append(sc.items, r.tableExpr(te, sc))
		}
		sc.from = false
	}
	walk := func(e tree.Expr) { r.sources(e, sc) }
	with := func(w *tree.With) {
		if w != nil {
			sc.parent = r.with(w, outer)
		}
	}
	where := func(w *tree.Where) {
		if w != nil {
			r.addPredicates(r.sources(w.Expr, sc))
		}
	}
	updates := func(exprs tree.UpdateExprs) {
		for _, e := range exprs {
			walk(e.Expr)
		}
	}
	returning := func(clause tree.ReturningClause) {
		if exprs, ok := clause.(*tree.ReturningExprs); ok {
			for _, e := range *exprs {
				walk(e.Expr)
			}
		}
	}
	switch n := stmt.(type) {
	case *tree.Select:
		r.selectQuery(n, outer)
	case *tree.ParenSelect:
		r.selectQuery(n.Select, outer)
	case *tree.Explain:
		r.statementTables(n.Statement, outer)
	case *tree.ExplainAnalyze:
		r.statementTables(n.Statement, outer)
	case *tree.Insert:
		// INSERT and UPSERT.
		with(n.With)
		if n.Rows != nil {
			r.selectQuery(n.Rows, sc.parent)
		}
		target(n.Table)
		if n.OnConflict != nil {
			// ON CONFLICT reads the target and the row proposed for it.
			if item := sc.items[0]; len(item.relations) == 1 {
				excluded := &spanRelation{name: "excluded", columns: item.columns}
				sc.items = append(sc.items, &fromItem{relations: []*spanRelation{excluded}})
			}
			// The arbiter predicate and DO UPDATE's WHERE decide which rows
			// conflict and which of them change.
			r.addPredicates(r.sources(n.OnConflict.ArbiterPredicate, sc))
			updates(n.OnConflict.Exprs)
			where(n.OnConflict.Where)
		}
		returning(n.Returning)
	case *tree.Update:
		with(n.With)
		target(n.Table)
		from(n.From)
		updates(n.Exprs)
		where(n.Where)
		r.orderBy(n.OrderBy, sc)
		r.limit(n.Limit, sc)
		returning(n.Returning)
	case *tree.Delete:
		with(n.With)
		target(n.Table)
		from(n.Using)
		if n.Batch != nil {
			for _, param := range n.Batch.Params {
				if size, ok := param.(*tree.SizeBatchParam); ok {
					walk(size.Size)
				}
			}
		}
		where(n.Where)
		r.orderBy(n.OrderBy, sc)
		r.limit(n.Limit, sc)
		returning(n.Returning)
	case *tree.CopyFrom:
		target(&n.Table)
	case *tree.CopyTo:
		if n.Statement != nil {
			r.statementTables(n.Statement, outer)
		} else {
			target(&n.Table)
		}
	case *tree.Import:
		if n.Table != nil {
			target(n.Table)
		}
	case *tree.Call:
		if n.Proc != nil {
			walk(n.Proc)
		}
	case *tree.CreateTable:
		if n.AsSource != nil {
			r.selectQuery(n.AsSource, outer)
		}
	case *tree.CreateView:
		if n.AsSource != nil {
			r.selectQuery(n.AsSource, outer)
		}
	}
}

// pgCatalogRelations are the relations of CockroachDB v25.2's pg_catalog,
// which an unqualified name reaches before the search path.
var pgCatalogRelations = setOf(
	"pg_aggregate", "pg_am", "pg_amop", "pg_amproc", "pg_attrdef", "pg_attribute", "pg_auth_members", "pg_authid",
	"pg_available_extension_versions", "pg_available_extensions", "pg_cast", "pg_class", "pg_collation", "pg_config",
	"pg_constraint", "pg_conversion", "pg_cursors", "pg_database", "pg_db_role_setting", "pg_default_acl", "pg_depend",
	"pg_description", "pg_enum", "pg_event_trigger", "pg_extension", "pg_file_settings", "pg_foreign_data_wrapper",
	"pg_foreign_server", "pg_foreign_table", "pg_group", "pg_hba_file_rules", "pg_index", "pg_indexes", "pg_inherits",
	"pg_init_privs", "pg_language", "pg_largeobject", "pg_largeobject_metadata", "pg_locks", "pg_matviews",
	"pg_namespace", "pg_opclass", "pg_operator", "pg_opfamily", "pg_partitioned_table", "pg_policies", "pg_policy",
	"pg_prepared_statements", "pg_prepared_xacts", "pg_proc", "pg_publication", "pg_publication_rel",
	"pg_publication_tables", "pg_range", "pg_replication_origin", "pg_replication_origin_status",
	"pg_replication_slots", "pg_rewrite", "pg_roles", "pg_rules", "pg_seclabel", "pg_seclabels", "pg_sequence",
	"pg_sequences", "pg_settings", "pg_shadow", "pg_shdepend", "pg_shdescription", "pg_shmem_allocations",
	"pg_shseclabel", "pg_stat_activity", "pg_stat_all_indexes", "pg_stat_all_tables", "pg_stat_archiver",
	"pg_stat_bgwriter", "pg_stat_database", "pg_stat_database_conflicts", "pg_stat_gssapi", "pg_stat_progress_analyze",
	"pg_stat_progress_basebackup", "pg_stat_progress_cluster", "pg_stat_progress_create_index",
	"pg_stat_progress_vacuum", "pg_stat_replication", "pg_stat_slru", "pg_stat_ssl", "pg_stat_subscription",
	"pg_stat_sys_indexes", "pg_stat_sys_tables", "pg_stat_user_functions", "pg_stat_user_indexes",
	"pg_stat_user_tables", "pg_stat_wal_receiver", "pg_stat_xact_all_tables", "pg_stat_xact_sys_tables",
	"pg_stat_xact_user_functions", "pg_stat_xact_user_tables", "pg_statio_all_indexes", "pg_statio_all_sequences",
	"pg_statio_all_tables", "pg_statio_sys_indexes", "pg_statio_sys_sequences", "pg_statio_sys_tables",
	"pg_statio_user_indexes", "pg_statio_user_sequences", "pg_statio_user_tables", "pg_statistic", "pg_statistic_ext",
	"pg_statistic_ext_data", "pg_stats", "pg_stats_ext", "pg_subscription", "pg_subscription_rel", "pg_tables",
	"pg_tablespace", "pg_timezone_abbrevs", "pg_timezone_names", "pg_transform", "pg_trigger", "pg_ts_config",
	"pg_ts_config_map", "pg_ts_dict", "pg_ts_parser", "pg_ts_template", "pg_type", "pg_user", "pg_user_mapping",
	"pg_user_mappings", "pg_views",
)

// pgExtensionRelations are the relations of CockroachDB's pg_extension, which
// an unqualified name reaches after the search path.
var pgExtensionRelations = setOf("geography_columns", "geometry_columns", "spatial_ref_sys")

func setOf(names ...string) map[string]bool {
	out := make(map[string]bool, len(names))
	for _, n := range names {
		out[n] = true
	}
	return out
}
