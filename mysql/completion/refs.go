package completion

import (
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/parser"
)

// TableRef is a table reference found in a SQL statement.
type TableRef struct {
	Database string // database/schema qualifier
	Table    string // table name
	Alias    string // AS alias
}

// extractTableRefs parses the SQL and returns table references visible at cursor.
func extractTableRefs(sql string, cursorOffset int) (refs []TableRef) {
	defer func() {
		if r := recover(); r != nil {
			refs = extractTableRefsLexer(sql, cursorOffset)
		}
	}()
	return extractTableRefsInner(sql, cursorOffset)
}

func extractTableRefsInner(sql string, cursorOffset int) []TableRef {
	list, err := parser.Parse(sql)
	if err != nil || list == nil || len(list.Items) == 0 {
		// Fallback: try with placeholder appended at cursor.
		if cursorOffset <= len(sql) {
			patched := sql[:cursorOffset] + " _x"
			if cursorOffset < len(sql) {
				patched += sql[cursorOffset:]
			}
			list, err = parser.Parse(patched)
			if err != nil || list == nil {
				return extractTableRefsLexer(sql, cursorOffset)
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
		return extractTableRefsLexer(sql, cursorOffset)
	}
	return refs
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
// when the SQL doesn't parse (e.g., incomplete SQL being edited).
func extractTableRefsLexer(sql string, cursorOffset int) []TableRef {
	lex := parser.NewLexer(sql)
	var tokens []parser.Token
	for {
		tok := lex.NextToken()
		if tok.Type == 0 || tok.Loc >= cursorOffset {
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
