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
func (p *Parser) parseCreateServerRoleStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "SERVER ROLE",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// role_name
	if name, ok := p.parseIdentifier(); ok {
		stmt.Name = name
	}

	// [ AUTHORIZATION server_principal ]
	if p.cur.Type == kwAUTHORIZATION {
		optLoc := p.pos()
		p.advance() // consume AUTHORIZATION
		var opts []nodes.Node
		if owner, ok := p.parseIdentifier(); ok {
			opts = append(opts, &nodes.ServerConfigOption{Name: "AUTHORIZATION", Value: owner, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
		}
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
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
func (p *Parser) parseAlterServerRoleStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "ALTER",
		ObjectType: "SERVER ROLE",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// server_role_name
	if name, ok := p.parseIdentifier(); ok {
		stmt.Name = name
	}

	var opts []nodes.Node

	// ADD MEMBER | DROP MEMBER | WITH NAME = new_name
	if p.cur.Type == kwADD {
		optLoc := p.pos()
		p.advance() // consume ADD
		if p.cur.Type == kwMEMBER {
			p.advance() // consume MEMBER
		}
		if member, ok := p.parseIdentifier(); ok {
			opts = append(opts, &nodes.ServerConfigOption{Name: "ADD MEMBER", Value: member, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
		}
	} else if p.cur.Type == kwDROP {
		optLoc := p.pos()
		p.advance() // consume DROP
		if p.cur.Type == kwMEMBER {
			p.advance() // consume MEMBER
		}
		if member, ok := p.parseIdentifier(); ok {
			opts = append(opts, &nodes.ServerConfigOption{Name: "DROP MEMBER", Value: member, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
		}
	} else if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		for {
			if !p.isAnyKeywordIdent() {
				break
			}
			optLoc := p.pos()
			key := strings.ToUpper(p.cur.Str)
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
				if val, ok := p.parseIdentifier(); ok {
					opts = append(opts, &nodes.ServerConfigOption{Name: key, Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
				}
			} else {
				opts = append(opts, &nodes.ServerConfigOption{Name: key, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
	}

	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseDropServerRoleStmt parses DROP SERVER ROLE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-server-role-transact-sql
//
//	DROP SERVER ROLE role_name
func (p *Parser) parseDropServerRoleStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "DROP",
		ObjectType: "SERVER ROLE",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	if name, ok := p.parseIdentifier(); ok {
		stmt.Name = name
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
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
//	   | EXTERNAL AUTHENTICATION { ON | OFF } [ ( USE_IDENTITY | CREDENTIAL_NAME = 'name' ) ]
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
func (p *Parser) parseAlterServerConfigurationStmt() (*nodes.AlterServerConfigurationStmt, error) {
	loc := p.pos()
	// ALTER, SERVER, CONFIGURATION already consumed by caller
	// Consume SET
	if p.cur.Type == kwSET {
		p.advance()
	}

	stmt := &nodes.AlterServerConfigurationStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Determine the option type by looking at the first keyword(s)
	var opts []nodes.Node

	if p.cur.Type == kwPROCESS {
		p.advance() // consume PROCESS
		if p.cur.Type == kwAFFINITY {
			p.advance() // consume AFFINITY
		}
		stmt.OptionType = "PROCESS AFFINITY"
		opts, _ = p.parseServerConfigProcessAffinity()
	} else if p.cur.Type == kwDIAGNOSTICS {
		p.advance() // consume DIAGNOSTICS
		if p.cur.Type == kwLOG {
			p.advance() // consume LOG
		}
		stmt.OptionType = "DIAGNOSTICS LOG"
		opts, _ = p.parseServerConfigDiagnosticsLog()
	} else if p.cur.Type == kwFAILOVER {
		p.advance() // consume FAILOVER
		if p.cur.Type == kwCLUSTER {
			p.advance() // consume CLUSTER
		}
		if p.cur.Type == kwPROPERTY {
			p.advance() // consume PROPERTY
		}
		stmt.OptionType = "FAILOVER CLUSTER PROPERTY"
		opts, _ = p.parseServerConfigFailoverClusterProperty()
	} else if p.cur.Type == kwHADR {
		p.advance() // consume HADR
		if p.cur.Type == kwCLUSTER {
			p.advance() // consume CLUSTER
		}
		if p.cur.Type == kwCONTEXT {
			p.advance() // consume CONTEXT
		}
		stmt.OptionType = "HADR CLUSTER CONTEXT"
		opts, _ = p.parseServerConfigHadrCluster()
	} else if p.cur.Type == kwBUFFER {
		p.advance() // consume BUFFER
		if p.cur.Type == kwPOOL {
			p.advance() // consume POOL
		}
		if p.cur.Type == kwEXTENSION {
			p.advance() // consume EXTENSION
		}
		stmt.OptionType = "BUFFER POOL EXTENSION"
		opts, _ = p.parseServerConfigBufferPoolExtension()
	} else if p.cur.Type == kwSOFTNUMA {
		p.advance() // consume SOFTNUMA
		stmt.OptionType = "SOFTNUMA"
		opts, _ = p.parseServerConfigOnOff()
	} else if p.cur.Type == kwMEMORY_OPTIMIZED {
		p.advance() // consume MEMORY_OPTIMIZED
		stmt.OptionType = "MEMORY_OPTIMIZED"
		opts, _ = p.parseServerConfigMemoryOptimized()
	} else if p.cur.Type == kwHARDWARE_OFFLOAD {
		p.advance() // consume HARDWARE_OFFLOAD
		stmt.OptionType = "HARDWARE_OFFLOAD"
		opts, _ = p.parseServerConfigOnOff()
	} else if p.cur.Type == kwEXTERNAL {
		p.advance() // consume EXTERNAL
		if p.cur.Type == kwAUTHENTICATION {
			p.advance() // consume AUTHENTICATION
		}
		stmt.OptionType = "EXTERNAL AUTHENTICATION"
		opts, _ = p.parseServerConfigExternalAuthentication()
	} else if p.cur.Type == kwSUSPEND_FOR_SNAPSHOT_BACKUP {
		p.advance() // consume SUSPEND_FOR_SNAPSHOT_BACKUP
		stmt.OptionType = "SUSPEND_FOR_SNAPSHOT_BACKUP"
		opts, _ = p.parseServerConfigSuspendForSnapshotBackup()
	} else {
		// Unknown option type - record the keyword and skip to statement boundary
		if p.isAnyKeywordIdent() {
			stmt.OptionType = strings.ToUpper(p.cur.Str)
			p.advance()
		} else {
			stmt.OptionType = "UNKNOWN"
		}
		// Skip remaining tokens until statement boundary
		for p.cur.Type != tokEOF && p.cur.Type != ';' {
			p.advance()
		}
	}

	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseServerConfigProcessAffinity parses PROCESS AFFINITY options.
//
//	PROCESS AFFINITY { CPU = { AUTO | <CPU_range_spec> } | NUMANODE = <NUMA_node_range_spec> }
//	<CPU_range_spec> ::= { CPU_ID | CPU_ID TO CPU_ID } [ ,...n ]
//	<NUMA_node_range_spec> ::= { NUMA_node_ID | NUMA_node_ID TO NUMA_node_ID } [ ,...n ]
func (p *Parser) parseServerConfigProcessAffinity() ([]nodes.Node, error) {
	var opts []nodes.Node

	// Expect CPU or NUMANODE
	if !p.isAnyKeywordIdent() {
		return opts, nil
	}

	optLoc := p.pos()
	key := strings.ToUpper(p.cur.Str) // CPU or NUMANODE
	p.advance()

	// Consume =
	if p.cur.Type == '=' {
		p.advance()
	}

	// Check for AUTO
	if p.cur.Type == kwAUTO {
		p.advance()
		opts = append(opts, &nodes.ServerConfigOption{Name: key, Value: "AUTO", Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
		return opts, nil
	}

	// Parse range spec: { ID | ID TO ID } [ ,...n ]
	var sb strings.Builder
	first := true
	for {
		if p.cur.Type != tokICONST {
			break
		}
		if !first {
			sb.WriteString(", ")
		}
		first = false
		startID := p.cur.Str
		p.advance()

		if p.cur.Type == kwTO {
			p.advance() // consume TO
			if p.cur.Type == tokICONST {
				sb.WriteString(startID + " TO " + p.cur.Str)
				p.advance()
			}
		} else {
			sb.WriteString(startID)
		}

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume comma
	}

	opts = append(opts, &nodes.ServerConfigOption{Name: key, Value: sb.String(), Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
	return opts, nil
}

// parseServerConfigDiagnosticsLog parses DIAGNOSTICS LOG options.
//
//	DIAGNOSTICS LOG { ON | OFF | PATH = { 'os_file_path' | DEFAULT } | MAX_SIZE = { 'log_max_size' MB | DEFAULT } | MAX_FILES = { 'max_file_count' | DEFAULT } }
func (p *Parser) parseServerConfigDiagnosticsLog() ([]nodes.Node, error) {
	var opts []nodes.Node

	// ON / OFF
	if p.cur.Type == kwON {
		optLoc := p.pos()
		p.advance()
		opts = append(opts, &nodes.ServerConfigOption{Name: "ON", Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
		return opts, nil
	}
	if p.cur.Type == kwOFF {
		optLoc := p.pos()
		p.advance()
		opts = append(opts, &nodes.ServerConfigOption{Name: "OFF", Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
		return opts, nil
	}

	// PATH / MAX_SIZE / MAX_FILES
	if !p.isAnyKeywordIdent() {
		return opts, nil
	}

	optLoc := p.pos()
	key := strings.ToUpper(p.cur.Str) // PATH, MAX_SIZE, MAX_FILES
	p.advance()

	// Consume =
	if p.cur.Type == '=' {
		p.advance()
	}

	// DEFAULT
	if p.cur.Type == kwDEFAULT {
		p.advance()
		opts = append(opts, &nodes.ServerConfigOption{Name: key, Value: "DEFAULT", Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
		return opts, nil
	}

	// String constant (for PATH)
	if p.cur.Type == tokSCONST {
		val := "'" + p.cur.Str + "'"
		p.advance()
		opts = append(opts, &nodes.ServerConfigOption{Name: key, Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
		return opts, nil
	}

	// Integer (for MAX_SIZE or MAX_FILES)
	if p.cur.Type == tokICONST {
		val := p.cur.Str
		p.advance()
		// Check for MB suffix (MAX_SIZE)
		if p.cur.Type == kwMB {
			val += " MB"
			p.advance()
		}
		opts = append(opts, &nodes.ServerConfigOption{Name: key, Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
		return opts, nil
	}

	return opts, nil
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
func (p *Parser) parseServerConfigFailoverClusterProperty() ([]nodes.Node, error) {
	var opts []nodes.Node

	if !p.isAnyKeywordIdent() {
		return opts, nil
	}

	optLoc := p.pos()
	key := p.cur.Str
	p.advance()

	// Consume =
	if p.cur.Type == '=' {
		p.advance()
	}

	// DEFAULT
	if p.cur.Type == kwDEFAULT {
		p.advance()
		opts = append(opts, &nodes.ServerConfigOption{Name: key, Value: "DEFAULT", Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
		return opts, nil
	}

	// String constant
	if p.cur.Type == tokSCONST {
		val := "'" + p.cur.Str + "'"
		p.advance()
		opts = append(opts, &nodes.ServerConfigOption{Name: key, Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
		return opts, nil
	}

	// Integer constant
	if p.cur.Type == tokICONST {
		val := p.cur.Str
		p.advance()
		opts = append(opts, &nodes.ServerConfigOption{Name: key, Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
		return opts, nil
	}

	return opts, nil
}

// parseServerConfigHadrCluster parses HADR CLUSTER CONTEXT = { 'remote_windows_cluster' | LOCAL }.
func (p *Parser) parseServerConfigHadrCluster() ([]nodes.Node, error) {
	var opts []nodes.Node

	optLoc := p.pos()
	// Consume =
	if p.cur.Type == '=' {
		p.advance()
	}

	// String constant or LOCAL
	if p.cur.Type == tokSCONST {
		val := "'" + p.cur.Str + "'"
		p.advance()
		opts = append(opts, &nodes.ServerConfigOption{Name: "CONTEXT", Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
	} else if p.cur.Type == kwLOCAL {
		p.advance()
		opts = append(opts, &nodes.ServerConfigOption{Name: "CONTEXT", Value: "LOCAL", Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
	}

	return opts, nil
}

// parseServerConfigBufferPoolExtension parses BUFFER POOL EXTENSION options.
//
//	BUFFER POOL EXTENSION { ON ( FILENAME = 'path', SIZE = <size_spec> ) | OFF }
//	<size_spec> ::= { size [ KB | MB | GB ] }
func (p *Parser) parseServerConfigBufferPoolExtension() ([]nodes.Node, error) {
	var opts []nodes.Node

	if p.cur.Type == kwOFF {
		optLoc := p.pos()
		p.advance()
		opts = append(opts, &nodes.ServerConfigOption{Name: "OFF", Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
		return opts, nil
	}

	if p.cur.Type == kwON {
		optLoc := p.pos()
		p.advance()
		opts = append(opts, &nodes.ServerConfigOption{Name: "ON", Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})

		// ( FILENAME = 'path', SIZE = <size_spec> )
		if p.cur.Type == '(' {
			p.advance() // consume (
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				if p.cur.Type == ',' {
					p.advance() // consume comma separator
					continue
				}
				if !p.isAnyKeywordIdent() {
					break // unexpected token - stop parsing options
				}
				subLoc := p.pos()
				keyType := p.cur.Type
				p.advance()

				if p.cur.Type == '=' {
					p.advance() // consume =
				}

				if keyType == kwFILENAME {
					if p.cur.Type == tokSCONST {
						val := "'" + p.cur.Str + "'"
						p.advance()
						opts = append(opts, &nodes.ServerConfigOption{Name: "FILENAME", Value: val, Loc: nodes.Loc{Start: subLoc, End: p.prevEnd()}})
					}
				} else if keyType == kwSIZE {
					if p.cur.Type == tokICONST || p.cur.Type == tokFCONST {
						val := p.cur.Str
						p.advance()
						// Check for KB/MB/GB suffix
						if p.isAnyKeywordIdent() {
							upper := strings.ToUpper(p.cur.Str)
							if upper == "KB" || upper == "MB" || upper == "GB" {
								val += " " + upper
								p.advance()
							}
						}
						opts = append(opts, &nodes.ServerConfigOption{Name: "SIZE", Value: val, Loc: nodes.Loc{Start: subLoc, End: p.prevEnd()}})
					}
				} else {
					break // unknown option key - stop parsing
				}
			}
			if p.cur.Type == ')' {
				p.advance() // consume )
			}
		}
	}

	return opts, nil
}

// parseServerConfigOnOff parses a simple ON/OFF option (SOFTNUMA, HARDWARE_OFFLOAD).
func (p *Parser) parseServerConfigOnOff() ([]nodes.Node, error) {
	var opts []nodes.Node
	if p.cur.Type == kwON {
		optLoc := p.pos()
		p.advance()
		opts = append(opts, &nodes.ServerConfigOption{Name: "ON", Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
	} else if p.cur.Type == kwOFF {
		optLoc := p.pos()
		p.advance()
		opts = append(opts, &nodes.ServerConfigOption{Name: "OFF", Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
	}
	return opts, nil
}

// parseServerConfigMemoryOptimized parses MEMORY_OPTIMIZED options.
//
//	MEMORY_OPTIMIZED { ON | OFF | TEMPDB_METADATA = { ON [(RESOURCE_POOL='pool')] | OFF } | HYBRID_BUFFER_POOL = { ON | OFF } }
func (p *Parser) parseServerConfigMemoryOptimized() ([]nodes.Node, error) {
	var opts []nodes.Node

	// Direct ON/OFF
	if p.cur.Type == kwON {
		optLoc := p.pos()
		p.advance()
		opts = append(opts, &nodes.ServerConfigOption{Name: "ON", Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
		return opts, nil
	}
	if p.cur.Type == kwOFF {
		optLoc := p.pos()
		p.advance()
		opts = append(opts, &nodes.ServerConfigOption{Name: "OFF", Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
		return opts, nil
	}

	// TEMPDB_METADATA or HYBRID_BUFFER_POOL
	if !p.isAnyKeywordIdent() {
		return opts, nil
	}

	optLoc := p.pos()
	key := strings.ToUpper(p.cur.Str)
	keyType := p.cur.Type
	p.advance()

	// Consume =
	if p.cur.Type == '=' {
		p.advance()
	}

	if p.cur.Type == kwON {
		p.advance()
		opts = append(opts, &nodes.ServerConfigOption{Name: key, Value: "ON", Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})

		// Check for (RESOURCE_POOL = 'pool_name') -- only for TEMPDB_METADATA
		if keyType == kwTEMPDB_METADATA && p.cur.Type == '(' {
			p.advance() // consume (
			if p.cur.Type == kwRESOURCE_POOL {
				subLoc := p.pos()
				p.advance() // consume RESOURCE_POOL
				if p.cur.Type == '=' {
					p.advance() // consume =
				}
				if p.cur.Type == tokSCONST {
					val := "'" + p.cur.Str + "'"
					p.advance()
					opts = append(opts, &nodes.ServerConfigOption{Name: "RESOURCE_POOL", Value: val, Loc: nodes.Loc{Start: subLoc, End: p.prevEnd()}})
				}
			}
			if p.cur.Type == ')' {
				p.advance() // consume )
			}
		}
	} else if p.cur.Type == kwOFF {
		p.advance()
		opts = append(opts, &nodes.ServerConfigOption{Name: key, Value: "OFF", Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
	}

	return opts, nil
}

// parseServerConfigSuspendForSnapshotBackup parses SUSPEND_FOR_SNAPSHOT_BACKUP options.
//
//	SUSPEND_FOR_SNAPSHOT_BACKUP = { ON | OFF } [ ( GROUP = ( <database>,...n ) [ , MODE = COPY_ONLY ] ) ]
func (p *Parser) parseServerConfigSuspendForSnapshotBackup() ([]nodes.Node, error) {
	var opts []nodes.Node

	// Consume =
	if p.cur.Type == '=' {
		p.advance()
	}

	if p.cur.Type == kwON {
		optLoc := p.pos()
		p.advance()
		opts = append(opts, &nodes.ServerConfigOption{Name: "ON", Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
	} else if p.cur.Type == kwOFF {
		optLoc := p.pos()
		p.advance()
		opts = append(opts, &nodes.ServerConfigOption{Name: "OFF", Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
		return opts, nil
	}

	// Optional ( GROUP = ( db1, db2 ) [ , MODE = COPY_ONLY ] )
	if p.cur.Type == '(' {
		p.advance() // consume outer (

		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			if p.cur.Type == ',' {
				p.advance() // consume comma between GROUP and MODE
				continue
			}
			if !p.isAnyKeywordIdent() {
				break // unexpected token
			}
			if p.cur.Type == kwGROUP {
				subLoc := p.pos()
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
					opts = append(opts, &nodes.ServerConfigOption{Name: "GROUP", Value: strings.Join(dbs, ", "), Loc: nodes.Loc{Start: subLoc, End: p.prevEnd()}})
				}
			} else if p.cur.Type == kwMODE {
				subLoc := p.pos()
				p.advance() // consume MODE
				if p.cur.Type == '=' {
					p.advance() // consume =
				}
				if p.isIdentLike() {
					val := p.cur.Str
					p.advance()
					opts = append(opts, &nodes.ServerConfigOption{Name: "MODE", Value: val, Loc: nodes.Loc{Start: subLoc, End: p.prevEnd()}})
				}
			} else {
				break // unknown option - stop parsing
			}
		}
		if p.cur.Type == ')' {
			p.advance() // consume outer )
		}
	}

	return opts, nil
}

// parseServerConfigExternalAuthentication parses EXTERNAL AUTHENTICATION options.
//
//	EXTERNAL AUTHENTICATION { ON | OFF } [ ( USE_IDENTITY | CREDENTIAL_NAME = 'name' ) ]
func (p *Parser) parseServerConfigExternalAuthentication() ([]nodes.Node, error) {
	var opts []nodes.Node

	// ON / OFF
	if p.cur.Type == kwON {
		optLoc := p.pos()
		p.advance()
		opts = append(opts, &nodes.ServerConfigOption{Name: "ON", Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
	} else if p.cur.Type == kwOFF {
		optLoc := p.pos()
		p.advance()
		opts = append(opts, &nodes.ServerConfigOption{Name: "OFF", Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
		return opts, nil
	}

	// Optional ( USE_IDENTITY | CREDENTIAL_NAME = 'name' )
	if p.cur.Type == '(' {
		p.advance() // consume (
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			if p.cur.Type == ',' {
				p.advance()
				continue
			}
			if !p.isAnyKeywordIdent() {
				break
			}
			subLoc := p.pos()
			key := p.cur.Str
			p.advance()
			if p.cur.Type == '=' {
				p.advance() // consume =
				if p.cur.Type == tokSCONST {
					val := "'" + p.cur.Str + "'"
					p.advance()
					opts = append(opts, &nodes.ServerConfigOption{Name: key, Value: val, Loc: nodes.Loc{Start: subLoc, End: p.prevEnd()}})
				} else if p.isAnyKeywordIdent() {
					opts = append(opts, &nodes.ServerConfigOption{Name: key, Value: p.cur.Str, Loc: nodes.Loc{Start: subLoc, End: p.prevEnd()}})
					p.advance()
				}
			} else {
				opts = append(opts, &nodes.ServerConfigOption{Name: key, Loc: nodes.Loc{Start: subLoc, End: p.prevEnd()}})
			}
		}
		if p.cur.Type == ')' {
			p.advance() // consume )
		}
	}

	return opts, nil
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
func (p *Parser) parseAlterServiceMasterKeyStmt() (*nodes.SecurityKeyStmt, error) {
	loc := p.pos()
	// SERVICE MASTER KEY already consumed by caller
	stmt := &nodes.SecurityKeyStmt{
		Action:     "ALTER",
		ObjectType: "SERVICE MASTER KEY",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	p.parseSecurityKeyOptions(stmt)

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseBackupServiceMasterKeyStmt parses BACKUP SERVICE MASTER KEY.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/backup-service-master-key-transact-sql
//
//	BACKUP SERVICE MASTER KEY TO FILE = 'path_to_file'
//	    ENCRYPTION BY PASSWORD = 'password'
func (p *Parser) parseBackupServiceMasterKeyStmt() (*nodes.SecurityKeyStmt, error) {
	loc := p.pos()
	// BACKUP already consumed, SERVICE MASTER KEY consumed by caller
	stmt := &nodes.SecurityKeyStmt{
		Action:     "BACKUP",
		ObjectType: "SERVICE MASTER KEY",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	p.parseSecurityKeyOptions(stmt)

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseRestoreServiceMasterKeyStmt parses RESTORE SERVICE MASTER KEY.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/restore-service-master-key-transact-sql
//
//	RESTORE SERVICE MASTER KEY FROM FILE = 'path_to_file'
//	    DECRYPTION BY PASSWORD = 'password' [FORCE]
func (p *Parser) parseRestoreServiceMasterKeyStmt() (*nodes.SecurityKeyStmt, error) {
	loc := p.pos()
	// RESTORE already consumed, SERVICE MASTER KEY consumed by caller
	stmt := &nodes.SecurityKeyStmt{
		Action:     "RESTORE",
		ObjectType: "SERVICE MASTER KEY",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	p.parseSecurityKeyOptions(stmt)

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}
