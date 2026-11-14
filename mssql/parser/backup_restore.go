// Package parser - backup_restore.go implements T-SQL BACKUP and RESTORE statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// The option words below are what SQL Server 2022 itself accepts. ScriptDOM
// reads every option as an Identifier and leaves its name to helpers
// (BackupOptionsNoValueHelper, BackupOptionsWithValueHelper,
// RestoreOptionNoValueHelper, DeviceTypesHelper, RestoreStatementKindsHelper)
// whose tables are not in TSql160.g, so each set was taken from the pinned
// server: every word of CodeGenerationSupporter.cs and of omni's keyword table
// was sent bare, as `= 1` and as `= 'a'`, and these are the ones it did not
// refuse with error 155 or 153. The documented lists are at
// https://learn.microsoft.com/en-us/sql/t-sql/statements/backup-transact-sql
// and https://learn.microsoft.com/en-us/sql/t-sql/statements/restore-statements-arguments-transact-sql.
var (
	// backupFlagOptions take no value: WITH INIT.
	backupFlagOptions = backupWordSet(
		"CHECKSUM", "COMPRESSION", "CONTINUE_AFTER_ERROR", "COPY_ONLY",
		"DIFFERENTIAL", "FILE_SNAPSHOT", "FORMAT", "INIT", "METADATA_ONLY",
		"NOFORMAT", "NOINIT", "NORECOVERY", "NOREWIND", "NOSKIP", "NOUNLOAD",
		"NO_CHECKSUM", "NO_COMPRESSION", "NO_LOG", "NO_TRUNCATE", "RESTART",
		"REWIND", "SKIP", "SNAPSHOT", "STATS", "STOP_ON_ERROR", "UNLOAD",
	)
	// backupValueOptions take `= value`: WITH NAME = 'a'.
	backupValueOptions = backupWordSet(
		"BACKUP_OPTIONS", "BLOCKSIZE", "BUFFERCOUNT", "CREDENTIAL",
		"DESCRIPTION", "ENHANCEDINTEGRITY", "EXPIREDATE", "MAXTRANSFERSIZE",
		"MEDIADESCRIPTION", "MEDIANAME", "MEDIAPASSWORD", "NAME", "PASSWORD",
		"RETAINDAYS", "STANDBY", "STATS",
	)
	// restoreFlagOptions take no value: WITH RECOVERY.
	restoreFlagOptions = backupWordSet(
		"CHECKSUM", "COMMIT_DIFFERENTIAL_BASE", "CONTINUE_AFTER_ERROR",
		"ENABLE_BROKER", "ERROR_BROKER_CONVERSATIONS", "KEEP_CDC",
		"KEEP_REPLICATION", "KEEP_TEMPORAL_RETENTION", "LOADHISTORY",
		"METADATA_ONLY", "NEW_BROKER", "NORECOVERY", "NOREWIND", "NOUNLOAD",
		"NO_CHECKSUM", "NO_LOG", "ONLINE", "PARTIAL", "RECOVERY", "REPLACE",
		"RESTART", "RESTRICTED_USER", "REWIND", "SNAPSHOT", "SNAPSHOT_IMPORT",
		"STATS", "STOP_ON_ERROR", "UNLOAD", "VERBOSE",
	)
	// restoreValueOptions take `= value`: WITH STOPAT = '...'.
	restoreValueOptions = backupWordSet(
		"BLOCKSIZE", "BUFFERCOUNT", "CREDENTIAL", "DBNAME", "ENHANCEDINTEGRITY",
		"MAXTRANSFERSIZE", "MEDIANAME", "MEDIAPASSWORD", "PASSWORD",
		"RESTORE_OPTIONS", "SNAPSHOTRESTOREPHASE", "STANDBY", "STATS",
		"STOPAT", "STOPATMARK", "STOPBEFOREMARK",
	)
	// backupDeviceTypes name a physical device: TO DISK = '...'.
	backupDeviceTypes = backupWordSet("DISK", "TAPE", "URL", "VIRTUAL_DEVICE", "DATABASE_SNAPSHOT")
	// restoreKinds read a backup set without restoring it: RESTORE HEADERONLY FROM ...
	restoreKinds = backupWordSet("FILELISTONLY", "HEADERONLY", "LABELONLY", "REWINDONLY", "VERIFYONLY")
	// backupEncryptionAlgorithms is ScriptDOM's EncryptionAlgorithmsHelper
	// for BACKUP ... WITH ENCRYPTION (ALGORITHM = ...).
	backupEncryptionAlgorithms = backupWordSet("AES_128", "AES_192", "AES_256", "TRIPLE_DES_3KEY")
	// backupCompressionAlgorithms is SQL Server 2022's
	// COMPRESSION (ALGORITHM = ...), which TSql160.g does not have.
	backupCompressionAlgorithms = backupWordSet("MS_XPRESS", "QAT_DEFLATE")
)

func backupWordSet(words ...string) map[string]bool {
	m := make(map[string]bool, len(words))
	for _, w := range words {
		m[w] = true
	}
	return m
}

// backupWord returns the upper-case text of the current token when it is a
// bare word: a keyword or an identifier that is not bracketed or quoted.
// SQL Server reads none of the words of BACKUP and RESTORE from a delimited
// identifier (`WITH [INIT]`, `TO [DISK] = '...'` are refused).
func (p *Parser) backupWord() string {
	if p.cur.Type == tokIDENT {
		if p.isDelimited(p.cur) {
			return ""
		}
		return strings.ToUpper(p.cur.Str)
	}
	if p.cur.Type >= kwABSENT && p.cur.Str != "" {
		return strings.ToUpper(p.cur.Str)
	}
	return ""
}

// backupExpectWord consumes the current token if it is the bare word.
func (p *Parser) backupExpectWord(word string) error {
	if p.backupWord() != word {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// backupIdentifierOrVariable parses identifierOrVariable: a name or a
// @variable (TSql160.g identifierOrVariable).
func (p *Parser) backupIdentifierOrVariable() (string, error) {
	if p.cur.Type == tokVARIABLE {
		name := p.cur.Str
		p.advance()
		return name, nil
	}
	if name, ok := p.parseIdentifier(); ok {
		return name, nil
	}
	return "", p.syntaxErrorAtCur()
}

// backupStringOrVariable parses stringOrVariable: a string literal or a
// @variable.
func (p *Parser) backupStringOrVariable() (string, error) {
	switch p.cur.Type {
	case tokSCONST, tokNSCONST, tokVARIABLE:
		s := p.cur.Str
		p.advance()
		return s, nil
	}
	return "", p.syntaxErrorAtCur()
}

// backupIsInteger reports whether the current token is an Integer: decimal
// digits, not hex (SQL Server refuses `STATS = 0x10`).
func (p *Parser) backupIsInteger() bool {
	if p.cur.Type != tokICONST {
		return false
	}
	return !strings.HasPrefix(p.cur.Str, "0x") && !strings.HasPrefix(p.cur.Str, "0X")
}

// backupSignedInteger parses signedInteger. SQL Server takes a minus sign
// and refuses a plus sign (`STATS = +5`).
func (p *Parser) backupSignedInteger() (string, error) {
	sign := ""
	if p.cur.Type == '-' {
		sign = "-"
		p.advance()
	}
	if !p.backupIsInteger() {
		return "", p.syntaxErrorAtCur()
	}
	s := sign + p.cur.Str
	p.advance()
	return s, nil
}

// backupOptionValue parses the value of a BACKUP or RESTORE option:
// signedIntegerOrVariable | stringLiteral (TSql160.g backupOption) and
// signedInteger | stringOrVariable (TSql160.g restoreOption) take the same
// tokens.
func (p *Parser) backupOptionValue() (string, error) {
	if p.cur.Type == '-' || p.backupIsInteger() {
		return p.backupSignedInteger()
	}
	return p.backupStringOrVariable()
}

// parseBackupStmt parses a BACKUP DATABASE or BACKUP LOG statement. BACKUP
// CERTIFICATE, MASTER KEY and SERVICE MASTER KEY are routed elsewhere.
//
// TSql160.g:5574, 5623-5654, 5749-5821
//
//	backupStatement:  backupMain ( backupDevices )? ( backupOptions )?
//	backupMain:       Database identifierOrVariable backupFileListOpt
//	                | 'LOG' identifierOrVariable
//	backupDevices:    To devList ( 'MIRROR' To devList )*
//	backupOptions:    With backupOption ( Comma backupOption )*
func (p *Parser) parseBackupStmt() (*nodes.BackupStmt, error) {
	loc := p.pos()
	p.advance() // consume BACKUP

	// Completion: after BACKUP → DATABASE keyword
	if p.collectMode() {
		p.addTokenCandidate(kwDATABASE)
		return nil, errCollecting
	}

	stmt := &nodes.BackupStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	switch {
	case p.cur.Type == kwDATABASE:
		stmt.Type = "DATABASE"
	case p.backupWord() == "LOG":
		stmt.Type = "LOG"
	default:
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	// Completion: after BACKUP DATABASE → database_ref
	if p.collectMode() {
		p.addRuleCandidate("database_ref")
		return nil, errCollecting
	}

	name, err := p.backupIdentifierOrVariable()
	if err != nil {
		return nil, err
	}
	stmt.Database = name

	// backupFileListOpt follows DATABASE only; the server refuses
	// BACKUP LOG a FILE = '...'.
	if stmt.Type == "DATABASE" {
		if stmt.FileSpecs, err = p.parseBackupRestoreFileList(true); err != nil {
			return nil, err
		}
	}

	if p.cur.Type == kwTO {
		p.advance() // consume TO
		if stmt.Devices, err = p.parseDeviceList(); err != nil {
			return nil, err
		}
		if s, ok := stmt.Devices.Items[0].(*nodes.String); ok {
			stmt.Target = extractDevicePath(s.Str)
		}

		// mirrorTo: 'MIRROR' To devList, repeated.
		for p.backupWord() == "MIRROR" {
			stmt.MirrorTo = true
			p.advance() // consume MIRROR
			if _, err := p.expect(kwTO); err != nil {
				return nil, err
			}
			mirrorDevs, err := p.parseDeviceList()
			if err != nil {
				return nil, err
			}
			if stmt.MirrorDevice == "" {
				if s, ok := mirrorDevs.Items[0].(*nodes.String); ok {
					stmt.MirrorDevice = extractDevicePath(s.Str)
				}
			}
		}
	}

	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		if stmt.Options, err = p.parseBackupRestoreOptions(true); err != nil {
			return nil, err
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseRestoreStmt parses RESTORE DATABASE / LOG and RESTORE HEADERONLY /
// FILELISTONLY / VERIFYONLY / LABELONLY / REWINDONLY. RESTORE MASTER KEY,
// SYMMETRIC KEY and SERVICE MASTER KEY are routed elsewhere.
//
// TSql160.g:5600, 5877
//
//	restoreStatement: ( restoreMain ( From devList )? | Identifier From devList ) ( restoreOptions )?
//	restoreMain:      Database identifierOrVariable restoreFileListOpt
//	                | 'LOG' identifierOrVariable restoreFileListOpt
func (p *Parser) parseRestoreStmt() (*nodes.RestoreStmt, error) {
	loc := p.pos()
	p.advance() // consume RESTORE

	// Completion: after RESTORE → DATABASE keyword
	if p.collectMode() {
		p.addTokenCandidate(kwDATABASE)
		return nil, errCollecting
	}

	stmt := &nodes.RestoreStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	word := p.backupWord()
	switch {
	case p.cur.Type == kwDATABASE:
		stmt.Type = "DATABASE"
	case word == "LOG":
		stmt.Type = "LOG"
	case restoreKinds[word]:
		stmt.Type = word
	default:
		// The server names the word: 'X' is not a recognized RESTORE option.
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	var err error
	if stmt.Type == "DATABASE" || stmt.Type == "LOG" {
		// Completion: after RESTORE DATABASE → database_ref
		if p.collectMode() {
			p.addRuleCandidate("database_ref")
			return nil, errCollecting
		}
		if stmt.Database, err = p.backupIdentifierOrVariable(); err != nil {
			return nil, err
		}
		// The server takes READ_WRITE_FILEGROUPS after RESTORE DATABASE
		// and refuses it after RESTORE LOG.
		if stmt.FileSpecs, err = p.parseBackupRestoreFileList(stmt.Type == "DATABASE"); err != nil {
			return nil, err
		}
		if p.cur.Type == kwFROM {
			p.advance() // consume FROM
			if err := p.parseRestoreDevices(stmt); err != nil {
				return nil, err
			}
		}
	} else {
		if _, err := p.expect(kwFROM); err != nil {
			return nil, err
		}
		if err := p.parseRestoreDevices(stmt); err != nil {
			return nil, err
		}
	}

	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		if stmt.Options, err = p.parseBackupRestoreOptions(false); err != nil {
			return nil, err
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseRestoreDevices parses the devList after RESTORE ... FROM. A
// DATABASE_SNAPSHOT = 'name' device also sets SnapshotName.
func (p *Parser) parseRestoreDevices(stmt *nodes.RestoreStmt) error {
	devs, err := p.parseDeviceList()
	if err != nil {
		return err
	}
	stmt.Devices = devs
	if s, ok := devs.Items[0].(*nodes.String); ok {
		stmt.Source = extractDevicePath(s.Str)
		if strings.HasPrefix(s.Str, "DATABASE_SNAPSHOT=") {
			stmt.SnapshotName = stmt.Source
		}
	}
	return nil
}

// parseBackupRestoreFileList parses backupFileListOpt / restoreFileListOpt:
// an optional list whose every comma is followed by another file.
//
// TSql160.g:5654-5730
//
//	backupRestoreFile:
//	    File EqualsSign ( stringOrVariable | backupRestoreFileNameList )
//	  | Identifier /* FILEGROUP */ EqualsSign ( stringOrVariable | backupRestoreFileNameList )
//	  | Identifier /* PAGE */ EqualsSign stringOrVariable
//	  | Identifier /* READ_WRITE_FILEGROUPS */
//	backupRestoreFileNameList:
//	    LeftParenthesis stringOrVariable ( Comma stringOrVariable )* RightParenthesis
func (p *Parser) parseBackupRestoreFileList(allowReadWriteFilegroups bool) (*nodes.List, error) {
	if !p.backupAtFile(allowReadWriteFilegroups) {
		return nil, nil
	}
	var specs []nodes.Node
	for {
		if !p.backupAtFile(allowReadWriteFilegroups) {
			return nil, p.syntaxErrorAtCur()
		}
		word := p.backupWord()
		p.advance()
		spec := word
		if word != "READ_WRITE_FILEGROUPS" {
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			val, err := p.parseFileSpecValue(word != "PAGE")
			if err != nil {
				return nil, err
			}
			spec += "=" + val
		}
		specs = append(specs, &nodes.String{Str: spec})
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ,
	}
	return &nodes.List{Items: specs}, nil
}

// backupAtFile reports whether the current token starts a backupRestoreFile.
func (p *Parser) backupAtFile(allowReadWriteFilegroups bool) bool {
	switch p.backupWord() {
	case "FILE", "FILEGROUP", "PAGE":
		return true
	case "READ_WRITE_FILEGROUPS":
		return allowReadWriteFilegroups
	}
	return false
}

// parseFileSpecValue parses the value after FILE= / FILEGROUP= / PAGE=:
// stringOrVariable, or for FILE and FILEGROUP a parenthesized list of them.
func (p *Parser) parseFileSpecValue(allowList bool) (string, error) {
	if p.cur.Type == '(' && allowList {
		p.advance() // consume (
		var parts []string
		for {
			s, err := p.backupStringOrVariable()
			if err != nil {
				return "", err
			}
			parts = append(parts, s)
			if p.cur.Type != ',' {
				break
			}
			p.advance() // consume ,
		}
		if _, err := p.expect(')'); err != nil {
			return "", err
		}
		return "(" + strings.Join(parts, ",") + ")", nil
	}
	return p.backupStringOrVariable()
}

// parseDeviceList parses a devList: one or more devices.
//
// TSql160.g:5769-5784
//
//	devList:    deviceInfo ( Comma deviceInfo )*
//	deviceInfo: identifierOrVariable
//	          | Identifier /* DISK | TAPE | URL | VIRTUAL_DEVICE | DATABASE_SNAPSHOT */ EqualsSign stringOrVariable
func (p *Parser) parseDeviceList() (*nodes.List, error) {
	var devs []nodes.Node
	for {
		dev, err := p.parseOneDevice()
		if err != nil {
			return nil, err
		}
		devs = append(devs, &nodes.String{Str: dev})
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ,
	}
	return &nodes.List{Items: devs}, nil
}

// parseOneDevice parses a single device. It returns "TYPE=path" for a
// physical device and the name for a logical one.
func (p *Parser) parseOneDevice() (string, error) {
	if word := p.backupWord(); backupDeviceTypes[word] && p.peekNext().Type == '=' {
		p.advance() // consume device type
		p.advance() // consume =
		path, err := p.backupStringOrVariable()
		if err != nil {
			return "", err
		}
		return word + "=" + path, nil
	}
	return p.backupIdentifierOrVariable()
}

// extractDevicePath extracts the path portion from a "TYPE=path" device string.
func extractDevicePath(dev string) string {
	if idx := strings.Index(dev, "="); idx >= 0 {
		return dev[idx+1:]
	}
	return dev
}

// parseBackupRestoreOptions parses the option list after WITH: one or more
// options, separated by commas.
//
// TSql160.g:5806-5975
//
//	backupOption:  backupEncryptionOption | Identifier
//	             | Identifier EqualsSign ( signedIntegerOrVariable | stringLiteral )
//	restoreOption: fileStreamRestoreOption | Identifier
//	             | Identifier EqualsSign ( stringOrVariable afterClause | signedInteger | stringOrVariable )
//	             | 'MOVE' stringOrVariable To stringOrVariable
//	             | File EqualsSign signedIntegerOrVariable
func (p *Parser) parseBackupRestoreOptions(backup bool) (*nodes.List, error) {
	var opts []nodes.Node
	for {
		opt, err := p.parseOneBackupRestoreOption(backup)
		if err != nil {
			return nil, err
		}
		opts = append(opts, opt)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ,
	}
	return &nodes.List{Items: opts}, nil
}

// parseOneBackupRestoreOption parses a single BACKUP (backup) or RESTORE
// WITH option.
func (p *Parser) parseOneBackupRestoreOption(backup bool) (*nodes.BackupRestoreOption, error) {
	optLoc := p.pos()
	name := p.backupWord()
	flags, values := restoreFlagOptions, restoreValueOptions
	if backup {
		flags, values = backupFlagOptions, backupValueOptions
	}

	switch {
	case backup && name == "ENCRYPTION":
		return p.parseBackupEncryptionOption()
	case backup && name == "COMPRESSION" && p.peekNext().Type == '(':
		return p.parseBackupCompressionOption()
	case !backup && name == "MOVE":
		return p.parseRestoreMoveOption()
	case !backup && name == "FILESTREAM":
		return p.parseRestoreFilestreamOption()
	case !backup && name == "FILE":
		// fileRestoreOption: File EqualsSign signedIntegerOrVariable
		p.advance() // consume FILE
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		var val string
		var err error
		if p.cur.Type == tokVARIABLE {
			val, err = p.backupStringOrVariable()
		} else {
			val, err = p.backupSignedInteger()
		}
		if err != nil {
			return nil, err
		}
		return &nodes.BackupRestoreOption{
			Name:  "FILE",
			Value: val,
			Loc:   nodes.Loc{Start: optLoc, End: p.prevEnd()},
		}, nil
	}

	if values[name] && p.peekNext().Type == '=' {
		p.advance() // consume option name
		p.advance() // consume =
		opt := &nodes.BackupRestoreOption{
			Name: name,
			Loc:  nodes.Loc{Start: optLoc, End: -1},
		}
		isString := p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST || p.cur.Type == tokVARIABLE
		val, err := p.backupOptionValue()
		if err != nil {
			return nil, err
		}
		opt.Value = val
		// afterClause: STOPATMARK and STOPBEFOREMARK take AFTER after a
		// string or variable (TSql160.g:5930 IsStopAtBeforeMarkRestoreOption).
		if !backup && isString && (name == "STOPATMARK" || name == "STOPBEFOREMARK") && p.cur.Type == kwAFTER {
			p.advance() // consume AFTER
			after, err := p.backupStringOrVariable()
			if err != nil {
				return nil, err
			}
			opt.Value += " AFTER " + after
		}
		opt.Loc.End = p.prevEnd()
		return opt, nil
	}

	if flags[name] {
		p.advance()
		return &nodes.BackupRestoreOption{
			Name: name,
			Loc:  nodes.Loc{Start: optLoc, End: p.prevEnd()},
		}, nil
	}
	return nil, p.syntaxErrorAtCur()
}

// parseBackupEncryptionOption parses the ENCRYPTION option.
//
// TSql160.g:5849, 5867, 5532
//
//	'ENCRYPTION' LeftParenthesis 'ALGORITHM' EqualsSign Identifier Comma
//	    'SERVER' ( 'CERTIFICATE' | 'ASYMMETRIC' Key ) EqualsSign identifier RightParenthesis
func (p *Parser) parseBackupEncryptionOption() (*nodes.BackupRestoreOption, error) {
	optLoc := p.pos()
	p.advance() // consume ENCRYPTION

	opt := &nodes.BackupRestoreOption{
		Name: "ENCRYPTION",
		Loc:  nodes.Loc{Start: optLoc, End: -1},
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	if err := p.backupExpectWord("ALGORITHM"); err != nil {
		return nil, err
	}
	if _, err := p.expect('='); err != nil {
		return nil, err
	}
	alg := p.backupWord()
	if !backupEncryptionAlgorithms[alg] {
		return nil, p.syntaxErrorAtCur()
	}
	opt.Algorithm = alg
	p.advance()
	if _, err := p.expect(','); err != nil {
		return nil, err
	}
	if err := p.backupExpectWord("SERVER"); err != nil {
		return nil, err
	}
	switch p.backupWord() {
	case "CERTIFICATE":
		p.advance()
		opt.EncryptorType = "SERVER CERTIFICATE"
	case "ASYMMETRIC":
		p.advance()
		if _, err := p.expect(kwKEY); err != nil {
			return nil, err
		}
		opt.EncryptorType = "ASYMMETRIC KEY"
	default:
		return nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect('='); err != nil {
		return nil, err
	}
	name, ok := p.parseIdentifier()
	if !ok {
		return nil, p.syntaxErrorAtCur()
	}
	opt.EncryptorName = name
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	opt.Loc.End = p.prevEnd()
	return opt, nil
}

// parseBackupCompressionOption parses SQL Server 2022's
// COMPRESSION ( ALGORITHM = { MS_XPRESS | QAT_DEFLATE } ), which TSql160.g
// predates (https://learn.microsoft.com/en-us/sql/t-sql/statements/backup-transact-sql).
// The algorithm is kept in Algorithm.
func (p *Parser) parseBackupCompressionOption() (*nodes.BackupRestoreOption, error) {
	optLoc := p.pos()
	p.advance() // consume COMPRESSION
	p.advance() // consume (
	if err := p.backupExpectWord("ALGORITHM"); err != nil {
		return nil, err
	}
	if _, err := p.expect('='); err != nil {
		return nil, err
	}
	alg := p.backupWord()
	if !backupCompressionAlgorithms[alg] {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.BackupRestoreOption{
		Name:      "COMPRESSION",
		Algorithm: alg,
		Loc:       nodes.Loc{Start: optLoc, End: p.prevEnd()},
	}, nil
}

// parseRestoreMoveOption parses MOVE stringOrVariable TO stringOrVariable
// (TSql160.g:5962).
func (p *Parser) parseRestoreMoveOption() (*nodes.BackupRestoreOption, error) {
	optLoc := p.pos()
	p.advance() // consume MOVE

	opt := &nodes.BackupRestoreOption{
		Name: "MOVE",
		Loc:  nodes.Loc{Start: optLoc, End: -1},
	}
	var err error
	if opt.MoveFrom, err = p.backupStringOrVariable(); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwTO); err != nil {
		return nil, err
	}
	if opt.MoveTo, err = p.backupStringOrVariable(); err != nil {
		return nil, err
	}
	opt.Loc.End = p.prevEnd()
	return opt, nil
}

// parseRestoreFilestreamOption parses FILESTREAM ( DIRECTORY_NAME = ... ).
// TSql160.g:5986 reuses the CREATE DATABASE FILESTREAM list, but the server
// takes only one DIRECTORY_NAME there, whose value is a string literal or
// NULL: it refuses NON_TRANSACTED_ACCESS, a second option and a variable.
func (p *Parser) parseRestoreFilestreamOption() (*nodes.BackupRestoreOption, error) {
	optLoc := p.pos()
	p.advance() // consume FILESTREAM

	opt := &nodes.BackupRestoreOption{
		Name: "FILESTREAM",
		Loc:  nodes.Loc{Start: optLoc, End: -1},
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	if err := p.backupExpectWord("DIRECTORY_NAME"); err != nil {
		return nil, err
	}
	if _, err := p.expect('='); err != nil {
		return nil, err
	}
	switch p.cur.Type {
	case tokSCONST, tokNSCONST:
		opt.Value = p.cur.Str
	case kwNULL:
		opt.Value = "NULL"
	default:
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	opt.Loc.End = p.prevEnd()
	return opt, nil
}
