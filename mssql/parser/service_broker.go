// Package parser - service_broker.go implements T-SQL Service Broker statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseCreateMessageTypeStmt parses CREATE MESSAGE TYPE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-message-type-transact-sql
//
//	CREATE MESSAGE TYPE message_type_name
//	    [ AUTHORIZATION owner_name ]
//	    [ VALIDATION = { NONE | EMPTY | WELL_FORMED_XML | VALID_XML WITH SCHEMA COLLECTION schema_collection_name } ]
func (p *Parser) parseCreateMessageTypeStmt() *nodes.ServiceBrokerStmt {
	loc := p.pos()
	// TYPE keyword already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "CREATE",
		ObjectType: "MESSAGE TYPE",
		Loc:        nodes.Loc{Start: loc},
	}

	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// Optional: AUTHORIZATION owner
	if p.cur.Type == kwAUTHORIZATION {
		p.advance()
		if p.isIdentLike() {
			p.advance()
		}
	}

	// Optional: VALIDATION = { NONE | EMPTY | WELL_FORMED_XML | VALID_XML WITH SCHEMA COLLECTION ... }
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "VALIDATION") {
		p.advance()
		if p.cur.Type == '=' {
			p.advance()
		}
		var opts []nodes.Node
		// consume the validation value
		for p.isIdentLike() || p.cur.Type == kwWITH {
			opt := strings.ToUpper(p.cur.Str)
			p.advance()
			if opt == "WITH" {
				// SCHEMA COLLECTION name
				for p.isIdentLike() {
					p.advance()
				}
				break
			}
			opts = append(opts, &nodes.String{Str: opt})
		}
		if len(opts) > 0 {
			stmt.Options = &nodes.List{Items: opts}
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseCreateContractStmt parses CREATE CONTRACT.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-contract-transact-sql
//
//	CREATE CONTRACT contract_name
//	    [ AUTHORIZATION owner_name ]
//	    ( {  { message_type_name | [ DEFAULT ] }
//	            SENT BY { INITIATOR | TARGET | ANY }
//	        } [ ,...n ] )
func (p *Parser) parseCreateContractStmt() *nodes.ServiceBrokerStmt {
	loc := p.pos()
	// CONTRACT keyword already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "CREATE",
		ObjectType: "CONTRACT",
		Loc:        nodes.Loc{Start: loc},
	}

	if p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// skip authorization and message type list
	if p.cur.Type == kwAUTHORIZATION {
		p.advance()
		if p.isIdentLike() {
			p.advance()
		}
	}
	if p.cur.Type == '(' {
		p.advance()
		depth := 1
		for depth > 0 && p.cur.Type != tokEOF {
			if p.cur.Type == '(' {
				depth++
			} else if p.cur.Type == ')' {
				depth--
			}
			if depth > 0 {
				p.advance()
			}
		}
		p.match(')')
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseCreateQueueStmt parses CREATE QUEUE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-queue-transact-sql
//
//	CREATE QUEUE <object>
//	    [ WITH
//	       [ STATUS = { ON | OFF } [ , ] ]
//	       [ RETENTION = { ON | OFF } [ , ] ]
//	       [ ACTIVATION (
//	           [ STATUS = { ON | OFF } , ]
//	           PROCEDURE_NAME = <procedure> ,
//	           MAX_QUEUE_READERS = max_readers ,
//	           EXECUTE AS { SELF | 'user_name' | OWNER }
//	           ) [ , ] ]
//	       [ POISON_MESSAGE_HANDLING ( STATUS = { ON | OFF } ) ]
//	    ]
//	    [ ON { filegroup | [DEFAULT] } ]
func (p *Parser) parseCreateQueueStmt() *nodes.ServiceBrokerStmt {
	loc := p.pos()
	// QUEUE keyword already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "CREATE",
		ObjectType: "QUEUE",
		Loc:        nodes.Loc{Start: loc},
	}

	if p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Options = p.parseServiceBrokerOptions()
	stmt.Loc.End = p.pos()
	return stmt
}

// parseCreateServiceStmt parses CREATE SERVICE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-service-transact-sql
//
//	CREATE SERVICE service_name
//	    [ AUTHORIZATION owner_name ]
//	    ON QUEUE [ schema . ] queue_name
//	    [ ( contract_name | [DEFAULT] [ ,...n ] ) ]
func (p *Parser) parseCreateServiceStmt() *nodes.ServiceBrokerStmt {
	loc := p.pos()
	// SERVICE keyword already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "CREATE",
		ObjectType: "SERVICE",
		Loc:        nodes.Loc{Start: loc},
	}

	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Options = p.parseServiceBrokerOptions()
	stmt.Loc.End = p.pos()
	return stmt
}

// parseSendStmt parses SEND ON CONVERSATION.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/send-transact-sql
//
//	SEND ON CONVERSATION conversation_handle
//	    [ MESSAGE TYPE message_type_name ]
//	    [ ( message_body_expression ) ]
func (p *Parser) parseSendStmt() *nodes.ServiceBrokerStmt {
	loc := p.pos()
	p.advance() // consume SEND

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "SEND",
		ObjectType: "CONVERSATION",
		Loc:        nodes.Loc{Start: loc},
	}

	// ON CONVERSATION handle
	if p.cur.Type == kwON {
		p.advance()
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "CONVERSATION") {
			p.advance()
		}
		// conversation handle (variable or expression)
		if p.cur.Type == tokVARIABLE || p.isIdentLike() {
			stmt.Name = p.cur.Str
			p.advance()
		}
	}

	// MESSAGE TYPE
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "MESSAGE") {
		p.advance()
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "TYPE") {
			p.advance()
		}
		if p.isIdentLike() {
			p.advance()
		}
	}

	// ( message_body )
	if p.cur.Type == '(' {
		p.advance()
		p.parseExpr()
		p.match(')')
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseReceiveStmt parses RECEIVE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/receive-transact-sql
//
//	[ WAITFOR ( ]
//	RECEIVE [ TOP ( n ) ]
//	    <column_specifier> [ ,...n ]
//	    FROM <queue>
//	    [ INTO table_variable ]
//	    [ WHERE { conversation_handle = @conversation_handle | conversation_group_id = @conversation_group_id } ]
//	[ ) ] [ , TIMEOUT timeout ]
func (p *Parser) parseReceiveStmt() *nodes.ServiceBrokerStmt {
	loc := p.pos()
	p.advance() // consume RECEIVE

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "RECEIVE",
		ObjectType: "QUEUE",
		Loc:        nodes.Loc{Start: loc},
	}

	// TOP (n)
	if p.cur.Type == kwTOP {
		p.advance()
		if p.cur.Type == '(' {
			p.advance()
			p.parseExpr()
			p.match(')')
		}
	}

	// column list until FROM
	for p.cur.Type != kwFROM && p.cur.Type != tokEOF && p.cur.Type != ';' {
		p.advance()
	}

	// FROM queue
	if _, ok := p.match(kwFROM); ok {
		ref := p.parseTableRef()
		if ref != nil {
			stmt.Name = ref.Object
		}
	}

	// WHERE clause
	if p.cur.Type == kwWHERE {
		p.advance()
		p.parseExpr()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseBeginConversationStmt parses BEGIN CONVERSATION TIMER or BEGIN DIALOG CONVERSATION.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/begin-dialog-conversation-transact-sql
//
//	BEGIN DIALOG [ CONVERSATION ] @dialog_handle
//	    FROM SERVICE initiator_service_name
//	    TO SERVICE 'target_service_name'
//	        [ , { 'service_broker_guid' | 'CURRENT DATABASE' }]
//	    [ ON CONTRACT contract_name ]
//	    [ WITH
//	        [ { RELATED_CONVERSATION = related_conversation_handle
//	          | RELATED_CONVERSATION_GROUP = related_conversation_group_id } , ]
//	        [ LIFETIME = dialog_lifetime , ]
//	        [ ENCRYPTION = { ON | OFF } ]
//	    ]
func (p *Parser) parseBeginConversationStmt() *nodes.ServiceBrokerStmt {
	loc := p.pos()
	// already consumed BEGIN (or DIALOG)

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "BEGIN",
		ObjectType: "CONVERSATION",
		Loc:        nodes.Loc{Start: loc},
	}

	// consume rest of line until end
	for p.cur.Type != ';' && p.cur.Type != tokEOF && p.cur.Type != kwGO {
		if p.isIdentLike() || p.cur.Type == tokVARIABLE || p.cur.Type == tokSCONST {
			if stmt.Name == "" && p.cur.Type == tokVARIABLE {
				stmt.Name = p.cur.Str
			}
		}
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseEndConversationStmt parses END CONVERSATION.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/end-conversation-transact-sql
//
//	END CONVERSATION conversation_handle
//	    [ WITH ERROR = failure_code DESCRIPTION = failure_text ]
//	    [ WITH CLEANUP ]
func (p *Parser) parseEndConversationStmt() *nodes.ServiceBrokerStmt {
	loc := p.pos()
	p.advance() // consume END
	// consume CONVERSATION
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "CONVERSATION") {
		p.advance()
	}

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "END",
		ObjectType: "CONVERSATION",
		Loc:        nodes.Loc{Start: loc},
	}

	// conversation handle
	if p.cur.Type == tokVARIABLE || p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// WITH ERROR/CLEANUP
	if p.cur.Type == kwWITH {
		p.advance()
		var opts []nodes.Node
		for p.cur.Type != ';' && p.cur.Type != tokEOF && p.cur.Type != kwGO {
			if p.isIdentLike() || p.cur.Type == kwNOT {
				opts = append(opts, &nodes.String{Str: strings.ToUpper(p.cur.Str)})
			}
			p.advance()
		}
		if len(opts) > 0 {
			stmt.Options = &nodes.List{Items: opts}
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
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
		p.advance()
		if p.isIdentLike() {
			opts = append(opts, &nodes.String{Str: "AUTHORIZATION=" + p.cur.Str})
			p.advance()
		}
	}

	if p.cur.Type == kwON {
		p.advance()
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "QUEUE") {
			p.advance()
			ref := p.parseTableRef()
			if ref != nil {
				opts = append(opts, &nodes.String{Str: "QUEUE=" + ref.Object})
			}
		}
	}

	if p.cur.Type == kwWITH {
		p.advance()
		if p.cur.Type == '(' {
			p.advance()
			depth := 1
			for depth > 0 && p.cur.Type != tokEOF {
				if p.cur.Type == '(' {
					depth++
				} else if p.cur.Type == ')' {
					depth--
				}
				if depth > 0 {
					p.advance()
				}
			}
			p.match(')')
		} else {
			for p.isIdentLike() {
				opt := strings.ToUpper(p.cur.Str)
				p.advance()
				if p.cur.Type == '=' {
					p.advance()
					if p.isIdentLike() || p.cur.Type == kwON || p.cur.Type == kwOFF {
						opt += "=" + strings.ToUpper(p.cur.Str)
						p.advance()
					}
				}
				opts = append(opts, &nodes.String{Str: opt})
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
				opts = append(opts, &nodes.String{Str: strings.ToUpper(p.cur.Str)})
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
