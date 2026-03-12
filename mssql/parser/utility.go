// Package parser - utility.go implements T-SQL utility statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseUseStmt parses a USE database statement.
//
//	USE database
func (p *Parser) parseUseStmt() *nodes.UseStmt {
	loc := p.pos()
	p.advance() // consume USE

	stmt := &nodes.UseStmt{
		Loc: nodes.Loc{Start: loc},
	}

	if p.isIdentLike() {
		stmt.Database = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parsePrintStmt parses a PRINT statement.
//
//	PRINT expr
func (p *Parser) parsePrintStmt() *nodes.PrintStmt {
	loc := p.pos()
	p.advance() // consume PRINT

	stmt := &nodes.PrintStmt{
		Loc: nodes.Loc{Start: loc},
	}

	stmt.Expr = p.parseExpr()

	stmt.Loc.End = p.pos()
	return stmt
}

// parseRaiseErrorStmt parses a RAISERROR statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/raiserror-transact-sql
//
//	RAISERROR (msg, severity, state [, args]) [WITH options]
func (p *Parser) parseRaiseErrorStmt() *nodes.RaiseErrorStmt {
	loc := p.pos()
	p.advance() // consume RAISERROR

	stmt := &nodes.RaiseErrorStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// RAISERROR can use parens or not: RAISERROR('msg', 16, 1)
	if _, err := p.expect('('); err == nil {
		// Message
		stmt.Message = p.parseExpr()

		// Severity
		if _, ok := p.match(','); ok {
			stmt.Severity = p.parseExpr()
		}

		// State
		if _, ok := p.match(','); ok {
			stmt.State = p.parseExpr()
		}

		// Optional args
		var args []nodes.Node
		for {
			if _, ok := p.match(','); !ok {
				break
			}
			arg := p.parseExpr()
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
			if p.isIdentLike() || p.cur.Type == kwNOWAIT {
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

	stmt.Loc.End = p.pos()
	return stmt
}

// parseThrowStmt parses a THROW statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/throw-transact-sql
//
//	THROW [number, message, state]
func (p *Parser) parseThrowStmt() *nodes.ThrowStmt {
	loc := p.pos()
	p.advance() // consume THROW

	stmt := &nodes.ThrowStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// THROW without arguments = rethrow
	if p.cur.Type == ';' || p.cur.Type == tokEOF || p.cur.Type == kwEND ||
		p.isStatementStart() {
		stmt.Loc.End = p.pos()
		return stmt
	}

	// Error number
	stmt.ErrorNumber = p.parseExpr()

	// Message
	if _, ok := p.match(','); ok {
		stmt.Message = p.parseExpr()
	}

	// State
	if _, ok := p.match(','); ok {
		stmt.State = p.parseExpr()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseCheckpointStmt parses a CHECKPOINT statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/checkpoint-transact-sql
//
//	CHECKPOINT [ checkpoint_duration ]
func (p *Parser) parseCheckpointStmt() *nodes.CheckpointStmt {
	loc := p.pos()
	p.advance() // consume CHECKPOINT

	stmt := &nodes.CheckpointStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Optional duration
	if p.cur.Type == tokICONST || p.cur.Type == tokVARIABLE {
		stmt.Duration = p.parseExpr()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseReconfigureStmt parses a RECONFIGURE statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/reconfigure-transact-sql
//
//	RECONFIGURE [ WITH OVERRIDE ]
func (p *Parser) parseReconfigureStmt() *nodes.ReconfigureStmt {
	loc := p.pos()
	p.advance() // consume RECONFIGURE

	stmt := &nodes.ReconfigureStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// WITH OVERRIDE
	if p.cur.Type == kwWITH {
		next := p.peekNext()
		if next.Str != "" && strings.EqualFold(next.Str, "OVERRIDE") {
			p.advance() // WITH
			p.advance() // OVERRIDE
			stmt.WithOverride = true
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseShutdownStmt parses a SHUTDOWN statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/shutdown-transact-sql
//
//	SHUTDOWN [ WITH NOWAIT ]
func (p *Parser) parseShutdownStmt() *nodes.ShutdownStmt {
	loc := p.pos()
	p.advance() // consume SHUTDOWN

	stmt := &nodes.ShutdownStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// WITH NOWAIT
	if p.cur.Type == kwWITH {
		p.advance()
		if p.cur.Type == kwNOWAIT {
			p.advance()
			stmt.WithNoWait = true
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseKillStmt parses a KILL statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/kill-transact-sql
//
//	KILL { session_id | UOW } [ WITH STATUSONLY ]
//	KILL STATS JOB job_id
func (p *Parser) parseKillStmt() *nodes.KillStmt {
	loc := p.pos()
	p.advance() // consume KILL

	stmt := &nodes.KillStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// session_id or STATS JOB n or UOW
	stmt.SessionID = p.parseExpr()

	// WITH STATUSONLY
	if p.cur.Type == kwWITH {
		p.advance()
		if p.isIdentLike() && strings.EqualFold(p.cur.Str, "STATUSONLY") {
			p.advance()
			stmt.StatusOnly = true
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseReadtextStmt parses a READTEXT statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/readtext-transact-sql
//
//	READTEXT { table.column text_ptr offset size } [ HOLDLOCK ]
func (p *Parser) parseReadtextStmt() *nodes.ReadtextStmt {
	loc := p.pos()
	p.advance() // consume READTEXT

	stmt := &nodes.ReadtextStmt{
		Loc: nodes.Loc{Start: loc},
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
				Loc:    nodes.Loc{Start: colLoc, End: p.pos()},
			}
		} else {
			stmt.Column = &nodes.ColumnRef{
				Column: table,
				Loc:    nodes.Loc{Start: colLoc, End: p.pos()},
			}
		}
	}

	// text_ptr
	stmt.TextPtr = p.parseExpr()
	// offset
	stmt.Offset = p.parseExpr()
	// size
	stmt.Size = p.parseExpr()

	// HOLDLOCK
	if p.cur.Type == kwHOLDLOCK {
		p.advance()
		stmt.HoldLock = true
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseWritetextStmt parses a WRITETEXT statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/writetext-transact-sql
//
//	WRITETEXT { table.column text_ptr } [ WITH LOG ] { data }
func (p *Parser) parseWritetextStmt() *nodes.WritetextStmt {
	loc := p.pos()
	p.advance() // consume WRITETEXT

	stmt := &nodes.WritetextStmt{
		Loc: nodes.Loc{Start: loc},
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
				Loc:    nodes.Loc{Start: loc},
			}
		} else {
			stmt.Column = &nodes.ColumnRef{
				Column: table,
				Loc:    nodes.Loc{Start: loc},
			}
		}
	}

	// text_ptr
	stmt.TextPtr = p.parseExpr()

	// WITH LOG
	if p.cur.Type == kwWITH {
		next := p.peekNext()
		if next.Str != "" && strings.EqualFold(next.Str, "LOG") {
			p.advance() // WITH
			p.advance() // LOG
			stmt.WithLog = true
		}
	}

	// data
	stmt.Data = p.parseExpr()

	stmt.Loc.End = p.pos()
	return stmt
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
func (p *Parser) parseUpdatetextStmt() *nodes.UpdatetextStmt {
	loc := p.pos()
	p.advance() // consume UPDATETEXT

	stmt := &nodes.UpdatetextStmt{
		Loc: nodes.Loc{Start: loc},
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
				Loc:    nodes.Loc{Start: loc},
			}
		} else {
			stmt.DestColumn = &nodes.ColumnRef{
				Column: table,
				Loc:    nodes.Loc{Start: loc},
			}
		}
	}

	// dest_text_ptr
	stmt.DestTextPtr = p.parseExpr()
	// insert_offset (NULL or n)
	stmt.InsertOffset = p.parseExpr()
	// delete_length (NULL or n)
	stmt.DeleteLength = p.parseExpr()

	// WITH LOG
	if p.cur.Type == kwWITH {
		next := p.peekNext()
		if next.Str != "" && strings.EqualFold(next.Str, "LOG") {
			p.advance()
			p.advance()
			stmt.WithLog = true
		}
	}

	// inserted_data
	if p.cur.Type != ';' && p.cur.Type != tokEOF && p.cur.Type != kwGO {
		stmt.InsertedData = p.parseExpr()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseTruncateStmt parses a TRUNCATE TABLE statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/truncate-table-transact-sql
//
//	TRUNCATE TABLE name
func (p *Parser) parseTruncateStmt() *nodes.TruncateStmt {
	loc := p.pos()
	p.advance() // consume TRUNCATE

	// TABLE
	p.match(kwTABLE)

	stmt := &nodes.TruncateStmt{
		Loc: nodes.Loc{Start: loc},
	}

	stmt.Table = p.parseTableRef()

	stmt.Loc.End = p.pos()
	return stmt
}

// parseCreateDefaultStmt parses a CREATE DEFAULT statement.
// Caller has consumed CREATE DEFAULT.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-default-transact-sql
//
//	CREATE DEFAULT [ schema_name . ] default_name
//	AS constant_expression [ ; ]
func (p *Parser) parseCreateDefaultStmt() *nodes.SecurityStmt {
	loc := p.pos()

	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "DEFAULT",
		Loc:        nodes.Loc{Start: loc},
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
		expr := p.parseExpr()
		if expr != nil {
			opts = append(opts, &nodes.String{Str: "AS"})
			opts = append(opts, expr)
		}
	}
	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseCreateRuleStmt parses a CREATE RULE statement.
// Caller has consumed CREATE RULE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-rule-transact-sql
//
//	CREATE RULE [ schema_name . ] rule_name
//	AS condition_expression [ ; ]
func (p *Parser) parseCreateRuleStmt() *nodes.SecurityStmt {
	loc := p.pos()

	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "RULE",
		Loc:        nodes.Loc{Start: loc},
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
		expr := p.parseExpr()
		if expr != nil {
			opts = append(opts, &nodes.String{Str: "AS"})
			opts = append(opts, expr)
		}
	}
	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseAlterDatabaseScopedConfigStmt parses ALTER DATABASE SCOPED CONFIGURATION.
// Caller has consumed ALTER DATABASE SCOPED CONFIGURATION.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-database-scoped-configuration-transact-sql
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
//	    | QUERY_OPTIMIZER_HOTFIXES = { ON | OFF | PRIMARY }
//	    | ROW_MODE_MEMORY_GRANT_FEEDBACK = { ON | OFF }
//	    | TSQL_SCALAR_UDF_INLINING = { ON | OFF }
//	    | VERBOSE_TRUNCATION_WARNINGS = { ON | OFF }
//	    | XTP_PROCEDURE_EXECUTION_STATISTICS = { ON | OFF }
//	    | XTP_QUERY_EXECUTION_STATISTICS = { ON | OFF }
//	}
func (p *Parser) parseAlterDatabaseScopedConfigStmt() *nodes.SecurityStmt {
	loc := p.pos()

	stmt := &nodes.SecurityStmt{
		Action:     "ALTER",
		ObjectType: "DATABASE SCOPED CONFIGURATION",
		Loc:        nodes.Loc{Start: loc},
	}

	var opts []nodes.Node

	// CLEAR PROCEDURE_CACHE [plan_handle]
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "CLEAR") {
		p.advance() // consume CLEAR
		optStr := "CLEAR"
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "PROCEDURE_CACHE") {
			optStr += " PROCEDURE_CACHE"
			p.advance()
		}
		// Optional plan_handle (hex constant)
		if p.cur.Type == tokICONST || p.cur.Type == tokSCONST || p.isIdentLike() {
			optStr += " " + p.cur.Str
			p.advance()
		}
		opts = append(opts, &nodes.String{Str: optStr})
	} else {
		// [ FOR SECONDARY ] SET option = value
		forSecondary := false
		if p.cur.Type == kwFOR {
			next := p.peekNext()
			if next.Str != "" && matchesKeywordCI(next.Str, "SECONDARY") {
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

		// Parse SET option = value pairs
		for p.cur.Type != ';' && p.cur.Type != tokEOF && p.cur.Type != kwGO {
			if p.isIdentLike() || p.cur.Type == kwON || p.cur.Type == kwOFF {
				key := strings.ToUpper(p.cur.Str)
				p.advance()
				if p.cur.Type == '=' {
					p.advance() // consume '='
					if p.isIdentLike() || p.cur.Type == tokSCONST || p.cur.Type == tokICONST ||
						p.cur.Type == tokFCONST || p.cur.Type == kwON || p.cur.Type == kwOFF ||
						p.cur.Type == kwNULL || p.cur.Type == kwPRIMARY {
						val := strings.ToUpper(p.cur.Str)
						if p.cur.Type == tokSCONST {
							val = "'" + p.cur.Str + "'"
						}
						p.advance()
						opts = append(opts, &nodes.String{Str: key + "=" + val})
					} else {
						opts = append(opts, &nodes.String{Str: key + "="})
					}
				} else {
					opts = append(opts, &nodes.String{Str: key})
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

	stmt.Loc.End = p.pos()
	return stmt
}

// parseEnableDisableTriggerStmt parses ENABLE TRIGGER or DISABLE TRIGGER.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/enable-trigger-transact-sql
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/disable-trigger-transact-sql
//
//	{ ENABLE | DISABLE } TRIGGER { [ schema_name . ] trigger_name [ , ...n ] | ALL }
//	    ON { object_name | DATABASE | ALL SERVER }
func (p *Parser) parseEnableDisableTriggerStmt(enable bool) *nodes.EnableDisableTriggerStmt {
	loc := p.pos()
	p.advance() // consume ENABLE or DISABLE
	p.advance() // consume TRIGGER

	stmt := &nodes.EnableDisableTriggerStmt{
		Enable: enable,
		Loc:    nodes.Loc{Start: loc},
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
			stmt.OnObject = p.parseTableRef()
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}
