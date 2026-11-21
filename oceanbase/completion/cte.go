package completion

import (
	"maps"
	"slices"
	"strings"

	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/oceanbase/catalog"
)

type resolvedTableSource struct {
	database string
	table    string
	alias    string
	cols     []string
	isCTE    bool
	remote   bool
}

func resolveColumnsASTScoped(
	cat *catalog.Catalog,
	sql string,
	cursorOffset int,
	caseSensitive bool,
	qDB string,
	qName string,
) (candidates []Candidate, handled bool) {
	stmt, cursor := cursorStatement(sql, cursorOffset)
	list, err := parseClosingTokens(stmt)
	if err != nil || list == nil || len(list.Items) == 0 {
		if cursor < 0 || cursor > len(stmt) {
			return nil, false
		}
		patched := stmt[:cursor] + " _x" + stmt[cursor:]
		list, err = parseClosingTokens(patched)
		if err != nil || list == nil || len(list.Items) == 0 {
			return nil, false
		}
	}

	var targetNode ast.Node
	for _, item := range list.Items {
		if item == nil {
			continue
		}
		targetNode = item
	}
	sel, ok := targetNode.(*ast.SelectStmt)
	if !ok || sel == nil {
		return nil, false
	}

	activeSel, ctes, outerScopes := findActiveSelectScope(sel, cursor, cat, nil, nil, caseSensitive)
	if activeSel == nil {
		return nil, false
	}

	sources := collectSelectFromSources(cat, activeSel.From, ctes, caseSensitive)
	scopes := make([][]resolvedTableSource, 0, 1+len(outerScopes))
	scopes = append(scopes, sources)
	scopes = append(scopes, outerScopes...)

	hasAnySource := false
	hasCTEOrDerived := false
	for _, scope := range scopes {
		if len(scope) > 0 {
			hasAnySource = true
		}
		for _, src := range scope {
			if src.isCTE {
				hasCTEOrDerived = true
				break
			}
		}
	}
	if !hasAnySource && len(ctes) == 0 {
		return nil, false
	}

	seen := make(map[string]bool)
	var result []Candidate
	appendCols := func(cols []string) {
		for _, col := range cols {
			if col == "" || seen[col] {
				continue
			}
			seen[col] = true
			result = append(result, Candidate{Text: col, Type: CandidateColumn})
		}
	}

	if qName != "" {
		matchedAny := false
		for _, scope := range scopes {
			scopeMatched := false
			for _, src := range scope {
				if src.remote {
					continue
				}
				if !strings.EqualFold(src.alias, qName) && !strings.EqualFold(src.table, qName) {
					continue
				}
				if qDB != "" {
					if src.isCTE {
						continue
					}
					srcDB := src.database
					if srcDB == "" {
						srcDB = cat.CurrentDatabase()
					}
					if !dbNameEquals(srcDB, qDB, caseSensitive) {
						continue
					}
				}
				scopeMatched = true
				matchedAny = true
				appendCols(src.cols)
			}
			// Inner aliases/tables shadow outer scopes of the same name.
			if scopeMatched {
				break
			}
		}
		if !matchedAny && qDB == "" {
			if cteCols, isCTE := ctes[strings.ToLower(qName)]; isCTE {
				appendCols(cteCols)
			}
		}
		return result, true
	}

	for _, scope := range scopes {
		for _, src := range scope {
			if src.remote {
				continue
			}
			appendCols(src.cols)
		}
	}
	if len(result) > 0 || hasCTEOrDerived {
		return result, true
	}
	return nil, false
}

func findActiveSelectScope(
	s *ast.SelectStmt,
	cursor int,
	cat *catalog.Catalog,
	outerCTEs map[string][]string,
	outerScopes [][]resolvedTableSource,
	caseSensitive bool,
) (*ast.SelectStmt, map[string][]string, [][]resolvedTableSource) {
	if s == nil {
		return nil, outerCTEs, outerScopes
	}
	ctes := maps.Clone(outerCTEs)
	if ctes == nil {
		ctes = make(map[string][]string)
	}
	for _, cte := range s.CTEs {
		if cte == nil {
			continue
		}
		if cte.Select != nil && locContains(cte.Select.Loc, cursor) {
			return findActiveSelectScope(cte.Select, cursor, cat, ctes, outerScopes, caseSensitive)
		}
		if cte.Name == "" {
			continue
		}
		if len(cte.Columns) > 0 {
			ctes[strings.ToLower(cte.Name)] = slices.Clone(cte.Columns)
		} else if cte.Select != nil {
			ctes[strings.ToLower(cte.Name)] = projectSelectColumns(cat, cte.Select, ctes, caseSensitive)
		}
	}

	if s.ParenSource != nil {
		return findActiveSelectScope(s.ParenSource, cursor, cat, ctes, outerScopes, caseSensitive)
	}
	if s.SetOp != ast.SetOpNone {
		if s.Right != nil && locContains(s.Right.Loc, cursor) {
			return findActiveSelectScope(s.Right, cursor, cat, ctes, outerScopes, caseSensitive)
		}
		if s.Left != nil {
			return findActiveSelectScope(s.Left, cursor, cat, ctes, outerScopes, caseSensitive)
		}
	}

	if sub, canSeeCurFrom := findSubqueryAtCursor(s, cursor); sub != nil {
		nextOuter := outerScopes
		if canSeeCurFrom && len(s.From) > 0 {
			curSources := collectSelectFromSources(cat, s.From, ctes, caseSensitive)
			if len(curSources) > 0 {
				nextOuter = make([][]resolvedTableSource, 0, 1+len(outerScopes))
				nextOuter = append(nextOuter, curSources)
				nextOuter = append(nextOuter, outerScopes...)
			}
		}
		return findActiveSelectScope(sub, cursor, cat, ctes, nextOuter, caseSensitive)
	}
	return s, ctes, outerScopes
}

func locContains(loc ast.Loc, cursor int) bool {
	return loc.Start >= 0 && loc.End >= loc.Start && cursor >= loc.Start && cursor <= loc.End
}

func findSubqueryAtCursor(s *ast.SelectStmt, cursor int) (*ast.SelectStmt, bool) {
	for _, te := range s.From {
		if sub, lateral := findSubqueryInTableExpr(te, cursor); sub != nil {
			return sub, lateral
		}
	}
	for _, e := range s.TargetList {
		if sub := findSubqueryInAST(e, cursor); sub != nil {
			return sub, true
		}
	}
	if sub := findSubqueryInAST(s.Where, cursor); sub != nil {
		return sub, true
	}
	for _, e := range s.GroupBy {
		if sub := findSubqueryInAST(e, cursor); sub != nil {
			return sub, true
		}
	}
	if sub := findSubqueryInAST(s.Having, cursor); sub != nil {
		return sub, true
	}
	for _, ob := range s.OrderBy {
		if ob != nil {
			if sub := findSubqueryInAST(ob.Expr, cursor); sub != nil {
				return sub, true
			}
		}
	}
	return nil, false
}

func findSubqueryInTableExpr(te ast.TableExpr, cursor int) (*ast.SelectStmt, bool) {
	switch v := te.(type) {
	case *ast.JoinClause:
		if sub, lateral := findSubqueryInTableExpr(v.Left, cursor); sub != nil {
			return sub, lateral
		}
		return findSubqueryInTableExpr(v.Right, cursor)
	case *ast.SubqueryExpr:
		if v.Select != nil && locContains(v.Loc, cursor) {
			return v.Select, v.Lateral
		}
	}
	return nil, false
}

func findSubqueryInAST(node ast.Node, cursor int) *ast.SelectStmt {
	if node == nil {
		return nil
	}
	var found *ast.SelectStmt
	ast.Inspect(node, func(n ast.Node) bool {
		if found != nil || n == nil {
			return false
		}
		switch v := n.(type) {
		case *ast.SubqueryExpr:
			if v.Select != nil && locContains(v.Loc, cursor) {
				found = v.Select
			}
			return false
		case *ast.ExistsExpr:
			if v.Select != nil && locContains(v.Loc, cursor) {
				found = v.Select
			}
			return false
		}
		return true
	})
	return found
}

const maxProjectedCTEColumns = 4096

func projectSelectColumns(
	cat *catalog.Catalog,
	s *ast.SelectStmt,
	ctes map[string][]string,
	caseSensitive bool,
) []string {
	leaf := ast.LeftmostQueryLeaf(s)
	if leaf == nil {
		return nil
	}
	localCTEs := maps.Clone(ctes)
	if localCTEs == nil {
		localCTEs = make(map[string][]string)
	}
	for _, cte := range leaf.CTEs {
		if cte == nil || cte.Name == "" {
			continue
		}
		if len(cte.Columns) > 0 {
			localCTEs[strings.ToLower(cte.Name)] = slices.Clone(cte.Columns)
		} else if cte.Select != nil {
			localCTEs[strings.ToLower(cte.Name)] = projectSelectColumns(cat, cte.Select, localCTEs, caseSensitive)
		}
	}

	sources := collectSelectFromSources(cat, leaf.From, localCTEs, caseSensitive)
	seen := make(map[string]bool)
	var cols []string
	addCol := func(name string) {
		if name == "" || len(cols) >= maxProjectedCTEColumns {
			return
		}
		key := strings.ToLower(name)
		if seen[key] {
			return
		}
		seen[key] = true
		cols = append(cols, name)
	}
	addCols := func(list []string) {
		for _, name := range list {
			addCol(name)
			if len(cols) >= maxProjectedCTEColumns {
				return
			}
		}
	}
	appendExpr := func(alias string, val ast.ExprNode) {
		if alias != "" {
			addCol(alias)
			return
		}
		switch v := val.(type) {
		case *ast.StarExpr:
			for _, src := range sources {
				addCols(src.cols)
			}
		case *ast.ColumnRef:
			if v.Star {
				for _, src := range sources {
					if v.Table == "" || strings.EqualFold(src.alias, v.Table) || strings.EqualFold(src.table, v.Table) {
						addCols(src.cols)
					}
				}
				return
			}
			if v.Column != "" && v.Column != "_x" {
				addCol(v.Column)
			}
		}
	}

	for _, target := range leaf.TargetList {
		switch t := target.(type) {
		case *ast.ResTarget:
			appendExpr(t.Name, t.Val)
		default:
			appendExpr("", target)
		}
	}
	return cols
}

func collectSelectFromSources(
	cat *catalog.Catalog,
	from []ast.TableExpr,
	ctes map[string][]string,
	caseSensitive bool,
) []resolvedTableSource {
	var sources []resolvedTableSource
	for _, te := range from {
		collectTableExprSources(cat, te, ctes, caseSensitive, &sources)
	}
	return sources
}

func collectTableExprSources(
	cat *catalog.Catalog,
	te ast.TableExpr,
	ctes map[string][]string,
	caseSensitive bool,
	out *[]resolvedTableSource,
) {
	if te == nil {
		return
	}
	switch v := te.(type) {
	case *ast.JoinClause:
		collectTableExprSources(cat, v.Left, ctes, caseSensitive, out)
		collectTableExprSources(cat, v.Right, ctes, caseSensitive, out)
	case *ast.SubqueryExpr:
		var cols []string
		if len(v.Columns) > 0 {
			cols = slices.Clone(v.Columns)
		} else if v.Select != nil {
			cols = projectSelectColumns(cat, v.Select, ctes, caseSensitive)
		}
		*out = append(*out, resolvedTableSource{
			alias: v.Alias,
			cols:  cols,
			isCTE: true,
		})
	case *ast.TableRef:
		if v.Name == "" {
			return
		}
		if isRemoteASTTableRef(v) {
			*out = append(*out, resolvedTableSource{
				database: v.Schema,
				table:    v.Name,
				alias:    v.Alias,
				remote:   true,
			})
			return
		}
		if v.Schema == "" && ctes != nil {
			if cteCols, ok := ctes[strings.ToLower(v.Name)]; ok {
				*out = append(*out, resolvedTableSource{
					table: v.Name,
					alias: v.Alias,
					cols:  cteCols,
					isCTE: true,
				})
				return
			}
		}
		var db *catalog.Database
		if v.Schema != "" {
			db = getDatabase(cat, v.Schema, caseSensitive)
		} else {
			db = currentDB(cat)
		}
		var cols []string
		if db != nil {
			cols = tableOrViewColumns(db, v.Name)
		}
		*out = append(*out, resolvedTableSource{
			database: v.Schema,
			table:    v.Name,
			alias:    v.Alias,
			cols:     cols,
		})
	}
}

func tableOrViewColumns(db *catalog.Database, name string) []string {
	for _, t := range db.Tables {
		if t.Name == name {
			cols := make([]string, 0, len(t.Columns))
			for _, col := range t.Columns {
				cols = append(cols, col.Name)
			}
			return cols
		}
	}
	for _, v := range db.Views {
		if v.Name == name {
			return slices.Clone(v.Columns)
		}
	}
	return nil
}

func isRemoteASTTableRef(v *ast.TableRef) bool {
	return v.Catalog != "" || v.Dblink != ""
}
