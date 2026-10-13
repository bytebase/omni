package review

import (
	"strings"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/tidb/parser"
)

// ident writes an identifier as MySQL would need it in a statement:
// bare when it is a plain name that is not a reserved word, otherwise
// quoted in backticks, with a backtick inside doubled and a line break
// written as \n or \r, so a message stays one line.
func ident(name string) string {
	if plainIdent(name) && !parser.IsReservedKeyword(name) {
		return name
	}
	return "`" + escapeLines(strings.ReplaceAll(name, "`", "``")) + "`"
}

// plainIdent reports whether a name needs no quoting: letters, digits,
// underscores, and dollar signs, not all digits.
func plainIdent(name string) bool {
	if name == "" {
		return false
	}
	digits := true
	for i := 0; i < len(name); i++ {
		c := name[i]
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c == '_', c == '$':
			digits = false
		case c >= '0' && c <= '9':
		default:
			return false
		}
	}
	return !digits
}

// table writes a table reference as written: its database, when the
// statement names one, and its name.
func table(ref *ast.TableRef) string {
	if ref == nil {
		return ""
	}
	if ref.Schema != "" {
		return ident(ref.Schema) + "." + ident(ref.Name)
	}
	return ident(ref.Name)
}

// tables writes a list of table references, comma separated.
func tables(refs []*ast.TableRef) string {
	names := make([]string, 0, len(refs))
	for _, ref := range refs {
		names = append(names, table(ref))
	}
	return strings.Join(names, ", ")
}

// tableExprs writes the tables a list of table expressions names, the
// sides of a join included.
func tableExprs(exprs []ast.TableExpr) string {
	var names []string
	var add func(e ast.TableExpr)
	add = func(e ast.TableExpr) {
		switch v := e.(type) {
		case *ast.TableRef:
			names = append(names, table(v))
		case *ast.JoinClause:
			add(v.Left)
			add(v.Right)
		}
	}
	for _, e := range exprs {
		add(e)
	}
	return strings.Join(names, ", ")
}

// sameName reports whether two names may name the same object. MySQL
// compares column, index, and constraint names, and on most servers
// table names too, without regard to case; comparing that way never
// misses a match.
func sameName(a, b string) bool {
	return strings.EqualFold(a, b)
}

// sameTable reports whether two table references may name the same
// table: the same name, in the same database when both name one. A
// reference without a database may mean the current one, which a USE in
// the change may have moved, so it matches a reference of any database.
func sameTable(a, b *ast.TableRef) bool {
	if a == nil || b == nil || !sameName(a.Name, b.Name) {
		return false
	}
	return a.Schema == "" || b.Schema == "" || sameName(a.Schema, b.Schema)
}
