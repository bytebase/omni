package parser

import (
	"strconv"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// parseChangeReplicationSourceStmtInner parses a CHANGE REPLICATION SOURCE TO statement.
// CHANGE REPLICATION have already been consumed; p.cur is SOURCE.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/change-replication-source-to.html
//
//	CHANGE REPLICATION SOURCE TO option [, option] ... [ channel_option ]
//
//	option: {
//	    SOURCE_BIND = 'interface_name'
//	  | SOURCE_HOST = 'host_name'
//	  | SOURCE_USER = 'user_name'
//	  | SOURCE_PASSWORD = 'password'
//	  | SOURCE_PORT = port_num
//	  | PRIVILEGE_CHECKS_USER = {NULL | 'account'}
//	  | REQUIRE_ROW_FORMAT = {0|1}
//	  | REQUIRE_TABLE_PRIMARY_KEY_CHECK = {STREAM | ON | OFF | GENERATE}
//	  | ASSIGN_GTIDS_TO_ANONYMOUS_TRANSACTIONS = {OFF | LOCAL | uuid}
//	  | SOURCE_LOG_FILE = 'source_log_name'
//	  | SOURCE_LOG_POS = source_log_pos
//	  | SOURCE_AUTO_POSITION = {0|1}
//	  | RELAY_LOG_FILE = 'relay_log_name'
//	  | RELAY_LOG_POS = relay_log_pos
//	  | SOURCE_HEARTBEAT_PERIOD = interval
//	  | SOURCE_CONNECT_RETRY = interval
//	  | SOURCE_RETRY_COUNT = count
//	  | SOURCE_CONNECTION_AUTO_FAILOVER = {0|1}
//	  | SOURCE_DELAY = interval
//	  | SOURCE_COMPRESSION_ALGORITHMS = 'algorithm[,algorithm][,algorithm]'
//	  | SOURCE_ZSTD_COMPRESSION_LEVEL = level
//	  | SOURCE_SSL = {0|1}
//	  | SOURCE_SSL_CA = 'ca_file_name'
//	  | SOURCE_SSL_CAPATH = 'ca_directory_name'
//	  | SOURCE_SSL_CERT = 'cert_file_name'
//	  | SOURCE_SSL_CRL = 'crl_file_name'
//	  | SOURCE_SSL_CRLPATH = 'crl_directory_name'
//	  | SOURCE_SSL_KEY = 'key_file_name'
//	  | SOURCE_SSL_CIPHER = 'cipher_list'
//	  | SOURCE_SSL_VERIFY_SERVER_CERT = {0|1}
//	  | SOURCE_TLS_VERSION = 'protocol_list'
//	  | SOURCE_TLS_CIPHERSUITES = 'ciphersuite_list'
//	  | SOURCE_PUBLIC_KEY_PATH = 'key_file_name'
//	  | GET_SOURCE_PUBLIC_KEY = {0|1}
//	  | NETWORK_NAMESPACE = 'namespace'
//	  | IGNORE_SERVER_IDS = (server_id_list)
//	  | GTID_ONLY = {0|1}
//	}
//
//	channel_option:
//	    FOR CHANNEL channel
//
//	server_id_list:
//	    [server_id [, server_id] ... ]
func (p *Parser) parseChangeReplicationSourceStmtInner(start int) (*nodes.ChangeReplicationSourceStmt, error) {
	// consume SOURCE
	p.advance()
	// consume TO
	p.match(kwTO)

	stmt := &nodes.ChangeReplicationSourceStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Parse options: option [, option] ...
	for {
		opt, err := p.parseReplicationSourceOption()
		if err != nil {
			return nil, err
		}
		stmt.Options = append(stmt.Options, opt)

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}

	// Optional: FOR CHANNEL channel
	if p.cur.Type == kwFOR {
		p.advance() // consume FOR
		p.advance() // consume CHANNEL (keyword token)
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Channel = name
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseReplicationSourceOption parses a single option in CHANGE REPLICATION SOURCE TO.
func (p *Parser) parseReplicationSourceOption() (*nodes.ReplicationOption, error) {
	start := p.pos()
	opt := &nodes.ReplicationOption{Loc: nodes.Loc{Start: start}}

	// The option name is an identifier (e.g., SOURCE_HOST, SOURCE_PORT, etc.)
	// Must accept any keyword since option names may collide with registered keywords.
	optTokType := p.cur.Type
	name, _, err := p.parseKeywordOrIdent()
	if err != nil {
		return nil, err
	}
	opt.Name = name

	// consume '='
	p.match('=')

	// Special case: IGNORE_SERVER_IDS = (id_list)
	if optTokType == kwIGNORE_SERVER_IDS {
		p.match('(')
		if p.cur.Type != ')' {
			for {
				if p.cur.Type != tokICONST {
					return nil, &ParseError{Message: "expected integer in server ID list", Position: p.cur.Loc}
				}
				id, err := strconv.ParseInt(p.cur.Str, 10, 64)
				if err != nil {
					return nil, &ParseError{Message: "invalid integer in server ID list", Position: p.cur.Loc}
				}
				opt.IDs = append(opt.IDs, id)
				p.advance() // consume the number
				if p.cur.Type != ',' {
					break
				}
				p.advance() // consume ','
			}
		}
		p.match(')')
		opt.Loc.End = p.pos()
		return opt, nil
	}

	// For all other options, consume the value token (string, number, identifier, or NULL)
	switch p.cur.Type {
	case tokSCONST:
		opt.Value = p.cur.Str
		p.advance()
	case tokICONST, tokFCONST:
		opt.Value = p.cur.Str
		p.advance()
	case kwNULL:
		opt.Value = "NULL"
		p.advance()
	default:
		// Identifier or keyword value (e.g., STREAM, ON, OFF, GENERATE, LOCAL).
		// Must accept reserved keywords like ON/OFF in this context.
		val, _, err := p.parseKeywordOrIdent()
		if err != nil {
			return nil, err
		}
		opt.Value = val
	}

	opt.Loc.End = p.pos()
	return opt, nil
}

// parseChangeMasterStmtInner parses a CHANGE MASTER TO statement (legacy alias).
// CHANGE MASTER have already been consumed; p.cur is TO.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/change-master-to.html
//
//	CHANGE MASTER TO option [, option] ... [ channel_option ]
//
//	option: {
//	    MASTER_BIND = 'interface_name'
//	  | MASTER_HOST = 'host_name'
//	  | MASTER_USER = 'user_name'
//	  | MASTER_PASSWORD = 'password'
//	  | MASTER_PORT = port_num
//	  | PRIVILEGE_CHECKS_USER = {NULL | 'account'}
//	  | REQUIRE_ROW_FORMAT = {0|1}
//	  | REQUIRE_TABLE_PRIMARY_KEY_CHECK = {STREAM | ON | OFF | GENERATE}
//	  | ASSIGN_GTIDS_TO_ANONYMOUS_TRANSACTIONS = {OFF | LOCAL | uuid}
//	  | MASTER_LOG_FILE = 'source_log_name'
//	  | MASTER_LOG_POS = source_log_pos
//	  | MASTER_AUTO_POSITION = {0|1}
//	  | RELAY_LOG_FILE = 'relay_log_name'
//	  | RELAY_LOG_POS = relay_log_pos
//	  | MASTER_HEARTBEAT_PERIOD = interval
//	  | MASTER_CONNECT_RETRY = interval
//	  | MASTER_RETRY_COUNT = count
//	  | SOURCE_CONNECTION_AUTO_FAILOVER = {0|1}
//	  | MASTER_DELAY = interval
//	  | MASTER_COMPRESSION_ALGORITHMS = 'algorithm[,algorithm][,algorithm]'
//	  | MASTER_ZSTD_COMPRESSION_LEVEL = level
//	  | MASTER_SSL = {0|1}
//	  | MASTER_SSL_CA = 'ca_file_name'
//	  | MASTER_SSL_CAPATH = 'ca_directory_name'
//	  | MASTER_SSL_CERT = 'cert_file_name'
//	  | MASTER_SSL_CRL = 'crl_file_name'
//	  | MASTER_SSL_CRLPATH = 'crl_directory_name'
//	  | MASTER_SSL_KEY = 'key_file_name'
//	  | MASTER_SSL_CIPHER = 'cipher_list'
//	  | MASTER_SSL_VERIFY_SERVER_CERT = {0|1}
//	  | MASTER_TLS_VERSION = 'protocol_list'
//	  | MASTER_TLS_CIPHERSUITES = 'ciphersuite_list'
//	  | MASTER_PUBLIC_KEY_PATH = 'key_file_name'
//	  | GET_MASTER_PUBLIC_KEY = {0|1}
//	  | NETWORK_NAMESPACE = 'namespace'
//	  | IGNORE_SERVER_IDS = (server_id_list)
//	  | GTID_ONLY = {0|1}
//	}
//
//	channel_option:
//	    FOR CHANNEL channel
//
//	server_id_list:
//	    [server_id [, server_id] ... ]
func (p *Parser) parseChangeMasterStmtInner(start int) (*nodes.ChangeReplicationSourceStmt, error) {
	// consume TO
	p.match(kwTO)

	stmt := &nodes.ChangeReplicationSourceStmt{
		Loc:    nodes.Loc{Start: start},
		Legacy: true,
	}

	// Parse options: option [, option] ...
	for {
		opt, err := p.parseReplicationSourceOption()
		if err != nil {
			return nil, err
		}
		stmt.Options = append(stmt.Options, opt)

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}

	// Optional: FOR CHANNEL channel
	if p.cur.Type == kwFOR {
		p.advance() // consume FOR
		p.advance() // consume CHANNEL (keyword token)
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Channel = name
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseChangeReplicationFilterStmtInner parses a CHANGE REPLICATION FILTER statement.
// CHANGE REPLICATION have already been consumed; p.cur is FILTER.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/change-replication-filter.html
//
//	CHANGE REPLICATION FILTER filter[, filter] [FOR CHANNEL channel]
//
//	filter: {
//	    REPLICATE_DO_DB = (db_list)
//	  | REPLICATE_IGNORE_DB = (db_list)
//	  | REPLICATE_DO_TABLE = (tbl_list)
//	  | REPLICATE_IGNORE_TABLE = (tbl_list)
//	  | REPLICATE_WILD_DO_TABLE = (wild_tbl_list)
//	  | REPLICATE_WILD_IGNORE_TABLE = (wild_tbl_list)
//	  | REPLICATE_REWRITE_DB = (db_pair_list)
//	}
//
//	db_list:
//	    db_name[, db_name][, ...]
//
//	tbl_list:
//	    db_name.table_name[, db_name.table_name][, ...]
//
//	wild_tbl_list:
//	    'db_pattern.table_pattern'[, 'db_pattern.table_pattern'][, ...]
//
//	db_pair_list:
//	    (db_pair)[, (db_pair)][, ...]
//
//	db_pair:
//	    from_db, to_db
func (p *Parser) parseChangeReplicationFilterStmtInner(start int) (*nodes.ChangeReplicationFilterStmt, error) {
	// consume FILTER
	p.advance()

	stmt := &nodes.ChangeReplicationFilterStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Parse filters: filter[, filter]
	for {
		f, err := p.parseReplicationFilter()
		if err != nil {
			return nil, err
		}
		stmt.Filters = append(stmt.Filters, f)

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}

	// Optional: FOR CHANNEL channel
	if p.cur.Type == kwFOR {
		p.advance() // consume FOR
		p.advance() // consume CHANNEL (keyword token)
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Channel = name
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseReplicationFilter parses a single filter in CHANGE REPLICATION FILTER.
func (p *Parser) parseReplicationFilter() (*nodes.ReplicationFilter, error) {
	start := p.pos()
	f := &nodes.ReplicationFilter{Loc: nodes.Loc{Start: start}}

	// Filter type is an identifier like REPLICATE_DO_DB.
	// Must accept any keyword since filter type names may collide with registered keywords.
	filterTokType := p.cur.Type
	name, _, err := p.parseKeywordOrIdent()
	if err != nil {
		return nil, err
	}
	f.Type = name

	// consume '='
	p.match('=')

	// consume '('
	p.match('(')

	// Special case for REPLICATE_REWRITE_DB: (from_db, to_db)[, ...]
	if filterTokType == kwREPLICATE_REWRITE_DB {
		if p.cur.Type != ')' {
			for {
				p.match('(') // inner paren
				fromDB, _, err := p.parseIdentifier()
				if err != nil {
					return nil, err
				}
				p.match(',')
				toDB, _, err := p.parseIdentifier()
				if err != nil {
					return nil, err
				}
				p.match(')') // close inner paren
				f.Values = append(f.Values, "("+fromDB+", "+toDB+")")
				if p.cur.Type != ',' {
					break
				}
				p.advance() // consume ','
			}
		}
	} else {
		// db_list, tbl_list, or wild_tbl_list
		if p.cur.Type != ')' {
			for {
				var val string
				if p.cur.Type == tokSCONST {
					// wild pattern like 'db%.tbl%'
					val = p.cur.Str
					p.advance()
				} else {
					// db_name or db_name.table_name
					ident, _, err := p.parseIdentifier()
					if err != nil {
						return nil, err
					}
					val = ident
					if p.cur.Type == '.' {
						p.advance() // consume '.'
						ident2, _, err := p.parseIdentifier()
						if err != nil {
							return nil, err
						}
						val += "." + ident2
					}
				}
				f.Values = append(f.Values, val)
				if p.cur.Type != ',' {
					break
				}
				p.advance() // consume ','
			}
		}
	}

	p.match(')') // close outer paren

	f.Loc.End = p.pos()
	return f, nil
}

// parseStartReplicaStmt parses a START REPLICA or START SLAVE statement.
// p.cur is START.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/start-replica.html
//
//	START REPLICA [thread_types] [until_option] [connection_options] [channel_option]
//
//	thread_types:
//	    [thread_type [, thread_type] ... ]
//
//	thread_type:
//	    IO_THREAD | SQL_THREAD
//
//	until_option:
//	    UNTIL {   {SQL_BEFORE_GTIDS | SQL_AFTER_GTIDS} = gtid_set
//	          |   MASTER_LOG_FILE = 'log_name', MASTER_LOG_POS = log_pos
//	          |   SOURCE_LOG_FILE = 'log_name', SOURCE_LOG_POS = log_pos
//	          |   RELAY_LOG_FILE = 'log_name', RELAY_LOG_POS = log_pos
//	          |   SQL_AFTER_MTS_GAPS  }
//
//	connection_options:
//	    [USER='user_name'] [PASSWORD='user_pass'] [DEFAULT_AUTH='plugin_name'] [PLUGIN_DIR='plugin_dir']
//
//	channel_option:
//	    FOR CHANNEL channel
func (p *Parser) parseStartReplicaStmt(start int) (*nodes.StartReplicaStmt, error) {
	// REPLICA or SLAVE already consumed
	stmt := &nodes.StartReplicaStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Optional thread types
	p.parseThreadTypes(stmt)

	// Optional UNTIL
	if p.cur.Type == kwUNTIL {
		p.advance() // consume UNTIL
		untilTokType := p.cur.Type
		untilName, _, err := p.parseKeywordOrIdent()
		if err != nil {
			return nil, err
		}
		stmt.UntilType = untilName

		if untilTokType == kwSQL_AFTER_MTS_GAPS {
			// no value needed
		} else {
			// consume '='
			p.match('=')
			// value is a string literal or gtid set
			if p.cur.Type == tokSCONST {
				stmt.UntilValue = p.cur.Str
				p.advance()
			} else {
				val, _, err := p.parseKeywordOrIdent()
				if err != nil {
					return nil, err
				}
				stmt.UntilValue = val
			}
			// For log file modes, also consume the POS part
			if untilTokType == kwSOURCE_LOG_FILE || untilTokType == kwMASTER_LOG_FILE || untilTokType == kwRELAY_LOG_FILE {
				if p.cur.Type == ',' {
					p.advance() // consume ','
					posName, _, err := p.parseKeywordOrIdent()
					if err != nil {
						return nil, err
					}
					_ = posName // SOURCE_LOG_POS, MASTER_LOG_POS, RELAY_LOG_POS
					p.match('=')
					if p.cur.Type == tokICONST {
						pos, _ := strconv.ParseInt(p.cur.Str, 10, 64)
						stmt.UntilPos = pos
						p.advance()
					}
				}
			}
		}
	}

	// Optional connection options
	for p.isIdentToken() || p.cur.Type == kwUSER || p.cur.Type == kwPASSWORD ||
		p.cur.Type == kwDEFAULT_AUTH || p.cur.Type == kwPLUGIN_DIR {
		switch p.cur.Type {
		case kwUSER:
			p.advance()
			p.match('=')
			stmt.User = p.cur.Str
			p.advance()
		case kwPASSWORD:
			p.advance()
			p.match('=')
			stmt.Password = p.cur.Str
			p.advance()
		case kwDEFAULT_AUTH:
			p.advance()
			p.match('=')
			stmt.DefaultAuth = p.cur.Str
			p.advance()
		case kwPLUGIN_DIR:
			p.advance()
			p.match('=')
			stmt.PluginDir = p.cur.Str
			p.advance()
		default:
			goto doneConnOpts
		}
	}
doneConnOpts:

	// Optional FOR CHANNEL
	if p.cur.Type == kwFOR {
		p.advance() // consume FOR
		p.advance() // consume CHANNEL
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Channel = name
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseThreadTypes parses optional IO_THREAD, SQL_THREAD for START/STOP REPLICA.
func (p *Parser) parseThreadTypes(stmt interface{}) {
	for p.cur.Type == kwIO_THREAD || p.cur.Type == kwSQL_THREAD {
		if p.cur.Type == kwIO_THREAD {
			switch s := stmt.(type) {
			case *nodes.StartReplicaStmt:
				s.IOThread = true
			case *nodes.StopReplicaStmt:
				s.IOThread = true
			}
			p.advance()
			if p.cur.Type == ',' {
				p.advance()
			}
		} else if p.cur.Type == kwSQL_THREAD {
			switch s := stmt.(type) {
			case *nodes.StartReplicaStmt:
				s.SQLThread = true
			case *nodes.StopReplicaStmt:
				s.SQLThread = true
			}
			p.advance()
			if p.cur.Type == ',' {
				p.advance()
			}
		} else {
			break
		}
	}
}

// parseStopReplicaStmt parses a STOP REPLICA or STOP SLAVE statement.
// p.cur is STOP.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/stop-replica.html
//
//	STOP REPLICA [thread_types] [channel_option]
//
//	thread_types:
//	    [thread_type [, thread_type] ... ]
//
//	thread_type:
//	    IO_THREAD | SQL_THREAD
//
//	channel_option:
//	    FOR CHANNEL channel
func (p *Parser) parseStopReplicaStmt(start int) (*nodes.StopReplicaStmt, error) {
	stmt := &nodes.StopReplicaStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Optional thread types
	p.parseThreadTypes(stmt)

	// Optional FOR CHANNEL
	if p.cur.Type == kwFOR {
		p.advance() // consume FOR
		p.advance() // consume CHANNEL
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Channel = name
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseResetReplicaStmt parses a RESET REPLICA or RESET SLAVE statement.
// RESET already consumed by caller. p.cur is REPLICA or SLAVE.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/reset-replica.html
//
//	RESET REPLICA [ALL] [channel_option]
//
//	channel_option:
//	    FOR CHANNEL channel
func (p *Parser) parseResetReplicaStmt(start int) (*nodes.ResetReplicaStmt, error) {
	p.advance() // consume REPLICA or SLAVE

	stmt := &nodes.ResetReplicaStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Optional ALL
	if p.cur.Type == kwALL {
		stmt.All = true
		p.advance()
	}

	// Optional FOR CHANNEL
	if p.cur.Type == kwFOR {
		p.advance() // consume FOR
		p.advance() // consume CHANNEL
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Channel = name
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parsePurgeBinaryLogsStmt parses a PURGE BINARY LOGS or PURGE MASTER LOGS statement.
// p.cur is PURGE.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/purge-binary-logs.html
//
//	PURGE { BINARY | MASTER } LOGS { TO 'log_name' | BEFORE datetime_expr }
func (p *Parser) parsePurgeBinaryLogsStmt() (*nodes.PurgeBinaryLogsStmt, error) {
	start := p.pos()
	p.advance() // consume PURGE
	p.advance() // consume BINARY or MASTER
	p.advance() // consume LOGS

	stmt := &nodes.PurgeBinaryLogsStmt{
		Loc: nodes.Loc{Start: start},
	}

	if p.cur.Type == kwTO {
		p.advance() // consume TO
		if p.cur.Type == tokSCONST {
			stmt.To = p.cur.Str
			p.advance()
		}
	} else if p.cur.Type == kwBEFORE {
		p.advance() // consume BEFORE
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if expr == nil {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.BeforeExpr = expr
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseResetMasterStmt parses a RESET MASTER statement.
// RESET already consumed. p.cur is MASTER.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/reset-master.html
//
//	RESET MASTER [TO binary_log_file_index_number]
func (p *Parser) parseResetMasterStmt(start int) (*nodes.ResetMasterStmt, error) {
	p.advance() // consume MASTER

	stmt := &nodes.ResetMasterStmt{
		Loc: nodes.Loc{Start: start},
	}

	if p.cur.Type == kwTO {
		p.advance() // consume TO
		if p.cur.Type == tokICONST {
			val, _ := strconv.ParseInt(p.cur.Str, 10, 64)
			stmt.To = val
			p.advance()
		}
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseStartGroupReplicationStmt parses a START GROUP_REPLICATION statement.
// START already consumed. GROUP_REPLICATION identifier already consumed.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/start-group-replication.html
//
//	START GROUP_REPLICATION
//	    [USER='user_name']
//	    [, PASSWORD='user_pass']
//	    [, DEFAULT_AUTH='plugin_name']
func (p *Parser) parseStartGroupReplicationStmt(start int) (*nodes.StartGroupReplicationStmt, error) {
	stmt := &nodes.StartGroupReplicationStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Optional connection options
	for p.isIdentToken() || p.cur.Type == kwUSER || p.cur.Type == kwPASSWORD || p.cur.Type == kwDEFAULT_AUTH {
		switch p.cur.Type {
		case kwUSER:
			p.advance()
			p.match('=')
			stmt.User = p.cur.Str
			p.advance()
		case kwPASSWORD:
			p.advance()
			p.match('=')
			stmt.Password = p.cur.Str
			p.advance()
		case kwDEFAULT_AUTH:
			p.advance()
			p.match('=')
			stmt.DefaultAuth = p.cur.Str
			p.advance()
		default:
			goto groupReplDone
		}
		if p.cur.Type == ',' {
			p.advance()
		}
	}
groupReplDone:
	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseStopGroupReplicationStmt parses a STOP GROUP_REPLICATION statement.
// p.cur is STOP. (dispatched when STOP followed by GROUP_REPLICATION identifier)
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/stop-group-replication.html
//
//	STOP GROUP_REPLICATION
func (p *Parser) parseStopGroupReplicationStmt(start int) (*nodes.StopGroupReplicationStmt, error) {
	stmt := &nodes.StopGroupReplicationStmt{
		Loc: nodes.Loc{Start: start, End: p.pos()},
	}
	return stmt, nil
}
