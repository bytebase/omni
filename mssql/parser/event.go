// Package parser - event.go implements T-SQL EVENT NOTIFICATION and EVENT SESSION statement parsing.
//
// The clauses follow ScriptDOM's grammar (TSql160.g:7586-8292, 15594-15633)
// and, where SQL Server 2022 differs, the pinned server under SET PARSEONLY
// ON: it takes event, action and target names of one part, refuses
// MAX_DURATION (SQL Server 2025's), and checks event notification types
// against its catalog.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// evIsWord reports whether the current token is the unquoted word word
// (upper case). A bracketed or quoted identifier never matches.
func (p *Parser) evIsWord(word string) bool {
	if kw := lookupKeyword(word); kw != tokIDENT {
		return p.cur.Type == kw
	}
	return p.cur.Type == tokIDENT && !p.isDelimited(p.cur) && strings.ToUpper(p.cur.Str) == word
}

// evExpectWordIn consumes one of the unquoted words words and returns it.
func (p *Parser) evExpectWordIn(words ...string) (string, error) {
	for _, w := range words {
		if p.evIsWord(w) {
			p.advance()
			return w, nil
		}
	}
	return "", p.syntaxErrorAtCur()
}

// evExpectIdentValue consumes an identifier, bracketed or not, whose text is
// one of words, and returns it. SQL Server reads EVENT_RETENTION_MODE and
// MEMORY_PARTITION_MODE values as identifiers (TSql160.g:8071, 8080).
func (p *Parser) evExpectIdentValue(words ...string) (string, error) {
	if p.cur.Type == tokIDENT || p.cur.Type >= kwABSENT {
		upper := strings.ToUpper(p.cur.Str)
		for _, w := range words {
			if upper == w {
				p.advance()
				return w, nil
			}
		}
	}
	return "", p.syntaxErrorAtCur()
}

// evExpectName consumes a one-part name: an identifier or a word that is not
// reserved.
func (p *Parser) evExpectName() (string, error) {
	name, ok := p.parseIdentifier()
	if !ok {
		return "", p.syntaxErrorAtCur()
	}
	return name, nil
}

// evExpectObjectName consumes a name of one to three parts with no empty
// part: [event_module_guid].event_package_name.event_name and the like. SQL
// Server takes one part too, where ScriptDOM wants two (TSql160.g:30433).
func (p *Parser) evExpectObjectName() (string, int, error) {
	var parts []string
	for {
		part, err := p.evExpectName()
		if err != nil {
			return "", 0, err
		}
		parts = append(parts, part)
		if p.cur.Type != '.' || len(parts) == 3 {
			break
		}
		p.advance()
	}
	return strings.Join(parts, "."), len(parts), nil
}

// evExpectQueueName consumes [ database . [ schema ] . | schema . ] queue.
func (p *Parser) evExpectQueueName() (string, error) {
	first, err := p.evExpectName()
	if err != nil {
		return "", err
	}
	parts := []string{first}
	for p.cur.Type == '.' && len(parts) < 3 {
		p.advance()
		if p.cur.Type == '.' && len(parts) == 1 {
			// database..queue: the schema is left out.
			parts = append(parts, "")
			continue
		}
		part, err := p.evExpectName()
		if err != nil {
			return "", err
		}
		parts = append(parts, part)
	}
	return strings.Join(parts, "."), nil
}

// evExpectString consumes a string literal and returns its value.
func (p *Parser) evExpectString() (string, error) {
	if p.cur.Type != tokSCONST && p.cur.Type != tokNSCONST {
		return "", p.syntaxErrorAtCur()
	}
	s := p.cur.Str
	p.advance()
	return s, nil
}

// evExpectInt consumes an unsigned integer that fits an int (32-bit).
func (p *Parser) evExpectInt() (string, error) {
	isBinary := len(p.cur.Str) > 1 && (p.cur.Str[1] == 'x' || p.cur.Str[1] == 'X')
	if p.cur.Type != tokICONST || isBinary || p.cur.Ival > 2147483647 {
		return "", p.syntaxErrorAtCur()
	}
	s := p.cur.Str
	p.advance()
	return s, nil
}

// parseCreateEventNotificationStmt parses CREATE EVENT NOTIFICATION.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-event-notification-transact-sql
// TSql160.g:7586
//
//	CREATE EVENT NOTIFICATION event_notification_name
//	ON { SERVER | DATABASE | QUEUE queue_name }
//	[ WITH FAN_IN ]
//	FOR { event_type | event_group } [ , ...n ]
//	TO SERVICE 'broker_service' , { 'broker_instance_specifier' | 'current database' }
//
// SQL Server takes only the event types and groups of
// sys.event_notification_event_types (1084), each once (1095), and refuses
// an empty broker service or instance string (1093).
func (p *Parser) parseCreateEventNotificationStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "EVENT NOTIFICATION",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}
	name, err := p.evExpectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	opt := &nodes.EventNotificationOption{Loc: nodes.Loc{Start: p.pos(), End: -1}}
	if err := p.parseEventNotificationScope(opt); err != nil {
		return nil, err
	}
	if p.cur.Type == kwWITH {
		p.advance()
		if _, err := p.evExpectWordIn("FAN_IN"); err != nil {
			return nil, err
		}
		opt.FanIn = true
	}
	if _, err := p.expect(kwFOR); err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	for {
		pos := p.pos()
		if (p.cur.Type != tokIDENT && p.cur.Type < kwABSENT) || p.isDelimited(p.cur) {
			return nil, p.syntaxErrorAtCur()
		}
		event := strings.ToUpper(p.cur.Str)
		if _, ok := evNotificationEventGroups[event]; !ok {
			return nil, p.newParseError(pos, "'"+p.cur.Str+"' is an invalid event type")
		}
		if evEventOverlaps(event, seen) {
			return nil, p.newParseError(pos, "\""+event+"\" has already been specified as an event type")
		}
		seen[event] = true
		p.advance()
		opt.Events = append(opt.Events, event)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if _, err := p.expect(kwTO); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwSERVICE); err != nil {
		return nil, err
	}
	for i, dst := range []*string{&opt.ServiceName, &opt.BrokerInstance} {
		if i == 1 {
			if _, err := p.expect(','); err != nil {
				return nil, err
			}
		}
		pos := p.pos()
		s, err := p.evExpectString()
		if err != nil {
			return nil, err
		}
		if s == "" {
			return nil, p.newParseError(pos, "\"\" is not a valid broker name")
		}
		*dst = s
	}
	opt.Loc.End = p.prevEnd()
	stmt.Options = &nodes.List{Items: []nodes.Node{opt}}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// evEventOverlaps reports whether event, or a group that holds it, or an
// event it holds, is among seen.
func evEventOverlaps(event string, seen map[string]bool) bool {
	for e := event; e != ""; e = evNotificationEventGroups[e] {
		if seen[e] {
			return true
		}
	}
	for s := range seen {
		for e := evNotificationEventGroups[s]; e != ""; e = evNotificationEventGroups[e] {
			if e == event {
				return true
			}
		}
	}
	return false
}

// parseEventNotificationScope parses ON { SERVER | DATABASE | QUEUE queue_name }
// (TSql160.g:7660).
func (p *Parser) parseEventNotificationScope(opt *nodes.EventNotificationOption) error {
	if _, err := p.expect(kwON); err != nil {
		return err
	}
	switch p.cur.Type {
	case kwSERVER:
		opt.Scope = "SERVER"
		p.advance()
	case kwDATABASE:
		opt.Scope = "DATABASE"
		p.advance()
	case kwQUEUE:
		opt.Scope = "QUEUE"
		p.advance()
		queue, err := p.evExpectQueueName()
		if err != nil {
			return err
		}
		opt.QueueName = queue
	default:
		return p.syntaxErrorAtCur()
	}
	return nil
}

// parseDropEventNotificationStmt parses DROP EVENT NOTIFICATION.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-event-notification-transact-sql
// TSql160.g:15594
//
//	DROP EVENT NOTIFICATION notification_name [ ,...n ]
//	ON { SERVER | DATABASE | QUEUE queue_name }
func (p *Parser) parseDropEventNotificationStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "DROP",
		ObjectType: "EVENT NOTIFICATION",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}
	name, err := p.evExpectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	opt := &nodes.EventNotificationOption{Loc: nodes.Loc{Start: p.pos(), End: -1}}
	for p.cur.Type == ',' {
		p.advance()
		extra, err := p.evExpectName()
		if err != nil {
			return nil, err
		}
		opt.ExtraNames = append(opt.ExtraNames, extra)
	}
	if err := p.parseEventNotificationScope(opt); err != nil {
		return nil, err
	}
	opt.Loc.End = p.prevEnd()
	stmt.Options = &nodes.List{Items: []nodes.Node{opt}}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCreateEventSessionStmt parses CREATE EVENT SESSION.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-event-session-transact-sql
// TSql160.g:7706
//
//	CREATE EVENT SESSION event_session_name
//	ON { SERVER | DATABASE }
//	<event_definition> [ , ...n ]
//	[ <event_target_definition> [ , ...n ] ]
//	[ WITH ( <event_session_options> [ , ...n ] ) ]
func (p *Parser) parseCreateEventSessionStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "EVENT SESSION",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}
	name, err := p.evExpectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	opts, err := p.parseEventSessionScope()
	if err != nil {
		return nil, err
	}
	if p.cur.Type != kwADD || p.peekNext().Type != kwEVENT {
		return nil, p.syntaxErrorAtCur()
	}
	decls, err := p.parseEventSessionDeclarations(kwADD, kwEVENT)
	if err != nil {
		return nil, err
	}
	opts = append(opts, decls...)
	if p.cur.Type == kwADD && p.peekNext().Type == kwTARGET {
		targets, err := p.parseEventSessionDeclarations(kwADD, kwTARGET)
		if err != nil {
			return nil, err
		}
		opts = append(opts, targets...)
	}
	if p.cur.Type == kwWITH {
		with, err := p.parseEventSessionOptions()
		if err != nil {
			return nil, err
		}
		opts = append(opts, with...)
	}
	stmt.Options = &nodes.List{Items: opts}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterEventSessionStmt parses ALTER EVENT SESSION.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-event-session-transact-sql
// TSql160.g:7725
//
//	ALTER EVENT SESSION event_session_name
//	ON { SERVER | DATABASE }
//	{
//	    { <add_event> [ , ...n ] | <add_target> [ , ...n ]
//	    | <drop_event> [ , ...n ] | <drop_target> [ , ...n ] }
//	    [ WITH ( <event_session_options> [ , ...n ] ) ]
//	  | WITH ( <event_session_options> [ , ...n ] )
//	  | STATE = { START | STOP }
//	}
//
// One statement adds or drops events, or targets, not both.
func (p *Parser) parseAlterEventSessionStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "ALTER",
		ObjectType: "EVENT SESSION",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}
	name, err := p.evExpectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	opts, err := p.parseEventSessionScope()
	if err != nil {
		return nil, err
	}
	switch {
	case p.cur.Type == kwADD || p.cur.Type == kwDROP:
		verb := p.cur.Type
		next := p.peekNext().Type
		if next != kwEVENT && next != kwTARGET {
			p.advance()
			return nil, p.syntaxErrorAtCur()
		}
		decls, err := p.parseEventSessionDeclarations(verb, next)
		if err != nil {
			return nil, err
		}
		opts = append(opts, decls...)
		if p.cur.Type == kwWITH {
			with, err := p.parseEventSessionOptions()
			if err != nil {
				return nil, err
			}
			opts = append(opts, with...)
		}
	case p.cur.Type == kwWITH:
		with, err := p.parseEventSessionOptions()
		if err != nil {
			return nil, err
		}
		opts = append(opts, with...)
	case p.cur.Type == kwSTATE:
		p.advance()
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		state, err := p.evExpectWordIn("START", "STOP")
		if err != nil {
			return nil, err
		}
		opts = append(opts, &nodes.String{Str: "STATE=" + state})
	default:
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Options = &nodes.List{Items: opts}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseDropEventSessionStmt parses DROP EVENT SESSION.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-event-session-transact-sql
// TSql160.g:15633
//
//	DROP EVENT SESSION event_session_name
//	ON { SERVER | DATABASE }
func (p *Parser) parseDropEventSessionStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "DROP",
		ObjectType: "EVENT SESSION",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}
	name, err := p.evExpectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	opts, err := p.parseEventSessionScope()
	if err != nil {
		return nil, err
	}
	stmt.Options = &nodes.List{Items: opts}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseEventSessionScope parses ON { SERVER | DATABASE } (TSql160.g:7688).
func (p *Parser) parseEventSessionScope() ([]nodes.Node, error) {
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	switch p.cur.Type {
	case kwSERVER:
		p.advance()
		return []nodes.Node{&nodes.String{Str: "ON=SERVER"}}, nil
	case kwDATABASE:
		p.advance()
		return []nodes.Node{&nodes.String{Str: "ON=DATABASE"}}, nil
	}
	return nil, p.syntaxErrorAtCur()
}

// parseEventSessionDeclarations parses a list of one kind of declaration,
// each starting with its own verb and noun: ADD EVENT x [ , ADD EVENT y ...]
// (TSql160.g:7797, 7816, 8273, 8292).
func (p *Parser) parseEventSessionDeclarations(verb, noun int) ([]nodes.Node, error) {
	prefix := "ADD "
	if verb == kwDROP {
		prefix = "DROP "
	}
	if noun == kwEVENT {
		prefix += "EVENT"
	} else {
		prefix += "TARGET"
	}
	var opts []nodes.Node
	for {
		if _, err := p.expect(verb); err != nil {
			return nil, err
		}
		if _, err := p.expect(noun); err != nil {
			return nil, err
		}
		name, _, err := p.evExpectObjectName()
		if err != nil {
			return nil, err
		}
		opts = append(opts, &nodes.String{Str: prefix + "=" + name})
		if verb == kwADD && p.cur.Type == '(' {
			var body []nodes.Node
			if noun == kwEVENT {
				body, err = p.parseEventSessionEventBody()
			} else {
				body, err = p.parseEventSessionTargetBody()
			}
			if err != nil {
				return nil, err
			}
			opts = append(opts, body...)
		}
		if p.cur.Type != ',' {
			return opts, nil
		}
		p.advance()
	}
}

// parseEventSessionEventBody parses an event declaration's parenthesized
// clauses, each optional, in this order (TSql160.g:7870):
//
//	( [ SET { event_customizable_attribute = <value> [ , ...n ] } ]
//	  [ ACTION ( { [event_module_guid].event_package_name.action_name [ , ...n ] } ) ]
//	  [ WHERE <predicate_expression> ] )
func (p *Parser) parseEventSessionEventBody() ([]nodes.Node, error) {
	p.advance() // consume '('
	var opts []nodes.Node
	if p.cur.Type == kwSET {
		sets, err := p.parseEventSessionSetList()
		if err != nil {
			return nil, err
		}
		opts = append(opts, sets...)
	}
	if p.cur.Type == kwACTION {
		p.advance()
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		for {
			name, _, err := p.evExpectObjectName()
			if err != nil {
				return nil, err
			}
			opts = append(opts, &nodes.String{Str: "ACTION=" + name})
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}
	if p.cur.Type == kwWHERE {
		p.advance()
		pred, err := p.parseEventSessionPredicate()
		if err != nil {
			return nil, err
		}
		opts = append(opts, &nodes.String{Str: "WHERE " + pred})
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return opts, nil
}

// parseEventSessionTargetBody parses a target declaration's
// ( SET { target_parameter_name = <value> [ , ...n ] } ) (TSql160.g:7943).
func (p *Parser) parseEventSessionTargetBody() ([]nodes.Node, error) {
	p.advance() // consume '('
	if p.cur.Type != kwSET {
		return nil, p.syntaxErrorAtCur()
	}
	opts, err := p.parseEventSessionSetList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return opts, nil
}

// parseEventSessionSetList parses SET name = <value> [ , ...n ]
// (TSql160.g:7897, 8089). SQL Server takes a bare name, not a bracketed one.
func (p *Parser) parseEventSessionSetList() ([]nodes.Node, error) {
	p.advance() // consume SET
	var opts []nodes.Node
	for {
		if p.isDelimited(p.cur) {
			return nil, p.syntaxErrorAtCur()
		}
		name, err := p.evExpectName()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		val, err := p.parseEventSessionValue()
		if err != nil {
			return nil, err
		}
		opts = append(opts, &nodes.String{Str: "SET " + name + "=" + val})
		if p.cur.Type != ',' {
			return opts, nil
		}
		p.advance()
	}
}

// parseEventSessionValue parses an event declaration value and returns its
// text (TSql160.g:8102):
//
//	<value> ::= - number | ( <value> ) | integer | numeric | real | binary | 'string'
func (p *Parser) parseEventSessionValue() (string, error) {
	switch p.cur.Type {
	case '-':
		p.advance()
		isBinary := len(p.cur.Str) > 1 && (p.cur.Str[1] == 'x' || p.cur.Str[1] == 'X')
		if (p.cur.Type != tokICONST && p.cur.Type != tokFCONST) || isBinary {
			return "", p.syntaxErrorAtCur()
		}
		val := "-" + p.cur.Str
		p.advance()
		return val, nil
	case '(':
		if err := p.enterNesting(); err != nil {
			return "", err
		}
		defer p.leaveNesting()
		p.advance()
		inner, err := p.parseEventSessionValue()
		if err != nil {
			return "", err
		}
		if _, err := p.expect(')'); err != nil {
			return "", err
		}
		return "(" + inner + ")", nil
	case tokICONST, tokFCONST, tokSCONST, tokNSCONST:
		val := p.cur.Str
		p.advance()
		return val, nil
	}
	return "", p.syntaxErrorAtCur()
}

// parseEventSessionPredicate parses an event's WHERE predicate and returns
// its text (TSql160.g:8142-8266):
//
//	<predicate_expression> ::= <and> [ OR <and> ...n ]
//	<and>    ::= <unary> [ AND <unary> ...n ]
//	<unary>  ::= NOT <unary> | ( <predicate_expression> ) | <leaf>
//	<leaf>   ::= <source> { = | <> | != | > | >= | !> | < | <= | !< | LIKE | NOT LIKE } <value>
//	           | package.predicate_compare_name ( <source> , <value> )
func (p *Parser) parseEventSessionPredicate() (string, error) {
	if err := p.enterNesting(); err != nil {
		return "", err
	}
	defer p.leaveNesting()
	var parts []string
	for {
		var ands []string
		for {
			term, err := p.parseEventSessionPredicateUnary()
			if err != nil {
				return "", err
			}
			ands = append(ands, term)
			if p.cur.Type != kwAND {
				break
			}
			p.advance()
		}
		parts = append(parts, strings.Join(ands, " AND "))
		if p.cur.Type != kwOR {
			return strings.Join(parts, " OR "), nil
		}
		p.advance()
	}
}

// parseEventSessionPredicateUnary parses NOT <unary>, a parenthesized
// predicate, or a leaf. A chain of NOT is read by a loop.
func (p *Parser) parseEventSessionPredicateUnary() (string, error) {
	prefix := ""
	for p.cur.Type == kwNOT {
		prefix += "NOT "
		p.advance()
	}
	if p.cur.Type == '(' {
		p.advance()
		inner, err := p.parseEventSessionPredicate()
		if err != nil {
			return "", err
		}
		if _, err := p.expect(')'); err != nil {
			return "", err
		}
		return prefix + "(" + inner + ")", nil
	}
	lhs, parts, err := p.evExpectObjectName()
	if err != nil {
		return "", err
	}
	if p.cur.Type == '(' {
		// A compare function names its package (TSql160.g:8220).
		if parts < 2 {
			return "", p.syntaxErrorAtCur()
		}
		p.advance()
		src, _, err := p.evExpectObjectName()
		if err != nil {
			return "", err
		}
		if _, err := p.expect(','); err != nil {
			return "", err
		}
		val, err := p.parseEventSessionValue()
		if err != nil {
			return "", err
		}
		if _, err := p.expect(')'); err != nil {
			return "", err
		}
		return prefix + lhs + "(" + src + ", " + val + ")", nil
	}
	var op string
	switch p.cur.Type {
	case '=':
		op = "="
	case tokNOTEQUAL:
		op = "<>"
	case tokLESSEQUAL:
		op = "<="
	case tokGREATEQUAL:
		op = ">="
	case tokNOTLESS:
		op = "!<"
	case tokNOTGREATER:
		op = "!>"
	case '<', '>':
		// SQL Server also reads < = and > = with a space between.
		op = string(rune(p.cur.Type))
		if p.peekNext().Type == '=' {
			p.advance()
			op += "="
		}
	case kwLIKE:
		op = "LIKE"
	case kwNOT:
		p.advance()
		if p.cur.Type != kwLIKE {
			return "", p.syntaxErrorAtCur()
		}
		op = "NOT LIKE"
	default:
		return "", p.syntaxErrorAtCur()
	}
	p.advance()
	val, err := p.parseEventSessionValue()
	if err != nil {
		return "", err
	}
	return prefix + lhs + " " + op + " " + val, nil
}

// parseEventSessionOptions parses WITH ( <event_session_options> [ , ...n ] )
// (TSql160.g:7835-8080).
//
//	<event_session_options> ::=
//	    MAX_MEMORY = size { KB | MB }
//	  | EVENT_RETENTION_MODE = { ALLOW_SINGLE_EVENT_LOSS | ALLOW_MULTIPLE_EVENT_LOSS | NO_EVENT_LOSS }
//	  | MAX_DISPATCH_LATENCY = { seconds SECONDS | INFINITE }
//	  | MAX_EVENT_SIZE = size { KB | MB }
//	  | MEMORY_PARTITION_MODE = { NONE | PER_NODE | PER_CPU }
//	  | TRACK_CAUSALITY = { ON | OFF }
//	  | STARTUP_STATE = { ON | OFF }
//
// SQL Server 2022 has no MAX_DURATION, and takes an option more than once.
func (p *Parser) parseEventSessionOptions() ([]nodes.Node, error) {
	p.advance() // consume WITH
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var opts []nodes.Node
	for {
		name, err := p.evExpectWordIn("MAX_MEMORY", "MAX_EVENT_SIZE", "TRACK_CAUSALITY", "STARTUP_STATE",
			"EVENT_RETENTION_MODE", "MEMORY_PARTITION_MODE", "MAX_DISPATCH_LATENCY")
		if err != nil {
			return nil, err
		}
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		var val string
		switch name {
		case "MAX_MEMORY", "MAX_EVENT_SIZE":
			var n, unit string
			if n, err = p.evExpectInt(); err == nil {
				unit, err = p.evExpectWordIn("KB", "MB")
				val = n + " " + unit
			}
		case "TRACK_CAUSALITY", "STARTUP_STATE":
			val, err = p.evExpectWordIn("ON", "OFF")
		case "EVENT_RETENTION_MODE":
			val, err = p.evExpectIdentValue("ALLOW_SINGLE_EVENT_LOSS", "ALLOW_MULTIPLE_EVENT_LOSS", "NO_EVENT_LOSS")
		case "MEMORY_PARTITION_MODE":
			val, err = p.evExpectIdentValue("NONE", "PER_NODE", "PER_CPU")
		case "MAX_DISPATCH_LATENCY":
			if p.evIsWord("INFINITE") {
				p.advance()
				val = "INFINITE"
			} else {
				var n string
				if n, err = p.evExpectInt(); err == nil {
					_, err = p.evExpectWordIn("SECONDS")
					val = n + " SECONDS"
				}
			}
		}
		if err != nil {
			return nil, err
		}
		opts = append(opts, &nodes.String{Str: name + "=" + val})
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return opts, nil
}
