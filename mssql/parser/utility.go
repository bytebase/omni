// Package parser - utility.go implements T-SQL utility statement parsing.
package parser

import (
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseUseStmt parses a USE database statement.
//
//	USE database
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

	if p.isIdentLike() {
		stmt.Database = p.cur.Str
		p.advance()
	}

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

	// RAISERROR can use parens or not: RAISERROR('msg', 16, 1)
	if _, err := p.expect('('); err == nil {
		// Completion: after RAISERROR( → expression context
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			return nil, errCollecting
		}
		// Message
		stmt.Message, _ = p.parseExpr()

		// Severity
		if _, ok := p.match(','); ok {
			stmt.Severity, _ = p.parseExpr()
		}

		// State
		if _, ok := p.match(','); ok {
			stmt.State, _ = p.parseExpr()
		}

		// Optional args
		var args []nodes.Node
		for {
			if _, ok := p.match(','); !ok {
				break
			}
			arg, _ := p.parseExpr()
			args = append(args, arg)
		}
		if len(args) > 0 {
			stmt.Args = &nodes.List{Items: args}
		}

		_, _ = p.expect(')')
	}

	// WITH options (LOG, NOWAIT, SETERROR)
	if p.cur.Type == kwWITH {
		p.advance()
		var opts []nodes.Node
		for {
			if p.isAnyKeywordIdent() || p.cur.Type == kwNOWAIT {
				opts = append(opts, &nodes.String{Str: strings.ToUpper(p.cur.Str)})
				p.advance()
			} else {
				break
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

// parseThrowStmt parses a THROW statement.
//
// BNF: mssql/parser/bnf/throw-transact-sql.bnf
//
//	THROW [ { error_number | @local_variable }
//	    , { message | @local_variable }
//	    , { state | @local_variable } ]
//	[ ; ]
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
		stmt.Loc.End = p.prevEnd()
		return stmt, nil
	}

	// Error number
	stmt.ErrorNumber, _ = p.parseExpr()

	// Message
	if _, ok := p.match(','); ok {
		stmt.Message, _ = p.parseExpr()
	}

	// State
	if _, ok := p.match(','); ok {
		stmt.State, _ = p.parseExpr()
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
func (p *Parser) parseShutdownStmt() (*nodes.ShutdownStmt, error) {
	loc := p.pos()
	p.advance() // consume SHUTDOWN

	stmt := &nodes.ShutdownStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// WITH NOWAIT
	if p.cur.Type == kwWITH {
		p.advance()
		if p.cur.Type == kwNOWAIT {
			p.advance()
			stmt.WithNoWait = true
		}
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
	stmt.SessionID, _ = p.parseExpr()

	// WITH { STATUSONLY | COMMIT | ROLLBACK }
	if p.cur.Type == kwWITH {
		p.advance()
		if p.cur.Type == kwSTATUSONLY {
			p.advance()
			stmt.StatusOnly = true
		} else if p.cur.Type == kwCOMMIT {
			p.advance()
			stmt.WithAction = "COMMIT"
		} else if p.cur.Type == kwROLLBACK {
			p.advance()
			stmt.WithAction = "ROLLBACK"
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseKillStatsJobStmt parses a KILL STATS JOB statement.
//
// BNF: mssql/parser/bnf/kill-stats-job-transact-sql.bnf
//
//	KILL STATS JOB job_id
func (p *Parser) parseKillStatsJobStmt(loc int) (*nodes.KillStatsJobStmt, error) {
	p.advance() // consume STATS
	p.advance() // consume JOB

	stmt := &nodes.KillStatsJobStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	stmt.JobID, _ = p.parseExpr()

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseKillQueryNotificationStmt parses a KILL QUERY NOTIFICATION SUBSCRIPTION statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/kill-query-notification-subscription-transact-sql
//
//	KILL QUERY NOTIFICATION SUBSCRIPTION { ALL | subscription_id }
func (p *Parser) parseKillQueryNotificationStmt(loc int) (*nodes.KillQueryNotificationStmt, error) {
	p.advance() // consume QUERY
	p.advance() // consume NOTIFICATION

	// consume SUBSCRIPTION
	if p.cur.Type == kwSUBSCRIPTION {
		p.advance()
	}

	stmt := &nodes.KillQueryNotificationStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// { ALL | subscription_id }
	if p.cur.Type == kwALL {
		stmt.All = true
		p.advance()
	} else {
		stmt.SubscriptionID, _ = p.parseExpr()
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseReadtextStmt parses a READTEXT statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/readtext-transact-sql
//
//	READTEXT { table.column text_ptr offset size } [ HOLDLOCK ]
func (p *Parser) parseReadtextStmt() (*nodes.ReadtextStmt, error) {
	loc := p.pos()
	p.advance() // consume READTEXT

	stmt := &nodes.ReadtextStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// table.column - parse as column ref
	if p.isIdentLike() {
		colLoc := p.pos()
		table := p.cur.Str
		p.advance()
		if p.cur.Type == '.' {
			p.advance() // consume .
			col := ""
			if p.isIdentLike() {
				col = p.cur.Str
				p.advance()
			}
			stmt.Column = &nodes.ColumnRef{
				Table:  table,
				Column: col,
				Loc:    nodes.Loc{Start: colLoc, End: p.prevEnd()},
			}
		} else {
			stmt.Column = &nodes.ColumnRef{
				Column: table,
				Loc:    nodes.Loc{Start: colLoc, End: p.prevEnd()},
			}
		}
	}

	// text_ptr
	stmt.TextPtr, _ = p.parseExpr()
	// offset
	stmt.Offset, _ = p.parseExpr()
	// size
	stmt.Size, _ = p.parseExpr()

	// HOLDLOCK
	if p.cur.Type == kwHOLDLOCK {
		p.advance()
		stmt.HoldLock = true
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseWritetextStmt parses a WRITETEXT statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/writetext-transact-sql
//
//	WRITETEXT { table.column text_ptr } [ WITH LOG ] { data }
func (p *Parser) parseWritetextStmt() (*nodes.WritetextStmt, error) {
	loc := p.pos()
	p.advance() // consume WRITETEXT

	stmt := &nodes.WritetextStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// table.column
	if p.isIdentLike() {
		table := p.cur.Str
		p.advance()
		if p.cur.Type == '.' {
			p.advance() // .
			col := ""
			if p.isIdentLike() {
				col = p.cur.Str
				p.advance()
			}
			stmt.Column = &nodes.ColumnRef{
				Table:  table,
				Column: col,
				Loc:    nodes.Loc{Start: loc, End: -1},
			}
		} else {
			stmt.Column = &nodes.ColumnRef{
				Column: table,
				Loc:    nodes.Loc{Start: loc, End: -1},
			}
		}
	}

	// text_ptr
	stmt.TextPtr, _ = p.parseExpr()

	// WITH LOG
	if p.cur.Type == kwWITH {
		next := p.peekNext()
		if next.Type == kwLOG {
			p.advance() // WITH
			p.advance() // LOG
			stmt.WithLog = true
		}
	}

	// data
	stmt.Data, _ = p.parseExpr()

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
func (p *Parser) parseUpdatetextStmt() (*nodes.UpdatetextStmt, error) {
	loc := p.pos()
	p.advance() // consume UPDATETEXT

	stmt := &nodes.UpdatetextStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// dest_table.dest_column dest_text_ptr
	if p.isIdentLike() {
		table := p.cur.Str
		p.advance()
		if p.cur.Type == '.' {
			p.advance()
			col := ""
			if p.isIdentLike() {
				col = p.cur.Str
				p.advance()
			}
			stmt.DestColumn = &nodes.ColumnRef{
				Table:  table,
				Column: col,
				Loc:    nodes.Loc{Start: loc, End: -1},
			}
		} else {
			stmt.DestColumn = &nodes.ColumnRef{
				Column: table,
				Loc:    nodes.Loc{Start: loc, End: -1},
			}
		}
	}

	// dest_text_ptr
	stmt.DestTextPtr, _ = p.parseExpr()
	// insert_offset (NULL or n)
	stmt.InsertOffset, _ = p.parseExpr()
	// delete_length (NULL or n)
	stmt.DeleteLength, _ = p.parseExpr()

	// WITH LOG
	if p.cur.Type == kwWITH {
		next := p.peekNext()
		if next.Type == kwLOG {
			p.advance()
			p.advance()
			stmt.WithLog = true
		}
	}

	// inserted_data
	if p.cur.Type != ';' && p.cur.Type != tokEOF && p.cur.Type != kwGO {
		stmt.InsertedData, _ = p.parseExpr()
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

	stmt.Table , _ = p.parseTableRef()

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
							p.advance() // consume TO
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

// parseAlterDatabaseScopedConfigStmt parses ALTER DATABASE SCOPED CONFIGURATION.
// Caller has consumed ALTER DATABASE SCOPED CONFIGURATION.
//
// BNF: mssql/parser/bnf/alter-database-scoped-configuration-transact-sql.bnf
//
//	ALTER DATABASE SCOPED CONFIGURATION
//	{
//	    { [ FOR SECONDARY ] SET <set_options> }
//	}
//	| CLEAR PROCEDURE_CACHE [plan_handle]
//	| SET < set_options >
//	[;]
//
//	< set_options > ::=
//	{
//	      ACCELERATED_PLAN_FORCING = { ON | OFF }
//	    | ALLOW_STALE_VECTOR_INDEX = { ON | OFF }
//	    | ASYNC_STATS_UPDATE_WAIT_AT_LOW_PRIORITY = { ON | OFF }
//	    | BATCH_MODE_ADAPTIVE_JOINS = { ON | OFF }
//	    | BATCH_MODE_MEMORY_GRANT_FEEDBACK = { ON | OFF }
//	    | BATCH_MODE_ON_ROWSTORE = { ON | OFF }
//	    | CE_FEEDBACK = { ON | OFF }
//	    | DEFERRED_COMPILATION_TV = { ON | OFF }
//	    | DOP_FEEDBACK = { ON | OFF }
//	    | ELEVATE_ONLINE = { OFF | WHEN_SUPPORTED | FAIL_UNSUPPORTED }
//	    | ELEVATE_RESUMABLE = { OFF | WHEN_SUPPORTED | FAIL_UNSUPPORTED }
//	    | EXEC_QUERY_STATS_FOR_SCALAR_FUNCTIONS = { ON | OFF }
//	    | FORCE_SHOWPLAN_RUNTIME_PARAMETER_COLLECTION = { ON | OFF }
//	    | FULLTEXT_INDEX_VERSION = <version>
//	    | IDENTITY_CACHE = { ON | OFF }
//	    | INTERLEAVED_EXECUTION_TVF = { ON | OFF }
//	    | ISOLATE_SECURITY_POLICY_CARDINALITY = { ON | OFF }
//	    | GLOBAL_TEMPORARY_TABLE_AUTO_DROP = { ON | OFF }
//	    | LAST_QUERY_PLAN_STATS = { ON | OFF }
//	    | LEDGER_DIGEST_STORAGE_ENDPOINT = { <endpoint URL string> | OFF }
//	    | LEGACY_CARDINALITY_ESTIMATION = { ON | OFF | PRIMARY }
//	    | LIGHTWEIGHT_QUERY_PROFILING = { ON | OFF }
//	    | MAXDOP = { <value> | PRIMARY }
//	    | MEMORY_GRANT_FEEDBACK_PERCENTILE_GRANT = { ON | OFF }
//	    | MEMORY_GRANT_FEEDBACK_PERSISTENCE = { ON | OFF }
//	    | OPTIMIZE_FOR_AD_HOC_WORKLOADS = { ON | OFF }
//	    | OPTIMIZED_PLAN_FORCING = { ON | OFF }
//	    | OPTIMIZED_SP_EXECUTESQL = { ON | OFF }
//	    | OPTIONAL_PARAMETER_OPTIMIZATION = { ON | OFF }
//	    | PARAMETER_SENSITIVE_PLAN_OPTIMIZATION = { ON | OFF }
//	    | PARAMETER_SNIFFING = { ON | OFF | PRIMARY }
//	    | PAUSED_RESUMABLE_INDEX_ABORT_DURATION_MINUTES = <time>
//	    | PREVIEW_FEATURES = { ON | OFF }
//	    | QUERY_OPTIMIZER_HOTFIXES = { ON | OFF | PRIMARY }
//	    | READABLE_SECONDARY_TEMPORARY_STATS_AUTO_CREATE = { ON | OFF | PRIMARY }
//	    | READABLE_SECONDARY_TEMPORARY_STATS_AUTO_UPDATE = { ON | OFF | PRIMARY }
//	    | ROW_MODE_MEMORY_GRANT_FEEDBACK = { ON | OFF }
//	    | TSQL_SCALAR_UDF_INLINING = { ON | OFF }
//	    | VERBOSE_TRUNCATION_WARNINGS = { ON | OFF }
//	    | XTP_PROCEDURE_EXECUTION_STATISTICS = { ON | OFF }
//	    | XTP_QUERY_EXECUTION_STATISTICS = { ON | OFF }
//	}
func (p *Parser) parseAlterDatabaseScopedConfigStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()

	stmt := &nodes.SecurityStmt{
		Action:     "ALTER",
		ObjectType: "DATABASE SCOPED CONFIGURATION",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	var opts []nodes.Node

	// CLEAR PROCEDURE_CACHE [plan_handle]
	if p.cur.Type == kwCLEAR {
		p.advance() // consume CLEAR
		optStr := "CLEAR"
		if p.cur.Type == kwPROCEDURE_CACHE {
			optStr += " PROCEDURE_CACHE"
			p.advance()
		}
		// Optional plan_handle (hex constant)
		if p.cur.Type == tokICONST || p.cur.Type == tokSCONST || p.isAnyKeywordIdent() {
			optStr += " " + p.cur.Str
			p.advance()
		}
		opts = append(opts, &nodes.String{Str: optStr})
	} else {
		// [ FOR SECONDARY ] SET option = value
		forSecondary := false
		if p.cur.Type == kwFOR {
			next := p.peekNext()
			if next.Type == kwSECONDARY {
				p.advance() // consume FOR
				p.advance() // consume SECONDARY
				forSecondary = true
			}
		}
		if forSecondary {
			opts = append(opts, &nodes.String{Str: "FOR SECONDARY"})
		}

		if p.cur.Type == kwSET {
			p.advance() // consume SET
		}

		// Parse SET option = value pairs as structured SecurityPrincipalOption nodes
		for p.cur.Type != ';' && p.cur.Type != tokEOF && p.cur.Type != kwGO {
			if p.isAnyKeywordIdent() || p.cur.Type == kwON || p.cur.Type == kwOFF {
				optLoc := p.pos()
				key := strings.ToUpper(p.cur.Str)
				p.advance()
				if p.cur.Type == '=' {
					p.advance() // consume '='
					val := ""
					if p.isAnyKeywordIdent() || p.cur.Type == tokSCONST || p.cur.Type == tokICONST ||
						p.cur.Type == tokFCONST || p.cur.Type == kwON || p.cur.Type == kwOFF ||
						p.cur.Type == kwNULL || p.cur.Type == kwPRIMARY {
						if p.cur.Type == tokSCONST {
							val = p.cur.Str
						} else {
							val = strings.ToUpper(p.cur.Str)
						}
						p.advance()
					}
					opts = append(opts, &nodes.SecurityPrincipalOption{
						Name:  key,
						Value: val,
						Loc:   nodes.Loc{Start: optLoc, End: p.prevEnd()},
					})
				} else {
					opts = append(opts, &nodes.SecurityPrincipalOption{
						Name: key,
						Loc:  nodes.Loc{Start: optLoc, End: p.prevEnd()},
					})
				}
			} else if p.cur.Type == ',' {
				p.advance() // skip commas
			} else {
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
			p.advance()
			if p.matchIdentCI("SERVER") {
				stmt.OnAllServer = true
			}
		} else {
			stmt.OnObject , _ = p.parseTableRef()
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

	// Optional 'username'
	if p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()

		// Optional WITH NORESET
		if p.cur.Type == kwWITH {
			p.advance() // consume WITH
			if p.matchIdentCI("NORESET") {
				opts := &nodes.List{}
				opts.Items = append(opts.Items, &nodes.String{Str: "NORESET"})
				stmt.Options = opts
			}
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCopyIntoStmt parses a COPY INTO statement (Azure Synapse / Fabric).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/copy-into-transact-sql
//
//	COPY INTO [ schema. ] table_name
//	[ ( Column_name [ DEFAULT Default_value ] [ Field_number ] [ ,...n ] ) ]
//	FROM '<external_location>' [ , ...n ]
//	WITH
//	(
//	  [ FILE_TYPE = { 'CSV' | 'PARQUET' | 'ORC' } ]
//	  [ , FILE_FORMAT = EXTERNAL FILE FORMAT OBJECT ]
//	  [ , CREDENTIAL = ( AZURE CREDENTIAL ) ]
//	  [ , ERRORFILE = '[http(s)://storageaccount/container]/errorfile_directory[/]' ]
//	  [ , ERRORFILE_CREDENTIAL = ( AZURE CREDENTIAL ) ]
//	  [ , MAXERRORS = max_errors ]
//	  [ , COMPRESSION = { 'Gzip' | 'DefaultCodec' | 'Snappy' } ]
//	  [ , FIELDQUOTE = 'string_delimiter' ]
//	  [ , FIELDTERMINATOR = 'field_terminator' ]
//	  [ , ROWTERMINATOR = 'row_terminator' ]
//	  [ , FIRSTROW = first_row ]
//	  [ , DATEFORMAT = 'date_format' ]
//	  [ , ENCODING = { 'UTF8' | 'UTF16' } ]
//	  [ , IDENTITY_INSERT = { 'ON' | 'OFF' } ]
//	  [ , AUTO_CREATE_TABLE = { 'ON' | 'OFF' } ]
//	)
func (p *Parser) parseCopyIntoStmt() (*nodes.CopyIntoStmt, error) {
	loc := p.pos()
	p.advance() // consume COPY

	// INTO keyword
	p.matchIdentCI("INTO")

	stmt := &nodes.CopyIntoStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Table name (possibly qualified: schema.table)
	stmt.Table , _ = p.parseTableRef()

	// Optional column list: ( Column_name [ DEFAULT value ] [ field_number ] [,...n] )
	if p.cur.Type == '(' {
		p.advance() // consume '('
		var cols []nodes.Node
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			// Each entry is: Column_name [ DEFAULT value ] [ field_number ]
			if p.isIdentLike() || p.cur.Type == '[' {
				col := &nodes.CopyIntoColumn{
					Loc: nodes.Loc{Start: p.pos(), End: -1},
				}
				col.Name = p.cur.Str
				p.advance()

				// Check for DEFAULT value
				if p.matchIdentCI("DEFAULT") {
					col.DefaultValue, _ = p.parseExpr()
				}

				// Check for field_number (integer)
				if p.cur.Type == tokICONST {
					if n, err := strconv.Atoi(p.cur.Str); err == nil {
						col.FieldNumber = n
					}
					p.advance()
				}

				col.Loc.End = p.prevEnd()
				cols = append(cols, col)
			} else {
				// unexpected token in column list - break to avoid silent consumption
				break
			}

			if _, ok := p.match(','); !ok {
				break
			}
		}
		p.match(')')
		if len(cols) > 0 {
			stmt.ColumnList = &nodes.List{Items: cols}
		}
	}

	// FROM '<external_location>' [ , ...n ]
	if p.cur.Type == kwFROM {
		p.advance() // consume FROM
		var sources []nodes.Node
		for {
			if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
				sources = append(sources, &nodes.String{Str: p.cur.Str})
				p.advance()
			} else {
				break
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
		if len(sources) > 0 {
			stmt.Sources = &nodes.List{Items: sources}
		}
	}

	// WITH ( option [,...n] )
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		if p.cur.Type == '(' {
			p.advance() // consume '('
			var opts []nodes.Node
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				opt := p.parseCopyIntoOption()
				if opt != nil {
					opts = append(opts, opt)
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
	if !p.isAnyKeywordIdent() && p.cur.Type != kwFILE {
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
	stmt.Name , _ = p.parseTableRef()

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

// parseCreateExternalTableAsSelectStmt parses CREATE EXTERNAL TABLE ... AS SELECT (CETAS).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-external-table-as-select-transact-sql
//
//	CREATE EXTERNAL TABLE { [ [ database_name . [ schema_name ] . ] | schema_name . ] table_name }
//	    [ (column_name [ , ...n ] ) ]
//	    WITH (
//	        LOCATION = 'hdfs_folder' | '<prefix>://<path>[:<port>]' ,
//	        DATA_SOURCE = external_data_source_name ,
//	        FILE_FORMAT = external_file_format_name
//	        [ , <reject_options> [ , ...n ] ]
//	    )
//	    AS <select_statement>
//
//	<reject_options> ::=
//	{
//	    | REJECT_TYPE = value | percentage
//	    | REJECT_VALUE = reject_value
//	    | REJECT_SAMPLE_VALUE = reject_sample_value
//	}
func (p *Parser) parseCreateExternalTableAsSelectStmt() (*nodes.CreateExternalTableAsSelectStmt, error) {
	loc := p.pos()
	// EXTERNAL TABLE already consumed by caller

	stmt := &nodes.CreateExternalTableAsSelectStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Table name
	stmt.Name , _ = p.parseTableRef()

	// Optional column list: parse as structured ColumnDef nodes
	if p.cur.Type == '(' {
		p.advance() // consume '('
		var cols []nodes.Node
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			col, _ := p.parseColumnDef()
			if col != nil {
				cols = append(cols, col)
			} else {
				break
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
		p.match(')') // consume ')'
		if len(cols) > 0 {
			stmt.Columns = &nodes.List{Items: cols}
		}
	}

	// WITH ( options )
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
			p.match(')')
			if len(opts) > 0 {
				stmt.Options = &nodes.List{Items: opts}
			}
		}
	}

	// AS <select_statement>
	if p.cur.Type == kwAS {
		p.advance() // consume AS
		if p.cur.Type == kwSELECT || p.cur.Type == kwWITH {
			stmt.Query, _ = p.parseSelectStmt()
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
	stmt.SourceName , _ = p.parseTableRef()

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

// parsePredictStmt parses a PREDICT statement.
// Caller has consumed PREDICT.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/predict-transact-sql
//
//	PREDICT (
//	  MODEL = @model | model_literal,
//	  DATA = object AS <table_alias>
//	  [, RUNTIME = ONNX ]
//	)
//	WITH ( <result_set_definition> )
//
//	<result_set_definition> ::=
//	  { column_name data_type [ COLLATE collation_name ] [ NULL | NOT NULL ] } [,...n]
func (p *Parser) parsePredictStmt() (*nodes.PredictStmt, error) {
	loc := p.pos()
	// PREDICT already consumed by caller

	stmt := &nodes.PredictStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Expect '('
	p.match('(')

	// Parse arguments: MODEL = expr, DATA = table AS alias [, RUNTIME = ONNX]
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		if p.cur.Type == ',' {
			p.advance()
			continue
		}

		if p.isAnyKeywordIdent() {
			key := strings.ToUpper(p.cur.Str)
			p.advance()

			if p.cur.Type == '=' {
				p.advance() // consume '='
			}

			switch key {
			case "MODEL":
				// MODEL = @variable | 'literal' | (subquery)
				stmt.Model, _ = p.parseExpr()
			case "DATA":
				// DATA = table_source AS alias
				stmt.Data, _ = p.parseExpr()
				// Optional AS alias
				if p.cur.Type == kwAS {
					p.advance() // consume AS
					if p.isIdentLike() {
						stmt.DataAlias = p.cur.Str
						p.advance()
					}
				}
			case "RUNTIME":
				// RUNTIME = ONNX
				if p.isAnyKeywordIdent() {
					stmt.Runtime = strings.ToUpper(p.cur.Str)
					p.advance()
				}
			default:
				// Skip unknown args
				if p.isAnyKeywordIdent() || p.cur.Type == tokSCONST || p.cur.Type == tokICONST {
					p.advance()
				}
			}
		} else {
			p.advance()
		}
	}

	p.match(')') // consume ')'

	// WITH ( result_set_definition )
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		if p.cur.Type == '(' {
			p.advance() // consume '('
			var cols []nodes.Node
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				if p.cur.Type == ',' {
					p.advance()
					continue
				}
				// column_name data_type [COLLATE collation] [NULL | NOT NULL]
				col := p.parsePredictColumnDef()
				if col != nil {
					cols = append(cols, col)
				}
			}
			p.match(')') // consume ')'
			if len(cols) > 0 {
				stmt.WithColumns = &nodes.List{Items: cols}
			}
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parsePredictColumnDef parses a column definition in PREDICT's WITH clause.
// column_name data_type [COLLATE collation_name] [NULL | NOT NULL]
func (p *Parser) parsePredictColumnDef() *nodes.ColumnDef {
	if !p.isIdentLike() {
		return nil
	}
	loc := p.pos()

	colName := p.cur.Str
	p.advance()

	// data_type
	dt , _ := p.parseDataType()

	col := &nodes.ColumnDef{
		Name:     colName,
		DataType: dt,
		Loc:      nodes.Loc{Start: loc, End: -1},
	}

	// Optional COLLATE
	if p.cur.Type == kwCOLLATE {
		p.advance()
		if p.isIdentLike() {
			col.Collation = p.cur.Str
			p.advance()
		}
	}

	// Optional NULL / NOT NULL
	if p.cur.Type == kwNOT {
		p.advance()
		if p.cur.Type == kwNULL {
			p.advance()
			col.Nullable = &nodes.NullableSpec{NotNull: true, Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}
		}
	} else if p.cur.Type == kwNULL {
		p.advance()
		col.Nullable = &nodes.NullableSpec{NotNull: false, Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}
	}

	col.Loc.End = p.prevEnd()
	return col
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
	stmt.Name , _ = p.parseTableRef()

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
			stmt.Query, _ = p.parseSelectStmt()
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
	stmt.DataType , _ = p.parseDataType()
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
	stmt.Table , _ = p.parseTableRef()

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

	col.DataType , _ = p.parseDataType()

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
