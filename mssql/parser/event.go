// Package parser - event.go implements T-SQL EVENT NOTIFICATION and EVENT SESSION statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// eventSessionStateValues: ALTER EVENT SESSION ... STATE = { START | STOP }.
var eventSessionStateValues = newOptionSet().withIdents("START", "STOP")

// parseCreateEventNotificationStmt parses CREATE EVENT NOTIFICATION.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-event-notification-transact-sql
//
//	CREATE EVENT NOTIFICATION event_notification_name
//	ON { SERVER | DATABASE | QUEUE queue_name }
//	[ WITH FAN_IN ]
//	FOR { event_type | event_group } [ , ...n ]
//	TO SERVICE 'broker_service' , { 'broker_instance_specifier' | 'current database' }
func (p *Parser) parseCreateEventNotificationStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	// EVENT NOTIFICATION keywords already consumed by caller

	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "EVENT NOTIFICATION",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// event_notification_name
	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// Consume rest of statement (ON, WITH FAN_IN, FOR, TO SERVICE clauses)
	stmt.Options = p.parseEventNotificationOptions()

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseDropEventNotificationStmt parses DROP EVENT NOTIFICATION.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-event-notification-transact-sql
//
//	DROP EVENT NOTIFICATION notification_name [ ,...n ]
//	ON { SERVER | DATABASE | QUEUE queue_name }
func (p *Parser) parseDropEventNotificationStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	// EVENT NOTIFICATION keywords already consumed by caller

	stmt := &nodes.SecurityStmt{
		Action:     "DROP",
		ObjectType: "EVENT NOTIFICATION",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// notification_name [ ,...n ]
	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// Consume rest: optional additional names and ON clause
	stmt.Options = p.parseEventNotificationOptions()

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseEventNotificationOptions consumes ON/WITH/FOR/TO clauses of EVENT NOTIFICATION.
//
//	ON { SERVER | DATABASE | QUEUE queue_name }
//	[ WITH FAN_IN ]
//	FOR { event_type | event_group } [ , ...n ]
//	TO SERVICE 'broker_service' , { 'broker_instance_specifier' | 'current database' }
func (p *Parser) parseEventNotificationOptions() *nodes.List {
	loc := p.pos()
	opt := &nodes.EventNotificationOption{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	for p.cur.Type != ';' && p.cur.Type != tokEOF && p.cur.Type != kwGO {
		switch {
		case p.cur.Type == kwON:
			p.advance()
			if p.cur.Type == kwSERVER {
				opt.Scope = "SERVER"
				p.advance()
			} else if p.cur.Type == kwDATABASE {
				opt.Scope = "DATABASE"
				p.advance()
			} else if p.cur.Type == kwQUEUE {
				opt.Scope = "QUEUE"
				p.advance()
				queueName := ""
				if p.isIdentLike() || p.cur.Type == tokSCONST {
					queueName = p.cur.Str
					p.advance()
				}
				for p.cur.Type == '.' {
					p.advance()
					if p.isIdentLike() || p.cur.Type == tokSCONST {
						queueName += "." + p.cur.Str
						p.advance()
					}
				}
				opt.QueueName = queueName
			}
		case p.cur.Type == kwWITH:
			p.advance()
			if p.cur.Type == kwFAN_IN {
				opt.FanIn = true
				p.advance()
			}
		case p.cur.Type == kwFOR:
			p.advance()
			for {
				if p.isIdentLike() || p.cur.Type == tokSCONST {
					opt.Events = append(opt.Events, strings.ToUpper(p.cur.Str))
					p.advance()
				}
				if _, ok := p.match(','); !ok {
					break
				}
			}
		case p.cur.Type == kwTO:
			p.advance()
			if p.cur.Type == kwSERVICE {
				p.advance()
			}
			if p.cur.Type == tokSCONST {
				opt.ServiceName = p.cur.Str
				p.advance()
			}
			p.match(',')
			if p.cur.Type == tokSCONST {
				opt.BrokerInstance = p.cur.Str
				p.advance()
			}
		default:
			// Handle commas between notification names for DROP
			if p.cur.Type == ',' {
				p.advance()
				if p.isIdentLike() || p.cur.Type == tokSCONST {
					opt.ExtraNames = append(opt.ExtraNames, p.cur.Str)
					p.advance()
				}
			} else if p.cur.Type != tokEOF && p.cur.Type != ';' {
				p.advance()
			} else {
				break
			}
		}
	}

	opt.Loc.End = p.prevEnd()
	return &nodes.List{Items: []nodes.Node{opt}}
}

// parseCreateEventSessionStmt parses CREATE EVENT SESSION.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-event-session-transact-sql
//
//	CREATE EVENT SESSION event_session_name
//	ON { SERVER | DATABASE }
//	{
//	    <event_definition> [ , ...n ]
//	    [ <event_target_definition> [ , ...n ] ]
//	    [ WITH ( <event_session_options> [ , ...n ] ) ]
//	}
//
//	<event_definition>::=
//	    ADD EVENT [event_module_guid].event_package_name.event_name
//	         [ ( {
//	                 [ SET { event_customizable_attribute = <value> [ , ...n ] } ]
//	                 [ ACTION ( { [event_module_guid].event_package_name.action_name [ , ...n ] } ) ]
//	                 [ WHERE <predicate_expression> ]
//	        } ) ]
//
//	<event_target_definition>::=
//	    ADD TARGET [event_module_guid].event_package_name.target_name
//	        [ ( SET { target_parameter_name = <value> [ , ...n ] } ) ]
//
//	<event_session_options>::=
//	{
//	    [       MAX_MEMORY = size [ KB | MB ] ]
//	    [ [ , ] EVENT_RETENTION_MODE = { ALLOW_SINGLE_EVENT_LOSS | ALLOW_MULTIPLE_EVENT_LOSS | NO_EVENT_LOSS } ]
//	    [ [ , ] MAX_DISPATCH_LATENCY = { seconds SECONDS | INFINITE } ]
//	    [ [ , ] MAX_EVENT_SIZE = size [ KB | MB ] ]
//	    [ [ , ] MEMORY_PARTITION_MODE = { NONE | PER_NODE | PER_CPU } ]
//	    [ [ , ] TRACK_CAUSALITY = { ON | OFF } ]
//	    [ [ , ] STARTUP_STATE = { ON | OFF } ]
//	    [ [ , ] MAX_DURATION = { <time duration> { SECONDS | MINUTES | HOURS | DAYS } | UNLIMITED } ]
//	}
func (p *Parser) parseCreateEventSessionStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	// EVENT SESSION keywords already consumed by caller

	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "EVENT SESSION",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// event_session_name
	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// Consume rest of statement (ON, ADD EVENT, ADD TARGET, WITH)
	stmt.Options = p.parseEventSessionBody()

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterEventSessionStmt parses ALTER EVENT SESSION.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-event-session-transact-sql
//
//	ALTER EVENT SESSION event_session_name
//	ON { SERVER | DATABASE }
//	{
//	    [ { <add_drop_event> [ , ...n ] }
//	       | { <add_drop_event_target> [ , ...n ] } ]
//	    [ WITH ( <event_session_options> [ , ...n ] ) ]
//	    ]
//	    | [ STATE = { START | STOP } ]
//	}
//
//	<add_drop_event>::=
//	    ADD EVENT <event_specifier> [ ( ... ) ]
//	    | DROP EVENT <event_specifier>
//
//	<add_drop_event_target>::=
//	    ADD TARGET <event_target_specifier> [ ( SET ... ) ]
//	    | DROP TARGET <event_target_specifier>
func (p *Parser) parseAlterEventSessionStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	// EVENT SESSION keywords already consumed by caller

	stmt := &nodes.SecurityStmt{
		Action:     "ALTER",
		ObjectType: "EVENT SESSION",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// event_session_name
	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// Consume rest of statement
	stmt.Options = p.parseEventSessionBody()

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseDropEventSessionStmt parses DROP EVENT SESSION.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-event-session-transact-sql
//
//	DROP EVENT SESSION event_session_name
//	ON { SERVER | DATABASE }
func (p *Parser) parseDropEventSessionStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	// EVENT SESSION keywords already consumed by caller

	stmt := &nodes.SecurityStmt{
		Action:     "DROP",
		ObjectType: "EVENT SESSION",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// event_session_name
	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// ON { SERVER | DATABASE }
	stmt.Options = p.parseEventSessionBody()

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseEventSessionBody consumes ON, ADD EVENT/TARGET, DROP EVENT/TARGET, WITH, STATE clauses.
func (p *Parser) parseEventSessionBody() *nodes.List {
	var opts []nodes.Node

	// ON { SERVER | DATABASE }
	if p.cur.Type == kwON {
		p.advance()
		if p.cur.Type == kwSERVER {
			opts = append(opts, &nodes.String{Str: "ON=SERVER"})
			p.advance()
		} else if p.cur.Type == kwDATABASE {
			opts = append(opts, &nodes.String{Str: "ON=DATABASE"})
			p.advance()
		}
	}

	// Rest of body: ADD EVENT, ADD TARGET, DROP EVENT, DROP TARGET, STATE, WITH
	for p.cur.Type != ';' && p.cur.Type != tokEOF && p.cur.Type != kwGO {
		if p.cur.Type == kwADD {
			p.advance()
			if p.cur.Type == kwEVENT {
				p.advance()
				opts = append(opts, p.parseEventSessionEventSpec("ADD EVENT")...)
			} else if p.cur.Type == kwTARGET {
				p.advance()
				opts = append(opts, p.parseEventSessionTargetSpec("ADD TARGET")...)
			}
			// Check for comma then another ADD
			if p.cur.Type == ',' {
				p.advance()
				continue
			}
		} else if p.cur.Type == kwDROP {
			p.advance()
			if p.cur.Type == kwEVENT {
				p.advance()
				eventName := p.parseEventSessionDottedName()
				opts = append(opts, &nodes.String{Str: "DROP EVENT=" + eventName})
			} else if p.cur.Type == kwTARGET {
				p.advance()
				targetName := p.parseEventSessionDottedName()
				opts = append(opts, &nodes.String{Str: "DROP TARGET=" + targetName})
			}
			// Check for comma then another DROP/ADD
			if p.cur.Type == ',' {
				p.advance()
				continue
			}
		} else if p.cur.Type == kwSTATE {
			p.advance()
			p.match('=')
			// ALTER EVENT SESSION STATE = { START | STOP }
			if p.isValidOption(eventSessionStateValues) {
				opts = append(opts, &nodes.String{Str: "STATE=" + strings.ToUpper(p.cur.Str)})
				p.advance()
			}
		} else if p.cur.Type == kwWITH {
			p.advance()
			if p.cur.Type == '(' {
				p.advance()
				for p.cur.Type != ')' && p.cur.Type != tokEOF {
					if p.isIdentLike() || p.cur.Type == kwON || p.cur.Type == kwOFF {
						optName := strings.ToUpper(p.cur.Str)
						p.advance()
						if p.cur.Type == '=' {
							p.advance()
							// Value could be: number [KB|MB|SECONDS], identifier, ON/OFF
							val := ""
							if p.isIdentLike() || p.cur.Type == tokICONST || p.cur.Type == tokSCONST ||
								p.cur.Type == kwON || p.cur.Type == kwOFF {
								val = strings.ToUpper(p.cur.Str)
								p.advance()
							}
							// Check for unit suffix (KB, MB, SECONDS, etc.)
							if p.cur.Type == kwKB || p.cur.Type == kwMB ||
								p.cur.Type == kwSECONDS || p.cur.Type == kwMINUTES ||
								p.cur.Type == kwHOURS || p.cur.Type == kwDAYS {
								val += " " + strings.ToUpper(p.cur.Str)
								p.advance()
							}
							optName += "=" + val
						}
						opts = append(opts, &nodes.String{Str: optName})
					} else {
						p.advance()
					}
					p.match(',')
				}
				p.match(')')
			}
		} else {
			break
		}
	}

	if len(opts) == 0 {
		return nil
	}
	return &nodes.List{Items: opts}
}

// parseEventSessionEventSpec parses an ADD EVENT specifier with optional (SET/ACTION/WHERE).
//
//	ADD EVENT [event_module_guid].event_package_name.event_name
//	     [ ( {
//	             [ SET { event_customizable_attribute = <value> [ , ...n ] } ]
//	             [ ACTION ( { [event_module_guid].event_package_name.action_name [ , ...n ] } ) ]
//	             [ WHERE <predicate_expression> ]
//	    } ) ]
func (p *Parser) parseEventSessionEventSpec(prefix string) []nodes.Node {
	var result []nodes.Node

	eventName := p.parseEventSessionDottedName()
	result = append(result, &nodes.String{Str: prefix + "=" + eventName})

	// Optional parenthesized SET/ACTION/WHERE
	if p.cur.Type == '(' {
		p.advance()
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			if p.cur.Type == kwSET {
				// SET { attribute = value [ , ...n ] }
				p.advance()
				for {
					if !p.isIdentLike() {
						break
					}
					attrName := p.cur.Str
					p.advance()
					if p.cur.Type == '=' {
						p.advance()
						val := p.parseEventSessionValue()
						result = append(result, &nodes.String{Str: "SET " + attrName + "=" + val})
					}
					if _, ok := p.match(','); !ok {
						break
					}
					// Stop if next token is ACTION or WHERE (not another SET pair)
					if p.cur.Type == kwACTION || p.cur.Type == kwWHERE {
						break
					}
				}
			} else if p.cur.Type == kwACTION {
				// ACTION ( { package.action_name [ , ...n ] } )
				p.advance()
				if p.cur.Type == '(' {
					p.advance()
					for p.cur.Type != ')' && p.cur.Type != tokEOF {
						actionName := p.parseEventSessionDottedName()
						if actionName != "" {
							result = append(result, &nodes.String{Str: "ACTION=" + actionName})
						}
						if _, ok := p.match(','); !ok {
							break
						}
					}
					p.match(')')
				}
			} else if p.cur.Type == kwWHERE {
				// WHERE <predicate_expression>
				p.advance()
				predicate := p.parseEventSessionPredicate()
				if predicate != "" {
					result = append(result, &nodes.String{Str: "WHERE " + predicate})
				}
			} else {
				// Skip unexpected tokens
				p.advance()
			}
			// Optional comma between clauses
			p.match(',')
		}
		p.match(')')
	}

	return result
}

// parseEventSessionTargetSpec parses an ADD TARGET specifier with optional (SET ...).
//
//	ADD TARGET [event_module_guid].event_package_name.target_name
//	    [ ( SET { target_parameter_name = <value> [ , ...n ] } ) ]
func (p *Parser) parseEventSessionTargetSpec(prefix string) []nodes.Node {
	var result []nodes.Node

	targetName := p.parseEventSessionDottedName()
	result = append(result, &nodes.String{Str: prefix + "=" + targetName})

	// Optional parenthesized SET
	if p.cur.Type == '(' {
		p.advance()
		if p.cur.Type == kwSET {
			p.advance()
			for {
				if !p.isIdentLike() && p.cur.Type != tokSCONST {
					break
				}
				paramName := p.cur.Str
				p.advance()
				if p.cur.Type == '=' {
					p.advance()
					val := p.parseEventSessionValue()
					result = append(result, &nodes.String{Str: "SET " + paramName + "=" + val})
				}
				if _, ok := p.match(','); !ok {
					break
				}
			}
		}
		p.match(')')
	}

	return result
}

// parseEventSessionValue consumes a value: number, 'string', N'string', or identifier.
func (p *Parser) parseEventSessionValue() string {
	switch {
	case p.cur.Type == tokICONST:
		val := p.cur.Str
		p.advance()
		return val
	case p.cur.Type == tokSCONST:
		val := p.cur.Str
		p.advance()
		return val
	case p.cur.Type == tokNSCONST:
		val := p.cur.Str
		p.advance()
		return val
	case p.cur.Type == kwON:
		p.advance()
		return "ON"
	case p.cur.Type == kwOFF:
		p.advance()
		return "OFF"
	case p.isIdentLike():
		val := p.cur.Str
		p.advance()
		return val
	default:
		return ""
	}
}

// parseEventSessionPredicate consumes a simple predicate expression for event WHERE clauses.
// Handles: field_name { = | <> | != | > | >= | < | <= } value
// and: AND/OR connectors, NOT, parenthesized sub-expressions.
func (p *Parser) parseEventSessionPredicate() string {
	var parts []string
	for {
		part := p.parseEventSessionPredicateFactor()
		if part == "" {
			break
		}
		parts = append(parts, part)
		// AND / OR
		if p.cur.Type == kwAND {
			parts = append(parts, "AND")
			p.advance()
		} else if p.cur.Type == kwOR {
			parts = append(parts, "OR")
			p.advance()
		} else {
			break
		}
	}
	return strings.Join(parts, " ")
}

// parseEventSessionPredicateFactor parses a single predicate factor.
//
// BNF (from create-event-session-transact-sql.bnf):
//
//	<predicate_factor> ::=
//	    <predicate_leaf> | ( <predicate_expression> )
//
//	<predicate_leaf> ::=
//	      <predicate_source_declaration> { = | <> | != | > | >= | < | <= } <value>
//	    | [event_module_guid].event_package_name.predicate_compare_name
//	        ( <predicate_source_declaration> , <value> )
func (p *Parser) parseEventSessionPredicateFactor() string {
	// NOT
	notPrefix := ""
	if p.cur.Type == kwNOT {
		notPrefix = "NOT "
		p.advance()
	}

	// Parenthesized sub-expression
	if p.cur.Type == '(' {
		p.advance()
		inner := p.parseEventSessionPredicate()
		p.match(')')
		return notPrefix + "(" + inner + ")"
	}

	// field_name or package.predicate_source_name or package.predicate_compare_name
	lhs := p.parseEventSessionDottedName()
	if lhs == "" {
		return ""
	}

	// Check for function call form: predicate_compare_name(source, value)
	if p.cur.Type == '(' {
		p.advance()
		arg1 := p.parseEventSessionDottedName()
		p.match(',')
		arg2 := p.parseEventSessionValue()
		p.match(')')
		return notPrefix + lhs + "(" + arg1 + ", " + arg2 + ")"
	}

	// Comparison operator
	op := ""
	switch p.cur.Type {
	case '=':
		op = "="
		p.advance()
	case '<':
		p.advance()
		if p.cur.Type == '=' {
			op = "<="
			p.advance()
		} else {
			op = "<"
		}
	case '>':
		p.advance()
		if p.cur.Type == '=' {
			op = ">="
			p.advance()
		} else {
			op = ">"
		}
	case tokNOTEQUAL: // != or <>
		op = "!="
		p.advance()
	default:
		return notPrefix + lhs
	}

	// RHS value
	rhs := p.parseEventSessionValue()
	return notPrefix + lhs + " " + op + " " + rhs
}

// parseEventSessionDottedName consumes a dotted name like package.event_name or [guid].package.name.
func (p *Parser) parseEventSessionDottedName() string {
	name := ""
	if p.isIdentLike() || p.cur.Type == tokSCONST || p.cur.Type == tokICONST {
		name = p.cur.Str
		p.advance()
	}
	for p.cur.Type == '.' {
		p.advance()
		if p.isIdentLike() || p.cur.Type == tokSCONST || p.cur.Type == tokICONST {
			name += "." + p.cur.Str
			p.advance()
		}
	}
	return name
}
