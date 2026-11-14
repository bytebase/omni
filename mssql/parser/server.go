// Package parser - server.go implements T-SQL server-level object statement parsing:
// CREATE/ALTER/DROP SERVER ROLE, ALTER SERVER CONFIGURATION,
// ALTER/BACKUP/RESTORE SERVICE MASTER KEY.
package parser

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseCreateServerRoleStmt parses CREATE SERVER ROLE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-server-role-transact-sql
// Grammar: TSql160.g:10300.
//
//	CREATE SERVER ROLE role_name [ AUTHORIZATION server_principal ]
func (p *Parser) parseCreateServerRoleStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "SERVER ROLE",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	name, err := p.prinName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	optLoc := p.pos()
	owner, ok, err := p.prinAuthorization()
	if err != nil {
		return nil, err
	}
	if ok {
		stmt.Options = &nodes.List{Items: []nodes.Node{&nodes.ServerConfigOption{Name: "AUTHORIZATION", Value: owner, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}}}}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterServerRoleStmt parses ALTER SERVER ROLE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-server-role-transact-sql
// Grammar: TSql160.g:12536, 12553-12590.
//
//	ALTER SERVER ROLE server_role_name
//	    { ADD MEMBER server_principal | DROP MEMBER server_principal | WITH NAME = new_server_role_name }
func (p *Parser) parseAlterServerRoleStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "ALTER",
		ObjectType: "SERVER ROLE",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	name, err := p.prinName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	optLoc := p.pos()
	action, value, err := p.prinRoleAction()
	if err != nil {
		return nil, err
	}
	stmt.Options = &nodes.List{Items: []nodes.Node{&nodes.ServerConfigOption{Name: action, Value: value, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}}}}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseDropServerRoleStmt parses DROP SERVER ROLE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-server-role-transact-sql
// Grammar: TSql160.g:14837.
//
//	DROP SERVER ROLE role_name
func (p *Parser) parseDropServerRoleStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "DROP",
		ObjectType: "SERVER ROLE",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	name, err := p.prinName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// SQL Server reads every word of ALTER SERVER CONFIGURATION by its text, so
// the sets below hold words, matched whether or not the lexer registers them
// as keywords. A bracketed or quoted word ([CPU], [LOCAL]) is refused, as SQL
// Server refuses it.
var (
	serverConfigOnOff   = newOptionSet().withIdents("ON", "OFF")
	serverConfigOn      = newOptionSet().withIdents("ON")
	serverConfigOff     = newOptionSet().withIdents("OFF")
	serverConfigDefault = newOptionSet().withIdents("DEFAULT")
	serverConfigAuto    = newOptionSet().withIdents("AUTO")
	serverConfigTo      = newOptionSet().withIdents("TO")
	serverConfigMB      = newOptionSet().withIdents("MB")
	serverConfigLocal   = newOptionSet().withIdents("LOCAL")

	// alterServerConfigurationStatement, TSql160.g:2247, plus the 2019 and
	// 2022 options ScriptDOM lacks (MEMORY_OPTIMIZED, HARDWARE_OFFLOAD,
	// SUSPEND_FOR_SNAPSHOT_BACKUP).
	serverConfigKinds = newOptionSet().withIdents(
		"PROCESS", "BUFFER", "DIAGNOSTICS", "FAILOVER", "HADR", "SOFTNUMA", "EXTERNAL",
		"MEMORY_OPTIMIZED", "HARDWARE_OFFLOAD", "SUSPEND_FOR_SNAPSHOT_BACKUP",
	)

	// alterServerConfigurationSetProcessAffinityStatement, TSql160.g:2617
	// (ProcessAffinityTypeHelper). Only CPU takes AUTO.
	serverConfigAffinityKinds = newOptionSet().withIdents("CPU", "NUMANODE")

	// alterServerConfigurationDiagnosticsLogOption, TSql160.g:2477.
	serverConfigDiagnosticsWords = newOptionSet().withIdents("PATH", "MAX_SIZE", "MAX_FILES")

	// alterServerConfigurationFailoverClusterPropertyOption, TSql160.g:2551
	// (AlterServerConfigurationFailoverClusterPropertyOptionHelper). SQL
	// Server 2022 refuses ClusterConnectionOptions.
	serverConfigFailoverProperties = newOptionSet().withIdents(
		"VerboseLogging", "SqlDumperDumpFlags", "SqlDumperDumpPath", "SqlDumperDumpTimeOut",
		"FailureConditionLevel", "HealthCheckTimeout",
	)

	// alterServerConfigurationHadrClusterOption, TSql160.g:2594.
	serverConfigContext = newOptionSet().withIdents("CONTEXT")

	// alterServerConfigurationBufferPoolExtensionSizeOption, TSql160.g:2444:
	// KB, MB or GB.
	serverConfigBPEUnits = newOptionSet().withIdents("KB", "MB", "GB")

	// alterServerConfigurationExternalAuthenticationOption, TSql160.g:2324.
	serverConfigUseIdentity    = newOptionSet().withIdents("USE_IDENTITY")
	serverConfigCredentialName = newOptionSet().withIdents("CREDENTIAL_NAME")

	// MEMORY_OPTIMIZED, HARDWARE_OFFLOAD and SUSPEND_FOR_SNAPSHOT_BACKUP words,
	// from the ALTER SERVER CONFIGURATION page.
	serverConfigMemoryOptimizedWords = newOptionSet().withIdents("TEMPDB_METADATA", "HYBRID_BUFFER_POOL")
	serverConfigResourcePool         = newOptionSet().withIdents("RESOURCE_POOL")
	serverConfigAccelerator          = newOptionSet().withIdents("ACCELERATOR")
	serverConfigQAT                  = newOptionSet().withIdents("QAT")
	serverConfigMode                 = newOptionSet().withIdents("MODE")
	serverConfigOffloadModes         = newOptionSet().withIdents("SOFTWARE", "HARDWARE")
	serverConfigSnapshotWords        = newOptionSet().withIdents("GROUP", "MODE")
	serverConfigCopyOnly             = newOptionSet().withIdents("COPY_ONLY")
)

// serverConfigMaxInt is the largest integer SQL Server reads in most
// ALTER SERVER CONFIGURATION options; a larger literal is a syntax error.
const serverConfigMaxInt = 2147483647

// parseAlterServerConfigurationStmt parses ALTER SERVER CONFIGURATION SET <optionspec>.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-server-configuration-transact-sql
// Grammar: alterServerConfigurationStatement, TSql160.g:2247.
//
//	ALTER SERVER CONFIGURATION
//	SET <optionspec>
//
//	<optionspec> ::=
//	{
//	     PROCESS AFFINITY { CPU = { AUTO | <CPU_range_spec> } | NUMANODE = <NUMA_node_range_spec> }
//	   | DIAGNOSTICS LOG { ON | OFF | PATH = { 'os_file_path' | DEFAULT } | MAX_SIZE = { log_max_size MB | DEFAULT } | MAX_FILES = { max_file_count | DEFAULT } }
//	   | FAILOVER CLUSTER PROPERTY <resource_property>
//	   | HADR CLUSTER CONTEXT = { 'remote_windows_cluster' | LOCAL }
//	   | BUFFER POOL EXTENSION { ON ( FILENAME = 'path', SIZE = size { KB | MB | GB } ) | OFF }
//	   | SOFTNUMA { ON | OFF }
//	   | MEMORY_OPTIMIZED { ON | OFF | TEMPDB_METADATA = { ON [ ( RESOURCE_POOL = 'pool' ) ] | OFF } | HYBRID_BUFFER_POOL = { ON | OFF } }
//	   | HARDWARE_OFFLOAD = { ON ( ACCELERATOR = QAT [ , MODE = { SOFTWARE | HARDWARE } ] ) | OFF ( ACCELERATOR = QAT ) }
//	   | SUSPEND_FOR_SNAPSHOT_BACKUP = { ON | OFF } [ ( { GROUP = ( database [ ,...n ] ) | MODE = COPY_ONLY } [ , ... ] ) ]
//	   | EXTERNAL AUTHENTICATION { ON ( { USE_IDENTITY | CREDENTIAL_NAME = 'name' } ) | OFF }
//	}
//
//	<CPU_range_spec> ::= { CPU_ID | CPU_ID TO CPU_ID } [ ,...n ]
//	<NUMA_node_range_spec> ::= { NUMA_node_ID | NUMA_node_ID TO NUMA_node_ID } [ ,...n ]
//
//	<resource_property> ::=
//	{
//	    VerboseLogging = { 0 | 1 | 2 | DEFAULT }
//	  | SqlDumperDumpFlags = { integer | DEFAULT }
//	  | SqlDumperDumpPath = { 'os_file_path' | DEFAULT }
//	  | SqlDumperDumpTimeOut = { integer | DEFAULT }
//	  | FailureConditionLevel = { 0 to 5 | DEFAULT }
//	  | HealthCheckTimeout = { 15000 or more | DEFAULT }
//	}
func (p *Parser) parseAlterServerConfigurationStmt() (*nodes.AlterServerConfigurationStmt, error) {
	loc := p.pos()
	// ALTER, SERVER, CONFIGURATION already consumed by caller
	if _, err := p.expect(kwSET); err != nil {
		return nil, err
	}

	stmt := &nodes.AlterServerConfigurationStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	kind, err := p.serverConfigExpect(serverConfigKinds)
	if err != nil {
		return nil, err
	}
	// Each kind's remaining words, written as the statement's OptionType.
	var words []string
	switch kind {
	case "PROCESS":
		words = []string{"AFFINITY"}
	case "DIAGNOSTICS":
		words = []string{"LOG"}
	case "FAILOVER":
		words = []string{"CLUSTER", "PROPERTY"}
	case "HADR":
		words = []string{"CLUSTER"}
	case "BUFFER":
		words = []string{"POOL", "EXTENSION"}
	case "EXTERNAL":
		words = []string{"AUTHENTICATION"}
	}
	optionType := kind
	for _, w := range words {
		if _, err := p.serverConfigExpect(newOptionSet().withIdents(w)); err != nil {
			return nil, err
		}
		optionType += " " + w
	}
	stmt.OptionType = optionType

	var opts []nodes.Node
	switch kind {
	case "PROCESS":
		opts, err = p.parseServerConfigProcessAffinity()
	case "DIAGNOSTICS":
		opts, err = p.parseServerConfigDiagnosticsLog()
	case "FAILOVER":
		opts, err = p.parseServerConfigFailoverClusterProperty()
	case "HADR":
		stmt.OptionType = "HADR CLUSTER CONTEXT"
		opts, err = p.parseServerConfigHadrCluster()
	case "BUFFER":
		opts, err = p.parseServerConfigBufferPoolExtension()
	case "SOFTNUMA":
		opts, err = p.parseServerConfigOnOff()
	case "MEMORY_OPTIMIZED":
		opts, err = p.parseServerConfigMemoryOptimized()
	case "HARDWARE_OFFLOAD":
		opts, err = p.parseServerConfigHardwareOffload()
	case "EXTERNAL":
		opts, err = p.parseServerConfigExternalAuthentication()
	case "SUSPEND_FOR_SNAPSHOT_BACKUP":
		opts, err = p.parseServerConfigSuspendForSnapshotBackup()
	}
	if err != nil {
		return nil, err
	}

	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// serverConfigAt reports whether the current token is an undelimited word of set.
func (p *Parser) serverConfigAt(set optionSet) bool {
	switch {
	case p.cur.Type == tokIDENT:
		if p.isDelimited(p.cur) {
			return false
		}
	case p.cur.Type < kwABSENT:
		return false
	}
	return p.isValidOption(set)
}

// serverConfigExpect consumes an undelimited word of set and returns it in upper case.
func (p *Parser) serverConfigExpect(set optionSet) (string, error) {
	if !p.serverConfigAt(set) {
		return "", p.syntaxErrorAtCur()
	}
	w := strings.ToUpper(p.cur.Str)
	p.advance()
	return w, nil
}

// serverConfigFlag consumes ON or OFF and returns it as a flag option.
func (p *Parser) serverConfigFlag(set optionSet) (*nodes.ServerConfigOption, error) {
	optLoc := p.pos()
	w, err := p.serverConfigExpect(set)
	if err != nil {
		return nil, err
	}
	return &nodes.ServerConfigOption{Name: w, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}}, nil
}

// serverConfigInteger reads an unsigned integer literal in [lo, hi]. SQL
// Server reads a literal above max as a syntax error and one outside [lo, hi]
// as out of range (Msg 1080).
func (p *Parser) serverConfigInteger(lo, hi, max int64) (string, error) {
	if p.cur.Type != tokICONST || p.cur.Ival > max {
		return "", p.syntaxErrorAtCur()
	}
	pos, v, text := p.pos(), p.cur.Ival, p.cur.Str
	p.advance()
	if v < lo || v > hi {
		return "", p.newParseError(pos, fmt.Sprintf("the integer value %d is out of range", v))
	}
	return text, nil
}

// serverConfigString reads a string literal and returns it quoted. SQL Server
// refuses ” where nonEmpty is set.
func (p *Parser) serverConfigString(nonEmpty bool) (string, error) {
	if (p.cur.Type != tokSCONST && p.cur.Type != tokNSCONST) || (nonEmpty && p.cur.Str == "") {
		return "", p.syntaxErrorAtCur()
	}
	s := "'" + p.cur.Str + "'"
	p.advance()
	return s, nil
}

// parseServerConfigProcessAffinity parses PROCESS AFFINITY options,
// TSql160.g:2617 and affinityRangeList, TSql160.g:2651.
//
//	PROCESS AFFINITY { CPU = { AUTO | <CPU_range_spec> } | NUMANODE = <NUMA_node_range_spec> }
//	<CPU_range_spec> ::= { CPU_ID | CPU_ID TO CPU_ID } [ ,...n ]
//	<NUMA_node_range_spec> ::= { NUMA_node_ID | NUMA_node_ID TO NUMA_node_ID } [ ,...n ]
func (p *Parser) parseServerConfigProcessAffinity() ([]nodes.Node, error) {
	optLoc := p.pos()
	key, err := p.serverConfigExpect(serverConfigAffinityKinds)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect('='); err != nil {
		return nil, err
	}

	if key == "CPU" && p.serverConfigAt(serverConfigAuto) {
		p.advance()
		return []nodes.Node{&nodes.ServerConfigOption{Name: key, Value: "AUTO", Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}}}, nil
	}

	var ranges []string
	for {
		from, err := p.serverConfigInteger(0, serverConfigMaxInt, serverConfigMaxInt)
		if err != nil {
			return nil, err
		}
		r := from
		if p.serverConfigAt(serverConfigTo) {
			p.advance()
			to, err := p.serverConfigInteger(0, serverConfigMaxInt, serverConfigMaxInt)
			if err != nil {
				return nil, err
			}
			r += " TO " + to
		}
		ranges = append(ranges, r)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return []nodes.Node{&nodes.ServerConfigOption{Name: key, Value: strings.Join(ranges, ", "), Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}}}, nil
}

// parseServerConfigDiagnosticsLog parses DIAGNOSTICS LOG options,
// alterServerConfigurationDiagnosticsLogOption TSql160.g:2477.
//
//	DIAGNOSTICS LOG { ON | OFF | PATH = { 'os_file_path' | DEFAULT }
//	    | MAX_SIZE = { log_max_size MB | DEFAULT } | MAX_FILES = { max_file_count | DEFAULT } }
//
// SQL Server refuses a size or file count of 0 (Msg 1080).
func (p *Parser) parseServerConfigDiagnosticsLog() ([]nodes.Node, error) {
	if p.serverConfigAt(serverConfigOnOff) {
		opt, err := p.serverConfigFlag(serverConfigOnOff)
		if err != nil {
			return nil, err
		}
		return []nodes.Node{opt}, nil
	}

	optLoc := p.pos()
	key, err := p.serverConfigExpect(serverConfigDiagnosticsWords)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect('='); err != nil {
		return nil, err
	}

	var val string
	switch {
	case p.serverConfigAt(serverConfigDefault):
		p.advance()
		val = "DEFAULT"
	case key == "PATH":
		val, err = p.serverConfigString(false)
	case key == "MAX_SIZE":
		val, err = p.serverConfigInteger(1, serverConfigMaxInt, serverConfigMaxInt)
		if err == nil {
			_, err = p.serverConfigExpect(serverConfigMB)
			val += " MB"
		}
	default: // MAX_FILES
		val, err = p.serverConfigInteger(1, serverConfigMaxInt, serverConfigMaxInt)
	}
	if err != nil {
		return nil, err
	}
	return []nodes.Node{&nodes.ServerConfigOption{Name: key, Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}}}, nil
}

// parseServerConfigFailoverClusterProperty parses FAILOVER CLUSTER PROPERTY
// <resource_property>, alterServerConfigurationFailoverClusterPropertyOption
// TSql160.g:2551. SqlDumperDumpPath takes a string and the others an
// integer; SQL Server checks each integer's range (Msg 1080).
func (p *Parser) parseServerConfigFailoverClusterProperty() ([]nodes.Node, error) {
	optLoc := p.pos()
	if !p.serverConfigAt(serverConfigFailoverProperties) {
		return nil, p.syntaxErrorAtCur()
	}
	key := p.cur.Str
	prop := strings.ToUpper(key)
	p.advance()
	if _, err := p.expect('='); err != nil {
		return nil, err
	}

	var val string
	var err error
	switch {
	case p.serverConfigAt(serverConfigDefault):
		p.advance()
		val = "DEFAULT"
	case prop == "SQLDUMPERDUMPPATH":
		val, err = p.serverConfigString(false)
	default:
		var lo, hi int64 = 0, 4294967295
		switch prop {
		case "VERBOSELOGGING":
			hi = 2
		case "FAILURECONDITIONLEVEL":
			hi = 5
		case "HEALTHCHECKTIMEOUT":
			lo = 15000
		}
		val, err = p.serverConfigInteger(lo, hi, 1<<62)
	}
	if err != nil {
		return nil, err
	}
	return []nodes.Node{&nodes.ServerConfigOption{Name: key, Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}}}, nil
}

// parseServerConfigHadrCluster parses CONTEXT = { 'remote_windows_cluster' | LOCAL }
// after HADR CLUSTER, alterServerConfigurationHadrClusterOption TSql160.g:2594.
func (p *Parser) parseServerConfigHadrCluster() ([]nodes.Node, error) {
	optLoc := p.pos()
	if _, err := p.serverConfigExpect(serverConfigContext); err != nil {
		return nil, err
	}
	if _, err := p.expect('='); err != nil {
		return nil, err
	}
	var val string
	if p.serverConfigAt(serverConfigLocal) {
		p.advance()
		val = "LOCAL"
	} else {
		var err error
		if val, err = p.serverConfigString(true); err != nil {
			return nil, err
		}
	}
	return []nodes.Node{&nodes.ServerConfigOption{Name: "CONTEXT", Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}}}, nil
}

// parseServerConfigBufferPoolExtension parses BUFFER POOL EXTENSION options,
// alterServerConfigurationBufferPoolExtensionContainerOption TSql160.g:2396.
//
//	BUFFER POOL EXTENSION { ON ( FILENAME = 'path' , SIZE = size { KB | MB | GB } ) | OFF }
func (p *Parser) parseServerConfigBufferPoolExtension() ([]nodes.Node, error) {
	if p.serverConfigAt(serverConfigOff) {
		opt, err := p.serverConfigFlag(serverConfigOff)
		if err != nil {
			return nil, err
		}
		return []nodes.Node{opt}, nil
	}
	on, err := p.serverConfigFlag(serverConfigOn)
	if err != nil {
		return nil, err
	}
	opts := []nodes.Node{on}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	fileLoc := p.pos()
	if _, err := p.expect(kwFILENAME); err != nil {
		return nil, err
	}
	if _, err := p.expect('='); err != nil {
		return nil, err
	}
	file, err := p.serverConfigString(false)
	if err != nil {
		return nil, err
	}
	opts = append(opts, &nodes.ServerConfigOption{Name: "FILENAME", Value: file, Loc: nodes.Loc{Start: fileLoc, End: p.prevEnd()}})
	if _, err := p.expect(','); err != nil {
		return nil, err
	}

	sizeLoc := p.pos()
	if _, err := p.expect(kwSIZE); err != nil {
		return nil, err
	}
	if _, err := p.expect('='); err != nil {
		return nil, err
	}
	size, err := p.serverConfigInteger(0, serverConfigMaxInt, serverConfigMaxInt)
	if err != nil {
		return nil, err
	}
	unit, err := p.serverConfigExpect(serverConfigBPEUnits)
	if err != nil {
		return nil, err
	}
	opts = append(opts, &nodes.ServerConfigOption{Name: "SIZE", Value: size + " " + unit, Loc: nodes.Loc{Start: sizeLoc, End: p.prevEnd()}})
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return opts, nil
}

// parseServerConfigOnOff parses SOFTNUMA { ON | OFF },
// alterServerConfigurationSoftNumaOption TSql160.g:2368.
func (p *Parser) parseServerConfigOnOff() ([]nodes.Node, error) {
	opt, err := p.serverConfigFlag(serverConfigOnOff)
	if err != nil {
		return nil, err
	}
	return []nodes.Node{opt}, nil
}

// parseServerConfigMemoryOptimized parses MEMORY_OPTIMIZED options.
//
//	MEMORY_OPTIMIZED { ON | OFF | TEMPDB_METADATA = { ON [ ( RESOURCE_POOL = 'pool' ) ] | OFF }
//	    | HYBRID_BUFFER_POOL = { ON | OFF } }
func (p *Parser) parseServerConfigMemoryOptimized() ([]nodes.Node, error) {
	if p.serverConfigAt(serverConfigOnOff) {
		opt, err := p.serverConfigFlag(serverConfigOnOff)
		if err != nil {
			return nil, err
		}
		return []nodes.Node{opt}, nil
	}

	optLoc := p.pos()
	key, err := p.serverConfigExpect(serverConfigMemoryOptimizedWords)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect('='); err != nil {
		return nil, err
	}
	val, err := p.serverConfigExpect(serverConfigOnOff)
	if err != nil {
		return nil, err
	}
	opts := []nodes.Node{&nodes.ServerConfigOption{Name: key, Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}}}

	// ( RESOURCE_POOL = 'pool' ) follows TEMPDB_METADATA = ON only.
	if key == "TEMPDB_METADATA" && val == "ON" && p.cur.Type == '(' {
		p.advance()
		subLoc := p.pos()
		if _, err := p.serverConfigExpect(serverConfigResourcePool); err != nil {
			return nil, err
		}
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		pool, err := p.serverConfigString(true)
		if err != nil {
			return nil, err
		}
		opts = append(opts, &nodes.ServerConfigOption{Name: "RESOURCE_POOL", Value: pool, Loc: nodes.Loc{Start: subLoc, End: p.prevEnd()}})
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}
	return opts, nil
}

// parseServerConfigHardwareOffload parses HARDWARE_OFFLOAD options.
//
//	HARDWARE_OFFLOAD = { ON ( ACCELERATOR = QAT [ , MODE = { SOFTWARE | HARDWARE } ] )
//	    | OFF ( ACCELERATOR = QAT ) }
//
// SQL Server 2022 refuses HARDWARE_OFFLOAD ON and = ON without the accelerator.
// The Linux server refuses MODE = HARDWARE, which the page documents for
// Windows; Parse takes it.
func (p *Parser) parseServerConfigHardwareOffload() ([]nodes.Node, error) {
	if _, err := p.expect('='); err != nil {
		return nil, err
	}
	flag, err := p.serverConfigFlag(serverConfigOnOff)
	if err != nil {
		return nil, err
	}
	opts := []nodes.Node{flag}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	accLoc := p.pos()
	if _, err := p.serverConfigExpect(serverConfigAccelerator); err != nil {
		return nil, err
	}
	if _, err := p.expect('='); err != nil {
		return nil, err
	}
	acc, err := p.serverConfigExpect(serverConfigQAT)
	if err != nil {
		return nil, err
	}
	opts = append(opts, &nodes.ServerConfigOption{Name: "ACCELERATOR", Value: acc, Loc: nodes.Loc{Start: accLoc, End: p.prevEnd()}})
	if flag.Name == "ON" && p.cur.Type == ',' {
		p.advance()
		modeLoc := p.pos()
		if _, err := p.serverConfigExpect(serverConfigMode); err != nil {
			return nil, err
		}
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		mode, err := p.serverConfigExpect(serverConfigOffloadModes)
		if err != nil {
			return nil, err
		}
		opts = append(opts, &nodes.ServerConfigOption{Name: "MODE", Value: mode, Loc: nodes.Loc{Start: modeLoc, End: p.prevEnd()}})
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return opts, nil
}

// parseServerConfigSuspendForSnapshotBackup parses SUSPEND_FOR_SNAPSHOT_BACKUP options.
//
//	SUSPEND_FOR_SNAPSHOT_BACKUP = { ON | OFF }
//	    [ ( { GROUP = ( database [ ,...n ] ) | MODE = COPY_ONLY } [ , ... ] ) ]
//
// SQL Server takes GROUP and MODE in either order, each once (Msg 1039).
func (p *Parser) parseServerConfigSuspendForSnapshotBackup() ([]nodes.Node, error) {
	if _, err := p.expect('='); err != nil {
		return nil, err
	}
	flag, err := p.serverConfigFlag(serverConfigOnOff)
	if err != nil {
		return nil, err
	}
	opts := []nodes.Node{flag}
	if p.cur.Type != '(' {
		return opts, nil
	}
	p.advance()
	seen := map[string]bool{}
	for {
		subLoc := p.pos()
		key, err := p.serverConfigExpect(serverConfigSnapshotWords)
		if err != nil {
			return nil, err
		}
		if seen[key] {
			return nil, p.newParseError(subLoc, fmt.Sprintf("option '%s' is specified more than once", key))
		}
		seen[key] = true
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		var val string
		if key == "GROUP" {
			if _, err := p.expect('('); err != nil {
				return nil, err
			}
			var dbs []string
			for {
				name, ok := p.parseIdentifier()
				if !ok {
					return nil, p.syntaxErrorAtCur()
				}
				dbs = append(dbs, name)
				if p.cur.Type != ',' {
					break
				}
				p.advance()
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			val = strings.Join(dbs, ", ")
		} else if val, err = p.serverConfigExpect(serverConfigCopyOnly); err != nil {
			return nil, err
		}
		opts = append(opts, &nodes.ServerConfigOption{Name: key, Value: val, Loc: nodes.Loc{Start: subLoc, End: p.prevEnd()}})
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return opts, nil
}

// parseServerConfigExternalAuthentication parses EXTERNAL AUTHENTICATION options,
// alterServerConfigurationExternalAuthenticationContainerOption TSql160.g:2293:
// ON takes exactly one option, OFF none.
//
//	EXTERNAL AUTHENTICATION { ON ( { USE_IDENTITY | CREDENTIAL_NAME = 'name' } ) | OFF }
func (p *Parser) parseServerConfigExternalAuthentication() ([]nodes.Node, error) {
	if p.serverConfigAt(serverConfigOff) {
		opt, err := p.serverConfigFlag(serverConfigOff)
		if err != nil {
			return nil, err
		}
		return []nodes.Node{opt}, nil
	}
	on, err := p.serverConfigFlag(serverConfigOn)
	if err != nil {
		return nil, err
	}
	opts := []nodes.Node{on}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	subLoc := p.pos()
	switch {
	case p.serverConfigAt(serverConfigUseIdentity):
		key := p.cur.Str
		p.advance()
		opts = append(opts, &nodes.ServerConfigOption{Name: key, Loc: nodes.Loc{Start: subLoc, End: p.prevEnd()}})
	case p.serverConfigAt(serverConfigCredentialName):
		key := p.cur.Str
		p.advance()
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		val, err := p.serverConfigString(false)
		if err != nil {
			return nil, err
		}
		opts = append(opts, &nodes.ServerConfigOption{Name: key, Value: val, Loc: nodes.Loc{Start: subLoc, End: p.prevEnd()}})
	default:
		return nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
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

	if err := p.keysServiceMasterKeyBody("ALTER", stmt); err != nil {
		return nil, err
	}

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

	if err := p.keysServiceMasterKeyBody("BACKUP", stmt); err != nil {
		return nil, err
	}

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

	if err := p.keysServiceMasterKeyBody("RESTORE", stmt); err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}
