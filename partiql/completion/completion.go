// Package completion provides auto-complete candidates for PartiQL queries.
package completion

import (
	"sort"
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/partiql"
	"github.com/bytebase/omni/partiql/catalog"
)

// Candidate represents one auto-complete suggestion.
type Candidate struct {
	Text string // the completion text
	Kind string // "keyword", "table", "column", etc.
}

// Complete returns auto-complete candidates for the given input at the
// cursor position. Uses the catalog for table names and a built-in
// keyword list for SQL keywords.
func Complete(input string, pos int, cat *catalog.Catalog) []Candidate {
	return CompleteWithMetadata(input, pos, cat, nil)
}

// CompleteWithMetadata returns auto-complete candidates for input at byte
// offset pos, resolving table names from cat and/or meta and column names from
// meta. Both cat and meta may be nil.
func CompleteWithMetadata(input string, pos int, cat *catalog.Catalog, meta *metadata.DatabaseSchemaMetadata) []Candidate {
	// Clamp pos into [0, len(input)] before any slicing so an out-of-range
	// caller offset never panics.
	if pos > len(input) {
		pos = len(input)
	}
	if pos < 0 {
		pos = 0
	}

	stmtSQL, localOffset := extractCurrentStatement(input, pos)
	prefix, prefixStart, prefixQuoted := extractPrefixAndQuote(stmtSQL, localOffset)
	inFromContext, inColumnContext, suppressed := analyzeCursorContext(stmtSQL, prefixStart)
	if suppressed {
		return nil
	}

	var candidates []Candidate
	seen := make(map[string]struct{})
	addCandidate := func(text, kind string) {
		key := kind + "\x00" + text
		if _, ok := seen[key]; ok {
			return
		}
		seen[key] = struct{}{}
		candidates = append(candidates, Candidate{Text: text, Kind: kind})
	}

	if inFromContext {
		for _, table := range collectTableNames(cat, meta) {
			if matchesPrefix(table, prefix, prefixQuoted) {
				addCandidate(table, "table")
			}
		}
	}

	if meta != nil && inColumnContext {
		for _, col := range resolveColumnNames(stmtSQL, localOffset, meta) {
			if matchesPrefix(col, prefix, prefixQuoted) {
				addCandidate(col, "column")
			}
		}
	}

	for _, kw := range keywords {
		if matchesPrefix(kw, prefix, prefixQuoted) {
			addCandidate(kw, "keyword")
		}
	}

	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].Kind != candidates[j].Kind {
			return candidates[i].Kind < candidates[j].Kind
		}
		return candidates[i].Text < candidates[j].Text
	})
	return candidates
}

func collectTableNames(cat *catalog.Catalog, meta *metadata.DatabaseSchemaMetadata) []string {
	seen := make(map[string]struct{})
	var tables []string
	add := func(name string) {
		if name == "" {
			return
		}
		if _, ok := seen[name]; ok {
			return
		}
		seen[name] = struct{}{}
		tables = append(tables, name)
	}
	if cat != nil {
		for _, t := range cat.Tables() {
			add(t)
		}
	}
	for _, s := range meta.GetSchemas() {
		for _, t := range s.GetTables() {
			add(t.GetName())
		}
	}
	return tables
}

func extractPrefixAndQuote(input string, pos int) (prefix string, start int, quoted bool) {
	if pos > len(input) {
		pos = len(input)
	}
	if pos < 0 {
		pos = 0
	}
	if qStart, ok := unclosedDoubleQuoteStart(input, pos); ok {
		return strings.ReplaceAll(input[qStart+1:pos], `""`, `"`), qStart, true
	}
	prefix = extractPrefix(input, pos)
	return prefix, pos - len(prefix), false
}

func extractPrefix(input string, pos int) string {
	if pos > len(input) {
		pos = len(input)
	}
	if pos < 0 {
		pos = 0
	}
	start := pos
	for start > 0 && isIdentChar(input[start-1]) {
		start--
	}
	return input[start:pos]
}

func isIdentChar(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') ||
		(ch >= '0' && ch <= '9') || ch == '_' || ch == '$'
}

type clauseTok struct {
	upper    string
	afterDot bool
}

// analyzeCursorContext determines whether the cursor is in a table-introducing
// clause (FROM / JOIN / INTO / UPDATE) or an expression/projection position
// where columns are valid, or suppressed inside an unfinished comment or
// single-quoted literal. Expression-nested FROM keywords (such as TRIM(' ' FROM ...)
// or SUBSTRING(x FROM ...)) stay inside expression parentheses and do not
// suppress column completion or trigger table suggestions.
func analyzeCursorContext(statement string, prefixStart int) (inFrom, inColumn, suppressed bool) {
	if prefixStart > len(statement) {
		prefixStart = len(statement)
	}
	if prefixStart < 0 {
		prefixStart = 0
	}
	activeStart := findActiveQueryScopeStart(statement, prefixStart)
	slice := isolateSetOperationArm(statement[activeStart:prefixStart], prefixStart-activeStart)

	var clauseToks []clauseTok
	type parenFrame struct {
		tableParen bool
	}
	var stack []parenFrame
	exprDepth := 0
	prevDot := false

	i := 0
	for i < len(slice) {
		c := slice[i]
		switch {
		case c == ' ' || c == '\t' || c == '\n' || c == '\r':
			i++
		case c == '-' && i+1 < len(slice) && slice[i+1] == '-':
			end := skipLineComment(slice, i)
			if slice[end] != '\n' && slice[end] != '\r' {
				return false, false, true
			}
			i = end + 1
		case c == '/' && i+1 < len(slice) && slice[i+1] == '*':
			end, closed := skipBlockComment(slice, i)
			if !closed {
				return false, false, true
			}
			i = end + 1
		case c == '\'' || c == '"':
			end, closed := skipQuotedLiteralWithState(slice, i, c)
			if !closed {
				if c == '\'' {
					// Unfinished single-quoted literal must suppress keyword/column suggestions.
					return false, false, true
				}
				// Unfinished double-quoted identifier being typed at the cursor
				// must preserve the preceding clause context (e.g. FROM "Mu|).
				i = end + 1
				continue
			}
			if exprDepth == 0 {
				clauseToks = append(clauseToks, clauseTok{upper: "_", afterDot: prevDot})
			}
			prevDot = false
			i = end + 1
		case c == '.':
			prevDot = true
			i++
		case c == '(':
			isTableParen := false
			if exprDepth == 0 && len(clauseToks) > 0 {
				last := clauseToks[len(clauseToks)-1]
				if !last.afterDot && (last.upper == "FROM" || last.upper == "JOIN") {
					isTableParen = true
				}
			}
			stack = append(stack, parenFrame{tableParen: isTableParen})
			if !isTableParen {
				exprDepth++
			}
			prevDot = false
			i++
		case c == ')':
			if len(stack) > 0 {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if !top.tableParen && exprDepth > 0 {
					exprDepth--
				}
			}
			if exprDepth == 0 {
				clauseToks = append(clauseToks, clauseTok{upper: ")"})
			}
			prevDot = false
			i++
		case c == ',':
			if exprDepth == 0 {
				clauseToks = append(clauseToks, clauseTok{upper: ","})
			}
			prevDot = false
			i++
		case isIdentChar(c):
			start := i
			for i < len(slice) && isIdentChar(slice[i]) {
				i++
			}
			if exprDepth == 0 {
				clauseToks = append(clauseToks, clauseTok{
					upper:    strings.ToUpper(slice[start:i]),
					afterDot: prevDot,
				})
			}
			prevDot = false
		default:
			prevDot = false
			i++
		}
	}

	if exprDepth > 0 && isDDLIndexPathContext(clauseToks) {
		return false, true, false
	}

	if !prevDot && exprDepth == 0 && len(clauseToks) > 0 {
		last := clauseToks[len(clauseToks)-1]
		if !last.afterDot {
			isConflictDoUpdate := last.upper == "UPDATE" && len(clauseToks) >= 2 &&
				!clauseToks[len(clauseToks)-2].afterDot && clauseToks[len(clauseToks)-2].upper == "DO"
			isDropTable := last.upper == "TABLE" && len(clauseToks) >= 2 &&
				!clauseToks[len(clauseToks)-2].afterDot && clauseToks[len(clauseToks)-2].upper == "DROP"
			if last.upper == "FROM" || last.upper == "JOIN" || last.upper == "INTO" ||
				(last.upper == "UPDATE" && !isConflictDoUpdate) ||
				isDropTable ||
				(last.upper == "ON" && isDDLIndexOn(clauseToks)) {
				return true, false, false
			}
			if last.upper == "," {
				for idx := len(clauseToks) - 2; idx >= 0; idx-- {
					tok := clauseToks[idx]
					if tok.afterDot {
						continue
					}
					switch tok.upper {
					case "FROM", "JOIN":
						return true, false, false
					case "SELECT", "PIVOT", "AT", "LET", "REMOVE", "WHERE", "SET", "GROUP", "HAVING", "ORDER", "VALUES", "RETURNING", "INTO", "UPDATE":
						idx = -1
					}
				}
			}
		}
	}

	queryHeadIdx := -1
	fromIdx := -1
	for idx, tok := range clauseToks {
		if tok.afterDot {
			continue
		}
		switch tok.upper {
		case "SELECT", "PIVOT":
			queryHeadIdx = idx
		case "FROM":
			fromIdx = idx
		}
	}
	if queryHeadIdx >= 0 && (fromIdx < 0 || fromIdx < queryHeadIdx) {
		return false, true, false
	}

	for idx := len(clauseToks) - 1; idx >= 0; idx-- {
		if idx <= fromIdx {
			break
		}
		tok := clauseToks[idx]
		if tok.afterDot {
			continue
		}
		if tok.upper == "ON" && (isDDLIndexOn(clauseToks[:idx+1]) ||
			(idx+1 < len(clauseToks) && !clauseToks[idx+1].afterDot &&
				(clauseToks[idx+1].upper == "CONFLICT" || clauseToks[idx+1].upper == "CONSTRAINT"))) {
			break
		}
		if tok.upper == "CONFLICT" || tok.upper == "DO" {
			break
		}
		switch tok.upper {
		case "LET", "REMOVE", "WHERE", "HAVING", "SET", "ORDER", "GROUP", "BY", "AND", "OR", "ON", "WHEN", "THEN", "ELSE", "RETURNING":
			return false, true, false
		}
	}
	return false, false, false
}

func isDDLIndexPathContext(toks []clauseTok) bool {
	return len(toks) == 4 && isDDLIndexOn(toks[:3])
}

func isDDLIndexOn(toks []clauseTok) bool {
	if len(toks) == 3 && !toks[0].afterDot && !toks[1].afterDot && !toks[2].afterDot &&
		toks[0].upper == "CREATE" && toks[1].upper == "INDEX" && toks[2].upper == "ON" {
		return true
	}
	if len(toks) == 4 && !toks[0].afterDot && !toks[1].afterDot && !toks[3].afterDot &&
		toks[0].upper == "DROP" && toks[1].upper == "INDEX" && toks[3].upper == "ON" {
		return true
	}
	return false
}

func matchesPrefix(text, prefix string, quoted bool) bool {
	if prefix == "" {
		return true
	}
	if quoted {
		return strings.HasPrefix(text, prefix)
	}
	return strings.HasPrefix(strings.ToUpper(text), strings.ToUpper(prefix))
}

func unclosedDoubleQuoteStart(s string, end int) (int, bool) {
	if end <= 0 {
		return 0, false
	}
	if end > len(s) {
		end = len(s)
	}
	i := 0
	for i < end {
		switch {
		case s[i] == '\'':
			i = skipQuotedLiteral(s[:end], i, '\'') + 1
		case s[i] == '"':
			next, closed := skipQuotedLiteralWithState(s[:end], i, '"')
			if !closed {
				return i, true
			}
			i = next + 1
		case s[i] == '-' && i+1 < end && s[i+1] == '-':
			i = skipLineComment(s[:end], i) + 1
		case s[i] == '/' && i+1 < end && s[i+1] == '*':
			next, _ := skipBlockComment(s[:end], i)
			i = next + 1
		default:
			i++
		}
	}
	return 0, false
}

type tableRef struct {
	name        string
	alias       string
	quoted      bool
	aliasQuoted bool
	subquery    []tableToken
}

// resolveColumnNames extracts column names from tables visible in the active
// query scope around pos, searching enclosing outer scopes when a qualified
// reference targets an outer alias or unqualified correlated subquery columns.
func resolveColumnNames(statement string, pos int, meta *metadata.DatabaseSchemaMetadata) []string {
	if meta == nil {
		return nil
	}
	qualifier, qualQuoted := columnQualifierBeforeCursor(statement, pos)
	scopeStatement := statement
	scopePos := pos
	if qStart, ok := unclosedDoubleQuoteStart(statement, pos); ok {
		scopeStatement = statement[:qStart] + statement[pos:]
		scopePos = qStart
	}
	scopes := extractQueryScopes(scopeStatement, scopePos)

	allTables := collectMetadataTables(meta)
	if qualifier != "" {
		for _, scopeText := range scopes {
			refs := extractTableRefsFromScope(scopeText)
			var cols []string
			matched := false
			for _, ref := range refs {
				if matchesTableQualifier(ref, qualifier, qualQuoted) {
					matched = true
					cols = append(cols, resolveTableRefColumns(ref, allTables)...)
				}
			}
			if matched {
				return cols
			}
		}
		return nil
	}

	hasAnyRef := false
	var cols []string
	seenCols := make(map[string]struct{})
	for _, scopeText := range scopes {
		refs := extractTableRefsFromScope(scopeText)
		if len(refs) == 0 {
			continue
		}
		hasAnyRef = true
		for _, ref := range refs {
			for _, col := range resolveTableRefColumns(ref, allTables) {
				if _, ok := seenCols[col]; ok {
					continue
				}
				seenCols[col] = struct{}{}
				cols = append(cols, col)
			}
		}
	}
	if !hasAnyRef {
		for _, table := range allTables {
			for _, col := range table.GetColumns() {
				if name := col.GetName(); name != "" {
					if _, ok := seenCols[name]; !ok {
						seenCols[name] = struct{}{}
						cols = append(cols, name)
					}
				}
			}
		}
	}
	return cols
}

func matchesTableQualifier(ref tableRef, qualifier string, qualQuoted bool) bool {
	if ref.alias != "" {
		return partiqlNameMatches(ref.alias, qualifier, ref.aliasQuoted || qualQuoted)
	}
	return ref.name != "" && partiqlNameMatches(ref.name, qualifier, ref.quoted || qualQuoted)
}

func partiqlNameMatches(candidate, query string, quoted bool) bool {
	if quoted {
		return candidate == query
	}
	return candidate == query || strings.EqualFold(candidate, query)
}

func resolveTableRefColumns(ref tableRef, allTables []*metadata.TableMetadata) []string {
	if len(ref.subquery) > 0 {
		return resolveSubqueryColumns(ref.subquery, allTables)
	}
	table := findMetadataTable(allTables, ref.name, ref.quoted)
	if table == nil {
		return nil
	}
	var cols []string
	for _, col := range table.GetColumns() {
		if col.GetName() != "" {
			cols = append(cols, col.GetName())
		}
	}
	return cols
}

func resolveSubqueryColumns(tokens []tableToken, allTables []*metadata.TableMetadata) []string {
	for len(tokens) >= 2 && tokens[0].paren == '(' && findMatchingTokenParen(tokens, 0) == len(tokens)-1 {
		tokens = tokens[1 : len(tokens)-1]
	}
	selIdx := findMainSelectIndex(tokens)
	if selIdx < 0 {
		return nil
	}
	listStart := selIdx + 1
	listEnd := len(tokens)
	nestDepth := 0
	for i := listStart; i < len(tokens); i++ {
		if delta := tokenNestDelta(tokens[i]); delta != 0 {
			if nestDepth+delta >= 0 {
				nestDepth += delta
			} else {
				nestDepth = 0
			}
			continue
		}
		if nestDepth == 0 && !tokens[i].quoted && isSelectListEndKeyword(strings.ToUpper(tokens[i].text)) {
			listEnd = i
			break
		}
	}

	var innerRefs []tableRef
	getInnerRefs := func() []tableRef {
		if innerRefs == nil {
			innerRefs = extractTableRefsFromTokens(tokens)
		}
		return innerRefs
	}

	var cols []string
	itemStart := listStart
	for itemStart < listEnd && !tokens[itemStart].quoted {
		firstUpper := strings.ToUpper(tokens[itemStart].text)
		if firstUpper == "VALUE" {
			// SELECT VALUE expr returns scalar/element values rather than named tuple attributes.
			return nil
		}
		if firstUpper == "DISTINCT" || firstUpper == "ALL" {
			itemStart++
			continue
		}
		break
	}
	projEnd := listEnd
	var excludeRules []excludeRule
	exclDepth := 0
	for i := itemStart; i < listEnd; i++ {
		if delta := tokenNestDelta(tokens[i]); delta != 0 {
			if exclDepth+delta >= 0 {
				exclDepth += delta
			} else {
				exclDepth = 0
			}
			continue
		}
		if exclDepth == 0 && !tokens[i].quoted && strings.ToUpper(tokens[i].text) == "EXCLUDE" {
			projEnd = i
			excludeRules = parseExcludeClause(tokens[i+1 : listEnd])
			break
		}
	}

	nestDepth = 0
	for i := itemStart; i <= projEnd; i++ {
		if i < projEnd {
			if delta := tokenNestDelta(tokens[i]); delta != 0 {
				if nestDepth+delta >= 0 {
					nestDepth += delta
				} else {
					nestDepth = 0
				}
				continue
			}
		}
		if i < projEnd && (nestDepth != 0 || !tokens[i].comma) {
			continue
		}
		item := tokens[itemStart:i]
		itemStart = i + 1
		if len(item) == 0 {
			continue
		}
		if len(item) == 1 && !item[0].quoted && item[0].text == "*" {
			for _, ref := range getInnerRefs() {
				cols = appendFilteredExpandedColumns(cols, ref, resolveTableRefColumns(ref, allTables), excludeRules)
			}
			continue
		}
		if len(item) == 1 && !item[0].quoted && strings.HasSuffix(item[0].text, ".*") {
			qual := strings.TrimSuffix(item[0].text, ".*")
			for _, ref := range getInnerRefs() {
				if matchesTableQualifier(ref, qual, false) {
					cols = appendFilteredExpandedColumns(cols, ref, resolveTableRefColumns(ref, allTables), excludeRules)
				}
			}
			continue
		}
		if len(item) == 2 && item[0].quoted && !item[1].quoted && item[1].text == ".*" {
			qual := item[0].text
			for _, ref := range getInnerRefs() {
				if matchesTableQualifier(ref, qual, true) {
					cols = appendFilteredExpandedColumns(cols, ref, resolveTableRefColumns(ref, allTables), excludeRules)
				}
			}
			continue
		}
		if col := selectItemOutputName(item); col != "" {
			cols = append(cols, col)
		}
	}
	return cols
}

type excludeRule struct {
	root       string
	rootQuoted bool
	col        string
	colQuoted  bool
	allCols    bool
}

func parseExcludeClause(tokens []tableToken) []excludeRule {
	var rules []excludeRule
	nestDepth := 0
	start := 0
	for i := 0; i <= len(tokens); i++ {
		if i < len(tokens) {
			if delta := tokenNestDelta(tokens[i]); delta != 0 {
				if nestDepth+delta >= 0 {
					nestDepth += delta
				} else {
					nestDepth = 0
				}
				continue
			}
			if nestDepth != 0 || !tokens[i].comma {
				continue
			}
		}
		if rule, ok := parseExcludeRule(tokens[start:i]); ok {
			rules = append(rules, rule)
		}
		start = i + 1
	}
	return rules
}

func parseExcludeRule(tokens []tableToken) (excludeRule, bool) {
	type pathSeg struct {
		name   string
		quoted bool
	}
	var segs []pathSeg
	needSeg := true
	for _, tok := range tokens {
		if tok.paren != 0 || tok.delim != 0 || tok.comma {
			return excludeRule{}, false
		}
		if tok.quoted {
			if !needSeg {
				return excludeRule{}, false
			}
			segs = append(segs, pathSeg{name: tok.text, quoted: true})
			needSeg = false
			continue
		}
		parts := strings.Split(tok.text, ".")
		for idx, p := range parts {
			if idx > 0 {
				needSeg = true
			}
			if p == "" {
				continue
			}
			if !needSeg {
				return excludeRule{}, false
			}
			segs = append(segs, pathSeg{name: p, quoted: false})
			needSeg = false
		}
	}
	if len(segs) != 2 {
		return excludeRule{}, false
	}
	rule := excludeRule{
		root:       segs[0].name,
		rootQuoted: segs[0].quoted,
	}
	if !segs[1].quoted && segs[1].name == "*" {
		rule.allCols = true
		return rule, true
	}
	rule.col = segs[1].name
	rule.colQuoted = segs[1].quoted
	return rule, true
}

func appendFilteredExpandedColumns(dst []string, ref tableRef, cols []string, rules []excludeRule) []string {
	for _, col := range cols {
		if isColumnExcluded(ref, col, rules) {
			continue
		}
		dst = append(dst, col)
	}
	return dst
}

func isColumnExcluded(ref tableRef, col string, rules []excludeRule) bool {
	for _, rule := range rules {
		if !matchesTableQualifier(ref, rule.root, rule.rootQuoted) {
			continue
		}
		if rule.allCols {
			return true
		}
		if rule.colQuoted {
			if col == rule.col {
				return true
			}
		} else if strings.EqualFold(col, rule.col) {
			return true
		}
	}
	return false
}

func findMainSelectIndex(tokens []tableToken) int {
	depth := 0
	for i, tok := range tokens {
		switch tok.paren {
		case '(':
			depth++
			continue
		case ')':
			if depth > 0 {
				depth--
			}
			continue
		}
		if depth == 0 && !tok.quoted && strings.ToUpper(tok.text) == "SELECT" {
			return i
		}
	}
	return -1
}

func isSelectListEndKeyword(upper string) bool {
	switch upper {
	case "FROM", "WHERE", "GROUP", "HAVING", "ORDER", "LIMIT", "OFFSET", "UNION", "INTERSECT", "EXCEPT":
		return true
	}
	return false
}

func selectItemOutputName(item []tableToken) string {
	if len(item) >= 2 && !item[len(item)-2].quoted && strings.ToUpper(item[len(item)-2].text) == "AS" {
		last := item[len(item)-1]
		if isValidProjectedAlias(last) {
			return last.text
		}
		return ""
	}
	if len(item) >= 2 {
		prev := item[len(item)-2]
		last := item[len(item)-1]
		if canPrecedeImplicitAlias(prev) && isValidProjectedAlias(last) {
			return last.text
		}
		return ""
	}
	if len(item) == 1 && item[0].paren == 0 && !item[0].comma && item[0].text != "" {
		name := item[0].text
		if item[0].quoted {
			return name
		}
		if dot := strings.LastIndexByte(name, '.'); dot >= 0 && dot+1 < len(name) {
			prefix := name[:dot]
			name = name[dot+1:]
			if !isPlainIdentifier(prefix) {
				return ""
			}
		}
		if !isPlainIdentifier(name) {
			return ""
		}
		if strings.EqualFold(name, "NULL") || strings.EqualFold(name, "MISSING") ||
			strings.EqualFold(name, "TRUE") || strings.EqualFold(name, "FALSE") {
			return ""
		}
		return name
	}
	return ""
}

func isValidProjectedAlias(tok tableToken) bool {
	if tok.paren != 0 || tok.comma || tok.text == "" {
		return false
	}
	if tok.quoted {
		return true
	}
	if !isPlainIdentifier(tok.text) || isKeyword(strings.ToUpper(tok.text)) {
		return false
	}
	// Reject pure numeric literals such as the "1" in "c1 + 1".
	allDigits := true
	for i := 0; i < len(tok.text); i++ {
		if tok.text[i] < '0' || tok.text[i] > '9' {
			allDigits = false
			break
		}
	}
	return !allDigits
}

func canPrecedeImplicitAlias(prev tableToken) bool {
	if prev.paren == ')' || prev.quoted {
		return true
	}
	if prev.paren != 0 || prev.comma || prev.text == "" {
		return false
	}
	upper := strings.ToUpper(prev.text)
	switch upper {
	case "AND", "OR", "NOT", "IS", "IN", "LIKE", "BETWEEN", "CASE", "WHEN",
		"THEN", "ELSE", "CAST", "AS", "EXISTS", "DISTINCT", "ALL", "VALUE", "FROM":
		return false
	}
	// Must be a plain identifier or dotted/numeric/literal operand, not an operator like +, -, *, /, ||, =.
	return isOperandToken(prev.text)
}

func isPlainIdentifier(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if !isIdentChar(s[i]) {
			return false
		}
	}
	return true
}

func isOperandToken(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if !isIdentChar(s[i]) && s[i] != '.' {
			return false
		}
	}
	return true
}

func columnQualifierBeforeCursor(statement string, pos int) (string, bool) {
	_, beforeEnd, _ := extractPrefixAndQuote(statement, pos)
	before := strings.TrimRight(statement[:beforeEnd], " \t\n\r")
	if !strings.HasSuffix(before, ".") {
		return "", false
	}
	before = strings.TrimRight(strings.TrimSuffix(before, "."), " \t\n\r")
	if before == "" {
		return "", false
	}
	if before[len(before)-1] == '"' {
		for i := len(before) - 2; i >= 0; i-- {
			if before[i] == '"' {
				if i > 0 && before[i-1] == '"' {
					i--
					continue
				}
				return strings.ReplaceAll(before[i+1:len(before)-1], `""`, `"`), true
			}
		}
		return "", false
	}
	start := len(before)
	for start > 0 && isIdentChar(before[start-1]) {
		start--
	}
	return before[start:], false
}

func collectMetadataTables(meta *metadata.DatabaseSchemaMetadata) []*metadata.TableMetadata {
	var tables []*metadata.TableMetadata
	for _, schema := range meta.GetSchemas() {
		tables = append(tables, schema.GetTables()...)
	}
	return tables
}

func findMetadataTable(tables []*metadata.TableMetadata, name string, quoted bool) *metadata.TableMetadata {
	for _, t := range tables {
		if t.GetName() == name {
			return t
		}
	}
	if !quoted {
		for _, t := range tables {
			if strings.EqualFold(t.GetName(), name) {
				return t
			}
		}
	}
	return nil
}

// extractCurrentStatement returns the statement segment containing pos and pos
// relative to that segment. Segment end offsets are exclusive at statement
// boundaries so a cursor immediately after ';' begins a new empty statement.
func extractCurrentStatement(statement string, pos int) (string, int) {
	segs := partiql.Split(statement)
	for i, seg := range segs {
		if pos < seg.ByteStart {
			continue
		}
		if pos < seg.ByteEnd {
			return seg.Text, pos - seg.ByteStart
		}
		if pos == seg.ByteEnd {
			if strings.HasSuffix(strings.TrimRight(seg.Text, " \t\r\n"), ";") {
				if i+1 < len(segs) {
					continue
				}
				return "", 0
			}
			return seg.Text, pos - seg.ByteStart
		}
	}
	return statement, pos
}

type scopeSpan struct {
	start int
	end   int
}

// extractQueryScopes returns query block slices enclosing pos, ordered from
// innermost to outermost, isolating the active UNION/INTERSECT/EXCEPT arm at pos.
// Enclosing scopes collapse their active inner subquery span so each byte of
// statement is scanned at most once across all levels.
func extractQueryScopes(statement string, pos int) []string {
	spans := findActiveSubquerySpans(statement, pos)
	spans = append(spans, scopeSpan{start: 0, end: len(statement)})

	scopes := make([]string, 0, len(spans))
	var prevChild *scopeSpan
	for idx := range spans {
		sp := spans[idx]
		if prevChild == nil {
			scopes = append(scopes, isolateSetOperationArm(statement[sp.start:sp.end], pos-sp.start))
		} else {
			raw := statement[sp.start:prevChild.start] + statement[prevChild.end:sp.end]
			scopes = append(scopes, isolateSetOperationArm(raw, prevChild.start-sp.start))
		}
		prevChild = &spans[idx]
	}
	return scopes
}

func isolateSetOperationArm(scope string, offset int) string {
	if offset < 0 {
		offset = 0
	}
	if offset > len(scope) {
		offset = len(scope)
	}
	armStart := 0
	depth := 0
	i := 0
	for i < len(scope) {
		c := scope[i]
		switch {
		case c == '\'':
			i = skipQuotedLiteral(scope, i, '\'') + 1
		case c == '"':
			i = skipQuotedLiteral(scope, i, '"') + 1
		case c == '-' && i+1 < len(scope) && scope[i+1] == '-':
			i = skipLineComment(scope, i) + 1
		case c == '/' && i+1 < len(scope) && scope[i+1] == '*':
			end, _ := skipBlockComment(scope, i)
			i = end + 1
		case c == '(':
			depth++
			i++
		case c == ')':
			if depth > 0 {
				depth--
			}
			i++
		case isIdentChar(c):
			start := i
			for i < len(scope) && isIdentChar(scope[i]) {
				i++
			}
			if depth == 0 && (start == 0 || scope[start-1] != '.') {
				upper := strings.ToUpper(scope[start:i])
				if upper == "UNION" || upper == "INTERSECT" || upper == "EXCEPT" {
					if offset < start {
						return scope[armStart:start]
					}
					armStart = i
				}
			}
		default:
			i++
		}
	}
	return scope[armStart:]
}

func findActiveQueryScopeStart(statement string, pos int) int {
	spans := findActiveSubquerySpans(statement, pos)
	if len(spans) == 0 {
		return 0
	}
	return spans[0].start
}

// findActiveSubquerySpans computes all SELECT/WITH subquery spans enclosing pos
// in a single linear pass, ordered innermost to outermost.
func findActiveSubquerySpans(statement string, pos int) []scopeSpan {
	if pos > len(statement) {
		pos = len(statement)
	}
	if pos < 0 {
		pos = 0
	}
	type parenFrame struct {
		start int
		sub   bool
	}
	var stack []parenFrame
	var pendingOpen []int
	var closedSpans []scopeSpan
	openSubCount := 0

	i := 0
	for i < len(statement) {
		c := statement[i]
		switch {
		case c == '\'':
			pendingOpen = pendingOpen[:0]
			i = skipQuotedLiteral(statement, i, '\'') + 1
		case c == '"':
			pendingOpen = pendingOpen[:0]
			i = skipQuotedLiteral(statement, i, '"') + 1
		case c == '-' && i+1 < len(statement) && statement[i+1] == '-':
			i = skipLineComment(statement, i) + 1
		case c == '/' && i+1 < len(statement) && statement[i+1] == '*':
			end, closed := skipBlockComment(statement, i)
			if !closed {
				i = len(statement)
				break
			}
			i = end + 1
		case c == '(':
			stack = append(stack, parenFrame{start: i + 1})
			if i < pos {
				pendingOpen = append(pendingOpen, len(stack)-1)
			} else {
				pendingOpen = pendingOpen[:0]
			}
			i++
		case c == ')':
			pendingOpen = pendingOpen[:0]
			if len(stack) > 0 {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if top.sub {
					openSubCount--
					if top.start <= pos && i >= pos {
						closedSpans = append(closedSpans, scopeSpan{start: top.start, end: i})
					}
				}
			}
			if i >= pos && openSubCount == 0 {
				return closedSpans
			}
			i++
		case c == ' ' || c == '\t' || c == '\n' || c == '\r':
			i++
		case isIdentChar(c):
			start := i
			for i < len(statement) && isIdentChar(statement[i]) {
				i++
			}
			if len(pendingOpen) > 0 {
				word := strings.ToUpper(statement[start:i])
				if word == "SELECT" || word == "WITH" || word == "PIVOT" {
					for _, idx := range pendingOpen {
						stack[idx].sub = true
						openSubCount++
					}
				}
				pendingOpen = pendingOpen[:0]
			}
		default:
			pendingOpen = pendingOpen[:0]
			i++
		}
	}

	spans := make([]scopeSpan, 0, len(closedSpans)+openSubCount)
	spans = append(spans, closedSpans...)
	for k := len(stack) - 1; k >= 0; k-- {
		if stack[k].sub && stack[k].start <= pos {
			spans = append(spans, scopeSpan{start: stack[k].start, end: len(statement)})
		}
	}
	return spans
}

func skipLineComment(s string, start int) int {
	i := start + 2
	for i < len(s) && s[i] != '\n' && s[i] != '\r' {
		i++
	}
	if i >= len(s) {
		return len(s) - 1
	}
	return i
}

func skipBlockComment(s string, start int) (int, bool) {
	i := start + 2
	depth := 1
	for i+1 < len(s) && depth > 0 {
		if s[i] == '/' && s[i+1] == '*' {
			depth++
			i += 2
		} else if s[i] == '*' && s[i+1] == '/' {
			depth--
			i += 2
		} else {
			i++
		}
	}
	if depth > 0 {
		return len(s) - 1, false
	}
	return i - 1, true
}

func skipQuotedLiteral(s string, start int, quote byte) int {
	end, _ := skipQuotedLiteralWithState(s, start, quote)
	return end
}

func skipQuotedLiteralWithState(s string, start int, quote byte) (int, bool) {
	i := start + 1
	for i < len(s) {
		if s[i] == quote {
			if i+1 < len(s) && s[i+1] == quote {
				i += 2
				continue
			}
			return i, true
		}
		i++
	}
	return len(s) - 1, false
}

// extractTableRefsFromScope scans the active query scope for table names and
// optional aliases after FROM, JOIN, INTO, or UPDATE keywords, unwrapping
// parenthesized relation expressions and recording subqueries in FROM.
func extractTableRefsFromScope(statement string) []tableRef {
	cleaned := maskStringsAndComments(statement)
	return extractTableRefsFromTokens(tokenizeForTableExtraction(cleaned))
}

func extractTableRefsFromTokens(tokens []tableToken) []tableRef {
	var refs []tableRef
	i := 0
	for i < len(tokens) {
		if tokens[i].paren == '(' {
			i = findMatchingTokenParen(tokens, i) + 1
			continue
		}
		if tokens[i].paren == ')' {
			i++
			continue
		}
		upper := strings.ToUpper(tokens[i].text)
		isDDLIndexTarget := upper == "ON" && !tokens[i].quoted && i == 2 &&
			!tokens[0].quoted && !tokens[1].quoted &&
			strings.ToUpper(tokens[0].text) == "CREATE" &&
			strings.ToUpper(tokens[1].text) == "INDEX"
		if ((upper == "FROM" || upper == "JOIN" || upper == "INTO" || upper == "UPDATE") && !tokens[i].quoted || isDDLIndexTarget) && i+1 < len(tokens) {
			i++
			for i < len(tokens) {
				parsed, next, ok := parseTableRefAt(tokens, i)
				if !ok {
					break
				}
				refs = append(refs, parsed...)
				i = next
				if i < len(tokens) && tokens[i].comma {
					i++
					continue
				}
				break
			}
		} else {
			i++
		}
	}
	return refs
}

func parseTableRefAt(tokens []tableToken, i int) ([]tableRef, int, bool) {
	if i >= len(tokens) {
		return nil, i, false
	}
	if tokens[i].paren == '(' {
		closeIdx := findMatchingTokenParen(tokens, i)
		innerEnd := closeIdx
		if innerEnd > len(tokens) {
			innerEnd = len(tokens)
		}
		inner := tokens[i+1 : innerEnd]
		nextIdx := innerEnd
		if nextIdx < len(tokens) && tokens[nextIdx].paren == ')' {
			nextIdx++
		}
		alias, aliasQuoted, nextIdx := parseOptionalAlias(tokens, nextIdx)
		if nextIsSubqueryStart(inner, 0) {
			return []tableRef{{alias: alias, aliasQuoted: aliasQuoted, subquery: inner}}, nextIdx, true
		}
		innerRefs := extractRefsFromParenthesizedExpr(inner)
		if len(innerRefs) == 1 && alias != "" && innerRefs[0].alias == "" {
			innerRefs[0].alias = alias
			innerRefs[0].aliasQuoted = aliasQuoted
		}
		return innerRefs, nextIdx, true
	}
	if tokens[i].paren == ')' || tokens[i].comma {
		return nil, i, false
	}
	tok := tokens[i]
	if !tok.quoted && isKeyword(strings.ToUpper(tok.text)) {
		return nil, i, false
	}
	alias, aliasQuoted, nextIdx := parseOptionalAlias(tokens, i+1)
	return []tableRef{{name: tok.text, alias: alias, quoted: tok.quoted, aliasQuoted: aliasQuoted}}, nextIdx, true
}

func parseOptionalAlias(tokens []tableToken, i int) (string, bool, int) {
	if i >= len(tokens) || tokens[i].comma || tokens[i].paren != 0 {
		return "", false, i
	}
	if !tokens[i].quoted && strings.ToUpper(tokens[i].text) == "AS" {
		if i+1 < len(tokens) && !tokens[i+1].comma && tokens[i+1].paren == 0 {
			return tokens[i+1].text, tokens[i+1].quoted, i + 2
		}
		return "", false, i + 1
	}
	if tokens[i].quoted || !isKeyword(strings.ToUpper(tokens[i].text)) {
		return tokens[i].text, tokens[i].quoted, i + 1
	}
	return "", false, i
}

func extractRefsFromParenthesizedExpr(tokens []tableToken) []tableRef {
	var refs []tableRef
	parsed, idx, ok := parseTableRefAt(tokens, 0)
	if ok {
		refs = append(refs, parsed...)
	} else {
		idx = 0
	}
	for idx < len(tokens) {
		if tokens[idx].paren == '(' {
			idx = findMatchingTokenParen(tokens, idx) + 1
			continue
		}
		upper := strings.ToUpper(tokens[idx].text)
		if tokens[idx].comma || (!tokens[idx].quoted && (upper == "JOIN" || upper == "FROM" || upper == "INTO" || upper == "UPDATE")) {
			idx++
			for idx < len(tokens) {
				p, next, matched := parseTableRefAt(tokens, idx)
				if !matched {
					break
				}
				refs = append(refs, p...)
				idx = next
				if idx < len(tokens) && tokens[idx].comma {
					idx++
					continue
				}
				break
			}
			continue
		}
		idx++
	}
	return refs
}

func findMatchingTokenParen(tokens []tableToken, openIdx int) int {
	depth := 1
	for i := openIdx + 1; i < len(tokens); i++ {
		switch tokens[i].paren {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return len(tokens)
}

func nextIsSubqueryStart(tokens []tableToken, idx int) bool {
	for idx < len(tokens) && tokens[idx].paren == '(' {
		idx++
	}
	if idx >= len(tokens) || tokens[idx].quoted {
		return false
	}
	upper := strings.ToUpper(tokens[idx].text)
	return upper == "SELECT" || upper == "WITH"
}

type tableToken struct {
	text   string
	quoted bool
	comma  bool
	paren  byte
	delim  byte
}

func tokenNestDelta(tok tableToken) int {
	switch tok.paren {
	case '(':
		return 1
	case ')':
		return -1
	}
	switch tok.delim {
	case '[', '{', '<':
		return 1
	case ']', '}', '>':
		return -1
	}
	return 0
}

func tokenizeForTableExtraction(s string) []tableToken {
	var tokens []tableToken
	i := 0
	for i < len(s) {
		if s[i] == ' ' || s[i] == '\t' || s[i] == '\n' || s[i] == '\r' {
			i++
			continue
		}
		if s[i] == ',' {
			tokens = append(tokens, tableToken{comma: true})
			i++
			continue
		}
		if s[i] == '(' || s[i] == ')' {
			tokens = append(tokens, tableToken{paren: s[i]})
			i++
			continue
		}
		if s[i] == '[' || s[i] == ']' || s[i] == '{' || s[i] == '}' {
			tokens = append(tokens, tableToken{delim: s[i]})
			i++
			continue
		}
		if (s[i] == '<' || s[i] == '>') && i+1 < len(s) && s[i+1] == s[i] {
			tokens = append(tokens, tableToken{delim: s[i]})
			i += 2
			continue
		}
		if s[i] == '"' {
			i++
			start := i
			for i < len(s) {
				if s[i] == '"' {
					if i+1 < len(s) && s[i+1] == '"' {
						i += 2
						continue
					}
					break
				}
				i++
			}
			text := strings.ReplaceAll(s[start:i], `""`, `"`)
			if i < len(s) {
				i++
			}
			tokens = append(tokens, tableToken{text: text, quoted: true})
			continue
		}
		start := i
		for i < len(s) && s[i] != ' ' && s[i] != '\t' && s[i] != '\n' &&
			s[i] != '\r' && s[i] != ',' && s[i] != '"' && s[i] != ';' &&
			s[i] != '(' && s[i] != ')' && s[i] != '[' && s[i] != ']' &&
			s[i] != '{' && s[i] != '}' &&
			((s[i] != '<' && s[i] != '>') || i+1 >= len(s) || s[i+1] != s[i]) {
			i++
		}
		if i > start {
			tokens = append(tokens, tableToken{text: s[start:i]})
		} else {
			i++
		}
	}
	return tokens
}

func maskStringsAndComments(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	i := 0
	for i < len(s) {
		switch {
		case s[i] == '\'':
			b.WriteByte('_')
			i++
			for i < len(s) {
				if s[i] == '\'' {
					i++
					if i < len(s) && s[i] == '\'' {
						b.WriteByte('_')
						i++
						continue
					}
					break
				}
				i++
			}
		case s[i] == '-' && i+1 < len(s) && s[i+1] == '-':
			b.WriteByte(' ')
			i += 2
			for i < len(s) && s[i] != '\n' && s[i] != '\r' {
				i++
			}
		case s[i] == '/' && i+1 < len(s) && s[i+1] == '*':
			b.WriteByte(' ')
			i += 2
			depth := 1
			for i+1 < len(s) && depth > 0 {
				if s[i] == '/' && s[i+1] == '*' {
					depth++
					i += 2
				} else if s[i] == '*' && s[i+1] == '/' {
					depth--
					i += 2
				} else {
					i++
				}
			}
			if depth > 0 {
				i = len(s)
			}
		default:
			b.WriteByte(s[i])
			i++
		}
	}
	return b.String()
}

func isKeyword(w string) bool {
	switch w {
	case "SELECT", "FROM", "WHERE", "AND", "OR", "NOT", "INSERT", "INTO",
		"VALUE", "VALUES", "UPDATE", "SET", "DELETE", "CREATE", "DROP",
		"TABLE", "INDEX", "AS", "AT", "BY", "ON", "IN", "IS", "LIKE",
		"BETWEEN", "ORDER", "GROUP", "HAVING", "LIMIT", "OFFSET",
		"JOIN", "LEFT", "RIGHT", "INNER", "OUTER", "CROSS",
		"UNION", "INTERSECT", "EXCEPT", "DISTINCT", "ALL",
		"NULL", "MISSING", "TRUE", "FALSE", "CAST", "CASE", "WHEN",
		"THEN", "ELSE", "END", "EXISTS", "ASC", "DESC", "EXPLAIN", "RETURNING",
		"PIVOT", "LET", "REMOVE", "CONFLICT", "DO", "NOTHING", "REPLACE",
		"EXCLUDE", "EXCLUDED":
		return true
	}
	return false
}

// keywords is the list of PartiQL keywords suggested for auto-completion.
// This is a subset of the full keyword list — only commonly-used keywords.
var keywords = []string{
	"SELECT", "FROM", "WHERE", "AND", "OR", "NOT",
	"INSERT", "INTO", "VALUE", "VALUES", "UPDATE", "SET",
	"DELETE", "CREATE", "DROP", "TABLE", "INDEX",
	"AS", "AT", "BY", "ON", "IN", "IS", "LIKE", "BETWEEN",
	"ORDER", "GROUP", "HAVING", "LIMIT", "OFFSET",
	"JOIN", "LEFT", "RIGHT", "INNER", "OUTER", "CROSS",
	"UNION", "INTERSECT", "EXCEPT", "DISTINCT", "ALL",
	"NULL", "MISSING", "TRUE", "FALSE",
	"CAST", "CASE", "WHEN", "THEN", "ELSE", "END",
	"EXISTS", "ASC", "DESC", "REMOVE", "EXCLUDED",
}
