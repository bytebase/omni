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

	prefix := extractPrefix(input, pos)

	// Strip the partial word being typed so context detection inspects the
	// token before it.
	beforeCursor := strings.ToUpper(input[:pos-len(prefix)])
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

	if meta != nil && isColumnContext(input, pos) {
		for _, col := range resolveColumnNames(input, pos, meta) {
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
	// Strip quoted identifiers before uppercasing so "order" does not match
	// the ORDER keyword in lastKeywordIndex.
	before := strings.ToUpper(stripQuotedIdentifiers(statement[:pos]))

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

// resolveColumnNames extracts column names from tables referenced in the active
// statement around pos. If no FROM/JOIN/INTO table is found, it falls back to
// all tables in meta.
func resolveColumnNames(statement string, pos int, meta *metadata.DatabaseSchemaMetadata) []string {
	if meta == nil {
		return nil
	}
	currentStmt := extractCurrentStatement(statement, pos)
	targetTables := extractTableNamesFromStatement(currentStmt)

	allTables := collectMetadataTables(meta)
	if len(targetTables) == 0 {
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
	for _, tableName := range targetTables {
		table := findMetadataTable(allTables, tableName)
		if table == nil {
			continue
		}
		for _, col := range table.GetColumns() {
			if col.GetName() != "" {
				cols = append(cols, col.GetName())
			}
		}
	}
	return cols
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

// extractCurrentStatement returns the statement segment containing pos using
// partiql.Split so semicolons inside strings, comments, or Ion literals are
// respected.
func extractCurrentStatement(statement string, pos int) string {
	segs := partiql.Split(statement)
	for _, seg := range segs {
		if pos >= seg.ByteStart && pos <= seg.ByteEnd {
			return seg.Text
		}
	}
	return statement
}

// extractTableNamesFromStatement scans a statement for table names after FROM,
// JOIN, or INTO keywords, handling comma-separated tables, quoted identifiers,
// aliases, and comments.
func extractTableNamesFromStatement(statement string) []string {
	cleaned := stripSQLComments(statement)
	tokens := tokenizeForTableExtraction(cleaned)

	var tables []string
	i := 0
	for i < len(tokens) {
		upper := strings.ToUpper(tokens[i].text)
		if (upper == "FROM" || upper == "JOIN" || upper == "INTO" || upper == "UPDATE") && i+1 < len(tokens) {
			i++
			for i < len(tokens) {
				tok := tokens[i]
				name := tok.text
				if tok.quoted {
					tables = append(tables, name)
				} else if isKeyword(strings.ToUpper(name)) {
					break
				} else {
					tables = append(tables, name)
				}
				i++
				if i < len(tokens) && strings.ToUpper(tokens[i].text) == "AS" {
					i += 2
				} else if i < len(tokens) && !tokens[i].comma && !isKeyword(strings.ToUpper(tokens[i].text)) {
					i++
				}
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
	return tables
}

type tableToken struct {
	text   string
	quoted bool
	comma  bool
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

func stripSQLComments(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	i := 0
	for i < len(s) {
		switch {
		case s[i] == '\'':
			b.WriteByte(s[i])
			i++
			for i < len(s) {
				b.WriteByte(s[i])
				if s[i] == '\'' {
					i++
					if i < len(s) && s[i] == '\'' {
						b.WriteByte(s[i])
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

func stripQuotedIdentifiers(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	i := 0
	for i < len(s) {
		if s[i] == '"' {
			b.WriteByte('_')
			i++
			for i < len(s) {
				if s[i] == '"' {
					i++
					if i < len(s) && s[i] == '"' {
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
		} else {
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
