// Package parser - utility.go implements T-SQL utility statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseUseStmt parses a USE database statement.
//
//	USE database
//
// useStatement (TSql160.g:13635): the database is required.
func (p *Parser) parseUseStmt() (*nodes.UseStmt, error) {
	loc := p.pos()
	p.advance() // consume USE

	// Completion: after USE → database_ref
	if p.collectMode() {
		p.addRuleCandidate("database_ref")
		return nil, errCollecting
	}

	stmt := &nodes.UseStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	name, ok := p.parseIdentifier()
	if !ok {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Database = name

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parsePrintStmt parses a PRINT statement.
//
// BNF: mssql/parser/bnf/print-transact-sql.bnf
//
//	PRINT msg_str | @local_variable | string_expr
func (p *Parser) parsePrintStmt() (*nodes.PrintStmt, error) {
	loc := p.pos()
	p.advance() // consume PRINT

	// Completion: after PRINT → expression context
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		p.addRuleCandidate("func_name")
		return nil, errCollecting
	}

	stmt := &nodes.PrintStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, p.newParseError(p.cur.Loc, "expected expression after PRINT")
	}
	stmt.Expr = expr

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseRaiseErrorStmt parses a RAISERROR statement.
//
// BNF: mssql/parser/bnf/raiserror-transact-sql.bnf
//
//	RAISERROR ( { msg_id | msg_str | @local_variable }
//	    { , severity , state }
//	    [ , argument [ , ...n ] ] )
//	    [ WITH option [ , ...n ] ]
//
// raiseErrorStatement (TSql160.g:19486): the message is a signed integer, a
// string or a variable, severity and state signed integers or variables, an
// argument a constant or a negative number, all required where written; the
// options are LOG, NOWAIT and SETERROR.
func (p *Parser) parseRaiseErrorStmt() (*nodes.RaiseErrorStmt, error) {
	loc := p.pos()
	p.advance() // consume RAISERROR

	// Completion: after RAISERROR( → expression context
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		p.addRuleCandidate("func_name")
		return nil, errCollecting
	}

	stmt := &nodes.RaiseErrorStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	// Completion: after RAISERROR( → expression context
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		p.addRuleCandidate("func_name")
		return nil, errCollecting
	}
	var err error
	stmt.Message, err = p.parseProcValue(procValueKinds{signedInt: true, str: true, variable: true})
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(','); err != nil {
		return nil, err
	}
	stmt.Severity, err = p.parseProcValue(procValueKinds{signedInt: true, variable: true})
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(','); err != nil {
		return nil, err
	}
	stmt.State, err = p.parseProcValue(procValueKinds{signedInt: true, variable: true})
	if err != nil {
		return nil, err
	}
	// Optional args: possibleNegativeConstant (TSql160.g:25159)
	var args []nodes.Node
	for {
		if _, ok := p.match(','); !ok {
			break
		}
		if p.isIdentLike() {
			return nil, p.syntaxErrorAtCur()
		}
		arg, err := p.parseProcConstantOrIdentifier()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
	if len(args) > 0 {
		stmt.Args = &nodes.List{Items: args}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	// WITH options (LOG, NOWAIT, SETERROR)
	if p.cur.Type == kwWITH {
		p.advance()
		var opts []nodes.Node
		for {
			if !p.procIsWord("LOG") && !p.procIsWord("NOWAIT") && !p.procIsWord("SETERROR") {
				return nil, p.syntaxErrorAtCur()
			}
			opts = append(opts, &nodes.String{Str: strings.ToUpper(p.cur.Str)})
			p.advance()
			if _, ok := p.match(','); !ok {
				break
			}
		}
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseThrowStmt parses a THROW statement.
//
// BNF: mssql/parser/bnf/throw-transact-sql.bnf
//
//	THROW [ { error_number | @local_variable }
//	    , { message | @local_variable }
//	    , { state | @local_variable } ]
//	[ ; ]
//
// throwStatement (TSql160.g:19537): all three arguments or none, and none
// only inside a CATCH block (Msg 10704).
func (p *Parser) parseThrowStmt() (*nodes.ThrowStmt, error) {
	loc := p.pos()
	p.advance() // consume THROW

	// Completion: after THROW → expression context
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		p.addRuleCandidate("func_name")
		return nil, errCollecting
	}

	stmt := &nodes.ThrowStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// THROW without arguments = rethrow
	if p.cur.Type == ';' || p.cur.Type == tokEOF || p.cur.Type == kwEND ||
		p.isStatementStart() {
		if p.procCatchDepth == 0 {
			return nil, p.newParseError(loc, "to rethrow an error, a THROW statement must be used inside a CATCH block")
		}
		stmt.Loc.End = p.prevEnd()
		return stmt, nil
	}

	var err error
	stmt.ErrorNumber, err = p.parseProcValue(procValueKinds{signedInt: true, variable: true})
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(','); err != nil {
		return nil, err
	}
	stmt.Message, err = p.parseProcValue(procValueKinds{str: true, variable: true})
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(','); err != nil {
		return nil, err
	}
	stmt.State, err = p.parseProcValue(procValueKinds{integer: true, variable: true})
	if err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCheckpointStmt parses a CHECKPOINT statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/checkpoint-transact-sql
//
//	CHECKPOINT [ checkpoint_duration ]
func (p *Parser) parseCheckpointStmt() (*nodes.CheckpointStmt, error) {
	loc := p.pos()
	p.advance() // consume CHECKPOINT

	stmt := &nodes.CheckpointStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Optional duration
	if p.cur.Type == tokICONST || p.cur.Type == tokVARIABLE {
		stmt.Duration, _ = p.parseExpr()
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseReconfigureStmt parses a RECONFIGURE statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/reconfigure-transact-sql
//
//	RECONFIGURE [ WITH OVERRIDE ]
func (p *Parser) parseReconfigureStmt() (*nodes.ReconfigureStmt, error) {
	loc := p.pos()
	p.advance() // consume RECONFIGURE

	stmt := &nodes.ReconfigureStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// WITH OVERRIDE
	if p.cur.Type == kwWITH {
		next := p.peekNext()
		if next.Type == kwOVERRIDE {
			p.advance() // WITH
			p.advance() // OVERRIDE
			stmt.WithOverride = true
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseShutdownStmt parses a SHUTDOWN statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/shutdown-transact-sql
//
//	SHUTDOWN [ WITH NOWAIT ]
//
// shutdownStatement (TSql160.g:13833).
func (p *Parser) parseShutdownStmt() (*nodes.ShutdownStmt, error) {
	loc := p.pos()
	p.advance() // consume SHUTDOWN

	stmt := &nodes.ShutdownStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// WITH NOWAIT
	if p.cur.Type == kwWITH {
		p.advance()
		if p.cur.Type != kwNOWAIT {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		stmt.WithNoWait = true
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseKillStmt parses a KILL statement.
//
// BNF: mssql/parser/bnf/kill-transact-sql.bnf
//
//	KILL { session_id [ WITH STATUSONLY ] | UOW [ WITH STATUSONLY | COMMIT | ROLLBACK ] }
//	KILL STATS JOB job_id
//	KILL QUERY NOTIFICATION SUBSCRIPTION { ALL | subscription_id }
//
// killStatement (TSql160.g:13723): a signed integer or a string, which SQL
// Server converts to a unit of work's GUID (Msg 8169 when it is not one).
func (p *Parser) parseKillStmt() (nodes.StmtNode, error) {
	loc := p.pos()
	p.advance() // consume KILL

	// Check for KILL QUERY NOTIFICATION SUBSCRIPTION
	if p.cur.Type == kwQUERY {
		next := p.peekNext()
		if next.Type == kwNOTIFICATION {
			return p.parseKillQueryNotificationStmt(loc)
		}
	}

	// Check for KILL STATS JOB job_id
	if p.cur.Type == kwSTATS {
		next := p.peekNext()
		if next.Type == kwJOB {
			return p.parseKillStatsJobStmt(loc)
		}
	}

	stmt := &nodes.KillStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// session_id or UOW
	if procIsString(p.cur) && !procIsGUID(p.cur.Str) {
		return nil, p.newParseError(p.cur.Loc, "conversion failed when converting from a character string to uniqueidentifier")
	}
	var err error
	stmt.SessionID, err = p.parseProcValue(procValueKinds{signedInt: true, str: true})
	if err != nil {
		return nil, err
	}

	// WITH { STATUSONLY | COMMIT | ROLLBACK }
	if p.cur.Type == kwWITH {
		p.advance()
		switch p.cur.Type {
		case kwSTATUSONLY:
			stmt.StatusOnly = true
		case kwCOMMIT:
			stmt.WithAction = "COMMIT"
		case kwROLLBACK:
			stmt.WithAction = "ROLLBACK"
		default:
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseKillStatsJobStmt parses a KILL STATS JOB statement.
//
// BNF: mssql/parser/bnf/kill-stats-job-transact-sql.bnf
//
//	KILL STATS JOB job_id
//
// killStatsJobStatement (TSql160.g:13779): the job is a signed integer.
func (p *Parser) parseKillStatsJobStmt(loc int) (*nodes.KillStatsJobStmt, error) {
	p.advance() // consume STATS
	p.advance() // consume JOB

	stmt := &nodes.KillStatsJobStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	var err error
	stmt.JobID, err = p.parseProcValue(procValueKinds{signedInt: true})
	if err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseKillQueryNotificationStmt parses a KILL QUERY NOTIFICATION SUBSCRIPTION statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/kill-query-notification-subscription-transact-sql
//
//	KILL QUERY NOTIFICATION SUBSCRIPTION { ALL | subscription_id }
//
// killQueryNotificationSubscriptionStatement (TSql160.g:13748): every word
// is required, and the subscription is ALL or an integer.
func (p *Parser) parseKillQueryNotificationStmt(loc int) (*nodes.KillQueryNotificationStmt, error) {
	p.advance() // consume QUERY
	p.advance() // consume NOTIFICATION
	if _, err := p.expect(kwSUBSCRIPTION); err != nil {
		return nil, err
	}

	stmt := &nodes.KillQueryNotificationStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// { ALL | subscription_id }
	if p.cur.Type == kwALL {
		stmt.All = true
		p.advance()
	} else {
		var err error
		stmt.SubscriptionID, err = p.parseProcValue(procValueKinds{integer: true})
		if err != nil {
			return nil, err
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseReadtextStmt parses a READTEXT statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/readtext-transact-sql
//
//	READTEXT { table.column text_ptr offset size } [ HOLDLOCK ]
//
// readTextStatement (TSql160.g:14213): the text pointer is a binary
// constant or a variable, offset and size integers or variables.
func (p *Parser) parseReadtextStmt() (*nodes.ReadtextStmt, error) {
	loc := p.pos()
	p.advance() // consume READTEXT

	stmt := &nodes.ReadtextStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	var err error
	if stmt.Column, err = p.parseProcTextColumn(); err != nil {
		return nil, err
	}
	if stmt.TextPtr, err = p.parseProcValue(procValueKinds{binary: true, variable: true}); err != nil {
		return nil, err
	}
	if stmt.Offset, err = p.parseProcValue(procValueKinds{integer: true, variable: true}); err != nil {
		return nil, err
	}
	if stmt.Size, err = p.parseProcValue(procValueKinds{integer: true, variable: true}); err != nil {
		return nil, err
	}

	// HOLDLOCK
	if p.cur.Type == kwHOLDLOCK {
		p.advance()
		stmt.HoldLock = true
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// procWriteStringKinds are the values of writeString (TSql160.g:14361).
var procWriteStringKinds = procValueKinds{null: true, str: true, binary: true, variable: true}

// parseWritetextStmt parses a WRITETEXT statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/writetext-transact-sql
//
//	WRITETEXT { table.column text_ptr } [ WITH LOG ] { data }
//
// writeTextStatement (TSql160.g:14294).
func (p *Parser) parseWritetextStmt() (*nodes.WritetextStmt, error) {
	loc := p.pos()
	p.advance() // consume WRITETEXT

	stmt := &nodes.WritetextStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	var err error
	if stmt.Column, stmt.TextPtr, err = p.parseProcTextTarget(); err != nil {
		return nil, err
	}
	if stmt.WithLog, err = p.parseProcTextWithLog(); err != nil {
		return nil, err
	}
	if stmt.Data, err = p.parseProcValue(procWriteStringKinds); err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseUpdatetextStmt parses an UPDATETEXT statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/updatetext-transact-sql
//
//	UPDATETEXT { table_name.dest_column_name dest_text_ptr }
//	    { NULL | insert_offset }
//	    { NULL | delete_length }
//	    [ WITH LOG ]
//	    [ inserted_data | { table_name.src_column_name src_text_ptr } ]
//
// updateTextStatement (TSql160.g:14249): offset and length are signed
// integers, variables or NULL; the source is a column and a text pointer
// (a binary constant or a variable), or a NULL, string, binary or variable.
func (p *Parser) parseUpdatetextStmt() (*nodes.UpdatetextStmt, error) {
	loc := p.pos()
	p.advance() // consume UPDATETEXT

	stmt := &nodes.UpdatetextStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	var err error
	if stmt.DestColumn, stmt.DestTextPtr, err = p.parseProcTextTarget(); err != nil {
		return nil, err
	}
	offsetKinds := procValueKinds{signedInt: true, variable: true, null: true}
	if stmt.InsertOffset, err = p.parseProcValue(offsetKinds); err != nil {
		return nil, err
	}
	if stmt.DeleteLength, err = p.parseProcValue(offsetKinds); err != nil {
		return nil, err
	}
	if stmt.WithLog, err = p.parseProcTextWithLog(); err != nil {
		return nil, err
	}

	// inserted_data | table_name.src_column_name src_text_ptr
	switch {
	case p.isIdentLike():
		col, err := p.parseProcTextColumn()
		if err != nil {
			return nil, err
		}
		if _, err := p.parseProcValue(procValueKinds{binary: true, variable: true}); err != nil {
			return nil, err
		}
		stmt.InsertedData = col
	case p.cur.Type == kwNULL || procIsString(p.cur) || procIsBinaryConst(p.cur) || p.cur.Type == tokVARIABLE:
		if stmt.InsertedData, err = p.parseProcValue(procWriteStringKinds); err != nil {
			return nil, err
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseTruncateStmt parses a TRUNCATE TABLE statement.
//
// BNF: mssql/parser/bnf/truncate-table-transact-sql.bnf
//
//	TRUNCATE TABLE
//	    { database_name.schema_name.table_name | schema_name.table_name | table_name }
//	    [ WITH ( PARTITIONS ( { <partition_number_expression> | <range> }
//	    [ , ...n ] ) ) ]
//
//	<range> ::=
//	<partition_number_expression> TO <partition_number_expression>
func (p *Parser) parseTruncateStmt() (*nodes.TruncateStmt, error) {
	loc := p.pos()
	p.advance() // consume TRUNCATE

	// TABLE
	p.match(kwTABLE)

	// Completion: after TRUNCATE TABLE → table_ref
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		return nil, errCollecting
	}

	stmt := &nodes.TruncateStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	stmt.Table, _ = p.parseTableRef()

	// WITH ( PARTITIONS ( range [,...n] ) )
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		if p.cur.Type == '(' {
			p.advance() // consume outer '('
			if p.matchIdentCI("PARTITIONS") {
				if p.cur.Type == '(' {
					p.advance() // consume inner '('
					var parts []nodes.Node
					for p.cur.Type != ')' && p.cur.Type != tokEOF {
						expr, _ := p.parseExpr()
						// Check for TO (range)
						if p.cur.Type == kwTO {
							p.advance()   // consume TO
							p.parseExpr() // end of range (consumed but range stored as start expr)
						}
						parts = append(parts, expr)
						if _, ok := p.match(','); !ok {
							break
						}
					}
					p.match(')') // consume inner ')'
					if len(parts) > 0 {
						stmt.Partitions = &nodes.List{Items: parts}
					}
				}
			}
			p.match(')') // consume outer ')'
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCreateDefaultStmt parses a CREATE DEFAULT statement.
// Caller has consumed CREATE DEFAULT.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-default-transact-sql
//
//	CREATE DEFAULT [ schema_name . ] default_name
//	AS constant_expression [ ; ]
func (p *Parser) parseCreateDefaultStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()

	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "DEFAULT",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// [ schema_name . ] default_name
	if p.isIdentLike() || p.cur.Type == tokSCONST {
		name := p.cur.Str
		p.advance()
		if p.cur.Type == '.' {
			p.advance() // consume '.'
			if p.isIdentLike() || p.cur.Type == tokSCONST {
				name = name + "." + p.cur.Str
				p.advance()
			}
		}
		stmt.Name = name
	}

	// AS constant_expression
	var opts []nodes.Node
	if p.cur.Type == kwAS {
		p.advance() // consume AS
		// Parse the expression that follows AS
		expr, _ := p.parseExpr()
		if expr != nil {
			opts = append(opts, &nodes.String{Str: "AS"})
			opts = append(opts, expr)
		}
	}
	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCreateRuleStmt parses a CREATE RULE statement.
// Caller has consumed CREATE RULE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-rule-transact-sql
//
//	CREATE RULE [ schema_name . ] rule_name
//	AS condition_expression [ ; ]
func (p *Parser) parseCreateRuleStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()

	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "RULE",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// [ schema_name . ] rule_name
	if p.isIdentLike() || p.cur.Type == tokSCONST {
		name := p.cur.Str
		p.advance()
		if p.cur.Type == '.' {
			p.advance() // consume '.'
			if p.isIdentLike() || p.cur.Type == tokSCONST {
				name = name + "." + p.cur.Str
				p.advance()
			}
		}
		stmt.Name = name
	}

	// AS condition_expression
	var opts []nodes.Node
	if p.cur.Type == kwAS {
		p.advance() // consume AS
		expr, _ := p.parseExpr()
		if expr != nil {
			opts = append(opts, &nodes.String{Str: "AS"})
			opts = append(opts, expr)
		}
	}
	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseEnableDisableTriggerStmt parses ENABLE TRIGGER or DISABLE TRIGGER.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/enable-trigger-transact-sql
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/disable-trigger-transact-sql
//
//	{ ENABLE | DISABLE } TRIGGER { [ schema_name . ] trigger_name [ , ...n ] | ALL }
//	    ON { object_name | DATABASE | ALL SERVER }
func (p *Parser) parseEnableDisableTriggerStmt(enable bool) (*nodes.EnableDisableTriggerStmt, error) {
	loc := p.pos()
	p.advance() // consume ENABLE or DISABLE
	p.advance() // consume TRIGGER

	stmt := &nodes.EnableDisableTriggerStmt{
		Enable: enable,
		Loc:    nodes.Loc{Start: loc, End: -1},
	}

	// { trigger_name [ , ...n ] | ALL }
	if p.cur.Type == kwALL {
		stmt.TriggerAll = true
		p.advance()
	} else {
		var triggers []nodes.Node
		for {
			if p.isIdentLike() || p.cur.Type == tokSCONST {
				// Possibly schema-qualified: schema.trigger_name
				name := p.cur.Str
				p.advance()
				if p.cur.Type == '.' {
					p.advance()
					if p.isIdentLike() {
						name = name + "." + p.cur.Str
						p.advance()
					}
				}
				triggers = append(triggers, &nodes.String{Str: name})
			} else {
				break
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
		if len(triggers) > 0 {
			stmt.Triggers = &nodes.List{Items: triggers}
		}
	}

	// ON { object_name | DATABASE | ALL SERVER }
	if _, ok := p.match(kwON); ok {
		if p.cur.Type == kwDATABASE {
			stmt.OnDatabase = true
			p.advance()
		} else if p.cur.Type == kwALL {
			// ALL SERVER (triggerScope, TSql160.g:25412)
			p.advance()
			if !p.matchIdentCI("SERVER") {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.OnAllServer = true
		} else {
			obj, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			if obj == nil {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.OnObject = obj
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseSetuserStmt parses a SETUSER statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/setuser-transact-sql
//
//	SETUSER [ 'username' [ WITH NORESET ] ]
func (p *Parser) parseSetuserStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	p.advance() // consume SETUSER

	stmt := &nodes.SecurityStmt{
		Action:     "SETUSER",
		ObjectType: "USER",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// Optional 'username' or variable, then optional WITH NORESET
	// (setUserStatement, TSql160.g:13851).
	if procIsString(p.cur) || p.cur.Type == tokVARIABLE {
		stmt.Name = p.cur.Str
		p.advance()

		if p.cur.Type == kwWITH {
			p.advance() // consume WITH
			if err := p.procExpectWord("NORESET"); err != nil {
				return nil, err
			}
			stmt.Options = &nodes.List{Items: []nodes.Node{&nodes.String{Str: "NORESET"}}}
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCopyIntoOption parses a single COPY INTO WITH option.
// Options can be:
//   - KEY = 'string_value'
//   - KEY = numeric_value
//   - KEY = identifier (e.g., FILE_FORMAT = myformat)
//   - CREDENTIAL = ( ... ) -- parenthesized credential
//   - ERRORFILE_CREDENTIAL = ( ... ) -- parenthesized credential
func (p *Parser) parseCopyIntoOption() nodes.Node {
	if !p.isKeywordOrIdent() && p.cur.Type != kwFILE {
		return nil
	}

	name := strings.ToUpper(p.cur.Str)
	p.advance()

	// Check for = sign
	if p.cur.Type == '=' {
		p.advance() // consume '='

		// Special handling for CREDENTIAL and ERRORFILE_CREDENTIAL: value in parens
		if (name == "CREDENTIAL" || name == "ERRORFILE_CREDENTIAL") && p.cur.Type == '(' {
			p.advance() // consume '('
			var parts []string
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				parts = append(parts, p.cur.Str)
				p.advance()
			}
			p.match(')')
			return &nodes.String{Str: name + "=(" + strings.Join(parts, " ") + ")"}
		}

		var valStr string
		switch p.cur.Type {
		case tokSCONST, tokNSCONST:
			valStr = "'" + p.cur.Str + "'"
			p.advance()
		case tokICONST:
			valStr = p.cur.Str
			p.advance()
		case tokFCONST:
			valStr = p.cur.Str
			p.advance()
		default:
			if p.isIdentLike() {
				valStr = p.cur.Str
				p.advance()
			}
		}
		return &nodes.String{Str: name + "=" + valStr}
	}

	// Plain flag option (unlikely for COPY INTO but handle gracefully)
	return &nodes.String{Str: name}
}

// parseRenameStmt parses a RENAME statement (Azure Synapse / PDW).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/rename-transact-sql
//
//	RENAME OBJECT [::] [ [ database_name . [ schema_name ] . ] | [ schema_name . ] ] table_name TO new_table_name
//	RENAME DATABASE [::] database_name TO new_database_name
//	RENAME OBJECT [::] [ [ database_name . [ schema_name ] . ] | [ schema_name . ] ] table_name COLUMN column_name TO new_column_name
func (p *Parser) parseRenameStmt() (*nodes.RenameStmt, error) {
	loc := p.pos()
	p.advance() // consume RENAME

	stmt := &nodes.RenameStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// OBJECT or DATABASE
	if p.cur.Type == kwOBJECT {
		stmt.ObjectType = "OBJECT"
		p.advance()
	} else if p.cur.Type == kwDATABASE {
		stmt.ObjectType = "DATABASE"
		p.advance()
	}

	// Optional :: separator
	if p.cur.Type == tokCOLONCOLON {
		p.advance() // consume ::
	}

	// Object/table/database name
	stmt.Name, _ = p.parseTableRef()

	// Check for COLUMN rename variant
	if p.cur.Type == kwCOLUMN {
		p.advance() // consume COLUMN
		if p.isIdentLike() || p.cur.Type == '[' {
			stmt.ColumnName = p.cur.Str
			p.advance()
		}
	}

	// TO new_name
	if p.cur.Type == kwTO {
		p.advance() // consume TO
		if p.isIdentLike() || p.cur.Type == '[' {
			if stmt.ColumnName != "" {
				stmt.NewColumnName = p.cur.Str
			} else {
				stmt.NewName = p.cur.Str
			}
			p.advance()
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCreateTableCloneStmt parses CREATE TABLE ... AS CLONE OF (Fabric).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-table-as-clone-of-transact-sql
//
//	CREATE TABLE
//	    { database_name.schema_name.table_name | schema_name.table_name | table_name }
//	AS CLONE OF
//	    { database_name.schema_name.table_name | schema_name.table_name | table_name }
//	    [ AT { point_in_time } ]
func (p *Parser) parseCreateTableCloneStmt(name *nodes.TableRef) (*nodes.CreateTableCloneStmt, error) {
	loc := p.pos()
	// AS CLONE OF already partially consumed; caller consumed AS, we consume CLONE OF

	stmt := &nodes.CreateTableCloneStmt{
		Name: name,
		Loc:  nodes.Loc{Start: loc, End: -1},
	}

	// consume CLONE
	p.matchIdentCI("CLONE")
	// consume OF
	p.matchIdentCI("OF")

	// Source table name
	stmt.SourceName, _ = p.parseTableRef()

	// Optional AT 'point_in_time'
	if p.cur.Type == kwAT {
		p.advance() // consume AT
		if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
			stmt.AtTime = p.cur.Str
			p.advance()
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCreateRemoteTableAsSelectStmt parses a CREATE REMOTE TABLE AS SELECT (CRTAS) statement.
// Caller has consumed CREATE REMOTE TABLE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-remote-table-as-select-parallel-data-warehouse
//
//	CREATE REMOTE TABLE { database_name.schema_name.table_name | schema_name.table_name | table_name }
//	    AT ('<connection_string>')
//	    [ WITH ( BATCH_SIZE = batch_size ) ]
//	    AS <select_statement>
//
//	<connection_string> ::=
//	    Data Source = { IP_address | hostname } [, port ]; User ID = user_name; Password = strong_password;
//
//	<select_statement> ::=
//	    [ WITH <common_table_expression> [ ,...n ] ]
//	    SELECT <select_criteria>
func (p *Parser) parseCreateRemoteTableAsSelectStmt() (*nodes.CreateRemoteTableAsSelectStmt, error) {
	loc := p.pos()
	// REMOTE TABLE already consumed by caller

	stmt := &nodes.CreateRemoteTableAsSelectStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Table name: { database_name.schema_name.table_name | schema_name.table_name | table_name }
	stmt.Name, _ = p.parseTableRef()

	// AT ('<connection_string>')
	if p.cur.Type == kwAT {
		p.advance() // consume AT
		if p.cur.Type == '(' {
			p.advance() // consume '('
			if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
				stmt.ConnectionString = p.cur.Str
				p.advance()
			}
			p.match(')') // consume ')'
		}
	}

	// [ WITH ( BATCH_SIZE = batch_size ) ]
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		if p.cur.Type == '(' {
			p.advance() // consume '('
			var opts []nodes.Node
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				opt := p.parseCopyIntoOption() // reuse generic KEY=VALUE option parser
				if opt != nil {
					opts = append(opts, opt)
				}
				if _, ok := p.match(','); !ok {
					break
				}
			}
			p.match(')') // consume ')'
			if len(opts) > 0 {
				stmt.Options = &nodes.List{Items: opts}
			}
		}
	}

	// AS <select_statement>
	if p.cur.Type == kwAS {
		p.advance() // consume AS
		if p.cur.Type == kwSELECT || p.cur.Type == kwWITH {
			stmt.Query, _ = p.parseSelectStmtWithCTE()
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// ---------- Batch 174: Federation Statements ----------

// parseCreateFederationStmt parses a CREATE FEDERATION statement.
//
// BNF: mssql/parser/bnf/create-federation-transact-sql.bnf
//
//	CREATE FEDERATION federation_name ( distribution_name data_type RANGE )
func (p *Parser) parseCreateFederationStmt() (*nodes.CreateFederationStmt, error) {
	stmt := &nodes.CreateFederationStmt{
		Loc: nodes.Loc{Start: p.pos(), End: -1},
	}

	// federation_name
	if name, ok := p.parseIdentifier(); ok {
		stmt.Name = name
	}

	// ( distribution_name data_type RANGE )
	p.match('(')
	if name, ok := p.parseIdentifier(); ok {
		stmt.DistributionName = name
	}
	// data_type - parse as type reference string
	stmt.DataType, _ = p.parseDataType()
	// RANGE
	p.matchIdentCI("RANGE")
	p.match(')')

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterFederationStmt parses an ALTER FEDERATION statement.
//
// BNF: mssql/parser/bnf/alter-federation-transact-sql.bnf
//
//	ALTER FEDERATION federation_name
//	{
//	    SPLIT AT ( distribution_name = value )
//	  | DROP AT ( { LOW | HIGH } distribution_name = value )
//	}
func (p *Parser) parseAlterFederationStmt() (*nodes.AlterFederationStmt, error) {
	stmt := &nodes.AlterFederationStmt{
		Loc: nodes.Loc{Start: p.pos(), End: -1},
	}

	// federation_name
	if name, ok := p.parseIdentifier(); ok {
		stmt.Name = name
	}

	// SPLIT AT or DROP AT
	if p.cur.Type == kwSPLIT {
		p.advance() // consume SPLIT
		p.matchIdentCI("AT")
		stmt.Kind = "SPLIT"
	} else if p.cur.Type == kwDROP {
		p.advance() // consume DROP
		p.matchIdentCI("AT")
		p.match('(')
		// LOW or HIGH
		if p.cur.Type == kwLOW {
			stmt.Kind = "DROP LOW"
			p.advance()
		} else if p.cur.Type == kwHIGH {
			stmt.Kind = "DROP HIGH"
			p.advance()
		}
		// distribution_name = value
		if name, ok := p.parseIdentifier(); ok {
			stmt.DistributionName = name
		}
		p.match('=')
		stmt.Boundary, _ = p.parseExpr()
		p.match(')')
		stmt.Loc.End = p.prevEnd()
		return stmt, nil
	}

	// For SPLIT: ( distribution_name = value )
	p.match('(')
	if name, ok := p.parseIdentifier(); ok {
		stmt.DistributionName = name
	}
	p.match('=')
	stmt.Boundary, _ = p.parseExpr()
	p.match(')')

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseDropFederationStmt parses a DROP FEDERATION statement.
//
// BNF: mssql/parser/bnf/drop-federation-transact-sql.bnf
//
//	DROP FEDERATION federation_name
func (p *Parser) parseDropFederationStmt() (*nodes.DropFederationStmt, error) {
	stmt := &nodes.DropFederationStmt{
		Loc: nodes.Loc{Start: p.pos(), End: -1},
	}

	if name, ok := p.parseIdentifier(); ok {
		stmt.Name = name
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseUseFederationStmt parses a USE FEDERATION statement.
//
// BNF: mssql/parser/bnf/use-federation-transact-sql.bnf
//
//	USE FEDERATION
//	{
//	    ROOT WITH
//	  | federation_name ( distribution_name = value )
//	    WITH FILTERING = { ON | OFF } ,
//	}
//	RESET
func (p *Parser) parseUseFederationStmt() (*nodes.UseFederationStmt, error) {
	stmt := &nodes.UseFederationStmt{
		Loc: nodes.Loc{Start: p.pos(), End: -1},
	}

	// ROOT or federation_name
	if p.cur.Type == kwROOT {
		stmt.IsRoot = true
		p.advance() // consume ROOT
		p.match(kwWITH)
	} else {
		// federation_name
		if name, ok := p.parseIdentifier(); ok {
			stmt.FederationName = name
		}
		// ( distribution_name = value )
		p.match('(')
		if name, ok := p.parseIdentifier(); ok {
			stmt.DistributionName = name
		}
		p.match('=')
		stmt.Value, _ = p.parseExpr()
		p.match(')')
		// WITH FILTERING = { ON | OFF }
		p.match(kwWITH)
		p.matchIdentCI("FILTERING")
		p.match('=')
		if p.cur.Type == kwON {
			stmt.Filtering = true
			p.advance()
		} else if p.cur.Type == kwOFF {
			p.advance()
		}
		// comma separator
		p.match(',')
	}

	// RESET
	p.matchIdentCI("RESET")

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// ---------- Batch 175: INSERT BULK and LINENO ----------

// parseInsertBulkStmt parses an INSERT BULK statement.
//
// BNF: mssql/parser/bnf/insert-bulk-transact-sql.bnf
//
//	INSERT BULK schemaObjectThreePartName
//	  [ ( column_name data_type [ NULL | NOT NULL ] [ ,...n ] ) ]
//	  [ WITH ( option [ = value ] [ ,...n ] ) ]
func (p *Parser) parseInsertBulkStmt() (*nodes.InsertBulkStmt, error) {
	stmt := &nodes.InsertBulkStmt{
		Loc: nodes.Loc{Start: p.pos(), End: -1},
	}

	// table name (three-part name)
	stmt.Table, _ = p.parseTableRef()

	// optional column definition list
	if p.cur.Type == '(' {
		p.advance() // consume (
		var cols []nodes.Node
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			col := p.parseInsertBulkColumnDef()
			if col != nil {
				cols = append(cols, col)
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
		p.match(')') // consume )
		if len(cols) > 0 {
			stmt.Columns = &nodes.List{Items: cols}
		}
	}

	// optional WITH ( options )
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		if p.cur.Type == '(' {
			p.advance() // consume (
			var opts []nodes.Node
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				opt := p.parseBulkInsertOption()
				if opt != nil {
					opts = append(opts, opt)
				}
				if _, ok := p.match(','); !ok {
					break
				}
			}
			p.match(')') // consume )
			if len(opts) > 0 {
				stmt.Options = &nodes.List{Items: opts}
			}
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseInsertBulkColumnDef parses a column definition in INSERT BULK.
//
//	column_name data_type [ NULL | NOT NULL ]
func (p *Parser) parseInsertBulkColumnDef() *nodes.InsertBulkColumnDef {
	col := &nodes.InsertBulkColumnDef{
		Loc: nodes.Loc{Start: p.pos(), End: -1},
	}

	if name, ok := p.parseIdentifier(); ok {
		col.Name = name
	}

	col.DataType, _ = p.parseDataType()

	// Optional NULL / NOT NULL
	if p.cur.Type == kwNOT {
		p.advance() // consume NOT
		p.match(kwNULL)
		f := false
		col.Nullable = &f
	} else if p.cur.Type == kwNULL {
		p.advance() // consume NULL
		t := true
		col.Nullable = &t
	}

	col.Loc.End = p.prevEnd()
	return col
}

// parseLinenoStmt parses a LINENO statement.
//
// BNF: mssql/parser/bnf/lineno-transact-sql.bnf
//
//	LINENO integer
func (p *Parser) parseLinenoStmt() (*nodes.LinenoStmt, error) {
	loc := p.pos()
	p.advance() // consume LINENO

	stmt := &nodes.LinenoStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	if p.cur.Type == tokICONST {
		stmt.LineNo = int(p.cur.Ival)
		p.advance()
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// procIsGUID reports whether s has the form of a uniqueidentifier,
// 8-4-4-4-12 hex digits, optionally in braces.
func procIsGUID(s string) bool {
	if len(s) == 38 && s[0] == '{' && s[37] == '}' {
		s = s[1:37]
	}
	if len(s) != 36 {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if c != '-' {
				return false
			}
			continue
		}
		if !isHexDigit(c) {
			return false
		}
	}
	return true
}

// parseProcTextColumn parses the column of READTEXT, WRITETEXT and
// UPDATETEXT (column, TSql160.g:32779): a name of two to four parts, the
// last of which may be IDENTITYCOL or ROWGUIDCOL. SQL Server needs the table
// and the column (Msg 182).
func (p *Parser) parseProcTextColumn() (*nodes.ColumnRef, error) {
	loc := p.pos()
	var parts []string
	for {
		if len(parts) > 0 && (p.cur.Type == kwIDENTITYCOL || p.cur.Type == kwROWGUIDCOL) {
			parts = append(parts, p.cur.Str)
			p.advance()
			break
		}
		name, ok := p.parseIdentifier()
		if !ok {
			return nil, p.syntaxErrorAtCur()
		}
		parts = append(parts, name)
		if p.cur.Type != '.' {
			break
		}
		p.advance() // consume .
	}
	if len(parts) < 2 {
		return nil, p.newParseError(loc, "table and column names must be supplied for the READTEXT or WRITETEXT utility")
	}
	if len(parts) > 4 {
		return nil, p.newParseError(loc, "the column name has too many parts")
	}
	ref := &nodes.ColumnRef{
		Column: parts[len(parts)-1],
		Table:  parts[len(parts)-2],
		Loc:    nodes.Loc{Start: loc, End: p.prevEnd()},
	}
	if len(parts) > 2 {
		ref.Schema = parts[len(parts)-3]
	}
	if len(parts) > 3 {
		ref.Database = parts[0]
	}
	return ref, nil
}

// parseProcTextTarget parses modificationTextStatement (TSql160.g:14323):
// [ BULK ] column text_ptr [ TIMESTAMP = binary ], the text pointer a
// binary constant, a variable or an integer.
func (p *Parser) parseProcTextTarget() (*nodes.ColumnRef, nodes.ExprNode, error) {
	p.match(kwBULK)
	col, err := p.parseProcTextColumn()
	if err != nil {
		return nil, nil, err
	}
	ptr, err := p.parseProcValue(procValueKinds{binary: true, variable: true, integer: true})
	if err != nil {
		return nil, nil, err
	}
	if p.procIsWord("TIMESTAMP") {
		p.advance()
		if _, err := p.expect('='); err != nil {
			return nil, nil, err
		}
		if !procIsBinaryConst(p.cur) {
			return nil, nil, p.syntaxErrorAtCur()
		}
		p.advance()
	}
	return col, ptr, nil
}

// parseProcTextWithLog consumes WITH LOG (modificationTextStatementWithLog,
// TSql160.g:14313) and reports whether it did.
func (p *Parser) parseProcTextWithLog() (bool, error) {
	if p.cur.Type != kwWITH {
		return false, nil
	}
	p.advance() // WITH
	if _, err := p.expect(kwLOG); err != nil {
		return false, err
	}
	return true, nil
}
