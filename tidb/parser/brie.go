package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// TiDB's backup, restore, recover and flashback statements.
//
// None of this is MySQL's. MySQL had BACKUP TABLE and RESTORE TABLE until
// 5.5 removed both, and mysql:8.0.46 answers 1064 for every form of this
// family, as mariadb:11.8.9 does (both read on 2026-11-04); OceanBase has no
// equivalent either, and omni's parser for each of the three refuses the lot.
// TiDB's are its own, and the grammar below is v8.5.5's, read from parser.y
// and confirmed form by form against a live v8.5.5 in
// brie_forms_oracle_test.go.
//
// Upstream calls the BACKUP and RESTORE family BRIE, for backup, restore,
// import and export, and parses all of it through one BRIEStmt rule. The
// log-backup management statements of that rule -- STOP, PAUSE, RESUME and
// PURGE BACKUP LOGS, SHOW BACKUP LOGS STATUS and METADATA, SHOW BACKUP
// METADATA, and SHOW and CANCEL BR JOB -- start with a different word each,
// so each is dispatched on its own and lands in the same node, with a Kind
// that says which. CANCEL BR JOB is read with the other CANCEL, in
// tidb_admin.go.
//
// Ref: TiDB v8.5.5 parser.y:5623-5721 (BRIEStmt), 3039-3059
// (RecoverTableStmt), 3067-3124 (FlashbackToTimestampStmt), 3134-3166
// (FlashbackTableStmt, FlashbackDatabaseStmt).
// Ref: https://docs.pingcap.com/tidb/stable/sql-statement-backup/

// brieWord returns the current token's word uppercased when it is written
// bare: a plain identifier or any keyword, never a backquoted name, a string
// literal, or a punctuation token. It returns "" otherwise, which is how the
// option list finds its end.
//
// Every keyword this family reads is read this way, which is what upstream's
// keyword tokens come to without three dozen single-use words in omni's
// keyword table. The two agree on the spelling that counts: v8.5.5 answers
// 1064 for `CONCURRENCY` 4, `backend` = 'local', snapshot 5 `minute` ago,
// online `true`, checksum `off` and rate_limit 500 `MB`/second, and takes
// the bare spelling of each. A backquoted name is taken only where the
// grammar names an Identifier, which is BACKEND's and ON_DUPLICATE's value
// and nothing else here.
func (p *Parser) brieWord() string {
	if p.cur.Type != tokIDENT && p.cur.Type < 700 {
		return ""
	}
	if strings.HasPrefix(p.inputText(p.cur.Loc, p.cur.End), "`") {
		return ""
	}
	return strings.ToUpper(p.cur.Str)
}

// matchBRIEWord consumes the current token when it is the bare word w.
func (p *Parser) matchBRIEWord(w string) bool {
	if p.brieWord() != w {
		return false
	}
	p.advance()
	return true
}

// brieIntegerOptions is BRIEIntegerOptionName, each taking a LengthNum.
//
// Ref: TiDB v8.5.5 parser.y:5759-5775.
var brieIntegerOptions = map[string]bool{
	"CONCURRENCY":          true,
	"RESUME":               true,
	"CHECKSUM_CONCURRENCY": true,
	"COMPRESSION_LEVEL":    true,
}

// brieBooleanOptions is BRIEBooleanOptionName, each taking a Boolean: an
// integer literal, or TRUE or FALSE.
//
// Ref: TiDB v8.5.5 parser.y:5777-5825.
var brieBooleanOptions = map[string]bool{
	"SEND_CREDENTIALS_TO_TIKV": true,
	"ONLINE":                   true,
	"CHECKPOINT":               true,
	"SKIP_SCHEMA_FILES":        true,
	"STRICT_FORMAT":            true,
	"CSV_NOT_NULL":             true,
	"CSV_BACKSLASH_ESCAPE":     true,
	"CSV_TRIM_LAST_SEPARATORS": true,
	"WAIT_TIFLASH_READY":       true,
	"WITH_SYS_TABLE":           true,
	"IGNORE_STATS":             true,
	"LOAD_STATS":               true,
}

// brieStringOptions is every option taking one string literal:
// BRIEStringOptionName, plus the five point-in-time options the grammar
// spells on their own.
//
// Ref: TiDB v8.5.5 parser.y:5827-5856 (BRIEStringOptionName), 6004-6038
// (FULL_BACKUP_STORAGE, RESTORED_TS, START_TS, UNTIL_TS, GC_TTL).
var brieStringOptions = map[string]bool{
	"TIKV_IMPORTER":       true,
	"CSV_SEPARATOR":       true,
	"CSV_DELIMITER":       true,
	"CSV_NULL":            true,
	"COMPRESSION_TYPE":    true,
	"ENCRYPTION_METHOD":   true,
	"ENCRYPTION_KEYFILE":  true,
	"FULL_BACKUP_STORAGE": true,
	"RESTORED_TS":         true,
	"START_TS":            true,
	"UNTIL_TS":            true,
	"GC_TTL":              true,
}

// brieSnapshotUnits are the TimestampUnit words SNAPSHOT n <unit> AGO takes.
// The grammar's TimestampUnit holds nine more -- MONTH, QUARTER and YEAR and
// their SQL_TSI_ spellings -- but the action converts the unit to a duration
// and fails the parse when that cannot be done, and a month has no fixed
// length: v8.5.5 answers 1064 for `SNAPSHOT 5 MONTH AGO` and takes
// `SNAPSHOT 5 WEEK AGO`.
//
// Ref: TiDB v8.5.5 parser.y:5903-5916 (the SNAPSHOT ... AGO action),
// parser.y:8768-8836 (TimestampUnit).
var brieSnapshotUnits = map[string]bool{
	"MICROSECOND":    true,
	"SECOND":         true,
	"MINUTE":         true,
	"HOUR":           true,
	"DAY":            true,
	"WEEK":           true,
	"SQL_TSI_SECOND": true,
	"SQL_TSI_MINUTE": true,
	"SQL_TSI_HOUR":   true,
	"SQL_TSI_DAY":    true,
	"SQL_TSI_WEEK":   true,
}

// brieOptionLevels are the OptionLevel words CHECKSUM and ANALYZE take
// instead of a boolean.
//
// Ref: TiDB v8.5.5 parser.y:6075-6087 (OptionLevel).
var brieOptionLevels = map[string]bool{
	"OFF":      true,
	"OPTIONAL": true,
	"REQUIRED": true,
}

// parseBackupStmt parses BACKUP, which is either a snapshot backup of
// databases or tables or the start of a log backup task.
//
//	BACKUP {DATABASE|SCHEMA} {* | db, ...} TO 'storage' [option ...]
//	BACKUP TABLE tbl, ...                  TO 'storage' [option ...]
//	BACKUP LOGS                            TO 'storage' [option ...]
//
// p.cur is kwBACKUP.
//
// Ref: TiDB v8.5.5 parser.y:5624-5639.
func (p *Parser) parseBackupStmt() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume BACKUP

	stmt := &nodes.BRIEStmt{Loc: nodes.Loc{Start: start}, Kind: nodes.BRIEBackup}
	if p.cur.Type == kwLOGS {
		p.advance()
		stmt.Kind = nodes.BRIEBackupLogs
		stmt.Target = nodes.BRIETargetNone
	} else if err := p.parseBRIETables(stmt); err != nil {
		return nil, err
	}

	if _, err := p.expect(kwTO); err != nil {
		return nil, err
	}
	if err := p.parseBRIEStorageAndOptions(stmt); err != nil {
		return nil, err
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseRestoreStmt parses RESTORE, either a restore of databases or tables
// from a snapshot backup or the point-in-time restore.
//
//	RESTORE {DATABASE|SCHEMA} {* | db, ...} FROM 'storage' [option ...]
//	RESTORE TABLE tbl, ...                  FROM 'storage' [option ...]
//	RESTORE POINT                           FROM 'storage' [option ...]
//
// p.cur is kwRESTORE.
//
// Ref: TiDB v8.5.5 parser.y:5706-5721.
func (p *Parser) parseRestoreStmt() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume RESTORE

	stmt := &nodes.BRIEStmt{Loc: nodes.Loc{Start: start}, Kind: nodes.BRIERestore}
	if p.cur.Type == kwPOINT {
		p.advance()
		stmt.Kind = nodes.BRIERestorePoint
		stmt.Target = nodes.BRIETargetNone
	} else if err := p.parseBRIETables(stmt); err != nil {
		return nil, err
	}

	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}
	if err := p.parseBRIEStorageAndOptions(stmt); err != nil {
		return nil, err
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// expectBackupLogs consumes the two words BACKUP LOGS.
func (p *Parser) expectBackupLogs() error {
	if _, err := p.expect(kwBACKUP); err != nil {
		return err
	}
	_, err := p.expect(kwLOGS)
	return err
}

// parseStopBackupLogsStmt parses STOP BACKUP LOGS. p.cur is kwBACKUP.
//
// Ref: TiDB v8.5.5 parser.y:5640-5645.
func (p *Parser) parseStopBackupLogsStmt(start int) (nodes.Node, error) {
	if err := p.expectBackupLogs(); err != nil {
		return nil, err
	}
	return &nodes.BRIEStmt{Loc: nodes.Loc{Start: start, End: p.prev.End}, Kind: nodes.BRIEStopLogs}, nil
}

// parsePauseBackupLogsStmt parses PAUSE BACKUP LOGS [option ...]. p.cur is the
// word PAUSE. The option list is the whole of BRIEOptions: v8.5.5 takes
// `PAUSE BACKUP LOGS concurrency 4` as readily as GC_TTL, the one it means.
//
// Ref: TiDB v8.5.5 parser.y:5646-5652.
func (p *Parser) parsePauseBackupLogsStmt() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume PAUSE
	if err := p.expectBackupLogs(); err != nil {
		return nil, err
	}
	opts, err := p.parseBRIEOptions()
	if err != nil {
		return nil, err
	}
	return &nodes.BRIEStmt{
		Loc:     nodes.Loc{Start: start, End: p.prev.End},
		Kind:    nodes.BRIEPauseLogs,
		Options: opts,
	}, nil
}

// parseResumeBackupLogsStmt parses RESUME BACKUP LOGS, which takes no option: the
// list PAUSE reads after the same words is a 1064 here. p.cur is kwRESUME.
//
// Ref: TiDB v8.5.5 parser.y:5653-5658.
func (p *Parser) parseResumeBackupLogsStmt() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume RESUME
	if err := p.expectBackupLogs(); err != nil {
		return nil, err
	}
	return &nodes.BRIEStmt{Loc: nodes.Loc{Start: start, End: p.prev.End}, Kind: nodes.BRIEResumeLogs}, nil
}

// parsePurgeStmt parses PURGE, and TiDB has one:
//
//	PURGE BACKUP LOGS FROM 'storage' [option ...]
//
// MySQL's PURGE {BINARY | MASTER} LOGS is a 1064 on v8.5.5, whose grammar names
// PURGE nowhere else. p.cur is kwPURGE.
//
// Ref: TiDB v8.5.5 parser.y:5659-5666.
func (p *Parser) parsePurgeStmt() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume PURGE
	if err := p.expectBackupLogs(); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}
	stmt := &nodes.BRIEStmt{Loc: nodes.Loc{Start: start}, Kind: nodes.BRIEPurgeLogs}
	if err := p.parseBRIEStorageAndOptions(stmt); err != nil {
		return nil, err
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// isShowBRIE reports whether the SHOW at p.cur starts one of BRIEStmt's SHOW arms.
// Both words are told apart without leaving SHOW, since every other target goes to
// parseShowStmt: BACKUPS is a word of its own and SHOW BACKUP alone is a 1064.
func (p *Parser) isShowBRIE() bool {
	next := p.peekNext()
	return next.Type == kwBACKUP ||
		(next.Type == tokIDENT && strings.EqualFold(next.Str, "BR") && !p.backtickQuotedAt(next.Loc))
}

// parseShowBRIEStmt parses BRIEStmt's SHOW arms. p.cur is kwSHOW, and isShowBRIE
// has said yes. None takes a LIKE, a WHERE or an option list.
//
//	SHOW BACKUP LOGS STATUS
//	SHOW BACKUP LOGS METADATA FROM 'storage'
//	SHOW BACKUP METADATA FROM 'storage'
//	SHOW BR JOB [QUERY] id
//
// Ref: TiDB v8.5.5 parser.y:5667-5693, 5701-5705.
func (p *Parser) parseShowBRIEStmt() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume SHOW
	stmt := &nodes.BRIEStmt{Loc: nodes.Loc{Start: start}}

	if p.cur.Type != kwBACKUP {
		p.advance() // consume BR
		if _, err := p.expect(kwJOB); err != nil {
			return nil, err
		}
		stmt.Kind = nodes.BRIEShowJob
		if _, ok := p.match(kwQUERY); ok {
			stmt.Kind = nodes.BRIEShowJobQuery
		}
		id, err := p.parseAdminInt64()
		if err != nil {
			return nil, err
		}
		stmt.JobID = id
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	p.advance() // consume BACKUP
	if _, ok := p.match(kwLOGS); ok {
		if _, ok := p.match(kwSTATUS); ok {
			stmt.Kind = nodes.BRIEShowLogsStatus
			stmt.Loc.End = p.prev.End
			return stmt, nil
		}
		stmt.Kind = nodes.BRIEShowLogsMetadata
	} else {
		stmt.Kind = nodes.BRIEShowBackupMetadata
	}
	if !p.matchBRIEWord("METADATA") {
		return nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}
	storage, loc, err := p.parseStringLit()
	if err != nil {
		return nil, err
	}
	stmt.Storage = storage
	stmt.StorageLoc = nodes.Loc{Start: loc, End: p.prev.End}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseBRIETables parses BRIETables, the target a BACKUP or RESTORE names.
//
//	{DATABASE|SCHEMA} *
//	{DATABASE|SCHEMA} db, ...
//	TABLE tbl, ...
//
// SCHEMA is DATABASE here, as everywhere TiDB's grammar writes DatabaseSym:
// its lexer maps the two words to one token. The asterisk is the bare token
// and is not a name -- v8.5.5 takes `DATABASE `*“ as a one-element list
// naming a database called `*` -- and it does not mix with names: `DATABASE
// *, a` and `DATABASE a, *` are both syntax errors, as are `DATABASE a.b`,
// `TABLE *` and `TABLE a.*`.
//
// Ref: TiDB v8.5.5 parser.y:5725-5747 (BRIETables, DBNameList),
// parser.y:5318 (DatabaseSym).
func (p *Parser) parseBRIETables(stmt *nodes.BRIEStmt) error {
	switch p.cur.Type {
	case kwDATABASE, kwSCHEMA:
		p.advance()
		if p.cur.Type == '*' {
			p.advance()
			stmt.Target = nodes.BRIETargetAllSchemas
			return nil
		}
		names, err := p.parseBRIEDBNameList()
		if err != nil {
			return err
		}
		stmt.Target = nodes.BRIETargetSchemas
		stmt.Schemas = names
		return nil

	case kwTABLE:
		p.advance()
		tables, err := p.parseBRIETableNameList()
		if err != nil {
			return err
		}
		stmt.Target = nodes.BRIETargetTables
		stmt.Tables = tables
		return nil
	}
	return p.syntaxErrorAtCur()
}

// parseBRIEDBNameList parses DBNameList, one or more DBName, each a bare
// identifier.
func (p *Parser) parseBRIEDBNameList() ([]string, error) {
	var names []string
	for {
		name, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		names = append(names, name)
		if p.cur.Type != ',' {
			return names, nil
		}
		p.advance()
	}
}

// parseBRIETableNameList parses TableNameList, one or more table names, each
// bare or qualified.
func (p *Parser) parseBRIETableNameList() ([]*nodes.TableRef, error) {
	var tables []*nodes.TableRef
	for {
		ref, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		tables = append(tables, ref)
		if p.cur.Type != ',' {
			return tables, nil
		}
		p.advance()
	}
}

// parseBRIEStorageAndOptions parses the storage URL and the option list that
// close every form, with TO or FROM already consumed.
//
// The URL is the grammar's stringLit, one literal: `TO 'noop' '://'` is a
// syntax error on v8.5.5, where the same two literals concatenate inside an
// expression.
func (p *Parser) parseBRIEStorageAndOptions(stmt *nodes.BRIEStmt) error {
	storage, loc, err := p.parseStringLit()
	if err != nil {
		return err
	}
	stmt.Storage = storage
	stmt.StorageLoc = nodes.Loc{Start: loc, End: p.prev.End}

	opts, err := p.parseBRIEOptions()
	if err != nil {
		return err
	}
	stmt.Options = opts
	return nil
}

// parseBRIEOptions parses BRIEOptions, a list with no separators that may be
// empty. A repeated option is not a syntax error and is kept twice; v8.5.5
// takes `concurrency 4 concurrency 8`.
//
// The list ends at the first token that is not a bare word, which is the
// statement's end: the production has nothing after it. A bare word that
// names no option is a syntax error here rather than the start of another
// statement, which is what v8.5.5 answers for `... TO 'noop://' nosuchopt 1`
// and for `... TO 'noop://' select 1`.
//
// Ref: TiDB v8.5.5 parser.y:5749-5757.
func (p *Parser) parseBRIEOptions() ([]*nodes.BRIEOption, error) {
	var opts []*nodes.BRIEOption
	for p.brieWord() != "" {
		opt, err := p.parseBRIEOption()
		if err != nil {
			return nil, err
		}
		opts = append(opts, opt)
	}
	return opts, nil
}

// parseBRIEOption parses one BRIEOption. The current token is a bare word.
//
// Ref: TiDB v8.5.5 parser.y:5871-6039.
func (p *Parser) parseBRIEOption() (*nodes.BRIEOption, error) {
	start := p.pos()
	name := p.brieWord()
	opt := &nodes.BRIEOption{Loc: nodes.Loc{Start: start}, Name: name}

	switch {
	case brieIntegerOptions[name]:
		p.advance()
		p.eqOpt()
		val, err := p.parseLengthNum()
		if err != nil {
			return nil, err
		}
		opt.Value = nodes.BRIEValueInt
		opt.IntValue = val

	case brieBooleanOptions[name]:
		p.advance()
		p.eqOpt()
		val, err := p.parseBRIEBoolean()
		if err != nil {
			return nil, err
		}
		opt.Value = nodes.BRIEValueInt
		opt.IntValue = val

	case brieStringOptions[name]:
		p.advance()
		p.eqOpt()
		val, _, err := p.parseStringLit()
		if err != nil {
			return nil, err
		}
		opt.Value = nodes.BRIEValueString
		opt.StrValue = val

	case name == "BACKEND", name == "ON_DUPLICATE":
		p.advance()
		p.eqOpt()
		if err := p.parseBRIEKeywordOptionValue(opt); err != nil {
			return nil, err
		}

	case name == "ON":
		// BRIEKeywordOptionName spells ON_DUPLICATE both ways; the two
		// are one option, so both arrive under the one name.
		p.advance()
		if _, err := p.expect(kwDUPLICATE); err != nil {
			return nil, err
		}
		opt.Name = "ON_DUPLICATE"
		p.eqOpt()
		if err := p.parseBRIEKeywordOptionValue(opt); err != nil {
			return nil, err
		}

	case name == "SNAPSHOT":
		p.advance()
		p.eqOpt()
		if err := p.parseBRIESnapshotValue(opt); err != nil {
			return nil, err
		}

	case name == "LAST_BACKUP":
		// LAST_BACKUP takes a timestamp or a TSO, and no `n unit AGO`.
		p.advance()
		p.eqOpt()
		if err := p.parseBRIEStringOrLengthNum(opt); err != nil {
			return nil, err
		}

	case name == "RATE_LIMIT":
		p.advance()
		p.eqOpt()
		val, err := p.parseLengthNum()
		if err != nil {
			return nil, err
		}
		// MB/SECOND is the only rate the grammar spells: v8.5.5 answers
		// 1064 for `RATE_LIMIT 500 GB/second` and for `RATE_LIMIT 500 MB`.
		if !p.matchBRIEWord("MB") {
			return nil, p.syntaxErrorAtCur()
		}
		if _, err := p.expect('/'); err != nil {
			return nil, err
		}
		if !p.matchBRIEWord("SECOND") {
			return nil, p.syntaxErrorAtCur()
		}
		opt.Value = nodes.BRIEValueRateLimit
		opt.IntValue = val

	case name == "CSV_HEADER":
		p.advance()
		p.eqOpt()
		// FieldsOrColumns, or a number.
		if word := p.brieWord(); word == "FIELDS" || word == "COLUMNS" {
			p.advance()
			opt.Value = nodes.BRIEValueCSVHeaderWords
			opt.StrValue = word
			break
		}
		val, err := p.parseLengthNum()
		if err != nil {
			return nil, err
		}
		opt.Value = nodes.BRIEValueInt
		opt.IntValue = val

	case name == "CHECKSUM", name == "ANALYZE":
		p.advance()
		p.eqOpt()
		if word := p.brieWord(); brieOptionLevels[word] {
			p.advance()
			opt.Value = nodes.BRIEValueLevel
			opt.StrValue = word
			break
		}
		val, err := p.parseBRIEBoolean()
		if err != nil {
			return nil, err
		}
		opt.Value = nodes.BRIEValueInt
		opt.IntValue = val

	default:
		return nil, p.syntaxErrorAtCur()
	}

	opt.Loc.End = p.prev.End
	return opt, nil
}

// parseBRIEKeywordOptionValue parses StringNameOrBRIEOptionKeyword, the value
// BACKEND and ON_DUPLICATE take: a string literal, an identifier, or the
// keyword IGNORE or REPLACE. An identifier may be backquoted here, the one
// place in the family where one is taken.
//
// Ref: TiDB v8.5.5 parser.y:13604-13611 (StringName,
// StringNameOrBRIEOptionKeyword).
func (p *Parser) parseBRIEKeywordOptionValue(opt *nodes.BRIEOption) error {
	opt.Value = nodes.BRIEValueString
	if p.cur.Type == tokSCONST {
		val, _, err := p.parseStringLit()
		if err != nil {
			return err
		}
		opt.StrValue = val
		return nil
	}
	// IGNORE and REPLACE are reserved, so Identifier does not reach them
	// and the rule names each on its own.
	if p.cur.Type == kwIGNORE || p.cur.Type == kwREPLACE {
		opt.StrValue = p.advance().Str
		return nil
	}
	val, _, err := p.parseIdent()
	if err != nil {
		return err
	}
	opt.StrValue = val
	return nil
}

// parseBRIESnapshotValue parses SNAPSHOT's value, which is a timestamp
// string, a TSO, or a duration back from now.
//
//	SNAPSHOT [=] 'ts'
//	SNAPSHOT [=] tso
//	SNAPSHOT [=] n <unit> AGO
//
// Ref: TiDB v8.5.5 parser.y:5900-5929.
func (p *Parser) parseBRIESnapshotValue(opt *nodes.BRIEOption) error {
	if p.cur.Type == tokSCONST {
		val, _, err := p.parseStringLit()
		if err != nil {
			return err
		}
		opt.Value = nodes.BRIEValueString
		opt.StrValue = val
		return nil
	}
	val, err := p.parseLengthNum()
	if err != nil {
		return err
	}
	// Only a unit opens the duration form. Any other word is the next
	// option's name, so the number stands on its own: v8.5.5 takes
	// `SNAPSHOT 1234567890 CONCURRENCY 4`, and a word that names no
	// option fails in parseBRIEOptions, where `SNAPSHOT 5 nosuchopt 4`
	// fails on the server too. Once a unit is read, AGO is required:
	// `SNAPSHOT 5 MINUTE` and `SNAPSHOT 5 MICROSECOND CONCURRENCY 4` are
	// both 1064.
	unit := p.brieWord()
	if !brieSnapshotUnits[unit] {
		opt.Value = nodes.BRIEValueInt
		opt.IntValue = val
		return nil
	}
	p.advance()
	if !p.matchBRIEWord("AGO") {
		return p.syntaxErrorAtCur()
	}
	opt.Value = nodes.BRIEValueTimeAgo
	opt.IntValue = val
	opt.Unit = unit
	return nil
}

// parseBRIEStringOrLengthNum parses a value that is either a string literal
// or an unsigned integer, which is what LAST_BACKUP takes.
//
// Ref: TiDB v8.5.5 parser.y:5931-5943.
func (p *Parser) parseBRIEStringOrLengthNum(opt *nodes.BRIEOption) error {
	if p.cur.Type == tokSCONST {
		val, _, err := p.parseStringLit()
		if err != nil {
			return err
		}
		opt.Value = nodes.BRIEValueString
		opt.StrValue = val
		return nil
	}
	val, err := p.parseLengthNum()
	if err != nil {
		return err
	}
	opt.Value = nodes.BRIEValueInt
	opt.IntValue = val
	return nil
}

// parseBRIEBoolean parses the grammar's Boolean: an unsigned integer
// literal, or TRUE or FALSE, which come back as 1 and 0. The literal carries
// no sign and must fit in a uint64, so v8.5.5 answers 1064 for `online -0`
// and for `online 18446744073709551616`, and takes `online 2`.
//
// Ref: TiDB v8.5.5 parser.y:6061-6073.
func (p *Parser) parseBRIEBoolean() (uint64, error) {
	switch p.brieWord() {
	case "TRUE":
		p.advance()
		return 1, nil
	case "FALSE":
		p.advance()
		return 0, nil
	}
	return p.parseLengthNum()
}

// parseRecoverTableStmt parses RECOVER TABLE, which brings back a table
// dropped within the GC lifetime.
//
//	RECOVER TABLE BY JOB job_id
//	RECOVER TABLE tbl [job_num]
//
// The two are separate productions. The first names one DDL job and no
// table, so `RECOVER TABLE BY JOB` and `RECOVER TABLE BY JOB 11,12,13` are
// both syntax errors; the second names exactly one table, so `RECOVER TABLE
// t1,t2` is one too, and takes at most one number after it, so `RECOVER
// TABLE t1 100 200` and `RECOVER TABLE t1 abc` are as well. Each number is
// the grammar's Int64Num: an unsigned literal that must fit in an int64,
// which is why v8.5.5 takes `BY JOB 0` and answers 1064 for `BY JOB -1` and
// for 9223372036854775808.
//
// p.cur is kwRECOVER.
//
// Ref: TiDB v8.5.5 parser.y:3039-3059.
// Ref: https://docs.pingcap.com/tidb/stable/sql-statement-recover-table/
func (p *Parser) parseRecoverTableStmt() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume RECOVER

	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}

	stmt := &nodes.RecoverTableStmt{Loc: nodes.Loc{Start: start}}

	if p.cur.Type == kwBY {
		p.advance()
		if _, err := p.expect(kwJOB); err != nil {
			return nil, err
		}
		// Int64Num, the same rule ADMIN's job ids read.
		id, err := p.parseAdminInt64()
		if err != nil {
			return nil, err
		}
		stmt.ByJob = true
		stmt.JobID = id
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	table, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Table = table

	if p.cur.Type == tokICONST {
		num, err := p.parseAdminInt64()
		if err != nil {
			return nil, err
		}
		stmt.JobNum = num
		stmt.HasJobNum = true
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseFlashbackStmt parses TiDB's four FLASHBACK statements.
//
//	FLASHBACK CLUSTER              TO {TIMESTAMP 'ts' | TSO n}
//	FLASHBACK TABLE tbl, ...       TO {TIMESTAMP 'ts' | TSO n}
//	FLASHBACK {DATABASE|SCHEMA} db TO {TIMESTAMP 'ts' | TSO n}
//	FLASHBACK TABLE tbl            [TO new_name]
//	FLASHBACK {DATABASE|SCHEMA} db [TO new_db]
//
// The last two bring back a dropped object, optionally under another name;
// the first three roll data back to a point in time. CLUSTER has no
// bring-back form, so `FLASHBACK CLUSTER` and `FLASHBACK CLUSTER TO c` are
// both syntax errors.
//
// Which form a TO introduces is decided the way upstream's lexer decides it:
// `TO TIMESTAMP` is the point-in-time form only when a string literal
// follows, and `TO TSO` only when an integer literal does. Everything else
// after TO is a new name. That is why v8.5.5 takes `FLASHBACK TABLE t TO
// timestamp`, renaming to a table called timestamp, and `FLASHBACK TABLE t
// TO tso` the same way, while `FLASHBACK CLUSTER TO timestamp` is a syntax
// error: the cluster form has no name to rename to.
//
// p.cur is kwFLASHBACK.
//
// Ref: TiDB v8.5.5 parser.y:3067-3124 (FlashbackToTimestampStmt),
// 3134-3156 (FlashbackTableStmt, FlashbackToNewName), 3160-3166
// (FlashbackDatabaseStmt); lexer.go:272-290 (the toTimestamp and toTSO
// tokens).
func (p *Parser) parseFlashbackStmt() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume FLASHBACK

	switch {
	case p.matchBRIEWord("CLUSTER"):
		// CLUSTER is read as a bare word for the reason brieWord gives:
		// upstream has it as a keyword, and nothing but FLASHBACK reads
		// it, so `FLASHBACK `cluster`` is a syntax error on both sides.
		ts, err := p.parseFlashbackPoint()
		if err != nil {
			return nil, err
		}
		ts.Loc.Start = start
		ts.Scope = nodes.FlashbackCluster
		ts.DatabaseLoc = nodes.NoLoc()
		ts.Loc.End = p.prev.End
		return ts, nil

	case p.cur.Type == kwTABLE:
		p.advance()
		tables, err := p.parseBRIETableNameList()
		if err != nil {
			return nil, err
		}
		if p.atFlashbackPoint() {
			ts, err := p.parseFlashbackPoint()
			if err != nil {
				return nil, err
			}
			ts.Loc = nodes.Loc{Start: start, End: p.prev.End}
			ts.Scope = nodes.FlashbackTables
			ts.Tables = tables
			ts.DatabaseLoc = nodes.NoLoc()
			return ts, nil
		}
		// The bring-back form names exactly one table.
		if len(tables) != 1 {
			return nil, p.syntaxErrorAtCur()
		}
		stmt := &nodes.FlashbackTableStmt{
			Loc:        nodes.Loc{Start: start},
			Table:      tables[0],
			NewNameLoc: nodes.NoLoc(),
		}
		if _, ok := p.match(kwTO); ok {
			// FlashbackToNewName reads an Identifier, never a
			// qualified name: `TO c.d` is a syntax error.
			name, loc, err := p.parseIdent()
			if err != nil {
				return nil, err
			}
			stmt.NewName = name
			stmt.NewNameLoc = nodes.Loc{Start: loc, End: p.prev.End}
		}
		stmt.Loc.End = p.prev.End
		return stmt, nil

	case p.cur.Type == kwDATABASE, p.cur.Type == kwSCHEMA:
		p.advance()
		name, nameLoc, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		nameRange := nodes.Loc{Start: nameLoc, End: p.prev.End}
		if p.atFlashbackPoint() {
			ts, err := p.parseFlashbackPoint()
			if err != nil {
				return nil, err
			}
			ts.Loc = nodes.Loc{Start: start, End: p.prev.End}
			ts.Scope = nodes.FlashbackDatabase
			ts.Database = name
			ts.DatabaseLoc = nameRange
			return ts, nil
		}
		stmt := &nodes.FlashbackDatabaseStmt{
			Loc:        nodes.Loc{Start: start},
			Name:       name,
			NameLoc:    nameRange,
			NewNameLoc: nodes.NoLoc(),
		}
		if _, ok := p.match(kwTO); ok {
			name, loc, err := p.parseIdent()
			if err != nil {
				return nil, err
			}
			stmt.NewName = name
			stmt.NewNameLoc = nodes.Loc{Start: loc, End: p.prev.End}
		}
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}
	return nil, p.syntaxErrorAtCur()
}

// atFlashbackPoint reports whether the tokens at p.cur open a point-in-time
// clause rather than a new name. It is upstream's three-token lexer rule:
// TO TIMESTAMP followed by a string literal, or TO TSO followed by an
// integer literal, and nothing less. A third token has to be seen to tell
// `TO TIMESTAMP '2021-05-26 16:45:26'` from a rename to a table called
// timestamp, which is one token further than peekNext reaches, so the scan
// runs on a snapshot and rewinds.
func (p *Parser) atFlashbackPoint() bool {
	if p.cur.Type != kwTO {
		return false
	}
	saved := p.snapshotTokenStream()
	defer p.restoreTokenStream(saved)
	p.advance() // consume TO
	switch {
	case p.cur.Type == kwTIMESTAMP:
		return p.peekNext().Type == tokSCONST
	case p.brieWord() == "TSO":
		return p.peekNext().Type == tokICONST
	}
	return false
}

// parseFlashbackPoint parses the TO TIMESTAMP and TO TSO clauses, filling
// only the point rolled back to; the caller fills the scope and its operands.
//
// The timestamp is one string literal the parser does not read: v8.5.5 takes
// `TO TIMESTAMP 'x'` and fails it at run time. The TSO is the grammar's
// LengthNum and must be above zero, which the action checks and reports as a
// parse failure: v8.5.5 answers 1064 for `TO TSO 0`, `TO TSO -100` and
// `TO TSO 18446744073709551616`.
func (p *Parser) parseFlashbackPoint() (*nodes.FlashbackToTimestampStmt, error) {
	if _, err := p.expect(kwTO); err != nil {
		return nil, err
	}
	stmt := &nodes.FlashbackToTimestampStmt{}
	if _, ok := p.match(kwTIMESTAMP); ok {
		ts, _, err := p.parseStringLit()
		if err != nil {
			return nil, err
		}
		stmt.Timestamp = ts
		return stmt, nil
	}
	if !p.matchBRIEWord("TSO") {
		return nil, p.syntaxErrorAtCur()
	}
	tsoTok := p.cur
	tso, err := p.parseLengthNum()
	if err != nil {
		return nil, err
	}
	if tso == 0 {
		return nil, &ParseError{
			Message:  "invalid TSO value: it must be above zero",
			Position: tsoTok.Loc,
		}
	}
	stmt.IsTSO = true
	stmt.TSO = tso
	return stmt, nil
}
