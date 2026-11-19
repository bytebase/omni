package parser

import (
	"math"

	"github.com/bytebase/omni/starrocks/ast"
)

// Resource groups (StarRocks.g4:1358-1373). The statements open with
// CREATE, ALTER or DROP and the word RESOURCE, which is also the object of
// CREATE EXTERNAL RESOURCE; GROUP, reserved, tells the two apart.

// atResourceGroup reports whether the token after RESOURCE, which cur names,
// opens a resource group statement.
func (p *Parser) atResourceGroup() bool { return p.peekNext().Kind == kwGROUP }

// parseCreateResourceGroup parses createResourceGroupStatement
// (StarRocks.g4:1358):
//
//	CREATE RESOURCE GROUP [IF NOT EXISTS] [OR REPLACE] identifier
//	    [TO classifier [',' classifier]*] WITH '(' property [',' property]* ')'
//
// On entry CREATE and RESOURCE have been consumed; cur is GROUP.
func (p *Parser) parseCreateResourceGroup(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume GROUP

	stmt := &ast.CreateResourceGroupStmt{}
	if p.cur.Kind == kwIF {
		p.advance() // consume IF
		if _, err := p.expect(kwNOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfNotExists = true
	}
	if p.cur.Kind == kwOR {
		p.advance() // consume OR
		if _, err := p.expect(kwREPLACE); err != nil {
			return nil, err
		}
		stmt.OrReplace = true
	}

	name, nameLoc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = &ast.ObjectName{Parts: []string{name}, Loc: nameLoc}
	endLoc := nameLoc

	if p.cur.Kind == kwTO {
		p.advance() // consume TO
		if stmt.Classifiers, endLoc, err = p.parseClassifiers(); err != nil {
			return nil, err
		}
	}

	if _, err := p.expect(kwWITH); err != nil {
		return nil, err
	}
	props, err := p.parsePropertyList()
	if err != nil {
		return nil, err
	}
	stmt.Properties = props
	// parsePropertyList leaves p.prev at the closing ')', which ends the
	// statement even though the property before it starts later.
	endLoc = p.prev.Loc

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseClassifiers parses classifier (',' classifier)* where a classifier is
// '(' expressionList ')' (StarRocks.g4:1370).
func (p *Parser) parseClassifiers() ([][]ast.Node, ast.Loc, error) {
	var out [][]ast.Node
	var endLoc ast.Loc
	for {
		if _, err := p.expect(int('(')); err != nil {
			return nil, ast.NoLoc(), err
		}
		exprs, err := p.parseExprList()
		if err != nil {
			return nil, ast.NoLoc(), err
		}
		// Each classifier expression is a Predicate: the AST builder casts it
		// and refuses `TO(1)`, `TO('x')` and `TO(a)`.
		for _, expr := range exprs {
			if !isPredicateNode(expr) {
				return nil, ast.NoLoc(), p.syntaxErrorAtCur()
			}
		}
		closeTok, err := p.expect(int(')'))
		if err != nil {
			return nil, ast.NoLoc(), err
		}
		out = append(out, exprs)
		endLoc = closeTok.Loc
		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}
	return out, endLoc, nil
}

// parseDropResourceGroup parses dropResourceGroupStatement (StarRocks.g4:1363):
//
//	DROP RESOURCE GROUP identifier
//
// On entry DROP and RESOURCE have been consumed; cur is GROUP.
func (p *Parser) parseDropResourceGroup(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume GROUP
	name, nameLoc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	return &ast.DropResourceGroupStmt{
		Name: &ast.ObjectName{Parts: []string{name}, Loc: nameLoc},
		Loc:  startLoc.Merge(nameLoc),
	}, nil
}

// parseAlterResourceGroup parses alterResourceGroupStatement
// (StarRocks.g4:1367):
//
//	ALTER RESOURCE GROUP identifier ADD classifier [',' classifier]*
//	ALTER RESOURCE GROUP identifier DROP '(' INTEGER_VALUE [',' INTEGER_VALUE]* ')'
//	ALTER RESOURCE GROUP identifier DROP ALL
//	ALTER RESOURCE GROUP identifier WITH '(' property [',' property]* ')'
//
// On entry ALTER and RESOURCE have been consumed; cur is GROUP.
func (p *Parser) parseAlterResourceGroup(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume GROUP
	name, nameLoc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt := &ast.AlterResourceGroupStmt{
		Name: &ast.ObjectName{Parts: []string{name}, Loc: nameLoc},
	}
	endLoc := nameLoc

	switch p.cur.Kind {
	case kwADD:
		p.advance() // consume ADD
		stmt.Action = "ADD"
		if stmt.Classifiers, endLoc, err = p.parseClassifiers(); err != nil {
			return nil, err
		}
	case kwDROP:
		p.advance() // consume DROP
		stmt.Action = "DROP"
		if p.cur.Kind == kwALL {
			stmt.Action = "DROP ALL"
			endLoc = p.cur.Loc
			p.advance()
			break
		}
		if _, err := p.expect(int('(')); err != nil {
			return nil, err
		}
		for {
			if p.cur.Kind != tokInt {
				return nil, p.syntaxErrorAtCur()
			}
			endLoc = p.cur.Loc
			stmt.ClassIDs = append(stmt.ClassIDs, p.advance().Str)
			if p.cur.Kind != int(',') {
				break
			}
			p.advance() // consume ','
		}
		closeTok, err := p.expect(int(')'))
		if err != nil {
			return nil, err
		}
		endLoc = closeTok.Loc
	case kwWITH:
		p.advance() // consume WITH
		stmt.Action = "WITH"
		props, err := p.parsePropertyList()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
		endLoc = p.prev.Loc // the closing ')'
	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// atClearDataCacheRules reports whether a CLEAR at cur clears the data cache
// rules: CLEAR DATACACHE RULES.
func (p *Parser) atClearDataCacheRules() bool {
	return p.peekNext().Kind == kwDATACACHE
}

// parseCreateDataCacheRule parses createDataCacheRuleStatement
// (StarRocks.g4:1851):
//
//	CREATE DATACACHE RULE dataCacheTarget [(WHERE expression)]
//	    PRIORITY '=' ['-'] INTEGER_VALUE properties?
//
// where dataCacheTarget (:1847) is identifierOrStringOrStar '.' identifierOrStringOrStar
// '.' identifierOrStringOrStar.
//
// On entry CREATE has been consumed; cur is DATACACHE.
func (p *Parser) parseCreateDataCacheRule(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume DATACACHE
	if _, err := p.expect(kwRULE); err != nil {
		return nil, err
	}

	stmt := &ast.CreateDataCacheRuleStmt{}

	var parts []string
	for {
		part, partLoc, err := p.parseIdentifierOrStringOrStar()
		if err != nil {
			return nil, err
		}
		parts = append(parts, part)
		endLoc := partLoc
		if len(parts) == 3 {
			stmt.Loc = startLoc.Merge(endLoc)
			break
		}
		if _, err := p.expect(int('.')); err != nil {
			return nil, err
		}
	}
	stmt.Target = parts
	endLoc := stmt.Loc

	if p.cur.Kind == kwWHERE {
		p.advance() // consume WHERE
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Where = where
		endLoc = ast.NodeLoc(where)
	}

	if _, err := p.expect(kwPRIORITY); err != nil {
		return nil, err
	}
	if _, err := p.expect(int('=')); err != nil {
		return nil, err
	}
	if p.cur.Kind == int('-') {
		stmt.Negative = true
		endLoc = p.cur.Loc
		p.advance() // consume '-'
	}
	if p.cur.Kind != tokInt {
		return nil, p.syntaxErrorAtCur()
	}
	// The AST builder parses the value with Integer.parseInt, so a priority
	// outside a signed 32-bit integer is refused ("For input string:
	// \"2147483648\"").
	if p.cur.Ival > math.MaxInt32 || p.cur.Ival < math.MinInt32 {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Priority = p.cur.Str
	endLoc = p.cur.Loc
	p.advance()

	if p.cur.Kind == kwPROPERTIES {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
		endLoc = p.prev.Loc // the closing ')'
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseIdentifierOrStringOrStar reads one part of a dataCacheTarget: an
// identifier, a string, or the '*' wildcard.
func (p *Parser) parseIdentifierOrStringOrStar() (string, ast.Loc, error) {
	if p.cur.Kind == int('*') {
		tok := p.advance()
		return "*", tok.Loc, nil
	}
	return p.parseIdentifierOrString()
}

// parseDropDataCacheRule parses dropDataCacheRuleStatement (StarRocks.g4:1859):
//
//	DROP DATACACHE RULE INTEGER_VALUE
//
// On entry DROP has been consumed; cur is DATACACHE.
func (p *Parser) parseDropDataCacheRule(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume DATACACHE
	if _, err := p.expect(kwRULE); err != nil {
		return nil, err
	}
	if p.cur.Kind != tokInt {
		return nil, p.syntaxErrorAtCur()
	}
	idTok := p.advance()
	return &ast.DropDataCacheRuleStmt{ID: idTok.Str, Loc: startLoc.Merge(idTok.Loc)}, nil
}

// parseClearDataCacheRules parses clearDataCacheRulesStatement
// (StarRocks.g4:1863): CLEAR DATACACHE RULES.
//
// On entry CLEAR has been consumed; cur is DATACACHE.
func (p *Parser) parseClearDataCacheRules(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume DATACACHE
	rulesTok, err := p.expect(kwRULES)
	if err != nil {
		return nil, err
	}
	return &ast.ClearDataCacheRulesStmt{Loc: startLoc.Merge(rulesTok.Loc)}, nil
}
