package parser

import "fmt"

// ResourcePlanParser.g statements: RESOURCE PLAN, WORKLOAD MANAGEMENT,
// TRIGGER, POOL, and MAPPING DDL.

func (p *parser) parseCreateResourcePlanStatement() error {
	if _, err := p.expect(KW_CREATE); err != nil {
		return err
	}
	if _, err := p.expect(KW_RESOURCE); err != nil {
		return err
	}
	if _, err := p.expect(KW_PLAN); err != nil {
		return err
	}
	p.try(p.parseIfNotExists)
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if p.match(KW_LIKE) {
		_, err := p.parseIdentifier()
		return err
	}
	if p.match(KW_WITH) {
		return p.parseRPAssignList()
	}
	return nil
}

func (p *parser) parseAlterResourcePlanStatement() error {
	// KW_ALTER already consumed by parseAlterDDL.
	if _, err := p.expect(KW_RESOURCE); err != nil {
		return err
	}
	if _, err := p.expect(KW_PLAN); err != nil {
		return err
	}
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	switch p.peek().Type {
	case KW_VALIDATE, KW_DISABLE:
		p.next()
		return nil
	case KW_SET:
		p.next()
		return p.parseRPAssignList()
	case KW_UNSET:
		p.next()
		return p.parseRPUnassignList()
	case KW_RENAME:
		p.next()
		if _, err := p.expect(KW_TO); err != nil {
			return err
		}
		_, err := p.parseIdentifier()
		return err
	case KW_ACTIVATE:
		if err := p.parseActivate(); err != nil {
			return err
		}
		p.match(KW_ENABLE)
		return nil
	case KW_ENABLE:
		p.next()
		if p.peek().Type == KW_ACTIVATE {
			return p.parseActivate()
		}
		return nil
	default:
		return p.errAt(p.peek(), fmt.Sprintf("unexpected token %q in ALTER RESOURCE PLAN", p.peek().Str))
	}
}

func (p *parser) parseActivate() error {
	if _, err := p.expect(KW_ACTIVATE); err != nil {
		return err
	}
	if p.match(KW_WITH) {
		_, err := p.expect(KW_REPLACE)
		return err
	}
	return nil
}

func (p *parser) parseGlobalWmStatement() error {
	if p.peek().Type != KW_ENABLE && p.peek().Type != KW_DISABLE {
		return p.errAt(p.peek(), "expected ENABLE or DISABLE")
	}
	p.next()
	if _, err := p.expect(KW_WORKLOAD); err != nil {
		return err
	}
	_, err := p.expect(KW_MANAGEMENT)
	return err
}

func (p *parser) parseReplaceResourcePlanStatement() error {
	if _, err := p.expect(KW_REPLACE); err != nil {
		return err
	}
	if p.match(KW_ACTIVE) {
		if _, err := p.expect(KW_RESOURCE); err != nil {
			return err
		}
		if _, err := p.expect(KW_PLAN); err != nil {
			return err
		}
		if _, err := p.expect(KW_WITH); err != nil {
			return err
		}
		_, err := p.parseIdentifier()
		return err
	}
	if _, err := p.expect(KW_RESOURCE); err != nil {
		return err
	}
	if _, err := p.expect(KW_PLAN); err != nil {
		return err
	}
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if _, err := p.expect(KW_WITH); err != nil {
		return err
	}
	_, err := p.parseIdentifier()
	return err
}

func (p *parser) parseDropResourcePlanStatement() error {
	if _, err := p.expect(KW_DROP); err != nil {
		return err
	}
	if _, err := p.expect(KW_RESOURCE); err != nil {
		return err
	}
	if _, err := p.expect(KW_PLAN); err != nil {
		return err
	}
	p.try(p.parseIfExists)
	_, err := p.parseIdentifier()
	return err
}

func (p *parser) parseRPAssignList() error {
	if err := p.parseRPAssign(); err != nil {
		return err
	}
	for p.match(TokComma) {
		if err := p.parseRPAssign(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseRPAssign() error {
	if p.match(KW_QUERY_PARALLELISM) {
		if _, err := p.expect(TokEqual); err != nil {
			return err
		}
		_, err := p.expect(TokNumber)
		return err
	}
	if _, err := p.expect(KW_DEFAULT); err != nil {
		return err
	}
	if _, err := p.expect(KW_POOL); err != nil {
		return err
	}
	if _, err := p.expect(TokEqual); err != nil {
		return err
	}
	return p.parsePoolPath()
}

func (p *parser) parseRPUnassignList() error {
	if err := p.parseRPUnassign(); err != nil {
		return err
	}
	for p.match(TokComma) {
		if err := p.parseRPUnassign(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseRPUnassign() error {
	if p.match(KW_QUERY_PARALLELISM) {
		return nil
	}
	if _, err := p.expect(KW_DEFAULT); err != nil {
		return err
	}
	_, err := p.expect(KW_POOL)
	return err
}

func (p *parser) parsePoolPath() error {
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	for p.match(TokDot) {
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseTriggerExpression() error {
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if _, err := p.expect(TokGreaterThan); err != nil {
		return err
	}
	if p.peek().Type == TokNumber || p.peek().Type == TokStringLiteral {
		p.next()
		return nil
	}
	return p.errAt(p.peek(), fmt.Sprintf("expected number or string literal in trigger expression, got %q", p.peek().Str))
}

func (p *parser) parseTriggerActionExpression() error {
	if p.match(KW_KILL) {
		return nil
	}
	if _, err := p.expect(KW_MOVE); err != nil {
		return err
	}
	if _, err := p.expect(KW_TO); err != nil {
		return err
	}
	return p.parsePoolPath()
}

func (p *parser) parseCreateTriggerStatement() error {
	if _, err := p.expect(KW_CREATE); err != nil {
		return err
	}
	if _, err := p.expect(KW_TRIGGER); err != nil {
		return err
	}
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if _, err := p.expect(TokDot); err != nil {
		return err
	}
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if _, err := p.expect(KW_WHEN); err != nil {
		return err
	}
	if err := p.parseTriggerExpression(); err != nil {
		return err
	}
	if _, err := p.expect(KW_DO); err != nil {
		return err
	}
	return p.parseTriggerActionExpression()
}

func (p *parser) parseAlterTriggerStatement() error {
	// KW_ALTER already consumed by parseAlterDDL.
	if _, err := p.expect(KW_TRIGGER); err != nil {
		return err
	}
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if _, err := p.expect(TokDot); err != nil {
		return err
	}
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	switch p.peek().Type {
	case KW_WHEN:
		p.next()
		if err := p.parseTriggerExpression(); err != nil {
			return err
		}
		if _, err := p.expect(KW_DO); err != nil {
			return err
		}
		return p.parseTriggerActionExpression()
	case KW_ADD:
		p.next()
		if _, err := p.expect(KW_TO); err != nil {
			return err
		}
		if p.match(KW_UNMANAGED) {
			return nil
		}
		if _, err := p.expect(KW_POOL); err != nil {
			return err
		}
		return p.parsePoolPath()
	case KW_DROP:
		p.next()
		if _, err := p.expect(KW_FROM); err != nil {
			return err
		}
		if p.match(KW_UNMANAGED) {
			return nil
		}
		if _, err := p.expect(KW_POOL); err != nil {
			return err
		}
		return p.parsePoolPath()
	default:
		return p.errAt(p.peek(), fmt.Sprintf("unexpected token %q in ALTER TRIGGER", p.peek().Str))
	}
}

func (p *parser) parseDropTriggerStatement() error {
	if _, err := p.expect(KW_DROP); err != nil {
		return err
	}
	if _, err := p.expect(KW_TRIGGER); err != nil {
		return err
	}
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if _, err := p.expect(TokDot); err != nil {
		return err
	}
	_, err := p.parseIdentifier()
	return err
}

func (p *parser) parsePoolAssignList() error {
	if err := p.parsePoolAssign(); err != nil {
		return err
	}
	for p.match(TokComma) {
		if err := p.parsePoolAssign(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parsePoolAssign() error {
	switch p.peek().Type {
	case KW_ALLOC_FRACTION, KW_QUERY_PARALLELISM:
		p.next()
		if _, err := p.expect(TokEqual); err != nil {
			return err
		}
		_, err := p.expect(TokNumber)
		return err
	case KW_SCHEDULING_POLICY:
		p.next()
		if _, err := p.expect(TokEqual); err != nil {
			return err
		}
		_, err := p.expect(TokStringLiteral)
		return err
	case KW_PATH:
		p.next()
		if _, err := p.expect(TokEqual); err != nil {
			return err
		}
		return p.parsePoolPath()
	default:
		return p.errAt(p.peek(), fmt.Sprintf("unexpected pool assignment %q", p.peek().Str))
	}
}

func (p *parser) parseCreatePoolStatement() error {
	if _, err := p.expect(KW_CREATE); err != nil {
		return err
	}
	if _, err := p.expect(KW_POOL); err != nil {
		return err
	}
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if _, err := p.expect(TokDot); err != nil {
		return err
	}
	if err := p.parsePoolPath(); err != nil {
		return err
	}
	if _, err := p.expect(KW_WITH); err != nil {
		return err
	}
	return p.parsePoolAssignList()
}

func (p *parser) parseAlterPoolStatement() error {
	// KW_ALTER already consumed by parseAlterDDL.
	if _, err := p.expect(KW_POOL); err != nil {
		return err
	}
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if _, err := p.expect(TokDot); err != nil {
		return err
	}
	if err := p.parsePoolPath(); err != nil {
		return err
	}
	switch p.peek().Type {
	case KW_SET:
		p.next()
		return p.parsePoolAssignList()
	case KW_UNSET:
		p.next()
		_, err := p.expect(KW_SCHEDULING_POLICY)
		return err
	case KW_ADD, KW_DROP:
		p.next()
		if _, err := p.expect(KW_TRIGGER); err != nil {
			return err
		}
		_, err := p.parseIdentifier()
		return err
	default:
		return p.errAt(p.peek(), fmt.Sprintf("unexpected token %q in ALTER POOL", p.peek().Str))
	}
}

func (p *parser) parseDropPoolStatement() error {
	if _, err := p.expect(KW_DROP); err != nil {
		return err
	}
	if _, err := p.expect(KW_POOL); err != nil {
		return err
	}
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if _, err := p.expect(TokDot); err != nil {
		return err
	}
	return p.parsePoolPath()
}

func (p *parser) parseCreateMappingStatement() error {
	if _, err := p.expect(KW_CREATE); err != nil {
		return err
	}
	return p.parseMappingBody()
}

func (p *parser) parseAlterMappingStatement() error {
	// KW_ALTER already consumed by parseAlterDDL.
	return p.parseMappingBody()
}

func (p *parser) parseMappingBody() error {
	switch p.peek().Type {
	case KW_USER, KW_GROUP, KW_APPLICATION:
		p.next()
	default:
		return p.errAt(p.peek(), fmt.Sprintf("expected USER, GROUP, or APPLICATION, got %q", p.peek().Str))
	}
	if _, err := p.expect(KW_MAPPING); err != nil {
		return err
	}
	if _, err := p.expect(TokStringLiteral); err != nil {
		return err
	}
	if _, err := p.expect(KW_IN); err != nil {
		return err
	}
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if p.match(KW_TO) {
		if err := p.parsePoolPath(); err != nil {
			return err
		}
	} else if !p.match(KW_UNMANAGED) {
		return p.errAt(p.peek(), fmt.Sprintf("expected TO or UNMANAGED in MAPPING, got %q", p.peek().Str))
	}
	if p.match(KW_WITH) {
		if _, err := p.expect(KW_ORDER); err != nil {
			return err
		}
		_, err := p.expect(TokNumber)
		return err
	}
	return nil
}

func (p *parser) parseDropMappingStatement() error {
	if _, err := p.expect(KW_DROP); err != nil {
		return err
	}
	switch p.peek().Type {
	case KW_USER, KW_GROUP, KW_APPLICATION:
		p.next()
	default:
		return p.errAt(p.peek(), fmt.Sprintf("expected USER, GROUP, or APPLICATION, got %q", p.peek().Str))
	}
	if _, err := p.expect(KW_MAPPING); err != nil {
		return err
	}
	if _, err := p.expect(TokStringLiteral); err != nil {
		return err
	}
	if _, err := p.expect(KW_IN); err != nil {
		return err
	}
	_, err := p.parseIdentifier()
	return err
}
