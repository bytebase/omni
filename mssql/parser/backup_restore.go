// Package parser - backup_restore.go implements T-SQL BACKUP and RESTORE statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseBackupStmt parses a BACKUP DATABASE or BACKUP LOG statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/backup-transact-sql
//
//	BACKUP { DATABASE | LOG } database_name
//	    TO { DISK | URL | TAPE } = 'path' [ , ...n ]
//	    [ <mirror_to_clause> ] [ ...next-mirror-to ]
//	    [ WITH { <general_WITH_options> | <backup_set_WITH_options> } [ ,...n ] ]
//
//	<general_WITH_options> ::=
//	    COPY_ONLY
//	  | { COMPRESSION | NO_COMPRESSION }
//	  | DESCRIPTION = { 'text' | @text_variable }
//	  | NAME = { backup_set_name | @backup_set_name_var }
//	  | CREDENTIAL
//	  | ENCRYPTION ( ALGORITHM = { AES_128 | AES_192 | AES_256 | TRIPLE_DES_3KEY },
//	        SERVER CERTIFICATE = cert_name | SERVER ASYMMETRIC KEY = key_name )
//	  | FILE_SNAPSHOT
//	  | { EXPIREDATE = { 'date' | @date_var } | RETAINDAYS = { days | @days_var } }
//	  | { NOINIT | INIT }
//	  | { NOSKIP | SKIP }
//	  | { NOFORMAT | FORMAT }
//	  | MEDIADESCRIPTION = { 'text' | @text_variable }
//	  | MEDIANAME = { media_name | @media_name_variable }
//	  | BLOCKSIZE = { blocksize | @blocksize_variable }
//	  | BUFFERCOUNT = { buffercount | @buffercount_variable }
//	  | MAXTRANSFERSIZE = { maxtransfersize | @maxtransfersize_variable }
//	  | { NO_CHECKSUM | CHECKSUM }
//	  | { STOP_ON_ERROR | CONTINUE_AFTER_ERROR }
//	  | RESTART
//	  | STATS [ = percentage ]
//	  | { REWIND | NOREWIND }
//	  | { UNLOAD | NOUNLOAD }
//	  | NORECOVERY
//	  | STANDBY = standby_file_name
//	  | NO_TRUNCATE
//	  | DIFFERENTIAL
func (p *Parser) parseBackupStmt() *nodes.BackupStmt {
	loc := p.pos()
	p.advance() // consume BACKUP

	stmt := &nodes.BackupStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// DATABASE or LOG (or identifier like LOG)
	if p.cur.Type == kwDATABASE {
		stmt.Type = "DATABASE"
		p.advance()
	} else if p.isIdentLike() && strings.EqualFold(p.cur.Str, "LOG") {
		stmt.Type = "LOG"
		p.advance()
	} else if p.isIdentLike() && strings.EqualFold(p.cur.Str, "CERTIFICATE") {
		stmt.Type = "CERTIFICATE"
		p.advance()
	} else {
		stmt.Type = "DATABASE"
	}

	// Database name (not for CERTIFICATE)
	if stmt.Type != "CERTIFICATE" {
		if p.isIdentLike() {
			stmt.Database = p.cur.Str
			p.advance()
		}
	}

	// TO { DISK | URL | TAPE | ... } = 'path'
	if p.cur.Type == kwTO {
		p.advance() // consume TO
		// consume DISK / URL / TAPE / identifier
		if p.isIdentLike() || p.cur.Type == kwFILE {
			p.advance()
		}
		// = 'path'
		if _, ok := p.match('='); ok {
			if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
				stmt.Target = p.cur.Str
				p.advance()
			}
		}
	}

	// WITH options — structured parsing
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		stmt.Options = p.parseBackupRestoreOptions()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseRestoreStmt parses a RESTORE DATABASE / LOG / HEADERONLY / FILELISTONLY statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/restore-statements-transact-sql
//
//	RESTORE { DATABASE | LOG } database_name
//	    FROM { DISK | URL | TAPE } = 'path' [ , ...n ]
//	    [ WITH
//	        { MOVE 'logical_file_name' TO 'operating_system_file_name' } [ ,...n ]
//	      | REPLACE
//	      | { RECOVERY | NORECOVERY | STANDBY = standby_file_name }
//	      | STOPAT = { 'datetime' | @datetime_var }
//	      | STOPATMARK = { 'lsn:lsn_number' | 'mark_name' } [ AFTER 'datetime' ]
//	      | STOPBEFOREMARK = { 'lsn:lsn_number' | 'mark_name' } [ AFTER 'datetime' ]
//	      | FILE = { backup_set_file_number | @backup_set_file_number }
//	      | MEDIANAME = { media_name | @media_name_variable }
//	      | MEDIAPASSWORD = { mediapassword | @mediapassword_variable }
//	      | ENABLE_BROKER
//	      | NEW_BROKER
//	      | ERROR_BROKER_CONVERSATIONS
//	      | { NO_CHECKSUM | CHECKSUM }
//	      | { STOP_ON_ERROR | CONTINUE_AFTER_ERROR }
//	      | STATS [ = percentage ]
//	      | { REWIND | NOREWIND }
//	      | { UNLOAD | NOUNLOAD }
//	      | RESTRICTED_USER
//	      | KEEP_REPLICATION
//	      | KEEP_CDC
//	      | BUFFERCOUNT = buffercount
//	      | MAXTRANSFERSIZE = maxtransfersize
//	      | BLOCKSIZE = blocksize
//	    ]
//
//	RESTORE { HEADERONLY | FILELISTONLY | VERIFYONLY | LABELONLY | REWINDONLY }
//	    FROM { DISK | URL | TAPE } = 'path'
//	    [ WITH options ]
func (p *Parser) parseRestoreStmt() *nodes.RestoreStmt {
	loc := p.pos()
	p.advance() // consume RESTORE

	stmt := &nodes.RestoreStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Determine restore type
	if p.cur.Type == kwDATABASE {
		stmt.Type = "DATABASE"
		p.advance()
	} else if p.isIdentLike() {
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

	// FROM { DISK | URL | TAPE | ... } = 'path'
	if p.cur.Type == kwFROM {
		p.advance() // consume FROM
		// consume DISK / URL / TAPE / identifier
		if p.isIdentLike() || p.cur.Type == kwFILE {
			p.advance()
		}
		// = 'path'
		if _, ok := p.match('='); ok {
			if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
				stmt.Source = p.cur.Str
				p.advance()
			}
		}
	}

	// WITH options — structured parsing
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		stmt.Options = p.parseBackupRestoreOptions()
	}

	stmt.Loc.End = p.pos()
	return stmt
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
//	  | STANDBY = standby_file_name
//	  | STOPAT = { 'datetime' | @var }
//	  | STOPATMARK = { 'mark' } [ AFTER 'datetime' ]
//	  | STOPBEFOREMARK = { 'mark' } [ AFTER 'datetime' ]
//	  | FILE = { n | @var }
//	  | ENCRYPTION ( ALGORITHM = alg, SERVER { CERTIFICATE | ASYMMETRIC KEY } = name )
//	  | MOVE 'logical_file_name' TO 'os_file_name'
func (p *Parser) parseBackupRestoreOptions() *nodes.List {
	var opts []nodes.Node

	for {
		if p.cur.Type == tokEOF || p.cur.Type == ';' || !p.isIdentLike() {
			break
		}

		opt := p.parseOneBackupRestoreOption()
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
		return nil
	}
	return &nodes.List{Items: opts}
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
}

// backupRestoreKVOptions lists option names that take = value.
var backupRestoreKVOptions = map[string]bool{
	"NAME": true, "DESCRIPTION": true,
	"EXPIREDATE": true, "RETAINDAYS": true,
	"BLOCKSIZE": true, "BUFFERCOUNT": true, "MAXTRANSFERSIZE": true,
	"MEDIADESCRIPTION": true, "MEDIANAME": true, "MEDIAPASSWORD": true,
	"STANDBY": true, "STOPAT": true,
	"STOPATMARK": true, "STOPBEFOREMARK": true,
	"FILE": true, "CREDENTIAL": true,
}

// parseOneBackupRestoreOption parses a single BACKUP/RESTORE WITH option.
func (p *Parser) parseOneBackupRestoreOption() *nodes.BackupRestoreOption {
	if !p.isIdentLike() {
		return nil
	}

	optLoc := p.pos()
	name := strings.ToUpper(p.cur.Str)

	// ENCRYPTION ( ALGORITHM = ..., SERVER CERTIFICATE|ASYMMETRIC KEY = ... )
	if name == "ENCRYPTION" {
		return p.parseBackupEncryptionOption()
	}

	// MOVE 'logical' TO 'physical'
	if name == "MOVE" {
		return p.parseRestoreMoveOption()
	}

	// STATS [ = percentage ] — special: '=' is optional
	if name == "STATS" {
		p.advance() // consume STATS
		opt := &nodes.BackupRestoreOption{
			Name: "STATS",
			Loc:  nodes.Loc{Start: optLoc},
		}
		if p.cur.Type == '=' {
			p.advance()
			if p.cur.Type == tokICONST || p.cur.Type == tokFCONST ||
				p.cur.Type == tokSCONST || (p.cur.Type == tokIDENT && p.cur.Str[0] == '@') {
				opt.Value = p.cur.Str
				p.advance()
			}
		}
		opt.Loc.End = p.pos()
		return opt
	}

	// Flag options (no value)
	if backupRestoreFlagOptions[name] {
		p.advance()
		return &nodes.BackupRestoreOption{
			Name: name,
			Loc:  nodes.Loc{Start: optLoc, End: p.pos()},
		}
	}

	// Key = value options
	if backupRestoreKVOptions[name] {
		p.advance() // consume option name
		opt := &nodes.BackupRestoreOption{
			Name: name,
			Loc:  nodes.Loc{Start: optLoc},
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
				p.isIdentLike() && strings.EqualFold(p.cur.Str, "AFTER") {
				p.advance()
				if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
					opt.Value = opt.Value + " AFTER " + p.cur.Str
					p.advance()
				}
			}
		}
		opt.Loc.End = p.pos()
		return opt
	}

	// Unknown option — consume name and optional = value structurally
	p.advance() // consume option name
	opt := &nodes.BackupRestoreOption{
		Name: name,
		Loc:  nodes.Loc{Start: optLoc},
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
	opt.Loc.End = p.pos()
	return opt
}

// parseBackupEncryptionOption parses the ENCRYPTION option.
//
//	ENCRYPTION ( ALGORITHM = { AES_128 | AES_192 | AES_256 | TRIPLE_DES_3KEY },
//	    SERVER CERTIFICATE = cert_name | SERVER ASYMMETRIC KEY = key_name )
func (p *Parser) parseBackupEncryptionOption() *nodes.BackupRestoreOption {
	optLoc := p.pos()
	p.advance() // consume ENCRYPTION

	opt := &nodes.BackupRestoreOption{
		Name: "ENCRYPTION",
		Loc:  nodes.Loc{Start: optLoc},
	}

	if p.cur.Type != '(' {
		opt.Loc.End = p.pos()
		return opt
	}
	p.advance() // consume (

	// ALGORITHM = alg
	if p.isIdentLike() && strings.EqualFold(p.cur.Str, "ALGORITHM") {
		p.advance() // consume ALGORITHM
		if _, ok := p.match('='); ok {
			if p.isIdentLike() {
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
	if p.isIdentLike() && strings.EqualFold(p.cur.Str, "SERVER") {
		p.advance() // consume SERVER
		if p.isIdentLike() {
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

	opt.Loc.End = p.pos()
	return opt
}

// parseRestoreMoveOption parses MOVE 'logical_file_name' TO 'os_file_name'.
func (p *Parser) parseRestoreMoveOption() *nodes.BackupRestoreOption {
	optLoc := p.pos()
	p.advance() // consume MOVE

	opt := &nodes.BackupRestoreOption{
		Name: "MOVE",
		Loc:  nodes.Loc{Start: optLoc},
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

	opt.Loc.End = p.pos()
	return opt
}
