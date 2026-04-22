// Package parser - backup_restore.go implements T-SQL BACKUP and RESTORE statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseBackupStmt parses a BACKUP DATABASE or BACKUP LOG statement.
//
// BNF: mssql/parser/bnf/backup-transact-sql.bnf
//
//	backupStatement:
//	    backupMain
//	    ( backupDevices )?
//	    ( backupOptions )?
//	    ;
//
//	backupMain:
//	    backupDatabase
//	    | backupTransactionLog
//	    ;
//
//	backupDatabase:
//	    Database identifierOrVariable
//	    backupFileListOpt
//	    ;
//
//	backupTransactionLog:
//	    Identifier /* LOG */ identifierOrVariable
//	    ;
//
//	backupFileListOpt:
//	    ( backupRestoreFile ( Comma backupRestoreFile )* )?
//	    ;
//
//	backupRestoreFile:
//	    File EqualsSign ( stringOrVariable | backupRestoreFileNameList )
//	    | Identifier /* FILEGROUP | PAGE */ EqualsSign ( stringOrVariable | backupRestoreFileNameList )
//	    | Identifier /* READ_WRITE_FILEGROUPS */
//	    ;
//
//	backupRestoreFileNameList:
//	    LeftParenthesis stringOrVariable ( Comma stringOrVariable )* RightParenthesis
//	    ;
//
//	backupDevices:
//	    To devList ( mirrorTo )*
//	    ;
//
//	mirrorTo:
//	    Identifier /* MIRROR */ To devList
//	    ;
//
//	devList:
//	    deviceInfo ( Comma deviceInfo )*
//	    ;
//
//	deviceInfo:
//	    identifierOrVariable
//	    | Identifier /* DISK | TAPE | URL */ EqualsSign stringOrVariable
//	    ;
//
//	backupOptions:
//	    With backupOption ( Comma backupOption )*
//	    ;
//
//	backupOption:
//	    backupEncryptionOption
//	    | Identifier
//	    | Identifier EqualsSign ( signedIntegerOrVariable | stringLiteral )
//	    ;
//
//	backupEncryptionOption:
//	    Identifier /* ENCRYPTION */ LeftParenthesis
//	        Identifier /* ALGORITHM */ EqualsSign Identifier /* algorithm kind */
//	        Comma backupEncryptor
//	    RightParenthesis
//	    ;
//
//	backupEncryptor:
//	    dekEncryptorType EqualsSign identifier
//	    ;
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

	// DATABASE or LOG (or identifier like LOG/CERTIFICATE)
	if p.cur.Type == kwDATABASE {
		stmt.Type = "DATABASE"
		p.advance()
	} else if p.cur.Type == kwLOG {
		stmt.Type = "LOG"
		p.advance()
	} else if p.cur.Type == kwCERTIFICATE {
		stmt.Type = "CERTIFICATE"
		p.advance()
	} else {
		stmt.Type = "DATABASE"
	}

	// Completion: after BACKUP DATABASE → database_ref
	if p.collectMode() {
		p.addRuleCandidate("database_ref")
		return nil, errCollecting
	}

	// Database name (not for CERTIFICATE)
	if stmt.Type != "CERTIFICATE" {
		if p.isIdentLike() {
			stmt.Database = p.cur.Str
			p.advance()
		}
	}

	// backupFileListOpt: FILE = ..., FILEGROUP = ..., READ_WRITE_FILEGROUPS (only for DATABASE)
	if stmt.Type == "DATABASE" {
		stmt.FileSpecs, _ = p.parseBackupRestoreFileList()
	}

	// TO devList [ MIRROR TO devList ]*
	if p.cur.Type == kwTO {
		p.advance() // consume TO
		stmt.Devices, _ = p.parseDeviceList()
		// Extract first device path for backward compat Target field
		if stmt.Devices != nil && len(stmt.Devices.Items) > 0 {
			if s, ok := stmt.Devices.Items[0].(*nodes.String); ok {
				// Extract path from "TYPE=path" format or use as-is
				stmt.Target = extractDevicePath(s.Str)
			}
		}

		// MIRROR TO devList (may repeat)
		for p.cur.Type == kwMIRROR {
			stmt.MirrorTo = true
			p.advance() // consume MIRROR
			if p.cur.Type == kwTO {
				p.advance() // consume TO
			}
			mirrorDevs, _ := p.parseDeviceList()
			if stmt.MirrorDevice == "" && mirrorDevs != nil && len(mirrorDevs.Items) > 0 {
				if s, ok := mirrorDevs.Items[0].(*nodes.String); ok {
					stmt.MirrorDevice = extractDevicePath(s.Str)
				}
			}
		}
	}

	// WITH options — structured parsing
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		var err error
		stmt.Options, err = p.parseBackupRestoreOptions()
		if err != nil {
			stmt.Loc.End = p.prevEnd()
			return stmt, err
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseRestoreStmt parses a RESTORE DATABASE / LOG / HEADERONLY / FILELISTONLY statement.
//
// BNF: mssql/parser/bnf/restore-statements-transact-sql.bnf
//
//	restoreStatement:
//	    (
//	        (
//	            restoreMain
//	            (From devList)?
//	        )
//	    |
//	        (
//	            Identifier /* HEADERONLY | FILELISTONLY | VERIFYONLY | LABELONLY | REWINDONLY */
//	            From devList
//	        )
//	    )
//	    (restoreOptions)?
//	    ;
//
//	restoreMain:
//	    Database identifierOrVariable restoreFileListOpt
//	    | Identifier /* LOG */ identifierOrVariable restoreFileListOpt
//	    ;
//
//	restoreFileListOpt:
//	    (backupRestoreFile (Comma backupRestoreFile)*)?
//	    ;
//
//	backupRestoreFile:
//	    File EqualsSign (stringOrVariable | backupRestoreFileNameList)
//	    | Identifier /* FILEGROUP | PAGE */ EqualsSign (stringOrVariable | backupRestoreFileNameList)
//	    ;
//
//	devList:
//	    deviceInfo (Comma deviceInfo)*
//	    ;
//
//	deviceInfo:
//	    identifierOrVariable
//	    | Identifier /* DISK | TAPE | URL | VIRTUAL_DEVICE */ EqualsSign stringOrVariable
//	    ;
//
//	restoreOptions:
//	    With restoreOptionsList
//	    ;
//
//	restoreOptionsList:
//	    restoreOption (Comma restoreOption)*
//	    ;
//
//	restoreOption:
//	    fileStreamRestoreOption
//	    | simpleRestoreOption
//	    | Identifier EqualsSign (stringOrVariable | signedInteger)
//	    | moveRestoreOption
//	    | fileRestoreOption
//	    ;
//
//	-- Also: RESTORE DATABASE db FROM DATABASE_SNAPSHOT = snapshot_name
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

	// Determine restore type
	if p.cur.Type == kwDATABASE {
		stmt.Type = "DATABASE"
		p.advance()
	} else if p.isAnyKeywordIdent() {
		upper := strings.ToUpper(p.cur.Str)
		switch upper {
		case "LOG":
			stmt.Type = "LOG"
			p.advance()
		case "HEADERONLY":
			stmt.Type = "HEADERONLY"
			p.advance()
		case "FILELISTONLY":
			stmt.Type = "FILELISTONLY"
			p.advance()
		case "VERIFYONLY":
			stmt.Type = "VERIFYONLY"
			p.advance()
		case "LABELONLY":
			stmt.Type = "LABELONLY"
			p.advance()
		case "REWINDONLY":
			stmt.Type = "REWINDONLY"
			p.advance()
		default:
			stmt.Type = "DATABASE"
		}
	} else {
		stmt.Type = "DATABASE"
	}

	// Completion: after RESTORE DATABASE → database_ref
	if p.collectMode() {
		p.addRuleCandidate("database_ref")
		return nil, errCollecting
	}

	// Database name (optional for HEADERONLY/FILELISTONLY/VERIFYONLY/LABELONLY)
	switch stmt.Type {
	case "HEADERONLY", "FILELISTONLY", "VERIFYONLY", "LABELONLY", "REWINDONLY":
		// no database name expected before FROM
	default:
		if p.isIdentLike() && p.cur.Type != kwFROM {
			stmt.Database = p.cur.Str
			p.advance()
		}
	}

	// File/filegroup list for DATABASE/LOG restores
	if stmt.Type == "DATABASE" || stmt.Type == "LOG" {
		stmt.FileSpecs, _ = p.parseBackupRestoreFileList()
	}

	// FROM clause
	if p.cur.Type == kwFROM {
		p.advance() // consume FROM

		// DATABASE_SNAPSHOT = snapshot_name
		if p.cur.Type == kwDATABASE_SNAPSHOT {
			p.advance() // consume DATABASE_SNAPSHOT
			if _, ok := p.match('='); ok {
				if p.isIdentLike() || p.cur.Type == tokSCONST {
					stmt.SnapshotName = p.cur.Str
					p.advance()
				}
			}
		} else {
			// devList: deviceInfo [ , deviceInfo ]*
			stmt.Devices, _ = p.parseDeviceList()
			if stmt.Devices != nil && len(stmt.Devices.Items) > 0 {
				if s, ok := stmt.Devices.Items[0].(*nodes.String); ok {
					stmt.Source = extractDevicePath(s.Str)
				}
			}
		}
	}

	// WITH options — structured parsing
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		var err error
		stmt.Options, err = p.parseBackupRestoreOptions()
		if err != nil {
			stmt.Loc.End = p.prevEnd()
			return stmt, err
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseBackupRestoreFileList parses the optional file/filegroup list.
//
//	backupRestoreFile:
//	    FILE = ( stringOrVariable | backupRestoreFileNameList )
//	    | Identifier /* FILEGROUP | PAGE */ = ( stringOrVariable | backupRestoreFileNameList )
//	    | Identifier /* READ_WRITE_FILEGROUPS */
//	    ;
func (p *Parser) parseBackupRestoreFileList() (*nodes.List, error) {
	var specs []nodes.Node

	for {
		if p.cur.Type == kwFILE {
			// FILE = value
			spec := "FILE="
			p.advance() // consume FILE
			if _, ok := p.match('='); ok {
				spec += p.parseFileSpecValue()
			}
			specs = append(specs, &nodes.String{Str: spec})
		} else if p.cur.Type == kwFILEGROUP {
			spec := "FILEGROUP="
			p.advance() // consume FILEGROUP
			if _, ok := p.match('='); ok {
				spec += p.parseFileSpecValue()
			}
			specs = append(specs, &nodes.String{Str: spec})
		} else if p.cur.Type == kwPAGE {
			spec := "PAGE="
			p.advance() // consume PAGE
			if _, ok := p.match('='); ok {
				spec += p.parseFileSpecValue()
			}
			specs = append(specs, &nodes.String{Str: spec})
		} else if p.cur.Type == kwREAD_WRITE_FILEGROUPS {
			specs = append(specs, &nodes.String{Str: "READ_WRITE_FILEGROUPS"})
			p.advance()
		} else {
			break
		}

		// comma to continue
		if p.cur.Type == ',' {
			// Peek: is this another file spec or the start of TO/FROM?
			// If next token after comma is FILE, FILEGROUP, PAGE, READ_WRITE_FILEGROUPS, consume comma
			next := p.peekNext()
			if next.Type == kwFILE ||
				next.Type == kwFILEGROUP ||
				next.Type == kwPAGE ||
				next.Type == kwREAD_WRITE_FILEGROUPS {
				p.advance() // consume comma
			} else {
				break
			}
		} else {
			break
		}
	}

	if len(specs) == 0 {
		return nil, nil
	}
	return &nodes.List{Items: specs}, nil
}

// parseFileSpecValue parses the value after FILE= / FILEGROUP= / PAGE=.
// This can be a string/variable or a parenthesized list.
func (p *Parser) parseFileSpecValue() string {
	if p.cur.Type == '(' {
		// backupRestoreFileNameList: ( stringOrVariable [ , stringOrVariable ]* )
		p.advance() // consume (
		var parts []string
		for {
			if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
				parts = append(parts, p.cur.Str)
				p.advance()
			} else if p.isAnyKeywordIdent() {
				parts = append(parts, p.cur.Str)
				p.advance()
			} else {
				break
			}
			if p.cur.Type == ',' {
				p.advance()
			} else {
				break
			}
		}
		if p.cur.Type == ')' {
			p.advance()
		}
		return "(" + strings.Join(parts, ",") + ")"
	}
	// stringOrVariable
	if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
		val := p.cur.Str
		p.advance()
		return val
	}
	if p.isAnyKeywordIdent() {
		val := p.cur.Str
		p.advance()
		return val
	}
	return ""
}

// parseDeviceList parses a comma-separated list of backup/restore devices.
//
//	devList: deviceInfo ( Comma deviceInfo )*
//
//	deviceInfo:
//	    identifierOrVariable
//	    | Identifier /* DISK | TAPE | URL | VIRTUAL_DEVICE */ EqualsSign stringOrVariable
func (p *Parser) parseDeviceList() (*nodes.List, error) {
	var devs []nodes.Node

	for {
		dev := p.parseOneDevice()
		if dev == "" {
			break
		}
		devs = append(devs, &nodes.String{Str: dev})

		if p.cur.Type == ',' {
			p.advance()
		} else {
			break
		}
	}

	if len(devs) == 0 {
		return nil, nil
	}
	return &nodes.List{Items: devs}, nil
}

// parseOneDevice parses a single device entry.
// Returns "TYPE=path" for physical devices or "logical_name" for logical devices.
func (p *Parser) parseOneDevice() string {
	if !p.isAnyKeywordIdent() && p.cur.Type != kwFILE {
		return ""
	}

	// Check if this is a known device type keyword followed by =
	upper := strings.ToUpper(p.cur.Str)
	switch upper {
	case "DISK", "TAPE", "URL", "VIRTUAL_DEVICE":
		devType := upper
		p.advance() // consume type keyword
		if _, ok := p.match('='); ok {
			if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
				path := p.cur.Str
				p.advance()
				return devType + "=" + path
			}
			if p.isAnyKeywordIdent() {
				// variable like @path_var
				path := p.cur.Str
				p.advance()
				return devType + "=" + path
			}
		}
		return devType
	default:
		// Logical device name
		name := p.cur.Str
		p.advance()
		// Check if followed by = (could be a device type we don't know)
		if p.cur.Type == '=' {
			p.advance()
			if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
				path := p.cur.Str
				p.advance()
				return name + "=" + path
			}
			if p.isAnyKeywordIdent() {
				path := p.cur.Str
				p.advance()
				return name + "=" + path
			}
		}
		return name
	}
}

// extractDevicePath extracts the path portion from a "TYPE=path" device string.
func extractDevicePath(dev string) string {
	if idx := strings.Index(dev, "="); idx >= 0 {
		return dev[idx+1:]
	}
	return dev
}

// parseBackupRestoreOptions parses structured BACKUP/RESTORE WITH options.
// Called after WITH has been consumed.
//
//	options ::= option [ , option ] ...
//	option ::=
//	    COMPRESSION | NO_COMPRESSION | DIFFERENTIAL | COPY_ONLY
//	  | INIT | NOINIT | NOSKIP | SKIP | FORMAT | NOFORMAT
//	  | NO_CHECKSUM | CHECKSUM | STOP_ON_ERROR | CONTINUE_AFTER_ERROR
//	  | RESTART | REPLACE | RECOVERY | NORECOVERY | NO_TRUNCATE | FILE_SNAPSHOT
//	  | ENABLE_BROKER | NEW_BROKER | ERROR_BROKER_CONVERSATIONS
//	  | REWIND | NOREWIND | UNLOAD | NOUNLOAD
//	  | RESTRICTED_USER | KEEP_REPLICATION | KEEP_CDC
//	  | PARTIAL | CREDENTIAL | METADATA_ONLY | SNAPSHOT
//	  | NAME = { 'name' | @var }
//	  | DESCRIPTION = { 'text' | @var }
//	  | EXPIREDATE = { 'date' | @var }
//	  | RETAINDAYS = { days | @var }
//	  | STATS [ = percentage ]
//	  | BLOCKSIZE = { n | @var }
//	  | BUFFERCOUNT = { n | @var }
//	  | MAXTRANSFERSIZE = { n | @var }
//	  | MEDIADESCRIPTION = { 'text' | @var }
//	  | MEDIANAME = { 'name' | @var }
//	  | MEDIAPASSWORD = { 'password' | @var }
//	  | PASSWORD = { 'password' | @var }
//	  | STANDBY = standby_file_name
//	  | STOPAT = { 'datetime' | @var }
//	  | STOPATMARK = { 'mark' } [ AFTER 'datetime' ]
//	  | STOPBEFOREMARK = { 'mark' } [ AFTER 'datetime' ]
//	  | FILE = { n | @var }
//	  | DBNAME = { database_name | @var }
//	  | ENCRYPTION ( ALGORITHM = alg, SERVER { CERTIFICATE | ASYMMETRIC KEY } = name )
//	  | MOVE 'logical_file_name' TO 'os_file_name'
//	  | FILESTREAM ( DIRECTORY_NAME = directory_name )
func (p *Parser) parseBackupRestoreOptions() (*nodes.List, error) {
	var opts []nodes.Node

	// After WITH, at least one valid option is required.
	if !p.isBackupRestoreOption() {
		return nil, p.syntaxErrorAtCur()
	}

	for {
		if p.cur.Type == tokEOF || p.cur.Type == ';' || !p.isBackupRestoreOption() {
			break
		}

		opt, _ := p.parseOneBackupRestoreOption()
		if opt != nil {
			opts = append(opts, opt)
		}

		if p.cur.Type == ',' {
			p.advance()
		} else {
			break
		}
	}

	if len(opts) == 0 {
		return nil, nil
	}
	return &nodes.List{Items: opts}, nil
}

// backupRestoreFlagOptions lists option names that take no value (bare flags).
var backupRestoreFlagOptions = map[string]bool{
	"COMPRESSION": true, "NO_COMPRESSION": true,
	"DIFFERENTIAL": true, "COPY_ONLY": true,
	"INIT": true, "NOINIT": true,
	"NOSKIP": true, "SKIP": true,
	"FORMAT": true, "NOFORMAT": true,
	"NO_CHECKSUM": true, "CHECKSUM": true,
	"STOP_ON_ERROR": true, "CONTINUE_AFTER_ERROR": true,
	"RESTART": true, "REPLACE": true,
	"RECOVERY": true, "NORECOVERY": true,
	"NO_TRUNCATE": true, "FILE_SNAPSHOT": true,
	"ENABLE_BROKER": true, "NEW_BROKER": true,
	"ERROR_BROKER_CONVERSATIONS": true,
	"REWIND": true, "NOREWIND": true,
	"UNLOAD": true, "NOUNLOAD": true,
	"RESTRICTED_USER": true,
	"KEEP_REPLICATION": true, "KEEP_CDC": true,
	"PARTIAL": true, "CREDENTIAL": true,
	"METADATA_ONLY": true, "SNAPSHOT": true,
}

// backupRestoreKVOptions lists option names that take = value.
var backupRestoreKVOptions = map[string]bool{
	"NAME": true, "DESCRIPTION": true,
	"EXPIREDATE": true, "RETAINDAYS": true,
	"BLOCKSIZE": true, "BUFFERCOUNT": true, "MAXTRANSFERSIZE": true,
	"MEDIADESCRIPTION": true, "MEDIANAME": true, "MEDIAPASSWORD": true,
	"PASSWORD": true,
	"STANDBY": true, "STOPAT": true,
	"STOPATMARK": true, "STOPBEFOREMARK": true,
	"FILE": true, "DBNAME": true,
}

// backupRestoreSpecialOptions lists option names handled by dedicated parsers.
var backupRestoreSpecialOptions = map[string]bool{
	"ENCRYPTION": true, "MOVE": true, "FILESTREAM": true, "STATS": true,
}

// isBackupRestoreOption checks whether the current token is a valid backup/restore option name.
func (p *Parser) isBackupRestoreOption() bool {
	if !p.isAnyKeywordIdent() {
		return false
	}
	name := strings.ToUpper(p.cur.Str)
	return backupRestoreFlagOptions[name] || backupRestoreKVOptions[name] || backupRestoreSpecialOptions[name]
}

// parseOneBackupRestoreOption parses a single BACKUP/RESTORE WITH option.
func (p *Parser) parseOneBackupRestoreOption() (*nodes.BackupRestoreOption, error) {
	if !p.isBackupRestoreOption() {
		return nil, nil
	}

	optLoc := p.pos()
	name := strings.ToUpper(p.cur.Str)

	// ENCRYPTION ( ALGORITHM = ..., SERVER CERTIFICATE|ASYMMETRIC KEY = ... )
	if name == "ENCRYPTION" {
		opt, err := p.parseBackupEncryptionOption()
		return opt, err
	}

	// MOVE 'logical' TO 'physical'
	if name == "MOVE" {
		opt, err := p.parseRestoreMoveOption()
		return opt, err
	}

	// FILESTREAM ( DIRECTORY_NAME = directory_name )
	if name == "FILESTREAM" {
		opt, err := p.parseRestoreFilestreamOption()
		return opt, err
	}

	// STATS [ = percentage ] — special: '=' is optional
	if name == "STATS" {
		p.advance() // consume STATS
		opt := &nodes.BackupRestoreOption{
			Name: "STATS",
			Loc:  nodes.Loc{Start: optLoc, End: -1},
		}
		if p.cur.Type == '=' {
			p.advance()
			if p.cur.Type == tokICONST || p.cur.Type == tokFCONST ||
				p.cur.Type == tokSCONST || (p.cur.Type == tokIDENT && p.cur.Str[0] == '@') {
				opt.Value = p.cur.Str
				p.advance()
			}
		}
		opt.Loc.End = p.prevEnd()
		return opt, nil
	}

	// Flag options (no value)
	if backupRestoreFlagOptions[name] {
		p.advance()
		return &nodes.BackupRestoreOption{
			Name: name,
			Loc:  nodes.Loc{Start: optLoc, End: p.prevEnd()},
		}, nil
	}

	// Key = value options
	if backupRestoreKVOptions[name] {
		p.advance() // consume option name
		opt := &nodes.BackupRestoreOption{
			Name: name,
			Loc:  nodes.Loc{Start: optLoc, End: -1},
		}
		if _, ok := p.match('='); ok {
			// Value can be string constant, number, or variable
			if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST ||
				p.cur.Type == tokICONST || p.cur.Type == tokFCONST {
				opt.Value = p.cur.Str
				p.advance()
			} else if p.isIdentLike() {
				opt.Value = p.cur.Str
				p.advance()
			}
			// For STOPATMARK / STOPBEFOREMARK: optional AFTER 'datetime'
			if (name == "STOPATMARK" || name == "STOPBEFOREMARK") &&
				p.cur.Type == kwAFTER {
				p.advance()
				if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
					opt.Value = opt.Value + " AFTER " + p.cur.Str
					p.advance()
				}
			}
		}
		opt.Loc.End = p.prevEnd()
		return opt, nil
	}

	// Unknown option — consume name and optional = value structurally
	p.advance() // consume option name
	opt := &nodes.BackupRestoreOption{
		Name: name,
		Loc:  nodes.Loc{Start: optLoc, End: -1},
	}
	if p.cur.Type == '=' {
		p.advance()
		if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST ||
			p.cur.Type == tokICONST || p.cur.Type == tokFCONST {
			opt.Value = p.cur.Str
			p.advance()
		} else if p.isIdentLike() {
			opt.Value = p.cur.Str
			p.advance()
		}
	}
	opt.Loc.End = p.prevEnd()
	return opt, nil
}

// parseBackupEncryptionOption parses the ENCRYPTION option.
//
//	ENCRYPTION ( ALGORITHM = { AES_128 | AES_192 | AES_256 | TRIPLE_DES_3KEY },
//	    SERVER CERTIFICATE = cert_name | SERVER ASYMMETRIC KEY = key_name )
func (p *Parser) parseBackupEncryptionOption() (*nodes.BackupRestoreOption, error) {
	optLoc := p.pos()
	p.advance() // consume ENCRYPTION

	opt := &nodes.BackupRestoreOption{
		Name: "ENCRYPTION",
		Loc:  nodes.Loc{Start: optLoc, End: -1},
	}

	if p.cur.Type != '(' {
		opt.Loc.End = p.prevEnd()
		return opt, nil
	}
	p.advance() // consume (

	// ALGORITHM = alg
	if p.cur.Type == kwALGORITHM {
		p.advance() // consume ALGORITHM
		if _, ok := p.match('='); ok {
			if p.isAnyKeywordIdent() {
				opt.Algorithm = strings.ToUpper(p.cur.Str)
				p.advance()
			}
		}
	}

	// comma separator
	if p.cur.Type == ',' {
		p.advance()
	}

	// SERVER CERTIFICATE = name | SERVER ASYMMETRIC KEY = name
	if p.cur.Type == kwSERVER {
		p.advance() // consume SERVER
		if p.isAnyKeywordIdent() {
			upper := strings.ToUpper(p.cur.Str)
			if upper == "CERTIFICATE" {
				opt.EncryptorType = "SERVER CERTIFICATE"
				p.advance()
			} else if upper == "ASYMMETRIC" {
				p.advance() // consume ASYMMETRIC
				if p.cur.Type == kwKEY {
					p.advance() // consume KEY
				}
				opt.EncryptorType = "ASYMMETRIC KEY"
			}
		}
		// = name
		if p.cur.Type == '=' {
			p.advance()
		}
		if p.isIdentLike() {
			opt.EncryptorName = p.cur.Str
			p.advance()
		}
	}

	// closing paren
	if p.cur.Type == ')' {
		p.advance()
	}

	opt.Loc.End = p.prevEnd()
	return opt, nil
}

// parseRestoreMoveOption parses MOVE 'logical_file_name' TO 'os_file_name'.
func (p *Parser) parseRestoreMoveOption() (*nodes.BackupRestoreOption, error) {
	optLoc := p.pos()
	p.advance() // consume MOVE

	opt := &nodes.BackupRestoreOption{
		Name: "MOVE",
		Loc:  nodes.Loc{Start: optLoc, End: -1},
	}

	// 'logical_file_name'
	if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
		opt.MoveFrom = p.cur.Str
		p.advance()
	}

	// TO
	if p.cur.Type == kwTO {
		p.advance()
	}

	// 'os_file_name'
	if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
		opt.MoveTo = p.cur.Str
		p.advance()
	}

	opt.Loc.End = p.prevEnd()
	return opt, nil
}

// parseRestoreFilestreamOption parses FILESTREAM ( DIRECTORY_NAME = directory_name ).
func (p *Parser) parseRestoreFilestreamOption() (*nodes.BackupRestoreOption, error) {
	optLoc := p.pos()
	p.advance() // consume FILESTREAM

	opt := &nodes.BackupRestoreOption{
		Name: "FILESTREAM",
		Loc:  nodes.Loc{Start: optLoc, End: -1},
	}

	if p.cur.Type != '(' {
		opt.Loc.End = p.prevEnd()
		return opt, nil
	}
	p.advance() // consume (

	// DIRECTORY_NAME = directory_name
	if p.cur.Type == kwDIRECTORY_NAME {
		p.advance() // consume DIRECTORY_NAME
		if _, ok := p.match('='); ok {
			if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
				opt.Value = p.cur.Str
				p.advance()
			} else if p.isIdentLike() {
				opt.Value = p.cur.Str
				p.advance()
			}
		}
	}

	if p.cur.Type == ')' {
		p.advance()
	}

	opt.Loc.End = p.prevEnd()
	return opt, nil
}
