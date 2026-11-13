package parser

import (
	"math"
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// replValue is what a replication option takes after its EQ, as source_def,
// source_file_def and filter_def in MySQL 8.0 sql/sql_yacc.yy spell it.
type replValue int

const (
	replString        replValue = iota // TEXT_STRING_sys_nonewline
	replStringNL                       // TEXT_STRING_sys, which may hold a newline
	replUlong                          // ulong_num
	replRealUlong                      // real_ulong_num: no decimal
	replUlonglong                      // ulonglong_num: no hex
	replNum                            // NUM_literal: no hex
	replServerIDs                      // '(' ignore_server_id_list ')'
	replPrivilegeUser                  // privilege_check_def
	replPKCheck                        // table_primary_key_check_def
	replAssignGtids                    // assign_gtids_to_anonymous_transactions_def
	replCiphersuites                   // source_tls_ciphersuites_def
)

// replicationSourceOptions is source_def and source_file_def: each option CHANGE
// REPLICATION SOURCE TO and CHANGE MASTER TO take, under either spelling the
// grammar gives it, and its value. The grammar does not tie a spelling to the
// statement: `CHANGE MASTER TO SOURCE_HOST = 'h'` parses.
var replicationSourceOptions = map[string]replValue{
	"NETWORK_NAMESPACE":                      replString,
	"PRIVILEGE_CHECKS_USER":                  replPrivilegeUser,
	"REQUIRE_ROW_FORMAT":                     replUlong,
	"REQUIRE_TABLE_PRIMARY_KEY_CHECK":        replPKCheck,
	"SOURCE_CONNECTION_AUTO_FAILOVER":        replRealUlong,
	"ASSIGN_GTIDS_TO_ANONYMOUS_TRANSACTIONS": replAssignGtids,
	"GTID_ONLY":                              replRealUlong,
	"IGNORE_SERVER_IDS":                      replServerIDs,
	"RELAY_LOG_FILE":                         replString,
	"RELAY_LOG_POS":                          replUlong,
	"GET_MASTER_PUBLIC_KEY":                  replUlong,
	"GET_SOURCE_PUBLIC_KEY":                  replUlong,
}

func init() {
	// The options spelled MASTER_x and SOURCE_x alike.
	for suffix, v := range map[string]replValue{
		"HOST":                   replString,
		"BIND":                   replString,
		"USER":                   replString,
		"PASSWORD":               replString,
		"PORT":                   replUlong,
		"CONNECT_RETRY":          replUlong,
		"RETRY_COUNT":            replUlong,
		"DELAY":                  replUlong,
		"SSL":                    replUlong,
		"SSL_CA":                 replString,
		"SSL_CAPATH":             replString,
		"TLS_VERSION":            replString,
		"TLS_CIPHERSUITES":       replCiphersuites,
		"SSL_CERT":               replString,
		"SSL_CIPHER":             replString,
		"SSL_KEY":                replString,
		"SSL_VERIFY_SERVER_CERT": replUlong,
		"SSL_CRL":                replString,
		"SSL_CRLPATH":            replString,
		"PUBLIC_KEY_PATH":        replString,
		"HEARTBEAT_PERIOD":       replNum,
		"COMPRESSION_ALGORITHMS": replStringNL,
		"ZSTD_COMPRESSION_LEVEL": replUlong,
		"AUTO_POSITION":          replUlong,
		"LOG_FILE":               replString,
		"LOG_POS":                replUlonglong,
	} {
		replicationSourceOptions["MASTER_"+suffix] = v
		replicationSourceOptions["SOURCE_"+suffix] = v
	}
}

// replicationWord returns the upper-cased word at the current token when it is
// a bare word, as the grammar's option keywords are, and "" otherwise: a
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

// expectReplicationEq consumes the EQ every replication option takes before its
// value. It is `=` alone: `:=` is SET_VAR.
func (p *Parser) expectReplicationEq() error {
	if p.cur.Type != '=' {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// parseReplicationString reads TEXT_STRING_sys, a quoted string. A
// TEXT_STRING_sys_nonewline refuses one that holds a newline while it parses
// (ER_WRONG_VALUE), so newline says whether one may.
func (p *Parser) parseReplicationString(newline bool) (string, error) {
	if p.cur.Type != tokSCONST {
		return "", p.syntaxErrorAtCur()
	}
	if !newline && strings.Contains(p.cur.Str, "\n") {
		return "", &ParseError{Message: "argument contains not-allowed LF", Position: p.cur.Loc}
	}
	return p.advance().Str, nil
}

// parseReplicationNumber reads ulong_num and its kin: an integer, a 0x hex
// number where the rule has HEX_NUM, and a decimal where it has DECIMAL_NUM and
// FLOAT_NUM. real_ulong_num refuses a decimal; ulonglong_num and NUM_literal
// take no hex number.
//
// Ref: ulong_num, real_ulong_num, ulonglong_num and NUM_literal in MySQL 8.0
// sql/sql_yacc.yy.
func (p *Parser) parseReplicationNumber(kind replValue) (Token, error) {
	switch p.cur.Type {
	case tokICONST:
		return p.advance(), nil
	case tokFCONST:
		if kind != replRealUlong {
			return p.advance(), nil
		}
	case tokXCONST:
		if (kind == replUlong || kind == replRealUlong) && strings.HasPrefix(strings.ToLower(p.cur.Str), "0x") {
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

// parseReplicationChannel reads opt_channel, `FOR CHANNEL 'name'`, when FOR is
// next. The channel is TEXT_STRING_sys_nonewline, a string and not a name.
//
// Ref: opt_channel in MySQL 8.0 sql/sql_yacc.yy.
func (p *Parser) parseReplicationChannel() (string, error) {
	if p.cur.Type != kwFOR {
		return "", nil
	}
	p.advance() // consume FOR
	if _, err := p.expect(kwCHANNEL); err != nil {
		return "", err
	}
	return p.parseReplicationString(false)
}

// parseChangeReplicationSourceStmtInner parses a CHANGE REPLICATION SOURCE TO
// statement. CHANGE REPLICATION have already been consumed; p.cur is SOURCE.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/change-replication-source-to.html
// Ref: change, `CHANGE change_replication_source TO_SYM source_defs
// opt_channel`, in MySQL 8.0 sql/sql_yacc.yy.
func (p *Parser) parseChangeReplicationSourceStmtInner(start int) (*nodes.ChangeReplicationSourceStmt, error) {
	p.advance() // consume SOURCE
	return p.parseSourceDefs(&nodes.ChangeReplicationSourceStmt{Loc: nodes.Loc{Start: start}})
}

// parseChangeMasterStmtInner parses CHANGE MASTER TO, the older spelling of
// CHANGE REPLICATION SOURCE TO. CHANGE MASTER have already been consumed; p.cur
// is TO.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/change-master-to.html
func (p *Parser) parseChangeMasterStmtInner(start int) (*nodes.ChangeReplicationSourceStmt, error) {
	return p.parseSourceDefs(&nodes.ChangeReplicationSourceStmt{Loc: nodes.Loc{Start: start}, Legacy: true})
}

// parseSourceDefs reads `TO source_defs opt_channel` into stmt.
func (p *Parser) parseSourceDefs(stmt *nodes.ChangeReplicationSourceStmt) (*nodes.ChangeReplicationSourceStmt, error) {
	if _, err := p.expect(kwTO); err != nil {
		return nil, err
	}
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
	channel, err := p.parseReplicationChannel()
	if err != nil {
		return nil, err
	}
	stmt.Channel = channel
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseReplicationSourceOption parses one source_def: an option keyword, `=`,
// and the value the keyword takes (replicationSourceOptions).
//
// Ref: source_def and source_file_def in MySQL 8.0 sql/sql_yacc.yy.
func (p *Parser) parseReplicationSourceOption() (*nodes.ReplicationOption, error) {
	opt := &nodes.ReplicationOption{Loc: nodes.Loc{Start: p.pos()}}
	name := p.replicationWord()
	kind, ok := replicationSourceOptions[name]
	if !ok {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume the option keyword
	opt.Name = name
	if err := p.expectReplicationEq(); err != nil {
		return nil, err
	}

	switch kind {
	case replString, replStringNL:
		val, err := p.parseReplicationString(kind == replStringNL)
		if err != nil {
			return nil, err
		}
		opt.Value = val
	case replUlong, replRealUlong, replUlonglong, replNum:
		tok, err := p.parseReplicationNumber(kind)
		if err != nil {
			return nil, err
		}
		// Two options hold their number to a range while the statement parses:
		// the delay to 31 bits (1729) and the heartbeat period to 4294967
		// seconds (1624). The delay is ulong_num, read by its whole part, so
		// 1e30 is 1.
		switch name {
		case "SOURCE_DELAY", "MASTER_DELAY":
			if replicationWhole(tok) > math.MaxInt32 {
				return nil, &ParseError{Message: "The requested value " + tok.Str + " for the source delay exceeds the maximum 2147483647", Position: tok.Loc}
			}
		case "SOURCE_HEARTBEAT_PERIOD", "MASTER_HEARTBEAT_PERIOD":
			if period, err := strconv.ParseFloat(tok.Str, 64); err != nil || period > 4294967 {
				return nil, &ParseError{Message: "The requested value for the heartbeat period exceeds the maximum 4294967", Position: tok.Loc}
			}
		}
		opt.Value = tok.Str
	case replServerIDs:
		// ignore_server_id_list may be empty before each comma: `(, 1)` parses.
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		if p.cur.Type != ',' && p.cur.Type != ')' {
			tok, err := p.parseReplicationNumber(replUlong)
			if err != nil {
				return nil, err
			}
			opt.IDs = append(opt.IDs, tok.Ival)
		}
		for p.cur.Type == ',' {
			p.advance() // consume ','
			tok, err := p.parseReplicationNumber(replUlong)
			if err != nil {
				return nil, err
			}
			opt.IDs = append(opt.IDs, tok.Ival)
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	case replPrivilegeUser:
		// privilege_check_def: NULL or user_ident_or_text, whose halves are
		// ident_or_text: CURRENT_USER is no user here.
		if p.cur.Type == kwNULL {
			p.advance()
			opt.Value = "NULL"
			break
		}
		if p.cur.Type == kwCURRENT_USER {
			return nil, p.syntaxErrorAtCur()
		}
		user, err := p.parseUserName()
		if err != nil {
			return nil, err
		}
		opt.Value = user
	case replPKCheck:
		val, err := p.consumeWordOptionValue("STREAM", "ON", "OFF", "GENERATE")
		if err != nil {
			return nil, err
		}
		opt.Value = strings.ToUpper(val)
	case replAssignGtids:
		// A string is a UUID, which the action checks by the 36 bytes it
		// opens with (binary_log::Uuid::is_valid at TEXT_LENGTH): ER_WRONG_VALUE
		// otherwise.
		if p.cur.Type == tokSCONST {
			if !isUUIDPrefix(p.cur.Str) {
				return nil, &ParseError{Message: "Incorrect UUID value: '" + p.cur.Str + "'", Position: p.cur.Loc}
			}
			opt.Value = p.advance().Str
			break
		}
		val, err := p.consumeWordOptionValue("OFF", "LOCAL")
		if err != nil {
			return nil, err
		}
		opt.Value = strings.ToUpper(val)
	case replCiphersuites:
		if p.cur.Type == kwNULL {
			p.advance()
			opt.Value = "NULL"
			break
		}
		val, err := p.parseReplicationString(false)
		if err != nil {
			return nil, err
		}
		opt.Value = val
	}

	opt.Loc.End = p.prev.End
	return opt, nil
}

// isUUIDPrefix reports whether s opens with a UUID in its text form,
// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx in hex digits.
func isUUIDPrefix(s string) bool {
	if len(s) < 36 {
		return false
	}
	for i := 0; i < 36; i++ {
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

// parseChangeReplicationFilterStmtInner parses a CHANGE REPLICATION FILTER
// statement. CHANGE REPLICATION have already been consumed; p.cur is FILTER.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/change-replication-filter.html
// Ref: change, `CHANGE REPLICATION FILTER_SYM filter_defs opt_channel`, in
// MySQL 8.0 sql/sql_yacc.yy.
func (p *Parser) parseChangeReplicationFilterStmtInner(start int) (*nodes.ChangeReplicationFilterStmt, error) {
	p.advance() // consume FILTER

	stmt := &nodes.ChangeReplicationFilterStmt{Loc: nodes.Loc{Start: start}}
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
	channel, err := p.parseReplicationChannel()
	if err != nil {
		return nil, err
	}
	stmt.Channel = channel
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseReplicationFilter parses one filter_def: a filter keyword, `=`, and a
// parenthesized list, which may be empty. The databases are names, the
// tables `db.table`, the wild tables strings that hold a dot
// (ER_INVALID_RPL_WILD_TABLE_FILTER_PATTERN otherwise), and the rewrites
// parenthesized pairs of names.
//
// Ref: filter_def, opt_filter_db_list, opt_filter_table_list,
// opt_filter_string_list and opt_filter_db_pair_list in MySQL 8.0
// sql/sql_yacc.yy.
func (p *Parser) parseReplicationFilter() (*nodes.ReplicationFilter, error) {
	f := &nodes.ReplicationFilter{Loc: nodes.Loc{Start: p.pos()}}
	name := p.replicationWord()
	switch name {
	case "REPLICATE_DO_DB", "REPLICATE_IGNORE_DB", "REPLICATE_DO_TABLE", "REPLICATE_IGNORE_TABLE",
		"REPLICATE_WILD_DO_TABLE", "REPLICATE_WILD_IGNORE_TABLE", "REPLICATE_REWRITE_DB":
	default:
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume the filter keyword
	f.Type = name
	if err := p.expectReplicationEq(); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	if p.cur.Type != ')' {
		for {
			val, err := p.parseReplicationFilterItem(name)
			if err != nil {
				return nil, err
			}
			f.Values = append(f.Values, val)
			if p.cur.Type != ',' {
				break
			}
			p.advance() // consume ','
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	f.Loc.End = p.prev.End
	return f, nil
}

// parseReplicationFilterItem reads one item of the list a filter takes.
func (p *Parser) parseReplicationFilterItem(filter string) (string, error) {
	switch filter {
	case "REPLICATE_WILD_DO_TABLE", "REPLICATE_WILD_IGNORE_TABLE":
		at := p.cur.Loc
		val, err := p.parseReplicationString(false)
		if err != nil {
			return "", err
		}
		if !strings.Contains(val, ".") {
			return "", &ParseError{Message: "Supplied filter list contains a value which is not in the required format 'db_pattern.table_pattern'", Position: at}
		}
		return val, nil
	case "REPLICATE_DO_TABLE", "REPLICATE_IGNORE_TABLE":
		db, _, err := p.parseIdent()
		if err != nil {
			return "", err
		}
		if _, err := p.expect('.'); err != nil {
			return "", err
		}
		table, _, err := p.parseIdent()
		if err != nil {
			return "", err
		}
		return db + "." + table, nil
	case "REPLICATE_REWRITE_DB":
		if _, err := p.expect('('); err != nil {
			return "", err
		}
		from, _, err := p.parseIdent()
		if err != nil {
			return "", err
		}
		if _, err := p.expect(','); err != nil {
			return "", err
		}
		to, _, err := p.parseIdent()
		if err != nil {
			return "", err
		}
		if _, err := p.expect(')'); err != nil {
			return "", err
		}
		return "(" + from + ", " + to + ")", nil
	default:
		name, _, err := p.parseIdent()
		return name, err
	}
}

// parseReplicaThreads reads opt_replica_thread_option_list: SQL_THREAD and
// RELAY_THREAD (IO_THREAD is the same token), comma-separated.
//
// Ref: replica_thread_option_list in MySQL 8.0 sql/sql_yacc.yy.
func (p *Parser) parseReplicaThreads() (io, sql bool, err error) {
	if p.cur.Type != kwIO_THREAD && p.cur.Type != kwRELAY_THREAD && p.cur.Type != kwSQL_THREAD {
		return false, false, nil
	}
	for {
		switch p.cur.Type {
		case kwIO_THREAD, kwRELAY_THREAD:
			io = true
		case kwSQL_THREAD:
			sql = true
		default:
			return false, false, p.syntaxErrorAtCur()
		}
		p.advance()
		if p.cur.Type != ',' {
			return io, sql, nil
		}
		p.advance() // consume ','
	}
}

// parseStartReplicaStmt parses a START REPLICA or START SLAVE statement. START
// and REPLICA or SLAVE have been consumed.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/start-replica.html
// Ref: start_replica_stmt, `START_SYM replica opt_replica_thread_option_list
// opt_replica_until opt_user_option opt_password_option
// opt_default_auth_option opt_plugin_dir_option opt_channel`, in MySQL 8.0
// sql/sql_yacc.yy.
func (p *Parser) parseStartReplicaStmt(start int) (*nodes.StartReplicaStmt, error) {
	stmt := &nodes.StartReplicaStmt{Loc: nodes.Loc{Start: start}}
	var err error
	if stmt.IOThread, stmt.SQLThread, err = p.parseReplicaThreads(); err != nil {
		return nil, err
	}
	if p.cur.Type == kwUNTIL {
		if err := p.parseReplicaUntil(stmt); err != nil {
			return nil, err
		}
	}

	// The connection options come in this order, each at most once.
	for _, o := range []struct {
		tok  int
		dest *string
	}{
		{kwUSER, &stmt.User},
		{kwPASSWORD, &stmt.Password},
		{kwDEFAULT_AUTH, &stmt.DefaultAuth},
		{kwPLUGIN_DIR, &stmt.PluginDir},
	} {
		if p.cur.Type != o.tok {
			continue
		}
		at := p.cur.Loc
		p.advance()
		if err := p.expectReplicationEq(); err != nil {
			return nil, err
		}
		val, err := p.parseReplicationString(true)
		if err != nil {
			return nil, err
		}
		// The SQL thread alone takes no account (ER_SQLTHREAD_WITH_SECURE_REPLICA).
		if stmt.SQLThread && !stmt.IOThread {
			return nil, &ParseError{Message: "Setting authentication options is not possible when only the Replication SQL Thread is being started", Position: at}
		}
		*o.dest = val
	}

	channel, err := p.parseReplicationChannel()
	if err != nil {
		return nil, err
	}
	stmt.Channel = channel
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseReplicaUntil reads `UNTIL replica_until`: a list of source_file_def, or
// one GTID set, or SQL_AFTER_MTS_GAPS. The action refuses a list that does not
// name a binary log file and its position, or a relay log file and its
// position (ER_BAD_REPLICA_UNTIL_COND).
//
// Ref: opt_replica_until and replica_until in MySQL 8.0 sql/sql_yacc.yy.
func (p *Parser) parseReplicaUntil(stmt *nodes.StartReplicaStmt) error {
	at := p.cur.Loc
	p.advance() // consume UNTIL
	switch word := p.replicationWord(); word {
	case "SQL_BEFORE_GTIDS", "SQL_AFTER_GTIDS":
		p.advance()
		if err := p.expectReplicationEq(); err != nil {
			return err
		}
		val, err := p.parseReplicationString(true)
		if err != nil {
			return err
		}
		stmt.UntilType, stmt.UntilValue = word, val
		return nil
	case "SQL_AFTER_MTS_GAPS":
		p.advance()
		stmt.UntilType = word
		return nil
	}
	seen := map[string]bool{}
	for {
		word := p.replicationWord()
		switch word {
		case "MASTER_LOG_FILE", "SOURCE_LOG_FILE", "RELAY_LOG_FILE":
			p.advance()
			if err := p.expectReplicationEq(); err != nil {
				return err
			}
			val, err := p.parseReplicationString(false)
			if err != nil {
				return err
			}
			stmt.UntilType, stmt.UntilValue = word, val
		case "MASTER_LOG_POS", "SOURCE_LOG_POS", "RELAY_LOG_POS":
			p.advance()
			if err := p.expectReplicationEq(); err != nil {
				return err
			}
			kind := replUlonglong
			if word == "RELAY_LOG_POS" {
				kind = replUlong
			}
			tok, err := p.parseReplicationNumber(kind)
			if err != nil {
				return err
			}
			stmt.UntilPos = tok.Ival
		default:
			return p.syntaxErrorAtCur()
		}
		seen[strings.TrimPrefix(strings.TrimPrefix(word, "MASTER_"), "SOURCE_")] = true
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	if !(seen["LOG_FILE"] && seen["LOG_POS"]) && !(seen["RELAY_LOG_FILE"] && seen["RELAY_LOG_POS"]) {
		return &ParseError{Message: "Incorrect parameter or combination of parameters for START REPLICA UNTIL", Position: at}
	}
	return nil
}

// parseStopReplicaStmt parses a STOP REPLICA or STOP SLAVE statement. STOP and
// REPLICA or SLAVE have been consumed.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/stop-replica.html
// Ref: stop_replica_stmt, `STOP_SYM replica opt_replica_thread_option_list
// opt_channel`, in MySQL 8.0 sql/sql_yacc.yy.
func (p *Parser) parseStopReplicaStmt(start int) (*nodes.StopReplicaStmt, error) {
	stmt := &nodes.StopReplicaStmt{Loc: nodes.Loc{Start: start}}
	var err error
	if stmt.IOThread, stmt.SQLThread, err = p.parseReplicaThreads(); err != nil {
		return nil, err
	}
	if stmt.Channel, err = p.parseReplicationChannel(); err != nil {
		return nil, err
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

	channel, err := p.parseReplicationChannel()
	if err != nil {
		return nil, err
	}
	stmt.Channel = channel

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
		if p.cur.Type == tokICONST {
			// The index is held to 1..2000000000 while the statement parses (3567).
			if p.cur.Ival < 1 || p.cur.Ival > 2000000000 {
				return nil, &ParseError{Message: "The requested value '" + p.cur.Str + "' is out of range for RESET MASTER TO", Position: p.cur.Loc}
			}
			stmt.To = p.cur.Ival
			p.advance()
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseStartGroupReplicationStmt parses a START GROUP_REPLICATION statement.
// START and GROUP_REPLICATION have been consumed.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/start-group-replication.html
// Ref: group_replication and group_replication_start_options in MySQL 8.0
// sql/sql_yacc.yy: a comma-separated list of USER, PASSWORD and DEFAULT_AUTH,
// each `= 'string'`. The actions refuse an empty user
// (ER_GROUP_REPLICATION_USER_EMPTY_MSG) and a password over 32 bytes
// (ER_GROUP_REPLICATION_PASSWORD_LENGTH).
func (p *Parser) parseStartGroupReplicationStmt(start int) (*nodes.StartGroupReplicationStmt, error) {
	stmt := &nodes.StartGroupReplicationStmt{Loc: nodes.Loc{Start: start}}
	if p.cur.Type != kwUSER && p.cur.Type != kwPASSWORD && p.cur.Type != kwDEFAULT_AUTH {
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}
	for {
		opt := p.cur
		switch opt.Type {
		case kwUSER, kwPASSWORD, kwDEFAULT_AUTH:
		default:
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		if err := p.expectReplicationEq(); err != nil {
			return nil, err
		}
		at := p.cur.Loc
		val, err := p.parseReplicationString(false)
		if err != nil {
			return nil, err
		}
		switch opt.Type {
		case kwUSER:
			if val == "" {
				return nil, &ParseError{Message: "The START GROUP_REPLICATION command failed since the username provided for recovery channel is empty.", Position: at}
			}
			stmt.User = val
		case kwPASSWORD:
			if len(val) > 32 {
				return nil, &ParseError{Message: "The START GROUP_REPLICATION command failed since the password provided for the recovery channel exceeds the maximum length of 32 characters.", Position: at}
			}
			stmt.Password = val
		default:
			stmt.DefaultAuth = val
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
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
