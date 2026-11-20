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

	// Strip the partial word being typed so context detection inspects the
	// token before it.
	beforeCursor := strings.ToUpper(stmtSQL[:localPos-len(prefix)])
	inFromContext := isInFromContext(beforeCursor)

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

	if meta != nil && isColumnContext(stmtSQL, localPos) {
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

func isInFromContext(upper string) bool {
	// The cursor is in a table-name position when the last completed clause
	// keyword is FROM / JOIN / INTO. Each must match as a STANDALONE keyword:
	// a dotted path step such as "x.from" (member access) also ends in "FROM",
	// but `from` after a '.' is a path key, not a FROM clause, so it must NOT
	// trigger table suggestions.
	upper = strings.TrimRight(upper, " \t\n\r")
	return endsWithKeyword(upper, "FROM") ||
		endsWithKeyword(upper, "JOIN") ||
		endsWithKeyword(upper, "INTO")
}

// endsWithKeyword reports whether s ends with kw as a standalone keyword: kw is
// a suffix of s, and either kw is the whole string or the character immediately
// before it is whitespace (a real clause boundary).
func endsWithKeyword(s, kw string) bool {
	if !strings.HasSuffix(s, kw) {
		return false
	}
	if len(s) == len(kw) {
		return true
	}
	switch s[len(s)-len(kw)-1] {
	case ' ', '\t', '\n', '\r':
		return true
	default:
		return false
	}
}

func matchesPrefix(text, prefix string) bool {
	if prefix == "" {
		return true
	}
	return strings.HasPrefix(strings.ToUpper(text), strings.ToUpper(prefix))
}

// isColumnContext reports whether pos appears to be in a context where column
// names should be suggested (between SELECT and FROM, or inside expression
// clauses after FROM such as WHERE, HAVING, SET, ORDER BY, GROUP BY).
// Table-only contexts immediately after FROM/JOIN/INTO return false.
func isColumnContext(statement string, pos int) bool {
	if pos > len(statement) {
		pos = len(statement)
	}
	// Mask comments, string literals, and quoted identifiers before uppercasing
	// so neither literals like 'FROM' nor identifiers like "order" match keywords.
	before := strings.ToUpper(maskCommentsStringsAndQuotedIdents(statement[:pos]))

	trimmed := strings.TrimRight(before, "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_$\"")
	trimmed = strings.TrimRight(trimmed, " \t\n\r")

	if strings.HasSuffix(trimmed, "FROM") ||
		strings.HasSuffix(trimmed, "JOIN") ||
		strings.HasSuffix(trimmed, "INTO") {
		return false
	}

	selectIdx := lastKeywordIndex(before, "SELECT")
	fromIdx := lastKeywordIndex(before, "FROM")
	if selectIdx >= 0 && (fromIdx < 0 || fromIdx < selectIdx) {
		return true
	}

	for _, kw := range []string{"WHERE", "HAVING", "SET", "ORDER", "GROUP", "BY", "AND", "OR", "ON", "WHEN", "THEN", "ELSE"} {
		kwIdx := lastKeywordIndex(before, kw)
		if kwIdx >= 0 && kwIdx > fromIdx {
			return true
		}
	}

	return false
}

// lastKeywordIndex returns the byte index of the last standalone occurrence of
// kw in s (not inside an identifier or after a '.' path step), or -1 if absent.
// Both s and kw must be uppercase.
func lastKeywordIndex(s, kw string) int {
	search := s
	for {
		idx := strings.LastIndex(search, kw)
		if idx < 0 {
			return -1
		}
		end := idx + len(kw)
		startOk := idx == 0 || (!isIdentRune(rune(search[idx-1])) && search[idx-1] != '.')
		endOk := end >= len(search) || !isIdentRune(rune(search[end]))
		if startOk && endOk {
			return idx
		}
		search = search[:idx]
	}
}

func isIdentRune(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') ||
		(r >= '0' && r <= '9') || r == '_' || r == '$'
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
		if last.paren == 0 && !last.comma && last.text != "" {
			return last.text
		}
	}
	if len(item) >= 2 {
		last := item[len(item)-1]
		if last.paren == 0 && !last.comma && last.text != "" && (last.quoted || !isKeyword(strings.ToUpper(last.text))) {
			return last.text
		}
	}
	if len(item) == 1 && item[0].paren == 0 && !item[0].comma && item[0].text != "" {
		name := item[0].text
		if !item[0].quoted {
			if dot := strings.LastIndexByte(name, '.'); dot >= 0 && dot+1 < len(name) {
				name = name[dot+1:]
			}
			if strings.EqualFold(name, "NULL") || strings.EqualFold(name, "MISSING") ||
				strings.EqualFold(name, "TRUE") || strings.EqualFold(name, "FALSE") {
				return ""
			}
		}
		return name
	}
	return ""
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
// innermost to outermost. Only parentheses that introduce a subquery (SELECT or
// WITH) start a new query scope; expression parentheses stay in the enclosing
// query scope.
func extractQueryScopes(statement string, pos int) []string {
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

	var scopes []string
	for k := len(subqueryStarts) - 1; k >= 0; k-- {
		start := subqueryStarts[k]
		scopes = append(scopes, statement[start:findSubqueryEnd(statement, start)])
	}
	scopes = append(scopes, statement)
	return scopes
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

func maskCommentsStringsAndQuotedIdents(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	i := 0
	for i < len(s) {
		switch {
		case s[i] == '"' || s[i] == '\'':
			q := s[i]
			b.WriteByte('_')
			i++
			for i < len(s) {
				if s[i] == q {
					i++
					if i < len(s) && s[i] == q {
						b.WriteByte('_')
						b.WriteByte('_')
						i++
						continue
					}
					b.WriteByte('_')
					break
				}
				b.WriteByte('_')
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
