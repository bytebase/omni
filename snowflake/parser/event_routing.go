package parser

import "github.com/bytebase/omni/snowflake/ast"

// ---------------------------------------------------------------------------
// EVENT ROUTING TABLE (an organization's event routing rules), read as an
// ast.ObjectStmt
// ---------------------------------------------------------------------------

func init() {
	registerObject(&objectForm{
		words:  []string{"EVENT", "ROUTING", "TABLE"},
		create: parseCreateEventRoutingTable,
		alter:  parseAlterEventRoutingTable,
		// The CREATE and ALTER pages write no IF [ NOT ] EXISTS.
		noIfNotExists: true,
		noIfExists:    true,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-event-routing-table
		// writes no IF EXISTS.
		drop:           true,
		dropNoIfExists: true,
	})
}

// parseCreateEventRoutingTable reads what follows the name of
//
//	CREATE EVENT ROUTING TABLE <table_name>
//	  WITH RULES
//	    <rule_name> = ( REGION_GROUP = <region_group>, REGIONS = ( '<region>' [ , ... ] ),
//	                    DESTINATION_ACCOUNT = <organization>.<account_name> )
//	    ...
//
// https://docs.snowflake.com/en/sql-reference/sql/create-event-routing-table
//
// The page writes the rules one after another, with no comma between them.
// The WITH RULES clause holds each rule as an option named by the rule whose
// Group holds its three settings.
func parseCreateEventRoutingTable(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	if _, err := p.expect(kwWITH); err != nil {
		return err
	}
	if err := p.expectWords("RULES"); err != nil {
		return err
	}
	c := p.addObjectKeywordsClause(s, "WITH RULES", start)
	for {
		rule, err := p.parseEventRule()
		if err != nil {
			return err
		}
		c.Options = append(c.Options, rule)
		if !p.startsEventRuleName() || !p.nextIs('=') {
			break
		}
	}
	c.Loc.End = p.prev.Loc.End
	return nil
}

// parseAlterEventRoutingTable reads what follows the name of
//
//	ALTER EVENT ROUTING TABLE <table_name> [ FORCE ] SET RULES
//	  <rule_name> = ( REGION_GROUP = <region_group>, REGIONS = ( '<region>' [ , ... ] ),
//	                  DESTINATION_ACCOUNT = <organization>.<account_name> ) [ , ... ]
//	ALTER EVENT ROUTING TABLE <table_name> [ FORCE ] SET RULE
//	  <rule_name> REGION_GROUP = <region_group> REGIONS = ( '<region>' [ , ... ] )
//	  DESTINATION_ACCOUNT = <organization>.<account_name>
//	ALTER EVENT ROUTING TABLE <table_name> [ FORCE ] UNSET RULE <rule_name>
//	ALTER EVENT ROUTING TABLE <table_name> MODIFY RULE <rule_name> RENAME TO <new_rule_name>
//	ALTER EVENT ROUTING TABLE <table_name> RENAME TO <new_table_name>
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-event-routing-table
func parseAlterEventRoutingTable(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	force := p.acceptWords("FORCE")
	if force {
		p.addObjectKeywordsClause(s, "FORCE", start)
		start = p.cur.Loc.Start
	}
	switch {
	case p.curIs(kwSET) && p.nextWord("RULES"):
		p.advance() // consume SET
		p.advance() // consume RULES
		c := p.addObjectKeywordsClause(s, "SET RULES", start)
		for {
			rule, err := p.parseEventRule()
			if err != nil {
				return err
			}
			c.Options = append(c.Options, rule)
			if !p.curIs(',') {
				break
			}
			p.advance() // consume ','
		}
		c.Loc.End = p.prev.Loc.End
		return nil
	case p.curIs(kwSET) && p.nextWord("RULE"):
		p.advance() // consume SET
		p.advance() // consume RULE
		name, err := p.parseEventRuleName()
		if err != nil {
			return err
		}
		var settings []*ast.CopyOption
		for range 3 {
			if !p.startsObjectProperty() {
				return p.syntaxErrorAtCur()
			}
			o, err := p.parseObjectProperty()
			if err != nil {
				return err
			}
			settings = append(settings, o)
		}
		if err := checkEventRuleSettings(settings, ast.Loc{Start: start, End: p.prev.Loc.End}); err != nil {
			return err
		}
		c := p.addObjectKeywordsClause(s, "SET RULE", start)
		c.Names = []*ast.ObjectName{name}
		c.Options = settings
		return nil
	case p.curIs(kwUNSET) && p.nextWord("RULE"):
		p.advance() // consume UNSET
		p.advance() // consume RULE
		name, err := p.parseEventRuleName()
		if err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "UNSET RULE", start).Names = []*ast.ObjectName{name}
		return nil
	case force:
		// FORCE goes only before SET RULES, SET RULE and UNSET RULE.
		return p.syntaxErrorAtCur()
	case p.acceptWords("MODIFY", "RULE"):
		name, err := p.parseEventRuleName()
		if err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "MODIFY RULE", start).Names = []*ast.ObjectName{name}
		renameStart := p.cur.Loc.Start
		if !p.curIs(kwRENAME) {
			return p.syntaxErrorAtCur()
		}
		p.advance() // consume RENAME
		if _, err := p.expect(kwTO); err != nil {
			return err
		}
		newName, err := p.parseEventRuleName()
		if err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "RENAME TO", renameStart).Names = []*ast.ObjectName{newName}
		return nil
	}
	return p.parseObjectAlterCommentOrRename(s, actRename)
}

// startsEventRuleName reports whether cur can be a rule's name: an
// identifier, a quoted identifier, or a keyword (the CREATE page's example
// names a rule default). It consumes nothing.
func (p *Parser) startsEventRuleName() bool {
	return p.curIs(tokIdent) || p.curIs(tokQuotedIdent) || p.cur.Type >= 700
}

// parseEventRuleName reads a rule's one-part name (startsEventRuleName).
func (p *Parser) parseEventRuleName() (*ast.ObjectName, error) {
	if !p.startsEventRuleName() {
		return nil, p.syntaxErrorAtCur()
	}
	tok := p.advance()
	id := ast.Ident{Name: tok.Str, Quoted: tok.Type == tokQuotedIdent, Loc: tok.Loc}
	return &ast.ObjectName{Name: id, Loc: tok.Loc}, nil
}

// parseEventRule reads <rule_name> = ( <settings> ) into an option named by
// the rule (uppercased; a quoted name keeps its quotes) whose Group holds the
// settings.
func (p *Parser) parseEventRule() (*ast.CopyOption, error) {
	name, err := p.parseEventRuleName()
	if err != nil {
		return nil, err
	}
	rule := &ast.CopyOption{Name: name.Name.Normalize(), Loc: ast.Loc{Start: name.Loc.Start}}
	if name.Name.Quoted {
		rule.Name = name.Name.String()
	}
	if _, err := p.expect('='); err != nil {
		return nil, err
	}
	if !p.curIs('(') {
		return nil, p.syntaxErrorAtCur()
	}
	if err := p.parseObjectValueParen(rule); err != nil {
		return nil, err
	}
	rule.Loc.End = p.prev.Loc.End
	if err := checkEventRuleSettings(rule.Group, rule.Loc); err != nil {
		return nil, err
	}
	return rule, nil
}

// checkEventRuleSettings holds a rule's settings to the pages' three, in
// their order: REGION_GROUP = <region_group> (a name; the CREATE page's
// example quotes it), REGIONS = ( '<region>' [ , ... ] ) and
// DESTINATION_ACCOUNT = <organization>.<account_name>.
func checkEventRuleSettings(settings []*ast.CopyOption, loc ast.Loc) error {
	want := []string{"REGION_GROUP", "REGIONS", "DESTINATION_ACCOUNT"}
	for i, o := range settings {
		ok := i < len(want) && o.Name == want[i]
		if ok {
			switch o.Name {
			case "REGION_GROUP":
				ok = propStringOr().matches(o) || propName.matches(o)
			case "REGIONS":
				ok = propStrings.matches(o)
			case "DESTINATION_ACCOUNT":
				ok = o.Words != "" && nameParts(o.Words) == 2
			}
		}
		if !ok {
			return &ParseError{Position: o.Loc.Start, End: o.Loc.End, Message: "syntax error at or near " + o.Name}
		}
	}
	if len(settings) != len(want) {
		return &ParseError{Position: loc.Start, End: loc.End, Message: "syntax error: a rule sets REGION_GROUP, REGIONS and DESTINATION_ACCOUNT"}
	}
	return nil
}
