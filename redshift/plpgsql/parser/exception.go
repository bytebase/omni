package parser

import (
	"strings"

	"github.com/bytebase/omni/redshift/plpgsql/ast"

	pgparser "github.com/bytebase/omni/redshift/parser"
)

// parseExceptionBlock parses the EXCEPTION section of a block:
//
//	EXCEPTION
//	  WHEN condition [OR condition ...] THEN stmts
//	  [WHEN condition [OR condition ...] THEN stmts ...]
//
// Conditions can be:
//   - Named conditions: division_by_zero, unique_violation, no_data_found, etc.
//   - SQLSTATE 'code': e.g. SQLSTATE '22012'
//   - OTHERS: catch-all
//
// Multiple conditions can be joined with OR.
func (p *Parser) parseExceptionBlock() ([]ast.Node, error) {
	// Consume EXCEPTION keyword (already verified by caller)
	p.advance()

	var whens []ast.Node

	for p.isKeyword("WHEN") {
		whenNode, err := p.parseExceptionWhen()
		if err != nil {
			return nil, err
		}
		whens = append(whens, whenNode)
	}

	if len(whens) == 0 {
		return nil, p.errorf("syntax error at or near %q, expected WHEN", p.tokenText(p.cur))
	}

	return whens, nil
}

// parseExceptionWhen parses a single WHEN clause in an EXCEPTION block:
//
//	WHEN condition [OR condition ...] THEN stmts
func (p *Parser) parseExceptionWhen() (*ast.PLExceptionWhen, error) {
	startPos := p.pos()
	p.advance() // consume WHEN

	// Parse condition list (joined by OR)
	conditions, err := p.parseExceptionConditions()
	if err != nil {
		return nil, err
	}

	// Expect THEN
	if err := p.expectKeyword("THEN"); err != nil {
		return nil, err
	}

	// Parse handler body
	body, err := p.parseStmtList()
	if err != nil {
		return nil, err
	}

	return &ast.PLExceptionWhen{
		Conditions: conditions,
		Body:       body,
		Loc:        ast.Loc{Start: startPos, End: p.prev.End},
	}, nil
}

// parseExceptionConditions parses a list of exception conditions joined by OR.
// Each condition is either:
//   - A named condition (identifier): division_by_zero, unique_violation, OTHERS, etc.
//   - SQLSTATE 'code': e.g. SQLSTATE '22012'
func (p *Parser) parseExceptionConditions() ([]string, error) {
	var conditions []string

	for {
		cond, err := p.parseExceptionCondition()
		if err != nil {
			return nil, err
		}
		conditions = append(conditions, cond)

		// Check for OR to continue
		if !p.isKeyword("OR") {
			break
		}
		p.advance() // consume OR
	}

	return conditions, nil
}

// parseExceptionCondition parses a single exception condition.
func (p *Parser) parseExceptionCondition() (string, error) {
	// Check for SQLSTATE 'code'
	if p.isKeyword("SQLSTATE") {
		p.advance() // consume SQLSTATE
		// Expect string constant for the SQLSTATE code
		if p.cur.Type != pgparser.SCONST {
			return "", p.errorf("syntax error at or near %q, expected SQLSTATE code string", p.tokenText(p.cur))
		}
		code := p.cur.Str
		p.advance()
		return "SQLSTATE '" + code + "'", nil
	}

	// Named condition (identifier) — OTHERS, division_by_zero, unique_violation, etc.
	if p.isIdent() || p.isAnyKeywordAsIdent() {
		name := strings.ToLower(p.identText())
		p.advance()
		return name, nil
	}

	return "", p.errorf("syntax error at or near %q, expected exception condition", p.tokenText(p.cur))
}
