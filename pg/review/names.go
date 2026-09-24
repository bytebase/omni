package review

import (
	"strings"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
)

// ident writes an identifier the way quote_identifier does: bare when it
// is a plain lower-case name that is not a keyword needing quotes,
// double-quoted otherwise. The parser folds unquoted names, so this is
// the closest the message gets to the SQL's spelling.
//
// pg: src/backend/utils/adt/ruleutils.c — quote_identifier
func ident(name string) string {
	if name == "" {
		return `""`
	}
	plain := true
	for i := 0; i < len(name); i++ {
		c := name[i]
		switch {
		case c >= 'a' && c <= 'z', c == '_':
		case c >= '0' && c <= '9' && i > 0:
		default:
			plain = false
		}
	}
	if plain {
		if kw := parser.LookupKeyword(name); kw == nil || kw.Category == parser.UnreservedKeyword {
			return name
		}
	}
	return `"` + escapeLines(strings.ReplaceAll(name, `"`, `""`)) + `"`
}

// escapeLines writes a line break inside quoted text as \n or \r, so a
// message stays one line whatever a name or a string holds.
func escapeLines(s string) string {
	return strings.NewReplacer("\n", `\n`, "\r", `\r`).Replace(s)
}

// qualified joins name parts with dots.
func qualified(parts []string) string {
	quoted := make([]string, len(parts))
	for i, p := range parts {
		quoted[i] = ident(p)
	}
	return strings.Join(quoted, ".")
}

// relation names a RangeVar as written: schema-qualified when the SQL was.
func relation(rv *ast.RangeVar) string {
	if rv == nil {
		return ""
	}
	parts := make([]string, 0, 3)
	if rv.Catalogname != "" {
		parts = append(parts, rv.Catalogname)
	}
	if rv.Schemaname != "" {
		parts = append(parts, rv.Schemaname)
	}
	parts = append(parts, rv.Relname)
	return qualified(parts)
}

// strings returns the String items of a name list.
func nameParts(list *ast.List) []string {
	if list == nil {
		return nil
	}
	parts := make([]string, 0, len(list.Items))
	for _, item := range list.Items {
		if s, ok := item.(*ast.String); ok {
			parts = append(parts, s.Str)
		}
	}
	return parts
}

// objectName writes a DROP object as the SQL named it. Each object kind
// has its own list shape: a cast is two types, a transform a type and a
// language, an operator class or family an access method then a name,
// a trigger, policy, or rule a relation then a name, and a routine or
// operator a name with argument types.
func objectName(kind ast.ObjectType, obj ast.Node) string {
	switch kind {
	case ast.OBJECT_CAST:
		types := typeNames(listOf(obj))
		if len(types) == 2 {
			return "from " + types[0] + " to " + types[1]
		}
	case ast.OBJECT_TRANSFORM:
		if l := listOf(obj); l != nil && len(l.Items) == 2 {
			tn, _ := l.Items[0].(*ast.TypeName)
			lang, _ := l.Items[1].(*ast.String)
			if tn != nil && lang != nil {
				return "for " + typeName(tn) + " language " + ident(lang.Str)
			}
		}
	case ast.OBJECT_OPCLASS, ast.OBJECT_OPFAMILY:
		if parts := nameParts(listOf(obj)); len(parts) > 1 {
			return qualified(parts[1:]) + " using " + ident(parts[0])
		}
	case ast.OBJECT_TRIGGER, ast.OBJECT_POLICY, ast.OBJECT_RULE:
		if parts := nameParts(listOf(obj)); len(parts) > 1 {
			return ident(parts[len(parts)-1]) + " on " + qualified(parts[:len(parts)-1])
		}
	}
	switch v := obj.(type) {
	case *ast.String:
		return ident(v.Str)
	case *ast.List:
		return qualified(nameParts(v))
	case *ast.TypeName:
		return typeName(v)
	case *ast.ObjectWithArgs:
		name := qualified(nameParts(v.Objname))
		switch {
		case v.ArgsUnspecified:
			return name
		case v.Objargs == nil && kind == ast.OBJECT_AGGREGATE:
			// An aggregate with no argument list is written (*); a
			// routine with none is written ().
			return name + "(*)"
		default:
			return name + "(" + strings.Join(typeNames(v.Objargs), ", ") + ")"
		}
	}
	return ""
}

// collapseSpace puts SQL text on one line: each run of whitespace outside
// a quoted identifier, string, or dollar-quoted string becomes one space,
// and leading and trailing whitespace goes. Whitespace inside quotes is
// part of the name or value and stays, except that a line break there is
// written as \n or \r so the result is still one line.
func collapseSpace(s string) string {
	var b strings.Builder
	pending := false
	writeQuoted := func(text string) {
		for i := 0; i < len(text); i++ {
			switch text[i] {
			case '\n':
				b.WriteString(`\n`)
			case '\r':
				b.WriteString(`\r`)
			default:
				b.WriteByte(text[i])
			}
		}
	}
	for i := 0; i < len(s); {
		c := s[i]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			pending = b.Len() > 0
			i++
			continue
		}
		if pending {
			b.WriteByte(' ')
			pending = false
		}
		end := i + 1
		switch {
		case c == '"' || c == '\'':
			// An escape string, E'...', hides a quote behind a backslash;
			// a doubled quote closes and reopens, which reads the same.
			escape := c == '\'' && i > 0 && (s[i-1] == 'E' || s[i-1] == 'e')
			end = quotedEnd(s, i, escape)
		case c == '$':
			if tag := dollarTag(s[i:]); tag != "" {
				if j := strings.Index(s[i+len(tag):], tag); j >= 0 {
					end = i + len(tag) + j + len(tag)
				} else {
					end = len(s)
				}
			}
		}
		writeQuoted(s[i:end])
		i = end
	}
	return b.String()
}

// quotedEnd returns the offset just past the quoted text that starts at
// s[i], or len(s) when it is not closed.
func quotedEnd(s string, i int, escape bool) int {
	quote := s[i]
	for j := i + 1; j < len(s); j++ {
		switch {
		case escape && s[j] == '\\':
			j++
		case s[j] == quote:
			return j + 1
		}
	}
	return len(s)
}

// dollarTag returns the $tag$ delimiter that s starts with, or "" when s
// does not start one: the tag is empty or an identifier that does not
// begin with a digit, and $1 is a parameter, not a delimiter.
func dollarTag(s string) string {
	for i := 1; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '$':
			return s[:i+1]
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c == '_', c >= 0x80:
		case c >= '0' && c <= '9' && i > 1:
		default:
			return ""
		}
	}
	return ""
}

func listOf(n ast.Node) *ast.List {
	l, _ := n.(*ast.List)
	return l
}

// typeNames writes the types of an argument list; a missing type, the
// NONE side of a unary operator, reads NONE.
func typeNames(args *ast.List) []string {
	if args == nil {
		return nil
	}
	names := make([]string, 0, len(args.Items))
	for _, item := range args.Items {
		if tn, ok := item.(*ast.TypeName); ok {
			names = append(names, typeName(tn))
		} else {
			names = append(names, "NONE")
		}
	}
	return names
}

// typeName writes a type name without the pg_catalog qualifier the parser
// adds to built-in types (int becomes pg_catalog.int4).
func typeName(tn *ast.TypeName) string {
	parts := nameParts(tn.Names)
	if len(parts) == 2 && parts[0] == "pg_catalog" {
		parts = parts[1:]
	}
	name := qualified(parts)
	if tn.ArrayBounds != nil {
		name += strings.Repeat("[]", len(tn.ArrayBounds.Items))
	}
	return name
}

// objectKind names an ObjectType the way the DROP keyword does.
func objectKind(t ast.ObjectType) string {
	switch t {
	case ast.OBJECT_TABLE:
		return "table"
	case ast.OBJECT_COLUMN:
		return "column"
	case ast.OBJECT_VIEW:
		return "view"
	case ast.OBJECT_MATVIEW:
		return "materialized view"
	case ast.OBJECT_INDEX:
		return "index"
	case ast.OBJECT_SEQUENCE:
		return "sequence"
	case ast.OBJECT_SCHEMA:
		return "schema"
	case ast.OBJECT_DATABASE:
		return "database"
	case ast.OBJECT_TYPE:
		return "type"
	case ast.OBJECT_DOMAIN:
		return "domain"
	case ast.OBJECT_FUNCTION:
		return "function"
	case ast.OBJECT_PROCEDURE:
		return "procedure"
	case ast.OBJECT_ROUTINE:
		return "routine"
	case ast.OBJECT_AGGREGATE:
		return "aggregate"
	case ast.OBJECT_OPERATOR:
		return "operator"
	case ast.OBJECT_OPCLASS:
		return "operator class"
	case ast.OBJECT_OPFAMILY:
		return "operator family"
	case ast.OBJECT_TRIGGER:
		return "trigger"
	case ast.OBJECT_EVENT_TRIGGER:
		return "event trigger"
	case ast.OBJECT_POLICY:
		return "policy"
	case ast.OBJECT_RULE:
		return "rule"
	case ast.OBJECT_EXTENSION:
		return "extension"
	case ast.OBJECT_FOREIGN_TABLE:
		return "foreign table"
	case ast.OBJECT_FOREIGN_SERVER:
		return "server"
	case ast.OBJECT_FDW:
		return "foreign data wrapper"
	case ast.OBJECT_COLLATION:
		return "collation"
	case ast.OBJECT_CONVERSION:
		return "conversion"
	case ast.OBJECT_CAST:
		return "cast"
	case ast.OBJECT_LANGUAGE:
		return "language"
	case ast.OBJECT_STATISTIC_EXT:
		return "statistics"
	case ast.OBJECT_PUBLICATION:
		return "publication"
	case ast.OBJECT_SUBSCRIPTION:
		return "subscription"
	case ast.OBJECT_TSCONFIGURATION:
		return "text search configuration"
	case ast.OBJECT_TSDICTIONARY:
		return "text search dictionary"
	case ast.OBJECT_TSPARSER:
		return "text search parser"
	case ast.OBJECT_TSTEMPLATE:
		return "text search template"
	case ast.OBJECT_ACCESS_METHOD:
		return "access method"
	case ast.OBJECT_TRANSFORM:
		return "transform"
	case ast.OBJECT_ROLE:
		return "role"
	case ast.OBJECT_TABLESPACE:
		return "tablespace"
	default:
		return "object"
	}
}

// roleNames lists the roles of a RoleSpec list; special roles read as
// their keyword.
func roleNames(list *ast.List) []string {
	if list == nil {
		return nil
	}
	names := make([]string, 0, len(list.Items))
	for _, item := range list.Items {
		rs, ok := item.(*ast.RoleSpec)
		if !ok {
			continue
		}
		names = append(names, roleName(rs))
	}
	return names
}

func roleName(rs *ast.RoleSpec) string {
	if rs == nil {
		return ""
	}
	switch ast.RoleSpecType(rs.Roletype) {
	case ast.ROLESPEC_CURRENT_ROLE:
		return "CURRENT_ROLE"
	case ast.ROLESPEC_CURRENT_USER:
		return "CURRENT_USER"
	case ast.ROLESPEC_SESSION_USER:
		return "SESSION_USER"
	case ast.ROLESPEC_PUBLIC:
		return "PUBLIC"
	default:
		return ident(rs.Rolename)
	}
}
