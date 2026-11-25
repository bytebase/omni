package completion

import (
	"slices"
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
	"github.com/bytebase/omni/snowflake/parser"
)

// placeholder stands in for the word at the caret while the statement is
// parsed. It is an identifier no keyword spells.
const placeholder = "omni_completion_caret"

// nameKind is what the word at the caret names.
type nameKind int

const (
	kindNone       nameKind = iota // nothing completion offers: an alias, a path step, a literal
	kindExpr                       // an expression: a column, a function, or a qualifier
	kindRelation                   // a table or view
	kindColumnName                 // a column of a target: INSERT (...), SET, USING (...)
	kindFunction                   // a function name before (
)

// nameContext is what the caret word names and what is in scope there.
type nameContext struct {
	kind nameKind
	// qualifier holds the resolved parts written before the word: t in t.c.
	qualifier []string
	scope     scope
	// targets are the relations whose columns a kindColumnName word names.
	targets []*rangeItem
	// kinds narrows a kindRelation word to these kinds of relation: a DROP
	// VIEW names a view. Any kind when empty.
	kinds []CandidateType
	// query is set when the statement is a query, where unqualified names
	// resolve through the search path.
	query bool
}

// nameContextAt parses stmt with the caret word replaced by the placeholder
// and reads what the placeholder names from where it lands in the AST.
func nameContextAt(stmt string, word caretWord, cat *catalog) *nameContext {
	path, text, ok := placeholderPath(stmt, word)
	if !ok {
		return lexicalContext(stmt[:word.start])
	}
	return contextFromPath(path, text, cat)
}

// relationWords are the keywords a table or view name follows.
var relationWords = map[string]bool{"FROM": true, "JOIN": true, "INTO": true, "UPDATE": true, "TABLE": true, "VIEW": true}

// lexicalContext reads a relation name's context from the tokens before the
// caret, for a statement no patch makes parse, such as ALTER TABLE cut short
// before its action: a table name follows FROM, JOIN, INTO, UPDATE, TABLE or
// VIEW, after any qualifier.
func lexicalContext(prefix string) *nameContext {
	toks, _ := parser.Tokenize(prefix)
	toks = toks[:len(toks)-1] // EOF
	var qualifier []string
	i := len(toks) - 1
	for i >= 1 && toks[i].Type == '.' {
		id, ok := identText(toks[i-1])
		if !ok {
			return nil
		}
		qualifier = append([]string{id}, qualifier...)
		i -= 2
	}
	if i < 0 || len(qualifier) > 2 || !isWordToken(toks[i]) || !relationWords[strings.ToUpper(toks[i].Str)] {
		return nil
	}
	return &nameContext{kind: kindRelation, qualifier: qualifier}
}

// identText resolves an identifier token: unquoted ones fold to upper case.
func identText(t parser.Token) (string, bool) {
	switch {
	case parser.TokenName(t.Type) == "QUOTED_IDENT":
		return t.Str, true
	case isWordToken(t):
		return strings.ToUpper(t.Str), true
	}
	return "", false
}

// patchTails close a statement cut short at the caret or left incomplete
// after it, tried in order until one parses: what may follow a name (an
// expression's operand, a SET's value), then each unclosed bracket, then a
// query for a statement that needs one.
var patchTails = []string{"", " 1", " = 1", " x"}

// placeholderPath returns the AST path from the statement to the node that
// holds the placeholder, innermost last. It tries the statement with its
// text after the caret, then the text before the caret alone, each closed by
// patchTails.
func placeholderPath(stmt string, word caretWord) ([]ast.Node, string, bool) {
	head := stmt[:word.start] + placeholder
	for _, base := range []string{head + stmt[word.end:], head} {
		for _, tail := range patchTails {
			text := base + tail
			closers := closersFor(text)
			for _, query := range []string{"", " SELECT 1"} {
				patched := text + closers + query
				if path, ok := parsePath(patched); ok {
					return path, patched, true
				}
			}
		}
	}
	return nil, "", false
}

// closersFor returns the brackets that close those text leaves open.
func closersFor(text string) string {
	toks, _ := parser.Tokenize(text)
	var open []byte
	for _, t := range toks {
		switch t.Type {
		case '(':
			open = append(open, ')')
		case '[':
			open = append(open, ']')
		case ')', ']':
			if n := len(open); n > 0 && open[n-1] == byte(t.Type) {
				open = open[:n-1]
			}
		}
	}
	var b strings.Builder
	for i := len(open) - 1; i >= 0; i-- {
		b.WriteByte(open[i])
	}
	return b.String()
}

// parsePath parses text and returns the path to the placeholder, when text
// parses without error and holds it.
func parsePath(text string) ([]ast.Node, bool) {
	res := parser.ParseBestEffort(text)
	if len(res.Errors) > 0 || len(res.File.Stmts) != 1 {
		return nil, false
	}
	f := &placeholderFinder{}
	ast.Walk(f, res.File.Stmts[0])
	return f.path, f.path != nil
}

// placeholderFinder walks to the node that holds the placeholder, keeping
// the path to it.
type placeholderFinder struct {
	stack []ast.Node
	path  []ast.Node
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
	if holdsPlaceholder(n) {
		f.path = append([]ast.Node(nil), f.stack...)
		return nil
	}
	return f
}

func isPlaceholder(id ast.Ident) bool {
	return strings.EqualFold(id.Name, placeholder)
}

func placeholderIndex(ids []ast.Ident) int {
	for i, id := range ids {
		if isPlaceholder(id) {
			return i
		}
	}
	return -1
}

// holdsPlaceholder reports whether n holds the placeholder in an identifier
// field of its own rather than in a child node.
func holdsPlaceholder(n ast.Node) bool {
	switch n := n.(type) {
	case *ast.ColumnRef:
		return placeholderIndex(n.Parts) >= 0
	case *ast.ObjectName:
		return placeholderIndex(n.Parts()) >= 0
	case *ast.AccessExpr:
		return isPlaceholder(n.Field)
	case *ast.TableRef:
		return isPlaceholder(n.Alias) || placeholderIndex(n.Columns) >= 0
	case *ast.InsertStmt:
		return placeholderIndex(n.Columns) >= 0
	case *ast.InsertMultiStmt:
		return slices.ContainsFunc(n.Branches, func(br *ast.InsertMultiBranch) bool { return placeholderIndex(br.Columns) >= 0 })
	case *ast.SelectStmt:
		return starListTarget(n) != nil
	case *ast.UnpivotColumn:
		return isPlaceholder(n.Column) || isPlaceholder(n.Alias)
	case *ast.JoinExpr:
		return placeholderIndex(n.Using) >= 0
	case *ast.MergeStmt:
		for _, w := range n.Whens {
			if placeholderIndex(w.InsertCols) >= 0 {
				return true
			}
		}
		return isPlaceholder(n.TargetAlias) || isPlaceholder(n.SourceAlias)
	case *ast.UpdateStmt:
		return isPlaceholder(n.Alias)
	case *ast.DeleteStmt:
		return isPlaceholder(n.Alias)
	}
	return false
}

// contextFromPath reads the name context from the path to the node holding
// the placeholder.
func contextFromPath(path []ast.Node, text string, cat *catalog) *nameContext {
	holder := path[len(path)-1]
	var parent ast.Node
	if len(path) > 1 {
		parent = path[len(path)-2]
	}
	nc := &nameContext{query: isQuery(path[0])}
	b := newScopeBuilder(cat, nc.query, text)
	ctes := ctesAlong(path)
	switch n := holder.(type) {
	case *ast.ColumnRef:
		i := placeholderIndex(n.Parts)
		nc.kind = kindExpr
		nc.qualifier = normalizeParts(n.Parts[:i])
	case *ast.ObjectName:
		parts := n.Parts()
		i := placeholderIndex(parts)
		nc.qualifier = normalizeParts(parts[:i])
		nc.kind, nc.targets, nc.kinds = objectNameKind(n, parent, b)
	case *ast.InsertStmt:
		nc.kind = kindColumnName
		nc.targets = []*rangeItem{b.relation(n.Target, ast.Ident{}, nil, 0)}
	case *ast.InsertMultiStmt:
		nc.kind = kindColumnName
		for _, br := range n.Branches {
			if placeholderIndex(br.Columns) >= 0 {
				nc.targets = []*rangeItem{b.relation(br.Target, ast.Ident{}, nil, 0)}
			}
		}
	case *ast.SelectStmt:
		// SELECT * EXCLUDE (...) or RENAME (... AS ...) names columns of
		// the star's relations.
		t := starListTarget(n)
		var qualifier []string
		if star, ok := t.Expr.(*ast.StarExpr); ok && star.Qualifier != nil {
			qualifier = normalizeParts(star.Qualifier.Parts())
		}
		nc.kind = kindColumnName
		for _, item := range b.leavesOf(n.From, ctes[len(ctes)-1]) {
			if qualifier == nil || item.addressedBy(qualifier) {
				nc.targets = append(nc.targets, item)
			}
		}
	case *ast.UnpivotColumn:
		// UNPIVOT ( v FOR n IN ( <column> [AS alias] ... ) ) names a column
		// of the relation the UNPIVOT reads.
		ref, ok := tableRefAbove(path)
		if !ok || isPlaceholder(n.Alias) {
			return nil
		}
		nc.kind = kindColumnName
		nc.targets = []*rangeItem{b.clauseSource(ref, ctes[len(ctes)-1], 0)}
	case *ast.MergeStmt:
		if isPlaceholder(n.TargetAlias) || isPlaceholder(n.SourceAlias) {
			return nil
		}
		nc.kind = kindColumnName
		nc.targets = []*rangeItem{b.relation(n.Target, n.TargetAlias, nil, 0)}
	case *ast.JoinExpr:
		nc.kind = kindColumnName
		nc.targets = b.leaves(n, ctes[len(ctes)-1])
	default:
		// An alias or a step of a semi-structured path names nothing in the
		// catalog.
		return nil
	}
	if nc.kind == kindNone {
		return nil
	}
	if nc.kind == kindExpr || nc.kind == kindFunction || nc.kind == kindRelation {
		nc.scope = b.scopeAt(path, ctes)
	}
	return nc
}

// starListTarget returns the star select-list item whose EXCLUDE list, or
// RENAME's source column, holds the placeholder.
func starListTarget(sel *ast.SelectStmt) *ast.SelectTarget {
	for _, t := range sel.Targets {
		if placeholderIndex(t.Exclude) >= 0 || slices.ContainsFunc(t.Rename, func(r ast.StarRename) bool { return isPlaceholder(r.Col) }) {
			return t
		}
	}
	return nil
}

// tableRefAbove returns the innermost table reference on path.
func tableRefAbove(path []ast.Node) (*ast.TableRef, bool) {
	for i := len(path) - 1; i >= 0; i-- {
		if ref, ok := path[i].(*ast.TableRef); ok {
			return ref, true
		}
	}
	return nil, false
}

// isQuery reports whether stmt is a query rather than DDL or DML.
func isQuery(stmt ast.Node) bool {
	switch stmt.(type) {
	case *ast.SelectStmt, *ast.SetOperationStmt, *ast.ResultScanStmt:
		return true
	}
	return false
}

// objectNameKind returns what an object name holding the placeholder names,
// from the node that holds the name: its kind, the targets of a column name,
// and the relation kinds a relation name is narrowed to.
func objectNameKind(n *ast.ObjectName, parent ast.Node, b *scopeBuilder) (nameKind, []*rangeItem, []CandidateType) {
	tables := []CandidateType{CandidateTable}
	switch p := parent.(type) {
	case *ast.TableRef:
		if p.Name == n {
			return kindRelation, nil, nil
		}
	case *ast.FuncCallExpr:
		return kindFunction, nil, nil
	case *ast.InsertStmt, *ast.InsertMultiStmt, *ast.DeleteStmt:
		return kindRelation, nil, nil
	case *ast.AlterTableStmt:
		return kindRelation, nil, tables
	case *ast.UndropStmt:
		if p.Kind == ast.UndropTable || p.Kind == ast.UndropDynamicTable {
			return kindRelation, nil, tables
		}
	case *ast.TruncateStmt:
		if p.MaterializedView {
			return kindRelation, nil, []CandidateType{CandidateMaterializedView}
		}
		return kindRelation, nil, tables
	case *ast.UpdateStmt:
		if p.Target == n {
			return kindRelation, nil, nil
		}
		// SET <column> = ...
		return kindColumnName, []*rangeItem{b.relation(p.Target, p.Alias, nil, 0)}, nil
	case *ast.MergeStmt:
		if p.Target == n {
			return kindRelation, nil, nil
		}
		return kindColumnName, []*rangeItem{b.relation(p.Target, p.TargetAlias, nil, 0)}, nil
	case *ast.DropStmt:
		switch p.Kind {
		case ast.DropTable, ast.DropDynamicTable, ast.DropExternalTable:
			return kindRelation, nil, tables
		case ast.DropView:
			return kindRelation, nil, []CandidateType{CandidateView}
		case ast.DropMaterializedView:
			return kindRelation, nil, []CandidateType{CandidateMaterializedView}
		}
	case *ast.DescribeStmt:
		switch p.ObjectType {
		case "TABLE", "VIEW", "MATERIALIZED VIEW", "EXTERNAL TABLE", "DYNAMIC TABLE":
			return kindRelation, nil, nil
		}
	case *ast.CommentStmt:
		switch p.ObjectType {
		case "TABLE", "VIEW", "MATERIALIZED VIEW", "EXTERNAL TABLE", "DYNAMIC TABLE":
			return kindRelation, nil, nil
		}
	}
	return kindNone, nil, nil
}

// normalizeParts resolves identifiers as Snowflake does: unquoted ones fold
// to upper case.
func normalizeParts(ids []ast.Ident) []string {
	out := make([]string, len(ids))
	for i, id := range ids {
		out[i] = id.Normalize()
	}
	return out
}
