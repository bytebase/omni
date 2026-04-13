package parser

import (
	"strings"

	"github.com/bytebase/omni/doris/ast"
)

// isIdentifierToken reports whether kind can form a standalone identifier in
// the first position of a qualified name (i.e., without dot context):
//   - tokIdent: unquoted identifier (e.g. myTable)
//   - tokQuotedIdent: backtick-quoted identifier (e.g. `my table`)
//   - a non-reserved keyword (e.g. COMMENT, BUCKETS) used as an identifier
func isIdentifierToken(kind TokenKind) bool {
	switch kind {
	case tokIdent, tokQuotedIdent:
		return true
	default:
		// A keyword token that is not reserved can serve as an identifier.
		if kind >= 700 && !IsReserved(kind) {
			return true
		}
		return false
	}
}

// isIdentifierTokenQualified reports whether kind is a valid identifier in a
// qualified (dot-following) position. After a '.', ALL keyword tokens —
// including reserved ones — are valid identifier parts because the context is
// syntactically unambiguous (e.g., db.table, catalog.select.foo).
func isIdentifierTokenQualified(kind TokenKind) bool {
	switch kind {
	case tokIdent, tokQuotedIdent:
		return true
	default:
		// Any keyword token (reserved or not) is valid after a dot.
		return kind >= 700
	}
}

// parseIdentifier parses a single identifier token and returns its string
// value and source location. Accepted tokens are:
//   - tokIdent: unquoted identifier
//   - tokQuotedIdent: backtick-quoted identifier (lexer has already stripped backticks)
//   - a non-reserved keyword used as an identifier
//
// Returns a syntax error if the current token is not a valid identifier.
func (p *Parser) parseIdentifier() (string, ast.Loc, error) {
	tok := p.cur
	switch tok.Kind {
	case tokIdent, tokQuotedIdent:
		p.advance()
		return tok.Str, tok.Loc, nil
	default:
		// Non-reserved keywords may be used as identifiers.
		if tok.Kind >= 700 && !IsReserved(tok.Kind) {
			p.advance()
			return tok.Str, tok.Loc, nil
		}
		return "", ast.Loc{}, p.syntaxErrorAtCur()
	}
}

// parseIdentifierQualified parses a single identifier token in a qualified
// (post-dot) position. In this context ALL keyword tokens are accepted as
// identifier parts because the syntax is unambiguous.
func (p *Parser) parseIdentifierQualified() (string, ast.Loc, error) {
	tok := p.cur
	switch tok.Kind {
	case tokIdent, tokQuotedIdent:
		p.advance()
		return tok.Str, tok.Loc, nil
	default:
		// After a '.', any keyword (reserved or not) is a valid identifier part.
		if tok.Kind >= 700 {
			p.advance()
			return tok.Str, tok.Loc, nil
		}
		return "", ast.Loc{}, p.syntaxErrorAtCur()
	}
}

// parseMultipartIdentifier parses a dot-separated qualified name:
//
//	identifier ('.' identifier)*
//
// Returns an *ast.ObjectName with Parts populated. The Loc on the returned
// node spans from the start of the first identifier to the end of the last.
//
// In the first position only non-reserved keywords are accepted as identifiers.
// In subsequent positions (after '.') all keywords are accepted because the
// dot context makes the parse unambiguous (e.g., db.select is valid).
func (p *Parser) parseMultipartIdentifier() (*ast.ObjectName, error) {
	first, loc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	name := &ast.ObjectName{
		Parts: []string{first},
		Loc:   loc,
	}

	// Consume additional '.identifier' segments.
	for p.cur.Kind == int('.') {
		// Peek past the dot to make sure we have an identifier; if not, stop
		// here rather than consuming the dot and getting a confusing error.
		// Note: peekNext looks at the token after cur (the '.').
		next := p.peekNext()
		if !isIdentifierTokenQualified(next.Kind) {
			break
		}

		p.advance() // consume '.'

		part, partLoc, err := p.parseIdentifierQualified()
		if err != nil {
			return nil, err
		}
		name.Parts = append(name.Parts, part)
		name.Loc.End = partLoc.End
	}

	return name, nil
}

// parseIdentifierOrString parses either an identifier or a string literal.
// Some Doris syntax allows either form in certain positions (e.g. catalog/db
// names in certain DDL contexts). Returns the string value and its location.
func (p *Parser) parseIdentifierOrString() (string, ast.Loc, error) {
	if p.cur.Kind == tokString {
		tok := p.advance()
		return tok.Str, tok.Loc, nil
	}
	return p.parseIdentifier()
}

// NormalizeIdentifier returns the canonical form of a Doris identifier.
// Doris is case-insensitive for unquoted identifiers, so unquoted identifiers
// are lowercased. Backtick-quoted identifiers preserve their original case
// (the lexer strips backticks and stores the raw content in Str, so the
// quoted flag drives the normalization decision).
func NormalizeIdentifier(name string, quoted bool) string {
	if quoted {
		return name
	}
	return strings.ToLower(name)
}

// NormalizeObjectName returns a normalized dot-joined representation of an
// ObjectName. Each part is lowercased to match Doris case-insensitive
// identifier comparison semantics. Useful for lookups and equality checks.
func NormalizeObjectName(name *ast.ObjectName) string {
	if name == nil {
		return ""
	}
	parts := make([]string, len(name.Parts))
	for i, p := range name.Parts {
		parts[i] = strings.ToLower(p)
	}
	return strings.Join(parts, ".")
}
