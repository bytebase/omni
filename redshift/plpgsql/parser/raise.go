package parser

import (
	"strings"

	"github.com/bytebase/omni/redshift/plpgsql/ast"

	pgparser "github.com/bytebase/omni/redshift/parser"
)

// --------------------------------------------------------------------------
// Section 5.1: RAISE
// --------------------------------------------------------------------------

// parseRaise parses a RAISE statement.
//
// Ref: https://www.postgresql.org/docs/17/plpgsql-errors-and-messages.html
//
//	RAISE ;
//	RAISE level 'format' [, param ...] [USING option = expr, ...] ;
//	RAISE level condition_name [USING ...] ;
//	RAISE level SQLSTATE 'code' [USING ...] ;
//	RAISE [level] USING option = expr [, ...] ;
//	RAISE 'format' [, param ...] [USING ...] ;
//
// Levels: DEBUG, LOG, INFO, NOTICE, WARNING, EXCEPTION (default)
// USING options: MESSAGE, DETAIL, HINT, ERRCODE, COLUMN, CONSTRAINT, DATATYPE, TABLE, SCHEMA
func (p *Parser) parseRaise() (ast.Node, error) {
	startPos := p.pos()
	p.advance() // consume RAISE

	node := &ast.PLRaise{}

	// Bare RAISE (re-raise): RAISE ;
	if p.cur.Type == ';' {
		p.advance()
		node.Loc = ast.Loc{Start: startPos, End: p.prev.End}
		return node, nil
	}

	// Try to parse an optional level keyword
	level, hasLevel := p.tryRaiseLevel()
	if hasLevel {
		node.Level = level
	}

	// After level (or with no level), check what follows:
	// 1. semicolon -> done (bare raise with level, or RAISE ; already handled)
	// 2. USING -> USING-only form
	// 3. SQLSTATE -> SQLSTATE form
	// 4. string literal -> format string
	// 5. identifier -> condition name

	if p.cur.Type == ';' {
		p.advance()
		node.Loc = ast.Loc{Start: startPos, End: p.prev.End}
		return node, nil
	}

	if p.isKeyword("USING") {
		// RAISE [level] USING option = expr, ...
		opts, err := p.parseRaiseUsing()
		if err != nil {
			return nil, err
		}
		node.Options = opts
	} else if p.isKeyword("SQLSTATE") {
		// RAISE [level] SQLSTATE 'code' [USING ...]
		p.advance() // consume SQLSTATE
		if p.cur.Type != pgparser.SCONST {
			return nil, p.errorf("syntax error at or near %q, expected SQLSTATE code string", p.tokenText(p.cur))
		}
		node.SQLState = p.cur.Str
		p.advance()
		// Optional USING
		if p.isKeyword("USING") {
			opts, err := p.parseRaiseUsing()
			if err != nil {
				return nil, err
			}
			node.Options = opts
		}
	} else if p.cur.Type == pgparser.SCONST {
		// RAISE [level] 'format' [, param ...] [USING ...]
		node.Message = p.cur.Str
		p.advance()
		// Parse optional comma-separated parameters
		for p.cur.Type == ',' {
			p.advance()
			param, err := p.collectRaiseParam()
			if err != nil {
				return nil, err
			}
			node.Params = append(node.Params, param)
		}
		// Optional USING
		if p.isKeyword("USING") {
			opts, err := p.parseRaiseUsing()
			if err != nil {
				return nil, err
			}
			node.Options = opts
		}
	} else if p.isIdent() || p.isAnyKeywordAsIdent() {
		// RAISE [level] condition_name [USING ...]
		// condition_name is an identifier (e.g., division_by_zero, unique_violation)
		node.CondName = p.identText()
		p.advance()
		// Optional USING
		if p.isKeyword("USING") {
			opts, err := p.parseRaiseUsing()
			if err != nil {
				return nil, err
			}
			node.Options = opts
		}
	}

	// Expect trailing semicolon
	if _, err := p.expect(';'); err != nil {
		return nil, err
	}

	node.Loc = ast.Loc{Start: startPos, End: p.prev.End}
	return node, nil
}

// tryRaiseLevel checks if the current token is a RAISE level keyword.
// If so, it consumes the token and returns the level and true.
// Otherwise returns (0, false) without consuming anything.
func (p *Parser) tryRaiseLevel() (ast.RaiseLevel, bool) {
	if p.isKeyword("DEBUG") {
		p.advance()
		return ast.RaiseLevelDebug, true
	}
	if p.isKeyword("LOG") {
		p.advance()
		return ast.RaiseLevelLog, true
	}
	if p.isKeyword("INFO") {
		p.advance()
		return ast.RaiseLevelInfo, true
	}
	if p.isKeyword("NOTICE") {
		p.advance()
		return ast.RaiseLevelNotice, true
	}
	if p.isKeyword("WARNING") {
		p.advance()
		return ast.RaiseLevelWarning, true
	}
	if p.isKeyword("EXCEPTION") {
		p.advance()
		return ast.RaiseLevelError, true
	}
	return ast.RaiseLevelNone, false
}

// parseRaiseUsing parses: USING option = expr [, option = expr ...]
// The USING keyword must be the current token.
func (p *Parser) parseRaiseUsing() ([]ast.Node, error) {
	p.advance() // consume USING

	var opts []ast.Node
	for {
		opt, err := p.parseRaiseOption()
		if err != nil {
			return nil, err
		}
		opts = append(opts, opt)

		if p.cur.Type == ',' {
			p.advance()
			continue
		}
		break
	}

	return opts, nil
}

// raiseOptionNames are the valid USING option names for RAISE.
var raiseOptionNames = []string{
	"MESSAGE", "DETAIL", "HINT", "ERRCODE",
	"COLUMN", "CONSTRAINT", "DATATYPE", "TABLE", "SCHEMA",
}

// parseRaiseOption parses: option_name = expression
func (p *Parser) parseRaiseOption() (*ast.PLRaiseOption, error) {
	startPos := p.pos()

	// Match option name
	var optName string
	for _, name := range raiseOptionNames {
		if p.isKeyword(name) {
			optName = strings.ToUpper(name)
			p.advance()
			break
		}
	}
	if optName == "" {
		return nil, p.errorf("syntax error at or near %q, expected RAISE USING option name", p.tokenText(p.cur))
	}

	// Expect =
	if p.cur.Type != '=' {
		return nil, p.errorf("syntax error at or near %q, expected =", p.tokenText(p.cur))
	}
	p.advance()

	// Collect expression until comma, semicolon, or end
	expr, err := p.collectRaiseOptionExpr()
	if err != nil {
		return nil, err
	}

	return &ast.PLRaiseOption{
		OptType: optName,
		Expr:    expr,
		Loc:     ast.Loc{Start: startPos, End: p.prev.End},
	}, nil
}

// collectRaiseParam collects a single RAISE format parameter expression.
// Parameters are separated by commas, terminated by USING or semicolon.
func (p *Parser) collectRaiseParam() (string, error) {
	start := p.pos()
	depth := 0

	for !p.isEOF() {
		if p.cur.Type == '(' {
			depth++
			p.advance()
			continue
		}
		if p.cur.Type == ')' {
			if depth > 0 {
				depth--
			}
			p.advance()
			continue
		}
		if depth == 0 {
			// Stop at comma (next param), USING, or semicolon
			if p.cur.Type == ',' || p.cur.Type == ';' || p.isKeyword("USING") {
				text := strings.TrimSpace(p.source[start:p.cur.Loc])
				if text == "" {
					return "", p.errorf("syntax error at or near %q, expected expression", p.tokenText(p.cur))
				}
				return text, nil
			}
		}
		p.advance()
	}

	text := strings.TrimSpace(p.source[start:p.pos()])
	return text, nil
}

// collectRaiseOptionExpr collects a RAISE USING option value expression.
// Terminated by comma (next option), semicolon, or EOF.
func (p *Parser) collectRaiseOptionExpr() (string, error) {
	start := p.pos()
	depth := 0

	for !p.isEOF() {
		if p.cur.Type == '(' {
			depth++
			p.advance()
			continue
		}
		if p.cur.Type == ')' {
			if depth > 0 {
				depth--
			}
			p.advance()
			continue
		}
		if depth == 0 {
			if p.cur.Type == ',' || p.cur.Type == ';' {
				text := strings.TrimSpace(p.source[start:p.cur.Loc])
				if text == "" {
					return "", p.errorf("syntax error at or near %q, expected expression", p.tokenText(p.cur))
				}
				return text, nil
			}
		}
		p.advance()
	}

	text := strings.TrimSpace(p.source[start:p.pos()])
	return text, nil
}
