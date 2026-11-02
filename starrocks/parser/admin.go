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
		endLoc, err = p.parseAdminShowBody(stmt, endLoc)
		if err != nil {
			return nil, err
		}
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
		if endLoc, err = p.parseAdminTarget(stmt, endLoc); err != nil {
			return nil, err
		}
		endLoc = p.collectRawArgs(stmt, endLoc)
	case "CHECK":
		// CHECK TABLET (id, ...) PROPERTIES(...)
		stmt.Object = "TABLET"
		if p.cur.Kind == kwTABLET {
			endLoc = p.cur.Loc
			p.advance()
		}
		// The tablet id list is raw args; the PROPERTIES tail is parsed.
		endLoc = p.collectRawArgsBefore(stmt, endLoc, kwPROPERTIES)
		if endLoc, err = p.parseAdminProperties(stmt, endLoc); err != nil {
			return nil, err
		}
	case "REPAIR":
		// REPAIR TABLE name [PARTITION(p)]
		stmt.Object = "TABLE"
		if p.cur.Kind == kwTABLE {
			endLoc = p.cur.Loc
			p.advance()
		}
		if endLoc, err = p.parseAdminTarget(stmt, endLoc); err != nil {
			return nil, err
		}
		endLoc = p.collectRawArgs(stmt, endLoc)
	case "CANCEL":
		if endLoc, err = p.parseAdminCancelBody(stmt, endLoc); err != nil {
			return nil, err
		}
	case "SET":
		if endLoc, err = p.parseAdminSetBody(stmt, endLoc); err != nil {
			return nil, err
		}
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
		endLoc = p.collectRawArgsBefore(stmt, endLoc, kwPROPERTIES)
		if endLoc, err = p.parseAdminProperties(stmt, endLoc); err != nil {
			return nil, err
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
func (p *Parser) parseAdminShowBody(stmt *ast.AdminStmt, endLoc ast.Loc) (ast.Loc, error) {
	var err error
	// Next token determines the object.
	obj := p.peekKeywordStr()
	switch obj {
	case "REPLICA":
		stmt.Object = "REPLICA"
		endLoc = p.cur.Loc
		p.advance() // consume REPLICA
		// sub-verb: DISTRIBUTION or STATUS
		switch p.peekKeywordStr() {
		case "DISTRIBUTION", "STATUS":
			endLoc = p.cur.Loc
			p.advance() // consume DISTRIBUTION / STATUS
			// optional FROM table [PARTITION p] / [WHERE ...]
			if p.cur.Kind == kwFROM {
				endLoc = p.cur.Loc
				p.advance()
				if endLoc, err = p.parseAdminTarget(stmt, endLoc); err != nil {
					return endLoc, err
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
	return endLoc, nil
}

// parseAdminCancelBody parses ADMIN CANCEL ... and fills stmt fields.
func (p *Parser) parseAdminCancelBody(stmt *ast.AdminStmt, endLoc ast.Loc) (ast.Loc, error) {
	var err error
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
		if endLoc, err = p.parseAdminTarget(stmt, endLoc); err != nil {
			return endLoc, err
		}
		endLoc = p.collectRawArgs(stmt, endLoc)
	default:
		endLoc = p.collectRawArgs(stmt, endLoc)
	}
	return endLoc, nil
}

// parseAdminSetBody parses ADMIN SET ... and fills stmt fields.
func (p *Parser) parseAdminSetBody(stmt *ast.AdminStmt, endLoc ast.Loc) (ast.Loc, error) {
	var err error
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
		if endLoc, err = p.parseAdminProperties(stmt, endLoc); err != nil {
			return endLoc, err
		}
	case "TABLE":
		stmt.Object = "TABLE STATUS"
		endLoc = p.cur.Loc
		p.advance() // consume TABLE
		if endLoc, err = p.parseAdminTarget(stmt, endLoc); err != nil {
			return endLoc, err
		}
		// STATUS PROPERTIES(...)
		if p.peekKeywordStr() == "STATUS" {
			endLoc = p.cur.Loc
			p.advance()
		}
		if endLoc, err = p.parseAdminProperties(stmt, endLoc); err != nil {
			return endLoc, err
		}
	case "PARTITION":
		stmt.Object = "PARTITION VERSION"
		endLoc = p.cur.Loc
		p.advance() // consume PARTITION
		endLoc = p.collectRawArgsBefore(stmt, endLoc, kwPROPERTIES)
		if endLoc, err = p.parseAdminProperties(stmt, endLoc); err != nil {
			return endLoc, err
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
		// The config is a bare ("key"="value") list, with no PROPERTIES
		// keyword in front of it.
		if p.cur.Kind == int('(') {
			props, err := p.parseFrontendConfigList()
			if err != nil {
				return endLoc, err
			}
			stmt.Properties = props
			endLoc = p.prev.Loc
		}
	default:
		endLoc = p.collectRawArgs(stmt, endLoc)
	}
	return endLoc, nil
}

// parseAdminTarget parses the optional table name an ADMIN statement names
// into stmt.Target, returning the updated endLoc. An error from the name is
// returned rather than dropped.
func (p *Parser) parseAdminTarget(stmt *ast.AdminStmt, endLoc ast.Loc) (ast.Loc, error) {
	name, nameLoc, err := p.parseMultipartIdentifierOptional()
	if err != nil {
		return endLoc, err
	}
	if name != nil {
		stmt.Target = name
		endLoc = nameLoc
	}
	return endLoc, nil
}

// parseAdminProperties parses an optional PROPERTIES(...) list at cur into
// stmt.Properties and returns the updated endLoc, which ends at the list's
// closing ')'. An error in the list is returned rather than dropped: at the
// end of the input nothing after it would catch it, so ADMIN SET REPLICA
// STATUS PROPERTIES( used to parse (container-verified reject).
func (p *Parser) parseAdminProperties(stmt *ast.AdminStmt, endLoc ast.Loc) (ast.Loc, error) {
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

// parseFrontendConfigList parses ("key"="value") without consuming a leading
// PROPERTIES keyword. cur is '(' on entry.
func (p *Parser) parseFrontendConfigList() ([]*ast.Property, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}

	// The engine takes exactly one config here -- ADMIN SET FRONTEND CONFIG
	// ("k"="v","k2"="v2") is a reject, where Doris accepts it (both
	// container-verified) -- which this still accepts. That is an arity
	// difference rather than a list-shape one, so it is left alone here.
	return p.parsePropertyItems(p.parseStringProperty)
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
// space-separated text. Returns the updated endLoc (last non-EOF token seen).
func (p *Parser) collectRawArgs(stmt *ast.AdminStmt, endLoc ast.Loc) ast.Loc {
	return p.collectRawArgsBefore(stmt, endLoc, tokEOF)
}

// collectRawArgsBefore is collectRawArgs stopping at the first stop token, so
// a structured tail the caller parses itself (the PROPERTIES list of ADMIN
// CHECK TABLET / COPY TABLET) is left for it. Collecting to EOF instead took
// that list as raw text, where nothing held it to the list's grammar.
func (p *Parser) collectRawArgsBefore(stmt *ast.AdminStmt, endLoc ast.Loc, stop TokenKind) ast.Loc {
	var parts []string
	for p.cur.Kind != tokEOF && p.cur.Kind != stop {
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
