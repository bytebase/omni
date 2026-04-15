// Package analysis extracts structural information from parsed Snowflake queries.
//
// The primary entry point is Extract (or ExtractSQL), which produces a QuerySpan
// describing what a SELECT/WITH/set-operation query reads and produces.
package analysis

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
	"github.com/bytebase/omni/snowflake/parser"
)

// ---------------------------------------------------------------------------
// Public types
// ---------------------------------------------------------------------------

// QuerySpan summarises what a query reads and produces.
type QuerySpan struct {
	Results []*ResultColumn // one entry per SELECT-list column
	Sources []*SourceColumn // every source column read by the query (flat union)
}

// ResultColumn represents one column in the query's result set.
type ResultColumn struct {
	Name      string          // output alias or derived name; "*" for unresolved star
	Sources   []*SourceColumn // source columns this result depends on (may be empty for constants)
	IsDerived bool            // true if computed from 2+ sources or involves a function/expression
}

// SourceColumn identifies a column read from a base table (or CTE).
type SourceColumn struct {
	Database string // may be empty
	Schema   string // may be empty
	Table    string // table name or CTE alias
	Column   string // column name; "*" for unresolved star
}

// ---------------------------------------------------------------------------
// Public API
// ---------------------------------------------------------------------------

// Extract runs the span extractor on a parsed query expression.
// The input must be a *ast.SelectStmt, *ast.SetOperationStmt, or *ast.File
// whose first stmt is one of those.
// Returns nil, error if the input isn't a query expression.
func Extract(node ast.Node) (*QuerySpan, error) {
	switch n := node.(type) {
	case *ast.File:
		if len(n.Stmts) == 0 {
			return nil, fmt.Errorf("analysis.Extract: empty file")
		}
		return Extract(n.Stmts[0])
	case *ast.SelectStmt:
		scope := newScope()
		return extractSelectStmt(n, scope), nil
	case *ast.SetOperationStmt:
		scope := newScope()
		return extractSetOperationStmt(n, scope), nil
	default:
		return nil, fmt.Errorf("analysis.Extract: unsupported node type %T", node)
	}
}

// ExtractSQL parses and extracts in one call.
func ExtractSQL(sql string) (*QuerySpan, error) {
	file, err := parser.Parse(sql)
	if err != nil {
		return nil, fmt.Errorf("analysis.ExtractSQL: parse error: %w", err)
	}
	return Extract(file)
}

// ---------------------------------------------------------------------------
// Internal scope
// ---------------------------------------------------------------------------

// queryScope carries CTE definitions visible at the current scope level.
type queryScope struct {
	ctes map[string]*QuerySpan // normalized CTE name → its span
}

func newScope() *queryScope {
	return &queryScope{ctes: make(map[string]*QuerySpan)}
}

// childScope returns a new scope that inherits all CTEs from the parent.
func childScope(parent *queryScope) *queryScope {
	child := newScope()
	for k, v := range parent.ctes {
		child.ctes[k] = v
	}
	return child
}

// ---------------------------------------------------------------------------
// tableEntry — a resolved FROM source with its known output columns
// ---------------------------------------------------------------------------

// tableEntry describes a single resolved FROM source (base table, subquery,
// or CTE) with the columns it exposes.
type tableEntry struct {
	alias   string          // the alias (or table name) used to qualify columns
	columns []*SourceColumn // nil means "unknown schema" (treat as *-source)
	// For a base table with unknown schema, we synthesize a * source column
	// rather than enumerating specific columns.
	isUnknown bool // true when we have no schema information
}

// ---------------------------------------------------------------------------
// Main extraction functions
// ---------------------------------------------------------------------------

// extractSpan dispatches to the correct extractor based on node type.
func extractSpan(node ast.Node, scope *queryScope) *QuerySpan {
	switch n := node.(type) {
	case *ast.SelectStmt:
		return extractSelectStmt(n, scope)
	case *ast.SetOperationStmt:
		return extractSetOperationStmt(n, scope)
	default:
		return &QuerySpan{}
	}
}

// extractSelectStmt extracts a span from a SELECT statement.
func extractSelectStmt(s *ast.SelectStmt, scope *queryScope) *QuerySpan {
	// If this SELECT has a WITH clause, build a child scope with the CTEs.
	if len(s.With) > 0 {
		child := childScope(scope)
		for _, cte := range s.With {
			cteSpan := extractSpan(cte.Query, child)
			// If the CTE has column aliases, rename the result columns.
			if len(cte.Columns) > 0 {
				for i, alias := range cte.Columns {
					if i < len(cteSpan.Results) {
						cteSpan.Results[i].Name = alias.Normalize()
					}
				}
			}
			child.ctes[cte.Name.Normalize()] = cteSpan
		}
		scope = child
	}

	// Build FROM scope entries.
	fromEntries := buildFromScope(s.From, scope)

	// Collect result columns.
	var results []*ResultColumn
	for _, target := range s.Targets {
		if target.Star {
			// SELECT * or SELECT t.* [EXCLUDE (cols)]
			results = append(results, resolveStarTarget(target, fromEntries))
		} else {
			rc := resolveExprTarget(target, fromEntries)
			results = append(results, rc)
		}
	}

	// Build the flat Sources list from all result column sources + any
	// additional table-level sources from FROM entries.
	sourcesMap := make(map[sourceKey]*SourceColumn)

	// Collect from result columns.
	for _, rc := range results {
		for _, sc := range rc.Sources {
			k := toSourceKey(sc)
			if _, exists := sourcesMap[k]; !exists {
				sourcesMap[k] = sc
			}
		}
	}

	// Also include "table-level" sources for every FROM entry even if no
	// columns are explicitly referenced (e.g. SELECT 1 FROM t still "touches" t).
	for _, entry := range fromEntries {
		sc := &SourceColumn{Table: entry.alias, Column: "*"}
		k := toSourceKey(sc)
		if _, exists := sourcesMap[k]; !exists {
			sourcesMap[k] = sc
		}
	}

	sources := flattenSources(sourcesMap)

	return &QuerySpan{
		Results: results,
		Sources: sources,
	}
}

// extractSetOperationStmt extracts a span from a UNION/INTERSECT/EXCEPT.
func extractSetOperationStmt(s *ast.SetOperationStmt, scope *queryScope) *QuerySpan {
	leftSpan := extractSpan(s.Left, scope)
	rightSpan := extractSpan(s.Right, scope)

	var mergedResults []*ResultColumn

	if s.ByName {
		// UNION ALL BY NAME: merge by column name.
		rightByName := make(map[string]*ResultColumn, len(rightSpan.Results))
		for _, rc := range rightSpan.Results {
			rightByName[strings.ToUpper(rc.Name)] = rc
		}
		seen := make(map[string]bool)
		for _, lrc := range leftSpan.Results {
			key := strings.ToUpper(lrc.Name)
			seen[key] = true
			merged := mergeResultColumns(lrc, rightByName[key])
			mergedResults = append(mergedResults, merged)
		}
		// Append right-only columns.
		for _, rrc := range rightSpan.Results {
			if !seen[strings.ToUpper(rrc.Name)] {
				mergedResults = append(mergedResults, rrc)
			}
		}
	} else {
		// Positional merge.
		maxLen := len(leftSpan.Results)
		if len(rightSpan.Results) > maxLen {
			maxLen = len(rightSpan.Results)
		}
		for i := 0; i < maxLen; i++ {
			var lrc, rrc *ResultColumn
			if i < len(leftSpan.Results) {
				lrc = leftSpan.Results[i]
			}
			if i < len(rightSpan.Results) {
				rrc = rightSpan.Results[i]
			}
			mergedResults = append(mergedResults, mergeResultColumns(lrc, rrc))
		}
	}

	// Sources = union of both branches.
	sourcesMap := make(map[sourceKey]*SourceColumn)
	for _, sc := range leftSpan.Sources {
		sourcesMap[toSourceKey(sc)] = sc
	}
	for _, sc := range rightSpan.Sources {
		k := toSourceKey(sc)
		if _, exists := sourcesMap[k]; !exists {
			sourcesMap[k] = sc
		}
	}

	return &QuerySpan{
		Results: mergedResults,
		Sources: flattenSources(sourcesMap),
	}
}

// ---------------------------------------------------------------------------
// FROM scope building
// ---------------------------------------------------------------------------

// buildFromScope converts a list of FROM items (TableRef / JoinExpr) into
// a flat list of tableEntry values.
func buildFromScope(froms []ast.Node, scope *queryScope) []tableEntry {
	var entries []tableEntry
	for _, item := range froms {
		entries = append(entries, resolveFromItem(item, scope)...)
	}
	return entries
}

// resolveFromItem resolves a single FROM item (TableRef or JoinExpr).
func resolveFromItem(item ast.Node, scope *queryScope) []tableEntry {
	switch n := item.(type) {
	case *ast.TableRef:
		return []tableEntry{resolveTableRef(n, scope)}
	case *ast.JoinExpr:
		return resolveJoin(n, scope)
	default:
		return nil
	}
}

// resolveTableRef resolves a single TableRef to a tableEntry.
func resolveTableRef(ref *ast.TableRef, scope *queryScope) tableEntry {
	// Subquery: (SELECT ...) AS alias
	if ref.Subquery != nil {
		subSpan := extractSpan(ref.Subquery, scope)
		alias := ref.Alias.Normalize()
		if alias == "" {
			alias = "_subquery"
		}
		// Convert subquery result columns to source columns for this virtual table.
		var cols []*SourceColumn
		for _, rc := range subSpan.Results {
			cols = append(cols, &SourceColumn{
				Table:  alias,
				Column: rc.Name,
			})
		}
		return tableEntry{alias: alias, columns: cols, isUnknown: false}
	}

	// Function call in FROM (TABLE(func(...)) etc.) — treat as unknown source.
	if ref.FuncCall != nil {
		alias := ref.Alias.Normalize()
		if alias == "" {
			alias = ref.FuncCall.Name.Name.Normalize()
		}
		return tableEntry{alias: alias, isUnknown: true}
	}

	// Base table or CTE reference.
	if ref.Name != nil {
		tableName := ref.Name.Name.Normalize()
		alias := ref.Alias.Normalize()
		if alias == "" {
			alias = tableName
		}

		// Check if this is a CTE reference.
		if cteSpan, ok := scope.ctes[tableName]; ok {
			var cols []*SourceColumn
			for _, rc := range cteSpan.Results {
				cols = append(cols, &SourceColumn{
					Table:  alias,
					Column: rc.Name,
				})
			}
			return tableEntry{alias: alias, columns: cols, isUnknown: false}
		}

		// Base table: build a tableEntry with the qualified name and mark as unknown.
		// We store database/schema in the entry so SourceColumns can be qualified.
		db := ref.Name.Database.Normalize()
		schema := ref.Name.Schema.Normalize()
		return tableEntry{
			alias:     alias,
			isUnknown: true,
			// Store db/schema info for use when building SourceColumns.
			// We'll use a special column to carry this.
			columns: []*SourceColumn{{Database: db, Schema: schema, Table: alias, Column: "*"}},
		}
	}

	return tableEntry{alias: "_unknown", isUnknown: true}
}

// resolveJoin flattens a JoinExpr into tableEntry slices.
func resolveJoin(join *ast.JoinExpr, scope *queryScope) []tableEntry {
	left := resolveFromItem(join.Left, scope)
	right := resolveFromItem(join.Right, scope)
	return append(left, right...)
}

// ---------------------------------------------------------------------------
// SelectTarget resolution
// ---------------------------------------------------------------------------

// resolveStarTarget handles a Star SelectTarget (SELECT * or SELECT t.*).
func resolveStarTarget(target *ast.SelectTarget, fromEntries []tableEntry) *ResultColumn {
	// Build a set of excluded column names (normalized).
	excluded := make(map[string]bool, len(target.Exclude))
	for _, ex := range target.Exclude {
		excluded[strings.ToUpper(ex.Normalize())] = true
	}

	var sources []*SourceColumn

	// Determine which tables to pull * from.
	// target.Expr holds a *ast.StarExpr when Star==true:
	//   - bare *:     StarExpr{Qualifier: nil}
	//   - qualified:  StarExpr{Qualifier: &ObjectName{...}}
	if star, ok := target.Expr.(*ast.StarExpr); ok && star != nil && star.Qualifier != nil {
		// SELECT t.* — qualified star; expand only for that table.
		qualifier := star.Qualifier.Name.Normalize()
		for _, entry := range fromEntries {
			if strings.EqualFold(entry.alias, qualifier) {
				sources = expandEntryStar(entry, excluded)
				break
			}
		}
	} else {
		// Bare SELECT * (or target.Expr is nil): pull from all FROM entries.
		for _, entry := range fromEntries {
			sources = append(sources, expandEntryStar(entry, excluded)...)
		}
	}

	return &ResultColumn{
		Name:      "*",
		Sources:   sources,
		IsDerived: false,
	}
}

// expandEntryStar expands a tableEntry for * selection, filtering excluded columns.
func expandEntryStar(entry tableEntry, excluded map[string]bool) []*SourceColumn {
	if entry.isUnknown || len(entry.columns) == 0 {
		// Unknown table: emit a * pseudo-source.
		return []*SourceColumn{{Table: entry.alias, Column: "*"}}
	}
	var cols []*SourceColumn
	for _, col := range entry.columns {
		if col.Column == "*" {
			// Unknown schema: just pass it through.
			cols = append(cols, col)
			continue
		}
		if excluded[strings.ToUpper(col.Column)] {
			continue
		}
		cols = append(cols, col)
	}
	return cols
}

// resolveExprTarget handles a non-star SelectTarget.
func resolveExprTarget(target *ast.SelectTarget, fromEntries []tableEntry) *ResultColumn {
	name := deriveResultName(target)
	refs := collectColumnRefs(target.Expr, fromEntries)
	derived := isDerivedExpr(target.Expr)

	return &ResultColumn{
		Name:      name,
		Sources:   refs,
		IsDerived: derived,
	}
}

// deriveResultName picks the output column name for a non-star target.
func deriveResultName(target *ast.SelectTarget) string {
	if !target.Alias.IsEmpty() {
		return target.Alias.Normalize()
	}
	return exprName(target.Expr)
}

// exprName returns a best-effort name for an expression.
func exprName(expr ast.Node) string {
	if expr == nil {
		return ""
	}
	switch e := expr.(type) {
	case *ast.ColumnRef:
		if len(e.Parts) > 0 {
			return e.Parts[len(e.Parts)-1].Normalize()
		}
	case *ast.FuncCallExpr:
		return e.Name.Name.Normalize()
	case *ast.CastExpr:
		return exprName(e.Expr)
	case *ast.ParenExpr:
		return exprName(e.Expr)
	case *ast.CollateExpr:
		return exprName(e.Expr)
	}
	return ""
}

// ---------------------------------------------------------------------------
// Column reference collection
// ---------------------------------------------------------------------------

// collectColumnRefs walks an expression and returns all SourceColumn references,
// resolved against the FROM entries.
func collectColumnRefs(expr ast.Node, fromEntries []tableEntry) []*SourceColumn {
	if expr == nil {
		return nil
	}
	var refs []*SourceColumn
	seen := make(map[sourceKey]bool)

	ast.Inspect(expr, func(node ast.Node) bool {
		if node == nil {
			return false
		}
		switch n := node.(type) {
		case *ast.ColumnRef:
			sc := resolveColumnRef(n, fromEntries)
			k := toSourceKey(sc)
			if !seen[k] {
				seen[k] = true
				refs = append(refs, sc)
			}
			return false // don't recurse into ColumnRef parts
		case *ast.StarExpr:
			// COUNT(*) or similar — pseudo source with no table.
			sc := &SourceColumn{Column: "*"}
			k := toSourceKey(sc)
			if !seen[k] {
				seen[k] = true
				refs = append(refs, sc)
			}
			return false
		case *ast.SubqueryExpr:
			// Subquery expressions in SELECT list: we don't trace their lineage
			// into the result column, but we do collect their table accesses.
			// For now, skip subquery lineage — just return false.
			return false
		}
		return true
	})

	return refs
}

// resolveColumnRef resolves a ColumnRef to a SourceColumn, looking up the
// table qualifier from the FROM entries when only a column name is provided.
func resolveColumnRef(ref *ast.ColumnRef, fromEntries []tableEntry) *SourceColumn {
	parts := ref.Parts
	switch len(parts) {
	case 1:
		col := parts[0].Normalize()
		// Try to find which table this column comes from.
		table := resolveColumnToTable(col, fromEntries)
		return &SourceColumn{Table: table, Column: col}
	case 2:
		return &SourceColumn{Table: parts[0].Normalize(), Column: parts[1].Normalize()}
	case 3:
		return &SourceColumn{Schema: parts[0].Normalize(), Table: parts[1].Normalize(), Column: parts[2].Normalize()}
	case 4:
		return &SourceColumn{Database: parts[0].Normalize(), Schema: parts[1].Normalize(), Table: parts[2].Normalize(), Column: parts[3].Normalize()}
	}
	return &SourceColumn{Column: strings.Join(identPartsToStrings(parts), ".")}
}

// resolveColumnToTable looks up which FROM entry contains the given column name.
// Returns the table alias if found, or empty string if ambiguous/not found.
func resolveColumnToTable(col string, fromEntries []tableEntry) string {
	colUpper := strings.ToUpper(col)
	var found string
	for _, entry := range fromEntries {
		if entry.isUnknown || len(entry.columns) == 0 {
			// Can't determine, but if there's only one FROM table assume it.
			if len(fromEntries) == 1 {
				return entry.alias
			}
			continue
		}
		for _, sc := range entry.columns {
			if sc.Column == "*" {
				// Unknown columns; can't match specifically.
				if len(fromEntries) == 1 {
					return entry.alias
				}
				continue
			}
			if strings.ToUpper(sc.Column) == colUpper {
				if found != "" {
					return "" // ambiguous
				}
				found = entry.alias
			}
		}
	}
	// If we found it in exactly one table, return it.
	// If not found at all but there's only one FROM table, assume it belongs there.
	if found == "" && len(fromEntries) == 1 {
		return fromEntries[0].alias
	}
	return found
}

// identPartsToStrings converts a slice of Ident to their normalized string forms.
func identPartsToStrings(parts []ast.Ident) []string {
	result := make([]string, len(parts))
	for i, p := range parts {
		result[i] = p.Normalize()
	}
	return result
}

// ---------------------------------------------------------------------------
// Expression classification
// ---------------------------------------------------------------------------

// isDerivedExpr returns true when the expression is not a simple column pass-through.
func isDerivedExpr(expr ast.Node) bool {
	if expr == nil {
		return false
	}
	switch e := expr.(type) {
	case *ast.ColumnRef:
		return false
	case *ast.ParenExpr:
		return isDerivedExpr(e.Expr)
	case *ast.Literal:
		return true // constants have no source
	case *ast.FuncCallExpr:
		return true
	case *ast.BinaryExpr:
		return true
	case *ast.UnaryExpr:
		return true
	case *ast.CaseExpr:
		return true
	case *ast.IffExpr:
		return true
	case *ast.CastExpr:
		// A plain cast of a column ref is not truly "derived" for lineage purposes.
		return isDerivedExpr(e.Expr)
	case *ast.CollateExpr:
		return isDerivedExpr(e.Expr)
	case *ast.SubqueryExpr:
		return true
	case *ast.StarExpr:
		return true // COUNT(*) context
	default:
		return true
	}
}

// ---------------------------------------------------------------------------
// Result column merging (set ops)
// ---------------------------------------------------------------------------

// mergeResultColumns combines two ResultColumns from set-op branches.
// If one side is nil (branch length mismatch), the other is used as-is.
func mergeResultColumns(left, right *ResultColumn) *ResultColumn {
	if left == nil {
		return right
	}
	if right == nil {
		return left
	}

	// Name comes from left branch.
	name := left.Name

	// Sources = union.
	seen := make(map[sourceKey]bool)
	var srcs []*SourceColumn
	for _, sc := range left.Sources {
		k := toSourceKey(sc)
		if !seen[k] {
			seen[k] = true
			srcs = append(srcs, sc)
		}
	}
	for _, sc := range right.Sources {
		k := toSourceKey(sc)
		if !seen[k] {
			seen[k] = true
			srcs = append(srcs, sc)
		}
	}

	derived := left.IsDerived || right.IsDerived

	return &ResultColumn{
		Name:      name,
		Sources:   srcs,
		IsDerived: derived,
	}
}

// ---------------------------------------------------------------------------
// Deduplication helpers
// ---------------------------------------------------------------------------

type sourceKey struct {
	database, schema, table, column string
}

func toSourceKey(sc *SourceColumn) sourceKey {
	if sc == nil {
		return sourceKey{}
	}
	return sourceKey{
		database: sc.Database,
		schema:   sc.Schema,
		table:    sc.Table,
		column:   sc.Column,
	}
}

func flattenSources(m map[sourceKey]*SourceColumn) []*SourceColumn {
	result := make([]*SourceColumn, 0, len(m))
	for _, sc := range m {
		result = append(result, sc)
	}
	return result
}
