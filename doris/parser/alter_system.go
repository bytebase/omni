package parser

import (
	"strings"

	"github.com/bytebase/omni/doris/ast"
)

// parseAlterSystemStmt parses ALTER SYSTEM. ALTER has been consumed; cur is
// SYSTEM. Doris reads ALTER SYSTEM RENAME COMPUTE GROUP itself
// (DorisParser.g4:207 alterSystemRenameComputeGroup) and hands every other
// form to its legacy CUP parser (DorisParser.g4:563 alterSystem), so those
// follow both alterSystemClause (DorisParser.g4:593) and alter_system_clause
// (DorisLegacyParser.g4:385). A best-effort parse the grammar refuses falls
// back to the loose reader.
func (p *Parser) parseAlterSystemStmt(alterLoc ast.Loc) (ast.Node, error) {
	cp := p.save()
	p.advance() // consume SYSTEM
	stmt := &ast.SystemAlterStmt{}
	if err := p.parseAlterSystemClause(stmt); err != nil {
		if p.strictTrailing {
			return nil, err
		}
		p.restore(cp)
		return p.parseSystemAlter(alterLoc)
	}
	stmt.Loc = ast.Loc{Start: alterLoc.Start, End: p.prev.Loc.End}
	return stmt, nil
}

// parseAlterSystemClause reads one ALTER SYSTEM clause into stmt:
//
//	ADD BACKEND string_list opt_properties
//	(DROP | DROPP) BACKEND string_list
//	DECOMMISSION BACKEND string_list
//	(ADD | DROP) (OBSERVER | FOLLOWER) STRING_LITERAL
//	(ADD | DROP) BROKER ident_or_text string_list
//	DROP ALL BROKER ident_or_text
//	SET LOAD ERRORS HUB opt_properties
//	MODIFY BACKEND string_list SET '(' key_value_map ')'
//	MODIFY (FRONTEND | BACKEND) STRING_LITERAL HOSTNAME STRING_LITERAL
//	RENAME COMPUTE GROUP identifier identifier
func (p *Parser) parseAlterSystemClause(stmt *ast.SystemAlterStmt) error {
	verb := p.cur
	stmt.Verb = strings.ToUpper(verb.Str)
	var err error
	switch verb.Kind {
	case kwADD, kwDROP, kwDROPP:
		p.advance()
		switch {
		case p.cur.Kind == kwBACKEND:
			// DorisLegacyParser.g4:1906-1908; DorisParser.g4:594-597
			p.advance()
			stmt.Object = "BACKEND"
			if stmt.Hosts, err = p.parseAlterStringList(); err != nil {
				return err
			}
			if verb.Kind == kwADD {
				stmt.Properties, err = p.alterOptProperties()
			}
			return err
		case verb.Kind == kwDROPP:
			return p.syntaxErrorAtCur()
		case p.cur.Kind == kwOBSERVER || p.cur.Kind == kwFOLLOWER:
			// DorisLegacyParser.g4:1910-1913; DorisParser.g4:600-603
			stmt.Object = strings.ToUpper(p.advance().Str)
			tok, err := p.alterExpectString()
			if err != nil {
				return err
			}
			stmt.Hosts = []string{tok.Str}
			return nil
		case p.cur.Kind == kwBROKER:
			// DorisLegacyParser.g4:1914-1915; DorisParser.g4:604-607
			p.advance()
			stmt.Object = "BROKER"
			if stmt.BrokerName, _, err = p.parseIdentifierOrString(); err != nil {
				return err
			}
			stmt.Hosts, err = p.parseAlterStringList()
			return err
		case verb.Kind == kwDROP && p.cur.Kind == kwALL:
			// DROP ALL BROKER ident_or_text (DorisLegacyParser.g4:1916;
			// DorisParser.g4:608)
			p.advance()
			if _, err = p.expect(kwBROKER); err != nil {
				return err
			}
			stmt.Object = "BROKER"
			stmt.DropAll = true
			stmt.BrokerName, _, err = p.parseIdentifierOrString()
			return err
		}
	case kwDECOMMISSION:
		// DorisLegacyParser.g4:1909; DorisParser.g4:598
		p.advance()
		if _, err = p.expect(kwBACKEND); err != nil {
			return err
		}
		stmt.Object = "BACKEND"
		stmt.Hosts, err = p.parseAlterStringList()
		return err
	case kwSET:
		// SET LOAD ERRORS HUB opt_properties (DorisLegacyParser.g4:1917;
		// DorisParser.g4:609)
		p.advance()
		for _, kw := range []int{kwLOAD, kwERRORS, kwHUB} {
			if _, err = p.expect(kw); err != nil {
				return err
			}
		}
		stmt.Object = "LOAD ERROR HUB"
		stmt.Properties, err = p.alterOptProperties()
		return err
	case kwMODIFY:
		p.advance()
		if p.cur.Kind != kwBACKEND && p.cur.Kind != kwFRONTEND {
			return p.syntaxErrorAtCur()
		}
		stmt.Object = strings.ToUpper(p.advance().Str)
		if stmt.Hosts, err = p.parseAlterStringList(); err != nil {
			return err
		}
		if p.cur.Kind == kwHOSTNAME && len(stmt.Hosts) == 1 {
			// MODIFY (FRONTEND | BACKEND) STRING_LITERAL HOSTNAME
			// STRING_LITERAL (DorisLegacyParser.g4:1919-1920;
			// DorisParser.g4:613)
			p.advance()
			tok, err := p.alterExpectString()
			if err != nil {
				return err
			}
			stmt.NewHost = tok.Str
			return nil
		}
		if stmt.Object != "BACKEND" {
			return p.syntaxErrorAtCur()
		}
		// MODIFY BACKEND string_list SET '(' key_value_map ')'
		// (DorisLegacyParser.g4:1918; DorisParser.g4:610)
		if _, err = p.expect(kwSET); err != nil {
			return err
		}
		stmt.SetClause, err = p.parseParenProperties()
		return err
	case kwRENAME:
		// RENAME COMPUTE GROUP identifier identifier (DorisParser.g4:207)
		p.advance()
		for _, kw := range []int{kwCOMPUTE, kwGROUP} {
			if _, err = p.expect(kw); err != nil {
				return err
			}
		}
		stmt.Object = "COMPUTE GROUP"
		if stmt.Name, _, err = p.parseIdentifier(); err != nil {
			return err
		}
		stmt.NewName, _, err = p.parseIdentifier()
		return err
	}
	return p.syntaxErrorAtCur()
}

// parseAlterStringList reads string_list: STRING_LITERAL (','
// STRING_LITERAL)* (DorisLegacyParser.g4 string_list). A string holding
// several comma-separated hosts is split, as Hosts always recorded them.
func (p *Parser) parseAlterStringList() ([]string, error) {
	var hosts []string
	for {
		tok, err := p.alterExpectString()
		if err != nil {
			return nil, err
		}
		for _, h := range strings.Split(tok.Str, ",") {
			if h = strings.TrimSpace(h); h != "" {
				hosts = append(hosts, h)
			}
		}
		if p.cur.Kind != int(',') {
			return hosts, nil
		}
		p.advance()
	}
}
