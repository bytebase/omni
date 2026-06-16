package parser

import (
	"strings"

	"github.com/bytebase/omni/pg/plpgsql/ast"

	pgparser "github.com/bytebase/omni/pg/parser"
)

// parseDeclareSection parses the DECLARE section of a block.
// The DECLARE keyword has already been consumed.
// Parses variable, cursor, and alias declarations until BEGIN is found.
func (p *Parser) parseDeclareSection() ([]ast.Node, error) {
	var decls []ast.Node

	for !p.isKeyword("BEGIN") && !p.isEOF() {
		// Check for label between DECLARE and BEGIN — this is an error
		if p.isOp("<<") {
			return nil, p.errorf("syntax error at or near %q", p.tokenText(p.cur))
		}

		// Extra DECLARE keyword within section — skip it
		if p.isKeyword("DECLARE") {
			p.advance()
			continue
		}

		// Parse a single declaration
		decl, err := p.parseDecl()
		if err != nil {
			return nil, err
		}
		decls = append(decls, decl)

		// Expect semicolon after declaration
		if p.cur.Type != ';' {
			return nil, p.errorf("syntax error at or near %q, expected ;", p.tokenText(p.cur))
		}
		p.advance() // consume ;
	}

	return decls, nil
}

// parseDecl parses a single declaration (variable, cursor, or alias).
// The declaration name has NOT been consumed yet.
func (p *Parser) parseDecl() (ast.Node, error) {
	startPos := p.pos()

	// Read the variable name
	if !p.isIdentOrUnreserved() {
		return nil, p.errorf("syntax error at or near %q, expected variable name", p.tokenText(p.cur))
	}
	name := p.identText()
	p.advance()

	// Check for ALIAS declaration
	if p.isKeyword("ALIAS") {
		return p.parseAliasDecl(name, startPos)
	}

	// Check for CONSTANT
	isConstant := false
	if p.isKeyword("CONSTANT") {
		isConstant = true
		p.advance()
	}

	// Check for CURSOR (direct cursor declaration)
	if p.isKeyword("CURSOR") {
		return p.parseCursorDecl(name, ast.ScrollNone, startPos)
	}

	// Check for SCROLL/NO SCROLL CURSOR
	if p.isKeyword("SCROLL") {
		p.advance()
		if !p.isKeyword("CURSOR") {
			return nil, p.errorf("syntax error at or near %q, expected CURSOR after SCROLL", p.tokenText(p.cur))
		}
		return p.parseCursorDecl(name, ast.ScrollYes, startPos)
	}
	if p.isKeyword("NO") {
		// Peek to see if next is SCROLL
		next := p.peekNext()
		if next.Type == pgparser.SCROLL || (next.Type == pgparser.IDENT && strings.EqualFold(next.Str, "SCROLL")) {
			p.advance() // consume NO
			p.advance() // consume SCROLL
			if !p.isKeyword("CURSOR") {
				return nil, p.errorf("syntax error at or near %q, expected CURSOR after NO SCROLL", p.tokenText(p.cur))
			}
			return p.parseCursorDecl(name, ast.ScrollNo, startPos)
		}
	}

	// Parse type name
	typeName, err := p.parseDeclTypeName()
	if err != nil {
		return nil, err
	}

	// Parse optional COLLATE
	collation := ""
	if p.isKeyword("COLLATE") {
		p.advance()
		// Collation name: can be a double-quoted identifier or unquoted identifier
		// Use Str directly for IDENT to preserve case of quoted identifiers
		if p.isIdent() {
			collation = p.cur.Str
			p.advance()
		} else if p.isAnyKeywordAsIdent() {
			collation = p.identText()
			p.advance()
		} else {
			return nil, p.errorf("syntax error at or near %q, expected collation name", p.tokenText(p.cur))
		}
	}

	// Parse optional NOT NULL
	isNotNull := false
	if p.isKeyword("NOT") {
		p.advance()
		if err := p.expectKeyword("NULL"); err != nil {
			return nil, err
		}
		isNotNull = true
	}

	// Parse optional default expression
	defaultExpr := ""
	hasDefault := false
	if p.isKeyword("DEFAULT") {
		p.advance()
		hasDefault = true
		expr, err := p.collectUntilSemicolon()
		if err != nil {
			return nil, err
		}
		defaultExpr = expr
	} else if p.cur.Type == pgparser.COLON_EQUALS || p.cur.Type == '=' {
		p.advance()
		hasDefault = true
		expr, err := p.collectUntilSemicolon()
		if err != nil {
			return nil, err
		}
		defaultExpr = expr
	}

	// Validate: NOT NULL requires a default
	if isNotNull && !hasDefault {
		return nil, &ParseError{
			Message:  "variable \"" + name + "\" with NOT NULL must have a default value",
			Position: startPos,
		}
	}

	// Validate: CONSTANT requires a default
	if isConstant && !hasDefault {
		return nil, &ParseError{
			Message:  "variable \"" + name + "\" declared CONSTANT must have a default value",
			Position: startPos,
		}
	}

	return &ast.PLDeclare{
		Name:      name,
		TypeName:  typeName,
		Constant:  isConstant,
		NotNull:   isNotNull,
		Collation: collation,
		Default:   defaultExpr,
		Loc:       ast.Loc{Start: startPos, End: p.pos()},
	}, nil
}

// parseDeclTypeName parses a type name in a variable declaration.
// Handles: simple types, qualified types (schema.type), %TYPE, %ROWTYPE, arrays.
func (p *Parser) parseDeclTypeName() (string, error) {
	if !p.isIdent() && !p.isAnyKeywordAsIdent() {
		return "", p.errorf("syntax error at or near %q, expected type name", p.tokenText(p.cur))
	}

	start := p.pos()

	var closers []int
	for !p.isEOF() {
		if len(closers) == 0 {
			if p.cur.Type == ',' || p.cur.Type == ')' || p.cur.Type == ';' ||
				p.cur.Type == pgparser.COLON_EQUALS || p.cur.Type == '=' ||
				p.isKeyword("COLLATE") || p.isKeyword("NOT") || p.isKeyword("DEFAULT") {
				break
			}
		}

		switch p.cur.Type {
		case '(':
			closers = append(closers, ')')
			p.advance()
			continue
		case '[':
			closers = append(closers, ']')
			p.advance()
			continue
		case ')', ']':
			if len(closers) == 0 {
				return "", p.errorf("syntax error at or near %q, unexpected closing delimiter", p.tokenText(p.cur))
			}
			want := closers[len(closers)-1]
			if p.cur.Type != want {
				return "", p.errorf("syntax error at or near %q, expected closing delimiter", p.tokenText(p.cur))
			}
			closers = closers[:len(closers)-1]
			p.advance()
			continue
		case '%':
			p.advance()
			if p.isKeyword("TYPE") || p.isKeyword("ROWTYPE") {
				p.advance()
				continue
			}
			return "", p.errorf("syntax error at or near %q, expected TYPE or ROWTYPE after %%", p.tokenText(p.cur))
		}

		p.advance()
	}
	if len(closers) > 0 {
		return "", p.errorf("syntax error at or near %q, expected closing delimiter", p.tokenText(p.cur))
	}

	typeName := strings.TrimSpace(p.source[start:p.prev.End])
	return typeName, nil
}

// parseAliasDecl parses an ALIAS FOR declaration.
// The variable name has been consumed; cur is on ALIAS.
func (p *Parser) parseAliasDecl(name string, startPos int) (ast.Node, error) {
	p.advance() // consume ALIAS

	if err := p.expectKeyword("FOR"); err != nil {
		return nil, err
	}

	// The reference can be $N (PARAM token) or an identifier
	var refName string
	if p.cur.Type == pgparser.PARAM {
		refName = p.source[p.cur.Loc:p.cur.End]
		p.advance()
	} else if p.isIdent() || p.isAnyKeywordAsIdent() {
		refName = p.identText()
		p.advance()
	} else {
		return nil, p.errorf("syntax error at or near %q, expected parameter reference", p.tokenText(p.cur))
	}

	return &ast.PLAliasDecl{
		Name:    name,
		RefName: refName,
		Loc:     ast.Loc{Start: startPos, End: p.prev.End},
	}, nil
}

// parseCursorDecl parses a cursor declaration after CURSOR keyword.
// The CURSOR keyword has NOT been consumed. name and scroll have been determined.
func (p *Parser) parseCursorDecl(name string, scroll ast.ScrollOption, startPos int) (ast.Node, error) {
	p.advance() // consume CURSOR

	// Optional parameter list: (p1 type1, p2 type2, ...)
	var args []ast.PLCursorArg
	if p.cur.Type == '(' {
		p.advance() // consume (
		for p.cur.Type != ')' && !p.isEOF() {
			// Parse argument name
			if !p.isIdent() && !p.isAnyKeywordAsIdent() {
				return nil, p.errorf("syntax error at or near %q, expected parameter name", p.tokenText(p.cur))
			}
			argName := p.identText()
			p.advance()

			// Parse argument type
			argType, err := p.parseDeclTypeName()
			if err != nil {
				return nil, err
			}

			args = append(args, ast.PLCursorArg{Name: argName, TypeName: argType})

			// Consume comma if present
			if p.cur.Type == ',' {
				p.advance()
			}
		}
		if p.cur.Type != ')' {
			return nil, p.errorf("syntax error at or near %q, expected )", p.tokenText(p.cur))
		}
		p.advance() // consume )
	}

	// Expect FOR or IS
	if !p.isKeyword("FOR") && !p.isKeyword("IS") {
		return nil, p.errorf("syntax error at or near %q, expected FOR or IS", p.tokenText(p.cur))
	}
	p.advance() // consume FOR/IS

	// Collect query text until semicolon
	query, err := p.collectUntilSemicolon()
	if err != nil {
		return nil, err
	}

	return &ast.PLCursorDecl{
		Name:   name,
		Scroll: scroll,
		Args:   args,
		Query:  query,
		Loc:    ast.Loc{Start: startPos, End: p.pos()},
	}, nil
}

// isIdentOrUnreserved checks if the current token can serve as a variable name.
// Allows IDENT tokens and any PL/pgSQL unreserved keyword.
// Also allows SQL keywords that are not PL/pgSQL reserved keywords.
func (p *Parser) isIdentOrUnreserved() bool {
	if p.cur.Type == pgparser.IDENT {
		// Check if it's a PL/pgSQL reserved keyword
		if cat, ok := LookupPLKeyword(p.cur.Str); ok && cat == PLReserved {
			return false
		}
		return true
	}
	// SQL keywords (type > 256) can be used as PL/pgSQL variable names
	// as long as they're not PL/pgSQL reserved keywords
	if p.isAnyKeywordAsIdent() {
		text := strings.ToLower(p.source[p.cur.Loc:p.cur.End])
		if cat, ok := LookupPLKeyword(text); ok && cat == PLReserved {
			return false
		}
		return true
	}
	return false
}

// Ensure imports are used.
var _ = pgparser.PARAM
