package review

import (
	"strings"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/oracle/parser"
)

// The parser's reserved-word lookup fills a keyword table on first use,
// without a lock; the first call here, before any review runs targets
// concurrently, keeps that write out of the goroutines.
var _ = parser.IsReservedKeyword("SELECT")

// ident writes an identifier as Oracle needs it in a statement: bare when
// it is a nonquoted identifier Oracle would fold to this very name, upper
// case and not reserved, and otherwise in double quotes, with a line break
// written as \n or \r, so a message stays one line. The parser hands the
// rules names already folded: an unquoted name upper-cased, a quoted one
// as written. Oracle identifiers cannot contain a double quote.
func ident(name string) string {
	if plainIdent(name) && !parser.IsReservedKeyword(name) {
		return name
	}
	return `"` + escapeLines(name) + `"`
}

// plainIdent reports whether a name reads the same unquoted: an upper-case
// letter, then upper-case letters, digits, _, $, or #.
func plainIdent(name string) bool {
	if name == "" {
		return false
	}
	for i := 0; i < len(name); i++ {
		c := name[i]
		switch {
		case c >= 'A' && c <= 'Z':
		case (c >= '0' && c <= '9' || c == '_' || c == '$' || c == '#') && i > 0:
		default:
			return false
		}
	}
	return true
}

// objectName writes an object reference: its schema when the statement
// names one, its name, and its database link.
func objectName(ref *ast.ObjectName) string {
	if ref == nil {
		return ""
	}
	s := ident(ref.Name)
	if ref.Schema != "" {
		s = ident(ref.Schema) + "." + s
	}
	if ref.DBLink != "" {
		s += "@" + escapeLines(ref.DBLink)
	}
	return s
}

// dmlTarget writes the target an UPDATE or DELETE names: its table, or
// what kind of target it is when it is not one.
func dmlTarget(e ast.TableExpr) string {
	switch v := e.(type) {
	case *ast.TableRef:
		return objectName(v.Name)
	case *ast.SubqueryRef, *ast.LateralRef:
		return "an inline view"
	case *ast.TableCollectionExpr:
		return "a TABLE collection"
	}
	return "the target"
}

// tableKey names a table of the change: its schema and its name, as
// Oracle stores them. The schema is "" when it is not known at all: the
// reference names none, after an ALTER SESSION SET CURRENT_SCHEMA the
// review cannot read; and it is unknownSchema for the target's database
// when its name is not known.
type tableKey struct {
	schema string
	name   string
}

// mayBe reports whether two keys may name the same table: the same name,
// in the same schema, or with either schema not known.
func (k tableKey) mayBe(o tableKey) bool {
	vague := func(schema string) bool { return schema == "" || schema == unknownSchema }
	return k.name == o.name && (k.schema == o.schema || vague(k.schema) || vague(o.schema))
}

// is reports whether two keys certainly name the same table.
func (k tableKey) is(o tableKey) bool {
	return k.schema != "" && k == o
}

// keyIn is the key of a reference resolved in the given current schema,
// "" when it is not known.
func keyIn(ref *ast.ObjectName, current string) tableKey {
	schema := ref.Schema
	if schema == "" {
		schema = current
	}
	return tableKey{schema: schema, name: ref.Name}
}

// sessionSchema returns the schema an ALTER SESSION SET CURRENT_SCHEMA
// moves unqualified names to, as Oracle folds it, and whether the
// statement sets it; the schema is "" when the value cannot be read.
func sessionSchema(n ast.Node) (string, bool) {
	v, ok := n.(*ast.AlterSessionStmt)
	if !ok || v.SetParams == nil {
		return "", false
	}
	for _, item := range v.SetParams.Items {
		p, ok := item.(*ast.SetParam)
		if !ok || !strings.EqualFold(p.Name, "CURRENT_SCHEMA") {
			continue
		}
		if c, ok := p.Value.(*ast.ColumnRef); ok && c.Table == "" && c.Schema == "" {
			return c.Column, true
		}
		return "", true
	}
	return "", false
}

// schemasAt returns, for each statement, the schema unqualified names
// resolve in when it runs: the target's database, db, until an ALTER
// SESSION SET CURRENT_SCHEMA moves it. "" is a schema not known.
func schemasAt(stmts []statement, db string) []string {
	out := make([]string, len(stmts))
	current := db
	for i := range stmts {
		out[i] = current
		if schema, ok := sessionSchema(stmts[i].node); ok {
			current = schema
		}
	}
	return out
}

// firstSessionChange returns the position of the first ALTER SESSION SET
// CURRENT_SCHEMA, the number of statements when there is none: from there
// on, an unqualified name of the change may mean another schema's object.
func firstSessionChange(stmts []statement) int {
	for i := range stmts {
		if _, ok := sessionSchema(stmts[i].node); ok {
			return i
		}
	}
	return len(stmts)
}
