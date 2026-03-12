// Package parser - server.go implements T-SQL server-level object statement parsing:
// CREATE/ALTER/DROP SERVER ROLE, ALTER SERVER CONFIGURATION,
// ALTER/BACKUP/RESTORE SERVICE MASTER KEY.
package parser

import (
	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseCreateServerRoleStmt parses CREATE SERVER ROLE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-server-role-transact-sql
//
//	CREATE SERVER ROLE role_name [ AUTHORIZATION server_principal ]
func (p *Parser) parseCreateServerRoleStmt() *nodes.SecurityStmt {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "SERVER ROLE",
		Loc:        nodes.Loc{Start: loc},
	}

	// role_name
	if name, ok := p.parseIdentifier(); ok {
		stmt.Name = name
	}

	// [ AUTHORIZATION server_principal ]
	if p.cur.Type == kwAUTHORIZATION {
		p.advance() // consume AUTHORIZATION
		var opts []nodes.Node
		if owner, ok := p.parseIdentifier(); ok {
			opts = append(opts, &nodes.String{Str: "AUTHORIZATION=" + owner})
		}
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseAlterServerRoleStmt parses ALTER SERVER ROLE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-server-role-transact-sql
//
//	ALTER SERVER ROLE server_role_name
//	{
//	    [ ADD MEMBER server_principal ]
//	  | [ DROP MEMBER server_principal ]
//	  | [ WITH NAME = new_server_role_name ]
//	}
func (p *Parser) parseAlterServerRoleStmt() *nodes.SecurityStmt {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "ALTER",
		ObjectType: "SERVER ROLE",
		Loc:        nodes.Loc{Start: loc},
	}

	// server_role_name
	if name, ok := p.parseIdentifier(); ok {
		stmt.Name = name
	}

	var opts []nodes.Node

	// ADD MEMBER | DROP MEMBER | WITH NAME = new_name
	if p.cur.Type == kwADD {
		p.advance() // consume ADD
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "MEMBER") {
			p.advance() // consume MEMBER
		}
		if member, ok := p.parseIdentifier(); ok {
			opts = append(opts, &nodes.String{Str: "ADD MEMBER=" + member})
		}
	} else if p.cur.Type == kwDROP {
		p.advance() // consume DROP
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "MEMBER") {
			p.advance() // consume MEMBER
		}
		if member, ok := p.parseIdentifier(); ok {
			opts = append(opts, &nodes.String{Str: "DROP MEMBER=" + member})
		}
	} else if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		for {
			if !p.isIdentLike() {
				break
			}
			key := p.cur.Str
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
				if val, ok := p.parseIdentifier(); ok {
					opts = append(opts, &nodes.String{Str: key + "=" + val})
				}
			} else {
				opts = append(opts, &nodes.String{Str: key})
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
	}

	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseDropServerRoleStmt parses DROP SERVER ROLE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-server-role-transact-sql
//
//	DROP SERVER ROLE role_name
func (p *Parser) parseDropServerRoleStmt() *nodes.SecurityStmt {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "DROP",
		ObjectType: "SERVER ROLE",
		Loc:        nodes.Loc{Start: loc},
	}

	if name, ok := p.parseIdentifier(); ok {
		stmt.Name = name
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseAlterServerConfigurationStmt parses ALTER SERVER CONFIGURATION SET <optionspec>.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-server-configuration-transact-sql
//
//	ALTER SERVER CONFIGURATION
//	SET <optionspec>
//
//	<optionspec> ::=
//	{
//	     PROCESS AFFINITY { CPU = { AUTO | <CPU_range_spec> } | NUMANODE = <NUMA_node_range_spec> }
//	   | DIAGNOSTICS LOG { ON | OFF | PATH = { 'os_file_path' | DEFAULT } | MAX_SIZE = { 'log_max_size' MB | DEFAULT } | MAX_FILES = { 'max_file_count' | DEFAULT } }
//	   | FAILOVER CLUSTER PROPERTY <resource_property>
//	   | HADR CLUSTER CONTEXT = { 'remote_windows_cluster' | LOCAL }
//	   | BUFFER POOL EXTENSION { ON ( FILENAME = 'path', SIZE = <size_spec> ) | OFF }
//	   | SOFTNUMA { ON | OFF }
//	   | MEMORY_OPTIMIZED { ON | OFF | TEMPDB_METADATA = { ON | OFF } | HYBRID_BUFFER_POOL = { ON | OFF } }
//	   | HARDWARE_OFFLOAD { ON | OFF }
//	   | SUSPEND_FOR_SNAPSHOT_BACKUP = { ON | OFF } [ ( GROUP = ( <database>,...n ) [ , MODE = COPY_ONLY ] ) ]
//	}
func (p *Parser) parseAlterServerConfigurationStmt() *nodes.AlterServerConfigurationStmt {
	loc := p.pos()
	// ALTER, SERVER, CONFIGURATION already consumed by caller
	// Consume SET
	if p.cur.Type == kwSET {
		p.advance()
	}

	stmt := &nodes.AlterServerConfigurationStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Determine the option type by looking at the first keyword(s)
	var opts []nodes.Node

	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "PROCESS") {
		// PROCESS AFFINITY ...
		p.advance() // consume PROCESS
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "AFFINITY") {
			p.advance() // consume AFFINITY
		}
		stmt.OptionType = "PROCESS AFFINITY"
		// Consume remaining tokens until statement boundary
		opts = p.consumeServerConfigOptions()
	} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "DIAGNOSTICS") {
		// DIAGNOSTICS LOG ...
		p.advance() // consume DIAGNOSTICS
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "LOG") {
			p.advance() // consume LOG
		}
		stmt.OptionType = "DIAGNOSTICS LOG"
		opts = p.consumeServerConfigOptions()
	} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "FAILOVER") {
		// FAILOVER CLUSTER PROPERTY ...
		p.advance() // consume FAILOVER
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "CLUSTER") {
			p.advance() // consume CLUSTER
		}
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "PROPERTY") {
			p.advance() // consume PROPERTY
		}
		stmt.OptionType = "FAILOVER CLUSTER PROPERTY"
		opts = p.consumeServerConfigOptions()
	} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "HADR") {
		// HADR CLUSTER CONTEXT = ...
		p.advance() // consume HADR
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "CLUSTER") {
			p.advance() // consume CLUSTER
		}
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "CONTEXT") {
			p.advance() // consume CONTEXT
		}
		stmt.OptionType = "HADR CLUSTER CONTEXT"
		opts = p.consumeServerConfigOptions()
	} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "BUFFER") {
		// BUFFER POOL EXTENSION ...
		p.advance() // consume BUFFER
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "POOL") {
			p.advance() // consume POOL
		}
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "EXTENSION") {
			p.advance() // consume EXTENSION
		}
		stmt.OptionType = "BUFFER POOL EXTENSION"
		opts = p.consumeServerConfigOptions()
	} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SOFTNUMA") {
		// SOFTNUMA ON/OFF
		p.advance() // consume SOFTNUMA
		stmt.OptionType = "SOFTNUMA"
		opts = p.consumeServerConfigOptions()
	} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "MEMORY_OPTIMIZED") {
		// MEMORY_OPTIMIZED ...
		p.advance() // consume MEMORY_OPTIMIZED
		stmt.OptionType = "MEMORY_OPTIMIZED"
		opts = p.consumeServerConfigOptions()
	} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "HARDWARE_OFFLOAD") {
		// HARDWARE_OFFLOAD ON/OFF
		p.advance() // consume HARDWARE_OFFLOAD
		stmt.OptionType = "HARDWARE_OFFLOAD"
		opts = p.consumeServerConfigOptions()
	} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SUSPEND_FOR_SNAPSHOT_BACKUP") {
		// SUSPEND_FOR_SNAPSHOT_BACKUP = ON/OFF
		p.advance() // consume SUSPEND_FOR_SNAPSHOT_BACKUP
		stmt.OptionType = "SUSPEND_FOR_SNAPSHOT_BACKUP"
		opts = p.consumeServerConfigOptions()
	} else {
		// Unknown option type - consume remaining tokens generically
		stmt.OptionType = "UNKNOWN"
		opts = p.consumeServerConfigOptions()
	}

	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// consumeServerConfigOptions consumes remaining tokens for an ALTER SERVER CONFIGURATION
// option until a statement boundary (semicolon or EOF), collecting them as string nodes.
func (p *Parser) consumeServerConfigOptions() []nodes.Node {
	var opts []nodes.Node
	depth := 0

	for p.cur.Type != tokEOF && p.cur.Type != ';' {
		if p.cur.Type == '(' {
			depth++
			opts = append(opts, &nodes.String{Str: "("})
			p.advance()
			continue
		}
		if p.cur.Type == ')' {
			depth--
			opts = append(opts, &nodes.String{Str: ")"})
			p.advance()
			if depth <= 0 {
				break
			}
			continue
		}
		if p.cur.Type == ',' {
			opts = append(opts, &nodes.String{Str: ","})
			p.advance()
			continue
		}

		// Collect the token value
		if p.cur.Type == tokSCONST {
			opts = append(opts, &nodes.String{Str: "'" + p.cur.Str + "'"})
		} else if p.cur.Type == tokICONST || p.cur.Type == tokFCONST {
			opts = append(opts, &nodes.String{Str: p.cur.Str})
		} else if p.cur.Type == '=' {
			opts = append(opts, &nodes.String{Str: "="})
		} else if p.cur.Type == ',' {
			opts = append(opts, &nodes.String{Str: ","})
		} else if p.isIdentLike() || p.cur.Type == kwON || p.cur.Type == kwOFF || p.cur.Type == kwNULL || p.cur.Type == kwDEFAULT || p.cur.Type == kwTO || p.cur.Type == kwKEY || p.cur.Type == kwFILE {
			opts = append(opts, &nodes.String{Str: p.cur.Str})
		} else {
			// Unknown token type - collect the string representation
			opts = append(opts, &nodes.String{Str: p.cur.Str})
		}
		p.advance()
	}
	return opts
}

// parseAlterServiceMasterKeyStmt parses ALTER SERVICE MASTER KEY.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-service-master-key-transact-sql
//
//	ALTER SERVICE MASTER KEY
//	    [ { <regenerate_option> | <recover_option> } ]
//
//	<regenerate_option> ::=
//	    [ FORCE ] REGENERATE
//
//	<recover_option> ::=
//	    { WITH OLD_ACCOUNT = 'account_name' , OLD_PASSWORD = 'password' }
//	    |
//	    { WITH NEW_ACCOUNT = 'account_name' , NEW_PASSWORD = 'password' }
func (p *Parser) parseAlterServiceMasterKeyStmt() *nodes.SecurityKeyStmt {
	loc := p.pos()
	// SERVICE MASTER KEY already consumed by caller
	stmt := &nodes.SecurityKeyStmt{
		Action:     "ALTER",
		ObjectType: "SERVICE MASTER KEY",
		Loc:        nodes.Loc{Start: loc},
	}

	p.parseSecurityKeyOptions(stmt)

	stmt.Loc.End = p.pos()
	return stmt
}

// parseBackupServiceMasterKeyStmt parses BACKUP SERVICE MASTER KEY.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/backup-service-master-key-transact-sql
//
//	BACKUP SERVICE MASTER KEY TO FILE = 'path_to_file'
//	    ENCRYPTION BY PASSWORD = 'password'
func (p *Parser) parseBackupServiceMasterKeyStmt() *nodes.SecurityKeyStmt {
	loc := p.pos()
	// BACKUP already consumed, SERVICE MASTER KEY consumed by caller
	stmt := &nodes.SecurityKeyStmt{
		Action:     "BACKUP",
		ObjectType: "SERVICE MASTER KEY",
		Loc:        nodes.Loc{Start: loc},
	}

	p.parseSecurityKeyOptions(stmt)

	stmt.Loc.End = p.pos()
	return stmt
}

// parseRestoreServiceMasterKeyStmt parses RESTORE SERVICE MASTER KEY.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/restore-service-master-key-transact-sql
//
//	RESTORE SERVICE MASTER KEY FROM FILE = 'path_to_file'
//	    DECRYPTION BY PASSWORD = 'password' [FORCE]
func (p *Parser) parseRestoreServiceMasterKeyStmt() *nodes.SecurityKeyStmt {
	loc := p.pos()
	// RESTORE already consumed, SERVICE MASTER KEY consumed by caller
	stmt := &nodes.SecurityKeyStmt{
		Action:     "RESTORE",
		ObjectType: "SERVICE MASTER KEY",
		Loc:        nodes.Loc{Start: loc},
	}

	p.parseSecurityKeyOptions(stmt)

	stmt.Loc.End = p.pos()
	return stmt
}
