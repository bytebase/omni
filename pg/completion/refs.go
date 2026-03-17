package completion

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
)

// TableRef is a table reference found in a FROM clause.
type TableRef struct {
	Schema string
	Table  string
	Alias  string
}

// extractTableRefs parses the SQL and returns table references visible at cursor.
func extractTableRefs(sql string, cursorOffset int) []TableRef {
	// Try parsing the full SQL.
	list, err := parser.Parse(sql)
	if err != nil || list == nil || len(list.Items) == 0 {
		// Fallback: try with placeholder
		if cursorOffset <= len(sql) {
			patched := sql[:cursorOffset] + " 1"
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
		raw, ok := item.(*ast.RawStmt)
		if !ok || raw.Stmt == nil {
			continue
		}
		ast.Inspect(raw.Stmt, func(n ast.Node) bool {
			switch v := n.(type) {
			case *ast.RangeVar:
				if v != nil && v.Relname != "" {
					ref := TableRef{Schema: v.Schemaname, Table: v.Relname}
					if v.Alias != nil {
						ref.Alias = v.Alias.Aliasname
					}
					refs = append(refs, ref)
				}
			case *ast.CommonTableExpr:
				if v != nil && v.Ctename != "" {
					refs = append(refs, TableRef{Table: v.Ctename})
				}
			}
			return true
		})
	}
	return refs
}

// extractTableRefsLexer is a fallback using lexer-based pattern matching.
func extractTableRefsLexer(sql string, cursorOffset int) []TableRef {
	lex := parser.NewLexer(sql)
	var tokens []parser.Token
	for {
		tok := lex.NextToken()
		if tok.Type == 0 || tok.Loc >= cursorOffset {
			break
		}
		mapped := parser.MapTokenType(tok.Type)
		tok.Type = mapped
		tokens = append(tokens, tok)
	}

	var refs []TableRef
	for i := 0; i < len(tokens); i++ {
		typ := tokens[i].Type
		if typ != parser.FROM && typ != parser.JOIN {
			continue
		}
		j := i + 1
		// Skip additional JOIN keyword if present (e.g. INNER JOIN, LEFT JOIN)
		if j < len(tokens) && tokens[j].Type == parser.JOIN {
			j++
		}
		if j >= len(tokens) || tokens[j].Type != parser.IDENT {
			continue
		}
		ref := TableRef{Table: tokens[j].Str}
		// Check for schema.table
		if j+2 < len(tokens) && tokens[j+1].Type == '.' && tokens[j+2].Type == parser.IDENT {
			ref.Schema = ref.Table
			ref.Table = tokens[j+2].Str
			j += 2
		}
		// Check for alias (AS alias or just alias)
		if j+1 < len(tokens) {
			if tokens[j+1].Type == parser.AS && j+2 < len(tokens) && tokens[j+2].Type == parser.IDENT {
				ref.Alias = tokens[j+2].Str
			} else if tokens[j+1].Type == parser.IDENT {
				ref.Alias = tokens[j+1].Str
			}
		}
		refs = append(refs, ref)
	}
	return refs
}
