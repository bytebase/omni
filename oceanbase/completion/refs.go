package completion

import (
	"errors"
	"strings"

	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/oceanbase/parser"
)

// TableRef is a table reference found in a SQL statement.
type TableRef struct {
	Database string // database/schema qualifier
	Table    string // table name
	Alias    string // AS alias
}

// extractTableRefs returns the table references visible at the cursor. It reads
// only the statement that holds the cursor, so another statement in the script
// that does not parse, such as one still being typed, cannot cost this one its
// table scope, and the tables another statement names do not leak into it.
func extractTableRefs(sql string, cursorOffset int) (refs []TableRef) {
	stmt, cursor := cursorStatement(sql, cursorOffset)
	defer func() {
		if r := recover(); r != nil {
			refs = extractTableRefsLexer(stmt)
		}
	}()
	return extractTableRefsInner(stmt, cursor)
}

// cursorStatement returns the statement that holds cursorOffset and the cursor's
// offset within it. It picks the text parser.Collect completes, so the table
// references come from the statement the candidates are for: the segment of
// parser.Split that contains the cursor, else the text after the last segment
// when the cursor is past it, else the whole input. A TableRef carries no
// position, so the statement is parsed on its own.
func cursorStatement(sql string, cursorOffset int) (string, int) {
	segments := parser.Split(sql)
	for _, seg := range segments {
		if seg.Empty() {
			continue
		}
		if cursorOffset >= seg.ByteStart && cursorOffset <= seg.ByteEnd {
			return seg.Text, cursorOffset - seg.ByteStart
		}
	}
	if len(segments) > 0 {
		last := segments[len(segments)-1]
		if last.ByteEnd < len(sql) && cursorOffset >= last.ByteEnd {
			return sql[last.ByteEnd:], cursorOffset - last.ByteEnd
		}
	}
	return sql, cursorOffset
}

func extractTableRefsInner(sql string, cursorOffset int) []TableRef {
	list, err := parseClosingTokens(sql)
	if err != nil || list == nil || len(list.Items) == 0 {
		// Fallback: try with placeholder appended at cursor.
		if cursorOffset <= len(sql) {
			patched := sql[:cursorOffset] + " _x"
			if cursorOffset < len(sql) {
				patched += sql[cursorOffset:]
			}
			list, err = parseClosingTokens(patched)
			if err != nil || list == nil {
				return extractTableRefsLexer(sql)
			}
		} else {
			return nil
		}
	}

	var refs []TableRef
	for _, item := range list.Items {
		refs = append(refs, extractRefsFromStmt(item)...)
	}
	if len(refs) == 0 {
		return extractTableRefsLexer(sql)
	}
	return refs
}

// maxClosedTokens bounds how many unterminated tokens parseClosingTokens closes.
// A string inside an executable comment that is itself unclosed takes two, and
// a string that ends in a lone backslash takes its quote twice, because the
// first quote is read as escaped.
const maxClosedTokens = 4

// parseClosingTokens parses a statement whose last token may still be open: a
// string, quoted identifier or comment that runs to the end of the text without
// its closing delimiter. parser.Parse rejects such a token, as the server does,
// but in a statement being edited it is the usual state of whatever the user is
// typing, often after the cursor. So each time Parse reports one, this appends
// the delimiter that closes it and parses again, completing the fragment rather
// than loosening Parse. The tokens before the open one, where the table
// references are, are lexed the same either way. Any other error is returned
// as is.
func parseClosingTokens(sql string) (*ast.List, error) {
	list, err := parser.Parse(sql)
	for range maxClosedTokens {
		var pe *parser.ParseError
		if err == nil || !errors.As(err, &pe) || !strings.HasPrefix(pe.Message, "unterminated ") {
			break
		}
		closer := closingDelimiter(sql, pe.Position)
		if closer == "" {
			break
		}
		sql += closer
		list, err = parser.Parse(sql)
	}
	return list, err
}

// closingDelimiter returns the text that closes the token starting at pos, which
// the lexer reported as unterminated and which therefore runs to the end of sql,
// or "" when no such token starts at pos.
func closingDelimiter(sql string, pos int) string {
	if pos < 0 || pos >= len(sql) {
		return ""
	}
	switch sql[pos] {
	case '\'', '"', '`':
		return sql[pos : pos+1]
	case '@': // a user variable quoted as @`name`
		return "`"
	case '/':
		if strings.HasPrefix(sql[pos:], "/*") {
			// The body of an executable comment can end in a line comment,
			// which would swallow a close on the same line.
			return "\n*/"
		}
	}
	return ""
}

// extractRefsFromStmt extracts table references from a single statement node.
func extractRefsFromStmt(n ast.Node) []TableRef {
	if n == nil {
		return nil
	}
	switch v := n.(type) {
	case *ast.SelectStmt:
		return extractRefsFromSelect(v)
	case *ast.InsertStmt:
		return extractRefsFromInsert(v)
	case *ast.UpdateStmt:
		return extractRefsFromUpdate(v)
	case *ast.DeleteStmt:
		return extractRefsFromDelete(v)
	}
	return nil
}

// extractRefsFromSelect extracts table references from a SELECT statement.
// Does not recurse into subqueries (their tables don't leak to the outer scope).
func extractRefsFromSelect(s *ast.SelectStmt) []TableRef {
	if s == nil {
		return nil
	}
	var refs []TableRef

	// Set operations: walk both sides.
	if s.SetOp != ast.SetOpNone {
		refs = append(refs, extractRefsFromSelect(s.Left)...)
		refs = append(refs, extractRefsFromSelect(s.Right)...)
		return refs
	}

	// Parenthesized query: walk the inner query.
	if s.ParenSource != nil {
		return append(refs, extractRefsFromSelect(s.ParenSource)...)
	}

	// CTEs.
	for _, cte := range s.CTEs {
		if cte != nil && cte.Name != "" {
			refs = append(refs, TableRef{Table: cte.Name})
		}
	}

	// FROM clause.
	for _, te := range s.From {
		refs = append(refs, extractRefsFromTableExpr(te)...)
	}
	return refs
}

// extractRefsFromTableExpr extracts table references from a TableExpr node.
func extractRefsFromTableExpr(te ast.TableExpr) []TableRef {
	if te == nil {
		return nil
	}
	switch v := te.(type) {
	case *ast.TableRef:
		if v.Name != "" {
			return []TableRef{{Database: v.Schema, Table: v.Name, Alias: v.Alias}}
		}
	case *ast.JoinClause:
		var refs []TableRef
		refs = append(refs, extractRefsFromTableExpr(v.Left)...)
		refs = append(refs, extractRefsFromTableExpr(v.Right)...)
		return refs
	case *ast.SubqueryExpr:
		// Subquery tables don't leak to the outer scope.
		return nil
	}
	return nil
}

// extractRefsFromInsert extracts the target table from an INSERT statement.
func extractRefsFromInsert(s *ast.InsertStmt) []TableRef {
	if s == nil || s.Table == nil {
		return nil
	}
	return []TableRef{{Database: s.Table.Schema, Table: s.Table.Name, Alias: s.Table.Alias}}
}

// extractRefsFromUpdate extracts table references from an UPDATE statement.
func extractRefsFromUpdate(s *ast.UpdateStmt) []TableRef {
	if s == nil {
		return nil
	}
	var refs []TableRef
	for _, te := range s.Tables {
		refs = append(refs, extractRefsFromTableExpr(te)...)
	}
	return refs
}

// extractRefsFromDelete extracts table references from a DELETE statement.
func extractRefsFromDelete(s *ast.DeleteStmt) []TableRef {
	if s == nil {
		return nil
	}
	var refs []TableRef
	for _, te := range s.Tables {
		refs = append(refs, extractRefsFromTableExpr(te)...)
	}
	for _, te := range s.Using {
		refs = append(refs, extractRefsFromTableExpr(te)...)
	}
	return refs
}

// extractTableRefsLexer is a fallback using lexer-based pattern matching
// when the SQL doesn't parse (e.g., incomplete SQL being edited). The caller
// passes only the statement that holds the cursor, so it reads all of it: the
// FROM clause of `SELECT | FROM t WHERE ...` comes after the cursor.
func extractTableRefsLexer(sql string) []TableRef {
	lex := parser.NewLexer(sql)
	var tokens []parser.Token
	for {
		tok := lex.NextToken()
		if tok.Type == 0 {
			break
		}
		tokens = append(tokens, tok)
	}

	var refs []TableRef

	for i := 0; i < len(tokens); i++ {
		typ := tokens[i].Type

		// FROM table, JOIN table
		if typ == parser.FROM || typ == parser.JOIN {
			ref, skip := lexerExtractTableAfter(tokens, i+1)
			if ref != nil {
				refs = append(refs, *ref)
			}
			i += skip
			continue
		}

		// UPDATE table
		if typ == parser.UPDATE {
			ref, skip := lexerExtractTableAfter(tokens, i+1)
			if ref != nil {
				refs = append(refs, *ref)
			}
			i += skip
			continue
		}

		// INSERT [INTO] table / REPLACE [INTO] table
		if typ == parser.INSERT || typ == parser.REPLACE {
			j := i + 1
			if j < len(tokens) && tokens[j].Type == parser.INTO {
				j++
			}
			ref, skip := lexerExtractTableAfter(tokens, j)
			if ref != nil {
				refs = append(refs, *ref)
			}
			i = j + skip
			continue
		}

		// DELETE [FROM] table
		if typ == parser.DELETE {
			j := i + 1
			if j < len(tokens) && tokens[j].Type == parser.FROM {
				j++
			}
			ref, skip := lexerExtractTableAfter(tokens, j)
			if ref != nil {
				refs = append(refs, *ref)
			}
			i = j + skip
			continue
		}
	}
	return refs
}

// lexerExtractTableAfter extracts a [db.]table [AS alias] reference starting at tokens[j].
// Returns the TableRef (or nil) and the number of tokens consumed.
func lexerExtractTableAfter(tokens []parser.Token, j int) (*TableRef, int) {
	if j >= len(tokens) || !parser.IsIdentTokenType(tokens[j].Type) {
		return nil, 0
	}
	ref := TableRef{Table: tokens[j].Str}
	consumed := j + 1
	// Check for db.table
	if j+2 < len(tokens) && tokens[j+1].Type == '.' && parser.IsIdentTokenType(tokens[j+2].Type) {
		ref.Database = ref.Table
		ref.Table = tokens[j+2].Str
		consumed = j + 3
	}
	// Check for alias (AS alias or bare ident)
	k := consumed
	if k < len(tokens) {
		if tokens[k].Type == parser.AS && k+1 < len(tokens) && parser.IsIdentTokenType(tokens[k+1].Type) {
			ref.Alias = tokens[k+1].Str
		} else if parser.IsIdentTokenType(tokens[k].Type) && !isClauseKeyword(tokens[k].Type) {
			ref.Alias = tokens[k].Str
		}
	}
	return &ref, consumed - j
}

// isClauseKeyword returns true for keywords that typically follow a table name
// and should not be treated as aliases.
func isClauseKeyword(typ int) bool {
	switch typ {
	case parser.SET, parser.WHERE, parser.ON, parser.VALUES, parser.JOIN,
		parser.INNER, parser.LEFT, parser.RIGHT, parser.CROSS, parser.NATURAL,
		parser.ORDER, parser.GROUP, parser.HAVING, parser.LIMIT, parser.UNION,
		parser.FOR, parser.USING, parser.FROM, parser.INTO,
		parser.SELECT, parser.INSERT, parser.UPDATE, parser.DELETE:
		return true
	}
	return false
}
