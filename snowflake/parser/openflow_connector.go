package parser

import "github.com/bytebase/omni/snowflake/ast"

// ---------------------------------------------------------------------------
// Openflow CONNECTOR, read as an ast.ObjectStmt
// ---------------------------------------------------------------------------

func init() {
	registerObject(&objectForm{
		words:  []string{"OPENFLOW", "CONNECTOR"},
		create: parseCreateOpenflowConnector,
		alter:  applicationTagless(parseAlterOpenflowConnector, "SET"),
		// https://docs.snowflake.com/en/sql-reference/sql/drop-openflow-connector
		drop: true,
		// ALTER ... SET documents DISPLAY_NAME and COMMENT together, or
		// DEFAULT_VERSION alone (parseAlterOpenflowConnector), and no UNSET.
		props: objectProps{
			"DISPLAY_NAME":    propString,
			"COMMENT":         propString,
			"DEFAULT_VERSION": propStringOr("FIRST", "LAST"),
		},
		unset: []string{},
	})
}

// parseCreateOpenflowConnector reads what follows the name of
//
//	CREATE OPENFLOW CONNECTOR [ IF NOT EXISTS ] <name>
//	  IN RUNTIME <runtime_name>
//	  FROM DEFINITION <connector_definition_id>
//	  [ DISPLAY_NAME = <string> ]
//	  [ COMMENT = <string> ]
//
//	CREATE OPENFLOW CONNECTOR [ IF NOT EXISTS ] <name>
//	  IN RUNTIME <runtime_name>
//	  FROM '<stage_reference>'
//	  [ COMMENT = <string> ]
//
// https://docs.snowflake.com/en/sql-reference/sql/create-openflow-connector
//
// The page's examples qualify the runtime (my_db.my_schema.my_runtime), so
// IN RUNTIME takes a qualified name. FROM DEFINITION holds the definition's
// one-part name in Names; FROM holds the stage reference as its one nameless
// option.
func parseCreateOpenflowConnector(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	if err := p.expectWords("IN", "RUNTIME"); err != nil {
		return err
	}
	runtime, err := p.parseObjectName()
	if err != nil {
		return err
	}
	p.addObjectKeywordsClause(s, "IN RUNTIME", start).Names = []*ast.ObjectName{runtime}

	start = p.cur.Loc.Start
	if _, err := p.expect(kwFROM); err != nil {
		return err
	}
	allowed := objectProps{"COMMENT": propString}
	if p.acceptWords("DEFINITION") {
		def, err := p.parseObjectSimpleName()
		if err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "FROM DEFINITION", start).Names = []*ast.ObjectName{def}
		allowed["DISPLAY_NAME"] = propString
	} else {
		ref, err := p.expect(tokString)
		if err != nil {
			return err
		}
		lit := &ast.CopyOption{Lit: &ast.Literal{Kind: ast.LitString, Value: ref.Str, Loc: ref.Loc}, Loc: ref.Loc}
		p.addObjectKeywordsClause(s, "FROM", start).Options = []*ast.CopyOption{lit}
	}
	if !p.startsObjectProperty() {
		return nil
	}
	props, err := p.parseObjectProperties(s, "", false)
	if err != nil {
		return err
	}
	return notebookCheckOptions(props, allowed)
}

// parseAlterOpenflowConnector reads what follows the name of
//
//	ALTER OPENFLOW CONNECTOR <name> ADD VERSION FROM '<stage_reference>'
//	ALTER OPENFLOW CONNECTOR <name> ADD LIVE VERSION FROM { LAST | <alias> } [ COMMENT = <string> ]
//	ALTER OPENFLOW CONNECTOR <name> COMMIT [ COMMENT = <string> ]
//	ALTER OPENFLOW CONNECTOR <name> ABORT
//	ALTER OPENFLOW CONNECTOR [ IF EXISTS ] <name> SET DEFAULT_VERSION = { FIRST | LAST | '<version>' }
//	ALTER OPENFLOW CONNECTOR [ IF EXISTS ] <name> { START | STOP }
//	ALTER OPENFLOW CONNECTOR [ IF EXISTS ] <name> TERMINATE [ FORCE ]
//	ALTER OPENFLOW CONNECTOR [ IF EXISTS ] <name> SET [ DISPLAY_NAME = <string> ] [ COMMENT = <string> ]
//	ALTER OPENFLOW CONNECTOR [ IF EXISTS ] <name> RENAME TO <new_name>
//	ALTER OPENFLOW CONNECTOR <name> PUSH [ TO '<git_branch_uri>' ]
//	  USERNAME = <string> PASSWORD = <string> NAME = <string> EMAIL = <string> [ COMMENT = <string> ]
//	ALTER OPENFLOW CONNECTOR <name> PULL
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-openflow-connector
//
// The forms the page writes without IF EXISTS reject it.
func parseAlterOpenflowConnector(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	noIfExists := func() error {
		if s.IfExists {
			return &ParseError{Position: start, End: p.prev.Loc.End, Message: "syntax error: IF EXISTS is not documented here"}
		}
		return nil
	}
	switch {
	case p.acceptWords("ADD", "VERSION", "FROM"):
		ref, err := p.expect(tokString)
		if err != nil {
			return err
		}
		lit := &ast.CopyOption{Lit: &ast.Literal{Kind: ast.LitString, Value: ref.Str, Loc: ref.Loc}, Loc: ref.Loc}
		p.addObjectKeywordsClause(s, "ADD VERSION FROM", start).Options = []*ast.CopyOption{lit}
		return noIfExists()
	case p.acceptWords("ADD", "LIVE", "VERSION", "FROM"):
		alias, err := p.parseObjectSimpleName()
		if err != nil {
			return err
		}
		c := p.addObjectKeywordsClause(s, "ADD LIVE VERSION FROM", start)
		c.Names = []*ast.ObjectName{alias}
		if err := p.parseOpenflowConnectorComment(c); err != nil {
			return err
		}
		return noIfExists()
	case p.acceptWords("COMMIT"):
		c := p.addObjectKeywordsClause(s, "COMMIT", start)
		if err := p.parseOpenflowConnectorComment(c); err != nil {
			return err
		}
		return noIfExists()
	case p.acceptWords("ABORT"):
		p.addObjectKeywordsClause(s, "ABORT", start)
		return noIfExists()
	case p.acceptWords("PULL"):
		p.addObjectKeywordsClause(s, "PULL", start)
		return noIfExists()
	case p.acceptWords("PUSH"):
		c := p.addObjectKeywordsClause(s, "PUSH", start)
		if p.cur.Type == kwTO {
			p.advance() // consume TO
			uri, err := p.expect(tokString)
			if err != nil {
				return err
			}
			c.Keywords = "PUSH TO"
			c.Options = []*ast.CopyOption{{Lit: &ast.Literal{Kind: ast.LitString, Value: uri.Str, Loc: uri.Loc}, Loc: uri.Loc}}
		}
		props, err := p.parseObjectProperties(s, "", true)
		if err != nil {
			return err
		}
		pushProps := objectProps{"USERNAME": propString, "PASSWORD": propString, "NAME": propString, "EMAIL": propString, "COMMENT": propString}
		if err := notebookCheckOptions(props, pushProps); err != nil {
			return err
		}
		if err := p.requireObjectProperties(props, "USERNAME", "PASSWORD", "NAME", "EMAIL"); err != nil {
			return err
		}
		// The run's options belong to PUSH: fold them into its clause.
		c.Options = append(c.Options, props.Options...)
		c.Loc.End = props.Loc.End
		s.Clauses = s.Clauses[:len(s.Clauses)-1]
		return noIfExists()
	case p.acceptWords("START"):
		p.addObjectKeywordsClause(s, "START", start)
		return nil
	case p.acceptWords("STOP"):
		p.addObjectKeywordsClause(s, "STOP", start)
		return nil
	case p.acceptWords("TERMINATE"):
		c := p.addObjectKeywordsClause(s, "TERMINATE", start)
		if p.acceptWords("FORCE") {
			c.Keywords = "TERMINATE FORCE"
			c.Loc.End = p.prev.Loc.End
		}
		return nil
	}
	if err := p.parseObjectAlterCommentOrRename(s, actSet|actRename); err != nil {
		return err
	}
	// SET DEFAULT_VERSION goes alone; DISPLAY_NAME and COMMENT go together.
	if c := s.Clauses[len(s.Clauses)-1]; c.Keywords == "SET" && len(c.Options) > 1 {
		for _, o := range c.Options {
			if o.Name == "DEFAULT_VERSION" {
				return &ParseError{Position: o.Loc.Start, End: o.Loc.End, Message: "syntax error: DEFAULT_VERSION is set alone"}
			}
		}
	}
	return nil
}

// parseOpenflowConnectorComment reads an optional COMMENT = <string> into
// clause c's options.
func (p *Parser) parseOpenflowConnectorComment(c *ast.ObjectClause) error {
	if !isWordToken(p.cur, "COMMENT") || p.peekNext().Type != '=' {
		return nil
	}
	o, err := p.parseObjectProperty()
	if err != nil {
		return err
	}
	if !propString.matches(o) {
		return &ParseError{Position: o.Loc.Start, End: o.Loc.End, Message: "syntax error: a comment is a string"}
	}
	c.Options = append(c.Options, o)
	c.Loc.End = p.prev.Loc.End
	return nil
}

// parseShowConnectorOrRoutingScope reads the scope after IN of
//
//	SHOW VERSIONS IN OPENFLOW CONNECTOR <name>
//	SHOW RULES IN EVENT ROUTING TABLE ( <event_routing_table_name> )
//
// (https://docs.snowflake.com/en/sql-reference/sql/show-versions-in-openflow-connector,
// https://docs.snowflake.com/en/sql-reference/sql/show-rules-in-event-routing-table)
// as a ShowScopeOther, and reports whether one started at cur.
func (p *Parser) parseShowConnectorOrRoutingScope(stmt *ast.ShowStmt) (bool, error) {
	switch {
	case stmt.ObjectClass == "VERSIONS" && p.acceptWords("OPENFLOW", "CONNECTOR"):
		stmt.Scope = ast.ShowScopeOther
		stmt.ScopeText = "OPENFLOW CONNECTOR"
		name, err := p.parseObjectName()
		stmt.ScopeName = name
		return true, err
	case stmt.ObjectClass == "RULES" && p.acceptWords("EVENT", "ROUTING", "TABLE"):
		stmt.Scope = ast.ShowScopeOther
		stmt.ScopeText = "EVENT ROUTING TABLE"
		if _, err := p.expect('('); err != nil {
			return true, err
		}
		name, err := p.parseObjectName()
		if err != nil {
			return true, err
		}
		stmt.ScopeName = name
		_, err = p.expect(')')
		return true, err
	}
	return false, nil
}

// checkShowConnectorOrRoutingScope fails a SHOW in an Openflow connector's or
// an event routing table's scope that carries any filter: neither page writes
// one.
func checkShowConnectorOrRoutingScope(stmt *ast.ShowStmt) error {
	if stmt.Scope != ast.ShowScopeOther || (stmt.ScopeText != "OPENFLOW CONNECTOR" && stmt.ScopeText != "EVENT ROUTING TABLE") {
		return nil
	}
	if stmt.Terse || stmt.History || stmt.HasLike || stmt.HasStarts || stmt.HasLimit || stmt.Privileges != nil {
		return &ParseError{
			Position: stmt.Loc.Start, End: stmt.Loc.End,
			Message: "syntax error: SHOW " + stmt.ObjectClass + " IN " + stmt.ScopeText + " takes no filter",
		}
	}
	return nil
}
