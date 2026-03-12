package completion

import (
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/yacc"
)

// collectAllRefs extracts table refs, CTEs and aliases from the SQL.
// It tries AST-based extraction first, then falls back to heuristics.
func collectAllRefs(fullSQL string, cursorOffset int) ([]tableRef, []cteInfo, []aliasInfo) {
	cleanSQL := removeCursorMarker(fullSQL)
	if cursorOffset > len(cleanSQL) {
		cursorOffset = len(cleanSQL)
	}

	// First, isolate the cursor's statement for multi-statement SQL
	stmtSQL, stmtOffset := isolateCursorStmt(cleanSQL, cursorOffset)

	// Strategy 1: repair and parse the isolated statement
	fixed := repairSQL(stmtSQL, stmtOffset)
	if fixed != "" {
		refs, ctes, ok := parseAndExtractRefs(fixed)
		if ok {
			return refs, ctes, buildAliases(refs)
		}
	}

	// Strategy 2: parse the isolated statement as-is
	refs, ctes, ok := parseAndExtractRefs(stmtSQL)
	if ok {
		return refs, ctes, buildAliases(refs)
	}

	// Strategy 3: heuristic extraction on the isolated statement
	refs = heuristicExtractRefs(stmtSQL)
	ctes = heuristicExtractCTEs(stmtSQL)
	return refs, ctes, buildAliases(refs)
}

// repairSQL attempts various fixes to make incomplete SQL parseable.
func repairSQL(sql string, cursorOffset int) string {
	if cursorOffset < 0 || cursorOffset > len(sql) {
		return ""
	}

	before := sql[:cursorOffset]
	after := sql[cursorOffset:]

	// Fix: "SELECT ... ORDER BY |" → "SELECT ... ORDER BY 1"
	trimBefore := strings.TrimRight(before, " \t\n")
	upperBefore := strings.ToUpper(trimBefore)
	if strings.HasSuffix(upperBefore, "ORDER BY") ||
		strings.HasSuffix(upperBefore, "WHERE") ||
		strings.HasSuffix(upperBefore, "AND") ||
		strings.HasSuffix(upperBefore, "OR") ||
		strings.HasSuffix(upperBefore, "ON") ||
		strings.HasSuffix(upperBefore, "HAVING") {
		return before + " 1 " + after
	}

	// Fix: "SELECT alias.| FROM ..." → "SELECT alias.__x__ FROM ..."
	if strings.HasSuffix(trimBefore, ".") {
		return before + "__x__" + after
	}

	// Fix: "UPDATE t1 SET |" → "UPDATE t1 SET __x__ = 1"
	if strings.HasSuffix(upperBefore, "SET") {
		return before + " __x__ = 1 " + after
	}

	// Fix: "SELECT | FROM ..." → "SELECT __x__ FROM ..."
	if strings.HasSuffix(upperBefore, "SELECT") ||
		strings.HasSuffix(upperBefore, ",") ||
		strings.HasSuffix(upperBefore, "(") {
		return before + " __x__ " + after
	}

	// Fix: "SELECT * FROM |" / "INSERT INTO |" → "... __x__"
	if strings.HasSuffix(upperBefore, "FROM") ||
		strings.HasSuffix(upperBefore, "JOIN") ||
		strings.HasSuffix(upperBefore, "INTO") ||
		strings.HasSuffix(upperBefore, "UPDATE") {
		return before + " __x__ " + after
	}

	// Generic: try inserting a placeholder
	return before + " __x__ " + after
}

// isolateCursorStmt extracts the statement containing the cursor.
// It splits on semicolons first, and additionally on newline-separated
// statement keywords to handle cases without semicolons between statements.
func isolateCursorStmt(sql string, cursorOffset int) (string, int) {
	// First split by semicolons
	stmts := splitStatements(sql)
	pos := 0
	targetStmt := ""
	targetOffset := cursorOffset
	for _, stmt := range stmts {
		stmtEnd := pos + len(stmt)
		if cursorOffset >= pos && cursorOffset <= stmtEnd {
			targetStmt = stmt
			targetOffset = cursorOffset - pos
			break
		}
		pos = stmtEnd + 1
	}
	if targetStmt == "" {
		if len(stmts) > 0 {
			targetStmt = stmts[len(stmts)-1]
			targetOffset = cursorOffset - (len(sql) - len(targetStmt))
		} else {
			return sql, cursorOffset
		}
	}

	// Within the semicolon-delimited segment, further split by newline + statement keyword.
	// This handles cases like "SELECT * FROM |\nselect * from ..."
	trimmed := targetStmt
	if targetOffset >= 0 && targetOffset <= len(trimmed) {
		// Find line breaks that start new statements after the cursor
		afterCursor := trimmed[targetOffset:]
		nlIdx := strings.IndexByte(afterCursor, '\n')
		if nlIdx >= 0 {
			// Check if the text after the newline starts a new statement
			rest := strings.TrimLeft(afterCursor[nlIdx+1:], " \t")
			upperRest := strings.ToUpper(rest)
			if strings.HasPrefix(upperRest, "SELECT") ||
				strings.HasPrefix(upperRest, "INSERT") ||
				strings.HasPrefix(upperRest, "UPDATE") ||
				strings.HasPrefix(upperRest, "DELETE") ||
				strings.HasPrefix(upperRest, "CREATE") ||
				strings.HasPrefix(upperRest, "ALTER") ||
				strings.HasPrefix(upperRest, "DROP") {
				// Truncate at the newline
				trimmed = trimmed[:targetOffset+nlIdx]
			}
		}
	}

	result := strings.TrimSpace(trimmed)
	// Adjust offset for any leading whitespace removed by TrimSpace.
	leading := len(trimmed) - len(strings.TrimLeft(trimmed, " \t\n\r"))
	adjustedOffset := targetOffset - leading
	if adjustedOffset < 0 {
		adjustedOffset = 0
	}
	return result, adjustedOffset
}

// heuristicExtractRefs extracts table references using text patterns.
func heuristicExtractRefs(sql string) []tableRef {
	var refs []tableRef
	tokens := tokenizeForRefs(sql)

	for i := 0; i < len(tokens); i++ {
		upper := strings.ToUpper(tokens[i])

		// INSERT INTO tablename / UPDATE tablename
		if upper == "INTO" || upper == "UPDATE" {
			if i+1 < len(tokens) {
				i++
				ref := parseTableName(tokens[i])
				// Check for alias
				if i+1 < len(tokens) {
					nextUp := strings.ToUpper(tokens[i+1])
					if nextUp == "AS" && i+2 < len(tokens) {
						ref.alias = tokens[i+2]
					} else if !isKeyword(nextUp) && nextUp != "(" && nextUp != "SET" && nextUp != "VALUES" {
						ref.alias = tokens[i+1]
					}
				}
				refs = append(refs, ref)
			}
			continue
		}

		if upper == "FROM" || upper == "JOIN" ||
			upper == "INNER" || upper == "LEFT" ||
			upper == "RIGHT" || upper == "CROSS" ||
			upper == "FULL" {
			// Skip JOIN keyword after direction
			if upper != "FROM" && upper != "JOIN" {
				// Look for JOIN
				if i+1 < len(tokens) && strings.ToUpper(tokens[i+1]) == "JOIN" {
					i++
				} else {
					continue
				}
			}

			// Next token should be a table reference
			i++
			if i >= len(tokens) {
				break
			}
			tableName := tokens[i]
			if tableName == "(" || tableName == "LATERAL" || strings.ToUpper(tableName) == "LATERAL" {
				// Subquery — skip, look for alias after closing paren
				if tableName == "(" || strings.ToUpper(tableName) == "LATERAL" {
					depth := 0
					if tableName == "(" {
						depth = 1
					}
					i++
					for i < len(tokens) {
						if tokens[i] == "(" {
							depth++
						} else if tokens[i] == ")" {
							depth--
							if depth == 0 {
								break
							}
						}
						i++
					}
					// Look for alias after )
					if i+1 < len(tokens) {
						next := tokens[i+1]
						nextUpper := strings.ToUpper(next)
						if !isKeyword(nextUpper) || isIdentKeyword(nextUpper) {
							if nextUpper != "ON" && nextUpper != "WHERE" &&
								nextUpper != "JOIN" && nextUpper != "LEFT" &&
								nextUpper != "RIGHT" && nextUpper != "INNER" &&
								nextUpper != "CROSS" && nextUpper != "FULL" &&
								nextUpper != "ORDER" && nextUpper != "GROUP" &&
								nextUpper != "HAVING" && nextUpper != "LIMIT" &&
								nextUpper != "UNION" && nextUpper != ";" {
								i++
								refs = append(refs, tableRef{
									table:  tokens[i],
									alias:  tokens[i],
									isSubq: true,
								})
							}
						}
					}
				}
				continue
			}

			ref := parseTableName(tableName)

			// Check for alias
			if i+1 < len(tokens) {
				next := tokens[i+1]
				nextUpper := strings.ToUpper(next)
				if nextUpper == "AS" {
					i += 2
					if i < len(tokens) {
						ref.alias = tokens[i]
					}
				} else if !isSQLKeyword(nextUpper) && next != "," && next != "(" && next != ")" && next != ";" {
					i++
					ref.alias = tokens[i]
				}
			}
			refs = append(refs, ref)
		}
	}

	return refs
}

func parseTableName(name string) tableRef {
	parts := strings.Split(name, ".")
	switch len(parts) {
	case 2:
		return tableRef{schema: strings.ToLower(parts[0]), table: strings.ToLower(parts[1])}
	case 1:
		return tableRef{table: strings.ToLower(parts[0])}
	}
	return tableRef{table: strings.ToLower(name)}
}

// heuristicExtractCTEs extracts CTE info using text patterns.
func heuristicExtractCTEs(sql string) []cteInfo {
	upperSQL := strings.ToUpper(sql)
	if !strings.Contains(upperSQL, "WITH") {
		return nil
	}

	// Try to parse just the WITH clause portion
	// Find "WITH ... AS (" and try parsing
	var ctes []cteInfo
	tokens := tokenizeForRefs(sql)

	for i := 0; i < len(tokens); i++ {
		if strings.ToUpper(tokens[i]) != "WITH" {
			continue
		}
		i++
		for i < len(tokens) {
			// CTE name
			if i >= len(tokens) {
				break
			}
			cteName := tokens[i]
			cte := cteInfo{name: cteName}
			i++

			// Optional column list: name(c1, c2)
			if i < len(tokens) && tokens[i] == "(" {
				i++
				for i < len(tokens) && tokens[i] != ")" {
					name := tokens[i]
					if name != "," {
						cte.columns = append(cte.columns, name)
					}
					i++
				}
				if i < len(tokens) {
					i++ // skip )
				}
			}

			// AS (
			if i < len(tokens) && strings.ToUpper(tokens[i]) == "AS" {
				i++
			}
			if i < len(tokens) && tokens[i] == "(" {
				// Find matching close paren
				depth := 1
				queryStart := i + 1
				i++
				for i < len(tokens) && depth > 0 {
					if tokens[i] == "(" {
						depth++
					} else if tokens[i] == ")" {
						depth--
					}
					i++
				}
				queryEnd := i - 1
				if queryEnd > queryStart {
					querySQL := strings.Join(tokens[queryStart:queryEnd], " ")
					// Try to parse the CTE query
					result, err := yacc.Parse(querySQL)
					if err == nil && result != nil && len(result.Items) > 0 {
						if raw, ok := result.Items[0].(*nodes.RawStmt); ok {
							if sel, ok := raw.Stmt.(*nodes.SelectStmt); ok {
								cte.selStmt = sel
							}
						}
					}
				}
			}

			ctes = append(ctes, cte)

			// Check for comma (another CTE) or SELECT
			if i < len(tokens) {
				if tokens[i] == "," {
					i++
					continue
				}
			}
			break
		}
		break
	}

	return ctes
}

// tokenizeForRefs does simple whitespace/punctuation tokenization for heuristic extraction.
func tokenizeForRefs(sql string) []string {
	var tokens []string
	var current strings.Builder
	inSingle := false
	inDouble := false

	flush := func() {
		if current.Len() > 0 {
			tokens = append(tokens, current.String())
			current.Reset()
		}
	}

	for i := 0; i < len(sql); i++ {
		ch := sql[i]
		switch {
		case ch == '\'' && !inDouble:
			inSingle = !inSingle
			current.WriteByte(ch)
		case ch == '"' && !inSingle:
			inDouble = !inDouble
			current.WriteByte(ch)
		case inSingle || inDouble:
			current.WriteByte(ch)
		case ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r':
			flush()
		case ch == '(' || ch == ')' || ch == ',' || ch == ';':
			flush()
			tokens = append(tokens, string(ch))
		case ch == '.':
			// Attach dot to the previous token (schema.table)
			if current.Len() > 0 {
				current.WriteByte(ch)
			} else if len(tokens) > 0 {
				// Attach to previous token
				tokens[len(tokens)-1] += "."
			}
		default:
			current.WriteByte(ch)
		}
	}
	flush()
	return tokens
}

func isSQLKeyword(s string) bool {
	switch s {
	case "SELECT", "FROM", "WHERE", "JOIN", "ON", "AND", "OR", "NOT",
		"ORDER", "BY", "GROUP", "HAVING", "LIMIT", "OFFSET",
		"UNION", "INTERSECT", "EXCEPT", "AS", "IN", "IS",
		"NULL", "TRUE", "FALSE", "BETWEEN", "LIKE", "ILIKE",
		"EXISTS", "CASE", "WHEN", "THEN", "ELSE", "END",
		"INSERT", "INTO", "VALUES", "UPDATE", "SET", "DELETE",
		"CREATE", "ALTER", "DROP", "TABLE", "INDEX", "VIEW",
		"LEFT", "RIGHT", "INNER", "OUTER", "CROSS", "FULL",
		"WITH", "RECURSIVE", "LATERAL", "DISTINCT",
		"ALL", "ANY", "SOME", "FOR", "FETCH", "WINDOW":
		return true
	}
	return false
}

func isKeyword(s string) bool {
	return isSQLKeyword(s)
}

func isIdentKeyword(s string) bool {
	// Keywords that can also be used as identifiers
	switch s {
	case "NAME", "TYPE", "VALUE", "DATA", "ROLE", "ACTION":
		return true
	}
	return false
}

// parseAndExtractRefs parses SQL and extracts table references and CTEs.
func parseAndExtractRefs(sql string) ([]tableRef, []cteInfo, bool) {
	result, err := yacc.Parse(sql)
	if err != nil || result == nil || len(result.Items) == 0 {
		return nil, nil, false
	}

	var allRefs []tableRef
	var allCTEs []cteInfo

	for _, item := range result.Items {
		// Items may be *nodes.RawStmt or directly *nodes.SelectStmt, etc.
		var stmt nodes.Node
		if raw, ok := item.(*nodes.RawStmt); ok {
			stmt = raw.Stmt
		} else {
			stmt = item
		}
		if stmt == nil {
			continue
		}
		refs, ctes := extractRefsFromStmt(stmt)
		allRefs = append(allRefs, refs...)
		allCTEs = append(allCTEs, ctes...)
	}

	return allRefs, allCTEs, len(allRefs) > 0 || len(allCTEs) > 0
}

func extractRefsFromStmt(stmt nodes.Node) ([]tableRef, []cteInfo) {
	switch s := stmt.(type) {
	case *nodes.SelectStmt:
		return extractRefsFromSelect(s)
	case *nodes.InsertStmt:
		return extractRefsFromInsert(s)
	case *nodes.UpdateStmt:
		return extractRefsFromUpdate(s)
	case *nodes.DeleteStmt:
		return extractRefsFromDelete(s)
	default:
		return nil, nil
	}
}

func extractRefsFromInsert(ins *nodes.InsertStmt) ([]tableRef, []cteInfo) {
	var refs []tableRef
	var ctes []cteInfo
	if ins.Relation != nil {
		refs = append(refs, rangeVarToRef(ins.Relation))
	}
	if ins.WithClause != nil {
		ctes = extractWithClauseCTEs(ins.WithClause)
	}
	return refs, ctes
}

func extractRefsFromUpdate(upd *nodes.UpdateStmt) ([]tableRef, []cteInfo) {
	var refs []tableRef
	var ctes []cteInfo
	if upd.Relation != nil {
		refs = append(refs, rangeVarToRef(upd.Relation))
	}
	if upd.FromClause != nil {
		refs = append(refs, extractFromClauseRefs(upd.FromClause)...)
	}
	if upd.WithClause != nil {
		ctes = extractWithClauseCTEs(upd.WithClause)
	}
	return refs, ctes
}

func extractRefsFromDelete(del *nodes.DeleteStmt) ([]tableRef, []cteInfo) {
	var refs []tableRef
	var ctes []cteInfo
	if del.Relation != nil {
		refs = append(refs, rangeVarToRef(del.Relation))
	}
	if del.UsingClause != nil {
		refs = append(refs, extractFromClauseRefs(del.UsingClause)...)
	}
	if del.WithClause != nil {
		ctes = extractWithClauseCTEs(del.WithClause)
	}
	return refs, ctes
}

func rangeVarToRef(rv *nodes.RangeVar) tableRef {
	ref := tableRef{table: rv.Relname}
	if rv.Schemaname != "" {
		ref.schema = rv.Schemaname
	}
	if rv.Alias != nil && rv.Alias.Aliasname != "" {
		ref.alias = rv.Alias.Aliasname
	}
	return ref
}

func extractWithClauseCTEs(wc *nodes.WithClause) []cteInfo {
	if wc == nil || wc.Ctes == nil {
		return nil
	}
	var ctes []cteInfo
	for _, item := range wc.Ctes.Items {
		cte, ok := item.(*nodes.CommonTableExpr)
		if !ok || cte == nil {
			continue
		}
		ci := cteInfo{name: cte.Ctename}
		if cte.Aliascolnames != nil {
			for _, col := range cte.Aliascolnames.Items {
				if s, ok := col.(*nodes.String); ok {
					ci.columns = append(ci.columns, s.Str)
				}
			}
		}
		if sel, ok := cte.Ctequery.(*nodes.SelectStmt); ok {
			ci.selStmt = sel
		}
		ctes = append(ctes, ci)
	}
	return ctes
}

func extractRefsFromSelect(sel *nodes.SelectStmt) ([]tableRef, []cteInfo) {
	var refs []tableRef
	var ctes []cteInfo

	if sel.WithClause != nil && sel.WithClause.Ctes != nil {
		for _, item := range sel.WithClause.Ctes.Items {
			cteNode, ok := item.(*nodes.CommonTableExpr)
			if !ok {
				continue
			}
			ci := cteInfo{name: cteNode.Ctename}
			if cteNode.Aliascolnames != nil {
				for _, col := range cteNode.Aliascolnames.Items {
					if s, ok := col.(*nodes.String); ok {
						ci.columns = append(ci.columns, s.Str)
					}
				}
			}
			if cteNode.Ctequery != nil {
				if innerSel, ok := cteNode.Ctequery.(*nodes.SelectStmt); ok {
					ci.selStmt = innerSel
				}
			}
			ctes = append(ctes, ci)
			refs = append(refs, tableRef{table: cteNode.Ctename})
		}
	}

	if sel.FromClause != nil {
		fromRefs := extractFromClauseRefs(sel.FromClause)
		refs = append(refs, fromRefs...)
	}

	return refs, ctes
}

func extractFromClauseRefs(fromClause *nodes.List) []tableRef {
	if fromClause == nil {
		return nil
	}
	var refs []tableRef
	for _, item := range fromClause.Items {
		refs = append(refs, extractNodeRefs(item)...)
	}
	return refs
}

func extractNodeRefs(n nodes.Node) []tableRef {
	if n == nil {
		return nil
	}
	switch v := n.(type) {
	case *nodes.RangeVar:
		ref := tableRef{schema: v.Schemaname, table: v.Relname}
		if v.Alias != nil {
			ref.alias = v.Alias.Aliasname
		}
		return []tableRef{ref}

	case *nodes.JoinExpr:
		var refs []tableRef
		refs = append(refs, extractNodeRefs(v.Larg)...)
		refs = append(refs, extractNodeRefs(v.Rarg)...)
		return refs

	case *nodes.RangeSubselect:
		ref := tableRef{isSubq: true}
		if v.Alias != nil {
			ref.alias = v.Alias.Aliasname
			ref.table = v.Alias.Aliasname
		}
		if v.Subquery != nil {
			if subSel, ok := v.Subquery.(*nodes.SelectStmt); ok {
				ref.subqCols = extractTargetNames(subSel)
				// Also capture inner FROM refs for resolving SELECT *
				if subSel.FromClause != nil {
					ref.subqFromRef = extractFromClauseRefs(subSel.FromClause)
				}
			}
		}
		return []tableRef{ref}

	case *nodes.RangeFunction:
		if v.Alias != nil {
			return []tableRef{{table: v.Alias.Aliasname, alias: v.Alias.Aliasname}}
		}
	}
	return nil
}

func extractTargetNames(sel *nodes.SelectStmt) []string {
	if sel == nil || sel.TargetList == nil {
		return nil
	}
	var names []string
	hasStar := false
	for _, item := range sel.TargetList.Items {
		rt, ok := item.(*nodes.ResTarget)
		if !ok {
			continue
		}
		if rt.Name != "" {
			names = append(names, rt.Name)
			continue
		}
		if cr, ok := rt.Val.(*nodes.ColumnRef); ok {
			if cr.Fields != nil && len(cr.Fields.Items) > 0 {
				if _, isStar := cr.Fields.Items[len(cr.Fields.Items)-1].(*nodes.A_Star); isStar {
					hasStar = true
					continue
				}
			}
			name := columnRefLastName(cr)
			if name != "" {
				names = append(names, name)
			}
		}
	}

	// For SELECT * subqueries, return nil to signal "all columns from FROM tables"
	if hasStar && len(names) == 0 {
		return nil
	}
	return names
}

func buildAliases(refs []tableRef) []aliasInfo {
	var aliases []aliasInfo
	for _, r := range refs {
		if r.alias != "" {
			aliases = append(aliases, aliasInfo{alias: r.alias, schema: r.schema, table: r.table})
		}
	}
	return aliases
}

func splitStatements(sql string) []string {
	var stmts []string
	var current strings.Builder
	inSingle := false
	inDouble := false

	for i := 0; i < len(sql); i++ {
		ch := sql[i]
		switch ch {
		case '\'':
			if !inDouble {
				inSingle = !inSingle
			}
			current.WriteByte(ch)
		case '"':
			if !inSingle {
				inDouble = !inDouble
			}
			current.WriteByte(ch)
		case ';':
			if !inSingle && !inDouble {
				stmts = append(stmts, current.String())
				current.Reset()
			} else {
				current.WriteByte(ch)
			}
		default:
			current.WriteByte(ch)
		}
	}
	if current.Len() > 0 {
		stmts = append(stmts, current.String())
	}
	return stmts
}
