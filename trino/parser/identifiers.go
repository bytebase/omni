package parser

import (
	"strings"

	"github.com/bytebase/omni/trino/ast"
)

// This file implements Trino's `identifier` and `qualifiedName` grammar rules
// over the token stream, plus the exported name helpers bytebase's parser
// consumers rely on (NormalizeTrinoIdentifier, ExtractQualifiedNameParts).
//
// Trino's grammar:
//
//	identifier
//	    : IDENTIFIER_            # unquotedIdentifier
//	    | QUOTED_IDENTIFIER_     # quotedIdentifier
//	    | nonReserved            # unquotedIdentifier
//	    | BACKQUOTED_IDENTIFIER_ # backQuotedIdentifier
//	    | DIGIT_IDENTIFIER_      # digitIdentifier
//	    ;
//	qualifiedName : identifier (DOT_ identifier)* ;
//
// Three oracle-adjudicated divergences from a literal reading of that rule are
// baked in (see the migration divergence ledger):
//
//  1. DIGIT_IDENTIFIER_ (e.g. 1abc) is NOT accepted, despite the grammar
//     alternative: Trino 481 rejects unquoted digit-leading identifiers in
//     every position (table/column/dotted), so isIdentifierStart excludes
//     tokDigitIdent. A digit-leading name must be double-quoted.
//  2. BACKQUOTED_IDENTIFIER_ (`backtick` quoting) is NOT accepted, despite the
//     grammar alternative: Trino 481 rejects backtick-quoted identifiers
//     everywhere (only "double-quotes" are valid identifier quoting). The lexer
//     still tokenizes backticks (tokBackquotedIdent) — that is the lexer node's
//     concern — but the grammar rejects them, so isIdentifierStart excludes it.
//  3. qualifiedName applies the SAME identifier rule to every part — a reserved
//     keyword is rejected as a name part whether it is first or follows a dot.
//     (doris/snowflake accept reserved words post-dot; Trino does not.)

// isIdentifierStart reports whether kind may begin an identifier per Trino's
// `identifier` rule, as adjudicated by the Trino 481 oracle:
//   - tokIdent: unquoted IDENTIFIER_
//   - tokQuotedIdent: "double-quoted" QUOTED_IDENTIFIER_
//   - a non-reserved keyword (kind >= 700 && !IsReserved)
//
// tokDigitIdent and tokBackquotedIdent are deliberately excluded — Trino 481
// rejects both despite their presence in the legacy grammar. See the file
// header divergence notes.
func isIdentifierStart(kind TokenKind) bool {
	switch kind {
	case tokIdent, tokQuotedIdent:
		return true
	default:
		return kind >= 700 && !IsReserved(kind)
	}
}

// identFromToken builds an ast.Identifier from a token already confirmed to be
// an identifier start (per isIdentifierStart). A double-quoted token carries
// the '"' quote rune; bare and non-reserved-keyword tokens are unquoted.
func identFromToken(tok Token) *ast.Identifier {
	if tok.Kind == tokQuotedIdent {
		return &ast.Identifier{Value: tok.Str, Quoted: true, QuoteRune: '"', Loc: tok.Loc}
	}
	// tokIdent or a non-reserved keyword: unquoted, source text in Str.
	return &ast.Identifier{Value: tok.Str, Quoted: false, Loc: tok.Loc}
}

// parseIdentifier parses one identifier (Trino's `identifier` rule) and returns
// it as an *ast.Identifier. Returns a *ParseError if the current token is not a
// valid identifier start.
//
// Quoting metadata is preserved so callers can normalize per Trino's rules
// (case-insensitive for unquoted, case-sensitive for quoted) via
// ast.Identifier.Normalize.
func (p *Parser) parseIdentifier() (*ast.Identifier, error) {
	if !isIdentifierStart(p.cur.Kind) {
		return nil, p.identifierError()
	}
	return identFromToken(p.advance()), nil
}

// identifierError returns a *ParseError describing a missing identifier at the
// current token, with a message distinct from the generic syntax error so
// diagnostics are actionable.
func (p *Parser) identifierError() *ParseError {
	if p.cur.Kind == tokEOF {
		return &ParseError{Loc: p.cur.Loc, Msg: "expected identifier, found end of input"}
	}
	text := p.cur.Str
	if text == "" {
		text = TokenName(p.cur.Kind)
	}
	return &ParseError{Loc: p.cur.Loc, Msg: "expected identifier, found " + text}
}

// parseQualifiedName parses a dot-separated chain of identifiers (Trino's
// `qualifiedName` rule) and returns an *ast.QualifiedName. Every part is parsed
// with the same identifier rule, so a reserved keyword is rejected in every
// position (see the file header divergence note).
//
// A '.' inside a qualifiedName always commits to another part: once the dot is
// consumed, a valid identifier MUST follow, otherwise parsing errors. This
// matches Trino, which rejects `s.from` / `t.1abc` (reserved or digit-leading
// name part) with a syntax error rather than stopping at the first part. The
// all-columns form `t.*` is a separate primaryExpression rule, not a
// qualifiedName, so it never reaches here.
func (p *Parser) parseQualifiedName() (*ast.QualifiedName, error) {
	first, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	qn := &ast.QualifiedName{
		Parts: []*ast.Identifier{first},
		Loc:   first.Loc,
	}
	for p.cur.Kind == int('.') {
		p.advance() // consume '.' — an identifier must follow
		part, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		qn.Parts = append(qn.Parts, part)
		qn.Loc.End = part.Loc.End
	}
	return qn, nil
}

// NormalizeTrinoIdentifier returns the canonical form of a single Trino
// identifier given its raw source spelling (with or without surrounding
// quotes). A "double-quoted" identifier has its surrounding quotes removed and
// its case preserved; any other spelling folds to lower case (Trino unquoted
// identifiers are case-insensitive).
//
// This is a byte-for-byte port of the legacy bytebase NormalizeTrinoIdentifier
// helper so the bytebase consumers (lineage, completion, data-access-control
// name resolution) compare names exactly as before. In particular it
// intentionally matches legacy by:
//   - stripping ONLY surrounding double-quotes (a backtick-quoted spelling is
//     lower-cased whole — backticks are not valid Trino identifier quoting
//     anyway; see the divergence ledger), and
//   - NOT collapsing a doubled-quote ("") escape, mirroring legacy's plain
//     slice of the inter-quote bytes.
//
// Callers that have an already-lexed token should use ast.Identifier.Normalize
// instead, which collapses the "" escape from the decoded Value.
func NormalizeTrinoIdentifier(raw string) string {
	if strings.HasPrefix(raw, `"`) && strings.HasSuffix(raw, `"`) && len(raw) >= 2 {
		return raw[1 : len(raw)-1]
	}
	return strings.ToLower(raw)
}

// ExtractQualifiedNameParts parses a dotted name string and returns its
// per-component normalized parts (see ast.QualifiedName.NormalizedParts). It is
// the string-input counterpart of parseQualifiedName, for callers (catalog
// lookups, lineage) that hold a name string rather than a token stream. Mirrors
// the legacy bytebase ExtractQualifiedNameParts helper.
func ExtractQualifiedNameParts(input string) ([]string, error) {
	qn, errs := ParseQualifiedName(input)
	if len(errs) > 0 {
		return nil, &errs[0]
	}
	return qn.NormalizedParts(), nil
}

// ParseQualifiedName parses a complete qualified name from a standalone string,
// returning the *ast.QualifiedName and any ParseErrors. Trailing tokens after
// the name are reported as an error. Useful for catalog lookups and tests that
// have a name string but no token stream.
func ParseQualifiedName(input string) (*ast.QualifiedName, []ParseError) {
	p := &Parser{lexer: NewLexer(input), input: input}
	p.advance()

	qn, err := p.parseQualifiedName()
	if err != nil {
		if pe, ok := err.(*ParseError); ok {
			return nil, []ParseError{*pe}
		}
		return nil, []ParseError{{Msg: err.Error()}}
	}
	if p.cur.Kind != tokEOF {
		text := p.cur.Str
		if text == "" {
			text = TokenName(p.cur.Kind)
		}
		return qn, []ParseError{{Loc: p.cur.Loc, Msg: "unexpected token after qualified name: " + text}}
	}
	return qn, nil
}
