package parser

import (
	"strings"

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

// parseWorkloadPolicyItemList parses a parenthesised list of items for
// CONDITIONS or ACTIONS. The list contents are captured as raw text because
// the individual condition/action syntax is complex and varied. Each comma-
// separated top-level item is returned as a separate WorkloadPolicyItem.
//
//	(item1, item2, ...)
//
// cur must be '(' on entry; it is consumed here.
func (p *Parser) parseWorkloadPolicyItemList() ([]*ast.WorkloadPolicyItem, error) {
	openTok, err := p.expect(int('('))
	if err != nil {
		return nil, err
	}
	_ = openTok

	var items []*ast.WorkloadPolicyItem
	var buf strings.Builder
	itemStart := p.cur.Loc

	depth := 0
	for p.cur.Kind != tokEOF {
		if p.cur.Kind == int(')') && depth == 0 {
			break
		}
		switch p.cur.Kind {
		case int('('):
			depth++
		case int(')'):
			depth--
		case int(','):
			if depth == 0 {
				// End of this item — save it and start a fresh one.
				raw := strings.TrimSpace(buf.String())
				items = append(items, &ast.WorkloadPolicyItem{
					RawText: raw,
					Loc:     ast.Loc{Start: itemStart.Start, End: p.prev.Loc.End},
				})
				buf.Reset()
				p.advance() // consume ','
				itemStart = p.cur.Loc
				continue
			}
		}
		if buf.Len() > 0 {
			buf.WriteByte(' ')
		}
		if p.cur.Str != "" {
			buf.WriteString(p.cur.Str)
		} else {
			buf.WriteString(TokenName(p.cur.Kind))
		}
		p.advance()
	}

	// Capture the last (or only) item before the closing ')'.
	raw := strings.TrimSpace(buf.String())
	if raw != "" {
		items = append(items, &ast.WorkloadPolicyItem{
			RawText: raw,
			Loc:     ast.Loc{Start: itemStart.Start, End: p.cur.Loc.Start},
		})
	}

	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	return items, nil
}

// parseCreateWorkloadPolicy parses:
//
//	CREATE WORKLOAD POLICY [IF NOT EXISTS] name
//	    CONDITIONS(condition_list)
//	    ACTIONS(action_list)
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
		conds, err := p.parseWorkloadPolicyItemList()
		if err != nil {
			return nil, err
		}
		stmt.Conditions = conds
	}

	// Optional ACTIONS(...)
	if p.cur.Kind == kwACTIONS {
		p.advance() // consume ACTIONS
		actions, err := p.parseWorkloadPolicyItemList()
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
