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

	stmtSQL, localPos := extractCurrentStatement(input, pos)
	prefix := extractPrefix(stmtSQL, localPos)
	inFromContext, inColumnContext := analyzeCursorContext(stmtSQL, localPos, len(prefix))

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
			if matchesPrefix(table, prefix) {
				addCandidate(table, "table")
			}
		}
	}

	if meta != nil && inColumnContext {
		for _, col := range resolveColumnNames(stmtSQL, localPos, meta) {
			if matchesPrefix(col, prefix) {
				addCandidate(col, "column")
			}
		}
	}

	for _, kw := range keywords {
		if matchesPrefix(kw, prefix) {
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

// analyzeCursorContext determines whether the cursor is in a table-introducing
// clause (FROM / JOIN / INTO) or an expression/projection position where columns
// are valid. Expression-nested FROM keywords (such as TRIM(' ' FROM ...) or
// SUBSTRING(x FROM ...)) stay inside expression parentheses and do not suppress
// column completion or trigger table suggestions.
func analyzeCursorContext(statement string, pos, prefixLen int) (inFrom, inColumn bool) {
	if pos > len(statement) {
		pos = len(statement)
	}
	if pos < 0 {
		pos = 0
	}
	beforeEnd := pos - prefixLen
	if beforeEnd < 0 {
		beforeEnd = 0
	}
	activeStart := findActiveQueryScopeStart(statement, beforeEnd)
	slice := isolateSetOperationArm(statement[activeStart:beforeEnd], beforeEnd-activeStart)

	type clauseTok struct {
		upper    string
		afterDot bool
	}
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
			i = skipLineComment(slice, i) + 1
		case c == '/' && i+1 < len(slice) && slice[i+1] == '*':
			end, _ := skipBlockComment(slice, i)
			i = end + 1
		case c == '\'' || c == '"':
			end := skipQuotedLiteral(slice, i, c)
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

	if !prevDot && exprDepth == 0 && len(clauseToks) > 0 {
		last := clauseToks[len(clauseToks)-1]
		if !last.afterDot && (last.upper == "FROM" || last.upper == "JOIN" || last.upper == "INTO") {
			return true, false
		}
	}

	selectIdx := -1
	fromIdx := -1
	for idx, tok := range clauseToks {
		if tok.afterDot {
			continue
		}
		switch tok.upper {
		case "SELECT":
			selectIdx = idx
		case "FROM":
			fromIdx = idx
		}
	}
	if selectIdx >= 0 && (fromIdx < 0 || fromIdx < selectIdx) {
		return false, true
	}

	for idx := len(clauseToks) - 1; idx >= 0; idx-- {
		if idx <= fromIdx {
			break
		}
		tok := clauseToks[idx]
		if tok.afterDot {
			continue
		}
		switch tok.upper {
		case "WHERE", "HAVING", "SET", "ORDER", "GROUP", "BY", "AND", "OR", "ON", "WHEN", "THEN", "ELSE":
			return false, true
		}
	}
	return false, false
}

func matchesPrefix(text, prefix string) bool {
	if prefix == "" {
		return true
	}
	return strings.HasPrefix(strings.ToUpper(text), strings.ToUpper(prefix))
}

type tableRef struct {
	name     string
	alias    string
	subquery []tableToken
}

// resolveColumnNames extracts column names from tables visible in the active
// query scope around pos, searching enclosing outer scopes when a qualified
// reference targets an outer alias in a correlated subquery.
func resolveColumnNames(statement string, pos int, meta *metadata.DatabaseSchemaMetadata) []string {
	if meta == nil {
		return nil
	}
	scopes := extractQueryScopes(statement, pos)
	qualifier := columnQualifierBeforeCursor(statement, pos)

	allTables := collectMetadataTables(meta)
	if qualifier != "" {
		for _, scopeText := range scopes {
			refs := extractTableRefsFromScope(scopeText)
			var cols []string
			matched := false
			for _, ref := range refs {
				if matchesTableQualifier(ref, qualifier) {
					matched = true
					cols = append(cols, resolveTableRefColumns(ref, allTables, 0)...)
				}
			}
			if matched {
				return cols
			}
		}
		return nil
	}

	var refs []tableRef
	if len(scopes) > 0 {
		refs = extractTableRefsFromScope(scopes[0])
	}
	if len(refs) == 0 {
		var cols []string
		for _, table := range allTables {
			for _, col := range table.GetColumns() {
				if col.GetName() != "" {
					cols = append(cols, col.GetName())
				}
			}
		}
		return cols
	}

	var cols []string
	for _, ref := range refs {
		cols = append(cols, resolveTableRefColumns(ref, allTables, 0)...)
	}
	return cols
}

func matchesTableQualifier(ref tableRef, qualifier string) bool {
	if ref.alias != "" && strings.EqualFold(ref.alias, qualifier) {
		return true
	}
	return ref.name != "" && strings.EqualFold(ref.name, qualifier)
}

func resolveTableRefColumns(ref tableRef, allTables []*metadata.TableMetadata, depth int) []string {
	if len(ref.subquery) > 0 {
		if depth >= 4 {
			return nil
		}
		return resolveSubqueryColumns(ref.subquery, allTables, depth+1)
	}
	table := findMetadataTable(allTables, ref.name)
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

func resolveSubqueryColumns(tokens []tableToken, allTables []*metadata.TableMetadata, depth int) []string {
	for len(tokens) >= 2 && tokens[0].paren == '(' && findMatchingTokenParen(tokens, 0) == len(tokens)-1 {
		tokens = tokens[1 : len(tokens)-1]
	}
	selIdx := findMainSelectIndex(tokens)
	if selIdx < 0 {
		return nil
	}
	listStart := selIdx + 1
	listEnd := len(tokens)
	parenDepth := 0
	for i := listStart; i < len(tokens); i++ {
		switch tokens[i].paren {
		case '(':
			parenDepth++
			continue
		case ')':
			if parenDepth > 0 {
				parenDepth--
			}
			continue
		}
		if parenDepth == 0 && !tokens[i].quoted && isSelectListEndKeyword(strings.ToUpper(tokens[i].text)) {
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
	if itemStart < listEnd && !tokens[itemStart].quoted {
		firstUpper := strings.ToUpper(tokens[itemStart].text)
		if firstUpper == "DISTINCT" || firstUpper == "ALL" || firstUpper == "VALUE" {
			itemStart++
		}
	}
	parenDepth = 0
	for i := itemStart; i <= listEnd; i++ {
		if i < listEnd {
			switch tokens[i].paren {
			case '(':
				parenDepth++
				continue
			case ')':
				if parenDepth > 0 {
					parenDepth--
				}
				continue
			}
		}
		if i < listEnd && (parenDepth != 0 || !tokens[i].comma) {
			continue
		}
		item := tokens[itemStart:i]
		itemStart = i + 1
		if len(item) == 0 {
			continue
		}
		if len(item) == 1 && !item[0].quoted && item[0].text == "*" {
			for _, ref := range getInnerRefs() {
				cols = append(cols, resolveTableRefColumns(ref, allTables, depth)...)
			}
			continue
		}
		if len(item) == 1 && !item[0].quoted && strings.HasSuffix(item[0].text, ".*") {
			qual := strings.TrimSuffix(item[0].text, ".*")
			for _, ref := range getInnerRefs() {
				if matchesTableQualifier(ref, qual) {
					cols = append(cols, resolveTableRefColumns(ref, allTables, depth)...)
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

func columnQualifierBeforeCursor(statement string, pos int) string {
	prefix := extractPrefix(statement, pos)
	before := strings.TrimRight(statement[:pos-len(prefix)], " \t\n\r")
	if !strings.HasSuffix(before, ".") {
		return ""
	}
	before = strings.TrimRight(strings.TrimSuffix(before, "."), " \t\n\r")
	if before == "" {
		return ""
	}
	if before[len(before)-1] == '"' {
		for i := len(before) - 2; i >= 0; i-- {
			if before[i] == '"' {
				if i > 0 && before[i-1] == '"' {
					i--
					continue
				}
				return strings.ReplaceAll(before[i+1:len(before)-1], `""`, `"`)
			}
		}
		return ""
	}
	start := len(before)
	for start > 0 && isIdentChar(before[start-1]) {
		start--
	}
	return before[start:]
}

func collectMetadataTables(meta *metadata.DatabaseSchemaMetadata) []*metadata.TableMetadata {
	var tables []*metadata.TableMetadata
	for _, schema := range meta.GetSchemas() {
		tables = append(tables, schema.GetTables()...)
	}
	return tables
}

func findMetadataTable(tables []*metadata.TableMetadata, name string) *metadata.TableMetadata {
	for _, t := range tables {
		if t.GetName() == name {
			return t
		}
	}
	for _, t := range tables {
		if strings.EqualFold(t.GetName(), name) {
			return t
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

// extractQueryScopes returns query block slices enclosing pos, ordered from
// innermost to outermost, isolating the active UNION/INTERSECT/EXCEPT arm at pos.
// Only parentheses that introduce a subquery (SELECT or WITH) start a new query
// scope; expression parentheses stay in the enclosing query scope.
func extractQueryScopes(statement string, pos int) []string {
	subqueryStarts := findActiveSubqueryStarts(statement, pos)
	var scopes []string
	for k := len(subqueryStarts) - 1; k >= 0; k-- {
		start := subqueryStarts[k]
		end := findSubqueryEnd(statement, start)
		scopes = append(scopes, isolateSetOperationArm(statement[start:end], pos-start))
	}
	scopes = append(scopes, isolateSetOperationArm(statement, pos))
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
	starts := findActiveSubqueryStarts(statement, pos)
	if len(starts) == 0 {
		return 0
	}
	return starts[len(starts)-1]
}

func findActiveSubqueryStarts(statement string, pos int) []int {
	var subqueryStarts []int
	var isSubqueryParen []bool
	for i := 0; i < len(statement) && i < pos; i++ {
		switch {
		case statement[i] == '\'':
			i = skipQuotedLiteral(statement, i, '\'')
		case statement[i] == '"':
			i = skipQuotedLiteral(statement, i, '"')
		case statement[i] == '-' && i+1 < len(statement) && statement[i+1] == '-':
			i = skipLineComment(statement, i)
		case statement[i] == '/' && i+1 < len(statement) && statement[i+1] == '*':
			i, _ = skipBlockComment(statement, i)
		case statement[i] == '(':
			sub := parenOpensSubquery(statement, i+1)
			isSubqueryParen = append(isSubqueryParen, sub)
			if sub {
				subqueryStarts = append(subqueryStarts, i+1)
			}
		case statement[i] == ')':
			if len(isSubqueryParen) > 0 {
				wasSub := isSubqueryParen[len(isSubqueryParen)-1]
				isSubqueryParen = isSubqueryParen[:len(isSubqueryParen)-1]
				if wasSub && len(subqueryStarts) > 0 {
					subqueryStarts = subqueryStarts[:len(subqueryStarts)-1]
				}
			}
		}
	}
	return subqueryStarts
}

func findSubqueryEnd(statement string, start int) int {
	depth := 0
	for i := start; i < len(statement); i++ {
		switch {
		case statement[i] == '\'':
			i = skipQuotedLiteral(statement, i, '\'')
		case statement[i] == '"':
			i = skipQuotedLiteral(statement, i, '"')
		case statement[i] == '-' && i+1 < len(statement) && statement[i+1] == '-':
			i = skipLineComment(statement, i)
		case statement[i] == '/' && i+1 < len(statement) && statement[i+1] == '*':
			i, _ = skipBlockComment(statement, i)
		case statement[i] == '(':
			depth++
		case statement[i] == ')':
			if depth == 0 {
				return i
			}
			depth--
		}
	}
	return len(statement)
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
	i := start + 1
	for i < len(s) {
		if s[i] == quote {
			if i+1 < len(s) && s[i+1] == quote {
				i += 2
				continue
			}
			return i
		}
		i++
	}
	return len(s) - 1
}

func parenOpensSubquery(s string, pos int) bool {
	for pos < len(s) {
		c := s[pos]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '(' {
			pos++
			continue
		}
		if c == '-' && pos+1 < len(s) && s[pos+1] == '-' {
			pos = skipLineComment(s, pos) + 1
			continue
		}
		if c == '/' && pos+1 < len(s) && s[pos+1] == '*' {
			end, closed := skipBlockComment(s, pos)
			if !closed {
				return false
			}
			pos = end + 1
			continue
		}
		break
	}
	end := pos
	for end < len(s) && isIdentChar(s[end]) {
		end++
	}
	word := strings.ToUpper(s[pos:end])
	return word == "SELECT" || word == "WITH"
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
		if (upper == "FROM" || upper == "JOIN" || upper == "INTO" || upper == "UPDATE") && i+1 < len(tokens) && !tokens[i].quoted {
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
		alias, nextIdx := parseOptionalAlias(tokens, nextIdx)
		if nextIsSubqueryStart(inner, 0) {
			return []tableRef{{alias: alias, subquery: inner}}, nextIdx, true
		}
		innerRefs := extractRefsFromParenthesizedExpr(inner)
		if len(innerRefs) == 1 && alias != "" && innerRefs[0].alias == "" {
			innerRefs[0].alias = alias
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
	alias, nextIdx := parseOptionalAlias(tokens, i+1)
	return []tableRef{{name: tok.text, alias: alias}}, nextIdx, true
}

func parseOptionalAlias(tokens []tableToken, i int) (string, int) {
	if i >= len(tokens) || tokens[i].comma || tokens[i].paren != 0 {
		return "", i
	}
	if !tokens[i].quoted && strings.ToUpper(tokens[i].text) == "AS" {
		if i+1 < len(tokens) && !tokens[i+1].comma && tokens[i+1].paren == 0 {
			return tokens[i+1].text, i + 2
		}
		return "", i + 1
	}
	if tokens[i].quoted || !isKeyword(strings.ToUpper(tokens[i].text)) {
		return tokens[i].text, i + 1
	}
	return "", i
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
			s[i] != '(' && s[i] != ')' {
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
		"THEN", "ELSE", "END", "EXISTS", "ASC", "DESC", "EXPLAIN":
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
	"EXISTS", "ASC", "DESC",
}
