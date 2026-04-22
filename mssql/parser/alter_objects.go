// Package parser - alter_objects.go implements T-SQL ALTER DATABASE and ALTER INDEX parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseAlterDatabaseStmt parses an ALTER DATABASE statement.
//
// BNF: mssql/parser/bnf/alter-database-transact-sql.bnf
//
//	ALTER DATABASE { database_name | CURRENT }
//	{
//	    SET <option_spec> [ ,...n ] [ WITH <termination> ]
//	  | ADD FILE <filespec> [ ,...n ] [ TO FILEGROUP filegroup_name ]
//	  | ADD LOG FILE <filespec> [ ,...n ]
//	  | REMOVE FILE logical_file_name
//	  | MODIFY FILE <filespec>
//	  | ADD FILEGROUP filegroup_name [ CONTAINS FILESTREAM | CONTAINS MEMORY_OPTIMIZED_DATA ]
//	  | REMOVE FILEGROUP filegroup_name
//	  | MODIFY FILEGROUP filegroup_name
//	      { <filegroup_updatability_option> | DEFAULT | NAME = new_name
//	        | AUTOGROW_SINGLE_FILE | AUTOGROW_ALL_FILES }
//	  | MODIFY NAME = new_database_name
//	  | MODIFY ( <azure_option> [ ,...n ] ) [ WITH MANUAL_CUTOVER ]
//	  | COLLATE collation_name
//	  | REBUILD LOG [ ON <filespec> ]
//	  | PERFORM_CUTOVER
//	}
//
//	<option_spec> ::=
//	{
//	    <db_state_option>         -- ONLINE | OFFLINE | EMERGENCY
//	  | <db_user_access_option>   -- SINGLE_USER | RESTRICTED_USER | MULTI_USER
//	  | <db_update_option>        -- READ_ONLY | READ_WRITE
//	  | <recovery_option>         -- RECOVERY { FULL | BULK_LOGGED | SIMPLE }
//	                              -- PAGE_VERIFY { CHECKSUM | TORN_PAGE_DETECTION | NONE }
//	                              -- TORN_PAGE_DETECTION { ON | OFF }
//	  | <auto_option>             -- AUTO_CLOSE|AUTO_SHRINK|AUTO_CREATE_STATISTICS|AUTO_UPDATE_STATISTICS[_ASYNC] { ON | OFF }
//	  | <sql_option>              -- ANSI_NULL_DEFAULT|ANSI_NULLS|ANSI_PADDING|ANSI_WARNINGS|ARITHABORT
//	                              --   |CONCAT_NULL_YIELDS_NULL|NUMERIC_ROUNDABORT|QUOTED_IDENTIFIER|RECURSIVE_TRIGGERS { ON | OFF }
//	  | <cursor_option>           -- CURSOR_CLOSE_ON_COMMIT|CURSOR_DEFAULT { ON|OFF|LOCAL|GLOBAL }
//	  | <containment_option>      -- CONTAINMENT = { NONE | PARTIAL }
//	  | <db_encryption_option>    -- ENCRYPTION { ON | OFF | SUSPEND | RESUME }
//	  | <change_tracking_option>  -- CHANGE_TRACKING = { ON [(opts)] | OFF }
//	  | <snapshot_option>         -- ALLOW_SNAPSHOT_ISOLATION|READ_COMMITTED_SNAPSHOT|MEMORY_OPTIMIZED_ELEVATE_TO_SNAPSHOT { ON | OFF }
//	  | <query_store_options>     -- QUERY_STORE = { ON [(opts)] | OFF [(FORCED)] } | QUERY_STORE CLEAR [ALL]
//	  | <parameterization_option> -- PARAMETERIZATION { SIMPLE | FORCED }
//	  | <external_access_option>  -- DB_CHAINING|TRUSTWORTHY { ON | OFF }
//	                              -- DEFAULT_FULLTEXT_LANGUAGE|DEFAULT_LANGUAGE = value
//	                              -- NESTED_TRIGGERS|TRANSFORM_NOISE_WORDS = { ON | OFF }
//	                              -- TWO_DIGIT_YEAR_CUTOFF = number
//	  | <delayed_durability_option> -- DELAYED_DURABILITY = { DISABLED | ALLOWED | FORCED }
//	  | <target_recovery_time_option> -- TARGET_RECOVERY_TIME = number { SECONDS | MINUTES }
//	  | <service_broker_option>   -- ENABLE_BROKER|DISABLE_BROKER|NEW_BROKER|ERROR_BROKER_CONVERSATIONS
//	                              -- HONOR_BROKER_PRIORITY { ON | OFF }
//	  | <accelerated_database_recovery> -- ACCELERATED_DATABASE_RECOVERY = { ON | OFF }
//	  | <mixed_page_allocation_option>  -- MIXED_PAGE_ALLOCATION { ON | OFF }
//	  | <optimized_locking>       -- OPTIMIZED_LOCKING = { ON | OFF }
//	  | COMPATIBILITY_LEVEL = number
//	  | <date_correlation_optimization_option> -- DATE_CORRELATION_OPTIMIZATION { ON | OFF }
//	  | <temporal_history_retention> -- TEMPORAL_HISTORY_RETENTION { ON | OFF }
//	  | <data_retention_policy>   -- DATA_RETENTION { ON | OFF }
//	  | FILESTREAM ( <FILESTREAM_option> )
//	  | <HADR_options>
//	  | <persistent_log_buffer_option>
//	  | <suspend_for_snapshot_backup>
//	}
//
//	<termination> ::=
//	{
//	    ROLLBACK AFTER number [ SECONDS ]
//	  | ROLLBACK IMMEDIATE
//	  | NO_WAIT
//	}
//
//	<filespec> ::=
//	(
//	    NAME = logical_file_name
//	    [ , NEWNAME = new_logical_name ]
//	    [ , FILENAME = 'os_file_name' ]
//	    [ , SIZE = size [ KB | MB | GB | TB ] ]
//	    [ , MAXSIZE = { max_size [ KB | MB | GB | TB ] | UNLIMITED } ]
//	    [ , FILEGROWTH = growth_increment [ KB | MB | GB | TB | % ] ]
//	    [ , OFFLINE ]
//	)
//
//	<filegroup_updatability_option> ::=
//	{
//	    { READONLY | READWRITE }
//	    | { READ_ONLY | READ_WRITE }
//	}
func (p *Parser) parseAlterDatabaseStmt() (*nodes.AlterDatabaseStmt, error) {
	loc := p.pos()

	stmt := &nodes.AlterDatabaseStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Database name or CURRENT. CURRENT is a Core keyword reserved as a
	// special "current database" token in this position per SqlScriptDOM.
	if p.isIdentLike() || p.cur.Type == kwCURRENT {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// Dispatch on action keyword
	if p.cur.Type == kwSET {
		p.advance() // consume SET
		stmt.Action = "SET"
		if err := p.parseAlterDatabaseSetOptions(stmt); err != nil {
			return nil, err
		}
	} else if p.cur.Type == kwCOLLATE {
		p.advance() // consume COLLATE
		stmt.Action = "COLLATE"
		if id, ok := p.parseIdentifier(); ok {
			stmt.TargetName = id
		}
	} else if p.cur.Type == kwADD {
		p.advance() // consume ADD
		stmt.Action = "ADD"
		p.parseAlterDatabaseAdd(stmt)
	} else if p.isAnyKeywordIdent() {
		action := strings.ToUpper(p.cur.Str)
		switch action {
		case "MODIFY":
			p.advance() // consume MODIFY
			stmt.Action = "MODIFY"
			p.parseAlterDatabaseModify(stmt)
		case "REMOVE":
			p.advance() // consume REMOVE
			stmt.Action = "REMOVE"
			p.parseAlterDatabaseRemove(stmt)
		case "REBUILD":
			// REBUILD LOG [ ON <filespec> ]
			p.advance() // consume REBUILD
			stmt.Action = "REBUILD"
			stmt.SubAction = "LOG"
			if p.cur.Type == kwLOG {
				p.advance() // consume LOG
			}
			// Optional ON <filespec>
			if p.cur.Type == kwON {
				p.advance() // consume ON
				if p.cur.Type == '(' {
					stmt.FileSpecs = &nodes.List{Items: []nodes.Node{p.parseDatabaseFileSpec()}}
				}
			}
		case "PERFORM_CUTOVER":
			p.advance() // consume PERFORM_CUTOVER
			stmt.Action = "PERFORM_CUTOVER"
		default:
			// Unknown action - record and collect remaining tokens as structured key=value options
			stmt.Action = action
			p.advance()
			var opts []nodes.Node
			for p.cur.Type != tokEOF && p.cur.Type != ';' && !p.isStatementStart() {
				if p.cur.Type == ',' {
					p.advance()
					continue
				}
				opt := p.parseAlterDatabaseUnknownOption()
				if opt != "" {
					opts = append(opts, &nodes.String{Str: opt})
				} else {
					break
				}
			}
			if len(opts) > 0 {
				stmt.Options = &nodes.List{Items: opts}
			}
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterDatabaseSetOptions parses SET <option_spec> [,...n] [WITH <termination>].
// Returns a syntax error if no valid option follows SET.
func (p *Parser) parseAlterDatabaseSetOptions(stmt *nodes.AlterDatabaseStmt) error {
	var opts []nodes.Node

	for {
		opt := p.parseAlterDatabaseSetOption()
		if opt == "" {
			break
		}
		opts = append(opts, &nodes.String{Str: opt})

		if p.cur.Type == ',' {
			p.advance()
		} else {
			break
		}
	}

	// ALTER DATABASE ... SET requires at least one valid option.
	if len(opts) == 0 {
		return p.syntaxErrorAtCur()
	}

	stmt.Options = &nodes.List{Items: opts}

	// WITH <termination>
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		stmt.Termination = p.parseAlterDatabaseTermination()
	}
	return nil
}

// dbSetOptions is the valid option set for ALTER DATABASE SET <option>.
// Derived from SqlScriptDOM: SimpleDbOptionsHelper, OnOffSimpleDbOptionsHelper,
// DatabaseOptionKindHelper, RecoveryDbOptionsHelper, PageVerifyDbOptionsHelper,
// PartnerDbOptionsHelper, and special-case options (HADR, WITNESS, etc.).
var dbSetOptions = newOptionSet(
	// Registered keywords that are valid database SET options.
	kwENCRYPTION,
	kwCHANGE_TRACKING,
	kwCONTAINMENT,
	kwPARAMETERIZATION,
	kwDELAYED_DURABILITY,
	kwFILESTREAM,
	kwHADR,
	kwREAD_ONLY,
	kwOFFLINE,
).withIdents(
	// SimpleDbOptionsHelper — standalone keywords (no value).
	"ONLINE", "EMERGENCY",
	"SINGLE_USER", "MULTI_USER", "RESTRICTED_USER",
	"READ_WRITE",
	"ENABLE_BROKER", "DISABLE_BROKER", "NEW_BROKER", "ERROR_BROKER_CONVERSATIONS",
	// OnOffSimpleDbOptionsHelper — ON/OFF options.
	"CURSOR_CLOSE_ON_COMMIT", "AUTO_CLOSE", "AUTO_CREATE_STATISTICS", "AUTO_SHRINK",
	"AUTO_UPDATE_STATISTICS", "ANSI_NULL_DEFAULT", "ANSI_NULLS", "ANSI_PADDING",
	"ANSI_WARNINGS", "ARITHABORT", "CONCAT_NULL_YIELDS_NULL", "NUMERIC_ROUNDABORT",
	"QUOTED_IDENTIFIER", "RECURSIVE_TRIGGERS", "TORN_PAGE_DETECTION",
	"DB_CHAINING", "TRUSTWORTHY", "AUTO_UPDATE_STATISTICS_ASYNC",
	"DATE_CORRELATION_OPTIMIZATION", "ALLOW_SNAPSHOT_ISOLATION", "READ_COMMITTED_SNAPSHOT",
	"SUPPLEMENTAL_LOGGING", "HONOR_BROKER_PRIORITY", "VARDECIMAL_STORAGE_FORMAT",
	"NESTED_TRIGGERS", "TRANSFORM_NOISE_WORDS", "MEMORY_OPTIMIZED_ELEVATE_TO_SNAPSHOT",
	"MIXED_PAGE_ALLOCATION", "TEMPORAL_HISTORY_RETENTION", "DATA_RETENTION", "LEDGER",
	// DatabaseOptionKindHelper — = value options.
	"COMPATIBILITY_LEVEL", "DEFAULT_FULLTEXT_LANGUAGE", "DEFAULT_LANGUAGE",
	"TWO_DIGIT_YEAR_CUTOFF", "EDITION", "SERVICE_OBJECTIVE",
	"QUERY_STORE", "CATALOG_COLLATION", "AUTOMATIC_TUNING",
	// Other special handlers.
	"RECOVERY", "PAGE_VERIFY", "TARGET_RECOVERY_TIME",
	"ACCELERATED_DATABASE_RECOVERY", "PARTNER", "WITNESS",
	"CURSOR_DEFAULT", "OPTIMIZED_LOCKING",
	// Additional (MANUAL_CUTOVER, PERFORM_CUTOVER are handled at ALTER DATABASE level, not SET level).
)

// Per-option value enums. When an option key is found here, its value must be
// one of the listed identifiers; any other keyword/ident is rejected (the
// parseAlterDatabaseSetOption function returns "" to skip this SET fragment,
// matching how it already treats unparseable options).
//
// This is the POC scope for the DDL option-value strictness plan
// (docs/plans/2026-04-22-ddl-option-strictness.md). Additional option → enum
// mappings belong in the follow-up batches.
var dbSetOptionValueEnums = map[string]optionSet{
	"RECOVERY":         newOptionSet().withIdents("FULL", "SIMPLE", "BULK_LOGGED"),
	"PAGE_VERIFY":      newOptionSet().withIdents("CHECKSUM", "TORN_PAGE_DETECTION", "NONE"),
	"PARAMETERIZATION": newOptionSet().withIdents("SIMPLE", "FORCED"),
	"CURSOR_DEFAULT":   newOptionSet().withIdents("LOCAL", "GLOBAL"),
}

// parseAlterDatabaseSetOption parses a single SET option and returns it as a string.
// Returns empty string if no option can be parsed.
func (p *Parser) parseAlterDatabaseSetOption() string {
	if p.cur.Type == tokEOF || p.cur.Type == ';' || p.cur.Type == kwWITH {
		return ""
	}
	if p.isStatementStart() && p.cur.Type != kwWITH && p.cur.Type != kwON {
		return ""
	}

	// Handle ON/OFF as identifier-like for options like ENABLE_BROKER
	if p.cur.Type == kwON {
		// Shouldn't happen at start of an option, but handle just in case
		return ""
	}

	// Validate option name against the known set from SqlScriptDOM.
	if !p.isValidOption(dbSetOptions) {
		return ""
	}

	key := strings.ToUpper(p.cur.Str)
	p.advance()

	// Options that are standalone keywords (no value):
	// ONLINE, OFFLINE, EMERGENCY, SINGLE_USER, MULTI_USER, RESTRICTED_USER,
	// READ_ONLY, READ_WRITE, ENABLE_BROKER, DISABLE_BROKER, NEW_BROKER, ERROR_BROKER_CONVERSATIONS
	switch key {
	case "ONLINE", "OFFLINE", "EMERGENCY",
		"SINGLE_USER", "MULTI_USER", "RESTRICTED_USER",
		"READ_ONLY", "READ_WRITE",
		"ENABLE_BROKER", "DISABLE_BROKER", "NEW_BROKER", "ERROR_BROKER_CONVERSATIONS":
		return key
	}

	// SET PARTNER (database mirroring)
	// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-database-transact-sql-database-mirroring
	//
	//   SET PARTNER { = 'partner_server' | FAILOVER | FORCE_SERVICE_ALLOW_DATA_LOSS
	//     | OFF | RESUME | SAFETY { FULL | OFF } | SUSPEND | TIMEOUT integer }
	if key == "PARTNER" {
		return p.parseAlterDatabaseSetPartner()
	}

	// SET WITNESS (database mirroring)
	//   SET WITNESS { = 'witness_server' | OFF }
	if key == "WITNESS" {
		if p.cur.Type == '=' {
			p.advance()
			val := p.parseAlterDatabaseOptionValue()
			return "WITNESS=" + val
		}
		if p.isAnyKeywordIdent() {
			val := strings.ToUpper(p.cur.Str)
			p.advance()
			return "WITNESS=" + val
		}
		return "WITNESS"
	}

	// SET HADR (Always On)
	// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-database-transact-sql-set-hadr
	//
	//   SET HADR { AVAILABILITY GROUP = group_name | OFF | SUSPEND | RESUME }
	if key == "HADR" {
		return p.parseAlterDatabaseSetHadr()
	}

	// QUERY_STORE special handling: CLEAR [ALL]
	if key == "QUERY_STORE" {
		if p.cur.Type == kwCLEAR {
			p.advance()
			if p.cur.Type == kwALL {
				p.advance()
				return "QUERY_STORE=CLEAR ALL"
			}
			return "QUERY_STORE=CLEAR"
		}
	}

	// Options with = value
	if p.cur.Type == '=' {
		p.advance() // consume =
		val := p.parseAlterDatabaseOptionValue()

		// Handle sub-options in parens: CHANGE_TRACKING = ON (...), QUERY_STORE = ON (...), etc.
		if (key == "CHANGE_TRACKING" || key == "QUERY_STORE" || key == "ACCELERATED_DATABASE_RECOVERY") &&
			val == "ON" && p.cur.Type == '(' {
			subOpts := p.parseAlterDatabaseSubOptions()
			return key + "=ON(" + subOpts + ")"
		}
		// QUERY_STORE = OFF (FORCED)
		if key == "QUERY_STORE" && val == "OFF" && p.cur.Type == '(' {
			subOpts := p.parseAlterDatabaseSubOptions()
			return key + "=OFF(" + subOpts + ")"
		}

		return key + "=" + val
	}

	// Options with value but no =: RECOVERY FULL, PAGE_VERIFY CHECKSUM, etc.
	// PARAMETERIZATION { SIMPLE | FORCED }, ENCRYPTION { ON | OFF | SUSPEND | RESUME }
	// Option NAME ON/OFF patterns: AUTO_CLOSE ON, ANSI_NULLS ON, etc.
	if p.cur.Type == kwON {
		p.advance()
		// Check for sub-options in parens after ON
		if (key == "CHANGE_TRACKING" || key == "QUERY_STORE" || key == "ACCELERATED_DATABASE_RECOVERY") && p.cur.Type == '(' {
			subOpts := p.parseAlterDatabaseSubOptions()
			return key + "=ON(" + subOpts + ")"
		}
		return key + "=ON"
	}
	// If this option has a declared value enum, restrict the value to that set.
	// Anything else is rejected by returning "" (which the outer loop treats
	// as "option not recognized here" — same error surface as an unknown key).
	if enum, ok := dbSetOptionValueEnums[key]; ok {
		if !p.isValidOption(enum) {
			return ""
		}
		val := strings.ToUpper(p.cur.Str)
		p.advance()
		return key + "=" + val
	}

	// Accept any keyword or identifier as option value (AUTO_CLOSE OFF, etc.)
	// TODO(ddl-option-strictness): migrate remaining keys (HADR, ENCRYPTION,
	// ACCELERATED_DATABASE_RECOVERY, ...) into dbSetOptionValueEnums per the
	// follow-up plan in docs/plans/2026-04-22-ddl-option-strictness.md.
	if p.isAnyKeywordIdent() {
		val := strings.ToUpper(p.cur.Str)
		p.advance()
		// Handle TARGET_RECOVERY_TIME = 60 SECONDS (already consumed the number via =)
		// Here it's like "AUTO_CLOSE OFF" or "TARGET_RECOVERY_TIME 60"
		if key == "TARGET_RECOVERY_TIME" {
			// val is the number, check for unit
			if p.isAnyKeywordIdent() {
				unit := strings.ToUpper(p.cur.Str)
				if unit == "SECONDS" || unit == "MINUTES" {
					p.advance()
					return key + "=" + val + " " + unit
				}
			}
			return key + "=" + val
		}
		return key + "=" + val
	}

	// FILESTREAM ( ... )
	if key == "FILESTREAM" && p.cur.Type == '(' {
		subOpts := p.parseAlterDatabaseSubOptions()
		return "FILESTREAM(" + subOpts + ")"
	}

	return key
}

// parseAlterDatabaseSetPartner parses the remainder of SET PARTNER options.
// Caller has consumed PARTNER.
//
//	SET PARTNER { = 'partner_server' | FAILOVER | FORCE_SERVICE_ALLOW_DATA_LOSS
//	  | OFF | RESUME | SAFETY { FULL | OFF } | SUSPEND | TIMEOUT integer }
func (p *Parser) parseAlterDatabaseSetPartner() string {
	if p.cur.Type == '=' {
		p.advance()
		val := p.parseAlterDatabaseOptionValue()
		return "PARTNER=" + val
	}
	if p.isAnyKeywordIdent() {
		sub := strings.ToUpper(p.cur.Str)
		p.advance()
		switch sub {
		case "SAFETY":
			// SAFETY { FULL | OFF }
			if p.isAnyKeywordIdent() {
				val := strings.ToUpper(p.cur.Str)
				p.advance()
				return "PARTNER SAFETY=" + val
			}
			return "PARTNER SAFETY"
		case "TIMEOUT":
			// TIMEOUT integer
			if p.cur.Type == tokICONST {
				val := p.cur.Str
				p.advance()
				return "PARTNER TIMEOUT=" + val
			}
			return "PARTNER TIMEOUT"
		case "OFF":
			return "PARTNER=OFF"
		default:
			// FAILOVER, FORCE_SERVICE_ALLOW_DATA_LOSS, RESUME, SUSPEND
			return "PARTNER " + sub
		}
	}
	return "PARTNER"
}

// parseAlterDatabaseSetHadr parses the remainder of SET HADR options.
// Caller has consumed HADR.
//
//	SET HADR { AVAILABILITY GROUP = group_name | OFF | SUSPEND | RESUME }
func (p *Parser) parseAlterDatabaseSetHadr() string {
	if p.isAnyKeywordIdent() {
		sub := strings.ToUpper(p.cur.Str)
		p.advance()
		switch sub {
		case "AVAILABILITY":
			// AVAILABILITY GROUP = group_name
			if p.cur.Type == kwGROUP {
				p.advance() // consume GROUP
			}
			if p.cur.Type == '=' {
				p.advance()
				val := p.parseAlterDatabaseOptionValue()
				return "HADR AVAILABILITY GROUP=" + val
			}
			return "HADR AVAILABILITY GROUP"
		case "OFF":
			return "HADR=OFF"
		default:
			// SUSPEND, RESUME
			return "HADR " + sub
		}
	}
	if p.cur.Type == kwOFF {
		p.advance()
		return "HADR=OFF"
	}
	return "HADR"
}

// parseAlterDatabaseOptionValue parses a value after = in a SET option.
func (p *Parser) parseAlterDatabaseOptionValue() string {
	if p.cur.Type == tokSCONST {
		val := p.cur.Str
		p.advance()
		return "'" + val + "'"
	}
	if p.cur.Type == tokICONST || p.cur.Type == tokFCONST {
		val := p.cur.Str
		p.advance()
		// Check for unit suffix (SECONDS, MINUTES, DAYS, HOURS)
		if p.isAnyKeywordIdent() {
			unit := strings.ToUpper(p.cur.Str)
			switch unit {
			case "SECONDS", "MINUTES", "DAYS", "HOURS":
				p.advance()
				return val + " " + unit
			}
		}
		return val
	}
	// Accept any keyword or identifier as an option value (ON, OFF, FULL, etc.)
	if p.isAnyKeywordIdent() {
		val := strings.ToUpper(p.cur.Str)
		p.advance()
		return val
	}
	return ""
}

// parseAlterDatabaseSubOptions parses parenthesized sub-options like (AUTO_CLEANUP = ON, CHANGE_RETENTION = 7 DAYS).
func (p *Parser) parseAlterDatabaseSubOptions() string {
	p.advance() // consume (
	var parts []string
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		// key
		key := ""
		if p.isAnyKeywordIdent() {
			key = strings.ToUpper(p.cur.Str)
			p.advance()
		} else if p.cur.Type == tokICONST || p.cur.Type == tokFCONST {
			// Skip bare numerics that aren't part of a key=value
			p.advance()
			continue
		}
		// = value
		if p.cur.Type == '=' {
			p.advance()
			key += "=" + p.parseAlterDatabaseOptionValue()
		} else if p.cur.Type == kwON {
			key += "=ON"
			p.advance()
		} else if p.cur.Type == kwOFF {
			key += "=OFF"
			p.advance()
		}
		if key != "" {
			parts = append(parts, key)
		}
		if p.cur.Type == ',' {
			p.advance()
		} else if p.cur.Type != ')' && p.cur.Type != tokEOF {
			// Skip unrecognized tokens (e.g., size units like GB, MB, TB)
			p.advance()
		}
	}
	p.match(')') // consume )
	return strings.Join(parts, ", ")
}

// parseAlterDatabaseTermination parses WITH <termination>.
//
//	<termination> ::=
//	    ROLLBACK AFTER number [ SECONDS ]
//	  | ROLLBACK IMMEDIATE
//	  | NO_WAIT
func (p *Parser) parseAlterDatabaseTermination() string {
	if p.isAnyKeywordIdent() {
		kw := strings.ToUpper(p.cur.Str)
		switch kw {
		case "ROLLBACK":
			p.advance() // consume ROLLBACK
			if p.cur.Type == kwIMMEDIATE {
				p.advance()
				return "ROLLBACK IMMEDIATE"
			}
			// ROLLBACK AFTER number [SECONDS]
			if p.cur.Type == kwAFTER {
				p.advance() // consume AFTER
				result := "ROLLBACK AFTER"
				if p.cur.Type == tokICONST {
					result += " " + p.cur.Str
					p.advance()
				}
				if p.cur.Type == kwSECONDS {
					result += " SECONDS"
					p.advance()
				}
				return result
			}
			return "ROLLBACK"
		case "NO_WAIT":
			p.advance()
			return "NO_WAIT"
		}
	}
	return ""
}

// parseAlterDatabaseUnknownOption parses a single token or key=value pair from an unknown ALTER DATABASE action.
// Returns the parsed option as a string, or empty string if nothing can be parsed.
func (p *Parser) parseAlterDatabaseUnknownOption() string {
	if p.cur.Type == tokEOF || p.cur.Type == ';' || p.isStatementStart() {
		return ""
	}

	// Parenthesized list of key=value pairs
	if p.cur.Type == '(' {
		p.advance() // consume (
		var parts []string
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			if p.cur.Type == ',' {
				p.advance()
				continue
			}
			key := p.parseAlterDatabaseUnknownOptionValue()
			if key == "" {
				break
			}
			// Check for = value
			if p.cur.Type == '=' {
				p.advance()
				val := p.parseAlterDatabaseUnknownOptionValue()
				parts = append(parts, key+"="+val)
			} else {
				parts = append(parts, key)
			}
		}
		p.match(')') // consume )
		return "(" + strings.Join(parts, ", ") + ")"
	}

	// Keyword or identifier
	if p.isAnyKeywordIdent() {
		key := strings.ToUpper(p.cur.Str)
		p.advance()
		// Check for = value
		if p.cur.Type == '=' {
			p.advance()
			val := p.parseAlterDatabaseUnknownOptionValue()
			return key + "=" + val
		}
		return key
	}

	// Numeric or string constant
	if p.cur.Type == tokICONST || p.cur.Type == tokFCONST {
		val := p.cur.Str
		p.advance()
		return val
	}
	if p.cur.Type == tokSCONST {
		val := "'" + p.cur.Str + "'"
		p.advance()
		return val
	}

	return ""
}

// parseAlterDatabaseUnknownOptionValue parses a single value in an unknown option context.
func (p *Parser) parseAlterDatabaseUnknownOptionValue() string {
	if p.cur.Type == tokICONST || p.cur.Type == tokFCONST {
		val := p.cur.Str
		p.advance()
		return val
	}
	if p.cur.Type == tokSCONST {
		val := "'" + p.cur.Str + "'"
		p.advance()
		return val
	}
	// Accept any keyword or identifier as an option value (ON, OFF, FULL, etc.)
	if p.isAnyKeywordIdent() {
		val := strings.ToUpper(p.cur.Str)
		p.advance()
		return val
	}
	return ""
}

// parseAlterDatabaseAdd parses ADD FILE/LOG FILE/FILEGROUP.
func (p *Parser) parseAlterDatabaseAdd(stmt *nodes.AlterDatabaseStmt) {
	if p.cur.Type == kwFILE {
		p.advance() // consume FILE
		stmt.SubAction = "FILE"
		// Parse file specs
		var files []nodes.Node
		if p.cur.Type == '(' {
			files = append(files, p.parseDatabaseFileSpec())
			for p.cur.Type == ',' {
				next := p.peekNext()
				if next.Type == '(' {
					p.advance() // consume comma
					files = append(files, p.parseDatabaseFileSpec())
				} else {
					break
				}
			}
		}
		stmt.FileSpecs = &nodes.List{Items: files}
		// TO FILEGROUP filegroup_name
		if p.cur.Type == kwTO {
			p.advance() // consume TO
			if p.cur.Type == kwFILEGROUP {
				p.advance() // consume FILEGROUP
				if id, ok := p.parseIdentifier(); ok {
					stmt.TargetName = id
				}
			}
		}
	} else if p.cur.Type == kwLOG {
		p.advance() // consume LOG
		if p.cur.Type == kwFILE {
			p.advance() // consume FILE
		}
		stmt.SubAction = "LOG FILE"
		var files []nodes.Node
		if p.cur.Type == '(' {
			files = append(files, p.parseDatabaseFileSpec())
			for p.cur.Type == ',' {
				next := p.peekNext()
				if next.Type == '(' {
					p.advance() // consume comma
					files = append(files, p.parseDatabaseFileSpec())
				} else {
					break
				}
			}
		}
		stmt.FileSpecs = &nodes.List{Items: files}
	} else if p.cur.Type == kwFILEGROUP {
		p.advance() // consume FILEGROUP
		stmt.SubAction = "FILEGROUP"
		if id, ok := p.parseIdentifier(); ok {
			stmt.TargetName = id
		}
		// CONTAINS FILESTREAM | CONTAINS MEMORY_OPTIMIZED_DATA
		if p.cur.Type == kwCONTAINS {
			p.advance() // consume CONTAINS
			if p.cur.Type == kwFILESTREAM {
				stmt.Options = &nodes.List{Items: []nodes.Node{&nodes.String{Str: "CONTAINS FILESTREAM"}}}
				p.advance()
			} else if p.cur.Type == kwMEMORY_OPTIMIZED_DATA {
				stmt.Options = &nodes.List{Items: []nodes.Node{&nodes.String{Str: "CONTAINS MEMORY_OPTIMIZED_DATA"}}}
				p.advance()
			}
		}
	}
}

// parseAlterDatabaseModify parses MODIFY FILE/NAME/FILEGROUP or azure options.
func (p *Parser) parseAlterDatabaseModify(stmt *nodes.AlterDatabaseStmt) {
	if p.cur.Type == kwFILE {
		p.advance() // consume FILE
		stmt.SubAction = "FILE"
		if p.cur.Type == '(' {
			stmt.FileSpecs = &nodes.List{Items: []nodes.Node{p.parseDatabaseFileSpec()}}
		}
	} else if p.cur.Type == kwNAME {
		p.advance() // consume NAME
		stmt.SubAction = "NAME"
		if p.cur.Type == '=' {
			p.advance() // consume =
			if id, ok := p.parseIdentifier(); ok {
				stmt.NewName = id
			}
		}
	} else if p.cur.Type == kwFILEGROUP {
		p.advance() // consume FILEGROUP
		stmt.SubAction = "FILEGROUP"
		// filegroup_name
		if id, ok := p.parseIdentifier(); ok {
			stmt.TargetName = id
		}
		// { <filegroup_updatability_option> | DEFAULT | NAME = new_name | AUTOGROW_SINGLE_FILE | AUTOGROW_ALL_FILES }
		p.parseAlterDatabaseModifyFilegroupOption(stmt)
	} else if p.cur.Type == '(' {
		// Azure options: MODIFY ( option = value [, ...] ) [ WITH MANUAL_CUTOVER ]
		stmt.SubAction = "AZURE_OPTIONS"
		subOpts := p.parseAlterDatabaseSubOptions()
		stmt.Options = &nodes.List{Items: []nodes.Node{&nodes.String{Str: subOpts}}}
		// [ WITH MANUAL_CUTOVER ]
		if p.cur.Type == kwWITH {
			p.advance() // consume WITH
			if p.cur.Type == kwMANUAL_CUTOVER {
				stmt.Termination = "MANUAL_CUTOVER"
				p.advance()
			}
		}
	}
}

// parseAlterDatabaseModifyFilegroupOption parses the option after MODIFY FILEGROUP fg_name.
func (p *Parser) parseAlterDatabaseModifyFilegroupOption(stmt *nodes.AlterDatabaseStmt) {
	if p.cur.Type == kwDEFAULT {
		p.advance()
		stmt.Options = &nodes.List{Items: []nodes.Node{&nodes.String{Str: "DEFAULT"}}}
		return
	}
	// NAME = new_name
	if p.cur.Type == kwNAME {
		p.advance() // consume NAME
		if p.cur.Type == '=' {
			p.advance() // consume =
			if id, ok := p.parseIdentifier(); ok {
				stmt.NewName = id
			}
		}
		return
	}
	if p.isAnyKeywordIdent() {
		opt := strings.ToUpper(p.cur.Str)
		switch opt {
		case "READ_ONLY", "READ_WRITE", "READONLY", "READWRITE",
			"AUTOGROW_SINGLE_FILE", "AUTOGROW_ALL_FILES":
			p.advance()
			stmt.Options = &nodes.List{Items: []nodes.Node{&nodes.String{Str: opt}}}
		}
	}
}

// parseAlterDatabaseRemove parses REMOVE FILE/FILEGROUP.
func (p *Parser) parseAlterDatabaseRemove(stmt *nodes.AlterDatabaseStmt) {
	if p.cur.Type == kwFILE {
		p.advance() // consume FILE
		stmt.SubAction = "FILE"
		if id, ok := p.parseIdentifier(); ok {
			stmt.TargetName = id
		}
	} else if p.cur.Type == kwFILEGROUP {
		p.advance() // consume FILEGROUP
		stmt.SubAction = "FILEGROUP"
		if id, ok := p.parseIdentifier(); ok {
			stmt.TargetName = id
		}
	}
}

// parseAlterIndexStmt parses an ALTER INDEX statement.
//
// BNF: mssql/parser/bnf/alter-index-transact-sql.bnf
//
//	ALTER INDEX { index_name | ALL } ON <object>
//	{
//	      REBUILD [
//	            [ WITH ( <rebuild_index_option> [ ,...n ] ) ]
//	          | [ PARTITION = ALL [ WITH ( <rebuild_index_option> [ ,...n ] ) ] ]
//	          | [ PARTITION = partition_number [ WITH ( <single_partition_rebuild_index_option> ) ] ]
//	      ]
//	    | DISABLE
//	    | REORGANIZE  [ PARTITION = partition_number ] [ WITH ( <reorganize_option> ) ]
//	    | SET ( <set_index_option> [ ,...n ] )
//	    | RESUME [ WITH ( <resumable_index_option> [ ,...n ] ) ]
//	    | PAUSE
//	    | ABORT
//	}
func (p *Parser) parseAlterIndexStmt() (*nodes.AlterIndexStmt, error) {
	loc := p.pos()

	stmt := &nodes.AlterIndexStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Index name or ALL. ALL is a Core keyword with special meaning
	// ("every index on the table") in this position per SqlScriptDOM.
	if p.isIdentLike() || p.cur.Type == kwALL {
		stmt.IndexName = p.cur.Str
		p.advance()
	}

	// ON table_name
	if p.cur.Type == kwON {
		p.advance() // consume ON
		var err error
		stmt.Table, err = p.parseTableRef()
		if err != nil {
			return nil, err
		}
	}

	// Action: REBUILD, REORGANIZE, DISABLE, SET, RESUME, PAUSE, ABORT
	if p.cur.Type == kwSET {
		stmt.Action = "SET"
		p.advance()
	} else if p.isAnyKeywordIdent() {
		stmt.Action = strings.ToUpper(p.cur.Str)
		p.advance()
	}

	// Parse PARTITION = { partition_number | ALL }
	if p.cur.Type == kwPARTITION {
		p.advance() // consume PARTITION
		p.match('=')
		if p.cur.Type == kwALL {
			stmt.Partition = "ALL"
			p.advance()
		} else {
			stmt.Partition = p.cur.Str
			p.advance()
		}
	}

	// Parse WITH ( option = value [, ...] ) for REBUILD/REORGANIZE/RESUME or SET ( ... )
	if p.cur.Type == kwWITH || stmt.Action == "SET" {
		if p.cur.Type == kwWITH {
			p.advance() // consume WITH
		}
		if p.cur.Type == '(' {
			var err error
			stmt.Options, err = p.parseAlterIndexOptions()
			if err != nil {
				return nil, err
			}
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// indexOptions defines the valid option names for index WITH clauses.
// Derived from SqlScriptDOM IndexOptionHelper + IndexStateOption + spatial/JSON/vector options.
var indexOptions = newOptionSet(
	kwFILLFACTOR,     // FILLFACTOR = N
	kwMAXDOP,         // MAXDOP = N
	kwBUCKET_COUNT,   // BUCKET_COUNT = N (hash indexes)
	kwLOB_COMPACTION, // LOB_COMPACTION = ON|OFF
	kwORDER,          // ORDER (columnstore)
	kwTYPE,           // TYPE (vector indexes)
).withIdents(
	// indexStateOption — on/off options parsed via IndexOptionHelper
	"PAD_INDEX",
	"SORT_IN_TEMPDB",
	"STATISTICS_NORECOMPUTE",
	"STATISTICS_INCREMENTAL",
	"DROP_EXISTING",
	"ALLOW_ROW_LOCKS",
	"ALLOW_PAGE_LOCKS",
	"IGNORE_DUP_KEY",
	"ONLINE",
	"RESUMABLE",
	"COMPRESS_ALL_ROW_GROUPS",
	"OPTIMIZE_FOR_SEQUENTIAL_KEY",
	// expression-valued options
	"MAX_DURATION",
	"DATA_COMPRESSION",
	"XML_COMPRESSION",
	"COMPRESSION_DELAY",
	"WAIT_AT_LOW_PRIORITY",
	// spatial index options
	"BOUNDING_BOX",
	"GRIDS",
	"CELLS_PER_OBJECT",
	// JSON index options
	"OPTIMIZE_FOR_ARRAY_SEARCH",
	// vector index options
	"METRIC",
)

// parseAlterIndexOptions parses a parenthesized list of index options.
//
// <rebuild_index_option> / <reorganize_option> / <set_index_option> ::=
//
//	PAD_INDEX = { ON | OFF }
//	| FILLFACTOR = fillfactor
//	| SORT_IN_TEMPDB = { ON | OFF }
//	| IGNORE_DUP_KEY = { ON | OFF }
//	| STATISTICS_NORECOMPUTE = { ON | OFF }
//	| STATISTICS_INCREMENTAL = { ON | OFF }
//	| ONLINE = { ON [ ( <low_priority_lock_wait> ) ] | OFF }
//	| RESUMABLE = { ON | OFF }
//	| MAX_DURATION = <time> [ MINUTES ]
//	| ALLOW_ROW_LOCKS = { ON | OFF }
//	| ALLOW_PAGE_LOCKS = { ON | OFF }
//	| MAXDOP = max_degree_of_parallelism
//	| DATA_COMPRESSION = { NONE | ROW | PAGE | COLUMNSTORE | COLUMNSTORE_ARCHIVE }
//	    [ ON PARTITIONS ( { <partition_number> [ TO <partition_number> ] } [ , ...n ] ) ]
//	| XML_COMPRESSION = { ON | OFF }
//	    [ ON PARTITIONS ( { <partition_number> [ TO <partition_number> ] } [ , ...n ] ) ]
//	| LOB_COMPACTION = { ON | OFF }
//	| COMPRESS_ALL_ROW_GROUPS = { ON | OFF }
func (p *Parser) parseAlterIndexOptions() (*nodes.List, error) {
	p.advance() // consume '('

	opts := &nodes.List{}
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		if p.cur.Type == ',' {
			p.advance()
			continue
		}

		// Parse option_name = value (validated against known index options)
		if p.isValidOption(indexOptions) {
			name := strings.ToUpper(p.cur.Str)
			p.advance()

			if p.cur.Type == '=' {
				p.advance() // consume '='
				// Value can be ON, OFF, an identifier, or a number
				val := ""
				if p.cur.Type == kwON {
					val = "ON"
					p.advance()
					// ONLINE = ON ( <low_priority_lock_wait> ) -- skip nested parens
					if p.cur.Type == '(' {
						depth := 1
						p.advance() // consume '('
						for depth > 0 && p.cur.Type != tokEOF {
							if p.cur.Type == '(' {
								depth++
							} else if p.cur.Type == ')' {
								depth--
								if depth == 0 {
									p.advance()
									break
								}
							}
							p.advance()
						}
					}
				} else if p.cur.Type == kwOFF {
					val = "OFF"
					p.advance()
				} else if p.cur.Type == tokICONST || p.cur.Type == tokFCONST {
					val = p.cur.Str
					p.advance()
				} else if p.isAnyKeywordIdent() {
					val = strings.ToUpper(p.cur.Str)
					p.advance()
				} else if p.cur.Type == '(' {
					// Nested parenthesized value like BOUNDING_BOX = (0, 0, 100, 100)
					depth := 1
					p.advance() // consume '('
					for depth > 0 && p.cur.Type != tokEOF {
						if p.cur.Type == '(' {
							depth++
						} else if p.cur.Type == ')' {
							depth--
							if depth == 0 {
								p.advance()
								break
							}
						}
						p.advance()
					}
				}
				opts.Items = append(opts.Items, &nodes.String{Str: name + "=" + val})
			} else {
				opts.Items = append(opts.Items, &nodes.String{Str: name})
			}

			// DATA_COMPRESSION / XML_COMPRESSION may have ON PARTITIONS (...)
			if p.cur.Type == kwON {
				next := p.peekNext()
				if next.Type == kwPARTITIONS {
					p.advance() // consume ON
					p.advance() // consume PARTITIONS
					if p.cur.Type == '(' {
						depth := 1
						p.advance() // consume '('
						for depth > 0 && p.cur.Type != tokEOF {
							if p.cur.Type == '(' {
								depth++
							} else if p.cur.Type == ')' {
								depth--
								if depth == 0 {
									p.advance()
									break
								}
							}
							p.advance()
						}
					}
				}
			}

			// MAX_DURATION may have MINUTES suffix
			if strings.HasPrefix(name, "MAX_DURATION") && p.cur.Type == kwMINUTES {
				p.advance()
			}
		} else if p.cur.Type == tokSCONST || p.cur.Type == tokICONST || p.cur.Type == tokFCONST {
			// Collect string/numeric literals as standalone option values
			opts.Items = append(opts.Items, &nodes.String{Str: p.cur.Str})
			p.advance()
		} else {
			// Unknown token in option list — reject
			return nil, p.unexpectedToken()
		}
	}

	p.match(')') // consume ')'
	return opts, nil
}
