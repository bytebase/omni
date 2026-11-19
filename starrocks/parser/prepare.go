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
		rel := p.cur.Loc.Start - p.baseOffset
		if rel < 0 || rel >= len(p.input) || p.input[rel] != '\'' {
			return nil, p.syntaxErrorAtCur()
		}
		// The string is SQL too: the server parses its text, refuses
		// `PREPARE a FROM 'a'` with the error inside it, `PREPARE p FROM ''`
		// and a comment for holding no statement, and
		// `PREPARE p FROM 'SELECT 1; SELECT 2'` with "only single statement is
		// supported". The nested parse starts at this parser's nesting depth,
		// so a string of wrappers cannot stack deeper than the limit.
		//
		// A literal written without escapes is parsed where it stands, so the
		// tree's Locs are the statement's and a query span reports the tables
		// it reads at their own bytes. An escaped literal's decoded text does
		// not map onto the source byte for byte, and rewriting a tree's Locs is
		// not a thing this parser does, so that form keeps the text and no
		// tree.
		node, err := p.parsePreparedString()
		if err != nil {
			return nil, err
		}
		stmt.Query = node
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

// parsePreparedString parses the statement inside PREPARE's string literal,
// which cur holds, and reports it at the literal's own bytes.
//
// The literal's text is not the statement's text when it holds escapes, so the
// decoded text is parsed and every node's Loc is then mapped back through the
// escape table: the tree answers where in the statement each table, column and
// operator stands, which is what a query span and a highlighter need.
func (p *Parser) parsePreparedString() (ast.Node, error) {
	rel := p.cur.Loc.Start - p.baseOffset
	end := p.cur.Loc.End - p.baseOffset
	if rel < 0 || end > len(p.input) || end-rel < 2 || p.input[rel] != '\'' || p.input[end-1] != '\'' {
		return nil, p.syntaxErrorAtCur()
	}
	text, at := decodePreparedString(p.input[rel+1:end-1], p.cur.Loc.Start+1)
	node, errs := parseSingle(text, 0, true, p.depth, false)
	if len(errs) > 0 || node == nil {
		return nil, &ParseError{
			Position: p.cur.Loc.Start, End: p.cur.Loc.End,
			Message: "PREPARE FROM string must hold one statement",
		}
	}
	ast.Inspect(node, func(n ast.Node) bool {
		loc := ast.NodeLoc(n)
		if loc.IsValid() {
			ast.SetNodeLoc(n, ast.Loc{Start: at(loc.Start), End: at(loc.End)})
		}
		return true
	})
	return node, nil
}

// decodePreparedString decodes a single-quoted literal's content and returns,
// for every byte of the result, the byte of the source it came from, base being
// the source offset of the content. A doubled quote decodes to one quote at the
// offset of the pair, and a backslash escape to the character after it.
func decodePreparedString(raw string, base int) (string, func(int) int) {
	var b strings.Builder
	b.Grow(len(raw))
	// starts[i] is the source offset of decoded byte i; starts[len] is the end.
	starts := make([]int, 0, len(raw)+1)
	for i := 0; i < len(raw); i++ {
		switch {
		case raw[i] == '\\' && i+1 < len(raw):
			starts = append(starts, base+i)
			b.WriteByte(raw[i+1])
			i++
		case raw[i] == '\'' && i+1 < len(raw) && raw[i+1] == '\'':
			starts = append(starts, base+i)
			b.WriteByte('\'')
			i++
		default:
			starts = append(starts, base+i)
			b.WriteByte(raw[i])
		}
	}
	starts = append(starts, base+len(raw))
	return b.String(), func(pos int) int {
		if pos < 0 {
			return -1
		}
		if pos >= len(starts) {
			return starts[len(starts)-1]
		}
		return starts[pos]
	}
} // parseExecute parses executeStatement (StarRocks.g4:2291):
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
