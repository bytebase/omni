// Package parser - service_broker.go implements T-SQL Service Broker statement parsing.
package parser

import (
	"fmt"
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

	// Optional: AUTHORIZATION owner
	if p.cur.Type == kwAUTHORIZATION {
		p.advance()
		if p.isIdentLike() {
			p.advance()
		}
	}

	// Parse message type definitions: ( message_type_name SENT BY { INITIATOR | TARGET | ANY } [,...n] )
	if p.cur.Type == '(' {
		p.advance()
		var opts []nodes.Node
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
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
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SENT") {
				p.advance()
				if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "BY") {
					p.advance()
				}
				if p.isIdentLike() || p.cur.Type == kwANY {
					sentBy = strings.ToUpper(p.cur.Str)
					p.advance()
				}
			}
			if msgTypeName != "" {
				entry := msgTypeName + " SENT BY " + sentBy
				opts = append(opts, &nodes.String{Str: entry})
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

	// Parse possibly schema-qualified queue name
	ref := p.parseTableRef()
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
		p.advance()
		if p.isIdentLike() || p.cur.Type == kwDEFAULT {
			opts = append(opts, &nodes.String{Str: "ON=" + strings.ToUpper(p.cur.Str)})
			p.advance()
		}
	}

	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}
	stmt.Loc.End = p.pos()
	return stmt
}

// parseQueueWithClause parses the WITH clause options for CREATE/ALTER QUEUE.
func (p *Parser) parseQueueWithClause(opts []nodes.Node) []nodes.Node {
	if p.cur.Type != kwWITH {
		return opts
	}
	p.advance()

	for {
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "STATUS") {
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			if p.cur.Type == kwON || p.cur.Type == kwOFF {
				opts = append(opts, &nodes.String{Str: "STATUS=" + strings.ToUpper(p.cur.Str)})
				p.advance()
			}
		} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "RETENTION") {
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			if p.cur.Type == kwON || p.cur.Type == kwOFF {
				opts = append(opts, &nodes.String{Str: "RETENTION=" + strings.ToUpper(p.cur.Str)})
				p.advance()
			}
		} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "ACTIVATION") {
			p.advance()
			opts = p.parseQueueActivationClause(opts)
		} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "POISON_MESSAGE_HANDLING") {
			p.advance()
			if p.cur.Type == '(' {
				p.advance()
				if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "STATUS") {
					p.advance()
					if p.cur.Type == '=' {
						p.advance()
					}
					if p.cur.Type == kwON || p.cur.Type == kwOFF {
						opts = append(opts, &nodes.String{Str: "POISON_MESSAGE_HANDLING=" + strings.ToUpper(p.cur.Str)})
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
	var activationOpts []string
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "STATUS") {
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			if p.cur.Type == kwON || p.cur.Type == kwOFF {
				activationOpts = append(activationOpts, "STATUS="+strings.ToUpper(p.cur.Str))
				p.advance()
			}
		} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "PROCEDURE_NAME") {
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			// Parse possibly schema-qualified procedure name
			ref := p.parseTableRef()
			if ref != nil {
				name := ref.Object
				if ref.Schema != "" {
					name = ref.Schema + "." + name
				}
				if ref.Database != "" {
					name = ref.Database + "." + name
				}
				activationOpts = append(activationOpts, "PROCEDURE_NAME="+name)
			}
		} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "MAX_QUEUE_READERS") {
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			if p.cur.Type == tokICONST {
				activationOpts = append(activationOpts, "MAX_QUEUE_READERS="+p.cur.Str)
				p.advance()
			}
		} else if p.cur.Type == kwEXECUTE || (p.isIdentLike() && matchesKeywordCI(p.cur.Str, "EXECUTE")) {
			p.advance()
			if p.cur.Type == kwAS {
				p.advance()
			}
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SELF") {
				activationOpts = append(activationOpts, "EXECUTE AS=SELF")
				p.advance()
			} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "OWNER") {
				activationOpts = append(activationOpts, "EXECUTE AS=OWNER")
				p.advance()
			} else if p.cur.Type == tokSCONST {
				activationOpts = append(activationOpts, "EXECUTE AS="+p.cur.Str)
				p.advance()
			}
		} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "DROP") {
			activationOpts = append(activationOpts, "DROP")
			p.advance()
		} else {
			break
		}
		p.match(',')
	}
	p.match(')')

	for _, ao := range activationOpts {
		opts = append(opts, &nodes.String{Str: "ACTIVATION:" + ao})
	}
	return opts
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
//
//	<column_specifier> ::=
//	{    *
//	  |  { column_name | [ ] expression } [ [ AS ] column_alias ]
//	}     [ ,...n ]
//
//	<queue> ::=
//	{ database_name.schema_name.queue_name | schema_name.queue_name | queue_name }
func (p *Parser) parseReceiveStmt() *nodes.ServiceBrokerStmt {
	loc := p.pos()
	p.advance() // consume RECEIVE

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "RECEIVE",
		ObjectType: "QUEUE",
		Loc:        nodes.Loc{Start: loc},
	}

	var opts []nodes.Node

	// TOP (n)
	if p.cur.Type == kwTOP {
		p.advance()
		if p.cur.Type == '(' {
			p.advance()
			topExpr := p.parseExpr()
			if topExpr != nil {
				opts = append(opts, &nodes.String{Str: "TOP=" + exprToString(topExpr)})
			}
			p.match(')')
		}
	}

	// column list: * or column_name [AS alias] [,...n]
	if p.cur.Type == '*' {
		opts = append(opts, &nodes.String{Str: "COLUMNS=*"})
		p.advance()
	} else {
		var cols []string
		for p.cur.Type != kwFROM && p.cur.Type != tokEOF && p.cur.Type != ';' {
			// Parse column expression (simple: column_name or expression)
			colStr := ""
			if p.isIdentLike() || p.cur.Type == tokVARIABLE {
				colStr = p.cur.Str
				p.advance()
			} else {
				// complex expression, skip until comma or AS or FROM
				expr := p.parseExpr()
				if expr != nil {
					colStr = exprToString(expr)
				}
			}
			// optional alias: [AS] alias
			if p.cur.Type == kwAS {
				p.advance()
				if p.isIdentLike() {
					colStr += " AS " + p.cur.Str
					p.advance()
				}
			} else if p.isIdentLike() && p.cur.Type != kwFROM && p.cur.Type != kwINTO && p.cur.Type != kwWHERE {
				// alias without AS - be careful not to consume FROM/INTO/WHERE
				colStr += " AS " + p.cur.Str
				p.advance()
			}
			if colStr != "" {
				cols = append(cols, colStr)
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
		if len(cols) > 0 {
			opts = append(opts, &nodes.String{Str: "COLUMNS=" + strings.Join(cols, ",")})
		}
	}

	// FROM queue
	if _, ok := p.match(kwFROM); ok {
		ref := p.parseTableRef()
		if ref != nil {
			queueName := ref.Object
			if ref.Schema != "" {
				queueName = ref.Schema + "." + queueName
			}
			if ref.Database != "" {
				queueName = ref.Database + "." + queueName
			}
			stmt.Name = ref.Object
			opts = append(opts, &nodes.String{Str: "FROM=" + queueName})
		}
	}

	// INTO table_variable
	if p.cur.Type == kwINTO {
		p.advance()
		if p.cur.Type == tokVARIABLE {
			opts = append(opts, &nodes.String{Str: "INTO=" + p.cur.Str})
			p.advance()
		} else if p.isIdentLike() {
			opts = append(opts, &nodes.String{Str: "INTO=" + p.cur.Str})
			p.advance()
		}
	}

	// WHERE clause
	if p.cur.Type == kwWHERE {
		p.advance()
		whereExpr := p.parseExpr()
		if whereExpr != nil {
			opts = append(opts, &nodes.String{Str: "WHERE"})
		}
	}

	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}
	stmt.Loc.End = p.pos()
	return stmt
}

// exprToString returns a simple string representation of an expression node.
func exprToString(n nodes.Node) string {
	if n == nil {
		return ""
	}
	switch v := n.(type) {
	case *nodes.String:
		return v.Str
	case *nodes.Integer:
		return fmt.Sprintf("%d", v.Ival)
	case *nodes.ColumnRef:
		var parts []string
		if v.Server != "" {
			parts = append(parts, v.Server)
		}
		if v.Database != "" {
			parts = append(parts, v.Database)
		}
		if v.Schema != "" {
			parts = append(parts, v.Schema)
		}
		if v.Table != "" {
			parts = append(parts, v.Table)
		}
		if v.Column != "" {
			parts = append(parts, v.Column)
		}
		return strings.Join(parts, ".")
	case *nodes.VariableRef:
		return v.Name
	}
	return "expr"
}

// parseBeginConversationStmt parses BEGIN DIALOG CONVERSATION.
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
//	          | RELATED_CONVERSATION_GROUP = related_conversation_group_id } ]
//	        [ [ , ] LIFETIME = dialog_lifetime ]
//	        [ [ , ] ENCRYPTION = { ON | OFF } ]
//	    ]
func (p *Parser) parseBeginConversationStmt() *nodes.ServiceBrokerStmt {
	loc := p.pos()
	// already consumed BEGIN, and DIALOG has been consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "BEGIN",
		ObjectType: "CONVERSATION",
		Loc:        nodes.Loc{Start: loc},
	}

	// optional CONVERSATION keyword
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "CONVERSATION") {
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
		p.advance()
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SERVICE") {
			p.advance()
		}
		if p.isIdentLike() || p.cur.Type == tokSCONST {
			opts = append(opts, &nodes.String{Str: "FROM SERVICE=" + p.cur.Str})
			p.advance()
		}
	}

	// TO SERVICE 'target_service_name' [, { 'service_broker_guid' | 'CURRENT DATABASE' }]
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "TO") {
		p.advance()
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SERVICE") {
			p.advance()
		}
		if p.cur.Type == tokSCONST || p.isIdentLike() {
			opts = append(opts, &nodes.String{Str: "TO SERVICE=" + p.cur.Str})
			p.advance()
		}
		// optional broker guid
		if _, ok := p.match(','); ok {
			if p.cur.Type == tokSCONST || p.isIdentLike() {
				opts = append(opts, &nodes.String{Str: "BROKER_INSTANCE=" + p.cur.Str})
				p.advance()
			}
		}
	}

	// ON CONTRACT contract_name
	if p.cur.Type == kwON {
		p.advance()
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "CONTRACT") {
			p.advance()
		}
		if p.isIdentLike() || p.cur.Type == tokSCONST {
			opts = append(opts, &nodes.String{Str: "ON CONTRACT=" + p.cur.Str})
			p.advance()
		}
	}

	// WITH options
	if p.cur.Type == kwWITH {
		p.advance()
		for {
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "RELATED_CONVERSATION_GROUP") {
				p.advance()
				if p.cur.Type == '=' {
					p.advance()
				}
				if p.cur.Type == tokVARIABLE || p.isIdentLike() || p.cur.Type == tokSCONST {
					opts = append(opts, &nodes.String{Str: "RELATED_CONVERSATION_GROUP=" + p.cur.Str})
					p.advance()
				}
			} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "RELATED_CONVERSATION") {
				p.advance()
				if p.cur.Type == '=' {
					p.advance()
				}
				if p.cur.Type == tokVARIABLE || p.isIdentLike() || p.cur.Type == tokSCONST {
					opts = append(opts, &nodes.String{Str: "RELATED_CONVERSATION=" + p.cur.Str})
					p.advance()
				}
			} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "LIFETIME") {
				p.advance()
				if p.cur.Type == '=' {
					p.advance()
				}
				if p.cur.Type == tokICONST || p.cur.Type == tokVARIABLE || p.isIdentLike() {
					opts = append(opts, &nodes.String{Str: "LIFETIME=" + p.cur.Str})
					p.advance()
				}
			} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "ENCRYPTION") {
				p.advance()
				if p.cur.Type == '=' {
					p.advance()
				}
				if p.cur.Type == kwON || p.cur.Type == kwOFF {
					opts = append(opts, &nodes.String{Str: "ENCRYPTION=" + strings.ToUpper(p.cur.Str)})
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
	stmt.Loc.End = p.pos()
	return stmt
}

// parseEndConversationStmt parses END CONVERSATION.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/end-conversation-transact-sql
//
//	END CONVERSATION conversation_handle
//	    [   [ WITH ERROR = failure_code DESCRIPTION = 'failure_text' ]
//	      | [ WITH CLEANUP ]
//	    ]
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

	// WITH ERROR = failure_code DESCRIPTION = 'failure_text'
	// WITH CLEANUP
	if p.cur.Type == kwWITH {
		p.advance()
		var opts []nodes.Node
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "ERROR") {
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
			opts = append(opts, &nodes.String{Str: "ERROR=" + errorCode})

			// DESCRIPTION = 'failure_text'
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "DESCRIPTION") {
				p.advance()
				if p.cur.Type == '=' {
					p.advance()
				}
				descVal := ""
				if p.cur.Type == tokSCONST || p.cur.Type == tokVARIABLE || p.isIdentLike() {
					descVal = p.cur.Str
					p.advance()
				}
				opts = append(opts, &nodes.String{Str: "DESCRIPTION=" + descVal})
			}
		} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "CLEANUP") {
			opts = append(opts, &nodes.String{Str: "CLEANUP"})
			p.advance()
		}
		if len(opts) > 0 {
			stmt.Options = &nodes.List{Items: opts}
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
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
func (p *Parser) parseCreateRouteStmt() *nodes.ServiceBrokerStmt {
	loc := p.pos()
	// ROUTE keyword already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "CREATE",
		ObjectType: "ROUTE",
		Loc:        nodes.Loc{Start: loc},
	}

	// route_name
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

	// WITH clause (required for CREATE ROUTE)
	if p.cur.Type == kwWITH {
		p.advance()
		var opts []nodes.Node
		for {
			if !p.isIdentLike() && p.cur.Type != tokSCONST {
				break
			}
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
				opts = append(opts, &nodes.String{Str: optName + "=" + val})
			} else {
				opts = append(opts, &nodes.String{Str: optName})
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
		if len(opts) > 0 {
			stmt.Options = &nodes.List{Items: opts}
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseCreateRemoteServiceBindingStmt parses CREATE REMOTE SERVICE BINDING.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-remote-service-binding-transact-sql
//
//	CREATE REMOTE SERVICE BINDING binding_name
//	    [ AUTHORIZATION owner_name ]
//	    TO SERVICE 'service_name'
//	    WITH USER = user_name [ , ANONYMOUS = { ON | OFF } ]
func (p *Parser) parseCreateRemoteServiceBindingStmt() *nodes.ServiceBrokerStmt {
	loc := p.pos()
	// REMOTE SERVICE BINDING keywords already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "CREATE",
		ObjectType: "REMOTE SERVICE BINDING",
		Loc:        nodes.Loc{Start: loc},
	}

	// binding_name
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

	var opts []nodes.Node

	// TO SERVICE 'service_name'
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "TO") {
		p.advance()
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SERVICE") {
			p.advance()
		}
		if p.cur.Type == tokSCONST {
			opts = append(opts, &nodes.String{Str: "SERVICE=" + p.cur.Str})
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
			optName := strings.ToUpper(p.cur.Str)
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
				var val string
				if p.isIdentLike() {
					val = p.cur.Str
					p.advance()
				} else if p.cur.Type == kwON {
					val = "ON"
					p.advance()
				} else if p.cur.Type == kwOFF {
					val = "OFF"
					p.advance()
				} else if p.cur.Type == tokSCONST {
					val = p.cur.Str
					p.advance()
				}
				opts = append(opts, &nodes.String{Str: optName + "=" + val})
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
	}

	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseGetConversationGroupStmt parses GET CONVERSATION GROUP.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/get-conversation-group-transact-sql
//
//	[ WAITFOR ( ]
//	    GET CONVERSATION GROUP @conversation_group_id
//	        FROM <queue>
//	[ ) ] [ , TIMEOUT timeout ]
func (p *Parser) parseGetConversationGroupStmt() *nodes.ServiceBrokerStmt {
	loc := p.pos()
	p.advance() // consume GET

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "GET",
		ObjectType: "CONVERSATION GROUP",
		Loc:        nodes.Loc{Start: loc},
	}

	// CONVERSATION GROUP
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "CONVERSATION") {
		p.advance()
	}
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "GROUP") {
		p.advance()
	}

	// @conversation_group_id
	if p.cur.Type == tokVARIABLE {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// FROM queue
	if _, ok := p.match(kwFROM); ok {
		ref := p.parseTableRef()
		if ref != nil {
			var opts []nodes.Node
			opts = append(opts, &nodes.String{Str: "QUEUE=" + ref.Object})
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

// parseAlterQueueStmt parses ALTER QUEUE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-queue-transact-sql
//
//	ALTER QUEUE [ database_name . [ schema_name ] . | schema_name . ] queue_name
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
//	    [ REBUILD [ WITH ( <queue_rebuild_options> ) ] ]
//	    [ REORGANIZE [ WITH ( LOB_COMPACTION = { ON | OFF } ) ] ]
//	    [ MOVE TO { file_group | [DEFAULT] } ]
func (p *Parser) parseAlterQueueStmt() *nodes.ServiceBrokerStmt {
	loc := p.pos()
	// QUEUE keyword already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "ALTER",
		ObjectType: "QUEUE",
		Loc:        nodes.Loc{Start: loc},
	}

	// Parse possibly schema-qualified queue name
	ref := p.parseTableRef()
	if ref != nil {
		if ref.Schema != "" {
			stmt.Name = ref.Schema + "." + ref.Object
		} else {
			stmt.Name = ref.Object
		}
	}

	var opts []nodes.Node
	opts = p.parseQueueWithClause(opts)

	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseAlterServiceStmt parses ALTER SERVICE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-service-transact-sql
//
//	ALTER SERVICE service_name
//	    [ ON QUEUE [ schema_name . ] queue_name ]
//	    [ ( { ADD CONTRACT contract_name } | { DROP CONTRACT contract_name } [ ,...n ] ) ]
func (p *Parser) parseAlterServiceStmt() *nodes.ServiceBrokerStmt {
	loc := p.pos()
	// SERVICE keyword already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "ALTER",
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
func (p *Parser) parseAlterRouteStmt() *nodes.ServiceBrokerStmt {
	loc := p.pos()
	// ROUTE keyword already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "ALTER",
		ObjectType: "ROUTE",
		Loc:        nodes.Loc{Start: loc},
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
				opts = append(opts, &nodes.String{Str: optName + "=" + val})
			} else {
				opts = append(opts, &nodes.String{Str: optName})
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
		if len(opts) > 0 {
			stmt.Options = &nodes.List{Items: opts}
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseAlterRemoteServiceBindingStmt parses ALTER REMOTE SERVICE BINDING.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-remote-service-binding-transact-sql
//
//	ALTER REMOTE SERVICE BINDING binding_name
//	    WITH USER = user_name [ , ANONYMOUS = { ON | OFF } ]
func (p *Parser) parseAlterRemoteServiceBindingStmt() *nodes.ServiceBrokerStmt {
	loc := p.pos()
	// REMOTE SERVICE BINDING keywords already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "ALTER",
		ObjectType: "REMOTE SERVICE BINDING",
		Loc:        nodes.Loc{Start: loc},
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
			optName := strings.ToUpper(p.cur.Str)
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
				var val string
				if p.isIdentLike() {
					val = p.cur.Str
					p.advance()
				} else if p.cur.Type == kwON {
					val = "ON"
					p.advance()
				} else if p.cur.Type == kwOFF {
					val = "OFF"
					p.advance()
				}
				opts = append(opts, &nodes.String{Str: optName + "=" + val})
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
	}
	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.pos()
	return stmt
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
func (p *Parser) parseDropServiceBrokerStmt(objectType string) *nodes.ServiceBrokerStmt {
	loc := p.pos()

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "DROP",
		ObjectType: objectType,
		Loc:        nodes.Loc{Start: loc},
	}

	// Parse possibly qualified name
	ref := p.parseTableRef()
	if ref != nil {
		if ref.Schema != "" {
			stmt.Name = ref.Schema + "." + ref.Object
		} else {
			stmt.Name = ref.Object
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseCreateBrokerPriorityStmt parses CREATE BROKER PRIORITY.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-broker-priority-transact-sql
//
//	CREATE BROKER PRIORITY ConversationPriorityName
//	    FOR CONVERSATION
//	    [ SET (
//	        [ CONTRACT_NAME = { ContractName | ANY } ]
//	        [ [ , ] LOCAL_SERVICE_NAME = { LocalServiceName | ANY } ]
//	        [ [ , ] REMOTE_SERVICE_NAME = { 'RemoteServiceName' | ANY } ]
//	        [ [ , ] PRIORITY_LEVEL = { PriorityValue | DEFAULT } ]
//	    ) ]
func (p *Parser) parseCreateBrokerPriorityStmt() *nodes.ServiceBrokerStmt {
	loc := p.pos()
	// BROKER PRIORITY keywords already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "CREATE",
		ObjectType: "BROKER PRIORITY",
		Loc:        nodes.Loc{Start: loc},
	}

	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// FOR CONVERSATION
	if p.cur.Type == kwFOR {
		p.advance()
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "CONVERSATION") {
			p.advance()
		}
	}

	// SET ( ... )
	if p.cur.Type == kwSET {
		p.advance()
		if p.cur.Type == '(' {
			p.advance()
			var opts []nodes.Node
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				if p.isIdentLike() {
					optName := strings.ToUpper(p.cur.Str)
					p.advance()
					if p.cur.Type == '=' {
						p.advance()
						var val string
						if p.cur.Type == tokSCONST {
							val = p.cur.Str
							p.advance()
						} else if p.isIdentLike() || p.cur.Type == kwDEFAULT || p.cur.Type == kwANY {
							val = strings.ToUpper(p.cur.Str)
							p.advance()
						} else if p.cur.Type == tokICONST {
							val = p.cur.Str
							p.advance()
						}
						opts = append(opts, &nodes.String{Str: optName + "=" + val})
					}
				} else {
					p.advance()
				}
				p.match(',')
			}
			p.match(')')
			if len(opts) > 0 {
				stmt.Options = &nodes.List{Items: opts}
			}
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseAlterBrokerPriorityStmt parses ALTER BROKER PRIORITY.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-broker-priority-transact-sql
//
//	ALTER BROKER PRIORITY ConversationPriorityName
//	    FOR CONVERSATION
//	    SET (
//	        [ CONTRACT_NAME = { ContractName | ANY } ]
//	        [ [ , ] LOCAL_SERVICE_NAME = { LocalServiceName | ANY } ]
//	        [ [ , ] REMOTE_SERVICE_NAME = { 'RemoteServiceName' | ANY } ]
//	        [ [ , ] PRIORITY_LEVEL = { PriorityValue | DEFAULT } ]
//	    )
func (p *Parser) parseAlterBrokerPriorityStmt() *nodes.ServiceBrokerStmt {
	loc := p.pos()
	// BROKER PRIORITY keywords already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "ALTER",
		ObjectType: "BROKER PRIORITY",
		Loc:        nodes.Loc{Start: loc},
	}

	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// FOR CONVERSATION
	if p.cur.Type == kwFOR {
		p.advance()
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "CONVERSATION") {
			p.advance()
		}
	}

	// SET ( ... )
	if p.cur.Type == kwSET {
		p.advance()
		if p.cur.Type == '(' {
			p.advance()
			var opts []nodes.Node
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				if p.isIdentLike() {
					optName := strings.ToUpper(p.cur.Str)
					p.advance()
					if p.cur.Type == '=' {
						p.advance()
						var val string
						if p.cur.Type == tokSCONST {
							val = p.cur.Str
							p.advance()
						} else if p.isIdentLike() || p.cur.Type == kwDEFAULT || p.cur.Type == kwANY {
							val = strings.ToUpper(p.cur.Str)
							p.advance()
						} else if p.cur.Type == tokICONST {
							val = p.cur.Str
							p.advance()
						}
						opts = append(opts, &nodes.String{Str: optName + "=" + val})
					}
				} else {
					p.advance()
				}
				p.match(',')
			}
			p.match(')')
			if len(opts) > 0 {
				stmt.Options = &nodes.List{Items: opts}
			}
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseMoveConversationStmt parses MOVE CONVERSATION.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/move-conversation-transact-sql
//
//	MOVE CONVERSATION conversation_handle
//	    TO conversation_group_id
func (p *Parser) parseMoveConversationStmt() *nodes.ServiceBrokerStmt {
	loc := p.pos()
	p.advance() // consume MOVE

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "MOVE",
		ObjectType: "CONVERSATION",
		Loc:        nodes.Loc{Start: loc},
	}

	// CONVERSATION
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "CONVERSATION") {
		p.advance()
	}

	// conversation_handle
	if p.cur.Type == tokVARIABLE || p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// TO conversation_group_id
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "TO") {
		p.advance()
		if p.cur.Type == tokVARIABLE || p.isIdentLike() || p.cur.Type == tokSCONST {
			var opts []nodes.Node
			opts = append(opts, &nodes.String{Str: "TO=" + p.cur.Str})
			stmt.Options = &nodes.List{Items: opts}
			p.advance()
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseBeginConversationTimerStmt parses a BEGIN CONVERSATION TIMER statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/begin-conversation-timer-transact-sql
//
//	BEGIN CONVERSATION TIMER ( conversation_handle )
//	   TIMEOUT = timeout
//	[ ; ]
func (p *Parser) parseBeginConversationTimerStmt() *nodes.ServiceBrokerStmt {
	// BEGIN and CONVERSATION TIMER already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "BEGIN",
		ObjectType: "CONVERSATION TIMER",
	}

	var opts []nodes.Node

	// ( conversation_handle )
	if p.cur.Type == '(' {
		p.advance()
		handle := p.parseExpr()
		if handle != nil {
			opts = append(opts, &nodes.String{Str: "HANDLE=" + nodes.NodeToString(handle)})
		}
		p.match(')')
	}

	// TIMEOUT = timeout
	if p.matchIdentCI("TIMEOUT") {
		p.match('=')
		timeout := p.parseExpr()
		if timeout != nil {
			opts = append(opts, &nodes.String{Str: "TIMEOUT=" + nodes.NodeToString(timeout)})
		}
	}

	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.pos()
	return stmt
}
