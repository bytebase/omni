package completion

import (
	"strings"

	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/mssql/parser"
)

// TableRef is a table reference found in a FROM clause or DML target.
type TableRef struct {
	Schema string
	Table  string
	Alias  string
}

// extractTableRefs parses the SQL and returns table references visible at cursor.
func extractTableRefs(sql string, cursorOffset int) (refs []TableRef) {
	// Guard against panics in AST walk.
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
		// Fallback: try appending a placeholder to make incomplete SQL parseable.
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
		refs = append(refs, collectTableRefsFromNode(item)...)
	}
	if len(refs) == 0 {
		return extractTableRefsLexer(sql, cursorOffset)
	}
	return refs
}

// collectTableRefsFromNode recursively collects TableRef values from an AST node.
func collectTableRefsFromNode(n ast.Node) []TableRef {
	if n == nil {
		return nil
	}
	var refs []TableRef

	switch v := n.(type) {
	case *ast.SelectStmt:
		if v == nil {
			return nil
		}
		refs = append(refs, collectFromWithClause(v.WithClause)...)
		refs = append(refs, collectFromList(v.FromClause)...)
		if v.IntoTable != nil {
			refs = append(refs, astTableRefToRef(v.IntoTable))
		}
		// Recurse into set operations.
		refs = append(refs, collectTableRefsFromNode(v.Larg)...)
		refs = append(refs, collectTableRefsFromNode(v.Rarg)...)

	case *ast.InsertStmt:
		if v == nil {
			return nil
		}
		refs = append(refs, collectFromWithClause(v.WithClause)...)
		if v.Relation != nil {
			refs = append(refs, astTableRefToRef(v.Relation))
		}
		if src, ok := v.Source.(ast.Node); ok {
			refs = append(refs, collectTableRefsFromNode(src)...)
		}

	case *ast.UpdateStmt:
		if v == nil {
			return nil
		}
		refs = append(refs, collectFromWithClause(v.WithClause)...)
		if v.Relation != nil {
			refs = append(refs, astTableRefToRef(v.Relation))
		}
		refs = append(refs, collectFromList(v.FromClause)...)

	case *ast.DeleteStmt:
		if v == nil {
			return nil
		}
		refs = append(refs, collectFromWithClause(v.WithClause)...)
		if v.Relation != nil {
			refs = append(refs, astTableRefToRef(v.Relation))
		}
		refs = append(refs, collectFromList(v.FromClause)...)

	case *ast.MergeStmt:
		if v == nil {
			return nil
		}
		refs = append(refs, collectFromWithClause(v.WithClause)...)
		if v.Target != nil {
			ref := astTableRefToRef(v.Target)
			// The parser may incorrectly assign "USING" as the target alias.
			if strings.EqualFold(ref.Alias, "USING") {
				ref.Alias = ""
			}
			refs = append(refs, ref)
		}
		refs = append(refs, collectFromTableExpr(v.Source)...)

	case *ast.TableRef:
		if v != nil && v.Object != "" {
			refs = append(refs, astTableRefToRef(v))
		}

	case *ast.JoinClause:
		if v == nil {
			return nil
		}
		refs = append(refs, collectFromTableExpr(v.Left)...)
		refs = append(refs, collectFromTableExpr(v.Right)...)

	case *ast.AliasedTableRef:
		if v == nil {
			return nil
		}
		inner := collectFromTableExpr(v.Table)
		// Apply the outer alias to the first ref if one is present.
		if v.Alias != "" && len(inner) > 0 {
			inner[0].Alias = v.Alias
		}
		refs = append(refs, inner...)

	case *ast.SubqueryExpr:
		if v != nil {
			refs = append(refs, collectTableRefsFromNode(v.Query)...)
		}

	case *ast.PivotExpr:
		if v != nil {
			refs = append(refs, collectFromTableExpr(v.Source)...)
		}

	case *ast.UnpivotExpr:
		if v != nil {
			refs = append(refs, collectFromTableExpr(v.Source)...)
		}

	case *ast.CommonTableExpr:
		if v != nil && v.Name != "" {
			refs = append(refs, TableRef{Table: v.Name})
		}

	case *ast.List:
		if v != nil {
			for _, item := range v.Items {
				refs = append(refs, collectTableRefsFromNode(item)...)
			}
		}
	}
	return refs
}

// collectFromTableExpr extracts refs from a TableExpr interface.
func collectFromTableExpr(te ast.TableExpr) []TableRef {
	if te == nil {
		return nil
	}
	// TableExpr implementations are also ast.Node.
	if n, ok := te.(ast.Node); ok {
		return collectTableRefsFromNode(n)
	}
	return nil
}

// collectFromList extracts refs from an ast.List (e.g. FROM clause items).
func collectFromList(l *ast.List) []TableRef {
	if l == nil {
		return nil
	}
	var refs []TableRef
	for _, item := range l.Items {
		refs = append(refs, collectTableRefsFromNode(item)...)
	}
	return refs
}

// collectFromWithClause extracts CTE names from a WITH clause.
func collectFromWithClause(w *ast.WithClause) []TableRef {
	if w == nil || w.CTEs == nil {
		return nil
	}
	var refs []TableRef
	for _, item := range w.CTEs.Items {
		if cte, ok := item.(*ast.CommonTableExpr); ok && cte != nil && cte.Name != "" {
			refs = append(refs, TableRef{Table: cte.Name})
		}
	}
	return refs
}

// astTableRefToRef converts an ast.TableRef to a completion TableRef.
func astTableRefToRef(r *ast.TableRef) TableRef {
	return TableRef{
		Schema: r.Schema,
		Table:  r.Object,
		Alias:  r.Alias,
	}
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
		tokens = append(tokens, tok)
	}

	var refs []TableRef
	for i := 0; i < len(tokens); i++ {
		kw := strings.ToLower(tokens[i].Str)

		// FROM table, JOIN table, INTO table, UPDATE table
		if kw == "from" || kw == "join" || kw == "into" || kw == "update" {
			ref, skip := lexerExtractTableAfter(tokens, i+1)
			if ref != nil {
				refs = append(refs, *ref)
			}
			i += skip
			continue
		}

		// MERGE [INTO] table — skip optional INTO, extract target.
		if kw == "merge" {
			j := i + 1
			if j < len(tokens) && strings.ToLower(tokens[j].Str) == "into" {
				j++
			}
			ref, skip := lexerExtractTableAfter(tokens, j)
			if ref != nil {
				refs = append(refs, *ref)
			}
			i = j + skip - 1
			continue
		}

		// Handle multi-word join keywords: INNER JOIN, LEFT JOIN, etc.
		if kw == "inner" || kw == "left" || kw == "right" || kw == "full" || kw == "cross" || kw == "outer" {
			if i+1 < len(tokens) && strings.ToLower(tokens[i+1].Str) == "join" {
				ref, skip := lexerExtractTableAfter(tokens, i+2)
				if ref != nil {
					refs = append(refs, *ref)
				}
				i += 1 + skip
				continue
			}
			// LEFT OUTER JOIN, RIGHT OUTER JOIN, FULL OUTER JOIN
			if i+2 < len(tokens) && strings.ToLower(tokens[i+1].Str) == "outer" && strings.ToLower(tokens[i+2].Str) == "join" {
				ref, skip := lexerExtractTableAfter(tokens, i+3)
				if ref != nil {
					refs = append(refs, *ref)
				}
				i += 2 + skip
				continue
			}
		}
	}
	return refs
}

// lexerExtractTableAfter extracts a [schema.]table reference starting at tokens[j].
// Returns the TableRef (or nil) and the number of tokens consumed beyond the trigger keyword.
func lexerExtractTableAfter(tokens []parser.Token, j int) (*TableRef, int) {
	if j >= len(tokens) || !isIdentLikeToken(tokens[j]) {
		return nil, 0
	}
	ref := TableRef{Table: tokens[j].Str}
	consumed := 1

	// Check for schema.table (dot-qualified name).
	if j+2 < len(tokens) && tokens[j+1].Type == '.' && isIdentLikeToken(tokens[j+2]) {
		ref.Schema = ref.Table
		ref.Table = tokens[j+2].Str
		consumed = 3
	}

	// Check for alias: AS alias, or bare identifier alias.
	k := j + consumed
	if k < len(tokens) {
		kw := strings.ToLower(tokens[k].Str)
		if kw == "as" && k+1 < len(tokens) && isIdentLikeToken(tokens[k+1]) {
			ref.Alias = tokens[k+1].Str
		} else if isIdentLikeToken(tokens[k]) && !isReservedTableTerminator(kw) {
			ref.Alias = tokens[k].Str
		}
	}

	return &ref, consumed
}

// isIdentLikeToken reports whether a token can serve as an identifier
// (either a real identifier or a non-reserved keyword used as a name).
func isIdentLikeToken(tok parser.Token) bool {
	// Real identifiers (including [bracketed] and "quoted").
	if tok.Type > 256 {
		// tokIDENT is the only > 256 literal token that acts as an identifier.
		// Keywords also have types > 256 and can serve as identifiers in T-SQL
		// when not in a reserved position. We accept all of them here; the caller
		// filters out reserved words that would terminate a table name.
		return tok.Str != ""
	}
	return false
}

// isReservedTableTerminator checks if a keyword typically terminates a table
// reference (so it shouldn't be treated as an alias).
func isReservedTableTerminator(kw string) bool {
	switch kw {
	case "on", "where", "set", "join", "inner", "left", "right", "full",
		"cross", "outer", "using", "when", "then", "values", "select",
		"from", "into", "group", "order", "having", "union", "except",
		"intersect", "for", "option", "with", "output", "merge":
		return true
	}
	return false
}
