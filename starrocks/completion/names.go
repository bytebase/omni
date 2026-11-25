package completion

import (
	"slices"
	"strings"

	"github.com/bytebase/omni/starrocks/ast"
	"github.com/bytebase/omni/starrocks/parser"
)

// placeholder is the identifier the caret word is replaced with to find
// where it lands in the AST.
const placeholder = "omni_caret_name"

// nameKind is what a name at the caret names.
type nameKind int

const (
	kindNone     nameKind = iota
	kindRelation          // a table, view or materialized view
	kindColumn            // a column of a table in scope
)

// nameContext is what the caret names and what is in scope there.
type nameContext struct {
	kind      nameKind
	qualifier []string    // the parts written before the caret word, as in db.
	scope     []tableItem // the query's FROM items and DML target
	outer     []tableItem // the FROM items of the statement around a subquery
	ctes      []string    // the names of the CTEs, the query's and the outer statement's
	aliases   []string    // the select list's aliases, where the caret's clause sees them
}

// tableItem is a table the query reads: its name parts and its alias.
type tableItem struct {
	parts []string
	alias string
}

// placeholderTails close a statement cut at the caret, tried in order when
// it does not parse with the rest of the text.
var placeholderTails = []string{"", " ON TRUE", " SET a = 1", " VALUES (1)", ")", " = 1", " = 1)"}

// tableKeywords are the keywords after which a statement names a table it
// reads or writes.
var tableKeywords = map[string]bool{
	"FROM": true, "JOIN": true, "INTO": true, "UPDATE": true, "TABLE": true,
	"DESC": true, "DESCRIBE": true,
}

// nameContextAt returns what the caret word in stmt names. outer, when not
// empty, is the statement around stmt, a subquery (see innermostQuery).
func nameContextAt(stmt string, word caretWord, outer string) *nameContext {
	nc := &nameContext{}
	var qualStart int
	nc.qualifier, qualStart = qualifierBefore(stmt, word.start)

	prefix := stmt[:word.start] + placeholder
	var root ast.Node
	rest := stmt[word.end:]
	truncated := false
	for i, tail := range append([]string{rest}, placeholderTails...) {
		if i > 0 && tail == "" && rest == "" {
			continue // the same text as the first try
		}
		if file, err := parser.Parse(prefix + tail); err == nil && len(file.Stmts) == 1 {
			root = file.Stmts[0]
			truncated = i > 0 && strings.TrimSpace(rest) != ""
			break
		}
	}
	if root != nil {
		var block ast.Node
		var path []ast.Node
		nc.kind, block, path = placeholderAt(root)
		collectScope(root, block, path, nc)
		nc.aliases = visibleAliases(block, path)
		if truncated {
			// The parse dropped the text after the caret, and with it the
			// FROM items a select list reads.
			nc.scope = append(nc.scope, lexicalScope(rest)...)
		}
	} else {
		nc.scope = lexicalScope(stmt[:qualStart])
		if expressionStart(stmt[:qualStart]) {
			nc.kind = kindColumn
		}
	}
	if outer != "" {
		if file, err := parser.Parse(outer); err == nil && len(file.Stmts) == 1 {
			outerScope := &nameContext{}
			collectScope(file.Stmts[0], nil, nil, outerScope)
			nc.outer = outerScope.scope
			nc.ctes = append(nc.ctes, outerScope.ctes...)
		} else {
			nc.outer = lexicalScope(outer)
		}
	}
	// CREATE TABLE names a new table; its AS query is raw text in the AST.
	if prev := previousKeyword(stmt[:qualStart]); nc.kind == kindNone && tableKeywords[prev] && (prev != "TABLE" || !startsWith(stmt, "CREATE")) {
		nc.kind = kindRelation
	}
	return nc
}

// placeholderAt finds the placeholder in root: a FROM item's name is a
// relation, a column reference or an UPDATE assignment's target a column.
// block is the query block holding it, whose FROM items are in scope: the
// innermost SELECT, UPDATE, DELETE or INSERT around it, or root. path is the
// nodes from root to the placeholder.
func placeholderAt(root ast.Node) (kind nameKind, block ast.Node, path []ast.Node) {
	f := &placeholderFinder{}
	ast.Walk(f, root)
	block = root
	for _, n := range f.path {
		switch n.(type) {
		case *ast.SelectStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.InsertStmt:
			block = n
		}
	}
	return f.kind, block, f.path
}

// placeholderFinder walks to the node holding the placeholder, keeping the
// path of nodes from the root to it.
type placeholderFinder struct {
	stack []ast.Node
	path  []ast.Node
	kind  nameKind
}

func (f *placeholderFinder) Visit(n ast.Node) ast.Visitor {
	if f.path != nil {
		return nil
	}
	if n == nil {
		f.stack = f.stack[:len(f.stack)-1]
		return nil
	}
	f.stack = append(f.stack, n)
	switch n := n.(type) {
	case *ast.TableRef:
		if n.Subquery == nil && endsWithPlaceholder(n.Name) {
			f.kind = kindRelation
		}
	case *ast.ColumnRef:
		if endsWithPlaceholder(n.Name) {
			f.kind = kindColumn
		}
	case *ast.Assignment:
		if endsWithPlaceholder(n.Column) {
			f.kind = kindColumn
		}
	case *ast.InsertStmt:
		if slices.Contains(n.Columns, placeholder) {
			f.kind = kindColumn
		}
	case *ast.JoinClause:
		if slices.Contains(n.Using, placeholder) {
			f.kind = kindColumn
		}
	}
	// A DML statement's target, in a form the relation keywords miss
	// (StarRocks's INSERT OVERWRITE t).
	switch n := n.(type) {
	case *ast.InsertStmt:
		if endsWithPlaceholder(n.Target) {
			f.kind = kindRelation
		}
	case *ast.UpdateStmt:
		if endsWithPlaceholder(n.Target) {
			f.kind = kindRelation
		}
	case *ast.DeleteStmt:
		if endsWithPlaceholder(n.Target) {
			f.kind = kindRelation
		}
	}
	if f.kind != kindNone {
		f.path = slices.Clone(f.stack)
		return nil
	}
	return f
}

func endsWithPlaceholder(name *ast.ObjectName) bool {
	return name != nil && len(name.Parts) > 0 && name.Parts[len(name.Parts)-1] == placeholder
}

// collectScope records the tables block reads and the CTE names visible
// on path, the nodes from root to the caret. A CTE's query and another arm
// of a set operation are other blocks, and a subquery is raw text in the
// AST, so only the block's own FROM items and DML target are collected. A
// nil block collects every block's tables, and a nil path the CTEs of
// root's own WITH.
func collectScope(root, block ast.Node, path []ast.Node, nc *nameContext) {
	ast.Inspect(root, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.CTE:
			nc.ctes = append(nc.ctes, n.Name)
			return slices.Contains(path, ast.Node(n))
		case *ast.SelectStmt, *ast.SetOpStmt:
			return n == root || slices.Contains(path, n)
		}
		return true
	})
	add := func(name *ast.ObjectName, alias string) {
		if name != nil && !endsWithPlaceholder(name) {
			nc.scope = append(nc.scope, tableItem{parts: name.Parts, alias: alias})
		}
	}
	all := block == nil
	if all {
		block = root
	}
	ast.Inspect(block, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.CTE:
			return all
		case *ast.SelectStmt, *ast.SetOpStmt:
			if n != block && !all {
				return false
			}
		case *ast.InsertStmt:
			// The target's columns are in scope only for its column list; a
			// query block below it is a block of its own.
			add(n.Target, "")
			return all
		case *ast.TableRef:
			if n.Subquery == nil {
				add(n.Name, n.Alias)
			}
		case *ast.UpdateStmt:
			add(n.Target, n.TargetAlias)
		case *ast.DeleteStmt:
			add(n.Target, n.TargetAlias)
		}
		return true
	})
}

// visibleAliases returns the select list's aliases when the caret is in
// block's GROUP BY, HAVING or ORDER BY, the clauses that read them.
func visibleAliases(block ast.Node, path []ast.Node) []string {
	sel, ok := block.(*ast.SelectStmt)
	i := slices.Index(path, block)
	if !ok || i < 0 || i+1 >= len(path) {
		return nil
	}
	child := path[i+1]
	visible := sel.Having != nil && child == sel.Having || slices.Contains(sel.GroupBy, child)
	for _, o := range sel.OrderBy {
		visible = visible || ast.Node(o) == child
	}
	if !visible {
		return nil
	}
	var out []string
	for _, item := range sel.Items {
		if item.Alias != "" {
			out = append(out, item.Alias)
		}
	}
	return out
}

// expressionKeywords are the keywords an expression, and so a column name,
// can follow.
var expressionKeywords = map[string]bool{
	"SELECT": true, "DISTINCT": true, "WHERE": true, "ON": true, "AND": true, "OR": true,
	"NOT": true, "BY": true, "HAVING": true, "SET": true, "CASE": true, "WHEN": true,
	"THEN": true, "ELSE": true, "IN": true, "BETWEEN": true, "LIKE": true, "IS": true,
}

// expressionStart reports whether a column name can start right after
// text, read from its last token: an operator, '(' or ',', or a keyword an
// expression follows. It places the caret of a statement that does not
// parse.
func expressionStart(text string) bool {
	toks, _ := parser.Tokenize(text)
	if len(toks) < 2 {
		return false
	}
	last := toks[len(toks)-2]
	if isKeyword(last) {
		return expressionKeywords[strings.ToUpper(last.Str)]
	}
	switch parser.TokenName(last.Kind) {
	case "(", ",", "=", "<", ">", "+", "-", "*", "/", "%", "<=", ">=", "<>", "<=>", "&&", "||":
		return true
	}
	return false
}

// qualifierBefore returns the dotted parts written right before offset end,
// as in db.t., and where they start.
func qualifierBefore(stmt string, end int) ([]string, int) {
	toks, _ := parser.Tokenize(stmt[:end])
	toks = toks[:len(toks)-1] // EOF
	var parts []string
	start := end
	for i := len(toks) - 1; i >= 1; i -= 2 {
		dot, part := toks[i], toks[i-1]
		// Tokens, not bytes, are adjacent: the parser reads db . t too.
		if dot.Kind != '.' || !isName(part) {
			break
		}
		parts = append([]string{part.Str}, parts...)
		start = part.Loc.Start
	}
	return parts, start
}

// isName reports whether t can be a name: an identifier, quoted or not, or
// a keyword word.
func isName(t parser.Token) bool {
	name := parser.TokenName(t.Kind)
	return name == "IDENT" || name == "QUOTED_IDENT" || isKeyword(t)
}

// previousKeyword returns the last token of text upper-cased when it is a
// keyword, and "" otherwise. An IF [NOT] EXISTS before the caret is passed
// over, so DROP TABLE IF EXISTS reads as DROP TABLE.
func previousKeyword(text string) string {
	toks, _ := parser.Tokenize(text)
	toks = toks[:len(toks)-1] // EOF
	word := func(i int, w string) bool {
		return i >= 0 && isKeyword(toks[i]) && strings.EqualFold(toks[i].Str, w)
	}
	n := len(toks)
	if word(n-1, "EXISTS") {
		n--
		if word(n-1, "NOT") {
			n--
		}
		if word(n-1, "IF") {
			n--
		}
	}
	if n < 1 || !isKeyword(toks[n-1]) {
		return ""
	}
	return strings.ToUpper(toks[n-1].Str)
}

func startsWith(stmt, keyword string) bool {
	toks, _ := parser.Tokenize(stmt)
	return isKeyword(toks[0]) && strings.EqualFold(toks[0].Str, keyword)
}

// lexicalScope reads the tables named after FROM, JOIN, UPDATE and INTO in
// text, a statement that does not parse, with their aliases.
func lexicalScope(text string) []tableItem {
	toks, _ := parser.Tokenize(text)
	var out []tableItem
	for i := 0; i < len(toks); i++ {
		if !isKeyword(toks[i]) || !tableKeywords[strings.ToUpper(toks[i].Str)] {
			continue
		}
		// A FROM list goes on past each comma.
		for j := i + 1; ; j++ {
			item, next, ok := lexicalItem(toks, j)
			if !ok {
				break
			}
			out = append(out, item)
			i = next - 1
			if next >= len(toks) || toks[next].Kind != ',' {
				break
			}
			j = next
		}
	}
	return out
}

// lexicalItem reads a table name and its alias at toks[j], and returns the
// index of the token after them.
func lexicalItem(toks []parser.Token, j int) (tableItem, int, bool) {
	var parts []string
	for j < len(toks) && isName(toks[j]) && !parser.IsReserved(toks[j].Kind) {
		parts = append(parts, toks[j].Str)
		if j+2 < len(toks) && toks[j+1].Kind == '.' {
			j += 2
			continue
		}
		j++
		break
	}
	if len(parts) == 0 {
		return tableItem{}, j, false
	}
	item := tableItem{parts: parts}
	k := j
	if k < len(toks) && isKeyword(toks[k]) && strings.EqualFold(toks[k].Str, "AS") {
		k++
	}
	if k < len(toks) && isName(toks[k]) && !parser.IsReserved(toks[k].Kind) {
		item.alias = toks[k].Str
		j = k + 1
	}
	return item, j, true
}

// nameCandidates returns the names nc offers, unquoted.
func nameCandidates(nc *nameContext, cat *catalog) []Candidate {
	var out []Candidate
	switch nc.kind {
	case kindRelation:
		switch len(nc.qualifier) {
		case 0:
			for _, cte := range nc.ctes {
				out = append(out, Candidate{Text: cte, Type: CandidateTable})
			}
			out = append(out, cat.relations(cat.currentDatabase())...)
			out = append(out, cat.databaseNames()...)
		case 1:
			out = append(out, cat.relations(cat.database(nc.qualifier[0]))...)
		case 2:
			if strings.EqualFold(nc.qualifier[0], nativeCatalog) {
				out = append(out, cat.relations(cat.database(nc.qualifier[1]))...)
			}
		}
	case kindColumn:
		switch len(nc.qualifier) {
		case 0:
			for _, item := range nc.scope {
				out = append(out, cat.columns(item.parts, nc.ctes)...)
			}
			for _, alias := range nc.aliases {
				out = append(out, Candidate{Text: alias, Type: CandidateColumn})
			}
		case 1:
			// t. is a table or alias in scope, else one of the statement
			// around a subquery, or else a database.
			hit := false
			for _, scope := range [][]tableItem{nc.scope, nc.outer} {
				for _, item := range scope {
					if item.addressedBy(nc.qualifier[0]) {
						hit = true
						out = append(out, cat.columns(item.parts, nc.ctes)...)
					}
				}
				if hit {
					break
				}
			}
			if !hit {
				out = append(out, cat.relations(cat.database(nc.qualifier[0]))...)
			}
		case 2, 3:
			// db.t. or catalog.db.t.
			out = append(out, cat.columns(nc.qualifier, nil)...)
		}
	}
	return out
}

// addressedBy reports whether a qualifier names the item: its alias, or
// with no alias its table name.
func (t tableItem) addressedBy(q string) bool {
	if t.alias != "" {
		return t.alias == q
	}
	return len(t.parts) > 0 && t.parts[len(t.parts)-1] == q
}
