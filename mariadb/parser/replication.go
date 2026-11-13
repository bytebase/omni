package parser

import (
	"math"
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
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

	stmt.Loc.End = p.prev.End
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
				// ignore_server_id is ulong_num: any number, read by its whole part.
				if p.cur.Type != tokICONST && p.cur.Type != tokFCONST {
					return nil, &ParseError{Message: "expected integer in server ID list", Position: p.cur.Loc}
				}
				opt.IDs = append(opt.IDs, p.cur.Ival)
				p.advance() // consume the number
				if p.cur.Type != ',' {
					break
				}
				p.advance() // consume ','
			}
		}
		p.match(')')
		opt.Loc.End = p.prev.End
		return opt, nil
	}

	// For all other options, consume the value token (string, number, identifier, or NULL)
	switch p.cur.Type {
	case tokSCONST:
		opt.Value = p.cur.Str
		p.advance()
	case tokICONST, tokFCONST:
		// Two options hold their number to a range while the statement parses: the
		// delay to 31 bits (1729) and the heartbeat period to 4294967 seconds (1624).
		// The delay is ulong_num, read by its whole part, so 1e30 is 1.
		//
		// Ref: master_def in MariaDB 11.8 sql/sql_yacc.yy.
		switch optTokType {
		case kwMASTER_DELAY:
			if p.cur.Ival > math.MaxInt32 {
				return nil, &ParseError{Message: "The requested value " + p.cur.Str + " for the master delay exceeds the maximum 2147483647", Position: p.cur.Loc}
			}
		case kwMASTER_HEARTBEAT_PERIOD:
			if period, err := strconv.ParseFloat(p.cur.Str, 64); err != nil || period > 4294967 {
				return nil, &ParseError{Message: "The requested value for the heartbeat period exceeds the maximum 4294967", Position: p.cur.Loc}
			}
		}
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

	opt.Loc.End = p.prev.End
	return opt, nil
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

	stmt.Loc.End = p.prev.End
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

	f.Loc.End = p.prev.End
	return f, nil
}

// masterDefValue is what a master_def option takes after its '='.
type masterDefValue int

const (
	masterString    masterDefValue = iota // TEXT_STRING_sys
	masterUlong                           // ulong_num
	masterUlonglong                       // ulonglong_num
	masterHeartbeat                       // opt_plus NUM_literal
	masterIDList                          // '(' ulong_num list ')', empty before each comma
	masterUseGtid                         // CURRENT_POS | SLAVE_POS | NO
	masterBool                            // bool: ulong_num | TRUE | FALSE
)

// masterDefs is master_def and master_file_def: each option CHANGE MASTER TO
// takes, under each spelling MariaDB's lexer gives its token, and its value.
//
// Ref: master_def and master_file_def in MariaDB 11.8 sql/sql_yacc.yy, and
// sql/lex.h for MASTER_DEMOTE_TO_REPLICA.
var masterDefs = map[string]masterDefValue{
	"MASTER_HOST":                   masterString,
	"MASTER_USER":                   masterString,
	"MASTER_PASSWORD":               masterString,
	"MASTER_PORT":                   masterUlong,
	"MASTER_CONNECT_RETRY":          masterUlong,
	"MASTER_DELAY":                  masterUlong,
	"MASTER_SSL":                    masterUlong,
	"MASTER_SSL_CA":                 masterString,
	"MASTER_SSL_CAPATH":             masterString,
	"MASTER_SSL_CERT":               masterString,
	"MASTER_SSL_CIPHER":             masterString,
	"MASTER_SSL_KEY":                masterString,
	"MASTER_SSL_VERIFY_SERVER_CERT": masterUlong,
	"MASTER_SSL_CRL":                masterString,
	"MASTER_SSL_CRLPATH":            masterString,
	"MASTER_HEARTBEAT_PERIOD":       masterHeartbeat,
	"IGNORE_SERVER_IDS":             masterIDList,
	"DO_DOMAIN_IDS":                 masterIDList,
	"IGNORE_DOMAIN_IDS":             masterIDList,
	"MASTER_LOG_FILE":               masterString,
	"MASTER_LOG_POS":                masterUlonglong,
	"RELAY_LOG_FILE":                masterString,
	"RELAY_LOG_POS":                 masterUlong,
	"MASTER_USE_GTID":               masterUseGtid,
	"MASTER_DEMOTE_TO_SLAVE":        masterBool,
	"MASTER_DEMOTE_TO_REPLICA":      masterBool,
}

// masterFileDefs are the master_def options that are master_file_def too,
// which START SLAVE UNTIL takes.
var masterFileDefs = map[string]bool{
	"MASTER_LOG_FILE": true, "MASTER_LOG_POS": true, "RELAY_LOG_FILE": true, "RELAY_LOG_POS": true,
	"MASTER_USE_GTID": true, "MASTER_DEMOTE_TO_SLAVE": true, "MASTER_DEMOTE_TO_REPLICA": true,
}

// replicationWord returns the upper-cased word at the current token when it is
// written bare, as the grammar's option keywords are, and "" otherwise: a
// backquoted name, a string or a symbol is no keyword.
func (p *Parser) replicationWord() string {
	if p.cur.Type != tokIDENT && p.cur.Type < 700 {
		return ""
	}
	text := p.inputText(p.cur.Loc, p.cur.End)
	if strings.HasPrefix(text, "`") || strings.HasPrefix(text, "@") {
		return ""
	}
	return strings.ToUpper(text)
}

// parseReplicationString reads TEXT_STRING_sys, one quoted string.
func (p *Parser) parseReplicationString() (string, error) {
	if p.cur.Type != tokSCONST {
		return "", p.syntaxErrorAtCur()
	}
	return p.advance().Str, nil
}

// parseReplicationNumber reads ulong_num (an integer, a 0x hex number or a
// decimal) or, with noHex, ulonglong_num and NUM_literal, which take no hex
// number.
//
// Ref: ulong_num, ulonglong_num and NUM_literal in MariaDB 11.8 sql/sql_yacc.yy.
func (p *Parser) parseReplicationNumber(noHex bool) (Token, error) {
	switch p.cur.Type {
	case tokICONST, tokFCONST:
		return p.advance(), nil
	case tokXCONST:
		if !noHex && strings.HasPrefix(strings.ToLower(p.cur.Str), "0x") {
			return p.advance(), nil
		}
	}
	return Token{}, p.syntaxErrorAtCur()
}

// replicationWhole is the whole part of a number parseReplicationNumber read,
// as the grammar's actions take it: a hex number in base 16, anything wider
// than 64 bits as the widest.
func replicationWhole(tok Token) uint64 {
	if tok.Type == tokXCONST {
		n, err := strconv.ParseUint(tok.Str[2:], 16, 64)
		if err != nil {
			return math.MaxUint64
		}
		return n
	}
	return uint64(tok.Ival)
}

// replicationInt is replicationWhole as the int64 the AST holds, the widest
// when it does not fit.
func replicationInt(tok Token) int64 {
	n := replicationWhole(tok)
	if n > math.MaxInt64 {
		return math.MaxInt64
	}
	return int64(n)
}

// replicationValueInt decodes a number parseMasterDef kept as its text.
func replicationValueInt(s string) int64 {
	if strings.HasPrefix(strings.ToLower(s), "0x") {
		return replicationInt(Token{Type: tokXCONST, Str: s})
	}
	return numberWhole(s)
}

// parseConnectionName reads optional_connection_name, a string, when one is
// next.
func (p *Parser) parseConnectionName() (string, bool) {
	if p.cur.Type != tokSCONST {
		return "", false
	}
	return p.advance().Str, true
}

// parseForChannel reads optional_for_channel, `FOR CHANNEL 'name'`. The name
// is the connection name by another spelling, and the action refuses both
// (ER_WRONG_ARGUMENTS), so named says whether one came before.
//
// Ref: optional_for_channel and for_channel in MariaDB 11.8 sql/sql_yacc.yy.
func (p *Parser) parseForChannel(named bool) (string, bool, error) {
	if p.cur.Type != kwFOR {
		return "", false, nil
	}
	at := p.cur.Loc
	p.advance() // consume FOR
	if _, err := p.expect(kwCHANNEL); err != nil {
		return "", false, err
	}
	name, err := p.parseReplicationString()
	if err != nil {
		return "", false, err
	}
	if named {
		return "", false, &ParseError{Message: "Incorrect arguments to CONNECTION_NAME AND FOR CHANNEL CAN NOT BE SPECIFIED AT THE SAME TIME)", Position: at}
	}
	return name, true, nil
}

// parseChangeMasterStmtInner parses CHANGE MASTER. CHANGE MASTER have already
// been consumed.
//
// Ref: https://mariadb.com/kb/en/change-master-to/
// Ref: change, `CHANGE MASTER_SYM optional_connection_name TO_SYM master_defs
// optional_for_channel`, in MariaDB 11.8 sql/sql_yacc.yy.
func (p *Parser) parseChangeMasterStmtInner(start int) (*nodes.ChangeReplicationSourceStmt, error) {
	stmt := &nodes.ChangeReplicationSourceStmt{Loc: nodes.Loc{Start: start}, Legacy: true}
	name, named := p.parseConnectionName()
	stmt.Channel = name
	if _, err := p.expect(kwTO); err != nil {
		return nil, err
	}
	for {
		opt, err := p.parseMasterDef(false)
		if err != nil {
			return nil, err
		}
		stmt.Options = append(stmt.Options, opt)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	channel, ok, err := p.parseForChannel(named)
	if err != nil {
		return nil, err
	}
	if ok {
		stmt.Channel = channel
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseMasterDef reads one master_def, or with fileOnly one master_file_def.
func (p *Parser) parseMasterDef(fileOnly bool) (*nodes.ReplicationOption, error) {
	opt := &nodes.ReplicationOption{Loc: nodes.Loc{Start: p.pos()}}
	name := p.replicationWord()
	kind, ok := masterDefs[name]
	if !ok || (fileOnly && !masterFileDefs[name]) {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume the option keyword
	opt.Name = name
	if _, err := p.expect('='); err != nil {
		return nil, err
	}

	switch kind {
	case masterString:
		val, err := p.parseReplicationString()
		if err != nil {
			return nil, err
		}
		opt.Value = val
	case masterUlong, masterUlonglong:
		tok, err := p.parseReplicationNumber(kind == masterUlonglong)
		if err != nil {
			return nil, err
		}
		// MASTER_DELAY is held to 31 bits (ER_MASTER_DELAY_VALUE_OUT_OF_RANGE).
		if name == "MASTER_DELAY" && replicationWhole(tok) > math.MaxInt32 {
			return nil, &ParseError{Message: "The requested value " + tok.Str + " for the master delay exceeds the maximum 2147483647", Position: tok.Loc}
		}
		opt.Value = tok.Str
	case masterHeartbeat:
		p.match('+')
		tok, err := p.parseReplicationNumber(true)
		if err != nil {
			return nil, err
		}
		// The period is held to SLAVE_MAX_HEARTBEAT_PERIOD
		// (ER_SLAVE_HEARTBEAT_VALUE_OUT_OF_RANGE).
		if period, err := strconv.ParseFloat(tok.Str, 64); err != nil || period > 4294967 {
			return nil, &ParseError{Message: "The requested value for the heartbeat period exceeds the maximum 4294967", Position: tok.Loc}
		}
		opt.Value = tok.Str
	case masterIDList:
		// Each list may be empty before each comma: `(, 1)` parses.
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		if p.cur.Type != ',' && p.cur.Type != ')' {
			tok, err := p.parseReplicationNumber(false)
			if err != nil {
				return nil, err
			}
			opt.IDs = append(opt.IDs, replicationInt(tok))
		}
		for p.cur.Type == ',' {
			p.advance() // consume ','
			tok, err := p.parseReplicationNumber(false)
			if err != nil {
				return nil, err
			}
			opt.IDs = append(opt.IDs, replicationInt(tok))
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	case masterUseGtid:
		switch word := p.replicationWord(); word {
		case "CURRENT_POS", "SLAVE_POS", "REPLICA_POS", "NO":
			p.advance()
			opt.Value = word
		default:
			return nil, p.syntaxErrorAtCur()
		}
	case masterBool:
		switch word := p.replicationWord(); word {
		case "TRUE", "FALSE":
			p.advance()
			opt.Value = word
		default:
			tok, err := p.parseReplicationNumber(false)
			if err != nil {
				return nil, err
			}
			opt.Value = tok.Str
		}
	}
	opt.Loc.End = p.prev.End
	return opt, nil
}

// parseSlaveThreads reads slave_thread_opts: a comma-separated list of
// SQL_THREAD and RELAY_THREAD (IO_THREAD is the same token), where each item
// may be empty, so `START SLAVE ,` parses.
//
// Ref: slave_thread_opt_list and slave_thread_opt in MariaDB 11.8
// sql/sql_yacc.yy.
func (p *Parser) parseSlaveThreads() (io, sql bool) {
	for {
		switch p.cur.Type {
		case kwIO_THREAD, kwRELAY_THREAD:
			io = true
			p.advance()
		case kwSQL_THREAD:
			sql = true
			p.advance()
		}
		if p.cur.Type != ',' {
			return io, sql
		}
		p.advance() // consume ','
	}
}

// atAllSlaves reports whether ALL SLAVES (or ALL REPLICAS) is next.
func (p *Parser) atAllSlaves() bool {
	if p.cur.Type != kwALL {
		return false
	}
	next := p.peekNext()
	if next.Type == kwREPLICAS {
		return true
	}
	// SLAVES is a keyword the lexer does not know: it must be the bare word,
	// not a string or a backquoted name.
	return next.Type == tokIDENT && strings.EqualFold(p.inputText(next.Loc, next.End), "SLAVES")
}

// parseStartReplicaStmt parses START SLAVE (or REPLICA, the same token) and
// START ALL SLAVES. START and SLAVE have been consumed, or START alone when
// all is true.
//
// Ref: https://mariadb.com/kb/en/start-replica/
// Ref: slave, `START_SYM SLAVE optional_connection_name slave_thread_opts
// optional_for_channel slave_until` and `START_SYM ALL SLAVES
// slave_thread_opts`, in MariaDB 11.8 sql/sql_yacc.yy.
func (p *Parser) parseStartReplicaStmt(start int, all bool) (*nodes.StartReplicaStmt, error) {
	stmt := &nodes.StartReplicaStmt{Loc: nodes.Loc{Start: start}, AllSlaves: all}
	if all {
		p.advance() // consume ALL
		p.advance() // consume SLAVES
		stmt.IOThread, stmt.SQLThread = p.parseSlaveThreads()
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}
	name, named := p.parseConnectionName()
	stmt.Channel = name
	stmt.IOThread, stmt.SQLThread = p.parseSlaveThreads()
	channel, ok, err := p.parseForChannel(named)
	if err != nil {
		return nil, err
	}
	if ok {
		stmt.Channel = channel
	}
	if p.cur.Type == kwUNTIL {
		if err := p.parseSlaveUntil(stmt); err != nil {
			return nil, err
		}
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseSlaveUntil reads slave_until: a GTID position, or a list of
// master_file_def that names a binary log file and its position or a relay
// log file and its position, and not both (ER_BAD_SLAVE_UNTIL_COND).
//
// Ref: slave_until in MariaDB 11.8 sql/sql_yacc.yy.
func (p *Parser) parseSlaveUntil(stmt *nodes.StartReplicaStmt) error {
	at := p.cur.Loc
	p.advance() // consume UNTIL
	switch word := p.replicationWord(); word {
	case "MASTER_GTID_POS", "SQL_AFTER_GTIDS", "SQL_BEFORE_GTIDS":
		p.advance()
		if _, err := p.expect('='); err != nil {
			return err
		}
		val, err := p.parseReplicationString()
		if err != nil {
			return err
		}
		stmt.UntilType, stmt.UntilValue = word, val
		return nil
	}
	seen := map[string]bool{}
	for {
		opt, err := p.parseMasterDef(true)
		if err != nil {
			return err
		}
		seen[opt.Name] = true
		switch opt.Name {
		case "MASTER_LOG_FILE", "RELAY_LOG_FILE":
			stmt.UntilType, stmt.UntilValue = opt.Name, opt.Value
		case "MASTER_LOG_POS", "RELAY_LOG_POS":
			stmt.UntilPos = replicationValueInt(opt.Value)
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	binlog := seen["MASTER_LOG_FILE"] || seen["MASTER_LOG_POS"]
	relay := seen["RELAY_LOG_FILE"] || seen["RELAY_LOG_POS"]
	whole := (seen["MASTER_LOG_FILE"] && seen["MASTER_LOG_POS"]) || (seen["RELAY_LOG_FILE"] && seen["RELAY_LOG_POS"])
	if (binlog && relay) || !whole {
		return &ParseError{Message: "Incorrect parameter or combination of parameters for START SLAVE UNTIL", Position: at}
	}
	return nil
}

// parseStopReplicaStmt parses STOP SLAVE (or REPLICA) and STOP ALL SLAVES.
// STOP and SLAVE have been consumed, or STOP alone when all is true.
//
// Ref: https://mariadb.com/kb/en/stop-replica/
// Ref: slave, `STOP_SYM SLAVE optional_connection_name slave_thread_opts
// optional_for_channel` and `STOP_SYM ALL SLAVES slave_thread_opts`, in
// MariaDB 11.8 sql/sql_yacc.yy.
func (p *Parser) parseStopReplicaStmt(start int, all bool) (*nodes.StopReplicaStmt, error) {
	stmt := &nodes.StopReplicaStmt{Loc: nodes.Loc{Start: start}, AllSlaves: all}
	if all {
		p.advance() // consume ALL
		p.advance() // consume SLAVES
		stmt.IOThread, stmt.SQLThread = p.parseSlaveThreads()
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}
	name, named := p.parseConnectionName()
	stmt.Channel = name
	stmt.IOThread, stmt.SQLThread = p.parseSlaveThreads()
	channel, ok, err := p.parseForChannel(named)
	if err != nil {
		return nil, err
	}
	if ok {
		stmt.Channel = channel
	}
	stmt.Loc.End = p.prev.End
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

	stmt.Loc.End = p.prev.End
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

	stmt.Loc.End = p.prev.End
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
		// ulonglong_num: any number, read by its whole part.
		if p.cur.Type == tokICONST || p.cur.Type == tokFCONST {
			stmt.To = p.cur.Ival
			p.advance()
		}
	}

	stmt.Loc.End = p.prev.End
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
	stmt.Loc.End = p.prev.End
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
		Loc: nodes.Loc{Start: start, End: p.prev.End},
	}
	return stmt, nil
}
