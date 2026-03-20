// Package parser - security_keys.go implements parsing for security key/cert/credential statements.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseSecurityKeyStmt parses CREATE/ALTER/DROP statements for security objects:
// MASTER KEY, SYMMETRIC KEY, ASYMMETRIC KEY, CERTIFICATE, CREDENTIAL,
// COLUMN ENCRYPTION KEY, COLUMN MASTER KEY, CRYPTOGRAPHIC PROVIDER.
// Also handles OPEN/CLOSE SYMMETRIC KEY and BACKUP CERTIFICATE.
//
// BNF: mssql/parser/bnf/create-symmetric-key-transact-sql.bnf
//
//	CREATE SYMMETRIC KEY key_name
//	    [ AUTHORIZATION owner_name ]
//	    [ FROM PROVIDER provider_name ]
//	    WITH
//	        [
//	            <key_options> [ , ... n ]
//	            | ENCRYPTION BY <encrypting_mechanism> [ , ... n ]
//	        ]
//	<key_options> ::=
//	    KEY_SOURCE = 'pass_phrase'
//	    | ALGORITHM = <algorithm>
//	    | IDENTITY_VALUE = 'identity_phrase'
//	    | PROVIDER_KEY_NAME = 'key_name_in_provider'
//	    | CREATION_DISPOSITION = { CREATE_NEW | OPEN_EXISTING }
//	<algorithm> ::=
//	    DES | TRIPLE_DES | TRIPLE_DES_3KEY | RC2 | RC4 | RC4_128
//	    | DESX | AES_128 | AES_192 | AES_256
//	<encrypting_mechanism> ::=
//	    CERTIFICATE certificate_name
//	    | PASSWORD = 'password'
//	    | SYMMETRIC KEY symmetric_key_name
//	    | ASYMMETRIC KEY asym_key_name
//
// BNF: mssql/parser/bnf/alter-symmetric-key-transact-sql.bnf
//
//	ALTER SYMMETRIC KEY Key_name <alter_option>
//	<alter_option> ::=
//	    ADD ENCRYPTION BY <encrypting_mechanism> [ , ... n ]
//	    | DROP ENCRYPTION BY <encrypting_mechanism> [ , ... n ]
//	<encrypting_mechanism> ::=
//	    CERTIFICATE certificate_name
//	    | PASSWORD = 'password'
//	    | SYMMETRIC KEY Symmetric_Key_Name
//	    | ASYMMETRIC KEY Asym_Key_Name
//
// BNF: mssql/parser/bnf/drop-symmetric-key-transact-sql.bnf
//
//	DROP SYMMETRIC KEY symmetric_key_name [REMOVE PROVIDER KEY]
//
// BNF: mssql/parser/bnf/create-asymmetric-key-transact-sql.bnf
//
//	CREATE ASYMMETRIC KEY asym_key_name
//	    [ AUTHORIZATION database_principal_name ]
//	    [ FROM <asym_key_source> ]
//	    [ WITH <key_option> ]
//	    [ ENCRYPTION BY <encrypting_mechanism> ]
//	<asym_key_source> ::=
//	    FILE = 'path_to_strong-name_file'
//	    | EXECUTABLE FILE = 'path_to_executable_file'
//	    | ASSEMBLY assembly_name
//	    | PROVIDER provider_name
//	<key_option> ::=
//	    ALGORITHM = <algorithm>
//	    | PROVIDER_KEY_NAME = 'key_name_in_provider'
//	    | CREATION_DISPOSITION = { CREATE_NEW | OPEN_EXISTING }
//	<algorithm> ::= RSA_4096 | RSA_3072 | RSA_2048 | RSA_1024 | RSA_512
//	<encrypting_mechanism> ::= PASSWORD = 'password'
//
// BNF: mssql/parser/bnf/alter-asymmetric-key-transact-sql.bnf
//
//	ALTER ASYMMETRIC KEY Asym_Key_Name <alter_option>
//	<alter_option> ::=
//	    <password_change_option>
//	    | REMOVE PRIVATE KEY
//	<password_change_option> ::=
//	    WITH PRIVATE KEY ( <password_option> [ , <password_option> ] )
//	<password_option> ::=
//	    ENCRYPTION BY PASSWORD = 'strongPassword'
//	    | DECRYPTION BY PASSWORD = 'oldPassword'
//
// BNF: mssql/parser/bnf/drop-asymmetric-key-transact-sql.bnf
//
//	DROP ASYMMETRIC KEY key_name [ REMOVE PROVIDER KEY ]
//
// BNF: mssql/parser/bnf/create-certificate-transact-sql.bnf
//
//	CREATE CERTIFICATE certificate_name [ AUTHORIZATION user_name ]
//	    { FROM <existing_keys> | <generate_new_keys> }
//	    [ ACTIVE FOR BEGIN_DIALOG = { ON | OFF } ]
//	<existing_keys> ::=
//	    ASSEMBLY assembly_name
//	    | { [ EXECUTABLE ] FILE = 'path_to_file'
//	        [ WITH [FORMAT = 'PFX',] PRIVATE KEY ( <private_key_options> ) ] }
//	    | { BINARY = asn_encoded_certificate
//	        [ WITH PRIVATE KEY ( <private_key_options> ) ] }
//	<generate_new_keys> ::=
//	    [ ENCRYPTION BY PASSWORD = 'password' ]
//	    WITH SUBJECT = 'certificate_subject_name'
//	    [ , <date_options> [ ,...n ] ]
//	<private_key_options> ::=
//	    { FILE = 'path_to_private_key'
//	      [ , DECRYPTION BY PASSWORD = 'password' ]
//	      [ , ENCRYPTION BY PASSWORD = 'password' ] }
//	    | { BINARY = private_key_bits
//	      [ , DECRYPTION BY PASSWORD = 'password' ]
//	      [ , ENCRYPTION BY PASSWORD = 'password' ] }
//	<date_options> ::= START_DATE = 'datetime' | EXPIRY_DATE = 'datetime'
//
// BNF: mssql/parser/bnf/alter-certificate-transact-sql.bnf
//
//	ALTER CERTIFICATE certificate_name
//	    REMOVE PRIVATE KEY
//	    | WITH PRIVATE KEY ( <private_key_spec> )
//	    | WITH ACTIVE FOR BEGIN_DIALOG = { ON | OFF }
//	<private_key_spec> ::=
//	    { FILE = 'path_to_private_key' | BINARY = private_key_bits }
//	    [ , DECRYPTION BY PASSWORD = 'current_password' ]
//	    [ , ENCRYPTION BY PASSWORD = 'new_password' ]
//
// BNF: mssql/parser/bnf/drop-certificate-transact-sql.bnf
//
//	DROP CERTIFICATE certificate_name
//
// BNF: mssql/parser/bnf/create-credential-transact-sql.bnf
//
//	CREATE CREDENTIAL credential_name
//	    WITH IDENTITY = 'identity_name'
//	    [ , SECRET = 'secret' ]
//	    [ FOR CRYPTOGRAPHIC PROVIDER cryptographic_provider_name ]
//
// BNF: mssql/parser/bnf/alter-credential-transact-sql.bnf
//
//	ALTER CREDENTIAL credential_name
//	    WITH IDENTITY = 'identity_name'
//	    [ , SECRET = 'secret' ]
//
// BNF: mssql/parser/bnf/drop-credential-transact-sql.bnf
//
//	DROP CREDENTIAL credential_name
//
// BNF: mssql/parser/bnf/create-master-key-transact-sql.bnf
//
//	CREATE MASTER KEY [ ENCRYPTION BY PASSWORD = 'password' ]
//
// BNF: mssql/parser/bnf/alter-master-key-transact-sql.bnf
//
//	ALTER MASTER KEY <alter_option>
//	<alter_option> ::= <regenerate_option> | <encryption_option>
//	<regenerate_option> ::= [ FORCE ] REGENERATE WITH ENCRYPTION BY PASSWORD = 'password'
//	<encryption_option> ::=
//	    ADD ENCRYPTION BY { SERVICE MASTER KEY | PASSWORD = 'password' }
//	    | DROP ENCRYPTION BY { SERVICE MASTER KEY | PASSWORD = 'password' }
//
// BNF: mssql/parser/bnf/drop-master-key-transact-sql.bnf
//
//	DROP MASTER KEY
func (p *Parser) parseSecurityKeyStmt(action string) *nodes.SecurityKeyStmt {
	loc := p.pos()
	stmt := &nodes.SecurityKeyStmt{
		Action: action,
		Loc:    nodes.Loc{Start: loc},
	}

	// Determine object type
	switch {
	case p.matchIdentCI("MASTER"):
		// MASTER KEY
		p.match(kwKEY)
		stmt.ObjectType = "MASTER KEY"
	case p.matchIdentCI("SYMMETRIC"):
		// SYMMETRIC KEY key_name
		p.match(kwKEY)
		stmt.ObjectType = "SYMMETRIC KEY"
		name, _ := p.parseIdentifier()
		stmt.Name = name
	case p.matchIdentCI("ASYMMETRIC"):
		// ASYMMETRIC KEY key_name
		p.match(kwKEY)
		stmt.ObjectType = "ASYMMETRIC KEY"
		name, _ := p.parseIdentifier()
		stmt.Name = name
	case p.matchIdentCI("CERTIFICATE"):
		stmt.ObjectType = "CERTIFICATE"
		name, _ := p.parseIdentifier()
		stmt.Name = name
	case p.matchIdentCI("CREDENTIAL"):
		stmt.ObjectType = "CREDENTIAL"
		name, _ := p.parseIdentifier()
		stmt.Name = name
	case p.matchIdentCI("CRYPTOGRAPHIC"):
		// CRYPTOGRAPHIC PROVIDER provider_name
		p.matchIdentCI("PROVIDER")
		stmt.ObjectType = "CRYPTOGRAPHIC PROVIDER"
		name, _ := p.parseIdentifier()
		stmt.Name = name
	default:
		stmt.Loc.End = p.pos()
		return stmt
	}

	// Dispatch to type-specific option parsers
	switch {
	case action == "ALTER" && stmt.ObjectType == "SYMMETRIC KEY":
		p.parseAlterSymmetricKeyOptions(stmt)
	case action == "ALTER" && stmt.ObjectType == "MASTER KEY":
		p.parseAlterMasterKeyOptions(stmt)
	default:
		p.parseSecurityKeyOptions(stmt)
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseAlterSymmetricKeyOptions parses ALTER SYMMETRIC KEY options.
//
//	ALTER SYMMETRIC KEY Key_name
//	    ADD ENCRYPTION BY <encrypting_mechanism> [ , ... n ]
//	    | DROP ENCRYPTION BY <encrypting_mechanism> [ , ... n ]
//	<encrypting_mechanism> ::= CERTIFICATE name | PASSWORD = 'password'
//	    | SYMMETRIC KEY name | ASYMMETRIC KEY name
func (p *Parser) parseAlterSymmetricKeyOptions(stmt *nodes.SecurityKeyStmt) {
	var items []nodes.Node

	// ADD ENCRYPTION BY ... or DROP ENCRYPTION BY ...
	var op string
	if p.cur.Type == kwADD {
		op = "ADD"
		p.advance()
	} else if p.cur.Type == kwDROP {
		op = "DROP"
		p.advance()
	}

	if op != "" {
		items = append(items, &nodes.String{Str: op})
		// ENCRYPTION BY
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "ENCRYPTION") {
			p.advance()
		}
		if p.cur.Type == kwBY {
			p.advance()
		}
		// Parse encrypting mechanisms: CERTIFICATE name | PASSWORD = 'pwd' | SYMMETRIC KEY name | ASYMMETRIC KEY name
		for {
			mech := p.parseEncryptingMechanism()
			if mech == "" {
				break
			}
			items = append(items, &nodes.String{Str: "ENCRYPTION BY " + mech})
			if _, ok := p.match(','); !ok {
				break
			}
		}
	}

	if len(items) > 0 {
		stmt.Options = &nodes.List{Items: items}
	}
}

// parseAlterMasterKeyOptions parses ALTER MASTER KEY options.
//
//	ALTER MASTER KEY <alter_option>
//	<alter_option> ::= <regenerate_option> | <encryption_option>
//	<regenerate_option> ::= [ FORCE ] REGENERATE WITH ENCRYPTION BY PASSWORD = 'password'
//	<encryption_option> ::=
//	    ADD ENCRYPTION BY { SERVICE MASTER KEY | PASSWORD = 'password' }
//	    | DROP ENCRYPTION BY { SERVICE MASTER KEY | PASSWORD = 'password' }
func (p *Parser) parseAlterMasterKeyOptions(stmt *nodes.SecurityKeyStmt) {
	var items []nodes.Node

	// Check for [FORCE] REGENERATE
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "FORCE") {
		items = append(items, &nodes.String{Str: "FORCE"})
		p.advance()
	}

	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "REGENERATE") {
		items = append(items, &nodes.String{Str: "REGENERATE"})
		p.advance()
		// WITH ENCRYPTION BY PASSWORD = 'password'
		if p.cur.Type == kwWITH {
			p.advance()
		}
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "ENCRYPTION") {
			p.advance()
		}
		if p.cur.Type == kwBY {
			p.advance()
		}
		mech := p.parseEncryptingMechanism()
		if mech != "" {
			items = append(items, &nodes.String{Str: "ENCRYPTION BY " + mech})
		}
	} else if p.cur.Type == kwADD || p.cur.Type == kwDROP {
		// ADD/DROP ENCRYPTION BY { SERVICE MASTER KEY | PASSWORD = 'password' }
		op := "ADD"
		if p.cur.Type == kwDROP {
			op = "DROP"
		}
		p.advance()
		items = append(items, &nodes.String{Str: op})

		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "ENCRYPTION") {
			p.advance()
		}
		if p.cur.Type == kwBY {
			p.advance()
		}
		mech := p.parseEncryptingMechanism()
		if mech != "" {
			items = append(items, &nodes.String{Str: "ENCRYPTION BY " + mech})
		}
	}

	if len(items) > 0 {
		stmt.Options = &nodes.List{Items: items}
	}
}

// parseEncryptingMechanism parses a single encrypting mechanism:
//
//	CERTIFICATE certificate_name | PASSWORD = 'password'
//	| SYMMETRIC KEY symmetric_key_name | ASYMMETRIC KEY asym_key_name
//	| SERVICE MASTER KEY
func (p *Parser) parseEncryptingMechanism() string {
	if p.matchIdentCI("CERTIFICATE") {
		name, _ := p.parseIdentifier()
		return "CERTIFICATE " + name
	}
	if p.matchIdentCI("PASSWORD") {
		if p.cur.Type == '=' {
			p.advance()
		}
		pwd := p.consumeAnyIdent()
		return "PASSWORD = " + pwd
	}
	if p.matchIdentCI("SYMMETRIC") {
		p.match(kwKEY)
		name, _ := p.parseIdentifier()
		return "SYMMETRIC KEY " + name
	}
	if p.matchIdentCI("ASYMMETRIC") {
		p.match(kwKEY)
		name, _ := p.parseIdentifier()
		return "ASYMMETRIC KEY " + name
	}
	if p.matchIdentCI("SERVICE") {
		// SERVICE MASTER KEY
		p.matchIdentCI("MASTER")
		p.match(kwKEY)
		return "SERVICE MASTER KEY"
	}
	if p.matchIdentCI("SERVER") {
		// SERVER CERTIFICATE name or SERVER ASYMMETRIC KEY name
		if p.matchIdentCI("CERTIFICATE") {
			name, _ := p.parseIdentifier()
			return "SERVER CERTIFICATE " + name
		}
		if p.matchIdentCI("ASYMMETRIC") {
			p.match(kwKEY)
			name, _ := p.parseIdentifier()
			return "SERVER ASYMMETRIC KEY " + name
		}
	}
	return ""
}

// parseSecurityKeyStmtColumn parses CREATE/ALTER/DROP COLUMN { ENCRYPTION KEY | MASTER KEY }.
// Called after "COLUMN" has been matched by the dispatcher.
//
// BNF: mssql/parser/bnf/create-column-encryption-key-transact-sql.bnf
//
//	CREATE COLUMN ENCRYPTION KEY key_name
//	WITH VALUES
//	  (
//	    COLUMN_MASTER_KEY = column_master_key_name,
//	    ALGORITHM = 'algorithm_name',
//	    ENCRYPTED_VALUE = varbinary_literal
//	  )
//	[, (
//	    COLUMN_MASTER_KEY = column_master_key_name,
//	    ALGORITHM = 'algorithm_name',
//	    ENCRYPTED_VALUE = varbinary_literal
//	  ) ]
//
// ALTER COLUMN ENCRYPTION KEY key_name
//
//	{ ADD | DROP } VALUE
//	  (
//	    COLUMN_MASTER_KEY = column_master_key_name
//	    [, ALGORITHM = 'algorithm_name', ENCRYPTED_VALUE = varbinary_literal ]
//	  )
//
// BNF: mssql/parser/bnf/create-column-master-key-transact-sql.bnf
//
//	CREATE COLUMN MASTER KEY key_name
//	  WITH (
//	    KEY_STORE_PROVIDER_NAME = 'key_store_provider_name',
//	    KEY_PATH = 'key_path'
//	    [, ENCLAVE_COMPUTATIONS (SIGNATURE = signature)]
//	  )
//
// DROP COLUMN ENCRYPTION KEY key_name
// DROP COLUMN MASTER KEY key_name
func (p *Parser) parseSecurityKeyStmtColumn(action string) *nodes.SecurityKeyStmt {
	loc := p.pos()
	stmt := &nodes.SecurityKeyStmt{
		Action: action,
		Loc:    nodes.Loc{Start: loc},
	}

	// Determine COLUMN ENCRYPTION KEY or COLUMN MASTER KEY
	if p.matchIdentCI("ENCRYPTION") {
		p.match(kwKEY)
		stmt.ObjectType = "COLUMN ENCRYPTION KEY"
	} else if p.matchIdentCI("MASTER") {
		p.match(kwKEY)
		stmt.ObjectType = "COLUMN MASTER KEY"
	} else {
		stmt.Loc.End = p.pos()
		return stmt
	}

	// Parse key name
	name, _ := p.parseIdentifier()
	stmt.Name = name

	// For ALTER COLUMN ENCRYPTION KEY, handle ADD/DROP VALUE specially
	// because DROP is a statement-start keyword that would cause parseSecurityKeyOptions to stop.
	if action == "ALTER" && stmt.ObjectType == "COLUMN ENCRYPTION KEY" {
		// Parse ADD VALUE or DROP VALUE as structured nodes
		var prefixItems []nodes.Node
		if p.cur.Type == kwADD {
			prefixItems = append(prefixItems, &nodes.String{Str: "ADD"})
			p.advance()
		} else if p.cur.Type == kwDROP {
			prefixItems = append(prefixItems, &nodes.String{Str: "DROP"})
			p.advance()
		}
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "VALUE") {
			prefixItems = append(prefixItems, &nodes.String{Str: "VALUE"})
			p.advance()
		}
		// Now consume the rest (parenthesized options)
		p.parseSecurityKeyOptions(stmt)
		// Prepend the ADD/DROP VALUE nodes to the options
		if len(prefixItems) > 0 {
			if stmt.Options != nil {
				stmt.Options.Items = append(prefixItems, stmt.Options.Items...)
			} else {
				stmt.Options = &nodes.List{Items: prefixItems}
			}
		}
	} else {
		// Consume remaining tokens as options
		p.parseSecurityKeyOptions(stmt)
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseSecurityKeyStmtDatabaseEncryption parses CREATE/ALTER/DROP DATABASE ENCRYPTION KEY
// and CREATE/ALTER/DROP DATABASE SCOPED CREDENTIAL.
// Called after "DATABASE" has been consumed and the next token is "ENCRYPTION" or "SCOPED".
//
// BNF: mssql/parser/bnf/create-database-encryption-key-transact-sql.bnf
//
//	CREATE DATABASE ENCRYPTION KEY
//	    WITH ALGORITHM = { AES_128 | AES_192 | AES_256 | TRIPLE_DES_3KEY }
//	    ENCRYPTION BY SERVER
//	        { CERTIFICATE Encryptor_Name | ASYMMETRIC KEY Encryptor_Name }
//
// BNF: mssql/parser/bnf/alter-database-encryption-key-transact-sql.bnf
//
//	ALTER DATABASE ENCRYPTION KEY
//	    REGENERATE WITH ALGORITHM = { AES_128 | AES_192 | AES_256 | TRIPLE_DES_3KEY }
//	        [ ENCRYPTION BY SERVER { CERTIFICATE Encryptor_Name | ASYMMETRIC KEY Encryptor_Name } ]
//	    | ENCRYPTION BY SERVER
//	        { CERTIFICATE Encryptor_Name | ASYMMETRIC KEY Encryptor_Name }
//
// BNF: mssql/parser/bnf/drop-database-encryption-key-transact-sql.bnf
//
//	DROP DATABASE ENCRYPTION KEY
//
// BNF: mssql/parser/bnf/create-database-scoped-credential-transact-sql.bnf
//
//	CREATE DATABASE SCOPED CREDENTIAL credential_name
//	    WITH IDENTITY = 'identity_name'
//	    [ , SECRET = 'secret' ]
//
//	ALTER DATABASE SCOPED CREDENTIAL credential_name
//	    WITH IDENTITY = 'identity_name'
//	    [ , SECRET = 'secret' ]
//
//	DROP DATABASE SCOPED CREDENTIAL credential_name
func (p *Parser) parseSecurityKeyStmtDatabaseEncryption(action string) *nodes.SecurityKeyStmt {
	loc := p.pos()
	stmt := &nodes.SecurityKeyStmt{
		Action: action,
		Loc:    nodes.Loc{Start: loc},
	}

	if p.matchIdentCI("ENCRYPTION") {
		// DATABASE ENCRYPTION KEY
		p.match(kwKEY)
		stmt.ObjectType = "DATABASE ENCRYPTION KEY"
		// No name for database encryption key
	} else if p.matchIdentCI("SCOPED") {
		// DATABASE SCOPED CREDENTIAL credential_name
		p.matchIdentCI("CREDENTIAL")
		stmt.ObjectType = "DATABASE SCOPED CREDENTIAL"
		name, _ := p.parseIdentifier()
		stmt.Name = name
	} else {
		stmt.Loc.End = p.pos()
		return stmt
	}

	// Consume remaining tokens as options
	p.parseSecurityKeyOptions(stmt)

	stmt.Loc.End = p.pos()
	return stmt
}

// parseOpenSymmetricKeyStmt parses OPEN SYMMETRIC KEY statements.
//
// BNF: mssql/parser/bnf/open-symmetric-key-transact-sql.bnf
//
//	OPEN SYMMETRIC KEY Key_name DECRYPTION BY <decryption_mechanism>
//	<decryption_mechanism> ::=
//	    CERTIFICATE certificate_name [ WITH PASSWORD = 'password' ]
//	    | ASYMMETRIC KEY asym_key_name [ WITH PASSWORD = 'password' ]
//	    | SYMMETRIC KEY decrypting_Key_name
//	    | PASSWORD = 'decryption_password'
func (p *Parser) parseOpenSymmetricKeyStmt() *nodes.SecurityKeyStmt {
	loc := p.pos()
	p.advance() // consume OPEN
	stmt := &nodes.SecurityKeyStmt{
		Action: "OPEN",
		Loc:    nodes.Loc{Start: loc},
	}

	if p.matchIdentCI("SYMMETRIC") {
		p.match(kwKEY)
		stmt.ObjectType = "SYMMETRIC KEY"
	}

	name, _ := p.parseIdentifier()
	stmt.Name = name

	p.parseSecurityKeyOptions(stmt)

	stmt.Loc.End = p.pos()
	return stmt
}

// parseCloseSymmetricKeyStmt parses CLOSE SYMMETRIC KEY or CLOSE ALL SYMMETRIC KEYS.
//
// BNF: mssql/parser/bnf/close-symmetric-key-transact-sql.bnf
//
//	CLOSE { SYMMETRIC KEY key_name | ALL SYMMETRIC KEYS }
func (p *Parser) parseCloseSymmetricKeyStmt() *nodes.SecurityKeyStmt {
	loc := p.pos()
	p.advance() // consume CLOSE
	stmt := &nodes.SecurityKeyStmt{
		Action: "CLOSE",
		Loc:    nodes.Loc{Start: loc},
	}

	if p.cur.Type == kwALL {
		p.advance() // consume ALL
		p.matchIdentCI("SYMMETRIC")
		p.matchIdentCI("KEYS")
		stmt.ObjectType = "ALL SYMMETRIC KEYS"
	} else if p.matchIdentCI("SYMMETRIC") {
		p.match(kwKEY)
		stmt.ObjectType = "SYMMETRIC KEY"
		name, _ := p.parseIdentifier()
		stmt.Name = name
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseBackupCertificateStmt parses BACKUP CERTIFICATE|MASTER KEY|SYMMETRIC KEY statements.
//
// BNF: mssql/parser/bnf/backup-certificate-transact-sql.bnf
//
//	BACKUP CERTIFICATE certname TO FILE = 'path_to_file'
//	    [ WITH
//	      [FORMAT = 'PFX',]
//	      PRIVATE KEY
//	      (
//	        FILE = 'path_to_private_key_file' ,
//	        ENCRYPTION BY PASSWORD = 'encryption_password'
//	        [ , DECRYPTION BY PASSWORD = 'decryption_password' ]
//	      )
//	    ]
//
// BNF: mssql/parser/bnf/backup-master-key-transact-sql.bnf
//
//	BACKUP MASTER KEY TO FILE = 'path_to_file'
//	    ENCRYPTION BY PASSWORD = 'password'
//
//	BACKUP SYMMETRIC KEY key_name TO { FILE = 'path' | URL = 'url' }
//	    ENCRYPTION BY PASSWORD = 'password'
func (p *Parser) parseBackupCertificateStmt() *nodes.SecurityKeyStmt {
	loc := p.pos()
	p.advance() // consume BACKUP
	stmt := &nodes.SecurityKeyStmt{
		Action: "BACKUP",
		Loc:    nodes.Loc{Start: loc},
	}

	if p.matchIdentCI("CERTIFICATE") {
		stmt.ObjectType = "CERTIFICATE"
		name, _ := p.parseIdentifier()
		stmt.Name = name
	} else if p.matchIdentCI("SYMMETRIC") {
		p.match(kwKEY)
		stmt.ObjectType = "SYMMETRIC KEY"
		name, _ := p.parseIdentifier()
		stmt.Name = name
	} else if p.matchIdentCI("MASTER") {
		p.match(kwKEY)
		stmt.ObjectType = "MASTER KEY"
	}

	p.parseSecurityKeyOptions(stmt)

	stmt.Loc.End = p.pos()
	return stmt
}

// parseOpenMasterKeyStmt parses OPEN MASTER KEY DECRYPTION BY PASSWORD = 'password'.
func (p *Parser) parseOpenMasterKeyStmt() *nodes.SecurityKeyStmt {
	loc := p.pos()
	p.advance() // consume OPEN
	stmt := &nodes.SecurityKeyStmt{
		Action:     "OPEN",
		ObjectType: "MASTER KEY",
		Loc:        nodes.Loc{Start: loc},
	}

	// consume MASTER KEY
	if p.matchIdentCI("MASTER") {
		p.match(kwKEY)
	}

	p.parseSecurityKeyOptions(stmt)
	stmt.Loc.End = p.pos()
	return stmt
}

// parseCloseMasterKeyStmt parses CLOSE MASTER KEY.
func (p *Parser) parseCloseMasterKeyStmt() *nodes.SecurityKeyStmt {
	loc := p.pos()
	p.advance() // consume CLOSE
	stmt := &nodes.SecurityKeyStmt{
		Action:     "CLOSE",
		ObjectType: "MASTER KEY",
		Loc:        nodes.Loc{Start: loc},
	}

	// consume MASTER KEY
	if p.matchIdentCI("MASTER") {
		p.match(kwKEY)
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseRestoreMasterKeyStmt parses RESTORE MASTER KEY FROM FILE = 'path' ...
//
// BNF: mssql/parser/bnf/restore-master-key-transact-sql.bnf
//
//	RESTORE MASTER KEY FROM { FILE = 'path_to_file' | URL = 'Azure Blob storage URL' }
//	    DECRYPTION BY PASSWORD = 'password'
//	    ENCRYPTION BY PASSWORD = 'password'
//	    [ FORCE ]
func (p *Parser) parseRestoreMasterKeyStmt() *nodes.SecurityKeyStmt {
	loc := p.pos()
	p.advance() // consume RESTORE
	stmt := &nodes.SecurityKeyStmt{
		Action:     "RESTORE",
		ObjectType: "MASTER KEY",
		Loc:        nodes.Loc{Start: loc},
	}

	// consume MASTER KEY
	if p.matchIdentCI("MASTER") {
		p.match(kwKEY)
	}

	p.parseSecurityKeyOptions(stmt)
	stmt.Loc.End = p.pos()
	return stmt
}

// parseRestoreSymmetricKeyStmt parses RESTORE SYMMETRIC KEY key_name FROM { FILE | URL } = '...' ...
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/restore-symmetric-key-transact-sql
//
//	RESTORE SYMMETRIC KEY key_name FROM
//	  {
//	    FILE = 'path_to_file'
//	  | URL = 'Azure Blob storage URL'
//	  }
//	      DECRYPTION BY PASSWORD = 'password'
//	      ENCRYPTION BY PASSWORD = 'password'
func (p *Parser) parseRestoreSymmetricKeyStmt() *nodes.SecurityKeyStmt {
	loc := p.pos()
	p.advance() // consume RESTORE
	stmt := &nodes.SecurityKeyStmt{
		Action:     "RESTORE",
		ObjectType: "SYMMETRIC KEY",
		Loc:        nodes.Loc{Start: loc},
	}

	// consume SYMMETRIC KEY
	if p.matchIdentCI("SYMMETRIC") {
		p.match(kwKEY)
	}

	// key name
	name, _ := p.parseIdentifier()
	stmt.Name = name

	p.parseSecurityKeyOptions(stmt)
	stmt.Loc.End = p.pos()
	return stmt
}

// parseSecurityKeyOptions consumes remaining tokens until a statement boundary,
// using structured key-value parsing when possible.
//
// Security key statements use several patterns:
//   - WITH ( KEY_SOURCE = ..., ALGORITHM = ..., ... )
//   - WITH SUBJECT = ..., EXPIRY_DATE = ...
//   - ENCRYPTION BY { PASSWORD = ... | CERTIFICATE ... | ASYMMETRIC KEY ... }
//   - DECRYPTION BY { PASSWORD = ... | CERTIFICATE ... | ASYMMETRIC KEY ... }
//   - WITH IDENTITY = '...', SECRET = '...'
//   - FROM FILE = '...'
//   - TO FILE = '...'
//
// We parse these as alternating name/value pairs in a List of String nodes.
func (p *Parser) parseSecurityKeyOptions(stmt *nodes.SecurityKeyStmt) {
	var items []nodes.Node
	depth := 0

	for p.cur.Type != tokEOF && p.cur.Type != ';' {
		// Break on keywords that start new statements (unless nested in parens)
		if depth == 0 && isStmtStart(p.cur.Type) {
			break
		}

		// WITH keyword: may introduce parenthesized or inline options
		if p.cur.Type == kwWITH {
			p.advance() // consume WITH
			if p.cur.Type == '(' {
				// WITH ( NAME = VALUE [, ...] )
				nested, _ := p.parseKeyValueOptionList()
				if nested != nil {
					items = append(items, nested.Items...)
				}
				continue
			}
			// WITH inline options: SUBJECT = ..., IDENTITY = ..., etc.
			// Fall through to parse as key=value below
		}

		// ENCRYPTION BY / DECRYPTION BY / FROM / TO patterns
		if p.isIdentLike() && (strings.EqualFold(p.cur.Str, "ENCRYPTION") ||
			strings.EqualFold(p.cur.Str, "DECRYPTION")) {
			keyword := strings.ToUpper(p.cur.Str)
			p.advance()
			if p.cur.Type == kwBY {
				p.advance() // consume BY
			}
			// Use parseEncryptingMechanism for structured parsing
			for {
				mech := p.parseEncryptingMechanism()
				if mech == "" {
					// Fallback: if no known mechanism matched, just record the keyword
					items = append(items, &nodes.String{Str: keyword + " BY"})
					break
				}
				items = append(items, &nodes.String{Str: keyword + " BY " + mech})
				if _, ok := p.match(','); !ok {
					break
				}
				// After comma, check if next token is another mechanism or a different clause
				if p.isIdentLike() && (strings.EqualFold(p.cur.Str, "ENCRYPTION") ||
					strings.EqualFold(p.cur.Str, "DECRYPTION")) {
					break
				}
				if p.cur.Type == kwWITH || p.cur.Type == tokEOF || p.cur.Type == ';' {
					break
				}
				if depth == 0 && isStmtStart(p.cur.Type) {
					break
				}
			}
			continue
		}

		if p.cur.Type == kwFROM || p.cur.Type == kwTO {
			keyword := strings.ToUpper(p.cur.Str)
			if keyword == "" {
				if p.cur.Type == kwFROM {
					keyword = "FROM"
				} else {
					keyword = "TO"
				}
			}
			p.advance()
			items = append(items, &nodes.String{Str: keyword})
			// Collect FILE = '...' or ASSEMBLY etc.
			continue
		}

		// General key = value parsing
		if p.isIdentLike() || p.cur.Type == tokICONST || p.cur.Type == tokSCONST {
			key := p.consumeAnyIdent()
			items = append(items, &nodes.String{Str: key})

			if p.cur.Type == '=' {
				p.advance()
				val := p.consumeAnyIdent()
				items = append(items, &nodes.String{Str: val})
			}
			p.match(',')
			continue
		}

		if p.cur.Type == '(' {
			depth++
			p.advance()
			continue
		}
		if p.cur.Type == ')' {
			if depth > 0 {
				depth--
				p.advance()
				continue
			}
			break
		}

		// Unknown token, consume to avoid infinite loop
		p.advance()
	}

	if len(items) > 0 {
		stmt.Options = &nodes.List{Items: items}
	}
}

// isStmtStart returns true if the token type could start a new statement.
func isStmtStart(t int) bool {
	switch t {
	case kwSELECT, kwINSERT, kwUPDATE, kwDELETE, kwCREATE, kwALTER, kwDROP,
		kwGRANT, kwREVOKE, kwDENY, kwEXEC, kwEXECUTE, kwDECLARE, kwSET,
		kwIF, kwWHILE, kwBEGIN, kwRETURN, kwGOTO, kwWAITFOR,
		kwPRINT, kwRAISERROR, kwTHROW, kwUSE, kwGO,
		kwCOMMIT, kwROLLBACK, kwSAVE, kwTRUNCATE, kwMERGE,
		kwOPEN, kwFETCH, kwCLOSE, kwDEALLOCATE, kwBACKUP:
		return true
	}
	return false
}
