// Package parser - server.go implements T-SQL server-level object statement parsing:
// CREATE/ALTER/DROP SERVER ROLE, ALTER SERVER CONFIGURATION,
// ALTER/BACKUP/RESTORE SERVICE MASTER KEY.
package parser

import (
	"strings"

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
//	   | MEMORY_OPTIMIZED { ON | OFF | TEMPDB_METADATA = { ON [(RESOURCE_POOL='pool')] | OFF } | HYBRID_BUFFER_POOL = { ON | OFF } }
//	   | HARDWARE_OFFLOAD { ON | OFF }
//	   | SUSPEND_FOR_SNAPSHOT_BACKUP = { ON | OFF } [ ( GROUP = ( <database>,...n ) [ , MODE = COPY_ONLY ] ) ]
//	}
//
//	<CPU_range_spec> ::=
//	    { CPU_ID | CPU_ID TO CPU_ID } [ ,...n ]
//
//	<NUMA_node_range_spec> ::=
//	    { NUMA_node_ID | NUMA_node_ID TO NUMA_node_ID } [ ,...n ]
//
//	<resource_property> ::=
//	{
//	    VerboseLogging = { 'logging_detail' | DEFAULT }
//	  | SqlDumperDumpFlags = { 'dump_file_type' | DEFAULT }
//	  | SqlDumperDumpPath = { 'os_file_path' | DEFAULT }
//	  | SqlDumperDumpTimeOut = { 'dump_time-out' | DEFAULT }
//	  | FailureConditionLevel = { 'failure_condition_level' | DEFAULT }
//	  | HealthCheckTimeout = { 'health_check_time-out' | DEFAULT }
//	  | ClusterConnectionOptions = '<key_value_pairs>[;...]'
//	}
//
//	<size_spec> ::=
//	    { size [ KB | MB | GB ] }
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
		p.advance() // consume PROCESS
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "AFFINITY") {
			p.advance() // consume AFFINITY
		}
		stmt.OptionType = "PROCESS AFFINITY"
		opts = p.parseServerConfigProcessAffinity()
	} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "DIAGNOSTICS") {
		p.advance() // consume DIAGNOSTICS
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "LOG") {
			p.advance() // consume LOG
		}
		stmt.OptionType = "DIAGNOSTICS LOG"
		opts = p.parseServerConfigDiagnosticsLog()
	} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "FAILOVER") {
		p.advance() // consume FAILOVER
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "CLUSTER") {
			p.advance() // consume CLUSTER
		}
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "PROPERTY") {
			p.advance() // consume PROPERTY
		}
		stmt.OptionType = "FAILOVER CLUSTER PROPERTY"
		opts = p.parseServerConfigFailoverClusterProperty()
	} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "HADR") {
		p.advance() // consume HADR
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "CLUSTER") {
			p.advance() // consume CLUSTER
		}
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "CONTEXT") {
			p.advance() // consume CONTEXT
		}
		stmt.OptionType = "HADR CLUSTER CONTEXT"
		opts = p.parseServerConfigHadrCluster()
	} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "BUFFER") {
		p.advance() // consume BUFFER
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "POOL") {
			p.advance() // consume POOL
		}
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "EXTENSION") {
			p.advance() // consume EXTENSION
		}
		stmt.OptionType = "BUFFER POOL EXTENSION"
		opts = p.parseServerConfigBufferPoolExtension()
	} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SOFTNUMA") {
		p.advance() // consume SOFTNUMA
		stmt.OptionType = "SOFTNUMA"
		opts = p.parseServerConfigOnOff()
	} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "MEMORY_OPTIMIZED") {
		p.advance() // consume MEMORY_OPTIMIZED
		stmt.OptionType = "MEMORY_OPTIMIZED"
		opts = p.parseServerConfigMemoryOptimized()
	} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "HARDWARE_OFFLOAD") {
		p.advance() // consume HARDWARE_OFFLOAD
		stmt.OptionType = "HARDWARE_OFFLOAD"
		opts = p.parseServerConfigOnOff()
	} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SUSPEND_FOR_SNAPSHOT_BACKUP") {
		p.advance() // consume SUSPEND_FOR_SNAPSHOT_BACKUP
		stmt.OptionType = "SUSPEND_FOR_SNAPSHOT_BACKUP"
		opts = p.parseServerConfigSuspendForSnapshotBackup()
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

// parseServerConfigProcessAffinity parses PROCESS AFFINITY options.
//
//	PROCESS AFFINITY { CPU = { AUTO | <CPU_range_spec> } | NUMANODE = <NUMA_node_range_spec> }
//	<CPU_range_spec> ::= { CPU_ID | CPU_ID TO CPU_ID } [ ,...n ]
//	<NUMA_node_range_spec> ::= { NUMA_node_ID | NUMA_node_ID TO NUMA_node_ID } [ ,...n ]
func (p *Parser) parseServerConfigProcessAffinity() []nodes.Node {
	var opts []nodes.Node

	// Expect CPU or NUMANODE
	if !p.isIdentLike() {
		return opts
	}

	key := p.cur.Str // CPU or NUMANODE
	p.advance()

	// Consume =
	if p.cur.Type == '=' {
		p.advance()
	}

	// Check for AUTO
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "AUTO") {
		opts = append(opts, &nodes.String{Str: key + "=AUTO"})
		p.advance()
		return opts
	}

	// Parse range spec: { ID | ID TO ID } [ ,...n ]
	var ranges []string
	for {
		if p.cur.Type != tokICONST {
			break
		}
		startID := p.cur.Str
		p.advance()

		if p.cur.Type == kwTO {
			p.advance() // consume TO
			if p.cur.Type == tokICONST {
				endID := p.cur.Str
				p.advance()
				ranges = append(ranges, startID+" TO "+endID)
			}
		} else {
			ranges = append(ranges, startID)
		}

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume comma
	}

	val := ""
	for i, r := range ranges {
		if i > 0 {
			val += ", "
		}
		val += r
	}
	opts = append(opts, &nodes.String{Str: key + "=" + val})
	return opts
}

// parseServerConfigDiagnosticsLog parses DIAGNOSTICS LOG options.
//
//	DIAGNOSTICS LOG { ON | OFF | PATH = { 'os_file_path' | DEFAULT } | MAX_SIZE = { 'log_max_size' MB | DEFAULT } | MAX_FILES = { 'max_file_count' | DEFAULT } }
func (p *Parser) parseServerConfigDiagnosticsLog() []nodes.Node {
	var opts []nodes.Node

	// ON / OFF
	if p.cur.Type == kwON {
		opts = append(opts, &nodes.String{Str: "ON"})
		p.advance()
		return opts
	}
	if p.cur.Type == kwOFF {
		opts = append(opts, &nodes.String{Str: "OFF"})
		p.advance()
		return opts
	}

	// PATH / MAX_SIZE / MAX_FILES
	if !p.isIdentLike() {
		return opts
	}

	key := strings.ToUpper(p.cur.Str) // PATH, MAX_SIZE, MAX_FILES
	p.advance()

	// Consume =
	if p.cur.Type == '=' {
		p.advance()
	}

	// DEFAULT
	if p.cur.Type == kwDEFAULT {
		opts = append(opts, &nodes.String{Str: key + "=DEFAULT"})
		p.advance()
		return opts
	}

	// String constant (for PATH)
	if p.cur.Type == tokSCONST {
		opts = append(opts, &nodes.String{Str: key + "='" + p.cur.Str + "'"})
		p.advance()
		return opts
	}

	// Integer (for MAX_SIZE or MAX_FILES)
	if p.cur.Type == tokICONST {
		val := p.cur.Str
		p.advance()
		// Check for MB suffix (MAX_SIZE)
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "MB") {
			val += " MB"
			p.advance()
		}
		opts = append(opts, &nodes.String{Str: key + "=" + val})
		return opts
	}

	return opts
}

// parseServerConfigFailoverClusterProperty parses FAILOVER CLUSTER PROPERTY <resource_property>.
//
//	<resource_property> ::=
//	{
//	    VerboseLogging = { 'logging_detail' | DEFAULT }
//	  | SqlDumperDumpFlags = { 'dump_file_type' | DEFAULT }
//	  | SqlDumperDumpPath = { 'os_file_path' | DEFAULT }
//	  | SqlDumperDumpTimeOut = { 'dump_time-out' | DEFAULT }
//	  | FailureConditionLevel = { 'failure_condition_level' | DEFAULT }
//	  | HealthCheckTimeout = { 'health_check_time-out' | DEFAULT }
//	  | ClusterConnectionOptions = '<key_value_pairs>[;...]'
//	}
func (p *Parser) parseServerConfigFailoverClusterProperty() []nodes.Node {
	var opts []nodes.Node

	if !p.isIdentLike() {
		return opts
	}

	key := p.cur.Str
	p.advance()

	// Consume =
	if p.cur.Type == '=' {
		p.advance()
	}

	// DEFAULT
	if p.cur.Type == kwDEFAULT {
		opts = append(opts, &nodes.String{Str: key + "=DEFAULT"})
		p.advance()
		return opts
	}

	// String constant
	if p.cur.Type == tokSCONST {
		opts = append(opts, &nodes.String{Str: key + "='" + p.cur.Str + "'"})
		p.advance()
		return opts
	}

	// Integer constant
	if p.cur.Type == tokICONST {
		opts = append(opts, &nodes.String{Str: key + "=" + p.cur.Str})
		p.advance()
		return opts
	}

	return opts
}

// parseServerConfigHadrCluster parses HADR CLUSTER CONTEXT = { 'remote_windows_cluster' | LOCAL }.
func (p *Parser) parseServerConfigHadrCluster() []nodes.Node {
	var opts []nodes.Node

	// Consume =
	if p.cur.Type == '=' {
		p.advance()
	}

	// String constant or LOCAL
	if p.cur.Type == tokSCONST {
		opts = append(opts, &nodes.String{Str: "CONTEXT='" + p.cur.Str + "'"})
		p.advance()
	} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "LOCAL") {
		opts = append(opts, &nodes.String{Str: "CONTEXT=LOCAL"})
		p.advance()
	}

	return opts
}

// parseServerConfigBufferPoolExtension parses BUFFER POOL EXTENSION options.
//
//	BUFFER POOL EXTENSION { ON ( FILENAME = 'path', SIZE = <size_spec> ) | OFF }
//	<size_spec> ::= { size [ KB | MB | GB ] }
func (p *Parser) parseServerConfigBufferPoolExtension() []nodes.Node {
	var opts []nodes.Node

	if p.cur.Type == kwOFF {
		opts = append(opts, &nodes.String{Str: "OFF"})
		p.advance()
		return opts
	}

	if p.cur.Type == kwON {
		opts = append(opts, &nodes.String{Str: "ON"})
		p.advance()

		// ( FILENAME = 'path', SIZE = <size_spec> )
		if p.cur.Type == '(' {
			p.advance() // consume (
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				if !p.isIdentLike() {
					p.advance()
					continue
				}
				key := p.cur.Str // FILENAME or SIZE
				p.advance()

				if p.cur.Type == '=' {
					p.advance() // consume =
				}

				if matchesKeywordCI(key, "FILENAME") {
					if p.cur.Type == tokSCONST {
						opts = append(opts, &nodes.String{Str: "FILENAME='" + p.cur.Str + "'"})
						p.advance()
					}
				} else if matchesKeywordCI(key, "SIZE") {
					if p.cur.Type == tokICONST || p.cur.Type == tokFCONST {
						val := p.cur.Str
						p.advance()
						// Check for KB/MB/GB suffix
						if p.isIdentLike() {
							upper := p.cur.Str
							if matchesKeywordCI(upper, "KB") || matchesKeywordCI(upper, "MB") || matchesKeywordCI(upper, "GB") {
								val += " " + upper
								p.advance()
							}
						}
						opts = append(opts, &nodes.String{Str: "SIZE=" + val})
					}
				}

				if p.cur.Type == ',' {
					p.advance() // consume comma
				}
			}
			if p.cur.Type == ')' {
				p.advance() // consume )
			}
		}
	}

	return opts
}

// parseServerConfigOnOff parses a simple ON/OFF option (SOFTNUMA, HARDWARE_OFFLOAD).
func (p *Parser) parseServerConfigOnOff() []nodes.Node {
	var opts []nodes.Node
	if p.cur.Type == kwON {
		opts = append(opts, &nodes.String{Str: "ON"})
		p.advance()
	} else if p.cur.Type == kwOFF {
		opts = append(opts, &nodes.String{Str: "OFF"})
		p.advance()
	}
	return opts
}

// parseServerConfigMemoryOptimized parses MEMORY_OPTIMIZED options.
//
//	MEMORY_OPTIMIZED { ON | OFF | TEMPDB_METADATA = { ON [(RESOURCE_POOL='pool')] | OFF } | HYBRID_BUFFER_POOL = { ON | OFF } }
func (p *Parser) parseServerConfigMemoryOptimized() []nodes.Node {
	var opts []nodes.Node

	// Direct ON/OFF
	if p.cur.Type == kwON {
		opts = append(opts, &nodes.String{Str: "ON"})
		p.advance()
		return opts
	}
	if p.cur.Type == kwOFF {
		opts = append(opts, &nodes.String{Str: "OFF"})
		p.advance()
		return opts
	}

	// TEMPDB_METADATA or HYBRID_BUFFER_POOL
	if !p.isIdentLike() {
		return opts
	}

	key := p.cur.Str
	p.advance()

	// Consume =
	if p.cur.Type == '=' {
		p.advance()
	}

	if p.cur.Type == kwON {
		opts = append(opts, &nodes.String{Str: key + "=ON"})
		p.advance()

		// Check for (RESOURCE_POOL = 'pool_name') -- only for TEMPDB_METADATA
		if matchesKeywordCI(key, "TEMPDB_METADATA") && p.cur.Type == '(' {
			p.advance() // consume (
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "RESOURCE_POOL") {
				p.advance() // consume RESOURCE_POOL
				if p.cur.Type == '=' {
					p.advance() // consume =
				}
				if p.cur.Type == tokSCONST {
					opts = append(opts, &nodes.String{Str: "RESOURCE_POOL='" + p.cur.Str + "'"})
					p.advance()
				}
			}
			if p.cur.Type == ')' {
				p.advance() // consume )
			}
		}
	} else if p.cur.Type == kwOFF {
		opts = append(opts, &nodes.String{Str: key + "=OFF"})
		p.advance()
	}

	return opts
}

// parseServerConfigSuspendForSnapshotBackup parses SUSPEND_FOR_SNAPSHOT_BACKUP options.
//
//	SUSPEND_FOR_SNAPSHOT_BACKUP = { ON | OFF } [ ( GROUP = ( <database>,...n ) [ , MODE = COPY_ONLY ] ) ]
func (p *Parser) parseServerConfigSuspendForSnapshotBackup() []nodes.Node {
	var opts []nodes.Node

	// Consume =
	if p.cur.Type == '=' {
		p.advance()
	}

	if p.cur.Type == kwON {
		opts = append(opts, &nodes.String{Str: "ON"})
		p.advance()
	} else if p.cur.Type == kwOFF {
		opts = append(opts, &nodes.String{Str: "OFF"})
		p.advance()
		return opts
	}

	// Optional ( GROUP = ( db1, db2 ) [ , MODE = COPY_ONLY ] )
	if p.cur.Type == '(' {
		p.advance() // consume outer (

		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "GROUP") {
				p.advance() // consume GROUP
				if p.cur.Type == '=' {
					p.advance() // consume =
				}
				// ( db1, db2, ... )
				if p.cur.Type == '(' {
					p.advance() // consume inner (
					var dbs []string
					for {
						if name, ok := p.parseIdentifier(); ok {
							dbs = append(dbs, name)
						}
						if p.cur.Type != ',' {
							break
						}
						p.advance() // consume comma
					}
					if p.cur.Type == ')' {
						p.advance() // consume inner )
					}
					val := ""
					for i, db := range dbs {
						if i > 0 {
							val += ", "
						}
						val += db
					}
					opts = append(opts, &nodes.String{Str: "GROUP=" + val})
				}
			} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "MODE") {
				p.advance() // consume MODE
				if p.cur.Type == '=' {
					p.advance() // consume =
				}
				if p.isIdentLike() {
					opts = append(opts, &nodes.String{Str: "MODE=" + p.cur.Str})
					p.advance()
				}
			} else {
				p.advance()
			}

			if p.cur.Type == ',' {
				p.advance() // consume comma between GROUP and MODE
			}
		}
		if p.cur.Type == ')' {
			p.advance() // consume outer )
		}
	}

	return opts
}

// consumeServerConfigOptions consumes remaining tokens for an ALTER SERVER CONFIGURATION
// option until a statement boundary (semicolon or EOF), collecting them as string nodes.
// Used as a fallback for unknown option types.
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
