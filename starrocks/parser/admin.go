package parser

import (
	"strings"

	"github.com/bytebase/omni/starrocks/ast"
)

// parseFrontendConfigList parses ("key"="value") without consuming a leading
// PROPERTIES keyword. cur is '(' on entry.
func (p *Parser) parseFrontendConfigList() ([]*ast.Property, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}

	return p.parsePropertyItems(p.parseStringProperty)
}

// ---------------------------------------------------------------------------
// ALTER SYSTEM statements
// ---------------------------------------------------------------------------

// parseSystemAlter parses ALTER SYSTEM ... statements.
//
// ALTER has already been consumed; cur is SYSTEM on entry. SYSTEM is consumed
// here.
//
// Supported forms:
//
//	ALTER SYSTEM ADD BACKEND 'host:port' [, ...] [PROPERTIES(...)]
//	ALTER SYSTEM DROP BACKEND 'host:port' [, ...]
//	ALTER SYSTEM DECOMMISSION BACKEND 'host:port' [, ...]
//	ALTER SYSTEM MODIFY BACKEND 'host:port' SET ("key"="value")
//	ALTER SYSTEM ADD FOLLOWER 'host:port'
//	ALTER SYSTEM DROP FOLLOWER 'host:port'
//	ALTER SYSTEM ADD OBSERVER 'host:port'
//	ALTER SYSTEM DROP OBSERVER 'host:port'
//	ALTER SYSTEM ADD BROKER broker_name 'host:port' [, ...]
//	ALTER SYSTEM DROP BROKER broker_name 'host:port' [, ...]
//	ALTER SYSTEM DROP ALL BROKER broker_name
//	ALTER SYSTEM SET LOAD ERROR HUB PROPERTIES(...)
func (p *Parser) parseSystemAlter(alterLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume SYSTEM

	stmt := &ast.SystemAlterStmt{}
	endLoc := p.prev.Loc

	// Verb: ADD, DROP, DECOMMISSION, MODIFY, SET
	verbTok := p.cur
	if verbTok.Kind == tokEOF {
		return nil, p.syntaxErrorAtCur()
	}
	verb := strings.ToUpper(verbTok.Str)
	if verb == "" {
		verb = strings.ToUpper(TokenName(verbTok.Kind))
	}
	stmt.Verb = verb
	endLoc = verbTok.Loc
	p.advance()

	var err error
	switch strings.ToUpper(verb) {
	case "ADD":
		if endLoc, err = p.parseSystemAlterAddBody(stmt, endLoc); err != nil {
			return nil, err
		}
	case "DROP":
		if endLoc, err = p.parseSystemAlterDropBody(stmt, endLoc); err != nil {
			return nil, err
		}
	case "DECOMMISSION":
		// DECOMMISSION BACKEND 'host:port' [, ...]
		stmt.Object = "BACKEND"
		if p.cur.Kind == kwBACKEND {
			endLoc = p.cur.Loc
			p.advance()
		}
		hosts, hostsEnd := p.parseHostList()
		stmt.Hosts = hosts
		if hostsEnd.IsValid() {
			endLoc = hostsEnd
		}
	case "MODIFY":
		// MODIFY BACKEND 'host:port' SET ("key"="value")
		stmt.Object = "BACKEND"
		if p.cur.Kind == kwBACKEND {
			endLoc = p.cur.Loc
			p.advance()
		}
		hosts, hostsEnd := p.parseHostList()
		stmt.Hosts = hosts
		if hostsEnd.IsValid() {
			endLoc = hostsEnd
		}
		// SET ("key"="value")
		if p.cur.Kind == kwSET {
			endLoc = p.cur.Loc
			p.advance()
			if p.cur.Kind == int('(') {
				props, err := p.parseFrontendConfigList()
				if err != nil {
					return nil, err
				}
				stmt.SetClause = props
				endLoc = p.prev.Loc
			}
		}
	case "SET":
		// SET LOAD ERROR HUB PROPERTIES(...)
		stmt.Object = "LOAD ERROR HUB"
		// consume LOAD ERROR HUB (best-effort)
		if p.cur.Kind == kwLOAD {
			endLoc = p.cur.Loc
			p.advance()
		}
		// ERROR and HUB are likely plain idents in the lexer
		for p.cur.Kind != kwPROPERTIES && p.cur.Kind != tokEOF {
			endLoc = p.cur.Loc
			p.advance()
		}
		if endLoc, err = p.parseSystemAlterProperties(stmt, endLoc); err != nil {
			return nil, err
		}
	default:
		// Unknown verb — collect as args
		for p.cur.Kind != tokEOF {
			endLoc = p.cur.Loc
			p.advance()
		}
	}

	stmt.Loc = alterLoc.Merge(endLoc)
	return stmt, nil
}

// parseSystemAlterProperties parses an optional PROPERTIES(...) list at cur
// into stmt.Properties and returns the updated endLoc, which ends at the
// list's closing ')'. An error in the list is returned, not dropped.
func (p *Parser) parseSystemAlterProperties(stmt *ast.SystemAlterStmt, endLoc ast.Loc) (ast.Loc, error) {
	if p.cur.Kind != kwPROPERTIES {
		return endLoc, nil
	}
	props, err := p.parseProperties()
	if err != nil {
		return endLoc, err
	}
	stmt.Properties = props
	return p.prev.Loc, nil
}

// parseSystemAlterBroker parses the broker name of ALTER SYSTEM ADD / DROP
// [ALL] BROKER into stmt.BrokerName, returning the updated endLoc. The name
// is required: a missing one used to be dropped silently, so ALTER SYSTEM ADD
// BROKER parsed (container-verified reject).
func (p *Parser) parseSystemAlterBroker(stmt *ast.SystemAlterStmt) (ast.Loc, error) {
	brokerName, brokerLoc, err := p.parseIdentifierOrString()
	if err != nil {
		return ast.Loc{}, err
	}
	stmt.BrokerName = brokerName
	return brokerLoc, nil
}

// parseSystemAlterAddBody parses the body of ALTER SYSTEM ADD ...
func (p *Parser) parseSystemAlterAddBody(stmt *ast.SystemAlterStmt, endLoc ast.Loc) (ast.Loc, error) {
	var err error
	obj := strings.ToUpper(p.cur.Str)
	if obj == "" {
		obj = strings.ToUpper(TokenName(p.cur.Kind))
	}
	switch obj {
	case "BACKEND":
		stmt.Object = "BACKEND"
		endLoc = p.cur.Loc
		p.advance()
		// 'host:port,host:port' or 'host:port', 'host:port'
		hosts, hostsEnd := p.parseHostList()
		stmt.Hosts = hosts
		if hostsEnd.IsValid() {
			endLoc = hostsEnd
		}
		// optional PROPERTIES
		if endLoc, err = p.parseSystemAlterProperties(stmt, endLoc); err != nil {
			return endLoc, err
		}
	case "FOLLOWER":
		stmt.Object = "FOLLOWER"
		endLoc = p.cur.Loc
		p.advance()
		hosts, hostsEnd := p.parseHostList()
		stmt.Hosts = hosts
		if hostsEnd.IsValid() {
			endLoc = hostsEnd
		}
	case "OBSERVER":
		stmt.Object = "OBSERVER"
		endLoc = p.cur.Loc
		p.advance()
		hosts, hostsEnd := p.parseHostList()
		stmt.Hosts = hosts
		if hostsEnd.IsValid() {
			endLoc = hostsEnd
		}
	case "BROKER":
		stmt.Object = "BROKER"
		endLoc = p.cur.Loc
		p.advance()
		// broker_name 'host:port' [, ...]
		if endLoc, err = p.parseSystemAlterBroker(stmt); err != nil {
			return endLoc, err
		}
		hosts, hostsEnd := p.parseHostList()
		stmt.Hosts = hosts
		if hostsEnd.IsValid() {
			endLoc = hostsEnd
		}
	default:
		// Unknown object — collect remaining as args
		for p.cur.Kind != tokEOF {
			endLoc = p.cur.Loc
			p.advance()
		}
	}
	return endLoc, nil
}

// parseSystemAlterDropBody parses the body of ALTER SYSTEM DROP ...
func (p *Parser) parseSystemAlterDropBody(stmt *ast.SystemAlterStmt, endLoc ast.Loc) (ast.Loc, error) {
	var err error
	// Check for DROP ALL BROKER
	if p.cur.Kind == kwALL {
		stmt.DropAll = true
		endLoc = p.cur.Loc
		p.advance()
		// expect BROKER
		if p.cur.Kind == kwBROKER {
			stmt.Object = "BROKER"
			endLoc = p.cur.Loc
			p.advance()
		}
		return p.parseSystemAlterBroker(stmt)
	}

	obj := strings.ToUpper(p.cur.Str)
	if obj == "" {
		obj = strings.ToUpper(TokenName(p.cur.Kind))
	}
	switch obj {
	case "BACKEND":
		stmt.Object = "BACKEND"
		endLoc = p.cur.Loc
		p.advance()
		hosts, hostsEnd := p.parseHostList()
		stmt.Hosts = hosts
		if hostsEnd.IsValid() {
			endLoc = hostsEnd
		}
	case "FOLLOWER":
		stmt.Object = "FOLLOWER"
		endLoc = p.cur.Loc
		p.advance()
		hosts, hostsEnd := p.parseHostList()
		stmt.Hosts = hosts
		if hostsEnd.IsValid() {
			endLoc = hostsEnd
		}
	case "OBSERVER":
		stmt.Object = "OBSERVER"
		endLoc = p.cur.Loc
		p.advance()
		hosts, hostsEnd := p.parseHostList()
		stmt.Hosts = hosts
		if hostsEnd.IsValid() {
			endLoc = hostsEnd
		}
	case "BROKER":
		stmt.Object = "BROKER"
		endLoc = p.cur.Loc
		p.advance()
		if endLoc, err = p.parseSystemAlterBroker(stmt); err != nil {
			return endLoc, err
		}
		hosts, hostsEnd := p.parseHostList()
		stmt.Hosts = hosts
		if hostsEnd.IsValid() {
			endLoc = hostsEnd
		}
	default:
		for p.cur.Kind != tokEOF {
			endLoc = p.cur.Loc
			p.advance()
		}
	}
	return endLoc, nil
}

// parseHostList parses a comma-separated list of string literal host:port
// specs. Handles both single "host:port" and multiple "host:port", "host:port"
// forms, as well as the Doris-specific "host1:port,host2:port" single-string
// form (split on comma internally).
// Returns the parsed hosts and the last Loc seen.
func (p *Parser) parseHostList() ([]string, ast.Loc) {
	var hosts []string
	var endLoc ast.Loc

	for p.cur.Kind == tokString {
		s := p.cur.Str
		endLoc = p.cur.Loc
		p.advance()
		// A single string may encode multiple hosts: "h1:p,h2:p"
		for _, h := range strings.Split(s, ",") {
			h = strings.TrimSpace(h)
			if h != "" {
				hosts = append(hosts, h)
			}
		}
		// Optional comma before next string
		if p.cur.Kind == int(',') {
			p.advance()
		} else {
			break
		}
	}
	return hosts, endLoc
}

// ---------------------------------------------------------------------------
// CANCEL DECOMMISSION BACKEND
// ---------------------------------------------------------------------------

// parseCancelDecommission parses CANCEL DECOMMISSION BACKEND 'host:port' [, ...]
//
// CANCEL has already been consumed; cur is DECOMMISSION on entry.
func (p *Parser) parseCancelDecommission(cancelLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume DECOMMISSION

	stmt := &ast.CancelDecommissionStmt{Object: "BACKEND"}
	endLoc := p.prev.Loc

	if p.cur.Kind == kwBACKEND {
		endLoc = p.cur.Loc
		p.advance()
	}

	hosts, hostsEnd := p.parseHostList()
	stmt.Hosts = hosts
	if hostsEnd.IsValid() {
		endLoc = hostsEnd
	}

	stmt.Loc = cancelLoc.Merge(endLoc)
	return stmt, nil
}
