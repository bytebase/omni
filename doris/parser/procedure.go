package parser

import (
	"strings"

	"github.com/bytebase/omni/doris/ast"
)

// parseCreateProcedure parses the CREATE and REPLACE forms of the grammar's
// createProcedure alternative (DorisParser.g4:43):
//
//	(ALTER | CREATE (OR REPLACE)? | REPLACE) (PROCEDURE | PROC)
//	    name=multipartIdentifier LEFT_PAREN .*? RIGHT_PAREN .*?
//
// On entry, CREATE (with OR REPLACE) or REPLACE has already been consumed;
// cur is PROCEDURE or PROC.
//
// The rule reads the parameter list and everything after it as free text, so
// Doris accepts the parameter list and body this parser reads as well as forms
// it cannot make sense of, and refuses a name without parentheses
// (compat.tsv's statement/(ALTER|CREATE_(OR_REPLACE)?|REPLACE)_(PROCEDURE|PROC)*
// rows). Both are read here: the parameter list and BEGIN...END body first, so
// the AST carries them, and the free text when that does not fit, in RawText.
// The rule has no IF NOT EXISTS: the name follows PROCEDURE, so
// CREATE PROCEDURE IF NOT EXISTS p() is a syntax error.
func (p *Parser) parseCreateProcedure(startLoc ast.Loc, orReplace, replace bool) (ast.Node, error) {
	if _, err := p.expectProcedureKeyword(); err != nil {
		return nil, err
	}

	stmt := &ast.CreateProcedureStmt{
		OrReplace: orReplace,
		Replace:   replace,
		Loc:       startLoc,
	}

	// Procedure name (may be qualified: db.proc_name)
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	if err := p.parseProcedureTail(stmt); err != nil {
		return nil, err
	}
	return stmt, nil
}

// parseAlterProcedure parses ALTER PROCEDURE, the grammar's third verb
// (DorisParser.g4:43):
//
//	ALTER (PROCEDURE | PROC) name=multipartIdentifier LEFT_PAREN .*? RIGHT_PAREN .*?
//
// Everything after the name is free text, which Doris keeps and never runs:
// the pinned server parses ALTER PROCEDURE a() and ALTER PROCEDURE a(a zzz)
// (compat.tsv's statement(.)*/%empty and statement(.)*/._(.)* rows) and
// refuses a name without parentheses. cur is PROCEDURE or PROC on entry.
func (p *Parser) parseAlterProcedure(startLoc ast.Loc) (ast.Node, error) {
	if _, err := p.expectProcedureKeyword(); err != nil {
		return nil, err
	}
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt := &ast.AlterProcedureStmt{Name: name, Loc: startLoc}
	if p.cur.Kind != int('(') {
		return nil, p.syntaxErrorAtCur()
	}
	raw, err := p.collectProcedureParens()
	if err != nil {
		return nil, err
	}
	stmt.RawText = strings.TrimSpace(raw + " " + p.collectProcedureTail())
	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// expectProcedureKeyword consumes PROCEDURE or PROC, the two spellings the
// grammar's procedure statements take (DorisParser.g4:43-45).
func (p *Parser) expectProcedureKeyword() (Token, error) {
	if p.cur.Kind == kwPROCEDURE || p.cur.Kind == kwPROC {
		return p.advance(), nil
	}
	return Token{}, p.syntaxErrorAtCur()
}

// parseProcedureTail reads what CREATE PROCEDURE states after the name: the
// parameter list, an optional COMMENT, and the body. A statement whose text
// does not fit keeps it in RawText instead, as the grammar's `.*?` takes it.
func (p *Parser) parseProcedureTail(stmt *ast.CreateProcedureStmt) error {
	if p.cur.Kind != int('(') {
		return p.syntaxErrorAtCur()
	}

	saved := p.save()
	if err := p.parseProcedureListAndBody(stmt); err == nil {
		// The trailing `.*?` takes what follows the body too, so
		// REPLACE PROCEDURE a() BEGIN END zzz is a statement (the pinned
		// engine parses whatever stands there).
		stmt.RawText = p.collectProcedureTail()
		stmt.Loc.End = p.prev.Loc.End
		return nil
	}
	p.restore(saved)

	// The grammar reads LEFT_PAREN .*? RIGHT_PAREN .*? (DorisParser.g4:43):
	// the parenthesis closes at the first ')' whatever it holds, and the rest
	// of the text is kept as written. Doris refuses none of it, so neither
	// does this (compat.tsv's statement(OR_REPLACE)?/%empty row,
	// CREATE PROCEDURE a()).
	stmt.Parameters, stmt.Comment, stmt.Body = nil, "", ""
	raw, err := p.collectProcedureParens()
	if err != nil {
		return err
	}
	stmt.RawText = strings.TrimSpace(raw + " " + p.collectProcedureTail())
	stmt.Loc.End = p.prev.Loc.End
	return nil
}

// parseProcedureListAndBody reads the parameter list, the optional COMMENT and
// the BEGIN...END body the engine runs CREATE PROCEDURE with.
func (p *Parser) parseProcedureListAndBody(stmt *ast.CreateProcedureStmt) error {
	if _, err := p.expect(int('(')); err != nil {
		return err
	}
	if p.cur.Kind != int(')') {
		params, err := p.parseProcedureParams()
		if err != nil {
			return err
		}
		stmt.Parameters = params
	}
	if _, err := p.expect(int(')')); err != nil {
		return err
	}

	// Optional COMMENT 'text'
	if p.cur.Kind == kwCOMMENT {
		p.advance() // consume COMMENT
		if p.cur.Kind != tokString {
			return p.syntaxErrorAtCur()
		}
		stmt.Comment = p.cur.Str
		p.advance()
	}

	if p.cur.Kind != kwBEGIN {
		return p.syntaxErrorAtCur()
	}
	body, _, err := p.collectProcedureBody()
	if err != nil {
		return err
	}
	stmt.Body = body
	return nil
}

// collectProcedureParens reads the grammar's LEFT_PAREN .*? RIGHT_PAREN at cur
// and returns the parenthesized text as written. The text is the shortest one
// that closes the parenthesis, as the grammar's non-greedy `.*?` takes it, and
// it does not cross a statement separator: the parser is handed one statement
// at a time.
func (p *Parser) collectProcedureParens() (string, error) {
	open, err := p.expect(int('('))
	if err != nil {
		return "", err
	}
	for {
		if p.cur.Kind == tokEOF || p.cur.Kind == int(';') {
			return "", &ParseError{
				Position: p.cur.Loc.Start, End: p.cur.Loc.End,
				Message: "unexpected end of statement inside a procedure's parentheses (missing ')')",
			}
		}
		tok := p.advance()
		if tok.Kind == int(')') {
			return p.textBetween(open.Loc.End, tok.Loc.Start), nil
		}
	}
}

// collectProcedureTail reads the grammar's trailing `.*?` after a procedure's
// parentheses: the rest of the statement, which nothing reads. It ends at a
// separator no compound block encloses — the boundary Split draws — so a tail
// that opens BEGIN...END without closing it takes the statements after it, as
// it does for Doris and for Split.
func (p *Parser) collectProcedureTail() string {
	start := p.cur.Loc.Start
	end := start
	depth := 0
	for p.cur.Kind != tokEOF {
		if depth == 0 && p.cur.Kind == int(';') {
			break
		}
		switch p.cur.Kind {
		case kwBEGIN:
			// Split decides the transaction forms only outside a block: inside
			// one, a BEGIN always deepens it.
			if depth > 0 || p.beginOpensCompoundBlock() {
				depth++
			}
		case kwEND:
			if depth > 0 {
				depth--
			}
		}
		end = p.advance().Loc.End
	}
	if end <= start {
		return ""
	}
	return strings.TrimSpace(p.textBetween(start, end))
}

// beginOpensCompoundBlock reports whether the BEGIN at cur opens a compound
// block rather than starting a transaction, with the test Split applies
// (isDorisTCLBeginFollower, and the label after BEGIN WITH).
func (p *Parser) beginOpensCompoundBlock() bool {
	next := p.peekNext()
	if isDorisTCLBeginFollower(next) {
		return false
	}
	if next.Kind != kwWITH {
		return true
	}
	saved := p.save()
	p.advance() // consume BEGIN
	labelled := p.peekNext().Kind == kwLABEL
	p.restore(saved)
	return !labelled
}

// textBetween returns the source text between two byte offsets.
func (p *Parser) textBetween(start, end int) string {
	if start < p.baseOffset || end < start || end > len(p.input)+p.baseOffset {
		return ""
	}
	return p.input[start-p.baseOffset : end-p.baseOffset]
}

// parseProcedureParams parses a non-empty comma-separated list of procedure
// parameters. Called after '(' has been consumed and ')' has not yet been seen.
func (p *Parser) parseProcedureParams() ([]*ast.ProcedureParam, error) {
	var params []*ast.ProcedureParam

	for {
		param, err := p.parseProcedureParam()
		if err != nil {
			return nil, err
		}
		params = append(params, param)

		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}

	return params, nil
}

// parseProcedureParam parses one procedure parameter:
//
//	[IN | OUT | INOUT] param_name type
func (p *Parser) parseProcedureParam() (*ast.ProcedureParam, error) {
	paramStart := p.cur.Loc
	param := &ast.ProcedureParam{}

	// Optional direction: IN, OUT, INOUT.
	// IN is a reserved keyword (kwIN). OUT and INOUT are not in the keyword
	// table and will be lexed as tokIdent. Detect them by their string value
	// only when the token AFTER them is a valid identifier (the param name),
	// using a lookahead check.
	switch {
	case p.cur.Kind == kwIN && isIdentifierToken(p.peekNext().Kind):
		param.Direction = "IN"
		p.advance()
	case p.cur.Kind == tokIdent && strings.EqualFold(p.cur.Str, "out") && isIdentifierToken(p.peekNext().Kind):
		param.Direction = "OUT"
		p.advance()
	case p.cur.Kind == tokIdent && strings.EqualFold(p.cur.Str, "inout") && isIdentifierToken(p.peekNext().Kind):
		param.Direction = "INOUT"
		p.advance()
	}

	// Parameter name
	name, nameLoc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	param.Name = name

	// Type
	typeName, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	param.Type = typeName
	param.Loc = paramStart.Merge(nameLoc).Merge(ast.NodeLoc(typeName))

	return param, nil
}

// collectProcedureBody collects the raw text of a BEGIN...END block,
// tracking nesting so inner BEGIN...END blocks are captured correctly.
// Returns the raw text and its source location.
// cur must be BEGIN on entry.
func (p *Parser) collectProcedureBody() (string, ast.Loc, error) {
	beginTok, err := p.expect(kwBEGIN)
	if err != nil {
		return "", ast.Loc{}, err
	}

	bodyStart := beginTok.Loc.Start
	depth := 1

	var endTok Token
	for depth > 0 {
		if p.cur.Kind == tokEOF {
			return "", ast.Loc{}, &ParseError{
				Position: p.cur.Loc.Start, End: p.cur.Loc.End,
				Message: "unexpected end of input inside procedure body (missing END)",
			}
		}
		tok := p.advance()
		switch tok.Kind {
		case kwBEGIN:
			depth++
		case kwEND:
			depth--
			if depth == 0 {
				endTok = tok
			}
		}
	}

	bodyEnd := endTok.Loc.End
	bodyText := strings.TrimSpace(p.input[bodyStart-p.baseOffset : bodyEnd-p.baseOffset])
	return bodyText, ast.Loc{Start: bodyStart, End: bodyEnd}, nil
}

// parseCallProcedure parses a CALL statement.
// On entry, cur is CALL (not yet consumed).
//
// Syntax:
//
//	CALL proc_name([args])
func (p *Parser) parseCallProcedure() (ast.Node, error) {
	callTok := p.advance() // consume CALL

	stmt := &ast.CallProcedureStmt{
		Loc: callTok.Loc,
	}

	// Procedure name (may be qualified: db.proc_name)
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Argument list: ( [expr, ...] )
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}

	var args []ast.Node
	if p.cur.Kind != int(')') {
		args, err = p.parseExprList()
		if err != nil {
			return nil, err
		}
	}
	stmt.Args = args

	closeParen, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}

	stmt.Loc = callTok.Loc.Merge(closeParen.Loc)
	return stmt, nil
}

// parseDropProcedure parses a DROP PROCEDURE statement (DorisParser.g4:44):
//
//	DROP (PROCEDURE | PROC) (IF EXISTS)? name=multipartIdentifier
//
// On entry, DROP has already been consumed; cur is PROCEDURE or PROC.
func (p *Parser) parseDropProcedure(startLoc ast.Loc) (ast.Node, error) {
	if _, err := p.expectProcedureKeyword(); err != nil {
		return nil, err
	}

	stmt := &ast.DropProcedureStmt{}

	// Optional IF EXISTS
	if p.cur.Kind == kwIF {
		p.advance() // consume IF
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}

	// Procedure name
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc = startLoc.Merge(ast.NodeLoc(name))
	return stmt, nil
}
