// Package parser - service_broker.go implements T-SQL Service Broker statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// validationValues is the set of valid values for MESSAGE TYPE VALIDATION.
// BNF: VALIDATION = { NONE | EMPTY | WELL_FORMED_XML | VALID_XML WITH SCHEMA COLLECTION ... }
var validationValues = newOptionSet(kwNONE).withIdents("EMPTY", "WELL_FORMED_XML", "VALID_XML")

// brokerWithOptions is the set of valid WITH option names for service broker statements.
// Covers CREATE/ALTER ROUTE, REMOTE SERVICE BINDING, and generic broker WITH clauses.
var brokerWithOptions = newOptionSet(
	kwENCRYPTION,
	kwSTATUS,
).withIdents(
	"SERVICE_NAME", "BROKER_INSTANCE", "LIFETIME", "ADDRESS",
	"MIRROR_ADDRESS", "CLASSIFIER_FUNCTION",
	"USER", "ANONYMOUS",
)

// parseCreateMessageTypeStmt parses CREATE MESSAGE TYPE.
//
// BNF: mssql/parser/bnf/create-message-type-transact-sql.bnf
//
//	CREATE MESSAGE TYPE message_type_name
//	    [ AUTHORIZATION owner_name ]
//	    [ VALIDATION = {  NONE
//	                    | EMPTY
//	                    | WELL_FORMED_XML
//	                    | VALID_XML WITH SCHEMA COLLECTION schema_collection_name
//	                   } ]
func (p *Parser) parseCreateMessageTypeStmt() (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	// TYPE keyword already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "CREATE",
		ObjectType: "MESSAGE TYPE",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// Optional: AUTHORIZATION owner
	if p.cur.Type == kwAUTHORIZATION {
		p.advance()
		if p.isAnyKeywordIdent() {
			p.advance()
		}
	}

	// Optional: VALIDATION = { NONE | EMPTY | WELL_FORMED_XML | VALID_XML WITH SCHEMA COLLECTION ... }
	if p.cur.Type == kwVALIDATION {
		valLoc := p.pos()
		p.advance()
		if p.cur.Type == '=' {
			p.advance()
		}
		var opts []nodes.Node
		if p.isValidOption(validationValues) {
			valType := strings.ToUpper(p.cur.Str)
			p.advance()
			if valType == "VALID_XML" && p.cur.Type == kwWITH {
				p.advance() // consume WITH
				if p.cur.Type == kwSCHEMA {
					p.advance()
				}
				if p.cur.Type == kwCOLLECTION {
					p.advance()
				}
				// schema collection name (possibly dot-qualified)
				ref , _ := p.parseTableRef()
				schemaName := ""
				if ref != nil {
					if ref.Schema != "" {
						schemaName = ref.Schema + "." + ref.Object
					} else {
						schemaName = ref.Object
					}
				}
				opts = append(opts, &nodes.ServiceBrokerOption{Name: "VALIDATION", Value: "VALID_XML WITH SCHEMA COLLECTION " + schemaName, Loc: nodes.Loc{Start: valLoc, End: p.prevEnd()}})
			} else {
				opts = append(opts, &nodes.ServiceBrokerOption{Name: "VALIDATION", Value: valType, Loc: nodes.Loc{Start: valLoc, End: p.prevEnd()}})
			}
		} else {
			return nil, p.syntaxErrorAtCur()
		}
		if len(opts) > 0 {
			stmt.Options = &nodes.List{Items: opts}
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCreateContractStmt parses CREATE CONTRACT.
//
// BNF: mssql/parser/bnf/create-contract-transact-sql.bnf
//
//	CREATE CONTRACT contract_name
//	   [ AUTHORIZATION owner_name ]
//	      (  {   { message_type_name | [ DEFAULT ] }
//	          SENT BY { INITIATOR | TARGET | ANY }
//	       } [ ,...n] )
func (p *Parser) parseCreateContractStmt() (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	// CONTRACT keyword already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "CREATE",
		ObjectType: "CONTRACT",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	if p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// Optional: AUTHORIZATION owner
	if p.cur.Type == kwAUTHORIZATION {
		p.advance()
		if p.isAnyKeywordIdent() {
			p.advance()
		}
	}

	// Parse message type definitions: ( message_type_name SENT BY { INITIATOR | TARGET | ANY } [,...n] )
	if p.cur.Type == '(' {
		p.advance()
		var opts []nodes.Node
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			entryLoc := p.pos()
			// message_type_name or [DEFAULT]
			var msgTypeName string
			if p.cur.Type == kwDEFAULT {
				msgTypeName = "DEFAULT"
				p.advance()
			} else if p.isIdentLike() || p.cur.Type == tokSCONST {
				msgTypeName = p.cur.Str
				p.advance()
			}
			// SENT BY { INITIATOR | TARGET | ANY }
			sentBy := ""
			if p.cur.Type == kwSENT {
				p.advance()
				if p.cur.Type == kwBY {
					p.advance()
				}
				if p.isAnyKeywordIdent() || p.cur.Type == kwANY {
					sentBy = strings.ToUpper(p.cur.Str)
					p.advance()
				}
			}
			if msgTypeName != "" {
				opts = append(opts, &nodes.ServiceBrokerOption{Name: msgTypeName, Value: "SENT BY " + sentBy, Loc: nodes.Loc{Start: entryLoc, End: p.prevEnd()}})
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
		p.match(')')
		if len(opts) > 0 {
			stmt.Options = &nodes.List{Items: opts}
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCreateQueueStmt parses CREATE QUEUE.
//
// BNF: mssql/parser/bnf/create-queue-transact-sql.bnf
//
//	CREATE QUEUE <object>
//	   [ WITH
//	     [ STATUS = { ON | OFF } [ , ] ]
//	     [ RETENTION = { ON | OFF } [ , ] ]
//	     [ ACTIVATION (
//	         [ STATUS = { ON | OFF } , ]
//	           PROCEDURE_NAME = <procedure> ,
//	           MAX_QUEUE_READERS = max_readers ,
//	           EXECUTE AS { SELF | 'user_name' | OWNER }
//	            ) [ , ] ]
//	     [ POISON_MESSAGE_HANDLING (
//	         [ STATUS = { ON | OFF } ] ) ]
//	    ]
//	     [ ON { filegroup | [ DEFAULT ] } ]
//
//	<object> ::=
//	{ database_name.schema_name.queue_name | schema_name.queue_name | queue_name }
//
//	<procedure> ::=
//	{ database_name.schema_name.stored_procedure_name | schema_name.stored_procedure_name | stored_procedure_name }
func (p *Parser) parseCreateQueueStmt() (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	// QUEUE keyword already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "CREATE",
		ObjectType: "QUEUE",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// Parse possibly schema-qualified queue name
	ref , _ := p.parseTableRef()
	if ref != nil {
		if ref.Schema != "" {
			stmt.Name = ref.Schema + "." + ref.Object
		} else {
			stmt.Name = ref.Object
		}
	}

	var opts []nodes.Node
	opts = p.parseQueueWithClause(opts)

	// ON { filegroup | [DEFAULT] }
	if p.cur.Type == kwON {
		onLoc := p.pos()
		p.advance()
		if p.isAnyKeywordIdent() || p.cur.Type == kwDEFAULT {
			opts = append(opts, &nodes.ServiceBrokerOption{Name: "ON", Value: strings.ToUpper(p.cur.Str), Loc: nodes.Loc{Start: onLoc, End: -1}})
			p.advance()
		}
	}

	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseQueueWithClause parses the WITH clause options for CREATE/ALTER QUEUE.
func (p *Parser) parseQueueWithClause(opts []nodes.Node) []nodes.Node {
	if p.cur.Type != kwWITH {
		return opts
	}
	p.advance()

	for {
		if p.cur.Type == kwSTATUS {
			optLoc := p.pos()
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			if p.cur.Type == kwON || p.cur.Type == kwOFF {
				opts = append(opts, &nodes.ServiceBrokerOption{Name: "STATUS", Value: strings.ToUpper(p.cur.Str), Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
				p.advance()
			}
		} else if p.cur.Type == kwRETENTION {
			optLoc := p.pos()
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			if p.cur.Type == kwON || p.cur.Type == kwOFF {
				opts = append(opts, &nodes.ServiceBrokerOption{Name: "RETENTION", Value: strings.ToUpper(p.cur.Str), Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
				p.advance()
			}
		} else if p.cur.Type == kwACTIVATION {
			p.advance()
			opts = p.parseQueueActivationClause(opts)
		} else if p.cur.Type == kwPOISON_MESSAGE_HANDLING {
			optLoc := p.pos()
			p.advance()
			if p.cur.Type == '(' {
				p.advance()
				if p.cur.Type == kwSTATUS {
					p.advance()
					if p.cur.Type == '=' {
						p.advance()
					}
					if p.cur.Type == kwON || p.cur.Type == kwOFF {
						opts = append(opts, &nodes.ServiceBrokerOption{Name: "POISON_MESSAGE_HANDLING", Value: strings.ToUpper(p.cur.Str), Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
						p.advance()
					}
				}
				p.match(')')
			}
		} else {
			break
		}
		if _, ok := p.match(','); !ok {
			break
		}
	}
	return opts
}

// parseQueueActivationClause parses ACTIVATION ( ... ) options.
func (p *Parser) parseQueueActivationClause(opts []nodes.Node) []nodes.Node {
	if p.cur.Type != '(' {
		return opts
	}
	p.advance()

	// Collect activation options
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		optLoc := p.pos()
		if p.cur.Type == kwSTATUS {
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			if p.cur.Type == kwON || p.cur.Type == kwOFF {
				opts = append(opts, &nodes.ServiceBrokerOption{Name: "ACTIVATION:STATUS", Value: strings.ToUpper(p.cur.Str), Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
				p.advance()
			}
		} else if p.cur.Type == kwPROCEDURE_NAME {
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			// Parse possibly schema-qualified procedure name
			ref , _ := p.parseTableRef()
			if ref != nil {
				name := ref.Object
				if ref.Schema != "" {
					name = ref.Schema + "." + name
				}
				if ref.Database != "" {
					name = ref.Database + "." + name
				}
				opts = append(opts, &nodes.ServiceBrokerOption{Name: "ACTIVATION:PROCEDURE_NAME", Value: name, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
			}
		} else if p.cur.Type == kwMAX_QUEUE_READERS {
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			if p.cur.Type == tokICONST {
				opts = append(opts, &nodes.ServiceBrokerOption{Name: "ACTIVATION:MAX_QUEUE_READERS", Value: p.cur.Str, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
				p.advance()
			}
		} else if p.cur.Type == kwEXECUTE {
			p.advance()
			if p.cur.Type == kwAS {
				p.advance()
			}
			if p.cur.Type == kwSELF {
				opts = append(opts, &nodes.ServiceBrokerOption{Name: "ACTIVATION:EXECUTE AS", Value: "SELF", Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
				p.advance()
			} else if p.cur.Type == kwOWNER {
				opts = append(opts, &nodes.ServiceBrokerOption{Name: "ACTIVATION:EXECUTE AS", Value: "OWNER", Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
				p.advance()
			} else if p.cur.Type == tokSCONST {
				opts = append(opts, &nodes.ServiceBrokerOption{Name: "ACTIVATION:EXECUTE AS", Value: p.cur.Str, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
				p.advance()
			}
		} else if p.cur.Type == kwDROP {
			opts = append(opts, &nodes.ServiceBrokerOption{Name: "ACTIVATION:DROP", Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
			p.advance()
		} else {
			break
		}
		p.match(',')
	}
	p.match(')')

	return opts
}

// parseCreateServiceStmt parses CREATE SERVICE.
//
// BNF: mssql/parser/bnf/create-service-transact-sql.bnf
//
//	CREATE SERVICE service_name
//	   [ AUTHORIZATION owner_name ]
//	   ON QUEUE [ schema_name. ]queue_name
//	   [ ( contract_name | [DEFAULT][ ,...n ] ) ]
func (p *Parser) parseCreateServiceStmt() (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	// SERVICE keyword already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "CREATE",
		ObjectType: "SERVICE",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Options = p.parseServiceBrokerOptions()
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseSendStmt parses SEND ON CONVERSATION.
//
// BNF: mssql/parser/bnf/send-transact-sql.bnf
//
//	SEND
//	   ON CONVERSATION [(]conversation_handle [,.. @conversation_handle_n][)]
//	   [ MESSAGE TYPE message_type_name ]
//	   [ ( message_body_expression ) ]
func (p *Parser) parseSendStmt() (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	p.advance() // consume SEND

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "SEND",
		ObjectType: "CONVERSATION",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// ON CONVERSATION [(]handle [,...n][)]
	if p.cur.Type == kwON {
		p.advance()
		if p.cur.Type == kwCONVERSATION {
			p.advance()
		}
		// Multiple handles may be wrapped in parentheses
		if p.cur.Type == '(' {
			p.advance()
			// First handle
			if p.cur.Type == tokVARIABLE || p.isIdentLike() {
				stmt.Name = p.cur.Str
				p.advance()
			}
			// Additional handles
			for {
				if _, ok := p.match(','); !ok {
					break
				}
				if p.cur.Type == tokVARIABLE || p.isAnyKeywordIdent() {
					p.advance()
				}
			}
			p.match(')')
		} else if p.cur.Type == tokVARIABLE || p.isIdentLike() {
			stmt.Name = p.cur.Str
			p.advance()
		}
	}

	// MESSAGE TYPE message_type_name
	if p.cur.Type == kwMESSAGE {
		p.advance()
		if p.cur.Type == kwTYPE {
			p.advance()
		}
		if p.isAnyKeywordIdent() || p.cur.Type == tokVARIABLE {
			p.advance()
		}
	}

	// ( message_body_expression )
	if p.cur.Type == '(' {
		p.advance()
		p.parseExpr()
		p.match(')')
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseReceiveStmt parses RECEIVE.
//
// BNF: mssql/parser/bnf/receive-transact-sql.bnf
//
//	[ WAITFOR ( ]
//	RECEIVE [ TOP ( n ) ]
//	    <column_specifier> [ ,...n ]
//	    FROM <queue>
//	    [ INTO table_variable ]
//	    [ WHERE { conversation_handle = @conversation_handle | conversation_group_id = @conversation_group_id } ]
//	[ ) ] [ , TIMEOUT timeout ]
//
//	<column_specifier> ::=
//	{    *
//	  |  { column_name | [ ] expression } [ [ AS ] column_alias ]
//	}     [ ,...n ]
//
//	<queue> ::=
//	{ database_name.schema_name.queue_name | schema_name.queue_name | queue_name }
func (p *Parser) parseReceiveStmt() (*nodes.ReceiveStmt, error) {
	loc := p.pos()
	p.advance() // consume RECEIVE

	stmt := &nodes.ReceiveStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// TOP (n)
	if p.cur.Type == kwTOP {
		p.advance()
		if p.cur.Type == '(' {
			p.advance()
			stmt.Top, _ = p.parseExpr()
			p.match(')')
		}
	}

	// column list: * or column_name [AS alias] [,...n]
	if p.cur.Type == '*' {
		stmt.AllColumns = true
		p.advance()
	} else {
		// RECEIVE column list is terminated by a DML/statement boundary
		// (FROM / INTO / WHERE / ; / EOF), so parseCommaList's single-end
		// contract does not fit. Keep open-coded but enforce trailing-comma
		// rejection to match SqlScriptDOM strictness.
		var cols []nodes.Node
		for p.cur.Type != kwFROM && p.cur.Type != tokEOF && p.cur.Type != ';' {
			col := &nodes.ReceiveColumn{
				Loc: nodes.Loc{Start: p.pos(), End: -1},
			}
			if p.isAnyKeywordIdent() || p.cur.Type == tokVARIABLE {
				col.Expr = &nodes.ColumnRef{Column: p.cur.Str, Loc: nodes.Loc{Start: p.pos(), End: p.prevEnd() + len(p.cur.Str)}}
				p.advance()
			} else {
				col.Expr, _ = p.parseExpr()
			}
			if p.cur.Type == kwAS {
				p.advance()
				if p.isIdentLike() {
					col.Alias = p.cur.Str
					p.advance()
				}
			} else if p.isIdentLike() && p.cur.Type != kwFROM && p.cur.Type != kwINTO && p.cur.Type != kwWHERE {
				col.Alias = p.cur.Str
				p.advance()
			}
			col.Loc.End = p.prevEnd()
			if col.Expr != nil {
				cols = append(cols, col)
			}
			if _, ok := p.match(','); !ok {
				break
			}
			// After consuming ',' there must be another column before the
			// statement terminator — reject trailing comma.
			if p.cur.Type == kwFROM || p.cur.Type == tokEOF || p.cur.Type == ';' {
				return nil, p.unexpectedToken()
			}
		}
		if len(cols) > 0 {
			stmt.Columns = &nodes.List{Items: cols}
		}
	}

	// FROM queue
	if _, ok := p.match(kwFROM); ok {
		stmt.Queue , _ = p.parseTableRef()
	}

	// INTO table_variable
	if p.cur.Type == kwINTO {
		p.advance()
		if p.cur.Type == tokVARIABLE {
			stmt.IntoVar = p.cur.Str
			p.advance()
		} else if p.isIdentLike() {
			stmt.IntoVar = p.cur.Str
			p.advance()
		}
	}

	// WHERE clause
	if p.cur.Type == kwWHERE {
		p.advance()
		stmt.WhereClause, _ = p.parseExpr()
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseBeginConversationStmt parses BEGIN DIALOG CONVERSATION.
//
// BNF: mssql/parser/bnf/begin-dialog-conversation-transact-sql.bnf
//
//	BEGIN DIALOG [ CONVERSATION ] @dialog_handle
//	    FROM SERVICE initiator_service_name
//	    TO SERVICE 'target_service_name'
//	        [ , { 'service_broker_guid' | 'CURRENT DATABASE' }]
//	    [ ON CONTRACT contract_name ]
//	    [ WITH
//	        [ { RELATED_CONVERSATION = related_conversation_handle
//	          | RELATED_CONVERSATION_GROUP = related_conversation_group_id } ]
//	        [ [ , ] LIFETIME = dialog_lifetime ]
//	        [ [ , ] ENCRYPTION = { ON | OFF } ]
//	    ]
func (p *Parser) parseBeginConversationStmt() (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	// already consumed BEGIN, and DIALOG has been consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "BEGIN",
		ObjectType: "CONVERSATION",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// optional CONVERSATION keyword
	if p.cur.Type == kwCONVERSATION {
		p.advance()
	}

	// @dialog_handle
	if p.cur.Type == tokVARIABLE {
		stmt.Name = p.cur.Str
		p.advance()
	}

	var opts []nodes.Node

	// FROM SERVICE initiator_service_name
	if p.cur.Type == kwFROM {
		optLoc := p.pos()
		p.advance()
		if p.cur.Type == kwSERVICE {
			p.advance()
		}
		if p.isAnyKeywordIdent() || p.cur.Type == tokSCONST {
			opts = append(opts, &nodes.ServiceBrokerOption{Name: "FROM SERVICE", Value: p.cur.Str, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
			p.advance()
		}
	}

	// TO SERVICE 'target_service_name' [, { 'service_broker_guid' | 'CURRENT DATABASE' }]
	if p.cur.Type == kwTO {
		optLoc := p.pos()
		p.advance()
		if p.cur.Type == kwSERVICE {
			p.advance()
		}
		if p.cur.Type == tokSCONST || p.isAnyKeywordIdent() {
			opts = append(opts, &nodes.ServiceBrokerOption{Name: "TO SERVICE", Value: p.cur.Str, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
			p.advance()
		}
		// optional broker guid
		if _, ok := p.match(','); ok {
			biLoc := p.pos()
			if p.cur.Type == tokSCONST || p.isAnyKeywordIdent() {
				opts = append(opts, &nodes.ServiceBrokerOption{Name: "BROKER_INSTANCE", Value: p.cur.Str, Loc: nodes.Loc{Start: biLoc, End: p.prevEnd()}})
				p.advance()
			}
		}
	}

	// ON CONTRACT contract_name
	if p.cur.Type == kwON {
		optLoc := p.pos()
		p.advance()
		if p.cur.Type == kwCONTRACT {
			p.advance()
		}
		if p.isAnyKeywordIdent() || p.cur.Type == tokSCONST {
			opts = append(opts, &nodes.ServiceBrokerOption{Name: "ON CONTRACT", Value: p.cur.Str, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
			p.advance()
		}
	}

	// WITH options
	if p.cur.Type == kwWITH {
		p.advance()
		for {
			optLoc := p.pos()
			if p.cur.Type == kwRELATED_CONVERSATION_GROUP {
				p.advance()
				if p.cur.Type == '=' {
					p.advance()
				}
				if p.cur.Type == tokVARIABLE || p.isAnyKeywordIdent() || p.cur.Type == tokSCONST {
					opts = append(opts, &nodes.ServiceBrokerOption{Name: "RELATED_CONVERSATION_GROUP", Value: p.cur.Str, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
					p.advance()
				}
			} else if p.cur.Type == kwRELATED_CONVERSATION {
				p.advance()
				if p.cur.Type == '=' {
					p.advance()
				}
				if p.cur.Type == tokVARIABLE || p.isAnyKeywordIdent() || p.cur.Type == tokSCONST {
					opts = append(opts, &nodes.ServiceBrokerOption{Name: "RELATED_CONVERSATION", Value: p.cur.Str, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
					p.advance()
				}
			} else if p.cur.Type == kwLIFETIME {
				p.advance()
				if p.cur.Type == '=' {
					p.advance()
				}
				if p.cur.Type == tokICONST || p.cur.Type == tokVARIABLE || p.isAnyKeywordIdent() {
					opts = append(opts, &nodes.ServiceBrokerOption{Name: "LIFETIME", Value: p.cur.Str, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
					p.advance()
				}
			} else if p.cur.Type == kwENCRYPTION {
				p.advance()
				if p.cur.Type == '=' {
					p.advance()
				}
				if p.cur.Type == kwON || p.cur.Type == kwOFF {
					opts = append(opts, &nodes.ServiceBrokerOption{Name: "ENCRYPTION", Value: strings.ToUpper(p.cur.Str), Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
					p.advance()
				}
			} else {
				break
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
	}

	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseEndConversationStmt parses END CONVERSATION.
//
// BNF: mssql/parser/bnf/end-conversation-transact-sql.bnf
//
//	END CONVERSATION conversation_handle
//	    [   [ WITH ERROR = failure_code DESCRIPTION = 'failure_text' ]
//	      | [ WITH CLEANUP ]
//	    ]
func (p *Parser) parseEndConversationStmt() (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	p.advance() // consume END
	// consume CONVERSATION
	if p.cur.Type == kwCONVERSATION {
		p.advance()
	}

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "END",
		ObjectType: "CONVERSATION",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// conversation handle
	if p.cur.Type == tokVARIABLE || p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// WITH ERROR = failure_code DESCRIPTION = 'failure_text'
	// WITH CLEANUP
	if p.cur.Type == kwWITH {
		p.advance()
		var opts []nodes.Node
		if p.cur.Type == kwERROR {
			errLoc := p.pos()
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			// failure_code: integer or variable
			errorCode := ""
			if p.cur.Type == tokICONST || p.cur.Type == tokVARIABLE || p.isIdentLike() {
				errorCode = p.cur.Str
				p.advance()
			}
			opts = append(opts, &nodes.ServiceBrokerOption{Name: "ERROR", Value: errorCode, Loc: nodes.Loc{Start: errLoc, End: p.prevEnd()}})

			// DESCRIPTION = 'failure_text'
			if p.cur.Type == kwDESCRIPTION {
				descLoc := p.pos()
				p.advance()
				if p.cur.Type == '=' {
					p.advance()
				}
				descVal := ""
				if p.cur.Type == tokSCONST || p.cur.Type == tokVARIABLE || p.isIdentLike() {
					descVal = p.cur.Str
					p.advance()
				}
				opts = append(opts, &nodes.ServiceBrokerOption{Name: "DESCRIPTION", Value: descVal, Loc: nodes.Loc{Start: descLoc, End: p.prevEnd()}})
			}
		} else if p.cur.Type == kwCLEANUP {
			cleanLoc := p.pos()
			opts = append(opts, &nodes.ServiceBrokerOption{Name: "CLEANUP", Loc: nodes.Loc{Start: cleanLoc, End: p.prevEnd()}})
			p.advance()
		}
		if len(opts) > 0 {
			stmt.Options = &nodes.List{Items: opts}
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCreateRouteStmt parses CREATE ROUTE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-route-transact-sql
//
//	CREATE ROUTE route_name
//	    [ AUTHORIZATION owner_name ]
//	    WITH
//	    [ SERVICE_NAME = 'service_name' , ]
//	    [ BROKER_INSTANCE = 'broker_instance_identifier' , ]
//	    [ LIFETIME = route_lifetime , ]
//	    ADDRESS = 'next_hop_address'
//	    [ , MIRROR_ADDRESS = 'next_hop_mirror_address' ]
func (p *Parser) parseCreateRouteStmt() (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	// ROUTE keyword already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "CREATE",
		ObjectType: "ROUTE",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// route_name
	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// Optional: AUTHORIZATION owner
	if p.cur.Type == kwAUTHORIZATION {
		p.advance()
		if p.isAnyKeywordIdent() {
			p.advance()
		}
	}

	// WITH clause (required for CREATE ROUTE)
	if p.cur.Type == kwWITH {
		p.advance()
		var opts []nodes.Node
		for {
			if !p.isIdentLike() && p.cur.Type != tokSCONST {
				break
			}
			optLoc := p.pos()
			optName := strings.ToUpper(p.cur.Str)
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
				var val string
				if p.cur.Type == tokSCONST {
					val = p.cur.Str
					p.advance()
				} else if p.cur.Type == tokICONST {
					val = p.cur.Str
					p.advance()
				} else if p.isIdentLike() {
					val = p.cur.Str
					p.advance()
				}
				opts = append(opts, &nodes.ServiceBrokerOption{Name: optName, Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
			} else {
				opts = append(opts, &nodes.ServiceBrokerOption{Name: optName, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
		if len(opts) > 0 {
			stmt.Options = &nodes.List{Items: opts}
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCreateRemoteServiceBindingStmt parses CREATE REMOTE SERVICE BINDING.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-remote-service-binding-transact-sql
//
//	CREATE REMOTE SERVICE BINDING binding_name
//	    [ AUTHORIZATION owner_name ]
//	    TO SERVICE 'service_name'
//	    WITH USER = user_name [ , ANONYMOUS = { ON | OFF } ]
func (p *Parser) parseCreateRemoteServiceBindingStmt() (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	// REMOTE SERVICE BINDING keywords already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "CREATE",
		ObjectType: "REMOTE SERVICE BINDING",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// binding_name
	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// Optional: AUTHORIZATION owner
	if p.cur.Type == kwAUTHORIZATION {
		p.advance()
		if p.isAnyKeywordIdent() {
			p.advance()
		}
	}

	var opts []nodes.Node

	// TO SERVICE 'service_name'
	if p.cur.Type == kwTO {
		optLoc := p.pos()
		p.advance()
		if p.cur.Type == kwSERVICE {
			p.advance()
		}
		if p.cur.Type == tokSCONST {
			opts = append(opts, &nodes.ServiceBrokerOption{Name: "SERVICE", Value: p.cur.Str, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
			p.advance()
		}
	}

	// WITH USER = user_name [, ANONYMOUS = { ON | OFF }]
	if p.cur.Type == kwWITH {
		p.advance()
		for {
			if !p.isIdentLike() && p.cur.Type != kwUSER {
				break
			}
			optLoc := p.pos()
			optName := strings.ToUpper(p.cur.Str)
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
				var val string
				if p.cur.Type == kwON {
					val = "ON"
					p.advance()
				} else if p.cur.Type == kwOFF {
					val = "OFF"
					p.advance()
				} else if p.isIdentLike() {
					val = p.cur.Str
					p.advance()
				} else if p.cur.Type == tokSCONST {
					val = p.cur.Str
					p.advance()
				}
				opts = append(opts, &nodes.ServiceBrokerOption{Name: optName, Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
	}

	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseGetConversationGroupStmt parses GET CONVERSATION GROUP.
//
// BNF: mssql/parser/bnf/get-conversation-group-transact-sql.bnf
//
//	[ WAITFOR ( ]
//	    GET CONVERSATION GROUP @conversation_group_id
//	        FROM <queue>
//	[ ) ] [ , TIMEOUT timeout ]
func (p *Parser) parseGetConversationGroupStmt() (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	p.advance() // consume GET

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "GET",
		ObjectType: "CONVERSATION GROUP",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// CONVERSATION GROUP
	if p.cur.Type == kwCONVERSATION {
		p.advance()
	}
	if p.cur.Type == kwGROUP {
		p.advance()
	}

	// @conversation_group_id
	if p.cur.Type == tokVARIABLE {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// FROM queue
	if _, ok := p.match(kwFROM); ok {
		qLoc := p.pos()
		ref , _ := p.parseTableRef()
		if ref != nil {
			var opts []nodes.Node
			opts = append(opts, &nodes.ServiceBrokerOption{Name: "QUEUE", Value: ref.Object, Loc: nodes.Loc{Start: qLoc, End: p.prevEnd()}})
			stmt.Options = &nodes.List{Items: opts}
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseServiceBrokerOptions consumes a generic WITH options clause.
func (p *Parser) parseServiceBrokerOptions() *nodes.List {
	// Skip any options
	if p.cur.Type != kwWITH && p.cur.Type != kwON &&
		p.cur.Type != kwAUTHORIZATION {
		return nil
	}

	var opts []nodes.Node
	if p.cur.Type == kwAUTHORIZATION {
		optLoc := p.pos()
		p.advance()
		if p.isIdentLike() {
			opts = append(opts, &nodes.ServiceBrokerOption{Name: "AUTHORIZATION", Value: p.cur.Str, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
			p.advance()
		}
	}

	if p.cur.Type == kwON {
		p.advance()
		if p.cur.Type == kwQUEUE {
			qLoc := p.pos()
			p.advance()
			ref , _ := p.parseTableRef()
			if ref != nil {
				opts = append(opts, &nodes.ServiceBrokerOption{Name: "QUEUE", Value: ref.Object, Loc: nodes.Loc{Start: qLoc, End: p.prevEnd()}})
			}
		}
	}

	if p.cur.Type == kwWITH {
		p.advance()
		if p.cur.Type == '(' {
			p.advance()
			// Structured parsing of parenthesized key=value options
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				if p.isIdentLike() || p.cur.Type == kwUSER {
					optLoc := p.pos()
					opt := strings.ToUpper(p.cur.Str)
					p.advance()
					if p.cur.Type == '=' {
						p.advance()
						var val string
						if p.cur.Type == tokSCONST {
							val = p.cur.Str
							p.advance()
						} else if p.cur.Type == tokICONST {
							val = p.cur.Str
							p.advance()
						} else if p.cur.Type == kwON {
							val = "ON"
							p.advance()
						} else if p.cur.Type == kwOFF {
							val = "OFF"
							p.advance()
						} else if p.isIdentLike() {
							val = p.cur.Str
							p.advance()
						}
						opts = append(opts, &nodes.ServiceBrokerOption{Name: opt, Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
					} else {
						opts = append(opts, &nodes.ServiceBrokerOption{Name: opt, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
					}
				} else {
					// Skip unexpected tokens
					break
				}
				if _, ok := p.match(','); !ok {
					if p.cur.Type != ')' {
						break
					}
				}
			}
			p.match(')')
		} else {
			for p.isIdentLike() {
				optLoc := p.pos()
				opt := strings.ToUpper(p.cur.Str)
				p.advance()
				val := ""
				if p.cur.Type == '=' {
					p.advance()
					if p.isIdentLike() || p.cur.Type == kwON || p.cur.Type == kwOFF {
						val = strings.ToUpper(p.cur.Str)
						p.advance()
					}
				}
				opts = append(opts, &nodes.ServiceBrokerOption{Name: opt, Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
				if _, ok := p.match(','); !ok {
					break
				}
			}
		}
	}

	// Contract list for CREATE SERVICE
	if p.cur.Type == '(' {
		p.advance()
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			if p.isIdentLike() || p.cur.Type == kwDEFAULT {
				optLoc := p.pos()
				opts = append(opts, &nodes.ServiceBrokerOption{Name: strings.ToUpper(p.cur.Str), Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
				p.advance()
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
		p.match(')')
	}

	if len(opts) == 0 {
		return nil
	}
	return &nodes.List{Items: opts}
}

// parseAlterQueueStmt parses ALTER QUEUE.
//
// BNF: mssql/parser/bnf/alter-queue-transact-sql.bnf
//
//	ALTER QUEUE <object>
//	    queue_settings
//	    | queue_action
//
//	<object> ::=
//	{ database_name.schema_name.queue_name | schema_name.queue_name | queue_name }
//
//	<queue_settings> ::=
//	WITH
//	   [ STATUS = { ON | OFF } [ , ] ]
//	   [ RETENTION = { ON | OFF } [ , ] ]
//	   [ ACTIVATION (
//	       { [ STATUS = { ON | OFF } [ , ] ]
//	         [ PROCEDURE_NAME = <procedure> [ , ] ]
//	         [ MAX_QUEUE_READERS = max_readers [ , ] ]
//	         [ EXECUTE AS { SELF | 'user_name'  | OWNER } ]
//	       |  DROP }
//	          ) [ , ]]
//	         [ POISON_MESSAGE_HANDLING (
//	          STATUS = { ON | OFF } )
//	         ]
//
//	<queue_action> ::=
//	   REBUILD [ WITH <query_rebuild_options> ]
//	   | REORGANIZE [ WITH (LOB_COMPACTION = { ON | OFF } ) ]
//	   | MOVE TO { file_group | "default" }
//
//	<queue_rebuild_options> ::=
//	{
//	   ( MAXDOP = max_degree_of_parallelism )
//	}
func (p *Parser) parseAlterQueueStmt() (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	// QUEUE keyword already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "ALTER",
		ObjectType: "QUEUE",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// Parse possibly schema-qualified queue name
	ref , _ := p.parseTableRef()
	if ref != nil {
		if ref.Schema != "" {
			stmt.Name = ref.Schema + "." + ref.Object
		} else {
			stmt.Name = ref.Object
		}
	}

	var opts []nodes.Node

	// queue_action: REBUILD | REORGANIZE | MOVE TO
	if p.cur.Type == kwREBUILD {
		optLoc := p.pos()
		p.advance()
		val := ""
		// [ WITH ( MAXDOP = n ) ]
		if p.cur.Type == kwWITH {
			p.advance()
			if p.cur.Type == '(' {
				p.advance()
				if p.cur.Type == kwMAXDOP {
					p.advance()
					if p.cur.Type == '=' {
						p.advance()
					}
					if p.cur.Type == tokICONST {
						val = "MAXDOP=" + p.cur.Str
						p.advance()
					}
				}
				p.match(')')
			}
		}
		opts = append(opts, &nodes.ServiceBrokerOption{Name: "REBUILD", Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
	} else if p.cur.Type == kwREORGANIZE {
		optLoc := p.pos()
		p.advance()
		val := ""
		// [ WITH ( LOB_COMPACTION = { ON | OFF } ) ]
		if p.cur.Type == kwWITH {
			p.advance()
			if p.cur.Type == '(' {
				p.advance()
				if p.cur.Type == kwLOB_COMPACTION {
					p.advance()
					if p.cur.Type == '=' {
						p.advance()
					}
					if p.cur.Type == kwON || p.cur.Type == kwOFF {
						val = "LOB_COMPACTION=" + strings.ToUpper(p.cur.Str)
						p.advance()
					}
				}
				p.match(')')
			}
		}
		opts = append(opts, &nodes.ServiceBrokerOption{Name: "REORGANIZE", Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
	} else if p.cur.Type == kwMOVE {
		optLoc := p.pos()
		p.advance()
		// TO { file_group | "default" }
		if p.cur.Type == kwTO {
			p.advance()
		}
		val := ""
		if p.isIdentLike() || p.cur.Type == kwDEFAULT || p.cur.Type == tokSCONST {
			val = p.cur.Str
			p.advance()
		}
		opts = append(opts, &nodes.ServiceBrokerOption{Name: "MOVE TO", Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
	} else {
		// queue_settings: WITH clause
		opts = p.parseQueueWithClause(opts)
	}

	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterServiceStmt parses ALTER SERVICE.
//
// BNF: mssql/parser/bnf/alter-service-transact-sql.bnf
//
//	ALTER SERVICE service_name
//	   [ ON QUEUE [ schema_name . ]queue_name ]
//	   [ ( < opt_arg > [ , ...n ] ) ]
//
//	<opt_arg> ::=
//	   ADD CONTRACT contract_name | DROP CONTRACT contract_name
func (p *Parser) parseAlterServiceStmt() (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	// SERVICE keyword already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "ALTER",
		ObjectType: "SERVICE",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	var opts []nodes.Node

	// [ ON QUEUE [ schema_name . ] queue_name ]
	if p.cur.Type == kwON {
		p.advance()
		if p.cur.Type == kwQUEUE {
			qLoc := p.pos()
			p.advance()
			ref , _ := p.parseTableRef()
			if ref != nil {
				opts = append(opts, &nodes.ServiceBrokerOption{Name: "QUEUE", Value: ref.Object, Loc: nodes.Loc{Start: qLoc, End: p.prevEnd()}})
			}
		}
	}

	// [ ( ADD CONTRACT name | DROP CONTRACT name [,...n] ) ]
	if p.cur.Type == '(' {
		p.advance()
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			optLoc := p.pos()
			if p.cur.Type == kwADD {
				p.advance()
				if p.cur.Type == kwCONTRACT {
					p.advance()
				}
				if p.isAnyKeywordIdent() || p.cur.Type == tokSCONST {
					opts = append(opts, &nodes.ServiceBrokerOption{Name: "ADD CONTRACT", Value: p.cur.Str, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
					p.advance()
				}
			} else if p.cur.Type == kwDROP {
				p.advance()
				if p.cur.Type == kwCONTRACT {
					p.advance()
				}
				if p.isAnyKeywordIdent() || p.cur.Type == tokSCONST {
					opts = append(opts, &nodes.ServiceBrokerOption{Name: "DROP CONTRACT", Value: p.cur.Str, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
					p.advance()
				}
			} else {
				break
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
		p.match(')')
	}

	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterRouteStmt parses ALTER ROUTE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-route-transact-sql
//
//	ALTER ROUTE route_name
//	    WITH
//	    [ SERVICE_NAME = 'service_name' , ]
//	    [ BROKER_INSTANCE = 'broker_instance_identifier' , ]
//	    [ LIFETIME = route_lifetime , ]
//	    ADDRESS = 'next_hop_address'
//	    [ , MIRROR_ADDRESS = 'next_hop_mirror_address' ]
func (p *Parser) parseAlterRouteStmt() (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	// ROUTE keyword already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "ALTER",
		ObjectType: "ROUTE",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// WITH clause (same as CREATE ROUTE)
	if p.cur.Type == kwWITH {
		p.advance()
		var opts []nodes.Node
		for {
			if !p.isIdentLike() && p.cur.Type != tokSCONST {
				break
			}
			optLoc := p.pos()
			optName := strings.ToUpper(p.cur.Str)
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
				var val string
				if p.cur.Type == tokSCONST {
					val = p.cur.Str
					p.advance()
				} else if p.cur.Type == tokICONST {
					val = p.cur.Str
					p.advance()
				} else if p.isIdentLike() {
					val = p.cur.Str
					p.advance()
				}
				opts = append(opts, &nodes.ServiceBrokerOption{Name: optName, Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
			} else {
				opts = append(opts, &nodes.ServiceBrokerOption{Name: optName, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
		if len(opts) > 0 {
			stmt.Options = &nodes.List{Items: opts}
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterRemoteServiceBindingStmt parses ALTER REMOTE SERVICE BINDING.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-remote-service-binding-transact-sql
//
//	ALTER REMOTE SERVICE BINDING binding_name
//	    WITH USER = user_name [ , ANONYMOUS = { ON | OFF } ]
func (p *Parser) parseAlterRemoteServiceBindingStmt() (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	// REMOTE SERVICE BINDING keywords already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "ALTER",
		ObjectType: "REMOTE SERVICE BINDING",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// WITH USER = user_name [, ANONYMOUS = { ON | OFF }]
	var opts []nodes.Node
	if p.cur.Type == kwWITH {
		p.advance()
		for {
			if !p.isIdentLike() && p.cur.Type != kwUSER {
				break
			}
			optLoc := p.pos()
			optName := strings.ToUpper(p.cur.Str)
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
				var val string
				if p.cur.Type == kwON {
					val = "ON"
					p.advance()
				} else if p.cur.Type == kwOFF {
					val = "OFF"
					p.advance()
				} else if p.isIdentLike() {
					val = p.cur.Str
					p.advance()
				}
				opts = append(opts, &nodes.ServiceBrokerOption{Name: optName, Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
	}
	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterMessageTypeStmt parses ALTER MESSAGE TYPE.
//
// BNF: mssql/parser/bnf/alter-message-type-transact-sql.bnf
//
//	ALTER MESSAGE TYPE message_type_name
//	    VALIDATION =
//	    {  NONE
//	     | EMPTY
//	     | WELL_FORMED_XML
//	     | VALID_XML WITH SCHEMA COLLECTION schema_collection_name }
func (p *Parser) parseAlterMessageTypeStmt() (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	// MESSAGE TYPE keywords already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "ALTER",
		ObjectType: "MESSAGE TYPE",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// Optional: VALIDATION = { NONE | EMPTY | WELL_FORMED_XML | VALID_XML WITH SCHEMA COLLECTION name }
	if p.cur.Type == kwVALIDATION {
		valLoc := p.pos()
		p.advance()
		if p.cur.Type == '=' {
			p.advance()
		}
		var opts []nodes.Node
		// Consume the validation value (NONE, EMPTY, WELL_FORMED_XML, or VALID_XML)
		if p.isValidOption(validationValues) {
			valType := strings.ToUpper(p.cur.Str)
			p.advance()
			// Check for VALID_XML WITH SCHEMA COLLECTION name
			if valType == "VALID_XML" && p.cur.Type == kwWITH {
				p.advance() // consume WITH
				if p.cur.Type == kwSCHEMA {
					p.advance() // consume SCHEMA
				}
				if p.cur.Type == kwCOLLECTION {
					p.advance() // consume COLLECTION
				}
				// schema collection name (possibly dot-qualified)
				ref , _ := p.parseTableRef()
				schemaName := ""
				if ref != nil {
					if ref.Schema != "" {
						schemaName = ref.Schema + "." + ref.Object
					} else {
						schemaName = ref.Object
					}
				}
				opts = append(opts, &nodes.ServiceBrokerOption{Name: "VALIDATION", Value: "VALID_XML WITH SCHEMA COLLECTION " + schemaName, Loc: nodes.Loc{Start: valLoc, End: p.prevEnd()}})
			} else {
				opts = append(opts, &nodes.ServiceBrokerOption{Name: "VALIDATION", Value: valType, Loc: nodes.Loc{Start: valLoc, End: p.prevEnd()}})
			}
		} else {
			return nil, p.syntaxErrorAtCur()
		}
		if len(opts) > 0 {
			stmt.Options = &nodes.List{Items: opts}
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterContractStmt parses ALTER CONTRACT.
//
// ALTER CONTRACT is limited in T-SQL. The message types and directions in a contract
// cannot be changed. To change AUTHORIZATION, use ALTER AUTHORIZATION.
// We support ADD/DROP MESSAGE TYPE as extensions that some tools may emit.
//
//	ALTER CONTRACT contract_name
//	    [ ADD MESSAGE TYPE message_type_name SENT BY { INITIATOR | TARGET | ANY } ]
//	    [ DROP MESSAGE TYPE message_type_name ]
func (p *Parser) parseAlterContractStmt() (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	// CONTRACT keyword already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "ALTER",
		ObjectType: "CONTRACT",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	var opts []nodes.Node

	// ADD MESSAGE TYPE name SENT BY { INITIATOR | TARGET | ANY }
	if p.cur.Type == kwADD {
		addLoc := p.pos()
		p.advance() // consume ADD
		if p.cur.Type == kwMESSAGE {
			p.advance() // consume MESSAGE
			if p.cur.Type == kwTYPE {
				p.advance() // consume TYPE
			}
		}
		var msgTypeName string
		if p.isIdentLike() || p.cur.Type == tokSCONST {
			msgTypeName = p.cur.Str
			p.advance()
		}
		sentBy := ""
		if p.cur.Type == kwSENT {
			p.advance() // consume SENT
			if p.cur.Type == kwBY {
				p.advance() // consume BY
			}
			if p.isAnyKeywordIdent() || p.cur.Type == kwANY {
				sentBy = strings.ToUpper(p.cur.Str)
				p.advance()
			}
		}
		if msgTypeName != "" {
			val := msgTypeName
			if sentBy != "" {
				val += " SENT BY " + sentBy
			}
			opts = append(opts, &nodes.ServiceBrokerOption{Name: "ADD", Value: val, Loc: nodes.Loc{Start: addLoc, End: p.prevEnd()}})
		}
	}

	// DROP MESSAGE TYPE name
	if p.cur.Type == kwDROP {
		dropLoc := p.pos()
		p.advance() // consume DROP
		if p.cur.Type == kwMESSAGE {
			p.advance() // consume MESSAGE
			if p.cur.Type == kwTYPE {
				p.advance() // consume TYPE
			}
		}
		if p.isAnyKeywordIdent() || p.cur.Type == tokSCONST {
			opts = append(opts, &nodes.ServiceBrokerOption{Name: "DROP", Value: p.cur.Str, Loc: nodes.Loc{Start: dropLoc, End: p.prevEnd()}})
			p.advance()
		}
	}

	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseDropServiceBrokerStmt parses DROP for Service Broker objects.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-message-type-transact-sql
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-contract-transact-sql
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-queue-transact-sql
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-service-transact-sql
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-route-transact-sql
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-remote-service-binding-transact-sql
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-broker-priority-transact-sql
//
//	DROP { MESSAGE TYPE | CONTRACT | QUEUE | SERVICE | ROUTE | REMOTE SERVICE BINDING | BROKER PRIORITY } name
func (p *Parser) parseDropServiceBrokerStmt(objectType string) (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "DROP",
		ObjectType: objectType,
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// Parse possibly qualified name
	ref , _ := p.parseTableRef()
	if ref != nil {
		if ref.Schema != "" {
			stmt.Name = ref.Schema + "." + ref.Object
		} else {
			stmt.Name = ref.Object
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCreateBrokerPriorityStmt parses CREATE BROKER PRIORITY.
//
// BNF: mssql/parser/bnf/create-broker-priority-transact-sql.bnf
//
//	CREATE BROKER PRIORITY ConversationPriorityName
//	    FOR CONVERSATION
//	    [ SET (
//	        [ CONTRACT_NAME = { ContractName | ANY } ]
//	        [ [ , ] LOCAL_SERVICE_NAME = { LocalServiceName | ANY } ]
//	        [ [ , ] REMOTE_SERVICE_NAME = { 'RemoteServiceName' | ANY } ]
//	        [ [ , ] PRIORITY_LEVEL = { PriorityValue | DEFAULT } ]
//	    ) ]
func (p *Parser) parseCreateBrokerPriorityStmt() (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	// BROKER PRIORITY keywords already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "CREATE",
		ObjectType: "BROKER PRIORITY",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// FOR CONVERSATION
	if p.cur.Type == kwFOR {
		p.advance()
		if p.cur.Type == kwCONVERSATION {
			p.advance()
		}
	}

	stmt.Options = p.parseBrokerPrioritySetOptions()

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterBrokerPriorityStmt parses ALTER BROKER PRIORITY.
//
// BNF: mssql/parser/bnf/alter-broker-priority-transact-sql.bnf
//
//	ALTER BROKER PRIORITY ConversationPriorityName
//	    FOR CONVERSATION
//	    SET (
//	        [ CONTRACT_NAME = { ContractName | ANY } ]
//	        [ [ , ] LOCAL_SERVICE_NAME = { LocalServiceName | ANY } ]
//	        [ [ , ] REMOTE_SERVICE_NAME = { 'RemoteServiceName' | ANY } ]
//	        [ [ , ] PRIORITY_LEVEL = { PriorityValue | DEFAULT } ]
//	    )
func (p *Parser) parseAlterBrokerPriorityStmt() (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	// BROKER PRIORITY keywords already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "ALTER",
		ObjectType: "BROKER PRIORITY",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// FOR CONVERSATION
	if p.cur.Type == kwFOR {
		p.advance()
		if p.cur.Type == kwCONVERSATION {
			p.advance()
		}
	}

	stmt.Options = p.parseBrokerPrioritySetOptions()

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseBrokerPrioritySetOptions parses SET ( key=value, ... ) for BROKER PRIORITY statements.
func (p *Parser) parseBrokerPrioritySetOptions() *nodes.List {
	if p.cur.Type != kwSET {
		return nil
	}
	p.advance()
	if p.cur.Type != '(' {
		return nil
	}
	p.advance()
	var opts []nodes.Node
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		if p.isAnyKeywordIdent() {
			optLoc := p.pos()
			optName := strings.ToUpper(p.cur.Str)
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
				var val string
				if p.cur.Type == tokSCONST {
					val = p.cur.Str
					p.advance()
				} else if p.isAnyKeywordIdent() || p.cur.Type == kwDEFAULT || p.cur.Type == kwANY {
					val = strings.ToUpper(p.cur.Str)
					p.advance()
				} else if p.cur.Type == tokICONST {
					val = p.cur.Str
					p.advance()
				}
				opts = append(opts, &nodes.ServiceBrokerOption{Name: optName, Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
			}
		} else {
			p.advance()
		}
		p.match(',')
	}
	p.match(')')
	if len(opts) > 0 {
		return &nodes.List{Items: opts}
	}
	return nil
}

// parseMoveConversationStmt parses MOVE CONVERSATION.
//
// BNF: mssql/parser/bnf/move-conversation-transact-sql.bnf
//
//	MOVE CONVERSATION conversation_handle
//	    TO conversation_group_id
func (p *Parser) parseMoveConversationStmt() (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	p.advance() // consume MOVE

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "MOVE",
		ObjectType: "CONVERSATION",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// CONVERSATION
	if p.cur.Type == kwCONVERSATION {
		p.advance()
	}

	// conversation_handle
	if p.cur.Type == tokVARIABLE || p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// TO conversation_group_id
	if p.cur.Type == kwTO {
		toLoc := p.pos()
		p.advance()
		if p.cur.Type == tokVARIABLE || p.isAnyKeywordIdent() || p.cur.Type == tokSCONST {
			var opts []nodes.Node
			opts = append(opts, &nodes.ServiceBrokerOption{Name: "TO", Value: p.cur.Str, Loc: nodes.Loc{Start: toLoc, End: p.prevEnd()}})
			stmt.Options = &nodes.List{Items: opts}
			p.advance()
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseBeginConversationTimerStmt parses a BEGIN CONVERSATION TIMER statement.
//
// BNF: mssql/parser/bnf/begin-conversation-timer-transact-sql.bnf
//
//	BEGIN CONVERSATION TIMER ( conversation_handle )
//	   TIMEOUT = timeout
//	[ ; ]
func (p *Parser) parseBeginConversationTimerStmt() (*nodes.ServiceBrokerStmt, error) {
	// BEGIN and CONVERSATION TIMER already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "BEGIN",
		ObjectType: "CONVERSATION TIMER",
	}

	var opts []nodes.Node

	// ( conversation_handle )
	if p.cur.Type == '(' {
		hLoc := p.pos()
		p.advance()
		handle, _ := p.parseExpr()
		if handle != nil {
			opts = append(opts, &nodes.ServiceBrokerOption{Name: "HANDLE", Value: nodes.NodeToString(handle), Loc: nodes.Loc{Start: hLoc, End: p.prevEnd()}})
		}
		p.match(')')
	}

	// TIMEOUT = timeout
	if p.matchIdentCI("TIMEOUT") {
		tLoc := p.pos()
		p.match('=')
		timeout, _ := p.parseExpr()
		if timeout != nil {
			opts = append(opts, &nodes.ServiceBrokerOption{Name: "TIMEOUT", Value: nodes.NodeToString(timeout), Loc: nodes.Loc{Start: tLoc, End: p.prevEnd()}})
		}
	}

	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}
