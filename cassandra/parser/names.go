package parser

import (
	"strings"

	"github.com/bytebase/omni/cassandra/ast"
)

// The keyword sets below are Parser.g's, at cassandra-5.0.9, for the name
// rules the DDL, DESCRIBE and type parsers read: which words a name may be.

// basicUnreservedKeywords is basic_unreserved_keyword (Parser.g:1956).
var basicUnreservedKeywords = map[string]bool{
	"access": true, "aggregate": true, "aggregates": true, "all": true,
	"ann": true, "as": true, "called": true, "cidrs": true, "cluster": true,
	"clustering": true, "compact": true, "contains": true, "custom": true,
	"datacenters": true, "default": true, "exists": true, "filtering": true,
	"finalfunc": true, "frozen": true, "function": true, "functions": true,
	"group": true, "hashed": true, "identity": true, "initcond": true,
	"input": true, "internals": true, "keys": true, "keyspaces": true,
	"language": true, "like": true, "list": true, "login": true, "map": true,
	"masked": true, "mbean": true, "mbeans": true, "nologin": true,
	"nosuperuser": true, "only": true, "options": true, "partition": true,
	"password": true, "per": true, "permission": true, "permissions": true,
	"replace": true, "returns": true, "role": true, "roles": true,
	"select_masked": true, "sfunc": true, "static": true, "storage": true,
	"stype": true, "superuser": true, "tables": true, "trigger": true,
	"tuple": true, "type": true, "types": true, "unmask": true, "unset": true,
	"user": true, "users": true, "values": true, "vector": true,
	"columnfamilies": true, // K_TABLES's other spelling (Lexer.g:96)
}

// otherUnreservedKeywords completes unreserved_keyword (Parser.g:1946): the
// native types (Parser.g:1873) and the words it adds.
var otherUnreservedKeywords = map[string]bool{
	"ascii": true, "bigint": true, "blob": true, "boolean": true, "cast": true,
	"count": true, "counter": true, "date": true, "decimal": true,
	"distinct": true, "double": true, "duration": true, "float": true,
	"inet": true, "int": true, "json": true, "key": true, "maxwritetime": true,
	"smallint": true, "text": true, "time": true, "timestamp": true,
	"timeuuid": true, "tinyint": true, "ttl": true, "uuid": true,
	"varchar": true, "varint": true, "writetime": true,
}

// reservedKeywords are the words Lexer.g reads as a keyword, a boolean, NaN
// or Infinity that no unreserved_keyword lists: never a name.
var reservedKeywords = map[string]bool{
	"add": true, "allow": true, "alter": true, "and": true, "apply": true,
	"asc": true, "authorize": true, "batch": true, "begin": true, "by": true,
	"columnfamily": true, "create": true, "delete": true, "desc": true,
	"describe": true, "drop": true, "entries": true, "execute": true,
	"from": true, "full": true, "grant": true, "if": true, "in": true,
	"index": true, "insert": true, "into": true, "is": true, "keyspace": true,
	"limit": true, "materialized": true, "modify": true, "norecursive": true,
	"not": true, "null": true, "of": true, "on": true, "or": true,
	"order": true, "primary": true, "rename": true, "revoke": true,
	"schema": true, "select": true, "set": true, "table": true, "to": true,
	"token": true, "truncate": true, "unlogged": true, "update": true,
	"use": true, "using": true, "view": true, "where": true, "with": true,
	"true": true, "false": true, "nan": true, "infinity": true,
}

// tokWord is the lower-case word an unquoted token spells, or "" for a
// token that is no word (a quoted name, a literal, punctuation).
func tokWord(tok Token) string {
	if tok.Type == tokIDENT || isKeyword(tok.Type) {
		return strings.ToLower(tok.Str)
	}
	return ""
}

// isUnreservedKeyword reports whether tok is an unreserved_keyword
// (Parser.g:1946).
func isUnreservedKeyword(tok Token) bool {
	w := tokWord(tok)
	return basicUnreservedKeywords[w] || otherUnreservedKeywords[w]
}

// isCQLName reports whether tok is ident (Parser.g:1454): IDENT,
// QUOTED_NAME or an unreserved keyword. A word Cassandra does not reserve is
// an IDENT to it, even where omni's lexer reads a keyword (REPLICATION).
func isCQLName(tok Token) bool {
	if tok.Type == tokQUOTED {
		return true
	}
	w := tokWord(tok)
	return w != "" && !reservedKeywords[w]
}

// isNonTypeIdent reports whether tok is non_type_ident (Parser.g:1939), the
// last part of a user-defined type's name: IDENT, QUOTED_NAME, a
// basic_unreserved_keyword or KEY. The native types and the other words of
// unreserved_keyword are not.
func isNonTypeIdent(tok Token) bool {
	if tok.Type == tokQUOTED {
		return true
	}
	w := tokWord(tok)
	if w == "" || reservedKeywords[w] {
		return false
	}
	return !otherUnreservedKeywords[w] || w == "key"
}

// parseNameIf parses one name token that ok accepts.
func (p *Parser) parseNameIf(ok func(Token) bool) (*ast.Identifier, error) {
	tok := p.cur
	if !ok(tok) {
		return nil, p.errorf("expected identifier, got %s", p.tokenDesc())
	}
	p.advance()
	return &ast.Identifier{Name: tok.Str, Quoted: tok.Type == tokQUOTED, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil
}

// parseCQLQualifiedName parses (ident '.')? ident, both parts by isCQLName.
func (p *Parser) parseCQLQualifiedName() (*ast.QualifiedName, error) {
	start := p.curLoc()
	first, err := p.parseNameIf(isCQLName)
	if err != nil {
		return nil, err
	}
	parts := []*ast.Identifier{first}
	if p.match(tokDOT) {
		second, err := p.parseNameIf(isCQLName)
		if err != nil {
			return nil, err
		}
		parts = append(parts, second)
	}
	return &ast.QualifiedName{Parts: parts, Loc: p.makeLoc(start)}, nil
}

// isLexerKeyword reports whether the lower-case word w is a keyword of
// Lexer.g, reserved or not (or a boolean, NaN or Infinity), as opposed to an
// IDENT. omni's lexer makes more words keywords than Cassandra's does.
func isLexerKeyword(w string) bool {
	return reservedKeywords[w] || basicUnreservedKeywords[w] || otherUnreservedKeywords[w]
}

// isStringLiteral reports whether typ is a STRING_LITERAL (Lexer.g:258):
// quoted with single quotes or with $$.
func isStringLiteral(typ int) bool {
	return typ == tokSTRING || typ == tokCODEBLOCK
}

// isAllowedFunctionName reports whether tok is allowedFunctionName
// (Parser.g:1613): IDENT, QUOTED_NAME, an unreserved_function_keyword (a
// basic_unreserved_keyword or a native type), TOKEN or COUNT. The other words
// of unreserved_keyword (TTL, WRITETIME, MAXWRITETIME, KEY, CAST, JSON,
// DISTINCT) name no function.
func isAllowedFunctionName(tok Token) bool {
	if tok.Type == tokQUOTED {
		return true
	}
	switch w := tokWord(tok); w {
	case "":
		return false
	case "token", "count":
		return true
	case "ttl", "writetime", "maxwritetime", "key", "cast", "json", "distinct":
		return false
	default:
		return !reservedKeywords[w]
	}
}

// parseFunctionQualifiedName parses functionName (Parser.g:1608):
// (keyspaceName '.')? allowedFunctionName.
func (p *Parser) parseFunctionQualifiedName() (*ast.QualifiedName, error) {
	start := p.curLoc()
	var parts []*ast.Identifier
	if p.peekNext().Type == tokDOT {
		ks, err := p.parseNameIf(isCQLName)
		if err != nil {
			return nil, err
		}
		p.advance() // '.'
		parts = append(parts, ks)
	}
	name, err := p.parseNameIf(isAllowedFunctionName)
	if err != nil {
		return nil, err
	}
	parts = append(parts, name)
	return &ast.QualifiedName{Parts: parts, Loc: p.makeLoc(start)}, nil
}
