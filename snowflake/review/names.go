package review

import (
	"github.com/bytebase/omni/snowflake/ast"
)

// ident writes an identifier as the SQL wrote it: bare when unquoted, in
// double quotes with a double quote inside doubled when quoted, and with a
// line break written as \n or \r, so a message stays one line.
func ident(id ast.Ident) string {
	return escapeLines(id.String())
}

// object writes a qualified name as the SQL wrote it.
func object(n *ast.ObjectName) string {
	if n == nil {
		return ""
	}
	return escapeLines(n.String())
}

// sameName reports whether two identifiers name the same object. Snowflake
// stores and resolves an unquoted identifier in upper case and a quoted one
// exactly as written, so t, T, and "T" are one name and "t" another.
func sameName(a, b ast.Ident) bool {
	return a.Normalize() == b.Normalize()
}

// sameObject reports whether two names may name the same object: the same
// name, in the same schema and database when both name one. A name without
// a schema may mean the current schema, which the other names, so it
// matches any. The rules use it where a match only holds a judgment back.
func sameObject(a, b *ast.ObjectName) bool {
	if a == nil || b == nil || !sameName(a.Name, b.Name) {
		return false
	}
	for _, pair := range [][2]ast.Ident{{a.Schema, b.Schema}, {a.Database, b.Database}} {
		if !pair[0].IsEmpty() && !pair[1].IsEmpty() && !sameName(pair[0], pair[1]) {
			return false
		}
	}
	return true
}

// sameWritten reports whether two names are written alike: the same parts,
// each the same name. Two names written alike in statements with no USE
// between them name the same object. The rules use it where a match makes
// a finding.
func sameWritten(a, b *ast.ObjectName) bool {
	if a == nil || b == nil {
		return false
	}
	for _, pair := range [][2]ast.Ident{{a.Name, b.Name}, {a.Schema, b.Schema}, {a.Database, b.Database}} {
		if pair[0].IsEmpty() != pair[1].IsEmpty() || !sameName(pair[0], pair[1]) {
			return false
		}
	}
	return true
}
