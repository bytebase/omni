package review

import (
	"strings"

	"github.com/bytebase/omni/mssql/ast"
)

// ident writes an identifier as SQL Server would need it in a statement:
// bare when it is a regular identifier, otherwise in brackets, with a
// closing bracket inside doubled and a line break written as \n or \r, so
// a message stays one line.
func ident(name string) string {
	if regularIdent(name) {
		return name
	}
	return "[" + escapeLines(strings.ReplaceAll(name, "]", "]]")) + "]"
}

// regularIdent reports whether a name needs no delimiters: a letter, an
// underscore, @, or #, then letters, digits, @, $, #, or underscores.
func regularIdent(name string) bool {
	if name == "" {
		return false
	}
	for i := 0; i < len(name); i++ {
		c := name[i]
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c == '_', c == '@', c == '#':
		case (c >= '0' && c <= '9' || c == '$') && i > 0:
		default:
			return false
		}
	}
	return true
}

// object writes a table reference as written: its server, database, and
// schema when the statement names them, and its name.
func object(ref *ast.TableRef) string {
	if ref == nil {
		return ""
	}
	// A part left out below a named one keeps its place: shop..t.
	var parts []string
	for _, p := range []string{ref.Server, ref.Database, ref.Schema} {
		switch {
		case p != "":
			parts = append(parts, ident(p))
		case len(parts) > 0:
			parts = append(parts, "")
		}
	}
	return strings.Join(append(parts, ident(ref.Object)), ".")
}

// tableExpr writes the table an UPDATE or DELETE names as its target.
func tableExpr(e ast.TableExpr) string {
	if ref, ok := e.(*ast.TableRef); ok {
		return object(ref)
	}
	return "the target"
}

// sameName reports whether two names may name the same object. SQL Server
// compares identifiers under the database's collation, which is usually
// case-insensitive; comparing without case never misses a match.
func sameName(a, b string) bool {
	return strings.EqualFold(a, b)
}

// sameObject reports whether two references may name the same object:
// the same name, in the same schema and database when both name one. A
// reference without a schema may mean the default schema, which another
// reference names, so it matches any.
func sameObject(a, b *ast.TableRef) bool {
	if a == nil || b == nil || !sameName(a.Object, b.Object) {
		return false
	}
	for _, pair := range [][2]string{{a.Schema, b.Schema}, {a.Database, b.Database}, {a.Server, b.Server}} {
		if pair[0] != "" && pair[1] != "" && !sameName(pair[0], pair[1]) {
			return false
		}
	}
	return true
}
