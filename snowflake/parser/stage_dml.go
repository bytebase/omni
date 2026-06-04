package parser

import (
	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Stage file DML — PUT / GET / LIST (LS) / REMOVE (RM) (T5.2)
// ---------------------------------------------------------------------------
//
// These statements move files between the client filesystem and an internal
// stage (PUT/GET) or enumerate/delete staged files (LIST/REMOVE). They share
// the stage-reference, local-file, and option-list helpers defined in copy.go
// (parseStageRef / parseLocalFile / parseCopyOption family), so the open-ended
// option-vocabulary approach is identical: transfer options (PARALLEL,
// AUTO_COMPRESS, SOURCE_COMPRESSION, OVERWRITE) and PATTERN are parsed as
// KEY = value pairs rather than enumerated.

// ---------------------------------------------------------------------------
// PUT
// ---------------------------------------------------------------------------

// parsePutStmt parses a PUT statement:
//
//	PUT file://<path> <stage>
//	    [ PARALLEL = n ] [ AUTO_COMPRESS = b ]
//	    [ SOURCE_COMPRESSION = kw ] [ OVERWRITE = b ]
//
// The PUT keyword has NOT been consumed when this function is called.
func (p *Parser) parsePutStmt() (ast.Node, error) {
	putTok := p.advance() // consume PUT

	stmt := &ast.PutStmt{Loc: ast.Loc{Start: putTok.Loc.Start}}

	// Source: a local file:// path (or quoted/bare local path).
	file, err := p.parseLocalFile()
	if err != nil {
		return nil, err
	}
	stmt.File = file

	// Destination: an internal stage.
	stage, err := p.parsePutGetStage()
	if err != nil {
		return nil, err
	}
	stmt.Stage = stage

	// Transfer options: PARALLEL / AUTO_COMPRESS / SOURCE_COMPRESSION / OVERWRITE.
	opts, err := p.parseCopyOptions()
	if err != nil {
		return nil, err
	}
	stmt.Options = opts

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// GET
// ---------------------------------------------------------------------------

// parseGetStmt parses a GET statement:
//
//	GET <stage> file://<dir> [ PARALLEL = n ] [ PATTERN = '...' ]
//
// The destination may be a file:// URI or a bare local path (e.g.
// `my_target_path`). The GET keyword has NOT been consumed.
func (p *Parser) parseGetStmt() (ast.Node, error) {
	getTok := p.advance() // consume GET

	stmt := &ast.GetStmt{Loc: ast.Loc{Start: getTok.Loc.Start}}

	// Source: an internal stage.
	stage, err := p.parsePutGetStage()
	if err != nil {
		return nil, err
	}
	stmt.Stage = stage

	// Destination: a local directory (file:// or bare path).
	target, err := p.parseLocalFile()
	if err != nil {
		return nil, err
	}
	stmt.Target = target

	// Transfer options: PARALLEL / PATTERN.
	opts, err := p.parseCopyOptions()
	if err != nil {
		return nil, err
	}
	stmt.Options = opts

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parsePutGetStage parses the stage operand of a PUT/GET. PUT/GET operate only
// on internal stages, which are always '@'-prefixed; this is a thin wrapper over
// parseStageRef that produces a clear error if the operand is not a stage.
func (p *Parser) parsePutGetStage() (*ast.StageLocation, error) {
	if p.cur.Type != '@' {
		return nil, p.syntaxErrorAtCur()
	}
	return p.parseStageRef()
}

// ---------------------------------------------------------------------------
// LIST / LS
// ---------------------------------------------------------------------------

// parseListStmt parses a LIST or LS statement:
//
//	{ LIST | LS } <stage> [ PATTERN = '<regex>' ]
//
// short is true when the statement used the LS alias. The leading keyword
// (LIST keyword, or LS as an identifier) has NOT been consumed.
func (p *Parser) parseListStmt(short bool) (ast.Node, error) {
	kwTok := p.advance() // consume LIST / LS

	stmt := &ast.ListStmt{Short: short, Loc: ast.Loc{Start: kwTok.Loc.Start}}

	stage, err := p.parseStageOperand()
	if err != nil {
		return nil, err
	}
	stmt.Stage = stage

	pat, err := p.parseOptionalPattern()
	if err != nil {
		return nil, err
	}
	stmt.Pattern = pat

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// REMOVE / RM
// ---------------------------------------------------------------------------

// parseRemoveStmt parses a REMOVE or RM statement:
//
//	{ REMOVE | RM } <stage> [ PATTERN = '<regex>' ]
//
// short is true when the statement used the RM alias. The leading keyword
// (REMOVE keyword, or RM as an identifier) has NOT been consumed.
func (p *Parser) parseRemoveStmt(short bool) (ast.Node, error) {
	kwTok := p.advance() // consume REMOVE / RM

	stmt := &ast.RemoveStmt{Short: short, Loc: ast.Loc{Start: kwTok.Loc.Start}}

	stage, err := p.parseStageOperand()
	if err != nil {
		return nil, err
	}
	stmt.Stage = stage

	pat, err := p.parseOptionalPattern()
	if err != nil {
		return nil, err
	}
	stmt.Pattern = pat

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// Shared: stage operand + optional PATTERN
// ---------------------------------------------------------------------------

// parseStageOperand parses the stage operand of LIST/REMOVE. LIST/REMOVE accept
// internal and external stages (both '@'-prefixed); a non-stage operand is an
// error.
func (p *Parser) parseStageOperand() (*ast.StageLocation, error) {
	if p.cur.Type != '@' {
		return nil, p.syntaxErrorAtCur()
	}
	return p.parseStageRef()
}

// parseOptionalPattern parses an optional trailing PATTERN = '<regex>' clause,
// returning the pattern literal (or nil if absent). Snowflake accepts both a
// single-quoted string and a double-quoted identifier for the pattern value
// (the GET docs example uses PATTERN = "tmp.parquet"), so a quoted-identifier
// token is accepted and normalized to a string literal carrying its text.
func (p *Parser) parseOptionalPattern() (*ast.Literal, error) {
	if p.cur.Type != kwPATTERN {
		return nil, nil
	}
	p.advance() // consume PATTERN
	if _, err := p.expect('='); err != nil {
		return nil, err
	}
	switch p.cur.Type {
	case tokString:
		tok := p.advance()
		return &ast.Literal{Kind: ast.LitString, Value: tok.Str, Loc: tok.Loc}, nil
	case tokQuotedIdent:
		// PATTERN = "regex": Snowflake permits a double-quoted value here; record
		// it as a string literal carrying the quoted text.
		tok := p.advance()
		return &ast.Literal{Kind: ast.LitString, Value: tok.Str, Loc: tok.Loc}, nil
	}
	return nil, p.syntaxErrorAtCur()
}
