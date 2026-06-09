package parser

import (
	"strings"

	"github.com/bytebase/omni/starrocks/ast"
)

// ---------------------------------------------------------------------------
// ADMIN statements
// ---------------------------------------------------------------------------

// parseAdminStmt parses ADMIN ... statements.
//
// ADMIN has already been consumed by the caller; cur is the verb keyword.
//
// Supported forms (best-effort — remaining tokens after verb+object are
// captured as raw text in Args or via specific structured fields):
//
//	ADMIN SHOW REPLICA DISTRIBUTION FROM table [PARTITION p]
//	ADMIN SHOW REPLICA STATUS FROM table [WHERE ...]
//	ADMIN SHOW CONFIG [LIKE 'pat']
//	ADMIN SHOW FRONTEND CONFIG [LIKE 'pat']
//	ADMIN SHOW TABLET STORAGE FORMAT [VERBOSE]
//	ADMIN SHOW DATA SKEW FROM table
//	ADMIN REBALANCE DISK [ON ('be1', ...)]
//	ADMIN CANCEL REBALANCE DISK [ON ('be1', ...)]
//	ADMIN DIAGNOSE TABLET tablet_id
//	ADMIN COMPACT TABLE name [WHERE ...]
//	ADMIN CHECK TABLET (id, ...) PROPERTIES(...)
//	ADMIN REPAIR TABLE name [PARTITION(p)]
//	ADMIN CANCEL REPAIR TABLE name [PARTITION(p)]
//	ADMIN SET REPLICA STATUS PROPERTIES(...)
//	ADMIN SET TABLE name STATUS PROPERTIES(...)
//	ADMIN SET PARTITION VERSION PROPERTIES(...)
//	ADMIN SET FRONTEND CONFIG ("key"="value")
//	ADMIN CLEAN TRASH [ON ('be1', ...)]
//	ADMIN COPY TABLET tablet_id PROPERTIES(...)
//	ADMIN DECOMMISSION BACKEND BY HOSTNAME 'host:port'
func (p *Parser) parseAdminStmt(adminLoc ast.Loc) (ast.Node, error) {
	stmt := &ast.AdminStmt{}
	startLoc := adminLoc

	// Capture the primary verb.
	verb, verbLoc, err := p.parseAdminVerb()
	if err != nil {
		return nil, err
	}
	stmt.Verb = verb
	endLoc := verbLoc

	// Capture the object keyword(s) and optional structured fields.
	switch strings.ToUpper(verb) {
	case "SHOW":
		endLoc = p.parseAdminShowBody(stmt, endLoc)
	case "REBALANCE":
		// REBALANCE DISK [ON (...)]
		stmt.Object = "DISK"
		if p.cur.Kind == kwDISK {
			endLoc = p.cur.Loc
			p.advance()
		}
		endLoc = p.skipAdminOnClause(endLoc)
	case "DIAGNOSE":
		// DIAGNOSE TABLET tablet_id
		stmt.Object = "TABLET"
		if p.cur.Kind == kwTABLET {
			endLoc = p.cur.Loc
			p.advance()
		}
		// tablet_id — consume remaining as args
		endLoc = p.collectRawArgs(stmt, endLoc)
	case "COMPACT":
		// COMPACT TABLE name [WHERE ...]
		stmt.Object = "TABLE"
		if p.cur.Kind == kwTABLE {
			endLoc = p.cur.Loc
			p.advance()
		}
		name, nameLoc, nameErr := p.parseMultipartIdentifierOptional()
		if nameErr == nil && name != nil {
			stmt.Target = name
			endLoc = nameLoc
		}
		endLoc = p.collectRawArgs(stmt, endLoc)
	case "CHECK":
		// CHECK TABLET (id, ...) PROPERTIES(...)
		stmt.Object = "TABLET"
		if p.cur.Kind == kwTABLET {
			endLoc = p.cur.Loc
			p.advance()
		}
		// consume tablet id list and optional PROPERTIES
		endLoc = p.collectRawArgs(stmt, endLoc)
		if p.cur.Kind == kwPROPERTIES {
			props, propErr := p.parseProperties()
			if propErr == nil {
				stmt.Properties = props
				if len(props) > 0 {
					endLoc = ast.NodeLoc(props[len(props)-1])
				}
			}
		}
	case "REPAIR":
		// REPAIR TABLE name [PARTITION(p)]
		stmt.Object = "TABLE"
		if p.cur.Kind == kwTABLE {
			endLoc = p.cur.Loc
			p.advance()
		}
		name, nameLoc, nameErr := p.parseMultipartIdentifierOptional()
		if nameErr == nil && name != nil {
			stmt.Target = name
			endLoc = nameLoc
		}
		endLoc = p.collectRawArgs(stmt, endLoc)
	case "CANCEL":
		endLoc = p.parseAdminCancelBody(stmt, endLoc)
	case "SET":
		endLoc = p.parseAdminSetBody(stmt, endLoc)
	case "CLEAN":
		// CLEAN TRASH [ON (...)]
		stmt.Object = "TRASH"
		if p.cur.Kind == kwTRASH {
			endLoc = p.cur.Loc
			p.advance()
		}
		endLoc = p.skipAdminOnClause(endLoc)
	case "COPY":
		// COPY TABLET tablet_id PROPERTIES(...)
		stmt.Object = "TABLET"
		if p.cur.Kind == kwTABLET {
			endLoc = p.cur.Loc
			p.advance()
		}
		endLoc = p.collectRawArgs(stmt, endLoc)
		if p.cur.Kind == kwPROPERTIES {
			props, propErr := p.parseProperties()
			if propErr == nil {
				stmt.Properties = props
				if len(props) > 0 {
					endLoc = ast.NodeLoc(props[len(props)-1])
				}
			}
		}
	case "DECOMMISSION":
		// DECOMMISSION BACKEND BY HOSTNAME 'host:port'
		stmt.Object = "BACKEND"
		if p.cur.Kind == kwBACKEND {
			endLoc = p.cur.Loc
			p.advance()
		}
		endLoc = p.collectRawArgs(stmt, endLoc)
	default:
		// Unknown ADMIN verb — collect everything as args.
		endLoc = p.collectRawArgs(stmt, endLoc)
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseAdminVerb consumes and returns the verb token after ADMIN.
// Verb tokens are keywords (SHOW, REBALANCE, DIAGNOSE, etc.) or bare idents.
func (p *Parser) parseAdminVerb() (string, ast.Loc, error) {
	tok := p.cur
	if tok.Kind == tokEOF {
		return "", ast.Loc{}, p.syntaxErrorAtCur()
	}
	p.advance()
	v := strings.ToUpper(tok.Str)
	if v == "" {
		v = strings.ToUpper(TokenName(tok.Kind))
	}
	return v, tok.Loc, nil
}

// parseAdminShowBody parses the body of ADMIN SHOW ... and fills stmt fields.
// Returns the updated endLoc.
func (p *Parser) parseAdminShowBody(stmt *ast.AdminStmt, endLoc ast.Loc) ast.Loc {
	// Next token determines the object.
	obj := p.peekKeywordStr()
	switch obj {
	case "REPLICA":
		stmt.Object = "REPLICA"
		endLoc = p.cur.Loc
		p.advance() // consume REPLICA
		// sub-verb: DISTRIBUTION or STATUS
		switch p.peekKeywordStr() {
		case "DISTRIBUTION":
			endLoc = p.cur.Loc
			p.advance() // consume DISTRIBUTION
			// optional FROM table [PARTITION p]
			if p.cur.Kind == kwFROM {
				endLoc = p.cur.Loc
				p.advance()
				name, nameLoc, err := p.parseMultipartIdentifierOptional()
				if err == nil && name != nil {
					stmt.Target = name
					endLoc = nameLoc
				}
			}
			endLoc = p.collectRawArgs(stmt, endLoc)
		case "STATUS":
			endLoc = p.cur.Loc
			p.advance() // consume STATUS
			if p.cur.Kind == kwFROM {
				endLoc = p.cur.Loc
				p.advance()
				name, nameLoc, err := p.parseMultipartIdentifierOptional()
				if err == nil && name != nil {
					stmt.Target = name
					endLoc = nameLoc
				}
			}
			endLoc = p.collectRawArgs(stmt, endLoc)
		default:
			endLoc = p.collectRawArgs(stmt, endLoc)
		}
	case "CONFIG":
		stmt.Object = "CONFIG"
		endLoc = p.cur.Loc
		p.advance() // consume CONFIG
		endLoc = p.collectRawArgs(stmt, endLoc)
	case "FRONTEND":
		endLoc = p.cur.Loc
		p.advance() // consume FRONTEND
		// ADMIN SHOW FRONTEND CONFIG [LIKE ...]
		stmt.Object = "FRONTEND CONFIG"
		if p.peekKeywordStr() == "CONFIG" {
			endLoc = p.cur.Loc
			p.advance() // consume CONFIG
		}
		endLoc = p.collectRawArgs(stmt, endLoc)
	case "TABLET":
		stmt.Object = "TABLET"
		endLoc = p.cur.Loc
		p.advance() // consume TABLET
		endLoc = p.collectRawArgs(stmt, endLoc)
	case "DATA":
		stmt.Object = "DATA"
		endLoc = p.cur.Loc
		p.advance() // consume DATA
		// ADMIN SHOW DATA SKEW FROM table
		endLoc = p.collectRawArgs(stmt, endLoc)
	default:
		endLoc = p.collectRawArgs(stmt, endLoc)
	}
	return endLoc
}

// parseAdminCancelBody parses ADMIN CANCEL ... and fills stmt fields.
func (p *Parser) parseAdminCancelBody(stmt *ast.AdminStmt, endLoc ast.Loc) ast.Loc {
	obj := p.peekKeywordStr()
	switch obj {
	case "REBALANCE":
		stmt.Object = "REBALANCE DISK"
		endLoc = p.cur.Loc
		p.advance() // consume REBALANCE
		if p.cur.Kind == kwDISK {
			endLoc = p.cur.Loc
			p.advance()
		}
		endLoc = p.skipAdminOnClause(endLoc)
	case "REPAIR":
		stmt.Object = "REPAIR TABLE"
		endLoc = p.cur.Loc
		p.advance() // consume REPAIR
		if p.cur.Kind == kwTABLE {
			endLoc = p.cur.Loc
			p.advance()
		}
		name, nameLoc, nameErr := p.parseMultipartIdentifierOptional()
		if nameErr == nil && name != nil {
			stmt.Target = name
			endLoc = nameLoc
		}
		endLoc = p.collectRawArgs(stmt, endLoc)
	default:
		endLoc = p.collectRawArgs(stmt, endLoc)
	}
	return endLoc
}

// parseAdminSetBody parses ADMIN SET ... and fills stmt fields.
func (p *Parser) parseAdminSetBody(stmt *ast.AdminStmt, endLoc ast.Loc) ast.Loc {
	obj := p.peekKeywordStr()
	switch obj {
	case "REPLICA":
		stmt.Object = "REPLICA STATUS"
		endLoc = p.cur.Loc
		p.advance() // consume REPLICA
		// STATUS
		if p.peekKeywordStr() == "STATUS" {
			endLoc = p.cur.Loc
			p.advance()
		}
		if p.cur.Kind == kwPROPERTIES {
			props, propErr := p.parseProperties()
			if propErr == nil {
				stmt.Properties = props
				if len(props) > 0 {
					endLoc = ast.NodeLoc(props[len(props)-1])
				}
			}
		}
	case "TABLE":
		stmt.Object = "TABLE STATUS"
		endLoc = p.cur.Loc
		p.advance() // consume TABLE
		name, nameLoc, nameErr := p.parseMultipartIdentifierOptional()
		if nameErr == nil && name != nil {
			stmt.Target = name
			endLoc = nameLoc
		}
		// STATUS PROPERTIES(...)
		if p.peekKeywordStr() == "STATUS" {
			endLoc = p.cur.Loc
			p.advance()
		}
		if p.cur.Kind == kwPROPERTIES {
			props, propErr := p.parseProperties()
			if propErr == nil {
				stmt.Properties = props
				if len(props) > 0 {
					endLoc = ast.NodeLoc(props[len(props)-1])
				}
			}
		}
	case "PARTITION":
		stmt.Object = "PARTITION VERSION"
		endLoc = p.cur.Loc
		p.advance() // consume PARTITION
		endLoc = p.collectRawArgs(stmt, endLoc)
		if p.cur.Kind == kwPROPERTIES {
			props, propErr := p.parseProperties()
			if propErr == nil {
				stmt.Properties = props
				if len(props) > 0 {
					endLoc = ast.NodeLoc(props[len(props)-1])
				}
			}
		}
	case "FRONTEND":
		// ADMIN SET FRONTEND CONFIG ("key"="value")
		stmt.Object = "FRONTEND CONFIG"
		endLoc = p.cur.Loc
		p.advance() // consume FRONTEND
		if p.peekKeywordStr() == "CONFIG" {
			endLoc = p.cur.Loc
			p.advance() // consume CONFIG
		}
		// The config is a ("key"="value") list — reuse parseProperties with '(' already expected
		if p.cur.Kind == int('(') {
			// parseProperties expects the PROPERTIES keyword has been consumed;
			// here we manually parse the parens list.
			props, propErr := p.parseFrontendConfigList()
			if propErr == nil {
				stmt.Properties = props
				if len(props) > 0 {
					endLoc = ast.NodeLoc(props[len(props)-1])
				}
			}
		}
	default:
		endLoc = p.collectRawArgs(stmt, endLoc)
	}
	return endLoc
}

// parseFrontendConfigList parses ("key"="value") without consuming a leading
// PROPERTIES keyword. cur is '(' on entry.
func (p *Parser) parseFrontendConfigList() ([]*ast.Property, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	var props []*ast.Property
	for p.cur.Kind != int(')') && p.cur.Kind != tokEOF {
		startLoc := p.cur.Loc
		// Key — string literal
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		key := p.cur.Str
		p.advance()
		if _, err := p.expect(int('=')); err != nil {
			return nil, err
		}
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		val := p.cur.Str
		endLoc := p.cur.Loc
		p.advance()
		props = append(props, &ast.Property{
			Key:   key,
			Value: val,
			Loc:   ast.Loc{Start: startLoc.Start, End: endLoc.End},
		})
		if p.cur.Kind == int(',') {
			p.advance()
		}
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	return props, nil
}

// skipAdminOnClause skips an optional ON ('be1', 'be2') clause.
// Returns the updated endLoc.
func (p *Parser) skipAdminOnClause(endLoc ast.Loc) ast.Loc {
	if p.cur.Kind != kwON {
		return endLoc
	}
	endLoc = p.cur.Loc
	p.advance() // consume ON
	if p.cur.Kind != int('(') {
		return endLoc
	}
	// consume the entire parenthesised list
	endLoc = p.cur.Loc
	p.advance() // consume '('
	depth := 1
	for p.cur.Kind != tokEOF && depth > 0 {
		switch p.cur.Kind {
		case int('('):
			depth++
		case int(')'):
			depth--
		}
		endLoc = p.cur.Loc
		p.advance()
	}
	return endLoc
}

// peekKeywordStr returns the string form of the current token, upper-cased.
// Used internally for switch dispatch without consuming.
func (p *Parser) peekKeywordStr() string {
	tok := p.cur
	s := strings.ToUpper(tok.Str)
	if s == "" {
		s = strings.ToUpper(TokenName(tok.Kind))
	}
	return s
}

// collectRawArgs collects all remaining tokens (up to EOF) into stmt.Args as
// space-separated text. Stops at EOF.  Does not consume PROPERTIES or other
// structured tails that callers want to handle separately.
// Returns the updated endLoc (last non-EOF token seen).
func (p *Parser) collectRawArgs(stmt *ast.AdminStmt, endLoc ast.Loc) ast.Loc {
	var parts []string
	for p.cur.Kind != tokEOF {
		text := p.cur.Str
		if text == "" {
			text = TokenName(p.cur.Kind)
		}
		parts = append(parts, text)
		endLoc = p.cur.Loc
		p.advance()
	}
	if len(parts) > 0 {
		if stmt.Args != "" {
			stmt.Args += " "
		}
		stmt.Args += strings.Join(parts, " ")
	}
	return endLoc
}

// parseMultipartIdentifierOptional tries to parse a multipart identifier.
// Returns (nil, zero, nil) gracefully if the current token is not an identifier.
func (p *Parser) parseMultipartIdentifierOptional() (*ast.ObjectName, ast.Loc, error) {
	if !isIdentifierToken(p.cur.Kind) {
		return nil, ast.Loc{}, nil
	}
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, ast.Loc{}, err
	}
	return name, name.Loc, nil
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

	switch strings.ToUpper(verb) {
	case "ADD":
		endLoc = p.parseSystemAlterAddBody(stmt, endLoc)
	case "DROP":
		endLoc = p.parseSystemAlterDropBody(stmt, endLoc)
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
				props, propErr := p.parseFrontendConfigList()
				if propErr == nil {
					stmt.SetClause = props
					if len(props) > 0 {
						endLoc = ast.NodeLoc(props[len(props)-1])
					}
				}
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
		if p.cur.Kind == kwPROPERTIES {
			props, propErr := p.parseProperties()
			if propErr == nil {
				stmt.Properties = props
				if len(props) > 0 {
					endLoc = ast.NodeLoc(props[len(props)-1])
				}
			}
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

// parseSystemAlterAddBody parses the body of ALTER SYSTEM ADD ...
func (p *Parser) parseSystemAlterAddBody(stmt *ast.SystemAlterStmt, endLoc ast.Loc) ast.Loc {
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
		if p.cur.Kind == kwPROPERTIES {
			props, propErr := p.parseProperties()
			if propErr == nil {
				stmt.Properties = props
				if len(props) > 0 {
					endLoc = ast.NodeLoc(props[len(props)-1])
				}
			}
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
		brokerName, brokerLoc, brokerErr := p.parseIdentifierOrString()
		if brokerErr == nil {
			stmt.BrokerName = brokerName
			endLoc = brokerLoc
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
	return endLoc
}

// parseSystemAlterDropBody parses the body of ALTER SYSTEM DROP ...
func (p *Parser) parseSystemAlterDropBody(stmt *ast.SystemAlterStmt, endLoc ast.Loc) ast.Loc {
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
		brokerName, brokerLoc, brokerErr := p.parseIdentifierOrString()
		if brokerErr == nil {
			stmt.BrokerName = brokerName
			endLoc = brokerLoc
		}
		return endLoc
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
		brokerName, brokerLoc, brokerErr := p.parseIdentifierOrString()
		if brokerErr == nil {
			stmt.BrokerName = brokerName
			endLoc = brokerLoc
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
	return endLoc
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
