package completion

import (
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
	ctes      []string    // the names of the query's CTEs
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

// nameContextAt returns what the caret word in stmt names.
func nameContextAt(stmt string, word caretWord) *nameContext {
	nc := &nameContext{}
	var qualStart int
	nc.qualifier, qualStart = qualifierBefore(stmt, word.start)

	prefix := stmt[:word.start] + placeholder
	var root ast.Node
	for i, tail := range append([]string{stmt[word.end:]}, placeholderTails...) {
		if i > 0 && tail == "" && word.end == len(stmt) {
			continue // the same text as the first try
		}
		if file, err := parser.Parse(prefix + tail); err == nil && len(file.Stmts) == 1 {
			root = file.Stmts[0]
			break
		}
	}
	if root != nil {
		nc.kind = placeholderKind(root)
		collectScope(root, nc)
	} else {
		nc.scope = lexicalScope(stmt[:qualStart])
	}
	// CREATE TABLE names a new table; its AS query is raw text in the AST.
	if prev := previousKeyword(stmt[:qualStart]); nc.kind == kindNone && tableKeywords[prev] && (prev != "TABLE" || !startsWith(stmt, "CREATE")) {
		nc.kind = kindRelation
	}
	return nc
}

// placeholderKind finds the placeholder in root: a FROM item's name is a
// relation, a column reference or an UPDATE assignment's target a column.
func placeholderKind(root ast.Node) nameKind {
	kind := kindNone
	ast.Inspect(root, func(n ast.Node) bool {
		if kind != kindNone {
			return false
		}
		switch n := n.(type) {
		case *ast.TableRef:
			if n.Subquery == nil && endsWithPlaceholder(n.Name) {
				kind = kindRelation
			}
		case *ast.ColumnRef:
			if endsWithPlaceholder(n.Name) {
				kind = kindColumn
			}
		case *ast.Assignment:
			if endsWithPlaceholder(n.Column) {
				kind = kindColumn
			}
		}
		return true
	})
	return kind
}

func endsWithPlaceholder(name *ast.ObjectName) bool {
	return name != nil && len(name.Parts) > 0 && name.Parts[len(name.Parts)-1] == placeholder
}

// collectScope records the tables and CTEs root reads. A subquery is raw
// text in the AST, so only the query holding the caret is collected.
func collectScope(root ast.Node, nc *nameContext) {
	add := func(name *ast.ObjectName, alias string) {
		if name != nil && !endsWithPlaceholder(name) {
			nc.scope = append(nc.scope, tableItem{parts: name.Parts, alias: alias})
		}
	}
	ast.Inspect(root, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.TableRef:
			if n.Subquery == nil {
				add(n.Name, n.Alias)
			}
		case *ast.UpdateStmt:
			add(n.Target, n.TargetAlias)
		case *ast.DeleteStmt:
			add(n.Target, n.TargetAlias)
		case *ast.CTE:
			nc.ctes = append(nc.ctes, n.Name)
		}
		return true
	})
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
		if dot.Kind != '.' || dot.Loc.End != start || !isName(part) || part.Loc.End != dot.Loc.Start {
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
// keyword, and "" otherwise.
func previousKeyword(text string) string {
	toks, _ := parser.Tokenize(text)
	if len(toks) < 2 || !isKeyword(toks[len(toks)-2]) {
		return ""
	}
	return strings.ToUpper(toks[len(toks)-2].Str)
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
		j := i + 1
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
			continue
		}
		item := tableItem{parts: parts}
		if j < len(toks) && isKeyword(toks[j]) && strings.EqualFold(toks[j].Str, "AS") {
			j++
		}
		if j < len(toks) && isName(toks[j]) && !parser.IsReserved(toks[j].Kind) {
			item.alias = toks[j].Str
		}
		out = append(out, item)
		i = j
	}
	return out
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
		}
	case kindColumn:
		switch len(nc.qualifier) {
		case 0:
			for _, item := range nc.scope {
				out = append(out, cat.columns(item.parts, nc.ctes)...)
			}
		case 1:
			// t. is a table or alias in scope, or else a database.
			hit := false
			for _, item := range nc.scope {
				if item.addressedBy(nc.qualifier[0]) {
					hit = true
					out = append(out, cat.columns(item.parts, nc.ctes)...)
				}
			}
			if !hit {
				out = append(out, cat.relations(cat.database(nc.qualifier[0]))...)
			}
		case 2:
			// db.t.
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
