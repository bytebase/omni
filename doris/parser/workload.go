package parser

import (
	"github.com/bytebase/omni/doris/ast"
)

// ---------------------------------------------------------------------------
// WORKLOAD GROUP
// ---------------------------------------------------------------------------

// parseCreateWorkloadGroup parses:
//
//	CREATE WORKLOAD GROUP [IF NOT EXISTS] name PROPERTIES(...)
//
// CREATE has already been consumed; cur is WORKLOAD. The WORKLOAD and GROUP
// tokens are consumed here.
func (p *Parser) parseCreateWorkloadGroup(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume WORKLOAD
	if _, err := p.expect(kwGROUP); err != nil {
		return nil, err
	}

	stmt := &ast.CreateWorkloadGroupStmt{}

	// Optional IF NOT EXISTS
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

	// Group name
	name, nameLoc, err := p.parseIdentifierOrString()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	endLoc := nameLoc

	// Optional PROPERTIES clause
	if p.cur.Kind == kwPROPERTIES {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
		if len(props) > 0 {
			endLoc = ast.NodeLoc(props[len(props)-1])
		}
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseAlterWorkloadGroup parses:
//
//	ALTER WORKLOAD GROUP name PROPERTIES(...)
//
// ALTER has already been consumed; cur is WORKLOAD. The WORKLOAD and GROUP
// tokens are consumed here.
func (p *Parser) parseAlterWorkloadGroup(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume WORKLOAD
	if _, err := p.expect(kwGROUP); err != nil {
		return nil, err
	}

	stmt := &ast.AlterWorkloadGroupStmt{}

	// Group name
	name, nameLoc, err := p.parseIdentifierOrString()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	endLoc := nameLoc

	// PROPERTIES clause
	if p.cur.Kind == kwPROPERTIES {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
		if len(props) > 0 {
			endLoc = ast.NodeLoc(props[len(props)-1])
		}
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseDropWorkloadGroup parses:
//
//	DROP WORKLOAD GROUP [IF EXISTS] name
//
// DROP has already been consumed; cur is WORKLOAD. The WORKLOAD and GROUP
// tokens are consumed here.
func (p *Parser) parseDropWorkloadGroup(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume WORKLOAD
	if _, err := p.expect(kwGROUP); err != nil {
		return nil, err
	}

	stmt := &ast.DropWorkloadGroupStmt{}

	// Optional IF EXISTS
	if p.cur.Kind == kwIF {
		p.advance() // consume IF
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}

	// Group name
	name, nameLoc, err := p.parseIdentifierOrString()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc = startLoc.Merge(nameLoc)
	return stmt, nil
}

// ---------------------------------------------------------------------------
// WORKLOAD POLICY
// ---------------------------------------------------------------------------

// parseWorkloadPolicyConditions reads CONDITIONS' list from its '(' through
// its ')' as both of Doris's parsers take it: DorisParser.g4's
//
//	workloadPolicyCondition : identifier comparisonOperator (number | STRING_LITERAL)
//
// (:852) and the CUP parser's workload_policy_condition_list
// (DorisLegacyParser.g4:1198), whose operator is =, >, <, > = or < = and
// whose value is an unsigned number or a string. The CUP lexer has no >= or
// <= token, so the two characters of either are its > = and < =; written
// apart, ANTLR refuses them.
func (p *Parser) parseWorkloadPolicyConditions() ([]*ast.WorkloadPolicyItem, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	var items []*ast.WorkloadPolicyItem
	for {
		start := p.cur.Loc
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		var op string
		switch p.cur.Kind {
		case int('='), int('<'), int('>'):
			op = string(rune(p.cur.Kind))
		case tokLessEq:
			op = "<="
		case tokGreaterEq:
			op = ">="
		default:
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		var value string
		switch p.cur.Kind {
		case tokString:
			value = p.cur.Str
		case tokInt, tokFloat:
			value = p.sourceText(p.cur.Loc)
		default:
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		items = append(items, p.workloadPolicyItem(name, op, value, start))
		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	return items, nil
}

// parseWorkloadPolicyActions reads ACTIONS' list from its '(' through its
// ')': DorisParser.g4's
//
//	workloadPolicyAction : SET_SESSION_VARIABLE STRING_LITERAL | identifier (STRING_LITERAL)?
//
// (:843), as the CUP parser's workload_policy_action_list
// (DorisLegacyParser.g4:1215) also takes it. Both take SET_SESSION_VARIABLE
// alone, as a name (testdata/compat.tsv).
func (p *Parser) parseWorkloadPolicyActions() ([]*ast.WorkloadPolicyItem, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	var items []*ast.WorkloadPolicyItem
	for {
		start := p.cur.Loc
		var name string
		if p.cur.Kind == kwSET_SESSION_VARIABLE {
			name = p.advance().Str
		} else {
			n, _, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			name = n
		}
		var value string
		if p.cur.Kind == tokString {
			value = p.advance().Str
		}
		items = append(items, p.workloadPolicyItem(name, "", value, start))
		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	return items, nil
}

// workloadPolicyItem builds an item that began at start and ends at the
// last token read.
func (p *Parser) workloadPolicyItem(name, op, value string, start ast.Loc) *ast.WorkloadPolicyItem {
	loc := ast.Loc{Start: start.Start, End: p.prev.Loc.End}
	return &ast.WorkloadPolicyItem{Name: name, Op: op, Value: value, RawText: p.sourceText(loc), Loc: loc}
}

// parseCreateWorkloadPolicy parses (DorisParser.g4:825):
//
//	CREATE WORKLOAD POLICY [IF NOT EXISTS] name
//	    [CONDITIONS (condition, ...)]
//	    [ACTIONS (action, ...)]
//	    [PROPERTIES(...)]
//
// CREATE has already been consumed; cur is WORKLOAD. The WORKLOAD and POLICY
// tokens are consumed here.
func (p *Parser) parseCreateWorkloadPolicy(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume WORKLOAD
	if _, err := p.expect(kwPOLICY); err != nil {
		return nil, err
	}

	stmt := &ast.CreateWorkloadPolicyStmt{}

	// Optional IF NOT EXISTS
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

	// Policy name
	name, nameLoc, err := p.parseIdentifierOrString()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	endLoc := nameLoc

	// Optional CONDITIONS(...)
	if p.cur.Kind == kwCONDITIONS {
		p.advance() // consume CONDITIONS
		conds, err := p.parseWorkloadPolicyConditions()
		if err != nil {
			return nil, err
		}
		stmt.Conditions = conds
	}

	// Optional ACTIONS(...)
	if p.cur.Kind == kwACTIONS {
		p.advance() // consume ACTIONS
		actions, err := p.parseWorkloadPolicyActions()
		if err != nil {
			return nil, err
		}
		stmt.Actions = actions
	}

	// Optional PROPERTIES clause
	if p.cur.Kind == kwPROPERTIES {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
		if len(props) > 0 {
			endLoc = ast.NodeLoc(props[len(props)-1])
		}
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseAlterWorkloadPolicy parses:
//
//	ALTER WORKLOAD POLICY name PROPERTIES(...)
//
// ALTER has already been consumed; cur is WORKLOAD. The WORKLOAD and POLICY
// tokens are consumed here.
func (p *Parser) parseAlterWorkloadPolicy(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume WORKLOAD
	if _, err := p.expect(kwPOLICY); err != nil {
		return nil, err
	}

	stmt := &ast.AlterWorkloadPolicyStmt{}

	// Policy name
	name, nameLoc, err := p.parseIdentifierOrString()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	endLoc := nameLoc

	// PROPERTIES clause
	if p.cur.Kind == kwPROPERTIES {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
		if len(props) > 0 {
			endLoc = ast.NodeLoc(props[len(props)-1])
		}
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseDropWorkloadPolicy parses:
//
//	DROP WORKLOAD POLICY [IF EXISTS] name
//
// DROP has already been consumed; cur is WORKLOAD. The WORKLOAD and POLICY
// tokens are consumed here.
func (p *Parser) parseDropWorkloadPolicy(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume WORKLOAD
	if _, err := p.expect(kwPOLICY); err != nil {
		return nil, err
	}

	stmt := &ast.DropWorkloadPolicyStmt{}

	// Optional IF EXISTS
	if p.cur.Kind == kwIF {
		p.advance() // consume IF
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}

	// Policy name
	name, nameLoc, err := p.parseIdentifierOrString()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc = startLoc.Merge(nameLoc)
	return stmt, nil
}

// ---------------------------------------------------------------------------
// RESOURCE
// ---------------------------------------------------------------------------

// parseCreateResource parses:
//
//	CREATE [EXTERNAL] RESOURCE [IF NOT EXISTS] name PROPERTIES(...)
//
// CREATE has already been consumed; cur is RESOURCE (for non-EXTERNAL) or
// EXTERNAL (the caller already peeked). The RESOURCE (and optional EXTERNAL)
// token(s) are consumed here. external indicates whether EXTERNAL was already
// seen/consumed by the caller.
func (p *Parser) parseCreateResource(startLoc ast.Loc, external bool) (ast.Node, error) {
	// cur is RESOURCE
	if _, err := p.expect(kwRESOURCE); err != nil {
		return nil, err
	}

	stmt := &ast.CreateResourceStmt{External: external}

	// Optional IF NOT EXISTS
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

	// Resource name — can be a string literal or identifier
	name, nameLoc, err := p.parseIdentifierOrString()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	endLoc := nameLoc

	// PROPERTIES clause
	if p.cur.Kind == kwPROPERTIES {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
		if len(props) > 0 {
			endLoc = ast.NodeLoc(props[len(props)-1])
		}
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseAlterResource parses:
//
//	ALTER RESOURCE name PROPERTIES(...)
//
// ALTER has already been consumed; cur is RESOURCE. RESOURCE is consumed here.
func (p *Parser) parseAlterResource(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume RESOURCE

	stmt := &ast.AlterResourceStmt{}

	// Resource name
	name, nameLoc, err := p.parseIdentifierOrString()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	endLoc := nameLoc

	// PROPERTIES clause
	if p.cur.Kind == kwPROPERTIES {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
		if len(props) > 0 {
			endLoc = ast.NodeLoc(props[len(props)-1])
		}
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseDropResource parses:
//
//	DROP RESOURCE [IF EXISTS] name
//
// DROP has already been consumed; cur is RESOURCE. RESOURCE is consumed here.
func (p *Parser) parseDropResource(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume RESOURCE

	stmt := &ast.DropResourceStmt{}

	// Optional IF EXISTS
	if p.cur.Kind == kwIF {
		p.advance() // consume IF
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}

	// Resource name
	name, nameLoc, err := p.parseIdentifierOrString()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc = startLoc.Merge(nameLoc)
	return stmt, nil
}

// ---------------------------------------------------------------------------
// SQL BLOCK RULE
// ---------------------------------------------------------------------------

// parseCreateSQLBlockRule parses:
//
//	CREATE SQL_BLOCK_RULE [IF NOT EXISTS] name PROPERTIES(...)
//
// CREATE has already been consumed; cur is SQL_BLOCK_RULE. SQL_BLOCK_RULE is
// consumed here.
func (p *Parser) parseCreateSQLBlockRule(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume SQL_BLOCK_RULE

	stmt := &ast.CreateSQLBlockRuleStmt{}

	// Optional IF NOT EXISTS
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

	// Rule name
	name, nameLoc, err := p.parseIdentifierOrString()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	endLoc := nameLoc

	// PROPERTIES clause
	if p.cur.Kind == kwPROPERTIES {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
		if len(props) > 0 {
			endLoc = ast.NodeLoc(props[len(props)-1])
		}
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseAlterSQLBlockRule parses:
//
//	ALTER SQL_BLOCK_RULE name PROPERTIES(...)
//
// ALTER has already been consumed; cur is SQL_BLOCK_RULE. SQL_BLOCK_RULE is
// consumed here.
func (p *Parser) parseAlterSQLBlockRule(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume SQL_BLOCK_RULE

	stmt := &ast.AlterSQLBlockRuleStmt{}

	// Rule name
	name, nameLoc, err := p.parseIdentifierOrString()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	endLoc := nameLoc

	// PROPERTIES clause
	if p.cur.Kind == kwPROPERTIES {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
		if len(props) > 0 {
			endLoc = ast.NodeLoc(props[len(props)-1])
		}
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseDropSQLBlockRule parses:
//
//	DROP SQL_BLOCK_RULE [IF EXISTS] name
//
// DROP has already been consumed; cur is SQL_BLOCK_RULE. SQL_BLOCK_RULE is
// consumed here.
func (p *Parser) parseDropSQLBlockRule(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume SQL_BLOCK_RULE

	stmt := &ast.DropSQLBlockRuleStmt{}

	// Optional IF EXISTS
	if p.cur.Kind == kwIF {
		p.advance() // consume IF
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}

	// Rule name
	name, nameLoc, err := p.parseIdentifierOrString()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc = startLoc.Merge(nameLoc)
	return stmt, nil
}
