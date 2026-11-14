// Package parser - alter_database.go implements T-SQL ALTER DATABASE and
// ALTER DATABASE SCOPED CONFIGURATION parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseAlterDatabaseStmt parses an ALTER DATABASE statement. Every action and
// option is read by its words, as SQL Server 2022 reads it.
//
// BNF: mssql/parser/bnf/alter-database-transact-sql.bnf
// Grammar: TSql160.g:2880 (alterDatabase)
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-database-transact-sql
//
//	ALTER DATABASE { database_name | CURRENT }
//	{
//	    SET <option_spec> [ ,...n ] [ WITH <termination> ]
//	  | ADD [ LOG ] FILE <filespec> [ ,...n ] [ TO FILEGROUP filegroup_name ]
//	  | ADD FILEGROUP filegroup_name [ CONTAINS { FILESTREAM | MEMORY_OPTIMIZED_DATA } ]
//	  | REMOVE FILE logical_file_name
//	  | REMOVE FILEGROUP filegroup_name
//	  | MODIFY FILE [ PRIMARY ] <filespec>
//	  | MODIFY FILEGROUP filegroup_name
//	      { { READ_ONLY | READ_WRITE | READONLY | READWRITE
//	        | AUTOGROW_SINGLE_FILE | AUTOGROW_ALL_FILES } [ WITH <termination> ]
//	      | DEFAULT | NAME = new_filegroup_name }
//	  | MODIFY NAME = new_database_name
//	  | COLLATE collation_name
//	  | REBUILD LOG [ ON [ PRIMARY ] <filespec> ]
//	}
//
// SQL Server 2022 refuses the Azure forms ScriptDOM also reads
// (MODIFY ( EDITION = ... ), PERFORM_CUTOVER).
func (p *Parser) parseAlterDatabaseStmt() (*nodes.AlterDatabaseStmt, error) {
	loc := p.pos()

	stmt := &nodes.AlterDatabaseStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Database name or CURRENT. CURRENT is a Core keyword reserved as a
	// special "current database" token in this position per SqlScriptDOM.
	if p.cur.Type == kwCURRENT {
		stmt.Name = p.cur.Str
		p.advance()
	} else {
		name, err := p.dbExpectName()
		if err != nil {
			return nil, err
		}
		stmt.Name = name
	}

	var err error
	switch {
	case p.cur.Type == kwSET:
		p.advance()
		stmt.Action = "SET"
		err = p.parseAlterDatabaseSetOptions(stmt)
	case p.cur.Type == kwCOLLATE:
		p.advance()
		stmt.Action = "COLLATE"
		stmt.TargetName, err = p.parseCollationName()
	case p.cur.Type == kwADD:
		p.advance()
		stmt.Action = "ADD"
		err = p.parseAlterDatabaseAdd(stmt)
	case p.dbWordIs("MODIFY"):
		p.advance()
		stmt.Action = "MODIFY"
		err = p.parseAlterDatabaseModify(stmt)
	case p.dbWordIs("REMOVE"):
		p.advance()
		stmt.Action = "REMOVE"
		err = p.parseAlterDatabaseRemove(stmt)
	case p.dbWordIs("REBUILD"):
		// REBUILD LOG [ ON [ PRIMARY ] <filespec> ] (TSql160.g:2930)
		p.advance()
		stmt.Action = "REBUILD"
		stmt.SubAction = "LOG"
		if err := p.dbExpectWord("LOG"); err != nil {
			return nil, err
		}
		if p.cur.Type == kwON {
			p.advance()
			if p.cur.Type == kwPRIMARY {
				p.advance()
			}
			spec, err := p.parseDatabaseFileSpec(dbFileAdd)
			if err != nil {
				return nil, err
			}
			stmt.FileSpecs = &nodes.List{Items: []nodes.Node{spec}}
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}
	if err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterDatabaseSetOptions parses SET <option_spec> [,...n] [WITH <termination>]
// (TSql160.g:3128, TSql160.g:3190).
func (p *Parser) parseAlterDatabaseSetOptions(stmt *nodes.AlterDatabaseStmt) error {
	var opts []nodes.Node
	for {
		opt, err := p.parseAlterDatabaseSetOption()
		if err != nil {
			return err
		}
		opts = append(opts, &nodes.String{Str: opt})
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	stmt.Options = &nodes.List{Items: opts}

	if p.cur.Type == kwWITH {
		p.advance()
		term, err := p.parseAlterDatabaseTermination()
		if err != nil {
			return err
		}
		stmt.Termination = term
	}
	return nil
}

// dbSetStandaloneOptions take no value (TSql160.g:3249, dbSingleIdentOption).
var dbSetStandaloneOptions = []string{
	"ONLINE", "OFFLINE", "EMERGENCY",
	"SINGLE_USER", "MULTI_USER", "RESTRICTED_USER",
	"READ_ONLY", "READ_WRITE",
	"ENABLE_BROKER", "DISABLE_BROKER", "NEW_BROKER", "ERROR_BROKER_CONVERSATIONS",
}

// dbSetOnOffOptions are written NAME { ON | OFF }, with no '='
// (TSql160.g:3352, alterDbOnOffOption); SQL Server refuses NAME = ON for them.
var dbSetOnOffOptions = []string{
	"CURSOR_CLOSE_ON_COMMIT", "AUTO_CLOSE", "AUTO_SHRINK", "AUTO_UPDATE_STATISTICS",
	"AUTO_UPDATE_STATISTICS_ASYNC", "ANSI_NULL_DEFAULT", "ANSI_NULLS", "ANSI_PADDING",
	"ANSI_WARNINGS", "ARITHABORT", "CONCAT_NULL_YIELDS_NULL", "NUMERIC_ROUNDABORT",
	"QUOTED_IDENTIFIER", "RECURSIVE_TRIGGERS", "TORN_PAGE_DETECTION", "DB_CHAINING",
	"TRUSTWORTHY", "DATE_CORRELATION_OPTIMIZATION", "ALLOW_SNAPSHOT_ISOLATION",
	"READ_COMMITTED_SNAPSHOT", "SUPPLEMENTAL_LOGGING", "HONOR_BROKER_PRIORITY",
	"VARDECIMAL_STORAGE_FORMAT", "MIXED_PAGE_ALLOCATION", "TEMPORAL_HISTORY_RETENTION",
}

// dbSetEqualsOnOffOptions are written NAME = { ON | OFF }.
var dbSetEqualsOnOffOptions = []string{"NESTED_TRIGGERS", "TRANSFORM_NOISE_WORDS"}

// dbSetKeywordOptions take one word of a closed set, with no '='.
var dbSetKeywordOptions = map[string][]string{
	"CURSOR_DEFAULT":   {"LOCAL", "GLOBAL"},                         // TSql160.g:3257
	"RECOVERY":         {"FULL", "BULK_LOGGED", "SIMPLE"},           // TSql160.g:3273
	"PAGE_VERIFY":      {"CHECKSUM", "TORN_PAGE_DETECTION", "NONE"}, // TSql160.g:3308
	"PARAMETERIZATION": {"SIMPLE", "FORCED"},                        // TSql160.g:3457
	"ENCRYPTION":       {"ON", "OFF", "SUSPEND", "RESUME"},
}

// dbSetEqualsKeywordOptions take '=' and one word of a closed set.
var dbSetEqualsKeywordOptions = map[string][]string{
	"CONTAINMENT":        {"NONE", "PARTIAL"},               // TSql160.g:4128
	"DELAYED_DURABILITY": {"DISABLED", "ALLOWED", "FORCED"}, // TSql160.g:4192
}

// dbSetSpecialOptions are the options parsed by their own branch below.
var dbSetSpecialOptions = []string{
	"AUTO_CREATE_STATISTICS", "MEMORY_OPTIMIZED_ELEVATE_TO_SNAPSHOT", "TARGET_RECOVERY_TIME",
	"COMPATIBILITY_LEVEL", "TWO_DIGIT_YEAR_CUTOFF", "DEFAULT_LANGUAGE", "DEFAULT_FULLTEXT_LANGUAGE",
	"ACCELERATED_DATABASE_RECOVERY", "SUSPEND_FOR_SNAPSHOT_BACKUP", "CHANGE_TRACKING",
	"QUERY_STORE", "AUTOMATIC_TUNING", "REMOTE_DATA_ARCHIVE", "PERSISTENT_LOG_BUFFER",
	"FILESTREAM", "HADR", "PARTNER", "WITNESS",
}

// dbSetOptionWords is every option name ALTER DATABASE SET takes.
var dbSetOptionWords = func() []string {
	words := append([]string{}, dbSetStandaloneOptions...)
	words = append(words, dbSetOnOffOptions...)
	words = append(words, dbSetEqualsOnOffOptions...)
	for k := range dbSetKeywordOptions {
		words = append(words, k)
	}
	for k := range dbSetEqualsKeywordOptions {
		words = append(words, k)
	}
	return append(words, dbSetSpecialOptions...)
}()

func dbWordListHas(words []string, w string) bool {
	for _, x := range words {
		if x == w {
			return true
		}
	}
	return false
}

// parseAlterDatabaseSetOption parses a single SET option (TSql160.g:3206,
// dbOptionStateItem) and returns it as a string. SQL Server 2022 refuses
// LEDGER, CATALOG_COLLATION, DATA_RETENTION, OPTIMIZED_LOCKING, and the Azure
// EDITION and SERVICE_OBJECTIVE here.
func (p *Parser) parseAlterDatabaseSetOption() (string, error) {
	key, err := p.dbExpectWordIn(dbSetOptionWords...)
	if err != nil {
		return "", err
	}

	switch {
	case dbWordListHas(dbSetStandaloneOptions, key):
		return key, nil
	case dbWordListHas(dbSetOnOffOptions, key):
		v, err := p.dbExpectOnOff()
		if err != nil {
			return "", err
		}
		return key + "=" + v, nil
	case dbWordListHas(dbSetEqualsOnOffOptions, key):
		if _, err := p.expect('='); err != nil {
			return "", err
		}
		v, err := p.dbExpectOnOff()
		if err != nil {
			return "", err
		}
		return key + "=" + v, nil
	}
	if vals, ok := dbSetKeywordOptions[key]; ok {
		v, err := p.dbExpectWordIn(vals...)
		if err != nil {
			return "", err
		}
		return key + "=" + v, nil
	}
	if vals, ok := dbSetEqualsKeywordOptions[key]; ok {
		if _, err := p.expect('='); err != nil {
			return "", err
		}
		v, err := p.dbExpectWordIn(vals...)
		if err != nil {
			return "", err
		}
		return key + "=" + v, nil
	}

	switch key {
	case "AUTO_CREATE_STATISTICS":
		// AUTO_CREATE_STATISTICS [ ON [ ( INCREMENTAL = { ON | OFF } ) ] | OFF ]
		// (TSql160.g:3318). SQL Server also takes the bare name.
		v, ok := p.dbWordIn("ON", "OFF")
		if !ok {
			return key, nil
		}
		p.advance()
		if v == "ON" && p.cur.Type == '(' {
			p.advance()
			if err := p.dbExpectWord("INCREMENTAL"); err != nil {
				return "", err
			}
			if _, err := p.expect('='); err != nil {
				return "", err
			}
			inc, err := p.dbExpectOnOff()
			if err != nil {
				return "", err
			}
			if _, err := p.expect(')'); err != nil {
				return "", err
			}
			return key + "=ON(INCREMENTAL=" + inc + ")", nil
		}
		return key + "=" + v, nil

	case "MEMORY_OPTIMIZED_ELEVATE_TO_SNAPSHOT":
		// Taken with or without '='.
		if p.cur.Type == '=' {
			p.advance()
		}
		v, err := p.dbExpectOnOff()
		if err != nil {
			return "", err
		}
		return key + "=" + v, nil

	case "TARGET_RECOVERY_TIME":
		// TARGET_RECOVERY_TIME = integer { SECONDS | MINUTES } (TSql160.g:3294)
		if _, err := p.expect('='); err != nil {
			return "", err
		}
		n, err := p.dbExpectInt()
		if err != nil {
			return "", err
		}
		unit, err := p.dbExpectWordIn("SECONDS", "MINUTES")
		if err != nil {
			return "", err
		}
		return key + "=" + n + " " + unit, nil

	case "COMPATIBILITY_LEVEL", "TWO_DIGIT_YEAR_CUTOFF":
		// TSql160.g:3473
		if _, err := p.expect('='); err != nil {
			return "", err
		}
		n, err := p.dbExpectInt()
		if err != nil {
			return "", err
		}
		return key + "=" + n, nil

	case "DEFAULT_LANGUAGE", "DEFAULT_FULLTEXT_LANGUAGE":
		if _, err := p.expect('='); err != nil {
			return "", err
		}
		v, err := p.dbExpectLanguage()
		if err != nil {
			return "", err
		}
		return key + "=" + strings.ToUpper(v), nil

	case "ACCELERATED_DATABASE_RECOVERY":
		// ACCELERATED_DATABASE_RECOVERY = { ON | OFF }
		//   [ ( PERSISTENT_VERSION_STORE_FILEGROUP = filegroup_name ) ] (TSql160.g:3517)
		if _, err := p.expect('='); err != nil {
			return "", err
		}
		v, err := p.dbExpectOnOff()
		if err != nil {
			return "", err
		}
		if p.cur.Type == '(' {
			p.advance()
			if err := p.dbExpectWord("PERSISTENT_VERSION_STORE_FILEGROUP"); err != nil {
				return "", err
			}
			if _, err := p.expect('='); err != nil {
				return "", err
			}
			fg, err := p.dbExpectName()
			if err != nil {
				return "", err
			}
			if _, err := p.expect(')'); err != nil {
				return "", err
			}
			return key + "=" + v + "(PERSISTENT_VERSION_STORE_FILEGROUP=" + strings.ToUpper(fg) + ")", nil
		}
		return key + "=" + v, nil

	case "SUSPEND_FOR_SNAPSHOT_BACKUP":
		// SUSPEND_FOR_SNAPSHOT_BACKUP = { ON [ ( MODE = COPY_ONLY ) ] | OFF }
		// Ref: https://learn.microsoft.com/en-us/sql/relational-databases/databases/create-a-transact-sql-snapshot-backup
		if _, err := p.expect('='); err != nil {
			return "", err
		}
		v, err := p.dbExpectOnOff()
		if err != nil {
			return "", err
		}
		if v == "ON" && p.cur.Type == '(' {
			p.advance()
			if err := p.dbExpectWord("MODE"); err != nil {
				return "", err
			}
			if _, err := p.expect('='); err != nil {
				return "", err
			}
			if err := p.dbExpectWord("COPY_ONLY"); err != nil {
				return "", err
			}
			if _, err := p.expect(')'); err != nil {
				return "", err
			}
			return key + "=ON(MODE=COPY_ONLY)", nil
		}
		return key + "=" + v, nil

	case "CHANGE_TRACKING":
		return p.parseAlterDatabaseOnOffWithList(key, p.parseAlterDatabaseChangeTrackingOptions)
	case "QUERY_STORE":
		// QUERY_STORE CLEAR [ ALL ] (TSql160.g:3599)
		if p.cur.Type == kwCLEAR {
			p.advance()
			if p.cur.Type == kwALL {
				p.advance()
				return "QUERY_STORE=CLEAR ALL", nil
			}
			return "QUERY_STORE=CLEAR", nil
		}
		if p.cur.Type == '=' && p.peekNext().Type == kwOFF {
			// QUERY_STORE = OFF [ ( FORCED [ , ...n ] ) ]
			p.advance()
			p.advance()
			if p.cur.Type != '(' {
				return key + "=OFF", nil
			}
			p.advance()
			for {
				if err := p.dbExpectWord("FORCED"); err != nil {
					return "", err
				}
				if p.cur.Type != ',' {
					break
				}
				p.advance()
			}
			if _, err := p.expect(')'); err != nil {
				return "", err
			}
			return key + "=OFF(FORCED)", nil
		}
		return p.parseAlterDatabaseOnOffWithList(key, p.parseAlterDatabaseQueryStoreOptions)
	case "REMOTE_DATA_ARCHIVE":
		return p.parseAlterDatabaseRemoteDataArchive()
	case "AUTOMATIC_TUNING":
		// AUTOMATIC_TUNING { = { AUTO | CUSTOM } | ( FORCE_LAST_GOOD_PLAN = { ON | OFF | DEFAULT } ) }
		// (TSql160.g:3903). SQL Server 2022 refuses INHERIT and the index options.
		if p.cur.Type == '=' {
			p.advance()
			v, err := p.dbExpectWordIn("AUTO", "CUSTOM")
			if err != nil {
				return "", err
			}
			return key + "=" + v, nil
		}
		if _, err := p.expect('('); err != nil {
			return "", err
		}
		if err := p.dbExpectWord("FORCE_LAST_GOOD_PLAN"); err != nil {
			return "", err
		}
		if _, err := p.expect('='); err != nil {
			return "", err
		}
		v, err := p.dbExpectWordIn("ON", "OFF", "DEFAULT")
		if err != nil {
			return "", err
		}
		if _, err := p.expect(')'); err != nil {
			return "", err
		}
		return key + "(FORCE_LAST_GOOD_PLAN=" + v + ")", nil
	case "PERSISTENT_LOG_BUFFER":
		// PERSISTENT_LOG_BUFFER = { OFF | ON ( DIRECTORY_NAME = 'path' ) }
		if _, err := p.expect('='); err != nil {
			return "", err
		}
		v, err := p.dbExpectOnOff()
		if err != nil {
			return "", err
		}
		if v == "OFF" {
			return key + "=OFF", nil
		}
		dir, err := p.parseDatabasePersistentLogBufferDir()
		if err != nil {
			return "", err
		}
		return key + "=ON(DIRECTORY_NAME='" + dir + "')", nil
	case "FILESTREAM":
		opt := &nodes.DatabaseOption{}
		if err := p.parseDatabaseFilestreamOption(opt, false); err != nil {
			return "", err
		}
		var parts []string
		if opt.FilestreamAccess != "" {
			parts = append(parts, "NON_TRANSACTED_ACCESS="+opt.FilestreamAccess)
		}
		if opt.FilestreamDirName != "" {
			parts = append(parts, "DIRECTORY_NAME="+opt.FilestreamDirName)
		}
		return "FILESTREAM(" + strings.Join(parts, ", ") + ")", nil
	case "HADR":
		return p.parseAlterDatabaseSetHadr()
	case "PARTNER":
		return p.parseAlterDatabaseSetPartner()
	case "WITNESS":
		// WITNESS { = 'witness_server' | OFF } (TSql160.g:3434)
		if p.cur.Type == '=' {
			p.advance()
			s, err := p.dbExpectString(false)
			if err != nil {
				return "", err
			}
			return "WITNESS='" + s + "'", nil
		}
		if err := p.dbExpectWord("OFF"); err != nil {
			return "", err
		}
		return "WITNESS=OFF", nil
	}
	return "", p.dbErrorAtPrev()
}

// parseAlterDatabaseOnOffWithList parses the shape CHANGE_TRACKING and
// QUERY_STORE share: = OFF | = ON [ ( list ) ] | ( list ).
func (p *Parser) parseAlterDatabaseOnOffWithList(key string, list func() (string, error)) (string, error) {
	if p.cur.Type == '(' {
		sub, err := list()
		if err != nil {
			return "", err
		}
		return key + "(" + sub + ")", nil
	}
	if _, err := p.expect('='); err != nil {
		return "", err
	}
	v, err := p.dbExpectOnOff()
	if err != nil {
		return "", err
	}
	if v == "ON" && p.cur.Type == '(' {
		sub, err := list()
		if err != nil {
			return "", err
		}
		return key + "=ON(" + sub + ")", nil
	}
	return key + "=" + v, nil
}

// parseAlterDatabaseChangeTrackingOptions parses the change tracking options
// (TSql160.g:3541, TSql160.g:3563), each at most once (Msg 5091).
//
//	( { AUTO_CLEANUP = { ON | OFF } | CHANGE_RETENTION = integer { DAYS | HOURS | MINUTES } } [ , ...n ] )
func (p *Parser) parseAlterDatabaseChangeTrackingOptions() (string, error) {
	return p.dbParseOptionList(func(key string) (string, error) {
		if key == "AUTO_CLEANUP" {
			return p.dbExpectOnOff()
		}
		n, err := p.dbExpectInt()
		if err != nil {
			return "", err
		}
		unit, err := p.dbExpectWordIn("DAYS", "HOURS", "MINUTES")
		if err != nil {
			return "", err
		}
		return n + " " + unit, nil
	}, "AUTO_CLEANUP", "CHANGE_RETENTION")
}

// parseAlterDatabaseQueryStoreOptions parses the query store options
// (TSql160.g:3644, TSql160.g:3663), each at most once (Msg 12401).
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-database-transact-sql-set-options#query-store
func (p *Parser) parseAlterDatabaseQueryStoreOptions() (string, error) {
	return p.dbParseOptionList(func(key string) (string, error) {
		switch key {
		case "OPERATION_MODE", "DESIRED_STATE":
			return p.dbExpectWordIn("READ_WRITE", "READ_ONLY")
		case "QUERY_CAPTURE_MODE":
			return p.dbExpectWordIn("ALL", "AUTO", "NONE", "CUSTOM")
		case "SIZE_BASED_CLEANUP_MODE":
			return p.dbExpectWordIn("OFF", "AUTO")
		case "WAIT_STATS_CAPTURE_MODE":
			return p.dbExpectOnOff()
		case "CLEANUP_POLICY":
			// CLEANUP_POLICY = ( STALE_QUERY_THRESHOLD_DAYS = integer )
			if p.cur.Type != '(' {
				return "", p.syntaxErrorAtCur()
			}
			sub, err := p.dbParseOptionList(func(string) (string, error) {
				return p.dbExpectInt()
			}, "STALE_QUERY_THRESHOLD_DAYS")
			if err != nil {
				return "", err
			}
			return "(" + sub + ")", nil
		case "QUERY_CAPTURE_POLICY":
			// QUERY_CAPTURE_POLICY = ( { STALE_CAPTURE_POLICY_THRESHOLD = integer { DAYS | HOURS }
			//   | EXECUTION_COUNT = integer | TOTAL_COMPILE_CPU_TIME_MS = integer
			//   | TOTAL_EXECUTION_CPU_TIME_MS = integer } [ , ...n ] )
			if p.cur.Type != '(' {
				return "", p.syntaxErrorAtCur()
			}
			sub, err := p.dbParseOptionList(func(k string) (string, error) {
				n, err := p.dbExpectInt()
				if err != nil || k != "STALE_CAPTURE_POLICY_THRESHOLD" {
					return n, err
				}
				unit, err := p.dbExpectWordIn("DAYS", "HOURS")
				if err != nil {
					return "", err
				}
				return n + " " + unit, nil
			}, "STALE_CAPTURE_POLICY_THRESHOLD", "EXECUTION_COUNT", "TOTAL_COMPILE_CPU_TIME_MS", "TOTAL_EXECUTION_CPU_TIME_MS")
			if err != nil {
				return "", err
			}
			return "(" + sub + ")", nil
		default:
			// FLUSH_INTERVAL_SECONDS, DATA_FLUSH_INTERVAL_SECONDS,
			// INTERVAL_LENGTH_MINUTES, MAX_STORAGE_SIZE_MB, MAX_PLANS_PER_QUERY
			return p.dbExpectInt()
		}
	}, "OPERATION_MODE", "DESIRED_STATE", "QUERY_CAPTURE_MODE", "SIZE_BASED_CLEANUP_MODE",
		"FLUSH_INTERVAL_SECONDS", "DATA_FLUSH_INTERVAL_SECONDS", "INTERVAL_LENGTH_MINUTES",
		"MAX_STORAGE_SIZE_MB", "MAX_PLANS_PER_QUERY", "CLEANUP_POLICY", "WAIT_STATS_CAPTURE_MODE",
		"QUERY_CAPTURE_POLICY")
}

// parseAlterDatabaseRemoteDataArchive parses REMOTE_DATA_ARCHIVE
// (TSql160.g:5171, TSql160.g:5225).
//
//	REMOTE_DATA_ARCHIVE { = OFF | = ON ( <setting> [ , ...n ] ) | ( <setting> [ , ...n ] ) }
//	<setting> ::= SERVER = 'server' | CREDENTIAL = name | FEDERATED_SERVICE_ACCOUNT = { ON | OFF }
//
// = ON needs SERVER and one of CREDENTIAL or FEDERATED_SERVICE_ACCOUNT = ON
// (Msg 10770), and CREDENTIAL does not go with FEDERATED_SERVICE_ACCOUNT = ON
// (Msg 10771).
func (p *Parser) parseAlterDatabaseRemoteDataArchive() (string, error) {
	const key = "REMOTE_DATA_ARCHIVE"
	on := false
	if p.cur.Type == '=' {
		p.advance()
		v, err := p.dbExpectOnOff()
		if err != nil {
			return "", err
		}
		if v == "OFF" {
			return key + "=OFF", nil
		}
		on = true
	}
	if p.cur.Type != '(' {
		return "", p.syntaxErrorAtCur()
	}
	vals := map[string]string{}
	sub, err := p.dbParseOptionList(func(k string) (string, error) {
		var v string
		var err error
		switch k {
		case "SERVER":
			var s string
			s, err = p.dbExpectString(false)
			v = "'" + s + "'"
		case "CREDENTIAL":
			v, err = p.dbExpectName()
		default:
			v, err = p.dbExpectOnOff()
		}
		vals[k] = v
		return v, err
	}, "SERVER", "CREDENTIAL", "FEDERATED_SERVICE_ACCOUNT")
	if err != nil {
		return "", err
	}
	fsa := vals["FEDERATED_SERVICE_ACCOUNT"] == "ON"
	_, cred := vals["CREDENTIAL"]
	if cred && fsa {
		return "", p.dbErrorAtPrev()
	}
	if on {
		_, server := vals["SERVER"]
		if !server || (!cred && !fsa) {
			return "", p.dbErrorAtPrev()
		}
		return key + "=ON(" + sub + ")", nil
	}
	return key + "(" + sub + ")", nil
}

// dbParseOptionList parses ( key = value [ , ...n ] ) where each key is one
// of the words and is given at most once; value parses what follows '='.
// It returns the items as "KEY=value" joined by ", ".
func (p *Parser) dbParseOptionList(value func(key string) (string, error), words ...string) (string, error) {
	if _, err := p.expect('('); err != nil {
		return "", err
	}
	seen := map[string]bool{}
	var parts []string
	for {
		key, err := p.dbExpectWordIn(words...)
		if err != nil {
			return "", err
		}
		if seen[key] {
			return "", p.dbErrorAtPrev()
		}
		seen[key] = true
		if _, err := p.expect('='); err != nil {
			return "", err
		}
		v, err := value(key)
		if err != nil {
			return "", err
		}
		parts = append(parts, key+"="+v)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if _, err := p.expect(')'); err != nil {
		return "", err
	}
	return strings.Join(parts, ", "), nil
}

// parseAlterDatabaseSetPartner parses the remainder of SET PARTNER (TSql160.g:3378).
// Caller has consumed PARTNER.
//
//	SET PARTNER { = 'partner_server' | FAILOVER | FORCE_SERVICE_ALLOW_DATA_LOSS
//	  | OFF | RESUME | SAFETY { FULL | OFF } | SUSPEND | TIMEOUT integer }
func (p *Parser) parseAlterDatabaseSetPartner() (string, error) {
	if p.cur.Type == '=' {
		p.advance()
		s, err := p.dbExpectString(false)
		if err != nil {
			return "", err
		}
		return "PARTNER='" + s + "'", nil
	}
	sub, err := p.dbExpectWordIn("FAILOVER", "FORCE_SERVICE_ALLOW_DATA_LOSS", "OFF", "RESUME", "SAFETY", "SUSPEND", "TIMEOUT")
	if err != nil {
		return "", err
	}
	switch sub {
	case "SAFETY":
		v, err := p.dbExpectWordIn("FULL", "OFF")
		if err != nil {
			return "", err
		}
		return "PARTNER SAFETY=" + v, nil
	case "TIMEOUT":
		n, err := p.dbExpectInt()
		if err != nil {
			return "", err
		}
		return "PARTNER TIMEOUT=" + n, nil
	case "OFF":
		return "PARTNER=OFF", nil
	}
	return "PARTNER " + sub, nil
}

// parseAlterDatabaseSetHadr parses the remainder of SET HADR (TSql160.g:4142).
// Caller has consumed HADR.
//
//	SET HADR { AVAILABILITY GROUP = group_name | OFF | SUSPEND | RESUME }
func (p *Parser) parseAlterDatabaseSetHadr() (string, error) {
	sub, err := p.dbExpectWordIn("AVAILABILITY", "OFF", "SUSPEND", "RESUME")
	if err != nil {
		return "", err
	}
	switch sub {
	case "AVAILABILITY":
		if _, err := p.expect(kwGROUP); err != nil {
			return "", err
		}
		if _, err := p.expect('='); err != nil {
			return "", err
		}
		name, err := p.dbExpectName()
		if err != nil {
			return "", err
		}
		return "HADR AVAILABILITY GROUP=" + strings.ToUpper(name), nil
	case "OFF":
		return "HADR=OFF", nil
	}
	return "HADR " + sub, nil
}

// parseAlterDatabaseTermination parses <termination> after WITH (TSql160.g:3150).
//
//	<termination> ::=
//	    ROLLBACK AFTER integer [ SECONDS ]
//	  | ROLLBACK IMMEDIATE
//	  | NO_WAIT
func (p *Parser) parseAlterDatabaseTermination() (string, error) {
	kw, err := p.dbExpectWordIn("ROLLBACK", "NO_WAIT")
	if err != nil {
		return "", err
	}
	if kw == "NO_WAIT" {
		return "NO_WAIT", nil
	}
	if p.dbWordIs("IMMEDIATE") {
		p.advance()
		return "ROLLBACK IMMEDIATE", nil
	}
	if err := p.dbExpectWord("AFTER"); err != nil {
		return "", err
	}
	n, err := p.dbExpectInt()
	if err != nil {
		return "", err
	}
	result := "ROLLBACK AFTER " + n
	if p.dbWordIs("SECONDS") {
		p.advance()
		result += " SECONDS"
	}
	return result, nil
}

// parseAlterDatabaseAdd parses ADD [LOG] FILE and ADD FILEGROUP
// (TSql160.g:2948, TSql160.g:2958, TSql160.g:2981).
func (p *Parser) parseAlterDatabaseAdd(stmt *nodes.AlterDatabaseStmt) error {
	if p.dbWordIs("FILEGROUP") {
		p.advance()
		stmt.SubAction = "FILEGROUP"
		name, err := p.dbExpectName()
		if err != nil {
			return err
		}
		stmt.TargetName = name
		if p.dbWordIs("CONTAINS") {
			p.advance()
			w, err := p.dbExpectWordIn("FILESTREAM", "MEMORY_OPTIMIZED_DATA")
			if err != nil {
				return err
			}
			stmt.Options = &nodes.List{Items: []nodes.Node{&nodes.String{Str: "CONTAINS " + w}}}
		}
		return nil
	}

	stmt.SubAction = "FILE"
	if p.dbWordIs("LOG") {
		p.advance()
		stmt.SubAction = "LOG FILE"
	}
	if _, err := p.expect(kwFILE); err != nil {
		return err
	}
	files, err := p.parseDatabaseFileSpecList(dbFileAdd)
	if err != nil {
		return err
	}
	stmt.FileSpecs = files
	// TO FILEGROUP filegroup_name (TSql160.g:3143)
	if p.cur.Type == kwTO {
		p.advance()
		if err := p.dbExpectWord("FILEGROUP"); err != nil {
			return err
		}
		name, err := p.dbExpectName()
		if err != nil {
			return err
		}
		stmt.TargetName = name
	}
	return nil
}

// parseAlterDatabaseModify parses MODIFY NAME, MODIFY FILE and MODIFY
// FILEGROUP (TSql160.g:3032).
func (p *Parser) parseAlterDatabaseModify(stmt *nodes.AlterDatabaseStmt) error {
	switch {
	case p.cur.Type == kwFILE:
		// MODIFY FILE [ PRIMARY ] <filespec> (TSql160.g:3069)
		p.advance()
		stmt.SubAction = "FILE"
		if p.cur.Type == kwPRIMARY {
			p.advance()
		}
		spec, err := p.parseDatabaseFileSpec(dbFileModify)
		if err != nil {
			return err
		}
		stmt.FileSpecs = &nodes.List{Items: []nodes.Node{spec}}
	case p.dbWordIs("NAME"):
		p.advance()
		stmt.SubAction = "NAME"
		if _, err := p.expect('='); err != nil {
			return err
		}
		name, err := p.dbExpectName()
		if err != nil {
			return err
		}
		stmt.NewName = name
	case p.dbWordIs("FILEGROUP"):
		p.advance()
		stmt.SubAction = "FILEGROUP"
		name, err := p.dbExpectName()
		if err != nil {
			return err
		}
		stmt.TargetName = name
		return p.parseAlterDatabaseModifyFilegroupOption(stmt)
	default:
		return p.syntaxErrorAtCur()
	}
	return nil
}

// parseAlterDatabaseModifyFilegroupOption parses the option after MODIFY
// FILEGROUP fg_name (TSql160.g:3083). Only an updatability or autogrow option
// takes a termination.
func (p *Parser) parseAlterDatabaseModifyFilegroupOption(stmt *nodes.AlterDatabaseStmt) error {
	if p.cur.Type == kwDEFAULT {
		p.advance()
		stmt.Options = &nodes.List{Items: []nodes.Node{&nodes.String{Str: "DEFAULT"}}}
		return nil
	}
	if p.dbWordIs("NAME") {
		p.advance()
		if _, err := p.expect('='); err != nil {
			return err
		}
		name, err := p.dbExpectName()
		if err != nil {
			return err
		}
		stmt.NewName = name
		return nil
	}
	opt, err := p.dbExpectWordIn("READ_ONLY", "READ_WRITE", "READONLY", "READWRITE",
		"AUTOGROW_SINGLE_FILE", "AUTOGROW_ALL_FILES")
	if err != nil {
		return err
	}
	stmt.Options = &nodes.List{Items: []nodes.Node{&nodes.String{Str: opt}}}
	if p.cur.Type == kwWITH {
		p.advance()
		term, err := p.parseAlterDatabaseTermination()
		if err != nil {
			return err
		}
		stmt.Termination = term
	}
	return nil
}

// parseAlterDatabaseRemove parses REMOVE FILE / REMOVE FILEGROUP (TSql160.g:3006).
func (p *Parser) parseAlterDatabaseRemove(stmt *nodes.AlterDatabaseStmt) error {
	if p.cur.Type == kwFILE {
		p.advance()
		stmt.SubAction = "FILE"
	} else if err := p.dbExpectWord("FILEGROUP"); err != nil {
		return err
	} else {
		stmt.SubAction = "FILEGROUP"
	}
	name, err := p.dbExpectName()
	if err != nil {
		return err
	}
	stmt.TargetName = name
	return nil
}

// parseAlterDatabaseScopedConfigStmt parses ALTER DATABASE SCOPED CONFIGURATION.
// Caller has consumed ALTER DATABASE SCOPED CONFIGURATION.
//
// BNF: mssql/parser/bnf/alter-database-scoped-configuration-transact-sql.bnf
// Grammar: TSql160.g:2703, TSql160.g:2733, TSql160.g:2767
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-database-scoped-configuration-transact-sql
//
//	ALTER DATABASE SCOPED CONFIGURATION
//	{ [ FOR SECONDARY ] SET <set_option>
//	| CLEAR PROCEDURE_CACHE [ plan_handle ] }
//
// SQL Server 2022 takes one option per statement, and PRIMARY only FOR
// SECONDARY (Msg 12109) and only for the options a secondary takes (Msg 12110).
func (p *Parser) parseAlterDatabaseScopedConfigStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()

	stmt := &nodes.SecurityStmt{
		Action:     "ALTER",
		ObjectType: "DATABASE SCOPED CONFIGURATION",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	var opts []nodes.Node

	if p.cur.Type == kwCLEAR {
		// CLEAR PROCEDURE_CACHE [ plan_handle ] — plan_handle is a binary constant.
		p.advance()
		if err := p.dbExpectWord("PROCEDURE_CACHE"); err != nil {
			return nil, err
		}
		optStr := "CLEAR PROCEDURE_CACHE"
		if (p.cur.Type == tokICONST || p.cur.Type == tokFCONST) &&
			(strings.HasPrefix(p.cur.Str, "0x") || strings.HasPrefix(p.cur.Str, "0X")) {
			optStr += " " + p.cur.Str
			p.advance()
		}
		opts = append(opts, &nodes.String{Str: optStr})
	} else {
		forSecondary := false
		if p.cur.Type == kwFOR {
			p.advance()
			if err := p.dbExpectWord("SECONDARY"); err != nil {
				return nil, err
			}
			forSecondary = true
			opts = append(opts, &nodes.String{Str: "FOR SECONDARY"})
		}
		if _, err := p.expect(kwSET); err != nil {
			return nil, err
		}
		opt, err := p.parseAlterDatabaseScopedConfigOption(forSecondary)
		if err != nil {
			return nil, err
		}
		opts = append(opts, opt)
	}

	stmt.Options = &nodes.List{Items: opts}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// dbScopedOnOffOptions take { ON | OFF }, and PRIMARY FOR SECONDARY.
var dbScopedOnOffOptions = []string{
	"ACCELERATED_PLAN_FORCING", "ASYNC_STATS_UPDATE_WAIT_AT_LOW_PRIORITY",
	"BATCH_MODE_ADAPTIVE_JOINS", "BATCH_MODE_MEMORY_GRANT_FEEDBACK", "BATCH_MODE_ON_ROWSTORE",
	"CE_FEEDBACK", "DEFERRED_COMPILATION_TV", "DOP_FEEDBACK",
	"EXEC_QUERY_STATS_FOR_SCALAR_FUNCTIONS", "FORCE_SHOWPLAN_RUNTIME_PARAMETER_COLLECTION",
	"INTERLEAVED_EXECUTION_TVF", "ISOLATE_SECURITY_POLICY_CARDINALITY", "LAST_QUERY_PLAN_STATS",
	"LEGACY_CARDINALITY_ESTIMATION", "LIGHTWEIGHT_QUERY_PROFILING",
	"MEMORY_GRANT_FEEDBACK_PERCENTILE_GRANT", "MEMORY_GRANT_FEEDBACK_PERSISTENCE",
	"OPTIMIZE_FOR_AD_HOC_WORKLOADS", "OPTIMIZED_PLAN_FORCING",
	"PARAMETER_SENSITIVE_PLAN_OPTIMIZATION", "PARAMETER_SNIFFING", "QUERY_OPTIMIZER_HOTFIXES",
	"ROW_MODE_MEMORY_GRANT_FEEDBACK", "TSQL_SCALAR_UDF_INLINING", "VERBOSE_TRUNCATION_WARNINGS",
	"XTP_PROCEDURE_EXECUTION_STATISTICS", "XTP_QUERY_EXECUTION_STATISTICS",
}

// dbScopedPrimaryOnlyOptions take { ON | OFF } and no FOR SECONDARY (Msg 12110).
var dbScopedPrimaryOnlyOptions = []string{"IDENTITY_CACHE", "GLOBAL_TEMPORARY_TABLE_AUTO_DROP"}

// parseAlterDatabaseScopedConfigOption parses option = value after SET.
func (p *Parser) parseAlterDatabaseScopedConfigOption(forSecondary bool) (*nodes.SecurityPrincipalOption, error) {
	optLoc := p.pos()
	words := append([]string{"MAXDOP", "ELEVATE_ONLINE", "ELEVATE_RESUMABLE",
		"PAUSED_RESUMABLE_INDEX_ABORT_DURATION_MINUTES", "LEDGER_DIGEST_STORAGE_ENDPOINT"},
		dbScopedOnOffOptions...)
	words = append(words, dbScopedPrimaryOnlyOptions...)
	key, err := p.dbExpectWordIn(words...)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect('='); err != nil {
		return nil, err
	}

	var val string
	switch {
	case key == "MAXDOP":
		// MAXDOP = { integer | PRIMARY } (TSql160.g:2788)
		if forSecondary && p.cur.Type == kwPRIMARY {
			val = "PRIMARY"
			p.advance()
		} else if val, err = p.dbExpectInt(); err != nil {
			return nil, err
		}
	case key == "ELEVATE_ONLINE" || key == "ELEVATE_RESUMABLE":
		if forSecondary {
			return nil, p.dbErrorAtPrev()
		}
		if val, err = p.dbExpectWordIn("OFF", "WHEN_SUPPORTED", "FAIL_UNSUPPORTED"); err != nil {
			return nil, err
		}
	case key == "PAUSED_RESUMABLE_INDEX_ABORT_DURATION_MINUTES":
		if forSecondary {
			return nil, p.dbErrorAtPrev()
		}
		if val, err = p.dbExpectInt(); err != nil {
			return nil, err
		}
	case key == "LEDGER_DIGEST_STORAGE_ENDPOINT":
		// LEDGER_DIGEST_STORAGE_ENDPOINT = { 'endpoint URL' | OFF }
		if p.dbWordIs("OFF") {
			val = "OFF"
			p.advance()
		} else if val, err = p.dbExpectString(false); err != nil {
			return nil, err
		}
	case dbWordListHas(dbScopedPrimaryOnlyOptions, key):
		if forSecondary {
			return nil, p.dbErrorAtPrev()
		}
		if val, err = p.dbExpectOnOff(); err != nil {
			return nil, err
		}
	default:
		// { ON | OFF | PRIMARY } (TSql160.g:5067)
		if forSecondary && p.cur.Type == kwPRIMARY {
			val = "PRIMARY"
			p.advance()
		} else if val, err = p.dbExpectOnOff(); err != nil {
			return nil, err
		}
	}

	return &nodes.SecurityPrincipalOption{
		Name:  key,
		Value: val,
		Loc:   nodes.Loc{Start: optLoc, End: p.prevEnd()},
	}, nil
}
