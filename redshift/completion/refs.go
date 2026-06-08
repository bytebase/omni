package completion

import (
	"github.com/bytebase/omni/redshift/ast"
	"github.com/bytebase/omni/redshift/parser"
)

// TableRef is a table reference found in a FROM clause.
type TableRef struct {
	Schema string
	Table  string
	Alias  string
}

// extractTableRefs parses the SQL and returns table references visible at cursor.
func extractTableRefs(sql string, cursorOffset int) (refs []TableRef) {
	// Guard against nil-typed AST nodes that cause panics in Walk.
	defer func() {
		if r := recover(); r != nil {
			// Fall back to lexer-based extraction on AST walk panic.
			refs = extractTableRefsLexer(sql, cursorOffset)
		}
	}()
	return extractTableRefsInner(sql, cursorOffset)
}

func extractTableRefsInner(sql string, cursorOffset int) []TableRef {
	// Try parsing the full SQL.
	list, err := parser.Parse(sql)
	if err != nil || list == nil || len(list.Items) == 0 {
		// Fallback: try with placeholder
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
	// AST walk may miss table references in DDL statements where the table
	// is not stored as a RangeVar (e.g. COMMENT ON COLUMN). Fall back to
	// lexer-based extraction when AST finds nothing.
	if len(refs) == 0 {
		return extractTableRefsLexer(sql, cursorOffset)
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

		// DML: FROM table, JOIN table
		if typ == parser.FROM || typ == parser.JOIN {
			ref, skip := lexerExtractTableAfter(tokens, i+1)
			if ref != nil {
				refs = append(refs, *ref)
			}
			i += skip
			continue
		}

		// DDL: ALTER TABLE [IF EXISTS] [schema.]table
		if typ == parser.ALTER && i+1 < len(tokens) && tokens[i+1].Type == parser.TABLE {
			j := i + 2
			// Skip optional IF EXISTS
			if j < len(tokens) && tokens[j].Type == parser.IF_P {
				j++ // IF
				if j < len(tokens) && tokens[j].Type == parser.EXISTS {
					j++ // EXISTS
				}
			}
			ref, skip := lexerExtractTableAfter(tokens, j)
			if ref != nil {
				refs = append(refs, *ref)
			}
			i = j + skip
			continue
		}

		// DDL: CREATE [UNIQUE] INDEX ... ON [schema.]table
		if typ == parser.CREATE {
			j := i + 1
			if j < len(tokens) && tokens[j].Type == parser.UNIQUE {
				j++
			}
			if j < len(tokens) && tokens[j].Type == parser.INDEX {
				// Scan forward to ON keyword
				for k := j + 1; k < len(tokens); k++ {
					if tokens[k].Type == parser.ON {
						ref, skip := lexerExtractTableAfter(tokens, k+1)
						if ref != nil {
							refs = append(refs, *ref)
						}
						i = k + skip
						break
					}
				}
			}
			continue
		}

		// DDL: COMMENT ON COLUMN [schema.]table.column
		// Handles both complete (table.col) and incomplete (table.) forms.
		if typ == parser.COMMENT {
			if i+2 < len(tokens) && tokens[i+1].Type == parser.ON && tokens[i+2].Type == parser.COLUMN {
				j := i + 3
				if j < len(tokens) && isIdentToken(tokens[j].Type) {
					name1 := tokens[j].Str
					if j+1 < len(tokens) && tokens[j+1].Type == '.' {
						if j+2 < len(tokens) && isIdentToken(tokens[j+2].Type) {
							name2 := tokens[j+2].Str
							if j+3 < len(tokens) && tokens[j+3].Type == '.' {
								// schema.table[.col] — name1=schema, name2=table
								refs = append(refs, TableRef{Schema: name1, Table: name2})
							} else {
								// table.col — name1=table
								refs = append(refs, TableRef{Table: name1})
							}
						} else {
							// table. (incomplete, cursor after dot) — name1=table
							refs = append(refs, TableRef{Table: name1})
						}
					}
				}
			}
			continue
		}
	}
	return refs
}

// lexerExtractTableAfter extracts a [schema.]table reference starting at tokens[j].
// Returns the TableRef (or nil) and the number of tokens consumed.
func lexerExtractTableAfter(tokens []parser.Token, j int) (*TableRef, int) {
	// Skip additional JOIN keyword if present (e.g. INNER JOIN, LEFT JOIN)
	if j < len(tokens) && tokens[j].Type == parser.JOIN {
		j++
	}
	if j >= len(tokens) || !isIdentToken(tokens[j].Type) {
		return nil, 0
	}
	ref := TableRef{Table: tokens[j].Str}
	consumed := j + 1
	// Check for schema.table
	if j+2 < len(tokens) && tokens[j+1].Type == '.' && isIdentToken(tokens[j+2].Type) {
		ref.Schema = ref.Table
		ref.Table = tokens[j+2].Str
		consumed = j + 3
	}
	// Check for alias (AS alias or just alias)
	k := consumed
	if k < len(tokens) {
		if tokens[k].Type == parser.AS && k+1 < len(tokens) && isIdentToken(tokens[k+1].Type) {
			ref.Alias = tokens[k+1].Str
		} else if isIdentToken(tokens[k].Type) {
			ref.Alias = tokens[k].Str
		}
	}
	return &ref, consumed - j
}

// isIdentToken reports whether a token type can be an identifier (IDENT or non-reserved keyword).
func isIdentToken(typ int) bool {
	return parser.IsIdentifierTokenType(typ)
}
