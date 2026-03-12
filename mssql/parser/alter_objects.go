// Package parser - alter_objects.go implements T-SQL ALTER DATABASE and ALTER INDEX parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseAlterDatabaseStmt parses an ALTER DATABASE statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-database-transact-sql
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-database-transact-sql-set-options
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-database-transact-sql-file-and-filegroup-options
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
//	  | COLLATE collation_name
//	  | MODIFY NAME = new_database_name
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
func (p *Parser) parseAlterDatabaseStmt() *nodes.AlterDatabaseStmt {
	loc := p.pos()

	stmt := &nodes.AlterDatabaseStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Database name or CURRENT
	if p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// Dispatch on action keyword
	if p.cur.Type == kwSET {
		p.advance() // consume SET
		stmt.Action = "SET"
		p.parseAlterDatabaseSetOptions(stmt)
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
	} else if p.isIdentLike() {
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
		default:
			// Unknown action - record and skip
			stmt.Action = action
			p.advance()
			for p.cur.Type != tokEOF && p.cur.Type != ';' && !p.isStatementStart() {
				p.advance()
			}
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseAlterDatabaseSetOptions parses SET <option_spec> [,...n] [WITH <termination>].
func (p *Parser) parseAlterDatabaseSetOptions(stmt *nodes.AlterDatabaseStmt) {
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

	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}

	// WITH <termination>
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		stmt.Termination = p.parseAlterDatabaseTermination()
	}
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

	if !p.isIdentLike() && p.cur.Type != kwSET {
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

	// QUERY_STORE special handling: CLEAR [ALL]
	if key == "QUERY_STORE" {
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "CLEAR") {
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
			strings.EqualFold(val, "ON") && p.cur.Type == '(' {
			subOpts := p.parseAlterDatabaseSubOptions()
			return key + "=ON(" + subOpts + ")"
		}
		// QUERY_STORE = OFF (FORCED)
		if key == "QUERY_STORE" && strings.EqualFold(val, "OFF") && p.cur.Type == '(' {
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
	if p.isIdentLike() {
		val := strings.ToUpper(p.cur.Str)
		p.advance()
		// Handle TARGET_RECOVERY_TIME = 60 SECONDS (already consumed the number via =)
		// Here it's like "RECOVERY FULL" or "AUTO_CLOSE OFF" or "TARGET_RECOVERY_TIME 60"
		if key == "TARGET_RECOVERY_TIME" {
			// val is the number, check for unit
			if p.isIdentLike() {
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
		if p.isIdentLike() {
			unit := strings.ToUpper(p.cur.Str)
			switch unit {
			case "SECONDS", "MINUTES", "DAYS", "HOURS":
				p.advance()
				return val + " " + unit
			}
		}
		return val
	}
	if p.cur.Type == kwON {
		p.advance()
		return "ON"
	}
	if p.cur.Type == kwOFF {
		p.advance()
		return "OFF"
	}
	if p.isIdentLike() {
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
		var sb strings.Builder
		// key
		if p.isIdentLike() {
			sb.WriteString(strings.ToUpper(p.cur.Str))
			p.advance()
		}
		// = value
		if p.cur.Type == '=' {
			sb.WriteString("=")
			p.advance()
			sb.WriteString(p.parseAlterDatabaseOptionValue())
		} else if p.cur.Type == kwON || (p.isIdentLike() && matchesKeywordCI(p.cur.Str, "ON")) {
			sb.WriteString("=ON")
			p.advance()
		} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "OFF") {
			sb.WriteString("=OFF")
			p.advance()
		}
		if sb.Len() > 0 {
			parts = append(parts, sb.String())
		}
		if p.cur.Type == ',' {
			p.advance()
		} else {
			break
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
	if p.isIdentLike() {
		kw := strings.ToUpper(p.cur.Str)
		switch kw {
		case "ROLLBACK":
			p.advance() // consume ROLLBACK
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "IMMEDIATE") {
				p.advance()
				return "ROLLBACK IMMEDIATE"
			}
			// ROLLBACK AFTER number [SECONDS]
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "AFTER") {
				p.advance() // consume AFTER
				var sb strings.Builder
				sb.WriteString("ROLLBACK AFTER ")
				if p.cur.Type == tokICONST {
					sb.WriteString(p.cur.Str)
					p.advance()
				}
				if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SECONDS") {
					sb.WriteString(" SECONDS")
					p.advance()
				}
				return sb.String()
			}
			return "ROLLBACK"
		case "NO_WAIT":
			p.advance()
			return "NO_WAIT"
		}
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
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "FILEGROUP") {
				p.advance() // consume FILEGROUP
				if id, ok := p.parseIdentifier(); ok {
					stmt.TargetName = id
				}
			}
		}
	} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "LOG") {
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
	} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "FILEGROUP") {
		p.advance() // consume FILEGROUP
		stmt.SubAction = "FILEGROUP"
		if id, ok := p.parseIdentifier(); ok {
			stmt.TargetName = id
		}
		// CONTAINS FILESTREAM | CONTAINS MEMORY_OPTIMIZED_DATA
		if p.cur.Type == kwCONTAINS {
			p.advance() // consume CONTAINS
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "FILESTREAM") {
				stmt.Options = &nodes.List{Items: []nodes.Node{&nodes.String{Str: "CONTAINS FILESTREAM"}}}
				p.advance()
			} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "MEMORY_OPTIMIZED_DATA") {
				stmt.Options = &nodes.List{Items: []nodes.Node{&nodes.String{Str: "CONTAINS MEMORY_OPTIMIZED_DATA"}}}
				p.advance()
			}
		}
	}
}

// parseAlterDatabaseModify parses MODIFY FILE/NAME/FILEGROUP.
func (p *Parser) parseAlterDatabaseModify(stmt *nodes.AlterDatabaseStmt) {
	if p.cur.Type == kwFILE {
		p.advance() // consume FILE
		stmt.SubAction = "FILE"
		if p.cur.Type == '(' {
			stmt.FileSpecs = &nodes.List{Items: []nodes.Node{p.parseDatabaseFileSpec()}}
		}
	} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "NAME") {
		p.advance() // consume NAME
		stmt.SubAction = "NAME"
		if p.cur.Type == '=' {
			p.advance() // consume =
			if id, ok := p.parseIdentifier(); ok {
				stmt.NewName = id
			}
		}
	} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "FILEGROUP") {
		p.advance() // consume FILEGROUP
		stmt.SubAction = "FILEGROUP"
		// filegroup_name
		if id, ok := p.parseIdentifier(); ok {
			stmt.TargetName = id
		}
		// { <filegroup_updatability_option> | DEFAULT | NAME = new_name | AUTOGROW_SINGLE_FILE | AUTOGROW_ALL_FILES }
		p.parseAlterDatabaseModifyFilegroupOption(stmt)
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
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "NAME") {
		p.advance() // consume NAME
		if p.cur.Type == '=' {
			p.advance() // consume =
			if id, ok := p.parseIdentifier(); ok {
				stmt.NewName = id
			}
		}
		return
	}
	if p.isIdentLike() {
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
	} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "FILEGROUP") {
		p.advance() // consume FILEGROUP
		stmt.SubAction = "FILEGROUP"
		if id, ok := p.parseIdentifier(); ok {
			stmt.TargetName = id
		}
	}
}

// parseAlterIndexStmt parses an ALTER INDEX statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-index-transact-sql
//
//	ALTER INDEX { index_name | ALL } ON <object>
//	    { REBUILD [ PARTITION = ... ] | REORGANIZE [ PARTITION = ... ] | DISABLE | SET ( ... ) }
func (p *Parser) parseAlterIndexStmt() *nodes.AlterIndexStmt {
	loc := p.pos()

	stmt := &nodes.AlterIndexStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Index name or ALL
	if p.isIdentLike() {
		stmt.IndexName = p.cur.Str
		p.advance()
	}

	// ON table_name
	if p.cur.Type == kwON {
		p.advance() // consume ON
		stmt.Table = p.parseTableRef()
	}

	// Action: REBUILD, REORGANIZE, DISABLE, SET
	if p.cur.Type == kwSET {
		stmt.Action = "SET"
		p.advance()
	} else if p.isIdentLike() {
		stmt.Action = strings.ToUpper(p.cur.Str)
		p.advance()
	}

	// Parse PARTITION = { partition_number | ALL }
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "PARTITION") {
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

	// Parse WITH ( option = value [, ...] ) or SET ( option = value [, ...] )
	if p.cur.Type == kwWITH || stmt.Action == "SET" {
		if p.cur.Type == kwWITH {
			p.advance() // consume WITH
		}
		if p.cur.Type == '(' {
			stmt.Options = p.parseAlterIndexOptions()
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

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
func (p *Parser) parseAlterIndexOptions() *nodes.List {
	p.advance() // consume '('

	opts := &nodes.List{}
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		if p.cur.Type == ',' {
			p.advance()
			continue
		}

		// Parse option_name = value
		if p.isIdentLike() || p.cur.Type == kwON || p.cur.Type == kwOFF {
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
				} else if p.isIdentLike() {
					val = strings.ToUpper(p.cur.Str)
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.String{Str: name + "=" + val})
			} else {
				opts.Items = append(opts.Items, &nodes.String{Str: name})
			}

			// DATA_COMPRESSION / XML_COMPRESSION may have ON PARTITIONS (...)
			if p.cur.Type == kwON {
				next := p.peekNext()
				if next.Str != "" && matchesKeywordCI(next.Str, "PARTITIONS") {
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
			if strings.HasPrefix(name, "MAX_DURATION") && p.isIdentLike() && matchesKeywordCI(p.cur.Str, "MINUTES") {
				p.advance()
			}
		} else {
			// Unknown token, skip
			p.advance()
		}
	}

	p.match(')') // consume ')'
	return opts
}
