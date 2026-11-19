package parser

import (
	"strings"

	"github.com/bytebase/omni/starrocks/ast"
)

// Prepared statements (StarRocks.g4:2282-2302). EXECUTE also opens
// executeAsStatement (:1646), which the word AS tells apart.

// parsePrepare parses prepareStatement (StarRocks.g4:2282):
//
//	PREPARE identifier FROM prepareSql
//
// prepareSql (:2300) is a statement or a single-quoted string. A statement is
// parsed in place, so a body that does not parse is a syntax error here, as it
// is on the pinned server.
//
// On entry PREPARE has NOT been consumed.
func (p *Parser) parsePrepare(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume PREPARE
	name, nameLoc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt := &ast.PrepareStmt{Name: &ast.ObjectName{Parts: []string{name}, Loc: nameLoc}}
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}
	if p.cur.Kind == tokString {
		// prepareSql admits SINGLE_QUOTED_TEXT only, and the server refuses
		// `PREPARE p FROM "SELECT 1"`. The quote is the source byte at the
		// token's start, which the token does not keep.
		if start := p.cur.Loc.Start - p.baseOffset; start < 0 || start >= len(p.input) || p.input[start] != '\'' {
			return nil, p.syntaxErrorAtCur()
		}
		// The string is SQL too: the server parses its text, refuses
		// `PREPARE a FROM 'a'` with the error inside it, `PREPARE p FROM ''`
		// and a comment for holding no statement, and
		// `PREPARE p FROM 'SELECT 1; SELECT 2'` with "only single statement is
		// supported". The nested parse starts at this parser's nesting depth,
		// so a string of wrappers cannot stack deeper than the input limit.
		// Positions in the string are not the statement's, so the error
		// carries none.
		nested := parseAll(p.cur.Str, true, p.depth)
		if len(nested.Errors) > 0 || nested.File == nil || len(nested.File.Stmts) != 1 {
			return nil, &ParseError{
				Position: p.cur.Loc.Start, End: p.cur.Loc.End,
				Message: "PREPARE FROM string must hold one statement",
			}
		}
		stmt.SQL = p.cur.Str
		stmt.Loc = startLoc.Merge(p.cur.Loc)
		p.advance()
		return stmt, nil
	}
	query, err := p.parseStmt()
	if err != nil {
		return nil, err
	}
	stmt.Query = query
	stmt.Loc = startLoc.Merge(ast.NodeLoc(query))
	return stmt, nil
}

// parseExecute parses executeStatement (StarRocks.g4:2291):
//
//	EXECUTE identifier [USING '@'identifierOrString [',' '@'identifierOrString]*]
//
// On entry EXECUTE has NOT been consumed.
func (p *Parser) parseExecute(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume EXECUTE
	name, nameLoc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt := &ast.ExecuteStmt{Name: &ast.ObjectName{Parts: []string{name}, Loc: nameLoc}}
	endLoc := nameLoc
	if p.cur.Kind == kwUSING {
		p.advance() // consume USING
		for {
			if _, err := p.expect(int('@')); err != nil {
				return nil, err
			}
			v, vLoc, err := p.parseIdentifierOrString()
			if err != nil {
				return nil, err
			}
			stmt.Using = append(stmt.Using, v)
			endLoc = vLoc
			if p.cur.Kind != int(',') {
				break
			}
			p.advance() // consume ','
		}
	}
	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseExecuteAs parses executeAsStatement (StarRocks.g4:1646):
//
//	EXECUTE AS user [WITH NO REVERT]
//
// On entry EXECUTE has NOT been consumed.
func (p *Parser) parseExecuteAs(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume EXECUTE
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}
	user, err := p.parseUserIdentity()
	if err != nil {
		return nil, err
	}
	stmt := &ast.ExecuteAsStmt{User: user}
	endLoc := user.Loc
	if p.cur.Kind == kwWITH {
		p.advance() // consume WITH
		if _, err := p.expect(kwNO); err != nil {
			return nil, err
		}
		revertTok, err := p.expect(kwREVERT)
		if err != nil {
			return nil, err
		}
		stmt.WithNoRevert = true
		endLoc = revertTok.Loc
	}
	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseDeallocate parses deallocateStatement (StarRocks.g4:2295):
//
//	(DEALLOCATE | DROP) PREPARE identifier
//
// On entry the verb has been consumed; cur is PREPARE.
func (p *Parser) parseDeallocate(startLoc ast.Loc) (ast.Node, error) {
	if _, err := p.expect(kwPREPARE); err != nil {
		return nil, err
	}
	name, nameLoc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	return &ast.DeallocateStmt{
		Name: &ast.ObjectName{Parts: []string{name}, Loc: nameLoc},
		Loc:  startLoc.Merge(nameLoc),
	}, nil
}

// parseTranslate parses translateStatement (StarRocks.g4:2073):
//
//	TRANSLATE dialect translateSQL
//
// translateSQL (:2077) is `.+`: every token to the end of the statement, which
// StarRocks hands to its SQL converter as text. It needs at least one token,
// and none of the tokens is read as SQL.
//
// On entry TRANSLATE has NOT been consumed.
func (p *Parser) parseTranslate() (ast.Node, error) {
	startLoc := p.cur.Loc
	p.advance() // consume TRANSLATE
	dialect, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt := &ast.TranslateStmt{Dialect: dialect, Loc: startLoc.Merge(p.prev.Loc)}
	// translateSQL's `.+` needs a token. At the end of the input there is none,
	// and the pinned server refuses a bare `TRANSLATE a` there; its terminator
	// is one, so `TRANSLATE a;` is accepted.
	if p.cur.Kind == tokEOF {
		if !p.terminated {
			return nil, p.syntaxErrorAtCur()
		}
		// The terminator is the token `.+` matched, and it is the source byte
		// at the end of the segment, so the span covers it even though the
		// translated text is empty.
		stmt.Loc.End = p.cur.Loc.Start + 1
		return stmt, nil
	}
	rawStart := p.cur.Loc.Start
	for p.cur.Kind != tokEOF {
		p.advance()
	}
	stmt.SQL = strings.TrimSpace(p.utilAdminRawText(rawStart))
	stmt.Loc = startLoc.Merge(p.prev.Loc)
	return stmt, nil
}
